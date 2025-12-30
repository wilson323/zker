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

package entity

import (
	"fmt"
	"time"
)

// AgentMetricType Agent指标类型
type AgentMetricType string

const (
	AgentMetricTypeQPS           AgentMetricType = "qps"             // Agent QPS
	AgentMetricTypeResponseTime  AgentMetricType = "response_time"   // Agent响应时间
	AgentMetricTypeErrorRate     AgentMetricType = "error_rate"      // Agent错误率
	AgentMetricTypeSatisfaction  AgentMetricType = "satisfaction"    // 用户满意度
	AgentMetricTypeTokenUsage    AgentMetricType = "token_usage"     // Token使用量
	AgentMetricTypeCost          AgentMetricType = "cost"            // 成本
	AgentMetricTypeConversation  AgentMetricType = "conversation"    // 对话数
	AgentMetricTypeActiveUsers   AgentMetricType = "active_users"    // 活跃用户数
)

// AgentMetrics Agent指标实体
type AgentMetrics struct {
	ID              string          `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string          `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1"`
	AgentID         string          `json:"agent_id" gorm:"type:varchar(36);not null;index:idx_agent_id,priority:1;index:idx_tenant_agent,priority:2"`
	AgentName       string          `json:"agent_name" gorm:"type:varchar(200);not null"`
	MetricType      AgentMetricType `json:"metric_type" gorm:"type:varchar(50);not null;index:idx_metric_type"`
	MetricValue     float64         `json:"metric_value" gorm:"type:decimal(15,4);not null"`
	MetricTimestamp time.Time       `json:"metric_timestamp" gorm:"not null;index:idx_timestamp"`
	Metadata        string          `json:"metadata" gorm:"type:json"` // 额外元数据
	CreatedAt       time.Time       `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (AgentMetrics) TableName() string {
	return "agent_metrics"
}

// NewAgentMetrics 创建Agent指标
func NewAgentMetrics(tenantID, agentID, agentName string, metricType AgentMetricType, value float64, metadata string) *AgentMetrics {
	return &AgentMetrics{
		ID:              generateID(),
		TenantID:        tenantID,
		AgentID:         agentID,
		AgentName:       agentName,
		MetricType:      metricType,
		MetricValue:     value,
		MetricTimestamp: time.Now(),
		Metadata:        metadata,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// IsExpired 检查指标是否过期
func (m *AgentMetrics) IsExpired(retentionDays int) bool {
	expirationTime := time.Now().AddDate(0, 0, -retentionDays)
	return m.MetricTimestamp.Before(expirationTime)
}

// ============================================================
// Agent指标聚合类型
// ============================================================

// AgentQPSMetrics Agent QPS指标聚合
type AgentQPSMetrics struct {
	AgentID      string  `json:"agent_id"`
	AgentName    string  `json:"agent_name"`
	Count        int64   `json:"count"`
	Average      float64 `json:"average"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	P50          float64 `json:"p50"`
	P95          float64 `json:"p95"`
	P99          float64 `json:"p99"`
	Trend        string  `json:"trend"`
	TrendChange  float64 `json:"trend_change"`
}

// AgentResponseTimeMetrics Agent响应时间指标聚合
type AgentResponseTimeMetrics struct {
	AgentID      string  `json:"agent_id"`
	AgentName    string  `json:"agent_name"`
	Average      float64 `json:"average"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	P50          float64 `json:"p50"`
	P95          float64 `json:"p95"`
	P99          float64 `json:"p99"`
	Trend        string  `json:"trend"`
	TrendChange  float64 `json:"trend_change"`
}

// AgentErrorRateMetrics Agent错误率指标聚合
type AgentErrorRateMetrics struct {
	AgentID         string            `json:"agent_id"`
	AgentName       string            `json:"agent_name"`
	TotalRequests   int64             `json:"total_requests"`
	ErrorRequests   int64             `json:"error_requests"`
	ErrorRate       float64           `json:"error_rate"`
	Average         float64           `json:"average"`
	Max             float64           `json:"max"`
	Trend           string            `json:"trend"`
	TrendChange     float64           `json:"trend_change"`
	ErrorBreakdown  map[string]int64  `json:"error_breakdown"`
}

// AgentSatisfactionMetrics Agent满意度指标聚合
type AgentSatisfactionMetrics struct {
	AgentID        string    `json:"agent_id"`
	AgentName      string    `json:"agent_name"`
	AverageScore   float64   `json:"average_score"`   // 平均满意度（0-5）
	MinScore       float64   `json:"min_score"`       // 最低评分
	MaxScore       float64   `json:"max_score"`       // 最高评分
	TotalRatings   int64     `json:"total_ratings"`   // 总评分数
	FiveStarCount  int64     `json:"five_star_count"` // 5星数量
	FourStarCount  int64     `json:"four_star_count"` // 4星数量
	ThreeStarCount int64     `json:"three_star_count"` // 3星数量
	TwoStarCount   int64     `json:"two_star_count"`   // 2星数量
	OneStarCount   int64     `json:"one_star_count"`   // 1星数量
	Trend          string    `json:"trend"`
	TrendChange    float64   `json:"trend_change"`
}

// AgentTokenUsageMetrics Agent Token使用指标
type AgentTokenUsageMetrics struct {
	AgentID          string  `json:"agent_id"`
	AgentName        string  `json:"agent_name"`
	TotalTokens      int64   `json:"total_tokens"`
	InputTokens      int64   `json:"input_tokens"`
	OutputTokens     int64   `json:"output_tokens"`
	AveragePerReq    float64 `json:"average_per_req"`
	EstimatedCost    float64 `json:"estimated_cost"`
	Trend            string  `json:"trend"`
	TrendChange      float64 `json:"trend_change"`
}

// AgentPerformanceReport Agent性能报告
type AgentPerformanceReport struct {
	AgentID            string                     `json:"agent_id"`
	AgentName          string                     `json:"agent_name"`
	ReportDate         string                     `json:"report_date"`         // YYYY-MM-DD
	TimeRange          TimeRange                  `json:"time_range"`
	QPSMetrics         *AgentQPSMetrics           `json:"qps_metrics"`
	ResponseTimeMetrics *AgentResponseTimeMetrics `json:"response_time_metrics"`
	ErrorRateMetrics   *AgentErrorRateMetrics     `json:"error_rate_metrics"`
	SatisfactionMetrics *AgentSatisfactionMetrics `json:"satisfaction_metrics"`
	TokenUsageMetrics  *AgentTokenUsageMetrics    `json:"token_usage_metrics"`
	OverallScore       float64                    `json:"overall_score"`       // 综合评分（0-100）
	PerformanceLevel   string                     `json:"performance_level"`   // 性能等级
	Recommendations    []string                   `json:"recommendations"`     // 优化建议
	ComparedWithPeriod string                     `json:"compared_with_period"` // 对比周期
	Improvements       []string                   `json:"improvements"`       // 改进项
	Concerns           []string                   `json:"concerns"`           // 关注项
	GeneratedAt        time.Time                  `json:"generated_at"`
}

// AgentComparison Agent对比数据
type AgentComparison struct {
	AgentID              string  `json:"agent_id"`
	AgentName            string  `json:"agent_name"`
	QPSRanking           int     `json:"qps_ranking"`
	ResponseTimeRanking  int     `json:"response_time_ranking"`
	ErrorRateRanking     int     `json:"error_rate_ranking"`
	SatisfactionRanking  int     `json:"satisfaction_ranking"`
	OverallRanking       int     `json:"overall_ranking"`
	QPSPercentile        float64 `json:"qps_percentile"`
	ResponseTimePercentile float64 `json:"response_time_percentile"`
	ErrorRatePercentile  float64 `json:"error_rate_percentile"`
	SatisfactionPercentile float64 `json:"satisfaction_percentile"`
}

// AgentHealthStatus Agent健康状态
type AgentHealthStatus struct {
	AgentID            string    `json:"agent_id"`
	AgentName          string    `json:"agent_name"`
	Status             string    `json:"status"`             // healthy/warning/critical/offline
	HealthScore        float64   `json:"health_score"`        // 健康评分（0-100）
	QPSStatus          string    `json:"qps_status"`          // normal/low/high
	ResponseTimeStatus string    `json:"response_time_status"` // normal/slow/very_slow
	ErrorRateStatus    string    `json:"error_rate_status"`   // normal/high/critical
	LastActiveTime     time.Time `json:"last_active_time"`
	Issues             []string  `json:"issues"`              // 发现的问题
	AlertCount         int64     `json:"alert_count"`         // 告警数量
}

// GetPerformanceLevel 获取性能等级
func GetPerformanceLevel(score float64) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 75:
		return "good"
	case score >= 60:
		return "fair"
	case score >= 40:
		return "poor"
	default:
		return "critical"
	}
}

// GetHealthLevel 获取健康等级
func GetHealthLevel(score float64) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 75:
		return "good"
	case score >= 60:
		return "warning"
	default:
		return "critical"
	}
}

// CalculateTrend 计算趋势
func CalculateTrend(current, previous float64) (trend string, change float64) {
	if previous == 0 {
		return "stable", 0
	}

	change = ((current - previous) / previous) * 100

	switch {
	case change > 10:
		return "up", change
	case change < -10:
		return "down", change
	default:
		return "stable", change
	}
}

// ValidateMetricType 验证Agent指标类型
func ValidateMetricType(metricType AgentMetricType) error {
	switch metricType {
	case AgentMetricTypeQPS, AgentMetricTypeResponseTime, AgentMetricTypeErrorRate,
		AgentMetricTypeSatisfaction, AgentMetricTypeTokenUsage, AgentMetricTypeCost,
		AgentMetricTypeConversation, AgentMetricTypeActiveUsers:
		return nil
	default:
		return fmt.Errorf("invalid agent metric type: %s", metricType)
	}
}
