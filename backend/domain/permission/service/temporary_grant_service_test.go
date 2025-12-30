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
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// MockTemporaryGrantRepository 临时授权仓储 Mock
type MockTemporaryGrantRepository struct {
	mock.Mock
}

func (m *MockTemporaryGrantRepository) Create(ctx context.Context, grant *entity.TemporaryGrant) error {
	args := m.Called(ctx, grant)
	return args.Error(0)
}

func (m *MockTemporaryGrantRepository) GetByID(ctx context.Context, grantID int64) (*entity.TemporaryGrant, error) {
	args := m.Called(ctx, grantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TemporaryGrant), args.Error(1)
}

func (m *MockTemporaryGrantRepository) GetByGrantCode(ctx context.Context, grantCode string) (*entity.TemporaryGrant, error) {
	args := m.Called(ctx, grantCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TemporaryGrant), args.Error(1)
}

func (m *MockTemporaryGrantRepository) Update(ctx context.Context, grant *entity.TemporaryGrant) error {
	args := m.Called(ctx, grant)
	return args.Error(0)
}

func (m *MockTemporaryGrantRepository) Delete(ctx context.Context, grantID int64) error {
	args := m.Called(ctx, grantID)
	return args.Error(0)
}

func (m *MockTemporaryGrantRepository) List(ctx context.Context, filter *repository.TemporaryGrantFilter) ([]*entity.TemporaryGrant, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.TemporaryGrant), args.Get(1).(int64), args.Error(2)
}

func (m *MockTemporaryGrantRepository) ListExpired(ctx context.Context, expiresBefore int64) ([]*entity.TemporaryGrant, error) {
	args := m.Called(ctx, expiresBefore)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.TemporaryGrant), args.Error(1)
}

func (m *MockTemporaryGrantRepository) DeleteExpired(ctx context.Context, expiresBefore int64) (int64, error) {
	args := m.Called(ctx, expiresBefore)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTemporaryGrantRepository) GetByGrantee(ctx context.Context, granteeID string) ([]*entity.TemporaryGrant, error) {
	args := m.Called(ctx, granteeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.TemporaryGrant), args.Error(1)
}

func (m *MockTemporaryGrantRepository) GetByGrantor(ctx context.Context, grantorID string) ([]*entity.TemporaryGrant, error) {
	args := m.Called(ctx, grantorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.TemporaryGrant), args.Error(1)
}

// MockTemporaryGrantHistoryRepository 临时授权历史仓储 Mock
type MockTemporaryGrantHistoryRepository struct {
	mock.Mock
}

func (m *MockTemporaryGrantHistoryRepository) Create(ctx context.Context, history *entity.TemporaryGrantHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockTemporaryGrantHistoryRepository) GetByGrantID(ctx context.Context, grantID int64) ([]*entity.TemporaryGrantHistory, error) {
	args := m.Called(ctx, grantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.TemporaryGrantHistory), args.Error(1)
}

func (m *MockTemporaryGrantHistoryRepository) GetByGrantCode(ctx context.Context, grantCode string) ([]*entity.TemporaryGrantHistory, error) {
	args := m.Called(ctx, grantCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.TemporaryGrantHistory), args.Error(1)
}

func (m *MockTemporaryGrantHistoryRepository) List(ctx context.Context, filter *repository.HistoryFilter) ([]*entity.TemporaryGrantHistory, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.TemporaryGrantHistory), args.Get(1).(int64), args.Error(2)
}

func (m *MockTemporaryGrantHistoryRepository) DeleteByGrantID(ctx context.Context, grantID int64) error {
	args := m.Called(ctx, grantID)
	return args.Error(0)
}

func (m *MockTemporaryGrantHistoryRepository) DeleteExpired(ctx context.Context, before int64) (int64, error) {
	args := m.Called(ctx, before)
	return args.Get(0).(int64), args.Error(1)
}

// TestTemporaryGrantService_CreateTemporaryGrant 测试创建临时授权
func TestTemporaryGrantService_CreateTemporaryGrant(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	permissionData := entity.PermissionData{"role_id": "role-123"}
	req := &CreateTemporaryGrantRequest{
		TenantID:       "tenant-123",
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		Duration:       24 * time.Hour,
		Reason:         "临时项目授权",
	}

	mockGrantRepo.On("Create", ctx, mock.AnythingOfType("*entity.TemporaryGrant")).Return(nil)
	mockHistoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.TemporaryGrantHistory")).Return(nil)

	// Act
	grant, err := service.CreateTemporaryGrant(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, grant)
	assert.NotEmpty(t, grant.GrantCode)
	assert.Equal(t, req.GranteeID, grant.GranteeID)
	assert.Equal(t, req.GrantorID, grant.GrantorID)
	assert.Equal(t, req.TenantID, grant.TenantID)
	assert.Equal(t, req.PermissionType, grant.PermissionType)
	assert.Greater(t, grant.ExpiresAt, int64(0))

	// 验证授权码格式（XXXX-XXXX-XXXX-XXXX）
	assert.Len(t, grant.GrantCode, 19) // 4组4位 + 3个连字符

	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_UseTemporaryGrant_Success 测试使用临时授权成功
func TestTemporaryGrantService_UseTemporaryGrant_Success(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      expiresAt,
		IsUsed:         false,
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)
	mockUserRoleRepo.On("Exists", ctx, "user-456", "tenant-123", "role-123").Return(false, nil)
	mockUserRoleRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserRole")).Return(nil)
	mockGrantRepo.On("Update", ctx, mock.MatchedBy(func(g *entity.TemporaryGrant) bool {
		return g.IsUsed == true
	})).Return(nil)
	mockHistoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.TemporaryGrantHistory")).Return(nil)

	// Act
	userRole, err := service.UseTemporaryGrant(ctx, grantCode, "user-456")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, userRole)
	assert.Equal(t, "user-456", userRole.UserID)
	assert.Equal(t, "tenant-123", userRole.TenantID)
	assert.Equal(t, "role-123", userRole.RoleID)
	assert.NotNil(t, userRole.ExpiresAt)

	mockGrantRepo.AssertExpectations(t)
	mockUserRoleRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_UseTemporaryGrant_Expired 测试使用已过期的授权
func TestTemporaryGrantService_UseTemporaryGrant_Expired(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expiresAt := time.Now().Add(-24 * time.Hour).UnixMilli() // 已过期

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      expiresAt,
		IsUsed:         false,
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)

	// Act
	userRole, err := service.UseTemporaryGrant(ctx, grantCode, "user-456")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, userRole)
	assert.Contains(t, err.Error(), "授权码已过期")

	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_UseTemporaryGrant_AlreadyUsed 测试使用已使用的授权
func TestTemporaryGrantService_UseTemporaryGrant_AlreadyUsed(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      expiresAt,
		IsUsed:         true, // 已使用
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)

	// Act
	userRole, err := service.UseTemporaryGrant(ctx, grantCode, "user-456")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, userRole)
	assert.Contains(t, err.Error(), "授权码已被使用")

	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_UseTemporaryGrant_WrongUser 测试错误的用户使用授权
func TestTemporaryGrantService_UseTemporaryGrant_WrongUser(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",      // 授权给 user-456
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      expiresAt,
		IsUsed:         false,
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)

	// Act - user-999 尝试使用授权给 user-456 的授权码
	userRole, err := service.UseTemporaryGrant(ctx, grantCode, "user-999")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, userRole)
	assert.Contains(t, err.Error(), "授权码不属于当前用户")

	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_UseTemporaryGrant_AlreadyHasRole 测试用户已拥有角色
func TestTemporaryGrantService_UseTemporaryGrant_AlreadyHasRole(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      expiresAt,
		IsUsed:         false,
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)
	mockUserRoleRepo.On("Exists", ctx, "user-456", "tenant-123", "role-123").Return(true, nil) // 已拥有角色

	// Act
	userRole, err := service.UseTemporaryGrant(ctx, grantCode, "user-456")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, userRole)
	assert.Contains(t, err.Error(), "用户已拥有该角色")

	mockGrantRepo.AssertExpectations(t)
	mockUserRoleRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_RevokeTemporaryGrant_Success 测试撤销临时授权成功
func TestTemporaryGrantService_RevokeTemporaryGrant(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      time.Now().Add(24 * time.Hour).UnixMilli(),
		IsUsed:         false,
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)
	mockGrantRepo.On("Update", ctx, mock.MatchedBy(func(g *entity.TemporaryGrant) bool {
		return g.IsRevoked == true
	})).Return(nil)
	mockHistoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.TemporaryGrantHistory")).Return(nil)

	// Act
	err := service.RevokeTemporaryGrant(ctx, grantCode, "user-789", "测试撤销")

	// Assert
	assert.NoError(t, err)
	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_RevokeTemporaryGrant_AlreadyUsed 测试撤销已使用的授权
func TestTemporaryGrantService_RevokeTemporaryGrant_AlreadyUsed(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"role_id": "role-123"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionTypeRole,
		PermissionData: permissionData,
		ExpiresAt:      time.Now().Add(24 * time.Hour).UnixMilli(),
		IsUsed:         true, // 已使用
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)

	// Act
	err := service.RevokeTemporaryGrant(ctx, grantCode, "user-789", "测试撤销")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "授权码已被使用")
	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_CleanupExpiredGrants 测试清理过期授权
func TestTemporaryGrantService_CleanupExpiredGrants(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	expiredGrants := []*entity.TemporaryGrant{
		{GrantID: 1, GrantCode: "EXPIRED-1"},
		{GrantID: 2, GrantCode: "EXPIRED-2"},
	}

	mockGrantRepo.On("ListExpired", ctx, mock.AnythingOfType("int64")).Return(expiredGrants, nil)
	mockHistoryRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.TemporaryGrantHistory")).Return(nil)
	mockGrantRepo.On("DeleteExpired", ctx, mock.AnythingOfType("int64")).Return(int64(2), nil)

	// Act
	count, err := service.CleanupExpiredGrants(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count) // 2个过期授权 + 2个被删除
	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_GetTemporaryGrant 测试获取临时授权
func TestTemporaryGrantService_GetTemporaryGrant(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expectedGrant := &entity.TemporaryGrant{
		GrantID:   1,
		GrantCode: grantCode,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(expectedGrant, nil)

	// Act
	grant, err := service.GetTemporaryGrant(ctx, grantCode)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedGrant, grant)
	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_GetTemporaryGrant_NotFound 测试获取不存在的授权
func TestTemporaryGrantService_GetTemporaryGrant_NotFound(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "INVALID-CODE"
	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return((*entity.TemporaryGrant)(nil), nil)

	// Act
	grant, err := service.GetTemporaryGrant(ctx, grantCode)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, grant)
	assert.Contains(t, err.Error(), "授权码不存在")
	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_GenerateGrantCode 测试生成授权码
func TestTemporaryGrantService_GenerateGrantCode(t *testing.T) {
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)

	// 生成多个授权码，验证格式和唯一性
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := service.generateGrantCode()
		require.NoError(t, err)
		require.NotEmpty(t, code)

		// 验证格式：XXXX-XXXX-XXXX-XXXX
		assert.Len(t, code, 19)
		assert.Equal(t, '-', code[4])
		assert.Equal(t, '-', code[9])
		assert.Equal(t, '-', code[14])

		// 验证唯一性
		_, exists := codes[code]
		assert.False(t, exists, "Generated duplicate code: %s", code)
		codes[code] = true
	}
}

// TestTemporaryGrantService_UseTemporaryGrant_UnsupportedPermissionType 测试不支持的权限类型
func TestTemporaryGrantService_UseTemporaryGrant_UnsupportedPermissionType(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	grantCode := "ABCD-EFGH-IJKL-MNOP"
	expiresAt := time.Now().Add(24 * time.Hour).UnixMilli()

	permissionDataJSON, _ := json.Marshal(entity.PermissionData{"field": "value"})
	var permissionData entity.PermissionData
	json.Unmarshal(permissionDataJSON, &permissionData)

	grant := &entity.TemporaryGrant{
		GrantID:        1,
		GrantCode:      grantCode,
		GranteeID:      "user-456",
		GrantorID:      "user-789",
		TenantID:       "tenant-123",
		PermissionType: entity.PermissionType("unsupported_type"), // 不支持的类型
		PermissionData: permissionData,
		ExpiresAt:      expiresAt,
		IsUsed:         false,
		IsRevoked:      false,
	}

	mockGrantRepo.On("GetByGrantCode", ctx, grantCode).Return(grant, nil)

	// Act
	userRole, err := service.UseTemporaryGrant(ctx, grantCode, "user-456")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, userRole)
	assert.Contains(t, err.Error(), "无效的权限类型")

	mockGrantRepo.AssertExpectations(t)
}

// TestTemporaryGrantService_ListTemporaryGrants 测试查询临时授权列表
func TestTemporaryGrantService_ListTemporaryGrants(t *testing.T) {
	// Arrange
	mockGrantRepo := new(MockTemporaryGrantRepository)
	mockHistoryRepo := new(MockTemporaryGrantHistoryRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	service := NewTemporaryGrantService(mockGrantRepo, mockHistoryRepo, mockUserRoleRepo)
	ctx := context.Background()

	expectedGrants := []*entity.TemporaryGrant{
		{GrantID: 1, GrantCode: "CODE-1"},
		{GrantID: 2, GrantCode: "CODE-2"},
	}

	filter := &repository.TemporaryGrantFilter{
		TenantID:  "tenant-123",
		GranteeID: "user-456",
	}

	mockGrantRepo.On("List", ctx, filter).Return(expectedGrants, int64(2), nil)

	// Act
	grants, total, err := service.ListTemporaryGrants(ctx, filter)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, grants, 2)
	assert.Equal(t, int64(2), total)
	mockGrantRepo.AssertExpectations(t)
}
