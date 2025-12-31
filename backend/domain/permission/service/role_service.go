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

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
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
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to check role code"),
            )
	}
	if existing != nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("field", "role_code"),
                errorx.KV("value", req.RoleCode),
                errorx.KV("reason", "already exists"),
            )
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
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to create role"),
            )
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
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to create data permission"),
            )
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
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to create field permission"),
            )
		}
	}

	// 5. 返回完整的角色信息
	return s.GetRole(ctx, role.RoleID)
}

// GetRole 获取角色信息
func (s *RoleService) GetRole(ctx context.Context, roleID string) (*entity.Role, error) {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get role"),
            )
	}
	if role == nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("resource", "role"),
                errorx.KV("role_id", roleID),
                errorx.KV("reason", "not found"),
            )
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
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get role"),
            )
	}
	if role == nil {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("resource", "role"),
                errorx.KV("role_id", req.RoleID),
                errorx.KV("reason", "not found"),
            )
	}

	// 2. 系统角色不允许修改核心属性
	if role.IsSystemRole() && len(req.DataPerms) > 0 {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "cannot modify system role permissions"),
            )
	}

	// 3. 更新基本信息
	if req.RoleName != "" {
		role.RoleName = req.RoleName
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to update role"),
            )
	}

	// 4. 更新数据权限（先删除旧的，再创建新的）
	if len(req.DataPerms) > 0 {
		if err := s.dataPermRepo.DeleteByRole(ctx, req.RoleID); err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to delete old data permissions"),
            )
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
				return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to create data permission"),
            )
			}
		}
	}

	// 5. 更新字段权限
	if len(req.FieldPerms) > 0 {
		if err := s.fieldPermRepo.DeleteByRole(ctx, req.RoleID); err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to delete old field permissions"),
            )
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
				return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to create field permission"),
            )
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
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get role"),
            )
	}
	if role == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("resource", "role"),
                errorx.KV("role_id", roleID),
                errorx.KV("reason", "not found"),
            )
	}

	// 2. 系统角色不允许删除
	if role.IsSystemRole() {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "cannot delete system role"),
            )
	}

	// 3. 删除角色的用户关联
	if err := s.userRoleRepo.DeleteByRole(ctx, roleID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to delete user role associations"),
            )
	}

	// 4. 删除数据权限
	if err := s.dataPermRepo.DeleteByRole(ctx, roleID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to delete data permissions"),
            )
	}

	// 5. 删除字段权限
	if err := s.fieldPermRepo.DeleteByRole(ctx, roleID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to delete field permissions"),
            )
	}

	// 6. 软删除角色
	if err := s.roleRepo.Delete(ctx, roleID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to delete role"),
            )
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
		return nil, 0, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "failed to list roles"),
            )
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
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get role"),
            )
	}
	if role == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("resource", "role"),
                errorx.KV("role_id", req.RoleID),
                errorx.KV("reason", "not found"),
            )
	}

	// 2. 检查是否已经分配
	exists, err := s.userRoleRepo.Exists(ctx, req.UserID, req.TenantID, req.RoleID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to check role assignment"),
            )
	}
	if exists {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "role already assigned to user"),
            )
	}

	// 3. 创建用户角色关联
	userRole := &entity.UserRole{
		UserID:   req.UserID,
		TenantID: req.TenantID,
		RoleID:   req.RoleID,
	}

	if err := s.userRoleRepo.Create(ctx, userRole); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to assign role"),
            )
	}

	return nil
}

// RevokeRole 撤销用户的角色
func (s *RoleService) RevokeRole(ctx context.Context, userID, tenantID, roleID string) error {
	if err := s.userRoleRepo.Delete(ctx, userID, tenantID, roleID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to revoke role"),
            )
	}
	return nil
}

// GetUserRoles 获取用户的所有角色
func (s *RoleService) GetUserRoles(ctx context.Context, userID, tenantID string) ([]*entity.Role, error) {
	roles, err := s.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get user roles"),
            )
	}
	return roles, nil
}

// CheckPermission 检查权限（对外统一接口）
func (s *RoleService) CheckPermission(
	ctx context.Context,
	userID, tenantID string,
	resourceType entity.ResourceType,
	action string,
	resourceID string,
) error {
	allowed, err := s.permChecker.CheckDataPermission(ctx, tenantID, userID, resourceType, action, resourceID)
	if err != nil {
		return err
	}
	if !allowed {
		return &PermissionDeniedError{
			UserID:       userID,
			ResourceType: string(resourceType),
			ResourceID:   resourceID,
			Reason:       "permission denied",
		}
	}
	return nil
}

// GetFieldPermissions 获取字段权限（对外统一接口）
func (s *RoleService) GetFieldPermissions(
	ctx context.Context,
	tenantID, userID string,
	resourceType string,
) (map[string]string, error) {
	return s.permChecker.GetFieldPermissions(ctx, tenantID, userID, resourceType)
}

// InitializeSystemRoles 初始化系统预置角色
// 调用时机: 新租户创建时
// 幂等性: 如果角色已存在则跳过
func (s *RoleService) InitializeSystemRoles(ctx context.Context, tenantID string) error {
	// 1. 检查是否已初始化
	filter := &repository.RoleFilter{
		TenantID: tenantID,
		PageSize:  100, // 足够大的数字以获取所有角色
	}
	existingRoles, _, err := s.roleRepo.List(ctx, filter)
	if err == nil && len(existingRoles) > 0 {
		// 已有角色,检查是否包含系统角色
		for _, role := range existingRoles {
			if role.IsSystemRole() {
				// 已有系统角色,无需重复初始化
				return nil
			}
		}
	}

	// 2. 创建4个系统预置角色
	roles := []*entity.Role{
		{
			RoleID:     generateRoleID("tenant_owner", tenantID),
			TenantID:   tenantID,
			RoleName:   "租户所有者",
			RoleCode:   "tenant_owner",
			RoleType:   entity.RoleTypeSystem,
			Description: "拥有租户内所有资源的完整权限",
		},
		{
			RoleID:     generateRoleID("tenant_admin", tenantID),
			TenantID:   tenantID,
			RoleName:   "租户管理员",
			RoleCode:   "tenant_admin",
			RoleType:   entity.RoleTypeSystem,
			Description: "拥有租户内部门及以下资源的完整权限",
		},
		{
			RoleID:     generateRoleID("tenant_member", tenantID),
			TenantID:   tenantID,
			RoleName:   "普通成员",
			RoleCode:   "tenant_member",
			RoleType:   entity.RoleTypeSystem,
			Description: "仅能访问自己创建的资源",
		},
		{
			RoleID:     generateRoleID("tenant_viewer", tenantID),
			TenantID:   tenantID,
			RoleName:   "只读成员",
			RoleCode:   "tenant_viewer",
			RoleType:   entity.RoleTypeSystem,
			Description: "仅能查看自己创建的资源",
		},
	}

	// 3. 批量创建角色
	for _, role := range roles {
		if err := s.roleRepo.Create(ctx, role); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "创建系统角色失败"),
                errorx.KV("role_code", role.RoleCode),
            )
		}

		// 4. 为每个角色初始化数据权限
		if err := s.initializeDataPermissions(ctx, role); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "初始化角色数据权限失败"),
                errorx.KV("role_code", role.RoleCode),
            )
		}

		// 5. 为每个角色初始化字段权限
		if err := s.initializeFieldPermissions(ctx, role); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "初始化角色字段权限失败"),
                errorx.KV("role_code", role.RoleCode),
            )
		}
	}

	return nil
}

// initializeDataPermissions 初始化角色的数据权限
func (s *RoleService) initializeDataPermissions(ctx context.Context, role *entity.Role) error {
	// 根据角色类型配置数据权限
	var scope entity.DataPermissionScope
	switch role.RoleCode {
	case "tenant_owner":
		scope = entity.DataPermissionScopeAll
	case "tenant_admin":
		scope = entity.DataPermissionScopeDepartment
	case "tenant_member", "tenant_viewer":
		scope = entity.DataPermissionScopeOwn
	default:
		scope = entity.DataPermissionScopeNone
	}

	// 为每种资源类型创建数据权限
	resourceTypes := []entity.ResourceType{
		entity.ResourceTypeBots,
		entity.ResourceTypeConversations,
		entity.ResourceTypeKnowledge,
		entity.ResourceTypeWorkflows,
		entity.ResourceTypePlugins,
	}

	for _, rt := range resourceTypes {
		perm := &entity.DataPermission{
			PermissionID: generatePermissionID(role.RoleID, string(rt)),
			RoleID:       role.RoleID,
			ResourceType: rt,
			Scope:        scope,
		}
		if err := s.dataPermRepo.Create(ctx, perm); err != nil {
			return err
		}
	}

	return nil
}

// initializeFieldPermissions 初始化角色的字段权限
func (s *RoleService) initializeFieldPermissions(ctx context.Context, role *entity.Role) error {
	// 定义敏感字段
	sensitiveFields := []struct {
		resourceType string
		fieldName     string
	}{
		{"bots", "api_key"},
		{"bots", "secret_key"},
		{"bots", "webhook_url"},
		{"workflows", "api_key"},
		{"workflows", "auth_token"},
		{"plugins", "api_key"},
		{"plugins", "secret_key"},
	}

	// 根据角色类型配置字段权限
	for _, sf := range sensitiveFields {
		var permLevel entity.FieldPermissionLevel
		switch role.RoleCode {
		case "tenant_owner":
			// 所有者可编辑所有字段,无需创建记录
			continue
		case "tenant_admin", "tenant_member":
			// 管理员和成员: 敏感字段只读
			permLevel = entity.FieldPermissionLevelReadonly
		case "tenant_viewer":
			// 查看者: 敏感字段隐藏
			permLevel = entity.FieldPermissionLevelHidden
		default:
			continue
		}

		fieldPerm := &entity.FieldPermission{
			PermissionID:    generateUUID(),
			RoleID:          role.RoleID,
			ResourceType:    sf.resourceType,
			FieldName:       sf.fieldName,
			PermissionLevel: permLevel,
		}
		if err := s.fieldPermRepo.Create(ctx, fieldPerm); err != nil {
			return err
		}
	}

	return nil
}

// generateRoleID 生成角色ID
func generateRoleID(roleCode, tenantID string) string {
	return fmt.Sprintf("role_%s_%s", roleCode, tenantID[:8])
}

// generatePermissionID 生成权限ID
func generatePermissionID(roleID, resourceType string) string {
	return fmt.Sprintf("perm_%s_%s", roleID, resourceType)
}

// generateUUID 生成UUID
func generateUUID() string {
	return uuid.New().String()
}
