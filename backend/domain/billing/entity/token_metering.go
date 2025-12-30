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

// 类型别名 - 用于向后兼容
// TokenLog 是 TokenUsageLog 的别名
type TokenLog = TokenUsageLog

// BudgetConfig 是 BudgetSettings 的别名
type BudgetConfig = BudgetSettings

// ============================================================
// Token计量实体 - 基于设计文档: 23-租户计费系统_TokenMetering补充_完整版.md
// ============================================================

// TokenUsageLog Token使用明细日志
type TokenUsageLog struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:日志ID"`

	// 租户与用户信息
	TenantID       string  `json:"tenant_id" gorm:"not null;index:idx_tenant_created,priority:1;size:64;comment:租户ID"`
	UserID         *uint64 `json:"user_id,omitempty" gorm:"comment:用户ID"`

	// 关联对象信息
	BotID          *string `json:"bot_id,omitempty" gorm:"index:idx_bot_created,priority:1;size:64;comment:Bot ID"`
	ConversationID *string `json:"conversation_id,omitempty" gorm:"size:64;comment:会话ID"`
	MessageID      *string `json:"message_id,omitempty" gorm:"size:64;comment:消息ID"`

	// Token统计
	InputTokens    int `json:"input_tokens" gorm:"not null;comment:输入Token数"`
	OutputTokens   int `json:"output_tokens" gorm:"not null;comment:输出Token数"`
	TotalTokens    int `json:"total_tokens" gorm:"not null;comment:总Token数"`

	// 模型信息
	ModelProvider string  `json:"model_provider" gorm:"not null;index:idx_model,priority:1;size:50;comment:模型提供商"`
	ModelName     string  `json:"model_name" gorm:"not null;index:idx_model,priority:2;size:50;comment:模型名称"`
	ModelVersion  *string `json:"model_version,omitempty" gorm:"size:20;comment:模型版本"`

	// 成本计算
	UnitPrice  float64 `json:"unit_price" gorm:"not null;type:decimal(10,6);comment:每1K Token单价(CNY)"`
	InputCost  float64 `json:"input_cost" gorm:"not null;type:decimal(10,6);comment:输入成本"`
	OutputCost float64 `json:"output_cost" gorm:"not null;type:decimal(10,6);comment:输出成本"`
	TotalCost  float64 `json:"total_cost" gorm:"not null;index:idx_cost;type:decimal(10,6);comment:总成本"`

	// 性能指标
	ResponseTimeMs *int   `json:"response_time_ms,omitempty" gorm:"comment:响应时间(毫秒)"`
	LatencyMs      *int   `json:"latency_ms,omitempty" gorm:"comment:首字延迟(毫秒)"`
	IsCached       bool   `json:"is_cached" gorm:"default:false;comment:是否缓存命中"`

	// 元数据
	RequestType string       `json:"request_type" gorm:"not null;size:20;comment:请求类型"`
	Metadata    string       `json:"metadata,omitempty" gorm:"type:json;comment:额外元数据"`
	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime;index:idx_tenant_created,priority:2;index:idx_bot_created,priority:2;comment:创建时间"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (TokenUsageLog) TableName() string {
	return "token_usage_logs"
}

// ============================================================

// TokenUsageSummary Token使用汇总表（按天/小时）
type TokenUsageSummary struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement"`

	// 分组维度
	TenantID    string  `json:"tenant_id" gorm:"not null;uniqueIndex:uk_tenant_bot_date_hour,priority:1;index:idx_tenant_date,priority:1;size:64;comment:租户ID"`
	BotID       *string `json:"bot_id,omitempty" gorm:"uniqueIndex:uk_tenant_bot_date_hour,priority:2;index:idx_bot_date,priority:1;size:64;comment:Bot ID (NULL表示整体汇总)"`
	SummaryDate string  `json:"summary_date" gorm:"not null;type:date;uniqueIndex:uk_tenant_bot_date_hour,priority:3;index:idx_tenant_date,priority:2;index:idx_bot_date,priority:2;comment:汇总日期"`
	SummaryHour *uint8  `json:"summary_hour,omitempty" gorm:"uniqueIndex:uk_tenant_bot_date_hour,priority:4;comment:小时 (0-23, NULL表示日汇总)"`

	// Token汇总
	TotalInputTokens  int64 `json:"total_input_tokens" gorm:"not null;default:0;comment:总输入Token数"`
	TotalOutputTokens int64 `json:"total_output_tokens" gorm:"not null;default:0;comment:总输出Token数"`
	TotalTokens       int64 `json:"total_tokens" gorm:"not null;default:0;comment:总Token数"`

	// 成本汇总
	TotalCost float64 `json:"total_cost" gorm:"not null;default:0;type:decimal(12,6);comment:总成本"`

	// 请求统计（企业级计费系统：支持int64类型以应对大规模请求）
	TotalRequests   int64   `json:"total_requests" gorm:"not null;default:0;comment:总请求数"`
	CachedRequests  int64   `json:"cached_requests" gorm:"not null;default:0;comment:缓存命中请求数"`
	AvgResponseTime *float64 `json:"avg_response_time,omitempty" gorm:"type:decimal(8,2);comment:平均响应时间(ms)"`

	// 模型分布
	ModelDistribution string  `json:"model_distribution,omitempty" gorm:"type:json;comment:模型使用分布"`

	// 时间戳
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (TokenUsageSummary) TableName() string {
	return "token_usage_summary"
}

// ============================================================

// BudgetSettings 预算配置表
type BudgetSettings struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID string `json:"tenant_id" gorm:"not null;uniqueIndex:uk_tenant_id;size:64;comment:租户ID"`

	// 预算配置
	BudgetType string  `json:"budget_type" gorm:"not null;size:20;default:'monthly';comment:预算周期"`
	BudgetAmount float64 `json:"budget_amount" gorm:"not null;type:decimal(12,2);comment:预算金额(CNY)"`
	Currency  string  `json:"currency" gorm:"size:3;default:'CNY';comment:货币"`

	// 告警阈值
	AlertThreshold1 int     `json:"alert_threshold_1" gorm:"default:80;comment:一级告警阈值(%)"`
	AlertThreshold2 int     `json:"alert_threshold_2" gorm:"default:95;comment:二级告警阈值(%)"`
	HardCapEnabled  bool    `json:"hard_cap_enabled" gorm:"default:false;comment:是否启用硬性上限"`
	HardCapAmount   *float64 `json:"hard_cap_amount,omitempty" gorm:"type:decimal(12,2);comment:硬性上限金额"`

	// 降级策略
	AutoDowngradeEnabled bool    `json:"auto_downgrade_enabled" gorm:"default:false;comment:是否启用自动降级"`
	DowngradeThreshold   int     `json:"downgrade_threshold" gorm:"default:90;comment:降级阈值(%)"`
	OriginalModel        *string `json:"original_model,omitempty" gorm:"size:100;comment:原模型"`
	FallbackModel        *string `json:"fallback_model,omitempty" gorm:"size:100;comment:降级模型"`

	// 通知配置
	NotificationChannels    string  `json:"notification_channels,omitempty" gorm:"type:json;comment:通知渠道"`
	NotificationRecipients  string  `json:"notification_recipients,omitempty" gorm:"type:json;comment:通知接收人列表"`

	// 时间戳
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (BudgetSettings) TableName() string {
	return "budget_settings"
}

// ============================================================

// BudgetAlert 预算告警历史表
type BudgetAlert struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID string  `json:"tenant_id" gorm:"not null;index:idx_tenant_created,priority:1;size:64;comment:租户ID"`

	// 告警类型
	AlertType string `json:"alert_type" gorm:"not null;index:idx_alert_type;size:30;comment:告警类型"` // threshold_1, threshold_2, hard_cap, downgrade

	// 预算状态
	BudgetAmount   float64 `json:"budget_amount" gorm:"not null;type:decimal(12,2);comment:预算金额"`
	UsedAmount     float64 `json:"used_amount" gorm:"not null;type:decimal(12,2);comment:已使用金额"`
	UsagePercent   float64 `json:"usage_percent" gorm:"not null;type:decimal(5,2);comment:使用率(%)"`

	// 告警信息
	AlertLevel   string  `json:"alert_level" gorm:"not null;size:20;comment:告警级别"` // warning, critical, emergency
	AlertMessage *string `json:"alert_message,omitempty" gorm:"type:text;comment:告警消息"`

	// 发送状态
	NotificationChannels string  `json:"notification_channels,omitempty" gorm:"type:json;comment:通知渠道"`
	NotificationStatus   string  `json:"notification_status" gorm:"default:'pending';index:idx_status;size:20;comment:通知状态"` // pending, sent, failed
	SentAt               *time.Time `json:"sent_at,omitempty" gorm:"comment:发送时间"`
	ErrorMessage         *string `json:"error_message,omitempty" gorm:"type:text;comment:错误信息"`

	// 时间戳
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime;index:idx_tenant_created,priority:2;comment:创建时间"`
}

// TableName 指定表名
func (BudgetAlert) TableName() string {
	return "budget_alerts_history"
}

// ============================================================

// CostOptimizationSuggestion 成本优化建议表
type CostOptimizationSuggestion struct {
	ID     uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID string  `json:"tenant_id" gorm:"not null;index:idx_tenant_status,priority:1;size:64;comment:租户ID"`

	// 建议类型
	SuggestionType string  `json:"suggestion_type" gorm:"not null;index:idx_type;size:50;comment:建议类型"` // model_downgrade, enable_cache, batch_request, prompt_optimization

	// 分析数据
	AnalysisPeriodStart string  `json:"analysis_period_start" gorm:"not null;type:date;comment:分析开始日期"`
	AnalysisPeriodEnd   string  `json:"analysis_period_end" gorm:"not null;type:date;comment:分析结束日期"`

	// 建议内容
	TargetBotID    *string `json:"target_bot_id,omitempty" gorm:"size:64;comment:目标Bot ID"`
	CurrentModel   *string `json:"current_model,omitempty" gorm:"size:100;comment:当前模型"`
	SuggestedModel *string `json:"suggested_model,omitempty" gorm:"size:100;comment:建议模型"`
	Reason         *string `json:"reason,omitempty" gorm:"type:text;comment:优化原因"`

	// 预期效果
	EstimatedMonthlySaving  *float64 `json:"estimated_monthly_saving,omitempty" gorm:"type:decimal(12,2);index:idx_saving;comment:预计月节省金额"`
	EstimatedSavingPercent  *float64 `json:"estimated_saving_percent,omitempty" gorm:"type:decimal(5,2);comment:预计节省比例(%)"`

	// 状态
	Status    string   `json:"status" gorm:"default:'pending';index:idx_tenant_status,priority:2;size:20;comment:状态"` // pending, approved, rejected, applied
	AppliedAt *time.Time `json:"applied_at,omitempty" gorm:"comment:应用时间"`
	AppliedBy *uint64   `json:"applied_by,omitempty" gorm:"comment:应用操作人ID"`

	// 时间戳
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (CostOptimizationSuggestion) TableName() string {
	return "cost_optimization_suggestions"
}

// ============================================================
// 辅助类型和常量
// ============================================================

// RequestType 请求类型枚举
const (
	RequestTypeChat        = "chat"
	RequestTypeCompletion  = "completion"
	RequestTypeEmbedding   = "embedding"
	RequestTypeRerank      = "rerank"
)

// BudgetType 预算类型枚举
const (
	BudgetTypeMonthly   = "monthly"
	BudgetTypeQuarterly = "quarterly"
	BudgetTypeYearly    = "yearly"
)

// AlertType 告警类型枚举
const (
	AlertTypeThreshold1 = "threshold_1"
	AlertTypeThreshold2 = "threshold_2"
	AlertTypeHardCap    = "hard_cap"
	AlertTypeDowngrade  = "downgrade"
)

// AlertLevel 告警级别枚举
const (
	AlertLevelWarning   = "warning"
	AlertLevelCritical  = "critical"
	AlertLevelEmergency = "emergency"
)

// NotificationStatus 通知状态枚举
const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

// SuggestionType 优化建议类型枚举
const (
	SuggestionTypeModelDowngrade     = "model_downgrade"
	SuggestionTypeEnableCache        = "enable_cache"
	SuggestionTypeBatchRequest       = "batch_request"
	SuggestionTypePromptOptimization = "prompt_optimization"
)

// SuggestionStatus 建议状态枚举
const (
	SuggestionStatusPending  = "pending"
	SuggestionStatusApproved = "approved"
	SuggestionStatusRejected = "rejected"
	SuggestionStatusApplied  = "applied"
)
