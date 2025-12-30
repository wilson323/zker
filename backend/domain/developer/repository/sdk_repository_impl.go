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

// sdkRepository SDK仓储实现
type sdkRepository struct {
	db *gorm.DB
}

// NewSDKRepository 创建SDK仓储实例
func NewSDKRepository(db *gorm.DB) SDKRepository {
	return &sdkRepository{db: db}
}

// Create 创建SDK
func (r *sdkRepository) Create(ctx context.Context, sdk *entity.SDK) error {
	sdk.CreatedAt = time.Now().UnixMilli()
	sdk.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(sdk).Error
}

// GetByID 根据ID获取SDK
func (r *sdkRepository) GetByID(ctx context.Context, sdkID string) (*entity.SDK, error) {
	var sdk entity.SDK
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("sdk_id = ?", sdkID).
		Where("deleted_at IS NULL").
		First(&sdk).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sdk, nil
}

// GetByProjectID 根据项目ID获取SDK列表
func (r *sdkRepository) GetByProjectID(ctx context.Context, projectID string) ([]*entity.SDK, error) {
	var sdks []*entity.SDK
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&sdks).Error

	return sdks, err
}

// GetByProjectIDAndLanguage 根据项目ID和语言获取SDK
func (r *sdkRepository) GetByProjectIDAndLanguage(ctx context.Context, projectID string, language entity.SDKLanguage) (*entity.SDK, error) {
	var sdk entity.SDK
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("language = ?", language).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		First(&sdk).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sdk, nil
}

// Update 更新SDK
func (r *sdkRepository) Update(ctx context.Context, sdk *entity.SDK) error {
	sdk.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.SDK{}).
		Where("sdk_id = ?", sdk.SDKID).
		Updates(sdk).Error
}

// Delete 软删除SDK
func (r *sdkRepository) Delete(ctx context.Context, sdkID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.SDK{}).
		Where("sdk_id = ?", sdkID).
		Update("deleted_at", now).Error
}

// List 分页查询SDK列表
func (r *sdkRepository) List(ctx context.Context, filter *SDKFilter) ([]*entity.SDK, int64, error) {
	var sdks []*entity.SDK
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.SDK{}).Where("deleted_at IS NULL")

	// 应用过滤器
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.ProjectID != "" {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Language != "" {
		query = query.Where("language = ?", filter.Language)
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
		Find(&sdks).Error

	return sdks, total, err
}

// UpdateStatus 更新SDK状态
func (r *sdkRepository) UpdateStatus(ctx context.Context, sdkID string, status entity.SDKStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.SDK{}).
		Where("sdk_id = ?", sdkID).
		Update("status", status).Error
}

// IncrementDownloadCount 增加下载次数
func (r *sdkRepository) IncrementDownloadCount(ctx context.Context, sdkID string) error {
	return r.db.WithContext(ctx).
		Model(&entity.SDK{}).
		Where("sdk_id = ?", sdkID).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

// GetLatestVersion 获取最新版本的SDK
func (r *sdkRepository) GetLatestVersion(ctx context.Context, projectID string, language entity.SDKLanguage) (*entity.SDK, error) {
	var sdk entity.SDK
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("language = ?", language).
		Where("status = ?", entity.SDKStatusPublished).
		Where("deleted_at IS NULL").
		Order("version DESC").
		First(&sdk).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &sdk, nil
}
