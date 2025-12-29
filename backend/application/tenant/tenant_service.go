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

	"github.com/coze-studio/backend/api/model/tenant"
	tenantentity "github.com/coze-studio/backend/domain/tenant/entity"
	tenantrepo "github.com/coze-studio/backend/domain/tenant/repository"
	tenantservice "github.com/coze-studio/backend/domain/tenant/service"
	"github.com/coze-studio/backend/infra/monitoring/metrics"
	"github.com/coze-studio/backend/pkg/errorx"
	"github.com/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-studio/backend/types/errno"
)

// TenantApplicationService 租户应用服务
type TenantApplicationService struct {
	tenantSvc     *tenantservice.TenantService
	subscriptionSvc *tenantservice.SubscriptionService
	quotaSvc      *tenantservice.QuotaService
	billingSvc    *tenantservice.BillingService
	quotaMonitor  *tenantservice.QuotaMonitorOptimized
}

// NewTenantApplicationService 创建租户应用服务
func NewTenantApplicationService(
	tenantSvc *tenantservice.TenantService,
	subscriptionSvc *tenantservice.SubscriptionService,
	quotaSvc *tenantservice.QuotaService,
	billingSvc *tenantservice.BillingService,
	quotaMonitor *tenantservice.QuotaMonitorOptimized,
) *TenantApplicationService {
	return &TenantApplicationService{
		tenantSvc:     tenantSvc,
		subscriptionSvc: subscriptionSvc,
		quotaSvc:      quotaSvc,
		billingSvc:    billingSvc,
		quotaMonitor:  quotaMonitor,
	}
}

// ==================== 租户管理用例 ====================

// CreateTenant 创建租户
func (s *TenantApplicationService) CreateTenant(ctx context.Context, req *tenant.CreateTenantRequest) (*tenant.TenantInfo, error) {
	// 1. 参数验证
	if req.TenantName == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_name is required"))
	}
	if req.TenantType == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_type is required"))
	}
	if req.ContactEmail == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "contact_email is required"))
	}

	// 2. 调用领域服务
	createReq := &tenantservice.CreateTenantRequest{
		TenantName:   req.TenantName,
		TenantType:   tenantentity.TenantType(req.TenantType),
		ContactEmail: req.ContactEmail,
		ContactPhone: req.ContactPhone,
	}

	tenantEntity, err := s.tenantSvc.CreateTenant(ctx, createReq)
	if err != nil {
		// 记录创建失败指标
		metrics.RecordTenantCreation(req.TenantType, "failure")
		return nil, err
	}

	// 3. 记录创建成功指标
	metrics.RecordTenantCreation(req.TenantType, "success")

	// 4. 返回DTO
	return s.entityToTenantInfo(tenantEntity), nil
}

// GetTenant 获取租户详情
func (s *TenantApplicationService) GetTenant(ctx context.Context, tenantID string) (*tenant.TenantDetail, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 获取租户信息
	tenantEntity, err := s.tenantSvc.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if tenantEntity == nil {
		return nil, errorx.New(errno.TenantNotFoundCode, errorx.KV("tenant_id", tenantID))
	}

	// 2. 获取配额信息
	quotas, err := s.quotaSvc.GetAllQuotaStatuses(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	quotaDTOs := make([]tenant.QuotaInfo, 0, len(quotas))
	for _, q := range quotas {
		quotaDTOs = append(quotaDTOs, s.quotaStatusToQuotaInfo(q))
	}

	return &tenant.TenantDetail{
		TenantID:         tenantEntity.TenantID,
		TenantName:       tenantEntity.TenantName,
		TenantType:       string(tenantEntity.TenantType),
		Status:           string(tenantEntity.Status),
		SubscriptionTier: string(tenantEntity.SubscriptionTier),
		ContactEmail:     tenantEntity.ContactEmail,
		ContactPhone:     tenantEntity.ContactPhone,
		Quotas:           quotaDTOs,
		CreatedAt:        tenantEntity.CreatedAt,
		UpdatedAt:        tenantEntity.UpdatedAt,
	}, nil
}

// UpdateTenant 更新租户信息
func (s *TenantApplicationService) UpdateTenant(ctx context.Context, tenantID string, req *tenant.UpdateTenantRequest) (*tenant.TenantInfo, error) {
	// 1. 获取现有租户
	tenantEntity, err := s.tenantSvc.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if tenantEntity == nil {
		return nil, errorx.New(errno.TenantNotFoundCode, errorx.KV("tenant_id", tenantID))
	}

	// 2. 应用更新
	if req.TenantName != nil {
		tenantEntity.TenantName = *req.TenantName
	}
	if req.ContactEmail != nil {
		tenantEntity.ContactEmail = *req.ContactEmail
	}
	if req.ContactPhone != nil {
		tenantEntity.ContactPhone = *req.ContactPhone
	}

	// 3. 保存更新
	if err := s.tenantSvc.Update(ctx, tenantEntity); err != nil {
		return nil, err
	}

	return s.entityToTenantInfo(tenantEntity), nil
}

// DeleteTenant 删除租户（软删除）
func (s *TenantApplicationService) DeleteTenant(ctx context.Context, tenantID string) error {
	if tenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 检查租户是否存在
	tenantEntity, err := s.tenantSvc.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	if tenantEntity == nil {
		return errorx.New(errno.TenantNotFoundCode, errorx.KV("tenant_id", tenantID))
	}

	// 2. 执行软删除
	err = s.tenantSvc.Delete(ctx, tenantID)
	if err != nil {
		metrics.RecordTenantDeletion("failure")
		return err
	}

	// 3. 记录删除成功指标
	metrics.RecordTenantDeletion("success")
	return nil
}

// ListTenants 列出租户
func (s *TenantApplicationService) ListTenants(ctx context.Context, req *tenant.ListTenantsRequest) (*tenant.ListTenantsData, error) {
	// 1. 构建查询条件
	listReq := &tenantservice.ListTenantsRequest{
		Status:           (*tenantentity.TenantStatus)(req.Status),
		TenantType:       (*tenantentity.TenantType)(req.TenantType),
		SubscriptionTier: (*tenantentity.SubscriptionTier)(req.SubscriptionTier),
		PageSize:         req.PageSize,
		PageToken:        req.PageToken,
	}

	// 2. 调用领域服务
	tenants, total, nextPageToken, err := s.tenantSvc.List(ctx, listReq)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	tenantDTOs := make([]tenant.TenantInfo, 0, len(tenants))
	for _, t := range tenants {
		tenantDTOs = append(tenantDTOs, *s.entityToTenantInfo(t))
	}

	return &tenant.ListTenantsData{
		Tenants:       tenantDTOs,
		TotalCount:    total,
		NextPageToken: nextPageToken,
	}, nil
}

// ==================== 订阅管理用例 ====================

// GetSubscription 获取订阅信息
func (s *TenantApplicationService) GetSubscription(ctx context.Context, tenantID string) (*tenant.SubscriptionInfo, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	subscription, err := s.subscriptionSvc.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, errorx.New(errno.SubscriptionNotFoundCode, errorx.KV("tenant_id", tenantID))
	}

	return s.entityToSubscriptionInfo(subscription), nil
}

// UpgradeSubscription 升级订阅等级
func (s *TenantApplicationService) UpgradeSubscription(ctx context.Context, tenantID string, req *tenant.UpgradeSubscriptionRequest) (*tenant.TenantInfo, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务升级订阅
	err := s.tenantSvc.UpgradeSubscription(ctx, tenantID, tenantentity.SubscriptionTier(req.TargetTier))
	if err != nil {
		return nil, err
	}

	// 2. 获取更新后的租户信息
	tenantEntity, err := s.tenantSvc.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	return s.entityToTenantInfo(tenantEntity), nil
}

// UpdateSubscription 更新订阅
func (s *TenantApplicationService) UpdateSubscription(ctx context.Context, tenantID string, req *tenant.UpdateSubscriptionRequest) (*tenant.SubscriptionInfo, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 获取现有订阅
	subscription, err := s.subscriptionSvc.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, errorx.New(errno.SubscriptionNotFoundCode, errorx.KV("tenant_id", tenantID))
	}

	// 2. 应用更新
	if req.BillingCycle != nil {
		subscription.BillingCycle = tenantentity.BillingCycle(*req.BillingCycle)
	}
	if req.AutoRenew != nil {
		subscription.AutoRenew = *req.AutoRenew
	}

	// 3. 保存更新
	if err := s.subscriptionSvc.Update(ctx, subscription); err != nil {
		return nil, err
	}

	return s.entityToSubscriptionInfo(subscription), nil
}

// ==================== 配额管理用例 ====================

// GetQuotas 获取配额状态
func (s *TenantApplicationService) GetQuotas(ctx context.Context, tenantID string) (*tenant.QuotasData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务获取所有配额状态
	quotaStatuses, err := s.quotaSvc.GetAllQuotaStatuses(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. 转换为DTO
	quotaDTOs := make([]tenant.QuotaInfo, 0, len(quotaStatuses))
	for _, qs := range quotaStatuses {
		quotaDTOs = append(quotaDTOs, s.quotaStatusToQuotaInfo(qs))
	}

	return &tenant.QuotasData{
		Quotas: quotaDTOs,
	}, nil
}

// CheckQuota 检查配额
func (s *TenantApplicationService) CheckQuota(ctx context.Context, tenantID string, req *tenant.CheckQuotaRequest) (*tenant.CheckQuotaData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务检查配额
	err := s.quotaSvc.CheckQuota(ctx, tenantID, tenantentity.ResourceType(req.ResourceType), req.RequiredCount)
	if err != nil {
		// 检查是否是配额超限错误
		if quotaExceededErr, ok := err.(*tenantservice.QuotaExceededError); ok {
			return &tenant.CheckQuotaData{
				Allowed:      false,
				CurrentUsage: quotaExceededErr.CurrentUsage,
				MaxLimit:     quotaExceededErr.MaxLimit,
				Remaining:    quotaExceededErr.MaxLimit - quotaExceededErr.CurrentUsage,
			}, nil
		}
		return nil, err
	}

	// 2. 获取当前配额状态
	quota, err := s.quotaSvc.GetByTenantAndResource(ctx, tenantID, tenantentity.ResourceType(req.ResourceType))
	if err != nil {
		return nil, err
	}

	return &tenant.CheckQuotaData{
		Allowed:      true,
		CurrentUsage: quota.UsedCount,
		MaxLimit:     quota.MaxLimit,
		Remaining:    quota.GetRemainingCount(),
	}, nil
}

// ConsumeQuota 消费配额
func (s *TenantApplicationService) ConsumeQuota(ctx context.Context, tenantID string, req *tenant.ConsumeQuotaRequest) (*tenant.ConsumeQuotaData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务消费配额
	previousCount, err := s.quotaSvc.ConsumeQuota(ctx, tenantID, tenantentity.ResourceType(req.ResourceType), req.Count)
	if err != nil {
		// 如果是配额超限错误，记录指标
		if errorx.IsErrorCode(err, errno.QuotaExceededCode) {
			metrics.RecordQuotaExceeded(tenantID, req.ResourceType)
		}
		return nil, err
	}

	// 2. 获取更新后的配额状态
	quota, err := s.quotaSvc.GetByTenantAndResource(ctx, tenantID, tenantentity.ResourceType(req.ResourceType))
	if err != nil {
		return nil, err
	}

	// 3. 计算使用百分比和告警级别
	usagePercent := quota.GetUsagePercent()
	alertLevel := "normal"
	if usagePercent >= 100 {
		alertLevel = "exceeded"
		metrics.RecordQuotaExceeded(tenantID, req.ResourceType)
	} else if usagePercent >= 90 {
		alertLevel = "critical"
	} else if usagePercent >= 75 {
		alertLevel = "warning"
	}

	// 4. 更新配额使用百分比指标
	metrics.UpdateQuotaUsagePercent(tenantID, req.ResourceType, alertLevel, usagePercent)

	return &tenant.ConsumeQuotaData{
		ResourceType:  req.ResourceType,
		PreviousCount: previousCount,
		CurrentCount:  quota.UsedCount,
		Remaining:     quota.GetRemainingCount(),
	}, nil
}

// RollbackQuota 回滚配额
func (s *TenantApplicationService) RollbackQuota(ctx context.Context, tenantID string, req *tenant.RollbackQuotaRequest) error {
	if tenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	return s.quotaSvc.RollbackQuota(ctx, tenantID, tenantentity.ResourceType(req.ResourceType), req.Count)
}

// ==================== 计费管理用例 ====================

// RecordUsage 记录使用量
func (s *TenantApplicationService) RecordUsage(ctx context.Context, tenantID string, req *tenant.RecordUsageRequest) error {
	if tenantID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 构建使用记录
	usage := &tenantentity.UsageRecord{
		TenantID:     tenantID,
		ResourceType: tenantentity.ResourceType(req.ResourceType),
		Action:       req.Action,
		Quantity:     req.Quantity,
		Metadata:     req.Metadata,
		RecordedAt:   time.Now().UnixMilli(),
	}

	// 2. 调用领域服务记录使用量
	return s.billingSvc.RecordUsage(ctx, usage)
}

// GetUsageSummary 获取使用量汇总
func (s *TenantApplicationService) GetUsageSummary(ctx context.Context, tenantID string, req *tenant.GetUsageSummaryRequest) (*tenant.GetUsageSummaryData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invalid start_date format"))
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invalid end_date format"))
	}

	// 2. 调用领域服务获取汇总
	summaries, err := s.billingSvc.GetUsageSummary(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	summaryDTOs := make([]tenant.UsageSummary, 0, len(summaries))
	for _, s := range summaries {
		summaryDTOs = append(summaryDTOs, tenant.UsageSummary{
			ResourceType:  string(s.ResourceType),
			TotalQuantity: s.TotalQuantity,
			ActionCounts:  s.ActionCounts,
			FirstUsage:    time.UnixMilli(s.FirstUsage).Format(time.RFC3339),
			LastUsage:     time.UnixMilli(s.LastUsage).Format(time.RFC3339),
		})
	}

	return &tenant.GetUsageSummaryData{
		Summaries: summaryDTOs,
	}, nil
}

// GenerateInvoice 生成账单
func (s *TenantApplicationService) GenerateInvoice(ctx context.Context, tenantID string, req *tenant.GenerateInvoiceRequest) (*tenant.InvoiceInfo, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invalid start_date format"))
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invalid end_date format"))
	}

	// 2. 调用领域服务生成账单
	invoice, err := s.billingSvc.GenerateInvoice(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	return s.entityToInvoiceInfo(invoice), nil
}

// GetInvoice 获取账单详情
func (s *TenantApplicationService) GetInvoice(ctx context.Context, tenantID, invoiceID string) (*tenant.InvoiceInfo, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if invoiceID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invoice_id is required"))
	}

	invoice, err := s.billingSvc.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}
	if invoice == nil {
		return nil, errorx.New(errno.InvoiceNotFoundCode, errorx.KV("invoice_id", invoiceID))
	}

	return s.entityToInvoiceInfo(invoice), nil
}

// ListInvoices 列出账单
func (s *TenantApplicationService) ListInvoices(ctx context.Context, tenantID string, status *string, limit, offset int) (*tenant.ListInvoicesData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务列出账单
	invoices, total, err := s.billingSvc.ListInvoices(ctx, tenantID, (*tenantentity.InvoiceStatus)(status), limit, offset)
	if err != nil {
		return nil, err
	}

	// 2. 转换为DTO
	invoiceDTOs := make([]tenant.InvoiceInfo, 0, len(invoices))
	for _, inv := range invoices {
		invoiceDTOs = append(invoiceDTOs, *s.entityToInvoiceInfo(inv))
	}

	return &tenant.ListInvoicesData{
		Invoices:   invoiceDTOs,
		TotalCount: total,
	}, nil
}

// PayInvoice 支付账单
func (s *TenantApplicationService) PayInvoice(ctx context.Context, tenantID, invoiceID string) (*tenant.PayInvoiceData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if invoiceID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invoice_id is required"))
	}

	// 1. 调用领域服务支付账单
	err := s.billingSvc.PayInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	// 2. 获取更新后的账单信息
	invoice, err := s.billingSvc.GetInvoice(ctx, invoiceID)
	if err != nil {
		return nil, err
	}

	return &tenant.PayInvoiceData{
		InvoiceID: invoice.InvoiceID,
		Status:    string(invoice.Status),
		PaidAt:    invoice.PaidAt,
	}, nil
}

// ==================== 监控用例 ====================

// GetMonitoringStatus 获取监控状态
func (s *TenantApplicationService) GetMonitoringStatus(ctx context.Context, tenantID string) (*tenant.MonitoringStatus, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 1. 调用领域服务获取配额监控状态
	quotaStatuses, alerts, err := s.quotaMonitor.GetTenantMonitoringStatus(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. 转换配额DTO
	quotaDTOs := make([]tenant.QuotaInfo, 0, len(quotaStatuses))
	for _, qs := range quotaStatuses {
		quotaDTOs = append(quotaDTOs, s.quotaStatusToQuotaInfo(qs))
	}

	// 3. 转换告警DTO
	alertDTOs := make([]tenant.AlertInfo, 0, len(alerts))
	for _, a := range alerts {
		alertDTOs = append(alertDTOs, tenant.AlertInfo{
			ResourceType: string(a.ResourceType),
			AlertType:    string(a.AlertLevel),
			Message:      a.Message,
		})
	}

	return &tenant.MonitoringStatus{
		TenantID: tenantID,
		Quotas:   quotaDTOs,
		Alerts:   alertDTOs,
	}, nil
}

// ==================== 转换方法 ====================

// entityToTenantInfo 实体转换为TenantInfo DTO
func (s *TenantApplicationService) entityToTenantInfo(entity *tenantentity.Tenant) *tenant.TenantInfo {
	return &tenant.TenantInfo{
		TenantID:         entity.TenantID,
		TenantName:       entity.TenantName,
		TenantType:       string(entity.TenantType),
		Status:           string(entity.Status),
		SubscriptionTier: string(entity.SubscriptionTier),
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

// entityToSubscriptionInfo 实体转换为SubscriptionInfo DTO
func (s *TenantApplicationService) entityToSubscriptionInfo(entity *tenantentity.Subscription) *tenant.SubscriptionInfo {
	return &tenant.SubscriptionInfo{
		SubscriptionID: entity.SubscriptionID,
		TenantID:       entity.TenantID,
		PlanTier:       string(entity.PlanTier),
		BillingCycle:   string(entity.BillingCycle),
		Status:         string(entity.Status),
		StartedAt:      entity.StartedAt,
		ExpiresAt:      entity.ExpiresAt,
		AutoRenew:      entity.AutoRenew,
	}
}

// quotaStatusToQuotaInfo QuotaStatus转换为QuotaInfo DTO
func (s *TenantApplicationService) quotaStatusToQuotaInfo(qs *tenantservice.QuotaStatus) tenant.QuotaInfo {
	return tenant.QuotaInfo{
		ResourceType:  string(qs.ResourceType),
		MaxLimit:      qs.MaxLimit,
		UsedCount:     qs.UsedCount,
		Remaining:     qs.GetRemainingCount(),
		UsagePercent:  qs.GetUsagePercent(),
		ResetCycle:    string(qs.ResetCycle),
		LastResetAt:   qs.LastResetAt,
		AlertLevel:    qs.GetAlertLevel(),
	}
}

// entityToInvoiceInfo 实体转换为InvoiceInfo DTO
func (s *TenantApplicationService) entityToInvoiceInfo(entity *tenantentity.Invoice) *tenant.InvoiceInfo {
	var paidAtStr string
	if entity.PaidAt != nil && *entity.PaidAt > 0 {
		paidAtStr = time.UnixMilli(*entity.PaidAt).Format(time.RFC3339)
	}

	return &tenant.InvoiceInfo{
		InvoiceID:    entity.InvoiceID,
		TenantID:     entity.TenantID,
		BillingCycle: string(entity.BillingCycle),
		StartDate:    time.UnixMilli(entity.StartDate).Format("2006-01-02"),
		EndDate:      time.UnixMilli(entity.EndDate).Format("2006-01-02"),
		TotalUsage:   entity.TotalUsage,
		TotalAmount:  entity.TotalAmount,
		Currency:     entity.Currency,
		Status:       string(entity.Status),
		DueDate:      time.UnixMilli(entity.DueDate).Format("2006-01-02"),
		CreatedAt:    time.UnixMilli(entity.CreatedAt).Format(time.RFC3339),
		PaidAt:       paidAtStr,
	}
}
