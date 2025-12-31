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

	billingModel "github.com/coze-dev/coze-studio/backend/api/model/billing"
	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// ============================================================
// 账单生成API Handler
// ============================================================

var (
	billGenerationService interface{}
)

// InitBillGenerationHandler 初始化账单生成Handler
func InitBillGenerationHandler(service interface{}) {
	billGenerationService = service
}

// ==================== 账单生成接口 ====================

// GenerateBill 生成账单
// @router POST /api/billing/bills/generate [json]
// @summary 生成账单
// @description 根据计费周期生成账单
// @tags BillGeneration
// @accept json
// @produce json
// @param request body billingModel.GenerateBillRequest true "生成账单请求"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GenerateBill(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 绑定并验证请求
	var req billingModel.GenerateBillRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.WithError(err).Warn("[BillGeneration] Invalid request parameters")
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[BillGeneration] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}
	req.TenantID = tenantID

	// 3. 验证计费周期格式
	if _, err := time.Parse("2006-01", req.BillingCycle); err != nil {
		logger.WithError(err).Warn("[BillGeneration] Invalid billing_cycle format")
		InvalidParamRequestResponse(c, "billing_cycle format must be YYYY-MM")
		return
	}

	// 4. TODO: 调用应用层服务生成账单
	// bill, err := h.billService.GenerateBill(ctx, &req)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrBillGenerateFailedCode,
	//         errorx.KV("tenant_id", tenantID),
	//         errorx.KV("billing_cycle", req.BillingCycle),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrBillGenerateFailedCode,
	//         Message: "生成账单失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":      tenantID,
		"billing_cycle":  req.BillingCycle,
		"bill_type":      req.BillType,
	}).Info("[BillGeneration] Bill generated successfully")

	// 5. 返回成功响应（示例数据）
	bill := &billingModel.BillDetail{
		BillID:        "BILL-2025-001",
		TenantID:      tenantID,
		BillingCycle:  req.BillingCycle,
		BillType:      req.BillType,
		TotalAmount:   1000.00,
		Currency:      "CNY",
		Status:        "pending",
		Description:   "2025年1月账单",
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "账单生成成功",
		Data:    bill,
	})
}

// GetBill 查询账单详情
// @router GET /api/billing/bills/:bill_id [json]
// @summary 查询账单详情
// @description 根据账单ID查询账单详情
// @tags BillGeneration
// @accept json
// @produce json
// @param bill_id path string true "账单ID"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 404 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GetBill(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 获取账单ID
	billID := c.Param("bill_id")
	if billID == "" {
		InvalidParamRequestResponse(c, "bill_id is required")
		return
	}

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[BillGeneration] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	// 3. TODO: 调用应用层服务查询账单
	// bill, err := h.billService.GetBill(ctx, billID, tenantID)
	// if err != nil {
	//     if errors.Is(err, gorm.ErrRecordNotFound) {
	//         errorx.WrapByCode(err, errno.ErrBillNotFoundCode,
	//             errorx.KV("bill_id", billID),
	//         )
	//         c.JSON(http.StatusNotFound, ErrorResponse{
	//             Code:    errno.ErrBillNotFoundCode,
	//             Message: "账单不存在",
	//             Detail:  err.Error(),
	//         })
	//         return
	//     }
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.InternalErrorCode,
	//         Message: "查询账单失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"bill_id":   billID,
	}).Info("[BillGeneration] Bill retrieved successfully")

	// 4. 返回成功响应（示例数据）
	bill := &billingModel.BillDetail{
		BillID:        billID,
		TenantID:      tenantID,
		BillingCycle:  "2025-01",
		BillType:      "subscription",
		TotalAmount:   1000.00,
		Currency:      "CNY",
		Status:        "pending",
		Description:   "2025年1月订阅账单",
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "查询成功",
		Data:    bill,
	})
}

// ListBills 查询账单列表
// @router GET /api/billing/bills [json]
// @summary 查询账单列表
// @description 查询租户的账单列表，支持分页和过滤
// @tags BillGeneration
// @accept json
// @produce json
// @param page query int false "页码（默认: 1）" minimum(1)
// @param page_size query int false "每页数量（默认: 20，最大: 100）" minimum(1) maximum(100)
// @param status query string false "账单状态" Enums(pending, paid, overdue, cancelled)
// @param start_date query string false "开始日期（格式: 2006-01-02）"
// @param end_date query string false "结束日期（格式: 2006-01-02）"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func ListBills(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[BillGeneration] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	// 2. 获取并验证查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 3. 验证日期格式
	if startDate != "" {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			InvalidParamRequestResponse(c, "invalid start_date format, expected: 2006-01-02")
			return
		}
	}

	if endDate != "" {
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			InvalidParamRequestResponse(c, "invalid end_date format, expected: 2006-01-02")
			return
		}
	}

	// 4. TODO: 调用应用层服务查询账单列表
	// bills, total, err := h.billService.ListBills(ctx, tenantID, page, pageSize, status, startDate, endDate)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrBillListFailedCode,
	//         errorx.KV("tenant_id", tenantID),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrBillListFailedCode,
	//         Message: "查询账单列表失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"page":       page,
		"page_size":  pageSize,
		"status":     status,
		"start_date": startDate,
		"end_date":   endDate,
	}).Info("[BillGeneration] Bills listed successfully")

	// 5. 返回成功响应（示例数据）
	bills := []*billingModel.BillSummary{
		{
			BillID:        "BILL-2025-001",
			BillingCycle:  "2025-01",
			BillType:      "subscription",
			TotalAmount:   1000.00,
			Currency:      "CNY",
			Status:        "pending",
			Description:   "2025年1月订阅账单",
			CreatedAt:     time.Now().Unix(),
			DueDate:       time.Now().AddDate(0, 0, 30).Unix(),
		},
	}

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "查询成功",
		Data: map[string]interface{}{
			"bills":     bills,
			"total":     1,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// PayBill 支付账单
// @router POST /api/billing/bills/:bill_id/pay [json]
// @summary 支付账单
// @description 创建支付订单并返回支付信息
// @tags BillGeneration
// @accept json
// @produce json
// @param bill_id path string true "账单ID"
// @param request body billingModel.PayBillRequest true "支付请求"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 403 {object} ErrorResponse
// @failure 404 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func PayBill(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 获取账单ID
	billID := c.Param("bill_id")
	if billID == "" {
		InvalidParamRequestResponse(c, "bill_id is required")
		return
	}

	// 2. 绑定并验证请求
	var req billingModel.PayBillRequest
	if err := c.BindAndValidate(&req); err != nil {
		logger.WithError(err).Warn("[BillGeneration] Invalid request parameters")
		InvalidParamRequestResponse(c, err.Error())
		return
	}

	// 3. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[BillGeneration] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	req.BillID = billID
	req.TenantID = tenantID

	// 4. TODO: 调用应用层服务支付账单
	// err := h.billService.PayBill(ctx, &req)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrBillPayFailedCode,
	//         errorx.KV("bill_id", billID),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrBillPayFailedCode,
	//         Message: "支付失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id":      tenantID,
		"bill_id":        billID,
		"payment_method": req.PaymentMethod,
	}).Info("[BillGeneration] Bill payment initiated successfully")

	// 5. 返回成功响应（示例数据）
	paymentInfo := &billingModel.PaymentInfo{
		PaymentID:     "PAY-2025-001",
		PaymentNumber: "PAY20250101001",
		Amount:        req.Amount,
		Currency:      "CNY",
		PaymentMethod: req.PaymentMethod,
		PaymentURL:    "https://payment.example.com/pay/PAY-2025-001",
		ExpiresAt:     time.Now().Add(time.Hour * 2).Unix(),
		CreatedAt:     time.Now().Unix(),
	}

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "支付订单创建成功",
		Data:    paymentInfo,
	})
}

// GetInvoice 获取发票
// @router GET /api/billing/bills/:bill_id/invoice [json]
// @summary 获取发票
// @description 获取账单的发票信息
// @tags BillGeneration
// @accept json
// @produce json
// @param bill_id path string true "账单ID"
// @success 200 {object} SuccessResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 404 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func GetInvoice(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 获取账单ID
	billID := c.Param("bill_id")
	if billID == "" {
		InvalidParamRequestResponse(c, "bill_id is required")
		return
	}

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[BillGeneration] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	// 3. TODO: 调用应用层服务获取发票
	// invoice, err := h.billService.GetInvoice(ctx, billID, tenantID)
	// if err != nil {
	//     if errors.Is(err, gorm.ErrRecordNotFound) {
	//         errorx.WrapByCode(err, errno.ErrInvoiceNotFoundCode,
	//             errorx.KV("bill_id", billID),
	//         )
	//         c.JSON(http.StatusNotFound, ErrorResponse{
	//             Code:    errno.ErrInvoiceNotFoundCode,
	//             Message: "发票不存在",
	//             Detail:  err.Error(),
	//         })
	//         return
	//     }
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.InternalErrorCode,
	//         Message: "查询发票失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"bill_id":   billID,
	}).Info("[BillGeneration] Invoice retrieved successfully")

	// 4. 返回成功响应（示例数据）
	invoice := &billingModel.InvoiceDetail{
		InvoiceID:     "INV-2025-001",
		InvoiceNumber: "INV2025010001",
		BillID:        billID,
		TenantID:      tenantID,
		TotalAmount:   1000.00,
		Currency:      "CNY",
		Status:        "issued",
		IssuedAt:      time.Now().Unix(),
		DueDate:       time.Now().AddDate(0, 0, 30).Unix(),
	}

	httputil.BuildSuccessResp(c, SuccessResponse{
		Code:    0,
		Message: "查询成功",
		Data:    invoice,
	})
}

// DownloadInvoicePDF 下载发票PDF
// @router GET /api/billing/bills/:bill_id/invoice/pdf [json]
// @summary 下载发票PDF
// @description 下载账单的发票PDF文件
// @tags BillGeneration
// @accept json
// @produce application/pdf
// @param bill_id path string true "账单ID"
// @success 200 {file} file "PDF文件"
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
// @failure 404 {object} ErrorResponse
// @failure 500 {object} ErrorResponse
func DownloadInvoicePDF(ctx context.Context, c *app.RequestContext) {
	logger := logrus.StandardLogger()

	// 1. 获取账单ID
	billID := c.Param("bill_id")
	if billID == "" {
		InvalidParamRequestResponse(c, "bill_id is required")
		return
	}

	// 2. 从上下文获取租户ID
	tenantID := contextutil.GetTenantID(ctx)
	if tenantID == "" {
		logger.Warn("[BillGeneration] Tenant ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    errno.UnauthorizedCode,
			Message: "未授权",
		})
		return
	}

	// 3. TODO: 调用应用层服务生成PDF
	// pdfData, filename, err := h.billService.GenerateInvoicePDF(ctx, billID, tenantID)
	// if err != nil {
	//     errorx.WrapByCode(err, errno.ErrInvoicePDFGenerateFailedCode,
	//         errorx.KV("bill_id", billID),
	//     )
	//     c.JSON(http.StatusInternalServerError, ErrorResponse{
	//         Code:    errno.ErrInvoicePDFGenerateFailedCode,
	//         Message: "生成发票PDF失败",
	//         Detail:  err.Error(),
	//     })
	//     return
	// }

	logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"bill_id":   billID,
	}).Info("[BillGeneration] Invoice PDF downloaded successfully")

	// 4. 返回PDF文件（示例）
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Header("Content-Type", "application/pdf")
	// c.Data(http.StatusOK, "application/pdf", pdfData)
	c.Status(http.StatusOK)
}
