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

// MetricType 指标类型
type MetricType string

const (
	MetricTypeQPS          MetricType = "qps"           // 每秒查询率
	MetricTypeResponseTime MetricType = "response_time" // 响应时间
	MetricTypeErrorRate    MetricType = "error_rate"    // 错误率
	MetricTypeConcurrency  MetricType = "concurrency"   // 并发数
	MetricTypeCPUUsage     MetricType = "cpu_usage"     // CPU使用率
	MetricTypeMemoryUsage  MetricType = "memory_usage"  // 内存使用率
	MetricTypeDiskIO       MetricType = "disk_io"       // 磁盘IO
	MetricTypeNetworkIO    MetricType = "network_io"    // 网络IO
)

// String 返回字符串表示
func (m MetricType) String() string {
	return string(m)
}

// IsValid 验证指标类型是否有效
func (m MetricType) IsValid() bool {
	switch m {
	case MetricTypeQPS, MetricTypeResponseTime, MetricTypeErrorRate,
		MetricTypeConcurrency, MetricTypeCPUUsage, MetricTypeMemoryUsage,
		MetricTypeDiskIO, MetricTypeNetworkIO:
		return true
	default:
		return false
	}
}

// MetricTags 指标标签（JSON格式）
type MetricTags map[string]string

// Scan 实现 sql.Scanner 接口
func (t *MetricTags) Scan(value interface{}) error {
	if value == nil {
		*t = make(MetricTags)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal MetricTags value: %v", value)
	}

	return json.Unmarshal(bytes, t)
}

// Value 实现 driver.Valuer 接口
func (t MetricTags) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}

	return json.Marshal(t)
}

// TenantMetrics 租户指标实体
type TenantMetrics struct {
	ID              string       `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string       `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1;index:idx_tenant_type_time,priority:1"`
	MetricType      MetricType   `json:"metric_type" gorm:"type:varchar(50);not null;index:idx_metric_type,priority:1;index:idx_tenant_type_time,priority:2"`
	MetricValue     float64      `json:"metric_value" gorm:"type:decimal(15,4);not null"`
	MetricTimestamp time.Time    `json:"metric_timestamp" gorm:"not null;index:idx_timestamp,priority:1;index:idx_tenant_type_time,priority:3"`
	Tags            MetricTags   `json:"tags" gorm:"type:json"`
	CreatedAt       time.Time    `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time    `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (TenantMetrics) TableName() string {
	return "tenant_metrics"
}

// NewTenantMetrics 创建租户指标
func NewTenantMetrics(tenantID string, metricType MetricType, value float64, tags MetricTags) *TenantMetrics {
	return &TenantMetrics{
		ID:              generateID(),
		TenantID:        tenantID,
		MetricType:      metricType,
		MetricValue:     value,
		MetricTimestamp: time.Now(),
		Tags:            tags,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// IsExpired 检查指标是否过期（超过保留期）
func (m *TenantMetrics) IsExpired(retentionDays int) bool {
	expirationTime := time.Now().AddDate(0, 0, -retentionDays)
	return m.MetricTimestamp.Before(expirationTime)
}

// GetTag 获取标签
func (m *TenantMetrics) GetTag(key string) (string, bool) {
	if m.Tags == nil {
		return "", false
	}
	val, ok := m.Tags[key]
	return val, ok
}

// SetTag 设置标签
func (m *TenantMetrics) SetTag(key, value string) {
	if m.Tags == nil {
		m.Tags = make(MetricTags)
	}
	m.Tags[key] = value
}

// GetTags 获取所有标签
func (m *TenantMetrics) GetTags() MetricTags {
	if m.Tags == nil {
		return make(MetricTags)
	}
	return m.Tags
}

// ============================================================
// 聚合指标类型
// ============================================================

// QPSMetrics QPS指标聚合
type QPSMetrics struct {
	Count       int64   `json:"count"`       // 总请求数
	Average     float64 `json:"average"`     // 平均QPS
	Min         float64 `json:"min"`         // 最小QPS
	Max         float64 `json:"max"`         // 最大QPS
	P50         float64 `json:"p50"`         // P50分位数
	P95         float64 `json:"p95"`         // P95分位数
	P99         float64 `json:"p99"`         // P99分位数
	Trend       string  `json:"trend"`       // 趋势：up/down/stable
	TrendChange float64 `json:"trend_change"` // 趋势变化百分比
}

// ResponseTimeMetrics 响应时间指标聚合
type ResponseTimeMetrics struct {
	Average     float64 `json:"average"`     // 平均响应时间（ms）
	Min         float64 `json:"min"`         // 最小响应时间（ms）
	Max         float64 `json:"max"`         // 最大响应时间（ms）
	P50         float64 `json:"p50"`         // P50分位数
	P95         float64 `json:"p95"`         // P95分位数
	P99         float64 `json:"p99"`         // P99分位数
	Trend       string  `json:"trend"`       // 趋势：up/down/stable
	TrendChange float64 `json:"trend_change"` // 趋势变化百分比
}

// ErrorRateMetrics 错误率指标聚合
type ErrorRateMetrics struct {
	TotalRequests  int64   `json:"total_requests"`  // 总请求数
	ErrorRequests  int64   `json:"error_requests"`  // 错误请求数
	ErrorRate      float64 `json:"error_rate"`      // 错误率（0-1）
	Average        float64 `json:"average"`         // 平均错误率
	Max            float64 `json:"max"`             // 最大错误率
	Trend          string  `json:"trend"`           // 趋势：up/down/stable
	TrendChange    float64 `json:"trend_change"`    // 趋势变化百分比
	ErrorBreakdown map[string]int64 `json:"error_breakdown"` // 错误类型分布
}

// ConcurrencyMetrics 并发指标聚合
type ConcurrencyMetrics struct {
	Current       int64   `json:"current"`        // 当前并发数
	Average       float64 `json:"average"`        // 平均并发数
	Peak          int64   `json:"peak"`           // 峰值并发数
	PeakTime      *time.Time `json:"peak_time"`   // 峰值发生时间
	Trend         string  `json:"trend"`          // 趋势：up/down/stable
	TrendChange   float64 `json:"trend_change"`   // 趋势变化百分比
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// NewTimeRange 创建时间范围
func NewTimeRange(startTime, endTime time.Time) (*TimeRange, error) {
	if endTime.Before(startTime) {
		return nil, fmt.Errorf("end_time must be after start_time")
	}
	return &TimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// NewTimeRangeFromDuration 从时间间隔创建范围
func NewTimeRangeFromDuration(duration time.Duration) *TimeRange {
	return &TimeRange{
		StartTime: time.Now().Add(-duration),
		EndTime:   time.Now(),
	}
}

// Duration 返回时间范围长度
func (tr *TimeRange) Duration() time.Duration {
	return tr.EndTime.Sub(tr.StartTime)
}

// Contains 检查时间是否在范围内
func (tr *TimeRange) Contains(t time.Time) bool {
	return (t.Equal(tr.StartTime) || t.After(tr.StartTime)) &&
		(t.Equal(tr.EndTime) || t.Before(tr.EndTime))
}

// GetDefaultTimeRanges 获取默认时间范围列表
func GetDefaultTimeRanges() map[string]time.Duration {
	return map[string]time.Duration{
		"5m":     5 * time.Minute,
		"15m":    15 * time.Minute,
		"1h":     time.Hour,
		"6h":     6 * time.Hour,
		"24h":    24 * time.Hour,
		"7d":     7 * 24 * time.Hour,
		"30d":    30 * 24 * time.Hour,
		"90d":    90 * 24 * time.Hour,
	}
}

// HealthScore 健康评分（0-100）
type HealthScore struct {
	Score         float64    `json:"score"`          // 总分（0-100）
	Level         string     `json:"level"`          // 等级：excellent/good/warning/critical
	QPSScore      float64    `json:"qps_score"`      // QPS评分
	ResponseScore float64    `json:"response_score"` // 响应时间评分
	ErrorScore    float64    `json:"error_score"`    // 错误率评分
	ResourceScore float64    `json:"resource_score"` // 资源使用评分
	Trend         string     `json:"trend"`          // 趋势：improving/stable/degrading
	AssessedAt    time.Time  `json:"assessed_at"`    // 评估时间
}

// TenantOverview 租户概览
type TenantOverview struct {
	TenantID            string          `json:"tenant_id"`
	TenantName          string          `json:"tenant_name"`
	HealthScore         HealthScore     `json:"health_score"`
	CurrentQPS          float64         `json:"current_qps"`
	AvgResponseTime     float64         `json:"avg_response_time"`
	ErrorRate           float64         `json:"error_rate"`
	ActiveAgents        int64           `json:"active_agents"`
	TotalRequests       int64           `json:"total_requests"`
	CPUUsage            float64         `json:"cpu_usage"`
	MemoryUsage         float64         `json:"memory_usage"`
	DiskUsage           float64         `json:"disk_usage"`
	AlertCount          int64           `json:"alert_count"`
	CriticalAlertCount  int64           `json:"critical_alert_count"`
	SubscriptionTier    string          `json:"subscription_tier"`
	QuotaUsagePercent   float64         `json:"quota_usage_percent"`
	LastUpdated         time.Time       `json:"last_updated"`
}

// helper function to generate ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
