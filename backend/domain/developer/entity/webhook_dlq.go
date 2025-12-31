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
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// WebhookDLQEntry Webhook死信队列条目
type WebhookDLQEntry struct {
	DLQID        string             `json:"dlq_id" gorm:"primaryKey;type:varchar(36);comment:死信队列ID"`
	WebhookID    string             `json:"webhook_id" gorm:"type:varchar(36);not null;index:idx_webhook_id;comment:Webhook ID"`
	EventType    WebhookEventType   `json:"event_type" gorm:"type:varchar(50);not null;index:idx_event_type;comment:事件类型"`
	Payload      *WebhookPayload    `json:"payload" gorm:"type:json;not null;comment:Webhook负载数据"`
	ErrorMessage string             `json:"error_message" gorm:"type:text;not null;comment:错误信息"`
	StatusCode   int                `json:"status_code" gorm:"comment:HTTP状态码"`
	RetryCount   int                `json:"retry_count" gorm:"default:0;comment:已重试次数"`
	LastRetryAt  *int64             `json:"last_retry_at,omitempty" gorm:"comment:最后重试时间"`
	ExpiresAt    *int64             `json:"expires_at,omitempty" gorm:"index:idx_expires_at;comment:过期时间（7天后）"`
	CreatedAt    int64              `json:"created_at" gorm:"not null;default:0;index:idx_created_at;comment:创建时间"`
}

// TableName 指定表名
func (WebhookDLQEntry) TableName() string {
	return "developer_webhook_dlq"
}

// IsExpired 是否已过期
func (e *WebhookDLQEntry) IsExpired() bool {
	if e.ExpiresAt == nil {
		// 默认7天过期
		expiry := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()
		return e.CreatedAt < expiry
	}
	return time.Now().UnixMilli() > *e.ExpiresAt
}

// CanRetry 是否可以重试
func (e *WebhookDLQEntry) CanRetry(maxRetries int) bool {
	return e.RetryCount < maxRetries && !e.IsExpired()
}

// GetLastRetryAtAsTime 获取最后重试时间
func (e *WebhookDLQEntry) GetLastRetryAtAsTime() *time.Time {
	if e.LastRetryAt == nil {
		return nil
	}
	t := time.Unix(*e.LastRetryAt/1000, 0)
	return &t
}

// Scan 实现 sql.Scanner 接口（用于Payload字段）
func (e *WebhookDLQEntry) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan into WebhookDLQEntry")
	}

	return json.Unmarshal(bytes, e)
}

// Value 实现 driver.Valuer 接口（用于Payload字段）
func (e WebhookDLQEntry) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// WebhookRetryConfig Webhook重试配置
type WebhookRetryConfig struct {
	MaxRetries      int           `json:"max_retries" yaml:"max_retries" default:"5"`                   // 最大重试次数
	BaseDelay       time.Duration `json:"base_delay" yaml:"base_delay" default:"1s"`                   // 基础延迟
	MaxDelay        time.Duration `json:"max_delay" yaml:"max_delay" default:"300s"`                   // 最大延迟（5分钟）
	RequestTimeout  time.Duration `json:"request_timeout" yaml:"request_timeout" default:"30s"`        // 请求超时
	DLQRetentionDays int          `json:"dlq_retention_days" yaml:"dlq_retention_days" default:"7"`    // 死信队列保留天数
}

// DefaultWebhookRetryConfig 默认重试配置
func DefaultWebhookRetryConfig() *WebhookRetryConfig {
	return &WebhookRetryConfig{
		MaxRetries:       5,
		BaseDelay:        time.Second,
		MaxDelay:         5 * time.Minute,
		RequestTimeout:   30 * time.Second,
		DLQRetentionDays: 7,
	}
}
