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

// organizationRepository 组织仓储实现
type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository 创建组织仓储实例
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// Create 创建组织
func (r *organizationRepository) Create(ctx context.Context, org *entity.Organization) error {
	org.CreatedAt = time.Now().UnixMilli()
	org.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(org).Error
}

// GetByID 根据ID获取组织
func (r *organizationRepository) GetByID(ctx context.Context, orgID string) (*entity.Organization, error) {
	var org entity.Organization
	err := r.db.WithContext(ctx).
		Preload("Leader").
		Where("org_id = ?", orgID).
		Where("deleted_at IS NULL").
		First(&org).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetByCode 根据组织编码获取组织
func (r *organizationRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Organization, error) {
	var org entity.Organization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("org_code = ?", code).
		Where("deleted_at IS NULL").
		First(&org).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetByTenantID 获取租户下的所有组织
func (r *organizationRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&orgs).Error
	return orgs, err
}

// GetTree 获取组织树
func (r *organizationRepository) GetTree(ctx context.Context, tenantID string) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&orgs).Error
	return orgs, err
}

// GetChildren 获取子组织
func (r *organizationRepository) GetChildren(ctx context.Context, parentID string) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Where("deleted_at IS NULL").
		Order("sort_order ASC").
		Find(&orgs).Error
	return orgs, err
}

// Update 更新组织
func (r *organizationRepository) Update(ctx context.Context, org *entity.Organization) error {
	org.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("org_id = ?", org.OrgID).
		Updates(org).Error
}

// Delete 软删除组织
func (r *organizationRepository) Delete(ctx context.Context, orgID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("org_id = ?", orgID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// ExistsByCode 检查编码是否存在
func (r *organizationRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("tenant_id = ?", tenantID).
		Where("org_code = ?", code).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("org_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// List 分页查询组织列表
func (r *organizationRepository) List(ctx context.Context, filter *OrganizationFilter) ([]*entity.Organization, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Organization{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.OrgType != "" {
		query = query.Where("org_type = ?", filter.OrgType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ParentID != nil {
		if *filter.ParentID == "" {
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", *filter.ParentID)
		}
	}
	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.Keyword != "" {
		query = query.Where("org_name LIKE ? OR org_code LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
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
	query = query.Order("level ASC, sort_order ASC")

	// 查询
	var orgs []*entity.Organization
	err := query.Find(&orgs).Error
	return orgs, total, err
}

// departmentRepository 部门仓储实现
type departmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository 创建部门仓储实例
func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

// Create 创建部门
func (r *departmentRepository) Create(ctx context.Context, dept *entity.Department) error {
	dept.CreatedAt = time.Now().UnixMilli()
	dept.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(dept).Error
}

// GetByID 根据ID获取部门
func (r *departmentRepository) GetByID(ctx context.Context, deptID string) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).
		Preload("Organization").
		Preload("Leader").
		Where("dept_id = ?", deptID).
		Where("deleted_at IS NULL").
		First(&dept).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// GetByCode 根据部门编码获取部门
func (r *departmentRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("dept_code = ?", code).
		Where("deleted_at IS NULL").
		First(&dept).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// GetByTenantID 获取租户下的所有部门
func (r *departmentRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&depts).Error
	return depts, err
}

// GetTree 获取部门树
func (r *departmentRepository) GetTree(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("level ASC, sort_order ASC").
		Find(&depts).Error
	return depts, err
}

// GetChildren 获取子部门
func (r *departmentRepository) GetChildren(ctx context.Context, parentID string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Where("deleted_at IS NULL").
		Order("sort_order ASC").
		Find(&depts).Error
	return depts, err
}

// Update 更新部门
func (r *departmentRepository) Update(ctx context.Context, dept *entity.Department) error {
	dept.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("dept_id = ?", dept.DeptID).
		Updates(dept).Error
}

// Delete 软删除部门
func (r *departmentRepository) Delete(ctx context.Context, deptID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("dept_id = ?", deptID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// Move 移动部门到新的父部门下
func (r *departmentRepository) Move(ctx context.Context, deptID, newParentID string) error {
	// 获取部门信息
	dept, err := r.GetByID(ctx, deptID)
	if err != nil {
		return err
	}
	if dept == nil {
		return errors.New("department not found")
	}

	// 计算新的层级和路径
	var newLevel int
	var newPath string
	if newParentID == "" {
		// 移动到根级别
		newLevel = 1
		newPath = "/" + deptID
	} else {
		// 获取父部门信息
		parent, err := r.GetByID(ctx, newParentID)
		if err != nil {
			return err
		}
		if parent == nil {
			return errors.New("parent department not found")
		}
		newLevel = parent.Level + 1
		newPath = parent.Path + "/" + deptID
	}

	// 更新部门
	dept.ParentID = &newParentID
	dept.Level = newLevel
	dept.Path = newPath
	dept.UpdatedAt = time.Now().UnixMilli()

	return r.Update(ctx, dept)
}

// ExistsByCode 检查编码是否存在
func (r *departmentRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("tenant_id = ?", tenantID).
		Where("dept_code = ?", code).
		Where("deleted_at IS NULL")

	if excludeID != "" {
		query = query.Where("dept_id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// List 分页查询部门列表
func (r *departmentRepository) List(ctx context.Context, filter *DepartmentFilter) ([]*entity.Department, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Department{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.OrgID != "" {
		query = query.Where("org_id = ?", filter.OrgID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ParentID != nil {
		if *filter.ParentID == "" {
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", *filter.ParentID)
		}
	}
	if filter.Level != nil {
		query = query.Where("level = ?", *filter.Level)
	}
	if filter.Keyword != "" {
		query = query.Where("dept_name LIKE ? OR dept_code LIKE ?", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")
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
	query = query.Order("level ASC, sort_order ASC")

	// 查询
	var depts []*entity.Department
	err := query.Find(&depts).Error
	return depts, total, err
}

// parseInt 解析整数字符串
func parseInt(s string) int {
	result := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}
