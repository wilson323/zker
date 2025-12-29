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

package event

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

// Producer NSQ生产者
type Producer struct {
	producer *nsq.Producer
	nsqdAddr string
}

// NewProducer 创建NSQ生产者
func NewProducer(nsqdAddr string) (*Producer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create nsq producer: %w", err)
	}

	// 测试连接
	if err := producer.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping nsqd: %w", err)
	}

	return &Producer{
		producer: producer,
		nsqdAddr: nsqdAddr,
	}, nil
}

// PublishTenantEvent 发布租户事件
func (p *Producer) PublishTenantEvent(event TenantEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal tenant event: %w", err)
	}

	err = p.producer.Publish("tenant-events", data)
	if err != nil {
		return fmt.Errorf("failed to publish tenant event: %w", err)
	}

	log.Printf("[Producer] Published tenant event: %s", event.EventType)
	return nil
}

// PublishQuotaEvent 发布配额事件
func (p *Producer) PublishQuotaEvent(event QuotaEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal quota event: %w", err)
	}

	err = p.producer.Publish("quota-events", data)
	if err != nil {
		return fmt.Errorf("failed to publish quota event: %w", err)
	}

	log.Printf("[Producer] Published quota event: %s", event.EventType)
	return nil
}

// PublishSubscriptionEvent 发布订阅事件
func (p *Producer) PublishSubscriptionEvent(event SubscriptionEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription event: %w", err)
	}

	err = p.producer.Publish("subscription-events", data)
	if err != nil {
		return fmt.Errorf("failed to publish subscription event: %w", err)
	}

	log.Printf("[Producer] Published subscription event: %s", event.EventType)
	return nil
}

// Close 关闭生产者
func (p *Producer) Close() {
	p.producer.Stop()
}

// ========== 事件定义 ==========

// TenantEvent 租户事件
type TenantEvent struct {
	EventType string    `json:"event_type"` // created, updated, deleted, suspended
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// QuotaEvent 配额事件
type QuotaEvent struct {
	EventType    string    `json:"event_type"` // checked, updated, exceeded, reset
	TenantID     string    `json:"tenant_id"`
	ResourceType string    `json:"resource_type"` // bot, knowledge, api_call
	Amount       int       `json:"amount"`
	CurrentUsage int       `json:"current_usage"`
	Limit        int       `json:"limit"`
	Timestamp    time.Time `json:"timestamp"`
}

// SubscriptionEvent 订阅事件
type SubscriptionEvent struct {
	EventType    string    `json:"event_type"` // created, renewed, upgraded, downgraded, expired
	TenantID     string    `json:"tenant_id"`
	SubscriptionID string  `json:"subscription_id"`
	FromTier     string    `json:"from_tier,omitempty"`
	ToTier       string    `json:"to_tier"`
	Timestamp    time.Time `json:"timestamp"`
}
