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

	"github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
)

// ConversationMemoryService 对话记忆服务接口
type ConversationMemoryService interface {
	// StoreMemory 存储对话记忆
	StoreMemory(ctx context.Context, memory *entity.ConversationMemory) error

	// RetrieveMemories 语义检索记忆
	RetrieveMemories(
		ctx context.Context,
		userID, query string,
		topK int,
	) ([]*entity.MemoryWithScore, error)

	// GenerateSummary 生成对话摘要
	GenerateSummary(
		ctx context.Context,
		conversationID string,
		targetLength int,
	) (*entity.SummaryResult, error)

	// ExtractEntities 提取实体
	ExtractEntities(
		ctx context.Context,
		conversationID string,
	) (*entity.ExtractEntityResult, error)

	// ExtractPreferences 提取偏好
	ExtractPreferences(
		ctx context.Context,
		userID string,
	) (*entity.ExtractPreferenceResult, error)

	// GetConversationHistory 获取对话历史
	GetConversationHistory(
		ctx context.Context,
		tenantID, userID, conversationID string,
		limit int,
	) ([]*entity.ConversationMemory, error)

	// UpdateMemory 更新记忆
	UpdateMemory(
		ctx context.Context,
		memoryID string,
		content string,
	) error

	// DeleteMemory 删除记忆
	DeleteMemory(ctx context.Context, memoryID string) error

	// DeleteConversationMemories 删除会话的所有记忆
	DeleteConversationMemories(
		ctx context.Context,
		conversationID string,
	) error

	// CompressMemories 压缩记忆
	CompressMemories(
		ctx context.Context,
		tenantID, userID string,
	) error

	// DeleteExpiredMemories 删除过期记忆
	DeleteExpiredMemories(ctx context.Context, tenantID string) error

	// GetMemoryStats 获取记忆统计
	GetMemoryStats(
		ctx context.Context,
		tenantID, userID string,
	) (*MemoryStats, error)

	// GetHotMemories 获取热点记忆
	GetHotMemories(
		ctx context.Context,
		tenantID, userID string,
		limit int,
	) ([]*entity.ConversationMemory, error)

	// BatchStoreMemories 批量存储记忆
	BatchStoreMemories(
		ctx context.Context,
		memories []*entity.ConversationMemory,
	) error
}

// MemoryStats 记忆统计
type MemoryStats struct {
	TotalMemories   int64                    `json:"total_memories"`
	ByType          map[entity.MemoryType]int64 `json:"by_type"`
	StorageUsedMB   float64                  `json:"storage_used_mb"`
	ExpiredCount    int64                    `json:"expired_count"`
	OldestMemory    string                   `json:"oldest_memory"`
	NewestMemory    string                   `json:"newest_memory"`
}

// MemoryRetrievalService 记忆检索服务接口
type MemoryRetrievalService interface {
	// SemanticSearch 语义搜索
	SemanticSearch(
		ctx context.Context,
		tenantID, userID, query string,
		memoryTypes []entity.MemoryType,
		topK int,
		scoreThreshold float64,
	) ([]*entity.MemoryWithScore, error)

	// KeywordSearch 关键词搜索
	KeywordSearch(
		ctx context.Context,
		tenantID, userID, keyword string,
		memoryTypes []entity.MemoryType,
		limit int,
	) ([]*entity.ConversationMemory, error)

	// TimeRangeSearch 时间范围搜索
	TimeRangeSearch(
		ctx context.Context,
		tenantID, userID string,
		startTime, endTime int64,
		memoryTypes []entity.MemoryType,
	) ([]*entity.ConversationMemory, error)

	// HybridSearch 混合搜索(语义+关键词+时间)
	HybridSearch(
		ctx context.Context,
		req *HybridSearchRequest,
	) ([]*entity.MemoryWithScore, error)
}

// HybridSearchRequest 混合搜索请求
type HybridSearchRequest struct {
	TenantID        string                 `json:"tenant_id"`
	UserID          string                 `json:"user_id"`
	Query           string                 `json:"query"`
	Keywords        []string               `json:"keywords"`
	MemoryTypes     []entity.MemoryType    `json:"memory_types"`
	StartTime       *int64                 `json:"start_time"`
	EndTime         *int64                 `json:"end_time"`
	TopK            int                    `json:"top_k"`
	ScoreThreshold  float64                `json:"score_threshold"`
	SemanticWeight  float64                `json:"semantic_weight"`  // 语义权重
	KeywordWeight   float64                `json:"keyword_weight"`   // 关键词权重
}
