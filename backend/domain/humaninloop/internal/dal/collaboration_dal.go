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

package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// CollaborationDAL 协同任务数据访问层
type CollaborationDAL struct {
	db *gorm.DB
}

// NewCollaborationDAL 创建协同任务DAL
func NewCollaborationDAL(db *gorm.DB) *CollaborationDAL {
	return &CollaborationDAL{db: db}
}

// ==================== 任务CRUD操作 ====================

// Create 创建协同任务
func (d *CollaborationDAL) Create(ctx context.Context, task *entity.CollaborationTask) error {
	now := time.Now().Unix() * 1000
	task.CreatedAt = now
	task.UpdatedAt = now

	if err := d.db.WithContext(ctx).Create(task).Error; err != nil {
		return fmt.Errorf("failed to create collaboration task: %w", err)
	}
	return nil
}

// GetByID 根据ID获取任务
func (d *CollaborationDAL) GetByID(ctx context.Context, taskID string) (*entity.CollaborationTask, error) {
	var task entity.CollaborationTask
	err := d.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Update 更新任务
func (d *CollaborationDAL) Update(ctx context.Context, task *entity.CollaborationTask) error {
	task.UpdatedAt = time.Now().Unix() * 1000

	err := d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Where("task_id = ?", task.TaskID).
		Updates(task).Error
	if err != nil {
		return fmt.Errorf("failed to update collaboration task: %w", err)
	}
	return nil
}

// UpdateStatus 更新任务状态
func (d *CollaborationDAL) UpdateStatus(ctx context.Context, taskID string, status entity.TaskStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Unix() * 1000,
	}

	if status == entity.TaskStatusCompleted {
		updates["completed_at"] = time.Now().Unix() * 1000
	}

	err := d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Where("task_id = ?", taskID).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}
	return nil
}

// AssignTask 分配任务
func (d *CollaborationDAL) AssignTask(ctx context.Context, taskID, assignedTo string) error {
	updates := map[string]interface{}{
		"assigned_to": assignedTo,
		"status":      entity.TaskStatusAssigned,
		"updated_at":  time.Now().Unix() * 1000,
	}

	err := d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Where("task_id = ?", taskID).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("failed to assign task: %w", err)
	}
	return nil
}

// ==================== 查询操作 ====================

// List 获取任务列表
func (d *CollaborationDAL) List(ctx context.Context, filter *TaskFilter) ([]*entity.CollaborationTask, int64, error) {
	var tasks []*entity.CollaborationTask
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.CollaborationTask{})

	// 应用过滤条件
	query = d.applyFilter(query, filter)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// 应用分页和排序
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 应用排序
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// 执行查询
	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, total, nil
}

// GetPendingTasks 获取待分配任务
func (d *CollaborationDAL) GetPendingTasks(ctx context.Context, tenantID string) ([]*entity.CollaborationTask, error) {
	var tasks []*entity.CollaborationTask
	err := d.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, entity.TaskStatusPending).
		Order("priority DESC, created_at ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pending tasks: %w", err)
	}
	return tasks, nil
}

// GetTasksByAssignee 获取分配给指定用户的任务
func (d *CollaborationDAL) GetTasksByAssignee(ctx context.Context, tenantID, assigneeID string) ([]*entity.CollaborationTask, error) {
	var tasks []*entity.CollaborationTask
	err := d.db.WithContext(ctx).
		Where("tenant_id = ? AND assigned_to = ?", tenantID, assigneeID).
		Order("priority DESC, sla_deadline ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by assignee: %w", err)
	}
	return tasks, nil
}

// GetOverdueTasks 获取超时任务
func (d *CollaborationDAL) GetOverdueTasks(ctx context.Context) ([]*entity.CollaborationTask, error) {
	var tasks []*entity.CollaborationTask
	now := time.Now()
	err := d.db.WithContext(ctx).
		Where("sla_deadline < ? AND status IN (?)", now, []entity.TaskStatus{
			entity.TaskStatusPending,
			entity.TaskStatusAssigned,
			entity.TaskStatusInProgress,
		}).
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}
	return tasks, nil
}

// ==================== 统计操作 ====================

// GetStats 获取任务统计
func (d *CollaborationDAL) GetStats(ctx context.Context, tenantID string) (*TaskStats, error) {
	stats := &TaskStats{}

	// 总任务数
	if err := d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Where("tenant_id = ?", tenantID).
		Count(&stats.TotalTasks).Error; err != nil {
		return nil, fmt.Errorf("failed to count total tasks: %w", err)
	}

	// 各状态统计
	statusStats := []struct {
		Status entity.TaskStatus
		Count  int64
	}{}

	err := d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Select("status, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Scan(&statusStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count tasks by status: %w", err)
	}

	stats.ByStatus = make(map[string]int64)
	for _, stat := range statusStats {
		stats.ByStatus[string(stat.Status)] = stat.Count
		switch stat.Status {
		case entity.TaskStatusPending:
			stats.PendingTasks = stat.Count
		case entity.TaskStatusInProgress:
			stats.InProgressTasks = stat.Count
		case entity.TaskStatusCompleted:
			stats.CompletedTasks = stat.Count
		}
	}

	// 优先级统计
	priorityStats := []struct {
		Priority entity.TaskPriority
		Count    int64
	}{}

	err = d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Select("priority, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("priority").
		Scan(&priorityStats).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count tasks by priority: %w", err)
	}

	stats.ByPriority = make(map[string]int64)
	for _, stat := range priorityStats {
		stats.ByPriority[string(stat.Priority)] = stat.Count
	}

	// 超时任务统计
	err = d.db.WithContext(ctx).
		Model(&entity.CollaborationTask{}).
		Where("tenant_id = ? AND sla_deadline < ? AND status NOT IN (?)",
			tenantID, time.Now(), []entity.TaskStatus{entity.TaskStatusCompleted, entity.TaskStatusCancelled}).
		Count(&stats.OverdueTasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count overdue tasks: %w", err)
	}

	return stats, nil
}

// ==================== 辅助方法 ====================

// applyFilter 应用查询过滤条件
func (d *CollaborationDAL) applyFilter(query *gorm.DB, filter *TaskFilter) *gorm.DB {
	if filter == nil {
		return query
	}

	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	if filter.TaskType != nil {
		query = query.Where("task_type = ?", *filter.TaskType)
	}
	if filter.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedTo)
	}
	if filter.IsOverdue != nil {
		if *filter.IsOverdue {
			query = query.Where("sla_deadline < ?", time.Now())
		} else {
			query = query.Where("sla_deadline >= ?", time.Now())
		}
	}

	return query
}

// ==================== 数据结构 ====================

// TaskFilter 任务查询过滤器
type TaskFilter struct {
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
