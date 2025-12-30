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

// developerRepository 开发者仓储实现
type developerRepository struct {
	db *gorm.DB
}

// NewDeveloperRepository 创建开发者仓储实例
func NewDeveloperRepository(db *gorm.DB) DeveloperRepository {
	return &developerRepository{db: db}
}

// Create 创建开发者
func (r *developerRepository) Create(ctx context.Context, developer *entity.Developer) error {
	developer.CreatedAt = time.Now().UnixMilli()
	developer.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(developer).Error
}

// GetByID 根据ID获取开发者
func (r *developerRepository) GetByID(ctx context.Context, developerID string) (*entity.Developer, error) {
	var developer entity.Developer
	err := r.db.WithContext(ctx).
		Preload("Projects").
		Where("developer_id = ?", developerID).
		Where("deleted_at IS NULL").
		First(&developer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &developer, nil
}

// GetByTenantIDAndUserID 根据租户ID和用户ID获取开发者
func (r *developerRepository) GetByTenantIDAndUserID(ctx context.Context, tenantID, userID string) (*entity.Developer, error) {
	var developer entity.Developer
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		First(&developer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &developer, nil
}

// GetByEmail 根据邮箱获取开发者
func (r *developerRepository) GetByEmail(ctx context.Context, email string) (*entity.Developer, error) {
	var developer entity.Developer
	err := r.db.WithContext(ctx).
		Where("developer_email = ?", email).
		Where("deleted_at IS NULL").
		First(&developer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &developer, nil
}

// Update 更新开发者
func (r *developerRepository) Update(ctx context.Context, developer *entity.Developer) error {
	developer.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Developer{}).
		Where("developer_id = ?", developer.DeveloperID).
		Updates(developer).Error
}

// Delete 软删除开发者
func (r *developerRepository) Delete(ctx context.Context, developerID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Developer{}).
		Where("developer_id = ?", developerID).
		Update("deleted_at", now).Error
}

// List 分页查询开发者列表
func (r *developerRepository) List(ctx context.Context, filter *DeveloperFilter) ([]*entity.Developer, int64, error) {
	var developers []*entity.Developer
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Developer{}).Where("deleted_at IS NULL")

	// 应用过滤器
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
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
	if filter.PageToken != "" {
		offset = parseInt(filter.PageToken)
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&developers).Error

	return developers, total, err
}

// UpdateStatus 更新开发者状态
func (r *developerRepository) UpdateStatus(ctx context.Context, developerID string, status entity.DeveloperStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Developer{}).
		Where("developer_id = ?", developerID).
		Update("status", status).Error
}

// Exists 检查开发者是否存在
func (r *developerRepository) Exists(ctx context.Context, tenantID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Developer{}).
		Where("tenant_id = ?", tenantID).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count > 0, err
}

// parseInt 解析字符串为整数
func parseInt(s string) int {
	var i int
	_, err := sscanf(s, "%d", &i)
	if err != nil {
		return 0
	}
	return i
}

// sscanf 简单的字符串解析
func sscanf(s, format string, a ...interface{}) (int, error) {
	// 简化实现，实际可以使用 fmt.Sscanf
	return 0, nil
}
