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

package logging

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

// ============================================================
// 计费系统结构化日志
// ============================================================

const (
	// 日志字段键
	LogKeyTenantID       = "tenant_id"
	LogKeyUserID         = "user_id"
	LogKeyBotID          = "bot_id"
	LogKeyConversationID = "conversation_id"
	LogKeyModelProvider  = "model_provider"
	LogKeyModelName      = "model_name"
	LogKeyRequestType    = "request_type"
	LogKeyTokensInput    = "tokens_input"
	LogKeyTokensOutput   = "tokens_output"
	LogKeyTokensTotal    = "tokens_total"
	LogKeyCostInput      = "cost_input"
	LogKeyCostOutput     = "cost_output"
	LogKeyCostTotal      = "cost_total"
	LogKeyIsCached       = "is_cached"
	LogKeyResponseTime   = "response_time_ms"
	LogKeyDuration       = "duration_ms"
	LogKeyOperation      = "operation"
	LogKeyResult         = "result"
	LogKeyError          = "error"
	LogKeyTraceID        = "trace_id"
	LogKeySpanID         = "span_id"
	LogKeyBudgetType     = "budget_type"
	LogKeyBudgetAmount   = "budget_amount"
	LogKeyUsedAmount     = "used_amount"
	LogKeyUsagePercent   = "usage_percent"
	LogKeyAlertLevel     = "alert_level"
	LogKeyAlertType      = "alert_type"
	LogKeyChannel        = "channel"
)

var logger *zap.Logger

// Init 初始化计费日志
func Init(globalLogger *zap.Logger) {
	logger = globalLogger.Named("billing")
}

// ============================================================
// Token计量日志
// ============================================================

// LogTokenRecord 记录Token记录日志
func LogTokenRecord(ctx context.Context, req *service.RecordTokenUsageRequest, resp *service.RecordTokenUsageResponse, duration time.Duration) {
	logger.Info("Token usage recorded",
		zap.String(LogKeyTenantID, req.TenantID),
		zap.String(LogKeyUserID, formatUint64Ptr(req.UserID)),
		zap.String(LogKeyBotID, derefString(req.BotID)),
		zap.String(LogKeyConversationID, derefString(req.ConversationID)),
		zap.String(LogKeyModelProvider, req.ModelProvider),
		zap.String(LogKeyModelName, req.ModelName),
		zap.String(LogKeyRequestType, req.RequestType),
		zap.Int(LogKeyTokensInput, req.InputTokens),
		zap.Int(LogKeyTokensOutput, req.OutputTokens),
		zap.Int(LogKeyTokensTotal, req.TotalTokens),
		zap.Float64(LogKeyCostInput, resp.InputCost),
		zap.Float64(LogKeyCostOutput, resp.OutputCost),
		zap.Float64(LogKeyCostTotal, resp.TotalCost),
		zap.Bool(LogKeyIsCached, req.IsCached),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// LogBatchTokenRecord 记录批量Token记录日志
func LogBatchTokenRecord(ctx context.Context, tenantID string, recordCount int, totalCost float64, duration time.Duration) {
	logger.Info("Batch token usage recorded",
		zap.String(LogKeyTenantID, tenantID),
		zap.Int("record_count", recordCount),
		zap.Float64("total_cost", totalCost),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// LogPricingCalculation 记录定价计算日志
func LogPricingCalculation(ctx context.Context, provider, model string, inputTokens, outputTokens int, inputCost, outputCost, totalCost float64, duration time.Duration, err error) {
	fields := []zapcore.Field{
		zap.String(LogKeyModelProvider, provider),
		zap.String(LogKeyModelName, model),
		zap.Int(LogKeyTokensInput, inputTokens),
		zap.Int(LogKeyTokensOutput, outputTokens),
		zap.Float64(LogKeyCostInput, inputCost),
		zap.Float64(LogKeyCostOutput, outputCost),
		zap.Float64(LogKeyCostTotal, totalCost),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	}

	if err != nil {
		logger.Error("Pricing calculation failed", append(fields, zap.Error(err))...)
	} else {
		logger.Debug("Pricing calculation completed", fields...)
	}
}

// LogUsageStatsRetrieval 记录使用统计查询日志
func LogUsageStatsRetrieval(ctx context.Context, tenantID, startDate, endDate string, recordCount int, duration time.Duration) {
	logger.Info("Usage statistics retrieved",
		zap.String(LogKeyTenantID, tenantID),
		zap.String("start_date", startDate),
		zap.String("end_date", endDate),
		zap.Int("record_count", recordCount),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// ============================================================
// 预算告警日志
// ============================================================

// LogBudgetCheck 记录预算检查日志
func LogBudgetCheck(ctx context.Context, tenantID, budgetType string, budgetAmount, usedAmount, usagePercent float64, triggered bool, duration time.Duration) {
	level := zap.DebugLevel
	if triggered {
		level = zap.InfoLevel
	}

	logger.Log(level, "Budget check completed",
		zap.String(LogKeyTenantID, tenantID),
		zap.String(LogKeyBudgetType, budgetType),
		zap.Float64(LogKeyBudgetAmount, budgetAmount),
		zap.Float64(LogKeyUsedAmount, usedAmount),
		zap.Float64(LogKeyUsagePercent, usagePercent),
		zap.Bool("triggered", triggered),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// LogBudgetAlert 记录预算告警日志
func LogBudgetAlert(ctx context.Context, alert *entity.BudgetAlert) {
	logger.Warn("Budget alert sent",
		zap.String("alert_id", fmt.Sprintf("%d", alert.ID)),
		zap.String(LogKeyTenantID, alert.TenantID),
		zap.String(LogKeyAlertType, alert.AlertType),
		zap.Float64(LogKeyBudgetAmount, alert.BudgetAmount),
		zap.Float64(LogKeyUsedAmount, alert.UsedAmount),
		zap.Float64(LogKeyUsagePercent, alert.UsagePercent),
		zap.String(LogKeyAlertLevel, alert.AlertLevel),
		zap.String(LogKeyChannel, alert.NotificationChannels),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// LogBatchBudgetCheck 记录批量预算检查日志
func LogBatchBudgetCheck(ctx context.Context, tenantCount, alertCount int, duration time.Duration) {
	logger.Info("Batch budget check completed",
		zap.Int("tenant_count", tenantCount),
		zap.Int("alert_count", alertCount),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// ============================================================
// 汇总更新日志
// ============================================================

// LogSummaryUpdate 记录汇总更新日志
func LogSummaryUpdate(ctx context.Context, tenantID, summaryDate string, dailyTokens int64, dailyCost float64, duration time.Duration, err error) {
	fields := []zapcore.Field{
		zap.String(LogKeyTenantID, tenantID),
		zap.String("summary_date", summaryDate),
		zap.Int64("daily_tokens", dailyTokens),
		zap.Float64("daily_cost", dailyCost),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	}

	if err != nil {
		logger.Error("Summary update failed", append(fields, zap.Error(err))...)
	} else {
		logger.Info("Summary update completed", fields...)
	}
}

// ============================================================
// 通知发送日志
// ============================================================

// LogAlertNotification 记录告警通知发送日志
func LogAlertNotification(ctx context.Context, tenantID, channel, alertID string, duration time.Duration, err error) {
	fields := []zapcore.Field{
		zap.String(LogKeyTenantID, tenantID),
		zap.String(LogKeyChannel, channel),
		zap.String("alert_id", alertID),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	}

	if err != nil {
		logger.Error("Alert notification failed", append(fields, zap.Error(err))...)
	} else {
		logger.Info("Alert notification sent", fields...)
	}
}

// ============================================================
// 数据库操作日志
// ============================================================

// LogDBOperation 记录数据库操作日志
func LogDBOperation(ctx context.Context, operation, table string, recordCount int, duration time.Duration, err error) {
	level := zap.DebugLevel
	if err != nil {
		level = zap.ErrorLevel
	}

	fields := []zapcore.Field{
		zap.String(LogKeyOperation, operation),
		zap.String("table", table),
		zap.Int("record_count", recordCount),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	}

	if err != nil {
		logger.Log(level, "Database operation failed", append(fields, zap.Error(err))...)
	} else {
		logger.Log(level, "Database operation completed", fields...)
	}
}

// ============================================================
// 缓存操作日志
// ============================================================

// LogCacheOperation 记录缓存操作日志
func LogCacheOperation(ctx context.Context, cacheType, operation, key string, hit bool, duration time.Duration) {
	result := "miss"
	if hit {
		result = "hit"
	}

	logger.Debug("Cache operation",
		zap.String("cache_type", cacheType),
		zap.String(LogKeyOperation, operation),
		zap.String("key", key),
		zap.String("result", result),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// ============================================================
// 错误日志
// ============================================================

// LogError 记录错误日志
func LogError(ctx context.Context, operation string, err error, fields ...zapcore.Field) {
	allFields := append([]zapcore.Field{
		zap.String(LogKeyOperation, operation),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
		zap.Error(err),
	}, fields...)

	logger.Error("Operation failed", allFields...)
}

// LogErrorWithFields 记录带字段的错误日志
func LogErrorWithFields(ctx context.Context, operation string, err error, extraFields map[string]interface{}) {
	fields := make([]zapcore.Field, 0, len(extraFields)+3)
	fields = append(fields,
		zap.String(LogKeyOperation, operation),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
		zap.Error(err),
	)

	for k, v := range extraFields {
		fields = append(fields, zap.Any(k, v))
	}

	logger.Error("Operation failed with details", fields...)
}

// ============================================================
// 性能日志
// ============================================================

// LogSlowOperation 记录慢操作日志
func LogSlowOperation(ctx context.Context, operation string, duration time.Duration, threshold time.Duration) {
	logger.Warn("Slow operation detected",
		zap.String(LogKeyOperation, operation),
		zap.Int64(LogKeyDuration, duration.Milliseconds()),
		zap.Int64("threshold_ms", threshold.Milliseconds()),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// ============================================================
// 业务事件日志
// ============================================================

// LogBudgetConfigured 记录预算配置日志
func LogBudgetConfigured(ctx context.Context, tenantID, budgetType string, amount float64, thresholds []float64) {
	logger.Info("Budget configured",
		zap.String(LogKeyTenantID, tenantID),
		zap.String(LogKeyBudgetType, budgetType),
		zap.Float64("amount", amount),
		zap.Any("thresholds", thresholds),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// LogBudgetThresholdExceeded 记录预算超限日志
func LogBudgetThresholdExceeded(ctx context.Context, tenantID, budgetType string, usagePercent float64, threshold float64) {
	logger.Warn("Budget threshold exceeded",
		zap.String(LogKeyTenantID, tenantID),
		zap.String(LogKeyBudgetType, budgetType),
		zap.Float64(LogKeyUsagePercent, usagePercent),
		zap.Float64("threshold", threshold),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// LogPricingRuleCreated 记录定价规则创建日志
func LogPricingRuleCreated(ctx context.Context, provider, model string, inputPrice, outputPrice float64) {
	logger.Info("Pricing rule created",
		zap.String(LogKeyModelProvider, provider),
		zap.String(LogKeyModelName, model),
		zap.Float64("input_price_per_1k", inputPrice),
		zap.Float64("output_price_per_1k", outputPrice),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	)
}

// ============================================================
// 审计日志
// ============================================================

// LogAuditEvent 记录审计事件
func LogAuditEvent(ctx context.Context, eventType string, actor string, action string, target string, details map[string]interface{}) {
	fields := []zapcore.Field{
		zap.String("event_type", eventType),
		zap.String("actor", actor),
		zap.String("action", action),
		zap.String("target", target),
		zap.String(LogKeyTraceID, service.GetTraceID(ctx)),
	}

	for k, v := range details {
		fields = append(fields, zap.Any("detail_"+k, v))
	}

	logger.Info("Audit event", fields...)
}

// ============================================================
// 辅助函数
// ============================================================

// FormatDuration 格式化持续时间
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dμs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// FormatCost 格式化成本
func FormatCost(cost float64) string {
	if cost < 1 {
		return fmt.Sprintf("¥%.4f", cost)
	} else if cost < 100 {
		return fmt.Sprintf("¥%.2f", cost)
	} else {
		return fmt.Sprintf("¥%.0f", cost)
	}
}

// FormatTokens 格式化Token数量
func FormatTokens(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	} else if tokens < 1000000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1000)
	} else {
		return fmt.Sprintf("%.2fM", float64(tokens)/1000000)
	}
}

// ============================================================
// 辅助函数 - 处理指针类型
// ============================================================

// derefString 安全地解引用字符串指针，nil 返回空字符串
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// formatUint64Ptr 安全地格式化 uint64 指针为字符串
func formatUint64Ptr(u *uint64) string {
	if u == nil {
		return ""
	}
	return fmt.Sprintf("%d", *u)
}
