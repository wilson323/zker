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
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// quotaRepositoryImpl 配额仓储实现
type quotaRepositoryImpl struct {
	db *gorm.DB
}

// NewQuotaRepository 创建配额仓储实例
func NewQuotaRepository(db *gorm.DB) QuotaRepository {
	return &quotaRepositoryImpl{db: db}
}

// Create 创建配额
func (r *quotaRepositoryImpl) Create(ctx context.Context, quota *entity.Quota) error {
	if err := r.db.WithContext(ctx).Create(quota).Error; err != nil {
		return fmt.Errorf("failed to create quota: %w", err)
	}
	return nil
}

// GetByID 根据ID获取配额
func (r *quotaRepositoryImpl) GetByID(ctx context.Context, quotaID string) (*entity.Quota, error) {
	var quota entity.Quota
	if err := r.db.WithContext(ctx).Where("quota_id = ?", quotaID).First(&quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get quota by ID: %w", err)
	}
	return &quota, nil
}

// GetByTenantAndResource 根据租户ID和资源类型获取配额
func (r *quotaRepositoryImpl) GetByTenantAndResource(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error) {
	var quota entity.Quota
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
		First(&quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get quota by tenant and resource: %w", err)
	}
	return &quota, nil
}

// GetByTenant 根据租户ID获取所有配额
func (r *quotaRepositoryImpl) GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	var quotas []*entity.Quota
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&quotas).Error; err != nil {
		return nil, fmt.Errorf("failed to get quotas by tenant: %w", err)
	}
	return quotas, nil
}

// List 获取租户的所有配额
func (r *quotaRepositoryImpl) List(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	var quotas []*entity.Quota
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&quotas).Error; err != nil {
		return nil, fmt.Errorf("failed to list quotas: %w", err)
	}
	return quotas, nil
}

// GetAll 获取所有配额
func (r *quotaRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Quota, error) {
	var quotas []*entity.Quota
	if err := r.db.WithContext(ctx).Find(&quotas).Error; err != nil {
		return nil, fmt.Errorf("failed to get all quotas: %w", err)
	}
	return quotas, nil
}

// Update 更新配额
func (r *quotaRepositoryImpl) Update(ctx context.Context, quota *entity.Quota) error {
	quota.UpdatedAt = time.Now().UnixMilli()
	if err := r.db.WithContext(ctx).Save(quota).Error; err != nil {
		return fmt.Errorf("failed to update quota: %w", err)
	}
	return nil
}

// UpdateUsedCount 更新使用计数（count可为正数或负数）
func (r *quotaRepositoryImpl) UpdateUsedCount(ctx context.Context, quotaID string, count int) error {
	now := time.Now().UnixMilli()
	result := r.db.WithContext(ctx).
		Model(&entity.Quota{}).
		Where("quota_id = ?", quotaID).
		Updates(map[string]interface{}{
			"used_count": gorm.Expr("used_count + ?", count),
			"updated_at": now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update used count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("quota not found: %s", quotaID)
	}

	return nil
}

// ResetUsage 重置使用计数
func (r *quotaRepositoryImpl) ResetUsage(ctx context.Context, quotaID string) error {
	now := time.Now().UnixMilli()
	result := r.db.WithContext(ctx).
		Model(&entity.Quota{}).
		Where("quota_id = ?", quotaID).
		Updates(map[string]interface{}{
			"used_count":   0,
			"last_reset_at": now,
			"updated_at":   now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to reset usage: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("quota not found: %s", quotaID)
	}

	return nil
}

// Delete 删除配额
func (r *quotaRepositoryImpl) Delete(ctx context.Context, quotaID string) error {
	result := r.db.WithContext(ctx).
		Where("quota_id = ?", quotaID).
		Delete(&entity.Quota{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete quota: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("quota not found: %s", quotaID)
	}

	return nil
}

// BatchCreate 批量创建配额
func (r *quotaRepositoryImpl) BatchCreate(ctx context.Context, quotas []*entity.Quota) error {
	if len(quotas) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(quotas, 100).Error; err != nil {
		return fmt.Errorf("failed to batch create quotas: %w", err)
	}

	return nil
}
