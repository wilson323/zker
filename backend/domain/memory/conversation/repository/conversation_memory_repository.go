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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
)

// ConversationMemoryRepository 对话记忆仓储接口
type ConversationMemoryRepository interface {
	// Create 创建记忆
	Create(ctx context.Context, memory *entity.ConversationMemory) error

	// FindByID 根据ID查找记忆
	FindByID(ctx context.Context, memoryID string) (*entity.ConversationMemory, error)

	// FindByConversationID 根据会话ID查找记忆
	FindByConversationID(
		ctx context.Context,
		tenantID, userID, conversationID string,
	) ([]*entity.ConversationMemory, error)

	// FindByType 根据类型查找记忆
	FindByType(
		ctx context.Context,
		tenantID, userID string,
		memoryType entity.MemoryType,
		limit int,
	) ([]*entity.ConversationMemory, error)

	// FindExpired 查找过期记忆
	FindExpired(ctx context.Context, tenantID string) ([]*entity.ConversationMemory, error)

	// FindByVectorIDs 根据向量ID列表查找记忆
	FindByVectorIDs(
		ctx context.Context,
		vectorIDs []string,
	) ([]*entity.ConversationMemory, error)

	// Update 更新记忆
	Update(ctx context.Context, memory *entity.ConversationMemory) error

	// UpdateAccessCount 更新访问次数和最后访问时间
	UpdateAccessCount(ctx context.Context, memoryID string) error

	// Delete 删除记忆
	Delete(ctx context.Context, memoryID string) error

	// DeleteByConversationID 删除会话的所有记忆
	DeleteByConversationID(ctx context.Context, conversationID string) error

	// DeleteExpired 删除过期记忆
	DeleteExpired(ctx context.Context, tenantID string) (int64, error)

	// BatchCreate 批量创建记忆
	BatchCreate(
		ctx context.Context,
		memories []*entity.ConversationMemory,
	) error

	// BatchDelete 批量删除记忆
	BatchDelete(ctx context.Context, memoryIDs []string) error

	// Count 统计记忆数量
	Count(
		ctx context.Context,
		tenantID, userID string,
		memoryType *entity.MemoryType,
	) (int64, error)

	// FindHotMemories 查找热点记忆(访问次数高)
	FindHotMemories(
		ctx context.Context,
		tenantID, userID string,
		limit int,
	) ([]*entity.ConversationMemory, error)

	// FindByTimeRange 按时间范围查找记忆
	FindByTimeRange(
		ctx context.Context,
		tenantID, userID string,
		startTime, endTime time.Time,
	) ([]*entity.ConversationMemory, error)
}
