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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
)

// MockEmployeeRepository 员工仓储Mock
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, emp *entity.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, empID string) (*entity.Employee, error) {
	args := m.Called(ctx, empID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Employee, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(2)
	}
	return args.Get(0).(*entity.Employee), args.Error(2)
}

func (m *MockEmployeeRepository) GetByUserID(ctx context.Context, userID string) (*entity.Employee, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Employee, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByDepartmentID(ctx context.Context, deptID string) ([]*entity.Employee, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByOrgID(ctx context.Context, orgID string) ([]*entity.Employee, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Update(ctx context.Context, emp *entity.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepository) UpdateStatus(ctx context.Context, empID string, status entity.EmployeeStatus) error {
	args := m.Called(ctx, empID, status)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, empID string) error {
	args := m.Called(ctx, empID)
	return args.Error(0)
}

func (m *MockEmployeeRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	args := m.Called(ctx, tenantID, code, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmployeeRepository) ExistsByEmail(ctx context.Context, tenantID, email string, excludeID string) (bool, error) {
	args := m.Called(ctx, tenantID, email, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmployeeRepository) List(ctx context.Context, filter *repository.EmployeeFilter) ([]*entity.Employee, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Employee), args.Int64(1), args.Error(2)
}

func (m *MockEmployeeRepository) Search(ctx context.Context, tenantID, keyword string, limit int) ([]*entity.Employee, error) {
	args := m.Called(ctx, tenantID, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(2)
	}
	return args.Get(0).([]*entity.Employee), args.Error(2)
}

func (m *MockEmployeeRepository) GetByPinyin(ctx context.Context, tenantID, pinyin string) ([]*entity.Employee, error) {
	args := m.Called(ctx, tenantID, pinyin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Employee), args.Error(1)
}

// MockDepartmentRepository 部门仓储Mock
type MockDepartmentRepository struct {
	mock.Mock
}

func (m *MockDepartmentRepository) Create(ctx context.Context, dept *entity.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepository) GetByID(ctx context.Context, deptID string) (*entity.Department, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Department, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(2)
	}
	return args.Get(0).(*entity.Department), args.Error(2)
}

func (m *MockDepartmentRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetTree(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

func (m *MockDepartmentRepository) GetChildren(ctx context.Context, parentID string) ([]*entity.Department, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

func (m *MockDepartmentRepository) Update(ctx context.Context, dept *entity.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepository) Delete(ctx context.Context, deptID string) error {
	args := m.Called(ctx, deptID)
	return args.Error(0)
}

func (m *MockDepartmentRepository) Move(ctx context.Context, deptID, newParentID string) error {
	args := m.Called(ctx, deptID, newParentID)
	return args.Error(0)
}

func (m *MockDepartmentRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	args := m.Called(ctx, tenantID, code, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepository) List(ctx context.Context, filter *repository.DepartmentFilter) ([]*entity.Department, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Department), args.Int64(1), args.Error(2)
}

// MockPositionRepository 岗位仓储Mock
type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) Create(ctx context.Context, pos *entity.Position) error {
	args := m.Called(ctx, pos)
	return args.Error(0)
}

func (m *MockPositionRepository) GetByID(ctx context.Context, positionID string) (*entity.Position, error) {
	args := m.Called(ctx, positionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByCode(ctx context.Context, tenantID, code string) (*entity.Position, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(2)
	}
	return args.Get(0).(*entity.Position), args.Error(2)
}

func (m *MockPositionRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Position, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByDepartmentID(ctx context.Context, deptID string) ([]*entity.Position, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Position), args.Error(1)
}

func (m *MockPositionRepository) Update(ctx context.Context, pos *entity.Position) error {
	args := m.Called(ctx, pos)
	return args.Error(0)
}

func (m *MockPositionRepository) Delete(ctx context.Context, positionID string) error {
	args := m.Called(ctx, positionID)
	return args.Error(0)
}

func (m *MockPositionRepository) ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error) {
	args := m.Called(ctx, tenantID, code, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPositionRepository) List(ctx context.Context, filter *repository.PositionFilter) ([]*entity.Position, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*entity.Position), args.Int64(1), args.Error(2)
}

// TestAddressBookService_GetAddressBook 测试获取通讯录
func TestAddressBookService_GetAddressBook(t *testing.T) {
	ctx := context.Background()

	mockEmpRepo := new(MockEmployeeRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockPosRepo := new(MockPositionRepository)

	service := NewAddressBookService(mockEmpRepo, mockDeptRepo, mockPosRepo)

	// 准备测试数据
	now := time.Now().UnixMilli()
	deptID := "dept001"
	positionID := "pos001"

	employees := []*entity.Employee{
		{
			EmpID:         "emp001",
			TenantID:      "tenant001",
			OrgID:         "org001",
			DeptID:        &deptID,
			PositionID:    &positionID,
			EmpName:       "张三",
			EmpCode:       "E001",
			JobLevel:      5,
			JobTitle:      "软件工程师",
			EmployeeType:  entity.EmpTypeFullTime,
			EmployeeStatus: entity.EmpStatusActive,
			Email:         stringPtr("zhangsan@example.com"),
			Phone:         stringPtr("13800138000"),
			CreatedAt:     now,
			UpdatedAt:     now,
			Department: &entity.Department{
				DeptID:   deptID,
				DeptName: "研发部",
			},
			Position: &entity.Position{
				PositionID:   positionID,
				PositionName: "P5工程师",
			},
		},
		{
			EmpID:         "emp002",
			TenantID:      "tenant001",
			OrgID:         "org001",
			DeptID:        &deptID,
			PositionID:    &positionID,
			EmpName:       "李四",
			EmpCode:       "E002",
			JobLevel:      6,
			JobTitle:      "高级软件工程师",
			EmployeeType:  entity.EmpTypeFullTime,
			EmployeeStatus: entity.EmpStatusActive,
			Email:         stringPtr("lisi@example.com"),
			Phone:         stringPtr("13800138001"),
			CreatedAt:     now,
			UpdatedAt:     now,
			Department: &entity.Department{
				DeptID:   deptID,
				DeptName: "研发部",
			},
			Position: &entity.Position{
				PositionID:   positionID,
				PositionName: "P6工程师",
			},
		},
	}

	// 设置Mock期望
	mockEmpRepo.On("List", ctx, mock.AnythingOfType("*repository.EmployeeFilter")).Return(employees, int64(2), nil)

	// 执行测试
	req := &GetAddressBookRequest{
		TenantID: "tenant001",
		Page:     1,
		PageSize: 20,
	}

	resp, err := service.GetAddressBook(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Contacts))
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, "张三", resp.Contacts[0].Name)
	assert.Equal(t, "研发部", resp.Contacts[0].DeptName)
	assert.Equal(t, "P5工程师", resp.Contacts[0].Position)

	mockEmpRepo.AssertExpectations(t)
}

// TestAddressBookService_SearchContacts 测试搜索联系人
func TestAddressBookService_SearchContacts(t *testing.T) {
	ctx := context.Background()

	mockEmpRepo := new(MockEmployeeRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockPosRepo := new(MockPositionRepository)

	service := NewAddressBookService(mockEmpRepo, mockDeptRepo, mockPosRepo)

	// 准备测试数据
	now := time.Now().UnixMilli()
	deptID := "dept001"
	positionID := "pos001"

	employees := []*entity.Employee{
		{
			EmpID:         "emp001",
			TenantID:      "tenant001",
			OrgID:         "org001",
			DeptID:        &deptID,
			PositionID:    &positionID,
			EmpName:       "张三",
			EmpCode:       "E001",
			JobLevel:      5,
			JobTitle:      "软件工程师",
			EmployeeType:  entity.EmpTypeFullTime,
			EmployeeStatus: entity.EmpStatusActive,
			Email:         stringPtr("zhangsan@example.com"),
			Phone:         stringPtr("13800138000"),
			CreatedAt:     now,
			UpdatedAt:     now,
			Department: &entity.Department{
				DeptID:   deptID,
				DeptName: "研发部",
			},
			Position: &entity.Position{
				PositionID:   positionID,
				PositionName: "P5工程师",
			},
		},
	}

	// 设置Mock期望
	mockEmpRepo.On("Search", ctx, "tenant001", "张三", 20).Return(employees, nil)

	// 执行测试
	req := &SearchContactsRequest{
		TenantID: "tenant001",
		Keyword:  "张三",
		Limit:    20,
	}

	contacts, err := service.SearchContacts(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, contacts)
	assert.Equal(t, 1, len(contacts))
	assert.Equal(t, "张三", contacts[0].Name)

	mockEmpRepo.AssertExpectations(t)
}

// TestAddressBookService_GetContactsByDepartment 测试按部门获取联系人
func TestAddressBookService_GetContactsByDepartment(t *testing.T) {
	ctx := context.Background()

	mockEmpRepo := new(MockEmployeeRepository)
	mockDeptRepo := new(MockDepartmentRepository)
	mockPosRepo := new(MockPositionRepository)

	service := NewAddressBookService(mockEmpRepo, mockDeptRepo, mockPosRepo)

	// 准备测试数据
	deptID := "dept001"
	positionID := "pos001"
	now := time.Now().UnixMilli()

	department := &entity.Department{
		DeptID:   deptID,
		TenantID: "tenant001",
		DeptName: "研发部",
	}

	employees := []*entity.Employee{
		{
			EmpID:         "emp001",
			TenantID:      "tenant001",
			OrgID:         "org001",
			DeptID:        &deptID,
			PositionID:    &positionID,
			EmpName:       "张三",
			EmpCode:       "E001",
			JobLevel:      5,
			JobTitle:      "软件工程师",
			EmployeeType:  entity.EmpTypeFullTime,
			EmployeeStatus: entity.EmpStatusActive,
			Email:         stringPtr("zhangsan@example.com"),
			CreatedAt:     now,
			UpdatedAt:     now,
			Department:    department,
		},
	}

	// 设置Mock期望
	mockDeptRepo.On("GetByID", ctx, deptID).Return(department, nil)
	mockEmpRepo.On("GetByDepartmentID", ctx, deptID).Return(employees, nil)

	// 执行测试
	req := &GetContactsByDepartmentRequest{
		TenantID:     "tenant001",
		DepartmentID: deptID,
		Page:         1,
		PageSize:     20,
	}

	resp, err := service.GetContactsByDepartment(ctx, req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Contacts))
	assert.Equal(t, "张三", resp.Contacts[0].Name)
	assert.Equal(t, "研发部", resp.Contacts[0].DeptName)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

// 辅助函数
func stringPtr(s string) *string {
	return &s
}
