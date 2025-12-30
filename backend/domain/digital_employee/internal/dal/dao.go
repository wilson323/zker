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
	"gorm.io/gorm/clause"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/repository"
)

var _ repository.ProfileRepository = (*profileDAL)(nil)
var _ repository.TaskAssignmentRepository = (*taskAssignmentDAL)(nil)
var _ repository.PerformanceRepository = (*performanceDAL)(nil)

// profileDAL 员工画像DAL实现
type profileDAL struct {
	db *gorm.DB
}

// NewProfileDAL 创建员工画像DAL
func NewProfileDAL(db *gorm.DB) repository.ProfileRepository {
	return &profileDAL{db: db}
}

// Create 创建员工画像
func (d *profileDAL) Create(ctx context.Context, profile *entity.EmployeeProfile) error {
	now := time.Now().UnixMilli()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	if err := d.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("failed to create employee profile: %w", err)
	}
	return nil
}

// Update 更新员工画像
func (d *profileDAL) Update(ctx context.Context, profile *entity.EmployeeProfile) error {
	profile.UpdatedAt = time.Now().UnixMilli()

	if err := d.db.WithContext(ctx).
		Model(&entity.EmployeeProfile{}).
		Where("employee_id = ?", profile.EmployeeID).
		Updates(profile).Error; err != nil {
		return fmt.Errorf("failed to update employee profile: %w", err)
	}
	return nil
}

// GetByID 根据ID获取员工画像
func (d *profileDAL) GetByID(ctx context.Context, employeeID string) (*entity.EmployeeProfile, error) {
	var profile entity.EmployeeProfile
	err := d.db.WithContext(ctx).
		Where("employee_id = ? AND deleted_at IS NULL", employeeID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee profile not found: %s", employeeID)
		}
		return nil, fmt.Errorf("failed to get employee profile: %w", err)
	}
	return &profile, nil
}

// GetByTenantID 根据租户ID列出员工画像
func (d *profileDAL) GetByTenantID(ctx context.Context, tenantID string, role *entity.EmployeeRole, status *entity.EmployeeStatus, page, pageSize int) ([]*entity.EmployeeProfile, int64, error) {
	var profiles []*entity.EmployeeProfile
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.EmployeeProfile{}).Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if role != nil {
		query = query.Where("role = ?", *role)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employee profiles: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&profiles).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list employee profiles: %w", err)
	}

	return profiles, total, nil
}

// Delete 删除员工画像（软删除）
func (d *profileDAL) Delete(ctx context.Context, employeeID string) error {
	now := time.Now().UnixMilli()
	err := d.db.WithContext(ctx).
		Model(&entity.EmployeeProfile{}).
		Where("employee_id = ?", employeeID).
		Updates(map[string]interface{}{
			"status":     entity.EmployeeStatusDeleted,
			"deleted_at": &now,
		}).Error
	if err != nil {
		return fmt.Errorf("failed to delete employee profile: %w", err)
	}
	return nil
}

// GetByBotID 根据Bot ID获取员工画像
func (d *profileDAL) GetByBotID(ctx context.Context, botID string) (*entity.EmployeeProfile, error) {
	var profile entity.EmployeeProfile
	err := d.db.WithContext(ctx).
		Where("bot_id = ? AND deleted_at IS NULL", botID).
		First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee profile not found for bot: %s", botID)
		}
		return nil, fmt.Errorf("failed to get employee profile by bot_id: %w", err)
	}
	return &profile, nil
}

// ExistsByName 检查名称是否存在
func (d *profileDAL) ExistsByName(ctx context.Context, tenantID, name string) (bool, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&entity.EmployeeProfile{}).
		Where("tenant_id = ? AND name = ? AND deleted_at IS NULL", tenantID, name).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check name existence: %w", err)
	}
	return count > 0, nil
}

// taskAssignmentDAL 任务分配DAL实现
type taskAssignmentDAL struct {
	db *gorm.DB
}

// NewTaskAssignmentDAL 创建任务分配DAL
func NewTaskAssignmentDAL(db *gorm.DB) repository.TaskAssignmentRepository {
	return &taskAssignmentDAL{db: db}
}

// Create 创建任务分配
func (d *taskAssignmentDAL) Create(ctx context.Context, assignment *entity.TaskAssignment) error {
	assignment.AssignedAt = time.Now().UnixMilli()

	if err := d.db.WithContext(ctx).Create(assignment).Error; err != nil {
		return fmt.Errorf("failed to create task assignment: %w", err)
	}
	return nil
}

// Update 更新任务分配
func (d *taskAssignmentDAL) Update(ctx context.Context, assignment *entity.TaskAssignment) error {
	if err := d.db.WithContext(ctx).
		Model(&entity.TaskAssignment{}).
		Where("assignment_id = ?", assignment.AssignmentID).
		Updates(assignment).Error; err != nil {
		return fmt.Errorf("failed to update task assignment: %w", err)
	}
	return nil
}

// GetByID 根据ID获取任务分配
func (d *taskAssignmentDAL) GetByID(ctx context.Context, assignmentID string) (*entity.TaskAssignment, error) {
	var assignment entity.TaskAssignment
	err := d.db.WithContext(ctx).
		Where("assignment_id = ?", assignmentID).
		First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task assignment not found: %s", assignmentID)
		}
		return nil, fmt.Errorf("failed to get task assignment: %w", err)
	}
	return &assignment, nil
}

// GetByEmployeeID 根据员工ID获取任务列表
func (d *taskAssignmentDAL) GetByEmployeeID(ctx context.Context, employeeID string, status *entity.TaskStatus, page, pageSize int) ([]*entity.TaskAssignment, int64, error) {
	var assignments []*entity.TaskAssignment
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.TaskAssignment{}).Where("employee_id = ?", employeeID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count task assignments: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("assigned_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&assignments).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list task assignments: %w", err)
	}

	return assignments, total, nil
}

// GetByTaskID 根据任务ID获取分配记录
func (d *taskAssignmentDAL) GetByTaskID(ctx context.Context, taskID string) (*entity.TaskAssignment, error) {
	var assignment entity.TaskAssignment
	err := d.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task assignment not found for task: %s", taskID)
		}
		return nil, fmt.Errorf("failed to get task assignment by task_id: %w", err)
	}
	return &assignment, nil
}

// GetActiveAssignmentsByTenant 获取租户下的活跃任务
func (d *taskAssignmentDAL) GetActiveAssignmentsByTenant(ctx context.Context, tenantID string, page, pageSize int) ([]*entity.TaskAssignment, int64, error) {
	var assignments []*entity.TaskAssignment
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.TaskAssignment{}).
		Where("tenant_id = ? AND status IN ?", tenantID, []entity.TaskStatus{entity.TaskStatusAssigned, entity.TaskStatusInProgress})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count active assignments: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("priority DESC, assigned_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&assignments).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list active assignments: %w", err)
	}

	return assignments, total, nil
}

// CountByEmployeeID 统计员工任务数量
func (d *taskAssignmentDAL) CountByEmployeeID(ctx context.Context, employeeID string) (total, completed, failed int64, err error) {
	var stats struct {
		Total     int64
		Completed int64
		Failed    int64
	}

	err = d.db.WithContext(ctx).
		Model(&entity.TaskAssignment{}).
		Where("employee_id = ?", employeeID).
		Select(
			"COUNT(*) as total",
			"SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as completed",
			"SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as failed",
			entity.TaskStatusCompleted, entity.TaskStatusFailed,
		).
		Scan(&stats).Error
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to count tasks by employee: %w", err)
	}

	return stats.Total, stats.Completed, stats.Failed, nil
}

// GetAvailableEmployees 获取可用员工（根据技能匹配）
func (d *taskAssignmentDAL) GetAvailableEmployees(ctx context.Context, tenantID string, taskType entity.TaskType, requiredSkills []string, limit int) ([]*entity.EmployeeProfile, error) {
	var employees []*entity.EmployeeProfile

	query := d.db.WithContext(ctx).
		Model(&entity.EmployeeProfile{}).
		Where("tenant_id = ? AND role = ? AND status = ? AND deleted_at IS NULL",
			tenantID, taskType, entity.EmployeeStatusActive)

	// 技能匹配（JSON包含查询）
	for _, skill := range requiredSkills {
		query = query.Where("JSON_CONTAINS(skills, ?)", fmt.Sprintf(`"%s"`, skill))
	}

	err := query.
		Order("created_at ASC").
		Limit(limit).
		Find(&employees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get available employees: %w", err)
	}

	return employees, nil
}

// performanceDAL 绩效DAL实现
type performanceDAL struct {
	db *gorm.DB
}

// NewPerformanceDAL 创建绩效DAL
func NewPerformanceDAL(db *gorm.DB) repository.PerformanceRepository {
	return &performanceDAL{db: db}
}

// GetOrCreate 获取或创建绩效记录
func (d *performanceDAL) GetOrCreate(ctx context.Context, employeeID string, tenantID string, period entity.PerformancePeriod, date string) (*entity.EmployeePerformance, error) {
	var performance entity.EmployeePerformance

	err := d.db.WithContext(ctx).
		Where("employee_id = ? AND period = ? AND date = ?", employeeID, period, date).
		First(&performance).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		performance = entity.EmployeePerformance{
			EmployeeID:     employeeID,
			TenantID:       tenantID,
			Period:         period,
			Date:           parseDate(date),
			TotalTasks:     0,
			CompletedTasks: 0,
			FailedTasks:    0,
			CompletionRate: 0.0,
			UpdatedAt:      time.Now().UnixMilli(),
		}
		if err := d.db.WithContext(ctx).Create(&performance).Error; err != nil {
			return nil, fmt.Errorf("failed to create performance record: %w", err)
		}
		return &performance, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get performance record: %w", err)
	}

	return &performance, nil
}

// Update 更新绩效记录
func (d *performanceDAL) Update(ctx context.Context, performance *entity.EmployeePerformance) error {
	performance.UpdatedAt = time.Now().UnixMilli()

	if err := d.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "employee_id"}, {Name: "period"}, {Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"total_tasks", "completed_tasks", "failed_tasks", "completion_rate", "avg_response_time", "customer_rating", "updated_at"}),
		}).
		Create(performance).Error; err != nil {
		return fmt.Errorf("failed to update performance record: %w", err)
	}
	return nil
}

// GetByEmployeeID 获取员工绩效列表
func (d *performanceDAL) GetByEmployeeID(ctx context.Context, employeeID string, period entity.PerformancePeriod, page, pageSize int) ([]*entity.EmployeePerformance, int64, error) {
	var performances []*entity.EmployeePerformance
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.EmployeePerformance{}).Where("employee_id = ?", employeeID)

	if period != "" {
		query = query.Where("period = ?", period)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count performance records: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("date DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&performances).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list performance records: %w", err)
	}

	return performances, total, nil
}

// GetByTenantID 获取租户下所有员工绩效
func (d *performanceDAL) GetByTenantID(ctx context.Context, tenantID string, period entity.PerformancePeriod, page, pageSize int) ([]*entity.EmployeePerformance, int64, error) {
	var performances []*entity.EmployeePerformance
	var total int64

	query := d.db.WithContext(ctx).Model(&entity.EmployeePerformance{}).Where("tenant_id = ?", tenantID)

	if period != "" {
		query = query.Where("period = ?", period)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count team performance records: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("date DESC, completion_rate DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&performances).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list team performance records: %w", err)
	}

	return performances, total, nil
}

// BatchUpdateStats 批量更新统计数据
func (d *performanceDAL) BatchUpdateStats(ctx context.Context) error {
	// TODO: 实现批量更新逻辑（定时任务调用）
	// 1. 查询所有员工
	// 2. 计算每个员工的绩效
	// 3. 更新或创建绩效记录
	return nil
}

// parseDate 解析日期字符串
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}
