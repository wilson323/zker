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

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// MockFieldPermissionRepository Mock字段权限仓储
type MockFieldPermissionRepository struct {
	mock.Mock
}

func (m *MockFieldPermissionRepository) GetByRoleAndResource(ctx context.Context, roleID string, resourceType string) ([]*entity.FieldPermission, error) {
	args := m.Called(ctx, roleID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FieldPermission), args.Error(1)
}

func (m *MockFieldPermissionRepository) Create(ctx context.Context, perm *entity.FieldPermission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockFieldPermissionRepository) GetByID(ctx context.Context, permissionID string) (*entity.FieldPermission, error) {
	args := m.Called(ctx, permissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.FieldPermission), args.Error(1)
}

func (m *MockFieldPermissionRepository) GetByRole(ctx context.Context, roleID string) ([]*entity.FieldPermission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FieldPermission), args.Error(1)
}

func (m *MockFieldPermissionRepository) Update(ctx context.Context, perm *entity.FieldPermission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockFieldPermissionRepository) Delete(ctx context.Context, permissionID string) error {
	args := m.Called(ctx, permissionID)
	return args.Error(0)
}

func (m *MockFieldPermissionRepository) DeleteByRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

// TestFieldPermissionChecker_CheckFieldPermission 测试检查字段权限
func TestFieldPermissionChecker_CheckFieldPermission(t *testing.T) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"
	field := "api_key"

	// Mock返回用户的角色和字段权限
	roles := []*entity.Role{
		{RoleID: "role1"},
	}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)

	fieldPerms := []*entity.FieldPermission{
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", resourceType).Return(fieldPerms, nil)

	perm, err := checker.CheckFieldPermission(ctx, userID, resourceType, field)

	assert.NoError(t, err)
	assert.Equal(t, VISIBLE, perm, "readonly字段应该只可见")

	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestFieldPermissionChecker_FilterVisibleFields 测试过滤可见字段
func TestFieldPermissionChecker_FilterVisibleFields(t *testing.T) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	fields := []string{"bot_id", "name", "api_key", "webhook_url"}

	// Mock返回用户的角色和字段权限
	roles := []*entity.Role{
		{RoleID: "role1"},
	}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)

	fieldPerms := []*entity.FieldPermission{
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelHidden,
		},
		{
			FieldName:       "webhook_url",
			PermissionLevel: entity.FieldPermissionLevelHidden,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", resourceType).Return(fieldPerms, nil)

	visibleFields, err := checker.FilterVisibleFields(ctx, userID, resourceType, fields)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(visibleFields))
	assert.Contains(t, visibleFields, "bot_id")
	assert.Contains(t, visibleFields, "name")
	assert.NotContains(t, visibleFields, "api_key")
	assert.NotContains(t, visibleFields, "webhook_url")

	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestFieldPermissionChecker_FilterEditableFields 测试过滤可编辑字段
func TestFieldPermissionChecker_FilterEditableFields(t *testing.T) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	fields := []string{"bot_id", "name", "api_key", "description"}

	// Mock返回用户的角色和字段权限
	roles := []*entity.Role{
		{RoleID: "role1"},
	}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)

	fieldPerms := []*entity.FieldPermission{
		{
			FieldName:       "bot_id",
			PermissionLevel: entity.FieldPermissionLevelReadonly,
		},
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", resourceType).Return(fieldPerms, nil)

	editableFields, err := checker.FilterEditableFields(ctx, userID, resourceType, fields)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(editableFields))
	assert.Contains(t, editableFields, "name")
	assert.Contains(t, editableFields, "description")
	assert.NotContains(t, editableFields, "bot_id")
	assert.NotContains(t, editableFields, "api_key")

	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestFieldPermissionChecker_MaskSensitiveFields 测试脱敏敏感字段
func TestFieldPermissionChecker_MaskSensitiveFields(t *testing.T) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"

	data := map[string]interface{}{
		"bot_id":     "bot123",
		"name":       "Test Bot",
		"api_key":    "sk-1234567890abcdef",
		"webhook_url": "https://example.com/webhook",
	}

	// Mock返回用户的角色和字段权限
	roles := []*entity.Role{
		{RoleID: "role1"},
	}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)

	fieldPerms := []*entity.FieldPermission{
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly,
			// 需要扩展FieldPermission实体以支持MaskRule字段
		},
		{
			FieldName:       "webhook_url",
			PermissionLevel: entity.FieldPermissionLevelHidden,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", "bots").Return(fieldPerms, nil)

	maskedData, err := checker.MaskSensitiveFields(ctx, userID, data)

	assert.NoError(t, err)
	assert.NotNil(t, maskedData)
	assert.Equal(t, "bot123", maskedData["bot_id"])
	assert.Equal(t, "Test Bot", maskedData["name"])
	// api_key可能被脱敏，取决于实现
	assert.NotNil(t, maskedData["api_key"])
	// webhook_url应该被删除
	_, exists := maskedData["webhook_url"]
	assert.False(t, exists, "hidden字段应该被删除")

	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestFieldPermissionChecker_ValidateFieldPermissions 测试验证字段权限
func TestFieldPermissionChecker_ValidateFieldPermissions(t *testing.T) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	data := map[string]interface{}{
		"name":        "New Bot",
		"api_key":     "sk-new-key",
		"description": "Test description",
	}

	// Mock返回用户的角色和字段权限
	roles := []*entity.Role{
		{RoleID: "role1"},
	}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)

	fieldPerms := []*entity.FieldPermission{
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", resourceType).Return(fieldPerms, nil)

	err := checker.ValidateFieldPermissions(ctx, userID, resourceType, data)

	assert.Error(t, err)
	assert.IsType(t, &FieldPermissionDeniedError{}, err)
	permErr := err.(*FieldPermissionDeniedError)
	assert.Equal(t, "api_key", permErr.FieldName)
	assert.Equal(t, "field is not editable", permErr.Reason)

	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestFieldPermissionChecker_GetFieldPermissions 测试获取所有字段权限
func TestFieldPermissionChecker_GetFieldPermissions(t *testing.T) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	// Mock返回用户的角色和字段权限
	roles := []*entity.Role{
		{RoleID: "role1"},
		{RoleID: "role2"},
	}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)

	// Role1的字段权限
	role1Perms := []*entity.FieldPermission{
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelReadonly,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", resourceType).Return(role1Perms, nil)

	// Role2的字段权限
	role2Perms := []*entity.FieldPermission{
		{
			FieldName:       "api_key",
			PermissionLevel: entity.FieldPermissionLevelEditable,
		},
		{
			FieldName:       "webhook_url",
			PermissionLevel: entity.FieldPermissionLevelHidden,
		},
	}
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role2", resourceType).Return(role2Perms, nil)

	fieldPerms, err := checker.GetFieldPermissions(ctx, userID, resourceType)

	assert.NoError(t, err)
	assert.NotNil(t, fieldPerms)

	// api_key应该取最大权限（editable）
	apiKeyPerm := fieldPerms["api_key"]
	assert.NotNil(t, apiKeyPerm)
	assert.True(t, apiKeyPerm.Permission.HasEditable(), "api_key应该是可编辑的")

	// webhook_url应该是hidden
	webhookPerm := fieldPerms["webhook_url"]
	assert.NotNil(t, webhookPerm)
	assert.False(t, webhookPerm.Permission.HasVisible(), "webhook_url应该是隐藏的")

	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestFieldPermission_String 测试字段权限字符串表示
func TestFieldPermission_String(t *testing.T) {
	tests := []struct {
		name     string
		perm     FieldPermission
		expected string
	}{
		{
			name:     "None",
			perm:     0,
			expected: "NONE",
		},
		{
			name:     "Visible only",
			perm:     VISIBLE,
			expected: "VISIBLE",
		},
		{
			name:     "Editable",
			perm:     VISIBLE | EDITABLE,
			expected: "VISIBLE|EDITABLE",
		},
		{
			name:     "Required",
			perm:     VISIBLE | EDITABLE | REQUIRED,
			expected: "VISIBLE|EDITABLE|REQUIRED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.perm.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFieldPermission_HasXXX 测试字段权限位操作
func TestFieldPermission_HasXXX(t *testing.T) {
	perm := VISIBLE | EDITABLE | REQUIRED

	assert.True(t, perm.HasVisible())
	assert.True(t, perm.HasEditable())
	assert.True(t, perm.HasRequired())

	perm2 := VISIBLE
	assert.True(t, perm2.HasVisible())
	assert.False(t, perm2.HasEditable())
	assert.False(t, perm2.HasRequired())
}

// BenchmarkFieldPermissionChecker_FilterVisibleFields 性能测试
func BenchmarkFieldPermissionChecker_FilterVisibleFields(b *testing.B) {
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)

	checker := NewFieldPermissionChecker(mockUserRoleRepo, mockFieldPermRepo)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"
	fields := []string{"bot_id", "name", "description", "api_key", "webhook_url"}

	// Mock setup
	roles := []*entity.Role{{RoleID: "role1"}}
	mockUserRoleRepo.On("GetRolesByUser", ctx, userID, "").Return(roles, nil)
	mockFieldPermRepo.On("GetByRoleAndResource", ctx, "role1", resourceType).Return([]*entity.FieldPermission{}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = checker.FilterVisibleFields(ctx, userID, resourceType, fields)
	}
}
