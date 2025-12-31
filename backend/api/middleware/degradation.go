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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// DegradationLevel 降级级别
type DegradationLevel int

const (
	DegradationLevelNone DegradationLevel = iota // 不降级
	DegradationLevelLow                          // 轻度降级
	DegradationLevelMedium                       // 中度降级
	DegradationLevelHigh                         // 重度降级
)

// DegradationConfig 降级配置
type DegradationConfig struct {
	// 当前降级级别
	CurrentLevel DegradationLevel

	// 降级策略
	Strategies map[string]*DegradationStrategy

	// 自动降级阈值
	AutoDegradeThreshold float64 // CPU/内存使用率阈值
	AutoRecoverThreshold  float64 // 自动恢复阈值
}

// DegradationStrategy 降级策略
type DegradationStrategy struct {
	// 策略名称
	Name string

	// 适用的降级级别
	Level DegradationLevel

	// 处理函数
	Handler func(ctx context.Context, c *app.RequestContext) (bool, interface{})

	// 是否启用
	Enabled bool
}

// degradationManager 降级管理器
type degradationManager struct {
	config *DegradationConfig
	mu     sync.RWMutex
}

// newDegradationManager 创建降级管理器
func newDegradationManager(config *DegradationConfig) *degradationManager {
	if config == nil {
		config = &DegradationConfig{
			CurrentLevel:         DegradationLevelNone,
			Strategies:           make(map[string]*DegradationStrategy),
			AutoDegradeThreshold: 0.8, // 80%
			AutoRecoverThreshold: 0.5, // 50%
		}
	}

	return &degradationManager{
		config: config,
	}
}

// setLevel 设置降级级别
func (dm *degradationManager) setLevel(level DegradationLevel) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	oldLevel := dm.config.CurrentLevel
	dm.config.CurrentLevel = level

	logs.CtxInfof(context.Background(), "[Degradation] Level changed: %d -> %d", oldLevel, level)
}

// getLevel 获取当前降级级别
func (dm *degradationManager) getLevel() DegradationLevel {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	return dm.config.CurrentLevel
}

// addStrategy 添加降级策略
func (dm *degradationManager) addStrategy(strategy *DegradationStrategy) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.config.Strategies[strategy.Name] = strategy
	logs.CtxInfof(context.Background(), "[Degradation] Strategy added: %s (level=%d)", strategy.Name, strategy.Level)
}

// shouldDegrade 判断是否应该降级
func (dm *degradationManager) shouldDegrade(strategyName string) bool {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	strategy, exists := dm.config.Strategies[strategyName]
	if !exists || !strategy.Enabled {
		return false
	}

	return dm.config.CurrentLevel >= strategy.Level
}

// execute 执行降级处理
func (dm *degradationManager) execute(ctx context.Context, c *app.RequestContext, strategyName string) (bool, interface{}) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	strategy, exists := dm.config.Strategies[strategyName]
	if !exists {
		return false, nil
	}

	return strategy.Handler(ctx, c)
}

// 全局降级管理器
var globalDegradationManager *degradationManager
var degradationOnce sync.Once

// initDegradationManager 初始化降级管理器
func initDegradationManager() *degradationManager {
	degradationOnce.Do(func() {
		globalDegradationManager = newDegradationManager(nil)
		logs.CtxInfof(context.Background(), "[Degradation] Manager initialized")
	})
	return globalDegradationManager
}

// Degradation 降级中间件
// 用途:在高负载时自动降级非核心功能
func Degradation(strategyName string, handler func(ctx context.Context, c *app.RequestContext) (bool, interface{})) app.HandlerFunc {
	dm := initDegradationManager()

	// 添加策略
	dm.addStrategy(&DegradationStrategy{
		Name:    strategyName,
		Level:   DegradationLevelLow,
		Handler: handler,
		Enabled: true,
	})

	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否需要降级
		if !dm.shouldDegrade(strategyName) {
			c.Next(ctx)
			return
		}

		// 执行降级处理
		handled, data := dm.execute(ctx, c, strategyName)
		if handled {
			httputil.SuccessRespWithData(c, data)
			c.Abort()
			return
		}

		// 继续正常流程
		c.Next(ctx)
	}
}

// SetDegradationLevel 手动设置降级级别
func SetDegradationLevel(level DegradationLevel) {
	dm := initDegradationManager()
	dm.setLevel(level)
}

// GetDegradationLevel 获取当前降级级别
func GetDegradationLevel() DegradationLevel {
	dm := initDegradationManager()
	return dm.getLevel()
}

// ========== 预定义降级策略 ==========

// MockSearchHandler 模拟搜索降级处理器
// 用途:当搜索服务过载时,返回空结果或缓存结果
func MockSearchHandler(ctx context.Context, c *app.RequestContext) (bool, interface{}) {
	// 检查是否有缓存
	if cached := c.GetHeader("X-Cached-Response"); len(cached) > 0 {
		var data interface{}
		if err := json.Unmarshal(cached, &data); err == nil {
			logs.CtxInfof(ctx, "[Degradation] Using cached search result")
			return true, data
		}
	}

	// 返回空结果
	logs.CtxInfof(ctx, "[Degradation] Returning empty search result")
	return true, map[string]interface{}{
		"items":      []interface{}{},
		"total":      0,
		"page":       1,
		"page_size":  20,
		"total_pages": 0,
		"message":    "Search service is busy, please try again later",
	}
}

// MockRecommendHandler 模拟推荐降级处理器
// 用途:当推荐服务过载时,返回默认推荐
func MockRecommendHandler(ctx context.Context, c *app.RequestContext) (bool, interface{}) {
	logs.CtxInfof(ctx, "[Degradation] Returning default recommendations")
	return true, map[string]interface{}{
		"items": []interface{}{
			map[string]string{"id": "default_1", "name": "Popular Bot 1"},
			map[string]string{"id": "default_2", "name": "Popular Bot 2"},
			map[string]string{"id": "default_3", "name": "Popular Bot 3"},
		},
		"total": 3,
		"message": "Recommendation service is busy, showing default recommendations",
	}
}

// MockAnalyticsHandler 模拟分析降级处理器
// 用途:当分析服务过载时,返回简化指标
func MockAnalyticsHandler(ctx context.Context, c *app.RequestContext) (bool, interface{}) {
	logs.CtxInfof(ctx, "[Degradation] Returning simplified analytics")
	return true, map[string]interface{}{
		"message": "Analytics service is busy, showing simplified data",
		"data": map[string]interface{}{
			"total_users":    0,
			"active_users":   0,
			"total_requests": 0,
		},
	}
}

// MockNotificationHandler 模拟通知降级处理器
// 用途:当通知服务过载时,直接返回成功(异步处理)
func MockNotificationHandler(ctx context.Context, c *app.RequestContext) (bool, interface{}) {
	logs.CtxInfof(ctx, "[Degradation] Notification will be sent asynchronously")
	return true, map[string]interface{}{
		"message": "Notification will be sent later",
		"status": "queued",
	}
}

// DegradationWithLevel 指定降级级别的降级中间件
func DegradationWithLevel(strategyName string, level DegradationLevel, handler func(ctx context.Context, c *app.RequestContext) (bool, interface{})) app.HandlerFunc {
	dm := initDegradationManager()

	dm.addStrategy(&DegradationStrategy{
		Name:    strategyName,
		Level:   level,
		Handler: handler,
		Enabled: true,
	})

	return func(ctx context.Context, c *app.RequestContext) {
		if !dm.shouldDegrade(strategyName) {
			c.Next(ctx)
			return
		}

		handled, data := dm.execute(ctx, c, strategyName)
		if handled {
			httputil.SuccessRespWithData(c, data)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// GetDegradationStats 获取降级统计信息
func GetDegradationStats() map[string]interface{} {
	dm := initDegradationManager()
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	levelName := "NONE"
	switch dm.config.CurrentLevel {
	case DegradationLevelLow:
		levelName = "LOW"
	case DegradationLevelMedium:
		levelName = "MEDIUM"
	case DegradationLevelHigh:
		levelName = "HIGH"
	}

	strategies := []string{}
	for name, strategy := range dm.config.Strategies {
		if strategy.Enabled {
			strategies = append(strategies, fmt.Sprintf("%s(level=%d)", name, strategy.Level))
		}
	}

	return map[string]interface{}{
		"current_level":     levelName,
		"enabled_strategies": strategies,
		"total_strategies":  len(dm.config.Strategies),
	}
}
