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

// CollaborationConfig 协同配置实体
type CollaborationConfig struct {
	ID                      int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID                string    `json:"tenant_id" gorm:"type:varchar(36);uniqueIndex;not null"`

	// 自动审核阈值
	AutoReviewThreshold     float64   `json:"auto_review_threshold" gorm:"type:decimal(3,2);default:0.70"`
	AutoEscalateThreshold   float64   `json:"auto_escalate_threshold" gorm:"type:decimal(3,2);default:0.30"`

	// SLA配置
	DefaultSLAMinutes       int       `json:"default_sla_minutes" gorm:"default:60"`
	UrgentSLAMinutes        int       `json:"urgent_sla_minutes" gorm:"default:15"`

	// 审核人员池配置
	ReviewPoolFilter        JSON      `json:"review_pool_filter" gorm:"type:json"`

	// 升级规则配置
	EscalationRules         JSON      `json:"escalation_rules" gorm:"type:json"`

	// 优先级规则配置
	PriorityRules           JSON      `json:"priority_rules" gorm:"type:json"`

	// 通知配置
	NotificationEnabled     bool      `json:"notification_enabled" gorm:"default:true"`
	NotificationChannels    JSON      `json:"notification_channels" gorm:"type:json"`

	// 时间戳
	CreatedAt               int64     `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt               int64     `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (CollaborationConfig) TableName() string {
	return "collaboration_configs"
}

// GetDefaultSLADeadline 获取默认SLA截止时间
func (c *CollaborationConfig) GetDefaultSLADeadline() time.Time {
	return time.Now().Add(time.Duration(c.DefaultSLAMinutes) * time.Minute)
}

// GetUrgentSLADeadline 获取紧急任务SLA截止时间
func (c *CollaborationConfig) GetUrgentSLADeadline() time.Time {
	return time.Now().Add(time.Duration(c.UrgentSLAMinutes) * time.Minute)
}

// ShouldAutoReview 判断是否应该自动审核
func (c *CollaborationConfig) ShouldAutoReview(confidence float64) bool {
	return confidence >= c.AutoReviewThreshold
}

// ShouldAutoEscalate 判断是否应该自动升级
func (c *CollaborationConfig) ShouldAutoEscalate(confidence float64) bool {
	return confidence < c.AutoEscalateThreshold
}

// EscalationRule 升级规则
type EscalationRule struct {
	Condition      string                 `json:"condition"`      // 触发条件
	Action         string                 `json:"action"`         // 升级动作
	TargetLevel    int                    `json:"target_level"`   // 目标级别
	NotifyUsers    []string               `json:"notify_users"`   // 通知用户列表
	AdditionalData map[string]interface{} `json:"additional_data"` // 附加数据
}

// PriorityRule 优先级规则
type PriorityRule struct {
	Name        string                 `json:"name"`
	Condition   string                 `json:"condition"`
	Priority    TaskPriority           `json:"priority"`
	SLAMinutes  int                    `json:"sla_minutes"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ReviewPoolFilter 审核人员池过滤规则
type ReviewPoolFilter struct {
	RoleIDs     []string `json:"role_ids,omitempty"`
	Department  string   `json:"department,omitempty"`
	Skills      []string `json:"skills,omitempty"`
	MaxLoad     int      `json:"max_load,omitempty"`
	OnlineOnly  bool     `json:"online_only,omitempty"`
}

// NotificationChannel 通知渠道配置
type NotificationChannel struct {
	Type     string                 `json:"type"`     // email, sms, webhook, in_app
	Enabled  bool                   `json:"enabled"`
	Config   map[string]interface{} `json:"config"`
}
