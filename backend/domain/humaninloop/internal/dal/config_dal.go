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

package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// ConfigDAL 配置数据访问层
type ConfigDAL struct {
	db *gorm.DB
}

// NewConfigDAL 创建配置DAL
func NewConfigDAL(db *gorm.DB) *ConfigDAL {
	return &ConfigDAL{db: db}
}

// GetByTenantID 获取租户配置
func (d *ConfigDAL) GetByTenantID(ctx context.Context, tenantID string) (*entity.CollaborationConfig, error) {
	var config entity.CollaborationConfig
	err := d.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认配置
			return d.getDefaultConfig(tenantID), nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	return &config, nil
}

// Create 创建配置
func (d *ConfigDAL) Create(ctx context.Context, config *entity.CollaborationConfig) error {
	now := time.Now().Unix() * 1000
	config.CreatedAt = now
	config.UpdatedAt = now

	if err := d.db.WithContext(ctx).Create(config).Error; err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}
	return nil
}

// Update 更新配置
func (d *ConfigDAL) Update(ctx context.Context, config *entity.CollaborationConfig) error {
	config.UpdatedAt = time.Now().Unix() * 1000

	err := d.db.WithContext(ctx).
		Model(&entity.CollaborationConfig{}).
		Where("tenant_id = ?", config.TenantID).
		Updates(config).Error
	if err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}
	return nil
}

// Upsert 创建或更新配置
func (d *ConfigDAL) Upsert(ctx context.Context, config *entity.CollaborationConfig) error {
	now := time.Now().Unix() * 1000
	config.UpdatedAt = now

	err := d.db.WithContext(ctx).
		Model(&entity.CollaborationConfig{}).
		Where("tenant_id = ?", config.TenantID).
		Assign(config).
		FirstOrCreate(config).Error
	if err != nil {
		return fmt.Errorf("failed to upsert config: %w", err)
	}
	return nil
}

// getDefaultConfig 获取默认配置
func (d *ConfigDAL) getDefaultConfig(tenantID string) *entity.CollaborationConfig {
	return &entity.CollaborationConfig{
		TenantID:              tenantID,
		AutoReviewThreshold:    0.70,
		AutoEscalateThreshold:  0.30,
		DefaultSLAMinutes:      60,
		UrgentSLAMinutes:       15,
		NotificationEnabled:    true,
		CreatedAt:              time.Now().Unix() * 1000,
		UpdatedAt:              time.Now().Unix() * 1000,
	}
}
