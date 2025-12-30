// backend/infra/mq/nsq_producer.go
package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

// NSQProducer NSQ生产者
type NSQProducer struct {
	producer *nsq.Producer
	topic    string
}

// NewNSQProducer 创建NSQ生产者
func NewNSQProducer(nsqdAddr, topic string) (*NSQProducer, error) {
	producer, err := nsq.NewProducer(nsqdAddr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	log.Printf("NSQ producer initialized: nsqd=%s, topic=%s", nsqdAddr, topic)

	return &NSQProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Publish 发布消息
func (p *NSQProducer) Publish(ctx context.Context, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.producer.Publish(p.topic, data)
}

// PublishDelayed 延迟发布消息
func (p *NSQProducer) PublishDelayed(ctx context.Context, message interface{}, delayMs int) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.producer.DeferredPublish(p.topic, time.Duration(delayMs)*time.Millisecond, data)
}

// Stop 停止生产者
func (p *NSQProducer) Stop() {
	p.producer.Stop()
}

// =====================================================================
// NSQ消费者
// =====================================================================

// NSQConsumer NSQ消费者
type NSQConsumer struct {
	consumer *nsq.Consumer
	handlers map[string]MessageHandler
}

// MessageHandler 消息处理器
type MessageHandler func(message *nsq.Message) error

// NewNSQConsumer 创建NSQ消费者
func NewNSQConsumer(nsqdAddr, lookupdAddr, topic, channel string) (*NSQConsumer, error) {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	// 连接到lookupd（自动发现）
	err = consumer.ConnectToNSQLookupd(lookupdAddr)
	if err != nil {
		return nil, err
	}

	log.Printf("NSQ consumer initialized: topic=%s, channel=%s", topic, channel)

	return &NSQConsumer{
		consumer: consumer,
		handlers: make(map[string]MessageHandler),
	}, nil
}

// RegisterHandler 注册消息处理器
func (c *NSQConsumer) RegisterHandler(topic string, handler MessageHandler) {
	c.handlers[topic] = handler
	c.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		return handler(message)
	}))
}

// Stop 停止消费者
func (c *NSQConsumer) Stop() {
	c.consumer.Stop()
}

// =====================================================================
// Token计量消息队列示例
// =====================================================================

// TokenUsageMessage Token使用消息
type TokenUsageMessage struct {
	TenantID         string  `json:"tenant_id"`
	UserID           string  `json:"user_id"`
	BotID            string  `json:"bot_id"`
	ModelID          string  `json:"model_id"`
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	CostUSD          float64 `json:"cost_usd"`
}

// TokenUsageProducer Token使用生产者
type TokenUsageProducer struct {
	nsqProducer *NSQProducer
}

// NewTokenUsageProducer 创建Token使用生产者
func NewTokenUsageProducer(nsqdAddr string) (*TokenUsageProducer, error) {
	producer, err := NewNSQProducer(nsqdAddr, "token_usage")
	if err != nil {
		return nil, err
	}

	return &TokenUsageProducer{nsqProducer: producer}, nil
}

// PublishTokenUsage 发布Token使用消息（异步）
func (p *TokenUsageProducer) PublishTokenUsage(ctx context.Context, msg *TokenUsageMessage) error {
	return p.nsqProducer.Publish(ctx, msg)
}

// TokenUsageConsumer Token使用消费者
type TokenUsageConsumer struct {
	nsqConsumer *NSQConsumer
	tokenRepo    TokenUsageRepository
}

// NewTokenUsageConsumer 创建Token使用消费者
func NewTokenUsageConsumer(nsqdAddr, lookupdAddr, channel string, tokenRepo TokenUsageRepository) (*TokenUsageConsumer, error) {
	consumer, err := NewNSQConsumer(nsqdAddr, lookupdAddr, "token_usage", channel)
	if err != nil {
		return nil, err
	}

	return &TokenUsageConsumer{
		nsqConsumer: consumer,
		tokenRepo:   tokenRepo,
	}, nil
}

// Start 启动消费者
func (c *TokenUsageConsumer) Start() {
	c.nsqConsumer.RegisterHandler("token_usage", func(message *nsq.Message) error {
		var msg TokenUsageMessage
		if err := json.Unmarshal(message.Body, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		// 写入数据库
		usage := &TokenUsage{
			TenantID:         msg.TenantID,
			UserID:           msg.UserID,
			BotID:            msg.BotID,
			ModelID:          msg.ModelID,
			PromptTokens:     msg.PromptTokens,
			CompletionTokens: msg.CompletionTokens,
			TotalTokens:      msg.TotalTokens,
			CostUSD:          msg.CostUSD,
		}

		if err := c.tokenRepo.Create(context.Background(), usage); err != nil {
			log.Printf("Failed to save token usage: %v", err)
			return err
		}

		log.Printf("Token usage saved: tenant=%s, model=%s, tokens=%d",
			msg.TenantID, msg.ModelID, msg.TotalTokens)

		return nil
	})
}

// 使用示例：
/*
// 生产者（API服务）
producer, _ := NewTokenUsageProducer("localhost:4150")
producer.PublishTokenUsage(ctx, &TokenUsageMessage{
	TenantID:    "tenant_1",
	UserID:      "user_1",
	BotID:       "bot_1",
	ModelID:     "gpt-4",
	TotalTokens: 1000,
	CostUSD:     0.03,
})

// 消费者（后台服务）
consumer, _ := NewTokenUsageConsumer("localhost:4150", "localhost:4161", "channel1", tokenRepo)
consumer.Start()

// 优雅停止
defer consumer.Stop()
select {}
*/
