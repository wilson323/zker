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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/config/entity"
	"github.com/coze-dev/coze-studio/backend/domain/config/repository"
)

// ConfigRepositoryImpl 配置仓储实现
type ConfigRepositoryImpl struct {
	db *gorm.DB
}

// NewConfigRepository 创建配置仓储
func NewConfigRepository(db *gorm.DB) repository.ConfigRepository {
	return &ConfigRepositoryImpl{db: db}
}

// Create 创建配置项
func (r *ConfigRepositoryImpl) Create(ctx context.Context, config *entity.ConfigItem) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// GetByID 根据ID获取配置
func (r *ConfigRepositoryImpl) GetByID(ctx context.Context, configID string) (*entity.ConfigItem, error) {
	var config entity.ConfigItem
	err := r.db.WithContext(ctx).
		Where("config_id = ? AND deleted_at IS NULL", configID).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetByKey 根据租户和key获取配置
func (r *ConfigRepositoryImpl) GetByKey(ctx context.Context, tenantID, key string) (*entity.ConfigItem, error) {
	var config entity.ConfigItem
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND config_key = ? AND deleted_at IS NULL", tenantID, key).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ListByTenant 获取租户所有配置
func (r *ConfigRepositoryImpl) ListByTenant(ctx context.Context, tenantID string) ([]*entity.ConfigItem, error) {
	var configs []*entity.ConfigItem
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// Update 更新配置
func (r *ConfigRepositoryImpl) Update(ctx context.Context, config *entity.ConfigItem) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// Delete 删除配置（软删除）
func (r *ConfigRepositoryImpl) Delete(ctx context.Context, configID string) error {
	return r.db.WithContext(ctx).
		Where("config_id = ?", configID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// GetByPrefix 根据前缀获取配置列表
func (r *ConfigRepositoryImpl) GetByPrefix(ctx context.Context, tenantID, prefix string) ([]*entity.ConfigItem, error) {
	var configs []*entity.ConfigItem
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND config_key LIKE ? AND deleted_at IS NULL", tenantID, prefix+"%").
		Order("config_key ASC").
		Find(&configs).Error
	if err != nil {
		return nil, err
	}
	return configs, nil
}

// ConfigHistoryRepositoryImpl 配置历史仓储实现
type ConfigHistoryRepositoryImpl struct {
	db *gorm.DB
}

// NewConfigHistoryRepository 创建配置历史仓储
func NewConfigHistoryRepository(db *gorm.DB) repository.ConfigHistoryRepository {
	return &ConfigHistoryRepositoryImpl{db: db}
}

// Create 创建历史记录
func (r *ConfigHistoryRepositoryImpl) Create(ctx context.Context, history *entity.ConfigHistory) error {
	return r.db.WithContext(ctx).Create(history).Error
}

// ListByConfigID 获取配置的所有历史记录
func (r *ConfigHistoryRepositoryImpl) ListByConfigID(
	ctx context.Context,
	configID string,
	limit,
	offset int,
) ([]*entity.ConfigHistory, error) {
	var histories []*entity.ConfigHistory
	err := r.db.WithContext(ctx).
		Where("config_id = ? AND deleted_at IS NULL", configID).
		Order("version_number DESC").
		Limit(limit).
		Offset(offset).
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

// GetLatestVersion 获取最新版本号
func (r *ConfigHistoryRepositoryImpl) GetLatestVersion(ctx context.Context, configID string) (int, error) {
	var version int
	err := r.db.WithContext(ctx).
		Model(&entity.ConfigHistory{}).
		Where("config_id = ? AND deleted_at IS NULL", configID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&version).Error
	return version, err
}

// ListByTenant 获取租户的所有配置变更历史
func (r *ConfigHistoryRepositoryImpl) ListByTenant(
	ctx context.Context,
	tenantID string,
	limit,
	offset int,
) ([]*entity.ConfigHistory, error) {
	var histories []*entity.ConfigHistory
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("changed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}
