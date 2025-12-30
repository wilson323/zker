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

// LoadMonitor 负载监控器
type LoadMonitor struct {
	mu           sync.RWMutex
	loadData     map[string]*ServiceLoadData
	maxCapacities map[string]int
}

// ServiceLoadData 服务负载数据
type ServiceLoadData struct {
	ServiceID      string
	CurrentLoad    int
	MaxCapacity    int
	LoadPercent    float64
	LastUpdateTime int64
	LoadHistory    []LoadSnapshot // 历史负载数据（用于趋势分析）
}

// LoadSnapshot 负载快照
type LoadSnapshot struct {
	Load       int
	Timestamp  int64
}

// NewLoadMonitor 创建负载监控器
func NewLoadMonitor() *LoadMonitor {
	return &LoadMonitor{
		loadData:     make(map[string]*ServiceLoadData),
		maxCapacities: make(map[string]int),
	}
}

// GetServiceLoad 获取服务负载
func (m *LoadMonitor) GetServiceLoad(ctx context.Context, serviceID string) (*ServiceLoad, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.loadData[serviceID]
	if !exists {
		// 服务不存在，返回默认负载
		return &ServiceLoad{
			ServiceID:   serviceID,
			CurrentLoad: 0,
			MaxCapacity: 100,
		}, nil
	}

	return &ServiceLoad{
		ServiceID:   data.ServiceID,
		CurrentLoad: data.CurrentLoad,
		MaxCapacity: data.MaxCapacity,
	}, nil
}

// IncrementLoad 增加负载
func (m *LoadMonitor) IncrementLoad(serviceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixMilli()

	// 初始化负载数据
	if _, exists := m.loadData[serviceID]; !exists {
		maxCapacity := 100 // 默认容量
		if cap, ok := m.maxCapacities[serviceID]; ok {
			maxCapacity = cap
		}

		m.loadData[serviceID] = &ServiceLoadData{
			ServiceID:      serviceID,
			MaxCapacity:    maxCapacity,
			LoadHistory:    make([]LoadSnapshot, 0, 100),
			LastUpdateTime: now,
		}
	}

	data := m.loadData[serviceID]
	data.CurrentLoad++
	data.LastUpdateTime = now

	// 更新负载百分比
	if data.MaxCapacity > 0 {
		data.LoadPercent = float64(data.CurrentLoad) / float64(data.MaxCapacity)
	}

	// 添加历史快照（每10秒记录一次）
	if len(data.LoadHistory) == 0 || now-data.LoadHistory[len(data.LoadHistory)-1].Timestamp > 10000 {
		data.LoadHistory = append(data.LoadHistory, LoadSnapshot{
			Load:      data.CurrentLoad,
			Timestamp: now,
		})

		// 限制历史记录数量
		if len(data.LoadHistory) > 100 {
			data.LoadHistory = data.LoadHistory[1:]
		}
	}
}

// DecrementLoad 减少负载
func (m *LoadMonitor) DecrementLoad(serviceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, exists := m.loadData[serviceID]
	if !exists {
		return
	}

	if data.CurrentLoad > 0 {
		data.CurrentLoad--
		data.LastUpdateTime = time.Now().UnixMilli()

		// 更新负载百分比
		if data.MaxCapacity > 0 {
			data.LoadPercent = float64(data.CurrentLoad) / float64(data.MaxCapacity)
		}
	}
}

// SetMaxCapacity 设置最大容量
func (m *LoadMonitor) SetMaxCapacity(serviceID string, capacity int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.maxCapacities[serviceID] = capacity

	if data, exists := m.loadData[serviceID]; exists {
		data.MaxCapacity = capacity
		if capacity > 0 {
			data.LoadPercent = float64(data.CurrentLoad) / float64(capacity)
		}
	}
}

// GetLoadTrend 获取负载趋势（最近5分钟）
func (m *LoadMonitor) GetLoadTrend(serviceID string) []LoadSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.loadData[serviceID]
	if !exists {
		return nil
	}

	// 返回最近5分钟的快照
	now := time.Now().UnixMilli()
	cutoff := now - 5*60*1000 // 5分钟

	result := make([]LoadSnapshot, 0)
	for _, snapshot := range data.LoadHistory {
		if snapshot.Timestamp >= cutoff {
			result = append(result, snapshot)
		}
	}

	return result
}

// GetLoadPercent 获取负载百分比
func (m *LoadMonitor) GetLoadPercent(serviceID string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, exists := m.loadData[serviceID]
	if !exists {
		return 0
	}

	return data.LoadPercent
}
