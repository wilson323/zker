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
	"time"

	"github.com/google/uuid"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// SubscriptionService 订阅管理服务
type SubscriptionService struct {
	subRepo   repository.SubscriptionRepository
	quotaRepo repository.QuotaRepository
}

// NewSubscriptionService 创建订阅管理服务实例
func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	quotaRepo repository.QuotaRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:   subRepo,
		quotaRepo: quotaRepo,
	}
}

// CreateSubscription 创建订阅
func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	tenantID string,
	planTier entity.SubscriptionTier,
	billingCycle entity.BillingCycle,
) (*entity.Subscription, error) {
	// 1. 计算订阅结束日期
	startDate := time.Now()
	var endDate *time.Time

	if billingCycle == entity.BillingCycleMonthly {
		ed := startDate.AddDate(0, 1, 0)
		endDate = &ed
	} else if billingCycle == entity.BillingCycleYearly {
		ed := startDate.AddDate(1, 0, 0)
		endDate = &ed
	}

	// 2. 创建订阅
	sub := &entity.Subscription{
		SubscriptionID: generateUUID(),
		TenantID:       tenantID,
		PlanTier:       planTier,
		BillingCycle:   billingCycle,
		StartDate:      startDate,
		EndDate:        endDate,
		AutoRenew:      true,
		Status:         entity.SubscriptionStatusActive,
	}

	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// 3. 创建配额
	quotas := s.getQuotasForTier(planTier)
	for _, quota := range quotas {
		quota.TenantID = tenantID
		if err := s.quotaRepo.Create(ctx, quota); err != nil {
			return nil, fmt.Errorf("failed to create quota: %w", err)
		}
	}

	return sub, nil
}

// getQuotasForTier 根据订阅等级获取配额配置
func (s *SubscriptionService) getQuotasForTier(planTier entity.SubscriptionTier) []*entity.Quota {
	quotas := make([]*entity.Quota, 0)

	switch planTier {
	case entity.SubscriptionTierFree:
		// 免费版：10个Bot，1000条消息/月
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeBots,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeMessages,
			MaxLimit:     1000,
			ResetCycle:   entity.ResetCycleMonthly,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeStorage,
			MaxLimit:     1024, // 1GB
			ResetCycle:   entity.ResetCycleNever,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeTeamMembers,
			MaxLimit:     1,
			ResetCycle:   entity.ResetCycleNever,
		})

	case entity.SubscriptionTierPro:
		// 专业版：100个Bot，100000条消息/月
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeBots,
			MaxLimit:     100,
			ResetCycle:   entity.ResetCycleNever,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeMessages,
			MaxLimit:     100000,
			ResetCycle:   entity.ResetCycleMonthly,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeStorage,
			MaxLimit:     10240, // 10GB
			ResetCycle:   entity.ResetCycleNever,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeTeamMembers,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
		})

	case entity.SubscriptionTierEnterprise:
		// 企业版：无限制
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeBots,
			MaxLimit:     -1, // 无限制
			ResetCycle:   entity.ResetCycleNever,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeMessages,
			MaxLimit:     -1, // 无限制
			ResetCycle:   entity.ResetCycleMonthly,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeStorage,
			MaxLimit:     -1, // 无限制
			ResetCycle:   entity.ResetCycleNever,
		})
		quotas = append(quotas, &entity.Quota{
			QuotaID:      generateUUID(),
			ResourceType: entity.ResourceTypeTeamMembers,
			MaxLimit:     -1, // 无限制
			ResetCycle:   entity.ResetCycleNever,
		})
	}

	return quotas
}

// UpdateSubscription 更新订阅
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, sub *entity.Subscription) error {
	return s.subRepo.Update(ctx, sub)
}

// UpgradeTier 升级订阅等级
func (s *SubscriptionService) UpgradeTier(
	ctx context.Context,
	tenantID string,
	newTier entity.SubscriptionTier,
) (*entity.Subscription, error) {
	// 1. 获取当前订阅
	sub, err := s.subRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub == nil {
		return nil, fmt.Errorf("subscription not found for tenant %s", tenantID)
	}

	// 2. 检查是否已经是该等级
	if sub.PlanTier == newTier {
		return sub, nil
	}

	// 3. 只能升级，不能降级（除非取消订阅）
	if !s.isUpgrade(sub.PlanTier, newTier) {
		return nil, fmt.Errorf("cannot downgrade from %s to %s", sub.PlanTier, newTier)
	}

	// 4. 更新订阅等级
	sub.PlanTier = newTier
	sub.UpdatedAt = time.Now().UnixMilli()

	if err := s.subRepo.Update(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// 5. 更新配额
	quotas := s.getQuotasForTier(newTier)
	for _, newQuota := range quotas {
		// 检查是否已存在该资源类型的配额
		existingQuota, _ := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, newQuota.ResourceType)

		if existingQuota != nil {
			// 更新现有配额的限制
			existingQuota.MaxLimit = newQuota.MaxLimit
			existingQuota.ResetCycle = newQuota.ResetCycle
			if err := s.quotaRepo.Update(ctx, existingQuota); err != nil {
				return nil, fmt.Errorf("failed to update quota: %w", err)
			}
		} else {
			// 创建新配额
			newQuota.TenantID = tenantID
			if err := s.quotaRepo.Create(ctx, newQuota); err != nil {
				return nil, fmt.Errorf("failed to create quota: %w", err)
			}
		}
	}

	return sub, nil
}

// isUpgrade 检查是否是升级
func (s *SubscriptionService) isUpgrade(currentTier, newTier entity.SubscriptionTier) bool {
	tierOrder := map[entity.SubscriptionTier]int{
		entity.SubscriptionTierFree:       1,
		entity.SubscriptionTierPro:        2,
		entity.SubscriptionTierEnterprise: 3,
	}

	return tierOrder[newTier] >= tierOrder[currentTier]
}

// CancelSubscription 取消订阅
func (s *SubscriptionService) CancelSubscription(ctx context.Context, tenantID string) error {
	sub, err := s.subRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub == nil {
		return fmt.Errorf("subscription not found for tenant %s", tenantID)
	}

	// 更新状态为已取消
	sub.Status = entity.SubscriptionStatusCancelled
	sub.AutoRenew = false
	sub.UpdatedAt = time.Now().UnixMilli()

	if err := s.subRepo.Update(ctx, sub); err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return nil
}

// CheckSubscriptionStatus 检查订阅状态（定时任务调用）
func (s *SubscriptionService) CheckSubscriptionStatus(ctx context.Context) error {
	// 获取所有激活的订阅
	subs, _, err := s.subRepo.List(ctx, &repository.SubscriptionFilter{
		Status: entity.SubscriptionStatusActive,
	})
	if err != nil {
		return err
	}

	now := time.Now()
	for _, sub := range subs {
		// 检查是否过期
		if sub.EndDate != nil && sub.EndDate.Before(now) {
			// 更新状态为过期
			if err := s.subRepo.UpdateStatus(ctx, sub.SubscriptionID, entity.SubscriptionStatusExpired); err != nil {
				fmt.Printf("failed to update subscription status for %s: %v\n", sub.SubscriptionID, err)
			}
		}
	}

	return nil
}

// GetSubscription 获取订阅信息
func (s *SubscriptionService) GetSubscription(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	sub, err := s.subRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// generateUUID 生成UUID
func generateUUID() string {
	return uuid.New().String()
}
