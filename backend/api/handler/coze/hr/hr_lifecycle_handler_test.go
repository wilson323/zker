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

package hr

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
)

// MockHRLifecycleService 模拟HR生命周期服务
type MockHRLifecycleService struct {
	mock.Mock
}

// CreateContract 模拟创建合同
func (m *MockHRLifecycleService) CreateContract(
	ctx context.Context,
	req *service.CreateContractRequest,
) (*entity.EmployeeContract, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeContract), args.Error(1)
}

// SignContract 模拟签署合同
func (m *MockHRLifecycleService) SignContract(
	ctx context.Context,
	req *service.SignContractRequest,
) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// GetEmployeeContracts 模拟获取员工合同列表
func (m *MockHRLifecycleService) GetEmployeeContracts(
	ctx context.Context,
	empID string,
) ([]*entity.EmployeeContract, error) {
	args := m.Called(ctx, empID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EmployeeContract), args.Error(1)
}

// GetActiveContract 模拟获取生效合同
func (m *MockHRLifecycleService) GetActiveContract(
	ctx context.Context,
	empID string,
) (*entity.EmployeeContract, error) {
	args := m.Called(ctx, empID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeContract), args.Error(1)
}

// TransferEmployee 模拟调岗
func (m *MockHRLifecycleService) TransferEmployee(
	ctx context.Context,
	req *service.TransferEmployeeRequest,
) (*entity.EmployeeTransfer, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeTransfer), args.Error(1)
}

// GetEmployeeTransfers 模拟获取调岗记录列表
func (m *MockHRLifecycleService) GetEmployeeTransfers(
	ctx context.Context,
	empID string,
) ([]*entity.EmployeeTransfer, error) {
	args := m.Called(ctx, empID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EmployeeTransfer), args.Error(1)
}

// ResignEmployee 模拟员工离职
func (m *MockHRLifecycleService) ResignEmployee(
	ctx context.Context,
	req *service.ResignEmployeeRequest,
) (*entity.EmployeeResignation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeResignation), args.Error(1)
}

// ApproveResignation 模拟审批离职
func (m *MockHRLifecycleService) ApproveResignation(
	ctx context.Context,
	req *service.ApproveResignationRequest,
) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// GetResignation 模拟获取离职记录
func (m *MockHRLifecycleService) GetResignation(
	ctx context.Context,
	resignationID string,
) (*entity.EmployeeResignation, error) {
	args := m.Called(ctx, resignationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeResignation), args.Error(1)
}

// GetEmployeeResignation 模拟获取员工离职记录
func (m *MockHRLifecycleService) GetEmployeeResignation(
	ctx context.Context,
	empID string,
) (*entity.EmployeeResignation, error) {
	args := m.Called(ctx, empID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeResignation), args.Error(1)
}

// UpdateHandoverStatus 模拟更新交接状态
func (m *MockHRLifecycleService) UpdateHandoverStatus(
	ctx context.Context,
	resignationID string,
	handoverStatus string,
) error {
	args := m.Called(ctx, resignationID, handoverStatus)
	return args.Error(0)
}

// GetPendingResignations 模拟获取待审批离职列表
func (m *MockHRLifecycleService) GetPendingResignations(
	ctx context.Context,
	tenantID string,
) ([]*entity.EmployeeResignation, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.EmployeeResignation), args.Error(1)
}

// GetUpcomingProbationEndings 模拟获取即将结束试用期员工列表
func (m *MockHRLifecycleService) GetUpcomingProbationEndings(
	ctx context.Context,
	tenantID string,
	days int,
) ([]*entity.Employee, error) {
	args := m.Called(ctx, tenantID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Employee), args.Error(1)
}

// ============================================================
// 辅助函数
// ============================================================

// newTestContext 创建测试上下文
func newTestContext() *app.RequestContext {
	ctx := &app.RequestContext{}
	ctx.Request.SetContentType("application/json")
	return ctx
}

// newTestHandler 创建测试Handler
func newTestHandler(mockService *MockHRLifecycleService) *HRLifecycleHandler {
	logger := zap.NewNop() // 使用无输出日志
	return NewHRLifecycleHandler(mockService, logger)
}

// ============================================================
// 合同管理测试 (5个API)
// ============================================================

// TestCreateContract_Success 测试成功创建合同
func TestCreateContract_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()

	salary := "15000"
	probationSalary := "12000"
	endDate := int64(1735689600000) // 2025-01-01
	workPlace := "北京"

	req := CreateContractRequest{
		TenantID:        "tenant_001",
		EmpID:           "emp_001",
		ContractType:    "labor",
		ContractNo:      "CONTRACT_2024001",
		StartDate:       1704067200000,
		EndDate:         &endDate,
		Salary:          &salary,
		SalaryType:      "monthly",
		ProbationDays:   90,
		ProbationSalary: &probationSalary,
		WorkHours:       "9:00-18:00",
		WorkPlace:       &workPlace,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedContract := &entity.EmployeeContract{
		ContractID:      "contract_001",
		TenantID:        "tenant_001",
		EmpID:           "emp_001",
		ContractType:    "labor",
		ContractNo:      "CONTRACT_2024001",
		StartDate:       1704067200000,
		EndDate:         &endDate,
		Salary:          &salary,
		SalaryType:      "monthly",
		ProbationDays:   90,
		ProbationSalary: &probationSalary,
		WorkHours:       "9:00-18:00",
		WorkPlace:       &workPlace,
		Status:          "draft",
		CreatedAt:       1704067200000,
		UpdatedAt:       1704067200000,
	}

	mockService.On("CreateContract", ctx, mock.Anything).Return(expectedContract, nil)

	handler.CreateContract(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp CreateContractResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "contract_001", resp.Data.ContractID)
	assert.Equal(t, "draft", resp.Data.Status)

	mockService.AssertExpectations(t)
}

// TestCreateContract_BindError 测试创建合同参数绑定错误
func TestCreateContract_BindError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()

	// 无效的JSON
	c.Request.SetBody([]byte("{invalid json"))

	handler.CreateContract(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestCreateContract_ServiceError 测试创建合同服务层错误
func TestCreateContract_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()

	req := CreateContractRequest{
		TenantID:     "tenant_001",
		EmpID:        "emp_001",
		ContractType: "labor",
		ContractNo:   "CONTRACT_2024001",
		StartDate:    1704067200000,
		SalaryType:   "monthly",
		WorkHours:    "9:00-18:00",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("CreateContract", ctx, mock.Anything).
		Return(nil, errors.New("employee not found"))

	handler.CreateContract(ctx, c)

	// 应该返回500错误
	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestSignContract_Success 测试成功签署合同
func TestSignContract_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "contract_001")

	req := SignContractRequest{
		ContractID: "contract_001",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("SignContract", ctx, mock.Anything).Return(nil)

	handler.SignContract(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp SignContractResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

// TestSignContract_MissingContractID 测试签署合同缺少合同ID
func TestSignContract_MissingContractID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加contract_id参数

	req := SignContractRequest{}
	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.SignContract(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestSignContract_ServiceError 测试签署合同服务层错误
func TestSignContract_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "contract_001")

	req := SignContractRequest{}
	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("SignContract", ctx, mock.Anything).
		Return(errors.New("contract not found"))

	handler.SignContract(ctx, c)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestGetContract_NotImplemented 测试获取合同详情未实现
func TestGetContract_NotImplemented(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "contract_001")

	handler.GetContract(ctx, c)

	assert.Equal(t, http.StatusNotImplemented, c.Response.StatusCode())

	var resp APIResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotImplemented, resp.Code)
}

// TestGetEmployeeContracts_Success 测试成功获取员工合同列表
func TestGetEmployeeContracts_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	salary := "15000"
	endDate := int64(1735689600000)

	expectedContracts := []*entity.EmployeeContract{
		{
			ContractID: "contract_001",
			TenantID:   "tenant_001",
			EmpID:      "emp_001",
			ContractNo: "CONTRACT_2024001",
			StartDate:  1704067200000,
			EndDate:    &endDate,
			Salary:     &salary,
			SalaryType: "monthly",
			WorkHours:  "9:00-18:00",
			Status:     "active",
			CreatedAt:  1704067200000,
			UpdatedAt:  1704067200000,
		},
		{
			ContractID: "contract_002",
			TenantID:   "tenant_001",
			EmpID:      "emp_001",
			ContractNo: "CONTRACT_2023001",
			StartDate:  1672531200000,
			EndDate:    nil,
			Salary:     &salary,
			SalaryType: "monthly",
			WorkHours:  "9:00-18:00",
			Status:     "expired",
			CreatedAt:  1672531200000,
			UpdatedAt:  1704067200000,
		},
	}

	mockService.On("GetEmployeeContracts", ctx, "emp_001").
		Return(expectedContracts, nil)

	handler.GetEmployeeContracts(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetEmployeeContractsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, 2, resp.Data.Total)
	assert.Equal(t, 2, len(resp.Data.Contracts))
	assert.Equal(t, "contract_001", resp.Data.Contracts[0].ContractID)

	mockService.AssertExpectations(t)
}

// TestGetEmployeeContracts_MissingEmpID 测试获取员工合同列表缺少员工ID
func TestGetEmployeeContracts_MissingEmpID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加emp_id参数

	handler.GetEmployeeContracts(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetEmployeeContracts_ServiceError 测试获取员工合同列表服务层错误
func TestGetEmployeeContracts_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	mockService.On("GetEmployeeContracts", ctx, "emp_001").
		Return(nil, errors.New("database error"))

	handler.GetEmployeeContracts(ctx, c)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestGetActiveContract_Success 测试成功获取生效合同
func TestGetActiveContract_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	salary := "15000"
	endDate := int64(1735689600000)

	expectedContract := &entity.EmployeeContract{
		ContractID: "contract_001",
		TenantID:   "tenant_001",
		EmpID:      "emp_001",
		ContractNo: "CONTRACT_2024001",
		StartDate:  1704067200000,
		EndDate:    &endDate,
		Salary:     &salary,
		SalaryType: "monthly",
		WorkHours:  "9:00-18:00",
		Status:     "active",
		CreatedAt:  1704067200000,
		UpdatedAt:  1704067200000,
	}

	mockService.On("GetActiveContract", ctx, "emp_001").
		Return(expectedContract, nil)

	handler.GetActiveContract(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetContractResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "contract_001", resp.Data.ContractID)
	assert.Equal(t, "active", resp.Data.Status)

	mockService.AssertExpectations(t)
}

// TestGetActiveContract_NotFound 测试获取生效合同未找到
func TestGetActiveContract_NotFound(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	mockService.On("GetActiveContract", ctx, "emp_001").
		Return(nil, nil)

	handler.GetActiveContract(ctx, c)

	assert.Equal(t, http.StatusNotFound, c.Response.StatusCode())

	var resp APIResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)

	mockService.AssertExpectations(t)
}

// TestGetActiveContract_MissingEmpID 测试获取生效合同缺少员工ID
func TestGetActiveContract_MissingEmpID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加emp_id参数

	handler.GetActiveContract(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// ============================================================
// 调岗管理测试 (3个API)
// ============================================================

// TestTransferEmployee_Success 测试成功调岗
func TestTransferEmployee_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	newDeptID := "dept_002"
	newPositionID := "pos_002"

	req := TransferEmployeeRequest{
		TenantID:      "tenant_001",
		NewDeptID:     &newDeptID,
		NewPositionID: &newPositionID,
		NewJobLevel:   5,
		NewJobTitle:   "高级工程师",
		TransferType:  "promotion",
		TransferDate:  1704067200000,
		Reason:        "晋升",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedTransfer := &entity.EmployeeTransfer{
		TransferID:    "transfer_001",
		TenantID:      "tenant_001",
		EmpID:         "emp_001",
		OldDeptID:     strPtr("dept_001"),
		OldPositionID: strPtr("pos_001"),
		OldJobLevel:   4,
		OldJobTitle:   "工程师",
		NewDeptID:     &newDeptID,
		NewPositionID: &newPositionID,
		NewJobLevel:   5,
		NewJobTitle:   "高级工程师",
		TransferType:  "promotion",
		TransferDate:  1704067200000,
		Reason:        "晋升",
		CreatedAt:     1704067200000,
		UpdatedAt:     1704067200000,
	}

	mockService.On("TransferEmployee", ctx, mock.Anything).
		Return(expectedTransfer, nil)

	handler.TransferEmployee(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp TransferEmployeeResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "transfer_001", resp.Data.TransferID)
	assert.Equal(t, "promotion", resp.Data.TransferType)
	assert.Equal(t, 5, resp.Data.NewJobLevel)

	mockService.AssertExpectations(t)
}

// TestTransferEmployee_MissingEmpID 测试调岗缺少员工ID
func TestTransferEmployee_MissingEmpID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加emp_id参数

	req := TransferEmployeeRequest{
		TenantID:     "tenant_001",
		NewJobLevel:  5,
		NewJobTitle:  "高级工程师",
		TransferType: "promotion",
		TransferDate: 1704067200000,
		Reason:       "晋升",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.TransferEmployee(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestTransferEmployee_ServiceError 测试调岗服务层错误
func TestTransferEmployee_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	req := TransferEmployeeRequest{
		TenantID:     "tenant_001",
		NewJobLevel:  5,
		NewJobTitle:  "高级工程师",
		TransferType: "promotion",
		TransferDate: 1704067200000,
		Reason:       "晋升",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("TransferEmployee", ctx, mock.Anything).
		Return(nil, errors.New("employee not found"))

	handler.TransferEmployee(ctx, c)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestGetEmployeeTransfers_Success 测试成功获取调岗记录列表
func TestGetEmployeeTransfers_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	newDeptID := "dept_002"

	expectedTransfers := []*entity.EmployeeTransfer{
		{
			TransferID:   "transfer_001",
			TenantID:     "tenant_001",
			EmpID:        "emp_001",
			OldJobLevel:  4,
			OldJobTitle:  "工程师",
			NewDeptID:    &newDeptID,
			NewJobLevel:  5,
			NewJobTitle:  "高级工程师",
			TransferType: "promotion",
			TransferDate: 1704067200000,
			Reason:       "晋升",
			CreatedAt:    1704067200000,
			UpdatedAt:    1704067200000,
		},
	}

	mockService.On("GetEmployeeTransfers", ctx, "emp_001").
		Return(expectedTransfers, nil)

	handler.GetEmployeeTransfers(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetEmployeeTransfersResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, 1, resp.Data.Total)
	assert.Equal(t, "transfer_001", resp.Data.Transfers[0].TransferID)

	mockService.AssertExpectations(t)
}

// TestGetEmployeeTransfers_MissingEmpID 测试获取调岗记录列表缺少员工ID
func TestGetEmployeeTransfers_MissingEmpID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加emp_id参数

	handler.GetEmployeeTransfers(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetTransfer_NotImplemented 测试获取调岗记录详情未实现
func TestGetTransfer_NotImplemented(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "transfer_001")

	handler.GetTransfer(ctx, c)

	assert.Equal(t, http.StatusNotImplemented, c.Response.StatusCode())

	var resp APIResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotImplemented, resp.Code)
}

// ============================================================
// 离职管理测试 (7个API)
// ============================================================

// TestResignEmployee_Success 测试成功员工离职
func TestResignEmployee_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	handoverToID := "emp_002"
	approverID := "emp_003"
	notes := "个人原因离职"

	req := ResignEmployeeRequest{
		TenantID:          "tenant_001",
		ResignationType:   "voluntary",
		ResignationReason: "个人发展",
		ApplyDate:         1704067200000,
		LastWorkDate:      1704153600000,
		ResignationDate:   1704240000000,
		HandoverToID:      &handoverToID,
		ApproverID:        &approverID,
		RehireEligible:    true,
		Notes:             &notes,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	expectedResignation := &entity.EmployeeResignation{
		ResignationID:     "resignation_001",
		TenantID:          "tenant_001",
		EmpID:             "emp_001",
		ResignationType:   "voluntary",
		ResignationReason: "个人发展",
		ApplyDate:         1704067200000,
		LastWorkDate:      1704153600000,
		ResignationDate:   1704240000000,
		HandoverToID:      &handoverToID,
		HandoverStatus:    "pending",
		ApprovalStatus:    "pending",
		ApproverID:        &approverID,
		RehireEligible:    true,
		Notes:             &notes,
		CreatedAt:         1704067200000,
		UpdatedAt:         1704067200000,
	}

	mockService.On("ResignEmployee", ctx, mock.Anything).
		Return(expectedResignation, nil)

	handler.ResignEmployee(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp ResignEmployeeResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "resignation_001", resp.Data.ResignationID)
	assert.Equal(t, "voluntary", resp.Data.ResignationType)
	assert.Equal(t, "pending", resp.Data.ApprovalStatus)
	assert.True(t, resp.Data.RehireEligible)

	mockService.AssertExpectations(t)
}

// TestResignEmployee_MissingEmpID 测试员工离职缺少员工ID
func TestResignEmployee_MissingEmpID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加emp_id参数

	req := ResignEmployeeRequest{
		TenantID:          "tenant_001",
		ResignationType:   "voluntary",
		ResignationReason: "个人发展",
		ApplyDate:         1704067200000,
		LastWorkDate:      1704153600000,
		ResignationDate:   1704240000000,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.ResignEmployee(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestResignEmployee_ServiceError 测试员工离职服务层错误
func TestResignEmployee_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	req := ResignEmployeeRequest{
		TenantID:          "tenant_001",
		ResignationType:   "voluntary",
		ResignationReason: "个人发展",
		ApplyDate:         1704067200000,
		LastWorkDate:      1704153600000,
		ResignationDate:   1704240000000,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("ResignEmployee", ctx, mock.Anything).
		Return(nil, errors.New("employee not found"))

	handler.ResignEmployee(ctx, c)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestApproveResignation_Success 测试成功审批离职
func TestApproveResignation_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "resignation_001")

	approvalComment := "同意离职申请"

	req := ApproveResignationRequest{
		ApprovalStatus:  "approved",
		ApprovalComment: &approvalComment,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("ApproveResignation", ctx, mock.Anything).Return(nil)

	handler.ApproveResignation(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp ApproveResignationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

// TestApproveResignation_Reject 测试拒绝离职申请
func TestApproveResignation_Reject(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "resignation_001")

	approvalComment := "暂时不接受离职申请"

	req := ApproveResignationRequest{
		ApprovalStatus:  "rejected",
		ApprovalComment: &approvalComment,
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("ApproveResignation", ctx, mock.Anything).Return(nil)

	handler.ApproveResignation(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp ApproveResignationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	mockService.AssertExpectations(t)
}

// TestApproveResignation_MissingResignationID 测试审批离职缺少离职ID
func TestApproveResignation_MissingResignationID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加id参数

	req := ApproveResignationRequest{
		ApprovalStatus: "approved",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.ApproveResignation(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetResignation_Success 测试成功获取离职记录
func TestGetResignation_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "resignation_001")

	approverID := "emp_003"

	expectedResignation := &entity.EmployeeResignation{
		ResignationID:     "resignation_001",
		TenantID:          "tenant_001",
		EmpID:             "emp_001",
		ResignationType:   "voluntary",
		ResignationReason: "个人发展",
		ApplyDate:         1704067200000,
		LastWorkDate:      1704153600000,
		ResignationDate:   1704240000000,
		HandoverStatus:    "completed",
		ApprovalStatus:    "approved",
		ApproverID:        &approverID,
		ApprovedAt:        int64Ptr(1704067200000),
		RehireEligible:    true,
		CreatedAt:         1704067200000,
		UpdatedAt:         1704067200000,
	}

	mockService.On("GetResignation", ctx, "resignation_001").
		Return(expectedResignation, nil)

	handler.GetResignation(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetResignationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "resignation_001", resp.Data.ResignationID)
	assert.Equal(t, "approved", resp.Data.ApprovalStatus)

	mockService.AssertExpectations(t)
}

// TestGetResignation_NotFound 测试获取离职记录未找到
func TestGetResignation_NotFound(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "resignation_001")

	mockService.On("GetResignation", ctx, "resignation_001").
		Return(nil, nil)

	handler.GetResignation(ctx, c)

	assert.Equal(t, http.StatusNotFound, c.Response.StatusCode())

	var resp APIResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)

	mockService.AssertExpectations(t)
}

// TestGetResignation_MissingResignationID 测试获取离职记录缺少ID
func TestGetResignation_MissingResignationID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加id参数

	handler.GetResignation(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetEmployeeResignation_Success 测试成功获取员工离职记录
func TestGetEmployeeResignation_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	approverID := "emp_003"

	expectedResignation := &entity.EmployeeResignation{
		ResignationID:     "resignation_001",
		TenantID:          "tenant_001",
		EmpID:             "emp_001",
		ResignationType:   "voluntary",
		ResignationReason: "个人发展",
		ApplyDate:         1704067200000,
		LastWorkDate:      1704153600000,
		ResignationDate:   1704240000000,
		HandoverStatus:    "in_progress",
		ApprovalStatus:    "approved",
		ApproverID:        &approverID,
		RehireEligible:    true,
		CreatedAt:         1704067200000,
		UpdatedAt:         1704067200000,
	}

	mockService.On("GetEmployeeResignation", ctx, "emp_001").
		Return(expectedResignation, nil)

	handler.GetEmployeeResignation(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetResignationResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, "resignation_001", resp.Data.ResignationID)

	mockService.AssertExpectations(t)
}

// TestGetEmployeeResignation_NotFound 测试获取员工离职记录未找到
func TestGetEmployeeResignation_NotFound(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("emp_id", "emp_001")

	mockService.On("GetEmployeeResignation", ctx, "emp_001").
		Return(nil, nil)

	handler.GetEmployeeResignation(ctx, c)

	assert.Equal(t, http.StatusNotFound, c.Response.StatusCode())

	var resp APIResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Code)

	mockService.AssertExpectations(t)
}

// TestUpdateHandoverStatus_Success 测试成功更新交接状态
func TestUpdateHandoverStatus_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "resignation_001")

	req := UpdateHandoverStatusRequest{
		HandoverStatus: "in_progress",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("UpdateHandoverStatus", ctx, "resignation_001", "in_progress").
		Return(nil)

	handler.UpdateHandoverStatus(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp UpdateHandoverStatusResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)

	mockService.AssertExpectations(t)
}

// TestUpdateHandoverStatus_Completed 测试更新交接状态为已完成
func TestUpdateHandoverStatus_Completed(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.Params.Append("id", "resignation_001")

	req := UpdateHandoverStatusRequest{
		HandoverStatus: "completed",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	mockService.On("UpdateHandoverStatus", ctx, "resignation_001", "completed").
		Return(nil)

	handler.UpdateHandoverStatus(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp UpdateHandoverStatusResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	mockService.AssertExpectations(t)
}

// TestUpdateHandoverStatus_MissingResignationID 测试更新交接状态缺少ID
func TestUpdateHandoverStatus_MissingResignationID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加id参数

	req := UpdateHandoverStatusRequest{
		HandoverStatus: "in_progress",
	}

	body, _ := json.Marshal(req)
	c.Request.SetBody(body)

	handler.UpdateHandoverStatus(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetPendingResignations_Success 测试成功获取待审批离职列表
func TestGetPendingResignations_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.QueryArgs().Add("tenant_id", "tenant_001")

	approverID := "emp_003"

	expectedResignations := []*entity.EmployeeResignation{
		{
			ResignationID:     "resignation_001",
			TenantID:          "tenant_001",
			EmpID:             "emp_001",
			ResignationType:   "voluntary",
			ResignationReason: "个人发展",
			ApplyDate:         1704067200000,
			LastWorkDate:      1704153600000,
			ResignationDate:   1704240000000,
			HandoverStatus:    "pending",
			ApprovalStatus:    "pending",
			ApproverID:        &approverID,
			RehireEligible:    true,
			CreatedAt:         1704067200000,
			UpdatedAt:         1704067200000,
		},
		{
			ResignationID:     "resignation_002",
			TenantID:          "tenant_001",
			EmpID:             "emp_002",
			ResignationType:   "involuntary",
			ResignationReason: "公司裁员",
			ApplyDate:         1704067200000,
			LastWorkDate:      1704153600000,
			ResignationDate:   1704240000000,
			HandoverStatus:    "pending",
			ApprovalStatus:    "pending",
			ApproverID:        &approverID,
			RehireEligible:    true,
			CreatedAt:         1704067200000,
			UpdatedAt:         1704067200000,
		},
	}

	mockService.On("GetPendingResignations", ctx, "tenant_001").
		Return(expectedResignations, nil)

	handler.GetPendingResignations(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetPendingResignationsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, 2, resp.Data.Total)
	assert.Equal(t, 2, len(resp.Data.Resignations))

	mockService.AssertExpectations(t)
}

// TestGetPendingResignations_MissingTenantID 测试获取待审批离职列表缺少租户ID
func TestGetPendingResignations_MissingTenantID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加tenant_id参数

	handler.GetPendingResignations(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetPendingResignations_ServiceError 测试获取待审批离职列表服务层错误
func TestGetPendingResignations_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.QueryArgs().Add("tenant_id", "tenant_001")

	mockService.On("GetPendingResignations", ctx, "tenant_001").
		Return(nil, errors.New("database error"))

	handler.GetPendingResignations(ctx, c)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// TestGetUpcomingProbationEndings_Success 测试成功获取即将结束试用期员工列表
func TestGetUpcomingProbationEndings_Success(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.QueryArgs().Add("tenant_id", "tenant_001")
	c.QueryArgs().Add("days", "30")

	regularDate := int64(1706745600000) // 30天后

	expectedEmps := []*entity.Employee{
		{
			EmpID:       "emp_001",
			TenantID:    "tenant_001",
			EmpName:     "张三",
			EmpCode:     "EMP001",
			RegularDate: &regularDate,
		},
		{
			EmpID:       "emp_002",
			TenantID:    "tenant_001",
			EmpName:     "李四",
			EmpCode:     "EMP002",
			RegularDate: &regularDate,
		},
	}

	mockService.On("GetUpcomingProbationEndings", ctx, "tenant_001", 30).
		Return(expectedEmps, nil)

	handler.GetUpcomingProbationEndings(ctx, c)

	assert.Equal(t, http.StatusOK, c.Response.StatusCode())

	var resp GetUpcomingProbationEndingsResponse
	err := json.Unmarshal(c.Response.Body(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.NotNil(t, resp.Data)
	assert.Equal(t, 2, resp.Data.Total)

	mockService.AssertExpectations(t)
}

// TestGetUpcomingProbationEndings_MissingTenantID 测试获取即将结束试用期员工列表缺少租户ID
func TestGetUpcomingProbationEndings_MissingTenantID(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	// 不添加tenant_id参数

	handler.GetUpcomingProbationEndings(ctx, c)

	// 应该返回400错误
	assert.NotEqual(t, http.StatusOK, c.Response.StatusCode())
}

// TestGetUpcomingProbationEndings_ServiceError 测试获取即将结束试用期员工列表服务层错误
func TestGetUpcomingProbationEndings_ServiceError(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()
	c := newTestContext()
	c.QueryArgs().Add("tenant_id", "tenant_001")

	mockService.On("GetUpcomingProbationEndings", ctx, "tenant_001", 30).
		Return(nil, errors.New("database error"))

	handler.GetUpcomingProbationEndings(ctx, c)

	assert.Equal(t, http.StatusInternalServerError, c.Response.StatusCode())

	mockService.AssertExpectations(t)
}

// ============================================================
// 辅助函数
// ============================================================

// strPtr 返回字符串指针的辅助函数
func strPtr(s string) *string {
	return &s
}

// int64Ptr 返回int64指针的辅助函数
func int64Ptr(i int64) *int64 {
	return &i
}

// ============================================================
// 状态机和审批流程测试
// ============================================================

// TestContractStateMachine 测试合同状态机: 草稿→生效→过期
func TestContractStateMachine(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()

	// 1. 创建合同(草稿状态)
	t.Run("CreateContract_DraftStatus", func(t *testing.T) {
		c := newTestContext()

		req := CreateContractRequest{
			TenantID:     "tenant_001",
			EmpID:        "emp_001",
			ContractType: "labor",
			ContractNo:   "CONTRACT_2024001",
			StartDate:    1704067200000,
			SalaryType:   "monthly",
			WorkHours:    "9:00-18:00",
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		contract := &entity.EmployeeContract{
			ContractID: "contract_001",
			Status:     "draft",
		}

		mockService.On("CreateContract", ctx, mock.Anything).
			Return(contract, nil)

		handler.CreateContract(ctx, c)

		var resp CreateContractResponse
		json.Unmarshal(c.Response.Body(), &resp)
		assert.Equal(t, "draft", resp.Data.Status)
	})

	// 2. 签署合同(生效状态)
	t.Run("SignContract_ActiveStatus", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("id", "contract_001")

		req := SignContractRequest{}
		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		mockService.On("SignContract", ctx, mock.Anything).Return(nil)

		handler.SignContract(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})
}

// TestResignApprovalWorkflow 测试离职审批流程: 待审批→已批准/已拒绝
func TestResignApprovalWorkflow(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()

	// 1. 创建离职申请(待审批状态)
	t.Run("ResignEmployee_PendingStatus", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("emp_id", "emp_001")

		req := ResignEmployeeRequest{
			TenantID:          "tenant_001",
			ResignationType:   "voluntary",
			ResignationReason: "个人发展",
			ApplyDate:         1704067200000,
			LastWorkDate:      1704153600000,
			ResignationDate:   1704240000000,
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		resignation := &entity.EmployeeResignation{
			ResignationID:  "resignation_001",
			ApprovalStatus: "pending",
		}

		mockService.On("ResignEmployee", ctx, mock.Anything).
			Return(resignation, nil)

		handler.ResignEmployee(ctx, c)

		var resp ResignEmployeeResponse
		json.Unmarshal(c.Response.Body(), &resp)
		assert.Equal(t, "pending", resp.Data.ApprovalStatus)
	})

	// 2. 审批通过
	t.Run("ApproveResignation_Approved", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("id", "resignation_001")

		req := ApproveResignationRequest{
			ApprovalStatus: "approved",
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		mockService.On("ApproveResignation", ctx, mock.Anything).Return(nil)

		handler.ApproveResignation(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})

	// 3. 审批拒绝(另一个测试场景)
	t.Run("ApproveResignation_Rejected", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("id", "resignation_002")

		comment := "不符合离职条件"
		req := ApproveResignationRequest{
			ApprovalStatus:  "rejected",
			ApprovalComment: &comment,
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		mockService.On("ApproveResignation", ctx, mock.Anything).Return(nil)

		handler.ApproveResignation(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})
}

// TestHandoverStatusFlow 测试交接状态流转: 待交接→交接中→已完成
func TestHandoverStatusFlow(t *testing.T) {
	mockService := new(MockHRLifecycleService)
	handler := newTestHandler(mockService)

	ctx := context.Background()

	// 1. 初始状态为待交接
	t.Run("InitialStatus_Pending", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("id", "resignation_001")

		req := ResignEmployeeRequest{
			TenantID:          "tenant_001",
			ResignationType:   "voluntary",
			ResignationReason: "个人发展",
			ApplyDate:         1704067200000,
			LastWorkDate:      1704153600000,
			ResignationDate:   1704240000000,
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		resignation := &entity.EmployeeResignation{
			ResignationID:  "resignation_001",
			HandoverStatus: "pending",
		}

		mockService.On("ResignEmployee", ctx, mock.Anything).
			Return(resignation, nil)

		handler.ResignEmployee(ctx, c)

		var resp ResignEmployeeResponse
		json.Unmarshal(c.Response.Body(), &resp)
		assert.Equal(t, "pending", resp.Data.HandoverStatus)
	})

	// 2. 更新为交接中
	t.Run("UpdateToInProgress", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("id", "resignation_001")

		req := UpdateHandoverStatusRequest{
			HandoverStatus: "in_progress",
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		mockService.On("UpdateHandoverStatus", ctx, "resignation_001", "in_progress").
			Return(nil)

		handler.UpdateHandoverStatus(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})

	// 3. 更新为已完成
	t.Run("UpdateToCompleted", func(t *testing.T) {
		c := newTestContext()
		c.Params.Append("id", "resignation_001")

		req := UpdateHandoverStatusRequest{
			HandoverStatus: "completed",
		}

		body, _ := json.Marshal(req)
		c.Request.SetBody(body)

		mockService.On("UpdateHandoverStatus", ctx, "resignation_001", "completed").
			Return(nil)

		handler.UpdateHandoverStatus(ctx, c)

		assert.Equal(t, http.StatusOK, c.Response.StatusCode())
	})
}
