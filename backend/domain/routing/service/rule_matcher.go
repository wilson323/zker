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
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/coze-studio/backend/domain/routing/entity"
	"github.com/coze-studio/backend/domain/routing/repository"
)

// RuleBasedMatcher 规则匹配器
type RuleBasedMatcher struct {
	ruleRepo repository.RoutingRuleRepository
}

// NewRuleBasedMatcher 创建规则匹配器实例
func NewRuleBasedMatcher(ruleRepo repository.RoutingRuleRepository) *RuleBasedMatcher {
	return &RuleBasedMatcher{ruleRepo: ruleRepo}
}

// MatchInput 匹配输入
type MatchInput struct {
	UserInput  string                 `json:"user_input"`
	TenantID   string                 `json:"tenant_id"`
	Context    map[string]interface{} `json:"context"` // 额外上下文信息
}

// MatchOutput 匹配输出
type MatchOutput struct {
	BotID       string  `json:"bot_id,omitempty"`
	WorkflowID  string  `json:"workflow_id,omitempty"`
	Confidence  float64 `json:"confidence"`
	MatchType   string  `json:"match_type"`
	RuleID      string  `json:"rule_id"`
}

// Match 规则匹配
func (m *RuleBasedMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	// 1. 获取所有启用的规则，按优先级排序
	rules, err := m.ruleRepo.GetActiveRulesByTenant(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}

	results := make([]*MatchOutput, 0)

	// 2. 遍历规则进行匹配
	for _, rule := range rules {
		matched, confidence := m.matchRule(ctx, rule, input)
		if matched {
			results = append(results, &MatchOutput{
				BotID:      m.getStringValue(rule.TargetBotID),
				WorkflowID: m.getStringValue(rule.TargetWorkflowID),
				Confidence: confidence,
				MatchType:  string(rule.RuleType),
				RuleID:     rule.RuleID,
			})
		}
	}

	// 3. 按置信度和优先级排序
	sort.Slice(results, func(i, j int) bool {
		if results[i].Confidence != results[j].Confidence {
			return results[i].Confidence > results[j].Confidence
		}
		// 置信度相同时，按规则ID排序（保证稳定性）
		return results[i].RuleID > results[j].RuleID
	})

	// 4. 返回 Top-3
	if len(results) > 3 {
		results = results[:3]
	}

	return results, nil
}

// matchRule 匹配单个规则
func (m *RuleBasedMatcher) matchRule(ctx context.Context, rule *entity.RoutingRule, input *MatchInput) (bool, float64) {
	var condition map[string]interface{}
	if err := json.Unmarshal([]byte(rule.Condition), &condition); err != nil {
		return false, 0
	}

	switch rule.RuleType {
	case entity.RuleTypeKeyword:
		return m.matchKeyword(condition, input.UserInput)
	case entity.RuleTypeRegex:
		return m.matchRegex(condition, input.UserInput)
	case entity.RuleTypeIntent:
		return m.matchIntent(condition, input.Context)
	case entity.RuleTypeCategory:
		return m.matchCategory(condition, input.Context)
	}

	return false, 0
}

// matchKeyword 关键词匹配
func (m *RuleBasedMatcher) matchKeyword(condition map[string]interface{}, userInput string) (bool, float64) {
	keywords, ok := condition["keywords"].([]interface{})
	if !ok {
		return false, 0
	}

	matchCount := 0
	userInputLower := strings.ToLower(userInput)

	for _, kw := range keywords {
		keyword, ok := kw.(string)
		if !ok {
			continue
		}
		if strings.Contains(userInputLower, strings.ToLower(keyword)) {
			matchCount++
		}
	}

	if matchCount == 0 {
		return false, 0
	}

	// 置信度 = 匹配的关键词数量 / 总关键词数量
	confidence := float64(matchCount) / float64(len(keywords))
	return true, confidence
}

// matchRegex 正则匹配
func (m *RuleBasedMatcher) matchRegex(condition map[string]interface{}, userInput string) (bool, float64) {
	pattern, ok := condition["pattern"].(string)
	if !ok {
		return false, 0
	}

	matched, err := regexp.MatchString(pattern, userInput)
	if err != nil || !matched {
		return false, 0
	}

	return true, 1.0
}

// matchIntent 意图匹配
func (m *RuleBasedMatcher) matchIntent(condition map[string]interface{}, context map[string]interface{}) (bool, float64) {
	requiredIntent, ok := condition["intent"].(string)
	if !ok {
		return false, 0
	}

	if context == nil {
		return false, 0
	}

	currentIntent, ok := context["intent"].(string)
	if !ok {
		return false, 0
	}

	if currentIntent == requiredIntent {
		return true, 1.0
	}

	return false, 0
}

// matchCategory 分类匹配
func (m *RuleBasedMatcher) matchCategory(condition map[string]interface{}, context map[string]interface{}) (bool, float64) {
	requiredCategory, ok := condition["category"].(string)
	if !ok {
		return false, 0
	}

	if context == nil {
		return false, 0
	}

	currentCategory, ok := context["category"].(string)
	if !ok {
		return false, 0
	}

	if currentCategory == requiredCategory {
		return true, 1.0
	}

	return false, 0
}

// getStringValue 安全获取字符串指针的值
func (m *RuleBasedMatcher) getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
