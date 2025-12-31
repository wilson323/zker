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

// Package developer 提供开发者平台 API 路由注册
package developer

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/coze-dev/coze-studio/backend/api/handler/developer"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
)

// Handlers 开发者平台 Handler 实例
var (
	projectHandler  *developer.ProjectManagementHandler
	apiKeyHandler   *developer.APIKeyManagementHandler
	sdkHandler      *developer.SDKGeneratorHandler
	webhookHandler  *developer.WebhookHandler
	userHandler     *developer.UserManagementHandler
)

// InitHandlers 初始化开发平台 Handler
// 需要在应用启动时调用，注入依赖
func InitHandlers(
	projectSvc interface{}, // devservice.ProjectManagementService
	apiKeySvc interface{},  // devservice.APIKeyManagementService
	sdkSvc interface{},     // devservice.SDKGeneratorService
	webhookSvc interface{}, // devservice.WebhookService
	userSvc interface{},    // *permservice.UserService
) {
	projectHandler = developer.NewProjectManagementHandler(projectSvc)
	apiKeyHandler = developer.NewAPIKeyManagementHandler(apiKeySvc)
	sdkHandler = developer.NewSDKGeneratorHandler(sdkSvc)
	webhookHandler = developer.NewWebhookHandler(webhookSvc)
	userHandler = developer.NewUserManagementHandler(userSvc)
}

// RegisterRoutes 注册开发者平台路由
//
// **API 前缀**: /api/developer
//
// **路由结构**:
//   /api/developer/projects       - 项目管理
//   /api/developer/api-keys       - API密钥管理
//   /api/developer/sdks           - SDK生成
//   /api/developer/webhooks       - Webhook管理
func RegisterRoutes(r *server.Hertz) {
	// 基础中间件：租户隔离 + 权限检查
	baseMiddleware := []app.HandlerFunc{
		middleware.TenantIsolationMiddleware(),
	}

	// 开发者平台路由组
	developerGroup := r.Group("/api/developer", baseMiddleware...)

	// ==================== 项目管理路由 ====================
	projectsGroup := developerGroup.Group("/projects")
	{
		// 项目列表
		projectsGroup.GET("", any(projectHandler.ListProjects))
		// 创建项目
		projectsGroup.POST("", any(projectHandler.CreateProject))
		// 获取项目详情
		projectsGroup.GET("/:id", any(projectHandler.GetProject))
		// 更新项目
		projectsGroup.PUT("/:id", any(projectHandler.UpdateProject))
		// 删除项目
		projectsGroup.DELETE("/:id", any(projectHandler.DeleteProject))
		// 发布项目
		projectsGroup.POST("/:id/publish", any(projectHandler.PublishProject))
		// 归档项目
		projectsGroup.POST("/:id/archive", any(projectHandler.ArchiveProject))
		// 获取项目统计
		projectsGroup.GET("/:id/stats", any(projectHandler.GetProjectStats))
	}

	// 开发者项目路由
	developersGroup := developerGroup.Group("/developers")
	{
		// 获取开发者的项目列表
		developersGroup.GET("/:id/projects", any(projectHandler.GetDeveloperProjects))
	}

	// 租户项目路由（内部接口）
	tenantGroup := developerGroup.Group("/tenant")
	{
		// 获取租户的项目列表
		tenantGroup.GET("/projects", any(projectHandler.GetTenantProjects))
	}

	// ==================== API密钥管理路由 ====================
	apiKeysGroup := developerGroup.Group("/api-keys")
	{
		// API密钥列表
		apiKeysGroup.GET("", any(apiKeyHandler.ListAPIKeys))
		// 创建API密钥
		apiKeysGroup.POST("", any(apiKeyHandler.CreateAPIKey))
		// 获取API密钥详情
		apiKeysGroup.GET("/:id", any(apiKeyHandler.GetAPIKey))
		// 更新API密钥
		apiKeysGroup.PUT("/:id", any(apiKeyHandler.UpdateAPIKey))
		// 删除API密钥
		apiKeysGroup.DELETE("/:id", any(apiKeyHandler.DeleteAPIKey))
		// 撤销API密钥
		apiKeysGroup.POST("/:id/revoke", any(apiKeyHandler.RevokeAPIKey))
		// 重新生成API密钥
		apiKeysGroup.POST("/:id/regenerate", any(apiKeyHandler.RegenerateAPIKey))
		// 获取即将过期的密钥
		apiKeysGroup.GET("/expiring", any(apiKeyHandler.GetExpiringKeys))
	}

	// 项目API密钥路由
	projectsAPIKeysGroup := developerGroup.Group("/projects")
	{
		// 获取项目的API密钥列表
		projectsAPIKeysGroup.GET("/:id/api-keys", any(apiKeyHandler.GetProjectAPIKeys))
	}

	// ==================== SDK生成路由 ====================
	sdksGroup := developerGroup.Group("/sdks")
	{
		// SDK列表
		sdksGroup.GET("", any(sdkHandler.ListSDKs))
		// 生成SDK
		sdksGroup.POST("/generate", any(sdkHandler.GenerateSDK))
		// 获取支持的SDK语言列表
		sdksGroup.GET("/languages", any(sdkHandler.GetSupportedLanguages))
	}

	// SDK详情路由
	sdkDetailGroup := developerGroup.Group("/sdks")
	{
		// 获取SDK详情
		sdkDetailGroup.GET("/:id", any(sdkHandler.GetSDK))
		// 删除SDK
		sdkDetailGroup.DELETE("/:id", any(sdkHandler.DeleteSDK))
		// 发布SDK
		sdkDetailGroup.POST("/:id/publish", any(sdkHandler.PublishSDK))
		// 下载SDK
		sdkDetailGroup.GET("/:id/download", any(sdkHandler.DownloadSDK))
		// 获取SDK代码预览
		sdkDetailGroup.GET("/:id/code", any(sdkHandler.GetSDKCode))
	}

	// 项目SDK路由
	projectsSDKsGroup := developerGroup.Group("/projects")
	{
		// 获取项目的SDK列表
		projectsSDKsGroup.GET("/:id/sdks", any(sdkHandler.GetProjectSDKs))
	}

	// ==================== Webhook管理路由 ====================
	webhooksGroup := developerGroup.Group("/webhooks")
	{
		// Webhook列表
		webhooksGroup.GET("", any(webhookHandler.ListWebhooks))
		// 创建Webhook
		webhooksGroup.POST("", any(webhookHandler.CreateWebhook))
		// 获取支持的事件类型列表
		webhooksGroup.GET("/events", any(webhookHandler.GetSupportedEvents))
	}

	// Webhook详情路由
	webhookDetailGroup := developerGroup.Group("/webhooks")
	{
		// 获取Webhook详情
		webhookDetailGroup.GET("/:id", any(webhookHandler.GetWebhook))
		// 更新Webhook
		webhookDetailGroup.PUT("/:id", any(webhookHandler.UpdateWebhook))
		// 删除Webhook
		webhookDetailGroup.DELETE("/:id", any(webhookHandler.DeleteWebhook))
		// 测试Webhook
		webhookDetailGroup.POST("/:id/test", any(webhookHandler.TestWebhook))
		// 暂停Webhook
		webhookDetailGroup.POST("/:id/pause", any(webhookHandler.PauseWebhook))
		// 恢复Webhook
		webhookDetailGroup.POST("/:id/resume", any(webhookHandler.ResumeWebhook))
		// 获取Webhook统计信息
		webhookDetailGroup.GET("/:id/stats", any(webhookHandler.GetWebhookStats))
	}

	// 项目Webhook路由
	projectsWebhooksGroup := developerGroup.Group("/projects")
	{
		// 获取项目的Webhook列表
		projectsWebhooksGroup.GET("/:id/webhooks", any(webhookHandler.GetProjectWebhooks))
	}

	// ==================== 用户管理路由 ====================
	usersGroup := developerGroup.Group("/users")
	{
		// 用户列表
		usersGroup.GET("", any(userHandler.ListUsers))
		// 创建用户
		usersGroup.POST("", any(userHandler.CreateUser))
		// 导出用户列表
		usersGroup.GET("/export", any(userHandler.ExportUsers))
		// 批量删除用户
		usersGroup.POST("/batch-delete", any(userHandler.BatchDeleteUsers))
		// 批量更新用户
		usersGroup.POST("/batch-update", any(userHandler.BatchUpdateUsers))
	}

	// 用户详情路由
	userDetailGroup := developerGroup.Group("/users")
	{
		// 获取用户详情
		userDetailGroup.GET("/:id", any(userHandler.GetUser))
		// 更新用户信息
		userDetailGroup.PUT("/:id", any(userHandler.UpdateUser))
		// 删除用户
		userDetailGroup.DELETE("/:id", any(userHandler.DeleteUser))
		// 更新用户状态
		userDetailGroup.PUT("/:id/status", any(userHandler.UpdateUserStatus))
		// 重置用户密码
		userDetailGroup.POST("/:id/reset-password", any(userHandler.ResetPassword))
		// 解锁用户账户
		userDetailGroup.POST("/:id/unlock", any(userHandler.UnlockUser))
	}

	// 用户角色路由
	userRolesGroup := developerGroup.Group("/users")
	{
		// 获取用户所有角色
		userRolesGroup.GET("/:id/roles", any(userHandler.GetUserRoles))
		// 分配角色给用户
		userRolesGroup.POST("/:id/roles", any(userHandler.AssignRoles))
		// 批量更新用户角色
		userRolesGroup.PUT("/:id/roles", any(userHandler.UpdateUserRoles))
		// 移除用户角色
		userRolesGroup.DELETE("/:id/roles/:role_id", any(userHandler.RemoveRole))
	}

	// 用户权限路由
	userPermissionsGroup := developerGroup.Group("/users")
	{
		// 获取用户所有权限
		userPermissionsGroup.GET("/:id/permissions", any(userHandler.GetUserPermissions))
		// 获取用户有效权限（含继承）
		userPermissionsGroup.GET("/:id/effective-permissions", any(userHandler.GetUserEffectivePermissions))
	}
}

// any 将 handler 转换为 app.HandlersFunc
// 用于解决类型兼容性问题
func any(handler interface{}) app.HandlerFunc {
	switch h := handler.(type) {
	case func(context.Context, *app.RequestContext):
		return h
	default:
		panic("unsupported handler type")
	}
}

// GetProjectHandler 获取项目Handler（用于测试）
func GetProjectHandler() *developer.ProjectManagementHandler {
	return projectHandler
}

// GetAPIKeyHandler 获取API密钥Handler（用于测试）
func GetAPIKeyHandler() *developer.APIKeyManagementHandler {
	return apiKeyHandler
}

// GetSDKHandler 获取SDK Handler（用于测试）
func GetSDKHandler() *developer.SDKGeneratorHandler {
	return sdkHandler
}

// GetWebhookHandler 获取Webhook Handler（用于测试）
func GetWebhookHandler() *developer.WebhookHandler {
	return webhookHandler
}

// GetUserHandler 获取用户管理Handler（用于测试）
func GetUserHandler() *developer.UserManagementHandler {
	return userHandler
}
