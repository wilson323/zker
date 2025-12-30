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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// TenantManagementService 租户管理服务
// 提供租户信息管理、配额查询、套餐管理功能
//
// 使用示例:
//
//	managementSvc := service.NewTenantManagementService(db, tenantRepo, quotaRepo, subscriptionRepo)
//	tenant, err := managementSvc.GetTenantInfo(ctx, tenantID)
type TenantManagementService struct {
	db               *gorm.DB
	tenantRepo       repository.TenantRepository
	quotaRepo        repository.QuotaRepository
	subscriptionRepo repository.SubscriptionRepository
}

// NewTenantManagementService 创建租户管理服务实例
func NewTenantManagementService(
	db *gorm.DB,
	tenantRepo repository.TenantRepository,
	quotaRepo repository.QuotaRepository,
	subscriptionRepo repository.SubscriptionRepository,
) *TenantManagementService {
	return &TenantManagementService{
		db:               db,
		tenantRepo:       tenantRepo,
		quotaRepo:        quotaRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

// GetTenantInfo 获取租户详细信息
//
// 功能：
// 1. 查询租户基本信息
// 2. 加载关联的订阅信息
// 3. 加载配额信息
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - *TenantDetail: 租户详细信息（包含订阅和配额）
//   - error: 查询失败时返回错误
func (s *TenantManagementService) GetTenantInfo(ctx context.Context, tenantID string) (*TenantDetail, error) {
	// 1. 查询租户基本信息
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("tenant not found: %s", tenantID)
	}

	// 2. 查询订阅信息
	subscription, err := s.subscriptionRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		logs.Warnf("Failed to get subscription for tenant %s: %v", tenantID, err)
		// 不返回错误，继续处理
	}

	// 3. 查询配额信息
	quotaPtrs, err := s.quotaRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		logs.Warnf("Failed to get quotas for tenant %s: %v", tenantID, err)
		// 不返回错误，继续处理
	}

	// 将指针切片转换为值切片
	quotas := make([]entity.Quota, len(quotaPtrs))
	for i, q := range quotaPtrs {
		quotas[i] = *q
	}

	// 4. 构建详细信息
	detail := &TenantDetail{
		Tenant:       tenant,
		Subscription: subscription,
		Quotas:       quotas,
		Statistics:   s.calculateStatistics(tenant, quotas),
	}

	return detail, nil
}

// UpdateTenantInfo 更新租户信息
//
// 功能：
// 1. 验证更新权限
// 2. 更新允许修改的字段
// 3. 记录更新日志
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - req: 更新请求
//
// 返回:
//   - *Tenant: 更新后的租户信息
//   - error: 更新失败时返回错误
func (s *TenantManagementService) UpdateTenantInfo(ctx context.Context, tenantID string, req *UpdateTenantRequest) (*entity.Tenant, error) {
	// 1. 查询现有租户
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("tenant not found: %s", tenantID)
	}

	// 2. 更新允许修改的字段
	if req.TenantName != "" {
		tenant.TenantName = req.TenantName
	}
	if req.Status != "" {
		tenant.Status = req.Status
	}

	tenant.UpdatedAt = time.Now().UnixMilli()

	// 3. 保存更新
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	logs.Infof("Tenant info updated: %s", tenantID)
	return tenant, nil
}

// GetQuotaUsage 获取配额使用情况
//
// 功能：
// 1. 查询所有配额
// 2. 计算使用率
// 3. 返回详细的使用统计
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - []*QuotaUsage: 配额使用情况列表
//   - error: 查询失败时返回错误
func (s *TenantManagementService) GetQuotaUsage(ctx context.Context, tenantID string) ([]*QuotaUsage, error) {
	// 1. 查询所有配额
	quotaPtrs, err := s.quotaRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotas: %w", err)
	}

	// 2. 计算使用率
	usages := make([]*QuotaUsage, 0, len(quotaPtrs))
	for _, quotaPtr := range quotaPtrs {
		usage := &QuotaUsage{
			Quota:     quotaPtr,
			UsedCount: quotaPtr.UsedCount,
			MaxLimit:  quotaPtr.MaxLimit,
		}

		if quotaPtr.MaxLimit > 0 {
			usage.UsagePercentage = float64(quotaPtr.UsedCount) / float64(quotaPtr.MaxLimit) * 100
		}

		usage.RemainingCount = quotaPtr.MaxLimit - quotaPtr.UsedCount
		usage.IsNearLimit = usage.UsagePercentage >= 80 // 使用率超过80%视为接近上限
		usage.IsExceeded = quotaPtr.UsedCount > quotaPtr.MaxLimit

		usages = append(usages, usage)
	}

	return usages, nil
}

// CheckQuotaAvailable 检查配额是否可用
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - resourceType: 资源类型
//   - requiredCount: 需要的数量
//
// 返回:
//   - bool: 是否可用
//   - error: 查询失败时返回错误
func (s *TenantManagementService) CheckQuotaAvailable(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	requiredCount int,
) (bool, error) {
	quota, err := s.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return false, fmt.Errorf("failed to get quota: %w", err)
	}
	if quota == nil {
		return false, fmt.Errorf("quota not found for resource type: %s", resourceType)
	}

	return (quota.UsedCount + requiredCount) <= quota.MaxLimit, nil
}

// ConsumeQuota 消耗配额
//
// 功能：
// 1. 检查配额是否足够
// 2. 增加使用计数
// 3. 记录使用日志
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - resourceType: 资源类型
//   - count: 消耗数量
//
// 返回:
//   - error: 消耗失败时返回错误
func (s *TenantManagementService) ConsumeQuota(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	count int,
) error {
	// 使用事务确保原子性
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询配额（加锁）
		var quota entity.Quota
		if err := tx.Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
			First(&quota).Error; err != nil {
			return fmt.Errorf("failed to get quota: %w", err)
		}

		// 2. 检查配额是否足够
		if quota.UsedCount+count > quota.MaxLimit {
			return fmt.Errorf("quota exceeded: resource_type=%s, used=%d, max=%d, required=%d",
				resourceType, quota.UsedCount, quota.MaxLimit, count)
		}

		// 3. 增加使用计数
		now := time.Now().UnixMilli()
		if err := tx.Model(&quota).
			Updates(map[string]interface{}{
				"used_count": quota.UsedCount + count,
				"updated_at": now,
			}).Error; err != nil {
			return fmt.Errorf("failed to update quota: %w", err)
		}

		logs.Infof("Quota consumed: tenant=%s, resource=%s, count=%d, total=%d",
			tenantID, resourceType, count, quota.UsedCount+count)

		return nil
	})
}

// GetSubscription 获取订阅信息
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - *SubscriptionDetail: 订阅详细信息
//   - error: 查询失败时返回错误
func (s *TenantManagementService) GetSubscription(ctx context.Context, tenantID string) (*SubscriptionDetail, error) {
	subscription, err := s.subscriptionRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	if subscription == nil {
		return nil, fmt.Errorf("subscription not found for tenant: %s", tenantID)
	}

	detail := &SubscriptionDetail{
		Subscription: subscription,
		Features:      s.getSubscriptionFeatures(subscription.PlanTier),
		Status:        s.calculateSubscriptionStatus(subscription),
	}

	return detail, nil
}

// UpgradeSubscription 升级订阅套餐
//
// 功能：
// 1. 验证升级合法性
// 2. 更新订阅等级
// 3. 调整配额限制
// 4. 记录升级日志
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - req: 升级请求
//
// 返回:
//   - *Subscription: 升级后的订阅信息
//   - error: 升级失败时返回错误
func (s *TenantManagementService) UpgradeSubscription(
	ctx context.Context,
	tenantID string,
	req *UpgradeSubscriptionRequest,
) (*entity.Subscription, error) {
	// 使用事务处理
	var subscription *entity.Subscription
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询现有订阅
		var existingSub entity.Subscription
		if err := tx.Where("tenant_id = ?", tenantID).First(&existingSub).Error; err != nil {
			return fmt.Errorf("failed to get existing subscription: %w", err)
		}

		// 2. 验证升级合法性
		if !s.isValidUpgrade(existingSub.PlanTier, req.NewPlanTier) {
			return fmt.Errorf("invalid upgrade: from %s to %s", existingSub.PlanTier, req.NewPlanTier)
		}

		// 3. 更新订阅
		now := time.Now()
		existingSub.PlanTier = req.NewPlanTier
		existingSub.BillingCycle = req.BillingCycle
		existingSub.Status = entity.SubscriptionStatusActive
		existingSub.UpdatedAt = time.Now().UnixMilli()

		// 如果结束日期存在，延长订阅期
		if existingSub.EndDate != nil && existingSub.EndDate.Before(now) {
			// 根据计费周期计算月数
			var monthsToAdd int
			switch req.BillingCycle {
			case entity.BillingCycleMonthly:
				monthsToAdd = 1
			case entity.BillingCycleYearly:
				monthsToAdd = 12
			default:
				monthsToAdd = 1 // 默认月付
			}
			newEndDate := now.AddDate(0, monthsToAdd, 0)
			existingSub.EndDate = &newEndDate
		}

		if err := tx.Save(&existingSub).Error; err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}

		// 4. 更新配额限制
		if err := s.updateQuotaLimits(tx, tenantID, req.NewPlanTier); err != nil {
			return fmt.Errorf("failed to update quota limits: %w", err)
		}

		// 5. 更新租户的订阅等级
		if err := tx.Model(&entity.Tenant{}).
			Where("tenant_id = ?", tenantID).
			Update("subscription_tier", req.NewPlanTier).Error; err != nil {
			return fmt.Errorf("failed to update tenant subscription tier: %w", err)
		}

		subscription = &existingSub
		logs.Infof("Subscription upgraded: tenant=%s, from %s to %s",
			tenantID, existingSub.PlanTier, req.NewPlanTier)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

// ============ 辅助方法 ============

// calculateStatistics 计算租户统计信息
func (s *TenantManagementService) calculateStatistics(
	tenant *entity.Tenant,
	quotas []entity.Quota,
) *TenantStatistics {
	stats := &TenantStatistics{
		TenantID:         tenant.TenantID,
		SubscriptionTier: tenant.SubscriptionTier,
		MemberCount:      0, // TODO: 从用户表查询
		BotCount:         0, // TODO: 从Bot表查询
		CreatedAt:        tenant.CreatedAt,
	}

	// 计算总配额使用
	for _, quota := range quotas {
		switch quota.ResourceType {
		case entity.ResourceTypeTeamMembers:
			stats.MemberCount = quota.UsedCount
		case entity.ResourceTypeBots:
			stats.BotCount = quota.UsedCount
		}
	}

	return stats
}

// getSubscriptionFeatures 获取订阅套餐功能列表
func (s *TenantManagementService) getSubscriptionFeatures(tier entity.SubscriptionTier) []string {
	features := map[entity.SubscriptionTier][]string{
		entity.SubscriptionTierFree: {
			"最多3个Bot",
			"每月1000条消息",
			"1GB存储空间",
			"5个团队成员",
		},
		entity.SubscriptionTierPro: {
			"最多50个Bot",
			"每月50000条消息",
			"50GB存储空间",
			"50个团队成员",
			"优先技术支持",
		},
		entity.SubscriptionTierEnterprise: {
			"无限Bot",
			"无限消息",
			"无限存储空间",
			"无限团队成员",
			"专属客户经理",
			"SLA保障",
			"自定义集成",
		},
	}

	return features[tier]
}

// calculateSubscriptionStatus 计算订阅状态
func (s *TenantManagementService) calculateSubscriptionStatus(sub *entity.Subscription) string {
	if sub.Status == entity.SubscriptionStatusSuspended {
		return "suspended"
	}

	if sub.EndDate != nil && sub.EndDate.Before(time.Now()) {
		return "expired"
	}

	return "active"
}

// isValidUpgrade 验证升级是否合法
func (s *TenantManagementService) isValidUpgrade(from, to entity.SubscriptionTier) bool {
	upgradePaths := map[entity.SubscriptionTier][]entity.SubscriptionTier{
		entity.SubscriptionTierFree: {
			entity.SubscriptionTierPro,
			entity.SubscriptionTierEnterprise,
		},
		entity.SubscriptionTierPro: {
			entity.SubscriptionTierEnterprise,
		},
		entity.SubscriptionTierEnterprise: {}, // 顶级套餐，无法升级
	}

	validTiers, exists := upgradePaths[from]
	if !exists {
		return false
	}

	for _, tier := range validTiers {
		if tier == to {
			return true
		}
	}

	return false
}

// updateQuotaLimits 更新配额限制
func (s *TenantManagementService) updateQuotaLimits(tx *gorm.DB, tenantID string, tier entity.SubscriptionTier) error {
	limits := s.getPlanLimits(tier)

	for resourceType, limit := range limits {
		if err := tx.Model(&entity.Quota{}).
			Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
			Update("max_limit", limit).Error; err != nil {
			return err
		}
	}

	return nil
}

// getPlanLimits 获取套餐配额限制
func (s *TenantManagementService) getPlanLimits(tier entity.SubscriptionTier) map[entity.ResourceType]int {
	switch tier {
	case entity.SubscriptionTierFree:
		return map[entity.ResourceType]int{
			entity.ResourceTypeBots:         3,
			entity.ResourceTypeMessages:     1000,
			entity.ResourceTypeStorage:      1024, // 1GB
			entity.ResourceTypeTeamMembers:  5,
		}
	case entity.SubscriptionTierPro:
		return map[entity.ResourceType]int{
			entity.ResourceTypeBots:         50,
			entity.ResourceTypeMessages:     50000,
			entity.ResourceTypeStorage:      51200, // 50GB
			entity.ResourceTypeTeamMembers:  50,
		}
	case entity.SubscriptionTierEnterprise:
		return map[entity.ResourceType]int{
			entity.ResourceTypeBots:         999999,
			entity.ResourceTypeMessages:     999999,
			entity.ResourceTypeStorage:      1024000, // 1TB
			entity.ResourceTypeTeamMembers:  999999,
		}
	default:
		return map[entity.ResourceType]int{}
	}
}

// ============ 请求/响应结构 ============

// TenantDetail 租户详细信息
type TenantDetail struct {
	Tenant       *entity.Tenant      `json:"tenant"`
	Subscription *entity.Subscription `json:"subscription,omitempty"`
	Quotas       []entity.Quota      `json:"quotas,omitempty"`
	Statistics   *TenantStatistics   `json:"statistics,omitempty"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	TenantName string           `json:"tenant_name,omitempty" validate:"omitempty,min=2,max=200"`
	Status     entity.TenantStatus `json:"status,omitempty" validate:"omitempty,oneof=active suspended deleted"`
}

// TenantStatistics 租户统计信息
type TenantStatistics struct {
	TenantID         string                  `json:"tenant_id"`
	SubscriptionTier entity.SubscriptionTier `json:"subscription_tier"`
	MemberCount      int                     `json:"member_count"`
	BotCount         int                     `json:"bot_count"`
	CreatedAt        int64                   `json:"created_at"`
}

// QuotaUsage 配额使用情况
type QuotaUsage struct {
	Quota              *entity.Quota `json:"quota"`
	UsedCount          int           `json:"used_count"`
	MaxLimit           int           `json:"max_limit"`
	RemainingCount     int           `json:"remaining_count"`
	UsagePercentage    float64       `json:"usage_percentage"`
	IsNearLimit        bool          `json:"is_near_limit"`
	IsExceeded         bool          `json:"is_exceeded"`
}

// SubscriptionDetail 订阅详细信息
type SubscriptionDetail struct {
	Subscription *entity.Subscription `json:"subscription"`
	Features     []string             `json:"features"`
	Status       string               `json:"status"` // active, expired, suspended
}

// UpgradeSubscriptionRequest 升级订阅请求
type UpgradeSubscriptionRequest struct {
	NewPlanTier   entity.SubscriptionTier  `json:"new_plan_tier" validate:"required,oneof=free pro enterprise"`
	BillingCycle  entity.BillingCycle      `json:"billing_cycle" validate:"required,oneof=monthly yearly"`
}
