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

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// DepartmentPermissionChecker 部门权限检查器
type DepartmentPermissionChecker struct {
	userDeptRepo   repository.UserDepartmentRepository
	departmentRepo repository.DepartmentRepository
}

// NewDepartmentPermissionChecker 创建部门权限检查器实例
func NewDepartmentPermissionChecker(
	userDeptRepo repository.UserDepartmentRepository,
	departmentRepo repository.DepartmentRepository,
) *DepartmentPermissionChecker {
	return &DepartmentPermissionChecker{
		userDeptRepo:   userDeptRepo,
		departmentRepo: departmentRepo,
	}
}

// GetAccessibleDepartmentIDs 获取用户可访问的部门ID列表
func (d *DepartmentPermissionChecker) GetAccessibleDepartmentIDs(ctx context.Context, userID, tenantID string) ([]string, error) {
	// 1. 获取用户所属部门
	userDepts, err := d.userDeptRepo.GetByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	deptMap := make(map[string]bool)
	deptIDs := make([]string, 0)

	for _, ud := range userDepts {
		if !deptMap[ud.DepartmentID] {
			deptMap[ud.DepartmentID] = true
			deptIDs = append(deptIDs, ud.DepartmentID)
		}

		// 2. 如果是部门领导，可以访问子部门
		if ud.IsLeader {
			childDepts, err := d.departmentRepo.GetDescendants(ctx, ud.DepartmentID)
			if err != nil {
				continue // 即使获取子部门失败，也不影响主流程
			}

			for _, child := range childDepts {
				if !deptMap[child.DepartmentID] {
					deptMap[child.DepartmentID] = true
					deptIDs = append(deptIDs, child.DepartmentID)
				}
			}
		}
	}

	return deptIDs, nil
}

// IsDepartmentLeader 检查用户是否是某个部门的领导
func (d *DepartmentPermissionChecker) IsDepartmentLeader(ctx context.Context, userID, departmentID string) (bool, error) {
	userDept, err := d.userDeptRepo.GetByUserAndDepartment(ctx, userID, departmentID)
	if err != nil {
		return false, err
	}

	return userDept != nil && userDept.IsLeader, nil
}

// FilterResourcesByDepartment 按部门过滤资源
func (d *DepartmentPermissionChecker) FilterResourcesByDepartment(
	ctx context.Context,
	userID, tenantID string,
	resources []ResourceWithDepartment,
) ([]ResourceWithDepartment, error) {
	accessibleDeptIDs, err := d.GetAccessibleDepartmentIDs(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	filtered := make([]ResourceWithDepartment, 0)

	for _, resource := range resources {
		// 如果资源有部门ID，且在可访问部门列表中，则保留
		if resource.GetDepartmentID() == "" {
			// 无部门信息的资源，默认保留（可能是系统资源）
			filtered = append(filtered, resource)
			continue
		}

		// 检查部门ID是否在可访问列表中
		for _, accessibleID := range accessibleDeptIDs {
			if resource.GetDepartmentID() == accessibleID {
				filtered = append(filtered, resource)
				break
			}
		}
	}

	return filtered, nil
}

// GetUserDepartments 获取用户的所有部门信息
func (d *DepartmentPermissionChecker) GetUserDepartments(ctx context.Context, userID, tenantID string) ([]*UserDepartmentInfo, error) {
	userDepts, err := d.userDeptRepo.GetByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]*UserDepartmentInfo, 0, len(userDepts))

	for _, ud := range userDepts {
		dept, err := d.departmentRepo.GetByID(ctx, ud.DepartmentID)
		if err != nil {
			continue
		}

		result = append(result, &UserDepartmentInfo{
			DepartmentID:   dept.DepartmentID,
			DepartmentName: dept.DepartmentName,
			IsLeader:       ud.IsLeader,
			ParentID:       dept.ParentDepartmentID,
		})
	}

	return result, nil
}

// GetDepartmentTree 获取完整的部门树
func (d *DepartmentPermissionChecker) GetDepartmentTree(ctx context.Context, tenantID string) ([]*DepartmentTreeNode, error) {
	// 获取所有部门
	depts, err := d.departmentRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 构建部门树
	return d.buildDepartmentTree(depts, ""), nil
}

// buildDepartmentTree 递归构建部门树
func (d *DepartmentPermissionChecker) buildDepartmentTree(depts []*entity.Department, parentID string) []*DepartmentTreeNode {
	nodes := make([]*DepartmentTreeNode, 0)

	for _, dept := range depts {
		// 跳过已删除的部门
		if dept.DeletedAt != nil {
			continue
		}

		// 匹配父部门ID
		hasParent := (parentID == "" && dept.ParentDepartmentID == nil) ||
			(parentID != "" && dept.ParentDepartmentID != nil && *dept.ParentDepartmentID == parentID)

		if !hasParent {
			continue
		}

		node := &DepartmentTreeNode{
			DepartmentID:   dept.DepartmentID,
			DepartmentName: dept.DepartmentName,
			ParentID:       dept.ParentDepartmentID,
			Children:       d.buildDepartmentTree(depts, dept.DepartmentID),
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// ResourceWithDepartment 带部门信息的资源接口
type ResourceWithDepartment interface {
	GetDepartmentID() string
}

// UserDepartmentInfo 用户部门信息
type UserDepartmentInfo struct {
	DepartmentID   string
	DepartmentName string
	IsLeader       bool
	ParentID       *string
}

// DepartmentTreeNode 部门树节点
type DepartmentTreeNode struct {
	DepartmentID   string
	DepartmentName string
	ParentID       *string
	Children       []*DepartmentTreeNode
}

// Department 部门信息（简化版）
type Department struct {
	DepartmentID         string
	TenantID             string
	DepartmentName       string
	ParentDepartmentID   *string
	CreatedAt            int64
	UpdatedAt            int64
	DeletedAt            *int64
}
