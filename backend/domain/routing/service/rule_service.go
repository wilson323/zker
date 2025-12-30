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

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
)

// RuleService 路由规则服务
type RuleService struct {
	ruleRepo repository.RoutingRuleRepository
}

// NewRuleService 创建规则服务
func NewRuleService(ruleRepo repository.RoutingRuleRepository) *RuleService {
	return &RuleService{
		ruleRepo: ruleRepo,
	}
}

// GetActiveRules 获取活跃规则（按优先级排序）
func (s *RuleService) GetActiveRules(ctx context.Context, tenantID string) ([]*entity.RoutingRule, error) {
	rules, err := s.ruleRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 过滤活跃规则
	activeRules := make([]*entity.RoutingRule, 0)
	for _, rule := range rules {
		if rule.IsActive {
			activeRules = append(activeRules, rule)
		}
	}

	return activeRules, nil
}

// GetRulesByType 按类型获取规则
func (s *RuleService) GetRulesByType(
	ctx context.Context,
	tenantID string,
	ruleType entity.RuleType,
) ([]*entity.RoutingRule, error) {
	rules, err := s.ruleRepo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	filtered := make([]*entity.RoutingRule, 0)
	for _, rule := range rules {
		if rule.RuleType == ruleType {
			filtered = append(filtered, rule)
		}
	}

	return filtered, nil
}

// EnableRule 启用规则
func (s *RuleService) EnableRule(ctx context.Context, ruleID string) error {
	return s.ruleRepo.UpdateActive(ctx, ruleID, true)
}

// DisableRule 禁用规则
func (s *RuleService) DisableRule(ctx context.Context, ruleID string) error {
	return s.ruleRepo.UpdateActive(ctx, ruleID, false)
}

// UpdatePriority 更新规则优先级
func (s *RuleService) UpdatePriority(ctx context.Context, ruleID string, priority int) error {
	return s.ruleRepo.UpdatePriority(ctx, ruleID, priority)
}

// ValidateRule 验证规则
func (s *RuleService) ValidateRule(rule *entity.RoutingRule) error {
	// TODO: 实现规则验证逻辑
	// 1. 验证条件格式（JSON）
	// 2. 验证目标ID（BotID或WorkflowID至少一个）
	// 3. 验证优先级范围
	// 4. 验证规则类型
	return nil
}
