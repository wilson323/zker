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
	"sync"
	"testing"
	"time"
)

/**
 * 性能基准测试 - 权限管理和智能路由
 *
 * 设计原则：
 * - SOLID: 单一职责，每个基准测试只测试一个功能点
 * - DRY: 复用测试数据，避免重复初始化
 * - KISS: 保持简单明了
 *
 * @author 研发B
 * @date 2025-01-04
 */

// ==================== 数据结构 ====================

/**
 * Mock 权限检查器
 */
type MockPermissionChecker struct {
	cache map[string]bool
	mu    sync.RWMutex
}

/**
 * Mock 路由匹配器
 */
type MockRouter struct {
	rules []RoutingRule
	mu    sync.RWMutex
}

/**
 * 路由规则
 */
type RoutingRule struct {
	RuleID      string
	Intent      string
	Priority    int
	TargetBotID string
	Conditions  map[string]string
}

// ==================== 权限检查基准测试 ====================

/**
 * BenchmarkPermissionCheck 单次权限检查性能
 *
 * 预期性能目标：
 * - ops/s: > 100,000 次/秒
 * - ns/op: < 10,000 ns/op
 */
func BenchmarkPermissionCheck(b *testing.B) {
	checker := &MockPermissionChecker{
		cache: map[string]bool{
			"user-123:tenant:read":   true,
			"user-123:tenant:write":  false,
			"user-123:bot:read":      true,
			"user-123:bot:write":     false,
			"user-456:tenant:read":   true,
			"user-456:tenant:write":  true,
			"user-456:bot:read":      true,
			"user-456:bot:write":     true,
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟权限检查
		userID := "user-123"
		resource := "tenant"
		action := "read"
		key := userID + ":" + resource + ":" + action

		checker.mu.RLock()
		_ = checker.cache[key]
		checker.mu.RUnlock()
	}
}

/**
 * BenchmarkPermissionCheckConcurrent 并发权限检查性能
 *
 * 预期性能目标：
 * - 并发数: 100
 * - ops/s: > 500,000 次/秒
 */
func BenchmarkPermissionCheckConcurrent(b *testing.B) {
	checker := &MockPermissionChecker{
		cache: map[string]bool{
			"user-123:tenant:read":   true,
			"user-123:tenant:write":  false,
			"user-123:bot:read":      true,
			"user-123:bot:write":     false,
			"user-456:tenant:read":   true,
			"user-456:tenant:write":  true,
			"user-456:bot:read":      true,
			"user-456:bot:write":     true,
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			userID := "user-" + strconv.Itoa(i%2+1)*100+23
			resource := "tenant"
			action := "read"
			key := userID + ":" + resource + ":" + action

			checker.mu.RLock()
			_ = checker.cache[key]
			checker.mu.RUnlock()

			i++
		}
	})
}

/**
 * BenchmarkRoleRetrieval 角色查询性能
 *
 * 预期性能目标：
 * - ops/s: > 50,000 次/秒
 * - ns/op: < 20,000 ns/op
 */
func BenchmarkRoleRetrieval(b *testing.B) {
	// Mock 角色数据
	roles := map[string][]string{
		"user-123": {"role-admin", "role-user"},
		"user-456": {"role-user"},
		"user-789": {"role-guest"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		userID := "user-" + strconv.Itoa(i%3+1)*100+23
		_ = roles[userID]
	}
}

// ==================== 智能路由基准测试 ====================

/**
 * BenchmarkIntentMatching 意图匹配性能
 *
 * 预期性能目标：
 * - ops/s: > 10,000 次/秒
 * - ns/op: < 100,000 ns/op
 */
func BenchmarkIntentMatching(b *testing.B) {
	router := &MockRouter{
		rules: []RoutingRule{
			{
				RuleID:      "rule-1",
				Intent:      "customer_service",
				Priority:    1,
				TargetBotID: "bot-cs-001",
				Conditions: map[string]string{
					"keyword": "咨询,问题,帮助",
				},
			},
			{
				RuleID:      "rule-2",
				Intent:      "sales",
				Priority:    2,
				TargetBotID: "bot-sales-001",
				Conditions: map[string]string{
					"keyword": "购买,价格,优惠",
				},
			},
			{
				RuleID:      "rule-3",
				Intent:      "technical",
				Priority:    3,
				TargetBotID: "bot-tech-001",
				Conditions: map[string]string{
					"keyword": "故障,错误,异常",
				},
			},
		},
	}

	testInputs := []string{
		"你好，我想咨询一下产品",
		"这个价格怎么样",
		"系统出现错误了",
		"帮我解决一下问题",
		"我想购买这个产品",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		input := testInputs[i%len(testInputs)]

		// 简单的关键词匹配逻辑
		router.mu.RLock()
		for _, rule := range router.rules {
			_ = rule
			_ = input
			// 模拟匹配逻辑
		}
		router.mu.RUnlock()
	}
}

/**
 * BenchmarkRoutingWithScoring 评分路由性能
 *
 * 预期性能目标：
 * - ops/s: > 5,000 次/秒
 * - ns/op: < 200,000 ns/op
 */
func BenchmarkRoutingWithScoring(b *testing.B) {
	router := &MockRouter{
		rules: []RoutingRule{
			{RuleID: "rule-1", Intent: "customer_service", Priority: 1, TargetBotID: "bot-cs-001"},
			{RuleID: "rule-2", Intent: "sales", Priority: 2, TargetBotID: "bot-sales-001"},
			{RuleID: "rule-3", Intent: "technical", Priority: 3, TargetBotID: "bot-tech-001"},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟评分路由
		router.mu.RLock()
		var bestRule *RoutingRule
		bestScore := 0.0

		for _, rule := range router.rules {
			// 模拟评分逻辑
			score := float64(rule.Priority) * 0.5
			if score > bestScore {
				bestScore = score
				bestRule = &rule
			}
		}
		router.mu.RUnlock()

		_ = bestRule
	}
}

/**
 * BenchmarkRoutingConcurrent 并发路由性能
 *
 * 预期性能目标：
 * - 并发数: 50
 * - ops/s: > 200,000 次/秒
 */
func BenchmarkRoutingConcurrent(b *testing.B) {
	router := &MockRouter{
		rules: []RoutingRule{
			{RuleID: "rule-1", Intent: "customer_service", Priority: 1, TargetBotID: "bot-cs-001"},
			{RuleID: "rule-2", Intent: "sales", Priority: 2, TargetBotID: "bot-sales-001"},
			{RuleID: "rule-3", Intent: "technical", Priority: 3, TargetBotID: "bot-tech-001"},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 模拟并发路由请求
			router.mu.RLock()
			if len(router.rules) > 0 {
				_ = router.rules[0]
			}
			router.mu.RUnlock()
		}
	})
}

// ==================== 缓存性能基准测试 ====================

/**
 * BenchmarkCacheHit 缓存命中性能
 *
 * 预期性能目标：
 * - ops/s: > 1,000,000 次/秒
 * - ns/op: < 1,000 ns/op
 */
func BenchmarkCacheHit(b *testing.B) {
	cache := make(map[string]string)
	cache["key-1"] = "value-1"
	cache["key-2"] = "value-2"
	cache["key-3"] = "value-3"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := "key-" + strconv.Itoa(i%3+1)
		_ = cache[key]
	}
}

/**
 * BenchmarkCacheMiss 缓存未命中性能
 *
 * 预期性能目标：
 * - ops/s: > 500,000 次/秒
 * - ns/op: < 2,000 ns/op
 */
func BenchmarkCacheMiss(b *testing.B) {
	cache := make(map[string]string)
	cache["key-1"] = "value-1"
	cache["key-2"] = "value-2"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := "key-not-exist-" + strconv.Itoa(i)
		_ = cache[key]
	}
}

/**
 * BenchmarkCacheWithLock 带锁的缓存性能
 *
 * 预期性能目标：
 * - ops/s: > 300,000 次/秒
 * - ns/op: < 5,000 ns/op
 */
func BenchmarkCacheWithLock(b *testing.B) {
	type SafeCache struct {
		data map[string]string
		mu   sync.RWMutex
	}

	cache := &SafeCache{
		data: map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
			"key-3": "value-3",
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := "key-" + strconv.Itoa(i%3+1)
		cache.mu.RLock()
		_ = cache.data[key]
		cache.mu.RUnlock()
	}
}

// ==================== 上下文传递基准测试 ====================

/**
 * BenchmarkContextValueRetrieval 上下文值获取性能
 *
 * 预期性能目标：
 * - ops/s: > 5,000,000 次/秒
 * - ns/op: < 300 ns/op
 */
func BenchmarkContextValueRetrieval(b *testing.B) {
	type contextKey string
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey("tenant_id"), "tenant-123")
	ctx = context.WithValue(ctx, contextKey("user_id"), "user-456")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ctx.Value(contextKey("tenant_id"))
		_ = ctx.Value(contextKey("user_id"))
	}
}

/**
 * BenchmarkContextCreation 上下文创建性能
 *
 * 预期性能目标：
 * - ops/s: > 1,000,000 次/秒
 * - ns/op: < 1,000 ns/op
 */
func BenchmarkContextCreation(b *testing.B) {
	type contextKey string

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKey("tenant_id"), "tenant-123")
		ctx = context.WithValue(ctx, contextKey("user_id"), "user-456")
		_ = ctx
	}
}

// ==================== JSON 序列化基准测试 ====================

/**
 * BenchmarkJSONMarshal JSON序列化性能
 *
 * 预期性能目标：
 * - ops/s: > 100,000 次/秒
 * - ns/op: < 10,000 ns/op
 */
func BenchmarkJSONMarshal(b *testing.B) {
	type TestData struct {
		TenantID   string            `json:"tenant_id"`
		UserID     string            `json:"user_id"`
		Action     string            `json:"action"`
		Resource   string            `json:"resource"`
		Metadata   map[string]string `json:"metadata"`
		Timestamp  time.Time         `json:"timestamp"`
	}

	data := TestData{
		TenantID:  "tenant-123",
		UserID:    "user-456",
		Action:    "read",
		Resource:  "bot",
		Metadata:  map[string]string{"key1": "value1", "key2": "value2"},
		Timestamp: time.Now(),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = data.MarshalJSON()
	}
}

/**
 * BenchmarkJSONUnmarshal JSON反序列化性能
 *
 * 预期性能目标：
 * - ops/s: > 100,000 次/秒
 * - ns/op: < 10,000 ns/op
 */
func BenchmarkJSONUnmarshal(b *testing.B) {
	type TestData struct {
		TenantID   string            `json:"tenant_id"`
		UserID     string            `json:"user_id"`
		Action     string            `json:"action"`
		Resource   string            `json:"resource"`
		Metadata   map[string]string `json:"metadata"`
		Timestamp  time.Time         `json:"timestamp"`
	}

	jsonData := `{"tenant_id":"tenant-123","user_id":"user-456","action":"read","resource":"bot","metadata":{"key1":"value1","key2":"value2"},"timestamp":"2025-01-04T10:00:00Z"}`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var data TestData
		_ = data.UnmarshalJSON([]byte(jsonData))
	}
}
