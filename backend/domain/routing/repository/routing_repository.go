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

	"github.com/coze-studio/backend/domain/routing/entity"
)

// RoutingRuleRepository 路由规则仓储接口
type RoutingRuleRepository interface {
	// Create 创建路由规则
	Create(ctx context.Context, rule *entity.RoutingRule) error

	// GetByID 根据ID获取路由规则
	GetByID(ctx context.Context, ruleID string) (*entity.RoutingRule, error)

	// Update 更新路由规则
	Update(ctx context.Context, rule *entity.RoutingRule) error

	// Delete 删除路由规则
	Delete(ctx context.Context, ruleID string) error

	// GetActiveRulesByTenant 获取租户的所有启用规则（按优先级排序）
	GetActiveRulesByTenant(ctx context.Context, tenantID string) ([]*entity.RoutingRule, error)

	// List 分页查询规则列表
	List(ctx context.Context, filter *RoutingRuleFilter) ([]*entity.RoutingRule, int64, error)
}

// RoutingRuleFilter 路由规则查询过滤器
type RoutingRuleFilter struct {
	TenantID string
	RuleType  *entity.RuleType
	IsActive  *bool
	PageToken string
	PageSize  int
}

// RoutingLogRepository 路由日志仓储接口
type RoutingLogRepository interface {
	// Create 创建路由日志
	Create(ctx context.Context, log *entity.RoutingLog) error

	// List 查询路由日志
	List(ctx context.Context, filter *RoutingLogFilter) ([]*entity.RoutingLog, int64, error)

	// DeleteOldLogs 删除旧日志
	DeleteOldLogs(ctx context.Context, beforeDate int64) error
}

// RoutingLogFilter 路由日志查询过滤器
type RoutingLogFilter struct {
	TenantID string
	PageToken string
	PageSize  int
}
