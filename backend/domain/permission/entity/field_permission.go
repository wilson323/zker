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

package entity

// FieldPermissionLevel 字段权限级别
type FieldPermissionLevel string

const (
	FieldPermissionLevelHidden   FieldPermissionLevel = "hidden"   // 隐藏
	FieldPermissionLevelReadonly FieldPermissionLevel = "readonly" // 只读
	FieldPermissionLevelEditable FieldPermissionLevel = "editable" // 可编辑
)

// FieldPermission 字段权限实体
type FieldPermission struct {
	PermissionID    string              `json:"permission_id" gorm:"primaryKey;type:varchar(36)"`
	RoleID          string              `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
	ResourceType    string              `json:"resource_type" gorm:"type:varchar(50);not null"`
	FieldName       string              `json:"field_name" gorm:"type:varchar(100);not null"`
	PermissionLevel FieldPermissionLevel `json:"permission_level" gorm:"type:enum('hidden','readonly','editable');not null"`
	CreatedAt       int64               `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt       int64               `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (FieldPermission) TableName() string {
	return "field_permissions"
}

// IsHidden 是否隐藏
func (fp *FieldPermission) IsHidden() bool {
	return fp.PermissionLevel == FieldPermissionLevelHidden
}

// IsReadonly 是否只读
func (fp *FieldPermission) IsReadonly() bool {
	return fp.PermissionLevel == FieldPermissionLevelReadonly
}

// IsEditable 是否可编辑
func (fp *FieldPermission) IsEditable() bool {
	return fp.PermissionLevel == FieldPermissionLevelEditable
}

// GetPriority 获取权限优先级（用于权限合并）
// editable > readonly > hidden
func (fp *FieldPermission) GetPriority() int {
	switch fp.PermissionLevel {
	case FieldPermissionLevelEditable:
		return 3
	case FieldPermissionLevelReadonly:
		return 2
	case FieldPermissionLevelHidden:
		return 1
	default:
		return 0
	}
}
