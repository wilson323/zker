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

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 通知服务 - 基于设计文档: 23-租户计费系统_TokenMetering补充_完整版.md
// ============================================================

// NotificationService 通知服务
type NotificationService struct {
	alertRepo      repository.BudgetAlertRepository
	emailService   *EmailService
	smsService     *SMSService
	webhookService *WebhookService
	logger         *logrus.Logger
	config         *viper.Viper
}

// NewNotificationService 创建通知服务实例
func NewNotificationService(
	alertRepo repository.BudgetAlertRepository,
	logger *logrus.Logger,
	config *viper.Viper,
) *NotificationService {
	return &NotificationService{
		alertRepo:      alertRepo,
		emailService:   NewEmailService(config),
		smsService:     NewSMSService(config),
		webhookService: NewWebhookService(config),
		logger:         logger,
		config:         config,
	}
}

// SendAlert 发送告警通知
// 参数:
//   - ctx: 上下文
//   - alert: 告警记录
// 返回: error
func (s *NotificationService) SendAlert(
	ctx context.Context,
	alert *entity.BudgetAlert,
) error {
	// 1. 解析通知渠道
	var channels []string
	if alert.NotificationChannels != "" {
		if err := json.Unmarshal([]byte(alert.NotificationChannels), &channels); err != nil {
			s.logger.WithError(err).Errorf("[Notification] Failed to parse notification channels: %s", alert.NotificationChannels)
			channels = []string{"email"} // 默认使用邮件
		}
	} else {
		channels = []string{"email"} // 默认使用邮件
	}

	// 2. 并发发送多个渠道
	var lastErr error
	successCount := 0

	for _, channel := range channels {
		var err error
		switch channel {
		case "email":
			err = s.emailService.SendAlertEmail(ctx, alert)
		case "sms":
			err = s.smsService.SendAlertSMS(ctx, alert)
		case "webhook":
			err = s.webhookService.SendWebhook(ctx, alert)
		default:
			s.logger.Warnf("[Notification] Unknown notification channel: %s", channel)
			continue
		}

		if err != nil {
			s.logger.WithError(err).Errorf("[Notification] Failed to send %s notification", channel)
			lastErr = err
		} else {
			successCount++
		}
	}

	// 3. 更新发送状态
	now := time.Now()
	if successCount > 0 {
		// 至少有一个渠道发送成功
		if err := s.alertRepo.UpdateNotificationStatus(ctx, alert.ID, entity.NotificationStatusSent, &now, nil); err != nil {
			s.logger.WithError(err).Error("[Notification] Failed to update notification status to sent")
		}
		return nil
	} else {
		// 所有渠道都失败
		errMsg := fmt.Sprintf("All notification channels failed: %v", lastErr)
		if err := s.alertRepo.UpdateNotificationStatus(ctx, alert.ID, entity.NotificationStatusFailed, &now, &errMsg); err != nil {
			s.logger.WithError(err).Error("[Notification] Failed to update notification status to failed")
		}
		return lastErr
	}
}

// ============================================================
// 邮件通知服务
// ============================================================

// EmailService 邮件通知服务
type EmailService struct {
	config *viper.Viper
}

// NewEmailService 创建邮件服务实例
func NewEmailService(config *viper.Viper) *EmailService {
	return &EmailService{config: config}
}

// SendAlertEmail 发送告警邮件
func (s *EmailService) SendAlertEmail(ctx context.Context, alert *entity.BudgetAlert) error {
	// TODO: 实现邮件发送逻辑
	// 示例：使用 SMTP、SendGrid、阿里云邮件服务等

	subject := fmt.Sprintf("[预算告警] %s - %.1f%%", alert.TenantID, alert.UsagePercent)
	body := s.buildEmailBody(alert)

	// 伪代码：
	// err := smtp.Send(
	//     s.config.GetString("email.smtp_host"),
	//     s.config.GetInt("email.smtp_port"),
	//     s.config.GetString("email.username"),
	//     s.config.GetString("email.password"),
	//     recipient,
	//     subject,
	//     body,
	// )
	// return err

	// 暂时记录日志
	logrus.WithFields(logrus.Fields{
		"tenant_id":     alert.TenantID,
		"alert_type":    alert.AlertType,
		"usage_percent": alert.UsagePercent,
	}).Infof("[Email] Would send alert email: %s\n%s", subject, body)

	return nil
}

// buildEmailBody 构建邮件正文
func (s *EmailService) buildEmailBody(alert *entity.BudgetAlert) string {
	return fmt.Sprintf(`
<html>
<body>
	<h2>⚠️ 预算告警通知</h2>
	<p><strong>租户ID:</strong> %s</p>
	<p><strong>告警类型:</strong> %s</p>
	<p><strong>告警级别:</strong> %s</p>
	<hr/>
	<p><strong>预算金额:</strong> ¥%.2f</p>
	<p><strong>已使用:</strong> ¥%.2f</p>
	<p><strong>使用率:</strong> %.1f%%</p>
	<hr/>
	<p>%s</p>
	<p style="color: gray; font-size: 12px;">发送时间: %s</p>
</body>
</html>
`,
		alert.TenantID,
		alert.AlertType,
		alert.AlertLevel,
		alert.BudgetAmount,
		alert.UsedAmount,
		alert.UsagePercent,
		*alert.AlertMessage,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

// ============================================================
// 短信通知服务
// ============================================================

// SMSService 短信通知服务
type SMSService struct {
	config *viper.Viper
}

// NewSMSService 创建短信服务实例
func NewSMSService(config *viper.Viper) *SMSService {
	return &SMSService{config: config}
}

// SendAlertSMS 发送告警短信
func (s *SMSService) SendAlertSMS(ctx context.Context, alert *entity.BudgetAlert) error {
	// TODO: 实现短信发送逻辑
	// 示例：使用阿里云短信、腾讯云短信等服务

	message := s.buildSMSMessage(alert)

	// 伪代码：
	// err := aliyunSms.Send(
	//     s.config.GetString("sms.access_key"),
	//     s.config.GetString("sms.access_secret"),
	//     phoneNumber,
	//     signName,
	//     templateCode,
	//     message,
	// )
	// return err

	// 暂时记录日志
	logrus.WithFields(logrus.Fields{
		"tenant_id":     alert.TenantID,
		"alert_type":    alert.AlertType,
		"usage_percent": alert.UsagePercent,
	}).Infof("[SMS] Would send alert SMS: %s", message)

	return nil
}

// buildSMSMessage 构建短信内容
func (s *SMSService) buildSMSMessage(alert *entity.BudgetAlert) string {
	return fmt.Sprintf(
		"【预算告警】租户%s已使用¥%.2f（%.1f%%），预算¥%.2f。%s",
		alert.TenantID,
		alert.UsedAmount,
		alert.UsagePercent,
		alert.BudgetAmount,
		*alert.AlertMessage,
	)
}

// ============================================================
// Webhook通知服务
// ============================================================

// WebhookService Webhook通知服务
type WebhookService struct {
	config *viper.Viper
	client *HTTPClient
}

// NewWebhookService 创建Webhook服务实例
func NewWebhookService(config *viper.Viper) *WebhookService {
	return &WebhookService{
		config: config,
		client: NewHTTPClient(),
	}
}

// SendWebhook 发送Webhook通知
func (s *WebhookService) SendWebhook(ctx context.Context, alert *entity.BudgetAlert) error {
	// TODO: 实现Webhook发送逻辑
	// 支持多种Webhook类型：钉钉、企业微信、飞书、Slack等

	webhookURL := s.config.GetString("webhook.url")
	if webhookURL == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	payload := s.buildWebhookPayload(alert)

	// 伪代码：
	// err := s.client.Post(ctx, webhookURL, payload, nil)
	// return err

	// 暂时记录日志
	payloadJSON, _ := json.Marshal(payload)
	logrus.WithFields(logrus.Fields{
		"tenant_id":     alert.TenantID,
		"alert_type":    alert.AlertType,
		"usage_percent": alert.UsagePercent,
	}).Infof("[Webhook] Would send webhook to %s: %s", webhookURL, string(payloadJSON))

	return nil
}

// buildWebhookPayload 构建Webhook负载
func (s *WebhookService) buildWebhookPayload(alert *entity.BudgetAlert) map[string]interface{} {
	return map[string]interface{}{
		"msg_type": "text",
		"text": map[string]interface{}{
			"content": fmt.Sprintf(
				"⚠️ 预算告警\n\n租户ID: %s\n预算金额: ¥%.2f\n已使用: ¥%.2f (%.1f%%)\n\n%s",
				alert.TenantID,
				alert.BudgetAmount,
				alert.UsedAmount,
				alert.UsagePercent,
				*alert.AlertMessage,
			),
		},
		"tenant_id":     alert.TenantID,
		"alert_type":    alert.AlertType,
		"alert_level":   alert.AlertLevel,
		"budget_amount": alert.BudgetAmount,
		"used_amount":   alert.UsedAmount,
		"usage_percent": alert.UsagePercent,
		"timestamp":     time.Now().Unix(),
	}
}

// ============================================================
// HTTP客户端（用于Webhook）
// ============================================================

// HTTPClient HTTP客户端
type HTTPClient struct{}

// NewHTTPClient 创建HTTP客户端实例
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// Post 发送POST请求
func (c *HTTPClient) Post(ctx context.Context, url string, payload map[string]interface{}, headers map[string]string) error {
	// TODO: 实现HTTP POST请求
	// 示例：使用 http.Client 或 resty.Client
	return nil
}

// ============================================================
// 辅助函数
// ============================================================

// GetNotificationRecipients 获取通知接收人列表
func GetNotificationRecipients(tenantID string) ([]string, error) {
	// TODO: 从数据库或配置中获取租户的通知接收人
	// 示例：从 tenant_settings 表中查询 notification_recipients 字段
	return []string{}, nil
}

// GetTenantAdminEmails 获取租户管理员邮箱列表
func GetTenantAdminEmails(tenantID string) ([]string, error) {
	// TODO: 从数据库中查询租户管理员的邮箱
	// 示例：
	// SELECT u.email FROM users u
	// JOIN user_roles ur ON u.user_id = ur.user_id
	// JOIN roles r ON ur.role_id = r.role_id
	// WHERE u.tenant_id = ? AND r.role_code = 'tenant_admin'
	return []string{}, nil
}

// GetTenantAdminPhones 获取租户管理员手机号列表
func GetTenantAdminPhones(tenantID string) ([]string, error) {
	// TODO: 从数据库中查询租户管理员的手机号
	return []string{}, nil
}
