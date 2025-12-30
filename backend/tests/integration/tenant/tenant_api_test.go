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

package tenant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== 类型定义 ====================

/**
 * 创建租户请求
 */
type CreateTenantRequest struct {
	TenantName    string `json:"tenant_name"`
	Description   string `json:"description,omitempty"`
	ContactEmail  string `json:"contact_email"`
	ContactPhone  string `json:"contact_phone,omitempty"`
	CompanyName   string `json:"company_name,omitempty"`
}

/**
 * 创建租户响应
 */
type CreateTenantResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TenantID   string    `json:"tenant_id"`
		TenantName string    `json:"tenant_name"`
		Status     string    `json:"status"`
		CreatedAt  time.Time `json:"created_at"`
	} `json:"data"`
}

/**
 * 列出租户响应
 */
type ListTenantsResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Total    int `json:"total"`
		Current  int `json:"current"`
		PageSize int `json:"page_size"`
		Tenants  []struct {
			TenantID   string    `json:"tenant_id"`
			TenantName string    `json:"tenant_name"`
			Status     string    `json:"status"`
			CreatedAt  time.Time `json:"created_at"`
		} `json:"tenants"`
	} `json:"data"`
}

// ==================== 测试套件 ====================

/**
 * 租户API集成测试套件
 */
func TestTenantAPI(t *testing.T) {
	// 跳过短测试
	if testing.Short() {
		t.Skip("跳过集成测试（使用 -short 标志）")
		return
	}

	// 设置测试容器
	testCtx, err := integration.SetupMySQLTestContainer(t)
	require.NoError(t, err, "测试容器设置失败")
	require.NotNil(t, testCtx.DB, "数据库连接不应为空")

	// 初始化数据库schema
	t.Run("初始化数据库Schema", func(t *testing.T) {
		err := initDatabaseSchema(testCtx.DB)
		require.NoError(t, err, "数据库Schema初始化失败")
	})

	// 运行API测试
	t.Run("创建租户", testCreateTenant(testCtx))
	t.Run("列出租户", testListTenants(testCtx))
	t.Run("获取租户详情", testGetTenant(testCtx))
	t.Run("更新租户", testUpdateTenant(testCtx))
	t.Run("删除租户", testDeleteTenant(testCtx))
}

// ==================== 测试用例 ====================

/**
 * 测试创建租户
 */
func testCreateTenant(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		// 创建测试服务器
		server := createTestServer(testCtx)
		defer server.Close()

		// 构造请求
		reqBody := CreateTenantRequest{
			TenantName:   "测试租户001",
			Description:  "这是一个测试租户",
			ContactEmail: "test@example.com",
			ContactPhone: "13800138000",
			CompanyName:  "测试公司",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err, "JSON序列化失败")

		// 发送请求
		req, err := http.NewRequest(
			"POST",
			server.URL+"/api/v1/tenants",
			bytes.NewBuffer(bodyBytes),
		)
		require.NoError(t, err, "创建HTTP请求失败")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err, "发送HTTP请求失败")
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode, "状态码应为200")

		var respBody CreateTenantResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err, "JSON解码失败")

		assert.Equal(t, 0, respBody.Code, "响应码应为0")
		assert.NotEmpty(t, respBody.Data.TenantID, "租户ID不应为空")
		assert.Equal(t, "测试租户001", respBody.Data.TenantName, "租户名称应匹配")
		assert.Equal(t, "active", respBody.Data.Status, "状态应为active")
	}
}

/**
 * 测试列出租户
 */
func testListTenants(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		server := createTestServer(testCtx)
		defer server.Close()

		// 发送请求
		req, err := http.NewRequest("GET", server.URL+"/api/v1/tenants?page=1&page_size=20", nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody ListTenantsResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		assert.Equal(t, 0, respBody.Code)
		assert.GreaterOrEqual(t, respBody.Data.Total, 0, "总数应>=0")
		assert.Equal(t, 1, respBody.Data.Current, "当前页应为1")
	}
}

/**
 * 测试获取租户详情
 */
func testGetTenant(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		// TODO: 实现获取租户详情测试
		t.Skip("TODO: 待实现")
	}
}

/**
 * 测试更新租户
 */
func testUpdateTenant(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		// TODO: 实现更新租户测试
		t.Skip("TODO: 待实现")
	}
}

/**
 * 测试删除租户
 */
func testDeleteTenant(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		// TODO: 实现删除租户测试
		t.Skip("TODO: 待实现")
	}
}

// ==================== 辅助函数 ====================

/**
 * 初始化数据库Schema
 *
 * 创建租户管理所需的表结构
 */
func initDatabaseSchema(db *integration.TestContext) error {
	schemaSQL := `
	-- 租户表
	CREATE TABLE IF NOT EXISTS tenants (
		tenant_id VARCHAR(36) PRIMARY KEY,
		tenant_name VARCHAR(100) NOT NULL,
		description VARCHAR(500),
		contact_email VARCHAR(255) NOT NULL,
		contact_phone VARCHAR(20),
		company_name VARCHAR(200),
		status VARCHAR(20) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		INDEX idx_status (status),
		INDEX idx_created_at (created_at)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := testCtx.DB.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("执行Schema失败: %w", err)
	}

	return nil
}

/**
 * 创建测试服务器
 *
 * 遵循依赖倒置原则：使用接口而非具体实现
 */
func createTestServer(testCtx *integration.TestContext) *httptest.Server {
	// TODO: 集成实际的HTTP服务器
	// 目前使用mock服务器
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/tenants":
			// 模拟创建租户响应
			response := CreateTenantResponse{
				Code:    0,
				Message: "success",
			}
			response.Data.TenantID = "tenant-123"
			response.Data.TenantName = "测试租户"
			response.Data.Status = "active"
			response.Data.CreatedAt = time.Now()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && r.URL.Path == "/api/v1/tenants":
			// 模拟列出租户响应
			response := ListTenantsResponse{
				Code:    0,
				Message: "success",
			}
			response.Data.Total = 1
			response.Data.Current = 1
			response.Data.PageSize = 20
			response.Data.Tenants = []struct {
				TenantID   string    `json:"tenant_id"`
				TenantName string    `json:"tenant_name"`
				Status     string    `json:"status"`
				CreatedAt  time.Time `json:"created_at"`
			}{
				{
					TenantID:   "tenant-123",
					TenantName: "测试租户",
					Status:     "active",
					CreatedAt:  time.Now(),
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		default:
			http.NotFound(w, r)
		}
	}))
}

/**
 * 创建租户并返回ID的辅助函数
 *
 * 遵循DRY原则：提取公共逻辑，避免重复
 */
func createTenantAndReturnID(t *testing.T, server *httptest.Server, name string) struct {
	TenantID string `json:"tenant_id"`
} {
	reqBody := CreateTenantRequest{
		TenantName:   name,
		Description:  "测试租户",
		ContactEmail:  "test@example.com",
		ContactPhone:  "13800138000",
		CompanyName:   "测试公司",
	}

	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err, "JSON序列化失败")

	req, err := http.NewRequest(
		"POST",
		server.URL+"/api/v1/tenants",
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err, "创建HTTP请求失败")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err, "发送HTTP请求失败")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "创建租户失败")

	var createResp CreateTenantResponse
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	require.NoError(t, err, "JSON解码失败")
	require.NotEmpty(t, createResp.Data.TenantID, "租户ID不应为空")

	return struct {
		TenantID string `json:"tenant_id"`
	}{
		TenantID: createResp.Data.TenantID,
	}
}
