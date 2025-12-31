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
	"sort"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// ScoringEnhancedService 增强评分服务
// 职责：提供多维度路由评分能力，包括语义相似度、功能匹配度、性能评分等
type ScoringEnhancedService struct {
	intentMatcher  IntentMatcher
	serviceRegistry ServiceRegistry
	loadMonitor    LoadMonitor
	logger         *zap.Logger

	// 评分权重配置
	semanticWeight     float64 // 语义相似度权重
	functionalWeight   float64 // 功能匹配度权重
	performanceWeight  float64 // 性能评分权重
	costWeight         float64 // 成本权重
	availabilityWeight float64 // 可用性权重

	// Top-K配置
	topK              int // Top-K选择数量
	minScoreThreshold float64
}

// NewScoringEnhancedService 创建增强评分服务实例
func NewScoringEnhancedService(
	intentMatcher IntentMatcher,
	serviceRegistry ServiceRegistry,
	loadMonitor LoadMonitor,
	logger *zap.Logger,
) *ScoringEnhancedService {
	return &ScoringEnhancedService{
		intentMatcher:       intentMatcher,
		serviceRegistry:     serviceRegistry,
		loadMonitor:         loadMonitor,
		logger:              logger,
		semanticWeight:      0.25, // 25%
		functionalWeight:    0.20, // 20%
		performanceWeight:   0.20, // 20%
		costWeight:          0.15, // 15%
		availabilityWeight:  0.20, // 20%
		topK:                5,
		minScoreThreshold:   0.5,
	}
}

// SetWeights 设置评分权重
func (s *ScoringEnhancedService) SetWeights(
	semantic, functional, performance, cost, availability float64,
) {
	s.semanticWeight = semantic
	s.functionalWeight = functional
	s.performanceWeight = performance
	s.costWeight = cost
	s.availabilityWeight = availability
}

// ScoreAndSelect 评分并选择最佳路由
func (s *ScoringEnhancedService) ScoreAndSelect(
	ctx context.Context,
	input *MatchInput,
) (*entity.RoutingDecision, []*ScoredCandidate, error) {
	// 1. 意图匹配获取候选列表
	candidates, err := s.intentMatcher.Match(ctx, input)
	if err != nil {
		return nil, nil, errorx.WrapByCode(err, errno.ErrRouteMatchFailedCode, errorx.KV("error", err.Error()))
	}

	if len(candidates) == 0 {
		return nil, nil, errorx.New(errno.ErrRouteNotFoundCode, errorx.KV("reason", "no candidates"))
	}

	// 2. 为每个候选计算多维度得分
	scoredCandidates := make([]*ScoredCandidate, 0)
	for _, candidate := range candidates {
		scored, err := s.scoreCandidate(ctx, candidate, input)
		if err != nil {
			s.logger.Warn("failed to score candidate",
				zap.String("bot_id", candidate.BotID),
				zap.Error(err))
			continue
		}
		scoredCandidates = append(scoredCandidates, scored)
	}

	if len(scoredCandidates) == 0 {
		return nil, nil, errorx.New(errno.ErrRouteScoreFailedCode, errorx.KV("reason", "no valid candidates"))
	}

	// 3. 按综合得分排序
	sort.Slice(scoredCandidates, func(i, j int) bool {
		return scoredCandidates[i].TotalScore > scoredCandidates[j].TotalScore
	})

	// 4. 选择最佳候选
	bestCandidate := scoredCandidates[0]

	// 5. 构建路由决策
	decision := &entity.RoutingDecision{
		AgentID:      bestCandidate.BotID,
		WorkflowID:   bestCandidate.WorkflowID,
		Confidence:   bestCandidate.Confidence,
		Score:        bestCandidate.TotalScore,
		MatchType:    bestCandidate.MatchType,
		RuleID:       bestCandidate.RuleID,
		Strategy:     "scoring_enhanced",
		Entities:     []*entity.Entity{}, // MatchOutput没有Entities字段，使用空切片
	}

	return decision, scoredCandidates, nil
}

// SelectTopK 选择Top-K候选
func (s *ScoringEnhancedService) SelectTopK(
	ctx context.Context,
	input *MatchInput,
	k int,
) ([]*entity.RoutingDecision, error) {
	if k <= 0 {
		k = s.topK
	}

	// 获取评分结果
	_, scoredCandidates, err := s.ScoreAndSelect(ctx, input)
	if err != nil {
		return nil, err
	}

	// 限制返回数量
	if len(scoredCandidates) > k {
		scoredCandidates = scoredCandidates[:k]
	}

	// 转换为路由决策列表
	decisions := make([]*entity.RoutingDecision, len(scoredCandidates))
	for i, candidate := range scoredCandidates {
		decisions[i] = &entity.RoutingDecision{
			AgentID:      candidate.BotID,
			WorkflowID:   candidate.WorkflowID,
			Confidence:   candidate.Confidence,
			Score:        candidate.TotalScore,
			MatchType:    candidate.MatchType,
			RuleID:       candidate.RuleID,
			Strategy:     "scoring_enhanced_topk",
		}
	}

	return decisions, nil
}

// scoreCandidate 评分单个候选
func (s *ScoringEnhancedService) scoreCandidate(
	ctx context.Context,
	candidate *MatchOutput,
	input *MatchInput,
) (*ScoredCandidate, error) {
	scored := &ScoredCandidate{
		BotID:      candidate.BotID,
		WorkflowID: candidate.WorkflowID,
		Confidence: candidate.Confidence,
		MatchType:  candidate.MatchType,
		RuleID:     candidate.RuleID,
		Entities:   []*entity.Entity{}, // MatchOutput没有Entities字段，使用空切片
	}

	// 1. 语义相似度得分
	scored.SemanticScore = s.calculateSemanticScore(candidate, input)

	// 2. 功能匹配度得分
	scored.FunctionalScore = s.calculateFunctionalScore(ctx, candidate)

	// 3. 性能评分
	scored.PerformanceScore = s.calculatePerformanceScore(ctx, candidate)

	// 4. 成本评分
	scored.CostScore = s.calculateCostScore(ctx, candidate)

	// 5. 可用性评分
	scored.AvailabilityScore = s.calculateAvailabilityScore(ctx, candidate)

	// 6. 计算综合得分
	scored.TotalScore =
		scored.SemanticScore*s.semanticWeight +
			scored.FunctionalScore*s.functionalWeight +
			scored.PerformanceScore*s.performanceWeight +
			scored.CostScore*s.costWeight +
			scored.AvailabilityScore*s.availabilityWeight

	return scored, nil
}

// calculateSemanticScore 计算语义相似度得分
func (s *ScoringEnhancedService) calculateSemanticScore(
	candidate *MatchOutput,
	input *MatchInput,
) float64 {
	// 基础得分：意图置信度
	baseScore := candidate.Confidence

	// MatchOutput没有Entities字段，无法计算实体加分
	// 未来如果需要实体信息，需要从input中提取

	score := baseScore
	if score > 1.0 {
		score = 1.0
	}

	return score
}

// calculateFunctionalScore 计算功能匹配度得分
func (s *ScoringEnhancedService) calculateFunctionalScore(
	ctx context.Context,
	candidate *MatchOutput,
) float64 {
	// 简化实现：基于规则ID和匹配类型
	// 实际应用中应该根据Agent/Workflow的功能标签计算

	score := 0.7 // 基础分

	// 根据匹配类型调整
	switch candidate.MatchType {
	case "exact":
		score = 1.0
	case "semantic":
		score = 0.9
	case "fuzzy":
		score = 0.7
	case "rule":
		score = 0.8
	default:
		score = 0.6
	}

	return score
}

// calculatePerformanceScore 计算性能评分
func (s *ScoringEnhancedService) calculatePerformanceScore(
	ctx context.Context,
	candidate *MatchOutput,
) float64 {
	// 获取服务健康状态
	var health *ServiceHealth
	var err error

	if candidate.BotID != "" {
		health, err = s.serviceRegistry.GetBotHealth(ctx, candidate.BotID)
	} else if candidate.WorkflowID != "" {
		health, err = s.serviceRegistry.GetWorkflowHealth(ctx, candidate.WorkflowID)
	}

	if err != nil || health == nil {
		// 无法获取健康状态，返回中等得分
		return 0.6
	}

	// 综合成功率、平均延迟计算得分
	successRateScore := health.SuccessRate
	latencyScore := s.calculateLatencyScore(health.AvgLatency)

	performanceScore := (successRateScore*0.6 + latencyScore*0.4)
	return performanceScore
}

// calculateLatencyScore 计算延迟得分
func (s *ScoringEnhancedService) calculateLatencyScore(latency time.Duration) float64 {
	// 延迟越低得分越高
	// <100ms: 1.0, 100-500ms: 0.8-1.0, 500-1000ms: 0.6-0.8, >1000ms: 0.4-0.6
	latencyMs := float64(latency.Milliseconds())

	switch {
	case latencyMs < 100:
		return 1.0
	case latencyMs < 500:
		return 1.0 - (latencyMs-100)/400*0.2
	case latencyMs < 1000:
		return 0.8 - (latencyMs-500)/500*0.2
	default:
		return 0.4
	}
}

// calculateCostScore 计算成本评分
func (s *ScoringEnhancedService) calculateCostScore(
	ctx context.Context,
	candidate *MatchOutput,
) float64 {
	// 简化实现：基于Agent/Workflow的成本等级
	// 实际应用中应该从配置或数据库读取成本信息

	// 默认给予中等成本得分
	return 0.7
}

// calculateAvailabilityScore 计算可用性评分
func (s *ScoringEnhancedService) calculateAvailabilityScore(
	ctx context.Context,
	candidate *MatchOutput,
) float64 {
	// 获取负载信息
	var serviceID string
	if candidate.BotID != "" {
		serviceID = candidate.BotID
	} else if candidate.WorkflowID != "" {
		serviceID = candidate.WorkflowID
	}

	load, err := s.loadMonitor.GetServiceLoad(ctx, serviceID)
	if err != nil || load == nil {
		// 无法获取负载信息，返回中等得分
		return 0.6
	}

	// 根据负载比例计算得分
	if load.MaxCapacity == 0 {
		return 0.5
	}

	loadRatio := float64(load.CurrentLoad) / float64(load.MaxCapacity)

	// 负载越低得分越高
	// <50%: 1.0, 50-80%: 0.7-1.0, 80-95%: 0.4-0.7, >95%: 0.2
	switch {
	case loadRatio < 0.5:
		return 1.0
	case loadRatio < 0.8:
		return 1.0 - (loadRatio-0.5)/0.3*0.3
	case loadRatio < 0.95:
		return 0.7 - (loadRatio-0.8)/0.15*0.3
	default:
		return 0.2
	}
}

// ScoredCandidate 评分后的候选
type ScoredCandidate struct {
	BotID             string     `json:"bot_id"`
	WorkflowID        string     `json:"workflow_id"`
	Confidence        float64    `json:"confidence"`
	SemanticScore     float64    `json:"semantic_score"`
	FunctionalScore   float64    `json:"functional_score"`
	PerformanceScore  float64    `json:"performance_score"`
	CostScore         float64    `json:"cost_score"`
	AvailabilityScore float64    `json:"availability_score"`
	TotalScore        float64    `json:"total_score"`
	MatchType         string     `json:"match_type"`
	RuleID            string     `json:"rule_id"`
	Entities          []*entity.Entity `json:"entities"`
}

// GetScoreDetails 获取评分详情
func (c *ScoredCandidate) GetScoreDetails() map[string]float64 {
	return map[string]float64{
		"semantic":     c.SemanticScore,
		"functional":   c.FunctionalScore,
		"performance":  c.PerformanceScore,
		"cost":         c.CostScore,
		"availability": c.AvailabilityScore,
		"total":        c.TotalScore,
	}
}

// NormalizeScores 归一化得分列表
func NormalizeScores(candidates []*ScoredCandidate) {
	if len(candidates) == 0 {
		return
	}

	// 找出最大和最小得分
	maxScore := candidates[0].TotalScore
	minScore := candidates[0].TotalScore

	for _, c := range candidates {
		if c.TotalScore > maxScore {
			maxScore = c.TotalScore
		}
		if c.TotalScore < minScore {
			minScore = c.TotalScore
		}
	}

	// 归一化到0-1范围
	if maxScore != minScore {
		for _, c := range candidates {
			c.TotalScore = (c.TotalScore - minScore) / (maxScore - minScore)
		}
	} else {
		// 所有得分相同，设为1.0
		for _, c := range candidates {
			c.TotalScore = 1.0
		}
	}
}

// CalculateSoftmax 计算Softmax概率分布
func CalculateSoftmax(candidates []*ScoredCandidate) []float64 {
	scores := make([]float64, len(candidates))
	for i, c := range candidates {
		scores[i] = c.TotalScore
	}

	// 减去最大值以提高数值稳定性
	maxScore := scores[0]
	for _, s := range scores {
		if s > maxScore {
			maxScore = s
		}
	}

	// 计算exp和sum
	expScores := make([]float64, len(scores))
	sum := 0.0
	for i, s := range scores {
		expScores[i] = math.Exp(s - maxScore)
		sum += expScores[i]
	}

	// 归一化
	probs := make([]float64, len(scores))
	for i, exp := range expScores {
		probs[i] = exp / sum
	}

	return probs
}

// ExplainScoring 解释评分结果
func (s *ScoringEnhancedService) ExplainScoring(candidate *ScoredCandidate) string {
	explanation := fmt.Sprintf("总得分: %.2f\n", candidate.TotalScore)
	explanation += fmt.Sprintf("- 语义相似度: %.2f (权重%.0f%%)\n", candidate.SemanticScore, s.semanticWeight*100)
	explanation += fmt.Sprintf("- 功能匹配度: %.2f (权重%.0f%%)\n", candidate.FunctionalScore, s.functionalWeight*100)
	explanation += fmt.Sprintf("- 性能评分: %.2f (权重%.0f%%)\n", candidate.PerformanceScore, s.performanceWeight*100)
	explanation += fmt.Sprintf("- 成本评分: %.2f (权重%.0f%%)\n", candidate.CostScore, s.costWeight*100)
	explanation += fmt.Sprintf("- 可用性评分: %.2f (权重%.0f%%)\n", candidate.AvailabilityScore, s.availabilityWeight*100)

	return explanation
}
