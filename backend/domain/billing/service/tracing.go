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

package service

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/coze-dev/coze-studio/backend/infra/tracing"
)

// ============================================================
// 计费系统 OpenTelemetry 分布式追踪
// ============================================================

// StartRecordSpan 开始Token记录Span
func StartRecordSpan(ctx context.Context, tenantID, modelProvider, modelName string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.WithTenantID(tenantID),
		attribute.String("model.provider", modelProvider),
		attribute.String("model.name", modelName),
		attribute.String("span.kind", "server"),
	}
	return tracing.StartSpanWithAttributes(ctx, "TokenMetering.Record", attrs...)
}

// StartBatchRecordSpan 开始批量Token记录Span
func StartBatchRecordSpan(ctx context.Context, tenantID string, recordCount int) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.WithTenantID(tenantID),
		attribute.Int("record.count", recordCount),
		attribute.String("span.kind", "server"),
	}
	return tracing.StartSpanWithAttributes(ctx, "TokenMetering.BatchRecord", attrs...)
}

// StartPricingCalculationSpan 开始定价计算Span
func StartPricingCalculationSpan(ctx context.Context, provider, model string, inputTokens, outputTokens int) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.String("model.provider", provider),
		attribute.String("model.name", model),
		attribute.Int("tokens.input", inputTokens),
		attribute.Int("tokens.output", outputTokens),
		attribute.String("span.kind", "internal"),
	}
	return tracing.StartSpanWithAttributes(ctx, "PricingEngine.Calculate", attrs...)
}

// StartBudgetCheckSpan 开始预算检查Span
func StartBudgetCheckSpan(ctx context.Context, tenantID string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.WithTenantID(tenantID),
		attribute.String("span.kind", "internal"),
	}
	return tracing.StartSpanWithAttributes(ctx, "BudgetAlertService.Check", attrs...)
}

// StartBatchBudgetCheckSpan 开始批量预算检查Span
func StartBatchBudgetCheckSpan(ctx context.Context, tenantCount int) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		attribute.Int("tenant.count", tenantCount),
		attribute.String("span.kind", "internal"),
	}
	return tracing.StartSpanWithAttributes(ctx, "BudgetAlertService.BatchCheck", attrs...)
}

// StartAlertSendSpan 开始告警发送Span
func StartAlertSendSpan(ctx context.Context, tenantID, channel string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.WithTenantID(tenantID),
		attribute.String("notification.channel", channel),
		attribute.String("span.kind", "client"),
	}
	return tracing.StartSpanWithAttributes(ctx, "NotificationService.SendAlert", attrs...)
}

// StartSummaryUpdateSpan 开始汇总更新Span
func StartSummaryUpdateSpan(ctx context.Context, tenantID, summaryDate string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.WithTenantID(tenantID),
		attribute.String("summary.date", summaryDate),
		attribute.String("span.kind", "internal"),
	}
	return tracing.StartSpanWithAttributes(ctx, "TokenMetering.UpdateSummary", attrs...)
}

// StartUsageStatsRetrievalSpan 开始使用统计查询Span
func StartUsageStatsRetrievalSpan(ctx context.Context, tenantID, startDate, endDate string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.WithTenantID(tenantID),
		attribute.String("query.start_date", startDate),
		attribute.String("query.end_date", endDate),
		attribute.String("span.kind", "server"),
	}
	return tracing.StartSpanWithAttributes(ctx, "TokenMetering.GetUsageStats", attrs...)
}

// StartDBQuerySpan 开始数据库查询Span
func StartDBQuerySpan(ctx context.Context, operation, table string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.AttrDBSystem.String("mysql"),
		tracing.AttrDBTable.String(table),
		tracing.AttrDBOperation.String(operation),
	}
	return tracing.StartSpanWithAttributes(ctx, "DB.Query", attrs...)
}

// StartCacheOperationSpan 开始缓存操作Span
func StartCacheOperationSpan(ctx context.Context, cacheType, operation, key string) (context.Context, trace.Span) {
	attrs := []attribute.KeyValue{
		tracing.AttrCacheType.String(cacheType),
		tracing.AttrCacheKey.String(key),
		attribute.String("cache.operation", operation),
		attribute.String("span.kind", "internal"),
	}
	return tracing.StartSpanWithAttributes(ctx, "Cache."+operation, attrs...)
}

// ============================================================
// Span 辅助函数
// ============================================================

// SetRecordTags 设置Token记录的Span标签
func SetRecordTags(span trace.Span, req *RecordTokenUsageRequest, result *RecordTokenUsageResponse) {
	span.SetAttributes(
		attribute.Int("tokens.input", req.InputTokens),
		attribute.Int("tokens.output", req.OutputTokens),
		attribute.Int("tokens.total", req.TotalTokens),
		attribute.Float64("cost.input", result.InputCost),
		attribute.Float64("cost.output", result.OutputCost),
		attribute.Float64("cost.total", result.TotalCost),
		attribute.String("request.type", req.RequestType),
		attribute.Bool("response.cached", req.IsCached),
	)
	if req.ResponseTimeMs != nil {
		span.SetAttributes(attribute.Int("response.time_ms", *req.ResponseTimeMs))
	}
}

// SetPricingTags 设置定价计算的Span标签
func SetPricingTags(span trace.Span, inputCost, outputCost, totalCost, unitPrice float64) {
	span.SetAttributes(
		attribute.Float64("cost.input", inputCost),
		attribute.Float64("cost.output", outputCost),
		attribute.Float64("cost.total", totalCost),
		attribute.Float64("unit.price", unitPrice),
	)
}

// SetBudgetCheckTags 设置预算检查的Span标签
func SetBudgetCheckTags(span trace.Span, budgetAmount, usedAmount, usagePercent float64, triggered bool) {
	span.SetAttributes(
		attribute.Float64("budget.amount", budgetAmount),
		attribute.Float64("budget.used", usedAmount),
		attribute.Float64("budget.usage_percent", usagePercent),
		attribute.Bool("budget.triggered", triggered),
	)
}

// SetAlertTags 设置告警的Span标签
func SetAlertTags(span trace.Span, alertType, alertLevel string, usagePercent float64) {
	span.SetAttributes(
		attribute.String("alert.type", alertType),
		attribute.String("alert.level", alertLevel),
		attribute.Float64("alert.usage_percent", usagePercent),
	)
}

// SetErrorTag 设置错误标签
func SetErrorTag(ctx context.Context, err error) {
	tracing.RecordError(ctx, err)
}

// SetSuccessTag 设置成功标签
func SetSuccessTag(span trace.Span, success bool) {
	span.SetAttributes(attribute.Bool("success", success))
	if !success {
		span.SetStatus(codes.Error, "operation failed")
	}
}

// AddEvent 添加事件到Span
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	tracing.AddSpanEvent(ctx, name, attrs...)
}

// ============================================================
// 性能追踪辅助函数
// ============================================================

// TraceOperation 追踪通用操作的执行时间
func TraceOperation(ctx context.Context, operationName string) func() {
	ctx, span := tracing.StartSpan(ctx, operationName)
	startTime := time.Now()

	return func() {
		duration := time.Since(startTime).Milliseconds()
		span.SetAttributes(attribute.Int64("duration_ms", duration))
		span.End()

		// 如果执行时间超过阈值，记录事件
		const threshold = 1000 // 1秒
		if duration > threshold {
			span.AddEvent("slow_operation", trace.WithAttributes(
				attribute.Int64("duration_ms", duration),
				attribute.Int64("threshold_ms", threshold),
			))
		}
	}
}

// TraceDBOperation 追踪数据库操作
func TraceDBOperation(ctx context.Context, operation, table string) func(error) {
	ctx, span := StartDBQuerySpan(ctx, operation, table)
	startTime := time.Now()

	return func(err error) {
		duration := time.Since(startTime).Milliseconds()
		span.SetAttributes(
			attribute.Int64("duration_ms", duration),
			attribute.Bool("db.query.success", err == nil),
		)

		if err != nil {
			tracing.RecordError(ctx, err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()
	}
}

// ============================================================
// 追踪上下文传播辅助函数
// ============================================================

// SpanFromContext 从上下文获取Span
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// ============================================================
// 日志与追踪集成
// ============================================================

// LogAndTraceToSpan 同时记录日志和追踪到Span
func LogAndTraceToSpan(span trace.Span, level string, message string, fields map[string]interface{}) {
	// 将关键字段记录到Span
	attrs := make([]attribute.KeyValue, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, attribute.String("log.fields."+k, toString(v)))
	}
	span.SetAttributes(attrs...)
	span.AddEvent(message, trace.WithAttributes(attrs...))
}

// toString 将任意值转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// GetTraceID 获取当前追踪ID
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	spanContext := span.SpanContext()
	return spanContext.TraceID().String()
}

// GetSpanID 获取当前Span ID
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	spanContext := span.SpanContext()
	return spanContext.SpanID().String()
}
