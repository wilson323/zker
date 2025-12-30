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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	convEntity "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
)

// MockVectorClient 模拟向量客户端
type MockVectorClient struct {
	mock.Mock
}

func (m *MockVectorClient) InsertVectors(ctx context.Context, collection string, ids []string, vectors [][]float32) error {
	args := m.Called(ctx, collection, ids, vectors)
	return args.Error(0)
}

func (m *MockVectorClient) SearchVectors(ctx context.Context, collection string, vector []float32, topK int) (*vectordb.SearchResult, error) {
	args := m.Called(ctx, collection, vector, topK)
	return args.Get(0).(*vectordb.SearchResult), args.Error(1)
}

func (m *MockVectorClient) DeleteVectors(ctx context.Context, collection string, ids []string) error {
	args := m.Called(ctx, collection, ids)
	return args.Error(0)
}

func (m *MockVectorClient) CreateCollection(ctx context.Context, collection string, dimension int) error {
	args := m.Called(ctx, collection, dimension)
	return args.Error(0)
}

func (m *MockVectorClient) DropCollection(ctx context.Context, collection string) error {
	args := m.Called(ctx, collection)
	return args.Error(0)
}

func (m *MockVectorClient) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestConversationMemoryService_StoreMemory 测试存储记忆
func TestConversationMemoryService_StoreMemory(t *testing.T) {
	// 1. 准备测试数据
	ctx := context.Background()
	db := setupTestDB(t)
	llmClient := llm.NewMockMemoryLLMClient()
	vectorClient := new(MockVectorClient)

	// 2. 设置Mock期望
	vectorClient.On("InsertVectors", ctx, "conversation_memories", mock.Anything, mock.Anything).Return(nil)

	// 3. 创建服务
	service := NewConversationMemoryService(db, vectorClient, llmClient)

	// 4. 执行测试
	memory := &convEntity.ConversationMemory{
		TenantID:       "tenant-001",
		UserID:         "user-001",
		ConversationID: "conv-001",
		MemoryType:     convEntity.MemoryTypeEntity,
		Content:        "用户是HR经理",
		ImportanceScore: 0.7,
	}

	err := service.StoreMemory(ctx, memory)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.NotEmpty(t, memory.MemoryID)
	assert.NotEmpty(t, memory.Embedding)
	assert.False(t, memory.CreatedAt.IsZero())

	// 6. 验证Mock调用
	vectorClient.AssertExpectations(t)
}

// TestConversationMemoryService_RetrieveMemories 测试检索记忆
func TestConversationMemoryService_RetrieveMemories(t *testing.T) {
	// 1. 准备测试数据
	ctx := context.Background()
	db := setupTestDB(t)
	llmClient := llm.NewMockMemoryLLMClient()
	vectorClient := new(MockVectorClient)

	service := NewConversationMemoryService(db, vectorClient, llmClient)

	// 2. 先存储记忆
	memory := &convEntity.ConversationMemory{
		TenantID:       "tenant-001",
		UserID:         "user-001",
		ConversationID: "conv-001",
		MemoryType:     convEntity.MemoryTypeEntity,
		Content:        "用户是HR经理，负责招聘",
		ImportanceScore: 0.7,
	}
	err := service.StoreMemory(ctx, memory)
	assert.NoError(t, err)

	// 3. 设置向量搜索Mock
	searchResult := &vectordb.SearchResult{
		IDs:    []string{memory.MemoryID},
		Scores: []float64{0.95},
	}
	vectorClient.On("SearchVectors", ctx, "conversation_memories", mock.Anything, 5).Return(searchResult, nil)

	// 4. 检索记忆
	memories, err := service.RetrieveMemories(ctx, "user-001", "用户的工作", 5)

	// 5. 验证结果
	assert.NoError(t, err)
	assert.Greater(t, len(memories), 0)
	if len(memories) > 0 {
		assert.Equal(t, memory.MemoryID, memories[0].Memory.MemoryID)
		assert.Greater(t, memories[0].Score, 0.7)
	}
}

// TestConversationMemoryService_DeleteExpiredMemories 测试删除过期记忆
func TestConversationMemoryService_DeleteExpiredMemories(t *testing.T) {
	// 1. 准备测试数据
	ctx := context.Background()
	db := setupTestDB(t)
	llmClient := llm.NewMockMemoryLLMClient()
	vectorClient := new(MockVectorClient)

	service := NewConversationMemoryService(db, vectorClient, llmClient)

	// 2. 创建过期记忆
	expiredTime := time.Now().Add(-1 * time.Hour)
	memory := &convEntity.ConversationMemory{
		TenantID:       "tenant-001",
		UserID:         "user-001",
		ConversationID: "conv-001",
		MemoryType:     convEntity.MemoryTypeEvent,
		Content:        "过期的事件",
		ImportanceScore: 0.3,
		ExpiresAt:      &expiredTime,
	}
	err := service.StoreMemory(ctx, memory)
	assert.NoError(t, err)

	// 3. 删除过期记忆
	err = service.DeleteExpiredMemories(ctx, "tenant-001")
	assert.NoError(t, err)
}

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	// TODO: 使用testcontainers创建MySQL实例
	// 这里返回nil，实际测试时需要实现
	return nil
}
