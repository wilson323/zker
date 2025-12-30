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

// WebhookStatus Webhook状态
type WebhookStatus string

const (
	WebhookStatusActive WebhookStatus = "active" // 激活
	WebhookStatusPaused WebhookStatus = "paused" // 暂停
)

// WebhookEventType Webhook事件类型
type WebhookEventType string

const (
	WebhookEventBotCreated         WebhookEventType = "bot.created"          // Bot创建
	WebhookEventBotUpdated         WebhookEventType = "bot.updated"          // Bot更新
	WebhookEventBotDeleted         WebhookEventType = "bot.deleted"          // Bot删除
	WebhookEventWorkflowCompleted  WebhookEventType = "workflow.completed"   // 工作流完成
	WebhookEventWorkflowFailed     WebhookEventType = "workflow.failed"      // 工作流失败
	WebhookEventConversationCreated WebhookEventType = "conversation.created" // 对话创建
	WebhookEventMessageReceived    WebhookEventType = "message.received"     // 消息接收
	WebhookEventAPICall            WebhookEventType = "api.call"             // API调用
	WebhookEventErrorOccurred      WebhookEventType = "error.occurred"       // 错误发生
)

// WebhookEvents 自定义类型，用于JSON处理
type WebhookEvents []WebhookEventType

// Value 实现 driver.Valuer 接口
func (e WebhookEvents) Value() (driver.Value, error) {
	if len(e) == 0 {
		return "[]", nil
	}
	return json.Marshal(e)
}

// Scan 实现 sql.Scanner 接口
func (e *WebhookEvents) Scan(value interface{}) error {
	if value == nil {
		*e = WebhookEvents{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan into WebhookEvents")
	}

	return json.Unmarshal(bytes, e)
}

// Contains 检查是否包含指定事件
func (e WebhookEvents) Contains(eventType WebhookEventType) bool {
	for _, et := range e {
		if et == eventType {
			return true
		}
	}
	return false
}

// Webhook Webhook实体
type Webhook struct {
	WebhookID    string        `json:"webhook_id" gorm:"primaryKey;type:varchar(36);comment:Webhook ID"`
	TenantID     string        `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id;comment:租户ID"`
	ProjectID    string        `json:"project_id" gorm:"type:varchar(36);not null;index:idx_project_id;comment:项目ID"`
	WebhookURL   string        `json:"webhook_url" gorm:"type:varchar(500);not null;comment:Webhook URL"`
	WebhookSecret string        `json:"-" gorm:"type:varchar(100);not null;comment:Webhook密钥（用于签名验证）"` // 不暴露给前端
	Events       WebhookEvents `json:"events" gorm:"type:json;not null;comment:事件类型"`
	Status       WebhookStatus `json:"status" gorm:"type:enum('active','paused');default:'active';not null;index:idx_status;comment:状态"`
	LastTriggerAt *int64       `json:"last_triggered_at,omitempty" gorm:"comment:最后触发时间"`
	SuccessCount int64         `json:"success_count" gorm:"default:0;comment:成功次数"`
	FailureCount int64         `json:"failure_count" gorm:"default:0;comment:失败次数"`
	CreatedAt    int64         `json:"created_at" gorm:"not null;default:0;comment:创建时间"`
	UpdatedAt    int64         `json:"updated_at" gorm:"not null;default:0;comment:更新时间"`
	DeletedAt    *int64        `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`

	// 关联
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID;references:ProjectID"`
}

// TableName 指定表名
func (Webhook) TableName() string {
	return "developer_webhooks"
}

// IsActive 是否激活
func (w *Webhook) IsActive() bool {
	return w.Status == WebhookStatusActive
}

// ShouldTrigger 检查是否应该触发指定事件
func (w *Webhook) ShouldTrigger(eventType WebhookEventType) bool {
	return w.IsActive() && w.Events.Contains(eventType)
}

// GetSuccessRate 获取成功率
func (w *Webhook) GetSuccessRate() float64 {
	total := w.SuccessCount + w.FailureCount
	if total == 0 {
		return 0
	}
	return float64(w.SuccessCount) / float64(total) * 100
}

// GetLastTriggeredAtAsTime 获取最后触发时间
func (w *Webhook) GetLastTriggeredAtAsTime() *time.Time {
	if w.LastTriggerAt == nil {
		return nil
	}
	t := time.Unix(*w.LastTriggerAt/1000, 0)
	return &t
}

// WebhookPayload Webhook负载数据
type WebhookPayload struct {
	EventID      string                 `json:"event_id"`
	EventType    WebhookEventType       `json:"event_type"`
	Timestamp    int64                  `json:"timestamp"`
	TenantID     string                 `json:"tenant_id"`
	ProjectID    string                 `json:"project_id"`
	EventTypeV1 string                 `json:"event_type_v1,omitempty"` // 兼容v1版本
	Data         map[string]interface{} `json:"data"`
}

// WebhookLog Webhook日志
type WebhookLog struct {
	LogID       string    `json:"log_id" gorm:"primaryKey;type:varchar(36);comment:日志ID"`
	WebhookID   string    `json:"webhook_id" gorm:"type:varchar(36);not null;index:idx_webhook_id;comment:Webhook ID"`
	EventType   string    `json:"event_type" gorm:"type:varchar(50);not null;index:idx_event_type;comment:事件类型"`
	StatusCode  int       `json:"status_code" gorm:"comment:HTTP状态码"`
	Response    string    `json:"response" gorm:"type:text;comment:响应内容"`
	DurationMs  int64     `json:"duration_ms" gorm:"comment:耗时（毫秒）"`
	Success     bool      `json:"success" gorm:"not null;default:false;comment:是否成功"`
	RetryCount  int       `json:"retry_count" gorm:"default:0;comment:重试次数"`
	CreatedAt   int64     `json:"created_at" gorm:"not null;default:0;index:idx_created_at;comment:创建时间"`
}

// TableName 指定表名
func (WebhookLog) TableName() string {
	return "developer_webhook_logs"
}
