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
)

/**
 * 数据库Repository性能基准测试
 *
 * 性能目标：
 * - 写入操作: < 5ms per operation
 * - 批量写入(100条): < 50ms
 * - 查询操作: < 10ms
 * - 聚合查询: < 100ms
 *
 * 设计原则：
 * - SOLID: 单一职责
 * - DRY: 复用测试数据生成器
 * - KISS: 保持简单
 *
 * @author 研发B
 * @date 2025-12-30
 */

// ==================== 1. TokenUsageLog Repository性能测试 ====================

/**
 * BenchmarkRepository_CreateUsageLog 创建TokenUsageLog性能
 *
 * 性能目标：
 * - ops/s: > 200 次/秒
 * - 每次操作: < 5ms
 */
func BenchmarkRepository_CreateUsageLog(b *testing.B) {
	repo := NewMockTokenUsageLogRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		userID := uint64(1000 + i)
		botID := fmt.Sprintf("bot-%d", i%10)
		log := &entity.TokenUsageLog{
			TenantID:       fmt.Sprintf("tenant-%d", i%10),
			UserID:         &userID,
			BotID:          &botID,
			InputTokens:    1000,
			OutputTokens:   500,
			TotalTokens:    1500,
			ModelProvider:  "openai",
			ModelName:      "gpt-4",
			UnitPrice:      0.03,
			InputCost:      0.03,
			OutputCost:     0.03,
			TotalCost:      0.06,
			RequestType:    entity.RequestTypeChat,
			ResponseTimeMs: intPtr(150),
		}

		err := repo.Create(ctx, log)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_BatchCreateUsageLog_100 批量创建100条TokenUsageLog性能
 *
 * 性能目标：
 * - 批量大小: 100条
 * - 总耗时: < 50ms
 * - 平均每条: < 0.5ms
 */
func BenchmarkRepository_BatchCreateUsageLog_100(b *testing.B) {
	repo := NewMockTokenUsageLogRepository()
	ctx := context.Background()
	batchSize := 100

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logs := make([]*entity.TokenUsageLog, batchSize)
		for j := 0; j < batchSize; j++ {
			userID := uint64(1000 + i*batchSize + j)
			botID := fmt.Sprintf("bot-%d", j%10)
			logs[j] = &entity.TokenUsageLog{
				TenantID:       fmt.Sprintf("tenant-%d", j%10),
				UserID:         &userID,
				BotID:          &botID,
				InputTokens:    1000,
				OutputTokens:   500,
				TotalTokens:    1500,
				ModelProvider:  "openai",
				ModelName:      "gpt-4",
				UnitPrice:      0.03,
				InputCost:      0.03,
				OutputCost:     0.03,
				TotalCost:      0.06,
				RequestType:    entity.RequestTypeChat,
				ResponseTimeMs: intPtr(150),
			}
		}

		err := repo.BatchCreate(ctx, logs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_BatchCreateUsageLog_1000 批量创建1000条TokenUsageLog性能
 *
 * 性能目标：
 * - 批量大小: 1000条
 * - 总耗时: < 500ms
 * - 平均每条: < 0.5ms
 */
func BenchmarkRepository_BatchCreateUsageLog_1000(b *testing.B) {
	repo := NewMockTokenUsageLogRepository()
	ctx := context.Background()
	batchSize := 1000

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logs := make([]*entity.TokenUsageLog, batchSize)
		for j := 0; j < batchSize; j++ {
			userID := uint64(1000 + i*batchSize + j)
			botID := fmt.Sprintf("bot-%d", j%10)
			logs[j] = &entity.TokenUsageLog{
				TenantID:       fmt.Sprintf("tenant-%d", j%10),
				UserID:         &userID,
				BotID:          &botID,
				InputTokens:    1000,
				OutputTokens:   500,
				TotalTokens:    1500,
				ModelProvider:  "openai",
				ModelName:      "gpt-4",
				UnitPrice:      0.03,
				InputCost:      0.03,
				OutputCost:     0.03,
				TotalCost:      0.06,
				RequestType:    entity.RequestTypeChat,
				ResponseTimeMs: intPtr(150),
			}
		}

		err := repo.BatchCreate(ctx, logs)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_CreateUsageLog_Concurrent 并发创建TokenUsageLog性能
 *
 * 性能目标：
 * - 并发数: 100
 * - 总ops/s: > 10,000
 * - 平均延迟: < 10ms
 */
func BenchmarkRepository_CreateUsageLog_Concurrent(b *testing.B) {
	repo := NewMockTokenUsageLogRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			userID := uint64(1000 + i)
			botID := fmt.Sprintf("bot-%d", i%10)
			log := &entity.TokenUsageLog{
				TenantID:       fmt.Sprintf("tenant-%d", i%10),
				UserID:         &userID,
				BotID:          &botID,
				InputTokens:    1000,
				OutputTokens:   500,
				TotalTokens:    1500,
				ModelProvider:  "openai",
				ModelName:      "gpt-4",
				UnitPrice:      0.03,
				InputCost:      0.03,
				OutputCost:     0.03,
				TotalCost:      0.06,
				RequestType:    entity.RequestTypeChat,
				ResponseTimeMs: intPtr(150),
			}

			err := repo.Create(ctx, log)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

/**
 * BenchmarkRepository_GetUsageByDateRange 按日期范围查询TokenUsageLog性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次查询: < 10ms
 */
func BenchmarkRepository_GetUsageByDateRange(b *testing.B) {
	repo := NewMockTokenUsageLogRepository()
	ctx := context.Background()

	// 预先生成30天的数据
	now := time.Now()
	for day := 0; day < 30; day++ {
		for i := 0; i < 100; i++ {
			userID := uint64(1000 + i)
			botID := fmt.Sprintf("bot-%d", i%10)
			date := now.AddDate(0, 0, -day).Format("2006-01-02")
			log := &entity.TokenUsageLog{
				TenantID:       fmt.Sprintf("tenant-%d", i%10),
				UserID:         &userID,
				BotID:          &botID,
				InputTokens:    1000,
				OutputTokens:   500,
				TotalTokens:    1500,
				ModelProvider:  "openai",
				ModelName:      "gpt-4",
				TotalCost:      0.06,
				RequestType:    entity.RequestTypeChat,
				CreatedAt:      now.AddDate(0, 0, -day),
			}
			repo.Create(ctx, log)
		}
	}

	startDate := now.AddDate(0, 0, -29).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		_, err := repo.GetByDateRange(ctx, tenantID, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== 2. TokenUsageSummary Repository性能测试 ====================

/**
 * BenchmarkRepository_UpsertSummary 更新或插入TokenUsageSummary性能
 *
 * 性能目标：
 * - ops/s: > 1,000 次/秒
 * - 每次操作: < 1ms
 */
func BenchmarkRepository_UpsertSummary(b *testing.B) {
	repo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		botID := fmt.Sprintf("bot-%d", i%10)
		summary := &entity.TokenUsageSummary{
			TenantID:         fmt.Sprintf("tenant-%d", i%10),
			BotID:            &botID,
			SummaryDate:      time.Now().Format("2006-01-02"),
			TotalInputTokens: 100000,
			TotalOutputTokens: 50000,
			TotalTokens:      150000,
			TotalCost:        450.0,
			TotalRequests:    1000,
			CachedRequests:   100,
		}

		err := repo.Upsert(ctx, summary)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_GetSummaryByTenantAndDate 获取指定租户和日期的汇总性能
 *
 * 性能目标：
 * - ops/s: > 10,000 次/秒
 * - ns/op: < 100,000 ns/op
 */
func BenchmarkRepository_GetSummaryByTenantAndDate(b *testing.B) {
	repo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	// 预先生成汇总数据
	for i := 0; i < 100; i++ {
		botID := fmt.Sprintf("bot-%d", i%10)
		summary := &entity.TokenUsageSummary{
			TenantID:         fmt.Sprintf("tenant-%d", i%10),
			BotID:            &botID,
			SummaryDate:      time.Now().Format("2006-01-02"),
			TotalInputTokens: 100000,
			TotalOutputTokens: 50000,
			TotalTokens:      150000,
			TotalCost:        450.0,
			TotalRequests:    1000,
		}
		repo.Upsert(ctx, summary)
	}

	date := time.Now().Format("2006-01-02")
	hour := uint8(10)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%100)
		_, err := repo.GetByTenantAndDate(ctx, tenantID, date, &hour)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_GetSummariesByDateRange 按日期范围查询汇总性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次查询: < 10ms
 */
func BenchmarkRepository_GetSummariesByDateRange(b *testing.B) {
	repo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	// 预先生成30天的汇总数据
	now := time.Now()
	for day := 0; day < 30; day++ {
		for i := 0; i < 10; i++ {
			botID := fmt.Sprintf("bot-%d", i)
			date := now.AddDate(0, 0, -day).Format("2006-01-02")
			summary := &entity.TokenUsageSummary{
				TenantID:         fmt.Sprintf("tenant-%d", i),
				BotID:            &botID,
				SummaryDate:      date,
				TotalInputTokens: 100000,
				TotalOutputTokens: 50000,
				TotalTokens:      150000,
				TotalCost:        450.0,
				TotalRequests:    1000,
			}
			repo.Upsert(ctx, summary)
		}
	}

	startDate := now.AddDate(0, 0, -29).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		_, err := repo.GetByDateRange(ctx, tenantID, nil, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_UpsertSummary_Concurrent 并发更新汇总性能
 *
 * 性能目标：
 * - 并发数: 100
 * - 总ops/s: > 10,000
 * - 平均延迟: < 10ms
 */
func BenchmarkRepository_UpsertSummary_Concurrent(b *testing.B) {
	repo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			botID := fmt.Sprintf("bot-%d", i%10)
			summary := &entity.TokenUsageSummary{
				TenantID:         fmt.Sprintf("tenant-%d", i%10),
				BotID:            &botID,
				SummaryDate:      time.Now().Format("2006-01-02"),
				TotalInputTokens: 100000,
				TotalOutputTokens: 50000,
				TotalTokens:      150000,
				TotalCost:        450.0,
				TotalRequests:    1000,
			}

			err := repo.Upsert(ctx, summary)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

// ==================== 3. BudgetSettings Repository性能测试 ====================

/**
 * BenchmarkRepository_CreateBudget 创建预算配置性能
 *
 * 性能目标：
 * - ops/s: > 1,000 次/秒
 * - 每次操作: < 1ms
 */
func BenchmarkRepository_CreateBudget(b *testing.B) {
	repo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		budget := &entity.BudgetSettings{
			TenantID:              fmt.Sprintf("tenant-%d", i),
			BudgetType:            entity.BudgetTypeMonthly,
			BudgetAmount:          10000.0,
			Currency:              "CNY",
			AlertThreshold1:       80,
			AlertThreshold2:       95,
			HardCapEnabled:        true,
			NotificationChannels:  `["email"]`,
			NotificationRecipients: `["admin@example.com"]`,
		}

		err := repo.Create(ctx, budget)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_GetBudgetByTenantID 获取租户预算配置性能
 *
 * 性能目标：
 * - ops/s: > 10,000 次/秒
 * - ns/op: < 100,000 ns/op
 */
func BenchmarkRepository_GetBudgetByTenantID(b *testing.B) {
	repo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	// 预先创建预算配置
	for i := 0; i < 100; i++ {
		budget := &entity.BudgetSettings{
			TenantID:              fmt.Sprintf("tenant-%d", i),
			BudgetType:            entity.BudgetTypeMonthly,
			BudgetAmount:          10000.0,
			Currency:              "CNY",
			AlertThreshold1:       80,
			AlertThreshold2:       95,
			HardCapEnabled:        true,
			NotificationChannels:  `["email"]`,
			NotificationRecipients: `["admin@example.com"]`,
		}
		repo.Create(ctx, budget)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%100)
		_, err := repo.GetByTenantID(ctx, tenantID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_UpdateBudget 更新预算配置性能
 *
 * 性能目标：
 * - ops/s: > 1,000 次/秒
 * - 每次操作: < 1ms
 */
func BenchmarkRepository_UpdateBudget(b *testing.B) {
	repo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	// 预先创建预算配置
	for i := 0; i < 100; i++ {
		budget := &entity.BudgetSettings{
			TenantID:              fmt.Sprintf("tenant-%d", i),
			BudgetType:            entity.BudgetTypeMonthly,
			BudgetAmount:          10000.0,
			Currency:              "CNY",
			AlertThreshold1:       80,
			AlertThreshold2:       95,
			HardCapEnabled:        true,
			NotificationChannels:  `["email"]`,
			NotificationRecipients: `["admin@example.com"]`,
		}
		repo.Create(ctx, budget)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%100)
		budget, _ := repo.GetByTenantID(ctx, tenantID)
		if budget != nil {
			budget.BudgetAmount = 15000.0
			err := repo.Update(ctx, budget)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

// ==================== 4. BudgetAlert Repository性能测试 ====================

/**
 * BenchmarkRepository_CreateAlert 创建告警记录性能
 *
 * 性能目标：
 * - ops/s: > 500 次/秒
 * - 每次操作: < 2ms
 */
func BenchmarkRepository_CreateAlert(b *testing.B) {
	repo := NewMockBudgetAlertRepository()
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		message := "预算告警测试"
		alert := &entity.BudgetAlert{
			TenantID:              fmt.Sprintf("tenant-%d", i%10),
			AlertType:             entity.AlertTypeThreshold1,
			BudgetAmount:          10000.0,
			UsedAmount:            8200.0,
			UsagePercent:          82.0,
			AlertLevel:            entity.AlertLevelWarning,
			AlertMessage:          &message,
			NotificationChannels:  `["email"]`,
			NotificationStatus:    entity.NotificationStatusPending,
		}

		err := repo.Create(ctx, alert)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_GetAlertsByTenantID 获取租户告警历史性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次查询: < 10ms
 */
func BenchmarkRepository_GetAlertsByTenantID(b *testing.B) {
	repo := NewMockBudgetAlertRepository()
	ctx := context.Background()

	// 预先生成告警历史
	now := time.Now()
	for i := 0; i < 100; i++ {
		message := "预算告警测试"
		alert := &entity.BudgetAlert{
			TenantID:              fmt.Sprintf("tenant-%d", i%10),
			AlertType:             entity.AlertTypeThreshold1,
			BudgetAmount:          10000.0,
			UsedAmount:            8200.0,
			UsagePercent:          82.0,
			AlertLevel:            entity.AlertLevelWarning,
			AlertMessage:          &message,
			NotificationChannels:  `["email"]`,
			NotificationStatus:    entity.NotificationStatusSent,
			SentAt:                &now,
		}
		repo.Create(ctx, alert)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		_, err := repo.GetByTenantID(ctx, tenantID, 20)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_CheckAlertSentToday 检查告警发送状态性能
 *
 * 性能目标：
 * - ops/s: > 100,000 次/秒
 * - ns/op: < 10,000 ns/op
 */
func BenchmarkRepository_CheckAlertSentToday(b *testing.B) {
	repo := NewMockBudgetAlertRepository()
	ctx := context.Background()

	today := time.Now().Format("2006-01-02")

	// 预先标记一些告警为已发送
	for i := 0; i < 10; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i)
		alertType := entity.AlertTypeThreshold1
		key := fmt.Sprintf("%s:%s:%s", tenantID, alertType, today)
		repo.alertMap[key] = true
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
		_, err := repo.CheckAlertSentToday(ctx, tenantID, alertType, today)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== 5. 聚合查询性能测试 ====================

/**
 * BenchmarkRepository_GetTotalCost 获取总成本查询性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次查询: < 10ms
 */
func BenchmarkRepository_GetTotalCost(b *testing.B) {
	repo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	// 预先生成汇总数据
	now := time.Now()
	for day := 0; day < 30; day++ {
		for i := 0; i < 10; i++ {
			botID := fmt.Sprintf("bot-%d", i)
			date := now.AddDate(0, 0, -day).Format("2006-01-02")
			summary := &entity.TokenUsageSummary{
				TenantID:         fmt.Sprintf("tenant-%d", i),
				BotID:            &botID,
				SummaryDate:      date,
				TotalInputTokens: 100000,
				TotalOutputTokens: 50000,
				TotalTokens:      150000,
				TotalCost:        450.0,
				TotalRequests:    1000,
			}
			repo.Upsert(ctx, summary)
		}
	}

	startDate := now.AddDate(0, 0, -29)
	endDate := now

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)
		_, err := repo.GetTotalCost(ctx, tenantID, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkRepository_ComplexQuery 复杂聚合查询性能
 *
 * 性能目标：
 * - ops/s: > 50 次/秒
 * - 每次查询: < 20ms
 */
func BenchmarkRepository_ComplexQuery(b *testing.B) {
	summaryRepo := NewMockTokenUsageSummaryRepository()
	ctx := context.Background()

	// 预先生成大量汇总数据
	now := time.Now()
	for day := 0; day < 30; day++ {
		for i := 0; i < 10; i++ {
			for hour := 0; hour < 24; hour++ {
				botID := fmt.Sprintf("bot-%d", i)
				h := uint8(hour)
				date := now.AddDate(0, 0, -day).Format("2006-01-02")
				summary := &entity.TokenUsageSummary{
					TenantID:         fmt.Sprintf("tenant-%d", i),
					BotID:            &botID,
					SummaryDate:      date,
					SummaryHour:      &h,
					TotalInputTokens: 10000,
					TotalOutputTokens: 5000,
					TotalTokens:      15000,
					TotalCost:        45.0,
					TotalRequests:    100,
				}
				summaryRepo.Upsert(ctx, summary)
			}
		}
	}

	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10)

		// 复杂查询：获取最近7天的数据并聚合
		summaries, err := summaryRepo.GetByDateRange(ctx, tenantID, nil, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}

		// 聚合计算
		totalCost := 0.0
		totalTokens := int64(0)
		for _, summary := range summaries {
			totalCost += summary.TotalCost
			totalTokens += summary.TotalTokens
		}
		_ = totalCost
		_ = totalTokens
	}
}

// ==================== 6. 并发性能测试 ====================

/**
 * BenchmarkRepository_MixedOperations 混合操作并发性能
 *
 * 场景：模拟真实场景的混合操作（读写混合）
 *
 * 性能目标：
 * - 并发数: 100
 * - 混合比例: 70%读, 30%写
 * - 总ops/s: > 5,000
 */
func BenchmarkRepository_MixedOperations(b *testing.B) {
	usageLogRepo := NewMockTokenUsageLogRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	budgetRepo := NewMockBudgetSettingsRepository()
	ctx := context.Background()

	// 预先生成一些数据
	for i := 0; i < 100; i++ {
		botID := fmt.Sprintf("bot-%d", i%10)
		summary := &entity.TokenUsageSummary{
			TenantID:         fmt.Sprintf("tenant-%d", i%10),
			BotID:            &botID,
			SummaryDate:      time.Now().Format("2006-01-02"),
			TotalInputTokens: 100000,
			TotalOutputTokens: 50000,
			TotalTokens:      150000,
			TotalCost:        450.0,
			TotalRequests:    1000,
		}
		summaryRepo.Upsert(ctx, summary)

		budget := &entity.BudgetSettings{
			TenantID:     fmt.Sprintf("tenant-%d", i%10),
			BudgetType:   entity.BudgetTypeMonthly,
			BudgetAmount: 10000.0,
			Currency:     "CNY",
		}
		budgetRepo.Create(ctx, budget)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 70%读操作，30%写操作
			if i%10 < 7 {
				// 读操作
				tenantID := fmt.Sprintf("tenant-%d", i%10)
				_, _ = summaryRepo.GetByTenantAndDate(
					ctx,
					tenantID,
					time.Now().Format("2006-01-02"),
					nil,
				)
				_, _ = budgetRepo.GetByTenantID(ctx, tenantID)
			} else {
				// 写操作
				userID := uint64(1000 + i)
				botID := fmt.Sprintf("bot-%d", i%10)
				log := &entity.TokenUsageLog{
					TenantID:      fmt.Sprintf("tenant-%d", i%10),
					UserID:        &userID,
					BotID:         &botID,
					InputTokens:   1000,
					OutputTokens:  500,
					TotalTokens:   1500,
					ModelProvider: "openai",
					ModelName:     "gpt-4",
					TotalCost:     0.06,
					RequestType:   entity.RequestTypeChat,
				}
				_ = usageLogRepo.Create(ctx, log)
			}
			i++
		}
	})
}
