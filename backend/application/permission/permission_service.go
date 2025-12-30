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
	"encoding/json"
	"time"

	"github.com/coze-dev/coze-studio/backend/api/model/permission"
	permissionentity "github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	permissionrepo "github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	permissionservice "github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/infra/monitoring/metrics"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// PermissionApplicationService 权限应用服务
type PermissionApplicationService struct {
	permissionChecker  *permissionservice.PermissionChecker
	roleSVC          *permissionservice.RoleService
	departmentSVC    *permissionservice.DepartmentService
}

// NewPermissionApplicationService 创建权限应用服务
func NewPermissionApplicationService(
	permissionChecker *permissionservice.PermissionChecker,
	roleSVC *permissionservice.RoleService,
	departmentSVC *permissionservice.DepartmentService,
) *PermissionApplicationService {
	return &PermissionApplicationService{
		permissionChecker: permissionChecker,
		roleSVC:          roleSVC,
		departmentSVC:    departmentSVC,
	}
}

// ==================== 角色管理用例 ====================

// CreateRole 创建角色
func (s *PermissionApplicationService) CreateRole(ctx context.Context, req *permission.CreateRoleRequest) (*permission.RoleInfo, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.RoleCode == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_code is required"))
	}
	if req.RoleName == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_name is required"))
	}
	if req.RoleType == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_type is required"))
	}

	// 2. 构建实体
	role := &permissionentity.Role{
		TenantID:     req.TenantID,
		RoleCode:     req.RoleCode,
		RoleName:     req.RoleName,
		RoleType:     permissionentity.RoleType(req.RoleType),
		Description:  req.Description,
		ParentRoleID: nil,
	}

	if req.ParentRoleID != "" {
		role.ParentRoleID = &req.ParentRoleID
	}

	// 3. 调用领域服务
	createdRole, err := s.roleSVC.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	// 4. 转换为DTO
	return s.entityToRoleInfo(createdRole, nil, nil), nil
}

// GetRole 获取角色
func (s *PermissionApplicationService) GetRole(ctx context.Context, roleID string) (*permission.RoleInfo, error) {
	if roleID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}

	role, err := s.roleSVC.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errorx.New(errno.RoleNotFoundCode, errorx.KV("role_id", roleID))
	}

	// 获取权限信息
	dataPerms, _ := s.roleSVC.GetAllDataPermissions(ctx, roleID)
	fieldPerms, _ := s.roleSVC.GetAllFieldPermissions(ctx, roleID)

	return s.entityToRoleInfo(role, dataPerms, fieldPerms), nil
}

// UpdateRole 更新角色
func (s *PermissionApplicationService) UpdateRole(ctx context.Context, roleID string, req *permission.UpdateRoleRequest) (*permission.RoleInfo, error) {
	if roleID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}

	// 1. 获取现有角色
	role, err := s.roleSVC.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errorx.New(errno.RoleNotFoundCode, errorx.KV("role_id", roleID))
	}

	// 2. 应用更新
	if req.RoleName != nil {
		role.RoleName = *req.RoleName
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	// 3. 保存更新
	updatedRole, err := s.roleSVC.UpdateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return s.entityToRoleInfo(updatedRole, nil, nil), nil
}

// DeleteRole 删除角色
func (s *PermissionApplicationService) DeleteRole(ctx context.Context, roleID string) error {
	if roleID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}

	return s.roleSVC.DeleteRole(ctx, roleID)
}

// ListRoles 列出角色
func (s *PermissionApplicationService) ListRoles(ctx context.Context, req *permission.ListRolesRequest) (*permission.ListRolesData, error) {
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 调用领域服务
	roles, total, err := s.roleSVC.ListRoles(ctx, req.TenantID, (*permissionentity.RoleType)(req.RoleType), pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	roleDTOs := make([]permission.RoleInfo, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, *s.entityToRoleInfo(role, nil, nil))
	}

	return &permission.ListRolesData{
		Roles:      roleDTOs,
		TotalCount: total,
	}, nil
}

// ==================== 用户角色分配用例 ====================

// AssignRole 分配角色
func (s *PermissionApplicationService) AssignRole(ctx context.Context, userID string, req *permission.AssignRoleRequest) (*permission.AssignRoleData, error) {
	if userID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.RoleID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}

	// 调用领域服务
	userRole, err := s.roleSVC.AssignRoleToUser(ctx, userID, req.TenantID, req.RoleID)
	if err != nil {
		return nil, err
	}

	return &permission.AssignRoleData{
		UserID:     userRole.UserID,
		RoleID:     userRole.RoleID,
		AssignedAt: userRole.CreatedAt,
	}, nil
}

// RevokeRole 撤销角色
func (s *PermissionApplicationService) RevokeRole(ctx context.Context, userID, tenantID, roleID string) error {
	if userID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if tenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if roleID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}

	return s.roleSVC.RevokeRoleFromUser(ctx, userID, tenantID, roleID)
}

// GetUserRoles 获取用户的所有角色
func (s *PermissionApplicationService) GetUserRoles(ctx context.Context, userID, tenantID string) (*permission.GetUserRolesData, error) {
	if userID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 调用领域服务
	roles, err := s.roleSVC.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	roleDTOs := make([]permission.RoleInfo, 0, len(roles))
	for _, role := range roles {
		roleDTOs = append(roleDTOs, *s.entityToRoleInfo(role, nil, nil))
	}

	return &permission.GetUserRolesData{
		Roles: roleDTOs,
	}, nil
}

// ==================== 数据权限用例 ====================

// SetDataPermission 设置角色数据权限
func (s *PermissionApplicationService) SetDataPermission(ctx context.Context, roleID string, req *permission.SetDataPermissionRequest) error {
	if roleID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}
	if req.TenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.ResourceType == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "resource_type is required"))
	}
	if req.Scope == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "scope is required"))
	}

	// 构建数据权限实体
	dataPerm := &permissionentity.DataPermission{
		RoleID:       roleID,
		ResourceType: permissionentity.ResourceType(req.ResourceType),
		Scope:        (*permissionentity.DataPermissionLevel)(&req.Scope),
	}

	// 根据scope设置其他字段
	if req.Scope == "dept_sub" || req.Scope == "dept" {
		if req.DepartmentID == "" {
			return errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required for dept/dept_sub scope"))
		}
		dataPerm.DepartmentID = &req.DepartmentID
	}

	if req.Scope == "custom" {
		if req.CustomFilter == nil {
			return errorx.New(errno.InvalidRequest, errorx.KV("msg", "custom_filter is required for custom scope"))
		}
		customFilterJSON, _ := json.Marshal(req.CustomFilter)
		dataPerm.CustomFilter = string(customFilterJSON)
	}

	// 调用领域服务
	return s.roleSVC.SetDataPermission(ctx, dataPerm)
}

// CheckDataPermission 检查数据权限
func (s *PermissionApplicationService) CheckDataPermission(ctx context.Context, req *permission.CheckDataPermissionRequest) (*permission.CheckDataPermissionData, error) {
	if req.UserID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.ResourceType == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "resource_type is required"))
	}
	if req.ResourceID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "resource_id is required"))
	}

	// 记录开始时间
	startTime := time.Now()

	// 调用领域服务
	err := s.permissionChecker.CheckDataPermission(ctx, req.UserID, req.TenantID, permissionentity.ResourceType(req.ResourceType), req.ResourceID)

	// 计算检查耗时
	duration := time.Since(startTime).Seconds()

	if err != nil {
		// 记录权限拒绝指标
		metrics.RecordPermissionCheck(req.TenantID, "data", "denied", duration)
		metrics.RecordPermissionDenied(req.TenantID, "data", "insufficient_scope")

		// 权限不足
		return &permission.CheckDataPermissionData{
			Allowed: false,
			Scope:   "none",
			Reasons: []string{err.Error()},
		}, nil
	}

	// 记录权限允许指标
	metrics.RecordPermissionCheck(req.TenantID, "data", "allowed", duration)

	// 有权限
	return &permission.CheckDataPermissionData{
		Allowed: true,
		Scope:   "all",
		Reasons: []string{"用户有权限访问该资源"},
	}, nil
}

// ==================== 字段权限用例 ====================

// SetFieldPermission 设置角色字段权限
func (s *PermissionApplicationService) SetFieldPermission(ctx context.Context, roleID string, req *permission.SetFieldPermissionRequest) error {
	if roleID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "role_id is required"))
	}
	if req.TenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.ResourceType == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "resource_type is required"))
	}
	if len(req.Permissions) == 0 {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "permissions is required"))
	}

	// 为每个字段设置权限
	for _, perm := range req.Permissions {
		fieldPerm := &permissionentity.FieldPermission{
			RoleID:         roleID,
			ResourceType:   permissionentity.ResourceType(req.ResourceType),
			FieldName:      perm.FieldName,
			PermissionLevel: permissionentity.FieldPermissionLevel(perm.PermissionLevel),
		}

		if err := s.roleSVC.SetFieldPermission(ctx, fieldPerm); err != nil {
			return err
		}
	}

	return nil
}

// GetFieldPermissions 获取字段权限
func (s *PermissionApplicationService) GetFieldPermissions(ctx context.Context, userID, tenantID, resourceType string) (*permission.GetFieldPermissionsData, error) {
	if userID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if resourceType == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "resource_type is required"))
	}

	// 调用领域服务获取字段权限
	permissions, err := s.permissionChecker.GetFieldPermissions(ctx, userID, tenantID, permissionentity.ResourceType(resourceType))
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	permMap := make(map[string]string)
	for field, level := range permissions {
		permMap[field] = string(level)
	}

	return &permission.GetFieldPermissionsData{
		Permissions: permMap,
	}, nil
}

// ==================== 部门管理用例 ====================

// CreateDepartment 创建部门
func (s *PermissionApplicationService) CreateDepartment(ctx context.Context, req *permission.CreateDepartmentRequest) (*permission.DepartmentInfo, error) {
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.DepartmentName == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_name is required"))
	}

	// 构建部门实体
	dept := &permissionentity.Department{
		TenantID:   req.TenantID,
		DeptName:   req.DepartmentName,
		ParentID:   nil,
	}

	if req.ParentDepartmentID != "" {
		dept.ParentID = &req.ParentDepartmentID
	}

	// 调用领域服务
	createdDept, err := s.departmentSVC.CreateDepartment(ctx, dept)
	if err != nil {
		return nil, err
	}

	return s.entityToDepartmentInfo(createdDept, false), nil
}

// GetDepartment 获取部门
func (s *PermissionApplicationService) GetDepartment(ctx context.Context, departmentID string) (*permission.DepartmentInfo, error) {
	if departmentID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}

	dept, err := s.departmentSVC.GetDepartmentByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, errorx.New(errno.DepartmentNotFoundCode, errorx.KV("department_id", departmentID))
	}

	// 检查是否有子部门
	hasChildren, _ := s.departmentSVC.HasChildren(ctx, departmentID)

	return s.entityToDepartmentInfo(dept, hasChildren), nil
}

// UpdateDepartment 更新部门
func (s *PermissionApplicationService) UpdateDepartment(ctx context.Context, departmentID string, req *permission.UpdateDepartmentRequest) (*permission.DepartmentInfo, error) {
	if departmentID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}

	// 1. 获取现有部门
	dept, err := s.departmentSVC.GetDepartmentByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, errorx.New(errno.DepartmentNotFoundCode, errorx.KV("department_id", departmentID))
	}

	// 2. 应用更新
	if req.DepartmentName != nil {
		dept.DeptName = *req.DepartmentName
	}

	// 3. 保存更新
	updatedDept, err := s.departmentSVC.UpdateDepartment(ctx, dept)
	if err != nil {
		return nil, err
	}

	hasChildren, _ := s.departmentSVC.HasChildren(ctx, departmentID)

	return s.entityToDepartmentInfo(updatedDept, hasChildren), nil
}

// DeleteDepartment 删除部门
func (s *PermissionApplicationService) DeleteDepartment(ctx context.Context, departmentID string) error {
	if departmentID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}

	return s.departmentSVC.DeleteDepartment(ctx, departmentID)
}

// ListDepartments 列出部门
func (s *PermissionApplicationService) ListDepartments(ctx context.Context, tenantID string) (*permission.ListDepartmentsData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 调用领域服务
	departments, err := s.departmentSVC.GetDepartmentsByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	deptDTOs := make([]permission.DepartmentInfo, 0, len(departments))
	for _, dept := range departments {
		hasChildren, _ := s.departmentSVC.HasChildren(ctx, dept.DepartmentID)
		deptDTOs = append(deptDTOs, *s.entityToDepartmentInfo(dept, hasChildren))
	}

	return &permission.ListDepartmentsData{
		Departments: deptDTOs,
	}, nil
}

// GetAccessibleDepartments 获取用户可访问的部门
func (s *PermissionApplicationService) GetAccessibleDepartments(ctx context.Context, userID, tenantID string) (*permission.GetAccessibleDepartmentsData, error) {
	if userID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务获取可访问的部门ID列表
	deptIDs, err := s.permissionChecker.GetAccessibleDepartmentIDs(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. 批量获取部门详情
	departments := make([]permission.DepartmentInfo, 0)
	for _, deptID := range deptIDs {
		dept, err := s.departmentSVC.GetDepartmentByID(ctx, deptID)
		if err == nil && dept != nil {
			hasChildren, _ := s.departmentSVC.HasChildren(ctx, deptID)
			departments = append(departments, *s.entityToDepartmentInfo(dept, hasChildren))
		}
	}

	return &permission.GetAccessibleDepartmentsData{
		AccessibleDepartmentIDs: deptIDs,
		Departments:             departments,
	}, nil
}

// ==================== 部门成员管理用例 ====================

// AddDepartmentMember 添加部门成员
func (s *PermissionApplicationService) AddDepartmentMember(ctx context.Context, departmentID string, req *permission.AddDepartmentMemberRequest) error {
	if departmentID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}
	if req.UserID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}
	if req.TenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 构建成员实体
	member := &permissionentity.UserDepartment{
		UserID:      req.UserID,
		DepartmentID: departmentID,
		TenantID:    req.TenantID,
		IsLeader:    req.IsLeader,
	}

	return s.departmentSVC.AddMember(ctx, member)
}

// UpdateDepartmentMember 更新部门成员
func (s *PermissionApplicationService) UpdateDepartmentMember(ctx context.Context, departmentID, userID string, req *permission.UpdateDepartmentMemberRequest) error {
	if departmentID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}
	if userID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}

	if req.IsLeader == nil {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "is_leader is required"))
	}

	return s.departmentSVC.UpdateMember(ctx, departmentID, userID, *req.IsLeader)
}

// RemoveDepartmentMember 移除部门成员
func (s *PermissionApplicationService) RemoveDepartmentMember(ctx context.Context, departmentID, userID string) error {
	if departmentID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}
	if userID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_id is required"))
	}

	return s.departmentSVC.RemoveMember(ctx, departmentID, userID)
}

// GetDepartmentMembers 获取部门成员
func (s *PermissionApplicationService) GetDepartmentMembers(ctx context.Context, departmentID string) (*permission.GetDepartmentMembersData, error) {
	if departmentID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "department_id is required"))
	}

	// 调用领域服务
	members, err := s.departmentSVC.GetMembers(ctx, departmentID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	memberDTOs := make([]permission.DepartmentMember, 0, len(members))
	for _, member := range members {
		memberDTOs = append(memberDTOs, permission.DepartmentMember{
			UserID:       member.UserID,
			DepartmentID: member.DepartmentID,
			IsLeader:     member.IsLeader,
			JoinedAt:     member.CreatedAt,
		})
	}

	return &permission.GetDepartmentMembersData{
		Members: memberDTOs,
	}, nil
}

// ==================== 转换方法 ====================

// entityToRoleInfo 实体转换为RoleInfo DTO
func (s *PermissionApplicationService) entityToRoleInfo(
	entity *permissionentity.Role,
	dataPerms []permissionentity.DataPermission,
	fieldPerms []permissionentity.FieldPermission,
) *permission.RoleInfo {
	roleInfo := &permission.RoleInfo{
		RoleID:       entity.RoleID,
		RoleCode:     entity.RoleCode,
		RoleName:     entity.RoleName,
		RoleType:     string(entity.RoleType),
		Description:  entity.Description,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}

	if entity.ParentRoleID != nil {
		roleInfo.ParentRoleID = *entity.ParentRoleID
	}

	// 转换数据权限
	if len(dataPerms) > 0 {
		roleInfo.Permissions.DataPermissions = make(map[string]string)
		for _, dp := range dataPerms {
			if dp.Scope != nil {
				roleInfo.Permissions.DataPermissions[string(dp.ResourceType)] = string(*dp.Scope)
			}
		}
	}

	// 转换字段权限
	if len(fieldPerms) > 0 {
		roleInfo.Permissions.FieldPermissions = make(map[string][]string)
		for _, fp := range fieldPerms {
			resourceType := string(fp.ResourceType)
			roleInfo.Permissions.FieldPermissions[resourceType] = append(
				roleInfo.Permissions.FieldPermissions[resourceType],
				fmt.Sprintf("%s:%s", fp.FieldName, fp.PermissionLevel),
			)
		}
	}

	return roleInfo
}

// entityToDepartmentInfo 实体转换为DepartmentInfo DTO
func (s *PermissionApplicationService) entityToDepartmentInfo(entity *permissionentity.Department, hasChildren bool) *permission.DepartmentInfo {
	deptInfo := &permission.DepartmentInfo{
		DepartmentID: entity.DepartmentID,
		DepartmentName: entity.DeptName,
		TenantID:      entity.TenantID,
		HasChildren:   hasChildren,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}

	if entity.ParentID != nil {
		deptInfo.ParentDepartmentID = *entity.ParentID
	}

	return deptInfo
}
