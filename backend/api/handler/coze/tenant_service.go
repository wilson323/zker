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
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/model/tenant"
	tenantapp "github.com/coze-dev/coze-studio/backend/application/tenant"
)

// ==================== 租户管理接口 ====================

// CreateTenant 创建租户
// @router /api/tenants [POST]
func CreateTenant(ctx context.Context, c *app.RequestContext) {
	var err error
	var req tenant.CreateTenantRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.CreateTenant(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTenant 获取租户详情
// @router /api/tenants/:tenant_id [GET]
func GetTenant(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := tenantapp.TenantAppSVC.GetTenant(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateTenant 更新租户信息
// @router /api/tenants/:tenant_id [PUT]
func UpdateTenant(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.UpdateTenantRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.UpdateTenant(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteTenant 删除租户
// @router /api/tenants/:tenant_id [DELETE]
func DeleteTenant(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	err := tenantapp.TenantAppSVC.DeleteTenant(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// ListTenants 列出租户
// @router /api/tenants [GET]
func ListTenants(ctx context.Context, c *app.RequestContext) {
	var req tenant.ListTenantsRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.ListTenants(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 订阅管理接口 ====================

// GetSubscription 获取订阅信息
// @router /api/tenants/:tenant_id/subscription [GET]
func GetSubscription(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := tenantapp.TenantAppSVC.GetSubscription(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateSubscription 更新订阅
// @router /api/tenants/:tenant_id/subscription [PUT]
func UpdateSubscription(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.UpdateSubscriptionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.UpdateSubscription(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpgradeSubscription 升级订阅等级
// @router /api/tenants/:tenant_id/upgrade [POST]
func UpgradeSubscription(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.UpgradeSubscriptionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.UpgradeSubscription(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 配额管理接口 ====================

// GetQuotas 获取配额状态
// @router /api/tenants/:tenant_id/quotas [GET]
func GetQuotas(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := tenantapp.TenantAppSVC.GetQuotas(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CheckQuota 检查配额
// @router /api/tenants/:tenant_id/quotas/check [POST]
func CheckQuota(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.CheckQuotaRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.CheckQuota(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ConsumeQuota 消费配额
// @router /api/tenants/:tenant_id/quotas/consume [POST]
func ConsumeQuota(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.ConsumeQuotaRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.ConsumeQuota(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RollbackQuota 回滚配额
// @router /api/tenants/:tenant_id/quotas/rollback [POST]
func RollbackQuota(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.RollbackQuotaRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = tenantapp.TenantAppSVC.RollbackQuota(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// ==================== 计费管理接口 ====================

// RecordUsage 记录使用量
// @router /api/tenants/:tenant_id/usage [POST]
func RecordUsage(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.RecordUsageRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	err = tenantapp.TenantAppSVC.RecordUsage(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &tenant.RecordUsageResponse{
		Code:    0,
		Message: "success",
	})
}

// GetUsageSummary 获取使用量汇总
// @router /api/tenants/:tenant_id/usage/summary [GET]
func GetUsageSummary(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var req tenant.GetUsageSummaryRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.GetUsageSummary(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GenerateInvoice 生成账单
// @router /api/tenants/:tenant_id/invoices/generate [POST]
func GenerateInvoice(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	var err error
	var req tenant.GenerateInvoiceRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := tenantapp.TenantAppSVC.GenerateInvoice(ctx, tenantID, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetInvoice 获取账单详情
// @router /api/tenants/:tenant_id/invoices/:invoice_id [GET]
func GetInvoice(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	invoiceID := c.Param("invoice_id")
	if invoiceID == "" {
		invalidParamRequestResponse(c, "invoice_id is required")
		return
	}

	resp, err := tenantapp.TenantAppSVC.GetInvoice(ctx, tenantID, invoiceID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListInvoices 列出账单
// @router /api/tenants/:tenant_id/invoices [GET]
func ListInvoices(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	status := c.Query("status")
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	var limitInt, offsetInt int
	_, err := fmt.Sscanf(limit, "%d", &limitInt)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}
	_, err = fmt.Sscanf(offset, "%d", &offsetInt)
	if err != nil || offsetInt < 0 {
		offsetInt = 0
	}

	resp, err := tenantapp.TenantAppSVC.ListInvoices(ctx, tenantID, &status, limitInt, offsetInt)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// PayInvoice 支付账单
// @router /api/tenants/:tenant_id/invoices/:invoice_id/pay [POST]
func PayInvoice(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	invoiceID := c.Param("invoice_id")
	if invoiceID == "" {
		invalidParamRequestResponse(c, "invoice_id is required")
		return
	}

	resp, err := tenantapp.TenantAppSVC.PayInvoice(ctx, tenantID, invoiceID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 监控接口 ====================

// GetMonitoringStatus 获取配额监控状态
// @router /api/tenants/:tenant_id/monitoring [GET]
func GetMonitoringStatus(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	resp, err := tenantapp.TenantAppSVC.GetMonitoringStatus(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ==================== 数据隔离接口 ====================

// ExportTenantData 导出租户数据
// @router /api/tenants/:tenant_id/export [POST]
func ExportTenantData(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 获取导出选项
	includeBots := c.DefaultQuery("include_bots", "true") == "true"
	includeWorkflows := c.DefaultQuery("include_workflows", "true") == "true"
	includeKnowledge := c.DefaultQuery("include_knowledge", "true") == "true"
	includeConversations := c.DefaultQuery("include_conversations", "false") == "true"

	exportID, err := tenantapp.TenantAppSVC.ExportTenantData(ctx, tenantID, &tenant.ExportDataRequest{
		IncludeBots:         includeBots,
		IncludeWorkflows:    includeWorkflows,
		IncludeKnowledge:    includeKnowledge,
		IncludeConversations: includeConversations,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &tenant.ExportDataResponse{
		Code:     0,
		Message:  "Export started",
		ExportID: exportID,
	})
}

// GetExportStatus 获取导出状态
// @router /api/tenants/:tenant_id/export/:export_id [GET]
func GetExportStatus(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	exportID := c.Param("export_id")
	if exportID == "" {
		invalidParamRequestResponse(c, "export_id is required")
		return
	}

	status, downloadURL, err := tenantapp.TenantAppSVC.GetExportStatus(ctx, tenantID, exportID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &tenant.ExportStatusResponse{
		Code:        0,
		Message:     "success",
		Status:      status,
		DownloadURL: downloadURL,
	})
}

// DeleteTenantData 删除租户数据（软删除）
// @router /api/tenants/:tenant_id/data [DELETE]
func DeleteTenantData(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 获取删除确认
	confirm := c.Query("confirm")
	if confirm != "true" {
		invalidParamRequestResponse(c, "confirmation required: add ?confirm=true to confirm deletion")
		return
	}

	err := tenantapp.TenantAppSVC.DeleteTenantData(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "Tenant data deletion initiated",
	})
}

// GetTenantStats 获取租户统计信息
// @router /api/tenants/:tenant_id/stats [GET]
func GetTenantStats(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		invalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	stats, err := tenantapp.TenantAppSVC.GetTenantStats(ctx, tenantID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, &tenant.TenantStatsResponse{
		Code:    0,
		Message: "success",
		Data:    stats,
	})
}
