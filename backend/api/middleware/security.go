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

// Package middleware 安全中间件(包含CSRF防护、安全头)
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// SecurityHeadersMiddleware 安全响应头中间件
//
// 功能:
//   - X-Content-Type-Options: nosniff (防止MIME类型嗅探)
//   - X-Frame-Options: DENY (防止点击劫持)
//   - X-XSS-Protection: 1; mode=block (启用浏览器XSS过滤)
//   - Referrer-Policy: strict-origin-when-cross-origin (限制引用来源)
//   - Content-Security-Policy: 内容安全策略
//   - Strict-Transport-Security: HSTS (仅HTTPS)
//
// 使用示例:
//
//	r.Use(middleware.SecurityHeadersMiddleware())
func SecurityHeadersMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 防止MIME类型嗅探
		c.Response.Header.Set("X-Content-Type-Options", "nosniff")

		// 防止点击劫持(禁止iframe嵌入)
		c.Response.Header.Set("X-Frame-Options", "DENY")

		// 启用浏览器XSS过滤
		c.Response.Header.Set("X-XSS-Protection", "1; mode=block")

		// 限制引用来源
		c.Response.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// 内容安全策略
		// 注意: 根据实际需求调整CSP策略
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none';"
		c.Response.Header.Set("Content-Security-Policy", csp)

		// HSTS (仅HTTPS)
		if string(c.Request.URI().Scheme()) == "https" {
			c.Response.Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		c.Next(ctx)
	}
}

// CSRFMiddleware CSRF防护中间件
//
// 功能:
//   - 验证POST/PUT/DELETE/PATCH请求的CSRF Token
//   - 跳过GET/HEAD/OPTIONS请求(只读操作)
//   - 支持通过Header或Cookie传递Token
//
// 使用示例:
//
//	r.POST("/api/bots",
//	    middleware.CSRFMiddleware(),
//	    handler.CreateBot)
func CSRFMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 跳过安全方法(GET, HEAD, OPTIONS)
		method := string(c.Request.Method())
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next(ctx)
			return
		}

		// 2. 从Session中获取CSRF Token
		sessionToken := getSessionCSRFToken(c)

		// 3. 从请求头中获取CSRF Token
		requestToken := string(c.GetHeader("X-CSRF-Token"))

		// 4. 验证Token
		if sessionToken == "" || requestToken == "" || sessionToken != requestToken {
			logs.Warnf("[CSRFMiddleware] CSRF token validation failed: sessionToken=%s, requestToken=%s",
				maskToken(sessionToken), maskToken(requestToken))

			c.JSON(http.StatusForbidden, map[string]string{
				"error": "CSRF token validation failed",
				"code":  "ERR_CSRF_TOKEN_INVALID",
			})
			c.Abort()
			return
		}

		logs.Infof("[CSRFMiddleware] CSRF token validated successfully")
		c.Next(ctx)
	}
}

// GenerateCSRFToken 生成CSRF Token
//
// 使用加密安全的随机数生成器生成32字节(64个十六进制字符)的Token
//
// 返回:
//   - token: 生成的CSRF Token
//   - error: 生成失败时返回错误
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32) // 32字节 = 64个十六进制字符
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// getSessionCSRFToken 从Session中获取CSRF Token
//
// 优先级:
//  1. 从Header中读取(用于AJAX请求)
//  2. 从Cookie中读取(用于表单提交)
func getSessionCSRFToken(c *app.RequestContext) string {
	// 尝试从Cookie中读取
	cookie := string(c.Cookie("csrf_token"))
	if cookie != "" {
		return cookie
	}

	// 尝试从Session中读取(如果存在Session机制)
	// 这里需要根据实际的Session实现调整
	return ""
}

// SetCSRFToken 设置CSRF Token到Cookie和Header
//
// 使用示例:
//
//	func GetCSRFToken(ctx context.Context, c *app.RequestContext) {
//	    token, _ := middleware.GenerateCSRFToken()
//	    middleware.SetCSRFToken(c, token)
//	    c.JSON(200, map[string]string{"token": token})
//	}
func SetCSRFToken(c *app.RequestContext, token string) {
	// 设置Cookie (HttpOnly=false, 允许JavaScript读取)
	// SameSite=Strict模式, 防止CSRF攻击
	// Secure=true, 仅HTTPS传输
	maxAge := int(24 * time.Hour / time.Second) // 24小时

	c.SetCookie("csrf_token",
		token,
		maxAge,
		"/",
		"",
		protocol.CookieSameSiteStrictMode,
		true,  // Secure
		false) // HttpOnly=false, 允许前端JS读取
}

// maskToken 掩码Token(用于日志记录,避免泄露完整Token)
func maskToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:4] + "***" + token[len(token)-4:]
}

// GetCSRFTokenHandler 获取CSRF Token的API Handler
//
// 前端调用此接口获取CSRF Token, 然后在后续请求中携带
//
// 路由:
//
//	GET /api/csrf_token
func GetCSRFTokenHandler(ctx context.Context, c *app.RequestContext) {
	token, err := GenerateCSRFToken()
	if err != nil {
		logs.Errorf("[GetCSRFTokenHandler] Failed to generate CSRF token: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate CSRF token",
			"code":  "ERR_CSRF_GENERATE_FAILED",
		})
		return
	}

	// 设置Token到Cookie
	SetCSRFToken(c, token)

	// 返回Token给前端
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

// ValidateCSRFToken 验证CSRF Token(用于手动验证)
//
// 使用场景:
//   - 在某些特殊接口中需要手动验证CSRF Token
//
// 参数:
//   - c: RequestContext
//
// 返回:
//   - bool: Token是否有效
func ValidateCSRFToken(c *app.RequestContext) bool {
	sessionToken := getSessionCSRFToken(c)
	requestToken := string(c.GetHeader("X-CSRF-Token"))

	return sessionToken != "" && sessionToken == requestToken
}

// SafeContentTypeMiddleware 安全Content-Type中间件
//
// 强制API响应使用JSON格式, 避免HTML响应(XSS攻击向量)
//
// 使用示例:
//
//	r.Use(middleware.SafeContentTypeMiddleware())
func SafeContentTypeMiddleware(skipPaths ...string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否跳过
		path := string(c.Request.URI().Path())
		for _, skipPath := range skipPaths {
			if strings.HasPrefix(path, skipPath) {
				c.Next(ctx)
				return
			}
		}

		// 强制JSON响应
		c.Response.Header.Set("Content-Type", "application/json; charset=utf-8")
		c.Response.Header.Set("X-Content-Type-Options", "nosniff")

		c.Next(ctx)
	}
}
