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

	"github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
)

// TenantMetricsRepository 租户指标仓储接口
type TenantMetricsRepository interface {
	// Create 创建指标记录
	Create(ctx context.Context, metrics *entity.TenantMetrics) error

	// CreateBatch 批量创建指标记录
	CreateBatch(ctx context.Context, metricsList []*entity.TenantMetrics) error

	// GetByID 根据ID获取指标
	GetByID(ctx context.Context, id string) (*entity.TenantMetrics, error)

	// GetByTenantAndType 根据租户ID和指标类型查询最新指标
	GetByTenantAndType(ctx context.Context, tenantID string, metricType entity.MetricType) (*entity.TenantMetrics, error)

	// GetByTenantAndTimeRange 根据租户ID和时间范围查询指标
	GetByTenantAndTimeRange(
		ctx context.Context,
		tenantID string,
		metricType entity.MetricType,
		startTime, endTime time.Time,
	) ([]*entity.TenantMetrics, error)

	// GetLatestByTenant 获取租户最新指标（所有类型）
	GetLatestByTenant(ctx context.Context, tenantID string) ([]*entity.TenantMetrics, error)

	// AggregateByTimeRange 按时间范围聚合指标
	AggregateByTimeRange(
		ctx context.Context,
		tenantID string,
		metricType entity.MetricType,
		startTime, endTime time.Time,
		interval string, // 1m, 5m, 1h, 1d
	) ([]*entity.TenantMetrics, error)

	// CalculateQPSMetrics 计算QPS指标聚合
	CalculateQPSMetrics(
		ctx context.Context,
		tenantID string,
		startTime, endTime time.Time,
	) (*entity.QPSMetrics, error)

	// CalculateResponseTimeMetrics 计算响应时间指标聚合
	CalculateResponseTimeMetrics(
		ctx context.Context,
		tenantID string,
		startTime, endTime time.Time,
	) (*entity.ResponseTimeMetrics, error)

	// CalculateErrorRateMetrics 计算错误率指标聚合
	CalculateErrorRateMetrics(
		ctx context.Context,
		tenantID string,
		startTime, endTime time.Time,
	) (*entity.ErrorRateMetrics, error)

	// CalculateConcurrencyMetrics 计算并发指标聚合
	CalculateConcurrencyMetrics(
		ctx context.Context,
		tenantID string,
		startTime, endTime time.Time,
	) (*entity.ConcurrencyMetrics, error)

	// GetTopTenantsByQPS 获取QPS最高的N个租户
	GetTopTenantsByQPS(
		ctx context.Context,
		limit int,
		startTime, endTime time.Time,
	) ([]*entity.TenantMetrics, error)

	// DeleteExpiredMetrics 删除过期指标（数据保留策略）
	DeleteExpiredMetrics(ctx context.Context, retentionDays int) (int64, error)

	// CountByTenant 统计租户指标数量
	CountByTenant(ctx context.Context, tenantID string) (int64, error)

	// GetMetricTypes 获取租户的所有指标类型
	GetMetricTypes(ctx context.Context, tenantID string) ([]entity.MetricType, error)

	// GetTimeSeriesData 获取时序数据（用于图表展示）
	GetTimeSeriesData(
		ctx context.Context,
		tenantID string,
		metricType entity.MetricType,
		startTime, endTime time.Time,
		interval time.Duration,
	) ([]TimeSeriesPoint, error)

	// BatchQueryLatestMetrics 批量查询多个租户的最新指标
	BatchQueryLatestMetrics(
		ctx context.Context,
		tenantIDs []string,
		metricType entity.MetricType,
	) (map[string]*entity.TenantMetrics, error)
}

// TimeSeriesPoint 时序数据点
type TimeSeriesPoint struct {
	Timestamp time.Time
	Value     float64
	Tags      entity.MetricTags
}

// AgentMetricsRepository Agent指标仓储接口
type AgentMetricsRepository interface {
	// Create 创建Agent指标记录
	Create(ctx context.Context, metrics *entity.AgentMetrics) error

	// CreateBatch 批量创建Agent指标记录
	CreateBatch(ctx context.Context, metricsList []*entity.AgentMetrics) error

	// GetByID 根据ID获取指标
	GetByID(ctx context.Context, id string) (*entity.AgentMetrics, error)

	// GetByAgentAndType 根据Agent ID和指标类型查询最新指标
	GetByAgentAndType(
		ctx context.Context,
		agentID string,
		metricType entity.AgentMetricType,
	) (*entity.AgentMetrics, error)

	// GetByTenantAndAgent 根据租户ID和Agent ID查询指标
	GetByTenantAndAgent(
		ctx context.Context,
		tenantID, agentID string,
		metricType entity.AgentMetricType,
		startTime, endTime time.Time,
	) ([]*entity.AgentMetrics, error)

	// GetByTenant 根据租户ID查询所有Agent指标
	GetByTenant(
		ctx context.Context,
		tenantID string,
		metricType entity.AgentMetricType,
		startTime, endTime time.Time,
	) ([]*entity.AgentMetrics, error)

	// CalculateAgentQPSMetrics 计算Agent QPS指标聚合
	CalculateAgentQPSMetrics(
		ctx context.Context,
		agentID string,
		startTime, endTime time.Time,
	) (*entity.AgentQPSMetrics, error)

	// CalculateAgentResponseTimeMetrics 计算Agent响应时间指标聚合
	CalculateAgentResponseTimeMetrics(
		ctx context.Context,
		agentID string,
		startTime, endTime time.Time,
	) (*entity.AgentResponseTimeMetrics, error)

	// CalculateAgentErrorRateMetrics 计算Agent错误率指标聚合
	CalculateAgentErrorRateMetrics(
		ctx context.Context,
		agentID string,
		startTime, endTime time.Time,
	) (*entity.AgentErrorRateMetrics, error)

	// CalculateAgentSatisfactionMetrics 计算Agent满意度指标聚合
	CalculateAgentSatisfactionMetrics(
		ctx context.Context,
		agentID string,
		startTime, endTime time.Time,
	) (*entity.AgentSatisfactionMetrics, error)

	// CalculateAgentTokenUsageMetrics 计算Agent Token使用指标
	CalculateAgentTokenUsageMetrics(
		ctx context.Context,
		agentID string,
		startTime, endTime time.Time,
	) (*entity.AgentTokenUsageMetrics, error)

	// GetTopAgentsByQPS 获取QPS最高的N个Agent
	GetTopAgentsByQPS(
		ctx context.Context,
		tenantID string,
		limit int,
		startTime, endTime time.Time,
	) ([]*entity.AgentMetrics, error)

	// GetAgentPerformanceReport 生成Agent性能报告
	GetAgentPerformanceReport(
		ctx context.Context,
		agentID string,
		startTime, endTime time.Time,
	) (*entity.AgentPerformanceReport, error)

	// CompareAgents 对比多个Agent的性能
	CompareAgents(
		ctx context.Context,
		tenantID string,
		agentIDs []string,
		startTime, endTime time.Time,
	) ([]*entity.AgentComparison, error)

	// DeleteExpiredMetrics 删除过期指标
	DeleteExpiredMetrics(ctx context.Context, retentionDays int) (int64, error)

	// GetActiveAgents 获取活跃Agent列表
	GetActiveAgents(
		ctx context.Context,
		tenantID string,
		startTime, endTime time.Time,
	) ([]string, error)

	// GetAgentHealthStatus 获取Agent健康状态
	GetAgentHealthStatus(
		ctx context.Context,
		tenantID string,
		agentID string,
	) (*entity.AgentHealthStatus, error)
}

// AlertRuleRepository 告警规则仓储接口
type AlertRuleRepository interface {
	// Create 创建告警规则
	Create(ctx context.Context, rule *entity.AlertRule) error

	// GetByID 根据ID获取规则
	GetByID(ctx context.Context, id string) (*entity.AlertRule, error)

	// GetByTenant 根据租户ID获取所有规则
	GetByTenant(ctx context.Context, tenantID string) ([]*entity.AlertRule, error)

	// GetEnabledRules 获取所有启用的规则
	GetEnabledRules(ctx context.Context) ([]*entity.AlertRule, error)

	// GetEnabledRulesByTenant 获取租户的所有启用规则
	GetEnabledRulesByTenant(ctx context.Context, tenantID string) ([]*entity.AlertRule, error)

	// Update 更新规则
	Update(ctx context.Context, rule *entity.AlertRule) error

	// Delete 删除规则（软删除）
	Delete(ctx context.Context, id string) error

	// EnableRule 启用规则
	EnableRule(ctx context.Context, id string) error

	// DisableRule 禁用规则
	DisableRule(ctx context.Context, id string) error

	// List 分页查询规则列表
	List(ctx context.Context, filter *AlertRuleFilter) ([]*entity.AlertRule, int64, error)
}

// AlertRuleFilter 告警规则查询过滤器
type AlertRuleFilter struct {
	TenantID   string
	MetricType string
	Severity   entity.AlertSeverity
	IsEnabled  *bool
	PageToken  string
	PageSize   int
}

// AlertHistoryRepository 告警历史仓储接口
type AlertHistoryRepository interface {
	// Create 创建告警历史
	Create(ctx context.Context, history *entity.AlertHistory) error

	// CreateBatch 批量创建告警历史
	CreateBatch(ctx context.Context, historyList []*entity.AlertHistory) error

	// GetByID 根据ID获取告警历史
	GetByID(ctx context.Context, id string) (*entity.AlertHistory, error)

	// GetByTenant 根据租户ID查询告警历史
	GetByTenant(ctx context.Context, filter *entity.AlertFilter) ([]*entity.AlertHistory, int64, error)

	// GetByAlertRule 根据告警规则ID查询历史
	GetByAlertRule(
		ctx context.Context,
		alertRuleID string,
		limit int,
	) ([]*entity.AlertHistory, error)

	// GetPendingAlerts 获取待处理的告警
	GetPendingAlerts(ctx context.Context) ([]*entity.AlertHistory, error)

	// GetPendingAlertsByTenant 获取租户的待处理告警
	GetPendingAlertsByTenant(ctx context.Context, tenantID string) ([]*entity.AlertHistory, error)

	// UpdateStatus 更新告警状态
	UpdateStatus(ctx context.Context, id string, status entity.AlertStatus) error

	// Acknowledge 确认告警
	Acknowledge(ctx context.Context, id, userID string) error

	// Resolve 解决告警
	Resolve(ctx context.Context, id, userID, note string) error

	// Silence 静默告警
	Silence(ctx context.Context, id string, duration time.Duration) error

	// GetAlertStatistics 获取告警统计信息
	GetAlertStatistics(
		ctx context.Context,
		tenantID string,
		startTime, endTime time.Time,
	) (*AlertStatistics, error)

	// GetTopAlerts 获取最频繁的告警
	GetTopAlerts(
		ctx context.Context,
		tenantID string,
		limit int,
		startTime, endTime time.Time,
	) ([]*AlertFrequency, error)
}

// AlertStatistics 告警统计信息
type AlertStatistics struct {
	TotalCount    int64                       `json:"total_count"`
	PendingCount  int64                       `json:"pending_count"`
	AcknowledgedCount int64                   `json:"acknowledged_count"`
	ResolvedCount int64                       `json:"resolved_count"`
	SeverityDistribution map[entity.AlertSeverity]int64 `json:"severity_distribution"`
	TrendData     []TrendDataPoint            `json:"trend_data"`
}

// TrendDataPoint 趋势数据点
type TrendDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
	Severity  entity.AlertSeverity `json:"severity"`
}

// AlertFrequency 告警频率
type AlertFrequency struct {
	AlertRuleID   string            `json:"alert_rule_id"`
	AlertRuleName string            `json:"alert_rule_name"`
	Count         int64             `json:"count"`
	Severity      entity.AlertSeverity `json:"severity"`
	LastOccurred  time.Time         `json:"last_occurred"`
}

// PerformanceReportRepository 性能报告仓储接口
type PerformanceReportRepository interface {
	// Create 创建性能报告
	Create(ctx context.Context, report *entity.PerformanceReport) error

	// GetByID 根据ID获取报告
	GetByID(ctx context.Context, id string) (*entity.PerformanceReport, error)

	// GetByTenantAndDate 根据租户ID和日期获取报告
	GetByTenantAndDate(
		ctx context.Context,
		tenantID, reportType, reportDate string,
	) (*entity.PerformanceReport, error)

	// GetByTenant 根据租户ID查询报告列表
	GetByTenant(ctx context.Context, filter *entity.ReportFilter) ([]*entity.PerformanceReport, int64, error)

	// GetLatestReport 获取最新报告
	GetLatestReport(
		ctx context.Context,
		tenantID string,
		reportType string,
	) (*entity.PerformanceReport, error)

	// Update 更新报告
	Update(ctx context.Context, report *entity.PerformanceReport) error

	// Delete 删除报告
	Delete(ctx context.Context, id string) error

	// List 分页查询报告列表
	List(ctx context.Context, filter *entity.ReportFilter) ([]*entity.PerformanceReport, int64, error)
}
