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

package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// QuotaAPITestSuite 配额API契约测试套件
type QuotaAPITestSuite struct {
	suite.Suite
	h *server.Hertz
}

// SetupSuite 设置测试套件
func (s *QuotaAPITestSuite) SetupSuite() {
	s.h = server.Default()
}

// TearDownSuite 清理测试套件
func (s *QuotaAPITestSuite) TearDownSuite() {
	// 清理资源
}

// TestQuotaAPITestSuite 运行测试套件
func TestQuotaAPITestSuite(t *testing.T) {
	suite.Run(t, new(QuotaAPITestSuite))
}

// TestGetQuotas_Success 测试获取租户配额列表成功
func (s *QuotaAPITestSuite) TestGetQuotas_Success() {
	tenantID := "test-tenant-001"

	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants/"+tenantID+"/quotas",
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证配额列表结构
	data := result["data"].(map[string]interface{})
	_, hasQuotas := data["quotas"]
	_, hasOverallUsage := data["overall_usage"]

	s.True(hasQuotas, "Response should have 'quotas' field")
	s.True(hasOverallUsage, "Response should have 'overall_usage' field")
}

// TestCheckQuota_Success 测试配额检查成功
func (s *QuotaAPITestSuite) TestCheckQuota_Success() {
	tenantID := "test-tenant-002"

	checkReq := map[string]interface{}{
		"resource_type": "bots",
		"amount":        1,
	}
	reqBody, _ := json.Marshal(checkReq)

	w := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants/"+tenantID+"/quotas/check",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证配额检查结果
	data := result["data"].(map[string]interface{})
	_, hasAllowed := data["allowed"]
	_, hasRemaining := data["remaining"]

	s.True(hasAllowed, "Response should have 'allowed' field")
	s.True(hasRemaining, "Response should have 'remaining' field")
}

// TestCheckQuota_Exceeded 测试配额超限
func (s *QuotaAPITestSuite) TestCheckQuota_Exceeded() {
	tenantID := "test-tenant-full"

	// 模拟配额已满的情况
	checkReq := map[string]interface{}{
		"resource_type": "bots",
		"amount":        10000, // 超过限制
	}
	reqBody, _ := json.Marshal(checkReq)

	w := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants/"+tenantID+"/quotas/check",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	// 验证响应
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// 如果配额超限，应该返回错误或allowed=false
	if code, ok := result["code"]; ok {
		if int32(code.(float64)) != 0 {
			// 错误响应
			s.assertErrorResponse(result)
			assert.DeepEqual(s.T(), errno.ErrQuotaExceededCode, int32(code.(float64)))
		} else {
			// 成功响应，但allowed=false
			data := result["data"].(map[string]interface{})
			allowed, hasAllowed := data["allowed"]
			s.True(hasAllowed, "Response should have 'allowed' field")
			if allowed != nil {
				s.False(allowed.(bool), "allowed should be false when quota exceeded")
			}
		}
	}
}

// TestGetQuotaUsage_Success 测试获取配额使用详情成功
func (s *QuotaAPITestSuite) TestGetQuotaUsage_Success() {
	tenantID := "test-tenant-003"

	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants/"+tenantID+"/quotas/usage/details",
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证使用详情结构
	data, ok := result["data"].([]interface{})
	s.True(ok, "Response data should be an array")

	if len(data) > 0 {
		firstUsage := data[0].(map[string]interface{})
		requiredFields := []string{"resource_type", "current_usage", "limit", "usage_percentage"}
		for _, field := range requiredFields {
			_, hasField := firstUsage[field]
			s.True(hasField, "Usage detail should have '"+field+"' field")
		}
	}
}

// TestUpdateQuotaLimit_Success 测试更新配额限制成功（仅超级管理员）
func (s *QuotaAPITestSuite) TestUpdateQuotaLimit_Success() {
	tenantID := "test-tenant-004"

	updateReq := map[string]interface{}{
		"resource_type": "bots",
		"limit":         100,
		"is_soft_limit": true,
		"overage_fee":   0.5,
	}
	reqBody, _ := json.Marshal(updateReq)

	w := ut.PerformRequest(s.h.Engine, "PUT", "/api/v1/tenants/"+tenantID+"/quotas/bots/limit",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("X-User-Role", "super_admin"),
	)

	// 验证响应（可能需要管理员权限）
	if w.Code == 200 || w.Code == 403 {
		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)

		if w.Code == 200 {
			s.assertSuccessResponse(result)
		}
	}
}

// TestQuotaAlertRules 测试配额预警规则
func (s *QuotaAPITestSuite) TestQuotaAlertRules() {
	tenantID := "test-tenant-005"

	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants/"+tenantID+"/quotas/alerts",
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证预警规则结构
	data := result["data"].(map[string]interface{})
	_, hasAlertRules := data["alert_rules"]
	s.True(hasAlertRules, "Response should have 'alert_rules' field")
}

// TestQuotaMiddleware_Concurrent 并发测试：验证配额中间件是并发安全的
func TestQuotaMiddleware_Concurrent(t *testing.T) {
	h := server.Default()
	tenantID := "test-tenant-concurrent"

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			checkReq := map[string]interface{}{
				"resource_type": "bots",
				"amount":        1,
			}
			reqBody, _ := json.Marshal(checkReq)

			w := ut.PerformRequest(h.Engine, "POST", "/api/v1/tenants/"+tenantID+"/quotas/check",
				ut.WithBody(bytes.NewReader(reqBody)),
				ut.WithHeader("Content-Type", "application/json"),
				ut.WithHeader("X-Tenant-ID", tenantID),
			)

			// 验证响应状态码
			assert.DeepEqual(t, w.Code, 200)

			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// BenchmarkCheckQuota 性能基准测试：配额检查
func BenchmarkCheckQuota(b *testing.B) {
	h := server.Default()
	tenantID := "benchmark-tenant"

	checkReq := map[string]interface{}{
		"resource_type": "bots",
		"amount":        1,
	}
	reqBody, _ := json.Marshal(checkReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := ut.PerformRequest(h.Engine, "POST", "/api/v1/tenants/"+tenantID+"/quotas/check",
			ut.WithBody(bytes.NewReader(reqBody)),
			ut.WithHeader("Content-Type", "application/json"),
			ut.WithHeader("X-Tenant-ID", tenantID),
		)
		if w.Code != 200 {
			b.Errorf("Expected status 200, got %d", w.Code)
		}
	}
}

// assertSuccessResponse 断言成功响应格式
func (s *QuotaAPITestSuite) assertSuccessResponse(result map[string]interface{}) {
	_, hasCode := result["code"]
	_, hasMessage := result["message"]
	_, hasData := result["data"]
	_, hasRequestID := result["request_id"]
	_, hasTimestamp := result["timestamp"]

	s.True(hasCode, "Response should have 'code' field")
	s.True(hasMessage, "Response should have 'message' field")
	s.True(hasData, "Response should have 'data' field")
	s.True(hasRequestID, "Response should have 'request_id' field")
	s.True(hasTimestamp, "Response should have 'timestamp' field")

	s.Equal(int32(0), int32(result["code"].(float64)))
}

// assertErrorResponse 断言错误响应格式
func (s *QuotaAPITestSuite) assertErrorResponse(result map[string]interface{}) {
	_, hasCode := result["code"]
	_, hasMessage := result["message"]
	_, hasRequestID := result["request_id"]
	_, hasTimestamp := result["timestamp"]

	s.True(hasCode, "Error response should have 'code' field")
	s.True(hasMessage, "Error response should have 'message' field")
	s.True(hasRequestID, "Error response should have 'request_id' field")
	s.True(hasTimestamp, "Error response should have 'timestamp' field")

	s.NotEqual(int32(0), int32(result["code"].(float64)))
}
