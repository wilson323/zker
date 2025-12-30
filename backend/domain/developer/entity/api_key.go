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
	"errors"
	"time"
)

// APIKeyStatus API密钥状态
type APIKeyStatus string

const (
	APIKeyStatusActive  APIKeyStatus = "active"  // 激活
	APIKeyStatusRevoked APIKeyStatus = "revoked" // 已撤销
	APIKeyStatusExpired APIKeyStatus = "expired" // 已过期
)

// APIKeyScope API密钥权限范围
type APIKeyScope string

const (
	APIKeyScopeRead  APIKeyScope = "read"  // 只读
	APIKeyScopeWrite APIKeyScope = "write" // 读写
	APIKeyScopeAdmin APIKeyScope = "admin" // 管理员
)

// APIKeyScopes 自定义类型，用于JSON处理
type APIKeyScopes []APIKeyScope

// Value 实现 driver.Valuer 接口
func (s APIKeyScopes) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 实现 sql.Scanner 接口
func (s *APIKeyScopes) Scan(value interface{}) error {
	if value == nil {
		*s = APIKeyScopes{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan into APIKeyScopes")
	}

	return json.Unmarshal(bytes, s)
}

// Contains 检查是否包含指定权限
func (s APIKeyScopes) Contains(scope APIKeyScope) bool {
	for _, sc := range s {
		if sc == scope {
			return true
		}
	}
	return false
}

// APIKeyPrefix API密钥前缀
type APIKeyPrefix string

const (
	APIKeyPrefixSecret APIKeyPrefix = "sk" // Secret Key
	APIKeyPrefixPublic APIKeyPrefix = "pk" // Public Key
	APIKeyPrefixTest   APIKeyPrefix = "tk" // Test Key
)

// APIKey API密钥实体
type APIKey struct {
	KeyID      string       `json:"key_id" gorm:"primaryKey;type:varchar(36);comment:密钥ID"`
	TenantID   string       `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id;comment:租户ID"`
	ProjectID  string       `json:"project_id" gorm:"type:varchar(36);not null;index:idx_project_id;comment:项目ID"`
	KeyName    string       `json:"key_name" gorm:"type:varchar(100);not null;comment:密钥名称"`
	KeySecret  string       `json:"-" gorm:"type:varchar(255);not null;uniqueIndex;comment:密钥（加密存储）"` // 不暴露给前端
	KeyPrefix  APIKeyPrefix `json:"key_prefix" gorm:"type:varchar(20);not null;comment:密钥前缀"`
	KeyMasked  string       `json:"key_masked" gorm:"type:varchar(50);not null;comment:脱敏密钥（用于显示）"` // sk_****1234
	Scopes     APIKeyScopes `json:"scopes" gorm:"type:json;not null;comment:权限范围"`
	ExpiresAt  *int64       `json:"expires_at" gorm:"index:idx_expires_at;comment:过期时间"`
	LastUsedAt *int64       `json:"last_used_at" gorm:"index:idx_last_used;comment:最后使用时间"`
	Status     APIKeyStatus `json:"status" gorm:"type:enum('active','revoked','expired');default:'active';not null;index:idx_status;comment:状态"`
	CreatedAt  int64        `json:"created_at" gorm:"not null;default:0;comment:创建时间"`
	UpdatedAt  int64        `json:"updated_at" gorm:"not null;default:0;comment:更新时间"`
	DeletedAt  *int64       `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`

	// 关联
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID;references:ProjectID"`
}

// TableName 指定表名
func (APIKey) TableName() string {
	return "developer_api_keys"
}

// IsActive 是否激活
func (k *APIKey) IsActive() bool {
	if k.Status != APIKeyStatusActive {
		return false
	}

	// 检查是否过期
	if k.ExpiresAt != nil && time.Now().UnixMilli() > *k.ExpiresAt {
		return false
	}

	return true
}

// IsExpired 是否已过期
func (k *APIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().UnixMilli() > *k.ExpiresAt
}

// HasScope 检查是否有指定权限
func (k *APIKey) HasScope(scope APIKeyScope) bool {
	return k.Scopes.Contains(scope)
}

// GetExpiresAtAsTime 获取过期时间
func (k *APIKey) GetExpiresAtAsTime() *time.Time {
	if k.ExpiresAt == nil {
		return nil
	}
	t := time.Unix(*k.ExpiresAt/1000, 0)
	return &t
}

// GetLastUsedAtAsTime 获取最后使用时间
func (k *APIKey) GetLastUsedAtAsTime() *time.Time {
	if k.LastUsedAt == nil {
		return nil
	}
	t := time.Unix(*k.LastUsedAt/1000, 0)
	return &t
}
