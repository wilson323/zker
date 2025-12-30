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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ============================================================
// 计费系统 Prometheus 指标定义
// ============================================================

var (
	// ========== Token计量指标 ==========

	// TokenRecordsTotal Token记录总数(Counter)
	TokenRecordsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_token_records_total",
			Help: "Total number of token usage records",
		},
		[]string{"tenant_id", "model_provider", "model_name"},
	)

	// TokenRecordDuration Token记录耗时(Histogram)
	TokenRecordDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_token_record_duration_seconds",
			Help:    "Duration of token usage record operations",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 2},
		},
		[]string{"operation"}, // record, batch_record
	)

	// TokensTotal Token总数(Counter)
	TokensTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_tokens_total",
			Help: "Total number of tokens processed",
		},
		[]string{"tenant_id", "token_type"}, // input, output
	)

	// CostTotal 成本总额(Counter)
	CostTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_cost_total",
			Help: "Total cost in CNY",
		},
		[]string{"tenant_id", "model_provider"},
	)

	// CachedRequestsTotal 缓存请求数(Counter)
	CachedRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_cached_requests_total",
			Help: "Total number of cached requests",
		},
		[]string{"tenant_id", "model_provider"},
	)

	// ========== 定价引擎指标 ==========

	// PricingCalculationsTotal 定价计算次数(Counter)
	PricingCalculationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_pricing_calculations_total",
			Help: "Total number of pricing calculations",
		},
		[]string{"model_provider", "model_name", "result"}, // result: success, failure
	)

	// PricingCalculationDuration 定价计算耗时(Histogram)
	PricingCalculationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_pricing_calculation_duration_seconds",
			Help:    "Duration of pricing calculation operations",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05},
		},
		[]string{"model_provider"},
	)

	// PricingErrorsTotal 定价错误数(Counter)
	PricingErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_pricing_errors_total",
			Help: "Total number of pricing calculation errors",
		},
		[]string{"model_provider", "model_name", "error_type"},
	)

	// UnitPriceGauge 单价(Gauge)
	UnitPriceGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_unit_price",
			Help: "Current unit price per 1K tokens",
		},
		[]string{"model_provider", "model_name", "token_type"},
	)

	// ========== 预算告警指标 ==========

	// BudgetChecksTotal 预算检查次数(Counter)
	BudgetChecksTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_budget_checks_total",
			Help: "Total number of budget checks",
		},
		[]string{"tenant_id", "triggered"}, // triggered: true/false
	)

	// BudgetCheckDuration 预算检查耗时(Histogram)
	BudgetCheckDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_budget_check_duration_seconds",
			Help:    "Duration of budget check operations",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
		},
		[]string{"tenant_id"},
	)

	// BudgetAlertsTotal 预算告警次数(Counter)
	BudgetAlertsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_budget_alerts_total",
			Help: "Total number of budget alerts",
		},
		[]string{"tenant_id", "alert_type", "alert_level"}, // alert_type: threshold1, threshold2, hard_cap, downgrade
	)

	// BudgetUsageGauge 当前使用率(Gauge)
	BudgetUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_budget_usage_percent",
			Help: "Current budget usage percentage",
		},
		[]string{"tenant_id", "budget_type", "alert_level"}, // budget_type: monthly, quarterly, yearly
	)

	// BudgetAmountGauge 预算金额(Gauge)
	BudgetAmountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_budget_amount",
			Help: "Budget amount in CNY",
		},
		[]string{"tenant_id", "budget_type"},
	)

	// BudgetUsedAmountGauge 已使用金额(Gauge)
	BudgetUsedAmountGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_budget_used_amount",
			Help: "Used budget amount in CNY",
		},
		[]string{"tenant_id", "budget_type"},
	)

	// ========== 汇总统计指标 ==========

	// SummaryUpdateTotal 汇总更新次数(Counter)
	SummaryUpdateTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_summary_update_total",
			Help: "Total number of summary updates",
		},
		[]string{"tenant_id", "result"}, // result: success, failure
	)

	// SummaryUpdateDuration 汇总更新耗时(Histogram)
	SummaryUpdateDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_summary_update_duration_seconds",
			Help:    "Duration of summary update operations",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
		},
		[]string{"tenant_id"},
	)

	// DailyTokensGauge 每日Token使用量(Gauge)
	DailyTokensGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_daily_tokens",
			Help: "Daily token usage",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// DailyCostGauge 每日成本(Gauge)
	DailyCostGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_daily_cost",
			Help: "Daily cost in CNY",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// ========== 数据库操作指标 ==========

	// TokenLogDBOperationsTotal Token日志数据库操作总数(Counter)
	TokenLogDBOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_token_log_db_operations_total",
			Help: "Total number of token log database operations",
		},
		[]string{"operation", "result"}, // operation: create, batch_create, get_by_date_range; result: success, failure
	)

	// TokenLogDBOperationDuration Token日志数据库操作耗时(Histogram)
	TokenLogDBOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_token_log_db_operation_duration_seconds",
			Help:    "Duration of token log database operations",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"operation"},
	)

	// ========== 通知发送指标 ==========

	// AlertNotificationTotal 告警通知发送总数(Counter)
	AlertNotificationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "billing_alert_notification_total",
			Help: "Total number of alert notifications sent",
		},
		[]string{"tenant_id", "channel", "result"}, // channel: email, webhook, sms; result: success, failure
	)

	// AlertNotificationDuration 告警通知发送耗时(Histogram)
	AlertNotificationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "billing_alert_notification_duration_seconds",
			Help:    "Duration of alert notification operations",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"channel"},
	)

	// ========== 业务指标 ==========

	// ModelUsageGauge 模型使用量(Gauge)
	ModelUsageGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_model_usage",
			Help: "Current model usage metrics",
		},
		[]string{"tenant_id", "model_provider", "model_name"},
	)

	// ActiveTenantsGauge 活跃租户数(Gauge)
	ActiveTenantsGauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "billing_active_tenants",
			Help: "Number of active tenants with billing activity",
		},
	)

	// TotalRevenueGauge 总收入(Gauge)
	TotalRevenueGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "billing_total_revenue",
			Help: "Total revenue in CNY",
		},
		[]string{"period"}, // period: daily, monthly, yearly
	)
)

// ============================================================
// 辅助函数 - 记录指标
// ============================================================

// RecordTokenRecord 记录Token记录指标
func RecordTokenRecord(tenantID, modelProvider, modelName string, inputTokens, outputTokens int, cost float64, duration float64, isCached bool) {
	TokenRecordsTotal.WithLabelValues(tenantID, modelProvider, modelName).Inc()
	TokensTotal.WithLabelValues(tenantID, "input").Add(float64(inputTokens))
	TokensTotal.WithLabelValues(tenantID, "output").Add(float64(outputTokens))
	CostTotal.WithLabelValues(tenantID, modelProvider).Add(cost)
	TokenRecordDuration.WithLabelValues("record").Observe(duration)

	if isCached {
		CachedRequestsTotal.WithLabelValues(tenantID, modelProvider).Inc()
	}
}

// RecordBatchTokenRecord 记录批量Token记录指标
func RecordBatchTokenRecord(tenantID string, count int, totalCost float64, duration float64) {
	for i := 0; i < count; i++ {
		TokenRecordDuration.WithLabelValues("batch_record").Observe(duration)
	}
}

// RecordPricingCalculation 记录定价计算指标
func RecordPricingCalculation(modelProvider, modelName, result string, duration float64, unitPrice float64, tokenType string) {
	PricingCalculationsTotal.WithLabelValues(modelProvider, modelName, result).Inc()
	PricingCalculationDuration.WithLabelValues(modelProvider).Observe(duration)
	UnitPriceGauge.WithLabelValues(modelProvider, modelName, tokenType).Set(unitPrice)
}

// RecordPricingError 记录定价错误指标
func RecordPricingError(modelProvider, modelName, errorType string) {
	PricingErrorsTotal.WithLabelValues(modelProvider, modelName, errorType).Inc()
}

// RecordBudgetCheck 记录预算检查指标
func RecordBudgetCheck(tenantID string, triggered bool, duration float64) {
	triggeredStr := "false"
	if triggered {
		triggeredStr = "true"
	}
	BudgetChecksTotal.WithLabelValues(tenantID, triggeredStr).Inc()
	BudgetCheckDuration.WithLabelValues(tenantID).Observe(duration)
}

// RecordBudgetAlert 记录预算告警指标
func RecordBudgetAlert(tenantID, alertType, alertLevel string) {
	BudgetAlertsTotal.WithLabelValues(tenantID, alertType, alertLevel).Inc()
}

// UpdateBudgetUsage 更新预算使用率指标
func UpdateBudgetUsage(tenantID, budgetType, alertLevel string, usagePercent, budgetAmount, usedAmount float64) {
	BudgetUsageGauge.WithLabelValues(tenantID, budgetType, alertLevel).Set(usagePercent)
	BudgetAmountGauge.WithLabelValues(tenantID, budgetType).Set(budgetAmount)
	BudgetUsedAmountGauge.WithLabelValues(tenantID, budgetType).Set(usedAmount)
}

// RecordSummaryUpdate 记录汇总更新指标
func RecordSummaryUpdate(tenantID, result string, duration float64) {
	SummaryUpdateTotal.WithLabelValues(tenantID, result).Inc()
	SummaryUpdateDuration.WithLabelValues(tenantID).Observe(duration)
}

// UpdateDailyMetrics 更新每日指标
func UpdateDailyMetrics(tenantID, botID string, tokens int64, cost float64) {
	if botID != "" {
		DailyTokensGauge.WithLabelValues(tenantID, botID).Set(float64(tokens))
		DailyCostGauge.WithLabelValues(tenantID, botID).Set(cost)
	}
}

// RecordTokenLogDBOperation 记录Token日志数据库操作指标
func RecordTokenLogDBOperation(operation, result string, duration float64) {
	TokenLogDBOperationsTotal.WithLabelValues(operation, result).Inc()
	TokenLogDBOperationDuration.WithLabelValues(operation).Observe(duration)
}

// RecordAlertNotification 记录告警通知发送指标
func RecordAlertNotification(tenantID, channel, result string, duration float64) {
	AlertNotificationTotal.WithLabelValues(tenantID, channel, result).Inc()
	AlertNotificationDuration.WithLabelValues(channel).Observe(duration)
}

// UpdateModelUsage 更新模型使用指标
func UpdateModelUsage(tenantID, modelProvider, modelName string, tokens int64, cost float64) {
	ModelUsageGauge.WithLabelValues(tenantID, modelProvider, modelName).Set(cost)
}

// UpdateActiveTenants 更新活跃租户数
func UpdateActiveTenants(count int) {
	ActiveTenantsGauge.Set(float64(count))
}

// UpdateTotalRevenue 更新总收入
func UpdateTotalRevenue(period string, revenue float64) {
	TotalRevenueGauge.WithLabelValues(period).Set(revenue)
}
