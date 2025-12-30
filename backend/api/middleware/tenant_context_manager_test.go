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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	tenantentity "github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	tenantrepo "github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
)

// MockTenantRepository Mock租户仓储
type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) Create(ctx context.Context, tenant *tenantentity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) GetByID(ctx context.Context, tenantID string) (*tenantentity.Tenant, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenantentity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetBySubdomain(ctx context.Context, subdomain string) (*tenantentity.Tenant, error) {
	args := m.Called(ctx, subdomain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenantentity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetByName(ctx context.Context, name string) (*tenantentity.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenantentity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Update(ctx context.Context, tenant *tenantentity.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepository) Delete(ctx context.Context, tenantID string) error {
	args := m.Called(ctx, tenantID)
	return args.Error(0)
}

func (m *MockTenantRepository) List(ctx context.Context, filter interface{}) ([]*tenantentity.Tenant, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*tenantentity.Tenant), args.Get(1).(int64), args.Error(2)
}

// TestNewTenantContextManager 测试创建租户上下文管理器
func TestNewTenantContextManager(t *testing.T) {
	mockRepo := new(MockTenantRepository)
	baseDomains := []string{"saas.coze.com", "app.example.com"}

	manager := NewTenantContextManager(mockRepo, baseDomains)

	assert.NotNil(t, manager)
}

// TestIdentifyFromSubdomain 测试从子域名识别租户
func TestIdentifyFromSubdomain(t *testing.T) {
	t.Run("成功识别子域名", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		expectedTenant := &tenantentity.Tenant{
			TenantID:   "tenant_123",
			TenantName: "测试租户",
			Subdomain:  "acme",
		}

		mockRepo.On("GetBySubdomain", ctx, "acme").Return(expectedTenant, nil).Once()

		manager := NewTenantContextManager(mockRepo, []string{"saas.coze.com"})

		tenant, err := manager.IdentifyFromSubdomain(ctx, "acme.saas.coze.com")

		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "tenant_123", tenant.TenantID)
		assert.Equal(t, "acme", tenant.Subdomain)

		mockRepo.AssertExpectations(t)
	})

	t.Run("子域名不存在", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		mockRepo.On("GetBySubdomain", ctx, "nonexistent").Return(nil, fmt.Errorf("record not found")).Once()

		manager := NewTenantContextManager(mockRepo, []string{"saas.coze.com"})

		tenant, err := manager.IdentifyFromSubdomain(ctx, "nonexistent.saas.coze.com")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "查询租户失败")

		mockRepo.AssertExpectations(t)
	})

	t.Run("无效的子域名格式", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		manager := NewTenantContextManager(mockRepo, []string{"saas.coze.com"})

		// 测试基础域名（没有子域名）
		tenant, err := manager.IdentifyFromSubdomain(ctx, "saas.coze.com")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "无效的子域名")
	})

	t.Run("缓存命中", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		expectedTenant := &tenantentity.Tenant{
			TenantID:   "tenant_123",
			TenantName: "测试租户",
			Subdomain:  "acme",
		}

		// 第一次调用会查询数据库
		mockRepo.On("GetBySubdomain", ctx, "acme").Return(expectedTenant, nil).Once()

		manager := NewTenantContextManager(mockRepo, []string{"saas.coze.com"})

		// 第一次调用
		tenant1, err1 := manager.IdentifyFromSubdomain(ctx, "acme.saas.coze.com")
		assert.NoError(t, err1)
		assert.Equal(t, "tenant_123", tenant1.TenantID)

		// 第二次调用应该从缓存获取，不再调用repo
		tenant2, err2 := manager.IdentifyFromSubdomain(ctx, "acme.saas.coze.com")
		assert.NoError(t, err2)
		assert.Equal(t, "tenant_123", tenant2.TenantID)

		mockRepo.AssertExpectations(t) // 应该只被调用一次
	})

	t.Run("清除缓存", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		expectedTenant := &tenantentity.Tenant{
			TenantID:   "tenant_123",
			TenantName: "测试租户",
			Subdomain:  "acme",
		}

		mockRepo.On("GetBySubdomain", ctx, "acme").Return(expectedTenant, nil).Times(2)

		manager := NewTenantContextManager(mockRepo, []string{"saas.coze.com"})

		// 第一次调用
		tenant1, _ := manager.IdentifyFromSubdomain(ctx, "acme.saas.coze.com")
		assert.Equal(t, "tenant_123", tenant1.TenantID)

		// 清除缓存
		ClearTenantCache()

		// 第二次调用（缓存已清除，应该再次查询数据库）
		tenant2, _ := manager.IdentifyFromSubdomain(ctx, "acme.saas.coze.com")
		assert.Equal(t, "tenant_123", tenant2.TenantID)

		mockRepo.AssertExpectations(t) // 应该被调用两次
	})
}

// TestIdentifyFromPath 测试从路径识别租户
func TestIdentifyFromPath(t *testing.T) {
	t.Run("成功从路径识别租户", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		expectedTenant := &tenantentity.Tenant{
			TenantID:   "tenant_abc",
			TenantName: "路径租户",
		}

		mockRepo.On("GetByID", ctx, "tenant-abc").Return(expectedTenant, nil).Once()

		manager := NewTenantContextManager(mockRepo, nil)

		tenant, err := manager.IdentifyFromPath(ctx, "/tenant-abc/bots")

		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "tenant_abc", tenant.TenantID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("路径格式无效", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		manager := NewTenantContextManager(mockRepo, nil)

		// 没有租户标识的路径
		tenant, err := manager.IdentifyFromPath(ctx, "/api/v1/bots")

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "无效的路径格式")
	})

	t.Run("租户不存在", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, "nonexistent").Return(nil, fmt.Errorf("not found")).Once()

		manager := NewTenantContextManager(mockRepo, nil)

		tenant, err := manager.IdentifyFromPath(ctx, "/nonexistent/bots")

		assert.Error(t, err)
		assert.Nil(t, tenant)

		mockRepo.AssertExpectations(t)
	})
}

// TestIdentifyFromHeader 测试从请求头识别租户
func TestIdentifyFromHeader(t *testing.T) {
	t.Run("成功从X-Tenant-ID识别", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		expectedTenant := &tenantentity.Tenant{
			TenantID:   "tenant_xyz",
			TenantName: "Header租户",
		}

		mockRepo.On("GetByID", ctx, "tenant-xyz").Return(expectedTenant, nil).Once()

		manager := NewTenantContextManager(mockRepo, nil)

		headers := map[string]string{
			"X-Tenant-ID": "tenant-xyz",
		}

		tenant, err := manager.IdentifyFromHeader(ctx, headers)

		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "tenant_xyz", tenant.TenantID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("兼容X-Tenant-Id", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		expectedTenant := &tenantentity.Tenant{
			TenantID:   "tenant_xyz",
			TenantName: "Header租户",
		}

		mockRepo.On("GetByID", ctx, "tenant-xyz").Return(expectedTenant, nil).Once()

		manager := NewTenantContextManager(mockRepo, nil)

		headers := map[string]string{
			"X-Tenant-Id": "tenant-xyz", // 旧格式
		}

		tenant, err := manager.IdentifyFromHeader(ctx, headers)

		assert.NoError(t, err)
		assert.NotNil(t, tenant)
		assert.Equal(t, "tenant_xyz", tenant.TenantID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("缺少租户Header", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		ctx := context.Background()

		manager := NewTenantContextManager(mockRepo, nil)

		headers := map[string]string{} // 空headers

		tenant, err := manager.IdentifyFromHeader(ctx, headers)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "缺少租户标识 Header")
	})
}

// TestSetAndGetTenantContext 测试设置和获取租户上下文
func TestSetAndGetTenantContext(t *testing.T) {
	t.Run("成功设置和获取", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		manager := NewTenantContextManager(mockRepo, nil)

		ctx := context.Background()
		tenant := &tenantentity.Tenant{
			TenantID:   "test_tenant",
			TenantName: "测试租户",
		}

		// 设置租户上下文
		newCtx := manager.SetTenantContext(ctx, tenant)

		// 获取租户上下文
		retrievedTenant, err := manager.GetTenantContext(newCtx)

		assert.NoError(t, err)
		assert.NotNil(t, retrievedTenant)
		assert.Equal(t, "test_tenant", retrievedTenant.TenantID)
	})

	t.Run("上下文中没有租户", func(t *testing.T) {
		mockRepo := new(MockTenantRepository)
		manager := NewTenantContextManager(mockRepo, nil)

		ctx := context.Background()

		// 尝试从不存在的上下文获取
		tenant, err := manager.GetTenantContext(ctx)

		assert.Error(t, err)
		assert.Nil(t, tenant)
		assert.Contains(t, err.Error(), "上下文中缺少租户信息")
	})
}

// TestIsIPAddress 测试IP地址验证
func TestIsIPAddress(t *testing.T) {
	tests := []struct {
		host     string
		expected bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"127.0.0.1", true},
		{"255.255.255.255", true},
		{"192.168.1", false},     // 只有3部分
		{"192.168.1.1.1", false}, // 5部分
		{"example.com", false},
		{"localhost", false},
		{"256.1.1.1", false}, // 无效IP
	}

	for _, tt := range tests {
		t.Run(tt.host, func(t *testing.T) {
			result := isIPAddress(tt.host)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestExtractSubdomain 测试子域名提取
func TestExtractSubdomain(t *testing.T) {
	mockRepo := new(MockTenantRepository)

	tests := []struct {
		name         string
		host         string
		baseDomains  []string
		expectedSub  string
	}{
		{
			name:         "简单子域名",
			host:         "acme.saas.coze.com",
			baseDomains:  []string{"saas.coze.com"},
			expectedSub:  "acme",
		},
		{
			name:         "带端口号",
			host:         "tenant-123.example.com:8080",
			baseDomains:  []string{},
			expectedSub:  "tenant-123",
		},
		{
			name:         "IP地址（跳过）",
			host:         "192.168.1.1",
			baseDomains:  []string{},
			expectedSub:  "",
		},
		{
			name:         "基础域名",
			host:         "saas.coze.com",
			baseDomains:  []string{"saas.coze.com"},
			expectedSub:  "",
		},
		{
			name:         "多级域名",
			host:         "my.app.saas.coze.com",
			baseDomains:  []string{},
			expectedSub:  "my",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewTenantContextManager(mockRepo, tt.baseDomains)
			subdomain := manager.(*tenantContextManagerImpl).extractSubdomain(tt.host)
			assert.Equal(t, tt.expectedSub, subdomain)
		})
	}
}
