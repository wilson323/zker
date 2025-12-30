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

// ==================== Request DTOs ====================

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	RoleCode     string `json:"role_code" binding:"required,min=1,max=100"`
	RoleName     string `json:"role_name" binding:"required,min=1,max=200"`
	RoleType     string `json:"role_type" binding:"required,oneof=system custom"`
	Description  string `json:"description" binding:"max=500"`
	ParentRoleID string `json:"parent_role_id,omitempty"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleName    *string `json:"role_name,omitempty" binding:"omitempty,min=1,max=200"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	RoleID   string `json:"role_id" binding:"required"`
}

// SetDataPermissionRequest 设置数据权限请求
type SetDataPermissionRequest struct {
	TenantID     string                 `json:"tenant_id" binding:"required"`
	ResourceType string                 `json:"resource_type" binding:"required"`
	Scope        string                 `json:"scope" binding:"required,oneof=all dept_sub dept self custom none"`
	DepartmentID string                 `json:"department_id,omitempty"`
	CustomFilter map[string]interface{} `json:"custom_filter,omitempty"`
}

// SetFieldPermissionRequest 设置字段权限请求
type SetFieldPermissionRequest struct {
	TenantID     string                `json:"tenant_id" binding:"required"`
	ResourceType string                `json:"resource_type" binding:"required"`
	Permissions  []FieldPermissionItem `json:"permissions" binding:"required"`
}

// FieldPermissionItem 字段权限项
type FieldPermissionItem struct {
	FieldName       string `json:"field_name" binding:"required"`
	PermissionLevel string `json:"permission_level" binding:"required,oneof=editable readonly hidden"`
}

// CheckDataPermissionRequest 检查数据权限请求
type CheckDataPermissionRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	TenantID     string `json:"tenant_id" binding:"required"`
	ResourceType string `json:"resource_type" binding:"required"`
	ResourceID   string `json:"resource_id" binding:"required"`
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	TenantID           string `json:"tenant_id" binding:"required"`
	DepartmentName     string `json:"department_name" binding:"required,min=1,max=200"`
	ParentDepartmentID string `json:"parent_department_id,omitempty"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	DepartmentName *string `json:"department_name,omitempty" binding:"omitempty,min=1,max=200"`
}

// AddDepartmentMemberRequest 添加部门成员请求
type AddDepartmentMemberRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	TenantID string `json:"tenant_id" binding:"required"`
	IsLeader bool   `json:"is_leader"`
}

// UpdateDepartmentMemberRequest 更新部门成员请求
type UpdateDepartmentMemberRequest struct {
	IsLeader *bool `json:"is_leader"`
}

// ListRolesRequest 列出角色请求
type ListRolesRequest struct {
	TenantID string  `form:"tenant_id" binding:"required"`
	RoleType *string `form:"role_type" binding:"omitempty,oneof=system custom"`
	PageSize int     `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// ==================== Response DTOs ====================

// RoleInfo 角色信息
type RoleInfo struct {
	RoleID       string          `json:"role_id"`
	RoleCode     string          `json:"role_code"`
	RoleName     string          `json:"role_name"`
	RoleType     string          `json:"role_type"`
	Description  string          `json:"description"`
	ParentRoleID string          `json:"parent_role_id,omitempty"`
	Permissions  RolePermissions `json:"permissions,omitempty"`
	CreatedAt    int64           `json:"created_at"`
	UpdatedAt    int64           `json:"updated_at"`
}

// RolePermissions 角色权限信息
type RolePermissions struct {
	DataPermissions  map[string]string   `json:"data_permissions,omitempty"`
	FieldPermissions map[string][]string `json:"field_permissions,omitempty"`
}

// CreateRoleResponse 创建角色响应
type CreateRoleResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    RoleInfo `json:"data"`
}

// GetRoleResponse 获取角色响应
type GetRoleResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    RoleInfo `json:"data"`
}

// ListRolesResponse 列出角色响应
type ListRolesResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    ListRolesData `json:"data"`
}

// ListRolesData 列出角色数据
type ListRolesData struct {
	Roles      []RoleInfo `json:"roles"`
	TotalCount int        `json:"total_count"`
}

// AssignRoleResponse 分配角色响应
type AssignRoleResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    AssignRoleData `json:"data"`
}

// AssignRoleData 分配角色数据
type AssignRoleData struct {
	UserID     string `json:"user_id"`
	RoleID     string `json:"role_id"`
	AssignedAt int64  `json:"assigned_at"`
}

// GetUserRolesResponse 获取用户角色响应
type GetUserRolesResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    GetUserRolesData `json:"data"`
}

// GetUserRolesData 获取用户角色数据
type GetUserRolesData struct {
	Roles []RoleInfo `json:"roles"`
}

// CheckDataPermissionResponse 检查数据权限响应
type CheckDataPermissionResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    CheckDataPermissionData `json:"data"`
}

// CheckDataPermissionData 检查数据权限数据
type CheckDataPermissionData struct {
	Allowed bool     `json:"allowed"`
	Scope   string   `json:"scope"`
	Reasons []string `json:"reasons"`
}

// SetDataPermissionResponse 设置数据权限响应
type SetDataPermissionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetFieldPermissionsResponse 获取字段权限响应
type GetFieldPermissionsResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    GetFieldPermissionsData `json:"data"`
}

// GetFieldPermissionsData 获取字段权限数据
type GetFieldPermissionsData struct {
	Permissions map[string]string `json:"permissions"`
}

// SetFieldPermissionResponse 设置字段权限响应
type SetFieldPermissionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DepartmentInfo 部门信息
type DepartmentInfo struct {
	DepartmentID       string `json:"department_id"`
	DepartmentName     string `json:"department_name"`
	ParentDepartmentID string `json:"parent_department_id,omitempty"`
	TenantID           string `json:"tenant_id"`
	IsLeader           bool   `json:"is_leader"`
	HasChildren        bool   `json:"has_children"`
	CreatedAt          int64  `json:"created_at"`
	UpdatedAt          int64  `json:"updated_at"`
}

// CreateDepartmentResponse 创建部门响应
type CreateDepartmentResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    DepartmentInfo `json:"data"`
}

// GetDepartmentResponse 获取部门响应
type GetDepartmentResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Data    DepartmentInfo `json:"data"`
}

// ListDepartmentsResponse 列出部门响应
type ListDepartmentsResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    ListDepartmentsData `json:"data"`
}

// ListDepartmentsData 列出部门数据
type ListDepartmentsData struct {
	Departments []DepartmentInfo `json:"departments"`
}

// GetAccessibleDepartmentsResponse 获取可访问部门响应
type GetAccessibleDepartmentsResponse struct {
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Data    GetAccessibleDepartmentsData `json:"data"`
}

// GetAccessibleDepartmentsData 获取可访问部门数据
type GetAccessibleDepartmentsData struct {
	AccessibleDepartmentIDs []string         `json:"accessible_department_ids"`
	Departments             []DepartmentInfo `json:"departments"`
}

// AddDepartmentMemberResponse 添加部门成员响应
type AddDepartmentMemberResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetDepartmentMembersResponse 获取部门成员响应
type GetDepartmentMembersResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    GetDepartmentMembersData `json:"data"`
}

// GetDepartmentMembersData 获取部门成员数据
type GetDepartmentMembersData struct {
	Members []DepartmentMember `json:"members"`
}

// DepartmentMember 部门成员
type DepartmentMember struct {
	UserID       string `json:"user_id"`
	DepartmentID string `json:"department_id"`
	IsLeader     bool   `json:"is_leader"`
	JoinedAt     int64  `json:"joined_at"`
}

// ==================== Temporary Grant Request/Response ====================

// CreateTemporaryGrantRequest 创建临时授权请求
type CreateTemporaryGrantRequest struct {
	TenantID       string `json:"tenant_id" binding:"required"`
	GranteeID      string `json:"grantee_id" binding:"required"`
	PermissionType string `json:"permission_type" binding:"required,oneof=role data_permission field_permission"`
	ResourceID     string `json:"resource_id,omitempty"`
	PermissionData string `json:"permission_data,omitempty"`
	ExpiresAt      int64  `json:"expires_at" binding:"required"`
	Reason         string `json:"reason,omitempty"`
}

// CreateTemporaryGrantResponse 创建临时授权响应
type CreateTemporaryGrantResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    TemporaryGrantDetail `json:"data"`
}

// TemporaryGrantDetail 临时授权详情
type TemporaryGrantDetail struct {
	GrantID        string `json:"grant_id"`
	GrantCode      string `json:"grant_code"`
	GranteeID      string `json:"grantee_id"`
	GrantorID      string `json:"grantor_id"`
	TenantID       string `json:"tenant_id"`
	PermissionType string `json:"permission_type"`
	ResourceID     string `json:"resource_id,omitempty"`
	PermissionData string `json:"permission_data,omitempty"`
	ExpiresAt      int64  `json:"expires_at"`
	IsUsed         bool   `json:"is_used"`
	IsRevoked      bool   `json:"is_revoked"`
	CreatedAt      int64  `json:"created_at"`
}

// GetTemporaryGrantResponse 获取临时授权响应
type GetTemporaryGrantResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    TemporaryGrantDetail `json:"data"`
}

// UseTemporaryGrantRequest 使用临时授权请求
type UseTemporaryGrantRequest struct {
	GrantCode string `json:"grant_code" binding:"required"`
}

// UseTemporaryGrantResponse 使用临时授权响应
type UseTemporaryGrantResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    TemporaryGrantDetail `json:"data"`
}

// RevokeTemporaryGrantRequest 撤销临时授权请求
type RevokeTemporaryGrantRequest struct {
	GrantCode string `json:"grant_code" binding:"required"`
	Reason    string `json:"reason,omitempty"`
}

// RevokeTemporaryGrantResponse 撤销临时授权响应
type RevokeTemporaryGrantResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListTemporaryGrantsRequest 查询临时授权列表请求
type ListTemporaryGrantsRequest struct {
	TenantID       string `form:"tenant_id" binding:"required"`
	GranteeID      string `form:"grantee_id,omitempty"`
	PermissionType string `form:"permission_type,omitempty"`
	IsUsed         *bool  `form:"is_used,omitempty"`
	IsRevoked      *bool  `form:"is_revoked,omitempty"`
	Page           int    `form:"page" binding:"min=1"`
	PageSize       int    `form:"page_size" binding:"min=1,max=100"`
}

// ListTemporaryGrantsResponse 查询临时授权列表响应
type ListTemporaryGrantsResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    ListTemporaryGrantsData `json:"data"`
}

// ListTemporaryGrantsData 临时授权列表数据
type ListTemporaryGrantsData struct {
	Total int64                  `json:"total"`
	Items []TemporaryGrantDetail `json:"items"`
}

// GetGrantHistoryResponse 获取授权历史响应
type GetGrantHistoryResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Data    []TemporaryGrantHistory `json:"data"`
}

// TemporaryGrantHistory 临时授权历史
type TemporaryGrantHistory struct {
	HistoryID  int64  `json:"history_id"`
	GrantID    int64  `json:"grant_id"`
	GrantCode  string `json:"grant_code"`
	Action     string `json:"action"`
	OperatorID string `json:"operator_id"`
	Reason     string `json:"reason,omitempty"`
	CreatedAt  int64  `json:"created_at"`
}
