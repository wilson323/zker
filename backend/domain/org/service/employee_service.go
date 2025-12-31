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
	"regexp"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// EmployeeService 员工管理服务
// 职责：员工的CRUD操作、员工查询、搜索、业务规则验证
type EmployeeService struct {
	empRepo    repository.EmployeeRepository
	orgRepo    repository.OrganizationRepository
	deptRepo   repository.DepartmentRepository
	posRepo    repository.PositionRepository
}

// NewEmployeeService 创建员工服务实例
func NewEmployeeService(
	empRepo repository.EmployeeRepository,
	orgRepo repository.OrganizationRepository,
	deptRepo repository.DepartmentRepository,
	posRepo repository.PositionRepository,
) *EmployeeService {
	return &EmployeeService{
		empRepo:  empRepo,
		orgRepo:  orgRepo,
		deptRepo: deptRepo,
		posRepo:  posRepo,
	}
}

// CreateEmployeeRequest 创建员工请求
type CreateEmployeeRequest struct {
	TenantID      string  `json:"tenant_id" binding:"required"`
	UserID        *string `json:"user_id,omitempty"`
	OrgID         string  `json:"org_id" binding:"required"`
	DeptID        *string `json:"dept_id,omitempty"`
	PositionID    *string `json:"position_id,omitempty"`

	// 基本信息
	EmpName      string       `json:"emp_name" binding:"required,min=1,max=100"`
	EmpCode      string       `json:"emp_code" binding:"required,min=1,max=50"`
	Gender       *string      `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Phone        *string      `json:"phone,omitempty"`
	Email        *string      `json:"email,omitempty"`

	// 职位信息
	EmployeeType entity.EmployeeType   `json:"employee_type" binding:"required,oneof=full_time part_time intern outsourcing contractor"`
	JobLevel     int                  `json:"job_level" binding:"required,min=1,max=10"`
	JobTitle     string               `json:"job_title" binding:"required,min=1,max=100"`

	// 入职信息
	HireDate      int64  `json:"hire_date" binding:"required"`      // 毫秒时间戳
	ProbationDays int    `json:"probation_days" binding:"min=0,max=180"`

	// 工作信息
	WorkLocation   *string `json:"work_location,omitempty"`
	DirectLeaderID *string `json:"direct_leader_id,omitempty"`

	// 个人信息
	IDCard          *string `json:"id_card,omitempty"`
	Birthday        *int64  `json:"birthday,omitempty"`
	Address         *string `json:"address,omitempty"`
	Education       *string `json:"education,omitempty"`
	GraduateSchool  *string `json:"graduate_school,omitempty"`
	Major           *string `json:"major,omitempty"`

	// 紧急联系人
	EmergencyContact *string `json:"emergency_contact,omitempty"`
	EmergencyPhone   *string `json:"emergency_phone,omitempty"`

	// 其他
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateEmployeeRequest 更新员工请求
type UpdateEmployeeRequest struct {
	EmpID         string                `json:"emp_id" binding:"required"`
	DeptID        *string               `json:"dept_id,omitempty"`
	PositionID    *string               `json:"position_id,omitempty"`

	// 基本信息
	EmpName      string  `json:"emp_name" binding:"omitempty,min=1,max=100"`
	Gender       *string `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`

	// 职位信息
	JobLevel int    `json:"job_level" binding:"omitempty,min=1,max=10"`
	JobTitle string `json:"job_title" binding:"omitempty,min=1,max=100"`

	// 工作信息
	WorkLocation   *string `json:"work_location,omitempty"`
	DirectLeaderID *string `json:"direct_leader_id,omitempty"`

	// 个人信息
	IDCard          *string `json:"id_card,omitempty"`
	Birthday        *int64  `json:"birthday,omitempty"`
	Address         *string `json:"address,omitempty"`
	Education       *string `json:"education,omitempty"`
	GraduateSchool  *string `json:"graduate_school,omitempty"`
	Major           *string `json:"major,omitempty"`

	// 紧急联系人
	EmergencyContact *string `json:"emergency_contact,omitempty"`
	EmergencyPhone   *string `json:"emergency_phone,omitempty"`

	// 其他
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CreateEmployee 创建员工
func (s *EmployeeService) CreateEmployee(ctx context.Context, req *CreateEmployeeRequest) (*entity.Employee, error) {
	// 1. 验证组织存在
	org, err := s.orgRepo.GetByID(ctx, req.OrgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get organization"),
            )
	}
	if org == nil {
		return nil, errorx.NewByErrorCode(errno.ErrOrgNotFound)
	}
	if org.TenantID != req.TenantID {
		return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
	}

	// 2. 如果有部门，验证部门存在
	if req.DeptID != nil && *req.DeptID != "" {
		dept, err := s.deptRepo.GetByID(ctx, *req.DeptID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get department"),
            )
		}
		if dept == nil {
			return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
		}
		if dept.OrgID != req.OrgID {
			return nil, errorx.NewByErrorCode(errno.ErrDeptOrgMismatch)
		}
	}

	// 3. 如果有岗位，验证岗位存在
	if req.PositionID != nil && *req.PositionID != "" {
		position, err := s.posRepo.GetByID(ctx, *req.PositionID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get position"),
            )
		}
		if position == nil {
			return nil, errorx.NewByErrorCode(errno.ErrPositionNotFound)
		}
		if position.TenantID != req.TenantID {
			return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
		}
	}

	// 4. 验证工号唯一性
	exists, err := s.empRepo.ExistsByCode(ctx, req.TenantID, req.EmpCode, "")
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "check emp code"),
            )
	}
	if exists {
		return nil, errorx.NewByErrorCode(errno.ErrEmpCodeAlreadyExists)
	}

	// 5. 验证邮箱唯一性（如果提供）
	if req.Email != nil && *req.Email != "" {
		if !s.isValidEmail(*req.Email) {
			return nil, errorx.NewByErrorCode(errno.ErrInvalidEmail)
		}
		exists, err := s.empRepo.ExistsByEmail(ctx, req.TenantID, *req.Email, "")
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "check email"),
            )
		}
		if exists {
			return nil, errorx.NewByErrorCode(errno.ErrEmailAlreadyExists)
		}
	}

	// 6. 验证手机号格式（如果提供）
	if req.Phone != nil && *req.Phone != "" {
		if !s.isValidPhone(*req.Phone) {
			return nil, errorx.NewByErrorCode(errno.ErrInvalidPhone)
		}
	}

	// 7. 验证身份证号格式（如果提供）
	if req.IDCard != nil && *req.IDCard != "" {
		if !s.isValidIDCard(*req.IDCard) {
			return nil, errorx.NewByErrorCode(errno.ErrInvalidIDCard)
		}
	}

	// 8. 生成员工ID
	empID := generateUUID()

	// 9. 计算转正日期
	var regularDate *int64
	if req.ProbationDays > 0 {
		rd := time.Unix(req.HireDate/1000, 0).AddDate(0, 0, req.ProbationDays).UnixMilli()
		regularDate = &rd
	}

	// 10. 创建员工实体
	emp := &entity.Employee{
		EmpID:         empID,
		TenantID:      req.TenantID,
		UserID:        req.UserID,
		OrgID:         req.OrgID,
		DeptID:        req.DeptID,
		PositionID:    req.PositionID,
		EmpName:       req.EmpName,
		EmpCode:       req.EmpCode,
		EmployeeType:  req.EmployeeType,
		EmployeeStatus: entity.EmpStatusTrial, // 默认试用期
		JobLevel:      req.JobLevel,
		JobTitle:      req.JobTitle,
		HireDate:      req.HireDate,
		RegularDate:   regularDate,
		ProbationDays: req.ProbationDays,
		WorkLocation:  req.WorkLocation,
		DirectLeaderID: req.DirectLeaderID,
		Status:        entity.EmpStatusTrial,
		CreatedAt:     time.Now().UnixMilli(),
		UpdatedAt:     time.Now().UnixMilli(),
	}

	// 设置可选字段
	if req.Gender != nil {
		gender := entity.Gender(*req.Gender)
		emp.Gender = &gender
	}
	emp.Phone = req.Phone
	emp.Email = req.Email
	emp.IDCard = req.IDCard
	emp.Birthday = req.Birthday
	emp.Address = req.Address
	emp.Education = req.Education
	emp.GraduateSchool = req.GraduateSchool
	emp.Major = req.Major
	emp.EmergencyContact = req.EmergencyContact
	emp.EmergencyPhone = req.EmergencyPhone
	emp.AvatarURL = req.AvatarURL
	emp.Description = req.Description

	// 11. 创建员工
	if err := s.empRepo.Create(ctx, emp); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create employee"),
            )
	}

	return emp, nil
}

// GetEmployee 获取员工详情
func (s *EmployeeService) GetEmployee(ctx context.Context, empID string) (*entity.Employee, error) {
	if empID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emp, err := s.empRepo.GetByID(ctx, empID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	return emp, nil
}

// UpdateEmployee 更新员工
func (s *EmployeeService) UpdateEmployee(ctx context.Context, req *UpdateEmployeeRequest) (*entity.Employee, error) {
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

	// 2. 如果有部门，验证部门存在且属于同一组织
	if req.DeptID != nil && *req.DeptID != "" {
		dept, err := s.deptRepo.GetByID(ctx, *req.DeptID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get department"),
            )
		}
		if dept == nil {
			return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
		}
		if dept.OrgID != emp.OrgID {
			return nil, errorx.NewByErrorCode(errno.ErrDeptOrgMismatch)
		}
		emp.DeptID = req.DeptID
	}

	// 3. 如果有岗位，验证岗位存在
	if req.PositionID != nil && *req.PositionID != "" {
		position, err := s.posRepo.GetByID(ctx, *req.PositionID)
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
		emp.PositionID = req.PositionID
	}

	// 4. 更新基本信息
	if req.EmpName != "" {
		emp.EmpName = req.EmpName
	}
	if req.Gender != nil {
		gender := entity.Gender(*req.Gender)
		emp.Gender = &gender
	}
	if req.Phone != nil {
		if *req.Phone != "" && !s.isValidPhone(*req.Phone) {
			return nil, errorx.NewByErrorCode(errno.ErrInvalidPhone)
		}
		emp.Phone = req.Phone
	}
	if req.Email != nil {
		if *req.Email != "" {
			if !s.isValidEmail(*req.Email) {
				return nil, errorx.NewByErrorCode(errno.ErrInvalidEmail)
			}
			// 检查邮箱是否被其他人使用
			exists, err := s.empRepo.ExistsByEmail(ctx, emp.TenantID, *req.Email, emp.EmpID)
			if err != nil {
				return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "check email"),
            )
			}
			if exists {
				return nil, errorx.NewByErrorCode(errno.ErrEmailAlreadyExists)
			}
		}
		emp.Email = req.Email
	}

	// 5. 更新职位信息
	if req.JobLevel != 0 {
		emp.JobLevel = req.JobLevel
	}
	if req.JobTitle != "" {
		emp.JobTitle = req.JobTitle
	}

	// 6. 更新工作信息
	emp.WorkLocation = req.WorkLocation
	emp.DirectLeaderID = req.DirectLeaderID

	// 7. 更新个人信息
	if req.IDCard != nil {
		if *req.IDCard != "" && !s.isValidIDCard(*req.IDCard) {
			return nil, errorx.NewByErrorCode(errno.ErrInvalidIDCard)
		}
		emp.IDCard = req.IDCard
	}
	emp.Birthday = req.Birthday
	emp.Address = req.Address
	emp.Education = req.Education
	emp.GraduateSchool = req.GraduateSchool
	emp.Major = req.Major
	emp.EmergencyContact = req.EmergencyContact
	emp.EmergencyPhone = req.EmergencyPhone

	// 8. 更新其他信息
	emp.AvatarURL = req.AvatarURL
	emp.Description = req.Description
	emp.UpdatedAt = time.Now().UnixMilli()

	// 9. 保存更新
	if err := s.empRepo.Update(ctx, emp); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update employee"),
            )
	}

	return emp, nil
}

// DeleteEmployee 删除员工（软删除）
func (s *EmployeeService) DeleteEmployee(ctx context.Context, empID string) error {
	if empID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 1. 获取员工
	emp, err := s.empRepo.GetByID(ctx, empID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}

	// 2. 在职员工不能删除，需要先离职
	if emp.IsActive() {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "CannotDeleteActiveEmployee"),
                )
	}

	// 3. 软删除员工
	if err := s.empRepo.Delete(ctx, empID); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "delete employee"),
            )
	}

	return nil
}

// UpdateEmployeeStatus 更新员工状态
func (s *EmployeeService) UpdateEmployeeStatus(ctx context.Context, empID string, status entity.EmployeeStatus) error {
	if empID == "" || status == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 1. 获取员工
	emp, err := s.empRepo.GetByID(ctx, empID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	// 2. 状态转换验证
	if !s.isValidStatusTransition(emp.EmployeeStatus, status) {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidStatusTransition"),
                )
	}

	// 3. 更新状态
	if err := s.empRepo.UpdateStatus(ctx, empID, status); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update employee status"),
            )
	}

	return nil
}

// isValidStatusTransition 验证状态转换是否合法
func (s *EmployeeService) isValidStatusTransition(oldStatus, newStatus entity.EmployeeStatus) bool {
	// 定义合法的状态转换
	validTransitions := map[entity.EmployeeStatus][]entity.EmployeeStatus{
		entity.EmpStatusTrial:     {entity.EmpStatusActive, entity.EmpStatusResigned},
		entity.EmpStatusProbation: {entity.EmpStatusActive, entity.EmpStatusResigned},
		entity.EmpStatusActive:    {entity.EmpStatusResigned, entity.EmpStatusSuspended, entity.EmpStatusRetired},
		entity.EmpStatusSuspended: {entity.EmpStatusActive, entity.EmpStatusResigned},
	}

	allowedStatuses, exists := validTransitions[oldStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// ListEmployees 分页查询员工列表
func (s *EmployeeService) ListEmployees(ctx context.Context, filter *repository.EmployeeFilter) ([]*entity.Employee, int64, error) {
	if filter.TenantID == "" {
		return nil, 0, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emps, total, err := s.empRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list employees"),
            )
	}

	return emps, total, nil
}

// GetEmployeeByCode 根据工号获取员工
func (s *EmployeeService) GetEmployeeByCode(ctx context.Context, tenantID, code string) (*entity.Employee, error) {
	if tenantID == "" || code == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emp, err := s.empRepo.GetByCode(ctx, tenantID, code)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee by code"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	return emp, nil
}

// GetEmployeeByUserID 根据用户ID获取员工
func (s *EmployeeService) GetEmployeeByUserID(ctx context.Context, userID string) (*entity.Employee, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emp, err := s.empRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee by user id"),
            )
	}

	return emp, nil
}

// GetEmployeesByDepartment 获取部门的所有员工
func (s *EmployeeService) GetEmployeesByDepartment(ctx context.Context, deptID string) ([]*entity.Employee, error) {
	if deptID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emps, err := s.empRepo.GetByDepartmentID(ctx, deptID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employees by department"),
            )
	}

	return emps, nil
}

// GetEmployeesByOrganization 获取组织的所有员工
func (s *EmployeeService) GetEmployeesByOrganization(ctx context.Context, orgID string) ([]*entity.Employee, error) {
	if orgID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emps, err := s.empRepo.GetByOrgID(ctx, orgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employees by organization"),
            )
	}

	return emps, nil
}

// SearchEmployees 搜索员工（按姓名、工号、手机号、邮箱）
func (s *EmployeeService) SearchEmployees(ctx context.Context, tenantID, keyword string, limit int) ([]*entity.Employee, error) {
	if tenantID == "" || keyword == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	if limit <= 0 || limit > 100 {
		limit = 20 // 默认限制20条
	}

	emps, err := s.empRepo.Search(ctx, tenantID, keyword, limit)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "search employees"),
            )
	}

	return emps, nil
}

// GetEmployeesByPinyin 按拼音首字母查询员工
func (s *EmployeeService) GetEmployeesByPinyin(ctx context.Context, tenantID, pinyin string) ([]*entity.Employee, error) {
	if tenantID == "" || pinyin == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emps, err := s.empRepo.GetByPinyin(ctx, tenantID, pinyin)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employees by pinyin"),
            )
	}

	return emps, nil
}

// isValidEmail 验证邮箱格式
func (s *EmployeeService) isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// isValidPhone 验证手机号格式（中国大陆）
func (s *EmployeeService) isValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// isValidIDCard 验证身份证号格式（中国大陆）
func (s *EmployeeService) isValidIDCard(idCard string) bool {
	// 18位身份证号
	pattern := `^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}
