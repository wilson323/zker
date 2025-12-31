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
	"time"
)

// CollaborationSession 人机协同会话
type CollaborationSession struct {
	SessionID      string    `json:"session_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID       string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID         string    `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	ConversationID string    `json:"conversation_id" gorm:"type:varchar(36);index:idx_conversation_id"`
	MessageID      string    `json:"message_id" gorm:"type:varchar(36);index:idx_message_id"`
	AgentID        string    `json:"agent_id" gorm:"type:varchar(36);index:idx_agent_id"`

	// 触发信息
	TriggerType    string    `json:"trigger_type" gorm:"type:varchar(50);not null"` // low_confidence, user_initiated, exception_detected, sensitive_content, security_risk
	TriggerReason  string    `json:"trigger_reason" gorm:"type:text"`
	Confidence     float64   `json:"confidence" gorm:"type:decimal(5,4)"`

	// 状态信息
	Status         string    `json:"status" gorm:"type:varchar(20);not null;index:idx_status"` // pending, assigned, completed, expired
	Priority       string    `json:"priority" gorm:"type:varchar(20);not null"` // critical, high, medium, low
	AssignedTo     string    `json:"assigned_to,omitempty" gorm:"type:varchar(36);index:idx_assigned_to"`

	// 协同内容
	InputText      string    `json:"input_text" gorm:"type:text"`
	AIResponse     string    `json:"ai_response" gorm:"type:text"`
	HumanResponse  string    `json:"human_response,omitempty" gorm:"type:text"`
	Resolution     string    `json:"resolution,omitempty" gorm:"type:varchar(50)" // approved, corrected, rejected, escalated

	// 时间信息
	CreatedAt      time.Time `json:"created_at" gorm:"not null;index:idx_created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`
	ExpiresAt      time.Time `json:"expires_at" gorm:"index:idx_expires_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`

	// 额外信息
	Metadata       string    `json:"metadata,omitempty" gorm:"type:json"`
}

// TableName 指定表名
func (CollaborationSession) TableName() string {
	return "collaboration_sessions"
}

// IsPending 是否待处理
func (c *CollaborationSession) IsPending() bool {
	return c.Status == "pending"
}

// IsExpired 是否已过期
func (c *CollaborationSession) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// CollaborationLog 协同日志
type CollaborationLog struct {
	LogID        string    `json:"log_id" gorm:"primaryKey;type:varchar(36)"`
	SessionID    string    `json:"session_id" gorm:"type:varchar(36);not null;index:idx_session_id"`
	TenantID     string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID       string    `json:"user_id" gorm:"type:varchar(36)"`
	Action       string    `json:"action" gorm:"type:varchar(50);not null"` // created, assigned, viewed, responded, completed, expired
	ActionDetail string    `json:"action_detail" gorm:"type:text"`
	Timestamp    time.Time `json:"timestamp" gorm:"not null;index:idx_timestamp"`
}

// TableName 指定表名
func (CollaborationLog) TableName() string {
	return "collaboration_logs"
}
