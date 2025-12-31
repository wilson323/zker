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
)

// QuotaCache 配额缓存接口
type QuotaCache interface {
	// GetQuota 获取配额信息（带缓存）
	GetQuota(ctx context.Context, tenantID string, resourceType string, fetchFn func() (*QuotaInfo, error)) (*QuotaInfo, error)

	// ConsumeQuota 消费配额（原子操作，更新后删除缓存）
	ConsumeQuota(ctx context.Context, tenantID string, resourceType string, count int, consumeFn func(int) error) error

	// InvalidateQuota 使配额缓存失效
	InvalidateQuota(ctx context.Context, tenantID string, resourceType string) error

	// GetCacheStats 获取缓存统计
	GetCacheStats() QuotaCacheStats
}

// QuotaInfo 配额信息
type QuotaInfo struct {
	TenantID      string    `json:"tenant_id"`
	ResourceType  string    `json:"resource_type"`
	MaxLimit      int       `json:"max_limit"`
	UsedCount     int       `json:"used_count"`
	IsUnlimited   bool      `json:"is_unlimited"`
	LastUpdated   time.Time `json:"last_updated"`
}

// HasRemaining 检查是否有剩余配额
func (q *QuotaInfo) HasRemaining(required int) bool {
	if q.IsUnlimited {
		return true
	}
	return q.UsedCount+required <= q.MaxLimit
}

// GetRemaining 获取剩余配额
func (q *QuotaInfo) GetRemaining() int {
	if q.IsUnlimited {
		return -1 // -1 表示无限制
	}
	remaining := q.MaxLimit - q.UsedCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetUsagePercentage 获取使用率百分比
func (q *QuotaInfo) GetUsagePercentage() float64 {
	if q.IsUnlimited || q.MaxLimit == 0 {
		return 0
	}
	return float64(q.UsedCount) / float64(q.MaxLimit) * 100
}

// QuotaCacheStats 配额缓存统计
type QuotaCacheStats struct {
	HitCount      int64     // 命中次数
	MissCount     int64     // 未命中次数
	ConsumeCount  int64     // 消费次数
	ErrorCount    int64     // 错误次数
	LastCleanupAt time.Time // 最后清理时间
}

// QuotaCacheImpl 配额缓存实现
type QuotaCacheImpl struct {
	cache    Cmdable
	stats    *QuotaCacheStats
	mu       sync.Mutex
	localCache *LocalCache // L1本地缓存
}

// NewQuotaCache 创建配额缓存
func NewQuotaCache(cache Cmdable) *QuotaCacheImpl {
	return &QuotaCacheImpl{
		cache:      cache,
		stats:      &QuotaCacheStats{},
		localCache: NewLocalCache(500, time.Minute), // L1: 500条，1分钟
	}
}

// GetQuota 获取配额信息（带多级缓存）
// 缓存策略：
//   L1: 本地内存缓存（1分钟TTL）
//   L2: Redis缓存（5分钟TTL）
//   L3: 数据库查询
func (c *QuotaCacheImpl) GetQuota(
	ctx context.Context,
	tenantID string,
	resourceType string,
	fetchFn func() (*QuotaInfo, error),
) (*QuotaInfo, error) {
	const (
		localTTL = time.Minute
		redisTTL = 5 * time.Minute
	)

	cacheKey := fmt.Sprintf("%s:%s", tenantID, resourceType)

	// L1: 本地缓存
	if val, ok := c.localCache.Get(cacheKey); ok {
		c.stats.HitCount++
		logs.Debugf("Quota cache L1 hit: tenant=%s, resource=%s", tenantID, resourceType)
		// ✅ 修复: 安全的类型断言
		if quota, ok := val.(*QuotaInfo); ok {
			return quota, nil
		}
		logs.Warnf("Quota cache L1 type assertion failed: tenant=%s, resource=%s", tenantID, resourceType)
	}

	// L2: Redis缓存
	val, err := c.cache.Get(ctx, buildKey("quota", cacheKey)).Result()
	if err == nil {
		c.stats.HitCount++
		logs.Debugf("Quota cache L2 hit: tenant=%s, resource=%s", tenantID, resourceType)

		// ✅ 修复: 使用JSON完整反序列化
		var quota QuotaInfo
		if err := json.Unmarshal([]byte(val), &quota); err == nil {
			// 回写L1缓存
			c.localCache.Set(cacheKey, &quota, localTTL)
			return &quota, nil
		}
		logs.Warnf("Failed to unmarshal quota cache: tenant=%s, resource=%s, error=%v",
			tenantID, resourceType, err)
	}

	// L3: 数据库查询
	c.stats.MissCount++
	logs.Debugf("Quota cache miss: tenant=%s, resource=%s", tenantID, resourceType)

	quota, err := fetchFn()
	if err != nil {
		c.stats.ErrorCount++
		return nil, fmt.Errorf("failed to fetch quota: %w", err)
	}

	// 写回缓存
	if err := c.setQuota(ctx, cacheKey, quota, redisTTL); err != nil {
		logs.Errorf("Failed to set quota cache: tenant=%s, resource=%s, error=%v",
			tenantID, resourceType, err)
	}

	return quota, nil
}

// ConsumeQuota 消费配额（原子操作）
// 关键：先更新数据库，再删除缓存（保证一致性）
func (c *QuotaCacheImpl) ConsumeQuota(
	ctx context.Context,
	tenantID string,
	resourceType string,
	count int,
	consumeFn func(int) error,
) error {
	cacheKey := fmt.Sprintf("%s:%s", tenantID, resourceType)

	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. 先执行数据库更新（事务保证原子性）
	if err := consumeFn(count); err != nil {
		c.stats.ErrorCount++
		return fmt.Errorf("failed to consume quota in database: %w", err)
	}

	c.stats.ConsumeCount++

	// 2. 删除本地缓存
	c.localCache.Delete(cacheKey)

	// 3. 删除Redis缓存（下次访问时重新加载）
	if err := c.cache.Del(ctx, buildKey("quota", cacheKey)).Err(); err != nil {
		logs.Errorf("Failed to delete quota cache: tenant=%s, resource=%s, error=%v",
			tenantID, resourceType, err)
		return err
	}

	logs.Infof("Quota consumed and cache invalidated: tenant=%s, resource=%s, count=%d",
		tenantID, resourceType, count)

	return nil
}

// InvalidateQuota 使配额缓存失效
func (c *QuotaCacheImpl) InvalidateQuota(ctx context.Context, tenantID string, resourceType string) error {
	cacheKey := fmt.Sprintf("%s:%s", tenantID, resourceType)

	// 删除本地缓存
	c.localCache.Delete(cacheKey)

	// 删除Redis缓存
	if err := c.cache.Del(ctx, buildKey("quota", cacheKey)).Err(); err != nil {
		logs.Errorf("Failed to delete quota cache: tenant=%s, resource=%s, error=%v",
			tenantID, resourceType, err)
		return err
	}

	logs.Infof("Quota cache invalidated: tenant=%s, resource=%s", tenantID, resourceType)
	return nil
}

// GetCacheStats 获取缓存统计
func (c *QuotaCacheImpl) GetCacheStats() QuotaCacheStats {
	c.mu.Lock()
	defer c.mu.Unlock()

	return *c.stats
}

// GetHitRate 获取缓存命中率
func (c *QuotaCacheImpl) GetHitRate() float64 {
	total := c.stats.HitCount + c.stats.MissCount
	if total == 0 {
		return 0
	}
	return float64(c.stats.HitCount) / float64(total)
}

// setQuota 设置配额缓存
func (c *QuotaCacheImpl) setQuota(ctx context.Context, cacheKey string, quota *QuotaInfo, ttl time.Duration) error {
	// ✅ 修复: 使用JSON完整序列化（而非简化字符串）
	data, err := json.Marshal(quota)
	if err != nil {
		return fmt.Errorf("marshal quota failed: %w", err)
	}

	// 设置Redis缓存
	if err := c.cache.Set(ctx, buildKey("quota", cacheKey), data, ttl).Err(); err != nil {
		return err
	}

	// 设置本地缓存
	c.localCache.Set(cacheKey, quota, time.Minute)

	return nil
}
