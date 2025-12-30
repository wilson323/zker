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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
)

// MockSubscriptionRepository 订阅仓储 Mock
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(ctx context.Context, sub *entity.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetByID(ctx context.Context, subscriptionID string) (*entity.Subscription, error) {
	args := m.Called(ctx, subscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetByTenantID(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetByTenant(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetAll(ctx context.Context) ([]*entity.Subscription, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(ctx context.Context, sub *entity.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) List(ctx context.Context, filter *repository.SubscriptionFilter) ([]*entity.Subscription, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Subscription), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubscriptionRepository) UpdateStatus(ctx context.Context, subscriptionID string, status entity.SubscriptionStatus) error {
	args := m.Called(ctx, subscriptionID, status)
	return args.Error(0)
}

// SubscriptionServiceTestSuite 测试套件
type SubscriptionServiceTestSuite struct {
	suite.Suite
	service       *SubscriptionService
	mockSubRepo   *MockSubscriptionRepository
	mockQuotaRepo *MockQuotaRepository
	ctx           context.Context
}

func (s *SubscriptionServiceTestSuite) SetupTest() {
	s.mockSubRepo = new(MockSubscriptionRepository)
	s.mockQuotaRepo = new(MockQuotaRepository)
	s.service = NewSubscriptionService(s.mockSubRepo, s.mockQuotaRepo)
	s.ctx = context.Background()
}

// TestCreateSubscription_FreeTier 测试创建免费版订阅
func (s *SubscriptionServiceTestSuite) TestCreateSubscription_FreeTier() {
	// Arrange
	tenantID := "tenant-123"
	planTier := entity.SubscriptionTierFree
	billingCycle := entity.BillingCycleMonthly

	s.mockSubRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(nil)
	s.mockQuotaRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Quota")).Return(nil).Times(4)

	// Act
	sub, err := s.service.CreateSubscription(s.ctx, tenantID, planTier, billingCycle)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), sub)
	assert.Equal(s.T(), tenantID, sub.TenantID)
	assert.Equal(s.T(), planTier, sub.PlanTier)
	assert.Equal(s.T(), billingCycle, sub.BillingCycle)
	assert.Equal(s.T(), entity.SubscriptionStatusActive, sub.Status)
	assert.True(s.T(), sub.AutoRenew)
	assert.NotNil(s.T(), sub.EndDate)
	s.mockSubRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestCreateSubscription_ProTier 测试创建专业版订阅
func (s *SubscriptionServiceTestSuite) TestCreateSubscription_ProTier() {
	// Arrange
	tenantID := "tenant-123"
	planTier := entity.SubscriptionTierPro
	billingCycle := entity.BillingCycleYearly

	s.mockSubRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(nil)
	s.mockQuotaRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Quota")).Return(nil).Times(4)

	// Act
	sub, err := s.service.CreateSubscription(s.ctx, tenantID, planTier, billingCycle)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), sub)
	assert.Equal(s.T(), planTier, sub.PlanTier)
	assert.Equal(s.T(), billingCycle, sub.BillingCycle)
	s.mockSubRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestCreateSubscription_EnterpriseTier 测试创建企业版订阅
func (s *SubscriptionServiceTestSuite) TestCreateSubscription_EnterpriseTier() {
	// Arrange
	tenantID := "tenant-123"
	planTier := entity.SubscriptionTierEnterprise
	billingCycle := entity.BillingCycleYearly

	s.mockSubRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(nil)
	s.mockQuotaRepo.On("Create", s.ctx, mock.MatchedBy(func(q *entity.Quota) bool {
		return q.MaxLimit == -1 // 企业版应该是无限制
	})).Return(nil).Times(4)

	// Act
	sub, err := s.service.CreateSubscription(s.ctx, tenantID, planTier, billingCycle)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), sub)
	s.mockSubRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestCreateSubscription_CreateError 测试创建订阅失败
func (s *SubscriptionServiceTestSuite) TestCreateSubscription_CreateError() {
	// Arrange
	tenantID := "tenant-123"
	planTier := entity.SubscriptionTierFree
	billingCycle := entity.BillingCycleMonthly

	s.mockSubRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(errors.New("database error"))

	// Act
	sub, err := s.service.CreateSubscription(s.ctx, tenantID, planTier, billingCycle)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), sub)
	assert.Contains(s.T(), err.Error(), "failed to create subscription")
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestCreateSubscription_QuotaCreateError 测试创建配额失败
func (s *SubscriptionServiceTestSuite) TestCreateSubscription_QuotaCreateError() {
	// Arrange
	tenantID := "tenant-123"
	planTier := entity.SubscriptionTierFree
	billingCycle := entity.BillingCycleMonthly

	s.mockSubRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(nil)
	s.mockQuotaRepo.On("Create", s.ctx, mock.AnythingOfType("*entity.Quota")).Return(errors.New("quota error")).Once()

	// Act
	sub, err := s.service.CreateSubscription(s.ctx, tenantID, planTier, billingCycle)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), sub)
	assert.Contains(s.T(), err.Error(), "failed to create quota")
	s.mockSubRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestUpdateSubscription 测试更新订阅
func (s *SubscriptionServiceTestSuite) TestUpdateSubscription() {
	// Arrange
	sub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       "tenant-123",
		PlanTier:       entity.SubscriptionTierFree,
		AutoRenew:      true,
	}

	s.mockSubRepo.On("Update", s.ctx, sub).Return(nil)

	// Act
	err := s.service.UpdateSubscription(s.ctx, sub)

	// Assert
	assert.NoError(s.T(), err)
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestUpgradeTier_FreeToPro 测试从免费版升级到专业版
func (s *SubscriptionServiceTestSuite) TestUpgradeTier_FreeToPro() {
	// Arrange
	tenantID := "tenant-123"
	newTier := entity.SubscriptionTierPro

	existingSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierFree,
		Status:         entity.SubscriptionStatusActive,
	}

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return(existingSub, nil)
	s.mockSubRepo.On("Update", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(nil)

	// Mock quota updates
	existingQuota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
	}

	s.mockQuotaRepo.On("GetByTenantAndResource", s.ctx, tenantID, mock.AnythingOfType("entity.ResourceType")).Return(existingQuota, nil).Times(4)
	s.mockQuotaRepo.On("Update", s.ctx, mock.AnythingOfType("*entity.Quota")).Return(nil).Times(4)

	// Act
	sub, err := s.service.UpgradeTier(s.ctx, tenantID, newTier)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), sub)
	assert.Equal(s.T(), newTier, sub.PlanTier)
	s.mockSubRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestUpgradeTier_ProToEnterprise 测试从专业版升级到企业版
func (s *SubscriptionServiceTestSuite) TestUpgradeTier_ProToEnterprise() {
	// Arrange
	tenantID := "tenant-123"
	newTier := entity.SubscriptionTierEnterprise

	existingSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
	}

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return(existingSub, nil)
	s.mockSubRepo.On("Update", s.ctx, mock.AnythingOfType("*entity.Subscription")).Return(nil)

	existingQuota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     100,
	}

	s.mockQuotaRepo.On("GetByTenantAndResource", s.ctx, tenantID, mock.AnythingOfType("entity.ResourceType")).Return(existingQuota, nil).Times(4)
	s.mockQuotaRepo.On("Update", s.ctx, mock.AnythingOfType("*entity.Quota")).Return(nil).Times(4)

	// Act
	sub, err := s.service.UpgradeTier(s.ctx, tenantID, newTier)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newTier, sub.PlanTier)
	s.mockSubRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestUpgradeTier_SameTier 测试升级到相同等级
func (s *SubscriptionServiceTestSuite) TestUpgradeTier_SameTier() {
	// Arrange
	tenantID := "tenant-123"
	newTier := entity.SubscriptionTierPro

	existingSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
	}

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return(existingSub, nil)

	// Act
	sub, err := s.service.UpgradeTier(s.ctx, tenantID, newTier)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), existingSub, sub)
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestUpgradeTier_SubscriptionNotFound 测试订阅不存在
func (s *SubscriptionServiceTestSuite) TestUpgradeTier_SubscriptionNotFound() {
	// Arrange
	tenantID := "tenant-123"
	newTier := entity.SubscriptionTierPro

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return((*entity.Subscription)(nil), nil)

	// Act
	sub, err := s.service.UpgradeTier(s.ctx, tenantID, newTier)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), sub)
	assert.Contains(s.T(), err.Error(), "subscription not found")
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestUpgradeTier_DowngradeNotAllowed 测试不允许降级
func (s *SubscriptionServiceTestSuite) TestUpgradeTier_DowngradeNotAllowed() {
	// Arrange
	tenantID := "tenant-123"
	newTier := entity.SubscriptionTierFree // 降级

	existingSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro, // 从高级别
		Status:         entity.SubscriptionStatusActive,
	}

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return(existingSub, nil)

	// Act
	sub, err := s.service.UpgradeTier(s.ctx, tenantID, newTier)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), sub)
	assert.Contains(s.T(), err.Error(), "cannot downgrade")
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestCancelSubscription 测试取消订阅
func (s *SubscriptionServiceTestSuite) TestCancelSubscription() {
	// Arrange
	tenantID := "tenant-123"

	existingSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
		AutoRenew:      true,
	}

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return(existingSub, nil)
	s.mockSubRepo.On("Update", s.ctx, mock.MatchedBy(func(sub *entity.Subscription) bool {
		return sub.Status == entity.SubscriptionStatusCancelled && !sub.AutoRenew
	})).Return(nil)

	// Act
	err := s.service.CancelSubscription(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestCancelSubscription_NotFound 测试取消不存在的订阅
func (s *SubscriptionServiceTestSuite) TestCancelSubscription_NotFound() {
	// Arrange
	tenantID := "tenant-123"

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return((*entity.Subscription)(nil), nil)

	// Act
	err := s.service.CancelSubscription(s.ctx, tenantID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "subscription not found")
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestGetSubscription 测试获取订阅
func (s *SubscriptionServiceTestSuite) TestGetSubscription() {
	// Arrange
	tenantID := "tenant-123"
	expectedSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
	}

	s.mockSubRepo.On("GetByTenantID", s.ctx, tenantID).Return(expectedSub, nil)

	// Act
	sub, err := s.service.GetSubscription(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedSub, sub)
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestCheckSubscriptionStatus 测试检查订阅状态
func (s *SubscriptionServiceTestSuite) TestCheckSubscriptionStatus() {
	// Arrange
	now := time.Now()
	expiredSub := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       "tenant-123",
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
		EndDate:        &now,
	}
	// 设置过期时间为昨天
	yesterday := now.Add(-24 * time.Hour)
	expiredSub.EndDate = &yesterday

	activeSub := &entity.Subscription{
		SubscriptionID: "sub-456",
		TenantID:       "tenant-456",
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
		EndDate:        nil,
	}

	subs := []*entity.Subscription{expiredSub, activeSub}

	s.mockSubRepo.On("List", s.ctx, mock.MatchedBy(func(filter *repository.SubscriptionFilter) bool {
		return filter.Status == entity.SubscriptionStatusActive
	})).Return(subs, int64(2), nil)
	s.mockSubRepo.On("UpdateStatus", s.ctx, "sub-123", entity.SubscriptionStatusExpired).Return(nil)

	// Act
	err := s.service.CheckSubscriptionStatus(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestCheckSubscriptionStatus_ListError 测试检查订阅状态失败
func (s *SubscriptionServiceTestSuite) TestCheckSubscriptionStatus_ListError() {
	// Arrange
	s.mockSubRepo.On("List", s.ctx, mock.AnythingOfType("*repository.SubscriptionFilter")).Return(
		([]*entity.Subscription)(nil), int64(0), errors.New("database error"))

	// Act
	err := s.service.CheckSubscriptionStatus(s.ctx)

	// Assert
	assert.Error(s.T(), err)
	s.mockSubRepo.AssertExpectations(s.T())
}

// TestGetQuotasForTier_FreeTier 测试获取免费版配额配置
func (s *SubscriptionServiceTestSuite) TestGetQuotasForTier_FreeTier() {
	// Act
	quotas := s.service.getQuotasForTier(entity.SubscriptionTierFree)

	// Assert
	assert.Len(s.T(), quotas, 4)
	// 验证资源类型
	resourceTypes := make(map[entity.ResourceType]bool)
	for _, q := range quotas {
		resourceTypes[q.ResourceType] = true
	}
	assert.True(s.T(), resourceTypes[entity.ResourceTypeBots])
	assert.True(s.T(), resourceTypes[entity.ResourceTypeMessages])
	assert.True(s.T(), resourceTypes[entity.ResourceTypeStorage])
	assert.True(s.T(), resourceTypes[entity.ResourceTypeTeamMembers])

	// 验证免费版限制
	for _, q := range quotas {
		if q.ResourceType == entity.ResourceTypeBots {
			assert.Equal(s.T(), 10, q.MaxLimit)
		}
		if q.ResourceType == entity.ResourceTypeMessages {
			assert.Equal(s.T(), 1000, q.MaxLimit)
		}
	}
}

// TestGetQuotasForTier_EnterpriseTier 测试获取企业版配额配置
func (s *SubscriptionServiceTestSuite) TestGetQuotasForTier_EnterpriseTier() {
	// Act
	quotas := s.service.getQuotasForTier(entity.SubscriptionTierEnterprise)

	// Assert
	assert.Len(s.T(), quotas, 4)
	// 验证企业版无限制
	for _, q := range quotas {
		assert.Equal(s.T(), -1, q.MaxLimit, "Enterprise tier should have unlimited quota for %s", q.ResourceType)
	}
}

// TestIsUpgrade 测试等级升级判断
func (s *SubscriptionServiceTestSuite) TestIsUpgrade() {
	tests := []struct {
		name       string
		current    entity.SubscriptionTier
		newTier    entity.SubscriptionTier
		wantResult bool
	}{
		{"Free to Pro", entity.SubscriptionTierFree, entity.SubscriptionTierPro, true},
		{"Free to Enterprise", entity.SubscriptionTierFree, entity.SubscriptionTierEnterprise, true},
		{"Pro to Enterprise", entity.SubscriptionTierPro, entity.SubscriptionTierEnterprise, true},
		{"Same Tier", entity.SubscriptionTierPro, entity.SubscriptionTierPro, true},
		{"Pro to Free", entity.SubscriptionTierPro, entity.SubscriptionTierFree, false},
		{"Enterprise to Free", entity.SubscriptionTierEnterprise, entity.SubscriptionTierFree, false},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			result := s.service.isUpgrade(tt.current, tt.newTier)
			assert.Equal(t, tt.wantResult, result)
		})
	}
}

// 运行测试套件
func TestSubscriptionServiceSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionServiceTestSuite))
}
