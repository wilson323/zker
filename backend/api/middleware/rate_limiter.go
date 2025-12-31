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
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/pkg/conv"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	// 全局限制
	GlobalRPS int // 每秒请求数

	// 租户级限制
	TenantRPS int // 每租户每秒请求数

	// 用户级限制
	UserRPS int // 每用户每秒请求数

	// IP级限制
	IPRPS int // 每IP每秒请求数

	// Burst 容量
	Burst int // 突发流量容量

	// 清理间隔
	CleanupInterval time.Duration
}

// DefaultRateLimiterConfig 默认限流配置
func DefaultRateLimiterConfig() *RateLimiterConfig {
	return &RateLimiterConfig{
		GlobalRPS:        10000, // 全局10000 RPS
		TenantRPS:        100,   // 每租户100 RPS
		UserRPS:          10,    // 每用户10 RPS
		IPRPS:            100,   // 每IP 100 RPS
		Burst:            5,     // 突发容量5
		CleanupInterval:  5 * time.Minute,
	}
}

// tokenBucket 令牌桶
type tokenBucket struct {
	tokens     int
	maxTokens  int
	refillRate int // 每秒填充令牌数
	lastRefill time.Time
	mu         sync.Mutex
}

// newTokenBucket 创建令牌桶
func newTokenBucket(maxTokens, refillRate int) *tokenBucket {
	return &tokenBucket{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// allow 检查是否允许请求
func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// 计算应该补充的令牌数
	tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
	tb.tokens += tokensToAdd
	if tb.tokens > tb.maxTokens {
		tb.tokens = tb.maxTokens
	}
	tb.lastRefill = now

	// 检查是否有足够的令牌
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

// rateLimiter 限流器
type rateLimiter struct {
	config *RateLimiterConfig

	// 全局限流
	globalLimiter *tokenBucket

	// 租户级限流
	tenantLimiters map[string]*tokenBucket
	tenantMu       sync.RWMutex

	// 用户级限流
	userLimiters map[string]*tokenBucket
	userMu       sync.RWMutex

	// IP级限流
	ipLimiters map[string]*tokenBucket
	ipMu       sync.RWMutex

	// 清理停止信号
	stopCleanup chan struct{}
}

// newRateLimiter 创建限流器
func newRateLimiter(config *RateLimiterConfig) *rateLimiter {
	rl := &rateLimiter{
		config:         config,
		globalLimiter:  newTokenBucket(config.GlobalRPS, config.GlobalRPS),
		tenantLimiters: make(map[string]*tokenBucket),
		userLimiters:   make(map[string]*tokenBucket),
		ipLimiters:     make(map[string]*tokenBucket),
		stopCleanup:    make(chan struct{}),
	}

	// 启动清理协程
	go rl.cleanup()

	return rl
}

// getTenantLimiter 获取租户限流器
func (rl *rateLimiter) getTenantLimiter(tenantID string) *tokenBucket {
	rl.tenantMu.Lock()
	defer rl.tenantMu.Unlock()

	if limiter, exists := rl.tenantLimiters[tenantID]; exists {
		return limiter
	}

	limiter := newTokenBucket(rl.config.Burst, rl.config.TenantRPS)
	rl.tenantLimiters[tenantID] = limiter
	return limiter
}

// getUserLimiter 获取用户限流器
func (rl *rateLimiter) getUserLimiter(userID string) *tokenBucket {
	rl.userMu.Lock()
	defer rl.userMu.Unlock()

	if limiter, exists := rl.userLimiters[userID]; exists {
		return limiter
	}

	limiter := newTokenBucket(rl.config.Burst, rl.config.UserRPS)
	rl.userLimiters[userID] = limiter
	return limiter
}

// getIPLimiter 获取IP限流器
func (rl *rateLimiter) getIPLimiter(ip string) *tokenBucket {
	rl.ipMu.Lock()
	defer rl.ipMu.Unlock()

	if limiter, exists := rl.ipLimiters[ip]; exists {
		return limiter
	}

	limiter := newTokenBucket(rl.config.Burst, rl.config.IPRPS)
	rl.ipLimiters[ip] = limiter
	return limiter
}

// allow 检查是否允许请求
func (rl *rateLimiter) allow(tenantID, userID, ip string) bool {
	// 全局限流
	if !rl.globalLimiter.allow() {
		logs.CtxWarnf(context.Background(), "[RateLimiter] Global rate limit exceeded")
		return false
	}

	// 租户级限流
	if tenantID != "" {
		if !rl.getTenantLimiter(tenantID).allow() {
			logs.CtxWarnf(context.Background(), "[RateLimiter] Tenant rate limit exceeded: tenant=%s", tenantID)
			return false
		}
	}

	// 用户级限流
	if userID != "" {
		if !rl.getUserLimiter(userID).allow() {
			logs.CtxWarnf(context.Background(), "[RateLimiter] User rate limit exceeded: user=%s", userID)
			return false
		}
	}

	// IP级限流
	if ip != "" {
		if !rl.getIPLimiter(ip).allow() {
			logs.CtxWarnf(context.Background(), "[RateLimiter] IP rate limit exceeded: ip=%s", ip)
			return false
		}
	}

	return true
}

// cleanup 定期清理未使用的限流器
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanUnusedLimiters()
		case <-rl.stopCleanup:
			return
		}
	}
}

// cleanUnusedLimiters 清理未使用的限流器
func (rl *rateLimiter) cleanUnusedLimiters() {
	// 简化实现:每5分钟清理一次
	// 实际应该记录最后使用时间,只清理超过阈值未使用的

	rl.tenantMu.Lock()
	if len(rl.tenantLimiters) > 10000 {
		// 清理一半
		count := 0
		for k := range rl.tenantLimiters {
			if count > 5000 {
				break
			}
			delete(rl.tenantLimiters, k)
			count++
		}
	}
	rl.tenantMu.Unlock()

	rl.userMu.Lock()
	if len(rl.userLimiters) > 100000 {
		count := 0
		for k := range rl.userLimiters {
			if count > 50000 {
				break
			}
			delete(rl.userLimiters, k)
			count++
		}
	}
	rl.userMu.Unlock()

	rl.ipMu.Lock()
	if len(rl.ipLimiters) > 100000 {
		count := 0
		for k := range rl.ipLimiters {
			if count > 50000 {
				break
			}
			delete(rl.ipLimiters, k)
			count++
		}
	}
	rl.ipMu.Unlock()
}

// stop 停止限流器
func (rl *rateLimiter) stop() {
	close(rl.stopCleanup)
}

// 全局限流器实例
var globalRateLimiter *rateLimiter
var rateLimiterOnce sync.Once

// InitRateLimiter 初始化限流器
func InitRateLimiter(config *RateLimiterConfig) {
	if config == nil {
		config = DefaultRateLimiterConfig()
	}

	rateLimiterOnce.Do(func() {
		globalRateLimiter = newRateLimiter(config)
		logs.CtxInfof(context.Background(), "[RateLimiter] Initialized with config: GlobalRPS=%d, TenantRPS=%d, UserRPS=%d",
			config.GlobalRPS, config.TenantRPS, config.UserRPS)
	})
}

// RateLimiter 限流中间件
// 遵循企业级规范:支持全局/租户/用户/IP四级限流
func RateLimiter() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if globalRateLimiter == nil {
			c.Next(ctx)
			return
		}

		// 获取租户ID、用户ID、IP
		tenantID := conv.BytesToString(c.GetHeader("X-Tenant-ID"))
		userID := conv.BytesToString(c.GetHeader("X-User-ID"))
		ip := c.ClientIP()

		// 检查限流
		if !globalRateLimiter.allow(tenantID, userID, ip) {
			httputil.RateLimitResp(c, 60) // 建议60秒后重试
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// GetRateLimitStats 获取限流统计信息
func GetRateLimitStats() map[string]interface{} {
	if globalRateLimiter == nil {
		return nil
	}

	return map[string]interface{}{
		"global_tokens":       globalRateLimiter.globalLimiter.tokens,
		"active_tenants":      len(globalRateLimiter.tenantLimiters),
		"active_users":        len(globalRateLimiter.userLimiters),
		"active_ips":          len(globalRateLimiter.ipLimiters),
		"global_rps":          globalRateLimiter.config.GlobalRPS,
		"tenant_rps":          globalRateLimiter.config.TenantRPS,
		"user_rps":            globalRateLimiter.config.UserRPS,
		"ip_rps":              globalRateLimiter.config.IPRPS,
	}
}

// ResetUserRateLimiter 重置用户限流器(管理员功能)
func ResetUserRateLimiter(userID string) error {
	if globalRateLimiter == nil {
		return fmt.Errorf("rate limiter not initialized")
	}

	globalRateLimiter.userMu.Lock()
	defer globalRateLimiter.userMu.Unlock()

	delete(globalRateLimiter.userLimiters, userID)
	logs.CtxInfof(context.Background(), "[RateLimiter] Reset rate limiter for user: %s", userID)

	return nil
}

// ResetTenantRateLimiter 重置租户限流器(管理员功能)
func ResetTenantRateLimiter(tenantID string) error {
	if globalRateLimiter == nil {
		return fmt.Errorf("rate limiter not initialized")
	}

	globalRateLimiter.tenantMu.Lock()
	defer globalRateLimiter.tenantMu.Unlock()

	delete(globalRateLimiter.tenantLimiters, tenantID)
	logs.CtxInfof(context.Background(), "[RateLimiter] Reset rate limiter for tenant: %s", tenantID)

	return nil
}
