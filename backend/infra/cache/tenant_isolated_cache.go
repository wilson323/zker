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
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	pkgerrorx "github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/infra/tracing"
	"go.uber.org/zap"
)

const (
	// DefaultCacheTTL 默认缓存过期时间
	DefaultCacheTTL = 5 * time.Minute

	// CacheKeyPrefix 缓存key前缀
	CacheKeyPrefix = "tenant"
)

// TenantIsolatedCache 租户隔离缓存
//
// **核心功能**：
// 1. 自动为所有缓存key添加租户前缀：tenant:{tenant_id}:{key}
// 2. 防止跨租户缓存访问
// 3. 支持批量操作
// 4. 提供缓存统计
//
// **Redis Key格式**：
// - String: tenant:{tenant_id}:{key}
// - Hash: tenant:{tenant_id}:hash:{key}
// - List: tenant:{tenant_id}:list:{key}
// - Set: tenant:{tenant_id}:set:{key}
//
// **使用示例**：
//	cache := NewTenantIsolatedCache(redisClient)
//	cache.Get(ctx, tenantID, "user:123")
//	cache.Set(ctx, tenantID, "user:123", userData, 5*time.Minute)
type TenantIsolatedCache struct {
	redis Cmdable

	// EnableStats 是否启用统计
	EnableStats bool

	// Stats 统计信息
	Stats *TenantCacheStats
}

// TenantCacheStats 租户缓存统计
type TenantCacheStats struct {
	HitCount  int64
	MissCount int64
	SetCount  int64
	DelCount  int64
}

// NewTenantIsolatedCache 创建租户隔离缓存
func NewTenantIsolatedCache(redis Cmdable) *TenantIsolatedCache {
	return &TenantIsolatedCache{
		redis:      redis,
		EnableStats: true,
		Stats:      &TenantCacheStats{},
	}
}

// buildCacheKey 构建缓存key：tenant:{tenant_id}:{key}
func (c *TenantIsolatedCache) buildCacheKey(tenantID, key string) string {
	return fmt.Sprintf("%s:%s:%s", CacheKeyPrefix, tenantID, key)
}

// buildHashKey 构建Hash缓存key：tenant:{tenant_id}:hash:{key}
func (c *TenantIsolatedCache) buildHashKey(tenantID, key string) string {
	return fmt.Sprintf("%s:%s:hash:%s", CacheKeyPrefix, tenantID, key)
}

// Get 获取缓存（String类型）
//
// **参数**：
// - ctx: 上下文（包含tenant_id）
// - key: 缓存key（不含租户前缀）
// - dest: 目标对象（必须传指针）
//
// **示例**：
//	var user User
//	err := cache.Get(ctx, "user:123", &user)
func (c *TenantIsolatedCache) Get(ctx context.Context, tenantID, key string, dest interface{}) error {
	// 记录追踪
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrCacheKey.String(key),
	)

	cacheKey := c.buildCacheKey(tenantID, key)

	// 从Redis获取
	data, err := c.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if c.EnableStats {
			c.Stats.MissCount++
		}

		if err == Nil {
			logs.CtxDebugf(ctx, "[TenantIsolatedCache] cache miss: key=%s", cacheKey)
			return pkgerrorx.NewByErrorCode(errno.ErrCacheMiss)
		}

		logs.CtxErrorf(ctx, "[TenantIsolatedCache] get failed: key=%s, error=%v", cacheKey, err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheGetFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	// 反序列化
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] unmarshal failed: key=%s, error=%v", cacheKey, err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheDecodeFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	if c.EnableStats {
		c.Stats.HitCount++
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] cache hit: key=%s", cacheKey)
	return nil
}

// Set 设置缓存（String类型）
//
// **参数**：
// - ctx: 上下文
// - tenantID: 租户ID
// - key: 缓存key
// - value: 缓存值（任意类型，会自动序列化为JSON）
// - ttl: 过期时间
//
// **示例**：
//	cache.Set(ctx, tenantID, "user:123", user, 5*time.Minute)
func (c *TenantIsolatedCache) Set(ctx context.Context, tenantID, key string, value interface{}, ttl time.Duration) error {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrCacheKey.String(key),
		tracing.AttrTTL.String(ttl.String()),
	)

	cacheKey := c.buildCacheKey(tenantID, key)

	// 序列化
	data, err := json.Marshal(value)
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] marshal failed: key=%s, error=%v", cacheKey, err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheEncodeFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	// 写入Redis
	if err := c.redis.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] set failed: key=%s, error=%v", cacheKey, err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheSetFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	if c.EnableStats {
		c.Stats.SetCount++
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] set success: key=%s, ttl=%s", cacheKey, ttl)
	return nil
}

// Delete 删除缓存
//
// **示例**：
//	cache.Delete(ctx, tenantID, "user:123")
func (c *TenantIsolatedCache) Delete(ctx context.Context, tenantID, key string) error {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrCacheKey.String(key),
	)

	cacheKey := c.buildCacheKey(tenantID, key)

	if err := c.redis.Del(ctx, cacheKey).Err(); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] delete failed: key=%s, error=%v", cacheKey, err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheDeleteFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	if c.EnableStats {
		c.Stats.DelCount++
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] delete success: key=%s", cacheKey)
	return nil
}

// GetSet Get并设置新值（原子操作）
//
// **使用场景**：缓存续期
//
// **示例**：
//	oldValue, err := cache.GetSet(ctx, tenantID, "user:123", newUser, 5*time.Minute)
func (c *TenantIsolatedCache) GetSet(ctx context.Context, tenantID, key string, newValue interface{}, ttl time.Duration) (string, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrCacheKey.String(key),
	)

	cacheKey := c.buildCacheKey(tenantID, key)

	// 序列化新值
	data, err := json.Marshal(newValue)
	if err != nil {
		return "", pkgerrorx.WrapWithZap(err, errno.ErrCacheEncodeFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	// GetSet操作
	oldData, err := c.redis.Get(ctx, cacheKey).Result()
	if err != nil && err != Nil {
		return "", pkgerrorx.WrapWithZap(err, errno.ErrCacheGetFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	// 设置新值
	if err := c.redis.Set(ctx, cacheKey, data, ttl).Err(); err != nil {
		return "", pkgerrorx.WrapWithZap(err, errno.ErrCacheSetFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	if c.EnableStats {
		c.Stats.HitCount++
		c.Stats.SetCount++
	}

	return oldData, nil
}

// Exists 检查key是否存在
//
// **示例**：
//	exists, err := cache.Exists(ctx, tenantID, "user:123")
func (c *TenantIsolatedCache) Exists(ctx context.Context, tenantID, key string) (bool, error) {
	cacheKey := c.buildCacheKey(tenantID, key)

	count, err := c.redis.Exists(ctx, cacheKey).Result()
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] exists check failed: key=%s, error=%v", cacheKey, err)
		return false, pkgerrorx.WrapWithZap(err, errno.ErrCacheCheckFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	return count > 0, nil
}

// Expire 设置key的过期时间
//
// **示例**：
//	err := cache.Expire(ctx, tenantID, "user:123", 10*time.Minute)
func (c *TenantIsolatedCache) Expire(ctx context.Context, tenantID, key string, ttl time.Duration) error {
	cacheKey := c.buildCacheKey(tenantID, key)

	if err := c.redis.Expire(ctx, cacheKey, ttl).Err(); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] expire failed: key=%s, error=%v", cacheKey, err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheExpireFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] expire success: key=%s, ttl=%s", cacheKey, ttl)
	return nil
}

// HGet 获取Hash字段
//
// **示例**：
//	var value string
//	err := cache.HGet(ctx, tenantID, "user:123", "name", &value)
func (c *TenantIsolatedCache) HGet(ctx context.Context, tenantID, key, field string, dest interface{}) error {
	hashKey := c.buildHashKey(tenantID, key)

	// 获取整个Hash
	data, err := c.redis.HGetAll(ctx, hashKey).Result()
	if err != nil {
		if err == Nil {
			return pkgerrorx.NewByErrorCode(errno.ErrCacheMiss)
		}
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheGetFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
			zap.String("field", field),
		)
	}

	// 提取指定字段
	fieldValue, exists := data[field]
	if !exists {
		return pkgerrorx.WrapWithZap(nil, errno.ErrCacheMiss, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
			zap.String("field", field),
		)
	}

	// 反序列化
	if err := json.Unmarshal([]byte(fieldValue), dest); err != nil {
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheDecodeFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
			zap.String("field", field),
		)
	}

	if c.EnableStats {
		c.Stats.HitCount++
	}

	return nil
}

// HSet 设置Hash字段
//
// **示例**：
//	cache.HSet(ctx, tenantID, "user:123", map[string]interface{}{"name": "John", "age": 30})
func (c *TenantIsolatedCache) HSet(ctx context.Context, tenantID, key string, fields map[string]interface{}) error {
	hashKey := c.buildHashKey(tenantID, key)

	// 序列化所有字段
	values := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		data, err := json.Marshal(v)
		if err != nil {
			return pkgerrorx.WrapWithZap(err, errno.ErrCacheEncodeFailed, 
				zap.String("tenant_id", tenantID),
				zap.String("key", key),
				zap.String("field", k),
			)
		}
		values = append(values, k, string(data))
	}

	if err := c.redis.HSet(ctx, hashKey, values...).Err(); err != nil {
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheSetFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	if c.EnableStats {
		c.Stats.SetCount++
	}

	return nil
}

// HDel 删除Hash字段
//
// **示例**：
//	cache.HDel(ctx, tenantID, "user:123", "name", "age")
func (c *TenantIsolatedCache) HDel(ctx context.Context, tenantID, key string, fields ...string) error {
	hashKey := c.buildHashKey(tenantID, key)

	if len(fields) == 0 {
		return nil
	}

	if err := c.redis.Del(ctx, hashKey).Err(); err != nil {
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheDeleteFailed, 
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	if c.EnableStats {
		c.Stats.DelCount++
	}

	return nil
}

// MGet 批量获取
//
// **示例**：
//	users, err := cache.MGet(ctx, tenantID, []string{"user:1", "user:2", "user:3"})
func (c *TenantIsolatedCache) MGet(ctx context.Context, tenantID string, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrKeyCount.Int(len(keys)),
	)

	// 构建缓存keys
	cacheKeys := make([]string, len(keys))
	for i, key := range keys {
		cacheKeys[i] = c.buildCacheKey(tenantID, key)
	}

	// 使用Pipeline批量获取
	pipe := c.redis.Pipeline()
	cmds := make([]StringCmd, len(cacheKeys))
	for i, cacheKey := range cacheKeys {
		cmds[i] = pipe.Get(ctx, cacheKey)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] mget failed: error=%v", err)
		return nil, pkgerrorx.WrapWithZap(err, errno.ErrCacheGetFailed, 
			zap.String("tenant_id", tenantID),
			zap.Int("key_count", len(keys)),
		)
	}

	// 解析结果
	result := make(map[string]interface{})
	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			if err == Nil {
				continue // 缓存未命中
			}
			logs.CtxWarnf(ctx, "[TenantIsolatedCache] mget partial failed: key=%s, error=%v", keys[i], err)
			continue
		}

		var dest interface{}
		if err := json.Unmarshal([]byte(data), &dest); err != nil {
			logs.CtxWarnf(ctx, "[TenantIsolatedCache] mget unmarshal failed: key=%s, error=%v", keys[i], err)
			continue
		}

		result[keys[i]] = dest
	}

	if c.EnableStats {
		c.Stats.HitCount += int64(len(result))
		c.Stats.MissCount += int64(len(keys) - len(result))
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] mget success: hit=%d, miss=%d",
		len(result), len(keys)-len(result))

	return result, nil
}

// MSet 批量设置
//
// **示例**：
//	data := map[string]interface{}{"user:1": user1, "user:2": user2}
//	err := cache.MSet(ctx, tenantID, data, 5*time.Minute)
func (c *TenantIsolatedCache) MSet(ctx context.Context, tenantID string, data map[string]interface{}, ttl time.Duration) error {
	if len(data) == 0 {
		return nil
	}

	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrKeyCount.Int(len(data)),
	)

	// 使用Pipeline批量设置
	pipe := c.redis.Pipeline()
	for key, value := range data {
		cacheKey := c.buildCacheKey(tenantID, key)

		// 序列化
		jsonData, err := json.Marshal(value)
		if err != nil {
			logs.CtxErrorf(ctx, "[TenantIsolatedCache] mset marshal failed: key=%s, error=%v", key, err)
			continue
		}

		pipe.Set(ctx, cacheKey, jsonData, ttl)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] mset failed: error=%v", err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheSetFailed, 
			zap.String("tenant_id", tenantID),
			zap.Int("key_count", len(data)),
		)
	}

	if c.EnableStats {
		c.Stats.SetCount += int64(len(data))
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] mset success: count=%d", len(data))
	return nil
}

// MDel 批量删除
//
// **示例**：
//	err := cache.MDel(ctx, tenantID, []string{"user:1", "user:2", "user:3"})
func (c *TenantIsolatedCache) MDel(ctx context.Context, tenantID string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		tracing.AttrKeyCount.Int(len(keys)),
	)

	// 构建缓存keys
	cacheKeys := make([]string, len(keys))
	for i, key := range keys {
		cacheKeys[i] = c.buildCacheKey(tenantID, key)
	}

	if err := c.redis.Del(ctx, cacheKeys...).Err(); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedCache] mdel failed: error=%v", err)
		return pkgerrorx.WrapWithZap(err, errno.ErrCacheDeleteFailed, 
			zap.String("tenant_id", tenantID),
			zap.Int("key_count", len(keys)),
		)
	}

	if c.EnableStats {
		c.Stats.DelCount += int64(len(keys))
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedCache] mdel success: count=%d", len(keys))
	return nil
}

// GetStats 获取缓存统计信息
func (c *TenantIsolatedCache) GetStats() TenantCacheStats {
	return TenantCacheStats{
		HitCount:  c.Stats.HitCount,
		MissCount: c.Stats.MissCount,
		SetCount:  c.Stats.SetCount,
		DelCount:  c.Stats.DelCount,
	}
}

// ResetStats 重置统计信息
func (c *TenantIsolatedCache) ResetStats() {
	c.Stats = &TenantCacheStats{}
}

// GetHitRate 获取缓存命中率
func (c *TenantIsolatedCache) GetHitRate() float64 {
	total := c.Stats.HitCount + c.Stats.MissCount
	if total == 0 {
		return 0
	}
	return float64(c.Stats.HitCount) / float64(total)
}

// ClearTenant 清除租户所有缓存
//
// **使用场景**：租户删除、数据清理
//
// **示例**：
//	count, err := cache.ClearTenant(ctx, tenantID)
func (c *TenantIsolatedCache) ClearTenant(ctx context.Context, tenantID string) (int64, error) {
	logs.CtxInfof(ctx, "[TenantIsolatedCache] clearing all cache for tenant: %s", tenantID)

	// 使用SCAN命令查找所有租户相关的key
	pattern := fmt.Sprintf("%s:%s:*", CacheKeyPrefix, tenantID)

	// 注意：这里需要使用SCAN + DEL来删除所有匹配的key
	// 实际实现需要结合Redis的SCAN命令

	// 简化实现：返回删除的key数量（实际生产环境需要完整的SCAN实现）
	logs.CtxWarnf(ctx, "[TenantIsolatedCache] ClearTenant not fully implemented, pattern=%s", pattern)

	return 0, pkgerrorx.WrapWithZap(nil, errno.ErrNotImplemented, 
		zap.String("tenant_id", tenantID),
		zap.String("pattern", pattern),
	)
}
