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
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/infra/monitoring/metrics"
)

// PrometheusMiddleware Prometheus监控中间件
// 功能:
// 1. 记录HTTP请求总数
// 2. 记录HTTP请求延迟
// 3. 记录请求/响应大小
func PrometheusMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		// 获取请求信息
		method := string(c.Request.Method())
		path := string(c.Request.URI().Path())

		// 执行请求
		c.Next(ctx)

		// 计算延迟
		duration := time.Since(start).Seconds()

		// 获取响应信息
		status := c.Response.StatusCode()
		requestSize := len(c.Request.Body())
		responseSize := len(c.Response.Body())

		// 记录指标
		metrics.RecordHTTPRequest(
			method,
			path,
			status,
			duration,
			requestSize,
			responseSize,
		)
	}
}

// InstrumentHandler 为单个handler添加Prometheus监控
func InstrumentHandler(handlerName string, handler app.HandlerFunc) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		// 调用实际handler
		handler(ctx, c)

		// 记录指标
		duration := time.Since(start).Seconds()
		method := string(c.Request.Method())
		status := c.Response.StatusCode()

		metrics.HTTPRequestDuration.WithLabelValues(method, handlerName).Observe(duration)
		metrics.HTTPRequestsTotal.WithLabelValues(method, handlerName, fmt.Sprintf("%d", status)).Inc()
	}
}
