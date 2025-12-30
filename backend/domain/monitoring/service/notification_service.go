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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
)

// NotificationService 通知服务
type NotificationService struct {
	httpClient *http.Client
	logger     *zap.Logger
}

// NewNotificationService 创建通知服务实例
func NewNotificationService(logger *zap.Logger) *NotificationService {
	return &NotificationService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// SendEmailNotification 发送邮件通知
func (s *NotificationService) SendEmailNotification(
	ctx context.Context,
	alert *entity.AlertHistory,
	rule *entity.AlertRule,
) error {
	s.logger.Info("Sending email notification",
		zap.String("alert_id", alert.ID),
		zap.String("tenant_id", alert.TenantID),
	)

	// 解析通知配置
	config := make(map[string]interface{})
	if rule.NotificationConfig != "" {
		if err := json.Unmarshal([]byte(rule.NotificationConfig), &config); err != nil {
			s.logger.Error("Failed to parse notification config",
				zap.String("rule_id", rule.ID),
				zap.Error(err),
			)
			return err
		}
	}

	// 获取收件人地址
	recipients, ok := config["email_recipients"].([]string)
	if !ok || len(recipients) == 0 {
		s.logger.Warn("No email recipients configured",
			zap.String("rule_id", rule.ID),
		)
		return nil
	}

	// 构建邮件内容
	subject := fmt.Sprintf("[%s] 告警通知: %s", alert.Severity, alert.AlertRuleName)
	body := s.buildEmailBody(alert, rule)

	// 调用邮件发送服务（这里简化处理，实际应集成邮件服务）
	s.logger.Info("Email notification prepared",
		zap.String("subject", subject),
		zap.Strings("recipients", recipients),
		zap.String("body", body),
	)

	// TODO: 实际发送邮件（可使用SendGrid、阿里云邮件服务等）
	return nil
}

// SendSMSNotification 发送短信通知
func (s *NotificationService) SendSMSNotification(
	ctx context.Context,
	alert *entity.AlertHistory,
	rule *entity.AlertRule,
) error {
	s.logger.Info("Sending SMS notification",
		zap.String("alert_id", alert.ID),
	)

	// 解析通知配置
	config := make(map[string]interface{})
	if rule.NotificationConfig != "" {
		if err := json.Unmarshal([]byte(rule.NotificationConfig), &config); err != nil {
			return err
		}
	}

	// 获取手机号
	phoneNumbers, ok := config["sms_recipients"].([]string)
	if !ok || len(phoneNumbers) == 0 {
		s.logger.Warn("No SMS recipients configured",
			zap.String("rule_id", rule.ID),
		)
		return nil
	}

	// 构建短信内容
	message := fmt.Sprintf("[%s] 告警: %s - %s",
		alert.Severity, alert.AlertRuleName, alert.AlertMessage)

	// 调用短信发送服务（简化处理）
	s.logger.Info("SMS notification prepared",
		zap.Strings("phone_numbers", phoneNumbers),
		zap.String("message", message),
	)

	// TODO: 实际发送短信（可使用阿里云短信、腾讯云短信等）
	return nil
}

// SendWebhookNotification 发送Webhook通知
func (s *NotificationService) SendWebhookNotification(
	ctx context.Context,
	alert *entity.AlertHistory,
	rule *entity.AlertRule,
) error {
	s.logger.Info("Sending webhook notification",
		zap.String("alert_id", alert.ID),
	)

	// 解析通知配置
	config := make(map[string]interface{})
	if rule.NotificationConfig != "" {
		if err := json.Unmarshal([]byte(rule.NotificationConfig), &config); err != nil {
			return err
		}
	}

	// 获取Webhook URL
	webhookURL, ok := config["webhook_url"].(string)
	if !ok || webhookURL == "" {
		s.logger.Warn("No webhook URL configured",
			zap.String("rule_id", rule.ID),
		)
		return nil
	}

	// 构建Webhook Payload
	payload := map[string]interface{}{
		"alert_id":      alert.ID,
		"tenant_id":     alert.TenantID,
		"rule_name":     alert.AlertRuleName,
		"severity":      string(alert.Severity),
		"status":        string(alert.Status),
		"message":       alert.AlertMessage,
		"alert_data":    alert.AlertData,
		"created_at":    alert.CreatedAt,
		"webhook_type":  "alert",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal webhook payload: %w", err)
	}

	// 发送HTTP POST请求
	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 添加自定义Headers
	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned error status %d: %s", resp.StatusCode, string(body))
	}

	s.logger.Info("Webhook notification sent successfully",
		zap.String("webhook_url", webhookURL),
		zap.Int("status_code", resp.StatusCode),
	)

	return nil
}

// ============================================================
// 私有方法
// ============================================================

// buildEmailBody 构建邮件内容
func (s *NotificationService) buildEmailBody(alert *entity.AlertHistory, rule *entity.AlertRule) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("<h2>告警通知: %s</h2>", alert.AlertRuleName))
	buf.WriteString(fmt.Sprintf("<p><strong>严重级别:</strong> %s</p>", alert.Severity))
	buf.WriteString(fmt.Sprintf("<p><strong>告警消息:</strong> %s</p>", alert.AlertMessage))
	buf.WriteString(fmt.Sprintf("<p><strong>租户ID:</strong> %s</p>", alert.TenantID))
	buf.WriteString(fmt.Sprintf("<p><strong>告警时间:</strong> %s</p>", alert.CreatedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("<p><strong>告警状态:</strong> %s</p>", alert.Status))

	if alert.AlertData != "" {
		buf.WriteString("<h3>告警详情</h3>")
		buf.WriteString(fmt.Sprintf("<pre>%s</pre>", alert.AlertData))
	}

	buf.WriteString("<p>请及时处理该告警。</p>")
	buf.WriteString("<hr>")
	buf.WriteString("<p>此邮件由系统自动发送，请勿回复。</p>")

	return buf.String()
}
