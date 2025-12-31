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
	"sync/atomic"
	"time"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// NSQConsumer NSQ消费者（增强版）
type NSQConsumer struct {
	consumer *nsq.Consumer
	config   *NSQConfig
	logger   *zap.Logger
	handler  MessageHandler
	dlq      *DeadLetterQueue

	// 控制状态
	stopCh    chan struct{}
	stopOnce  sync.Once
	running   atomic.Bool

	// 指标
	metrics *ConsumerMetrics

	// 重试追踪
	retryTracker map[string]*RetryInfo
	retryMu      sync.RWMutex
	// Topic信息（用于死信队列）
	topic string
}

// RetryInfo 重试信息
type RetryInfo struct {
	AttemptCount int
	FirstAttempt time.Time
	LastError    string
	NextRetry    time.Time
}

// ConsumerMetrics 消费者指标
type ConsumerMetrics struct {
	mu                sync.RWMutex
	MessagesReceived  int64
	MessagesSuccess   int64
	MessagesFailed    int64
	MessagesRetried   int64
	MessagesToDLQ     int64
	ProcessingTime    []time.Duration
}

// MessageHandler 消息处理接口
type MessageHandler interface {
	HandleMessage(ctx context.Context, message *nsq.Message) error
}

// MessageHandlerFunc 函数式消息处理器
type MessageHandlerFunc func(ctx context.Context, message *nsq.Message) error

// HandleMessage 实现MessageHandler接口
func (f MessageHandlerFunc) HandleMessage(ctx context.Context, message *nsq.Message) error {
	return f(ctx, message)
}

// NewNSQConsumer 创建NSQ消费者
func NewNSQConsumer(
	topic string,
	channel string,
	config *NSQConfig,
	handler MessageHandler,
	logger *zap.Logger,
	dlq *DeadLetterQueue,
) (*NSQConsumer, error) {
	if config == nil {
		config = DefaultNSQConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid NSQ config: %w", err)
	}

	if handler == nil {
		return nil, fmt.Errorf("handler cannot be nil")
	}

	// 配置NSQ消费者
	nsqConfig := nsq.NewConfig()
	nsqConfig.MaxAttempts = uint16(config.ConsumerMaxRetries + 1) // +1 因为第一次尝试不算重试
	nsqConfig.MaxInFlight = config.ConsumerMaxInFlight
	nsqConfig.MsgTimeout = config.ConsumerMessageTimeout
	nsqConfig.HeartbeatInterval = config.ConsumerHeartbeatInterval

	// 创建消费者
	consumer, err := nsq.NewConsumer(topic, channel, nsqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ consumer: %w", err)
	}

	// 设置日志
	consumer.SetLogger(newNSQLogger(logger, LogLevelInfo), toNSQLogLevel(LogLevelInfo))

	c := &NSQConsumer{
		consumer:     consumer,
		config:       config,
		logger:       logger,
		handler:      handler,
		dlq:          dlq,
		stopCh:       make(chan struct{}),
		topic:        topic,
		metrics:      &ConsumerMetrics{},
		retryTracker: make(map[string]*RetryInfo),
	}

	// 注册处理器
	consumer.AddHandler(nsq.HandlerFunc(c.handleMessage))

	c.logger.Info("NSQ consumer initialized",
		zap.String("topic", topic),
		zap.String("channel", channel),
		zap.Int("max_retries", config.ConsumerMaxRetries),
		zap.Int("max_in_flight", config.ConsumerMaxInFlight),
	)

	return c, nil
}

// Consume 开始消费消息（连接到NSQLookupd）
func (c *NSQConsumer) Consume() error {
	if !c.running.CompareAndSwap(false, true) {
		return fmt.Errorf("consumer is already running")
	}

	// 连接到所有NSQLookupd地址
	for _, lookupdAddr := range c.config.NSQLookupdAddresses {
		err := c.consumer.ConnectToNSQLookupd(lookupdAddr)
		if err != nil {
			c.running.Store(false)
			return fmt.Errorf("failed to connect to NSQLookupd %s: %w", lookupdAddr, err)
		}
	}

	c.logger.Info("NSQ consumer started",
		zap.Strings("nsqlookupd_addresses", c.config.NSQLookupdAddresses),
	)

	return nil
}

// ConsumeFromNSQD 直接连接到NSQD（用于测试或单节点）
func (c *NSQConsumer) ConsumeFromNSQD(nsqdAddr string) error {
	if !c.running.CompareAndSwap(false, true) {
		return fmt.Errorf("consumer is already running")
	}

	err := c.consumer.ConnectToNSQD(nsqdAddr)
	if err != nil {
		c.running.Store(false)
		return fmt.Errorf("failed to connect to NSQD %s: %w", nsqdAddr, err)
	}

	c.logger.Info("NSQ consumer connected to NSQD",
		zap.String("nsqd_address", nsqdAddr),
	)

	return nil
}

// handleMessage 处理单条消息
func (c *NSQConsumer) handleMessage(message *nsq.Message) error {
	startTime := time.Now()
	ctx := context.Background()

	// 增加接收计数
	atomic.AddInt64(&c.metrics.MessagesReceived, 1)

	// 记录处理时间
	defer func() {
		processingTime := time.Since(startTime)
		c.metrics.mu.Lock()
		if len(c.metrics.ProcessingTime) >= 1000 {
			c.metrics.ProcessingTime = c.metrics.ProcessingTime[1:]
		}
		c.metrics.ProcessingTime = append(c.metrics.ProcessingTime, processingTime)
		c.metrics.mu.Unlock()
	}()

	// 检查重试次数
	msgID := string(message.ID[:])
	attempt := c.getAttemptCount(msgID)

	if attempt > c.config.ConsumerMaxRetries {
		c.logger.Error("message exceeded max retries, sending to DLQ",
			zap.String("message_id", msgID),
			zap.Int("attempt", attempt),
			zap.Int("max_retries", c.config.ConsumerMaxRetries),
		)

		// 发送到死信队列
		if c.dlq != nil && c.config.DLQEnabled {
			if err := c.dlq.MoveToDLQ(ctx, c.topic, message.Body, fmt.Errorf("max retries exceeded")); err != nil {
				c.logger.Error("failed to send message to DLQ",
					zap.String("message_id", msgID),
					zap.Error(err),
				)
			} else {
				atomic.AddInt64(&c.metrics.MessagesToDLQ, 1)
			}
		}

		c.clearRetryInfo(msgID)
		atomic.AddInt64(&c.metrics.MessagesFailed, 1)
		return nil // 返回nil以避免NSQ重试
	}

	// 处理消息
	err := c.handler.HandleMessage(ctx, message)
	if err != nil {
		atomic.AddInt64(&c.metrics.MessagesRetried, 1)

		// 更新重试信息
		c.updateRetryInfo(msgID, attempt+1, err.Error())

		// 计算重试延迟（指数退避）
		retryDelay := c.calculateRetryDelay(attempt)

		c.logger.Warn("message processing failed, will retry",
			zap.String("message_id", msgID),
			zap.Int("attempt", attempt+1),
			zap.Error(err),
			zap.Duration("retry_delay", retryDelay),
		)

		// 返回错误让NSQ重新入队
		return err
	}

	// 成功处理
	atomic.AddInt64(&c.metrics.MessagesSuccess, 1)
	c.clearRetryInfo(msgID)

	c.logger.Debug("message processed successfully",
		zap.String("message_id", msgID),
		zap.Duration("processing_time", time.Since(startTime)),
	)

	return nil
}

// calculateRetryDelay 计算重试延迟（指数退避）
func (c *NSQConsumer) calculateRetryDelay(attempt int) time.Duration {
	// 指数退避：2^attempt * base_delay
	delay := c.config.ConsumerRetryDelay * time.Duration(1<<uint(attempt))

	// 最大延迟不超过60秒
	maxDelay := 60 * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// getAttemptCount 获取消息的尝试次数
func (c *NSQConsumer) getAttemptCount(msgID string) int {
	c.retryMu.RLock()
	defer c.retryMu.RUnlock()

	if info, ok := c.retryTracker[msgID]; ok {
		return info.AttemptCount
	}

	// 从NSQ消息属性中获取尝试次数
	return 0
}

// updateRetryInfo 更新重试信息
func (c *NSQConsumer) updateRetryInfo(msgID string, attempt int, lastError string) {
	c.retryMu.Lock()
	defer c.retryMu.Unlock()

	if info, ok := c.retryTracker[msgID]; ok {
		info.AttemptCount = attempt
		info.LastError = lastError
		info.NextRetry = time.Now().Add(c.calculateRetryDelay(attempt))
	} else {
		c.retryTracker[msgID] = &RetryInfo{
			AttemptCount: attempt,
			FirstAttempt: time.Now(),
			LastError:    lastError,
			NextRetry:    time.Now().Add(c.calculateRetryDelay(attempt)),
		}
	}
}

// clearRetryInfo 清除重试信息
func (c *NSQConsumer) clearRetryInfo(msgID string) {
	c.retryMu.Lock()
	defer c.retryMu.Unlock()

	delete(c.retryTracker, msgID)
}

// Stop 优雅停止消费者
func (c *NSQConsumer) Stop() {
	c.stopOnce.Do(func() {
		c.logger.Info("stopping NSQ consumer")

		c.running.Store(false)

		// 停止消费者
		c.consumer.Stop()

		// 等待处理完成
		select {
		case <-c.stopCh:
		case <-time.After(30 * time.Second):
			c.logger.Warn("timeout waiting for consumer to stop")
		}

		c.logger.Info("NSQ consumer stopped")
	})
}

// GetMetrics 获取消费者指标
func (c *NSQConsumer) GetMetrics() *ConsumerMetrics {
	c.metrics.mu.RLock()
	defer c.metrics.mu.RUnlock()

	// 返回副本
	return &ConsumerMetrics{
		MessagesReceived: atomic.LoadInt64(&c.metrics.MessagesReceived),
		MessagesSuccess:  atomic.LoadInt64(&c.metrics.MessagesSuccess),
		MessagesFailed:   atomic.LoadInt64(&c.metrics.MessagesFailed),
		MessagesRetried:  atomic.LoadInt64(&c.metrics.MessagesRetried),
		MessagesToDLQ:    atomic.LoadInt64(&c.metrics.MessagesToDLQ),
	}
}

// IsRunning 检查消费者是否正在运行
func (c *NSQConsumer) IsRunning() bool {
	return c.running.Load()
}

// GetRetryInfo 获取重试信息
func (c *NSQConsumer) GetRetryInfo(msgID string) (*RetryInfo, bool) {
	c.retryMu.RLock()
	defer c.retryMu.RUnlock()

	info, ok := c.retryTracker[msgID]
	if !ok {
		return nil, false
	}

	// 返回副本
	return &RetryInfo{
		AttemptCount: info.AttemptCount,
		FirstAttempt: info.FirstAttempt,
		LastError:    info.LastError,
		NextRetry:    info.NextRetry,
	}, true
}

// GetAllRetryInfo 获取所有重试信息
func (c *NSQConsumer) GetAllRetryInfo() map[string]*RetryInfo {
	c.retryMu.RLock()
	defer c.retryMu.RUnlock()

	result := make(map[string]*RetryInfo, len(c.retryTracker))
	for msgID, info := range c.retryTracker {
		result[msgID] = &RetryInfo{
			AttemptCount: info.AttemptCount,
			FirstAttempt: info.FirstAttempt,
			LastError:    info.LastError,
			NextRetry:    info.NextRetry,
		}
	}

	return result
}

// =====================================================================
// 便捷的消费者包装器
// =====================================================================

// JSONMessageHandler JSON消息处理器
type JSONMessageHandler struct {
	ProcessFunc func(ctx context.Context, msg map[string]interface{}) error
}

// HandleMessage 实现MessageHandler接口
func (h *JSONMessageHandler) HandleMessage(ctx context.Context, message *nsq.Message) error {
	var msg map[string]interface{}
	if err := json.Unmarshal(message.Body, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal JSON message: %w", err)
	}

	return h.ProcessFunc(ctx, msg)
}

// TypedMessageHandler 类型化消息处理器
type TypedMessageHandler[T any] struct {
	ProcessFunc func(ctx context.Context, msg *T) error
}

// HandleMessage 实现MessageHandler接口
func (h *TypedMessageHandler[T]) HandleMessage(ctx context.Context, message *nsq.Message) error {
	var msg T
	if err := json.Unmarshal(message.Body, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return h.ProcessFunc(ctx, &msg)
}

// NewTypedMessageHandler 创建类型化消息处理器
func NewTypedMessageHandler[T any](
	processFunc func(ctx context.Context, msg *T) error,
) MessageHandler {
	return &TypedMessageHandler[T]{
		ProcessFunc: processFunc,
	}
}
