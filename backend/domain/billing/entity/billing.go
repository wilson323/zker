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

package entity

import (
	"time"
)

// ============================================================
// 计费账单实体 - 基于设计文档: 23-MultiTenant SaaS核心_租户计费系统.md
// ============================================================

// BillingAccount 计费账户
type BillingAccount struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:账户ID"`

	// 租户信息
	TenantID string `json:"tenant_id" gorm:"not null;uniqueIndex:uk_tenant_id;size:64;comment:租户ID"`

	// 账户状态
	Status string `json:"status" gorm:"not null;size:20;default:'active';index:idx_status;comment:账户状态"` // active, suspended, closed

	// 信用额度
	CreditLimit         float64  `json:"credit_limit" gorm:"not null;type:decimal(12,2);default:0;comment:信用额度"`
	AvailableCredit    float64  `json:"available_credit" gorm:"not null;type:decimal(12,2);default:0;comment:可用信用额度"`
	Currency           string   `json:"currency" gorm:"not null;size:3;default:'CNY';comment:货币"`
	AutoRecharge       bool     `json:"auto_recharge" gorm:"default:false;comment:是否自动充值"`
	RechargeThreshold  *float64 `json:"recharge_threshold,omitempty" gorm:"type:decimal(12,2);comment:自动充值阈值"`
	RechargeAmount     *float64 `json:"recharge_amount,omitempty" gorm:"type:decimal(12,2);comment:自动充值金额"`

	// 账单周期
	BillingCycle       string `json:"billing_cycle" gorm:"not null;size:20;default:'monthly';comment:计费周期"` // monthly, quarterly, yearly
	BillingDayOfMonth  *uint8  `json:"billing_day_of_month,omitempty" gorm:"comment:每月账单日(1-28)"`

	// 支付方式
	PaymentMethods     string  `json:"payment_methods,omitempty" gorm:"type:json;comment:支付方式列表"`
	DefaultPaymentMethod *string `json:"default_payment_method,omitempty" gorm:"size:64;comment:默认支付方式ID"`

	// 备注
	Notes              *string `json:"notes,omitempty" gorm:"type:text;comment:备注"`

	// 时间戳
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (BillingAccount) TableName() string {
	return "billing_accounts"
}

// ============================================================

// Invoice 发票/账单
type Invoice struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:发票ID"`

	// 租户信息
	TenantID         string  `json:"tenant_id" gorm:"not null;index:idx_tenant_status,priority:1;index:idx_tenant_period,priority:1;size:64;comment:租户ID"`
	BillingAccountID uint64  `json:"billing_account_id" gorm:"not null;index:idx_account;comment:计费账户ID"`

	// 发票编号
	InvoiceNumber    string  `json:"invoice_number" gorm:"not null;uniqueIndex:uk_invoice_number;size:64;comment:发票编号"`
	InvoiceType      string  `json:"invoice_type" gorm:"not null;size:20;comment:发票类型"` // subscription, overage, one_time, credit

	// 账单周期
	PeriodStart      time.Time `json:"period_start" gorm:"not null;index:idx_period_start;comment:账期开始"`
	PeriodEnd        time.Time `json:"period_end" gorm:"not null;index:idx_period_end;comment:账期结束"`

	// 金额明细
	Subtotal         float64  `json:"subtotal" gorm:"not null;type:decimal(12,2);comment:小计金额"`
	TaxAmount        float64  `json:"tax_amount" gorm:"not null;type:decimal(12,2);default:0;comment:税额"`
	DiscountAmount   float64  `json:"discount_amount" gorm:"not null;type:decimal(12,2);default:0;comment:折扣金额"`
	TotalAmount      float64  `json:"total_amount" gorm:"not null;type:decimal(12,2);index:idx_amount;comment:总金额"`
	Currency         string   `json:"currency" gorm:"not null;size:3;default:'CNY';comment:货币"`

	// 费用明细
	LineItems        string   `json:"line_items,omitempty" gorm:"type:json;comment:费用明细"`

	// 发票状态
	Status           string   `json:"status" gorm:"not null;size:20;default:'draft';index:idx_tenant_status,priority:2;index:idx_status_period;comment:发票状态"` // draft, sent, viewed, paid, partial_paid, overdue, cancelled, void

	// 支付信息
	DueDate          *time.Time `json:"due_date,omitempty" gorm:"index:idx_due_date;comment:到期日期"`
	PaidAmount       float64  `json:"paid_amount" gorm:"not null;type:decimal(12,2);default:0;comment:已支付金额"`
	PaidAt           *time.Time `json:"paid_at,omitempty" gorm:"comment:支付时间"`

	// 发票PDF
	PdfURL           *string `json:"pdf_url,omitempty" gorm:"type:text;comment:发票PDF地址"`
	PdfGeneratedAt   *time.Time `json:"pdf_generated_at,omitempty" gorm:"comment:PDF生成时间"`

	// 备注
	Metadata         string  `json:"metadata,omitempty" gorm:"type:json;comment:元数据"`
	Notes            *string `json:"notes,omitempty" gorm:"type:text;comment:备注"`

	// 时间戳
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;index:idx_created_at;comment:创建时间"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Invoice) TableName() string {
	return "invoices"
}

// ============================================================

// Payment 支付记录
type Payment struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:支付ID"`

	// 租户信息
	TenantID         string  `json:"tenant_id" gorm:"not null;index:idx_tenant_created,priority:1;size:64;comment:租户ID"`
	BillingAccountID uint64  `json:"billing_account_id" gorm:"not null;index:idx_account;comment:计费账户ID"`
	InvoiceID        *uint64 `json:"invoice_id,omitempty" gorm:"index:idx_invoice;comment:关联发票ID"`

	// 支付单号
	PaymentNumber    string  `json:"payment_number" gorm:"not null;uniqueIndex:uk_payment_number;size:64;comment:支付单号"`
	TransactionID    *string `json:"transaction_id,omitempty" gorm:"size:128;index:idx_transaction;comment:第三方交易ID"`

	// 支付金额
	Amount           float64  `json:"amount" gorm:"not null;type:decimal(12,2);index:idx_amount;comment:支付金额"`
	Currency         string   `json:"currency" gorm:"not null;size:3;default:'CNY';comment:货币"`

	// 支付方式
	PaymentMethod    string  `json:"payment_method" gorm:"not null;size:50;index:idx_method;comment:支付方式"` // alipay, wechat, credit_card, bank_transfer, paypal
	PaymentMethodID  *string `json:"payment_method_id,omitempty" gorm:"size:64;comment:支付方式ID"`

	// 支付状态
	Status           string  `json:"status" gorm:"not null;size:20;default:'pending';index:idx_status;comment:支付状态"` // pending, processing, success, failed, cancelled, refunded

	// 第三方信息
	Provider         string  `json:"provider" gorm:"not null;size:50;comment:支付提供商"` // alipay, wechat, stripe, paypal
	ProviderResponse string  `json:"provider_response,omitempty" gorm:"type:text;comment:第三方响应"`

	// 退款信息
	RefundAmount     *float64 `json:"refund_amount,omitempty" gorm:"type:decimal(12,2);comment:退款金额"`
	RefundedAt       *time.Time `json:"refunded_at,omitempty" gorm:"comment:退款时间"`
	RefundReason     *string `json:"refund_reason,omitempty" gorm:"type:text;comment:退款原因"`

	// 失败原因
	FailureReason    *string `json:"failure_reason,omitempty" gorm:"type:text;comment:失败原因"`

	// 元数据
	Metadata         string  `json:"metadata,omitempty" gorm:"type:json;comment:元数据"`
	Notes            *string `json:"notes,omitempty" gorm:"type:text;comment:备注"`

	// 时间戳
	ProcessedAt      *time.Time `json:"processed_at,omitempty" gorm:"comment:处理时间"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;index:idx_tenant_created,priority:2;comment:创建时间"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

// ============================================================

// InvoiceLineItem 发票明细项
type InvoiceLineItem struct {
	ID       uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:明细ID"`
	InvoiceID uint64 `json:"invoice_id" gorm:"not null;index:idx_invoice;comment:发票ID"`

	// 明细信息
	Description string  `json:"description" gorm:"not null;type:text;comment:描述"`
	Quantity    int     `json:"quantity" gorm:"not null;default:1;comment:数量"`
	UnitPrice   float64 `json:"unit_price" gorm:"not null;type:decimal(12,2);comment:单价"`
	Amount      float64 `json:"amount" gorm:"not null;type:decimal(12,2);comment:金额"`

	// 分类
	ItemType    string  `json:"item_type" gorm:"not null;size:50;comment:项目类型"` // subscription, overage_tokens, overage_storage, setup_fee, discount
	ItemCode    string  `json:"item_code" gorm:"not null;size:50;comment:项目代码"`

	// 元数据
	Metadata    string  `json:"metadata,omitempty" gorm:"type:json;comment:元数据"`

	// 时间戳
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (InvoiceLineItem) TableName() string {
	return "invoice_line_items"
}

// ============================================================
// 辅助类型和常量
// ============================================================

// BillingAccountStatus 账户状态枚举
const (
	BillingAccountStatusActive    = "active"
	BillingAccountStatusSuspended = "suspended"
	BillingAccountStatusClosed    = "closed"
)

// BillingCycle 计费周期枚举
const (
	BillingCycleMonthly   = "monthly"
	BillingCycleQuarterly = "quarterly"
	BillingCycleYearly    = "yearly"
)

// InvoiceType 发票类型枚举
const (
	InvoiceTypeSubscription = "subscription"
	InvoiceTypeOverage      = "overage"
	InvoiceTypeOneTime      = "one_time"
	InvoiceTypeCredit       = "credit"
)

// InvoiceStatus 发票状态枚举
const (
	InvoiceStatusDraft       = "draft"
	InvoiceStatusSent        = "sent"
	InvoiceStatusViewed      = "viewed"
	InvoiceStatusPaid        = "paid"
	InvoiceStatusPartialPaid = "partial_paid"
	InvoiceStatusOverdue     = "overdue"
	InvoiceStatusCancelled   = "cancelled"
	InvoiceStatusVoid        = "void"
)

// PaymentMethod 支付方式枚举
const (
	PaymentMethodAlipay       = "alipay"
	PaymentMethodWechat       = "wechat"
	PaymentMethodCreditCard   = "credit_card"
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodPaypal       = "paypal"
	PaymentMethodBalance      = "balance"
)

// PaymentStatus 支付状态枚举
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSuccess    = "success"
	PaymentStatusFailed     = "failed"
	PaymentStatusCancelled  = "cancelled"
	PaymentStatusRefunded   = "refunded"
)

// PaymentProvider 支付提供商枚举
const (
	PaymentProviderAlipay = "alipay"
	PaymentProviderWechat = "wechat"
	PaymentProviderStripe = "stripe"
	PaymentProviderPaypal = "paypal"
)

// LineItemType 明细项类型枚举
const (
	LineItemTypeSubscription    = "subscription"
	LineItemTypeOverageTokens   = "overage_tokens"
	LineItemTypeOverageStorage  = "overage_storage"
	LineItemTypeSetupFee        = "setup_fee"
	LineItemTypeDiscount        = "discount"
	LineItemTypeTax             = "tax"
)
