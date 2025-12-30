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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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
}

// CreateWebhookRequest 创建Webhook请求
type CreateWebhookRequest struct {
	TenantID    string                    `json:"tenant_id" validate:"required"`
	ProjectID   string                    `json:"project_id" validate:"required"`
	WebhookURL  string                    `json:"webhook_url" validate:"required,url,max=500"`
	Events      entity.WebhookEvents      `json:"events" validate:"required"`
}

// UpdateWebhookRequest 更新Webhook请求
type UpdateWebhookRequest struct {
	WebhookID   string                `json:"webhook_id" validate:"required"`
	WebhookURL  *string               `json:"webhook_url" validate:"omitempty,url,max=500"`
	Events      *entity.WebhookEvents `json:"events"`
	Status      *entity.WebhookStatus `json:"status"`
}

// WebhookListFilter Webhook列表过滤器
type WebhookListFilter struct {
	TenantID  string                 `validate:"required"`
	ProjectID string
	Status    entity.WebhookStatus
	PageToken string
	PageSize  int
}

// WebhookDetailedStats Webhook详细统计信息
type WebhookDetailedStats struct {
	WebhookID      string    `json:"webhook_id"`
	TotalCalls     int64     `json:"total_calls"`
	SuccessCalls   int64     `json:"success_calls"`
	FailureCalls   int64     `json:"failure_calls"`
	SuccessRate    float64   `json:"success_rate"`
	AvgDuration    float64   `json:"avg_duration"`
	LastTriggeredAt *int64   `json:"last_triggered_at,omitempty"`
	LastStatusCode int       `json:"last_status_code"`
}

// webhookService Webhook服务实现
type webhookService struct {
	webhookRepo repository.WebhookRepository
	logRepo     repository.WebhookLogRepository
	httpClient  *http.Client
}

// NewWebhookService 创建Webhook服务实例
func NewWebhookService(
	webhookRepo repository.WebhookRepository,
	logRepo repository.WebhookLogRepository,
) WebhookService {
	return &webhookService{
		webhookRepo: webhookRepo,
		logRepo:     logRepo,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateWebhook 创建Webhook
func (s *webhookService) CreateWebhook(ctx context.Context, req *CreateWebhookRequest) (*entity.Webhook, string, error) {
	// TODO: 验证项目是否存在且属于该租户

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

	// 异步触发所有匹配的Webhook
	for _, webhook := range webhooks {
		if webhook.ShouldTrigger(eventType) {
			go s.triggerWebhookAsync(context.Background(), webhook, eventType, payload)
		}
	}

	return nil
}

// triggerWebhookAsync 异步触发Webhook
func (s *webhookService) triggerWebhookAsync(ctx context.Context, webhook *entity.Webhook, eventType entity.WebhookEventType, payload *entity.WebhookPayload) {
	startTime := time.Now()
	eventID := generateID("event")

	// 更新payload
	payload.EventID = eventID
	payload.EventType = eventType
	payload.Timestamp = time.Now().UnixMilli()

	// 序列化payload
	body, err := json.Marshal(payload)
	if err != nil {
		s.logWebhookFailure(ctx, webhook.WebhookID, eventType, err, 0)
		return
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", webhook.WebhookURL, nil)
	if err != nil {
		s.logWebhookFailure(ctx, webhook.WebhookID, eventType, err, 0)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Coze-Webhook/1.0")

	// 生成签名
	signature := s.generateSignature(body, webhook.WebhookSecret)
	req.Header.Set("X-Coze-Signature", signature)
	req.Header.Set("X-Coze-Event-ID", eventID)
	req.Header.Set("X-Coze-Event-Type", string(eventType))
	req.Header.Set("X-Coze-Timestamp", fmt.Sprintf("%d", payload.Timestamp))

	// 发送请求
	resp, err := s.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		s.logWebhookFailure(ctx, webhook.WebhookID, eventType, err, duration.Milliseconds())
		s.webhookRepo.UpdateStats(context.Background(), webhook.WebhookID, false)
		return
	}
	defer resp.Body.Close()

	// 记录日志
	log := &entity.WebhookLog{
		LogID:      generateID("log"),
		WebhookID:  webhook.WebhookID,
		EventType:  string(eventType),
		StatusCode: resp.StatusCode,
		DurationMs: duration.Milliseconds(),
		Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
		RetryCount: 0,
	}

	_ = s.logRepo.Create(context.Background(), log)

	// 更新统计
	s.webhookRepo.UpdateStats(context.Background(), webhook.WebhookID, log.Success)
	s.webhookRepo.UpdateLastTriggerAt(context.Background(), webhook.WebhookID)
}

// logWebhookFailure 记录Webhook失败日志
func (s *webhookService) logWebhookFailure(ctx context.Context, webhookID string, eventType entity.WebhookEventType, err error, duration int64) {
	log := &entity.WebhookLog{
		LogID:      generateID("log"),
		WebhookID:  webhookID,
		EventType:  string(eventType),
		StatusCode: 0,
		Response:   err.Error(),
		DurationMs: duration,
		Success:    false,
		RetryCount: 0,
	}

	_ = s.logRepo.Create(ctx, log)
}

// VerifyWebhookSignature 验证Webhook签名
func (s *webhookService) VerifyWebhookSignature(ctx context.Context, payload []byte, signature string, secret string) bool {
	expectedSignature := s.generateSignature(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// generateSignature 生成签名
func (s *webhookService) generateSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
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
	// TODO: 使用加密安全的随机数生成器
	return "wh_secret_" + randomString(32)
}
