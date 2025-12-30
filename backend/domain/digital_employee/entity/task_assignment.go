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

// TaskType 任务类型
type TaskType string

const (
	TaskTypeCustomerService TaskType = "customer_service" // 客服咨询
	TaskTypeSalesFollow     TaskType = "sales_follow"     // 销售跟进
	TaskTypeTechSupport     TaskType = "tech_support"     // 技术支持
	TaskTypeConsultation    TaskType = "consultation"     // 业务咨询
	TaskTypeTraining        TaskType = "training"         // 培训指导
)

// TaskPriority 任务优先级
type TaskPriority string

const (
	TaskPriorityHigh   TaskPriority = "high"   // 高优先级
	TaskPriorityMedium TaskPriority = "medium" // 中优先级
	TaskPriorityLow    TaskPriority = "low"    // 低优先级
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusAssigned    TaskStatus = "assigned"     // 已分配
	TaskStatusInProgress  TaskStatus = "in_progress"  // 进行中
	TaskStatusCompleted   TaskStatus = "completed"    // 已完成
	TaskStatusFailed      TaskStatus = "failed"       // 失败
	TaskStatusCancelled   TaskStatus = "cancelled"    // 已取消
)

// TaskAssignment 任务分配实体
type TaskAssignment struct {
	AssignmentID  string       `json:"assignment_id" gorm:"primaryKey;type:varchar(36)"`
	TaskID        string       `json:"task_id" gorm:"type:varchar(36);not null;index:idx_task_id"`
	EmployeeID    string       `json:"employee_id" gorm:"type:varchar(36);not null;index:idx_employee_id"`
	TenantID      string       `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	TaskType      TaskType     `json:"task_type" gorm:"type:enum('customer_service','sales_follow','tech_support','consultation','training');not null;index:idx_task_type"`
	Priority      TaskPriority `json:"priority" gorm:"type:enum('high','medium','low');not null;index:idx_priority"`
	Status        TaskStatus   `json:"status" gorm:"type:enum('assigned','in_progress','completed','failed','cancelled');default:'assigned';index:idx_status"`
	AssignedAt    int64        `json:"assigned_at" gorm:"not null;default:0"`
	CompletedAt   *int64       `json:"completed_at,omitempty"`
	Result        string       `json:"result" gorm:"type:text"`
	FailureReason string       `json:"failure_reason,omitempty" gorm:"type:varchar(255)"`
}

// TableName 指定表名
func (TaskAssignment) TableName() string {
	return "digital_employee_task_assignments"
}

// IsCompleted 是否已完成
func (t *TaskAssignment) IsCompleted() bool {
	return t.Status == TaskStatusCompleted
}

// IsFailed 是否已失败
func (t *TaskAssignment) IsFailed() bool {
	return t.Status == TaskStatusFailed
}

// GetAssignedAtAsTime 获取分配时间
func (t *TaskAssignment) GetAssignedAtAsTime() time.Time {
	return time.Unix(t.AssignedAt/1000, 0)
}

// GetCompletedAtAsTime 获取完成时间
func (t *TaskAssignment) GetCompletedAtAsTime() *time.Time {
	if t.CompletedAt == nil {
		return nil
	}
	completed := time.Unix(*t.CompletedAt/1000, 0)
	return &completed
}

// AssignTaskRequest 分配任务请求
type AssignTaskRequest struct {
	TaskID     string       `json:"task_id" binding:"required,uuid"`
	EmployeeID string       `json:"employee_id" binding:"required,uuid"`
	TaskType   TaskType     `json:"task_type" binding:"required,oneof=customer_service sales_follow tech_support consultation training"`
	Priority   TaskPriority `json:"priority" binding:"required,oneof=high medium low"`
}

// AutoAssignTaskRequest 自动分配任务请求
type AutoAssignTaskRequest struct {
	TaskID        string   `json:"task_id" binding:"required,uuid"`
	TaskType      TaskType `json:"task_type" binding:"required,oneof=customer_service sales_follow tech_support consultation training"`
	Priority      TaskPriority `json:"priority" binding:"required,oneof=high medium low"`
	RequiredSkills []string `json:"required_skills" binding:"required"`
}

// CompleteTaskRequest 完成任务请求
type CompleteTaskRequest struct {
	AssignmentID string `json:"assignment_id" binding:"required,uuid"`
	Result       string `json:"result" binding:"required"`
}

// FailTaskRequest 标记任务失败请求
type FailTaskRequest struct {
	AssignmentID  string `json:"assignment_id" binding:"required,uuid"`
	FailureReason string `json:"failure_reason" binding:"required,max=255"`
}

// GetAssignmentsRequest 获取员工任务列表请求
type GetAssignmentsRequest struct {
	EmployeeID string     `json:"employee_id" binding:"required,uuid"`
	Status     *TaskStatus `json:"status" binding:"omitempty,oneof=assigned in_progress completed failed cancelled"`
	Page       int        `json:"page" binding:"required,min=1"`
	PageSize   int        `json:"page_size" binding:"required,min=1,max=100"`
}

// GetAssignmentsResponse 获取员工任务列表响应
type GetAssignmentsResponse struct {
	Assignments []*TaskAssignment `json:"assignments"`
	Total       int64             `json:"total"`
	Page        int               `json:"page"`
	PageSize    int               `json:"page_size"`
}
