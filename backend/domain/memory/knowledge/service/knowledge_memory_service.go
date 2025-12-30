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

package service

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
)

// KnowledgeMemoryService 知识记忆服务接口
type KnowledgeMemoryService interface {
	// StoreKnowledge 存储知识记忆
	StoreKnowledge(ctx context.Context, knowledge *entity.KnowledgeMemory) error

	// RetrieveKnowledge 语义检索知识
	RetrieveKnowledge(
		ctx context.Context,
		tenantID, query string,
		topK int,
	) ([]*entity.KnowledgeWithScore, error)

	// StoreDocument 存储文档(自动分块)
	StoreDocument(
		ctx context.Context,
		docID, title, content string,
		chunkSize int,
	) ([]*entity.KnowledgeMemory, error)

	// FullTextSearch 全文检索
	FullTextSearch(
		ctx context.Context,
		tenantID, query string,
		limit int,
	) ([]*entity.KnowledgeMemory, error)

	// UpdateKnowledge 更新知识
	UpdateKnowledge(
		ctx context.Context,
		memoryID string,
		content string,
	) error

	// DeleteKnowledge 删除知识
	DeleteKnowledge(ctx context.Context, memoryID string) error

	// RateKnowledge 质量评分
	RateKnowledge(
		ctx context.Context,
		memoryID string,
		rating float64,
	) error

	// GetKnowledgeByID 根据ID获取知识
	GetKnowledgeByID(
		ctx context.Context,
		memoryID string,
	) (*entity.KnowledgeMemory, error)

	// ListKnowledgeByType 根据类型列出知识
	ListKnowledgeByType(
		ctx context.Context,
		tenantID string,
		knowledgeType entity.KnowledgeType,
		limit int,
	) ([]*entity.KnowledgeMemory, error)

	// GetKnowledgeStats 获取知识统计
	GetKnowledgeStats(
		ctx context.Context,
		tenantID string,
	) (*KnowledgeStats, error)

	// BatchStoreKnowledge 批量存储知识
	BatchStoreKnowledge(
		ctx context.Context,
		knowledges []*entity.KnowledgeMemory,
	) error

	// GetHighQualityKnowledge 获取高质量知识
	GetHighQualityKnowledge(
		ctx context.Context,
		tenantID string,
		minScore float64,
		limit int,
	) ([]*entity.KnowledgeMemory, error)
}

// KnowledgeAssociationService 知识关联服务接口
type KnowledgeAssociationService interface {
	// CreateAssociation 创建知识关联
	CreateAssociation(
		ctx context.Context,
		sourceMemoryID, targetMemoryID string,
		associationType string,
		strength float64,
	) error

	// GetAssociations 获取知识关联
	GetAssociations(
		ctx context.Context,
		memoryID string,
	) ([]*entity.KnowledgeAssociation, error)

	// BuildKnowledgeGraph 构建知识图谱
	BuildKnowledgeGraph(
		ctx context.Context,
		tenantID string,
	) (*entity.KnowledgeGraph, error)

	// UpdateAssociationStrength 更新关联强度
	UpdateAssociationStrength(
		ctx context.Context,
		sourceMemoryID, targetMemoryID string,
		strength float64,
	) error

	// DeleteAssociation 删除关联
	DeleteAssociation(
		ctx context.Context,
		sourceMemoryID, targetMemoryID string,
	) error

	// FindPath 查找知识之间的路径
	FindPath(
		ctx context.Context,
		startMemoryID, endMemoryID string,
		maxDepth int,
	) ([]string, error)
}

// KnowledgeStats 知识统计
type KnowledgeStats struct {
	TotalKnowledge  int64                        `json:"total_knowledge"`
	ByType          map[entity.KnowledgeType]int64 `json:"by_type"`
	AvgQualityScore float64                      `json:"avg_quality_score"`
	StorageUsedMB   float64                      `json:"storage_used_mb"`
	HighQualityCount int64                       `json:"high_quality_count"`
	LowQualityCount int64                        `json:"low_quality_count"`
}

// DocumentChunkRequest 文档分块请求
type DocumentChunkRequest struct {
	DocumentID    string `json:"document_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	ChunkSize     int    `json:"chunk_size"`     // 分块大小(字符数)
	OverlapSize   int    `json:"overlap_size"`   // 重叠大小
	TenantID      string `json:"tenant_id"`
}

// DocumentChunkResponse 文档分块响应
type DocumentChunkResponse struct {
	Chunks     []*entity.KnowledgeMemory `json:"chunks"`
	TotalCount int                       `json:"total_count"`
	DocumentID string                    `json:"document_id"`
}
