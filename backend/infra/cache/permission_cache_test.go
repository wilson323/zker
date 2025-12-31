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

// MockCmdable 模拟Redis客户端
type MockCmdable struct {
	getFunc    func(ctx context.Context, key string) *MockStringCmd
	setFunc    func(ctx context.Context, key string, value interface{}, expiration time.Duration) *MockStatusCmd
	delFunc    func(ctx context.Context, keys ...string) *MockIntCmd
	dataStore  map[string]string
	callCount  map[string]int
}

type MockStringCmd struct {
	value string
	err   error
}

func (m *MockStringCmd) Result() (string, error) {
	return m.value, m.err
}

func (m *MockStringCmd) Val() string {
	return m.value
}

type MockStatusCmd struct {
	err error
}

func (m *MockStatusCmd) Err() error {
	return m.err
}

type MockIntCmd struct {
	err error
}

func (m *MockIntCmd) Err() error {
	return m.err
}

func NewMockCmdable() *MockCmdable {
	return &MockCmdable{
		dataStore: make(map[string]string),
		callCount: make(map[string]int),
	}
}

func (m *MockCmdable) Get(ctx context.Context, key string) StringCmd {
	m.callCount["get"]++
	val, ok := m.dataStore[key]
	if !ok {
		return &MockStringCmd{err: fmt.Errorf("not found")}
	}
	return &MockStringCmd{value: val}
}

func (m *MockCmdable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) StatusCmd {
	m.callCount["set"]++
	m.dataStore[key] = value.(string)
	return &MockStatusCmd{}
}

func (m *MockCmdable) Del(ctx context.Context, keys ...string) IntCmd {
	m.callCount["del"]++
	for _, key := range keys {
		delete(m.dataStore, key)
	}
	return &MockIntCmd{}
}

func (m *MockCmdable) GetCallCount(operation string) int {
	return m.callCount[operation]
}

// TestPermissionCache_BasicOperations 测试基本缓存操作
func TestPermissionCache_BasicOperations(t *testing.T) {
	mockRedis := NewMockCmdable()
	permCache := NewPermissionCache(mockRedis)
	ctx := context.Background()

	// 第一次查询 - 缓存未命中，应调用数据库查询
	cacheKey := "tenant1:user1:bots:read:bot123"
	dbQueryCount := 0

	result, err := permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		dbQueryCount++
		return true, nil // 模拟数据库查询返回true
	})

	require.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, 1, dbQueryCount, "第一次应该查询数据库")

	// 第二次查询 - 缓存命中，不应调用数据库查询
	result, err = permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		dbQueryCount++
		return true, nil
	})

	require.NoError(t, err)
	assert.True(t, result)
	assert.Equal(t, 1, dbQueryCount, "第二次应该从缓存获取，不查询数据库")
}

// TestPermissionCache_CacheInvalidation 测试缓存失效
func TestPermissionCache_CacheInvalidation(t *testing.T) {
	mockRedis := NewMockCmdable()
	permCache := NewPermissionCache(mockRedis)
	ctx := context.Background()

	cacheKey := "tenant1:user1:bots:read:bot123"
	dbQueryCount := 0

	// 第一次查询 - 缓存未命中
	permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		dbQueryCount++
		return true, nil
	})
	assert.Equal(t, 1, dbQueryCount)

	// 第二次查询 - 缓存命中
	permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		dbQueryCount++
		return true, nil
	})
	assert.Equal(t, 1, dbQueryCount)

	// 失效缓存
	err := permCache.InvalidateDataPermission(ctx, cacheKey)
	require.NoError(t, err)

	// 第三次查询 - 缓存已失效，重新查询数据库
	permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		dbQueryCount++
		return true, nil
	})
	assert.Equal(t, 2, dbQueryCount, "缓存失效后应该重新查询数据库")
}

// TestPermissionCache_CacheStats 测试缓存统计
func TestPermissionCache_CacheStats(t *testing.T) {
	mockRedis := NewMockCmdable()
	permCache := NewPermissionCache(mockRedis)
	ctx := context.Background()

	// 执行10次查询（相同key）
	for i := 0; i < 10; i++ {
		cacheKey := fmt.Sprintf("tenant1:user1:bots:read:bot%d", i%3) // 3个不同的key
		permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
			return true, nil
		})
	}

	stats := permCache.GetCacheStats()
	assert.Equal(t, int64(3), stats.MissCount, "应该有3次缓存未命中")
	assert.Equal(t, int64(7), stats.HitCount, "应该有7次缓存命中")

	hitRate := permCache.GetCacheStats().HitCount / float64(permCache.GetCacheStats().HitCount+permCache.GetCacheStats().MissCount)
	assert.Greater(t, hitRate, 0.5, "缓存命中率应该>50%")
}

// TestLocalCache_LRUEviction 测试LRU淘汰
func TestLocalCache_LRUEviction(t *testing.T) {
	localCache := NewLocalCache(3, time.Minute) // 最多3条

	// 添加3个key
	localCache.Set("key1", "value1", time.Minute)
	localCache.Set("key2", "value2", time.Minute)
	localCache.Set("key3", "value3", time.Minute)

	// 验证3个key都存在
	val, ok := localCache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// 添加第4个key（应该淘汰key1）
	localCache.Set("key4", "value4", time.Minute)

	// key1应该被淘汰
	_, ok = localCache.Get("key1")
	assert.False(t, ok, "key1应该被LRU淘汰")

	// key2应该还存在
	val, ok = localCache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", val)
}

// TestLocalCache_TTLExpiration 测试TTL过期
func TestLocalCache_TTLExpiration(t *testing.T) {
	localCache := NewLocalCache(100, 100*time.Millisecond) // 100ms过期

	// 添加key
	localCache.Set("key1", "value1", 100*time.Millisecond)

	// 立即读取 - 应该存在
	val, ok := localCache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// 等待过期
	time.Sleep(150 * time.Millisecond)

	// 再次读取 - 应该已过期
	_, ok = localCache.Get("key1")
	assert.False(t, ok, "key1应该已过期")
}

// BenchmarkPermissionCache_WithCache 基准测试：使用缓存
func BenchmarkPermissionCache_WithCache(b *testing.B) {
	mockRedis := NewMockCmdable()
	permCache := NewPermissionCache(mockRedis)
	ctx := context.Background()

	// 预热缓存
	cacheKey := "tenant1:user1:bots:read:bot123"
	permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
		return true, nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		permCache.CheckDataPermission(ctx, cacheKey, func() (bool, error) {
			return true, nil
		})
	}
}

// BenchmarkPermissionCache_WithoutCache 基准测试：不使用缓存
func BenchmarkPermissionCache_WithoutCache(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 模拟直接查询数据库
		time.Sleep(20 * time.Millisecond) // 假设数据库查询需要20ms
	}
}
