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

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
)

// apiKeyRepository API密钥仓储实现
type apiKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository 创建API密钥仓储实例
func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

// Create 创建API密钥
func (r *apiKeyRepository) Create(ctx context.Context, apiKey *entity.APIKey) error {
	apiKey.CreatedAt = time.Now().UnixMilli()
	apiKey.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(apiKey).Error
}

// GetByID 根据ID获取API密钥
func (r *apiKeyRepository) GetByID(ctx context.Context, keyID string) (*entity.APIKey, error) {
	var apiKey entity.APIKey
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("key_id = ?", keyID).
		Where("deleted_at IS NULL").
		First(&apiKey).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// GetByProjectID 根据项目ID获取API密钥列表
func (r *apiKeyRepository) GetByProjectID(ctx context.Context, projectID string) ([]*entity.APIKey, error) {
	var apiKeys []*entity.APIKey
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&apiKeys).Error

	return apiKeys, err
}

// GetByKeySecret 根据密钥获取API密钥（用于验证）
func (r *apiKeyRepository) GetByKeySecret(ctx context.Context, keySecret string) (*entity.APIKey, error) {
	var apiKey entity.APIKey
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("key_secret = ?", keySecret).
		Where("deleted_at IS NULL").
		First(&apiKey).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// Update 更新API密钥
func (r *apiKeyRepository) Update(ctx context.Context, apiKey *entity.APIKey) error {
	apiKey.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.APIKey{}).
		Where("key_id = ?", apiKey.KeyID).
		Updates(apiKey).Error
}

// Delete 软删除API密钥
func (r *apiKeyRepository) Delete(ctx context.Context, keyID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.APIKey{}).
		Where("key_id = ?", keyID).
		Update("deleted_at", now).Error
}

// List 分页查询API密钥列表
func (r *apiKeyRepository) List(ctx context.Context, filter *APIKeyFilter) ([]*entity.APIKey, int64, error) {
	var apiKeys []*entity.APIKey
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.APIKey{}).Where("deleted_at IS NULL")

	// 应用过滤器
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.ProjectID != "" {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.KeyPrefix != "" {
		query = query.Where("key_prefix = ?", filter.KeyPrefix)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := 0
	limit := filter.PageSize
	if limit <= 0 {
		limit = 20
	}

	err := query.
		Preload("Project").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&apiKeys).Error

	return apiKeys, total, err
}

// Revoke 撤销API密钥
func (r *apiKeyRepository) Revoke(ctx context.Context, keyID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.APIKey{}).
		Where("key_id = ?", keyID).
		Update("status", entity.APIKeyStatusRevoked).Error
}

// UpdateLastUsedAt 更新最后使用时间
func (r *apiKeyRepository) UpdateLastUsedAt(ctx context.Context, keyID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.APIKey{}).
		Where("key_id = ?", keyID).
		Update("last_used_at", now).Error
}

// Expire 过期API密钥
func (r *apiKeyRepository) Expire(ctx context.Context, keyID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.APIKey{}).
		Where("key_id = ?", keyID).
		Update("status", entity.APIKeyStatusExpired).Error
}

// CountByProjectID 统计项目的API密钥数量
func (r *apiKeyRepository) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.APIKey{}).
		Where("project_id = ?", projectID).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

// GetExpiringKeys 获取即将过期的密钥（7天内）
func (r *apiKeyRepository) GetExpiringKeys(ctx context.Context, tenantID string) ([]*entity.APIKey, error) {
	var apiKeys []*entity.APIKey
	sevenDaysLater := time.Now().Add(7 * 24 * time.Hour).UnixMilli()

	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("expires_at IS NOT NULL").
		Where("expires_at <= ?", sevenDaysLater).
		Where("expires_at > ?", time.Now().UnixMilli()).
		Where("status = ?", entity.APIKeyStatusActive).
		Where("deleted_at IS NULL").
		Find(&apiKeys).Error

	return apiKeys, err
}
