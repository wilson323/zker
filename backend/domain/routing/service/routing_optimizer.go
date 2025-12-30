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
	"math"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// RoutingRequest 路由请求
type RoutingRequest struct {
	TenantID     string            `json:"tenant_id"`
	UserID       string            `json:"user_id"`
	Query        string            `json:"query"`
	Context      map[string]string `json:"context"`
	Preferences  map[string]string `json:"preferences"` // 用户偏好
	Region       string            `json:"region"`       // 区域
	DeviceType   string            `json:"device_type"`  // 设备类型
	TimeWindow   TimeWindow        `json:"time_window"`  // 时间窗口
}

// TimeWindow 时间窗口
type TimeWindow struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// RoutingContext 路由上下文
type RoutingContext struct {
	Request       *RoutingRequest
	Intent        string
	Entities      []*entity.Entity
	CandidateAgents []string
	CurrentLoad   map[string]int
	UserHistory   []string // 历史路由记录
}

// LoadPrediction 负载预测
type LoadPrediction struct {
	AgentID       string    `json:"agent_id"`
	PredictedLoad float64   `json:"predicted_load"`
	Confidence    float64   `json:"confidence"`
	TimeWindow    TimeWindow `json:"time_window"`
}

// RoutingOptimizer 路由优化器
// 职责：根据多维度评分选择最优Agent，实现智能分发优化
type RoutingOptimizer struct {
	routingEngine    *ScoreBasedRouter
	scoringEngine    *ScoringEngine
	ruleRepo         repository.RoutingRuleRepository
	optLogRepo       repository.RoutingOptimizeLogRepository
	logger           *zap.Logger

	// 优化器配置
	enableLoadPrediction bool
	enableCostOptimization bool
	enableUserPreference  bool
}

// NewRoutingOptimizer 创建路由优化器实例
func NewRoutingOptimizer(
	routingEngine *ScoreBasedRouter,
	scoringEngine *ScoringEngine,
	ruleRepo repository.RoutingRuleRepository,
	optLogRepo repository.RoutingOptimizeLogRepository,
	logger *zap.Logger,
) *RoutingOptimizer {
	return &RoutingOptimizer{
		routingEngine:         routingEngine,
		scoringEngine:         scoringEngine,
		ruleRepo:              ruleRepo,
		optLogRepo:            optLogRepo,
		logger:                logger,
		enableLoadPrediction:  true,
		enableCostOptimization: true,
		enableUserPreference:   true,
	}
}

// OptimizeRouting 优化路由决策
// 根据多维度评分选择最优Agent
func (o *RoutingOptimizer) OptimizeRouting(ctx context.Context, req *RoutingRequest) (*entity.RoutingDecision, error) {
	// 1. 构建路由上下文
	routingCtx, err := o.buildRoutingContext(ctx, req)
	if err != nil {
		o.logger.Error("failed to build routing context",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		return nil, errno.ROUTING500001.WithDetail("error", err.Error())
	}

	// 2. 获取候选Agent列表
	candidates, err := o.getCandidateAgents(ctx, routingCtx)
	if err != nil {
		o.logger.Error("failed to get candidate agents",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		return nil, errno.ROUTING500001.WithDetail("error", err.Error())
	}

	if len(candidates) == 0 {
		return nil, errno.ROUTING404001.WithDetail("reason", "no available agents")
	}

	// 3. 为每个候选Agent评分
	scores := make([]*AgentScore, 0, len(candidates))
	for _, agentID := range candidates {
		score, err := o.scoringEngine.ScoreAgent(ctx, agentID, routingCtx)
		if err != nil {
			o.logger.Warn("failed to score agent",
				zap.String("agent_id", agentID),
				zap.Error(err))
			continue
		}
		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return nil, errno.ROUTING500001.WithDetail("reason", "no valid agent scores")
	}

	// 4. 选择得分最高的Agent
	bestScore := scores[0]
	for _, score := range scores {
		if score.TotalScore > bestScore.TotalScore {
			bestScore = score
		}
	}

	// 5. 构建路由决策
	decision := &entity.RoutingDecision{
		AgentID:      bestScore.AgentID,
		Confidence:   bestScore.Confidence,
		Score:        bestScore.TotalScore,
		Reasons:      bestScore.Reasons,
		MatchType:    "optimized",
		IntentID:     routingCtx.Intent,
		Entities:     routingCtx.Entities,
		Strategy:     "multi_dimension_score",
		PredictedLoad: bestScore.PredictedLoad,
	}

	// 6. 异步记录路由日志
	go o.logRouting(context.Background(), req, decision, bestScore)

	return decision, nil
}

// GetOptimalAgent 获取最优Agent（简化版）
func (o *RoutingOptimizer) GetOptimalAgent(ctx context.Context, intent string, agentCandidates []string) (string, error) {
	if len(agentCandidates) == 0 {
		return "", errno.ROUTING404001.WithDetail("reason", "empty candidate list")
	}

	// 如果只有一个候选，直接返回
	if len(agentCandidates) == 1 {
		return agentCandidates[0], nil
	}

	// TODO: 实现基于意图的Agent选择逻辑
	// 简化实现：返回第一个候选
	return agentCandidates[0], nil
}

// PredictAgentLoad 预测Agent负载
func (o *RoutingOptimizer) PredictAgentLoad(ctx context.Context, agentID string, timeWindow TimeWindow) (*LoadPrediction, error) {
	// 简化实现：使用历史数据预测
	// 实际应该使用时间序列预测模型（如ARIMA、LSTM）

	// 1. 获取历史负载数据
	historyLogs, err := o.optLogRepo.GetByAgent(ctx, "", agentID, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent history: %w", err)
	}

	// 2. 简单平均预测
	if len(historyLogs) == 0 {
		return &LoadPrediction{
			AgentID:       agentID,
			PredictedLoad: 0.5,
			Confidence:    0.5,
			TimeWindow:    timeWindow,
		}, nil
	}

	// 计算平均负载和成功率
	totalLoad := 0.0
	successCount := 0
	for _, log := range historyLogs {
		// 假设response_time越短，负载越小
		load := 1.0 - math.Min(float64(log.ResponseTime)/5000.0, 1.0)
		totalLoad += load
		if log.UserRating >= 4 {
			successCount++
		}
	}

	avgLoad := totalLoad / float64(len(historyLogs))
	successRate := float64(successCount) / float64(len(historyLogs))

	return &LoadPrediction{
		AgentID:       agentID,
		PredictedLoad: avgLoad,
		Confidence:    successRate,
		TimeWindow:    timeWindow,
	}, nil
}

// buildRoutingContext 构建路由上下文
func (o *RoutingOptimizer) buildRoutingContext(ctx context.Context, req *RoutingRequest) (*RoutingContext, error) {
	routingCtx := &RoutingContext{
		Request:    req,
		CandidateAgents: make([]string, 0),
		CurrentLoad: make(map[string]int),
		UserHistory: make([]string, 0),
	}

	// 获取激活的路由规则
	rules, err := o.ruleRepo.GetActiveRulesByTenant(ctx, req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active rules: %w", err)
	}

	// 提取候选Agent
	for _, rule := range rules {
		if rule.HasBotTarget() {
			routingCtx.CandidateAgents = append(routingCtx.CandidateAgents, rule.TargetBotID)
		}
	}

	return routingCtx, nil
}

// getCandidateAgents 获取候选Agent列表
func (o *RoutingOptimizer) getCandidateAgents(ctx context.Context, routingCtx *RoutingContext) ([]string, error) {
	// 去重
	uniqueAgents := make(map[string]bool)
	for _, agentID := range routingCtx.CandidateAgents {
		uniqueAgents[agentID] = true
	}

	candidates := make([]string, 0, len(uniqueAgents))
	for agentID := range uniqueAgents {
		candidates = append(candidates, agentID)
	}

	return candidates, nil
}

// logRouting 记录路由日志
func (o *RoutingOptimizer) logRouting(ctx context.Context, req *RoutingRequest, decision *entity.RoutingDecision, score *AgentScore) {
	log := &entity.RoutingOptimizeLog{
		LogID:        generateID(),
		TenantID:     req.TenantID,
		UserID:       req.UserID,
		Query:        req.Query,
		Intent:       decision.IntentID,
		AgentID:      decision.AgentID,
		Strategy:     decision.Strategy,
		ResponseTime: 0, // 需要从外部传入
		UserRating:   0,
	}

	if err := o.optLogRepo.Create(ctx, log); err != nil {
		o.logger.Error("failed to create routing optimize log",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
	}
}

// generateID 生成唯一ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// AgentScore Agent评分结果
type AgentScore struct {
	AgentID        string    `json:"agent_id"`
	TotalScore     float64   `json:"total_score"`
	Confidence     float64   `json:"confidence"`
	PredictedLoad  float64   `json:"predicted_load"`
	LoadScore      float64   `json:"load_score"`
	HealthScore    float64   `json:"health_score"`
	CostScore      float64   `json:"cost_score"`
	PreferenceScore float64  `json:"preference_score"`
	Reasons        []string  `json:"reasons"`
}
