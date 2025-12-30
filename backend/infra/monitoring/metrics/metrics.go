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

// ========== HTTP指标 ==========

var (
	// HTTPRequestsTotal HTTP请求总数
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration HTTP请求延迟
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	// HTTPRequestSize HTTP请求大小
	HTTPRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"endpoint"},
	)

	// HTTPResponseSize HTTP响应大小
	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 1000, 10000, 100000, 1000000},
		},
		[]string{"endpoint"},
	)
)

// ========== 业务指标 ==========

var (
	// QuotaCheckDuration 配额检查延迟
	QuotaCheckDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "quota_check_duration_seconds",
			Help:    "Quota check latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
		},
		[]string{"resource_type", "result"},
	)

	// QuotaUsage 配额使用量(Gauge)
	QuotaUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_usage",
			Help: "Current quota usage",
		},
		[]string{"tenant_id", "resource_type"},
	)

	// QuotaLimit 配额限制(Gauge)
	QuotaLimit = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "quota_limit",
			Help: "Quota limit",
		},
		[]string{"tenant_id", "resource_type"},
	)

	// QuotaCheckTotal 配额检查总数
	QuotaCheckTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "quota_check_total",
			Help: "Total number of quota checks",
		},
		[]string{"resource_type", "result"},
	)

	// TenantActive 活跃租户数(Gauge)
	TenantActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "tenant_active_count",
			Help: "Number of active tenants",
		},
	)

	// TenantSuspended 暂停租户数(Gauge)
	TenantSuspended = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "tenant_suspended_count",
			Help: "Number of suspended tenants",
		},
	)

	// SubscriptionExpiring 订阅即将过期数量(Gauge)
	SubscriptionExpiring = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subscription_expiring_count",
			Help: "Number of subscriptions expiring soon",
		},
		[]string{"tier", "days"}, // days: 7, 14, 30
	)

	// SubscriptionRevenue 订阅收入(Gauge)
	SubscriptionRevenue = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "subscription_revenue_total",
			Help: "Total subscription revenue",
		},
		[]string{"tier", "currency"},
	)
)

// ========== 数据库指标 ==========

var (
	// DBConnectionsActive 数据库活跃连接数(Gauge)
	DBConnectionsActive = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
		[]string{"database"},
	)

	// DBConnectionsIdle 数据库空闲连接数(Gauge)
	DBConnectionsIdle = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
		[]string{"database"},
	)

	// DBQueryDuration 数据库查询延迟
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"database", "operation", "table"},
	)

	// DBQueryTotal 数据库查询总数
	DBQueryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_total",
			Help: "Total number of database queries",
		},
		[]string{"database", "operation", "table", "status"},
	)

	// DBTransactionDuration 数据库事务延迟
	DBTransactionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_transaction_duration_seconds",
			Help:    "Database transaction latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"database", "result"},
	)

	// DBTransactionTotal 数据库事务总数
	DBTransactionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_transaction_total",
			Help: "Total number of database transactions",
		},
		[]string{"database", "result"},
	)
)

// ========== 缓存指标 ==========

var (
	// CacheHitTotal 缓存命中总数
	CacheHitTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type", "key_prefix"},
	)

	// CacheMissTotal 缓存未命中总数
	CacheMissTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type", "key_prefix"},
	)

	// CacheDuration 缓存操作延迟
	CacheDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cache_duration_seconds",
			Help:    "Cache operation latency in seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05},
		},
		[]string{"cache_type", "operation"},
	)

	// CacheSize 缓存大小(Gauge)
	CacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Cache size in bytes",
		},
		[]string{"cache_type"},
	)

	// CacheItemCount 缓存项数量(Gauge)
	CacheItemCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_item_count",
			Help: "Number of items in cache",
		},
		[]string{"cache_type"},
	)
)

// ========== 错误指标 ==========

var (
	// ErrorTotal 错误总数
	ErrorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_total",
			Help: "Total number of errors",
		},
		[]string{"error_code", "error_type", "endpoint"},
	)

	// PanicTotal panic总数
	PanicTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "panic_total",
			Help: "Total number of panics",
		},
	)
)

// ========== 系统指标 ==========

var (
	// GoRoutines Goroutine数量(Gauge)
	GoRoutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sys_goroutines_count",
			Help: "Number of goroutines",
		},
	)

	// MemoryHeap 堆内存使用(Gauge)
	MemoryHeap = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "sys_memory_heap_bytes",
			Help: "Heap memory usage in bytes",
		},
	)

	// MemoryGC 垃圾回收次数(Counter)
	MemoryGC = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "sys_gc_total",
			Help: "Total number of garbage collections",
		},
	)
)

// ========== 辅助函数 ==========

// RecordHTTPRequest 记录HTTP请求
func RecordHTTPRequest(method, endpoint string, status int, duration float64, requestSize, responseSize int) {
	HTTPRequestsTotal.WithLabelValues(method, endpoint, formatStatus(status)).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	HTTPRequestSize.WithLabelValues(endpoint).Observe(float64(requestSize))
	HTTPResponseSize.WithLabelValues(endpoint).Observe(float64(responseSize))
}

// RecordQuotaCheck 记录配额检查
func RecordQuotaCheck(tenantID, resourceType string, duration float64, allowed bool, currentUsage, limit int) {
	result := "allowed"
	if !allowed {
		result = "denied"
	}

	QuotaCheckDuration.WithLabelValues(resourceType, result).Observe(duration)
	QuotaCheckTotal.WithLabelValues(resourceType, result).Inc()

	// 更新配额使用量
	QuotaUsage.WithLabelValues(tenantID, resourceType).Set(float64(currentUsage))
	QuotaLimit.WithLabelValues(tenantID, resourceType).Set(float64(limit))
}

// RecordDBQuery 记录数据库查询
func RecordDBQuery(database, operation, table string, duration float64, status string) {
	DBQueryDuration.WithLabelValues(database, operation, table).Observe(duration)
	DBQueryTotal.WithLabelValues(database, operation, table, status).Inc()
}

// RecordCacheOperation 记录缓存操作
func RecordCacheOperation(cacheType, operation string, duration float64, hit bool, keyPrefix string) {
	CacheDuration.WithLabelValues(cacheType, operation).Observe(duration)

	if hit {
		CacheHitTotal.WithLabelValues(cacheType, keyPrefix).Inc()
	} else {
		CacheMissTotal.WithLabelValues(cacheType, keyPrefix).Inc()
	}
}

// RecordError 记录错误
func RecordError(errorCode int32, errorType, endpoint string) {
	ErrorTotal.WithLabelValues(formatErrorCode(errorCode), errorType, endpoint).Inc()
}

// formatStatus 格式化HTTP状态码
func formatStatus(status int) string {
	if status >= 200 && status < 300 {
		return "2xx"
	} else if status >= 300 && status < 400 {
		return "3xx"
	} else if status >= 400 && status < 500 {
		return "4xx"
	} else if status >= 500 {
		return "5xx"
	}
	return "unknown"
}

// formatErrorCode 格式化错误码
func formatErrorCode(code int32) string {
	// 根据错误码范围分类
	if code >= 200000000 && code < 300000000 {
		return "2xx_tenant"
	} else if code >= 300000000 && code < 400000000 {
		return "3xx_quota"
	} else if code >= 400000000 && code < 500000000 {
		return "4xx_subscription"
	} else if code >= 700000000 && code < 800000000 {
		return "7xx_user"
	}
	return "other"
}

// ========== 初始化函数 ==========

// init 强制初始化所有指标变量
//
// 注意：在 Go 中，同一个包的多个文件的包级变量初始化顺序是不确定的。
// 这个函数通过引用其他文件中的指标变量，确保它们在 Prometheus 注册表中完成注册。
func init() {
	// 引用基础指标文件中的变量（metrics.go本身）
	_ = HTTPRequestsTotal
	_ = HTTPRequestDuration
	_ = HTTPRequestSize
	_ = HTTPResponseSize
	_ = QuotaCheckDuration
	_ = QuotaUsage
	_ = QuotaLimit
	_ = QuotaCheckTotal
	_ = TenantActive
	_ = TenantSuspended
	_ = SubscriptionExpiring
	_ = SubscriptionRevenue
	_ = DBConnectionsActive
	_ = DBConnectionsIdle
	_ = DBQueryDuration
	_ = DBQueryTotal
	_ = DBTransactionDuration
	_ = DBTransactionTotal
	_ = CacheHitTotal
	_ = CacheMissTotal
	_ = CacheDuration
	_ = CacheSize
	_ = CacheItemCount
	_ = ErrorTotal
	_ = PanicTotal
	_ = GoRoutines
	_ = MemoryHeap
	_ = MemoryGC

	// 引用业务指标文件中的变量（business_metrics.go）
	_ = BotTotal
	_ = BotInvocationTotal
	_ = BotInvocationDuration
	_ = BotErrorRate
	_ = BotTokenUsage
	_ = BotCost
	_ = WorkflowTotal
	_ = WorkflowExecutionTotal
	_ = WorkflowExecutionDuration
	_ = WorkflowNodeExecutionTotal
	_ = WorkflowStepDuration
	_ = ConversationTotal
	_ = MessageTotal
	_ = MessageLatency
	_ = ConversationLength
	_ = KnowledgeBaseTotal
	_ = KnowledgeBaseDocumentTotal
	_ = KnowledgeBaseSearchTotal
	_ = KnowledgeBaseSearchDuration
	_ = KnowledgeBaseIndexSize
	_ = PermissionCheckTotal
	_ = PermissionCheckDuration
	_ = RoleAssignmentTotal
	_ = UserRoleTotal
	_ = StorageUsageTotal
	_ = StorageQuota
	_ = StorageUsageRate
	_ = ObjectUploadTotal
	_ = ObjectDownloadTotal
	_ = ObjectTransferDuration
	_ = ActiveUserTotal
	_ = UserSessionTotal
	_ = UserLoginTotal
	_ = UserOperationTotal
	_ = ModelInvocationTotal
	_ = ModelTokenUsageTotal
	_ = ModelInvocationDuration
	_ = ModelCostTotal
	_ = ModelErrorRate
	_ = RateLimitExceededTotal
	_ = RateLimitRemaining
	_ = ThrottleRequestTotal
	_ = BusinessHealthScore
	_ = SLAComplianceRate
	_ = IncidentTotal

	// 引用数据库读写指标文件中的变量（db_readwrite_metrics.go）
	_ = DBWriteQPS
	_ = DBReadQPS
	_ = DBReplicationLag
	_ = DBWriteLatency
	_ = DBReadLatency
	_ = DBConnectionPoolUsage
	_ = DBReplicationStatus

	// 引用企业级指标文件中的变量（enterprise_metrics.go）
	_ = TenantCreationTotal
	_ = TenantDeletionTotal
	_ = TenantUpgradeTotal
	_ = SubscriptionStatusCurrent
	_ = QuotaUsagePercent
	_ = QuotaExceededTotal
	_ = InvoiceCreationTotal
	_ = InvoicePaymentTotal
	_ = InvoiceOverdueTotal
	_ = RoutingRuleTotal
	_ = RoutingExecutionTotal
	_ = RoutingDuration
	_ = RoutingConfidence
	_ = RoutingScore
	_ = RoutingNoMatchTotal
	_ = BotHealthStatus
	_ = BotSuccessRate
	_ = BotCurrentLoad
	_ = BotMaxCapacity
	_ = CircuitBreakerStatus
	_ = CircuitBreakerTripTotal
	_ = RoleTotal
	_ = UserRoleAssignmentTotal
	_ = PermissionDeniedTotal
	_ = DataPermissionScopeTotal

	// 引用组织架构指标文件中的变量（org_metrics.go）
	_ = OrganizationTotal
	_ = OrganizationCreationTotal
	_ = OrganizationDeletionTotal
	_ = OrganizationUpdateTotal
	_ = OrganizationMoveTotal
	_ = DepartmentTotal
	_ = DepartmentCreationTotal
	_ = DepartmentDeletionTotal
	_ = DepartmentMemberTotal
	_ = EmployeeTotal
	_ = EmployeeCreationTotal
	_ = EmployeeDeletionTotal
	_ = EmployeeUpdateTotal
	_ = EmployeeTransferTotal
	_ = EmployeePromotionTotal
	_ = EmployeeResignationTotal
	_ = EmployeeOnboardingTotal
	_ = EmployeeLoginTotal
	_ = PositionTotal
	_ = PositionCreationTotal
	_ = PositionDeletionTotal
	_ = PositionOccupiedTotal
	_ = PositionVacantTotal
	_ = OrgOrphanTotal
	_ = HRProcessTotal

	// 引用组件指标文件中的变量（component_metrics.go）
	_ = TenantIsolationTotal
	_ = TenantIsolationDuration
	_ = TenantValidationTotal
	_ = TenantValidationDuration
	_ = TenantCacheHitRate
	_ = TenantIDExtractionTotal
	_ = ErrorCodeTotal
	_ = ErrorCodeByModule
	_ = DeprecatedErrorCodeUsage
	_ = ErrorCreationDuration
	_ = ErrorDistribution
	_ = SessionCreationTotal
	_ = SessionCreationDuration
	_ = SessionValidationTotal
	_ = SessionValidationDuration
	_ = SessionActiveCount
	_ = SessionExpiredCount
	_ = SessionWithTenantIDCount
	_ = SessionAccessTotal
	_ = BotOperationTotal
	_ = BotOperationDuration
	_ = BotExecutionTotal
	_ = BotExecutionDuration
	_ = BotActiveCount
	_ = BotQuotaUsage
	_ = ConversationOperationTotal
	_ = ConversationOperationDuration
	_ = MessageDuration
	_ = AgentRunTotal
	_ = AgentRunDuration
	_ = ConversationActiveCount
	_ = WorkflowOperationTotal
	_ = WorkflowOperationDuration
	_ = WorkflowNodeExecutionDuration
	_ = WorkflowActiveCount
	_ = WorkflowVariableTotal
	_ = WorkflowVersionTotal
	_ = ContextCacheOperationTotal
}
