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
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// CORSConfig CORS配置
type CORSConfig struct {
	// AllowOrigins 允许的源
	// 支持: * (所有源), 具体域名 (如: https://example.com), 前缀匹配 (如: https://*.example.com)
	AllowOrigins []string

	// AllowMethods 允许的HTTP方法
	AllowMethods []string

	// AllowHeaders 允许的请求头
	AllowHeaders []string

	// ExposeHeaders 暴露的响应头
	ExposeHeaders []string

	// AllowCredentials 是否允许携带凭证
	AllowCredentials bool

	// MaxAge 预检请求缓存时间(秒)
	MaxAge int
}

// DefaultCORSConfig 默认CORS配置
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8888",
			"https://*.coze.com",
		},
		AllowMethods: []string{
			consts.MethodGet,
			consts.MethodPost,
			consts.MethodPut,
			consts.MethodDelete,
			consts.MethodPatch,
			consts.MethodOptions,
			consts.MethodHead,
		},
		AllowHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Authorization",
			"Content-Type",
			"Origin",
			"User-Agent",
			"X-Tenant-ID",
			"X-Request-ID",
			"X-Trace-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
			"X-Trace-ID",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24小时
	}
}

// CORS CORS中间件
// 遵循企业级规范:支持动态源检查、预检请求、凭证处理
func CORS(config *CORSConfig) app.HandlerFunc {
	if config == nil {
		config = DefaultCORSConfig()
	}

	// 预处理方法名为大写
	methodsMap := make(map[string]bool)
	for _, method := range config.AllowMethods {
		methodsMap[method] = true
	}

	return func(ctx context.Context, c *app.RequestContext) {
		origin := string(c.GetHeader("Origin"))
		method := string(c.Method())

		// 检查源是否允许
		allowed := isOriginAllowed(origin, config.AllowOrigins)
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 允许的方法
		if len(config.AllowMethods) > 0 {
			c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowMethods, ", "))
		}

		// 允许的请求头
		if len(config.AllowHeaders) > 0 {
			c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowHeaders, ", "))
		}

		// 暴露的响应头
		if len(config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposeHeaders, ", "))
		}

		// 是否允许携带凭证
		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 预检请求缓存时间
		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", string(rune(config.MaxAge)))
		}

		// 处理预检请求
		if method == consts.MethodOptions {
			logs.CtxInfof(ctx, "[CORS] Preflight request from origin: %s", origin)
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}

		// 如果源不允许,返回403
		if !allowed && origin != "" {
			logs.CtxWarnf(ctx, "[CORS] Origin not allowed: %s", origin)
			c.AbortWithStatus(consts.StatusForbidden)
			return
		}

		c.Next(ctx)
	}
}

// isOriginAllowed 检查源是否允许
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return true // 没有Origin头的请求(如同源请求)总是允许
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}

		if allowed == origin {
			return true
		}

		// 支持通配符前缀匹配
		if strings.HasSuffix(allowed, "*") {
			prefix := strings.TrimSuffix(allowed, "*")
			if strings.HasPrefix(origin, prefix) {
				return true
			}
		}
	}

	return false
}

// CORSPermissive 宽松CORS配置(开发环境)
func CORSPermissive() app.HandlerFunc {
	config := &CORSConfig{
		AllowOrigins:      []string{"*"},
		AllowMethods:      []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:      []string{"*"},
		AllowCredentials:  false,
		MaxAge:            86400,
	}

	return CORS(config)
}

// CORSStrict 严格CORS配置(生产环境)
func CORSStrict(allowedOrigins []string) app.HandlerFunc {
	config := &CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{
			consts.MethodGet,
			consts.MethodPost,
			consts.MethodPut,
			consts.MethodDelete,
			consts.MethodPatch,
			consts.MethodOptions,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Tenant-ID",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	return CORS(config)
}
