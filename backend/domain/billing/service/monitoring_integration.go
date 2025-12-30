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
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// 监控集成辅助类
// ============================================================

// MonitoringIntegration 监控集成辅助类
type MonitoringIntegration struct {
	service *TokenMeteringService
}

// NewMonitoringIntegration 创建监控集成实例
func NewMonitoringIntegration(service *TokenMeteringService) *MonitoringIntegration {
	return &MonitoringIntegration{
		service: service,
	}
}

// ============================================================
// Token记录监控集成
// ============================================================

// RecordTokenWithMonitoring 集成监控的Token记录
func (m *MonitoringIntegration) RecordTokenWithMonitoring(
	ctx context.Context,
	req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
	// 1. 开始追踪
	spanCtx, span := StartRecordSpan(ctx, req.TenantID, req.ModelProvider, req.ModelName)
	defer span.End()

	startTime := time.Now()

	// 2. 调用原始方法
	result, err := m.service.RecordTokenUsage(spanCtx, req)
	duration := time.Since(startTime)

	// 3. 记录指标
	if err != nil {
		LogTokenRecordError(ctx, req, err)
		SetErrorTag(spanCtx, err)
		return nil, err
	}

	// 4. 记录成功指标
	RecordTokenRecord(
		req.TenantID,
		req.ModelProvider,
		req.ModelName,
		req.InputTokens,
		req.OutputTokens,
		result.TotalCost,
		duration.Seconds(),
		req.IsCached,
	)

	// 5. 设置Span标签
	SetRecordTags(span, req, result)

	// 6. 记录日志
	LogTokenRecord(ctx, req, result, duration)

	// 7. 记录慢操作警告
	if duration > 100*time.Millisecond {
		LogPerformanceWarning(ctx, "TokenRecord", duration, 100*time.Millisecond)
	}

	return result, nil
}

// BatchRecordTokenWithMonitoring 集成监控的批量Token记录
func (m *MonitoringIntegration) BatchRecordTokenWithMonitoring(
	ctx context.Context,
	records []*RecordTokenUsageRequest,
) (*BatchRecordResponse, error) {
	// 1. 开始追踪
	spanCtx, span := StartBatchRecordSpan(ctx, records[0].TenantID, len(records))
	defer span.End()

	startTime := time.Now()

	// 2. 调用原始方法
	result, err := m.service.BatchRecordTokenUsage(spanCtx, records)
	duration := time.Since(startTime)

	// 3. 记录指标
	RecordBatchTokenRecord(
		records[0].TenantID,
		result.SuccessCount,
		result.TotalCost,
		duration.Seconds(),
	)

	// 4. 设置Span标签
	span.SetAttributes(
		attribute.Int("success_count", result.SuccessCount),
		attribute.Int("failed_count", result.FailedCount),
		attribute.Float64("total_cost", result.TotalCost),
	)

	// 5. 记录日志
	if err != nil {
		LogBatchTokenRecordError(ctx, len(records), err)
		SetErrorTag(spanCtx, err)
		return nil, err
	}

	LogBatchTokenRecord(ctx, len(records), result.SuccessCount, result.FailedCount, result.TotalCost, duration)

	return result, nil
}

// ============================================================
// 预算检查监控集成
// ============================================================

// BudgetAlertMonitoringIntegration 预算告警监控集成
type BudgetAlertMonitoringIntegration struct {
	service *BudgetAlertService
}

// NewBudgetAlertMonitoringIntegration 创建预算告警监控集成实例
func NewBudgetAlertMonitoringIntegration(service *BudgetAlertService) *BudgetAlertMonitoringIntegration {
	return &BudgetAlertMonitoringIntegration{
		service: service,
	}
}

// CheckBudgetWithMonitoring 集成监控的预算检查
func (m *BudgetAlertMonitoringIntegration) CheckBudgetWithMonitoring(
	ctx context.Context,
	tenantID string,
) error {
	// 1. 开始追踪
	spanCtx, span := StartBudgetCheckSpan(ctx, tenantID)
	defer span.End()

	startTime := time.Now()

	// 2. 调用原始方法
	err := m.service.CheckBudget(spanCtx, tenantID)
	duration := time.Since(startTime)

	// 3. 获取预算信息（用于记录指标）
	budget, budgetErr := m.service.budgetRepo.GetByTenantID(ctx, tenantID)
	if budgetErr == nil && budget != nil {
		// 计算使用率
		now := time.Now()
		var startDate time.Time
		switch budget.BudgetType {
		case entity.BudgetTypeMonthly:
			startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		case entity.BudgetTypeQuarterly:
			quarter := (now.Month() - 1) / 3
			startDate = time.Date(now.Year(), quarter*3+1, 1, 0, 0, 0, 0, now.Location())
		case entity.BudgetTypeYearly:
			startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		}

		usedAmount, usedErr := m.service.summaryRepo.GetTotalCost(ctx, tenantID, startDate, now)
		if usedErr == nil {
			usagePercent := (usedAmount / budget.BudgetAmount) * 100
			triggered := usedAmount >= budget.BudgetAmount

			// 记录指标
			RecordBudgetCheck(tenantID, triggered, duration.Seconds())
			UpdateBudgetUsage(tenantID, string(budget.BudgetType), getAlertLevel(usagePercent), usagePercent, budget.BudgetAmount, usedAmount)

			// 设置Span标签
			SetBudgetCheckTags(span, budget.BudgetAmount, usedAmount, usagePercent, triggered)

			// 记录日志
			LogBudgetCheck(ctx, tenantID, budget.BudgetAmount, usedAmount, usagePercent, triggered)
		}
	}

	if err != nil {
		LogBudgetCheckError(ctx, tenantID, err)
		SetErrorTag(spanCtx, err)
		return err
	}

	return nil
}

// SendBudgetAlertWithMonitoring 集成监控的预算告警发送
func (m *BudgetAlertMonitoringIntegration) SendBudgetAlertWithMonitoring(
	ctx context.Context,
	tenantID string,
	alertType, alertLevel string,
	usagePercent, budgetAmount, usedAmount float64,
) error {
	// 1. 开始追踪
	_, span := StartAlertSendSpan(ctx, tenantID, "webhook") // 假设使用webhook
	defer span.End()

	startTime := time.Now()

	// 2. 记录指标
	RecordBudgetAlert(tenantID, alertType, alertLevel)

	// 3. 设置Span标签
	SetAlertTags(span, alertType, alertLevel, usagePercent)

	// 4. 记录日志
	LogBudgetAlert(ctx, tenantID, alertType, alertLevel, budgetAmount, usedAmount, usagePercent)

	// 5. 模拟发送通知（实际应该调用NotificationService）
	duration := time.Since(startTime)
	RecordAlertNotification(tenantID, "webhook", "success", duration.Seconds())
	LogAlertNotification(ctx, tenantID, "webhook", true, duration)

	return nil
}

// ============================================================
// 辅助函数
// ============================================================

// getAlertLevel 根据使用率获取告警级别
func getAlertLevel(usagePercent float64) string {
	if usagePercent >= 100 {
		return "exceeded"
	} else if usagePercent >= 95 {
		return "critical"
	} else if usagePercent >= 80 {
		return "warning"
	}
	return "normal"
}

// ============================================================
// 定价引擎监控集成
// ============================================================

// PricingEngineMonitoringIntegration 定价引擎监控集成
type PricingEngineMonitoringIntegration struct {
	engine *PricingEngine
}

// NewPricingEngineMonitoringIntegration 创建定价引擎监控集成实例
func NewPricingEngineMonitoringIntegration(engine *PricingEngine) *PricingEngineMonitoringIntegration {
	return &PricingEngineMonitoringIntegration{
		engine: engine,
	}
}

// CalculateCostWithMonitoring 集成监控的成本计算
func (p *PricingEngineMonitoringIntegration) CalculateCostWithMonitoring(
	ctx context.Context,
	modelProvider, modelName string,
	inputTokens, outputTokens int,
) (*CostBreakdown, error) {
	// 1. 开始追踪
	_, span := StartPricingCalculationSpan(ctx, modelProvider, modelName, inputTokens, outputTokens)
	defer span.End()

	startTime := time.Now()

	// 2. 调用原始方法
	result := p.engine.CalculateCost(modelProvider, modelName, inputTokens, outputTokens)
	duration := time.Since(startTime)

	// 3. 记录指标
	// 提取成本值
	inputCost, _ := result.InputCost.Float64()
	outputCost, _ := result.OutputCost.Float64()
	totalCost, _ := result.TotalCost.Float64()
	unitPrice, _ := result.UnitPrice.Float64()
	
	RecordPricingCalculation(modelProvider, modelName, "success", duration.Seconds(), unitPrice, "input")
	SetPricingTags(span, inputCost, outputCost, totalCost, unitPrice)
	LogPricingCalculation(ctx, modelProvider, modelName, inputTokens, outputTokens, inputCost, outputCost, totalCost)

	return result, nil
}

// ============================================================
// 汇总更新监控集成
// ============================================================

// SummaryUpdateWithMonitoring 集成监控的汇总更新
func SummaryUpdateWithMonitoring(
	ctx context.Context,
	tenantID, summaryDate string,
	summaryHour *uint8,
	totalTokens int64,
	totalCost float64,
	updateFunc func(context.Context) error,
) error {
	// 1. 开始追踪
	spanCtx, span := StartSummaryUpdateSpan(ctx, tenantID, summaryDate)
	defer span.End()

	startTime := time.Now()

	// 2. 调用更新函数
	err := updateFunc(spanCtx)
	duration := time.Since(startTime)

	// 3. 记录指标
	result := "success"
	if err != nil {
		result = "failure"
	}

	RecordSummaryUpdate(tenantID, result, duration.Seconds())

	// 4. 设置Span标签
	span.SetAttributes(
		attribute.Int64("total_tokens", totalTokens),
		attribute.Float64("total_cost", totalCost),
		attribute.String("update_result", result),
	)

	// 5. 记录日志
	if err != nil {
		LogSummaryUpdateError(ctx, tenantID, summaryDate, err)
		SetErrorTag(spanCtx, err)
		return err
	}

	LogSummaryUpdate(ctx, tenantID, summaryDate, summaryHour, totalTokens, totalCost, duration)

	return nil
}
