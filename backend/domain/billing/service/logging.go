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

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ============================================================
// 计费系统结构化日志
// ============================================================

// LogTokenRecord 记录Token使用日志
func LogTokenRecord(ctx context.Context, req *RecordTokenUsageRequest, result *RecordTokenUsageResponse, duration time.Duration) {
	logs.CtxInfof(ctx,
		"[TokenMetering] Record completed: tenant=%s model=%s/%s tokens=%d/%d/%d cost=%.4f duration=%dms log_id=%d",
		req.TenantID,
		req.ModelProvider,
		req.ModelName,
		req.InputTokens,
		req.OutputTokens,
		req.TotalTokens,
		result.TotalCost,
		duration.Milliseconds(),
		result.LogID,
	)
}

// LogTokenRecordError 记录Token记录错误日志
func LogTokenRecordError(ctx context.Context, req *RecordTokenUsageRequest, err error) {
	logs.CtxErrorf(ctx,
		"[TokenMetering] Record failed: tenant=%s model=%s/%s error=%v",
		req.TenantID,
		req.ModelProvider,
		req.ModelName,
		err,
	)
}

// LogBatchTokenRecord 记录批量Token使用日志
func LogBatchTokenRecord(ctx context.Context, count int, successCount, failedCount int, totalCost float64, duration time.Duration) {
	logs.CtxInfof(ctx,
		"[TokenMetering] Batch record completed: count=%d success=%d failed=%d cost=%.4f duration=%dms avg=%.2fms",
		count,
		successCount,
		failedCount,
		totalCost,
		duration.Milliseconds(),
		float64(duration.Milliseconds())/float64(count),
	)
}

// LogBatchTokenRecordError 记录批量Token记录错误日志
func LogBatchTokenRecordError(ctx context.Context, count int, err error) {
	logs.CtxErrorf(ctx,
		"[TokenMetering] Batch record failed: count=%d error=%v",
		count,
		err,
	)
}

// LogPricingCalculation 记录定价计算日志
func LogPricingCalculation(ctx context.Context, modelProvider, modelName string, inputTokens, outputTokens int, inputCost, outputCost, totalCost float64) {
	logs.CtxDebugf(ctx,
		"[PricingEngine] Calculation completed: model=%s/%s input_tokens=%d output_tokens=%d input_cost=%.6f output_cost=%.6f total_cost=%.6f",
		modelProvider,
		modelName,
		inputTokens,
		outputTokens,
		inputCost,
		outputCost,
		totalCost,
	)
}

// LogPricingCalculationError 记录定价计算错误日志
func LogPricingCalculationError(ctx context.Context, modelProvider, modelName string, err error) {
	logs.CtxErrorf(ctx,
		"[PricingEngine] Calculation failed: model=%s/%s error=%v",
		modelProvider,
		modelName,
		err,
	)
}

// LogBudgetCheck 记录预算检查日志
func LogBudgetCheck(ctx context.Context, tenantID string, budgetAmount, usedAmount, usagePercent float64, triggered bool) {
	if triggered {
		logs.CtxWarnf(ctx,
			"[BudgetAlert] Budget exceeded: tenant=%s budget=%.2f used=%.2f usage=%.2f%%",
			tenantID,
			budgetAmount,
			usedAmount,
			usagePercent,
		)
	} else {
		logs.CtxDebugf(ctx,
			"[BudgetAlert] Budget OK: tenant=%s budget=%.2f used=%.2f usage=%.2f%%",
			tenantID,
			budgetAmount,
			usedAmount,
			usagePercent,
		)
	}
}

// LogBudgetCheckError 记录预算检查错误日志
func LogBudgetCheckError(ctx context.Context, tenantID string, err error) {
	logs.CtxErrorf(ctx,
		"[BudgetAlert] Budget check failed: tenant=%s error=%v",
		tenantID,
		err,
	)
}

// LogBudgetAlert 记录预算告警日志
func LogBudgetAlert(ctx context.Context, tenantID string, alertType, alertLevel string, budgetAmount, usedAmount, usagePercent float64) {
	levelStr := "INFO"
	if alertLevel == "critical" || alertLevel == "emergency" {
		levelStr = "ERROR"
	} else if alertLevel == "warning" {
		levelStr = "WARN"
	}

	message := fmt.Sprintf(
		"[BudgetAlert] Alert triggered: tenant=%s type=%s level=%s budget=%.2f used=%.2f usage=%.2f%%",
		tenantID,
		alertType,
		alertLevel,
		budgetAmount,
		usedAmount,
		usagePercent,
	)

	if levelStr == "ERROR" {
		logs.CtxErrorf(ctx, message)
	} else if levelStr == "WARN" {
		logs.CtxWarnf(ctx, message)
	} else {
		logs.CtxInfof(ctx, message)
	}
}

// LogBudgetHardCap 记录硬性上限触发日志
func LogBudgetHardCap(ctx context.Context, tenantID string, budgetAmount, usedAmount float64) {
	logs.CtxErrorf(ctx,
		"[BudgetAlert] HARD CAP TRIGGERED: tenant=%s budget=%.2f used=%.2f service_paused=true",
		tenantID,
		budgetAmount,
		usedAmount,
	)
}

// LogBudgetDowngrade 记录模型降级日志
func LogBudgetDowngrade(ctx context.Context, tenantID, originalModel, fallbackModel string) {
	logs.CtxWarnf(ctx,
		"[BudgetAlert] Model downgrade executed: tenant=%s original=%s fallback=%s",
		tenantID,
		originalModel,
		fallbackModel,
	)
}

// LogAlertNotification 记录告警通知发送日志
func LogAlertNotification(ctx context.Context, tenantID, channel string, success bool, duration time.Duration) {
	if success {
		logs.CtxInfof(ctx,
			"[NotificationService] Alert sent: tenant=%s channel=%s duration=%dms",
			tenantID,
			channel,
			duration.Milliseconds(),
		)
	} else {
		logs.CtxErrorf(ctx,
			"[NotificationService] Alert send failed: tenant=%s channel=%s duration=%dms",
			tenantID,
			channel,
			duration.Milliseconds(),
		)
	}
}

// LogSummaryUpdate 记录汇总更新日志
func LogSummaryUpdate(ctx context.Context, tenantID string, summaryDate string, summaryHour *uint8, totalTokens int64, totalCost float64, duration time.Duration) {
	hourStr := "all"
	if summaryHour != nil {
		hourStr = fmt.Sprintf("%d", *summaryHour)
	}

	logs.CtxDebugf(ctx,
		"[TokenMetering] Summary updated: tenant=%s date=%s hour=%s tokens=%d cost=%.4f duration=%dms",
		tenantID,
		summaryDate,
		hourStr,
		totalTokens,
		totalCost,
		duration.Milliseconds(),
	)
}

// LogSummaryUpdateError 记录汇总更新错误日志
func LogSummaryUpdateError(ctx context.Context, tenantID, summaryDate string, err error) {
	logs.CtxErrorf(ctx,
		"[TokenMetering] Summary update failed: tenant=%s date=%s error=%v",
		tenantID,
		summaryDate,
		err,
	)
}

// LogUsageStatsRetrieval 记录使用统计查询日志
func LogUsageStatsRetrieval(ctx context.Context, tenantID, startDate, endDate string, totalTokens int64, totalCost float64, duration time.Duration) {
	logs.CtxInfof(ctx,
		"[TokenMetering] Usage stats retrieved: tenant=%s period=%s~%s tokens=%d cost=%.2f duration=%dms",
		tenantID,
		startDate,
		endDate,
		totalTokens,
		totalCost,
		duration.Milliseconds(),
	)
}

// LogUsageStatsRetrievalError 记录使用统计查询错误日志
func LogUsageStatsRetrievalError(ctx context.Context, tenantID string, err error) {
	logs.CtxErrorf(ctx,
		"[TokenMetering] Usage stats retrieval failed: tenant=%s error=%v",
		tenantID,
		err,
	)
}

// LogModelUsageStats 记录模型使用统计日志
func LogModelUsageStats(ctx context.Context, tenantID string, modelCount int, startDate, endDate string, duration time.Duration) {
	logs.CtxInfof(ctx,
		"[TokenMetering] Model usage stats retrieved: tenant=%s models=%d period=%s~%s duration=%dms",
		tenantID,
		modelCount,
		startDate,
		endDate,
		duration.Milliseconds(),
	)
}

// LogDailyUsageStats 记录每日使用统计日志
func LogDailyUsageStats(ctx context.Context, tenantID string, daysCount int, startDate, endDate string) {
	logs.CtxInfof(ctx,
		"[TokenMetering] Daily usage stats retrieved: tenant=%s days=%d period=%s~%s",
		tenantID,
		daysCount,
		startDate,
		endDate,
	)
}

// LogDBOperation 记录数据库操作日志
func LogDBOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
	if err != nil {
		logs.CtxErrorf(ctx,
			"[DB] Operation failed: operation=%s table=%s duration=%dms error=%v",
			operation,
			table,
			duration.Milliseconds(),
			err,
		)
	} else {
		logs.CtxDebugf(ctx,
			"[DB] Operation succeeded: operation=%s table=%s duration=%dms",
			operation,
			table,
			duration.Milliseconds(),
		)
	}
}

// LogSlowDBOperation 记录慢数据库操作日志
func LogSlowDBOperation(ctx context.Context, operation, table string, duration time.Duration) {
	logs.CtxWarnf(ctx,
		"[DB] Slow operation detected: operation=%s table=%s duration=%dms threshold=100ms",
		operation,
		table,
		duration.Milliseconds(),
	)
}

// LogBatchBudgetCheck 记录批量预算检查日志
func LogBatchBudgetCheck(ctx context.Context, totalTenants, successCount, failureCount int, duration time.Duration) {
	logs.CtxInfof(ctx,
		"[BudgetAlert] Batch check completed: total=%d success=%d failed=%d duration=%dms",
		totalTenants,
		successCount,
		failureCount,
		duration.Milliseconds(),
	)
}

// LogCacheHit 记录缓存命中日志
func LogCacheHit(ctx context.Context, cacheType, key string) {
	logs.CtxDebugf(ctx,
		"[Cache] Hit: type=%s key=%s",
		cacheType,
		key,
	)
}

// LogCacheMiss 记录缓存未命中日志
func LogCacheMiss(ctx context.Context, cacheType, key string) {
	logs.CtxDebugf(ctx,
		"[Cache] Miss: type=%s key=%s",
		cacheType,
		key,
	)
}

// LogAsyncTaskStart 记录异步任务启动日志
func LogAsyncTaskStart(ctx context.Context, taskName string) {
	logs.CtxInfof(ctx,
		"[AsyncTask] Task started: name=%s",
		taskName,
	)
}

// LogAsyncTaskComplete 记录异步任务完成日志
func LogAsyncTaskComplete(ctx context.Context, taskName string, duration time.Duration, err error) {
	if err != nil {
		logs.CtxErrorf(ctx,
			"[AsyncTask] Task failed: name=%s duration=%dms error=%v",
			taskName,
			duration.Milliseconds(),
			err,
		)
	} else {
		logs.CtxInfof(ctx,
			"[AsyncTask] Task completed: name=%s duration=%dms",
			taskName,
			duration.Milliseconds(),
		)
	}
}

// ============================================================
// 性能日志辅助函数
// ============================================================

// LogPerformanceWarning 记录性能警告
func LogPerformanceWarning(ctx context.Context, operation string, duration time.Duration, threshold time.Duration) {
	logs.CtxWarnf(ctx,
		"[Performance] Slow operation: operation=%s duration=%dms threshold=%dms",
		operation,
		duration.Milliseconds(),
		threshold.Milliseconds(),
	)
}

// LogBusinessMetrics 记录业务指标日志
func LogBusinessMetrics(ctx context.Context, metricName string, value float64, labels map[string]string) {
	labelStr := ""
	for k, v := range labels {
		labelStr += fmt.Sprintf("%s=%s ", k, v)
	}

	logs.CtxInfof(ctx,
		"[BusinessMetric] %s: value=%.2f %s",
		metricName,
		value,
		labelStr,
	)
}
