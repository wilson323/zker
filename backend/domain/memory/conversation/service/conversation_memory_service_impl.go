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

	"github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	"github.com/coze-dev/coze-studio/backend/domain/memory/internal/dal"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
	vectordb "github.com/coze-dev/coze-studio/backend/infra/vectordb"
)

// NewConversationMemoryService 创建对话记忆服务
func NewConversationMemoryService(
	db *gorm.DB,
	vectorClient vectordb.VectorClient,
	llmClient llm.MemoryLLMClient,
) ConversationMemoryService {
	return &conversationMemoryServiceImpl{
		memoryDAL:    dal.NewConversationMemoryDAO(db),
		vectorClient: vectorClient,
		llmClient:    llmClient,
	}
}

type conversationMemoryServiceImpl struct {
	memoryDAL    *dal.ConversationMemoryDAO
	vectorClient vectordb.VectorClient
	llmClient    llm.MemoryLLMClient
}

// StoreMemory 存储对话记忆
func (s *conversationMemoryServiceImpl) StoreMemory(
	ctx context.Context,
	memory *entity.ConversationMemory,
) error {
	// 1. 生成向量
	vector, err := s.llmClient.Embed(ctx, memory.Content)
	if err != nil {
		return fmt.Errorf("failed to embed content: %w", err)
	}
	memory.Embedding = vector

	// 2. 生成MemoryID
	if memory.MemoryID == "" {
		memory.MemoryID = uuid.New().String()
	}

	// 3. 设置默认值
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}
	memory.UpdatedAt = time.Now()

	// 4. 计算重要性评分
	if memory.ImportanceScore == 0 {
		memory.ImportanceScore = s.calculateImportance(memory)
	}

	// 5. 存储到向量数据库
	vectorID := memory.MemoryID
	err = s.vectorClient.InsertVectors(
		ctx,
		"conversation_memories",
		[]string{vectorID},
		[][]float32{memory.Embedding},
	)
	if err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}

	// 6. 存储到MySQL
	return s.memoryDAL.Create(ctx, memory)
}

// RetrieveMemories 语义检索记忆
func (s *conversationMemoryServiceImpl) RetrieveMemories(
	ctx context.Context,
	userID, query string,
	topK int,
) ([]*entity.MemoryWithScore, error) {
	// 1. 向量化查询
	vector, err := s.llmClient.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// 2. 向量搜索
	searchResult, err := s.vectorClient.SearchVectors(
		ctx,
		"conversation_memories",
		vector,
		topK,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	// 3. 根据vectorID查找记忆
	memories, err := s.memoryDAL.FindByVectorIDs(ctx, searchResult.IDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find memories: %w", err)
	}

	// 4. 组合结果
	result := make([]*entity.MemoryWithScore, len(memories))
	for i, mem := range memories {
		result[i] = &entity.MemoryWithScore{
			Memory: mem,
			Score:  searchResult.Scores[i],
		}
		// 更新访问统计
		go s.memoryDAL.UpdateAccessCount(context.Background(), mem.MemoryID)
	}

	return result, nil
}

// GenerateSummary 生成对话摘要
func (s *conversationMemoryServiceImpl) GenerateSummary(
	ctx context.Context,
	conversationID string,
	targetLength int,
) (*entity.SummaryResult, error) {
	// TODO: 实现对话摘要
	return &entity.SummaryResult{
		Summary:          "对话摘要",
		OriginalCount:    10,
		CompressedCount:  1,
		KeyFacts:         []string{},
		CompressionRatio: 0.9,
		NewMemoryID:      uuid.New().String(),
	}, nil
}

// ExtractEntities 提取实体
func (s *conversationMemoryServiceImpl) ExtractEntities(
	ctx context.Context,
	conversationID string,
) (*entity.ExtractEntityResult, error) {
	// TODO: 实现实体提取
	return &entity.ExtractEntityResult{
		Entities:   []entity.Entity{},
		Confidence: 0.9,
	}, nil
}

// ExtractPreferences 提取偏好
func (s *conversationMemoryServiceImpl) ExtractPreferences(
	ctx context.Context,
	userID string,
) (*entity.ExtractPreferenceResult, error) {
	// TODO: 实现偏好提取
	return &entity.ExtractPreferenceResult{
		Preferences: []entity.Preference{},
		Confidence:  0.85,
	}, nil
}

// GetConversationHistory 获取对话历史
func (s *conversationMemoryServiceImpl) GetConversationHistory(
	ctx context.Context,
	tenantID, userID, conversationID string,
	limit int,
) ([]*entity.ConversationMemory, error) {
	return s.memoryDAL.FindByConversationID(ctx, tenantID, userID, conversationID)
}

// UpdateMemory 更新记忆
func (s *conversationMemoryServiceImpl) UpdateMemory(
	ctx context.Context,
	memoryID string,
	content string,
) error {
	memory, err := s.memoryDAL.FindByID(ctx, memoryID)
	if err != nil {
		return err
	}
	memory.Content = content
	memory.UpdatedAt = time.Now()
	return s.memoryDAL.Update(ctx, memory)
}

// DeleteMemory 删除记忆
func (s *conversationMemoryServiceImpl) DeleteMemory(ctx context.Context, memoryID string) error {
	// 1. 删除向量
	err := s.vectorClient.DeleteVectors(ctx, "conversation_memories", []string{memoryID})
	if err != nil {
		return fmt.Errorf("failed to delete vector: %w", err)
	}
	// 2. 删除记忆
	return s.memoryDAL.Delete(ctx, memoryID)
}

// DeleteConversationMemories 删除会话的所有记忆
func (s *conversationMemoryServiceImpl) DeleteConversationMemories(
	ctx context.Context,
	conversationID string,
) error {
	// TODO: 批量删除向量
	return s.memoryDAL.DeleteByConversationID(ctx, conversationID)
}

// CompressMemories 压缩记忆
func (s *conversationMemoryServiceImpl) CompressMemories(
	ctx context.Context,
	tenantID, userID string,
) error {
	// TODO: 实现记忆压缩
	return nil
}

// DeleteExpiredMemories 删除过期记忆
func (s *conversationMemoryServiceImpl) DeleteExpiredMemories(ctx context.Context, tenantID string) error {
	_, err := s.memoryDAL.DeleteExpired(ctx, tenantID)
	return err
}

// GetMemoryStats 获取记忆统计
func (s *conversationMemoryServiceImpl) GetMemoryStats(
	ctx context.Context,
	tenantID, userID string,
) (*MemoryStats, error) {
	total, err := s.memoryDAL.Count(ctx, tenantID, userID, nil)
	if err != nil {
		return nil, err
	}
	return &MemoryStats{
		TotalMemories: total,
		ByType:         make(map[entity.MemoryType]int64),
		StorageUsedMB:  0,
		ExpiredCount:   0,
	}, nil
}

// GetHotMemories 获取热点记忆
func (s *conversationMemoryServiceImpl) GetHotMemories(
	ctx context.Context,
	tenantID, userID string,
	limit int,
) ([]*entity.ConversationMemory, error) {
	return s.memoryDAL.FindHotMemories(ctx, tenantID, userID, limit)
}

// BatchStoreMemories 批量存储记忆
func (s *conversationMemoryServiceImpl) BatchStoreMemories(
	ctx context.Context,
	memories []*entity.ConversationMemory,
) error {
	for _, mem := range memories {
		if err := s.StoreMemory(ctx, mem); err != nil {
			return err
		}
	}
	return nil
}

// calculateImportance 计算重要性评分
func (s *conversationMemoryServiceImpl) calculateImportance(memory *entity.ConversationMemory) float64 {
	// 简单的重要性评分算法
	score := 0.5
	switch memory.MemoryType {
	case entity.MemoryTypeSummary:
		score = 0.9
	case entity.MemoryTypeEntity:
		score = 0.7
	case entity.MemoryTypePreference:
		score = 0.8
	case entity.MemoryTypeEvent:
		score = 0.6
	}
	return score
}
