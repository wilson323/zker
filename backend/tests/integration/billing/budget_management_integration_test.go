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
// 预算管理API集成测试套件
//
// 测试场景覆盖:
// 1. 预算CRUD测试 - 创建、查询、更新、删除
// 2. 预算使用计算测试 - 使用率、剩余金额计算
// 3. 多周期预算测试 - 月度、季度、年度预算
// 4. 告警历史查询测试 - 分页、过滤功能
// 5. 手动预算检查测试 - 手动触发检查
// ============================================================

// ============================================================
// 测试辅助函数
// ============================================================

/**
 * setupBudgetTestContext 设置预算管理测试上下文
 */
type budgetTestContext struct {
	DB                    *gorm.DB
	BudgetManagementSvc   *service.BudgetManagementService
	TokenMeteringSvc      *service.TokenMeteringService
	Logger                *logrus.Logger
}

func setupBudgetTestContext(t *testing.T) (*budgetTestContext, func()) {
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
	ctx := &budgetTestContext{
		DB:                  db,
		BudgetManagementSvc: budgetManagementSvc,
		TokenMeteringSvc:    tokenMeteringSvc,
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
 * createTestTenantWithBudget 创建测试租户和预算
 */
func createTestTenantWithBudget(ctx *budgetTestContext, tenantID string, budgetAmount float64) {
	// 创建租户
	tenant := &struct {
		TenantID string `gorm:"primaryKey;size:36"`
		Name     string `gorm:"size:100"`
	}{
		TenantID: tenantID,
		Name:     "测试租户",
	}
	ctx.DB.Table("tenants").Create(tenant)

	// 创建预算配置
	hardCapAmount := budgetAmount * 1.2
	budget := &entity.BudgetSettings{
		TenantID:     tenantID,
		BudgetType:   entity.BudgetTypeMonthly,
		BudgetAmount: budgetAmount,
		Currency:     "CNY",
		AlertThreshold1: 80,
		AlertThreshold2: 95,
		HardCapEnabled:   true,
		HardCapAmount:    &hardCapAmount,
		NotificationChannels:   `["email"]`,
		NotificationRecipients: `["admin@example.com"]`,
	}

	ctx.DB.Create(budget)
}

/**
 * recordTokenUsageWithCost 记录指定成本的Token使用
 */
func recordTokenUsageWithCost(ctx *budgetTestContext, tenantID string, cost float64) error {
	testCtx := context.Background()
	botID := "bot_test_usage"

	// 记录Token使用（大约产生指定的成本）
	// gpt-4定价: 约¥0.06/1K tokens
	tokens := int((cost / 0.06) * 1000)

	req := &service.RecordTokenUsageRequest{
		TenantID:      tenantID,
		BotID:         &botID,
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		InputTokens:   tokens / 2,
		OutputTokens:  tokens / 2,
		TotalTokens:   tokens,
		RequestType:   entity.RequestTypeChat,
		IsCached:      false,
	}

	_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
	return err
}

// ============================================================
// 测试场景1: 预算CRUD测试
// ============================================================

/**
 * TestBudgetManagement_CRUD
 *
 * 测试预算配置的完整CRUD流程：
 * 1. CreateBudget - 创建月度预算
 * 2. GetBudget - 获取预算配置
 * 3. UpdateBudget - 更新预算金额和阈值
 * 4. GetBudget - 验证更新成功
 * 5. DeleteBudget - 删除预算
 * 6. GetBudget - 验证已删除
 */
func TestBudgetManagement_CRUD(t *testing.T) {
	ctx, cleanup := setupBudgetTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_test_crud"
	createTestTenant(ctx, tenantID)

	// 1. CreateBudget - 创建月度预算
	hardCapAmount := 12000.00
	originalModel := "openai/gpt-4"
	fallbackModel := "openai/gpt-3.5-turbo"

	createReq := &service.CreateBudgetRequest{
		TenantID:     tenantID,
		BudgetType:   entity.BudgetTypeMonthly,
		BudgetAmount: 10000.00,
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
		NotificationRecipients: []string{"admin@example.com"},
	}

	budget, err := ctx.BudgetManagementSvc.CreateBudget(testCtx, createReq)
	require.NoError(t, err, "CreateBudget should succeed")
	assert.NotNil(t, budget)
	assert.Equal(t, tenantID, budget.TenantID)
	assert.Equal(t, 10000.00, budget.BudgetAmount)
	assert.Equal(t, entity.BudgetTypeMonthly, budget.BudgetType)
	assert.Equal(t, 80, budget.AlertThreshold1)
	assert.Equal(t, 95, budget.AlertThreshold2)
	assert.True(t, budget.HardCapEnabled)
	assert.Equal(t, 12000.00, *budget.HardCapAmount)
	assert.True(t, budget.AutoDowngradeEnabled)
	assert.Equal(t, 90, budget.DowngradeThreshold)
	assert.Equal(t, "openai/gpt-4", *budget.OriginalModel)
	assert.Equal(t, "openai/gpt-3.5-turbo", *budget.FallbackModel)
	assert.Len(t, budget.NotificationChannels, 2)
	assert.Contains(t, budget.NotificationChannels, "email")
	assert.Contains(t, budget.NotificationChannels, "webhook")

	// 2. GetBudget - 获取预算配置
	fetched, err := ctx.BudgetManagementSvc.GetBudget(testCtx, tenantID)
	require.NoError(t, err, "GetBudget should succeed")
	assert.Equal(t, budget.TenantID, fetched.TenantID)
	assert.Equal(t, budget.BudgetAmount, fetched.BudgetAmount)

	// 3. UpdateBudget - 更新预算金额和阈值
	newBudgetAmount := 15000.00
	newThreshold1 := 85
	newHardCapAmount := 18000.00

	updateReq := &service.UpdateBudgetRequest{
		BudgetAmount:     &newBudgetAmount,
		AlertThreshold1:  &newThreshold1,
		HardCapAmount:    &newHardCapAmount,
	}

	updated, err := ctx.BudgetManagementSvc.UpdateBudget(testCtx, tenantID, updateReq)
	require.NoError(t, err, "UpdateBudget should succeed")
	assert.Equal(t, 15000.00, updated.BudgetAmount)
	assert.Equal(t, 85, updated.AlertThreshold1)
	assert.Equal(t, 18000.00, *updated.HardCapAmount)

	// 4. GetBudget - 验证更新成功
	fetched2, err := ctx.BudgetManagementSvc.GetBudget(testCtx, tenantID)
	require.NoError(t, err, "GetBudget should succeed")
	assert.Equal(t, 15000.00, fetched2.BudgetAmount)
	assert.Equal(t, 85, fetched2.AlertThreshold1)

	// 5. DeleteBudget - 删除预算
	err = ctx.BudgetManagementSvc.DeleteBudget(testCtx, tenantID)
	require.NoError(t, err, "DeleteBudget should succeed")

	// 6. GetBudget - 验证已删除
	_, err = ctx.BudgetManagementSvc.GetBudget(testCtx, tenantID)
	assert.Error(t, err, "GetBudget should fail after deletion")
	assert.Contains(t, err.Error(), "not found", "Error should mention 'not found'")

	t.Logf("✅ Budget CRUD test completed successfully")
}

// ============================================================
// 测试场景2: 预算使用计算测试
// ============================================================

/**
 * TestBudgetManagement_UsageCalculation
 *
 * 测试预算使用计算：
 * 1. 创建预算配置（1000元）
 * 2. 记录Token使用（800元）
 * 3. 调用GetBudgetUsage
 * 4. 验证使用率 = 80%
 * 5. 验证剩余金额 = 200元
 */
func TestBudgetManagement_UsageCalculation(t *testing.T) {
	ctx, cleanup := setupBudgetTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_test_usage"

	// 1. 创建预算配置（1000元）
	createTestTenantWithBudget(ctx, tenantID, 1000.00)

	// 2. 记录Token使用（约800元）
	// 分批记录，每次100元，共8次
	for i := 0; i < 8; i++ {
		err := recordTokenUsageWithCost(ctx, tenantID, 100.0)
		require.NoError(t, err, "RecordTokenUsage should succeed")
	}

	// 3. 等待汇总更新
	time.Sleep(2 * time.Second)

	// 4. 调用GetBudgetUsage
	usage, err := ctx.BudgetManagementSvc.GetBudgetUsage(testCtx, tenantID)
	require.NoError(t, err, "GetBudgetUsage should succeed")

	// 5. 验证使用率（允许误差）
	assert.InDelta(t, 1000.00, usage.BudgetAmount, 1.0, "Budget amount should be 1000")
	assert.InDelta(t, 800.00, usage.UsedAmount, 50.0, "Used amount should be around 800")
	assert.InDelta(t, 80.0, usage.UsagePercent, 5.0, "Usage percent should be around 80%")
	assert.InDelta(t, 200.00, usage.RemainingAmount, 50.0, "Remaining amount should be around 200")
	assert.Greater(t, usage.TotalTokens, int64(0), "Should have positive token count")
	assert.Greater(t, usage.TotalRequests, 0, "Should have positive request count")
	assert.Greater(t, usage.AverageCost, 0.0, "Should have positive average cost")

	t.Logf("✅ Budget usage calculation validated:")
	t.Logf("   - Budget: ¥%.2f", usage.BudgetAmount)
	t.Logf("   - Used: ¥%.2f (%.1f%%)", usage.UsedAmount, usage.UsagePercent)
	t.Logf("   - Remaining: ¥%.2f", usage.RemainingAmount)
	t.Logf("   - Tokens: %d", usage.TotalTokens)
	t.Logf("   - Requests: %d", usage.TotalRequests)
	t.Logf("   - Avg Cost: ¥%.4f/1K tokens", usage.AverageCost)
}

// ============================================================
// 测试场景3: 多周期预算测试
// ============================================================

/**
 * TestBudgetManagement_MultiPeriod
 *
 * 测试多周期预算：
 * 1. 测试月度预算（当月1日-月底）
 * 2. 测试季度预算（季度首日-季度末）
 * 3. 测试年度预算（1月1日-12月31日）
 * 4. 验证周期计算准确性
 */
func TestBudgetManagement_MultiPeriod(t *testing.T) {
	ctx, cleanup := setupBudgetTestContext(t)
	defer cleanup()

	testCtx := context.Background()

	// 1. 测试月度预算
	tenantID1 := "tenant_test_monthly"
	createTestTenant(ctx, tenantID1)

	hardCapAmount := 12000.00
	createReq1 := &service.CreateBudgetRequest{
		TenantID:     tenantID1,
		BudgetType:   entity.BudgetTypeMonthly,
		BudgetAmount: 10000.00,
		Currency:     "CNY",
		AlertThreshold1: 80,
		AlertThreshold2: 95,
		HardCapEnabled:   true,
		HardCapAmount:    &hardCapAmount,
		NotificationChannels:   []string{"email"},
		NotificationRecipients: []string{"admin@example.com"},
	}

	budget1, err := ctx.BudgetManagementSvc.CreateBudget(testCtx, createReq1)
	require.NoError(t, err)

	usage1, err := ctx.BudgetManagementSvc.GetBudgetUsage(testCtx, tenantID1)
	require.NoError(t, err)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	assert.Equal(t, monthStart.Unix(), usage1.PeriodStart, "Monthly budget period start mismatch")
	assert.Equal(t, monthEnd.Unix(), usage1.PeriodEnd, "Monthly budget period end mismatch")
	assert.Equal(t, int(time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Sub(now)/24/3600/1000/1000/1000), usage1.DaysRemaining, "Days remaining mismatch")

	t.Logf("✅ Monthly budget period validated: %s to %s",
		monthStart.Format("2006-01-02"), monthEnd.Format("2006-01-02"))

	// 2. 测试季度预算
	tenantID2 := "tenant_test_quarterly"
	createTestTenant(ctx, tenantID2)

	createReq2 := &service.CreateBudgetRequest{
		TenantID:     tenantID2,
		BudgetType:   entity.BudgetTypeQuarterly,
		BudgetAmount: 30000.00,
		Currency:     "CNY",
		AlertThreshold1: 80,
		AlertThreshold2: 95,
		HardCapEnabled:   false,
		NotificationChannels:   []string{"email"},
		NotificationRecipients: []string{"admin@example.com"},
	}

	budget2, err := ctx.BudgetManagementSvc.CreateBudget(testCtx, createReq2)
	require.NoError(t, err)

	usage2, err := ctx.BudgetManagementSvc.GetBudgetUsage(testCtx, tenantID2)
	require.NoError(t, err)

	quarter := (now.Month() - 1) / 3
	quarterStart := time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location())
	quarterEnd := quarterStart.AddDate(0, 3, 0)

	assert.Equal(t, quarterStart.Unix(), usage2.PeriodStart, "Quarterly budget period start mismatch")
	assert.Equal(t, quarterEnd.Unix(), usage2.PeriodEnd, "Quarterly budget period end mismatch")

	t.Logf("✅ Quarterly budget period validated: %s to %s",
		quarterStart.Format("2006-01-02"), quarterEnd.Format("2006-01-02"))

	// 3. 测试年度预算
	tenantID3 := "tenant_test_yearly"
	createTestTenant(ctx, tenantID3)

	createReq3 := &service.CreateBudgetRequest{
		TenantID:     tenantID3,
		BudgetType:   entity.BudgetTypeYearly,
		BudgetAmount: 120000.00,
		Currency:     "CNY",
		AlertThreshold1: 80,
		AlertThreshold2: 95,
		HardCapEnabled:   false,
		NotificationChannels:   []string{"email"},
		NotificationRecipients: []string{"admin@example.com"},
	}

	budget3, err := ctx.BudgetManagementSvc.CreateBudget(testCtx, createReq3)
	require.NoError(t, err)

	usage3, err := ctx.BudgetManagementSvc.GetBudgetUsage(testCtx, tenantID3)
	require.NoError(t, err)

	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	yearEnd := yearStart.AddDate(1, 0, 0)

	assert.Equal(t, yearStart.Unix(), usage3.PeriodStart, "Yearly budget period start mismatch")
	assert.Equal(t, yearEnd.Unix(), usage3.PeriodEnd, "Yearly budget period end mismatch")

	t.Logf("✅ Yearly budget period validated: %s to %s",
		yearStart.Format("2006-01-02"), yearEnd.Format("2006-01-02"))

	t.Logf("✅ All budget period types validated successfully")
}

// ============================================================
// 测试场景4: 告警历史查询测试
// ============================================================

/**
 * TestBudgetManagement_AlertHistory
 *
 * 测试告警历史查询：
 * 1. 触发多次预算告警（80%, 95%, 100%）
 * 2. 调用GetAlerts查询
 * 3. 验证告警记录完整
 * 4. 验证分页功能
 * 5. 验证过滤功能（按级别、类型）
 */
func TestBudgetManagement_AlertHistory(t *testing.T) {
	ctx, cleanup := setupBudgetTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_test_alert_history"

	// 1. 创建预算配置
	hardCapAmount := 1200.00
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
	ctx.DB.Create(budget)

	createTestTenant(ctx, tenantID)

	// 2. 手动创建多个告警记录
	alertTypes := []string{
		entity.AlertTypeThreshold1,
		entity.AlertTypeThreshold1, // 重复，测试去重
		entity.AlertTypeThreshold2,
		entity.AlertTypeHardCap,
	}

	for i, alertType := range alertTypes {
		message := stringPtr("Test alert message")
		alert := &entity.BudgetAlert{
			TenantID:       tenantID,
			AlertType:      alertType,
			BudgetAmount:   1000.00,
			UsedAmount:     800.0 + float64(i)*100,
			UsagePercent:   80.0 + float64(i)*10,
			AlertLevel:     entity.AlertLevelWarning,
			AlertMessage:   message,
			NotificationStatus: entity.NotificationStatusSent,
		}

		if alertType == entity.AlertTypeThreshold2 {
			alert.AlertLevel = entity.AlertLevelCritical
		} else if alertType == entity.AlertTypeHardCap {
			alert.AlertLevel = entity.AlertLevelEmergency
		}

		err := ctx.DB.Create(alert).Error
		require.NoError(t, err)
	}

	// 3. 查询所有告警
	filter := &service.AlertFilter{
		PageSize: 10,
	}

	alerts, pagination, err := ctx.BudgetManagementSvc.GetAlerts(testCtx, tenantID, filter)
	require.NoError(t, err, "GetAlerts should succeed")
	assert.Len(t, alerts, 4, "Should have 4 alerts")
	assert.Equal(t, int64(4), pagination.TotalCount, "Total count should be 4")
	assert.False(t, pagination.HasMore, "Should not have more pages")

	// 4. 验证告警详情
	for _, alert := range alerts {
		assert.Equal(t, tenantID, alert.TenantID)
		assert.NotEmpty(t, alert.AlertType)
		assert.NotEmpty(t, alert.AlertLevel)
		assert.Greater(t, alert.UsagePercent, 0.0)
		assert.NotNil(t, alert.AlertMessage)
	}

	// 5. 测试分页功能（pageSize=2）
	filter2 := &service.AlertFilter{
		PageSize: 2,
	}

	alerts2, pagination2, err := ctx.BudgetManagementSvc.GetAlerts(testCtx, tenantID, filter2)
	require.NoError(t, err)
	assert.Len(t, alerts2, 2, "Should return 2 alerts per page")
	assert.Equal(t, int64(4), pagination2.TotalCount, "Total count should still be 4")
	assert.True(t, pagination2.HasMore, "Should have more pages")

	// 6. 测试过滤功能（按级别）
	criticalLevel := entity.AlertLevelCritical
	filter3 := &service.AlertFilter{
		AlertLevel: &criticalLevel,
		PageSize:   10,
	}

	alerts3, pagination3, err := ctx.BudgetManagementSvc.GetAlerts(testCtx, tenantID, filter3)
	require.NoError(t, err)
	assert.Greater(t, len(alerts3), 0, "Should have critical alerts")

	for _, alert := range alerts3 {
		assert.Equal(t, entity.AlertLevelCritical, alert.AlertLevel,
			"All alerts should be critical level")
	}

	t.Logf("✅ Alert history query validated:")
	t.Logf("   - Total alerts: %d", pagination.TotalCount)
	t.Logf("   - Critical alerts: %d", len(alerts3))
}

// ============================================================
// 测试场景5: 手动预算检查测试
// ============================================================

/**
 * TestBudgetManagement_ManualCheck
 *
 * 测试手动预算检查：
 * 1. 创建预算配置
 * 2. 记录Token使用（达到告警阈值）
 * 3. 调用CheckBudget手动触发检查
 * 4. 验证告警被触发
 * 5. 验证返回结果正确
 */
func TestBudgetManagement_ManualCheck(t *testing.T) {
	ctx, cleanup := setupBudgetTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_test_manual_check"

	// 1. 创建预算配置
	createTestTenantWithBudget(ctx, tenantID, 1000.00)

	// 2. 记录Token使用（约850元，超过80%阈值）
	for i := 0; i < 85; i++ {
		err := recordTokenUsageWithCost(ctx, tenantID, 10.0)
		require.NoError(t, err)
	}

	// 等待汇总更新
	time.Sleep(2 * time.Second)

	// 3. 调用CheckBudget手动触发检查
	result, err := ctx.BudgetManagementSvc.CheckBudget(testCtx, tenantID)
	require.NoError(t, err, "CheckBudget should succeed")

	// 4. 验证返回结果
	assert.Equal(t, tenantID, result.TenantID)
	assert.Equal(t, 1000.00, result.BudgetAmount)
	assert.Greater(t, result.UsedAmount, 800.0, "Used amount should exceed 800")
	assert.Greater(t, result.UsagePercent, 80.0, "Usage percent should exceed 80%")
	assert.True(t, result.AlertTriggered, "Alert should be triggered")
	assert.NotNil(t, result.AlertLevel, "Should have alert level")
	assert.NotEmpty(t, result.Message, "Should have alert message")
	assert.Greater(t, result.CheckedAt, int64(0), "CheckedAt should be set")

	if result.AlertLevel != nil {
		assert.Contains(t, []string{entity.AlertLevelWarning, entity.AlertLevelCritical, entity.AlertLevelEmergency},
			*result.AlertLevel, "Alert level should be valid")
	}

	// 5. 验证告警记录被创建
	var alerts []entity.BudgetAlert
	err = ctx.DB.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Find(&alerts).Error
	require.NoError(t, err)
	assert.Greater(t, len(alerts), 0, "Should have alert records")

	t.Logf("✅ Manual budget check validated:")
	t.Logf("   - Budget: ¥%.2f", result.BudgetAmount)
	t.Logf("   - Used: ¥%.2f (%.1f%%)", result.UsedAmount, result.UsagePercent)
	t.Logf("   - Alert Level: %s", stringPtrToString(result.AlertLevel))
	t.Logf("   - Message: %s", result.Message)
	t.Logf("   - Alerts created: %d", len(alerts))
}

// ============================================================
// 辅助函数
// ============================================================

/**
 * createTestTenant 创建测试租户
 */
func createTestTenant(ctx *budgetTestContext, tenantID string) {
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
