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

package middleware

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	"github.com/coze-studio/backend/pkg/errorx"
	"github.com/coze-studio/backend/pkg/loggerx"
	tenantservice "github.com/coze-studio/backend/domain/tenant/service"
	berrno "github.com/coze-studio/backend/types/errno"
)

var (
	// 配额服务（在应用启动时注入）
	quotaService *tenantservice.QuotaService
)

// InitQuotaMiddleware 初始化配额中间件依赖
func InitQuotaMiddleware(qs *tenantservice.QuotaService) {
	quotaService = qs
}

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeBots        ResourceType = "bots"
	ResourceTypeWorkflows   ResourceType = "workflows"
	ResourceTypeMessages    ResourceType = "messages"
	ResourceTypeStorage     ResourceType = "storage"
	ResourceTypeTeamMembers ResourceType = "team_members"
)

// QuotaCheckConfig 配额检查配置
type QuotaCheckConfig struct {
	ResourceType  ResourceType // 资源类型
	RequiredCount int          // 需要的数量
	CheckOnly     bool         // 仅检查不扣减（默认false，会扣减配额）
}

// QuotaCheck 配额检查中间件
// 用于在处理请求前检查租户配额是否充足
//
// 使用示例：
//
//	r.POST("/api/bots", middleware.QuotaCheck(middleware.QuotaCheckConfig{
//	    ResourceType:  middleware.ResourceTypeBots,
//	    RequiredCount: 1,
//	    CheckOnly:     false, // 会扣减配额
//	}), handler.CreateBot)
func QuotaCheck(config QuotaCheckConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取tenant_id
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == string("") {
			// 尝试从上下文获取
			if tid, exists := c.Get("tenant_id"); exists {
				tenantID = string(tid)
			} else {
				c.JSON(consts.StatusBadRequest, responseWithError(berrno.ErrTenantNotFoundCode, "tenant_id is required"))
				c.Abort()
				return
			}
		}

		// 2. 检查是否已初始化配额服务
		if quotaService == nil {
			loggerx.CtxWarnf(ctx, "[QuotaCheck] quotaService not initialized, skipping quota check")
			c.Next(ctx)
			return
		}

		// 3. 检查配额
		err := quotaService.CheckQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
		if err != nil {
			// 检查是否是配额超限错误
			var quotaExceededErr *tenantservice.QuotaExceededError
			if errors.As(err, &quotaExceededErr) {
				loggerx.CtxWarnf(ctx, "[QuotaCheck] quota exceeded: tenant_id=%s, resource_type=%s, required=%d, current=%d, limit=%d",
					tenantID,
					config.ResourceType,
					config.RequiredCount,
					quotaExceededErr.CurrentUsage,
					quotaExceededErr.MaxLimit,
				)

				c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrQuotaExceededCode,
					errorx.KV(
						"resource_type", config.ResourceType,
						"current_usage", quotaExceededErr.CurrentUsage,
						"max_limit", quotaExceededErr.MaxLimit,
						"remaining", quotaExceededErr.MaxLimit-quotaExceededErr.CurrentUsage,
					),
				))
				c.Abort()
				return
			}

			// 其他错误
			loggerx.CtxErrorf(ctx, "[QuotaCheck] quota check failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaCheckFailedCode))
			c.Abort()
			return
		}

		// 4. 如果不是仅检查模式，则扣减配额
		if !config.CheckOnly {
			_, err = quotaService.ConsumeQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
			if err != nil {
				loggerx.CtxErrorf(ctx, "[QuotaCheck] consume quota failed: %v", err)

				// 配额检查通过但消费失败，返回错误
				c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaConsumeFailedCode))
				c.Abort()
				return
			}

			loggerx.CtxInfof(ctx, "[QuotaCheck] quota consumed: tenant_id=%s, resource_type=%s, count=%d",
				tenantID, config.ResourceType, config.RequiredCount)
		} else {
			loggerx.CtxInfof(ctx, "[QuotaCheck] quota check only (no consume): tenant_id=%s, resource_type=%s, count=%d",
				tenantID, config.ResourceType, config.RequiredCount)
		}

		// 5. 配额检查通过，继续处理请求
		c.Next(ctx)
	}
}

// QuotaCheckWithRollback 带回滚的配额检查中间件
// 用于在请求失败时自动回滚配额
//
// 使用场景：创建资源时，如果后续步骤失败，需要回滚已扣减的配额
//
// 使用示例：
//
//	r.POST("/api/bots",
//	    middleware.QuotaCheckWithRollback(middleware.QuotaCheckConfig{
//	        ResourceType:  middleware.ResourceTypeBots,
//	        RequiredCount: 1,
//	    }),
//	    handler.CreateBot)
func QuotaCheckWithRollback(config QuotaCheckConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
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

		if quotaService == nil {
			loggerx.CtxWarnf(ctx, "[QuotaCheck] quotaService not initialized, skipping quota check")
			c.Next(ctx)
			return
		}

		// 1. 检查配额
		err := quotaService.CheckQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
		if err != nil {
			var quotaExceededErr *tenantservice.QuotaExceededError
			if errors.As(err, &quotaExceededErr) {
				c.JSON(consts.StatusForbidden, responseWithError(berrno.ErrQuotaExceededCode,
					errorx.KV(
						"resource_type", config.ResourceType,
						"current_usage", quotaExceededErr.CurrentUsage,
						"max_limit", quotaExceededErr.MaxLimit,
					),
				))
				c.Abort()
				return
			}

			loggerx.CtxErrorf(ctx, "[QuotaCheck] quota check failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaCheckFailedCode))
			c.Abort()
			return
		}

		// 2. 扣减配额
		previousCount, err := quotaService.ConsumeQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
		if err != nil {
			loggerx.CtxErrorf(ctx, "[QuotaCheck] consume quota failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaConsumeFailedCode))
			c.Abort()
			return
		}

		loggerx.CtxInfof(ctx, "[QuotaCheckWithRollback] quota consumed: tenant_id=%s, resource_type=%s, count=%d, previous=%d",
			tenantID, config.ResourceType, config.RequiredCount, previousCount)

		// 3. 将配额信息保存到上下文，以便在失败时回滚
		c.Set("quota_rollback_info", map[string]interface{}{
			"tenant_id":     tenantID,
			"resource_type": config.ResourceType,
			"count":         config.RequiredCount,
		})

		// 4. 设置响应回滚拦截器
		c.Next(ctx)

		// 5. 检查响应状态码，如果失败则回滚配额
		statusCode := c.Response.StatusCode()
		if statusCode >= 400 {
			// 请求失败，回滚配额
			rollbackErr := quotaService.RollbackQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
			if rollbackErr != nil {
				loggerx.CtxErrorf(ctx, "[QuotaCheckWithRollback] rollback quota failed: %v", rollbackErr)
			} else {
				loggerx.CtxInfof(ctx, "[QuotaCheckWithRollback] quota rolled back: tenant_id=%s, resource_type=%s, count=%d",
					tenantID, config.ResourceType, config.RequiredCount)
			}
		}
	}
}

// ManualQuotaCheck 手动配额检查辅助函数
// 用于在Handler内部进行配额检查
//
// 使用示例：
//
//	func CreateBot(ctx context.Context, c *app.RequestContext) {
//	    err := middleware.ManualQuotaCheck(ctx, c, middleware.QuotaCheckConfig{
//	        ResourceType:  middleware.ResourceTypeBots,
//	        RequiredCount: 1,
//	    })
//	    if err != nil {
//	        // 配额检查失败
//	        return
//	    }
//	    // 继续处理...
//	}
func ManualQuotaCheck(ctx context.Context, c *app.RequestContext, config QuotaCheckConfig) error {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == string("") {
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = string(tid)
		} else {
			return errorx.New(berrno.ErrTenantNotFoundCode, errorx.KV("msg", "tenant_id is required"))
		}
	}

	if quotaService == nil {
		loggerx.CtxWarnf(ctx, "[ManualQuotaCheck] quotaService not initialized")
		return nil
	}

	// 检查配额
	err := quotaService.CheckQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
	if err != nil {
		return err
	}

	// 如果不是仅检查模式，则扣减配额
	if !config.CheckOnly {
		_, err = quotaService.ConsumeQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
		if err != nil {
			return err
		}
	}

	return nil
}

// ManualQuotaRollback 手动配额回滚辅助函数
// 用于在Handler内部手动回滚配额
//
// 使用示例：
//
//	func CreateBot(ctx context.Context, c *app.RequestContext) {
//	    // ... 处理逻辑 ...
//	    if err != nil {
//	        middleware.ManualQuotaRollback(ctx, c, middleware.QuotaCheckConfig{
//	            ResourceType:  middleware.ResourceTypeBots,
//	            RequiredCount: 1,
//	        })
//	        return
//	    }
//	}
func ManualQuotaRollback(ctx context.Context, c *app.RequestContext, config QuotaCheckConfig) error {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == string("") {
		if tid, exists := c.Get("tenant_id"); exists {
			tenantID = string(tid)
		} else {
			return errorx.New(berrno.ErrTenantNotFoundCode, errorx.KV("msg", "tenant_id is required"))
		}
	}

	if quotaService == nil {
		loggerx.CtxWarnf(ctx, "[ManualQuotaRollback] quotaService not initialized")
		return nil
	}

	return quotaService.RollbackQuota(ctx, tenantID, tenantservice.ResourceType(config.ResourceType), config.RequiredCount)
}

// responseWithError 构造错误响应
func responseWithError(code int, kv ...zap.Field) map[string]interface{} {
	kv = append(kv, zap.Int("code", code))
	errMsg := berrno.ErrMsgByCode(code)
	return map[string]interface{}{
		"code":    code,
		"message": errMsg,
		"data":    nil,
	}
}
