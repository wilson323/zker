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

	"github.com/coze-studio/backend/domain/permission/entity"
)

// roleRepository 角色仓储实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	role.CreatedAt = time.Now().UnixMilli()
	role.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(role).Error
}

// GetByID 根据ID获取角色
func (r *roleRepository) GetByID(ctx context.Context, roleID string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).
		Preload("DataPerms").
		Preload("FieldPerms").
		Where("role_id = ?", roleID).
		Where("deleted_at IS NULL").
		First(&role).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据租户ID和角色编码获取角色
func (r *roleRepository) GetByCode(ctx context.Context, tenantID string, roleCode string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("role_code = ?", roleCode).
		Where("deleted_at IS NULL").
		First(&role).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update 更新角色
func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	role.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("role_id = ?", role.RoleID).
		Updates(role).Error
}

// Delete 软删除角色
func (r *roleRepository) Delete(ctx context.Context, roleID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Role{}).
		Where("role_id = ?", roleID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// List 分页查询角色列表
func (r *roleRepository) List(ctx context.Context, filter *RoleFilter) ([]*entity.Role, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.Role{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.RoleType != nil {
		query = query.Where("role_type = ?", *filter.RoleType)
	}

	// 软删除过滤
	query = query.Where("deleted_at IS NULL")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 查询数据
	var roles []*entity.Role
	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Find(&roles).Error

	return roles, total, err
}

// GetByUser 获取用户的所有角色
func (r *roleRepository) GetByUser(ctx context.Context, userID, tenantID string) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.tenant_id = ?", tenantID).
		Where("roles.deleted_at IS NULL").
		Find(&roles).Error

	return roles, err
}

// dataPermissionRepository 数据权限仓储实现
type dataPermissionRepository struct {
	db *gorm.DB
}

// NewDataPermissionRepository 创建数据权限仓储实例
func NewDataPermissionRepository(db *gorm.DB) DataPermissionRepository {
	return &dataPermissionRepository{db: db}
}

// Create 创建数据权限
func (r *dataPermissionRepository) Create(ctx context.Context, perm *entity.DataPermission) error {
	perm.CreatedAt = time.Now().UnixMilli()
	perm.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(perm).Error
}

// GetByID 根据ID获取数据权限
func (r *dataPermissionRepository) GetByID(ctx context.Context, permissionID string) (*entity.DataPermission, error) {
	var perm entity.DataPermission
	err := r.db.WithContext(ctx).
		Where("permission_id = ?", permissionID).
		First(&perm).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetByRoleAndResource 根据角色ID和资源类型获取数据权限
func (r *dataPermissionRepository) GetByRoleAndResource(ctx context.Context, roleID string, resourceType entity.ResourceType) (*entity.DataPermission, error) {
	var perm entity.DataPermission
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Where("resource_type = ?", resourceType).
		First(&perm).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetByRole 获取角色的所有数据权限
func (r *dataPermissionRepository) GetByRole(ctx context.Context, roleID string) ([]*entity.DataPermission, error) {
	var perms []*entity.DataPermission
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

// Update 更新数据权限
func (r *dataPermissionRepository) Update(ctx context.Context, perm *entity.DataPermission) error {
	perm.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.DataPermission{}).
		Where("permission_id = ?", perm.PermissionID).
		Updates(perm).Error
}

// Delete 删除数据权限
func (r *dataPermissionRepository) Delete(ctx context.Context, permissionID string) error {
	return r.db.WithContext(ctx).
		Where("permission_id = ?", permissionID).
		Delete(&entity.DataPermission{}).Error
}

// DeleteByRole 删除角色的所有数据权限
func (r *dataPermissionRepository) DeleteByRole(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&entity.DataPermission{}).Error
}

// fieldPermissionRepository 字段权限仓储实现
type fieldPermissionRepository struct {
	db *gorm.DB
}

// NewFieldPermissionRepository 创建字段权限仓储实例
func NewFieldPermissionRepository(db *gorm.DB) FieldPermissionRepository {
	return &fieldPermissionRepository{db: db}
}

// Create 创建字段权限
func (r *fieldPermissionRepository) Create(ctx context.Context, perm *entity.FieldPermission) error {
	perm.CreatedAt = time.Now().UnixMilli()
	perm.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(perm).Error
}

// GetByID 根据ID获取字段权限
func (r *fieldPermissionRepository) GetByID(ctx context.Context, permissionID string) (*entity.FieldPermission, error) {
	var perm entity.FieldPermission
	err := r.db.WithContext(ctx).
		Where("permission_id = ?", permissionID).
		First(&perm).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetByRoleAndResource 根据角色ID和资源类型获取字段权限列表
func (r *fieldPermissionRepository) GetByRoleAndResource(ctx context.Context, roleID string, resourceType string) ([]*entity.FieldPermission, error) {
	var perms []*entity.FieldPermission
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Where("resource_type = ?", resourceType).
		Find(&perms).Error
	return perms, err
}

// GetByRole 获取角色的所有字段权限
func (r *fieldPermissionRepository) GetByRole(ctx context.Context, roleID string) ([]*entity.FieldPermission, error) {
	var perms []*entity.FieldPermission
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Find(&perms).Error
	return perms, err
}

// Update 更新字段权限
func (r *fieldPermissionRepository) Update(ctx context.Context, perm *entity.FieldPermission) error {
	perm.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.FieldPermission{}).
		Where("permission_id = ?", perm.PermissionID).
		Updates(perm).Error
}

// Delete 删除字段权限
func (r *fieldPermissionRepository) Delete(ctx context.Context, permissionID string) error {
	return r.db.WithContext(ctx).
		Where("permission_id = ?", permissionID).
		Delete(&entity.FieldPermission{}).Error
}

// DeleteByRole 删除角色的所有字段权限
func (r *fieldPermissionRepository) DeleteByRole(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&entity.FieldPermission{}).Error
}

// userRoleRepository 用户角色仓储实现
type userRoleRepository struct {
	db *gorm.DB
}

// NewUserRoleRepository 创建用户角色仓储实例
func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}

// Create 创建用户角色关联
func (r *userRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	userRole.CreatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(userRole).Error
}

// Delete 删除用户角色关联
func (r *userRoleRepository) Delete(ctx context.Context, userID, tenantID, roleID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", tenantID).
		Where("role_id = ?", roleID).
		Delete(&entity.UserRole{}).Error
}

// DeleteByUser 删除用户的所有角色
func (r *userRoleRepository) DeleteByUser(ctx context.Context, userID, tenantID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.UserRole{}).Error
}

// DeleteByRole 删除角色的所有用户关联
func (r *userRoleRepository) DeleteByRole(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&entity.UserRole{}).Error
}

// GetRolesByUser 获取用户的所有角色
func (r *userRoleRepository) GetRolesByUser(ctx context.Context, userID, tenantID string) ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.tenant_id = ?", tenantID).
		Where("roles.deleted_at IS NULL").
		Find(&roles).Error

	return roles, err
}

// GetUsersByRole 获取角色的所有用户
func (r *userRoleRepository) GetUsersByRole(ctx context.Context, roleID string) ([]*entity.UserRole, error) {
	var userRoles []*entity.UserRole
	err := r.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Find(&userRoles).Error
	return userRoles, err
}

// Exists 检查用户角色关联是否存在
func (r *userRoleRepository) Exists(ctx context.Context, userID, tenantID, roleID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.UserRole{}).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", tenantID).
		Where("role_id = ?", roleID).
		Count(&count).Error
	return count > 0, err
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
func (r *departmentRepository) GetByID(ctx context.Context, departmentID string) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
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

// GetByName 根据租户ID和名称获取部门
func (r *departmentRepository) GetByName(ctx context.Context, tenantID string, name string) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("department_name = ?", name).
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

// Update 更新部门
func (r *departmentRepository) Update(ctx context.Context, dept *entity.Department) error {
	dept.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("department_id = ?", dept.DepartmentID).
		Updates(dept).Error
}

// Delete 软删除部门
func (r *departmentRepository) Delete(ctx context.Context, departmentID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Department{}).
		Where("department_id = ?", departmentID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// List 获取租户的所有部门
func (r *departmentRepository) List(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Find(&depts).Error
	return depts, err
}

// GetDescendants 获取子部门列表
func (r *departmentRepository) GetDescendants(ctx context.Context, departmentID string) ([]*entity.Department, error) {
	var depts []*entity.Department
	err := r.db.WithContext(ctx).
		Where("parent_department_id = ?", departmentID).
		Where("deleted_at IS NULL").
		Find(&depts).Error
	return depts, err
}

// GetByTenant 获取租户的所有部门（别名方法）
func (r *departmentRepository) GetByTenant(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	return r.List(ctx, tenantID)
}

// userDepartmentRepository 用户部门仓储实现
type userDepartmentRepository struct {
	db *gorm.DB
}

// NewUserDepartmentRepository 创建用户部门仓储实例
func NewUserDepartmentRepository(db *gorm.DB) UserDepartmentRepository {
	return &userDepartmentRepository{db: db}
}

// Create 创建用户部门关联
func (r *userDepartmentRepository) Create(ctx context.Context, userDept *entity.UserDepartment) error {
	userDept.CreatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(userDept).Error
}

// Delete 删除用户部门关联
func (r *userDepartmentRepository) Delete(ctx context.Context, userID, tenantID, departmentID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", tenantID).
		Where("department_id = ?", departmentID).
		Delete(&entity.UserDepartment{}).Error
}

// GetByUser 获取用户的所有部门
func (r *userDepartmentRepository) GetByUser(ctx context.Context, userID, tenantID string) ([]*entity.UserDepartment, error) {
	var userDepts []*entity.UserDepartment
	err := r.db.WithContext(ctx).
		Preload("Department").
		Where("user_id = ?", userID).
		Where("tenant_id = ?", tenantID).
		Find(&userDepts).Error
	return userDepts, err
}

// GetByDepartment 获取部门的所有用户
func (r *userDepartmentRepository) GetByDepartment(ctx context.Context, departmentID string) ([]*entity.UserDepartment, error) {
	var userDepts []*entity.UserDepartment
	err := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Find(&userDepts).Error
	return userDepts, err
}

// GetByUserAndDepartment 获取用户在特定部门的关联
func (r *userDepartmentRepository) GetByUserAndDepartment(ctx context.Context, userID, departmentID string) (*entity.UserDepartment, error) {
	var userDept entity.UserDepartment
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("department_id = ?", departmentID).
		First(&userDept).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &userDept, nil
}
