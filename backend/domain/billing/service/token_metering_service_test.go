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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// 测试用例
// ============================================================

// TestRecordTokenUsage_Success 测试成功记录Token使用
func TestRecordTokenUsage_Success(t *testing.T) {
	// 1. 准备测试数据
	ctx := context.Background()
	mockDB := &gorm.DB{} // 简化的mock

	mockLogRepo := new(MockTokenUsageLogRepository)
	mockSummaryRepo := new(MockTokenUsageSummaryRepository)
	mockAlertRepo := new(MockBudgetAlertRepository)
	mockBudgetRepo := new(MockBudgetSettingsRepository)
	mockNotifier := &NotificationService{}

	pricingEngine := NewPricingEngine()
	budgetAlertSvc := NewBudgetAlertService(mockBudgetRepo, mockSummaryRepo, mockAlertRepo, mockNotifier, nil)

	service := NewTokenMeteringService(pricingEngine, mockLogRepo, mockSummaryRepo, budgetAlertSvc, mockDB, nil)

	req := &RecordTokenUsageRequest{
		TenantID:       "tenant-123",
		UserID:         uint64Ptr(1001),
		BotID:          strPtr("bot-456"),
		ConversationID: strPtr("conv-789"),
		MessageID:      strPtr("msg-001"),
		ModelProvider:  "openai",
		ModelName:      "gpt-4",
		InputTokens:    1000,
		OutputTokens:   500,
		TotalTokens:    1500,
		RequestType:    entity.RequestTypeChat,
		ResponseTimeMs: intPtr(1200),
		IsCached:       false,
	}

	// 2. 设置Mock期望
	mockLogRepo.On("Create", ctx, mock.AnythingOfType("*entity.TokenUsageLog")).Return(nil).Run(func(args mock.Arguments) {
		log := args.Get(1).(*entity.TokenUsageLog)
		log.ID = 1 // 模拟数据库分配的ID
		log.CreatedAt = time.Now()
	})

	// 3. 执行测试
	resp, err := service.RecordTokenUsage(ctx, req)

	// 4. 验证结果
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.LogID)
	assert.Greater(t, resp.InputCost, 0.0)
	assert.Greater(t, resp.OutputCost, 0.0)
	assert.Greater(t, resp.TotalCost, 0.0)
	assert.Greater(t, resp.CreatedAt, int64(0))

	// 验证成本计算
	expectedInputCost := 0.03 * 1000 / 1000  // ¥0.03/1K tokens
	expectedOutputCost := 0.06 * 500 / 1000  // ¥0.06/1K tokens
	assert.InDelta(t, expectedInputCost, resp.InputCost, 0.001)
	assert.InDelta(t, expectedOutputCost, resp.OutputCost, 0.001)

	mockLogRepo.AssertExpectations(t)
}

// TestRecordTokenUsage_ValidationError 测试参数验证错误
func TestRecordTokenUsage_ValidationError(t *testing.T) {
	ctx := context.Background()
	mockDB := &gorm.DB{}

	mockLogRepo := new(MockTokenUsageLogRepository)
	mockSummaryRepo := new(MockTokenUsageSummaryRepository)
	mockAlertRepo := new(MockBudgetAlertRepository)
	mockBudgetRepo := new(MockBudgetSettingsRepository)
	mockNotifier := &NotificationService{}

	pricingEngine := NewPricingEngine()
	budgetAlertSvc := NewBudgetAlertService(mockBudgetRepo, mockSummaryRepo, mockAlertRepo, mockNotifier, nil)

	service := NewTokenMeteringService(pricingEngine, mockLogRepo, mockSummaryRepo, budgetAlertSvc, mockDB, nil)

	tests := []struct {
		name    string
		req     *RecordTokenUsageRequest
		errMsg  string
	}{
		{
			name: "缺少tenant_id",
			req: &RecordTokenUsageRequest{
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   entity.RequestTypeChat,
			},
			errMsg: "tenant_id is required",
		},
		{
			name: "Token数不匹配",
			req: &RecordTokenUsageRequest{
				TenantID:      "tenant-123",
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   2000, // 错误：应该是1500
				RequestType:   entity.RequestTypeChat,
			},
			errMsg: "total_tokens must equal",
		},
		{
			name: "无效的request_type",
			req: &RecordTokenUsageRequest{
				TenantID:      "tenant-123",
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   1000,
				OutputTokens:  500,
				TotalTokens:   1500,
				RequestType:   "invalid_type",
			},
			errMsg: "invalid request_type",
		},
		{
			name: "负数Token",
			req: &RecordTokenUsageRequest{
				TenantID:      "tenant-123",
				ModelProvider: "openai",
				ModelName:     "gpt-4",
				InputTokens:   -100,
				OutputTokens:  500,
				TotalTokens:   400,
				RequestType:   entity.RequestTypeChat,
			},
			errMsg: "cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.RecordTokenUsage(ctx, tt.req)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

// TestBatchRecordTokenUsage_Success 测试批量记录成功
func TestBatchRecordTokenUsage_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &gorm.DB{}

	mockLogRepo := new(MockTokenUsageLogRepository)
	mockSummaryRepo := new(MockTokenUsageSummaryRepository)
	mockAlertRepo := new(MockBudgetAlertRepository)
	mockBudgetRepo := new(MockBudgetSettingsRepository)
	mockNotifier := &NotificationService{}

	pricingEngine := NewPricingEngine()
	budgetAlertSvc := NewBudgetAlertService(mockBudgetRepo, mockSummaryRepo, mockAlertRepo, mockNotifier, nil)

	service := NewTokenMeteringService(pricingEngine, mockLogRepo, mockSummaryRepo, budgetAlertSvc, mockDB, nil)

	// 准备10条记录
	records := make([]*RecordTokenUsageRequest, 10)
	for i := 0; i < 10; i++ {
		records[i] = &RecordTokenUsageRequest{
			TenantID:      "tenant-123",
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			RequestType:   entity.RequestTypeChat,
		}
	}

	// 设置Mock期望
	mockLogRepo.On("BatchCreate", ctx, mock.AnythingOfType("[]*entity.TokenUsageLog")).Return(nil).Run(func(args mock.Arguments) {
		logs := args.Get(1).([]*entity.TokenUsageLog)
		for i, log := range logs {
			log.ID = uint64(i + 1)
		}
	})

	// 执行测试
	startTime := time.Now()
	resp, err := service.BatchRecordTokenUsage(ctx, records)
	duration := time.Since(startTime)

	// 验证结果
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 10, resp.SuccessCount)
	assert.Equal(t, 0, resp.FailedCount)
	assert.Len(t, resp.LogIDs, 10)
	assert.Greater(t, resp.TotalCost, 0.0)
	assert.Less(t, duration.Milliseconds(), int64(1000), "Batch recording should complete within 1 second")

	mockLogRepo.AssertExpectations(t)
}

// TestBatchRecordTokenUsage_ExceedLimit 测试批量超过限制
func TestBatchRecordTokenUsage_ExceedLimit(t *testing.T) {
	ctx := context.Background()
	mockDB := &gorm.DB{}

	mockLogRepo := new(MockTokenUsageLogRepository)
	mockSummaryRepo := new(MockTokenUsageSummaryRepository)
	mockAlertRepo := new(MockBudgetAlertRepository)
	mockBudgetRepo := new(MockBudgetSettingsRepository)
	mockNotifier := &NotificationService{}

	pricingEngine := NewPricingEngine()
	budgetAlertSvc := NewBudgetAlertService(mockBudgetRepo, mockSummaryRepo, mockAlertRepo, mockNotifier, nil)

	service := NewTokenMeteringService(pricingEngine, mockLogRepo, mockSummaryRepo, budgetAlertSvc, mockDB, nil)

	// 准备1001条记录（超过限制）
	records := make([]*RecordTokenUsageRequest, 1001)
	for i := 0; i < 1001; i++ {
		records[i] = &RecordTokenUsageRequest{
			TenantID:      "tenant-123",
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			RequestType:   entity.RequestTypeChat,
		}
	}

	// 执行测试
	_, err := service.BatchRecordTokenUsage(ctx, records)

	// 验证结果
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds limit")
}

// TestGetUsageStats_Success 测试获取使用统计成功
func TestGetUsageStats_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &gorm.DB{}

	mockLogRepo := new(MockTokenUsageLogRepository)
	mockSummaryRepo := new(MockTokenUsageSummaryRepository)
	mockAlertRepo := new(MockBudgetAlertRepository)
	mockBudgetRepo := new(MockBudgetSettingsRepository)
	mockNotifier := &NotificationService{}

	pricingEngine := NewPricingEngine()
	budgetAlertSvc := NewBudgetAlertService(mockBudgetRepo, mockSummaryRepo, mockAlertRepo, mockNotifier, nil)

	service := NewTokenMeteringService(pricingEngine, mockLogRepo, mockSummaryRepo, budgetAlertSvc, mockDB, nil)

	// 准备Mock数据
	summaries := []*entity.TokenUsageSummary{
		{
			TenantID:         "tenant-123",
			TotalInputTokens: 10000,
			TotalOutputTokens: 5000,
			TotalTokens:      15000,
			TotalCost:        1.5,
			TotalRequests:    100,
			CachedRequests:   20,
			AvgResponseTime:  float64Ptr(500.0),
		},
	}

	mockSummaryRepo.On("GetByDateRange", ctx, "tenant-123", (*string)(nil), mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(summaries, nil)

	// 执行测试
	filter := &UsageStatsFilter{}
	resp, err := service.GetUsageStats(ctx, "tenant-123", filter)

	// 验证结果
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, int64(10000), resp.TotalInputTokens)
	assert.Equal(t, int64(5000), resp.TotalOutputTokens)
	assert.Equal(t, int64(15000), resp.TotalTokens)
	assert.Equal(t, 1.5, resp.TotalCost)
	assert.Equal(t, int64(100), resp.TotalRequests)
	assert.Equal(t, int64(20), resp.CachedRequests)
	assert.Equal(t, 500.0, resp.AvgResponseTime)

	mockSummaryRepo.AssertExpectations(t)
}

// TestGetModelUsageStats_Success 测试获取模型统计成功
func TestGetModelUsageStats_Success(t *testing.T) {
	ctx := context.Background()
	mockDB := &gorm.DB{}

	mockLogRepo := new(MockTokenUsageLogRepository)
	mockSummaryRepo := new(MockTokenUsageSummaryRepository)
	mockAlertRepo := new(MockBudgetAlertRepository)
	mockBudgetRepo := new(MockBudgetSettingsRepository)
	mockNotifier := &NotificationService{}

	pricingEngine := NewPricingEngine()
	budgetAlertSvc := NewBudgetAlertService(mockBudgetRepo, mockSummaryRepo, mockAlertRepo, mockNotifier, nil)

	service := NewTokenMeteringService(pricingEngine, mockLogRepo, mockSummaryRepo, budgetAlertSvc, mockDB, nil)

	// 准备Mock数据
	logs := []*entity.TokenUsageLog{
		{
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			TotalCost:     0.09,
		},
		{
			ModelProvider: "anthropic",
			ModelName:     "claude-3-opus",
			InputTokens:   2000,
			OutputTokens:  1000,
			TotalTokens:   3000,
			TotalCost:     0.45,
		},
	}

	mockLogRepo.On("GetByDateRange", ctx, "tenant-123", "2024-01-01", "2024-01-07").
		Return(logs, nil)

	// 执行测试
	stats, err := service.GetModelUsageStats(ctx, "tenant-123", "2024-01-01", "2024-01-07")

	// 验证结果
	require.NoError(t, err)
	require.Len(t, stats, 2)

	// 验证第一个模型（GPT-4）
	gpt4Stat := stats[0]
	assert.Equal(t, "openai", gpt4Stat.ModelProvider)
	assert.Equal(t, "gpt-4", gpt4Stat.ModelName)
	assert.Equal(t, int64(1500), gpt4Stat.TotalTokens)
	assert.Equal(t, 0.09, gpt4Stat.TotalCost)
	assert.Equal(t, int64(1), gpt4Stat.TotalRequests)

	mockLogRepo.AssertExpectations(t)
}

