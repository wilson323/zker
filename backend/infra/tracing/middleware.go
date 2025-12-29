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

package tracing

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-studio/coze-studio/backend/infra/tracing"
)

// TracingMiddleware 分布式追踪中间件
// 功能:
// 1. 从HTTP请求头中提取trace context
// 2. 为每个请求创建span
// 3. 自动记录请求/响应信息
// 4. 将trace context注入到response header
func TracingMiddleware(serviceName string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头中提取trace context
		carrier := make(propagation.HeaderCarrier)
		carrier.Set("traceparent", string(c.Request.Header.Get("traceparent")))
		carrier.Set("tracestate", string(c.Request.Header.Get("tracestate")))

		// 提取或创建span context
		spanCtx := propagator.Extract(ctx, carrier)

		// 创建span
		method := string(c.Request.Method())
		path := string(c.Request.URI().Path())
		spanName := method + " " + path

		spanCtx, span := tracing.StartSpan(
			spanCtx,
			spanName,
		)

		// 添加span属性
		span.SetAttributes(
			tracing.AttrHTTPMethod.String(method),
			tracing.AttrHTTPPath.String(path),
			tracing.AttrHTTPKey.String(c.Request.Host()),
			tracing.AttrRemoteAddr.String(c.ClientIP()),
			tracing.AttrUserAgent.String(string(c.Request.Header.UserAgent())),
		)

		// 从context中提取的业务属性
		if tenantID := ctx.Value("tenant_id"); tenantID != nil {
			if id, ok := tenantID.(string); ok {
				span.SetAttributes(tracing.AttrTenantID.String(id))
			}
		}

		if userID := ctx.Value("user_id"); userID != nil {
			if id, ok := userID.(string); ok {
				span.SetAttributes(tracing.AttrUserID.String(id))
			}
		}

		// 将span context注入到context
		ctx = trace.ContextWithSpan(spanCtx, span)

		// 执行请求
		c.Next(ctx)

		// 记录响应状态
		status := c.Response.StatusCode()
		span.SetAttributes(
			tracing.AttrHTTPStatus.Int(status),
		)

		// 根据状态码设置span状态
		if status >= 400 && status < 500 {
			span.SetStatus(trace.StatusCodeError, "Client Error")
		} else if status >= 500 {
			span.SetStatus(trace.StatusCodeError, "Server Error")
		} else {
			span.SetStatus(trace.StatusCodeOK, "OK")
		}

		// 结束span
		span.End()

		// 将trace context注入到response header
		propagator.Inject(spanCtx, carrier)
		c.Response.Header.Set("traceparent", string(carrier.Get("traceparent")))
		c.Response.Header.Set("tracestate", string(carrier.Get("tracestate")))
	}
}

// ========== 传播器 ==========

// W3CTracePropagator W3C trace传播器
type W3CTracePropagator struct{}

// Extract 从carrier中提取trace context
func (p *W3CTracePropagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	traceParent := carrier.Get("traceparent")
	if traceParent == "" {
		return ctx
	}

	// 解析traceparent header
	// 格式: version-trace_id-span_id-trace_flags
	// 示例: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
	// TODO: 实际解析逻辑
	return ctx
}

// Inject 将trace context注入到carrier
func (p *W3CTracePropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return
	}

	spanContext := span.SpanContext()
	if !spanContext.IsValid() {
		return
	}

	// 注入traceparent header
	// TODO: 实际生成traceparent header
	carrier.Set("traceparent", "00-00000000000000000000000000000000-0000000000000000-00")
}

// Fields 返回传播器支持的字段
func (p *W3CTracePropagator) Fields() []string {
	return []string{"traceparent", "tracestate"}
}

var propagator = &W3CTracePropagator{}

// ========== 数据库追踪辅助函数 ==========

// TraceDBQuery 追踪数据库查询
func TraceDBQuery(ctx context.Context, dbSystem, dbName, table, operation string, fn func() error) error {
	attrs := []attribute.KeyValue{
		tracing.AttrDBSystem.String(dbSystem),
		tracing.AttrDBName.String(dbName),
		tracing.AttrDBTable.String(table),
		tracing.AttrDBOperation.String(operation),
	}

	ctx, span := tracing.StartSpanWithAttributes(ctx, "db.query", attrs...)

	err := fn()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetStatus(trace.StatusCodeOK, "OK")
	}

	span.End()
	return err
}

// TraceDBTransaction 追踪数据库事务
func TraceDBTransaction(ctx context.Context, dbSystem, dbName string, fn func() error) error {
	attrs := []attribute.KeyValue{
		tracing.AttrDBSystem.String(dbSystem),
		tracing.AttrDBName.String(dbName),
	}

	ctx, span := tracing.StartSpanWithAttributes(ctx, "db.transaction", attrs...)

	err := fn()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetStatus(trace.StatusCodeOK, "OK")
	}

	span.End()
	return err
}

// ========== 缓存追踪辅助函数 ==========

// TraceCacheOperation 追踪缓存操作
func TraceCacheOperation(ctx context.Context, cacheType, operation, key string, fn func() (interface{}, error)) (interface{}, error) {
	attrs := []attribute.KeyValue{
		tracing.AttrCacheType.String(cacheType),
		attribute.String("cache.operation", operation),
		tracing.AttrCacheKey.String(key),
	}

	ctx, span := tracing.StartSpanWithAttributes(ctx, "cache.operation", attrs...)

	result, err := fn()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetStatus(trace.StatusCodeOK, "OK")
	}

	span.End()
	return result, err
}

// ========== HTTP客户端追踪辅助函数 ==========

// TraceHTTPClientRequest 追踪HTTP客户端请求
func TraceHTTPClientRequest(ctx context.Context, method, url string, fn func() error) error {
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.url", url),
	}

	ctx, span := tracing.StartSpanWithAttributes(ctx, "http.request", attrs...)

	err := fn()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetStatus(trace.StatusCodeOK, "OK")
	}

	span.End()
	return err
}

// ========== 业务操作追踪辅助函数 ==========

// TraceTenantOperation 追踪租户操作
func TraceTenantOperation(ctx context.Context, operation, tenantID string, fn func() error) error {
	attrs := []attribute.KeyValue{
		attribute.String("tenant.operation", operation),
		tracing.AttrTenantID.String(tenantID),
	}

	ctx, span := tracing.StartSpanWithAttributes(ctx, "tenant.operation", attrs...)

	err := fn()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetStatus(trace.StatusCodeOK, "OK")
	}

	span.End()
	return err
}

// TraceQuotaCheck 追踪配额检查
func TraceQuotaCheck(ctx context.Context, resourceType string, amount int, fn func() (bool, error)) (bool, error) {
	attrs := []attribute.KeyValue{
		tracing.AttrResourceType.String(resourceType),
		attribute.Int("quota.amount", amount),
	}

	ctx, span := tracing.StartSpanWithAttributes(ctx, "quota.check", attrs...)

	allowed, err := fn()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusCodeError, err.Error())
	} else {
		span.SetAttributes(attribute.Bool("quota.allowed", allowed))
		if !allowed {
			span.SetStatus(trace.StatusCodeError, "Quota Exceeded")
		} else {
			span.SetStatus(trace.StatusCodeOK, "OK")
		}
	}

	span.End()
	return allowed, err
}
