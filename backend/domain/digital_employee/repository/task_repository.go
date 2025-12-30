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

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
)

// TaskAssignmentRepository 任务分配仓储接口
type TaskAssignmentRepository interface {
	// Create 创建任务分配
	Create(ctx context.Context, assignment *entity.TaskAssignment) error

	// Update 更新任务分配
	Update(ctx context.Context, assignment *entity.TaskAssignment) error

	// GetByID 根据ID获取任务分配
	GetByID(ctx context.Context, assignmentID string) (*entity.TaskAssignment, error)

	// GetByEmployeeID 根据员工ID获取任务列表
	GetByEmployeeID(ctx context.Context, employeeID string, status *entity.TaskStatus, page, pageSize int) ([]*entity.TaskAssignment, int64, error)

	// GetByTaskID 根据任务ID获取分配记录
	GetByTaskID(ctx context.Context, taskID string) (*entity.TaskAssignment, error)

	// GetActiveAssignmentsByTenant 获取租户下的活跃任务
	GetActiveAssignmentsByTenant(ctx context.Context, tenantID string, page, pageSize int) ([]*entity.TaskAssignment, int64, error)

	// CountByEmployeeID 统计员工任务数量（按状态）
	CountByEmployeeID(ctx context.Context, employeeID string) (total int64, completed int64, failed int64, err error)

	// GetAvailableEmployees 获取可用员工（根据技能匹配）
	GetAvailableEmployees(ctx context.Context, tenantID string, taskType entity.TaskType, requiredSkills []string, limit int) ([]*entity.EmployeeProfile, error)
}
