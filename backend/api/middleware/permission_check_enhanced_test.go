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
	"errors"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// MockPermissionCheckerService 模拟权限检查服务
type MockPermissionCheckerService struct {
	mock.Mock
}

func (m *MockPermissionCheckerService) CheckDataPermission(ctx context.Context, userID, tenantID string, resourceType entity.ResourceType, resourceID string) error {
	args := m.Called(ctx, userID, tenantID, resourceType, resourceID)
	return args.Error(0)
}

func (m *MockPermissionCheckerService) GetFieldPermissions(ctx context.Context, userID, tenantID string, resourceType string) (map[string]string, error) {
	args := m.Called(ctx, userID, tenantID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

// MockUserRoleRepository 模拟用户角色仓储
type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) GetRolesByUser(ctx context.Context, userID, tenantID string) ([]*entity.Role, error) {
	args := m.Called(ctx, userID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Role), args.Error(1)
}

// ================================================================================
// RequirePermission 测试
// ================================================================================

func TestRequirePermission_Success(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType: entity.ResourceTypeBots,
		Action:       "read",
		GetResourceID: func(c *app.RequestContext) string {
			return "bot_123"
		},
	}

	// 设置期望：权限检查通过
	mockSvc.On("CheckDataPermission", ctx, "user_123", "tenant_123", entity.ResourceTypeBots, "bot_123").Return(nil)

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该继续执行（没有被Abort）
	assert.False(t, c.Response.StatusCode() >= 400)
	mockSvc.AssertExpectations(t)
}

func TestRequirePermission_PermissionDenied(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType: entity.ResourceTypeBots,
		Action:       "read",
		GetResourceID: func(c *app.RequestContext) string {
			return "bot_123"
		},
	}

	// 设置期望：权限拒绝
	permissionDeniedErr := errors.New("permission denied: user has no roles")
	mockSvc.On("CheckDataPermission", ctx, "user_123", "tenant_123", entity.ResourceTypeBots, "bot_123").Return(permissionDeniedErr)

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该返回403
	assert.Equal(t, consts.StatusForbidden, c.Response.StatusCode())

	var responseBody map[string]interface{}
	err := c.JSON.Bind(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, berrno.ErrDataPermissionDeniedCode, responseBody["code"])

	mockSvc.AssertExpectations(t)
}

func TestRequirePermission_MissingUserID(t *testing.T) {
	enhancedPermissionChecker = nil // 不需要权限检查器

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	// 不设置 X-User-ID

	config := PermissionCheckConfig{
		ResourceType: entity.ResourceTypeBots,
		Action:       "read",
	}

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该返回401
	assert.Equal(t, consts.StatusUnauthorized, c.Response.StatusCode())
}

func TestRequirePermission_MissingTenantID(t *testing.T) {
	enhancedPermissionChecker = nil

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	// 不设置 X-Tenant-ID

	config := PermissionCheckConfig{
		ResourceType: entity.ResourceTypeBots,
		Action:       "read",
	}

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该返回400
	assert.Equal(t, consts.StatusBadRequest, c.Response.StatusCode())
}

// ================================================================================
// RequireRole 测试
// ================================================================================

func TestRequireRole_Success(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodPost, "/api/v1/admin/users")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	// 设置期望：用户有admin角色
	roles := []*entity.Role{
		{
			RoleID:   "role_123",
			RoleCode: "admin",
		},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	middleware := RequireRole("admin")
	middleware(ctx, c)

	// 验证：应该继续执行
	assert.False(t, c.Response.StatusCode() >= 400)
	mockRoleRepo.AssertExpectations(t)
}

func TestRequireRole_RoleRequired(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodPost, "/api/v1/admin/users")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	// 设置期望：用户没有admin角色（有其他角色）
	roles := []*entity.Role{
		{
			RoleID:   "role_456",
			RoleCode: "user",
		},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	middleware := RequireRole("admin")
	middleware(ctx, c)

	// 验证：应该返回403
	assert.Equal(t, consts.StatusForbidden, c.Response.StatusCode())

	var responseBody map[string]interface{}
	err := c.JSON.Bind(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, berrno.ErrPermissionDeniedCode, responseBody["code"])
	assert.Equal(t, "admin", responseBody["required_role"])

	mockRoleRepo.AssertExpectations(t)
}

// ================================================================================
// RequireAnyRole 测试
// ================================================================================

func TestRequireAnyRole_Success(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodPost, "/api/v1/admin/users")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	// 设置期望：用户有super_admin角色
	roles := []*entity.Role{
		{
			RoleID:   "role_123",
			RoleCode: "super_admin",
		},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	middleware := RequireAnyRole("admin", "super_admin")
	middleware(ctx, c)

	// 验证：应该继续执行
	assert.False(t, c.Response.StatusCode() >= 400)
	mockRoleRepo.AssertExpectations(t)
}

func TestRequireAnyRole_NoRoles(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodPost, "/api/v1/admin/users")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	// 设置期望：用户没有任一角色
	roles := []*entity.Role{
		{
			RoleID:   "role_789",
			RoleCode: "guest",
		},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	middleware := RequireAnyRole("admin", "super_admin")
	middleware(ctx, c)

	// 验证：应该返回403
	assert.Equal(t, consts.StatusForbidden, c.Response.StatusCode())

	var responseBody map[string]interface{}
	err := c.JSON.Bind(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, berrno.ErrPermissionDeniedCode, responseBody["code"])

	mockRoleRepo.AssertExpectations(t)
}

// ================================================================================
// RequireFieldPermission 测试
// ================================================================================

func TestRequireFieldPermission_Success(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	enhancedPermissionChecker = NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	// 设置期望：获取字段权限
	fieldPerms := map[string]string{
		"name":  "editable",
		"desc":  "editable",
		"token": "hidden",
	}
	mockSvc.On("GetFieldPermissions", ctx, "user_123", "tenant_123", "bots").Return(fieldPerms, nil)

	middleware := RequireFieldPermission("bots", []string{"id", "name"})
	middleware(ctx, c)

	// 验证：字段权限应该保存到上下文
	fieldPermsObj, exists := c.Get("field_permissions")
	assert.True(t, exists)
	assert.Equal(t, fieldPerms, fieldPermsObj)

	mockSvc.AssertExpectations(t)
}

// ================================================================================
// EnhancedPermissionChecker 测试
// ================================================================================

func TestEnhancedPermissionChecker_UserHasRole(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	checker := NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()

	// 测试：有角色
	roles := []*entity.Role{
		{RoleID: "role_123", RoleCode: "admin"},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	hasRole, err := checker.UserHasRole(ctx, "tenant_123", "user_123", "admin")
	assert.NoError(t, err)
	assert.True(t, hasRole)

	mockRoleRepo.AssertExpectations(t)
}

func TestEnhancedPermissionChecker_UserHasAnyRole(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	checker := NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()

	// 测试：有任一角色
	roles := []*entity.Role{
		{RoleID: "role_123", RoleCode: "admin"},
		{RoleID: "role_456", RoleCode: "user"},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	hasAnyRole, err := checker.UserHasAnyRole(ctx, "tenant_123", "user_123", "admin", "super_admin")
	assert.NoError(t, err)
	assert.True(t, hasAnyRole)

	mockRoleRepo.AssertExpectations(t)
}

func TestEnhancedPermissionChecker_GetUserPermissions(t *testing.T) {
	mockSvc := new(MockPermissionCheckerService)
	mockRoleRepo := new(MockUserRoleRepository)

	checker := NewEnhancedPermissionChecker(mockSvc, mockRoleRepo)

	ctx := context.Background()

	// 测试：获取用户权限
	roles := []*entity.Role{
		{RoleID: "role_123", RoleCode: "bot.read"},
		{RoleID: "role_456", RoleCode: "bot.write"},
	}
	mockRoleRepo.On("GetRolesByUser", ctx, "user_123", "tenant_123").Return(roles, nil)

	perms, err := checker.GetUserPermissions(ctx, "tenant_123", "user_123")
	assert.NoError(t, err)
	assert.Contains(t, perms, "bot.read")
	assert.Contains(t, perms, "bot.write")

	mockRoleRepo.AssertExpectations(t)
}

// ================================================================================
// 辅助函数测试
// ================================================================================

func TestParseResourceIDFromPath(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Params.Set("bot_id", "bot_123")

	getter := ParseResourceIDFromPath("bot_id")
	resourceID := getter(c)

	assert.Equal(t, "bot_123", resourceID)
}

func TestParseResourceIDFromQuery(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots?bot_id=bot_123")
	c.QueryArgs().Add("bot_id", "bot_123")

	getter := ParseResourceIDFromQuery("bot_id")
	resourceID := getter(c)

	assert.Equal(t, "bot_123", resourceID)
}
