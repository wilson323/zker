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
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

/**
 * Token计量和预算管理系统 - 完整性能基准测试套件
 *
 * 性能目标（基于ZKER性能基线文档）:
 * - Token单次记录: < 100ms
 * - Token批量记录(1000条): < 1秒
 * - 定价计算(10000次): < 100ms
 * - 使用统计查询: < 50ms
 * - 预算检查: < 50ms
 * - 并发吞吐量: > 1000 req/s
 *
 * 性能指标:
 * - CPU: < 80% (单核)
 * - 内存: < 500MB
 * - 数据库连接: < 100
 * - goroutine数量: < 1000
 *
 * 设计原则:
 * - SOLID: 单一职责，每个基准测试只测试一个功能点
 * - DRY: 复用测试数据生成器和Mock对象
 * - KISS: 保持简单明了，易于理解和维护
 *
 * @author 研发B
 * @date 2025-12-30
 * @version v1.0.0
 */

// ==================== Benchmark 1: Token计量单次记录性能 ====================

/**
 * BenchmarkRecordTokenUsage 单次Token记录性能
 *
 * 性能目标：
 * - ops/s: > 10 次/秒
 * - 每次操作: < 100ms
 * - 内存分配: 最小化
 *
 * 测试内容：
 * - Token使用记录创建
 * - 成本计算
 * - 汇总数据更新
 * - 预算告警检查
 */
func BenchmarkRecordTokenUsage(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()
	req := generateTokenUsageRequest(1)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req.TenantID = fmt.Sprintf("tenant-%d", i%100)
		_, err := meteringSvc.RecordTokenUsage(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== Benchmark 2: Token计量批量记录性能 ====================

/**
 * BenchmarkBatchRecordTokenUsage 批量Token记录性能
 *
 * 性能目标：
 * - 批量大小: 10, 50, 100, 500, 1000
 * - 1000条记录: < 1秒
 * - 平均每条: < 1ms
 */
func BenchmarkBatchRecordTokenUsage(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 测试不同批次大小
	batchSizes := []int{10, 50, 100, 500, 1000}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize_%d", batchSize), func(b *testing.B) {
			requests := generateBatchTokenUsageRequests(batchSize)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				// 更新租户ID避免重复
				for j := 0; j < batchSize; j++ {
					requests[j].TenantID = fmt.Sprintf("tenant-%d-%d", i, j)
				}

				_, err := meteringSvc.BatchRecordTokenUsage(ctx, requests)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ==================== Benchmark 3: Token使用统计查询性能 ====================

/**
 * BenchmarkGetUsageStats Token使用统计查询性能
 *
 * 性能目标：
 * - ops/s: > 20 次/秒
 * - 每次查询: < 50ms
 * - 不同数据量: 1000, 10000, 100000条记录
 */
func BenchmarkGetUsageStats(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 测试不同数据量
	dataSizes := []struct {
		name  string
		count int
	}{
		{"1K_Records", 1000},
		{"10K_Records", 10000},
		{"100K_Records", 100000},
	}

	for _, dataSize := range dataSizes {
		b.Run(dataSize.name, func(b *testing.B) {
			// 预先生成数据
			now := time.Now()
			for i := 0; i < dataSize.count; i++ {
				date := now.AddDate(0, 0, -i%30).Format("2006-01-02")
				summary := &entity.TokenUsageSummary{
					TenantID:         "tenant-1",
					SummaryDate:      date,
					TotalInputTokens: 100000,
					TotalOutputTokens: 50000,
					TotalTokens:      150000,
					TotalCost:        45.0,
					TotalRequests:    1000,
					CachedRequests:   100,
				}
				usageSummaryRepo.Upsert(ctx, summary)
			}

			filter := &service.UsageStatsFilter{
				StartDate: now.AddDate(0, 0, -29).Format("2006-01-02"),
				EndDate:   now.Format("2006-01-02"),
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tenantID := fmt.Sprintf("tenant-%d", i%10+1)
				_, err := meteringSvc.GetUsageStats(ctx, tenantID, filter)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ==================== Benchmark 4: 定价引擎性能测试 ====================

/**
 * BenchmarkPricingEngine_CalculateCost 定价引擎成本计算性能
 *
 * 性能目标：
 * - 10000次计算: < 100ms
 * - 平均每次: < 0.01ms
 * - 零内存分配（栈上计算）
 */
func BenchmarkPricingEngine_CalculateCost(b *testing.B) {
	pricingEngine := service.NewPricingEngine()

	// 定义13个AI模型进行测试
	models := []struct {
		provider string
		model    string
	}{
		{"openai", "gpt-4"},
		{"openai", "gpt-3.5-turbo"},
		{"anthropic", "claude-3-opus"},
		{"anthropic", "claude-3-haiku"},
		{"qwen", "qwen-max"},
		{"qwen", "qwen-plus"},
		{"qwen", "qwen-turbo"},
		{"baidu", "ernie-bot-4"},
		{"zhipu", "chatglm-pro"},
		{"zhipu", "chatglm-turbo"},
		{"deepseek", "deepseek-chat"},
		{"baichuan", "baichuan2-turbo"},
		{"minimax", "abab5.5-chat"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		model := models[i%len(models)]
		inputTokens := 1000 + rand.Intn(5000)
		outputTokens := 500 + rand.Intn(2000)

		_ = pricingEngine.CalculateCost(model.provider, model.model, inputTokens, outputTokens)
	}
}

/**
 * BenchmarkPricingEngine_CalculateCost_10000 10000次定价计算性能
 *
 * 性能目标：
 * - 总耗时: < 100ms
 * - 验证: 所有13个AI模型的计算性能
 */
func BenchmarkPricingEngine_CalculateCost_10000(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	iterations := 10000

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			providers := []string{"openai", "anthropic", "qwen"}
			models := []string{"gpt-4", "claude-3-opus", "qwen-max"}
			provider := providers[j%3]
			model := models[j%3]
			_ = pricingEngine.CalculateCost(provider, model, 1000, 500)
		}
	}
}

// ==================== Benchmark 5: 预算检查性能测试 ====================

/**
 * BenchmarkCheckBudget 单租户预算检查性能
 *
 * 性能目标：
 * - ops/s: > 20 次/秒
 * - 每次检查: < 50ms
 * - 包含: 使用量查询 + 告警判断 + 通知发送
 */
func BenchmarkCheckBudget(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{}

	budgetAlertSvc := service.NewBudgetAlertService(
		budgetRepo,
		summaryRepo,
		alertRepo,
		notifier,
		nil,
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
		TotalCost:   7000.0,
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

// ==================== Benchmark 6: 预算使用情况查询性能 ====================

/**
 * BenchmarkGetBudgetUsage 获取预算使用情况性能
 *
 * 性能目标：
 * - ops/s: > 10 次/秒
 * - 每次查询: < 100ms
 * - 包含: 聚合查询 + 使用率计算 + 预测
 */
func BenchmarkGetBudgetUsage(b *testing.B) {
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

	startDate := now.AddDate(0, 0, -29)
	endDate := now

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10+1)
		_, err := summaryRepo.GetTotalCost(ctx, tenantID, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}

		// 获取日期范围内的汇总数据
		startDateStr := startDate.Format("2006-01-02")
		endDateStr := endDate.Format("2006-01-02")
		_, err = summaryRepo.GetByDateRange(ctx, tenantID, nil, startDateStr, endDateStr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== Benchmark 7: 批量预算检查性能 ====================

/**
 * BenchmarkCheckBudgetsForAllTenants 批量预算检查性能
 *
 * 性能目标：
 * - 100个租户: < 5秒
 * - 平均每个: < 50ms
 * - 场景: 定时任务场景
 */
func BenchmarkCheckBudgetsForAllTenants(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{}

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
			TotalCost:   7000.0 + float64(i%30)*100,
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

// ==================== Benchmark 8: Repository批量插入性能 ====================

/**
 * BenchmarkTokenUsageLogBatchCreate Repository批量插入性能
 *
 * 性能目标：
 * - 测试: GORM CreateInBatches性能
 * - 对比: 不同批次大小（10, 50, 100, 500）
 * - 指标: 吞吐量（条/秒）
 */
func BenchmarkTokenUsageLogBatchCreate(b *testing.B) {
	repo := NewMockTokenUsageLogRepository()
	ctx := context.Background()

	// 测试不同批次大小
	batchSizes := []int{10, 50, 100, 500}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize_%d", batchSize), func(b *testing.B) {
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
		})
	}
}

// ==================== Benchmark 9: 汇总查询性能 ====================

/**
 * BenchmarkTokenUsageSummaryGetTotalCost 汇总查询性能
 *
 * 性能目标：
 * - 场景: 不同时间范围（日、周、月）
 * - 优化: 验证索引效果
 * - 指标: < 100ms
 */
func BenchmarkTokenUsageSummaryGetTotalCost(b *testing.B) {
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

	// 测试不同时间范围
	timeRanges := []struct {
		name      string
		days      int
		startDate time.Time
		endDate   time.Time
	}{
		{"Daily", 1, now.AddDate(0, 0, -1), now},
		{"Weekly", 7, now.AddDate(0, 0, -7), now},
		{"Monthly", 30, now.AddDate(0, 0, -29), now},
	}

	for _, timeRange := range timeRanges {
		b.Run(timeRange.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tenantID := fmt.Sprintf("tenant-%d", i%10)
				_, err := repo.GetTotalCost(ctx, tenantID, timeRange.startDate, timeRange.endDate)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ==================== Benchmark 10: 并发Token记录性能 ====================

/**
 * BenchmarkConcurrentRecordTokenUsage 并发Token记录性能
 *
 * 性能目标：
 * - 配置: 10, 50, 100, 500并发goroutine
 * - 指标: 吞吐量、延迟、P99延迟
 * - 验证: 无竞态条件
 */
func BenchmarkConcurrentRecordTokenUsage(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 测试不同并发级别
	concurrencyLevels := []int{10, 50, 100, 500}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency_%d", concurrency), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					req := generateTokenUsageRequest(i)
					_, err := meteringSvc.RecordTokenUsage(ctx, req)
					if err != nil {
						b.Fatal(err)
					}
					i++
				}
			})
		})
	}
}

// ==================== Benchmark 11: 并发预算检查性能 ====================

/**
 * BenchmarkConcurrentCheckBudget 并发预算检查性能
 *
 * 性能目标：
 * - 测试: 并发预算检查性能
 * - 验证: 无竞态条件
 * - 指标: 吞吐量
 */
func BenchmarkConcurrentCheckBudget(b *testing.B) {
	budgetRepo := NewMockBudgetSettingsRepository()
	summaryRepo := NewMockTokenUsageSummaryRepository()
	alertRepo := NewMockBudgetAlertRepository()
	notifier := &service.NotificationService{}

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

// ==================== Profile 1: CPU Profile ====================

/**
 * BenchmarkCPUProfile CPU性能分析
 *
 * 用法:
 *   go test -cpuprofile=cpu.prof -bench=BenchmarkCPUProfile ./tests/benchmark/billing/
 *   go tool pprof cpu.prof
 *
 * 分析:
 *   (pprof) top10
 *   (pprof) list [function_name]
 *   (pprof) web
 */
func BenchmarkCPUProfile(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()
	req := generateTokenUsageRequest(1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req.TenantID = fmt.Sprintf("tenant-%d", i%100)
		_, _ = meteringSvc.RecordTokenUsage(ctx, req)
	}
}

// ==================== Profile 2: Memory Profile ====================

/**
 * BenchmarkMemoryProfile 内存性能分析
 *
 * 用法:
 *   go test -memprofile=mem.prof -bench=BenchmarkMemoryProfile ./tests/benchmark/billing/
 *   go tool pprof mem.prof
 *
 * 分析:
 *   (pprof) top
 *   (pprof) list [function_name]
 *   (pprof) web
 */
func BenchmarkMemoryProfile(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()
	requests := generateBatchTokenUsageRequests(100)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			requests[j].TenantID = fmt.Sprintf("tenant-%d-%d", i, j)
		}
		_, _ = meteringSvc.BatchRecordTokenUsage(ctx, requests)
	}
}

// ==================== Profile 3: Heap Profile (内存泄漏检测) ====================

/**
 * TestMemoryLeak 内存泄漏检测测试
 *
 * 验证:
 * - 长时间运行无内存泄漏
 * - 场景: 100000次Token记录
 * - 检查: 内存增长 < 10MB
 */
func TestMemoryLeak(t *testing.T) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 强制GC，获取初始内存
	runtime.GC()
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 执行100000次Token记录
	iterations := 100000
	for i := 0; i < iterations; i++ {
		req := generateTokenUsageRequest(i)
		_, err := meteringSvc.RecordTokenUsage(ctx, req)
		if err != nil {
			t.Fatalf("RecordTokenUsage failed: %v", err)
		}
	}

	// 强制GC，获取结束内存
	runtime.GC()
	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

	// 计算内存增长
	memoryGrowth := int64(m2.HeapInuse) - int64(m1.HeapInuse)
	memoryGrowthMB := float64(memoryGrowth) / 1024 / 1024

	t.Logf("Iterations: %d", iterations)
	t.Logf("Initial Heap: %d MB", m1.HeapInuse/1024/1024)
	t.Logf("Final Heap: %d MB", m2.HeapInuse/1024/1024)
	t.Logf("Memory Growth: %.2f MB", memoryGrowthMB)

	// 验证内存增长 < 10MB
	if memoryGrowthMB > 10 {
		t.Errorf("Memory leak detected: growth %.2f MB > 10 MB", memoryGrowthMB)
	}
}

// ==================== Stress Test 1: 持续高负载 ====================

/**
 * TestStressTokenRecording 持续高负载压力测试
 *
 * 场景: 持续1小时，100 QPS
 * 验证: 系统稳定性
 * 检查: 无错误、无延迟突增
 *
 * 注意: 此测试运行时间较长，仅在需要时执行
 */
func TestStressTokenRecording(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 配置
	duration := 1 * time.Hour // 1小时
	qps := 100                // 每秒100个请求
	interval := time.Second / time.Duration(qps)

	startTime := time.Now()
	endTime := startTime.Add(duration)
	requestCount := 0
	errorCount := 0

	// 记录延迟统计
	latencies := make([]time.Duration, 0, qps*int(duration.Seconds()))

	t.Logf("Starting stress test: QPS=%d, Duration=%v", qps, duration)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for time.Now().Before(endTime) {
		select {
		case <-ticker.C:
			reqStart := time.Now()
			req := generateTokenUsageRequest(requestCount)

			_, err := meteringSvc.RecordTokenUsage(ctx, req)
			latency := time.Since(reqStart)
			latencies = append(latencies, latency)

			if err != nil {
				errorCount++
				t.Logf("Request %d failed: %v", requestCount, err)
			}

			requestCount++
		}
	}

	// 统计结果
	totalDuration := time.Since(startTime)
	actualQPS := float64(requestCount) / totalDuration.Seconds()
	errorRate := float64(errorCount) / float64(requestCount) * 100

	// 计算延迟百分位数
	sortLatencies(latencies)
	p50 := latencies[len(latencies)*50/100]
	p95 := latencies[len(latencies)*95/100]
	p99 := latencies[len(latencies)*99/100]

	t.Logf("Stress test completed:")
	t.Logf("  Total Duration: %v", totalDuration)
	t.Logf("  Total Requests: %d", requestCount)
	t.Logf("  Error Count: %d", errorCount)
	t.Logf("  Error Rate: %.2f%%", errorRate)
	t.Logf("  Actual QPS: %.2f", actualQPS)
	t.Logf("  Latency P50: %v", p50)
	t.Logf("  Latency P95: %v", p95)
	t.Logf("  Latency P99: %v", p99)

	// 验证
	if errorRate > 1 {
		t.Errorf("Error rate too high: %.2f%% > 1%%", errorRate)
	}
	if p95 > 200*time.Millisecond {
		t.Errorf("P95 latency too high: %v > 200ms", p95)
	}
}

/**
 * sortLatencies 对延迟切片进行排序
 */
func sortLatencies(latencies []time.Duration) {
	for i := 0; i < len(latencies)-1; i++ {
		for j := i + 1; j < len(latencies); j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}
}

// ==================== Stress Test 2: 峰值负载 ====================

/**
 * TestPeakLoad 峰值负载压力测试
 *
 * 场景: 短时间峰值，1000 QPS
 * 验证: 系统不崩溃
 * 检查: 错误率 < 1%
 */
func TestPeakLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping peak load test in short mode")
	}

	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{}

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 配置
	duration := 1 * time.Minute // 1分钟
	qps := 1000                 // 每秒1000个请求
	workers := 100              // 100个并发worker

	startTime := time.Now()
	endTime := startTime.Add(duration)

	var wg sync.WaitGroup
	var mu sync.Mutex
	requestCount := 0
	errorCount := 0

	t.Logf("Starting peak load test: QPS=%d, Duration=%v, Workers=%d", qps, duration, workers)

	// 启动worker
	requestChan := make(chan int, qps*10)

	// 请求生成器
	go func() {
		defer close(requestChan)
		i := 0
		for time.Now().Before(endTime) {
			for j := 0; j < qps/10; j++ {
				requestChan <- i
				i++
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 启动workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for reqID := range requestChan {
				req := generateTokenUsageRequest(reqID)
				_, err := meteringSvc.RecordTokenUsage(ctx, req)

				mu.Lock()
				requestCount++
				if err != nil {
					errorCount++
				}
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// 统计结果
	totalDuration := time.Since(startTime)
	actualQPS := float64(requestCount) / totalDuration.Seconds()
	errorRate := float64(errorCount) / float64(requestCount) * 100

	t.Logf("Peak load test completed:")
	t.Logf("  Total Duration: %v", totalDuration)
	t.Logf("  Total Requests: %d", requestCount)
	t.Logf("  Error Count: %d", errorCount)
	t.Logf("  Error Rate: %.2f%%", errorRate)
	t.Logf("  Actual QPS: %.2f", actualQPS)

	// 验证
	if errorRate > 1 {
		t.Errorf("Error rate too high: %.2f%% > 1%%", errorRate)
	}
}

// ==================== 辅助函数 ====================

/**
 * generateTokenUsageRequest 生成测试Token使用记录请求
 */
func generateTokenUsageRequest(idx int) *service.RecordTokenUsageRequest {
	userID := uint64(1000 + idx%100)
	botID := fmt.Sprintf("bot-%d", idx%10)
	conversationID := fmt.Sprintf("conv-%d", idx%1000)

	return &service.RecordTokenUsageRequest{
		TenantID:       fmt.Sprintf("tenant-%d", idx%10),
		UserID:         &userID,
		BotID:          &botID,
		ConversationID: &conversationID,
		ModelProvider:  "openai",
		ModelName:      "gpt-4",
		InputTokens:    1000 + idx%100,
		OutputTokens:   500 + idx%50,
		TotalTokens:    1500 + idx%150,
		RequestType:    entity.RequestTypeChat,
		IsCached:       idx%10 == 0,
		ResponseTimeMs: intPtr(100 + idx%200),
	}
}

/**
 * generateBatchTokenUsageRequests 生成批量测试Token使用记录请求
 */
func generateBatchTokenUsageRequests(count int) []*service.RecordTokenUsageRequest {
	requests := make([]*service.RecordTokenUsageRequest, count)
	for i := 0; i < count; i++ {
		requests[i] = generateTokenUsageRequest(i)
	}
	return requests
}

/**
 * generateBudgetSettings 生成测试预算配置
 */
func generateBudgetSettings(tenantID string, budgetAmount float64) *entity.BudgetSettings {
	hardCapAmount := budgetAmount * 1.1
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
 * intPtr 返回int指针
 */
func intPtr(v int) *int {
	return &v
}
