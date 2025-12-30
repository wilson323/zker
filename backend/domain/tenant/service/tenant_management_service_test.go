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

// MockTenantRepositoryForManagement 租户仓储 Mock（用于管理服务）
type MockTenantRepositoryForManagement struct {
	mock.Mock
}

func (m *MockTenantRepositoryForManagement) Create(ctx context.Context, tenant *entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepositoryForManagement) GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepositoryForManagement) GetByName(ctx context.Context, name string) (*entity.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepositoryForManagement) GetBySubdomain(ctx context.Context, subdomain string) (*entity.Tenant, error) {
	args := m.Called(ctx, subdomain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepositoryForManagement) Update(ctx context.Context, tenant *entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepositoryForManagement) Delete(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockTenantRepositoryForManagement) List(ctx context.Context, filter *repository.TenantFilter) ([]*entity.Tenant, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Tenant), args.Get(1).(int64), args.Error(2)
}

func (m *MockTenantRepositoryForManagement) UpdateStatus(ctx context.Context, tenantID string, status entity.TenantStatus) error {
	args := m.Called(ctx, tenantID, status)
	return args.Error(0)
}

func (m *MockTenantRepositoryForManagement) UpdateSubscriptionTier(ctx context.Context, tenantID string, tier entity.SubscriptionTier) error {
	args := m.Called(ctx, tenantID, tier)
	return args.Error(0)
}

// MockSubscriptionRepositoryForManagement 订阅仓储 Mock（用于管理服务）
type MockSubscriptionRepositoryForManagement struct {
	mock.Mock
}

func (m *MockSubscriptionRepositoryForManagement) Create(ctx context.Context, sub *entity.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepositoryForManagement) GetByID(ctx context.Context, subscriptionID string) (*entity.Subscription, error) {
	args := m.Called(ctx, subscriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepositoryForManagement) GetByTenantID(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepositoryForManagement) GetByTenant(ctx context.Context, tenantID string) (*entity.Subscription, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepositoryForManagement) GetAll(ctx context.Context) ([]*entity.Subscription, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepositoryForManagement) Update(ctx context.Context, sub *entity.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockSubscriptionRepositoryForManagement) List(ctx context.Context, filter *repository.SubscriptionFilter) ([]*entity.Subscription, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Subscription), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubscriptionRepositoryForManagement) UpdateStatus(ctx context.Context, subscriptionID string, status entity.SubscriptionStatus) error {
	args := m.Called(ctx, subscriptionID, status)
	return args.Error(0)
}

// TenantManagementServiceTestSuite 测试套件
type TenantManagementServiceTestSuite struct {
	suite.Suite
	service               *TenantManagementService
	mockTenantRepo        *MockTenantRepositoryForManagement
	mockQuotaRepo         *MockQuotaRepository
	mockSubscriptionRepo  *MockSubscriptionRepositoryForManagement
	ctx                   context.Context
}

func (s *TenantManagementServiceTestSuite) SetupTest() {
	s.mockTenantRepo = new(MockTenantRepositoryForManagement)
	s.mockQuotaRepo = new(MockQuotaRepository)
	s.mockSubscriptionRepo = new(MockSubscriptionRepositoryForManagement)
	s.service = NewTenantManagementService(nil, s.mockTenantRepo, s.mockQuotaRepo, s.mockSubscriptionRepo)
	s.ctx = context.Background()
}

// TestGetTenantInfo_Success 测试获取租户信息成功
func (s *TenantManagementServiceTestSuite) TestGetTenantInfo_Success() {
	// Arrange
	tenantID := "tenant-123"
	tenant := &entity.Tenant{
		TenantID:         tenantID,
		TenantName:       "Test Tenant",
		Status:           entity.TenantStatusActive,
		SubscriptionTier: entity.SubscriptionTierPro,
	}

	subscription := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
	}

	quotas := []*entity.Quota{
		{QuotaID: "quota-1", TenantID: tenantID, ResourceType: entity.ResourceTypeBots, MaxLimit: 50, UsedCount: 10},
		{QuotaID: "quota-2", TenantID: tenantID, ResourceType: entity.ResourceTypeMessages, MaxLimit: 50000, UsedCount: 1000},
	}

	s.mockTenantRepo.On("GetByID", s.ctx, tenantID).Return(tenant, nil)
	s.mockSubscriptionRepo.On("GetByTenant", s.ctx, tenantID).Return(subscription, nil)
	s.mockQuotaRepo.On("GetByTenant", s.ctx, tenantID).Return(quotas, nil)

	// Act
	detail, err := s.service.GetTenantInfo(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), detail)
	assert.Equal(s.T(), tenant, detail.Tenant)
	assert.Equal(s.T(), subscription, detail.Subscription)
	assert.Len(s.T(), detail.Quotas, 2)
	assert.NotNil(s.T(), detail.Statistics)
	s.mockTenantRepo.AssertExpectations(s.T())
	s.mockSubscriptionRepo.AssertExpectations(s.T())
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestGetTenantInfo_TenantNotFound 测试租户不存在
func (s *TenantManagementServiceTestSuite) TestGetTenantInfo_TenantNotFound() {
	// Arrange
	tenantID := "tenant-123"

	s.mockTenantRepo.On("GetByID", s.ctx, tenantID).Return((*entity.Tenant)(nil), nil)

	// Act
	detail, err := s.service.GetTenantInfo(s.ctx, tenantID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), detail)
	assert.Contains(s.T(), err.Error(), "tenant not found")
	s.mockTenantRepo.AssertExpectations(s.T())
}

// TestUpdateTenantInfo_Success 测试更新租户信息成功
func (s *TenantManagementServiceTestSuite) TestUpdateTenantInfo_Success() {
	// Arrange
	tenantID := "tenant-123"
	tenant := &entity.Tenant{
		TenantID:         tenantID,
		TenantName:       "Old Name",
		Status:           entity.TenantStatusActive,
		SubscriptionTier: entity.SubscriptionTierPro,
	}

	req := &UpdateTenantRequest{
		TenantName: "New Name",
	}

	s.mockTenantRepo.On("GetByID", s.ctx, tenantID).Return(tenant, nil)
	s.mockTenantRepo.On("Update", s.ctx, mock.MatchedBy(func(t *entity.Tenant) bool {
		return t.TenantName == "New Name"
	})).Return(nil)

	// Act
	updatedTenant, err := s.service.UpdateTenantInfo(s.ctx, tenantID, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), updatedTenant)
	assert.Equal(s.T(), "New Name", updatedTenant.TenantName)
	s.mockTenantRepo.AssertExpectations(s.T())
}

// TestUpdateTenantInfo_StatusUpdate 测试更新租户状态
func (s *TenantManagementServiceTestSuite) TestUpdateTenantInfo_StatusUpdate() {
	// Arrange
	tenantID := "tenant-123"
	tenant := &entity.Tenant{
		TenantID:         tenantID,
		TenantName:       "Test Tenant",
		Status:           entity.TenantStatusActive,
		SubscriptionTier: entity.SubscriptionTierPro,
	}

	req := &UpdateTenantRequest{
		Status: entity.TenantStatusSuspended,
	}

	s.mockTenantRepo.On("GetByID", s.ctx, tenantID).Return(tenant, nil)
	s.mockTenantRepo.On("Update", s.ctx, mock.MatchedBy(func(t *entity.Tenant) bool {
		return t.Status == entity.TenantStatusSuspended
	})).Return(nil)

	// Act
	updatedTenant, err := s.service.UpdateTenantInfo(s.ctx, tenantID, req)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), entity.TenantStatusSuspended, updatedTenant.Status)
	s.mockTenantRepo.AssertExpectations(s.T())
}

// TestGetQuotaUsage_Success 测试获取配额使用情况成功
func (s *TenantManagementServiceTestSuite) TestGetQuotaUsage_Success() {
	// Arrange
	tenantID := "tenant-123"
	quotas := []*entity.Quota{
		{QuotaID: "quota-1", TenantID: tenantID, ResourceType: entity.ResourceTypeBots, MaxLimit: 100, UsedCount: 80},
		{QuotaID: "quota-2", TenantID: tenantID, ResourceType: entity.ResourceTypeMessages, MaxLimit: 1000, UsedCount: 200},
	}

	s.mockQuotaRepo.On("GetByTenant", s.ctx, tenantID).Return(quotas, nil)

	// Act
	usages, err := s.service.GetQuotaUsage(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), usages, 2)

	// 验证第一个配额（接近上限）
	botUsage := usages[0]
	assert.Equal(s.T(), 80, botUsage.UsedCount)
	assert.Equal(s.T(), 100, botUsage.MaxLimit)
	assert.Equal(s.T(), 20, botUsage.RemainingCount)
	assert.InDelta(s.T(), 80.0, botUsage.UsagePercentage, 0.1)
	assert.True(s.T(), botUsage.IsNearLimit)

	// 验证第二个配额（正常）
	msgUsage := usages[1]
	assert.Equal(s.T(), 200, msgUsage.UsedCount)
	assert.Equal(s.T(), 1000, msgUsage.MaxLimit)
	assert.Equal(s.T(), 800, msgUsage.RemainingCount)
	assert.InDelta(s.T(), 20.0, msgUsage.UsagePercentage, 0.1)
	assert.False(s.T(), msgUsage.IsNearLimit)

	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestCheckQuotaAvailable_Available 测试配额可用
func (s *TenantManagementServiceTestSuite) TestCheckQuotaAvailable_Available() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	requiredCount := 10
	quota := &entity.Quota{
		QuotaID:      "quota-1",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    50,
	}

	s.mockQuotaRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	available, err := s.service.CheckQuotaAvailable(s.ctx, tenantID, resourceType, requiredCount)

	// Assert
	assert.NoError(s.T(), err)
	assert.True(s.T(), available)
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestCheckQuotaAvailable_NotAvailable 测试配额不可用
func (s *TenantManagementServiceTestSuite) TestCheckQuotaAvailable_NotAvailable() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	requiredCount := 60
	quota := &entity.Quota{
		QuotaID:      "quota-1",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    50,
	}

	s.mockQuotaRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	available, err := s.service.CheckQuotaAvailable(s.ctx, tenantID, resourceType, requiredCount)

	// Assert
	assert.NoError(s.T(), err)
	assert.False(s.T(), available)
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestCheckQuotaAvailable_NotFound 测试配额不存在
func (s *TenantManagementServiceTestSuite) TestCheckQuotaAvailable_NotFound() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots

	s.mockQuotaRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return((*entity.Quota)(nil), nil)

	// Act
	available, err := s.service.CheckQuotaAvailable(s.ctx, tenantID, resourceType, 10)

	// Assert
	assert.Error(s.T(), err)
	assert.False(s.T(), available)
	assert.Contains(s.T(), err.Error(), "quota not found")
	s.mockQuotaRepo.AssertExpectations(s.T())
}

// TestGetSubscription_Success 测试获取订阅成功
func (s *TenantManagementServiceTestSuite) TestGetSubscription_Success() {
	// Arrange
	tenantID := "tenant-123"
	subscription := &entity.Subscription{
		SubscriptionID: "sub-123",
		TenantID:       tenantID,
		PlanTier:       entity.SubscriptionTierPro,
		Status:         entity.SubscriptionStatusActive,
	}

	s.mockSubscriptionRepo.On("GetByTenant", s.ctx, tenantID).Return(subscription, nil)

	// Act
	detail, err := s.service.GetSubscription(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), detail)
	assert.Equal(s.T(), subscription, detail.Subscription)
	assert.NotEmpty(s.T(), detail.Features)
	assert.Equal(s.T(), "active", detail.Status)
	s.mockSubscriptionRepo.AssertExpectations(s.T())
}

// TestGetSubscription_NotFound 测试订阅不存在
func (s *TenantManagementServiceTestSuite) TestGetSubscription_NotFound() {
	// Arrange
	tenantID := "tenant-123"

	s.mockSubscriptionRepo.On("GetByTenant", s.ctx, tenantID).Return((*entity.Subscription)(nil), nil)

	// Act
	detail, err := s.service.GetSubscription(s.ctx, tenantID)

	// Assert
	assert.Error(s.T(), err)
	assert.Nil(s.T(), detail)
	assert.Contains(s.T(), err.Error(), "subscription not found")
	s.mockSubscriptionRepo.AssertExpectations(s.T())
}

// TestGetSubscriptionFeatures 测试获取订阅功能列表
func (s *TenantManagementServiceTestSuite) TestGetSubscriptionFeatures() {
	tests := []struct {
		name          string
		tier          entity.SubscriptionTier
		minFeatures   int
		expectFeature string
	}{
		{
			name:          "Free Tier",
			tier:          entity.SubscriptionTierFree,
			minFeatures:   3,
			expectFeature: "最多3个Bot",
		},
		{
			name:          "Pro Tier",
			tier:          entity.SubscriptionTierPro,
			minFeatures:   4,
			expectFeature: "优先技术支持",
		},
		{
			name:          "Enterprise Tier",
			tier:          entity.SubscriptionTierEnterprise,
			minFeatures:   6,
			expectFeature: "SLA保障",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			features := s.service.getSubscriptionFeatures(tt.tier)
			assert.GreaterOrEqual(t, len(features), tt.minFeatures)
			assert.Contains(t, features, tt.expectFeature)
		})
	}
}

// TestCalculateSubscriptionStatus 测试计算订阅状态
func (s *TenantManagementServiceTestSuite) TestCalculateSubscriptionStatus() {
	tests := []struct {
		name     string
		sub      *entity.Subscription
		expected string
	}{
		{
			name: "Active Subscription",
			sub: &entity.Subscription{
				Status:  entity.SubscriptionStatusActive,
				EndDate: nil,
			},
			expected: "active",
		},
		{
			name: "Suspended Subscription",
			sub: &entity.Subscription{
				Status:  entity.SubscriptionStatusSuspended,
				EndDate: nil,
			},
			expected: "suspended",
		},
		{
			name: "Expired Subscription",
			sub: &entity.Subscription{
				Status:  entity.SubscriptionStatusActive,
				EndDate: func() *time.Time { t := time.Now().Add(-24 * time.Hour); return &t }(),
			},
			expected: "expired",
		},
		{
			name: "Future End Date",
			sub: &entity.Subscription{
				Status:  entity.SubscriptionStatusActive,
				EndDate: func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
			},
			expected: "active",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			result := s.service.calculateSubscriptionStatus(tt.sub)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsValidUpgrade 测试升级路径验证
func (s *TenantManagementServiceTestSuite) TestIsValidUpgrade() {
	tests := []struct {
		name     string
		from     entity.SubscriptionTier
		to       entity.SubscriptionTier
		expected bool
	}{
		{"Free to Pro", entity.SubscriptionTierFree, entity.SubscriptionTierPro, true},
		{"Free to Enterprise", entity.SubscriptionTierFree, entity.SubscriptionTierEnterprise, true},
		{"Pro to Enterprise", entity.SubscriptionTierPro, entity.SubscriptionTierEnterprise, true},
		{"Enterprise to any", entity.SubscriptionTierEnterprise, entity.SubscriptionTierPro, false},
		{"Pro to Free", entity.SubscriptionTierPro, entity.SubscriptionTierFree, false},
		{"Free to Free", entity.SubscriptionTierFree, entity.SubscriptionTierFree, false},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			result := s.service.isValidUpgrade(tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetPlanLimits 测试获取套餐配额限制
func (s *TenantManagementServiceTestSuite) TestGetPlanLimits() {
	tests := []struct {
		name          string
		tier          entity.SubscriptionTier
		expectedBots  int
		expectedMsgs  int
	}{
		{
			name:         "Free Tier",
			tier:         entity.SubscriptionTierFree,
			expectedBots: 3,
			expectedMsgs: 1000,
		},
		{
			name:         "Pro Tier",
			tier:         entity.SubscriptionTierPro,
			expectedBots: 50,
			expectedMsgs: 50000,
		},
		{
			name:         "Enterprise Tier",
			tier:         entity.SubscriptionTierEnterprise,
			expectedBots: 999999,
			expectedMsgs: 999999,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			limits := s.service.getPlanLimits(tt.tier)
			assert.Equal(t, tt.expectedBots, limits[entity.ResourceTypeBots])
			assert.Equal(t, tt.expectedMsgs, limits[entity.ResourceTypeMessages])
		})
	}
}

// TestCalculateStatistics 测试计算统计信息
func (s *TenantManagementServiceTestSuite) TestCalculateStatistics() {
	// Arrange
	tenant := &entity.Tenant{
		TenantID:         "tenant-123",
		TenantName:       "Test Tenant",
		Status:           entity.TenantStatusActive,
		SubscriptionTier: entity.SubscriptionTierPro,
		CreatedAt:        time.Now().UnixMilli(),
	}

	quotas := []entity.Quota{
		{QuotaID: "quota-1", ResourceType: entity.ResourceTypeBots, UsedCount: 10},
		{QuotaID: "quota-2", ResourceType: entity.ResourceTypeTeamMembers, UsedCount: 5},
	}

	// Act
	stats := s.service.calculateStatistics(tenant, quotas)

	// Assert
	assert.NotNil(s.T(), stats)
	assert.Equal(s.T(), tenant.TenantID, stats.TenantID)
	assert.Equal(s.T(), tenant.SubscriptionTier, stats.SubscriptionTier)
	assert.Equal(s.T(), 5, stats.MemberCount)
	assert.Equal(s.T(), 10, stats.BotCount)
}

// 运行测试套件
func TestTenantManagementServiceSuite(t *testing.T) {
	suite.Run(t, new(TenantManagementServiceTestSuite))
}

// TestQuotaService_ErrorWrapper 测试错误包装
func TestQuotaService_ErrorWrapper(t *testing.T) {
	// 测试 CheckQuota 错误包装
	t.Run("CheckQuota wraps repository error", func(t *testing.T) {
		mockRepo := new(MockQuotaRepository)
		service := NewQuotaService(mockRepo)
		ctx := context.Background()

		mockRepo.On("GetByTenantAndResource", ctx, "tenant-123", entity.ResourceTypeBots).
			Return((*entity.Quota)(nil), errors.New("db connection failed"))

		err := service.CheckQuota(ctx, "tenant-123", entity.ResourceTypeBots, 10)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get quota")
		mockRepo.AssertExpectations(t)
	})
}

// TestSubscriptionService_ErrorWrapper 测试订阅服务错误处理
func TestSubscriptionService_ErrorWrapper(t *testing.T) {
	// 测试 UpgradeTier 错误包装
	t.Run("UpgradeTier returns error when subscription not found", func(t *testing.T) {
		mockSubRepo := new(MockSubscriptionRepository)
		mockQuotaRepo := new(MockQuotaRepository)
		service := NewSubscriptionService(mockSubRepo, mockQuotaRepo)
		ctx := context.Background()

		mockSubRepo.On("GetByTenantID", ctx, "tenant-123").Return((*entity.Subscription)(nil), nil)

		sub, err := service.UpgradeTier(ctx, "tenant-123", entity.SubscriptionTierPro)

		assert.Error(t, err)
		assert.Nil(t, sub)
		assert.Contains(t, err.Error(), "subscription not found")
		mockSubRepo.AssertExpectations(t)
	})
}

// TestTenantManagementService_ConsumeQuota_Error 测试消耗配额错误情况
func TestTenantManagementService_ConsumeQuota_Error(t *testing.T) {
	t.Skip("Skipping test that requires database connection")
	service := &TenantManagementService{db: nil}

	// 没有真实数据库，我们只能测试错误路径
	// 实际项目中应该使用 sqlmock 或 testcontainers

	t.Run("ConsumeQuota requires valid database connection", func(t *testing.T) {
		ctx := context.Background()
		err := service.ConsumeQuota(ctx, "tenant-123", entity.ResourceTypeBots, 10)
		assert.Error(t, err)
	})
}
