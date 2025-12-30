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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// MockQuotaRepository 配额仓储 Mock
type MockQuotaRepository struct {
	mock.Mock
}

func (m *MockQuotaRepository) Create(ctx context.Context, quota *entity.Quota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) GetByID(ctx context.Context, quotaID string) (*entity.Quota, error) {
	args := m.Called(ctx, quotaID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) GetByTenantAndResource(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error) {
	args := m.Called(ctx, tenantID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) List(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) GetAll(ctx context.Context) ([]*entity.Quota, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) Update(ctx context.Context, quota *entity.Quota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) UpdateUsedCount(ctx context.Context, quotaID string, count int) error {
	args := m.Called(ctx, quotaID, count)
	return args.Error(0)
}

func (m *MockQuotaRepository) ResetUsage(ctx context.Context, quotaID string) error {
	args := m.Called(ctx, quotaID)
	return args.Error(0)
}

func (m *MockQuotaRepository) Delete(ctx context.Context, quotaID string) error {
	args := m.Called(ctx, quotaID)
	return args.Error(0)
}

func (m *MockQuotaRepository) BatchCreate(ctx context.Context, quotas []*entity.Quota) error {
	args := m.Called(ctx, quotas)
	return args.Error(0)
}

// QuotaServiceTestSuite 测试套件
type QuotaServiceTestSuite struct {
	suite.Suite
	service  *QuotaService
	mockRepo *MockQuotaRepository
	ctx      context.Context
}

func (s *QuotaServiceTestSuite) SetupTest() {
	s.mockRepo = new(MockQuotaRepository)
	s.service = NewQuotaService(s.mockRepo)
	s.ctx = context.Background()
}

// TestCheckQuota_Success 测试配额检查成功
func (s *QuotaServiceTestSuite) TestCheckQuota_Success() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	requiredCount := 5
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    10,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	err := s.service.CheckQuota(s.ctx, tenantID, resourceType, requiredCount)

	// Assert
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
}

// TestCheckQuota_Unlimited 测试无限制配额
func (s *QuotaServiceTestSuite) TestCheckQuota_Unlimited() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	requiredCount := 10000
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     -1, // 无限制
		UsedCount:    10,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	err := s.service.CheckQuota(s.ctx, tenantID, resourceType, requiredCount)

	// Assert
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
}

// TestCheckQuota_Exceeded 测试配额超额
func (s *QuotaServiceTestSuite) TestCheckQuota_Exceeded() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	requiredCount := 50
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    60,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	err := s.service.CheckQuota(s.ctx, tenantID, resourceType, requiredCount)

	// Assert
	assert.Error(s.T(), err)
	var quotaErr *QuotaExceededError
	assert.True(s.T(), errors.As(err, &quotaErr))
	assert.Equal(s.T(), tenantID, quotaErr.TenantID)
	assert.Equal(s.T(), resourceType, quotaErr.ResourceType)
	assert.Equal(s.T(), 60, quotaErr.CurrentUsage)
	assert.Equal(s.T(), 100, quotaErr.MaxLimit)
	assert.Equal(s.T(), 50, quotaErr.Required)
	assert.Equal(s.T(), 40, quotaErr.GetRemaining()) // MaxLimit - CurrentUsage = 100 - 60 = 40
	assert.InDelta(s.T(), 60.0, quotaErr.UsagePercent, 0.01) // 60/100 * 100 = 60%
	s.mockRepo.AssertExpectations(s.T())
}

// TestCheckQuota_QuotaNotFound 测试配额不存在
func (s *QuotaServiceTestSuite) TestCheckQuota_QuotaNotFound() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return((*entity.Quota)(nil), nil)

	// Act
	err := s.service.CheckQuota(s.ctx, tenantID, resourceType, 5)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "quota not found")
	s.mockRepo.AssertExpectations(s.T())
}

// TestCheckQuota_RepositoryError 测试仓储错误
func (s *QuotaServiceTestSuite) TestCheckQuota_RepositoryError() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return((*entity.Quota)(nil), errors.New("database error"))

	// Act
	err := s.service.CheckQuota(s.ctx, tenantID, resourceType, 5)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to get quota")
	s.mockRepo.AssertExpectations(s.T())
}

// TestConsumeQuota_Success 测试消费配额成功
func (s *QuotaServiceTestSuite) TestConsumeQuota_Success() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	count := 5
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    10,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil).Twice()
	s.mockRepo.On("UpdateUsedCount", s.ctx, "quota-123", count).Return(nil)

	// Act
	previousUsage, err := s.service.ConsumeQuota(s.ctx, tenantID, resourceType, count)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 10, previousUsage)
	s.mockRepo.AssertExpectations(s.T())
}

// TestConsumeQuota_QuotaExceeded 测试消费配额超额
func (s *QuotaServiceTestSuite) TestConsumeQuota_QuotaExceeded() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	count := 50
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    60,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	previousUsage, err := s.service.ConsumeQuota(s.ctx, tenantID, resourceType, count)

	// Assert
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 0, previousUsage)
	s.mockRepo.AssertExpectations(s.T())
}

// TestRollbackQuota_Success 测试回滚配额成功
func (s *QuotaServiceTestSuite) TestRollbackQuota_Success() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	count := 5
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    10,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)
	s.mockRepo.On("UpdateUsedCount", s.ctx, "quota-123", -count).Return(nil)

	// Act
	err := s.service.RollbackQuota(s.ctx, tenantID, resourceType, count)

	// Assert
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
}

// TestRollbackQuota_NotFound 测试回滚配额不存在
func (s *QuotaServiceTestSuite) TestRollbackQuota_NotFound() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return((*entity.Quota)(nil), nil)

	// Act
	err := s.service.RollbackQuota(s.ctx, tenantID, resourceType, 5)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "quota not found")
	s.mockRepo.AssertExpectations(s.T())
}

// TestGetQuota_Success 测试获取配额成功
func (s *QuotaServiceTestSuite) TestGetQuota_Success() {
	// Arrange
	tenantID := "tenant-123"
	resourceType := entity.ResourceTypeBots
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     tenantID,
		ResourceType: resourceType,
		MaxLimit:     100,
		UsedCount:    10,
	}

	s.mockRepo.On("GetByTenantAndResource", s.ctx, tenantID, resourceType).Return(quota, nil)

	// Act
	result, err := s.service.GetQuota(s.ctx, tenantID, resourceType)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), quota, result)
	s.mockRepo.AssertExpectations(s.T())
}

// TestGetAllQuotas_Success 测试获取所有配额成功
func (s *QuotaServiceTestSuite) TestGetAllQuotas_Success() {
	// Arrange
	tenantID := "tenant-123"
	quotas := []*entity.Quota{
		{QuotaID: "quota-1", TenantID: tenantID, ResourceType: entity.ResourceTypeBots},
		{QuotaID: "quota-2", TenantID: tenantID, ResourceType: entity.ResourceTypeMessages},
	}

	s.mockRepo.On("List", s.ctx, tenantID).Return(quotas, nil)

	// Act
	result, err := s.service.GetAllQuotas(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.mockRepo.AssertExpectations(s.T())
}

// TestResetQuota_Success 测试重置配额成功
func (s *QuotaServiceTestSuite) TestResetQuota_Success() {
	// Arrange
	quotaID := "quota-123"

	s.mockRepo.On("ResetUsage", s.ctx, quotaID).Return(nil)

	// Act
	err := s.service.ResetQuota(s.ctx, quotaID)

	// Assert
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
}

// TestCheckAndResetQuotas_Success 测试检查并重置配额
func (s *QuotaServiceTestSuite) TestCheckAndResetQuotas_Success() {
	// Arrange - 创建需要重置的配额（月度重置，上月重置）
	quotaToReset := &entity.Quota{
		QuotaID:     "quota-123",
		TenantID:    "tenant-123",
		ResourceType: entity.ResourceTypeMessages,
		MaxLimit:    1000,
		UsedCount:   500,
		ResetCycle:  entity.ResetCycleMonthly,
		LastResetAt: 1000000, // 很久之前
	}

	// 创建不需要重置的配额（永不重置）
	quotaNoReset := &entity.Quota{
		QuotaID:      "quota-456",
		TenantID:     "tenant-123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
		UsedCount:    5,
		ResetCycle:   entity.ResetCycleNever,
		LastResetAt:  1000000,
	}

	quotas := []*entity.Quota{quotaToReset, quotaNoReset}

	s.mockRepo.On("GetAll", s.ctx).Return(quotas, nil)
	s.mockRepo.On("ResetUsage", s.ctx, "quota-123").Return(nil)

	// Act
	err := s.service.CheckAndResetQuotas(s.ctx)

	// Assert
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
}

// TestCheckAndResetQuotas_ResetError 测试重置配额时部分失败
func (s *QuotaServiceTestSuite) TestCheckAndResetQuotas_ResetError() {
	// Arrange
	quota := &entity.Quota{
		QuotaID:      "quota-123",
		TenantID:     "tenant-123",
		ResourceType: entity.ResourceTypeMessages,
		MaxLimit:     1000,
		UsedCount:    500,
		ResetCycle:   entity.ResetCycleMonthly,
		LastResetAt:  1000000,
	}

	quotas := []*entity.Quota{quota}

	s.mockRepo.On("GetAll", s.ctx).Return(quotas, nil)
	s.mockRepo.On("ResetUsage", s.ctx, "quota-123").Return(errors.New("reset error"))

	// Act
	err := s.service.CheckAndResetQuotas(s.ctx)

	// Assert - 函数不应该返回错误，应该继续处理其他配额
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
}

// TestQuotaExceededError_GetRemaining 测试获取剩余配额
func (s *QuotaServiceTestSuite) TestQuotaExceededError_GetRemaining() {
	// Arrange
	err := &QuotaExceededError{
		CurrentUsage: 80,
		MaxLimit:     100,
	}

	// Act
	remaining := err.GetRemaining()

	// Assert
	assert.Equal(s.T(), 20, remaining)
}

// TestQuotaExceededError_GetRemainingNegative 测试获取剩余配额（负数）
func (s *QuotaServiceTestSuite) TestQuotaExceededError_GetRemainingNegative() {
	// Arrange
	err := &QuotaExceededError{
		CurrentUsage: 120,
		MaxLimit:     100,
	}

	// Act
	remaining := err.GetRemaining()

	// Assert
	assert.Equal(s.T(), 0, remaining)
}

// 运行测试套件
func TestQuotaServiceSuite(t *testing.T) {
	suite.Run(t, new(QuotaServiceTestSuite))
}
