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
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"

	"github.com/coze-dev/coze-studio/backend/api/model/billing"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var budgetManagementService *service.BudgetManagementService

// InitBudgetManagementService 初始化预算管理服务
func InitBudgetManagementService(svc *service.BudgetManagementService, logger *logrus.Logger) {
	budgetManagementService = svc
	logger.Info("[BudgetManagement] Budget management service initialized")
}

// ============================================================
// Budget Management API Handlers
// ============================================================

// GetBudget 获取预算配置
// @router /api/v1/tenants/:tenant_id/budget [GET]
func GetBudget(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, &billing.BudgetResponse{
			Code:    400,
			Message: "tenant_id is required",
		})
		return
	}

	// 2. 调用服务
	budgetDTO, err := budgetManagementService.GetBudget(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to get budget: %v", err)
		c.JSON(http.StatusInternalServerError, &billing.BudgetResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// 3. 转换DTO为API响应
	budgetData := &billing.BudgetSettingsData{
		TenantID:              budgetDTO.TenantID,
		BudgetType:            budgetDTO.BudgetType,
		BudgetAmount:          budgetDTO.BudgetAmount,
		Currency:              budgetDTO.Currency,
		AlertThreshold1:       budgetDTO.AlertThreshold1,
		AlertThreshold2:       budgetDTO.AlertThreshold2,
		HardCapEnabled:        budgetDTO.HardCapEnabled,
		HardCapAmount:         budgetDTO.HardCapAmount,
		AutoDowngradeEnabled:  budgetDTO.AutoDowngradeEnabled,
		DowngradeThreshold:    budgetDTO.DowngradeThreshold,
		OriginalModel:         budgetDTO.OriginalModel,
		FallbackModel:         budgetDTO.FallbackModel,
		NotificationChannels:   budgetDTO.NotificationChannels,
		NotificationRecipients: budgetDTO.NotificationRecipients,
		CreatedAt:             budgetDTO.CreatedAt,
		UpdatedAt:             budgetDTO.UpdatedAt,
	}

	// 4. 返回结果
	httputil.BuildSuccessResp(c, &billing.BudgetResponse{
		Code:    0,
		Message: "success",
		Data:    budgetData,
	})
}

// CreateBudget 创建预算配置
// @router /api/v1/tenants/:tenant_id/budget [POST]
func CreateBudget(ctx context.Context, c *app.RequestContext) {
	// 1. 解析请求
	var req service.CreateBudgetRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 2. 获取租户ID（从路径参数或请求体）
	tenantID := c.Param("tenant_id")
	if tenantID != "" {
		req.TenantID = tenantID
	}

	if req.TenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 3. 调用服务
	budget, err := budgetManagementService.CreateBudget(ctx, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to create budget: %v", err)

		// 判断错误类型
		if err.Error() == "budget already exists for tenant "+req.TenantID {
			c.JSON(http.StatusConflict, &billing.BudgetResponse{
				Code:    409,
				Message: err.Error(),
			})
			return
		}

		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 4. 返回结果
	c.JSON(http.StatusCreated, &billing.BudgetResponse{
		Code:    0,
		Message: "Budget created successfully",
		Data:    budget,
	})
}

// UpdateBudget 更新预算配置
// @router /api/v1/tenants/:tenant_id/budget [PUT]
func UpdateBudget(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 解析请求
	var req service.UpdateBudgetRequest
	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. 调用服务
	budget, err := budgetManagementService.UpdateBudget(ctx, tenantID, &req)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to update budget: %v", err)

		// 判断错误类型
		if err.Error() == "budget not found for tenant "+tenantID {
			c.JSON(http.StatusNotFound, &billing.BudgetResponse{
				Code:    404,
				Message: err.Error(),
			})
			return
		}

		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 4. 返回结果
	httputil.BuildSuccessResp(c, &billing.BudgetResponse{
		Code:    0,
		Message: "Budget updated successfully",
		Data:    budget,
	})
}

// DeleteBudget 删除预算配置
// @router /api/v1/tenants/:tenant_id/budget [DELETE]
func DeleteBudget(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 调用服务
	err := budgetManagementService.DeleteBudget(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to delete budget: %v", err)

		// 判断错误类型
		if err.Error() == "budget not found for tenant "+tenantID {
			c.JSON(http.StatusNotFound, &billing.BudgetResponse{
				Code:    404,
				Message: err.Error(),
			})
			return
		}

		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 3. 返回结果
	httputil.BuildSuccessResp(c, &billing.BudgetResponse{
		Code:    0,
		Message: "Budget deleted successfully",
	})
}

// GetBudgetUsage 获取预算使用情况
// @router /api/v1/tenants/:tenant_id/budget/usage [GET]
func GetBudgetUsage(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 调用服务
	usage, err := budgetManagementService.GetBudgetUsage(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to get budget usage: %v", err)

		// 判断错误类型
		if err.Error() == "budget not found for tenant "+tenantID {
			c.JSON(http.StatusNotFound, &billing.BudgetUsageResponse{
				Code:    404,
				Message: err.Error(),
			})
			return
		}

		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 3. 返回结果
	httputil.BuildSuccessResp(c, &billing.BudgetUsageResponse{
		Code:    0,
		Message: "success",
		Data:    usage,
	})
}

// CheckBudget 手动触发预算检查
// @router /api/v1/tenants/:tenant_id/budget/check [POST]
func CheckBudget(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 调用服务
	result, err := budgetManagementService.CheckBudget(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to check budget: %v", err)
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 3. 返回结果
	httputil.BuildSuccessResp(c, &billing.BudgetCheckResponse{
		Code:    0,
		Message: "success",
		Data:    result,
	})
}

// GetBudgetAlerts 获取告警历史
// @router /api/v1/tenants/:tenant_id/budget/alerts [GET]
func GetBudgetAlerts(ctx context.Context, c *app.RequestContext) {
	// 1. 获取租户ID
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 解析查询参数
	filter := &service.AlertFilter{}

	// 告警类型
	if alertType := c.Query("alert_type"); alertType != "" {
		filter.AlertType = &alertType
	}

	// 告警级别
	if alertLevel := c.Query("alert_level"); alertLevel != "" {
		filter.AlertLevel = &alertLevel
	}

	// 通知状态
	if notifStatus := c.Query("notification_status"); notifStatus != "" {
		filter.NotificationStatus = &notifStatus
	}

	// 开始日期
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			filter.StartDate = &startDate
		}
	}

	// 结束日期
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			filter.EndDate = &endDate
		}
	}

	// 分页参数
	filter.PageToken = c.Query("page_token")

	if pageSize := c.Query("page_size"); pageSize != "" {
		if size, err := strconv.Atoi(pageSize); err == nil {
			filter.PageSize = size
		}
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// 3. 调用服务
	alerts, pagination, err := budgetManagementService.GetAlerts(ctx, tenantID, filter)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetAPI] Failed to get alerts: %v", err)
		internalServerErrorResponse(ctx, c, err)
		return
	}

	// 4. 返回结果
	httputil.BuildSuccessResp(c, &billing.BudgetAlertsResponse{
		Code:       0,
		Message:    "success",
		Data:       alerts,
		Pagination: pagination,
	})
}
