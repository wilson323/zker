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

package permission

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// DepartmentServiceAdapter 部门服务适配器接口
// 这个接口定义了权限应用服务需要的部门服务方法
type DepartmentServiceAdapter interface {
	CreateDepartment(ctx context.Context, dept *entity.Department) (*entity.Department, error)
	GetDepartmentByID(ctx context.Context, departmentID string) (*entity.Department, error)
	UpdateDepartment(ctx context.Context, dept *entity.Department) (*entity.Department, error)
	DeleteDepartment(ctx context.Context, departmentID string) error
	GetDepartmentsByTenant(ctx context.Context, tenantID string) ([]*entity.Department, error)
	HasChildren(ctx context.Context, departmentID string) (bool, error)
	AddMember(ctx context.Context, member *entity.UserDepartment) error
	UpdateMember(ctx context.Context, departmentID, userID string, isLeader bool) error
	RemoveMember(ctx context.Context, departmentID, userID string) error
	GetMembers(ctx context.Context, departmentID string) ([]*entity.UserDepartment, error)
}

// DefaultDepartmentAdapter 默认的部门服务适配器实现
// 这是一个临时的实现，实际应该由 org 模块提供真实实现
type DefaultDepartmentAdapter struct{}

// NewDefaultDepartmentAdapter 创建默认部门服务适配器
func NewDefaultDepartmentAdapter() *DefaultDepartmentAdapter {
	return &DefaultDepartmentAdapter{}
}

// CreateDepartment 创建部门
func (a *DefaultDepartmentAdapter) CreateDepartment(ctx context.Context, dept *entity.Department) (*entity.Department, error) {
	return nil, errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// GetDepartmentByID 获取部门
func (a *DefaultDepartmentAdapter) GetDepartmentByID(ctx context.Context, departmentID string) (*entity.Department, error) {
	return nil, errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// UpdateDepartment 更新部门
func (a *DefaultDepartmentAdapter) UpdateDepartment(ctx context.Context, dept *entity.Department) (*entity.Department, error) {
	return nil, errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// DeleteDepartment 删除部门
func (a *DefaultDepartmentAdapter) DeleteDepartment(ctx context.Context, departmentID string) error {
	return errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// GetDepartmentsByTenant 获取租户的所有部门
func (a *DefaultDepartmentAdapter) GetDepartmentsByTenant(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	return nil, errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// HasChildren 检查是否有子部门
func (a *DefaultDepartmentAdapter) HasChildren(ctx context.Context, departmentID string) (bool, error) {
	return false, errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// AddMember 添加部门成员
func (a *DefaultDepartmentAdapter) AddMember(ctx context.Context, member *entity.UserDepartment) error {
	return errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// UpdateMember 更新部门成员
func (a *DefaultDepartmentAdapter) UpdateMember(ctx context.Context, departmentID, userID string, isLeader bool) error {
	return errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// RemoveMember 移除部门成员
func (a *DefaultDepartmentAdapter) RemoveMember(ctx context.Context, departmentID, userID string) error {
	return errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}

// GetMembers 获取部门成员
func (a *DefaultDepartmentAdapter) GetMembers(ctx context.Context, departmentID string) ([]*entity.UserDepartment, error) {
	return nil, errorx.New(errno.ErrPermissionCheckFailedCode,
		errorx.KV("reason", "department service not implemented in permission module"),
	)
}
