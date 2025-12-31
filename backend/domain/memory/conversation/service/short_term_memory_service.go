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

// ShortTermMemoryService 短期记忆服务接口
type ShortTermMemoryService interface {
	// Store 存储短期记忆上下文
	Store(ctx context.Context, req *entity.StoreShortTermMemoryRequest) error

	// Get 获取短期记忆上下文
	Get(ctx context.Context, conversationID string) (*entity.GetShortTermMemoryResponse, error)

	// Delete 删除短期记忆
	Delete(ctx context.Context, conversationID string) error

	// Compress 压缩上下文
	Compress(ctx context.Context, req *entity.CompressionRequest) (*entity.CompressionResult, error)

	// AutoCompress 自动压缩(当超过窗口大小时)
	AutoCompress(ctx context.Context, conversationID string) (*entity.CompressionResult, error)

	// AddMessage 添加消息到上下文
	AddMessage(ctx context.Context, conversationID string, msg *entity.Message) error

	// ClearExpired 清理过期的短期记忆
	ClearExpired(ctx context.Context) error

	// GetWindowStats 获取窗口统计信息
	GetWindowStats(ctx context.Context, conversationID string) (*WindowStats, error)
}

// WindowStats 窗口统计信息
type WindowStats struct {
	ConversationID    string    `json:"conversation_id"`
	MessageCount      int       `json:"message_count"`
	TotalTokens       int       `json:"total_tokens"`
	MaxTokens         int       `json:"max_tokens"`
	UtilizationRatio  float64   `json:"utilization_ratio"` // 利用率
	NeedsCompression  bool      `json:"needs_compression"`
	ExpiresAt         *int64    `json:"expires_at,omitempty"`
}
