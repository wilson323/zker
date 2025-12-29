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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// RedisClient Redis客户端接口
type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

// TenantServiceCached 带缓存的租户服务
type TenantServiceCached struct {
	tenantRepo    repository.TenantRepository
	subscription  repository.SubscriptionRepository
	redis         RedisClient
	cacheDuration time.Duration
}

// NewTenantServiceCached 创建带缓存的租户服务实例
func NewTenantServiceCached(
	tenantRepo repository.TenantRepository,
	subscription repository.SubscriptionRepository,
	redis RedisClient,
	cacheDuration time.Duration,
) *TenantServiceCached {
	if cacheDuration == 0 {
		cacheDuration = 5 * time.Minute // 默认5分钟
	}

	return &TenantServiceCached{
		tenantRepo:    tenantRepo,
		subscription:  subscription,
		redis:         redis,
		cacheDuration: cacheDuration,
	}
}

// GetByID 根据ID获取租户（带缓存）
func (s *TenantServiceCached) GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	// 1. 先查缓存
	cacheKey := fmt.Sprintf("tenant:%s", tenantID)
	cached, err := s.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var tenant entity.Tenant
		if err := json.Unmarshal([]byte(cached), &tenant); err == nil {
			return &tenant, nil
		}
	}

	// 2. 缓存未命中，查数据库
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存
	if tenant != nil {
		data, _ := json.Marshal(tenant)
		_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)
	}

	return tenant, nil
}

// CreateTenant 创建租户（带缓存失效）
func (s *TenantServiceCached) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*entity.Tenant, error) {
	tenant := &entity.Tenant{
		TenantID:   generateTenantID(),
		TenantName: req.TenantName,
		TenantType: req.TenantType,
		Status:     entity.TenantStatusActive,
	}

	// 1. 创建租户
	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, err
	}

	// 2. 删除可能的缓存（防止脏读）
	cacheKey := fmt.Sprintf("tenant:%s", tenant.TenantID)
	_ = s.redis.Del(ctx, cacheKey)

	// 3. 预热缓存
	data, _ := json.Marshal(tenant)
	_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)

	return tenant, nil
}

// Update 更新租户（带缓存更新）
func (s *TenantServiceCached) Update(ctx context.Context, tenant *entity.Tenant) error {
	// 1. 更新数据库
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return err
	}

	// 2. 更新缓存
	cacheKey := fmt.Sprintf("tenant:%s", tenant.TenantID)
	data, _ := json.Marshal(tenant)
	_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)

	return nil
}

// Delete 删除租户（带缓存失效）
func (s *TenantServiceCached) Delete(ctx context.Context, tenantID string) error {
	// 1. 删除数据库记录
	if err := s.tenantRepo.Delete(ctx, tenantID); err != nil {
		return err
	}

	// 2. 删除缓存
	cacheKey := fmt.Sprintf("tenant:%s", tenantID)
	_ = s.redis.Del(ctx, cacheKey)

	return nil
}

// GetSubscription 获取订阅信息（带缓存）
func (s *TenantServiceCached) GetSubscription(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	// 1. 先查缓存
	cacheKey := fmt.Sprintf("subscription:%s", tenantID)
	cached, err := s.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var subscription entity.Subscription
		if err := json.Unmarshal([]byte(cached), &subscription); err == nil {
			return &subscription, nil
		}
	}

	// 2. 缓存未命中，查数据库
	subscription, err := s.subscription.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 3. 写入缓存
	if subscription != nil {
		data, _ := json.Marshal(subscription)
		_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)
	}

	return subscription, nil
}

// InvalidateTenantCache 失效租户相关缓存
func (s *TenantServiceCached) InvalidateTenantCache(ctx context.Context, tenantID string) error {
	// 删除租户缓存
	tenantKey := fmt.Sprintf("tenant:%s", tenantID)
	// 删除订阅缓存
	subKey := fmt.Sprintf("subscription:%s", tenantID)
	// 删除配额缓存
	quotaKey := fmt.Sprintf("quota:%s:*", tenantID)

	return s.redis.Del(ctx, tenantKey, subKey, quotaKey)
}

// BatchGetTenants 批量获取租户（缓存优化）
func (s *TenantServiceCached) BatchGetTenants(ctx context.Context, tenantIDs []string) ([]*entity.Tenant, error) {
	tenants := make([]*entity.Tenant, 0, len(tenantIDs))
	missedIDs := make([]string, 0)

	// 1. 批量查缓存
	for _, tenantID := range tenantIDs {
		cacheKey := fmt.Sprintf("tenant:%s", tenantID)
		cached, err := s.redis.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var tenant entity.Tenant
			if err := json.Unmarshal([]byte(cached), &tenant); err == nil {
				tenants = append(tenants, &tenant)
				continue
			}
		}

		// 缓存未命中，记录ID
		missedIDs = append(missedIDs, tenantID)
	}

	// 2. 批量查询未命中的租户
	if len(missedIDs) > 0 {
		for _, tenantID := range missedIDs {
			tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
			if err == nil && tenant != nil {
				tenants = append(tenants, tenant)

				// 预热缓存
				cacheKey := fmt.Sprintf("tenant:%s", tenantID)
				data, _ := json.Marshal(tenant)
				_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)
			}
		}
	}

	return tenants, nil
}

// WarmupCache 预热缓存（用于系统启动时）
func (s *TenantServiceCached) WarmupCache(ctx context.Context, tenantIDs []string) error {
	for _, tenantID := range tenantIDs {
		// 查询并缓存租户信息
		tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
		if err == nil && tenant != nil {
			cacheKey := fmt.Sprintf("tenant:%s", tenantID)
			data, _ := json.Marshal(tenant)
			_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)
		}

		// 查询并缓存订阅信息
		subscription, err := s.subscription.GetByTenant(ctx, tenantID)
		if err == nil && subscription != nil {
			cacheKey := fmt.Sprintf("subscription:%s", tenantID)
			data, _ := json.Marshal(subscription)
			_ = s.redis.Set(ctx, cacheKey, data, s.cacheDuration)
		}
	}

	return nil
}

// GetCacheStats 获取缓存统计信息
func (s *TenantServiceCached) GetCacheStats(ctx context.Context) (*CacheStats, error) {
	// 这里需要实现实际的统计逻辑
	// 可能需要维护hit/miss计数器

	return &CacheStats{
		TotalTenants:    0,
		CachedTenants:   0,
		HitRate:         0.0,
		AvgResponseTime: 0,
	}, nil
}

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalTenants    int64
	CachedTenants   int64
	HitRate         float64
	AvgResponseTime time.Duration
}

// generateTenantID 生成租户ID
func generateTenantID() string {
	return fmt.Sprintf("tenant_%d", time.Now().UnixNano())
}
