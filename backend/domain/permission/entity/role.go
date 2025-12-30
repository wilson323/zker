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

import (
	"time"
)

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
	Priority      int              `json:"priority" gorm:"not null;default:0;index:idx_priority"` // 角色优先级，数值越大优先级越高
	IsSystem      bool             `json:"is_system" gorm:"not null;default:false"`               // 是否是系统角色（不可删除）
	CreatedAt     int64            `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt     int64            `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt     *int64           `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	ParentRole   *Role                 `json:"parent_role,omitempty" gorm:"foreignKey:ParentRoleID"`
	DataPerms    []*DataPermission     `json:"data_permissions,omitempty" gorm:"foreignKey:RoleID"`
	FieldPerms   []*FieldPermission    `json:"field_permissions,omitempty" gorm:"foreignKey:RoleID"`
	UserRoles    []*UserRole           `json:"user_roles,omitempty" gorm:"foreignKey:RoleID"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// IsSystemRole 是否是系统角色
func (r *Role) IsSystemRole() bool {
	return r.RoleType == RoleTypeSystem || r.IsSystem
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
	ID            int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        string  `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID      string  `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	RoleID        string  `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
	GrantedBy     string  `json:"granted_by" gorm:"type:varchar(36);not null;default:'system'"`              // 授权人ID
	GrantedReason string  `json:"granted_reason" gorm:"type:varchar(500)"`                                   // 授权原因
	ExpiresAt     *int64  `json:"expires_at,omitempty" gorm:"type:bigint unsigned;index:idx_expires_at"`     // 过期时间（毫秒时间戳，NULL表示永久）
	CreatedAt     int64   `json:"created_at" gorm:"not null;default:0"`

	// 关联
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// IsExpired 检查角色是否已过期
func (ur *UserRole) IsExpired() bool {
	if ur.ExpiresAt == nil {
		return false // 永久角色
	}

	now := time.Now().UnixMilli()
	return *ur.ExpiresAt < now
}

// GetRemainingTime 获取剩余有效时间（毫秒）
// 返回值: -1表示永久有效，0表示已过期，>0表示剩余毫秒数
func (ur *UserRole) GetRemainingTime() int64 {
	if ur.ExpiresAt == nil {
		return -1 // 永久有效
	}

	now := time.Now().UnixMilli()
	remaining := *ur.ExpiresAt - now

	if remaining < 0 {
		return 0
	}

	return remaining
}

// IsPermanent 检查是否为永久角色
func (ur *UserRole) IsPermanent() bool {
	return ur.ExpiresAt == nil
}

// GetExpiresAtTime 获取过期时间的time.Time对象
func (ur *UserRole) GetExpiresAtTime() *time.Time {
	if ur.ExpiresAt == nil {
		return nil
	}

	t := time.UnixMilli(*ur.ExpiresAt)
	return &t
}
