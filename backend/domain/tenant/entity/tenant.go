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

// TenantType 租户类型
type TenantType string

const (
	TenantTypeIndividual TenantType = "individual" // 个人租户
	TenantTypeTeam       TenantType = "team"       // 团队租户
	TenantTypeEnterprise TenantType = "enterprise" // 企业租户
)

// TenantStatus 租户状态
type TenantStatus string

const (
	TenantStatusActive    TenantStatus = "active"    // 激活
	TenantStatusSuspended TenantStatus = "suspended" // 暂停
	TenantStatusDeleted   TenantStatus = "deleted"   // 已删除
)

// SubscriptionTier 订阅等级
type SubscriptionTier string

const (
	SubscriptionTierFree       SubscriptionTier = "free"       // 免费版
	SubscriptionTierPro        SubscriptionTier = "pro"        // 专业版
	SubscriptionTierEnterprise SubscriptionTier = "enterprise" // 企业版
)

// IsolationStrategy 隔离策略类型
type IsolationStrategy string

const (
	// StrategyRowLevel 行级隔离（共享表，通过tenant_id区分）
	// 适用场景：小型租户，数据量小，成本敏感
	IsolationStrategyRowLevel IsolationStrategy = "row_level"

	// StrategySchemaLevel Schema级隔离（独立Schema）
	// 适用场景：中型租户，需要更好的性能和隔离性
	IsolationStrategySchemaLevel IsolationStrategy = "schema_level"

	// StrategyDatabaseLevel 数据库级隔离（独立数据库实例）
	// 适用场景：大型企业租户，要求最高隔离性和性能
	IsolationStrategyDatabaseLevel IsolationStrategy = "database_level"
)

// Tenant 租户实体
type Tenant struct {
	TenantID           string             `json:"tenant_id" gorm:"primaryKey;type:varchar(36)"`
	TenantName         string             `json:"tenant_name" gorm:"type:varchar(200);not null"`
	TenantType         TenantType         `json:"tenant_type" gorm:"type:enum('individual','team','enterprise');not null"`
	Subdomain          string             `json:"subdomain" gorm:"type:varchar(64);uniqueIndex;not null"` // 租户子域名（唯一）
	Status             TenantStatus       `json:"status" gorm:"type:enum('active','suspended','deleted');default:'active'"`
	SubscriptionTier   SubscriptionTier   `json:"subscription_tier" gorm:"type:enum('free','pro','enterprise');default:'free'"`
	IsolationStrategy  IsolationStrategy  `json:"isolation_strategy" gorm:"type:enum('row_level','schema_level','database_level');default:'row_level'"` // 隔离策略
	ContactEmail       string             `json:"contact_email,omitempty" gorm:"type:varchar(255)"`            // 联系邮箱
	ContactPhone       string             `json:"contact_phone,omitempty" gorm:"type:varchar(32)"`             // 联系电话
	CreatedAt          int64              `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt          int64              `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt          *int64             `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:TenantID;references:TenantID"`
	Quotas       []Quota       `json:"quotas,omitempty" gorm:"foreignKey:TenantID;references:TenantID"`
}

// TableName 指定表名
func (Tenant) TableName() string {
	return "tenants"
}

// IsActive 是否激活
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

// IsDeleted 是否已删除
func (t *Tenant) IsDeleted() bool {
	return t.Status == TenantStatusDeleted || t.DeletedAt != nil
}

// GetCreatedAtAsTime 获取创建时间
func (t *Tenant) GetCreatedAtAsTime() time.Time {
	return time.Unix(t.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (t *Tenant) GetUpdatedAtAsTime() time.Time {
	return time.Unix(t.UpdatedAt/1000, 0)
}
