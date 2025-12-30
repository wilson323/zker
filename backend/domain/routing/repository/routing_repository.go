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

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
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

// IntentRepository 意图仓储接口
type IntentRepository interface {
	// Create 创建意图
	Create(ctx context.Context, intent *entity.Intent) error

	// GetByID 根据ID获取意图
	GetByID(ctx context.Context, intentID string) (*entity.Intent, error)

	// GetByTenantAndName 根据租户和名称获取意图
	GetByTenantAndName(ctx context.Context, tenantID, intentName string) (*entity.Intent, error)

	// Update 更新意图
	Update(ctx context.Context, intent *entity.Intent) error

	// Delete 删除意图（软删除）
	Delete(ctx context.Context, intentID string) error

	// List 列出租户的所有意图
	List(ctx context.Context, filter *IntentFilter) ([]*entity.Intent, int64, error)

	// GetActiveIntentsByTenant 获取租户的所有激活意图
	GetActiveIntentsByTenant(ctx context.Context, tenantID string) ([]*entity.Intent, error)
}

// IntentFilter 意图查询过滤器
type IntentFilter struct {
	TenantID   string
	AgentID    string
	WorkflowID string
	IsActive   *bool
	PageToken  string
	PageSize   int
}

// IntentSampleRepository 意图样本仓储接口
type IntentSampleRepository interface {
	// CreateBatch 批量创建样本
	CreateBatch(ctx context.Context, samples []*entity.IntentSample) error

	// GetByIntentID 获取意图的所有样本
	GetByIntentID(ctx context.Context, intentID string) ([]*entity.IntentSample, error)

	// DeleteByIntentID 删除意图的所有样本
	DeleteByIntentID(ctx context.Context, intentID string) error

	// SearchByText 使用向量相似度搜索样本
	SearchByText(ctx context.Context, tenantID string, embedding []float64, limit int) ([]*entity.IntentSample, error)
}

// ABTestRepository A/B测试仓储接口
type ABTestRepository interface {
	// Create 创建A/B测试
	Create(ctx context.Context, test *entity.ABTest) error

	// GetByID 根据ID获取测试
	GetByID(ctx context.Context, testID string) (*entity.ABTest, error)

	// Update 更新测试
	Update(ctx context.Context, test *entity.ABTest) error

	// Delete 删除测试（软删除）
	Delete(ctx context.Context, testID string) error

	// List 列出租户的所有测试
	List(ctx context.Context, filter *ABTestFilter) ([]*entity.ABTest, int64, error)

	// GetActiveTestsByTenant 获取租户的所有运行中的测试
	GetActiveTestsByTenant(ctx context.Context, tenantID string) ([]*entity.ABTest, error)
}

// ABTestFilter A/B测试查询过滤器
type ABTestFilter struct {
	TenantID string
	Status   *entity.ABTestStatus
	PageToken string
	PageSize  int
}

// ABTestRecordRepository A/B测试记录仓储接口
type ABTestRecordRepository interface {
	// Create 创建测试记录
	Create(ctx context.Context, record *entity.ABTestRecord) error

	// CreateBatch 批量创建测试记录
	CreateBatch(ctx context.Context, records []*entity.ABTestRecord) error

	// GetByTestID 获取测试的所有记录
	GetByTestID(ctx context.Context, testID string) ([]*entity.ABTestRecord, error)

	// GetStatsByTestID 获取测试的统计数据
	GetStatsByTestID(ctx context.Context, testID string) (map[string]*entity.TestStrategyStats, error)

	// DeleteByTestID 删除测试的所有记录
	DeleteByTestID(ctx context.Context, testID string) error
}

// RoutingOptimizeLogRepository 路由优化日志仓储接口
type RoutingOptimizeLogRepository interface {
	// Create 创建优化日志
	Create(ctx context.Context, log *entity.RoutingOptimizeLog) error

	// CreateBatch 批量创建优化日志
	CreateBatch(ctx context.Context, logs []*entity.RoutingOptimizeLog) error

	// GetByTimeRange 获取时间范围内的日志
	GetByTimeRange(ctx context.Context, tenantID string, startTime, endTime int64) ([]*entity.RoutingOptimizeLog, error)

	// GetByIntent 获取特定意图的日志
	GetByIntent(ctx context.Context, tenantID, intent string, limit int) ([]*entity.RoutingOptimizeLog, error)

	// GetByAgent 获取特定Agent的日志
	GetByAgent(ctx context.Context, tenantID, agentID string, limit int) ([]*entity.RoutingOptimizeLog, error)

	// DeleteOldLogs 删除旧日志
	DeleteOldLogs(ctx context.Context, beforeDate int64) error
}
