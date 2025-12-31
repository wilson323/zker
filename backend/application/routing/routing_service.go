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

import (
	"context"
	"encoding/json"
	"time"

	"github.com/coze-dev/coze-studio/backend/api/model/routing"
	routingentity "github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	routingservice "github.com/coze-dev/coze-studio/backend/domain/routing/service"
	"github.com/coze-dev/coze-studio/backend/infra/monitoring/metrics"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// stringPtr 返回字符串指针的辅助函数
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// stringVal 安全地从指针获取字符串值
func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// float64Val 安全地从指针获取 float64 值
func float64Val(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

// RoutingApplicationService 路由应用服务
type RoutingApplicationService struct {
	routingEngine    *routingservice.RoutingEngine
	ruleSVC         *routingservice.RuleService
	logSVC          *routingservice.RoutingLogService
	healthSVC       *routingservice.ServiceHealthMonitor
	loadBalancer    *routingservice.LoadBalancer
	circuitBreaker  *routingservice.CircuitBreakerService
}

// NewRoutingApplicationService 创建路由应用服务
func NewRoutingApplicationService(
	routingEngine *routingservice.RoutingEngine,
	ruleSVC *routingservice.RuleService,
	logSVC *routingservice.RoutingLogService,
	healthSVC *routingservice.ServiceHealthMonitor,
	loadBalancer *routingservice.LoadBalancer,
	circuitBreaker *routingservice.CircuitBreakerService,
) *RoutingApplicationService {
	return &RoutingApplicationService{
		routingEngine:   routingEngine,
		ruleSVC:        ruleSVC,
		logSVC:         logSVC,
		healthSVC:      healthSVC,
		loadBalancer:   loadBalancer,
		circuitBreaker: circuitBreaker,
	}
}

// ==================== 路由规则管理用例 ====================

// CreateRoutingRule 创建路由规则
func (s *RoutingApplicationService) CreateRoutingRule(ctx context.Context, req *routing.CreateRoutingRuleRequest) (*routing.RoutingRuleInfo, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}
	if req.RuleName == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "rule_name is required"))
	}
	if req.RuleType == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "rule_type is required"))
	}
	if req.Condition == nil {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "condition is required"))
	}
	if req.TargetBotID == "" && req.TargetWorkflowID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "target_bot_id or target_workflow_id is required"))
	}

	// 2. 序列化条件
	conditionJSON, err := json.Marshal(req.Condition)
	if err != nil {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invalid condition format"))
	}

	// 3. 构建实体
	rule := &routingentity.RoutingRule{
		TenantID:         req.TenantID,
		RuleName:         req.RuleName,
		RuleType:         routingentity.RuleType(req.RuleType),
		Priority:         req.Priority,
		Condition:        string(conditionJSON),
		TargetBotID:      stringPtr(req.TargetBotID),
		TargetWorkflowID: stringPtr(req.TargetWorkflowID),
		IsActive:         req.IsActive,
	}

	// 4. 调用领域服务
	createdRule, err := s.ruleSVC.CreateRule(ctx, rule)
	if err != nil {
		return nil, err
	}

	// 5. 转换为DTO
	return s.entityToRoutingRuleInfo(createdRule), nil
}

// GetRoutingRule 获取路由规则
func (s *RoutingApplicationService) GetRoutingRule(ctx context.Context, ruleID string) (*routing.RoutingRuleInfo, error) {
	if ruleID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "rule_id is required"))
	}

	rule, err := s.ruleSVC.GetRuleByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errorx.New(errno.ErrRouteNotFoundCode, errorx.KV("rule_id", ruleID))
	}

	return s.entityToRoutingRuleInfo(rule), nil
}

// UpdateRoutingRule 更新路由规则
func (s *RoutingApplicationService) UpdateRoutingRule(ctx context.Context, ruleID string, req *routing.UpdateRoutingRuleRequest) (*routing.RoutingRuleInfo, error) {
	if ruleID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "rule_id is required"))
	}

	// 1. 获取现有规则
	rule, err := s.ruleSVC.GetRuleByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errorx.New(errno.ErrRouteNotFoundCode, errorx.KV("rule_id", ruleID))
	}

	// 2. 应用更新
	if req.RuleName != nil {
		rule.RuleName = *req.RuleName
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.Condition != nil {
		conditionJSON, err := json.Marshal(*req.Condition)
		if err != nil {
			return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "invalid condition format"))
		}
		rule.Condition = string(conditionJSON)
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	// 3. 保存更新
	updatedRule, err := s.ruleSVC.UpdateRule(ctx, rule)
	if err != nil {
		return nil, err
	}

	return s.entityToRoutingRuleInfo(updatedRule), nil
}

// DeleteRoutingRule 删除路由规则
func (s *RoutingApplicationService) DeleteRoutingRule(ctx context.Context, ruleID string) error {
	if ruleID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "rule_id is required"))
	}

	return s.ruleSVC.DeleteRule(ctx, ruleID)
}

// ListRoutingRules 列出路由规则
func (s *RoutingApplicationService) ListRoutingRules(ctx context.Context, req *routing.ListRoutingRulesRequest) (*routing.ListRoutingRulesData, error) {
	// 1. 构建查询条件
	listReq := &routingservice.ListRulesRequest{
		TenantID: req.TenantID,
		RuleType: (*routingentity.RuleType)(req.RuleType),
		IsActive: req.IsActive,
		PageSize: req.PageSize,
		PageToken: req.PageToken,
	}

	// 2. 调用领域服务
	rules, total, nextPageToken, err := s.ruleSVC.ListRules(ctx, listReq)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	ruleDTOs := make([]routing.RoutingRuleInfo, 0, len(rules))
	for _, rule := range rules {
		ruleDTOs = append(ruleDTOs, *s.entityToRoutingRuleInfo(rule))
	}

	return &routing.ListRoutingRulesData{
		Rules:         ruleDTOs,
		TotalCount:    int(total),
		NextPageToken: nextPageToken,
	}, nil
}

// ==================== 路由决策用例 ====================

// ExecuteRouting 执行路由
func (s *RoutingApplicationService) ExecuteRouting(ctx context.Context, req *routing.ExecuteRoutingRequest) (*routing.RoutingDecision, error) {
	// 1. 参数验证
	if req.UserInput == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "user_input is required"))
	}
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 2. 构建匹配输入
	matchInput := &routingservice.MatchInput{
		UserInput: req.UserInput,
		TenantID:  req.TenantID,
		Context:   req.Context,
	}

	// 3. 记录开始时间
	startTime := time.Now()

	// 4. 调用路由引擎
	decision, err := s.routingEngine.Route(ctx, matchInput)
	if err != nil {
		return nil, err
	}

	// 5. 计算路由耗时
	duration := time.Since(startTime).Seconds()

	// 6. 确定匹配类型
	matchType := decision.MatchType
	if matchType == "" {
		matchType = "none"
	}

	// 7. 确定结果
	result := "success"
	if decision.BotID == "" {
		result = "no_match"
		metrics.RecordRoutingNoMatch(req.TenantID)
	}

	// 8. 记录路由执行指标
	confidence := decision.Confidence
	score := decision.Score

	metrics.RecordRoutingExecution(
		req.TenantID,
		matchType,
		result,
		duration,
		confidence,
		score,
	)

	// 9. 转换为DTO
	return s.matchOutputToRoutingDecision(decision), nil
}

// ==================== 路由日志用例 ====================

// ListRoutingLogs 列出路由日志
func (s *RoutingApplicationService) ListRoutingLogs(ctx context.Context, req *routing.ListRoutingLogsRequest) (*routing.ListRoutingLogsData, error) {
	if req.TenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	limit := req.Limit
	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// 调用领域服务获取日志
	logs, err := s.logSVC.GetRoutingLogs(ctx, req.TenantID, limit, 0)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	logDTOs := make([]routing.RoutingLogInfo, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, routing.RoutingLogInfo{
			LogID:         log.LogID,
			UserInput:     log.UserInput,
			MatchedBotID:  stringVal(log.MatchedBotID),
			MatchedRuleID: stringVal(log.MatchedRuleID),
			Confidence:    float64Val(log.Confidence),
			RoutingScore:  float64Val(log.RoutingScore),
			CreatedAt:     log.CreatedAt,
		})
	}

	return &routing.ListRoutingLogsData{
		Logs: logDTOs,
	}, nil
}

// ==================== 配置管理用例 ====================

// SetMatcherWeights 设置匹配器权重
func (s *RoutingApplicationService) SetMatcherWeights(ctx context.Context, req *routing.SetMatcherWeightsRequest) error {
	// 参数验证
	totalWeight := req.RuleWeight + req.SimilarityWeight + req.ModelWeight
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "weights must sum to 1.0"))
	}

	// 调用领域服务设置权重
	// 注意：ScoreBasedRouter 有 SetWeights 方法，不是 SetMatcherWeights
	s.routingEngine.SetWeights(req.RuleWeight, req.SimilarityWeight+req.ModelWeight, 0, 0, 0)
	return nil
}

// SetScorerWeights 设置评分权重
func (s *RoutingApplicationService) SetScorerWeights(ctx context.Context, req *routing.SetScorerWeightsRequest) error {
	// 参数验证
	totalWeight := req.IntentWeight + req.HealthWeight + req.LoadWeight + req.CostWeight + req.RegionWeight
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "weights must sum to 1.0"))
	}

	// 调用领域服务设置权重
	s.routingEngine.SetWeights(req.IntentWeight, req.HealthWeight, req.LoadWeight, req.CostWeight, req.RegionWeight)
	return nil
}

// SetSimilarityThreshold 设置相似度阈值
func (s *RoutingApplicationService) SetSimilarityThreshold(ctx context.Context, req *routing.SetSimilarityThresholdRequest) error {
	if req.Threshold < 0 || req.Threshold > 1 {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "threshold must be between 0 and 1"))
	}

	// TODO: 需要在 SimilarityMatcher 中实现 SetThreshold 方法
	return errorx.New(errno.ErrRouteParamInvalidCode, errorx.KV("msg", "method not implemented yet"))
}

// GenerateBotEmbedding 生成Bot向量嵌入
func (s *RoutingApplicationService) GenerateBotEmbedding(ctx context.Context, botID string, req *routing.GenerateBotEmbeddingRequest) error {
	if botID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "bot_id is required"))
	}
	if req.Description == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "description is required"))
	}

	// TODO: 需要在 RoutingEngine 中实现 GenerateBotEmbedding 方法
	return errorx.New(errno.ErrRouteParamInvalidCode, errorx.KV("msg", "method not implemented yet"))
}

// ==================== 服务健康检查用例 ====================

// HealthCheck 健康检查
func (s *RoutingApplicationService) HealthCheck(ctx context.Context, req *routing.HealthCheckRequest) (*routing.ServiceHealth, error) {
	// 1. 参数验证
	if req.BotID == "" && req.WorkflowID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "bot_id or workflow_id is required"))
	}

	var health *routingservice.ServiceHealth
	var err error

	// 2. 调用领域服务进行健康检查
	if req.BotID != "" {
		health, err = s.healthSVC.GetBotHealth(ctx, req.BotID)
	} else {
		health, err = s.healthSVC.GetWorkflowHealth(ctx, req.WorkflowID)
	}

	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	return &routing.ServiceHealth{
		ServiceID:   health.ServiceID,
		IsHealthy:   health.IsHealthy,
		SuccessRate: health.SuccessRate,
		AvgLatency:  int(health.AvgLatency.Milliseconds()),
		CurrentLoad: 0, // ServiceHealth 中没有 CurrentLoad 字段
		MaxCapacity: 0, // ServiceHealth 中没有 MaxCapacity 字段
	}, nil
}

// ==================== 负载均衡用例 ====================

// ConfigureLoadBalance 配置负载均衡
func (s *RoutingApplicationService) ConfigureLoadBalance(ctx context.Context, req *routing.LoadBalanceConfigRequest) error {
	// 调用 LoadBalancer 的 SetStrategy 和 SetThreshold 方法
	s.loadBalancer.SetStrategy(req.Strategy)
	s.loadBalancer.SetThreshold(req.MaxLoadPercent)
	return nil
}

// ==================== 熔断器用例 ====================

// ConfigureCircuitBreaker 配置熔断器
func (s *RoutingApplicationService) ConfigureCircuitBreaker(ctx context.Context, req *routing.CircuitBreakerConfigRequest) error {
	// TODO: CircuitBreakerService 尚未实现 Configure 方法
	return errorx.New(errno.ErrRouteParamInvalidCode, errorx.KV("msg", "method not implemented yet"))
}

// ==================== 路由统计用例 ====================

// GetRoutingStats 获取路由统计
func (s *RoutingApplicationService) GetRoutingStats(ctx context.Context, tenantID, startDate, endDate string) (*routing.RoutingStatsData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 调用领域服务获取统计
	stats, err := s.logSVC.GetRoutingStats(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return &routing.RoutingStatsData{
		TotalRoutes:           stats.TotalRoutings,
		SuccessRate:           0, // RoutingStats 中没有 SuccessRate 字段，暂时设为0
		AvgLatency:            0, // RoutingStats 中没有 AvgLatency 字段，暂时设为0
		BotDistribution:       make(map[string]int), // 暂时返回空map
		MatchTypeDistribution: make(map[string]int), // 暂时返回空map
	}, nil
}

// ==================== 转换方法 ====================

// entityToRoutingRuleInfo 实体转换为RoutingRuleInfo DTO
func (s *RoutingApplicationService) entityToRoutingRuleInfo(entity *routingentity.RoutingRule) *routing.RoutingRuleInfo {
	var condition map[string]interface{}
	_ = json.Unmarshal([]byte(entity.Condition), &condition)

	return &routing.RoutingRuleInfo{
		RuleID:           entity.RuleID,
		RuleName:         entity.RuleName,
		RuleType:         string(entity.RuleType),
		Priority:         entity.Priority,
		Condition:        condition,
		TargetBotID:      stringVal(entity.TargetBotID),
		TargetWorkflowID: stringVal(entity.TargetWorkflowID),
		IsActive:         entity.IsActive,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
	}
}

// matchOutputToRoutingDecision MatchOutput转换为RoutingDecision DTO
func (s *RoutingApplicationService) matchOutputToRoutingDecision(output *routingservice.RoutingDecision) *routing.RoutingDecision {
	return &routing.RoutingDecision{
		BotID:      output.BotID,
		WorkflowID: output.WorkflowID,
		Confidence: output.Confidence,
		Score:      output.Score,
		Reasons:    output.Reasons,
		MatchType:  output.MatchType,
		RuleID:     output.RuleID,
	}
}
