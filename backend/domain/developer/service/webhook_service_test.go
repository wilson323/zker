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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
)

// TestCalculateBackoff 测试指数退避计算
func TestCalculateBackoff(t *testing.T) {
	service := &webhookService{
		config: &entity.WebhookRetryConfig{
			BaseDelay: time.Second,
			MaxDelay:  300 * time.Second,
		},
	}

	testCases := []struct {
		name            string
		attempt         int
		expectedMin     time.Duration
		expectedMax     time.Duration
	}{
		{
			name:        "Attempt 0",
			attempt:     0,
			expectedMin: 800 * time.Millisecond,  // 1s - 20%
			expectedMax: 1200 * time.Millisecond, // 1s + 20%
		},
		{
			name:        "Attempt 1",
			attempt:     1,
			expectedMin: 1600 * time.Millisecond, // 2s - 20%
			expectedMax: 2400 * time.Millisecond, // 2s + 20%
		},
		{
			name:        "Attempt 2",
			attempt:     2,
			expectedMin: 3200 * time.Millisecond, // 4s - 20%
			expectedMax: 4800 * time.Millisecond, // 4s + 20%
		},
		{
			name:        "Attempt 8",
			attempt:     8,
			expectedMin: 240 * time.Second, // 300s - 20%
			expectedMax: 300 * time.Second, // capped at MaxDelay
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 运行多次确保抖动在合理范围内
			for i := 0; i < 100; i++ {
				delay := service.calculateBackoff(tc.attempt)

				if delay < tc.expectedMin || delay > tc.expectedMax {
					t.Errorf("Delay %v outside expected range [%v, %v]", delay, tc.expectedMin, tc.expectedMax)
				}
			}
		})
	}
}

// TestCalculateSignature 测试签名计算
func TestCalculateSignature(t *testing.T) {
	service := &webhookService{}

	payload := []byte(`{"test": "data"}`)
	secret := "test_secret"

	signature := service.calculateSignature(payload, secret)

	// 签名应该是固定格式
	if !strings.HasPrefix(signature, "sha256=") {
		t.Errorf("Signature should start with 'sha256=', got: %s", signature)
	}

	// 相同输入应该产生相同签名
	signature2 := service.calculateSignature(payload, secret)
	if signature != signature2 {
		t.Error("Same input should produce same signature")
	}

	// 不同输入应该产生不同签名
	differentPayload := []byte(`{"test": "different"}`)
	signature3 := service.calculateSignature(differentPayload, secret)
	if signature == signature3 {
		t.Error("Different inputs should produce different signatures")
	}

	// 不同密钥应该产生不同签名
	signature4 := service.calculateSignature(payload, "different_secret")
	if signature == signature4 {
		t.Error("Different secrets should produce different signatures")
	}
}

// TestVerifyWebhookSignature 测试签名验证
func TestVerifyWebhookSignature(t *testing.T) {
	service := &webhookService{}

	payload := []byte(`{"test": "data"}`)
	secret := "test_secret"

	// 正确的签名
	correctSignature := service.calculateSignature(payload, secret)
	if !service.VerifyWebhookSignature(context.Background(), payload, correctSignature, secret) {
		t.Error("Correct signature should verify successfully")
	}

	// 错误的签名
	wrongSignature := "sha256=wrong"
	if service.VerifyWebhookSignature(context.Background(), payload, wrongSignature, secret) {
		t.Error("Wrong signature should fail verification")
	}

	// 错误的密钥
	wrongSecret := "wrong_secret"
	if service.VerifyWebhookSignature(context.Background(), payload, correctSignature, wrongSecret) {
		t.Error("Signature with wrong secret should fail verification")
	}
}

// TestDoTrigger 测试HTTP请求触发
func TestDoTrigger(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Content-Type should be application/json")
		}

		if r.Header.Get("User-Agent") != "ZKER-Webhook/1.0" {
			t.Error("User-Agent should be ZKER-Webhook/1.0")
		}

		signature := r.Header.Get("X-ZKER-Signature")
		if signature == "" || !strings.HasPrefix(signature, "sha256=") {
			t.Error("X-ZKER-Signature header should be present and valid")
		}

		// 读取body
		var payload entity.WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode payload: %v", err)
		}

		// 返回成功
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	service := &webhookService{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	webhook := &entity.Webhook{
		WebhookID:    "test_webhook",
		WebhookURL:   server.URL,
		WebhookSecret: "test_secret",
	}

	payload := &entity.WebhookPayload{
		EventID:    "test_event",
		EventType:  entity.WebhookEventBotCreated,
		Timestamp:  time.Now().UnixMilli(),
		TenantID:   "tenant_1",
		ProjectID:  "project_1",
		Data:       map[string]interface{}{"test": "data"},
	}

	// 执行触发
	err := service.doTrigger(context.Background(), webhook, payload, 0)
	if err != nil {
		t.Errorf("doTrigger failed: %v", err)
	}
}

// TestDoTriggerHTTPError 测试HTTP错误处理
func TestDoTriggerHTTPError(t *testing.T) {
	// 创建返回错误的服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	service := &webhookService{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	webhook := &entity.Webhook{
		WebhookID:    "test_webhook",
		WebhookURL:   server.URL,
		WebhookSecret: "test_secret",
	}

	payload := &entity.WebhookPayload{
		EventID:    "test_event",
		EventType:  entity.WebhookEventBotCreated,
		Timestamp:  time.Now().UnixMilli(),
		TenantID:   "tenant_1",
		ProjectID:  "project_1",
	}

	// 执行触发（应该失败）
	err := service.doTrigger(context.Background(), webhook, payload, 0)
	if err == nil {
		t.Error("doTrigger should return error for 500 status")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf("Error should mention status code 500, got: %v", err)
	}
}

// TestGenerateWebhookSecret 测试Webhook密钥生成
func TestGenerateWebhookSecret(t *testing.T) {
	// 生成多个密钥
	secret1 := generateWebhookSecret()
	secret2 := generateWebhookSecret()

	// 密钥应该不同
	if secret1 == secret2 {
		t.Error("Generated secrets should be different")
	}

	// 密钥应该有前缀
	if !strings.HasPrefix(secret1, "whsec_") {
		t.Errorf("Secret should have 'whsec_' prefix, got: %s", secret1)
	}

	// 密钥长度应该合理（whsec_ + base64编码的32字节 = 6 + 44 = 50字符）
	if len(secret1) < 40 {
		t.Errorf("Secret length should be at least 40, got: %d", len(secret1))
	}
}

// TestWebhookDLQEntry 测试死信队列条目
func TestWebhookDLQEntry(t *testing.T) {
	entry := &entity.WebhookDLQEntry{
		DLQID:        "dlq_1",
		WebhookID:    "webhook_1",
		EventType:    entity.WebhookEventBotCreated,
		Payload:      &entity.WebhookPayload{},
		ErrorMessage: "Test error",
		RetryCount:   0,
		CreatedAt:    time.Now().UnixMilli(),
	}

	// 测试CanRetry
	if !entry.CanRetry(5) {
		t.Error("Entry with retry_count=0 should be retryable")
	}

	entry.RetryCount = 5
	if entry.CanRetry(5) {
		t.Error("Entry with retry_count=5 should not be retryable with max_retries=5")
	}

	// 测试IsExpired
	oldEntry := &entity.WebhookDLQEntry{
		CreatedAt: time.Now().Add(-10 * 24 * time.Hour).UnixMilli(),
	}
	if !oldEntry.IsExpired() {
		t.Error("Entry older than 7 days should be expired")
	}

	newEntry := &entity.WebhookDLQEntry{
		CreatedAt: time.Now().UnixMilli(),
	}
	if newEntry.IsExpired() {
		t.Error("New entry should not be expired")
	}
}

// TestWebhookRetryConfig 测试重试配置
func TestWebhookRetryConfig(t *testing.T) {
	config := entity.DefaultWebhookRetryConfig()

	if config.MaxRetries != 5 {
		t.Errorf("Default MaxRetries should be 5, got: %d", config.MaxRetries)
	}

	if config.BaseDelay != time.Second {
		t.Errorf("Default BaseDelay should be 1s, got: %v", config.BaseDelay)
	}

	if config.MaxDelay != 5*time.Minute {
		t.Errorf("Default MaxDelay should be 5m, got: %v", config.MaxDelay)
	}

	if config.RequestTimeout != 30*time.Second {
		t.Errorf("Default RequestTimeout should be 30s, got: %v", config.RequestTimeout)
	}

	if config.DLQRetentionDays != 7 {
		t.Errorf("Default DLQRetentionDays should be 7, got: %d", config.DLQRetentionDays)
	}
}
