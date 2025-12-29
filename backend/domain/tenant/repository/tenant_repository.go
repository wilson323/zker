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

	"github.com/coze-studio/backend/domain/tenant/entity"
)

// TenantRepository 租户仓储接口
type TenantRepository interface {
	// Create 创建租户
	Create(ctx context.Context, tenant *entity.Tenant) error

	// GetByID 根据ID获取租户
	GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error)

	// GetByName 根据名称获取租户
	GetByName(ctx context.Context, name string) (*entity.Tenant, error)

	// Update 更新租户
	Update(ctx context.Context, tenant *entity.Tenant) error

	// Delete 软删除租户
	Delete(ctx context.Context, tenantID string) error

	// List 分页查询租户列表
	List(ctx context.Context, filter *TenantFilter) ([]*entity.Tenant, int64, error)

	// UpdateStatus 更新租户状态
	UpdateStatus(ctx context.Context, tenantID string, status entity.TenantStatus) error

	// UpdateSubscriptionTier 更新订阅等级
	UpdateSubscriptionTier(ctx context.Context, tenantID string, tier entity.SubscriptionTier) error
}

// TenantFilter 租户查询过滤器
type TenantFilter struct {
	Status           entity.TenantStatus
	TenantType       entity.TenantType
	SubscriptionTier entity.SubscriptionTier
	PageToken        string
	PageSize         int
}

// SubscriptionRepository 订阅仓储接口
type SubscriptionRepository interface {
	// Create 创建订阅
	Create(ctx context.Context, sub *entity.Subscription) error

	// GetByID 根据ID获取订阅
	GetByID(ctx context.Context, subscriptionID string) (*entity.Subscription, error)

	// GetByTenantID 根据租户ID获取订阅
	GetByTenantID(ctx context.Context, tenantID string) (*entity.Subscription, error)

	// GetByTenant 根据租户ID获取订阅（别名方法）
	GetByTenant(ctx context.Context, tenantID string) (*entity.Subscription, error)

	// GetAll 获取所有订阅
	GetAll(ctx context.Context) ([]*entity.Subscription, error)

	// Update 更新订阅
	Update(ctx context.Context, sub *entity.Subscription) error

	// List 分页查询订阅列表
	List(ctx context.Context, filter *SubscriptionFilter) ([]*entity.Subscription, int64, error)

	// UpdateStatus 更新订阅状态
	UpdateStatus(ctx context.Context, subscriptionID string, status entity.SubscriptionStatus) error
}

// SubscriptionFilter 订阅查询过滤器
type SubscriptionFilter struct {
	TenantID string
	Status   entity.SubscriptionStatus
	PlanTier entity.SubscriptionTier
	PageToken string
	PageSize  int
}

// QuotaRepository 配额仓储接口
type QuotaRepository interface {
	// Create 创建配额
	Create(ctx context.Context, quota *entity.Quota) error

	// GetByID 根据ID获取配额
	GetByID(ctx context.Context, quotaID string) (*entity.Quota, error)

	// GetByTenantAndResource 根据租户ID和资源类型获取配额
	GetByTenantAndResource(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error)

	// Update 更新配额
	Update(ctx context.Context, quota *entity.Quota) error

	// UpdateUsedCount 更新使用计数
	UpdateUsedCount(ctx context.Context, quotaID string, delta int) error

	// ResetUsage 重置使用量
	ResetUsage(ctx context.Context, quotaID string) error

	// List 获取租户的所有配额
	List(ctx context.Context, tenantID string) ([]*entity.Quota, error)

	// GetAll 获取所有配额（用于监控）
	GetAll(ctx context.Context) ([]*entity.Quota, error)

	// GetByTenant 获取租户的所有配额
	GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error)
}

// UsageLogRepository 使用日志仓储接口
type UsageLogRepository interface {
	// Create 创建使用日志
	Create(ctx context.Context, log interface{}) error

	// CreateBatch 批量创建使用日志
	CreateBatch(ctx context.Context, logs []interface{}) error

	// GetByTenantAndDateRange 根据租户ID和日期范围获取日志
	GetByTenantAndDateRange(ctx context.Context, tenantID string, startDate, endDate interface{}) ([]interface{}, error)

	// GetByID 根据ID获取日志
	GetByID(ctx context.Context, logID string) (interface{}, error)

	// List 分页查询日志
	List(ctx context.Context, filter *UsageLogFilter) ([]interface{}, int64, error)
}

// UsageLogFilter 使用日志查询过滤器
type UsageLogFilter struct {
	TenantID     string
	ResourceType entity.ResourceType
	Action       string
	StartDate    interface{}
	EndDate      interface{}
	PageToken    string
	PageSize     int
}

// InvoiceRepository 账单仓储接口
type InvoiceRepository interface {
	// Create 创建账单
	Create(ctx context.Context, invoice interface{}) error

	// GetByID 根据ID获取账单
	GetByID(ctx context.Context, invoiceID string) (interface{}, error)

	// GetByTenantAndDateRange 根据租户ID和日期范围获取账单
	GetByTenantAndDateRange(ctx context.Context, tenantID string, startDate, endDate interface{}) (interface{}, error)

	// GetByTenant 获取租户的账单列表
	GetByTenant(ctx context.Context, tenantID string, limit, offset int) ([]interface{}, int64, error)

	// Update 更新账单
	Update(ctx context.Context, invoice interface{}) error

	// List 分页查询账单
	List(ctx context.Context, filter *InvoiceFilter) ([]interface{}, int64, error)
}

// InvoiceFilter 账单查询过滤器
type InvoiceFilter struct {
	TenantID       string
	SubscriptionID string
	Status         string
	StartDate      interface{}
	EndDate        interface{}
	PageToken      string
	PageSize       int
}

