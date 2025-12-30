// backend/domain/billing/service/token_metering_async.go
package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// AsyncTokenRecorder 异步Token记录器
// 使用channel实现异步写入，避免阻塞API请求
type AsyncTokenRecorder struct {
	tokenUsageRepo repository.TokenUsageRepository
	usageChan      chan *entity.TokenUsage
	wg             sync.WaitGroup
	stopChan       chan struct{}
	bufferSize     int
}

// NewAsyncTokenRecorder 创建异步Token记录器
func NewAsyncTokenRecorder(
	tokenUsageRepo repository.TokenUsageRepository,
	bufferSize int,
) *AsyncTokenRecorder {
	if bufferSize <= 0 {
		bufferSize = 1000 // 默认缓冲区大小
	}

	return &AsyncTokenRecorder{
		tokenUsageRepo: tokenUsageRepo,
		usageChan:      make(chan *entity.TokenUsage, bufferSize),
		stopChan:       make(chan struct{}),
		bufferSize:     bufferSize,
	}
}

// Start 启动异步处理goroutine
func (r *AsyncTokenRecorder) Start(ctx context.Context) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		logs.CtxInfof(ctx, "[AsyncTokenRecorder] started")

		for {
			select {
			case usage := <-r.usageChan:
				// 异步写入数据库
				if err := r.tokenUsageRepo.Create(ctx, usage); err != nil {
					logs.CtxErrorf(ctx,
						"[AsyncTokenRecorder] failed to record token usage: tenant=%s, model=%s, tokens=%d, error=%v",
						usage.TenantID, usage.ModelID, usage.TotalTokens, err)
				} else {
					logs.CtxDebugf(ctx,
						"[AsyncTokenRecorder] token usage recorded: tenant=%s, model=%s, tokens=%d",
						usage.TenantID, usage.ModelID, usage.TotalTokens)
				}

			case <-r.stopChan:
				logs.CtxInfof(ctx, "[AsyncTokenRecorder] stopping...")
				// 处理剩余数据
				r.flushRemaining(ctx)
				return
			}
		}
	}()
}

// Stop 停止异步处理器
func (r *AsyncTokenRecorder) Stop() {
	close(r.stopChan)
	r.wg.Wait()
	// 注意：这里没有context，使用不带ctx的日志
	// logs.Infof("[AsyncTokenRecorder] stopped")
}

// RecordAsync 异步记录Token使用
func (r *AsyncTokenRecorder) RecordAsync(usage *entity.TokenUsage) error {
	select {
	case r.usageChan <- usage:
		return nil
	default:
		// channel满了，记录日志但不阻塞请求
		// 注意：这里没有context，无法使用Ctx日志
		// logs.Warnf("[AsyncTokenRecorder] token usage channel full (size=%d), dropping record",
		// 	len(r.usageChan))
		return fmt.Errorf("token usage channel full, record dropped")
	}
}

// GetStats 获取异步记录器统计信息
func (r *AsyncTokenRecorder) GetStats() AsyncRecorderStats {
	return AsyncRecorderStats{
		BufferSize:    r.bufferSize,
		CurrentUsage:  len(r.usageChan),
		UsagePercent:  float64(len(r.usageChan)) / float64(r.bufferSize) * 100,
	}
}

// flushRemaining 处理channel中剩余的数据
func (r *AsyncTokenRecorder) flushRemaining(ctx context.Context) {
	logs.CtxInfof(ctx, "[AsyncTokenRecorder] flushing %d remaining records...", len(r.usageChan))

	flushed := 0
	for {
		select {
		case usage := <-r.usageChan:
			if err := r.tokenUsageRepo.Create(ctx, usage); err != nil {
				logs.CtxErrorf(ctx, "[AsyncTokenRecorder] failed to flush record: %v", err)
			}
			flushed++
		default:
			// channel已空
			logs.CtxInfof(ctx, "[AsyncTokenRecorder] flushed %d records", flushed)
			return
		}
	}
}

// AsyncRecorderStats 异步记录器统计信息
type AsyncRecorderStats struct {
	BufferSize   int     // 缓冲区大小
	CurrentUsage int     // 当前使用量
	UsagePercent float64 // 使用率（百分比）
}

// =====================================================================
// 使用示例
// =====================================================================

// TokenMeteringServiceWithAsync 支持异步的Token计量服务
type TokenMeteringServiceWithAsync struct {
	asyncRecorder *AsyncTokenRecorder
	budgetRepo    repository.TokenBudgetRepository
}

// NewTokenMeteringServiceWithAsync 创建支持异步的Token计量服务
func NewTokenMeteringServiceWithAsync(
	asyncRecorder *AsyncTokenRecorder,
	budgetRepo repository.TokenBudgetRepository,
) *TokenMeteringServiceWithAsync {
	return &TokenMeteringServiceWithAsync{
		asyncRecorder: asyncRecorder,
		budgetRepo:    budgetRepo,
	}
}

// 注意：TokenMeteringServiceWithAsync 保留作为参考，实际使用请参考 token_metering_service.go 中的完整实现

// RecordTokenUsageAsync 异步记录Token使用（推荐方法）
// 注意：此方法需要配合PricingEngine使用，这里简化实现
func (s *TokenMeteringServiceWithAsync) RecordTokenUsageAsync(
	ctx context.Context,
	req *RecordTokenUsageRequest,
) error {
	// 计算成本 - 简化版本，实际应该使用 PricingEngine
	cost := s.calculateCost(req.ModelName, req.TotalTokens)

	// 创建使用记录 - 注意：entity.TokenUsage 是旧版本，应该使用 TokenUsageLog
	userID := ""
	if req.UserID != nil {
		userID = fmt.Sprintf("%d", *req.UserID)
	}

	usage := &entity.TokenUsage{
		TenantID:        req.TenantID,
		UserID:          userID,
		BotID:           stringValue(req.BotID),
		ModelID:         req.ModelName,
		PromptTokens:    req.InputTokens,
		CompletionTokens: req.OutputTokens,
		TotalTokens:     req.TotalTokens,
		CostUSD:         cost,
	}

	// 异步写入（不阻塞API请求）
	return s.asyncRecorder.RecordAsync(usage)
}

// calculateCost 计算Token成本（简化版本）
func (s *TokenMeteringServiceWithAsync) calculateCost(modelID string, totalTokens int) float64 {
	pricing := map[string]float64{
		"gpt-4":          0.03,
		"gpt-4-32k":      0.06,
		"gpt-3.5-turbo":  0.002,
		"text-ada-001":   0.0004,
		"text-babbage-001": 0.0005,
		"text-curie-001": 0.002,
	}

	pricePer1K, ok := pricing[modelID]
	if !ok {
		pricePer1K = 0.002 // 默认价格
	}

	return (float64(totalTokens) / 1000.0) * pricePer1K
}

// stringValue 辅助函数：将字符串指针转换为字符串
func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
