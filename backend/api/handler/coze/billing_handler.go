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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

// ============================================================
// 计费系统 API Handler
// ============================================================

// BillingHandler 计费处理器
type BillingHandler struct {
	billingEngine   *service.BillingEngine
	paymentService  *service.PaymentService
	invoiceService  *service.InvoiceService
	logger          *logrus.Logger
}

// NewBillingHandler 创建计费处理器
func NewBillingHandler(
	billingEngine *service.BillingEngine,
	paymentService *service.PaymentService,
	invoiceService *service.InvoiceService,
	logger *logrus.Logger,
) *BillingHandler {
	return &BillingHandler{
		billingEngine:  billingEngine,
		paymentService: paymentService,
		invoiceService: invoiceService,
		logger:         logger,
	}
}

// ============================================================
// 请求和响应类型
// ============================================================

// GenerateInvoiceRequest 生成发票请求
type GenerateInvoiceRequest struct {
	TenantID    string `json:"tenant_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate     string `json:"end_date" binding:"required"`   // YYYY-MM-DD
	InvoiceType string `json:"invoice_type" binding:"required"`
}

// GenerateInvoiceResponse 生成发票响应
type GenerateInvoiceResponse struct {
	InvoiceID     uint64  `json:"invoice_id"`
	InvoiceNumber string  `json:"invoice_number"`
	TotalAmount   float64 `json:"total_amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	DueDate       *string `json:"due_date,omitempty"`
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	TenantID      string  `json:"tenant_id" binding:"required"`
	InvoiceID     *uint64 `json:"invoice_id,omitempty"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=alipay wechat"`
	ReturnURL     *string `json:"return_url,omitempty"`
}

// SendInvoiceEmailRequest 发送发票邮件请求
type SendInvoiceEmailRequest struct {
	InvoiceID     uint64 `json:"invoice_id" binding:"required"`
	RecipientEmail string `json:"recipient_email" binding:"required,email"`
}

// APIResponse 统一API响应
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ============================================================
// 发票相关API
// ============================================================

// GenerateInvoice 生成发票
// POST /api/v1/billing/invoices/generate
func (h *BillingHandler) GenerateInvoice(c *gin.Context) {
	ctx := c.Request.Context()
	startTime := time.Now()

	// 1. 绑定请求参数
	var req GenerateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// 2. 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid start_date format, expected YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid end_date format, expected YYYY-MM-DD",
		})
		return
	}

	// 3. 验证日期范围
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "end_date must be after start_date",
		})
		return
	}

	// 4. 调用计费引擎生成发票
	billingCycle := service.BillingCycle{
		StartDate: startDate,
		EndDate:   endDate,
	}

	invoice, err := h.billingEngine.GenerateInvoice(ctx, req.TenantID, billingCycle)
	if err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to generate invoice")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to generate invoice: %v", err),
		})
		return
	}

	// 5. 构建响应
	var dueDateStr *string
	if invoice.DueDate != nil {
		dueDate := invoice.DueDate.Format("2006-01-02")
		dueDateStr = &dueDate
	}

	response := GenerateInvoiceResponse{
		InvoiceID:     invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		TotalAmount:   invoice.TotalAmount,
		Currency:      invoice.Currency,
		Status:        invoice.Status,
		DueDate:       dueDateStr,
	}

	h.logger.WithFields(logrus.Fields{
		"tenant_id":   req.TenantID,
		"invoice_id":  invoice.ID,
		"duration_ms": time.Since(startTime).Milliseconds(),
	}).Info("[BillingHandler] Invoice generated")

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Invoice generated successfully",
		Data:    response,
	})
}

// GetInvoice 获取发票详情
// GET /api/v1/billing/invoices/:invoice_id
func (h *BillingHandler) GetInvoice(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 获取发票ID
	invoiceIDStr := c.Param("invoice_id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid invoice_id",
		})
		return
	}

	// 2. 获取发票详情（通过InvoiceService）
	// TODO: 需要在InvoiceService中添加GetInvoiceWithDetails方法
	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Invoice retrieved",
		Data:    gin.H{"invoice_id": invoiceID},
	})
}

// ListInvoices 列出发票
// GET /api/v1/billing/invoices
func (h *BillingHandler) ListInvoices(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 获取查询参数
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "tenant_id is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 2. 查询发票列表
	// TODO: 需要在InvoiceService中添加ListInvoicesByTenant方法

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Invoices retrieved",
		Data:    gin.H{"tenant_id": tenantID, "page": page, "page_size": pageSize},
	})
}

// GenerateInvoicePDF 生成发票PDF
// POST /api/v1/billing/invoices/:invoice_id/pdf
func (h *BillingHandler) GenerateInvoicePDF(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 获取发票ID
	invoiceIDStr := c.Param("invoice_id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid invoice_id",
		})
		return
	}

	// 2. 生成PDF
	pdfPath, err := h.invoiceService.GenerateInvoicePDF(ctx, invoiceID)
	if err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to generate PDF")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to generate PDF: %v", err),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"invoice_id": invoiceID,
		"pdf_path":   pdfPath,
	}).Info("[BillingHandler] Invoice PDF generated")

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "PDF generated successfully",
		Data:    gin.H{"pdf_path": pdfPath},
	})
}

// SendInvoiceEmail 发送发票邮件
// POST /api/v1/billing/invoices/:invoice_id/send
func (h *BillingHandler) SendInvoiceEmail(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 获取发票ID
	invoiceIDStr := c.Param("invoice_id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid invoice_id",
		})
		return
	}

	// 2. 绑定请求参数
	var req SendInvoiceEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	req.InvoiceID = invoiceID

	// 3. 发送邮件
	if err := h.invoiceService.SendInvoiceEmail(ctx, req.InvoiceID, req.RecipientEmail); err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to send email")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to send email: %v", err),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"invoice_id":     invoiceID,
		"recipient_email": req.RecipientEmail,
	}).Info("[BillingHandler] Invoice email sent")

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Email sent successfully",
	})
}

// ============================================================
// 支付相关API
// ============================================================

// CreatePayment 创建支付
// POST /api/v1/billing/payments
func (h *BillingHandler) CreatePayment(c *gin.Context) {
	ctx := c.Request.Context()
	startTime := time.Now()

	// 1. 绑定请求参数
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// 2. 构建支付请求
	clientIP := c.ClientIP()
	paymentReq := &service.CreatePaymentRequest{
		TenantID:      req.TenantID,
		InvoiceID:     req.InvoiceID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		ClientIP:      clientIP,
		ReturnURL:     req.ReturnURL,
	}

	// 3. 创建支付
	resp, err := h.paymentService.CreatePayment(ctx, paymentReq)
	if err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to create payment")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create payment: %v", err),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"tenant_id":     req.TenantID,
		"payment_id":    resp.PaymentID,
		"payment_number": resp.PaymentNumber,
		"amount":        resp.Amount,
		"duration_ms":   time.Since(startTime).Milliseconds(),
	}).Info("[BillingHandler] Payment created")

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Payment created successfully",
		Data:    resp,
	})
}

// QueryPaymentStatus 查询支付状态
// GET /api/v1/billing/payments/:payment_id/status
func (h *BillingHandler) QueryPaymentStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 获取支付ID
	paymentIDStr := c.Param("payment_id")
	paymentID, err := strconv.ParseUint(paymentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid payment_id",
		})
		return
	}

	// 2. 查询支付状态
	payment, err := h.paymentService.QueryPaymentStatus(ctx, paymentID)
	if err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to query payment")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to query payment: %v", err),
		})
		return
	}

	if payment == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    404,
			Message: "Payment not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Payment status retrieved",
		Data:    payment,
	})
}

// HandlePaymentCallback 处理支付回调
// POST /api/v1/billing/payments/callback/:payment_method
func (h *BillingHandler) HandlePaymentCallback(c *gin.Context) {
	ctx := c.Request.Context()
	startTime := time.Now()

	// 1. 获取支付方式
	paymentMethod := c.Param("payment_method")

	// 2. 解析回调数据
	var callbackData map[string]interface{}
	if err := c.ShouldBindJSON(&callbackData); err != nil {
		// 某些支付方式使用表单数据
		if err := c.Bind(&callbackData); err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Code:    400,
				Message: "Invalid callback data",
			})
			return
		}
	}

	h.logger.WithFields(logrus.Fields{
		"payment_method": paymentMethod,
		"callback_data":  callbackData,
	}).Info("[BillingHandler] Received payment callback")

	// 3. 处理回调
	if err := h.paymentService.HandlePaymentCallback(ctx, paymentMethod, callbackData); err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to handle callback")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: "Failed to process callback",
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"duration_ms": time.Since(startTime).Milliseconds(),
	}).Info("[BillingHandler] Payment callback processed")

	// 返回成功（支付宝/微信要求返回特定格式）
	c.String(http.StatusOK, "success")
}

// ============================================================
// 账单相关API
// ============================================================

// CalculateBill 计算账单
// POST /api/v1/billing/bills/calculate
func (h *BillingHandler) CalculateBill(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 绑定请求参数
	var req struct {
		TenantID  string `json:"tenant_id" binding:"required"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// 2. 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid start_date format",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Code:    400,
			Message: "Invalid end_date format",
		})
		return
	}

	// 3. 计算订阅费用
	subscriptionBill, err := h.billingEngine.CalculateSubscriptionFee(ctx, req.TenantID, startDate, endDate)
	if err != nil {
		h.logger.WithError(err).Errorf("[BillingHandler] Failed to calculate bill")
		c.JSON(http.StatusInternalServerError, APIResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to calculate bill: %v", err),
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"tenant_id":     req.TenantID,
		"subscription_fee": subscriptionBill.TotalAmount.String(),
	}).Info("[BillingHandler] Bill calculated")

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Bill calculated successfully",
		Data:    subscriptionBill,
	})
}

// GetBillingStatus 获取计费状态
// GET /api/v1/billing/accounts/:tenant_id/status
func (h *BillingHandler) GetBillingStatus(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID := c.Param("tenant_id")

	// TODO: 实现获取计费状态的逻辑

	c.JSON(http.StatusOK, APIResponse{
		Code:    200,
		Message: "Billing status retrieved",
		Data:    gin.H{"tenant_id": tenantID, "status": "active"},
	})
}

// ============================================================
// 路由注册
// ============================================================

// RegisterBillingRoutes 注册计费相关路由
func RegisterBillingRoutes(r *gin.RouterGroup, handler *BillingHandler) {
	billing := r.Group("/billing")
	{
		// 发票相关
		invoices := billing.Group("/invoices")
		{
			invoices.POST("/generate", handler.GenerateInvoice)
			invoices.GET("/:invoice_id", handler.GetInvoice)
			invoices.GET("", handler.ListInvoices)
			invoices.POST("/:invoice_id/pdf", handler.GenerateInvoicePDF)
			invoices.POST("/:invoice_id/send", handler.SendInvoiceEmail)
		}

		// 支付相关
		payments := billing.Group("/payments")
		{
			payments.POST("", handler.CreatePayment)
			payments.GET("/:payment_id/status", handler.QueryPaymentStatus)
			payments.POST("/callback/:payment_method", handler.HandlePaymentCallback)
		}

		// 账单相关
		bills := billing.Group("/bills")
		{
			bills.POST("/calculate", handler.CalculateBill)
		}

		// 账户相关
		accounts := billing.Group("/accounts")
		{
			accounts.GET("/:tenant_id/status", handler.GetBillingStatus)
		}
	}
}
