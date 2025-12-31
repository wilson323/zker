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

// UserTenantRelation 用户-租户关联关系
// 用于支持用户多租户场景
type UserTenantRelation struct {
	UserID    int64     `json:"user_id" gorm:"column:user_id;primaryKey"`
	TenantID  string    `json:"tenant_id" gorm:"column:tenant_id;primaryKey"`
	Role      string    `json:"role" gorm:"column:role;not null"` // 用户在租户中的角色
	IsDefault bool      `json:"is_default" gorm:"column:is_default;default:false"` // 是否为默认租户

	// 时间戳
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at;index"`
}

// TableName 指定表名
func (UserTenantRelation) TableName() string {
	return "user_tenant"
}
