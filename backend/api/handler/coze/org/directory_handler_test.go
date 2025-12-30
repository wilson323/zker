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

package org

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	orgentity "github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
)

// MockDirectoryService 模拟通讯录服务
type MockDirectoryService struct {
	mock.Mock
}

// GetOrganizationDirectory 模拟获取组织架构目录
func (m *MockDirectoryService) GetOrganizationDirectory(ctx context.Context, tenantID string) (*service.OrganizationTreeNode, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.OrganizationTreeNode), args.Error(1)
}

// GetDepartmentDirectory 模拟获取部门目录
func (m *MockDirectoryService) GetDepartmentDirectory(ctx context.Context, tenantID string) (*service.DepartmentTreeNode, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.DepartmentTreeNode), args.Error(1)
}

// GetDepartmentEmployees 模拟获取部门员工
func (m *MockDirectoryService) GetDepartmentEmployees(ctx context.Context, deptID string) ([]*service.EmployeeListItem, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.EmployeeListItem), args.Error(1)
}

// GetOrganizationEmployees 模拟获取组织员工
func (m *MockDirectoryService) GetOrganizationEmployees(ctx context.Context, orgID string) ([]*service.EmployeeListItem, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.EmployeeListItem), args.Error(1)
}

// SearchEmployees 模拟搜索员工
func (m *MockDirectoryService) SearchEmployees(ctx context.Context, tenantID, keyword string, limit int) ([]*service.EmployeeListItem, error) {
	args := m.Called(ctx, tenantID, keyword, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*service.EmployeeListItem), args.Error(1)
}

// GetEmployeeByCode 模拟根据工号获取员工
func (m *MockDirectoryService) GetEmployeeByCode(ctx context.Context, tenantID, code string) (*service.EmployeeListItem, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.EmployeeListItem), args.Error(1)
}

// ==================== 组织目录接口测试 ====================

// TestGetOrganizationDirectory_Success 测试成功获取组织架构目录
func TestGetOrganizationDirectory_Success(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/organization?tenant_id=tenant_001")

	expectedResult := &service.OrganizationTreeNode{
		OrgID:         "org_001",
		OrgName:       "技术公司",
		OrgType:       orgentity.OrgTypeCompany,
		Level:         1,
		EmployeeCount: 100,
		Children: []*service.OrganizationTreeNode{
			{
				OrgID:         "org_002",
				OrgName:       "技术部",
				OrgType:       orgentity.OrgTypeDepartment,
				Level:         2,
				EmployeeCount: 50,
			},
		},
	}

	mockService.On("GetOrganizationDirectory", ctx, "tenant_001").Return(expectedResult, nil)

	// 注入模拟服务
	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetOrganizationDirectory(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// TestGetOrganizationDirectory_ServiceNotInitialized 测试服务未初始化
func TestGetOrganizationDirectory_ServiceNotInitialized(t *testing.T) {
	ctx := context.Background()
	c := &app.RequestContext{}

	oldService := directoryService
	directoryService = nil
	defer func() { directoryService = oldService }()

	GetOrganizationDirectory(ctx, c)

	// 应该返回500错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetOrganizationDirectory_MissingTenantID 测试缺少tenant_id参数
func TestGetOrganizationDirectory_MissingTenantID(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置tenant_id

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetOrganizationDirectory(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, float64(0), resp["code"])
}

// TestGetOrganizationDirectory_ServiceError 测试服务层错误
func TestGetOrganizationDirectory_ServiceError(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/organization?tenant_id=tenant_001")

	mockService.On("GetOrganizationDirectory", ctx, "tenant_001").Return(
		nil,
		assert.AnError,
	)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetOrganizationDirectory(ctx, c)

	// 应该返回500错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== 部门目录接口测试 ====================

// TestGetDepartmentDirectory_Success 测试成功获取部门目录
func TestGetDepartmentDirectory_Success(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/department?tenant_id=tenant_001")

	expectedResult := &service.DepartmentTreeNode{
		DeptID:        "dept_001",
		DeptName:      "研发中心",
		Level:         1,
		EmployeeCount: 80,
		LeaderID:      stringPtr("emp_001"),
		LeaderName:    stringPtr("张三"),
		Children: []*service.DepartmentTreeNode{
			{
				DeptID:        "dept_002",
				DeptName:      "后端开发组",
				Level:         2,
				EmployeeCount: 30,
			},
		},
	}

	mockService.On("GetDepartmentDirectory", ctx, "tenant_001").Return(expectedResult, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetDepartmentDirectory(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])
	assert.NotNil(t, resp["data"])

	mockService.AssertExpectations(t)
}

// TestGetDepartmentDirectory_MissingTenantID 测试缺少tenant_id
func TestGetDepartmentDirectory_MissingTenantID(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetDepartmentDirectory(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// ==================== 部门员工列表接口测试 ====================

// TestGetDepartmentEmployees_Success 测试成功获取部门员工列表
func TestGetDepartmentEmployees_Success(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("dept_id", "dept_001")

	expectedEmployees := []*service.EmployeeListItem{
		{
			EmpID:          "emp_001",
			EmpName:        "张三",
			EmpCode:        "E001",
			OrgName:        "技术公司",
			DeptName:       stringPtr("研发中心"),
			JobTitle:       "高级工程师",
			JobLevel:       7,
			Phone:          stringPtr("13800138000"),
			Email:          stringPtr("zhangsan@example.com"),
			EmployeeStatus: orgentity.EmpStatusActive,
		},
		{
			EmpID:          "emp_002",
			EmpName:        "李四",
			EmpCode:        "E002",
			OrgName:        "技术公司",
			DeptName:       stringPtr("研发中心"),
			JobTitle:       "工程师",
			JobLevel:       5,
			EmployeeStatus: orgentity.EmpStatusTrial,
		},
	}

	mockService.On("GetDepartmentEmployees", ctx, "dept_001").Return(expectedEmployees, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetDepartmentEmployees(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["employees"])

	mockService.AssertExpectations(t)
}

// TestGetDepartmentEmployees_MissingDeptID 测试缺少dept_id参数
func TestGetDepartmentEmployees_MissingDeptID(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置dept_id

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetDepartmentEmployees(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetDepartmentEmployees_EmptyList 测试空员工列表
func TestGetDepartmentEmployees_EmptyList(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("dept_id", "dept_001")

	emptyList := []*service.EmployeeListItem{}
	mockService.On("GetDepartmentEmployees", ctx, "dept_001").Return(emptyList, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetDepartmentEmployees(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])

	mockService.AssertExpectations(t)
}

// ==================== 组织员工列表接口测试 ====================

// TestGetOrganizationEmployees_Success 测试成功获取组织员工列表
func TestGetOrganizationEmployees_Success(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("org_id", "org_001")

	expectedEmployees := []*service.EmployeeListItem{
		{
			EmpID:          "emp_001",
			EmpName:        "张三",
			EmpCode:        "E001",
			OrgName:        "技术公司",
			JobTitle:       "高级工程师",
			JobLevel:       7,
			EmployeeStatus: orgentity.EmpStatusActive,
		},
		{
			EmpID:          "emp_002",
			EmpName:        "李四",
			EmpCode:        "E002",
			OrgName:        "技术公司",
			JobTitle:       "产品经理",
			JobLevel:       6,
			EmployeeStatus: orgentity.EmpStatusActive,
		},
		{
			EmpID:          "emp_003",
			EmpName:        "王五",
			EmpCode:        "E003",
			OrgName:        "技术公司",
			JobTitle:       "设计师",
			JobLevel:       5,
			EmployeeStatus: orgentity.EmpStatusActive,
		},
	}

	mockService.On("GetOrganizationEmployees", ctx, "org_001").Return(expectedEmployees, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetOrganizationEmployees(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(3), data["total"])

	employees := data["employees"].([]interface{})
	assert.Equal(t, 3, len(employees))

	mockService.AssertExpectations(t)
}

// TestGetOrganizationEmployees_MissingOrgID 测试缺少org_id参数
func TestGetOrganizationEmployees_MissingOrgID(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetOrganizationEmployees(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// ==================== 员工搜索接口测试 ====================

// TestSearchEmployees_Success 测试成功搜索员工
func TestSearchEmployees_Success(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/search?keyword=张三&limit=10&tenant_id=tenant_001")

	expectedEmployees := []*service.EmployeeListItem{
		{
			EmpID:          "emp_001",
			EmpName:        "张三",
			EmpCode:        "E001",
			OrgName:        "技术公司",
			DeptName:       stringPtr("研发中心"),
			JobTitle:       "高级工程师",
			JobLevel:       7,
			Phone:          stringPtr("13800138000"),
			Email:          stringPtr("zhangsan@example.com"),
			EmployeeStatus: orgentity.EmpStatusActive,
		},
	}

	mockService.On("SearchEmployees", ctx, "tenant_001", "张三", 10).Return(expectedEmployees, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	SearchEmployees(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.Equal(t, "张三", data["keyword"])

	mockService.AssertExpectations(t)
}

// TestSearchEmployees_DefaultLimit 测试使用默认limit
func TestSearchEmployees_DefaultLimit(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/search?keyword=李四&tenant_id=tenant_001")
	// 不设置limit，应该使用默认值20

	expectedEmployees := []*service.EmployeeListItem{
		{
			EmpID:          "emp_002",
			EmpName:        "李四",
			EmpCode:        "E002",
			OrgName:        "技术公司",
			JobTitle:       "工程师",
			JobLevel:       5,
			EmployeeStatus: orgentity.EmpStatusActive,
		},
	}

	mockService.On("SearchEmployees", ctx, "tenant_001", "李四", 20).Return(expectedEmployees, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	SearchEmployees(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestSearchEmployees_MissingKeyword 测试缺少keyword参数
func TestSearchEmployees_MissingKeyword(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/search?tenant_id=tenant_001")
	// 不设置keyword

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	SearchEmployees(ctx, c)

	// 应该返回400错误（BindAndValidate失败）
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestSearchEmployees_MissingTenantID 测试缺少tenant_id
func TestSearchEmployees_MissingTenantID(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/search?keyword=张三")
	// 不设置tenant_id

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	SearchEmployees(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestSearchEmployees_NoResults 测试搜索无结果
func TestSearchEmployees_NoResults(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.URI().SetQueryString("/api/directory/search?keyword=不存在&tenant_id=tenant_001")

	emptyList := []*service.EmployeeListItem{}
	mockService.On("SearchEmployees", ctx, "tenant_001", "不存在", 20).Return(emptyList, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	SearchEmployees(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(0), data["total"])

	mockService.AssertExpectations(t)
}

// ==================== 员工详情接口测试 ====================

// TestGetEmployeeByCode_Success 测试成功根据工号获取员工
func TestGetEmployeeByCode_Success(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("code", "E001")
	c.URI().SetQueryString("/api/directory/employee/by-code/E001?tenant_id=tenant_001")

	expectedEmployee := &service.EmployeeListItem{
		EmpID:          "emp_001",
		EmpName:        "张三",
		EmpCode:        "E001",
		OrgName:        "技术公司",
		DeptName:       stringPtr("研发中心"),
		PositionName:   stringPtr("高级工程师"),
		JobTitle:       "高级工程师",
		JobLevel:       7,
		Phone:          stringPtr("13800138000"),
		Email:          stringPtr("zhangsan@example.com"),
		AvatarURL:      stringPtr("https://example.com/avatar.jpg"),
		EmployeeStatus: orgentity.EmpStatusActive,
	}

	mockService.On("GetEmployeeByCode", ctx, "tenant_001", "E001").Return(expectedEmployee, nil)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetEmployeeByCode(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), resp["code"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "emp_001", data["emp_id"])
	assert.Equal(t, "张三", data["emp_name"])
	assert.Equal(t, "E001", data["emp_code"])

	mockService.AssertExpectations(t)
}

// TestGetEmployeeByCode_MissingCode 测试缺少code参数
func TestGetEmployeeByCode_MissingCode(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetEmployeeByCode(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetEmployeeByCode_MissingTenantID 测试缺少tenant_id
func TestGetEmployeeByCode_MissingTenantID(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("code", "E001")
	// 不设置tenant_id

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetEmployeeByCode(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetEmployeeByCode_NotFound 测试员工不存在
func TestGetEmployeeByCode_NotFound(t *testing.T) {
	mockService := new(MockDirectoryService)
	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("code", "E999")
	c.URI().SetQueryString("/api/directory/employee/by-code/E999?tenant_id=tenant_001")

	mockService.On("GetEmployeeByCode", ctx, "tenant_001", "E999").Return(
		nil,
		assert.AnError,
	)

	oldService := directoryService
	directoryService = mockService
	defer func() { directoryService = oldService }()

	GetEmployeeByCode(ctx, c)

	// 应该返回500错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== 辅助函数 ====================

// stringPtr 返回字符串指针的辅助函数
func stringPtr(s string) *string {
	return &s
}
