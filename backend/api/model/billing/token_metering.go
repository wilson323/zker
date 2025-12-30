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

package billing

// ============================================================
// Token计量API Model定义
// ============================================================

// RecordTokenUsageRequest Token使用记录请求
type RecordTokenUsageRequest struct {
	TenantID       string  `json:"tenant_id" binding:"required"`
	UserID         *uint64 `json:"user_id,omitempty"`
	BotID          *string `json:"bot_id,omitempty"`
	ConversationID *string `json:"conversation_id,omitempty"`
	MessageID      *string `json:"message_id,omitempty"`

	ModelProvider string `json:"model_provider" binding:"required"` // openai, anthropic, etc
	ModelName     string `json:"model_name" binding:"required"`    // gpt-4, claude-3-opus, etc
	ModelVersion  *string `json:"model_version,omitempty"`

	InputTokens  int    `json:"input_tokens" binding:"required,min=0"`
	OutputTokens int    `json:"output_tokens" binding:"required,min=0"`
	TotalTokens  int    `json:"total_tokens" binding:"required,min=0"`

	RequestType    string `json:"request_type" binding:"required"`    // chat, completion, embedding
	ResponseTimeMs *int   `json:"response_time_ms,omitempty"`
	LatencyMs      *int   `json:"latency_ms,omitempty"`
	IsCached       bool   `json:"is_cached"`
	Metadata       string `json:"metadata,omitempty"`
}

// RecordTokenUsageResponse Token使用记录响应
type RecordTokenUsageResponse struct {
	LogID      uint64  `json:"log_id"`
	InputCost  float64 `json:"input_cost"`
	OutputCost float64 `json:"output_cost"`
	TotalCost  float64 `json:"total_cost"`
	CreatedAt  int64   `json:"created_at"`
}

// BatchRecordTokenUsageRequest 批量记录请求
type BatchRecordTokenUsageRequest struct {
	Records []*RecordTokenUsageRequest `json:"records" binding:"required,min=1,max=1000"`
}

// BatchRecordTokenUsageResponse 批量记录响应
type BatchRecordTokenUsageResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	LogIDs       []uint64 `json:"log_ids"`
	TotalCost    float64  `json:"total_cost"`
}

// GetUsageStatsRequest 获取使用统计请求
type GetUsageStatsRequest struct {
	BotID       *string `form:"bot_id,omitempty"`
	ModelName   *string `form:"model_name,omitempty"`
	StartDate   string  `form:"start_date,omitempty" binding:"omitempty,date_format=2006-01-02"`
	EndDate     string  `form:"end_date,omitempty" binding:"omitempty,date_format=2006-01-02"`
	Granularity string  `form:"granularity,omitempty" binding:"omitempty,oneof=daily hourly"`
}

// UsageStatsResponse 使用统计响应
type UsageStatsResponse struct {
	TenantID          string  `json:"tenant_id"`
	TotalInputTokens  int64   `json:"total_input_tokens"`
	TotalOutputTokens int64   `json:"total_output_tokens"`
	TotalTokens       int64   `json:"total_tokens"`
	TotalCost         float64 `json:"total_cost"`
	TotalRequests     int64   `json:"total_requests"`
	CachedRequests    int64   `json:"cached_requests"`
	AvgResponseTime   float64 `json:"avg_response_time_ms"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
}

// ModelUsageStatsResponse 模型使用统计响应
type ModelUsageStatsResponse struct {
	ModelStats []*ModelUsageStats `json:"model_stats"`
	Total      *AggregateStats    `json:"total,omitempty"`
}

// ModelUsageStats 模型使用统计
type ModelUsageStats struct {
	ModelProvider     string  `json:"model_provider"`
	ModelName         string  `json:"model_name"`
	TotalTokens       int64   `json:"total_tokens"`
	TotalCost         float64 `json:"total_cost"`
	TotalRequests     int64   `json:"total_requests"`
	AvgCostPer1KToken float64 `json:"avg_cost_per_1k_token"`
	CostPercentage   float64 `json:"cost_percentage,omitempty"` // 占总成本的百分比
}

// AggregateStats 汇总统计
type AggregateStats struct {
	TotalTokens      int64   `json:"total_tokens"`
	TotalCost        float64 `json:"total_cost"`
	TotalRequests    int64   `json:"total_requests"`
	TopModel         string  `json:"top_model,omitempty"`
	TopModelCost     float64 `json:"top_model_cost,omitempty"`
}

// DailyUsageStatsResponse 每日使用趋势响应
type DailyUsageStatsResponse struct {
	DailyStats []*DailyUsageStats `json:"daily_stats"`
	Total      *AggregateStats    `json:"total,omitempty"`
}

// DailyUsageStats 每日使用统计
type DailyUsageStats struct {
	Date            string  `json:"date"`
	TotalTokens     int64   `json:"total_tokens"`
	TotalCost       float64 `json:"total_cost"`
	TotalRequests   int64   `json:"total_requests"`
	CachedRequests  int64   `json:"cached_requests"`
	AvgResponseTime float64 `json:"avg_response_time_ms"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ============================================================
// 预算管理API Model定义
// ============================================================

// BudgetResponse 预算配置响应
type BudgetResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    *BudgetSettingsData  `json:"data,omitempty"`
}

// BudgetSettingsData 预算配置数据
type BudgetSettingsData struct {
	TenantID     string   `json:"tenant_id"`
	BudgetType   string   `json:"budget_type"`   // monthly, quarterly, yearly
	BudgetAmount float64  `json:"budget_amount"`
	Currency     string   `json:"currency"`

	// 告警配置
	AlertThreshold1 int      `json:"alert_threshold_1"` // 80%
	AlertThreshold2 int      `json:"alert_threshold_2"` // 95%
	HardCapEnabled  bool     `json:"hard_cap_enabled"`
	HardCapAmount   *float64 `json:"hard_cap_amount,omitempty"`

	// 降级配置
	AutoDowngradeEnabled bool    `json:"auto_downgrade_enabled"`
	DowngradeThreshold   int     `json:"downgrade_threshold"` // 90%
	OriginalModel        *string `json:"original_model,omitempty"`
	FallbackModel        *string `json:"fallback_model,omitempty"`

	// 通知配置
	NotificationChannels   []string `json:"notification_channels"`
	NotificationRecipients []string `json:"notification_recipients"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// BudgetUsageResponse 预算使用情况响应
type BudgetUsageResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *BudgetUsageData  `json:"data,omitempty"`
}

// BudgetUsageData 预算使用情况数据
type BudgetUsageData struct {
	// 预算信息
	BudgetAmount    float64  `json:"budget_amount"`
	UsedAmount      float64  `json:"used_amount"`
	RemainingAmount float64  `json:"remaining_amount"`
	UsagePercent    float64  `json:"usage_percent"`

	// 当前周期统计
	PeriodStart    int64  `json:"period_start"`
	PeriodEnd      int64  `json:"period_end"`
	TotalTokens    int64  `json:"total_tokens"`
	TotalRequests  int    `json:"total_requests"`
	AverageCost    float64 `json:"average_cost"`

	// 告警统计
	AlertCount     int     `json:"alert_count"`
	LastAlertAt    *int64  `json:"last_alert_at,omitempty"`
	LastAlertLevel *string `json:"last_alert_level,omitempty"`

	// 预测
	EstimatedDailyUsage float64 `json:"estimated_daily_usage"`
	DaysRemaining       int     `json:"days_remaining"`
	WillExceedBudget    bool    `json:"will_exceed_budget"`

	UpdatedAt int64 `json:"updated_at"`
}

// BudgetCheckResponse 预算检查响应
type BudgetCheckResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    *BudgetCheckData   `json:"data,omitempty"`
}

// BudgetCheckData 预算检查数据
type BudgetCheckData struct {
	TenantID       string  `json:"tenant_id"`
	BudgetAmount   float64 `json:"budget_amount"`
	UsedAmount     float64 `json:"used_amount"`
	UsagePercent   float64 `json:"usage_percent"`
	AlertTriggered bool    `json:"alert_triggered"`
	AlertLevel     *string `json:"alert_level,omitempty"`
	Message        string  `json:"message"`
	CheckedAt      int64   `json:"checked_at"`
}

// BudgetAlertsResponse 告警列表响应
type BudgetAlertsResponse struct {
	Code       int                  `json:"code"`
	Message    string               `json:"message"`
	Data       []*BudgetAlertData   `json:"data,omitempty"`
	Pagination *PaginationInfo      `json:"pagination,omitempty"`
}

// BudgetAlertData 告警数据
type BudgetAlertData struct {
	ID                  uint64    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	AlertType           string    `json:"alert_type"`
	BudgetAmount        float64   `json:"budget_amount"`
	UsedAmount          float64   `json:"used_amount"`
	UsagePercent        float64   `json:"usage_percent"`
	AlertLevel          string    `json:"alert_level"`
	AlertMessage        *string   `json:"alert_message,omitempty"`
	NotificationChannels []string `json:"notification_channels,omitempty"`
	NotificationStatus  string    `json:"notification_status"`
	SentAt              *int64    `json:"sent_at,omitempty"`
	ErrorMessage        *string   `json:"error_message,omitempty"`
	CreatedAt           int64     `json:"created_at"`
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	NextPageToken string `json:"next_page_token,omitempty"`
	TotalCount    int64  `json:"total_count"`
	HasMore       bool   `json:"has_more"`
}
