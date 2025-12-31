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
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// DataPermissionService 数据权限服务
// 职责：五级数据权限控制、数据过滤、权限检查
type DataPermissionService struct {
	permRepo repository.DataPermissionRepository
	orgRepo  repository.OrganizationRepository
	empRepo  repository.EmployeeRepository
	deptRepo repository.DepartmentRepository
}

// NewDataPermissionService 创建数据权限服务实例
func NewDataPermissionService(
	permRepo repository.DataPermissionRepository,
	orgRepo repository.OrganizationRepository,
	empRepo repository.EmployeeRepository,
	deptRepo repository.DepartmentRepository,
) *DataPermissionService {
	return &DataPermissionService{
		permRepo: permRepo,
		orgRepo:  orgRepo,
		empRepo:  empRepo,
		deptRepo: deptRepo,
	}
}

// GetUserPermissionLevel 获取用户的数据权限级别
func (s *DataPermissionService) GetUserPermissionLevel(ctx context.Context, userID, resourceType string) (entity.DataPermissionLevel, error) {
	level, err := s.permRepo.GetUserPermissionLevel(ctx, userID, resourceType)
	if err != nil {
		return entity.DataPermLevelNone, errorx.Wrapf(err, "get user permission level failed: user_id=%s, resource_type=%s", userID, resourceType)
	}

	return entity.DataPermissionLevel(level), nil
}

// FilterDataByPermission 根据权限过滤数据
func (s *DataPermissionService) FilterDataByPermission(
	ctx context.Context,
	userID string,
	resourceType string,
	data []*entity.Employee,
) ([]*entity.Employee, error) {
	// 1. 获取用户权限级别
	level, err := s.GetUserPermissionLevel(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	// 2. 根据权限级别过滤数据
	switch level {
	case entity.DataPermLevelAll:
		// 全部数据权限，返回所有数据
		return data, nil

	case entity.DataPermLevelDeptAndBelow:
		// 本部门及下级部门数据
		return s.filterByDeptAndBelow(ctx, userID, data)

	case entity.DataPermLevelDept:
		// 本部门数据
		return s.filterByDept(ctx, userID, data)

	case entity.DataPermLevelSelf:
		// 本人数据
		return s.filterBySelf(ctx, userID, data)

	case entity.DataPermLevelNone:
		// 无数据权限
		return []*entity.Employee{}, nil

	default:
		return nil, errorx.New(errno.ErrPermissionCheckFailedCode)
	}
}

// filterByDeptAndBelow 过滤本部门及下级部门数据
func (s *DataPermissionService) filterByDeptAndBelow(
	ctx context.Context,
	userID string,
	data []*entity.Employee,
) ([]*entity.Employee, error) {
	// 1. 获取用户可访问的部门列表（包含下级部门）
	deptIDs, err := s.permRepo.GetAccessibleDepartments(ctx, userID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get accessible departments failed: user_id=%s", userID)
	}

	// 2. 过滤数据
	result := make([]*entity.Employee, 0)
	for _, emp := range data {
		if emp.DeptID != nil {
			for _, deptID := range deptIDs {
				if *emp.DeptID == deptID {
					result = append(result, emp)
					break
				}
			}
		}
	}

	return result, nil
}

// filterByDept 过滤本部门数据
func (s *DataPermissionService) filterByDept(
	ctx context.Context,
	userID string,
	data []*entity.Employee,
) ([]*entity.Employee, error) {
	// 1. 获取员工信息
	emp, err := s.empRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get employee info failed: user_id=%s", userID)
	}

	// 2. 如果员工没有部门，返回空
	if emp.DeptID == nil || *emp.DeptID == "" {
		return []*entity.Employee{}, nil
	}

	// 3. 过滤同部门数据
	result := make([]*entity.Employee, 0)
	for _, e := range data {
		if e.DeptID != nil && *e.DeptID == *emp.DeptID {
			result = append(result, e)
		}
	}

	return result, nil
}

// filterBySelf 过滤本人数据
func (s *DataPermissionService) filterBySelf(
	ctx context.Context,
	userID string,
	data []*entity.Employee,
) ([]*entity.Employee, error) {
	result := make([]*entity.Employee, 0)
	for _, emp := range data {
		if emp.EmpID == userID {
			result = append(result, emp)
			break
		}
	}
	return result, nil
}

// CheckDataPermission 检查数据权限
func (s *DataPermissionService) CheckDataPermission(
	ctx context.Context,
	userID, resourceType, resourceID string,
) (bool, error) {
	// 1. 检查数据权限
	hasPermission, err := s.permRepo.CheckDataPermission(ctx, userID, resourceType, resourceID)
	if err != nil {
		return false, errorx.Wrapf(err, "check data permission failed: user_id=%s, resource_type=%s, resource_id=%s", userID, resourceType, resourceID)
	}

	return hasPermission, nil
}

// GetAccessibleDepartments 获取用户可访问的部门列表
func (s *DataPermissionService) GetAccessibleDepartments(
	ctx context.Context,
	userID string,
) ([]string, error) {
	deptIDs, err := s.permRepo.GetAccessibleDepartments(ctx, userID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get accessible departments failed: user_id=%s", userID)
	}

	return deptIDs, nil
}

// GetAccessibleOrganizations 获取用户可访问的组织列表
func (s *DataPermissionService) GetAccessibleOrganizations(
	ctx context.Context,
	userID string,
) ([]string, error) {
	orgIDs, err := s.permRepo.GetAccessibleOrganizations(ctx, userID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get accessible organizations failed: user_id=%s", userID)
	}

	return orgIDs, nil
}

// SetDataPermissionLevel 设置用户的数据权限级别
func (s *DataPermissionService) SetDataPermissionLevel(
	ctx context.Context,
	req *SetDataPermissionRequest,
) error {
	// 1. 验证权限级别
	if req.Level < entity.DataPermLevelNone || req.Level > entity.DataPermLevelAll {
		return errorx.New(errno.ErrInvalidParamCode)
	}

	// 2. 验证用户存在
	_, err := s.empRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return errorx.Wrapf(err, "validate employee failed: user_id=%s", req.UserID)
	}

	// 3. 设置权限（实际实现需要调用权限领域服务）
	// 这里简化处理，实际应该调用 permission domain 的服务

	return nil
}

// GetDataPermissionSummary 获取数据权限摘要
func (s *DataPermissionService) GetDataPermissionSummary(
	ctx context.Context,
	userID string,
) (*DataPermissionSummary, error) {
	// 1. 获取权限级别
	level, err := s.GetUserPermissionLevel(ctx, userID, "employee")
	if err != nil {
		return nil, err
	}

	// 2. 获取可访问的组织和部门
	orgIDs, err := s.GetAccessibleOrganizations(ctx, userID)
	if err != nil {
		return nil, err
	}

	deptIDs, err := s.GetAccessibleDepartments(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. 构建摘要
	summary := &DataPermissionSummary{
		UserID:             userID,
		PermissionLevel:    level,
		AccessibleOrgCount: len(orgIDs),
		AccessibleDeptCount: len(deptIDs),
		AccessibleOrgs:     orgIDs,
		AccessibleDepts:    deptIDs,
	}

	return summary, nil
}

// BatchCheckPermission 批量检查权限
func (s *DataPermissionService) BatchCheckPermission(
	ctx context.Context,
	req *BatchCheckPermissionRequest,
) (map[string]bool, error) {
	result := make(map[string]bool)

	for _, resourceID := range req.ResourceIDs {
		hasPermission, err := s.CheckDataPermission(ctx, req.UserID, req.ResourceType, resourceID)
		if err != nil {
			return nil, err
		}
		result[resourceID] = hasPermission
	}

	return result, nil
}

// SetDataPermissionRequest 设置数据权限请求
type SetDataPermissionRequest struct {
	UserID       string                    `json:"user_id" binding:"required"`
	Level        entity.DataPermissionLevel `json:"level" binding:"required,min=1,max=5"`
	ResourceType string                    `json:"resource_type" binding:"required"`
	Reason       string                    `json:"reason,omitempty"`
}

// BatchCheckPermissionRequest 批量检查权限请求
type BatchCheckPermissionRequest struct {
	UserID       string   `json:"user_id" binding:"required"`
	ResourceType string   `json:"resource_type" binding:"required"`
	ResourceIDs  []string `json:"resource_ids" binding:"required"`
}

// DataPermissionSummary 数据权限摘要
type DataPermissionSummary struct {
	UserID             string                    `json:"user_id"`
	PermissionLevel    entity.DataPermissionLevel `json:"permission_level"`
	LevelName          string                    `json:"level_name"`
	AccessibleOrgCount int                       `json:"accessible_org_count"`
	AccessibleDeptCount int                       `json:"accessible_dept_count"`
	AccessibleOrgs     []string                  `json:"accessible_orgs,omitempty"`
	AccessibleDepts    []string                  `json:"accessible_depts,omitempty"`
}

// GetLevelName 获取权限级别名称
func (s *DataPermissionService) GetLevelName(level entity.DataPermissionLevel) string {
	names := map[entity.DataPermissionLevel]string{
		entity.DataPermLevelAll:          "全部数据",
		entity.DataPermLevelDeptAndBelow: "本部门及下级部门",
		entity.DataPermLevelDept:         "本部门",
		entity.DataPermLevelSelf:         "本人数据",
		entity.DataPermLevelNone:         "无数据权限",
	}

	if name, ok := names[level]; ok {
		return name
	}
	return "未知级别"
}

// ValidateDataAccess 验证数据访问权限（用于中间件）
func (s *DataPermissionService) ValidateDataAccess(
	ctx context.Context,
	userID, resourceType, resourceID string,
) error {
	hasPermission, err := s.CheckDataPermission(ctx, userID, resourceType, resourceID)
	if err != nil {
		return err
	}

	if !hasPermission {
		return errorx.New(errno.ErrPermissionCheckFailedCode, errorx.KV("resource_id", resourceID))
	}

	return nil
}
