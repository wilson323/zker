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

// GraphRelationshipRepository 图关系仓储接口
//go:generate mockgen -destination ../internal/mock/repository/graph_relationship_repository.go -source graph_relationship_repository.go
type GraphRelationshipRepository interface {
	// Create 创建关系
	Create(ctx context.Context, relationship *entity.GraphRelationship) error

	// CreateBatch 批量创建关系
	CreateBatch(ctx context.Context, relationships []*entity.GraphRelationship) error

	// Update 更新关系
	Update(ctx context.Context, relationship *entity.GraphRelationship) error

	// Delete 软删除关系
	Delete(ctx context.Context, relationshipID string) error

	// GetByID 根据ID获取关系
	GetByID(ctx context.Context, relationshipID string) (*entity.GraphRelationship, error)

	// GetByEntities 获取两个实体之间的关系
	GetByEntities(ctx context.Context, sourceID, targetID string, relationType entity.RelationType) (*entity.GraphRelationship, error)

	// ListByTenantID 列出租户的所有关系
	ListByTenantID(ctx context.Context, tenantID string, relationType entity.RelationType, limit, offset int) ([]*entity.GraphRelationship, int64, error)

	// ListBySourceEntity 列出实体的所有出边关系
	ListBySourceEntity(ctx context.Context, sourceEntityID string, limit int) ([]*entity.GraphRelationship, error)

	// ListByTargetEntity 列出实体的所有入边关系
	ListByTargetEntity(ctx context.Context, targetEntityID string, limit int) ([]*entity.GraphRelationship, error)

	// ListNeighbors 获取实体的邻居（双向）
	ListNeighbors(ctx context.Context, entityID string, limit int) ([]*entity.GraphRelationship, []*entity.GraphRelationship, error)

	// SearchByType 根据关系类型搜索
	SearchByType(ctx context.Context, tenantID string, relationType entity.RelationType, limit, offset int) ([]*entity.GraphRelationship, error)

	// CountByTenantID 统计租户的关系数量
	CountByTenantID(ctx context.Context, tenantID string, relationType entity.RelationType) (int64, error)

	// DeleteByEntity 删除实体的所有关系
	DeleteByEntity(ctx context.Context, entityID string) error

	// BatchDelete 批量软删除关系
	BatchDelete(ctx context.Context, relationshipIDs []string) error

	// GetRelationshipsByTypeRange 根据权重范围查询关系
	GetRelationshipsByTypeRange(ctx context.Context, tenantID string, relationType entity.RelationType, minWeight, maxWeight float32, limit int) ([]*entity.GraphRelationship, error)
}

// graphRelationshipRepository 图关系仓储实现
type graphRelationshipRepository struct {
	db *gorm.DB
}

// NewGraphRelationshipRepository 创建图关系仓储
func NewGraphRelationshipRepository(db *gorm.DB) GraphRelationshipRepository {
	return &graphRelationshipRepository{db: db}
}

// Create 创建关系
func (r *graphRelationshipRepository) Create(ctx context.Context, rel *entity.GraphRelationship) error {
	if err := r.db.WithContext(ctx).Create(rel).Error; err != nil {
		return fmt.Errorf("failed to create graph relationship: %w", err)
	}
	return nil
}

// CreateBatch 批量创建关系
func (r *graphRelationshipRepository) CreateBatch(ctx context.Context, relationships []*entity.GraphRelationship) error {
	if len(relationships) == 0 {
		return nil
	}

	batchSize := 1000
	for i := 0; i < len(relationships); i += batchSize {
		end := i + batchSize
		if end > len(relationships) {
			end = len(relationships)
		}

		batch := relationships[i:end]
		if err := r.db.WithContext(ctx).CreateInBatches(batch, batchSize).Error; err != nil {
			return fmt.Errorf("failed to batch create graph relationships: %w", err)
		}
	}

	return nil
}

// Update 更新关系
func (r *graphRelationshipRepository) Update(ctx context.Context, rel *entity.GraphRelationship) error {
	result := r.db.WithContext(ctx).
		Model(&entity.GraphRelationship{}).
		Where("id = ? AND tenant_id = ?", rel.ID, rel.TenantID).
		Updates(rel)

	if result.Error != nil {
		return fmt.Errorf("failed to update graph relationship: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Delete 软删除关系
func (r *graphRelationshipRepository) Delete(ctx context.Context, relationshipID string) error {
	now := int64(0)
	result := r.db.WithContext(ctx).
		Model(&entity.GraphRelationship{}).
		Where("id = ?", relationshipID).
		Update("deleted_at", &now)

	if result.Error != nil {
		return fmt.Errorf("failed to delete graph relationship: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// GetByID 根据ID获取关系
func (r *graphRelationshipRepository) GetByID(ctx context.Context, relationshipID string) (*entity.GraphRelationship, error) {
	var rel entity.GraphRelationship
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", relationshipID).
		First(&rel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("graph relationship not found: %s", relationshipID)
		}
		return nil, fmt.Errorf("failed to get graph relationship: %w", err)
	}

	return &rel, nil
}

// GetByEntities 获取两个实体之间的关系
func (r *graphRelationshipRepository) GetByEntities(ctx context.Context, sourceID, targetID string, relationType entity.RelationType) (*entity.GraphRelationship, error) {
	var rel entity.GraphRelationship
	err := r.db.WithContext(ctx).
		Where("source_entity_id = ? AND target_entity_id = ? AND relation_type = ? AND deleted_at IS NULL",
			sourceID, targetID, relationType).
		First(&rel).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("graph relationship not found: source=%s, target=%s, type=%s", sourceID, targetID, relationType)
		}
		return nil, fmt.Errorf("failed to get graph relationship: %w", err)
	}

	return &rel, nil
}

// ListByTenantID 列出租户的所有关系
func (r *graphRelationshipRepository) ListByTenantID(ctx context.Context, tenantID string, relationType entity.RelationType, limit, offset int) ([]*entity.GraphRelationship, int64, error) {
	var relationships []*entity.GraphRelationship
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.GraphRelationship{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count graph relationships: %w", err)
	}

	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&relationships).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list graph relationships: %w", err)
	}

	return relationships, total, nil
}

// ListBySourceEntity 列出实体的所有出边关系
func (r *graphRelationshipRepository) ListBySourceEntity(ctx context.Context, sourceEntityID string, limit int) ([]*entity.GraphRelationship, error) {
	var relationships []*entity.GraphRelationship

	err := r.db.WithContext(ctx).
		Where("source_entity_id = ? AND deleted_at IS NULL", sourceEntityID).
		Order("weight DESC").
		Limit(limit).
		Find(&relationships).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list source relationships: %w", err)
	}

	return relationships, nil
}

// ListByTargetEntity 列出实体的所有入边关系
func (r *graphRelationshipRepository) ListByTargetEntity(ctx context.Context, targetEntityID string, limit int) ([]*entity.GraphRelationship, error) {
	var relationships []*entity.GraphRelationship

	err := r.db.WithContext(ctx).
		Where("target_entity_id = ? AND deleted_at IS NULL", targetEntityID).
		Order("weight DESC").
		Limit(limit).
		Find(&relationships).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list target relationships: %w", err)
	}

	return relationships, nil
}

// ListNeighbors 获取实体的邻居（双向）
func (r *graphRelationshipRepository) ListNeighbors(ctx context.Context, entityID string, limit int) ([]*entity.GraphRelationship, []*entity.GraphRelationship, error) {
	outbound, err := r.ListBySourceEntity(ctx, entityID, limit)
	if err != nil {
		return nil, nil, err
	}

	inbound, err := r.ListByTargetEntity(ctx, entityID, limit)
	if err != nil {
		return nil, nil, err
	}

	return outbound, inbound, nil
}

// SearchByType 根据关系类型搜索
func (r *graphRelationshipRepository) SearchByType(ctx context.Context, tenantID string, relationType entity.RelationType, limit, offset int) ([]*entity.GraphRelationship, error) {
	var relationships []*entity.GraphRelationship

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND relation_type = ? AND deleted_at IS NULL", tenantID, relationType).
		Order("weight DESC").
		Limit(limit).
		Offset(offset).
		Find(&relationships).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search relationships by type: %w", err)
	}

	return relationships, nil
}

// CountByTenantID 统计租户的关系数量
func (r *graphRelationshipRepository) CountByTenantID(ctx context.Context, tenantID string, relationType entity.RelationType) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&entity.GraphRelationship{}).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count graph relationships: %w", err)
	}

	return count, nil
}

// DeleteByEntity 删除实体的所有关系
func (r *graphRelationshipRepository) DeleteByEntity(ctx context.Context, entityID string) error {
	now := int64(0)
	result := r.db.WithContext(ctx).
		Model(&entity.GraphRelationship{}).
		Where("source_entity_id = ? OR target_entity_id = ?", entityID, entityID).
		Update("deleted_at", &now)

	if result.Error != nil {
		return fmt.Errorf("failed to delete relationships by entity: %w", result.Error)
	}

	return nil
}

// BatchDelete 批量软删除关系
func (r *graphRelationshipRepository) BatchDelete(ctx context.Context, relationshipIDs []string) error {
	if len(relationshipIDs) == 0 {
		return nil
	}

	now := int64(0)
	result := r.db.WithContext(ctx).
		Model(&entity.GraphRelationship{}).
		Where("id IN ?", relationshipIDs).
		Update("deleted_at", &now)

	if result.Error != nil {
		return fmt.Errorf("failed to batch delete relationships: %w", result.Error)
	}

	return nil
}

// GetRelationshipsByTypeRange 根据权重范围查询关系
func (r *graphRelationshipRepository) GetRelationshipsByTypeRange(ctx context.Context, tenantID string, relationType entity.RelationType, minWeight, maxWeight float32, limit int) ([]*entity.GraphRelationship, error) {
	var relationships []*entity.GraphRelationship

	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND relation_type = ? AND weight >= ? AND weight <= ? AND deleted_at IS NULL",
			tenantID, relationType, minWeight, maxWeight).
		Order("weight DESC").
		Limit(limit).
		Find(&relationships).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get relationships by weight range: %w", err)
	}

	return relationships, nil
}
