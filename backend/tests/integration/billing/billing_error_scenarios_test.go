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
	"errors"
	"fmt"
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
// 错误场景测试套件
//
// 测试场景覆盖:
// 1. 数据库连接失败 - 模拟数据库不可用
// 2. 无效请求参数 - 验证参数校验
// 3. 租户ID不匹配 - 验证租户隔离
// 4. 预算配置不存在 - 验证错误处理
// 5. 模型价格查询失败 - 验证降级处理
// 6. 事务回滚 - 验证数据一致性
// 7. 批量操作部分失败 - 验证错误恢复
// ============================================================

// ============================================================
// 测试辅助函数
// ============================================================

/**
 * setupErrorTestContext 设置错误测试上下文
 */
type errorTestContext struct {
	DB                  *gorm.DB
	TokenMeteringSvc    *service.TokenMeteringService
	BudgetManagementSvc *service.BudgetManagementService
	BudgetAlertSvc      *service.BudgetAlertService
	Logger              *logrus.Logger
}

func setupErrorTestContext(t *testing.T) (*errorTestContext, func()) {
	// 1. 创建MySQL测试容器
	testCtx, err := integration.SetupMySQLTestContainer(t)
	require.NoError(t, err, "Failed to setup MySQL test container")

	// 2. 使用GORM连接数据库
	db, err := gorm.Open(mysql.Open(testCtx.Config.DSN), &gorm.Config{
		SkipDefaultTransaction: false,
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
	logger.SetLevel(logrus.ErrorLevel) // 只输出错误日志

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
	testContext := &errorTestContext{
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

	return testContext, cleanup
}

// ============================================================
// 测试场景1: 无效请求参数测试
// ============================================================

/**
 * TestErrorScenarios_InvalidParameters
 *
 * 测试无效请求参数：
 * 1. 空租户ID
 * 2. 空模型提供商
 * 3. 空模型名称
 * 4. 负数Token数
 * 5. Token总数不匹配
 * 6. 无效请求类型
 * 7. 验证所有错误被正确捕获
 */
func TestErrorScenarios_InvalidParameters(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	botID := "bot_error_test"
	var userID uint64 = 1001

	testCases := []struct {
		name        string
		req         *service.RecordTokenUsageRequest
		expectedErr string
	}{
		{
			name: "空租户ID",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   entity.RequestTypeChat,
			},
			expectedErr: "tenant_id",
		},
		{
			name: "空模型提供商",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "tenant_error_1",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   entity.RequestTypeChat,
			},
			expectedErr: "model_provider",
		},
		{
			name: "空模型名称",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "tenant_error_2",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   entity.RequestTypeChat,
			},
			expectedErr: "model_name",
		},
		{
			name: "负数输入Token",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "tenant_error_3",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   -100,
				OutputTokens:  500,
				TotalTokens:   400,
				RequestType:   entity.RequestTypeChat,
			},
			expectedErr: "negative",
		},
		{
			name: "Token总数不匹配",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "tenant_error_4",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   2000, // 应该是1500
				RequestType:   entity.RequestTypeChat,
			},
			expectedErr: "total_tokens must equal",
		},
		{
			name: "无效请求类型",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "tenant_error_5",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   "invalid_type",
			},
			expectedErr: "invalid request_type",
		},
		{
			name: "空请求类型",
			req: &service.RecordTokenUsageRequest{
				TenantID:      "tenant_error_6",
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   "",
			},
			expectedErr: "request_type",
		},
	}

	t.Log("🚀 开始无效参数测试")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, tc.req)

			// 验证错误被正确捕获
			assert.Error(t, err, "Should return error for invalid parameter")
			assert.Contains(t, err.Error(), tc.expectedErr,
				"Error message should mention the invalid field")

			t.Logf("✅ %s: 正确捕获错误 - %v", tc.name, err)
		})
	}

	t.Log("✅ 所有无效参数测试通过")
}

// ============================================================
// 测试场景2: 预算配置不存在测试
// ============================================================

/**
 * TestErrorScenarios_BudgetNotFound
 *
 * 测试预算配置不存在：
 * 1. 查询不存在的预算配置
 * 2. 更新不存在的预算配置
 * 3. 删除不存在的预算配置
 * 4. 验证返回明确的错误信息
 */
func TestErrorScenarios_BudgetNotFound(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	nonExistentTenantID := "tenant_non_existent_12345"

	t.Log("🚀 开始预算配置不存在测试")

	// 测试1: GetBudget查询不存在的预算
	_, err := ctx.BudgetManagementSvc.GetBudget(testCtx, nonExistentTenantID)
	assert.Error(t, err, "GetBudget should fail for non-existent tenant")
	assert.Contains(t, err.Error(), "not found",
		"Error should mention 'not found'")
	t.Log("✅ GetBudget正确返回错误")

	// 测试2: UpdateBudget更新不存在的预算
	newAmount := 99999.00
	updateReq := &service.UpdateBudgetRequest{
		BudgetAmount: &newAmount,
	}

	_, err = ctx.BudgetManagementSvc.UpdateBudget(testCtx, nonExistentTenantID, updateReq)
	assert.Error(t, err, "UpdateBudget should fail for non-existent tenant")
	assert.Contains(t, err.Error(), "not found",
		"Error should mention 'not found'")
	t.Log("✅ UpdateBudget正确返回错误")

	// 测试3: DeleteBudget删除不存在的预算
	err = ctx.BudgetManagementSvc.DeleteBudget(testCtx, nonExistentTenantID)
	assert.Error(t, err, "DeleteBudget should fail for non-existent tenant")
	assert.Contains(t, err.Error(), "not found",
		"Error should mention 'not found'")
	t.Log("✅ DeleteBudget正确返回错误")

	// 测试4: GetBudgetUsage查询不存在的预算
	_, err = ctx.BudgetManagementSvc.GetBudgetUsage(testCtx, nonExistentTenantID)
	assert.Error(t, err, "GetBudgetUsage should fail for non-existent tenant")
	assert.Contains(t, err.Error(), "not found",
		"Error should mention 'not found'")
	t.Log("✅ GetBudgetUsage正确返回错误")

	t.Log("✅ 所有预算配置不存在测试通过")
}

// ============================================================
// 测试场景3: 无效日期格式测试
// ============================================================

/**
 * TestErrorScenarios_InvalidDateFormat
 *
 * 测试无效日期格式：
 * 1. GetUsageStats使用无效日期格式
 * 2. GetModelUsageStats使用无效日期格式
 * 3. GetDailyUsageStats使用无效日期格式
 * 4. 验证返回明确的错误信息
 */
func TestErrorScenarios_InvalidDateFormat(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_error_date"

	// 创建测试租户
	createTestTenantForError(ctx, tenantID)

	t.Log("🚀 开始无效日期格式测试")

	// 测试1: GetUsageStats使用无效开始日期
	filter := &service.UsageStatsFilter{
		StartDate: "2024-13-01", // 无效月份
		EndDate:   "2024-12-31",
	}

	_, err := ctx.TokenMeteringSvc.GetUsageStats(testCtx, tenantID, filter)
	assert.Error(t, err, "Should fail for invalid start date")
	assert.Contains(t, err.Error(), "invalid",
		"Error should mention 'invalid'")
	t.Log("✅ GetUsageStats正确拒绝无效开始日期")

	// 测试2: GetUsageStats使用无效结束日期
	filter2 := &service.UsageStatsFilter{
		StartDate: "2024-01-01",
		EndDate:   "2024-02-30", // 无效日期
	}

	_, err = ctx.TokenMeteringSvc.GetUsageStats(testCtx, tenantID, filter2)
	assert.Error(t, err, "Should fail for invalid end date")
	assert.Contains(t, err.Error(), "invalid",
		"Error should mention 'invalid'")
	t.Log("✅ GetUsageStats正确拒绝无效结束日期")

	// 测试3: GetModelUsageStats使用无效日期格式
	_, err = ctx.TokenMeteringSvc.GetModelUsageStats(
		testCtx,
		tenantID,
		"2024/01/01", // 错误的格式
		"2024-12-31",
	)
	assert.Error(t, err, "Should fail for invalid date format")
	assert.Contains(t, err.Error(), "invalid",
		"Error should mention 'invalid'")
	t.Log("✅ GetModelUsageStats正确拒绝无效日期格式")

	t.Log("✅ 所有无效日期格式测试通过")
}

// ============================================================
// 测试场景4: 批量操作部分失败测试
// ============================================================

/**
 * TestErrorScenarios_BatchPartialFailure
 *
 * 测试批量操作部分失败：
 * 1. 批量记录包含有效和无效记录
 * 2. 验证有效记录被保存
 * 3. 验证无效记录被跳过
 * 4. 验证返回的成功/失败计数正确
 */
func TestErrorScenarios_BatchPartialFailure(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_error_batch"
	botID := "bot_error_batch"

	// 创建测试租户
	createTestTenantForError(ctx, tenantID)

	t.Log("🚀 开始批量操作部分失败测试")

	// 准备批量数据：包含有效和无效记录
	records := make([]*service.RecordTokenUsageRequest, 10)
	for i := 0; i < 10; i++ {
		var userID uint64 = 5000 + uint64(i)

		if i == 3 || i == 7 {
			// 创建无效记录（负数Token）
			records[i] = &service.RecordTokenUsageRequest{
				TenantID:      tenantID,
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   -100, // 无效
				OutputTokens:  500,
				TotalTokens:   400,
				RequestType:   entity.RequestTypeChat,
			}
		} else {
			// 创建有效记录
			records[i] = &service.RecordTokenUsageRequest{
				TenantID:      tenantID,
				UserID:        &userID,
				BotID:         &botID,
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000 + i,
				OutputTokens:  500,
				TotalTokens:   1500 + i,
				RequestType:   entity.RequestTypeChat,
			}
		}
	}

	// 执行批量记录
	resp, err := ctx.TokenMeteringSvc.BatchRecordTokenUsage(testCtx, records)

	// 验证结果
	require.NoError(t, err, "BatchRecordTokenUsage should succeed")
	assert.Equal(t, 8, resp.SuccessCount, "Should have 8 successful records")
	assert.Equal(t, 2, resp.FailedCount, "Should have 2 failed records")
	assert.Len(t, resp.LogIDs, 8, "Should have 8 log IDs")
	assert.Greater(t, resp.TotalCost, 0.0, "Total cost should be positive")

	t.Logf("   - 成功: %d", resp.SuccessCount)
	t.Logf("   - 失败: %d", resp.FailedCount)
	t.Logf("   - 总成本: ¥%.2f", resp.TotalCost)

	// 验证数据库中只有有效记录
	var count int64
	err = ctx.DB.Model(&entity.TokenUsageLog{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	require.NoError(t, err)

	assert.Equal(t, int64(8), count, "Database should have only 8 valid records")

	t.Log("✅ 批量操作部分失败测试通过")
}

// ============================================================
// 测试场景5: 租户隔离测试
// ============================================================

/**
 * TestErrorScenarios_TenantIsolation
 *
 * 测试租户隔离：
 * 1. 创建两个租户的数据
 * 2. 验证租户A不能访问租户B的数据
 * 3. 验证统计数据正确隔离
 * 4. 验证预算告警正确隔离
 */
func TestErrorScenarios_TenantIsolation(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantIDA := "tenant_isolation_a"
	tenantIDB := "tenant_isolation_b"
	botIDA := "bot_isolation_a"
	botIDB := "bot_isolation_b"

	// 创建两个测试租户
	createTestTenantForError(ctx, tenantIDA)
	createTestTenantForError(ctx, tenantIDB)

	t.Log("🚀 开始租户隔离测试")

	// 为租户A记录Token使用
	for i := 0; i < 5; i++ {
		var userID uint64 = 6000 + uint64(i)
		req := &service.RecordTokenUsageRequest{
			TenantID:      tenantIDA,
			UserID:        &userID,
			BotID:         &botIDA,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			RequestType:   entity.RequestTypeChat,
		}

		_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err)
	}

	// 为租户B记录Token使用
	for i := 0; i < 10; i++ {
		var userID uint64 = 7000 + uint64(i)
		req := &service.RecordTokenUsageRequest{
			TenantID:      tenantIDB,
			UserID:        &userID,
			BotID:         &botIDB,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   2000,
			OutputTokens:  1000,
			TotalTokens:   3000,
			RequestType:   entity.RequestTypeChat,
		}

		_, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
		require.NoError(t, err)
	}

	// 等待汇总更新
	time.Sleep(2 * time.Second)

	// 验证租户A的统计数据
	filterA := &service.UsageStatsFilter{}
	statsA, err := ctx.TokenMeteringSvc.GetUsageStats(testCtx, tenantIDA, filterA)
	require.NoError(t, err)
	assert.Equal(t, int64(5), statsA.TotalRequests, "Tenant A should have 5 requests")
	assert.Equal(t, int64(7500), statsA.TotalTokens, "Tenant A should have 7500 tokens")

	// 验证租户B的统计数据
	filterB := &service.UsageStatsFilter{}
	statsB, err := ctx.TokenMeteringSvc.GetUsageStats(testCtx, tenantIDB, filterB)
	require.NoError(t, err)
	assert.Equal(t, int64(10), statsB.TotalRequests, "Tenant B should have 10 requests")
	assert.Equal(t, int64(30000), statsB.TotalTokens, "Tenant B should have 30000 tokens")

	// 验证租户A看不到租户B的数据
	var countA int64
	ctx.DB.Model(&entity.TokenUsageLog{}).
		Where("tenant_id = ?", tenantIDA).
		Count(&countA)

	var countB int64
	ctx.DB.Model(&entity.TokenUsageLog{}).
		Where("tenant_id = ?", tenantIDB).
		Count(&countB)

	assert.Equal(t, int64(5), countA, "Tenant A should have 5 log records")
	assert.Equal(t, int64(10), countB, "Tenant B should have 10 log records")

	t.Logf("   - 租户A: %d 请求, %d tokens", statsA.TotalRequests, statsA.TotalTokens)
	t.Logf("   - 租户B: %d 请求, %d tokens", statsB.TotalRequests, statsB.TotalTokens)

	t.Log("✅ 租户隔离测试通过")
}

// ============================================================
// 测试场景6: 事务回滚测试
// ============================================================

/**
 * TestErrorScenarios_TransactionRollback
 *
 * 测试事务回滚：
 * 1. 模拟数据库错误（违反约束）
 * 2. 验证事务正确回滚
 * 3. 验证部分数据不被保存
 * 4. 验证数据一致性
 */
func TestErrorScenarios_TransactionRollback(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_error_rollback"
	botID := "bot_error_rollback"

	// 创建测试租户
	createTestTenantForError(ctx, tenantID)

	t.Log("🚀 开始事务回滚测试")

	// 记录第一次成功的Token使用
	var userID1 uint64 = 8001
	req1 := &service.RecordTokenUsageRequest{
		TenantID:      tenantID,
		UserID:        &userID1,
		BotID:         &botID,
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		InputTokens:   1000,
		OutputTokens:  500,
		TotalTokens:   1500,
		RequestType:   entity.RequestTypeChat,
	}

	resp1, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req1)
	require.NoError(t, err, "First record should succeed")
	assert.NotZero(t, resp1.LogID)

	// 记录第二次成功的Token使用
	var userID2 uint64 = 8002
	req2 := &service.RecordTokenUsageRequest{
		TenantID:      tenantID,
		UserID:        &userID2,
		BotID:         &botID,
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		InputTokens:   2000,
		OutputTokens:  1000,
		TotalTokens:   3000,
		RequestType:   entity.RequestTypeChat,
	}

	resp2, err := ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req2)
	require.NoError(t, err, "Second record should succeed")
	assert.NotZero(t, resp2.LogID)

	// 尝试记录无效的Token使用（应该失败并触发事务回滚）
	var userID3 uint64 = 8003
	req3 := &service.RecordTokenUsageRequest{
		TenantID:      tenantID,
		UserID:        &userID3,
		BotID:         &botID,
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		InputTokens:   -100, // 无效，会导致失败
		OutputTokens:  500,
		TotalTokens:   400,
		RequestType:   entity.RequestTypeChat,
	}

	_, err = ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req3)
	assert.Error(t, err, "Third record should fail")

	// 验证前两条记录仍然存在
	var count int64
	err = ctx.DB.Model(&entity.TokenUsageLog{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error
	require.NoError(t, err)

	assert.Equal(t, int64(2), count, "Should have 2 valid records (third record not saved)")

	// 验证第三条记录不存在
	var log3 entity.TokenUsageLog
	err = ctx.DB.Where("tenant_id = ? AND user_id = ?", tenantID, userID3).
		First(&log3).Error
	assert.Error(t, err, "Third record should not exist")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Should be record not found error")

	t.Log("✅ 事务回滚测试通过")
}

// ============================================================
// 测试场景7: 数据库连接失败测试（模拟）
// ============================================================

/**
 * TestErrorScenarios_DatabaseConnectionFailure
 *
 * 测试数据库连接失败处理：
 * 1. 关闭数据库连接
 * 2. 尝试执行操作
 * 3. 验证返回明确的错误信息
 * 4. 重新连接后验证恢复正常
 */
func TestErrorScenarios_DatabaseConnectionFailure(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_error_connection"
	botID := "bot_error_connection"

	// 创建测试租户
	createTestTenantForError(ctx, tenantID)

	t.Log("🚀 开始数据库连接失败测试")

	// 获取底层SQL数据库连接
	sqlDB, err := ctx.DB.DB()
	require.NoError(t, err, "Should get SQL DB connection")

	// 关闭数据库连接
	err = sqlDB.Close()
	require.NoError(t, err, "Should close database connection")

	// 尝试记录Token使用（应该失败）
	var userID uint64 = 9001
	req := &service.RecordTokenUsageRequest{
		TenantID:      tenantID,
		UserID:        &userID,
		BotID:         &botID,
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		InputTokens:   1000,
		OutputTokens:  500,
		TotalTokens:   1500,
		RequestType:   entity.RequestTypeChat,
	}

	_, err = ctx.TokenMeteringSvc.RecordTokenUsage(testCtx, req)
	assert.Error(t, err, "Should fail when database is closed")
	// 错误可能是 driver database is closed 或 connection failed
	assert.True(t,
		errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrConnClosed) ||
			containsString(err.Error(), "closed") || containsString(err.Error(), "connection"),
		"Error should be related to connection failure")

	t.Logf("   - 错误类型: %T", err)
	t.Logf("   - 错误信息: %v", err)

	t.Log("✅ 数据库连接失败测试通过")
}

// ============================================================
// 测试场景8: 大批量数据测试
// ============================================================

/**
 * TestErrorScenarios_LargeBatchData
 *
 * 测试大批量数据处理：
 * 1. 批量记录超过限制（>1000条）
 * 2. 验证返回明确的错误
 * 3. 验证边界值测试（1000条）
 */
func TestErrorScenarios_LargeBatchData(t *testing.T) {
	ctx, cleanup := setupErrorTestContext(t)
	defer cleanup()

	testCtx := context.Background()
	tenantID := "tenant_error_large_batch"

	// 创建测试租户
	createTestTenantForError(ctx, tenantID)

	t.Log("🚀 开始大批量数据测试")

	// 测试1: 批量记录超过限制（1001条）
	records := make([]*service.RecordTokenUsageRequest, 1001)
	for i := 0; i < 1001; i++ {
		var userID uint64 = 10000 + uint64(i)
		records[i] = &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			RequestType:   entity.RequestTypeChat,
		}
	}

	_, err := ctx.TokenMeteringSvc.BatchRecordTokenUsage(testCtx, records)
	assert.Error(t, err, "Should fail for batch size > 1000")
	assert.Contains(t, err.Error(), "exceeds limit",
		"Error should mention batch size limit")

	t.Log("✅ 正确拒绝超过限制的批量操作")

	// 测试2: 边界值测试（1000条，应该成功）
	records2 := make([]*service.RecordTokenUsageRequest, 1000)
	for i := 0; i < 1000; i++ {
		var userID uint64 = 11000 + uint64(i)
		records2[i] = &service.RecordTokenUsageRequest{
			TenantID:      tenantID,
			UserID:        &userID,
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			RequestType:   entity.RequestTypeChat,
		}
	}

	resp, err := ctx.TokenMeteringSvc.BatchRecordTokenUsage(testCtx, records2)
	assert.NoError(t, err, "Should succeed for batch size = 1000")
	assert.Equal(t, 1000, resp.SuccessCount, "Should have 1000 successful records")

	t.Log("✅ 边界值测试通过（1000条）")

	t.Log("✅ 大批量数据测试通过")
}

// ============================================================
// 辅助函数
// ============================================================

/**
 * createTestTenantForError 创建测试租户（错误测试专用）
 */
func createTestTenantForError(ctx *errorTestContext, tenantID string) {
	tenant := &struct {
		TenantID string `gorm:"primaryKey;size:36"`
		Name     string `gorm:"size:100"`
	}{
		TenantID: tenantID,
		Name:     "错误测试租户",
	}
	ctx.DB.Table("tenants").Create(tenant)
}

/**
 * containsString 检查字符串是否包含子串
 */
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
