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

package tenant

import (
	"context"
	"fmt"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	billingentity "github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingservice "github.com/coze-dev/coze-studio/backend/domain/billing/service"
	tenantservice "github.com/coze-dev/coze-studio/backend/domain/tenant/service"
)

// ============================================================
// 服务别名 - 用于兼容导入
// ============================================================

// TenantService 租户服务别名
type TenantService = tenantservice.TenantManagementService

// SubscriptionService 订阅服务别名
type SubscriptionService = tenantservice.SubscriptionService

// QuotaService 配额服务别名
type QuotaService = tenantservice.QuotaService

// BillingService 计费服务别名
type BillingService = billingservice.BillingEngine

// ============================================================
// QuotaStatus 配额状态
// ============================================================

// QuotaStatus 配额状态
type QuotaStatus struct {
	ResourceType  entity.ResourceType `json:"resource_type"`
	MaxLimit      int                `json:"max_limit"`
	UsedCount     int                `json:"used_count"`
	ResetCycle    entity.ResetCycle  `json:"reset_cycle"`
	LastResetAt   int64              `json:"last_reset_at"`
}

// GetRemainingCount 获取剩余配额
func (qs *QuotaStatus) GetRemainingCount() int {
	if qs.MaxLimit < 0 {
		return -1 // 无限制
	}
	remaining := qs.MaxLimit - qs.UsedCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetUsagePercent 获取使用率百分比
func (qs *QuotaStatus) GetUsagePercent() float64 {
	if qs.MaxLimit < 0 || qs.MaxLimit == 0 {
		return 0
	}
	return float64(qs.UsedCount) / float64(qs.MaxLimit) * 100
}

// GetAlertLevel 获取告警级别
func (qs *QuotaStatus) GetAlertLevel() string {
	if qs.MaxLimit < 0 {
		return "normal" // 无限制
	}

	usagePercent := qs.GetUsagePercent()
	if usagePercent >= 100 {
		return "exceeded"
	} else if usagePercent >= 90 {
		return "critical"
	} else if usagePercent >= 75 {
		return "warning"
	}
	return "normal"
}

// ============================================================
// QuotaMonitorOptimized 配额监控器
// ============================================================

// QuotaMonitorOptimized 优化的配额监控器
type QuotaMonitorOptimized struct {
	quotaRepo   repository.QuotaRepository
	billingSvc  *billingservice.BillingEngine
	alertThreshold int
}

// NewQuotaMonitorOptimized 创建配额监控器
func NewQuotaMonitorOptimized(
	quotaRepo repository.QuotaRepository,
	billingSvc *billingservice.BillingEngine,
	alertThreshold int,
) *QuotaMonitorOptimized {
	return &QuotaMonitorOptimized{
		quotaRepo:   quotaRepo,
		billingSvc:  billingSvc,
		alertThreshold: alertThreshold,
	}
}

// GetTenantMonitoringStatus 获取租户监控状态
func (m *QuotaMonitorOptimized) GetTenantMonitoringStatus(
	ctx context.Context,
	tenantID string,
) ([]*QuotaStatus, []*QuotaAlert, error) {
	// 获取所有配额
	quotas, err := m.quotaRepo.List(ctx, tenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get quotas: %w", err)
	}

	// 转换为QuotaStatus
	quotaStatuses := make([]*QuotaStatus, 0, len(quotas))
	for _, q := range quotas {
		quotaStatuses = append(quotaStatuses, &QuotaStatus{
			ResourceType:  q.ResourceType,
			MaxLimit:      q.MaxLimit,
			UsedCount:     q.UsedCount,
			ResetCycle:    q.ResetCycle,
			LastResetAt:   q.LastResetAt,
		})
	}

	// 生成告警
	alerts := make([]*QuotaAlert, 0)
	for _, qs := range quotaStatuses {
		alertLevel := qs.GetAlertLevel()
		if alertLevel == "warning" || alertLevel == "critical" || alertLevel == "exceeded" {
			alerts = append(alerts, &QuotaAlert{
				ResourceType: qs.ResourceType,
				AlertLevel:   alertLevel,
				Message:      fmt.Sprintf("Quota usage at %.1f%% for %s", qs.GetUsagePercent(), qs.ResourceType),
			})
		}
	}

	return quotaStatuses, alerts, nil
}

// QuotaAlert 配额告警
type QuotaAlert struct {
	ResourceType entity.ResourceType `json:"resource_type"`
	AlertLevel   string              `json:"alert_level"`
	Message      string              `json:"message"`
}

// ============================================================
// 服务适配器
// ============================================================

// NewTenantService 创建租户服务
func NewTenantService(repo repository.TenantRepository) *tenantservice.TenantManagementService {
	return tenantservice.NewTenantManagementService(nil, repo, nil, nil)
}

// NewSubscriptionService 创建订阅服务
func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	quotaRepo repository.QuotaRepository,
) *tenantservice.SubscriptionService {
	return tenantservice.NewSubscriptionService(subRepo, quotaRepo)
}

// NewQuotaService 创建配额服务
func NewQuotaService(quotaRepo repository.QuotaRepository) *tenantservice.QuotaService {
	return tenantservice.NewQuotaService(quotaRepo)
}

// ============================================================
// Invoice 相关类型和适配器
// ============================================================

// Invoice 发票实体别名
type Invoice = billingentity.Invoice

// UsageRecord 使用记录
type UsageRecord struct {
	TenantID     string              `json:"tenant_id"`
	ResourceType entity.ResourceType `json:"resource_type"`
	Action       string              `json:"action"`
	Quantity     int                 `json:"quantity"`
	Metadata     map[string]string   `json:"metadata,omitempty"`
	RecordedAt   int64               `json:"recorded_at"`
}

// UsageSummary 使用汇总
type UsageSummary struct {
	ResourceType  entity.ResourceType `json:"resource_type"`
	TotalQuantity int64               `json:"total_quantity"`
	ActionCounts  map[string]int      `json:"action_counts"`
	FirstUsage    int64               `json:"first_usage"`
	LastUsage     int64               `json:"last_usage"`
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	TenantName   string `json:"tenant_name"`
	TenantType   string `json:"tenant_type"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone,omitempty"`
}

// ListTenantsRequest 列出租户请求
type ListTenantsRequest struct {
	Status           *string `json:"status,omitempty"`
	TenantType       *string `json:"tenant_type,omitempty"`
	SubscriptionTier *string `json:"subscription_tier,omitempty"`
	PageSize         int     `json:"page_size"`
	PageToken        string  `json:"page_token,omitempty"`
}

// UpgradeSubscriptionRequest 升级订阅请求
type UpgradeSubscriptionRequest struct {
	TargetTier string `json:"target_tier"`
}

// UpdateSubscriptionRequest 更新订阅请求
type UpdateSubscriptionRequest struct {
	BillingCycle *string `json:"billing_cycle,omitempty"`
	AutoRenew    *bool   `json:"auto_renew,omitempty"`
}

// CheckQuotaRequestApp 应用层检查配额请求（避免与quota_app.go中的冲突）
type CheckQuotaRequestApp struct {
	ResourceType string `json:"resource_type"`
	RequiredCount int   `json:"required_count"`
}

// ConsumeQuotaRequest 消费配额请求
type ConsumeQuotaRequest struct {
	ResourceType string `json:"resource_type"`
	Count        int    `json:"count"`
}

// RollbackQuotaRequest 回滚配额请求
type RollbackQuotaRequest struct {
	ResourceType string `json:"resource_type"`
	Count        int    `json:"count"`
}

// RecordUsageRequest 记录使用量请求
type RecordUsageRequest struct {
	ResourceType string              `json:"resource_type"`
	Action       string              `json:"action"`
	Quantity     int64               `json:"quantity"`
	Metadata     map[string]string   `json:"metadata,omitempty"`
}

// GetUsageSummaryRequest 获取使用量汇总请求
type GetUsageSummaryRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// GenerateInvoiceRequest 生成账单请求
type GenerateInvoiceRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	TenantName   *string `json:"tenant_name,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	ContactPhone *string `json:"contact_phone,omitempty"`
}

// ============================================================
// 计费服务适配器
// ============================================================

// BillingServiceAdapter 计费服务适配器
type BillingServiceAdapter struct {
	engine *billingservice.BillingEngine
}

// NewBillingServiceAdapter 创建计费服务适配器
func NewBillingServiceAdapter(engine *billingservice.BillingEngine) *BillingServiceAdapter {
	return &BillingServiceAdapter{engine: engine}
}

// RecordUsage 记录使用量
func (a *BillingServiceAdapter) RecordUsage(ctx context.Context, usage *UsageRecord) error {
	// TODO: 实现使用量记录逻辑
	return fmt.Errorf("not implemented")
}

// GetUsageSummary 获取使用量汇总
func (a *BillingServiceAdapter) GetUsageSummary(
	ctx context.Context,
	tenantID string,
	startDate, endDate time.Time,
) ([]*UsageSummary, error) {
	// TODO: 实现使用量汇总逻辑
	return []*UsageSummary{}, nil
}

// GenerateInvoice 生成账单
func (a *BillingServiceAdapter) GenerateInvoice(
	ctx context.Context,
	tenantID string,
	startDate, endDate time.Time,
) (*billingentity.Invoice, error) {
	billingCycle := billingservice.BillingCycle{
		StartDate: startDate,
		EndDate:   endDate,
	}
	return a.engine.GenerateInvoice(ctx, tenantID, billingCycle)
}

// GetInvoice 获取账单详情
func (a *BillingServiceAdapter) GetInvoice(ctx context.Context, invoiceID string) (*billingentity.Invoice, error) {
	// TODO: 实现获取账单详情逻辑
	return nil, fmt.Errorf("not implemented")
}

// ListInvoices 列出账单
func (a *BillingServiceAdapter) ListInvoices(
	ctx context.Context,
	tenantID string,
	status *string,
	limit, offset int,
) ([]*billingentity.Invoice, int64, error) {
	// TODO: 实现列出账单逻辑
	return []*billingentity.Invoice{}, 0, nil
}

// PayInvoice 支付账单
func (a *BillingServiceAdapter) PayInvoice(ctx context.Context, invoiceID string) error {
	// TODO: 实现支付账单逻辑
	return fmt.Errorf("not implemented")
}

// QuotaServiceAdapter 配额服务适配器
type QuotaServiceAdapter struct {
	service *tenantservice.QuotaService
}

// NewQuotaServiceAdapter 创建配额服务适配器
func NewQuotaServiceAdapter(service *tenantservice.QuotaService) *QuotaServiceAdapter {
	return &QuotaServiceAdapter{service: service}
}

// GetAllQuotaStatuses 获取所有配额状态
func (a *QuotaServiceAdapter) GetAllQuotaStatuses(
	ctx context.Context,
	tenantID string,
) ([]*QuotaStatus, error) {
	quotas, err := a.service.GetAllQuotas(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	statuses := make([]*QuotaStatus, 0, len(quotas))
	for _, q := range quotas {
		statuses = append(statuses, &QuotaStatus{
			ResourceType:  q.ResourceType,
			MaxLimit:      q.MaxLimit,
			UsedCount:     q.UsedCount,
			ResetCycle:    q.ResetCycle,
			LastResetAt:   q.LastResetAt,
		})
	}

	return statuses, nil
}

// GetByTenantAndResource 根据租户和资源类型获取配额
func (a *QuotaServiceAdapter) GetByTenantAndResource(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
) (*entity.Quota, error) {
	return a.service.GetQuota(ctx, tenantID, resourceType)
}

// SubscriptionServiceAdapter 订阅服务适配器
type SubscriptionServiceAdapter struct {
	service *tenantservice.SubscriptionService
}

// NewSubscriptionServiceAdapter 创建订阅服务适配器
func NewSubscriptionServiceAdapter(service *tenantservice.SubscriptionService) *SubscriptionServiceAdapter {
	return &SubscriptionServiceAdapter{service: service}
}

// GetByTenant 根据租户ID获取订阅
func (a *SubscriptionServiceAdapter) GetByTenant(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	return a.service.GetSubscription(ctx, tenantID)
}

// Update 更新订阅
func (a *SubscriptionServiceAdapter) Update(ctx context.Context, sub *entity.Subscription) error {
	return a.service.UpdateSubscription(ctx, sub)
}

// TenantServiceAdapter 租户服务适配器
type TenantServiceAdapter struct {
	service *tenantservice.TenantManagementService
}

// NewTenantServiceAdapter 创建租户服务适配器
func NewTenantServiceAdapter(service *tenantservice.TenantManagementService) *TenantServiceAdapter {
	return &TenantServiceAdapter{service: service}
}

// GetByID 根据ID获取租户
func (a *TenantServiceAdapter) GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	info, err := a.service.GetTenantInfo(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return info.Tenant, nil
}

// CreateTenant 创建租户
func (a *TenantServiceAdapter) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*entity.Tenant, error) {
	// 使用UpdateTenantInfo方法的请求格式
	updateReq := &tenantservice.UpdateTenantRequest{
		TenantName: req.TenantName,
		Status:     "active", // 使用字符串而非类型转换
	}

	tenant, err := a.service.UpdateTenantInfo(ctx, req.TenantName, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	return tenant, nil
}

// List 列出租户
func (a *TenantServiceAdapter) List(ctx context.Context, req *ListTenantsRequest) ([]*entity.Tenant, int64, string, error) {
	// TODO: 实现列出租户逻辑 - 需要在TenantManagementService中添加List方法
	// 暂时返回空列表
	return []*entity.Tenant{}, 0, "", nil
}

// Delete 删除租户
func (a *TenantServiceAdapter) Delete(ctx context.Context, tenantID string) error {
	// TODO: 实现删除租户逻辑
	// 可以通过将status设置为deleted来实现软删除
	return fmt.Errorf("not implemented: use UpdateTenantInfo to set status to deleted")
}

// Update 更新租户
func (a *TenantServiceAdapter) Update(ctx context.Context, tenant *entity.Tenant) error {
	updateReq := &tenantservice.UpdateTenantRequest{
		TenantName: tenant.TenantName,
		Status:     tenant.Status, // 直接使用entity.TenantStatus类型
	}
	_, err := a.service.UpdateTenantInfo(ctx, tenant.TenantID, updateReq)
	return err
}

// UpgradeSubscription 升级订阅
func (a *TenantServiceAdapter) UpgradeSubscription(ctx context.Context, tenantID string, tier entity.SubscriptionTier) error {
	req := &tenantservice.UpgradeSubscriptionRequest{
		NewPlanTier:   tier,
		BillingCycle:  entity.BillingCycleMonthly, // 默认月付
	}
	_, err := a.service.UpgradeSubscription(ctx, tenantID, req)
	return err
}
