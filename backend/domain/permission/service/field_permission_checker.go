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
	"regexp"
	"strings"

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// FieldPermission 3级字段权限（位标志）
type FieldPermission int

const (
	// VISIBLE 可见权限
	VISIBLE FieldPermission = 1 << iota
	// EDITABLE 可编辑权限
	EDITABLE
	// REQUIRED 必填权限
	REQUIRED
)

// HasVisible 是否可见
func (f FieldPermission) HasVisible() bool {
	return f&VISIBLE != 0
}

// HasEditable 是否可编辑
func (f FieldPermission) HasEditable() bool {
	return f&EDITABLE != 0
}

// HasRequired 是否必填
func (f FieldPermission) HasRequired() bool {
	return f&REQUIRED != 0
}

// String 返回权限的字符串表示
func (f FieldPermission) String() string {
	parts := []string{}
	if f.HasVisible() {
		parts = append(parts, "VISIBLE")
	}
	if f.HasEditable() {
		parts = append(parts, "EDITABLE")
	}
	if f.HasRequired() {
		parts = append(parts, "REQUIRED")
	}
	if len(parts) == 0 {
		return "NONE"
	}
	return strings.Join(parts, "|")
}

// FieldPermissionInfo 字段权限信息
type FieldPermissionInfo struct {
	FieldName       string            // 字段名
	Permission      FieldPermission   // 权限级别
	PermissionLevel string            // 权限级别字符串（hidden/readonly/editable）
	MaskPattern     string            // 脱敏规则（可选）
}

// FieldPermissionChecker 字段权限检查器接口
type FieldPermissionChecker interface {
	// CheckFieldPermission 检查字段权限
	CheckFieldPermission(ctx context.Context, userID, resourceType, field string) (FieldPermission, error)

	// FilterVisibleFields 过滤可见字段
	FilterVisibleFields(ctx context.Context, userID, resourceType string, fields []string) ([]string, error)

	// FilterEditableFields 过滤可编辑字段
	FilterEditableFields(ctx context.Context, userID, resourceType string, fields []string) ([]string, error)

	// MaskSensitiveFields 脱敏敏感字段
	MaskSensitiveFields(ctx context.Context, userID string, data map[string]interface{}) (map[string]interface{}, error)

	// GetFieldPermissions 获取所有字段权限
	GetFieldPermissions(ctx context.Context, userID, resourceType string) (map[string]*FieldPermissionInfo, error)

	// ValidateFieldPermissions 验证字段权限（用于写入操作）
	ValidateFieldPermissions(ctx context.Context, userID, resourceType string, data map[string]interface{}) error
}

// FieldPermissionCheckerImpl 字段权限检查器实现
type FieldPermissionCheckerImpl struct {
	userRoleRepo  repository.UserRoleRepository
	fieldPermRepo repository.FieldPermissionRepository
}

// NewFieldPermissionChecker 创建字段权限检查器
func NewFieldPermissionChecker(
	userRoleRepo repository.UserRoleRepository,
	fieldPermRepo repository.FieldPermissionRepository,
) FieldPermissionChecker {
	return &FieldPermissionCheckerImpl{
		userRoleRepo:  userRoleRepo,
		fieldPermRepo: fieldPermRepo,
	}
}

// CheckFieldPermission 检查字段权限
func (c *FieldPermissionCheckerImpl) CheckFieldPermission(
	ctx context.Context,
	userID, resourceType, field string,
) (FieldPermission, error) {
	// 获取用户的字段权限
	fieldPerms, err := c.GetFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return 0, err
	}

	// 如果字段没有配置权限，默认为可见+可编辑
	if permInfo, ok := fieldPerms[field]; ok {
		return permInfo.Permission, nil
	}

	// 默认权限：可见+可编辑
	return VISIBLE | EDITABLE, nil
}

// FilterVisibleFields 过滤可见字段
func (c *FieldPermissionCheckerImpl) FilterVisibleFields(
	ctx context.Context,
	userID, resourceType string,
	fields []string,
) ([]string, error) {
	visibleFields := make([]string, 0)

	// 获取字段权限
	fieldPerms, err := c.GetFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		perm := VISIBLE | EDITABLE // 默认可见

		if permInfo, ok := fieldPerms[field]; ok {
			perm = permInfo.Permission
		}

		if perm.HasVisible() {
			visibleFields = append(visibleFields, field)
		}
	}

	return visibleFields, nil
}

// FilterEditableFields 过滤可编辑字段
func (c *FieldPermissionCheckerImpl) FilterEditableFields(
	ctx context.Context,
	userID, resourceType string,
	fields []string,
) ([]string, error) {
	editableFields := make([]string, 0)

	// 获取字段权限
	fieldPerms, err := c.GetFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	for _, field := range fields {
		perm := VISIBLE | EDITABLE // 默认可编辑

		if permInfo, ok := fieldPerms[field]; ok {
			perm = permInfo.Permission
		}

		if perm.HasEditable() {
			editableFields = append(editableFields, field)
		}
	}

	return editableFields, nil
}

// MaskSensitiveFields 脱敏敏感字段
func (c *FieldPermissionCheckerImpl) MaskSensitiveFields(
	ctx context.Context,
	userID string,
	data map[string]interface{},
) (map[string]interface{}, error) {
	if data == nil {
		return nil, nil
	}

	// 推断资源类型（从数据中获取）
	resourceType := c.inferResourceType(data)
	if resourceType == "" {
		return data, nil
	}

	// 获取字段权限
	fieldPerms, err := c.GetFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return nil, err
	}

	// 复制数据避免修改原数据
	result := make(map[string]interface{})
	for key, value := range data {
		result[key] = value
	}

	// 脱敏处理
	for field, permInfo := range fieldPerms {
		if !permInfo.Permission.HasVisible() {
			// 不可见：删除字段
			delete(result, field)
		} else if !permInfo.Permission.HasEditable() && permInfo.MaskPattern != "" {
			// 可见但不可编辑且有脱敏规则：脱敏处理
			if value, ok := result[field]; ok {
				result[field] = c.maskValue(value, permInfo.MaskPattern)
			}
		}
	}

	return result, nil
}

// GetFieldPermissions 获取所有字段权限
func (c *FieldPermissionCheckerImpl) GetFieldPermissions(
	ctx context.Context,
	userID, resourceType string,
) (map[string]*FieldPermissionInfo, error) {
	// 1. 获取用户的所有角色
	roles, err := c.userRoleRepo.GetRolesByUser(ctx, userID, "")
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("reason", "failed to get user roles"),
            )
	}

	// 2. 合并所有角色的字段权限（取最大权限）
	fieldPerms := make(map[string]*FieldPermissionInfo)

	for _, role := range roles {
		perms, err := c.fieldPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType)
		if err != nil {
			continue
		}

		for _, perm := range perms {
			fieldName := perm.FieldName
			newPerm := c.convertFieldPermission(perm)

			// 如果已有权限，合并权限（取最大权限）
			if existing, ok := fieldPerms[fieldName]; ok {
				// 合并权限：位运算OR
				existing.Permission = existing.Permission | newPerm.Permission
				// 优先保留脱敏规则
				if newPerm.MaskPattern != "" {
					existing.MaskPattern = newPerm.MaskPattern
				}
			} else {
				fieldPerms[fieldName] = newPerm
			}
		}
	}

	return fieldPerms, nil
}

// ValidateFieldPermissions 验证字段权限（用于写入操作）
func (c *FieldPermissionCheckerImpl) ValidateFieldPermissions(
	ctx context.Context,
	userID, resourceType string,
	data map[string]interface{},
) error {
	if data == nil || len(data) == 0 {
		return nil
	}

	// 获取字段权限
	fieldPerms, err := c.GetFieldPermissions(ctx, userID, resourceType)
	if err != nil {
		return err
	}

	// 检查每个字段
	for field := range data {
		perm := VISIBLE | EDITABLE // 默认可编辑

		if permInfo, ok := fieldPerms[field]; ok {
			perm = permInfo.Permission
		}

		// 检查是否可编辑
		if !perm.HasEditable() {
			return &FieldPermissionDeniedError{
				UserID:       userID,
				ResourceType: resourceType,
				FieldName:    field,
				Reason:       "field is not editable",
			}
		}

		// 检查是否必填
		if perm.HasRequired() {
			if value, ok := data[field]; ok && isEmpty(value) {
				return &FieldPermissionDeniedError{
					UserID:       userID,
					ResourceType: resourceType,
					FieldName:    field,
					Reason:       "required field is empty",
				}
			}
		}
	}

	return nil
}

// convertFieldPermission 转换字段权限实体到权限信息
func (c *FieldPermissionCheckerImpl) convertFieldPermission(perm *entity.FieldPermission) *FieldPermissionInfo {
	var fp FieldPermission

	// 根据PermissionLevel转换为位标志
	switch perm.PermissionLevel {
	case entity.FieldPermissionLevelHidden:
		fp = 0 // 不可见
	case entity.FieldPermissionLevelReadonly:
		fp = VISIBLE // 仅可见
	case entity.FieldPermissionLevelEditable:
		fp = VISIBLE | EDITABLE // 可见+可编辑
	}

	return &FieldPermissionInfo{
		FieldName:       perm.FieldName,
		Permission:      fp,
		PermissionLevel: string(perm.PermissionLevel),
		MaskPattern:     "", // 可从扩展字段获取
	}
}

// inferResourceType 推断资源类型
func (c *FieldPermissionCheckerImpl) inferResourceType(data map[string]interface{}) string {
	// 根据数据中的特定字段推断资源类型
	if _, ok := data["bot_id"]; ok {
		return "bots"
	}
	if _, ok := data["conversation_id"]; ok {
		return "conversations"
	}
	if _, ok := data["knowledge_id"]; ok {
		return "knowledge"
	}
	if _, ok := data["workflow_id"]; ok {
		return "workflows"
	}
	return ""
}

// maskValue 脱敏处理
func (c *FieldPermissionCheckerImpl) maskValue(value interface{}, pattern string) interface{} {
	if value == nil {
		return nil
	}

	strValue := fmt.Sprintf("%v", value)

	// 支持多种脱敏规则
	switch pattern {
	case "phone":
		// 手机号脱敏：138****1234
		if len(strValue) == 11 {
			return strValue[:3] + "****" + strValue[7:]
		}
	case "email":
		// 邮箱脱敏：u***@example.com
		parts := strings.Split(strValue, "@")
		if len(parts) == 2 && len(parts[0]) > 1 {
			return parts[0][:1] + "***@" + parts[1]
		}
	case "idcard":
		// 身份证脱敏：110101********1234
		if len(strValue) == 18 {
			return strValue[:6] + "********" + strValue[14:]
		}
	case "name":
		// 姓名脱敏：张**
		if len(strValue) >= 2 {
			runes := []rune(strValue)
			if len(runes) == 2 {
				return string(runes[:1]) + "*"
			} else if len(runes) > 2 {
				return string(runes[:1]) + "**" + string(runes[len(runes)-1:])
			}
		}
	case "card":
		// 银行卡脱敏：6222***********1234
		if len(strValue) >= 16 {
			return strValue[:4] + strings.Repeat("*", len(strValue)-8) + strValue[len(strValue)-4:]
		}
	default:
		// 自定义正则脱敏
		if strings.HasPrefix(pattern, "regex:") {
			regexPattern := strings.TrimPrefix(pattern, "regex:")
			re, err := regexp.Compile(regexPattern)
			if err == nil {
				return re.ReplaceAllString(strValue, "***")
			}
		}
		// 默认完全脱敏
		if len(strValue) > 0 {
			return strings.Repeat("*", len(strValue))
		}
	}

	return value
}

// FieldPermissionDeniedError 字段权限拒绝错误
type FieldPermissionDeniedError struct {
	UserID       string
	ResourceType string
	FieldName    string
	Reason       string
}

func (e *FieldPermissionDeniedError) Error() string {
	return fmt.Sprintf("field permission denied for user %s on %s.%s: %s",
		e.UserID, e.ResourceType, e.FieldName, e.Reason)
}

// isEmpty 检查值是否为空
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	case int, int64, float64:
		return v == 0 || v == 0.0
	case bool:
		return !v
	default:
		return false
	}
}
