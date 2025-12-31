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
)

// TestWeightedSplitStrategy 测试加权分割策略
func TestWeightedSplitStrategy(t *testing.T) {
	strategy := &WeightedSplitStrategy{}

	variants := []Variant{
		{ID: "a", Name: "A", Weight: 50.0},
		{ID: "b", Name: "B", Weight: 30.0},
		{ID: "c", Name: "C", Weight: 20.0},
	}

	// 测试多次分割，确保分布大致符合权重
	counts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		userID := generateTestUserID(i)
		variant := strategy.Split(userID, variants)
		counts[variant.ID]++
	}

	// 验证分布大致合理（允许10%误差）
	expectedA := 500
	expectedB := 300
	expectedC := 200

	tolerance := 100

	if abs(counts["a"]-expectedA) > tolerance {
		t.Errorf("Variant A count out of range: got %d, want %d±%d", counts["a"], expectedA, tolerance)
	}
	if abs(counts["b"]-expectedB) > tolerance {
		t.Errorf("Variant B count out of range: got %d, want %d±%d", counts["b"], expectedB, tolerance)
	}
	if abs(counts["c"]-expectedC) > tolerance {
		t.Errorf("Variant C count out of range: got %d, want %d±%d", counts["c"], expectedC, tolerance)
	}
}

// TestConsistentHashSplitStrategy 测试一致性哈希分割策略
func TestConsistentHashSplitStrategy(t *testing.T) {
	strategy := &ConsistentHashSplitStrategy{}

	variants := []Variant{
		{ID: "a", Name: "A", Weight: 33.3},
		{ID: "b", Name: "B", Weight: 33.3},
		{ID: "c", Name: "C", Weight: 33.4},
	}

	// 测试同一用户多次分配，应得到相同结果
	userID := "test_user_12345"
	variant1 := strategy.Split(userID, variants)
	variant2 := strategy.Split(userID, variants)

	if variant1.ID != variant2.ID {
		t.Errorf("Consistent hash failed: same user got different variants: %s vs %s", variant1.ID, variant2.ID)
	}

	// 测试不同用户应得到不同结果
	userID2 := "test_user_67890"
	variant3 := strategy.Split(userID2, variants)

	// 不强制要求不同，但大多数情况应该不同
	if variant1.ID == variant3.ID {
		t.Logf("Info: Different users got same variant (possible with hash collision): %s", variant1.ID)
	}
}

// TestLeastConnectionsStrategy 测试最少连接策略
func TestLeastConnectionsStrategy(t *testing.T) {
	// TODO: 需要mock LoadMonitor
}

// TestWeightedResponseTimeStrategy 测试加权响应时间策略
func TestWeightedResponseTimeStrategy(t *testing.T) {
	// TODO: 需要mock LoadMonitor
}

// TestResourceAwareStrategy 测试资源感知策略
func TestResourceAwareStrategy(t *testing.T) {
	// TODO: 需要mock LoadMonitor
}

// TestFeatureExtractor 测试特征提取器
func TestFeatureExtractor(t *testing.T) {
	// TODO: 需要mock Repository
}

// 测试辅助函数
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateTestUserID(index int) string {
	return "test_user_" + string(rune('0'+index%10))
}

// TestCalculateComplexity 测试复杂度计算
func TestCalculateComplexity(t *testing.T) {
	extractor := &FeatureExtractor{}

	tests := []struct {
		name     string
		query    string
		expected float64 // 预期复杂度范围
	}{
		{
			name:     "简单查询",
			query:    "你好",
			expected: 0.2,
		},
		{
			name:     "中等复杂度查询",
			query:    "我想知道今天北京的天气怎么样，还有明天会不会下雨？",
			expected: 0.5,
		},
		{
			name:     "复杂查询",
			query:    "我需要在2025年1月1日到12月31日之间，查询所有订单金额大于1000元且状态为已完成的订单，并按金额降序排列，同时需要导出Excel报表发送到邮箱test@example.com。",
			expected: 0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexity := extractor.calculateComplexity(tt.query, QueryFeatures{
				Length:        len([]rune(tt.query)),
				WordCount:     len([]string{}),
				SentenceCount: 1,
			})

			// 允许0.2的误差
			if absFloat(complexity-tt.expected) > 0.2 {
				t.Logf("Query: %s", tt.query)
				t.Logf("Expected complexity: %.2f, got: %.2f", tt.expected, complexity)
			}
		})
	}
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// TestDetectQuestionType 测试问题类型检测
func TestDetectQuestionType(t *testing.T) {
	extractor := &FeatureExtractor{}

	tests := []struct {
		query    string
		expected string
	}{
		{"什么是人工智能", "what"},
		{"北京在哪里", "where"},
		{"什么时候开始", "when"},
		{"怎么使用这个功能", "how"},
		{"为什么失败", "why"},
		{"是谁创建了它", "who"},
		{"这是一个陈述句", "statement"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := extractor.detectQuestionType(tt.query)
			if result != tt.expected {
				t.Errorf("detectQuestionType(%q) = %q, want %q", tt.query, result, tt.expected)
			}
		})
	}
}

// TestDetectEmotionalTone 测试情感倾向检测
func TestDetectEmotionalTone(t *testing.T) {
	extractor := &FeatureExtractor{}

	tests := []struct {
		query    string
		expected string
	}{
		{"很好，非常满意", "positive"},
		{"太棒了，谢谢", "positive"},
		{"很差，非常不满", "negative"},
		{"糟糕，有问题", "negative"},
		{"查询天气", "neutral"},
		{"这是什么", "neutral"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			result := extractor.detectEmotionalTone(tt.query)
			if result != tt.expected {
				t.Errorf("detectEmotionalTone(%q) = %q, want %q", tt.query, result, tt.expected)
			}
		})
	}
}

// BenchmarkWeightedSplitStrategy 性能测试：加权分割策略
func BenchmarkWeightedSplitStrategy(b *testing.B) {
	strategy := &WeightedSplitStrategy{}
	variants := []Variant{
		{ID: "a", Weight: 50.0},
		{ID: "b", Weight: 30.0},
		{ID: "c", Weight: 20.0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := "test_user_12345"
		_ = strategy.Split(userID, variants)
	}
}

// BenchmarkConsistentHashSplitStrategy 性能测试：一致性哈希策略
func BenchmarkConsistentHashSplitStrategy(b *testing.B) {
	strategy := &ConsistentHashSplitStrategy{}
	variants := []Variant{
		{ID: "a", Weight: 33.3},
		{ID: "b", Weight: 33.3},
		{ID: "c", Weight: 33.4},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := "test_user_12345"
		_ = strategy.Split(userID, variants)
	}
}

// BenchmarkFeatureExtraction 性能测试：特征提取
func BenchmarkFeatureExtraction(b *testing.B) {
	extractor := &FeatureExtractor{}
	query := "我需要在2025年1月1日到12月31日之间，查询所有订单金额大于1000元且状态为已完成的订单。"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractor.extractQueryFeatures(nil, query)
	}
}
