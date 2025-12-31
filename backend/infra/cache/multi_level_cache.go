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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/patrickmn/go-cache"
)

// MultiLevelCache 多级缓存实现
// L1: 本地内存缓存 (快速,容量小)
// L2: Redis缓存 (较慢,容量大)
type MultiLevelCache struct {
	l1        *cache.Cache          // 本地缓存
	l2        *RedisCache           // Redis缓存
	l1TTL     time.Duration         // L1过期时间
	l2TTL     time.Duration         // L2过期时间
	mu        sync.RWMutex          // 保护stats的互斥锁
	stats     *MultiLevelCacheStats  // 统计信息
}

// MultiLevelCacheStats 多级缓存统计
type MultiLevelCacheStats struct {
	L1Hits   int64
	L1Misses int64
	L2Hits   int64
	L2Misses int64
}

// NewMultiLevelCache 创建多级缓存
func NewMultiLevelCache(redis *RedisCache) *MultiLevelCache {
	return &MultiLevelCache{
		l1:    cache.New(5*time.Minute, 10*time.Minute), // L1: 5分钟过期
		l2:    redis,
		l1TTL: 5 * time.Minute,
		l2TTL: 1 * time.Hour, // L2: 1小时过期
		stats: &MultiLevelCacheStats{},
	}
}

// Get 获取缓存 (先查L1,再查L2)
func (m *MultiLevelCache) Get(ctx context.Context, key string, dest interface{}) error {
	// L1: 本地缓存
	if x, found := m.l1.Get(key); found {
		m.mu.Lock()
		m.stats.L1Hits++
		m.mu.Unlock()

		// 类型断言
		data, ok := x.([]byte)
		if !ok {
			return fmt.Errorf("invalid cache data type")
		}

		if err := json.Unmarshal(data, dest); err != nil {
			logs.Errorf("Failed to unmarshal L1 cache: key=%s, error=%v", key, err)
			return err
		}

		logs.Debugf("L1 cache hit: %s", key)
		return nil
	}

	// L1未命中
	m.mu.Lock()
	m.stats.L1Misses++
	m.mu.Unlock()

	// L2: Redis缓存
	if err := m.l2.Get(ctx, key, dest); err == nil {
		m.mu.Lock()
		m.stats.L2Hits++
		m.mu.Unlock()

		// 回写L1缓存
		data, _ := json.Marshal(dest)
		m.l1.Set(key, data, m.l1TTL)

		logs.Debugf("L2 cache hit: %s", key)
		return nil
	}

	// L2也未命中
	m.mu.Lock()
	m.stats.L2Misses++
	m.mu.Unlock()

	logs.Debugf("Cache miss: %s", key)
	return ErrCacheNotFound
}

// Set 设置缓存 (同时设置L1和L2)
func (m *MultiLevelCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value failed: %w", err)
	}

	// 设置L1缓存
	m.l1.Set(key, data, m.l1TTL)

	// 设置L2缓存
	if err := m.l2.Set(ctx, key, value); err != nil {
		logs.Errorf("Failed to set L2 cache: key=%s, error=%v", key, err)
		return err
	}

	logs.Debugf("Cache set: %s", key)
	return nil
}

// SetWithTTL 设置缓存并指定TTL
func (m *MultiLevelCache) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal cache value failed: %w", err)
	}

	// L1 TTL不超过5分钟
	l1TTL := ttl
	if l1TTL > 5*time.Minute {
		l1TTL = 5 * time.Minute
	}
	m.l1.Set(key, data, l1TTL)

	// L2使用指定TTL
	if err := m.l2.SetWithTTL(ctx, key, value, ttl); err != nil {
		logs.Errorf("Failed to set L2 cache: key=%s, error=%v", key, err)
		return err
	}

	return nil
}

// Delete 删除缓存 (同时删除L1和L2)
func (m *MultiLevelCache) Delete(ctx context.Context, key string) error {
	// 删除L1
	m.l1.Delete(key)

	// 删除L2
	if err := m.l2.Delete(ctx, key); err != nil {
		logs.Errorf("Failed to delete L2 cache: key=%s, error=%v", key, err)
		return err
	}

	logs.Debugf("Cache deleted: %s", key)
	return nil
}

// GetBatch 批量获取缓存
func (m *MultiLevelCache) GetBatch(ctx context.Context, keys []string) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	missingKeys := make([]string, 0, len(keys))

	// 先从L1获取
	for _, key := range keys {
		if x, found := m.l1.Get(key); found {
			var value interface{}
			if data, ok := x.([]byte); ok {
				if err := json.Unmarshal(data, &value); err == nil {
					results[key] = value
					m.mu.Lock()
					m.stats.L1Hits++
					m.mu.Unlock()
					continue
				}
			}
		}
		missingKeys = append(missingKeys, key)
	}

	// 从L2获取缺失的key
	if len(missingKeys) > 0 {
		l2Results, err := m.l2.GetBatch(ctx, missingKeys)
		if err != nil {
			logs.Errorf("Failed to get L2 batch cache: error=%v", err)
		} else {
			for key, value := range l2Results {
				results[key] = value
				// 回写L1
				data, _ := json.Marshal(value)
				m.l1.Set(key, data, m.l1TTL)
				m.mu.Lock()
				m.stats.L2Hits++
				m.mu.Unlock()
			}
		}
	}

	return results, nil
}

// SetBatch 批量设置缓存
func (m *MultiLevelCache) SetBatch(ctx context.Context, items map[string]interface{}) error {
	// 设置L1
	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			logs.Errorf("Failed to marshal batch cache: key=%s, error=%v", key, err)
			continue
		}
		m.l1.Set(key, data, m.l1TTL)
	}

	// 设置L2
	if err := m.l2.SetBatch(ctx, items); err != nil {
		logs.Errorf("Failed to set L2 batch cache: error=%v", err)
		return err
	}

	return nil
}

// DeleteByPattern 根据模式删除缓存
func (m *MultiLevelCache) DeleteByPattern(ctx context.Context, pattern string) error {
	// 删除L1中匹配的key
	for key := range m.l1.Items() {
		matched, err := matchPattern(key, pattern)
		if err != nil {
			continue
		}
		if matched {
			m.l1.Delete(key)
		}
	}

	// 删除L2中匹配的key
	if err := m.l2.DeleteByPattern(ctx, pattern); err != nil {
		logs.Errorf("Failed to delete L2 pattern cache: pattern=%s, error=%v", pattern, err)
		return err
	}

	return nil
}

// GetStats 获取缓存统计
func (m *MultiLevelCache) GetStats() *MultiLevelCacheStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本
	return &MultiLevelCacheStats{
		L1Hits:   m.stats.L1Hits,
		L1Misses: m.stats.L1Misses,
		L2Hits:   m.stats.L2Hits,
		L2Misses: m.stats.L2Misses,
	}
}

// GetHitRate 获取缓存命中率
func (m *MultiLevelCache) GetHitRate() (l1Rate, l2Rate float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalL1 := m.stats.L1Hits + m.stats.L1Misses
	if totalL1 > 0 {
		l1Rate = float64(m.stats.L1Hits) / float64(totalL1) * 100
	}

	totalL2 := m.stats.L2Hits + m.stats.L2Misses
	if totalL2 > 0 {
		l2Rate = float64(m.stats.L2Hits) / float64(totalL2) * 100
	}

	return
}

// Clear 清空所有缓存
func (m *MultiLevelCache) Clear(ctx context.Context) error {
	// 清空L1
	m.l1.Flush()

	// 清空L2 (通过删除模式)
	if err := m.l2.DeleteByPattern(ctx, "*"); err != nil {
		logs.Errorf("Failed to clear L2 cache: error=%v", err)
		return err
	}

	logs.Info("Multi-level cache cleared")
	return nil
}

// matchPattern 简单的模式匹配
func matchPattern(key, pattern string) (bool, error) {
	// 简化实现: 只支持 * 通配符
	// 生产环境应使用正则表达式或专门的匹配库
	if pattern == "*" {
		return true, nil
	}
	return false, nil
}
