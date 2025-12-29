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
	"fmt"
	"time"

	"github.com/coze-studio/backend/api/model/routing"
	routingentity "github.com/coze-studio/backend/domain/routing/entity"
	routingrepo "github.com/coze-studio/backend/domain/routing/repository"
	routingservice "github.com/coze-studio/backend/domain/routing/service"
	"github.com/coze-studio/backend/infra/monitoring/metrics"
	"github.com/coze-studio/backend/pkg/errorx"
	"github.com/coze-studio/backend/types/errno"
)

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
		TenantID:        req.TenantID,
		RuleName:        req.RuleName,
		RuleType:        routingentity.RuleType(req.RuleType),
		Priority:        req.Priority,
		Condition:       string(conditionJSON),
		TargetBotID:     req.TargetBotID,
		TargetWorkflowID: req.TargetWorkflowID,
		IsActive:        req.IsActive,
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
		return nil, errorx.New(errno.RoutingRuleNotFoundCode, errorx.KV("rule_id", ruleID))
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
		return nil, errorx.New(errno.RoutingRuleNotFoundCode, errorx.KV("rule_id", ruleID))
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
		TotalCount:    total,
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
	matchType := "none"
	if decision.MatchedRule != nil {
		switch decision.MatchedRule.RuleType {
		case routingentity.RuleTypeKeyword:
			matchType = "keyword"
		case routingentity.RuleTypeRegex:
			matchType = "regex"
		case routingentity.RuleTypeIntent:
			matchType = "similarity"
		case routingentity.RuleTypeCategory:
			matchType = "category"
		}
	}

	// 7. 确定结果
	result := "success"
	if decision.BotID == "" {
		result = "no_match"
		metrics.RecordRoutingNoMatch(req.TenantID)
	}

	// 8. 记录路由执行指标
	confidence := 0.0
	if decision.Confidence != nil {
		confidence = *decision.Confidence
	}

	score := 0.0
	if decision.Score != nil {
		score = *decision.Score
	}

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
	logs, err := s.logSVC.GetLogsByTenant(ctx, req.TenantID, limit)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	logDTOs := make([]routing.RoutingLogInfo, 0, len(logs))
	for _, log := range logs {
		logDTOs = append(logDTOs, routing.RoutingLogInfo{
			LogID:         log.LogID,
			UserInput:     log.UserInput,
			MatchedBotID:  log.MatchedBotID,
			MatchedRuleID: log.MatchedRuleID,
			Confidence:    log.Confidence,
			RoutingScore:  log.RoutingScore,
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
	return s.routingEngine.SetMatcherWeights(ctx, req.RuleWeight, req.SimilarityWeight, req.ModelWeight)
}

// SetScorerWeights 设置评分权重
func (s *RoutingApplicationService) SetScorerWeights(ctx context.Context, req *routing.SetScorerWeightsRequest) error {
	// 参数验证
	totalWeight := req.IntentWeight + req.HealthWeight + req.LoadWeight + req.CostWeight + req.RegionWeight
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "weights must sum to 1.0"))
	}

	// 调用领域服务设置权重
	return s.routingEngine.SetScorerWeights(ctx, req.IntentWeight, req.HealthWeight, req.LoadWeight, req.CostWeight, req.RegionWeight)
}

// SetSimilarityThreshold 设置相似度阈值
func (s *RoutingApplicationService) SetSimilarityThreshold(ctx context.Context, req *routing.SetSimilarityThresholdRequest) error {
	if req.Threshold < 0 || req.Threshold > 1 {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "threshold must be between 0 and 1"))
	}

	return s.routingEngine.SetSimilarityThreshold(ctx, req.Threshold)
}

// GenerateBotEmbedding 生成Bot向量嵌入
func (s *RoutingApplicationService) GenerateBotEmbedding(ctx context.Context, botID string, req *routing.GenerateBotEmbeddingRequest) error {
	if botID == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "bot_id is required"))
	}
	if req.Description == "" {
		return errorx.New(errno.InvalidRequest, errorx.KV("msg", "description is required"))
	}

	return s.routingEngine.GenerateBotEmbedding(ctx, botID, req.Description)
}

// ==================== 服务健康检查用例 ====================

// HealthCheck 健康检查
func (s *RoutingApplicationService) HealthCheck(ctx context.Context, req *routing.HealthCheckRequest) (*routing.ServiceHealth, error) {
	// 1. 参数验证
	if req.BotID == "" && req.WorkflowID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "bot_id or workflow_id is required"))
	}

	var serviceID string
	if req.BotID != "" {
		serviceID = fmt.Sprintf("bot:%s", req.BotID)
	} else {
		serviceID = fmt.Sprintf("workflow:%s", req.WorkflowID)
	}

	// 2. 调用领域服务进行健康检查
	checkType := routingservice.CheckType(req.CheckType)
	windowMinutes := req.WindowMinutes

	health, err := s.healthSVC.CheckServiceHealth(ctx, serviceID, checkType, windowMinutes)
	if err != nil {
		return nil, err
	}

	// 3. 转换为DTO
	return &routing.ServiceHealth{
		ServiceID:   health.ServiceID,
		IsHealthy:   health.IsHealthy,
		SuccessRate: health.SuccessRate,
		AvgLatency:  int(health.AvgLatency.Milliseconds()),
		CurrentLoad: health.CurrentLoad,
		MaxCapacity: health.MaxCapacity,
	}, nil
}

// ==================== 负载均衡用例 ====================

// ConfigureLoadBalance 配置负载均衡
func (s *RoutingApplicationService) ConfigureLoadBalance(ctx context.Context, req *routing.LoadBalanceConfigRequest) error {
	strategy := routingservice.LoadBalanceStrategy(req.Strategy)

	return s.loadBalancer.Configure(ctx, strategy, req.MaxLoadPercent)
}

// ==================== 熔断器用例 ====================

// ConfigureCircuitBreaker 配置熔断器
func (s *RoutingApplicationService) ConfigureCircuitBreaker(ctx context.Context, req *routing.CircuitBreakerConfigRequest) error {
	config := &routingservice.CircuitBreakerConfig{
		FailureThreshold: req.FailureThreshold,
		SuccessThreshold: req.SuccessThreshold,
		TimeoutSeconds:   req.TimeoutSeconds,
		HalfOpenMaxCalls: req.HalfOpenMaxCalls,
	}

	return s.circuitBreaker.Configure(ctx, config)
}

// ==================== 路由统计用例 ====================

// GetRoutingStats 获取路由统计
func (s *RoutingApplicationService) GetRoutingStats(ctx context.Context, tenantID, startDate, endDate string) (*routing.RoutingStatsData, error) {
	if tenantID == "" {
		return nil, errorx.New(errno.InvalidRequest, errorx.KV("msg", "tenant_id is required"))
	}

	// 调用领域服务获取统计
	stats, err := s.logSVC.GetStatistics(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 转换为DTO
	return &routing.RoutingStatsData{
		TotalRoutes:           stats.TotalRoutes,
		SuccessRate:           stats.SuccessRate,
		AvgLatency:            stats.AvgLatency,
		BotDistribution:       stats.BotDistribution,
		MatchTypeDistribution: stats.MatchTypeDistribution,
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
		TargetBotID:      entity.TargetBotID,
		TargetWorkflowID: entity.TargetWorkflowID,
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
