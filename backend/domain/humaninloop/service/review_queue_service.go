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

package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/repository"
)

// QueueService 审核队列服务接口
type QueueService interface {
	// GetQueue 获取审核队列
	GetQueue(ctx context.Context, tenantID, reviewerID string, filter *QueueFilter) ([]*entity.CollaborationTask, int64, error)

	// SortByPriority 按优先级排序
	SortByPriority(tasks []*entity.CollaborationTask) []*entity.CollaborationTask

	// AdjustPriorityBySLA 根据SLA调整优先级
	AdjustPriorityBySLA(ctx context.Context) error

	// BatchAssign 批量分配任务
	BatchAssign(ctx context.Context, taskIDs, assigneeIDs []string) error

	// GetQueueStats 获取队列统计
	GetQueueStats(ctx context.Context, tenantID string) (*QueueStats, error)
}

// QueueFilter 队列过滤器
type QueueFilter struct {
	Status     *entity.TaskStatus
	Priority   *entity.TaskPriority
	TaskType   *entity.TaskType
	IsOverdue  *bool
	Page       int
	PageSize   int
	SortBy     string
	SortOrder  string
}

// QueueStats 队列统计
type QueueStats struct {
	TotalTasks       int64                 `json:"total_tasks"`
	PendingTasks     int64                 `json:"pending_tasks"`
	InProgressTasks  int64                 `json:"in_progress_tasks"`
	CompletedTasks   int64                 `json:"completed_tasks"`
	OverdueTasks     int64                 `json:"overdue_tasks"`
	ByPriority       map[string]int64      `json:"by_priority"`
	ByStatus         map[string]int64      `json:"by_status"`
	ByType           map[string]int64      `json:"by_type"`
	AverageWaitTime  float64               `json:"average_wait_time"`
	AverageHandleTime float64              `json:"average_handle_time"`
}

// queueService 审核队列服务实现
type queueService struct {
	taskRepo    repository.CollaborationRepository
	historyRepo repository.HistoryRepository
	configRepo  repository.ConfigRepository
}

// NewQueueService 创建审核队列服务
func NewQueueService(
	taskRepo repository.CollaborationRepository,
	historyRepo repository.HistoryRepository,
	configRepo repository.ConfigRepository,
) QueueService {
	return &queueService{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
		configRepo:  configRepo,
	}
}

// GetQueue 获取审核队列
func (s *queueService) GetQueue(
	ctx context.Context,
	tenantID, reviewerID string,
	filter *QueueFilter,
) ([]*entity.CollaborationTask, int64, error) {
	// 构建查询过滤器
	listFilter := &repository.ListFilter{
		TenantID:  tenantID,
		Status:    filter.Status,
		Priority:  filter.Priority,
		TaskType:  filter.TaskType,
		IsOverdue: filter.IsOverdue,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		SortBy:    filter.SortBy,
		SortOrder: filter.SortOrder,
	}

	// 如果指定了审核人，只返回分配给该审核人的任务
	if reviewerID != "" {
		listFilter.AssignedTo = &reviewerID
	}

	// 查询任务列表
	tasks, total, err := s.taskRepo.List(ctx, listFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get queue: %w", err)
	}

	return tasks, total, nil
}

// SortByPriority 按优先级排序
func (s *queueService) SortByPriority(tasks []*entity.CollaborationTask) []*entity.CollaborationTask {
	// 优先级权重
	priorityWeight := map[entity.TaskPriority]int{
		entity.PriorityUrgent:  4,
		entity.PriorityHigh:    3,
		entity.PriorityMedium:  2,
		entity.PriorityLow:     1,
	}

	sortedTasks := make([]*entity.CollaborationTask, len(tasks))
	copy(sortedTasks, tasks)

	sort.SliceStable(sortedTasks, func(i, j int) bool {
		// 1. 优先级高的在前
		if sortedTasks[i].Priority != sortedTasks[j].Priority {
			return priorityWeight[sortedTasks[i].Priority] > priorityWeight[sortedTasks[j].Priority]
		}

		// 2. 优先级相同时，SLA即将到期的在前
		return sortedTasks[i].GetRemainingSLA() < sortedTasks[j].GetRemainingSLA()
	})

	return sortedTasks
}

// AdjustPriorityBySLA 根据SLA调整优先级
func (s *queueService) AdjustPriorityBySLA(ctx context.Context) error {
	// 1. 获取所有进行中的任务
	filter := &repository.ListFilter{
		Status: (*entity.TaskStatus)(&entity.TaskStatusInProgress),
	}

	tasks, _, err := s.taskRepo.List(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	// 2. 检查SLA并调整优先级
	for _, task := range tasks {
		// 计算剩余SLA时间
		remainingSLA := task.GetRemainingSLA()

		// 如果SLA剩余时间少于15分钟，提升到紧急
		if remainingSLA < 900 && task.Priority != entity.PriorityUrgent {
			task.Priority = entity.PriorityUrgent
			if err := s.taskRepo.Update(ctx, task); err != nil {
				continue
			}
		}
		// 如果SLA剩余时间少于1小时，提升到高
		if remainingSLA < 3600 && task.Priority == entity.PriorityLow {
			task.Priority = entity.PriorityHigh
			if err := s.taskRepo.Update(ctx, task); err != nil {
				continue
			}
		}
	}

	return nil
}

// BatchAssign 批量分配任务
func (s *queueService) BatchAssign(ctx context.Context, taskIDs, assigneeIDs []string) error {
	if len(taskIDs) != len(assigneeIDs) {
		return fmt.Errorf("taskIDs and assigneeIDs length mismatch")
	}

	for i := 0; i < len(taskIDs); i++ {
		if err := s.taskRepo.AssignTask(ctx, taskIDs[i], assigneeIDs[i]); err != nil {
			return fmt.Errorf("failed to assign task %s: %w", taskIDs[i], err)
		}
	}

	return nil
}

// GetQueueStats 获取队列统计
func (s *queueService) GetQueueStats(ctx context.Context, tenantID string) (*QueueStats, error) {
	stats, err := s.taskRepo.GetStats(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &QueueStats{
		TotalTasks:       stats.TotalTasks,
		PendingTasks:     stats.PendingTasks,
		InProgressTasks:  stats.InProgressTasks,
		CompletedTasks:   stats.CompletedTasks,
		OverdueTasks:     stats.OverdueTasks,
		ByPriority:       stats.ByPriority,
		ByStatus:         stats.ByStatus,
		ByType:           make(map[string]int64),
		AverageWaitTime:  0,
		AverageHandleTime: 0,
	}, nil
}
