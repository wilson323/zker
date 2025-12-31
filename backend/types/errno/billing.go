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

package errno

import (
	"net/http"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// Billing: 400 000 000 ~ 400 999 999
const (
	// 实时计费错误 (400 000 000 ~ 400 009 999)
	ErrBillingChargeFailedCode      = 400000001 // 扣费失败
	ErrBillingBalanceNotFoundCode   = 400000002 // 余额不存在
	ErrBillingInsufficientBalanceCode = 400000003 // 余额不足
	ErrBillingRecordFailedCode      = 400000004 // 记录失败
	ErrBillingCalculateFailedCode   = 400000005 // 计算失败

	// 账单生成错误 (400 010 000 ~ 400 019 999)
	ErrBillGenerateFailedCode       = 400010001 // 生成账单失败
	ErrBillNotFoundCode             = 400010002 // 账单不存在
	ErrBillListFailedCode           = 400010003 // 查询账单列表失败
	ErrBillPayFailedCode            = 400010004 // 支付账单失败
	ErrBillAlreadyPaidCode          = 400010005 // 账单已支付
	ErrBillInvalidStatusCode        = 400010006 // 账单状态无效
	ErrBillExpiredCode              = 400010007 // 账单已过期

	// 发票错误 (400 020 000 ~ 400 029 999)
	ErrInvoiceGenerateFailedCode    = 400020001 // 生成发票失败
	ErrInvoiceNotFoundCode          = 400020002 // 发票不存在
	ErrInvoiceSendFailedCode        = 400020003 // 发送发票失败
	ErrInvoicePDFGenerateFailedCode = 400020004 // 生成发票PDF失败

	// 支付错误 (400 030 000 ~ 400 039 999)
	ErrPaymentCreateFailedCode      = 400030001 // 创建支付失败
	ErrPaymentNotFoundCode          = 400030002 // 支付不存在
	ErrPaymentCallbackFailedCode    = 400030003 // 支付回调处理失败
	ErrPaymentExpiredCode           = 400030004 // 支付已过期
	ErrPaymentAlreadyPaidCode       = 400030005 // 支付已完成
	ErrPaymentMethodInvalidCode     = 400030006 // 支付方式无效

	// 预算错误 (400 040 000 ~ 400 049 999)
	ErrBudgetNotFoundCode           = 400040001 // 预算不存在
	ErrBudgetExceededCode           = 400040002 // 预算超限
	ErrBudgetAlertFailedCode        = 400040003 // 预算告警失败
	ErrBudgetUpdateFailedCode       = 400040004 // 更新预算失败
	ErrBudgetInvalidCode            = 400040005 // 预算无效
)

func init() {
	// 实时计费错误注册
	code.Register(
		ErrBillingChargeFailedCode,
		"Failed to charge: tenant_id={tenant_id}, operation={operation}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBillingBalanceNotFoundCode,
		"Balance not found: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBillingInsufficientBalanceCode,
		"Insufficient balance: tenant_id={tenant_id}, required={required}, current={current}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBillingRecordFailedCode,
		"Failed to record billing: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBillingCalculateFailedCode,
		"Failed to calculate billing: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 账单生成错误注册
	code.Register(
		ErrBillGenerateFailedCode,
		"Failed to generate bill: tenant_id={tenant_id}, billing_cycle={billing_cycle}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBillNotFoundCode,
		"Bill not found: bill_id={bill_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBillListFailedCode,
		"Failed to list bills: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBillPayFailedCode,
		"Failed to pay bill: bill_id={bill_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBillAlreadyPaidCode,
		"Bill already paid: bill_id={bill_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBillInvalidStatusCode,
		"Invalid bill status: bill_id={bill_id}, current_status={current_status}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBillExpiredCode,
		"Bill expired: bill_id={bill_id}, due_date={due_date}",
		code.WithAffectStability(false),
	)

	// 发票错误注册
	code.Register(
		ErrInvoiceGenerateFailedCode,
		"Failed to generate invoice: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrInvoiceNotFoundCode,
		"Invoice not found: invoice_id={invoice_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrInvoiceSendFailedCode,
		"Failed to send invoice: invoice_id={invoice_id}, recipient={recipient}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrInvoicePDFGenerateFailedCode,
		"Failed to generate invoice PDF: invoice_id={invoice_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 支付错误注册
	code.Register(
		ErrPaymentCreateFailedCode,
		"Failed to create payment: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrPaymentNotFoundCode,
		"Payment not found: payment_id={payment_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrPaymentCallbackFailedCode,
		"Failed to handle payment callback: payment_method={payment_method}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrPaymentExpiredCode,
		"Payment expired: payment_id={payment_id}, expired_at={expired_at}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrPaymentAlreadyPaidCode,
		"Payment already completed: payment_id={payment_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrPaymentMethodInvalidCode,
		"Invalid payment method: payment_method={payment_method}",
		code.WithAffectStability(false),
	)

	// 预算错误注册
	code.Register(
		ErrBudgetNotFoundCode,
		"Budget not found: tenant_id={tenant_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBudgetExceededCode,
		"Budget exceeded: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBudgetAlertFailedCode,
		"Failed to send budget alert: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBudgetUpdateFailedCode,
		"Failed to update budget: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBudgetInvalidCode,
		"Invalid budget: tenant_id={tenant_id}, reason: {reason}",
		code.WithAffectStability(false),
	)
}

// Billing 便捷错误码变量（用于直接使用）
var (
	// 实时计费错误
	ErrBillingChargeFailed = &BaseErrorCode{
		code:       "BILLING_CHARGE_FAILED",
		message:    "Failed to charge",
		messageZH:  "扣费失败",
		messageEN:  "Failed to charge",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBillingBalanceNotFound = &BaseErrorCode{
		code:       "BILLING_BALANCE_NOT_FOUND",
		message:    "Balance not found",
		messageZH:  "余额不存在",
		messageEN:  "Balance not found",
		httpStatus: http.StatusNotFound,
	}
	ErrBillingInsufficientBalance = &BaseErrorCode{
		code:       "BILLING_INSUFFICIENT_BALANCE",
		message:    "Insufficient balance",
		messageZH:  "余额不足",
		messageEN:  "Insufficient balance",
		httpStatus: http.StatusPaymentRequired,
	}
	ErrBillingRecordFailed = &BaseErrorCode{
		code:       "BILLING_RECORD_FAILED",
		message:    "Failed to record billing",
		messageZH:  "记录失败",
		messageEN:  "Failed to record billing",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBillingCalculateFailed = &BaseErrorCode{
		code:       "BILLING_CALCULATE_FAILED",
		message:    "Failed to calculate billing",
		messageZH:  "计算失败",
		messageEN:  "Failed to calculate billing",
		httpStatus: http.StatusInternalServerError,
	}

	// 账单生成错误
	ErrBillGenerateFailed = &BaseErrorCode{
		code:       "BILL_GENERATE_FAILED",
		message:    "Failed to generate bill",
		messageZH:  "生成账单失败",
		messageEN:  "Failed to generate bill",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBillNotFound = &BaseErrorCode{
		code:       "BILL_NOT_FOUND",
		message:    "Bill not found",
		messageZH:  "账单不存在",
		messageEN:  "Bill not found",
		httpStatus: http.StatusNotFound,
	}
	ErrBillListFailed = &BaseErrorCode{
		code:       "BILL_LIST_FAILED",
		message:    "Failed to list bills",
		messageZH:  "查询账单列表失败",
		messageEN:  "Failed to list bills",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBillPayFailed = &BaseErrorCode{
		code:       "BILL_PAY_FAILED",
		message:    "Failed to pay bill",
		messageZH:  "支付账单失败",
		messageEN:  "Failed to pay bill",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBillAlreadyPaid = &BaseErrorCode{
		code:       "BILL_ALREADY_PAID",
		message:    "Bill already paid",
		messageZH:  "账单已支付",
		messageEN:  "Bill already paid",
		httpStatus: http.StatusConflict,
	}
	ErrBillInvalidStatus = &BaseErrorCode{
		code:       "BILL_INVALID_STATUS",
		message:    "Invalid bill status",
		messageZH:  "账单状态无效",
		messageEN:  "Invalid bill status",
		httpStatus: http.StatusBadRequest,
	}
	ErrBillExpired = &BaseErrorCode{
		code:       "BILL_EXPIRED",
		message:    "Bill expired",
		messageZH:  "账单已过期",
		messageEN:  "Bill expired",
		httpStatus: http.StatusGone,
	}

	// 发票错误
	ErrInvoiceGenerateFailed = &BaseErrorCode{
		code:       "INVOICE_GENERATE_FAILED",
		message:    "Failed to generate invoice",
		messageZH:  "生成发票失败",
		messageEN:  "Failed to generate invoice",
		httpStatus: http.StatusInternalServerError,
	}
	ErrInvoiceNotFound = &BaseErrorCode{
		code:       "INVOICE_NOT_FOUND",
		message:    "Invoice not found",
		messageZH:  "发票不存在",
		messageEN:  "Invoice not found",
		httpStatus: http.StatusNotFound,
	}
	ErrInvoiceSendFailed = &BaseErrorCode{
		code:       "INVOICE_SEND_FAILED",
		message:    "Failed to send invoice",
		messageZH:  "发送发票失败",
		messageEN:  "Failed to send invoice",
		httpStatus: http.StatusInternalServerError,
	}
	ErrInvoicePDFGenerateFailed = &BaseErrorCode{
		code:       "INVOICE_PDF_GENERATE_FAILED",
		message:    "Failed to generate invoice PDF",
		messageZH:  "生成发票PDF失败",
		messageEN:  "Failed to generate invoice PDF",
		httpStatus: http.StatusInternalServerError,
	}

	// 支付错误
	ErrPaymentCreateFailed = &BaseErrorCode{
		code:       "PAYMENT_CREATE_FAILED",
		message:    "Failed to create payment",
		messageZH:  "创建支付失败",
		messageEN:  "Failed to create payment",
		httpStatus: http.StatusInternalServerError,
	}
	ErrPaymentNotFound = &BaseErrorCode{
		code:       "PAYMENT_NOT_FOUND",
		message:    "Payment not found",
		messageZH:  "支付不存在",
		messageEN:  "Payment not found",
		httpStatus: http.StatusNotFound,
	}
	ErrPaymentCallbackFailed = &BaseErrorCode{
		code:       "PAYMENT_CALLBACK_FAILED",
		message:    "Failed to handle payment callback",
		messageZH:  "支付回调处理失败",
		messageEN:  "Failed to handle payment callback",
		httpStatus: http.StatusInternalServerError,
	}
	ErrPaymentExpired = &BaseErrorCode{
		code:       "PAYMENT_EXPIRED",
		message:    "Payment expired",
		messageZH:  "支付已过期",
		messageEN:  "Payment expired",
		httpStatus: http.StatusGone,
	}
	ErrPaymentAlreadyPaid = &BaseErrorCode{
		code:       "PAYMENT_ALREADY_PAID",
		message:    "Payment already completed",
		messageZH:  "支付已完成",
		messageEN:  "Payment already completed",
		httpStatus: http.StatusConflict,
	}
	ErrPaymentMethodInvalid = &BaseErrorCode{
		code:       "PAYMENT_METHOD_INVALID",
		message:    "Invalid payment method",
		messageZH:  "支付方式无效",
		messageEN:  "Invalid payment method",
		httpStatus: http.StatusBadRequest,
	}

	// 预算错误
	ErrBudgetNotFound = &BaseErrorCode{
		code:       "BUDGET_NOT_FOUND",
		message:    "Budget not found",
		messageZH:  "预算不存在",
		messageEN:  "Budget not found",
		httpStatus: http.StatusNotFound,
	}
	ErrBudgetExceeded = &BaseErrorCode{
		code:       "BUDGET_EXCEEDED",
		message:    "Budget exceeded",
		messageZH:  "预算超限",
		messageEN:  "Budget exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrBudgetAlertFailed = &BaseErrorCode{
		code:       "BUDGET_ALERT_FAILED",
		message:    "Failed to send budget alert",
		messageZH:  "预算告警失败",
		messageEN:  "Failed to send budget alert",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBudgetUpdateFailed = &BaseErrorCode{
		code:       "BUDGET_UPDATE_FAILED",
		message:    "Failed to update budget",
		messageZH:  "更新预算失败",
		messageEN:  "Failed to update budget",
		httpStatus: http.StatusInternalServerError,
	}
	ErrBudgetInvalid = &BaseErrorCode{
		code:       "BUDGET_INVALID",
		message:    "Invalid budget",
		messageZH:  "预算无效",
		messageEN:  "Invalid budget",
		httpStatus: http.StatusBadRequest,
	}
)
