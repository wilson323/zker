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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// NSQProducer NSQ生产者（增强版）
type NSQProducer struct {
	producer *nsq.Producer
	config   *NSQConfig
	logger   *zap.Logger
	mu       sync.RWMutex

	// 连接状态
	connected bool

	// 指标
	metrics *ProducerMetrics
}

// ProducerMetrics 生产者指标
type ProducerMetrics struct {
	mu                sync.RWMutex
	PublishCount      int64
	PublishSuccess    int64
	PublishFailed     int64
	PublishRetryCount int64
	PublishLatency    []time.Duration
}

// NewNSQProducer 创建NSQ生产者
func NewNSQProducer(config *NSQConfig, logger *zap.Logger) (*NSQProducer, error) {
	if config == nil {
		config = DefaultNSQConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid NSQ config: %w", err)
	}

	// 配置NSQ
	nsqConfig := nsq.NewConfig()
	nsqConfig.DialTimeout = config.ProducerTimeout
	nsqConfig.WriteTimeout = config.ProducerTimeout

	// 创建生产者（连接到第一个NSQD地址）
	producer, err := nsq.NewProducer(config.NSQDAddresses[0], nsqConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ producer: %w", err)
	}

	// 设置Ping回调，监控连接状态
	producer.SetLogger(newNSQLogger(logger, LogLevelInfo), toNSQLogLevel(LogLevelInfo))

	p := &NSQProducer{
		producer:  producer,
		config:    config,
		logger:    logger,
		connected: true,
		metrics:   &ProducerMetrics{},
	}

	p.logger.Info("NSQ producer initialized",
		zap.Strings("nsqd_addresses", config.NSQDAddresses),
		zap.Duration("timeout", config.ProducerTimeout),
	)

	return p, nil
}

// Publish 发布消息到指定topic
func (p *NSQProducer) Publish(ctx context.Context, topic string, message interface{}) error {
	startTime := time.Now()

	// 序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal message",
			zap.String("topic", topic),
			zap.Error(err),
		)
		p.incrementMetricFailed()
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 重试发布
	var lastErr error
	for attempt := 0; attempt <= p.config.ProducerRetry; attempt++ {
		if attempt > 0 {
			p.metrics.mu.Lock()
			p.metrics.PublishRetryCount++
			p.metrics.mu.Unlock()

			p.logger.Warn("retrying publish",
				zap.String("topic", topic),
				zap.Int("attempt", attempt),
				zap.Error(lastErr),
			)

			// 等待后重试
			time.Sleep(p.config.ProducerRetryInterval)
		}

		// 发布消息
		err = p.producer.Publish(topic, data)
		if err == nil {
			// 成功
			latency := time.Since(startTime)
			p.incrementMetricSuccess(latency)

			p.logger.Debug("message published successfully",
				zap.String("topic", topic),
				zap.Duration("latency", latency),
				zap.Int("attempt", attempt+1),
			)
			return nil
		}

		lastErr = err
	}

	// 所有重试都失败
	p.incrementMetricFailed()
	p.logger.Error("failed to publish message after retries",
		zap.String("topic", topic),
		zap.Int("retries", p.config.ProducerRetry),
		zap.Error(lastErr),
	)

	return fmt.Errorf("failed to publish after %d retries: %w", p.config.ProducerRetry, lastErr)
}

// PublishDelayed 延迟发布消息
func (p *NSQProducer) PublishDelayed(ctx context.Context, topic string, delay time.Duration, message interface{}) error {
	startTime := time.Now()

	// 序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal message for delayed publish",
			zap.String("topic", topic),
			zap.Error(err),
		)
		p.incrementMetricFailed()
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 延迟发布
	err = p.producer.DeferredPublish(topic, delay, data)
	if err != nil {
		p.incrementMetricFailed()
		p.logger.Error("failed to publish delayed message",
			zap.String("topic", topic),
			zap.Duration("delay", delay),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish delayed message: %w", err)
	}

	// 成功
	latency := time.Since(startTime)
	p.incrementMetricSuccess(latency)

	p.logger.Debug("delayed message published successfully",
		zap.String("topic", topic),
		zap.Duration("delay", delay),
		zap.Duration("latency", latency),
	)

	return nil
}

// PublishBatch 批量发布消息
func (p *NSQProducer) PublishBatch(ctx context.Context, topic string, messages []interface{}) error {
	if len(messages) == 0 {
		return nil
	}

	startTime := time.Now()
	successCount := 0
	failCount := 0

	// 限制并发数
	sem := make(chan struct{}, 10) // 最多10个并发
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for _, msg := range messages {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(message interface{}) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			err := p.Publish(ctx, topic, message)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failCount++
				if firstErr == nil {
					firstErr = err
				}
			} else {
				successCount++
			}
		}(msg)
	}

	wg.Wait()

	latency := time.Since(startTime)

	p.logger.Info("batch publish completed",
		zap.String("topic", topic),
		zap.Int("total", len(messages)),
		zap.Int("success", successCount),
		zap.Int("failed", failCount),
		zap.Duration("total_latency", latency),
	)

	if firstErr != nil {
		return fmt.Errorf("batch publish partially failed: %d/%d succeeded: %w",
			successCount, len(messages), firstErr)
	}

	return nil
}

// PublishMulti 发布多条消息（事务性）
func (p *NSQProducer) PublishMulti(ctx context.Context, topic string, messages [][]byte) error {
	if len(messages) == 0 {
		return nil
	}

	startTime := time.Now()

	// 重试发布
	var lastErr error
	for attempt := 0; attempt <= p.config.ProducerRetry; attempt++ {
		if attempt > 0 {
			p.metrics.mu.Lock()
			p.metrics.PublishRetryCount++
			p.metrics.mu.Unlock()

			time.Sleep(p.config.ProducerRetryInterval)
		}

		// 批量发布
		err := p.producer.MultiPublish(topic, messages)
		if err == nil {
			latency := time.Since(startTime)
			p.incrementMetricSuccess(latency)

			p.logger.Debug("multi-message published successfully",
				zap.String("topic", topic),
				zap.Int("count", len(messages)),
				zap.Duration("latency", latency),
			)
			return nil
		}

		lastErr = err
	}

	p.incrementMetricFailed()
	p.logger.Error("failed to multi-publish after retries",
		zap.String("topic", topic),
		zap.Int("count", len(messages)),
		zap.Error(lastErr),
	)

	return fmt.Errorf("failed to multi-publish after %d retries: %w", p.config.ProducerRetry, lastErr)
}

// IsConnected 检查连接状态
func (p *NSQProducer) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.connected
}

// Reconnect 重新连接
func (p *NSQProducer) Reconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.logger.Info("attempting to reconnect NSQ producer")

	// 停止旧的producer
	if p.producer != nil {
		p.producer.Stop()
	}

	// 创建新的producer
	nsqConfig := nsq.NewConfig()
	nsqConfig.DialTimeout = p.config.ProducerTimeout
	nsqConfig.WriteTimeout = p.config.ProducerTimeout

	producer, err := nsq.NewProducer(p.config.NSQDAddresses[0], nsqConfig)
	if err != nil {
		p.connected = false
		return fmt.Errorf("failed to recreate NSQ producer: %w", err)
	}

	p.producer = producer
	p.connected = true

	p.logger.Info("NSQ producer reconnected successfully")
	return nil
}

// Stop 优雅停止生产者
func (p *NSQProducer) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.connected {
		return
	}

	p.logger.Info("stopping NSQ producer")

	// 停止producer
	p.producer.Stop()

	p.connected = false

	p.logger.Info("NSQ producer stopped")
}

// GetMetrics 获取生产者指标
func (p *NSQProducer) GetMetrics() *ProducerMetrics {
	p.metrics.mu.RLock()
	defer p.metrics.mu.RUnlock()

	// 返回副本，避免并发问题
	return &ProducerMetrics{
		PublishCount:      p.metrics.PublishCount,
		PublishSuccess:    p.metrics.PublishSuccess,
		PublishFailed:     p.metrics.PublishFailed,
		PublishRetryCount: p.metrics.PublishRetryCount,
	}
}

// incrementMetricSuccess 增加成功计数
func (p *NSQProducer) incrementMetricSuccess(latency time.Duration) {
	p.metrics.mu.Lock()
	defer p.metrics.mu.Unlock()

	p.metrics.PublishCount++
	p.metrics.PublishSuccess++

	// 保留最近1000条延迟记录
	if len(p.metrics.PublishLatency) >= 1000 {
		p.metrics.PublishLatency = p.metrics.PublishLatency[1:]
	}
	p.metrics.PublishLatency = append(p.metrics.PublishLatency, latency)
}

// incrementMetricFailed 增加失败计数
func (p *NSQProducer) incrementMetricFailed() {
	p.metrics.mu.Lock()
	defer p.metrics.mu.Unlock()

	p.metrics.PublishCount++
	p.metrics.PublishFailed++
}

// PublishBatchJSON 批量发布JSON消息（便捷方法）
func (p *NSQProducer) PublishBatchJSON(ctx context.Context, topic string, messages []interface{}) error {
	if len(messages) == 0 {
		return nil
	}

	// 序列化所有消息
	dataList := make([][]byte, 0, len(messages))
	for _, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}
		dataList = append(dataList, data)
	}

	// 使用MultiPublish批量发布
	return p.PublishMulti(ctx, topic, dataList)
}

// PublishWithRetry 带自定义重试策略的发布
func (p *NSQProducer) PublishWithRetry(ctx context.Context, topic string, message interface{}, retryStrategy RetryStrategy) error {
	startTime := time.Now()

	// 序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 使用自定义重试策略
	attempt := 0
	var lastErr error

	for {
		err := p.producer.Publish(topic, data)
		if err == nil {
			latency := time.Since(startTime)
			p.incrementMetricSuccess(latency)

			p.logger.Debug("message published with custom retry strategy",
				zap.String("topic", topic),
				zap.Int("attempt", attempt+1),
				zap.Duration("latency", latency),
			)
			return nil
		}

		lastErr = err
		attempt++

		// 检查是否应该重试
		shouldRetry, delay := retryStrategy.ShouldRetry(ctx, attempt, err)
		if !shouldRetry {
			break
		}

		p.metrics.mu.Lock()
		p.metrics.PublishRetryCount++
		p.metrics.mu.Unlock()

		p.logger.Warn("retrying publish with custom strategy",
			zap.String("topic", topic),
			zap.Int("attempt", attempt),
			zap.Duration("delay", delay),
			zap.Error(err),
		)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// 继续重试
		}
	}

	p.incrementMetricFailed()
	p.logger.Error("failed to publish after custom retry strategy",
		zap.String("topic", topic),
		zap.Int("attempts", attempt),
		zap.Error(lastErr),
	)

	return fmt.Errorf("failed to publish after %d attempts: %w", attempt, lastErr)
}

// RetryStrategy 自定义重试策略接口
type RetryStrategy interface {
	ShouldRetry(ctx context.Context, attempt int, err error) (shouldRetry bool, delay time.Duration)
}

// ExponentialBackoffStrategy 指数退避重试策略
type ExponentialBackoffStrategy struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	MaxAttempts  int
}

// ShouldRetry 实现RetryStrategy接口
func (s *ExponentialBackoffStrategy) ShouldRetry(ctx context.Context, attempt int, err error) (bool, time.Duration) {
	if attempt >= s.MaxAttempts {
		return false, 0
	}

	// 计算延迟时间：2^attempt * InitialDelay
	delay := s.InitialDelay * time.Duration(1<<uint(attempt))
	if delay > s.MaxDelay {
		delay = s.MaxDelay
	}

	return true, delay
}

// FixedDelayStrategy 固定延迟重试策略
type FixedDelayStrategy struct {
	Delay        time.Duration
	MaxAttempts  int
}

// ShouldRetry 实现RetryStrategy接口
func (s *FixedDelayStrategy) ShouldRetry(ctx context.Context, attempt int, err error) (bool, time.Duration) {
	if attempt >= s.MaxAttempts {
		return false, 0
	}

	return true, s.Delay
}

// ErrNotConnected 生产者未连接错误
var ErrNotConnected = errors.New("NSQ producer is not connected")
