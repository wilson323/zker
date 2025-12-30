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
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/repository"
)

// ProtocolService 交互协议服务接口
type ProtocolService interface {
	// BuildTaskContext 构建任务上下文
	BuildTaskContext(ctx context.Context, task *entity.CollaborationTask) (*TaskContext, error)

	// FormatResult 格式化任务结果
	FormatResult(ctx context.Context, task *entity.CollaborationTask) (*FormattedResult, error)

	// RecordHistory 记录协同历史
	RecordHistory(ctx context.Context, history *entity.CollaborationHistory) error

	// GetHistory 获取协同历史
	GetHistory(ctx context.Context, taskID string) ([]*entity.CollaborationHistory, error)

	// GetHistoryWithPagination 分页获取历史
	GetHistoryWithPagination(
		ctx context.Context,
		taskID string,
		page, pageSize int,
	) ([]*entity.CollaborationHistory, int64, error)

	// CalculateMetrics 计算性能指标
	CalculateMetrics(ctx context.Context, taskID string) (*PerformanceMetrics, error)
}

// TaskContext 任务上下文
type TaskContext struct {
	TaskID          string                 `json:"task_id"`
	TaskType        entity.TaskType        `json:"task_type"`
	Source          entity.TaskSource      `json:"source"`
	Priority        entity.TaskPriority    `json:"priority"`
	Status          entity.TaskStatus      `json:"status"`
	Context         map[string]interface{} `json:"context"`
	CreatedAt       int64                  `json:"created_at"`
	SLADeadline     int64                  `json:"sla_deadline"`
	RemainingSLA    int64                  `json:"remaining_sla"`
	AssignedTo      *string                `json:"assigned_to,omitempty"`
}

// FormattedResult 格式化结果
type FormattedResult struct {
	TaskID          string                 `json:"task_id"`
	Decision        string                 `json:"decision"`
	HumanDecision   map[string]interface{} `json:"human_decision"`
	DecisionReason  string                 `json:"decision_reason"`
	FormattedOutput string                 `json:"formatted_output"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	TaskID              string  `json:"task_id"`
	WaitTime            float64 `json:"wait_time"`             // 等待时间（秒）
	HandleTime          float64 `json:"handle_time"`           // 处理时间（秒）
	TotalTime           float64 `json:"total_time"`            // 总时间（秒）
	SLACompliance       bool    `json:"sla_compliance"`        // 是否符合SLA
	ActionsCount        int     `json:"actions_count"`         // 操作次数
	EscalationCount      int     `json:"escalation_count"`      // 升级次数
	ReviewerChangedCount int     `json:"reviewer_changed_count"` // 审核人变更次数
}

// protocolService 交互协议服务实现
type protocolService struct {
	taskRepo    repository.CollaborationRepository
	historyRepo repository.HistoryRepository
}

// NewProtocolService 创建交互协议服务
func NewProtocolService(
	taskRepo repository.CollaborationRepository,
	historyRepo repository.HistoryRepository,
) ProtocolService {
	return &protocolService{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
	}
}

// BuildTaskContext 构建任务上下文
func (s *protocolService) BuildTaskContext(
	ctx context.Context,
	task *entity.CollaborationTask,
) (*TaskContext, error) {
	// 1. 解析任务上下文
	var contextData map[string]interface{}
	if err := json.Unmarshal(task.Context, &contextData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal context: %w", err)
	}

	// 2. 构建任务上下文
	return &TaskContext{
		TaskID:       task.TaskID,
		TaskType:     task.TaskType,
		Source:       task.Source,
		Priority:     task.Priority,
		Status:       task.Status,
		Context:      contextData,
		CreatedAt:    task.CreatedAt,
		SLADeadline:  task.SLADeadline.Unix() * 1000,
		RemainingSLA: task.GetRemainingSLA(),
		AssignedTo:   task.AssignedTo,
	}, nil
}

// FormatResult 格式化任务结果
func (s *protocolService) FormatResult(
	ctx context.Context,
	task *entity.CollaborationTask,
) (*FormattedResult, error) {
	// 1. 解析结果
	var resultData map[string]interface{}
	if len(task.Result) > 0 {
		if err := json.Unmarshal(task.Result, &resultData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	// 2. 格式化输出
	decision := "pending"
	if task.Status == entity.TaskStatusCompleted {
		decision = "approved"
	} else if task.Status == entity.TaskStatusCancelled {
		decision = "rejected"
	}

	return &FormattedResult{
		TaskID:          task.TaskID,
		Decision:        decision,
		HumanDecision:   resultData,
		DecisionReason:  "",
		FormattedOutput: fmt.Sprintf("Task %s: %s", task.TaskID, decision),
	}, nil
}

// RecordHistory 记录协同历史
func (s *protocolService) RecordHistory(ctx context.Context, history *entity.CollaborationHistory) error {
	return s.historyRepo.Create(ctx, history)
}

// GetHistory 获取协同历史
func (s *protocolService) GetHistory(ctx context.Context, taskID string) ([]*entity.CollaborationHistory, error) {
	return s.historyRepo.GetByTaskID(ctx, taskID)
}

// GetHistoryWithPagination 分页获取历史
func (s *protocolService) GetHistoryWithPagination(
	ctx context.Context,
	taskID string,
	page, pageSize int,
) ([]*entity.CollaborationHistory, int64, error) {
	return s.historyRepo.GetByTaskIDWithPagination(ctx, taskID, page, pageSize)
}

// CalculateMetrics 计算性能指标
func (s *protocolService) CalculateMetrics(
	ctx context.Context,
	taskID string,
) (*PerformanceMetrics, error) {
	// 1. 获取任务
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// 2. 获取历史记录
	histories, err := s.historyRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	// 3. 计算指标
	metrics := &PerformanceMetrics{
		TaskID:        taskID,
		TotalTime:     float64(time.Now().Unix()*1000-task.CreatedAt) / 1000,
		ActionsCount:  len(histories),
	}

	// 4. 分析历史记录
	var assignedAt, completedAt int64
	escalationCount := 0
	reviewerChangedCount := 0

	for _, history := range histories {
		switch history.Action {
		case entity.ActionAssigned:
			if assignedAt == 0 {
				assignedAt = history.Timestamp
			} else {
				reviewerChangedCount++
			}
		case entity.ActionCompleted:
			completedAt = history.Timestamp
		case entity.ActionEscalated:
			escalationCount++
		}
	}

	// 5. 计算等待时间和处理时间
	if assignedAt > 0 {
		metrics.WaitTime = float64(assignedAt-task.CreatedAt) / 1000
	}
	if completedAt > 0 && assignedAt > 0 {
		metrics.HandleTime = float64(completedAt-assignedAt) / 1000
	}

	// 6. 判断SLA符合性
	if completedAt > 0 {
		metrics.SLACompliance = completedAt <= task.SLADeadline.Unix()*1000
	}

	metrics.EscalationCount = escalationCount
	metrics.ReviewerChangedCount = reviewerChangedCount

	return metrics, nil
}
