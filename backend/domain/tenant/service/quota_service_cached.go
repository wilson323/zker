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
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
)

// CachedQuotaService 带缓存的配额服务
type CachedQuotaService struct {
	*QuotaService        // 嵌入原始配额服务
	quotaCache    cache.QuotaCache // 配额缓存
}

// NewCachedQuotaService 创建带缓存的配额服务实例
func NewCachedQuotaService(
	quotaRepo repository.QuotaRepository,
	quotaCache cache.QuotaCache,
) *CachedQuotaService {
	baseService := &QuotaService{
		quotaRepo: quotaRepo,
	}

	return &CachedQuotaService{
		QuotaService: baseService,
		quotaCache:   quotaCache,
	}
}

// CheckQuota 检查配额是否充足（带缓存）
// 性能优化：
//   - L1: 本地缓存（1分钟）- 响应时间 < 1ms
//   - L2: Redis缓存（5分钟）- 响应时间 < 5ms
//   - L3: 数据库查询 - 响应时间 15ms
// 预期缓存命中率：95%+
func (s *CachedQuotaService) CheckQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType, requiredCount int) error {
	// 1. 从缓存获取配额信息
	quotaInfo, err := s.quotaCache.GetQuota(ctx, tenantID, string(resourceType), func() (*cache.QuotaInfo, error) {
		// 缓存未命中时从数据库查询
		quota, err := s.QuotaService.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
		if err != nil {
			return nil, err
		}

		if quota == nil {
			return nil, fmt.Errorf("quota not found for tenant %s, resource %s", tenantID, resourceType)
		}

		// 转换为缓存对象
		return &cache.QuotaInfo{
			TenantID:     quota.TenantID,
			ResourceType: string(quota.ResourceType),
			MaxLimit:     quota.MaxLimit,
			UsedCount:    quota.UsedCount,
			IsUnlimited:  quota.MaxLimit == -1,
		}, nil
	})

	if err != nil {
		return fmt.Errorf("failed to get quota: %w", err)
	}

	// 2. 检查配额是否充足
	if !quotaInfo.HasRemaining(requiredCount) {
		return &QuotaExceededError{
			TenantID:     tenantID,
			ResourceType: resourceType,
			CurrentUsage: quotaInfo.UsedCount,
			MaxLimit:     quotaInfo.MaxLimit,
			Required:     requiredCount,
			UsagePercent: quotaInfo.GetUsagePercentage(),
		}
	}

	return nil
}

// ConsumeQuota 消费配额（带缓存自动失效）
// 关键：先更新数据库，再删除缓存（保证一致性）
func (s *CachedQuotaService) ConsumeQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType, count int) (int, error) {
	// 1. 先检查配额是否充足
	if err := s.CheckQuota(ctx, tenantID, resourceType, count); err != nil {
		return 0, err
	}

	// 2. 查询当前配额（获取previous usage）
	quota, err := s.QuotaService.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return 0, err
	}

	if quota == nil {
		return 0, fmt.Errorf("quota not found for tenant %s, resource %s", tenantID, resourceType)
	}

	previousUsage := quota.UsedCount

	// 3. 使用缓存的ConsumeQuota方法（会自动失效缓存）
	err = s.quotaCache.ConsumeQuota(ctx, tenantID, string(resourceType), count, func(actualCount int) error {
		return s.QuotaService.quotaRepo.UpdateUsedCount(ctx, quota.QuotaID, actualCount)
	})

	if err != nil {
		return previousUsage, err
	}

	return previousUsage, nil
}

// RollbackQuota 回滚配额（带缓存自动失效）
func (s *CachedQuotaService) RollbackQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType, count int) error {
	// 回滚就是消费负数配额
	_, err := s.ConsumeQuota(ctx, tenantID, resourceType, -count)
	return err
}

// GetQuota 获取配额信息（带缓存）
func (s *CachedQuotaService) GetQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error) {
	quotaInfo, err := s.quotaCache.GetQuota(ctx, tenantID, string(resourceType), func() (*cache.QuotaInfo, error) {
		quota, err := s.QuotaService.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
		if err != nil {
			return nil, err
		}

		if quota == nil {
			return nil, fmt.Errorf("quota not found")
		}

		return &cache.QuotaInfo{
			TenantID:     quota.TenantID,
			ResourceType: string(quota.ResourceType),
			MaxLimit:     quota.MaxLimit,
			UsedCount:    quota.UsedCount,
			IsUnlimited:  quota.MaxLimit == -1,
		}, nil
	})

	if err != nil {
		return nil, err
	}

	// 转换回entity.Quota
	return &entity.Quota{
		TenantID:     quotaInfo.TenantID,
		ResourceType: entity.ResourceType(quotaInfo.ResourceType),
		MaxLimit:     quotaInfo.MaxLimit,
		UsedCount:    quotaInfo.UsedCount,
	}, nil
}

// InvalidateQuota 使配额缓存失效
// 调用时机：
//   1. 管理员手动调整配额时
//   2. 订阅计划变更时
//   3. 配额重置时
func (s *CachedQuotaService) InvalidateQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType) error {
	return s.quotaCache.InvalidateQuota(ctx, tenantID, string(resourceType))
}

// GetCacheStats 获取缓存统计信息
func (s *CachedQuotaService) GetCacheStats() cache.QuotaCacheStats {
	return s.quotaCache.GetCacheStats()
}
