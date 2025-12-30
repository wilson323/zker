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

package repository

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// CollaborationRepository 协同任务仓储接口
type CollaborationRepository interface {
	// Create 创建协同任务
	Create(ctx context.Context, task *entity.CollaborationTask) error

	// GetByID 根据ID获取任务
	GetByID(ctx context.Context, taskID string) (*entity.CollaborationTask, error)

	// Update 更新任务
	Update(ctx context.Context, task *entity.CollaborationTask) error

	// UpdateStatus 更新任务状态
	UpdateStatus(ctx context.Context, taskID string, status entity.TaskStatus) error

	// AssignTask 分配任务
	AssignTask(ctx context.Context, taskID, assignedTo string) error

	// List 获取任务列表
	List(ctx context.Context, filter *ListFilter) ([]*entity.CollaborationTask, int64, error)

	// GetPendingTasks 获取待分配任务
	GetPendingTasks(ctx context.Context, tenantID string) ([]*entity.CollaborationTask, error)

	// GetTasksByAssignee 获取分配给指定用户的任务
	GetTasksByAssignee(ctx context.Context, tenantID, assigneeID string) ([]*entity.CollaborationTask, error)

	// GetOverdueTasks 获取超时任务
	GetOverdueTasks(ctx context.Context) ([]*entity.CollaborationTask, error)

	// GetStats 获取任务统计
	GetStats(ctx context.Context, tenantID string) (*TaskStats, error)
}

// HistoryRepository 历史记录仓储接口
type HistoryRepository interface {
	// Create 创建历史记录
	Create(ctx context.Context, history *entity.CollaborationHistory) error

	// GetByTaskID 获取任务的历史记录
	GetByTaskID(ctx context.Context, taskID string) ([]*entity.CollaborationHistory, error)

	// GetByTaskIDWithPagination 分页获取任务历史
	GetByTaskIDWithPagination(
		ctx context.Context,
		taskID string,
		page, pageSize int,
	) ([]*entity.CollaborationHistory, int64, error)

	// GetByAction 获取指定操作类型的历史记录
	GetByAction(
		ctx context.Context,
		taskID string,
		action entity.HistoryAction,
	) ([]*entity.CollaborationHistory, error)

	// DeleteOldHistory 删除旧的历史记录
	DeleteOldHistory(ctx context.Context, beforeDate int64) (int64, error)
}

// ConfigRepository 配置仓储接口
type ConfigRepository interface {
	// GetByTenantID 获取租户配置
	GetByTenantID(ctx context.Context, tenantID string) (*entity.CollaborationConfig, error)

	// Create 创建配置
	Create(ctx context.Context, config *entity.CollaborationConfig) error

	// Update 更新配置
	Update(ctx context.Context, config *entity.CollaborationConfig) error

	// Upsert 创建或更新配置
	Upsert(ctx context.Context, config *entity.CollaborationConfig) error
}

// ==================== 数据结构 ====================

// ListFilter 列表查询过滤器
type ListFilter struct {
	TenantID   string
	Status     *entity.TaskStatus
	Priority   *entity.TaskPriority
	TaskType   *entity.TaskType
	AssignedTo *string
	IsOverdue  *bool
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
}

// TaskStats 任务统计
type TaskStats struct {
	TotalTasks      int64            `json:"total_tasks"`
	PendingTasks    int64            `json:"pending_tasks"`
	InProgressTasks int64            `json:"in_progress_tasks"`
	CompletedTasks  int64            `json:"completed_tasks"`
	OverdueTasks    int64            `json:"overdue_tasks"`
	ByPriority      map[string]int64 `json:"by_priority"`
	ByStatus        map[string]int64 `json:"by_status"`
}
