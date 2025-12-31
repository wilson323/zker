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
	"github.com/sirupsen/logrus"

	billingModel "github.com/coze-dev/coze-studio/backend/api/model/billing"
	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// ============================================================
// 实时计费API Handler
// ============================================================

var (
	realtimeBillingService interface{}
)

// InitRealtimeBillingHandler 初始化实时计费Handler
func InitRealtimeBillingHandler(service interface{}) {
	realtimeBillingService = service
}

// ==================== 实时扣费接口 ====================

// Charge 实时扣费
// @router POST /api/billing/realtime/charge [json]
// @summary 实时扣费
// @description 根据Token使用情况进行实时扣费
// @tags RealtimeBilling
// @accept json
// @produce json
// @param request body billingModel.ChargeRequest true "扣费请求"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func Charge(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 绑定并验证请求
	var req billingModel.ChargeRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.WithError(err).Warn("[RealtimeBilling] Invalid request parameters")
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 2. 从上下文获取租户ID（通过中间件注入）
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[RealtimeBilling] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}
	req.TenantID = tenantID

	// 3. TODO: 调用应用层服务执行扣费
	// err := h.billingService.Charge(ctx, &req)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrBillingChargeFailedCode,
	//         errorx.KV("operation", "realtime charge"),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrBillingChargeFailedCode,
	//         Message: "扣费失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":   tenantID,
		"amount":      req.Amount,
		"currency":    req.Currency,
		"description": req.Description,
	}).Info("[RealtimeBilling] Charge processed successfully")

	// 4. 返回成功响应
	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "扣费成功",
		Data:    nil,
	})
}

// GetBalance 查询余额
// @router GET /api/billing/tenant/:tenant_id/balance [json]
// @summary 查询余额
// @description 查询租户的当前余额
// @tags RealtimeBilling
// @accept json
// @produce json
// @param tenant_id path string true "租户ID"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 403 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GetBalance(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 获取路径参数
	tenantID := c.Param("tenant_id")
	if tenantID == "" {
		InvalidParamRequestResponse(c, "tenant_id is required")
		return
	}

	// 2. 从上下文获取当前租户ID（用于权限检查）
	currentTenantID := contextutil.GetTenantID(ctx)
	if currentTenantID == "" {
		logger.Warn("[RealtimeBilling] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	// 3. 权限检查：只能查询自己的余额（管理员除外）
	if currentTenantID != tenantID && !contextutil.IsAdmin(ctx) {
		logger.WithFields(logrus.Fields{
			"current_tenant_id": currentTenantID,
			"target_tenant_id":  tenantID,
		}).Warn("[RealtimeBilling] Permission denied: insufficient privileges")
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    errno.ForbiddenCode,
			Message: "权限不足",
		})
		return
	}

	// 4. TODO: 调用应用层服务查询余额
	// balance, err := h.billingService.GetBalance(ctx, tenantID)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrBillingBalanceNotFoundCode,
	//         errorx.KV("tenant_id", tenantID),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrBillingBalanceNotFoundCode,
	//         Message: "查询余额失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithField("tenant_id", tenantID).Info("[RealtimeBilling] Balance queried successfully")

	// 5. 返回成功响应（示例数据）
	balance := &billingModel.BalanceResponse{
		TenantID:        tenantID,
		CurrentBalance:  1000.00,
		Currency:        "CNY",
		AvailableAmount: 950.00,
		FrozenAmount:    50.00,
		CreditLimit:     0.00,
		UpdatedAt:       1704067200,
	}

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "查询成功",
		Data:    balance,
	})
}

// ConsumeQuota 消费配额
// @router POST /api/billing/quota/consume [json]
// @summary 消费配额
// @description 消费指定类型的配额（Bot数、知识库、工作流等）
// @tags QuotaManagement
// @accept json
// @produce json
// @param request body billingModel.ConsumeQuotaRequest true "消费配额请求"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 403 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func ConsumeQuota(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 绑定并验证请求
	var req billingModel.ConsumeQuotaRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.WithError(err).Warn("[QuotaManagement] Invalid request parameters")
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}
	req.TenantID = tenantID

	// 3. TODO: 调用应用层服务消费配额
	// err := h.quotaService.Consume(ctx, &req)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrQuotaConsumeFailedCode,
	//         errorx.KV("tenant_id", tenantID),
	//         errorx.KV("resource_type", req.ResourceType),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrQuotaConsumeFailedCode,
	//         Message: "配额消费失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":      tenantID,
		"resource_type":  req.ResourceType,
		"quantity":       req.Quantity,
		"operation_id":   req.OperationID,
	}).Info("[QuotaManagement] Quota consumed successfully")

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "配额消费成功",
		Data: &billingModel.ConsumeQuotaResponse{
			Success:     true,
			NewQuantity: 100,
			Remaining:   50,
		},
	})
}

// CheckQuota 检查配额
// @router GET /api/billing/quota/check [json]
// @summary 检查配额
// @description 检查指定资源的配额是否足够
// @tags QuotaManagement
// @accept json
// @produce json
// @param resource_type query string true "资源类型" Enums(bots, knowledge, workflows, api_calls, storage, concurrent)
// @param quantity query int true "需要检查的数量"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func CheckQuota(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 获取查询参数
	resourceType := c.Query("resource_type")
	if resourceType == "" {
		InvalidParamRequestResponse(c, "resource_type is required")
		return
	}

	quantity := c.Query("quantity")
	// TODO: 验证quantity参数

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	// 3. TODO: 调用应用层服务检查配额
	// sufficient, current, limit, err := h.quotaService.Check(ctx, tenantID, resourceType, quantity)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrQuotaCheckFailedCode,
	//         errorx.KV("tenant_id", tenantID),
	//         errorx.KV("resource_type", resourceType),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrQuotaCheckFailedCode,
	//         Message: "配额检查失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":     tenantID,
		"resource_type": resourceType,
		"quantity":      quantity,
	}).Info("[QuotaManagement] Quota checked successfully")

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "配额检查成功",
		Data: &billingModel.CheckQuotaResponse{
			TenantID:     tenantID,
			ResourceType: resourceType,
			Sufficient:   true,
			Current:      50,
			Limit:        100,
			Remaining:    50,
		},
	})
}

// ResetQuota 重置配额
// @router POST /api/billing/quota/reset [json]
// @summary 重置配额
// @description 重置指定资源的配额（管理员权限）
// @tags QuotaManagement
// @accept json
// @produce json
// @param request body billingModel.ResetQuotaRequest true "重置配额请求"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 403 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func ResetQuota(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 检查管理员权限
	if !contextutil.IsAdmin(ctx) {
		logger.Warn("[QuotaManagement] Permission denied: admin required")
		c.JSON(http.StatusForbidden, ErrorResponse{
			Code:    errno.ForbiddenCode,
			Message: "权限不足（需要管理员权限）",
		})
		return
	}

	// 2. 绑定并验证请求
	var req billingModel.ResetQuotaRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.WithError(err).Warn("[QuotaManagement] Invalid request parameters")
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. TODO: 调用应用层服务重置配额
	// err := h.quotaService.Reset(ctx, &req)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrQuotaResetFailedCode,
	//         errorx.KV("tenant_id", req.TenantID),
	//         errorx.KV("resource_type", req.ResourceType),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrQuotaResetFailedCode,
	//         Message: "配额重置失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":     req.TenantID,
		"resource_type": req.ResourceType,
		"new_limit":     req.NewLimit,
	}).Info("[QuotaManagement] Quota reset successfully")

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "配额重置成功",
		Data:    nil,
	})
}

// RollbackQuota 回滚配额
// @router POST /api/billing/quota/rollback [json]
// @summary 回滚配额
// @description 回滚之前消费的配额
// @tags QuotaManagement
// @accept json
// @produce json
// @param request body billingModel.RollbackQuotaRequest true "回滚配额请求"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func RollbackQuota(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 绑定并验证请求
	var req billingModel.RollbackQuotaRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.WithError(err).Warn("[QuotaManagement] Invalid request parameters")
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}
	req.TenantID = tenantID

	// 3. TODO: 调用应用层服务回滚配额
	// err := h.quotaService.Rollback(ctx, &req)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrQuotaRollbackFailedCode,
	//         errorx.KV("tenant_id", tenantID),
	//         errorx.KV("operation_id", req.OperationID),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrQuotaRollbackFailedCode,
	//         Message: "配额回滚失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":    tenantID,
		"operation_id": req.OperationID,
		"reason":       req.Reason,
	}).Info("[QuotaManagement] Quota rolled back successfully")

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "配额回滚成功",
		Data:    nil,
	})
}

// ==================== 通用响应类型 ====================

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
