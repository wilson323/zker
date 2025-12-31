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

package entity

import (
	"encoding/json"
	"time"
)

// Message 消息结构
type Message struct {
	Role      string                 `json:"role"`      // user, assistant, system
	Content   string                 `json:"content"`
	TokenCount int                   `json:"token_count"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// ShortTermMemoryContext 短期记忆上下文
type ShortTermMemoryContext struct {
	ConversationID string     `json:"conversation_id"`
	Messages       []*Message `json:"messages"`
	TotalTokens    int        `json:"total_tokens"`
	MaxTokens      int        `json:"max_tokens"`
	TTL            int        `json:"ttl"`              // 过期时间(秒)
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CompressionStrategy 压缩策略
type CompressionStrategy string

const (
	CompressionStrategyRecent      CompressionStrategy = "recent"       // 保留最近消息
	CompressionStrategySummary     CompressionStrategy = "summary"      // 摘要压缩
	CompressionStrategySemantic    CompressionStrategy = "semantic"     // 语义重要性
	CompressionStrategyHybrid      CompressionStrategy = "hybrid"       // 混合策略
)

// CompressionRequest 压缩请求
type CompressionRequest struct {
	ConversationID    string             `json:"conversation_id"`
	Messages          []*Message         `json:"messages"`
	MaxTokens         int                `json:"max_tokens"`
	Strategy          CompressionStrategy `json:"strategy"`
	TargetSummaryLength int              `json:"target_summary_length"`
}

// CompressionResult 压缩结果
type CompressionResult struct {
	OriginalCount    int                `json:"original_count"`
	CompressedCount  int                `json:"compressed_count"`
	OriginalTokens   int                `json:"original_tokens"`
	CompressedTokens int                `json:"compressed_tokens"`
	CompressionRatio float64            `json:"compression_ratio"`
	Messages         []*Message         `json:"messages"`
	Summary          string             `json:"summary,omitempty"`
	StrategyUsed     CompressionStrategy `json:"strategy_used"`
	CompressedAt     time.Time          `json:"compressed_at"`
}

// StoreShortTermMemoryRequest 存储短期记忆请求
type StoreShortTermMemoryRequest struct {
	ConversationID string     `json:"conversation_id"`
	Messages       []*Message `json:"messages"`
	TTL            int        `json:"ttl"` // 过期时间(秒), 0表示不过期
	MaxTokens      int        `json:"max_tokens"` // 最大Token数
}

// GetShortTermMemoryResponse 获取短期记忆响应
type GetShortTermMemoryResponse struct {
	ConversationID string     `json:"conversation_id"`
	Messages       []*Message `json:"messages"`
	TotalTokens    int        `json:"total_tokens"`
	MessageCount   int        `json:"message_count"`
	CreatedAt      time.Time  `json:"created_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

// GetMessagesJSON 获取消息JSON
func (ctx *ShortTermMemoryContext) GetMessagesJSON() (string, error) {
	data, err := json.Marshal(ctx.Messages)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetMessagesFromJSON 从JSON设置消息
func (ctx *ShortTermMemoryContext) SetMessagesFromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), &ctx.Messages)
}

// AddMessage 添加消息
func (ctx *ShortTermMemoryContext) AddMessage(msg *Message) {
	msg.Timestamp = time.Now()
	ctx.Messages = append(ctx.Messages, msg)
	ctx.TotalTokens += msg.TokenCount
}

// CalculateTokens 计算Token数量(简化版本)
func (msg *Message) CalculateTokens() int {
	// 简化计算: 中文按1.5字符/token, 英文按4字符/token
	// 实际应使用tokenizer
	if msg.TokenCount > 0 {
		return msg.TokenCount
	}

	// 简单估算
	runes := []rune(msg.Content)
	charCount := len(runes)
	estimatedTokens := int(float64(charCount) / 2.0) // 粗略估算

	// 最少1个token
	if estimatedTokens < 1 {
		estimatedTokens = 1
	}

	return estimatedTokens
}

// ShouldCompress 判断是否需要压缩
func (ctx *ShortTermMemoryContext) ShouldCompress() bool {
	return ctx.TotalTokens > ctx.MaxTokens
}

// IsExpired 判断是否过期
func (ctx *ShortTermMemoryContext) IsExpired() bool {
	if ctx.TTL == 0 {
		return false
	}
	expireTime := ctx.CreatedAt.Add(time.Duration(ctx.TTL) * time.Second)
	return time.Now().After(expireTime)
}
