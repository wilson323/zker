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

// HistoryAction 历史操作类型
type HistoryAction string

const (
	ActionCreated     HistoryAction = "CREATED"      // 任务创建
	ActionAssigned    HistoryAction = "ASSIGNED"     // 任务分配
	ActionAccepted    HistoryAction = "ACCEPTED"     // 接受任务
	ActionRejected    HistoryAction = "REJECTED"     // 拒绝任务
	ActionInProgress  HistoryAction = "IN_PROGRESS"  // 开始处理
	ActionCompleted   HistoryAction = "COMPLETED"    // 任务完成
	ActionEscalated   HistoryAction = "ESCALATED"    // 任务升级
	ActionCancelled   HistoryAction = "CANCELLED"    // 任务取消
	ActionModified    HistoryAction = "MODIFIED"     // 内容修改
	ActionCommented   HistoryAction = "COMMENTED"    // 添加评论
	ActionSLAWarning  HistoryAction = "SLA_WARNING"  // SLA预警
)

// CollaborationHistory 协同历史记录实体
type CollaborationHistory struct {
	ID         int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID     string        `json:"task_id" gorm:"type:varchar(36);index:idx_task_id;not null"`
	TenantID   string        `json:"tenant_id" gorm:"type:varchar(36);not null"`

	// 操作者信息
	ActorType  TaskSource    `json:"actor_type" gorm:"type:enum('AI','HUMAN','SYSTEM');not null"`
	ActorID    string        `json:"actor_id" gorm:"type:varchar(36);not null"`

	// 操作信息
	Action     HistoryAction `json:"action" gorm:"type:varchar(64);not null"`
	Details    JSON          `json:"details" gorm:"type:json"`

	// 时间戳
	Timestamp  int64         `json:"timestamp" gorm:"index;not null;default:0"`
}

// TableName 指定表名
func (CollaborationHistory) TableName() string {
	return "collaboration_history"
}

// GetTimestampAsTime 获取操作时间
func (h *CollaborationHistory) GetTimestampAsTime() int64 {
	return h.Timestamp
}

// HistoryDetail 历史详情
type HistoryDetail struct {
	OldStatus      *TaskStatus     `json:"old_status,omitempty"`
	NewStatus      *TaskStatus     `json:"new_status,omitempty"`
	AssignedFrom   *string         `json:"assigned_from,omitempty"`
	AssignedTo     *string         `json:"assigned_to,omitempty"`
	Comment        string          `json:"comment,omitempty"`
	Reason         string          `json:"reason,omitempty"`
	Changes        map[string]interface{} `json:"changes,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}
