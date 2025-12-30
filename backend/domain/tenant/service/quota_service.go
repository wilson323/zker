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
)

// QuotaExceededError 配额超额错误
type QuotaExceededError struct {
	TenantID      string
	ResourceType  entity.ResourceType
	CurrentUsage  int // 当前使用量
	MaxLimit      int // 最大限制
	Required      int // 需要的额外量
	UsagePercent  float64 // 使用率百分比
}

func (e *QuotaExceededError) Error() string {
	return fmt.Sprintf("quota exceeded for tenant %s, resource %s: used %d/%d, required %d",
		e.TenantID, e.ResourceType, e.CurrentUsage, e.MaxLimit, e.Required)
}

// GetRemaining 获取剩余配额
func (e *QuotaExceededError) GetRemaining() int {
	remaining := e.MaxLimit - e.CurrentUsage
	if remaining < 0 {
		return 0
	}
	return remaining
}

// QuotaService 配额服务
type QuotaService struct {
	quotaRepo repository.QuotaRepository
}

// NewQuotaService 创建配额服务实例
func NewQuotaService(quotaRepo repository.QuotaRepository) *QuotaService {
	return &QuotaService{
		quotaRepo: quotaRepo,
	}
}

// CheckQuota 检查配额是否充足
func (s *QuotaService) CheckQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType, requiredCount int) error {
	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return fmt.Errorf("failed to get quota: %w", err)
	}

	if quota == nil {
		return fmt.Errorf("quota not found for tenant %s, resource %s", tenantID, resourceType)
	}

	// 检查是否无限制
	if quota.IsUnlimited() {
		return nil
	}

	// 检查是否超额
	if quota.UsedCount+requiredCount > quota.MaxLimit {
		return &QuotaExceededError{
			TenantID:     tenantID,
			ResourceType: resourceType,
			CurrentUsage: quota.UsedCount,
			MaxLimit:     quota.MaxLimit,
			Required:     requiredCount,
			UsagePercent: quota.GetUsagePercentage(),
		}
	}

	return nil
}

// ConsumeQuota 消费配额
// 返回值: (previousUsage, error)
func (s *QuotaService) ConsumeQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType, count int) (int, error) {
	if err := s.CheckQuota(ctx, tenantID, resourceType, count); err != nil {
		return 0, err
	}

	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return 0, err
	}

	if quota == nil {
		return 0, fmt.Errorf("quota not found for tenant %s, resource %s", tenantID, resourceType)
	}

	previousUsage := quota.UsedCount
	if err := s.quotaRepo.UpdateUsedCount(ctx, quota.QuotaID, count); err != nil {
		return previousUsage, err
	}

	return previousUsage, nil
}

// RollbackQuota 回滚配额
func (s *QuotaService) RollbackQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType, count int) error {
	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return err
	}

	if quota == nil {
		return fmt.Errorf("quota not found for tenant %s, resource %s", tenantID, resourceType)
	}

	// 回滚就是减少使用计数
	return s.quotaRepo.UpdateUsedCount(ctx, quota.QuotaID, -count)
}

// GetQuota 获取配额信息
func (s *QuotaService) GetQuota(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error) {
	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return nil, err
	}
	return quota, nil
}

// GetAllQuotas 获取租户的所有配额
func (s *QuotaService) GetAllQuotas(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	quotas, err := s.quotaRepo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return quotas, nil
}

// ResetQuota 重置配额使用量
func (s *QuotaService) ResetQuota(ctx context.Context, quotaID string) error {
	return s.quotaRepo.ResetUsage(ctx, quotaID)
}

// CheckAndResetQuotas 检查并重置需要重置的配额
func (s *QuotaService) CheckAndResetQuotas(ctx context.Context) error {
	quotas, err := s.quotaRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, quota := range quotas {
		if quota.ShouldReset() {
			if err := s.quotaRepo.ResetUsage(ctx, quota.QuotaID); err != nil {
				// 记录错误，但继续处理其他配额
				fmt.Printf("failed to reset quota %s: %v\n", quota.QuotaID, err)
			}
		}
	}

	return nil
}
