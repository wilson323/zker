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

	"github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
)

// MockRedisClient Mock Redis客户端
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisClient) Get(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

// MockLLMClient Mock LLM客户端
type MockLLMClient struct {
	mock.Mock
}

func (m *MockLLMClient) Embed(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), nil
}

func (m *MockLLMClient) Chat(ctx context.Context, req *llm.ChatRequest) (string, error) {
	args := m.Called(ctx, req)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return args.String(0), nil
}

// TestShortTermMemoryService_Store 测试存储短期记忆
func TestShortTermMemoryService_Store(t *testing.T) {
	// 1. 创建Mock客户端
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)

	// 2. 设置期望
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// 3. 创建服务
	service := NewShortTermMemoryService(mockRedis, mockLLM)

	// 4. 准备测试数据
	req := &entity.StoreShortTermMemoryRequest{
		ConversationID: "test-conv-001",
		Messages: []*entity.Message{
			{
				Role:      "user",
				Content:   "你好",
				TokenCount: 5,
				Timestamp: time.Now(),
			},
			{
				Role:      "assistant",
				Content:   "你好!有什么可以帮助你的吗?",
				TokenCount: 20,
				Timestamp: time.Now(),
			},
		},
		TTL:       3600,
		MaxTokens: 8000,
	}

	// 5. 执行测试
	err := service.Store(context.Background(), req)

	// 6. 验证结果
	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

// TestShortTermMemoryService_Get 测试获取短期记忆
func TestShortTermMemoryService_Get(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)

	// 设置期望: Set用于准备数据, Get用于测试
	testContext := &entity.ShortTermMemoryContext{
		ConversationID: "test-conv-001",
		Messages: []*entity.Message{
			{
				Role:        "user",
				Content:     "测试消息",
				TokenCount:  10,
				Timestamp:   time.Now(),
			},
		},
		TotalTokens: 10,
		MaxTokens:   8000,
		TTL:         3600,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		key := args.String(1)
		value := args.Get(2).(*entity.ShortTermMemoryContext)

		// 模拟Redis Get行为
		*value = *testContext
	}).Return(nil)

	service := NewShortTermMemoryService(mockRedis, mockLLM)

	resp, err := service.Get(context.Background(), "test-conv-001")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-conv-001", resp.ConversationID)
	assert.Equal(t, 1, resp.MessageCount)
	assert.Equal(t, 10, resp.TotalTokens)
	mockRedis.AssertExpectations(t)
}

// TestShortTermMemoryService_CompressByRecent 测试按最近消息压缩
func TestShortTermMemoryService_CompressByRecent(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewShortTermMemoryService(mockRedis, mockLLM)

	// 创建超过窗口大小的消息列表
	messages := []*entity.Message{}
	for i := 0; i < 20; i++ {
		messages = append(messages, &entity.Message{
			Role:        "user",
			Content:     "测试消息",
			TokenCount:  500,
			Timestamp:   time.Now(),
		})
	}

	req := &entity.CompressionRequest{
		ConversationID: "test-conv-001",
		Messages:       messages,
		MaxTokens:      5000,
		Strategy:       entity.CompressionStrategyRecent,
	}

	result, err := service.Compress(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 20, result.OriginalCount)
	assert.LessOrEqual(t, result.CompressedTokens, 5000)
	assert.Greater(t, result.CompressionRatio, 0.0)
	assert.Equal(t, entity.CompressionStrategyRecent, result.StrategyUsed)
}

// TestShortTermMemoryService_CompressBySemantic 测试按语义压缩
func TestShortTermMemoryService_CompressBySemantic(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewShortTermMemoryService(mockRedis, mockLLM)

	messages := []*entity.Message{
		{Role: "system", Content: "系统消息", TokenCount: 100, Timestamp: time.Now()},
		{Role: "user", Content: "用户消息1", TokenCount: 200, Timestamp: time.Now()},
		{Role: "assistant", Content: "助手回复1", TokenCount: 300, Timestamp: time.Now()},
		{Role: "user", Content: "用户消息2", TokenCount: 200, Timestamp: time.Now()},
		{Role: "assistant", Content: "助手回复2", TokenCount: 300, Timestamp: time.Now()},
	}

	req := &entity.CompressionRequest{
		ConversationID: "test-conv-001",
		Messages:       messages,
		MaxTokens:      500,
		Strategy:       entity.CompressionStrategySemantic,
	}

	result, err := service.Compress(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Semantic策略应该优先保留system和user消息
	assert.LessOrEqual(t, result.CompressedTokens, 500)
}

// TestShortTermMemoryService_AddMessage 测试添加消息
func TestShortTermMemoryService_AddMessage(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)

	// 第一次Get返回不存在(创建新上下文)
	mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError).Once()
	// 第一次Set创建新上下文
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	service := NewShortTermMemoryService(mockRedis, mockLLM)

	msg := &entity.Message{
		Role:      "user",
		Content:   "新消息",
		Timestamp: time.Now(),
	}

	err := service.AddMessage(context.Background(), "new-conv-001", msg)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

// TestShortTermMemoryService_Delete 测试删除短期记忆
func TestShortTermMemoryService_Delete(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)
	mockRedis.On("Del", mock.Anything, mock.Anything).Return(nil)

	service := NewShortTermMemoryService(mockRedis, mockLLM)

	err := service.Delete(context.Background(), "test-conv-001")

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

// TestShortTermMemoryService_GetWindowStats 测试获取窗口统计
func TestShortTermMemoryService_GetWindowStats(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)

	testContext := &entity.ShortTermMemoryContext{
		ConversationID: "test-conv-001",
		Messages: []*entity.Message{
			{Role: "user", Content: "测试", TokenCount: 10, Timestamp: time.Now()},
		},
		TotalTokens: 10,
		MaxTokens:   8000,
		TTL:         3600,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		value := args.Get(2).(*entity.ShortTermMemoryContext)
		*value = *testContext
	}).Return(nil)

	service := NewShortTermMemoryService(mockRedis, mockLLM)

	stats, err := service.GetWindowStats(context.Background(), "test-conv-001")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, "test-conv-001", stats.ConversationID)
	assert.Equal(t, 1, stats.MessageCount)
	assert.Equal(t, 10, stats.TotalTokens)
	assert.Equal(t, 8000, stats.MaxTokens)
	assert.False(t, stats.NeedsCompression)
	assert.InDelta(t, 0.00125, stats.UtilizationRatio, 0.0001)
}

// BenchmarkCompressByRecent 性能测试
func BenchmarkCompressByRecent(b *testing.B) {
	mockRedis := new(MockRedisClient)
	mockLLM := new(MockLLMClient)
	service := NewShortTermMemoryService(mockRedis, mockLLM)

	// 创建100条消息
	messages := make([]*entity.Message, 100)
	for i := 0; i < 100; i++ {
		messages[i] = &entity.Message{
			Role:        "user",
			Content:     "测试消息内容",
			TokenCount:  100,
			Timestamp:   time.Now(),
		}
	}

	req := &entity.CompressionRequest{
		Messages:  messages,
		MaxTokens: 5000,
		Strategy:  entity.CompressionStrategyRecent,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Compress(context.Background(), req)
	}
}
