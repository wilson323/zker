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

// =====================================================
// 租户隔离中间件指标
// =====================================================

var (
	// TenantIsolationTotal 租户隔离中间件处理总数
	TenantIsolationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tenant_isolation_total",
			Help: "Total number of tenant isolation middleware executions",
		},
		[]string{"source"}, // source: header, session, jwt, default
	)

	// TenantIsolationDuration 租户隔离中间件处理延迟
	TenantIsolationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tenant_isolation_duration_seconds",
			Help:    "Tenant isolation middleware processing latency in seconds",
			Buckets: []float64{0.000001, 0.000005, 0.00001, 0.00005, 0.0001, 0.0005, 0.001},
		},
		[]string{"source"},
	)

	// TenantValidationTotal 租户验证总数
	TenantValidationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tenant_validation_total",
			Help: "Total number of tenant validations",
		},
		[]string{"result"}, // result: success, not_found, suspended, deleted
	)

	// TenantValidationDuration 租户验证延迟
	TenantValidationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tenant_validation_duration_seconds",
			Help:    "Tenant validation latency in seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1},
		},
		[]string{"cache_hit"}, // cache_hit: true, false
	)

	// TenantCacheHitRate 租户缓存命中率(Gauge)
	TenantCacheHitRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "tenant_cache_hit_rate",
			Help: "Tenant cache hit rate (0-1)",
		},
	)

	// TenantIDExtractionTotal tenant_id提取总数
	TenantIDExtractionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tenant_id_extraction_total",
			Help: "Total number of tenant_id extractions",
		},
		[]string{"source"}, // source: header, session, jwt, default
	)
)

// =====================================================
// 错误码系统指标
// =====================================================

var (
	// ErrorCodeTotal 错误码出现总数（按错误码段分类）
	ErrorCodeTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_code_total",
			Help: "Total number of errors by error code range",
		},
		[]string{"module", "code_range", "error_name"}, // module: bot, conversation, workflow, tenant
	)

	// ErrorCodeByModule 按模块统计的错误总数
	ErrorCodeByModule = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_code_by_module_total",
			Help: "Total number of errors by module",
		},
		[]string{"module", "severity"}, // severity: error, warning, info
	)

	// DeprecatedErrorCodeUsage 废弃错误码使用次数
	DeprecatedErrorCodeUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "deprecated_error_code_usage_total",
			Help: "Total number of deprecated error code usages",
		},
		[]string{"old_code_range", "new_code_range"},
	)

	// ErrorCreationDuration 错误码创建延迟
	ErrorCreationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "error_creation_duration_seconds",
			Help:    "Error creation latency in seconds",
			Buckets: []float64{0.0000001, 0.0000005, 0.000001, 0.000005, 0.00001, 0.00005},
		},
		[]string{"with_params"}, // with_params: true, false
	)

	// ErrorDistribution 错误分布（按错误码）
	ErrorDistribution = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_distribution_total",
			Help: "Error distribution by specific error code",
		},
		[]string{"module", "error_code", "error_name"},
	)
)

// =====================================================
// Session 操作指标
// =====================================================

var (
	// SessionCreationTotal Session创建总数
	SessionCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_creation_total",
			Help: "Total number of session creations",
		},
		[]string{"result"}, // result: success, failed
	)

	// SessionCreationDuration Session创建延迟
	SessionCreationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "session_creation_duration_seconds",
			Help:    "Session creation latency in seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05},
		},
		[]string{"result"},
	)

	// SessionValidationTotal Session验证总数
	SessionValidationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_validation_total",
			Help: "Total number of session validations",
		},
		[]string{"result"}, // result: valid, invalid, expired, not_found
	)

	// SessionValidationDuration Session验证延迟
	SessionValidationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "session_validation_duration_seconds",
			Help:    "Session validation latency in seconds",
			Buckets: []float64{0.00001, 0.00005, 0.0001, 0.0005, 0.001, 0.005},
		},
		[]string{"result"},
	)

	// SessionActiveCount 活跃Session数量(Gauge)
	SessionActiveCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "session_active_count",
			Help: "Number of active sessions",
		},
	)

	// SessionExpiredCount 过期Session数量(Gauge)
	SessionExpiredCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "session_expired_count",
			Help: "Number of expired sessions",
		},
	)

	// SessionWithTenantID 有tenant_id的Session数量(Gauge)
	SessionWithTenantIDCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "session_with_tenant_id_count",
			Help: "Number of sessions with tenant_id",
		},
		[]string{"tenant_type"}, // tenant_type: individual, team, enterprise
	)

	// SessionAccessTotal Session访问总数
	SessionAccessTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_access_total",
			Help: "Total number of session accesses",
		},
		[]string{"operation"}, // operation: get_tenant_id, has_tenant_id, validate
	)
)

// =====================================================
// Bot 模块指标
// =====================================================

var (
	// BotOperationTotal Bot操作总数
	BotOperationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_operation_total",
			Help: "Total number of Bot operations",
		},
		[]string{"operation", "result"}, // operation: create, update, delete, publish, unpublish
	)

	// BotOperationDuration Bot操作延迟
	BotOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bot_operation_duration_seconds",
			Help:    "Bot operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"operation"},
	)

	// BotExecutionTotal Bot执行总数
	BotExecutionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_execution_total",
			Help: "Total number of Bot executions",
		},
		[]string{"result"}, // result: success, failed, timeout, quota_exceeded
	)

	// BotExecutionDuration Bot执行延迟
	BotExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bot_execution_duration_seconds",
			Help:    "Bot execution latency in seconds",
			Buckets: []float64{0.1, 0.5, 1, 5, 10, 30, 60, 300},
		},
		[]string{"result"},
	)

	// BotActiveCount 活跃Bot数量(Gauge)
	BotActiveCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_active_count",
			Help: "Number of active Bots",
		},
		[]string{"status"}, // status: draft, published, archived
	)

	// BotQuotaUsage Bot配额使用量(Gauge)
	BotQuotaUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_quota_usage",
			Help: "Bot quota usage",
		},
		[]string{"tenant_id", "quota_type"}, // quota_type: bots, executions, messages
	)
)

// =====================================================
// Conversation 模块指标
// =====================================================

var (
	// ConversationOperationTotal Conversation操作总数
	ConversationOperationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "conversation_operation_total",
			Help: "Total number of Conversation operations",
		},
		[]string{"operation", "result"}, // operation: create, update, delete, close
	)

	// ConversationOperationDuration Conversation操作延迟
	ConversationOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conversation_operation_duration_seconds",
			Help:    "Conversation operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"operation"},
	)

	// MessageDuration 消息处理延迟
	MessageDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "message_duration_seconds",
			Help:    "Message processing latency in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
		[]string{"direction"},
	)

	// AgentRunTotal Agent运行总数
	AgentRunTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "agent_run_total",
			Help: "Total number of Agent runs",
		},
		[]string{"result"}, // result: success, failed, timeout, cancelled
	)

	// AgentRunDuration Agent运行延迟
	AgentRunDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "agent_run_duration_seconds",
			Help:    "Agent run latency in seconds",
			Buckets: []float64{0.5, 1, 5, 10, 30, 60, 300, 600},
		},
		[]string{"result"},
	)

	// ConversationActiveCount 活跃Conversation数量(Gauge)
	ConversationActiveCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conversation_active_count",
			Help: "Number of active conversations",
		},
		[]string{"status"}, // status: active, closed
	)
)

// =====================================================
// Workflow 模块指标
// =====================================================

var (
	// WorkflowOperationTotal Workflow操作总数
	WorkflowOperationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_operation_total",
			Help: "Total number of Workflow operations",
		},
		[]string{"operation", "result"}, // operation: create, update, delete, publish, unpublish
	)

	// WorkflowOperationDuration Workflow操作延迟
	WorkflowOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "workflow_operation_duration_seconds",
			Help:    "Workflow operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"operation"},
	)

	// WorkflowNodeExecutionDuration Workflow节点执行延迟
	WorkflowNodeExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "workflow_node_execution_duration_seconds",
			Help:    "Workflow node execution latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30},
		},
		[]string{"node_type"},
	)

	// WorkflowActiveCount 活跃Workflow数量(Gauge)
	WorkflowActiveCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "workflow_active_count",
			Help: "Number of active workflows",
		},
		[]string{"status"}, // status: draft, published, archived
	)

	// WorkflowVariableTotal Workflow变量操作总数
	WorkflowVariableTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_variable_operation_total",
			Help: "Total number of Workflow variable operations",
		},
		[]string{"operation", "result"}, // operation: get, set, delete
	)

	// WorkflowVersionTotal Workflow版本操作总数
	WorkflowVersionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_version_operation_total",
			Help: "Total number of Workflow version operations",
		},
		[]string{"operation"}, // operation: create, rollback, delete
	)
)

// =====================================================
// Context 操作指标
// =====================================================

var (
	// ContextCacheOperationTotal Context缓存操作总数
	ContextCacheOperationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "context_cache_operation_total",
			Help: "Total number of context cache operations",
		},
		[]string{"operation", "result"}, // operation: store, get, delete
	)

	// ContextCacheOperationDuration Context缓存操作延迟
	ContextCacheOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "context_cache_operation_duration_seconds",
			Help:    "Context cache operation latency in seconds",
			Buckets: []float64{0.0000001, 0.0000005, 0.000001, 0.000005, 0.00001, 0.00005},
		},
		[]string{"operation"},
	)

	// ContextCacheSize Context缓存大小(Gauge)
	ContextCacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "context_cache_size_bytes",
			Help: "Context cache size in bytes",
		},
	)

	// ContextCacheItemCount Context缓存项数量(Gauge)
	ContextCacheItemCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "context_cache_item_count",
			Help: "Number of items in context cache",
		},
	)
)

// =====================================================
// 辅助函数
// =====================================================

// RecordTenantIsolation 记录租户隔离中间件处理
func RecordTenantIsolation(source string, duration float64) {
	TenantIsolationTotal.WithLabelValues(source).Inc()
	TenantIsolationDuration.WithLabelValues(source).Observe(duration)
	TenantIDExtractionTotal.WithLabelValues(source).Inc()
}

// RecordTenantValidation 记录租户验证
func RecordTenantValidation(result string, duration float64, cacheHit bool) {
	TenantValidationTotal.WithLabelValues(result).Inc()
	TenantValidationDuration.WithLabelValues(formatCacheHit(cacheHit)).Observe(duration)
}

// RecordErrorCode 记录错误码
func RecordErrorCode(module, codeRange, errorName string) {
	ErrorCodeTotal.WithLabelValues(module, codeRange, errorName).Inc()
	ErrorDistribution.WithLabelValues(module, codeRange, errorName).Inc()
}

// RecordSessionCreation 记录Session创建
func RecordSessionCreation(result string, duration float64) {
	SessionCreationTotal.WithLabelValues(result).Inc()
	SessionCreationDuration.WithLabelValues(result).Observe(duration)
}

// RecordSessionValidation 记录Session验证
func RecordSessionValidation(result string, duration float64) {
	SessionValidationTotal.WithLabelValues(result).Inc()
	SessionValidationDuration.WithLabelValues(result).Observe(duration)
}

// RecordBotOperation 记录Bot操作
func RecordBotOperation(operation, result string, duration float64) {
	BotOperationTotal.WithLabelValues(operation, result).Inc()
	BotOperationDuration.WithLabelValues(operation).Observe(duration)
}

// RecordBotExecution 记录Bot执行
func RecordBotExecution(result string, duration float64) {
	BotExecutionTotal.WithLabelValues(result).Inc()
	BotExecutionDuration.WithLabelValues(result).Observe(duration)
}

// RecordAgentRun 记录Agent运行
func RecordAgentRun(result string, duration float64) {
	AgentRunTotal.WithLabelValues(result).Inc()
	AgentRunDuration.WithLabelValues(result).Observe(duration)
}

// RecordWorkflowOperation 记录Workflow操作
func RecordWorkflowOperation(operation, result string, duration float64) {
	WorkflowOperationTotal.WithLabelValues(operation, result).Inc()
	WorkflowOperationDuration.WithLabelValues(operation).Observe(duration)
}

// RecordWorkflowNodeExecution 记录Workflow节点执行
func RecordWorkflowNodeExecution(nodeType, result string, duration float64) {
	WorkflowNodeExecutionTotal.WithLabelValues(nodeType, result).Inc()
	WorkflowNodeExecutionDuration.WithLabelValues(nodeType).Observe(duration)
}

// RecordContextCacheOperation 记录Context缓存操作
func RecordContextCacheOperation(operation string, duration float64) {
	ContextCacheOperationTotal.WithLabelValues(operation, "success").Inc()
	ContextCacheOperationDuration.WithLabelValues(operation).Observe(duration)
}

// formatCacheHit 格式化缓存命中状态
func formatCacheHit(hit bool) string {
	if hit {
		return "true"
	}
	return "false"
}
