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
)

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	strategy     string // "round_robin", "least_loaded", "random"
	threshold    int    // 负载阈值（0-100）
	roundRobinIdx map[string]int // 每个服务集群的轮询索引
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer(strategy string, threshold int) *LoadBalancer {
	return &LoadBalancer{
		strategy:     strategy,
		threshold:    threshold,
		roundRobinIdx: make(map[string]int),
	}
}

// SelectBestCandidate 选择最佳候选服务
func (lb *LoadBalancer) SelectBestCandidate(
	ctx context.Context,
	candidates []*RoutingDecision,
) (*RoutingDecision, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates available")
	}

	// 过滤掉高负载的服务
	filtered := lb.filterByLoad(candidates)
	if len(filtered) == 0 {
		// 如果所有服务都高负载，返回得分最高的
		return candidates[0], nil
	}

	switch lb.strategy {
	case "least_loaded":
		return lb.selectLeastLoaded(ctx, filtered)
	case "round_robin":
		return lb.selectRoundRobin(ctx, filtered)
	case "random":
		return lb.selectRandom(ctx, filtered)
	default:
		// 默认使用 least_loaded
		return lb.selectLeastLoaded(ctx, filtered)
	}
}

// filterByLoad 过滤高负载服务
func (lb *LoadBalancer) filterByLoad(candidates []*RoutingDecision) []*RoutingDecision {
	filtered := make([]*RoutingDecision, 0)

	for _, candidate := range candidates {
		// 检查负载
		serviceID := candidate.BotID
		if serviceID == "" {
			serviceID = candidate.WorkflowID
		}

		// 这里应该从 LoadMonitor 获取实际负载
		// 暂时使用简单的启发式规则
		if lb.isServiceAvailable(serviceID) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}

// isServiceAvailable 检查服务是否可用
func (lb *LoadBalancer) isServiceAvailable(serviceID string) bool {
	// TODO: 从 LoadMonitor 获取实际负载
	// 暂时返回 true
	return true
}

// selectLeastLoaded 选择负载最低的服务
func (lb *LoadBalancer) selectLeastLoaded(
	ctx context.Context,
	candidates []*RoutingDecision,
) (*RoutingDecision, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates after filtering")
	}

	// 按得分排序（已经按得分排序了，直接取第一个）
	return candidates[0], nil
}

// selectRoundRobin 轮询选择
func (lb *LoadBalancer) selectRoundRobin(
	ctx context.Context,
	candidates []*RoutingDecision,
) (*RoutingDecision, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates after filtering")
	}

	// 使用服务集群标识
	clusterKey := lb.getClusterKey(candidates[0])

	// 获取并更新索引
	idx := lb.roundRobinIdx[clusterKey]
	if idx >= len(candidates) {
		idx = 0
	}
	lb.roundRobinIdx[clusterKey] = idx + 1

	return candidates[idx], nil
}

// selectRandom 随机选择
func (lb *LoadBalancer) selectRandom(
	ctx context.Context,
	candidates []*RoutingDecision,
) (*RoutingDecision, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no candidates after filtering")
	}

	// 简单随机选择：基于时间戳的伪随机
	now := ctx.Value("timestamp")
	if now == nil {
		return candidates[0], nil
	}

	// 使用分数作为权重进行加权随机选择
	totalScore := 0.0
	for _, c := range candidates {
		totalScore += c.Score
	}

	if totalScore == 0 {
		return candidates[0], nil
	}

	// 加权随机选择
	r := float64(time.Now().UnixNano()%100) / 100.0
	accumulated := 0.0
	for _, c := range candidates {
		weight := c.Score / totalScore
		accumulated += weight
		if r <= accumulated {
			return c, nil
		}
	}

	return candidates[len(candidates)-1], nil
}

// getClusterKey 获取服务集群标识
func (lb *LoadBalancer) getClusterKey(decision *RoutingDecision) string {
	// 使用BotID或WorkflowID作为集群标识
	if decision.BotID != "" {
		return "bot_cluster"
	}
	return "workflow_cluster"
}

// SetStrategy 设置负载均衡策略
func (lb *LoadBalancer) SetStrategy(strategy string) {
	lb.strategy = strategy
}

// SetThreshold 设置负载阈值
func (lb *LoadBalancer) SetThreshold(threshold int) {
	lb.threshold = threshold
}

// SortByScore 按得分排序候选
func (lb *LoadBalancer) SortByScore(candidates []*RoutingDecision) []*RoutingDecision {
	sorted := make([]*RoutingDecision, len(candidates))
	copy(sorted, candidates)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})

	return sorted
}
