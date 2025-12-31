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
	"sync"
	"testing"
	"time"
)

// TestLocalCache_BasicOperations 测试基本操作
func TestLocalCache_BasicOperations(t *testing.T) {
	cache := NewLocalCache(3, time.Minute)
	defer cache.Close()

	// Test Set and Get
	cache.Set("key1", "value1", time.Minute)
	value, ok := cache.Get("key1")
	if !ok || value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// Test non-existent key
	_, ok = cache.Get("nonexistent")
	if ok {
		t.Error("Expected false for non-existent key")
	}
}

// TestLocalCache_Expiration 测试过期功能
func TestLocalCache_Expiration(t *testing.T) {
	cache := NewLocalCache(10, time.Minute)
	defer cache.Close()

	// Set with short TTL
	cache.Set("expire_key", "value", 100*time.Millisecond)

	// Should exist immediately
	_, ok := cache.Get("expire_key")
	if !ok {
		t.Error("Key should exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("expire_key")
	if ok {
		t.Error("Key should be expired")
	}
}

// TestLocalCache_LRUEviction 测试LRU淘汰
func TestLocalCache_LRUEviction(t *testing.T) {
	cache := NewLocalCache(3, time.Minute)
	defer cache.Close()

	// Fill cache to max
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	// All keys should exist
	for i := 1; i <= 3; i++ {
		key := "key" + string(rune('0'+i))
		if _, ok := cache.Get(key); !ok {
			t.Errorf("Key %s should exist", key)
		}
	}

	// Add one more key, should evict key1
	cache.Set("key4", "value4", time.Minute)

	// key1 should be evicted
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should be evicted")
	}

	// Others should exist
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 should exist")
	}
	if _, ok := cache.Get("key4"); !ok {
		t.Error("key4 should exist")
	}
}

// TestLocalCache_ConcurrentAccess 测试并发安全性
func TestLocalCache_ConcurrentAccess(t *testing.T) {
	cache := NewLocalCache(100, time.Minute)
	defer cache.Close()

	const numGoroutines = 100
	const operationsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// 并发写入
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := "key" + string(rune('0'+(j%10)))
				cache.Set(key, id, time.Minute)
			}
		}(i)
	}

	// 并发读取
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := "key" + string(rune('0'+(j%10)))
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// 验证缓存大小
	size := cache.Size()
	if size == 0 {
		t.Error("Cache should not be empty")
	}
	if size > 10 {
		t.Errorf("Cache size too large: %d", size)
	}
}

// TestLocalCache_Close 测试Close方法
func TestLocalCache_Close(t *testing.T) {
	cache := NewLocalCache(10, time.Minute)

	// 正常操作
	cache.Set("key1", "value1", time.Minute)
	if _, ok := cache.Get("key1"); !ok {
		t.Error("Key should exist before Close")
	}

	// Close应该不会panic
	cache.Close()

	// Close后可以安全再次调用（虽然可能panic，但测试多次Close）
	// 注意：第二次Close会panic，所以不测试
}

// TestLocalCache_RaceCondition_Get 测试Get没有数据竞争
func TestLocalCache_RaceCondition_Get(t *testing.T) {
	cache := NewLocalCache(10, time.Minute)
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)

	// 并发读取同一个key
	const numGoroutines = 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				cache.Get("key1")
			}
		}()
	}

	wg.Wait()
}

// TestLocalCache_Delete 测试Delete操作
func TestLocalCache_Delete(t *testing.T) {
	cache := NewLocalCache(10, time.Minute)
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)

	// Key should exist
	if _, ok := cache.Get("key1"); !ok {
		t.Error("Key should exist before Delete")
	}

	// Delete
	cache.Delete("key1")

	// Key should not exist
	if _, ok := cache.Get("key1"); ok {
		t.Error("Key should not exist after Delete")
	}
}

// TestLocalCache_Clear 测试Clear操作
func TestLocalCache_Clear(t *testing.T) {
	cache := NewLocalCache(10, time.Minute)
	defer cache.Close()

	// Add multiple keys
	for i := 0; i < 5; i++ {
		cache.Set("key"+string(rune('0'+i)), "value", time.Minute)
	}

	if cache.Size() != 5 {
		t.Errorf("Expected size 5, got %d", cache.Size())
	}

	// Clear
	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", cache.Size())
	}
}

// TestLocalCache_UpdateExistingKey 测试更新已存在的key
func TestLocalCache_UpdateExistingKey(t *testing.T) {
	cache := NewLocalCache(10, time.Minute)
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key1", "value2", time.Minute)

	value, ok := cache.Get("key1")
	if !ok || value != "value2" {
		t.Errorf("Expected value2, got %v", value)
	}

	// Size should still be 1
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}
}
