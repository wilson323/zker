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
	"time"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/repository"
)

//go:generate mockgen -destination ../../../internal/mock/domain/digital_employee/service/task_service_mock.go --package mockTaskService -source task_service.go
type TaskService interface {
	// AssignTask 分配任务给员工
	AssignTask(ctx context.Context, req *entity.AssignTaskRequest) (*entity.TaskAssignment, error)

	// AutoAssignTask 自动分配任务（基于技能匹配）
	AutoAssignTask(ctx context.Context, req *entity.AutoAssignTaskRequest) (*entity.TaskAssignment, error)

	// CompleteTask 完成任务
	CompleteTask(ctx context.Context, req *entity.CompleteTaskRequest) error

	// FailTask 标记任务失败
	FailTask(ctx context.Context, req *entity.FailTaskRequest) error

	// GetAssignments 获取员工任务列表
	GetAssignments(ctx context.Context, req *entity.GetAssignmentsRequest) (*entity.GetAssignmentsResponse, error)
}

type taskServiceImpl struct {
	taskRepo   repository.TaskAssignmentRepository
	profileRepo repository.ProfileRepository
}

// NewTaskService 创建任务服务
func NewTaskService(taskRepo repository.TaskAssignmentRepository, profileRepo repository.ProfileRepository) TaskService {
	return &taskServiceImpl{
		taskRepo:   taskRepo,
		profileRepo: profileRepo,
	}
}

// AssignTask 分配任务给员工
func (s *taskServiceImpl) AssignTask(ctx context.Context, req *entity.AssignTaskRequest) (*entity.TaskAssignment, error) {
	// 验证员工是否存在
	_, err := s.profileRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	// 创建任务分配
	assignment := &entity.TaskAssignment{
		AssignmentID: uuid.New().String(),
		TaskID:       req.TaskID,
		EmployeeID:   req.EmployeeID,
		TaskType:     req.TaskType,
		Priority:     req.Priority,
		Status:       entity.TaskStatusAssigned,
	}

	if err := s.taskRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return assignment, nil
}

// AutoAssignTask 自动分配任务（基于技能匹配）
func (s *taskServiceImpl) AutoAssignTask(ctx context.Context, req *entity.AutoAssignTaskRequest) (*entity.TaskAssignment, error) {
	// TODO: 从上下文获取tenantID
	// 这里需要从session或其他地方获取租户ID
	// 为了示例，我们假设从某个地方获取
	tenantID := "" // 需要从上下文获取

	// 查找匹配的员工
	employees, err := s.taskRepo.GetAvailableEmployees(ctx, tenantID, req.TaskType, req.RequiredSkills, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to find available employees: %w", err)
	}

	if len(employees) == 0 {
		return nil, fmt.Errorf("no available employees found for the task")
	}

	// 选择第一个匹配的员工（后续可以根据负载均衡算法优化）
	employee := employees[0]

	// 创建任务分配
	assignment := &entity.TaskAssignment{
		AssignmentID: uuid.New().String(),
		TaskID:       req.TaskID,
		EmployeeID:   employee.EmployeeID,
		TaskType:     req.TaskType,
		Priority:     req.Priority,
		Status:       entity.TaskStatusAssigned,
	}

	if err := s.taskRepo.Create(ctx, assignment); err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return assignment, nil
}

// CompleteTask 完成任务
func (s *taskServiceImpl) CompleteTask(ctx context.Context, req *entity.CompleteTaskRequest) error {
	// 获取任务分配
	assignment, err := s.taskRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// 更新状态
	now := time.Now().UnixMilli()
	assignment.Status = entity.TaskStatusCompleted
	assignment.CompletedAt = &now
	assignment.Result = req.Result

	if err := s.taskRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	return nil
}

// FailTask 标记任务失败
func (s *taskServiceImpl) FailTask(ctx context.Context, req *entity.FailTaskRequest) error {
	// 获取任务分配
	assignment, err := s.taskRepo.GetByID(ctx, req.AssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get assignment: %w", err)
	}

	// 更新状态
	assignment.Status = entity.TaskStatusFailed
	assignment.FailureReason = req.FailureReason

	if err := s.taskRepo.Update(ctx, assignment); err != nil {
		return fmt.Errorf("failed to update assignment: %w", err)
	}

	return nil
}

// GetAssignments 获取员工任务列表
func (s *taskServiceImpl) GetAssignments(ctx context.Context, req *entity.GetAssignmentsRequest) (*entity.GetAssignmentsResponse, error) {
	assignments, total, err := s.taskRepo.GetByEmployeeID(ctx, req.EmployeeID, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %w", err)
	}

	return &entity.GetAssignmentsResponse{
		Assignments: assignments,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}, nil
}
