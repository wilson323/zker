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

	"github.com/coze-studio/backend/domain/tenant/entity"
)

// MockQuotaRepository Mock配额仓储
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

func (m *MockQuotaRepository) Update(ctx context.Context, quota *entity.Quota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) UpdateUsedCount(ctx context.Context, quotaID string, delta int) error {
	args := m.Called(ctx, quotaID, delta)
	return args.Error(0)
}

func (m *MockQuotaRepository) ResetUsage(ctx context.Context, quotaID string) error {
	args := m.Called(ctx, quotaID)
	return args.Error(0)
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

func (m *MockQuotaRepository) GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Quota), args.Error(1)
}

// TestQuotaService_CheckQuota_Success 测试配额检查成功
func TestQuotaService_CheckQuota_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	quota := &entity.Quota{
		QuotaID:      "quota_123",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
		UsedCount:    5,
	}

	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(quota, nil)

	// Act
	err := service.CheckQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 3)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_CheckQuota_Exceeded 测试配额超限
func TestQuotaService_CheckQuota_Exceeded(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	quota := &entity.Quota{
		QuotaID:      "quota_123",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
		UsedCount:    9,
	}

	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(quota, nil)

	// Act
	err := service.CheckQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 2)

	// Assert
	assert.Error(t, err)
	var quotaErr *QuotaExceededError
	assert.True(t, errors.As(err, &quotaErr))
	assert.Equal(t, "tenant_123", quotaErr.TenantID)
	assert.Equal(t, entity.ResourceTypeBots, quotaErr.ResourceType)
	assert.Equal(t, 9, quotaErr.CurrentUsage)
	assert.Equal(t, 10, quotaErr.MaxLimit)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_CheckQuota_Unlimited 测试无限制配额
func TestQuotaService_CheckQuota_Unlimited(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	quota := &entity.Quota{
		QuotaID:      "quota_123",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     0, // 无限制
		UsedCount:    9999,
	}

	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(quota, nil)

	// Act
	err := service.CheckQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 10000)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_CheckQuota_AutoReset 测试自动重置
func TestQuotaService_CheckQuota_AutoReset(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	quota := &entity.Quota{
		QuotaID:      "quota_123",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
		UsedCount:    9,
		ResetCycle:   entity.ResetCycleDaily,
		LastResetAt:  yesterday.UnixMilli(),
	}

	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(quota, nil).Once()
	mockRepo.On("ResetUsage", mock.Anything, "quota_123").
		Return(nil)
	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(&entity.Quota{
			QuotaID:      "quota_123",
			TenantID:     "tenant_123",
			ResourceType: entity.ResourceTypeBots,
			MaxLimit:     10,
			UsedCount:    0, // 重置后
			ResetCycle:   entity.ResetCycleDaily,
			LastResetAt:  now.UnixMilli(),
		}, nil).Once()

	// Act
	err := service.CheckQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 5)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_ConsumeQuota 测试消费配额
func TestQuotaService_ConsumeQuota(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	quota := &entity.Quota{
		QuotaID:      "quota_123",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
		UsedCount:    5,
	}

	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(quota, nil)
	mockRepo.On("UpdateUsedCount", mock.Anything, "quota_123", 3).
		Return(nil)

	// Act
	err := service.ConsumeQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 3)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_ConsumeQuota_Exceeded 测试消费超限配额
func TestQuotaService_ConsumeQuota_Exceeded(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	quota := &entity.Quota{
		QuotaID:      "quota_123",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		MaxLimit:     10,
		UsedCount:    9,
	}

	mockRepo.On("GetByTenantAndResource", mock.Anything, "tenant_123", entity.ResourceTypeBots).
		Return(quota, nil)

	// Act
	err := service.ConsumeQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 2)

	// Assert
	assert.Error(t, err)
	var quotaErr *QuotaExceededError
	assert.True(t, errors.As(err, &quotaErr))
	mockRepo.AssertNotCalled(t, "UpdateUsedCount")
}

// TestQuotaService_RollbackQuota 测试回滚配额
func TestQuotaService_RollbackQuota(t *testing.T) {
	// Arrange
	mockRepo := new(MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	mockRepo.On("UpdateUsedCount", mock.Anything, "quota_123", -3).
		Return(nil)

	// Act
	err := service.RollbackQuota(context.Background(), "tenant_123", entity.ResourceTypeBots, 3)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaExceededError_Error 测试错误信息
func TestQuotaExceededError_Error(t *testing.T) {
	err := &QuotaExceededError{
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeBots,
		CurrentUsage: 9,
		MaxLimit:     10,
	}

	expected := "quota exceeded for tenant tenant_123, resource bots: usage 9, limit 10"
	assert.Equal(t, expected, err.Error())
}
