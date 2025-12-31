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

package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// DeadLetterQueue 死信队列
type DeadLetterQueue struct {
	producer *NSQProducer
	config   *NSQConfig
	logger   *zap.Logger
	storage  DLQStorage // 可选的持久化存储

	// 指标
	metrics *DLQMetrics
	mu      sync.RWMutex
}

// DLQMetrics 死信队列指标
type DLQMetrics struct {
	MessagesStored   int64
	MessagesReplayed int64
	MessagesDeleted  int64
	MessagesExpired  int64
}

// DLQMessage 死信队列消息
type DLQMessage struct {
	// 原始消息信息
	OriginalTopic string    `json:"original_topic"`
	OriginalBody  []byte    `json:"original_body"`
	FailedAt      time.Time `json:"failed_at"`

	// 失败原因
	LastError     string    `json:"last_error"`
	AttemptCount  int       `json:"attempt_count"`

	// 元数据
	MessageID     string    `json:"message_id"`
	ReplayCount   int       `json:"replay_count"`
	NextReplay    time.Time `json:"next_replay,omitempty"`

	// 附加信息
	Headers       map[string]string `json:"headers,omitempty"`
}

// DLQStorage 死信队列存储接口（可选）
type DLQStorage interface {
	Store(ctx context.Context, topic string, message *DLQMessage) error
	Get(ctx context.Context, messageID string) (*DLQMessage, error)
	List(ctx context.Context, topic string, limit int) ([]*DLQMessage, error)
	Delete(ctx context.Context, messageID string) error
	MarkReplayed(ctx context.Context, messageID string) error
}

// NewDeadLetterQueue 创建死信队列
func NewDeadLetterQueue(producer *NSQProducer, config *NSQConfig, logger *zap.Logger) *DeadLetterQueue {
	if config == nil {
		config = DefaultNSQConfig()
	}

	dlq := &DeadLetterQueue{
		producer: producer,
		config:   config,
		logger:   logger,
		metrics:  &DLQMetrics{},
	}

	return dlq
}

// SetStorage 设置存储（可选）
func (dlq *DeadLetterQueue) SetStorage(storage DLQStorage) {
	dlq.mu.Lock()
	defer dlq.mu.Unlock()
	dlq.storage = storage
}

// MoveToDLQ 将失败消息移动到死信队列
func (dlq *DeadLetterQueue) MoveToDLQ(ctx context.Context, originalTopic string, messageBody []byte, err error) error {
	now := time.Now()

	// 创建死信消息
	dlqMsg := &DLQMessage{
		OriginalTopic: originalTopic,
		OriginalBody:  messageBody,
		FailedAt:      now,
		LastError:     err.Error(),
		AttemptCount:  0,
		MessageID:     generateMessageID(),
		ReplayCount:   0,
	}

	// 序列化死信消息
	data, err := json.Marshal(dlqMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal DLQ message: %w", err)
	}

	// 发布到死信队列Topic
	dlqTopic := originalTopic + dlq.config.DLQTopicSuffix
	if publishErr := dlq.producer.Publish(ctx, dlqTopic, data); publishErr != nil {
		return fmt.Errorf("failed to publish to DLQ: %w", publishErr)
	}

	// 可选：持久化到存储
	if dlq.storage != nil {
		if storeErr := dlq.storage.Store(ctx, dlqTopic, dlqMsg); storeErr != nil {
			dlq.logger.Error("failed to persist DLQ message to storage",
				zap.String("message_id", dlqMsg.MessageID),
				zap.Error(storeErr),
			)
			// 不返回错误，因为消息已经发送到NSQ
		}
	}

	// 更新指标
	dlq.mu.Lock()
	dlq.metrics.MessagesStored++
	dlq.mu.Unlock()

	dlq.logger.Warn("message moved to DLQ",
		zap.String("original_topic", originalTopic),
		zap.String("dlq_topic", dlqTopic),
		zap.String("message_id", dlqMsg.MessageID),
		zap.Error(err),
	)

	return nil
}

// ReplayDLQ 重放死信队列中的消息
func (dlq *DeadLetterQueue) ReplayDLQ(ctx context.Context, topic string) error {
	dlqTopic := topic + dlq.config.DLQTopicSuffix

	// 如果有持久化存储，从存储中读取
	if dlq.storage != nil {
		return dlq.replayFromStorage(ctx, topic, dlqTopic)
	}

	// 否则，直接从NSQ重放（需要消费者处理）
	dlq.logger.Info("DLQ replay requested (from NSQ)",
		zap.String("topic", topic),
		zap.String("dlq_topic", dlqTopic),
	)

	return fmt.Errorf("DLQ replay requires DLQStorage to be configured")
}

// replayFromStorage 从持久化存储重放消息
func (dlq *DeadLetterQueue) replayFromStorage(ctx context.Context, topic, dlqTopic string) error {
	// 获取DLQ消息列表
	messages, err := dlq.storage.List(ctx, topic, 100) // 每次最多100条
	if err != nil {
		return fmt.Errorf("failed to list DLQ messages: %w", err)
	}

	replayedCount := 0
	failedCount := 0

	for _, msg := range messages {
		// 检查重放次数限制
		if msg.ReplayCount >= dlq.config.DLQMaxReplays {
			dlq.logger.Warn("DLQ message exceeded max replays, skipping",
				zap.String("message_id", msg.MessageID),
				zap.Int("replay_count", msg.ReplayCount),
				zap.Int("max_replays", dlq.config.DLQMaxReplays),
			)
			failedCount++
			continue
		}

		// 重新发布到原始Topic
		if err := dlq.producer.Publish(ctx, topic, msg.OriginalBody); err != nil {
			dlq.logger.Error("failed to replay DLQ message",
				zap.String("message_id", msg.MessageID),
				zap.Error(err),
			)
			failedCount++
			continue
		}

		// 标记为已重放
		if err := dlq.storage.MarkReplayed(ctx, msg.MessageID); err != nil {
			dlq.logger.Error("failed to mark message as replayed",
				zap.String("message_id", msg.MessageID),
				zap.Error(err),
			)
		}

		replayedCount++
	}

	// 更新指标
	dlq.mu.Lock()
	dlq.metrics.MessagesReplayed += int64(replayedCount)
	dlq.mu.Unlock()

	dlq.logger.Info("DLQ replay completed",
		zap.String("topic", topic),
		zap.Int("replayed", replayedCount),
		zap.Int("failed", failedCount),
	)

	return nil
}

// ReplayMessage 重放单条消息
func (dlq *DeadLetterQueue) ReplayMessage(ctx context.Context, messageID string) error {
	if dlq.storage == nil {
		return fmt.Errorf("DLQStorage not configured")
	}

	// 获取消息
	msg, err := dlq.storage.Get(ctx, messageID)
	if err != nil {
		return fmt.Errorf("failed to get DLQ message: %w", err)
	}

	// 检查重放次数限制
	if msg.ReplayCount >= dlq.config.DLQMaxReplays {
		return fmt.Errorf("message exceeded max replays (%d)", dlq.config.DLQMaxReplays)
	}

	// 重新发布到原始Topic
	if err := dlq.producer.Publish(ctx, msg.OriginalTopic, msg.OriginalBody); err != nil {
		return fmt.Errorf("failed to replay message: %w", err)
	}

	// 标记为已重放
	if err := dlq.storage.MarkReplayed(ctx, messageID); err != nil {
		dlq.logger.Error("failed to mark message as replayed",
			zap.String("message_id", messageID),
			zap.Error(err),
		)
	}

	// 更新指标
	dlq.mu.Lock()
	dlq.metrics.MessagesReplayed++
	dlq.mu.Unlock()

	dlq.logger.Info("DLQ message replayed",
		zap.String("message_id", messageID),
		zap.String("original_topic", msg.OriginalTopic),
	)

	return nil
}

// DeleteMessage 从死信队列删除消息
func (dlq *DeadLetterQueue) DeleteMessage(ctx context.Context, messageID string) error {
	if dlq.storage == nil {
		return fmt.Errorf("DLQStorage not configured")
	}

	if err := dlq.storage.Delete(ctx, messageID); err != nil {
		return fmt.Errorf("failed to delete DLQ message: %w", err)
	}

	// 更新指标
	dlq.mu.Lock()
	dlq.metrics.MessagesDeleted++
	dlq.mu.Unlock()

	dlq.logger.Info("DLQ message deleted",
		zap.String("message_id", messageID),
	)

	return nil
}

// GetMetrics 获取死信队列指标
func (dlq *DeadLetterQueue) GetMetrics() *DLQMetrics {
	dlq.mu.RLock()
	defer dlq.mu.RUnlock()

	return &DLQMetrics{
		MessagesStored:   dlq.metrics.MessagesStored,
		MessagesReplayed: dlq.metrics.MessagesReplayed,
		MessagesDeleted:  dlq.metrics.MessagesDeleted,
		MessagesExpired:  dlq.metrics.MessagesExpired,
	}
}

// CleanupExpiredMessages 清理过期的死信消息
func (dlq *DeadLetterQueue) CleanupExpiredMessages(ctx context.Context) error {
	if dlq.storage == nil {
		return fmt.Errorf("DLQStorage not configured")
	}

	// TODO: 实现清理逻辑
	// 1. 从存储中查询所有DLQ消息
	// 2. 检查FailedAt + RetentionDuration < Now
	// 3. 删除过期消息

	return fmt.Errorf("not implemented")
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return fmt.Sprintf("dlq_%d", time.Now().UnixNano())
}

// =====================================================================
// 内存死信队列存储（用于测试或开发）
// =====================================================================

// MemoryDLQStorage 内存死信队列存储
type MemoryDLQStorage struct {
	mu       sync.RWMutex
	messages map[string]*DLQMessage
	topics   map[string][]string // topic -> message IDs
}

// NewMemoryDLQStorage 创建内存死信队列存储
func NewMemoryDLQStorage() *MemoryDLQStorage {
	return &MemoryDLQStorage{
		messages: make(map[string]*DLQMessage),
		topics:   make(map[string][]string),
	}
}

// Store 存储消息
func (s *MemoryDLQStorage) Store(ctx context.Context, topic string, message *DLQMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages[message.MessageID] = message
	s.topics[topic] = append(s.topics[topic], message.MessageID)

	return nil
}

// Get 获取消息
func (s *MemoryDLQStorage) Get(ctx context.Context, messageID string) (*DLQMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg, ok := s.messages[messageID]
	if !ok {
		return nil, fmt.Errorf("message not found: %s", messageID)
	}

	// 返回副本
	msgCopy := *msg
	return &msgCopy, nil
}

// List 列出消息
func (s *MemoryDLQStorage) List(ctx context.Context, topic string, limit int) ([]*DLQMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	messageIDs, ok := s.topics[topic]
	if !ok {
		return []*DLQMessage{}, nil
	}

	result := make([]*DLQMessage, 0, limit)
	for i := 0; i < len(messageIDs) && i < limit; i++ {
		msg, ok := s.messages[messageIDs[i]]
		if ok {
			msgCopy := *msg
			result = append(result, &msgCopy)
		}
	}

	return result, nil
}

// Delete 删除消息
func (s *MemoryDLQStorage) Delete(ctx context.Context, messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.messages, messageID)

	// 从topics中移除
	for topic, ids := range s.topics {
		for i, id := range ids {
			if id == messageID {
				s.topics[topic] = append(ids[:i], ids[i+1:]...)
				break
			}
		}
	}

	return nil
}

// MarkReplayed 标记消息已重放
func (s *MemoryDLQStorage) MarkReplayed(ctx context.Context, messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[messageID]
	if !ok {
		return fmt.Errorf("message not found: %s", messageID)
	}

	msg.ReplayCount++
	return nil
}

// GetMessageCount 获取消息数量（用于监控）
func (s *MemoryDLQStorage) GetMessageCount(topic string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, ok := s.topics[topic]
	if !ok {
		return 0
	}

	return len(ids)
}
