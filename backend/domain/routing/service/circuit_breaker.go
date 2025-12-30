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

// CircuitBreakerService 熔断器服务
type CircuitBreakerService struct {
	mu          sync.RWMutex
	breakers    map[string]*CircuitBreaker
}

// CircuitBreaker 熔断器状态
type CircuitBreaker struct {
	mu              sync.Mutex
	ServiceID       string
	State           CircuitState // CLOSED, OPEN, HALF_OPEN
	FailureCount    int
	SuccessCount    int
	LastFailureTime int64
	LastStateChange int64

	// 配置
	FailureThreshold int   // 失败阈值
	SuccessThreshold int   // 成功阈值（半开状态）
	Timeout          int64 // 熔断超时时间（毫秒）
	HalfOpenCalls    int   // 半开状态允许的调用次数
}

// CircuitState 熔断器状态
type CircuitState string

const (
	CircuitClosed   CircuitState = "CLOSED"   // 正常状态
	CircuitOpen     CircuitState = "OPEN"     // 熔断状态
	CircuitHalfOpen CircuitState = "HALF_OPEN" // 半开状态
)

// NewCircuitBreakerService 创建熔断器服务
func NewCircuitBreakerService() *CircuitBreakerService {
	return &CircuitBreakerService{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// getOrCreateBreaker 获取或创建熔断器
func (s *CircuitBreakerService) getOrCreateBreaker(serviceID string) *CircuitBreaker {
	s.mu.Lock()
	defer s.mu.Unlock()

	if breaker, exists := s.breakers[serviceID]; exists {
		return breaker
	}

	breaker := &CircuitBreaker{
		ServiceID:       serviceID,
		State:           CircuitClosed,
		FailureThreshold: 5,  // 默认5次失败后熔断
		SuccessThreshold: 2,  // 半开状态需要2次成功
		Timeout:         60000, // 默认60秒超时
		HalfOpenCalls:    3,  // 半开状态允许3次调用
	}

	s.breakers[serviceID] = breaker
	return breaker
}

// AllowRequest 检查是否允许请求
func (s *CircuitBreakerService) AllowRequest(ctx context.Context, serviceID string) bool {
	breaker := s.getOrCreateBreaker(serviceID)
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	now := time.Now().UnixMilli()

	switch breaker.State {
	case CircuitClosed:
		// 正常状态，允许请求
		return true

	case CircuitOpen:
		// 熔断状态，检查是否可以半开
		if now-breaker.LastStateChange > breaker.Timeout {
			// 转换到半开状态
			breaker.State = CircuitHalfOpen
			breaker.LastStateChange = now
			breaker.SuccessCount = 0
			return true
		}
		// 仍在熔断中，拒绝请求
		return false

	case CircuitHalfOpen:
		// 半开状态，允许有限次数的请求
		return true

	default:
		return true
	}
}

// RecordSuccess 记录成功调用
func (s *CircuitBreakerService) RecordSuccess(serviceID string) {
	breaker := s.getOrCreateBreaker(serviceID)
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	now := time.Now().UnixMilli()

	switch breaker.State {
	case CircuitClosed:
		// 正常状态，重置失败计数
		breaker.FailureCount = 0

	case CircuitHalfOpen:
		// 半开状态，增加成功计数
		breaker.SuccessCount++
		if breaker.SuccessCount >= breaker.SuccessThreshold {
			// 足够成功次数，恢复到正常状态
			breaker.State = CircuitClosed
			breaker.LastStateChange = now
			breaker.FailureCount = 0
			breaker.SuccessCount = 0
		}
	}
}

// RecordFailure 记录失败调用
func (s *CircuitBreakerService) RecordFailure(serviceID string) {
	breaker := s.getOrCreateBreaker(serviceID)
	breaker.mu.Lock()
	defer breaker.mu.Unlock()

	now := time.Now().UnixMilli()

	breaker.FailureCount++
	breaker.LastFailureTime = now

	switch breaker.State {
	case CircuitClosed:
		// 正常状态，检查是否需要熔断
		if breaker.FailureCount >= breaker.FailureThreshold {
			// 达到失败阈值，进入熔断状态
			breaker.State = CircuitOpen
			breaker.LastStateChange = now
		}

	case CircuitHalfOpen:
		// 半开状态，失败则重新进入熔断
		breaker.State = CircuitOpen
		breaker.LastStateChange = now
		breaker.SuccessCount = 0
	}
}

// GetState 获取熔断器状态
func (s *CircuitBreakerService) GetState(serviceID string) CircuitState {
	breaker := s.getOrCreateBreaker(serviceID)
	breaker.mu.RLock()
	defer breaker.mu.RUnlock()

	return breaker.State
}

// GetBreakerInfo 获取熔断器信息
func (s *CircuitBreakerService) GetBreakerInfo(serviceID string) *CircuitBreaker {
	breaker := s.getOrCreateBreaker(serviceID)
	breaker.mu.RLock()
	defer breaker.mu.RUnlock()

	// 返回副本
	return &CircuitBreaker{
		ServiceID:        breaker.ServiceID,
		State:            breaker.State,
		FailureCount:     breaker.FailureCount,
		SuccessCount:     breaker.SuccessCount,
		LastFailureTime:  breaker.LastFailureTime,
		LastStateChange:  breaker.LastStateChange,
		FailureThreshold: breaker.FailureThreshold,
		SuccessThreshold: breaker.SuccessThreshold,
		Timeout:          breaker.Timeout,
		HalfOpenCalls:    breaker.HalfOpenCalls,
	}
}

// Reset 重置熔断器
func (s *CircuitBreakerService) Reset(serviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if breaker, exists := s.breakers[serviceID]; exists {
		breaker.mu.Lock()
		breaker.State = CircuitClosed
		breaker.FailureCount = 0
		breaker.SuccessCount = 0
		breaker.LastStateChange = time.Now().UnixMilli()
		breaker.mu.Unlock()
	}
}
