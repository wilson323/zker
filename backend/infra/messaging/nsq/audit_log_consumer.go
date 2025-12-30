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

package nsq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/coze-dev/coze-studio/backend/domain/audit/entity"
	auditrepo "github.com/coze-dev/coze-studio/backend/domain/audit/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// AuditLogConsumer 审计日志消费者
type AuditLogConsumer struct {
	producer      *NSQProducer
	repo          auditrepo.AuditLogRepository
	maxRetries    int
	retryDelay    time.Duration
	batchSize     int
	flushInterval time.Duration
}

// NewAuditLogConsumer 创建审计日志消费者实例
func NewAuditLogConsumer(
	producer *NSQProducer,
	repo auditrepo.AuditLogRepository,
	maxRetries int,
) *AuditLogConsumer {
	return &AuditLogConsumer{
		producer:      producer,
		repo:          repo,
		maxRetries:    maxRetries,
		retryDelay:    5 * time.Second,
		batchSize:     100,
		flushInterval: 5 * time.Second,
	}
}

// HandleMessage 处理NSQ消息
func (c *AuditLogConsumer) HandleMessage(message *nsq.Message) error {
	ctx := context.Background()

	// 1. 解析消息
	var log entity.AuditLog
	if err := json.Unmarshal(message.Body, &log); err != nil {
		logs.Errorf("[AuditLogConsumer] failed to unmarshal message: %v", err)
		// 消息格式错误，不重试
		return nil
	}

	// 2. 验证必填字段
	if err := c.validateLog(&log); err != nil {
		logs.Errorf("[AuditLogConsumer] invalid audit log: %v", err)
		return nil
	}

	// 3. 保存到数据库(带重试)
	var lastErr error
	for i := 0; i < c.maxRetries; i++ {
		err := c.repo.Create(ctx, &log)
		if err == nil {
			logs.Debugf("[AuditLogConsumer] successfully saved audit log: %s", log.LogID)
			return nil
		}

		lastErr = err
		logs.Warnf("[AuditLogConsumer] failed to save audit log (attempt %d/%d): %v",
			i+1, c.maxRetries, err)

		if i < c.maxRetries-1 {
			time.Sleep(c.retryDelay)
		}
	}

	// 所有重试都失败
	logs.Errorf("[AuditLogConsumer] failed to save audit log after %d attempts: %v",
		c.maxRetries, lastErr)

	// 返回错误会让NSQ重新投递消息
	return fmt.Errorf("failed to save audit log: %w", lastErr)
}

// validateLog 验证日志
func (c *AuditLogConsumer) validateLog(log *entity.AuditLog) error {
	if log.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	if log.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if log.UserEmail == "" {
		return fmt.Errorf("user_email is required")
	}
	if log.Action == "" {
		return fmt.Errorf("action is required")
	}
	if log.Resource == "" {
		return fmt.Errorf("resource is required")
	}
	return nil
}

// Start 启动消费者
func (c *AuditLogConsumer) Start(nsqdAddr, topic, channel string) error {
	config := nsq.NewConfig()
	config.MaxAttempts = c.maxRetries
	config.DefaultRequeueDelay = c.retryDelay

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// 设置消息处理器
	consumer.AddHandler(nsq.HandlerFunc(c.HandleMessage))

	// 连接到NSQ
	if err := consumer.ConnectToNSQD(nsqdAddr); err != nil {
		return fmt.Errorf("failed to connect to NSQ: %w", err)
	}

	logs.Infof("[AuditLogConsumer] started consuming topic=%s, channel=%s", topic, channel)

	return nil
}

// Stop 停止消费者
func (c *AuditLogConsumer) Stop() {
	logs.Infof("[AuditLogConsumer] stopping consumer...")
	// NSQ consumer会在程序退出时自动停止
}

// NSQProducer NSQ生产者(简化实现)
type NSQProducer struct {
	producer *nsq.Producer
}

// NewNSQProducer 创建NSQ生产者
func NewNSQProducer(nsqdAddr string) (*NSQProducer, error) {
	producer, err := nsq.NewProducer(nsqdAddr, nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ producer: %w", err)
	}

	return &NSQProducer{producer: producer}, nil
}

// Publish 发布消息
func (p *NSQProducer) Publish(topic string, data []byte) error {
	return p.producer.Publish(topic, data)
}

// DeferredPublish 延迟发布消息
func (p *NSQProducer) DeferredPublish(topic string, delay time.Duration, data []byte) error {
	return p.producer.DeferredPublish(topic, delay, data)
}

// Stop 停止生产者
func (p *NSQProducer) Stop() {
	p.producer.Stop()
}
