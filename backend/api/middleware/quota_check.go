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
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	tenantentity "github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	tenantservice "github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

var (
	// quotaService 配额服务（在应用启动时注入）
	quotaService *tenantservice.QuotaService
)

// InitQuotaMiddleware 初始化配额中间件依赖
func InitQuotaMiddleware(qs *tenantservice.QuotaService) {
	quotaService = qs
}

// ResourceType 资源类型（中间件层定义，对应 domain 层的 ResourceType）
type ResourceType string

const (
	ResourceTypeBots        ResourceType = "bots"         // Bot数量
	ResourceTypeMessages    ResourceType = "messages"     // 消息数量
	ResourceTypeStorage     ResourceType = "storage"      // 存储空间
	ResourceTypeTeamMembers ResourceType = "team_members" // 团队成员数量
	ResourceTypeWorkflows   ResourceType = "workflows"    // 工作流数量
)

// ToDomainResourceType 转换为 domain 层的 ResourceType
func (rt ResourceType) ToDomainResourceType() tenantentity.ResourceType {
	return tenantentity.ResourceType(rt)
}

// QuotaCheckConfig 配额检查配置
type QuotaCheckConfig struct {
	ResourceType  ResourceType // 资源类型
	RequiredCount int          // 需要的数量（默认1）
	CheckOnly     bool         // 仅检查不扣减（默认false，会扣减配额）
	SkipRollback  bool         // 跳过自动回滚（默认false，失败时自动回滚）
}

// QuotaCheck 配额检查中间件
// 用于在处理请求前检查租户配额是否充足
//
// 使用示例：
//
//	r.POST("/api/v1/bots",
//	    middleware.QuotaCheck(middleware.QuotaCheckConfig{
//	        ResourceType:  middleware.ResourceTypeBots,
//	        RequiredCount: 1,
//	        CheckOnly:     false, // 会扣减配额
//	    }),
//	    handler.CreateBot)
func QuotaCheck(config QuotaCheckConfig) app.HandlerFunc {
	// 设置默认值
	if config.RequiredCount <= 0 {
		config.RequiredCount = 1
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 获取 tenant_id
		tenantID, err := extractTenantID(c)
		if err != nil {
			c.JSON(consts.StatusBadRequest, responseWithError(berrno.ErrTenantNotFoundCode,
				errorx.KV("msg", "failed to extract tenant_id")))
			c.Abort()
			return
		}

		// 2. 检查是否已初始化配额服务
		if quotaService == nil {
			logs.CtxWarnf(ctx, "[QuotaCheck] quotaService not initialized, skipping quota check")
			c.Next(ctx)
			return
		}

		// 3. 检查配额
		domainResourceType := config.ResourceType.ToDomainResourceType()
		err = quotaService.CheckQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
		if err != nil {
			// 检查是否是配额超限错误
			var quotaExceededErr *tenantservice.QuotaExceededError
			if errors.As(err, &quotaExceededErr) {
				logs.CtxWarnf(ctx, "[QuotaCheck] quota exceeded: tenant_id=%s, resource_type=%s, required=%d, used=%d, limit=%d, usage=%.2f%%",
					tenantID,
					config.ResourceType,
					config.RequiredCount,
					quotaExceededErr.Used,
					quotaExceededErr.MaxLimit,
					quotaExceededErr.UsagePercent,
				)

				c.JSON(consts.StatusPaymentRequired, responseWithError(berrno.ErrQuotaExceededCode,
					errorx.KV(
						"msg", fmt.Sprintf("quota exceeded for %s", config.ResourceType),
						"resource_type", config.ResourceType,
						"used", quotaExceededErr.Used,
						"limit", quotaExceededErr.MaxLimit,
						"required", config.RequiredCount,
						"remaining", quotaExceededErr.MaxLimit - quotaExceededErr.Used,
						"usage_percent", fmt.Sprintf("%.2f%%", quotaExceededErr.UsagePercent),
					),
				))
				c.Abort()
				return
			}

			// 其他错误
			logs.CtxErrorf(ctx, "[QuotaCheck] quota check failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaCheckFailedCode))
			c.Abort()
			return
		}

		// 4. 如果不是仅检查模式，则扣减配额
		if !config.CheckOnly {
			err = quotaService.ConsumeQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
			if err != nil {
				logs.CtxErrorf(ctx, "[QuotaCheck] consume quota failed: %v", err)

				// 配额检查通过但消费失败，返回错误
				c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaConsumeFailedCode))
				c.Abort()
				return
			}

			logs.CtxInfof(ctx, "[QuotaCheck] quota consumed: tenant_id=%s, resource_type=%s, count=%d",
				tenantID, config.ResourceType, config.RequiredCount)
		} else {
			logs.CtxInfof(ctx, "[QuotaCheck] quota check only (no consume): tenant_id=%s, resource_type=%s, count=%d",
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
//	r.POST("/api/v1/bots",
//	    middleware.QuotaCheckWithRollback(middleware.QuotaCheckConfig{
//	        ResourceType:  middleware.ResourceTypeBots,
//	        RequiredCount: 1,
//	    }),
//	    handler.CreateBot)
func QuotaCheckWithRollback(config QuotaCheckConfig) app.HandlerFunc {
	// 设置默认值
	if config.RequiredCount <= 0 {
		config.RequiredCount = 1
	}

	return func(ctx context.Context, c *app.RequestContext) {
		tenantID, err := extractTenantID(c)
		if err != nil {
			c.JSON(consts.StatusBadRequest, responseWithError(berrno.ErrTenantNotFoundCode,
				errorx.KV("msg", "failed to extract tenant_id")))
			c.Abort()
			return
		}

		if quotaService == nil {
			logs.CtxWarnf(ctx, "[QuotaCheckWithRollback] quotaService not initialized, skipping quota check")
			c.Next(ctx)
			return
		}

		domainResourceType := config.ResourceType.ToDomainResourceType()

		// 1. 检查配额
		err = quotaService.CheckQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
		if err != nil {
			var quotaExceededErr *tenantservice.QuotaExceededError
			if errors.As(err, &quotaExceededErr) {
				c.JSON(consts.StatusPaymentRequired, responseWithError(berrno.ErrQuotaExceededCode,
					errorx.KV(
						"msg", fmt.Sprintf("quota exceeded for %s", config.ResourceType),
						"resource_type", config.ResourceType,
						"used", quotaExceededErr.Used,
						"limit", quotaExceededErr.MaxLimit,
						"remaining", quotaExceededErr.MaxLimit - quotaExceededErr.Used,
					),
				))
				c.Abort()
				return
			}

			logs.CtxErrorf(ctx, "[QuotaCheckWithRollback] quota check failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaCheckFailedCode))
			c.Abort()
			return
		}

		// 2. 扣减配额
		err = quotaService.ConsumeQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
		if err != nil {
			logs.CtxErrorf(ctx, "[QuotaCheckWithRollback] consume quota failed: %v", err)
			c.JSON(consts.StatusInternalServerError, responseWithError(berrno.ErrQuotaConsumeFailedCode))
			c.Abort()
			return
		}

		logs.CtxInfof(ctx, "[QuotaCheckWithRollback] quota consumed: tenant_id=%s, resource_type=%s, count=%d",
			tenantID, config.ResourceType, config.RequiredCount)

		// 3. 将配额信息保存到上下文，以便在失败时回滚
		c.Set("quota_rollback_info", map[string]interface{}{
			"tenant_id":     tenantID,
			"resource_type": config.ResourceType,
			"count":         config.RequiredCount,
			"skip_rollback": config.SkipRollback,
		})

		// 4. 处理请求
		c.Next(ctx)

		// 5. 检查响应状态码，如果失败则回滚配额
		if !config.SkipRollback {
			statusCode := c.Response.StatusCode()
			if statusCode >= 400 {
				// 请求失败，回滚配额
				rollbackErr := quotaService.RollbackQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
				if rollbackErr != nil {
					logs.CtxErrorf(ctx, "[QuotaCheckWithRollback] rollback quota failed: %v", rollbackErr)
				} else {
					logs.CtxInfof(ctx, "[QuotaCheckWithRollback] quota rolled back: tenant_id=%s, resource_type=%s, count=%d, status_code=%d",
						tenantID, config.ResourceType, config.RequiredCount, statusCode)
				}
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
	if config.RequiredCount <= 0 {
		config.RequiredCount = 1
	}

	tenantID, err := extractTenantID(c)
	if err != nil {
		return err
	}

	if quotaService == nil {
		logs.CtxWarnf(ctx, "[ManualQuotaCheck] quotaService not initialized")
		return nil
	}

	domainResourceType := config.ResourceType.ToDomainResourceType()

	// 检查配额
	err = quotaService.CheckQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
	if err != nil {
		return err
	}

	// 如果不是仅检查模式，则扣减配额
	if !config.CheckOnly {
		err = quotaService.ConsumeQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
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
//	    err := middleware.ManualQuotaCheck(ctx, c, middleware.QuotaCheckConfig{
//	        ResourceType:  middleware.ResourceTypeBots,
//	        RequiredCount: 1,
//	    })
//	    if err != nil {
//	        return
//	    }
//
//	    // ... 业务逻辑 ...
//	    if businessErr != nil {
//	        // 业务失败，回滚配额
//	        middleware.ManualQuotaRollback(ctx, c, middleware.QuotaCheckConfig{
//	            ResourceType:  middleware.ResourceTypeBots,
//	            RequiredCount: 1,
//	        })
//	        return businessErr
//	    }
//	}
func ManualQuotaRollback(ctx context.Context, c *app.RequestContext, config QuotaCheckConfig) error {
	if config.RequiredCount <= 0 {
		config.RequiredCount = 1
	}

	tenantID, err := extractTenantID(c)
	if err != nil {
		return err
	}

	if quotaService == nil {
		logs.CtxWarnf(ctx, "[ManualQuotaRollback] quotaService not initialized")
		return nil
	}

	domainResourceType := config.ResourceType.ToDomainResourceType()
	return quotaService.RollbackQuota(ctx, tenantID, domainResourceType, config.RequiredCount)
}

// GetQuotaInfo 获取配额信息辅助函数
// 用于查询租户的配额使用情况
//
// 使用示例：
//
//	func GetQuotaStatus(ctx context.Context, c *app.RequestContext) {
//	    quota, err := middleware.GetQuotaInfo(ctx, c, middleware.ResourceTypeBots)
//	    if err != nil {
//	        // 处理错误
//	        return
//	    }
//	    // 使用 quota 信息...
//	}
func GetQuotaInfo(ctx context.Context, c *app.RequestContext, resourceType ResourceType) (*tenantentity.Quota, error) {
	tenantID, err := extractTenantID(c)
	if err != nil {
		return nil, err
	}

	if quotaService == nil {
		return nil, fmt.Errorf("quotaService not initialized")
	}

	domainResourceType := resourceType.ToDomainResourceType()
	return quotaService.GetQuota(ctx, tenantID, domainResourceType)
}

// extractTenantID 从请求上下文中提取 tenant_id
func extractTenantID(c *app.RequestContext) (string, error) {
	// 1. 尝试从 Header 获取
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID != "" {
		return tenantID, nil
	}

	// 2. 尝试从 Query 参数获取
	tenantID = c.Query("tenant_id")
	if tenantID != "" {
		return tenantID, nil
	}

	// 3. 尝试从上下文获取（如果之前有中间件设置）
	if tid, exists := c.Get("tenant_id"); exists {
		if tenantIDStr, ok := tid.(string); ok {
			return tenantIDStr, nil
		}
	}

	return "", errorx.New(berrno.ErrTenantNotFoundCode, errorx.KV("msg", "tenant_id not found in request"))
}

// responseWithError 构造错误响应
func responseWithError(code int, kv ...errorx.KV) map[string]interface{} {
	result := map[string]interface{}{
		"code":    code,
		"message": berrno.ErrMsgByCode(code),
		"data":    nil,
	}

	// 添加额外的字段
	for _, item := range kv {
		result[item.Key] = item.Value
	}

	return result
}
