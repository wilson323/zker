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

// Package middleware 安全中间件集成测试
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/test/mock"
	"github.com/stretchr/testify/require"
)

// TestCSRFMiddlewareIntegration 集成测试：CSRF中间件完整流程
func TestCSRFMiddlewareIntegration(t *testing.T) {
	// 创建测试服务器
	h := server.NewHFServer(t)

	// 注册中间件
	h.Use(SecurityHeadersMiddleware())
	h.Use(CSRFMiddleware())

	// 注册测试路由
	h.POST("/api/bot/create", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	h.GET("/api/bot/list", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	h.GET("/api/csrf_token", GetCSRFTokenHandler)

	t.Run("完整流程：获取Token -> 使用Token创建Bot", func(t *testing.T) {
		// Step 1: 获取CSRF Token
		req1 := httptest.NewRequest("GET", "/api/csrf_token", nil)
		w1 := httptest.NewRecorder()
		h.ServeHTTP(w1, req1)

		assert.DeepEqual(t, http.StatusOK, w1.Code)
		resp1 := mock.GetResponse(w1)
		tokenData := map[string]string{}
		err := resp1.BindJSON(&tokenData)
		require.NoError(t, err)
		require.NotEmpty(t, tokenData["token"])

		// 验证Cookie已设置
		cookies := w1.Result().Cookies()
		var csrfCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "csrf_token" {
				csrfCookie = c
				break
			}
		}
		require.NotNil(t, csrfCookie)
		assert.DeepEqual(t, tokenData["token"], csrfCookie.Value)

		// Step 2: 使用Token创建Bot（成功）
		req2 := httptest.NewRequest("POST", "/api/bot/create", strings.NewReader(`{"name":"test"}`))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("X-CSRF-Token", tokenData["token"])
		// 携带Cookie
		req2.AddCookie(csrfCookie)

		w2 := httptest.NewRecorder()
		h.ServeHTTP(w2, req2)

		assert.DeepEqual(t, http.StatusOK, w2.Code)
	})

	t.Run("完整流程：GET请求不需要CSRF Token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/bot/list", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.DeepEqual(t, http.StatusOK, w.Code)
	})

	t.Run("完整流程：没有CSRF Token的POST请求被拒绝", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/bot/create", strings.NewReader(`{"name":"test"}`))
		req.Header.Set("Content-Type", "application/json")
		// 故意不设置X-CSRF-Token

		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		assert.DeepEqual(t, http.StatusForbidden, w.Code)
		resp := mock.GetResponse(w)
		errorData := map[string]string{}
		err := resp.BindJSON(&errorData)
		require.NoError(t, err)
		assert.DeepEqual(t, "ERR_CSRF_TOKEN_INVALID", errorData["code"])
	})

	t.Run("完整流程：CSRF Token不匹配被拒绝", func(t *testing.T) {
		// Step 1: 获取Token
		req1 := httptest.NewRequest("GET", "/api/csrf_token", nil)
		w1 := httptest.NewRecorder()
		h.ServeHTTP(w1, req1)

		resp1 := mock.GetResponse(w1)
		tokenData := map[string]string{}
		err := resp1.BindJSON(&tokenData)
		require.NoError(t, err)

		// Step 2: 使用错误的Token
		req2 := httptest.NewRequest("POST", "/api/bot/create", strings.NewReader(`{"name":"test"}`))
		req2.Header.Set("Content-Type", "application/json")
		req2.Header.Set("X-CSRF-Token", "wrong-token-"+tokenData["token"])
		// Cookie中的Token是正确的
		cookies := w1.Result().Cookies()
		for _, c := range cookies {
			req2.AddCookie(c)
		}

		w2 := httptest.NewRecorder()
		h.ServeHTTP(w2, req2)

		assert.DeepEqual(t, http.StatusForbidden, w2.Code)
	})
}

// TestCSRFMiddlewareWithAllHTTPMethods 测试所有HTTP方法
func TestCSRFMiddlewareWithAllHTTPMethods(t *testing.T) {
	h := server.NewHFServer(t)
	h.Use(CSRFMiddleware())

	// 注册所有HTTP方法的测试路由
	h.GET("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"method": "GET"})
	})
	h.POST("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"method": "POST"})
	})
	h.PUT("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"method": "PUT"})
	})
	h.PATCH("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"method": "PATCH"})
	})
	h.DELETE("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"method": "DELETE"})
	})
	h.HEAD("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.Status(200)
	})
	h.OPTIONS("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.Status(200)
	})

	testCases := []struct {
		method       string
		needsCSRF    bool
		expectedCode int
	}{
		{"GET", false, 200},
		{"HEAD", false, 200},
		{"OPTIONS", false, 200},
		{"POST", true, 403},
		{"PUT", true, 403},
		{"PATCH", true, 403},
		{"DELETE", true, 403},
	}

	for _, tc := range testCases {
		t.Run(tc.method+"请求", func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/test", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			if tc.needsCSRF {
				// 写操作没有CSRF Token应该返回403
				assert.DeepEqual(t, http.StatusForbidden, w.Code)
			} else {
				// 读操作应该正常返回
				assert.DeepEqual(t, tc.expectedCode, w.Code)
			}
		})
	}
}

// TestSecurityHeadersMiddlewareIntegration 测试安全响应头
func TestSecurityHeadersMiddlewareIntegration(t *testing.T) {
	h := server.NewHFServer(t)
	h.Use(SecurityHeadersMiddleware())

	h.GET("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()

	// 验证所有安全头
	assert.DeepEqual(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.DeepEqual(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.DeepEqual(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.DeepEqual(t, "strict-origin-when-cross-origin", resp.Header.Get("Referrer-Policy"))

	// 验证CSP头（应该包含default-src 'self'）
	csp := resp.Header.Get("Content-Security-Policy")
	assert.Contains(t, csp, "default-src 'self'")
}

// BenchmarkCSRFMiddleware 性能测试：CSRF中间件开销
func BenchmarkCSRFMiddleware(b *testing.B) {
	h := server.NewHFServer(b)
	h.Use(CSRFMiddleware())

	h.POST("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	// 预热
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("POST", "/api/test", strings.NewReader(`{}`))
		req.Header.Set("X-CSRF-Token", "test-token")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
	}

	// 基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/test", strings.NewReader(`{}`))
		req.Header.Set("X-CSRF-Token", "test-token")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
	}
}

// BenchmarkCSRFMiddlewareWithTokenGeneration 性能测试：包含Token生成
func BenchmarkCSRFMiddlewareWithTokenGeneration(b *testing.B) {
	b.Run("GenerateCSRFToken", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GenerateCSRFToken()
		}
	})
}

// TestCSRFMiddlewareConcurrency 并发测试
func TestCSRFMiddlewareConcurrency(t *testing.T) {
	h := server.NewHFServer(t)
	h.Use(CSRFMiddleware())

	h.POST("/api/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	// 并发发送100个请求
	done := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func() {
			req := httptest.NewRequest("POST", "/api/test", strings.NewReader(`{}`))
			req.Header.Set("X-CSRF-Token", "test-token")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			// 所有请求应该被拒绝（没有有效的Cookie）
			assert.DeepEqual(t, http.StatusForbidden, w.Code)

			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 100; i++ {
		<-done
	}
}

// 辅助函数：创建测试HTTP服务器
type server struct {
	*app.Server
}

// NewHFServer 创建Hertz测试服务器
func NewHFServer(t testing.TB) *server {
	s := &server{}
	s.Server = app.NewServer()
	return s
}
