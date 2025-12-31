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

package collector

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// EnterpriseMetricsCollector 企业级指标收集器
// 用于收集和暴露企业级监控指标到 Prometheus
type EnterpriseMetricsCollector struct {
	logger *zap.Logger
	mu     sync.RWMutex

	// 性能指标
	requestDuration    *prometheus.HistogramVec
	requestCount       *prometheus.CounterVec
	errorCount         *prometheus.CounterVec

	// 业务指标
	activeUsersGauge   prometheus.Gauge
	conversationCount  *prometheus.CounterVec

	// 资源指标
	tokenUsageGauge    *prometheus.GaugeVec
	quotaUsageGauge    *prometheus.GaugeVec

	// Agent指标
	agentQPSGauge              *prometheus.GaugeVec
	agentResponseTimeGauge     *prometheus.GaugeVec
	agentErrorRateGauge        *prometheus.GaugeVec
	agentSatisfactionGauge     *prometheus.GaugeVec

	// 租户隔离指标
	crossTenantAccessCount     *prometheus.CounterVec
	tenantIsolationViolation   *prometheus.CounterVec

	// RBAC指标
	permissionCheckCount       *prometheus.CounterVec
	permissionDenyCount        *prometheus.CounterVec
	roleCheckDuration          *prometheus.HistogramVec

	// 数据库指标
	dbQueryDuration            *prometheus.HistogramVec
	dbQueryCount               *prometheus.CounterVec
	dbConnectionPoolGauge      *prometheus.GaugeVec

	// 缓存指标
	cacheHitCount              *prometheus.CounterVec
	cacheMissCount             *prometheus.CounterVec
	cacheDuration              *prometheus.HistogramVec

	// 工作流指标
	workflowExecutionCount     *prometheus.CounterVec
	workflowExecutionDuration  *prometheus.HistogramVec
	workflowNodeExecutionCount *prometheus.CounterVec

	// 知识库指标
	knowledgeQueryCount         *prometheus.CounterVec
	knowledgeQueryDuration      *prometheus.HistogramVec
	knowledgeDocumentCountGauge *prometheus.GaugeVec

	// 配额指标
	quotaUsagePercent          *prometheus.GaugeVec
	quotaExceededCount          *prometheus.CounterVec
}

// NewEnterpriseMetricsCollector 创建企业级指标收集器
func NewEnterpriseMetricsCollector(logger *zap.Logger) *EnterpriseMetricsCollector {
	c := &EnterpriseMetricsCollector{
		logger: logger,

		// 性能指标
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "zker_api_request_duration_seconds",
				Help:    "API request duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"endpoint", "method", "tenant_id", "status_code"},
		),

		requestCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_api_requests_total",
				Help: "Total number of API requests",
			},
			[]string{"endpoint", "method", "tenant_id", "status_code"},
		),

		errorCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_api_errors_total",
				Help: "Total number of API errors",
			},
			[]string{"endpoint", "method", "tenant_id", "error_type"},
		),

		// 业务指标
		activeUsersGauge: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "zker_active_users",
				Help: "Number of active users",
			},
		),

		conversationCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_conversations_total",
				Help: "Total number of conversations",
			},
			[]string{"tenant_id", "agent_id"},
		),

		// 资源指标
		tokenUsageGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_token_usage",
				Help: "Token usage count",
			},
			[]string{"tenant_id", "agent_id", "model"},
		),

		quotaUsageGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_quota_usage",
				Help: "Quota usage count",
			},
			[]string{"tenant_id", "resource_type"},
		),

		// Agent指标
		agentQPSGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_agent_qps",
				Help: "Agent queries per second",
			},
			[]string{"tenant_id", "agent_id"},
		),

		agentResponseTimeGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_agent_response_time_ms",
				Help: "Agent response time in milliseconds",
			},
			[]string{"tenant_id", "agent_id"},
		),

		agentErrorRateGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_agent_error_rate",
				Help: "Agent error rate (0-1)",
			},
			[]string{"tenant_id", "agent_id"},
		),

		agentSatisfactionGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_agent_satisfaction_score",
				Help: "Agent satisfaction score (1-5)",
			},
			[]string{"tenant_id", "agent_id"},
		),

		// 租户隔离指标
		crossTenantAccessCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_cross_tenant_access_total",
				Help: "Total number of cross-tenant access attempts",
			},
			[]string{"source_tenant", "target_tenant", "access_type"},
		),

		tenantIsolationViolation: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_tenant_isolation_violations_total",
				Help: "Total number of tenant isolation violations",
			},
			[]string{"violation_type", "severity"},
		),

		// RBAC指标
		permissionCheckCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_permission_checks_total",
				Help: "Total number of permission checks",
			},
			[]string{"tenant_id", "permission", "result"},
		),

		permissionDenyCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_permission_denied_total",
				Help: "Total number of denied permissions",
			},
			[]string{"tenant_id", "permission", "reason"},
		),

		roleCheckDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "zker_role_check_duration_seconds",
				Help:    "Role check duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"tenant_id"},
		),

		// 数据库指标
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "zker_db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5},
			},
			[]string{"database", "operation", "table"},
		),

		dbQueryCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"database", "operation", "table", "status"},
		),

		dbConnectionPoolGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_db_connection_pool",
				Help: "Database connection pool size",
			},
			[]string{"database", "state"}, // state: idle, in_use, max
		),

		// 缓存指标
		cacheHitCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_cache_hits_total",
				Help: "Total number of cache hits",
			},
			[]string{"cache_type", "tenant_id"},
		),

		cacheMissCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_cache_misses_total",
				Help: "Total number of cache misses",
			},
			[]string{"cache_type", "tenant_id"},
		),

		cacheDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "zker_cache_duration_seconds",
				Help:    "Cache operation duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"cache_type", "operation"}, // operation: get, set, delete
		),

		// 工作流指标
		workflowExecutionCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_workflow_executions_total",
				Help: "Total number of workflow executions",
			},
			[]string{"tenant_id", "workflow_id", "status"},
		),

		workflowExecutionDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "zker_workflow_execution_duration_seconds",
				Help:    "Workflow execution duration in seconds",
				Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60, 120, 300},
			},
			[]string{"tenant_id", "workflow_id"},
		),

		workflowNodeExecutionCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_workflow_node_executions_total",
				Help: "Total number of workflow node executions",
			},
			[]string{"tenant_id", "workflow_id", "node_type", "status"},
		),

		// 知识库指标
		knowledgeQueryCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_knowledge_queries_total",
				Help: "Total number of knowledge base queries",
			},
			[]string{"tenant_id", "knowledge_id", "status"},
		),

		knowledgeQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "zker_knowledge_query_duration_seconds",
				Help:    "Knowledge base query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"tenant_id", "knowledge_id"},
		),

		knowledgeDocumentCountGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_knowledge_documents",
				Help: "Number of documents in knowledge base",
			},
			[]string{"tenant_id", "knowledge_id"},
		),

		// 配额指标
		quotaUsagePercent: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "zker_quota_usage_percentage",
				Help: "Quota usage percentage",
			},
			[]string{"tenant_id", "resource_type", "alert_level"}, // alert_level: normal, warning, critical, exceeded
		),

		quotaExceededCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "zker_quota_exceeded_total",
				Help: "Total number of quota exceeded events",
			},
			[]string{"tenant_id", "resource_type"},
		),
	}

	return c
}

// ============================================================
// API指标记录方法
// ============================================================

// RecordRequest 记录API请求
func (c *EnterpriseMetricsCollector) RecordRequest(
	endpoint, method, tenantID, statusCode string,
	duration time.Duration,
) {
	c.requestCount.WithLabelValues(endpoint, method, tenantID, statusCode).Inc()
	c.requestDuration.WithLabelValues(endpoint, method, tenantID, statusCode).Observe(duration.Seconds())
}

// RecordError 记录API错误
func (c *EnterpriseMetricsCollector) RecordError(
	endpoint, method, tenantID, errorType string,
) {
	c.errorCount.WithLabelValues(endpoint, method, tenantID, errorType).Inc()
}

// ============================================================
// 业务指标记录方法
// ============================================================

// SetActiveUsers 设置活跃用户数
func (c *EnterpriseMetricsCollector) SetActiveUsers(count float64) {
	c.activeUsersGauge.Set(count)
}

// RecordConversation 记录对话
func (c *EnterpriseMetricsCollector) RecordConversation(tenantID, agentID string) {
	c.conversationCount.WithLabelValues(tenantID, agentID).Inc()
}

// ============================================================
// Agent指标记录方法
// ============================================================

// RecordAgentMetrics 记录Agent指标
func (c *EnterpriseMetricsCollector) RecordAgentMetrics(
	tenantID, agentID string,
	qps, responseTime, errorRate, satisfaction float64,
) {
	c.agentQPSGauge.WithLabelValues(tenantID, agentID).Set(qps)
	c.agentResponseTimeGauge.WithLabelValues(tenantID, agentID).Set(responseTime)
	c.agentErrorRateGauge.WithLabelValues(tenantID, agentID).Set(errorRate)
	c.agentSatisfactionGauge.WithLabelValues(tenantID, agentID).Set(satisfaction)
}

// ============================================================
// 租户隔离指标记录方法
// ============================================================

// RecordCrossTenantAccess 记录跨租户访问
func (c *EnterpriseMetricsCollector) RecordCrossTenantAccess(
	sourceTenant, targetTenant, accessType string,
) {
	c.crossTenantAccessCount.WithLabelValues(sourceTenant, targetTenant, accessType).Inc()
}

// RecordTenantIsolationViolation 记录租户隔离违规
func (c *EnterpriseMetricsCollector) RecordTenantIsolationViolation(
	violationType, severity string,
) {
	c.tenantIsolationViolation.WithLabelValues(violationType, severity).Inc()
}

// ============================================================
// RBAC指标记录方法
// ============================================================

// RecordPermissionCheck 记录权限检查
func (c *EnterpriseMetricsCollector) RecordPermissionCheck(
	tenantID, permission, result string,
	duration time.Duration,
) {
	c.permissionCheckCount.WithLabelValues(tenantID, permission, result).Inc()
	c.roleCheckDuration.WithLabelValues(tenantID).Observe(duration.Seconds())

	if result == "deny" {
		c.permissionDenyCount.WithLabelValues(tenantID, permission, "").Inc()
	}
}

// ============================================================
// 数据库指标记录方法
// ============================================================

// RecordDBQuery 记录数据库查询
func (c *EnterpriseMetricsCollector) RecordDBQuery(
	database, operation, table, status string,
	duration time.Duration,
) {
	c.dbQueryCount.WithLabelValues(database, operation, table, status).Inc()
	c.dbQueryDuration.WithLabelValues(database, operation, table).Observe(duration.Seconds())
}

// SetDBConnectionPool 设置数据库连接池大小
func (c *EnterpriseMetricsCollector) SetDBConnectionPool(
	database, state string,
	size float64,
) {
	c.dbConnectionPoolGauge.WithLabelValues(database, state).Set(size)
}

// ============================================================
// 缓存指标记录方法
// ============================================================

// RecordCacheHit 记录缓存命中
func (c *EnterpriseMetricsCollector) RecordCacheHit(cacheType, tenantID string) {
	c.cacheHitCount.WithLabelValues(cacheType, tenantID).Inc()
}

// RecordCacheMiss 记录缓存未命中
func (c *EnterpriseMetricsCollector) RecordCacheMiss(cacheType, tenantID string) {
	c.cacheMissCount.WithLabelValues(cacheType, tenantID).Inc()
}

// RecordCacheOperation 记录缓存操作
func (c *EnterpriseMetricsCollector) RecordCacheOperation(
	cacheType, operation string,
	duration time.Duration,
) {
	c.cacheDuration.WithLabelValues(cacheType, operation).Observe(duration.Seconds())
}

// ============================================================
// 工作流指标记录方法
// ============================================================

// RecordWorkflowExecution 记录工作流执行
func (c *EnterpriseMetricsCollector) RecordWorkflowExecution(
	tenantID, workflowID, status string,
	duration time.Duration,
) {
	c.workflowExecutionCount.WithLabelValues(tenantID, workflowID, status).Inc()
	c.workflowExecutionDuration.WithLabelValues(tenantID, workflowID).Observe(duration.Seconds())
}

// RecordWorkflowNodeExecution 记录工作流节点执行
func (c *EnterpriseMetricsCollector) RecordWorkflowNodeExecution(
	tenantID, workflowID, nodeType, status string,
) {
	c.workflowNodeExecutionCount.WithLabelValues(tenantID, workflowID, nodeType, status).Inc()
}

// ============================================================
// 知识库指标记录方法
// ============================================================

// RecordKnowledgeQuery 记录知识库查询
func (c *EnterpriseMetricsCollector) RecordKnowledgeQuery(
	tenantID, knowledgeID, status string,
	duration time.Duration,
) {
	c.knowledgeQueryCount.WithLabelValues(tenantID, knowledgeID, status).Inc()
	c.knowledgeQueryDuration.WithLabelValues(tenantID, knowledgeID).Observe(duration.Seconds())
}

// SetKnowledgeDocumentCount 设置知识库文档数量
func (c *EnterpriseMetricsCollector) SetKnowledgeDocumentCount(
	tenantID, knowledgeID string,
	count float64,
) {
	c.knowledgeDocumentCountGauge.WithLabelValues(tenantID, knowledgeID).Set(count)
}

// ============================================================
// 配额指标记录方法
// ============================================================

// UpdateQuotaUsagePercent 更新配额使用百分比
func (c *EnterpriseMetricsCollector) UpdateQuotaUsagePercent(
	tenantID, resourceType, alertLevel string,
	percentage float64,
) {
	c.quotaUsagePercent.WithLabelValues(tenantID, resourceType, alertLevel).Set(percentage)
}

// RecordQuotaExceeded 记录配额超限
func (c *EnterpriseMetricsCollector) RecordQuotaExceeded(
	tenantID, resourceType string,
) {
	c.quotaExceededCount.WithLabelValues(tenantID, resourceType).Inc()
}

// ============================================================
// Token使用指标记录方法
// ============================================================

// RecordTokenUsage 记录Token使用
func (c *EnterpriseMetricsCollector) RecordTokenUsage(
	tenantID, agentID, model string,
	count float64,
) {
	c.tokenUsageGauge.WithLabelValues(tenantID, agentID, model).Set(count)
}

// ============================================================
// 资源配额指标记录方法
// ============================================================

// RecordQuotaUsage 记录配额使用
func (c *EnterpriseMetricsCollector) RecordQuotaUsage(
	tenantID, resourceType string,
	usage float64,
) {
	c.quotaUsageGauge.WithLabelValues(tenantID, resourceType).Set(usage)
}
