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
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/repository"
)

// collaborationRepositoryImpl 协同任务仓储实现
type collaborationRepositoryImpl struct {
	dal *CollaborationDAL
}

// NewCollaborationRepository 创建协同任务仓储
func NewCollaborationRepository(db *gorm.DB) repository.CollaborationRepository {
	return &collaborationRepositoryImpl{
		dal: NewCollaborationDAL(db),
	}
}

// Create 创建协同任务
func (r *collaborationRepositoryImpl) Create(ctx context.Context, task *entity.CollaborationTask) error {
	return r.dal.Create(ctx, task)
}

// GetByID 根据ID获取任务
func (r *collaborationRepositoryImpl) GetByID(ctx context.Context, taskID string) (*entity.CollaborationTask, error) {
	return r.dal.GetByID(ctx, taskID)
}

// Update 更新任务
func (r *collaborationRepositoryImpl) Update(ctx context.Context, task *entity.CollaborationTask) error {
	return r.dal.Update(ctx, task)
}

// UpdateStatus 更新任务状态
func (r *collaborationRepositoryImpl) UpdateStatus(ctx context.Context, taskID string, status entity.TaskStatus) error {
	return r.dal.UpdateStatus(ctx, taskID, status)
}

// AssignTask 分配任务
func (r *collaborationRepositoryImpl) AssignTask(ctx context.Context, taskID, assignedTo string) error {
	return r.dal.AssignTask(ctx, taskID, assignedTo)
}

// List 获取任务列表
func (r *collaborationRepositoryImpl) List(ctx context.Context, filter *repository.ListFilter) ([]*entity.CollaborationTask, int64, error) {
	dalFilter := &TaskFilter{
		TenantID:   filter.TenantID,
		Status:     filter.Status,
		Priority:   filter.Priority,
		TaskType:   filter.TaskType,
		AssignedTo: filter.AssignedTo,
		IsOverdue:  filter.IsOverdue,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		SortBy:     filter.SortBy,
		SortOrder:  filter.SortOrder,
	}
	return r.dal.List(ctx, dalFilter)
}

// GetPendingTasks 获取待分配任务
func (r *collaborationRepositoryImpl) GetPendingTasks(ctx context.Context, tenantID string) ([]*entity.CollaborationTask, error) {
	return r.dal.GetPendingTasks(ctx, tenantID)
}

// GetTasksByAssignee 获取分配给指定用户的任务
func (r *collaborationRepositoryImpl) GetTasksByAssignee(ctx context.Context, tenantID, assigneeID string) ([]*entity.CollaborationTask, error) {
	return r.dal.GetTasksByAssignee(ctx, tenantID, assigneeID)
}

// GetOverdueTasks 获取超时任务
func (r *collaborationRepositoryImpl) GetOverdueTasks(ctx context.Context) ([]*entity.CollaborationTask, error) {
	return r.dal.GetOverdueTasks(ctx)
}

// GetStats 获取任务统计
func (r *collaborationRepositoryImpl) GetStats(ctx context.Context, tenantID string) (*repository.TaskStats, error) {
	stats, err := r.dal.GetStats(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return &repository.TaskStats{
		TotalTasks:      stats.TotalTasks,
		PendingTasks:    stats.PendingTasks,
		InProgressTasks: stats.InProgressTasks,
		CompletedTasks:  stats.CompletedTasks,
		OverdueTasks:    stats.OverdueTasks,
		ByPriority:      stats.ByPriority,
		ByStatus:        stats.ByStatus,
	}, nil
}

// historyRepositoryImpl 历史记录仓储实现
type historyRepositoryImpl struct {
	dal *HistoryDAL
}

// NewHistoryRepository 创建历史记录仓储
func NewHistoryRepository(db *gorm.DB) repository.HistoryRepository {
	return &historyRepositoryImpl{
		dal: NewHistoryDAL(db),
	}
}

// Create 创建历史记录
func (r *historyRepositoryImpl) Create(ctx context.Context, history *entity.CollaborationHistory) error {
	return r.dal.Create(ctx, history)
}

// GetByTaskID 获取任务的历史记录
func (r *historyRepositoryImpl) GetByTaskID(ctx context.Context, taskID string) ([]*entity.CollaborationHistory, error) {
	return r.dal.GetByTaskID(ctx, taskID)
}

// GetByTaskIDWithPagination 分页获取任务历史
func (r *historyRepositoryImpl) GetByTaskIDWithPagination(
	ctx context.Context,
	taskID string,
	page, pageSize int,
) ([]*entity.CollaborationHistory, int64, error) {
	return r.dal.GetByTaskIDWithPagination(ctx, taskID, page, pageSize)
}

// GetByAction 获取指定操作类型的历史记录
func (r *historyRepositoryImpl) GetByAction(
	ctx context.Context,
	taskID string,
	action entity.HistoryAction,
) ([]*entity.CollaborationHistory, error) {
	return r.dal.GetByAction(ctx, taskID, action)
}

// DeleteOldHistory 删除旧的历史记录
func (r *historyRepositoryImpl) DeleteOldHistory(ctx context.Context, beforeDate int64) (int64, error) {
	return r.dal.DeleteOldHistory(ctx, time.Unix(beforeDate/1000, 0))
}

// configRepositoryImpl 配置仓储实现
type configRepositoryImpl struct {
	dal *ConfigDAL
}

// NewConfigRepository 创建配置仓储
func NewConfigRepository(db *gorm.DB) repository.ConfigRepository {
	return &configRepositoryImpl{
		dal: NewConfigDAL(db),
	}
}

// GetByTenantID 获取租户配置
func (r *configRepositoryImpl) GetByTenantID(ctx context.Context, tenantID string) (*entity.CollaborationConfig, error) {
	return r.dal.GetByTenantID(ctx, tenantID)
}

// Create 创建配置
func (r *configRepositoryImpl) Create(ctx context.Context, config *entity.CollaborationConfig) error {
	return r.dal.Create(ctx, config)
}

// Update 更新配置
func (r *configRepositoryImpl) Update(ctx context.Context, config *entity.CollaborationConfig) error {
	return r.dal.Update(ctx, config)
}

// Upsert 创建或更新配置
func (r *configRepositoryImpl) Upsert(ctx context.Context, config *entity.CollaborationConfig) error {
	return r.dal.Upsert(ctx, config)
}
