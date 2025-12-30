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

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
)

// KnowledgeMemoryRepository 知识记忆仓储接口
type KnowledgeMemoryRepository interface {
	// Create 创建知识
	Create(ctx context.Context, knowledge *entity.KnowledgeMemory) error

	// FindByID 根据ID查找知识
	FindByID(ctx context.Context, memoryID string) (*entity.KnowledgeMemory, error)

	// FindByTenantID 根据租户ID查找知识
	FindByTenantID(
		ctx context.Context,
		tenantID string,
		knowledgeType *entity.KnowledgeType,
		limit int,
	) ([]*entity.KnowledgeMemory, error)

	// FindByType 根据类型查找知识
	FindByType(
		ctx context.Context,
		tenantID string,
		knowledgeType entity.KnowledgeType,
		limit int,
	) ([]*entity.KnowledgeMemory, error)

	// FindByVectorIDs 根据向量ID列表查找知识
	FindByVectorIDs(
		ctx context.Context,
		vectorIDs []string,
	) ([]*entity.KnowledgeMemory, error)

	// FullTextSearch 全文检索
	FullTextSearch(
		ctx context.Context,
		tenantID, query string,
		limit int,
	) ([]*entity.KnowledgeMemory, error)

	// FindHighQuality 查找高质量知识
	FindHighQuality(
		ctx context.Context,
		tenantID string,
		minQualityScore float64,
		limit int,
	) ([]*entity.KnowledgeMemory, error)

	// Update 更新知识
	Update(ctx context.Context, knowledge *entity.KnowledgeMemory) error

	// UpdateQualityScore 更新质量评分
	UpdateQualityScore(
		ctx context.Context,
		memoryID string,
		score float64,
	) error

	// Delete 删除知识
	Delete(ctx context.Context, memoryID string) error

	// BatchCreate 批量创建知识
	BatchCreate(
		ctx context.Context,
		knowledges []*entity.KnowledgeMemory,
	) error

	// BatchDelete 批量删除知识
	BatchDelete(ctx context.Context, memoryIDs []string) error

	// Count 统计知识数量
	Count(
		ctx context.Context,
		tenantID string,
		knowledgeType *entity.KnowledgeType,
	) (int64, error)

	// IncrementAccessCount 增加访问次数
	IncrementAccessCount(ctx context.Context, memoryID string) error

	// FindBySourceURI 根据来源URI查找知识
	FindBySourceURI(
		ctx context.Context,
		tenantID, sourceURI string,
	) ([]*entity.KnowledgeMemory, error)
}

// KnowledgeAssociationRepository 知识关联仓储接口
type KnowledgeAssociationRepository interface {
	// Create 创建关联
	Create(ctx context.Context, association *entity.KnowledgeAssociation) error

	// FindBySourceMemoryID 查找源知识关联
	FindBySourceMemoryID(
		ctx context.Context,
		sourceMemoryID string,
	) ([]*entity.KnowledgeAssociation, error)

	// FindByTargetMemoryID 查找目标知识关联
	FindByTargetMemoryID(
		ctx context.Context,
		targetMemoryID string,
	) ([]*entity.KnowledgeAssociation, error)

	// FindByType 根据类型查找关联
	FindByType(
		ctx context.Context,
		associationType string,
	) ([]*entity.KnowledgeAssociation, error)

	// Delete 删除关联
	Delete(ctx context.Context, sourceMemoryID, targetMemoryID string) error

	// UpdateStrength 更新关联强度
	UpdateStrength(
		ctx context.Context,
		sourceMemoryID, targetMemoryID string,
		strength float64,
	) error

	// BatchCreate 批量创建关联
	BatchCreate(
		ctx context.Context,
		associations []*entity.KnowledgeAssociation,
	) error
}
