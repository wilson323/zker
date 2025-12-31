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
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
	"github.com/coze-dev/coze-studio/backend/domain/memory/internal/dal"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
	vectordb "github.com/coze-dev/coze-studio/backend/infra/vectordb"
)

// NewKnowledgeMemoryService 创建知识记忆服务
func NewKnowledgeMemoryService(
	db *gorm.DB,
	vectorClient vectordb.VectorClient,
	llmClient llm.MemoryLLMClient,
) KnowledgeMemoryService {
	return &knowledgeMemoryServiceImpl{
		knowledgeDAL: dal.NewKnowledgeMemoryDAO(db),
		vectorClient: vectorClient,
		llmClient:     llmClient,
	}
}

type knowledgeMemoryServiceImpl struct {
	knowledgeDAL *dal.KnowledgeMemoryDAO
	vectorClient vectordb.VectorClient
	llmClient     llm.MemoryLLMClient
}

// StoreKnowledge 存储知识记忆
func (s *knowledgeMemoryServiceImpl) StoreKnowledge(
	ctx context.Context,
	knowledge *entity.KnowledgeMemory,
) error {
	// 1. 生成向量
	vector, err := s.llmClient.Embed(ctx, knowledge.Content)
	if err != nil {
		return fmt.Errorf("failed to embed content: %w", err)
	}
	knowledge.Embedding = vector

	// 2. 生成MemoryID
	if knowledge.MemoryID == "" {
		knowledge.MemoryID = uuid.New().String()
	}

	// 3. 设置默认值
	if knowledge.CreatedAt.IsZero() {
		knowledge.CreatedAt = time.Now()
	}
	knowledge.UpdatedAt = time.Now()

	// 4. 存储到向量数据库
	err = s.vectorClient.InsertVectors(
		ctx,
		"knowledge_memories",
		[]string{knowledge.MemoryID},
		[][]float32{knowledge.Embedding},
	)
	if err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}

	// 5. 存储到MySQL
	return s.knowledgeDAL.Create(ctx, knowledge)
}

// RetrieveKnowledge 语义检索知识
func (s *knowledgeMemoryServiceImpl) RetrieveKnowledge(
	ctx context.Context,
	tenantID, query string,
	topK int,
) ([]*entity.KnowledgeWithScore, error) {
	// 1. 向量化查询
	vector, err := s.llmClient.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// 2. 向量搜索
	searchResult, err := s.vectorClient.SearchVectors(
		ctx,
		"knowledge_memories",
		vector,
		topK,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	// 3. 根据vectorID查找知识
	knowledges, err := s.knowledgeDAL.FindByVectorIDs(ctx, searchResult.IDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledges: %w", err)
	}

	// 4. 组合结果
	result := make([]*entity.KnowledgeWithScore, len(knowledges))
	for i, know := range knowledges {
		result[i] = &entity.KnowledgeWithScore{
			Knowledge: know,
			Score:     searchResult.Scores[i],
		}
		// 更新访问统计
		go s.knowledgeDAL.IncrementAccessCount(context.Background(), know.MemoryID)
	}

	return result, nil
}

// StoreDocument 存储文档(自动分块)
func (s *knowledgeMemoryServiceImpl) StoreDocument(
	ctx context.Context,
	docID, title, content string,
	chunkSize int,
) ([]*entity.KnowledgeMemory, error) {
	// 1. 文档分块
	chunks := s.chunkDocument(content, chunkSize)

	// 2. 为每个分块创建知识记忆
	knowledges := make([]*entity.KnowledgeMemory, len(chunks))
	for i, chunk := range chunks {
		knowledges[i] = &entity.KnowledgeMemory{
			MemoryID:      uuid.New().String(),
			TenantID:      "", // 需要从context获取
			KnowledgeType: entity.KnowledgeTypeDocument,
			Title:         fmt.Sprintf("%s (Part %d)", title, i+1),
			Content:       chunk,
			SourceURI:     docID,
			QualityScore:  0.5,
			Version:       1,
		}
	}

	// 3. 批量存储
	for _, know := range knowledges {
		if err := s.StoreKnowledge(ctx, know); err != nil {
			return nil, err
		}
	}

	return knowledges, nil
}

// FullTextSearch 全文检索
func (s *knowledgeMemoryServiceImpl) FullTextSearch(
	ctx context.Context,
	tenantID, query string,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	return s.knowledgeDAL.FullTextSearch(ctx, tenantID, query, limit)
}

// UpdateKnowledge 更新知识
func (s *knowledgeMemoryServiceImpl) UpdateKnowledge(
	ctx context.Context,
	memoryID string,
	content string,
) error {
	knowledge, err := s.knowledgeDAL.FindByID(ctx, memoryID)
	if err != nil {
		return err
	}
	knowledge.Content = content
	knowledge.UpdatedAt = time.Now()
	knowledge.Version++

	// 重新生成向量
	vector, err := s.llmClient.Embed(ctx, content)
	if err != nil {
		return err
	}
	knowledge.Embedding = vector

	return s.knowledgeDAL.Update(ctx, knowledge)
}

// DeleteKnowledge 删除知识
func (s *knowledgeMemoryServiceImpl) DeleteKnowledge(ctx context.Context, memoryID string) error {
	// 1. 删除向量
	err := s.vectorClient.DeleteVectors(ctx, "knowledge_memories", []string{memoryID})
	if err != nil {
		return fmt.Errorf("failed to delete vector: %w", err)
	}
	// 2. 删除知识
	return s.knowledgeDAL.Delete(ctx, memoryID)
}

// RateKnowledge 质量评分
func (s *knowledgeMemoryServiceImpl) RateKnowledge(
	ctx context.Context,
	memoryID string,
	rating float64,
) error {
	return s.knowledgeDAL.UpdateQualityScore(ctx, memoryID, rating)
}

// GetKnowledgeByID 根据ID获取知识
func (s *knowledgeMemoryServiceImpl) GetKnowledgeByID(
	ctx context.Context,
	memoryID string,
) (*entity.KnowledgeMemory, error) {
	return s.knowledgeDAL.FindByID(ctx, memoryID)
}

// ListKnowledgeByType 根据类型列出知识
func (s *knowledgeMemoryServiceImpl) ListKnowledgeByType(
	ctx context.Context,
	tenantID string,
	knowledgeType entity.KnowledgeType,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	return s.knowledgeDAL.FindByType(ctx, tenantID, knowledgeType, limit)
}

// GetKnowledgeStats 获取知识统计
func (s *knowledgeMemoryServiceImpl) GetKnowledgeStats(
	ctx context.Context,
	tenantID string,
) (*KnowledgeStats, error) {
	total, err := s.knowledgeDAL.Count(ctx, tenantID, nil)
	if err != nil {
		return nil, err
	}
	return &KnowledgeStats{
		TotalKnowledge:  total,
		ByType:          make(map[entity.KnowledgeType]int64),
		AvgQualityScore: 0.5,
		StorageUsedMB:   0,
	}, nil
}

// BatchStoreKnowledge 批量存储知识
func (s *knowledgeMemoryServiceImpl) BatchStoreKnowledge(
	ctx context.Context,
	knowledges []*entity.KnowledgeMemory,
) error {
	for _, know := range knowledges {
		if err := s.StoreKnowledge(ctx, know); err != nil {
			return err
		}
	}
	return nil
}

// GetHighQualityKnowledge 获取高质量知识
func (s *knowledgeMemoryServiceImpl) GetHighQualityKnowledge(
	ctx context.Context,
	tenantID string,
	minScore float64,
	limit int,
) ([]*entity.KnowledgeMemory, error) {
	return s.knowledgeDAL.FindHighQuality(ctx, tenantID, minScore, limit)
}

// chunkDocument 文档分块
func (s *knowledgeMemoryServiceImpl) chunkDocument(content string, chunkSize int) []string {
	// 简单分块算法
	chunks := []string{}
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}
		chunks = append(chunks, content[i:end])
	}
	return chunks
}
