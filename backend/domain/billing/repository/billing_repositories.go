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
// Token计量相关Repository接口
// ============================================================

// TokenUsageLogRepository Token使用明细日志仓储接口
type TokenUsageLogRepository interface {
	// Create 创建Token使用日志记录
	Create(ctx context.Context, log *entity.TokenUsageLog) error

	// GetByID 根据ID获取日志记录
	GetByID(ctx context.Context, id uint64) (*entity.TokenUsageLog, error)

	// GetByDateRange 按日期范围查询日志
	GetByDateRange(
		ctx context.Context,
		tenantID string,
		startDate, endDate string,
	) ([]*entity.TokenUsageLog, error)

	// GetByBotAndDateRange 按Bot和日期范围查询
	GetByBotAndDateRange(
		ctx context.Context,
		tenantID, botID string,
		startDate, endDate string,
	) ([]*entity.TokenUsageLog, error)

	// List 分页查询日志
	List(ctx context.Context, filter *TokenUsageLogFilter) ([]*entity.TokenUsageLog, int64, error)

	// BatchCreate 批量创建日志记录
	BatchCreate(ctx context.Context, logs []*entity.TokenUsageLog) error
}

// TokenUsageLogFilter Token日志查询过滤器
type TokenUsageLogFilter struct {
	TenantID     string
	BotID        *string
	ModelProvider *string
	ModelName    *string
	StartDate    *time.Time
	EndDate      *time.Time
	PageToken    string
	PageSize     int
}

// TokenUsageSummaryRepository Token使用汇总仓储接口
type TokenUsageSummaryRepository interface {
	// Upsert 创建或更新汇总记录
	Upsert(ctx context.Context, summary *entity.TokenUsageSummary) error

	// GetByTenantAndDate 获取租户在指定日期的汇总
	GetByTenantAndDate(
		ctx context.Context,
		tenantID string,
		summaryDate string,
		summaryHour *uint8,
	) (*entity.TokenUsageSummary, error)

	// GetByDateRange 按日期范围查询汇总数据
	GetByDateRange(
		ctx context.Context,
		tenantID string,
		botID *string,
		startDate, endDate string,
	) ([]*entity.TokenUsageSummary, error)

	// GetTotalCost 获取指定时间范围内的总成本
	GetTotalCost(
		ctx context.Context,
		tenantID string,
		startDate, endDate time.Time,
	) (float64, error)

	// List 分页查询汇总数据
	List(ctx context.Context, filter *TokenUsageSummaryFilter) ([]*entity.TokenUsageSummary, int64, error)
}

// TokenUsageSummaryFilter Token汇总查询过滤器
type TokenUsageSummaryFilter struct {
	TenantID     string
	BotID        *string
	StartDate    string
	EndDate      string
	Granularity  string // daily, hourly
	PageToken    string
	PageSize     int
}

// ============================================================
// 预算相关Repository接口
// ============================================================

// BudgetSettingsRepository 预算配置仓储接口
type BudgetSettingsRepository interface {
	// Create 创建预算配置
	Create(ctx context.Context, budget *entity.BudgetSettings) error

	// GetByTenantID 获取租户的预算配置
	GetByTenantID(ctx context.Context, tenantID string) (*entity.BudgetSettings, error)

	// Update 更新预算配置
	Update(ctx context.Context, budget *entity.BudgetSettings) error

	// Delete 删除预算配置
	Delete(ctx context.Context, tenantID string) error

	// List 分页查询预算配置
	List(ctx context.Context, filter *BudgetSettingsFilter) ([]*entity.BudgetSettings, int64, error)
}

// BudgetSettingsFilter 预算配置查询过滤器
type BudgetSettingsFilter struct {
	BudgetType   *string
	HardCapEnabled *bool
	PageToken    string
	PageSize     int
}

// BudgetAlertRepository 预算告警仓储接口
type BudgetAlertRepository interface {
	// Create 创建告警记录
	Create(ctx context.Context, alert *entity.BudgetAlert) error

	// GetByID 根据ID获取告警
	GetByID(ctx context.Context, id uint64) (*entity.BudgetAlert, error)

	// GetByTenant 获取租户的告警历史
	GetByTenant(ctx context.Context, tenantID string, limit int) ([]*entity.BudgetAlert, error)

	// GetByTenantAndType 获取租户特定类型的告警
	GetByTenantAndType(
		ctx context.Context,
		tenantID string,
		alertType string,
		limit int,
	) ([]*entity.BudgetAlert, error)

	// UpdateNotificationStatus 更新通知状态
	UpdateNotificationStatus(
		ctx context.Context,
		id uint64,
		status string,
		sentAt *time.Time,
		errorMessage *string,
	) error

	// CheckAlertSentToday 检查今天是否已发送过特定类型的告警
	CheckAlertSentToday(
		ctx context.Context,
		tenantID string,
		alertType string,
		today string,
	) (bool, error)

	// List 分页查询告警
	List(ctx context.Context, filter *BudgetAlertFilter) ([]*entity.BudgetAlert, int64, error)
}

// BudgetAlertFilter 预算告警查询过滤器
type BudgetAlertFilter struct {
	TenantID         string
	AlertType        *string
	AlertLevel       *string
	NotificationStatus *string
	StartDate        *time.Time
	EndDate          *time.Time
	PageToken        string
	PageSize         int
}

// ============================================================
// 成本优化相关Repository接口
// ============================================================

// CostOptimizationSuggestionRepository 成本优化建议仓储接口
type CostOptimizationSuggestionRepository interface {
	// Create 创建优化建议
	Create(ctx context.Context, suggestion *entity.CostOptimizationSuggestion) error

	// GetByID 根据ID获取建议
	GetByID(ctx context.Context, id uint64) (*entity.CostOptimizationSuggestion, error)

	// GetByTenant 获取租户的优化建议
	GetByTenant(
		ctx context.Context,
		tenantID string,
		status *string,
		limit int,
	) ([]*entity.CostOptimizationSuggestion, error)

	// UpdateStatus 更新建议状态
	UpdateStatus(
		ctx context.Context,
		id uint64,
		status string,
		appliedAt *time.Time,
		appliedBy *uint64,
	) error

	// Delete 删除建议
	Delete(ctx context.Context, id uint64) error

	// List 分页查询建议
	List(ctx context.Context, filter *OptimizationSuggestionFilter) ([]*entity.CostOptimizationSuggestion, int64, error)

	// GetEstimatedMonthlySaving 获取租户预计月节省金额
	GetEstimatedMonthlySaving(
		ctx context.Context,
		tenantID string,
		status *string,
	) (float64, error)
}

// OptimizationSuggestionFilter 优化建议查询过滤器
type OptimizationSuggestionFilter struct {
	TenantID        string
	SuggestionType  *string
	Status          *string
	MinSaving       *float64
	StartDate       *time.Time
	EndDate         *time.Time
	PageToken       string
	PageSize        int
}
