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
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
)

// positionRepository 岗位仓储实现
type positionRepository struct {
	db *gorm.DB
}

// NewPositionRepository 创建岗位仓储实例
func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &positionRepository{db: db}
}

// Create 创建岗位
func (r *positionRepository) Create(ctx context.Context, pos *entity.Position) error {
	pos.CreatedAt = time.Now().UnixMilli()
	pos.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(pos).Error
}

// GetByID 根据ID获取岗位
func (r *positionRepository) GetByID(ctx context.Context, positionID string) (*entity.Position, error) {
	var pos entity.Position
	err := r.db.WithContext(ctx).
		Preload("Department").
		Where("position_id = ?", positionID).
		Where("deleted_at IS NULL").
		First(&pos).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pos, nil
}

// GetByCode 根据岗位编码获取岗位
func (r *positionRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Position, error) {
	var pos entity.Position
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("position_code = ?", code).
		Where("deleted_at IS NULL").
		First(&pos).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &pos, nil
}

// GetByTenantID 获取租户下的所有岗位
func (r *positionRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Position, error) {
	var positions []*entity.Position
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level DESC, sort_order ASC").
		Find(&positions).Error
	return positions, err
}

// GetByDepartmentID 获取部门下的岗位
func (r *positionRepository) GetByDepartmentID(ctx context.Context, deptID string) ([]*entity.Position, error) {
	var positions []*entity.Position
	err := r.db.WithContext(ctx).
		Where("dept_id = ?", deptID).
		Where("deleted_at IS NULL").
		Order("level DESC, sort_order ASC").
		Find(&positions).Error
	return positions, err
}

// Update 更新岗位
func (r *positionRepository) Update(ctx context.Context, pos *entity.Position) error {
	pos.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Position{}).
		Where("position_id = ?", pos.PositionID).
		Updates(pos).Error
}

// Delete 软删除岗位
func (r *positionRepository) Delete(ctx context.Context, positionID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Position{}).
		Where("position_id = ?", positionID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// ExistsByCode 检查编码是否存在
func (r *positionRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Position{}).
		Where("tenant_id = ?", tenantID).
		Where("position_code = ?", code).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("position_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// List 分页查询岗位列表
func (r *positionRepository) List(ctx context.Context, filter *PositionFilter) ([]*entity.Position, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Position{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.DeptID != nil {
		if *filter.DeptID == "" {
			query = query.Where("dept_id IS NULL")
		} else {
			query = query.Where("dept_id = ?", *filter.DeptID)
		}
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.Keyword != "" {
		query = query.Where("position_name LIKE ? OR position_code LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	// 软删除过滤
	query = query.Where("deleted_at IS NULL")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parseInt(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	query = query.Order("level DESC, sort_order ASC")

	// 查询
	var positions []*entity.Position
	err := query.Find(&positions).Error
	return positions, total, err
}

// employeeRepository 员工仓储实现
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository 创建员工仓储实例
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

// Create 创建员工
func (r *employeeRepository) Create(ctx context.Context, emp *entity.Employee) error {
	emp.CreatedAt = time.Now().UnixMilli()
	emp.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(emp).Error
}

// GetByID 根据ID获取员工
func (r *employeeRepository) GetByID(ctx context.Context, empID string) (*entity.Employee, error) {
	var emp entity.Employee
	err := r.db.WithContext(ctx).
		Preload("Organization").
		Preload("Department").
		Preload("Position").
		Preload("DirectLeader").
		Where("emp_id = ?", empID).
		Where("deleted_at IS NULL").
		First(&emp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetByCode 根据工号获取员工
func (r *employeeRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Employee, error) {
	var emp entity.Employee
	err := r.db.WithContext(ctx).
		Preload("Organization").
		Preload("Department").
		Preload("Position").
		Where("tenant_id = ?", tenantID).
		Where("emp_code = ?", code).
		Where("deleted_at IS NULL").
		First(&emp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetByUserID 根据用户ID获取员工
func (r *employeeRepository) GetByUserID(ctx context.Context, userID string) (*entity.Employee, error) {
	var emp entity.Employee
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		First(&emp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetByTenantID 获取租户下的所有员工
func (r *employeeRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Employee, error) {
	var emps []*entity.Employee
	err := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Position").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("hire_date DESC").
		Find(&emps).Error
	return emps, err
}

// GetByDepartmentID 获取部门下的员工
func (r *employeeRepository) GetByDepartmentID(ctx context.Context, deptID string) ([]*entity.Employee, error) {
	var emps []*entity.Employee
	err := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Position").
		Where("dept_id = ?", deptID).
		Where("deleted_at IS NULL").
		Order("hire_date DESC").
		Find(&emps).Error
	return emps, err
}

// GetByOrgID 获取组织下的员工
func (r *employeeRepository) GetByOrgID(ctx context.Context, orgID string) ([]*entity.Employee, error) {
	var emps []*entity.Employee
	err := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Position").
		Where("org_id = ?", orgID).
		Where("deleted_at IS NULL").
		Order("hire_date DESC").
		Find(&emps).Error
	return emps, err
}

// Update 更新员工
func (r *employeeRepository) Update(ctx context.Context, emp *entity.Employee) error {
	emp.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Where("emp_id = ?", emp.EmpID).
		Updates(emp).Error
}

// UpdateStatus 更新员工状态
func (r *employeeRepository) UpdateStatus(ctx context.Context, empID string, status entity.EmployeeStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Where("emp_id = ?", empID).
		Updates(map[string]interface{}{
			"employee_status": status,
			"status":          status,
			"updated_at":      time.Now().UnixMilli(),
		}).Error
}

// Delete 软删除员工
func (r *employeeRepository) Delete(ctx context.Context, empID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Where("emp_id = ?", empID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// ExistsByCode 检查工号是否存在
func (r *employeeRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Where("tenant_id = ?", tenantID).
		Where("emp_code = ?", code).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("emp_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *employeeRepository) ExistsByEmail(ctx context.Context, tenantID, email string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Where("tenant_id = ?", tenantID).
		Where("email = ?", email).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("emp_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// List 分页查询员工列表
func (r *employeeRepository) List(ctx context.Context, filter *EmployeeFilter) ([]*entity.Employee, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Employee{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.OrgID != "" {
		query = query.Where("org_id = ?", filter.OrgID)
	}
	if filter.DeptID != nil {
		if *filter.DeptID == "" {
			query = query.Where("dept_id IS NULL")
		} else {
			query = query.Where("dept_id = ?", *filter.DeptID)
		}
	}
	if filter.PositionID != nil {
		query = query.Where("position_id = ?", *filter.PositionID)
	}
	if filter.EmployeeType != "" {
		query = query.Where("employee_type = ?", filter.EmployeeType)
	}
	if filter.EmployeeStatus != "" {
		query = query.Where("employee_status = ?", filter.EmployeeStatus)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.JobLevel != nil {
		query = query.Where("job_level = ?", *filter.JobLevel)
	}
	if filter.Keyword != "" {
		query = query.Where("emp_name LIKE ? OR emp_code LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
	}

	// 软删除过滤
	query = query.Where("deleted_at IS NULL")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parseInt(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	query = query.Order("hire_date DESC")

	// 查询
	var emps []*entity.Employee
	err := query.Preload("Department").Preload("Position").Find(&emps).Error
	return emps, total, err
}

// Search 搜索员工（按姓名、工号、手机号、邮箱）
func (r *employeeRepository) Search(ctx context.Context, tenantID, keyword string, limit int) ([]*entity.Employee, error) {
	var emps []*entity.Employee
	query := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Position").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Where("emp_name LIKE ? OR emp_code LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("hire_date DESC").Find(&emps).Error
	return emps, err
}

// GetByPinyin 按拼音首字母查询员工
func (r *employeeRepository) GetByPinyin(ctx context.Context, tenantID, pinyin string) ([]*entity.Employee, error) {
	// TODO: 实现拼音搜索功能
	// 可以使用 pinyin 库将 emp_name 转换为拼音，然后进行匹配
	// 或者在数据库中增加一个拼音首字母字段
	var emps []*entity.Employee
	err := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Position").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("emp_name ASC").
		Find(&emps).Error
	return emps, err
}
