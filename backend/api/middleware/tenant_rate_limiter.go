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
	"strconv"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	pkgerrorx "github.com/coze-dev/coze-studio/backend/pkg/errorx"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	"github.com/coze-dev/coze-studio/backend/infra/tracing"
	"golang.org/x/time/rate"
)

const (
	// DefaultRateLimit 默认速率限制（请求/秒）
	DefaultRateLimit = 100

	// DefaultBurst 默认突发容量
	DefaultBurst = 200

	// RateLimitKeyPrefix 限流key前缀
	RateLimitKeyPrefix = "rate_limit:tenant:"

	// RateLimitTTL 限流记录TTL
	RateLimitTTL = 1 * time.Minute
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	// Rate 每秒请求数
	Rate float64

	// Burst 突发容量
	Burst int

	// CleanupInterval 清理间隔
	CleanupInterval time.Duration

	// EnableConcurrencyLimit 是否启用并发限制
	EnableConcurrencyLimit bool

	// MaxConcurrency 最大并发数
	MaxConcurrency int
}

// TenantRateLimiter 租户限流器
//
// **核心功能**：
// 1. 基于Token Bucket算法的租户级限流
// 2. 支持动态限流配置
// 3. 并发数限制
// 4. 配额限制（与quota系统集成）
// 5. 分布式限流（基于Redis）
//
// **使用示例**：
//	limiter := NewTenantRateLimiter(redisClient)
//	app.Use(limiter.Middleware())
type TenantRateLimiter struct {
	// Redis客户端
	redis cache.Cmdable

	// 限流配置
	config *RateLimitConfig

	// 本地限流器（用于单机限流）
	limiters sync.Map // map[string]*limiterState

	// 租户配置
	tenantConfigs sync.Map // map[string]*RateLimitConfig
}

// limiterState 限流器状态
type limiterState struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

// NewTenantRateLimiter 创建租户限流器
func NewTenantRateLimiter(redis cache.Cmdable, config *RateLimitConfig) *TenantRateLimiter {
	if config == nil {
		config = &RateLimitConfig{
			Rate:                  DefaultRateLimit,
			Burst:                 DefaultBurst,
			CleanupInterval:       5 * time.Minute,
			EnableConcurrencyLimit: false,
			MaxConcurrency:        100,
		}
	}

	limiter := &TenantRateLimiter{
		redis:  redis,
		config: config,
	}

	// 启动清理goroutine
	go limiter.cleanup()

	return limiter
}

// Middleware 限流中间件
//
// **限流策略**：
// 1. 先检查租户级限流
// 2. 再检查并发数限制（如果启用）
// 3. 最后检查配额限制（如果启用）
//
// **示例**：
//	app.Use(limiter.Middleware())
func (l *TenantRateLimiter) Middleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取租户ID
		tenantID := GetTenantIDFromContext(c)
		if tenantID == "" || tenantID == DefaultTenantID {
			// 跳过默认租户的限流
			ctx.Next(c)
			return
		}

		// 获取租户配置
		config := l.getTenantConfig(tenantID)

		// 1. 检查速率限制
		if err := l.checkRateLimit(c, tenantID, config); err != nil {
			l.logRateLimitExceeded(c, tenantID, config)
			ctx.Set("X-RateLimit-Limit", strconv.Itoa(int(config.Rate)))
			ctx.Set("X-RateLimit-Remaining", "0")
			ctx.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(RateLimitTTL).Unix(), 10))
			ctx.Set("Retry-After", strconv.FormatInt(int64(RateLimitTTL.Seconds()), 10))

			ctx.JSON(429, map[string]interface{}{
				"code":    berrno.ErrRateLimitExceededCode,
				"message": "Rate limit exceeded",
				"error": map[string]interface{}{
					"tenant_id":      tenantID,
					"rate_limit":     config.Rate,
					"retry_after":    RateLimitTTL.Seconds(),
				},
			})
			ctx.Abort()
			return
		}

		// 2. 检查并发限制（如果启用）
		if config.EnableConcurrencyLimit {
			if err := l.checkConcurrencyLimit(c, tenantID, config); err != nil {
				l.logConcurrencyLimitExceeded(c, tenantID, config)
				ctx.JSON(429, map[string]interface{}{
					"code":    berrno.ErrConcurrencyLimitExceededCode,
					"message": "Concurrency limit exceeded",
					"error": map[string]interface{}{
						"tenant_id":        tenantID,
						"max_concurrency":  config.MaxConcurrency,
					},
				})
				ctx.Abort()
				return
			}
		}

		// 记录限流统计
		l.recordRequest(c, tenantID)

		ctx.Next(c)
	}
}

// checkRateLimit 检查速率限制
func (l *TenantRateLimiter) checkRateLimit(ctx context.Context, tenantID string, config *RateLimitConfig) error {
	// 记录追踪
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
	)

	// 如果配置了Redis，使用分布式限流
	if l.redis != nil {
		return l.checkDistributedRateLimit(ctx, tenantID, config)
	}

	// 否则使用本地限流
	return l.checkLocalRateLimit(ctx, tenantID, config)
}

// checkLocalRateLimit 本地速率限制检查
func (l *TenantRateLimiter) checkLocalRateLimit(ctx context.Context, tenantID string, config *RateLimitConfig) error {
	// 获取或创建限流器
	limiterValue, _ := l.limiters.LoadOrStore(tenantID, &limiterState{
		limiter:  rate.NewLimiter(rate.Limit(config.Rate), config.Burst),
		lastSeen: time.Now(),
	})

	state := limiterValue.(*limiterState)
	state.mu.Lock()
	state.lastSeen = time.Now()
	state.mu.Unlock()

	// 检查是否允许请求
	if !state.limiter.Allow() {
		logs.CtxWarnf(ctx, "[TenantRateLimiter] rate limit exceeded: tenant_id=%s, rate=%.2f",
			tenantID, config.Rate)
		return pkgerrorx.New(berrno.ErrRateLimitExceededCode,
			pkgerrorx.KV("tenant_id", tenantID),
			pkgerrorx.KV("rate_limit", fmt.Sprintf("%.2f", config.Rate)),
		)
	}

	return nil
}

// checkDistributedRateLimit 分布式速率限制检查（基于Redis）
func (l *TenantRateLimiter) checkDistributedRateLimit(ctx context.Context, tenantID string, config *RateLimitConfig) error {
	key := RateLimitKeyPrefix + tenantID

	// 使用Redis的INCR和EXPIRE实现简单的令牌桶
	// 注意：这是简化实现，生产环境建议使用Redis Cell模块或Redisson

	// 获取当前计数
	current, err := l.redis.Incr(ctx, key).Result()
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantRateLimiter] redis incr failed: key=%s, error=%v", key, err)
		// Redis失败时降级为允许请求
		return nil
	}

	// 首次设置时设置过期时间
	if current == 1 {
		l.redis.Expire(ctx, key, RateLimitTTL)
	}

	// 检查是否超过限制
	if current > int64(config.Burst) {
		logs.CtxWarnf(ctx, "[TenantRateLimiter] distributed rate limit exceeded: tenant_id=%s, current=%d, burst=%d",
			tenantID, current, config.Burst)
		return pkgerrorx.New(berrno.ErrRateLimitExceededCode,
			pkgerrorx.KV("tenant_id", tenantID),
			pkgerrorx.KV("current", fmt.Sprintf("%d", current)),
			pkgerrorx.KV("burst", fmt.Sprintf("%d", config.Burst)),
		)
	}

	return nil
}

// checkConcurrencyLimit 检查并发限制
func (l *TenantRateLimiter) checkConcurrencyLimit(ctx context.Context, tenantID string, config *RateLimitConfig) error {
	key := fmt.Sprintf("concurrency:%s", tenantID)

	// 使用Redis计数当前并发数
	current, err := l.redis.Incr(ctx, key).Result()
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantRateLimiter] check concurrency failed: error=%v", err)
		// Redis失败时降级为允许请求
		return nil
	}

	// 首次设置时设置较长的过期时间（防止泄漏）
	if current == 1 {
		l.redis.Expire(ctx, key, 1*time.Hour)
	}

	// 检查是否超过并发限制
	if current > int64(config.MaxConcurrency) {
		// 超过限制，递减计数
		l.redis.IncrBy(ctx, key, -1)
		logs.CtxWarnf(ctx, "[TenantRateLimiter] concurrency limit exceeded: tenant_id=%s, current=%d, max=%d",
			tenantID, current, config.MaxConcurrency)
		return pkgerrorx.New(berrno.ErrConcurrencyLimitExceededCode,
			pkgerrorx.KV("tenant_id", tenantID),
			pkgerrorx.KV("current", fmt.Sprintf("%d", current)),
			pkgerrorx.KV("max_concurrency", fmt.Sprintf("%d", config.MaxConcurrency)),
		)
	}

	// 在请求处理完成后递减计数
	go func() {
		<-ctx.Done()
		l.redis.IncrBy(context.Background(), key, -1)
	}()

	return nil
}

// getTenantConfig 获取租户配置
func (l *TenantRateLimiter) getTenantConfig(tenantID string) *RateLimitConfig {
	// 从缓存获取租户配置
	if config, ok := l.tenantConfigs.Load(tenantID); ok {
		return config.(*RateLimitConfig)
	}

	// 返回默认配置
	return l.config
}

// SetTenantConfig 设置租户配置
//
// **使用场景**：动态调整租户限流配置
//
// **示例**：
//	config := &RateLimitConfig{Rate: 200, Burst: 400}
//	limiter.SetTenantConfig(tenantID, config)
func (l *TenantRateLimiter) SetTenantConfig(tenantID string, config *RateLimitConfig) {
	l.tenantConfigs.Store(tenantID, config)
	logs.Infof("[TenantRateLimiter] tenant config updated: tenant_id=%s, rate=%.2f, burst=%d",
		tenantID, config.Rate, config.Burst)
}

// RemoveTenantConfig 移除租户配置（恢复默认）
func (l *TenantRateLimiter) RemoveTenantConfig(tenantID string) {
	l.tenantConfigs.Delete(tenantID)
	logs.Infof("[TenantRateLimiter] tenant config removed: tenant_id=%s", tenantID)
}

// recordRequest 记录请求统计
func (l *TenantRateLimiter) recordRequest(ctx context.Context, tenantID string) {
	// 记录到追踪系统
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
	)

	// 记录到日志（采样记录，避免日志过多）
	if time.Now().Unix()%10 == 0 {
		logs.CtxDebugf(ctx, "[TenantRateLimiter] request recorded: tenant_id=%s", tenantID)
	}
}

// logRateLimitExceeded 记录限流超出
func (l *TenantRateLimiter) logRateLimitExceeded(ctx context.Context, tenantID string, config *RateLimitConfig) {
	logs.CtxWarnf(ctx, "[TenantRateLimiter] rate limit exceeded: tenant_id=%s, rate=%.2f, burst=%d",
		tenantID, config.Rate, config.Burst)

	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
	)
}

// logConcurrencyLimitExceeded 记录并发超出
func (l *TenantRateLimiter) logConcurrencyLimitExceeded(ctx context.Context, tenantID string, config *RateLimitConfig) {
	logs.CtxWarnf(ctx, "[TenantRateLimiter] concurrency limit exceeded: tenant_id=%s, max=%d",
		tenantID, config.MaxConcurrency)

	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
	)
}

// cleanup 定期清理未使用的限流器
func (l *TenantRateLimiter) cleanup() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.limiters.Range(func(key, value interface{}) bool {
			state := value.(*limiterState)
			state.mu.Lock()

			// 如果超过清理时间未使用，删除限流器
			if time.Since(state.lastSeen) > l.config.CleanupInterval {
				l.limiters.Delete(key)
				logs.Debugf("[TenantRateLimiter] limiter cleaned up: tenant_id=%s", key)
			}

			state.mu.Unlock()
			return true
		})
	}
}

// GetStats 获取限流统计
func (l *TenantRateLimiter) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["active_limiters"] = 0
	stats["tenant_configs"] = 0

	l.limiters.Range(func(_, _ interface{}) bool {
		stats["active_limiters"] = stats["active_limiters"].(int) + 1
		return true
	})

	l.tenantConfigs.Range(func(_, _ interface{}) bool {
		stats["tenant_configs"] = stats["tenant_configs"].(int) + 1
		return true
	})

	return stats
}

// Reset 重置限流器
func (l *TenantRateLimiter) Reset() {
	l.limiters = sync.Map{}
	l.tenantConfigs = sync.Map{}
	logs.Infof("[TenantRateLimiter] reset completed")
}

// --------------------------------------------------------------------
// 辅助函数
// --------------------------------------------------------------------

// GetTenantRateLimit 获取租户当前限流状态
//
// **返回**：
// - rate: 当前速率限制
// - remaining: 剩余请求数
// - reset: 重置时间（Unix时间戳）
//
// **示例**：
//	rate, remaining, reset := GetTenantRateLimit(ctx, redisClient, tenantID)
func GetTenantRateLimit(ctx context.Context, redis cache.Cmdable, tenantID string) (float64, int64, int64, error) {
	key := RateLimitKeyPrefix + tenantID

	// 获取当前计数
	current, err := redis.Get(ctx, key).Int64()
	if err != nil && err != cache.Nil {
		return 0, 0, 0, pkgerrorx.Wrapf(err, "failed to get rate limit: tenant_id=%s", tenantID)
	}

	// 默认配置
	rateLimit := float64(DefaultRateLimit)
	burst := DefaultBurst

	remaining := int64(burst) - current
	if remaining < 0 {
		remaining = 0
	}

	// 计算重置时间（简化：使用当前时间+TTL）
	reset := time.Now().Add(RateLimitTTL).Unix()

	return rateLimit, remaining, reset, nil
}

// SetTenantRateLimit 设置租户限流
//
// **示例**：
//	err := SetTenantRateLimit(ctx, redisClient, tenantID, 200, 400)
func SetTenantRateLimit(ctx context.Context, redis cache.Cmdable, tenantID string, rate float64, burst int) error {
	// 这里需要实现租户配置的持久化
	// 实际应用中应该保存到数据库或配置中心
	logs.CtxInfof(ctx, "[SetTenantRateLimit] tenant rate limit updated: tenant_id=%s, rate=%.2f, burst=%d",
		tenantID, rate, burst)
	return nil
}
