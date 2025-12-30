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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
	"github.com/coze-dev/coze-studio/backend/domain/memory/internal/dal/model"
)

// NewKnowledgeMemoryDAO 创建知识记忆DAO
func NewKnowledgeMemoryDAO(db *gorm.DB) *KnowledgeMemoryDAO {
	return &KnowledgeMemoryDAO{db: db}
}

// KnowledgeMemoryDAO 知识记忆数据访问对象
type KnowledgeMemoryDAO struct {
	db *gorm.DB
}

// Create 创建知识
func (dao *KnowledgeMemoryDAO) Create(ctx context.Context, knowledge *entity.KnowledgeMemory) error {
	daoModel := dao.toModel(knowledge)
	result := dao.db.WithContext(ctx).Create(daoModel)
	if result.Error != nil {
		return result.Error
	}
	knowledge.ID = daoModel.ID
	return nil
}

// FindByID 根据ID查找知识
func (dao *KnowledgeMemoryDAO) FindByID(ctx context.Context, memoryID string) (*entity.KnowledgeMemory, error) {
	var daoModel model.KnowledgeMemoryDAO
	err := dao.db.WithContext(ctx).Where("memory_id = ?", memoryID).First(&daoModel).Error
	if err != nil {
		return nil, err
	}
	return dao.toEntity(&daoModel), nil
}

// FindByTenantID 根据租户ID查找知识
func (dao *KnowledgeMemoryDAO) FindByTenantID(
	ctx context.Context,
	tenantID string,
	knowledgeType *entity.KnowledgeType,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	var daoModels []model.KnowledgeMemoryDAO
	query := dao.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
	if knowledgeType != nil {
		query = query.Where("knowledge_type = ?", *knowledgeType)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Order("quality_score DESC, created_at DESC").Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FindByType 根据类型查找知识
func (dao *KnowledgeMemoryDAO) FindByType(
	ctx context.Context,
	tenantID string,
	knowledgeType entity.KnowledgeType,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	var daoModels []model.KnowledgeMemoryDAO
	query := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND knowledge_type = ?", tenantID, knowledgeType).
		Order("quality_score DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FullTextSearch 全文检索
func (dao *KnowledgeMemoryDAO) FullTextSearch(
	ctx context.Context,
	tenantID, query string,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	var daoModels []model.KnowledgeMemoryDAO
	searchQuery := dao.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("MATCH(title, content) AGAINST(? IN NATURAL LANGUAGE MODE)", query).
		Order("quality_score DESC")
	if limit > 0 {
		searchQuery = searchQuery.Limit(limit)
	}
	err := searchQuery.Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FindHighQuality 查找高质量知识
func (dao *KnowledgeMemoryDAO) FindHighQuality(
	ctx context.Context,
	tenantID string,
	minQualityScore float64,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	var daoModels []model.KnowledgeMemoryDAO
	query := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND quality_score >= ?", tenantID, minQualityScore).
		Order("quality_score DESC, access_count DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// Update 更新知识
func (dao *KnowledgeMemoryDAO) Update(ctx context.Context, knowledge *entity.KnowledgeMemory) error {
	daoModel := dao.toModel(knowledge)
	return dao.db.WithContext(ctx).
		Model(&model.KnowledgeMemoryDAO{}).
		Where("memory_id = ?", knowledge.MemoryID).
		Updates(daoModel).Error
}

// UpdateQualityScore 更新质量评分
func (dao *KnowledgeMemoryDAO) UpdateQualityScore(ctx context.Context, memoryID string, score float64) error {
	return dao.db.WithContext(ctx).
		Model(&model.KnowledgeMemoryDAO{}).
		Where("memory_id = ?", memoryID).
		Update("quality_score", score).Error
}

// Delete 删除知识
func (dao *KnowledgeMemoryDAO) Delete(ctx context.Context, memoryID string) error {
	return dao.db.WithContext(ctx).
		Where("memory_id = ?", memoryID).
		Delete(&model.KnowledgeMemoryDAO{}).Error
}

// BatchCreate 批量创建知识
func (dao *KnowledgeMemoryDAO) BatchCreate(ctx context.Context, knowledges []*entity.KnowledgeMemory) error {
	if len(knowledges) == 0 {
		return nil
	}
	daoModels := make([]*model.KnowledgeMemoryDAO, len(knowledges))
	for i, know := range knowledges {
		daoModels[i] = dao.toModel(know)
	}
	return dao.db.WithContext(ctx).CreateInBatches(daoModels, 100).Error
}

// BatchDelete 批量删除知识
func (dao *KnowledgeMemoryDAO) BatchDelete(ctx context.Context, memoryIDs []string) error {
	if len(memoryIDs) == 0 {
		return nil
	}
	return dao.db.WithContext(ctx).
		Where("memory_id IN ?", memoryIDs).
		Delete(&model.KnowledgeMemoryDAO{}).Error
}

// Count 统计知识数量
func (dao *KnowledgeMemoryDAO) Count(
	ctx context.Context,
	tenantID string,
	knowledgeType *entity.KnowledgeType,
) (int64, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&model.KnowledgeMemoryDAO{}).
		Where("tenant_id = ?", tenantID)
	if knowledgeType != nil {
		query = query.Where("knowledge_type = ?", *knowledgeType)
	}
	err := query.Count(&count).Error
	return count, err
}

// IncrementAccessCount 增加访问次数
func (dao *KnowledgeMemoryDAO) IncrementAccessCount(ctx context.Context, memoryID string) error {
	return dao.db.WithContext(ctx).
		Model(&model.KnowledgeMemoryDAO{}).
		Where("memory_id = ?", memoryID).
		UpdateColumn("access_count", gorm.Expr("access_count + 1")).Error
}

// FindBySourceURI 根据来源URI查找知识
func (dao *KnowledgeMemoryDAO) FindBySourceURI(
	ctx context.Context,
	tenantID, sourceURI string,
) ([]*entity.KnowledgeMemory, error) {
	var daoModels []model.KnowledgeMemoryDAO
	err := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND source_uri = ?", tenantID, sourceURI).
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// toModel 转换为DAO Model
func (dao *KnowledgeMemoryDAO) toModel(ent *entity.KnowledgeMemory) *model.KnowledgeMemoryDAO {
	metadataJSON, _ := ent.GetMetadataJSON()
	return &model.KnowledgeMemoryDAO{
		MemoryID:      ent.MemoryID,
		TenantID:      ent.TenantID,
		KnowledgeType: string(ent.KnowledgeType),
		Title:         ent.Title,
		Content:       ent.Content,
		SourceURI:     ent.SourceURI,
		Metadata:      metadataJSON,
		QualityScore:  ent.QualityScore,
		Version:       ent.Version,
		AccessCount:   ent.AccessCount,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}

// toEntity 转换为Entity
func (dao *KnowledgeMemoryDAO) toEntity(daoModel *model.KnowledgeMemoryDAO) *entity.KnowledgeMemory {
	ent := &entity.KnowledgeMemory{
		ID:            daoModel.ID,
		MemoryID:      daoModel.MemoryID,
		TenantID:      daoModel.TenantID,
		KnowledgeType: entity.KnowledgeType(daoModel.KnowledgeType),
		Title:         daoModel.Title,
		Content:       daoModel.Content,
		SourceURI:     daoModel.SourceURI,
		QualityScore:  daoModel.QualityScore,
		Version:       daoModel.Version,
		AccessCount:   daoModel.AccessCount,
		CreatedAt:     daoModel.CreatedAt,
		UpdatedAt:     daoModel.UpdatedAt,
	}
	ent.SetMetadataFromJSON(daoModel.Metadata)
	return ent
}

// batchToEntity 批量转换为Entity
func (dao *KnowledgeMemoryDAO) batchToEntity(daoModels []model.KnowledgeMemoryDAO) []*entity.KnowledgeMemory {
	entities := make([]*entity.KnowledgeMemory, len(daoModels))
	for i, daoModel := range daoModels {
		entities[i] = dao.toEntity(&daoModel)
	}
	return entities
}

// NewMemoryAssociationDAO 创建知识关联DAO
func NewMemoryAssociationDAO(db *gorm.DB) *MemoryAssociationDAO {
	return &MemoryAssociationDAO{db: db}
}

// MemoryAssociationDAO 知识关联数据访问对象
type MemoryAssociationDAO struct {
	db *gorm.DB
}

// Create 创建关联
func (dao *MemoryAssociationDAO) Create(ctx context.Context, association *entity.KnowledgeAssociation) error {
	daoModel := &model.MemoryAssociationDAO{
		SourceMemoryID:  association.SourceMemoryID,
		TargetMemoryID:  association.TargetMemoryID,
		AssociationType: association.AssociationType,
		Strength:        association.Strength,
		CreatedAt:       association.CreatedAt,
	}
	return dao.db.WithContext(ctx).Create(daoModel).Error
}

// FindBySourceMemoryID 查找源知识关联
func (dao *MemoryAssociationDAO) FindBySourceMemoryID(
	ctx context.Context,
	sourceMemoryID string,
) ([]*entity.KnowledgeAssociation, error) {
	var daoModels []model.MemoryAssociationDAO
	err := dao.db.WithContext(ctx).
		Where("source_memory_id = ?", sourceMemoryID).
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToAssociationEntity(daoModels), nil
}

// FindByTargetMemoryID 查找目标知识关联
func (dao *MemoryAssociationDAO) FindByTargetMemoryID(
	ctx context.Context,
	targetMemoryID string,
) ([]*entity.KnowledgeAssociation, error) {
	var daoModels []model.MemoryAssociationDAO
	err := dao.db.WithContext(ctx).
		Where("target_memory_id = ?", targetMemoryID).
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToAssociationEntity(daoModels), nil
}

// Delete 删除关联
func (dao *MemoryAssociationDAO) Delete(ctx context.Context, sourceMemoryID, targetMemoryID string) error {
	return dao.db.WithContext(ctx).
		Where("source_memory_id = ? AND target_memory_id = ?", sourceMemoryID, targetMemoryID).
		Delete(&model.MemoryAssociationDAO{}).Error
}

// UpdateStrength 更新关联强度
func (dao *MemoryAssociationDAO) UpdateStrength(
	ctx context.Context,
	sourceMemoryID, targetMemoryID string,
	strength float64,
) error {
	return dao.db.WithContext(ctx).
		Model(&model.MemoryAssociationDAO{}).
		Where("source_memory_id = ? AND target_memory_id = ?", sourceMemoryID, targetMemoryID).
		Update("strength", strength).Error
}

// BatchCreate 批量创建关联
func (dao *MemoryAssociationDAO) BatchCreate(ctx context.Context, associations []*entity.KnowledgeAssociation) error {
	if len(associations) == 0 {
		return nil
	}
	daoModels := make([]*model.MemoryAssociationDAO, len(associations))
	for i, assoc := range associations {
		daoModels[i] = &model.MemoryAssociationDAO{
			SourceMemoryID:  assoc.SourceMemoryID,
			TargetMemoryID:  assoc.TargetMemoryID,
			AssociationType: assoc.AssociationType,
			Strength:        assoc.Strength,
			CreatedAt:       assoc.CreatedAt,
		}
	}
	return dao.db.WithContext(ctx).CreateInBatches(daoModels, 100).Error
}

// batchToAssociationEntity 批量转换为Entity
func (dao *MemoryAssociationDAO) batchToAssociationEntity(daoModels []model.MemoryAssociationDAO) []*entity.KnowledgeAssociation {
	entities := make([]*entity.KnowledgeAssociation, len(daoModels))
	for i, daoModel := range daoModels {
		entities[i] = &entity.KnowledgeAssociation{
			ID:              daoModel.ID,
			SourceMemoryID:  daoModel.SourceMemoryID,
			TargetMemoryID:  daoModel.TargetMemoryID,
			AssociationType: daoModel.AssociationType,
			Strength:        daoModel.Strength,
			CreatedAt:       daoModel.CreatedAt,
		}
	}
	return entities
}
