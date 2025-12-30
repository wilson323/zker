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

// ========== Bot相关指标 ==========

var (
	// BotTotal Bot总数(Gauge)
	BotTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_total",
			Help: "Total number of bots",
		},
		[]string{"tenant_id", "status"}, // status: active, inactive, deleted
	)

	// BotInvocationTotal Bot调用总数(Counter)
	BotInvocationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_invocation_total",
			Help: "Total number of bot invocations",
		},
		[]string{"tenant_id", "bot_id", "status"},
	)

	// BotInvocationDuration Bot调用延迟(Histogram)
	BotInvocationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bot_invocation_duration_seconds",
			Help:    "Bot invocation latency in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
		[]string{"tenant_id", "bot_id"},
	)

	// BotErrorRate Bot错误率(Gauge)
	BotErrorRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_error_rate",
			Help: "Bot error rate (percentage)",
		},
		[]string{"tenant_id", "bot_id", "error_type"},
	)

	// BotTokenUsage Bot Token使用量(Counter)
	BotTokenUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bot_token_usage_total",
			Help: "Total token usage by bots",
		},
		[]string{"tenant_id", "bot_id", "model", "token_type"}, // token_type: prompt, completion
	)

	// BotCost Bot成本(Gauge)
	BotCost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bot_cost_total",
			Help: "Total bot cost in cents",
		},
		[]string{"tenant_id", "bot_id", "currency"},
	)
)

// ========== 工作流相关指标 ==========

var (
	// WorkflowTotal 工作流总数(Gauge)
	WorkflowTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "workflow_total",
			Help: "Total number of workflows",
		},
		[]string{"tenant_id", "status"},
	)

	// WorkflowExecutionTotal 工作流执行总数(Counter)
	WorkflowExecutionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_execution_total",
			Help: "Total number of workflow executions",
		},
		[]string{"tenant_id", "workflow_id", "status"}, // status: success, failed, cancelled
	)

	// WorkflowExecutionDuration 工作流执行延迟(Histogram)
	WorkflowExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "workflow_execution_duration_seconds",
			Help:    "Workflow execution latency in seconds",
			Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60, 300, 600},
		},
		[]string{"tenant_id", "workflow_id"},
	)

	// WorkflowNodeExecutionTotal 工作流节点执行总数(Counter)
	WorkflowNodeExecutionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_node_execution_total",
			Help: "Total number of workflow node executions",
		},
		[]string{"tenant_id", "workflow_id", "node_type", "status"},
	)

	// WorkflowStepDuration 工作流步骤延迟(Histogram)
	WorkflowStepDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "workflow_step_duration_seconds",
			Help:    "Workflow step latency in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"tenant_id", "workflow_id", "node_type"},
	)
)

// ========== 对话相关指标 ==========

var (
	// ConversationTotal 对话总数(Gauge)
	ConversationTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "conversation_total",
			Help: "Total number of conversations",
		},
		[]string{"tenant_id", "bot_id"},
	)

	// MessageTotal 消息总数(Gauge)
	MessageTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "message_total",
			Help: "Total number of messages",
		},
		[]string{"tenant_id", "bot_id", "message_type"}, // message_type: user, assistant, system
	)

	// MessageLatency 消息响应延迟(Histogram)
	MessageLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "message_latency_seconds",
			Help:    "Message response latency in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"tenant_id", "bot_id", "message_type"},
	)

	// ConversationLength 对话长度分布(Histogram)
	ConversationLength = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "conversation_length_messages",
			Help:    "Conversation length in number of messages",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100},
		},
		[]string{"tenant_id", "bot_id"},
	)

)

// ========== 知识库相关指标 ==========

var (
	// KnowledgeBaseTotal 知识库总数(Gauge)
	KnowledgeBaseTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "knowledge_base_total",
			Help: "Total number of knowledge bases",
		},
		[]string{"tenant_id"},
	)

	// KnowledgeBaseDocumentTotal 知识库文档总数(Gauge)
	KnowledgeBaseDocumentTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "knowledge_base_document_total",
			Help: "Total number of documents in knowledge base",
		},
		[]string{"tenant_id", "knowledge_base_id"},
	)

	// KnowledgeBaseSearchTotal 知识库搜索总数(Counter)
	KnowledgeBaseSearchTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "knowledge_base_search_total",
			Help: "Total number of knowledge base searches",
		},
		[]string{"tenant_id", "knowledge_base_id", "status"},
	)

	// KnowledgeBaseSearchDuration 知识库搜索延迟(Histogram)
	KnowledgeBaseSearchDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "knowledge_base_search_duration_seconds",
			Help:    "Knowledge base search latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "knowledge_base_id"},
	)

	// KnowledgeBaseIndexSize 知识库索引大小(Gauge)
	KnowledgeBaseIndexSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "knowledge_base_index_size_bytes",
			Help: "Knowledge base index size in bytes",
		},
		[]string{"tenant_id", "knowledge_base_id"},
	)
)

// ========== 权限相关指标 ==========

var (
	// PermissionCheckTotal 权限检查总数(Counter)
	PermissionCheckTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "permission_check_total",
			Help: "Total number of permission checks",
		},
		[]string{"tenant_id", "permission_type", "result"}, // permission_type: data, field, role; result: allowed, denied
	)

	// PermissionCheckDuration 权限检查延迟(Histogram)
	PermissionCheckDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "permission_check_duration_seconds",
			Help:    "Permission check latency in seconds",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05},
		},
		[]string{"tenant_id", "permission_type"},
	)

	// RoleAssignmentTotal 角色分配总数(Counter)
	RoleAssignmentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "role_assignment_total",
			Help: "Total number of role assignments",
		},
		[]string{"tenant_id", "role_code", "action"}, // action: assign, revoke
	)

	// UserRoleTotal 用户角色数(Gauge)
	UserRoleTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_role_total",
			Help: "Total number of user roles",
		},
		[]string{"tenant_id", "role_code"},
	)
)

// ========== 存储相关指标 ==========

var (
	// StorageUsageTotal 存储使用量(Gauge)
	StorageUsageTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "storage_usage_bytes",
			Help: "Total storage usage in bytes",
		},
		[]string{"tenant_id", "storage_type"}, // storage_type: database, object_store, cache
	)

	// StorageQuota 存储配额(Gauge)
	StorageQuota = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "storage_quota_bytes",
			Help: "Storage quota in bytes",
		},
		[]string{"tenant_id", "storage_type"},
	)

	// StorageUsageRate 存储使用率(Gauge)
	StorageUsageRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "storage_usage_rate",
			Help: "Storage usage rate (percentage)",
		},
		[]string{"tenant_id", "storage_type"},
	)

	// ObjectUploadTotal 对象上传总数(Counter)
	ObjectUploadTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "object_upload_total",
			Help: "Total number of object uploads",
		},
		[]string{"tenant_id", "object_type", "status"},
	)

	// ObjectDownloadTotal 对象下载总数(Counter)
	ObjectDownloadTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "object_download_total",
			Help: "Total number of object downloads",
		},
		[]string{"tenant_id", "object_type", "status"},
	)

	// ObjectTransferDuration 对象传输延迟(Histogram)
	ObjectTransferDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "object_transfer_duration_seconds",
			Help:    "Object transfer latency in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
		[]string{"tenant_id", "object_type", "operation"}, // operation: upload, download
	)
)

// ========== 用户活跃度指标 ==========

var (
	// ActiveUserTotal 活跃用户数(Gauge)
	ActiveUserTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_active_total",
			Help: "Total number of active users",
		},
		[]string{"tenant_id", "period"}, // period: daily, weekly, monthly
	)

	// UserSessionTotal 用户会话数(Gauge)
	UserSessionTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "user_session_total",
			Help: "Total number of active user sessions",
		},
		[]string{"tenant_id"},
	)

	// UserLoginTotal 用户登录总数(Counter)
	UserLoginTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_login_total",
			Help: "Total number of user logins",
		},
		[]string{"tenant_id", "status"}, // status: success, failed
	)

	// UserOperationTotal 用户操作总数(Counter)
	UserOperationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_operation_total",
			Help: "Total number of user operations",
		},
		[]string{"tenant_id", "operation_type"}, // operation_type: create, update, delete, view
	)
)

// ========== 模型调用指标 ==========

var (
	// ModelInvocationTotal 模型调用总数(Counter)
	ModelInvocationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "model_invocation_total",
			Help: "Total number of model invocations",
		},
		[]string{"tenant_id", "model", "model_type", "status"}, // model_type: chat, completion, embedding
	)

	// ModelTokenUsageTotal 模型Token使用量(Counter)
	ModelTokenUsageTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "model_token_usage_total",
			Help: "Total token usage by models",
		},
		[]string{"tenant_id", "model", "token_type"}, // token_type: prompt, completion
	)

	// ModelInvocationDuration 模型调用延迟(Histogram)
	ModelInvocationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "model_invocation_duration_seconds",
			Help:    "Model invocation latency in seconds",
			Buckets: []float64{0.5, 1, 2, 5, 10, 20, 30, 60},
		},
		[]string{"tenant_id", "model"},
	)

	// ModelCostTotal 模型成本(Gauge)
	ModelCostTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "model_cost_total",
			Help: "Total model cost in cents",
		},
		[]string{"tenant_id", "model", "currency"},
	)

	// ModelErrorRate 模型错误率(Gauge)
	ModelErrorRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "model_error_rate",
			Help: "Model error rate (percentage)",
		},
		[]string{"tenant_id", "model", "error_type"},
	)
)

// ========== API限流指标 ==========

var (
	// RateLimitExceededTotal 限流触发总数(Counter)
	RateLimitExceededTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded events",
		},
		[]string{"tenant_id", "user_id", "limit_type"},
	)

	// RateLimitRemaining 剩余请求数(Gauge)
	RateLimitRemaining = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rate_limit_remaining",
			Help: "Remaining number of requests in rate limit window",
		},
		[]string{"tenant_id", "user_id", "limit_type"},
	)

	// ThrottleRequestTotal 节流请求数(Counter)
	ThrottleRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "throttle_request_total",
			Help: "Total number of throttled requests",
		},
		[]string{"tenant_id", "user_id", "throttle_reason"},
	)
)

// ========== 业务健康度指标 ==========

var (
	// BusinessHealthScore 业务健康度评分(Gauge)
	BusinessHealthScore = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "business_health_score",
			Help: "Business health score (0-100)",
		},
		[]string{"tenant_id", "health_dimension"}, // health_dimension: availability, performance, quality
	)

	// SLAComplianceRate SLA合规率(Gauge)
	SLAComplianceRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sla_compliance_rate",
			Help: "SLA compliance rate (percentage)",
		},
		[]string{"tenant_id", "sla_type"}, // sla_type: uptime, response_time, error_rate
	)

	// IncidentTotal 事故总数(Counter)
	IncidentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "incident_total",
			Help: "Total number of incidents",
		},
		[]string{"tenant_id", "severity", "incident_type"}, // severity: critical, high, medium, low
	)
)

// ========== 辅助函数 ==========

// RecordBotInvocation 记录Bot调用
func RecordBotInvocation(tenantID, botID string, duration float64, status string, tokenUsage map[string]int64, cost float64) {
	BotInvocationTotal.WithLabelValues(tenantID, botID, status).Inc()
	BotInvocationDuration.WithLabelValues(tenantID, botID).Observe(duration)

	if tokenUsage != nil {
		for tokenType, count := range tokenUsage {
			BotTokenUsage.WithLabelValues(tenantID, botID, "", tokenType).Add(float64(count))
		}
	}

	BotCost.WithLabelValues(tenantID, botID, "USD").Set(cost)
}

// RecordWorkflowExecution 记录工作流执行
func RecordWorkflowExecution(tenantID, workflowID string, duration float64, status string) {
	WorkflowExecutionTotal.WithLabelValues(tenantID, workflowID, status).Inc()
	WorkflowExecutionDuration.WithLabelValues(tenantID, workflowID).Observe(duration)
}

// RecordMessage 记录消息
func RecordMessage(tenantID, botID, messageType string, duration float64) {
	MessageTotal.WithLabelValues(tenantID, botID, messageType).Inc()
	MessageLatency.WithLabelValues(tenantID, botID, messageType).Observe(duration)
}

// RecordPermissionCheck 记录权限检查
func RecordPermissionCheck(tenantID, permissionType string, duration float64, result string) {
	PermissionCheckTotal.WithLabelValues(tenantID, permissionType, result).Inc()
	PermissionCheckDuration.WithLabelValues(tenantID, permissionType).Observe(duration)
}

// RecordModelInvocation 记录模型调用
func RecordModelInvocation(tenantID, model, modelType, status string, duration float64, promptTokens, completionTokens int64, cost float64) {
	ModelInvocationTotal.WithLabelValues(tenantID, model, modelType, status).Inc()
	ModelInvocationDuration.WithLabelValues(tenantID, model).Observe(duration)
	ModelTokenUsageTotal.WithLabelValues(tenantID, model, "prompt").Add(float64(promptTokens))
	ModelTokenUsageTotal.WithLabelValues(tenantID, model, "completion").Add(float64(completionTokens))
	ModelCostTotal.WithLabelValues(tenantID, model, "USD").Set(cost)
}

// UpdateHealthScore 更新健康度评分
func UpdateHealthScore(tenantID, dimension string, score float64) {
	BusinessHealthScore.WithLabelValues(tenantID, dimension).Set(score)
}

// RecordIncident 记录事故
func RecordIncident(tenantID, severity, incidentType string) {
	IncidentTotal.WithLabelValues(tenantID, severity, incidentType).Inc()
}
