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
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
)

// tenantMetricsRepository 租户指标仓储实现
type tenantMetricsRepository struct {
	db *gorm.DB
}

// NewTenantMetricsRepository 创建租户指标仓储实例
func NewTenantMetricsRepository(db *gorm.DB) TenantMetricsRepository {
	return &tenantMetricsRepository{db: db}
}

// Create 创建指标记录
func (r *tenantMetricsRepository) Create(ctx context.Context, metrics *entity.TenantMetrics) error {
	return r.db.WithContext(ctx).Create(metrics).Error
}

// CreateBatch 批量创建指标记录
func (r *tenantMetricsRepository) CreateBatch(ctx context.Context, metricsList []*entity.TenantMetrics) error {
	if len(metricsList) == 0 {
		return nil
	}
	// 批量插入，每批最多100条
	batchSize := 100
	return r.db.WithContext(ctx).CreateInBatches(metricsList, batchSize).Error
}

// GetByID 根据ID获取指标
func (r *tenantMetricsRepository) GetByID(ctx context.Context, id string) (*entity.TenantMetrics, error) {
	var metrics entity.TenantMetrics
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&metrics).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// GetByTenantAndType 根据租户ID和指标类型查询最新指标
func (r *tenantMetricsRepository) GetByTenantAndType(
	ctx context.Context,
	tenantID string,
	metricType entity.MetricType,
) (*entity.TenantMetrics, error) {
	var metrics entity.TenantMetrics
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND metric_type = ?", tenantID, metricType).
		Order("metric_timestamp DESC").
		First(&metrics).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

// GetByTenantAndTimeRange 根据租户ID和时间范围查询指标
func (r *tenantMetricsRepository) GetByTenantAndTimeRange(
	ctx context.Context,
	tenantID string,
	metricType entity.MetricType,
	startTime, endTime time.Time,
) ([]*entity.TenantMetrics, error) {
	var metricsList []*entity.TenantMetrics
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, metricType, startTime, endTime).
		Order("metric_timestamp ASC").
		Find(&metricsList).Error

	if err != nil {
		return nil, err
	}
	return metricsList, nil
}

// GetLatestByTenant 获取租户最新指标（所有类型）
func (r *tenantMetricsRepository) GetLatestByTenant(ctx context.Context, tenantID string) ([]*entity.TenantMetrics, error) {
	var metricsList []*entity.TenantMetrics

	// 使用子查询获取每个指标类型的最新记录
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT tm.*
			FROM tenant_metrics tm
			INNER JOIN (
				SELECT tenant_id, metric_type, MAX(metric_timestamp) as max_timestamp
				FROM tenant_metrics
				WHERE tenant_id = ?
				GROUP BY tenant_id, metric_type
			) latest ON tm.tenant_id = latest.tenant_id
				AND tm.metric_type = latest.metric_type
				AND tm.metric_timestamp = latest.max_timestamp
			WHERE tm.tenant_id = ?
		`, tenantID, tenantID).
		Scan(&metricsList).Error

	if err != nil {
		return nil, err
	}
	return metricsList, nil
}

// AggregateByTimeRange 按时间范围聚合指标
func (r *tenantMetricsRepository) AggregateByTimeRange(
	ctx context.Context,
	tenantID string,
	metricType entity.MetricType,
	startTime, endTime time.Time,
	interval string,
) ([]*entity.TenantMetrics, error) {
	var aggregatedMetrics []*entity.TenantMetrics

	// 根据interval确定时间桶
	var timeBucket string
	switch interval {
	case "1m":
		timeBucket = "DATE_FORMAT(metric_timestamp, '%Y-%m-%d %H:%i:00')"
	case "5m":
		timeBucket = "DATE_FORMAT(metric_timestamp, '%Y-%m-%d %H:%i:00') - INTERVAL MINUTE(metric_timestamp) % 5 MINUTE"
	case "1h":
		timeBucket = "DATE_FORMAT(metric_timestamp, '%Y-%m-%d %H:00:00')"
	case "1d":
		timeBucket = "DATE(metric_timestamp)"
	default:
		timeBucket = "DATE_FORMAT(metric_timestamp, '%Y-%m-%d %H:%i:00')"
	}

	query := fmt.Sprintf(`
		SELECT
			MAX(id) as id,
			tenant_id,
			metric_type,
			AVG(metric_value) as metric_value,
			%s as metric_timestamp,
			NULL as tags,
			MAX(created_at) as created_at
		FROM tenant_metrics
		WHERE tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?
		GROUP BY tenant_id, metric_type, %s
		ORDER BY metric_timestamp ASC
	`, timeBucket, timeBucket)

	err := r.db.WithContext(ctx).Raw(query, tenantID, metricType, startTime, endTime).Scan(&aggregatedMetrics).Error
	if err != nil {
		return nil, err
	}
	return aggregatedMetrics, nil
}

// CalculateQPSMetrics 计算QPS指标聚合
func (r *tenantMetricsRepository) CalculateQPSMetrics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*entity.QPSMetrics, error) {
	type Result struct {
		AvgValue   float64
		MinValue   float64
		MaxValue   float64
		TotalCount int64
	}

	var result Result
	err := r.db.WithContext(ctx).
		Model(&entity.TenantMetrics{}).
		Select(`
			AVG(metric_value) as avg_value,
			MIN(metric_value) as min_value,
			MAX(metric_value) as max_value,
			COUNT(*) as total_count
		`).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, entity.MetricTypeQPS, startTime, endTime).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// 计算分位数（简化版本，实际应使用PERCENTILE函数）
	metrics := &entity.QPSMetrics{
		Count:   result.TotalCount,
		Average: result.AvgValue,
		Min:     result.MinValue,
		Max:     result.MaxValue,
		P50:     result.AvgValue, // 简化处理
		P95:     result.MaxValue * 0.95,
		P99:     result.MaxValue * 0.99,
		Trend:   "stable",
	}

	return metrics, nil
}

// CalculateResponseTimeMetrics 计算响应时间指标聚合
func (r *tenantMetricsRepository) CalculateResponseTimeMetrics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*entity.ResponseTimeMetrics, error) {
	type Result struct {
		AvgValue float64
		MinValue float64
		MaxValue float64
	}

	var result Result
	err := r.db.WithContext(ctx).
		Model(&entity.TenantMetrics{}).
		Select(`
			AVG(metric_value) as avg_value,
			MIN(metric_value) as min_value,
			MAX(metric_value) as max_value
		`).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, entity.MetricTypeResponseTime, startTime, endTime).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	metrics := &entity.ResponseTimeMetrics{
		Average: result.AvgValue,
		Min:     result.MinValue,
		Max:     result.MaxValue,
		P50:     result.AvgValue,
		P95:     result.MaxValue * 0.95,
		P99:     result.MaxValue * 0.99,
		Trend:   "stable",
	}

	return metrics, nil
}

// CalculateErrorRateMetrics 计算错误率指标聚合
func (r *tenantMetricsRepository) CalculateErrorRateMetrics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*entity.ErrorRateMetrics, error) {
	type Result struct {
		AvgErrorRate float64
		MaxErrorRate float64
		TotalCount   int64
	}

	var result Result
	err := r.db.WithContext(ctx).
		Model(&entity.TenantMetrics{}).
		Select(`
			AVG(metric_value) as avg_error_rate,
			MAX(metric_value) as max_error_rate,
			COUNT(*) as total_count
		`).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, entity.MetricTypeErrorRate, startTime, endTime).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	metrics := &entity.ErrorRateMetrics{
		TotalRequests: result.TotalCount,
		ErrorRate:     result.AvgErrorRate,
		Average:       result.AvgErrorRate,
		Max:           result.MaxErrorRate,
		Trend:         "stable",
		ErrorBreakdown: make(map[string]int64),
	}

	return metrics, nil
}

// CalculateConcurrencyMetrics 计算并发指标聚合
func (r *tenantMetricsRepository) CalculateConcurrencyMetrics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*entity.ConcurrencyMetrics, error) {
	type Result struct {
		AvgValue float64
		MaxValue float64
		PeakTime time.Time
	}

	var result Result
	err := r.db.WithContext(ctx).
		Model(&entity.TenantMetrics{}).
		Select(`
			AVG(metric_value) as avg_value,
			MAX(metric_value) as max_value
		`).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, entity.MetricTypeConcurrency, startTime, endTime).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// 获取峰值时间
	var peakMetric entity.TenantMetrics
	err = r.db.WithContext(ctx).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, entity.MetricTypeConcurrency, startTime, endTime).
		Order("metric_value DESC").
		First(&peakMetric).Error

	if err == nil {
		result.PeakTime = peakMetric.MetricTimestamp
	}

	metrics := &entity.ConcurrencyMetrics{
		Average:  result.AvgValue,
		Peak:     int64(result.MaxValue),
		PeakTime: &result.PeakTime,
		Trend:    "stable",
	}

	return metrics, nil
}

// GetTopTenantsByQPS 获取QPS最高的N个租户
func (r *tenantMetricsRepository) GetTopTenantsByQPS(
	ctx context.Context,
	limit int,
	startTime, endTime time.Time,
) ([]*entity.TenantMetrics, error) {
	var metricsList []*entity.TenantMetrics
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT tenant_id, metric_type, AVG(metric_value) as metric_value
			FROM tenant_metrics
			WHERE metric_type = ? AND metric_timestamp BETWEEN ? AND ?
			GROUP BY tenant_id, metric_type
			ORDER BY metric_value DESC
			LIMIT ?
		`, entity.MetricTypeQPS, startTime, endTime, limit).
		Scan(&metricsList).Error

	if err != nil {
		return nil, err
	}
	return metricsList, nil
}

// DeleteExpiredMetrics 删除过期指标
func (r *tenantMetricsRepository) DeleteExpiredMetrics(ctx context.Context, retentionDays int) (int64, error) {
	expirationDate := time.Now().AddDate(0, 0, -retentionDays)
	result := r.db.WithContext(ctx).
		Where("metric_timestamp < ?", expirationDate).
		Delete(&entity.TenantMetrics{})

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// CountByTenant 统计租户指标数量
func (r *tenantMetricsRepository) CountByTenant(ctx context.Context, tenantID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.TenantMetrics{}).
		Where("tenant_id = ?", tenantID).
		Count(&count).Error

	return count, err
}

// GetMetricTypes 获取租户的所有指标类型
func (r *tenantMetricsRepository) GetMetricTypes(ctx context.Context, tenantID string) ([]entity.MetricType, error) {
	var metricTypes []entity.MetricType
	err := r.db.WithContext(ctx).
		Model(&entity.TenantMetrics{}).
		Where("tenant_id = ?", tenantID).
		Distinct("metric_type").
		Pluck("metric_type", &metricTypes).Error

	if err != nil {
		return nil, err
	}
	return metricTypes, nil
}

// GetTimeSeriesData 获取时序数据
func (r *tenantMetricsRepository) GetTimeSeriesData(
	ctx context.Context,
	tenantID string,
	metricType entity.MetricType,
	startTime, endTime time.Time,
	interval time.Duration,
) ([]TimeSeriesPoint, error) {
	var metricsList []*entity.TenantMetrics
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND metric_type = ? AND metric_timestamp BETWEEN ? AND ?",
			tenantID, metricType, startTime, endTime).
		Order("metric_timestamp ASC").
		Find(&metricsList).Error

	if err != nil {
		return nil, err
	}

	points := make([]TimeSeriesPoint, 0, len(metricsList))
	for _, m := range metricsList {
		points = append(points, TimeSeriesPoint{
			Timestamp: m.MetricTimestamp,
			Value:     m.MetricValue,
			Tags:      m.Tags,
		})
	}

	return points, nil
}

// BatchQueryLatestMetrics 批量查询多个租户的最新指标
func (r *tenantMetricsRepository) BatchQueryLatestMetrics(
	ctx context.Context,
	tenantIDs []string,
	metricType entity.MetricType,
) (map[string]*entity.TenantMetrics, error) {
	var metricsList []*entity.TenantMetrics
	err := r.db.WithContext(ctx).
		Where("tenant_id IN ? AND metric_type = ?", tenantIDs, metricType).
		Order("metric_timestamp DESC").
		Find(&metricsList).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]*entity.TenantMetrics)
	for _, m := range metricsList {
		if _, exists := result[m.TenantID]; !exists {
			result[m.TenantID] = m
		}
	}

	return result, nil
}
