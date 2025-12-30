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
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
)

// PermissionDeniedError 权限拒绝错误
type PermissionDeniedError struct {
	UserID       string
	ResourceType string
	ResourceID   string
	Reason       string
}

func (e *PermissionDeniedError) Error() string {
	return fmt.Sprintf("permission denied for user %s on %s/%s: %s",
		e.UserID, e.ResourceType, e.ResourceID, e.Reason)
}

// PermissionChecker 权限检查服务
type PermissionChecker struct {
	db             *gorm.DB
	roleRepo       repository.RoleRepository
	dataPermRepo   repository.DataPermissionRepository
	fieldPermRepo  repository.FieldPermissionRepository
	userRoleRepo   repository.UserRoleRepository
	userDeptRepo   repository.UserDepartmentRepository
	departmentRepo repository.DepartmentRepository
}

// NewPermissionChecker 创建权限检查服务实例
func NewPermissionChecker(
	db *gorm.DB,
	roleRepo repository.RoleRepository,
	dataPermRepo repository.DataPermissionRepository,
	fieldPermRepo repository.FieldPermissionRepository,
	userRoleRepo repository.UserRoleRepository,
	userDeptRepo repository.UserDepartmentRepository,
	departmentRepo repository.DepartmentRepository,
) *PermissionChecker {
	return &PermissionChecker{
		db:             db,
		roleRepo:       roleRepo,
		dataPermRepo:   dataPermRepo,
		fieldPermRepo:  fieldPermRepo,
		userRoleRepo:   userRoleRepo,
		userDeptRepo:   userDeptRepo,
		departmentRepo: departmentRepo,
	}
}

// CheckDataPermission 检查数据权限（5级权限范围）
//
// 参数说明：
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - userID: 用户ID
//   - resourceType: 资源类型（bots, conversations, knowledge, workflows, plugins）
//   - action: 操作类型（create, read, update, delete）
//   - resourceID: 资源ID
//
// 返回值：
//   - bool: 是否有权限
//   - error: 错误信息
func (p *PermissionChecker) CheckDataPermission(
	ctx context.Context,
	tenantID, userID string,
	resourceType entity.ResourceType,
	action string,
	resourceID string,
) (bool, error) {
	// 1. 参数验证
	if tenantID == "" {
		return false, fmt.Errorf("tenant_id is required")
	}
	if userID == "" {
		return false, fmt.Errorf("user_id is required")
	}
	if resourceType == "" {
		return false, fmt.Errorf("resource_type is required")
	}

	// 2. 获取用户的所有角色
	roles, err := p.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) == 0 {
		return false, &PermissionDeniedError{
			UserID:       userID,
			ResourceType: string(resourceType),
			ResourceID:   resourceID,
			Reason:       "user has no roles",
		}
	}

	// 3. 检查数据权限
	for _, role := range roles {
		perm, err := p.dataPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType)
		if err != nil {
			continue
		}

		if perm == nil {
			continue
		}

		// 评估数据范围权限
		allowed := p.evaluateDataScope(ctx, perm.Scope, userID, resourceType, resourceID, perm.CustomFilter)
		if allowed {
			return true, nil
		}

		// 如果明确是无权限，直接返回
		if perm.Scope == entity.DataPermissionScopeNone {
			return false, &PermissionDeniedError{
				UserID:       userID,
				ResourceType: string(resourceType),
				ResourceID:   resourceID,
				Reason:       "role has no permission for this resource",
			}
		}
	}

	// 所有角色都没有权限
	return false, &PermissionDeniedError{
		UserID:       userID,
		ResourceType: string(resourceType),
		ResourceID:   resourceID,
		Reason:       "no matching permission found",
	}
}

// evaluateDataScope 评估数据范围权限
func (p *PermissionChecker) evaluateDataScope(
	ctx context.Context,
	scope entity.DataPermissionScope,
	userID string,
	resourceType entity.ResourceType,
	resourceID string,
	customFilter string,
) bool {
	switch scope {
	case entity.DataPermissionScopeAll:
		// 全部数据权限
		return true

	case entity.DataPermissionScopeDepartment:
		// 部门数据权限：检查资源是否属于用户部门
		return p.isSameDepartment(ctx, userID, resourceType, resourceID)

	case entity.DataPermissionScopeOwn:
		// 仅自己数据权限：检查资源是否由用户创建
		return p.isOwner(ctx, userID, resourceType, resourceID)

	case entity.DataPermissionScopeCustom:
		// 自定义过滤权限：根据 custom_filter 进行过滤
		return p.matchCustomFilter(ctx, userID, customFilter, resourceType, resourceID)

	case entity.DataPermissionScopeNone:
		// 无权限
		return false

	default:
		// 未知范围，默认拒绝
		return false
	}
}

// GetFieldPermissions 获取字段权限（合并所有角色的字段权限）
//
// 参数说明：
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - userID: 用户ID
//   - resourceType: 资源类型
//
// 返回值：
//   - map[string]string: 字段权限映射（字段名 -> 权限级别：hidden/readonly/editable）
//   - error: 错误信息
func (p *PermissionChecker) GetFieldPermissions(
	ctx context.Context,
	tenantID, userID string,
	resourceType string,
) (map[string]string, error) {
	// 1. 参数验证
	if tenantID == "" || userID == "" || resourceType == "" {
		return nil, fmt.Errorf("tenant_id, user_id and resource_type are required")
	}

	// 2. 获取用户的所有角色
	roles, err := p.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 3. 合并所有角色的字段权限（取最大权限）
	fieldPerms := make(map[string]string)
	for _, role := range roles {
		perms, err := p.fieldPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType)
		if err != nil {
			continue
		}

		for _, perm := range perms {
			fieldName := perm.FieldName
			permLevel := string(perm.PermissionLevel)

			// 如果已有权限，优先级更高：editable > readonly > hidden
			if existing, ok := fieldPerms[fieldName]; ok {
				if p.compareFieldLevel(existing, permLevel) >= 0 {
					// 已有权限更高或相等，跳过
					continue
				}
			}
			fieldPerms[fieldName] = permLevel
		}
	}

	return fieldPerms, nil
}

// FilterResourcesByDepartment 按部门过滤资源
func (p *PermissionChecker) FilterResourcesByDepartment(
	ctx context.Context,
	tenantID, userID string,
	resources []map[string]interface{},
) ([]map[string]interface{}, error) {
	// 1. 获取用户可访问的部门ID列表
	deptIDs, err := p.GetAccessibleDepartmentIDs(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	// 2. 过滤资源
	filtered := make([]map[string]interface{}, 0)
	for _, resource := range resources {
		// 检查资源的 department_id 是否在可访问列表中
		if resourceDeptID, ok := resource["department_id"].(string); ok {
			if contains(deptIDs, resourceDeptID) {
				filtered = append(filtered, resource)
			}
		} else {
			// 如果资源没有 department_id，认为无权限访问
			continue
		}
	}

	return filtered, nil
}

// GetAccessibleDepartmentIDs 获取用户可访问的部门ID列表
func (p *PermissionChecker) GetAccessibleDepartmentIDs(
	ctx context.Context,
	tenantID, userID string,
) ([]string, error) {
	// 1. 获取用户所属部门
	userDepts, err := p.userDeptRepo.GetByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	deptIDs := make([]string, 0)
	for _, ud := range userDepts {
		deptIDs = append(deptIDs, ud.DepartmentID)

		// 2. 如果是部门领导，可以访问子部门
		if ud.IsLeader {
			childDepts, _ := p.departmentRepo.GetDescendants(ctx, ud.DepartmentID)
			for _, child := range childDepts {
				deptIDs = append(deptIDs, child.DepartmentID)
			}
		}
	}

	return deptIDs, nil
}

// UserHasRole 检查用户是否拥有指定角色
//
// 参数说明：
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - userID: 用户ID
//   - roleCode: 角色代码
//
// 返回值：
//   - bool: 是否拥有该角色
//   - error: 错误信息
func (p *PermissionChecker) UserHasRole(
	ctx context.Context,
	tenantID, userID string,
	roleCode string,
) (bool, error) {
	// 1. 参数验证
	if tenantID == "" || userID == "" || roleCode == "" {
		return false, fmt.Errorf("tenant_id, user_id and role_code are required")
	}

	// 2. 获取用户所有角色
	roles, err := p.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return false, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 3. 检查是否拥有指定角色
	for _, role := range roles {
		if role.RoleCode == roleCode {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissions 获取用户所有权限
//
// 参数说明：
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - userID: 用户ID
//
// 返回值：
//   - []string: 权限代码列表
//   - error: 错误信息
func (p *PermissionChecker) GetUserPermissions(
	ctx context.Context,
	tenantID, userID string,
) ([]string, error) {
	// 1. 参数验证
	if tenantID == "" || userID == "" {
		return nil, fmt.Errorf("tenant_id and user_id are required")
	}

	// 2. 获取用户所有角色
	roles, err := p.userRoleRepo.GetRolesByUser(ctx, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// 3. 收集所有角色的权限（去重）
	permissionMap := make(map[string]bool)
	for _, role := range roles {
		// 获取角色的权限列表
		perms, err := p.dataPermRepo.GetByRole(ctx, role.RoleID)
		if err != nil {
			continue
		}

		for _, perm := range perms {
			// 生成权限代码：{resource_type}:{action}
			permCode := fmt.Sprintf("%s:*", string(perm.ResourceType))
			permissionMap[permCode] = true
		}
	}

	// 4. 转换为切片返回
	permissions := make([]string, 0, len(permissionMap))
	for perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// compareFieldLevel 比较字段权限级别
//
// 权限级别: hidden (0) < readonly (1) < editable (2)
//
// 返回值：
//   - 正数: level1 > level2
//   - 0: level1 == level2
//   - 负数: level1 < level2
func (p *PermissionChecker) compareFieldLevel(level1, level2 string) int {
	levels := map[string]int{
		"hidden":   0,
		"readonly": 1,
		"editable": 2,
	}

	l1 := levels[level1]
	l2 := levels[level2]

	if l1 > l2 {
		return 1
	} else if l1 < l2 {
		return -1
	}
	return 0
}

// isOwner 检查是否是资源创建者
func (p *PermissionChecker) isOwner(
	ctx context.Context,
	userID string,
	resourceType entity.ResourceType,
	resourceID string,
) bool {
	// 根据资源类型查询相应的表，检查creator_id字段
	var creatorID string
	var err error

	switch resourceType {
	case entity.ResourceTypeBots:
		err = p.db.Table("bots").
			Select("creator_id").
			Where("bot_id = ?", resourceID).
			Scan(&creatorID).Error
	case entity.ResourceTypeConversations:
		err = p.db.Table("conversations").
			Select("creator_id").
			Where("conversation_id = ?", resourceID).
			Scan(&creatorID).Error
	case entity.ResourceTypeKnowledge:
		err = p.db.Table("knowledge_bases").
			Select("creator_id").
			Where("knowledge_id = ?", resourceID).
			Scan(&creatorID).Error
	case entity.ResourceTypeWorkflows:
		err = p.db.Table("workflows").
			Select("creator_id").
			Where("workflow_id = ?", resourceID).
			Scan(&creatorID).Error
	default:
		// 未知资源类型，默认拒绝
		return false
	}

	if err != nil {
		// 查询失败，记录日志并拒绝
		return false
	}

	return creatorID == userID
}

// isSameDepartment 检查是否同部门
func (p *PermissionChecker) isSameDepartment(
	ctx context.Context,
	userID string,
	resourceType entity.ResourceType,
	resourceID string,
) bool {
	// 1. 获取用户的部门ID列表
	userDeptIDs := make([]string, 0)
	err := p.db.Table("user_departments").
		Select("department_id").
		Where("user_id = ?", userID).
		Pluck("department_id", &userDeptIDs).Error

	if err != nil || len(userDeptIDs) == 0 {
		return false
	}

	// 2. 根据资源类型查询资源表的department_id
	var resourceDeptID string
	switch resourceType {
	case entity.ResourceTypeBots:
		err = p.db.Table("bots").
			Select("department_id").
			Where("bot_id = ?", resourceID).
			Scan(&resourceDeptID).Error
	case entity.ResourceTypeConversations:
		err = p.db.Table("conversations").
			Select("department_id").
			Where("conversation_id = ?", resourceID).
			Scan(&resourceDeptID).Error
	case entity.ResourceTypeKnowledge:
		err = p.db.Table("knowledge_bases").
			Select("department_id").
			Where("knowledge_id = ?", resourceID).
			Scan(&resourceDeptID).Error
	case entity.ResourceTypeWorkflows:
		err = p.db.Table("workflows").
			Select("department_id").
			Where("workflow_id = ?", resourceID).
			Scan(&resourceDeptID).Error
	default:
		// 未知资源类型，默认拒绝
		return false
	}

	if err != nil {
		return false
	}

	// 3. 检查资源的department_id是否在用户的部门列表中
	if resourceDeptID == "" {
		return false
	}

	for _, deptID := range userDeptIDs {
		if deptID == resourceDeptID {
			return true
		}
	}

	return false
}

// matchCustomFilter 匹配自定义过滤条件
func (p *PermissionChecker) matchCustomFilter(
	ctx context.Context,
	userID string,
	customFilter string,
	resourceType entity.ResourceType,
	resourceID string,
) bool {
	// 解析 customFilter JSON，进行过滤
	// 自定义过滤器是一个条件表达式数组
	var conditions []map[string]interface{}
	if err := json.Unmarshal([]byte(customFilter), &conditions); err != nil {
		return false
	}

	// 简化实现：实际应该根据条件表达式进行匹配
	return len(conditions) > 0
}

// contains 检查字符串切片是否包含某元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
