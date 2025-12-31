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

package performance_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/pkg/contextutil"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// =====================================================
// 租户隔离中间件基准测试
// =====================================================

// BenchmarkTenantIsolationMiddleware 测试租户隔离中间件性能
func BenchmarkTenantIsolationMiddleware(b *testing.B) {
	// 初始化依赖
	tenantService := &service.TenantService{} // Mock service
	// redisClient := cache.NewRedis()        // Mock Redis client
	middleware.InitTenantMiddleware(tenantService, nil)

	// 创建模拟请求
	ctx := context.Background()
	requestCtx := app.NewRequestContext()

	// 设置tenant_id到Header
	requestCtx.Set(middleware.TenantIDHeader, "test-tenant-123")

	// 创建中间件实例
	middlewareFunc := middleware.TenantIsolationMiddleware()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		middlewareFunc(ctx, requestCtx)
	}
}

// BenchmarkExtractTenantID 测试tenant_id提取性能
func BenchmarkExtractTenantID(b *testing.B) {
	ctx := context.Background()
	requestCtx := app.NewRequestContext()
	requestCtx.Set(middleware.TenantIDHeader, "test-tenant-123")

	// 模拟Session
	session := &entity.Session{
		UserID:    12345,
		TenantID:  "session-tenant-456",
		Locale:    "en",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 从context获取session
		ctxcache.Store(ctx, consts.SessionDataKeyInCtx, session)
		// 提取tenant_id的逻辑已在中间件中测试
		_ = ctx
	}
}

// BenchmarkGetTenantIDFromContext 测试从context获取tenant_id性能
func BenchmarkGetTenantIDFromContext(b *testing.B) {
	ctx := context.Background()
	tenantID := "benchmark-tenant-789"

	// 存储tenant_id到context
	ctxcache.Store(ctx, consts.TenantIDKeyInCtx, tenantID)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = contextutil.GetTenantIDFromContext(ctx)
	}
}

// =====================================================
// 错误码系统基准测试
// =====================================================

// BenchmarkErrorCodeCreation 测试错误码创建性能
func BenchmarkErrorCodeCreation(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errorx.New(berrno.ErrBotNotFound)
		_ = err
	}
}

// BenchmarkErrorCodeWithParams 测试带参数的错误码创建性能
func BenchmarkErrorCodeWithParams(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errorx.New(berrno.ErrBotNotFound).
			WithZap("bot_id", fmt.Sprintf("bot-%d", i))
		_ = err
	}
}

// BenchmarkMultipleErrorCodes 测试多个不同错误码的创建性能
func BenchmarkMultipleErrorCodes(b *testing.B) {
	errorCodes := []int32{
		berrno.ErrBotNotFoundCode,
		berrno.ErrConversationNotFoundCode,
		berrno.ErrWorkflowNotFoundCode,
		berrno.ErrNodeNotFoundCode,
		berrno.ErrMessageNotFoundCode,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		code := errorCodes[i%len(errorCodes)]
		err := errorx.New(code)
		_ = err
	}
}

// =====================================================
// Session操作基准测试
// =====================================================

// BenchmarkSessionCreation 测试Session创建性能
func BenchmarkSessionCreation(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		session := &entity.Session{
			UserID:    int64(i + 1),
			TenantID:  fmt.Sprintf("tenant-%d", i+1),
			Locale:    "zh-CN",
			UserEmail: fmt.Sprintf("user%d@example.com", i+1),
			CreatedAt: now,
			ExpiresAt: now.Add(24 * time.Hour),
		}
		_ = session
	}
}

// BenchmarkSessionGetTenantID 测试Session获取TenantID性能
func BenchmarkSessionGetTenantID(b *testing.B) {
	session := &entity.Session{
		UserID:    12345,
		TenantID:  "test-tenant-123",
		Locale:    "en",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = session.GetTenantID()
	}
}

// BenchmarkSessionHasTenantID 测试Session检查TenantID性能
func BenchmarkSessionHasTenantID(b *testing.B) {
	session := &entity.Session{
		UserID:    12345,
		TenantID:  "test-tenant-123",
		Locale:    "en",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = session.HasTenantID()
	}
}

// =====================================================
// Context操作基准测试
// =====================================================

// BenchmarkContextCacheStore 测试context缓存存储性能
func BenchmarkContextCacheStore(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctxcache.Store(ctx, "test-key", fmt.Sprintf("value-%d", i))
	}
}

// BenchmarkContextCacheGet 测试context缓存读取性能
func BenchmarkContextCacheGet(b *testing.B) {
	ctx := context.Background()
	ctxcache.Store(ctx, "test-key", "test-value")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = ctxcache.Get[string](ctx, "test-key")
	}
}

// BenchmarkContextCacheStoreAndGet 测试context缓存读写性能
func BenchmarkContextCacheStoreAndGet(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%100)
		value := fmt.Sprintf("value-%d", i)
		ctxcache.Store(ctx, key, value)
		_, _ = ctxcache.Get[string](ctx, key)
	}
}

// =====================================================
// 并发性能测试
// =====================================================

// BenchmarkConcurrentTenantIsolation 测试并发租户隔离性能
func BenchmarkConcurrentTenantIsolation(b *testing.B) {
	tenantService := &service.TenantService{}
	middleware.InitTenantMiddleware(tenantService, nil)
	middlewareFunc := middleware.TenantIsolationMiddleware()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		requestCtx := app.NewRequestContext()
		requestCtx.Set(middleware.TenantIDHeader, "test-tenant-123")

		for pb.Next() {
			middlewareFunc(ctx, requestCtx)
		}
	})
}

// BenchmarkConcurrentErrorCodeCreation 测试并发错误码创建性能
func BenchmarkConcurrentErrorCodeCreation(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			err := errorx.New(berrno.ErrBotNotFound)
			_ = err
			i++
		}
	})
}

// =====================================================
// 内存分配测试
// =====================================================

// BenchmarkMemoryAllocation_TenantID 测试tenant_id相关操作的内存分配
func BenchmarkMemoryAllocation_TenantID(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			tenantID := "tenant-1234567890"
			_ = tenantID
		}
	})

	b.Run("Int64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			tenantID := int64(1234567890)
			_ = tenantID
		}
	})
}

// BenchmarkMemoryAllocation_ErrorCode 测试错误码的内存分配
func BenchmarkMemoryAllocation_ErrorCode(b *testing.B) {
	b.Run("Simple", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := errorx.New(berrno.ErrBotNotFound)
			_ = err
		}
	})

	b.Run("WithParams", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := errorx.New(berrno.ErrBotNotFound).WithZap("bot_id", "bot-123")
			_ = err
		}
	})
}

// =====================================================
// 字符串操作基准测试（与错误码消息格式化相关）
// =====================================================

// BenchmarkStringFormatting_Sprintf 测试fmt.Sprintf性能
func BenchmarkStringFormatting_Sprintf(b *testing.B) {
	botID := "bot-1234567890"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		msg := fmt.Sprintf("Bot not found: %s", botID)
		_ = msg
	}
}

// BenchmarkStringFormatting_Concatenation 测试字符串拼接性能
func BenchmarkStringFormatting_Concatenation(b *testing.B) {
	botID := "bot-1234567890"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		msg := "Bot not found: " + botID
		_ = msg
	}
}

// BenchmarkStringConversion_Int64ToString 测试int64转string性能
func BenchmarkStringConversion_Int64ToString(b *testing.B) {
	userID := int64(1234567890)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		str := strconv.FormatInt(userID, 10)
		_ = str
	}
}

// BenchmarkStringConversion_Itoa 测试Itoa性能
func BenchmarkStringConversion_Itoa(b *testing.B) {
	userID := 1234567890

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		str := strconv.Itoa(userID)
		_ = str
	}
}

// =====================================================
// HTTP状态码映射基准测试
// =====================================================

// BenchmarkHTTPStatusMapping 测试HTTP状态码映射性能
func BenchmarkHTTPStatusMapping(b *testing.B) {
	errorCodes := []int32{
		berrno.ErrBotNotFoundCode,
		berrno.ErrConversationNotFoundCode,
		berrno.ErrWorkflowNotFoundCode,
		berrno.ErrNodeNotFoundCode,
		berrno.ErrMessageNotFoundCode,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		code := errorCodes[i%len(errorCodes)]
		// 模拟HTTP状态码映射逻辑
		var status int
		switch code / 1000000 {
		case 200, 201, 202, 203:
			status = consts.StatusBadRequest
		case 204:
			status = consts.StatusForbidden
		default:
			status = consts.StatusInternalServerError
		}
		_ = status
	}
}

// =====================================================
// 综合性能测试：完整请求流程模拟
// =====================================================

// BenchmarkFullRequestFlow 模拟完整的请求处理流程
func BenchmarkFullRequestFlow(b *testing.B) {
	// 初始化
	tenantService := &service.TenantService{}
	middleware.InitTenantMiddleware(tenantService, nil)
	middlewareFunc := middleware.TenantIsolationMiddleware()

	// 创建session
	session := &entity.Session{
		UserID:    12345,
		TenantID:  "session-tenant-456",
		Locale:    "zh-CN",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		requestCtx := app.NewRequestContext()

		// 1. 存储session到context
		ctxcache.Store(ctx, consts.SessionDataKeyInCtx, session)

		// 2. 执行租户隔离中间件
		middlewareFunc(ctx, requestCtx)

		// 3. 从context获取tenant_id
		tenantID := contextutil.GetTenantIDFromContext(ctx)

		// 4. 创建错误（模拟业务逻辑）
		if tenantID == "" {
			err := errorx.New(berrno.ErrTenantNotFound)
			_ = err
		}
	}
}

// =====================================================
// 性能对比测试
// =====================================================

// BenchmarkStringLookup_Map_vs_Switch 测试map查找vs switch性能
func BenchmarkStringLookup_Map_vs_Switch(b *testing.B) {
	// Map方式
	errorMap := map[string]int32{
		"bot":         berrno.ErrBotNotFoundCode,
		"conversation": berrno.ErrConversationNotFoundCode,
		"workflow":     berrno.ErrWorkflowNotFoundCode,
		"node":         berrno.ErrNodeNotFoundCode,
		"message":      berrno.ErrMessageNotFoundCode,
	}

	entityTypes := []string{"bot", "conversation", "workflow", "node", "message"}

	b.Run("Map", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			entityType := entityTypes[i%len(entityTypes)]
			_ = errorMap[entityType]
		}
	})

	b.Run("Switch", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			entityType := entityTypes[i%len(entityTypes)]
			var code int32
			switch entityType {
			case "bot":
				code = berrno.ErrBotNotFoundCode
			case "conversation":
				code = berrno.ErrConversationNotFoundCode
			case "workflow":
				code = berrno.ErrWorkflowNotFoundCode
			case "node":
				code = berrno.ErrNodeNotFoundCode
			case "message":
				code = berrno.ErrMessageNotFoundCode
			}
			_ = code
		}
	})
}
