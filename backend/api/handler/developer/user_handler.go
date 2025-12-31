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

package developer

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	permservice "github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// UserManagementHandler 用户管理Handler
type UserManagementHandler struct {
	userService *permservice.UserService
}

// NewUserManagementHandler 创建用户管理Handler实例
func NewUserManagementHandler(userService *permservice.UserService) *UserManagementHandler {
	return &UserManagementHandler{
		userService: userService,
	}
}

// ========== 请求/响应结构 ==========

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email,max=255"`
	UniqueName  string `json:"unique_name" validate:"required,min=2,max=50"`
	Description string `json:"description" validate:"omitempty,max=500"`
	IconURI     string `json:"icon_uri" validate:"omitempty,url,max=500"`
	Locale      string `json:"locale" validate:"omitempty,max=10"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	IconURI     *string `json:"icon_uri" validate:"omitempty,url,max=500"`
	Locale      *string `json:"locale" validate:"omitempty,max=10"`
}

// ListUsersRequest 获取用户列表请求
type ListUsersRequest struct {
	Keyword   string `json:"keyword" query:"keyword"`      // 搜索关键词（姓名、邮箱）
	RoleID    string `json:"role_id" query:"role_id"`      // 按角色筛选
	Status    string `json:"status" query:"status"`        // 按状态筛选（active/inactive）
	PageToken string `json:"page_token" query:"page_token"`
	PageSize  int    `json:"page_size" query:"page_size"`
}

// ListUsersResponse 获取用户列表响应
type ListUsersResponse struct {
	Users     []*UserListItem `json:"users"`
	Total     int64           `json:"total"`
	PageToken string          `json:"page_token,omitempty"`
}

// UserListItem 用户列表项
type UserListItem struct {
	UserID      int64  `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	UniqueName  string `json:"unique_name"`
	Description string `json:"description"`
	IconURI     string `json:"icon_uri"`
	Locale      string `json:"locale"`
	IsEnabled   bool   `json:"is_enabled"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

// UserDetailResponse 用户详情响应
type UserDetailResponse struct {
	UserID         int64                `json:"user_id"`
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	UniqueName     string               `json:"unique_name"`
	Description    string               `json:"description"`
	IconURI        string               `json:"icon_uri"`
	Locale         string               `json:"locale"`
	IsEnabled      bool                 `json:"is_enabled"`
	UserVerified   bool                 `json:"user_verified"`
	CreatedAt      int64                `json:"created_at"`
	UpdatedAt      int64                `json:"updated_at"`
	Roles          []*RoleSummary       `json:"roles"`
	Permissions    []*PermissionSummary `json:"permissions"`
}

// RoleSummary 角色摘要
type RoleSummary struct {
	RoleID      string `json:"role_id"`
	RoleName    string `json:"role_name"`
	RoleCode    string `json:"role_code"`
	RoleType    string `json:"role_type"`
	Description string `json:"description"`
}

// PermissionSummary 权限摘要
type PermissionSummary struct {
	ResourceType string `json:"resource_type"`
	Actions      []string `json:"actions"`
	Scope        string `json:"scope"`
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids" validate:"required,min=1"`
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	IsEnabled bool `json:"is_enabled"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

// BatchDeleteUsersRequest 批量删除用户请求
type BatchDeleteUsersRequest struct {
	UserIDs []int64 `json:"user_ids" validate:"required,min=1,max=100"`
}

// BatchUpdateUsersRequest 批量更新用户请求
type BatchUpdateUsersRequest struct {
	UserIDs    []int64  `json:"user_ids" validate:"required,min=1,max=100"`
	IsEnabled  *bool    `json:"is_enabled,omitempty"`
	RoleIDs    []string `json:"role_ids,omitempty"`
}

// ========== 1. 用户基础CRUD ==========

// CreateUser 创建用户
// @router /api/developer/users [POST]
func (h *UserManagementHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req CreateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 构建创建请求
	createReq := &permservice.CreateUserRequest{
		TenantID:    tenantID,
		Name:        req.Name,
		Email:       req.Email,
		UniqueName:  req.UniqueName,
		Description: req.Description,
		IconURI:     req.IconURI,
		Locale:      req.Locale,
	}

	// 创建用户
	user, err := h.userService.CreateUser(ctx, createReq)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, user)
}

// ListUsers 获取用户列表（支持分页、搜索）
// @router /api/developer/users [GET]
func (h *UserManagementHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req ListUsersRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分页参数
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建过滤器
	filter := &permservice.UserListFilter{
		TenantID:  tenantID,
		Keyword:   req.Keyword,
		RoleID:    req.RoleID,
		Status:    req.Status,
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	// 查询用户列表
	users, total, err := h.userService.ListUsers(ctx, filter)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	items := make([]*UserListItem, 0, len(users))
	for _, u := range users {
		items = append(items, &UserListItem{
			UserID:      u.UserID,
			Name:        u.Name,
			Email:       u.Email,
			UniqueName:  u.UniqueName,
			Description: u.Description,
			IconURI:     u.IconURI,
			Locale:      u.Locale,
			IsEnabled:   u.IsEnabled,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		})
	}

	httputil.BuildSuccessResp(c, &ListUsersResponse{
		Users:     items,
		Total:     total,
		PageToken: req.PageToken,
	})
}

// GetUser 获取用户详情
// @router /api/developer/users/:id [GET]
func (h *UserManagementHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询用户
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 验证租户隔离
	if user == nil || user.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrUserNotFound.Int32Code(), berrno.ErrUserNotFound.Message(), berrno.ErrUserNotFound.MessageZH(), nil)
		return
	}

	// 获取用户角色
	roles, _ := h.userService.GetUserRoles(ctx, userID, tenantID)

	// 获取用户权限
	permissions, _ := h.userService.GetUserPermissions(ctx, userID, tenantID)

	// 转换响应
	roleSummaries := make([]*RoleSummary, 0, len(roles))
	for _, r := range roles {
		roleSummaries = append(roleSummaries, &RoleSummary{
			RoleID:      r.RoleID,
			RoleName:    r.RoleName,
			RoleCode:    r.RoleCode,
			RoleType:    string(r.RoleType),
			Description: r.Description,
		})
	}

	permissionSummaries := make([]*PermissionSummary, 0, len(permissions))
	for _, p := range permissions {
		permissionSummaries = append(permissionSummaries, &PermissionSummary{
			ResourceType: string(p.ResourceType),
			Actions:      p.Actions,
			Scope:        string(p.Scope),
		})
	}

	httputil.BuildSuccessResp(c, &UserDetailResponse{
		UserID:         user.UserID,
		Name:           user.Name,
		Email:          user.Email,
		UniqueName:     user.UniqueName,
		Description:    user.Description,
		IconURI:        user.IconURI,
		Locale:         user.Locale,
		IsEnabled:      user.IsEnabled,
		UserVerified:   user.UserVerified,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Roles:          roleSummaries,
		Permissions:    permissionSummaries,
	})
}

// UpdateUser 更新用户信息
// @router /api/developer/users/:id [PUT]
func (h *UserManagementHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证用户存在且属于该租户
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if user == nil || user.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrUserNotFound.Int32Code(), berrno.ErrUserNotFound.Message(), berrno.ErrUserNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req UpdateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 构建更新请求
	updateReq := &permservice.UpdateUserRequest{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IconURI:     req.IconURI,
		Locale:      req.Locale,
	}

	// 更新用户
	if err := h.userService.UpdateUser(ctx, updateReq); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":  userID,
		"message": "User updated successfully",
	})
}

// DeleteUser 删除用户（软删除）
// @router /api/developer/users/:id [DELETE]
func (h *UserManagementHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证用户存在且属于该租户
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if user == nil || user.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrUserNotFound.Int32Code(), berrno.ErrUserNotFound.Message(), berrno.ErrUserNotFound.MessageZH(), nil)
		return
	}

	// 删除用户
	if err := h.userService.DeleteUser(ctx, userID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "User deleted successfully",
	})
}

// ========== 2. 用户角色管理 ==========

// AssignRoles 分配角色给用户
// @router /api/developer/users/:id/roles [POST]
func (h *UserManagementHandler) AssignRoles(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req AssignRolesRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 分配角色
	if err := h.userService.AssignRoles(ctx, userID, tenantID, req.RoleIDs); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":  userID,
		"role_ids": req.RoleIDs,
		"message":  "Roles assigned successfully",
	})
}

// RemoveRole 移除用户角色
// @router /api/developer/users/:id/roles/:role_id [DELETE]
func (h *UserManagementHandler) RemoveRole(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	roleID := c.Param("role_id")
	if roleID == "" {
		httputil.BadRequest(c, "role_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 移除角色
	if err := h.userService.RemoveRole(ctx, userID, tenantID, roleID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":  userID,
		"role_id":  roleID,
		"message":  "Role removed successfully",
	})
}

// GetUserRoles 获取用户所有角色
// @router /api/developer/users/:id/roles [GET]
func (h *UserManagementHandler) GetUserRoles(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 获取用户角色
	roles, err := h.userService.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, roles)
}

// UpdateUserRoles 批量更新用户角色
// @router /api/developer/users/:id/roles [PUT]
func (h *UserManagementHandler) UpdateUserRoles(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req AssignRolesRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 更新用户角色（先清除所有角色，再分配新角色）
	if err := h.userService.UpdateUserRoles(ctx, userID, tenantID, req.RoleIDs); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":  userID,
		"role_ids": req.RoleIDs,
		"message":  "User roles updated successfully",
	})
}

// ========== 3. 用户权限查询 ==========

// GetUserPermissions 获取用户所有权限
// @router /api/developer/users/:id/permissions [GET]
func (h *UserManagementHandler) GetUserPermissions(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 获取用户权限
	permissions, err := h.userService.GetUserPermissions(ctx, userID, tenantID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, permissions)
}

// GetUserEffectivePermissions 获取用户有效权限（含继承）
// @router /api/developer/users/:id/effective-permissions [GET]
func (h *UserManagementHandler) GetUserEffectivePermissions(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 获取用户有效权限
	permissions, err := h.userService.GetUserEffectivePermissions(ctx, userID, tenantID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, permissions)
}

// ========== 4. 用户状态管理 ==========

// UpdateUserStatus 启用/禁用用户
// @router /api/developer/users/:id/status [PUT]
func (h *UserManagementHandler) UpdateUserStatus(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证用户存在且属于该租户
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if user == nil || user.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrUserNotFound.Int32Code(), berrno.ErrUserNotFound.Message(), berrno.ErrUserNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req UpdateUserStatusRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 更新用户状态
	if err := h.userService.UpdateUserStatus(ctx, userID, req.IsEnabled); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":    userID,
		"is_enabled": req.IsEnabled,
		"message":    "User status updated successfully",
	})
}

// ResetPassword 重置用户密码
// @router /api/developer/users/:id/reset-password [POST]
func (h *UserManagementHandler) ResetPassword(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证用户存在且属于该租户
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if user == nil || user.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrUserNotFound.Int32Code(), berrno.ErrUserNotFound.Message(), berrno.ErrUserNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req ResetPasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 重置密码
	if err := h.userService.ResetPassword(ctx, userID, req.NewPassword); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":  userID,
		"message":  "Password reset successfully",
	})
}

// UnlockUser 解锁用户账户
// @router /api/developer/users/:id/unlock [POST]
func (h *UserManagementHandler) UnlockUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		httputil.BadRequest(c, "user_id is required")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		httputil.BadRequest(c, "invalid user_id")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证用户存在且属于该租户
	user, err := h.userService.GetUser(ctx, userID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if user == nil || user.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrUserNotFound.Int32Code(), berrno.ErrUserNotFound.Message(), berrno.ErrUserNotFound.MessageZH(), nil)
		return
	}

	// 解锁用户
	if err := h.userService.UnlockUser(ctx, userID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"user_id":  userID,
		"message":  "User unlocked successfully",
	})
}

// ========== 5. 用户批量操作 ==========

// BatchDeleteUsers 批量删除用户
// @router /api/developer/users/batch-delete [POST]
func (h *UserManagementHandler) BatchDeleteUsers(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req BatchDeleteUsersRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 批量删除用户
	results, err := h.userService.BatchDeleteUsers(ctx, tenantID, req.UserIDs)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, results)
}

// BatchUpdateUsers 批量更新用户
// @router /api/developer/users/batch-update [POST]
func (h *UserManagementHandler) BatchUpdateUsers(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req BatchUpdateUsersRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 批量更新用户
	results, err := h.userService.BatchUpdateUsers(ctx, tenantID, req.UserIDs, req.IsEnabled, req.RoleIDs)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, results)
}

// ExportUsers 导出用户列表
// @router /api/developer/users/export [GET]
func (h *UserManagementHandler) ExportUsers(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req ListUsersRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 构建过滤器
	filter := &permservice.UserListFilter{
		TenantID: tenantID,
		Keyword:  req.Keyword,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}

	// 导出用户列表
	fileURL, err := h.userService.ExportUsers(ctx, filter)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"file_url": fileURL,
		"message":  "Users exported successfully",
	})
}
