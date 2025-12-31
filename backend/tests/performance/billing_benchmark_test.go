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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// =====================================================
// 计费系统基准测试
// =====================================================

// BenchmarkRealtimeCharge 测试实时扣费性能
// 目标: ns/op < 100000 (0.1ms/op), 0 B/op, 0 allocs/op
func BenchmarkRealtimeCharge(b *testing.B) {
	ctx := context.Background()

	// 模拟扣费请求
	chargeRequest := struct {
		TenantID   string
		UsageType  string
		Amount     int64
		Metadata   map[string]string
	}{
		TenantID:  "benchmark-tenant-001",
		UsageType: "token",
		Amount:    1000,
		Metadata: map[string]string{
			"bot_id":     "bot-123",
			"user_id":    "user-456",
			"request_id": "req-789",
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟实时扣费逻辑（无数据库操作）
		err := simulateRealtimeCharge(ctx, chargeRequest)
		require.NoError(b, err)
	}
}

// simulateRealtimeCharge 模拟实时扣费逻辑
func simulateRealtimeCharge(ctx context.Context, req struct {
	TenantID  string
	UsageType string
	Amount    int64
	Metadata  map[string]string
}) error {
	// 1. 验证租户状态
	if req.TenantID == "" {
		return errorx.New(berrno.ErrTenantNotFound)
	}

	// 2. 验证使用类型
	if req.UsageType != "token" && req.UsageType != "api_call" {
		return errorx.New(berrno.ErrInvalidParameter)
	}

	// 3. 验证金额
	if req.Amount <= 0 {
		return errorx.New(berrno.ErrInvalidParameter)
	}

	// 4. 计算费用（简化逻辑）
	cost := req.Amount

	// 5. 检查余额（模拟）
	balance := int64(1000000) // 模拟余额
	if balance < cost {
		return errorx.New(berrno.ErrInsufficientQuota)
	}

	// 6. 扣除余额（模拟）
	_ = cost

	return nil
}

// BenchmarkGetBalance 测试查询余额性能
// 目标: ns/op < 50000 (0.05ms/op), 0 B/op, 0 allocs/op
func BenchmarkGetBalance(b *testing.B) {
	ctx := context.Background()
	tenantID := "benchmark-tenant-001"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		balance, err := simulateGetBalance(ctx, tenantID)
		require.NoError(b, err)
		require.NotNil(b, balance)
	}
}

// simulateGetBalance 模拟查询余额逻辑
func simulateGetBalance(ctx context.Context, tenantID string) (*struct {
	TenantID      string
	TokenBalance  int64
	ApiCallBalance int64
	UpdatedAt     time.Time
}, error) {
	if tenantID == "" {
		return nil, errorx.New(berrno.ErrTenantNotFound)
	}

	return &struct {
		TenantID       string
		TokenBalance   int64
		ApiCallBalance int64
		UpdatedAt      time.Time
	}{
		TenantID:       tenantID,
		TokenBalance:   1000000,
		ApiCallBalance: 50000,
		UpdatedAt:      time.Now(),
	}, nil
}

// BenchmarkGenerateBill 测试生成账单性能
// 目标: ns/op < 5000000 (5ms/op), 10000 B/op, 100 allocs/op
func BenchmarkGenerateBill(b *testing.B) {
	ctx := context.Background()

	request := struct {
		TenantID     string
		BillingCycle string
	}{
		TenantID:     "benchmark-tenant-001",
		BillingCycle: "2025-01",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bill, err := simulateGenerateBill(ctx, request)
		require.NoError(b, err)
		require.NotNil(b, bill)
	}
}

// simulateGenerateBill 模拟生成账单逻辑
func simulateGenerateBill(ctx context.Context, req struct {
	TenantID     string
	BillingCycle string
}) (*struct {
	BillID       string
	TenantID     string
	BillingCycle string
	TokenUsage   int64
	ApiCallUsage int64
	TotalAmount  int64
	CreatedAt    time.Time
}, error) {
	if req.TenantID == "" {
		return nil, errorx.New(berrno.ErrTenantNotFound)
	}

	if req.BillingCycle == "" {
		return nil, errorx.New(berrno.ErrInvalidParameter)
	}

	// 模拟聚合使用数据
	tokenUsage := int64(500000) // 50万 tokens
	apiCallUsage := int64(1000)  // 1000次API调用

	// 计算费用
	tokenAmount := tokenUsage * 2 / 1000        // 每个token 0.002元
	apiCallAmount := apiCallUsage * 5           // 每次API调用 0.005元
	totalAmount := tokenAmount + apiCallAmount

	return &struct {
		BillID       string
		TenantID     string
		BillingCycle string
		TokenUsage   int64
		ApiCallUsage int64
		TotalAmount  int64
		CreatedAt    time.Time
	}{
		BillID:       fmt.Sprintf("bill-%s-%s", req.TenantID, req.BillingCycle),
		TenantID:     req.TenantID,
		BillingCycle: req.BillingCycle,
		TokenUsage:   tokenUsage,
		ApiCallUsage: apiCallUsage,
		TotalAmount:  totalAmount,
		CreatedAt:    time.Now(),
	}, nil
}

// BenchmarkConcurrentRealtimeCharge 测试并发实时扣费性能
// 目标: 1000 并发下，ns/op < 200000 (0.2ms/op)
func BenchmarkConcurrentRealtimeCharge(b *testing.B) {
	ctx := context.Background()
	chargeRequest := struct {
		TenantID  string
		UsageType string
		Amount    int64
		Metadata  map[string]string
	}{
		TenantID:  "benchmark-tenant-001",
		UsageType: "token",
		Amount:    1000,
		Metadata: map[string]string{
			"bot_id":     "bot-123",
			"user_id":    "user-456",
			"request_id": "req-789",
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			req := chargeRequest
			req.Metadata["request_id"] = fmt.Sprintf("req-%d", i)
			err := simulateRealtimeCharge(ctx, req)
			require.NoError(b, err)
			i++
		}
	})
}

// =====================================================
// 批量操作基准测试
// =====================================================

// BenchmarkBatchRealtimeCharge 测试批量实时扣费性能
// 目标: 100次批量操作，平均每次 < 0.15ms
func BenchmarkBatchRealtimeCharge(b *testing.B) {
	ctx := context.Background()
	batchSize := 100

	requests := make([]struct {
		TenantID  string
		UsageType string
		Amount    int64
		Metadata  map[string]string
	}, batchSize)

	for i := 0; i < batchSize; i++ {
		requests[i] = struct {
			TenantID  string
			UsageType string
			Amount    int64
			Metadata  map[string]string
		}{
			TenantID:  fmt.Sprintf("benchmark-tenant-%d", i%10),
			UsageType: "token",
			Amount:    1000,
			Metadata: map[string]string{
				"bot_id":     fmt.Sprintf("bot-%d", i%50),
				"user_id":    fmt.Sprintf("user-%d", i%100),
				"request_id": fmt.Sprintf("req-%d", i),
			},
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, req := range requests {
			err := simulateRealtimeCharge(ctx, req)
			require.NoError(b, err)
		}
	}
}

// BenchmarkBatchGetBalance 测试批量查询余额性能
// 目标: 100次批量操作，平均每次 < 0.05ms
func BenchmarkBatchGetBalance(b *testing.B) {
	ctx := context.Background()
	batchSize := 100

	tenantIDs := make([]string, batchSize)
	for i := 0; i < batchSize; i++ {
		tenantIDs[i] = fmt.Sprintf("benchmark-tenant-%d", i%10)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, tenantID := range tenantIDs {
			balance, err := simulateGetBalance(ctx, tenantID)
			require.NoError(b, err)
			require.NotNil(b, balance)
		}
	}
}

// =====================================================
// 缓存性能测试
// =====================================================

// BenchmarkCacheHit 测试缓存命中性能
// 目标: ns/op < 10000 (0.01ms/op), 0 B/op, 0 allocs/op
func BenchmarkCacheHit(b *testing.B) {
	ctx := context.Background()
	tenantID := "benchmark-tenant-001"

	// 模拟缓存
	cache := make(map[string]*struct {
		TokenBalance   int64
		ApiCallBalance int64
		UpdatedAt      time.Time
	})

	// 预热缓存
	cache[tenantID] = &struct {
		TokenBalance   int64
		ApiCallBalance int64
		UpdatedAt      time.Time
	}{
		TokenBalance:   1000000,
		ApiCallBalance: 50000,
		UpdatedAt:      time.Now(),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		data, ok := cache[tenantID]
		require.True(b, ok)
		require.NotNil(b, data)
	}
}

// BenchmarkCacheMiss 测试缓存未命中性能
// 目标: ns/op < 5000 (0.005ms/op), 0 B/op, 0 allocs/op
func BenchmarkCacheMiss(b *testing.B) {
	ctx := context.Background()
	tenantID := "benchmark-tenant-not-exist"

	cache := make(map[string]*struct {
		TokenBalance   int64
		ApiCallBalance int64
		UpdatedAt      time.Time
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, ok := cache[tenantID]
		require.False(b, ok)
	}
}

// =====================================================
// 错误码性能测试
// =====================================================

// BenchmarkBillingErrorCodeCreation 测试计费相关错误码创建性能
// 目标: ns/op < 5000 (0.005ms/op), 0 B/op, 0 allocs/op
func BenchmarkBillingErrorCodeCreation(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errorx.New(berrno.ErrInsufficientQuota)
		_ = err
	}
}

// BenchmarkBillingErrorCodeWithParams 测试带参数的计费错误码创建性能
// 目标: ns/op < 10000 (0.01ms/op), < 100 B/op, < 5 allocs/op
func BenchmarkBillingErrorCodeWithParams(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := errorx.New(berrno.ErrInsufficientQuota).
			WithZap("tenant_id", fmt.Sprintf("tenant-%d", i)).
			WithZap("required", 1000).
			WithZap("available", 500)
		_ = err
	}
}

// =====================================================
// 数据转换性能测试
// =====================================================

// BenchmarkBillingRecordToDTO 测试计费记录转DTO性能
// 目标: ns/op < 50000 (0.05ms/op), < 1000 B/op, < 50 allocs/op
func BenchmarkBillingRecordToDTO(b *testing.B) {
	type BillingRecord struct {
		TenantID     string
		UsageType    string
		Amount       int64
		Cost         int64
		CreatedAt    time.Time
	}

	type BillingDTO struct {
		TenantID     string  `json:"tenant_id"`
		UsageType    string  `json:"usage_type"`
		Amount       int64   `json:"amount"`
		Cost         float64 `json:"cost"`
		CostYuan     float64 `json:"cost_yuan"`
		CreatedAt    string  `json:"created_at"`
	}

	record := BillingRecord{
		TenantID:  "tenant-001",
		UsageType: "token",
		Amount:    1000,
		Cost:      2000, // 0.002元 = 2000微元
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		dto := BillingDTO{
			TenantID:  record.TenantID,
			UsageType: record.UsageType,
			Amount:    record.Amount,
			Cost:      float64(record.Cost) / 1000000, // 转换为元
			CostYuan:  float64(record.Cost) / 1000000,
			CreatedAt: record.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		_ = dto
	}
}
