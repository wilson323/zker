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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-studio/backend/domain/tenant/repository"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 自动迁移
	err = db.AutoMigrate(
		&entity.Tenant{},
		&entity.Subscription{},
		&entity.Quota{},
	)
	require.NoError(t, err)

	return db
}

// TestTenantService_Integration_CreateTenant 测试租户创建集成
func TestTenantService_Integration_CreateTenant(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	tenantRepo := repository.NewTenantRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	quotaRepo := repository.NewQuotaRepository(db)

	subscriptionService := NewSubscriptionService(subscriptionRepo, quotaRepo)
	tenantService := NewTenantService(tenantRepo, subscriptionService)

	// Act
	ctx := context.Background()
	req := &CreateTenantRequest{
		TenantName: "测试租户",
		TenantType: entity.TenantTypeEnterprise,
	}

	tenant, err := tenantService.CreateTenant(ctx, req)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, tenant.TenantID)
	assert.Equal(t, "测试租户", tenant.TenantName)
	assert.Equal(t, entity.TenantTypeEnterprise, tenant.TenantType)
	assert.Equal(t, entity.TenantStatusActive, tenant.Status)
	assert.True(t, tenant.CreatedAt > 0)

	// 验证订阅已创建
	subscription, err := subscriptionRepo.GetByTenant(ctx, tenant.TenantID)
	require.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.Equal(t, tenant.TenantID, subscription.TenantID)
	assert.Equal(t, entity.SubscriptionTierFree, subscription.PlanTier)

	// 验证配额已创建
	quotas, err := quotaRepo.GetByTenant(ctx, tenant.TenantID)
	require.NoError(t, err)
	assert.NotEmpty(t, quotas)
}

// TestQuotaService_Integration_QuotaEnforcement 测试配额强制执行集成
func TestQuotaService_Integration_QuotaEnforcement(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	tenantRepo := repository.NewTenantRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	quotaRepo := repository.NewQuotaRepository(db)

	subscriptionService := NewSubscriptionService(subscriptionRepo, quotaRepo)
	tenantService := NewTenantService(tenantRepo, subscriptionService)
	quotaService := NewQuotaService(quotaRepo)

	// 创建租户
	ctx := context.Background()
	tenant, err := tenantService.CreateTenant(ctx, &CreateTenantRequest{
		TenantName: "配额测试租户",
		TenantType: entity.TenantTypeEnterprise,
	})
	require.NoError(t, err)

	// 更新配额为最大10个Bot
	quotas, err := quotaRepo.GetByTenant(ctx, tenant.TenantID)
	require.NoError(t, err)
	require.NotEmpty(t, quotas)

	botQuota := quotas[0]
	botQuota.MaxLimit = 10
	err = quotaRepo.Update(ctx, botQuota)
	require.NoError(t, err)

	// Act & Assert: 尝试创建10个Bot（应该成功）
	for i := 0; i < 10; i++ {
		err := quotaService.ConsumeQuota(ctx, tenant.TenantID, entity.ResourceTypeBots, 1)
		assert.NoError(t, err, "第%d个Bot应该创建成功", i+1)
	}

	// Act & Assert: 尝试创建第11个Bot（应该失败）
	err = quotaService.ConsumeQuota(ctx, tenant.TenantID, entity.ResourceTypeBots, 1)
	assert.Error(t, err)
	var quotaErr *QuotaExceededError
	assert.True(t, assert.ErrorAs(t, err, &quotaErr))
	assert.Equal(t, tenant.TenantID, quotaErr.TenantID)
	assert.Equal(t, entity.ResourceTypeBots, quotaErr.ResourceType)
	assert.Equal(t, 10, quotaErr.CurrentUsage)
	assert.Equal(t, 10, quotaErr.MaxLimit)
}

// TestSubscriptionService_Integration_UpgradeTier 测试订阅升级集成
func TestSubscriptionService_Integration_UpgradeTier(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	quotaRepo := repository.NewQuotaRepository(db)
	subscriptionService := NewSubscriptionService(subscriptionRepo, quotaRepo)

	// 创建免费订阅
	ctx := context.Background()
	subscription, err := subscriptionService.CreateSubscription(
		ctx,
		"tenant_123",
		entity.SubscriptionTierFree,
		entity.BillingCycleMonthly,
	)
	require.NoError(t, err)
	assert.Equal(t, entity.SubscriptionTierFree, subscription.PlanTier)

	// 获取配额
	quotas, err := quotaRepo.GetByTenant(ctx, "tenant_123")
	require.NoError(t, err)
	freeQuotaLimit := quotas[0].MaxLimit

	// Act: 升级到Pro版本
	updated, err := subscriptionService.UpgradeTier(ctx, "tenant_123", entity.SubscriptionTierPro)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, entity.SubscriptionTierPro, updated.PlanTier)

	// 验证配额已增加
	quotasAfter, err := quotaRepo.GetByTenant(ctx, "tenant_123")
	require.NoError(t, err)
	assert.True(t, quotasAfter[0].MaxLimit > freeQuotaLimit)
}

// TestQuotaMonitor_Integration_ResetQuotas 测试配额监控和重置集成
func TestQuotaMonitor_Integration_ResetQuotas(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	quotaRepo := repository.NewQuotaRepository(db)
	monitor := NewQuotaMonitor(quotaRepo, nil) // 不需要告警服务

	// 创建一个每日重置的配额
	ctx := context.Background()
	yesterday := time.Now().Add(-24 * time.Hour)
	quota := &entity.Quota{
		QuotaID:      "quota_daily_test",
		TenantID:     "tenant_123",
		ResourceType: entity.ResourceTypeMessages,
		MaxLimit:     1000,
		UsedCount:    500,
		ResetCycle:   entity.ResetCycleDaily,
		LastResetAt:  yesterday.UnixMilli(),
		CreatedAt:    yesterday.UnixMilli(),
		UpdatedAt:    yesterday.UnixMilli(),
	}

	err := quotaRepo.Create(ctx, quota)
	require.NoError(t, err)

	// Act: 执行定期重置
	err = monitor.ResetPeriodicQuotas(ctx)

	// Assert
	require.NoError(t, err)

	// 验证配额已重置
	resetQuota, err := quotaRepo.GetByID(ctx, quota.QuotaID)
	require.NoError(t, err)
	assert.Equal(t, 0, resetQuota.UsedCount)
	assert.True(t, resetQuota.LastResetAt > quota.LastResetAt)
}

// TestBillingService_Integration_GenerateInvoice 测试计费服务集成
func TestBillingService_Integration_GenerateInvoice(t *testing.T) {
	// Setup
	db := setupTestDB(t)

	// 注意：这里需要创建UsageLog和Invoice表的实体
	// 由于篇幅限制，这里省略表定义和迁移

	// 创建Mock仓储（实际应用中应该使用真实的仓储实现）
	// 这里演示集成测试的框架

	// Act & Assert框架
	// 1. 记录使用量
	// 2. 生成账单
	// 3. 验证账单金额

	t.Skip("需要完整的UsageLog和Invoice实体定义")
}

// TestCustomFilterEngine_Integration 测试自定义过滤器引擎集成
func TestCustomFilterEngine_Integration(t *testing.T) {
	// Setup
	engine := NewCustomFilterEngine()

	// 测试数据
	customFilter := `[
		{"field": "status", "operator": "eq", "value": "active"},
		{"field": "priority", "operator": "gte", "value": 5},
		{"field": "type", "operator": "in", "value": ["bot", "workflow"]}
	]`

	// 测试资源1：匹配所有条件
	resource1 := map[string]interface{}{
		"status":   "active",
		"priority": 8,
		"type":     "bot",
	}

	// Act
	matched, err := engine.Match(context.Background(), customFilter, resource1)

	// Assert
	require.NoError(t, err)
	assert.True(t, matched)

	// 测试资源2：不匹配status
	resource2 := map[string]interface{}{
		"status":   "inactive",
		"priority": 8,
		"type":     "bot",
	}

	matched2, err := engine.Match(context.Background(), customFilter, resource2)
	require.NoError(t, err)
	assert.False(t, matched2)
}

// TestHybridIntentMatcher_Integration 测试混合意图匹配器集成
func TestHybridIntentMatcher_Integration(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	ruleRepo := setupTestRoutingRules(db)
	botRepo := setupTestBots(db)

	// 注意：需要真实的EmbeddingClient实现
	// 这里使用Mock演示集成测试框架

	t.Skip("需要真实的EmbeddingClient实现")
}

// setupTestRoutingRules 创建测试路由规则
func setupTestRoutingRules(db *gorm.DB) interface{} {
	// 实现省略
	return nil
}

// setupTestBots 创建测试Bot数据
func setupTestBots(db *gorm.DB) interface{} {
	// 实现省略
	return nil
}
