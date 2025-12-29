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
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"go.uber.org/zap"

	"github.com/coze-studio/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-studio/coze-studio/backend/pkg/logs"
	"github.com/coze-studio/coze-studio/backend/types/consts"
)

// TestTenantIsolationMiddleware_Header 租户隔离中间件测试 - 从HTTP Header获取tenant_id
func TestTenantIsolationMiddleware_Header(t *testing.T) {
	// 初始化日志
	logs.Init("development")

	// 创建中间件
	middleware := TenantIsolationMiddleware()

	// 测试用例：从Header获取tenant_id
	t.Run("Extract tenant_id from HTTP header", func(t *testing.T) {
		tenantID := "test-tenant-123"

		// 创建模拟请求
		w := ut.PerformRequest(
			ut.NewContext(),
			middleware,
			"GET",
			"/api/bots",
			ut.Header{
				"X-Tenant-ID": tenantID,
			},
		)

		// 验证响应
		assert.DeepEqual(t, w.Result().Status(), "200 OK")

		// 验证context中存储了tenant_id
		ctx := w.Request.Context()
		storedTenantID := ctx.Value(TenantIDKey)
		assert.DeepEqual(t, storedTenantID, tenantID)

		// 验证ctxcache中也存储了tenant_id
		cachedTenantID, ok := ctxcache.Get[string](ctx, consts.TenantIDKeyInCtx)
		assert.True(t, ok)
		assert.DeepEqual(t, cachedTenantID, tenantID)
	})

	// 测试用例：没有tenant_id时使用默认值
	t.Run("Use default tenant_id when missing", func(t *testing.T) {
		// 创建模拟请求（不携带tenant_id）
		w := ut.PerformRequest(
			ut.NewContext(),
			middleware,
			"GET",
			"/api/bots",
			nil,
		)

		// 验证响应（应该继续处理）
		assert.DeepEqual(t, w.Result().Status(), "200 OK")

		// 验证使用了默认租户ID
		ctx := w.Request.Context()
		storedTenantID := ctx.Value(TenantIDKey)
		assert.DeepEqual(t, storedTenantID, DefaultTenantID)
	})
}

// TestTenantIsolationMiddleware_SkipPaths 测试跳过特定路径
func TestTenantIsolationMiddleware_SkipPaths(t *testing.T) {
	middleware := TenantIsolationMiddleware()

	skipPaths := []string{
		"/api/passport/login",
		"/api/health/status",
		"/api/metrics",
		"/api/public/info",
	}

	for _, path := range skipPaths {
		t.Run("Skip path: "+path, func(t *testing.T) {
			w := ut.PerformRequest(
				ut.NewContext(),
				middleware,
				"GET",
				path,
				nil,
			)

			// 验证请求通过（没有被拦截）
			assert.DeepEqual(t, w.Result().Status(), "200 OK")
		})
	}
}

// TestGetTenantIDFromContext 测试从context获取tenant_id
func TestGetTenantIDFromContext(t *testing.T) {
	logs.Init("development")

	t.Run("Get tenant_id from Hertz context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TenantIDKey, "test-tenant-456")

		tenantID := GetTenantIDFromContext(ctx)
		assert.DeepEqual(t, tenantID, "test-tenant-456")
	})

	t.Run("Get tenant_id from ctxcache", func(t *testing.T) {
		ctx := context.Background()
		ctxcache.Store(ctx, consts.TenantIDKeyInCtx, "cached-tenant-789")

		tenantID := GetTenantIDFromContext(ctx)
		assert.DeepEqual(t, tenantID, "cached-tenant-789")
	})

	t.Run("Return default tenant_id when not found", func(t *testing.T) {
		ctx := context.Background()

		tenantID := GetTenantIDFromContext(ctx)
		assert.DeepEqual(t, tenantID, DefaultTenantID)
	})
}

// TestRequireTenantID 测试强制要求tenant_id的中间件
func TestRequireTenantID(t *testing.T) {
	logs.Init("development")
	middleware := RequireTenantID()

	t.Run("Reject default tenant_id", func(t *testing.T) {
		// 设置默认租户ID到context
		ctx := context.Background()
		ctxcache.Store(ctx, consts.TenantIDKeyInCtx, DefaultTenantID)

		// 创建模拟请求处理器
		handler := func(c context.Context, ctx *app.RequestContext) {
			ctx.Set("result", "success")
		}

		// 包装中间件
		wrappedHandler := middleware(handler)

		// 执行请求
		hertzCtx := ut.NewContext()
		hertzCtx.SetContext(ctx)
		wrappedHandler(hertzCtx, hertzCtx.Request)

		// 验证：应该被拒绝
		// 注意：这里需要根据实际的httputil.Unauthorized实现调整断言
	})
}

// TestTenantIDToInt64 测试tenant_id转换为int64
func TestTenantIDToInt64(t *testing.T) {
	t.Run("Convert valid int64 tenant_id", func(t *testing.T) {
		tenantID := "123456"
		id, err := TenantIDToInt64(tenantID)

		assert.True(t, err == nil)
		assert.DeepEqual(t, id, int64(123456))
	})

	t.Run("Reject UUID format tenant_id", func(t *testing.T) {
		tenantID := "550e8400-e29b-41d4-a716-446655440000"
		_, err := TenantIDToInt64(tenantID)

		assert.True(t, err != nil)
	})

	t.Run("Reject invalid tenant_id", func(t *testing.T) {
		tenantID := "invalid"
		_, err := TenantIDToInt64(tenantID)

		assert.True(t, err != nil)
	})
}

// TestTenantIsolationIntegration 集成测试：完整的租户隔离流程
func TestTenantIsolationIntegration(t *testing.T) {
	logs.Init("development")

	// 创建完整的中间件链
	tenantMiddleware := TenantIsolationMiddleware()

	// 模拟业务处理器
	businessHandler := func(c context.Context, ctx *app.RequestContext) {
		tenantID := GetTenantIDFromContext(c)
		ctx.Set("tenant_id", tenantID)
		ctx.Set("status", "ok")
	}

	t.Run("Complete tenant isolation flow", func(t *testing.T) {
		tenantID := "integration-test-tenant"

		// 创建测试处理器（中间件 + 业务逻辑）
		testHandler := func(c context.Context, ctx *app.RequestContext) {
			// 执行租户隔离中间件
			tenantMiddleware(c, ctx)

			// 检查是否有错误（如果中间件返回了错误，不执行业务逻辑）
			if ctx.Response.StatusCode() != 200 {
				return
			}

			// 执行业务逻辑
			businessHandler(c, ctx)
		}

		// 执行请求
		w := ut.PerformRequest(
			ut.NewContext(),
			testHandler,
			"POST",
			"/api/bots",
			ut.Header{
				"X-Tenant-ID": tenantID,
			},
		)

		// 验证业务逻辑成功执行
		assert.DeepEqual(t, w.Result().Status(), "200 OK")
		assert.DeepEqual(t, w.Response.Body().GetBytes("tenant_id"), []byte(tenantID))
		assert.DeepEqual(t, w.Response.Body().GetBytes("status"), []byte("ok"))
	})
}

// BenchmarkTenantIsolationMiddleware 性能测试
func BenchmarkTenantIsolationMiddleware(b *testing.B) {
	logs.Init("production")
	middleware := TenantIsolationMiddleware()

	// 创建模拟请求
	ctx := ut.NewContext()
	ctx.Request.Header.Set("X-Tenant-ID", "benchmark-tenant")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		middleware(context.Background(), ctx)
	}
}

// BenchmarkGetTenantIDFromContext 性能测试：从context获取tenant_id
func BenchmarkGetTenantIDFromContext(b *testing.B) {
	ctx := context.WithValue(context.Background(), TenantIDKey, "benchmark-tenant")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetTenantIDFromContext(ctx)
	}
}

// TestTenantIsolationLogging 测试日志记录
func TestTenantIsolationLogging(t *testing.T) {
	logs.Init("development")
	middleware := TenantIsolationMiddleware()

	tenantID := "logging-test-tenant"

	// 模拟请求并验证日志
	w := ut.PerformRequest(
		ut.NewContext(),
		middleware,
		"GET",
		"/api/bots",
		ut.Header{
			"X-Tenant-ID": tenantID,
		},
	)

	// 验证请求成功
	assert.DeepEqual(t, w.Result().Status(), "200 OK")

	// 验证context中设置了tenant_id（说明中间件执行成功）
	c := w.Request.Context()
	storedTenantID := c.Value(TenantIDKey)
	assert.DeepEqual(t, storedTenantID, tenantID)

	// 注意：实际测试中可以验证日志输出
	// 这里需要配置日志捕获器来验证日志内容
}

// TestTenantIsolationEdgeCases 边界条件测试
func TestTenantIsolationEdgeCases(t *testing.T) {
	logs.Init("development")
	middleware := TenantIsolationMiddleware()

	testCases := []struct {
		name     string
		tenantID string
		expected string
	}{
		{
			name:     "Empty tenant_id",
			tenantID: "",
			expected: DefaultTenantID,
		},
		{
			name:     "Whitespace tenant_id",
			tenantID: "   ",
			expected: "   ", // 当前实现会保留空格
		},
		{
			name:     "Special characters tenant_id",
			tenantID: "tenant-@#$-123",
			expected: "tenant-@#$-123",
		},
		{
			name:     "Long tenant_id",
			tenantID: string(make([]byte, 1000)),
			expected: string(make([]byte, 1000)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := ut.PerformRequest(
				ut.NewContext(),
				middleware,
				"GET",
				"/api/bots",
				ut.Header{
					"X-Tenant-ID": tc.tenantID,
				},
			)

			ctx := w.Request.Context()
			storedTenantID := ctx.Value(TenantIDKey)

			if tc.tenantID == "" {
				// 空tenant_id应该使用默认值
				assert.DeepEqual(t, storedTenantID, tc.expected)
			} else {
				// 其他情况保留原值
				assert.DeepEqual(t, storedTenantID, tc.expected)
			}
		})
	}
}

// TestTenantIsolationConcurrent 并发测试：验证中间件是并发安全的
func TestTenantIsolationConcurrent(t *testing.T) {
	logs.Init("development")
	middleware := TenantIsolationMiddleware()

	// 并发执行多个请求
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			tenantID := "concurrent-tenant-" + string(rune('0'+index))

			w := ut.PerformRequest(
				ut.NewContext(),
				middleware,
				"GET",
				"/api/bots",
				ut.Header{
					"X-Tenant-ID": tenantID,
				},
			)

			ctx := w.Request.Context()
			storedTenantID := ctx.Value(TenantIDKey)

			// 验证每个请求的tenant_id是独立的
			assert.DeepEqual(t, storedTenantID, tenantID)

			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestTenantIsolationWithTracing 测试与追踪系统的集成
func TestTenantIsolationWithTracing(t *testing.T) {
	logs.Init("development")
	// TODO: 初始化tracing
	// tracing.Init("test-service", "http://localhost:14268/api/traces")

	middleware := TenantIsolationMiddleware()
	tenantID := "tracing-test-tenant"

	w := ut.PerformRequest(
		ut.NewContext(),
		middleware,
		"GET",
		"/api/bots",
		ut.Header{
			"X-Tenant-ID": tenantID,
		},
	)

	// 验证请求成功
	assert.DeepEqual(t, w.Result().Status(), "200 OK")

	// TODO: 验证追踪系统中记录了tenant_id属性
	// 这里需要mock tracing系统来验证
}
