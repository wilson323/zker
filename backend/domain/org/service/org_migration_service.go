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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// OrgMigrationService 组织迁移服务
// 职责：组织数据迁移、增量迁移、全量迁移、迁移进度追踪、迁移回滚
type OrgMigrationService struct {
	orgRepo      repository.OrganizationRepository
	employeeRepo repository.EmployeeRepository
	departmentRepo repository.DepartmentRepository
	db           *gorm.DB
}

// NewOrgMigrationService 创建组织迁移服务实例
func NewOrgMigrationService(
	orgRepo repository.OrganizationRepository,
	employeeRepo repository.EmployeeRepository,
	departmentRepo repository.DepartmentRepository,
	db *gorm.DB,
) *OrgMigrationService {
	return &OrgMigrationService{
		orgRepo:      orgRepo,
		employeeRepo: employeeRepo,
		departmentRepo: departmentRepo,
		db:           db,
	}
}

// StartMigrationRequest 启动迁移请求
type StartMigrationRequest struct {
	TenantID       string                      `json:"tenant_id" binding:"required"`
	SourceOrgID    string                      `json:"source_org_id" binding:"required"`
	TargetTenantID string                      `json:"target_tenant_id" binding:"required"`
	TargetOrgID    *string                     `json:"target_org_id,omitempty"`
	Type           entity.MigrationType        `json:"type" binding:"required,oneof=full incremental"`
	DataTypes      []entity.MigrationDataType  `json:"data_types" binding:"required,min=1"`
	CreatedBy      string                      `json:"created_by" binding:"required"`
}

// StartMigrationResponse 启动迁移响应
type StartMigrationResponse struct {
	TaskID string                     `json:"task_id"`
	Status entity.MigrationTaskStatus `json:"status"`
}

// StartMigration 启动迁移
func (s *OrgMigrationService) StartMigration(
	ctx context.Context,
	req *StartMigrationRequest,
) (*StartMigrationResponse, error) {
	// 1. 验证源组织存在
	sourceOrg, err := s.orgRepo.GetByID(ctx, req.SourceOrgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get source org"),
		)
	}
	if sourceOrg == nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "SourceOrgNotFound"),
		)
	}

	// 2. 验证租户匹配
	if sourceOrg.TenantID != req.TenantID {
		return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
	}

	// 3. 验证目标租户存在（假设租户已经存在）
	// TODO: 添加租户仓储验证

	// 4. 生成迁移任务ID
	taskID := generateUUID()

	// 5. 创建迁移任务
	now := time.Now().UnixMilli()
	task := &entity.OrgMigrationTask{
		TaskID:         taskID,
		TenantID:       req.TenantID,
		SourceOrgID:    req.SourceOrgID,
		TargetTenantID: req.TargetTenantID,
		TargetOrgID:    req.TargetOrgID,
		Type:           req.Type,
		DataTypes:      req.DataTypes,
		Status:         entity.MigrationTaskStatusPending,
		Progress:       0,
		CreatedBy:      req.CreatedBy,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 6. 在事务中创建任务
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建任务记录
		if err := tx.Create(task).Error; err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
				errorx.KV("operation", "create migration task"),
			)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 7. 异步执行迁移
	go s.executeMigration(context.Background(), task)

	return &StartMigrationResponse{
		TaskID: taskID,
		Status: task.Status,
	}, nil
}

// executeMigration 执行迁移
func (s *OrgMigrationService) executeMigration(ctx context.Context, task *entity.OrgMigrationTask) {
	// 1. 更新状态为运行中
	now := time.Now().UnixMilli()
	task.Status = entity.MigrationTaskStatusRunning
	task.StartedAt = &now
	task.UpdatedAt = now

	s.db.Save(task)

	// 2. 根据迁移类型执行
	var err error
	switch task.Type {
	case entity.MigrationTypeFull:
		err = s.executeFullMigration(ctx, task)
	case entity.MigrationTypeIncremental:
		err = s.executeIncrementalMigration(ctx, task)
	default:
		err = fmt.Errorf("unsupported migration type: %s", task.Type)
	}

	// 3. 更新任务状态
	now = time.Now().UnixMilli()
	task.UpdatedAt = now

	if err != nil {
		task.Status = entity.MigrationTaskStatusFailed
		task.ErrorMessages = append(task.ErrorMessages, err.Error())
	} else {
		task.Status = entity.MigrationTaskStatusCompleted
		task.CompletedAt = &now
		task.Progress = 100
	}

	s.db.Save(task)
}

// executeFullMigration 执行全量迁移
func (s *OrgMigrationService) executeFullMigration(
	ctx context.Context,
	task *entity.OrgMigrationTask,
) error {
	// 1. 统计需要迁移的数据量
	totalItems, err := s.countMigrationData(ctx, task)
	if err != nil {
		return fmt.Errorf("count migration data failed: %w", err)
	}
	task.TotalItems = totalItems
	s.db.Save(task)

	// 2. 迁移员工
	if containsDataType(task.DataTypes, entity.MigrationDataEmployees) {
		if err := s.migrateEmployees(ctx, task); err != nil {
			return fmt.Errorf("migrate employees failed: %w", err)
		}
	}

	// 3. 迁移部门
	if containsDataType(task.DataTypes, entity.MigrationDataDepartments) {
		if err := s.migrateDepartments(ctx, task); err != nil {
			return fmt.Errorf("migrate departments failed: %w", err)
		}
	}

	// 4. 迁移角色（TODO: 需要角色仓储）
	if containsDataType(task.DataTypes, entity.MigrationDataRoles) {
		// TODO: 实现角色迁移
	}

	// 5. 迁移权限（TODO: 需要权限仓储）
	if containsDataType(task.DataTypes, entity.MigrationDataPermissions) {
		// TODO: 实现权限迁移
	}

	return nil
}

// executeIncrementalMigration 执行增量迁移
func (s *OrgMigrationService) executeIncrementalMigration(
	ctx context.Context,
	task *entity.OrgMigrationTask,
) error {
	// TODO: 实现增量迁移逻辑
	// 增量迁移只迁移上次迁移后新增或修改的数据
	return s.executeFullMigration(ctx, task)
}

// countMigrationData 统计需要迁移的数据量
func (s *OrgMigrationService) countMigrationData(
	ctx context.Context,
	task *entity.OrgMigrationTask,
) (int, error) {
	total := 0

	// 统计员工数量
	if containsDataType(task.DataTypes, entity.MigrationDataEmployees) {
		employees, err := s.employeeRepo.GetByOrgID(ctx, task.SourceOrgID)
		if err != nil {
			return 0, err
		}
		total += len(employees)
	}

	// 统计部门数量
	if containsDataType(task.DataTypes, entity.MigrationDataDepartments) {
		departments, err := s.departmentRepo.GetByTenantID(ctx, task.TenantID)
		if err != nil {
			return 0, err
		}
		total += len(departments)
	}

	return total, nil
}

// migrateEmployees 迁移员工
func (s *OrgMigrationService) migrateEmployees(
	ctx context.Context,
	task *entity.OrgMigrationTask,
) error {
	// 1. 获取源组织员工
	employees, err := s.employeeRepo.GetByOrgID(ctx, task.SourceOrgID)
	if err != nil {
		return err
	}

	// 2. 逐个迁移员工
	for _, emp := range employees {
		// 处理目标组织ID
		targetOrgID := ""
		if task.TargetOrgID != nil {
			targetOrgID = *task.TargetOrgID
		}

		// 创建新员工记录（属于目标租户）
		newEmp := &entity.Employee{
			EmpID:         generateUUID(),
			TenantID:      task.TargetTenantID,
			OrgID:         targetOrgID, // 如果为空，需要创建目标组织
			DeptID:        emp.DeptID,       // TODO: 需要映射到目标部门ID
			PositionID:    emp.PositionID,   // TODO: 需要映射到目标岗位ID
			EmpName:       emp.EmpName,
			EmpCode:       emp.EmpCode,
			Gender:        emp.Gender,
			Phone:         emp.Phone,
			Email:         emp.Email,
			EmployeeType:  emp.EmployeeType,
			EmployeeStatus: emp.EmployeeStatus,
			JobLevel:      emp.JobLevel,
			JobTitle:      emp.JobTitle,
			HireDate:      emp.HireDate,
			RegularDate:   emp.RegularDate,
			ProbationDays: emp.ProbationDays,
			WorkLocation:  emp.WorkLocation,
			Status:        emp.Status,
			CreatedAt:     time.Now().UnixMilli(),
			UpdatedAt:     time.Now().UnixMilli(),
		}

		// 创建员工
		if err := s.employeeRepo.Create(ctx, newEmp); err != nil {
			task.FailedItems++
			task.ErrorMessages = append(task.ErrorMessages,
				fmt.Sprintf("Failed to migrate employee %s: %v", emp.EmpID, err))
			continue
		}

		task.MigratedItems++
		task.Progress = task.GetProgressPercent()
		s.db.Save(task)

		// 记录日志
		s.logMigration(ctx, task, entity.MigrationDataEmployees, "migrate",
			"employee", emp.EmpID, "success", "")
	}

	return nil
}

// migrateDepartments 迁移部门
func (s *OrgMigrationService) migrateDepartments(
	ctx context.Context,
	task *entity.OrgMigrationTask,
) error {
	// 1. 获取源组织部门
	departments, err := s.departmentRepo.GetByTenantID(ctx, task.TenantID)
	if err != nil {
		return err
	}

	// 2. 逐个迁移部门
	for _, dept := range departments {
		// 处理目标组织ID
		targetOrgID := ""
		if task.TargetOrgID != nil {
			targetOrgID = *task.TargetOrgID
		}

		// 创建新部门记录
		newDept := &entity.Department{
			DeptID:      generateUUID(),
			TenantID:    task.TargetTenantID,
			OrgID:       targetOrgID, // 如果为空，需要创建目标组织
			ParentID:    dept.ParentID,    // TODO: 需要映射到目标父部门ID
			DeptName:    dept.DeptName,
			DeptCode:    dept.DeptCode,
			Level:       dept.Level,
			Path:        dept.Path,
			SortOrder:   dept.SortOrder,
			Status:      dept.Status,
			LeaderID:    dept.LeaderID,     // TODO: 需要映射到目标员工ID
			Description: dept.Description,
			CreatedAt:   time.Now().UnixMilli(),
			UpdatedAt:   time.Now().UnixMilli(),
		}

		// 创建部门
		if err := s.departmentRepo.Create(ctx, newDept); err != nil {
			task.FailedItems++
			task.ErrorMessages = append(task.ErrorMessages,
				fmt.Sprintf("Failed to migrate department %s: %v", dept.DeptID, err))
			continue
		}

		task.MigratedItems++
		task.Progress = task.GetProgressPercent()
		s.db.Save(task)

		// 记录日志
		s.logMigration(ctx, task, entity.MigrationDataDepartments, "migrate",
			"department", dept.DeptID, "success", "")
	}

	return nil
}

// GetMigrationProgress 获取迁移进度
func (s *OrgMigrationService) GetMigrationProgress(
	ctx context.Context,
	taskID string,
) (*entity.OrgMigrationTask, error) {
	var task entity.OrgMigrationTask
	err := s.db.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
				errorx.KV("reason", "TaskNotFound"),
			)
		}
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get migration task"),
		)
	}

	return &task, nil
}

// RollbackMigration 回滚迁移
func (s *OrgMigrationService) RollbackMigration(
	ctx context.Context,
	taskID string,
) error {
	// 1. 获取迁移任务
	var task entity.OrgMigrationTask
	err := s.db.Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorx.New(errno.ErrPermissionInvalidParamCode,
				errorx.KV("reason", "TaskNotFound"),
			)
		}
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get migration task"),
		)
	}

	// 2. 验证任务状态
	if !task.IsCompleted() && !task.IsFailed() {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "InvalidTaskStatus"),
		)
	}

	// 3. 执行回滚（TODO: 实现实际回滚逻辑）
	now := time.Now().UnixMilli()
	task.Status = entity.MigrationTaskStatusRolledBack
	task.RolledBackAt = &now
	task.UpdatedAt = now

	if err := s.db.Save(&task).Error; err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "update migration task"),
		)
	}

	return nil
}

// logMigration 记录迁移日志
func (s *OrgMigrationService) logMigration(
	ctx context.Context,
	task *entity.OrgMigrationTask,
	dataType entity.MigrationDataType,
	action, entityType, entityID, status, message string,
) {
	log := &entity.OrgMigrationLog{
		LogID:      generateUUID(),
		TaskID:     task.TaskID,
		DataType:   dataType,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Status:     status,
		Message:    message,
		CreatedAt:  time.Now().UnixMilli(),
	}

	s.db.Create(log)
}

// 辅助函数
func containsDataType(dataTypes []entity.MigrationDataType, dataType entity.MigrationDataType) bool {
	for _, dt := range dataTypes {
		if dt == dataType {
			return true
		}
	}
	return false
}
