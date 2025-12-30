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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// AlertSeverity 告警严重级别
type AlertSeverity string

const (
	AlertSeverityWarning   AlertSeverity = "warning"   // 警告
	AlertSeverityCritical  AlertSeverity = "critical"  // 严重
	AlertSeverityEmergency AlertSeverity = "emergency" // 紧急
)

// AlertStatus 告警状态
type AlertStatus string

const (
	AlertStatusPending       AlertStatus = "pending"        // 待处理
	AlertStatusAcknowledged  AlertStatus = "acknowledged"   // 已确认
	AlertStatusResolved      AlertStatus = "resolved"       // 已解决
	AlertStatusSilenced      AlertStatus = "silenced"       // 已静默
)

// AlertComparisonOperator 告警比较操作符
type AlertComparisonOperator string

const (
	AlertComparisonGT AlertComparisonOperator = "gt" // 大于
	AlertComparisonLT AlertComparisonOperator = "lt" // 小于
	AlertComparisonEQ AlertComparisonOperator = "eq" // 等于
	AlertComparisonGTE AlertComparisonOperator = "gte" // 大于等于
	AlertComparisonLTE AlertComparisonOperator = "lte" // 小于等于
)

// NotificationChannel 通知渠道
type NotificationChannel string

const (
	NotificationChannelEmail   NotificationChannel = "email"
	NotificationChannelSMS     NotificationChannel = "sms"
	NotificationChannelWebhook NotificationChannel = "webhook"
	NotificationChannelSlack   NotificationChannel = "slack"
	NotificationChannelDingTalk NotificationChannel = "dingtalk"
	NotificationChannelWeChat  NotificationChannel = "wechat"
)

// NotificationChannels 通知渠道列表（JSON格式）
type NotificationChannels []NotificationChannel

// Scan 实现 sql.Scanner 接口
func (nc *NotificationChannels) Scan(value interface{}) error {
	if value == nil {
		*nc = make(NotificationChannels, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal NotificationChannels value: %v", value)
	}

	return json.Unmarshal(bytes, nc)
}

// Value 实现 driver.Valuer 接口
func (nc NotificationChannels) Value() (driver.Value, error) {
	if nc == nil {
		return "[]", nil
	}
	return json.Marshal(nc)
}

// Contains 检查是否包含指定渠道
func (nc NotificationChannels) Contains(channel NotificationChannel) bool {
	for _, c := range nc {
		if c == channel {
			return true
		}
	}
	return false
}

// ============================================================
// AlertRule 告警规则实体
// ============================================================

// AlertRule 告警规则
type AlertRule struct {
	ID                   string                   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID             string                   `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	RuleName             string                   `json:"rule_name" gorm:"type:varchar(100);not null"`
	Description          string                   `json:"description" gorm:"type:text"`
	MetricType           string                   `json:"metric_type" gorm:"type:varchar(50);not null"` // qps, response_time, error_rate
	Threshold            float64                  `json:"threshold" gorm:"type:decimal(15,4);not null"`
	Comparison           AlertComparisonOperator  `json:"comparison" gorm:"type:varchar(10);not null"`
	Severity             AlertSeverity            `json:"severity" gorm:"type:varchar(20);not null"`
	NotificationChannels NotificationChannels     `json:"notification_channels" gorm:"type:json"`
	NotificationConfig   string                   `json:"notification_config" gorm:"type:json"` // 通知配置
	EvaluationInterval   int                      `json:"evaluation_interval"`  // 评估间隔（秒）
	ForDuration          int                      `json:"for_duration"`         // 持续时间（秒）
	IsEnabled            bool                     `json:"is_enabled" gorm:"default:true"`
	CreatedAt            time.Time                `json:"created_at" gorm:"not null"`
	UpdatedAt            time.Time                `json:"updated_at" gorm:"not null"`
	DeletedAt            *time.Time               `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (AlertRule) TableName() string {
	return "alert_rules"
}

// NewAlertRule 创建告警规则
func NewAlertRule(tenantID, ruleName, metricType string, threshold float64, comparison AlertComparisonOperator, severity AlertSeverity) *AlertRule {
	return &AlertRule{
		ID:                   generateID(),
		TenantID:             tenantID,
		RuleName:             ruleName,
		MetricType:           metricType,
		Threshold:            threshold,
		Comparison:           comparison,
		Severity:             severity,
		NotificationChannels: make(NotificationChannels, 0),
		EvaluationInterval:   60,  // 默认60秒
		ForDuration:          300, // 默认持续5分钟
		IsEnabled:            true,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// ShouldAlert 判断是否应该触发告警
func (r *AlertRule) ShouldAlert(value float64) bool {
	switch r.Comparison {
	case AlertComparisonGT:
		return value > r.Threshold
	case AlertComparisonLT:
		return value < r.Threshold
	case AlertComparisonEQ:
		return value == r.Threshold
	case AlertComparisonGTE:
		return value >= r.Threshold
	case AlertComparisonLTE:
		return value <= r.Threshold
	default:
		return false
	}
}

// AddNotificationChannel 添加通知渠道
func (r *AlertRule) AddNotificationChannel(channel NotificationChannel) {
	if r.NotificationChannels == nil {
		r.NotificationChannels = make(NotificationChannels, 0)
	}
	if !r.NotificationChannels.Contains(channel) {
		r.NotificationChannels = append(r.NotificationChannels, channel)
	}
}

// RemoveNotificationChannel 移除通知渠道
func (r *AlertRule) RemoveNotificationChannel(channel NotificationChannel) {
	for i, c := range r.NotificationChannels {
		if c == channel {
			r.NotificationChannels = append(r.NotificationChannels[:i], r.NotificationChannels[i+1:]...)
			break
		}
	}
}

// ============================================================
// AlertHistory 告警历史实体
// ============================================================

// AlertHistory 告警历史
type AlertHistory struct {
	ID              string              `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string              `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1"`
	AlertRuleID     string              `json:"alert_rule_id" gorm:"type:varchar(36);not null;index:idx_alert_rule,priority:1"`
	AlertRuleName   string              `json:"alert_rule_name" gorm:"type:varchar(100);not null"`
	Severity        AlertSeverity       `json:"severity" gorm:"type:varchar(20);not null;index:idx_severity,priority:1"`
	AlertMessage    string              `json:"alert_message" gorm:"type:text;not null"`
	AlertData       string              `json:"alert_data" gorm:"type:json"` // {"metric_value": 95.5, "threshold": 80}
	AgentID         string              `json:"agent_id,omitempty" gorm:"type:varchar(36);index:idx_agent_id"`
	Status          AlertStatus         `json:"status" gorm:"type:varchar(20);not null;index:idx_status,priority:1"`
	AcknowledgedBy  string              `json:"acknowledged_by" gorm:"type:varchar(36)"`
	AcknowledgedAt  *time.Time          `json:"acknowledged_at"`
	ResolvedAt      *time.Time          `json:"resolved_at"`
	SilencedUntil   *time.Time          `json:"silenced_until"`
	ResolvedBy      string              `json:"resolved_by" gorm:"type:varchar(36)"`
	ResolutionNote  string              `json:"resolution_note" gorm:"type:text"`
	CreatedAt       time.Time           `json:"created_at" gorm:"not null;index:idx_created_at,priority:2"`
	UpdatedAt       time.Time           `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (AlertHistory) TableName() string {
	return "alert_history"
}

// NewAlertHistory 创建告警历史
func NewAlertHistory(tenantID, alertRuleID, alertRuleName string, severity AlertSeverity, message, data string) *AlertHistory {
	return &AlertHistory{
		ID:            generateID(),
		TenantID:      tenantID,
		AlertRuleID:   alertRuleID,
		AlertRuleName: alertRuleName,
		Severity:      severity,
		AlertMessage:  message,
		AlertData:     data,
		Status:        AlertStatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// Acknowledge 确认告警
func (h *AlertHistory) Acknowledge(userID string) {
	h.Status = AlertStatusAcknowledged
	h.AcknowledgedBy = userID
	now := time.Now()
	h.AcknowledgedAt = &now
	h.UpdatedAt = now
}

// Resolve 解决告警
func (h *AlertHistory) Resolve(userID, note string) {
	h.Status = AlertStatusResolved
	h.ResolvedBy = userID
	h.ResolutionNote = note
	now := time.Now()
	h.ResolvedAt = &now
	h.UpdatedAt = now
}

// Silence 静默告警
func (h *AlertHistory) Silence(duration time.Duration) {
	h.Status = AlertStatusSilenced
	silencedUntil := time.Now().Add(duration)
	h.SilencedUntil = &silencedUntil
	h.UpdatedAt = time.Now()
}

// IsSilenced 检查告警是否已静默
func (h *AlertHistory) IsSilenced() bool {
	if h.Status != AlertStatusSilenced || h.SilencedUntil == nil {
		return false
	}
	return time.Now().Before(*h.SilencedUntil)
}

// GetDuration 获取告警持续时间
func (h *AlertHistory) GetDuration() time.Duration {
	endTime := h.ResolvedAt
	if endTime == nil {
		endTime = &time.Time{}
		*endTime = time.Now()
	}
	return endTime.Sub(h.CreatedAt)
}

// ============================================================
// PerformanceReport 性能报告实体
// ============================================================

// PerformanceReport 性能报告
type PerformanceReport struct {
	ID                   string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID             string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1"`
	AgentID              *string   `json:"agent_id,omitempty" gorm:"type:varchar(36);index:idx_agent_id"`
	ReportType           string    `json:"report_type" gorm:"type:varchar(50);not null;index:idx_report_type"` // daily, weekly, monthly
	ReportDate           string    `json:"report_date" gorm:"type:varchar(20);not null;index:idx_report_date"` // YYYY-MM-DD
	QPSMetrics           string    `json:"qps_metrics" gorm:"type:json"`       // {"avg": 1000, "p50": 950, "p95": 1100, "p99": 1500}
	ResponseTimeMetrics  string    `json:"response_time_metrics" gorm:"type:json"` // {"avg": 50, "p50": 45, "p95": 80, "p99": 120}
	ErrorRate            float64   `json:"error_rate" gorm:"type:decimal(5,4);not null"`
	Satisfaction         *float64  `json:"satisfaction" gorm:"type:decimal(3,2)"` // 0.00-1.00
	TokenUsage           *int64    `json:"token_usage"`
	EstimatedCost        *float64  `json:"estimated_cost" gorm:"type:decimal(15,4)"`
	Recommendations      string    `json:"recommendations" gorm:"type:text"` // JSON数组格式的建议
	Summary              string    `json:"summary" gorm:"type:text"`         // 报告摘要
	TopIssues            string    `json:"top_issues" gorm:"type:json"`      // 顶级问题列表
	Improvements         string    `json:"improvements" gorm:"type:json"`    // 改进点
	GeneratedBy          string    `json:"generated_by" gorm:"type:varchar(36)"` // 生成者（用户ID或system）
	CreatedAt            time.Time `json:"created_at" gorm:"not null"`
}

// TableName 指定表名
func (PerformanceReport) TableName() string {
	return "performance_reports"
}

// NewPerformanceReport 创建性能报告
func NewPerformanceReport(tenantID string, reportType, reportDate string) *PerformanceReport {
	return &PerformanceReport{
		ID:           generateID(),
		TenantID:     tenantID,
		ReportType:   reportType,
		ReportDate:   reportDate,
		ErrorRate:    0,
		GeneratedBy:  "system",
		CreatedAt:    time.Now(),
	}
}

// ============================================================
// AlertFilter 告警过滤器
// ============================================================

// AlertFilter 告警查询过滤器
type AlertFilter struct {
	TenantID     string              `json:"tenant_id"`
	AlertRuleID  string              `json:"alert_rule_id,omitempty"`
	Severity     AlertSeverity       `json:"severity,omitempty"`
	Status       AlertStatus         `json:"status,omitempty"`
	AgentID      string              `json:"agent_id,omitempty"`
	StartTime    *time.Time          `json:"start_time,omitempty"`
	EndTime      *time.Time          `json:"end_time,omitempty"`
	PageToken    string              `json:"page_token,omitempty"`
	PageSize     int                 `json:"page_size,omitempty"`
	SortBy       string              `json:"sort_by,omitempty"`       // created_at, severity
	SortOrder    string              `json:"sort_order,omitempty"`   // asc, desc
}

// ReportFilter 报告查询过滤器
type ReportFilter struct {
	TenantID   string    `json:"tenant_id"`
	AgentID    string    `json:"agent_id,omitempty"`
	ReportType string    `json:"report_type,omitempty"`
	StartDate  string    `json:"start_date,omitempty"` // YYYY-MM-DD
	EndDate    string    `json:"end_date,omitempty"`   // YYYY-MM-DD
	PageToken  string    `json:"page_token,omitempty"`
	PageSize   int       `json:"page_size,omitempty"`
}

// ValidateSeverity 验证告警严重级别
func ValidateSeverity(severity AlertSeverity) error {
	switch severity {
	case AlertSeverityWarning, AlertSeverityCritical, AlertSeverityEmergency:
		return nil
	default:
		return fmt.Errorf("invalid alert severity: %s", severity)
	}
}

// ValidateStatus 验证告警状态
func ValidateStatus(status AlertStatus) error {
	switch status {
	case AlertStatusPending, AlertStatusAcknowledged, AlertStatusResolved, AlertStatusSilenced:
		return nil
	default:
		return fmt.Errorf("invalid alert status: %s", status)
	}
}

// GetSeverityLevel 获取严重级别数值
func GetSeverityLevel(severity AlertSeverity) int {
	switch severity {
	case AlertSeverityEmergency:
		return 3
	case AlertSeverityCritical:
		return 2
	case AlertSeverityWarning:
		return 1
	default:
		return 0
	}
}

// GetSeverityPriority 获取严重级别优先级（用于排序）
func GetSeverityPriority(severity AlertSeverity) string {
	switch severity {
	case AlertSeverityEmergency:
		return "0"
	case AlertSeverityCritical:
		return "1"
	case AlertSeverityWarning:
		return "2"
	default:
		return "9"
	}
}
