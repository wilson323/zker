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
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// ConfigCache 系统配置缓存接口
type ConfigCache interface {
	// GetConfig 获取配置值（带缓存）
	GetConfig(ctx context.Context, key string, fetchFn func() (string, error)) (string, error)

	// SetConfig 设置配置值并失效缓存
	SetConfig(ctx context.Context, key string, value string, updateFn func(string) error) error

	// DeleteConfig 删除配置值并失效缓存
	DeleteConfig(ctx context.Context, key string, deleteFn func() error) error

	// GetAllConfigs 获取所有配置（带缓存）
	GetAllConfigs(ctx context.Context, fetchFn func() (map[string]string, error)) (map[string]string, error)

	// InvalidateAllConfigs 使所有配置缓存失效
	InvalidateAllConfigs(ctx context.Context) error

	// GetCacheStats 获取缓存统计
	GetCacheStats() ConfigCacheStats
}

// ConfigCacheStats 配置缓存统计
type ConfigCacheStats struct {
	HitCount       int64     // 命中次数
	MissCount      int64     // 未命中次数
	UpdateCount    int64     // 更新次数
	DeleteCount    int64     // 删除次数
	LastRefreshAt  time.Time // 最后刷新时间
}

// ConfigCacheImpl 系统配置缓存实现
type ConfigCacheImpl struct {
	cache       Cmdable
	localCache  *LocalCache // L1本地缓存
	stats       *ConfigCacheStats
	mu          sync.RWMutex
	allConfigs  map[string]string // 本地全量配置缓存
	initialized bool
}

// NewConfigCache 创建配置缓存
func NewConfigCache(cache Cmdable) *ConfigCacheImpl {
	return &ConfigCacheImpl{
		cache:      cache,
		localCache: NewLocalCache(200, 10*time.Minute), // L1: 200条，10分钟
		stats:      &ConfigCacheStats{},
		allConfigs: make(map[string]string),
	}
}

// GetConfig 获取配置值（带多级缓存）
// 缓存策略：
//   L1: 本地全量配置缓存（应用启动时加载，定时刷新）
//   L2: 本内存缓存（10分钟TTL）
//   L3: Redis缓存（10分钟TTL）
//   L4: 数据库查询
func (c *ConfigCacheImpl) GetConfig(
	ctx context.Context,
	key string,
	fetchFn func() (string, error),
) (string, error) {
	// L1: 尝试从全量配置缓存读取
	c.mu.RLock()
	if value, ok := c.allConfigs[key]; ok && c.initialized {
		c.mu.RUnlock()
		c.stats.HitCount++
		logs.Debugf("Config cache L1 hit: key=%s", key)
		return value, nil
	}
	c.mu.RUnlock()

	// L2: 本地内存缓存
	if value, ok := c.localCache.Get(key); ok {
		// ✅ 修复: 安全的类型断言
		if strValue, ok := value.(string); ok {
			c.mu.Lock()
			c.allConfigs[key] = strValue
			c.mu.Unlock()
			c.stats.HitCount++
			logs.Debugf("Config cache L2 hit: key=%s", key)
			return strValue, nil
		}
		logs.Warnf("Config cache L2 type assertion failed: key=%s", key)
	}

	// L3: Redis缓存
	val, err := c.cache.Get(ctx, buildKey("config", key)).Result()
	if err == nil {
		c.mu.Lock()
		c.allConfigs[key] = val
		c.mu.Unlock()
		c.stats.HitCount++
		c.localCache.Set(key, val, 10*time.Minute)
		logs.Debugf("Config cache L3 hit: key=%s", key)
		return val, nil
	}

	// L4: 数据库查询
	c.stats.MissCount++
	logs.Debugf("Config cache miss: key=%s", key)

	value, err := fetchFn()
	if err != nil {
		return "", fmt.Errorf("failed to fetch config: %w", err)
	}

	// 写回缓存
	c.setConfig(ctx, key, value)

	return value, nil
}

// SetConfig 设置配置值并失效缓存
func (c *ConfigCacheImpl) SetConfig(
	ctx context.Context,
	key string,
	value string,
	updateFn func(string) error,
) error {
	// 1. 先更新数据库
	if err := updateFn(value); err != nil {
		return fmt.Errorf("failed to update config in database: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.UpdateCount++

	// 2. 更新本地全量缓存
	c.allConfigs[key] = value

	// 3. 更新本地缓存
	c.localCache.Set(key, value, 10*time.Minute)

	// 4. 更新Redis缓存
	if err := c.cache.Set(ctx, buildKey("config", key), value, 10*time.Minute).Err(); err != nil {
		logs.Errorf("Failed to set config cache: key=%s, error=%v", key, err)
		return err
	}

	logs.Infof("Config updated and cache refreshed: key=%s", key)
	return nil
}

// DeleteConfig 删除配置值并失效缓存
func (c *ConfigCacheImpl) DeleteConfig(
	ctx context.Context,
	key string,
	deleteFn func() error,
) error {
	// 1. 先删除数据库记录
	if err := deleteFn(); err != nil {
		return fmt.Errorf("failed to delete config from database: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats.DeleteCount++

	// 2. 删除本地全量缓存
	delete(c.allConfigs, key)

	// 3. 删除本地缓存
	c.localCache.Delete(key)

	// 4. 删除Redis缓存
	if err := c.cache.Del(ctx, buildKey("config", key)).Err(); err != nil {
		logs.Errorf("Failed to delete config cache: key=%s, error=%v", key, err)
		return err
	}

	logs.Infof("Config deleted and cache invalidated: key=%s", key)
	return nil
}

// GetAllConfigs 获取所有配置（带缓存）
func (c *ConfigCacheImpl) GetAllConfigs(
	ctx context.Context,
	fetchFn func() (map[string]string, error),
) (map[string]string, error) {
	c.mu.RLock()
	if c.initialized && len(c.allConfigs) > 0 {
		c.mu.RUnlock()
		c.stats.HitCount++
		logs.Debugf("All configs cache hit: count=%d", len(c.allConfigs))
		// 返回副本
		result := make(map[string]string, len(c.allConfigs))
		for k, v := range c.allConfigs {
			result[k] = v
		}
		return result, nil
	}
	c.mu.RUnlock()

	// 从数据库加载所有配置
	c.stats.MissCount++
	logs.Debugf("All configs cache miss, loading from database")

	configs, err := fetchFn()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all configs: %w", err)
	}

	// 更新缓存
	c.mu.Lock()
	c.allConfigs = configs
	c.initialized = true
	c.stats.LastRefreshAt = time.Now()
	c.mu.Unlock()

	logs.Infof("All configs loaded and cached: count=%d", len(configs))
	return configs, nil
}

// InvalidateAllConfigs 使所有配置缓存失效
func (c *ConfigCacheImpl) InvalidateAllConfigs(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 清空本地全量缓存
	c.allConfigs = make(map[string]string)
	c.initialized = false

	// 清空本地缓存
	c.localCache.Clear()

	// 清空Redis配置缓存（使用SCAN查找所有config:开头的key）
	pattern := buildKey("config", "*")
	scanResult := c.cache.(interface {
		Scan(ctx context.Context, cursor uint64, match string, count int64) ScanCmd
	}).Scan(ctx, 0, pattern, 100)

	iter := scanResult.Iterator()

	keys := []string{}
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if len(keys) > 0 {
		if err := c.cache.Del(ctx, keys...).Err(); err != nil {
			logs.Errorf("Failed to delete all config cache: error=%v", err)
			return err
		}
	}

	logs.Infof("All config cache invalidated: count=%d", len(keys))
	return nil
}

// GetCacheStats 获取缓存统计
func (c *ConfigCacheImpl) GetCacheStats() ConfigCacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return *c.stats
}

// GetHitRate 获取缓存命中率
func (c *ConfigCacheImpl) GetHitRate() float64 {
	total := c.stats.HitCount + c.stats.MissCount
	if total == 0 {
		return 0
	}
	return float64(c.stats.HitCount) / float64(total)
}

// RefreshAllConfigs 刷新所有配置（定时任务调用）
func (c *ConfigCacheImpl) RefreshAllConfigs(ctx context.Context, fetchFn func() (map[string]string, error)) error {
	configs, err := fetchFn()
	if err != nil {
		return fmt.Errorf("failed to refresh all configs: %w", err)
	}

	c.mu.Lock()
	c.allConfigs = configs
	c.initialized = true
	c.stats.LastRefreshAt = time.Now()
	c.mu.Unlock()

	logs.Infof("All configs refreshed: count=%d", len(configs))
	return nil
}

// setConfig 设置配置缓存（内部方法）
func (c *ConfigCacheImpl) setConfig(ctx context.Context, key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 更新本地全量缓存
	c.allConfigs[key] = value

	// 更新本地缓存
	c.localCache.Set(key, value, 10*time.Minute)

	// 更新Redis缓存
	if err := c.cache.Set(ctx, buildKey("config", key), value, 10*time.Minute).Err(); err != nil {
		logs.Errorf("Failed to set config cache: key=%s, error=%v", key, err)
	}
}
