// Package middleware 提供HTTP中间件
//
// 本文件展示权限检查中间件的实现，包括：
// - RBAC权限验证
// - 数据权限过滤
// - 临时授权检查
// - 权限缓存优化
package middleware

import (
	"fmt"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-studio/crossdomain/permission"
)

// PermissionMiddleware 权限检查中间件
type PermissionMiddleware struct {
	permissionSvc *permission.PermissionService
}

// NewPermissionMiddleware 创建权限检查中间件
func NewPermissionMiddleware(
	permissionSvc *permission.PermissionService,
) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionSvc: permissionSvc,
	}
}

// RequirePermission 要求特定权限
//
// 用法示例：
//
//	router.POST("/api/v1/bots",
//	    authMiddleware.Middleware(),
//	    permissionMiddleware.RequirePermission("bot.create"),
//	    botHandler.CreateBot,
//	)
func (m *PermissionMiddleware) RequirePermission(
	permission string,
) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		// 步骤1: 获取用户上下文（由认证中间件设置）
		userContext := GetUserContext(ctx)
		if userContext == nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 步骤2: 检查权限
		hasPermission, err := m.permissionSvc.CheckPermission(
			ctx,
			userContext,
			permission,
		)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "failed to check permission",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		if !hasPermission {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": fmt.Sprintf("permission denied: %s", permission),
			})
			ctx.Abort()
			return
		}

		// 继续处理请求
		ctx.Next(ctx)
	}
}

// RequireRole 要求特定角色
//
// 用法示例：
//
//	router.DELETE("/api/v1/users/:id",
//	    authMiddleware.Middleware(),
//	    permissionMiddleware.RequireRole("admin"),
//	    userHandler.DeleteUser,
//	)
func (m *PermissionMiddleware) RequireRole(
	role string,
) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		userContext := GetUserContext(ctx)
		if userContext == nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 检查角色
		hasRole := false
		for _, r := range userContext.Roles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": fmt.Sprintf("role required: %s", role),
			})
			ctx.Abort()
			return
		}

		ctx.Next(ctx)
	}
}

// RequireAnyRole 要求任一角色
//
// 用法示例：
//
//	router.GET("/api/v1/admin/dashboard",
//	    authMiddleware.Middleware(),
//	    permissionMiddleware.RequireAnyRole("admin", "super_admin"),
//	    adminHandler.Dashboard,
//	)
func (m *PermissionMiddleware) RequireAnyRole(
	roles ...string,
) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		userContext := GetUserContext(ctx)
		if userContext == nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 检查是否有任一角色
		hasAnyRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userContext.Roles {
				if userRole == requiredRole {
					hasAnyRole = true
					break
				}
			}
			if hasAnyRole {
				break
			}
		}

		if !hasAnyRole {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": fmt.Sprintf("one of roles required: %s", strings.Join(roles, ", ")),
			})
			ctx.Abort()
			return
		}

		ctx.Next(ctx)
	}
}

// RequireOwnership 要求资源所有权或管理员权限
//
// 用法示例：
//
//	router.PUT("/api/v1/bots/:id",
//	    authMiddleware.Middleware(),
//	    permissionMiddleware.RequireOwnership("bot"),
//	    botHandler.UpdateBot,
//	)
func (m *PermissionMiddleware) RequireOwnership(
	resourceType string,
) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		userContext := GetUserContext(ctx)
		if userContext == nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 步骤1: 检查是否是管理员
		isAdmin := false
		for _, role := range userContext.Roles {
			if role == "admin" || role == "super_admin" {
				isAdmin = true
				break
			}
		}

		if isAdmin {
			// 管理员跳过所有权检查
			ctx.Next(ctx)
			return
		}

		// 步骤2: 提取资源ID
		resourceID := ctx.Param("id")
		if resourceID == "" {
			ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
				"code":    consts.StatusBadRequest,
				"message": "resource ID is required",
			})
			ctx.Abort()
			return
		}

		// 步骤3: 检查所有权
		isOwner, err := m.permissionSvc.CheckResourceOwnership(
			ctx,
			userContext,
			resourceType,
			resourceID,
		)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "failed to check ownership",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		if !isOwner {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": "permission denied: not the owner of this resource",
			})
			ctx.Abort()
			return
		}

		ctx.Next(ctx)
	}
}

// RequireDataLevel 要求特定数据级别权限
//
// 数据级别定义：
// - level_1: 仅本人数据
// - level_2: 本部门数据
// - level_3: 本部门及子部门数据
// - level_4: 全租户数据
// - level_5: 全平台数据（仅超级管理员）
//
// 用法示例：
//
//	router.GET("/api/v1/users",
//	    authMiddleware.Middleware(),
//	    permissionMiddleware.RequireDataLevel("level_3"),
//	    userHandler.ListUsers,
//	)
func (m *PermissionMiddleware) RequireDataLevel(
	requiredLevel string,
) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		userContext := GetUserContext(ctx)
		if userContext == nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 检查数据级别权限
		hasPermission, err := m.permissionSvc.CheckDataLevelPermission(
			ctx,
			userContext,
			requiredLevel,
		)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "failed to check data level permission",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		if !hasPermission {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": fmt.Sprintf("data level permission required: %s", requiredLevel),
			})
			ctx.Abort()
			return
		}

		ctx.Next(ctx)
	}
}

// RequireFieldPermission 要求字段级权限
//
// 用法示例（在处理器内部使用）：
//
//	func (h *UserHandler) UpdateUser(ctx *app.RequestContext) {
//	    userContext := middleware.GetUserContext(ctx)
//
//	    // 检查是否有修改薪资字段的权限
//	    if !middleware.CheckFieldPermission(userContext, "user", "salary") {
//	        // 移除薪资字段
//	        delete(req, "salary")
//	    }
//
//	    // 继续处理...
//	}
func (m *PermissionMiddleware) RequireFieldPermission(
	resourceType string,
	field string,
) app.HandlerFunc {
	return func(ctx *app.RequestContext) {
		userContext := GetUserContext(ctx)
		if userContext == nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    consts.StatusUnauthorized,
				"message": "unauthorized",
			})
			ctx.Abort()
			return
		}

		// 检查字段权限
		hasPermission, err := m.permissionSvc.CheckFieldPermission(
			ctx,
			userContext,
			resourceType,
			field,
		)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "failed to check field permission",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		if !hasPermission {
			ctx.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": fmt.Sprintf("field permission denied: %s.%s", resourceType, field),
			})
			ctx.Abort()
			return
		}

		ctx.Next(ctx)
	}
}
