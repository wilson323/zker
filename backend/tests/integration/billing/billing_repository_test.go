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
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
	"github.com/coze-dev/coze-studio/backend/tests/integration"
)

// ============================================================
// 测试套件
// ============================================================

/**
 * 计费系统Repository集成测试套件
 *
 * 使用MySQL 8.4.5测试容器，确保与生产环境完全一致
 * 遵循企业级开发规范和全局一致性原则
 */

// ============================================================
// 测试辅助函数
// ============================================================

/**
 * 设置测试数据库
 *
 * 遵循SOLID原则：单一职责
 *
 * @param t 测试对象
 * @return *gorm.DB GORM数据库实例
 * @return func() 清理函数
 */
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
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

	// 4. 返回清理函数
	cleanup := func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		if testCtx.Cleanup != nil {
			testCtx.Cleanup()
		}
	}

	return db, cleanup
}

// ============================================================
// TokenUsageLogRepository 测试
// ============================================================

/**
 * TestTokenUsageLogRepository_Create
 *
 * 测试创建Token使用日志记录
 * 验证：主键自动生成、数据持久化
 */
func TestTokenUsageLogRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewTokenUsageLogRepository(db)
	ctx := context.Background()

	// 准备测试数据
	var userID uint64 = 1001
	botID := "bot_test_001"
	convID := "conv_test_001"
	msgID := "msg_test_001"

	log := &entity.TokenUsageLog{
		TenantID:       "tenant_test_001",
		UserID:         &userID,
		BotID:          &botID,
		ConversationID: &convID,
		MessageID:      &msgID,
		InputTokens:    1000,
		OutputTokens:   500,
		TotalTokens:    1500,
		ModelProvider:  "openai",
		ModelName:      "gpt-4",
		UnitPrice:      0.03,
		InputCost:      0.03,
		OutputCost:     0.03,
		TotalCost:      0.06,
		RequestType:    "chat",
	}

	// 执行测试
	err := repo.Create(ctx, log)

	// 验证结果
	assert.NoError(t, err, "Create should succeed")
	assert.NotZero(t, log.ID, "ID should be auto-generated")

	// 验证数据库持久化
	var savedLog entity.TokenUsageLog
	err = db.First(&savedLog, log.ID).Error
	require.NoError(t, err, "Should find the created log")
	assert.Equal(t, "tenant_test_001", savedLog.TenantID)
	assert.Equal(t, 1000, savedLog.InputTokens)
	assert.Equal(t, 0.06, savedLog.TotalCost)
}

/**
 * TestTokenUsageLogRepository_BatchCreate
 *
 * 测试批量创建Token使用日志
 * 验证：批量插入性能（250条记录）
 */
func TestTokenUsageLogRepository_BatchCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewTokenUsageLogRepository(db)
	ctx := context.Background()

	// 准备批量测试数据（250条）
	logs := make([]*entity.TokenUsageLog, 250)
	for i := 0; i < 250; i++ {
		logs[i] = &entity.TokenUsageLog{
			TenantID:      "tenant_test_001",
			InputTokens:   1000 + i,
			OutputTokens:  500,
			TotalTokens:   1500 + i,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			UnitPrice:     0.03,
			InputCost:     0.03,
			OutputCost:    0.015,
			TotalCost:     0.045,
			RequestType:   "chat",
		}
	}

	// 执行测试
	start := time.Now()
	err := repo.BatchCreate(ctx, logs)
	duration := time.Since(start)

	// 验证结果
	assert.NoError(t, err, "BatchCreate should succeed")

	// 验证数据库记录数
	var count int64
	err = db.Model(&entity.TokenUsageLog{}).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(250), count, "Should insert all 250 logs")

	// 性能验证：批量插入应该在合理时间内完成
	assert.Less(t, duration, 5*time.Second, "Batch insert should complete within 5 seconds")
	t.Logf("Batch inserted 250 logs in %v", duration)
}

/**
 * TestTokenUsageLogRepository_GetByDateRange
 *
 * 测试按日期范围查询Token使用日志
 * 验证：日期范围过滤、排序
 */
func TestTokenUsageLogRepository_GetByDateRange(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewTokenUsageLogRepository(db)
	ctx := context.Background()

	// 准备测试数据（跨3天）
	baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10; i++ {
		log := &entity.TokenUsageLog{
			TenantID:      "tenant_test_001",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			UnitPrice:     0.03,
			InputCost:     0.03,
			OutputCost:    0.015,
			TotalCost:     0.045,
			RequestType:   "chat",
			CreatedAt:     baseDate.Add(time.Duration(i) * 24 * time.Hour),
		}
		require.NoError(t, db.Create(log).Error)
	}

	// 测试用例1：查询前3天的数据
	startDate := "2025-01-01"
	endDate := "2025-01-03"
	logs, err := repo.GetByDateRange(ctx, "tenant_test_001", startDate, endDate)

	assert.NoError(t, err, "GetByDateRange should succeed")
	assert.Len(t, logs, 4, "Should return 4 logs (day 1,2,3,4)")
	assert.Equal(t, "2025-01-01", logs[0].CreatedAt.Format("2006-01-02"))

	// 测试用例2：查询不存在的租户
	logs, err = repo.GetByDateRange(ctx, "tenant_nonexistent", startDate, endDate)
	assert.NoError(t, err, "Should succeed even for non-existent tenant")
	assert.Len(t, logs, 0, "Should return empty list")
}

/**
 * TestTokenUsageLogRepository_List
 *
 * 测试分页查询Token使用日志
 * 验证：游标分页、排序
 */
func TestTokenUsageLogRepository_List(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewTokenUsageLogRepository(db)
	ctx := context.Background()

	// 准备测试数据（25条）
	for i := 0; i < 25; i++ {
		log := &entity.TokenUsageLog{
			TenantID:      "tenant_test_001",
			InputTokens:   1000 + i,
			OutputTokens:  500,
			TotalTokens:   1500 + i,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			UnitPrice:     0.03,
			InputCost:     0.03,
			OutputCost:    0.015,
			TotalCost:     0.045,
			RequestType:   "chat",
			CreatedAt:     time.Now().Add(time.Duration(i) * time.Second),
		}
		require.NoError(t, db.Create(log).Error)
	}

	// 测试用例1：第一页（pageSize=10）
	filter := &repository.TokenUsageLogFilter{
		TenantID:  "tenant_test_001",
		PageSize:  10,
		PageToken: "",
	}
	logs, total, err := repo.List(ctx, filter)

	assert.NoError(t, err, "List should succeed")
	assert.Len(t, logs, 10, "Should return 10 logs")
	assert.Equal(t, int64(25), total, "Should have total 25 logs")

	// 测试用例2：第二页
	if len(logs) > 0 {
		nextPageToken := logs[len(logs)-1].CreatedAt.Format(time.RFC3339Nano)
		filter.PageToken = nextPageToken
		logs, total, err = repo.List(ctx, filter)

		assert.NoError(t, err)
		assert.True(t, len(logs) > 0, "Should return more logs")
		assert.Equal(t, int64(25), total)
	}
}

// ============================================================
// TokenUsageSummaryRepository 测试
// ============================================================

/**
 * TestTokenUsageSummaryRepository_Upsert
 *
 * 测试创建或更新Token使用汇总
 * 验证：UPSERT幂等性、累加逻辑
 */
func TestTokenUsageSummaryRepository_Upsert(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewTokenUsageSummaryRepository(db)
	ctx := context.Background()

	summaryDate := "2025-01-01"
	botID := "bot_test_001"
	var summaryHour uint8 = 14
	var avgResponseTime float64 = 1200.5

	// 准备测试数据
	summary := &entity.TokenUsageSummary{
		TenantID:           "tenant_test_001",
		BotID:              &botID,
		SummaryDate:        summaryDate,
		SummaryHour:        &summaryHour,
		TotalInputTokens:   100000,
		TotalOutputTokens:  50000,
		TotalTokens:        150000,
		TotalCost:          4500.00,
		TotalRequests:      100,
		CachedRequests:     10,
		AvgResponseTime:    &avgResponseTime,
		ModelDistribution:  `{"openai/gpt-4": {"tokens": 150000, "cost": 4500.00}}`,
	}

	// 第一次插入
	err := repo.Upsert(ctx, summary)
	assert.NoError(t, err, "First Upsert should succeed")

	// 验证数据库中的记录
	var saved entity.TokenUsageSummary
	err = db.Where("tenant_id = ? AND bot_id = ? AND summary_date = ? AND summary_hour = ?",
		"tenant_test_001", botID, summaryDate, 14).First(&saved).Error
	require.NoError(t, err)
	assert.Equal(t, int64(100000), saved.TotalInputTokens)
	assert.Equal(t, 4500.00, saved.TotalCost)

	// 第二次Upsert（测试累加逻辑）
	summary.TotalInputTokens = 50000
	summary.TotalOutputTokens = 25000
	summary.TotalTokens = 75000
	summary.TotalCost = 2250.00
	summary.TotalRequests = 50

	err = repo.Upsert(ctx, summary)
	assert.NoError(t, err, "Second Upsert should succeed")

	// 验证累加后的数据
	err = db.Where("tenant_id = ? AND bot_id = ? AND summary_date = ? AND summary_hour = ?",
		"tenant_test_001", botID, summaryDate, 14).First(&saved).Error
	require.NoError(t, err)
	assert.Equal(t, int64(150000), saved.TotalInputTokens, "Should accumulate: 100000 + 50000")
	assert.Equal(t, 6750.00, saved.TotalCost, "Should accumulate: 4500.00 + 2250.00")
	assert.Equal(t, 150, saved.TotalRequests, "Should accumulate: 100 + 50")
}

// ============================================================
// BudgetSettingsRepository 测试
// ============================================================

/**
 * TestBudgetSettingsRepository_CRUD
 *
 * 测试预算配置的CRUD操作
 * 验证：创建、查询、更新、唯一约束
 */
func TestBudgetSettingsRepository_CRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewBudgetSettingsRepository(db)
	ctx := context.Background()

	// 测试创建
	var hardCapAmount float64 = 12000.00
	settings := &entity.BudgetSettings{
		TenantID:             "tenant_test_001",
		BudgetType:           "monthly",
		BudgetAmount:         10000.00,
		Currency:             "CNY",
		AlertThreshold1:      80,
		AlertThreshold2:      95,
		HardCapEnabled:       true,
		HardCapAmount:        &hardCapAmount,
		AutoDowngradeEnabled: true,
		DowngradeThreshold:   90,
		OriginalModel:        stringPtr("openai/gpt-4"),
		FallbackModel:        stringPtr("openai/gpt-3.5-turbo"),
		NotificationChannels: `["email","sms"]`,
		NotificationRecipients: `["admin@example.com"]`,
	}

	err := repo.Create(ctx, settings)
	assert.NoError(t, err, "Create should succeed")
	assert.NotZero(t, settings.ID, "ID should be auto-generated")

	// 测试查询
	result, err := repo.GetByTenantID(ctx, "tenant_test_001")
	assert.NoError(t, err, "GetByTenantID should succeed")
	assert.NotNil(t, result, "Should return settings")
	assert.Equal(t, "tenant_test_001", result.TenantID)
	assert.Equal(t, 10000.00, result.BudgetAmount)

	// 测试更新
	settings.BudgetAmount = 15000.00
	settings.AlertThreshold1 = 85
	err = repo.Update(ctx, settings)
	assert.NoError(t, err, "Update should succeed")

	// 验证更新后的数据
	var updated entity.BudgetSettings
	err = db.Where("tenant_id = ?", "tenant_test_001").First(&updated).Error
	require.NoError(t, err)
	assert.Equal(t, 15000.00, updated.BudgetAmount)
	assert.Equal(t, 85, updated.AlertThreshold1)

	// 测试唯一约束（同一租户不能创建第二条预算配置）
	duplicate := &entity.BudgetSettings{
		TenantID:     "tenant_test_001",
		BudgetType:   "monthly",
		BudgetAmount: 20000.00,
	}
	err = repo.Create(ctx, duplicate)
	assert.Error(t, err, "Should fail due to unique constraint")
}

// ============================================================
// BudgetAlertRepository 测试
// ============================================================

/**
 * TestBudgetAlertRepository_CheckAlertSentToday
 *
 * 测试告警去重逻辑
 * 验证：同一天同一类型告警只发送一次
 */
func TestBudgetAlertRepository_CheckAlertSentToday(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewBudgetAlertRepository(db)
	ctx := context.Background()

	// 准备测试数据：创建今天的告警
	today := time.Now()
	message := "测试告警"
	alert := &entity.BudgetAlert{
		TenantID:            "tenant_test_001",
		AlertType:           "threshold_1",
		BudgetAmount:         10000.00,
		UsedAmount:           8000.00,
		UsagePercent:         80.0,
		AlertLevel:           "warning",
		AlertMessage:         &message,
		NotificationStatus:   "sent",
		CreatedAt:            today,
	}
	require.NoError(t, db.Create(alert).Error)

	// 测试用例1：检查今天是否已发送告警
	todayStr := today.Format("2006-01-02")
	sent, err := repo.CheckAlertSentToday(ctx, "tenant_test_001", "threshold_1", todayStr)

	assert.NoError(t, err, "CheckAlertSentToday should succeed")
	assert.True(t, sent, "Should return true for existing alert today")

	// 测试用例2：检查不同类型的告警
	sent, err = repo.CheckAlertSentToday(ctx, "tenant_test_001", "threshold_2", todayStr)
	assert.NoError(t, err)
	assert.False(t, sent, "Should return false for different alert type")

	// 测试用例3：检查不存在的租户
	sent, err = repo.CheckAlertSentToday(ctx, "tenant_nonexistent", "threshold_1", todayStr)
	assert.NoError(t, err)
	assert.False(t, sent, "Should return false for non-existent tenant")
}

// ============================================================
// CostOptimizationSuggestionRepository 测试
// ============================================================

/**
 * TestCostOptimizationSuggestionRepository_CRUD
 *
 * 测试成本优化建议的CRUD操作
 * 验证：创建、查询、状态更新、聚合查询
 */
func TestCostOptimizationSuggestionRepository_CRUD(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repository.NewCostOptimizationSuggestionRepository(db)
	ctx := context.Background()

	// 测试创建
	reason := "将gpt-4降级到gpt-3.5-turbo可节省60%成本"
	targetBotID := "bot_test_001"
	currentModel := "openai/gpt-4"
	suggestedModel := "openai/gpt-3.5-turbo"
	var estimatedSaving float64 = 3000.00
	var estimatedPercent float64 = 60.0

	suggestion := &entity.CostOptimizationSuggestion{
		TenantID:               "tenant_test_001",
		SuggestionType:         "model_downgrade",
		AnalysisPeriodStart:    "2025-01-01",
		AnalysisPeriodEnd:      "2025-01-31",
		TargetBotID:            &targetBotID,
		CurrentModel:           &currentModel,
		SuggestedModel:         &suggestedModel,
		Reason:                 &reason,
		EstimatedMonthlySaving: &estimatedSaving,
		EstimatedSavingPercent: &estimatedPercent,
		Status:                 "pending",
	}

	err := repo.Create(ctx, suggestion)
	assert.NoError(t, err, "Create should succeed")
	assert.NotZero(t, suggestion.ID, "ID should be auto-generated")

	// 测试查询
	pendingStatus := "pending"
	suggestions, err := repo.GetByTenant(ctx, "tenant_test_001", &pendingStatus, 10)
	assert.NoError(t, err, "GetByTenant should succeed")
	assert.Len(t, suggestions, 1, "Should return 1 suggestion")

	// 测试状态更新
	var appliedBy uint64 = 1001
	appliedAt := time.Now()
	err = repo.UpdateStatus(ctx, suggestion.ID, "applied", &appliedAt, &appliedBy)
	assert.NoError(t, err, "UpdateStatus should succeed")

	// 验证更新后的状态
	var updated entity.CostOptimizationSuggestion
	err = db.First(&updated, suggestion.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "applied", updated.Status)
	assert.NotNil(t, updated.AppliedAt)
	assert.NotNil(t, updated.AppliedBy)
	assert.Equal(t, uint64(1001), *updated.AppliedBy)

	// 测试聚合查询：创建多个建议
	for i := 0; i < 5; i++ {
		estSaving := 1000.00 * float64(i+1)
		estPercent := 50.0

		s := &entity.CostOptimizationSuggestion{
			TenantID:               "tenant_test_001",
			SuggestionType:         "model_downgrade",
			AnalysisPeriodStart:    "2025-01-01",
			AnalysisPeriodEnd:      "2025-01-31",
			CurrentModel:           &currentModel,
			SuggestedModel:         &suggestedModel,
			Reason:                 &reason,
			EstimatedMonthlySaving: &estSaving,
			EstimatedSavingPercent: &estPercent,
			Status:                 "pending",
		}
		require.NoError(t, db.Create(s).Error)
	}

	// 查询总预计节省
	totalSaving, err := repo.GetEstimatedMonthlySaving(ctx, "tenant_test_001", nil)
	assert.NoError(t, err, "GetEstimatedMonthlySaving should succeed")
	assert.Equal(t, 15000.00, totalSaving, "Should sum all savings: 1000+2000+3000+4000+5000=15000")
}

// ============================================================
// 辅助函数
// ============================================================

/**
 * stringPtr 返回字符串指针
 *
 * 遵循DRY原则：避免重复创建指针的代码
 */
func stringPtr(s string) *string {
	return &s
}
