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

	"github.com/google/uuid"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// BillingService 计费统计服务
type BillingService struct {
	quotaRepo    repository.QuotaRepository
	usageLogRepo repository.UsageLogRepository
	invoiceRepo  repository.InvoiceRepository
	subscription repository.SubscriptionRepository
}

// UsageLog 使用日志
type UsageLog struct {
	LogID        string    `json:"log_id"`
	TenantID     string    `json:"tenant_id"`
	ResourceType entity.ResourceType `json:"resource_type"`
	Action       string    `json:"action"` // "create", "update", "delete", "query"
	Quantity     int       `json:"quantity"`
	Timestamp    time.Time `json:"timestamp"`
	Metadata     string    `json:"metadata"` // JSON格式的额外信息
}

// Invoice 账单
type Invoice struct {
	InvoiceID      string              `json:"invoice_id"`
	TenantID       string              `json:"tenant_id"`
	SubscriptionID string              `json:"subscription_id"`
	BillingCycle   string              `json:"billing_cycle"` // "monthly", "yearly"
	StartDate      time.Time           `json:"start_date"`
	EndDate        time.Time           `json:"end_date"`
	TotalUsage     int                 `json:"total_usage"`
	TotalAmount    float64             `json:"total_amount"`
	Currency       string              `json:"currency"` // "USD", "CNY", etc.
	Status         entity.InvoiceStatus `json:"status"` // "pending", "paid", "overdue", "cancelled"
	CreatedAt      time.Time           `json:"created_at"`
	PaidAt         *time.Time          `json:"paid_at,omitempty"`
	DueDate        time.Time           `json:"due_date"`
}

// UsageSummary 使用量汇总
type UsageSummary struct {
	TenantID       string
	ResourceType   entity.ResourceType
	TotalQuantity  int
	ActionCounts   map[string]int // 每种操作类型的计数
	FirstUsage     time.Time
	LastUsage      time.Time
}

// NewBillingService 创建计费服务实例
func NewBillingService(
	quotaRepo repository.QuotaRepository,
	usageLogRepo repository.UsageLogRepository,
	invoiceRepo repository.InvoiceRepository,
	subscription repository.SubscriptionRepository,
) *BillingService {
	return &BillingService{
		quotaRepo:    quotaRepo,
		usageLogRepo: usageLogRepo,
		invoiceRepo:  invoiceRepo,
		subscription: subscription,
	}
}

// RecordUsage 记录使用量
func (s *BillingService) RecordUsage(ctx context.Context, log *UsageLog) error {
	if log.LogID == "" {
		log.LogID = uuid.New().String()
	}

	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	return s.usageLogRepo.Create(ctx, log)
}

// RecordUsageBatch 批量记录使用量
func (s *BillingService) RecordUsageBatch(ctx context.Context, logs []*UsageLog) error {
	for _, log := range logs {
		if log.LogID == "" {
			log.LogID = uuid.New().String()
		}

		if log.Timestamp.IsZero() {
			log.Timestamp = time.Now()
		}
	}

	return s.usageLogRepo.CreateBatch(ctx, logs)
}

// GetUsageSummary 获取使用量汇总
func (s *BillingService) GetUsageSummary(ctx context.Context, tenantID string, startDate, endDate time.Time) ([]*UsageSummary, error) {
	logs, err := s.usageLogRepo.GetByTenantAndDateRange(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	summaryMap := make(map[entity.ResourceType]*UsageSummary)

	for _, log := range logs {
		if summary, ok := summaryMap[log.ResourceType]; ok {
			summary.TotalQuantity += log.Quantity
			summary.ActionCounts[log.Action]++

			if log.Timestamp.Before(summary.FirstUsage) {
				summary.FirstUsage = log.Timestamp
			}
			if log.Timestamp.After(summary.LastUsage) {
				summary.LastUsage = log.Timestamp
			}
		} else {
			summaryMap[log.ResourceType] = &UsageSummary{
				TenantID:      tenantID,
				ResourceType:  log.ResourceType,
				TotalQuantity: log.Quantity,
				ActionCounts:  map[string]int{log.Action: 1},
				FirstUsage:    log.Timestamp,
				LastUsage:     log.Timestamp,
			}
		}
	}

	summaries := make([]*UsageSummary, 0, len(summaryMap))
	for _, summary := range summaryMap {
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// GenerateInvoice 生成账单
func (s *BillingService) GenerateInvoice(ctx context.Context, tenantID string, startDate, endDate time.Time) (*Invoice, error) {
	// 1. 检查是否已存在账单
	existing, err := s.invoiceRepo.GetByTenantAndDateRange(ctx, tenantID, startDate, endDate)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("invoice already exists for this period")
	}

	// 2. 获取订阅信息
	subscription, err := s.subscription.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	// 3. 获取使用记录
	logs, err := s.usageLogRepo.GetByTenantAndDateRange(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage logs: %w", err)
	}

	// 4. 统计总使用量
	totalUsage := 0
	for _, log := range logs {
		totalUsage += log.Quantity
	}

	// 5. 计算费用
	totalAmount := s.calculateAmount(ctx, subscription, logs)

	// 6. 确定账单周期
	billingCycle := "monthly"
	if subscription != nil && subscription.BillingCycle == entity.BillingCycleYearly {
		billingCycle = "yearly"
	}

	// 7. 计算到期日期（30天后）
	dueDate := time.Now().AddDate(0, 0, 30)

	// 8. 创建账单
	invoice := &Invoice{
		InvoiceID:      uuid.New().String(),
		TenantID:       tenantID,
		SubscriptionID: getSubscriptionID(subscription),
		BillingCycle:   billingCycle,
		StartDate:      startDate,
		EndDate:        endDate,
		TotalUsage:     totalUsage,
		TotalAmount:    totalAmount,
		Currency:       "CNY", // 默认人民币
		Status:         entity.InvoiceStatusPending,
		CreatedAt:      time.Now(),
		DueDate:        dueDate,
	}

	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoice, nil
}

// calculateAmount 计算费用
func (s *BillingService) calculateAmount(ctx context.Context, subscription *entity.Subscription, logs []*UsageLog) float64 {
	if subscription == nil {
		// 免费版，按默认费率计算
		return s.calculateByDefaultRate(logs)
	}

	// 根据订阅等级和使用量计算费用
	// 这里简化实现，实际应结合定价策略
	baseRate := s.getBaseRate(subscription.PlanTier)
	totalQuantity := 0

	for _, log := range logs {
		totalQuantity += log.Quantity
	}

	return float64(totalQuantity) * baseRate
}

// calculateByDefaultRate 按默认费率计算
func (s *BillingService) calculateByDefaultRate(logs []*UsageLog) float64 {
	totalQuantity := 0
	for _, log := range logs {
		totalQuantity += log.Quantity
	}

	// 免费版每1000次0.1元
	return float64(totalQuantity) / 1000.0 * 0.1
}

// getBaseRate 获取基础费率
func (s *BillingService) getBaseRate(tier entity.SubscriptionTier) float64 {
	switch tier {
	case entity.SubscriptionTierFree:
		return 0.0001 // 免费版：每次0.0001元
	case entity.SubscriptionTierPro:
		return 0.00005 // 专业版：每次0.00005元
	case entity.SubscriptionTierEnterprise:
		return 0.00002 // 企业版：每次0.00002元
	default:
		return 0.0001
	}
}

// GetInvoice 获取账单详情
func (s *BillingService) GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error) {
	return s.invoiceRepo.GetByID(ctx, invoiceID)
}

// ListInvoices 列出租户的账单
func (s *BillingService) ListInvoices(ctx context.Context, tenantID string, limit, offset int) ([]*Invoice, int64, error) {
	return s.invoiceRepo.GetByTenant(ctx, tenantID, limit, offset)
}

// PayInvoice 支付账单
func (s *BillingService) PayInvoice(ctx context.Context, invoiceID string) error {
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if invoice.Status != entity.InvoiceStatusPending {
		return fmt.Errorf("invoice is not in pending status")
	}

	now := time.Now()
	invoice.Status = entity.InvoiceStatusPaid
	invoice.PaidAt = &now

	return s.invoiceRepo.Update(ctx, invoice)
}

// CancelInvoice 取消账单
func (s *BillingService) CancelInvoice(ctx context.Context, invoiceID string) error {
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if invoice.Status != entity.InvoiceStatusPending {
		return fmt.Errorf("can only cancel pending invoices")
	}

	invoice.Status = entity.InvoiceStatusCancelled

	return s.invoiceRepo.Update(ctx, invoice)
}

// GenerateMonthlyInvoices 生成月度账单（定时任务）
func (s *BillingService) GenerateMonthlyInvoices(ctx context.Context) error {
	// 获取所有订阅
	subscriptions, err := s.subscription.GetAll(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	// 上个月的第一天
	startDate := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local)
	// 上个月最后一天
	endDate := startDate.AddDate(0, 1, -1)

	for _, sub := range subscriptions {
		// 只处理月度订阅
		if sub.BillingCycle == entity.BillingCycleMonthly {
			_, err := s.GenerateInvoice(ctx, sub.TenantID, startDate, endDate)
			if err != nil {
				// 记录错误但继续处理其他订阅
				continue
			}
		}
	}

	return nil
}

// getSubscriptionID 获取订阅ID
func getSubscriptionID(subscription *entity.Subscription) string {
	if subscription == nil {
		return ""
	}
	return subscription.SubscriptionID
}
