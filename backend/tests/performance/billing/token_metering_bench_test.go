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
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/service"
)

/**
 * Token计量性能基准测试
 *
 * 性能目标：
 * - 单次记录: < 10ms (P95)
 * - 批量记录(100条): < 100ms
 * - 批量记录(1000条): < 1秒
 * - 定价计算: < 0.01ms (10,000次 < 100ms)
 * - 汇总更新: < 50ms
 *
 * 设计原则：
 * - SOLID: 单一职责，每个基准测试只测试一个功能点
 * - DRY: 复用测试数据生成器
 * - KISS: 保持简单明了
 *
 * @author 研发B
 * @date 2025-12-30
 */

// ==================== Mock Repository ====================

/**
 * MockTokenUsageLogRepository Mock日志仓储
 */
type MockTokenUsageLogRepository struct {
	logs []*entity.TokenUsageLog
	mu   sync.RWMutex
}

func NewMockTokenUsageLogRepository() *MockTokenUsageLogRepository {
	return &MockTokenUsageLogRepository{
		logs: make([]*entity.TokenUsageLog, 0),
	}
}

func (r *MockTokenUsageLogRepository) Create(ctx context.Context, log *entity.TokenUsageLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	log.ID = uint64(len(r.logs) + 1)
	r.logs = append(r.logs, log)
	return nil
}

func (r *MockTokenUsageLogRepository) BatchCreate(ctx context.Context, logs []*entity.TokenUsageLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, log := range logs {
		log.ID = uint64(len(r.logs) + i + 1)
	}
	r.logs = append(r.logs, logs...)
	return nil
}

func (r *MockTokenUsageLogRepository) GetByDateRange(ctx context.Context, tenantID, startDate, endDate string) ([]*entity.TokenUsageLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.logs, nil
}

/**
 * MockTokenUsageSummaryRepository Mock汇总仓储
 */
type MockTokenUsageSummaryRepository struct {
	summaries map[string]*entity.TokenUsageSummary
	mu        sync.RWMutex
}

func NewMockTokenUsageSummaryRepository() *MockTokenUsageSummaryRepository {
	return &MockTokenUsageSummaryRepository{
		summaries: make(map[string]*entity.TokenUsageSummary),
	}
}

func (r *MockTokenUsageSummaryRepository) GetByTenantAndDate(
	ctx context.Context,
	tenantID, date string,
	hour *uint8,
) (*entity.TokenUsageSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key := fmt.Sprintf("%s:%s:%v", tenantID, date, hour)
	if summary, ok := r.summaries[key]; ok {
		return summary, nil
	}
	return nil, nil
}

func (r *MockTokenUsageSummaryRepository) Upsert(ctx context.Context, summary *entity.TokenUsageSummary) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s:%s:%v", summary.TenantID, summary.SummaryDate, summary.SummaryHour)
	r.summaries[key] = summary
	return nil
}

func (r *MockTokenUsageSummaryRepository) GetByDateRange(
	ctx context.Context,
	tenantID string,
	botID *string,
	startDate, endDate string,
) ([]*entity.TokenUsageSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	summaries := make([]*entity.TokenUsageSummary, 0, len(r.summaries))
	for _, summary := range r.summaries {
		if summary.TenantID == tenantID {
			summaries = append(summaries, summary)
		}
	}
	return summaries, nil
}

func (r *MockTokenUsageSummaryRepository) GetTotalCost(
	ctx context.Context,
	tenantID string,
	startDate, endDate time.Time,
) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	totalCost := 0.0
	for _, summary := range r.summaries {
		if summary.TenantID == tenantID {
			totalCost += summary.TotalCost
		}
	}
	return totalCost, nil
}

/**
 * MockBudgetSettingsRepository Mock预算仓储
 */
type MockBudgetSettingsRepository struct {
	budgets map[string]*entity.BudgetSettings
	mu      sync.RWMutex
}

func NewMockBudgetSettingsRepository() *MockBudgetSettingsRepository {
	return &MockBudgetSettingsRepository{
		budgets: make(map[string]*entity.BudgetSettings),
	}
}

func (r *MockBudgetSettingsRepository) GetByTenantID(ctx context.Context, tenantID string) (*entity.BudgetSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.budgets[tenantID], nil
}

func (r *MockBudgetSettingsRepository) Create(ctx context.Context, budget *entity.BudgetSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	budget.ID = uint64(len(r.budgets) + 1)
	r.budgets[budget.TenantID] = budget
	return nil
}

func (r *MockBudgetSettingsRepository) Update(ctx context.Context, budget *entity.BudgetSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.budgets[budget.TenantID] = budget
	return nil
}

// ==================== 测试数据生成器 ====================

/**
 * 生成测试Token使用记录请求
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
		IsCached:       idx%10 == 0, // 10%缓存命中率
		ResponseTimeMs: intPtr(100 + idx%200),
	}
}

/**
 * 生成批量测试Token使用记录请求
 */
func generateBatchTokenUsageRequests(count int) []*service.RecordTokenUsageRequest {
	requests := make([]*service.RecordTokenUsageRequest, count)
	for i := 0; i < count; i++ {
		requests[i] = generateTokenUsageRequest(i)
	}
	return requests
}

func intPtr(v int) *int {
	return &v
}

// ==================== 1. Token计量性能测试 ====================

/**
 * BenchmarkTokenMetering_RecordSingle 单次记录性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒 (每次操作 < 10ms)
 * - ns/op: < 10,000,000 ns/op
 * - 内存分配: 最小化
 */
func BenchmarkTokenMetering_RecordSingle(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil, // db (not used in benchmark)
		nil, // logger
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

/**
 * BenchmarkTokenMetering_RecordSingle_Parallel 并发单次记录性能
 *
 * 性能目标：
 * - 并发数: 100
 * - ops/s: > 1,000 次/秒
 * - 平均延迟: < 50ms
 */
func BenchmarkTokenMetering_RecordSingle_Parallel(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

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
}

/**
 * BenchmarkTokenMetering_BatchRecord_100 批量记录100条性能
 *
 * 性能目标：
 * - 批量大小: 100条
 * - 总耗时: < 100ms
 * - 平均每条: < 1ms
 */
func BenchmarkTokenMetering_BatchRecord_100(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()
	batchSize := 100
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
}

/**
 * BenchmarkTokenMetering_BatchRecord_1000 批量记录1000条性能
 *
 * 性能目标：
 * - 批量大小: 1000条
 * - 总耗时: < 1秒
 * - 平均每条: < 1ms
 */
func BenchmarkTokenMetering_BatchRecord_1000(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()
	batchSize := 1000
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
}

/**
 * BenchmarkTokenMetering_BatchRecord_Concurrent 批量记录并发性能
 *
 * 性能目标：
 * - 并发数: 10
 * - 每个并发批量大小: 100
 * - 总ops/s: > 1,000
 */
func BenchmarkTokenMetering_BatchRecord_Concurrent(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		batchSize := 100
		i := 0
		for pb.Next() {
			requests := generateBatchTokenUsageRequests(batchSize)
			for j := 0; j < batchSize; j++ {
				requests[j].TenantID = fmt.Sprintf("tenant-%d-%d", i, j)
			}

			_, err := meteringSvc.BatchRecordTokenUsage(ctx, requests)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

// ==================== 2. 定价引擎性能测试 ====================

/**
 * BenchmarkPricingEngine_CalculateCost 定价计算性能
 *
 * 性能目标：
 * - ops/s: > 100,000 次/秒
 * - ns/op: < 10,000 ns/op (0.01ms)
 * - 零内存分配（栈上计算）
 */
func BenchmarkPricingEngine_CalculateCost(b *testing.B) {
	pricingEngine := service.NewPricingEngine()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		provider := []string{"openai", "anthropic", "qwen", "baidu", "zhipu"}[i%5]
		model := []string{"gpt-4", "claude-3-opus", "qwen-max", "ernie-bot-4", "chatglm-pro"}[i%5]
		inputTokens := 1000 + i%100
		outputTokens := 500 + i%50

		_ = pricingEngine.CalculateCost(provider, model, inputTokens, outputTokens)
	}
}

/**
 * BenchmarkPricingEngine_CalculateCost_10000 10000次计算性能
 *
 * 性能目标：
 * - 总耗时: < 100ms
 * - 平均每次: < 0.01ms
 */
func BenchmarkPricingEngine_CalculateCost_10000(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	iterations := 10000

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < iterations; j++ {
			provider := []string{"openai", "anthropic", "qwen"}[j%3]
			model := []string{"gpt-4", "claude-3-opus", "qwen-max"}[j%3]
			_ = pricingEngine.CalculateCost(provider, model, 1000, 500)
		}
	}
}

/**
 * BenchmarkPricingEngine_EstimateCost 成本估算性能
 *
 * 性能目标：
 * - ops/s: > 50,000 次/秒
 * - ns/op: < 20,000 ns/op
 */
func BenchmarkPricingEngine_EstimateCost(b *testing.B) {
	pricingEngine := service.NewPricingEngine()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		provider := "openai"
		model := "gpt-4"
		dailyTokens := 100000
		days := 30

		_ = pricingEngine.EstimateCost(provider, model, dailyTokens, days)
	}
}

/**
 * BenchmarkPricingEngine_CalculateSavings 节省计算性能
 *
 * 性能目标：
 * - ops/s: > 50,000 次/秒
 * - ns/op: < 20,000 ns/op
 */
func BenchmarkPricingEngine_CalculateSavings(b *testing.B) {
	pricingEngine := service.NewPricingEngine()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = pricingEngine.CalculateSavings(
			"openai", "gpt-4",
			"anthropic", "claude-3-haiku",
			1000, 500,
		)
	}
}

// ==================== 3. Token使用统计查询性能测试 ====================

/**
 * BenchmarkTokenMetering_GetUsageStats 使用统计查询性能
 *
 * 性能目标：
 * - ops/s: > 100 次/秒
 * - 每次查询: < 100ms
 */
func BenchmarkTokenMetering_GetUsageStats(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 预先生成一些汇总数据
	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
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
}

/**
 * BenchmarkTokenMetering_GetDailyUsageStats 每日趋势查询性能
 *
 * 性能目标：
 * - 查询天数: 30天
 * - 总耗时: < 200ms
 * - 平均每天: < 7ms
 */
func BenchmarkTokenMetering_GetDailyUsageStats(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

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
			TotalCost:        45.0,
			TotalRequests:    1000,
			CachedRequests:   100,
		}
		usageSummaryRepo.Upsert(ctx, summary)
	}

	startDate := now.AddDate(0, 0, -29).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10+1)
		_, err := meteringSvc.GetDailyUsageStats(ctx, tenantID, startDate, endDate, nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkTokenMetering_GetModelUsageStats 模型统计查询性能
 *
 * 性能目标：
 * - ops/s: > 50 次/秒
 * - 每次查询: < 150ms
 */
func BenchmarkTokenMetering_GetModelUsageStats(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	// 预先生成日志数据
	now := time.Now()
	for i := 0; i < 1000; i++ {
		userID := uint64(1000 + i)
		botID := fmt.Sprintf("bot-%d", i%10)
		log := &entity.TokenUsageLog{
			TenantID:       "tenant-1",
			UserID:         &userID,
			BotID:          &botID,
			ModelProvider:  []string{"openai", "anthropic", "qwen"}[i%3],
			ModelName:      []string{"gpt-4", "claude-3-opus", "qwen-max"}[i%3],
			InputTokens:    1000 + i%100,
			OutputTokens:   500 + i%50,
			TotalTokens:    1500 + i%150,
			TotalCost:      0.45,
			RequestType:    entity.RequestTypeChat,
			CreatedAt:      now,
		}
		usageLogRepo.Create(ctx, log)
	}

	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")
	endDate := now.Format("2006-01-02")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%d", i%10+1)
		_, err := meteringSvc.GetModelUsageStats(ctx, tenantID, startDate, endDate)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ==================== 4. Token使用汇总更新性能测试 ====================

/**
 * BenchmarkTokenMetering_UpdateSummary 汇总更新性能
 *
 * 性能目标：
 * - ops/s: > 20 次/秒
 * - 每次更新: < 50ms
 */
func BenchmarkTokenMetering_UpdateSummary(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		userID := uint64(1000 + i%100)
		botID := fmt.Sprintf("bot-%d", i%10)
		log := &entity.TokenUsageLog{
			TenantID:       fmt.Sprintf("tenant-%d", i%10),
			UserID:         &userID,
			BotID:          &botID,
			InputTokens:    1000,
			OutputTokens:   500,
			TotalTokens:    1500,
			TotalCost:      0.45,
			RequestType:    entity.RequestTypeChat,
			ResponseTimeMs: intPtr(150),
		}

		// 直接调用私有方法（通过反射或导出测试）
		// 这里我们简化：只测试Upsert性能
		summary := &entity.TokenUsageSummary{
			TenantID:         log.TenantID,
			BotID:            log.BotID,
			SummaryDate:      time.Now().Format("2006-01-02"),
			TotalInputTokens: int64(log.InputTokens),
			TotalOutputTokens: int64(log.OutputTokens),
			TotalTokens:      int64(log.TotalTokens),
			TotalCost:        log.TotalCost,
			TotalRequests:    1,
		}

		err := usageSummaryRepo.Upsert(ctx, summary)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkTokenMetering_UpdateSummary_Concurrent 并发汇总更新性能
 *
 * 性能目标：
 * - 并发数: 50
 * - 总ops/s: > 1,000
 * - 平均延迟: < 50ms
 */
func BenchmarkTokenMetering_UpdateSummary_Concurrent(b *testing.B) {
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
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
				TotalInputTokens: 1000,
				TotalOutputTokens: 500,
				TotalTokens:      1500,
				TotalCost:        0.45,
				TotalRequests:    1,
			}

			err := usageSummaryRepo.Upsert(ctx, summary)
			if err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

// ==================== 5. 综合场景性能测试 ====================

/**
 * BenchmarkTokenMetering_CompleteFlow 完整流程性能（记录+统计+查询）
 *
 * 性能目标：
 * - 完整流程: < 200ms
 * - 包含：记录Token → 更新汇总 → 查询统计
 */
func BenchmarkTokenMetering_CompleteFlow(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// 1. 记录Token使用
		req := generateTokenUsageRequest(i)
		_, err := meteringSvc.RecordTokenUsage(ctx, req)
		if err != nil {
			b.Fatal(err)
		}

		// 2. 查询使用统计（简化，只查询当天的）
		now := time.Now()
		today := now.Format("2006-01-02")
		filter := &service.UsageStatsFilter{
			StartDate: today,
			EndDate:   today,
		}

		_, err = meteringSvc.GetUsageStats(ctx, req.TenantID, filter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

/**
 * BenchmarkTokenMetering_HighFrequency 高频记录场景
 *
 * 场景：模拟高并发场景下的Token记录
 *
 * 性能目标：
 * - 100并发
 * - 持续记录1000次
 * - 无错误
 * - 平均延迟 < 50ms
 */
func BenchmarkTokenMetering_HighFrequency(b *testing.B) {
	pricingEngine := service.NewPricingEngine()
	usageLogRepo := NewMockTokenUsageLogRepository()
	usageSummaryRepo := NewMockTokenUsageSummaryRepository()
	budgetAlertSvc := &service.BudgetAlertService{} // Mock

	meteringSvc := service.NewTokenMeteringService(
		pricingEngine,
		usageLogRepo,
		usageSummaryRepo,
		budgetAlertSvc,
		nil,
		nil,
	)

	ctx := context.Background()

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
}
