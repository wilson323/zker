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
)

// MockPermissionChecker 模拟权限检查器
type MockPermissionChecker struct {
	mock.Mock
}

func (m *MockPermissionChecker) CheckDataPermission(ctx context.Context, tenantID, userID string, resourceType entity.ResourceType, action, resourceID string) (bool, error) {
	args := m.Called(ctx, tenantID, userID, resourceType, action, resourceID)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionChecker) CheckFieldPermission(ctx context.Context, tenantID, userID, resourceType string, fields []string) (map[string]string, error) {
	args := m.Called(ctx, tenantID, userID, resourceType, fields)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockPermissionChecker) UserHasRole(ctx context.Context, tenantID, userID, roleCode string) (bool, error) {
	args := m.Called(ctx, tenantID, userID, roleCode)
	if args.Get(0) == nil {
		return false, args.Error(1)
	}
	return args.Bool(0), args.Error(1)
}

// ================================================================================
// RequirePermission 测试
// ================================================================================

func TestRequirePermission_Success(t *testing.T) {
	mockChecker := new(MockPermissionChecker)

	// 初始化权限检查器
	permissionChecker = mockChecker

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType:  entity.ResourceTypeBots,
		Action:        "read",
		GetResourceID: func(c *app.RequestContext) string { return "bot_123" },
	}

	// 设置期望：权限检查通过
	mockChecker.On("CheckDataPermission", ctx, "tenant_123", "user_123", entity.ResourceTypeBots, "read", "bot_123").Return(true, nil)

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该继续执行（没有被Abort）
	assert.False(t, c.Response.StatusCode() >= 400)
	mockChecker.AssertExpectations(t)
}

func TestRequirePermission_PermissionDenied(t *testing.T) {
	mockChecker := new(MockPermissionChecker)

	// 初始化权限检查器
	permissionChecker = mockChecker

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType:  entity.ResourceTypeBots,
		Action:        "read",
		GetResourceID: func(c *app.RequestContext) string { return "bot_123" },
	}

	// 设置期望：权限拒绝
	mockChecker.On("CheckDataPermission", ctx, "tenant_123", "user_123", entity.ResourceTypeBots, "read", "bot_123").Return(false, nil)

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该返回403
	assert.Equal(t, consts.StatusForbidden, c.Response.StatusCode())
	mockChecker.AssertExpectations(t)
}

func TestRequirePermission_MissingUserID(t *testing.T) {
	mockChecker := new(MockPermissionChecker)

	// 初始化权限检查器
	permissionChecker = mockChecker

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
	mockChecker := new(MockPermissionChecker)

	// 初始化权限检查器
	permissionChecker = mockChecker

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

func TestRequirePermission_CheckerError(t *testing.T) {
	mockChecker := new(MockPermissionChecker)

	// 初始化权限检查器
	permissionChecker = mockChecker

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType:  entity.ResourceTypeBots,
		Action:        "read",
		GetResourceID: func(c *app.RequestContext) string { return "bot_123" },
	}

	// 设置期望：检查器返回错误
	checkErr := errors.New("database connection failed")
	mockChecker.On("CheckDataPermission", ctx, "tenant_123", "user_123", entity.ResourceTypeBots, "read", "bot_123").Return(false, checkErr)

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该返回500
	assert.Equal(t, consts.StatusInternalServerError, c.Response.StatusCode())
	mockChecker.AssertExpectations(t)
}

func TestRequirePermission_NoChecker(t *testing.T) {
	// 不初始化权限检查器
	permissionChecker = nil

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType: entity.ResourceTypeBots,
		Action:       "read",
	}

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该跳过检查并继续执行（没有Abort）
	assert.False(t, c.Response.StatusCode() >= 400)
}

func TestRequirePermission_NoResourceID(t *testing.T) {
	mockChecker := new(MockPermissionChecker)

	// 初始化权限检查器
	permissionChecker = mockChecker

	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodPost, "/api/v1/bots")
	c.Request.Header.Set("X-User-ID", "user_123")
	c.Request.Header.Set("X-Tenant-ID", "tenant_123")

	config := PermissionCheckConfig{
		ResourceType:  entity.ResourceTypeBots,
		Action:        "create",
		GetResourceID: nil, // 没有GetResourceID函数
	}

	// 设置期望：resourceID为空字符串
	mockChecker.On("CheckDataPermission", ctx, "tenant_123", "user_123", entity.ResourceTypeBots, "create", "").Return(true, nil)

	middleware := RequirePermission(config)
	middleware(ctx, c)

	// 验证：应该继续执行
	assert.False(t, c.Response.StatusCode() >= 400)
	mockChecker.AssertExpectations(t)
}

// ================================================================================
// ParseResourceIDFromPath 测试
// ================================================================================

func TestParseResourceIDFromPath_Success(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	c.Params.Set("bot_id", "bot_123")

	getter := ParseResourceIDFromPath("bot_id")
	resourceID := getter(c)

	assert.Equal(t, "bot_123", resourceID)
}

func TestParseResourceIDFromPath_NotFound(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots/bot_123")
	// 不设置参数

	getter := ParseResourceIDFromPath("bot_id")
	resourceID := getter(c)

	assert.Equal(t, "", resourceID)
}

// ================================================================================
// ParseResourceIDFromQuery 测试
// ================================================================================

func TestParseResourceIDFromQuery_Success(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots?bot_id=bot_123")
	c.QueryArgs().Add("bot_id", "bot_123")

	getter := ParseResourceIDFromQuery("bot_id")
	resourceID := getter(c)

	assert.Equal(t, "bot_123", resourceID)
}

func TestParseResourceIDFromQuery_NotFound(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots")
	// 不设置查询参数

	getter := ParseResourceIDFromQuery("bot_id")
	resourceID := getter(c)

	assert.Equal(t, "", resourceID)
}

func TestParseResourceIDFromQuery_MultipleValues(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots?bot_id=bot_123&bot_id=bot_456")
	c.QueryArgs().Add("bot_id", "bot_123")
	c.QueryArgs().Add("bot_id", "bot_456")

	getter := ParseResourceIDFromQuery("bot_id")
	resourceID := getter(c)

	// 应该返回第一个值
	assert.Equal(t, "bot_123", resourceID)
}

// ================================================================================
// ParseResourceIDFromHeader 测试
// ================================================================================

func TestParseResourceIDFromHeader_Success(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots")
	c.Request.Header.Set("X-Bot-ID", "bot_123")

	getter := ParseResourceIDFromHeader("X-Bot-ID")
	resourceID := getter(c)

	assert.Equal(t, "bot_123", resourceID)
}

func TestParseResourceIDFromHeader_NotFound(t *testing.T) {
	ctx := context.Background()
	c := app.NewRequestContext(ctx, consts.MethodGet, "/api/v1/bots")
	// 不设置header

	getter := ParseResourceIDFromHeader("X-Bot-ID")
	resourceID := getter(c)

	assert.Equal(t, "", resourceID)
}

// ================================================================================
// InitPermissionMiddleware 测试
// ================================================================================

func TestInitPermissionMiddleware(t *testing.T) {
	mockChecker := new(MockPermissionChecker)

	// 调用初始化函数
	InitPermissionMiddleware(mockChecker)

	// 验证：permissionChecker应该被设置
	assert.Equal(t, mockChecker, permissionChecker)

	// 清理
	permissionChecker = nil
}
