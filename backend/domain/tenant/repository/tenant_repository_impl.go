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

	"github.com/coze-studio/backend/domain/tenant/entity"
)

// tenantRepository 租户仓储实现
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository 创建租户仓储实例
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

// Create 创建租户
func (r *tenantRepository) Create(ctx context.Context, tenant *entity.Tenant) error {
	tenant.CreatedAt = time.Now().UnixMilli()
	tenant.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(tenant).Error
}

// GetByID 根据ID获取租户
func (r *tenantRepository) GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.db.WithContext(ctx).
		Preload("Subscription").
		Preload("Quotas").
		Where("tenant_id = ?", tenantID).
		First(&tenant).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetByName 根据名称获取租户
func (r *tenantRepository) GetByName(ctx context.Context, name string) (*entity.Tenant, error) {
	var tenant entity.Tenant
	err := r.db.WithContext(ctx).
		Where("tenant_name = ?", name).
		Where("deleted_at IS NULL").
		First(&tenant).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Update 更新租户
func (r *tenantRepository) Update(ctx context.Context, tenant *entity.Tenant) error {
	tenant.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Tenant{}).
		Where("tenant_id = ?", tenant.TenantID).
		Updates(tenant).Error
}

// Delete 软删除租户
func (r *tenantRepository) Delete(ctx context.Context, tenantID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Tenant{}).
		Where("tenant_id = ?", tenantID).
		Updates(map[string]interface{}{
			"status":     entity.TenantStatusDeleted,
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// List 分页查询租户列表
func (r *tenantRepository) List(ctx context.Context, filter *TenantFilter) ([]*entity.Tenant, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Tenant{})

	// 应用过滤条件
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.TenantType != "" {
		query = query.Where("tenant_type = ?", filter.TenantType)
	}
	if filter.SubscriptionTier != "" {
		query = query.Where("subscription_tier = ?", filter.SubscriptionTier)
	}

	// 软删除过滤
	query = query.Where("deleted_at IS NULL")

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
	var tenants []*entity.Tenant
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Find(&tenants).Error

	return tenants, total, err
}

// UpdateStatus 更新租户状态
func (r *tenantRepository) UpdateStatus(ctx context.Context, tenantID string, status entity.TenantStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Tenant{}).
		Where("tenant_id = ?", tenantID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().UnixMilli(),
		}).Error
}

// UpdateSubscriptionTier 更新订阅等级
func (r *tenantRepository) UpdateSubscriptionTier(ctx context.Context, tenantID string, tier entity.SubscriptionTier) error {
	return r.db.WithContext(ctx).
		Model(&entity.Tenant{}).
		Where("tenant_id = ?", tenantID).
		Updates(map[string]interface{}{
			"subscription_tier": tier,
			"updated_at":        time.Now().UnixMilli(),
		}).Error
}

// subscriptionRepository 订阅仓储实现
type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository 创建订阅仓储实例
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// Create 创建订阅
func (r *subscriptionRepository) Create(ctx context.Context, sub *entity.Subscription) error {
	sub.CreatedAt = time.Now().UnixMilli()
	sub.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(sub).Error
}

// GetByID 根据ID获取订阅
func (r *subscriptionRepository) GetByID(ctx context.Context, subscriptionID string) (*entity.Subscription, error) {
	var sub entity.Subscription
	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Where("subscription_id = ?", subscriptionID).
		First(&sub).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetByTenantID 根据租户ID获取订阅
func (r *subscriptionRepository) GetByTenantID(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	var sub entity.Subscription
	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Where("tenant_id = ?", tenantID).
		First(&sub).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetByTenant 根据租户ID获取订阅（别名方法）
func (r *subscriptionRepository) GetByTenant(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	return r.GetByTenantID(ctx, tenantID)
}

// GetAll 获取所有订阅
func (r *subscriptionRepository) GetAll(ctx context.Context) ([]*entity.Subscription, error) {
	var subs []*entity.Subscription
	err := r.db.WithContext(ctx).
		Preload("Tenant").
		Find(&subs).Error
	return subs, err
}

// Update 更新订阅
func (r *subscriptionRepository) Update(ctx context.Context, sub *entity.Subscription) error {
	sub.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("subscription_id = ?", sub.SubscriptionID).
		Updates(sub).Error
}

// List 分页查询订阅列表
func (r *subscriptionRepository) List(ctx context.Context, filter *SubscriptionFilter) ([]*entity.Subscription, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Subscription{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.PlanTier != "" {
		query = query.Where("plan_tier = ?", filter.PlanTier)
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
	var subs []*entity.Subscription
	err := query.
		Preload("Tenant").
		Order("created_at DESC").
		Limit(pageSize).
		Find(&subs).Error

	return subs, total, err
}

// UpdateStatus 更新订阅状态
func (r *subscriptionRepository) UpdateStatus(ctx context.Context, subscriptionID string, status entity.SubscriptionStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Subscription{}).
		Where("subscription_id = ?", subscriptionID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now().UnixMilli(),
		}).Error
}

// quotaRepository 配额仓储实现
type quotaRepository struct {
	db *gorm.DB
}

// NewQuotaRepository 创建配额仓储实例
func NewQuotaRepository(db *gorm.DB) QuotaRepository {
	return &quotaRepository{db: db}
}

// Create 创建配额
func (r *quotaRepository) Create(ctx context.Context, quota *entity.Quota) error {
	quota.CreatedAt = time.Now().UnixMilli()
	quota.UpdatedAt = time.Now().UnixMilli()
	quota.LastResetAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(quota).Error
}

// GetByID 根据ID获取配额
func (r *quotaRepository) GetByID(ctx context.Context, quotaID string) (*entity.Quota, error) {
	var quota entity.Quota
	err := r.db.WithContext(ctx).
		Where("quota_id = ?", quotaID).
		First(&quota).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// GetByTenantAndResource 根据租户ID和资源类型获取配额
func (r *quotaRepository) GetByTenantAndResource(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error) {
	var quota entity.Quota
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("resource_type = ?", resourceType).
		First(&quota).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

// Update 更新配额
func (r *quotaRepository) Update(ctx context.Context, quota *entity.Quota) error {
	quota.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Quota{}).
		Where("quota_id = ?", quota.QuotaID).
		Updates(quota).Error
}

// UpdateUsedCount 更新使用计数
func (r *quotaRepository) UpdateUsedCount(ctx context.Context, quotaID string, delta int) error {
	return r.db.WithContext(ctx).
		Model(&entity.Quota{}).
		Where("quota_id = ?", quotaID).
		Updates(map[string]interface{}{
			"used_count":  gorm.Expr("used_count + ?", delta),
			"updated_at":  time.Now().UnixMilli(),
		}).Error
}

// ResetUsage 重置使用量
func (r *quotaRepository) ResetUsage(ctx context.Context, quotaID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Quota{}).
		Where("quota_id = ?", quotaID).
		Updates(map[string]interface{}{
			"used_count":   0,
			"last_reset_at": time.Now().UnixMilli(),
			"updated_at":    time.Now().UnixMilli(),
		}).Error
}

// List 获取租户的所有配额
func (r *quotaRepository) List(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	var quotas []*entity.Quota
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&quotas).Error
	return quotas, err
}

// GetByTenant 获取租户的所有配额（别名方法）
func (r *quotaRepository) GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	return r.List(ctx, tenantID)
}

// GetAll 获取所有配额（用于监控）
func (r *quotaRepository) GetAll(ctx context.Context) ([]*entity.Quota, error) {
	var quotas []*entity.Quota
	err := r.db.WithContext(ctx).
		Find(&quotas).Error
	return quotas, err
}
