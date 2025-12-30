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
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
	"github.com/coze-dev/coze-studio/backend/domain/monitoring/repository"
)

// AlertService 告警服务
type AlertService struct {
	alertRuleRepo    repository.AlertRuleRepository
	alertHistoryRepo repository.AlertHistoryRepository
	metricsRepo      repository.TenantMetricsRepository
	notificationSvc  *NotificationService
	logger           *zap.Logger
}

// NewAlertService 创建告警服务实例
func NewAlertService(
	alertRuleRepo repository.AlertRuleRepository,
	alertHistoryRepo repository.AlertHistoryRepository,
	metricsRepo repository.TenantMetricsRepository,
	notificationSvc *NotificationService,
	logger *zap.Logger,
) *AlertService {
	return &AlertService{
		alertRuleRepo:    alertRuleRepo,
		alertHistoryRepo: alertHistoryRepo,
		metricsRepo:      metricsRepo,
		notificationSvc:  notificationSvc,
		logger:           logger,
	}
}

// CreateAlertRule 创建告警规则
func (s *AlertService) CreateAlertRule(ctx context.Context, rule *entity.AlertRule) error {
	s.logger.Info("Creating alert rule",
		zap.String("tenant_id", rule.TenantID),
		zap.String("rule_name", rule.RuleName),
		zap.String("metric_type", rule.MetricType),
	)

	// 验证规则
	if err := s.validateAlertRule(rule); err != nil {
		return fmt.Errorf("invalid alert rule: %w", err)
	}

	err := s.alertRuleRepo.Create(ctx, rule)
	if err != nil {
		s.logger.Error("Failed to create alert rule",
			zap.String("rule_name", rule.RuleName),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create alert rule: %w", err)
	}

	return nil
}

// UpdateAlertRule 更新告警规则
func (s *AlertService) UpdateAlertRule(ctx context.Context, rule *entity.AlertRule) error {
	s.logger.Info("Updating alert rule",
		zap.String("rule_id", rule.ID),
	)

	if err := s.validateAlertRule(rule); err != nil {
		return fmt.Errorf("invalid alert rule: %w", err)
	}

	err := s.alertRuleRepo.Update(ctx, rule)
	if err != nil {
		s.logger.Error("Failed to update alert rule",
			zap.String("rule_id", rule.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update alert rule: %w", err)
	}

	return nil
}

// DeleteAlertRule 删除告警规则
func (s *AlertService) DeleteAlertRule(ctx context.Context, ruleID string) error {
	s.logger.Info("Deleting alert rule",
		zap.String("rule_id", ruleID),
	)

	err := s.alertRuleRepo.Delete(ctx, ruleID)
	if err != nil {
		s.logger.Error("Failed to delete alert rule",
			zap.String("rule_id", ruleID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}

	return nil
}

// CheckAlerts 检查告警规则（定时任务）
func (s *AlertService) CheckAlerts(ctx context.Context) ([]*entity.AlertHistory, error) {
	s.logger.Info("Checking alert rules")

	// 获取所有启用的告警规则
	rules, err := s.alertRuleRepo.GetEnabledRules(ctx)
	if err != nil {
		s.logger.Error("Failed to get enabled alert rules", zap.Error(err))
		return nil, fmt.Errorf("failed to get enabled alert rules: %w", err)
	}

	triggeredAlerts := make([]*entity.AlertHistory, 0)

	// 检查每个规则
	for _, rule := range rules {
		alert, err := s.checkRule(ctx, rule)
		if err != nil {
			s.logger.Error("Failed to check alert rule",
				zap.String("rule_id", rule.ID),
				zap.Error(err),
			)
			continue
		}

		if alert != nil {
			triggeredAlerts = append(triggeredAlerts, alert)

			// 发送告警通知
			if err := s.SendAlert(ctx, alert, rule); err != nil {
				s.logger.Error("Failed to send alert notification",
					zap.String("alert_id", alert.ID),
					zap.Error(err),
				)
			}
		}
	}

	s.logger.Info("Alert check completed",
		zap.Int("triggered_alerts", len(triggeredAlerts)),
	)

	return triggeredAlerts, nil
}

// checkRule 检查单个告警规则
func (s *AlertService) checkRule(ctx context.Context, rule *entity.AlertRule) (*entity.AlertHistory, error) {
	// 获取最新指标
	var metricType entity.MetricType
	switch rule.MetricType {
	case "qps":
		metricType = entity.MetricTypeQPS
	case "response_time":
		metricType = entity.MetricTypeResponseTime
	case "error_rate":
		metricType = entity.MetricTypeErrorRate
	default:
		return nil, fmt.Errorf("unsupported metric type: %s", rule.MetricType)
	}

	latestMetric, err := s.metricsRepo.GetByTenantAndType(ctx, rule.TenantID, metricType)
	if err != nil {
		return nil, err
	}

	if latestMetric == nil {
		// 没有指标数据
		return nil, nil
	}

	// 判断是否触发告警
	if !rule.ShouldAlert(latestMetric.MetricValue) {
		return nil, nil
	}

	// 检查是否已存在未解决的告警
	pendingAlerts, err := s.alertHistoryRepo.GetPendingAlertsByTenant(ctx, rule.TenantID)
	if err != nil {
		return nil, err
	}

	// 如果已有相同规则的未解决告警，则不重复创建
	for _, alert := range pendingAlerts {
		if alert.AlertRuleID == rule.ID && !alert.IsSilenced() {
			return nil, nil
		}
	}

	// 创建告警历史记录
	alertData := map[string]interface{}{
		"metric_value":  latestMetric.MetricValue,
		"threshold":     rule.Threshold,
		"comparison":    string(rule.Comparison),
		"metric_type":   rule.MetricType,
		"timestamp":     latestMetric.MetricTimestamp,
	}

	alertDataJSON, _ := json.Marshal(alertData)
	alertMessage := fmt.Sprintf("告警规则 [%s] 触发: 当前值 %.2f, 阈值 %.2f",
		rule.RuleName, latestMetric.MetricValue, rule.Threshold)

	alert := entity.NewAlertHistory(
		rule.TenantID,
		rule.ID,
		rule.RuleName,
		rule.Severity,
		alertMessage,
		string(alertDataJSON),
	)

	err = s.alertHistoryRepo.Create(ctx, alert)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert history: %w", err)
	}

	s.logger.Warn("Alert triggered",
		zap.String("tenant_id", rule.TenantID),
		zap.String("rule_name", rule.RuleName),
		zap.String("severity", string(rule.Severity)),
		zap.Float64("value", latestMetric.MetricValue),
		zap.Float64("threshold", rule.Threshold),
	)

	return alert, nil
}

// SendAlert 发送告警通知
func (s *AlertService) SendAlert(ctx context.Context, alert *entity.AlertHistory, rule *entity.AlertRule) error {
	s.logger.Info("Sending alert notification",
		zap.String("alert_id", alert.ID),
		zap.String("severity", string(alert.Severity)),
	)

	// 根据通知渠道发送通知
	for _, channel := range rule.NotificationChannels {
		switch channel {
		case entity.NotificationChannelEmail:
			if err := s.notificationSvc.SendEmailNotification(ctx, alert, rule); err != nil {
				s.logger.Error("Failed to send email notification",
					zap.String("alert_id", alert.ID),
					zap.Error(err),
				)
			}
		case entity.NotificationChannelSMS:
			if err := s.notificationSvc.SendSMSNotification(ctx, alert, rule); err != nil {
				s.logger.Error("Failed to send SMS notification",
					zap.String("alert_id", alert.ID),
					zap.Error(err),
				)
			}
		case entity.NotificationChannelWebhook:
			if err := s.notificationSvc.SendWebhookNotification(ctx, alert, rule); err != nil {
				s.logger.Error("Failed to send webhook notification",
					zap.String("alert_id", alert.ID),
					zap.Error(err),
				)
			}
		default:
			s.logger.Warn("Unsupported notification channel",
				zap.String("channel", string(channel)),
			)
		}
	}

	return nil
}

// AcknowledgeAlert 确认告警
func (s *AlertService) AcknowledgeAlert(ctx context.Context, alertID, userID string) error {
	s.logger.Info("Acknowledging alert",
		zap.String("alert_id", alertID),
		zap.String("user_id", userID),
	)

	err := s.alertHistoryRepo.Acknowledge(ctx, alertID, userID)
	if err != nil {
		s.logger.Error("Failed to acknowledge alert",
			zap.String("alert_id", alertID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}

// ResolveAlert 解决告警
func (s *AlertService) ResolveAlert(ctx context.Context, alertID, userID, note string) error {
	s.logger.Info("Resolving alert",
		zap.String("alert_id", alertID),
		zap.String("user_id", userID),
	)

	err := s.alertHistoryRepo.Resolve(ctx, alertID, userID, note)
	if err != nil {
		s.logger.Error("Failed to resolve alert",
			zap.String("alert_id", alertID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to resolve alert: %w", err)
	}

	return nil
}

// SilenceAlert 静默告警
func (s *AlertService) SilenceAlert(ctx context.Context, alertID string, duration time.Duration) error {
	s.logger.Info("Silencing alert",
		zap.String("alert_id", alertID),
		zap.Duration("duration", duration),
	)

	err := s.alertHistoryRepo.Silence(ctx, alertID, duration)
	if err != nil {
		s.logger.Error("Failed to silence alert",
			zap.String("alert_id", alertID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to silence alert: %w", err)
	}

	return nil
}

// GetAlertHistory 获取告警历史
func (s *AlertService) GetAlertHistory(
	ctx context.Context,
	filter *entity.AlertFilter,
) ([]*entity.AlertHistory, int64, error) {
	s.logger.Info("Getting alert history",
		zap.String("tenant_id", filter.TenantID),
	)

	history, total, err := s.alertHistoryRepo.GetByTenant(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to get alert history",
			zap.String("tenant_id", filter.TenantID),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get alert history: %w", err)
	}

	return history, total, nil
}

// GetAlertStatistics 获取告警统计信息
func (s *AlertService) GetAlertStatistics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*repository.AlertStatistics, error) {
	s.logger.Info("Getting alert statistics",
		zap.String("tenant_id", tenantID),
	)

	stats, err := s.alertHistoryRepo.GetAlertStatistics(ctx, tenantID, startTime, endTime)
	if err != nil {
		s.logger.Error("Failed to get alert statistics",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get alert statistics: %w", err)
	}

	return stats, nil
}

// ============================================================
// 私有方法
// ============================================================

// validateAlertRule 验证告警规则
func (s *AlertService) validateAlertRule(rule *entity.AlertRule) error {
	if rule.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	if rule.RuleName == "" {
		return fmt.Errorf("rule_name is required")
	}

	if rule.MetricType == "" {
		return fmt.Errorf("metric_type is required")
	}

	if err := entity.ValidateSeverity(rule.Severity); err != nil {
		return err
	}

	if rule.Threshold < 0 {
		return fmt.Errorf("threshold must be non-negative")
	}

	if rule.EvaluationInterval <= 0 {
		return fmt.Errorf("evaluation_interval must be positive")
	}

	if rule.ForDuration < 0 {
		return fmt.Errorf("for_duration must be non-negative")
	}

	return nil
}
