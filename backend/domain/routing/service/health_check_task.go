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
	"net/http"
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// HealthCheckBotRepository Bot健康检查仓储接口
type HealthCheckBotRepository interface {
	// GetAllActiveBots 获取所有活跃的Bot（用于健康检查）
	GetAllActiveBots(ctx context.Context) ([]*HealthCheckBotInfo, error)
}

// HealthCheckBotInfo Bot基本信息（用于健康检查）
type HealthCheckBotInfo struct {
	BotID     string
	TenantID  string
	Name      string
	Status    string
	HealthURL string // 健康检查URL
}

// BotHealthCheckTask Bot健康检查定时任务
type BotHealthCheckTask struct {
	healthMonitor    *ServiceHealthMonitor
	botRepo          HealthCheckBotRepository // Bot健康检查仓储
	checkInterval    time.Duration      // 检查间隔
	checkTimeout     time.Duration      // 单次检查超时
	maxConcurrency   int                // 最大并发数
	httpClient       *http.Client       // HTTP客户端
	stopChan         chan struct{}      // 停止信号
	wg               sync.WaitGroup     // 等待组
	enabled          bool               // 是否启用
}

// BotHealthCheckConfig Bot健康检查配置
type BotHealthCheckConfig struct {
	CheckInterval  time.Duration // 检查间隔（默认30秒）
	CheckTimeout   time.Duration // 单次检查超时（默认10秒）
	MaxConcurrency int           // 最大并发数（默认20）
}

// DefaultBotHealthCheckConfig 默认配置
var DefaultBotHealthCheckConfig = BotHealthCheckConfig{
	CheckInterval:  30 * time.Second,
	CheckTimeout:   10 * time.Second,
	MaxConcurrency: 20,
}

// NewBotHealthCheckTask 创建Bot健康检查定时任务
func NewBotHealthCheckTask(
	healthMonitor *ServiceHealthMonitor,
	botRepo HealthCheckBotRepository,
	config BotHealthCheckConfig,
) *BotHealthCheckTask {
	if config.CheckInterval == 0 {
		config.CheckInterval = DefaultBotHealthCheckConfig.CheckInterval
	}
	if config.CheckTimeout == 0 {
		config.CheckTimeout = DefaultBotHealthCheckConfig.CheckTimeout
	}
	if config.MaxConcurrency == 0 {
		config.MaxConcurrency = DefaultBotHealthCheckConfig.MaxConcurrency
	}

	return &BotHealthCheckTask{
		healthMonitor:  healthMonitor,
		botRepo:        botRepo,
		checkInterval:  config.CheckInterval,
		checkTimeout:   config.CheckTimeout,
		maxConcurrency: config.MaxConcurrency,
		httpClient: &http.Client{
			Timeout: config.CheckTimeout,
		},
		stopChan: make(chan struct{}),
		enabled:  false,
	}
}

// Start 启动健康检查定时任务
func (t *BotHealthCheckTask) Start(ctx context.Context) {
	t.wg.Add(1)
	t.enabled = true

	logs.CtxInfof(ctx, "[BotHealthCheckTask] starting health check task with interval %v", t.checkInterval)

	go func() {
		defer t.wg.Done()

		ticker := time.NewTicker(t.checkInterval)
		defer ticker.Stop()

		// 立即执行一次
		t.executeHealthCheck(ctx)

		for {
			select {
			case <-ctx.Done():
				logs.CtxInfof(ctx, "[BotHealthCheckTask] context canceled, stopping health check task")
				return
			case <-t.stopChan:
				logs.CtxInfof(ctx, "[BotHealthCheckTask] received stop signal")
				return
			case <-ticker.C:
				t.executeHealthCheck(ctx)
			}
		}
	}()
}

// Stop 停止健康检查定时任务
func (t *BotHealthCheckTask) Stop() {
	if !t.enabled {
		return
	}

	logs.Infof("[BotHealthCheckTask] stopping health check task")
	close(t.stopChan)
	t.wg.Wait()
	t.enabled = false
	logs.Infof("[BotHealthCheckTask] health check task stopped")
}

// IsEnabled 检查任务是否启用
func (t *BotHealthCheckTask) IsEnabled() bool {
	return t.enabled
}

// executeHealthCheck 执行健康检查
func (t *BotHealthCheckTask) executeHealthCheck(ctx context.Context) {
	startTime := time.Now()
	logs.CtxInfof(ctx, "[BotHealthCheckTask] executing scheduled health check")

	// 从数据库获取所有活跃的Bot
	bots, err := t.botRepo.GetAllActiveBots(ctx)
	if err != nil {
		logs.CtxErrorf(ctx, "[BotHealthCheckTask] failed to get active bots: %v", err)
		return
	}

	if len(bots) == 0 {
		logs.CtxDebugf(ctx, "[BotHealthCheckTask] no active bots to check")
		return
	}

	// 提取Bot ID列表
	botIDs := make([]string, len(bots))
	for i, bot := range bots {
		botIDs[i] = bot.BotID
	}

	// 并发执行健康检查
	successCount, failureCount := t.checkBotsConcurrently(ctx, botIDs)

	duration := time.Since(startTime)
	logs.CtxInfof(ctx, "[BotHealthCheckTask] health check completed: "+
		"total=%d, success=%d, failure=%d, duration=%v",
		len(botIDs), successCount, failureCount, duration)

	// 记录Prometheus指标
	// TODO: 实现metrics.RecordHealthCheckExecution
	_ = duration // 避免未使用变量警告
}

// checkBotsConcurrently 并发检查多个Bot的健康状态
func (t *BotHealthCheckTask) checkBotsConcurrently(
	ctx context.Context,
	botIDs []string,
) (successCount, failureCount int) {
	if len(botIDs) == 0 {
		return 0, 0
	}

	// 使用信号量控制并发数
	sem := make(chan struct{}, t.maxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var successes, failures int

	for _, botID := range botIDs {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(id string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			// 执行健康检查
			if err := t.checkSingleBot(ctx, id); err != nil {
				mu.Lock()
				failures++
				mu.Unlock()

				// 记录失败
				logs.CtxWarnf(ctx, "[BotHealthCheckTask] bot %s health check failed: %v", id, err)
				t.healthMonitor.RecordFailure("bot:"+id, err.Error())
			} else {
				mu.Lock()
				successes++
				mu.Unlock()

				logs.CtxDebugf(ctx, "[BotHealthCheckTask] bot %s health check passed", id)
			}
		}(botID)
	}

	wg.Wait()
	return successes, failures
}

// checkSingleBot 检查单个Bot的健康状态
func (t *BotHealthCheckTask) checkSingleBot(ctx context.Context, botID string) error {
	startTime := time.Now()

	// 1. 设置超时
	_, cancel := context.WithTimeout(ctx, t.checkTimeout)
	defer cancel()

	// 2. 构建健康检查URL
	// TODO: 根据实际的Bot API URL配置
	// healthCheckURL := fmt.Sprintf("http://bot-%s/health", botID)
	//
	// 示例：发送HTTP请求到Bot的健康检查端点
	// req, err := http.NewRequestWithContext(checkCtx, "GET", healthCheckURL, nil)
	// if err != nil {
	//     return fmt.Errorf("create request failed: %w", err)
	// }
	//
	// resp, err := t.httpClient.Do(req)
	// if err != nil {
	//     return fmt.Errorf("http request failed: %w", err)
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	//     return fmt.Errorf("health check returned status %d", resp.StatusCode)
	// }

	// 临时模拟：假设健康检查通过
	// 实际实现时需要删除下面的模拟代码
	latency := time.Since(startTime)

	// 记录成功
	t.healthMonitor.RecordSuccess("bot:"+botID, latency)

	return nil
}

// GetHealthStatus 获取Bot的健康状态
func (t *BotHealthCheckTask) GetHealthStatus(ctx context.Context, botID string) (bool, error) {
	health, err := t.healthMonitor.GetBotHealth(ctx, botID)
	if err != nil {
		return false, err
	}
	return health.IsHealthy, nil
}

// GetHealthStats 获取健康统计信息
func (t *BotHealthCheckTask) GetHealthStats(botID string) (*ServiceHealthStats, bool) {
	return t.healthMonitor.GetHealthStats("bot:" + botID)
}
