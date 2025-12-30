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

package monitoring

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/coze-dev/coze-studio/backend/infra/monitoring/metrics"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// init 强制初始化所有指标
//
// 通过下划线导入确保所有指标文件都被加载和初始化。
// 这是必须的，因为 Go 只初始化被直接引用的包级变量。
func init() {
	// 通过引用metrics包中的变量来触发所有指标的初始化
	// metrics.go中的init()函数会引用所有其他文件的指标变量
	var _ = metrics.HTTPRequestsTotal
}

// ========== HTTP监控中间件 ==========

// HTTPMiddleware HTTP请求监控中间件
func HTTPMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		method := string(c.Request.Method())
		path := string(c.Request.URI().Path())

		// 执行请求
		c.Next(ctx)

		// 请求结束
		duration := time.Since(start).Seconds()
		status := c.Response.StatusCode()

		// 使用 metrics 包中的函数记录指标
		metrics.RecordHTTPRequest(method, path, status, duration,
			len(c.Request.Body()), len(c.Response.Body()))
	}
}

// ========== Prometheus端点注册 ==========

// RegisterMetricsHandler 注册Prometheus指标端点
func RegisterMetricsHandler(r *server.Hertz) {
	// GET /metrics 端点
	r.GET("/metrics", func(ctx context.Context, c *app.RequestContext) {
		// 创建标准http.ResponseWriter适配器
		writer := &hertzResponseWriter{ctx: c}

		// 创建标准http.Request适配器
		body := c.Request.Body()
		req, err := http.NewRequest(
			string(c.Request.Method()),
			string(c.Request.URI().RequestURI()),
			bytes.NewReader(body),
		)
		if err != nil {
			c.String(consts.StatusInternalServerError, "Failed to create request: "+err.Error())
			return
		}

		// 复制请求头
		c.Request.Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// 使用prometheus标准处理器
		promhttp.Handler().ServeHTTP(writer, req)
	})

	// GET /health 健康检查端点
	r.GET("/health", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "coze-api",
		})
	})

	logs.Info("[Monitoring] Prometheus metrics endpoint registered at /metrics")
	logs.Info("[Monitoring] Health check endpoint registered at /health")
}

// ========== Hertz ResponseWriter适配器 ==========

// hertzResponseWriter 将Hertz的Response适配到http.ResponseWriter
type hertzResponseWriter struct {
	ctx        *app.RequestContext
	statusCode int
	written    bool
}

func (w *hertzResponseWriter) Header() http.Header {
	// 转换Hertz的Response.Header到http.Header
	h := make(http.Header)
	w.ctx.Response.Header.VisitAll(func(key, value []byte) {
		h.Set(string(key), string(value))
	})
	return h
}

func (w *hertzResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.statusCode = consts.StatusOK
		w.written = true
	}
	// Hertz使用SetBody而不是Write
	w.ctx.Response.SetBody(b)
	return len(b), nil
}

func (w *hertzResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.written = true
	w.ctx.Response.SetStatusCode(statusCode)
}

// ========== 便捷函数（使用 metrics 包） ==========

// RecordDBQuery 记录数据库查询指标
func RecordDBQuery(database, operation, table string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	metrics.RecordDBQuery(database, operation, table, duration.Seconds(), status)
}

// RecordCacheHit 记录缓存命中
func RecordCacheHit(cacheType string) {
	metrics.RecordCacheOperation(cacheType, "get", 0, true, "")
}

// RecordCacheMiss 记录缓存未命中
func RecordCacheMiss(cacheType string) {
	metrics.RecordCacheOperation(cacheType, "get", 0, false, "")
}

// RecordQuotaExceeded 记录配额超限
func RecordQuotaExceeded(tenantID, resourceType string) {
	metrics.RecordQuotaExceeded(tenantID, resourceType)
}

// UpdateQuotaUsage 更新配额使用情况
func UpdateQuotaUsage(tenantID, resourceType string, used, maxLimit int) {
	metrics.QuotaUsage.WithLabelValues(tenantID, resourceType).Set(float64(used))

	if maxLimit > 0 {
		percentage := float64(used) / float64(maxLimit) * 100
		alertLevel := "normal"
		if percentage >= 100 {
			alertLevel = "exceeded"
		} else if percentage >= 90 {
			alertLevel = "critical"
		} else if percentage >= 75 {
			alertLevel = "warning"
		}
		metrics.UpdateQuotaUsagePercent(tenantID, resourceType, alertLevel, percentage)
	}
}

// RecordIsolationUpgrade 记录隔离策略升级
func RecordIsolationUpgrade(fromStrategy, toStrategy string) {
	// 使用 metrics 包中的相关函数
	metrics.RecordRoutingExecution("", "upgrade", "success", 0, 0, 0)
}

// RecordRegistration 记录租户注册
func RecordRegistration() {
	metrics.RecordTenantCreation("new", "success")
}
