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
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/domain/developer/repository"
)

// WebhookService Webhook服务接口
type WebhookService interface {
	// CreateWebhook 创建Webhook
	CreateWebhook(ctx context.Context, req *CreateWebhookRequest) (*entity.Webhook, string, error)

	// GetWebhook 获取Webhook详情
	GetWebhook(ctx context.Context, webhookID string) (*entity.Webhook, error)

	// ListWebhooks 查询Webhook列表
	ListWebhooks(ctx context.Context, filter *WebhookListFilter) ([]*entity.Webhook, int64, error)

	// UpdateWebhook 更新Webhook
	UpdateWebhook(ctx context.Context, req *UpdateWebhookRequest) error

	// DeleteWebhook 删除Webhook
	DeleteWebhook(ctx context.Context, webhookID string) error

	// TriggerWebhook 触发Webhook
	TriggerWebhook(ctx context.Context, eventType entity.WebhookEventType, payload *entity.WebhookPayload) error

	// VerifyWebhookSignature 验证Webhook签名
	VerifyWebhookSignature(ctx context.Context, payload []byte, signature string, secret string) bool

	// PauseWebhook 暂停Webhook
	PauseWebhook(ctx context.Context, webhookID string) error

	// ResumeWebhook 恢复Webhook
	ResumeWebhook(ctx context.Context, webhookID string) error

	// GetProjectWebhooks 获取项目的Webhook列表
	GetProjectWebhooks(ctx context.Context, projectID string) ([]*entity.Webhook, error)

	// GetWebhookStats 获取Webhook统计信息
	GetWebhookStats(ctx context.Context, webhookID string) (*WebhookDetailedStats, error)

	// RetryDeadLetterQueue 重试死信队列条目
	RetryDeadLetterQueue(ctx context.Context) error
}

// CreateWebhookRequest 创建Webhook请求
type CreateWebhookRequest struct {
	TenantID   string                    `json:"tenant_id" validate:"required"`
	ProjectID  string                    `json:"project_id" validate:"required"`
	WebhookURL string                    `json:"webhook_url" validate:"required,url,max=500"`
	Events     entity.WebhookEvents      `json:"events" validate:"required"`
}

// UpdateWebhookRequest 更新Webhook请求
type UpdateWebhookRequest struct {
	WebhookID  string                `json:"webhook_id" validate:"required"`
	WebhookURL *string               `json:"webhook_url" validate:"omitempty,url,max=500"`
	Events     *entity.WebhookEvents `json:"events"`
	Status     *entity.WebhookStatus `json:"status"`
}

// WebhookListFilter Webhook列表过滤器
type WebhookListFilter struct {
	TenantID  string
	ProjectID string
	Status    entity.WebhookStatus
	PageToken string
	PageSize  int
}

// WebhookDetailedStats Webhook详细统计信息
type WebhookDetailedStats struct {
	WebhookID       string  `json:"webhook_id"`
	TotalCalls      int64   `json:"total_calls"`
	SuccessCalls    int64   `json:"success_calls"`
	FailureCalls    int64   `json:"failure_calls"`
	SuccessRate     float64 `json:"success_rate"`
	AvgDuration     float64 `json:"avg_duration"`
	LastTriggeredAt *int64  `json:"last_triggered_at,omitempty"`
	LastStatusCode  int     `json:"last_status_code"`
}

// webhookService Webhook服务实现
type webhookService struct {
	webhookRepo repository.WebhookRepository
	logRepo     repository.WebhookLogRepository
	dlqRepo     repository.WebhookDLQRepository
	httpClient  *http.Client
	config      *entity.WebhookRetryConfig
}

// NewWebhookService 创建Webhook服务实例
func NewWebhookService(
	webhookRepo repository.WebhookRepository,
	logRepo repository.WebhookLogRepository,
	dlqRepo repository.WebhookDLQRepository,
) WebhookService {
	return &webhookService{
		webhookRepo: webhookRepo,
		logRepo:     logRepo,
		dlqRepo:     dlqRepo,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		config:      entity.DefaultWebhookRetryConfig(),
	}
}

// SetRetryConfig 设置重试配置
func (s *webhookService) SetRetryConfig(config *entity.WebhookRetryConfig) {
	s.config = config
	s.httpClient.Timeout = config.RequestTimeout
}

// CreateWebhook 创建Webhook
func (s *webhookService) CreateWebhook(ctx context.Context, req *CreateWebhookRequest) (*entity.Webhook, string, error) {
	// 生成Webhook密钥（用于签名验证）
	secret := generateWebhookSecret()

	webhook := &entity.Webhook{
		WebhookID:    generateID("webhook"),
		TenantID:     req.TenantID,
		ProjectID:    req.ProjectID,
		WebhookURL:   req.WebhookURL,
		WebhookSecret: secret,
		Events:       req.Events,
		Status:       entity.WebhookStatusActive,
	}

	if err := s.webhookRepo.Create(ctx, webhook); err != nil {
		return nil, "", err
	}

	return webhook, secret, nil
}

// GetWebhook 获取Webhook详情
func (s *webhookService) GetWebhook(ctx context.Context, webhookID string) (*entity.Webhook, error) {
	return s.webhookRepo.GetByID(ctx, webhookID)
}

// ListWebhooks 查询Webhook列表
func (s *webhookService) ListWebhooks(ctx context.Context, filter *WebhookListFilter) ([]*entity.Webhook, int64, error) {
	repoFilter := &repository.WebhookFilter{
		TenantID:  filter.TenantID,
		ProjectID: filter.ProjectID,
		Status:    filter.Status,
		PageToken: filter.PageToken,
		PageSize:  filter.PageSize,
	}

	return s.webhookRepo.List(ctx, repoFilter)
}

// UpdateWebhook 更新Webhook
func (s *webhookService) UpdateWebhook(ctx context.Context, req *UpdateWebhookRequest) error {
	webhook, err := s.webhookRepo.GetByID(ctx, req.WebhookID)
	if err != nil {
		return err
	}
	if webhook == nil {
		return errors.New("webhook not found")
	}

	// 更新字段
	if req.WebhookURL != nil {
		webhook.WebhookURL = *req.WebhookURL
	}
	if req.Events != nil {
		webhook.Events = *req.Events
	}
	if req.Status != nil {
		webhook.Status = *req.Status
	}

	return s.webhookRepo.Update(ctx, webhook)
}

// DeleteWebhook 删除Webhook
func (s *webhookService) DeleteWebhook(ctx context.Context, webhookID string) error {
	return s.webhookRepo.Delete(ctx, webhookID)
}

// TriggerWebhook 触发Webhook
func (s *webhookService) TriggerWebhook(ctx context.Context, eventType entity.WebhookEventType, payload *entity.WebhookPayload) error {
	// 查找所有需要触发该事件的Webhook
	webhooks, err := s.webhookRepo.GetByProjectID(ctx, payload.ProjectID)
	if err != nil {
		return err
	}

	// 异步触发所有匹配的Webhook（带重试机制）
	for _, webhook := range webhooks {
		if webhook.ShouldTrigger(eventType) {
			go s.triggerWebhookWithRetry(context.Background(), webhook, eventType, payload)
		}
	}

	return nil
}

// triggerWebhookWithRetry 带重试机制的Webhook触发
func (s *webhookService) triggerWebhookWithRetry(ctx context.Context, webhook *entity.Webhook, eventType entity.WebhookEventType, payload *entity.WebhookPayload) {
	var lastErr error

	// 指数退避重试
	for attempt := 0; attempt <= s.config.MaxRetries; attempt++ {
		// 克隆payload避免修改原始数据
		payloadCopy := *payload
		payloadCopy.EventID = generateID("event")
		payloadCopy.EventType = eventType
		payloadCopy.Timestamp = time.Now().UnixMilli()

		// 执行HTTP请求
		err := s.doTrigger(ctx, webhook, &payloadCopy, attempt)
		if err == nil {
			// 成功，更新统计
			s.updateSuccessStats(ctx, webhook.WebhookID)
			return
		}

		lastErr = err

		// 如果是最后一次尝试，不再重试
		if attempt == s.config.MaxRetries {
			break
		}

		// 计算退避时间
		delay := s.calculateBackoff(attempt)

		// 记录重试日志
		s.logWebhookRetry(ctx, webhook.WebhookID, eventType, attempt+1, delay, err)

		// 等待后重试
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			s.updateFailureStats(ctx, webhook.WebhookID)
			return
		}
	}

	// 所有重试都失败，移动到死信队列
	s.updateFailureStats(ctx, webhook.WebhookID)
	if err := s.moveToDLQ(ctx, webhook, payload, lastErr); err != nil {
		// 记录日志但不中断流程
		fmt.Printf("⚠️  Failed to move to DLQ: %v\n", err)
	}
}

// doTrigger 实际执行HTTP请求
func (s *webhookService) doTrigger(ctx context.Context, webhook *entity.Webhook, payload *entity.WebhookPayload, attempt int) error {
	startTime := time.Now()

	// 序列化payload
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", webhook.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ZKER-Webhook/1.0")

	// 计算签名
	signature := s.calculateSignature(body, webhook.WebhookSecret)
	req.Header.Set("X-ZKER-Signature", signature)
	req.Header.Set("X-ZKER-Event-ID", payload.EventID)
	req.Header.Set("X-ZKER-Event-Type", string(payload.EventType))
	req.Header.Set("X-ZKER-Timestamp", fmt.Sprintf("%d", payload.Timestamp))

	if attempt > 0 {
		req.Header.Set("X-ZKER-Retry-Count", fmt.Sprintf("%d", attempt))
	}

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(startTime)

	// 记录日志
	s.logWebhookAttempt(ctx, webhook.WebhookID, payload.EventType, resp.StatusCode, duration, attempt)

	// 检查响应状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned error status: %d, body: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

// calculateBackoff 计算指数退避时间
func (s *webhookService) calculateBackoff(attempt int) time.Duration {
	// 指数退避: 2^attempt
	delay := s.config.BaseDelay * time.Duration(math.Pow(2, float64(attempt)))

	// 限制最大延迟
	if delay > s.config.MaxDelay {
		delay = s.config.MaxDelay
	}

	// 添加随机抖动（±20%）避免惊群效应
	jitter := time.Duration(mathrand.IntN(int(delay) / 5))
	delay += jitter - (jitter / 2)

	return delay
}

// calculateSignature 计算HMAC-SHA256签名
func (s *webhookService) calculateSignature(body []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

// moveToDLQ 移动到死信队列
func (s *webhookService) moveToDLQ(ctx context.Context, webhook *entity.Webhook, payload *entity.WebhookPayload, lastErr error) error {
	dlqEntry := &entity.WebhookDLQEntry{
		DLQID:        generateID("dlq"),
		WebhookID:    webhook.WebhookID,
		EventType:    payload.EventType,
		Payload:      payload,
		ErrorMessage: lastErr.Error(),
		RetryCount:   0,
	}

	return s.dlqRepo.Create(ctx, dlqEntry)
}

// logWebhookAttempt 记录Webhook尝试日志
func (s *webhookService) logWebhookAttempt(ctx context.Context, webhookID string, eventType entity.WebhookEventType, statusCode int, duration time.Duration, retryCount int) {
	log := &entity.WebhookLog{
		LogID:      generateID("log"),
		WebhookID:  webhookID,
		EventType:  string(eventType),
		StatusCode: statusCode,
		DurationMs: duration.Milliseconds(),
		Success:    statusCode >= 200 && statusCode < 300,
		RetryCount: retryCount,
	}

	_ = s.logRepo.Create(ctx, log)
}

// logWebhookRetry 记录Webhook重试日志
func (s *webhookService) logWebhookRetry(ctx context.Context, webhookID string, eventType entity.WebhookEventType, attempt int, delay time.Duration, err error) {
	log := &entity.WebhookLog{
		LogID:      generateID("log"),
		WebhookID:  webhookID,
		EventType:  string(eventType),
		StatusCode: 0,
		Response:   fmt.Sprintf("Retry %d after %v: %v", attempt, delay, err.Error()),
		DurationMs: 0,
		Success:    false,
		RetryCount: attempt,
	}

	_ = s.logRepo.Create(ctx, log)
}

// updateSuccessStats 更新成功统计
func (s *webhookService) updateSuccessStats(ctx context.Context, webhookID string) {
	_ = s.webhookRepo.UpdateStats(ctx, webhookID, true)
	_ = s.webhookRepo.UpdateLastTriggerAt(ctx, webhookID)
}

// updateFailureStats 更新失败统计
func (s *webhookService) updateFailureStats(ctx context.Context, webhookID string) {
	_ = s.webhookRepo.UpdateStats(ctx, webhookID, false)
}

// RetryDeadLetterQueue 重试死信队列条目
func (s *webhookService) RetryDeadLetterQueue(ctx context.Context) error {
	// 获取可重试的条目
	entries, err := s.dlqRepo.GetRetryable(ctx, s.config.MaxRetries)
	if err != nil {
		return fmt.Errorf("failed to get retryable entries: %w", err)
	}

	// 逐个重试
	for _, entry := range entries {
		webhook, err := s.webhookRepo.GetByID(ctx, entry.WebhookID)
		if err != nil || webhook == nil {
			// Webhook不存在，跳过
			continue
		}

		// 执行重试
		err = s.doTrigger(ctx, webhook, entry.Payload, entry.RetryCount)
		if err == nil {
			// 成功，从死信队列删除
			_ = s.dlqRepo.Delete(ctx, entry.DLQID)
			s.updateSuccessStats(ctx, webhook.WebhookID)
		} else {
			// 失败，更新重试次数
			entry.RetryCount++
			_ = s.dlqRepo.UpdateRetryCount(ctx, entry.DLQID, entry.RetryCount)
		}
	}

	return nil
}

// VerifyWebhookSignature 验证Webhook签名
func (s *webhookService) VerifyWebhookSignature(ctx context.Context, payload []byte, signature string, secret string) bool {
	expectedSignature := s.calculateSignature(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// PauseWebhook 暂停Webhook
func (s *webhookService) PauseWebhook(ctx context.Context, webhookID string) error {
	return s.webhookRepo.UpdateStatus(ctx, webhookID, entity.WebhookStatusPaused)
}

// ResumeWebhook 恢复Webhook
func (s *webhookService) ResumeWebhook(ctx context.Context, webhookID string) error {
	return s.webhookRepo.UpdateStatus(ctx, webhookID, entity.WebhookStatusActive)
}

// GetProjectWebhooks 获取项目的Webhook列表
func (s *webhookService) GetProjectWebhooks(ctx context.Context, projectID string) ([]*entity.Webhook, error) {
	return s.webhookRepo.GetByProjectID(ctx, projectID)
}

// GetWebhookStats 获取Webhook统计信息
func (s *webhookService) GetWebhookStats(ctx context.Context, webhookID string) (*WebhookDetailedStats, error) {
	// 获取最近7天的统计
	startTime := time.Now().AddDate(0, 0, -7).UnixMilli()
	endTime := time.Now().UnixMilli()

	stats, err := s.logRepo.GetStats(ctx, webhookID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
	if err != nil {
		return nil, err
	}

	detailedStats := &WebhookDetailedStats{
		WebhookID:       webhookID,
		TotalCalls:      stats.TotalCalls,
		SuccessCalls:    stats.SuccessCalls,
		FailureCalls:    stats.FailureCalls,
		SuccessRate:     stats.SuccessRate,
		AvgDuration:     stats.AvgDuration,
		LastTriggeredAt: webhook.LastTriggerAt,
	}

	return detailedStats, nil
}

// generateWebhookSecret 生成Webhook密钥
func generateWebhookSecret() string {
	// 生成32字节随机密钥
	key := make([]byte, 32)
	if _, err := cryptorand.Read(key); err != nil {
		// 备用方案：使用时间戳
		return fmt.Sprintf("whsec_%d_%s", time.Now().UnixNano(), randomString(16))
	}
	return "whsec_" + base64.StdEncoding.EncodeToString(key)
}
