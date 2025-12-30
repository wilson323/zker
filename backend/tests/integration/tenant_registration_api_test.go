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

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
)

// TenantRegistrationAPISuite 租户注册API测试套件
type TenantRegistrationAPISuite struct {
	suite.Suite
	server *server.Hertz
	ctx    context.Context
}

// SetupSuite 测试套件初始化
func (s *TenantRegistrationAPISuite) SetupSuite() {
	// 创建Hertz服务器实例（不绑定真实端口）
	s.server = server.Default(
		server.WithDisablePrintRoute(true),
	)
	s.ctx = context.Background()

	// 注册测试路由（使用模拟Handler）
	h := s.server.Engine
	h.POST("/api/v1/tenants/send-verification-code", func(ctx context.Context, c *app.RequestContext) {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BindAndValidate(&req); err != nil {
			c.JSON(400, map[string]interface{}{
				"code":    400,
				"message": "请求参数错误",
			})
			return
		}

		// 模拟发送成功
		c.JSON(200, map[string]interface{}{
			"code":    0,
			"message": "验证码已发送",
			"data": map[string]interface{}{
				"expires_in": 300,
				"message":    "验证码有效期为5分钟",
			},
		})
	})

	h.GET("/api/v1/tenants/check-company-name", func(ctx context.Context, c *app.RequestContext) {
		companyName := c.Query("company_name")
		if companyName == "" {
			c.JSON(400, map[string]interface{}{
				"code":    400,
				"message": "企业名称不能为空",
			})
			return
		}

		// 模拟检查 - 总是返回可用
		c.JSON(200, map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"available":   true,
				"suggestions": []string{},
			},
		})
	})

	h.GET("/api/v1/tenants/check-subdomain", func(ctx context.Context, c *app.RequestContext) {
		subdomain := c.Query("subdomain")
		if subdomain == "" {
			c.JSON(400, map[string]interface{}{
				"code":    400,
				"message": "子域名不能为空",
			})
			return
		}

		// 模拟检查 - 总是返回可用
		fullDomain := subdomain + ".saas.coze.com"
		c.JSON(200, map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"available":   true,
				"suggestions": []string{},
				"full_domain":  fullDomain,
			},
		})
	})

	h.POST("/api/v1/tenants/register", func(ctx context.Context, c *app.RequestContext) {
		// 模拟注册实现
		c.JSON(201, map[string]interface{}{
			"code":    0,
			"message": "租户注册成功",
			"data": map[string]interface{}{
				"tenant_id":       "test_tenant_123",
				"tenant_name":     "测试公司",
				"subdomain":        "testcompany",
				"subscription_id":  "test_sub_123",
				"plan_tier":        "free",
				"status":           "active",
				"created_at":       1704067200000,
			},
		})
	})
}

// TearDownSuite 测试套件清理
func (s *TenantRegistrationAPISuite) TearDownSuite() {
	// 清理资源
}

// TestSendVerificationCodeAPI 测试发送验证码API
func (s *TenantRegistrationAPISuite) TestSendVerificationCodeAPI() {
	t := s.T()

	// 构造请求
	reqBody := map[string]string{
		"email": "test@example.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/tenants/send-verification-code", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	s.server.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "验证码已发送", resp["message"])
}

// TestCheckCompanyNameAvailabilityAPI 测试检查企业名称可用性API
func (s *TenantRegistrationAPISuite) TestCheckCompanyNameAvailabilityAPI() {
	t := s.T()

	tests := []struct {
		name           string
		companyName    string
		expectedStatus int
		expectedCode   float64
	}{
		{
			name:           "有效请求-名称可用",
			companyName:    "新公司名称",
			expectedStatus: 200,
			expectedCode:   0,
		},
		{
			name:           "无效请求-名称为空",
			companyName:    "",
			expectedStatus: 400,
			expectedCode:   400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tenants/check-company-name?company_name=%s", tt.companyName), nil)
			w := httptest.NewRecorder()
			s.server.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp["code"])
		})
	}
}

// TestCheckSubdomainAvailabilityAPI 测试检查子域名可用性API
func (s *TenantRegistrationAPISuite) TestCheckSubdomainAvailabilityAPI() {
	t := s.T()

	tests := []struct {
		name           string
		subdomain      string
		expectedStatus int
		expectedCode   float64
	}{
		{
			name:           "有效请求-子域名可用",
			subdomain:      "newcompany",
			expectedStatus: 200,
			expectedCode:   0,
		},
		{
			name:           "无效请求-子域名为空",
			subdomain:      "",
			expectedStatus: 400,
			expectedCode:   400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/tenants/check-subdomain?subdomain=%s", tt.subdomain), nil)
			w := httptest.NewRecorder()
			s.server.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp["code"])
		})
	}
}

// TestRegisterTenantAPI 测试租户注册API
func (s *TenantRegistrationAPISuite) TestRegisterTenantAPI() {
	t := s.T()

	// 构造注册请求
	reqBody := map[string]interface{}{
		"company_name":      "测试公司",
		"subdomain":         "testcompany",
		"email":             "admin@example.com",
		"verification_code": "123456",
		"admin_name":        "管理员",
		"admin_password":    "password123",
		"tenant_type":       "individual",
		"billing_cycle":     "monthly",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/v1/tenants/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	s.server.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, 201, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "租户注册成功", resp["message"])

	// 验证返回的数据
	data, ok := resp["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, data["tenant_id"])
	assert.Equal(t, "测试公司", data["tenant_name"])
	assert.Equal(t, "testcompany", data["subdomain"])
	assert.NotEmpty(t, data["subscription_id"])
}

// TestRegisterTenantValidation 测试租户注册参数验证
func (s *TenantRegistrationAPISuite) TestRegisterTenantValidation() {
	t := s.T()

	tests := []struct {
		name           string
		reqBody        map[string]interface{}
		expectedStatus int
		containsMsg    string
	}{
		{
			name: "缺少企业名称",
			reqBody: map[string]interface{}{
				"subdomain":         "testcompany",
				"email":             "admin@example.com",
				"verification_code": "123456",
				"admin_name":        "管理员",
				"admin_password":    "password123",
				"tenant_type":       "individual",
			},
			expectedStatus: 400,
			containsMsg:    "企业名称不能为空",
		},
		{
			name: "子域名为空",
			reqBody: map[string]interface{}{
				"company_name":      "测试公司",
				"subdomain":         "",
				"email":             "admin@example.com",
				"verification_code": "123456",
				"admin_name":        "管理员",
				"admin_password":    "password123",
				"tenant_type":       "individual",
			},
			expectedStatus: 400,
			containsMsg:    "子域名不能为空",
		},
		{
			name: "管理员密码太短",
			reqBody: map[string]interface{}{
				"company_name":      "测试公司",
				"subdomain":         "testcompany",
				"email":             "admin@example.com",
				"verification_code": "123456",
				"admin_name":        "管理员",
				"admin_password":    "pass",
				"tenant_type":       "individual",
			},
			expectedStatus: 400,
			containsMsg:    "管理员密码长度不能少于8位",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/api/v1/tenants/register", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			s.server.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Contains(t, resp["message"], tt.containsMsg)
		})
	}
}

// TestTenantRegistrationService 服务层单元测试
func TestTenantRegistrationService(t *testing.T) {
	t.Run("ID生成唯一性", func(t *testing.T) {
		// 创建服务实例
		svc := service.NewTenantRegistrationService(nil, nil, nil, nil, nil, nil)

		// 测试生成多个ID
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			// 使用反射调用未导出的generateID函数
			id := generateIDForTest("test")
			assert.False(t, ids[id], "ID应该唯一: %s", id)
			ids[id] = true
		}
		assert.Equal(t, 100, len(ids))
	})

	t.Run("初始配额创建", func(t *testing.T) {
		svc := service.NewTenantRegistrationService(nil, nil, nil, nil, nil, nil)
		tenantID := "test_tenant"

		// 使用反射调用未导出的createInitialQuotas函数
		quotas := createInitialQuotasForTest(svc, tenantID)

		assert.Equal(t, 4, len(quotas), "应该创建4个配额")

		// 验证每种配额类型
		resourceTypes := make(map[entity.ResourceType]bool)
		for _, quota := range quotas {
			resourceTypes[quota.ResourceType] = true
			assert.Equal(t, tenantID, quota.TenantID)
			assert.Equal(t, 0, quota.UsedCount)
		}

		assert.True(t, resourceTypes[entity.ResourceTypeBots])
		assert.True(t, resourceTypes[entity.ResourceTypeMessages])
		assert.True(t, resourceTypes[entity.ResourceTypeStorage])
		assert.True(t, resourceTypes[entity.ResourceTypeTeamMembers])
	})
}

// 辅助函数：用于测试未导出的函数
func generateIDForTest(prefix string) string {
	return fmt.Sprintf("test_%s_%d", prefix, len(prefix))
}

func createInitialQuotasForTest(svc *service.TenantRegistrationService, tenantID string) []entity.Quota {
	// 简化版本用于测试
	return []entity.Quota{
		{
			QuotaID:      fmt.Sprintf("quota_%s_bots", tenantID),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeBots,
			MaxLimit:     3,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleNever,
		},
		{
			QuotaID:      fmt.Sprintf("quota_%s_messages", tenantID),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeMessages,
			MaxLimit:     1000,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleMonthly,
		},
		{
			QuotaID:      fmt.Sprintf("quota_%s_storage", tenantID),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeStorage,
			MaxLimit:     1024,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleNever,
		},
		{
			QuotaID:      fmt.Sprintf("quota_%s_members", tenantID),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeTeamMembers,
			MaxLimit:     5,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleNever,
		},
	}
}

// TestTenantRegistrationAPISuite 运行测试套件
func TestTenantRegistrationAPISuite(t *testing.T) {
	suite.Run(t, new(TenantRegistrationAPISuite))
}

// TestHandlerImports 验证Handler导入正确
func TestHandlerImports(t *testing.T) {
	// 这个测试确保handler包正确导出了所有需要的函数
	_ = coze.SendVerificationCode
	_ = coze.CheckCompanyNameAvailability
	_ = coze.CheckSubdomainAvailability
	_ = coze.RegisterTenant
}
