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
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// HRLifecycleService 人力资源生命周期服务
// 职责：员工入职、转正、调岗、离职、合同管理
type HRLifecycleService struct {
	empRepo          repository.EmployeeRepository
	contractRepo     repository.EmployeeContractRepository
	transferRepo     repository.EmployeeTransferRepository
	resignationRepo  repository.EmployeeResignationRepository
	deptRepo         repository.DepartmentRepository
	posRepo          repository.PositionRepository
	db               *gorm.DB
}

// NewHRLifecycleService 创建HR生命周期服务实例
func NewHRLifecycleService(
	empRepo repository.EmployeeRepository,
	contractRepo repository.EmployeeContractRepository,
	transferRepo repository.EmployeeTransferRepository,
	resignationRepo repository.EmployeeResignationRepository,
	deptRepo repository.DepartmentRepository,
	posRepo repository.PositionRepository,
	db *gorm.DB,
) *HRLifecycleService {
	return &HRLifecycleService{
		empRepo:         empRepo,
		contractRepo:    contractRepo,
		transferRepo:    transferRepo,
		resignationRepo: resignationRepo,
		deptRepo:        deptRepo,
		posRepo:         posRepo,
		db:              db,
	}
}

// CreateContractRequest 创建合同请求
type CreateContractRequest struct {
	TenantID          string  `json:"tenant_id" binding:"required"`
	EmpID             string  `json:"emp_id" binding:"required"`
	ContractType      string  `json:"contract_type" binding:"required"`
	ContractNo        string  `json:"contract_no" binding:"required"`
	StartDate         int64   `json:"start_date" binding:"required"`
	EndDate           *int64  `json:"end_date,omitempty"`
	Salary            *string `json:"salary,omitempty"`
	SalaryType        string  `json:"salary_type" binding:"required,oneof=monthly yearly"`
	ProbationDays     int     `json:"probation_days" binding:"min=0"`
	ProbationSalary   *string `json:"probation_salary,omitempty"`
	WorkHours         string  `json:"work_hours" binding:"required"`
	WorkPlace         *string `json:"work_place,omitempty"`
	ContractFileURL   *string `json:"contract_file_url,omitempty"`
}

// SignContractRequest 签署合同请求
type SignContractRequest struct {
	ContractID string `json:"contract_id" binding:"required"`
}

// TransferEmployeeRequest 调岗请求
type TransferEmployeeRequest struct {
	TenantID       string  `json:"tenant_id" binding:"required"`
	EmpID          string  `json:"emp_id" binding:"required"`
	NewDeptID      *string `json:"new_dept_id,omitempty"`
	NewPositionID  *string `json:"new_position_id,omitempty"`
	NewJobLevel    int     `json:"new_job_level" binding:"required,min=1,max=10"`
	NewJobTitle    string  `json:"new_job_title" binding:"required"`
	TransferType   string  `json:"transfer_type" binding:"required"`
	TransferDate   int64   `json:"transfer_date" binding:"required"`
	Reason         string  `json:"reason" binding:"required"`
	ApproverID     *string `json:"approver_id,omitempty"`
}

// ResignEmployeeRequest 离职请求
type ResignEmployeeRequest struct {
	TenantID          string  `json:"tenant_id" binding:"required"`
	EmpID             string  `json:"emp_id" binding:"required"`
	ResignationType   string  `json:"resignation_type" binding:"required"`
	ResignationReason string  `json:"resignation_reason" binding:"required"`
	ApplyDate         int64   `json:"apply_date" binding:"required"`
	LastWorkDate       int64  `json:"last_work_date" binding:"required"`
	ResignationDate   int64   `json:"resignation_date" binding:"required"`
	HandoverToID      *string `json:"handover_to_id,omitempty"`
	ApproverID        *string `json:"approver_id,omitempty"`
	RehireEligible    bool    `json:"rehire_eligible"`
	Notes             *string `json:"notes,omitempty"`
}

// ApproveResignationRequest 审批离职请求
type ApproveResignationRequest struct {
	ResignationID   string  `json:"resignation_id" binding:"required"`
	ApprovalStatus  string  `json:"approval_status" binding:"required,oneof=approved rejected"`
	ApprovalComment *string `json:"approval_comment,omitempty"`
}

// CreateContract 创建合同
func (s *HRLifecycleService) CreateContract(ctx context.Context, req *CreateContractRequest) (*entity.EmployeeContract, error) {
	// 1. 验证员工存在
	emp, err := s.empRepo.GetByID(ctx, req.EmpID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	// 2. 验证合同编号唯一性
	// TODO: 需要在EmployeeContractRepository中添加ExistsByNo方法

	// 3. 生成合同ID
	contractID := generateUUID()

	// 4. 创建合同实体
	contract := &entity.EmployeeContract{
		ContractID:      contractID,
		TenantID:        emp.TenantID,
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
		Status:          "draft", // 初始状态为草稿
		CreatedAt:       time.Now().UnixMilli(),
		UpdatedAt:       time.Now().UnixMilli(),
	}

	// 5. 创建合同
	if err := s.contractRepo.Create(ctx, contract); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create contract"),
            )
	}

	return contract, nil
}

// SignContract 签署合同
func (s *HRLifecycleService) SignContract(ctx context.Context, req *SignContractRequest) error {
	if req.ContractID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 1. 获取合同
	contract, err := s.contractRepo.GetByID(ctx, req.ContractID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get contract"),
            )
	}
	if contract == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "ContractNotFound"),
                )
	}

	// 2. 只有草稿状态可以签署
	if contract.Status != "draft" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "CannotSignContract"),
                )
	}

	// 3. 更新合同状态
	now := time.Now().UnixMilli()
	contract.Status = "active"
	contract.SignedAt = &now
	contract.UpdatedAt = now

	if err := s.contractRepo.Update(ctx, contract); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update contract"),
            )
	}

	return nil
}

// GetActiveContract 获取员工的生效合同
func (s *HRLifecycleService) GetActiveContract(ctx context.Context, empID string) (*entity.EmployeeContract, error) {
	if empID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	contract, err := s.contractRepo.GetActiveByEmployeeID(ctx, empID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get active contract"),
            )
	}

	return contract, nil
}

// GetEmployeeContracts 获取员工的所有合同
func (s *HRLifecycleService) GetEmployeeContracts(ctx context.Context, empID string) ([]*entity.EmployeeContract, error) {
	if empID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	contracts, err := s.contractRepo.GetByEmployeeID(ctx, empID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee contracts"),
            )
	}

	return contracts, nil
}

// TransferEmployee 调岗
func (s *HRLifecycleService) TransferEmployee(ctx context.Context, req *TransferEmployeeRequest) (*entity.EmployeeTransfer, error) {
	// 1. 获取员工
	emp, err := s.empRepo.GetByID(ctx, req.EmpID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	// 2. 验证新部门存在（如果提供）
	if req.NewDeptID != nil && *req.NewDeptID != "" {
		dept, err := s.deptRepo.GetByID(ctx, *req.NewDeptID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get department"),
            )
		}
		if dept == nil {
			return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
		}
		if dept.TenantID != emp.TenantID {
			return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
		}
	}

	// 3. 验证新岗位存在（如果提供）
	if req.NewPositionID != nil && *req.NewPositionID != "" {
		position, err := s.posRepo.GetByID(ctx, *req.NewPositionID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get position"),
            )
		}
		if position == nil {
			return nil, errorx.NewByErrorCode(errno.ErrPositionNotFound)
		}
		if position.TenantID != emp.TenantID {
			return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
		}
	}

	// 4. 生成调岗记录ID
	transferID := generateUUID()

	// 5. 创建调岗记录
	transfer := &entity.EmployeeTransfer{
		TransferID:    transferID,
		TenantID:      emp.TenantID,
		EmpID:         req.EmpID,
		OldDeptID:     emp.DeptID,
		OldPositionID: emp.PositionID,
		OldJobLevel:   emp.JobLevel,
		OldJobTitle:   emp.JobTitle,
		NewDeptID:     req.NewDeptID,
		NewPositionID: req.NewPositionID,
		NewJobLevel:   req.NewJobLevel,
		NewJobTitle:   req.NewJobTitle,
		TransferType:  req.TransferType,
		TransferDate:  req.TransferDate,
		Reason:        req.Reason,
		ApproverID:    req.ApproverID,
		CreatedAt:     time.Now().UnixMilli(),
		UpdatedAt:     time.Now().UnixMilli(),
	}

	// 6. 使用事务创建调岗记录并更新员工信息
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建调岗记录
		if err := s.transferRepo.Create(ctx, transfer); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create transfer record"),
            )
		}

		// 更新员工信息
		emp.DeptID = req.NewDeptID
		emp.PositionID = req.NewPositionID
		emp.JobLevel = req.NewJobLevel
		emp.JobTitle = req.NewJobTitle
		emp.UpdatedAt = time.Now().UnixMilli()

		if err := s.empRepo.Update(ctx, emp); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update employee"),
            )
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

// GetEmployeeTransfers 获取员工的调岗记录
func (s *HRLifecycleService) GetEmployeeTransfers(ctx context.Context, empID string) ([]*entity.EmployeeTransfer, error) {
	if empID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	transfers, err := s.transferRepo.GetByEmployeeID(ctx, empID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee transfers"),
            )
	}

	return transfers, nil
}

// ResignEmployee 员工离职
func (s *HRLifecycleService) ResignEmployee(ctx context.Context, req *ResignEmployeeRequest) (*entity.EmployeeResignation, error) {
	// 1. 获取员工
	emp, err := s.empRepo.GetByID(ctx, req.EmpID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	// 2. 验证员工状态（只有在职员工可以离职）
	if !emp.IsActive() {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotActive)
	}

	// 3. 验证交接人存在（如果提供）
	if req.HandoverToID != nil && *req.HandoverToID != "" {
		handoverTo, err := s.empRepo.GetByID(ctx, *req.HandoverToID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get handover employee"),
            )
		}
		if handoverTo == nil {
			return nil, errorx.NewByErrorCode(errno.ErrHandoverEmployeeNotFound)
		}
		if handoverTo.EmpID == emp.EmpID {
			return nil, errorx.NewByErrorCode(errno.ErrCannotHandoverToSelf)
		}
	}

	// 4. 生成离职记录ID
	resignationID := generateUUID()

	// 5. 创建离职记录
	resignation := &entity.EmployeeResignation{
		ResignationID:    resignationID,
		TenantID:         emp.TenantID,
		EmpID:            req.EmpID,
		ResignationType:  req.ResignationType,
		ResignationReason: req.ResignationReason,
		ApplyDate:        req.ApplyDate,
		LastWorkDate:     req.LastWorkDate,
		ResignationDate:  req.ResignationDate,
		HandoverToID:     req.HandoverToID,
		HandoverStatus:   "pending",
		ApprovalStatus:   "pending",
		ApproverID:       req.ApproverID,
		RehireEligible:   req.RehireEligible,
		Notes:            req.Notes,
		CreatedAt:        time.Now().UnixMilli(),
		UpdatedAt:        time.Now().UnixMilli(),
	}

	// 6. 创建离职记录
	if err := s.resignationRepo.Create(ctx, resignation); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create resignation"),
            )
	}

	return resignation, nil
}

// ApproveResignation 审批离职
func (s *HRLifecycleService) ApproveResignation(ctx context.Context, req *ApproveResignationRequest) error {
	if req.ResignationID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 1. 获取离职记录
	resignation, err := s.resignationRepo.GetByID(ctx, req.ResignationID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get resignation"),
            )
	}
	if resignation == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "ResignationNotFound"),
                )
	}

	// 2. 验证状态（只有待审批状态可以审批）
	if resignation.ApprovalStatus != "pending" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "ResignationAlreadyProcessed"),
                )
	}

	// 3. 更新审批状态
	now := time.Now().UnixMilli()
	resignation.ApprovalStatus = req.ApprovalStatus
	resignation.ApprovalComment = req.ApprovalComment
	resignation.ApprovedAt = &now
	resignation.UpdatedAt = now

	// 4. 如果审批通过，更新员工状态为已离职
	if req.ApprovalStatus == "approved" {
		return s.db.Transaction(func(tx *gorm.DB) error {
			// 更新离职记录
			if err := s.resignationRepo.Update(ctx, resignation); err != nil {
				return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update resignation"),
            )
			}

			// 更新员工状态
			if err := s.empRepo.UpdateStatus(ctx, resignation.EmpID, entity.EmpStatusResigned); err != nil {
				return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update employee status"),
            )
			}

			return nil
		})
	} else {
		// 审批拒绝，只更新离职记录
		if err := s.resignationRepo.Update(ctx, resignation); err != nil {
			return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update resignation"),
            )
		}
	}

	return nil
}

// GetResignation 获取离职记录
func (s *HRLifecycleService) GetResignation(ctx context.Context, resignationID string) (*entity.EmployeeResignation, error) {
	if resignationID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	resignation, err := s.resignationRepo.GetByID(ctx, resignationID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get resignation"),
            )
	}

	return resignation, nil
}

// GetEmployeeResignation 获取员工的离职记录
func (s *HRLifecycleService) GetEmployeeResignation(ctx context.Context, empID string) (*entity.EmployeeResignation, error) {
	if empID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	resignation, err := s.resignationRepo.GetByEmployeeID(ctx, empID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee resignation"),
            )
	}

	return resignation, nil
}

// UpdateHandoverStatus 更新交接状态
func (s *HRLifecycleService) UpdateHandoverStatus(ctx context.Context, resignationID, handoverStatus string) error {
	if resignationID == "" || handoverStatus == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 1. 获取离职记录
	resignation, err := s.resignationRepo.GetByID(ctx, resignationID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get resignation"),
            )
	}
	if resignation == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "ResignationNotFound"),
                )
	}

	// 2. 更新交接状态
	resignation.HandoverStatus = handoverStatus
	resignation.UpdatedAt = time.Now().UnixMilli()

	if err := s.resignationRepo.Update(ctx, resignation); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update resignation"),
            )
	}

	return nil
}

// GetPendingResignations 获取待审批的离职列表
func (s *HRLifecycleService) GetPendingResignations(ctx context.Context, tenantID string) ([]*entity.EmployeeResignation, error) {
	if tenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	filter := &repository.ResignationFilter{
		TenantID:       tenantID,
		ApprovalStatus: "pending",
	}

	resignations, _, err := s.resignationRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get pending resignations"),
            )
	}

	return resignations, nil
}

// GetUpcomingProbationEndings 获取即将结束试用期的员工列表
func (s *HRLifecycleService) GetUpcomingProbationEndings(ctx context.Context, tenantID string, days int) ([]*entity.Employee, error) {
	if tenantID == "" || days <= 0 {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 计算时间范围
	now := time.Now()
	startDate := now.AddDate(0, 0, days).UnixMilli()
	endDate := now.AddDate(0, 0, days+7).UnixMilli() // 7天内

	// TODO: 需要在EmployeeRepository中添加GetByProbationDateRange方法
	// 这里先用简单实现

	filter := &repository.EmployeeFilter{
		TenantID:      tenantID,
		EmployeeStatus: entity.EmpStatusTrial,
	}

	emps, _, err := s.empRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employees"),
            )
	}

	// 过滤出即将结束试用期的员工
	result := make([]*entity.Employee, 0)
	for _, emp := range emps {
		if emp.RegularDate != nil {
			if *emp.RegularDate >= startDate && *emp.RegularDate <= endDate {
				result = append(result, emp)
			}
		}
	}

	return result, nil
}
