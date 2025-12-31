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

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// DirectoryService 通讯录服务
// 职责：组织架构查询、员工搜索、通讯录视图生成
type DirectoryService struct {
	orgRepo    repository.OrganizationRepository
	deptRepo   repository.DepartmentRepository
	empRepo    repository.EmployeeRepository
	posRepo    repository.PositionRepository
}

// NewDirectoryService 创建通讯录服务实例
func NewDirectoryService(
	orgRepo repository.OrganizationRepository,
	deptRepo repository.DepartmentRepository,
	empRepo repository.EmployeeRepository,
	posRepo repository.PositionRepository,
) *DirectoryService {
	return &DirectoryService{
		orgRepo:  orgRepo,
		deptRepo: deptRepo,
		empRepo:  empRepo,
		posRepo:  posRepo,
	}
}

// OrganizationTreeNode 组织架构树节点
type OrganizationTreeNode struct {
	OrgID       string                  `json:"org_id"`
	OrgName     string                  `json:"org_name"`
	OrgType     entity.OrganizationType `json:"org_type"`
	Level       int                     `json:"level"`
	EmployeeCount int                   `json:"employee_count"`
	Children    []*OrganizationTreeNode `json:"children,omitempty"`
}

// DepartmentTreeNode 部门树节点
type DepartmentTreeNode struct {
	DeptID        string                `json:"dept_id"`
	DeptName      string                `json:"dept_name"`
	Level         int                   `json:"level"`
	EmployeeCount int                   `json:"employee_count"`
	LeaderID      *string               `json:"leader_id,omitempty"`
	LeaderName    *string               `json:"leader_name,omitempty"`
	Children      []*DepartmentTreeNode `json:"children,omitempty"`
}

// EmployeeListItem 员工列表项（用于通讯录展示）
type EmployeeListItem struct {
	EmpID       string               `json:"emp_id"`
	EmpName     string               `json:"emp_name"`
	EmpCode     string               `json:"emp_code"`
	OrgName     string               `json:"org_name"`
	DeptName    *string              `json:"dept_name,omitempty"`
	PositionName *string             `json:"position_name,omitempty"`
	JobTitle    string               `json:"job_title"`
	JobLevel    int                  `json:"job_level"`
	Phone       *string              `json:"phone,omitempty"`
	Email       *string              `json:"email,omitempty"`
	AvatarURL   *string              `json:"avatar_url,omitempty"`
	EmployeeStatus entity.EmployeeStatus `json:"employee_status"`
}

// GetOrganizationDirectory 获取组织架构目录
func (s *DirectoryService) GetOrganizationDirectory(ctx context.Context, tenantID string) (*OrganizationTreeNode, error) {
	if tenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取所有组织
	orgs, err := s.orgRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get organizations"),
            )
	}

	// 2. 构建组织树并统计员工数
	tree, err := s.buildOrgDirectory(ctx, orgs, "")
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "build org directory"),
		)
	}
	return tree, nil
}

// buildOrgDirectory 构建组织目录树
func (s *DirectoryService) buildOrgDirectory(ctx context.Context, orgs []*entity.Organization, parentID string) (*OrganizationTreeNode, error) {
	// 查找子组织
	var children []*entity.Organization
	for _, org := range orgs {
		isChild := false
		if parentID == "" {
			// 根节点
			isChild = org.ParentID == nil || *org.ParentID == ""
		} else {
			// 子节点
			isChild = org.ParentID != nil && *org.ParentID == parentID
		}

		if isChild {
			children = append(children, org)
		}
	}

	// 如果没有子组织，返回nil
	if len(children) == 0 {
		return nil, nil
	}

	// 如果有多个子组织，返回虚拟根节点
	if len(children) > 1 || parentID != "" {
		nodes := make([]*OrganizationTreeNode, 0, len(children))
		for _, org := range children {
			// 统计员工数
			empCount, err := s.countEmployeesByOrg(ctx, org.OrgID)
			if err != nil {
				return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "count employees"),
            )
			}

			// 递归构建子节点
			childNode, err := s.buildOrgDirectory(ctx, orgs, org.OrgID)
			if err != nil {
				return nil, err
			}

			node := &OrganizationTreeNode{
				OrgID:         org.OrgID,
				OrgName:       org.OrgName,
				OrgType:       org.OrgType,
				Level:         org.Level,
				EmployeeCount: empCount,
			}

			if childNode != nil {
				node.Children = []*OrganizationTreeNode{childNode}
			}

			nodes = append(nodes, node)
		}

		// 如果是根级别，返回第一个节点作为根（通常只有一个公司）
		if parentID == "" && len(nodes) == 1 {
			return nodes[0], nil
		}

		// 否则返回虚拟根节点
		return &OrganizationTreeNode{
			OrgID:         "root",
			OrgName:       "Root",
			OrgType:       "",
			Level:         0,
			EmployeeCount: 0,
			Children:      nodes,
		}, nil
	}

	// 单个组织节点
	org := children[0]
	empCount, err := s.countEmployeesByOrg(ctx, org.OrgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "count employees"),
            )
	}

	node := &OrganizationTreeNode{
		OrgID:         org.OrgID,
		OrgName:       org.OrgName,
		OrgType:       org.OrgType,
		Level:         org.Level,
		EmployeeCount: empCount,
	}

	// 递归处理子节点
	childNode, err := s.buildOrgDirectory(ctx, orgs, org.OrgID)
	if err != nil {
		return nil, err
	}
	if childNode != nil {
		node.Children = []*OrganizationTreeNode{childNode}
	}

	return node, nil
}

// GetDepartmentDirectory 获取部门目录
func (s *DirectoryService) GetDepartmentDirectory(ctx context.Context, tenantID string) (*DepartmentTreeNode, error) {
	if tenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取所有部门
	depts, err := s.deptRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get departments"),
            )
	}

	// 2. 构建部门树并统计员工数
	tree, err := s.buildDeptDirectory(ctx, depts, "")
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "build dept directory"),
		)
	}
	return tree, nil
}

// buildDeptDirectory 构建部门目录树
func (s *DirectoryService) buildDeptDirectory(ctx context.Context, depts []*entity.Department, parentID string) (*DepartmentTreeNode, error) {
	// 查找子部门
	var children []*entity.Department
	for _, dept := range depts {
		isChild := false
		if parentID == "" {
			// 根节点
			isChild = dept.ParentID == nil || *dept.ParentID == ""
		} else {
			// 子节点
			isChild = dept.ParentID != nil && *dept.ParentID == parentID
		}

		if isChild {
			children = append(children, dept)
		}
	}

	// 如果没有子部门，返回nil
	if len(children) == 0 {
		return nil, nil
	}

	// 多个根级别部门
	if len(children) > 1 || parentID != "" {
		nodes := make([]*DepartmentTreeNode, 0, len(children))
		for _, dept := range children {
			// 统计员工数
			empCount, err := s.countEmployeesByDept(ctx, dept.DeptID)
			if err != nil {
				return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "count employees"),
            )
			}

			// 递归构建子节点
			childNode, err := s.buildDeptDirectory(ctx, depts, dept.DeptID)
			if err != nil {
				return nil, err
			}

			// 获取负责人信息
			var leaderID, leaderName *string
			if dept.LeaderID != nil {
				leaderID = dept.LeaderID
				leader, err := s.empRepo.GetByID(ctx, *dept.LeaderID)
				if err == nil && leader != nil {
					leaderName = &leader.EmpName
				}
			}

			node := &DepartmentTreeNode{
				DeptID:        dept.DeptID,
				DeptName:      dept.DeptName,
				Level:         dept.Level,
				EmployeeCount: empCount,
				LeaderID:      leaderID,
				LeaderName:    leaderName,
			}

			if childNode != nil {
				node.Children = []*DepartmentTreeNode{childNode}
			}

			nodes = append(nodes, node)
		}

		// 返回第一个节点作为根（简化处理）
		if parentID == "" && len(nodes) > 0 {
			return nodes[0], nil
		}

		// 虚拟根节点
		return &DepartmentTreeNode{
			DeptID:        "root",
			DeptName:      "Root",
			Level:         0,
			EmployeeCount: 0,
			Children:      nodes,
		}, nil
	}

	// 单个部门节点
	dept := children[0]
	empCount, err := s.countEmployeesByDept(ctx, dept.DeptID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "count employees"),
            )
	}

	// 获取负责人信息
	var leaderID, leaderName *string
	if dept.LeaderID != nil {
		leaderID = dept.LeaderID
		leader, err := s.empRepo.GetByID(ctx, *dept.LeaderID)
		if err == nil && leader != nil {
			leaderName = &leader.EmpName
		}
	}

	node := &DepartmentTreeNode{
		DeptID:        dept.DeptID,
		DeptName:      dept.DeptName,
		Level:         dept.Level,
		EmployeeCount: empCount,
		LeaderID:      leaderID,
		LeaderName:    leaderName,
	}

	// 递归处理子节点
	childNode, err := s.buildDeptDirectory(ctx, depts, dept.DeptID)
	if err != nil {
		return nil, err
	}
	if childNode != nil {
		node.Children = []*DepartmentTreeNode{childNode}
	}

	return node, nil
}

// GetDepartmentEmployees 获取部门员工列表
func (s *DirectoryService) GetDepartmentEmployees(ctx context.Context, deptID string) ([]*EmployeeListItem, error) {
	if deptID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取部门信息
	dept, err := s.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get department"),
            )
	}
	if dept == nil {
		return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
	}

	// 2. 获取员工
	emps, err := s.empRepo.GetByDepartmentID(ctx, deptID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employees"),
            )
	}

	// 3. 转换为列表项
	return s.convertToEmployeeListItems(ctx, emps)
}

// GetOrganizationEmployees 获取组织员工列表
func (s *DirectoryService) GetOrganizationEmployees(ctx context.Context, orgID string) ([]*EmployeeListItem, error) {
	if orgID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取组织信息
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get organization"),
            )
	}
	if org == nil {
		return nil, errorx.NewByErrorCode(errno.ErrOrgNotFound)
	}

	// 2. 获取员工
	emps, err := s.empRepo.GetByOrgID(ctx, orgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employees"),
            )
	}

	// 3. 转换为列表项
	return s.convertToEmployeeListItems(ctx, emps)
}

// SearchEmployees 搜索员工（用于通讯录搜索）
func (s *DirectoryService) SearchEmployees(ctx context.Context, tenantID, keyword string, limit int) ([]*EmployeeListItem, error) {
	if tenantID == "" || keyword == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 1. 搜索员工
	emps, err := s.empRepo.Search(ctx, tenantID, keyword, limit)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "search employees"),
            )
	}

	// 2. 转换为列表项
	return s.convertToEmployeeListItems(ctx, emps)
}

// convertToEmployeeListItems 转换为员工列表项
func (s *DirectoryService) convertToEmployeeListItems(ctx context.Context, emps []*entity.Employee) ([]*EmployeeListItem, error) {
	items := make([]*EmployeeListItem, 0, len(emps))

	for _, emp := range emps {
		// 获取组织名称
		org, _ := s.orgRepo.GetByID(ctx, emp.OrgID)
		orgName := ""
		if org != nil {
			orgName = org.OrgName
		}

		// 获取部门名称
		var deptName *string
		if emp.DeptID != nil && *emp.DeptID != "" {
			dept, _ := s.deptRepo.GetByID(ctx, *emp.DeptID)
			if dept != nil {
				deptName = &dept.DeptName
			}
		}

		// 获取岗位名称
		var positionName *string
		if emp.PositionID != nil && *emp.PositionID != "" {
			position, _ := s.posRepo.GetByID(ctx, *emp.PositionID)
			if position != nil {
				positionName = &position.PositionName
			}
		}

		items = append(items, &EmployeeListItem{
			EmpID:          emp.EmpID,
			EmpName:        emp.EmpName,
			EmpCode:        emp.EmpCode,
			OrgName:        orgName,
			DeptName:       deptName,
			PositionName:   positionName,
			JobTitle:       emp.JobTitle,
			JobLevel:       emp.JobLevel,
			Phone:          emp.Phone,
			Email:          emp.Email,
			AvatarURL:      emp.AvatarURL,
			EmployeeStatus: emp.EmployeeStatus,
		})
	}

	return items, nil
}

// countEmployeesByOrg 统计组织的员工数
func (s *DirectoryService) countEmployeesByOrg(ctx context.Context, orgID string) (int, error) {
	emps, err := s.empRepo.GetByOrgID(ctx, orgID)
	if err != nil {
		return 0, err
	}

	// 只统计在职员工
	count := 0
	for _, emp := range emps {
		if emp.IsActive() {
			count++
		}
	}

	return count, nil
}

// countEmployeesByDept 统计部门的员工数
func (s *DirectoryService) countEmployeesByDept(ctx context.Context, deptID string) (int, error) {
	emps, err := s.empRepo.GetByDepartmentID(ctx, deptID)
	if err != nil {
		return 0, err
	}

	// 只统计在职员工
	count := 0
	for _, emp := range emps {
		if emp.IsActive() {
			count++
		}
	}

	return count, nil
}

// GetEmployeeByCode 获取员工信息（用于通讯录）
func (s *DirectoryService) GetEmployeeByCode(ctx context.Context, tenantID, code string) (*EmployeeListItem, error) {
	if tenantID == "" || code == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	emp, err := s.empRepo.GetByCode(ctx, tenantID, code)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	items, err := s.convertToEmployeeListItems(ctx, []*entity.Employee{emp})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	return items[0], nil
}
