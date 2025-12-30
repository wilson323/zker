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
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var (
	upgradeSvc *service.IsolationUpgradeService
)

// InitIsolationUpgradeHandler 初始化隔离升级Handler
func InitIsolationUpgradeHandler(db *gorm.DB, tenantRepo repository.TenantRepository) {
	upgradeSvc = service.NewIsolationUpgradeService(db, tenantRepo, service.DefaultUpgradeConfig())
}

// GetIsolationStatus 获取租户隔离状态
// @router /api/v1/tenants/:tenant_id/isolation [GET]
func GetIsolationStatus(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
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

	status, err := upgradeSvc.GetUpgradeStatus(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "Failed to get isolation status: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to get isolation status",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    status,
	})
}

// UpgradeIsolationStrategy 手动升级隔离策略
// @router /api/v1/tenants/:tenant_id/isolation/upgrade [POST]
func UpgradeIsolationStrategy(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "tenant_id is required",
		})
		return
	}

	var req struct {
		TargetStrategy string `json:"target_strategy" validate:"required,oneof=row_level schema_level database_level"`
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

	targetStrategy := service.IsolationStrategy(req.TargetStrategy)

	// 执行升级
	if err := upgradeSvc.ManualUpgrade(ctx, tenantID, targetStrategy); err != nil {
		logs.CtxErrorf(ctx, "Failed to upgrade isolation: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Failed to upgrade isolation strategy",
			"error": map[string]interface{}{
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Isolation strategy upgraded successfully",
		"data": map[string]interface{}{
			"tenant_id":       tenantID,
			"target_strategy": targetStrategy,
		},
	})
}
