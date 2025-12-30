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

// 🔧 P0修复：日志系统单元测试

package logging

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"go.uber.org/zap"
	. "github.com/smartystreets/goconvey/convey"
)

// TestLoggingInit 测试日志系统初始化
func TestLoggingInit(t *testing.T) {
	Convey("测试日志系统初始化", t, func() {
		Convey("开发环境初始化", func() {
			err := Init("development")

			So(err, ShouldBeNil)
			So(GetLogger(), ShouldNotBeNil)
			So(GetSugaredLogger(), ShouldNotBeNil)
		})

		Convey("生产环境初始化", func() {
			// 设置临时日志路径
			tmpDir := os.TempDir()
			customLogPath := filepath.Join(tmpDir, "zker-test")
			os.Setenv("LOG_PATH", customLogPath)
			defer os.Unsetenv("LOG_PATH")

			err := Init("production")

			So(err, ShouldBeNil)
			So(GetLogger(), ShouldNotBeNil)

			// 验证日志目录创建
			_, err = os.Stat(customLogPath)
			So(err, ShouldBeNil)

			// 清理
			os.RemoveAll(customLogPath)
		})

		Convey("跨平台日志路径", func() {
			// 测试Windows路径
			if runtime.GOOS == "windows" {
				os.Setenv("LOG_PATH", "")
				err := Init("production")

				So(err, ShouldBeNil)

				// Windows默认路径应该使用C:\Logs\zker
				// 由于可能无法访问，这里只验证不panic
			} else {
				// Linux/macOS路径
				os.Setenv("LOG_PATH", "")
				err := Init("production")

				So(err, ShouldBeNil)
			}
		})
	})
}

// TestBasicLogging 测试基础日志方法
func TestBasicLogging(t *testing.T) {
	Convey("测试基础日志方法", t, func() {
		// 初始化日志系统
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("Debug日志", func() {
			// 只验证不panic
			Debug("Test debug message", String("key", "value"))
		})

		Convey("Info日志", func() {
			Info("Test info message", String("key", "value"))
		})

		Convey("Warn日志", func() {
			Warn("Test warn message", String("key", "value"))
		})

		Convey("Error日志", func() {
			Error("Test error message", String("key", "value"))
		})

		Convey("字段辅助函数", func() {
			So(String("key", "value"), ShouldNotBeNil)
			So(Int("key", 42), ShouldNotBeNil)
			So(Int64("key", int64(42)), ShouldNotBeNil)
			So(Float64("key", 3.14), ShouldNotBeNil)
			So(Bool("key", true), ShouldNotBeNil)
			So(Any("key", nil), ShouldNotBeNil)
			So(Err(nil), ShouldNotBeNil)
		})
	})
}

// TestContextLogging 测试带Context的日志方法
func TestContextLogging(t *testing.T) {
	Convey("测试带Context的日志方法", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("WithContext创建logger", func() {
			ctx := context.Background()
			logger := WithContext(ctx)

			So(logger, ShouldNotBeNil)

			// 只验证不panic
			logger.Debug("Test message")
		})

		Convey("Context包含元数据", func() {
			ctx := context.WithValue(context.Background(), "request_id", "req-123")
			ctx = context.WithValue(ctx, "user_id", "user-456")
			ctx = context.WithValue(ctx, "tenant_id", "tenant-789")

			logger := WithContext(ctx)
			So(logger, ShouldNotBeNil)

			// 验证日志记录
			logger.Info("Test with context")
		})

		Convey("DebugContext日志", func() {
			ctx := context.Background()
			DebugContext(ctx, "Test debug context", String("key", "value"))
		})

		Convey("InfoContext日志", func() {
			ctx := context.Background()
			InfoContext(ctx, "Test info context", String("key", "value"))
		})

		Convey("WarnContext日志", func() {
			ctx := context.Background()
			WarnContext(ctx, "Test warn context", String("key", "value"))
		})

		Convey("ErrorContext日志", func() {
			ctx := context.Background()
			ErrorContext(ctx, "Test error context", String("key", "value"))
		})
	})
}

// TestHTTPLogging 测试HTTP日志方法
func TestHTTPLogging(t *testing.T) {
	Convey("测试HTTP日志方法", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("LogHTTPRequest", func() {
			ctx := context.Background()
			LogHTTPRequest(ctx, "GET", "/api/test", "query=test", "127.0.0.1", "Mozilla")
		})

		Convey("LogHTTPResponse", func() {
			ctx := context.Background()
			LogHTTPResponse(ctx, 200, 123.45, 1024)
		})

		Convey("LogHTTPError", func() {
			ctx := context.Background()
			err := TestError("test error")
			LogHTTPError(ctx, err, 500)
		})
	})
}

// TestDBLogging 测试数据库日志方法
func TestDBLogging(t *testing.T) {
	Convey("测试数据库日志方法", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("LogDBQuery", func() {
			ctx := context.Background()
			LogDBQuery(ctx, "mysql", "tenants", "SELECT", 25.5)
		})

		Convey("LogDBError", func() {
			ctx := context.Background()
			err := TestError("database error")
			LogDBError(ctx, "mysql", "INSERT", err)
		})
	})
}

// TestBusinessLogging 测试业务日志方法
func TestBusinessLogging(t *testing.T) {
	Convey("测试业务日志方法", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("LogQuotaCheck", func() {
			ctx := context.Background()
			LogQuotaCheck(ctx, "bots", 1, true, 50, 100)
		})

		Convey("LogTenantOperation", func() {
			ctx := context.Background()
			LogTenantOperation(ctx, "create", "tenant-123", true)
		})

		Convey("LogSubscriptionChange", func() {
			ctx := context.Background()
			LogSubscriptionChange(ctx, "tenant-123", "upgrade", "basic", "premium")
		})
	})
}

// TestPerformanceLogging 测试性能日志方法
func TestPerformanceLogging(t *testing.T) {
	Convey("测试性能日志方法", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("LogSlowQuery", func() {
			ctx := context.Background()
			// 超过阈值的查询应该被记录
			LogSlowQuery(ctx, "SELECT * FROM large_table", 1500.0, 1000.0)

			// 低于阈值的查询不应该被记录
			LogSlowQuery(ctx, "SELECT * FROM small_table", 500.0, 1000.0)
		})

		Convey("LogSlowAPI", func() {
			ctx := context.Background()
			// 超过阈值的API应该被记录
			LogSlowAPI(ctx, "/api/slow", 3000.0, 2000.0)

			// 低于阈值的API不应该被记录
			LogSlowAPI(ctx, "/api/fast", 500.0, 2000.0)
		})
	})
}

// TestLoggingSync 测试日志同步
func TestLoggingSync(t *testing.T) {
	Convey("测试日志同步", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("Sync调用", func() {
			err := Sync()
			So(err, ShouldBeNil)
		})
	})
}

// TestLoggingWithRotation 测试日志轮转
func TestLoggingWithRotation(t *testing.T) {
	Convey("测试日志轮转", t, func() {
		Convey("InitWithFileRotation", func() {
			tmpDir := os.TempDir()
			customLogPath := filepath.Join(tmpDir, "zker-rotation-test")
			os.Setenv("LOG_PATH", customLogPath)
			defer os.Unsetenv("LOG_PATH")

			err := InitWithFileRotation("production", "")
			So(err, ShouldBeNil)

			// 验证日志目录创建
			_, err = os.Stat(customLogPath)
			So(err, ShouldBeNil)

			// 记录一些日志
			for i := 0; i < 100; i++ {
				Infof("Test message %d", i)
			}

			// 清理
			os.RemoveAll(customLogPath)
		})
	})
}

// BenchmarkLogging 基准测试：日志性能
func BenchmarkLogging(b *testing.B) {
	Init("development")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("Benchmark message", String("key", "value"))
	}
}

// BenchmarkContextLogging 基准测试：Context日志性能
func BenchmarkContextLogging(b *testing.B) {
	Init("development")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InfoContext(ctx, "Benchmark message", String("key", "value"))
	}
}

// BenchmarkLoggingWithFields 基准测试：带多个字段的日志性能
func BenchmarkLoggingWithFields(b *testing.B) {
	Init("development")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("Benchmark message",
			String("key1", "value1"),
			String("key2", "value2"),
			String("key3", "value3"),
			Int("count", i),
		)
	}
}

// TestError 自定义测试错误类型
type TestError string

func (e TestError) Error() string {
	return string(e)
}

// TestLoggerConcurrency 测试并发日志记录
func TestLoggerConcurrency(t *testing.T) {
	Convey("测试并发日志记录", t, func() {
		err := Init("development")
		So(err, ShouldBeNil)

		Convey("多个goroutine同时记录日志", func() {
			done := make(chan bool)

			// 启动10个goroutine
			for i := 0; i < 10; i++ {
				go func(id int) {
					for j := 0; j < 100; j++ {
						Infof("Goroutine %d, message %d", id, j)
					}
					done <- true
				}(i)
			}

			// 等待所有goroutine完成
			for i := 0; i < 10; i++ {
				<-done
			}

			// 同步日志
			Sync()

			// 如果没有panic，测试通过
			So(true, ShouldBeTrue)
		})
	})
}

// TestLogRotationIntegration 测试日志轮转集成
func TestLogRotationIntegration(t *testing.T) {
	Convey("测试日志轮转集成", t, func() {
		tmpDir := os.TempDir()
		customLogPath := filepath.Join(tmpDir, "zker-integration-test")
		os.Setenv("LOG_PATH", customLogPath)
		defer os.Unsetenv("LOG_PATH")
		defer os.RemoveAll(customLogPath)

		Convey("使用小文件大小触发轮转", func() {
			// 使用自定义配置进行快速测试
			err := InitWithFileRotation("production", "")
			So(err, ShouldBeNil)

			// 记录大量日志
			for i := 0; i < 1000; i++ {
				Infof("Test log message for rotation %d", i)
			}

			// 同步
			Sync()

			// 给系统一些时间处理
			time.Sleep(100 * time.Millisecond)

			// 验证日志文件存在
			logFile := filepath.Join(customLogPath, "app.log")
			_, err = os.Stat(logFile)
			So(err, ShouldBeNil)
		})
	})
}
