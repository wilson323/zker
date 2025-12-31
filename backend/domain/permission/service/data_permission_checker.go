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
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// DataPermissionLevel 5级数据权限级别
type DataPermissionLevel int

const (
	// ALL 全部数据权限
	ALL DataPermissionLevel = iota
	// DEPARTMENT_AND_SUB 本部门及子部门数据
	DEPARTMENT_AND_SUB
	// DEPARTMENT 本部门数据
	DEPARTMENT
	// SELF 仅本人数据
	SELF
	// CUSTOM 自定义过滤
	CUSTOM
)

// String 返回权限级别的字符串表示
func (l DataPermissionLevel) String() string {
	switch l {
	case ALL:
		return "ALL"
	case DEPARTMENT_AND_SUB:
		return "DEPARTMENT_AND_SUB"
	case DEPARTMENT:
		return "DEPARTMENT"
	case SELF:
		return "SELF"
	case CUSTOM:
		return "CUSTOM"
	default:
		return "UNKNOWN"
	}
}

// DataPermissionFilter 数据权限过滤条件
type DataPermissionFilter struct {
	WhereClause string        // WHERE子句
	Args        []interface{} // 参数
}

// DataPermissionChecker 数据权限检查器接口
type DataPermissionChecker interface {
	// Filter 生成数据权限过滤条件
	Filter(ctx context.Context, userID string, level DataPermissionLevel, resourceType string) (*DataPermissionFilter, error)

	// CheckAccess 检查对指定资源的访问权限
	CheckAccess(ctx context.Context, userID, resourceID string, level DataPermissionLevel, resourceType string) (bool, error)

	// GetAccessibleResourceIDs 获取可访问的资源ID列表
	GetAccessibleResourceIDs(ctx context.Context, userID string, level DataPermissionLevel, resourceType string) ([]string, error)

	// GetAccessibleDepartmentIDs 获取可访问的部门ID列表
	GetAccessibleDepartmentIDs(ctx context.Context, userID string, level DataPermissionLevel) ([]string, error)
}

// DataPermissionCheckerImpl 数据权限检查器实现
type DataPermissionCheckerImpl struct {
	db             *gorm.DB
	userDeptRepo   repository.UserDepartmentRepository
	departmentRepo repository.DepartmentRepository
	userRoleRepo   repository.UserRoleRepository
	dataPermRepo   repository.DataPermissionRepository
}

// NewDataPermissionChecker 创建数据权限检查器
func NewDataPermissionChecker(
	db *gorm.DB,
	userDeptRepo repository.UserDepartmentRepository,
	departmentRepo repository.DepartmentRepository,
	userRoleRepo repository.UserRoleRepository,
	dataPermRepo repository.DataPermissionRepository,
) DataPermissionChecker {
	return &DataPermissionCheckerImpl{
		db:             db,
		userDeptRepo:   userDeptRepo,
		departmentRepo: departmentRepo,
		userRoleRepo:   userRoleRepo,
		dataPermRepo:   dataPermRepo,
	}
}

// Filter 生成数据权限过滤条件
func (c *DataPermissionCheckerImpl) Filter(
	ctx context.Context,
	userID string,
	level DataPermissionLevel,
	resourceType string,
) (*DataPermissionFilter, error) {
	switch level {
	case ALL:
		// 全部数据：不添加过滤条件
		return &DataPermissionFilter{
			WhereClause: "",
			Args:        nil,
		}, nil

	case DEPARTMENT_AND_SUB:
		// 本部门及子部门
		deptIDs, err := c.GetAccessibleDepartmentIDs(ctx, userID, level)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get department IDs"),
            )
		}

		if len(deptIDs) == 0 {
			// 无部门，返回1=0使查询结果为空
			return &DataPermissionFilter{
				WhereClause: "1 = 0",
				Args:        nil,
			}, nil
		}

		// 生成IN查询
		placeholders := make([]string, len(deptIDs))
		args := make([]interface{}, len(deptIDs))
		for i, deptID := range deptIDs {
			placeholders[i] = "?"
			args[i] = deptID
		}

		whereClause := fmt.Sprintf("department_id IN (%s)", joinStrings(placeholders, ","))
		return &DataPermissionFilter{
			WhereClause: whereClause,
			Args:        args,
		}, nil

	case DEPARTMENT:
		// 仅本部门
		userDepts, err := c.userDeptRepo.GetByUser(ctx, userID, "")
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get user departments"),
            )
		}

		if len(userDepts) == 0 {
			// 无部门，返回1=0使查询结果为空
			return &DataPermissionFilter{
				WhereClause: "1 = 0",
				Args:        nil,
			}, nil
		}

		deptIDs := make([]string, len(userDepts))
		for i, ud := range userDepts {
			deptIDs[i] = ud.DepartmentID
		}

		// 生成IN查询
		placeholders := make([]string, len(deptIDs))
		args := make([]interface{}, len(deptIDs))
		for i, deptID := range deptIDs {
			placeholders[i] = "?"
			args[i] = deptID
		}

		whereClause := fmt.Sprintf("department_id IN (%s)", joinStrings(placeholders, ","))
		return &DataPermissionFilter{
			WhereClause: whereClause,
			Args:        args,
		}, nil

	case SELF:
		// 仅本人数据
		return &DataPermissionFilter{
			WhereClause: "creator_id = ?",
			Args:        []interface{}{userID},
		}, nil

	case CUSTOM:
		// 自定义过滤：从角色权限中获取custom_filter
		roles, err := c.userRoleRepo.GetRolesByUser(ctx, userID, "")
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get user roles"),
            )
		}

		for _, role := range roles {
			perm, err := c.dataPermRepo.GetByRoleAndResource(ctx, role.RoleID, entity.ResourceType(resourceType))
			if err != nil {
				continue
			}

			if perm != nil && perm.CustomFilter != "" {
				// 解析custom_filter JSON
				var customFilter map[string]interface{}
				if err := json.Unmarshal([]byte(perm.CustomFilter), &customFilter); err != nil {
					continue
				}

				// 根据custom_filter生成WHERE条件
				return c.buildCustomFilter(customFilter)
			}
		}

		// 未找到自定义过滤器，默认返回1=0
		return &DataPermissionFilter{
			WhereClause: "1 = 0",
			Args:        nil,
		}, nil

	default:
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "unknown permission level"),
                errorx.KV("level", fmt.Sprintf("%d", level)),
            )
	}
}

// CheckAccess 检查对指定资源的访问权限
func (c *DataPermissionCheckerImpl) CheckAccess(
	ctx context.Context,
	userID, resourceID string,
	level DataPermissionLevel,
	resourceType string,
) (bool, error) {
	switch level {
	case ALL:
		// 全部数据权限：直接返回true
		return true, nil

	case DEPARTMENT_AND_SUB:
		// 检查资源是否属于用户的部门或子部门
		return c.checkDepartmentAccess(ctx, userID, resourceID, resourceType, true)

	case DEPARTMENT:
		// 检查资源是否属于用户的部门
		return c.checkDepartmentAccess(ctx, userID, resourceID, resourceType, false)

	case SELF:
		// 检查资源是否由用户创建
		return c.checkOwnership(ctx, userID, resourceID, resourceType)

	case CUSTOM:
		// 根据自定义过滤器检查
		return c.checkCustomFilter(ctx, userID, resourceID, resourceType)

	default:
		return false, errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "unknown permission level"),
                errorx.KV("level", fmt.Sprintf("%d", level)),
            )
	}
}

// GetAccessibleResourceIDs 获取可访问的资源ID列表
func (c *DataPermissionCheckerImpl) GetAccessibleResourceIDs(
	ctx context.Context,
	userID string,
	level DataPermissionLevel,
	resourceType string,
) ([]string, error) {
	// 生成过滤条件
	filter, err := c.Filter(ctx, userID, level, resourceType)
	if err != nil {
		return nil, err
	}

	// 构建查询
	query := c.db.Table(resourceType).Select("id")
	if filter.WhereClause != "" {
		query = query.Where(filter.WhereClause, filter.Args...)
	}

	// 执行查询
	var resourceIDs []string
	if err := query.Pluck("id", &resourceIDs).Error; err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to query resource IDs"),
            )
	}

	return resourceIDs, nil
}

// GetAccessibleDepartmentIDs 获取可访问的部门ID列表
func (c *DataPermissionCheckerImpl) GetAccessibleDepartmentIDs(
	ctx context.Context,
	userID string,
	level DataPermissionLevel,
) ([]string, error) {
	// 获取用户的部门
	userDepts, err := c.userDeptRepo.GetByUser(ctx, userID, "")
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get user departments"),
            )
	}

	deptIDs := make([]string, 0)

	for _, ud := range userDepts {
		deptIDs = append(deptIDs, ud.DepartmentID)

		// 如果是DEPARTMENT_AND_SUB级别或者是部门领导，获取子部门
		if level == DEPARTMENT_AND_SUB || ud.IsLeader {
			childDepts, err := c.departmentRepo.GetDescendants(ctx, ud.DepartmentID)
			if err != nil {
				continue
			}

			for _, child := range childDepts {
				deptIDs = append(deptIDs, child.DepartmentID)
			}
		}
	}

	return deptIDs, nil
}

// checkDepartmentAccess 检查部门访问权限
func (c *DataPermissionCheckerImpl) checkDepartmentAccess(
	ctx context.Context,
	userID, resourceID, resourceType string,
	includeSubDept bool,
) (bool, error) {
	// 1. 获取用户的部门ID列表
	deptIDs, err := c.userDeptRepo.GetByUser(ctx, userID, "")
	if err != nil {
		return false, errorx.WrapByCode(err, errno.ErrPermissionInvalidParamCode,
                errorx.KV("operation", "get user departments"),
            )
	}

	if len(deptIDs) == 0 {
		return false, nil
	}

	// 2. 获取资源的部门ID
	var resourceDeptID string
	err = c.db.Table(resourceType).
		Select("department_id").
		Where("id = ?", resourceID).
		Scan(&resourceDeptID).Error

	if err != nil {
		return false, errorx.WrapByCode(err, errno.ErrPermissionInvalidParamCode,
                errorx.KV("operation", "get resource department"),
            )
	}

	if resourceDeptID == "" {
		return false, nil
	}

	// 3. 检查资源的部门ID是否在用户的部门列表中
	for _, ud := range deptIDs {
		if ud.DepartmentID == resourceDeptID {
			return true, nil
		}

		// 如果需要包含子部门，检查是否是子部门
		if includeSubDept {
			childDepts, err := c.departmentRepo.GetDescendants(ctx, ud.DepartmentID)
			if err != nil {
				continue
			}

			for _, child := range childDepts {
				if child.DepartmentID == resourceDeptID {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// checkOwnership 检查所有权
func (c *DataPermissionCheckerImpl) checkOwnership(
	ctx context.Context,
	userID, resourceID, resourceType string,
) (bool, error) {
	var creatorID string
	err := c.db.Table(resourceType).
		Select("creator_id").
		Where("id = ?", resourceID).
		Scan(&creatorID).Error

	if err != nil {
		return false, errorx.WrapByCode(err, errno.ErrPermissionInvalidParamCode,
                errorx.KV("operation", "check ownership"),
            )
	}

	return creatorID == userID, nil
}

// checkCustomFilter 检查自定义过滤器
func (c *DataPermissionCheckerImpl) checkCustomFilter(
	ctx context.Context,
	userID, resourceID, resourceType string,
) (bool, error) {
	// 获取自定义过滤器
	roles, err := c.userRoleRepo.GetRolesByUser(ctx, userID, "")
	if err != nil {
		return false, errorx.WrapByCode(err, errno.ErrPermissionInvalidParamCode,
                errorx.KV("operation", "get user roles"),
            )
	}

	for _, role := range roles {
		perm, err := c.dataPermRepo.GetByRoleAndResource(ctx, role.RoleID, entity.ResourceType(resourceType))
		if err != nil {
			continue
		}

		if perm != nil && perm.CustomFilter != "" {
			// 解析custom_filter JSON
			var customFilter map[string]interface{}
			if err := json.Unmarshal([]byte(perm.CustomFilter), &customFilter); err != nil {
				continue
			}

			// 根据custom_filter检查资源
			return c.matchCustomFilter(ctx, resourceID, resourceType, customFilter)
		}
	}

	return false, nil
}

// buildCustomFilter 根据自定义过滤器构建WHERE条件
func (c *DataPermissionCheckerImpl) buildCustomFilter(customFilter map[string]interface{}) (*DataPermissionFilter, error) {
	// custom_filter格式示例：
	// {
	//   "conditions": [
	//     {"field": "status", "operator": "=", "value": "published"},
	//     {"field": "created_at", "operator": ">", "value": "1234567890"}
	//   ],
	//   "logic": "AND"
	// }

	conditions, ok := customFilter["conditions"].([]interface{})
	if !ok || len(conditions) == 0 {
		return &DataPermissionFilter{
			WhereClause: "1 = 0",
			Args:        nil,
		}, nil
	}

	logic := "AND"
	if logicVal, ok := customFilter["logic"].(string); ok {
		logic = logicVal
	}

	whereParts := make([]string, 0)
	args := make([]interface{}, 0)

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		field, _ := condMap["field"].(string)
		operator, _ := condMap["operator"].(string)
		value := condMap["value"]

		whereParts = append(whereParts, fmt.Sprintf("%s %s ?", field, operator))
		args = append(args, value)
	}

	if len(whereParts) == 0 {
		return &DataPermissionFilter{
			WhereClause: "1 = 0",
			Args:        nil,
		}, nil
	}

	whereClause := joinStrings(whereParts, " "+logic+" ")
	return &DataPermissionFilter{
		WhereClause: whereClause,
		Args:        args,
	}, nil
}

// matchCustomFilter 匹配自定义过滤器
func (c *DataPermissionCheckerImpl) matchCustomFilter(
	ctx context.Context,
	resourceID, resourceType string,
	customFilter map[string]interface{},
) (bool, error) {
	// 获取资源数据
	var resourceData map[string]interface{}
	err := c.db.Table(resourceType).
		Where("id = ?", resourceID).
		First(&resourceData).Error

	if err != nil {
		return false, err
	}

	// 解析条件
	conditions, ok := customFilter["conditions"].([]interface{})
	if !ok || len(conditions) == 0 {
		return false, nil
	}

	logic := "AND"
	if logicVal, ok := customFilter["logic"].(string); ok {
		logic = logicVal
	}

	results := make([]bool, 0)

	for _, cond := range conditions {
		condMap, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}

		field, _ := condMap["field"].(string)
		operator, _ := condMap["operator"].(string)
		expectedValue := condMap["value"]

		// 获取资源字段值
		actualValue, ok := resourceData[field]
		if !ok {
			results = append(results, false)
			continue
		}

		// 比较值
		result := c.compareValues(actualValue, operator, expectedValue)
		results = append(results, result)
	}

	// 根据逻辑运算符组合结果
	if logic == "AND" {
		for _, r := range results {
			if !r {
				return false, nil
			}
		}
		return true, nil
	} else { // OR
		for _, r := range results {
			if r {
				return true, nil
			}
		}
		return false, nil
	}
}

// compareValues 比较值
func (c *DataPermissionCheckerImpl) compareValues(actual interface{}, operator string, expected interface{}) bool {
	switch operator {
	case "=":
		return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected)
	case "!=":
		return fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected)
	case ">":
		return c.compareNumeric(actual, expected, ">")
	case ">=":
		return c.compareNumeric(actual, expected, ">=")
	case "<":
		return c.compareNumeric(actual, expected, "<")
	case "<=":
		return c.compareNumeric(actual, expected, "<=")
	case "IN":
		return c.compareIn(actual, expected)
	case "NOT IN":
		return !c.compareIn(actual, expected)
	default:
		return false
	}
}

// compareNumeric 比较数值
func (c *DataPermissionCheckerImpl) compareNumeric(actual, expected interface{}, operator string) bool {
	// 简化实现：转换为float64比较
	var actualFloat, expectedFloat float64

	switch v := actual.(type) {
	case int:
		actualFloat = float64(v)
	case int64:
		actualFloat = float64(v)
	case float64:
		actualFloat = v
	case float32:
		actualFloat = float64(v)
	case string:
		// 尝试解析时间戳
		if timestamp, err := time.Parse(time.RFC3339, v); err == nil {
			actualFloat = float64(timestamp.Unix())
		}
	default:
		return false
	}

	switch v := expected.(type) {
	case int:
		expectedFloat = float64(v)
	case int64:
		expectedFloat = float64(v)
	case float64:
		expectedFloat = v
	case float32:
		expectedFloat = float64(v)
	case string:
		if timestamp, err := time.Parse(time.RFC3339, v); err == nil {
			expectedFloat = float64(timestamp.Unix())
		}
	default:
		return false
	}

	switch operator {
	case ">":
		return actualFloat > expectedFloat
	case ">=":
		return actualFloat >= expectedFloat
	case "<":
		return actualFloat < expectedFloat
	case "<=":
		return actualFloat <= expectedFloat
	default:
		return false
	}
}

// compareIn 比较IN操作
func (c *DataPermissionCheckerImpl) compareIn(actual, expected interface{}) bool {
	expectedSlice, ok := expected.([]interface{})
	if !ok {
		return false
	}

	for _, item := range expectedSlice {
		if fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", item) {
			return true
		}
	}

	return false
}

// joinStrings 连接字符串切片
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}

	return result
}
