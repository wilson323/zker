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
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// Cache key prefixes
	cacheKeyPrefixUserRoles       = "perm:user:roles:"        // 用户角色缓存
	cacheKeyPrefixRolePermissions = "perm:role:perms:"       // 角色权限缓存
	cacheKeyPrefixDataFilter      = "perm:data:filter:"      // 数据权限过滤器缓存
	cacheKeyPrefixFieldPerms      = "perm:field:perms:"      // 字段权限缓存

	// Cache TTL durations
	cacheTTLUserRoles       = 1 * time.Hour  // 用户角色缓存：1小时
	cacheTTLRolePermissions = 30 * time.Minute // 角色权限缓存：30分钟
	cacheTTLDataFilter      = 15 * time.Minute // 数据权限过滤器缓存：15分钟
	cacheTTLFieldPerms      = 30 * time.Minute // 字段权限缓存：30分钟
)

// RoleInfo 角色信息
type RoleInfo struct {
	RoleID   string `json:"role_id"`
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

// PermissionRule 权限规则
type PermissionRule struct {
	ResourceType string `json:"resource_type"`
	Scope        string `json:"scope"`
	CustomFilter string `json:"custom_filter,omitempty"`
}

// DataFilterInfo 数据权限过滤器信息
type DataFilterInfo struct {
	WhereClause string        `json:"where_clause"`
	Args        []interface{} `json:"args"`
	CachedAt    time.Time     `json:"cached_at"`
}

// FieldPermInfo 字段权限信息
type FieldPermInfo struct {
	FieldName       string `json:"field_name"`
	PermissionLevel string `json:"permission_level"`
	Permission      int    `json:"permission"`
}

// PermissionCache 权限缓存接口
type PermissionCache interface {
	// GetUserRoles 获取用户角色缓存
	GetUserRoles(ctx context.Context, userID string) ([]*RoleInfo, error)

	// SetUserRoles 设置用户角色缓存
	SetUserRoles(ctx context.Context, userID string, roles []*RoleInfo) error

	// GetPermissionRules 获取权限规则缓存
	GetPermissionRules(ctx context.Context, roleID string) ([]*PermissionRule, error)

	// SetPermissionRules 设置权限规则缓存
	SetPermissionRules(ctx context.Context, roleID string, rules []*PermissionRule) error

	// GetDataPermissionFilter 获取数据权限过滤器缓存
	GetDataPermissionFilter(ctx context.Context, userID, resourceType string) (*DataFilterInfo, error)

	// SetDataPermissionFilter 设置数据权限过滤器缓存
	SetDataPermissionFilter(ctx context.Context, userID, resourceType string, filter *DataFilterInfo) error

	// GetFieldPermissions 获取字段权限缓存（返回 map[string]string 格式）
	GetFieldPermissions(ctx context.Context, cacheKey string, loader func() (map[string]string, error)) (map[string]string, error)

	// SetFieldPermissions 设置字段权限缓存
	SetFieldPermissions(ctx context.Context, userID, resourceType string, perms map[string]*FieldPermInfo) error

	// CheckDataPermission 检查数据权限（带缓存）
	CheckDataPermission(ctx context.Context, cacheKey string, loader func() (bool, error)) (bool, error)

	// InvalidateUserPermissions 使用户权限缓存失效
	InvalidateUserPermissions(ctx context.Context, userID string) error

	// InvalidateDataPermission 使特定数据权限缓存失效
	InvalidateDataPermission(ctx context.Context, cacheKey string) error

	// InvalidateUser 失效用户相关缓存
	InvalidateUser(ctx context.Context, userID string) error

	// InvalidateRole 失效角色相关缓存
	InvalidateRole(ctx context.Context, roleID string) error

	// InvalidateRoleUsers 失效角色的所有用户缓存
	InvalidateRoleUsers(ctx context.Context, roleID string, userIDs []string) error

	// GetCacheStats 获取缓存统计信息
	GetCacheStats() PermissionCacheStats
}

// PermissionCacheStats 权限缓存统计信息
type PermissionCacheStats struct {
	HitCount  uint64 `json:"hit_count"`
	MissCount uint64 `json:"miss_count"`
	EvictCount uint64 `json:"evict_count"`
	L1Hits    uint64 `json:"l1_hits"` // L1缓存命中
	L2Hits    uint64 `json:"l2_hits"` // L2缓存命中
	mu        sync.RWMutex
}

// HitRate 缓存命中率
func (s *PermissionCacheStats) HitRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := s.HitCount + s.MissCount
	if total == 0 {
		return 0
	}
	return float64(s.HitCount) / float64(total)
}

// L1HitRate L1缓存命中率
func (s *PermissionCacheStats) L1HitRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.HitCount == 0 {
		return 0
	}
	return float64(s.L1Hits) / float64(s.HitCount)
}

// PermissionCacheImpl 权限缓存实现
type PermissionCacheImpl struct {
	redis  *redis.Client
	logger *zap.Logger
}

// NewPermissionCache 创建权限缓存实例
func NewPermissionCache(redisClient *redis.Client, logger *zap.Logger) PermissionCache {
	return &PermissionCacheImpl{
		redis:  redisClient,
		logger: logger,
	}
}

// GetUserRoles 获取用户角色缓存
func (c *PermissionCacheImpl) GetUserRoles(ctx context.Context, userID string) ([]*RoleInfo, error) {
	key := c.buildCacheKey(cacheKeyPrefixUserRoles, userID)

	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		c.logger.Error("failed to get user roles from cache",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	var roles []*RoleInfo
	if err := json.Unmarshal([]byte(val), &roles); err != nil {
		c.logger.Error("failed to unmarshal user roles",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	c.logger.Debug("cache hit for user roles",
		zap.String("user_id", userID),
		zap.Int("count", len(roles)))

	return roles, nil
}

// SetUserRoles 设置用户角色缓存
func (c *PermissionCacheImpl) SetUserRoles(ctx context.Context, userID string, roles []*RoleInfo) error {
	key := c.buildCacheKey(cacheKeyPrefixUserRoles, userID)

	data, err := json.Marshal(roles)
	if err != nil {
		c.logger.Error("failed to marshal user roles",
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}

	if err := c.redis.Set(ctx, key, data, cacheTTLUserRoles).Err(); err != nil {
		c.logger.Error("failed to set user roles to cache",
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}

	c.logger.Debug("cached user roles",
		zap.String("user_id", userID),
		zap.Int("count", len(roles)))

	return nil
}

// GetPermissionRules 获取权限规则缓存
func (c *PermissionCacheImpl) GetPermissionRules(ctx context.Context, roleID string) ([]*PermissionRule, error) {
	key := c.buildCacheKey(cacheKeyPrefixRolePermissions, roleID)

	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		c.logger.Error("failed to get permission rules from cache",
			zap.String("role_id", roleID),
			zap.Error(err))
		return nil, err
	}

	var rules []*PermissionRule
	if err := json.Unmarshal([]byte(val), &rules); err != nil {
		c.logger.Error("failed to unmarshal permission rules",
			zap.String("role_id", roleID),
			zap.Error(err))
		return nil, err
	}

	c.logger.Debug("cache hit for permission rules",
		zap.String("role_id", roleID),
		zap.Int("count", len(rules)))

	return rules, nil
}

// SetPermissionRules 设置权限规则缓存
func (c *PermissionCacheImpl) SetPermissionRules(ctx context.Context, roleID string, rules []*PermissionRule) error {
	key := c.buildCacheKey(cacheKeyPrefixRolePermissions, roleID)

	data, err := json.Marshal(rules)
	if err != nil {
		c.logger.Error("failed to marshal permission rules",
			zap.String("role_id", roleID),
			zap.Error(err))
		return err
	}

	if err := c.redis.Set(ctx, key, data, cacheTTLRolePermissions).Err(); err != nil {
		c.logger.Error("failed to set permission rules to cache",
			zap.String("role_id", roleID),
			zap.Error(err))
		return err
	}

	c.logger.Debug("cached permission rules",
		zap.String("role_id", roleID),
		zap.Int("count", len(rules)))

	return nil
}

// GetDataPermissionFilter 获取数据权限过滤器缓存
func (c *PermissionCacheImpl) GetDataPermissionFilter(ctx context.Context, userID, resourceType string) (*DataFilterInfo, error) {
	key := c.buildCacheKey(cacheKeyPrefixDataFilter, userID, ":", resourceType)

	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		c.logger.Error("failed to get data filter from cache",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return nil, err
	}

	var filter DataFilterInfo
	if err := json.Unmarshal([]byte(val), &filter); err != nil {
		c.logger.Error("failed to unmarshal data filter",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return nil, err
	}

	c.logger.Debug("cache hit for data filter",
		zap.String("user_id", userID),
		zap.String("resource_type", resourceType))

	return &filter, nil
}

// SetDataPermissionFilter 设置数据权限过滤器缓存
func (c *PermissionCacheImpl) SetDataPermissionFilter(ctx context.Context, userID, resourceType string, filter *DataFilterInfo) error {
	key := c.buildCacheKey(cacheKeyPrefixDataFilter, userID, ":", resourceType)

	filter.CachedAt = time.Now()

	data, err := json.Marshal(filter)
	if err != nil {
		c.logger.Error("failed to marshal data filter",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return err
	}

	if err := c.redis.Set(ctx, key, data, cacheTTLDataFilter).Err(); err != nil {
		c.logger.Error("failed to set data filter to cache",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return err
	}

	c.logger.Debug("cached data filter",
		zap.String("user_id", userID),
		zap.String("resource_type", resourceType))

	return nil
}

// GetFieldPermissions 获取字段权限缓存（新版本：带 loader 的缓存模式）
func (c *PermissionCacheImpl) GetFieldPermissions(ctx context.Context, cacheKey string, loader func() (map[string]string, error)) (map[string]string, error) {
	// 尝试从缓存获取
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中，调用 loader 加载数据
			result, loadErr := loader()
			if loadErr != nil {
				return nil, loadErr
			}
			// 将加载的数据写入缓存
			data, marshalErr := json.Marshal(result)
			if marshalErr != nil {
				c.logger.Error("failed to marshal field permissions",
					zap.String("cache_key", cacheKey),
					zap.Error(marshalErr))
				return result, nil // 即使缓存失败，也返回加载的数据
			}
			if setErr := c.redis.Set(ctx, cacheKey, data, cacheTTLFieldPerms).Err(); setErr != nil {
				c.logger.Error("failed to set field permissions to cache",
					zap.String("cache_key", cacheKey),
					zap.Error(setErr))
			}
			return result, nil
		}
		c.logger.Error("failed to get field permissions from cache",
			zap.String("cache_key", cacheKey),
			zap.Error(err))
		return nil, err
	}

	// 缓存命中，反序列化数据
	var perms map[string]string
	if err := json.Unmarshal([]byte(val), &perms); err != nil {
		c.logger.Error("failed to unmarshal field permissions",
			zap.String("cache_key", cacheKey),
			zap.Error(err))
		return nil, err
	}

	c.logger.Debug("cache hit for field permissions",
		zap.String("cache_key", cacheKey),
		zap.Int("count", len(perms)))

	return perms, nil
}

// GetFieldPermissionsByUser 获取字段权限缓存（原版本：保留用于兼容）
func (c *PermissionCacheImpl) GetFieldPermissionsByUser(ctx context.Context, userID, resourceType string) (map[string]*FieldPermInfo, error) {
	key := c.buildCacheKey(cacheKeyPrefixFieldPerms, userID, ":", resourceType)

	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		c.logger.Error("failed to get field permissions from cache",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return nil, err
	}

	var perms map[string]*FieldPermInfo
	if err := json.Unmarshal([]byte(val), &perms); err != nil {
		c.logger.Error("failed to unmarshal field permissions",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return nil, err
	}

	c.logger.Debug("cache hit for field permissions",
		zap.String("user_id", userID),
		zap.String("resource_type", resourceType),
		zap.Int("count", len(perms)))

	return perms, nil
}

// SetFieldPermissions 设置字段权限缓存
func (c *PermissionCacheImpl) SetFieldPermissions(ctx context.Context, userID, resourceType string, perms map[string]*FieldPermInfo) error {
	key := c.buildCacheKey(cacheKeyPrefixFieldPerms, userID, ":", resourceType)

	data, err := json.Marshal(perms)
	if err != nil {
		c.logger.Error("failed to marshal field permissions",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return err
	}

	if err := c.redis.Set(ctx, key, data, cacheTTLFieldPerms).Err(); err != nil {
		c.logger.Error("failed to set field permissions to cache",
			zap.String("user_id", userID),
			zap.String("resource_type", resourceType),
			zap.Error(err))
		return err
	}

	c.logger.Debug("cached field permissions",
		zap.String("user_id", userID),
		zap.String("resource_type", resourceType),
		zap.Int("count", len(perms)))

	return nil
}

// InvalidateUser 失效用户相关缓存
func (c *PermissionCacheImpl) InvalidateUser(ctx context.Context, userID string) error {
	pattern := c.buildCacheKey("*", userID, "*")

	keys, err := c.redis.Keys(ctx, pattern).Result()
	if err != nil {
		c.logger.Error("failed to scan keys for user invalidation",
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := c.redis.Del(ctx, keys...).Err(); err != nil {
		c.logger.Error("failed to delete keys for user",
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}

	c.logger.Info("invalidated user cache",
		zap.String("user_id", userID),
		zap.Int("keys_deleted", len(keys)))

	return nil
}

// CheckDataPermission 检查数据权限（带缓存）
func (c *PermissionCacheImpl) CheckDataPermission(ctx context.Context, cacheKey string, loader func() (bool, error)) (bool, error) {
	// 尝试从缓存获取
	val, err := c.redis.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中，调用 loader 加载数据
			result, loadErr := loader()
			if loadErr != nil {
				return false, loadErr
			}
			// 将加载的结果写入缓存（使用较短的有效期，5分钟）
			strResult := "false"
			if result {
				strResult = "true"
			}
			data, _ := json.Marshal(strResult)
			if setErr := c.redis.Set(ctx, cacheKey, data, 5*time.Minute).Err(); setErr != nil {
				c.logger.Error("failed to set data permission to cache",
					zap.String("cache_key", cacheKey),
					zap.Error(setErr))
			}
			return result, nil
		}
		c.logger.Error("failed to get data permission from cache",
			zap.String("cache_key", cacheKey),
			zap.Error(err))
		return false, err
	}

	// 缓存命中，解析结果
	var strResult string
	if err := json.Unmarshal([]byte(val), &strResult); err != nil {
		c.logger.Error("failed to unmarshal data permission",
			zap.String("cache_key", cacheKey),
			zap.Error(err))
		return false, err
	}

	c.logger.Debug("cache hit for data permission",
		zap.String("cache_key", cacheKey),
		zap.String("result", strResult))

	return strResult == "true", nil
}

// InvalidateUserPermissions 使用户权限缓存失效
func (c *PermissionCacheImpl) InvalidateUserPermissions(ctx context.Context, userID string) error {
	// 复用 InvalidateUser 方法
	return c.InvalidateUser(ctx, userID)
}

// InvalidateDataPermission 使特定数据权限缓存失效
func (c *PermissionCacheImpl) InvalidateDataPermission(ctx context.Context, cacheKey string) error {
	if err := c.redis.Del(ctx, cacheKey).Err(); err != nil {
		c.logger.Error("failed to delete data permission cache",
			zap.String("cache_key", cacheKey),
			zap.Error(err))
		return err
	}

	c.logger.Debug("invalidated data permission cache",
		zap.String("cache_key", cacheKey))

	return nil
}

// InvalidateRole 失效角色相关缓存
func (c *PermissionCacheImpl) InvalidateRole(ctx context.Context, roleID string) error {
	// 失效角色权限缓存
	permKey := c.buildCacheKey(cacheKeyPrefixRolePermissions, roleID)
	if err := c.redis.Del(ctx, permKey).Err(); err != nil {
		c.logger.Error("failed to delete role permissions cache",
			zap.String("role_id", roleID),
			zap.Error(err))
		return err
	}

	c.logger.Info("invalidated role cache", zap.String("role_id", roleID))

	return nil
}

// InvalidateRoleUsers 失效角色的所有用户缓存
func (c *PermissionCacheImpl) InvalidateRoleUsers(ctx context.Context, roleID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}

	// 构建所有需要删除的缓存键
	keys := make([]string, 0, len(userIDs)*3) // 预估每个用户有3个缓存

	for _, userID := range userIDs {
		// 用户角色缓存
		keys = append(keys, c.buildCacheKey(cacheKeyPrefixUserRoles, userID))

		// 数据权限过滤器缓存（所有资源类型）
		filterPattern := c.buildCacheKey(cacheKeyPrefixDataFilter, userID, ":*")
		filterKeys, err := c.redis.Keys(ctx, filterPattern).Result()
		if err == nil {
			keys = append(keys, filterKeys...)
		}

		// 字段权限缓存（所有资源类型）
		fieldPattern := c.buildCacheKey(cacheKeyPrefixFieldPerms, userID, ":*")
		fieldKeys, err := c.redis.Keys(ctx, fieldPattern).Result()
		if err == nil {
			keys = append(keys, fieldKeys...)
		}
	}

	if len(keys) == 0 {
		return nil
	}

	// 批量删除
	if err := c.redis.Del(ctx, keys...).Err(); err != nil {
		c.logger.Error("failed to delete role users cache",
			zap.String("role_id", roleID),
			zap.Error(err))
		return err
	}

	c.logger.Info("invalidated role users cache",
		zap.String("role_id", roleID),
		zap.Int("users_count", len(userIDs)),
		zap.Int("keys_deleted", len(keys)))

	return nil
}

// buildCacheKey 构建缓存键
func (c *PermissionCacheImpl) buildCacheKey(parts ...string) string {
	return fmt.Sprintf("%s%s", "coze:studio:", strings.Join(parts, ""))
}

// GetCacheStats 获取缓存统计信息
// 注意：由于当前使用Redis客户端，未实现本地统计，返回空统计
func (c *PermissionCacheImpl) GetCacheStats() PermissionCacheStats {
	// TODO: 实现缓存统计功能
	// 可以在 PermissionCacheImpl 中添加统计字段，
	// 或使用 Redis 的 INFO 命令获取统计信息
	return PermissionCacheStats{}
}
