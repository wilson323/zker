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

package middleware

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

// TestGenerateCSRFToken 测试CSRF Token生成
func TestGenerateCSRFToken(t *testing.T) {
	token, err := GenerateCSRFToken()

	// 验证生成成功
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 验证Token长度(32字节 = 64个十六进制字符)
	assert.Len(t, token, 64)

	// 验证Token格式(应该是十六进制)
	for _, c := range token {
		assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'))
	}

	// 验证Token唯一性
	token2, err := GenerateCSRFToken()
	assert.NoError(t, err)
	assert.NotEqual(t, token, token2, "生成的Token应该是唯一的")
}

// TestCSRFMiddleware 测试CSRF中间件
func TestCSRFMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		sessionToken   string
		requestToken   string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET请求应该通过",
			method:         "GET",
			sessionToken:   "",
			requestToken:   "",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "HEAD请求应该通过",
			method:         "HEAD",
			sessionToken:   "",
			requestToken:   "",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "OPTIONS请求应该通过",
			method:         "OPTIONS",
			sessionToken:   "",
			requestToken:   "",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "POST请求缺少Token",
			method:         "POST",
			sessionToken:   "",
			requestToken:   "",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "CSRF token validation failed",
		},
		{
			name:           "POST请求Token不匹配",
			method:         "POST",
			sessionToken:   "token123",
			requestToken:   "token456",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "CSRF token validation failed",
		},
		{
			name:           "POST请求Token匹配",
			method:         "POST",
			sessionToken:   "token123",
			requestToken:   "token123",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "PUT请求Token匹配",
			method:         "PUT",
			sessionToken:   "token123",
			requestToken:   "token123",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "DELETE请求Token匹配",
			method:         "DELETE",
			sessionToken:   "token123",
			requestToken:   "token123",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name:           "PATCH请求Token匹配",
			method:         "PATCH",
			sessionToken:   "token123",
			requestToken:   "token123",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试上下文
			ctx := context.Background()
			c := &app.RequestContext{
				Request:  &protocol.Request{},
				Response: &protocol.Response{},
			}

			// 设置请求方法和Header
			c.Request.SetMethod([]byte(tt.method))
			if tt.requestToken != "" {
				c.Request.Header.Set("X-CSRF-Token", tt.requestToken)
			}

			// 设置Session Token(通过Cookie模拟)
			if tt.sessionToken != "" {
				c.Request.SetCookie("csrf_token", tt.sessionToken)
			}

			// 创建测试Handler
			handlerCalled := false
			testHandler := func(ctx context.Context, c *app.RequestContext) {
				handlerCalled = true
				c.String(http.StatusOK, "success")
			}

			// 包装CSRF中间件
			middleware := CSRFMiddleware()
			middleware(ctx, c)

			// 如果中间件没有终止,调用测试Handler
			if !c.Response.Header.IsStatusCode(http.StatusForbidden) {
				testHandler(ctx, c)
			}

			// 验证结果
			assert.Equal(t, tt.expectedStatus, c.Response.StatusCode())
			if tt.expectedBody != "" {
				body := string(c.Response.Body())
				assert.Contains(t, body, tt.expectedBody)
			}

			// 验证Handler是否被调用
			if tt.expectedStatus == http.StatusOK {
				assert.True(t, handlerCalled, "Handler应该被调用")
			} else {
				assert.False(t, handlerCalled, "Handler不应该被调用")
			}
		})
	}
}

// TestSetCSRFToken 测试设置CSRF Token
func TestSetCSRFToken(t *testing.T) {
	c := &app.RequestContext{
		Request:  &protocol.Request{},
		Response: &protocol.Response{},
	}

	token := "test_token_12345"
	SetCSRFToken(c, token)

	// 验证Cookie是否设置
	cookie := c.Response.Header.Get("Set-Cookie")
	assert.NotEmpty(t, cookie)
	assert.Contains(t, cookie, "csrf_token="+token)
	assert.Contains(t, cookie, "SameSite=Strict")
	assert.Contains(t, cookie, "Secure")
	assert.Contains(t, cookie, "HttpOnly=false")
}

// TestMaskToken 测试Token掩码
func TestMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "正常Token",
			token:    "abcdef1234567890",
			expected: "abcd***7890",
		},
		{
			name:     "短Token",
			token:    "abc",
			expected: "***",
		},
		{
			name:     "空Token",
			token:    "",
			expected: "***",
		},
		{
			name:     "正好8个字符",
			token:    "12345678",
			expected: "***",
		},
		{
			name:     "长Token",
			token:    "0123456789abcdefghijklmnopqrstuv",
			expected: "0123***rstuv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskToken(tt.token)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSecurityHeadersMiddleware 测试安全头中间件
func TestSecurityHeadersMiddleware(t *testing.T) {
	ctx := context.Background()
	c := &app.RequestContext{
		Request:  &protocol.Request{},
		Response: &protocol.Response{},
	}

	// 设置HTTPS
	c.Request.SetScheme("https")

	middleware := SecurityHeadersMiddleware()
	middleware(ctx, c)

	// 验证安全头
	headers := c.Response.Header

	assert.Equal(t, "nosniff", string(headers.Get("X-Content-Type-Options")))
	assert.Equal(t, "DENY", string(headers.Get("X-Frame-Options")))
	assert.Equal(t, "1; mode=block", string(headers.Get("X-XSS-Protection")))
	assert.Equal(t, "strict-origin-when-cross-origin", string(headers.Get("Referrer-Policy")))

	csp := string(headers.Get("Content-Security-Policy"))
	assert.NotEmpty(t, csp)
	assert.Contains(t, csp, "default-src 'self'")
	assert.Contains(t, csp, "frame-ancestors 'none'")

	hsts := string(headers.Get("Strict-Transport-Security"))
	assert.NotEmpty(t, hsts)
	assert.Contains(t, hsts, "max-age=31536000")
}

// TestSafeContentTypeMiddleware 测试安全Content-Type中间件
func TestSafeContentTypeMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		skipPaths    []string
		expectHeader bool
	}{
		{
			name:         "正常路径应该设置Header",
			path:         "/api/bots",
			skipPaths:    []string{},
			expectHeader: true,
		},
		{
			name:         "跳过指定路径",
			path:         "/api/public/health",
			skipPaths:    []string{"/api/public"},
			expectHeader: false,
		},
		{
			name:         "前缀匹配",
			path:         "/api/public/health",
			skipPaths:    []string{"/api/public"},
			expectHeader: false,
		},
		{
			name:         "不匹配跳过路径",
			path:         "/api/bots",
			skipPaths:    []string{"/api/public"},
			expectHeader: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := &app.RequestContext{
				Request:  &protocol.Request{},
				Response: &protocol.Response{},
			}

			c.Request.SetRequestURI([]byte(tt.path))

			middleware := SafeContentTypeMiddleware(tt.skipPaths...)
			middleware(ctx, c)

			contentType := string(c.Response.Header.Get("Content-Type"))
			if tt.expectHeader {
				assert.Equal(t, "application/json; charset=utf-8", contentType)
			} else {
				assert.Empty(t, contentType)
			}
		})
	}
}

// TestCSRFAttackVectors 测试CSRF攻击向量
func TestCSRFAttackVectors(t *testing.T) {
	attackVectors := []struct {
		name      string
		token     string
		shouldFail bool
	}{
		{
			name:      "空Token",
			token:     "",
			shouldFail: true,
		},
		{
			name:      "XSS攻击尝试",
			token:     "<script>alert(1)</script>",
			shouldFail: true,
		},
		{
			name:      "SQL注入尝试",
			token:     "' OR '1'='1",
			shouldFail: true,
		},
		{
			name:      "Null字节注入",
			token:     "token\x00",
			shouldFail: true,
		},
		{
			name:      "超长Token",
			token:     strings.Repeat("a", 1000),
			shouldFail: true,
		},
	}

	for _, tt := range attackVectors {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := &app.RequestContext{
				Request:  &protocol.Request{},
				Response: &protocol.Response{},
			}

			c.Request.SetMethod([]byte("POST"))
			c.Request.Header.Set("X-CSRF-Token", tt.token)
			c.Request.SetCookie("csrf_token", tt.token)

			middleware := CSRFMiddleware()
			middleware(ctx, c)

			if tt.shouldFail {
				assert.Equal(t, http.StatusForbidden, c.Response.StatusCode())
			}
		})
	}
}

// BenchmarkGenerateCSRFToken 性能测试
func BenchmarkGenerateCSRFToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateCSRFToken()
	}
}

// BenchmarkCSRFMiddleware 性能测试(GET请求)
func BenchmarkCSRFMiddleware_GET(b *testing.B) {
	ctx := context.Background()
	middleware := CSRFMiddleware()

	for i := 0; i < b.N; i++ {
		c := &app.RequestContext{
			Request: &protocol.Request{},
		}
		c.Request.SetMethod([]byte("GET"))
		middleware(ctx, c)
	}
}

// BenchmarkCSRFMiddleware_POST 性能测试(POST请求)
func BenchmarkCSRFMiddleware_POST(b *testing.B) {
	ctx := context.Background()
	middleware := CSRFMiddleware()

	for i := 0; i < b.N; i++ {
		c := &app.RequestContext{
			Request: &protocol.Request{},
		}
		c.Request.SetMethod([]byte("POST"))
		c.Request.Header.Set("X-CSRF-Token", "test_token")
		c.Request.SetCookie("csrf_token", "test_token")
		middleware(ctx, c)
	}
}
