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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
)

// MockTenantRepository 模拟租户仓储
type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) GetByID(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetByName(ctx context.Context, name string) (*entity.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetBySubdomain(ctx context.Context, subdomain string) (*entity.Tenant, error) {
	args := m.Called(ctx, subdomain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Update(ctx context.Context, tenant *entity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) Delete(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockTenantRepository) List(ctx context.Context, filter *repository.TenantFilter) ([]*entity.Tenant, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.Tenant), args.Get(1).(int64), args.Error(2)
}

func (m *MockTenantRepository) UpdateStatus(ctx context.Context, tenantID string, status entity.TenantStatus) error {
	args := m.Called(ctx, tenantID, status)
	return args.Error(0)
}

func (m *MockTenantRepository) UpdateSubscriptionTier(ctx context.Context, tenantID string, tier entity.SubscriptionTier) error {
	args := m.Called(ctx, tenantID, tier)
	return args.Error(0)
}

// 其他必需的方法...

func TestTenantValidationService_CheckCompanyNameAvailability(t *testing.T) {
	t.Run("名称可用", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		// Mock: 数据库中不存在该名称
		mockRepo.On("GetByName", ctx, "测试公司").Return(nil, repository.ErrRecordNotFound)

		available, suggestions, err := service.CheckCompanyNameAvailability(ctx, "测试公司")

		assert.NoError(t, err)
		assert.True(t, available)
		assert.Empty(t, suggestions)
		mockRepo.AssertExpectations(t)
	})

	t.Run("名称已被占用", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		existingTenant := &entity.Tenant{
			TenantName: "测试公司",
		}

		// Mock: 数据库中已存在该名称
		mockRepo.On("GetByName", ctx, "测试公司").Return(existingTenant, nil)

		available, suggestions, err := service.CheckCompanyNameAvailability(ctx, "测试公司")

		assert.NoError(t, err)
		assert.False(t, available)
		assert.NotEmpty(t, suggestions)
		assert.True(t, len(suggestions) > 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("名称长度验证", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		// 名称太短
		_, _, err := service.CheckCompanyNameAvailability(ctx, "a")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "长度必须在2-200个字符之间")

		// 名称太长
		longName := string(make([]byte, 201))
		_, _, err = service.CheckCompanyNameAvailability(ctx, longName)
		assert.Error(t, err)
	})
}

func TestTenantValidationService_CheckSubdomainAvailability(t *testing.T) {
	t.Run("子域名可用", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		// Mock: 数据库中不存在该子域名
		mockRepo.On("GetBySubdomain", ctx, "testcompany").Return(nil, repository.ErrRecordNotFound)

		available, suggestions, err := service.CheckSubdomainAvailability(ctx, "testcompany")

		assert.NoError(t, err)
		assert.True(t, available)
		assert.Empty(t, suggestions)
		mockRepo.AssertExpectations(t)
	})

	t.Run("子域名已被占用", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		existingTenant := &entity.Tenant{
			Subdomain: "testcompany",
		}

		// Mock: 数据库中已存在该子域名
		mockRepo.On("GetBySubdomain", ctx, "testcompany").Return(existingTenant, nil)

		available, suggestions, err := service.CheckSubdomainAvailability(ctx, "testcompany")

		assert.NoError(t, err)
		assert.False(t, available)
		assert.NotEmpty(t, suggestions)
		assert.True(t, len(suggestions) > 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("子域名格式验证-长度", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		// 子域名太短
		_, _, err := service.CheckSubdomainAvailability(ctx, "ab")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "长度必须在3-63个字符之间")

		// 子域名太长
		longSubdomain := string(make([]byte, 64))
		_, _, err = service.CheckSubdomainAvailability(ctx, longSubdomain)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "长度必须在3-63个字符之间")
	})

	t.Run("子域名格式验证-连字符", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		// 以连字符开头
		_, _, err := service.CheckSubdomainAvailability(ctx, "-test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能以连字符开头或结尾")

		// 以连字符结尾
		_, _, err = service.CheckSubdomainAvailability(ctx, "test-")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不能以连字符开头或结尾")
	})

	t.Run("子域名转换小写", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		service := NewTenantValidationService(mockRepo)
		ctx := context.Background()

		// Mock: 数据库中不存在该子域名
		mockRepo.On("GetBySubdomain", ctx, "testcompany").Return(nil, repository.ErrRecordNotFound)

		// 输入大写字母，应该被转换为小写
		available, _, err := service.CheckSubdomainAvailability(ctx, "TESTCOMPANY")

		assert.NoError(t, err)
		assert.True(t, available)
		mockRepo.AssertExpectations(t)
	})
}

func TestTenantValidationService_GenerateCompanyNameSuggestions(t *testing.T) {
	mockRepo := new(MockTenantRepository)
	service := NewTenantValidationService(mockRepo)

	t.Run("生成建议数量", func(t *testing.T) {
		suggestions := service.generateCompanyNameSuggestions("测试公司")

		assert.NotEmpty(t, suggestions)
		assert.True(t, len(suggestions) >= 3, "应该至少生成3个建议")
		assert.True(t, len(suggestions) <= 8, "建议数量不应超过8个")
	})

	t.Run("建议内容多样性", func(t *testing.T) {
		suggestions := service.generateCompanyNameSuggestions("测试公司")

		// 检查是否有不同类型的后缀
		hasSuffix := false
		for _, suggestion := range suggestions {
			if len(suggestion) > len("测试公司") {
				suffix := suggestion[len("测试公司"):]
				if suffix != "" {
					hasSuffix = true
					break
				}
			}
		}
		assert.True(t, hasSuffix, "建议应该包含不同的后缀")
	})
}

func TestTenantValidationService_GenerateSubdomainSuggestions(t *testing.T) {
	mockRepo := new(MockTenantRepository)
	service := NewTenantValidationService(mockRepo)

	t.Run("生成建议数量", func(t *testing.T) {
		suggestions := service.generateSubdomainSuggestions("testcompany")

		assert.NotEmpty(t, suggestions)
		assert.True(t, len(suggestions) >= 5, "应该至少生成5个建议")
		assert.True(t, len(suggestions) <= 8, "建议数量不应超过8个")
	})

	t.Run("建议内容多样性", func(t *testing.T) {
		suggestions := service.generateSubdomainSuggestions("testcompany")

		// 检查是否有数字后缀
		hasNumberSuffix := false
		for _, suggestion := range suggestions {
			// 检查是否以数字结尾
			if len(suggestion) > 0 && suggestion[len(suggestion)-1] >= '0' && suggestion[len(suggestion)-1] <= '9' {
				hasNumberSuffix = true
				break
			}
		}

		// 检查是否有行业后缀
		hasIndustrySuffix := false
		industrySuffixes := []string{"-tech", "-cloud", "-ai", "-data", "-sys"}
		for _, suggestion := range suggestions {
			for _, suffix := range industrySuffixes {
				if len(suggestion) > len(suffix) && suggestion[len(suggestion)-len(suffix):] == suffix {
					hasIndustrySuffix = true
					break
				}
			}
			if hasIndustrySuffix {
				break
			}
		}

		assert.True(t, hasNumberSuffix || hasIndustrySuffix, "建议应该包含数字或行业后缀")
	})
}

func TestTenantValidationService_CleanCompanyName(t *testing.T) {
	mockRepo := new(MockTenantRepository)
	service := NewTenantValidationService(mockRepo)

	t.Run("去除首尾空格", func(t *testing.T) {
		cleaned := service.cleanCompanyName("  测试公司  ")
		assert.Equal(t, "测试公司", cleaned)
	})

	t.Run("去除多余空格", func(t *testing.T) {
		cleaned := service.cleanCompanyName("测试   公司")
		assert.Equal(t, "测试 公司", cleaned)
	})

	t.Run("去除不可见字符", func(t *testing.T) {
		cleaned := service.cleanCompanyName("测试\x00公司")
		assert.Equal(t, "测试公司", cleaned)
	})
}
