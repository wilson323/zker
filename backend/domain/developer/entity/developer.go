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

// DeveloperStatus 开发者状态
type DeveloperStatus string

const (
	DeveloperStatusActive    DeveloperStatus = "active"    // 激活
	DeveloperStatusSuspended DeveloperStatus = "suspended" // 暂停
	DeveloperStatusDeleted   DeveloperStatus = "deleted"   // 已删除
)

// Developer 开发者实体
type Developer struct {
	DeveloperID    string          `json:"developer_id" gorm:"primaryKey;type:varchar(36);comment:开发者ID"`
	TenantID       string          `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id;comment:租户ID"`
	UserID         string          `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id;comment:用户ID"`
	DeveloperName  string          `json:"developer_name" gorm:"type:varchar(100);not null;comment:开发者名称"`
	DeveloperEmail string          `json:"developer_email" gorm:"type:varchar(100);not null;comment:开发者邮箱"`
	Status         DeveloperStatus `json:"status" gorm:"type:enum('active','suspended','deleted');default:'active';not null;index:idx_status;comment:状态"`
	CreatedAt      int64           `json:"created_at" gorm:"not null;default:0;comment:创建时间"`
	UpdatedAt      int64           `json:"updated_at" gorm:"not null;default:0;comment:更新时间"`
	DeletedAt      *int64          `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`

	// 关联
	Projects []Project `json:"projects,omitempty" gorm:"foreignKey:DeveloperID;references:DeveloperID"`
}

// TableName 指定表名
func (Developer) TableName() string {
	return "developers"
}

// IsActive 是否激活
func (d *Developer) IsActive() bool {
	return d.Status == DeveloperStatusActive
}

// IsDeleted 是否已删除
func (d *Developer) IsDeleted() bool {
	return d.Status == DeveloperStatusDeleted || d.DeletedAt != nil
}

// GetCreatedAtAsTime 获取创建时间
func (d *Developer) GetCreatedAtAsTime() time.Time {
	return time.Unix(d.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (d *Developer) GetUpdatedAtAsTime() time.Time {
	return time.Unix(d.UpdatedAt/1000, 0)
}
