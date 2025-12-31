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

package billing

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/sirupsen/logrus"

	"github.com/coze-dev/coze-studio/backend/core/bizerror"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

// ============================================================
// 提现管理Handler
// ============================================================

// WithdrawalHandler 提现处理器
type WithdrawalHandler struct {
	withdrawalService *service.WithdrawalService
	logger            *logrus.Logger
}

// NewWithdrawalHandler 创建提现处理器实例
func NewWithdrawalHandler(
	withdrawalService *service.WithdrawalService,
	logger *logrus.Logger,
) *WithdrawalHandler {
	return &WithdrawalHandler{
		withdrawalService: withdrawalService,
		logger:            logger,
	}
}

// ============================================================
// API实现
// ============================================================

// CreateWithdrawalRequest 创建提现申请
// @router POST /api/billing/withdrawals
// @summary 创建提现申请
// @description 用户创建提现申请，系统自动冻结余额并计算手续费
// @tags billing/withdrawal
// @accept json
// @produce json
// @param request body service.CreateWithdrawalRequest true "提现申请信息"
// @success 200 {object} WithdrawalRequestResponse
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) CreateWithdrawalRequest(ctx context.Context, c *app.RequestContext) {
	var req service.CreateWithdrawalRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Invalid request parameters")
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid request parameters: "+err.Error()))
		return
	}

	// 从上下文获取租户ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		c.JSON(consts.StatusUnauthorized, bizerror.NewBizError(consts.StatusUnauthorized, "Unauthorized: missing tenant_id"))
		return
	}
	req.TenantID = tenantID

	// 调用服务
	resp, err := h.withdrawalService.CreateWithdrawal(ctx, &req)
	if err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Failed to create withdrawal request")
		c.JSON(consts.StatusInternalServerError, bizerror.NewBizError(consts.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// ReviewWithdrawal 审核提现申请
// @router POST /api/billing/withdrawals/:id/review
// @summary 审核提现申请
// @description 管理员审核提现申请，批准或拒绝
// @tags billing/withdrawal
// @accept json
// @produce json
// @param id path int true "提现ID"
// @param request body service.ReviewWithdrawalRequest true "审核信息"
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) ReviewWithdrawal(ctx context.Context, c *app.RequestContext) {
	// 获取提现ID
	idStr := c.Param("id")
	withdrawalID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid withdrawal ID"))
		return
	}

	var req service.ReviewWithdrawalRequest
	if err := c.BindAndValidate(&req); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Invalid review parameters")
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid review parameters: "+err.Error()))
		return
	}

	req.WithdrawalID = withdrawalID

	// 从上下文获取审核人ID
	reviewerID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		c.JSON(consts.StatusUnauthorized, bizerror.NewBizError(consts.StatusUnauthorized, "Unauthorized: missing user_id"))
		return
	}
	req.ReviewerID = reviewerID

	// 调用服务
	if err := h.withdrawalService.ReviewWithdrawal(ctx, &req); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Failed to review withdrawal")
		c.JSON(consts.StatusInternalServerError, bizerror.NewBizError(consts.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Withdrawal reviewed successfully",
	})
}

// ExecuteWithdrawal 执行提现
// @router POST /api/billing/withdrawals/:id/execute
// @summary 执行提现
// @description 系统自动执行已批准的提现申请
// @tags billing/withdrawal
// @accept json
// @produce json
// @param id path int true "提现ID"
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) ExecuteWithdrawal(ctx context.Context, c *app.RequestContext) {
	// 获取提现ID
	idStr := c.Param("id")
	withdrawalID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid withdrawal ID"))
		return
	}

	// 调用服务
	if err := h.withdrawalService.ExecuteWithdrawal(ctx, withdrawalID); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Failed to execute withdrawal")
		c.JSON(consts.StatusInternalServerError, bizerror.NewBizError(consts.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Withdrawal executed successfully",
	})
}

// GetWithdrawalByID 获取提现详情
// @router GET /api/billing/withdrawals/:id
// @summary 获取提现详情
// @description 根据ID获取提现申请详细信息
// @tags billing/withdrawal
// @accept json
// @produce json
// @param id path int true "提现ID"
// @success 200 {object} entity.WithdrawalRequest
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) GetWithdrawalByID(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现获取提现详情
	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Get withdrawal by ID - to be implemented",
	})
}

// ListWithdrawals 查询提现列表
// @router GET /api/billing/withdrawals
// @summary 查询提现列表
// @description 分页查询提现申请列表
// @tags billing/withdrawal
// @accept json
// @produce json
// @param status query string false "状态过滤"
// @param withdrawal_method query string false "提现方式"
// @param start_date query string false "开始日期"
// @param end_date query string false "结束日期"
// @param page query int false "页码" default(1)
// @param page_size query int false "每页数量" default(20)
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) ListWithdrawals(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现查询提现列表
	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "List withdrawals - to be implemented",
	})
}

// GetWithdrawalStatistics 获取提现统计
// @router GET /api/billing/withdrawals/statistics
// @summary 获取提现统计
// @description 获取指定时间范围内的提现统计数据
// @tags billing/withdrawal
// @accept json
// @produce json
// @param start_date query string false "开始日期(YYYY-MM-DD)"
// @param end_date query string false "结束日期(YYYY-MM-DD)"
// @success 200 {object} service.WithdrawalStatisticsResponse
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) GetWithdrawalStatistics(ctx context.Context, c *app.RequestContext) {
	// 从上下文获取租户ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		c.JSON(consts.StatusUnauthorized, bizerror.NewBizError(consts.StatusUnauthorized, "Unauthorized: missing tenant_id"))
		return
	}

	// 获取查询参数
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	// 调用服务
	stats, err := h.withdrawalService.GetWithdrawalStatistics(ctx, tenantID, startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Failed to get withdrawal statistics")
		c.JSON(consts.StatusInternalServerError, bizerror.NewBizError(consts.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, stats)
}

// GetWithdrawalLimits 获取提现限额
// @router GET /api/billing/withdrawals/limits
// @summary 获取提现限额
// @description 获取各种提现方式的限额配置
// @tags billing/withdrawal
// @accept json
// @produce json
// @param withdrawal_method query string false "提现方式"
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) GetWithdrawalLimits(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现获取提现限额
	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Get withdrawal limits - to be implemented",
	})
}

// GetWithdrawalRecords 获取提现流水
// @router GET /api/billing/withdrawals/records
// @summary 获取提现流水
// @description 分页查询提现流水记录
// @tags billing/withdrawal
// @accept json
// @produce json
// @param transaction_type query string false "交易类型"
// @param start_date query string false "开始日期"
// @param end_date query string false "结束日期"
// @param page query int false "页码" default(1)
// @param page_size query int false "每页数量" default(20)
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) GetWithdrawalRecords(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现获取提现流水
	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Get withdrawal records - to be implemented",
	})
}

// Reconciliation 执行对账
// @router POST /api/billing/withdrawals/reconciliation
// @summary 执行对账
// @description 执行指定日期的提现对账
// @tags billing/withdrawal
// @accept json
// @produce json
// @param request body map[string]string true "对账日期"
// @success 200 {object} entity.WithdrawalReconciliation
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) Reconciliation(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ReconciliationDate string `json:"reconciliation_date" validate:"required"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Invalid reconciliation parameters")
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid parameters: "+err.Error()))
		return
	}

	// 调用服务
	reconciliation, err := h.withdrawalService.Reconciliation(ctx, req.ReconciliationDate)
	if err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Failed to execute reconciliation")
		c.JSON(consts.StatusInternalServerError, bizerror.NewBizError(consts.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, reconciliation)
}

// GetReconciliationList 获取对账列表
// @router GET /api/billing/withdrawals/reconciliations
// @summary 获取对账列表
// @description 分页查询对账记录列表
// @tags billing/withdrawal
// @accept json
// @produce json
// @param start_date query string false "开始日期"
// @param end_date query string false "结束日期"
// @param status query string false "状态"
// @param page query int false "页码" default(1)
// @param page_size query int false "每页数量" default(20)
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) GetReconciliationList(ctx context.Context, c *app.RequestContext) {
	// TODO: 实现获取对账列表
	c.JSON(consts.StatusOK, map[string]interface{}{
		"message": "Get reconciliation list - to be implemented",
	})
}

// BatchExecuteWithdrawals 批量执行提现
// @router POST /api/billing/withdrawals/batch-execute
// @summary 批量执行提现
// @description 批量自动执行已批准的提现申请
// @tags billing/withdrawal
// @accept json
// @produce json
// @param request body map[string][]int true "提现ID列表"
// @success 200 {object} map[string]interface{}
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) BatchExecuteWithdrawals(ctx context.Context, c *app.RequestContext) {
	var req struct {
		WithdrawalIDs []uint64 `json:"withdrawal_ids" validate:"required"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Invalid batch execute parameters")
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid parameters: "+err.Error()))
		return
	}

	// 批量执行
	successCount := 0
	failedCount := 0
	results := make([]map[string]interface{}, 0)

	for _, id := range req.WithdrawalIDs {
		if err := h.withdrawalService.ExecuteWithdrawal(ctx, id); err != nil {
			failedCount++
			results = append(results, map[string]interface{}{
				"id":     id,
				"status": "failed",
				"error":  err.Error(),
			})
		} else {
			successCount++
			results = append(results, map[string]interface{}{
				"id":     id,
				"status": "success",
			})
		}
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"total":        len(req.WithdrawalIDs),
		"success_count": successCount,
		"failed_count":  failedCount,
		"results":       results,
	})
}

// CallbackWithdrawal 提现回调
// @router POST /api/billing/withdrawals/callback/:provider
// @summary 提现回调
// @description 接收第三方支付平台的提现回调通知
// @tags billing/withdrawal
// @accept json
// @produce json
// @param provider path string true "支付提供商(alipay/wechat)"
// @success 200 {string} string "success"
// @failure 400 {object} bizerror.BizError
// @failure 500 {object} bizerror.BizError
func (h *WithdrawalHandler) CallbackWithdrawal(ctx context.Context, c *app.RequestContext) {
	provider := c.Param("provider")

	// 解析回调数据
	var callbackData map[string]interface{}
	if err := c.BindAndValidate(&callbackData); err != nil {
		h.logger.WithError(err).Error("[WithdrawalHandler] Invalid callback data")
		c.JSON(consts.StatusBadRequest, bizerror.NewBizError(consts.StatusBadRequest, "Invalid callback data"))
		return
	}

	// TODO: 实现回调处理逻辑
	h.logger.WithFields(logrus.Fields{
		"provider": provider,
		"data":     callbackData,
	}).Info("[WithdrawalHandler] Received withdrawal callback")

	// 返回成功响应
	c.String(consts.StatusOK, "success")
}
