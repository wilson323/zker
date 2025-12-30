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

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/repository"
)

// CollaborationOrchestrator 协同编排器接口
type CollaborationOrchestrator interface {
	// CreateTask 创建协同任务
	CreateTask(ctx context.Context, req *CreateTaskRequest) (*entity.CollaborationTask, error)

	// AssignTask 分配任务
	AssignTask(ctx context.Context, taskID, assigneeID string, reason string) error

	// AutoAssignTask 自动分配任务
	AutoAssignTask(ctx context.Context, taskID string) error

	// SubmitResult 提交任务结果
	SubmitResult(ctx context.Context, taskID string, result *SubmitResultRequest) error

	// EscalateTask 升级任务
	EscalateTask(ctx context.Context, taskID string, escalateTo string, reason string) error

	// MonitorSLA 监控SLA
	MonitorSLA(ctx context.Context) ([]*OverdueTaskInfo, error)

	// GetTask 获取任务详情
	GetTask(ctx context.Context, taskID string) (*entity.CollaborationTask, error)
}

// orchestrator 协同编排器实现
type orchestrator struct {
	taskRepo    repository.CollaborationRepository
	historyRepo repository.HistoryRepository
	configRepo  repository.ConfigRepository
	assigner    TaskAssigner
	notifier    NotificationService
	logger      *zap.Logger
}

// NewCollaborationOrchestrator 创建协同编排器
func NewCollaborationOrchestrator(
	taskRepo repository.CollaborationRepository,
	historyRepo repository.HistoryRepository,
	configRepo repository.ConfigRepository,
	assigner TaskAssigner,
	notifier NotificationService,
	logger *zap.Logger,
) CollaborationOrchestrator {
	return &orchestrator{
		taskRepo:    taskRepo,
		historyRepo: historyRepo,
		configRepo:  configRepo,
		assigner:    assigner,
		notifier:    notifier,
		logger:      logger,
	}
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	TenantID       string
	TaskType       entity.TaskType
	Source         entity.TaskSource
	Priority       entity.TaskPriority
	ConversationID *string
	MessageID      *string
	BotID          *string
	WorkflowID     *string
	Context        map[string]interface{}
	TriggerReason  string
	AssignTo       *string
	SLAMinutes     int
}

// SubmitResultRequest 提交结果请求
type SubmitResultRequest struct {
	ReviewerID      string
	Decision        string // approved, rejected, modified
	HumanDecision   map[string]interface{}
	DecisionReason  string
	ModifiedContent *string
}

// OverdueTaskInfo 超时任务信息
type OverdueTaskInfo struct {
	TaskID          string              `json:"task_id"`
	TenantID        string              `json:"tenant_id"`
	TaskType        entity.TaskType     `json:"task_type"`
	Priority        entity.TaskPriority `json:"priority"`
	Status          entity.TaskStatus   `json:"status"`
	AssignedTo      *string             `json:"assigned_to,omitempty"`
	OverdueDuration int64               `json:"overdue_duration"` // 超时时长（秒）
}

// ==================== 核心方法实现 ====================

// CreateTask 创建协同任务
func (s *orchestrator) CreateTask(ctx context.Context, req *CreateTaskRequest) (*entity.CollaborationTask, error) {
	// 1. 获取租户配置
	config, err := s.configRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		s.logger.Error("failed to get config",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// 2. 确定SLA截止时间
	slaDeadline := s.determineSLADeadline(req.Priority, req.SLAMinutes, config)

	// 3. 构建任务上下文
	contextJSON, err := json.Marshal(req.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %w", err)
	}

	// 4. 创建任务实体
	task := &entity.CollaborationTask{
		TaskID:       uuid.New().String(),
		TenantID:     req.TenantID,
		TaskType:     req.TaskType,
		Source:       req.Source,
		Priority:     req.Priority,
		Status:       entity.TaskStatusPending,
		ConversationID: req.ConversationID,
		MessageID:    req.MessageID,
		BotID:        req.BotID,
		WorkflowID:   req.WorkflowID,
		Context:      entity.JSON(contextJSON),
		SLADeadline:  slaDeadline,
	}

	// 5. 保存任务
	if err := s.taskRepo.Create(ctx, task); err != nil {
		s.logger.Error("failed to create task",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// 6. 记录历史
	s.recordHistory(ctx, task.TaskID, req.TenantID, entity.SourceSystem, entity.ActionCreated,
		&entity.HistoryDetail{
			Comment: req.TriggerReason,
		})

	// 7. 自动分配（如果指定）
	if req.AssignTo != nil {
		if err := s.AssignTask(ctx, task.TaskID, *req.AssignTo, "创建时指定"); err != nil {
			s.logger.Warn("failed to auto assign task",
				zap.String("task_id", task.TaskID),
				zap.Error(err))
		}
	}

	s.logger.Info("collaboration task created",
		zap.String("task_id", task.TaskID),
		zap.String("tenant_id", req.TenantID),
		zap.String("task_type", string(req.TaskType)))

	return task, nil
}

// AssignTask 分配任务
func (s *orchestrator) AssignTask(ctx context.Context, taskID, assigneeID string, reason string) error {
	// 1. 获取任务
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2. 检查是否可分配
	if !task.CanBeAssigned() {
		return fmt.Errorf("task cannot be assigned in current status: %s", task.Status)
	}

	// 3. 更新任务分配
	if err := s.taskRepo.AssignTask(ctx, taskID, assigneeID); err != nil {
		return fmt.Errorf("failed to assign task: %w", err)
	}

	// 4. 记录历史
	detail := &entity.HistoryDetail{
		NewStatus: &task.Status,
		AssignedTo: &assigneeID,
		Comment:   reason,
	}
	if task.AssignedTo != nil {
		detail.AssignedFrom = task.AssignedTo
	}

	s.recordHistory(ctx, taskID, task.TenantID, entity.SourceHuman, entity.ActionAssigned, detail)

	// 5. 发送通知
	s.notifier.NotifyTaskAssigned(ctx, task, assigneeID)

	s.logger.Info("task assigned",
		zap.String("task_id", taskID),
		zap.String("assignee_id", assigneeID),
		zap.String("reason", reason))

	return nil
}

// AutoAssignTask 自动分配任务
func (s *orchestrator) AutoAssignTask(ctx context.Context, taskID string) error {
	// 1. 获取任务
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2. 选择审核人
	assigneeID, err := s.assigner.SelectAssignee(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to select assignee: %w", err)
	}

	// 3. 分配任务
	return s.AssignTask(ctx, taskID, assigneeID, "自动分配")
}

// SubmitResult 提交任务结果
func (s *orchestrator) SubmitResult(ctx context.Context, taskID string, result *SubmitResultRequest) error {
	// 1. 获取任务
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2. 检查是否可完成
	if !task.CanBeCompleted() {
		return fmt.Errorf("task cannot be completed in current status: %s", task.Status)
	}

	// 3. 确定新状态
	var newStatus entity.TaskStatus
	var action entity.HistoryAction

	switch result.Decision {
	case "approved":
		newStatus = entity.TaskStatusCompleted
		action = entity.ActionCompleted
	case "rejected":
		newStatus = entity.TaskStatusCancelled
		action = entity.ActionRejected
	case "modified":
		newStatus = entity.TaskStatusCompleted
		action = entity.ActionModified
	default:
		return fmt.Errorf("invalid decision: %s", result.Decision)
	}

	// 4. 构建结果JSON
	resultJSON, err := json.Marshal(result.HumanDecision)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	task.Result = entity.JSON(resultJSON)

	// 5. 更新任务状态
	if err := s.taskRepo.UpdateStatus(ctx, taskID, newStatus); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	// 6. 记录历史
	detail := &entity.HistoryDetail{
		OldStatus:   &task.Status,
		NewStatus:   &newStatus,
		Comment:     result.DecisionReason,
	}

	s.recordHistory(ctx, taskID, task.TenantID, entity.SourceHuman, action, detail)

	// 7. 发送通知
	s.notifier.NotifyTaskCompleted(ctx, task, result.Decision)

	s.logger.Info("task completed",
		zap.String("task_id", taskID),
		zap.String("decision", result.Decision),
		zap.String("reviewer_id", result.ReviewerID))

	return nil
}

// EscalateTask 升级任务
func (s *orchestrator) EscalateTask(ctx context.Context, taskID string, escalateTo string, reason string) error {
	// 1. 获取任务
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// 2. 检查是否可升级
	if !task.CanBeEscalated() {
		return fmt.Errorf("task cannot be escalated in current status: %s", task.Status)
	}

	// 3. 更新任务状态和分配
	oldStatus := task.Status
	task.Status = entity.TaskStatusEscalated

	if err := s.taskRepo.AssignTask(ctx, taskID, escalateTo); err != nil {
		return fmt.Errorf("failed to escalate task: %w", err)
	}

	// 4. 记录历史
	s.recordHistory(ctx, taskID, task.TenantID, entity.SourceHuman, entity.ActionEscalated,
		&entity.HistoryDetail{
			OldStatus: &oldStatus,
			NewStatus: &task.Status,
			AssignedTo: &escalateTo,
			Reason:    reason,
		})

	// 5. 发送通知
	s.notifier.NotifyTaskEscalated(ctx, task, escalateTo, reason)

	s.logger.Info("task escalated",
		zap.String("task_id", taskID),
		zap.String("escalate_to", escalateTo),
		zap.String("reason", reason))

	return nil
}

// MonitorSLA 监控SLA
func (s *orchestrator) MonitorSLA(ctx context.Context) ([]*OverdueTaskInfo, error) {
	// 1. 获取超时任务
	overdueTasks, err := s.taskRepo.GetOverdueTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}

	// 2. 构建超时任务信息
	overdueInfos := make([]*OverdueTaskInfo, 0, len(overdueTasks))
	for _, task := range overdueTasks {
		overdueDuration := int64(time.Since(task.SLADeadline).Seconds())

		overdueInfos = append(overdueInfos, &OverdueTaskInfo{
			TaskID:          task.TaskID,
			TenantID:        task.TenantID,
			TaskType:        task.TaskType,
			Priority:        task.Priority,
			Status:          task.Status,
			AssignedTo:      task.AssignedTo,
			OverdueDuration: overdueDuration,
		})

		// 3. 记录SLA预警历史
		s.recordHistory(ctx, task.TaskID, task.TenantID, entity.SourceSystem, entity.ActionSLAWarning,
			&entity.HistoryDetail{
				Comment: fmt.Sprintf("任务超时 %d 秒", overdueDuration),
			})
	}

	s.logger.Warn("found overdue tasks",
		zap.Int("count", len(overdueInfos)))

	return overdueInfos, nil
}

// GetTask 获取任务详情
func (s *orchestrator) GetTask(ctx context.Context, taskID string) (*entity.CollaborationTask, error) {
	return s.taskRepo.GetByID(ctx, taskID)
}

// ==================== 辅助方法 ====================

// determineSLADeadline 确定SLA截止时间
func (s *orchestrator) determineSLADeadline(
	priority entity.TaskPriority,
	slaMinutes int,
	config *entity.CollaborationConfig,
) time.Time {
	// 如果指定了SLA，使用指定的
	if slaMinutes > 0 {
		return time.Now().Add(time.Duration(slaMinutes) * time.Minute)
	}

	// 根据优先级确定SLA
	switch priority {
	case entity.PriorityUrgent:
		return config.GetUrgentSLADeadline()
	default:
		return config.GetDefaultSLADeadline()
	}
}

// recordHistory 记录历史
func (s *orchestrator) recordHistory(
	ctx context.Context,
	taskID, tenantID string,
	actorType entity.TaskSource,
	action entity.HistoryAction,
	detail *entity.HistoryDetail,
) {
	detailJSON, err := json.Marshal(detail)
	if err != nil {
		s.logger.Error("failed to marshal history detail",
			zap.String("task_id", taskID),
			zap.Error(err))
		return
	}

	history := &entity.CollaborationHistory{
		TaskID:    taskID,
		TenantID:  tenantID,
		ActorType: actorType,
		ActorID:   "system",
		Action:    action,
		Details:   entity.JSON(detailJSON),
	}

	if err := s.historyRepo.Create(ctx, history); err != nil {
		s.logger.Error("failed to record history",
			zap.String("task_id", taskID),
			zap.Error(err))
	}
}
