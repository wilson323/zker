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
	"database/sql/driver"
	"encoding/json"
	"time"
)

// PermissionType 权限类型
type PermissionType string

const (
	PermissionTypeRole            PermissionType = "role"             // 角色权限
	PermissionTypeDataPermission  PermissionType = "data_permission"  // 数据权限
	PermissionTypeFieldPermission PermissionType = "field_permission" // 字段权限
)

// TemporaryGrant 临时授权实体
type TemporaryGrant struct {
	GrantID int64 `json:"grant_id" gorm:"primaryKey;autoIncrement"`

	// 授权码信息
	GrantCode string `json:"grant_code" gorm:"uniqueIndex:uk_grant_code;size:100;not null"`

	// 授权信息
	GranteeID string `json:"grantee_id" gorm:"size:36;not null;index:idx_grantee_id"` // 被授权人用户ID
	GrantorID string `json:"grantor_id" gorm:"size:36;not null;index:idx_grantor_id"` // 授权人用户ID
	TenantID  string `json:"tenant_id" gorm:"size:36;not null;index:idx_tenant_id"`   // 租户ID

	// 权限类型和数据
	PermissionType  PermissionType `json:"permission_type" gorm:"type:enum('role','data_permission','field_permission');not null"`
	PermissionData  PermissionData `json:"permission_data" gorm:"type:json;not null"` // JSON数据

	// 有效期
	ExpiresAt int64 `json:"expires_at" gorm:"not null;index:idx_expires_at"` // 过期时间（毫秒时间戳）

	// 状态
	IsUsed    bool   `json:"is_used" gorm:"default:false;index:idx_is_used"`
	IsRevoked bool   `json:"is_revoked" gorm:"default:false;index:idx_is_revoked"`
	UsedAt    *int64 `json:"used_at,omitempty"`    // 使用时间（毫秒时间戳）
	RevokedAt *int64 `json:"revoked_at,omitempty"` // 撤销时间（毫秒时间戳）

	// 元数据
	Reason    string `json:"reason,omitempty" gorm:"size:500"`   // 授权/撤销原因
	RequestID string `json:"request_id,omitempty" gorm:"size:36"` // 关联请求ID

	// 审计字段
	CreatedAt int64 `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt int64 `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (TemporaryGrant) TableName() string {
	return "temporary_grants"
}

// PermissionData 权限数据（支持JSON序列化）
type PermissionData map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (pd *PermissionData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, pd)
}

// Value 实现 driver.Valuer 接口
func (pd PermissionData) Value() (driver.Value, error) {
	if len(pd) == 0 {
		return nil, nil
	}

	return json.Marshal(pd)
}

// IsExpired 检查授权是否已过期
func (tg *TemporaryGrant) IsExpired() bool {
	return time.Now().UnixMilli() > tg.ExpiresAt
}

// IsValid 检查授权是否有效（未过期、未使用、未撤销）
func (tg *TemporaryGrant) IsValid() bool {
	if tg.IsExpired() {
		return false
	}
	if tg.IsUsed {
		return false
	}
	if tg.IsRevoked {
		return false
	}
	return true
}

// GetRemainingTime 获取剩余有效时间（毫秒）
// 返回值: 0表示已过期，>0表示剩余毫秒数
func (tg *TemporaryGrant) GetRemainingTime() int64 {
	now := time.Now().UnixMilli()
	remaining := tg.ExpiresAt - now

	if remaining < 0 {
		return 0
	}

	return remaining
}

// GetExpiresAtTime 获取过期时间的time.Time对象
func (tg *TemporaryGrant) GetExpiresAtTime() time.Time {
	return time.UnixMilli(tg.ExpiresAt)
}

// GetUsedAtTime 获取使用时间的time.Time对象
func (tg *TemporaryGrant) GetUsedAtTime() *time.Time {
	if tg.UsedAt == nil {
		return nil
	}

	t := time.UnixMilli(*tg.UsedAt)
	return &t
}

// GetRevokedAtTime 获取撤销时间的time.Time对象
func (tg *TemporaryGrant) GetRevokedAtTime() *time.Time {
	if tg.RevokedAt == nil {
		return nil
	}

	t := time.UnixMilli(*tg.RevokedAt)
	return &t
}

// MarkAsUsed 标记为已使用
func (tg *TemporaryGrant) MarkAsUsed() {
	now := time.Now().UnixMilli()
	tg.IsUsed = true
	tg.UsedAt = &now
	tg.UpdatedAt = now
}

// MarkAsRevoked 标记为已撤销
func (tg *TemporaryGrant) MarkAsRevoked(reason string) {
	now := time.Now().UnixMilli()
	tg.IsRevoked = true
	tg.RevokedAt = &now
	tg.Reason = reason
	tg.UpdatedAt = now
}
