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

package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
)

// routingRuleRepository 路由规则仓储实现
type routingRuleRepository struct {
	db *gorm.DB
}

// NewRoutingRuleRepository 创建路由规则仓储实例
func NewRoutingRuleRepository(db *gorm.DB) RoutingRuleRepository {
	return &routingRuleRepository{db: db}
}

// Create 创建路由规则
func (r *routingRuleRepository) Create(ctx context.Context, rule *entity.RoutingRule) error {
	rule.CreatedAt = time.Now().UnixMilli()
	rule.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(rule).Error
}

// GetByID 根据ID获取路由规则
func (r *routingRuleRepository) GetByID(ctx context.Context, ruleID string) (*entity.RoutingRule, error) {
	var rule entity.RoutingRule
	err := r.db.WithContext(ctx).
		Where("rule_id = ?", ruleID).
		First(&rule).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update 更新路由规则
func (r *routingRuleRepository) Update(ctx context.Context, rule *entity.RoutingRule) error {
	rule.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.RoutingRule{}).
		Where("rule_id = ?", rule.RuleID).
		Updates(rule).Error
}

// Delete 删除路由规则
func (r *routingRuleRepository) Delete(ctx context.Context, ruleID string) error {
	return r.db.WithContext(ctx).
		Where("rule_id = ?", ruleID).
		Delete(&entity.RoutingRule{}).Error
}

// GetActiveRulesByTenant 获取租户的所有启用规则（按优先级排序）
func (r *routingRuleRepository) GetActiveRulesByTenant(ctx context.Context, tenantID string) ([]*entity.RoutingRule, error) {
	var rules []*entity.RoutingRule
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("is_active = ?", true).
		Order("priority DESC").
		Find(&rules).Error

	return rules, err
}

// List 分页查询规则列表
func (r *routingRuleRepository) List(ctx context.Context, filter *RoutingRuleFilter) ([]*entity.RoutingRule, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.RoutingRule{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.RuleType != nil {
		query = query.Where("rule_type = ?", *filter.RuleType)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询数据
	var rules []*entity.RoutingRule
	err := query.
		Order("priority DESC").
		Order("created_at DESC").
		Limit(pageSize).
		Find(&rules).Error

	return rules, total, err
}

// routingLogRepository 路由日志仓储实现
type routingLogRepository struct {
	db *gorm.DB
}

// NewRoutingLogRepository 创建路由日志仓储实例
func NewRoutingLogRepository(db *gorm.DB) RoutingLogRepository {
	return &routingLogRepository{db: db}
}

// Create 创建路由日志
func (r *routingLogRepository) Create(ctx context.Context, log *entity.RoutingLog) error {
	log.CreatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(log).Error
}

// List 查询路由日志
func (r *routingLogRepository) List(ctx context.Context, filter *RoutingLogFilter) ([]*entity.RoutingLog, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.RoutingLog{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询数据
	var logs []*entity.RoutingLog
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// DeleteOldLogs 删除旧日志
func (r *routingLogRepository) DeleteOldLogs(ctx context.Context, beforeDate int64) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", beforeDate).
		Delete(&entity.RoutingLog{}).Error
}
