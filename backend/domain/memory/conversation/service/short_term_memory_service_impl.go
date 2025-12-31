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

	"github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
	"github.com/coze-dev/coze-studio/backend/pkg/redis"
)

const (
	// DefaultMaxTokens 默认最大Token数
	DefaultMaxTokens = 8000
	// ShortTermMemoryPrefix Redis key前缀
	ShortTermMemoryPrefix = "short_term_memory:"
	// DefaultTTL 默认过期时间(秒)
	DefaultTTL = 3600 // 1小时
)

// NewShortTermMemoryService 创建短期记忆服务
func NewShortTermMemoryService(
	redisClient *redis.Client,
	llmClient llm.MemoryLLMClient,
) ShortTermMemoryService {
	return &shortTermMemoryServiceImpl{
		redisClient: redisClient,
		llmClient:   llmClient,
	}
}

type shortTermMemoryServiceImpl struct {
	redisClient *redis.Client
	llmClient   llm.MemoryLLMClient
}

// Store 存储短期记忆上下文
func (s *shortTermMemoryServiceImpl) Store(
	ctx context.Context,
	req *entity.StoreShortTermMemoryRequest,
) error {
	// 1. 计算Token总数
	totalTokens := 0
	for _, msg := range req.Messages {
		if msg.TokenCount == 0 {
			msg.TokenCount = msg.CalculateTokens()
		}
		totalTokens += msg.TokenCount
	}

	// 2. 设置默认值
	if req.MaxTokens == 0 {
		req.MaxTokens = DefaultMaxTokens
	}
	if req.TTL == 0 {
		req.TTL = DefaultTTL
	}

	// 3. 检查是否超过窗口大小
	messages := req.Messages
	originalTokens := totalTokens

	if totalTokens > req.MaxTokens {
		// 自动压缩
		compressReq := &entity.CompressionRequest{
			ConversationID: req.ConversationID,
			Messages:       req.Messages,
			MaxTokens:      req.MaxTokens,
			Strategy:       entity.CompressionStrategyHybrid,
		}
		result, err := s.Compress(ctx, compressReq)
		if err != nil {
			return fmt.Errorf("failed to compress context: %w", err)
		}
		messages = result.Messages
		totalTokens = result.CompressedTokens
	}

	// 4. 构建上下文对象
	now := time.Now()
	context := &entity.ShortTermMemoryContext{
		ConversationID: req.ConversationID,
		Messages:       messages,
		TotalTokens:    totalTokens,
		MaxTokens:      req.MaxTokens,
		TTL:            req.TTL,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 5. 存储到Redis
	key := s.getRedisKey(req.ConversationID)
	err := s.redisClient.Set(ctx, key, context, time.Duration(req.TTL)*time.Second)
	if err != nil {
		return fmt.Errorf("failed to store in redis: %w", err)
	}

	return nil
}

// Get 获取短期记忆上下文
func (s *shortTermMemoryServiceImpl) Get(
	ctx context.Context,
	conversationID string,
) (*entity.GetShortTermMemoryResponse, error) {
	// 1. 从Redis获取
	key := s.getRedisKey(conversationID)
	var context entity.ShortTermMemoryContext
	err := s.redisClient.Get(ctx, key, &context)
	if err != nil {
		return nil, fmt.Errorf("failed to get from redis: %w", err)
	}

	// 2. 构建响应
	var expiresAt *time.Time
	if context.TTL > 0 {
		expireTime := context.UpdatedAt.Add(time.Duration(context.TTL) * time.Second)
		expiresAt = &expireTime
	}

	return &entity.GetShortTermMemoryResponse{
		ConversationID: context.ConversationID,
		Messages:       context.Messages,
		TotalTokens:    context.TotalTokens,
		MessageCount:   len(context.Messages),
		CreatedAt:      context.CreatedAt,
		ExpiresAt:      expiresAt,
	}, nil
}

// Delete 删除短期记忆
func (s *shortTermMemoryServiceImpl) Delete(
	ctx context.Context,
	conversationID string,
) error {
	key := s.getRedisKey(conversationID)
	err := s.redisClient.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete from redis: %w", err)
	}
	return nil
}

// Compress 压缩上下文
func (s *shortTermMemoryServiceImpl) Compress(
	ctx context.Context,
	req *entity.CompressionRequest,
) (*entity.CompressionResult, error) {
	// 计算原始Token数
	originalTokens := 0
	for _, msg := range req.Messages {
		if msg.TokenCount == 0 {
			msg.TokenCount = msg.CalculateTokens()
		}
		originalTokens += msg.TokenCount
	}

	var compressedMessages []*entity.Message
	var summary string
	var err error

	// 根据策略选择压缩算法
	switch req.Strategy {
	case entity.CompressionStrategyRecent:
		compressedMessages = s.compressByRecent(req.Messages, req.MaxTokens)
	case entity.CompressionStrategySummary:
		compressedMessages, summary, err = s.compressBySummary(ctx, req.Messages, req.MaxTokens)
		if err != nil {
			return nil, fmt.Errorf("failed to compress by summary: %w", err)
		}
	case entity.CompressionStrategySemantic:
		compressedMessages = s.compressBySemantic(req.Messages, req.MaxTokens)
	case entity.CompressionStrategyHybrid:
		compressedMessages, summary, err = s.compressHybrid(ctx, req.Messages, req.MaxTokens)
		if err != nil {
			return nil, fmt.Errorf("failed to compress by hybrid: %w", err)
		}
	default:
		compressedMessages = s.compressByRecent(req.Messages, req.MaxTokens)
	}

	// 计算压缩后Token数
	compressedTokens := 0
	for _, msg := range compressedMessages {
		compressedTokens += msg.TokenCount
	}

	// 计算压缩比
	compressionRatio := 0.0
	if originalTokens > 0 {
		compressionRatio = 1.0 - float64(compressedTokens)/float64(originalTokens)
	}

	return &entity.CompressionResult{
		OriginalCount:    len(req.Messages),
		CompressedCount:  len(compressedMessages),
		OriginalTokens:   originalTokens,
		CompressedTokens: compressedTokens,
		CompressionRatio: compressionRatio,
		Messages:         compressedMessages,
		Summary:          summary,
		StrategyUsed:     req.Strategy,
		CompressedAt:     time.Now(),
	}, nil
}

// AutoCompress 自动压缩
func (s *shortTermMemoryServiceImpl) AutoCompress(
	ctx context.Context,
	conversationID string,
) (*entity.CompressionResult, error) {
	// 1. 获取当前上下文
	resp, err := s.Get(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get context: %w", err)
	}

	// 2. 检查是否需要压缩
	if resp.TotalTokens <= resp.TotalTokens {
		return nil, nil // 不需要压缩
	}

	// 3. 执行压缩
	req := &entity.CompressionRequest{
		ConversationID: conversationID,
		Messages:       resp.Messages,
		MaxTokens:      resp.TotalTokens,
		Strategy:       entity.CompressionStrategyHybrid,
	}

	result, err := s.Compress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to compress: %w", err)
	}

	// 4. 存储压缩后的上下文
	storeReq := &entity.StoreShortTermMemoryRequest{
		ConversationID: conversationID,
		Messages:       result.Messages,
		MaxTokens:      resp.TotalTokens,
		TTL:            DefaultTTL,
	}

	err = s.Store(ctx, storeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to store compressed context: %w", err)
	}

	return result, nil
}

// AddMessage 添加消息到上下文
func (s *shortTermMemoryServiceImpl) AddMessage(
	ctx context.Context,
	conversationID string,
	msg *entity.Message,
) error {
	// 1. 获取当前上下文
	resp, err := s.Get(ctx, conversationID)
	if err != nil {
		// 如果不存在,创建新上下文
		if err.Error() == "redis: key not found" {
			req := &entity.StoreShortTermMemoryRequest{
				ConversationID: conversationID,
				Messages:       []*entity.Message{msg},
				TTL:            DefaultTTL,
				MaxTokens:      DefaultMaxTokens,
			}
			return s.Store(ctx, req)
		}
		return fmt.Errorf("failed to get context: %w", err)
	}

	// 2. 计算Token
	if msg.TokenCount == 0 {
		msg.TokenCount = msg.CalculateTokens()
	}

	// 3. 添加消息
	messages := append(resp.Messages, msg)

	// 4. 存储更新后的上下文
	req := &entity.StoreShortTermMemoryRequest{
		ConversationID: conversationID,
		Messages:       messages,
		TTL:            DefaultTTL,
		MaxTokens:      resp.TotalTokens, // 保持原来的maxTokens
	}

	return s.Store(ctx, req)
}

// ClearExpired 清理过期的短期记忆
func (s *shortTermMemoryServiceImpl) ClearExpired(ctx context.Context) error {
	// Redis会自动清理过期的key
	// 这里可以添加额外的逻辑,比如记录日志
	return nil
}

// GetWindowStats 获取窗口统计信息
func (s *shortTermMemoryServiceImpl) GetWindowStats(
	ctx context.Context,
	conversationID string,
) (*WindowStats, error) {
	resp, err := s.Get(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get context: %w", err)
	}

	utilizationRatio := 0.0
	if resp.TotalTokens > 0 {
		utilizationRatio = float64(resp.TotalTokens) / float64(DefaultMaxTokens)
	}

	var expiresAt *int64
	if resp.ExpiresAt != nil {
		timestamp := resp.ExpiresAt.Unix()
		expiresAt = &timestamp
	}

	return &WindowStats{
		ConversationID:   conversationID,
		MessageCount:     resp.MessageCount,
		TotalTokens:      resp.TotalTokens,
		MaxTokens:        DefaultMaxTokens,
		UtilizationRatio: utilizationRatio,
		NeedsCompression: resp.TotalTokens > DefaultMaxTokens,
		ExpiresAt:        expiresAt,
	}, nil
}

// compressByRecent 按最近消息压缩
func (s *shortTermMemoryServiceImpl) compressByRecent(
	messages []*entity.Message,
	maxTokens int,
) []*entity.Message {
	result := make([]*entity.Message, 0, len(messages))
	currentTokens := 0

	// 从后往前保留消息
	for i := len(messages) - 1; i >= 0; i-- {
		msgTokens := messages[i].TokenCount
		if currentTokens+msgTokens > maxTokens {
			break
		}
		result = append([]*entity.Message{messages[i]}, result...)
		currentTokens += msgTokens
	}

	return result
}

// compressBySummary 按摘要压缩
func (s *shortTermMemoryServiceImpl) compressBySummary(
	ctx context.Context,
	messages []*entity.Message,
	maxTokens int,
) ([]*entity.Message, string, error) {
	// 1. 构建消息文本
	text := ""
	for _, msg := range messages {
		text += fmt.Sprintf("[%s]: %s\n", msg.Role, msg.Content)
	}

	// 2. 调用LLM生成摘要
	prompt := fmt.Sprintf(
		"请将以下对话内容压缩为简洁的摘要(保留关键信息):\n%s\n\n请只返回摘要内容,不要有任何额外说明。",
		text,
	)

	summary, err := s.llmClient.Chat(ctx, &llm.ChatRequest{
		Messages: []*llm.Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 500,
	})
	if err != nil {
		return nil, "", err
	}

	// 3. 构建摘要消息
	summaryMsg := &entity.Message{
		Role:      "system",
		Content:   fmt.Sprintf("[对话摘要]: %s", summary),
		TokenCount: len(summary) / 2, // 估算
		Timestamp: time.Now(),
	}

	// 4. 保留最近的消息(预留空间给摘要)
	remainingTokens := maxTokens - summaryMsg.TokenCount
	recentMessages := s.compressByRecent(messages, remainingTokens)

	// 5. 组合结果
	result := append([]*entity.Message{summaryMsg}, recentMessages...)

	return result, summary, nil
}

// compressBySemantic 按语义重要性压缩
func (s *shortTermMemoryServiceImpl) compressBySemantic(
	messages []*entity.Message,
	maxTokens int,
) []*entity.Message {
	// 简化实现: 优先保留system和user消息
	result := make([]*entity.Message, 0, len(messages))
	currentTokens := 0

	// 优先级: system > user > assistant
	priority := map[string]int{
		"system":     3,
		"user":       2,
		"assistant":  1,
	}

	// 按优先级排序
	sortedMessages := make([]*entity.Message, len(messages))
	copy(sortedMessages, messages)

	// 简单冒泡排序(按优先级)
	for i := 0; i < len(sortedMessages)-1; i++ {
		for j := i + 1; j < len(sortedMessages); j++ {
			if priority[sortedMessages[i].Role] < priority[sortedMessages[j].Role] {
				sortedMessages[i], sortedMessages[j] = sortedMessages[j], sortedMessages[i]
			}
		}
	}

	// 选择高优先级消息
	for _, msg := range sortedMessages {
		if currentTokens+msg.TokenCount > maxTokens {
			break
		}
		result = append(result, msg)
		currentTokens += msg.TokenCount
	}

	return result
}

// compressHybrid 混合压缩策略
func (s *shortTermMemoryServiceImpl) compressHybrid(
	ctx context.Context,
	messages []*entity.Message,
	maxTokens int,
) ([]*entity.Message, string, error) {
	// 如果消息数量不多,直接使用recent策略
	if len(messages) <= 10 {
		return s.compressByRecent(messages, maxTokens), "", nil
	}

	// 否则使用summary+recent策略
	return s.compressBySummary(ctx, messages, maxTokens)
}

// getRedisKey 获取Redis key
func (s *shortTermMemoryServiceImpl) getRedisKey(conversationID string) string {
	return ShortTermMemoryPrefix + conversationID
}
