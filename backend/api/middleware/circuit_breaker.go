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

package middleware

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// CircuitBreakerState 熔断器状态
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota // 关闭状态(正常)
	StateOpen                               // 开启状态(熔断)
	StateHalfOpen                           // 半开状态(尝试恢复)
)

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	// 触发熔穿的错误率阈值
	ErrorThreshold float64 // 错误率阈值(0.0-1.0)

	// 触发熔穿的连续错误数
	ConsecutiveErrors int

	// 触发熔穿的窗口期
	Window time.Duration

	// 熔断恢复的超时时间
	Timeout time.Duration

	// 半开状态的请求数
	HalfOpenMaxRequests int

	// 成功阈值(半开状态下的成功请求数)
	SuccessThreshold int
}

// DefaultCircuitBreakerConfig 默认熔断器配置
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		ErrorThreshold:       0.5,              // 50%错误率
		ConsecutiveErrors:    5,                // 连续5次错误
		Window:               time.Minute,      // 1分钟窗口
		Timeout:              time.Second * 30, // 30秒超时
		HalfOpenMaxRequests:  3,                // 半开状态允许3个请求
		SuccessThreshold:     2,                // 2个成功即恢复
	}
}

// circuitBreaker 熔断器
type circuitBreaker struct {
	config *CircuitBreakerConfig

	state         CircuitBreakerState
	stateMu       sync.RWMutex
	lastStateChange time.Time

	// 统计信息
	requests     int
	errors       int
	consecutiveErrors int
	windowStart  time.Time
	windowMu     sync.RWMutex

	// 半开状态计数
	halfOpenRequests int
	halfOpenSuccesses int
	halfOpenMu       sync.RWMutex
}

// newCircuitBreaker 创建熔断器
func newCircuitBreaker(config *CircuitBreakerConfig) *circuitBreaker {
	return &circuitBreaker{
		config:           config,
		state:            StateClosed,
		lastStateChange:  time.Now(),
		windowStart:      time.Now(),
	}
}

// allow 检查是否允许请求
func (cb *circuitBreaker) allow() bool {
	cb.stateMu.RLock()
	state := cb.state
	cb.stateMu.RUnlock()

	switch state {
	case StateClosed:
		return true
	case StateOpen:
		// 检查是否超时,可以尝试恢复
		if time.Since(cb.lastStateChange) > cb.config.Timeout {
			cb.stateMu.Lock()
			cb.state = StateHalfOpen
			cb.lastStateChange = time.Now()
			cb.halfOpenRequests = 0
			cb.halfOpenSuccesses = 0
			cb.stateMu.Unlock()

			logs.CtxInfof(context.Background(), "[CircuitBreaker] Transition to HALF_OPEN")
			return true
		}
		return false
	case StateHalfOpen:
		cb.halfOpenMu.RLock()
		allowed := cb.halfOpenRequests < cb.config.HalfOpenMaxRequests
		cb.halfOpenMu.RUnlock()
		return allowed
	default:
		return false
	}
}

// recordSuccess 记录成功
func (cb *circuitBreaker) recordSuccess() {
	cb.stateMu.RLock()
	state := cb.state
	cb.stateMu.RUnlock()

	switch state {
	case StateClosed:
		cb.recordRequestInWindow(false)
	case StateHalfOpen:
		cb.halfOpenMu.Lock()
		cb.halfOpenRequests++
		cb.halfOpenSuccesses++
		cb.halfOpenMu.Unlock()

		// 检查是否达到成功阈值
		if cb.halfOpenSuccesses >= cb.config.SuccessThreshold {
			cb.stateMu.Lock()
			cb.state = StateClosed
			cb.lastStateChange = time.Now()
			cb.stateMu.Unlock()

			logs.CtxInfof(context.Background(), "[CircuitBreaker] Transition to CLOSED")
		}
	}
}

// recordFailure 记录失败
func (cb *circuitBreaker) recordFailure() {
	cb.stateMu.RLock()
	state := cb.state
	cb.stateMu.RUnlock()

	switch state {
	case StateClosed:
		cb.recordRequestInWindow(true)

		// 检查是否需要熔断
		if cb.shouldOpen() {
			cb.stateMu.Lock()
			cb.state = StateOpen
			cb.lastStateChange = time.Now()
			cb.stateMu.Unlock()

			logs.CtxWarnf(context.Background(), "[CircuitBreaker] Transition to OPEN")
		}
	case StateHalfOpen:
		cb.halfOpenMu.Lock()
		cb.halfOpenRequests++
		cb.halfOpenMu.Unlock()

		// 半开状态失败,重新熔断
		cb.stateMu.Lock()
		cb.state = StateOpen
		cb.lastStateChange = time.Now()
		cb.stateMu.Unlock()

		logs.CtxWarnf(context.Background(), "[CircuitBreaker] HALF_OPEN failed, transition to OPEN")
	}
}

// recordRequestInWindow 在窗口期内记录请求
func (cb *circuitBreaker) recordRequestInWindow(isError bool) {
	cb.windowMu.Lock()
	defer cb.windowMu.Unlock()

	// 检查是否需要重置窗口
	if time.Since(cb.windowStart) > cb.config.Window {
		cb.requests = 0
		cb.errors = 0
		cb.consecutiveErrors = 0
		cb.windowStart = time.Now()
	}

	cb.requests++
	if isError {
		cb.errors++
		cb.consecutiveErrors++
	} else {
		cb.consecutiveErrors = 0
	}
}

// shouldOpen 判断是否应该开启熔断
func (cb *circuitBreaker) shouldOpen() bool {
	cb.windowMu.RLock()
	defer cb.windowMu.RUnlock()

	// 检查连续错误数
	if cb.consecutiveErrors >= cb.config.ConsecutiveErrors {
		return true
	}

	// 检查错误率
	if cb.requests >= 10 { // 至少10个请求才计算错误率
		errorRate := float64(cb.errors) / float64(cb.requests)
		if errorRate >= cb.config.ErrorThreshold {
			return true
		}
	}

	return false
}

// GetStats 获取统计信息
func (cb *circuitBreaker) GetStats() map[string]interface{} {
	cb.stateMu.RLock()
	state := cb.state
	cb.windowMu.RLock()
	cb.halfOpenMu.RLock()
	defer cb.stateMu.RUnlock()
	defer cb.windowMu.RUnlock()
	defer cb.halfOpenMu.RUnlock()

	stateName := "CLOSED"
	switch state {
	case StateOpen:
		stateName = "OPEN"
	case StateHalfOpen:
		stateName = "HALF_OPEN"
	}

	return map[string]interface{}{
		"state":               stateName,
		"requests":            cb.requests,
		"errors":              cb.errors,
		"consecutive_errors":  cb.consecutiveErrors,
		"half_open_requests":  cb.halfOpenRequests,
		"half_open_successes": cb.halfOpenSuccesses,
		"error_rate":          float64(cb.errors) / float64(cb.requests),
		"last_state_change":   cb.lastStateChange,
	}
}

// 全局熔断器实例
var globalCircuitBreakers map[string]*circuitBreaker
var cbMutex sync.RWMutex

// InitCircuitBreaker 初始化熔断器
func InitCircuitBreaker() {
	cbMutex.Lock()
	defer cbMutex.Unlock()

	if globalCircuitBreakers == nil {
		globalCircuitBreakers = make(map[string]*circuitBreaker)
		logs.CtxInfof(context.Background(), "[CircuitBreaker] Initialized")
	}
}

// getCircuitBreaker 获取熔断器
func getCircuitBreaker(name string, config *CircuitBreakerConfig) *circuitBreaker {
	cbMutex.RLock()
	cb, exists := globalCircuitBreakers[name]
	cbMutex.RUnlock()

	if exists {
		return cb
	}

	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}

	cbMutex.Lock()
	defer cbMutex.Unlock()

	// 双重检查
	if cb, exists := globalCircuitBreakers[name]; exists {
		return cb
	}

	cb = newCircuitBreaker(config)
	globalCircuitBreakers[name] = cb

	logs.CtxInfof(context.Background(), "[CircuitBreaker] Created: %s", name)

	return cb
}

// CircuitBreaker 熔断中间件
// 用途:保护下游服务,防止级联故障
func CircuitBreaker(name string, config *CircuitBreakerConfig) app.HandlerFunc {
	InitCircuitBreaker()
	cb := getCircuitBreaker(name, config)

	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否允许请求
		if !cb.allow() {
			httputil.ServiceUnavailableResp(c, "Circuit breaker is OPEN")
			c.Abort()
			return
		}

		// 执行请求
		c.Next(ctx)

		// 根据响应结果记录成功或失败
		if c.Response.StatusCode() >= 500 {
			cb.recordFailure()
		} else {
			cb.recordSuccess()
		}
	}
}

// CircuitBreakerByPath 基于路径的熔断中间件
// 用途:为不同的API路径配置独立的熔断器
func CircuitBreakerByPath(config *CircuitBreakerConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Path())
		cb := getCircuitBreaker("path:"+path, config)

		if !cb.allow() {
			httputil.ServiceUnavailableResp(c, "Circuit breaker is OPEN")
			c.Abort()
			return
		}

		c.Next(ctx)

		if c.Response.StatusCode() >= 500 {
			cb.recordFailure()
		} else {
			cb.recordSuccess()
		}
	}
}

// GetCircuitBreakerStats 获取熔断器统计信息
func GetCircuitBreakerStats(name string) (map[string]interface{}, error) {
	cbMutex.RLock()
	defer cbMutex.RUnlock()

	cb, exists := globalCircuitBreakers[name]
	if !exists {
		return nil, errors.New("circuit breaker not found")
	}

	return cb.GetStats(), nil
}

// GetAllCircuitBreakersStats 获取所有熔断器统计信息
func GetAllCircuitBreakersStats() map[string]map[string]interface{} {
	cbMutex.RLock()
	defer cbMutex.RUnlock()

	result := make(map[string]map[string]interface{})
	for name, cb := range globalCircuitBreakers {
		result[name] = cb.GetStats()
	}

	return result
}

// ResetCircuitBreaker 重置熔断器(管理员功能)
func ResetCircuitBreaker(name string) error {
	cbMutex.Lock()
	defer cbMutex.Unlock()

	cb, exists := globalCircuitBreakers[name]
	if !exists {
		return errors.New("circuit breaker not found")
	}

	cb.stateMu.Lock()
	cb.state = StateClosed
	cb.lastStateChange = time.Now()
	cb.stateMu.Unlock()

	cb.windowMu.Lock()
	cb.requests = 0
	cb.errors = 0
	cb.consecutiveErrors = 0
	cb.windowStart = time.Now()
	cb.windowMu.Unlock()

	logs.CtxInfof(context.Background(), "[CircuitBreaker] Reset: %s", name)

	return nil
}
