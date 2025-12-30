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
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingRepository "github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 计费引擎 - 基于设计文档: 23-MultiTenant SaaS核心_租户计费系统.md
// ============================================================

// BillingEngine 计费引擎
type BillingEngine struct {
	pricingEngine     *PricingEngine
	subscriptionRepo  billingRepository.SubscriptionRepository
	quotaRepo         billingRepository.QuotaRepository
	usageLogRepo      billingRepository.TokenUsageLogRepository
	usageSummaryRepo  billingRepository.TokenUsageSummaryRepository // 企业级计费：使用汇总数据提高性能
	billingAccountRepo billingRepository.BillingAccountRepository
	invoiceRepo       billingRepository.InvoiceRepository
	db                *gorm.DB
	logger            *logrus.Logger
}

// NewBillingEngine 创建计费引擎实例
func NewBillingEngine(
	pricingEngine *PricingEngine,
	subscriptionRepo billingRepository.SubscriptionRepository,
	quotaRepo billingRepository.QuotaRepository,
	usageLogRepo billingRepository.TokenUsageLogRepository,
	usageSummaryRepo billingRepository.TokenUsageSummaryRepository,
	billingAccountRepo billingRepository.BillingAccountRepository,
	invoiceRepo billingRepository.InvoiceRepository,
	db *gorm.DB,
	logger *logrus.Logger,
) *BillingEngine {
	return &BillingEngine{
		pricingEngine:     pricingEngine,
		subscriptionRepo:  subscriptionRepo,
		quotaRepo:         quotaRepo,
		usageLogRepo:      usageLogRepo,
		usageSummaryRepo:  usageSummaryRepo,
		billingAccountRepo: billingAccountRepo,
		invoiceRepo:       invoiceRepo,
		db:                db,
		logger:            logger,
	}
}

// ============================================================
// 请求和响应类型
// ============================================================

// BillingCycle 账单周期
type BillingCycle struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Bill 账单信息
type Bill struct {
	TenantID         string            `json:"tenant_id"`
	BillingCycle     BillingCycle      `json:"billing_cycle"`
	SubscriptionFee  decimal.Decimal   `json:"subscription_fee"`
	OverageFee       decimal.Decimal   `json:"overage_fee"`
	TaxAmount        decimal.Decimal   `json:"tax_amount"`
	DiscountAmount   decimal.Decimal   `json:"discount_amount"`
	TotalAmount      decimal.Decimal   `json:"total_amount"`
	Currency         string            `json:"currency"`
	LineItems        []LineItem        `json:"line_items"`
}

// LineItem 账单明细
type LineItem struct {
	Description string          `json:"description"`
	Quantity    int             `json:"quantity"`
	UnitPrice   decimal.Decimal `json:"unit_price"`
	Amount      decimal.Decimal `json:"amount"`
	ItemType    string          `json:"item_type"`
	ItemCode    string          `json:"item_code"`
}

// Usage 使用情况
type Usage struct {
	TokenUsage    int64           `json:"token_usage"`
	StorageUsage  int64           `json:"storage_usage"`
	RequestCount  int64           `json:"request_count"`
	CustomMetrics map[string]int64 `json:"custom_metrics"`
}

// CalculateSubscriptionFeeRequest 计算订阅费用请求
type CalculateSubscriptionFeeRequest struct {
	TenantID    string     `json:"tenant_id" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     time.Time  `json:"end_date" validate:"required"`
}

// CalculateOverageRequest 计算超额费用请求
type CalculateOverageRequest struct {
	TenantID string    `json:"tenant_id" validate:"required"`
	Usage    *Usage    `json:"usage" validate:"required"`
	EndDate  time.Time `json:"end_date" validate:"required"`
}

// GenerateInvoiceRequest 生成发票请求
type GenerateInvoiceRequest struct {
	TenantID      string       `json:"tenant_id" validate:"required"`
	BillingCycle  BillingCycle `json:"billing_cycle" validate:"required"`
	InvoiceType   string       `json:"invoice_type" validate:"required"`
}

// ============================================================
// 核心方法实现
// ============================================================

// CalculateSubscriptionFee 计算订阅费用
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - startDate: 账期开始日期
//   - endDate: 账期结束日期
// 返回: 账单信息或错误
func (e *BillingEngine) CalculateSubscriptionFee(
	ctx context.Context,
	tenantID string,
	startDate, endDate time.Time,
) (*Bill, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordBillingCalculation("subscription_fee", duration)
	}()

	e.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	}).Info("[BillingEngine] Calculating subscription fee")

	// 1. 查询租户的订阅信息
	subscription, err := e.subscriptionRepo.GetActiveByTenant(ctx, tenantID)
	if err != nil {
		e.logger.WithError(err).Errorf("[BillingEngine] Failed to get subscription")
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if subscription == nil {
		return nil, fmt.Errorf("no active subscription found for tenant: %s", tenantID)
	}

	// 2. 计算账期天数
	daysInPeriod := int(endDate.Sub(startDate).Hours() / 24)
	if daysInPeriod <= 0 {
		return nil, fmt.Errorf("invalid billing period: days must be positive")
	}

	// 3. 计算订阅费用（按比例）
	monthlyFee := decimal.NewFromFloat(subscription.Price)
	daysInMonth := daysInMonth(startDate)

	// 按比例计算订阅费用
	subscriptionFee := monthlyFee.Mul(decimal.NewFromInt(int64(daysInPeriod))).
		Div(decimal.NewFromInt(int64(daysInMonth)))

	// 4. 创建账单
	bill := &Bill{
		TenantID:        tenantID,
		BillingCycle:    BillingCycle{StartDate: startDate, EndDate: endDate},
		SubscriptionFee: subscriptionFee.Round(2),
		OverageFee:      decimal.Zero,
		TaxAmount:       decimal.Zero,
		DiscountAmount:  decimal.Zero,
		TotalAmount:     subscriptionFee.Round(2),
		Currency:        subscription.Currency,
		LineItems: []LineItem{
			{
				Description: fmt.Sprintf("Subscription Fee (%s)", subscription.PlanID),
				Quantity:    daysInPeriod,
				UnitPrice:   monthlyFee.Div(decimal.NewFromInt(int64(daysInMonth))).Round(2),
				Amount:      subscriptionFee.Round(2),
				ItemType:    entity.LineItemTypeSubscription,
				ItemCode:    "subscription_monthly",
			},
		},
	}

	e.logger.WithFields(logrus.Fields{
		"tenant_id":          tenantID,
		"subscription_fee":   bill.SubscriptionFee.String(),
		"total_amount":       bill.TotalAmount.String(),
		"calculation_time_s": time.Since(startTime).Seconds(),
	}).Info("[BillingEngine] Subscription fee calculated")

	return bill, nil
}

// CalculateOverage 计算超额费用
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - usage: 使用情况
// 返回: 超额费用或错误
func (e *BillingEngine) CalculateOverage(
	ctx context.Context,
	tenantID string,
	usage *Usage,
) (decimal.Decimal, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordBillingCalculation("overage_fee", duration)
	}()

	e.logger.WithFields(logrus.Fields{
		"tenant_id":    tenantID,
		"token_usage":  usage.TokenUsage,
		"storage_usage": usage.StorageUsage,
	}).Info("[BillingEngine] Calculating overage fee")

	// 1. 查询租户的配额信息
	quota, err := e.quotaRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		e.logger.WithError(err).Errorf("[BillingEngine] Failed to get quota")
		return decimal.Zero, fmt.Errorf("failed to get quota: %w", err)
	}

	if quota == nil {
		return decimal.Zero, fmt.Errorf("no quota found for tenant: %s", tenantID)
	}

	// 2. 计算超额Token费用
	tokenOverageFee := decimal.Zero
	if quota.TokenLimit != nil && usage.TokenUsage > *quota.TokenLimit {
		overageTokens := usage.TokenUsage - *quota.TokenLimit

		// 使用定价引擎计算超额成本
		// 假设超额价格: ¥0.05/1K tokens (高于标准价格)
		overageFee := decimal.NewFromInt(overageTokens).
			Div(decimal.NewFromInt(1000)).
			Mul(decimal.NewFromFloat(0.05))

		tokenOverageFee = overageFee

		e.logger.WithFields(logrus.Fields{
			"tenant_id":      tenantID,
			"token_limit":    *quota.TokenLimit,
			"token_usage":    usage.TokenUsage,
			"overage_tokens": overageTokens,
			"overage_fee":    overageFee.String(),
		}).Info("[BillingEngine] Token overage calculated")
	}

	// 3. 计算超额存储费用
	storageOverageFee := decimal.Zero
	if quota.StorageLimit != nil && usage.StorageUsage > *quota.StorageLimit {
		overageStorage := usage.StorageUsage - *quota.StorageLimit

		// 超额存储价格: ¥0.1/GB/月
		overageFee := decimal.NewFromInt(overageStorage).
			Div(decimal.NewFromInt(1024 * 1024 * 1024)). // 转换为GB
			Mul(decimal.NewFromFloat(0.1))

		storageOverageFee = overageFee

		e.logger.WithFields(logrus.Fields{
			"tenant_id":       tenantID,
			"storage_limit":   *quota.StorageLimit,
			"storage_usage":   usage.StorageUsage,
			"overage_storage": overageStorage,
			"overage_fee":     overageFee.String(),
		}).Info("[BillingEngine] Storage overage calculated")
	}

	// 4. 计算总超额费用
	totalOverageFee := tokenOverageFee.Add(storageOverageFee).Round(2)

	e.logger.WithFields(logrus.Fields{
		"tenant_id":          tenantID,
		"token_overage_fee":  tokenOverageFee.String(),
		"storage_overage_fee": storageOverageFee.String(),
		"total_overage_fee":  totalOverageFee.String(),
		"calculation_time_s": time.Since(startTime).Seconds(),
	}).Info("[BillingEngine] Overage fee calculated")

	return totalOverageFee, nil
}

// GenerateInvoice 生成发票
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - billingCycle: 账单周期
// 返回: 发票实体或错误
func (e *BillingEngine) GenerateInvoice(
	ctx context.Context,
	tenantID string,
	billingCycle BillingCycle,
) (*entity.Invoice, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordBillingCalculation("generate_invoice", duration)
	}()

	e.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"start_date": billingCycle.StartDate.Format("2006-01-02"),
		"end_date":   billingCycle.EndDate.Format("2006-01-02"),
	}).Info("[BillingEngine] Generating invoice")

	// 1. 查询计费账户
	account, err := e.billingAccountRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		e.logger.WithError(err).Errorf("[BillingEngine] Failed to get billing account")
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	if account == nil {
		return nil, fmt.Errorf("no billing account found for tenant: %s", tenantID)
	}

	// 2. 计算订阅费用
	subscriptionBill, err := e.CalculateSubscriptionFee(
		ctx,
		tenantID,
		billingCycle.StartDate,
		billingCycle.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate subscription fee: %w", err)
	}

	// 3. 获取实际使用情况
	usage, err := e.getUsageForPeriod(ctx, tenantID, billingCycle.StartDate, billingCycle.EndDate)
	if err != nil {
		e.logger.WithError(err).Warnf("[BillingEngine] Failed to get usage, skipping overage calculation")
		usage = &Usage{} // 使用空值继续
	}

	// 4. 计算超额费用
	overageFee, err := e.CalculateOverage(ctx, tenantID, usage)
	if err != nil {
		e.logger.WithError(err).Warnf("[BillingEngine] Failed to calculate overage, using zero")
		overageFee = decimal.Zero
	}

	// 5. 生成发票编号
	invoiceNumber, err := e.generateInvoiceNumber(ctx, tenantID, billingCycle.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %w", err)
	}

	// 6. 构建费用明细
	lineItems := make([]entity.InvoiceLineItem, 0, len(subscriptionBill.LineItems)+2)

	// 添加订阅费用明细
	for i, item := range subscriptionBill.LineItems {
		unitPrice, _ := item.UnitPrice.Float64()
		amount, _ := item.Amount.Float64()

		lineItem := entity.InvoiceLineItem{
			InvoiceID:   0, // 稍后设置
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			Amount:      amount,
			ItemType:    item.ItemType,
			ItemCode:    item.ItemCode,
		}
		lineItems = append(lineItems, lineItem)

		// 设置临时ID用于后续关联
		subscriptionBill.LineItems[i].ItemCode = fmt.Sprintf("line_item_%d", i)
	}

	// 添加超额费用明细（如果有）
	if overageFee.IsPositive() {
		overageAmount, _ := overageFee.Float64()

		lineItem := entity.InvoiceLineItem{
			Description: "Overage Fee (Tokens & Storage)",
			Quantity:    1,
			UnitPrice:   overageAmount,
			Amount:      overageAmount,
			ItemType:    entity.LineItemTypeOverageTokens,
			ItemCode:    "overage",
		}
		lineItems = append(lineItems, lineItem)
	}

	// 7. 计算总金额
	subtotal := subscriptionBill.SubscriptionFee.Add(overageFee)
	taxAmount := decimal.Zero
	discountAmount := decimal.Zero
	totalAmount := subtotal.Add(taxAmount).Sub(discountAmount)

	// 8. 创建发票
	subtotalFloat, _ := subtotal.Float64()
	totalFloat, _ := totalAmount.Float64()

	invoice := &entity.Invoice{
		TenantID:        tenantID,
		BillingAccountID: account.ID,
		InvoiceNumber:   invoiceNumber,
		InvoiceType:     entity.InvoiceTypeSubscription,
		PeriodStart:     billingCycle.StartDate,
		PeriodEnd:       billingCycle.EndDate,
		Subtotal:        subtotalFloat,
		TaxAmount:       0,
		DiscountAmount:  0,
		TotalAmount:     totalFloat,
		Currency:        subscriptionBill.Currency,
		Status:          entity.InvoiceStatusDraft,
	}

	// 设置到期日（默认账期结束后15天）
	dueDate := billingCycle.EndDate.AddDate(0, 0, 15)
	invoice.DueDate = &dueDate

	// 9. 保存发票（使用事务）
	err = e.db.Transaction(func(tx *gorm.DB) error {
		// 9.1 保存发票主记录
		if err := e.invoiceRepo.Create(ctx, invoice); err != nil {
			e.logger.WithError(err).Errorf("[BillingEngine] Failed to create invoice")
			return fmt.Errorf("failed to create invoice: %w", err)
		}

		// 9.2 保存明细项
		for i := range lineItems {
			lineItems[i].InvoiceID = invoice.ID
			if err := e.invoiceRepo.CreateLineItem(ctx, &lineItems[i]); err != nil {
				e.logger.WithError(err).Errorf("[BillingEngine] Failed to create line item")
				return fmt.Errorf("failed to create line item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 10. 记录日志
	e.logger.WithFields(logrus.Fields{
		"tenant_id":       tenantID,
		"invoice_id":      invoice.ID,
		"invoice_number":  invoice.InvoiceNumber,
		"subtotal":        subtotalFloat,
		"total_amount":    totalFloat,
		"generation_time_s": time.Since(startTime).Seconds(),
	}).Info("[BillingEngine] Invoice generated successfully")

	// 11. 更新Prometheus指标
	UpdateInvoiceMetrics(tenantID, invoice.InvoiceType, invoice.Status, totalFloat)

	return invoice, nil
}

// ============================================================
// 私有辅助方法
// ============================================================

// daysInMonth 计算月份天数
func daysInMonth(date time.Time) int {
	year, month, _ := date.Date()
	// 下个月的第一天减1天就是本月最后一天
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.Add(-time.Second * 24 * 60 * 60)
	return lastOfMonth.Day()
}

// generateInvoiceNumber 生成发票编号
func (e *BillingEngine) generateInvoiceNumber(
	ctx context.Context,
	tenantID string,
	periodEnd time.Time,
) (string, error) {
	// 格式: INV-{TenantID}-{YYYYMMDD}-{序号}
	prefix := fmt.Sprintf("INV-%s-%s", tenantID, periodEnd.Format("20060102"))

	// 查询当天已有的发票数量
	count, err := e.invoiceRepo.CountByTenantAndDate(ctx, tenantID, periodEnd)
	if err != nil {
		return "", fmt.Errorf("failed to count invoices: %w", err)
	}

	// 生成序号（从0001开始）
	sequence := count + 1
	invoiceNumber := fmt.Sprintf("%s-%04d", prefix, sequence)

	return invoiceNumber, nil
}

// getUsageForPeriod 获取指定期间的使用情况（企业级计费：使用汇总数据提高性能）
func (e *BillingEngine) getUsageForPeriod(
	ctx context.Context,
	tenantID string,
	startDate, endDate time.Time,
) (*Usage, error) {
	// 1. 查询Token使用汇总数据（使用TokenUsageSummary而非TokenUsageLog）
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	summaries, err := e.usageSummaryRepo.GetByDateRange(ctx, tenantID, nil, startDateStr, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage summaries: %w", err)
	}

	// 2. 聚合汇总数据（企业级计费系统：精确的int64类型）
	usage := &Usage{
		CustomMetrics: make(map[string]int64),
	}

	for _, summary := range summaries {
		// TokenUsageSummary.TotalTokens 是 int64
		usage.TokenUsage += summary.TotalTokens
		// TokenUsageSummary.TotalRequests 是 int64（企业级计费修复）
		usage.RequestCount += summary.TotalRequests
	}

	// TODO: 添加存储使用量的查询
	// TODO: 添加自定义指标的查询

	return usage, nil
}

// ============================================================
// 账单周期计算辅助方法
// ============================================================

// GetCurrentBillingCycle 获取当前账单周期
func (e *BillingEngine) GetCurrentBillingCycle(
	ctx context.Context,
	tenantID string,
) (*BillingCycle, error) {
	// 1. 查询计费账户
	account, err := e.billingAccountRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	if account == nil {
		return nil, fmt.Errorf("no billing account found for tenant: %s", tenantID)
	}

	// 2. 根据计费周期计算当前账期
	now := time.Now()
	var startDate, endDate time.Time

	switch account.BillingCycle {
	case entity.BillingCycleMonthly:
		// 月度账期
		year, month, _ := now.Date()
		if account.BillingDayOfMonth != nil && *account.BillingDayOfMonth > 0 && *account.BillingDayOfMonth <= 28 {
			// 自定义账单日
			billingDay := int(*account.BillingDayOfMonth)
			if now.Day() >= billingDay {
				// 本月账单日已过，当前账期从本月账单日开始
				startDate = time.Date(year, month, billingDay, 0, 0, 0, 0, time.UTC)
				endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
			} else {
				// 本月账单日未到，当前账期从上月账单日开始
				startDate = time.Date(year, month, billingDay, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
				endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
			}
		} else {
			// 默认：自然月
			startDate = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
			endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
		}

	case entity.BillingCycleQuarterly:
		// 季度账期
		year, month, _ := now.Date()
		quarter := (month-1)/3 + 1
		startMonth := (quarter-1)*3 + 1
		startDate = time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 3, 0).Add(-time.Second)

	case entity.BillingCycleYearly:
		// 年度账期
		year, _, _ := now.Date()
		startDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(1, 0, 0).Add(-time.Second)

	default:
		return nil, fmt.Errorf("unsupported billing cycle: %s", account.BillingCycle)
	}

	return &BillingCycle{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}

// GetNextBillingCycle 获取下一个账单周期
func (e *BillingEngine) GetNextBillingCycle(
	ctx context.Context,
	tenantID string,
) (*BillingCycle, error) {
	current, err := e.GetCurrentBillingCycle(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 下一个账期从当前账期结束的次日开始
	startDate := current.EndDate.Add(time.Second)

	// 查询计费账户以确定账期类型
	account, err := e.billingAccountRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	var endDate time.Time
	switch account.BillingCycle {
	case entity.BillingCycleMonthly:
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)
	case entity.BillingCycleQuarterly:
		endDate = startDate.AddDate(0, 3, 0).Add(-time.Second)
	case entity.BillingCycleYearly:
		endDate = startDate.AddDate(1, 0, 0).Add(-time.Second)
	default:
		return nil, fmt.Errorf("unsupported billing cycle: %s", account.BillingCycle)
	}

	return &BillingCycle{
		StartDate: startDate,
		EndDate:   endDate,
	}, nil
}
