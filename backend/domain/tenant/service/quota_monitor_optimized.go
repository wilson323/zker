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
	"sync"
	"time"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// QuotaMonitorOptimized 优化后的配额监控服务（并发处理）
type QuotaMonitorOptimized struct {
	quotaRepo       repository.QuotaRepository
	alertService    AlertService
	maxConcurrency  int // 最大并发数
}

// NewQuotaMonitorOptimized 创建优化版配额监控实例
func NewQuotaMonitorOptimized(
	quotaRepo repository.QuotaRepository,
	alertService AlertService,
	maxConcurrency int,
) *QuotaMonitorOptimized {
	if maxConcurrency <= 0 {
		maxConcurrency = 10 // 默认最多10个并发
	}

	return &QuotaMonitorOptimized{
		quotaRepo:      quotaRepo,
		alertService:   alertService,
		maxConcurrency: maxConcurrency,
	}
}

// MonitorQuotas 并发监控所有租户的配额
func (m *QuotaMonitorOptimized) MonitorQuotas(ctx context.Context) error {
	// 1. 获取所有配额
	quotas, err := m.quotaRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get quotas: %w", err)
	}

	// 2. 并发检查每个配额
	return m.processQuotasConcurrently(ctx, quotas, m.checkQuota)
}

// ResetPeriodicQuotas 并发重置到期配额
func (m *QuotaMonitorOptimized) ResetPeriodicQuotas(ctx context.Context) error {
	quotas, err := m.quotaRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get quotas: %w", err)
	}

	now := time.Now()

	// 筛选需要重置的配额
	quotasToReset := make([]*entity.Quota, 0)
	for _, quota := range quotas {
		if m.shouldResetQuota(quota, now) {
			quotasToReset = append(quotasToReset, quota)
		}
	}

	// 并发重置
	return m.processQuotasConcurrently(ctx, quotasToReset, func(ctx context.Context, quota *entity.Quota) {
		_ = m.quotaRepo.ResetUsage(ctx, quota.QuotaID)
	})
}

// processQuotasConcurrently 并发处理配额
func (m *QuotaMonitorOptimized) processQuotasConcurrently(
	ctx context.Context,
	quotas []*entity.Quota,
	processFunc func(context.Context, *entity.Quota),
) error {
	if len(quotas) == 0 {
		return nil
	}

	// 使用信号量控制并发数
	sem := make(chan struct{}, m.maxConcurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for _, quota := range quotas {
		wg.Add(1)
		sem <- struct{}{} // 获取信号量

		go func(q *entity.Quota) {
			defer wg.Done()
			defer func() { <-sem }() // 释放信号量

			// 处理配额
			processFunc(ctx, q)

		}(quota)
	}

	wg.Wait()
	return firstErr
}

// shouldResetQuota 判断是否应该重置配额
func (m *QuotaMonitorOptimized) shouldResetQuota(quota *entity.Quota, now time.Time) bool {
	if quota.ResetCycle == entity.ResetCycleNever {
		return false
	}

	// 如果从未重置过，需要重置
	if quota.LastResetAt == 0 {
		return true
	}

	lastReset := time.UnixMilli(quota.LastResetAt)

	switch quota.ResetCycle {
	case entity.ResetCycleDaily:
		// 如果上次重置不是今天
		return lastReset.Day() != now.Day() ||
			lastReset.Month() != now.Month() ||
			lastReset.Year() != now.Year()

	case entity.ResetCycleMonthly:
		// 如果上次重置不是本月
		return lastReset.Month() != now.Month() ||
			lastReset.Year() != now.Year()

	case entity.ResetCycleYearly:
		// 如果上次重置不是本年
		return lastReset.Year() != now.Year()

	default:
		return false
	}
}

// checkQuota 检查单个配额（线程安全）
func (m *QuotaMonitorOptimized) checkQuota(ctx context.Context, quota *entity.Quota) {
	if quota.MaxLimit <= 0 {
		return // 无限制配额，不需要监控
	}

	usagePercent := float64(quota.UsedCount) / float64(quota.MaxLimit) * 100

	// 1. 警告：使用率超过 80%
	if usagePercent >= 80 && usagePercent < 95 {
		_ = m.alertService.SendAlert(ctx, &QuotaAlert{
			TenantID:     quota.TenantID,
			ResourceType: quota.ResourceType,
			AlertType:    "warning",
			Usage:        quota.UsedCount,
			MaxLimit:     quota.MaxLimit,
			UsagePercent: usagePercent,
			AlertTime:    time.Now(),
		})
	}

	// 2. 严重：使用率超过 95%
	if usagePercent >= 95 && usagePercent < 100 {
		_ = m.alertService.SendAlert(ctx, &QuotaAlert{
			TenantID:     quota.TenantID,
			ResourceType: quota.ResourceType,
			AlertType:    "critical",
			Usage:        quota.UsedCount,
			MaxLimit:     quota.MaxLimit,
			UsagePercent: usagePercent,
			AlertTime:    time.Now(),
		})
	}

	// 3. 超额：使用率超过 100%
	if usagePercent >= 100 {
		_ = m.alertService.SendAlert(ctx, &QuotaAlert{
			TenantID:     quota.TenantID,
			ResourceType: quota.ResourceType,
			AlertType:    "exceeded",
			Usage:        quota.UsedCount,
			MaxLimit:     quota.MaxLimit,
			UsagePercent: usagePercent,
			AlertTime:    time.Now(),
		})
	}
}

// MonitorTenantQuotas 监控指定租户的配额
func (m *QuotaMonitorOptimized) MonitorTenantQuotas(ctx context.Context, tenantID string) error {
	quotas, err := m.quotaRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant quotas: %w", err)
	}

	return m.processQuotasConcurrently(ctx, quotas, m.checkQuota)
}

// GetQuotaStatus 获取配额状态
func (m *QuotaMonitorOptimized) GetQuotaStatus(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*QuotaStatus, error) {
	quota, err := m.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return nil, err
	}

	if quota == nil {
		return &QuotaStatus{
			TenantID:      tenantID,
			ResourceType:  resourceType,
			HasQuota:      false,
			IsUnlimited:   true,
			Usage:         0,
			MaxLimit:      0,
			Remaining:     0,
			UsagePercent:  0,
			AlertLevel:    "none",
			ShouldReset:   false,
			LastResetTime: nil,
		}, nil
	}

	status := &QuotaStatus{
		TenantID:      tenantID,
		ResourceType:  resourceType,
		HasQuota:      true,
		IsUnlimited:   quota.MaxLimit <= 0,
		Usage:         quota.UsedCount,
		MaxLimit:      quota.MaxLimit,
		Remaining:     quota.MaxLimit - quota.UsedCount,
		UsagePercent:  float64(quota.UsedCount) / float64(quota.MaxLimit) * 100,
		ShouldReset:   m.shouldResetQuota(quota, time.Now()),
	}

	if quota.LastResetAt > 0 {
		resetTime := time.UnixMilli(quota.LastResetAt)
		status.LastResetTime = &resetTime
	}

	// 确定告警级别
	if quota.MaxLimit <= 0 {
		status.AlertLevel = "none"
	} else if status.UsagePercent >= 100 {
		status.AlertLevel = "exceeded"
	} else if status.UsagePercent >= 95 {
		status.AlertLevel = "critical"
	} else if status.UsagePercent >= 80 {
		status.AlertLevel = "warning"
	} else {
		status.AlertLevel = "normal"
	}

	return status, nil
}

// GetAllQuotaStatuses 获取租户的所有配额状态
func (m *QuotaMonitorOptimized) GetAllQuotaStatuses(ctx context.Context, tenantID string) ([]*QuotaStatus, error) {
	quotas, err := m.quotaRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	statuses := make([]*QuotaStatus, 0, len(quotas))

	for _, quota := range quotas {
		status, err := m.GetQuotaStatus(ctx, tenantID, quota.ResourceType)
		if err != nil {
			continue
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// StartMonitoringDaemon 启动监控守护进程
func (m *QuotaMonitorOptimized) StartMonitoringDaemon(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 并发监控所有配额
			_ = m.MonitorQuotas(ctx)

			// 并发重置到期配额
			_ = m.ResetPeriodicQuotas(ctx)
		}
	}
}

// BatchGetQuotaStatuses 批量获取多个租户的配额状态（并发优化）
func (m *QuotaMonitorOptimized) BatchGetQuotaStatuses(ctx context.Context, tenantIDs []string) (map[string][]*QuotaStatus, error) {
	if len(tenantIDs) == 0 {
		return make(map[string][]*QuotaStatus), nil
	}

	result := make(map[string][]*QuotaStatus)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 限制并发数
	sem := make(chan struct{}, m.maxConcurrency)

	for _, tenantID := range tenantIDs {
		wg.Add(1)
		sem <- struct{}{}

		go func(tid string) {
			defer wg.Done()
			defer func() { <-sem }()

			statuses, err := m.GetAllQuotaStatuses(ctx, tid)
			if err == nil {
				mu.Lock()
				result[tid] = statuses
				mu.Unlock()
			}
		}(tenantID)
	}

	wg.Wait()
	return result, nil
}

// BenchmarkQuotaMonitor 性能基准测试
func BenchmarkQuotaMonitor(b *testing.B) {
	// Setup测试数据
	// ...

	b.Run("Sequential", func(b *testing.B) {
		// 测试串行版本
	})

	b.Run("Concurrent", func(b *testing.B) {
		// 测试并发版本
	})
}
