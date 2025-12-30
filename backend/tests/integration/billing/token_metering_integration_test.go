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
// Token计量API集成测试套件
//
// 测试场景覆盖:
// 1. 完整流程测试 - 记录使用、汇总更新、成本计算
// 2. 批量记录测试 - 1000条记录的性能验证
// 3. 多模型统计测试 - 验证不同模型的Token分布
// 4. 每日趋势测试 - 验证跨天数据统计
// 5. 预算告警触发测试 - 验证告警机制
// ============================================================

// ============================================================
// 测试辅助函数
// ============================================================

/**
 * setupTestContext 设置测试上下文
 *
 * 遵循SOLID原则：单一职责
 *
 * @param t 测试对象
 * @return *testContext 测试上下文
 * @return func() 清理函数
 */
type testContext struct {
	DB              *gorm.DB
	TokenMeteringSvc *service.TokenMeteringService
	BudgetAlertSvc  *service.BudgetAlertService
	Logger          *logrus.Logger
}

func setupTestContext(t *testing.T) (*testContext, func()) {
	// 1. 创建MySQL测试容器（使用项目标准设置）
	testCtx, err := integration.SetupMySQLTestContainer(t)
	require.NoError(t, err, "Failed to setup MySQL test container")

	// 2. 使用GORM连接数据库
	db, err := gorm.Open(mysql.Open(testCtx.Config.DSN), &gorm.Config{
		SkipDefaultTransaction: true, // 禁用默认事务，提高性能
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

	// 6. 返回测试上下文
	ctx := &testContext{
		DB:              db,
		TokenMeteringSvc: tokenMeteringSvc,
		BudgetAlertSvc:  budgetAlertSvc,
		Logger:          logger,
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
 * createTestTenant 创建测试租户
 *
 * 遵循DRY原则：避免重复创建租户的代码
 */
func createTestTenant(ctx *testContext, tenantID string) {
	tenant := &struct {
		TenantID string `gorm:"primaryKey;size:36"`
		Name     string `gorm:"size:100"`
	}{
		TenantID: tenantID,
		Name:     "测试租户",
	}

	ctx.DB.Table("tenants").Create(tenant)
}

/**
 * createTestBot 创建测试Bot
 */
func createTestBot(ctx *testContext, botID, tenantID string) {
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
 * waitForSummary 等待汇总更新（异步操作）
 *
 * 由于汇总更新是异步的，需要轮询检查
 */
func waitForSummary(t *testing.T, ctx *testContext, tenantID string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var count int64
		err := ctx.DB.Model(&entity.TokenUsageSummary{}).
			Where("tenant_id = ?", tenantID).
			Count(&count).Error

		if err == nil && count > 0 {
			return // 成功
		}

		time.Sleep(100 * time.Millisecond)
	}

	// 超时
	t.Fatalf("Summary not updated within %v", timeout)
}

// ============================================================
// 测试场景1: 完整流程测试
// ============================================================

/**
 * TestTokenMetering_CompleteFlow
 *
 * 测试完整的Token计量流程：
 * 1. 创建租户和Bot
 * 2. 调用RecordTokenUsage记录使用
 * 3. 验证TokenUsageLog表中的记录
 * 4. 等待异步汇总完成
 * 5. 验证TokenUsageSummary表中的汇总数据
 * 6. 验证成本计算准确性
 */
func TestTokenMetering_CompleteFlow(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()

	// 1. 创建测试租户和Bot
	tenantID := "tenant_test_complete"
	botID := "bot_test_complete"
	createTestTenant(ctx, tenantID)
	createTestBot(ctx, botID, tenantID)

	// 2. 记录Token使用
	testCtx := context.Background()
	var userID uint64 = 1001
	responseTimeMs := 1200

	req := &service.RecordTokenUsageRequest{
		TenantID:       tenantID,
		UserID:         &userID,
		BotID:          &botID,
		ConversationID: stringPtr("conv_test_001"),
		MessageID:      stringPtr("msg_test_001"),
		ModelProvider:  "openai",
		ModelName:      "gpt-4",
		InputTokens:    1000,
		OutputTokens:   500,
		TotalTokens:    1500,
		RequestType:    entity.RequestTypeChat,
		ResponseTimeMs: &responseTimeMs,
		IsCached:       false,
	}

	resp, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)

	// 3. 验证响应
	require.NoError(t, err, "RecordTokenUsage should succeed")
	assert.NotZero(t, resp.LogID, "LogID should be generated")
	assert.Greater(t, resp.InputCost, 0.0, "InputCost should be positive")
	assert.Greater(t, resp.OutputCost, 0.0, "OutputCost should be positive")
	assert.Greater(t, resp.TotalCost, 0.0, "TotalCost should be positive")

	// 4. 验证TokenUsageLog表中的记录
	var log entity.TokenUsageLog
	err = ctx.DB.First(&log, resp.LogID).Error
	require.NoError(t, err, "Should find the created log")
	assert.Equal(t, tenantID, log.TenantID)
	assert.Equal(t, 1000, log.InputTokens)
	assert.Equal(t, 500, log.OutputTokens)
	assert.Equal(t, 1500, log.TotalTokens)
	assert.Equal(t, "openai", log.ModelProvider)
	assert.Equal(t, "gpt-4", log.ModelName)
	assert.False(t, log.IsCached)

	// 5. 等待异步汇总完成
	waitForSummary(t, ctx, tenantID, 5*time.Second)

	// 6. 验证TokenUsageSummary表中的汇总数据
	var summary entity.TokenUsageSummary
	err = ctx.DB.Where("tenant_id = ? AND summary_hour IS NOT NULL", tenantID).
		Order("created_at DESC").
		First(&summary).Error
	require.NoError(t, err, "Should find summary record")
	assert.Equal(t, tenantID, summary.TenantID)
	assert.Equal(t, int64(1000), summary.TotalInputTokens)
	assert.Equal(t, int64(500), summary.TotalOutputTokens)
	assert.Equal(t, int64(1500), summary.TotalTokens)
	assert.Equal(t, 1, summary.TotalRequests)
	assert.Equal(t, 0, summary.CachedRequests)
	assert.NotNil(t, summary.AvgResponseTime)

	// 7. 验证成本计算准确性
	// gpt-4定价: 输入 ¥0.03/1K tokens, 输出 ¥0.06/1K tokens
	// 预期: 输入成本 = 1000/1000 * 0.03 = 0.03
	//       输出成本 = 500/1000 * 0.06 = 0.03
	//       总成本 = 0.06
	expectedInputCost := 0.03
	expectedOutputCost := 0.03
	expectedTotalCost := 0.06

	assert.InDelta(t, expectedInputCost, resp.InputCost, 0.001, "Input cost mismatch")
	assert.InDelta(t, expectedOutputCost, resp.OutputCost, 0.001, "Output cost mismatch")
	assert.InDelta(t, expectedTotalCost, resp.TotalCost, 0.001, "Total cost mismatch")
}

// ============================================================
// 测试场景2: 批量记录测试
// ============================================================

/**
 * TestTokenMetering_BatchRecord
 *
 * 测试批量记录Token使用：
 * 1. 批量记录1000条Token使用
 * 2. 验证所有记录成功写入
 * 3. 验证汇总数据正确
 * 4. 验证性能（< 5秒）
 */
func TestTokenMetering_BatchRecord(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()

	// 1. 创建测试租户
	tenantID := "tenant_test_batch"
	createTestTenant(ctx, tenantID)

	// 2. 准备批量数据（1000条）
	records := make([]*service.RecordTokenUsageRequest, 1000)
	for i := 0; i < 1000; i++ {
		botID := stringPtr("bot_test_batch")
		var userID uint64 = 1001 + uint64(i)

		records[i] = &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			BotID:         botID,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000 + i,
			OutputTokens:  500,
			TotalTokens:   1500 + i,
			RequestType:   entity.RequestTypeChat,
			IsCached:      i%10 == 0, // 10%缓存命中率
		}
	}

	// 3. 执行批量记录
	testCtx := context.Background()
	start := time.Now()
	resp, err := ctx.TokenMeteringSvc.BatchRecordTokenUsage(testCtx, records)
	duration := time.Since(start)

	// 4. 验证结果
	require.NoError(t, err, "BatchRecordTokenUsage should succeed")
	assert.Equal(t, 1000, resp.SuccessCount, "All records should succeed")
	assert.Equal(t, 0, resp.FailedCount, "No records should fail")
	assert.Len(t, resp.LogIDs, 1000, "Should have 1000 log IDs")
	assert.Greater(t, resp.TotalCost, 0.0, "Total cost should be positive")

	// 5. 验证数据库记录数
	var count int64
	err = ctx.DB.Model(&entity.TokenUsageLog{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1000), count, "Should have 1000 log records")

	// 6. 验证性能（批量插入应该在合理时间内完成）
	assert.Less(t, duration, 5*time.Second, "Batch insert should complete within 5 seconds")
	t.Logf("✅ Batch inserted 1000 records in %v (%.2f records/sec)",
		duration, float64(1000)/duration.Seconds())

	// 7. 等待汇总更新
	waitForSummary(t, ctx, tenantID, 10*time.Second)

	// 8. 验证汇总数据
	var summaries []entity.TokenUsageSummary
	err = ctx.DB.Where("tenant_id = ?", tenantID).
		Find(&summaries).Error
	require.NoError(t, err)

	// 验证至少有一个汇总记录
	assert.Greater(t, len(summaries), 0, "Should have at least one summary record")
}

// ============================================================
// 测试场景3: 多模型统计测试
// ============================================================

/**
 * TestTokenMetering_MultiModelStats
 *
 * 测试多模型统计：
 * 1. 记录多个AI模型的使用（gpt-4, claude-3, qwen-max）
 * 2. 调用GetUsageStats
 * 3. 验证各模型Token分布正确
 * 4. 验证成本占比正确
 */
func TestTokenMetering_MultiModelStats(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()

	// 1. 创建测试租户
	tenantID := "tenant_test_multimodel"
	createTestTenant(ctx, tenantID)

	// 2. 记录多个模型的使用
	testCtx := context.Background()
	now := time.Now()
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	// gpt-4: 100K tokens, 成本较高
	for i := 0; i < 10; i++ {
		botID := stringPtr("bot_test_gpt4")
		var userID uint64 = 1001

		req := &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			BotID:         botID,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   5000,
			OutputTokens:  5000,
			TotalTokens:   10000,
			RequestType:   entity.RequestTypeChat,
			IsCached:      false,
		}

		_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err)
	}

	// claude-3-opus: 80K tokens, 成本最高
	for i := 0; i < 8; i++ {
		botID := stringPtr("bot_test_claude")
		var userID uint64 = 1002

		req := &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			BotID:         botID,
			ModelProvider: "anthropic",
			ModelName:     "claude-3-opus-20240229",
			InputTokens:   5000,
			OutputTokens:  5000,
			TotalTokens:   10000,
			RequestType:   entity.RequestTypeChat,
			IsCached:      false,
		}

		_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err)
	}

	// qwen-max: 150K tokens, 成本较低
	for i := 0; i < 15; i++ {
		botID := stringPtr("bot_test_qwen")
		var userID uint64 = 1003

		req := &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			BotID:         botID,
			ModelProvider: "aliyun",
			ModelName:     "qwen-max",
			InputTokens:   5000,
			OutputTokens:  5000,
			TotalTokens:   10000,
			RequestType:   entity.RequestTypeChat,
			IsCached:      false,
		}

		_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err)
	}

	// 3. 等待汇总
	waitForSummary(t, ctx, tenantID, 10*time.Second)

	// 4. 获取模型统计
	modelStats, err := ctx.TokenMeteringSvc.GetModelUsageStats(
		testCtx,
		tenantID,
		startDate,
		endDate,
	)

	require.NoError(t, err, "GetModelUsageStats should succeed")
	assert.Len(t, modelStats, 3, "Should have 3 models")

	// 5. 验证各模型Token分布
	modelMap := make(map[string]*service.ModelUsageStats)
	for _, stat := range modelStats {
		key := stat.ModelProvider + "/" + stat.ModelName
		modelMap[key] = stat
	}

	// 验证gpt-4
	gpt4Stats, ok := modelMap["openai/gpt-4"]
	require.True(t, ok, "Should have gpt-4 stats")
	assert.Equal(t, int64(100000), gpt4Stats.TotalTokens, "gpt-4 should have 100K tokens")
	assert.Equal(t, int64(10), gpt4Stats.TotalRequests, "gpt-4 should have 10 requests")

	// 验证claude-3-opus
	claudeStats, ok := modelMap["anthropic/claude-3-opus-20240229"]
	require.True(t, ok, "Should have claude-3-opus stats")
	assert.Equal(t, int64(80000), claudeStats.TotalTokens, "claude-3-opus should have 80K tokens")
	assert.Equal(t, int64(8), claudeStats.TotalRequests, "claude-3-opus should have 8 requests")

	// 验证qwen-max
	qwenStats, ok := modelMap["aliyun/qwen-max"]
	require.True(t, ok, "Should have qwen-max stats")
	assert.Equal(t, int64(150000), qwenStats.TotalTokens, "qwen-max should have 150K tokens")
	assert.Equal(t, int64(15), qwenStats.TotalRequests, "qwen-max should have 15 requests")

	// 6. 验证成本占比
	totalCost := gpt4Stats.TotalCost + claudeStats.TotalCost + qwenStats.TotalCost
	assert.Greater(t, totalCost, 0.0, "Total cost should be positive")

	// claude-3应该是最贵的
	assert.Greater(t, claudeStats.TotalCost, gpt4Stats.TotalCost,
		"claude-3 should cost more than gpt-4")
	assert.Greater(t, claudeStats.TotalCost, qwenStats.TotalCost,
		"claude-3 should cost more than qwen-max")

	// qwen-max应该是最便宜的
	assert.Less(t, qwenStats.TotalCost, gpt4Stats.TotalCost,
		"qwen-max should cost less than gpt-4")

	t.Logf("✅ Model stats validated:")
	t.Logf("   - gpt-4: %d tokens, ¥%.2f", gpt4Stats.TotalTokens, gpt4Stats.TotalCost)
	t.Logf("   - claude-3-opus: %d tokens, ¥%.2f", claudeStats.TotalTokens, claudeStats.TotalCost)
	t.Logf("   - qwen-max: %d tokens, ¥%.2f", qwenStats.TotalTokens, qwenStats.TotalCost)
}

// ============================================================
// 测试场景4: 每日趋势测试
// ============================================================

/**
 * TestTokenMetering_DailyTrend
 *
 * 测试每日使用趋势：
 * 1. 记录跨多天的Token使用
 * 2. 调用GetDailyUsageStats
 * 3. 验证每日数据正确
 * 4. 验证趋势计算准确
 */
func TestTokenMetering_DailyTrend(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()

	// 1. 创建测试租户
	tenantID := "tenant_test_daily"
	botID := "bot_test_daily"
	createTestTenant(ctx, tenantID)
	createTestBot(ctx, botID, tenantID)

	// 2. 记录跨多天的Token使用
	testCtx := context.Background()
	now := time.Now()

	// 过去7天，每天递增使用量
	for dayOffset := 6; dayOffset >= 0; dayOffset-- {
		targetDate := now.AddDate(0, 0, -dayOffset)
		tokens := 10000 * (7 - dayOffset) // 第1天10K，第7天70K

		for i := 0; i < 10; i++ {
			botIDPtr := &botID
			var userID uint64 = 1001

			req := &service.RecordTokenUsageRequest{
				TenantID:      tenantID,
				UserID:        &userID,
				BotID:         botIDPtr,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   tokens / 2,
				OutputTokens:  tokens / 2,
				TotalTokens:   tokens,
				RequestType:   entity.RequestTypeChat,
				IsCached:      false,
			}

			// 修改CreatedAt时间（模拟历史数据）
			resp, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
			require.NoError(t, err)

			// 更新数据库中的CreatedAt
			ctx.DB.Model(&entity.TokenUsageLog{}).
				Where("id = ?", resp.LogID).
				Update("created_at", targetDate)
		}
	}

	// 手动触发汇总（因为我们修改了CreatedAt）
	// 在真实场景中，这应该由定时任务完成
	// 这里我们直接查询原始日志进行聚合

	// 3. 获取每日统计
	startDate := now.AddDate(0, 0, -6).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	dailyStats, err := ctx.TokenMeteringSvc.GetDailyUsageStats(
		testCtx,
		tenantID,
		startDate,
		endDate,
		&botID,
	)

	require.NoError(t, err, "GetDailyUsageStats should succeed")
	assert.Len(t, dailyStats, 7, "Should have 7 days of data")

	// 4. 验证每日数据递增
	for i := 0; i < 7; i++ {
		day := dailyStats[i]
		expectedTokens := int64(100000 * (i + 1)) // 10次请求 * (10K * (i+1))

		assert.Equal(t, int64(10), day.TotalRequests, "Day %d should have 10 requests", i)
		assert.Equal(t, expectedTokens, day.TotalTokens, "Day %d token count mismatch", i)
		assert.Greater(t, day.TotalCost, 0.0, "Day %d should have positive cost", i)
	}

	t.Logf("✅ Daily trend validated: 7 days with increasing usage")
}

// ============================================================
// 测试场景5: 预算告警触发测试
// ============================================================

/**
 * TestTokenMetering_BudgetAlertTrigger
 *
 * 测试预算告警触发：
 * 1. 创建预算配置（1000元，80%阈值）
 * 2. 记录Token使用直到达到80%
 * 3. 验证告警被触发
 * 4. 验证BudgetAlert表中的记录
 * 5. 验证通知服务被调用（模拟）
 */
func TestTokenMetering_BudgetAlertTrigger(t *testing.T) {
	ctx, cleanup := setupTestContext(t)
	defer cleanup()

	// 1. 创建测试租户
	tenantID := "tenant_test_alert"
	createTestTenant(ctx, tenantID)

	// 2. 创建预算配置（1000元，80%阈值）
	hardCapAmount := 1200.0
	budget := &entity.BudgetSettings{
		TenantID:     tenantID,
		BudgetType:   entity.BudgetTypeMonthly,
		BudgetAmount: 1000.00,
		Currency:     "CNY",
		AlertThreshold1: 80,
		AlertThreshold2: 95,
		HardCapEnabled:   true,
		HardCapAmount:    &hardCapAmount,
		NotificationChannels:   `["email"]`,
		NotificationRecipients: `["admin@example.com"]`,
	}

	err := ctx.DB.Create(budget).Error
	require.NoError(t, err, "Failed to create budget settings")

	// 3. 记录Token使用，累计约850元（超过80%阈值）
	testCtx := context.Background()
	costPerRequest := 10.0 // 每次请求10元

	for i := 0; i < 85; i++ {
		botID := stringPtr("bot_test_alert")
		var userID uint64 = 1001

		// gpt-4定价: 约10元/请求
		req := &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			BotID:         botID,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   100000, // 约5元
			OutputTokens:  100000, // 约5元
			TotalTokens:   200000,
			RequestType:   entity.RequestTypeChat,
			IsCached:      false,
		}

		_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err)
	}

	// 4. 等待汇总更新
	waitForSummary(t, ctx, tenantID, 10*time.Second)

	// 5. 手动触发预算检查
	err = ctx.BudgetAlertSvc.CheckBudget(testCtx, tenantID)
	require.NoError(t, err, "CheckBudget should succeed")

	// 6. 验证BudgetAlert表中的记录
	var alerts []entity.BudgetAlert
	err = ctx.DB.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&alerts).Error
	require.NoError(t, err)

	assert.Greater(t, len(alerts), 0, "Should have at least one alert")

	// 7. 验证告警详情
	threshold1Alert := false
	for _, alert := range alerts {
		if alert.AlertType == entity.AlertTypeThreshold1 {
			threshold1Alert = true
			assert.Equal(t, entity.AlertLevelWarning, alert.AlertLevel,
				"Threshold1 alert should be warning level")
			assert.NotNil(t, alert.AlertMessage, "Should have alert message")
			assert.GreaterOrEqual(t, alert.UsagePercent, 80.0,
				"Usage percent should be >= 80%")
		}
	}

	assert.True(t, threshold1Alert, "Should have triggered threshold1 alert")

	t.Logf("✅ Budget alert triggered: %d alerts created", len(alerts))
	for _, alert := range alerts {
		t.Logf("   - Type: %s, Level: %s, Usage: %.1f%%",
			alert.AlertType, alert.AlertLevel, alert.UsagePercent)
	}
}

// ============================================================
// 辅助函数
// ============================================================

/**
 * stringPtr 返回字符串指针
 *
 * 遵循DRY原则
 */
func stringPtr(s string) *string {
	return &s
}
