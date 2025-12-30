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

// BillingCycle 计费周期
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly" // 月付
	BillingCycleYearly  BillingCycle = "yearly"  // 年付
)

// SubscriptionStatus 订阅状态
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"    // 激活
	SubscriptionStatusExpired   SubscriptionStatus = "expired"   // 过期
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled" // 取消
	SubscriptionStatusSuspended SubscriptionStatus = "suspended" // 暂停
)

// Subscription 订阅实体
type Subscription struct {
	SubscriptionID string              `json:"subscription_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID       string              `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	PlanTier       SubscriptionTier    `json:"plan_tier" gorm:"type:enum('free','pro','enterprise');not null"`
	BillingCycle   BillingCycle        `json:"billing_cycle" gorm:"type:enum('monthly','yearly');not null"`
	StartDate      time.Time           `json:"start_date" gorm:"type:date;not null"`
	EndDate        *time.Time          `json:"end_date,omitempty" gorm:"type:date"`
	AutoRenew      bool                `json:"auto_renew" gorm:"default:true"`
	Status         SubscriptionStatus  `json:"status" gorm:"type:enum('active','expired','cancelled');default:'active'"`
	CreatedAt      int64               `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt      int64               `json:"updated_at" gorm:"not null;default:0"`

	// 关联
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID;references:TenantID"`
}

// TableName 指定表名
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsActive 是否激活
func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive
}

// IsExpired 是否过期
func (s *Subscription) IsExpired() bool {
	if s.Status == SubscriptionStatusExpired {
		return true
	}
	if s.EndDate != nil && s.EndDate.Before(time.Now()) {
		return true
	}
	return false
}

// GetCreatedAtAsTime 获取创建时间
func (s *Subscription) GetCreatedAtAsTime() time.Time {
	return time.Unix(s.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (s *Subscription) GetUpdatedAtAsTime() time.Time {
	return time.Unix(s.UpdatedAt/1000, 0)
}
