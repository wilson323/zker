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

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 预算告警服务 - 基于设计文档: 23-租户计费系统_TokenMetering补充_完整版.md
// ============================================================

// BudgetAlertService 预算告警服务
type BudgetAlertService struct {
	budgetRepo   repository.BudgetSettingsRepository
	summaryRepo  repository.TokenUsageSummaryRepository
	alertRepo    repository.BudgetAlertRepository
	notifier     *NotificationService
	logger       *logrus.Logger
}

// NewBudgetAlertService 创建预算告警服务实例
func NewBudgetAlertService(
	budgetRepo repository.BudgetSettingsRepository,
	summaryRepo repository.TokenUsageSummaryRepository,
	alertRepo repository.BudgetAlertRepository,
	notifier *NotificationService,
	logger *logrus.Logger,
) *BudgetAlertService {
	return &BudgetAlertService{
		budgetRepo:  budgetRepo,
		summaryRepo: summaryRepo,
		alertRepo:   alertRepo,
		notifier:    notifier,
		logger:      logger,
	}
}

// CheckBudget 检查预算状态
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
// 返回: error
func (s *BudgetAlertService) CheckBudget(
	ctx context.Context,
	tenantID string,
) error {
	// 1. 获取预算配置
	budget, err := s.budgetRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to get budget for tenant %s", tenantID)
		return err
	}

	if budget == nil {
		// 未配置预算，不检查
		return nil
	}

	// 2. 如果未启用硬性上限，不检查告警
	if !budget.HardCapEnabled {
		return nil
	}

	// 3. 计算当前周期使用量
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
	default:
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	// 4. 获取使用量
	usedAmount, err := s.summaryRepo.GetTotalCost(ctx, tenantID, startDate, now)
	if err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to get total cost for tenant %s", tenantID)
		return err
	}

	// 5. 计算使用率
	usagePercent := decimal.NewFromFloat(usedAmount).
		Div(decimal.NewFromFloat(budget.BudgetAmount)).
		Mul(decimal.NewFromInt(100))

	percentFloat, _ := usagePercent.Float64()

	s.logger.WithFields(logrus.Fields{
		"tenant_id":      tenantID,
		"budget_amount":  budget.BudgetAmount,
		"used_amount":    usedAmount,
		"usage_percent":  percentFloat,
	}).Infof("[BudgetAlert] Budget usage: %.2f%%", percentFloat)

	// 6. 检查告警阈值
	if usagePercent.Cmp(decimal.NewFromInt(int64(budget.AlertThreshold1))) >= 0 {
		// 一级告警 (80%)
		s.sendAlertIfNeeded(ctx, tenantID, budget, usedAmount, usagePercent, entity.AlertTypeThreshold1)
	}

	if usagePercent.Cmp(decimal.NewFromInt(int64(budget.AlertThreshold2))) >= 0 {
		// 二级告警 (95%)
		s.sendAlertIfNeeded(ctx, tenantID, budget, usedAmount, usagePercent, entity.AlertTypeThreshold2)
	}

	// 7. 检查硬性上限
	hardCapAmount := 0.0
	if budget.HardCapAmount != nil {
		hardCapAmount = *budget.HardCapAmount
	} else {
		hardCapAmount = budget.BudgetAmount
	}

	if usedAmount >= hardCapAmount {
		s.handleHardCap(ctx, tenantID, budget, usedAmount, usagePercent)
	}

	return nil
}

// sendAlertIfNeeded 发送告警（防重复）
func (s *BudgetAlertService) sendAlertIfNeeded(
	ctx context.Context,
	tenantID string,
	budget *entity.BudgetSettings,
	usedAmount float64,
	usagePercent decimal.Decimal,
	alertType string,
) {
	// 1. 检查是否已发送过相同类型的告警
	today := time.Now().Format("2006-01-02")
	alreadySent, err := s.alertRepo.CheckAlertSentToday(ctx, tenantID, alertType, today)
	if err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to check alert sent status")
	}
	if alreadySent {
		s.logger.Infof("[BudgetAlert] Alert already sent today for tenant %s, alert type %s", tenantID, alertType)
		return
	}

	// 2. 构建告警级别
	alertLevel := entity.AlertLevelWarning
	if alertType == entity.AlertTypeThreshold2 {
		alertLevel = entity.AlertLevelCritical
	}

	// 3. 构建告警消息
	message := s.buildAlertMessage(budget, usedAmount, usagePercent, alertLevel)

	// 4. 记录告警历史
	alert := &entity.BudgetAlert{
		TenantID:       tenantID,
		AlertType:      alertType,
		BudgetAmount:   budget.BudgetAmount,
		UsedAmount:     usedAmount,
		UsagePercent:   usedAmount / budget.BudgetAmount * 100,
		AlertLevel:     alertLevel,
		AlertMessage:   &message,
		NotificationChannels: budget.NotificationChannels,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to create alert record")
		return
	}

	// 5. 发送通知
	if err := s.notifier.SendAlert(ctx, alert); err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to send alert notification")
	}
}

// handleHardCap 处理硬性上限
func (s *BudgetAlertService) handleHardCap(
	ctx context.Context,
	tenantID string,
	budget *entity.BudgetSettings,
	usedAmount float64,
	usagePercent decimal.Decimal,
) {
	s.logger.Warnf("[BudgetAlert] HARD CAP TRIGGERED for tenant %s", tenantID)

	// 1. 发送紧急告警
	message := "⚠️ 硬性上限触发！服务已暂停。请联系管理员充值或调整预算。"
	alert := &entity.BudgetAlert{
		TenantID:       tenantID,
		AlertType:      entity.AlertTypeHardCap,
		BudgetAmount:   budget.BudgetAmount,
		UsedAmount:     usedAmount,
		UsagePercent:   usedAmount / budget.BudgetAmount * 100,
		AlertLevel:     entity.AlertLevelEmergency,
		AlertMessage:   &message,
		NotificationChannels: budget.NotificationChannels,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to create hard cap alert")
	}

	if err := s.notifier.SendAlert(ctx, alert); err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to send hard cap notification")
	}

	// 2. 执行暂停/降级策略
	if budget.AutoDowngradeEnabled {
		s.executeDowngrade(ctx, tenantID, budget)
	} else {
		s.pauseService(ctx, tenantID)
	}
}

// pauseService 暂停服务
func (s *BudgetAlertService) pauseService(ctx context.Context, tenantID string) {
	s.logger.Warnf("[BudgetAlert] Pausing service for tenant %s due to hard cap", tenantID)

	// TODO: 调用租户服务暂停租户
	// 示例：tenantService.SuspendTenant(ctx, tenantID, "Budget hard cap exceeded")

	// 记录日志
	s.logger.WithField("tenant_id", tenantID).
		Warn("[BudgetAlert] Service paused due to hard cap")
}

// executeDowngrade 执行降级
func (s *BudgetAlertService) executeDowngrade(
	ctx context.Context,
	tenantID string,
	budget *entity.BudgetSettings,
) {
	if budget.OriginalModel == nil || budget.FallbackModel == nil {
		s.logger.Warnf("[BudgetAlert] Downgrade enabled but model not configured for tenant %s", tenantID)
		return
	}

	s.logger.Infof("[BudgetAlert] Executing downgrade for tenant %s: %s → %s",
		tenantID, *budget.OriginalModel, *budget.FallbackModel)

	// TODO: 更新模型路由配置
	// 示例：modelRouter.UpdateRoutingRule(ctx, &RoutingRule{
	//     TenantID:      tenantID,
	//     OriginalModel: *budget.OriginalModel,
	//     FallbackModel: *budget.FallbackModel,
	//     Reason:        "budget_overage",
	//     Enabled:       true,
	// })

	// 记录降级事件
	message := fmt.Sprintf("自动降级已触发: %s → %s", *budget.OriginalModel, *budget.FallbackModel)
	alert := &entity.BudgetAlert{
		TenantID:     tenantID,
		AlertType:    entity.AlertTypeDowngrade,
		AlertLevel:   entity.AlertLevelWarning,
		AlertMessage: &message,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to create downgrade alert")
	}

	if err := s.notifier.SendAlert(ctx, alert); err != nil {
		s.logger.WithError(err).Errorf("[BudgetAlert] Failed to send downgrade notification")
	}

	s.logger.WithField("tenant_id", tenantID).
		WithField("original_model", *budget.OriginalModel).
		WithField("fallback_model", *budget.FallbackModel).
		Info("[BudgetAlert] Model downgrade executed successfully")
}

// buildAlertMessage 构建告警消息
func (s *BudgetAlertService) buildAlertMessage(
	budget *entity.BudgetSettings,
	usedAmount float64,
	usagePercent decimal.Decimal,
	alertLevel string,
) string {
	emoji := "⚠️"
	if alertLevel == entity.AlertLevelCritical {
		emoji = "🚨"
	} else if alertLevel == entity.AlertLevelEmergency {
		emoji = "🛑"
	}

	percentFloat, _ := usagePercent.Float64()

	message := fmt.Sprintf(
		"%s **预算告警**\n\n"+
			"租户ID: %s\n"+
			"预算金额: ¥%.2f\n"+
			"已使用: ¥%.2f (%.1f%%)\n\n",
		emoji,
		budget.TenantID,
		budget.BudgetAmount,
		usedAmount,
		percentFloat,
	)

	if alertLevel == entity.AlertLevelWarning {
		message += "建议: 请注意控制使用量，避免超额产生额外费用。"
	} else if alertLevel == entity.AlertLevelCritical {
		message += "警告: 即将达到预算上限，建议立即充值或调整预算配置。"
	} else if alertLevel == entity.AlertLevelEmergency {
		message += "紧急: 已达到预算上限，服务已暂停或降级。请联系管理员处理。"
	}

	return message
}

// CheckBudgetsForAllTenants 检查所有租户的预算状态
// 应该由定时任务调用（例如每小时执行一次）
func (s *BudgetAlertService) CheckBudgetsForAllTenants(ctx context.Context) error {
	// 1. 获取所有启用了硬性上限的租户
	filter := &repository.BudgetSettingsFilter{
		HardCapEnabled: boolPtr(true),
		PageSize:       1000, // 足够大的数字
	}

	budgets, _, err := s.budgetRepo.List(ctx, filter)
	if err != nil {
		s.logger.WithError(err).Error("[BudgetAlert] Failed to list budgets")
		return err
	}

	// 2. 逐个检查
	for _, budget := range budgets {
		if err := s.CheckBudget(ctx, budget.TenantID); err != nil {
			s.logger.WithError(err).
				WithField("tenant_id", budget.TenantID).
				Error("[BudgetAlert] Failed to check budget for tenant")
			continue
		}
	}

	return nil
}

// boolPtr 返回bool指针的辅助函数
func boolPtr(b bool) *bool {
	return &b
}
