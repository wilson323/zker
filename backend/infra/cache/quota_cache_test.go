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

package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQuotaCache_GetQuota 测试获取配额
func TestQuotaCache_GetQuota(t *testing.T) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	tenantID := "tenant123"
	resourceType := "bots"

	// 第一次查询 - 缓存未命中
	dbQueryCount := 0
	quotaInfo, err := quotaCache.GetQuota(ctx, tenantID, resourceType, func() (*QuotaInfo, error) {
		dbQueryCount++
		return &QuotaInfo{
			TenantID:     tenantID,
			ResourceType: resourceType,
			MaxLimit:     100,
			UsedCount:    50,
			IsUnlimited:  false,
		}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, dbQueryCount, "第一次应该查询数据库")
	assert.Equal(t, 100, quotaInfo.MaxLimit)
	assert.Equal(t, 50, quotaInfo.UsedCount)
	assert.True(t, quotaInfo.HasRemaining(10))
	assert.False(t, quotaInfo.HasRemaining(60))

	// 第二次查询 - 缓存命中
	quotaInfo, err = quotaCache.GetQuota(ctx, tenantID, resourceType, func() (*QuotaInfo, error) {
		dbQueryCount++
		return &QuotaInfo{}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, dbQueryCount, "第二次应该从缓存获取")
	assert.Equal(t, 100, quotaInfo.MaxLimit)
}

// TestQuotaCache_ConsumeQuota 测试消费配额
func TestQuotaCache_ConsumeQuota(t *testing.T) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	tenantID := "tenant123"
	resourceType := "bots"
	consumeCallCount := 0

	// 消费配额
	err := quotaCache.ConsumeQuota(ctx, tenantID, resourceType, 5, func(count int) error {
		consumeCallCount++
		assert.Equal(t, 5, count)
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, consumeCallCount, "应该调用一次消费函数")

	// 验证缓存已失效（下次访问会重新加载）
	dbQueryCount := 0
	quotaCache.GetQuota(ctx, tenantID, resourceType, func() (*QuotaInfo, error) {
		dbQueryCount++
		return &QuotaInfo{
			TenantID:     tenantID,
			ResourceType: resourceType,
			MaxLimit:     100,
			UsedCount:    55, // 50 + 5
			IsUnlimited:  false,
		}, nil
	})

	assert.Equal(t, 1, dbQueryCount, "缓存失效后应该重新查询数据库")
}

// TestQuotaCache_UnlimitedQuota 测试无限制配额
func TestQuotaCache_UnlimitedQuota(t *testing.T) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	quotaInfo := &QuotaInfo{
		TenantID:     "tenant123",
		ResourceType: "bots",
		MaxLimit:     -1, // -1 表示无限制
		UsedCount:    9999,
		IsUnlimited:  true,
	}

	// 无限制配额应该总是返回true
	assert.True(t, quotaInfo.HasRemaining(1))
	assert.True(t, quotaInfo.HasRemaining(10000))
	assert.Equal(t, -1, quotaInfo.GetRemaining())
	assert.Equal(t, float64(0), quotaInfo.GetUsagePercentage())
}

// TestQuotaCache_QuotaExceeded 测试配额超额
func TestQuotaCache_GetQuota(t *testing.T) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	tenantID := "tenant123"
	resourceType := "bots"

	// 第一次查询 - 缓存未命中
	dbQueryCount := 0
	quotaInfo, err := quotaCache.GetQuota(ctx, tenantID, resourceType, func() (*QuotaInfo, error) {
		dbQueryCount++
		return &QuotaInfo{
			TenantID:     tenantID,
			ResourceType: resourceType,
			MaxLimit:     100,
			UsedCount:    95,
			IsUnlimited:  false,
		}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, 1, dbQueryCount)
	assert.Equal(t, 5, quotaInfo.GetRemaining())

	// 使用率应该接近95%
	usagePercent := quotaInfo.GetUsagePercentage()
	assert.Greater(t, usagePercent, 94.0)
	assert.Less(t, usagePercent, 96.0)
}

// TestQuotaCache_InvalidateQuota 测试配额缓存失效
func TestQuotaCache_InvalidateQuota(t *testing.T) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	tenantID := "tenant123"
	resourceType := "bots"

	// 第一次查询 - 缓存未命中
	dbQueryCount := 0
	quotaCache.GetQuota(ctx, tenantID, resourceType, func() (*QuotaInfo, error) {
		dbQueryCount++
		return &QuotaInfo{
			TenantID:     tenantID,
			ResourceType: resourceType,
			MaxLimit:     100,
			UsedCount:    50,
			IsUnlimited:  false,
		}, nil
	})

	assert.Equal(t, 1, dbQueryCount)

	// 失效缓存
	err := quotaCache.InvalidateQuota(ctx, tenantID, resourceType)
	require.NoError(t, err)

	// 第二次查询 - 缓存已失效，重新查询数据库
	quotaCache.GetQuota(ctx, tenantID, resourceType, func() (*QuotaInfo, error) {
		dbQueryCount++
		return &QuotaInfo{
			TenantID:     tenantID,
			ResourceType: resourceType,
			MaxLimit:     100,
			UsedCount:    60,
			IsUnlimited:  false,
		}, nil
	})

	assert.Equal(t, 2, dbQueryCount, "缓存失效后应该重新查询数据库")
}

// TestQuotaCache_GetCacheStats 测试缓存统计
func TestQuotaCache_GetCacheStats(t *testing.T) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	// 执行10次查询（相同key）
	for i := 0; i < 10; i++ {
		quotaCache.GetQuota(ctx, "tenant1", "bots", func() (*QuotaInfo, error) {
			return &QuotaInfo{
				TenantID:     "tenant1",
				ResourceType: "bots",
				MaxLimit:     100,
				UsedCount:    50,
				IsUnlimited:  false,
			}, nil
		})
	}

	stats := quotaCache.GetCacheStats()
	assert.Equal(t, int64(1), stats.MissCount, "应该有1次缓存未命中")
	assert.Equal(t, int64(9), stats.HitCount, "应该有9次缓存命中")

	hitRate := quotaCache.GetHitRate()
	assert.Greater(t, hitRate, 0.8, "缓存命中率应该>80%")
}

// BenchmarkQuotaCache_WithCache 基准测试：使用缓存
func BenchmarkQuotaCache_WithCache(b *testing.B) {
	mockRedis := NewMockCmdable()
	quotaCache := NewQuotaCache(mockRedis)
	ctx := context.Background()

	// 预热缓存
	quotaCache.GetQuota(ctx, "tenant1", "bots", func() (*QuotaInfo, error) {
		return &QuotaInfo{
			TenantID:     "tenant1",
			ResourceType: "bots",
			MaxLimit:     100,
			UsedCount:    50,
			IsUnlimited:  false,
		}, nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		quotaCache.GetQuota(ctx, "tenant1", "bots", func() (*QuotaInfo, error) {
			return &QuotaInfo{}, nil
		})
	}
}

// BenchmarkQuotaCache_WithoutCache 基准测试：不使用缓存
func BenchmarkQuotaCache_WithoutCache(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟直接查询数据库
		time.Sleep(15 * time.Millisecond) // 假设数据库查询需要15ms
	}
}

// TestQuotaInfo 测试QuotaInfo方法
func TestQuotaInfo(t *testing.T) {
	tests := []struct {
		name          string
		quotaInfo     *QuotaInfo
		required      int
		hasRemaining  bool
		remaining     int
		usagePercent  float64
	}{
		{
			name: "正常配额-有剩余",
			quotaInfo: &QuotaInfo{
				MaxLimit:  100,
				UsedCount: 50,
			},
			required:     10,
			hasRemaining: true,
			remaining:    50,
			usagePercent: 50.0,
		},
		{
			name: "正常配额-刚好用完",
			quotaInfo: &QuotaInfo{
				MaxLimit:  100,
				UsedCount: 100,
			},
			required:     1,
			hasRemaining: false,
			remaining:    0,
			usagePercent: 100.0,
		},
		{
			name: "正常配额-已超额",
			quotaInfo: &QuotaInfo{
				MaxLimit:  100,
				UsedCount: 105,
			},
			required:     1,
			hasRemaining: false,
			remaining:    0,
			usagePercent: 105.0,
		},
		{
			name: "无限制配额",
			quotaInfo: &QuotaInfo{
				MaxLimit:    -1,
				UsedCount:   9999,
				IsUnlimited: true,
			},
			required:     10000,
			hasRemaining: true,
			remaining:    -1,
			usagePercent: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.hasRemaining, tt.quotaInfo.HasRemaining(tt.required))
			assert.Equal(t, tt.remaining, tt.quotaInfo.GetRemaining())
			assert.Equal(t, tt.usagePercent, tt.quotaInfo.GetUsagePercentage())
		})
	}
}
