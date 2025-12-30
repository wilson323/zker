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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ========== 租户管理模块指标 ==========

var (
	// TenantCreationTotal 租户创建总数(Counter)
	TenantCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tenant_creation_total",
			Help: "Total number of tenant creations",
		},
		[]string{"tenant_type", "result"}, // tenant_type: individual, team, enterprise; result: success, failure
	)

	// TenantDeletionTotal 租户删除总数(Counter)
	TenantDeletionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tenant_deletion_total",
			Help: "Total number of tenant deletions",
		},
		[]string{"result"},
	)

	// TenantUpgradeTotal 租户订阅升级总数(Counter)
	TenantUpgradeTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tenant_upgrade_total",
			Help: "Total number of tenant subscription upgrades",
		},
		[]string{"from_tier", "to_tier", "result"},
	)

	// SubscriptionStatusCurrent 订阅状态(Gauge)
	SubscriptionStatusCurrent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subscription_status_current",
			Help: "Current subscription status",
		},
		[]string{"tier", "status"}, // tier: free, pro, enterprise; status: active, suspended, cancelled, expired
	)

	// QuotaUsagePercent 配额使用百分比(Gauge)
	QuotaUsagePercent = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_usage_percent",
			Help: "Quota usage percentage",
		},
		[]string{"tenant_id", "resource_type", "alert_level"}, // alert_level: normal, warning, critical, exceeded
	)

	// QuotaExceededTotal 配额超限次数(Counter)
	QuotaExceededTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "quota_exceeded_total",
			Help: "Total number of quota exceeded events",
		},
		[]string{"tenant_id", "resource_type"},
	)

	// InvoiceCreationTotal 账单创建总数(Counter)
	InvoiceCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "invoice_creation_total",
			Help: "Total number of invoice creations",
		},
		[]string{"tenant_id", "currency", "result"},
	)

	// InvoicePaymentTotal 账单支付总数(Counter)
	InvoicePaymentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "invoice_payment_total",
			Help: "Total number of invoice payments",
		},
		[]string{"tenant_id", "result"}, // result: success, failure
	)

	// InvoiceOverdueTotal 账单逾期总数(Gauge)
	InvoiceOverdueTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "invoice_overdue_total",
			Help: "Total number of overdue invoices",
		},
		[]string{"tenant_id", "tier"},
	)
)

// ========== 路由模块指标 ==========

var (
	// RoutingRuleTotal 路由规则总数(Gauge)
	RoutingRuleTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "routing_rule_total",
			Help: "Total number of routing rules",
		},
		[]string{"tenant_id", "rule_type", "is_active"}, // rule_type: keyword, regex, intent, category
	)

	// RoutingExecutionTotal 路由执行总数(Counter)
	RoutingExecutionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "routing_execution_total",
			Help: "Total number of routing executions",
		},
		[]string{"tenant_id", "match_type", "result"}, // match_type: keyword, similarity, model, none; result: success, no_match
	)

	// RoutingDuration 路由决策延迟(Histogram)
	RoutingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "routing_duration_seconds",
			Help:    "Routing decision latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"tenant_id", "match_type"},
	)

	// RoutingConfidence 路由置信度(Histogram)
	RoutingConfidence = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "routing_confidence",
			Help:    "Routing decision confidence",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.8, 0.9, 0.95, 0.99, 1.0},
		},
		[]string{"tenant_id", "match_type"},
	)

	// RoutingScore 路由评分(Histogram)
	RoutingScore = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "routing_score",
			Help:    "Routing score",
			Buckets: []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0},
		},
		[]string{"tenant_id"},
	)

	// RoutingNoMatchTotal 路由无匹配次数(Counter)
	RoutingNoMatchTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "routing_no_match_total",
			Help: "Total number of routing executions with no match",
		},
		[]string{"tenant_id"},
	)

	// BotHealthStatus Bot健康状态(Gauge)
	BotHealthStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_health_status",
			Help: "Bot health status (1=healthy, 0=unhealthy)",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// BotSuccessRate Bot成功率(Gauge)
	BotSuccessRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_success_rate",
			Help: "Bot success rate (percentage)",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// BotCurrentLoad Bot当前负载(Gauge)
	BotCurrentLoad = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_current_load",
			Help: "Bot current load",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// BotMaxCapacity Bot最大容量(Gauge)
	BotMaxCapacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_max_capacity",
			Help: "Bot max capacity",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// CircuitBreakerStatus 熔断器状态(Gauge)
	CircuitBreakerStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_status",
			Help: "Circuit breaker status (0=closed, 1=open, 2=half_open)",
		},
		[]string{"service_id"},
	)

	// CircuitBreakerTripTotal 熔断器跳闸次数(Counter)
	CircuitBreakerTripTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_trip_total",
			Help: "Total number of circuit breaker trips",
		},
		[]string{"service_id"},
	)
)

// ========== 权限模块指标 ==========

var (
	// RoleTotal 角色总数(Gauge)
	RoleTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "role_total",
			Help: "Total number of roles",
		},
		[]string{"tenant_id", "role_type"}, // role_type: system, custom
	)

	// UserRoleAssignmentTotal 用户角色分配总数(Counter)
	UserRoleAssignmentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_role_assignment_total",
			Help: "Total number of user role assignments",
		},
		[]string{"tenant_id", "role_id", "result"}, // result: success, failure
	)

	// PermissionDeniedTotal 权限拒绝次数(Counter)
	PermissionDeniedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "permission_denied_total",
			Help: "Total number of permission denials",
		},
		[]string{"tenant_id", "permission_type", "reason"}, // reason: insufficient_scope, no_role, custom_filter_failed
	)

	// DataPermissionScopeTotal 数据权限范围分布(Gauge)
	DataPermissionScopeTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "data_permission_scope_total",
			Help: "Total number of data permissions by scope",
		},
		[]string{"tenant_id", "scope"}, // scope: all, dept_sub, dept, self, custom
	)
)

// ========== 辅助函数 ==========

// RecordTenantCreation 记录租户创建
func RecordTenantCreation(tenantType, result string) {
	TenantCreationTotal.WithLabelValues(tenantType, result).Inc()
}

// RecordTenantDeletion 记录租户删除
func RecordTenantDeletion(result string) {
	TenantDeletionTotal.WithLabelValues(result).Inc()
}

// RecordTenantUpgrade 记录租户升级
func RecordTenantUpgrade(fromTier, toTier, result string) {
	TenantUpgradeTotal.WithLabelValues(fromTier, toTier, result).Inc()
}

// RecordQuotaExceeded 记录配额超限
func RecordQuotaExceeded(tenantID, resourceType string) {
	QuotaExceededTotal.WithLabelValues(tenantID, resourceType).Inc()
}

// UpdateQuotaUsagePercent 更新配额使用百分比
func UpdateQuotaUsagePercent(tenantID, resourceType, alertLevel string, percent float64) {
	QuotaUsagePercent.WithLabelValues(tenantID, resourceType, alertLevel).Set(percent)
}

// RecordRoutingExecution 记录路由执行
func RecordRoutingExecution(tenantID, matchType, result string, duration float64, confidence, score float64) {
	RoutingExecutionTotal.WithLabelValues(tenantID, matchType, result).Inc()
	RoutingDuration.WithLabelValues(tenantID, matchType).Observe(duration)
	RoutingConfidence.WithLabelValues(tenantID, matchType).Observe(confidence)
	RoutingScore.WithLabelValues(tenantID).Observe(score)
}

// RecordRoutingNoMatch 记录路由无匹配
func RecordRoutingNoMatch(tenantID string) {
	RoutingNoMatchTotal.WithLabelValues(tenantID).Inc()
}

// UpdateBotHealth 更新Bot健康状态
func UpdateBotHealth(tenantID, botID string, isHealthy bool, successRate, currentLoad, maxCapacity int) {
	status := 0.0
	if isHealthy {
		status = 1.0
	}
	BotHealthStatus.WithLabelValues(tenantID, botID).Set(status)
	BotSuccessRate.WithLabelValues(tenantID, botID).Set(float64(successRate))
	BotCurrentLoad.WithLabelValues(tenantID, botID).Set(float64(currentLoad))
	BotMaxCapacity.WithLabelValues(tenantID, botID).Set(float64(maxCapacity))
}

// UpdateCircuitBreakerStatus 更新熔断器状态
func UpdateCircuitBreakerStatus(serviceID string, status int) {
	// status: 0=closed, 1=open, 2=half_open
	CircuitBreakerStatus.WithLabelValues(serviceID).Set(float64(status))
}

// RecordCircuitBreakerTrip 记录熔断器跳闸
func RecordCircuitBreakerTrip(serviceID string) {
	CircuitBreakerTripTotal.WithLabelValues(serviceID).Inc()
}

// RecordUserRoleAssignment 记录用户角色分配
func RecordUserRoleAssignment(tenantID, roleID, result string) {
	UserRoleAssignmentTotal.WithLabelValues(tenantID, roleID, result).Inc()
}

// RecordPermissionDenied 记录权限拒绝
func RecordPermissionDenied(tenantID, permissionType, reason string) {
	PermissionDeniedTotal.WithLabelValues(tenantID, permissionType, reason).Inc()
}

// RecordInvoiceCreation 记录账单创建
func RecordInvoiceCreation(tenantID, currency string, result string) {
	InvoiceCreationTotal.WithLabelValues(tenantID, currency, result).Inc()
}

// RecordInvoicePayment 记录账单支付
func RecordInvoicePayment(tenantID, result string) {
	InvoicePaymentTotal.WithLabelValues(tenantID, result).Inc()
}

// UpdateInvoiceOverdueCount 更新逾期账单数
func UpdateInvoiceOverdueCount(tenantID, tier string, count int) {
	InvoiceOverdueTotal.WithLabelValues(tenantID, tier).Set(float64(count))
}
