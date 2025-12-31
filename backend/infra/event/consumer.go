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

// Consumer NSQ消费者
type Consumer struct {
	consumer    *nsq.Consumer
	nsqdAddr    string
	lookupdAddr string
	topic       string
	channel     string
	handler     MessageHandler
}

// MessageHandler 消息处理函数
type MessageHandler func(message *nsq.Message) error

// NewConsumer 创建NSQ消费者
func NewConsumer(lookupdAddr, topic, channel string, handler MessageHandler) (*Consumer, error) {
	config := nsq.NewConfig()

	// 配置消费参数
	config.MaxAttempts = 5              // 最大重试次数
	config.DefaultRequeueDelay = 30 * time.Second  // 重试延迟
	config.MaxInFlight = 10               // 最大并发处理数

	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create nsq consumer: %w", err)
	}

	// 注册消息处理器
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		return handler(message)
	}))

	return &Consumer{
		consumer:    consumer,
		lookupdAddr: lookupdAddr,
		topic:       topic,
		channel:     channel,
		handler:     handler,
	}, nil
}

// Start 启动消费者
func (c *Consumer) Start() error {
	log.Printf("[Consumer] Starting consumer for topic=%s channel=%s", c.topic, c.channel)

	err := c.consumer.ConnectToNSQLookupd(c.lookupdAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to nsqlookupd: %w", err)
	}

	return nil
}

// Stop 停止消费者
func (c *Consumer) Stop() {
	log.Printf("[Consumer] Stopping consumer for topic=%s channel=%s", c.topic, c.channel)
	c.consumer.Stop()
}

// ========== 租户事件处理器 ==========

// TenantEventHandler 租户事件处理器
func TenantEventHandler(message *nsq.Message) error {
	var event TenantEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		log.Printf("[TenantEventHandler] Failed to unmarshal message: %v", err)
		message.Requeue(time.Second * 5)
		return err
	}

	log.Printf("[TenantEventHandler] Processing event: %s for tenant: %s", event.EventType, event.TenantID)

	switch event.EventType {
	case "created":
		return onTenantCreated(event)
	case "updated":
		return onTenantUpdated(event)
	case "deleted":
		return onTenantDeleted(event)
	case "suspended":
		return onTenantSuspended(event)
	default:
		log.Printf("[TenantEventHandler] Unknown event type: %s", event.EventType)
		return nil
	}
}

// onTenantCreated 处理租户创建事件
func onTenantCreated(event TenantEvent) error {
	log.Printf("[TenantEventHandler] Tenant created: %s (%s)", event.TenantID, event.Name)
	// TODO: 初始化租户数据、发送欢迎邮件等
	return nil
}

// onTenantUpdated 处理租户更新事件
func onTenantUpdated(event TenantEvent) error {
	log.Printf("[TenantEventHandler] Tenant updated: %s", event.TenantID)
	// TODO: 更新缓存、发送通知等
	return nil
}

// onTenantDeleted 处理租户删除事件
func onTenantDeleted(event TenantEvent) error {
	log.Printf("[TenantEventHandler] Tenant deleted: %s", event.TenantID)
	// TODO: 清理租户数据、归档等
	return nil
}

// onTenantSuspended 处理租户暂停事件
func onTenantSuspended(event TenantEvent) error {
	log.Printf("[TenantEventHandler] Tenant suspended: %s", event.TenantID)
	// TODO: 暂停服务、通知用户等
	return nil
}

// ========== 配额事件处理器 ==========

// QuotaEventHandler 配额事件处理器
func QuotaEventHandler(message *nsq.Message) error {
	var event QuotaEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		log.Printf("[QuotaEventHandler] Failed to unmarshal message: %v", err)
		message.Requeue(time.Second * 5)
		return err
	}

	log.Printf("[QuotaEventHandler] Processing event: %s for tenant: %s", event.EventType, event.TenantID)

	switch event.EventType {
	case "checked":
		return onQuotaChecked(event)
	case "updated":
		return onQuotaUpdated(event)
	case "exceeded":
		return onQuotaExceeded(event)
	case "reset":
		return onQuotaReset(event)
	default:
		log.Printf("[QuotaEventHandler] Unknown event type: %s", event.EventType)
		return nil
	}
}

// onQuotaChecked 处理配额检查事件
func onQuotaChecked(event QuotaEvent) error {
	log.Printf("[QuotaEventHandler] Quota checked: tenant=%s resource=%s usage=%d/%d",
		event.TenantID, event.ResourceType, event.CurrentUsage, event.Limit)
	// TODO: 记录配额使用统计
	return nil
}

// onQuotaUpdated 处理配额更新事件
func onQuotaUpdated(event QuotaEvent) error {
	log.Printf("[QuotaEventHandler] Quota updated: tenant=%s resource=%s amount=%d",
		event.TenantID, event.ResourceType, event.Amount)
	// TODO: 更新配额缓存、发送通知等
	return nil
}

// onQuotaExceeded 处理配额超限事件
func onQuotaExceeded(event QuotaEvent) error {
	log.Printf("[QuotaEventHandler] Quota exceeded: tenant=%s resource=%s usage=%d/%d",
		event.TenantID, event.ResourceType, event.CurrentUsage, event.Limit)
	// TODO: 阻止操作、发送告警、通知用户等
	return nil
}

// onQuotaReset 处理配额重置事件
func onQuotaReset(event QuotaEvent) error {
	log.Printf("[QuotaEventHandler] Quota reset: tenant=%s resource=%s",
		event.TenantID, event.ResourceType)
	// TODO: 重置配额计数
	return nil
}

// ========== 订阅事件处理器 ==========

// SubscriptionEventHandler 订阅事件处理器
func SubscriptionEventHandler(message *nsq.Message) error {
	var event SubscriptionEvent
	if err := json.Unmarshal(message.Body, &event); err != nil {
		log.Printf("[SubscriptionEventHandler] Failed to unmarshal message: %v", err)
		message.Requeue(time.Second * 5)
		return err
	}

	log.Printf("[SubscriptionEventHandler] Processing event: %s for tenant: %s", event.EventType, event.TenantID)

	switch event.EventType {
	case "created":
		return onSubscriptionCreated(event)
	case "renewed":
		return onSubscriptionRenewed(event)
	case "upgraded":
		return onSubscriptionUpgraded(event)
	case "downgraded":
		return onSubscriptionDowngraded(event)
	case "expired":
		return onSubscriptionExpired(event)
	default:
		log.Printf("[SubscriptionEventHandler] Unknown event type: %s", event.EventType)
		return nil
	}
}

// onSubscriptionCreated 处理订阅创建事件
func onSubscriptionCreated(event SubscriptionEvent) error {
	log.Printf("[SubscriptionEventHandler] Subscription created: tenant=%s tier=%s",
		event.TenantID, event.ToTier)
	// TODO: 初始化订阅权益、发送欢迎邮件等
	return nil
}

// onSubscriptionRenewed 处理订阅续费事件
func onSubscriptionRenewed(event SubscriptionEvent) error {
	log.Printf("[SubscriptionEventHandler] Subscription renewed: tenant=%s",
		event.TenantID)
	// TODO: 更新到期时间、发送确认邮件等
	return nil
}

// onSubscriptionUpgraded 处理订阅升级事件
func onSubscriptionUpgraded(event SubscriptionEvent) error {
	log.Printf("[SubscriptionEventHandler] Subscription upgraded: tenant=%s %s->%s",
		event.TenantID, event.FromTier, event.ToTier)
	// TODO: 开通新权益、发送通知等
	return nil
}

// onSubscriptionDowngraded 处理订阅降级事件
func onSubscriptionDowngraded(event SubscriptionEvent) error {
	log.Printf("[SubscriptionEventHandler] Subscription downgraded: tenant=%s %s->%s",
		event.TenantID, event.FromTier, event.ToTier)
	// TODO: 限制权益、发送通知等
	return nil
}

// onSubscriptionExpired 处理订阅过期事件
func onSubscriptionExpired(event SubscriptionEvent) error {
	log.Printf("[SubscriptionEventHandler] Subscription expired: tenant=%s",
		event.TenantID)
	// TODO: 降级到免费版、发送通知等
	return nil
}
