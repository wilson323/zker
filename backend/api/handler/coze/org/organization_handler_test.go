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
	"errors"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	"github.com/coze-dev/coze-studio/backend/api/model/org"
	orgentity "github.com/coze-dev/coze-studio/backend/domain/org/entity"
	orgservice "github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// MockOrganizationService 模拟组织服务
// 实现所有必要的接口方法用于测试
type MockOrganizationService struct {
	mock.Mock
}

// CreateOrganization 模拟创建组织
func (m *MockOrganizationService) CreateOrganization(ctx context.Context, req *orgservice.CreateOrganizationRequest) (*orgentity.Organization, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Organization), args.Error(1)
}

// GetOrganization 模拟获取组织
func (m *MockOrganizationService) GetOrganization(ctx context.Context, orgID string) (*orgentity.Organization, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Organization), args.Error(1)
}

// UpdateOrganization 模拟更新组织
func (m *MockOrganizationService) UpdateOrganization(ctx context.Context, req *orgservice.UpdateOrganizationRequest) (*orgentity.Organization, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Organization), args.Error(1)
}

// DeleteOrganization 模拟删除组织
func (m *MockOrganizationService) DeleteOrganization(ctx context.Context, orgID string) error {
	args := m.Called(ctx, orgID)
	return args.Error(0)
}

// GetOrganizationTree 模拟获取组织树
func (m *MockOrganizationService) GetOrganizationTree(ctx context.Context, tenantID string) ([]*orgentity.Organization, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Organization), args.Error(1)
}

// ListOrganizations 模拟分页查询组织
func (m *MockOrganizationService) ListOrganizations(ctx context.Context, filter *orgservice.OrganizationFilter) ([]*orgentity.Organization, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*orgentity.Organization), args.Get(1).(int64), args.Error(2)
}

// MoveOrganization 模拟移动组织
func (m *MockOrganizationService) MoveOrganization(ctx context.Context, req *orgservice.MoveOrganizationRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// GetChildren 模拟获取子组织
func (m *MockOrganizationService) GetChildren(ctx context.Context, parentID string) ([]*orgentity.Organization, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Organization), args.Error(1)
}

// GetAncestors 模拟获取祖先组织
func (m *MockOrganizationService) GetAncestors(ctx context.Context, orgID string) ([]*orgentity.Organization, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Organization), args.Error(1)
}

// GetDescendants 模拟获取后代组织
func (m *MockOrganizationService) GetDescendants(ctx context.Context, orgID string) ([]*orgentity.Organization, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Organization), args.Error(1)
}

// ==================== 测试辅助函数 ====================

// createTestRequestContext 创建测试用的请求上下文
func createTestRequestContext() *app.RequestContext {
	c := &app.RequestContext{}
	c.Request.SetMethod("POST")
	c.Request.SetRequestURI("/api/organizations")
	return c
}

// createMockOrg 创建模拟组织实体
func createMockOrg() *orgentity.Organization {
	return &orgentity.Organization{
		OrgID:       "org_001",
		TenantID:    "tenant_001",
		OrgName:     "技术部",
		OrgType:     orgentity.OrgTypeDepartment,
		OrgCode:     "TECH",
		Level:       2,
		Path:        "/root/org_001",
		SortOrder:   1,
		Status:      orgentity.OrgStatusActive,
		Description: "技术研发部门",
		LeaderID:    strPtr("leader_001"),
		CreatedAt:   1704067200000,
		UpdatedAt:   1704067200000,
	}
}

// strPtr 辅助函数:字符串指针
func strPtr(s string) *string {
	return &s
}

// ==================== CreateOrganization 测试 ====================

// TestCreateOrganization_Success 测试成功创建组织
func TestCreateOrganization_Success(t *testing.T) {
	// 准备测试数据
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := createTestRequestContext()

	req := org.CreateOrganizationRequest{
		TenantID:    "tenant_001",
		OrgName:     "技术部",
		OrgType:     "department",
		OrgCode:     "TECH",
		Description: "技术研发部门",
		SortOrder:   1,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedOrg := createMockOrg()

	// 设置Mock期望
	mockService.On("CreateOrganization", ctx, mock.AnythingOfType("*orgservice.CreateOrganizationRequest")).Return(expectedOrg, nil)

	// 执行测试
	handler.CreateOrganization(ctx, c)

	// 验证响应
	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.CreateOrganizationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "org_001", resp.Data.OrgID)
	assert.Equal(t, "技术部", resp.Data.OrgName)
	assert.Equal(t, "department", resp.Data.OrgType)

	mockService.AssertExpectations(t)
}

// TestCreateOrganization_InvalidParam 测试参数验证失败
func TestCreateOrganization_InvalidParam(t *testing.T) {
	testCases := []struct {
		name        string
		req         org.CreateOrganizationRequest
		expectCode  int32
		expectMsg   string
	}{
		{
			name: "缺少必填字段-tenant_id",
			req: org.CreateOrganizationRequest{
				OrgName: "技术部",
				OrgType: "department",
				OrgCode: "TECH",
			},
			expectCode: coze.InvalidParam,
			expectMsg:  "invalid parameter",
		},
		{
			name: "缺少必填字段-org_name",
			req: org.CreateOrganizationRequest{
				TenantID: "tenant_001",
				OrgType:  "department",
				OrgCode:  "TECH",
			},
			expectCode: coze.InvalidParam,
			expectMsg:  "invalid parameter",
		},
		{
			name: "组织类型无效",
			req: org.CreateOrganizationRequest{
				TenantID: "tenant_001",
				OrgName:  "技术部",
				OrgType:  "invalid_type",
				OrgCode:  "TECH",
			},
			expectCode: coze.InvalidParam,
			expectMsg:  "invalid parameter",
		},
		{
			name: "组织名称过长",
			req: org.CreateOrganizationRequest{
				TenantID: "tenant_001",
				OrgName:  string(make([]byte, 201)), // 超过200字符
				OrgType:  "department",
				OrgCode:  "TECH",
			},
			expectCode: coze.InvalidParam,
			expectMsg:  "invalid parameter",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			handler := NewOrganizationHandler(mockService)

			ctx := context.Background()
			c := createTestRequestContext()

			body, _ := json.Marshal(tc.req)
			c.Request.SetBody(body)

			// 执行测试
			handler.CreateOrganization(ctx, c)

			// 验证响应
			assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

			var resp map[string]interface{}
			err := json.Unmarshal(c.Response.Body(), &resp)
			require.NoError(t, err)
			assert.Contains(t, resp, "code")
		})
	}
}

// TestCreateOrganization_ServiceError 测试服务层返回错误
func TestCreateOrganization_ServiceError(t *testing.T) {
	testCases := []struct {
		name        string
		serviceErr  error
		expectHTTP  int
	}{
		{
			name:       "组织已存在",
			serviceErr: errno.ErrOrgAlreadyExists,
			expectHTTP: http.StatusConflict,
		},
		{
			name:       "组织编码已存在",
			serviceErr: errno.ErrOrgCodeExists,
			expectHTTP: http.StatusConflict,
		},
		{
			name:       "父组织不存在",
			serviceErr: errno.ErrParentOrgNotFound,
			expectHTTP: http.StatusNotFound,
		},
		{
			name:       "组织层级超出限制",
			serviceErr: errno.ErrOrgLevelExceeded,
			expectHTTP: http.StatusBadRequest,
		},
		{
			name:       "租户不匹配",
			serviceErr: errno.ErrTenantMismatch,
			expectHTTP: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			handler := NewOrganizationHandler(mockService)

			ctx := context.Background()
			c := createTestRequestContext()

			req := org.CreateOrganizationRequest{
				TenantID: "tenant_001",
				OrgName:  "技术部",
				OrgType:  "department",
				OrgCode:  "TECH",
			}

			body, _ := json.Marshal(req)
			c.Request.SetBody(body)

			mockService.On("CreateOrganization", ctx, mock.Anything).Return(nil, tc.serviceErr)

			handler.CreateOrganization(ctx, c)

			assert.Equal(t, tc.expectHTTP, c.Response.StatusCode())

			mockService.AssertExpectations(t)
		})
	}
}

// ==================== GetOrganization 测试 ====================

// TestGetOrganization_Success 测试成功获取组织
func TestGetOrganization_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	expectedOrg := createMockOrg()

	mockService.On("GetOrganization", ctx, "org_001").Return(expectedOrg, nil)

	handler.GetOrganization(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetOrganizationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "org_001", resp.Data.OrgID)
	assert.Equal(t, "技术部", resp.Data.OrgName)

	mockService.AssertExpectations(t)
}

// TestGetOrganization_MissingID 测试缺少组织ID参数
func TestGetOrganization_MissingID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置id参数

	handler.GetOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// TestGetOrganization_NotFound 测试组织不存在
func TestGetOrganization_NotFound(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "nonexistent_org")

	mockService.On("GetOrganization", ctx, "nonexistent_org").Return(nil, errno.ErrOrgNotFound)

	handler.GetOrganization(ctx, c)

	assert.Equal(t, http.StatusNotFound, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== UpdateOrganization 测试 ====================

// TestUpdateOrganization_Success 测试成功更新组织
func TestUpdateOrganization_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	req := org.UpdateOrganizationRequest{
		OrgName:     "新技术部",
		Description: "更新后的描述",
		SortOrder:   2,
		Status:      "active",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	updatedOrg := createMockOrg()
	updatedOrg.OrgName = "新技术部"
	updatedOrg.Description = "更新后的描述"
	updatedOrg.SortOrder = 2

	mockService.On("UpdateOrganization", ctx, mock.AnythingOfType("*orgservice.UpdateOrganizationRequest")).Return(updatedOrg, nil)

	handler.UpdateOrganization(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.UpdateOrganizationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "新技术部", resp.Data.OrgName)

	mockService.AssertExpectations(t)
}

// TestUpdateOrganization_InvalidParam 测试参数验证失败
func TestUpdateOrganization_InvalidParam(t *testing.T) {
	testCases := []struct {
		name       string
		setupReq   func(*app.RequestContext)
		expectCode int
	}{
		{
			name: "缺少组织ID",
			setupReq: func(c *app.RequestContext) {
				// 不设置id参数
			},
			expectCode: http.StatusBadRequest,
		},
		{
			name: "状态值无效",
			setupReq: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
				req := org.UpdateOrganizationRequest{
					Status: "invalid_status",
				}
				body, _ := json.Marshal(req)
				c.Request.SetBody(body)
			},
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			handler := NewOrganizationHandler(mockService)

			ctx := context.Background()
			c := createTestRequestContext()
			tc.setupReq(c)

			handler.UpdateOrganization(ctx, c)

			assert.Equal(t, tc.expectCode, c.Response.StatusCode())
		})
	}
}

// TestUpdateOrganization_NotFound 测试更新不存在的组织
func TestUpdateOrganization_NotFound(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "nonexistent_org")

	req := org.UpdateOrganizationRequest{
		OrgName: "更新名称",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("UpdateOrganization", ctx, mock.Anything).Return(nil, errno.ErrOrgNotFound)

	handler.UpdateOrganization(ctx, c)

	assert.Equal(t, http.StatusNotFound, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== DeleteOrganization 测试 ====================

// TestDeleteOrganization_Success 测试成功删除组织
func TestDeleteOrganization_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	mockService.On("DeleteOrganization", ctx, "org_001").Return(nil)

	handler.DeleteOrganization(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.DeleteOrganizationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "organization deleted successfully", resp.Message)

	mockService.AssertExpectations(t)
}

// TestDeleteOrganization_MissingID 测试缺少组织ID
func TestDeleteOrganization_MissingID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	handler.DeleteOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// TestDeleteOrganization_NotFound 测试删除不存在的组织
func TestDeleteOrganization_NotFound(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "nonexistent_org")

	mockService.On("DeleteOrganization", ctx, "nonexistent_org").Return(errno.ErrOrgNotFound)

	handler.DeleteOrganization(ctx, c)

	assert.Equal(t, http.StatusNotFound, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestDeleteOrganization_HasChildren 测试删除有子组织的组织
func TestDeleteOrganization_HasChildren(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	mockService.On("DeleteOrganization", ctx, "org_001").Return(errno.ErrOrgHasChildren)

	handler.DeleteOrganization(ctx, c)

	assert.Equal(t, http.StatusForbidden, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== GetOrganizationTree 测试 ====================

// TestGetOrganizationTree_Success 测试成功获取组织树
func TestGetOrganizationTree_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Add("tenant_id", "tenant_001")

	orgTree := []*orgentity.Organization{
		createMockOrg(),
		{
			OrgID:    "org_002",
			TenantID: "tenant_001",
			OrgName:  "市场部",
			OrgType:  orgentity.OrgTypeDepartment,
			OrgCode:  "MARKET",
			Level:    2,
			Path:     "/root/org_002",
			Status:   orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetOrganizationTree", ctx, "tenant_001").Return(orgTree, nil)

	handler.GetOrganizationTree(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetOrganizationTreeResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "技术部", resp.Data[0].OrgName)
	assert.Equal(t, "市场部", resp.Data[1].OrgName)

	mockService.AssertExpectations(t)
}

// TestGetOrganizationTree_MissingTenantID 测试缺少租户ID
func TestGetOrganizationTree_MissingTenantID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置tenant_id

	handler.GetOrganizationTree(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// TestGetOrganizationTree_EmptyTree 测试获取空组织树
func TestGetOrganizationTree_EmptyTree(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Add("tenant_id", "tenant_001")

	mockService.On("GetOrganizationTree", ctx, "tenant_001").Return([]*orgentity.Organization{}, nil)

	handler.GetOrganizationTree(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetOrganizationTreeResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 0, len(resp.Data))

	mockService.AssertExpectations(t)
}

// ==================== ListOrganizations 测试 ====================

// TestListOrganizations_Success 测试成功分页查询组织
func TestListOrganizations_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Add("tenant_id", "tenant_001")
	c.QueryArgs().Add("page", "1")
	c.QueryArgs().Add("page_size", "20")

	orgs := []*orgentity.Organization{
		createMockOrg(),
	}
	total := int64(1)

	mockService.On("ListOrganizations", ctx, mock.AnythingOfType("*orgservice.OrganizationFilter")).Return(orgs, total, nil)

	handler.ListOrganizations(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.ListOrganizationsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.PageSize)

	mockService.AssertExpectations(t)
}

// TestListOrganizations_BoundaryConditions 测试分页边界条件
func TestListOrganizations_BoundaryConditions(t *testing.T) {
	testCases := []struct {
		name            string
		page            string
		pageSize        string
		expectedPage    int
		expectedPageSize int
	}{
		{
			name:            "page=0自动修正为1",
			page:            "0",
			pageSize:        "20",
			expectedPage:    1,
			expectedPageSize: 20,
		},
		{
			name:            "page<0自动修正为1",
			page:            "-1",
			pageSize:        "20",
			expectedPage:    1,
			expectedPageSize: 20,
		},
		{
			name:            "pageSize>100自动修正为20",
			page:            "1",
			pageSize:        "101",
			expectedPage:    1,
			expectedPageSize: 20,
		},
		{
			name:            "pageSize=0自动修正为20",
			page:            "1",
			pageSize:        "0",
			expectedPage:    1,
			expectedPageSize: 20,
		},
		{
			name:            "pageSize<0自动修正为20",
			page:            "1",
			pageSize:        "-1",
			expectedPage:    1,
			expectedPageSize: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			handler := NewOrganizationHandler(mockService)

			ctx := context.Background()
			c := &app.RequestContext{}
			c.QueryArgs().Add("tenant_id", "tenant_001")
			c.QueryArgs().Add("page", tc.page)
			c.QueryArgs().Add("page_size", tc.pageSize)

			mockService.On("ListOrganizations", ctx, mock.MatchedBy(func(filter *orgservice.OrganizationFilter) bool {
				return filter.Page == tc.expectedPage && filter.PageSize == tc.expectedPageSize
			})).Return([]*orgentity.Organization{}, int64(0), nil)

			handler.ListOrganizations(ctx, c)

			mockService.AssertExpectations(t)
		})
	}
}

// TestListOrganizations_WithFilters 测试带过滤条件的查询
func TestListOrganizations_WithFilters(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.QueryArgs().Add("tenant_id", "tenant_001")
	c.QueryArgs().Add("org_type", "department")
	c.QueryArgs().Add("status", "active")
	c.QueryArgs().Add("parent_id", "parent_001")
	c.QueryArgs().Add("keyword", "技术")
	c.QueryArgs().Add("sort_by", "org_name")
	c.QueryArgs().Add("sort_order", "asc")

	mockService.On("ListOrganizations", ctx, mock.MatchedBy(func(filter *orgservice.OrganizationFilter) bool {
		return filter.TenantID == "tenant_001" &&
			filter.OrgType == "department" &&
			filter.Status == "active" &&
			filter.ParentID == "parent_001" &&
			filter.Keyword == "技术" &&
			filter.SortBy == "org_name" &&
			filter.SortOrder == "asc"
	})).Return([]*orgentity.Organization{}, int64(0), nil)

	handler.ListOrganizations(ctx, c)

	mockService.AssertExpectations(t)
}

// TestListOrganizations_MissingTenantID 测试缺少租户ID
func TestListOrganizations_MissingTenantID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	// 不设置tenant_id

	handler.ListOrganizations(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// ==================== MoveOrganization 测试 ====================

// TestMoveOrganization_Success 测试成功移动组织
func TestMoveOrganization_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	newParentID := "new_parent_001"
	req := org.MoveOrganizationRequest{
		OrgID:       "org_001",
		NewParentID: &newParentID,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveOrganization", ctx, mock.AnythingOfType("*orgservice.MoveOrganizationRequest")).Return(nil)

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.MoveOrganizationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "organization moved successfully", resp.Message)

	mockService.AssertExpectations(t)
}

// TestMoveOrganization_MoveToRoot 测试移动到根节点
func TestMoveOrganization_MoveToRoot(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	req := org.MoveOrganizationRequest{
		OrgID:       "org_001",
		NewParentID: nil, // 移动到根节点
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveOrganization", ctx, mock.MatchedBy(func(req *orgservice.MoveOrganizationRequest) bool {
		return req.OrgID == "org_001" && req.NewParentID == nil
	})).Return(nil)

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestMoveOrganization_CannotMoveToSelf 测试不能移动到自己
func TestMoveOrganization_CannotMoveToSelf(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	sameParentID := "org_001"
	req := org.MoveOrganizationRequest{
		OrgID:       "org_001",
		NewParentID: &sameParentID,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveOrganization", ctx, mock.Anything).Return(errno.ErrCannotMoveToSelf)

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestMoveOrganization_CannotMoveToDescendant 测试不能移动到后代节点
func TestMoveOrganization_CannotMoveToDescendant(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	descendantID := "org_001_child"
	req := org.MoveOrganizationRequest{
		OrgID:       "org_001",
		NewParentID: &descendantID,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveOrganization", ctx, mock.Anything).Return(errno.ErrCannotMoveToDescendant)

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestMoveOrganization_CycleDetected 测试检测到循环依赖
func TestMoveOrganization_CycleDetected(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	newParentID := "parent_001"
	req := org.MoveOrganizationRequest{
		OrgID:       "org_001",
		NewParentID: &newParentID,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveOrganization", ctx, mock.Anything).Return(errno.ErrOrgCycleDetected)

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestMoveOrganization_LevelExceeded 测试组织层级超出限制
func TestMoveOrganization_LevelExceeded(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	newParentID := "deep_parent"
	req := org.MoveOrganizationRequest{
		OrgID:       "org_001",
		NewParentID: &newParentID,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("MoveOrganization", ctx, mock.Anything).Return(errno.ErrOrgLevelExceeded)

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestMoveOrganization_MissingID 测试缺少组织ID
func TestMoveOrganization_MissingID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	handler.MoveOrganization(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// ==================== GetChildren 测试 ====================

// TestGetChildren_Success 测试成功获取子组织
func TestGetChildren_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	children := []*orgentity.Organization{
		{
			OrgID:    "child_001",
			TenantID: "tenant_001",
			OrgName:  "子组织1",
			OrgType:  orgentity.OrgTypeDepartment,
			OrgCode:  "CHILD1",
			Level:    3,
			Path:     "/root/org_001/child_001",
			Status:   orgentity.OrgStatusActive,
		},
		{
			OrgID:    "child_002",
			TenantID: "tenant_001",
			OrgName:  "子组织2",
			OrgType:  orgentity.OrgTypeDepartment,
			OrgCode:  "CHILD2",
			Level:    3,
			Path:     "/root/org_001/child_002",
			Status:   orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetChildren", ctx, "org_001").Return(children, nil)

	handler.GetChildren(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetChildrenResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "子组织1", resp.Data[0].OrgName)
	assert.Equal(t, "子组织2", resp.Data[1].OrgName)

	mockService.AssertExpectations(t)
}

// TestGetChildren_Empty 测试获取空子组织列表
func TestGetChildren_Empty(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	mockService.On("GetChildren", ctx, "org_001").Return([]*orgentity.Organization{}, nil)

	handler.GetChildren(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetChildrenResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 0, len(resp.Data))

	mockService.AssertExpectations(t)
}

// TestGetChildren_MissingID 测试缺少组织ID
func TestGetChildren_MissingID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	handler.GetChildren(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// ==================== GetAncestors 测试 ====================

// TestGetAncestors_Success 测试成功获取祖先组织
func TestGetAncestors_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	ancestors := []*orgentity.Organization{
		{
			OrgID:    "root",
			TenantID: "tenant_001",
			OrgName:  "根组织",
			OrgType:  orgentity.OrgTypeCompany,
			OrgCode:  "ROOT",
			Level:    1,
			Path:     "/root",
			Status:   orgentity.OrgStatusActive,
		},
		{
			OrgID:    "parent_001",
			TenantID: "tenant_001",
			OrgName:  "父组织",
			OrgType:  orgentity.OrgTypeDivision,
			OrgCode:  "PARENT",
			Level:    2,
			Path:     "/root/parent_001",
			Status:   orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetAncestors", ctx, "org_001").Return(ancestors, nil)

	handler.GetAncestors(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetAncestorsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "根组织", resp.Data[0].OrgName)
	assert.Equal(t, "父组织", resp.Data[1].OrgName)

	mockService.AssertExpectations(t)
}

// TestGetAncestors_Empty 测试获取空祖先列表
func TestGetAncestors_Empty(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	mockService.On("GetAncestors", ctx, "org_001").Return([]*orgentity.Organization{}, nil)

	handler.GetAncestors(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetAncestorsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 0, len(resp.Data))

	mockService.AssertExpectations(t)
}

// TestGetAncestors_MissingID 测试缺少组织ID
func TestGetAncestors_MissingID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	handler.GetAncestors(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// ==================== GetDescendants 测试 ====================

// TestGetDescendants_Success 测试成功获取后代组织
func TestGetDescendants_Success(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	descendants := []*orgentity.Organization{
		{
			OrgID:    "child_001",
			TenantID: "tenant_001",
			OrgName:  "子组织1",
			OrgType:  orgentity.OrgTypeDepartment,
			OrgCode:  "CHILD1",
			Level:    3,
			Path:     "/root/org_001/child_001",
			Status:   orgentity.OrgStatusActive,
		},
		{
			OrgID:    "grandchild_001",
			TenantID: "tenant_001",
			OrgName:  "孙组织1",
			OrgType:  orgentity.OrgTypeProject,
			OrgCode:  "GRANDCHILD1",
			Level:    4,
			Path:     "/root/org_001/child_001/grandchild_001",
			Status:   orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetDescendants", ctx, "org_001").Return(descendants, nil)

	handler.GetDescendants(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetDescendantsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "子组织1", resp.Data[0].OrgName)
	assert.Equal(t, "孙组织1", resp.Data[1].OrgName)

	mockService.AssertExpectations(t)
}

// TestGetDescendants_Empty 测试获取空后代列表
func TestGetDescendants_Empty(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}
	c.Params.Append("id", "org_001")

	mockService.On("GetDescendants", ctx, "org_001").Return([]*orgentity.Organization{}, nil)

	handler.GetDescendants(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp org.GetDescendantsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	require.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 0, len(resp.Data))

	mockService.AssertExpectations(t)
}

// TestGetDescendants_MissingID 测试缺少组织ID
func TestGetDescendants_MissingID(t *testing.T) {
	mockService := new(MockOrganizationService)
	handler := NewOrganizationHandler(mockService)

	ctx := context.Background()
	c := &app.RequestContext{}

	handler.GetDescendants(ctx, c)

	assert.Equal(t, http.StatusBadRequest, c.Response.StatusCode())
}

// ==================== 数据转换函数测试 ====================

// TestToOrganizationData 测试数据转换
func TestToOrganizationData(t *testing.T) {
	entity := &orgentity.Organization{
		OrgID:       "org_001",
		TenantID:    "tenant_001",
		OrgName:     "技术部",
		OrgType:     orgentity.OrgTypeDepartment,
		OrgCode:     "TECH",
		Level:       2,
		Path:        "/root/org_001",
		SortOrder:   1,
		Status:      orgentity.OrgStatusActive,
		Description: "技术研发部门",
		LeaderID:    strPtr("leader_001"),
		CreatedAt:   1704067200000,
		UpdatedAt:   1704067200000,
	}

	data := toOrganizationData(entity)

	assert.NotNil(t, data)
	assert.Equal(t, "org_001", data.OrgID)
	assert.Equal(t, "tenant_001", data.TenantID)
	assert.Equal(t, "技术部", data.OrgName)
	assert.Equal(t, "department", data.OrgType)
	assert.Equal(t, "TECH", data.OrgCode)
	assert.Equal(t, 2, data.Level)
	assert.Equal(t, "/root/org_001", data.Path)
	assert.Equal(t, "active", data.Status)
	assert.Equal(t, "leader_001", *data.LeaderID)
}

// TestToOrganizationDataList 测试批量数据转换
func TestToOrganizationDataList(t *testing.T) {
	entities := []*orgentity.Organization{
		{
			OrgID:   "org_001",
			OrgName: "技术部",
			OrgType: orgentity.OrgTypeDepartment,
			Status:  orgentity.OrgStatusActive,
		},
		{
			OrgID:   "org_002",
			OrgName: "市场部",
			OrgType: orgentity.OrgTypeDepartment,
			Status:  orgentity.OrgStatusActive,
		},
	}

	dataList := toOrganizationDataList(entities)

	assert.NotNil(t, dataList)
	assert.Equal(t, 2, len(dataList))
	assert.Equal(t, "org_001", dataList[0].OrgID)
	assert.Equal(t, "org_002", dataList[1].OrgID)
}

// TestToOrganizationData_Nil 测试nil实体转换
func TestToOrganizationData_Nil(t *testing.T) {
	data := toOrganizationData(nil)
	assert.Nil(t, data)
}

// TestToOrganizationDataList_Nil 测试nil列表转换
func TestToOrganizationDataList_Nil(t *testing.T) {
	data := toOrganizationDataList(nil)
	assert.Nil(t, data)
}

// TestToOrganizationDataList_Empty 测试空列表转换
func TestToOrganizationDataList_Empty(t *testing.T) {
	data := toOrganizationDataList([]*orgentity.Organization{})
	assert.NotNil(t, data)
	assert.Equal(t, 0, len(data))
}

// ==================== 错误场景综合测试 ====================

// TestOrganizationHandler_ServiceErrors 测试所有服务层错误
func TestOrganizationHandler_ServiceErrors(t *testing.T) {
	testCases := []struct {
		name           string
		handlerFunc    func(*OrganizationHandler, context.Context, *app.RequestContext)
		setupMock      func(*MockOrganizationService)
		setupContext   func(*app.RequestContext)
		expectedStatus int
	}{
		{
			name: "CreateOrganization-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.CreateOrganization(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("CreateOrganization", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				req := org.CreateOrganizationRequest{
					TenantID: "tenant_001",
					OrgName:  "技术部",
					OrgType:  "department",
					OrgCode:  "TECH",
				}
				body, _ := json.Marshal(req)
				c.Request.SetBody(body)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "GetOrganization-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.GetOrganization(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("GetOrganization", mock.Anything, "org_001").Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "UpdateOrganization-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.UpdateOrganization(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("UpdateOrganization", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
				req := org.UpdateOrganizationRequest{OrgName: "更新"}
				body, _ := json.Marshal(req)
				c.Request.SetBody(body)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "DeleteOrganization-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.DeleteOrganization(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("DeleteOrganization", mock.Anything, "org_001").Return(errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "GetOrganizationTree-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.GetOrganizationTree(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("GetOrganizationTree", mock.Anything, "tenant_001").Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.QueryArgs().Add("tenant_id", "tenant_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "ListOrganizations-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.ListOrganizations(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("ListOrganizations", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.QueryArgs().Add("tenant_id", "tenant_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "MoveOrganization-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.MoveOrganization(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("MoveOrganization", mock.Anything, mock.Anything).Return(errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
				newParentID := "new_parent"
				req := org.MoveOrganizationRequest{NewParentID: &newParentID}
				body, _ := json.Marshal(req)
				c.Request.SetBody(body)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "GetChildren-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.GetChildren(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("GetChildren", mock.Anything, "org_001").Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "GetAncestors-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.GetAncestors(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("GetAncestors", mock.Anything, "org_001").Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "GetDescendants-内部错误",
			handlerFunc: func(h *OrganizationHandler, ctx context.Context, c *app.RequestContext) {
				h.GetDescendants(ctx, c)
			},
			setupMock: func(m *MockOrganizationService) {
				m.On("GetDescendants", mock.Anything, "org_001").Return(nil, errors.New("internal error"))
			},
			setupContext: func(c *app.RequestContext) {
				c.Params.Append("id", "org_001")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockOrganizationService)
			handler := NewOrganizationHandler(mockService)

			ctx := context.Background()
			c := createTestRequestContext()
			tc.setupContext(c)
			tc.setupMock(mockService)

			tc.handlerFunc(handler, ctx, c)

			assert.Equal(t, tc.expectedStatus, c.Response.StatusCode())

			mockService.AssertExpectations(t)
		})
	}
}

// ==================== 并发测试 ====================

// TestToOrganizationDataList_Concurrent 测试并发安全的数据转换
func TestToOrganizationDataList_Concurrent(t *testing.T) {
	entities := make([]*orgentity.Organization, 100)
	for i := 0; i < 100; i++ {
		entities[i] = &orgentity.Organization{
			OrgID:   "org_" + string(rune(i)),
			OrgName: "组织" + string(rune(i)),
			OrgType: orgentity.OrgTypeDepartment,
			Status:  orgentity.OrgStatusActive,
		}
	}

	// 并发执行多次转换
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			data := toOrganizationDataList(entities)
			assert.Equal(t, 100, len(data))
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}
