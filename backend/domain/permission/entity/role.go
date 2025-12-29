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

// RoleType 角色类型
type RoleType string

const (
	RoleTypeSystem RoleType = "system" // 系统预置角色
	RoleTypeCustom RoleType = "custom" // 用户自定义角色
)

// Role 角色实体
type Role struct {
	RoleID        string           `json:"role_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string           `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	RoleName      string           `json:"role_name" gorm:"type:varchar(100);not null"`
	RoleCode      string           `json:"role_code" gorm:"type:varchar(50);not null"`
	RoleType      RoleType         `json:"role_type" gorm:"type:enum('system','custom');not null;index:idx_role_type"`
	ParentRoleID  *string          `json:"parent_role_id,omitempty" gorm:"type:varchar(36)"`
	Description   string           `json:"description" gorm:"type:text"`
	CreatedAt     int64            `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt     int64            `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt     *int64           `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	ParentRole   *Role               `json:"parent_role,omitempty" gorm:"foreignKey:ParentRoleID"`
	DataPerms    []DataPermission    `json:"data_permissions,omitempty" gorm:"foreignKey:RoleID"`
	FieldPerms   []FieldPermission   `json:"field_permissions,omitempty" gorm:"foreignKey:RoleID"`
	UserRoles    []UserRole          `json:"user_roles,omitempty" gorm:"foreignKey:RoleID"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// IsSystemRole 是否是系统角色
func (r *Role) IsSystemRole() bool {
	return r.RoleType == RoleTypeSystem
}

// IsDeleted 是否已删除
func (r *Role) IsDeleted() bool {
	return r.DeletedAt != nil
}

// HasParent 是否有父角色
func (r *Role) HasParent() bool {
	return r.ParentRoleID != nil && *r.ParentRoleID != ""
}

// UserRole 用户角色关联实体
type UserRole struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    string `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID  string `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	RoleID    string `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
	CreatedAt int64  `json:"created_at" gorm:"not null;default:0"`

	// 关联
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
