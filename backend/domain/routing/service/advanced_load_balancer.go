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
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// AdvancedLoadBalancer 高级负载均衡器
// 职责：提供多种高级负载均衡策略，优化路由性能
type AdvancedLoadBalancer struct {
	strategies map[string]AdvancedLoadBalanceStrategy
	monitor    *LoadMonitor
	logger     *zap.Logger
	mu         sync.RWMutex
}

// AdvancedLoadBalanceStrategy 高级负载均衡策略接口
type AdvancedLoadBalanceStrategy interface {
	// Name 策略名称
	Name() string

	// SelectAgent 选择最佳Agent
	SelectAgent(ctx context.Context, req *RoutingRequest, candidates []*CandidateAgent) (*CandidateAgent, error)

	// Validate 验证候选Agent
	Validate(ctx context.Context, candidates []*CandidateAgent) error
}

// CandidateAgent 候选Agent
type CandidateAgent struct {
	AgentID           string    `json:"agent_id"`
	BotID             string    `json:"bot_id"`
	WorkflowID        string    `json:"workflow_id"`
	Score             float64   `json:"score"`
	HealthScore       float64   `json:"health_score"`
	ActiveConnections int64     `json:"active_connections"`
	AvgResponseTime   float64   `json:"avg_response_time"` // ms
	CPUUsage          float64   `json:"cpu_usage"`          // 0-1
	MemoryUsage       float64   `json:"memory_usage"`       // 0-1
	ErrorRate         float64   `json:"error_rate"`         // 0-1
	LastUpdated       time.Time `json:"last_updated"`
	Tags              []string  `json:"tags"`
}

// NewAdvancedLoadBalancer 创建高级负载均衡器
func NewAdvancedLoadBalancer(
	monitor *LoadMonitor,
	logger *zap.Logger,
) *AdvancedLoadBalancer {
	balancer := &AdvancedLoadBalancer{
		strategies: make(map[string]AdvancedLoadBalanceStrategy),
		monitor:    monitor,
		logger:     logger,
	}

	// 注册内置策略
	balancer.RegisterStrategy(&LeastConnectionsStrategy{
		monitor: monitor,
		logger:  logger,
	})

	balancer.RegisterStrategy(&WeightedResponseTimeStrategy{
		monitor: monitor,
		logger:  logger,
	})

	balancer.RegisterStrategy(&ResourceAwareStrategy{
		monitor: monitor,
		logger:  logger,
	})

	balancer.RegisterStrategy(&PredictiveScalingStrategy{
		predictor: NewTrafficPredictor(),
		monitor:   monitor,
		logger:    logger,
	})

	balancer.RegisterStrategy(&IPHashStrategy{})

	balancer.RegisterStrategy(&LeastLatencyStrategy{
		monitor: monitor,
		logger:  logger,
	})

	return balancer
}

// RegisterStrategy 注册负载均衡策略
func (b *AdvancedLoadBalancer) RegisterStrategy(strategy AdvancedLoadBalanceStrategy) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.strategies[strategy.Name()] = strategy
}

// GetStrategy 获取负载均衡策略
func (b *AdvancedLoadBalancer) GetStrategy(name string) (AdvancedLoadBalanceStrategy, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	strategy, exists := b.strategies[name]
	if !exists {
		return nil, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("strategy", name),
			errorx.KV("reason", "strategy not found"))
	}

	return strategy, nil
}

// SelectBestAgent 使用指定策略选择最佳Agent
func (b *AdvancedLoadBalancer) SelectBestAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
	strategyName string,
) (*CandidateAgent, error) {
	if len(candidates) == 0 {
		return nil, errorx.New(errno.ErrRouteNotFoundCode,
			errorx.KV("reason", "no candidates available"))
	}

	// 获取策略
	strategy, err := b.GetStrategy(strategyName)
	if err != nil {
		// 使用默认策略
		b.logger.Warn("strategy not found, using default",
			zap.String("requested", strategyName),
			zap.String("default", "least_connections"))
		strategy, _ = b.GetStrategy("least_connections")
	}

	// 验证候选Agent
	if err := strategy.Validate(ctx, candidates); err != nil {
		return nil, err
	}

	// 选择最佳Agent
	agent, err := strategy.SelectAgent(ctx, req, candidates)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrRoutingDecisionFailedCode)
	}

	return agent, nil
}

// LeastConnectionsStrategy 最少连接策略
// 选择当前活动连接数最少的Agent
type LeastConnectionsStrategy struct {
	monitor *LoadMonitor
	logger  *zap.Logger
}

func (s *LeastConnectionsStrategy) Name() string {
	return "least_connections"
}

func (s *LeastConnectionsStrategy) SelectAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
) (*CandidateAgent, error) {
	var selected *CandidateAgent
	minConns := int64(1<<63 - 1) // MaxInt64

	for _, agent := range candidates {
		// 优先使用实时监控数据
		metrics := s.monitor.GetRealTimeMetrics(agent.AgentID)
		if metrics != nil {
			if metrics.ActiveConnections < minConns {
				minConns = metrics.ActiveConnections
				selected = agent
			}
		} else {
			// 回退到候选数据
			if agent.ActiveConnections < minConns {
				minConns = agent.ActiveConnections
				selected = agent
			}
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("failed to select agent")
	}

	return selected, nil
}

func (s *LeastConnectionsStrategy) Validate(
	ctx context.Context,
	candidates []*CandidateAgent,
) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates")
	}
	return nil
}

// WeightedResponseTimeStrategy 加权响应时间策略
// 响应时间越短，权重越高
type WeightedResponseTimeStrategy struct {
	monitor *LoadMonitor
	logger  *zap.Logger
}

func (s *WeightedResponseTimeStrategy) Name() string {
	return "weighted_response_time"
}

type scoredAgent struct {
	agent *CandidateAgent
	score float64
}

func (s *WeightedResponseTimeStrategy) SelectAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
) (*CandidateAgent, error) {
	scoredAgents := make([]scoredAgent, len(candidates))

	// 计算每个Agent的得分
	for i, agent := range candidates {
		// 获取实时指标
		metrics := s.monitor.GetRealTimeMetrics(agent.AgentID)
		responseTime := agent.AvgResponseTime

		if metrics != nil && metrics.AvgResponseTime > 0 {
			responseTime = metrics.AvgResponseTime
		}

		// 分数 = 1 / 响应时间（响应时间越短，分数越高）
		// 避免除以0，添加最小值0.001
		score := 1.0 / (responseTime/1000.0 + 0.001)

		scoredAgents[i] = scoredAgent{
			agent: agent,
			score: score,
		}
	}

	// 按得分排序
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].score > scoredAgents[j].score
	})

	// 返回得分最高的
	return scoredAgents[0].agent, nil
}

func (s *WeightedResponseTimeStrategy) Validate(
	ctx context.Context,
	candidates []*CandidateAgent,
) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates")
	}
	return nil
}

// ResourceAwareStrategy 资源感知策略
// 综合考虑CPU、内存、网络等资源使用情况
type ResourceAwareStrategy struct {
	monitor *LoadMonitor
	logger  *zap.Logger
}

func (s *ResourceAwareStrategy) Name() string {
	return "resource_aware"
}

func (s *ResourceAwareStrategy) SelectAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
) (*CandidateAgent, error) {
	var selected *CandidateAgent
	bestScore := -1.0

	for _, agent := range candidates {
		// 获取实时指标
		metrics := s.monitor.GetRealTimeMetrics(agent.AgentID)

		var cpuUsage, memUsage float64
		var activeConns int64

		if metrics != nil {
			cpuUsage = metrics.CPUUsage
			memUsage = metrics.MemoryUsage
			activeConns = metrics.ActiveConnections
		} else {
			cpuUsage = agent.CPUUsage
			memUsage = agent.MemoryUsage
			activeConns = agent.ActiveConnections
		}

		// 综合评分（资源使用率越低，分数越高）
		cpuScore := 1.0 - cpuUsage
		memScore := 1.0 - memUsage
		connScore := 1.0 - math.Min(float64(activeConns)/1000.0, 1.0)

		// 加权平均
		totalScore := (cpuScore*0.4 + memScore*0.3 + connScore*0.3)

		// 考虑健康状态
		healthScore := agent.HealthScore
		if healthScore > 0 {
			totalScore = totalScore * healthScore
		}

		if totalScore > bestScore {
			bestScore = totalScore
			selected = agent
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("failed to select agent")
	}

	return selected, nil
}

func (s *ResourceAwareStrategy) Validate(
	ctx context.Context,
	candidates []*CandidateAgent,
) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates")
	}
	return nil
}

// PredictiveScalingStrategy 预测扩容策略
// 基于流量预测动态选择Agent
type PredictiveScalingStrategy struct {
	predictor *TrafficPredictor
	monitor   *LoadMonitor
	logger    *zap.Logger
}

func (s *PredictiveScalingStrategy) Name() string {
	return "predictive_scaling"
}

func (s *PredictiveScalingStrategy) SelectAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
) (*CandidateAgent, error) {
	// 预测未来5分钟的负载
	predictedLoad := s.predictor.Predict(ctx, 5*time.Minute)

	s.logger.Debug("traffic prediction",
		zap.Float64("predicted_load", predictedLoad),
		zap.Time("prediction_time", time.Now()))

	// 如果预测高负载，选择资源最充足的Agent
	if predictedLoad > 0.8 {
		s.logger.Info("high load predicted, using resource-aware selection",
			zap.Float64("load", predictedLoad))
		strategy := &ResourceAwareStrategy{monitor: s.monitor, logger: s.logger}
		return strategy.SelectAgent(ctx, req, candidates)
	}

	// 否则使用加权响应时间策略
	strategy := &WeightedResponseTimeStrategy{monitor: s.monitor, logger: s.logger}
	return strategy.SelectAgent(ctx, req, candidates)
}

func (s *PredictiveScalingStrategy) Validate(
	ctx context.Context,
	candidates []*CandidateAgent,
) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates")
	}
	return nil
}

// IPHashStrategy IP哈希策略
// 基于用户IP的哈希值选择Agent，确保同一用户总是路由到同一Agent
type IPHashStrategy struct{}

func (s *IPHashStrategy) Name() string {
	return "ip_hash"
}

func (s *IPHashStrategy) SelectAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
) (*CandidateAgent, error) {
	if req.ClientIP == "" {
		// 如果没有IP信息，使用轮询
		return candidates[0], nil
	}

	// 计算IP哈希
	hash := simpleHash(req.ClientIP)
	index := int(hash) % len(candidates)

	return candidates[index], nil
}

func (s *IPHashStrategy) Validate(
	ctx context.Context,
	candidates []*CandidateAgent,
) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates")
	}
	return nil
}

// simpleHash 简单哈希函数
func simpleHash(s string) uint32 {
	h := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(s); i++ {
		h *= prime32
		h ^= uint32(s[i])
	}
	return h
}

// LeastLatencyStrategy 最低延迟策略
// 选择网络延迟最低的Agent
type LeastLatencyStrategy struct {
	monitor *LoadMonitor
	logger  *zap.Logger
}

func (s *LeastLatencyStrategy) Name() string {
	return "least_latency"
}

func (s *LeastLatencyStrategy) SelectAgent(
	ctx context.Context,
	req *RoutingRequest,
	candidates []*CandidateAgent,
) (*CandidateAgent, error) {
	var selected *CandidateAgent
	minLatency := math.MaxFloat64

	for _, agent := range candidates {
		// 获取实时延迟数据
		metrics := s.monitor.GetRealTimeMetrics(agent.AgentID)

		var latency float64
		if metrics != nil && metrics.Latency > 0 {
			latency = metrics.Latency
		} else {
			// 回退到响应时间估算延迟
			latency = agent.AvgResponseTime
		}

		if latency < minLatency {
			minLatency = latency
			selected = agent
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("failed to select agent")
	}

	return selected, nil
}

func (s *LeastLatencyStrategy) Validate(
	ctx context.Context,
	candidates []*CandidateAgent,
) error {
	if len(candidates) == 0 {
		return fmt.Errorf("no candidates")
	}
	return nil
}

// TrafficPredictor 流量预测器
type TrafficPredictor struct {
	historicalData []TrafficDataPoint
	mu             sync.RWMutex
}

// TrafficDataPoint 流量数据点
type TrafficDataPoint struct {
	Timestamp time.Time
	Load      float64 // 0-1
}

// NewTrafficPredictor 创建流量预测器
func NewTrafficPredictor() *TrafficPredictor {
	return &TrafficPredictor{
		historicalData: make([]TrafficDataPoint, 0),
	}
}

// Predict 预测未来负载
func (p *TrafficPredictor) Predict(ctx context.Context, duration time.Duration) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.historicalData) < 2 {
		// 数据不足，返回默认值
		return 0.5
	}

	// 简单线性预测
	// 实际应用中应使用更复杂的算法（如ARIMA、LSTM等）
	last := p.historicalData[len(p.historicalData)-1]
	secondLast := p.historicalData[len(p.historicalData)-2]

	// 计算变化率
	changeRate := last.Load - secondLast.Load

	// 预测未来负载
	predictedLoad := last.Load + changeRate*float64(duration.Minutes())/5.0

	// 限制在0-1范围
	if predictedLoad < 0 {
		predictedLoad = 0
	} else if predictedLoad > 1 {
		predictedLoad = 1
	}

	return predictedLoad
}

// RecordTraffic 记录流量数据
func (p *TrafficPredictor) RecordTraffic(load float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	dataPoint := TrafficDataPoint{
		Timestamp: time.Now(),
		Load:      load,
	}

	p.historicalData = append(p.historicalData, dataPoint)

	// 保留最近1000个数据点
	if len(p.historicalData) > 1000 {
		p.historicalData = p.historicalData[len(p.historicalData)-1000:]
	}
}
