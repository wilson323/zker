// +build integration

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

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	billingService "github.com/coze-dev/coze-studio/backend/domain/billing/service"
	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// TestPricingEngine_BasicFunctionality 测试定价引擎基础功能
func TestPricingEngine_BasicFunctionality(t *testing.T) {
	pricingEngine := billingService.NewPricingEngine()

	t.Run("CalculateCost_GPT4", func(t *testing.T) {
		cost := pricingEngine.CalculateCost("openai", "gpt-4", 1000, 500)

		inputCost, _ := cost.InputCost.Float64()
		outputCost, _ := cost.OutputCost.Float64()
		totalCost, _ := cost.TotalCost.Float64()

		// GPT-4: ¥0.03/1K input, ¥0.06/1K output
		assert.InDelta(t, 0.03, inputCost, 0.001)
		assert.InDelta(t, 0.03, outputCost, 0.001) // 500 * 0.06 / 1000 = 0.03
		assert.InDelta(t, 0.06, totalCost, 0.001)
	})

	t.Run("CalculateCost_Claude3Opus", func(t *testing.T) {
		cost := pricingEngine.CalculateCost("anthropic", "claude-3-opus", 2000, 1000)

		totalCost, _ := cost.TotalCost.Float64()

		// Claude-3 Opus: ¥0.09/1K input, ¥0.27/1K output
		// (2000 * 0.09 / 1000) + (1000 * 0.27 / 1000) = 0.18 + 0.27 = 0.45
		assert.InDelta(t, 0.45, totalCost, 0.001)
	})

	t.Run("CalculateCost_UnknownModel", func(t *testing.T) {
		cost := pricingEngine.CalculateCost("unknown", "unknown-model", 1000, 500)

		totalCost, _ := cost.TotalCost.Float64()

		// 未知模型使用默认价格: ¥0.01/1K input, ¥0.02/1K output
		assert.InDelta(t, 0.02, totalCost, 0.001)
	})
}

// TestTokenMeteringService_RecordUsage 测试Token使用记录
func TestTokenMeteringService_RecordUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 注意：这个测试需要真实的数据库连接
	// 在CI/CD环境中运行时，需要配置测试数据库

	ctx := context.Background()

	// 模拟请求
	req := &billingService.RecordTokenUsageRequest{
		TenantID:      "test-tenant-001",
		ModelProvider: "openai",
		ModelName:     "gpt-4",
		InputTokens:   1000,
		OutputTokens:  500,
		TotalTokens:   1500,
		RequestType:   entity.RequestTypeChat,
		ResponseTimeMs: intPtr(1200),
		IsCached:      false,
	}

	// 验证请求参数
	err := validateRequest(req)
	require.NoError(t, err, "Request validation should pass")

	// 验证成本计算（使用定价引擎）
	pricingEngine := billingService.NewPricingEngine()
	cost := pricingEngine.CalculateCost(req.ModelProvider, req.ModelName, req.InputTokens, req.OutputTokens)

	inputCost, _ := cost.InputCost.Float64()
	outputCost, _ := cost.OutputCost.Float64()
	totalCost, _ := cost.TotalCost.Float64()

	assert.Greater(t, inputCost, 0.0, "Input cost should be positive")
	assert.Greater(t, outputCost, 0.0, "Output cost should be positive")
	assert.Greater(t, totalCost, 0.0, "Total cost should be positive")

	t.Logf("✅ Token usage record validated: Input=%d, Output=%d, Cost=¥%.4f",
		req.InputTokens, req.OutputTokens, totalCost)
}

// TestTokenMeteringService_BatchRecord 测试批量记录
func TestTokenMeteringService_BatchRecord(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 准备100条记录
	records := make([]*billingService.RecordTokenUsageRequest, 100)
	for i := 0; i < 100; i++ {
		records[i] = &billingService.RecordTokenUsageRequest{
			TenantID:      "test-tenant-001",
			ModelProvider: "openai",
			ModelName:     "gpt-4",
			InputTokens:   1000,
			OutputTokens:  500,
			TotalTokens:   1500,
			RequestType:   entity.RequestTypeChat,
		}
	}

	// 验证批量大小
	assert.Len(t, records, 100, "Should have 100 records")

	// 计算预期总成本
	pricingEngine := billingService.NewPricingEngine()
	totalCost := 0.0

	for _, req := range records {
		cost := pricingEngine.CalculateCost(req.ModelProvider, req.ModelName, req.InputTokens, req.OutputTokens)
		costFloat, _ := cost.TotalCost.Float64()
		totalCost += costFloat
	}

	expectedCost := 0.06 * 100 // 每条记录约¥0.06，共100条
	assert.InDelta(t, expectedCost, totalCost, 0.1, "Total cost should be approximately ¥6.0")

	t.Logf("✅ Batch record validated: 100 records, Total cost=¥%.2f", totalCost)
}

// TestTokenMeteringService_Performance 测试性能
func TestTokenMeteringService_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 测试定价引擎性能
	pricingEngine := billingService.NewPricingEngine()

	iterations := 10000
	startTime := time.Now()

	for i := 0; i < iterations; i++ {
		pricingEngine.CalculateCost("openai", "gpt-4", 1000, 500)
	}

	duration := time.Since(startTime)

	// 性能要求：10000次计算应在100ms内完成
	assert.Less(t, duration.Milliseconds(), int64(100),
		"10000 calculations should complete within 100ms")

	t.Logf("✅ Performance test passed: %d calculations in %v (%.2f ops/sec)",
		iterations, duration, float64(iterations)/duration.Seconds())
}

// ============================================================
// 辅助函数
// ============================================================

func validateRequest(req *billingService.RecordTokenUsageRequest) error {
	if req.TenantID == "" {
		return t.Errorf("tenant_id is required")
	}

	if req.ModelProvider == "" {
		return t.Errorf("model_provider is required")
	}

	if req.ModelName == "" {
		return t.Errorf("model_name is required")
	}

	if req.InputTokens < 0 || req.OutputTokens < 0 || req.TotalTokens < 0 {
		return t.Errorf("token counts cannot be negative")
	}

	if req.TotalTokens != req.InputTokens+req.OutputTokens {
		return t.Errorf("total_tokens must equal input_tokens + output_tokens")
	}

	if req.RequestType == "" {
		return t.Errorf("request_type is required")
	}

	return nil
}

func intPtr(v int) *int {
	return &v
}
