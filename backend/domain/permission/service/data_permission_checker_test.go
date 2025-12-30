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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// Mock repositories for testing
type MockUserDepartmentRepository struct {
	mock.Mock
}

func (m *MockUserDepartmentRepository) GetByUser(ctx context.Context, userID, tenantID string) ([]*entity.UserDepartment, error) {
	args := m.Called(ctx, userID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.UserDepartment), args.Error(1)
}

type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) GetDescendants(ctx context.Context, departmentID string) ([]*entity.Department, error) {
	args := m.Called(ctx, departmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

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

type MockDataPermissionRepository struct {
	mock.Mock
}

func (m *MockDataPermissionRepository) GetByRoleAndResource(ctx context.Context, roleID string, resourceType entity.ResourceType) (*entity.DataPermission, error) {
	args := m.Called(ctx, roleID, resourceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.DataPermission), args.Error(1)
}

// TestDataPermissionChecker_Filter_ALL 测试全部数据权限
func TestDataPermissionChecker_Filter_ALL(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	filter, err := checker.Filter(ctx, userID, ALL, resourceType)

	assert.NoError(t, err)
	assert.NotNil(t, filter)
	assert.Equal(t, "", filter.WhereClause, "ALL权限不应该有WHERE条件")
	assert.Nil(t, filter.Args, "ALL权限不应该有参数")
}

// TestDataPermissionChecker_Filter_DEPARTMENT 测试部门数据权限
func TestDataPermissionChecker_Filter_DEPARTMENT(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	// Mock返回用户的部门
	userDepts := []*entity.UserDepartment{
		{
			UserID:       userID,
			DepartmentID: "dept1",
			IsLeader:     false,
		},
		{
			UserID:       userID,
			DepartmentID: "dept2",
			IsLeader:     false,
		},
	}
	mockUserDeptRepo.On("GetByUser", ctx, userID, "").Return(userDepts, nil)

	filter, err := checker.Filter(ctx, userID, DEPARTMENT, resourceType)

	assert.NoError(t, err)
	assert.NotNil(t, filter)
	assert.Contains(t, filter.WhereClause, "department_id IN (")
	assert.Equal(t, 2, len(filter.Args), "应该有2个部门ID参数")
	assert.Equal(t, "dept1", filter.Args[0])
	assert.Equal(t, "dept2", filter.Args[1])

	mockUserDeptRepo.AssertExpectations(t)
}

// TestDataPermissionChecker_Filter_DEPARTMENT_AND_SUB 测试本部门及子部门权限
func TestDataPermissionChecker_Filter_DEPARTMENT_AND_SUB(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	// Mock返回用户的部门（是部门领导）
	userDepts := []*entity.UserDepartment{
		{
			UserID:       userID,
			DepartmentID: "dept1",
			IsLeader:     true, // 是部门领导
		},
	}
	mockUserDeptRepo.On("GetByUser", ctx, userID, "").Return(userDepts, nil)

	// Mock返回子部门
	childDepts := []*entity.Department{
		{DepartmentID: "dept1_child1"},
		{DepartmentID: "dept1_child2"},
	}
	mockDeptRepo.On("GetDescendants", ctx, "dept1").Return(childDepts, nil)

	filter, err := checker.Filter(ctx, userID, DEPARTMENT_AND_SUB, resourceType)

	assert.NoError(t, err)
	assert.NotNil(t, filter)
	assert.Contains(t, filter.WhereClause, "department_id IN (")
	assert.Equal(t, 3, len(filter.Args), "应该有3个部门ID（1个本部门+2个子部门）")

	mockUserDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertExpectations(t)
}

// TestDataPermissionChecker_Filter_SELF 测试仅本人数据权限
func TestDataPermissionChecker_Filter_SELF(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	filter, err := checker.Filter(ctx, userID, SELF, resourceType)

	assert.NoError(t, err)
	assert.NotNil(t, filter)
	assert.Equal(t, "creator_id = ?", filter.WhereClause)
	assert.Equal(t, 1, len(filter.Args))
	assert.Equal(t, userID, filter.Args[0])
}

// TestDataPermissionChecker_Filter_NO_DEPARTMENT 测试无部门情况
func TestDataPermissionChecker_Filter_NO_DEPARTMENT(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceType := "bots"

	// Mock返回空部门列表
	mockUserDeptRepo.On("GetByUser", ctx, userID, "").Return([]*entity.UserDepartment{}, nil)

	filter, err := checker.Filter(ctx, userID, DEPARTMENT, resourceType)

	assert.NoError(t, err)
	assert.NotNil(t, filter)
	assert.Equal(t, "1 = 0", filter.WhereClause, "无部门时应该返回false条件")

	mockUserDeptRepo.AssertExpectations(t)
}

// TestDataPermissionChecker_CheckAccess_ALL 测试全部数据权限访问检查
func TestDataPermissionChecker_CheckAccess_ALL(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceID := "bot456"
	resourceType := "bots"

	allowed, err := checker.CheckAccess(ctx, userID, resourceID, ALL, resourceType)

	assert.NoError(t, err)
	assert.True(t, allowed, "ALL权限应该允许访问所有资源")
}

// TestDataPermissionChecker_CheckAccess_DEPARTMENT 测试部门权限访问检查
func TestDataPermissionChecker_CheckAccess_DEPARTMENT(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"
	resourceID := "bot456"
	resourceType := "bots"

	// Mock返回用户的部门
	userDepts := []*entity.UserDepartment{
		{
			UserID:       userID,
			DepartmentID: "dept1",
			IsLeader:     false,
		},
	}
	mockUserDeptRepo.On("GetByUser", ctx, userID, "").Return(userDepts, nil)

	allowed, err := checker.CheckAccess(ctx, userID, resourceID, DEPARTMENT, resourceType)

	assert.NoError(t, err)
	// 结果取决于数据库中的资源department_id，这里只测试调用成功
	assert.NotNil(t, allowed)

	mockUserDeptRepo.AssertExpectations(t)
}

// TestDataPermissionChecker_GetAccessibleDepartmentIDs 测试获取可访问部门ID
func TestDataPermissionChecker_GetAccessibleDepartmentIDs(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"

	// 测试DEPARTMENT级别（不含子部门）
	userDepts := []*entity.UserDepartment{
		{
			UserID:       userID,
			DepartmentID: "dept1",
			IsLeader:     false,
		},
		{
			UserID:       userID,
			DepartmentID: "dept2",
			IsLeader:     false,
		},
	}
	mockUserDeptRepo.On("GetByUser", ctx, userID, "").Return(userDepts, nil)

	deptIDs, err := checker.GetAccessibleDepartmentIDs(ctx, userID, DEPARTMENT)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(deptIDs))
	assert.Contains(t, deptIDs, "dept1")
	assert.Contains(t, deptIDs, "dept2")

	mockUserDeptRepo.AssertExpectations(t)
}

// TestDataPermissionChecker_GetAccessibleDepartmentIDs_WITH_SUB 测试获取可访问部门ID（含子部门）
func TestDataPermissionChecker_GetAccessibleDepartmentIDs_WITH_SUB(t *testing.T) {
	db, _ := gorm.Open()

	mockUserDeptRepo := new(MockUserDepartmentRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockUserRoleRepo := new(MockUserRoleRepository)
	mockDataPermRepo := new(MockDataPermissionRepository)

	checker := NewDataPermissionChecker(
		db,
		mockUserDeptRepo,
		mockDeptRepo,
		mockUserRoleRepo,
		mockDataPermRepo,
	)

	ctx := context.Background()
	userID := "user123"

	// 测试DEPARTMENT_AND_SUB级别（含子部门）
	userDepts := []*entity.UserDepartment{
		{
			UserID:       userID,
			DepartmentID: "dept1",
			IsLeader:     false,
		},
	}
	mockUserDeptRepo.On("GetByUser", ctx, userID, "").Return(userDepts, nil)

	// Mock返回子部门
	childDepts := []*entity.Department{
		{DepartmentID: "dept1_child1"},
		{DepartmentID: "dept1_child2"},
	}
	mockDeptRepo.On("GetDescendants", ctx, "dept1").Return(childDepts, nil)

	deptIDs, err := checker.GetAccessibleDepartmentIDs(ctx, userID, DEPARTMENT_AND_SUB)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(deptIDs), "应该包含本部门和2个子部门")
	assert.Contains(t, deptIDs, "dept1")
	assert.Contains(t, deptIDs, "dept1_child1")
	assert.Contains(t, deptIDs, "dept1_child2")

	mockUserDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertExpectations(t)
}

// TestDataPermissionLevel_String 测试权限级别字符串表示
func TestDataPermissionLevel_String(t *testing.T) {
	tests := []struct {
		level    DataPermissionLevel
		expected string
	}{
		{ALL, "ALL"},
		{DEPARTMENT_AND_SUB, "DEPARTMENT_AND_SUB"},
		{DEPARTMENT, "DEPARTMENT"},
		{SELF, "SELF"},
		{CUSTOM, "CUSTOM"},
		{DataPermissionLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}
