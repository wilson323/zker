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
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	"github.com/coze-dev/coze-studio/backend/domain/memory/internal/dal/model"
)

// NewConversationMemoryDAO 创建对话记忆DAO
func NewConversationMemoryDAO(db *gorm.DB) *ConversationMemoryDAO {
	return &ConversationMemoryDAO{db: db}
}

// ConversationMemoryDAO 对话记忆数据访问对象
type ConversationMemoryDAO struct {
	db *gorm.DB
}

// Create 创建记忆
func (dao *ConversationMemoryDAO) Create(ctx context.Context, memory *entity.ConversationMemory) error {
	daoModel := dao.toModel(memory)
	result := dao.db.WithContext(ctx).Create(daoModel)
	if result.Error != nil {
		return result.Error
	}
	memory.ID = daoModel.ID
	return nil
}

// FindByID 根据ID查找记忆
func (dao *ConversationMemoryDAO) FindByID(ctx context.Context, memoryID string) (*entity.ConversationMemory, error) {
	var daoModel model.ConversationMemoryDAO
	err := dao.db.WithContext(ctx).Where("memory_id = ?", memoryID).First(&daoModel).Error
	if err != nil {
		return nil, err
	}
	return dao.toEntity(&daoModel), nil
}

// FindByConversationID 根据会话ID查找记忆
func (dao *ConversationMemoryDAO) FindByConversationID(
	ctx context.Context,
	tenantID, userID, conversationID string,
) ([]*entity.ConversationMemory, error) {
	var daoModels []model.ConversationMemoryDAO
	err := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ? AND conversation_id = ?", tenantID, userID, conversationID).
		Order("created_at DESC").
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FindByType 根据类型查找记忆
func (dao *ConversationMemoryDAO) FindByType(
	ctx context.Context,
	tenantID, userID string,
	memoryType entity.MemoryType,
	limit int,
) ([]*entity.ConversationMemory, error) {
	var daoModels []model.ConversationMemoryDAO
	query := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ? AND memory_type = ?", tenantID, userID, memoryType).
		Order("importance_score DESC, created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FindExpired 查找过期记忆
func (dao *ConversationMemoryDAO) FindExpired(ctx context.Context, tenantID string) ([]*entity.ConversationMemory, error) {
	var daoModels []model.ConversationMemoryDAO
	err := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND expires_at < NOW()", tenantID).
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FindByVectorIDs 根据向量ID列表查找记忆
// 保持向量搜索结果的顺序返回
func (dao *ConversationMemoryDAO) FindByVectorIDs(
	ctx context.Context,
	vectorIDs []string,
) ([]*entity.ConversationMemory, error) {
	if len(vectorIDs) == 0 {
		return []*entity.ConversationMemory{}, nil
	}

	var daoModels []model.ConversationMemoryDAO
	err := dao.db.WithContext(ctx).
		Where("memory_id IN ?", vectorIDs).
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}

	// 构建内存ID到实体的映射，用于保持vectorIDs的顺序
	memoryMap := make(map[string]*entity.ConversationMemory, len(daoModels))
	for i := range daoModels {
		memory := dao.toEntity(&daoModels[i])
		memoryMap[memory.MemoryID] = memory
	}

	// 按照vectorIDs的顺序返回结果
	result := make([]*entity.ConversationMemory, 0, len(vectorIDs))
	for _, vectorID := range vectorIDs {
		if memory, exists := memoryMap[vectorID]; exists {
			result = append(result, memory)
		}
	}

	return result, nil
}

// Update 更新记忆
func (dao *ConversationMemoryDAO) Update(ctx context.Context, memory *entity.ConversationMemory) error {
	daoModel := dao.toModel(memory)
	return dao.db.WithContext(ctx).
		Model(&model.ConversationMemoryDAO{}).
		Where("memory_id = ?", memory.MemoryID).
		Updates(daoModel).Error
}

// UpdateAccessCount 更新访问次数
func (dao *ConversationMemoryDAO) UpdateAccessCount(ctx context.Context, memoryID string) error {
	return dao.db.WithContext(ctx).
		Model(&model.ConversationMemoryDAO{}).
		Where("memory_id = ?", memoryID).
		Updates(map[string]interface{}{
			"access_count":      gorm.Expr("access_count + 1"),
			"last_accessed_at": time.Now(),
		}).Error
}

// Delete 删除记忆
func (dao *ConversationMemoryDAO) Delete(ctx context.Context, memoryID string) error {
	return dao.db.WithContext(ctx).
		Where("memory_id = ?", memoryID).
		Delete(&model.ConversationMemoryDAO{}).Error
}

// DeleteByConversationID 删除会话的所有记忆
func (dao *ConversationMemoryDAO) DeleteByConversationID(ctx context.Context, conversationID string) error {
	return dao.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Delete(&model.ConversationMemoryDAO{}).Error
}

// DeleteExpired 删除过期记忆
func (dao *ConversationMemoryDAO) DeleteExpired(ctx context.Context, tenantID string) (int64, error) {
	result := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND expires_at < NOW()", tenantID).
		Delete(&model.ConversationMemoryDAO{})
	return result.RowsAffected, result.Error
}

// BatchCreate 批量创建记忆
func (dao *ConversationMemoryDAO) BatchCreate(ctx context.Context, memories []*entity.ConversationMemory) error {
	if len(memories) == 0 {
		return nil
	}
	daoModels := make([]*model.ConversationMemoryDAO, len(memories))
	for i, mem := range memories {
		daoModels[i] = dao.toModel(mem)
	}
	return dao.db.WithContext(ctx).CreateInBatches(daoModels, 100).Error
}

// BatchDelete 批量删除记忆
func (dao *ConversationMemoryDAO) BatchDelete(ctx context.Context, memoryIDs []string) error {
	if len(memoryIDs) == 0 {
		return nil
	}
	return dao.db.WithContext(ctx).
		Where("memory_id IN ?", memoryIDs).
		Delete(&model.ConversationMemoryDAO{}).Error
}

// Count 统计记忆数量
func (dao *ConversationMemoryDAO) Count(
	ctx context.Context,
	tenantID, userID string,
	memoryType *entity.MemoryType,
) (int64, error) {
	var count int64
	query := dao.db.WithContext(ctx).Model(&model.ConversationMemoryDAO{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID)
	if memoryType != nil {
		query = query.Where("memory_type = ?", *memoryType)
	}
	err := query.Count(&count).Error
	return count, err
}

// FindHotMemories 查找热点记忆
func (dao *ConversationMemoryDAO) FindHotMemories(
	ctx context.Context,
	tenantID, userID string,
	limit int,
) ([]*entity.ConversationMemory, error) {
	var daoModels []model.ConversationMemoryDAO
	err := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Order("access_count DESC, last_accessed_at DESC").
		Limit(limit).
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// FindByTimeRange 按时间范围查找记忆
func (dao *ConversationMemoryDAO) FindByTimeRange(
	ctx context.Context,
	tenantID, userID string,
	startTime, endTime time.Time,
) ([]*entity.ConversationMemory, error) {
	var daoModels []model.ConversationMemoryDAO
	err := dao.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ? AND created_at >= ? AND created_at <= ?", tenantID, userID, startTime, endTime).
		Order("created_at DESC").
		Find(&daoModels).Error
	if err != nil {
		return nil, err
	}
	return dao.batchToEntity(daoModels), nil
}

// toModel 转换为DAO Model
func (dao *ConversationMemoryDAO) toModel(entity *entity.ConversationMemory) *model.ConversationMemoryDAO {
	metadataJSON, _ := entity.GetMetadataJSON()
	return &model.ConversationMemoryDAO{
		MemoryID:        entity.MemoryID,
		TenantID:        entity.TenantID,
		UserID:          entity.UserID,
		ConversationID:  entity.ConversationID,
		MemoryType:      string(entity.MemoryType),
		Content:         entity.Content,
		ImportanceScore: entity.ImportanceScore,
		AccessCount:     entity.AccessCount,
		LastAccessedAt:  entity.LastAccessedAt,
		ExpiresAt:       entity.ExpiresAt,
		Metadata:        metadataJSON,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
	}
}

// toEntity 转换为Entity
func (dao *ConversationMemoryDAO) toEntity(daoModel *model.ConversationMemoryDAO) *entity.ConversationMemory {
	entity := &entity.ConversationMemory{
		ID:              daoModel.ID,
		MemoryID:        daoModel.MemoryID,
		TenantID:        daoModel.TenantID,
		UserID:          daoModel.UserID,
		ConversationID:  daoModel.ConversationID,
		MemoryType:      entity.MemoryType(daoModel.MemoryType),
		Content:         daoModel.Content,
		ImportanceScore: daoModel.ImportanceScore,
		AccessCount:     daoModel.AccessCount,
		LastAccessedAt:  daoModel.LastAccessedAt,
		ExpiresAt:       daoModel.ExpiresAt,
		CreatedAt:       daoModel.CreatedAt,
		UpdatedAt:       daoModel.UpdatedAt,
	}
	entity.SetMetadataFromJSON(daoModel.Metadata)
	return entity
}

// batchToEntity 批量转换为Entity
func (dao *ConversationMemoryDAO) batchToEntity(daoModels []model.ConversationMemoryDAO) []*entity.ConversationMemory {
	entities := make([]*entity.ConversationMemory, len(daoModels))
	for i, daoModel := range daoModels {
		entities[i] = dao.toEntity(&daoModel)
	}
	return entities
}
