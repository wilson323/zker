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
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"

	monitoringEntity "github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
	monitoringRepo "github.com/coze-dev/coze-studio/backend/domain/monitoring/repository"
)

// MetricsCollector 监控指标收集器
type MetricsCollector struct {
	tenantMetricsRepo monitoringRepo.TenantMetricsRepository
	logger            *zap.Logger

	// Prometheus指标
	requestsTotal     *prometheus.CounterVec
	requestDuration   *prometheus.HistogramVec
	responseTimeGauge *prometheus.GaugeVec
	qpsGauge          *prometheus.GaugeVec
	errorRateGauge    *prometheus.GaugeVec
}

// NewMetricsCollector 创建监控指标收集器
func NewMetricsCollector(
	tenantMetricsRepo monitoringRepo.TenantMetricsRepository,
	logger *zap.Logger,
) *MetricsCollector {
	collector := &MetricsCollector{
		tenantMetricsRepo: tenantMetricsRepo,
		logger:            logger,

		// 初始化Prometheus指标
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "coze_http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"tenant_id", "method", "endpoint", "status"},
		),

		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "coze_http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"tenant_id", "method", "endpoint"},
		),

		responseTimeGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "coze_response_time_milliseconds",
				Help: "Current response time in milliseconds",
			},
			[]string{"tenant_id", "endpoint"},
		),

		qpsGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "coze_qps",
				Help: "Current queries per second",
			},
			[]string{"tenant_id"},
		),

		errorRateGauge: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "coze_error_rate",
				Help: "Current error rate",
			},
			[]string{"tenant_id"},
		),
	}

	return collector
}

// CollectTenantMetrics 收集租户指标
func (c *MetricsCollector) CollectTenantMetrics(ctx context.Context, tenantID string) error {
	c.logger.Info("Collecting tenant metrics",
		zap.String("tenant_id", tenantID),
	)

	// 收集QPS指标
	if err := c.collectQPS(ctx, tenantID); err != nil {
		c.logger.Error("Failed to collect QPS",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
	}

	// 收集响应时间指标
	if err := c.collectResponseTime(ctx, tenantID); err != nil {
		c.logger.Error("Failed to collect response time",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
	}

	// 收集错误率指标
	if err := c.collectErrorRate(ctx, tenantID); err != nil {
		c.logger.Error("Failed to collect error rate",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
	}

	return nil
}

// CollectAllTenantsMetrics 收集所有租户指标（定时任务）
func (c *MetricsCollector) CollectAllTenantsMetrics(ctx context.Context, tenantIDs []string) error {
	c.logger.Info("Collecting metrics for all tenants",
		zap.Int("tenant_count", len(tenantIDs)),
	)

	for _, tenantID := range tenantIDs {
		if err := c.CollectTenantMetrics(ctx, tenantID); err != nil {
			c.logger.Error("Failed to collect metrics for tenant",
				zap.String("tenant_id", tenantID),
				zap.Error(err),
			)
			continue
		}
	}

	return nil
}

// RecordHTTPRequest 记录HTTP请求指标
func (c *MetricsCollector) RecordHTTPRequest(
	tenantID, method, endpoint, status string,
	duration time.Duration,
) {
	// 增加请求总数
	c.requestsTotal.WithLabelValues(tenantID, method, endpoint, status).Inc()

	// 记录请求持续时间
	c.requestDuration.WithLabelValues(tenantID, method, endpoint).Observe(duration.Seconds())
}

// RecordQPS 记录QPS指标
func (c *MetricsCollector) RecordQPS(tenantID string, qps float64) {
	c.qpsGauge.WithLabelValues(tenantID).Set(qps)
}

// RecordResponseTime 记录响应时间指标
func (c *MetricsCollector) RecordResponseTime(tenantID, endpoint string, responseTime float64) {
	c.responseTimeGauge.WithLabelValues(tenantID, endpoint).Set(responseTime)
}

// RecordErrorRate 记录错误率指标
func (c *MetricsCollector) RecordErrorRate(tenantID string, errorRate float64) {
	c.errorRateGauge.WithLabelValues(tenantID).Set(errorRate)
}

// ============================================================
// 私有方法
// ============================================================

// collectQPS 收集QPS指标
func (c *MetricsCollector) collectQPS(ctx context.Context, tenantID string) error {
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Minute)

	metrics, err := c.tenantMetricsRepo.CalculateQPSMetrics(ctx, tenantID, startTime, endTime)
	if err != nil {
		return err
	}

	if metrics != nil {
		// 更新Prometheus指标
		c.RecordQPS(tenantID, metrics.Average)

		// 记录到数据库
		metric := monitoringEntity.NewTenantMetrics(
			tenantID,
			monitoringEntity.MetricTypeQPS,
			metrics.Average,
			monitoringEntity.MetricTags{"source": "collector"},
		)
		return c.tenantMetricsRepo.Create(ctx, metric)
	}

	return nil
}

// collectResponseTime 收集响应时间指标
func (c *MetricsCollector) collectResponseTime(ctx context.Context, tenantID string) error {
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Minute)

	metrics, err := c.tenantMetricsRepo.CalculateResponseTimeMetrics(ctx, tenantID, startTime, endTime)
	if err != nil {
		return err
	}

	if metrics != nil {
		c.RecordResponseTime(tenantID, "all", metrics.Average)

		metric := monitoringEntity.NewTenantMetrics(
			tenantID,
			monitoringEntity.MetricTypeResponseTime,
			metrics.Average,
			monitoringEntity.MetricTags{"source": "collector"},
		)
		return c.tenantMetricsRepo.Create(ctx, metric)
	}

	return nil
}

// collectErrorRate 收集错误率指标
func (c *MetricsCollector) collectErrorRate(ctx context.Context, tenantID string) error {
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Minute)

	metrics, err := c.tenantMetricsRepo.CalculateErrorRateMetrics(ctx, tenantID, startTime, endTime)
	if err != nil {
		return err
	}

	if metrics != nil {
		c.RecordErrorRate(tenantID, metrics.ErrorRate)

		metric := monitoringEntity.NewTenantMetrics(
			tenantID,
			monitoringEntity.MetricTypeErrorRate,
			metrics.ErrorRate,
			monitoringEntity.MetricTags{"source": "collector"},
		)
		return c.tenantMetricsRepo.Create(ctx, metric)
	}

	return nil
}
