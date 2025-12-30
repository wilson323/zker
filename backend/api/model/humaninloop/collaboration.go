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

package humaninloop

import (
	"encoding/json"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// ==================== 请求模型 ====================

// CreateTaskRequest 创建协同任务请求
type CreateTaskRequest struct {
	TaskType        entity.TaskType     `json:"task_type" binding:"required"`
	Source          entity.TaskSource   `json:"source" binding:"required"`
	Priority        entity.TaskPriority `json:"priority"`
	ConversationID  *string             `json:"conversation_id,omitempty"`
	MessageID       *string             `json:"message_id,omitempty"`
	BotID           *string             `json:"bot_id,omitempty"`
	WorkflowID      *string             `json:"workflow_id,omitempty"`
	Context         json.RawMessage     `json:"context" binding:"required"`
	TriggerReason   string              `json:"trigger_reason,omitempty"`
	AssignTo        *string             `json:"assign_to,omitempty"`
	SLAMinutes      int                 `json:"sla_minutes,omitempty"`
}

// AssignTaskRequest 分配任务请求
type AssignTaskRequest struct {
	AssignTo string `json:"assign_to" binding:"required"`
	Reason   string `json:"reason,omitempty"`
}

// ReviewTaskRequest 审核任务请求
type ReviewTaskRequest struct {
	Decision        string          `json:"decision" binding:"required,oneof=approved rejected modified"`
	HumanDecision   json.RawMessage `json:"human_decision" binding:"required"`
	DecisionReason  string          `json:"decision_reason,omitempty"`
	ModifiedContent *string         `json:"modified_content,omitempty"`
}

// EscalateTaskRequest 升级任务请求
type EscalateTaskRequest struct {
	Reason       string `json:"reason" binding:"required"`
	EscalateTo   string `json:"escalate_to" binding:"required"`
	Priority     *entity.TaskPriority `json:"priority,omitempty"`
}

// ListTasksRequest 任务列表请求
type ListTasksRequest struct {
	Status      *entity.TaskStatus     `json:"status,omitempty"`
	Priority    *entity.TaskPriority   `json:"priority,omitempty"`
	TaskType    *entity.TaskType       `json:"task_type,omitempty"`
	AssignedTo  *string                `json:"assigned_to,omitempty"`
	IsOverdue   *bool                  `json:"is_overdue,omitempty"`
	Page        int                    `json:"page" binding:"min=1"`
	PageSize    int                    `json:"page_size" binding:"min=1,max=100"`
	SortBy      string                 `json:"sort_by,omitempty" binding:"omitempty,oneof=created_at updated_at sla_deadline priority"`
	SortOrder   string                 `json:"sort_order,omitempty" binding:"omitempty,oneof=asc desc"`
}

// CreateRuleRequest 创建规则请求
type CreateRuleRequest struct {
	Name        string          `json:"name" binding:"required,max=100"`
	RuleType    string          `json:"rule_type" binding:"required,oneof=confidence keyword category custom"`
	Config      json.RawMessage `json:"config" binding:"required"`
	AssignTo    *string         `json:"assign_to,omitempty"`
	Priority    entity.TaskPriority `json:"priority" binding:"required"`
	IsActive    bool            `json:"is_active"`
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	AutoReviewThreshold   *float64       `json:"auto_review_threshold,omitempty" binding:"omitempty,min=0,max=1"`
	AutoEscalateThreshold *float64       `json:"auto_escalate_threshold,omitempty" binding:"omitempty,min=0,max=1"`
	DefaultSLAMinutes     *int           `json:"default_sla_minutes,omitempty" binding:"omitempty,min=1"`
	UrgentSLAMinutes      *int           `json:"urgent_sla_minutes,omitempty" binding:"omitempty,min=1"`
	ReviewPoolFilter      *json.RawMessage `json:"review_pool_filter,omitempty"`
	EscalationRules       *json.RawMessage `json:"escalation_rules,omitempty"`
	PriorityRules         *json.RawMessage `json:"priority_rules,omitempty"`
	NotificationEnabled   *bool          `json:"notification_enabled,omitempty"`
	NotificationChannels  *json.RawMessage `json:"notification_channels,omitempty"`
}

// GetAnalyticsRequest 获取分析数据请求
type GetAnalyticsRequest struct {
	StartDate int64 `json:"start_date" binding:"required"`
	EndDate   int64 `json:"end_date" binding:"required"`
	TaskType  *entity.TaskType `json:"task_type,omitempty"`
}

// ==================== 响应模型 ====================

// TaskResponse 任务响应
type TaskResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *entity.CollaborationTask `json:"data,omitempty"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Code        int                      `json:"code"`
	Message     string                   `json:"message"`
	Data        *TaskListData            `json:"data,omitempty"`
}

type TaskListData struct {
	Tasks      []*entity.CollaborationTask `json:"tasks"`
	Total      int64                       `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
}

// HistoryResponse 历史记录响应
type HistoryResponse struct {
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Data    *HistoryData                 `json:"data,omitempty"`
}

type HistoryData struct {
	TaskID   string                     `json:"task_id"`
	Histories []*entity.CollaborationHistory `json:"histories"`
	Total    int                        `json:"total"`
}

// QueueStatsResponse 队列统计响应
type QueueStatsResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *QueueStats `json:"data,omitempty"`
}

type QueueStats struct {
	TotalTasks      int64                 `json:"total_tasks"`
	PendingTasks    int64                 `json:"pending_tasks"`
	InProgressTasks int64                 `json:"in_progress_tasks"`
	CompletedTasks  int64                 `json:"completed_tasks"`
	OverdueTasks    int64                 `json:"overdue_tasks"`
	ByPriority      map[string]int64      `json:"by_priority"`
	ByStatus        map[string]int64      `json:"by_status"`
	ByType          map[string]int64      `json:"by_type"`
	AverageWaitTime float64               `json:"average_wait_time"`
	AverageHandleTime float64             `json:"average_handle_time"`
}

// AnalyticsResponse 分析数据响应
type AnalyticsResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    *AnalyticsData  `json:"data,omitempty"`
}

type AnalyticsData struct {
	TotalTasks        int64                      `json:"total_tasks"`
	CompletedTasks    int64                      `json:"completed_tasks"`
	PendingTasks      int64                      `json:"pending_tasks"`
	EscalatedTasks    int64                      `json:"escalated_tasks"`
	CompletionRate    float64                    `json:"completion_rate"`
	AverageHandleTime float64                    `json:"average_handle_time"`
	ByPriority        map[string]*PriorityStats  `json:"by_priority"`
	ByType            map[string]*TypeStats      `json:"by_type"`
	TopReviewers      []*ReviewerStats           `json:"top_reviewers"`
}

type PriorityStats struct {
	Total          int64   `json:"total"`
	Completed      int64   `json:"completed"`
	AverageTime    float64 `json:"average_time"`
	OverdueCount   int64   `json:"overdue_count"`
}

type TypeStats struct {
	Total       int64   `json:"total"`
	Completed   int64   `json:"completed"`
	AverageTime float64 `json:"average_time"`
}

type ReviewerStats struct {
	ReviewerID      string  `json:"reviewer_id"`
	ReviewerName    string  `json:"reviewer_name"`
	TaskCount       int64   `json:"task_count"`
	AverageTime     float64 `json:"average_time"`
	ApprovalRate    float64 `json:"approval_rate"`
}

// OverdueTaskResponse 超时任务响应
type OverdueTaskResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    []*OverdueTask `json:"data,omitempty"`
}

type OverdueTask struct {
	TaskID       string              `json:"task_id"`
	TenantID     string              `json:"tenant_id"`
	TaskType     entity.TaskType     `json:"task_type"`
	Priority     entity.TaskPriority `json:"priority"`
	Status       entity.TaskStatus   `json:"status"`
	AssignedTo   *string             `json:"assigned_to,omitempty"`
	OverdueDuration int64            `json:"overdue_duration"` // 超时时长（秒）
}

// ConfigResponse 配置响应
type ConfigResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    *entity.CollaborationConfig `json:"data,omitempty"`
}

// ==================== 通用响应 ====================

// BaseResponse 基础响应
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse 成功响应
func SuccessResponse(message string) *BaseResponse {
	return &BaseResponse{
		Code:    0,
		Message: message,
	}
}
