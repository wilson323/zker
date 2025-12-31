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
	"time"
)

// TraceID 链路追踪ID类型
type TraceID string

// SpanID 跨度ID类型
type SpanID string

// Trace 链路追踪实体
type Trace struct {
	TraceID      TraceID              `json:"trace_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID     string               `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID       string               `json:"user_id" gorm:"type:varchar(36);index:idx_user_id"`
	AgentID      string               `json:"agent_id" gorm:"type:varchar(36);index:idx_agent_id"`
	WorkflowID   string               `json:"workflow_id,omitempty" gorm:"type:varchar(36);index:idx_workflow_id"`
	ConversationID string             `json:"conversation_id" gorm:"type:varchar(36);index:idx_conversation_id"`
	MessageID    string               `json:"message_id" gorm:"type:varchar(36)"`
	ParentSpanID SpanID               `json:"parent_span_id,omitempty" gorm:"type:varchar(36)"`
	RootSpanID   SpanID               `json:"root_span_id" gorm:"type:varchar(36)"`
	SpanName     string               `json:"span_name" gorm:"type:varchar(200);not null"`
	SpanKind     string               `json:"span_kind" gorm:"type:varchar(20);not null"` // internal, server, client, producer, consumer
	StartTime    time.Time            `json:"start_time" gorm:"not null;index:idx_start_time"`
	EndTime      time.Time            `json:"end_time"`
	Duration     int64                `json:"duration" gorm:"not null"` // 毫秒
	Status       string               `json:"status" gorm:"type:varchar(20);not null"` // success, error, timeout
	StatusCode   string               `json:"status_code,omitempty" gorm:"type:varchar(50)"`
	StatusMessage string              `json:"status_message,omitempty" gorm:"type:text"`
	Attributes   string               `json:"attributes" gorm:"type:json"` // JSON格式的属性
	Events       string               `json:"events" gorm:"type:json"`     // JSON格式的事件列表
	Links        string               `json:"links" gorm:"type:json"`     // JSON格式的关联链接
	CreatedAt    time.Time            `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (Trace) TableName() string {
	return "agent_traces"
}

// IsSuccess 是否成功
func (t *Trace) IsSuccess() bool {
	return t.Status == "success"
}

// GetDurationMs 获取持续时间(毫秒)
func (t *Trace) GetDurationMs() int64 {
	return t.Duration
}

// GetDurationSec 获取持续时间(秒)
func (t *Trace) GetDurationSec() float64 {
	return float64(t.Duration) / 1000.0
}

// TraceSpan 链路跨度
type TraceSpan struct {
	SpanID       SpanID              `json:"span_id"`
	ParentSpanID SpanID              `json:"parent_span_id"`
	SpanName     string              `json:"span_name"`
	StartTime    time.Time           `json:"start_time"`
	EndTime      time.Time           `json:"end_time"`
	Duration     time.Duration       `json:"duration"`
	Attributes   map[string]string   `json:"attributes"`
	Events       []TraceEvent        `json:"events"`
	Status       string              `json:"status"`
}

// TraceEvent 链路事件
type TraceEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Name      string                 `json:"name"`
	Attributes map[string]string     `json:"attributes"`
}

// QualityMetric 质量指标
type QualityMetric struct {
	MetricID    string    `json:"metric_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID    string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	AgentID     string    `json:"agent_id" gorm:"type:varchar(36);index:idx_agent_id"`
	ConversationID string `json:"conversation_id" gorm:"type:varchar(36);index:idx_conversation_id"`
	MessageID   string    `json:"message_id" gorm:"type:varchar(36);index:idx_message_id"`
	UserID      string    `json:"user_id" gorm:"type:varchar(36)"`
	ResponsivenessScore float64 `json:"responsiveness_score"` // 响应质量评分
	RelevanceScore      float64 `json:"relevance_score"`      // 相关性评分
	AccuracyScore       float64 `json:"accuracy_score"`       // 准确性评分
	CompletenessScore   float64 `json:"completeness_score"`   // 完整性评分
	OverallScore        float64 `json:"overall_score"`        // 综合评分
	UserRating          int     `json:"user_rating"`          // 用户评分(1-5)
	UserFeedback        string  `json:"user_feedback" gorm:"type:text"`
	EvaluatedAt         time.Time `json:"evaluated_at" gorm:"not null"`
	CreatedAt           time.Time `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (QualityMetric) TableName() string {
	return "agent_quality_metrics"
}

// CostMetric 成本指标
type CostMetric struct {
	MetricID      string    `json:"metric_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	AgentID       string    `json:"agent_id" gorm:"type:varchar(36);index:idx_agent_id"`
	UserID        string    `json:"user_id" gorm:"type:varchar(36);index:idx_user_id"`
	ConversationID string  `json:"conversation_id" gorm:"type:varchar(36)"`
	MessageID     string    `json:"message_id" gorm:"type:varchar(36)"`
	ModelName     string    `json:"model_name" gorm:"type:varchar(100);index:idx_model_name"`
	PromptTokens  int       `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens   int       `json:"total_tokens" gorm:"index:idx_total_tokens"`
	Cost          float64   `json:"cost" gorm:"type:decimal(10,4);index:idx_cost"`
	Currency      string    `json:"currency" gorm:"type:varchar(10);default:'USD'"`
	Timestamp     time.Time `json:"timestamp" gorm:"not null;index:idx_timestamp"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (CostMetric) TableName() string {
	return "agent_cost_metrics"
}

// GetTotalCostUSD 获取总成本(美元)
func (c *CostMetric) GetTotalCostUSD() float64 {
	if c.Currency == "USD" {
		return c.Cost
	}
	// 简化实现，实际应用中应该使用汇率转换
	return c.Cost
}

// PerformanceMetric 性能指标
type PerformanceMetric struct {
	MetricID      string    `json:"metric_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	AgentID       string    `json:"agent_id" gorm:"type:varchar(36);index:idx_agent_id"`
	ResponseTime  int64     `json:"response_time" gorm:"not null"` // 毫秒
	Throughput    float64   `json:"throughput"`                    // 每秒请求数
	ErrorRate     float64   `json:"error_rate"`                    // 错误率
	SuccessRate   float64   `json:"success_rate"`                  // 成功率
	Timestamp     time.Time `json:"timestamp" gorm:"not null;index:idx_timestamp"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (PerformanceMetric) TableName() string {
	return "agent_performance_metrics"
}

// TraceFilter 链路追踪过滤器
type TraceFilter struct {
	TenantID       string    `json:"tenant_id"`
	UserID         string    `json:"user_id,omitempty"`
	AgentID        string    `json:"agent_id,omitempty"`
	WorkflowID     string    `json:"workflow_id,omitempty"`
	ConversationID string    `json:"conversation_id,omitempty"`
	TraceID        string    `json:"trace_id,omitempty"`
	StartTime      time.Time `json:"start_time,omitempty"`
	EndTime        time.Time `json:"end_time,omitempty"`
	MinDuration    int64     `json:"min_duration,omitempty"` // 毫秒
	MaxDuration    int64     `json:"max_duration,omitempty"` // 毫秒
	Status         string    `json:"status,omitempty"`
	Limit          int       `json:"limit,omitempty"`
	Offset         int       `json:"offset,omitempty"`
}

// QualityFilter 质量指标过滤器
type QualityFilter struct {
	TenantID        string    `json:"tenant_id"`
	AgentID         string    `json:"agent_id,omitempty"`
	ConversationID  string    `json:"conversation_id,omitempty"`
	MinOverallScore float64   `json:"min_overall_score,omitempty"`
	MaxOverallScore float64   `json:"max_overall_score,omitempty"`
	StartTime       time.Time `json:"start_time,omitempty"`
	EndTime         time.Time `json:"end_time,omitempty"`
	Limit           int       `json:"limit,omitempty"`
	Offset          int       `json:"offset,omitempty"`
}

// CostFilter 成本指标过滤器
type CostFilter struct {
	TenantID        string    `json:"tenant_id"`
	AgentID         string    `json:"agent_id,omitempty"`
	UserID          string    `json:"user_id,omitempty"`
	ModelName       string    `json:"model_name,omitempty"`
	MinCost         float64   `json:"min_cost,omitempty"`
	MaxCost         float64   `json:"max_cost,omitempty"`
	StartTime       time.Time `json:"start_time,omitempty"`
	EndTime         time.Time `json:"end_time,omitempty"`
	Limit           int       `json:"limit,omitempty"`
	Offset          int       `json:"offset,omitempty"`
}

// CostAnalysis 成本分析结果
type CostAnalysis struct {
	TotalCost        float64            `json:"total_cost"`
	TotalTokens      int                `json:"total_tokens"`
	AvgCostPerToken  float64            `json:"avg_cost_per_token"`
	ModelCosts       map[string]float64 `json:"model_costs"`       // 各模型成本
	AgentCosts       map[string]float64 `json:"agent_costs"`       // 各Agent成本
	UserCosts        map[string]float64 `json:"user_costs"`        // 各用户成本
	DailyCosts       []DailyCost        `json:"daily_costs"`       // 每日成本
	CostTrend        []TrendPoint       `json:"cost_trend"`        // 成本趋势
	OptimizationTips []string           `json:"optimization_tips"` // 优化建议
}

// DailyCost 每日成本
type DailyCost struct {
	Date  string  `json:"date"`
	Cost  float64 `json:"cost"`
	Tokens int    `json:"tokens"`
}

// TrendPoint 趋势点
type TrendPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// PerformanceAnalysis 性能分析结果
type PerformanceAnalysis struct {
	AvgResponseTime  float64         `json:"avg_response_time"`  // 平均响应时间
	P50ResponseTime  float64         `json:"p50_response_time"`
	P95ResponseTime  float64         `json:"p95_response_time"`
	P99ResponseTime  float64         `json:"p99_response_time"`
	Throughput       float64         `json:"throughput"`
	ErrorRate        float64         `json:"error_rate"`
	SuccessRate      float64         `json:"success_rate"`
	PerformanceTrend []TrendPoint    `json:"performance_trend"`
	Bottlenecks      []string        `json:"bottlenecks"`
}

// EvaluationReport 评估报告
type EvaluationReport struct {
	ReportID       string            `json:"report_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID       string            `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	ReportType     string            `json:"report_type" gorm:"type:varchar(50);not null"` // quality, cost, performance, comprehensive
	StartTime      time.Time         `json:"start_time" gorm:"not null"`
	EndTime        time.Time         `json:"end_time" gorm:"not null"`
	QualitySummary *QualityAnalysis  `json:"quality_summary,omitempty" gorm:"type:json"`
	CostSummary    *CostAnalysis     `json:"cost_summary,omitempty" gorm:"type:json"`
	PerformanceSummary *PerformanceAnalysis `json:"performance_summary,omitempty" gorm:"type:json"`
	Recommendations []string         `json:"recommendations" gorm:"type:json"`
	CreatedAt      time.Time         `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (EvaluationReport) TableName() string {
	return "agent_evaluation_reports"
}

// QualityAnalysis 质量分析结果
type QualityAnalysis struct {
	AvgOverallScore      float64            `json:"avg_overall_score"`
	AvgResponsiveness    float64            `json:"avg_responsiveness"`
	AvgRelevance         float64            `json:"avg_relevance"`
	AvgAccuracy          float64            `json:"avg_accuracy"`
	AvgCompleteness      float64            `json:"avg_completeness"`
	ScoreDistribution    map[string]int     `json:"score_distribution"`    // 分数分布
	UserSatisfactionRate float64            `json:"user_satisfaction_rate"` // 用户满意度
	QualityTrend         []TrendPoint       `json:"quality_trend"`
}
