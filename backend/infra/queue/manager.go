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
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

// QueueManager 队列管理器
type QueueManager struct {
	config   *NSQConfig
	logger   *zap.Logger

	// 生产者和死信队列
	producer *NSQProducer
	dlq      *DeadLetterQueue

	// 消费者注册表
	consumers map[string]*NSQConsumer
	consumersMu sync.RWMutex

	// 任务类型注册表
	taskTypes map[string]TaskHandler
	taskTypesMu sync.RWMutex

	// 状态
	started   atomic.Bool
	shutdown chan struct{}
	once     sync.Once
}

// TaskHandler 任务处理器
type TaskHandler interface {
	HandleTask(ctx context.Context, taskType string, taskData []byte) error
}

// TaskHandlerFunc 函数式任务处理器
type TaskHandlerFunc func(ctx context.Context, taskType string, taskData []byte) error

// HandleTask 实现TaskHandler接口
func (f TaskHandlerFunc) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
	return f(ctx, taskType, taskData)
}

// NewQueueManager 创建队列管理器
func NewQueueManager(config *NSQConfig, logger *zap.Logger) (*QueueManager, error) {
	if config == nil {
		config = DefaultNSQConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid NSQ config: %w", err)
	}

	// 创建生产者
	producer, err := NewNSQProducer(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// 创建死信队列
	dlq := NewDeadLetterQueue(producer, config, logger)

	manager := &QueueManager{
		config:    config,
		logger:    logger,
		producer:  producer,
		dlq:       dlq,
		consumers: make(map[string]*NSQConsumer),
		taskTypes: make(map[string]TaskHandler),
		shutdown:  make(chan struct{}),
	}

	logger.Info("Queue manager initialized",
		zap.Strings("nsqd_addresses", config.NSQDAddresses),
		zap.Strings("nsqlookupd_addresses", config.NSQLookupdAddresses),
	)

	return manager, nil
}

// RegisterTask 注册任务类型
func (m *QueueManager) RegisterTask(taskType string, handler TaskHandler) error {
	m.taskTypesMu.Lock()
	defer m.taskTypesMu.Unlock()

	if _, ok := m.taskTypes[taskType]; ok {
		return fmt.Errorf("task type already registered: %s", taskType)
	}

	m.taskTypes[taskType] = handler

	m.logger.Info("task type registered",
		zap.String("task_type", taskType),
	)

	return nil
}

// RegisterTaskWithConsumer 注册任务类型并创建消费者
func (m *QueueManager) RegisterTaskWithConsumer(taskType string, handler TaskHandler) error {
	m.taskTypesMu.Lock()
	defer m.taskTypesMu.Unlock()

	if _, ok := m.taskTypes[taskType]; ok {
		return fmt.Errorf("task type already registered: %s", taskType)
	}

	// 创建消费者
	consumer, err := m.createConsumer(taskType, handler)
	if err != nil {
		return fmt.Errorf("failed to create consumer for task %s: %w", taskType, err)
	}

	m.taskTypes[taskType] = handler
	m.consumers[taskType] = consumer

	m.logger.Info("task type registered with consumer",
		zap.String("task_type", taskType),
	)

	return nil
}

// createConsumer 创建消费者
func (m *QueueManager) createConsumer(taskType string, handler TaskHandler) (*NSQConsumer, error) {
	// 创建消息处理器
	messageHandler := MessageHandlerFunc(func(ctx context.Context, message *nsq.Message) error {
		return handler.HandleTask(ctx, taskType, message.Body)
	})

	// 创建消费者
	consumer, err := NewNSQConsumer(
		taskType,
		"default", // 默认channel
		m.config,
		messageHandler,
		m.logger,
		m.dlq,
	)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// EnqueueTask 入队任务
func (m *QueueManager) EnqueueTask(ctx context.Context, taskType string, task interface{}) error {
	return m.producer.Publish(ctx, taskType, task)
}

// EnqueueTaskDelayed 入队延迟任务
func (m *QueueManager) EnqueueTaskDelayed(ctx context.Context, taskType string, delay time.Duration, task interface{}) error {
	return m.producer.PublishDelayed(ctx, taskType, delay, task)
}

// EnqueueTaskBatch 批量入队任务
func (m *QueueManager) EnqueueTaskBatch(ctx context.Context, taskType string, tasks []interface{}) error {
	return m.producer.PublishBatch(ctx, taskType, tasks)
}

// Start 启动队列管理器
func (m *QueueManager) Start() error {
	if !m.started.CompareAndSwap(false, true) {
		return fmt.Errorf("queue manager already started")
	}

	m.logger.Info("starting queue manager")

	// 启动所有消费者
	m.consumersMu.RLock()
	consumers := make([]*NSQConsumer, 0, len(m.consumers))
	for _, consumer := range m.consumers {
		consumers = append(consumers, consumer)
	}
	m.consumersMu.RUnlock()

	// 连接所有消费者到NSQLookupd
	for _, consumer := range consumers {
		if err := consumer.Consume(); err != nil {
			m.logger.Error("failed to start consumer",
				zap.Error(err),
			)
			// 继续启动其他消费者
		}
	}

	m.logger.Info("queue manager started",
		zap.Int("consumers", len(consumers)),
	)

	return nil
}

// Stop 优雅关闭队列管理器
func (m *QueueManager) Stop() error {
	m.once.Do(func() {
		m.logger.Info("stopping queue manager")

		close(m.shutdown)

		// 停止所有消费者
		m.consumersMu.Lock()
		for taskType, consumer := range m.consumers {
			m.logger.Info("stopping consumer",
				zap.String("task_type", taskType),
			)
			consumer.Stop()
		}
		m.consumersMu.Unlock()

		// 停止生产者
		m.producer.Stop()

		m.logger.Info("queue manager stopped")
	})

	return nil
}

// GetProducer 获取生产者
func (m *QueueManager) GetProducer() *NSQProducer {
	return m.producer
}

// GetDLQ 获取死信队列
func (m *QueueManager) GetDLQ() *DeadLetterQueue {
	return m.dlq
}

// GetConsumer 获取消费者
func (m *QueueManager) GetConsumer(taskType string) (*NSQConsumer, bool) {
	m.consumersMu.RLock()
	defer m.consumersMu.RUnlock()

	consumer, ok := m.consumers[taskType]
	return consumer, ok
}

// GetAllConsumers 获取所有消费者
func (m *QueueManager) GetAllConsumers() map[string]*NSQConsumer {
	m.consumersMu.RLock()
	defer m.consumersMu.RUnlock()

	result := make(map[string]*NSQConsumer, len(m.consumers))
	for taskType, consumer := range m.consumers {
		result[taskType] = consumer
	}

	return result
}

// GetTaskTypes 获取所有注册的任务类型
func (m *QueueManager) GetTaskTypes() []string {
	m.taskTypesMu.RLock()
	defer m.taskTypesMu.RUnlock()

	taskTypes := make([]string, 0, len(m.taskTypes))
	for taskType := range m.taskTypes {
		taskTypes = append(taskTypes, taskType)
	}

	return taskTypes
}

// GetStats 获取队列统计信息
func (m *QueueManager) GetStats() *QueueStats {
	producerMetrics := m.producer.GetMetrics()
	dlqMetrics := m.dlq.GetMetrics()

	m.consumersMu.RLock()
	consumers := make(map[string]*ConsumerMetrics, len(m.consumers))
	for taskType, consumer := range m.consumers {
		consumers[taskType] = consumer.GetMetrics()
	}
	m.consumersMu.RUnlock()

	return &QueueStats{
		ProducerMetrics:  producerMetrics,
		ConsumerMetrics:  consumers,
		DLQMetrics:       dlqMetrics,
		RegisteredTasks:  m.GetTaskTypes(),
	}
}

// QueueStats 队列统计信息
type QueueStats struct {
	ProducerMetrics  *ProducerMetrics            `json:"producer_metrics"`
	ConsumerMetrics  map[string]*ConsumerMetrics `json:"consumer_metrics"`
	DLQMetrics       *DLQMetrics                 `json:"dlq_metrics"`
	RegisteredTasks  []string                    `json:"registered_tasks"`
}

// IsStarted 检查管理器是否已启动
func (m *QueueManager) IsStarted() bool {
	return m.started.Load()
}

// AddDynamicConsumer 动态添加消费者
func (m *QueueManager) AddDynamicConsumer(taskType string, handler TaskHandler) error {
	m.consumersMu.Lock()
	defer m.consumersMu.Unlock()

	if _, ok := m.consumers[taskType]; ok {
		return fmt.Errorf("consumer already exists for task type: %s", taskType)
	}

	// 创建消费者
	consumer, err := m.createConsumer(taskType, handler)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// 如果管理器已启动，立即启动消费者
	if m.started.Load() {
		if err := consumer.Consume(); err != nil {
			return fmt.Errorf("failed to start consumer: %w", err)
		}
	}

	m.consumers[taskType] = consumer

	m.logger.Info("dynamic consumer added",
		zap.String("task_type", taskType),
	)

	return nil
}

// RemoveConsumer 移除消费者
func (m *QueueManager) RemoveConsumer(taskType string) error {
	m.consumersMu.Lock()
	defer m.consumersMu.Unlock()

	consumer, ok := m.consumers[taskType]
	if !ok {
		return fmt.Errorf("consumer not found for task type: %s", taskType)
	}

	// 停止消费者
	consumer.Stop()

	// 移除
	delete(m.consumers, taskType)

	m.logger.Info("consumer removed",
		zap.String("task_type", taskType),
	)

	return nil
}

// ReplayDLQTasks 重放死信队列中的任务
func (m *QueueManager) ReplayDLQTasks(ctx context.Context, taskType string) error {
	return m.dlq.ReplayDLQ(ctx, taskType)
}

// ReplayDLQTask 重放单条死信任务
func (m *QueueManager) ReplayDLQTask(ctx context.Context, messageID string) error {
	return m.dlq.ReplayMessage(ctx, messageID)
}
