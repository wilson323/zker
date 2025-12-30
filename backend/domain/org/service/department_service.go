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
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// DepartmentService 部门管理服务
// 职责：部门的CRUD操作、层级管理、移动操作、业务规则验证
type DepartmentService struct {
	deptRepo    repository.DepartmentRepository
	treeRepo    repository.DepartmentTreeRepository
	orgRepo     repository.OrganizationRepository
	db          *gorm.DB
}

// NewDepartmentService 创建部门服务实例
func NewDepartmentService(
	deptRepo repository.DepartmentRepository,
	treeRepo repository.DepartmentTreeRepository,
	orgRepo repository.OrganizationRepository,
	db *gorm.DB,
) *DepartmentService {
	return &DepartmentService{
		deptRepo: deptRepo,
		treeRepo: treeRepo,
		orgRepo:  orgRepo,
		db:       db,
	}
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	TenantID      string  `json:"tenant_id" binding:"required"`
	OrgID         string  `json:"org_id" binding:"required"`
	ParentID      *string `json:"parent_id,omitempty"`
	DeptName      string  `json:"dept_name" binding:"required,min=1,max=200"`
	DeptCode      string  `json:"dept_code" binding:"required,min=1,max=50"`
	LeaderID      *string `json:"leader_id,omitempty"`
	ParentLeader  *string `json:"parent_leader,omitempty"`
	Description   string  `json:"description,omitempty"`
	SortOrder     int     `json:"sort_order,omitempty"`
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	DeptID       string  `json:"dept_id" binding:"required"`
	DeptName     string  `json:"dept_name" binding:"omitempty,min=1,max=200"`
	LeaderID     *string `json:"leader_id,omitempty"`
	ParentLeader *string `json:"parent_leader,omitempty"`
	Description  string  `json:"description,omitempty"`
	SortOrder    int     `json:"sort_order,omitempty"`
	Status       string  `json:"status" binding:"omitempty,oneof=active inactive frozen"`
}

// MoveDepartmentRequest 移动部门请求
type MoveDepartmentRequest struct {
	DeptID      string  `json:"dept_id" binding:"required"`
	NewParentID *string `json:"new_parent_id,omitempty"`
}

// CreateDepartment 创建部门
func (s *DepartmentService) CreateDepartment(ctx context.Context, req *CreateDepartmentRequest) (*entity.Department, error) {
	// 1. 验证组织存在
	org, err := s.orgRepo.GetByID(ctx, req.OrgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	if org == nil {
		return nil, errno.ErrOrgNotFound
	}
	if org.TenantID != req.TenantID {
		return nil, errno.ErrTenantMismatch
	}

	// 2. 如果有父部门，验证父部门存在且属于同一组织
	var parent *entity.Department
	var level int
	var path string

	if req.ParentID != nil && *req.ParentID != "" {
		parent, err = s.deptRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent department: %w", err)
		}
		if parent == nil {
			return nil, errno.ErrParentDeptNotFound
		}
		if parent.OrgID != req.OrgID {
			return nil, errno.ErrDeptOrgMismatch
		}

		level = parent.Level + 1
		path = parent.Path + "/"
	} else {
		level = 1
		path = "/"
	}

	// 3. 验证编码唯一性
	exists, err := s.deptRepo.ExistsByCode(ctx, req.TenantID, req.DeptCode, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check dept code: %w", err)
	}
	if exists {
		return nil, errno.ErrDeptCodeAlreadyExists
	}

	// 4. 生成部门ID
	deptID := generateUUID()

	// 5. 创建部门实体
	dept := &entity.Department{
		DeptID:       deptID,
		TenantID:     req.TenantID,
		OrgID:        req.OrgID,
		ParentID:     req.ParentID,
		DeptName:     req.DeptName,
		DeptCode:     req.DeptCode,
		Level:        level,
		Path:         path + deptID,
		SortOrder:    req.SortOrder,
		Status:       entity.OrgStatusActive,
		LeaderID:     req.LeaderID,
		ParentLeader: req.ParentLeader,
		Description:  req.Description,
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	// 6. 使用事务创建部门和路径
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建部门
		if err := s.deptRepo.Create(ctx, dept); err != nil {
			return fmt.Errorf("failed to create department: %w", err)
		}

		// 创建闭包表路径
		if err := s.createDepartmentPaths(ctx, dept, parent); err != nil {
			return fmt.Errorf("failed to create department paths: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dept, nil
}

// createDepartmentPaths 创建部门的所有闭包表路径
func (s *DepartmentService) createDepartmentPaths(ctx context.Context, dept *entity.Department, parent *entity.Department) error {
	paths := make([]*entity.DepartmentTree, 0, dept.Level)

	// 添加自己到自己的路径（深度0）
	paths = append(paths, &entity.DepartmentTree{
		TenantID:     dept.TenantID,
		AncestorID:   dept.DeptID,
		DescendantID: dept.DeptID,
		Depth:        0,
	})

	// 如果有父部门，添加所有祖先的路径
	if parent != nil {
		ancestorPaths, err := s.treeRepo.GetAncestors(ctx, parent.DeptID)
		if err != nil {
			return fmt.Errorf("failed to get ancestor paths: %w", err)
		}

		for _, ancestorPath := range ancestorPaths {
			paths = append(paths, &entity.DepartmentTree{
				TenantID:     dept.TenantID,
				AncestorID:   ancestorPath.AncestorID,
				DescendantID: dept.DeptID,
				Depth:        ancestorPath.Depth + 1,
			})
		}
	}

	// 批量创建路径
	return s.treeRepo.CreatePaths(ctx, paths)
}

// GetDepartment 获取部门详情
func (s *DepartmentService) GetDepartment(ctx context.Context, deptID string) (*entity.Department, error) {
	if deptID == "" {
		return nil, errno.ErrInvalidParam
	}

	dept, err := s.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	if dept == nil {
		return nil, errno.ErrDeptNotFound
	}

	return dept, nil
}

// UpdateDepartment 更新部门
func (s *DepartmentService) UpdateDepartment(ctx context.Context, req *UpdateDepartmentRequest) (*entity.Department, error) {
	// 1. 获取部门
	dept, err := s.deptRepo.GetByID(ctx, req.DeptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	if dept == nil {
		return nil, errno.ErrDeptNotFound
	}

	// 2. 更新字段
	if req.DeptName != "" {
		dept.DeptName = req.DeptName
	}
	if req.LeaderID != nil {
		dept.LeaderID = req.LeaderID
	}
	if req.ParentLeader != nil {
		dept.ParentLeader = req.ParentLeader
	}
	if req.Description != "" {
		dept.Description = req.Description
	}
	if req.SortOrder != 0 {
		dept.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		dept.Status = entity.OrganizationStatus(req.Status)
	}

	// 3. 保存更新
	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return dept, nil
}

// DeleteDepartment 删除部门
func (s *DepartmentService) DeleteDepartment(ctx context.Context, deptID string) error {
	if deptID == "" {
		return errno.ErrInvalidParam
	}

	// 1. 获取部门
	dept, err := s.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return fmt.Errorf("failed to get department: %w", err)
	}
	if dept == nil {
		return errno.ErrDeptNotFound
	}

	// 2. 检查是否有子部门
	children, err := s.deptRepo.GetChildren(ctx, deptID)
	if err != nil {
		return fmt.Errorf("failed to check children: %w", err)
	}
	if len(children) > 0 {
		return errno.ErrDeptHasChildren
	}

	// 3. 使用事务删除部门和路径
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 软删除部门
		if err := s.deptRepo.Delete(ctx, deptID); err != nil {
			return fmt.Errorf("failed to delete department: %w", err)
		}

		// 删除闭包表路径
		if err := s.treeRepo.DeletePaths(ctx, deptID); err != nil {
			return fmt.Errorf("failed to delete department paths: %w", err)
		}

		return nil
	})
}

// MoveDepartment 移动部门
func (s *DepartmentService) MoveDepartment(ctx context.Context, req *MoveDepartmentRequest) error {
	if req.DeptID == "" {
		return errno.ErrInvalidParam
	}

	// 1. 获取要移动的部门
	dept, err := s.deptRepo.GetByID(ctx, req.DeptID)
	if err != nil {
		return fmt.Errorf("failed to get department: %w", err)
	}
	if dept == nil {
		return errno.ErrDeptNotFound
	}

	// 2. 验证不能移动到自己
	if req.NewParentID != nil && *req.NewParentID == dept.DeptID {
		return errno.ErrCannotMoveToSelf
	}

	// 3. 验证新父部门存在且属于同一组织
	var newParent *entity.Department
	var newLevel int
	var newPath string

	if req.NewParentID != nil && *req.NewParentID != "" {
		newParent, err = s.deptRepo.GetByID(ctx, *req.NewParentID)
		if err != nil {
			return fmt.Errorf("failed to get new parent: %w", err)
		}
		if newParent == nil {
			return errno.ErrParentDeptNotFound
		}
		if newParent.OrgID != dept.OrgID {
			return errno.ErrDeptOrgMismatch
		}

		// 验证不能移动到自己的后代
		isDescendant, err := s.isDescendant(ctx, dept.DeptID, newParent.DeptID)
		if err != nil {
			return fmt.Errorf("failed to check descendant: %w", err)
		}
		if isDescendant {
			return errno.ErrCannotMoveToDescendant
		}

		newLevel = newParent.Level + 1
		newPath = newParent.Path + "/" + dept.DeptID
	} else {
		// 移动到根级别
		newLevel = 1
		newPath = "/" + dept.DeptID
	}

	// 4. 使用事务移动部门和更新所有后代的路径
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新部门的层级和路径
		dept.ParentID = req.NewParentID
		dept.Level = newLevel
		dept.Path = newPath
		dept.UpdatedAt = time.Now().UnixMilli()

		if err := s.deptRepo.Update(ctx, dept); err != nil {
			return fmt.Errorf("failed to update department: %w", err)
		}

		// 更新闭包表路径
		newParentID := ""
		if req.NewParentID != nil {
			newParentID = *req.NewParentID
		}
		if err := s.treeRepo.MoveSubtree(ctx, dept.DeptID, newParentID); err != nil {
			return fmt.Errorf("failed to move subtree: %w", err)
		}

		return nil
	})
}

// isDescendant 检查是否是后代节点
func (s *DepartmentService) isDescendant(ctx context.Context, ancestorID, descendantID string) (bool, error) {
	paths, err := s.treeRepo.GetAncestors(ctx, descendantID)
	if err != nil {
		return false, err
	}

	for _, path := range paths {
		if path.AncestorID == ancestorID && path.Depth > 0 {
			return true, nil
		}
	}

	return false, nil
}

// GetDepartmentTree 获取部门树
func (s *DepartmentService) GetDepartmentTree(ctx context.Context, tenantID string) ([]*entity.Department, error) {
	if tenantID == "" {
		return nil, errno.ErrInvalidParam
	}

	depts, err := s.deptRepo.GetTree(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get department tree: %w", err)
	}

	// 构建树形结构
	return s.buildDeptTree(depts, ""), nil
}

// buildDeptTree 构建部门树
func (s *DepartmentService) buildDeptTree(depts []*entity.Department, parentID string) []*entity.Department {
	var tree []*entity.Department

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
			dept.Children = s.buildDeptTree(depts, dept.DeptID)
			tree = append(tree, dept)
		}
	}

	return tree
}

// ListDepartments 分页查询部门列表
func (s *DepartmentService) ListDepartments(ctx context.Context, filter *repository.DepartmentFilter) ([]*entity.Department, int64, error) {
	if filter.TenantID == "" {
		return nil, 0, errno.ErrInvalidParam
	}

	depts, total, err := s.deptRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list departments: %w", err)
	}

	return depts, total, nil
}

// GetChildren 获取子部门
func (s *DepartmentService) GetChildren(ctx context.Context, parentID string) ([]*entity.Department, error) {
	if parentID == "" {
		return nil, errno.ErrInvalidParam
	}

	depts, err := s.deptRepo.GetChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}

	return depts, nil
}

// GetAncestors 获取部门的所有祖先
func (s *DepartmentService) GetAncestors(ctx context.Context, deptID string) ([]*entity.Department, error) {
	if deptID == "" {
		return nil, errno.ErrInvalidParam
	}

	paths, err := s.treeRepo.GetAncestors(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ancestor paths: %w", err)
	}

	// 转换为部门实体
	var ancestors []*entity.Department
	for _, path := range paths {
		if path.Depth > 0 { // 排除自己
			dept, err := s.deptRepo.GetByID(ctx, path.AncestorID)
			if err != nil {
				return nil, fmt.Errorf("failed to get ancestor department: %w", err)
			}
			if dept != nil {
				ancestors = append(ancestors, dept)
			}
		}
	}

	return ancestors, nil
}

// GetDescendants 获取部门的所有后代
func (s *DepartmentService) GetDescendants(ctx context.Context, deptID string) ([]*entity.Department, error) {
	if deptID == "" {
		return nil, errno.ErrInvalidParam
	}

	paths, err := s.treeRepo.GetDescendants(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get descendant paths: %w", err)
	}

	// 转换为部门实体
	var descendants []*entity.Department
	for _, path := range paths {
		if path.Depth > 0 { // 排除自己
			dept, err := s.deptRepo.GetByID(ctx, path.DescendantID)
			if err != nil {
				return nil, fmt.Errorf("failed to get descendant department: %w", err)
			}
			if dept != nil {
				descendants = append(descendants, dept)
			}
		}
	}

	return descendants, nil
}

// GetDepartmentsByOrg 获取组织的所有部门
func (s *DepartmentService) GetDepartmentsByOrg(ctx context.Context, orgID string) ([]*entity.Department, error) {
	if orgID == "" {
		return nil, errno.ErrInvalidParam
	}

	filter := &repository.DepartmentFilter{
		OrgID: orgID,
	}

	depts, _, err := s.deptRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments by org: %w", err)
	}

	return depts, nil
}
