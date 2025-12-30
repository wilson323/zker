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

// GraphEntityRepository 图实体仓储接口
//go:generate mockgen -destination ../internal/mock/repository/graph_entity_repository.go -source graph_entity_repository.go
type GraphEntityRepository interface {
	// Create 创建实体
	Create(ctx context.Context, entity *entity.GraphEntity) error

	// CreateBatch 批量创建实体
	CreateBatch(ctx context.Context, entities []*entity.GraphEntity) error

	// Update 更新实体
	Update(ctx context.Context, entity *entity.GraphEntity) error

	// Delete 软删除实体
	Delete(ctx context.Context, entityID string) error

	// GetByID 根据ID获取实体
	GetByID(ctx context.Context, entityID string) (*entity.GraphEntity, error)

	// MGetByID 批量获取实体
	MGetByID(ctx context.Context, entityIDs []string) ([]*entity.GraphEntity, error)

	// GetByTenantIDAndName 根据租户ID和实体名称获取
	GetByTenantIDAndName(ctx context.Context, tenantID, entityName string) (*entity.GraphEntity, error)

	// ListByTenantID 列出租户的所有实体
	ListByTenantID(ctx context.Context, tenantID string, entityType entity.EntityType, limit, offset int) ([]*entity.GraphEntity, int64, error)

	// SearchByName 根据名称模糊搜索
	SearchByName(ctx context.Context, tenantID string, keyword string, limit int) ([]*entity.GraphEntity, error)

	// SearchByType 根据类型搜索
	SearchByType(ctx context.Context, tenantID string, entityType entity.EntityType, limit, offset int) ([]*entity.GraphEntity, error)

	// SearchByProperty 根据属性搜索
	SearchByProperty(ctx context.Context, tenantID string, propertyKey, propertyValue string, limit int) ([]*entity.GraphEntity, error)

	// Upsert 创建或更新实体（根据ID）
	Upsert(ctx context.Context, entity *entity.GraphEntity) error

	// CountByTenantID 统计租户的实体数量
	CountByTenantID(ctx context.Context, tenantID string, entityType entity.EntityType) (int64, error)

	// GetEntitiesWithEmbedding 获取有向量嵌入的实体
	GetEntitiesWithEmbedding(ctx context.Context, tenantID string, limit int) ([]*entity.GraphEntity, error)
}

// graphEntityRepository 图实体仓储实现
type graphEntityRepository struct {
	db *gorm.DB
}

// NewGraphEntityRepository 创建图实体仓储
func NewGraphEntityRepository(db *gorm.DB) GraphEntityRepository {
	return &graphEntityRepository{db: db}
}

// Create 创建实体
func (r *graphEntityRepository) Create(ctx context.Context, ent *entity.GraphEntity) error {
	if err := r.db.WithContext(ctx).Create(ent).Error; err != nil {
		return fmt.Errorf("failed to create graph entity: %w", err)
	}
	return nil
}

// CreateBatch 批量创建实体
func (r *graphEntityRepository) CreateBatch(ctx context.Context, entities []*entity.GraphEntity) error {
	if len(entities) == 0 {
		return nil
	}

	// 批量插入，每批最多1000条
	batchSize := 1000
	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}

		batch := entities[i:end]
		if err := r.db.WithContext(ctx).CreateInBatches(batch, batchSize).Error; err != nil {
			return fmt.Errorf("failed to batch create graph entities: %w", err)
		}
	}

	return nil
}

// Update 更新实体
func (r *graphEntityRepository) Update(ctx context.Context, ent *entity.GraphEntity) error {
	result := r.db.WithContext(ctx).
		Model(&entity.GraphEntity{}).
		Where("id = ? AND tenant_id = ?", ent.ID, ent.TenantID).
		Updates(ent)

	if result.Error != nil {
		return fmt.Errorf("failed to update graph entity: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete 软删除实体
func (r *graphEntityRepository) Delete(ctx context.Context, entityID string) error {
	now := int64(0)
	result := r.db.WithContext(ctx).
		Model(&entity.GraphEntity{}).
		Where("id = ?", entityID).
		Update("deleted_at", &now)

	if result.Error != nil {
		return fmt.Errorf("failed to delete graph entity: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetByID 根据ID获取实体
func (r *graphEntityRepository) GetByID(ctx context.Context, entityID string) (*entity.GraphEntity, error) {
	var ent entity.GraphEntity
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", entityID).
		First(&ent).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("graph entity not found: %s", entityID)
		}
		return nil, fmt.Errorf("failed to get graph entity: %w", err)
	}

	return &ent, nil
}

// MGetByID 批量获取实体
func (r *graphEntityRepository) MGetByID(ctx context.Context, entityIDs []string) ([]*entity.GraphEntity, error) {
	if len(entityIDs) == 0 {
		return []*entity.GraphEntity{}, nil
	}

	var entities []*entity.GraphEntity
	err := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", entityIDs).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to batch get graph entities: %w", err)
	}

	return entities, nil
}

// GetByTenantIDAndName 根据租户ID和实体名称获取
func (r *graphEntityRepository) GetByTenantIDAndName(ctx context.Context, tenantID, entityName string) (*entity.GraphEntity, error) {
	var ent entity.GraphEntity
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND entity_name = ? AND deleted_at IS NULL", tenantID, entityName).
		First(&ent).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("graph entity not found: tenant_id=%s, name=%s", tenantID, entityName)
		}
		return nil, fmt.Errorf("failed to get graph entity: %w", err)
	}

	return &ent, nil
}

// ListByTenantID 列出租户的所有实体
func (r *graphEntityRepository) ListByTenantID(ctx context.Context, tenantID string, entityType entity.EntityType, limit, offset int) ([]*entity.GraphEntity, int64, error) {
	var entities []*entity.GraphEntity
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.GraphEntity{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	// 先统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count graph entities: %w", err)
	}

	// 分页查询
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list graph entities: %w", err)
	}

	return entities, total, nil
}

// SearchByName 根据名称模糊搜索
func (r *graphEntityRepository) SearchByName(ctx context.Context, tenantID string, keyword string, limit int) ([]*entity.GraphEntity, error) {
	var entities []*entity.GraphEntity

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND entity_name LIKE ? AND deleted_at IS NULL", tenantID, "%"+keyword+"%").
		Order("created_at DESC").
		Limit(limit).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search graph entities by name: %w", err)
	}

	return entities, nil
}

// SearchByType 根据类型搜索
func (r *graphEntityRepository) SearchByType(ctx context.Context, tenantID string, entityType entity.EntityType, limit, offset int) ([]*entity.GraphEntity, error) {
	var entities []*entity.GraphEntity

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND entity_type = ? AND deleted_at IS NULL", tenantID, entityType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search graph entities by type: %w", err)
	}

	return entities, nil
}

// SearchByProperty 根据属性搜索（JSON查询）
func (r *graphEntityRepository) SearchByProperty(ctx context.Context, tenantID string, propertyKey, propertyValue string, limit int) ([]*entity.GraphEntity, error) {
	var entities []*entity.GraphEntity

	// MySQL JSON_EXTRACT函数
	jsonPath := fmt.Sprintf("$.%s", propertyKey)
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND JSON_EXTRACT(properties, ?) = ? AND deleted_at IS NULL", tenantID, jsonPath, propertyValue).
		Order("created_at DESC").
		Limit(limit).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search graph entities by property: %w", err)
	}

	return entities, nil
}

// Upsert 创建或更新实体
func (r *graphEntityRepository) Upsert(ctx context.Context, ent *entity.GraphEntity) error {
	// 使用MySQL的ON DUPLICATE KEY UPDATE
	result := r.db.WithContext(ctx).
		Model(&entity.GraphEntity{}).
		Where("id = ?", ent.ID).
		Assign(ent).
		FirstOrCreate(ent)

	if result.Error != nil {
		return fmt.Errorf("failed to upsert graph entity: %w", result.Error)
	}

	return nil
}

// CountByTenantID 统计租户的实体数量
func (r *graphEntityRepository) CountByTenantID(ctx context.Context, tenantID string, entityType entity.EntityType) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&entity.GraphEntity{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count graph entities: %w", err)
	}

	return count, nil
}

// GetEntitiesWithEmbedding 获取有向量嵌入的实体
func (r *graphEntityRepository) GetEntitiesWithEmbedding(ctx context.Context, tenantID string, limit int) ([]*entity.GraphEntity, error) {
	var entities []*entity.GraphEntity

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND embedding IS NOT NULL AND JSON_LENGTH(embedding) > 0 AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&entities).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get entities with embedding: %w", err)
	}

	return entities, nil
}
