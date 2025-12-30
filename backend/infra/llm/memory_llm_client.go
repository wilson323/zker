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

package llm

import (
	"context"
)

// MemoryLLMClient 记忆引擎LLM客户端接口
type MemoryLLMClient interface {
	// Embed 文本向量化
	Embed(ctx context.Context, text string) ([]float32, error)

	// BatchEmbed 批量文本向量化
	BatchEmbed(ctx context.Context, texts []string) ([][]float32, error)

	// Summarize 文本摘要
	Summarize(ctx context.Context, text string, targetLength int) (*SummaryResult, error)

	// ExtractEntities 实体提取
	ExtractEntities(ctx context.Context, text string) (*ExtractEntityResult, error)

	// ExtractPreferences 偏好提取
	ExtractPreferences(ctx context.Context, text string) (*ExtractPreferenceResult, error)

	// ExtractKeyFacts 提取关键事实
	ExtractKeyFacts(ctx context.Context, text string) ([]string, error)
}

// SummaryResult 摘要结果
type SummaryResult struct {
	Text      string   `json:"text"`      // 摘要文本
	KeyFacts  []string `json:"key_facts"` // 关键事实
	TokensUsed int     `json:"tokens_used"` // 使用的token数
}

// ExtractEntityResult 实体提取结果
type ExtractEntityResult struct {
	Entities  []Entity `json:"entities"`  // 提取的实体
	Confidence float64  `json:"confidence"` // 整体置信度
}

// Entity 提取的实体
type Entity struct {
	EntityType string            `json:"entity_type"` // 实体类型
	EntityValue string           `json:"entity_value"` // 实体值
	Attributes  map[string]string `json:"attributes"`  // 实体属性
	Confidence  float64           `json:"confidence"`  // 置信度
}

// ExtractPreferenceResult 偏好提取结果
type ExtractPreferenceResult struct {
	Preferences []Preference `json:"preferences"` // 提取的偏好
	Confidence  float64      `json:"confidence"`  // 整体置信度
}

// Preference 用户偏好
type Preference struct {
	PreferenceType  string `json:"preference_type"`  // 偏好类型
	PreferenceValue string `json:"preference_value"` // 偏好值
	Source         string `json:"source"`          // 来源
	Confidence     float64 `json:"confidence"`      // 置信度
}

// MockMemoryLLMClient 模拟LLM客户端(用于开发测试)
type MockMemoryLLMClient struct{}

// NewMockMemoryLLMClient 创建模拟LLM客户端
func NewMockMemoryLLMClient() *MockMemoryLLMClient {
	return &MockMemoryLLMClient{}
}

// Embed 文本向量化
func (c *MockMemoryLLMClient) Embed(ctx context.Context, text string) ([]float32, error) {
	// 返回模拟的1536维向量
	vector := make([]float32, 1536)
	for i := range vector {
		vector[i] = 0.1
	}
	return vector, nil
}

// BatchEmbed 批量文本向量化
func (c *MockMemoryLLMClient) BatchEmbed(ctx context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for i := range texts {
		vectors[i], _ = c.Embed(ctx, texts[i])
	}
	return vectors, nil
}

// Summarize 文本摘要
func (c *MockMemoryLLMClient) Summarize(ctx context.Context, text string, targetLength int) (*SummaryResult, error) {
	// 返回模拟摘要
	summary := "这是对话的摘要：用户讨论了关于AI记忆引擎的话题。"
	return &SummaryResult{
		Text:      summary,
		KeyFacts:  []string{"用户讨论AI记忆引擎", "关注向量检索技术"},
		TokensUsed: 100,
	}, nil
}

// ExtractEntities 实体提取
func (c *MockMemoryLLMClient) ExtractEntities(ctx context.Context, text string) (*ExtractEntityResult, error) {
	// 返回模拟实体
	return &ExtractEntityResult{
		Entities: []Entity{
			{
				EntityType:  "PERSON",
				EntityValue: "用户",
				Attributes:  map[string]string{"role": "user"},
				Confidence:  0.9,
			},
		},
		Confidence: 0.9,
	}, nil
}

// ExtractPreferences 偏好提取
func (c *MockMemoryLLMClient) ExtractPreferences(ctx context.Context, text string) (*ExtractPreferenceResult, error) {
	// 返回模拟偏好
	return &ExtractPreferenceResult{
		Preferences: []Preference{
			{
				PreferenceType:  "communication_style",
				PreferenceValue: "简洁",
				Source:         "inferred",
				Confidence:     0.85,
			},
		},
		Confidence: 0.85,
	}, nil
}

// ExtractKeyFacts 提取关键事实
func (c *MockMemoryLLMClient) ExtractKeyFacts(ctx context.Context, text string) ([]string, error) {
	// 返回模拟关键事实
	return []string{
		"用户讨论AI记忆引擎",
		"关注向量检索技术",
	}, nil
}
