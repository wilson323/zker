/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package billing

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
	"github.com/coze-dev/coze-studio/backend/tests/integration"
)

// ============================================================
// 端到端集成测试套件
//
// 测试完整的计费流程，验证所有组件协同工作：
// 1. Token记录 → 成本计算 → 汇总更新
// 2. 预算配置 → 使用统计 → 告警触发 → 通知发送
// 3. 数据一致性验证
// ============================================================

// ============================================================
// 测试辅助函数
// ============================================================

/**
 * setupE2ETestContext 设置端到端测试上下文
 */
type e2eTestContext struct {
	DB                  *gorm.DB
	TokenMeteringSvc    *service.TokenMeteringService
	BudgetManagementSvc *service.BudgetManagementService
	BudgetAlertSvc      *service.BudgetAlertService
	Logger              *logrus.Logger
}

func setupE2ETestContext(t *testing.T) (*e2eTestContext, func()) {
	// 1. 创建MySQL测试容器
	testCtx, err := integration.SetupMySQLTestContainer(t)
	require.NoError(t, err, "Failed to setup MySQL test container")

	// 2. 使用GORM连接数据库
	db, err := gorm.Open(mysql.Open(testCtx.Config.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	require.NoError(t, err, "Failed to open GORM connection")

	// 3. 自动迁移表结构
	err = db.AutoMigrate(
		&entity.TokenUsageLog{},
		&entity.TokenUsageSummary{},
		&entity.BudgetSettings{},
		&entity.BudgetAlert{},
		&entity.CostOptimizationSuggestion{},
	)
	require.NoError(t, err, "Failed to migrate tables")

	// 4. 创建Repository实例
	usageLogRepo := repository.NewTokenUsageLogRepository(db)
	usageSummaryRepo := repository.NewTokenUsageSummaryRepository(db)
	budgetRepo := repository.NewBudgetSettingsRepository(db)
	alertRepo := repository.NewBudgetAlertRepository(db)

	// 5. 创建Service实例
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)
	logger.SetLevel(logrus.DebugLevel)

	pricingEngine := service.NewPricingEngine(logger)
	notifier := service.NewNotificationService(logger)

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		usageSummaryRepo,
		alertRepo,
		notifier,
		logger,
	)

	tokenMeteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		db,
		logger,
	)

	budgetManagementSvc := service.NewBudgetManagementService(
		budgetRepo,
		alertRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		logger,
	)

	// 6. 返回测试上下文
	ctx := &e2eTestContext{
		DB:                  db,
		TokenMeteringSvc:    tokenMeteringSvc,
		BudgetManagementSvc: budgetManagementSvc,
		BudgetAlertSvc:      budgetAlertSvc,
		Logger:              logger,
	}

	// 7. 清理函数
	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		if testCtx.Cleanup != nil {
			testCtx.Cleanup()
		}
	}

	return ctx, cleanup
}

/**
 * createTestTenantAndBot 创建测试租户和Bot
 */
func createTestTenantAndBot(ctx *e2eTestContext, tenantID, botID string) {
	tenant := &struct {
		TenantID string `gorm:"primaryKey;size:36"`
		Name     string `gorm:"size:100"`
	}{
		TenantID: tenantID,
		Name:     "测试租户",
	}
	ctx.DB.Table("tenants").Create(tenant)

	bot := &struct {
		BotID    string `gorm:"primaryKey;size:36"`
		TenantID string `gorm:"size:36;not null"`
		Name     string `gorm:"size:100"`
	}{
		BotID:    botID,
		TenantID: tenantID,
		Name:     "测试Bot",
	}
	ctx.DB.Table("bots").Create(bot)
}

/**
 * waitForAsyncOperations 等待异步操作完成
 */
func waitForAsyncOperations(t *testing.T, ctx *e2eTestContext, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// 检查是否有汇总数据
		var summaryCount int64
		ctx.DB.Model(&entity.TokenUsageSummary{}).Count(&summaryCount)

		// 检查是否有告警数据
		var alertCount int64
		ctx.DB.Model(&entity.BudgetAlert{}).Count(&alertCount)

		if summaryCount > 0 && alertCount > 0 {
			return // 成功
		}

		time.Sleep(100 * time.Millisecond)
	}

	t.Logf("⚠️  Async operations not completed within %v (continuing anyway)", timeout)
}

// ============================================================
// 测试场景: 完整计费流程端到端测试
// ============================================================

/**
 * TestBilling_E2E_Flow
 *
 * 测试完整的计费流程：
 * 1. 创建租户和Bot
 * 2. 配置月度预算
 * 3. 模拟AI模型调用（记录Token使用）
 * 4. 验证成本计算
 * 5. 验证汇总更新
 * 6. 验证预算告警（达到阈值）
 * 7. 验证通知发送（模拟）
 * 8. 查询使用统计
 * 9. 查询预算使用情况
 * 10. 查询告警历史
 */
func TestBilling_E2E_Flow(t *testing.T) {
	ctx, cleanup := setupE2ETestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_e2e_test"
	botID := "bot_e2e_test"

	t.Log("🚀 开始端到端测试: 完整计费流程")

	// ========== 步骤1: 创建租户和Bot ==========
	t.Log("📝 步骤1: 创建租户和Bot")
	createTestTenantAndBot(ctx, tenantID, botID)

	// 验证租户创建成功
	var tenantCount int64
	ctx.DB.Table("tenants").Where("tenant_id = ?", tenantID).Count(&tenantCount)
	assert.Equal(t, int64(1), tenantCount, "Tenant should be created")

	// 验证Bot创建成功
	var botCount int64
	ctx.DB.Table("bots").Where("bot_id = ? AND tenant_id = ?", botID, tenantID).Count(&botCount)
	assert.Equal(t, int64(1), botCount, "Bot should be created")

	t.Log("✅ 租户和Bot创建成功")

	// ========== 步骤2: 配置月度预算 ==========
	t.Log("📝 步骤2: 配置月度预算")

	hardCapAmount := 1200.00
	originalModel := "openai/gpt-4"
	fallbackModel := "openai/gpt-3.5-turbo"

	createBudgetReq := &service.CreateBudgetRequest{
		TenantID:     tenantID,
		BudgetType:   entity.BudgetTypeMonthly,
		BudgetAmount: 1000.00,
		Currency:     "CNY",
		AlertThreshold1: 80,
		AlertThreshold2: 95,
		HardCapEnabled:   true,
		HardCapAmount:    &hardCapAmount,
		AutoDowngradeEnabled: true,
		DowngradeThreshold:   90,
		OriginalModel:        &originalModel,
		FallbackModel:        &fallbackModel,
		NotificationChannels: []string{"email", "webhook"},
		NotificationRecipients: []string{"admin@example.com", "ops@example.com"},
	}

	budget, err := ctx.BudgetManagementSvc.CreateBudget(testCtx, createBudgetReq)
	require.NoError(t, err, "Failed to create budget")

	assert.Equal(t, tenantID, budget.TenantID)
	assert.Equal(t, 1000.00, budget.BudgetAmount)
	assert.Equal(t, 80, budget.AlertThreshold1)
	assert.Equal(t, 95, budget.AlertThreshold2)

	t.Log("✅ 预算配置成功")
	t.Logf("   - 预算金额: ¥%.2f", budget.BudgetAmount)
	t.Logf("   - 一级告警阈值: %d%%", budget.AlertThreshold1)
	t.Logf("   - 二级告警阈值: %d%%", budget.AlertThreshold2)

	// ========== 步骤3: 模拟AI模型调用（记录Token使用） ==========
	t.Log("📝 步骤3: 模拟AI模型调用")

	// 模拟多次AI模型调用，累计使用约850元（超过80%阈值）
	// 每次调用约10元，共85次
	var totalRequests = 85
	var totalCost float64

	for i := 0; i < totalRequests; i++ {
		var userID uint64 = 1001 + uint64(i)
		convID := stringPtr("conv_e2e_001")
		msgID := stringPtr("msg_e2e_001")
		responseTimeMs := 1200 + i*10

		req := &service.RecordTokenUsageRequest{
			TenantID:       tenantID,
			UserID:         &userID,
			BotID:          &botID,
			ConversationID: convID,
			MessageID:      msgID,
			ModelProvider:  "openai",
			ModelName:      "gpt-4",
			InputTokens:    100000, // 约5元
			OutputTokens:   100000, // 约5元
			TotalTokens:    200000,
			RequestType:    entity.RequestTypeChat,
			ResponseTimeMs: &responseTimeMs,
			IsCached:       i%10 == 0, // 10%缓存命中率
		}

		resp, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err, "RecordTokenUsage should succeed")

		totalCost += resp.TotalCost

		// 每20次记录一次进度
		if (i+1)%20 == 0 {
			t.Logf("   已记录 %d/%d 次调用，累计成本: ¥%.2f", i+1, totalRequests, totalCost)
		}
	}

	t.Log("✅ AI模型调用记录完成")
	t.Logf("   - 总请求数: %d", totalRequests)
	t.Logf("   - 累计成本: ¥%.2f", totalCost)

	// ========== 步骤4: 验证成本计算 ==========
	t.Log("📝 步骤4: 验证成本计算")

	// 从TokenUsageLog表验证成本计算
	var logs []entity.TokenUsageLog
	err = ctx.DB.Where("tenant_id = ?", tenantID).
		Find(&logs).Error
	require.NoError(t, err)

	var calculatedTotalCost float64
	for _, log := range logs {
		calculatedTotalCost += log.TotalCost
	}

	assert.Greater(t, calculatedTotalCost, 800.0, "Total cost should exceed 800")
	assert.InDelta(t, totalCost, calculatedTotalCost, 1.0, "Cost calculation should be accurate")

	t.Log("✅ 成本计算验证通过")
	t.Logf("   - 日志记录数: %d", len(logs))
	t.Logf("   - 总成本: ¥%.2f", calculatedTotalCost)

	// ========== 步骤5: 验证汇总更新 ==========
	t.Log("📝 步骤5: 验证汇总更新")

	// 等待异步汇总更新
	time.Sleep(3 * time.Second)

	var summaries []entity.TokenUsageSummary
	err = ctx.DB.Where("tenant_id = ?", tenantID).
		Find(&summaries).Error
	require.NoError(t, err)

	assert.Greater(t, len(summaries), 0, "Should have summary records")

	// 验证汇总数据
	var totalInputTokens, totalOutputTokens, totalTokens int64
	var totalRequests int
	var totalCostFromSummary float64

	for _, summary := range summaries {
		totalInputTokens += summary.TotalInputTokens
		totalOutputTokens += summary.TotalOutputTokens
		totalTokens += summary.TotalTokens
		totalRequests += summary.TotalRequests
		totalCostFromSummary += summary.TotalCost
	}

	assert.Equal(t, int64(totalRequests), totalTokens/200000, "Total tokens should match")
	assert.Greater(t, totalCostFromSummary, 0.0, "Total cost from summary should be positive")

	t.Log("✅ 汇总更新验证通过")
	t.Logf("   - 汇总记录数: %d", len(summaries))
	t.Logf("   - 总Token数: %d", totalTokens)
	t.Logf("   - 总成本: ¥%.2f", totalCostFromSummary)

	// ========== 步骤6: 验证预算告警（达到阈值） ==========
	t.Log("📝 步骤6: 验证预算告警")

	// 手动触发预算检查
	err = ctx.BudgetAlertSvc.CheckBudget(testCtx, tenantID)
	require.NoError(t, err, "CheckBudget should succeed")

	// 等待异步告警处理
	time.Sleep(1 * time.Second)

	// 验证告警记录
	var alerts []entity.BudgetAlert
	err = ctx.DB.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&alerts).Error
	require.NoError(t, err)

	assert.Greater(t, len(alerts), 0, "Should have alert records")

	// 验证告警级别和类型
	hasThreshold1Alert := false
	hasThreshold2Alert := false

	for _, alert := range alerts {
		if alert.AlertType == entity.AlertTypeThreshold1 {
			hasThreshold1Alert = true
			assert.Equal(t, entity.AlertLevelWarning, alert.AlertLevel)
			assert.GreaterOrEqual(t, alert.UsagePercent, 80.0)
		}
		if alert.AlertType == entity.AlertTypeThreshold2 {
			hasThreshold2Alert = true
			assert.Equal(t, entity.AlertLevelCritical, alert.AlertLevel)
			assert.GreaterOrEqual(t, alert.UsagePercent, 95.0)
		}
	}

	assert.True(t, hasThreshold1Alert, "Should have threshold1 alert")
	// hasThreshold2Alert可能为false，因为使用量在80-95%之间

	t.Log("✅ 预算告警验证通过")
	t.Logf("   - 告警记录数: %d", len(alerts))
	for _, alert := range alerts {
		t.Logf("   - 告警类型: %s, 级别: %s, 使用率: %.1f%%",
			alert.AlertType, alert.AlertLevel, alert.UsagePercent)
	}

	// ========== 步骤7: 验证通知发送（模拟） ==========
	t.Log("📝 步骤7: 验证通知发送")

	// NotificationService会在日志中记录发送的通知
	// 这里我们验证通知状态已更新
	var sentAlerts int64
	ctx.DB.Model(&entity.BudgetAlert{}).
		Where("tenant_id = ? AND notification_status = ?", tenantID, entity.NotificationStatusSent).
		Count(&sentAlerts)

	assert.Greater(t, sentAlerts, int64(0), "Should have sent notifications")

	t.Log("✅ 通知发送验证通过")
	t.Logf("   - 已发送通知数: %d", sentAlerts)

	// ========== 步骤8: 查询使用统计 ==========
	t.Log("📝 步骤8: 查询使用统计")

	now := time.Now()
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	filter := &service.UsageStatsFilter{
		StartDate: startDate,
		EndDate:   endDate,
		BotID:     &botID,
	}

	stats, err := ctx.TokenMeteringSvc.GetUsageStats(testCtx, tenantID, filter)
	require.NoError(t, err, "GetUsageStats should succeed")

	assert.Equal(t, tenantID, stats.TenantID)
	assert.Greater(t, stats.TotalTokens, int64(0), "Should have positive total tokens")
	assert.Greater(t, stats.TotalCost, 0.0, "Should have positive total cost")
	assert.Greater(t, stats.TotalRequests, int64(0), "Should have positive total requests")

	t.Log("✅ 使用统计查询成功")
	t.Logf("   - 总Token数: %d", stats.TotalTokens)
	t.Logf("   - 总成本: ¥%.2f", stats.TotalCost)
	t.Logf("   - 总请求数: %d", stats.TotalRequests)
	t.Logf("   - 缓存命中: %d", stats.CachedRequests)
	t.Logf("   - 平均响应时间: %.2fms", stats.AvgResponseTime)

	// ========== 步骤9: 查询预算使用情况 ==========
	t.Log("📝 步骤9: 查询预算使用情况")

	usage, err := ctx.BudgetManagementSvc.GetBudgetUsage(testCtx, tenantID)
	require.NoError(t, err, "GetBudgetUsage should succeed")

	assert.Equal(t, 1000.00, usage.BudgetAmount)
	assert.Greater(t, usage.UsedAmount, 800.0, "Used amount should exceed 800")
	assert.Greater(t, usage.UsagePercent, 80.0, "Usage percent should exceed 80%")
	assert.Less(t, usage.RemainingAmount, 200.0, "Remaining amount should be less than 200")
	assert.Greater(t, usage.TotalTokens, int64(0), "Should have positive total tokens")
	assert.Greater(t, usage.TotalRequests, 0, "Should have positive total requests")
	assert.Greater(t, usage.AverageCost, 0.0, "Should have positive average cost")
	assert.Greater(t, usage.AlertCount, 0, "Should have alerts")
	assert.NotNil(t, usage.LastAlertAt, "Should have last alert time")
	assert.NotNil(t, usage.LastAlertLevel, "Should have last alert level")
	assert.NotNil(t, usage.EstimatedDailyUsage, "Should have estimated daily usage")
	assert.Greater(t, usage.DaysRemaining, 0, "Should have days remaining")

	t.Log("✅ 预算使用情况查询成功")
	t.Logf("   - 预算金额: ¥%.2f", usage.BudgetAmount)
	t.Logf("   - 已使用: ¥%.2f (%.1f%%)", usage.UsedAmount, usage.UsagePercent)
	t.Logf("   - 剩余: ¥%.2f", usage.RemainingAmount)
	t.Logf("   - 总Token数: %d", usage.TotalTokens)
	t.Logf("   - 总请求数: %d", usage.TotalRequests)
	t.Logf("   - 平均成本: ¥%.4f/1K tokens", usage.AverageCost)
	t.Logf("   - 告警数: %d", usage.AlertCount)
	t.Logf("   - 最后告警级别: %s", stringPtrToString(usage.LastAlertLevel))
	t.Logf("   - 预计每日使用: ¥%.2f", usage.EstimatedDailyUsage)
	t.Logf("   - 剩余天数: %d", usage.DaysRemaining)
	t.Logf("   - 是否会超预算: %v", usage.WillExceedBudget)

	// ========== 步骤10: 查询告警历史 ==========
	t.Log("📝 步骤10: 查询告警历史")

	alertFilter := &service.AlertFilter{
		PageSize: 10,
	}

	alertHistory, pagination, err := ctx.BudgetManagementSvc.GetAlerts(testCtx, tenantID, alertFilter)
	require.NoError(t, err, "GetAlerts should succeed")

	assert.Greater(t, len(alertHistory), 0, "Should have alert history")
	assert.Equal(t, pagination.TotalCount, int64(len(alertHistory)), "Total count should match")

	for i, alert := range alertHistory {
		assert.Equal(t, tenantID, alert.TenantID)
		assert.NotEmpty(t, alert.AlertType)
		assert.NotEmpty(t, alert.AlertLevel)
		assert.Greater(t, alert.UsagePercent, 0.0)
		assert.NotEmpty(t, alert.NotificationStatus)

		if i < 3 { // 只打印前3条
			t.Logf("   - 告警%d: %s/%s, 使用率%.1f%%, 状态:%s",
				i+1, alert.AlertType, alert.AlertLevel, alert.UsagePercent, alert.NotificationStatus)
		}
	}

	t.Log("✅ 告警历史查询成功")
	t.Logf("   - 告警记录数: %d", len(alertHistory))

	// ========== 数据一致性验证 ==========
	t.Log("📝 数据一致性验证")

	// 验证TokenUsageLog、TokenUsageSummary、BudgetSettings的数据一致性
	var logCount int64
	ctx.DB.Model(&entity.TokenUsageLog{}).Where("tenant_id = ?", tenantID).Count(&logCount)

	var summaryCount int64
	ctx.DB.Model(&entity.TokenUsageSummary{}).Where("tenant_id = ?", tenantID).Count(&summaryCount)

	var budgetCount int64
	ctx.DB.Model(&entity.BudgetSettings{}).Where("tenant_id = ?", tenantID).Count(&budgetCount)

	assert.Equal(t, int64(totalRequests), logCount, "Log count should match requests")
	assert.Greater(t, summaryCount, int64(0), "Should have summary records")
	assert.Equal(t, int64(1), budgetCount, "Should have one budget record")

	t.Log("✅ 数据一致性验证通过")
	t.Logf("   - 日志记录: %d", logCount)
	t.Logf("   - 汇总记录: %d", summaryCount)
	t.Logf("   - 预算记录: %d", budgetCount)

	// ========== 测试总结 ==========
	t.Log("🎉 端到端测试完成!")
	t.Log("==========================================")
	t.Log("✅ 所有步骤验证通过:")
	t.Log("  1. ✅ 租户和Bot创建")
	t.Log("  2. ✅ 预算配置")
	t.Log("  3. ✅ AI模型调用记录")
	t.Log("  4. ✅ 成本计算")
	t.Log("  5. ✅ 汇总更新")
	t.Log("  6. ✅ 预算告警")
	t.Log("  7. ✅ 通知发送")
	t.Log("  8. ✅ 使用统计")
	t.Log("  9. ✅ 预算使用情况")
	t.Log(" 10. ✅ 告警历史")
	t.Log(" 11. ✅ 数据一致性")
	t.Log("==========================================")
}

// ============================================================
// 辅助函数
// ============================================================

/**
 * stringPtr 返回字符串指针
 */
func stringPtr(s string) *string {
	return &s
}

/**
 * stringPtrToString 安全地将字符串指针转换为字符串
 */
func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
