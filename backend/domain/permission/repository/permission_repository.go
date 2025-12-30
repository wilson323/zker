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

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) error

	// GetByID 根据ID获取角色
	GetByID(ctx context.Context, roleID string) (*entity.Role, error)

	// GetByCode 根据租户ID和角色编码获取角色
	GetByCode(ctx context.Context, tenantID string, roleCode string) (*entity.Role, error)

	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) error

	// Delete 软删除角色
	Delete(ctx context.Context, roleID string) error

	// List 分页查询角色列表
	List(ctx context.Context, filter *RoleFilter) ([]*entity.Role, int64, error)

	// GetByUser 获取用户的所有角色
	GetByUser(ctx context.Context, userID, tenantID string) ([]*entity.Role, error)
}

// RoleFilter 角色查询过滤器
type RoleFilter struct {
	TenantID string
	RoleType *entity.RoleType
	PageToken string
	PageSize  int
}

// DataPermissionRepository 数据权限仓储接口
type DataPermissionRepository interface {
	// Create 创建数据权限
	Create(ctx context.Context, perm *entity.DataPermission) error

	// GetByID 根据ID获取数据权限
	GetByID(ctx context.Context, permissionID string) (*entity.DataPermission, error)

	// GetByRoleAndResource 根据角色ID和资源类型获取数据权限
	GetByRoleAndResource(ctx context.Context, roleID string, resourceType entity.ResourceType) (*entity.DataPermission, error)

	// GetByRole 获取角色的所有数据权限
	GetByRole(ctx context.Context, roleID string) ([]*entity.DataPermission, error)

	// Update 更新数据权限
	Update(ctx context.Context, perm *entity.DataPermission) error

	// Delete 删除数据权限
	Delete(ctx context.Context, permissionID string) error

	// DeleteByRole 删除角色的所有数据权限
	DeleteByRole(ctx context.Context, roleID string) error
}

// FieldPermissionRepository 字段权限仓储接口
type FieldPermissionRepository interface {
	// Create 创建字段权限
	Create(ctx context.Context, perm *entity.FieldPermission) error

	// GetByID 根据ID获取字段权限
	GetByID(ctx context.Context, permissionID string) (*entity.FieldPermission, error)

	// GetByRoleAndResource 根据角色ID和资源类型获取字段权限列表
	GetByRoleAndResource(ctx context.Context, roleID string, resourceType string) ([]*entity.FieldPermission, error)

	// GetByRole 获取角色的所有字段权限
	GetByRole(ctx context.Context, roleID string) ([]*entity.FieldPermission, error)

	// Update 更新字段权限
	Update(ctx context.Context, perm *entity.FieldPermission) error

	// Delete 删除字段权限
	Delete(ctx context.Context, permissionID string) error

	// DeleteByRole 删除角色的所有字段权限
	DeleteByRole(ctx context.Context, roleID string) error
}

// UserRoleRepository 用户角色仓储接口
type UserRoleRepository interface {
	// Create 创建用户角色关联
	Create(ctx context.Context, userRole *entity.UserRole) error

	// Delete 删除用户角色关联
	Delete(ctx context.Context, userID, tenantID, roleID string) error

	// DeleteByUser 删除用户的所有角色
	DeleteByUser(ctx context.Context, userID, tenantID string) error

	// DeleteByRole 删除角色的所有用户关联
	DeleteByRole(ctx context.Context, roleID string) error

	// GetRolesByUser 获取用户的所有角色
	GetRolesByUser(ctx context.Context, userID, tenantID string) ([]*entity.Role, error)

	// GetUsersByRole 获取角色的所有用户
	GetUsersByRole(ctx context.Context, roleID string) ([]*entity.UserRole, error)

	// Exists 检查用户角色关联是否存在
	Exists(ctx context.Context, userID, tenantID, roleID string) (bool, error)
}

// DepartmentRepository 部门仓储接口
type DepartmentRepository interface {
	// Create 创建部门
	Create(ctx context.Context, dept *entity.Department) error

	// GetByID 根据ID获取部门
	GetByID(ctx context.Context, departmentID string) (*entity.Department, error)

	// GetByName 根据租户ID和名称获取部门
	GetByName(ctx context.Context, tenantID string, name string) (*entity.Department, error)

	// Update 更新部门
	Update(ctx context.Context, dept *entity.Department) error

	// Delete 软删除部门
	Delete(ctx context.Context, departmentID string) error

	// List 获取租户的所有部门
	List(ctx context.Context, tenantID string) ([]*entity.Department, error)

	// GetDescendants 获取子部门列表
	GetDescendants(ctx context.Context, departmentID string) ([]*entity.Department, error)

	// GetByTenant 获取租户的所有部门（别名方法）
	GetByTenant(ctx context.Context, tenantID string) ([]*entity.Department, error)
}

// UserDepartmentRepository 用户部门仓储接口
type UserDepartmentRepository interface {
	// Create 创建用户部门关联
	Create(ctx context.Context, userDept *entity.UserDepartment) error

	// Delete 删除用户部门关联
	Delete(ctx context.Context, userID, tenantID, departmentID string) error

	// GetByUser 获取用户的所有部门
	GetByUser(ctx context.Context, userID, tenantID string) ([]*entity.UserDepartment, error)

	// GetByDepartment 获取部门的所有用户
	GetByDepartment(ctx context.Context, departmentID string) ([]*entity.UserDepartment, error)

	// GetByUserAndDepartment 获取用户在特定部门的关联
	GetByUserAndDepartment(ctx context.Context, userID, departmentID string) (*entity.UserDepartment, error)
}

// TemporaryGrantRepository 临时授权仓储接口
type TemporaryGrantRepository interface {
	// Create 创建临时授权
	Create(ctx context.Context, grant *entity.TemporaryGrant) error

	// GetByID 根据ID获取临时授权
	GetByID(ctx context.Context, grantID int64) (*entity.TemporaryGrant, error)

	// GetByGrantCode 根据授权码获取临时授权
	GetByGrantCode(ctx context.Context, grantCode string) (*entity.TemporaryGrant, error)

	// Update 更新临时授权
	Update(ctx context.Context, grant *entity.TemporaryGrant) error

	// Delete 删除临时授权
	Delete(ctx context.Context, grantID int64) error

	// List 获取租户的所有临时授权
	List(ctx context.Context, filter *TemporaryGrantFilter) ([]*entity.TemporaryGrant, int64, error)

	// ListExpired 获取所有过期的临时授权
	ListExpired(ctx context.Context, expiresBefore int64) ([]*entity.TemporaryGrant, error)

	// DeleteExpired 删除过期的临时授权
	DeleteExpired(ctx context.Context, expiresBefore int64) (int64, error)

	// GetByGrantee 获取被授权人的所有临时授权
	GetByGrantee(ctx context.Context, granteeID string) ([]*entity.TemporaryGrant, error)

	// GetByGrantor 获取授权人的所有临时授权
	GetByGrantor(ctx context.Context, grantorID string) ([]*entity.TemporaryGrant, error)
}

// TemporaryGrantFilter 临时授权查询过滤器
type TemporaryGrantFilter struct {
	TenantID       string
	GranteeID      string
	GrantorID      string
	PermissionType *entity.PermissionType
	IsUsed         *bool
	IsRevoked      *bool
	PageToken      string
	PageSize       int
}

// TemporaryGrantHistoryRepository 临时授权历史仓储接口
type TemporaryGrantHistoryRepository interface {
	// Create 创建历史记录
	Create(ctx context.Context, history *entity.TemporaryGrantHistory) error

	// GetByGrantID 获取临时授权的所有历史记录
	GetByGrantID(ctx context.Context, grantID int64) ([]*entity.TemporaryGrantHistory, error)

	// GetByGrantCode 根据授权码获取历史记录
	GetByGrantCode(ctx context.Context, grantCode string) ([]*entity.TemporaryGrantHistory, error)

	// List 获取租户的所有历史记录
	List(ctx context.Context, filter *HistoryFilter) ([]*entity.TemporaryGrantHistory, int64, error)

	// DeleteByGrantID 删除临时授权的所有历史记录
	DeleteByGrantID(ctx context.Context, grantID int64) error

	// DeleteExpired 删除过期的历史记录（例如90天前）
	DeleteExpired(ctx context.Context, before int64) (int64, error)
}

// HistoryFilter 历史记录查询过滤器
type HistoryFilter struct {
	TenantID      string
	GrantID       int64
	GrantCode     string
	ActionType    *entity.ActionType
	OperatorID    string
	PageToken     string
	PageSize      int
	Before        int64 // 查询指定时间之前的记录
	After         int64 // 查询指定时间之后的记录
}
