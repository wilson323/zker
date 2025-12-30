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
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// MockDepartmentService 模拟部门服务
// 职责：为DepartmentHandler提供测试用的模拟服务层
type MockDepartmentService struct {
	mock.Mock
}

// CreateDepartment 模拟创建部门
func (m *MockDepartmentService) CreateDepartment(ctx context.Context, req *service.CreateDepartmentRequest) (*entity.Department, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Department), args.Error(1)
}

// GetDepartment 模拟获取部门
func (m *MockDepartmentService) GetDepartment(ctx context.Context, deptID string) (*entity.Department, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Department), args.Error(1)
}

// UpdateDepartment 模拟更新部门
func (m *MockDepartmentService) UpdateDepartment(ctx context.Context, req *service.UpdateDepartmentRequest) (*entity.Department, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Department), args.Error(1)
}

// DeleteDepartment 模拟删除部门
func (m *MockDepartmentService) DeleteDepartment(ctx context.Context, deptID string) error {
	args := m.Called(ctx, deptID)
	return args.Error(0)
}

// GetDepartmentTree 模拟获取部门树
func (m *MockDepartmentService) GetDepartmentTree(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

// ListDepartments 模拟分页查询部门列表
func (m *MockDepartmentService) ListDepartments(ctx context.Context, filter *repository.DepartmentFilter) ([]*entity.Department, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Department), args.Get(1).(int64), args.Error(2)
}

// MoveDepartment 模拟移动部门
func (m *MockDepartmentService) MoveDepartment(ctx context.Context, req *service.MoveDepartmentRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// GetChildren 模拟获取子部门
func (m *MockDepartmentService) GetChildren(ctx context.Context, parentID string) ([]*entity.Department, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

// GetAncestors 模拟获取祖先部门
func (m *MockDepartmentService) GetAncestors(ctx context.Context, deptID string) ([]*entity.Department, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

// GetDescendants 模拟获取后代部门
func (m *MockDepartmentService) GetDescendants(ctx context.Context, deptID string) ([]*entity.Department, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

// GetDepartmentsByOrg 模拟获取组织的所有部门
func (m *MockDepartmentService) GetDepartmentsByOrg(ctx context.Context, orgID string) ([]*entity.Department, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Department), args.Error(1)
}

// TestCreateDepartment_Success 测试成功创建部门
func TestCreateDepartment_Success(t *testing.T) {
	// 准备测试数据
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	leaderID := "leader_001"
	req := CreateDepartmentRequest{
		TenantID:    "tenant_001",
		OrgID:       "org_001",
		DeptName:    "技术部",
		DeptCode:    "TECH",
		LeaderID:    &leaderID,
		Description: "技术研发部门",
		SortOrder:   1,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedDept := &entity.Department{
		DeptID:      "dept_001",
		TenantID:    "tenant_001",
		OrgID:       "org_001",
		DeptName:    "技术部",
		DeptCode:    "TECH",
		Level:       1,
		Path:        "/dept_001",
		SortOrder:   1,
		Status:      entity.OrgStatusActive,
		LeaderID:    &leaderID,
		Description: "技术研发部门",
	}

	// 设置Mock期望
	mockService.On("CreateDepartment", ctx, mock.AnythingOfType("*service.CreateDepartmentRequest")).Return(expectedDept, nil)

	// 执行测试
	handler.CreateDepartment(ctx, c)

	// 验证响应
	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))
	assert.Equal(t, "success", resp["message"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "dept_001", data["dept_id"])
	assert.Equal(t, "技术部", data["dept_name"])

	mockService.AssertExpectations(t)
}

// TestCreateDepartment_InvalidParams 测试创建部门参数验证失败
func TestCreateDepartment_InvalidParams(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	// 测试缺少必填字段
	req := CreateDepartmentRequest{
		TenantID: "tenant_001",
		// 缺少 OrgID, DeptName, DeptCode
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.CreateDepartment(ctx, c)

	// 验证响应
	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestCreateDepartment_ServiceError 测试创建部门服务层错误
func TestCreateDepartment_ServiceError(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	req := CreateDepartmentRequest{
		TenantID: "tenant_001",
		OrgID:    "org_001",
		DeptName: "技术部",
		DeptCode: "TECH",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	// 设置Mock期望返回错误
	mockService.On("CreateDepartment", ctx, mock.Anything).Return(nil, errno.ErrDeptCodeAlreadyExists)

	handler.CreateDepartment(ctx, c)

	// 验证响应
	assert.Equal(t, errno.ErrDeptCodeAlreadyExists.HTTPStatus(), c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, errno.ErrDeptCodeAlreadyExists.Code(), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetDepartment_Success 测试成功获取部门详情
func TestGetDepartment_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	expectedDept := &entity.Department{
		DeptID:   "dept_001",
		TenantID: "tenant_001",
		OrgID:    "org_001",
		DeptName: "技术部",
		DeptCode: "TECH",
		Level:    1,
		Path:     "/dept_001",
		Status:   entity.OrgStatusActive,
	}

	mockService.On("GetDepartment", ctx, "dept_001").Return(expectedDept, nil)

	handler.GetDepartment(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "dept_001", data["dept_id"])

	mockService.AssertExpectations(t)
}

// TestGetDepartment_MissingID 测试获取部门缺少ID参数
func TestGetDepartment_MissingID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置ID参数

	handler.GetDepartment(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestGetDepartment_NotFound 测试获取部门不存在
func TestGetDepartment_NotFound(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_not_exist")

	mockService.On("GetDepartment", ctx, "dept_not_exist").Return(nil, errno.ErrDeptNotFound)

	handler.GetDepartment(ctx, c)

	assert.Equal(t, errno.ErrDeptNotFound.HTTPStatus(), c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, errno.ErrDeptNotFound.Code(), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestUpdateDepartment_Success 测试成功更新部门
func TestUpdateDepartment_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	leaderID := "leader_002"
	req := UpdateDepartmentRequest{
		DeptName:    "研发部",
		LeaderID:    &leaderID,
		Description: "核心研发部门",
		Status:      "active",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedDept := &entity.Department{
		DeptID:      "dept_001",
		DeptName:    "研发部",
		LeaderID:    &leaderID,
		Description: "核心研发部门",
		Status:      entity.OrgStatusActive,
	}

	mockService.On("UpdateDepartment", ctx, mock.AnythingOfType("*service.UpdateDepartmentRequest")).Return(expectedDept, nil)

	handler.UpdateDepartment(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "研发部", data["dept_name"])

	mockService.AssertExpectations(t)
}

// TestUpdateDepartment_InvalidStatus 测试更新部门状态无效
func TestUpdateDepartment_InvalidStatus(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	req := UpdateDepartmentRequest{
		DeptName: "研发部",
		Status:   "invalid_status", // 无效状态
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.UpdateDepartment(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestDeleteDepartment_Success 测试成功删除部门
func TestDeleteDepartment_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	mockService.On("DeleteDepartment", ctx, "dept_001").Return(nil)

	handler.DeleteDepartment(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestDeleteDepartment_HasChildren 测试删除有子部门的部门
func TestDeleteDepartment_HasChildren(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	mockService.On("DeleteDepartment", ctx, "dept_001").Return(errno.ErrDeptHasChildren)

	handler.DeleteDepartment(ctx, c)

	assert.Equal(t, errno.ErrDeptHasChildren.HTTPStatus(), c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, errno.ErrDeptHasChildren.Code(), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetDepartmentTree_Success 测试成功获取部门树
func TestGetDepartmentTree_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Add("tenant_id", "tenant_001")

	expectedTree := []*entity.Department{
		{
			DeptID:   "dept_001",
			DeptName: "技术部",
			Level:    1,
			Children: []entity.Department{
				{
					DeptID:   "dept_002",
					DeptName: "前端组",
					Level:    2,
				},
			},
		},
	}

	mockService.On("GetDepartmentTree", ctx, "tenant_001").Return(expectedTree, nil)

	handler.GetDepartmentTree(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetDepartmentTree_MissingTenantID 测试获取部门树缺少租户ID
func TestGetDepartmentTree_MissingTenantID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置 tenant_id

	handler.GetDepartmentTree(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestListDepartments_Success 测试成功分页查询部门列表
func TestListDepartments_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Set("tenant_id", "tenant_001")
	c.QueryArgs().Set("page", "1")
	c.QueryArgs().Set("page_size", "20")

	expectedDepts := []*entity.Department{
		{DeptID: "dept_001", DeptName: "技术部"},
		{DeptID: "dept_002", DeptName: "市场部"},
	}

	mockService.On("ListDepartments", ctx, mock.AnythingOfType("*repository.DepartmentFilter")).Return(expectedDepts, int64(2), nil)

	handler.ListDepartments(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(2), data["total"])
	assert.NotNil(t, data["departments"])

	mockService.AssertExpectations(t)
}

// TestListDepartments_WithFilter 测试带过滤条件的分页查询
func TestListDepartments_WithFilter(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Set("tenant_id", "tenant_001")
	c.QueryArgs().Set("org_id", "org_001")
	c.QueryArgs().Set("status", "active")
	c.QueryArgs().Set("parent_id", "parent_001")
	c.QueryArgs().Set("level", "2")
	c.QueryArgs().Set("keyword", "技术")
	c.QueryArgs().Set("page", "1")
	c.QueryArgs().Set("page_size", "10")

	expectedDepts := []*entity.Department{
		{DeptID: "dept_001", DeptName: "技术部"},
	}

	mockService.On("ListDepartments", ctx, mock.MatchedBy(func(filter *repository.DepartmentFilter) bool {
		return filter.TenantID == "tenant_001" &&
			filter.OrgID == "org_001" &&
			filter.Status == "active" &&
			*filter.ParentID == "parent_001" &&
			*filter.Level == 2 &&
			filter.Keyword == "技术" &&
			filter.Page == 1 &&
			filter.PageSize == 10
	})).Return(expectedDepts, int64(1), nil)

	handler.ListDepartments(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestListDepartments_DefaultPagination 测试使用默认分页参数
func TestListDepartments_DefaultPagination(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Set("tenant_id", "tenant_001")
	// 不设置page和page_size

	mockService.On("ListDepartments", ctx, mock.MatchedBy(func(filter *repository.DepartmentFilter) bool {
		return filter.Page == 1 && filter.PageSize == 20
	})).Return([]*entity.Department{}, int64(0), nil)

	handler.ListDepartments(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestMoveDepartment_Success 测试成功移动部门
func TestMoveDepartment_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	newParentID := "new_parent_001"
	req := MoveDepartmentRequest{
		NewParentID: &newParentID,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveDepartment", ctx, mock.AnythingOfType("*service.MoveDepartmentRequest")).Return(nil)

	handler.MoveDepartment(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestMoveDepartment_MissingID 测试移动部门缺少ID
func TestMoveDepartment_MissingID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置ID

	handler.MoveDepartment(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestMoveDepartment_CannotMoveToSelf 测试不能移动到自己
func TestMoveDepartment_CannotMoveToSelf(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	deptID := "dept_001"
	req := MoveDepartmentRequest{
		NewParentID: &deptID, // 尝试移动到自己
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveDepartment", ctx, mock.Anything).Return(errno.ErrCannotMoveToSelf)

	handler.MoveDepartment(ctx, c)

	assert.Equal(t, errno.ErrCannotMoveToSelf.HTTPStatus(), c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, errno.ErrCannotMoveToSelf.Code(), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetChildren_Success 测试成功获取子部门
func TestGetChildren_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	expectedChildren := []*entity.Department{
		{DeptID: "dept_002", DeptName: "前端组"},
		{DeptID: "dept_003", DeptName: "后端组"},
	}

	mockService.On("GetChildren", ctx, "dept_001").Return(expectedChildren, nil)

	handler.GetChildren(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetChildren_MissingID 测试获取子部门缺少ID
func TestGetChildren_MissingID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置ID

	handler.GetChildren(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestGetAncestors_Success 测试成功获取祖先部门
func TestGetAncestors_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_003")

	expectedAncestors := []*entity.Department{
		{DeptID: "dept_001", DeptName: "技术部", Level: 1},
		{DeptID: "dept_002", DeptName: "研发中心", Level: 2},
	}

	mockService.On("GetAncestors", ctx, "dept_003").Return(expectedAncestors, nil)

	handler.GetAncestors(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetAncestors_MissingID 测试获取祖先部门缺少ID
func TestGetAncestors_MissingID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置ID

	handler.GetAncestors(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestGetDescendants_Success 测试成功获取后代部门
func TestGetDescendants_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "dept_001")

	expectedDescendants := []*entity.Department{
		{DeptID: "dept_002", DeptName: "前端组", Level: 2},
		{DeptID: "dept_003", DeptName: "前端一组", Level: 3},
	}

	mockService.On("GetDescendants", ctx, "dept_001").Return(expectedDescendants, nil)

	handler.GetDescendants(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetDescendants_MissingID 测试获取后代部门缺少ID
func TestGetDescendants_MissingID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置ID

	handler.GetDescendants(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestGetDepartmentsByOrg_Success 测试成功获取组织的所有部门
func TestGetDepartmentsByOrg_Success(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("org_id", "org_001")

	expectedDepts := []*entity.Department{
		{DeptID: "dept_001", DeptName: "技术部"},
		{DeptID: "dept_002", DeptName: "市场部"},
		{DeptID: "dept_003", DeptName: "人力资源部"},
	}

	mockService.On("GetDepartmentsByOrg", ctx, "org_001").Return(expectedDepts, nil)

	handler.GetDepartmentsByOrg(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), int32(resp["code"].(float64)))

	mockService.AssertExpectations(t)
}

// TestGetDepartmentsByOrg_MissingOrgID 测试获取组织部门缺少组织ID
func TestGetDepartmentsByOrg_MissingOrgID(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置org_id

	handler.GetDepartmentsByOrg(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.NotEqual(t, int32(0), int32(resp["code"].(float64)))
}

// TestHandleError_WithErrorCode 测试错误处理（带错误码）
func TestHandleError_WithErrorCode(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	handler.handleError(ctx, c, errno.ErrDeptNotFound)

	assert.Equal(t, errno.ErrDeptNotFound.HTTPStatus(), c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, errno.ErrDeptNotFound.Code(), int32(resp["code"].(float64)))
	assert.Equal(t, errno.ErrDeptNotFound.Message(), resp["message"])
}

// TestHandleError_WithGenericError 测试错误处理（通用错误）
func TestHandleError_WithGenericError(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	genericErr := assert.AnError

	handler.handleError(ctx, c, genericErr)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	var resp map[string]interface{}
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, errno.InternalError.Code(), int32(resp["code"].(float64)))
}

// TestSuccess 测试成功响应
func TestSuccess(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := Success(data)

	assert.Equal(t, int32(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
	assert.Equal(t, data, resp["data"])
}

// TestSuccess_NilData 测试成功响应（nil数据）
func TestSuccess_NilData(t *testing.T) {
	resp := Success(nil)

	assert.Equal(t, int32(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
	assert.Nil(t, resp["data"])
}

// TestFail 测试失败响应
func TestFail(t *testing.T) {
	code := int32(400)
	message := "Bad Request"
	resp := Fail(code, message)

	assert.Equal(t, code, resp["code"])
	assert.Equal(t, message, resp["message"])
	assert.Nil(t, resp["data"])
}

// TestParseStringPtr_ValidString 测试解析非空字符串指针
func TestParseStringPtr_ValidString(t *testing.T) {
	s := "test_string"
	ptr := parseStringPtr(s)

	assert.NotNil(t, ptr)
	assert.Equal(t, s, *ptr)
}

// TestParseStringPtr_EmptyString 测试解析空字符串指针
func TestParseStringPtr_EmptyString(t *testing.T) {
	ptr := parseStringPtr("")

	assert.Nil(t, ptr)
}

// TestParseIntPtr_PositiveInt 测试解析正整数指针
func TestParseIntPtr_PositiveInt(t *testing.T) {
	i := 5
	ptr := parseIntPtr(i)

	assert.NotNil(t, ptr)
	assert.Equal(t, i, *ptr)
}

// TestParseIntPtr_ZeroInt 测试解析零整数指针
func TestParseIntPtr_ZeroInt(t *testing.T) {
	ptr := parseIntPtr(0)

	assert.Nil(t, ptr)
}

// TestParseIntPtr_NegativeInt 测试解析负整数指针
func TestParseIntPtr_NegativeInt(t *testing.T) {
	ptr := parseIntPtr(-1)

	assert.Nil(t, ptr)
}

// TestDepartmentHandlerIntegration 集成测试：完整流程
func TestDepartmentHandlerIntegration(t *testing.T) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()

	// 1. 创建部门
	t.Run("CreateDepartment", func(t *testing.T) {
		c := &app.RequestContext{}
		leaderID := "leader_001"
		req := CreateDepartmentRequest{
			TenantID:    "tenant_001",
			OrgID:       "org_001",
			DeptName:    "技术部",
			DeptCode:    "TECH",
			LeaderID:    &leaderID,
			Description: "技术研发部门",
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		createdDept := &entity.Department{
			DeptID:      "dept_001",
			TenantID:    "tenant_001",
			OrgID:       "org_001",
			DeptName:    "技术部",
			DeptCode:    "TECH",
			Level:       1,
			Path:        "/dept_001",
			Status:      entity.OrgStatusActive,
			LeaderID:    &leaderID,
			Description: "技术研发部门",
		}

		mockService.On("CreateDepartment", ctx, mock.Anything).Return(createdDept, nil).Once()

		handler.CreateDepartment(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
		var resp map[string]interface{}
		json.Unmarshal(c.Response.Body(), &resp)
		assert.Equal(t, int32(0), int32(resp["code"].(float64)))
	})

	// 2. 获取部门
	t.Run("GetDepartment", func(t *testing.T) {
		c := &app.RequestContext{}
		c.Params.Append("id", "dept_001")

		dept := &entity.Department{
			DeptID:   "dept_001",
			DeptName: "技术部",
		}

		mockService.On("GetDepartment", ctx, "dept_001").Return(dept, nil).Once()

		handler.GetDepartment(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})

	// 3. 更新部门
	t.Run("UpdateDepartment", func(t *testing.T) {
		c := &app.RequestContext{}
		c.Params.Append("id", "dept_001")

		req := UpdateDepartmentRequest{
			DeptName: "研发部",
			Status:   "active",
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		updatedDept := &entity.Department{
			DeptID:   "dept_001",
			DeptName: "研发部",
			Status:   entity.OrgStatusActive,
		}

		mockService.On("UpdateDepartment", ctx, mock.Anything).Return(updatedDept, nil).Once()

		handler.UpdateDepartment(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})

	// 4. 删除部门
	t.Run("DeleteDepartment", func(t *testing.T) {
		c := &app.RequestContext{}
		c.Params.Append("id", "dept_001")

		mockService.On("DeleteDepartment", ctx, "dept_001").Return(nil).Once()

		handler.DeleteDepartment(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})

	mockService.AssertExpectations(t)
}

// BenchmarkCreateDepartment 性能测试：创建部门
func BenchmarkCreateDepartment(b *testing.B) {
	mockService := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockService)

	ctx := context.Background()
	leaderID := "leader_001"

	req := CreateDepartmentRequest{
		TenantID:    "tenant_001",
		OrgID:       "org_001",
		DeptName:    "技术部",
		DeptCode:    "TECH",
		LeaderID:    &leaderID,
		Description: "技术研发部门",
	}

	expectedDept := &entity.Department{
		DeptID:   "dept_001",
		DeptName: "技术部",
		DeptCode: "TECH",
	}

	mockService.On("CreateDepartment", ctx, mock.Anything).Return(expectedDept, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := &app.RequestContext{}
		body, _ := json.Marshal(req)
		c.Request.SetBody(body)
		handler.CreateDepartment(ctx, c)
	}
}
