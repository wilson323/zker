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
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// SecurityHeadersConfig 安全头配置
type SecurityHeadersConfig struct {
	// XSS 保护
	XSSProtection string // "1; mode=block"

	// 内容类型选项
	ContentTypeOptions string // "nosniff"

	// X-Frame-Options (点击劫持保护)
	XFrameOptions string // "DENY" or "SAMEORIGIN"

	// 严格传输安全
	HSTSMaxAge           int
	HSTSIncludeSubDomains bool

	// 内容安全策略
	ContentSecurityPolicy string

	// Referrer策略
	ReferrerPolicy string

	// 权限策略
	PermissionsPolicy string

	// 跨域资源隔离
	CrossOriginEmbedderPolicy string
	CrossOriginOpenerPolicy   string
	CrossOriginResourcePolicy string

	// 自定义安全头
	CustomHeaders map[string]string
}

// DefaultSecurityHeadersConfig 默认安全头配置
func DefaultSecurityHeadersConfig() *SecurityHeadersConfig {
	return &SecurityHeadersConfig{
		XSSProtection: "1; mode=block",
		ContentTypeOptions: "nosniff",
		XFrameOptions: "SAMEORIGIN",
		HSTSMaxAge: 31536000, // 1年
		HSTSIncludeSubDomains: true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'",
		ReferrerPolicy: "strict-origin-when-cross-origin",
		PermissionsPolicy: "geolocation=(), microphone=(), camera=()",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy: "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		CustomHeaders: make(map[string]string),
	}
}

// SecurityHeaders 安全头中间件
// 遵循OWASP安全最佳实践
func SecurityHeaders(config *SecurityHeadersConfig) app.HandlerFunc {
	if config == nil {
		config = DefaultSecurityHeadersConfig()
	}

	return func(ctx context.Context, c *app.RequestContext) {
		// XSS保护
		if config.XSSProtection != "" {
			c.Header("X-XSS-Protection", config.XSSProtection)
		}

		// 内容类型嗅探保护
		if config.ContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", config.ContentTypeOptions)
		}

		// 点击劫持保护
		if config.XFrameOptions != "" {
			c.Header("X-Frame-Options", config.XFrameOptions)
		}

		// HTTP严格传输安全(HSTS)
		if config.HSTSMaxAge > 0 {
			hstsValue := "max-age=" + strconv.Itoa(config.HSTSMaxAge)
			if config.HSTSIncludeSubDomains {
				hstsValue += "; includeSubDomains"
			}
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// 内容安全策略
		if config.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
		}

		// Referrer策略
		if config.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", config.ReferrerPolicy)
		}

		// 权限策略
		if config.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", config.PermissionsPolicy)
		}

		// 跨域资源隔离(COOP/COEP/COEP)
		if config.CrossOriginEmbedderPolicy != "" {
			c.Header("Cross-Origin-Embedder-Policy", config.CrossOriginEmbedderPolicy)
		}
		if config.CrossOriginOpenerPolicy != "" {
			c.Header("Cross-Origin-Opener-Policy", config.CrossOriginOpenerPolicy)
		}
		if config.CrossOriginResourcePolicy != "" {
			c.Header("Cross-Origin-Resource-Policy", config.CrossOriginResourcePolicy)
		}

		// 自定义安全头
		for key, value := range config.CustomHeaders {
			c.Header(key, value)
		}

		// 移除敏感信息
		removeSensitiveHeaders(c)

		c.Next(ctx)
	}
}

// removeSensitiveHeaders 移除可能泄露敏感信息的响应头
func removeSensitiveHeaders(c *app.RequestContext) {
	// 移除服务器版本信息
	c.DelHeader("Server")
	c.DelHeader("X-Powered-By")

	// 移除后端技术栈信息
	c.DelHeader("X-AspNet-Version")
	c.DelHeader("X-AspNetMvc-Version")

	// 移除开发服务器信息
	c.DelHeader("X-Debug")
}

// SecurityHeadersStrict 严格安全头配置(生产环境)
func SecurityHeadersStrict() app.HandlerFunc {
	config := &SecurityHeadersConfig{
		XSSProtection: "1; mode=block",
		ContentTypeOptions: "nosniff",
		XFrameOptions: "DENY", // 完全禁止iframe嵌入
		HSTSMaxAge: 63072000, // 2年
		HSTSIncludeSubDomains: true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self'; object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'; upgrade-insecure-requests",
		ReferrerPolicy: "no-referrer",
		PermissionsPolicy: "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=()",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy: "same-origin",
		CrossOriginResourcePolicy: "same-origin",
	}

	return SecurityHeaders(config)
}

// SecurityHeadersDev 开发环境安全头配置
func SecurityHeadersDev() app.HandlerFunc {
	config := &SecurityHeadersConfig{
		XSSProtection: "1; mode=block",
		ContentTypeOptions: "nosniff",
		XFrameOptions: "SAMEORIGIN",
		HSTSMaxAge: 0, // 开发环境不启用HSTS
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' 'unsafe-eval'; img-src 'self' data: http://localhost:*",
		ReferrerPolicy: "strict-origin-when-cross-origin",
	}

	return SecurityHeaders(config)
}

// SecurityHeadersMinimal 最小安全头配置
func SecurityHeadersMinimal() app.HandlerFunc {
	config := &SecurityHeadersConfig{
		XSSProtection: "1; mode=block",
		ContentTypeOptions: "nosniff",
		XFrameOptions: "DENY",
	}

	return SecurityHeaders(config)
}

// NoIndex 禁止搜索引擎索引
// 用途:测试环境、管理后台等不需要被搜索引擎索引的页面
func NoIndex() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("X-Robots-Tag", "noindex, nofollow, nosnippet, noarchive")
		c.Next(ctx)
	}
}

// XContentTypeOptions 设置X-Content-Type-Options头
func XContentTypeOptions() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Next(ctx)
	}
}

// XFrameOptions 设置X-Frame-Options头
func XFrameOptions(options string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if options == "" {
			options = "DENY"
		}
		c.Header("X-Frame-Options", options)
		c.Next(ctx)
	}
}

// ContentSecurityPolicy 设置内容安全策略
func ContentSecurityPolicy(policy string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if policy == "" {
			policy = "default-src 'self'"
		}
		c.Header("Content-Security-Policy", policy)
		c.Next(ctx)
	}
}

// CustomSecurityHeader 添加自定义安全头
func CustomSecurityHeader(key, value string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header(key, value)
		c.Next(ctx)
	}
}

// SecurityAuditLog 安全审计日志
// 用途:记录可疑请求、安全事件
func SecurityAuditLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Path())
		method := string(c.Method())
		userAgent := string(c.GetHeader("User-Agent"))
		referer := string(c.GetHeader("Referer"))

		// 检测可疑请求
		isSuspicious := detectSuspiciousRequest(path, userAgent, referer)

		if isSuspicious {
			logs.CtxWarnf(ctx, "[SecurityAudit] Suspicious request: method=%s path=%s ua=%s referer=%s ip=%s",
				method,
				path,
				userAgent,
				referer,
				c.ClientIP(),
			)
		}

		c.Next(ctx)
	}
}

// detectSuspiciousRequest 检测可疑请求
func detectSuspiciousRequest(path, userAgent, referer string) bool {
	// 检测路径遍历攻击
	if strings.Contains(path, "..") || strings.Contains(path, "%2e%2e") {
		return true
	}

	// 检测SQL注入特征
	sqlKeywords := []string{"union", "select", "insert", "update", "delete", "drop", "exec"}
	pathLower := strings.ToLower(path)
	for _, keyword := range sqlKeywords {
		if strings.Contains(pathLower, keyword) {
			return true
		}
	}

	// 检测空User-Agent(可能是爬虫)
	if userAgent == "" || userAgent == "-" {
		return true
	}

	// 检测常见攻击工具的特征
	attackTools := []string{
		"sqlmap", "nikto", "nmap", "masscan", "zap", "burp",
		"metasploit", "w3af", "hydra", "medusa", "john",
	}
	userAgentLower := strings.ToLower(userAgent)
	for _, tool := range attackTools {
		if strings.Contains(userAgentLower, tool) {
			return true
		}
	}

	return false
}
