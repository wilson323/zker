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

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// TaskAssigner 任务分配器接口
type TaskAssigner interface {
	// SelectAssignee 选择审核人
	SelectAssignee(ctx context.Context, task *entity.CollaborationTask) (string, error)

	// GetAvailableReviewers 获取可用的审核人列表
	GetAvailableReviewers(ctx context.Context, tenantID string) ([]*Reviewer, error)

	// GetReviewerLoad 获取审核人当前负载
	GetReviewerLoad(ctx context.Context, reviewerID string) (*ReviewerLoad, error)
}

// Reviewer 审核人信息
type Reviewer struct {
	ReviewerID      string   `json:"reviewer_id"`
	ReviewerName    string   `json:"reviewer_name"`
	RoleIDs         []string `json:"role_ids,omitempty"`
	Department      string   `json:"department,omitempty"`
	Skills          []string `json:"skills,omitempty"`
	IsOnline        bool     `json:"is_online"`
	MaxTasks        int      `json:"max_tasks"`
	CurrentTasks    int      `json:"current_tasks"`
}

// ReviewerLoad 审核人负载
type ReviewerLoad struct {
	ReviewerID       string  `json:"reviewer_id"`
	PendingTasks     int     `json:"pending_tasks"`
	InProgressTasks  int     `json:"in_progress_tasks"`
	TotalLoad        float64 `json:"total_load"` // 0-1之间，1表示满载
	AverageHandleTime float64 `json:"average_handle_time"` // 平均处理时间（秒）
}

// taskAssigner 任务分配器实现
type taskAssigner struct {
	// 依赖的其他服务
	// userRepo   UserRepository
	// roleRepo   RoleRepository
}

// NewTaskAssigner 创建任务分配器
func NewTaskAssigner() TaskAssigner {
	return &taskAssigner{}
}

// SelectAssignee 选择审核人（基于负载均衡）
func (a *taskAssigner) SelectAssignee(ctx context.Context, task *entity.CollaborationTask) (string, error) {
	// 1. 获取可用审核人
	reviewers, err := a.GetAvailableReviewers(ctx, task.TenantID)
	if err != nil {
		return "", fmt.Errorf("failed to get available reviewers: %w", err)
	}

	if len(reviewers) == 0 {
		return "", fmt.Errorf("no available reviewers for tenant: %s", task.TenantID)
	}

	// 2. 选择负载最低的审核人
	var selectedReviewer *Reviewer
	minLoad := float64(1.0)

	for _, reviewer := range reviewers {
		load, err := a.GetReviewerLoad(ctx, reviewer.ReviewerID)
		if err != nil {
			continue
		}

		if load.TotalLoad < minLoad {
			minLoad = load.TotalLoad
			selectedReviewer = reviewer
		}
	}

	if selectedReviewer == nil {
		// 如果无法获取负载，选择第一个
		selectedReviewer = reviewers[0]
	}

	return selectedReviewer.ReviewerID, nil
}

// GetAvailableReviewers 获取可用的审核人列表
func (a *taskAssigner) GetAvailableReviewers(ctx context.Context, tenantID string) ([]*Reviewer, error) {
	// TODO: 实现从数据库获取审核人列表
	// 这里需要根据租户配置中的审核池过滤规则来查询
	return []*Reviewer{}, nil
}

// GetReviewerLoad 获取审核人当前负载
func (a *taskAssigner) GetReviewerLoad(ctx context.Context, reviewerID string) (*ReviewerLoad, error) {
	// TODO: 实现从数据库获取审核人负载
	// 需要查询该审核人当前的任务数量
	return &ReviewerLoad{
		ReviewerID:       reviewerID,
		PendingTasks:     0,
		InProgressTasks:  0,
		TotalLoad:        0,
		AverageHandleTime: 0,
	}, nil
}
