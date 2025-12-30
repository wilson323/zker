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
	"fmt"
	"math"

	"github.com/shopspring/decimal"
)

// ============================================================
// 定价引擎 - 基于设计文档: 23-租户计费系统_TokenMetering补充_完整版.md
// ============================================================

// PricingEngine 定价引擎
type PricingEngine struct {
	// 模型价格配置 (每1K tokens价格, CNY)
	prices map[string]ModelPrice
}

// ModelPrice 模型价格
type ModelPrice struct {
	InputPrice  decimal.Decimal // 输入价格（每1K tokens，CNY）
	OutputPrice decimal.Decimal // 输出价格（每1K tokens，CNY）
}

// CostBreakdown 成本明细
type CostBreakdown struct {
	UnitPrice  decimal.Decimal // 平均单价（每1K tokens）
	InputCost  decimal.Decimal // 输入成本
	OutputCost decimal.Decimal // 输出成本
	TotalCost  decimal.Decimal // 总成本
}

// NewPricingEngine 创建定价引擎实例
func NewPricingEngine() *PricingEngine {
	return &PricingEngine{
		prices: map[string]ModelPrice{
			// OpenAI 模型定价
			"openai/gpt-4": {
				InputPrice:  decimal.NewFromFloat(0.03),  // ¥0.03/1K tokens
				OutputPrice: decimal.NewFromFloat(0.06),  // ¥0.06/1K tokens
			},
			"openai/gpt-4-32k": {
				InputPrice:  decimal.NewFromFloat(0.06),
				OutputPrice: decimal.NewFromFloat(0.12),
			},
			"openai/gpt-3.5-turbo": {
				InputPrice:  decimal.NewFromFloat(0.003), // ¥0.003/1K tokens
				OutputPrice: decimal.NewFromFloat(0.006), // ¥0.006/1K tokens
			},
			"openai/gpt-3.5-turbo-16k": {
				InputPrice:  decimal.NewFromFloat(0.004),
				OutputPrice: decimal.NewFromFloat(0.008),
			},

			// Anthropic 模型定价
			"anthropic/claude-3-opus": {
				InputPrice:  decimal.NewFromFloat(0.09),
				OutputPrice: decimal.NewFromFloat(0.27),
			},
			"anthropic/claude-3-sonnet": {
				InputPrice:  decimal.NewFromFloat(0.015),
				OutputPrice: decimal.NewFromFloat(0.045),
			},
			"anthropic/claude-3-haiku": {
				InputPrice:  decimal.NewFromFloat(0.0025),
				OutputPrice: decimal.NewFromFloat(0.0125),
			},

			// 通义千问 模型定价
			"qwen/qwen-max": {
				InputPrice:  decimal.NewFromFloat(0.02),
				OutputPrice: decimal.NewFromFloat(0.06),
			},
			"qwen/qwen-plus": {
				InputPrice:  decimal.NewFromFloat(0.004),
				OutputPrice: decimal.NewFromFloat(0.012),
			},
			"qwen/qwen-turbo": {
				InputPrice:  decimal.NewFromFloat(0.001),
				OutputPrice: decimal.NewFromFloat(0.002),
			},

			// 百度文心 模型定价
			"baidu/ernie-bot-4": {
				InputPrice:  decimal.NewFromFloat(0.012),
				OutputPrice: decimal.NewFromFloat(0.012),
			},
			"baidu/ernie-bot-turbo": {
				InputPrice:  decimal.NewFromFloat(0.008),
				OutputPrice: decimal.NewFromFloat(0.008),
			},

			// 智谱 ChatGLM 模型定价
			"zhipu/chatglm-turbo": {
				InputPrice:  decimal.NewFromFloat(0.005),
				OutputPrice: decimal.NewFromFloat(0.005),
			},
			"zhipu/chatglm-pro": {
				InputPrice:  decimal.NewFromFloat(0.01),
				OutputPrice: decimal.NewFromFloat(0.01),
			},
		},
	}
}

// CalculateCost 计算成本
// 参数:
//   - provider: 模型提供商 (openai, anthropic, qwen, baidu, zhipu)
//   - model: 模型名称 (gpt-4, claude-3-opus, qwen-max等)
//   - inputTokens: 输入Token数量
//   - outputTokens: 输出Token数量
// 返回: 成本明细
func (e *PricingEngine) CalculateCost(
	provider, model string,
	inputTokens, outputTokens int,
) *CostBreakdown {
	// 1. 构建模型key
	key := fmt.Sprintf("%s/%s", provider, model)

	// 2. 查找价格配置
	price, ok := e.prices[key]
	if !ok {
		// 如果模型不存在，使用默认价格
		price = ModelPrice{
			InputPrice:  decimal.NewFromFloat(0.01),
			OutputPrice: decimal.NewFromFloat(0.02),
		}
	}

	// 3. 计算输入成本 (Token数 / 1000 * 单价)
	inputCost := decimal.NewFromInt(int64(inputTokens)).
		Div(decimal.NewFromInt(1000)).
		Mul(price.InputPrice)

	// 4. 计算输出成本
	outputCost := decimal.NewFromInt(int64(outputTokens)).
		Div(decimal.NewFromInt(1000)).
		Mul(price.OutputPrice)

	// 5. 计算总成本
	totalCost := inputCost.Add(outputCost)

	// 6. 计算平均单价（每1K tokens的平均价格）
	totalTokens := inputTokens + outputTokens
	var unitPrice decimal.Decimal
	if totalTokens > 0 {
		avgPrice := decimal.NewFromInt(int64(totalTokens)).
			Div(decimal.NewFromInt(1000))
		unitPrice = totalCost.Div(avgPrice)
	} else {
		unitPrice = decimal.Zero
	}

	return &CostBreakdown{
		UnitPrice:  unitPrice,
		InputCost:  inputCost,
		OutputCost: outputCost,
		TotalCost:  totalCost,
	}
}

// GetModelPrice 获取模型价格配置
func (e *PricingEngine) GetModelPrice(provider, model string) (*ModelPrice, error) {
	key := fmt.Sprintf("%s/%s", provider, model)
	price, ok := e.prices[key]
	if !ok {
		return nil, fmt.Errorf("model price not found: %s", key)
	}
	return &price, nil
}

// ListModels 列出所有支持的模型
func (e *PricingEngine) ListModels() []string {
	models := make([]string, 0, len(e.prices))
	for model := range e.prices {
		models = append(models, model)
	}
	return models
}

// EstimateCost 估算成本（用于预算规划）
// 参数:
//   - provider: 模型提供商
//   - model: 模型名称
//   - estimatedDailyTokens: 预估每日Token使用量
//   - days: 天数
// 返回: 预估总成本
func (e *PricingEngine) EstimateCost(
	provider, model string,
	estimatedDailyTokens int,
	days int,
) decimal.Decimal {
	// 假设输入输出Token比例为 7:3（根据行业经验）
	inputTokens := int(float64(estimatedDailyTokens) * 0.7)
	outputTokens := int(float64(estimatedDailyTokens) * 0.3)

	// 计算单日成本
	dailyCost := e.CalculateCost(provider, model, inputTokens, outputTokens).TotalCost

	// 计算总成本
	totalCost := dailyCost.Mul(decimal.NewFromInt(int64(days)))

	return totalCost
}

// CalculateSavings 计算切换模型后的节省金额
// 参数:
//   - currentProvider, currentModel: 当前模型
//   - newProvider, newModel: 新模型
//   - inputTokens, outputTokens: Token使用量
// 返回: 节省的金额和节省百分比
func (e *PricingEngine) CalculateSavings(
	currentProvider, currentModel string,
	newProvider, newModel string,
	inputTokens, outputTokens int,
) (saving decimal.Decimal, savingPercent float64) {
	// 1. 计算当前模型成本
	currentCost := e.CalculateCost(currentProvider, currentModel, inputTokens, outputTokens)

	// 2. 计算新模型成本
	newCost := e.CalculateCost(newProvider, newModel, inputTokens, outputTokens)

	// 3. 计算节省金额
	savingAmount := currentCost.TotalCost.Sub(newCost.TotalCost)

	// 4. 计算节省百分比
	var percent float64
	if currentCost.TotalCost.IsPositive() {
		percentStr := savingAmount.Div(currentCost.TotalCost).
			Mul(decimal.NewFromInt(100)).
			String()
		fmt.Sscanf(percentStr, "%f", &percent)
	}

	return savingAmount, percent
}

// GetCheaperModel 获取更便宜的替代模型
// 参数:
//   - currentProvider, currentModel: 当前模型
//   - maxCostIncreasePercent: 最大成本增加百分比（负数表示只能更便宜）
// 返回: 推荐的替代模型
func (e *PricingEngine) GetCheaperModel(
	currentProvider, currentModel string,
	maxCostIncreasePercent float64,
) []string {
	currentKey := fmt.Sprintf("%s/%s", currentProvider, currentModel)
	currentPrice, ok := e.prices[currentKey]
	if !ok {
		return nil
	}

	// 计算当前平均价格
	currentAvgPrice := currentPrice.InputPrice.Add(currentPrice.OutputPrice).
		Div(decimal.NewFromInt(2))

	cheaperModels := make([]string, 0)

	for key, price := range e.prices {
		// 跳过当前模型
		if key == currentKey {
			continue
		}

		// 计算新模型的平均价格
		newAvgPrice := price.InputPrice.Add(price.OutputPrice).
			Div(decimal.NewFromInt(2))

		// 计算价格变化百分比
		priceChangePercent := newAvgPrice.Sub(currentAvgPrice).
			Div(currentAvgPrice).
			Mul(decimal.NewFromInt(100))

		// 检查是否符合条件
		priceChangeFloat, _ := parseFloat(priceChangePercent.String())
		if priceChangeFloat <= maxCostIncreasePercent {
			cheaperModels = append(cheaperModels, key)
		}
	}

	return cheaperModels
}

// parseFloat 安全地将字符串转换为float64
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// roundFloat 保留指定小数位数
func roundFloat(f float64, places int) float64 {
	shift := math.Pow10(places)
	return math.Round(f*shift) / shift
}
