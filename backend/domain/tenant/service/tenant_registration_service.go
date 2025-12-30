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
	"math/rand"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

const (
	// 免费版默认配额
	freeTierBotsLimit        = 3    // 3个Bot
	freeTierMessagesLimit    = 1000 // 每月1000条消息
	freeTierStorageLimit     = 1024 // 1GB存储 (MB)
	freeTierTeamMembersLimit = 5    // 5个团队成员
)

var (
	idCounter int64
	idMutex   sync.Mutex
)

// RegisterTenantRequest 租户注册请求
type RegisterTenantRequest struct {
	CompanyName       string                   `json:"company_name" validate:"required,min=2,max=200"`
	Subdomain         string                   `json:"subdomain" validate:"required,min=3,max=63"`
	Email             string                   `json:"email" validate:"required,email"`
	VerificationCode  string                   `json:"verification_code" validate:"required,len=6"`
	AdminName         string                   `json:"admin_name" validate:"required,min=1,max=50"`
	AdminPassword     string                   `json:"admin_password" validate:"required,min=8,max=100"`
	TenantType        entity.TenantType        `json:"tenant_type" validate:"required"`
	BillingCycle      entity.BillingCycle      `json:"billing_cycle"`
}

// RegisterTenantResult 租户注册结果
type RegisterTenantResult struct {
	TenantID         string                    `json:"tenant_id"`
	TenantName       string                    `json:"tenant_name"`
	Subdomain        string                    `json:"subdomain"`
	SubscriptionID   string                    `json:"subscription_id"`
	AdminUserID      string                    `json:"admin_user_id"`
	PlanTier         entity.SubscriptionTier   `json:"plan_tier"`
	Status           entity.TenantStatus       `json:"status"`
	CreatedAt        int64                     `json:"created_at"`
	AccessToken      string                    `json:"access_token,omitempty"` // 管理员访问令牌
	ExpiresAt        int64                     `json:"expires_at,omitempty"`
}

// TenantRegistrationService 租户注册服务
// 负责处理租户自注册的完整流程，包括验证、事务处理、初始化
//
// 使用示例:
//
//	regSvc := service.NewTenantRegistrationService(db, tenantRepo, subscriptionRepo, validationSvc, codeSvc, emailSvc)
//	result, err := regSvc.RegisterTenant(ctx, &service.RegisterTenantRequest{...})
type TenantRegistrationService struct {
	db                 *gorm.DB
	tenantRepo         repository.TenantRepository
	subscriptionRepo   repository.SubscriptionRepository
	validationSvc      *TenantValidationService
	codeSvc            *VerificationCodeService
	emailSvc           *EmailService
}

// NewTenantRegistrationService 创建租户注册服务实例
func NewTenantRegistrationService(
	db *gorm.DB,
	tenantRepo repository.TenantRepository,
	subscriptionRepo repository.SubscriptionRepository,
	validationSvc *TenantValidationService,
	codeSvc *VerificationCodeService,
	emailSvc *EmailService,
) *TenantRegistrationService {
	return &TenantRegistrationService{
		db:               db,
		tenantRepo:       tenantRepo,
		subscriptionRepo: subscriptionRepo,
		validationSvc:    validationSvc,
		codeSvc:          codeSvc,
		emailSvc:         emailSvc,
	}
}

// RegisterTenant 注册新租户
//
// 功能：
// 1. 验证邮箱验证码
// 2. 检查企业名称和子域名可用性
// 3. 使用事务创建租户、订阅、配额
// 4. 创建管理员账户
// 5. 发送欢迎邮件
// 6. 返回注册结果
//
// 参数:
//   - ctx: 上下文
//   - req: 注册请求
//
// 返回:
//   - result: 注册结果（包含租户ID、订阅ID、管理员信息等）
//   - error: 注册失败时返回错误
func (s *TenantRegistrationService) RegisterTenant(ctx context.Context, req *RegisterTenantRequest) (*RegisterTenantResult, error) {
	// 1. 参数验证
	if err := s.validateRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	// 2. 验证邮箱验证码
	valid, err := s.codeSvc.VerifyCode(ctx, req.Email, req.VerificationCode)
	if err != nil {
		logs.Errorf("Failed to verify code for %s: %v", req.Email, err)
		return nil, fmt.Errorf("验证码验证失败: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("验证码错误或已过期")
	}

	// 3. 检查企业名称和子域名可用性
	available, _, err := s.validationSvc.CheckCompanyNameAvailability(ctx, req.CompanyName)
	if err != nil {
		return nil, fmt.Errorf("failed to check company name: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("企业名称 '%s' 已被占用", req.CompanyName)
	}

	available, _, err = s.validationSvc.CheckSubdomainAvailability(ctx, req.Subdomain)
	if err != nil {
		return nil, fmt.Errorf("failed to check subdomain: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("子域名 '%s' 已被占用", req.Subdomain)
	}

	// 4. 生成租户ID和管理员ID
	tenantID := generateTenantID()
	subscriptionID := generateSubscriptionID()
	adminUserID := generateUserID()

	// 5. 使用事务创建租户、订阅、配额
	var result *RegisterTenantResult
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 5.1 创建租户
		tenant := &entity.Tenant{
			TenantID:         tenantID,
			TenantName:       req.CompanyName,
			TenantType:       req.TenantType,
			Subdomain:        req.Subdomain,
			Status:           entity.TenantStatusActive,
			SubscriptionTier: entity.SubscriptionTierFree,
			CreatedAt:        time.Now().UnixMilli(),
			UpdatedAt:        time.Now().UnixMilli(),
		}

		if err := tx.Create(tenant).Error; err != nil {
			logs.Errorf("Failed to create tenant: %v", err)
			return fmt.Errorf("failed to create tenant: %w", err)
		}
		logs.Infof("Tenant created successfully: %s (%s)", tenant.TenantID, tenant.TenantName)

		// 5.2 创建免费订阅
		subscription := &entity.Subscription{
			SubscriptionID: subscriptionID,
			TenantID:       tenantID,
			PlanTier:       entity.SubscriptionTierFree,
			BillingCycle:   entity.BillingCycleMonthly,
			StartDate:      time.Now(),
			EndDate:        nil, // 免费版无结束日期
			AutoRenew:      true,
			Status:         entity.SubscriptionStatusActive,
			CreatedAt:      time.Now().UnixMilli(),
			UpdatedAt:      time.Now().UnixMilli(),
		}

		if err := tx.Create(subscription).Error; err != nil {
			logs.Errorf("Failed to create subscription: %v", err)
			return fmt.Errorf("failed to create subscription: %w", err)
		}
		logs.Infof("Subscription created successfully: %s for tenant %s", subscription.SubscriptionID, tenantID)

		// 5.3 创建初始配额
		quotas := s.createInitialQuotas(tenantID)
		for _, quota := range quotas {
			if err := tx.Create(&quota).Error; err != nil {
				logs.Errorf("Failed to create quota: %v", err)
				return fmt.Errorf("failed to create quota: %w", err)
			}
		}
		logs.Infof("Initial quotas created successfully for tenant %s: %d quotas", tenantID, len(quotas))

		// 5.4 创建管理员账户 (TODO: 需要User实体和Repository)
		// adminUser := &entity.User{
		// 	UserID:     adminUserID,
		// 	TenantID:   tenantID,
		// 	Username:   req.AdminName,
		// 	Email:      req.Email,
		// 	Password:   hashPassword(req.AdminPassword),
		// 	Role:       "admin",
		// 	Status:     "active",
		// 	CreatedAt:  time.Now().UnixMilli(),
		// }
		// if err := tx.Create(adminUser).Error; err != nil {
		// 	return fmt.Errorf("failed to create admin user: %w", err)
		// }

		// 6. 构建返回结果
		result = &RegisterTenantResult{
			TenantID:       tenant.TenantID,
			TenantName:     tenant.TenantName,
			Subdomain:      tenant.Subdomain,
			SubscriptionID: subscription.SubscriptionID,
			AdminUserID:    adminUserID, // TODO: 实际创建用户后返回
			PlanTier:       tenant.SubscriptionTier,
			Status:         tenant.Status,
			CreatedAt:      tenant.CreatedAt,
			// AccessToken 和 ExpiresAt 需要JWT服务生成
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 7. 发送欢迎邮件（异步，不阻塞）
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.emailSvc.SendWelcomeEmail(ctx, req.Email, req.CompanyName); err != nil {
			logs.Warnf("Failed to send welcome email to %s: %v", req.Email, err)
		}
	}()

	logs.Infof("Tenant registration completed successfully: %s (%s)", result.TenantID, result.Subdomain)
	return result, nil
}

// validateRequest 验证注册请求参数
func (s *TenantRegistrationService) validateRequest(ctx context.Context, req *RegisterTenantRequest) error {
	if req == nil {
		return fmt.Errorf("请求不能为空")
	}

	// 企业名称验证
	if req.CompanyName == "" {
		return fmt.Errorf("企业名称不能为空")
	}
	if len(req.CompanyName) < 2 || len(req.CompanyName) > 200 {
		return fmt.Errorf("企业名称长度必须在2-200个字符之间")
	}

	// 子域名验证
	if req.Subdomain == "" {
		return fmt.Errorf("子域名不能为空")
	}
	if len(req.Subdomain) < 3 || len(req.Subdomain) > 63 {
		return fmt.Errorf("子域名长度必须在3-63个字符之间")
	}

	// 邮箱验证
	if req.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	// 验证码验证
	if req.VerificationCode == "" {
		return fmt.Errorf("验证码不能为空")
	}
	if len(req.VerificationCode) != 6 {
		return fmt.Errorf("验证码格式不正确")
	}

	// 管理员信息验证
	if req.AdminName == "" {
		return fmt.Errorf("管理员姓名不能为空")
	}
	if req.AdminPassword == "" {
		return fmt.Errorf("管理员密码不能为空")
	}
	if len(req.AdminPassword) < 8 {
		return fmt.Errorf("管理员密码长度不能少于8位")
	}

	// 租户类型验证
	if req.TenantType == "" {
		return fmt.Errorf("租户类型不能为空")
	}
	if req.TenantType != entity.TenantTypeIndividual &&
	   req.TenantType != entity.TenantTypeTeam &&
	   req.TenantType != entity.TenantTypeEnterprise {
		return fmt.Errorf("无效的租户类型: %s", req.TenantType)
	}

	return nil
}

// createInitialQuotas 创建初始配额（免费版）
func (s *TenantRegistrationService) createInitialQuotas(tenantID string) []entity.Quota {
	now := time.Now().UnixMilli()

	return []entity.Quota{
		{
			QuotaID:      generateQuotaID(),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeBots,
			MaxLimit:     freeTierBotsLimit,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  now,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			QuotaID:      generateQuotaID(),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeMessages,
			MaxLimit:     freeTierMessagesLimit,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleMonthly,
			LastResetAt:  now,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			QuotaID:      generateQuotaID(),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeStorage,
			MaxLimit:     freeTierStorageLimit,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  now,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			QuotaID:      generateQuotaID(),
			TenantID:     tenantID,
			ResourceType: entity.ResourceTypeTeamMembers,
			MaxLimit:     freeTierTeamMembersLimit,
			UsedCount:    0,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  now,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}

// generateTenantID 生成租户ID
func generateTenantID() string {
	return generateID("tenant")
}

// generateSubscriptionID 生成订阅ID
func generateSubscriptionID() string {
	return generateID("sub")
}

// generateQuotaID 生成配额ID
func generateQuotaID() string {
	return generateID("quota")
}

// generateUserID 生成用户ID
func generateUserID() string {
	return generateID("user")
}

// generateID 生成通用ID (TODO: 使用UUID库)
func generateID(prefix string) string {
	idMutex.Lock()
	idCounter++
	counter := idCounter
	idMutex.Unlock()

	timestamp := time.Now().UnixNano()
	randomPart := rand.Int63n(1000000)
	return fmt.Sprintf("%s_%d_%d_%d", prefix, timestamp, counter, randomPart)
}
