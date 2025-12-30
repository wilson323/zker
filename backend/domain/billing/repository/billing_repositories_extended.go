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

package repository

import (
	"context"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// 计费系统相关Repository接口扩展
// ============================================================

// ============================================================
// SubscriptionRepository 订阅仓储接口
// 注意：这个接口假设在tenant或subscription包中已定义
// 这里仅作为引用，实际使用时需要导入正确的包
// ============================================================

type SubscriptionRepository interface {
	// GetActiveByTenant 获取租户的活跃订阅
	GetActiveByTenant(ctx context.Context, tenantID string) (*Subscription, error)
}

// Subscription 订阅信息（简化版）
type Subscription struct {
	ID              uint64
	TenantID        string
	PlanID          string
	Price           float64
	Currency        string
	BillingCycle    string
	Status          string
	StartDate       time.Time
	EndDate         *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ============================================================
// QuotaRepository 配额仓储接口
// 注意：这个接口假设在tenant或quota包中已定义
// ============================================================

type QuotaRepository interface {
	// GetByTenant 获取租户的配额信息
	GetByTenant(ctx context.Context, tenantID string) (*Quota, error)
}

// Quota 配额信息（简化版）
type Quota struct {
	ID                uint64
	TenantID          string
	TokenLimit        *int64
	StorageLimit      *int64
	RequestCountLimit *int64
	MonthlyBudget     *float64
	Currency          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ============================================================
// BillingAccountRepository 计费账户仓储接口
// ============================================================

// BillingAccountRepository 计费账户仓储接口
type BillingAccountRepository interface {
	// Create 创建计费账户
	Create(ctx context.Context, account *entity.BillingAccount) error

	// GetByID 根据ID获取账户
	GetByID(ctx context.Context, id uint64) (*entity.BillingAccount, error)

	// GetByTenant 获取租户的计费账户
	GetByTenant(ctx context.Context, tenantID string) (*entity.BillingAccount, error)

	// Update 更新账户
	Update(ctx context.Context, account *entity.BillingAccount) error

	// Delete 删除账户
	Delete(ctx context.Context, id uint64) error

	// List 分页查询账户
	List(ctx context.Context, filter *BillingAccountFilter) ([]*entity.BillingAccount, int64, error)

	// UpdateBalance 更新账户余额
	UpdateBalance(ctx context.Context, tenantID string, amount float64) error
}

// BillingAccountFilter 计费账户查询过滤器
type BillingAccountFilter struct {
	Status     *string
	Currency   *string
	PageToken  string
	PageSize   int
}

// ============================================================
// InvoiceRepository 发票仓储接口
// ============================================================

// InvoiceRepository 发票仓储接口
type InvoiceRepository interface {
	// Create 创建发票
	Create(ctx context.Context, invoice *entity.Invoice) error

	// GetByID 根据ID获取发票
	GetByID(ctx context.Context, id uint64) (*entity.Invoice, error)

	// GetByInvoiceNumber 根据发票编号获取发票
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*entity.Invoice, error)

	// Update 更新发票
	Update(ctx context.Context, invoice *entity.Invoice) error

	// Delete 删除发票
	Delete(ctx context.Context, id uint64) error

	// ListByTenant 列出租户的发票（分页）
	ListByTenant(
		ctx context.Context,
		tenantID string,
		page, pageSize int,
	) ([]*entity.Invoice, int64, error)

	// List 分页查询发票
	List(ctx context.Context, filter *InvoiceFilter) ([]*entity.Invoice, int64, error)

	// GetOverdueInvoices 获取逾期发票
	GetOverdueInvoices(ctx context.Context, beforeDate time.Time) ([]*entity.Invoice, error)

	// CountByTenantAndDate 统计租户在指定日期的发票数量
	CountByTenantAndDate(ctx context.Context, tenantID string, date time.Time) (int, error)

	// ========== 明细项操作 ==========

	// CreateLineItem 创建发票明细
	CreateLineItem(ctx context.Context, item *entity.InvoiceLineItem) error

	// GetLineItems 获取发票的所有明细
	GetLineItems(ctx context.Context, invoiceID uint64) ([]entity.InvoiceLineItem, error)

	// UpdateLineItem 更新发票明细
	UpdateLineItem(ctx context.Context, item *entity.InvoiceLineItem) error

	// DeleteLineItem 删除发票明细
	DeleteLineItem(ctx context.Context, id uint64) error
}

// InvoiceFilter 发票查询过滤器
type InvoiceFilter struct {
	TenantID     string
	Status       *string
	InvoiceType  *string
	StartDate    *time.Time
	EndDate      *time.Time
	MinAmount    *float64
	MaxAmount    *float64
	Overdue      bool
	PageToken    string
	PageSize     int
}

// ============================================================
// PaymentRepository 支付仓储接口
// ============================================================

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	// Create 创建支付记录
	Create(ctx context.Context, payment *entity.Payment) error

	// GetByID 根据ID获取支付记录
	GetByID(ctx context.Context, id uint64) (*entity.Payment, error)

	// GetByPaymentNumber 根据支付单号获取支付记录
	GetByPaymentNumber(ctx context.Context, paymentNumber string) (*entity.Payment, error)

	// Update 更新支付记录
	Update(ctx context.Context, payment *entity.Payment) error

	// Delete 删除支付记录
	Delete(ctx context.Context, id uint64) error

	// ListByTenant 列出租户的支付记录（分页）
	ListByTenant(
		ctx context.Context,
		tenantID string,
		page, pageSize int,
	) ([]*entity.Payment, int64, error)

	// List 分页查询支付记录
	List(ctx context.Context, filter *PaymentFilter) ([]*entity.Payment, int64, error)

	// CountByTenantAndTimestamp 统计租户在指定时间的支付数量
	CountByTenantAndTimestamp(ctx context.Context, tenantID string, timestamp time.Time) (int, error)

	// GetPaymentsByInvoice 获取发票关联的所有支付
	GetPaymentsByInvoice(ctx context.Context, invoiceID uint64) ([]*entity.Payment, error)
}

// PaymentFilter 支付查询过滤器
type PaymentFilter struct {
	TenantID      string
	Status        *string
	PaymentMethod *string
	Provider      *string
	StartDate     *time.Time
	EndDate       *time.Time
	MinAmount     *float64
	MaxAmount     *float64
	InvoiceID     *uint64
	PageToken     string
	PageSize      int
}
