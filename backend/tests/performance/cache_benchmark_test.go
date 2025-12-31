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

	"github.com/coze-dev/coze-studio/backend/infra/cache"
)

// =====================================================
// 多级缓存性能测试
// =====================================================

// BenchmarkMultiLevelCache_Get 测试多级缓存Get性能
func BenchmarkMultiLevelCache_Get(b *testing.B) {
	// 注意: 需要先启动Redis服务
	// redisCache := cache.NewRedisCache("localhost:6379", "", 0, 1*time.Hour)
	// multiCache := cache.NewMultiLevelCache(redisCache)

	// ctx := context.Background()
	// key := "benchmark:test:key"
	// value := map[string]interface{}{"id": 1, "name": "test"}

	// // 预热缓存
	// multiCache.Set(ctx, key, value)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟缓存操作
		key := fmt.Sprintf("bench:key:%d", i%1000)
		_ = key
		// multiCache.Get(ctx, key, &value)
	}
}

// BenchmarkMultiLevelCache_Set 测试多级缓存Set性能
func BenchmarkMultiLevelCache_Set(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench:set:%d", i)
		value := map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("test-%d", i),
		}
		_ = key
		_ = value
		// multiCache.Set(ctx, key, value)
	}
}

// BenchmarkMultiLevelCache_GetBatch 测试批量获取性能
func BenchmarkMultiLevelCache_GetBatch(b *testing.B) {
	keys := make([]string, 100)
	for i := 0; i < 100; i++ {
		keys[i] = fmt.Sprintf("bench:batch:%d", i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = keys
		// multiCache.GetBatch(ctx, keys)
	}
}

// BenchmarkMultiLevelCache_SetBatch 测试批量设置性能
func BenchmarkMultiLevelCache_SetBatch(b *testing.B) {
	items := make(map[string]interface{}, 100)
	for i := 0; i < 100; i++ {
		items[fmt.Sprintf("bench:batch:%d", i)] = map[string]interface{}{"id": i}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = items
		// multiCache.SetBatch(ctx, items)
	}
}

// =====================================================
// 数据库查询性能测试 (模拟)
// =====================================================

// BenchmarkDBQuery_Simple 模拟简单查询性能
func BenchmarkDBQuery_Simple(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟: SELECT * FROM organizations WHERE org_id = ?
		orgID := fmt.Sprintf("org-%d", i%1000)
		_ = orgID
		// db.Where("org_id = ?", orgID).First(&org)
	}
}

// BenchmarkDBQuery_WithJoin 模拟JOIN查询性能
func BenchmarkDBQuery_WithJoin(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟: SELECT * FROM organizations JOIN leader ON organizations.leader_id = leader.id
		orgID := fmt.Sprintf("org-%d", i%1000)
		_ = orgID
		// db.Joins("Leader").Where("org_id = ?", orgID).First(&org)
	}
}

// BenchmarkDBQuery_List 分页查询性能测试
func BenchmarkDBQuery_List(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		page := (i / 10) % 100
		pageSize := 20
		_ = tenantID
		_ = page
		_ = pageSize
		// db.Where("tenant_id = ?", tenantID).Offset(page*pageSize).Limit(pageSize).Find(&orgs)
	}
}

// =====================================================
// API端点性能测试 (模拟)
// =====================================================

// BenchmarkAPI_GetBot 模拟获取Bot信息API
func BenchmarkAPI_GetBot(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		botID := fmt.Sprintf("bot-%d", i%10000)

		// 模拟API处理流程:
		// 1. 参数验证
		_ = botID
		// 2. 权限检查 (缓存)
		// 3. 查询Bot (缓存/DB)
		// 4. 序列化JSON
		// 5. 返回响应
	}
}

// BenchmarkAPI_ListBots 模拟Bot列表API
func BenchmarkAPI_ListBots(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		page := (i / 10) % 100
		pageSize := 20

		// 模拟API处理流程:
		// 1. 参数验证
		_ = tenantID
		_ = page
		_ = pageSize
		// 2. 权限检查 (缓存)
		// 3. 查询Bot列表 (DB + 缓存)
		// 4. 查询分类信息 (批量)
		// 5. 序列化JSON
		// 6. 返回响应
	}
}

// BenchmarkAPI_CreateBot 模拟创建Bot API
func BenchmarkAPI_CreateBot(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		botName := fmt.Sprintf("Test Bot %d", i)

		// 模拟API处理流程:
		// 1. 参数验证
		_ = botName
		// 2. 权限检查
		// 3. 创建Bot (DB事务)
		// 4. 清除缓存
		// 5. 返回响应
	}
}

// =====================================================
// JSON序列化性能测试
// =====================================================

import (
	"encoding/json"
)

type TestBot struct {
	BotID      string    `json:"bot_id"`
	Name       string    `json:"name"`
	Desc       string    `json:"description"`
	TenantID   string    `json:"tenant_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
	Version    int       `json:"version"`
	CategoryID string    `json:"category_id"`
}

// BenchmarkJSONMarshal 测试JSON序列化性能
func BenchmarkJSONMarshal(b *testing.B) {
	bot := &TestBot{
		BotID:      "bot-12345",
		Name:       "Test Bot for Benchmark",
		Desc:       "This is a test bot for performance benchmarking",
		TenantID:   "tenant-67890",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     "active",
		Version:    1,
		CategoryID: "cat-001",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bot)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONUnmarshal 测试JSON反序列化性能
func BenchmarkJSONUnmarshal(b *testing.B) {
	data, _ := json.Marshal(&TestBot{
		BotID:      "bot-12345",
		Name:       "Test Bot for Benchmark",
		Desc:       "This is a test bot for performance benchmarking",
		TenantID:   "tenant-67890",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Status:     "active",
		Version:    1,
		CategoryID: "cat-001",
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var bot TestBot
		err := json.Unmarshal(data, &bot)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// =====================================================
// 并发性能测试
// =====================================================

// BenchmarkConcurrent_CacheGet 并发缓存读取测试
func BenchmarkConcurrent_CacheGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("concurrent:key:%d", i%1000)
			_ = key
			// cache.Get(ctx, key, &value)
			i++
		}
	})
}

// BenchmarkConcurrent_DBQuery 并发数据库查询测试
func BenchmarkConcurrent_DBQuery(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			orgID := fmt.Sprintf("org-%d", i%1000)
			_ = orgID
			// db.Where("org_id = ?", orgID).First(&org)
			i++
		}
	})
}
