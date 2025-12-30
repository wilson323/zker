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

// 🔧 P0修复：分布式追踪单元测试

package tracing

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	. "github.com/smartystreets/goconvey/convey"
)

// TestTracingInit 测试追踪系统初始化
func TestTracingInit(t *testing.T) {
	Convey("测试追踪系统初始化", t, func() {
		Convey("基础初始化", func() {
			// 注意：这里使用无效的Jaeger端点，只测试初始化逻辑
			err := Init("test-service", "http://localhost:14268/api/traces")

			// 在测试环境中，Jaeger可能不可用，所以这里只验证不panic
			// 实际环境中应该正确处理连接失败
			So(err, ShouldBeNil) // 或者根据实际实现调整
		})

		Convey("自定义采样率初始化", func() {
			err := InitWithSampler("test-service", "http://localhost:14268/api/traces", 0.5)

			So(err, ShouldBeNil)
		})

		Convey("无效端点处理", func() {
			err := Init("test-service", "http://invalid-host:9999/api/traces")

			// 应该优雅地处理错误，不panic
			// 根据实际实现，这里可能返回错误
			So(err, ShouldNotBeNil)
		})
	})
}

// TestSpanOperations 测试Span操作
func TestSpanOperations(t *testing.T) {
	Convey("测试Span操作", t, func() {
		// 初始化（可能失败，但继续测试）
		_ = Init("test-service", "http://localhost:14268/api/traces")

		Convey("启动和结束Span", func() {
			ctx := context.Background()
			ctx, span := StartSpan(ctx, "test-operation")

			So(ctx, ShouldNotBeNil)
			So(span, ShouldNotBeNil)

			// 模拟一些工作
			time.Sleep(10 * time.Millisecond)

			// 结束span
			span.End()
		})

		Convey("启动带属性的Span", func() {
			ctx := context.Background()
			attrs := []attribute.KeyValue{
				attribute.String("key1", "value1"),
				attribute.String("key2", "value2"),
			}

			ctx, span := StartSpanWithAttributes(ctx, "test-attributes", attrs...)

			So(ctx, ShouldNotBeNil)
			So(span, ShouldNotBeNil)

			span.End()
		})

		Convey("嵌套Span", func() {
			ctx := context.Background()

			// 父span
			ctx, parentSpan := StartSpan(ctx, "parent-operation")
			So(parentSpan, ShouldNotBeNil)

			// 子span
			ctx, childSpan := StartSpan(ctx, "child-operation")
			So(childSpan, ShouldNotBeNil)

			// 模拟工作
			time.Sleep(5 * time.Millisecond)

			childSpan.End()
			parentSpan.End()
		})

		Convey("添加Span属性", func() {
			ctx := context.Background()
			ctx, span := StartSpan(ctx, "test-span")

			// 添加属性
			AddSpanAttributes(ctx,
				attribute.String("user_id", "user-123"),
				attribute.String("tenant_id", "tenant-456"),
				attribute.Int("request_count", 10),
			)

			span.End()
		})

		Convey("添加Span事件", func() {
			ctx := context.Background()
			ctx, span := StartSpan(ctx, "test-span")

			// 添加事件
			AddSpanEvent(ctx, "cache-miss",
				attribute.String("cache_key", "user:123"),
			)

			AddSpanEvent(ctx, "db-query-start",
				attribute.String("query", "SELECT * FROM users"),
			)

			span.End()
		})

		Convey("记录错误", func() {
			ctx := context.Background()
			ctx, span := StartSpan(ctx, "test-span")

			// 模拟错误
			testErr := TestError("test error")
			RecordError(ctx, testErr)

			// 设置错误状态
			SetSpanStatus(ctx, "operation failed")

			span.End()
		})
	})
}

// TestTracingHelpers 测试追踪辅助函数
func TestTracingHelpers(t *testing.T) {
	Convey("测试追踪辅助函数", t, func() {
		Convey("常用属性创建", func() {
			So(WithTenantID("tenant-123"), ShouldNotBeNil)
			So(WithUserID("user-456"), ShouldNotBeNil)
			So(WithRequestID("req-789"), ShouldNotBeNil)

			testErr := TestError("test error")
			So(WithError(testErr), ShouldNotBeNil)
		})

		Convey("常用属性常量", func() {
			So(AttrTenantID, ShouldNotBeNil)
			So(AttrUserID, ShouldNotBeNil)
			So(AttrRequestID, ShouldNotBeNil)
			So(AttrHTTPMethod, ShouldNotBeNil)
			So(AttrHTTPPath, ShouldNotBeNil)
			So(AttrDBSystem, ShouldNotBeNil)
			So(AttrDBName, ShouldNotBeNil)
		})
	})
}

// TestTracingShutdown 测试追踪系统关闭
func TestTracingShutdown(t *testing.T) {
	Convey("测试追踪系统关闭", t, func() {
		Convey("正常关闭", func() {
			_ = Init("test-service", "http://localhost:14268/api/traces")

			ctx := context.Background()
			err := Shutdown(ctx)

			So(err, ShouldBeNil)
		})

		Convey("Flush操作", func() {
			_ = Init("test-service", "http://localhost:14268/api/traces")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := Flush(ctx)
			So(err, ShouldBeNil)
		})
	})
}

// TestRealWorldScenarios 测试真实场景
func TestRealWorldScenarios(t *testing.T) {
	Convey("测试真实场景", t, func() {
		_ = Init("test-service", "http://localhost:14268/api/traces")

		Convey("HTTP请求追踪", func() {
			ctx := context.Background()

			// 启动HTTP请求span
			ctx, span := StartSpanWithAttributes(ctx, "HTTP GET /api/users",
				attribute.String("http.method", "GET"),
				attribute.String("http.path", "/api/users"),
				attribute.String("user_id", "user-123"),
			)

			// 模拟数据库查询
			ctx, dbSpan := StartSpan(ctx, "Database Query")
			AddSpanAttributes(ctx,
				attribute.String("db.system", "mysql"),
				attribute.String("db.name", "opencoze"),
				attribute.String("db.table", "users"),
				attribute.String("db.operation", "SELECT"),
			)

			time.Sleep(10 * time.Millisecond)
			dbSpan.End()

			// 模拟缓存查询
			ctx, cacheSpan := StartSpan(ctx, "Cache Get")
			AddSpanAttributes(ctx,
				attribute.String("cache.type", "redis"),
				attribute.String("cache.key", "user:123"),
			)

			time.Sleep(5 * time.Millisecond)
			cacheSpan.End()

			// 添加事件
			AddSpanEvent(ctx, "request-completed",
				attribute.Int("http.status_code", 200),
			)

			span.End()
		})

		Convey("租户操作追踪", func() {
			ctx := context.Background()

			// 租户上下文
			ctx = context.WithValue(ctx, "tenant_id", "tenant-456")
			ctx = context.WithValue(ctx, "user_id", "user-789")

			// 启动span
			ctx, span := StartSpanWithAttributes(ctx, "TenantOperation",
				WithTenantID("tenant-456"),
				WithUserID("user-789"),
				attribute.String("operation", "create_bot"),
			)

			// 数据库操作
			ctx, dbSpan := StartSpan(ctx, "Create Bot")
			AddSpanAttributes(ctx,
				AttrDBSystem.String("mysql"),
				AttrDBName.String("opencoze"),
				AttrDBTable.String("bots"),
				AttrDBOperation.String("INSERT"),
			)

			time.Sleep(15 * time.Millisecond)

			// 检查配额
			AddSpanEvent(ctx, "quota-check",
				attribute.String("resource_type", "bots"),
				attribute.Bool("allowed", true),
				attribute.Int("current_usage", 5),
				attribute.Int("quota_limit", 10),
			)

			dbSpan.End()
			span.End()
		})
	})
}

// TestConcurrentSpans 测试并发Span
func TestConcurrentSpans(t *testing.T) {
	Convey("测试并发Span", t, func() {
		_ = Init("test-service", "http://localhost:14268/api/traces")

		Convey("多个并发Span", func() {
			ctx := context.Background()
			done := make(chan bool)

			// 启动10个并发操作
			for i := 0; i < 10; i++ {
				go func(id int) {
					// 每个goroutine有自己的span
					_, span := StartSpan(ctx, "concurrent-operation")
					AddSpanAttributes(ctx, attribute.Int("worker_id", id))

					time.Sleep(10 * time.Millisecond)
					span.End()

					done <- true
				}(i)
			}

			// 等待所有goroutine完成
			for i := 0; i < 10; i++ {
				<-done
			}
		})
	})
}

// BenchmarkSpanCreation 基准测试：Span创建性能
func BenchmarkSpanCreation(b *testing.B) {
	_ = Init("test-service", "http://localhost:14268/api/traces")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := StartSpan(ctx, "benchmark-operation")
		span.End()
	}
}

// BenchmarkSpanWithAttributes 基准测试：带属性Span创建性能
func BenchmarkSpanWithAttributes(b *testing.B) {
	_ = Init("test-service", "http://localhost:14268/api/traces")
	ctx := context.Background()

	attrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.Int("key3", 42),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := StartSpanWithAttributes(ctx, "benchmark-operation", attrs...)
		span.End()
	}
}

// BenchmarkSpanEvents 基准测试：Span事件性能
func BenchmarkSpanEvents(b *testing.B) {
	_ = Init("test-service", "http://localhost:14268/api/traces")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := StartSpan(ctx, "benchmark-operation")
		AddSpanEvent(ctx, "test-event", attribute.Int("iteration", i))
		span.End()
	}
}

// TestError 自定义测试错误类型
type TestError string

func (e TestError) Error() string {
	return string(e)
}
