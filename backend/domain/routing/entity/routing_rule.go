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

// RuleType 规则类型
type RuleType string

const (
	RuleTypeKeyword  RuleType = "keyword"  // 关键词匹配
	RuleTypeRegex    RuleType = "regex"    // 正则表达式匹配
	RuleTypeIntent   RuleType = "intent"   // 意图匹配
	RuleTypeCategory RuleType = "category" // 分类匹配
)

// RoutingRule 路由规则实体
type RoutingRule struct {
	RuleID          string          `json:"rule_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string          `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_priority"`
	RuleName        string          `json:"rule_name" gorm:"type:varchar(200);not null"`
	RuleType        RuleType        `json:"rule_type" gorm:"type:enum('keyword','regex','intent','category');not null;index:idx_rule_type"`
	Priority        int             `json:"priority" gorm:"not null;default:0"` // 数字越大优先级越高
	Condition       string          `json:"condition" gorm:"type:json;not null"` // JSON格式的匹配条件
	TargetBotID     *string         `json:"target_bot_id,omitempty" gorm:"type:varchar(36)"`
	TargetWorkflowID *string         `json:"target_workflow_id,omitempty" gorm:"type:varchar(36)"`
	IsActive        bool            `json:"is_active" gorm:"default:true;index:idx_is_active"`
	CreatedAt       int64           `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt       int64           `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (RoutingRule) TableName() string {
	return "routing_rules"
}

// IsActive 是否启用
func (r *RoutingRule) IsRuleActive() bool {
	return r.IsActive
}

// HasBotTarget 是否有Bot目标
func (r *RoutingRule) HasBotTarget() bool {
	return r.TargetBotID != nil && *r.TargetBotID != ""
}

// HasWorkflowTarget 是否有工作流目标
func (r *RoutingRule) HasWorkflowTarget() bool {
	return r.TargetWorkflowID != nil && *r.TargetWorkflowID != ""
}

// GetTargetID 获取目标ID（Bot或Workflow）
func (r *RoutingRule) GetTargetID() string {
	if r.HasBotTarget() {
		return *r.TargetBotID
	}
	if r.HasWorkflowTarget() {
		return *r.TargetWorkflowID
	}
	return ""
}

// RoutingLog 路由日志实体
type RoutingLog struct {
	LogID            string  `json:"log_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID         string  `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserInput        string  `json:"user_input" gorm:"type:text;not null"`
	MatchedRuleID    *string `json:"matched_rule_id,omitempty" gorm:"type:varchar(36)"`
	MatchedBotID     *string `json:"matched_bot_id,omitempty" gorm:"type:varchar(36)"`
	MatchedWorkflowID *string `json:"matched_workflow_id,omitempty" gorm:"type:varchar(36)"`
	Confidence       *float64 `json:"confidence,omitempty"`
	MatchType        *string `json:"match_type,omitempty" gorm:"type:varchar(50)"`
	RoutingScore     *float64 `json:"routing_score,omitempty"`
	CreatedAt        int64   `json:"created_at" gorm:"not null;index:idx_created_at"`
}

// TableName 指定表名
func (RoutingLog) TableName() string {
	return "routing_logs"
}
