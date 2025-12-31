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

// EmailService 邮件服务接口
// 用于发送验证码、通知等邮件
//
// 使用示例:
//
//	emailSvc := service.NewEmailService()
//	err := emailSvc.SendVerificationCode(ctx, "user@example.com", "123456")
type EmailService struct {
	enabled      bool
	fromAddress  string
	fromName     string
	// 可扩展：添加SMTP客户端配置
}

// NewEmailService 创建邮件服务实例
func NewEmailService(enabled bool, fromAddress, fromName string) *EmailService {
	return &EmailService{
		enabled:     enabled,
		fromAddress: fromAddress,
		fromName:    fromName,
	}
}

// SendVerificationCode 发送验证码邮件
//
// 参数:
//   - ctx: 上下文
//   - to: 收件人邮箱地址
//   - code: 6位数字验证码
//
// 返回:
//   - error: 发送失败时返回错误
func (s *EmailService) SendVerificationCode(ctx context.Context, to, code string) error {
	if !s.enabled {
		logs.Warnf("Email service is disabled, skipping send to %s", to)
		return nil
	}

	// 构建邮件内容
	_ = "【Coze Studio】邮箱验证码" // subject (TODO: 实际发送时使用)
	_ = s.buildVerificationCodeEmail(code) // body (TODO: 实际发送时使用)

	// 记录日志（实际发送可接入SMTP服务如：阿里云邮件、SendGrid）
	logs.Infof("Sending verification code email to %s: code=%s", to, code)

	// TODO: 接入实际邮件发送服务
	// 示例：使用阿里云邮件、SendGrid、腾讯云邮件等
	// 1. 构建SMTP请求
	// 2. 调用邮件服务API
	// 3. 处理发送结果

	// 模拟发送延迟
	time.Sleep(50 * time.Millisecond)

	return nil
}

// buildVerificationCodeEmail 构建验证码邮件内容
func (s *EmailService) buildVerificationCodeEmail(code string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #1890ff; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 8px 8px; }
        .code { font-size: 32px; font-weight: bold; color: #1890ff; text-align: center; padding: 20px; background: white; border-radius: 4px; margin: 20px 0; letter-spacing: 5px; }
        .footer { text-align: center; color: #999; font-size: 12px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Coze Studio</h2>
        </div>
        <div class="content">
            <h3>邮箱验证码</h3>
            <p>您好，</p>
            <p>您正在进行 Coze Studio 租户注册操作，验证码如下：</p>
            <div class="code">%s</div>
            <p><strong>有效期为 5 分钟</strong>，请尽快完成验证。</p>
            <p>如果这不是您的操作，请忽略此邮件。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复。</p>
            <p>&copy; 2025 Coze Studio. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, code)
}

// SendWelcomeEmail 发送欢迎邮件（租户注册成功后）
func (s *EmailService) SendWelcomeEmail(ctx context.Context, to, tenantName string) error {
	if !s.enabled {
		return nil
	}

	_ = "欢迎加入 Coze Studio" // subject (TODO: 实际发送时使用)
	logs.Infof("Sending welcome email to %s for tenant: %s", to, tenantName)

	// TODO: 实现欢迎邮件发送
	return nil
}

// SendPaymentSuccessEmail 发送支付成功邮件
func (s *EmailService) SendPaymentSuccessEmail(ctx context.Context, to, invoiceNumber string, amount float64) error {
	if !s.enabled {
		return nil
	}

	_ = "支付成功确认" // subject (TODO: 实际发送时使用)
	logs.Infof("Sending payment success email to %s: invoice=%s, amount=%.2f", to, invoiceNumber, amount)

	// TODO: 实现支付成功邮件发送
	return nil
}

// SendApprovalEmail 发送审核通过邮件
func (s *EmailService) SendApprovalEmail(ctx context.Context, to, tenantName string) error {
	if !s.enabled {
		return nil
	}

	_ = "租户审核通过通知" // subject (TODO: 实际发送时使用)
	logs.Infof("Sending approval email to %s for tenant: %s", to, tenantName)

	// TODO: 实现审核通过邮件发送
	return nil
}

// SendRejectionEmail 发送审核拒绝邮件
func (s *EmailService) SendRejectionEmail(ctx context.Context, to, tenantName, reason string) error {
	if !s.enabled {
		return nil
	}

	_ = "租户审核拒绝通知" // subject (TODO: 实际发送时使用)
	logs.Infof("Sending rejection email to %s for tenant: %s, reason: %s", to, tenantName, reason)

	// TODO: 实现审核拒绝邮件发送
	return nil
}
