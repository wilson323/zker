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
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// ============================================================
// HR 生命周期管理 API Handler
// ============================================================

// HRLifecycleHandler HR生命周期处理器
// 职责：处理合同管理、调岗管理、离职管理相关的HTTP请求
type HRLifecycleHandler struct {
	hrService *service.HRLifecycleService
	logger    *zap.Logger
}

// NewHRLifecycleHandler 创建HR生命周期Handler实例
func NewHRLifecycleHandler(
	hrService *service.HRLifecycleService,
	logger *zap.Logger,
) *HRLifecycleHandler {
	return &HRLifecycleHandler{
		hrService: hrService,
		logger:    logger,
	}
}

// ============================================================
// 统一响应结构
// ============================================================

// APIResponse 统一API响应
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ============================================================
// 请求和响应类型 - 合同管理
// ============================================================

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	TenantID        string  `json:"tenant_id" binding:"required"`
	EmpID           string  `json:"emp_id" binding:"required"`
	ContractType    string  `json:"contract_type" binding:"required,oneof=labor internship internship_renewal project"`
	ContractNo      string  `json:"contract_no" binding:"required"`
	StartDate       int64   `json:"start_date" binding:"required"`
	EndDate         *int64  `json:"end_date,omitempty"`
	Salary          *string `json:"salary,omitempty"`
	SalaryType      string  `json:"salary_type" binding:"required,oneof=monthly yearly"`
	ProbationDays   int     `json:"probation_days" binding:"min=0"`
	ProbationSalary *string `json:"probation_salary,omitempty"`
	WorkHours       string  `json:"work_hours" binding:"required"`
	WorkPlace       *string `json:"work_place,omitempty"`
	ContractFileURL *string `json:"contract_file_url,omitempty"`
}

// CreateContractResponse 创建合同响应
type CreateContractResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *ContractInfo `json:"data,omitempty"`
}

// ContractInfo 合同信息
type ContractInfo struct {
	ContractID      string  `json:"contract_id"`
	TenantID        string  `json:"tenant_id"`
	EmpID           string  `json:"emp_id"`
	ContractType    string  `json:"contract_type"`
	ContractNo      string  `json:"contract_no"`
	StartDate       int64   `json:"start_date"`
	EndDate         *int64  `json:"end_date,omitempty"`
	Salary          *string `json:"salary,omitempty"`
	SalaryType      string  `json:"salary_type"`
	ProbationDays   int     `json:"probation_days"`
	ProbationSalary *string `json:"probation_salary,omitempty"`
	WorkHours       string  `json:"work_hours"`
	WorkPlace       *string `json:"work_place,omitempty"`
	ContractFileURL *string `json:"contract_file_url,omitempty"`
	Status          string  `json:"status"`
	SignedAt        *int64  `json:"signed_at,omitempty"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
}

// SignContractRequest 签署合同请求
type SignContractRequest struct {
	ContractID string `json:"contract_id" binding:"required"`
}

// SignContractResponse 签署合同响应
type SignContractResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetContractResponse 获取合同响应
type GetContractResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *ContractInfo `json:"data,omitempty"`
}

// GetEmployeeContractsResponse 获取员工合同列表响应
type GetEmployeeContractsResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *ContractListData `json:"data,omitempty"`
}

// ContractListData 合同列表数据
type ContractListData struct {
	Contracts []ContractInfo `json:"contracts"`
	Total     int            `json:"total"`
}

// ============================================================
// 请求和响应类型 - 调岗管理
// ============================================================

// TransferEmployeeRequest 调岗请求
type TransferEmployeeRequest struct {
	TenantID      string  `json:"tenant_id" binding:"required"`
	EmpID         string  `json:"emp_id" binding:"required"`
	NewDeptID     *string `json:"new_dept_id,omitempty"`
	NewPositionID *string `json:"new_position_id,omitempty"`
	NewJobLevel   int     `json:"new_job_level" binding:"required,min=1,max=10"`
	NewJobTitle   string  `json:"new_job_title" binding:"required"`
	TransferType  string  `json:"transfer_type" binding:"required,oneof=promotion demotion lateral transfer"`
	TransferDate  int64   `json:"transfer_date" binding:"required"`
	Reason        string  `json:"reason" binding:"required"`
	ApproverID    *string `json:"approver_id,omitempty"`
}

// TransferEmployeeResponse 调岗响应
type TransferEmployeeResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *TransferInfo `json:"data,omitempty"`
}

// TransferInfo 调岗信息
type TransferInfo struct {
	TransferID    string  `json:"transfer_id"`
	TenantID      string  `json:"tenant_id"`
	EmpID         string  `json:"emp_id"`
	OldDeptID     *string `json:"old_dept_id,omitempty"`
	OldPositionID *string `json:"old_position_id,omitempty"`
	OldJobLevel   int     `json:"old_job_level"`
	OldJobTitle   string  `json:"old_job_title"`
	NewDeptID     *string `json:"new_dept_id,omitempty"`
	NewPositionID *string `json:"new_position_id,omitempty"`
	NewJobLevel   int     `json:"new_job_level"`
	NewJobTitle   string  `json:"new_job_title"`
	TransferType  string  `json:"transfer_type"`
	TransferDate  int64   `json:"transfer_date"`
	Reason        string  `json:"reason"`
	ApproverID    *string `json:"approver_id,omitempty"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

// GetTransferResponse 获取调岗记录响应
type GetTransferResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *TransferInfo `json:"data,omitempty"`
}

// GetEmployeeTransfersResponse 获取员工调岗记录列表响应
type GetEmployeeTransfersResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *TransferListData `json:"data,omitempty"`
}

// TransferListData 调岗记录列表数据
type TransferListData struct {
	Transfers []TransferInfo `json:"transfers"`
	Total     int            `json:"total"`
}

// ============================================================
// 请求和响应类型 - 离职管理
// ============================================================

// ResignEmployeeRequest 离职请求
type ResignEmployeeRequest struct {
	TenantID          string  `json:"tenant_id" binding:"required"`
	EmpID             string  `json:"emp_id" binding:"required"`
	ResignationType   string  `json:"resignation_type" binding:"required,oneof=voluntary involuntary dismissal retirement"`
	ResignationReason string  `json:"resignation_reason" binding:"required"`
	ApplyDate         int64   `json:"apply_date" binding:"required"`
	LastWorkDate      int64   `json:"last_work_date" binding:"required"`
	ResignationDate   int64   `json:"resignation_date" binding:"required"`
	HandoverToID      *string `json:"handover_to_id,omitempty"`
	ApproverID        *string `json:"approver_id,omitempty"`
	RehireEligible    bool    `json:"rehire_eligible"`
	Notes             *string `json:"notes,omitempty"`
}

// ResignEmployeeResponse 离职响应
type ResignEmployeeResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *ResignationInfo `json:"data,omitempty"`
}

// ResignationInfo 离职信息
type ResignationInfo struct {
	ResignationID     string  `json:"resignation_id"`
	TenantID          string  `json:"tenant_id"`
	EmpID             string  `json:"emp_id"`
	ResignationType   string  `json:"resignation_type"`
	ResignationReason string  `json:"resignation_reason"`
	ApplyDate         int64   `json:"apply_date"`
	LastWorkDate      int64   `json:"last_work_date"`
	ResignationDate   int64   `json:"resignation_date"`
	HandoverToID      *string `json:"handover_to_id,omitempty"`
	HandoverStatus    string  `json:"handover_status"`
	ApprovalStatus    string  `json:"approval_status"`
	ApproverID        *string `json:"approver_id,omitempty"`
	ApprovalComment   *string `json:"approval_comment,omitempty"`
	ApprovedAt        *int64  `json:"approved_at,omitempty"`
	RehireEligible    bool    `json:"rehire_eligible"`
	Notes             *string `json:"notes,omitempty"`
	CreatedAt         int64   `json:"created_at"`
	UpdatedAt         int64   `json:"updated_at"`
}

// ApproveResignationRequest 审批离职请求
type ApproveResignationRequest struct {
	ApprovalStatus  string  `json:"approval_status" binding:"required,oneof=approved rejected"`
	ApprovalComment *string `json:"approval_comment,omitempty"`
}

// ApproveResignationResponse 审批离职响应
type ApproveResignationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetResignationResponse 获取离职记录响应
type GetResignationResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *ResignationInfo `json:"data,omitempty"`
}

// UpdateHandoverStatusRequest 更新交接状态请求
type UpdateHandoverStatusRequest struct {
	HandoverStatus string `json:"handover_status" binding:"required,oneof=pending in_progress completed"`
}

// UpdateHandoverStatusResponse 更新交接状态响应
type UpdateHandoverStatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetPendingResignationsResponse 获取待审批离职列表响应
type GetPendingResignationsResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    *ResignationListData `json:"data,omitempty"`
}

// ResignationListData 离职记录列表数据
type ResignationListData struct {
	Resignations []ResignationInfo `json:"resignations"`
	Total        int               `json:"total"`
}

// GetUpcomingProbationEndingsResponse 获取即将结束试用期员工列表响应
type GetUpcomingProbationEndingsResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    *EmployeeProbationListData `json:"data,omitempty"`
}

// EmployeeProbationInfo 员工试用期信息
type EmployeeProbationInfo struct {
	EmpID         string  `json:"emp_id"`
	EmpName       string  `json:"emp_name"`
	EmpCode       string  `json:"emp_code"`
	DeptID        *string `json:"dept_id,omitempty"`
	DeptName      *string `json:"dept_name,omitempty"`
	PositionID    *string `json:"position_id,omitempty"`
	PositionName  *string `json:"position_name,omitempty"`
	HireDate      int64   `json:"hire_date"`
	RegularDate   *int64  `json:"regular_date,omitempty"`
	ProbationDays int     `json:"probation_days"`
	DaysRemaining int     `json:"days_remaining"`
}

// EmployeeProbationListData 员工试用期列表数据
type EmployeeProbationListData struct {
	Employees []EmployeeProbationInfo `json:"employees"`
	Total     int                     `json:"total"`
}

// ============================================================
// 合同管理 API
// ============================================================

// CreateContract 创建合同
// POST /api/v1/hr/contracts
func (h *HRLifecycleHandler) CreateContract(ctx context.Context, c *app.RequestContext) {
	var req CreateContractRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 转换为Service层请求
	serviceReq := &service.CreateContractRequest{
		TenantID:        req.TenantID,
		EmpID:           req.EmpID,
		ContractType:    req.ContractType,
		ContractNo:      req.ContractNo,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Salary:          req.Salary,
		SalaryType:      req.SalaryType,
		ProbationDays:   req.ProbationDays,
		ProbationSalary: req.ProbationSalary,
		WorkHours:       req.WorkHours,
		WorkPlace:       req.WorkPlace,
		ContractFileURL: req.ContractFileURL,
	}

	// 调用Service层创建合同
	contract, err := h.hrService.CreateContract(ctx, serviceReq)
	if err != nil {
		h.logger.Error("failed to create contract",
			zap.String("tenant_id", req.TenantID),
			zap.String("emp_id", req.EmpID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, CreateContractResponse{
		Code:    0,
		Message: "success",
		Data: &ContractInfo{
			ContractID:      contract.ContractID,
			TenantID:        contract.TenantID,
			EmpID:           contract.EmpID,
			ContractType:    contract.ContractType,
			ContractNo:      contract.ContractNo,
			StartDate:       contract.StartDate,
			EndDate:         contract.EndDate,
			Salary:          contract.Salary,
			SalaryType:      contract.SalaryType,
			ProbationDays:   contract.ProbationDays,
			ProbationSalary: contract.ProbationSalary,
			WorkHours:       contract.WorkHours,
			WorkPlace:       contract.WorkPlace,
			ContractFileURL: contract.ContractFileURL,
			Status:          contract.Status,
			SignedAt:        contract.SignedAt,
			CreatedAt:       contract.CreatedAt,
			UpdatedAt:       contract.UpdatedAt,
		},
	})
}

// SignContract 签署合同
// POST /api/v1/hr/contracts/:id/sign
func (h *HRLifecycleHandler) SignContract(ctx context.Context, c *app.RequestContext) {
	contractID := c.Param("id")
	if contractID == "" {
		httputil.BadRequest(c, "contract_id is required")
		return
	}

	var req SignContractRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	req.ContractID = contractID

	// 转换为Service层请求
	serviceReq := &service.SignContractRequest{
		ContractID: req.ContractID,
	}

	// 调用Service层签署合同
	if err := h.hrService.SignContract(ctx, serviceReq); err != nil {
		h.logger.Error("failed to sign contract",
			zap.String("contract_id", contractID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, SignContractResponse{
		Code:    0,
		Message: "success",
	})
}

// GetContract 获取合同详情
// GET /api/v1/hr/contracts/:id
func (h *HRLifecycleHandler) GetContract(ctx context.Context, c *app.RequestContext) {
	// TODO: 需要在HRLifecycleService中添加GetContract方法
	// 暂时返回未实现
	c.JSON(http.StatusNotImplemented, APIResponse{
		Code:    http.StatusNotImplemented,
		Message: "not implemented yet",
	})
}

// GetEmployeeContracts 获取员工的所有合同
// GET /api/v1/hr/employees/:emp_id/contracts
func (h *HRLifecycleHandler) GetEmployeeContracts(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("emp_id")
	if empID == "" {
		httputil.BadRequest(c, "emp_id is required")
		return
	}

	// 调用Service层获取员工合同
	contracts, err := h.hrService.GetEmployeeContracts(ctx, empID)
	if err != nil {
		h.logger.Error("failed to get employee contracts",
			zap.String("emp_id", empID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	contractList := make([]ContractInfo, len(contracts))
	for i, contract := range contracts {
		contractList[i] = ContractInfo{
			ContractID:      contract.ContractID,
			TenantID:        contract.TenantID,
			EmpID:           contract.EmpID,
			ContractType:    contract.ContractType,
			ContractNo:      contract.ContractNo,
			StartDate:       contract.StartDate,
			EndDate:         contract.EndDate,
			Salary:          contract.Salary,
			SalaryType:      contract.SalaryType,
			ProbationDays:   contract.ProbationDays,
			ProbationSalary: contract.ProbationSalary,
			WorkHours:       contract.WorkHours,
			WorkPlace:       contract.WorkPlace,
			ContractFileURL: contract.ContractFileURL,
			Status:          contract.Status,
			SignedAt:        contract.SignedAt,
			CreatedAt:       contract.CreatedAt,
			UpdatedAt:       contract.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, GetEmployeeContractsResponse{
		Code:    0,
		Message: "success",
		Data: &ContractListData{
			Contracts: contractList,
			Total:     len(contractList),
		},
	})
}

// GetActiveContract 获取员工的生效合同
// GET /api/v1/hr/employees/:emp_id/contracts/active
func (h *HRLifecycleHandler) GetActiveContract(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("emp_id")
	if empID == "" {
		httputil.BadRequest(c, "emp_id is required")
		return
	}

	// 调用Service层获取生效合同
	contract, err := h.hrService.GetActiveContract(ctx, empID)
	if err != nil {
		h.logger.Error("failed to get active contract",
			zap.String("emp_id", empID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	if contract == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    http.StatusNotFound,
			Message: "active contract not found",
		})
		return
	}

	c.JSON(http.StatusOK, GetContractResponse{
		Code:    0,
		Message: "success",
		Data: &ContractInfo{
			ContractID:      contract.ContractID,
			TenantID:        contract.TenantID,
			EmpID:           contract.EmpID,
			ContractType:    contract.ContractType,
			ContractNo:      contract.ContractNo,
			StartDate:       contract.StartDate,
			EndDate:         contract.EndDate,
			Salary:          contract.Salary,
			SalaryType:      contract.SalaryType,
			ProbationDays:   contract.ProbationDays,
			ProbationSalary: contract.ProbationSalary,
			WorkHours:       contract.WorkHours,
			WorkPlace:       contract.WorkPlace,
			ContractFileURL: contract.ContractFileURL,
			Status:          contract.Status,
			SignedAt:        contract.SignedAt,
			CreatedAt:       contract.CreatedAt,
			UpdatedAt:       contract.UpdatedAt,
		},
	})
}

// ============================================================
// 调岗管理 API
// ============================================================

// TransferEmployee 调岗
// POST /api/v1/hr/employees/:emp_id/transfer
func (h *HRLifecycleHandler) TransferEmployee(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("emp_id")
	if empID == "" {
		httputil.BadRequest(c, "emp_id is required")
		return
	}

	var req TransferEmployeeRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	req.EmpID = empID

	// 转换为Service层请求
	serviceReq := &service.TransferEmployeeRequest{
		TenantID:      req.TenantID,
		EmpID:         req.EmpID,
		NewDeptID:     req.NewDeptID,
		NewPositionID: req.NewPositionID,
		NewJobLevel:   req.NewJobLevel,
		NewJobTitle:   req.NewJobTitle,
		TransferType:  req.TransferType,
		TransferDate:  req.TransferDate,
		Reason:        req.Reason,
		ApproverID:    req.ApproverID,
	}

	// 调用Service层调岗
	transfer, err := h.hrService.TransferEmployee(ctx, serviceReq)
	if err != nil {
		h.logger.Error("failed to transfer employee",
			zap.String("tenant_id", req.TenantID),
			zap.String("emp_id", req.EmpID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, TransferEmployeeResponse{
		Code:    0,
		Message: "success",
		Data: &TransferInfo{
			TransferID:    transfer.TransferID,
			TenantID:      transfer.TenantID,
			EmpID:         transfer.EmpID,
			OldDeptID:     transfer.OldDeptID,
			OldPositionID: transfer.OldPositionID,
			OldJobLevel:   transfer.OldJobLevel,
			OldJobTitle:   transfer.OldJobTitle,
			NewDeptID:     transfer.NewDeptID,
			NewPositionID: transfer.NewPositionID,
			NewJobLevel:   transfer.NewJobLevel,
			NewJobTitle:   transfer.NewJobTitle,
			TransferType:  transfer.TransferType,
			TransferDate:  transfer.TransferDate,
			Reason:        transfer.Reason,
			ApproverID:    transfer.ApproverID,
			CreatedAt:     transfer.CreatedAt,
			UpdatedAt:     transfer.UpdatedAt,
		},
	})
}

// GetEmployeeTransfers 获取员工的调岗记录
// GET /api/v1/hr/employees/:emp_id/transfers
func (h *HRLifecycleHandler) GetEmployeeTransfers(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("emp_id")
	if empID == "" {
		httputil.BadRequest(c, "emp_id is required")
		return
	}

	// 调用Service层获取调岗记录
	transfers, err := h.hrService.GetEmployeeTransfers(ctx, empID)
	if err != nil {
		h.logger.Error("failed to get employee transfers",
			zap.String("emp_id", empID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	transferList := make([]TransferInfo, len(transfers))
	for i, transfer := range transfers {
		transferList[i] = TransferInfo{
			TransferID:    transfer.TransferID,
			TenantID:      transfer.TenantID,
			EmpID:         transfer.EmpID,
			OldDeptID:     transfer.OldDeptID,
			OldPositionID: transfer.OldPositionID,
			OldJobLevel:   transfer.OldJobLevel,
			OldJobTitle:   transfer.OldJobTitle,
			NewDeptID:     transfer.NewDeptID,
			NewPositionID: transfer.NewPositionID,
			NewJobLevel:   transfer.NewJobLevel,
			NewJobTitle:   transfer.NewJobTitle,
			TransferType:  transfer.TransferType,
			TransferDate:  transfer.TransferDate,
			Reason:        transfer.Reason,
			ApproverID:    transfer.ApproverID,
			CreatedAt:     transfer.CreatedAt,
			UpdatedAt:     transfer.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, GetEmployeeTransfersResponse{
		Code:    0,
		Message: "success",
		Data: &TransferListData{
			Transfers: transferList,
			Total:     len(transferList),
		},
	})
}

// GetTransfer 获取调岗记录详情
// GET /api/v1/hr/transfers/:id
func (h *HRLifecycleHandler) GetTransfer(ctx context.Context, c *app.RequestContext) {
	// TODO: 需要在HRLifecycleService中添加GetTransfer方法
	// 暂时返回未实现
	c.JSON(http.StatusNotImplemented, APIResponse{
		Code:    http.StatusNotImplemented,
		Message: "not implemented yet",
	})
}

// ============================================================
// 离职管理 API
// ============================================================

// ResignEmployee 员工离职
// POST /api/v1/hr/employees/:emp_id/resign
func (h *HRLifecycleHandler) ResignEmployee(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("emp_id")
	if empID == "" {
		httputil.BadRequest(c, "emp_id is required")
		return
	}

	var req ResignEmployeeRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	req.EmpID = empID

	// 转换为Service层请求
	serviceReq := &service.ResignEmployeeRequest{
		TenantID:          req.TenantID,
		EmpID:             req.EmpID,
		ResignationType:   req.ResignationType,
		ResignationReason: req.ResignationReason,
		ApplyDate:         req.ApplyDate,
		LastWorkDate:      req.LastWorkDate,
		ResignationDate:   req.ResignationDate,
		HandoverToID:      req.HandoverToID,
		ApproverID:        req.ApproverID,
		RehireEligible:    req.RehireEligible,
		Notes:             req.Notes,
	}

	// 调用Service层处理离职
	resignation, err := h.hrService.ResignEmployee(ctx, serviceReq)
	if err != nil {
		h.logger.Error("failed to resign employee",
			zap.String("tenant_id", req.TenantID),
			zap.String("emp_id", req.EmpID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, ResignEmployeeResponse{
		Code:    0,
		Message: "success",
		Data: &ResignationInfo{
			ResignationID:     resignation.ResignationID,
			TenantID:          resignation.TenantID,
			EmpID:             resignation.EmpID,
			ResignationType:   resignation.ResignationType,
			ResignationReason: resignation.ResignationReason,
			ApplyDate:         resignation.ApplyDate,
			LastWorkDate:      resignation.LastWorkDate,
			ResignationDate:   resignation.ResignationDate,
			HandoverToID:      resignation.HandoverToID,
			HandoverStatus:    resignation.HandoverStatus,
			ApprovalStatus:    resignation.ApprovalStatus,
			ApproverID:        resignation.ApproverID,
			ApprovalComment:   resignation.ApprovalComment,
			ApprovedAt:        resignation.ApprovedAt,
			RehireEligible:    resignation.RehireEligible,
			Notes:             resignation.Notes,
			CreatedAt:         resignation.CreatedAt,
			UpdatedAt:         resignation.UpdatedAt,
		},
	})
}

// ApproveResignation 审批离职
// POST /api/v1/hr/resignations/:id/approve
func (h *HRLifecycleHandler) ApproveResignation(ctx context.Context, c *app.RequestContext) {
	resignationID := c.Param("id")
	if resignationID == "" {
		httputil.BadRequest(c, "resignation_id is required")
		return
	}

	var req ApproveResignationRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 转换为Service层请求
	serviceReq := &service.ApproveResignationRequest{
		ResignationID:   resignationID,
		ApprovalStatus:  req.ApprovalStatus,
		ApprovalComment: req.ApprovalComment,
	}

	// 调用Service层审批离职
	if err := h.hrService.ApproveResignation(ctx, serviceReq); err != nil {
		h.logger.Error("failed to approve resignation",
			zap.String("resignation_id", resignationID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, ApproveResignationResponse{
		Code:    0,
		Message: "success",
	})
}

// GetResignation 获取离职记录详情
// GET /api/v1/hr/resignations/:id
func (h *HRLifecycleHandler) GetResignation(ctx context.Context, c *app.RequestContext) {
	resignationID := c.Param("id")
	if resignationID == "" {
		httputil.BadRequest(c, "resignation_id is required")
		return
	}

	// 调用Service层获取离职记录
	resignation, err := h.hrService.GetResignation(ctx, resignationID)
	if err != nil {
		h.logger.Error("failed to get resignation",
			zap.String("resignation_id", resignationID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	if resignation == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    http.StatusNotFound,
			Message: "resignation not found",
		})
		return
	}

	c.JSON(http.StatusOK, GetResignationResponse{
		Code:    0,
		Message: "success",
		Data: &ResignationInfo{
			ResignationID:     resignation.ResignationID,
			TenantID:          resignation.TenantID,
			EmpID:             resignation.EmpID,
			ResignationType:   resignation.ResignationType,
			ResignationReason: resignation.ResignationReason,
			ApplyDate:         resignation.ApplyDate,
			LastWorkDate:      resignation.LastWorkDate,
			ResignationDate:   resignation.ResignationDate,
			HandoverToID:      resignation.HandoverToID,
			HandoverStatus:    resignation.HandoverStatus,
			ApprovalStatus:    resignation.ApprovalStatus,
			ApproverID:        resignation.ApproverID,
			ApprovalComment:   resignation.ApprovalComment,
			ApprovedAt:        resignation.ApprovedAt,
			RehireEligible:    resignation.RehireEligible,
			Notes:             resignation.Notes,
			CreatedAt:         resignation.CreatedAt,
			UpdatedAt:         resignation.UpdatedAt,
		},
	})
}

// GetEmployeeResignation 获取员工的离职记录
// GET /api/v1/hr/employees/:emp_id/resignation
func (h *HRLifecycleHandler) GetEmployeeResignation(ctx context.Context, c *app.RequestContext) {
	empID := c.Param("emp_id")
	if empID == "" {
		httputil.BadRequest(c, "emp_id is required")
		return
	}

	// 调用Service层获取员工离职记录
	resignation, err := h.hrService.GetEmployeeResignation(ctx, empID)
	if err != nil {
		h.logger.Error("failed to get employee resignation",
			zap.String("emp_id", empID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	if resignation == nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Code:    http.StatusNotFound,
			Message: "resignation not found",
		})
		return
	}

	c.JSON(http.StatusOK, GetResignationResponse{
		Code:    0,
		Message: "success",
		Data: &ResignationInfo{
			ResignationID:     resignation.ResignationID,
			TenantID:          resignation.TenantID,
			EmpID:             resignation.EmpID,
			ResignationType:   resignation.ResignationType,
			ResignationReason: resignation.ResignationReason,
			ApplyDate:         resignation.ApplyDate,
			LastWorkDate:      resignation.LastWorkDate,
			ResignationDate:   resignation.ResignationDate,
			HandoverToID:      resignation.HandoverToID,
			HandoverStatus:    resignation.HandoverStatus,
			ApprovalStatus:    resignation.ApprovalStatus,
			ApproverID:        resignation.ApproverID,
			ApprovalComment:   resignation.ApprovalComment,
			ApprovedAt:        resignation.ApprovedAt,
			RehireEligible:    resignation.RehireEligible,
			Notes:             resignation.Notes,
			CreatedAt:         resignation.CreatedAt,
			UpdatedAt:         resignation.UpdatedAt,
		},
	})
}

// UpdateHandoverStatus 更新交接状态
// PUT /api/v1/hr/resignations/:id/handover
func (h *HRLifecycleHandler) UpdateHandoverStatus(ctx context.Context, c *app.RequestContext) {
	resignationID := c.Param("id")
	if resignationID == "" {
		httputil.BadRequest(c, "resignation_id is required")
		return
	}

	var req UpdateHandoverStatusRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层更新交接状态
	if err := h.hrService.UpdateHandoverStatus(ctx, resignationID, req.HandoverStatus); err != nil {
		h.logger.Error("failed to update handover status",
			zap.String("resignation_id", resignationID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, UpdateHandoverStatusResponse{
		Code:    0,
		Message: "success",
	})
}

// GetPendingResignations 获取待审批的离职列表
// GET /api/v1/hr/resignations/pending
func (h *HRLifecycleHandler) GetPendingResignations(ctx context.Context, c *app.RequestContext) {
	// 获取tenant_id
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.BadRequest(c, "tenant_id is required")
		return
	}

	// 调用Service层获取待审批离职列表
	resignations, err := h.hrService.GetPendingResignations(ctx, tenantID)
	if err != nil {
		h.logger.Error("failed to get pending resignations",
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	resignationList := make([]ResignationInfo, len(resignations))
	for i, resignation := range resignations {
		resignationList[i] = ResignationInfo{
			ResignationID:     resignation.ResignationID,
			TenantID:          resignation.TenantID,
			EmpID:             resignation.EmpID,
			ResignationType:   resignation.ResignationType,
			ResignationReason: resignation.ResignationReason,
			ApplyDate:         resignation.ApplyDate,
			LastWorkDate:      resignation.LastWorkDate,
			ResignationDate:   resignation.ResignationDate,
			HandoverToID:      resignation.HandoverToID,
			HandoverStatus:    resignation.HandoverStatus,
			ApprovalStatus:    resignation.ApprovalStatus,
			ApproverID:        resignation.ApproverID,
			ApprovalComment:   resignation.ApprovalComment,
			ApprovedAt:        resignation.ApprovedAt,
			RehireEligible:    resignation.RehireEligible,
			Notes:             resignation.Notes,
			CreatedAt:         resignation.CreatedAt,
			UpdatedAt:         resignation.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, GetPendingResignationsResponse{
		Code:    0,
		Message: "success",
		Data: &ResignationListData{
			Resignations: resignationList,
			Total:        len(resignationList),
		},
	})
}

// GetUpcomingProbationEndings 获取即将结束试用期的员工列表
// GET /api/v1/hr/employees/probation/upcoming
func (h *HRLifecycleHandler) GetUpcomingProbationEndings(ctx context.Context, c *app.RequestContext) {
	// 获取tenant_id和days
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		httputil.BadRequest(c, "tenant_id is required")
		return
	}

	days := 30 // 默认30天
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := c.Query("days"); err == nil && d != "" {
			// TODO: 解析days参数
		}
	}

	// 调用Service层获取即将结束试用期的员工列表
	emps, err := h.hrService.GetUpcomingProbationEndings(ctx, tenantID, days)
	if err != nil {
		h.logger.Error("failed to get upcoming probation endings",
			zap.String("tenant_id", tenantID),
			zap.Int("days", days),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// TODO: 转换响应，需要获取员工详细信息
	// 目前先返回基本信息
	c.JSON(http.StatusOK, GetUpcomingProbationEndingsResponse{
		Code:    0,
		Message: "success",
		Data: &EmployeeProbationListData{
			Employees: []EmployeeProbationInfo{},
			Total:     len(emps),
		},
	})
}
