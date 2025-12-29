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

package routing

// ==================== Request DTOs ====================

// CreateRoutingRuleRequest 创建路由规则请求
type CreateRoutingRuleRequest struct {
	TenantID        string                 `json:"tenant_id" binding:"required"`
	RuleName        string                 `json:"rule_name" binding:"required,min=1,max=200"`
	RuleType        string                 `json:"rule_type" binding:"required,oneof=keyword regex intent category"`
	Priority        int                    `json:"priority" binding:"required,min=0,max=1000"`
	Condition       map[string]interface{} `json:"condition" binding:"required"`
	TargetBotID     string                 `json:"target_bot_id,omitempty"`
	TargetWorkflowID string                 `json:"target_workflow_id,omitempty"`
	IsActive        bool                   `json:"is_active"`
}

// UpdateRoutingRuleRequest 更新路由规则请求
type UpdateRoutingRuleRequest struct {
	RuleName   *string `json:"rule_name,omitempty" binding:"omitempty,min=1,max=200"`
	Priority   *int    `json:"priority,omitempty" binding:"omitempty,min=0,max=1000"`
	Condition  *map[string]interface{} `json:"condition,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// ListRoutingRulesRequest 列出路由规则请求
type ListRoutingRulesRequest struct {
	TenantID  string  `form:"tenant_id" binding:"required"`
	RuleType  *string `form:"rule_type" binding:"omitempty,oneof=keyword regex intent category"`
	IsActive  *bool   `form:"is_active" binding:"omitempty"`
	PageSize  int     `form:"page_size" binding:"omitempty,min=1,max=100"`
	PageToken string  `form:"page_token" binding:"omitempty"`
}

// ExecuteRoutingRequest 执行路由请求
type ExecuteRoutingRequest struct {
	UserInput string                 `json:"user_input" binding:"required,min=1,max=5000"`
	TenantID  string                 `json:"tenant_id" binding:"required"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// ListRoutingLogsRequest 列出路由日志请求
type ListRoutingLogsRequest struct {
	TenantID string `form:"tenant_id" binding:"required"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// SetMatcherWeightsRequest 设置匹配器权重请求
type SetMatcherWeightsRequest struct {
	RuleWeight       float64 `json:"rule_weight" binding:"required,min=0,max=1"`
	SimilarityWeight float64 `json:"similarity_weight" binding:"required,min=0,max=1"`
	ModelWeight      float64 `json:"model_weight" binding:"required,min=0,max=1"`
}

// SetScorerWeightsRequest 设置评分权重请求
type SetScorerWeightsRequest struct {
	IntentWeight  float64 `json:"intent_weight" binding:"required,min=0,max=1"`
	HealthWeight  float64 `json:"health_weight" binding:"required,min=0,max=1"`
	LoadWeight    float64 `json:"load_weight" binding:"required,min=0,max=1"`
	CostWeight    float64 `json:"cost_weight" binding:"required,min=0,max=1"`
	RegionWeight  float64 `json:"region_weight" binding:"required,min=0,max=1"`
}

// SetSimilarityThresholdRequest 设置相似度阈值请求
type SetSimilarityThresholdRequest struct {
	Threshold float64 `json:"threshold" binding:"required,min=0,max=1"`
}

// GenerateBotEmbeddingRequest 生成Bot向量嵌入请求
type GenerateBotEmbeddingRequest struct {
	Description string `json:"description" binding:"required,min=1,max=2000"`
}

// HealthCheckRequest 健康检查请求
type HealthCheckRequest struct {
	BotID        string `json:"bot_id,omitempty"`
	WorkflowID   string `json:"workflow_id,omitempty"`
	CheckType    string `json:"check_type" binding:"required,oneof=success_rate latency"`
	WindowMinutes int   `json:"window_minutes" binding:"required,min=1,max=60"`
}

// LoadBalanceConfigRequest 负载均衡配置请求
type LoadBalanceConfigRequest struct {
	Strategy         string `json:"strategy" binding:"required,oneof=least_loaded round_robin random"`
	MaxLoadPercent   int    `json:"max_load_percent" binding:"required,min=1,max=100"`
}

// CircuitBreakerConfigRequest 熔断配置请求
type CircuitBreakerConfigRequest struct {
	FailureThreshold   int `json:"failure_threshold" binding:"required,min=1,max=100"`
	SuccessThreshold   int `json:"success_threshold" binding:"required,min=1,max=100"`
	TimeoutSeconds     int `json:"timeout_seconds" binding:"required,min=1,max=3600"`
	HalfOpenMaxCalls   int `json:"half_open_max_calls" binding:"required,min=1,max=100"`
}

// ==================== Response DTOs ====================

// RoutingRuleInfo 路由规则信息
type RoutingRuleInfo struct {
	RuleID            string                 `json:"rule_id"`
	RuleName          string                 `json:"rule_name"`
	RuleType          string                 `json:"rule_type"`
	Priority          int                    `json:"priority"`
	Condition         map[string]interface{} `json:"condition"`
	TargetBotID       string                 `json:"target_bot_id,omitempty"`
	TargetWorkflowID  string                 `json:"target_workflow_id,omitempty"`
	IsActive          bool                   `json:"is_active"`
	CreatedAt         int64                  `json:"created_at"`
	UpdatedAt         int64                  `json:"updated_at"`
}

// CreateRoutingRuleResponse 创建路由规则响应
type CreateRoutingRuleResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    RoutingRuleInfo      `json:"data"`
}

// GetRoutingRuleResponse 获取路由规则响应
type GetRoutingRuleResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    RoutingRuleInfo      `json:"data"`
}

// ListRoutingRulesResponse 列出路由规则响应
type ListRoutingRulesResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    ListRoutingRulesData      `json:"data"`
}

// ListRoutingRulesData 列出路由规则数据
type ListRoutingRulesData struct {
	Rules       []RoutingRuleInfo `json:"rules"`
	TotalCount  int               `json:"total_count"`
	NextPageToken string          `json:"next_page_token,omitempty"`
}

// RoutingDecision 路由决策
type RoutingDecision struct {
	BotID       string   `json:"bot_id,omitempty"`
	WorkflowID  string   `json:"workflow_id,omitempty"`
	Confidence  float64  `json:"confidence"`
	Score       float64  `json:"score"`
	Reasons     []string `json:"reasons"`
	MatchType   string   `json:"match_type"`
	RuleID      string   `json:"rule_id,omitempty"`
}

// ExecuteRoutingResponse 执行路由响应
type ExecuteRoutingResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    RoutingDecision   `json:"data"`
}

// RoutingLogInfo 路由日志信息
type RoutingLogInfo struct {
	LogID          string  `json:"log_id"`
	UserInput      string  `json:"user_input"`
	MatchedBotID   string  `json:"matched_bot_id,omitempty"`
	MatchedRuleID  string  `json:"matched_rule_id,omitempty"`
	Confidence     float64 `json:"confidence"`
	RoutingScore   float64 `json:"routing_score"`
	CreatedAt      int64   `json:"created_at"`
}

// ListRoutingLogsResponse 列出路由日志响应
type ListRoutingLogsResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    ListRoutingLogsData  `json:"data"`
}

// ListRoutingLogsData 列出路由日志数据
type ListRoutingLogsData struct {
	Logs []RoutingLogInfo `json:"logs"`
}

// SetWeightsResponse 设置权重响应
type SetWeightsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GenerateBotEmbeddingResponse 生成Bot向量嵌入响应
type GenerateBotEmbeddingResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServiceHealth 服务健康信息
type ServiceHealth struct {
	ServiceID    string  `json:"service_id"`
	IsHealthy    bool    `json:"is_healthy"`
	SuccessRate  float64 `json:"success_rate"`
	AvgLatency   int     `json:"avg_latency_ms"`
	CurrentLoad  int     `json:"current_load"`
	MaxCapacity  int     `json:"max_capacity"`
}

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    ServiceHealth     `json:"data"`
}

// RoutingStatsResponse 路由统计响应
type RoutingStatsResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    RoutingStatsData  `json:"data"`
}

// RoutingStatsData 路由统计数据
type RoutingStatsData struct {
	TotalRoutes              int                    `json:"total_routes"`
	SuccessRate              float64                `json:"success_rate"`
	AvgLatency               float64                `json:"avg_latency_ms"`
	BotDistribution          map[string]int         `json:"bot_distribution"`
	MatchTypeDistribution    map[string]int         `json:"match_type_distribution"`
}
