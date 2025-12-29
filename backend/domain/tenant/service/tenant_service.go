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

	"github.com/google/uuid"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// TenantService 租户管理服务
type TenantService struct {
	tenantRepo    repository.TenantRepository
	subService    *SubscriptionService
	quotaService  *QuotaService
}

// NewTenantService 创建租户管理服务实例
func NewTenantService(
	tenantRepo repository.TenantRepository,
	subRepo repository.SubscriptionRepository,
	quotaRepo repository.QuotaRepository,
) *TenantService {
	subService := NewSubscriptionService(subRepo, quotaRepo)
	quotaService := NewQuotaService(quotaRepo)

	return &TenantService{
		tenantRepo:   tenantRepo,
		subService:   subService,
		quotaService: quotaService,
	}
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	TenantName string            `json:"tenant_name" validate:"required,min=2,max=200"`
	TenantType entity.TenantType `json:"tenant_type" validate:"required"`
	PlanTier   entity.SubscriptionTier `json:"plan_tier,omitempty"`
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*entity.Tenant, error) {
	// 1. 验证租户名称是否已存在
	existing, err := s.tenantRepo.GetByName(ctx, req.TenantName)
	if err != nil {
		return nil, fmt.Errorf("failed to check tenant name: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("tenant name %s already exists", req.TenantName)
	}

	// 2. 确定订阅等级
	planTier := req.PlanTier
	if planTier == "" {
		planTier = entity.SubscriptionTierFree // 默认免费版
	}

	// 3. 创建租户
	tenant := &entity.Tenant{
		TenantID:         generateUUID(),
		TenantName:       req.TenantName,
		TenantType:       req.TenantType,
		Status:           entity.TenantStatusActive,
		SubscriptionTier: planTier,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// 4. 创建订阅
	billingCycle := entity.BillingCycleMonthly // 默认月付
	_, err = s.subService.CreateSubscription(ctx, tenant.TenantID, planTier, billingCycle)
	if err != nil {
		// 回滚租户创建
		_ = s.tenantRepo.Delete(ctx, tenant.TenantID)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return tenant, nil
}

// GetTenant 获取租户信息
func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("tenant not found: %s", tenantID)
	}
	return tenant, nil
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	TenantID   string `json:"tenant_id" validate:"required"`
	TenantName string `json:"tenant_name,omitempty" validate:"omitempty,min=2,max=200"`
}

// UpdateTenant 更新租户信息
func (s *TenantService) UpdateTenant(ctx context.Context, req *UpdateTenantRequest) (*entity.Tenant, error) {
	// 1. 获取现有租户
	tenant, err := s.tenantRepo.GetByID(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return nil, fmt.Errorf("tenant not found: %s", req.TenantID)
	}

	// 2. 如果更新名称，检查新名称是否已被占用
	if req.TenantName != "" && req.TenantName != tenant.TenantName {
		existing, err := s.tenantRepo.GetByName(ctx, req.TenantName)
		if err != nil {
			return nil, fmt.Errorf("failed to check tenant name: %w", err)
		}
		if existing != nil {
			return nil, fmt.Errorf("tenant name %s already exists", req.TenantName)
		}
		tenant.TenantName = req.TenantName
	}

	// 3. 更新租户
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}

// DeleteTenant 删除租户（软删除）
func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	// 1. 检查租户是否存在
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}
	if tenant == nil {
		return fmt.Errorf("tenant not found: %s", tenantID)
	}

	// 2. 软删除租户
	if err := s.tenantRepo.Delete(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	return nil
}

// SuspendTenant 暂停租户
func (s *TenantService) SuspendTenant(ctx context.Context, tenantID string) error {
	if err := s.tenantRepo.UpdateStatus(ctx, tenantID, entity.TenantStatusSuspended); err != nil {
		return fmt.Errorf("failed to suspend tenant: %w", err)
	}
	return nil
}

// ActivateTenant 激活租户
func (s *TenantService) ActivateTenant(ctx context.Context, tenantID string) error {
	if err := s.tenantRepo.UpdateStatus(ctx, tenantID, entity.TenantStatusActive); err != nil {
		return fmt.Errorf("failed to activate tenant: %w", err)
	}
	return nil
}

// ListTenantsRequest 查询租户列表请求
type ListTenantsRequest struct {
	Status           *entity.TenantStatus     `json:"status,omitempty"`
	TenantType       *entity.TenantType       `json:"tenant_type,omitempty"`
	SubscriptionTier *entity.SubscriptionTier `json:"subscription_tier,omitempty"`
	PageToken         string                  `json:"page_token,omitempty"`
	PageSize          int                     `json:"page_size,omitempty"`
}

// ListTenantsResponse 查询租户列表响应
type ListTenantsResponse struct {
	Tenants    []*entity.Tenant `json:"tenants"`
	TotalCount int64            `json:"total_count"`
	NextToken  string           `json:"next_token,omitempty"`
}

// ListTenants 查询租户列表
func (s *TenantService) ListTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error) {
	// 构建过滤器
	filter := &repository.TenantFilter{
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	if req.Status != nil {
		filter.Status = *req.Status
	}
	if req.TenantType != nil {
		filter.TenantType = *req.TenantType
	}
	if req.SubscriptionTier != nil {
		filter.SubscriptionTier = *req.SubscriptionTier
	}

	// 查询数据
	tenants, total, err := s.tenantRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return &ListTenantsResponse{
		Tenants:    tenants,
		TotalCount: total,
	}, nil
}

// UpgradeTenantTier 升级租户订阅等级
func (s *TenantService) UpgradeTenantTier(ctx context.Context, tenantID string, newTier entity.SubscriptionTier) (*entity.Tenant, error) {
	// 1. 升级订阅
	sub, err := s.subService.UpgradeTier(ctx, tenantID, newTier)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade subscription: %w", err)
	}

	// 2. 更新租户的订阅等级
	if err := s.tenantRepo.UpdateSubscriptionTier(ctx, tenantID, newTier); err != nil {
		return nil, fmt.Errorf("failed to update tenant subscription tier: %w", err)
	}

	// 3. 返回更新后的租户信息
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated tenant: %w", err)
	}

	tenant.Subscription = sub
	return tenant, nil
}

// GetTenantQuotas 获取租户的所有配额
func (s *TenantService) GetTenantQuotas(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	quotas, err := s.quotaService.GetAllQuotas(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotas: %w", err)
	}
	return quotas, nil
}

// CheckTenantQuota 检查租户配额
func (s *TenantService) CheckTenantQuota(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	requiredCount int,
) error {
	return s.quotaService.CheckQuota(ctx, tenantID, resourceType, requiredCount)
}

// ConsumeTenantQuota 消费租户配额
func (s *TenantService) ConsumeTenantQuota(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	count int,
) error {
	return s.quotaService.ConsumeQuota(ctx, tenantID, resourceType, count)
}

// RollbackTenantQuota 回滚租户配额
func (s *TenantService) RollbackTenantQuota(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	count int,
) error {
	return s.quotaService.RollbackQuota(ctx, tenantID, resourceType, count)
}

// GetTenantSubscription 获取租户订阅信息
func (s *TenantService) GetTenantSubscription(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	sub, err := s.subService.GetSubscription(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return sub, nil
}
