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
	"io"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// TenantAPITestSuite 租户API契约测试套件
type TenantAPITestSuite struct {
	suite.Suite
	h *server.Hertz
}

// SetupSuite 设置测试套件
func (s *TenantAPITestSuite) SetupSuite() {
	// 创建Hertz服务器(仅用于测试)
	s.h = server.Default()
}

// TearDownSuite 清理测试套件
func (s *TenantAPITestSuite) TearDownSuite() {
	// 清理资源
}

// TestTenantAPITestSuite 运行测试套件
func TestTenantAPITestSuite(t *testing.T) {
	suite.Run(t, new(TenantAPITestSuite))
}

// TestCreateTenant_Success 测试创建租户成功
func (s *TenantAPITestSuite) TestCreateTenant_Success() {
	// 准备请求数据
	createReq := map[string]interface{}{
		"tenant_name":   "Test Tenant",
		"tenant_type":   "individual",
		"admin_email":   "test@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	// 创建HTTP请求
	w := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应状态码
	assert.DeepEqual(s.T(), 201, w.Code)

	// 验证响应格式
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// 验证必填字段
	s.assertSuccessResponse(result)

	// 验证数据结构
	data := result["data"].(map[string]interface{})
	s.assertTenantData(data)
}

// TestCreateTenant_Conflict 测试租户已存在的冲突
func (s *TenantAPITestSuite) TestCreateTenant_Conflict() {
	// 第一次创建成功
	createReq := map[string]interface{}{
		"tenant_name":   "Duplicate Tenant",
		"tenant_type":   "individual",
		"admin_email":   "duplicate@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	// 创建第一个租户
	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)
	assert.DeepEqual(s.T(), 201, w1.Code)

	// 尝试创建同名租户
	w2 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	// 验证冲突响应
	assert.DeepEqual(s.T(), 409, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertErrorResponse(result)
	assert.DeepEqual(s.T(), errno.ErrTenantAlreadyExistsCode, int32(result["code"].(float64)))
}

// TestGetTenant_Success 测试获取租户详情成功
func (s *TenantAPITestSuite) TestGetTenant_Success() {
	// 先创建一个租户
	createReq := map[string]interface{}{
		"tenant_name":   "Get Test Tenant",
		"tenant_type":   "enterprise",
		"admin_email":   "gettest@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	tenantID := createResp["data"].(map[string]interface{})["tenant_id"].(string)

	// 获取租户详情
	w2 := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants/"+tenantID,
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	data := result["data"].(map[string]interface{})
	assert.DeepEqual(s.T(), tenantID, data["tenant_id"])
	s.assertTenantData(data)
}

// TestGetTenant_NotFound 测试租户不存在
func (s *TenantAPITestSuite) TestGetTenant_NotFound() {
	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants/nonexistent-tenant-id",
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证404响应
	assert.DeepEqual(s.T(), 404, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertErrorResponse(result)
	assert.DeepEqual(s.T(), errno.ErrTenantNotFoundCode, int32(result["code"].(float64)))
}

// TestListTenants_Success 测试获取租户列表成功
func (s *TenantAPITestSuite) TestListTenants_Success() {
	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants",
		ut.WithHeader("Authorization", "Bearer test-token"),
		ut.WithQuery("page_size", "20"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证列表数据结构
	data := result["data"].(map[string]interface{})
	_, hasTenants := data["tenants"]
	_, hasTotalCount := data["total_count"]

	s.True(hasTenants, "Response should have 'tenants' field")
	s.True(hasTotalCount, "Response should have 'total_count' field")
}

// TestListTenants_WithFilters 测试带过滤条件的列表查询
func (s *TenantAPITestSuite) TestListTenants_WithFilters() {
	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants",
		ut.WithHeader("Authorization", "Bearer test-token"),
		ut.WithQuery("status", "active"),
		ut.WithQuery("tenant_type", "enterprise"),
		ut.WithQuery("subscription_tier", "pro"),
		ut.WithQuery("page_size", "10"),
	)

	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)
}

// TestUpdateTenant_Success 测试更新租户成功
func (s *TenantAPITestSuite) TestUpdateTenant_Success() {
	// 先创建租户
	createReq := map[string]interface{}{
		"tenant_name":   "Update Test Tenant",
		"tenant_type":   "individual",
		"admin_email":   "updatetest@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	tenantID := createResp["data"].(map[string]interface{})["tenant_id"].(string)

	// 更新租户
	updateReq := map[string]interface{}{
		"tenant_name": "Updated Tenant Name",
	}
	updateBody, _ := json.Marshal(updateReq)

	w2 := ut.PerformRequest(s.h.Engine, "PUT", "/api/v1/tenants/"+tenantID,
		ut.WithBody(bytes.NewReader(updateBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	assert.DeepEqual(s.T(), 200, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertSuccessResponse(result)
	data := result["data"].(map[string]interface{})
	assert.DeepEqual(s.T(), "Updated Tenant Name", data["tenant_name"])
}

// TestDeleteTenant_Success 测试删除租户成功(软删除)
func (s *TenantAPITestSuite) TestDeleteTenant_Success() {
	// 先创建租户
	createReq := map[string]interface{}{
		"tenant_name":   "Delete Test Tenant",
		"tenant_type":   "individual",
		"admin_email":   "deletetest@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	tenantID := createResp["data"].(map[string]interface{})["tenant_id"].(string)

	// 删除租户
	w2 := ut.PerformRequest(s.h.Engine, "DELETE", "/api/v1/tenants/"+tenantID,
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证204响应
	assert.DeepEqual(s.T(), 204, w2.Code)
}

// TestQuotaCheck_Success 测试配额检查成功
func (s *TenantAPITestSuite) TestQuotaCheck_Success() {
	// 先创建租户
	createReq := map[string]interface{}{
		"tenant_name":   "Quota Test Tenant",
		"tenant_type":   "enterprise",
		"admin_email":   "quotatest@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	tenantID := createResp["data"].(map[string]interface{})["tenant_id"].(string)

	// 检查Bot配额
	checkReq := map[string]interface{}{
		"amount": 1,
	}
	checkBody, _ := json.Marshal(checkReq)

	w2 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/tenants/"+tenantID+"/quotas/bots/check",
		ut.WithBody(bytes.NewReader(checkBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	data := result["data"].(map[string]interface{})
	_, hasAllowed := data["allowed"]
	_, hasRemaining := data["remaining"]

	s.True(hasAllowed, "Response should have 'allowed' field")
	s.True(hasRemaining, "Response should have 'remaining' field")
}

// TestQuotaCheck_Exceeded 测试配额超限
func (s *TenantAPITestSuite) TestQuotaCheck_Exceeded() {
	// TODO: 模拟配额已满的情况
	// 需要先创建租户并使用完配额,然后再次检查
	s.T().Skip("需要模拟配额已满的场景")
}

// TestErrorResponse_Format 测试错误响应格式
func (s *TenantAPITestSuite) TestErrorResponse_Format() {
	// 发送一个无效请求
	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/tenants/invalid-id",
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	// 验证错误响应必填字段
	s.assertErrorResponse(result)
}

// assertSuccessResponse 断言成功响应格式
func (s *TenantAPITestSuite) assertSuccessResponse(result map[string]interface{}) {
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

	// 验证成功码
	s.Equal(int32(0), int32(result["code"].(float64)))
}

// assertErrorResponse 断言错误响应格式
func (s *TenantAPITestSuite) assertErrorResponse(result map[string]interface{}) {
	_, hasCode := result["code"]
	_, hasMessage := result["message"]
	_, hasRequestID := result["request_id"]
	_, hasTimestamp := result["timestamp"]

	s.True(hasCode, "Error response should have 'code' field")
	s.True(hasMessage, "Error response should have 'message' field")
	s.True(hasRequestID, "Error response should have 'request_id' field")
	s.True(hasTimestamp, "Error response should have 'timestamp' field")

	// 错误码应该非0
	s.NotEqual(int32(0), int32(result["code"].(float64)))
}

// assertTenantData 断言租户数据结构
func (s *TenantAPITestSuite) assertTenantData(data map[string]interface{}) {
	requiredFields := []string{
		"tenant_id",
		"tenant_name",
		"tenant_type",
		"status",
		"subscription_tier",
		"created_at",
		"updated_at",
	}

	for _, field := range requiredFields {
		_, hasField := data[field]
		s.True(hasField, "Tenant data should have '"+field+"' field")
	}

	// 验证枚举值
	validTypes := []string{"individual", "team", "enterprise"}
	validStatuses := []string{"active", "suspended", "deleted"}
	validTiers := []string{"free", "pro", "enterprise"}

	tenantType := data["tenant_type"].(string)
	status := data["status"].(string)
	subscriptionTier := data["subscription_tier"].(string)

	s.Contains(validTypes, tenantType, "tenant_type should be valid")
	s.Contains(validStatuses, status, "status should be valid")
	s.Contains(validTiers, subscriptionTier, "subscription_tier should be valid")
}

// BenchmarkCreateTenant 创建租户的性能基准测试
func BenchmarkCreateTenant(b *testing.B) {
	h := server.Default()

	createReq := map[string]interface{}{
		"tenant_name":   "Benchmark Tenant",
		"tenant_type":   "individual",
		"admin_email":   "benchmark@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := ut.PerformRequest(h.Engine, "POST", "/api/v1/tenants",
			ut.WithBody(bytes.NewReader(reqBody)),
			ut.WithHeader("Content-Type", "application/json"),
		)
		if w.Code != 201 {
			b.Errorf("Expected status 201, got %d", w.Code)
		}
	}
}

// BenchmarkGetTenant 获取租户的性能基准测试
func BenchmarkGetTenant(b *testing.B) {
	h := server.Default()

	// 先创建一个租户用于测试
	createReq := map[string]interface{}{
		"tenant_name":   "Benchmark Get Tenant",
		"tenant_type":   "individual",
		"admin_email":   "benchmarkget@example.com",
		"admin_password": "SecurePassword123!",
	}
	reqBody, _ := json.Marshal(createReq)

	w := ut.PerformRequest(h.Engine, "POST", "/api/v1/tenants",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	tenantID := createResp["data"].(map[string]interface{})["tenant_id"].(string)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := ut.PerformRequest(h.Engine, "GET", "/api/v1/tenants/"+tenantID)
		if w.Code != 200 {
			b.Errorf("Expected status 200, got %d", w.Code)
		}
	}
}
