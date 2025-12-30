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
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ============================================================
// 预算管理服务 - 提供预算配置的完整CRUD和查询功能
// ============================================================

// BudgetManagementService 预算管理服务
type BudgetManagementService struct {
	budgetRepo   repository.BudgetSettingsRepository
	alertRepo    repository.BudgetAlertRepository
	summaryRepo  repository.TokenUsageSummaryRepository
	alertService interface{} // *BudgetAlertService - 使用interface{}避免循环依赖
	logger       *logrus.Logger
}

// NewBudgetManagementService 创建预算管理服务实例
func NewBudgetManagementService(
	budgetRepo repository.BudgetSettingsRepository,
	alertRepo repository.BudgetAlertRepository,
	summaryRepo repository.TokenUsageSummaryRepository,
	alertService interface{}, // *BudgetAlertService
	logger *logrus.Logger,
) *BudgetManagementService {
	if logger == nil {
		logger = logrus.New()
	}
	return &BudgetManagementService{
		budgetRepo:   budgetRepo,
		alertRepo:    alertRepo,
		summaryRepo:  summaryRepo,
		alertService: alertService,
		logger:       logger,
	}
}

// ============================================================
// DTO定义
// ============================================================

// BudgetSettingsDTO 预算配置DTO
type BudgetSettingsDTO struct {
	TenantID     string   `json:"tenant_id"`
	BudgetType   string   `json:"budget_type"`   // monthly, quarterly, yearly
	BudgetAmount float64  `json:"budget_amount"`
	Currency     string   `json:"currency"`

	// 告警配置
	AlertThreshold1 int      `json:"alert_threshold_1"` // 80%
	AlertThreshold2 int      `json:"alert_threshold_2"` // 95%
	HardCapEnabled  bool     `json:"hard_cap_enabled"`
	HardCapAmount   *float64 `json:"hard_cap_amount,omitempty"`

	// 降级配置
	AutoDowngradeEnabled bool    `json:"auto_downgrade_enabled"`
	DowngradeThreshold   int     `json:"downgrade_threshold"` // 90%
	OriginalModel        *string `json:"original_model,omitempty"`
	FallbackModel        *string `json:"fallback_model,omitempty"`

	// 通知配置
	NotificationChannels   []string `json:"notification_channels"`
	NotificationRecipients []string `json:"notification_recipients"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// CreateBudgetRequest 创建预算请求
type CreateBudgetRequest struct {
	TenantID string `json:"tenant_id" validate:"required"`

	// 预算配置
	BudgetType   string  `json:"budget_type" validate:"required,oneof=monthly quarterly yearly"`
	BudgetAmount float64 `json:"budget_amount" validate:"required,gt=0"`
	Currency     string  `json:"currency" validate:"required,len=3"`

	// 告警配置
	AlertThreshold1 int     `json:"alert_threshold_1" validate:"required,min=0,max=100"`
	AlertThreshold2 int     `json:"alert_threshold_2" validate:"required,min=0,max=100"`
	HardCapEnabled  bool    `json:"hard_cap_enabled"`
	HardCapAmount   *float64 `json:"hard_cap_amount,omitempty" validate:"omitempty,gt=0"`

	// 降级配置
	AutoDowngradeEnabled bool    `json:"auto_downgrade_enabled"`
	DowngradeThreshold   int     `json:"downgrade_threshold" validate:"required,min=0,max=100"`
	OriginalModel        *string `json:"original_model,omitempty" validate:"omitempty,min=1,max=100"`
	FallbackModel        *string `json:"fallback_model,omitempty" validate:"omitempty,min=1,max=100"`

	// 通知配置
	NotificationChannels   []string `json:"notification_channels" validate:"required,dive,oneof=email webhook sms"`
	NotificationRecipients []string `json:"notification_recipients" validate:"required,dive,email"`
}

// UpdateBudgetRequest 更新预算请求
type UpdateBudgetRequest struct {
	BudgetAmount          *float64 `json:"budget_amount,omitempty" validate:"omitempty,gt=0"`
	AlertThreshold1       *int     `json:"alert_threshold_1,omitempty" validate:"omitempty,min=0,max=100"`
	AlertThreshold2       *int     `json:"alert_threshold_2,omitempty" validate:"omitempty,min=0,max=100"`
	HardCapEnabled        *bool    `json:"hard_cap_enabled,omitempty"`
	HardCapAmount         *float64 `json:"hard_cap_amount,omitempty" validate:"omitempty,gt=0"`
	AutoDowngradeEnabled  *bool    `json:"auto_downgrade_enabled,omitempty"`
	DowngradeThreshold    *int     `json:"downgrade_threshold,omitempty" validate:"omitempty,min=0,max=100"`
	OriginalModel         *string  `json:"original_model,omitempty" validate:"omitempty,min=1,max=100"`
	FallbackModel         *string  `json:"fallback_model,omitempty" validate:"omitempty,min=1,max=100"`
	NotificationChannels  []string `json:"notification_channels,omitempty" validate:"omitempty,dive,oneof=email webhook sms"`
	NotificationRecipients []string `json:"notification_recipients,omitempty" validate:"omitempty,dive,email"`
}

// AlertFilter 告警查询过滤器
type AlertFilter struct {
	AlertType          *string    `json:"alert_type,omitempty"`
	AlertLevel         *string    `json:"alert_level,omitempty"`
	NotificationStatus *string    `json:"notification_status,omitempty"`
	StartDate          *time.Time `json:"start_date,omitempty"`
	EndDate            *time.Time `json:"end_date,omitempty"`
	PageToken          string     `json:"page_token,omitempty"`
	PageSize           int        `json:"page_size,omitempty"`
}

// BudgetAlertDTO 预算告警DTO
type BudgetAlertDTO struct {
	ID                  uint64    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	AlertType           string    `json:"alert_type"`
	BudgetAmount        float64   `json:"budget_amount"`
	UsedAmount          float64   `json:"used_amount"`
	UsagePercent        float64   `json:"usage_percent"`
	AlertLevel          string    `json:"alert_level"`
	AlertMessage        *string   `json:"alert_message,omitempty"`
	NotificationChannels []string `json:"notification_channels,omitempty"`
	NotificationStatus  string    `json:"notification_status"`
	SentAt              *int64    `json:"sent_at,omitempty"`
	ErrorMessage        *string   `json:"error_message,omitempty"`
	CreatedAt           int64     `json:"created_at"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount    int64  `json:"total_count"`
	HasMore       bool   `json:"has_more"`
}

// BudgetCheckResult 预算检查结果
type BudgetCheckResult struct {
	TenantID       string  `json:"tenant_id"`
	BudgetAmount   float64 `json:"budget_amount"`
	UsedAmount     float64 `json:"used_amount"`
	UsagePercent   float64 `json:"usage_percent"`
	AlertTriggered bool    `json:"alert_triggered"`
	AlertLevel     *string `json:"alert_level,omitempty"`
	Message        string  `json:"message"`
	CheckedAt      int64   `json:"checked_at"`
}

// BudgetUsageDTO 预算使用情况DTO
type BudgetUsageDTO struct {
	// 预算信息
	BudgetAmount    float64  `json:"budget_amount"`
	UsedAmount      float64  `json:"used_amount"`
	RemainingAmount float64  `json:"remaining_amount"`
	UsagePercent    float64  `json:"usage_percent"`

	// 当前周期统计
	PeriodStart    int64  `json:"period_start"`
	PeriodEnd      int64  `json:"period_end"`
	TotalTokens    int64  `json:"total_tokens"`
	TotalRequests  int64  `json:"total_requests"` // 企业级计费系统：使用int64支持大规模请求
	AverageCost    float64 `json:"average_cost"`

	// 告警统计
	AlertCount     int     `json:"alert_count"`
	LastAlertAt    *int64  `json:"last_alert_at,omitempty"`
	LastAlertLevel *string `json:"last_alert_level,omitempty"`

	// 预测
	EstimatedDailyUsage float64 `json:"estimated_daily_usage"`
	DaysRemaining       int     `json:"days_remaining"`
	WillExceedBudget    bool    `json:"will_exceed_budget"`

	UpdatedAt int64 `json:"updated_at"`
}

// ============================================================
// 核心方法
// ============================================================

// GetBudget 获取预算配置
func (s *BudgetManagementService) GetBudget(
	ctx context.Context,
	tenantID string,
) (*BudgetSettingsDTO, error) {
	// 1. 获取预算配置
	budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("budget not found for tenant %s", tenantID)
		}
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get budget: %v", err)
		return nil, err
	}

	if budget == nil {
		return nil, fmt.Errorf("budget not found for tenant %s", tenantID)
	}

	// 2. 转换为DTO
	return s.toBudgetDTO(budget), nil
}

// CreateBudget 创建预算配置
func (s *BudgetManagementService) CreateBudget(
	ctx context.Context,
	req *CreateBudgetRequest,
) (*BudgetSettingsDTO, error) {
	// 1. 验证请求
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// 2. 检查是否已存在
	existing, err := s.budgetRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to check existing budget: %v", err)
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("budget already exists for tenant %s", req.TenantID)
	}

	// 3. 转换JSON字段
	channelsJSON, _ := json.Marshal(req.NotificationChannels)
	recipientsJSON, _ := json.Marshal(req.NotificationRecipients)

	// 4. 创建预算实体
	budget := &entity.BudgetSettings{
		TenantID:              req.TenantID,
		BudgetType:            req.BudgetType,
		BudgetAmount:          req.BudgetAmount,
		Currency:              req.Currency,
		AlertThreshold1:       req.AlertThreshold1,
		AlertThreshold2:       req.AlertThreshold2,
		HardCapEnabled:        req.HardCapEnabled,
		HardCapAmount:         req.HardCapAmount,
		AutoDowngradeEnabled:  req.AutoDowngradeEnabled,
		DowngradeThreshold:    req.DowngradeThreshold,
		OriginalModel:         req.OriginalModel,
		FallbackModel:         req.FallbackModel,
		NotificationChannels:  string(channelsJSON),
		NotificationRecipients: string(recipientsJSON),
	}

	// 5. 保存到数据库
	if err := s.budgetRepo.Create(ctx, budget); err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to create budget: %v", err)
		return nil, err
	}

	logs.CtxInfof(ctx, "[BudgetManagement] Budget created for tenant %s", req.TenantID)

	return s.toBudgetDTO(budget), nil
}

// UpdateBudget 更新预算配置
func (s *BudgetManagementService) UpdateBudget(
	ctx context.Context,
	tenantID string,
	req *UpdateBudgetRequest,
) (*BudgetSettingsDTO, error) {
	// 1. 获取现有配置
	budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("budget not found for tenant %s", tenantID)
		}
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get budget: %v", err)
		return nil, err
	}

	if budget == nil {
		return nil, fmt.Errorf("budget not found for tenant %s", tenantID)
	}

	// 2. 更新字段（仅更新非零值）
	s.updateBudgetFields(budget, req)

	// 3. 验证更新后的配置
	if err := s.validateBudgetSettings(budget); err != nil {
		return nil, err
	}

	// 4. 保存更新
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to update budget: %v", err)
		return nil, err
	}

	logs.CtxInfof(ctx, "[BudgetManagement] Budget updated for tenant %s", tenantID)

	return s.toBudgetDTO(budget), nil
}

// DeleteBudget 删除预算配置
func (s *BudgetManagementService) DeleteBudget(
	ctx context.Context,
	tenantID string,
) error {
	// 1. 检查是否存在
	budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("budget not found for tenant %s", tenantID)
		}
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get budget: %v", err)
		return err
	}

	if budget == nil {
		return fmt.Errorf("budget not found for tenant %s", tenantID)
	}

	// 2. 删除配置
	if err := s.budgetRepo.Delete(ctx, tenantID); err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to delete budget: %v", err)
		return err
	}

	logs.CtxInfof(ctx, "[BudgetManagement] Budget deleted for tenant %s", tenantID)

	return nil
}

// GetAlerts 获取告警历史
func (s *BudgetManagementService) GetAlerts(
	ctx context.Context,
	tenantID string,
	filter *AlertFilter,
) ([]*BudgetAlertDTO, *PaginationInfo, error) {
	// 1. 构建查询过滤器
	repoFilter := &repository.BudgetAlertFilter{
		TenantID:           tenantID,
		AlertType:          filter.AlertType,
		AlertLevel:         filter.AlertLevel,
		NotificationStatus: filter.NotificationStatus,
		StartDate:          filter.StartDate,
		EndDate:            filter.EndDate,
		PageToken:          filter.PageToken,
		PageSize:           filter.PageSize,
	}

	if repoFilter.PageSize == 0 {
		repoFilter.PageSize = 20
	}
	if repoFilter.PageSize > 100 {
		repoFilter.PageSize = 100
	}

	// 2. 查询告警历史
	alerts, total, err := s.alertRepo.List(ctx, repoFilter)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to list alerts: %v", err)
		return nil, nil, err
	}

	// 3. 转换为DTO
	alertDTOs := make([]*BudgetAlertDTO, 0, len(alerts))
	for _, alert := range alerts {
		alertDTOs = append(alertDTOs, s.toAlertDTO(alert))
	}

	// 4. 构建分页信息
	pagination := &PaginationInfo{
		TotalCount: total,
		HasMore:    int64(len(alertDTOs)) < total,
	}

	return alertDTOs, pagination, nil
}

// CheckBudget 手动触发预算检查
func (s *BudgetManagementService) CheckBudget(
	ctx context.Context,
	tenantID string,
) (*BudgetCheckResult, error) {
	// 1. 获取预算配置
	budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &BudgetCheckResult{
				TenantID: tenantID,
				Message:  "Budget not configured",
				CheckedAt: time.Now().Unix(),
			}, nil
		}
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get budget: %v", err)
		return nil, err
	}

	if budget == nil {
		return &BudgetCheckResult{
			TenantID: tenantID,
			Message:  "Budget not configured",
			CheckedAt: time.Now().Unix(),
		}, nil
	}

	// 2. 计算当前周期使用量
	now := time.Now()
	startDate, err := s.calculatePeriodStart(now, budget.BudgetType)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to calculate period start: %v", err)
		return nil, err
	}

	usedAmount, err := s.summaryRepo.GetTotalCost(ctx, tenantID, startDate, now)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get total cost: %v", err)
		return nil, err
	}

	// 3. 计算使用率
	usagePercent := 0.0
	if budget.BudgetAmount > 0 {
		usagePercent = usedAmount / budget.BudgetAmount * 100
	}

	// 4. 检查告警状态
	result := &BudgetCheckResult{
		TenantID:     tenantID,
		BudgetAmount: budget.BudgetAmount,
		UsedAmount:   usedAmount,
		UsagePercent: usagePercent,
		Message:      "Budget usage normal",
		CheckedAt:    now.Unix(),
	}

	if usagePercent >= float64(budget.AlertThreshold2) {
		result.AlertTriggered = true
		level := entity.AlertLevelCritical
		result.AlertLevel = &level
		result.Message = fmt.Sprintf("Critical: Usage at %.1f%% (≥%d%%)", usagePercent, budget.AlertThreshold2)
	} else if usagePercent >= float64(budget.AlertThreshold1) {
		result.AlertTriggered = true
		level := entity.AlertLevelWarning
		result.AlertLevel = &level
		result.Message = fmt.Sprintf("Warning: Usage at %.1f%% (≥%d%%)", usagePercent, budget.AlertThreshold1)
	}

	// 5. 检查硬性上限
	hardCapAmount := budget.BudgetAmount
	if budget.HardCapAmount != nil {
		hardCapAmount = *budget.HardCapAmount
	}

	if budget.HardCapEnabled && usedAmount >= hardCapAmount {
		result.AlertTriggered = true
		level := entity.AlertLevelEmergency
		result.AlertLevel = &level
		result.Message = "EMERGENCY: Hard cap exceeded!"
	}

	logs.CtxInfof(ctx, "[BudgetManagement] Budget check completed for tenant %s: %.1f%%", tenantID, usagePercent)

	return result, nil
}

// GetBudgetUsage 获取预算使用情况
func (s *BudgetManagementService) GetBudgetUsage(
	ctx context.Context,
	tenantID string,
) (*BudgetUsageDTO, error) {
	// 1. 获取预算配置
	budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("budget not found for tenant %s", tenantID)
		}
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get budget: %v", err)
		return nil, err
	}

	if budget == nil {
		return nil, fmt.Errorf("budget not found for tenant %s", tenantID)
	}

	// 2. 计算当前周期
	now := time.Now()
	startDate, err := s.calculatePeriodStart(now, budget.BudgetType)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to calculate period start: %v", err)
		return nil, err
	}

	endDate, err := s.calculatePeriodEnd(now, budget.BudgetType)
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to calculate period end: %v", err)
		return nil, err
	}

	// 3. 获取使用汇总数据
	summaries, err := s.summaryRepo.GetByDateRange(ctx, tenantID, nil, startDate.Format("2006-01-02"), now.Format("2006-01-02"))
	if err != nil {
		logs.CtxErrorf(ctx, "[BudgetManagement] Failed to get usage summaries: %v", err)
		return nil, err
	}

	// 4. 计算统计数据（企业级计费系统：统一使用int64类型）
	totalTokens := int64(0)
	totalRequests := int64(0) // 修复：使用int64类型
	usedAmount := 0.0

	for _, summary := range summaries {
		totalTokens += summary.TotalTokens
		totalRequests += summary.TotalRequests // 现在是int64类型，无需转换
		usedAmount += summary.TotalCost
	}

	// 5. 计算平均值
	averageCost := 0.0
	if totalTokens > 0 {
		averageCost = usedAmount / float64(totalTokens) * 1000 // 每1K tokens平均成本
	}

	// 6. 计算使用率
	usagePercent := 0.0
	if budget.BudgetAmount > 0 {
		usagePercent = usedAmount / budget.BudgetAmount * 100
	}

	// 7. 获取告警统计
	alerts, err := s.alertRepo.GetByTenant(ctx, tenantID, 100)
	if err != nil {
		logs.CtxWarnf(ctx, "[BudgetManagement] Failed to get alerts: %v", err)
	}

	alertCount := len(alerts)
	var lastAlertAt *int64
	var lastAlertLevel *string

	if len(alerts) > 0 {
		ts := alerts[0].CreatedAt.Unix()
		lastAlertAt = &ts
		lastAlertLevel = &alerts[0].AlertLevel
	}

	// 8. 计算预测数据
	daysRemaining := int(endDate.Sub(now).Hours() / 24)
	estimatedDailyUsage := 0.0
	willExceedBudget := false

	if daysRemaining > 0 && len(summaries) > 0 {
		// 基于已过天数计算平均每日使用量
		daysPassed := int(now.Sub(startDate).Hours()/24) + 1
		estimatedDailyUsage = usedAmount / float64(daysPassed)

		// 预测是否会超预算
		predictedUsage := usedAmount + (estimatedDailyUsage * float64(daysRemaining))
		willExceedBudget = predictedUsage > budget.BudgetAmount
	}

	// 9. 构建返回结果
	return &BudgetUsageDTO{
		BudgetAmount:         budget.BudgetAmount,
		UsedAmount:           usedAmount,
		RemainingAmount:      budget.BudgetAmount - usedAmount,
		UsagePercent:         usagePercent,
		PeriodStart:          startDate.Unix(),
		PeriodEnd:            endDate.Unix(),
		TotalTokens:          totalTokens,
		TotalRequests:        totalRequests,
		AverageCost:          averageCost,
		AlertCount:           alertCount,
		LastAlertAt:          lastAlertAt,
		LastAlertLevel:       lastAlertLevel,
		EstimatedDailyUsage:  estimatedDailyUsage,
		DaysRemaining:        daysRemaining,
		WillExceedBudget:     willExceedBudget,
		UpdatedAt:            now.Unix(),
	}, nil
}

// ============================================================
// 辅助方法
// ============================================================

// toBudgetDTO 转换为BudgetDTO
func (s *BudgetManagementService) toBudgetDTO(budget *entity.BudgetSettings) *BudgetSettingsDTO {
	dto := &BudgetSettingsDTO{
		TenantID:              budget.TenantID,
		BudgetType:            budget.BudgetType,
		BudgetAmount:          budget.BudgetAmount,
		Currency:              budget.Currency,
		AlertThreshold1:       budget.AlertThreshold1,
		AlertThreshold2:       budget.AlertThreshold2,
		HardCapEnabled:        budget.HardCapEnabled,
		HardCapAmount:         budget.HardCapAmount,
		AutoDowngradeEnabled:  budget.AutoDowngradeEnabled,
		DowngradeThreshold:    budget.DowngradeThreshold,
		OriginalModel:         budget.OriginalModel,
		FallbackModel:         budget.FallbackModel,
		CreatedAt:             budget.CreatedAt.Unix(),
		UpdatedAt:             budget.UpdatedAt.Unix(),
	}

	// 解析JSON字段
	if budget.NotificationChannels != "" {
		var channels []string
		json.Unmarshal([]byte(budget.NotificationChannels), &channels)
		dto.NotificationChannels = channels
	}

	if budget.NotificationRecipients != "" {
		var recipients []string
		json.Unmarshal([]byte(budget.NotificationRecipients), &recipients)
		dto.NotificationRecipients = recipients
	}

	return dto
}

// toAlertDTO 转换为AlertDTO
func (s *BudgetManagementService) toAlertDTO(alert *entity.BudgetAlert) *BudgetAlertDTO {
	dto := &BudgetAlertDTO{
		ID:              alert.ID,
		TenantID:        alert.TenantID,
		AlertType:       alert.AlertType,
		BudgetAmount:    alert.BudgetAmount,
		UsedAmount:      alert.UsedAmount,
		UsagePercent:    alert.UsagePercent,
		AlertLevel:      alert.AlertLevel,
		AlertMessage:    alert.AlertMessage,
		NotificationStatus: alert.NotificationStatus,
		ErrorMessage:    alert.ErrorMessage,
		CreatedAt:       alert.CreatedAt.Unix(),
	}

	// 解析JSON字段
	if alert.NotificationChannels != "" {
		var channels []string
		json.Unmarshal([]byte(alert.NotificationChannels), &channels)
		dto.NotificationChannels = channels
	}

	// 转换时间戳
	if alert.SentAt != nil {
		ts := alert.SentAt.Unix()
		dto.SentAt = &ts
	}

	return dto
}

// validateCreateRequest 验证创建请求
func (s *BudgetManagementService) validateCreateRequest(req *CreateBudgetRequest) error {
	// 1. 基本验证
	if req.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	// 2. 阈值验证
	if req.AlertThreshold1 >= req.AlertThreshold2 {
		return fmt.Errorf("alert_threshold_1 must be less than alert_threshold_2")
	}

	// 3. 降级配置验证
	if req.AutoDowngradeEnabled {
		if req.OriginalModel == nil || req.FallbackModel == nil {
			return fmt.Errorf("original_model and fallback_model are required when auto_downgrade_enabled is true")
		}
		if *req.OriginalModel == *req.FallbackModel {
			return fmt.Errorf("original_model and fallback_model must be different")
		}
	}

	// 4. 硬性上限验证
	if req.HardCapEnabled && req.HardCapAmount != nil {
		if *req.HardCapAmount < req.BudgetAmount {
			return fmt.Errorf("hard_cap_amount cannot be less than budget_amount")
		}
	}

	// 5. 通知配置验证
	if len(req.NotificationChannels) == 0 {
		return fmt.Errorf("at least one notification channel is required")
	}

	if len(req.NotificationRecipients) == 0 {
		return fmt.Errorf("at least one notification recipient is required")
	}

	return nil
}

// updateBudgetFields 更新预算字段
func (s *BudgetManagementService) updateBudgetFields(budget *entity.BudgetSettings, req *UpdateBudgetRequest) {
	if req.BudgetAmount != nil {
		budget.BudgetAmount = *req.BudgetAmount
	}
	if req.AlertThreshold1 != nil {
		budget.AlertThreshold1 = *req.AlertThreshold1
	}
	if req.AlertThreshold2 != nil {
		budget.AlertThreshold2 = *req.AlertThreshold2
	}
	if req.HardCapEnabled != nil {
		budget.HardCapEnabled = *req.HardCapEnabled
	}
	if req.HardCapAmount != nil {
		budget.HardCapAmount = req.HardCapAmount
	}
	if req.AutoDowngradeEnabled != nil {
		budget.AutoDowngradeEnabled = *req.AutoDowngradeEnabled
	}
	if req.DowngradeThreshold != nil {
		budget.DowngradeThreshold = *req.DowngradeThreshold
	}
	if req.OriginalModel != nil {
		budget.OriginalModel = req.OriginalModel
	}
	if req.FallbackModel != nil {
		budget.FallbackModel = req.FallbackModel
	}
	if req.NotificationChannels != nil {
		channelsJSON, _ := json.Marshal(req.NotificationChannels)
		budget.NotificationChannels = string(channelsJSON)
	}
	if req.NotificationRecipients != nil {
		recipientsJSON, _ := json.Marshal(req.NotificationRecipients)
		budget.NotificationRecipients = string(recipientsJSON)
	}
}

// validateBudgetSettings 验证预算配置
func (s *BudgetManagementService) validateBudgetSettings(budget *entity.BudgetSettings) error {
	// 1. 阈值验证
	if budget.AlertThreshold1 >= budget.AlertThreshold2 {
		return fmt.Errorf("alert_threshold_1 must be less than alert_threshold_2")
	}

	// 2. 降级配置验证
	if budget.AutoDowngradeEnabled {
		if budget.OriginalModel == nil || budget.FallbackModel == nil {
			return fmt.Errorf("original_model and fallback_model are required when auto_downgrade_enabled is true")
		}
		if *budget.OriginalModel == *budget.FallbackModel {
			return fmt.Errorf("original_model and fallback_model must be different")
		}
	}

	// 3. 硬性上限验证
	if budget.HardCapEnabled && budget.HardCapAmount != nil {
		if *budget.HardCapAmount < budget.BudgetAmount {
			return fmt.Errorf("hard_cap_amount cannot be less than budget_amount")
		}
	}

	// 4. 通知配置验证
	if budget.NotificationChannels == "" {
		return fmt.Errorf("at least one notification channel is required")
	}

	if budget.NotificationRecipients == "" {
		return fmt.Errorf("at least one notification recipient is required")
	}

	return nil
}

// calculatePeriodStart 计算周期开始时间
func (s *BudgetManagementService) calculatePeriodStart(now time.Time, budgetType string) (time.Time, error) {
	switch budgetType {
	case entity.BudgetTypeMonthly:
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), nil
	case entity.BudgetTypeQuarterly:
		quarter := (now.Month() - 1) / 3
		return time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location()), nil
	case entity.BudgetTypeYearly:
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()), nil
	default:
		return time.Time{}, fmt.Errorf("invalid budget type: %s", budgetType)
	}
}

// calculatePeriodEnd 计算周期结束时间
func (s *BudgetManagementService) calculatePeriodEnd(now time.Time, budgetType string) (time.Time, error) {
	switch budgetType {
	case entity.BudgetTypeMonthly:
		return now.AddDate(0, 1, 0), nil
	case entity.BudgetTypeQuarterly:
		return now.AddDate(0, 3, 0), nil
	case entity.BudgetTypeYearly:
		return now.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid budget type: %s", budgetType)
	}
}
