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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var (
	managementSvc *service.TenantManagementService
)

// InitTenantManagementHandler 初始化租户管理Handler
func InitTenantManagementHandler(db *gorm.DB, tenantRepo repository.TenantRepository, quotaRepo repository.QuotaRepository, subscriptionRepo repository.SubscriptionRepository) {
	managementSvc = service.NewTenantManagementService(db, tenantRepo, quotaRepo, subscriptionRepo)
}

// GetTenantInfo 获取租户详细信息
// @router /api/v1/tenants/:tenant_id [GET]
func GetTenantInfo(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")

	// 如果没有传tenant_id，尝试从上下文获取
	if tenantID == "" {
		if tenant, err := middleware.GetTenantFromContext(c); err == nil {
			tenantID = tenant.TenantID
		} else {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "tenant_id is required",
			})
			return
		}
	}

	detail, err := managementSvc.GetTenantInfo(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to get tenant info: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to get tenant info",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    detail,
	})
}

// UpdateTenantInfo 更新租户信息
// @router /api/v1/tenants/:tenant_id [PUT]
func UpdateTenantInfo(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tenant_id is required",
		})
		return
	}

	var req service.UpdateTenantRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid request parameters",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	tenant, err := managementSvc.UpdateTenantInfo(ctx, tenantID, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to update tenant info: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to update tenant info",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Tenant info updated successfully",
		"data":    tenant,
	})
}

// GetQuotaUsage 获取配额使用情况
// @router /api/v1/tenants/:tenant_id/quotas [GET]
func GetQuotaUsage(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		// 尝试从上下文获取
		if tenant, err := middleware.GetTenantFromContext(c); err == nil {
			tenantID = tenant.TenantID
		} else {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "tenant_id is required",
			})
			return
		}
	}

	usages, err := managementSvc.GetQuotaUsage(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to get quota usage: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to get quota usage",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"quotas": usages,
		},
	})
}

// CheckQuotaAvailable 检查配额是否可用
// @router /api/v1/tenants/:tenant_id/quotas/check [POST]
func CheckQuotaAvailable(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tenant_id is required",
		})
		return
	}

	var req struct {
		ResourceType string `json:"resource_type" validate:"required"`
		RequiredCount int    `json:"required_count" validate:"required,min=1"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid request parameters",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	available, err := managementSvc.CheckQuotaAvailable(ctx, tenantID, req.ResourceType, req.RequiredCount)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to check quota: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to check quota",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"available": available,
		},
	})
}

// GetSubscription 获取订阅信息
// @router /api/v1/tenants/:tenant_id/subscription [GET]
func GetSubscription(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		// 尝试从上下文获取
		if tenant, err := middleware.GetTenantFromContext(c); err == nil {
			tenantID = tenant.TenantID
		} else {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": "tenant_id is required",
			})
			return
		}
	}

	detail, err := managementSvc.GetSubscription(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to get subscription: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to get subscription",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    detail,
	})
}

// UpgradeSubscription 升级订阅套餐
// @router /api/v1/tenants/:tenant_id/subscription/upgrade [POST]
func UpgradeSubscription(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tenant_id is required",
		})
		return
	}

	var req service.UpgradeSubscriptionRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Invalid request parameters",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	subscription, err := managementSvc.UpgradeSubscription(ctx, tenantID, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to upgrade subscription: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to upgrade subscription",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Subscription upgraded successfully",
		"data":    subscription,
	})
}
