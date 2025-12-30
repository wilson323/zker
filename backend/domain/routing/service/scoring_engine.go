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

	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	"go.uber.org/zap"
)

// ScoringEngine 评分引擎
// 职责：根据多维度指标对Agent进行评分
type ScoringEngine struct {
	loadMonitor     LoadMonitor
	serviceRegistry ServiceRegistry
	optLogRepo      repository.RoutingOptimizeLogRepository
	logger          *zap.Logger

	// 评分权重配置
	loadWeight      float64 // 负载权重
	healthWeight    float64 // 健康度权重
	responseWeight  float64 // 响应时间权重
	errorWeight     float64 // 错误率权重
	costWeight      float64 // 成本权重
	preferenceWeight float64 // 用户偏好权重
}

// NewScoringEngine 创建评分引擎实例
func NewScoringEngine(
	loadMonitor LoadMonitor,
	serviceRegistry ServiceRegistry,
	optLogRepo repository.RoutingOptimizeLogRepository,
	logger *zap.Logger,
) *ScoringEngine {
	return &ScoringEngine{
		loadMonitor:      loadMonitor,
		serviceRegistry:  serviceRegistry,
		optLogRepo:       optLogRepo,
		logger:           logger,
		loadWeight:       0.25,  // 25%
		healthWeight:     0.20,  // 20%
		responseWeight:   0.20,  // 20%
		errorWeight:      0.15,  // 15%
		costWeight:       0.10,  // 10%
		preferenceWeight: 0.10,  // 10%
	}
}

// SetWeights 设置评分权重
func (e *ScoringEngine) SetWeights(load, health, response, error, cost, preference float64) {
	e.loadWeight = load
	e.healthWeight = health
	e.responseWeight = response
	e.errorWeight = error
	e.costWeight = cost
	e.preferenceWeight = preference
}

// ScoreAgent 对Agent进行综合评分
// 评分维度：负载、响应时间、错误率、成本、用户偏好
func (e *ScoringEngine) ScoreAgent(ctx context.Context, agentID string, routingCtx *RoutingContext) (*AgentScore, error) {
	score := &AgentScore{
		AgentID:    agentID,
		TotalScore: 0.0,
		Reasons:    make([]string, 0),
	}

	// 1. 负载评分（权重25%）
	loadScore, err := e.scoreLoad(ctx, agentID)
	if err != nil {
		e.logger.Warn("failed to score load", zap.String("agent_id", agentID), zap.Error(err))
		loadScore = 0.5 // 默认中等分数
	}
	score.LoadScore = loadScore
	score.TotalScore += loadScore * e.loadWeight
	score.Reasons = append(score.Reasons, fmt.Sprintf("Load: %.2f", loadScore))

	// 2. 健康度评分（权重20%）
	healthScore, err := e.scoreHealth(ctx, agentID)
	if err != nil {
		e.logger.Warn("failed to score health", zap.String("agent_id", agentID), zap.Error(err))
		healthScore = 0.5
	}
	score.HealthScore = healthScore
	score.TotalScore += healthScore * e.healthWeight
	score.Reasons = append(score.Reasons, fmt.Sprintf("Health: %.2f", healthScore))

	// 3. 响应时间评分（权重20%）
	responseScore, err := e.scoreResponse(ctx, agentID)
	if err != nil {
		e.logger.Warn("failed to score response", zap.String("agent_id", agentID), zap.Error(err))
		responseScore = 0.5
	}
	score.Reasons = append(score.Reasons, fmt.Sprintf("Response: %.2f", responseScore))
	score.TotalScore += responseScore * e.responseWeight

	// 4. 错误率评分（权重15%）
	errorScore, err := e.scoreError(ctx, agentID)
	if err != nil {
		e.logger.Warn("failed to score error", zap.String("agent_id", agentID), zap.Error(err))
		errorScore = 0.5
	}
	score.Reasons = append(score.Reasons, fmt.Sprintf("Error: %.2f", errorScore))
	score.TotalScore += errorScore * e.errorWeight

	// 5. 成本评分（权重10%）
	costScore := e.scoreCost(ctx, agentID)
	score.CostScore = costScore
	score.TotalScore += costScore * e.costWeight
	score.Reasons = append(score.Reasons, fmt.Sprintf("Cost: %.2f", costScore))

	// 6. 用户偏好评分（权重10%）
	preferenceScore := e.scorePreference(ctx, agentID, routingCtx)
	score.PreferenceScore = preferenceScore
	score.TotalScore += preferenceScore * e.preferenceWeight
	score.Reasons = append(score.Reasons, fmt.Sprintf("Preference: %.2f", preferenceScore))

	// 7. 计算置信度
	score.Confidence = e.calculateConfidence(score)
	score.PredictedLoad = 1.0 - loadScore // 负载预测

	return score, nil
}

// scoreLoad 负载评分
// 负载越低，分数越高
func (e *ScoringEngine) scoreLoad(ctx context.Context, agentID string) (float64, error) {
	load, err := e.loadMonitor.GetServiceLoad(ctx, agentID)
	if err != nil {
		return 0.5, err
	}

	if load == nil || load.MaxCapacity == 0 {
		return 0.5, nil // 默认中等分数
	}

	// 计算负载率（0-1）
	loadRatio := float64(load.CurrentLoad) / float64(load.MaxCapacity)

	// 负载率越低，分数越高
	// 0%负载 → 1.0分
	// 100%负载 → 0.0分
	return 1.0 - loadRatio, nil
}

// scoreHealth 健康度评分
// 健康度越高，分数越高
func (e *ScoringEngine) scoreHealth(ctx context.Context, agentID string) (float64, error) {
	health, err := e.serviceRegistry.GetBotHealth(ctx, agentID)
	if err != nil {
		return 0.5, err
	}

	if health == nil {
		return 0.5, nil
	}

	// 如果不健康，返回低分
	if !health.IsHealthy {
		return 0.1, nil
	}

	// 健康度 = 成功率
	return health.SuccessRate, nil
}

// scoreResponse 响应时间评分
// 响应时间越短，分数越高
func (e *ScoringEngine) scoreResponse(ctx context.Context, agentID string) (float64, error) {
	// 获取历史路由日志
	logs, err := e.optLogRepo.GetByAgent(ctx, "", agentID, 100)
	if err != nil || len(logs) == 0 {
		return 0.5, err // 默认中等分数
	}

	// 计算平均响应时间
	totalTime := 0
	for _, log := range logs {
		totalTime += log.ResponseTime
	}
	avgTime := float64(totalTime) / float64(len(logs))

	// 响应时间评分（ms）
	// < 1000ms → 1.0分
	// 1000-3000ms → 0.5-1.0分
	// > 3000ms → 0.0-0.5分
	if avgTime < 1000 {
		return 1.0, nil
	} else if avgTime < 3000 {
		return 1.0 - (avgTime-1000)/2000, nil
	} else {
		return math.Max(0.0, 1.0-(avgTime/5000)), nil
	}
}

// scoreError 错误率评分
// 错误率越低，分数越高
func (e *ScoringEngine) scoreError(ctx context.Context, agentID string) (float64, error) {
	// 获取历史路由日志
	logs, err := e.optLogRepo.GetByAgent(ctx, "", agentID, 100)
	if err != nil || len(logs) == 0 {
		return 0.5, err // 默认中等分数
	}

	// 计算错误率（低评分视为错误）
	errorCount := 0
	for _, log := range logs {
		if log.UserRating > 0 && log.UserRating <= 2 {
			errorCount++
		}
	}
	errorRate := float64(errorCount) / float64(len(logs))

	// 错误率评分
	// 0% → 1.0分
	// 50% → 0.5分
	// 100% → 0.0分
	return 1.0 - errorRate, nil
}

// scoreCost 成本评分
// 简化实现：可以根据Agent类型、模型等级等计算成本
func (e *ScoringEngine) scoreCost(ctx context.Context, agentID string) float64 {
	// 简化实现：返回默认成本分数
	// 实际应该从Agent配置中获取成本信息
	return 0.8
}

// scorePreference 用户偏好评分
// 根据用户历史偏好评分
func (e *ScoringEngine) scorePreference(ctx context.Context, agentID string, routingCtx *RoutingContext) float64 {
	if routingCtx == nil || routingCtx.Request == nil {
		return 0.5
	}

	// 检查用户历史路由记录
	for _, historyAgentID := range routingCtx.UserHistory {
		if historyAgentID == agentID {
			// 用户曾经使用过这个Agent，给予加分
			return 0.9
		}
	}

	// 检查用户偏好设置
	if routingCtx.Request.Preferences != nil {
		if preferredAgent, ok := routingCtx.Request.Preferences["preferred_agent"]; ok {
			if preferredAgent == agentID {
				return 1.0
			}
		}
	}

	return 0.5
}

// calculateConfidence 计算评分置信度
func (e *ScoringEngine) calculateConfidence(score *AgentScore) float64 {
	// 简化实现：基于评分的标准差计算置信度
	// 实际应该使用更复杂的统计方法

	// 如果各项评分差异不大，置信度高
	scores := []float64{
		score.LoadScore,
		score.HealthScore,
		score.PreferenceScore,
	}

	mean := (scores[0] + scores[1] + scores[2]) / 3.0
	variance := 0.0
	for _, s := range scores {
		variance += math.Pow(s-mean, 2)
	}
	variance /= 3.0
	stdDev := math.Sqrt(variance)

	// 标准差越小，置信度越高
	confidence := 1.0 / (1.0 + stdDev)
	return confidence
}

// BatchScoreAgents 批量对Agent评分
func (e *ScoringEngine) BatchScoreAgents(ctx context.Context, agentIDs []string, routingCtx *RoutingContext) ([]*AgentScore, error) {
	scores := make([]*AgentScore, 0, len(agentIDs))

	for _, agentID := range agentIDs {
		score, err := e.ScoreAgent(ctx, agentID, routingCtx)
		if err != nil {
			e.logger.Warn("failed to score agent",
				zap.String("agent_id", agentID),
				zap.Error(err))
			continue
		}
		scores = append(scores, score)
	}

	return scores, nil
}
