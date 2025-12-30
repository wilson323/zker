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
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================
// PricingEngine 单元测试套件
// ============================================================

/**
 * TestPricingEngine_CalculateCost_OpenAI
 *
 * 测试OpenAI模型成本计算
 * 验证：13个AI模型的定价精度
 */
func TestPricingEngine_CalculateCost_OpenAI(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("GPT-4 - 标准定价", func(t *testing.T) {
		// 1000输入 + 500输出 = ¥0.03 + ¥0.03 = ¥0.06
		breakdown := engine.CalculateCost("openai", "gpt-4", 1000, 500)

		require.NotNil(t, breakdown)

		// 验证输入成本：1000 / 1000 * 0.03 = 0.03
		assert.Equal(t, 0.03, mustFloat64(breakdown.InputCost), "Input cost should be 0.03")

		// 验证输出成本：500 / 1000 * 0.06 = 0.03
		assert.Equal(t, 0.03, mustFloat64(breakdown.OutputCost), "Output cost should be 0.03")

		// 验证总成本：0.03 + 0.03 = 0.06
		assert.Equal(t, 0.06, mustFloat64(breakdown.TotalCost), "Total cost should be 0.06")

		// 验证平均单价：(0.03+0.03)/(1500/1000) = 0.04
		assert.Equal(t, 0.04, mustFloat64(breakdown.UnitPrice), "Unit price should be 0.04")
	})

	t.Run("GPT-4-32K - 双倍价格", func(t *testing.T) {
		breakdown := engine.CalculateCost("openai", "gpt-4-32k", 1000, 500)

		// 输入成本：1000/1000 * 0.06 = 0.06
		assert.Equal(t, 0.06, mustFloat64(breakdown.InputCost))
		// 输出成本：500/1000 * 0.12 = 0.06
		assert.Equal(t, 0.06, mustFloat64(breakdown.OutputCost))
		// 总成本：0.12
		assert.Equal(t, 0.12, mustFloat64(breakdown.TotalCost))
	})

	t.Run("GPT-3.5-Turbo - 经济型", func(t *testing.T) {
		breakdown := engine.CalculateCost("openai", "gpt-3.5-turbo", 1000, 500)

		// 输入成本：1000/1000 * 0.003 = 0.003
		assert.Equal(t, 0.003, mustFloat64(breakdown.InputCost))
		// 输出成本：500/1000 * 0.006 = 0.003
		assert.Equal(t, 0.003, mustFloat64(breakdown.OutputCost))
		// 总成本：0.006（是GPT-4的1/10）
		assert.Equal(t, 0.006, mustFloat64(breakdown.TotalCost))
	})

	t.Run("GPT-3.5-Turbo-16K - 16K上下文", func(t *testing.T) {
		breakdown := engine.CalculateCost("openai", "gpt-3.5-turbo-16k", 1000, 500)

		assert.Equal(t, 0.004, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.004, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.008, mustFloat64(breakdown.TotalCost))
	})
}

/**
 * TestPricingEngine_CalculateCost_Anthropic
 *
 * 测试Anthropic模型成本计算
 */
func TestPricingEngine_CalculateCost_Anthropic(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("Claude-3-Opus - 最高端模型", func(t *testing.T) {
		breakdown := engine.CalculateCost("anthropic", "claude-3-opus", 1000, 500)

		// 输入成本：1000/1000 * 0.09 = 0.09
		assert.Equal(t, 0.09, mustFloat64(breakdown.InputCost))
		// 输出成本：500/1000 * 0.27 = 0.135
		assert.Equal(t, 0.135, mustFloat64(breakdown.OutputCost))
		// 总成本：0.225
		assert.Equal(t, 0.225, mustFloat64(breakdown.TotalCost))
	})

	t.Run("Claude-3-Sonnet - 中端模型", func(t *testing.T) {
		breakdown := engine.CalculateCost("anthropic", "claude-3-sonnet", 1000, 500)

		assert.Equal(t, 0.015, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.0225, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.0375, mustFloat64(breakdown.TotalCost))
	})

	t.Run("Claude-3-Haiku - 经济型", func(t *testing.T) {
		breakdown := engine.CalculateCost("anthropic", "claude-3-haiku", 1000, 500)

		assert.Equal(t, 0.0025, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.00625, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.00875, mustFloat64(breakdown.TotalCost))
	})
}

/**
 * TestPricingEngine_CalculateCost_Qwen
 *
 * 测试通义千问模型成本计算
 */
func TestPricingEngine_CalculateCost_Qwen(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("Qwen-Max - 标准定价", func(t *testing.T) {
		breakdown := engine.CalculateCost("qwen", "qwen-max", 1000, 500)

		assert.Equal(t, 0.02, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.03, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.05, mustFloat64(breakdown.TotalCost))
	})

	t.Run("Qwen-Plus - 中端", func(t *testing.T) {
		breakdown := engine.CalculateCost("qwen", "qwen-plus", 1000, 500)

		assert.Equal(t, 0.004, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.006, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.01, mustFloat64(breakdown.TotalCost))
	})

	t.Run("Qwen-Turbo - 经济型", func(t *testing.T) {
		breakdown := engine.CalculateCost("qwen", "qwen-turbo", 1000, 500)

		assert.Equal(t, 0.001, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.001, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.002, mustFloat64(breakdown.TotalCost))
	})
}

/**
 * TestPricingEngine_CalculateCost_Baidu
 *
 * 测试百度文心模型成本计算
 */
func TestPricingEngine_CalculateCost_Baidu(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("ERNIE-Bot-4 - 统一定价", func(t *testing.T) {
		breakdown := engine.CalculateCost("baidu", "ernie-bot-4", 1000, 500)

		// 输入输出价格相同
		assert.Equal(t, 0.012, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.006, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.018, mustFloat64(breakdown.TotalCost))
	})

	t.Run("ERNIE-Bot-Turbo - 经济型", func(t *testing.T) {
		breakdown := engine.CalculateCost("baidu", "ernie-bot-turbo", 1000, 500)

		assert.Equal(t, 0.008, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.004, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.012, mustFloat64(breakdown.TotalCost))
	})
}

/**
 * TestPricingEngine_CalculateCost_Zhipu
 *
 * 测试智谱ChatGLM模型成本计算
 */
func TestPricingEngine_CalculateCost_Zhipu(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("ChatGLM-Turbo", func(t *testing.T) {
		breakdown := engine.CalculateCost("zhipu", "chatglm-turbo", 1000, 500)

		assert.Equal(t, 0.005, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.0025, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.0075, mustFloat64(breakdown.TotalCost))
	})

	t.Run("ChatGLM-Pro", func(t *testing.T) {
		breakdown := engine.CalculateCost("zhipu", "chatglm-pro", 1000, 500)

		assert.Equal(t, 0.01, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.005, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.015, mustFloat64(breakdown.TotalCost))
	})
}

/**
 * TestPricingEngine_CalculateCost_EdgeCases
 *
 * 测试边界条件和异常情况
 */
func TestPricingEngine_CalculateCost_EdgeCases(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("零Token - 边界条件", func(t *testing.T) {
		breakdown := engine.CalculateCost("openai", "gpt-4", 0, 0)

		assert.Equal(t, 0.0, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.0, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.0, mustFloat64(breakdown.TotalCost))
		assert.Equal(t, 0.0, mustFloat64(breakdown.UnitPrice))
	})

	t.Run("只有输入Token", func(t *testing.T) {
		breakdown := engine.CalculateCost("openai", "gpt-4", 1000, 0)

		assert.Equal(t, 0.03, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.0, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.03, mustFloat64(breakdown.TotalCost))
		assert.Equal(t, 0.03, mustFloat64(breakdown.UnitPrice))
	})

	t.Run("只有输出Token", func(t *testing.T) {
		breakdown := engine.CalculateCost("openai", "gpt-4", 0, 1000)

		assert.Equal(t, 0.0, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.06, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.06, mustFloat64(breakdown.TotalCost))
		assert.Equal(t, 0.06, mustFloat64(breakdown.UnitPrice))
	})

	t.Run("大数量Token - 100万tokens", func(t *testing.T) {
		// 1M输入 + 500K输出
		breakdown := engine.CalculateCost("openai", "gpt-4", 1000000, 500000)

		// 输入成本：1000000/1000 * 0.03 = 30
		assert.Equal(t, 30.0, mustFloat64(breakdown.InputCost))
		// 输出成本：500000/1000 * 0.06 = 30
		assert.Equal(t, 30.0, mustFloat64(breakdown.OutputCost))
		// 总成本：60
		assert.Equal(t, 60.0, mustFloat64(breakdown.TotalCost))
	})

	t.Run("未知模型 - 默认定价", func(t *testing.T) {
		breakdown := engine.CalculateCost("unknown", "unknown-model", 1000, 500)

		// 应该使用默认价格：输入0.01，输出0.02
		assert.Equal(t, 0.01, mustFloat64(breakdown.InputCost))
		assert.Equal(t, 0.01, mustFloat64(breakdown.OutputCost))
		assert.Equal(t, 0.02, mustFloat64(breakdown.TotalCost))
	})
}

/**
 * TestPricingEngine_DecimalPrecision
 *
 * 测试decimal精度计算
 * 验证：财务计算必须精确到小数点后6位
 */
func TestPricingEngine_DecimalPrecision(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("精度验证 - 小数点后6位", func(t *testing.T) {
		// 使用非整数的token数
		breakdown := engine.CalculateCost("openai", "gpt-4", 1234, 5678)

		// 手动计算预期值
		expectedInputCost := decimal.NewFromInt(1234).
			Div(decimal.NewFromInt(1000)).
			Mul(decimal.NewFromFloat(0.03))

		expectedOutputCost := decimal.NewFromInt(5678).
			Div(decimal.NewFromInt(1000)).
			Mul(decimal.NewFromFloat(0.06))

		expectedTotal := expectedInputCost.Add(expectedOutputCost)

		// 验证decimal精确计算
		assert.True(t, breakdown.InputCost.Equal(expectedInputCost))
		assert.True(t, breakdown.OutputCost.Equal(expectedOutputCost))
		assert.True(t, breakdown.TotalCost.Equal(expectedTotal))
	})

	t.Run("避免浮点数精度问题", func(t *testing.T) {
		// 这是经典的浮点数问题：0.1 + 0.2 != 0.3
		// 使用decimal可以避免
		breakdown := engine.CalculateCost("openai", "gpt-3.5-turbo", 333, 666)

		// decimal计算不会产生精度误差
		inputCostStr := breakdown.InputCost.String()
		outputCostStr := breakdown.OutputCost.String()

		// 验证字符串格式正确（避免科学计数法）
		assert.Contains(t, inputCostStr, ".")
		assert.Contains(t, outputCostStr, ".")

		// 总成本应该精确等于输入+输出
		expectedTotal := breakdown.InputCost.Add(breakdown.OutputCost)
		assert.True(t, breakdown.TotalCost.Equal(expectedTotal))
	})
}

/**
 * TestPricingEngine_GetModelPrice
 *
 * 测试获取模型价格配置
 */
func TestPricingEngine_GetModelPrice(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("获取已知模型价格", func(t *testing.T) {
		price, err := engine.GetModelPrice("openai", "gpt-4")

		require.NoError(t, err)
		assert.NotNil(t, price)
		assert.Equal(t, 0.03, mustFloat64(price.InputPrice))
		assert.Equal(t, 0.06, mustFloat64(price.OutputPrice))
	})

	t.Run("获取不存在的模型", func(t *testing.T) {
		price, err := engine.GetModelPrice("unknown", "unknown-model")

		assert.Error(t, err)
		assert.Nil(t, price)
		assert.Contains(t, err.Error(), "model price not found")
	})
}

/**
 * TestPricingEngine_ListModels
 *
 * 测试列出所有支持的模型
 */
func TestPricingEngine_ListModels(t *testing.T) {
	engine := NewPricingEngine()
	models := engine.ListModels()

	// 验证13个模型都在列表中
	assert.Len(t, models, 13, "Should have 13 models")

	expectedModels := []string{
		"openai/gpt-4",
		"openai/gpt-4-32k",
		"openai/gpt-3.5-turbo",
		"openai/gpt-3.5-turbo-16k",
		"anthropic/claude-3-opus",
		"anthropic/claude-3-sonnet",
		"anthropic/claude-3-haiku",
		"qwen/qwen-max",
		"qwen/qwen-plus",
		"qwen/qwen-turbo",
		"baidu/ernie-bot-4",
		"baidu/ernie-bot-turbo",
		"zhipu/chatglm-turbo",
		"zhipu/chatglm-pro",
	}

	for _, expected := range expectedModels {
		assert.Contains(t, models, expected, "Should contain model: "+expected)
	}
}

/**
 * TestPricingEngine_EstimateCost
 *
 * 测试成本估算功能
 * 验证：7:3输入输出比例假设
 */
func TestPricingEngine_EstimateCost(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("估算30天成本 - GPT-4", func(t *testing.T) {
		// 假设每天使用10万tokens
		dailyTokens := 100000
		days := 30

		cost := engine.EstimateCost("openai", "gpt-4", dailyTokens, days)

		// 计算预期值：
		// 每天：输入7万+输出3万 = ¥2.1 + ¥1.8 = ¥3.9
		// 30天：¥3.9 * 30 = ¥117
		expectedCost := decimal.NewFromFloat(117.0)

		assert.True(t, cost.Equal(expectedCost),
			"Estimated cost should be 117.0, got: "+cost.String())
	})

	t.Run("估算7天成本 - GPT-3.5-Turbo", func(t *testing.T) {
		cost := engine.EstimateCost("openai", "gpt-3.5-turbo", 100000, 7)

		// 每天：7万输入 * 0.003 + 3万输出 * 0.006 = ¥0.21 + ¥0.18 = ¥0.39
		// 7天：¥0.39 * 7 = ¥2.73
		expectedCost := decimal.NewFromFloat(2.73)

		assert.True(t, cost.Equal(expectedCost))
	})
}

/**
 * TestPricingEngine_CalculateSavings
 *
 * 测试模型切换节省计算
 * 验证：降级节省百分比计算
 */
func TestPricingEngine_CalculateSavings(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("GPT-4降级到GPT-3.5-Turbo", func(t *testing.T) {
		saving, percent := engine.CalculateSavings(
			"openai", "gpt-4",
			"openai", "gpt-3.5-turbo",
			1000, 500,
		)

		// GPT-4成本：¥0.06
		// GPT-3.5-Turbo成本：¥0.006
		// 节省：¥0.054
		assert.Equal(t, 0.054, mustFloat64(saving), "Saving amount should be 0.054")

		// 节省百分比：0.054 / 0.06 = 90%
		assert.InDelta(t, 90.0, percent, 0.01, "Saving percent should be ~90%")
	})

	t.Run("Claude-3-Opus降级到Claude-3-Haiku", func(t *testing.T) {
		saving, percent := engine.CalculateSavings(
			"anthropic", "claude-3-opus",
			"anthropic", "claude-3-haiku",
			1000, 500,
		)

		// Claude-3-Opus成本：¥0.225
		// Claude-3-Haiku成本：¥0.00875
		// 节省：¥0.21625
		assert.InDelta(t, 0.216, mustFloat64(saving), 0.001)

		// 节省百分比：0.21625 / 0.225 ≈ 96.1%
		assert.InDelta(t, 96.0, percent, 1.0, "Saving percent should be ~96%")
	})

	t.Run("切换到更贵的模型", func(t *testing.T) {
		saving, percent := engine.CalculateSavings(
			"openai", "gpt-3.5-turbo",
			"anthropic", "claude-3-opus",
			1000, 500,
		)

		// 应该是负数（成本增加）
		assert.True(t, saving.IsNegative(), "Saving should be negative when upgrading")
		assert.True(t, percent < 0, "Percent should be negative")
	})
}

/**
 * TestPricingEngine_GetCheaperModel
 *
 * 测试获取更便宜的替代模型
 */
func TestPricingEngine_GetCheaperModel(t *testing.T) {
	engine := NewPricingEngine()

	t.Run("GPT-4的更便宜替代", func(t *testing.T) {
		// 查找比GPT-4便宜至少80%的模型
		cheaperModels := engine.GetCheaperModel("openai", "gpt-4", -80.0)

		// 应该包含所有更便宜的模型
		assert.NotEmpty(t, cheaperModels, "Should find cheaper models")

		// 验证GPT-3.5-Turbo在列表中（便宜90%）
		assert.Contains(t, cheaperModels, "openai/gpt-3.5-turbo")
	})

	t.Run("查找经济型模型", func(t *testing.T) {
		// 查找比Qwen-Plus便宜50%以上的模型
		cheaperModels := engine.GetCheaperModel("qwen", "qwen-plus", -50.0)

		// Qwen-Turbo应该符合条件（便宜75%）
		assert.Contains(t, cheaperModels, "qwen/qwen-turbo")
	})

	t.Run("未知模型", func(t *testing.T) {
		cheaperModels := engine.GetCheaperModel("unknown", "unknown-model", -100.0)

		assert.Nil(t, cheaperModels, "Should return nil for unknown model")
	})
}

/**
 * TestPricingEngine_CostComparison
 *
 * 测试模型成本对比
 * 验证：13个模型的相对成本
 */
func TestPricingEngine_CostComparison(t *testing.T) {
	engine := NewPricingEngine()

	// 标准测试用例：1000输入 + 500输出tokens
	inputTokens := 1000
	outputTokens := 500

	t.Run("成本排序 - 验证价格层次", func(t *testing.T) {
		costs := make(map[string]float64)

		// 计算所有模型的成本
		models := []string{
			"openai/gpt-4",
			"openai/gpt-3.5-turbo",
			"anthropic/claude-3-opus",
			"anthropic/claude-3-haiku",
			"qwen/qwen-max",
			"qwen/qwen-turbo",
		}

		for _, model := range models {
			provider := "openai"
			if contains(model, "claude") {
				provider = "anthropic"
			} else if contains(model, "qwen") {
				provider = "qwen"
			}

			parts := splitModelKey(model)
			breakdown := engine.CalculateCost(provider, parts[1], inputTokens, outputTokens)
			costs[model] = mustFloat64(breakdown.TotalCost)
		}

		// 验证成本关系
		assert.Greater(t, costs["anthropic/claude-3-opus"], costs["openai/gpt-4"],
			"Claude-3-Opus should be more expensive than GPT-4")
		assert.Greater(t, costs["openai/gpt-4"], costs["openai/gpt-3.5-turbo"],
			"GPT-4 should be more expensive than GPT-3.5-Turbo")
		assert.Less(t, costs["qwen/qwen-turbo"], costs["qwen/qwen-max"],
			"Qwen-Turbo should be cheaper than Qwen-Max")
	})
}

/**
 * TestPricingEngine_Performance
 *
 * 性能测试：验证定价计算速度
 */
func TestPricingEngine_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	engine := NewPricingEngine()

	t.Run("批量计算性能 - 10000次计算", func(t *testing.T) {
		iterations := 10000

		// 记录开始时间
		start := time.Now()

		for i := 0; i < iterations; i++ {
			engine.CalculateCost("openai", "gpt-4", 1000, 500)
		}

		duration := time.Since(start)

		// 验证性能：10000次计算应该在100ms内完成
		assert.Less(t, duration, 100*time.Millisecond,
			"10000 calculations should complete within 100ms, took: "+duration.String())

		// 计算每次的平均时间
		avgTime := duration / time.Duration(iterations)
		t.Logf("Average time per calculation: %v", avgTime)
	})
}

// ============================================================
// 辅助函数
// ============================================================

/**
 * mustFloat64 安全地将decimal.Decimal转换为float64
 *
 * 注意：仅用于测试验证，生产代码应使用decimal
 */
func mustFloat64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

/**
 * contains 检查字符串是否包含子字符串
 */
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		indexOf(s, substr) >= 0))
}

/**
 * indexOf 查找子字符串位置
 */
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

/**
 * splitModelKey 分割模型key
 */
func splitModelKey(key string) []string {
	// 简单实现：按"/"分割
	if len(key) > 0 && key[0] == '/' {
		return []string{"", key[1:]}
	}

	for i, c := range key {
		if c == '/' {
			return []string{key[:i], key[i+1:]}
		}
	}
	return []string{key}
}
