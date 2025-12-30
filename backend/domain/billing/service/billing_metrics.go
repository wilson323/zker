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

package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ============================================================
// 计费、支付、发票系统 Prometheus 指标
// ============================================================

var (
	// ========== 计费引擎指标 ==========

	// BillingCalculationsTotal 计费计算总数(Counter)
	BillingCalculationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_calculations_total",
			Help: "Total number of billing calculations",
		},
		[]string{"calculation_type"}, // subscription_fee, overage_fee, generate_invoice
	)

	// BillingCalculationDuration 计费计算耗时(Histogram)
	BillingCalculationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_calculation_duration_seconds",
			Help:    "Duration of billing calculation operations",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"calculation_type"},
	)

	// InvoiceGenerationTotal 发票生成总数(Counter)
	InvoiceGenerationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_invoice_generation_total",
			Help: "Total number of invoices generated",
		},
		[]string{"tenant_id", "invoice_type"},
	)

	// ========== 支付指标 ==========

	// PaymentOperationsTotal 支付操作总数(Counter)
	PaymentOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_payment_operations_total",
			Help: "Total number of payment operations",
		},
		[]string{"operation", "payment_method", "result"}, // operation: create, callback, query; payment_method: alipay, wechat; result: success, failure
	)

	// PaymentOperationDuration 支付操作耗时(Histogram)
	PaymentOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_payment_operation_duration_seconds",
			Help:    "Duration of payment operations",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"operation", "payment_method"},
	)

	// PaymentAmountTotal 支付总金额(Gauge)
	PaymentAmountTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_payment_amount_total",
			Help: "Total payment amount",
		},
		[]string{"tenant_id", "currency"},
	)

	// ========== 发票指标 ==========

	// InvoiceOperationsTotal 发票操作总数(Counter)
	InvoiceOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_invoice_operations_total",
			Help: "Total number of invoice operations",
		},
		[]string{"operation", "result"}, // operation: generate_pdf, send_email; result: success, failure
	)

	// InvoiceOperationDuration 发票操作耗时(Histogram)
	InvoiceOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_invoice_operation_duration_seconds",
			Help:    "Duration of invoice operations",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"operation"},
	)

	// InvoiceAmountGauge 发票金额(Gauge)
	InvoiceAmountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_invoice_amount",
			Help: "Invoice amount",
		},
		[]string{"tenant_id", "invoice_type", "status"},
	)

	// InvoiceEmailsTotal 发票邮件发送总数(Counter)
	InvoiceEmailsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_invoice_emails_total",
			Help: "Total number of invoice emails sent",
		},
		[]string{"tenant_id", "result"},
	)

	// ========== 账户指标 ==========

	// AccountBalanceGauge 账户余额(Gauge)
	AccountBalanceGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_account_balance",
			Help: "Account balance",
		},
		[]string{"tenant_id", "currency"},
	)

	// AccountCreditLimitGauge 账户信用额度(Gauge)
	AccountCreditLimitGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_account_credit_limit",
			Help: "Account credit limit",
		},
		[]string{"tenant_id", "currency"},
	)

	// ========== 业务指标 ==========

	// OverdueInvoicesGauge 逾期发票数量(Gauge)
	OverdueInvoicesGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_overdue_invoices",
			Help: "Number of overdue invoices",
		},
		[]string{"tenant_id"},
	)

	// PendingPaymentsGauge 待支付金额(Gauge)
	PendingPaymentsGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_pending_payments",
			Help: "Total pending payment amount",
		},
		[]string{"tenant_id", "currency"},
	)
)

// ============================================================
// 辅助函数 - 记录计费系统指标
// ============================================================

// RecordBillingCalculation 记录计费计算指标
func RecordBillingCalculation(calculationType string, duration float64) {
	BillingCalculationsTotal.WithLabelValues(calculationType).Inc()
	BillingCalculationDuration.WithLabelValues(calculationType).Observe(duration)
}

// UpdateInvoiceMetrics 更新发票指标
func UpdateInvoiceMetrics(tenantID, invoiceType, status string, amount float64) {
	InvoiceGenerationTotal.WithLabelValues(tenantID, invoiceType).Inc()
	InvoiceAmountGauge.WithLabelValues(tenantID, invoiceType, status).Set(amount)
}

// RecordPaymentOperation 记录支付操作指标
func RecordPaymentOperation(operation, paymentMethod string, duration float64) {
	result := "success" // 简化处理，实际应根据结果设置
	PaymentOperationsTotal.WithLabelValues(operation, paymentMethod, result).Inc()
	PaymentOperationDuration.WithLabelValues(operation, paymentMethod).Observe(duration)
}

// UpdatePaymentMetrics 更新支付指标
func UpdatePaymentMetrics(tenantID, paymentMethod, status string, amount float64) {
	result := "success"
	if status != "success" {
		result = "failure"
	}
	PaymentOperationsTotal.WithLabelValues("create", paymentMethod, result).Inc()
	PaymentAmountTotal.WithLabelValues(tenantID, "CNY").Add(amount)
}

// RecordInvoiceOperation 记录发票操作指标
func RecordInvoiceOperation(operation string, duration float64) {
	result := "success" // 简化处理
	InvoiceOperationsTotal.WithLabelValues(operation, result).Inc()
	InvoiceOperationDuration.WithLabelValues(operation).Observe(duration)
}

// UpdateAccountMetrics 更新账户指标
func UpdateAccountMetrics(tenantID, currency string, balance, creditLimit float64) {
	AccountBalanceGauge.WithLabelValues(tenantID, currency).Set(balance)
	AccountCreditLimitGauge.WithLabelValues(tenantID, currency).Set(creditLimit)
}

// UpdatePendingPaymentMetrics 更新待支付指标
func UpdatePendingPaymentMetrics(tenantID, currency string, amount float64) {
	PendingPaymentsGauge.WithLabelValues(tenantID, currency).Set(amount)
}

// IncrementOverdueInvoices 增加逾期发票计数
func IncrementOverdueInvoices(tenantID string) {
	OverdueInvoicesGauge.WithLabelValues(tenantID).Inc()
}
