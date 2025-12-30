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
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
)

// DepartmentPermissionCheckerTestSuite 测试套件
type DepartmentPermissionCheckerTestSuite struct {
	suite.Suite
	checker               *DepartmentPermissionChecker
	mockUserDeptRepo      *MockUserDepartmentRepository
	mockDepartmentRepo    *MockDepartmentRepository
	ctx                   context.Context
}

func (s *DepartmentPermissionCheckerTestSuite) SetupTest() {
	s.mockUserDeptRepo = new(MockUserDepartmentRepository)
	s.mockDepartmentRepo = new(MockDepartmentRepository)
	s.checker = NewDepartmentPermissionChecker(s.mockUserDeptRepo, s.mockDepartmentRepo)
	s.ctx = context.Background()
}

// TestGetAccessibleDepartmentIDs_Success 测试获取可访问部门成功
func (s *DepartmentPermissionCheckerTestSuite) TestGetAccessibleDepartmentIDs_Success() {
	// Arrange
	userID := "user-123"
	tenantID := "tenant-123"

	userDepts := []*entity.UserDepartment{
		{UserID: userID, DepartmentID: "dept-1", IsLeader: false},
		{UserID: userID, DepartmentID: "dept-2", IsLeader: true},
	}

	s.mockUserDeptRepo.On("GetByUser", s.ctx, userID, tenantID).Return(userDepts, nil)

	// 部门2有子部门
	childDepts := []*entity.Department{
		{DepartmentID: "dept-2-1", ParentDepartmentID: func() *string { s := "dept-2"; return &s }()},
		{DepartmentID: "dept-2-2", ParentDepartmentID: func() *string { s := "dept-2"; return &s }()},
	}

	s.mockDepartmentRepo.On("GetDescendants", s.ctx, "dept-2").Return(childDepts, nil)

	// Act
	deptIDs, err := s.checker.GetAccessibleDepartmentIDs(s.ctx, userID, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), []string{"dept-1", "dept-2", "dept-2-1", "dept-2-2"}, deptIDs)
	s.mockUserDeptRepo.AssertExpectations(s.T())
	s.mockDepartmentRepo.AssertExpectations(s.T())
}

// TestGetAccessibleDepartmentIDs_NoLeader 测试非领导用户
func (s *DepartmentPermissionCheckerTestSuite) TestGetAccessibleDepartmentIDs_NoLeader() {
	// Arrange
	userID := "user-123"
	tenantID := "tenant-123"

	userDepts := []*entity.UserDepartment{
		{UserID: userID, DepartmentID: "dept-1", IsLeader: false},
		{UserID: userID, DepartmentID: "dept-2", IsLeader: false},
	}

	s.mockUserDeptRepo.On("GetByUser", s.ctx, userID, tenantID).Return(userDepts, nil)

	// Act
	deptIDs, err := s.checker.GetAccessibleDepartmentIDs(s.ctx, userID, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), []string{"dept-1", "dept-2"}, deptIDs)
	// 不应该调用 GetDescendants
	s.mockDepartmentRepo.AssertNotCalled(s.T(), "GetDescendants", mock.Anything, mock.Anything)
	s.mockUserDeptRepo.AssertExpectations(s.T())
}

// TestGetAccessibleDepartmentIDs_DescendantsError 测试获取子部门失败
func (s *DepartmentPermissionCheckerTestSuite) TestGetAccessibleDepartmentIDs_DescendantsError() {
	// Arrange
	userID := "user-123"
	tenantID := "tenant-123"

	userDepts := []*entity.UserDepartment{
		{UserID: userID, DepartmentID: "dept-1", IsLeader: false},
		{UserID: userID, DepartmentID: "dept-2", IsLeader: true},
	}

	s.mockUserDeptRepo.On("GetByUser", s.ctx, userID, tenantID).Return(userDepts, nil)
	s.mockDepartmentRepo.On("GetDescendants", s.ctx, "dept-2").Return(
		([]*entity.Department)(nil), errors.New("database error"))

	// Act
	deptIDs, err := s.checker.GetAccessibleDepartmentIDs(s.ctx, userID, tenantID)

	// Assert - 即使获取子部门失败，也应该返回基础部门
	assert.NoError(s.T(), err)
	assert.ElementsMatch(s.T(), []string{"dept-1", "dept-2"}, deptIDs)
	s.mockUserDeptRepo.AssertExpectations(s.T())
	s.mockDepartmentRepo.AssertExpectations(s.T())
}

// TestIsDepartmentLeader_Success 测试检查部门领导成功
func (s *DepartmentPermissionCheckerTestSuite) TestIsDepartmentLeader_Success() {
	// Arrange
	userID := "user-123"
	departmentID := "dept-123"

	userDept := &entity.UserDepartment{
		UserID:       userID,
		DepartmentID: departmentID,
		IsLeader:     true,
	}

	s.mockUserDeptRepo.On("GetByUserAndDepartment", s.ctx, userID, departmentID).Return(userDept, nil)

	// Act
	isLeader, err := s.checker.IsDepartmentLeader(s.ctx, userID, departmentID)

	// Assert
	assert.NoError(s.T(), err)
	assert.True(s.T(), isLeader)
	s.mockUserDeptRepo.AssertExpectations(s.T())
}

// TestIsDepartmentLeader_NotLeader 测试不是部门领导
func (s *DepartmentPermissionCheckerTestSuite) TestIsDepartmentLeader_NotLeader() {
	// Arrange
	userID := "user-123"
	departmentID := "dept-123"

	userDept := &entity.UserDepartment{
		UserID:       userID,
		DepartmentID: departmentID,
		IsLeader:     false,
	}

	s.mockUserDeptRepo.On("GetByUserAndDepartment", s.ctx, userID, departmentID).Return(userDept, nil)

	// Act
	isLeader, err := s.checker.IsDepartmentLeader(s.ctx, userID, departmentID)

	// Assert
	assert.NoError(s.T(), err)
	assert.False(s.T(), isLeader)
	s.mockUserDeptRepo.AssertExpectations(s.T())
}

// TestIsDepartmentLeader_NotMember 测试不是部门成员
func (s *DepartmentPermissionCheckerTestSuite) TestIsDepartmentLeader_NotMember() {
	// Arrange
	userID := "user-123"
	departmentID := "dept-123"

	s.mockUserDeptRepo.On("GetByUserAndDepartment", s.ctx, userID, departmentID).Return((*entity.UserDepartment)(nil), nil)

	// Act
	isLeader, err := s.checker.IsDepartmentLeader(s.ctx, userID, departmentID)

	// Assert
	assert.NoError(s.T(), err)
	assert.False(s.T(), isLeader)
	s.mockUserDeptRepo.AssertExpectations(s.T())
}

// TestFilterResourcesByDepartment_Success 测试按部门过滤资源成功
func (s *DepartmentPermissionCheckerTestSuite) TestFilterResourcesByDepartment_Success() {
	// Arrange
	userID := "user-123"
	tenantID := "tenant-123"

	accessibleDepts := []string{"dept-1", "dept-2"}
	s.mockUserDeptRepo.On("GetByUser", s.ctx, userID, tenantID).Return(
		[]*entity.UserDepartment{
			{UserID: userID, DepartmentID: "dept-1", IsLeader: false},
			{UserID: userID, DepartmentID: "dept-2", IsLeader: false},
		}, nil)

	resources := []ResourceWithDepartment{
		&mockResource{deptID: "dept-1"}, // 可访问
		&mockResource{deptID: "dept-2"}, // 可访问
		&mockResource{deptID: "dept-3"}, // 不可访问
		&mockResource{deptID: ""},       // 无部门信息，默认保留
	}

	// Act
	filtered, err := s.checker.FilterResourcesByDepartment(s.ctx, userID, tenantID, resources)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), filtered, 3) // dept-1, dept-2, 无部门资源
	s.mockUserDeptRepo.AssertExpectations(s.T())
}

// TestGetUserDepartments_Success 测试获取用户部门信息成功
func (s *DepartmentPermissionCheckerTestSuite) TestGetUserDepartments_Success() {
	// Arrange
	userID := "user-123"
	tenantID := "tenant-123"

	userDepts := []*entity.UserDepartment{
		{UserID: userID, DepartmentID: "dept-1", IsLeader: true},
		{UserID: userID, DepartmentID: "dept-2", IsLeader: false},
	}

	s.mockUserDeptRepo.On("GetByUser", s.ctx, userID, tenantID).Return(userDepts, nil)

	dept1 := &entity.Department{
		DepartmentID:       "dept-1",
		DepartmentName:     "Department 1",
		ParentDepartmentID: nil,
	}

	dept2 := &entity.Department{
		DepartmentID:       "dept-2",
		DepartmentName:     "Department 2",
		ParentDepartmentID: func() *string { s := "dept-1"; return &s }(),
	}

	s.mockDepartmentRepo.On("GetByID", s.ctx, "dept-1").Return(dept1, nil)
	s.mockDepartmentRepo.On("GetByID", s.ctx, "dept-2").Return(dept2, nil)

	// Act
	deptInfos, err := s.checker.GetUserDepartments(s.ctx, userID, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), deptInfos, 2)

	assert.Equal(s.T(), "dept-1", deptInfos[0].DepartmentID)
	assert.Equal(s.T(), "Department 1", deptInfos[0].DepartmentName)
	assert.True(s.T(), deptInfos[0].IsLeader)

	assert.Equal(s.T(), "dept-2", deptInfos[1].DepartmentID)
	assert.Equal(s.T(), "Department 2", deptInfos[1].DepartmentName)
	assert.False(s.T(), deptInfos[1].IsLeader)

	s.mockUserDeptRepo.AssertExpectations(s.T())
	s.mockDepartmentRepo.AssertExpectations(s.T())
}

// TestGetDepartmentTree_Success 测试获取部门树成功
func (s *DepartmentPermissionCheckerTestSuite) TestGetDepartmentTree_Success() {
	// Arrange
	tenantID := "tenant-123"

	parentID := "root"
	depts := []*entity.Department{
		{
			DepartmentID:       "dept-1",
			DepartmentName:     "Department 1",
			ParentDepartmentID: nil,
		},
		{
			DepartmentID:       "dept-2",
			DepartmentName:     "Department 2",
			ParentDepartmentID: &parentID,
		},
		{
			DepartmentID:       "dept-1-1",
			DepartmentName:     "Department 1-1",
			ParentDepartmentID: func() *string { s := "dept-1"; return &s }(),
		},
		{
			DepartmentID:       "dept-deleted",
			DepartmentName:     "Deleted Department",
			ParentDepartmentID: nil,
			DeletedAt:         func() *int64 { i := int64(123); return &i }(),
		},
	}

	s.mockDepartmentRepo.On("GetByTenant", s.ctx, tenantID).Return(depts, nil)

	// Act
	tree, err := s.checker.GetDepartmentTree(s.ctx, tenantID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Len(s.T(), tree, 2) // 两个根节点（dept-1 和 dept-2）

	// 验证子节点结构
	var dept1Node *DepartmentTreeNode
	for _, node := range tree {
		if node.DepartmentID == "dept-1" {
			dept1Node = node
			break
		}
	}
	assert.NotNil(s.T(), dept1Node)
	assert.Len(s.T(), dept1Node.Children, 1) // dept-1 有一个子节点
	assert.Equal(s.T(), "dept-1-1", dept1Node.Children[0].DepartmentID)

	s.mockDepartmentRepo.AssertExpectations(s.T())
}

// TestBuildDepartmentTree 测试构建部门树
func (s *DepartmentPermissionCheckerTestSuite) TestBuildDepartmentTree() {
	// Arrange
	parentID := "root"
	depts := []*entity.Department{
		{
			DepartmentID:       "dept-1",
			DepartmentName:     "Department 1",
			ParentDepartmentID: nil,
		},
		{
			DepartmentID:       "dept-1-1",
			DepartmentName:     "Department 1-1",
			ParentDepartmentID: func() *string { s := "dept-1"; return &s }(),
		},
		{
			DepartmentID:       "dept-1-1-1",
			DepartmentName:     "Department 1-1-1",
			ParentDepartmentID: func() *string { s := "dept-1-1"; return &s }(),
		},
	}

	// Act
	tree := s.checker.buildDepartmentTree(depts, "")

	// Assert
	assert.Len(s.T(), tree, 1)
	assert.Equal(s.T(), "dept-1", tree[0].DepartmentID)
	assert.Len(s.T(), tree[0].Children, 1)
	assert.Equal(s.T(), "dept-1-1", tree[0].Children[0].DepartmentID)
	assert.Len(s.T(), tree[0].Children[0].Children, 1)
	assert.Equal(s.T(), "dept-1-1-1", tree[0].Children[0].Children[0].DepartmentID)
}

// mockResource 用于测试的资源实现
type mockResource struct {
	deptID string
}

func (m *mockResource) GetDepartmentID() string {
	return m.deptID
}

// 运行测试套件
func TestDepartmentPermissionCheckerSuite(t *testing.T) {
	suite.Run(t, new(DepartmentPermissionCheckerTestSuite))
}

// TestDepartmentPermissionChecker_ErrorHandling 测试错误处理
func TestDepartmentPermissionChecker_ErrorHandling(t *testing.T) {
	t.Run("GetAccessibleDepartmentIDs returns error on repository failure", func(t *testing.T) {
		mockUserDeptRepo := new(MockUserDepartmentRepository)
		mockDepartmentRepo := new(MockDepartmentRepository)
		checker := NewDepartmentPermissionChecker(mockUserDeptRepo, mockDepartmentRepo)
		ctx := context.Background()

		mockUserDeptRepo.On("GetByUser", ctx, "user-123", "tenant-123").
			Return(([]*entity.UserDepartment)(nil), errors.New("database error"))

		_, err := checker.GetAccessibleDepartmentIDs(ctx, "user-123", "tenant-123")

		assert.Error(t, err)
		mockUserDeptRepo.AssertExpectations(t)
	})
}
