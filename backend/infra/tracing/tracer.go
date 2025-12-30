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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerProvider *sdktrace.TracerProvider
	tracer        trace.Tracer
)

// Init 初始化OpenTelemetry追踪
func Init(serviceName, jaegerEndpoint string) error {
	// 创建Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(jaegerEndpoint),
	))
	if err != nil {
		return err
	}

	// 创建resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		return err
	}

	// 创建tracer provider
	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // 10%采样率
	)

	// 设置全局tracer provider
	otel.SetTracerProvider(tracerProvider)

	// 获取tracer
	tracer = otel.Tracer(serviceName)

	return nil
}

// InitWithSampler 使用自定义采样率初始化
func InitWithSampler(serviceName, jaegerEndpoint string, sampleRatio float64) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(jaegerEndpoint),
	))
	if err != nil {
		return err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return err
	}

	tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRatio)),
	)

	otel.SetTracerProvider(tracerProvider)
	tracer = otel.Tracer(serviceName)

	return nil
}

// StartSpan 启动一个span
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName)
}

// StartSpanWithAttributes 启动带属性的span
func StartSpanWithAttributes(ctx context.Context, spanName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
	return ctx, span
}

// AddSpanAttributes 添加属性到当前span
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// AddSpanEvent 添加事件到当前span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// RecordError 记录错误到span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
}

// SetSpanStatus 设置span状态
func SetSpanStatus(ctx context.Context, desc string) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetStatus(codes.Error, desc)
	}
}

// GetTracer 获取tracer
func GetTracer() trace.Tracer {
	return tracer
}

// Shutdown 关闭tracer provider
func Shutdown(ctx context.Context) error {
	if tracerProvider != nil {
		return tracerProvider.Shutdown(ctx)
	}
	return nil
}

// Flush 刷新所有pending的spans
func Flush(ctx context.Context) error {
	if tracerProvider != nil {
		return tracerProvider.ForceFlush(ctx)
	}
	return nil
}

// ========== 常用属性 ==========

// 常用属性键
var (
	AttrTenantID     = attribute.Key("tenant_id")
	AttrUserID       = attribute.Key("user_id")
	AttrRequestID    = attribute.Key("request_id")
	AttrResourceType = attribute.Key("resource_type")
	AttrErrorCode    = attribute.Key("error_code")
	AttrHTTPMethod   = attribute.Key("http.method")
	AttrHTTPPath     = attribute.Key("http.path")
	AttrHTTPStatus   = attribute.Key("http.status_code")
	AttrDBSystem     = attribute.Key("db.system")
	AttrDBName       = attribute.Key("db.name")
	AttrDBTable      = attribute.Key("db.table")
	AttrDBOperation  = attribute.Key("db.operation")
	AttrCacheType    = attribute.Key("cache.type")
	AttrCacheKey     = attribute.Key("cache.key")
	AttrTTL          = attribute.Key("ttl")
	AttrKeyCount     = attribute.Key("key_count")
)

// ========== 辅助函数 ==========

// WithTenantID 添加租户ID属性
func WithTenantID(tenantID string) attribute.KeyValue {
	return AttrTenantID.String(tenantID)
}

// WithUserID 添加用户ID属性
func WithUserID(userID string) attribute.KeyValue {
	return AttrUserID.String(userID)
}

// WithRequestID 添加请求ID属性
func WithRequestID(requestID string) attribute.KeyValue {
	return AttrRequestID.String(requestID)
}

// WithError 添加错误属性
func WithError(err error) attribute.KeyValue {
	return attribute.String("error", err.Error())
}
