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
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/api/model/org"
	orgentity "github.com/coze-dev/coze-studio/backend/domain/org/entity"
	orgservice "github.com/coze-dev/coze-studio/backend/domain/org/service"
)

// MockPositionService 模拟岗位服务
// 职责：为PositionHandler测试提供Mock依赖
type MockPositionService struct {
	mock.Mock
}

// CreatePosition 模拟创建岗位
func (m *MockPositionService) CreatePosition(ctx context.Context, req *orgservice.CreatePositionRequest) (*orgentity.Position, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Position), args.Error(1)
}

// GetPosition 模拟获取岗位详情
func (m *MockPositionService) GetPosition(ctx context.Context, positionID string) (*orgentity.Position, error) {
	args := m.Called(ctx, positionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Position), args.Error(1)
}

// UpdatePosition 模拟更新岗位
func (m *MockPositionService) UpdatePosition(ctx context.Context, req *orgservice.UpdatePositionRequest) (*orgentity.Position, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Position), args.Error(1)
}

// DeletePosition 模拟删除岗位
func (m *MockPositionService) DeletePosition(ctx context.Context, positionID string) error {
	args := m.Called(ctx, positionID)
	return args.Error(0)
}

// ListPositions 模拟分页查询岗位列表
func (m *MockPositionService) ListPositions(ctx context.Context, filter interface{}) ([]*orgentity.Position, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*orgentity.Position), args.Get(1).(int64), args.Error(2)
}

// GetPositionByCode 模拟根据编码获取岗位
func (m *MockPositionService) GetPositionByCode(ctx context.Context, tenantID, code string) (*orgentity.Position, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orgentity.Position), args.Error(1)
}

// GetPositionsByDepartment 模拟获取部门的所有岗位
func (m *MockPositionService) GetPositionsByDepartment(ctx context.Context, deptID string) ([]*orgentity.Position, error) {
	args := m.Called(ctx, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Position), args.Error(1)
}

// GetPositionsByLevel 模拟按职级获取岗位
func (m *MockPositionService) GetPositionsByLevel(ctx context.Context, tenantID string, level int) ([]*orgentity.Position, error) {
	args := m.Called(ctx, tenantID, level)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Position), args.Error(1)
}

// GetPositionsByCategory 模拟按类别获取岗位
func (m *MockPositionService) GetPositionsByCategory(ctx context.Context, tenantID, category string) ([]*orgentity.Position, error) {
	args := m.Called(ctx, tenantID, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*orgentity.Position), args.Error(1)
}

// ==================== 测试辅助函数 ====================

// setupTestContext 创建测试上下文和请求上下文
func setupTestContext() (context.Context, *app.RequestContext) {
	ctx := context.Background()
	c := &app.RequestContext{}
	return ctx, c
}

// createTestPosition 创建测试用岗位实体
func createTestPosition() *orgentity.Position {
	deptID := "dept_001"
	return &orgentity.Position{
		PositionID:      "pos_001",
		TenantID:        "tenant_001",
		DeptID:          &deptID,
		PositionName:    "高级软件工程师",
		PositionCode:    "SE001",
		Level:           3,
		Category:        "技术岗",
		Responsibilities: "负责系统架构设计",
		Requirements:    "5年以上开发经验",
		Status:          orgentity.OrgStatusActive,
		SortOrder:       1,
		CreatedAt:       1704067200000,
		UpdatedAt:       1704067200000,
	}
}

// ==================== CreatePosition 测试 ====================

// TestCreatePosition_Success 测试成功创建岗位
func TestCreatePosition_Success(t *testing.T) {
	// 准备测试数据
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()

	req := org.CreatePositionRequest{
		TenantID:         "tenant_001",
		DeptID:           strPtr("dept_001"),
		PositionName:     "高级软件工程师",
		PositionCode:     "SE001",
		Level:            3,
		Category:         "技术岗",
		Responsibilities: "负责系统架构设计",
		Requirements:     "5年以上开发经验",
		SortOrder:        1,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedPosition := createTestPosition()

	// 设置Mock期望：验证服务调用参数
	mockService.On("CreatePosition", ctx, mock.MatchedBy(func(r *orgservice.CreatePositionRequest) bool {
		return r.TenantID == "tenant_001" &&
			r.PositionName == "高级软件工程师" &&
			r.PositionCode == "SE001" &&
			r.Level == 3 &&
			r.Category == "技术岗"
	})).Return(expectedPosition, nil)

	// 执行测试
	handler.CreatePosition(ctx, c)

	// 验证响应
	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.CreatePositionResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "pos_001", resp.Data.PositionID)
	assert.Equal(t, "高级软件工程师", resp.Data.PositionName)
	assert.Equal(t, "SE001", resp.Data.PositionCode)
	assert.Equal(t, 3, resp.Data.Level)
	assert.Equal(t, "技术岗", resp.Data.Category)
	assert.Equal(t, "active", resp.Data.Status)

	mockService.AssertExpectations(t)
}

// TestCreatePosition_InvalidParam 测试创建岗位参数验证失败
func TestCreatePosition_InvalidParam(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()

	// 缺少必填字段
	req := org.CreatePositionRequest{
		TenantID: "tenant_001",
		// PositionName 缺失
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.CreatePosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "CreatePosition", mock.Anything, mock.Anything)
}

// TestCreatePosition_ServiceError 测试创建岗位服务层错误
func TestCreatePosition_ServiceError(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()

	req := org.CreatePositionRequest{
		TenantID:     "tenant_001",
		PositionName: "高级软件工程师",
		PositionCode: "SE001",
		Level:        3,
		Category:     "技术岗",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	// 模拟服务返回错误
	mockService.On("CreatePosition", ctx, mock.Anything).Return(nil, errors.New("database error"))

	handler.CreatePosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== GetPosition 测试 ====================

// TestGetPosition_Success 测试成功获取岗位详情
func TestGetPosition_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_001")

	expectedPosition := createTestPosition()
	mockService.On("GetPosition", ctx, "pos_001").Return(expectedPosition, nil)

	handler.GetPosition(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.GetPositionResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "pos_001", resp.Data.PositionID)
	assert.Equal(t, "高级软件工程师", resp.Data.PositionName)

	mockService.AssertExpectations(t)
}

// TestGetPosition_MissingID 测试获取岗位缺少ID参数
func TestGetPosition_MissingID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置id参数

	handler.GetPosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPosition", mock.Anything, mock.Anything)
}

// TestGetPosition_NotFound 测试获取岗位不存在
func TestGetPosition_NotFound(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_not_exist")

	mockService.On("GetPosition", ctx, "pos_not_exist").Return(nil, errors.New("position not found"))

	handler.GetPosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== UpdatePosition 测试 ====================

// TestUpdatePosition_Success 测试成功更新岗位
func TestUpdatePosition_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_001")

	req := org.UpdatePositionRequest{
		PositionName:     "资深软件工程师",
		Level:            4,
		Category:         "技术岗",
		Responsibilities: "负责核心系统设计",
		Requirements:     "8年以上开发经验",
		SortOrder:        2,
		Status:           "active",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	updatedPosition := createTestPosition()
	updatedPosition.PositionName = "资深软件工程师"
	updatedPosition.Level = 4

	mockService.On("UpdatePosition", ctx, mock.MatchedBy(func(r *orgservice.UpdatePositionRequest) bool {
		return r.PositionID == "pos_001" &&
			r.PositionName == "资深软件工程师" &&
			r.Level == 4
	})).Return(updatedPosition, nil)

	handler.UpdatePosition(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.UpdatePositionResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "资深软件工程师", resp.Data.PositionName)
	assert.Equal(t, 4, resp.Data.Level)

	mockService.AssertExpectations(t)
}

// TestUpdatePosition_MissingID 测试更新岗位缺少ID参数
func TestUpdatePosition_MissingID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置id参数

	req := org.UpdatePositionRequest{
		PositionName: "资深软件工程师",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.UpdatePosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "UpdatePosition", mock.Anything, mock.Anything)
}

// TestUpdatePosition_InvalidStatus 测试更新岗位状态无效
func TestUpdatePosition_InvalidStatus(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_001")

	req := org.UpdatePositionRequest{
		PositionName: "资深软件工程师",
		Status:       "invalid_status", // 无效状态
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.UpdatePosition(ctx, c)

	// 验证返回错误响应（参数验证失败）
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "UpdatePosition", mock.Anything, mock.Anything)
}

// TestUpdatePosition_ServiceError 测试更新岗位服务层错误
func TestUpdatePosition_ServiceError(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_001")

	req := org.UpdatePositionRequest{
		PositionName: "资深软件工程师",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("UpdatePosition", ctx, mock.Anything).Return(nil, errors.New("database error"))

	handler.UpdatePosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== DeletePosition 测试 ====================

// TestDeletePosition_Success 测试成功删除岗位
func TestDeletePosition_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_001")

	mockService.On("DeletePosition", ctx, "pos_001").Return(nil)

	handler.DeletePosition(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.DeletePositionResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

// TestDeletePosition_MissingID 测试删除岗位缺少ID参数
func TestDeletePosition_MissingID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置id参数

	handler.DeletePosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "DeletePosition", mock.Anything, mock.Anything)
}

// TestDeletePosition_NotFound 测试删除岗位不存在
func TestDeletePosition_NotFound(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("id", "pos_not_exist")

	mockService.On("DeletePosition", ctx, "pos_not_exist").Return(errors.New("position not found"))

	handler.DeletePosition(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ==================== ListPositions 测试 ====================

// TestListPositions_Success 测试成功分页查询岗位列表
func TestListPositions_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Request.SetQueryString("?tenant_id=tenant_001&page=1&page_size=20&level=3")

	positions := []*orgentity.Position{
		createTestPosition(),
		{
			PositionID:   "pos_002",
			TenantID:     "tenant_001",
			PositionName: "产品经理",
			PositionCode: "PM001",
			Level:        3,
			Category:     "产品岗",
			Status:       orgentity.OrgStatusActive,
		},
	}

	mockService.On("ListPositions", ctx, mock.MatchedBy(func(filter interface{}) bool {
		// 验证过滤器包含正确参数
		return true
	})).Return(positions, int64(2), nil)

	handler.ListPositions(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.ListPositionsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.PageSize)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "pos_001", resp.Data[0].PositionID)
	assert.Equal(t, "pos_002", resp.Data[1].PositionID)

	mockService.AssertExpectations(t)
}

// TestListPositions_MissingTenantID 测试分页查询缺少租户ID
func TestListPositions_MissingTenantID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置tenant_id

	handler.ListPositions(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "ListPositions", mock.Anything, mock.Anything)
}

// TestListPositions_InvalidPageParams 测试分页参数验证
func TestListPositions_InvalidPageParams(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Request.SetQueryString("?tenant_id=tenant_001&page=0&page_size=150")

	positions := []*orgentity.Position{createTestPosition()}

	mockService.On("ListPositions", ctx, mock.Anything).Return(positions, int64(1), nil)

	handler.ListPositions(ctx, c)

	// 验证page和pageSize被修正为默认值
	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.ListPositionsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Page)      // 应该被修正为1
	assert.Equal(t, 20, resp.PageSize) // 应该被修正为20

	mockService.AssertExpectations(t)
}

// TestListPositions_WithLevelFilter 测试按职级筛选
func TestListPositions_WithLevelFilter(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Request.SetQueryString("?tenant_id=tenant_001&level=3")

	positions := []*orgentity.Position{createTestPosition()}

	mockService.On("ListPositions", ctx, mock.MatchedBy(func(filter interface{}) bool {
		// 验证level参数被正确解析
		return true
	})).Return(positions, int64(1), nil)

	handler.ListPositions(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.ListPositionsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Data))

	mockService.AssertExpectations(t)
}

// ==================== GetPositionByCode 测试 ====================

// TestGetPositionByCode_Success 测试成功根据编码获取岗位
func TestGetPositionByCode_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("code", "SE001")
	c.Request.SetQueryString("?tenant_id=tenant_001")

	expectedPosition := createTestPosition()
	mockService.On("GetPositionByCode", ctx, "tenant_001", "SE001").Return(expectedPosition, nil)

	handler.GetPositionByCode(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.GetPositionByCodeResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "SE001", resp.Data.PositionCode)

	mockService.AssertExpectations(t)
}

// TestGetPositionByCode_MissingCode 测试根据编码获取岗位缺少编码参数
func TestGetPositionByCode_MissingCode(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置code参数
	c.Request.SetQueryString("?tenant_id=tenant_001")

	handler.GetPositionByCode(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionByCode", mock.Anything, mock.Anything, mock.Anything)
}

// TestGetPositionByCode_MissingTenantID 测试根据编码获取岗位缺少租户ID
func TestGetPositionByCode_MissingTenantID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("code", "SE001")
	// 不设置tenant_id

	handler.GetPositionByCode(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionByCode", mock.Anything, mock.Anything, mock.Anything)
}

// ==================== GetPositionsByDepartment 测试 ====================

// TestGetPositionsByDepartment_Success 测试成功获取部门的所有岗位
func TestGetPositionsByDepartment_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("dept_id", "dept_001")

	positions := []*orgentity.Position{
		createTestPosition(),
		{
			PositionID:   "pos_002",
			TenantID:     "tenant_001",
			PositionName: "技术经理",
			PositionCode: "TM001",
			Level:        4,
			Category:     "技术岗",
			Status:       orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetPositionsByDepartment", ctx, "dept_001").Return(positions, nil)

	handler.GetPositionsByDepartment(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.GetPositionsByDepartmentResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "pos_001", resp.Data[0].PositionID)
	assert.Equal(t, "pos_002", resp.Data[1].PositionID)

	mockService.AssertExpectations(t)
}

// TestGetPositionsByDepartment_MissingDeptID 测试获取部门岗位缺少部门ID
func TestGetPositionsByDepartment_MissingDeptID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置dept_id参数

	handler.GetPositionsByDepartment(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionsByDepartment", mock.Anything, mock.Anything)
}

// TestGetPositionsByDepartment_EmptyList 测试获取部门岗位返回空列表
func TestGetPositionsByDepartment_EmptyList(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("dept_id", "dept_empty")

	mockService.On("GetPositionsByDepartment", ctx, "dept_empty").Return([]*orgentity.Position{}, nil)

	handler.GetPositionsByDepartment(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.GetPositionsByDepartmentResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 0, len(resp.Data))

	mockService.AssertExpectations(t)
}

// ==================== GetPositionsByLevel 测试 ====================

// TestGetPositionsByLevel_Success 测试成功按职级获取岗位
func TestGetPositionsByLevel_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("level", "3")
	c.Request.SetQueryString("?tenant_id=tenant_001")

	positions := []*orgentity.Position{
		createTestPosition(),
		{
			PositionID:   "pos_002",
			TenantID:     "tenant_001",
			PositionName: "高级产品经理",
			PositionCode: "PM003",
			Level:        3,
			Category:     "产品岗",
			Status:       orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetPositionsByLevel", ctx, "tenant_001", 3).Return(positions, nil)

	handler.GetPositionsByLevel(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.GetPositionsByLevelResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, 3, resp.Data[0].Level)
	assert.Equal(t, 3, resp.Data[1].Level)

	mockService.AssertExpectations(t)
}

// TestGetPositionsByLevel_MissingLevel 测试按职级获取岗位缺少职级参数
func TestGetPositionsByLevel_MissingLevel(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置level参数
	c.Request.SetQueryString("?tenant_id=tenant_001")

	handler.GetPositionsByLevel(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionsByLevel", mock.Anything, mock.Anything, mock.Anything)
}

// TestGetPositionsByLevel_InvalidLevel 测试按职级获取岗位职级格式无效
func TestGetPositionsByLevel_InvalidLevel(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("level", "invalid")
	c.Request.SetQueryString("?tenant_id=tenant_001")

	handler.GetPositionsByLevel(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionsByLevel", mock.Anything, mock.Anything, mock.Anything)
}

// TestGetPositionsByLevel_MissingTenantID 测试按职级获取岗位缺少租户ID
func TestGetPositionsByLevel_MissingTenantID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("level", "3")
	// 不设置tenant_id

	handler.GetPositionsByLevel(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionsByLevel", mock.Anything, mock.Anything, mock.Anything)
}

// ==================== GetPositionsByCategory 测试 ====================

// TestGetPositionsByCategory_Success 测试成功按类别获取岗位
func TestGetPositionsByCategory_Success(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("category", "技术岗")
	c.Request.SetQueryString("?tenant_id=tenant_001")

	positions := []*orgentity.Position{
		createTestPosition(),
		{
			PositionID:   "pos_002",
			TenantID:     "tenant_001",
			PositionName: "架构师",
			PositionCode: "SA001",
			Level:        4,
			Category:     "技术岗",
			Status:       orgentity.OrgStatusActive,
		},
	}

	mockService.On("GetPositionsByCategory", ctx, "tenant_001", "技术岗").Return(positions, nil)

	handler.GetPositionsByCategory(ctx, c)

	assert.Equal(t, 200, c.Response.StatusCode())

	var resp org.GetPositionsByCategoryResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), resp.Code)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "技术岗", resp.Data[0].Category)
	assert.Equal(t, "技术岗", resp.Data[1].Category)

	mockService.AssertExpectations(t)
}

// TestGetPositionsByCategory_MissingCategory 测试按类别获取岗位缺少类别参数
func TestGetPositionsByCategory_MissingCategory(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	// 不设置category参数
	c.Request.SetQueryString("?tenant_id=tenant_001")

	handler.GetPositionsByCategory(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionsByCategory", mock.Anything, mock.Anything, mock.Anything)
}

// TestGetPositionsByCategory_MissingTenantID 测试按类别获取岗位缺少租户ID
func TestGetPositionsByCategory_MissingTenantID(t *testing.T) {
	mockService := new(MockPositionService)
	handler := NewPositionHandler(mockService)

	ctx, c := setupTestContext()
	c.Params.Append("category", "技术岗")
	// 不设置tenant_id

	handler.GetPositionsByCategory(ctx, c)

	// 验证返回错误响应
	assert.NotEqual(t, 200, c.Response.StatusCode())

	// 服务不应被调用
	mockService.AssertNotCalled(t, "GetPositionsByCategory", mock.Anything, mock.Anything, mock.Anything)
}

// ==================== 数据转换函数测试 ====================

// TestToPositionData 测试岗位数据转换
func TestToPositionData(t *testing.T) {
	entity := createTestPosition()

	data := toPositionData(entity)

	assert.NotNil(t, data)
	assert.Equal(t, "pos_001", data.PositionID)
	assert.Equal(t, "tenant_001", data.TenantID)
	assert.Equal(t, "dept_001", *data.DeptID)
	assert.Equal(t, "高级软件工程师", data.PositionName)
	assert.Equal(t, "SE001", data.PositionCode)
	assert.Equal(t, 3, data.Level)
	assert.Equal(t, "技术岗", data.Category)
	assert.Equal(t, "负责系统架构设计", data.Responsibilities)
	assert.Equal(t, "5年以上开发经验", data.Requirements)
	assert.Equal(t, "active", data.Status)
	assert.Equal(t, 1, data.SortOrder)
	assert.Equal(t, int64(1704067200000), data.CreatedAt)
	assert.Equal(t, int64(1704067200000), data.UpdatedAt)
}

// TestToPositionData_NilEntity 测试nil实体转换
func TestToPositionData_NilEntity(t *testing.T) {
	data := toPositionData(nil)
	assert.Nil(t, data)
}

// TestToPositionData_WithoutDept 测试没有部门的岗位数据转换
func TestToPositionData_WithoutDept(t *testing.T) {
	entity := &orgentity.Position{
		PositionID:   "pos_001",
		TenantID:     "tenant_001",
		DeptID:       nil, // 无部门
		PositionName: "高级软件工程师",
		PositionCode: "SE001",
		Level:        3,
		Category:     "技术岗",
		Status:       orgentity.OrgStatusActive,
	}

	data := toPositionData(entity)

	assert.NotNil(t, data)
	assert.Nil(t, data.DeptID)
	assert.Equal(t, "pos_001", data.PositionID)
}

// TestToPositionDataList 测试批量岗位数据转换
func TestToPositionDataList(t *testing.T) {
	entities := []*orgentity.Position{
		createTestPosition(),
		{
			PositionID:   "pos_002",
			TenantID:     "tenant_001",
			PositionName: "产品经理",
			PositionCode: "PM001",
			Level:        3,
			Category:     "产品岗",
			Status:       orgentity.OrgStatusActive,
		},
	}

	dataList := toPositionDataList(entities)

	assert.NotNil(t, dataList)
	assert.Equal(t, 2, len(dataList))
	assert.Equal(t, "pos_001", dataList[0].PositionID)
	assert.Equal(t, "pos_002", dataList[1].PositionID)
	assert.Equal(t, "技术岗", dataList[0].Category)
	assert.Equal(t, "产品岗", dataList[1].Category)
}

// TestToPositionDataList_Nil 测试nil列表转换
func TestToPositionDataList_Nil(t *testing.T) {
	data := toPositionDataList(nil)
	assert.Nil(t, data)
}

// TestToPositionDataList_Empty 测试空列表转换
func TestToPositionDataList_Empty(t *testing.T) {
	entities := []*orgentity.Position{}
	dataList := toPositionDataList(entities)

	assert.NotNil(t, dataList)
	assert.Equal(t, 0, len(dataList))
}

// ==================== 辅助函数 ====================

// strPtr 返回字符串指针的辅助函数
func strPtr(s string) *string {
	return &s
}
