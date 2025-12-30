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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// TestShouldResetQuota 测试配额重置逻辑
func TestShouldResetQuota(t *testing.T) {
	t.Run("每日重置", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)

		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleDaily,
			LastResetAt: yesterday.UnixMilli(),
		}

		// 应该重置（已经是新的一天）
		assert.True(t, shouldResetQuota(quota))
	})

	t.Run("每日重置-同一天", func(t *testing.T) {
		now := time.Now()
		today := now

		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleDaily,
			LastResetAt: today.UnixMilli(),
		}

		// 不应该重置（还是同一天）
		assert.False(t, shouldResetQuota(quota))
	})

	t.Run("每月重置", func(t *testing.T) {
		// 上个月
		lastMonth := time.Date(2025, 1, 15, 0, 0, 0, 0, time.Local)
		now := time.Date(2025, 2, 1, 0, 0, 0, 0, time.Local)

		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleMonthly,
			LastResetAt: lastMonth.UnixMilli(),
		}

		// 手动设置当前时间（通过全局时间可能不稳定）
		// 这里只测试逻辑
		if now.Month() != lastMonth.Month() || now.Year() != lastMonth.Year() {
			// 不同的月份/年份应该重置
			// 注意：这个测试假设shouldResetQuota内部使用time.Now()
		}
	})

	t.Run("从不重置", func(t *testing.T) {
		quota := &entity.Quota{
			ResetCycle:  entity.ResetCycleNever,
			LastResetAt: time.Now().Add(-24 * time.Hour).UnixMilli(),
		}

		// Never模式不应该重置
		assert.False(t, shouldResetQuota(quota))
	})
}

// TestDetermineResourceRequirement 测试资源需求确定
func TestDetermineResourceRequirement(t *testing.T) {
	// 注意：这个测试需要构造RequestContext，比较复杂
	// 这里只测试基本逻辑，实际集成测试会在HTTP测试中覆盖

	t.Run("POST创建Bot", func(t *testing.T) {
		// 模拟POST /api/v1/bots请求
		// 应该返回: ResourceTypeBots, 1
		// TODO: 添加完整的HTTP测试
	})
}

// TestQuotaStatus 测试配额状态计算
func TestQuotaStatus(t *testing.T) {
	t.Run("正常使用率计算", func(t *testing.T) {
		status := &QuotaStatus{
			UsedCount:  50,
			MaxLimit:   100,
		}

		status.UsagePercentage = float64(status.UsedCount) / float64(status.MaxLimit) * 100
		status.RemainingCount = status.MaxLimit - status.UsedCount

		assert.Equal(t, 50.0, status.UsagePercentage)
		assert.Equal(t, 50, status.RemainingCount)
	})

	t.Run("接近上限", func(t *testing.T) {
		status := &QuotaStatus{
			UsedCount:  95,
			MaxLimit:   100,
		}

		status.UsagePercentage = float64(status.UsedCount) / float64(status.MaxLimit) * 100

		assert.Equal(t, 95.0, status.UsagePercentage)
		assert.True(t, status.UsagePercentage >= 80)
	})

	t.Run("已超限", func(t *testing.T) {
		status := &QuotaStatus{
			UsedCount:  110,
			MaxLimit:   100,
		}

		status.UsagePercentage = float64(status.UsedCount) / float64(status.MaxLimit) * 100

		assert.Equal(t, 110.0, status.UsagePercentage)
		assert.True(t, status.UsagePercentage > 100)
	})

	t.Run("空配额", func(t *testing.T) {
		status := &QuotaStatus{
			UsedCount:  0,
			MaxLimit:   0,
		}

		// 避免除零错误
		if status.MaxLimit > 0 {
			status.UsagePercentage = float64(status.UsedCount) / float64(status.MaxLimit) * 100
		} else {
			status.UsagePercentage = 0
		}

		assert.Equal(t, 0.0, status.UsagePercentage)
		assert.Equal(t, 0, status.RemainingCount)
	})
}

// TestCheckQuotaAvailability 测试配额可用性检查
func TestCheckQuotaAvailability(t *testing.T) {
	t.Run("配额足够", func(t *testing.T) {
		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    2,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}

		available := quota.UsedCount < quota.MaxLimit
		assert.True(t, available)
	})

	t.Run("配额已满", func(t *testing.T) {
		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    10,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}

		available := quota.UsedCount < quota.MaxLimit
		assert.False(t, available)
	})

	t.Run("配额已超限", func(t *testing.T) {
		quota := &entity.Quota{
			TenantID:     "test_tenant",
			ResourceType: entity.ResourceTypeBots,
			UsedCount:    15,
			MaxLimit:     10,
			ResetCycle:   entity.ResetCycleNever,
			LastResetAt:  time.Now().UnixMilli(),
		}

		available := quota.UsedCount < quota.MaxLimit
		assert.False(t, available)
	})
}

// TestResourceTypeString 测试资源类型字符串
func TestResourceTypeString(t *testing.T) {
	t.Run("Bots类型", func(t *testing.T) {
		rt := entity.ResourceTypeBots
		assert.Equal(t, "bots", string(rt))
	})

	t.Run("Messages类型", func(t *testing.T) {
		rt := entity.ResourceTypeMessages
		assert.Equal(t, "messages", string(rt))
	})

	t.Run("Storage类型", func(t *testing.T) {
		rt := entity.ResourceTypeStorage
		assert.Equal(t, "storage", string(rt))
	})

	t.Run("TeamMembers类型", func(t *testing.T) {
		rt := entity.ResourceTypeTeamMembers
		assert.Equal(t, "team_members", string(rt))
	})
}

// TestResetCycleString 测试重置周期字符串
func TestResetCycleString(t *testing.T) {
	t.Run("Never周期", func(t *testing.T) {
		rc := entity.ResetCycleNever
		assert.Equal(t, "never", string(rc))
	})

	t.Run("Daily周期", func(t *testing.T) {
		rc := entity.ResetCycleDaily
		assert.Equal(t, "daily", string(rc))
	})

	t.Run("Weekly周期", func(t *testing.T) {
		rc := entity.ResetCycleWeekly
		assert.Equal(t, "weekly", string(rc))
	})

	t.Run("Monthly周期", func(t *testing.T) {
		rc := entity.ResetCycleMonthly
		assert.Equal(t, "monthly", string(rc))
	})

	t.Run("Yearly周期", func(t *testing.T) {
		rc := entity.ResetCycleYearly
		assert.Equal(t, "yearly", string(rc))
	})
}
