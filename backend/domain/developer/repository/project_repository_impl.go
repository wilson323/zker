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

// projectRepository 项目仓储实现
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓储实例
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// Create 创建项目
func (r *projectRepository) Create(ctx context.Context, project *entity.Project) error {
	project.CreatedAt = time.Now().UnixMilli()
	project.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(project).Error
}

// GetByID 根据ID获取项目
func (r *projectRepository) GetByID(ctx context.Context, projectID string) (*entity.Project, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).
		Preload("Developer").
		Preload("APIKeys").
		Preload("Webhooks").
		Where("project_id = ?", projectID).
		Where("deleted_at IS NULL").
		First(&project).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByTenantIDAndName 根据租户ID和名称获取项目
func (r *projectRepository) GetByTenantIDAndName(ctx context.Context, tenantID, name string) (*entity.Project, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("project_name = ?", name).
		Where("deleted_at IS NULL").
		First(&project).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByDeveloperID 根据开发者ID获取项目列表
func (r *projectRepository) GetByDeveloperID(ctx context.Context, developerID string) ([]*entity.Project, error) {
	var projects []*entity.Project
	err := r.db.WithContext(ctx).
		Where("developer_id = ?", developerID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&projects).Error

	return projects, err
}

// GetByTenantID 根据租户ID获取项目列表
func (r *projectRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Project, error) {
	var projects []*entity.Project
	err := r.db.WithContext(ctx).
		Preload("Developer").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&projects).Error

	return projects, err
}

// Update 更新项目
func (r *projectRepository) Update(ctx context.Context, project *entity.Project) error {
	project.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("project_id = ?", project.ProjectID).
		Updates(project).Error
}

// Delete 软删除项目
func (r *projectRepository) Delete(ctx context.Context, projectID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("project_id = ?", projectID).
		Update("deleted_at", now).Error
}

// List 分页查询项目列表
func (r *projectRepository) List(ctx context.Context, filter *ProjectFilter) ([]*entity.Project, int64, error) {
	var projects []*entity.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Project{}).Where("deleted_at IS NULL")

	// 应用过滤器
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.DeveloperID != "" {
		query = query.Where("developer_id = ?", filter.DeveloperID)
	}
	if filter.ProjectType != "" {
		query = query.Where("project_type = ?", filter.ProjectType)
	}
	if filter.ProjectStatus != "" {
		query = query.Where("status = ?", filter.ProjectStatus)
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
		Preload("Developer").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&projects).Error

	return projects, total, err
}

// UpdateStatus 更新项目状态
func (r *projectRepository) UpdateStatus(ctx context.Context, projectID string, status entity.ProjectStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("project_id = ?", projectID).
		Update("status", status).Error
}

// UpdateConfig 更新项目配置
func (r *projectRepository) UpdateConfig(ctx context.Context, projectID string, config string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("project_id = ?", projectID).
		Update("config", config).Error
}

// CountByDeveloperID 统计开发者的项目数量
func (r *projectRepository) CountByDeveloperID(ctx context.Context, developerID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("developer_id = ?", developerID).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

// CountByTenantID 统计租户的项目数量
func (r *projectRepository) CountByTenantID(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}
