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
// 提现管理实体 - 基于企业级支付规范
// ============================================================

// WithdrawalRequest 提现申请
type WithdrawalRequest struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:提现ID"`

	// 租户信息
	TenantID         string  `json:"tenant_id" gorm:"not null;index:idx_tenant_status,priority:1;index:idx_tenant_created,priority:1;size:64;comment:租户ID"`
	BillingAccountID uint64  `json:"billing_account_id" gorm:"not null;index:idx_account;comment:计费账户ID"`

	// 提现单号
	WithdrawalNumber string  `json:"withdrawal_number" gorm:"not null;uniqueIndex:uk_withdrawal_number;size:64;comment:提现单号"`

	// 提现金额
	Amount           float64 `json:"amount" gorm:"not null;type:decimal(12,2);index:idx_amount;comment:提现金额"`
	Fee              float64 `json:"fee" gorm:"not null;type:decimal(12,2);default:0;comment:手续费"`
	ActualAmount     float64 `json:"actual_amount" gorm:"not null;type:decimal(12,2);comment:实际到账金额"`
	Currency         string  `json:"currency" gorm:"not null;size:3;default:'CNY';comment:货币"`

	// 提现方式
	WithdrawalMethod string  `json:"withdrawal_method" gorm:"not null;size:50;index:idx_method;comment:提现方式"` // alipay, wechat, bank_transfer

	// 收款账户信息（加密存储）
	BankAccount      string  `json:"bank_account,omitempty" gorm:"type:varchar(255);comment:银行账号(加密)"`
	BankName         *string `json:"bank_name,omitempty" gorm:"size:100;comment:开户银行"`
	AccountName      *string `json:"account_name,omitempty" gorm:"size:100;comment:账户名"`
	AccountType      string  `json:"account_type" gorm:"size:20;default:'personal';comment:账户类型"` // personal, company

	// 第三方信息
	Provider         string  `json:"provider" gorm:"not null;size:50;comment:支付提供商"` // alipay, wechat, bank
	TransactionID    *string `json:"transaction_id,omitempty" gorm:"size:128;index:idx_transaction;comment:第三方交易ID"`
	BatchNo          *string `json:"batch_no,omitempty" gorm:"size:64;comment:批次号"`

	// 提现状态
	Status           string  `json:"status" gorm:"not null;size:20;default:'pending';index:idx_status;index:idx_tenant_status,priority:2;comment:提现状态"` // pending, reviewing, approved, rejected, processing, success, failed, cancelled

	// 审核信息
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty" gorm:"comment:审核时间"`
	ReviewedBy       *uint64    `json:"reviewed_by,omitempty" gorm:"comment:审核人ID"`
	ReviewComment    *string    `json:"review_comment,omitempty" gorm:"type:text;comment:审核意见"`
	AutoApproved     bool       `json:"auto_approved" gorm:"default:false;comment:是否自动审核通过"`

	// 处理信息
	ProcessedAt      *time.Time `json:"processed_at,omitempty" gorm:"comment:处理时间"`
	SuccessAt        *time.Time `json:"success_at,omitempty" gorm:"index:idx_success_at;comment:成功时间"`
	FailedReason     *string    `json:"failed_reason,omitempty" gorm:"type:text;comment:失败原因"`

	// 第三方响应
	ProviderResponse string     `json:"provider_response,omitempty" gorm:"type:text;comment:第三方响应"`

	// 风控信息
	RiskScore        *float64   `json:"risk_score,omitempty" gorm:"type:decimal(5,2);comment:风控评分"`
	RiskReason       *string    `json:"risk_reason,omitempty" gorm:"type:text;comment:风控原因"`

	// 元数据
	Metadata         string     `json:"metadata,omitempty" gorm:"type:json;comment:元数据"`
	Notes            *string    `json:"notes,omitempty" gorm:"type:text;comment:备注"`

	// 时间戳
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime;index:idx_tenant_created,priority:2;index:idx_created_at;comment:创建时间"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (WithdrawalRequest) TableName() string {
	return "withdrawal_requests"
}

// ============================================================

// WithdrawalLimit 提现限额配置
type WithdrawalLimit struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:配置ID"`

	// 限额类型
	LimitType        string  `json:"limit_type" gorm:"not null;size:50;index:idx_type;comment:限额类型"` // daily, weekly, monthly, single_transaction

	// 限额规则
	MinAmount        float64 `json:"min_amount" gorm:"not null;type:decimal(12,2);default:0;comment:最小金额"`
	MaxAmount        float64 `json:"max_amount" gorm:"not null;type:decimal(12,2);comment:最大金额"`
	DailyLimit       *float64 `json:"daily_limit,omitempty" gorm:"type:decimal(12,2);comment:每日限额"`
	WeeklyLimit      *float64 `json:"weekly_limit,omitempty" gorm:"type:decimal(12,2);comment:每周限额"`
	MonthlyLimit     *float64 `json:"monthly_limit,omitempty" gorm:"type:decimal(12,2);comment:每月限额"`

	// 提现方式
	WithdrawalMethod *string `json:"withdrawal_method,omitempty" gorm:"size:50;comment:提现方式"` // null表示全部方式

	// 费率配置
	FeeType          string  `json:"fee_type" gorm:"not null;size:20;default:'fixed';comment:费率类型"` // fixed, percentage, tiered
	FeeAmount        *float64 `json:"fee_amount,omitempty" gorm:"type:decimal(12,2);comment:固定手续费"`
	FeeRate          *float64 `json:"fee_rate,omitempty" gorm:"type:decimal(5,4);comment:费率(万分之几)"`
	MinFee           *float64 `json:"min_fee,omitempty" gorm:"type:decimal(12,2);comment:最低手续费"`
	MaxFee           *float64 `json:"max_fee,omitempty" gorm:"type:decimal(12,2);comment:最高手续费"`

	// 审核规则
	AutoApproveLimit *float64 `json:"auto_approve_limit,omitempty" gorm:"type:decimal(12,2);comment:自动审核限额"`

	// 状态
	Enabled          bool     `json:"enabled" gorm:"not null;default:true;index:idx_enabled;comment:是否启用"`

	// 时间戳
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (WithdrawalLimit) TableName() string {
	return "withdrawal_limits"
}

// ============================================================

// WithdrawalRecord 提现流水记录（用于对账）
type WithdrawalRecord struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:流水ID"`

	// 关联信息
	TenantID         string  `json:"tenant_id" gorm:"not null;index:idx_tenant_date,priority:1;size:64;comment:租户ID"`
	WithdrawalID     uint64  `json:"withdrawal_id" gorm:"not null;index:idx_withdrawal;comment:提现申请ID"`
	WithdrawalNumber string  `json:"withdrawal_number" gorm:"not null;size:64;comment:提现单号"`

	// 交易信息
	TransactionType  string  `json:"transaction_type" gorm:"not null;size:50;index:idx_type;comment:交易类型"` // freeze, unfreeze, withdraw, refund, fee

	// 金额变动
	Amount           float64 `json:"amount" gorm:"not null;type:decimal(12,2);comment:金额"`
	BalanceBefore    float64 `json:"balance_before" gorm:"not null;type:decimal(12,2);comment:变动前余额"`
	BalanceAfter     float64 `json:"balance_after" gorm:"not null;type:decimal(12,2);comment:变动后余额"`
	Currency         string  `json:"currency" gorm:"not null;size:3;default:'CNY';comment:货币"`

	// 关联凭证
	RelatedPaymentID *uint64 `json:"related_payment_id,omitempty" gorm:"index:idx_payment;comment:关联支付ID"`
	RelatedInvoiceID *uint64 `json:"related_invoice_id,omitempty" gorm:"index:idx_invoice;comment:关联发票ID"`

	// 备注
	Description      string  `json:"description" gorm:"type:text;comment:描述"`
	Metadata         string  `json:"metadata,omitempty" gorm:"type:json;comment:元数据"`

	// 时间戳
	TransactionDate  string  `json:"transaction_date" gorm:"not null;index:idx_tenant_date,priority:2;size:10;comment:交易日期(YYYY-MM-DD)"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (WithdrawalRecord) TableName() string {
	return "withdrawal_records"
}

// ============================================================

// WithdrawalReconciliation 提现对账记录
type WithdrawalReconciliation struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:对账ID"`

	// 对账批次
	ReconciliationBatch string  `json:"reconciliation_batch" gorm:"not null;uniqueIndex:uk_batch;size:64;comment:对账批次号"`

	// 对账日期
	ReconciliationDate string  `json:"reconciliation_date" gorm:"not null;index:idx_date;size:10;comment:对账日期(YYYY-MM-DD)"`

	// 对账范围
	StartDate          string  `json:"start_date" gorm:"not null;size:10;comment:开始日期"`
	EndDate            string  `json:"end_date" gorm:"not null;size:10;comment:结束日期"`

	// 对账统计
	TotalCount         int     `json:"total_count" gorm:"not null;comment:总笔数"`
	TotalAmount        float64 `json:"total_amount" gorm:"not null;type:decimal(12,2);comment:总金额"`
	SuccessCount       int     `json:"success_count" gorm:"not null;comment:成功笔数"`
	SuccessAmount      float64 `json:"success_amount" gorm:"not null;type:decimal(12,2);comment:成功金额"`
	FailedCount        int     `json:"failed_count" gorm:"not null;comment:失败笔数"`
	FailedAmount       float64 `json:"failed_amount" gorm:"not null;type:decimal(12,2);comment:失败金额"`
	PendingCount       int     `json:"pending_count" gorm:"not null;comment:待处理笔数"`
	PendingAmount      float64 `json:"pending_amount" gorm:"not null;type:decimal(12,2);comment:待处理金额"`

	// 对账结果
	Status             string  `json:"status" gorm:"not null;size:20;default:'completed';index:idx_status;comment:对账状态"` // in_progress, completed, failed

	// 差异记录
	DifferenceCount    int     `json:"difference_count" gorm:"not null;default:0;comment:差异笔数"`
	DifferenceDetails  string  `json:"difference_details,omitempty" gorm:"type:json;comment:差异明细"`

	// 执行信息
	ExecutedAt         *time.Time `json:"executed_at,omitempty" gorm:"comment:执行时间"`
	ExecutedBy         string  `json:"executed_by" gorm:"not null;size:64;comment:执行人"`
	Duration           *float64 `json:"duration,omitempty" gorm:"type:decimal(10,2);comment:执行时长(秒)"`

	// 备注
	Notes              *string `json:"notes,omitempty" gorm:"type:text;comment:备注"`

	// 时间戳
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (WithdrawalReconciliation) TableName() string {
	return "withdrawal_reconciliations"
}

// ============================================================
// 常量定义
// ============================================================

// WithdrawalStatus 提现状态枚举
const (
	WithdrawalStatusPending    = "pending"     // 待审核
	WithdrawalStatusReviewing  = "reviewing"   // 审核中
	WithdrawalStatusApproved   = "approved"    // 已批准
	WithdrawalStatusRejected   = "rejected"    // 已拒绝
	WithdrawalStatusProcessing = "processing"  // 处理中
	WithdrawalStatusSuccess    = "success"     // 成功
	WithdrawalStatusFailed     = "failed"      // 失败
	WithdrawalStatusCancelled  = "cancelled"   // 已取消
)

// WithdrawalMethod 提现方式枚举
const (
	WithdrawalMethodAlipay       = "alipay"
	WithdrawalMethodWechat       = "wechat"
	WithdrawalMethodBankTransfer = "bank_transfer"
)

// WithdrawalProvider 提现提供商枚举
const (
	WithdrawalProviderAlipay = "alipay"
	WithdrawalProviderWechat = "wechat"
	WithdrawalProviderBank   = "bank"
)

// WithdrawalTransactionType 交易类型枚举
const (
	WithdrawalTransactionTypeFreeze   = "freeze"    // 冻结金额
	WithdrawalTransactionTypeUnfreeze = "unfreeze"  // 解冻金额
	WithdrawalTransactionTypeWithdraw = "withdraw"  // 提现
	WithdrawalTransactionTypeRefund   = "refund"    // 退款
	WithdrawalTransactionTypeFee      = "fee"       // 手续费
)

// AccountType 账户类型枚举
const (
	AccountTypePersonal = "personal"
	AccountTypeCompany  = "company"
)

// FeeType 费率类型枚举
const (
	FeeTypeFixed    = "fixed"
	FeeTypePercentage = "percentage"
	FeeTypeTiered  = "tiered"
)

// LimitType 限额类型枚举
const (
	LimitTypeDaily           = "daily"
	LimitTypeWeekly          = "weekly"
	LimitTypeMonthly         = "monthly"
	LimitTypeSingleTransaction = "single_transaction"
)

// ReconciliationStatus 对账状态枚举
const (
	ReconciliationStatusInProgress = "in_progress"
	ReconciliationStatusCompleted  = "completed"
	ReconciliationStatusFailed     = "failed"
)
