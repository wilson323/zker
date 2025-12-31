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
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// QuotaEnforcementConfig 配额强制配置
type QuotaEnforcementConfig struct {
	DB             *gorm.DB
	QuotaRepo      repository.QuotaRepository
	StrictMode     bool            // 严格模式：超过配额直接拒绝
	WarningOnly    bool            // 仅警告模式：记录日志但不拒绝
	WhitelistPaths map[string]bool // 跳过检查的路径白名单
}

// QuotaEnforcementMiddleware 配额强制中间件
//
// 功能：
// 1. 在请求处理前检查配额
// 2. 超过配额时拒绝请求（403）
// 3. 记录配额使用情况
//
// 使用示例:
//
//	config := QuotaEnforcementConfig{
//	    DB:         db,
//	    QuotaRepo:  quotaRepo,
//	    StrictMode: true,
//	}
//	h.Use(QuotaEnforcementMiddleware(config))
func QuotaEnforcementMiddleware(config QuotaEnforcementConfig) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 1. 检查是否在白名单中
		if config.WhitelistPaths != nil {
			path := string(c.Path())
			if config.WhitelistPaths[path] {
				c.Next(ctx)
				return
			}
		}

		// 2. 获取租户ID
		tenantIDBytes := c.GetHeader(TenantIDHeader)
		if len(tenantIDBytes) == 0 {
			// 如果没有租户ID，跳过检查（可能公开端点）
			c.Next(ctx)
			return
		}
		tenantID := string(tenantIDBytes)

		// 3. 根据请求路径确定资源类型
		resourceType, requiredCount := determineResourceRequirement(c)
		if resourceType == "" {
			// 无法确定资源类型，跳过检查
			c.Next(ctx)
			return
		}

		// 4. 检查配额
		available, usage, err := checkQuotaAvailability(ctx, config.DB, config.QuotaRepo, tenantID, resourceType)
		if err != nil {
			logs.CtxErrorf(ctx, "[QuotaMiddleware] Failed to check quota: %v", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "Failed to check quota",
			})
			c.Abort()
			return
		}

		// 5. 计算操作后是否超过配额
		wouldExceed := (usage.UsedCount + requiredCount) > usage.MaxLimit

		// 6. 处理配额超限
		if wouldExceed {
			if config.WarningOnly {
				// 仅警告模式，记录日志但继续
				logs.CtxWarnf(ctx, "[QuotaMiddleware] Quota would be exceeded: tenant=%s, resource=%s, used=%d, max=%d, required=%d",
					tenantID, resourceType, usage.UsedCount, usage.MaxLimit, requiredCount)
				c.Next(ctx)
				return
			}

			if config.StrictMode {
				// 严格模式，拒绝请求
				logs.CtxWarnf(ctx, "[QuotaMiddleware] Quota exceeded, rejecting request: tenant=%s, resource=%s, used=%d, max=%d",
					tenantID, resourceType, usage.UsedCount, usage.MaxLimit)

				c.JSON(http.StatusForbidden, map[string]interface{}{
					"code":    403,
					"message": "Quota exceeded",
					"error": map[string]interface{}{
						"code":          "QUOTA_EXCEEDED",
						"resource_type": resourceType,
						"used":          usage.UsedCount,
						"max_limit":     usage.MaxLimit,
						"description":   fmt.Sprintf("You have exceeded your %s quota. Please upgrade your subscription.", resourceType),
					},
				})
				c.Abort()
				return
			}
		}

		// 7. 记录配额使用（异步）
		if available && !wouldExceed {
			go recordQuotaUsage(context.Background(), config.DB, tenantID, resourceType, requiredCount)
		}

		// 8. 继续处理请求
		c.Next(ctx)
	}
}

// determineResourceRequirement 根据请求确定资源类型和所需数量
func determineResourceRequirement(c *app.RequestContext) (entity.ResourceType, int) {
	path := string(c.Path())
	method := string(c.Method())

	// POST请求通常是创建资源，消耗1个配额
	if method == "POST" {
		switch {
		case path == "/api/v1/bots":
			return entity.ResourceTypeBots, 1
		default:
			return "", 0 // 无法确定
		}
	}

	// PUT/DELETE请求不增加配额（已有资源）
	return "", 0
}

// checkQuotaAvailability 检查配额可用性
func checkQuotaAvailability(
	ctx context.Context,
	db *gorm.DB,
	quotaRepo repository.QuotaRepository,
	tenantID string,
	resourceType entity.ResourceType,
) (bool, *entity.Quota, error) {
	// 1. 查询配额
	quota, err := quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return false, nil, err
	}
	if quota == nil {
		// 如果配额不存在，返回默认值（0限制）
		return false, &entity.Quota{
			TenantID:     tenantID,
			ResourceType: resourceType,
			UsedCount:    0,
			MaxLimit:     0,
		}, nil
	}

	// 2. 检查是否需要重置
	if quota.ResetCycle != entity.ResetCycleNever {
		if shouldResetQuota(quota) {
			// 执行重置
			now := time.Now().UnixMilli()
			if err := db.Model(&entity.Quota{}).
				Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
				Updates(map[string]interface{}{
					"used_count":    0,
					"last_reset_at": now,
					"updated_at":    now,
				}).Error; err != nil {
				return false, quota, err
			}
			quota.UsedCount = 0
			quota.LastResetAt = now
		}
	}

	// 3. 检查是否可用
	available := quota.UsedCount < quota.MaxLimit

	return available, quota, nil
}

// shouldResetQuota 判断是否应该重置配额
func shouldResetQuota(quota *entity.Quota) bool {
	now := time.Now()
	lastReset := time.Unix(quota.LastResetAt/1000, 0)

	switch quota.ResetCycle {
	case entity.ResetCycleDaily:
		// 每日重置：检查是否是新的一天
		return now.YearDay() != lastReset.YearDay() || now.Year() != lastReset.Year()

	case entity.ResetCycleMonthly:
		// 每月重置：检查是否是新月份
		return now.Month() != lastReset.Month() || now.Year() != lastReset.Year()

	case entity.ResetCycleYearly:
		// 每年重置：检查是否是新年份
		return now.Year() != lastReset.Year()

	default:
		return false
	}
}

// recordQuotaUsage 记录配额使用（异步）
func recordQuotaUsage(ctx context.Context, db *gorm.DB, tenantID string, resourceType entity.ResourceType, count int) {
	// 使用事务确保原子性
	err := db.Transaction(func(tx *gorm.DB) error {
		var quota entity.Quota
		if err := tx.Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
			First(&quota).Error; err != nil {
			return fmt.Errorf("failed to query quota: %w", err)
		}

		// 增加使用计数
		now := time.Now().UnixMilli()
		if err := tx.Model(&quota).
			Updates(map[string]interface{}{
				"used_count": quota.UsedCount + count,
				"updated_at": now,
			}).Error; err != nil {
			return fmt.Errorf("failed to update quota: %w", err)
		}

		logs.Infof("[QuotaMiddleware] Quota recorded: tenant=%s, resource=%s, count=%d, total=%d",
			tenantID, resourceType, count, quota.UsedCount+count)

		return nil
	})

	if err != nil {
		logs.Errorf("[QuotaMiddleware] Failed to record quota usage: %v", err)
	}
}

// QuotaChecker 配额检查器（供业务代码使用）
type QuotaChecker struct {
	db        *gorm.DB
	quotaRepo repository.QuotaRepository
}

// NewQuotaChecker 创建配额检查器
func NewQuotaChecker(db *gorm.DB, quotaRepo repository.QuotaRepository) *QuotaChecker {
	return &QuotaChecker{
		db:        db,
		quotaRepo: quotaRepo,
	}
}

// CheckAndConsume 检查并消耗配额（原子操作）
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - resourceType: 资源类型
//   - count: 需要的数量
//
// 返回:
//   - bool: 是否成功（配额足够且已消耗）
//   - error: 错误信息
func (qc *QuotaChecker) CheckAndConsume(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	count int,
) (bool, error) {
	var success bool

	err := qc.db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询配额（加锁）
		var quota entity.Quota
		if err := tx.Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
			First(&quota).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("quota not found for resource type: %s", resourceType)
			}
			return err
		}

		// 2. 检查是否需要重置
		if shouldResetQuota(&quota) {
			now := time.Now().UnixMilli()
			quota.UsedCount = 0
			quota.LastResetAt = now
		}

		// 3. 检查配额是否足够
		if quota.UsedCount+count > quota.MaxLimit {
			return fmt.Errorf("quota exceeded: resource_type=%s, used=%d, max=%d, required=%d",
				resourceType, quota.UsedCount, quota.MaxLimit, count)
		}

		// 4. 消耗配额
		now := time.Now().UnixMilli()
		if err := tx.Model(&quota).
			Updates(map[string]interface{}{
				"used_count":    quota.UsedCount + count,
				"last_reset_at": quota.LastResetAt,
				"updated_at":    now,
			}).Error; err != nil {
			return err
		}

		success = true
		logs.Infof("[QuotaChecker] Quota consumed: tenant=%s, resource=%s, count=%d, total=%d",
			tenantID, resourceType, count, quota.UsedCount+count)

		return nil
	})

	if err != nil {
		return false, err
	}

	return success, nil
}

// CheckOnly 仅检查配额，不消耗
func (qc *QuotaChecker) CheckOnly(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	count int,
) (bool, *QuotaStatus, error) {
	quota, err := qc.quotaRepo.GetByTenantAndResource(ctx, tenantID, resourceType)
	if err != nil {
		return false, nil, err
	}
	if quota == nil {
		return false, nil, fmt.Errorf("quota not found for resource type: %s", resourceType)
	}

	status := &QuotaStatus{
		ResourceType:    quota.ResourceType,
		UsedCount:       quota.UsedCount,
		MaxLimit:        quota.MaxLimit,
		RemainingCount:  quota.MaxLimit - quota.UsedCount,
		UsagePercentage: float64(quota.UsedCount) / float64(quota.MaxLimit) * 100,
		ResetCycle:      quota.ResetCycle,
		LastResetAt:     quota.LastResetAt,
	}

	available := quota.UsedCount+count <= quota.MaxLimit
	return available, status, nil
}

// RollbackConsumption 回滚配额消耗
//
// 场景：业务操作失败时，需要回滚已消耗的配额
func (qc *QuotaChecker) RollbackConsumption(
	ctx context.Context,
	tenantID string,
	resourceType entity.ResourceType,
	count int,
) error {
	return qc.db.Transaction(func(tx *gorm.DB) error {
		var quota entity.Quota
		if err := tx.Where("tenant_id = ? AND resource_type = ?", tenantID, resourceType).
			First(&quota).Error; err != nil {
			return err
		}

		// 确保不会回滚成负数
		newCount := quota.UsedCount - count
		if newCount < 0 {
			newCount = 0
		}

		now := time.Now().UnixMilli()
		if err := tx.Model(&quota).
			Updates(map[string]interface{}{
				"used_count": newCount,
				"updated_at": now,
			}).Error; err != nil {
			return err
		}

		logs.Infof("[QuotaChecker] Quota rolled back: tenant=%s, resource=%s, count=%d, total=%d",
			tenantID, resourceType, count, newCount)

		return nil
	})
}

// QuotaStatus 配额状态
type QuotaStatus struct {
	ResourceType    entity.ResourceType `json:"resource_type"`
	UsedCount       int                 `json:"used_count"`
	MaxLimit        int                 `json:"max_limit"`
	RemainingCount  int                 `json:"remaining_count"`
	UsagePercentage float64             `json:"usage_percentage"`
	ResetCycle      entity.ResetCycle   `json:"reset_cycle"`
	LastResetAt     int64               `json:"last_reset_at"`
}
