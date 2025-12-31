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
)

// 以下类型和常量定义在 permission_cache.go 中，已在此引用

// PermissionCacheUnified 统一的权限缓存实现
// ✅ 整合改进:
//   1. 使用 Cmdable 接口而非 *redis.Client (统一Redis访问)
//   2. 使用统一的 "zker:" 前缀而非 "coze:studio:"
//   3. 添加多级缓存支持 (L1本地 + L2 Redis)
//   4. 完整的统计信息
// 注意: PermissionCacheStats 类型已在 permission_cache.go 中定义
type PermissionCacheUnified struct {
	cache      Cmdable            // 使用Cmdable接口
	localCache *LocalCache        // L1本地缓存
	stats      *PermissionCacheStats
	l1TTL      time.Duration      // L1缓存TTL
}

// NewPermissionCacheUnified 创建统一权限缓存实例
func NewPermissionCacheUnified(cache Cmdable) *PermissionCacheUnified {
	return &PermissionCacheUnified{
		cache:      cache,
		localCache: NewLocalCache(1000, time.Minute), // L1: 1000条，1分钟
		stats:      &PermissionCacheStats{},
		l1TTL:      time.Minute,
	}
}

// GetUserRoles 获取用户角色缓存 (多级缓存)
func (c *PermissionCacheUnified) GetUserRoles(ctx context.Context, userID string) ([]*RoleInfo, error) {
	cacheKey := buildKey(cacheKeyPrefixUserRoles, userID)

	// L1: 本地缓存
	if val, ok := c.localCache.Get(cacheKey); ok {
		c.stats.mu.Lock()
		c.stats.L1Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		if roles, ok := val.([]*RoleInfo); ok {
			logs.Debugf("Permission cache L1 hit: user_roles, user=%s", userID)
			return roles, nil
		}
	}

	// L2: Redis缓存
	val, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		c.stats.mu.Lock()
		c.stats.L2Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		var roles []*RoleInfo
		if err := json.Unmarshal([]byte(val), &roles); err == nil {
			// 回写L1缓存
			c.localCache.Set(cacheKey, roles, c.l1TTL)
			logs.Debugf("Permission cache L2 hit: user_roles, user=%s", userID)
			return roles, nil
		}
	}

	// 缓存未命中
	c.stats.mu.Lock()
	c.stats.MissCount++
	c.stats.mu.Unlock()

	return nil, nil
}

// SetUserRoles 设置用户角色缓存
func (c *PermissionCacheUnified) SetUserRoles(ctx context.Context, userID string, roles []*RoleInfo) error {
	cacheKey := buildKey(cacheKeyPrefixUserRoles, userID)

	data, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("marshal user roles failed: %w", err)
	}

	// 设置L2缓存
	if err := c.cache.Set(ctx, cacheKey, data, cacheTTLUserRoles).Err(); err != nil {
		return fmt.Errorf("set Redis cache failed: %w", err)
	}

	// 设置L1缓存
	c.localCache.Set(cacheKey, roles, c.l1TTL)

	logs.Debugf("Cached user roles: user=%s, count=%d", userID, len(roles))
	return nil
}

// GetPermissionRules 获取权限规则缓存
func (c *PermissionCacheUnified) GetPermissionRules(ctx context.Context, roleID string) ([]*PermissionRule, error) {
	cacheKey := buildKey(cacheKeyPrefixRolePermissions, roleID)

	// L1: 本地缓存
	if val, ok := c.localCache.Get(cacheKey); ok {
		c.stats.mu.Lock()
		c.stats.L1Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		if rules, ok := val.([]*PermissionRule); ok {
			logs.Debugf("Permission cache L1 hit: role_perms, role=%s", roleID)
			return rules, nil
		}
	}

	// L2: Redis缓存
	val, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		c.stats.mu.Lock()
		c.stats.L2Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		var rules []*PermissionRule
		if err := json.Unmarshal([]byte(val), &rules); err == nil {
			// 回写L1缓存
			c.localCache.Set(cacheKey, rules, c.l1TTL)
			logs.Debugf("Permission cache L2 hit: role_perms, role=%s", roleID)
			return rules, nil
		}
	}

	c.stats.mu.Lock()
	c.stats.MissCount++
	c.stats.mu.Unlock()

	return nil, nil
}

// SetPermissionRules 设置权限规则缓存
func (c *PermissionCacheUnified) SetPermissionRules(ctx context.Context, roleID string, rules []*PermissionRule) error {
	cacheKey := buildKey(cacheKeyPrefixRolePermissions, roleID)

	data, err := json.Marshal(rules)
	if err != nil {
		return fmt.Errorf("marshal permission rules failed: %w", err)
	}

	if err := c.cache.Set(ctx, cacheKey, data, cacheTTLRolePermissions).Err(); err != nil {
		return fmt.Errorf("set Redis cache failed: %w", err)
	}

	c.localCache.Set(cacheKey, rules, c.l1TTL)

	logs.Debugf("Cached permission rules: role=%s, count=%d", roleID, len(rules))
	return nil
}

// GetDataPermissionFilter 获取数据权限过滤器缓存
func (c *PermissionCacheUnified) GetDataPermissionFilter(ctx context.Context, userID, resourceType string) (*DataFilterInfo, error) {
	cacheKey := buildKey(cacheKeyPrefixDataFilter, userID, ":", resourceType)

	// L1: 本地缓存
	if val, ok := c.localCache.Get(cacheKey); ok {
		c.stats.mu.Lock()
		c.stats.L1Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		if filter, ok := val.(*DataFilterInfo); ok {
			logs.Debugf("Permission cache L1 hit: data_filter, user=%s, resource=%s", userID, resourceType)
			return filter, nil
		}
	}

	// L2: Redis缓存
	val, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		c.stats.mu.Lock()
		c.stats.L2Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		var filter DataFilterInfo
		if err := json.Unmarshal([]byte(val), &filter); err == nil {
			// 回写L1缓存
			c.localCache.Set(cacheKey, &filter, c.l1TTL)
			logs.Debugf("Permission cache L2 hit: data_filter, user=%s, resource=%s", userID, resourceType)
			return &filter, nil
		}
	}

	c.stats.mu.Lock()
	c.stats.MissCount++
	c.stats.mu.Unlock()

	return nil, nil
}

// SetDataPermissionFilter 设置数据权限过滤器缓存
func (c *PermissionCacheUnified) SetDataPermissionFilter(ctx context.Context, userID, resourceType string, filter *DataFilterInfo) error {
	cacheKey := buildKey(cacheKeyPrefixDataFilter, userID, ":", resourceType)

	filter.CachedAt = time.Now()

	data, err := json.Marshal(filter)
	if err != nil {
		return fmt.Errorf("marshal data filter failed: %w", err)
	}

	if err := c.cache.Set(ctx, cacheKey, data, cacheTTLDataFilter).Err(); err != nil {
		return fmt.Errorf("set Redis cache failed: %w", err)
	}

	c.localCache.Set(cacheKey, filter, c.l1TTL)

	logs.Debugf("Cached data filter: user=%s, resource=%s", userID, resourceType)
	return nil
}

// GetFieldPermissions 获取字段权限缓存
func (c *PermissionCacheUnified) GetFieldPermissions(ctx context.Context, userID, resourceType string) (map[string]*FieldPermInfo, error) {
	cacheKey := buildKey(cacheKeyPrefixFieldPerms, userID, ":", resourceType)

	// L1: 本地缓存
	if val, ok := c.localCache.Get(cacheKey); ok {
		c.stats.mu.Lock()
		c.stats.L1Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		if perms, ok := val.(map[string]*FieldPermInfo); ok {
			logs.Debugf("Permission cache L1 hit: field_perms, user=%s, resource=%s", userID, resourceType)
			return perms, nil
		}
	}

	// L2: Redis缓存
	val, err := c.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		c.stats.mu.Lock()
		c.stats.L2Hits++
		c.stats.HitCount++
		c.stats.mu.Unlock()

		var perms map[string]*FieldPermInfo
		if err := json.Unmarshal([]byte(val), &perms); err == nil {
			// 回写L1缓存
			c.localCache.Set(cacheKey, perms, c.l1TTL)
			logs.Debugf("Permission cache L2 hit: field_perms, user=%s, resource=%s", userID, resourceType)
			return perms, nil
		}
	}

	c.stats.mu.Lock()
	c.stats.MissCount++
	c.stats.mu.Unlock()

	return nil, nil
}

// SetFieldPermissions 设置字段权限缓存
func (c *PermissionCacheUnified) SetFieldPermissions(ctx context.Context, userID, resourceType string, perms map[string]*FieldPermInfo) error {
	cacheKey := buildKey(cacheKeyPrefixFieldPerms, userID, ":", resourceType)

	data, err := json.Marshal(perms)
	if err != nil {
		return fmt.Errorf("marshal field permissions failed: %w", err)
	}

	if err := c.cache.Set(ctx, cacheKey, data, cacheTTLFieldPerms).Err(); err != nil {
		return fmt.Errorf("set Redis cache failed: %w", err)
	}

	c.localCache.Set(cacheKey, perms, c.l1TTL)

	logs.Debugf("Cached field permissions: user=%s, resource=%s, count=%d", userID, resourceType, len(perms))
	return nil
}

// InvalidateUser 失效用户相关缓存
func (c *PermissionCacheUnified) InvalidateUser(ctx context.Context, userID string) error {
	// 删除L1本地缓存
	c.localCache.Clear()

	// 删除L2 Redis缓存（使用SCAN查找所有相关key）
	pattern := buildKey(cacheKeyPrefixUserRoles, userID)
	keys, _ := c.scanKeys(ctx, pattern)

	// 数据权限过滤器
	filterPattern := buildKey(cacheKeyPrefixDataFilter, userID, ":*")
	filterKeys, _ := c.scanKeys(ctx, filterPattern)
	keys = append(keys, filterKeys...)

	// 字段权限
	fieldPattern := buildKey(cacheKeyPrefixFieldPerms, userID, ":*")
	fieldKeys, _ := c.scanKeys(ctx, fieldPattern)
	keys = append(keys, fieldKeys...)

	if len(keys) > 0 {
		if err := c.cache.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("delete Redis cache failed: %w", err)
		}
		logs.Infof("Invalidated user cache: user=%s, keys=%d", userID, len(keys))
	}

	return nil
}

// InvalidateRole 失效角色相关缓存
func (c *PermissionCacheUnified) InvalidateRole(ctx context.Context, roleID string) error {
	// 删除L1本地缓存
	c.localCache.Clear()

	// 删除L2 Redis缓存
	cacheKey := buildKey(cacheKeyPrefixRolePermissions, roleID)
	if err := c.cache.Del(ctx, cacheKey).Err(); err != nil {
		return fmt.Errorf("delete Redis cache failed: %w", err)
	}

	logs.Infof("Invalidated role cache: role=%s", roleID)
	return nil
}

// InvalidateRoleUsers 失效角色的所有用户缓存
func (c *PermissionCacheUnified) InvalidateRoleUsers(ctx context.Context, roleID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	// 删除L1本地缓存
	c.localCache.Clear()

	// 构建所有需要删除的缓存键
	keys := make([]string, 0, len(userIDs)*3)

	for _, userID := range userIDs {
		// 用户角色缓存
		keys = append(keys, buildKey(cacheKeyPrefixUserRoles, userID))

		// 数据权限过滤器缓存（所有资源类型）
		filterPattern := buildKey(cacheKeyPrefixDataFilter, userID, ":*")
		filterKeys, _ := c.scanKeys(ctx, filterPattern)
		keys = append(keys, filterKeys...)

		// 字段权限缓存（所有资源类型）
		fieldPattern := buildKey(cacheKeyPrefixFieldPerms, userID, ":*")
		fieldKeys, _ := c.scanKeys(ctx, fieldPattern)
		keys = append(keys, fieldKeys...)
	}

	// 批量删除
	if len(keys) > 0 {
		if err := c.cache.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("delete Redis cache failed: %w", err)
		}
		logs.Infof("Invalidated role users cache: role=%s, users=%d, keys=%d", roleID, len(userIDs), len(keys))
	}

	return nil
}

// GetCacheStats 获取缓存统计信息
func (c *PermissionCacheUnified) GetCacheStats() *PermissionCacheStats {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	// 返回副本
	return &PermissionCacheStats{
		HitCount:  c.stats.HitCount,
		MissCount: c.stats.MissCount,
		EvictCount: c.stats.EvictCount,
		L1Hits:    c.stats.L1Hits,
		L2Hits:    c.stats.L2Hits,
	}
}

// scanKeys 扫描匹配的键 (辅助方法)
func (c *PermissionCacheUnified) scanKeys(ctx context.Context, pattern string) ([]string, error) {
	// 使用SCAN命令避免阻塞
	scanCmd := c.cache.(interface {
		Scan(ctx context.Context, cursor uint64, match string, count int64) ScanCmd
	}).Scan(ctx, 0, pattern, 100)

	keys := []string{}
	iter := scanCmd.Iterator()

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		logs.Errorf("Failed to scan keys: pattern=%s, error=%v", pattern, err)
		return nil, err
	}

	return keys, nil
}
