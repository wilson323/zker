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
	// 直接使用仓储的GetActiveRulesByTenant方法
	return s.ruleRepo.GetActiveRulesByTenant(ctx, tenantID)
}

// GetRulesByType 按类型获取规则
func (s *RuleService) GetRulesByType(
	ctx context.Context,
	tenantID string,
	ruleType entity.RuleType,
) ([]*entity.RoutingRule, error) {
	filter := &repository.RoutingRuleFilter{
		TenantID: tenantID,
		RuleType: &ruleType,
		IsActive: nil, // 获取所有，不过滤激活状态
	}
	rules, _, err := s.ruleRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// EnableRule 启用规则
func (s *RuleService) EnableRule(ctx context.Context, ruleID string) error {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}
	rule.IsActive = true
	return s.ruleRepo.Update(ctx, rule)
}

// DisableRule 禁用规则
func (s *RuleService) DisableRule(ctx context.Context, ruleID string) error {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}
	rule.IsActive = false
	return s.ruleRepo.Update(ctx, rule)
}

// UpdatePriority 更新规则优先级
func (s *RuleService) UpdatePriority(ctx context.Context, ruleID string, priority int) error {
	rule, err := s.ruleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}
	rule.Priority = priority
	return s.ruleRepo.Update(ctx, rule)
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

// CreateRule 创建路由规则
func (s *RuleService) CreateRule(ctx context.Context, rule *entity.RoutingRule) (*entity.RoutingRule, error) {
	if err := s.ValidateRule(rule); err != nil {
		return nil, err
	}
	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// GetRuleByID 根据ID获取规则
func (s *RuleService) GetRuleByID(ctx context.Context, ruleID string) (*entity.RoutingRule, error) {
	return s.ruleRepo.GetByID(ctx, ruleID)
}

// UpdateRule 更新路由规则
func (s *RuleService) UpdateRule(ctx context.Context, rule *entity.RoutingRule) (*entity.RoutingRule, error) {
	if err := s.ValidateRule(rule); err != nil {
		return nil, err
	}
	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return nil, err
	}
	return rule, nil
}

// DeleteRule 删除路由规则
func (s *RuleService) DeleteRule(ctx context.Context, ruleID string) error {
	return s.ruleRepo.Delete(ctx, ruleID)
}

// ListRulesRequest 列出规则请求
type ListRulesRequest struct {
	TenantID  string
	RuleType  *entity.RuleType
	IsActive  *bool
	PageSize  int
	PageToken string
}

// ListRules 列出路由规则
func (s *RuleService) ListRules(ctx context.Context, req *ListRulesRequest) ([]*entity.RoutingRule, int64, string, error) {
	// 构建过滤条件
	filter := &repository.RoutingRuleFilter{
		TenantID:  req.TenantID,
		RuleType:  req.RuleType,
		IsActive:  req.IsActive,
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	}

	// 调用仓储获取规则列表（只有3个返回值）
	rules, total, err := s.ruleRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, "", err
	}

	// 简化实现：暂时不返回分页token
	return rules, total, "", nil
}
