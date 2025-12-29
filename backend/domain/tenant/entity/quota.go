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

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeBots         ResourceType = "bots"         // Bot数量
	ResourceTypeMessages     ResourceType = "messages"     // 消息数量
	ResourceTypeStorage      ResourceType = "storage"      // 存储空间
	ResourceTypeTeamMembers  ResourceType = "team_members" // 团队成员数量
)

// ResetCycle 重置周期
type ResetCycle string

const (
	ResetCycleDaily   ResetCycle = "daily"   // 每日重置
	ResetCycleMonthly ResetCycle = "monthly" // 每月重置
	ResetCycleYearly  ResetCycle = "yearly"  // 每年重置
	ResetCycleNever   ResetCycle = "never"   // 不重置
)

// Quota 配额实体
type Quota struct {
	QuotaID      string        `json:"quota_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID     string        `json:"tenant_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_tenant_resource"`
	ResourceType ResourceType  `json:"resource_type" gorm:"type:enum('bots','messages','storage','team_members');not null;uniqueIndex:uk_tenant_resource"`
	MaxLimit     int           `json:"max_limit" gorm:"not null"` // -1表示无限制
	UsedCount    int           `json:"used_count" gorm:"default:0"`
	ResetCycle   ResetCycle    `json:"reset_cycle" gorm:"type:enum('daily','monthly','yearly','never');default:'monthly'"`
	LastResetAt  int64         `json:"last_reset_at" gorm:"not null;default:0"`
	CreatedAt    int64         `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64         `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (Quota) TableName() string {
	return "quotas"
}

// IsUnlimited 是否无限制
func (q *Quota) IsUnlimited() bool {
	return q.MaxLimit < 0
}

// IsExceeded 是否已超额
func (q *Quota) IsExceeded() bool {
	if q.IsUnlimited() {
		return false
	}
	return q.UsedCount >= q.MaxLimit
}

// GetRemainingCount 获取剩余配额
func (q *Quota) GetRemainingCount() int {
	if q.IsUnlimited() {
		return -1 // -1表示无限
	}
	remaining := q.MaxLimit - q.UsedCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetUsagePercentage 获取使用率百分比
func (q *Quota) GetUsagePercentage() float64 {
	if q.IsUnlimited() || q.MaxLimit == 0 {
		return 0
	}
	return float64(q.UsedCount) / float64(q.MaxLimit) * 100
}

// ShouldReset 是否应该重置
func (q *Quota) ShouldReset() bool {
	if q.ResetCycle == ResetCycleNever {
		return false
	}

	now := time.Now()
	lastReset := time.Unix(q.LastResetAt/1000, 0)

	switch q.ResetCycle {
	case ResetCycleDaily:
		// 如果上次重置不是今天
		return lastReset.Day() != now.Day() ||
			lastReset.Month() != now.Month() ||
			lastReset.Year() != now.Year()
	case ResetCycleMonthly:
		// 如果上次重置不是本月
		return lastReset.Month() != now.Month() ||
			lastReset.Year() != now.Year()
	case ResetCycleYearly:
		// 如果上次重置不是本年
		return lastReset.Year() != now.Year()
	}

	return false
}

// GetLastResetAsTime 获取上次重置时间
func (q *Quota) GetLastResetAsTime() time.Time {
	return time.Unix(q.LastResetAt/1000, 0)
}

// GetCreatedAtAsTime 获取创建时间
func (q *Quota) GetCreatedAtAsTime() time.Time {
	return time.Unix(q.CreatedAt/1000, 0)
}
