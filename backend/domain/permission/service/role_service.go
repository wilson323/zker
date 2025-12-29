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

	"github.com/google/uuid"

	"github.com/coze-studio/backend/domain/permission/entity"
	"github.com/coze-studio/backend/domain/permission/repository"
)

// RoleService 角色管理服务
type RoleService struct {
	roleRepo       repository.RoleRepository
	dataPermRepo   repository.DataPermissionRepository
	fieldPermRepo  repository.FieldPermissionRepository
	userRoleRepo   repository.UserRoleRepository
	permChecker    *PermissionChecker
}

// NewRoleService 创建角色管理服务实例
func NewRoleService(
	roleRepo repository.RoleRepository,
	dataPermRepo repository.DataPermissionRepository,
	fieldPermRepo repository.FieldPermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	permChecker *PermissionChecker,
) *RoleService {
	return &RoleService{
		roleRepo:      roleRepo,
		dataPermRepo:  dataPermRepo,
		fieldPermRepo: fieldPermRepo,
		userRoleRepo:  userRoleRepo,
		permChecker:   permChecker,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	TenantID     string                  `json:"tenant_id" validate:"required"`
	RoleName     string                  `json:"role_name" validate:"required,min=2,max=100"`
	RoleCode     string                  `json:"role_code" validate:"required,min=2,max=50"`
	RoleType     entity.RoleType         `json:"role_type" validate:"required"`
	ParentRoleID *string                 `json:"parent_role_id,omitempty"`
	Description  string                  `json:"description,omitempty"`
	DataPerms    []*CreateDataPermRequest `json:"data_permissions,omitempty"`
	FieldPerms   []*CreateFieldPermRequest `json:"field_permissions,omitempty"`
}

// CreateDataPermRequest 创建数据权限请求
type CreateDataPermRequest struct {
	ResourceType entity.ResourceType        `json:"resource_type" validate:"required"`
	Scope         entity.DataPermissionScope `json:"scope" validate:"required"`
	CustomFilter  string                    `json:"custom_filter,omitempty"`
}

// CreateFieldPermRequest 创建字段权限请求
type CreateFieldPermRequest struct {
	ResourceType    string                       `json:"resource_type" validate:"required"`
	FieldName       string                       `json:"field_name" validate:"required"`
	PermissionLevel entity.FieldPermissionLevel `json:"permission_level" validate:"required"`
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, req *CreateRoleRequest) (*entity.Role, error) {
	// 1. 检查角色编码是否已存在
	existing, err := s.roleRepo.GetByCode(ctx, req.TenantID, req.RoleCode)
	if err != nil {
		return nil, fmt.Errorf("failed to check role code: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("role code %s already exists", req.RoleCode)
	}

	// 2. 创建角色
	role := &entity.Role{
		RoleID:       generateUUID(),
		TenantID:     req.TenantID,
		RoleName:     req.RoleName,
		RoleCode:     req.RoleCode,
		RoleType:     req.RoleType,
		ParentRoleID: req.ParentRoleID,
		Description:  req.Description,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// 3. 创建数据权限
	for _, dataPermReq := range req.DataPerms {
		dataPerm := &entity.DataPermission{
			PermissionID: generateUUID(),
			RoleID:       role.RoleID,
			ResourceType: dataPermReq.ResourceType,
			Scope:        dataPermReq.Scope,
			CustomFilter: dataPermReq.CustomFilter,
		}
		if err := s.dataPermRepo.Create(ctx, dataPerm); err != nil {
			return nil, fmt.Errorf("failed to create data permission: %w", err)
		}
	}

	// 4. 创建字段权限
	for _, fieldPermReq := range req.FieldPerms {
		fieldPerm := &entity.FieldPermission{
			PermissionID:    generateUUID(),
			RoleID:          role.RoleID,
			ResourceType:    fieldPermReq.ResourceType,
			FieldName:       fieldPermReq.FieldName,
			PermissionLevel: fieldPermReq.PermissionLevel,
		}
		if err := s.fieldPermRepo.Create(ctx, fieldPerm); err != nil {
			return nil, fmt.Errorf("failed to create field permission: %w", err)
		}
	}

	// 5. 返回完整的角色信息
	return s.GetRole(ctx, role.RoleID)
}

// GetRole 获取角色信息
func (s *RoleService) GetRole(ctx context.Context, roleID string) (*entity.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("role not found: %s", roleID)
	}

	// 加载数据权限和字段权限
	dataPerms, _ := s.dataPermRepo.GetByRole(ctx, roleID)
	role.DataPerms = dataPerms

	fieldPerms, _ := s.fieldPermRepo.GetByRole(ctx, roleID)
	role.FieldPerms = fieldPerms

	return role, nil
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleID      string                  `json:"role_id" validate:"required"`
	RoleName    string                  `json:"role_name,omitempty" validate:"omitempty,min=2,max=100"`
	Description string                  `json:"description,omitempty"`
	DataPerms   []*CreateDataPermRequest `json:"data_permissions,omitempty"`
	FieldPerms  []*CreateFieldPermRequest `json:"field_permissions,omitempty"`
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*entity.Role, error) {
	// 1. 获取现有角色
	role, err := s.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("role not found: %s", req.RoleID)
	}

	// 2. 系统角色不允许修改核心属性
	if role.IsSystemRole() && len(req.DataPerms) > 0 {
		return nil, fmt.Errorf("cannot modify system role permissions")
	}

	// 3. 更新基本信息
	if req.RoleName != "" {
		role.RoleName = req.RoleName
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// 4. 更新数据权限（先删除旧的，再创建新的）
	if len(req.DataPerms) > 0 {
		if err := s.dataPermRepo.DeleteByRole(ctx, req.RoleID); err != nil {
			return nil, fmt.Errorf("failed to delete old data permissions: %w", err)
		}

		for _, dataPermReq := range req.DataPerms {
			dataPerm := &entity.DataPermission{
				PermissionID: generateUUID(),
				RoleID:       req.RoleID,
				ResourceType: dataPermReq.ResourceType,
				Scope:        dataPermReq.Scope,
				CustomFilter: dataPermReq.CustomFilter,
			}
			if err := s.dataPermRepo.Create(ctx, dataPerm); err != nil {
				return nil, fmt.Errorf("failed to create data permission: %w", err)
			}
		}
	}

	// 5. 更新字段权限
	if len(req.FieldPerms) > 0 {
		if err := s.fieldPermRepo.DeleteByRole(ctx, req.RoleID); err != nil {
			return nil, fmt.Errorf("failed to delete old field permissions: %w", err)
		}

		for _, fieldPermReq := range req.FieldPerms {
			fieldPerm := &entity.FieldPermission{
				PermissionID:    generateUUID(),
				RoleID:          req.RoleID,
				ResourceType:    fieldPermReq.ResourceType,
				FieldName:       fieldPermReq.FieldName,
				PermissionLevel: fieldPermReq.PermissionLevel,
			}
			if err := s.fieldPermRepo.Create(ctx, fieldPerm); err != nil {
				return nil, fmt.Errorf("failed to create field permission: %w", err)
			}
		}
	}

	// 6. 返回更新后的角色信息
	return s.GetRole(ctx, req.RoleID)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, roleID string) error {
	// 1. 获取角色
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found: %s", roleID)
	}

	// 2. 系统角色不允许删除
	if role.IsSystemRole() {
		return fmt.Errorf("cannot delete system role")
	}

	// 3. 删除角色的用户关联
	if err := s.userRoleRepo.DeleteByRole(ctx, roleID); err != nil {
		return fmt.Errorf("failed to delete user role associations: %w", err)
	}

	// 4. 删除数据权限
	if err := s.dataPermRepo.DeleteByRole(ctx, roleID); err != nil {
		return fmt.Errorf("failed to delete data permissions: %w", err)
	}

	// 5. 删除字段权限
	if err := s.fieldPermRepo.DeleteByRole(ctx, roleID); err != nil {
		return fmt.Errorf("failed to delete field permissions: %w", err)
	}

	// 6. 软删除角色
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// ListRolesRequest 查询角色列表请求
type ListRolesRequest struct {
	TenantID string           `json:"tenant_id" validate:"required"`
	RoleType *entity.RoleType `json:"role_type,omitempty"`
	PageToken string           `json:"page_token,omitempty"`
	PageSize  int              `json:"page_size,omitempty"`
}

// ListRoles 查询角色列表
func (s *RoleService) ListRoles(ctx context.Context, req *ListRolesRequest) ([]*entity.Role, int64, error) {
	filter := &repository.RoleFilter{
		TenantID: req.TenantID,
		RoleType: req.RoleType,
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	roles, total, err := s.roleRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	return roles, total, nil
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	TenantID string `json:"tenant_id" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
}

// AssignRole 分配角色给用户
func (s *RoleService) AssignRole(ctx context.Context, req *AssignRoleRequest) error {
	// 1. 检查角色是否存在
	role, err := s.roleRepo.GetByID(ctx, req.RoleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found: %s", req.RoleID)
	}

	// 2. 检查是否已经分配
	exists, err := s.userRoleRepo.Exists(ctx, req.UserID, req.TenantID, req.RoleID)
	if err != nil {
		return fmt.Errorf("failed to check role assignment: %w", err)
	}
	if exists {
		return fmt.Errorf("role already assigned to user")
	}

	// 3. 创建用户角色关联
	userRole := &entity.UserRole{
		UserID:   req.UserID,
		TenantID: req.TenantID,
		RoleID:   req.RoleID,
	}

	if err := s.userRoleRepo.Create(ctx, userRole); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RevokeRole 撤销用户的角色
func (s *RoleService) RevokeRole(ctx context.Context, userID, tenantID, roleID string) error {
	if err := s.userRoleRepo.Delete(ctx, userID, tenantID, roleID); err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}
	return nil
}

// GetUserRoles 获取用户的所有角色
func (s *RoleService) GetUserRoles(ctx context.Context, userID, tenantID string) ([]*entity.Role, error) {
	roles, err := s.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	return roles, nil
}

// CheckPermission 检查权限（对外统一接口）
func (s *RoleService) CheckPermission(
	ctx context.Context,
	userID, tenantID string,
	resourceType entity.ResourceType,
	resourceID string,
) error {
	return s.permChecker.CheckDataPermission(ctx, userID, tenantID, resourceType, resourceID)
}

// GetFieldPermissions 获取字段权限（对外统一接口）
func (s *RoleService) GetFieldPermissions(
	ctx context.Context,
	userID, tenantID string,
	resourceType string,
) (map[string]string, error) {
	return s.permChecker.GetFieldPermissions(ctx, userID, tenantID, resourceType)
}

// generateUUID 生成UUID
func generateUUID() string {
	return uuid.New().String()
}
