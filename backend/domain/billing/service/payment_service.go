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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingRepository "github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// 支付集成服务 - 基于设计文档: 23-MultiTenant SaaS核心_租户计费系统.md
// ============================================================

// PaymentService 支付服务
type PaymentService struct {
	db                *gorm.DB
	paymentRepo       billingRepository.PaymentRepository
	invoiceRepo       billingRepository.InvoiceRepository
	billingAccountRepo billingRepository.BillingAccountRepository
	logger            *logrus.Logger

	// 支付配置
	alipayConfig      *AlipayConfig
	wechatConfig      *WechatConfig
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID            string `json:"app_id"`
	PrivateKey        string `json:"private_key"`
	AlipayPublicKey   string `json:"alipay_public_key"`
	GatewayURL       string `json:"gateway_url"`
	NotifyURL        string `json:"notify_url"`
	ReturnURL        string `json:"return_url"`
	Sandbox          bool   `json:"sandbox"`
}

// WechatConfig 微信支付配置
type WechatConfig struct {
	AppID       string `json:"app_id"`
	MchID       string `json:"mch_id"`
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	NotifyURL   string `json:"notify_url"`
	CertPath    string `json:"cert_path"`
	KeyPath     string `json:"key_path"`
	Sandbox     bool   `json:"sandbox"`
}

// NewPaymentService 创建支付服务实例
func NewPaymentService(
	db *gorm.DB,
	paymentRepo billingRepository.PaymentRepository,
	invoiceRepo billingRepository.InvoiceRepository,
	billingAccountRepo billingRepository.BillingAccountRepository,
	alipayConfig *AlipayConfig,
	wechatConfig *WechatConfig,
	logger *logrus.Logger,
) *PaymentService {
	return &PaymentService{
		db:                db,
		paymentRepo:       paymentRepo,
		invoiceRepo:       invoiceRepo,
		billingAccountRepo: billingAccountRepo,
		alipayConfig:      alipayConfig,
		wechatConfig:      wechatConfig,
		logger:            logger,
	}
}

// ============================================================
// 请求和响应类型
// ============================================================

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	TenantID        string  `json:"tenant_id" validate:"required"`
	InvoiceID       *uint64 `json:"invoice_id,omitempty"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod   string  `json:"payment_method" validate:"required"` // alipay, wechat
	ClientIP        string  `json:"client_ip" validate:"required"`
	ReturnURL       *string `json:"return_url,omitempty"`
	Description     *string `json:"description,omitempty"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentID      uint64  `json:"payment_id"`
	PaymentNumber  string  `json:"payment_number"`
	PaymentURL     string  `json:"payment_url,omitempty"`
	QRCodeData     string  `json:"qr_code_data,omitempty"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	CreatedAt      int64   `json:"created_at"`
}

// PaymentCallbackRequest 支付回调请求
type PaymentCallbackRequest struct {
	PaymentMethod string                 `json:"payment_method"`
	TradeNo       string                 `json:"trade_no"`
	OutTradeNo    string                 `json:"out_trade_no"`
	TradeStatus   string                 `json:"trade_status"`
	TotalAmount   string                 `json:"total_amount"`
	RawData       map[string]interface{} `json:"-"`
}

// ============================================================
// 核心方法实现
// ============================================================

// CreatePayment 创建支付
// 参数:
//   - ctx: 上下文
//   - req: 创建支付请求
// 返回: 创建支付响应或错误
func (s *PaymentService) CreatePayment(
	ctx context.Context,
	req *CreatePaymentRequest,
) (*CreatePaymentResponse, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordPaymentOperation("create", req.PaymentMethod, duration)
	}()

	s.logger.WithFields(logrus.Fields{
		"tenant_id":      req.TenantID,
		"amount":         req.Amount,
		"payment_method": req.PaymentMethod,
	}).Info("[PaymentService] Creating payment")

	// 1. 验证金额
	if req.Amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// 2. 查询计费账户
	account, err := s.billingAccountRepo.GetByTenant(ctx, req.TenantID)
	if err != nil {
		s.logger.WithError(err).Errorf("[PaymentService] Failed to get billing account")
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	if account == nil {
		return nil, fmt.Errorf("no billing account found for tenant: %s", req.TenantID)
	}

	// 3. 验证发票（如果有）
	if req.InvoiceID != nil {
		invoice, err := s.invoiceRepo.GetByID(ctx, *req.InvoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get invoice: %w", err)
		}

		if invoice == nil {
			return nil, fmt.Errorf("invoice not found: %d", *req.InvoiceID)
		}

		if invoice.TenantID != req.TenantID {
			return nil, fmt.Errorf("invoice does not belong to tenant")
		}

		// 检查支付金额
		remainingAmount := invoice.TotalAmount - invoice.PaidAmount
		if req.Amount > remainingAmount {
			return nil, fmt.Errorf("payment amount exceeds remaining invoice amount")
		}
	}

	// 4. 生成支付单号
	paymentNumber, err := s.generatePaymentNumber(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate payment number: %w", err)
	}

	// 5. 创建支付记录
	payment := &entity.Payment{
		TenantID:         req.TenantID,
		BillingAccountID: account.ID,
		InvoiceID:        req.InvoiceID,
		PaymentNumber:    paymentNumber,
		Amount:           req.Amount,
		Currency:         account.Currency,
		PaymentMethod:    req.PaymentMethod,
		Status:           entity.PaymentStatusPending,
		Provider:         s.getProvider(req.PaymentMethod),
	}

	if req.Description != nil {
		payment.Notes = req.Description
	}

	// 6. 根据支付方式创建支付订单
	var paymentURL, qrCodeData string

	switch req.PaymentMethod {
	case entity.PaymentMethodAlipay:
		paymentURL, qrCodeData, err = s.createAlipayPayment(ctx, payment, req.ClientIP)
	case entity.PaymentMethodWechat:
		paymentURL, qrCodeData, err = s.createWechatPayment(ctx, payment, req.ClientIP)
	default:
		return nil, fmt.Errorf("unsupported payment method: %s", req.PaymentMethod)
	}

	if err != nil {
		s.logger.WithError(err).Errorf("[PaymentService] Failed to create payment with provider")
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 7. 保存支付记录
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		s.logger.WithError(err).Errorf("[PaymentService] Failed to save payment")
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	// 8. 记录日志
	s.logger.WithFields(logrus.Fields{
		"tenant_id":      req.TenantID,
		"payment_id":     payment.ID,
		"payment_number": paymentNumber,
		"amount":         req.Amount,
		"provider":       payment.Provider,
		"creation_time_s": time.Since(startTime).Seconds(),
	}).Info("[PaymentService] Payment created successfully")

	// 9. 更新Prometheus指标
	UpdatePaymentMetrics(req.TenantID, req.PaymentMethod, "created", req.Amount)

	return &CreatePaymentResponse{
		PaymentID:     payment.ID,
		PaymentNumber: paymentNumber,
		PaymentURL:    paymentURL,
		QRCodeData:    qrCodeData,
		Amount:        req.Amount,
		Currency:      payment.Currency,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt.Unix(),
	}, nil
}

// HandlePaymentCallback 处理支付回调
// 参数:
//   - ctx: 上下文
//   - paymentMethod: 支付方式
//   - callbackData: 回调数据
// 返回: 错误
func (s *PaymentService) HandlePaymentCallback(
	ctx context.Context,
	paymentMethod string,
	callbackData map[string]interface{},
) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		RecordPaymentOperation("callback", paymentMethod, duration)
	}()

	s.logger.WithFields(logrus.Fields{
		"payment_method": paymentMethod,
	}).Info("[PaymentService] Handling payment callback")

	// 1. 根据支付方式解析回调
	var callbackReq *PaymentCallbackRequest
	var err error

	switch paymentMethod {
	case entity.PaymentMethodAlipay:
		callbackReq, err = s.parseAlipayCallback(callbackData)
	case entity.PaymentMethodWechat:
		callbackReq, err = s.parseWechatCallback(callbackData)
	default:
		return fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}

	if err != nil {
		s.logger.WithError(err).Errorf("[PaymentService] Failed to parse callback")
		return fmt.Errorf("failed to parse callback: %w", err)
	}

	// 2. 查询支付记录
	payment, err := s.paymentRepo.GetByPaymentNumber(ctx, callbackReq.OutTradeNo)
	if err != nil {
		s.logger.WithError(err).Errorf("[PaymentService] Failed to get payment")
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		return fmt.Errorf("payment not found: %s", callbackReq.OutTradeNo)
	}

	// 3. 检查支付状态（避免重复处理）
	if payment.Status == entity.PaymentStatusSuccess {
		s.logger.Infof("[PaymentService] Payment already processed: %s", callbackReq.OutTradeNo)
		return nil
	}

	// 4. 验证回调签名
	var valid bool
	switch paymentMethod {
	case entity.PaymentMethodAlipay:
		valid = s.verifyAlipayCallback(callbackData)
	case entity.PaymentMethodWechat:
		valid = s.verifyWechatCallback(callbackData)
	default:
		return fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}

	if !valid {
		s.logger.Errorf("[PaymentService] Invalid callback signature")
		return fmt.Errorf("invalid callback signature")
	}

	// 5. 更新支付状态
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 5.1 根据交易状态更新支付记录
		shouldUpdate := false

		switch callbackReq.TradeStatus {
		case "TRADE_SUCCESS", "SUCCESS":
			payment.Status = entity.PaymentStatusSuccess
			payment.TransactionID = &callbackReq.TradeNo
			now := time.Now()
			payment.ProcessedAt = &now
			shouldUpdate = true

			// 5.2 更新发票支付状态（如果有）
			if payment.InvoiceID != nil {
				if err := s.updateInvoicePayment(ctx, *payment.InvoiceID, payment.Amount); err != nil {
					s.logger.WithError(err).Errorf("[PaymentService] Failed to update invoice")
					return fmt.Errorf("failed to update invoice: %w", err)
				}
			}

			// 5.3 更新账户余额
			if err := s.updateAccountBalance(ctx, payment.TenantID, payment.Amount); err != nil {
				s.logger.WithError(err).Errorf("[PaymentService] Failed to update account balance")
				return fmt.Errorf("failed to update account balance: %w", err)
			}

		case "TRADE_CLOSED", "CLOSED":
			payment.Status = entity.PaymentStatusCancelled
			shouldUpdate = true

		case "WAIT_BUYER_PAY":
			payment.Status = entity.PaymentStatusProcessing
			shouldUpdate = true

		case "TRADE_FINISHED":
			// 支付已完成（不可退款）
			payment.Status = entity.PaymentStatusSuccess
			payment.TransactionID = &callbackReq.TradeNo
			now := time.Now()
			payment.ProcessedAt = &now
			shouldUpdate = true
		}

		if shouldUpdate {
			// 保存第三方响应
			responseJSON, _ := json.Marshal(callbackData)
			payment.ProviderResponse = string(responseJSON)

			if err := s.paymentRepo.Update(ctx, payment); err != nil {
				return fmt.Errorf("failed to update payment: %w", err)
			}
		}

		// 6. 记录日志
		s.logger.WithFields(logrus.Fields{
			"tenant_id":        payment.TenantID,
			"payment_id":       payment.ID,
			"payment_number":   payment.PaymentNumber,
			"status":           payment.Status,
			"transaction_id":   callbackReq.TradeNo,
			"callback_time_s":  time.Since(startTime).Seconds(),
		}).Info("[PaymentService] Payment callback processed")

		// 7. 更新Prometheus指标
		UpdatePaymentMetrics(payment.TenantID, paymentMethod, payment.Status, payment.Amount)

		return nil
	})
}

// QueryPaymentStatus 查询支付状态
// 参数:
//   - ctx: 上下文
//   - paymentID: 支付ID
// 返回: 支付记录或错误
func (s *PaymentService) QueryPaymentStatus(
	ctx context.Context,
	paymentID uint64,
) (*entity.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		return nil, fmt.Errorf("payment not found: %d", paymentID)
	}

	// 如果支付状态为pending或processing，主动查询第三方
	if payment.Status == entity.PaymentStatusPending || payment.Status == entity.PaymentStatusProcessing {
		s.syncPaymentStatus(ctx, payment)
	}

	return payment, nil
}

// ============================================================
// 支付宝支付
// ============================================================

// createAlipayPayment 创建支付宝支付
func (s *PaymentService) createAlipayPayment(
	ctx context.Context,
	payment *entity.Payment,
	clientIP string,
) (string, string, error) {
	if s.alipayConfig == nil {
		return "", "", fmt.Errorf("alipay config not found")
	}

	// 构建支付宝支付参数
	params := map[string]string{
		"app_id":         s.alipayConfig.AppID,
		"method":         "alipay.trade.page.pay",
		"charset":        "utf-8",
		"sign_type":      "RSA2",
		"timestamp":      time.Now().Format("2006-01-02 15:04:05"),
		"version":        "1.0",
		"out_trade_no":   payment.PaymentNumber,
		"total_amount":   fmt.Sprintf("%.2f", payment.Amount),
		"subject":        "Account Recharge",
		"product_code":   "FAST_INSTANT_TRADE_PAY",
	}

	// 添加回调URL
	if s.alipayConfig.NotifyURL != "" {
		params["notify_url"] = s.alipayConfig.NotifyURL
	}

	// 生成签名
	sign, err := s.generateAlipaySign(params)
	if err != nil {
		return "", "", err
	}
	params["sign"] = sign

	// 构建支付URL
	paymentURL := s.alipayConfig.GatewayURL + "?" + s.buildAlipayQuery(params)

	// 生成二维码数据（用于扫码支付）
	qrCodeData := paymentURL

	return paymentURL, qrCodeData, nil
}

// generateAlipaySign 生成支付宝签名
func (s *PaymentService) generateAlipaySign(params map[string]string) (string, error) {
	// 1. 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接参数
	var builder strings.Builder
	for _, k := range keys {
		if params[k] != "" && k != "sign" {
			if builder.Len() > 0 {
				builder.WriteString("&")
			}
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(params[k])
		}
	}

	// 3. RSA签名（这里简化处理，实际应使用RSA签名）
	// TODO: 实现真实的RSA2签名
	signData := builder.String()
	hash := md5.Sum([]byte(signData))
	return fmt.Sprintf("%x", hash), nil
}

// verifyAlipayCallback 验证支付宝回调签名
func (s *PaymentService) verifyAlipayCallback(params map[string]interface{}) bool {
	// TODO: 实现真实的签名验证
	// 1. 获取支付宝公钥
	// 2. 验证签名
	return true // 简化处理
}

// parseAlipayCallback 解析支付宝回调
func (s *PaymentService) parseAlipayCallback(data map[string]interface{}) (*PaymentCallbackRequest, error) {
	tradeNo, _ := data["trade_no"].(string)
	outTradeNo, _ := data["out_trade_no"].(string)
	tradeStatus, _ := data["trade_status"].(string)
	totalAmount, _ := data["total_amount"].(string)

	return &PaymentCallbackRequest{
		PaymentMethod: entity.PaymentMethodAlipay,
		TradeNo:       tradeNo,
		OutTradeNo:    outTradeNo,
		TradeStatus:   tradeStatus,
		TotalAmount:   totalAmount,
		RawData:       data,
	}, nil
}

// buildAlipayQuery 构建支付宝查询字符串
func (s *PaymentService) buildAlipayQuery(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

// ============================================================
// 微信支付
// ============================================================

// createWechatPayment 创建微信支付
func (s *PaymentService) createWechatPayment(
	ctx context.Context,
	payment *entity.Payment,
	clientIP string,
) (string, string, error) {
	if s.wechatConfig == nil {
		return "", "", fmt.Errorf("wechat config not found")
	}

	// 构建微信支付参数（Native支付）
	params := map[string]string{
		"appid":            s.wechatConfig.AppID,
		"mch_id":           s.wechatConfig.MchID,
		"nonce_str":        generateNonceStr(),
		"body":             "Account Recharge",
		"out_trade_no":     payment.PaymentNumber,
		"total_fee":        fmt.Sprintf("%.0f", payment.Amount*100), // 单位：分
		"spbill_create_ip": clientIP,
		"notify_url":       s.wechatConfig.NotifyURL,
		"trade_type":       "NATIVE",
	}

	// 生成签名
	sign, err := s.generateWechatSign(params)
	if err != nil {
		return "", "", err
	}
	params["sign"] = sign

	// 构建XML请求
	_ = s.buildWechatXML(params)

	// 发送请求到微信支付API
	// TODO: 实现真实的HTTP请求
	// resp, err := s.sendWechatRequest("https://api.mch.weixin.qq.com/pay/unifiedorder", xmlData)

	// 模拟返回二维码数据
	qrCodeData := "weixin://wxpay/bizpayurl?pr=mock_qr_code_data"

	return "", qrCodeData, nil
}

// generateWechatSign 生成微信签名
func (s *PaymentService) generateWechatSign(params map[string]string) (string, error) {
	// 1. 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 拼接参数
	var builder strings.Builder
	for _, k := range keys {
		if params[k] != "" {
			if builder.Len() > 0 {
				builder.WriteString("&")
			}
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(params[k])
		}
	}

	// 3. 添加API密钥
	builder.WriteString("&key=")
	builder.WriteString(s.wechatConfig.APIKey)

	// 4. MD5签名
	signData := builder.String()
	hash := md5.Sum([]byte(signData))
	return fmt.Sprintf("%X", hash), nil
}

// verifyWechatCallback 验证微信回调签名
func (s *PaymentService) verifyWechatCallback(params map[string]interface{}) bool {
	// TODO: 实现真实的签名验证
	return true // 简化处理
}

// parseWechatCallback 解析微信回调
func (s *PaymentService) parseWechatCallback(data map[string]interface{}) (*PaymentCallbackRequest, error) {
	tradeNo, _ := data["transaction_id"].(string)
	outTradeNo, _ := data["out_trade_no"].(string)
	tradeStatus := "SUCCESS"
	resultCode, _ := data["result_code"].(string)
	if resultCode != "SUCCESS" {
		tradeStatus = "FAILED"
	}
	totalFee, _ := data["total_fee"].(string)
	totalAmount := fmt.Sprintf("%.2f", parseFloatValue(totalFee)/100)

	return &PaymentCallbackRequest{
		PaymentMethod: entity.PaymentMethodWechat,
		TradeNo:       tradeNo,
		OutTradeNo:    outTradeNo,
		TradeStatus:   tradeStatus,
		TotalAmount:   totalAmount,
		RawData:       data,
	}, nil
}

// buildWechatXML 构建微信支付XML
func (s *PaymentService) buildWechatXML(params map[string]string) string {
	var builder strings.Builder
	builder.WriteString("<xml>")
	for k, v := range params {
		builder.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	builder.WriteString("</xml>")
	return builder.String()
}

// ============================================================
// 私有辅助方法
// ============================================================

// generatePaymentNumber 生成支付单号
func (s *PaymentService) generatePaymentNumber(
	ctx context.Context,
	tenantID string,
) (string, error) {
	// 格式: PAY-{TenantID}-{YYYYMMDDHHMMSS}-{序号}
	now := time.Now()
	prefix := fmt.Sprintf("PAY-%s-%s", tenantID, now.Format("20060102150405"))

	// 查询当前秒内已有的支付数量
	count, err := s.paymentRepo.CountByTenantAndTimestamp(ctx, tenantID, now)
	if err != nil {
		return "", fmt.Errorf("failed to count payments: %w", err)
	}

	sequence := count + 1
	paymentNumber := fmt.Sprintf("%s-%03d", prefix, sequence)

	return paymentNumber, nil
}

// getProvider 根据支付方式获取提供商
func (s *PaymentService) getProvider(paymentMethod string) string {
	switch paymentMethod {
	case entity.PaymentMethodAlipay:
		return entity.PaymentProviderAlipay
	case entity.PaymentMethodWechat:
		return entity.PaymentProviderWechat
	default:
		return "unknown"
	}
}

// updateInvoicePayment 更新发票支付状态
func (s *PaymentService) updateInvoicePayment(
	ctx context.Context,
	invoiceID uint64,
	amount float64,
) error {
	invoice, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if invoice == nil {
		return fmt.Errorf("invoice not found: %d", invoiceID)
	}

	// 更新支付金额
	invoice.PaidAmount += amount

	// 更新状态
	if invoice.PaidAmount >= invoice.TotalAmount {
		invoice.Status = entity.InvoiceStatusPaid
		now := time.Now()
		invoice.PaidAt = &now
	} else if invoice.PaidAmount > 0 {
		invoice.Status = entity.InvoiceStatusPartialPaid
	}

	return s.invoiceRepo.Update(ctx, invoice)
}

// updateAccountBalance 更新账户余额
func (s *PaymentService) updateAccountBalance(
	ctx context.Context,
	tenantID string,
	amount float64,
) error {
	account, err := s.billingAccountRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("billing account not found: %s", tenantID)
	}

	// 增加可用额度
	account.AvailableCredit += amount

	return s.billingAccountRepo.Update(ctx, account)
}

// syncPaymentStatus 同步支付状态
func (s *PaymentService) syncPaymentStatus(ctx context.Context, payment *entity.Payment) error {
	// TODO: 实现主动查询第三方支付状态
	// 1. 根据支付方式调用查询接口
	// 2. 更新本地支付状态
	s.logger.Infof("[PaymentService] Syncing payment status for: %s", payment.PaymentNumber)
	return nil
}

// generateNonceStr 生成随机字符串
func generateNonceStr() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// parseFloat 解析浮点数
func parseFloatValue(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// sendWechatRequest 发送微信支付请求
func (s *PaymentService) sendWechatRequest(url, xmlData string) ([]byte, error) {
	// TODO: 实现真实的HTTP请求
	resp, err := http.Post(url, "application/xml", strings.NewReader(xmlData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
