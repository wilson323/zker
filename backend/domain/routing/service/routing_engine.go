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

package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
)

// IntentMatcher 意图匹配器接口
type IntentMatcher interface {
	Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error)
}

// ServiceHealth 服务健康状态
type ServiceHealth struct {
	ServiceID   string
	IsHealthy    bool
	SuccessRate  float64
	AvgLatency   time.Duration
	CurrentLoad  int
	MaxCapacity  int
}

// ServiceLoad 服务负载信息
type ServiceLoad struct {
	ServiceID   string
	CurrentLoad int
	MaxCapacity int
}

// ServiceRegistry 服务注册中心接口
type ServiceRegistry interface {
	GetBotHealth(ctx context.Context, botID string) (*ServiceHealth, error)
	GetWorkflowHealth(ctx context.Context, workflowID string) (*ServiceHealth, error)
}

// LoadMonitor 负载监控接口
type LoadMonitor interface {
	GetServiceLoad(ctx context.Context, serviceID string) (*ServiceLoad, error)
}

// RoutingDecision 路由决策结果
type RoutingDecision struct {
	BotID       string
	WorkflowID  string
	Confidence  float64
	Score       float64
	Reasons     []string
	MatchType   string
	RuleID      string
}

// ScoreBasedRouter 基于评分的路由决策引擎
type ScoreBasedRouter struct {
	intentMatcher   IntentMatcher
	serviceRegistry ServiceRegistry
	loadMonitor     LoadMonitor
	logRepo         repository.RoutingLogRepository

	// 评分权重配置
	intentWeight    float64 // 意图匹配权重
	healthWeight    float64 // 健康度权重
	loadWeight      float64 // 负载权重
	costWeight      float64 // 成本权重
	regionWeight    float64 // 区域权重
}

// NewScoreBasedRouter 创建评分路由器实例
func NewScoreBasedRouter(
	intentMatcher IntentMatcher,
	serviceRegistry ServiceRegistry,
	loadMonitor LoadMonitor,
	logRepo repository.RoutingLogRepository,
) *ScoreBasedRouter {
	return &ScoreBasedRouter{
		intentMatcher:   intentMatcher,
		serviceRegistry: serviceRegistry,
		loadMonitor:     loadMonitor,
		logRepo:         logRepo,
		intentWeight:    0.3, // 30%
		healthWeight:    0.2, // 20%
		loadWeight:      0.2, // 20%
		costWeight:      0.1, // 10%
		regionWeight:    0.1, // 10%
	}
}

// SetWeights 设置评分权重
func (r *ScoreBasedRouter) SetWeights(intent, health, load, cost, region float64) {
	r.intentWeight = intent
	r.healthWeight = health
	r.loadWeight = load
	r.costWeight = cost
	r.regionWeight = region
}

// Route 路由决策
func (r *ScoreBasedRouter) Route(ctx context.Context, input *MatchInput) (*RoutingDecision, error) {
	// 1. 意图匹配，获取候选列表
	candidates, err := r.intentMatcher.Match(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("intent matching failed: %w", err)
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no matching candidates found")
	}

	// 2. 为每个候选计算综合得分
	scoredCandidates := make([]*RoutingDecision, 0)

	for _, candidate := range candidates {
		score, reasons := r.calculateScore(ctx, candidate, input)

		scoredCandidates = append(scoredCandidates, &RoutingDecision{
			BotID:      candidate.BotID,
			WorkflowID: candidate.WorkflowID,
			Confidence: candidate.Confidence,
			Score:      score,
			Reasons:    reasons,
			MatchType:  candidate.MatchType,
			RuleID:     candidate.RuleID,
		})
	}

	// 3. 按得分排序
	sort.Slice(scoredCandidates, func(i, j int) bool {
		return scoredCandidates[i].Score > scoredCandidates[j].Score
	})

	// 4. 选择最佳候选
	bestCandidate := scoredCandidates[0]

	// 5. 记录路由日志
	r.logRouting(ctx, input, bestCandidate)

	return bestCandidate, nil
}

// calculateScore 计算综合得分
func (r *ScoreBasedRouter) calculateScore(ctx context.Context, candidate *MatchOutput, input *MatchInput) (float64, []string) {
	reasons := make([]string, 0)
	score := 0.0

	// 1. 意图匹配得分（权重30%）
	intentScore := candidate.Confidence * r.intentWeight
	score += intentScore
	reasons = append(reasons, fmt.Sprintf("Intent: %.2f", intentScore))

	// 2. 服务健康度得分（权重20%）
	var serviceID string
	if candidate.BotID != "" {
		serviceID = candidate.BotID
	} else if candidate.WorkflowID != "" {
		serviceID = candidate.WorkflowID
	}

	var health *ServiceHealth
	if candidate.BotID != "" {
		health, _ = r.serviceRegistry.GetBotHealth(ctx, candidate.BotID)
	} else if candidate.WorkflowID != "" {
		health, _ = r.serviceRegistry.GetWorkflowHealth(ctx, candidate.WorkflowID)
	}

	healthScore := 0.0
	if health != nil && health.IsHealthy {
		healthScore = health.SuccessRate * r.healthWeight
	} else {
		// 如果健康检查失败，给予较低分数
		healthScore = 0.1 * r.healthWeight
	}
	score += healthScore
	reasons = append(reasons, fmt.Sprintf("Health: %.2f", healthScore))

	// 3. 负载得分（权重20%）
	load, _ := r.loadMonitor.GetServiceLoad(ctx, serviceID)
	loadScore := 0.0
	if load != nil {
		if load.MaxCapacity > 0 {
			loadRatio := float64(load.CurrentLoad) / float64(load.MaxCapacity)
			loadScore = (1 - loadRatio) * r.loadWeight
		}
	} else {
		// 如果无法获取负载信息，给予中等分数
		loadScore = 0.5 * r.loadWeight
	}
	score += loadScore
	reasons = append(reasons, fmt.Sprintf("Load: %.2f", loadScore))

	// 4. 区域亲和性得分（权重10%）
	regionScore := 0.1 * r.regionWeight // 简化实现，假设都符合
	score += regionScore
	reasons = append(reasons, fmt.Sprintf("Region: %.2f", regionScore))

	// 5. 成本得分（权重10%）
	costScore := 0.1 * r.costWeight // 简化实现
	score += costScore
	reasons = append(reasons, fmt.Sprintf("Cost: %.2f", costScore))

	return score, reasons
}

// logRouting 记录路由日志
func (r *ScoreBasedRouter) logRouting(ctx context.Context, input *MatchInput, decision *RoutingDecision) {
	log := &entity.RoutingLog{
		LogID:            uuid.New().String(),
		TenantID:         input.TenantID,
		UserInput:        input.UserInput,
		MatchedRuleID:    &decision.RuleID,
		MatchedBotID:     stringPtr(decision.BotID),
		MatchedWorkflowID: stringPtr(decision.WorkflowID),
		Confidence:       &decision.Confidence,
		MatchType:        &decision.MatchType,
		RoutingScore:     &decision.Score,
	}

	// 异步记录日志，不阻塞主流程
	go func() {
		_ = r.logRepo.Create(context.Background(), log)
	}()
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetRoutingHistory 获取路由历史记录
func (r *ScoreBasedRouter) GetRoutingHistory(ctx context.Context, tenantID string, limit int) ([]*entity.RoutingLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	filter := &repository.RoutingLogFilter{
		TenantID: tenantID,
		PageSize:  limit,
	}

	logs, _, err := r.logRepo.List(ctx, filter)
	return logs, err
}

// RoutingEngine 路由引擎类型别名
type RoutingEngine = ScoreBasedRouter

// NewRoutingEngine 创建路由引擎实例（扩展版本）
func NewRoutingEngine(
	intentMatcher IntentMatcher,
	serviceRegistry ServiceRegistry,
	loadMonitor LoadMonitor,
	logRepo repository.RoutingLogRepository,
	intentWeight, healthWeight, loadWeight, costWeight, regionWeight float64,
) *RoutingEngine {
	engine := &ScoreBasedRouter{
		intentMatcher:   intentMatcher,
		serviceRegistry: serviceRegistry,
		loadMonitor:     loadMonitor,
		logRepo:         logRepo,
	}
	engine.SetWeights(intentWeight, healthWeight, loadWeight, costWeight, regionWeight)
	return engine
}
