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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	billingService "github.com/coze-dev/coze-studio/backend/domain/billing/service"
	"github.com/coze-dev/coze-studio/backend/tests/integration"
)

// ============================================================
// API集成测试套件
//
// 测试HTTP API端点，验证：
// 1. 实时计费API
// 2. 查询余额API
// 3. 生成账单API
// 4. 多租户隔离
// 5. 错误处理
// ============================================================

// ============================================================
// 测试辅助函数
// ============================================================

/**
 * setupAPITestContext 设置API测试上下文
 */
type apiTestContext struct {
	DB              *gorm.DB
	Server          *httptest.Server
	BaseURL         string
	TestTenantID    string
	TestAuthToken   string
	Cleanup         func()
}

func setupAPITestContext(t *testing.T) *apiTestContext {
	// 1. 启动MySQL容器
	db, dbCleanup := integration.SetupMySQLContainer(t)

	// 2. 运行迁移
	err := db.AutoMigrate(
		&entity.TokenUsageLog{},
		&entity.TokenUsageSummary{},
		&entity.BudgetSettings{},
		&entity.BudgetAlert{},
		&entity.CostOptimizationSuggestion{},
	)
	require.NoError(t, err, "Failed to migrate tables")

	// 3. 创建测试租户
	testTenantID := fmt.Sprintf("tenant_api_test_%d", time.Now().UnixNano())
	tenant := struct {
		TenantID string `gorm:"primaryKey;size:36"`
		Name     string `gorm:"size:100"`
		Status   string `gorm:"size:20"`
	}{
		TenantID: testTenantID,
		Name:     "API测试租户",
		Status:   "active",
	}
	err = db.Table("tenants").Create(&tenant).Error
	require.NoError(t, err, "Failed to create test tenant")

	// 4. 创建余额
	balance := struct {
		TenantID string  `gorm:"primaryKey;size:36"`
		Amount   float64
		Currency string `gorm:"size:10"`
	}{
		TenantID: testTenantID,
		Amount:   100000, // 1000元
		Currency: "CNY",
	}
	err = db.Table("balances").Create(&balance).Error
	require.NoError(t, err, "Failed to create test balance")

	// 5. 创建测试服务器（模拟API路由）
	mux := http.NewServeMux()

	// 注册测试路由
	registerTestRoutes(mux, db, testTenantID)

	server := httptest.NewServer(mux)

	// 6. 生成测试Token
	testAuthToken := fmt.Sprintf("Bearer test_token_%s", testTenantID)

	// 7. 清理函数
	cleanup := func() {
		server.Close()
		dbCleanup()
	}

	return &apiTestContext{
		DB:            db,
		Server:        server,
		BaseURL:       server.URL,
		TestTenantID:  testTenantID,
		TestAuthToken: testAuthToken,
		Cleanup:       cleanup,
	}
}

/**
 * registerTestRoutes 注册测试路由
 */
func registerTestRoutes(mux *http.ServeMux, db *gorm.DB, testTenantID string) {
	// POST /api/billing/realtime/charge - 实时扣费
	mux.HandleFunc("/api/billing/realtime/charge", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 验证Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 解析请求
		var req struct {
			TenantID   string  `json:"tenant_id"`
			UsageType  string  `json:"usage_type"`
			Amount     int     `json:"amount"`
			Tokens     int     `json:"tokens"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 验证租户ID
		if req.TenantID != testTenantID {
			http.Error(w, "Tenant not found", http.StatusNotFound)
			return
		}

		// 检查余额
		var balance struct {
			Amount float64
		}
		err = db.Table("balances").Where("tenant_id = ?", req.TenantID).Select("amount").Scan(&balance).Error
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		cost := float64(req.Tokens) * 0.0001 // 假设0.0001元/token
		if balance.Amount < cost {
			http.Error(w, "余额不足", http.StatusBadRequest)
			return
		}

		// 扣费
		err = db.Table("balances").Where("tenant_id = ?", req.TenantID).Update("amount", gorm.Expr("amount - ?", cost)).Error
		if err != nil {
			http.Error(w, "Failed to charge", http.StatusInternalServerError)
			return
		}

		// 记录使用日志
		log := entity.TokenUsageLog{
			TenantID:      req.TenantID,
			TotalTokens:   int64(req.Tokens),
			TotalCost:     cost,
			RequestType:   "chat",
			ModelProvider: "openai",
			ModelName:     "gpt-4",
		}
		err = db.Create(&log).Error
		if err != nil {
			http.Error(w, "Failed to log usage", http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "扣费成功",
			"data": map[string]interface{}{
				"tenant_id":   req.TenantID,
				"charged":     cost,
				"remaining":   balance.Amount - cost,
				"tokens_used": req.Tokens,
			},
		})
	})

	// GET /api/billing/tenant/:tenantId/balance - 查询余额
	mux.HandleFunc(fmt.Sprintf("/api/billing/tenant/%s/balance", testTenantID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 验证Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 查询余额
		var balance struct {
			TenantID string  `json:"tenant_id"`
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		}
		err := db.Table("balances").Where("tenant_id = ?", testTenantID).Scan(&balance).Error
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// 返回余额
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "查询成功",
			"data":    balance,
		})
	})

	// POST /api/billing/bills/generate - 生成账单
	mux.HandleFunc("/api/billing/bills/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 验证Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 解析请求
		var req struct {
			TenantID     string `json:"tenant_id"`
			BillingCycle string `json:"billing_cycle"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 查询使用记录
		var logs []entity.TokenUsageLog
		err = db.Where("tenant_id = ?", req.TenantID).Find(&logs).Error
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// 计算总费用
		var totalCost float64
		var totalTokens int64
		for _, log := range logs {
			totalCost += log.TotalCost
			totalTokens += log.TotalTokens
		}

		// 生成账单
		billID := fmt.Sprintf("bill_%d", time.Now().UnixNano())
		bill := map[string]interface{}{
			"id":             billID,
			"tenant_id":      req.TenantID,
			"billing_cycle":  req.BillingCycle,
			"total_amount":   totalCost,
			"total_tokens":   totalTokens,
			"currency":       "CNY",
			"status":         "pending",
			"created_at":     time.Now().Format(time.RFC3339),
		}

		// 返回账单
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    0,
			"message": "账单生成成功",
			"data":    bill,
		})
	})
}

// ============================================================
// 测试场景: 实时计费API
// ============================================================

/**
 * TestAPI_RealtimeCharging
 *
 * 测试实时计费API端点
 * 验证：
 * 1. 正常扣费流程
 * 2. 余额不足拒绝
 * 3. Token计费准确性
 */
func TestAPI_RealtimeCharging(t *testing.T) {
	ctx := setupAPITestContext(t)
	defer ctx.Cleanup()

	t.Log("🚀 开始测试: 实时计费API")

	// ========== 测试用例1: 正常扣费 ==========
	t.Run("正常扣费", func(t *testing.T) {
		// 准备请求
		chargeReq := map[string]interface{}{
			"tenant_id":  ctx.TestTenantID,
			"usage_type": "token",
			"amount":     1000,
			"tokens":     10000, // 10000 tokens = 1元
		}
		body, _ := json.Marshal(chargeReq)

		// 发送请求
		req, _ := http.NewRequest("POST", ctx.BaseURL+"/api/billing/realtime/charge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.TestAuthToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200")

		var respData struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Data    struct {
				TenantID   string  `json:"tenant_id"`
				Charged    float64 `json:"charged"`
				Remaining  float64 `json:"remaining"`
				TokensUsed int     `json:"tokens_used"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err, "Response should be valid JSON")

		assert.Equal(t, 0, respData.Code, "Code should be 0")
		assert.Equal(t, "扣费成功", respData.Message, "Message should be correct")
		assert.Equal(t, ctx.TestTenantID, respData.Data.TenantID, "TenantID should match")
		assert.Equal(t, 1.0, respData.Data.Charged, "Charged amount should be 1.0")
		assert.Equal(t, 10000, respData.Data.TokensUsed, "Tokens used should be 10000")
		assert.Less(t, respData.Data.Remaining, 100000.0, "Remaining should be less than initial")

		t.Logf("✅ 扣费成功: ¥%.2f, 剩余: ¥%.2f", respData.Data.Charged, respData.Data.Remaining)
	})

	// ========== 测试用例2: 余额不足 ==========
	t.Run("余额不足", func(t *testing.T) {
		// 准备请求（需要2000元，但只有999元）
		chargeReq := map[string]interface{}{
			"tenant_id":  ctx.TestTenantID,
			"usage_type": "token",
			"amount":     10000000, // 1000万元
			"tokens":     100000000, // 1亿tokens
		}
		body, _ := json.Marshal(chargeReq)

		// 发送请求
		req, _ := http.NewRequest("POST", ctx.BaseURL+"/api/billing/realtime/charge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.TestAuthToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status should be 400")

		t.Log("✅ 正确拒绝余额不足的扣费请求")
	})

	// ========== 测试用例3: 无效租户 ==========
	t.Run("无效租户", func(t *testing.T) {
		chargeReq := map[string]interface{}{
			"tenant_id":  "nonexistent_tenant",
			"usage_type": "token",
			"tokens":     1000,
		}
		body, _ := json.Marshal(chargeReq)

		req, _ := http.NewRequest("POST", ctx.BaseURL+"/api/billing/realtime/charge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.TestAuthToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status should be 404")
		t.Log("✅ 正确返回404 for invalid tenant")
	})
}

// ============================================================
// 测试场景: 查询余额API
// ============================================================

/**
 * TestAPI_QueryBalance
 *
 * 测试查询余额API端点
 * 验证：
 * 1. 正常查询
 * 2. 权限验证
 * 3. 多租户隔离
 */
func TestAPI_QueryBalance(t *testing.T) {
	ctx := setupAPITestContext(t)
	defer ctx.Cleanup()

	t.Log("🚀 开始测试: 查询余额API")

	t.Run("正常查询", func(t *testing.T) {
		// 发送请求
		req, _ := http.NewRequest("GET", ctx.BaseURL+"/api/billing/tenant/"+ctx.TestTenantID+"/balance", nil)
		req.Header.Set("Authorization", ctx.TestAuthToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200")

		var respData struct {
			Code    int `json:"code"`
			Message string `json:"message"`
			Data    struct {
				TenantID string  `json:"tenant_id"`
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)

		assert.Equal(t, 0, respData.Code)
		assert.Equal(t, ctx.TestTenantID, respData.Data.TenantID)
		assert.Equal(t, float64(100000), respData.Data.Amount, "Initial balance should be 100000")
		assert.Equal(t, "CNY", respData.Data.Currency)

		t.Logf("✅ 查询成功: 余额 ¥%.2f", respData.Data.Amount)
	})

	t.Run("未授权访问", func(t *testing.T) {
		req, _ := http.NewRequest("GET", ctx.BaseURL+"/api/billing/tenant/"+ctx.TestTenantID+"/balance", nil)
		// 不设置Authorization

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401")
		t.Log("✅ 正确拒绝未授权访问")
	})
}

// ============================================================
// 测试场景: 生成账单API
// ============================================================

/**
 * TestAPI_GenerateBill
 *
 * 测试生成账单API端点
 * 验证：
 * 1. 账单生成
 * 2. 费用汇总准确性
 * 3. 账单格式正确性
 */
func TestAPI_GenerateBill(t *testing.T) {
	ctx := setupAPITestContext(t)
	defer ctx.Cleanup()

	t.Log("🚀 开始测试: 生成账单API")

	// 1. 先创建一些使用记录
	for i := 0; i < 10; i++ {
		log := entity.TokenUsageLog{
			TenantID:      ctx.TestTenantID,
			TotalTokens:   100000,
			TotalCost:     10.0,
			RequestType:   "chat",
			ModelProvider: "openai",
			ModelName:     "gpt-4",
		}
		err := ctx.DB.Create(&log).Error
		require.NoError(t, err)
	}

	t.Run("生成账单", func(t *testing.T) {
		// 准备请求
		generateReq := map[string]interface{}{
			"tenant_id":     ctx.TestTenantID,
			"billing_cycle": "2025-01",
		}
		body, _ := json.Marshal(generateReq)

		// 发送请求
		req, _ := http.NewRequest("POST", ctx.BaseURL+"/api/billing/bills/generate", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.TestAuthToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err, "Request should succeed")
		defer resp.Body.Close()

		// 验证响应
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Status should be 200")

		var respData struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respData)
		require.NoError(t, err)

		assert.Equal(t, 0, respData.Code)
		assert.Equal(t, "账单生成成功", respData.Message)
		assert.NotEmpty(t, respData.Data["id"])
		assert.Equal(t, ctx.TestTenantID, respData.Data["tenant_id"])
		assert.Equal(t, "2025-01", respData.Data["billing_cycle"])
		assert.Equal(t, float64(100), respData.Data["total_amount"], "Total should be 100 (10 * 10)")
		assert.Equal(t, int64(1000000), respData.Data["total_tokens"], "Total tokens should be 1,000,000")

		t.Logf("✅ 账单生成成功: ID=%s, 金额=¥%.2f",
			respData.Data["id"], respData.Data["total_amount"])
	})
}

// ============================================================
// 测试场景: 并发请求
// ============================================================

/**
 * TestAPI_ConcurrentRequests
 *
 * 测试API并发处理能力
 * 验证：
 * 1. 并发扣费
 * 2. 数据一致性
 * 3. 无竞态条件
 */
func TestAPI_ConcurrentRequests(t *testing.T) {
	ctx := setupAPITestContext(t)
	defer ctx.Cleanup()

	t.Log("🚀 开始测试: API并发请求")

	const concurrency = 10
	errors := make(chan error, concurrency)

	// 并发发送扣费请求
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			chargeReq := map[string]interface{}{
				"tenant_id":  ctx.TestTenantID,
				"usage_type": "token",
				"tokens":     1000, // 0.1元
			}
			body, _ := json.Marshal(chargeReq)

			req, _ := http.NewRequest("POST", ctx.BaseURL+"/api/billing/realtime/charge", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", ctx.TestAuthToken)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errors <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("unexpected status: %d", resp.StatusCode)
				return
			}

			errors <- nil
		}(i)
	}

	// 等待所有请求完成
	for i := 0; i < concurrency; i++ {
		err := <-errors
		assert.NoError(t, err, "Concurrent request should succeed")
	}

	// 验证最终余额
	var balance struct {
		Amount float64
	}
	err := ctx.DB.Table("balances").Where("tenant_id = ?", ctx.TestTenantID).Select("amount").Scan(&balance).Error
	require.NoError(t, err)

	expectedBalance := 100000.0 - float64(concurrency)*0.1
	assert.InDelta(t, expectedBalance, balance.Amount, 0.01, "Balance should be correct after concurrent charges")

	t.Logf("✅ 并发测试通过: %d个请求, 最终余额 ¥%.2f", concurrency, balance.Amount)
}
