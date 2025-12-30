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

import "time"

// UserRole 用户在租户中的角色
type UserRole string

const (
	UserRoleOwner   UserRole = "owner"   // 租户拥有者
	UserRoleAdmin   UserRole = "admin"   // 租户管理员
	UserRoleMember  UserRole = "member"  // 普通成员
)

// UserTenantRelation 用户-租户关联实体
type UserTenantRelation struct {
	UserID    int64     `json:"user_id"`
	TenantID  string    `json:"tenant_id"`
	Role      UserRole  `json:"role"`
	IsDefault bool      `json:"is_default"`
}

// UserTenant 用户-租户关联表实体（完整）
type UserTenant struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	TenantID  string     `json:"tenant_id"`
	Role      UserRole   `json:"role"`
	IsDefault bool       `json:"is_default"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (UserTenant) TableName() string {
	return "user_tenant"
}

// IsOwner 是否为租户拥有者
func (ut *UserTenant) IsOwner() bool {
	return ut.Role == UserRoleOwner
}

// IsAdmin 是否为租户管理员
func (ut *UserTenant) IsAdmin() bool {
	return ut.Role == UserRoleAdmin
}

// CanManage 是否可以管理租户
func (ut *UserTenant) CanManage() bool {
	return ut.Role == UserRoleOwner || ut.Role == UserRoleAdmin
}

// IsActive 是否激活（未删除）
func (ut *UserTenant) IsActive() bool {
	return ut.DeletedAt == nil
}
