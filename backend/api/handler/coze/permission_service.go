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

package coze

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/model/permission"
	permissionapp "github.com/coze-studio/backend/application/permission"
)

// ==================== 角色管理接口 ====================

// CreateRole 创建角色
// @router /api/roles [POST]
func CreateRole(ctx context.Context, c *app.RequestContext) {
	var err error
	var req permission.CreateRoleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.CreateRole(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRole 获取角色
// @router /api/roles/:role_id [GET]
func GetRole(ctx context.Context, c *app.RequestContext) {
	roleID := c.Param("role_id")
	if roleID == "" {
		invalidParamRequestResponse(c, "role_id is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.GetRole(ctx, roleID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRole 更新角色
// @router /api/roles/:role_id [PUT]
func UpdateRole(ctx context.Context, c *app.RequestContext) {
	roleID := c.Param("role_id")
	if roleID == "" {
		invalidParamRequestResponse(c, "role_id is required")
		return
	}

	var err error
	var req permission.UpdateRoleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.UpdateRole(ctx, roleID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteRole 删除角色
// @router /api/roles/:role_id [DELETE]
func DeleteRole(ctx context.Context, c *app.RequestContext) {
	roleID := c.Param("role_id")
	if roleID == "" {
		invalidParamRequestResponse(c, "role_id is required")
		return
	}

	err := permissionapp.PermissionAppSVC.DeleteRole(ctx, roleID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// ListRoles 列出角色
// @router /api/roles [GET]
func ListRoles(ctx context.Context, c *app.RequestContext) {
	var req permission.ListRolesRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.ListRoles(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 用户角色分配接口 ====================

// AssignRole 分配角色
// @router /api/users/:user_id/roles [POST]
func AssignRole(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	var err error
	var req permission.AssignRoleRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.AssignRole(ctx, userID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.AssignRoleResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// RevokeRole 撤销角色
// @router /api/users/:user_id/roles/:role_id [DELETE]
func RevokeRole(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	roleID := c.Param("role_id")
	if roleID == "" {
		invalidParamRequestResponse(c, "role_id is required")
		return
	}

	err := permissionapp.PermissionAppSVC.RevokeRole(ctx, userID, tenantID, roleID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// GetUserRoles 获取用户的所有角色
// @router /api/users/:user_id/roles [GET]
func GetUserRoles(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.GetUserRolesResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// ==================== 数据权限接口 ====================

// SetDataPermission 设置角色数据权限
// @router /api/roles/:role_id/data-permissions [POST]
func SetDataPermission(ctx context.Context, c *app.RequestContext) {
	roleID := c.Param("role_id")
	if roleID == "" {
		invalidParamRequestResponse(c, "role_id is required")
		return
	}

	var err error
	var req permission.SetDataPermissionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = permissionapp.PermissionAppSVC.SetDataPermission(ctx, roleID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.SetDataPermissionResponse{
		Code:    0,
		Message: "success",
	})
}

// CheckDataPermission 检查数据权限
// @router /api/permissions/check-data [POST]
func CheckDataPermission(ctx context.Context, c *app.RequestContext) {
	var err error
	var req permission.CheckDataPermissionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.CheckDataPermission(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.CheckDataPermissionResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// ==================== 字段权限接口 ====================

// SetFieldPermission 设置角色字段权限
// @router /api/roles/:role_id/field-permissions [POST]
func SetFieldPermission(ctx context.Context, c *app.RequestContext) {
	roleID := c.Param("role_id")
	if roleID == "" {
		invalidParamRequestResponse(c, "role_id is required")
		return
	}

	var err error
	var req permission.SetFieldPermissionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = permissionapp.PermissionAppSVC.SetFieldPermission(ctx, roleID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.SetFieldPermissionResponse{
		Code:    0,
		Message: "success",
	})
}

// GetFieldPermissions 获取字段权限
// @router /api/users/:user_id/field-permissions [GET]
func GetFieldPermissions(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resourceType := c.Query("resource_type")
	if resourceType == "" {
		invalidParamRequestResponse(c, "resource_type is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.GetFieldPermissions(ctx, userID, tenantID, resourceType)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.GetFieldPermissionsResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// ==================== 部门管理接口 ====================

// CreateDepartment 创建部门
// @router /api/departments [POST]
func CreateDepartment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req permission.CreateDepartmentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.CreateDepartment(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.CreateDepartmentResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// GetDepartment 获取部门
// @router /api/departments/:department_id [GET]
func GetDepartment(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.GetDepartment(ctx, departmentID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.GetDepartmentResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// UpdateDepartment 更新部门
// @router /api/departments/:department_id [PUT]
func UpdateDepartment(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	var err error
	var req permission.UpdateDepartmentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := permissionapp.PermissionAppSVC.UpdateDepartment(ctx, departmentID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.GetDepartmentResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// DeleteDepartment 删除部门
// @router /api/departments/:department_id [DELETE]
func DeleteDepartment(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	err := permissionapp.PermissionAppSVC.DeleteDepartment(ctx, departmentID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// ListDepartments 列出部门
// @router /api/departments [GET]
func ListDepartments(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.ListDepartments(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.ListDepartmentsResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// GetAccessibleDepartments 获取用户可访问的部门
// @router /api/users/:user_id/accessible-departments [GET]
func GetAccessibleDepartments(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.GetAccessibleDepartments(ctx, userID, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.GetAccessibleDepartmentsResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}

// ==================== 部门成员管理接口 ====================

// AddDepartmentMember 添加部门成员
// @router /api/departments/:department_id/members [POST]
func AddDepartmentMember(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	var err error
	var req permission.AddDepartmentMemberRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = permissionapp.PermissionAppSVC.AddDepartmentMember(ctx, departmentID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.AddDepartmentMemberResponse{
		Code:    0,
		Message: "success",
	})
}

// UpdateDepartmentMember 更新部门成员
// @router /api/departments/:department_id/members/:user_id [PUT]
func UpdateDepartmentMember(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	var err error
	var req permission.UpdateDepartmentMemberRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = permissionapp.PermissionAppSVC.UpdateDepartmentMember(ctx, departmentID, userID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.AddDepartmentMemberResponse{
		Code:    0,
		Message: "success",
	})
}

// RemoveDepartmentMember 移除部门成员
// @router /api/departments/:department_id/members/:user_id [DELETE]
func RemoveDepartmentMember(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	userID := c.Param("user_id")
	if userID == "" {
		invalidParamRequestResponse(c, "user_id is required")
		return
	}

	err := permissionapp.PermissionAppSVC.RemoveDepartmentMember(ctx, departmentID, userID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// GetDepartmentMembers 获取部门成员
// @router /api/departments/:department_id/members [GET]
func GetDepartmentMembers(ctx context.Context, c *app.RequestContext) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		invalidParamRequestResponse(c, "department_id is required")
		return
	}

	resp, err := permissionapp.PermissionAppSVC.GetDepartmentMembers(ctx, departmentID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &permission.GetDepartmentMembersResponse{
		Code:    0,
		Message: "success",
		Data:    *resp,
	})
}
