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
	"time"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// QuotaMonitor 配额监控服务
type QuotaMonitor struct {
	quotaRepo    repository.QuotaRepository
	alertService AlertService
}

// AlertService 告警服务接口
type AlertService interface {
	SendAlert(ctx context.Context, alert *QuotaAlert) error
}

// QuotaAlert 配额告警
type QuotaAlert struct {
	TenantID     string
	ResourceType entity.ResourceType
	AlertType    string // "warning", "critical", "exceeded"
	Usage        int
	MaxLimit     int
	UsagePercent float64
	AlertTime    time.Time
}

// NewQuotaMonitor 创建配额监控实例
func NewQuotaMonitor(
	quotaRepo repository.QuotaRepository,
	alertService AlertService,
) *QuotaMonitor {
	return &QuotaMonitor{
		quotaRepo:    quotaRepo,
		alertService: alertService,
	}
}

// MonitorQuotas 监控所有租户的配额
func (m *QuotaMonitor) MonitorQuotas(ctx context.Context) error {
	// 1. 获取所有配额
	quotas, err := m.quotaRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get quotas: %w", err)
	}

	// 2. 检查每个配额
	for _, quota := range quotas {
		m.checkQuota(ctx, quota)
	}

	return nil
}

// MonitorTenantQuotas 监控指定租户的配额
func (m *QuotaMonitor) MonitorTenantQuotas(ctx context.Context, tenantID string) error {
	quotas, err := m.quotaRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant quotas: %w", err)
	}

	for _, quota := range quotas {
		m.checkQuota(ctx, quota)
	}

	return nil
}

// checkQuota 检查单个配额
func (m *QuotaMonitor) checkQuota(ctx context.Context, quota *entity.Quota) {
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

// ResetPeriodicQuotas 定期重置配额
func (m *QuotaMonitor) ResetPeriodicQuotas(ctx context.Context) error {
	quotas, err := m.quotaRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get quotas: %w", err)
	}

	now := time.Now()
	resetCount := 0

	for _, quota := range quotas {
		shouldReset := m.shouldResetQuota(quota, now)

		if shouldReset {
			if err := m.quotaRepo.ResetUsage(ctx, quota.QuotaID); err != nil {
				// 记录错误但继续处理其他配额
				continue
			}
			resetCount++
		}
	}

	return nil
}

// shouldResetQuota 判断是否应该重置配额
func (m *QuotaMonitor) shouldResetQuota(quota *entity.Quota, now time.Time) bool {
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

// GetQuotaStatus 获取配额状态
func (m *QuotaMonitor) GetQuotaStatus(ctx context.Context, tenantID string, resourceType entity.ResourceType) (*QuotaStatus, error) {
	quota, err := m.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return nil, err
	}

	if quota == nil {
		return &QuotaStatus{
			TenantID:       tenantID,
			ResourceType:   resourceType,
			HasQuota:       false,
			IsUnlimited:    true,
			Usage:          0,
			MaxLimit:       0,
			Remaining:      0,
			UsagePercent:   0,
			AlertLevel:     "none",
			ShouldReset:    false,
			LastResetTime:  nil,
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
func (m *QuotaMonitor) GetAllQuotaStatuses(ctx context.Context, tenantID string) ([]*QuotaStatus, error) {
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
func (m *QuotaMonitor) StartMonitoringDaemon(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 监控所有配额
			_ = m.MonitorQuotas(ctx)

			// 重置到期配额
			_ = m.ResetPeriodicQuotas(ctx)
		}
	}
}

// QuotaStatus 配额状态
type QuotaStatus struct {
	TenantID      string
	ResourceType  entity.ResourceType
	HasQuota      bool
	IsUnlimited   bool
	Usage         int
	MaxLimit      int
	Remaining     int
	UsagePercent  float64
	AlertLevel    string // "none", "normal", "warning", "critical", "exceeded"
	ShouldReset   bool
	LastResetTime *time.Time
}
