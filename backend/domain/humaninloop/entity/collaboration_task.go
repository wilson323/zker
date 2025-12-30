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
	"time"
)

// TaskType 任务类型
type TaskType string

const (
	TaskTypeReview       TaskType = "REVIEW"        // 内容审核
	TaskTypeCorrection   TaskType = "CORRECTION"    // 错误修正
	TaskTypeValidation   TaskType = "VALIDATION"    // 决策验证
	TaskTypeEscalation   TaskType = "ESCALATION"    // 问题升级
)

// TaskPriority 任务优先级
type TaskPriority string

const (
	PriorityLow    TaskPriority = "LOW"    // 低优先级
	PriorityMedium TaskPriority = "MEDIUM" // 中优先级
	PriorityHigh   TaskPriority = "HIGH"   // 高优先级
	PriorityUrgent TaskPriority = "URGENT" // 紧急
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "PENDING"     // 待分配
	TaskStatusAssigned   TaskStatus = "ASSIGNED"    // 已分配
	TaskStatusInProgress TaskStatus = "IN_PROGRESS" // 处理中
	TaskStatusCompleted  TaskStatus = "COMPLETED"   // 已完成
	TaskStatusEscalated  TaskStatus = "ESCALATED"   // 已升级
	TaskStatusCancelled  TaskStatus = "CANCELLED"   // 已取消
)

// TaskSource 任务来源
type TaskSource string

const (
	SourceAI     TaskSource = "AI"     // AI触发
	SourceHuman  TaskSource = "HUMAN"  // 人工触发
	SourceSystem TaskSource = "SYSTEM" // 系统触发
)

// JSON JSON类型，用于存储复杂数据
type JSON json.RawMessage

// Value 实现 driver.Valuer 接口
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Scan 实现 sql.Scanner 接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	*j = JSON(bytes)
	return nil
}

// String 返回JSON字符串
func (j JSON) String() string {
	if len(j) == 0 {
		return "null"
	}
	return string(j)
}

// CollaborationTask 协同任务实体
type CollaborationTask struct {
	ID           int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID       string       `json:"task_id" gorm:"type:varchar(36);uniqueIndex;not null"`
	TenantID     string       `json:"tenant_id" gorm:"type:varchar(36);index:idx_tenant_status;not null"`
	TaskType     TaskType     `json:"task_type" gorm:"type:enum('REVIEW','CORRECTION','VALIDATION','ESCALATION');not null"`
	Source       TaskSource   `json:"source" gorm:"type:enum('AI','HUMAN','SYSTEM');not null"`
	Priority     TaskPriority `json:"priority" gorm:"type:enum('LOW','MEDIUM','HIGH','URGENT');default:'MEDIUM'"`
	Status       TaskStatus   `json:"status" gorm:"type:enum('PENDING','ASSIGNED','IN_PROGRESS','COMPLETED','ESCALATED','CANCELLED');index:idx_tenant_status;default:'PENDING'"`

	// 分配信息
	AssignedTo   *string      `json:"assigned_to,omitempty" gorm:"type:varchar(36);index:idx_assigned_to"`

	// SLA信息
	SLADeadline  time.Time    `json:"sla_deadline" gorm:"index:idx_sla;not null"`

	// 上下文和结果
	Context      JSON         `json:"context" gorm:"type:json"`
	Result       JSON         `json:"result,omitempty" gorm:"type:json"`

	// 时间戳
	CreatedAt    int64        `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64        `json:"updated_at" gorm:"not null;default:0"`
	CompletedAt  *int64       `json:"completed_at,omitempty"`

	// 关联信息
	ConversationID *string `json:"conversation_id,omitempty" gorm:"type:varchar(36)"`
	MessageID      *string `json:"message_id,omitempty" gorm:"type:varchar(36)"`
	BotID          *string `json:"bot_id,omitempty" gorm:"type:varchar(36)"`
	WorkflowID     *string `json:"workflow_id,omitempty" gorm:"type:varchar(36)"`
}

// TableName 指定表名
func (CollaborationTask) TableName() string {
	return "collaboration_tasks"
}

// IsOverdue 检查是否超时
func (t *CollaborationTask) IsOverdue() bool {
	return time.Now().Unix()*1000 > t.SLADeadline.Unix()*1000
}

// GetRemainingSLA 获取剩余SLA时间（秒）
func (t *CollaborationTask) GetRemainingSLA() int64 {
	deadline := t.SLADeadline.Unix() * 1000
	now := time.Now().Unix() * 1000
	remaining := deadline - now
	if remaining < 0 {
		return 0
	}
	return remaining / 1000
}

// IsAssigned 是否已分配
func (t *CollaborationTask) IsAssigned() bool {
	return t.AssignedTo != nil && *t.AssignedTo != ""
}

// GetCompletedAtAsTime 获取完成时间
func (t *CollaborationTask) GetCompletedAtAsTime() *time.Time {
	if t.CompletedAt == nil {
		return nil
	}
	ts := time.Unix(*t.CompletedAt/1000, 0)
	return &ts
}

// CanBeAssigned 是否可以分配
func (t *CollaborationTask) CanBeAssigned() bool {
	return t.Status == TaskStatusPending || t.Status == TaskStatusCancelled
}

// CanBeCompleted 是否可以完成
func (t *CollaborationTask) CanBeCompleted() bool {
	return t.Status == TaskStatusAssigned || t.Status == TaskStatusInProgress
}

// CanBeEscalated 是否可以升级
func (t *CollaborationTask) CanBeEscalated() bool {
	return t.Status != TaskStatusCompleted && t.Status != TaskStatusCancelled
}
