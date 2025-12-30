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
	"sync"
	"time"
)

// ServiceHealthMonitor 服务健康监控器
type ServiceHealthMonitor struct {
	mu               sync.RWMutex
	healthData       map[string]*ServiceHealthStats
	successWindows   map[string]*SuccessWindow
	failureWindows   map[string]*FailureWindow
}

// ServiceHealthStats 服务健康统计数据
type ServiceHealthStats struct {
	ServiceID       string
	TotalRequests   int64
	SuccessRequests int64
	FailureRequests int64
	SuccessRate     float64
	AvgLatency      time.Duration
	P95Latency      time.Duration
	P99Latency      time.Duration
	LastError       string
	LastErrorTime   int64
	LastUpdateTime  int64
	IsHealthy       bool
}

// SuccessWindow 成功请求窗口（用于计算延迟）
type SuccessWindow struct {
	latencies []time.Duration
	maxSize   int
}

// FailureWindow 失败请求窗口
type FailureWindow struct {
	failures    []time.Time
	maxSize     int
	timeWindow  time.Duration
}

// NewServiceHealthMonitor 创建服务健康监控器
func NewServiceHealthMonitor() *ServiceHealthMonitor {
	return &ServiceHealthMonitor{
		healthData:     make(map[string]*ServiceHealthStats),
		successWindows: make(map[string]*SuccessWindow),
		failureWindows:  make(map[string]*FailureWindow),
	}
}

// GetBotHealth 获取Bot健康状态
func (m *ServiceHealthMonitor) GetBotHealth(ctx context.Context, botID string) (*ServiceHealth, error) {
	return m.getServiceHealth("bot:" + botID)
}

// GetWorkflowHealth 获取工作流健康状态
func (m *ServiceHealthMonitor) GetWorkflowHealth(ctx context.Context, workflowID string) (*ServiceHealth, error) {
	return m.getServiceHealth("workflow:" + workflowID)
}

// getServiceHealth 获取服务健康状态
func (m *ServiceHealthMonitor) getServiceHealth(serviceID string) (*ServiceHealth, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats, exists := m.healthData[serviceID]
	if !exists {
		// 服务不存在，返回默认健康状态
		return &ServiceHealth{
			ServiceID:  serviceID,
			IsHealthy:  true,
			SuccessRate: 1.0,
			AvgLatency:  0,
		}, nil
	}

	// 判断是否健康（成功率 > 0.8 且 平均延迟 < 5s）
	isHealthy := stats.SuccessRate > 0.8 && stats.AvgLatency < 5*time.Second

	return &ServiceHealth{
		ServiceID:   serviceID,
		IsHealthy:   isHealthy,
		SuccessRate: stats.SuccessRate,
		AvgLatency:  stats.AvgLatency,
	}, nil
}

// RecordSuccess 记录成功请求
func (m *ServiceHealthMonitor) RecordSuccess(serviceID string, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixMilli()

	// 初始化统计数据
	if _, exists := m.healthData[serviceID]; !exists {
		m.healthData[serviceID] = &ServiceHealthStats{
			ServiceID:      serviceID,
			IsHealthy:      true,
			LastUpdateTime: now,
		}
		m.successWindows[serviceID] = &SuccessWindow{
			latencies: make([]time.Duration, 0, 1000),
			maxSize:   1000,
		}
		m.failureWindows[serviceID] = &FailureWindow{
			failures:   make([]time.Time, 0, 100),
			maxSize:    100,
			timeWindow: 5 * time.Minute,
		}
	}

	stats := m.healthData[serviceID]
	stats.TotalRequests++
	stats.SuccessRequests++
	stats.LastUpdateTime = now

	// 更新延迟窗口
	window := m.successWindows[serviceID]
	window.latencies = append(window.latencies, latency)
	if len(window.latencies) > window.maxSize {
		window.latencies = window.latencies[1:]
	}

	// 计算延迟统计
	stats.AvgLatency = m.calculateAvgLatency(window.latencies)
	stats.P95Latency = m.calculatePercentile(window.latencies, 0.95)
	stats.P99Latency = m.calculatePercentile(window.latencies, 0.99)

	// 更新成功率
	stats.SuccessRate = float64(stats.SuccessRequests) / float64(stats.TotalRequests)
	stats.IsHealthy = stats.SuccessRate > 0.8 && stats.AvgLatency < 5*time.Second
}

// RecordFailure 记录失败请求
func (m *ServiceHealthMonitor) RecordFailure(serviceID string, errMsg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixMilli()

	// 初始化统计数据
	if _, exists := m.healthData[serviceID]; !exists {
		m.healthData[serviceID] = &ServiceHealthStats{
			ServiceID:      serviceID,
			IsHealthy:      true,
			LastUpdateTime: now,
		}
		m.failureWindows[serviceID] = &FailureWindow{
			failures:   make([]time.Time, 0, 100),
			maxSize:    100,
			timeWindow: 5 * time.Minute,
		}
	}

	stats := m.healthData[serviceID]
	stats.TotalRequests++
	stats.FailureRequests++
	stats.LastError = errMsg
	stats.LastErrorTime = now
	stats.LastUpdateTime = now

	// 更新失败窗口
	window := m.failureWindows[serviceID]
	nowTime := time.Now()
	window.failures = append(window.failures, nowTime)

	// 清理过期的失败记录
	cutoff := nowTime.Add(-window.timeWindow)
	validFailures := make([]time.Time, 0)
	for _, f := range window.failures {
		if f.After(cutoff) {
			validFailures = append(validFailures, f)
		}
	}
	window.failures = validFailures

	// 更新成功率
	stats.SuccessRate = float64(stats.SuccessRequests) / float64(stats.TotalRequests)
	stats.IsHealthy = stats.SuccessRate > 0.8 && stats.AvgLatency < 5*time.Second
}

// calculateAvgLatency 计算平均延迟
func (m *ServiceHealthMonitor) calculateAvgLatency(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var sum time.Duration
	for _, lat := range latencies {
		sum += lat
	}
	return sum / time.Duration(len(latencies))
}

// calculatePercentile 计算百分位数
func (m *ServiceHealthMonitor) calculatePercentile(latencies []time.Duration, percentile float64) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	// 复制并排序
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)

	// 简单冒泡排序
	for i := 0; i < len(sorted)-1; i++ {
		for j := 0; j < len(sorted)-i-1; j++ {
			if sorted[j] > sorted[j+1] {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}

	index := int(float64(len(sorted)) * percentile)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// GetHealthStats 获取健康统计信息
func (m *ServiceHealthMonitor) GetHealthStats(serviceID string) (*ServiceHealthStats, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats, exists := m.healthData[serviceID]
	return stats, exists
}
