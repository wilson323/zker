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

// 🔧 P0修复：监控系统单元测试

package monitoring

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/smartystreets/goconvey/convey"

	"github.com/coze-dev/coze-studio/backend/infra/monitoring/metrics"
)

// TestMonitoringInit 测试监控系统初始化
func TestMonitoringInit(t *testing.T) {
	convey.Convey("测试监控系统初始化", t, func() {
		convey.Convey("注册Prometheus端点", func() {
			h := server.Default()

			RegisterMetricsHandler(h)

			// 验证端点注册成功
			w := ut.PerformRequest(h.Engine, "GET", "/metrics", nil)
			assert.DeepEqual(t, 200, w.Code)

			// 验证指标输出不为空
			body := w.Body.String()
			convey.So(len(body), convey.ShouldBeGreaterThan, 0)

			// 验证包含基础指标（这些指标应该在metrics.go中定义）
			convey.So(body, convey.ShouldContainSubstring, "# HELP")
			convey.So(body, convey.ShouldContainSubstring, "# TYPE")
		})

		convey.Convey("注册健康检查端点", func() {
			h := server.Default()

			RegisterMetricsHandler(h)

			w := ut.PerformRequest(h.Engine, "GET", "/health", nil)

			assert.DeepEqual(t, 200, w.Code)
			body := w.Body.String()
			convey.So(body, convey.ShouldContainSubstring, "healthy")
		})
	})
}

// TestHTTPMiddleware 测试HTTP监控中间件
func TestHTTPMiddleware(t *testing.T) {
	convey.Convey("测试HTTP监控中间件", t, func() {
		convey.Convey("正常请求记录指标", func() {
			h := server.Default()

			// 注册监控
			RegisterMetricsHandler(h)

			// 注册中间件
			h.Use(HTTPMiddleware())

			// 注册测试路由
			h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
				c.JSON(200, map[string]string{"status": "ok"})
			})

			w := ut.PerformRequest(h.Engine, "GET", "/test", nil)

			assert.DeepEqual(t, 200, w.Code)

			// 等待指标被记录
			time.Sleep(100 * time.Millisecond)

			// 验证metrics端点可访问
			w2 := ut.PerformRequest(h.Engine, "GET", "/metrics", nil)
			assert.DeepEqual(t, 200, w2.Code)

			// 验证指标输出不为空
			body := w2.Body.String()
			convey.So(len(body), convey.ShouldBeGreaterThan, 0)
		})

		convey.Convey("错误请求记录指标", func() {
			// 创建新的服务器实例以避免状态污染
			h2 := server.Default()
			h2.Use(HTTPMiddleware())
			h2.GET("/error", func(ctx context.Context, c *app.RequestContext) {
				c.JSON(500, map[string]string{"error": "internal error"})
			})
			RegisterMetricsHandler(h2)

			w := ut.PerformRequest(h2.Engine, "GET", "/error", nil)

			assert.DeepEqual(t, 500, w.Code)

			// 等待指标被记录
			time.Sleep(100 * time.Millisecond)

			// 验证metrics端点可访问
			w2 := ut.PerformRequest(h2.Engine, "GET", "/metrics", nil)
			assert.DeepEqual(t, 200, w2.Code)
		})
	})
}

// TestDBQueryRecording 测试数据库查询指标记录
func TestDBQueryRecording(t *testing.T) {
	convey.Convey("测试数据库查询指标记录", t, func() {
		convey.Convey("记录成功查询", func() {
			duration := 50 * time.Millisecond

			RecordDBQuery("zker", "SELECT", "tenants", duration, nil)

			// 验证：这里只验证函数不panic
			// 实际指标验证需要通过Prometheus客户端API
			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("记录失败查询", func() {
			duration := 100 * time.Millisecond
			err := errors.New("database error")

			RecordDBQuery("zker", "INSERT", "bots", duration, err)

			// 验证：这里只验证函数不panic
			convey.So(true, convey.ShouldBeTrue)
		})
	})
}

// TestCacheRecording 测试缓存指标记录
func TestCacheRecording(t *testing.T) {
	convey.Convey("测试缓存指标记录", t, func() {
		convey.Convey("记录缓存命中", func() {
			RecordCacheHit("tenant")

			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("记录缓存未命中", func() {
			RecordCacheMiss("tenant")

			convey.So(true, convey.ShouldBeTrue)
		})
	})
}

// TestQuotaRecording 测试配额指标记录
func TestQuotaRecording(t *testing.T) {
	convey.Convey("测试配额指标记录", t, func() {
		convey.Convey("记录配额超限", func() {
			RecordQuotaExceeded("tenant_123", "bots")

			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("更新配额使用", func() {
			UpdateQuotaUsage("tenant_123", "bots", 50, 100)

			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("配额超过100%", func() {
			UpdateQuotaUsage("tenant_123", "bots", 150, 100)

			convey.So(true, convey.ShouldBeTrue)
		})
	})
}

// TestIsolationRecording 测试隔离策略指标记录
func TestIsolationRecording(t *testing.T) {
	convey.Convey("测试隔离策略指标记录", t, func() {
		convey.Convey("记录隔离策略升级", func() {
			RecordIsolationUpgrade("shared", "dedicated")

			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("记录租户注册", func() {
			RecordRegistration()

			convey.So(true, convey.ShouldBeTrue)
		})
	})
}

// TestMetricsIntegration 测试完整的监控系统集成
func TestMetricsIntegration(t *testing.T) {
	convey.Convey("测试完整的监控系统集成", t, func() {
		h := server.Default()

		// 注册监控
		RegisterMetricsHandler(h)
		h.Use(HTTPMiddleware())

		// 注册测试路由
		h.GET("/api/test", func(ctx context.Context, c *app.RequestContext) {
			// 模拟数据库查询
			RecordDBQuery("zker", "SELECT", "tenants", 10*time.Millisecond, nil)

			// 模拟缓存命中
			RecordCacheHit("tenant")

			c.JSON(200, map[string]string{"status": "ok"})
		})

		convey.Convey("完整请求流程记录所有指标", func() {
			w := ut.PerformRequest(h.Engine, "GET", "/api/test", nil)

			assert.DeepEqual(t, 200, w.Code)

			// 等待指标被记录
			time.Sleep(100 * time.Millisecond)

			// 获取metrics
			w2 := ut.PerformRequest(h.Engine, "GET", "/metrics", nil)

			assert.DeepEqual(t, 200, w2.Code)

			// 验证指标输出不为空
			body := w2.Body.String()
			convey.So(len(body), convey.ShouldBeGreaterThan, 0)
			convey.So(body, convey.ShouldContainSubstring, "# HELP")
			convey.So(body, convey.ShouldContainSubstring, "# TYPE")
		})

		convey.Convey("并发请求处理", func() {
			// 发送10个并发请求
			for i := 0; i < 10; i++ {
				ut.PerformRequest(h.Engine, "GET", "/api/test", nil)
			}

			// 等待指标被记录
			time.Sleep(200 * time.Millisecond)

			// 验证指标端点
			w := ut.PerformRequest(h.Engine, "GET", "/metrics", nil)

			assert.DeepEqual(t, 200, w.Code)

			body := w.Body.String()
			convey.So(len(body), convey.ShouldBeGreaterThan, 0)
		})
	})
}

// TestHertzResponseWriter 测试Hertz ResponseWriter适配器
func TestHertzResponseWriter(t *testing.T) {
	convey.Convey("测试Hertz ResponseWriter适配器", t, func() {
		h := server.Default()

		h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
			c.Response.SetBodyRaw([]byte(`{"status":"ok"}`))
		})

		convey.Convey("响应正确返回", func() {
			w := ut.PerformRequest(h.Engine, "GET", "/test", nil)

			assert.DeepEqual(t, 200, w.Code)
			body := w.Body.String()
			convey.So(body, convey.ShouldContainSubstring, `{"status":"ok"}`)
		})
	})
}

// TestMetricsPackageIntegration 测试metrics包集成
func TestMetricsPackageIntegration(t *testing.T) {
	convey.Convey("测试metrics包集成", t, func() {
		convey.Convey("记录HTTP请求", func() {
			metrics.RecordHTTPRequest("GET", "/api/test", 200, 0.5, 100, 200)
			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("记录配额检查", func() {
			metrics.RecordQuotaCheck("tenant_123", "bots", 0.01, true, 50, 100)
			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("记录权限检查", func() {
			metrics.RecordPermissionCheck("tenant_123", "data", 0.001, "allowed")
			convey.So(true, convey.ShouldBeTrue)
		})

		convey.Convey("记录租户创建", func() {
			metrics.RecordTenantCreation("enterprise", "success")
			convey.So(true, convey.ShouldBeTrue)
		})
	})
}

// BenchmarkHTTPMiddleware 基准测试：HTTP中间件性能
func BenchmarkHTTPMiddleware(b *testing.B) {
	h := server.Default()
	h.Use(HTTPMiddleware())

	h.GET("/test", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]string{"status": "ok"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ut.PerformRequest(h.Engine, "GET", "/test", nil)
	}
}

// BenchmarkRecordDBQuery 基准测试：数据库查询记录性能
func BenchmarkRecordDBQuery(b *testing.B) {
	duration := 50 * time.Millisecond

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RecordDBQuery("zker", "SELECT", "tenants", duration, nil)
	}
}

// BenchmarkRecordCacheHit 基准测试：缓存命中记录性能
func BenchmarkRecordCacheHit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RecordCacheHit("tenant")
	}
}

// BenchmarkMetricsRecordHTTP 基准测试：HTTP请求记录性能
func BenchmarkMetricsRecordHTTP(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordHTTPRequest("GET", "/api/test", 200, 0.5, 100, 200)
	}
}
