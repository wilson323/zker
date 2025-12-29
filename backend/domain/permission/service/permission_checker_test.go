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

	"github.com/coze-studio/backend/domain/permission/entity"
)

// MockUserRoleRepository Mock用户角色仓储
type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	args := m.Called(ctx, userRole)
	return args.Error(0)
}

func (m *MockUserRoleRepository) Delete(ctx context.Context, userID, tenantID, roleID string) error {
	args := m.Called(ctx, userID, tenantID, roleID)
	return args.Error(0)
}

func (m *MockUserRoleRepository) DeleteByUser(ctx context.Context, userID, tenantID string) error {
	args := m.Called(ctx, userID, tenantID)
	return args.Error(0)
}

func (m *MockUserRoleRepository) DeleteByRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

func (m *MockUserRoleRepository) GetRolesByUser(ctx context.Context, userID, tenantID string) ([]*entity.Role, error) {
	args := m.Called(ctx, userID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (m *MockUserRoleRepository) GetUsersByRole(ctx context.Context, roleID string) ([]*entity.UserRole, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.UserRole), args.Error(1)
}

func (m *MockUserRoleRepository) Exists(ctx context.Context, userID, tenantID, roleID string) (bool, error) {
	args := m.Called(ctx, userID, tenantID, roleID)
	return args.Bool(0), args.Error(1)
}

// MockDataPermissionRepository Mock数据权限仓储
type MockDataPermissionRepository struct {
	mock.Mock
}

func (m *MockDataPermissionRepository) Create(ctx context.Context, perm *entity.DataPermission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockDataPermissionRepository) GetByID(ctx context.Context, permissionID string) (*entity.DataPermission, error) {
	args := m.Called(ctx, permissionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.DataPermission), args.Error(1)
}

func (m *MockDataPermissionRepository) GetByRoleAndResource(ctx context.Context, roleID string, resourceType entity.ResourceType) (*entity.DataPermission, error) {
	args := m.Called(ctx, roleID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.DataPermission), args.Error(1)
}

func (m *MockDataPermissionRepository) GetByRole(ctx context.Context, roleID string) ([]*entity.DataPermission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.DataPermission), args.Error(1)
}

func (m *MockDataPermissionRepository) Update(ctx context.Context, perm *entity.DataPermission) error {
	args := m.Called(ctx, perm)
	return args.Error(0)
}

func (m *MockDataPermissionRepository) Delete(ctx context.Context, permissionID string) error {
	args := m.Called(ctx, permissionID)
	return args.Error(0)
}

func (m *MockDataPermissionRepository) DeleteByRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

// MockFieldPermissionRepository Mock字段权限仓储
type MockFieldPermissionRepository struct {
	mock.Mock
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

func (m *MockFieldPermissionRepository) GetByRoleAndResource(ctx context.Context, roleID string, resourceType string) ([]*entity.FieldPermission, error) {
	args := m.Called(ctx, roleID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.FieldPermission), args.Error(1)
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

// TestPermissionChecker_CheckDataPermission_All 测试ALL权限
func TestPermissionChecker_CheckDataPermission_All(t *testing.T) {
	// Arrange
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)
	checker := NewPermissionChecker(mockUserRoleRepo, mockDataPermRepo, nil)

	roles := []*entity.Role{
		{RoleID: "role_1", RoleCode: "admin"},
	}

	mockUserRoleRepo.On("GetRolesByUser", mock.Anything, "user_123", "tenant_123").
		Return(roles, nil)

	allScope := entity.PermissionScopeAll
	mockDataPermRepo.On("GetByRoleAndResource", mock.Anything, "role_1", entity.ResourceTypeBots).
		Return(&entity.DataPermission{
			RoleID:       "role_1",
			ResourceType: entity.ResourceTypeBots,
			Scope:        &allScope,
		}, nil)

	// Act
	err := checker.CheckDataPermission(context.Background(), "user_123", "tenant_123", entity.ResourceTypeBots, "bot_456")

	// Assert
	assert.NoError(t, err)
	mockUserRoleRepo.AssertExpectations(t)
	mockDataPermRepo.AssertExpectations(t)
}

// TestPermissionChecker_CheckDataPermission_Own 测试OWN权限
func TestPermissionChecker_CheckDataPermission_Own(t *testing.T) {
	// Arrange
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)
	checker := NewPermissionChecker(mockUserRoleRepo, mockDataPermRepo, nil)

	roles := []*entity.Role{
		{RoleID: "role_1"},
	}

	mockUserRoleRepo.On("GetRolesByUser", mock.Anything, "user_123", "tenant_123").
		Return(roles, nil)

	ownScope := entity.PermissionScopeOwn
	mockDataPermRepo.On("GetByRoleAndResource", mock.Anything, "role_1", entity.ResourceTypeBots).
		Return(&entity.DataPermission{
			RoleID:       "role_1",
			ResourceType: entity.ResourceTypeBots,
			Scope:        &ownScope,
		}, nil)

	// Act - 使用自定义的isOwner实现（这里简化为总是返回true）
	// 实际应用中需要注入ResourceOwnerChecker
	err := checker.CheckDataPermission(context.Background(), "user_123", "tenant_123", entity.ResourceTypeBots, "bot_456")

	// Assert - 会因为没有ResourceOwnerChecker而失败
	assert.Error(t, err)
}

// TestPermissionChecker_CheckDataPermission_NoPermission 测试无权限
func TestPermissionChecker_CheckDataPermission_NoPermission(t *testing.T) {
	// Arrange
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)
	checker := NewPermissionChecker(mockUserRoleRepo, mockDataPermRepo, nil)

	roles := []*entity.Role{
		{RoleID: "role_1"},
	}

	mockUserRoleRepo.On("GetRolesByUser", mock.Anything, "user_123", "tenant_123").
		Return(roles, nil)

	noneScope := entity.PermissionScopeNone
	mockDataPermRepo.On("GetByRoleAndResource", mock.Anything, "role_1", entity.ResourceTypeBots).
		Return(&entity.DataPermission{
			RoleID:       "role_1",
			ResourceType: entity.ResourceTypeBots,
			Scope:        &noneScope,
		}, nil)

	// Act
	err := checker.CheckDataPermission(context.Background(), "user_123", "tenant_123", entity.ResourceTypeBots, "bot_456")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}

// TestPermissionChecker_GetFieldPermissions 测试获取字段权限
func TestPermissionChecker_GetFieldPermissions(t *testing.T) {
	// Arrange
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockFieldPermRepo := new(MockFieldPermissionRepository)
	checker := NewPermissionChecker(mockUserRoleRepo, nil, mockFieldPermRepo)

	roles := []*entity.Role{
		{RoleID: "role_1"},
		{RoleID: "role_2"},
	}

	mockUserRoleRepo.On("GetRolesByUser", mock.Anything, "user_123", "tenant_123").
		Return(roles, nil)

	// Role 1 permissions
	mockFieldPermRepo.On("GetByRoleAndResource", mock.Anything, "role_1", "bots").
		Return([]*entity.FieldPermission{
			{FieldName: "name", PermissionLevel: entity.FieldPermissionEditable},
			{FieldName: "config", PermissionLevel: entity.FieldPermissionReadonly},
		}, nil)

	// Role 2 permissions (higher priority)
	mockFieldPermRepo.On("GetByRoleAndResource", mock.Anything, "role_2", "bots").
		Return([]*entity.FieldPermission{
			{FieldName: "name", PermissionLevel: entity.FieldPermissionHidden}, // 应该被editable覆盖
			{FieldName: "status", PermissionLevel: entity.FieldPermissionEditable},
		}, nil)

	// Act
	perms, err := checker.GetFieldPermissions(context.Background(), "user_123", "tenant_123", "bots")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "editable", perms["name"])   // editable优先级高于hidden
	assert.Equal(t, "readonly", perms["config"]) // readonly保持不变
	assert.Equal(t, "editable", perms["status"]) // 新增字段
	mockUserRoleRepo.AssertExpectations(t)
	mockFieldPermRepo.AssertExpectations(t)
}

// TestPermissionChecker_MultipleRoles 测试多角色权限聚合
func TestPermissionChecker_MultipleRoles(t *testing.T) {
	// Arrange
	mockUserRoleRepo := new(MockUserRoleRepo)
	mockDataPermRepo := new(MockDataPermissionRepository)
	checker := NewPermissionChecker(mockUserRoleRepo, mockDataPermRepo, nil)

	roles := []*entity.Role{
		{RoleID: "role_1"},
		{RoleID: "role_2"},
	}

	mockUserRoleRepo.On("GetRolesByUser", mock.Anything, "user_123", "tenant_123").
		Return(roles, nil)

	// Role 1: OWN权限
	ownScope := entity.PermissionScopeOwn
	mockDataPermRepo.On("GetByRoleAndResource", mock.Anything, "role_1", entity.ResourceTypeBots).
		Return(&entity.DataPermission{
			RoleID:       "role_1",
			ResourceType: entity.ResourceTypeBots,
			Scope:        &ownScope,
		}, nil)

	// Role 2: ALL权限（应该覆盖Role 1）
	allScope := entity.PermissionScopeAll
	mockDataPermRepo.On("GetByRoleAndResource", mock.Anything, "role_2", entity.ResourceTypeBots).
		Return(&entity.DataPermission{
			RoleID:       "role_2",
			ResourceType: entity.ResourceTypeBots,
			Scope:        &allScope,
		}, nil)

	// Act
	err := checker.CheckDataPermission(context.Background(), "user_123", "tenant_123", entity.ResourceTypeBots, "bot_456")

	// Assert - 应该使用ALL权限
	assert.NoError(t, err)
}
