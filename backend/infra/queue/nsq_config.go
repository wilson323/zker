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

import "time"

// NSQConfig NSQ配置
type NSQConfig struct {
	// NSQD地址列表（支持集群）
	NSQDAddresses []string `yaml:"nsqd_addresses" json:"nsqd_addresses" default:"localhost:4150"`

	// NSQLookupD地址列表（支持集群）
	NSQLookupdAddresses []string `yaml:"nsqlookupd_addresses" json:"nsqlookupd_addresses" default:"localhost:4161"`

	// ============================================================
	// 生产者配置
	// ============================================================

	// 生产者超时时间
	ProducerTimeout time.Duration `yaml:"producer_timeout" json:"producer_timeout" default:"5s"`

	// 生产者重试次数
	ProducerRetry int `yaml:"producer_retry" json:"producer_retry" default:"3"`

	// 生产者重试间隔
	ProducerRetryInterval time.Duration `yaml:"producer_retry_interval" json:"producer_retry_interval" default:"1s"`

	// 连接池最大空闲连接数
	ProducerMaxIdleConns int `yaml:"producer_max_idle_conns" json:"producer_max_idle_conns" default:"10"`

	// ============================================================
	// 消费者配置
	// ============================================================

	// 消费者最大重试次数（超过此次数进入死信队列）
	ConsumerMaxRetries int `yaml:"consumer_max_retries" json:"consumer_max_retries" default:"3"`

	// 消费者重试延迟（指数退避初始值）
	ConsumerRetryDelay time.Duration `yaml:"consumer_retry_delay" json:"consumer_retry_delay" default:"1s"`

	// 消费者最大并发消息数
	ConsumerMaxInFlight int `yaml:"consumer_max_in_flight" json:"consumer_max_in_flight" default:"10"`

	// 消费者消息处理超时
	ConsumerMessageTimeout time.Duration `yaml:"consumer_message_timeout" json:"consumer_message_timeout" default:"60s"`

	// 消费者心跳间隔
	ConsumerHeartbeatInterval time.Duration `yaml:"consumer_heartbeat_interval" json:"consumer_heartbeat_interval" default:"30s"`

	// ============================================================
	// 死信队列配置
	// ============================================================

	// 是否启用死信队列
	DLQEnabled bool `yaml:"dlq_enabled" json:"dlq_enabled" default:"true"`

	// 死信队列Topic后缀（例如：-dlq）
	DLQTopicSuffix string `yaml:"dlq_topic_suffix" json:"dlq_topic_suffix" default:"-dlq"`

	// 死信队列保留时长
	DLQRetentionDuration time.Duration `yaml:"dlq_retention_duration" json:"dlq_retention_duration" default:"72h"`

	// 死信队列最大重放次数
	DLQMaxReplays int `yaml:"dlq_max_replays" json:"dlq_max_replays" default:"3"`

	// ============================================================
	// 监控配置
	// ============================================================

	// 是否启用Prometheus监控
	MetricsEnabled bool `yaml:"metrics_enabled" json:"metrics_enabled" default:"true"`

	// Metrics地址
	MetricsAddress string `yaml:"metrics_address" json:"metrics_address" default:":2112"`

	// ============================================================
	// 日志配置
	// ============================================================

	// 日志级别（debug, info, warn, error）
	LogLevel string `yaml:"log_level" json:"log_level" default:"info"`

	// 是否输出详细日志
	VerboseLogging bool `yaml:"verbose_logging" json:"verbose_logging" default:"false"`
}

// DefaultNSQConfig 返回默认配置
func DefaultNSQConfig() *NSQConfig {
	return &NSQConfig{
		NSQDAddresses:          []string{"localhost:4150"},
		NSQLookupdAddresses:    []string{"localhost:4161"},
		ProducerTimeout:        5 * time.Second,
		ProducerRetry:          3,
		ProducerRetryInterval:  1 * time.Second,
		ProducerMaxIdleConns:   10,
		ConsumerMaxRetries:     3,
		ConsumerRetryDelay:     1 * time.Second,
		ConsumerMaxInFlight:    10,
		ConsumerMessageTimeout: 60 * time.Second,
		ConsumerHeartbeatInterval: 30 * time.Second,
		DLQEnabled:             true,
		DLQTopicSuffix:         "-dlq",
		DLQRetentionDuration:   72 * time.Hour,
		DLQMaxReplays:          3,
		MetricsEnabled:         true,
		MetricsAddress:         ":2112",
		LogLevel:               "info",
		VerboseLogging:         false,
	}
}

// Validate 验证配置
func (c *NSQConfig) Validate() error {
	if len(c.NSQDAddresses) == 0 {
		return ErrInvalidNSQConfig{Field: "nsqd_addresses", Reason: "至少需要一个NSQD地址"}
	}

	if len(c.NSQLookupdAddresses) == 0 {
		return ErrInvalidNSQConfig{Field: "nsqlookupd_addresses", Reason: "至少需要一个NSQLookupd地址"}
	}

	if c.ProducerRetry < 0 {
		return ErrInvalidNSQConfig{Field: "producer_retry", Reason: "重试次数不能为负数"}
	}

	if c.ConsumerMaxRetries < 0 {
		return ErrInvalidNSQConfig{Field: "consumer_max_retries", Reason: "重试次数不能为负数"}
	}

	if c.ConsumerMaxInFlight <= 0 {
		return ErrInvalidNSQConfig{Field: "consumer_max_in_flight", Reason: "最大并发消息数必须大于0"}
	}

	if c.DLQEnabled && c.DLQTopicSuffix == "" {
		return ErrInvalidNSQConfig{Field: "dlq_topic_suffix", Reason: "死信队列Topic后缀不能为空"}
	}

	return nil
}

// ErrInvalidNSQConfig 无效的NSQ配置错误
type ErrInvalidNSQConfig struct {
	Field  string
	Reason string
}

func (e ErrInvalidNSQConfig) Error() string {
	return "invalid NSQ config: " + e.Field + ", " + e.Reason
}
