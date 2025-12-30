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

package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// PrometheusMiddleware Prometheus监控中间件
type PrometheusMiddleware struct {
	logger *zap.Logger

	// HTTP请求总数
	requestsTotal *prometheus.CounterVec

	// HTTP请求持续时间
	requestDuration *prometheus.HistogramVec

	// HTTP响应大小
	responseSize *prometheus.HistogramVec

	// 当前正在处理的请求数
	inFlightRequests *prometheus.GaugeVec
}

// NewPrometheusMiddleware 创建Prometheus中间件
func NewPrometheusMiddleware(logger *zap.Logger) *PrometheusMiddleware {
	return &PrometheusMiddleware{
		logger: logger,

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
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0, 10.0},
			},
			[]string{"tenant_id", "method", "endpoint"},
		),

		responseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "coze_http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"tenant_id", "method", "endpoint"},
		),

		inFlightRequests: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "coze_http_in_flight_requests",
				Help: "Current number of in-flight HTTP requests",
			},
			[]string{"tenant_id"},
		),
	}
}

// Middleware 返回Gin中间件函数
func (m *PrometheusMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过metrics端点的监控
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		// 从上下文获取tenant_id
		tenantID := m.getTenantID(c)

		// 增加正在处理的请求计数
		m.inFlightRequests.WithLabelValues(tenantID).Inc()
		defer m.inFlightRequests.WithLabelValues(tenantID).Dec()

		// 处理请求
		c.Next()

		// 记录指标
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		endpoint := m.normalizeEndpoint(c.FullPath())

		// 记录请求总数
		m.requestsTotal.WithLabelValues(tenantID, method, endpoint, status).Inc()

		// 记录请求持续时间
		m.requestDuration.WithLabelValues(tenantID, method, endpoint).Observe(duration)

		// 记录响应大小
		if c.Writer.Size() > 0 {
			m.responseSize.WithLabelValues(tenantID, method, endpoint).Observe(float64(c.Writer.Size()))
		}

		m.logger.Debug("HTTP request metrics recorded",
			zap.String("tenant_id", tenantID),
			zap.String("method", method),
			zap.String("endpoint", endpoint),
			zap.String("status", status),
			zap.Float64("duration", duration),
		)
	}
}

// getTenantID 从上下文获取租户ID
func (m *PrometheusMiddleware) getTenantID(c *gin.Context) string {
	// 尝试从header获取
	if tenantID := c.GetHeader("X-Tenant-ID"); tenantID != "" {
		return tenantID
	}

	// 尝试从查询参数获取
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		return tenantID
	}

	// 尝试从上下文获取
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}

	// 返回默认值
	return "unknown"
}

// normalizeEndpoint 标准化端点路径
func (m *PrometheusMiddleware) normalizeEndpoint(endpoint string) string {
	if endpoint == "" {
		return "/unknown"
	}
	return endpoint
}
