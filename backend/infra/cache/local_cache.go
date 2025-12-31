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
	"time"
)

// LocalCache 本地内存缓存（LRU + TTL）
// ✅ 安全修复：
// 1. 修复数据竞争：Get操作在读锁下不调用updateLRU
// 2. 修复Goroutine泄漏：添加停止机制和Close方法
type LocalCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem

	// 配置
	maxSize int
	ttl     time.Duration

	// LRU相关
	lruList  []string
	lruIndex map[string]int

	// ✅ Goroutine生命周期管理
	stopCh chan struct{}
	done   sync.WaitGroup
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewLocalCache 创建本地缓存
func NewLocalCache(maxSize int, ttl time.Duration) *LocalCache {
	cache := &LocalCache{
		items:    make(map[string]*cacheItem),
		maxSize:  maxSize,
		ttl:      ttl,
		lruList:  make([]string, 0, maxSize),
		lruIndex: make(map[string]int),
		stopCh:   make(chan struct{}),
	}

	// ✅ 修复Goroutine泄漏：使用WaitGroup跟踪
	cache.done.Add(1)
	go func() {
		defer cache.done.Done()
		cache.cleanupExpired()
	}()

	return cache
}

// Get 获取缓存值
// ✅ 修复数据竞争：不在读锁下调写锁方法
func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// 检查是否过期
	if time.Now().After(item.expiration) {
		return nil, false
	}

	// ✅ 不更新LRU，避免在读锁下调写锁方法
	// LRU淘汰策略可以仅在Set时触发

	return item.value, true
}

// Set 设置缓存值
func (c *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果已存在，更新值并返回
	if _, ok := c.items[key]; ok {
		c.items[key] = &cacheItem{
			value:      value,
			expiration: time.Now().Add(ttl),
		}
		c.updateLRUUnsafe(key)
		return
	}

	// 检查是否需要淘汰
	if len(c.items) >= c.maxSize {
		c.evictLRUUnsafe()
	}

	// 添加新项
	c.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}

	// 更新LRU
	c.lruList = append(c.lruList, key)
	c.lruIndex[key] = len(c.lruList) - 1
}

// Delete 删除缓存值
func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	c.removeLRUUnsafe(key)
}

// Clear 清空缓存
func (c *LocalCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
	c.lruList = make([]string, 0, c.maxSize)
	c.lruIndex = make(map[string]int)
}

// Size 获取缓存大小
func (c *LocalCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Close 关闭缓存，停止后台goroutine
// ✅ 修复Goroutine泄漏：提供优雅停止机制
func (c *LocalCache) Close() {
	close(c.stopCh)
	c.done.Wait()
}

// updateLRUUnsafe 更新LRU位置（将访问的key移到末尾）
// ✅ 必须在持有写锁时调用
func (c *LocalCache) updateLRUUnsafe(key string) {
	idx, ok := c.lruIndex[key]
	if !ok {
		return
	}

	// 移动到末尾
	c.lruList = append(c.lruList[:idx], c.lruList[idx+1:]...)
	c.lruList = append(c.lruList, key)
	c.lruIndex[key] = len(c.lruList) - 1
}

// evictLRUUnsafe 淘汰最久未使用的项
// ✅ 必须在持有写锁时调用
func (c *LocalCache) evictLRUUnsafe() {
	if len(c.lruList) == 0 {
		return
	}

	// 淘汰第一个（最久未使用）
	key := c.lruList[0]
	delete(c.items, key)
	c.removeLRUUnsafe(key)
}

// removeLRUUnsafe 从LRU链表中移除
// ✅ 必须在持有写锁时调用
func (c *LocalCache) removeLRUUnsafe(key string) {
	idx, ok := c.lruIndex[key]
	if !ok {
		return
	}

	c.lruList = append(c.lruList[:idx], c.lruList[idx+1:]...)
	delete(c.lruIndex, key)
}

// cleanupExpired 定期清理过期项
// ✅ 修复Goroutine泄漏：监听停止信号，支持优雅退出
func (c *LocalCache) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()

			for key, item := range c.items {
				if now.After(item.expiration) {
					delete(c.items, key)
					c.removeLRUUnsafe(key)
				}
			}

			c.mu.Unlock()
		case <-c.stopCh:
			// ✅ 收到停止信号，退出goroutine
			return
		}
	}
}
