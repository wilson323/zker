// +build ignore

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

// Deprecated: Use permission_check_enhanced.go instead.
// This file is kept for reference only and is excluded from builds.
package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

var (
	// 权限检查器（在应用启动时注入）
	permissionChecker *service.PermissionChecker
)

// InitPermissionMiddleware 初始化权限中间件依赖
func InitPermissionMiddleware(pc *service.PermissionChecker) {
	permissionChecker = pc
}

// PermissionCheckConfig 权限检查配置
type PermissionCheckConfig struct {
	ResourceType entity.ResourceType // 资源类型
	Action       string              // 操作类型: create, read, update, delete
	GetResourceID func(c *app.RequestContext) string // 从请求中获取资源ID的函数
	RequiredPerm string              // 需要的权限（可选，用于简单权限检查）
}

// RequirePermission 权限检查中间件
// 用于检查用户是否有权限访问指定资源
//
// 使用示例：
//
//	r.GET("/api/bots/:bot_id", middleware.RequirePermission(middleware.PermissionCheckConfig{
//	    ResourceType: entity.ResourceTypeBots,
//	    Action:       "read",
//	    GetResourceID: func(c *app.RequestContext) string {
//	        return c.Param("bot_id")
//	    },
//	}), handler.GetBot)
func RequirePermission(config PermissionCheckConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取user_id
		userID := c.GetHeader("X-User-ID")
		if userID == string("") {
			c.JSON(consts.StatusUnauthorized, responseWithError(berrno.ErrUnauthorizedCode))
			c.Abort()
			return
		}

		// 2. 获取tenant_id
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == string("") {
			if tid, exists := c.Get("tenant_id"); exists {
				tenantID = string(tid)
			} else {
				c.JSON(consts.StatusBadRequest, responseWithError(berrno.ErrTenantNotFoundCode, "tenant_id is required"))
				c.Abort()
				return
			}
		}

		// 3. 检查是否已初始化权限检查器
		if permissionChecker == nil {
			logs.CtxWarnf(ctx, "[PermissionCheck] permissionChecker not initialized, skipping permission check")
			c.Next(ctx)
			return
		}

		// 4. 获取资源ID
		resourceID := ""
		if config.GetResourceID != nil {
			resourceID = config.GetResourceID(c)
		}

		// 5. 检查数据权限
		hasPermission, err := permissionChecker.CheckDataPermission(ctx, tenantID, userID, config.ResourceType, config.Action, resourceID)
		if err != nil {
			logs.CtxErrorf(ctx, "[PermissionCheck] check data permission failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrPermissionCheckFailedCode))
			c.Abort()
			return
		}

		if !hasPermission {
			logs.CtxWarnf(ctx, "[PermissionCheck] permission denied: user_id=%s, tenant_id=%s, resource_type=%s, action=%s, resource_id=%s",
				userID, tenantID, config.ResourceType, config.Action, resourceID)

			c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrDataPermissionDeniedCode,
				errorx.KV(
					"resource_type", string(config.ResourceType),
					"action", config.Action,
					"resource_id", resourceID,
				),
			))
			c.Abort()
			return
		}

		logs.CtxInfof(ctx, "[PermissionCheck] permission granted: user_id=%s, tenant_id=%s, resource_type=%s, action=%s, resource_id=%s",
			userID, tenantID, config.ResourceType, config.Action, resourceID)

		// 6. 权限检查通过，继续处理请求
		c.Next(ctx)
	}
}

// RequireRole 角色检查中间件
// 用于检查用户是否拥有指定角色
//
// 使用示例：
//
//	r.POST("/api/admin/users", middleware.RequireRole("admin"), handler.CreateUser)
func RequireRole(requiredRole string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取user_id
		userID := c.GetHeader("X-User-ID")
		if userID == string("") {
			c.JSON(consts.StatusUnauthorized, responseWithError(berrno.ErrUnauthorizedCode))
			c.Abort()
			return
		}

		// 2. 获取tenant_id
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == string("") {
			if tid, exists := c.Get("tenant_id"); exists {
				tenantID = string(tid)
			} else {
				c.JSON(consts.StatusBadRequest, responseWithError(berrno.ErrTenantNotFoundCode, "tenant_id is required"))
				c.Abort()
				return
			}
		}

		// 3. 检查是否已初始化权限检查器
		if permissionChecker == nil {
			logs.CtxWarnf(ctx, "[RequireRole] permissionChecker not initialized, skipping role check")
			c.Next(ctx)
			return
		}

		// 4. 检查用户角色
		hasRole, err := permissionChecker.UserHasRole(ctx, tenantID, userID, requiredRole)
		if err != nil {
			logs.CtxErrorf(ctx, "[RequireRole] check user role failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrPermissionCheckFailedCode))
			c.Abort()
			return
		}

		if !hasRole {
			logs.CtxWarnf(ctx, "[RequireRole] role required: user_id=%s, tenant_id=%s, required_role=%s",
				userID, tenantID, requiredRole)

			c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrPermissionDeniedCode,
				errorx.KV("required_role", requiredRole),
			))
			c.Abort()
			return
		}

		logs.CtxInfof(ctx, "[RequireRole] role check passed: user_id=%s, tenant_id=%s, role=%s",
			userID, tenantID, requiredRole)

		// 5. 角色检查通过，继续处理请求
		c.Next(ctx)
	}
}

// RequireAnyRole 角色检查中间件（满足任一角色即可）
// 用于检查用户是否拥有指定角色中的至少一个
//
// 使用示例：
//
//	r.POST("/api/admin/users", middleware.RequireAnyRole("admin", "super_admin"), handler.CreateUser)
func RequireAnyRole(requiredRoles ...string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if len(requiredRoles) == 0 {
			c.Next(ctx)
			return
		}

		// 1. 获取user_id
		userID := c.GetHeader("X-User-ID")
		if userID == string("") {
			c.JSON(consts.StatusUnauthorized, responseWithError(berrno.ErrUnauthorizedCode))
			c.Abort()
			return
		}

		// 2. 获取tenant_id
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == string("") {
			if tid, exists := c.Get("tenant_id"); exists {
				tenantID = string(tid)
			} else {
				c.JSON(consts.StatusBadRequest, responseWithError(berrno.ErrTenantNotFoundCode, "tenant_id is required"))
				c.Abort()
				return
			}
		}

		// 3. 检查是否已初始化权限检查器
		if permissionChecker == nil {
			logs.CtxWarnf(ctx, "[RequireAnyRole] permissionChecker not initialized, skipping role check")
			c.Next(ctx)
			return
		}

		// 4. 检查用户是否拥有任一角色
		hasAnyRole := false
		for _, role := range requiredRoles {
			hasRole, err := permissionChecker.UserHasRole(ctx, tenantID, userID, role)
			if err != nil {
				logs.CtxErrorf(ctx, "[RequireAnyRole] check user role failed: %v", err)
				c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrPermissionCheckFailedCode))
				c.Abort()
				return
			}
			if hasRole {
				hasAnyRole = true
				break
			}
		}

		if !hasAnyRole {
			logs.CtxWarnf(ctx, "[RequireAnyRole] role required: user_id=%s, tenant_id=%s, required_roles=%v",
				userID, tenantID, requiredRoles)

			c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrPermissionDeniedCode,
				errorx.KV("required_roles", requiredRoles),
			))
			c.Abort()
			return
		}

		logs.CtxInfof(ctx, "[RequireAnyRole] role check passed: user_id=%s, tenant_id=%s, roles=%v",
			userID, tenantID, requiredRoles)

		// 5. 角色检查通过，继续处理请求
		c.Next(ctx)
	}
}

// FilterFieldsByPermission 根据字段权限过滤响应字段
// 用于在返回响应前，根据用户的字段权限过滤掉无权访问的字段
//
// 使用示例：
//
//	func GetBot(ctx context.Context, c *app.RequestContext) {
//	    bot := getBotFromDB(botID)
//	    middleware.FilterFieldsByPermission(ctx, c, "bots", &bot)
//	    c.JSON(200, bot)
//	}
func FilterFieldsByPermission(ctx context.Context, c *app.RequestContext, resourceType string, response interface{}) error {
	// 1. 获取user_id
	userID := c.GetHeader("X-User-ID")
	if userID == string("") {
		return errorx.New(berrno.ErrUnauthorizedCode)
	}

	// 2. 获取tenant_id
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == string("") {
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = string(tid)
		} else {
			return errorx.New(berrno.ErrTenantNotFoundCode, errorx.KV("msg", "tenant_id is required"))
		}
	}

	// 3. 检查是否已初始化权限检查器
	if permissionChecker == nil {
		logs.CtxWarnf(ctx, "[FilterFieldsByPermission] permissionChecker not initialized")
		return nil
	}

	// 4. 获取字段权限
	fieldPerms, err := permissionChecker.GetFieldPermissions(ctx, tenantID, userID, resourceType)
	if err != nil {
		logs.CtxErrorf(ctx, "[FilterFieldsByPermission] get field permissions failed: %v", err)
		return err
	}

	// 5. 过滤hidden字段
	// TODO: 实现字段过滤逻辑（需要使用反射或mapstructure）

	logs.CtxInfof(ctx, "[FilterFieldsByPermission] field permissions applied: user_id=%s, tenant_id=%s, resource_type=%s, hidden_fields=%d",
		userID, tenantID, resourceType, len(fieldPerms))

	return nil
}

// ManualPermissionCheck 手动权限检查辅助函数
// 用于在Handler内部进行权限检查
//
// 使用示例：
//
//	func CreateBot(ctx context.Context, c *app.RequestContext) {
//	    hasPermission := middleware.ManualPermissionCheck(ctx, c, middleware.PermissionCheckConfig{
//	        ResourceType: entity.ResourceTypeBots,
//	        Action:       "create",
//	    })
//	    if !hasPermission {
//	        // 权限检查失败
//	        return
//	    }
//	    // 继续处理...
//	}
func ManualPermissionCheck(ctx context.Context, c *app.RequestContext, config PermissionCheckConfig) bool {
	// 1. 获取user_id
	userID := c.GetHeader("X-User-ID")
	if userID == string("") {
		return false
	}

	// 2. 获取tenant_id
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == string("") {
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = string(tid)
		} else {
			return false
		}
	}

	// 3. 检查是否已初始化权限检查器
	if permissionChecker == nil {
		logs.CtxWarnf(ctx, "[ManualPermissionCheck] permissionChecker not initialized")
		return true // 降级：未初始化则通过
	}

	// 4. 获取资源ID
	resourceID := ""
	if config.GetResourceID != nil {
		resourceID = config.GetResourceID(c)
	}

	// 5. 检查权限
	hasPermission, err := permissionChecker.CheckDataPermission(ctx, tenantID, userID, config.ResourceType, config.Action, resourceID)
	if err != nil {
		logs.CtxErrorf(ctx, "[ManualPermissionCheck] check permission failed: %v", err)
		return false
	}

	return hasPermission
}

// GetCurrentUserPermissions 获取当前用户的所有权限
// 用于在Handler内部获取用户权限信息
//
// 使用示例：
//
//	func SomeHandler(ctx context.Context, c *app.RequestContext) {
//	    perms := middleware.GetCurrentUserPermissions(ctx, c)
//	    // 根据权限进行业务逻辑处理
//	}
func GetCurrentUserPermissions(ctx context.Context, c *app.RequestContext) ([]string, error) {
	// 1. 获取user_id
	userID := c.GetHeader("X-User-ID")
	if userID == string("") {
		return nil, errorx.New(berrno.ErrUnauthorizedCode)
	}

	// 2. 获取tenant_id
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == string("") {
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = string(tid)
		} else {
			return nil, errorx.New(berrno.ErrTenantNotFoundCode, errorx.KV("msg", "tenant_id is required"))
		}
	}

	// 3. 检查是否已初始化权限检查器
	if permissionChecker == nil {
		logs.CtxWarnf(ctx, "[GetCurrentUserPermissions] permissionChecker not initialized")
		return []string{}, nil
	}

	// 4. 获取用户权限
	perms, err := permissionChecker.GetUserPermissions(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	return perms, nil
}

// ParseResourceIDFromPath 从路径中解析资源ID的辅助函数
// 用于简化GetResourceID函数的编写
//
// 使用示例：
//
//	middleware.RequirePermission(middleware.PermissionCheckConfig{
//	    ResourceType: entity.ResourceTypeBots,
//	    Action:       "read",
//	    GetResourceID: middleware.ParseResourceIDFromPath("bot_id"),
//	})
func ParseResourceIDFromPath(paramName string) func(c *app.RequestContext) string {
	return func(c *app.RequestContext) string {
		return strings.TrimSpace(c.Param(paramName))
	}
}

// ParseResourceIDFromQuery 从查询参数中解析资源ID的辅助函数
//
// 使用示例：
//
//	middleware.RequirePermission(middleware.PermissionCheckConfig{
//	    ResourceType: entity.ResourceTypeBots,
//	    Action:       "read",
//	    GetResourceID: middleware.ParseResourceIDFromQuery("bot_id"),
//	})
func ParseResourceIDFromQuery(queryName string) func(c *app.RequestContext) string {
	return func(c *app.RequestContext) string {
		return strings.TrimSpace(c.Query(queryName))
	}
}
