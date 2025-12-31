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
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// AddressBookService 企业通讯录服务
// 职责：企业通讯录查询、按部门筛选、按姓名搜索、按职位筛选
type AddressBookService struct {
	employeeRepo  repository.EmployeeRepository
	departmentRepo repository.DepartmentRepository
	positionRepo  repository.PositionRepository
}

// NewAddressBookService 创建通讯录服务实例
func NewAddressBookService(
	employeeRepo repository.EmployeeRepository,
	departmentRepo repository.DepartmentRepository,
	positionRepo repository.PositionRepository,
) *AddressBookService {
	return &AddressBookService{
		employeeRepo:  employeeRepo,
		departmentRepo: departmentRepo,
		positionRepo:  positionRepo,
	}
}

// Contact 通讯录联系人
type Contact struct {
	EmpID      string  `json:"emp_id"`
	Name       string  `json:"name"`
	DeptID     *string `json:"dept_id,omitempty"`
	DeptName   string  `json:"dept_name"`
	PositionID *string `json:"position_id,omitempty"`
	Position   string  `json:"position"`
	JobLevel   int     `json:"job_level"`
	JobTitle   string  `json:"job_title"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	EmployeeType entity.EmployeeType `json:"employee_type"`
	EmployeeStatus entity.EmployeeStatus `json:"employee_status"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
}

// GetAddressBookRequest 查询通讯录请求
type GetAddressBookRequest struct {
	TenantID     string  `json:"tenant_id" binding:"required"`
	OrgID        *string `json:"org_id,omitempty"`
	DepartmentID *string `json:"department_id,omitempty"`
	Name         string  `json:"name,omitempty"`
	Position     string  `json:"position,omitempty"`
	EmployeeType entity.EmployeeType `json:"employee_type,omitempty"`
	EmployeeStatus entity.EmployeeStatus `json:"employee_status,omitempty"`
	Page         int     `json:"page" binding:"min=1"`
	PageSize     int     `json:"page_size" binding:"min=1,max=100"`
}

// GetAddressBookResponse 查询通讯录响应
type GetAddressBookResponse struct {
	Contacts []*Contact `json:"contacts"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// GetAddressBook 获取企业通讯录
func (s *AddressBookService) GetAddressBook(
	ctx context.Context,
	req *GetAddressBookRequest,
) (*GetAddressBookResponse, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 2. 验证组织存在（如果指定）
	if req.OrgID != nil && *req.OrgID != "" {
		// TODO: 需要OrganizationRepository来验证
	}

	// 3. 验证部门存在（如果指定）
	if req.DepartmentID != nil && *req.DepartmentID != "" {
		dept, err := s.departmentRepo.GetByID(ctx, *req.DepartmentID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
				errorx.KV("operation", "get department"),
			)
		}
		if dept == nil {
			return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
		}
		if dept.TenantID != req.TenantID {
			return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
		}
	}

	// 4. 构建查询条件
	filter := &repository.EmployeeFilter{
		TenantID:      req.TenantID,
		DeptID:        req.DepartmentID,
		EmployeeType:  req.EmployeeType,
		EmployeeStatus: req.EmployeeStatus,
		Keyword:       req.Name,
		PageToken:     fmt.Sprintf("%d", req.Page),
		PageSize:      req.PageSize,
	}

	// 5. 查询员工
	employees, total, err := s.employeeRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "list employees"),
		)
	}

	// 6. 转换为通讯录格式
	contacts := make([]*Contact, 0, len(employees))
	for _, emp := range employees {
		contact := s.convertToContact(emp)

		// 按职位筛选
		if req.Position != "" && emp.JobTitle != req.Position {
			continue
		}

		contacts = append(contacts, contact)
	}

	return &GetAddressBookResponse{
		Contacts: contacts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// convertToContact 将员工实体转换为通讯录联系人
func (s *AddressBookService) convertToContact(emp *entity.Employee) *Contact {
	contact := &Contact{
		EmpID:      emp.EmpID,
		Name:       emp.EmpName,
		DeptID:     emp.DeptID,
		PositionID: emp.PositionID,
		JobLevel:   emp.JobLevel,
		JobTitle:   emp.JobTitle,
		Email:      emp.Email,
		Phone:      emp.Phone,
		EmployeeType: emp.EmployeeType,
		EmployeeStatus: emp.EmployeeStatus,
		AvatarURL:  emp.AvatarURL,
	}

	// 获取部门名称
	if emp.DeptID != nil && emp.Department != nil {
		contact.DeptName = emp.Department.DeptName
	}

	// 获取职位名称
	if emp.PositionID != nil && emp.Position != nil {
		contact.Position = emp.Position.PositionName
	}

	return contact
}

// SearchContactsRequest 搜索联系人请求
type SearchContactsRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	Keyword  string `json:"keyword" binding:"required,min=1"`
	Limit    int    `json:"limit" binding:"min=1,max=100"`
}

// SearchContacts 搜索联系人
func (s *AddressBookService) SearchContacts(
	ctx context.Context,
	req *SearchContactsRequest,
) ([]*Contact, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}
	if req.Keyword == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 设置默认限制
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 2. 搜索员工
	employees, err := s.employeeRepo.Search(ctx, req.TenantID, req.Keyword, req.Limit)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "search employees"),
		)
	}

	// 3. 转换为通讯录格式
	contacts := make([]*Contact, len(employees))
	for i, emp := range employees {
		contacts[i] = s.convertToContact(emp)
	}

	return contacts, nil
}

// GetContactsByDepartmentRequest 按部门获取联系人请求
type GetContactsByDepartmentRequest struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	DepartmentID string `json:"department_id" binding:"required"`
	Page         int    `json:"page" binding:"min=1"`
	PageSize     int    `json:"page_size" binding:"min=1,max=100"`
}

// GetContactsByDepartment 按部门获取联系人
func (s *AddressBookService) GetContactsByDepartment(
	ctx context.Context,
	req *GetContactsByDepartmentRequest,
) (*GetAddressBookResponse, error) {
	// 1. 验证部门存在
	dept, err := s.departmentRepo.GetByID(ctx, req.DepartmentID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get department"),
		)
	}
	if dept == nil {
		return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
	}
	if dept.TenantID != req.TenantID {
		return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
	}

	// 2. 查询部门员工
	employees, err := s.employeeRepo.GetByDepartmentID(ctx, req.DepartmentID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get department employees"),
		)
	}

	// 3. 转换为通讯录格式
	contacts := make([]*Contact, 0, len(employees))
	for _, emp := range employees {
		contact := s.convertToContact(emp)
		contacts = append(contacts, contact)
	}

	// 4. 分页处理
	total := int64(len(contacts))
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize

	if start >= len(contacts) {
		return &GetAddressBookResponse{
			Contacts: []*Contact{},
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, nil
	}

	if end > len(contacts) {
		end = len(contacts)
	}

	pagedContacts := contacts[start:end]

	return &GetAddressBookResponse{
		Contacts: pagedContacts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetContactsByPositionRequest 按职位获取联系人请求
type GetContactsByPositionRequest struct {
	TenantID   string `json:"tenant_id" binding:"required"`
	PositionID string `json:"position_id" binding:"required"`
	Page       int    `json:"page" binding:"min=1"`
	PageSize   int    `json:"page_size" binding:"min=1,max=100"`
}

// GetContactsByPosition 按职位获取联系人
func (s *AddressBookService) GetContactsByPosition(
	ctx context.Context,
	req *GetContactsByPositionRequest,
) (*GetAddressBookResponse, error) {
	// 1. 验证职位存在
	position, err := s.positionRepo.GetByID(ctx, req.PositionID)
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

	// 2. 查询该职位的员工
	filter := &repository.EmployeeFilter{
		TenantID:   req.TenantID,
		PositionID: &req.PositionID,
		PageToken:  fmt.Sprintf("%d", req.Page),
		PageSize:   req.PageSize,
	}

	employees, total, err := s.employeeRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "list employees by position"),
		)
	}

	// 3. 转换为通讯录格式
	contacts := make([]*Contact, len(employees))
	for i, emp := range employees {
		contacts[i] = s.convertToContact(emp)
	}

	return &GetAddressBookResponse{
		Contacts: contacts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
