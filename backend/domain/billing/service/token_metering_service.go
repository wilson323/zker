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

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/billing/repository"
)

// ============================================================
// Token计量服务 - 基于设计文档: 23-租户计费系统_TokenMetering补充_完整版.md
// ============================================================

// TokenMeteringService Token计量服务
type TokenMeteringService struct {
	pricingEngine    *PricingEngine
	usageLogRepo     repository.TokenUsageLogRepository
	usageSummaryRepo repository.TokenUsageSummaryRepository
	budgetAlertSvc   *BudgetAlertService
	db               *gorm.DB
	logger           *logrus.Logger
}

// NewTokenMeteringService 创建Token计量服务实例
func NewTokenMeteringService(
	pricingEngine *PricingEngine,
	usageLogRepo repository.TokenUsageLogRepository,
	usageSummaryRepo repository.TokenUsageSummaryRepository,
	budgetAlertSvc *BudgetAlertService,
	db *gorm.DB,
	logger *logrus.Logger,
) *TokenMeteringService {
	return &TokenMeteringService{
		pricingEngine:    pricingEngine,
		usageLogRepo:     usageLogRepo,
		usageSummaryRepo: usageSummaryRepo,
		budgetAlertSvc:   budgetAlertSvc,
		db:               db,
		logger:           logger,
	}
}

// ============================================================
// 请求和响应类型
// ============================================================

// RecordTokenUsageRequest Token使用记录请求
type RecordTokenUsageRequest struct {
	TenantID       string  `json:"tenant_id" validate:"required"`
	UserID         *uint64 `json:"user_id,omitempty"`
	BotID          *string `json:"bot_id,omitempty"`
	ConversationID *string `json:"conversation_id,omitempty"`
	MessageID      *string `json:"message_id,omitempty"`

	ModelProvider string  `json:"model_provider" validate:"required"` // openai, anthropic, etc
	ModelName     string  `json:"model_name" validate:"required"`    // gpt-4, claude-3-opus, etc
	ModelVersion  *string `json:"model_version,omitempty"`

	InputTokens  int    `json:"input_tokens" validate:"required,min=0"`
	OutputTokens int    `json:"output_tokens" validate:"required,min=0"`
	TotalTokens  int    `json:"total_tokens" validate:"required,min=0"`

	RequestType    string `json:"request_type" validate:"required"`    // chat, completion, embedding
	ResponseTimeMs *int   `json:"response_time_ms,omitempty"`
	LatencyMs      *int   `json:"latency_ms,omitempty"`
	IsCached       bool   `json:"is_cached"`
	Metadata       string `json:"metadata,omitempty"`
}

// RecordTokenUsageResponse Token使用记录响应
type RecordTokenUsageResponse struct {
	LogID      uint64  `json:"log_id"`
	InputCost  float64 `json:"input_cost"`
	OutputCost float64 `json:"output_cost"`
	TotalCost  float64 `json:"total_cost"`
	CreatedAt  int64   `json:"created_at"`
}

// BatchRecordResponse 批量记录响应
type BatchRecordResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	LogIDs       []uint64 `json:"log_ids"`
	TotalCost    float64  `json:"total_cost"`
}

// UsageStatsFilter 使用统计过滤器
type UsageStatsFilter struct {
	BotID        *string `json:"bot_id,omitempty"`
	ModelName    *string `json:"model_name,omitempty"`
	StartDate    string  `json:"start_date,omitempty"`
	EndDate      string  `json:"end_date,omitempty"`
	Granularity  string  `json:"granularity,omitempty"` // daily, hourly
}

// UsageStatsResponse 使用统计响应
type UsageStatsResponse struct {
	TenantID          string  `json:"tenant_id"`
	TotalInputTokens  int64   `json:"total_input_tokens"`
	TotalOutputTokens int64   `json:"total_output_tokens"`
	TotalTokens       int64   `json:"total_tokens"`
	TotalCost         float64 `json:"total_cost"`
	TotalRequests     int64   `json:"total_requests"`
	CachedRequests    int64   `json:"cached_requests"`
	AvgResponseTime   float64 `json:"avg_response_time_ms"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
}

// ModelUsageStats 模型使用统计
type ModelUsageStats struct {
	ModelProvider     string  `json:"model_provider"`
	ModelName         string  `json:"model_name"`
	TotalTokens       int64   `json:"total_tokens"`
	TotalCost         float64 `json:"total_cost"`
	TotalRequests     int64   `json:"total_requests"`
	AvgCostPer1KToken float64 `json:"avg_cost_per_1k_token"`
}

// DailyUsageStats 每日使用趋势
type DailyUsageStats struct {
	Date             string  `json:"date"`
	TotalTokens      int64   `json:"total_tokens"`
	TotalCost        float64 `json:"total_cost"`
	TotalRequests    int64   `json:"total_requests"`
	CachedRequests   int64   `json:"cached_requests"`
	AvgResponseTime  float64 `json:"avg_response_time_ms"`
}

// ============================================================
// 核心方法实现
// ============================================================

// RecordTokenUsage 记录Token使用
// 参数:
//   - ctx: 上下文
//   - req: Token使用记录请求
// 返回: 记录响应或错误
func (s *TokenMeteringService) RecordTokenUsage(
	ctx context.Context,
	req *RecordTokenUsageRequest,
) (*RecordTokenUsageResponse, error) {
	// 1. 验证请求参数
	if err := s.validateRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 2. 计算成本
	cost := s.pricingEngine.CalculateCost(
		req.ModelProvider,
		req.ModelName,
		req.InputTokens,
		req.OutputTokens,
	)

	inputCost, _ := cost.InputCost.Float64()
	outputCost, _ := cost.OutputCost.Float64()
	totalCost, _ := cost.TotalCost.Float64()
	unitPrice, _ := cost.UnitPrice.Float64()

	// 3. 创建TokenUsageLog记录
	log := &entity.TokenUsageLog{
		TenantID:       req.TenantID,
		UserID:         req.UserID,
		BotID:          req.BotID,
		ConversationID: req.ConversationID,
		MessageID:      req.MessageID,
		InputTokens:    req.InputTokens,
		OutputTokens:   req.OutputTokens,
		TotalTokens:    req.TotalTokens,
		ModelProvider:  req.ModelProvider,
		ModelName:      req.ModelName,
		ModelVersion:   req.ModelVersion,
		UnitPrice:      unitPrice,
		InputCost:      inputCost,
		OutputCost:     outputCost,
		TotalCost:      totalCost,
		ResponseTimeMs: req.ResponseTimeMs,
		LatencyMs:      req.LatencyMs,
		IsCached:       req.IsCached,
		RequestType:    req.RequestType,
		Metadata:       req.Metadata,
	}

	// 4. 保存日志（使用事务）
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 4.1 保存日志记录
		if err := s.usageLogRepo.Create(ctx, log); err != nil {
			s.logger.WithError(err).Errorf("[TokenMetering] Failed to create usage log")
			return fmt.Errorf("failed to create usage log: %w", err)
		}

		// 4.2 异步更新汇总数据
		go s.updateSummaryAsync(context.Background(), req.TenantID, req.BotID, log)

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 5. 触发预算告警检查（异步）
	go func() {
		if checkErr := s.budgetAlertSvc.CheckBudget(context.Background(), req.TenantID); checkErr != nil {
			s.logger.WithError(checkErr).
				WithField("tenant_id", req.TenantID).
				Warn("[TokenMetering] Failed to check budget after recording usage")
		}
	}()

	// 6. 记录日志
	s.logger.WithFields(logrus.Fields{
		"tenant_id":      req.TenantID,
		"log_id":         log.ID,
		"model":          fmt.Sprintf("%s/%s", req.ModelProvider, req.ModelName),
		"input_tokens":   req.InputTokens,
		"output_tokens":  req.OutputTokens,
		"total_cost":     totalCost,
	}).Info("[TokenMetering] Token usage recorded")

	// 7. 返回响应
	return &RecordTokenUsageResponse{
		LogID:      log.ID,
		InputCost:  inputCost,
		OutputCost: outputCost,
		TotalCost:  totalCost,
		CreatedAt:  log.CreatedAt.Unix(),
	}, nil
}

// BatchRecordTokenUsage 批量记录Token使用
// 参数:
//   - ctx: 上下文
//   - records: Token使用记录请求列表
// 返回: 批量记录响应或错误
func (s *TokenMeteringService) BatchRecordTokenUsage(
	ctx context.Context,
	records []*RecordTokenUsageRequest,
) (*BatchRecordResponse, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no records to batch")
	}

	if len(records) > 1000 {
		return nil, fmt.Errorf("batch size exceeds limit of 1000")
	}

	startTime := time.Now()
	successCount := 0
	failedCount := 0
	logIDs := make([]uint64, 0, len(records))
	totalCost := 0.0

	// 1. 准备日志记录
	logs := make([]*entity.TokenUsageLog, 0, len(records))
	for _, req := range records {
		// 验证请求
		if err := s.validateRequest(req); err != nil {
			s.logger.WithError(err).Warnf("[TokenMetering] Invalid record in batch, skipping")
			failedCount++
			continue
		}

		// 计算成本
		cost := s.pricingEngine.CalculateCost(
			req.ModelProvider,
			req.ModelName,
			req.InputTokens,
			req.OutputTokens,
		)

		inputCost, _ := cost.InputCost.Float64()
		outputCost, _ := cost.OutputCost.Float64()
		totalCost, _ := cost.TotalCost.Float64()
		unitPrice, _ := cost.UnitPrice.Float64()

		log := &entity.TokenUsageLog{
			TenantID:       req.TenantID,
			UserID:         req.UserID,
			BotID:          req.BotID,
			ConversationID: req.ConversationID,
			MessageID:      req.MessageID,
			InputTokens:    req.InputTokens,
			OutputTokens:   req.OutputTokens,
			TotalTokens:    req.TotalTokens,
			ModelProvider:  req.ModelProvider,
			ModelName:      req.ModelName,
			ModelVersion:   req.ModelVersion,
			UnitPrice:      unitPrice,
			InputCost:      inputCost,
			OutputCost:     outputCost,
			TotalCost:      totalCost,
			ResponseTimeMs: req.ResponseTimeMs,
			LatencyMs:      req.LatencyMs,
			IsCached:       req.IsCached,
			RequestType:    req.RequestType,
			Metadata:       req.Metadata,
		}

		logs = append(logs, log)
		totalCost += totalCost
	}

	// 2. 批量保存
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.usageLogRepo.BatchCreate(ctx, logs); err != nil {
			s.logger.WithError(err).Errorf("[TokenMetering] Failed to batch create usage logs")
			return fmt.Errorf("failed to batch create logs: %w", err)
		}

		// 收集ID
		for _, log := range logs {
			logIDs = append(logIDs, log.ID)
		}
		successCount = len(logs)

		// 异步更新汇总
		for _, log := range logs {
			go s.updateSummaryAsync(context.Background(), log.TenantID, log.BotID, log)
		}

		return nil
	})

	if err != nil {
		failedCount = len(records)
		return &BatchRecordResponse{
			SuccessCount: 0,
			FailedCount:  failedCount,
			LogIDs:       []uint64{},
			TotalCost:    0,
		}, err
	}

	// 3. 记录性能指标
	duration := time.Since(startTime).Milliseconds()
	s.logger.WithFields(logrus.Fields{
		"total_records":  len(records),
		"success_count":  successCount,
		"failed_count":   failedCount,
		"total_cost":     totalCost,
		"duration_ms":    duration,
	}).Infof("[TokenMetering] Batch recording completed in %dms", duration)

	// 4. 触发预算告警检查（仅针对第一个租户）
	if len(logs) > 0 {
		go func() {
			if checkErr := s.budgetAlertSvc.CheckBudget(context.Background(), logs[0].TenantID); checkErr != nil {
				s.logger.WithError(checkErr).Warn("[TokenMetering] Failed to check budget after batch recording")
			}
		}()
	}

	return &BatchRecordResponse{
		SuccessCount: successCount,
		FailedCount:  failedCount,
		LogIDs:       logIDs,
		TotalCost:    totalCost,
	}, nil
}

// GetUsageStats 获取使用统计
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - filter: 统计过滤器
// 返回: 统计响应或错误
func (s *TokenMeteringService) GetUsageStats(
	ctx context.Context,
	tenantID string,
	filter *UsageStatsFilter,
) (*UsageStatsResponse, error) {
	// 1. 设置默认时间范围（最近7天）
	startDate := filter.StartDate
	endDate := filter.EndDate

	if startDate == "" || endDate == "" {
		now := time.Now()
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
	}

	// 2. 查询汇总数据
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}

	summaries, err := s.usageSummaryRepo.GetByDateRange(
		ctx,
		tenantID,
		filter.BotID,
		startDate,
		endDate,
	)

	if err != nil {
		s.logger.WithError(err).Errorf("[TokenMetering] Failed to get usage summaries")
		return nil, fmt.Errorf("failed to get usage summaries: %w", err)
	}

	// 3. 聚合统计数据
	stats := &UsageStatsResponse{
		TenantID:  tenantID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	var totalResponseTime float64
	responseTimeCount := 0

	for _, summary := range summaries {
		stats.TotalInputTokens += summary.TotalInputTokens
		stats.TotalOutputTokens += summary.TotalOutputTokens
		stats.TotalTokens += summary.TotalTokens
		stats.TotalCost += summary.TotalCost
		stats.TotalRequests += int64(summary.TotalRequests)
		stats.CachedRequests += int64(summary.CachedRequests)

		if summary.AvgResponseTime != nil {
			totalResponseTime += *summary.AvgResponseTime
			responseTimeCount++
		}
	}

	// 4. 计算平均响应时间
	if responseTimeCount > 0 {
		stats.AvgResponseTime = totalResponseTime / float64(responseTimeCount)
	}

	s.logger.WithFields(logrus.Fields{
		"tenant_id":       tenantID,
		"total_tokens":    stats.TotalTokens,
		"total_cost":      stats.TotalCost,
		"total_requests":  stats.TotalRequests,
	}).Info("[TokenMetering] Usage stats retrieved")

	return stats, nil
}

// GetModelUsageStats 获取模型使用统计
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
// 返回: 模型统计列表或错误
func (s *TokenMeteringService) GetModelUsageStats(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
) ([]*ModelUsageStats, error) {
	// 1. 验证日期格式
	_, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	_, err = time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date: %w", err)
	}

	// 2. 查询原始日志
	logs, err := s.usageLogRepo.GetByDateRange(ctx, tenantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage logs: %w", err)
	}

	// 2. 按模型聚合统计
	modelStatsMap := make(map[string]*ModelUsageStats)

	for _, log := range logs {
		key := fmt.Sprintf("%s/%s", log.ModelProvider, log.ModelName)

		if stats, ok := modelStatsMap[key]; ok {
			stats.TotalTokens += int64(log.TotalTokens)
			stats.TotalCost += log.TotalCost
			stats.TotalRequests++
		} else {
			modelStatsMap[key] = &ModelUsageStats{
				ModelProvider: log.ModelProvider,
				ModelName:     log.ModelName,
				TotalTokens:   int64(log.TotalTokens),
				TotalCost:     log.TotalCost,
				TotalRequests: 1,
			}
		}
	}

	// 3. 转换为数组并计算平均成本
	result := make([]*ModelUsageStats, 0, len(modelStatsMap))
	for _, stats := range modelStatsMap {
		// 计算每1K token的平均成本
		if stats.TotalTokens > 0 {
			stats.AvgCostPer1KToken = (stats.TotalCost / float64(stats.TotalTokens)) * 1000
		}
		result = append(result, stats)
	}

	// 4. 按总成本降序排序（这里简化处理，实际可以使用sort包）

	s.logger.WithFields(logrus.Fields{
		"tenant_id":     tenantID,
		"model_count":   len(result),
		"start_date":    startDate,
		"end_date":      endDate,
	}).Info("[TokenMetering] Model usage stats retrieved")

	return result, nil
}

// GetDailyUsageStats 获取每日使用趋势
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - startDate: 开始日期
//   - endDate: 结束日期
//   - botID: Bot ID（可选）
// 返回: 每日统计列表或错误
func (s *TokenMeteringService) GetDailyUsageStats(
	ctx context.Context,
	tenantID string,
	startDate, endDate string,
	botID *string,
) ([]*DailyUsageStats, error) {
	// 1. 查询汇总数据
	summaries, err := s.usageSummaryRepo.GetByDateRange(
		ctx,
		tenantID,
		botID,
		startDate,
		endDate,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get daily summaries: %w", err)
	}

	// 2. 构建每日统计
	result := make([]*DailyUsageStats, 0, len(summaries))
	for _, summary := range summaries {
		daily := &DailyUsageStats{
			Date:           summary.SummaryDate,
			TotalTokens:    summary.TotalTokens,
			TotalCost:      summary.TotalCost,
			TotalRequests:  int64(summary.TotalRequests),
			CachedRequests: int64(summary.CachedRequests),
		}

		if summary.AvgResponseTime != nil {
			daily.AvgResponseTime = *summary.AvgResponseTime
		}

		result = append(result, daily)
	}

	s.logger.WithFields(logrus.Fields{
		"tenant_id":    tenantID,
		"days_count":   len(result),
		"start_date":   startDate,
		"end_date":     endDate,
	}).Info("[TokenMetering] Daily usage stats retrieved")

	return result, nil
}

// ============================================================
// 私有辅助方法
// ============================================================

// validateRequest 验证请求参数
func (s *TokenMeteringService) validateRequest(req *RecordTokenUsageRequest) error {
	if req.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	if req.ModelProvider == "" {
		return fmt.Errorf("model_provider is required")
	}

	if req.ModelName == "" {
		return fmt.Errorf("model_name is required")
	}

	if req.InputTokens < 0 || req.OutputTokens < 0 || req.TotalTokens < 0 {
		return fmt.Errorf("token counts cannot be negative")
	}

	if req.TotalTokens != req.InputTokens+req.OutputTokens {
		return fmt.Errorf("total_tokens must equal input_tokens + output_tokens")
	}

	if req.RequestType == "" {
		return fmt.Errorf("request_type is required")
	}

	// 验证请求类型
	validTypes := map[string]bool{
		entity.RequestTypeChat:       true,
		entity.RequestTypeCompletion: true,
		entity.RequestTypeEmbedding:  true,
		entity.RequestTypeRerank:     true,
	}

	if !validTypes[req.RequestType] {
		return fmt.Errorf("invalid request_type: %s", req.RequestType)
	}

	return nil
}

// updateSummaryAsync 异步更新汇总数据
func (s *TokenMeteringService) updateSummaryAsync(
	ctx context.Context,
	tenantID string,
	botID *string,
	log *entity.TokenUsageLog,
) {
	// 获取当前日期和小时
	now := time.Now()
	summaryDate := now.Format("2006-01-02")
	summaryHour := uint8(now.Hour())

	// 查询或创建汇总记录
	summary, err := s.usageSummaryRepo.GetByTenantAndDate(ctx, tenantID, summaryDate, &summaryHour)
	if err != nil {
		s.logger.WithError(err).Warnf("[TokenMetering] Failed to get summary, will create new")
		summary = &entity.TokenUsageSummary{
			TenantID:    tenantID,
			BotID:       botID,
			SummaryDate: summaryDate,
			SummaryHour: &summaryHour,
		}
	}

	// 更新统计数据
	summary.TotalInputTokens += int64(log.InputTokens)
	summary.TotalOutputTokens += int64(log.OutputTokens)
	summary.TotalTokens += int64(log.TotalTokens)
	summary.TotalCost += log.TotalCost
	summary.TotalRequests++

	if log.IsCached {
		summary.CachedRequests++
	}

	// 更新平均响应时间
	if log.ResponseTimeMs != nil {
		if summary.AvgResponseTime == nil {
			avg := float64(*log.ResponseTimeMs)
			summary.AvgResponseTime = &avg
		} else {
			// 简单的移动平均
			*summary.AvgResponseTime = (*summary.AvgResponseTime*float64(summary.TotalRequests-1) + float64(*log.ResponseTimeMs)) / float64(summary.TotalRequests)
		}
	}

	// 保存汇总记录
	if err := s.usageSummaryRepo.Upsert(ctx, summary); err != nil {
		s.logger.WithError(err).Errorf("[TokenMetering] Failed to upsert summary")
	}

	s.logger.WithFields(logrus.Fields{
		"tenant_id":     tenantID,
		"summary_date":  summaryDate,
		"summary_hour":  summaryHour,
		"total_tokens":  summary.TotalTokens,
		"total_cost":    summary.TotalCost,
	}).Debug("[TokenMetering] Summary updated")
}
