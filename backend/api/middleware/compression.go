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
	"bytes"
	"compress/gzip"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

const (
	gzipHeader = "Accept-Encoding"
	minSize    = 512 // 最小压缩大小 (字节)
)

// GzipMiddleware Gzip压缩中间件
// 自动压缩大于512字节的JSON响应
// 注意: Hertz框架中Response没有Writer字段,这个中间件采用响应后压缩的方式
func GzipMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查客户端是否支持gzip
		acceptEncoding := string(c.GetHeader(gzipHeader))
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next(ctx)
			return
		}

		// 执行后续处理器
		c.Next(ctx)

		// 如果有错误或响应太小,不压缩
		if c.Response.StatusCode() >= 400 || len(c.Response.Body()) < minSize {
			c.Response.Header.Set("Vary", gzipHeader)
			return
		}

		// 设置gzip响应头
		c.Response.Header.Set("Content-Encoding", "gzip")
		c.Response.Header.Set("Vary", gzipHeader)

		// 执行压缩
		body := c.Response.Body()
		if len(body) > 0 {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)

			if _, err := gz.Write(body); err != nil {
				logs.Errorf("Gzip compression failed: %v", err)
				return
			}

			if err := gz.Close(); err != nil {
				logs.Errorf("Gzip close failed: %v", err)
				return
			}

			// 替换响应体
			c.Response.SetBody(buf.Bytes())
		}
	}
}

// =====================================================
// 轻量级压缩中间件 (仅压缩JSON响应)
// =====================================================

// JSONCompressionMiddleware JSON响应压缩中间件
// 只压缩Content-Type为application/json的响应
func JSONCompressionMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查客户端是否支持gzip
		acceptEncoding := string(c.GetHeader(gzipHeader))
		if !strings.Contains(acceptEncoding, "gzip") {
			c.Next(ctx)
			return
		}

		c.Next(ctx)

		// 只压缩JSON响应
		contentType := string(c.Response.Header.Peek("Content-Type"))
		if !strings.Contains(contentType, "application/json") {
			return
		}

		// 检查响应大小
		body := c.Response.Body()
		if len(body) < minSize {
			return
		}

		// 执行压缩
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		if _, err := gz.Write(body); err != nil {
			logs.Errorf("JSON compression failed: %v", err)
			return
		}

		if err := gz.Close(); err != nil {
			logs.Errorf("JSON gzip close failed: %v", err)
			return
		}

		// 设置压缩响应头
		c.Response.Header.Set("Content-Encoding", "gzip")
		c.Response.Header.Set("Vary", gzipHeader)
		c.Response.SetBody(buf.Bytes())

		logs.Debugf("JSON compressed: %d -> %d bytes (%.1f%%)",
			len(body), buf.Len(), float64(buf.Len())/float64(len(body))*100)
	}
}

// =====================================================
// 使用示例
// =====================================================

/*
// 在路由中注册中间件

import (
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/coze-dev/coze-studio/backend/api/middleware"
)

func main() {
    h := server.Default()

    // 全局启用JSON压缩
    h.Use(middleware.JSONCompressionMiddleware())

    // 或仅对特定路由启用
    api := h.Group("/api")
    api.Use(middleware.JSONCompressionMiddleware())

    api.GET("/bot-store/list", handlers.ListBotStoreItems)

    h.Run(":8888")
}
*/
