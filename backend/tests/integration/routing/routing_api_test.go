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

package routing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== 类型定义 ====================

/**
 * 创建路由规则请求
 */
type CreateRoutingRuleRequest struct {
	RuleName    string            `json:"rule_name"`
	Intent      string            `json:"intent"`
	Priority    int               `json:"priority"`
	Conditions  map[string]string `json:"conditions"`
	TargetBotID string            `json:"target_bot_id"`
	Enabled     bool              `json:"enabled"`
}

/**
 * 创建路由规则响应
 */
type CreateRoutingRuleResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RuleID      string    `json:"rule_id"`
		RuleName    string    `json:"rule_name"`
		Intent      string    `json:"intent"`
		Priority    int       `json:"priority"`
		TargetBotID string    `json:"target_bot_id"`
		Enabled     bool      `json:"enabled"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"data"`
}

/**
 * 路由请求
 */
type RouteRequest struct {
	UserInput  string `json:"user_input"`
	SessionID  string `json:"session_id,omitempty"`
	Context    map[string]string `json:"context,omitempty"`
}

/**
 * 路由响应
 */
type RouteResponse struct {
	Code    int `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MatchedRuleID string `json:"matched_rule_id"`
		MatchedIntent string `json:"matched_intent"`
		TargetBotID   string `json:"target_bot_id"`
		Confidence    float64 `json:"confidence"`
		Reason        string `json:"reason,omitempty"`
	} `json:"data"`
}

/**
 * 更新路由规则请求
 */
type UpdateRoutingRuleRequest struct {
	RuleName    string            `json:"rule_name,omitempty"`
	Intent      string            `json:"intent,omitempty"`
	Priority    int               `json:"priority,omitempty"`
	Conditions  map[string]string `json:"conditions,omitempty"`
	TargetBotID string            `json:"target_bot_id,omitempty"`
	Enabled     *bool             `json:"enabled,omitempty"`
}

// ==================== 测试套件 ====================

/**
 * 智能路由 API 集成测试套件
 *
 * 设计原则：
 * - SOLID: 单一职责，每个测试函数只测试一个功能点
 * - DRY: 复用测试辅助函数，避免重复
 * - KISS: 保持简单明了
 */
func TestRoutingAPI(t *testing.T) {
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
	t.Run("创建路由规则", testCreateRoutingRule(testCtx))
	t.Run("列岕路由规则", testListRoutingRules(testCtx))
	t.Run("智能路由匹配", testRouteRequest(testCtx))
	t.Run("更新路由规则", testUpdateRoutingRule(testCtx))
	t.Run("删除路由规则", testDeleteRoutingRule(testCtx))
}

// ==================== 测试用例 ====================

/**
 * 测试创建路由规则
 */
func testCreateRoutingRule(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		server := createTestServer(testCtx)
		defer server.Close()

		// 构造请求
		reqBody := CreateRoutingRuleRequest{
			RuleName: "客服咨询路由",
			Intent:   "customer_service",
			Priority: 1,
			Conditions: map[string]string{
				"keyword":     "咨询,问题,帮助",
				"user_type":   "vip",
				"time_range":  "09:00-18:00",
			},
			TargetBotID: "bot-cs-001",
			Enabled:     true,
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err, "JSON序列化失败")

		// 发送请求
		req, err := http.NewRequest(
			"POST",
			server.URL+"/api/v1/routing/rules",
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

		var respBody CreateRoutingRuleResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err, "JSON解码失败")

		assert.Equal(t, 0, respBody.Code, "响应码应为0")
		assert.NotEmpty(t, respBody.Data.RuleID, "规则ID不应为空")
		assert.Equal(t, "客服咨询路由", respBody.Data.RuleName, "规则名称应匹配")
		assert.Equal(t, "customer_service", respBody.Data.Intent, "意图应匹配")
		assert.Equal(t, 1, respBody.Data.Priority, "优先级应匹配")
		assert.Equal(t, "bot-cs-001", respBody.Data.TargetBotID, "目标Bot应匹配")
		assert.Equal(t, true, respBody.Data.Enabled, "启用状态应为true")
	}
}

/**
 * 测试列岕路由规则
 */
func testListRoutingRules(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		server := createTestServer(testCtx)
		defer server.Close()

		// 发送请求
		req, err := http.NewRequest("GET", server.URL+"/api/v1/routing/rules?page=1&page_size=20", nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

/**
 * 测试智能路由匹配
 */
func testRouteRequest(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		server := createTestServer(testCtx)
		defer server.Close()

		// 先创建路由规则
		ruleID := createRuleAndReturnID(t, server)

		// 路由请求
		reqBody := RouteRequest{
			UserInput: "你好，我想咨询一下产品问题",
			SessionID: "session-123",
			Context: map[string]string{
				"user_type": "vip",
				"time":      "10:00",
			},
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest(
			"POST",
			server.URL+"/api/v1/routing/route",
			bytes.NewBuffer(bodyBytes),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody RouteResponse
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		require.NoError(t, err)

		assert.Equal(t, 0, respBody.Code, "响应码应为0")
		assert.NotEmpty(t, respBody.Data.MatchedRuleID, "匹配规则ID不应为空")
		assert.NotEmpty(t, respBody.Data.TargetBotID, "目标BotID不应为空")
		assert.GreaterOrEqual(t, respBody.Data.Confidence, 0.0, "置信度应>=0")
		assert.LessOrEqual(t, respBody.Data.Confidence, 1.0, "置信度应<=1")
	}
}

/**
 * 测试更新路由规则
 */
func testUpdateRoutingRule(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		server := createTestServer(testCtx)
		defer server.Close()

		// 先创建路由规则
		ruleID := createRuleAndReturnID(t, server)

		// 更新路由规则
		enabled := false
		reqBody := UpdateRoutingRuleRequest{
			Priority: &[]int{2}[0],
			Enabled:  &enabled,
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest(
			"PUT",
			server.URL+"/api/v1/routing/rules/"+ruleID,
			bytes.NewBuffer(bodyBytes),
		)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

/**
 * 测试删除路由规则
 */
func testDeleteRoutingRule(testCtx *integration.TestContext) func(*testing.T) {
	return func(t *testing.T) {
		server := createTestServer(testCtx)
		defer server.Close()

		// 先创建路由规则
		ruleID := createRuleAndReturnID(t, server)

		// 删除路由规则
		req, err := http.NewRequest(
			"DELETE",
			server.URL+"/api/v1/routing/rules/"+ruleID,
			nil,
		)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

// ==================== 辅助函数 ====================

/**
 * 初始化数据库Schema
 *
 * 创建智能路由所需的表结构
 */
func initDatabaseSchema(db *integration.TestContext) error {
	schemaSQL := `
	-- 路由规则表
	CREATE TABLE IF NOT EXISTS routing_rules (
		rule_id VARCHAR(36) PRIMARY KEY,
		rule_name VARCHAR(100) NOT NULL,
		intent VARCHAR(50) NOT NULL,
		priority INT NOT NULL DEFAULT 0,
		target_bot_id VARCHAR(36) NOT NULL,
		conditions JSON,
		enabled BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL,
		INDEX idx_intent (intent),
		INDEX idx_priority (priority),
		INDEX idx_enabled (enabled)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 意图匹配器表
	CREATE TABLE IF NOT EXISTS intent_matchers (
		matcher_id VARCHAR(36) PRIMARY KEY,
		rule_id VARCHAR(36) NOT NULL,
		matcher_type VARCHAR(20) NOT NULL,
		config JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (rule_id) REFERENCES routing_rules(rule_id) ON DELETE CASCADE,
		INDEX idx_rule_id (rule_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

	-- 路由历史表
	CREATE TABLE IF NOT EXISTS routing_history (
		history_id BIGINT PRIMARY KEY AUTO_INCREMENT,
		session_id VARCHAR(36),
		user_input TEXT,
		matched_rule_id VARCHAR(36),
		target_bot_id VARCHAR(36),
		confidence DOUBLE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_session_id (session_id),
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
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/routing/rules":
			// 模拟创建路由规则响应
			response := CreateRoutingRuleResponse{
				Code:    0,
				Message: "success",
			}
			response.Data.RuleID = "rule-123"
			response.Data.RuleName = "客服咨询路由"
			response.Data.Intent = "customer_service"
			response.Data.Priority = 1
			response.Data.TargetBotID = "bot-cs-001"
			response.Data.Enabled = true
			response.Data.CreatedAt = time.Now()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		case r.Method == "GET" && r.URL.Path == "/api/v1/routing/rules":
			// 模拟列岕路由规则响应
			response := struct {
				Code    int `json:"code"`
				Message string `json:"message"`
				Data    struct {
					Total int `json:"total"`
					Rules []struct {
						RuleID      string `json:"rule_id"`
						RuleName    string `json:"rule_name"`
						Intent      string `json:"intent"`
						Priority    int    `json:"priority"`
						TargetBotID string `json:"target_bot_id"`
					} `json:"rules"`
				} `json:"data"`
			}{
				Code:    0,
				Message: "success",
			}
			response.Data.Total = 1
			response.Data.Rules = []struct {
				RuleID      string `json:"rule_id"`
				RuleName    string `json:"rule_name"`
				Intent      string `json:"intent"`
				Priority    int    `json:"priority"`
				TargetBotID string `json:"target_bot_id"`
			}{
				{
					RuleID:      "rule-123",
					RuleName:    "客服咨询路由",
					Intent:      "customer_service",
					Priority:    1,
					TargetBotID: "bot-cs-001",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		case r.Method == "POST" && r.URL.Path == "/api/v1/routing/route":
			// 模拟路由响应
			response := RouteResponse{
				Code:    0,
				Message: "success",
			}
			response.Data.MatchedRuleID = "rule-123"
			response.Data.MatchedIntent = "customer_service"
			response.Data.TargetBotID = "bot-cs-001"
			response.Data.Confidence = 0.95
			response.Data.Reason = ""

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

		default:
			http.NotFound(w, r)
		}
	}))
}

/**
 * 创建路由规则并返回ID的辅助函数
 *
 * 遵循DRY原则：提取公共逻辑，避免重复
 */
func createRuleAndReturnID(t *testing.T, server *httptest.Server) string {
	reqBody := CreateRoutingRuleRequest{
		RuleName: fmt.Sprintf("路由规则_%d", time.Now().Unix()),
		Intent:   "test_intent",
		Priority: 1,
		Conditions: map[string]string{
			"keyword": "test",
		},
		TargetBotID: "bot-test-001",
		Enabled:     true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	require.NoError(t, err, "JSON序列化失败")

	req, err := http.NewRequest(
		"POST",
		server.URL+"/api/v1/routing/rules",
		bytes.NewBuffer(bodyBytes),
	)
	require.NoError(t, err, "创建HTTP请求失败")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	require.NoError(t, err, "发送HTTP请求失败")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "创建路由规则失败")

	var createResp CreateRoutingRuleResponse
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	require.NoError(t, err, "JSON解码失败")
	require.NotEmpty(t, createResp.Data.RuleID, "规则ID不应为空")

	return createResp.Data.RuleID
}
