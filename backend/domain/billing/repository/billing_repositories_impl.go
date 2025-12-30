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
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// TokenUsageLogRepository实现
// ============================================================

// tokenUsageLogRepository Token使用明细日志仓储实现
type tokenUsageLogRepository struct {
	db *gorm.DB
}

// NewTokenUsageLogRepository 创建Token使用明细日志仓储实例
func NewTokenUsageLogRepository(db *gorm.DB) TokenUsageLogRepository {
	return &tokenUsageLogRepository{db: db}
}

// Create 创建Token使用日志记录
func (r *tokenUsageLogRepository) Create(ctx context.Context, log *entity.TokenUsageLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByID 根据ID获取日志记录
func (r *tokenUsageLogRepository) GetByID(ctx context.Context, id uint64) (*entity.TokenUsageLog, error) {
	var log entity.TokenUsageLog
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&log).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByDateRange 按日期范围查询日志
func (r *tokenUsageLogRepository) GetByDateRange(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
) ([]*entity.TokenUsageLog, error) {
	var logs []*entity.TokenUsageLog
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

// GetByBotAndDateRange 按Bot和日期范围查询
func (r *tokenUsageLogRepository) GetByBotAndDateRange(
	ctx context.Context,
	tenantID, botID string,
	startDate, endDate string,
) ([]*entity.TokenUsageLog, error) {
	var logs []*entity.TokenUsageLog
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("bot_id = ?", botID).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

// List 分页查询日志
func (r *tokenUsageLogRepository) List(
	ctx context.Context,
	filter *TokenUsageLogFilter,
) ([]*entity.TokenUsageLog, int64, error) {
	var logs []*entity.TokenUsageLog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TokenUsageLog{})

	// 租户ID过滤
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	// BotID过滤
	if filter.BotID != nil {
		query = query.Where("bot_id = ?", *filter.BotID)
	}

	// 模型提供商过滤
	if filter.ModelProvider != nil {
		query = query.Where("model_provider = ?", *filter.ModelProvider)
	}

	// 模型名称过滤
	if filter.ModelName != nil {
		query = query.Where("model_name = ?", *filter.ModelName)
	}

	// 日期范围过滤
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parsePageToken(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 查询
	if err := query.Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// BatchCreate 批量创建日志记录
func (r *tokenUsageLogRepository) BatchCreate(ctx context.Context, logs []*entity.TokenUsageLog) error {
	if len(logs) == 0 {
		return nil
	}

	// GORM批量插入
	return r.db.WithContext(ctx).CreateInBatches(logs, 100).Error
}

// ============================================================
// TokenUsageSummaryRepository实现
// ============================================================

// tokenUsageSummaryRepository Token使用汇总仓储实现
type tokenUsageSummaryRepository struct {
	db *gorm.DB
}

// NewTokenUsageSummaryRepository 创建Token使用汇总仓储实例
func NewTokenUsageSummaryRepository(db *gorm.DB) TokenUsageSummaryRepository {
	return &tokenUsageSummaryRepository{db: db}
}

// Upsert 创建或更新汇总记录
func (r *tokenUsageSummaryRepository) Upsert(ctx context.Context, summary *entity.TokenUsageSummary) error {
	// 使用MySQL的ON DUPLICATE KEY UPDATE实现upsert
	sql := `
		INSERT INTO token_usage_summary (
			tenant_id, bot_id, summary_date, summary_hour,
			total_input_tokens, total_output_tokens, total_tokens,
			total_cost, total_requests, cached_requests, avg_response_time,
			model_distribution, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			total_input_tokens = total_input_tokens + VALUES(total_input_tokens),
			total_output_tokens = total_output_tokens + VALUES(total_output_tokens),
			total_tokens = total_tokens + VALUES(total_tokens),
			total_cost = total_cost + VALUES(total_cost),
			total_requests = total_requests + VALUES(total_requests),
			cached_requests = cached_requests + VALUES(cached_requests),
			avg_response_time = (avg_response_time + VALUES(avg_response_time)) / 2,
			model_distribution = IFNULL(model_distribution, '{}'),
			updated_at = NOW()
	`

	return r.db.WithContext(ctx).Exec(sql,
		summary.TenantID, summary.BotID, summary.SummaryDate, summary.SummaryHour,
		summary.TotalInputTokens, summary.TotalOutputTokens, summary.TotalTokens,
		summary.TotalCost, summary.TotalRequests, summary.CachedRequests, summary.AvgResponseTime,
		summary.ModelDistribution,
	).Error
}

// GetByTenantAndDate 获取租户在指定日期的汇总
func (r *tokenUsageSummaryRepository) GetByTenantAndDate(
	ctx context.Context,
	tenantID string,
	summaryDate string,
	summaryHour *uint8,
) (*entity.TokenUsageSummary, error) {
	var summary entity.TokenUsageSummary
	query := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("summary_date = ?", summaryDate)

	if summaryHour != nil {
		query = query.Where("summary_hour = ?", *summaryHour)
	} else {
		query = query.Where("summary_hour IS NULL")
	}

	err := query.First(&summary).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetByDateRange 按日期范围查询汇总数据
func (r *tokenUsageSummaryRepository) GetByDateRange(
	ctx context.Context,
	tenantID string,
	botID *string,
	startDate, endDate string,
) ([]*entity.TokenUsageSummary, error) {
	var summaries []*entity.TokenUsageSummary
	query := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("summary_date >= ? AND summary_date <= ?", startDate, endDate)

	if botID != nil {
		query = query.Where("bot_id = ?", *botID)
	}

	err := query.Order("summary_date DESC, summary_hour DESC").
		Find(&summaries).Error

	return summaries, err
}

// GetTotalCost 获取指定时间范围内的总成本
func (r *tokenUsageSummaryRepository) GetTotalCost(
	ctx context.Context,
	tenantID string,
	startDate, endDate time.Time,
) (float64, error) {
	var totalCost float64
	err := r.db.WithContext(ctx).
		Model(&entity.TokenUsageSummary{}).
		Where("tenant_id = ?", tenantID).
		Where("summary_date >= ? AND summary_date <= ?",
			startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02")).
		Select("COALESCE(SUM(total_cost), 0)").
		Scan(&totalCost).Error

	return totalCost, err
}

// List 分页查询汇总数据
func (r *tokenUsageSummaryRepository) List(
	ctx context.Context,
	filter *TokenUsageSummaryFilter,
) ([]*entity.TokenUsageSummary, int64, error) {
	var summaries []*entity.TokenUsageSummary
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TokenUsageSummary{})

	// 租户ID过滤
	query = query.Where("tenant_id = ?", filter.TenantID)

	// BotID过滤
	if filter.BotID != nil {
		query = query.Where("bot_id = ?", *filter.BotID)
	}

	// 日期范围过滤
	if filter.StartDate != "" {
		query = query.Where("summary_date >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("summary_date <= ?", filter.EndDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parsePageToken(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 查询并排序
	if err := query.Order("summary_date DESC, summary_hour DESC").
		Find(&summaries).Error; err != nil {
		return nil, 0, err
	}

	return summaries, total, nil
}

// ============================================================
// BudgetSettingsRepository实现
// ============================================================

// budgetSettingsRepository 预算配置仓储实现
type budgetSettingsRepository struct {
	db *gorm.DB
}

// NewBudgetSettingsRepository 创建预算配置仓储实例
func NewBudgetSettingsRepository(db *gorm.DB) BudgetSettingsRepository {
	return &budgetSettingsRepository{db: db}
}

// Create 创建预算配置
func (r *budgetSettingsRepository) Create(ctx context.Context, budget *entity.BudgetSettings) error {
	return r.db.WithContext(ctx).Create(budget).Error
}

// GetByTenantID 获取租户的预算配置
func (r *budgetSettingsRepository) GetByTenantID(ctx context.Context, tenantID string) (*entity.BudgetSettings, error) {
	var budget entity.BudgetSettings
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		First(&budget).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

// Update 更新预算配置
func (r *budgetSettingsRepository) Update(ctx context.Context, budget *entity.BudgetSettings) error {
	return r.db.WithContext(ctx).
		Model(&entity.BudgetSettings{}).
		Where("tenant_id = ?", budget.TenantID).
		Updates(budget).Error
}

// Delete 删除预算配置
func (r *budgetSettingsRepository) Delete(ctx context.Context, tenantID string) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&entity.BudgetSettings{}).Error
}

// List 分页查询预算配置
func (r *budgetSettingsRepository) List(
	ctx context.Context,
	filter *BudgetSettingsFilter,
) ([]*entity.BudgetSettings, int64, error) {
	var budgets []*entity.BudgetSettings
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BudgetSettings{})

	// 预算类型过滤
	if filter.BudgetType != nil {
		query = query.Where("budget_type = ?", *filter.BudgetType)
	}

	// 硬性上限过滤
	if filter.HardCapEnabled != nil {
		query = query.Where("hard_cap_enabled = ?", *filter.HardCapEnabled)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parsePageToken(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 查询
	if err := query.Order("created_at DESC").Find(&budgets).Error; err != nil {
		return nil, 0, err
	}

	return budgets, total, nil
}

// ============================================================
// BudgetAlertRepository实现
// ============================================================

// budgetAlertRepository 预算告警仓储实现
type budgetAlertRepository struct {
	db *gorm.DB
}

// NewBudgetAlertRepository 创建预算告警仓储实例
func NewBudgetAlertRepository(db *gorm.DB) BudgetAlertRepository {
	return &budgetAlertRepository{db: db}
}

// Create 创建告警记录
func (r *budgetAlertRepository) Create(ctx context.Context, alert *entity.BudgetAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// GetByID 根据ID获取告警
func (r *budgetAlertRepository) GetByID(ctx context.Context, id uint64) (*entity.BudgetAlert, error) {
	var alert entity.BudgetAlert
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&alert).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// GetByTenant 获取租户的告警历史
func (r *budgetAlertRepository) GetByTenant(
	ctx context.Context,
	tenantID string,
	limit int,
) ([]*entity.BudgetAlert, error) {
	var alerts []*entity.BudgetAlert
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Find(&alerts).Error

	return alerts, err
}

// GetByTenantAndType 获取租户特定类型的告警
func (r *budgetAlertRepository) GetByTenantAndType(
	ctx context.Context,
	tenantID string,
	alertType string,
	limit int,
) ([]*entity.BudgetAlert, error) {
	var alerts []*entity.BudgetAlert
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("alert_type = ?", alertType).
		Order("created_at DESC").
		Limit(limit).
		Find(&alerts).Error

	return alerts, err
}

// UpdateNotificationStatus 更新通知状态
func (r *budgetAlertRepository) UpdateNotificationStatus(
	ctx context.Context,
	id uint64,
	status string,
	sentAt *time.Time,
	errorMessage *string,
) error {
	updates := map[string]interface{}{
		"notification_status": status,
	}
	if sentAt != nil {
		updates["sent_at"] = *sentAt
	}
	if errorMessage != nil {
		updates["error_message"] = *errorMessage
	}

	return r.db.WithContext(ctx).
		Model(&entity.BudgetAlert{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// CheckAlertSentToday 检查今天是否已发送过特定类型的告警
func (r *budgetAlertRepository) CheckAlertSentToday(
	ctx context.Context,
	tenantID string,
	alertType string,
	today string,
) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.BudgetAlert{}).
		Where("tenant_id = ?", tenantID).
		Where("alert_type = ?", alertType).
		Where("DATE_FORMAT(FROM_UNIXTIME(created_at/1000), '%Y-%m-%d') = ?", today).
		Where("notification_status IN ?", []string{"sent", "pending"}).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// List 分页查询告警
func (r *budgetAlertRepository) List(
	ctx context.Context,
	filter *BudgetAlertFilter,
) ([]*entity.BudgetAlert, int64, error) {
	var alerts []*entity.BudgetAlert
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.BudgetAlert{})

	// 租户ID过滤
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}

	// 告警类型过滤
	if filter.AlertType != nil {
		query = query.Where("alert_type = ?", *filter.AlertType)
	}

	// 告警级别过滤
	if filter.AlertLevel != nil {
		query = query.Where("alert_level = ?", *filter.AlertLevel)
	}

	// 通知状态过滤
	if filter.NotificationStatus != nil {
		query = query.Where("notification_status = ?", *filter.NotificationStatus)
	}

	// 日期范围过滤
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parsePageToken(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 查询
	if err := query.Order("created_at DESC").Find(&alerts).Error; err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// ============================================================
// CostOptimizationSuggestionRepository实现
// ============================================================

// costOptimizationSuggestionRepository 成本优化建议仓储实现
type costOptimizationSuggestionRepository struct {
	db *gorm.DB
}

// NewCostOptimizationSuggestionRepository 创建成本优化建议仓储实例
func NewCostOptimizationSuggestionRepository(db *gorm.DB) CostOptimizationSuggestionRepository {
	return &costOptimizationSuggestionRepository{db: db}
}

// Create 创建优化建议
func (r *costOptimizationSuggestionRepository) Create(
	ctx context.Context,
	suggestion *entity.CostOptimizationSuggestion,
) error {
	return r.db.WithContext(ctx).Create(suggestion).Error
}

// GetByID 根据ID获取建议
func (r *costOptimizationSuggestionRepository) GetByID(
	ctx context.Context,
	id uint64,
) (*entity.CostOptimizationSuggestion, error) {
	var suggestion entity.CostOptimizationSuggestion
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&suggestion).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &suggestion, nil
}

// GetByTenant 获取租户的优化建议
func (r *costOptimizationSuggestionRepository) GetByTenant(
	ctx context.Context,
	tenantID string,
	status *string,
	limit int,
) ([]*entity.CostOptimizationSuggestion, error) {
	var suggestions []*entity.CostOptimizationSuggestion
	query := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Find(&suggestions).Error

	return suggestions, err
}

// UpdateStatus 更新建议状态
func (r *costOptimizationSuggestionRepository) UpdateStatus(
	ctx context.Context,
	id uint64,
	status string,
	appliedAt *time.Time,
	appliedBy *uint64,
) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if appliedAt != nil {
		updates["applied_at"] = *appliedAt
	}
	if appliedBy != nil {
		updates["applied_by"] = *appliedBy
	}

	return r.db.WithContext(ctx).
		Model(&entity.CostOptimizationSuggestion{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除建议
func (r *costOptimizationSuggestionRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.CostOptimizationSuggestion{}).Error
}

// List 分页查询建议
func (r *costOptimizationSuggestionRepository) List(
	ctx context.Context,
	filter *OptimizationSuggestionFilter,
) ([]*entity.CostOptimizationSuggestion, int64, error) {
	var suggestions []*entity.CostOptimizationSuggestion
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.CostOptimizationSuggestion{})

	// 租户ID过滤
	query = query.Where("tenant_id = ?", filter.TenantID)

	// 建议类型过滤
	if filter.SuggestionType != nil {
		query = query.Where("suggestion_type = ?", *filter.SuggestionType)
	}

	// 状态过滤
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	// 最小节省金额过滤
	if filter.MinSaving != nil {
		query = query.Where("estimated_monthly_saving >= ?", *filter.MinSaving)
	}

	// 日期范围过滤
	if filter.StartDate != nil {
		query = query.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("created_at <= ?", *filter.EndDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parsePageToken(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 查询
	if err := query.Order("estimated_monthly_saving DESC").
		Find(&suggestions).Error; err != nil {
		return nil, 0, err
	}

	return suggestions, total, nil
}

// GetEstimatedMonthlySaving 获取租户预计月节省金额
func (r *costOptimizationSuggestionRepository) GetEstimatedMonthlySaving(
	ctx context.Context,
	tenantID string,
	status *string,
) (float64, error) {
	var totalSaving float64
	query := r.db.WithContext(ctx).
		Model(&entity.CostOptimizationSuggestion{}).
		Where("tenant_id = ?", tenantID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.
		Select("COALESCE(SUM(estimated_monthly_saving), 0)").
		Scan(&totalSaving).Error

	return totalSaving, err
}

// ============================================================
// 辅助函数
// ============================================================

// parsePageToken 解析分页token（简单的偏移量解析）
func parsePageToken(token string) int {
	// 简化实现：假设token就是偏移量
	// 实际应用中应该使用更安全的编码方式
	offset := 0
	if _, err := fmt.Sscanf(token, "%d", &offset); err == nil {
		return offset
	}
	return 0
}
