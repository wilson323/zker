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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/entity"
)

// GraphQueryRepository 图查询仓储接口
//go:generate mockgen -destination ../internal/mock/repository/graph_query_repository.go -source graph_query_repository.go
type GraphQueryRepository interface {
	// Create 创建查询
	Create(ctx context.Context, query *entity.GraphQuery) error

	// Update 更新查询
	Update(ctx context.Context, query *entity.GraphQuery) error

	// Delete 软删除查询
	Delete(ctx context.Context, queryID string) error

	// GetByID 根据ID获取查询
	GetByID(ctx context.Context, queryID string) (*entity.GraphQuery, error)

	// ListByTenantID 列出租户的所有查询
	ListByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*entity.GraphQuery, int64, error)

	// ListByCreator 列出创建者的查询
	ListByCreator(ctx context.Context, tenantID, creatorID string, limit, offset int) ([]*entity.GraphQuery, error)

	// ListPublic 列出公开查询
	ListPublic(ctx context.Context, limit, offset int) ([]*entity.GraphQuery, int64, error)

	// SearchByName 根据名称搜索
	SearchByName(ctx context.Context, tenantID string, keyword string, limit int) ([]*entity.GraphQuery, error)

	// CountByTenantID 统计租户的查询数量
	CountByTenantID(ctx context.Context, tenantID string) (int64, error)

	// GetByName 根据名称获取查询
	GetByName(ctx context.Context, tenantID, queryName string) (*entity.GraphQuery, error)
}

// graphQueryRepository 图查询仓储实现
type graphQueryRepository struct {
	db *gorm.DB
}

// NewGraphQueryRepository 创建图查询仓储
func NewGraphQueryRepository(db *gorm.DB) GraphQueryRepository {
	return &graphQueryRepository{db: db}
}

// Create 创建查询
func (r *graphQueryRepository) Create(ctx context.Context, query *entity.GraphQuery) error {
	if err := r.db.WithContext(ctx).Create(query).Error; err != nil {
		return fmt.Errorf("failed to create graph query: %w", err)
	}
	return nil
}

// Update 更新查询
func (r *graphQueryRepository) Update(ctx context.Context, query *entity.GraphQuery) error {
	result := r.db.WithContext(ctx).
		Model(&entity.GraphQuery{}).
		Where("id = ? AND tenant_id = ?", query.ID, query.TenantID).
		Updates(query)

	if result.Error != nil {
		return fmt.Errorf("failed to update graph query: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete 软删除查询
func (r *graphQueryRepository) Delete(ctx context.Context, queryID string) error {
	now := int64(0)
	result := r.db.WithContext(ctx).
		Model(&entity.GraphQuery{}).
		Where("id = ?", queryID).
		Update("deleted_at", &now)

	if result.Error != nil {
		return fmt.Errorf("failed to delete graph query: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetByID 根据ID获取查询
func (r *graphQueryRepository) GetByID(ctx context.Context, queryID string) (*entity.GraphQuery, error) {
	var query entity.GraphQuery
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", queryID).
		First(&query).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("graph query not found: %s", queryID)
		}
		return nil, fmt.Errorf("failed to get graph query: %w", err)
	}

	return &query, nil
}

// ListByTenantID 列出租户的所有查询
func (r *graphQueryRepository) ListByTenantID(ctx context.Context, tenantID string, limit, offset int) ([]*entity.GraphQuery, int64, error) {
	var queries []*entity.GraphQuery
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.GraphQuery{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count graph queries: %w", err)
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&queries).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list graph queries: %w", err)
	}

	return queries, total, nil
}

// ListByCreator 列出创建者的查询
func (r *graphQueryRepository) ListByCreator(ctx context.Context, tenantID, creatorID string, limit, offset int) ([]*entity.GraphQuery, error) {
	var queries []*entity.GraphQuery

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND created_by = ? AND deleted_at IS NULL", tenantID, creatorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&queries).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list queries by creator: %w", err)
	}

	return queries, nil
}

// ListPublic 列出公开查询
func (r *graphQueryRepository) ListPublic(ctx context.Context, limit, offset int) ([]*entity.GraphQuery, int64, error) {
	var queries []*entity.GraphQuery
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.GraphQuery{}).
		Where("is_public = ? AND deleted_at IS NULL", true)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count public queries: %w", err)
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&queries).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list public queries: %w", err)
	}

	return queries, total, nil
}

// SearchByName 根据名称搜索
func (r *graphQueryRepository) SearchByName(ctx context.Context, tenantID string, keyword string, limit int) ([]*entity.GraphQuery, error) {
	var queries []*entity.GraphQuery

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND query_name LIKE ? AND deleted_at IS NULL", tenantID, "%"+keyword+"%").
		Order("created_at DESC").
		Limit(limit).
		Find(&queries).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search queries by name: %w", err)
	}

	return queries, nil
}

// CountByTenantID 统计租户的查询数量
func (r *graphQueryRepository) CountByTenantID(ctx context.Context, tenantID string) (int64, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&entity.GraphQuery{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count graph queries: %w", err)
	}

	return count, nil
}

// GetByName 根据名称获取查询
func (r *graphQueryRepository) GetByName(ctx context.Context, tenantID, queryName string) (*entity.GraphQuery, error) {
	var query entity.GraphQuery
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND query_name = ? AND deleted_at IS NULL", tenantID, queryName).
		First(&query).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("graph query not found: tenant_id=%s, name=%s", tenantID, queryName)
		}
		return nil, fmt.Errorf("failed to get graph query: %w", err)
	}

	return &query, nil
}
