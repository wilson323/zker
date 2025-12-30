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

// RoutingLearningService 路由学习服务
// 职责：从历史路由日志中学习，优化路由策略
type RoutingLearningService struct {
	abTestRepo       repository.ABTestRepository
	routingLogRepo   repository.RoutingOptimizeLogRepository
	ruleRepo         repository.RoutingRuleRepository
	intentRepo       repository.IntentRepository
	optimizer        *RoutingOptimizer
	logger           *zap.Logger
}

// NewRoutingLearningService 创建路由学习服务实例
func NewRoutingLearningService(
	abTestRepo repository.ABTestRepository,
	routingLogRepo repository.RoutingOptimizeLogRepository,
	ruleRepo repository.RoutingRuleRepository,
	intentRepo repository.IntentRepository,
	optimizer *RoutingOptimizer,
	logger *zap.Logger,
) *RoutingLearningService {
	return &RoutingLearningService{
		abTestRepo:     abTestRepo,
		routingLogRepo: routingLogRepo,
		ruleRepo:       ruleRepo,
		intentRepo:     intentRepo,
		optimizer:      optimizer,
		logger:         logger,
	}
}

// RoutingInsight 路由洞察
type RoutingInsight struct {
	TenantID          string                 `json:"tenant_id"`
	TimeRange         TimeRange              `json:"time_range"`
	TotalRequests     int                    `json:"total_requests"`
	AvgResponseTime   float64                `json:"avg_response_time"`
	SuccessRate       float64                `json:"success_rate"`
	IntentDistribution map[string]int        `json:"intent_distribution"`
	AgentPerformance  map[string]*AgentPerformanceStats `json:"agent_performance"`
	Recommendations   []string               `json:"recommendations"`
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// AgentPerformanceStats Agent性能统计
type AgentPerformanceStats struct {
	AgentID         string  `json:"agent_id"`
	RequestCount    int     `json:"request_count"`
	AvgResponseTime float64 `json:"avg_response_time"`
	SuccessRate     float64 `json:"success_rate"`
	AvgUserRating   float64 `json:"avg_user_rating"`
	ErrorRate       float64 `json:"error_rate"`
}

// RoutingRuleSuggestion 路由规则建议
type RoutingRuleSuggestion struct {
	Intent         string  `json:"intent"`
	SuggestedAgent string  `json:"suggested_agent"`
	Confidence     float64 `json:"confidence"`
	Reason         string  `json:"reason"`
	ExpectedImprovement float64 `json:"expected_improvement"`
}

// OptimizedStrategy 优化后的策略
type OptimizedStrategy struct {
	StrategyID        string                 `json:"strategy_id"`
	StrategyName      string                 `json:"strategy_name"`
	Config            map[string]interface{} `json:"config"`
	ExpectedMetrics   map[string]float64     `json:"expected_metrics"`
}

// LearnFromRoutingLogs 从路由日志中学习
func (s *RoutingLearningService) LearnFromRoutingLogs(ctx context.Context, tenantID string, timeRange TimeRange) (*RoutingInsight, error) {
	// 1. 获取时间范围内的日志
	logs, err := s.routingLogRepo.GetByTimeRange(ctx, tenantID, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		return nil, errno.ROUTING500001.WithDetail("error", err.Error())
	}

	if len(logs) == 0 {
		return nil, errno.ROUTING404001.WithDetail("reason", "no logs found in time range")
	}

	// 2. 计算总体统计
	insight := &RoutingInsight{
		TenantID:          tenantID,
		TimeRange:         timeRange,
		TotalRequests:     len(logs),
		IntentDistribution: make(map[string]int),
		AgentPerformance:  make(map[string]*AgentPerformanceStats),
		Recommendations:   make([]string, 0),
	}

	// 3. 分析意图分布
	intentCounts := make(map[string]int)
	agentPerformanceMap := make(map[string]*agentPerfAccumulator)

	for _, log := range logs {
		// 意图统计
		if log.Intent != "" {
			intentCounts[log.Intent]++
		}

		// Agent性能统计
		if _, exists := agentPerformanceMap[log.AgentID]; !exists {
			agentPerformanceMap[log.AgentID] = &agentPerfAccumulator{
				AgentID:       log.AgentID,
				ResponseTimes: make([]int, 0),
				Ratings:       make([]int, 0),
			}
		}
		acc := agentPerformanceMap[log.AgentID]
		acc.ResponseTimes = append(acc.ResponseTimes, log.ResponseTime)
		if log.UserRating > 0 {
			acc.Ratings = append(acc.Ratings, log.UserRating)
		}
	}

	// 4. 计算平均指标
	totalTime := 0
	successCount := 0
	for _, log := range logs {
		totalTime += log.ResponseTime
		if log.UserRating >= 3 {
			successCount++
		}
	}

	insight.AvgResponseTime = float64(totalTime) / float64(len(logs))
	insight.SuccessRate = float64(successCount) / float64(len(logs))
	insight.IntentDistribution = intentCounts

	// 5. 计算Agent性能
	for agentID, acc := range agentPerformanceMap {
		insight.AgentPerformance[agentID] = s.calculateAgentPerformance(acc)
	}

	// 6. 生成推荐
	insight.Recommendations = s.generateRecommendations(insight)

	return insight, nil
}

// SuggestRoutingRules 建议路由规则
func (s *RoutingLearningService) SuggestRoutingRules(ctx context.Context, tenantID string) ([]*RoutingRuleSuggestion, error) {
	// 1. 获取最近的日志
	endTime := time.Now().Unix()
	startTime := endTime - 30*24*3600 // 最近30天

	logs, err := s.routingLogRepo.GetByTimeRange(ctx, tenantID, startTime, endTime)
	if err != nil {
		return nil, errno.ROUTING500001.WithDetail("error", err.Error())
	}

	if len(logs) == 0 {
		return []*RoutingRuleSuggestion{}, nil
	}

	// 2. 按意图分组分析
	intentAgentMap := make(map[string]*intentAgentAnalysis)

	for _, log := range logs {
		if log.Intent == "" {
			continue
		}

		if _, exists := intentAgentMap[log.Intent]; !exists {
			intentAgentMap[log.Intent] = &intentAgentAnalysis{
				Intent:      log.Intent,
				AgentScores: make(map[string]*agentScore),
			}
		}

		analysis := intentAgentMap[log.Intent]
		if _, exists := analysis.AgentScores[log.AgentID]; !exists {
			analysis.AgentScores[log.AgentID] = &agentScore{
				AgentID:       log.AgentID,
				ResponseTimes: make([]int, 0),
				Ratings:       make([]int, 0),
				Count:         0,
			}
		}

		score := analysis.AgentScores[log.AgentID]
		score.ResponseTimes = append(score.ResponseTimes, log.ResponseTime)
		if log.UserRating > 0 {
			score.Ratings = append(score.Ratings, log.UserRating)
		}
		score.Count++
	}

	// 3. 生成建议
	suggestions := make([]*RoutingRuleSuggestion, 0)

	for intent, analysis := range intentAgentMap {
		// 找出最佳Agent
		var bestAgent *agentScore
		for _, score := range analysis.AgentScores {
			if bestAgent == nil || s.calculateAgentScore(score) > s.calculateAgentScore(bestAgent) {
				bestAgent = score
			}
		}

		if bestAgent != nil && bestAgent.Count >= 10 { // 至少10个样本
			suggestions = append(suggestions, &RoutingRuleSuggestion{
				Intent:              intent,
				SuggestedAgent:      bestAgent.AgentID,
				Confidence:          s.calculateConfidence(bestAgent),
				Reason:              fmt.Sprintf("该Agent在意图'%s'上的表现最佳（%d次调用）", intent, bestAgent.Count),
				ExpectedImprovement: s.calculateExpectedImprovement(analysis, bestAgent),
			})
		}
	}

	return suggestions, nil
}

// OptimizeRoutingStrategy 优化路由策略
func (s *RoutingLearningService) OptimizeRoutingStrategy(ctx context.Context, strategyID string) (*OptimizedStrategy, error) {
	// 1. 获取历史性能数据
	// 简化实现：返回默认优化策略
	return &OptimizedStrategy{
		StrategyID:   strategyID,
		StrategyName: "optimized_strategy",
		Config: map[string]interface{}{
			"load_weight":       0.25,
			"health_weight":     0.20,
			"response_weight":   0.20,
			"error_weight":      0.15,
			"cost_weight":       0.10,
			"preference_weight": 0.10,
		},
		ExpectedMetrics: map[string]float64{
			"avg_response_time": 1500.0,
			"success_rate":      0.92,
			"user_satisfaction": 4.3,
		},
	}, nil
}

// 辅助结构和方法

type agentPerfAccumulator struct {
	AgentID       string
	ResponseTimes []int
	Ratings       []int
}

func (s *RoutingLearningService) calculateAgentPerformance(acc *agentPerfAccumulator) *AgentPerformanceStats {
	stats := &AgentPerformanceStats{
		AgentID:      acc.AgentID,
		RequestCount: len(acc.ResponseTimes),
	}

	// 平均响应时间
	totalTime := 0
	for _, t := range acc.ResponseTimes {
		totalTime += t
	}
	if len(acc.ResponseTimes) > 0 {
		stats.AvgResponseTime = float64(totalTime) / float64(len(acc.ResponseTimes))
	}

	// 平均评分
	totalRating := 0
	for _, r := range acc.Ratings {
		totalRating += r
	}
	if len(acc.Ratings) > 0 {
		stats.AvgUserRating = float64(totalRating) / float64(len(acc.Ratings))
		stats.SuccessRate = float64(len(acc.Ratings)) / float64(len(acc.ResponseTimes))
	}

	return stats
}

func (s *RoutingLearningService) generateRecommendations(insight *RoutingInsight) []string {
	recommendations := make([]string, 0)

	// 1. 响应时间推荐
	if insight.AvgResponseTime > 3000 {
		recommendations = append(recommendations, "平均响应时间较高，建议优化Agent性能或增加资源")
	}

	// 2. 成功率推荐
	if insight.SuccessRate < 0.8 {
		recommendations = append(recommendations, "成功率较低，建议调整路由策略或改进Agent配置")
	}

	// 3. Agent性能推荐
	for agentID, perf := range insight.AgentPerformance {
		if perf.ErrorRate > 0.2 {
			recommendations = append(recommendations,
				fmt.Sprintf("Agent %s 错误率较高(%.2f%%)，建议检查配置或减少负载", agentID, perf.ErrorRate*100))
		}
		if perf.AvgResponseTime > 5000 {
			recommendations = append(recommendations,
				fmt.Sprintf("Agent %s 响应时间较长(%.2fms)，建议优化", agentID, perf.AvgResponseTime))
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "路由性能良好，建议保持当前配置")
	}

	return recommendations
}

type intentAgentAnalysis struct {
	Intent      string
	AgentScores map[string]*agentScore
}

type agentScore struct {
	AgentID       string
	ResponseTimes []int
	Ratings       []int
	Count         int
}

func (s *RoutingLearningService) calculateAgentScore(score *agentScore) float64 {
	if len(score.ResponseTimes) == 0 {
		return 0.0
	}

	// 计算响应时间得分（越短越好）
	totalTime := 0
	for _, t := range score.ResponseTimes {
		totalTime += t
	}
	avgTime := float64(totalTime) / float64(len(score.ResponseTimes))
	timeScore := 1.0 / (1.0 + avgTime/1000.0)

	// 计算评分得分（越高越好）
	var avgRating float64
	if len(score.Ratings) > 0 {
		totalRating := 0
		for _, r := range score.Ratings {
			totalRating += r
		}
		avgRating = float64(totalRating) / float64(len(score.Ratings))
	}
	ratingScore := avgRating / 5.0

	// 综合得分
	return timeScore*0.4 + ratingScore*0.6
}

func (s *RoutingLearningService) calculateConfidence(score *agentScore) float64 {
	// 基于样本量计算置信度
	// 样本越多，置信度越高
	return math.Min(1.0, float64(score.Count)/100.0)
}

func (s *RoutingLearningService) calculateExpectedImprovement(analysis *intentAgentAnalysis, bestAgent *agentScore) float64 {
	// 计算相比其他Agent的改进幅度
	var totalScore float64
	var count int

	for _, score := range analysis.AgentScores {
		if score.AgentID != bestAgent.AgentID {
			totalScore += s.calculateAgentScore(score)
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	avgScore := totalScore / float64(count)
	bestScore := s.calculateAgentScore(bestAgent)

	improvement := (bestScore - avgScore) / avgScore
	return math.Max(0.0, improvement)
}
