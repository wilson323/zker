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

package billing_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

/**
 * 预算管理性能基准测试
 *
 * 性能目标：
 * - 单租户预算检查: < 50ms
 * - 批量检查(100租户): < 1秒
 * - 使用统计查询: < 100ms
 * - 告警历史查询: < 150ms
 *
 * 设计原则：
 * - SOLID: 单一职责
 * - DRY: 复用测试数据生成器
 * - KISS: 保持简单
 *
 * @author 研发B
 * @date 2025-12-30
 */

// ==================== Mock Repository ====================

/**
 * MockBudgetAlertRepository Mock告警仓储
 */
type MockBudgetAlertRepository struct {
	alerts   []*entity.BudgetAlert
	alertMap map[string]bool // tenantID:alertType:date -> sent
	mu       sync.RWMutex
}

func NewMockBudgetAlertRepository() *MockBudgetAlertRepository {
	return &MockBudgetAlertRepository{
		alerts:   make([]*entity.BudgetAlert, 0),
		alertMap: make(map[string]bool),
	}
}

func (r *MockBudgetAlertRepository) Create(ctx context.Context, alert *entity.BudgetAlert) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	alert.ID = uint64(len(r.alerts) + 1)
	r.alerts = append(r.alerts, alert)
	return nil
}

func (r *MockBudgetAlertRepository) CheckAlertSentToday(
	ctx context.Context,
	tenantID, alertType, date string,
) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%s:%s", tenantID, alertType, date)
	return r.alertMap[key], nil
}

func (r *MockBudgetAlertRepository) GetByTenantID(
	ctx context.Context,
	tenantID string,
	limit int,
) ([]*entity.BudgetAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*entity.BudgetAlert, 0)
	for _, alert := range r.alerts {
		if alert.TenantID == tenantID {
			result = append(result, alert)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

// ==================== 测试数据生成器 ====================

/**
 * 生成测试预算配置
 */
func generateBudgetSettings(tenantID string, budgetAmount float64) *entity.BudgetSettings {
	hardCapAmount := budgetAmount * 1.1 // 允许超10%
	return &entity.BudgetSettings{
		TenantID:              tenantID,
		BudgetType:            entity.BudgetTypeMonthly,
		BudgetAmount:          budgetAmount,
		Currency:              "CNY",
		AlertThreshold1:       80,
		AlertThreshold2:       95,
		HardCapEnabled:        true,
		HardCapAmount:         &hardCapAmount,
		NotificationChannels:  `["email", "webhook"]`,
		NotificationRecipients: `["admin@example.com"]`,
	}
}

/**
 * 生成测试预算告警
 */
func generateBudgetAlert(tenantID string, alertType string, usedAmount, budgetAmount float64) *entity.BudgetAlert {
	usagePercent := (usedAmount / budgetAmount) * 100

	alertLevel := entity.AlertLevelWarning
	if alertType == entity.AlertTypeThreshold2 {
		alertLevel = entity.AlertLevelCritical
	} else if alertType == entity.AlertTypeHardCap {
		alertLevel = entity.AlertLevelEmergency
	}

	message := fmt.Sprintf("预算告警: 已使用 %.2f%%", usagePercent)

	return &entity.BudgetAlert{
		TenantID:              tenantID,
		AlertType:             alertType,
		BudgetAmount:          budgetAmount,
		UsedAmount:            usedAmount,
		UsagePercent:          usagePercent,
		AlertLevel:            alertLevel,
		AlertMessage:          &message,
		NotificationChannels:  `["email"]`,
		NotificationStatus:    entity.NotificationStatusSent,
		SentAt:                &[]time.Time{time.Now()}[0],
	}
}

// ==================== 1. 预算检查性能测试 ====================

/**
 * BenchmarkBudgetAlertService_CheckBudget 单租户预算检查性能
 *
 * 性能目标：
 * - ops/s: > 20 次/秒
 * - 每次检查: < 50ms
 */
func BenchmarkBudgetAlertService_CheckBudget(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil, // logger
	)

	ctx := context.Background()

	// 预先创建预算配置
	budget := generateBudgetSettings("tenant-1", 10000.0)
	budgetRepo.Create(ctx, budget)

	// 预先创建使用汇总数据（模拟已使用70%）
	now := time.Now()
	summary := &entity.TokenUsageSummary{
		TenantID:    "tenant-1",
		SummaryDate: now.Format("2006-01-02"),
		TotalCost:   7000.0, // 已使用70%
	}
	summaryRepo.Upsert(ctx, summary)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := budgetAlertSvc.CheckBudget(ctx, "tenant-1")
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetAlertService_CheckBudget_Threshold1 一级告警检查性能
 *
 * 性能目标：
 * - ops/s: > 20 次/秒
 * - 每次检查: < 50ms
 */
func BenchmarkBudgetAlertService_CheckBudget_Threshold1(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
	)

	ctx := context.Background()

	budget := generateBudgetSettings("tenant-1", 10000.0)
	budgetRepo.Create(ctx, budget)

	// 模拟使用82%（触发一级告警）
	now := time.Now()
	summary := &entity.TokenUsageSummary{
		TenantID:    "tenant-1",
		SummaryDate: now.Format("2006-01-02"),
		TotalCost:   8200.0,
	}
	summaryRepo.Upsert(ctx, summary)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := budgetAlertSvc.CheckBudget(ctx, "tenant-1")
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetAlertService_CheckBudget_Threshold2 二级告警检查性能
 *
 * 性能目标：
 * - ops/s: > 20 次/秒
 * - 每次检查: < 50ms
 */
func BenchmarkBudgetAlertService_CheckBudget_Threshold2(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
	)

	ctx := context.Background()

	budget := generateBudgetSettings("tenant-1", 10000.0)
	budgetRepo.Create(ctx, budget)

	// 模拟使用96%（触发二级告警）
	now := time.Now()
	summary := &entity.TokenUsageSummary{
		TenantID:    "tenant-1",
		SummaryDate: now.Format("2006-01-02"),
		TotalCost:   9600.0,
	}
	summaryRepo.Upsert(ctx, summary)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := budgetAlertSvc.CheckBudget(ctx, "tenant-1")
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetAlertService_CheckBatch_100 批量检查100个租户性能
 *
 * 性能目标：
 * - 批量大小: 100个租户
 * - 总耗时: < 1秒
 * - 平均每个: < 10ms
 */
func BenchmarkBudgetAlertService_CheckBatch_100(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
	)

	ctx := context.Background()
	tenantCount := 100

	// 预先创建100个租户的预算配置和使用数据
	now := time.Now()
	for i := 0; i < tenantCount; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		budget := generateBudgetSettings(tenantID, 10000.0)
		budgetRepo.Create(ctx, budget)

		summary := &entity.TokenUsageSummary{
			TenantID:    tenantID,
			SummaryDate: now.Format("2006-01-02"),
			TotalCost:   7000.0 + float64(i%30)*100, // 70% - 99%使用率
		}
		summaryRepo.Upsert(ctx, summary)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 逐个检查所有租户
		for j := 0; j < tenantCount; j++ {
			tenantID := fmt.Sprintf("tenant-%d", j)
			err := budgetAlertSvc.CheckBudget(ctx, tenantID)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

/**
 * BenchmarkBudgetAlertService_CheckBudget_Concurrent 并发预算检查性能
 *
 * 性能目标：
 * - 并发数: 50
 * - 总ops/s: > 1,000
 * - 平均延迟: < 50ms
 */
func BenchmarkBudgetAlertService_CheckBudget_Concurrent(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
	)

	ctx := context.Background()

	// 预先创建100个租户的数据
	now := time.Now()
	for i := 0; i < 100; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		budget := generateBudgetSettings(tenantID, 10000.0)
		budgetRepo.Create(ctx, budget)

		summary := &entity.TokenUsageSummary{
			TenantID:    tenantID,
			SummaryDate: now.Format("2006-01-02"),
			TotalCost:   7000.0,
		}
		summaryRepo.Upsert(ctx, summary)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			tenantID := fmt.Sprintf("tenant-%d", i%100)
			err := budgetAlertSvc.CheckBudget(ctx, tenantID)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

// ==================== 2. 使用统计查询性能测试 ====================

/**
 * BenchmarkBudgetManagement_GetUsage 获取预算使用情况性能
 *
 * 性能目标：
 * - ops/s: > 10 次/秒
 * - 每次查询: < 100ms
 */
func BenchmarkBudgetManagement_GetUsage(b *testing.B) {
	summaryRepo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	// 预先生成30天的汇总数据
	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		summary := &entity.TokenUsageSummary{
			TenantID:         "tenant-1",
			SummaryDate:      date,
			TotalInputTokens: 100000,
			TotalOutputTokens: 50000,
			TotalTokens:      150000,
			TotalCost:        450.0,
			TotalRequests:    1000,
			CachedRequests:   100,
		}
		summaryRepo.Upsert(ctx, summary)
	}

	startDate := now.AddDate(0, 0, -29).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10+1)
		_, err := summaryRepo.GetTotalCost(ctx, tenantID, now.AddDate(0, 0, -29), now)
		if err != nil {
			b.Fatal(err)
		}

		// 获取日期范围内的汇总数据
		_, err = summaryRepo.GetByDateRange(ctx, tenantID, nil, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetManagement_GetDailyUsage 每日使用趋势查询性能
 *
 * 性能目标：
 * - 查询天数: 30天
 * - 总耗时: < 200ms
 */
func BenchmarkBudgetManagement_GetDailyUsage(b *testing.B) {
	summaryRepo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	// 预先生成30天的汇总数据
	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		summary := &entity.TokenUsageSummary{
			TenantID:         "tenant-1",
			SummaryDate:      date,
			TotalInputTokens: 100000,
			TotalOutputTokens: 50000,
			TotalTokens:      150000,
			TotalCost:        450.0,
			TotalRequests:    1000,
			CachedRequests:   100,
		}
		summaryRepo.Upsert(ctx, summary)
	}

	startDate := now.AddDate(0, 0, -29).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10+1)
		_, err := summaryRepo.GetByDateRange(ctx, tenantID, nil, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== 3. 告警历史查询性能测试 ====================

/**
 * BenchmarkBudgetManagement_GetAlerts 告警历史查询性能
 *
 * 性能目标：
 * - ops/s: > 10 次/秒
 * - 每次查询: < 150ms
 */
func BenchmarkBudgetManagement_GetAlerts(b *testing.B) {
	alertRepo := NewMockBudgetAlertRepository()
	ctx := context.Background()

	// 预先生成100条告警历史
	now := time.Now()
	for i := 0; i < 100; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		alertType := []string{
			entity.AlertTypeThreshold1,
			entity.AlertTypeThreshold2,
			entity.AlertTypeHardCap,
		}[i % 3]
		usedAmount := 8000.0 + float64(i%20)*100
		alert := generateBudgetAlert(tenantID, alertType, usedAmount, 10000.0)
		alert.CreatedAt = now.Add(time.Duration(-i) * time.Hour)
		alertRepo.Create(ctx, alert)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		_, err := alertRepo.GetByTenantID(ctx, tenantID, 20)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetManagement_CheckAlertSent 检查告警发送状态性能
 *
 * 性能目标：
 * - ops/s: > 50,000 次/秒
 * - ns/op: < 20,000 ns/op
 */
func BenchmarkBudgetManagement_CheckAlertSent(b *testing.B) {
	alertRepo := NewMockBudgetAlertRepository()
	ctx := context.Background()

	// 预先标记一些告警为已发送
	today := time.Now().Format("2006-01-02")
	for i := 0; i < 10; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		alertType := entity.AlertTypeThreshold1
		key := fmt.Sprintf("%s:%s:%s", tenantID, alertType, today)
		alertRepo.alertMap[key] = true
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%20)
		alertType := []string{
			entity.AlertTypeThreshold1,
			entity.AlertTypeThreshold2,
			entity.AlertTypeHardCap,
		}[i % 3]
		_, err := alertRepo.CheckAlertSentToday(ctx, tenantID, alertType, today)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== 4. 预算配置CRUD性能测试 ====================

/**
 * BenchmarkBudgetManagement_Create 创建预算配置性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次操作: < 10ms
 */
func BenchmarkBudgetManagement_Create(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		budget := generateBudgetSettings(tenantID, 10000.0)
		err := budgetRepo.Create(ctx, budget)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetManagement_Update 更新预算配置性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次操作: < 10ms
 */
func BenchmarkBudgetManagement_Update(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	// 预先创建预算配置
	for i := 0; i < 100; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		budget := generateBudgetSettings(tenantID, 10000.0)
		budgetRepo.Create(ctx, budget)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%100)
		budget, _ := budgetRepo.GetByTenantID(ctx, tenantID)
		if budget != nil {
			budget.BudgetAmount = 15000.0
			err := budgetRepo.Update(ctx, budget)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

/**
 * BenchmarkBudgetManagement_Get 获取预算配置性能
 *
 * 性能目标：
 * - ops/s: > 10,000 次/秒
 * - ns/op: < 100,000 ns/op
 */
func BenchmarkBudgetManagement_Get(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	// 预先创建预算配置
	for i := 0; i < 100; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		budget := generateBudgetSettings(tenantID, 10000.0)
		budgetRepo.Create(ctx, budget)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%100)
		_, err := budgetRepo.GetByTenantID(ctx, tenantID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== 5. 综合场景性能测试 ====================

/**
 * BenchmarkBudgetManagement_CompleteFlow 预算管理完整流程性能
 *
 * 场景：创建预算 → 记录使用 → 检查预算 → 生成告警 → 查询历史
 *
 * 性能目标：
 * - 完整流程: < 300ms
 */
func BenchmarkBudgetManagement_CompleteFlow(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
	)

	ctx := context.Background()
	now := time.Now()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%100)

		// 1. 创建预算配置
		budget := generateBudgetSettings(tenantID, 10000.0)
		err := budgetRepo.Create(ctx, budget)
		if err != nil {
			b.Fatal(err)
		}

		// 2. 记录使用（通过汇总）
		summary := &entity.TokenUsageSummary{
			TenantID:         tenantID,
			SummaryDate:      now.Format("2006-01-02"),
			TotalInputTokens: 100000,
			TotalOutputTokens: 50000,
			TotalTokens:      150000,
			TotalCost:        8200.0, // 82%使用率
			TotalRequests:    1000,
			CachedRequests:   100,
		}
		err = summaryRepo.Upsert(ctx, summary)
		if err != nil {
			b.Fatal(err)
		}

		// 3. 检查预算
		err = budgetAlertSvc.CheckBudget(ctx, tenantID)
		if err != nil {
			b.Fatal(err)
		}

		// 4. 查询告警历史
		_, err = alertRepo.GetByTenantID(ctx, tenantID, 10)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkBudgetManagement_AllTenantsCheck 全租户预算检查性能
 *
 * 场景：定时任务检查所有租户的预算状态
 *
 * 性能目标：
 * - 租户数: 1000个
 * - 总耗时: < 5秒
 * - 平均每个: < 5ms
 */
func BenchmarkBudgetManagement_AllTenantsCheck(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{} // Mock

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
	)

	ctx := context.Background()
	tenantCount := 1000
	now := time.Now()

	// 预先创建1000个租户的数据
	for i := 0; i < tenantCount; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		budget := generateBudgetSettings(tenantID, 10000.0)
		budgetRepo.Create(ctx, budget)

		summary := &entity.TokenUsageSummary{
			TenantID:    tenantID,
			SummaryDate: now.Format("2006-01-02"),
			TotalCost:   7000.0 + float64(i%30)*100,
		}
		summaryRepo.Upsert(ctx, summary)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 模拟定时任务检查所有租户
		for j := 0; j < tenantCount; j++ {
			tenantID := fmt.Sprintf("tenant-%d", j)
			err := budgetAlertSvc.CheckBudget(ctx, tenantID)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
