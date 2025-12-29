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
	"sort"
)

// HybridIntentMatcher 混合意图匹配器
type HybridIntentMatcher struct {
	ruleMatcher       *RuleBasedMatcher
	similarityMatcher *SimilarityMatcher
	ruleWeight         float64 // 规则匹配权重
	similarityWeight   float64 // 相似度匹配权重
}

// NewHybridIntentMatcher 创建混合意图匹配器实例
func NewHybridIntentMatcher(
	ruleMatcher *RuleBasedMatcher,
	similarityMatcher *SimilarityMatcher,
) *HybridIntentMatcher {
	return &HybridIntentMatcher{
		ruleMatcher:       ruleMatcher,
		similarityMatcher: similarityMatcher,
		ruleWeight:         0.5, // 默认权重50%
		similarityWeight:   0.5, // 默认权重50%
	}
}

// SetWeights 设置匹配权重
func (m *HybridIntentMatcher) SetWeights(ruleWeight, similarityWeight float64) {
	total := ruleWeight + similarityWeight
	if total > 0 {
		m.ruleWeight = ruleWeight / total
		m.similarityWeight = similarityWeight / total
	}
}

// Match 混合意图匹配
func (m *HybridIntentMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	results := make([]*MatchOutput, 0)

	// 1. 规则匹配（权重默认50%）
	ruleResults, err := m.ruleMatcher.Match(ctx, input)
	if err == nil {
		for _, r := range ruleResults {
			r.Confidence *= m.ruleWeight
			results = append(results, r)
		}
	}

	// 2. 相似度匹配（权重默认50%）
	simResults, err := m.similarityMatcher.Match(ctx, input)
	if err == nil {
		for _, r := range simResults {
			r.Confidence *= m.similarityWeight
			results = append(results, r)
		}
	}

	// 3. 聚合结果（相同Bot的置信度累加）
	aggregated := m.aggregateResults(results)

	// 4. 按置信度排序
	sort.Slice(aggregated, func(i, j int) bool {
		return aggregated[i].Confidence > aggregated[j].Confidence
	})

	// 5. 返回 Top-3
	if len(aggregated) > 3 {
		aggregated = aggregated[:3]
	}

	return aggregated, nil
}

// aggregateResults 聚合结果
func (m *HybridIntentMatcher) aggregateResults(results []*MatchOutput) []*MatchOutput {
	aggMap := make(map[string]*MatchOutput)

	for _, r := range results {
		// 使用BotID或WorkflowID作为唯一标识
		key := r.BotID
		if key == "" {
			key = r.WorkflowID
		}

		if key == "" {
			continue // 跳过无效的结果
		}

		if existing, ok := aggMap[key]; ok {
			// 累加置信度
			existing.Confidence += r.Confidence
			// 合并匹配类型
			if existing.MatchType != "" && r.MatchType != "" {
				existing.MatchType += "," + r.MatchType
			}
		} else {
			// 新建条目
			aggMap[key] = r
		}
	}

	// 转换为切片
	aggregated := make([]*MatchOutput, 0, len(aggMap))
	for _, v := range aggMap {
		aggregated = append(aggregated, v)
	}

	return aggregated
}
