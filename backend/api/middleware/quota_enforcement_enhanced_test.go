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

package middleware

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
)

// MockQuotaRepository 配额仓储Mock
type MockQuotaRepository struct {
	mock.Mock
}

func (m *MockQuotaRepository) GetByTenant(ctx context.Context, tenantID string) ([]*entity.Quota, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) GetByTenantAndResource(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*entity.Quota, error) {
	args := m.Called(ctx, tenantID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Quota), args.Error(1)
}

func (m *MockQuotaRepository) Create(ctx context.Context, quota *entity.Quota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) Update(ctx context.Context, quota *entity.Quota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) Delete(ctx context.Context, quotaID string) error {
	args := m.Called(ctx, quotaID)
	return args.Error(0)
}

// TestQuotaChecker_CheckAndConsume_BoundaryConditions 测试边界条件
func TestQuotaChecker_CheckAndConsume_BoundaryConditions(t *testing.T) {
	t.Run("刚好达到配额上限", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		// 创建配额：已用9个，限制10个
		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    9,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		checker := NewQuotaChecker(db, mockRepo)

		// 消耗1个，应该成功
		success, err := checker.CheckAndConsume(context.Background(), "test_tenant", entity.ResourceTypeBots, 1)
		assert.NoError(t, err)
		assert.True(t, success)

		// 验证已用计数为10
		var updatedQuota entity.Quota
		err = db.Where("tenant_id = ? AND resource_type = ?", "test_tenant", entity.ResourceTypeBots).First(&updatedQuota).Error
		assert.NoError(t, err)
		assert.Equal(t, 10, updatedQuota.UsedCount)
	})

	t.Run("刚好超过配额上限", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		// 创建配额：已用10个，限制10个
		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    10,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		checker := NewQuotaChecker(db, mockRepo)

		// 尝试消耗1个，应该失败
		success, err := checker.CheckAndConsume(context.Background(), "test_tenant", entity.ResourceTypeBots, 1)
		assert.Error(t, err)
		assert.False(t, success)
		assert.Contains(t, err.Error(), "quota exceeded")
	})

	t.Run("消耗0个配额", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    5,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		checker := NewQuotaChecker(db, mockRepo)

		// 消耗0个，应该成功
		success, err := checker.CheckAndConsume(context.Background(), "test_tenant", entity.ResourceTypeBots, 0)
		assert.NoError(t, err)
		assert.True(t, success)

		// 验证计数不变
		var updatedQuota entity.Quota
		err = db.Where("tenant_id = ? AND resource_type = ?", "test_tenant", entity.ResourceTypeBots).First(&updatedQuota).Error
		assert.NoError(t, err)
		assert.Equal(t, 5, updatedQuota.UsedCount)
	})

	t.Run("配额不存在", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		checker := NewQuotaChecker(db, mockRepo)

		// 尝试消耗不存在的配额
		success, err := checker.CheckAndConsume(context.Background(), "non_existent_tenant", entity.ResourceTypeBots, 1)
		assert.Error(t, err)
		assert.False(t, success)
		assert.Contains(t, err.Error(), "quota not found")
	})
}

// TestQuotaChecker_CheckOnly_Complete 测试CheckOnly完整场景
func TestQuotaChecker_CheckOnly_Complete(t *testing.T) {
	t.Run("配额充足", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    3,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleDaily,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		// Mock返回
		mockRepo.On("GetByTenantAndResource", mock.Anything, "test_tenant", entity.ResourceTypeBots).Return(quota, nil)

		checker := NewQuotaChecker(db, mockRepo)

		available, status, err := checker.CheckOnly(context.Background(), "test_tenant", entity.ResourceTypeBots, 2)
		assert.NoError(t, err)
		assert.True(t, available)
		assert.NotNil(t, status)
		assert.Equal(t, entity.ResourceTypeBots, status.ResourceType)
		assert.Equal(t, 3, status.UsedCount)
		assert.Equal(t, 10, status.MaxLimit)
		assert.Equal(t, 7, status.RemainingCount)
		assert.InDelta(t, 30.0, status.UsagePercentage, 0.1)
	})

	t.Run("配额不足", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    9,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleDaily,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		mockRepo.On("GetByTenantAndResource", mock.Anything, "test_tenant", entity.ResourceTypeBots).Return(quota, nil)

		checker := NewQuotaChecker(db, mockRepo)

		available, status, err := checker.CheckOnly(context.Background(), "test_tenant", entity.ResourceTypeBots, 2)
		assert.NoError(t, err)
		assert.False(t, available)
		assert.NotNil(t, status)
		assert.Equal(t, 9, status.UsedCount)
		assert.Equal(t, 1, status.RemainingCount)
	})

	t.Run("配额刚好等于需求", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    8,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleDaily,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		mockRepo.On("GetByTenantAndResource", mock.Anything, "test_tenant", entity.ResourceTypeBots).Return(quota, nil)

		checker := NewQuotaChecker(db, mockRepo)

		available, status, err := checker.CheckOnly(context.Background(), "test_tenant", entity.ResourceTypeBots, 2)
		assert.NoError(t, err)
		assert.True(t, available)
		assert.Equal(t, 2, status.RemainingCount)
	})
}

// TestQuotaChecker_RollbackConsumption_BoundaryConditions 测试回滚边界条件
func TestQuotaChecker_RollbackConsumption_BoundaryConditions(t *testing.T) {
	t.Run("正常回滚", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    5,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		checker := NewQuotaChecker(db, mockRepo)

		err := checker.RollbackConsumption(context.Background(), "test_tenant", entity.ResourceTypeBots, 2)
		assert.NoError(t, err)

		// 验证计数减少
		var updatedQuota entity.Quota
		err = db.Where("tenant_id = ? AND resource_type = ?", "test_tenant", entity.ResourceTypeBots).First(&updatedQuota).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, updatedQuota.UsedCount)
	})

	t.Run("回滚到负数时应设为0", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    1,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		checker := NewQuotaChecker(db, mockRepo)

		// 回滚5个，但只有1个已用
		err := checker.RollbackConsumption(context.Background(), "test_tenant", entity.ResourceTypeBots, 5)
		assert.NoError(t, err)

		// 验证计数为0（不是负数）
		var updatedQuota entity.Quota
		err = db.Where("tenant_id = ? AND resource_type = ?", "test_tenant", entity.ResourceTypeBots).First(&updatedQuota).Error
		assert.NoError(t, err)
		assert.Equal(t, 0, updatedQuota.UsedCount)
	})

	t.Run("回滚全部配额", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    10,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}
		db.Create(quota)

		checker := NewQuotaChecker(db, mockRepo)

		err := checker.RollbackConsumption(context.Background(), "test_tenant", entity.ResourceTypeBots, 10)
		assert.NoError(t, err)

		// 验证计数为0
		var updatedQuota entity.Quota
		err = db.Where("tenant_id = ? AND resource_type = ?", "test_tenant", entity.ResourceTypeBots).First(&updatedQuota).Error
		assert.NoError(t, err)
		assert.Equal(t, 0, updatedQuota.UsedCount)
	})
}

// TestShouldResetQuota_AllCycles 测试所有重置周期
func TestShouldResetQuota_AllCycles(t *testing.T) {
	now := time.Now()

	t.Run("每日重置-同一天", func(t *testing.T) {
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleDaily,
			LastResetAt: now.UnixMilli(),
		}
		assert.False(t, shouldResetQuota(quota))
	})

	t.Run("每日重置-跨天", func(t *testing.T) {
		yesterday := now.Add(-24 * time.Hour)
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleDaily,
			LastResetAt: yesterday.UnixMilli(),
		}
		assert.True(t, shouldResetQuota(quota))
	})

	t.Run("每周重置-同一周", func(t *testing.T) {
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleWeekly,
			LastResetAt: now.UnixMilli(),
		}
		assert.False(t, shouldResetQuota(quota))
	})

	t.Run("每周重置-跨周", func(t *testing.T) {
		lastWeek := now.Add(-7 * 24 * time.Hour)
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleWeekly,
			LastResetAt: lastWeek.UnixMilli(),
		}
		assert.True(t, shouldResetQuota(quota))
	})

	t.Run("每月重置-同一个月", func(t *testing.T) {
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleMonthly,
			LastResetAt: now.UnixMilli(),
		}
		assert.False(t, shouldResetQuota(quota))
	})

	t.Run("每月重置-跨月", func(t *testing.T) {
		lastMonth := now.AddDate(0, -1, 0)
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleMonthly,
			LastResetAt: lastMonth.UnixMilli(),
		}
		assert.True(t, shouldResetQuota(quota))
	})

	t.Run("每年重置-同一年", func(t *testing.T) {
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleYearly,
			LastResetAt: now.UnixMilli(),
		}
		assert.False(t, shouldResetQuota(quota))
	})

	t.Run("每年重置-跨年", func(t *testing.T) {
		lastYear := now.AddDate(-1, 0, 0)
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleYearly,
			LastResetAt: lastYear.UnixMilli(),
		}
		assert.True(t, shouldResetQuota(quota))
	})

	t.Run("从不重置", func(t *testing.T) {
		oldTime := now.Add(-365 * 24 * time.Hour)
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleNever,
			LastResetAt: oldTime.UnixMilli(),
		}
		assert.False(t, shouldResetQuota(quota))
	})
}

// TestDetermineResourceRequirement_AllMethods 测试所有HTTP方法和路径组合
func TestDetermineResourceRequirement_AllMethods(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedType   entity.ResourceType
		expectedCount  int
	}{
		{
			name:          "POST创建Bot",
			method:        "POST",
			path:          "/api/v1/bots",
			expectedType:  entity.ResourceTypeBots,
			expectedCount: 1,
		},
		{
			name:          "POST创建对话",
			method:        "POST",
			path:          "/api/v1/conversations",
			expectedType:  entity.ResourceTypeMessages,
			expectedCount: 1,
		},
		{
			name:          "GET请求不消耗配额",
			method:        "GET",
			path:          "/api/v1/bots",
			expectedType:  "",
			expectedCount: 0,
		},
		{
			name:          "PUT请求不消耗配额",
			method:        "PUT",
			path:          "/api/v1/bots/123",
			expectedType:  "",
			expectedCount: 0,
		},
		{
			name:          "DELETE请求不消耗配额",
			method:        "DELETE",
			path:          "/api/v1/bots/123",
			expectedType:  "",
			expectedCount: 0,
		},
		{
			name:          "未知POST路径",
			method:        "POST",
			path:          "/api/v1/unknown",
			expectedType:  "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟RequestContext
			ctx := &mockRequestContext{
				method: tt.method,
				path:   tt.path,
			}

			resourceType, count := determineResourceRequirement(ctx)
			assert.Equal(t, tt.expectedType, resourceType)
			assert.Equal(t, tt.expectedCount, count)
		})
	}
}

// TestCheckQuotaAvailability_ErrorScenarios 测试错误场景
func TestCheckQuotaAvailability_ErrorScenarios(t *testing.T) {
	t.Run("配额不存在时返回默认值", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		// Mock返回nil（配额不存在）
		mockRepo.On("GetByTenantAndResource", mock.Anything, "test_tenant", entity.ResourceTypeBots).
			Return((*entity.Quota)(nil), nil)

		available, quota, err := checkQuotaAvailability(context.Background(), db, mockRepo, "test_tenant", entity.ResourceTypeBots)
		assert.NoError(t, err)
		assert.False(t, available)
		assert.NotNil(t, quota)
		assert.Equal(t, "test_tenant", quota.TenantID)
		assert.Equal(t, entity.ResourceTypeBots, quota.ResourceType)
		assert.Equal(t, 0, quota.UsedCount)
		assert.Equal(t, 0, quota.MaxLimit)
	})

	t.Run("数据库错误", func(t *testing.T) {
		db := setupTestDB(t)
		mockRepo := new(MockQuotaRepository)

		// Mock返回错误
		mockRepo.On("GetByTenantAndResource", mock.Anything, "test_tenant", entity.ResourceTypeBots).
			Return((*entity.Quota)(nil), syscall.ECONNREFUSED)

		available, quota, err := checkQuotaAvailability(context.Background(), db, mockRepo, "test_tenant", entity.ResourceTypeBots)
		assert.Error(t, err)
		assert.False(t, available)
		assert.Nil(t, quota)
	})
}

// mockRequestContext 模拟RequestContext
type mockRequestContext struct {
	method string
	path   string
}

func (m *mockRequestContext) Method() []byte {
	return []byte(m.method)
}

func (m *mockRequestContext) Path() []byte {
	return []byte(m.path)
}

// setupTestDB 设置测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	// 这里应该返回一个内存数据库或测试数据库连接
	// 由于需要真实数据库连接，这里使用mock
	// 实际实现中应该使用testify/mock或sqlmock
	return nil
}

// 注意：以上测试需要真实的数据库连接才能运行
// 在CI/CD环境中应该使用测试数据库
// 或者使用sqlmock进行Mock
