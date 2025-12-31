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
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/agent_monitoring/entity"
	"github.com/coze-dev/coze-studio/backend/domain/agent_monitoring/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// TraceService 链路追踪服务
// 职责：实现Agent调用链追踪,集成OpenTelemetry
type TraceService struct {
	traceRepo   repository.TraceRepository
	logger      *zap.Logger
	batchSize   int
	flushInterval time.Duration
}

// NewTraceService 创建链路追踪服务实例
func NewTraceService(
	traceRepo repository.TraceRepository,
	logger *zap.Logger,
) *TraceService {
	return &TraceService{
		traceRepo:     traceRepo,
		logger:        logger,
		batchSize:     100,
		flushInterval: 5 * time.Second,
	}
}

// StartTrace 开始链路追踪
func (s *TraceService) StartTrace(
	ctx context.Context,
	tenantID, userID, agentID, conversationID, messageID, spanName string,
) (*entity.Trace, context.Context) {
	traceID := entity.TraceID(uuid.New().String())
	rootSpanID := entity.SpanID(uuid.New().String())

	trace := &entity.Trace{
		TraceID:        traceID,
		TenantID:       tenantID,
		UserID:         userID,
		AgentID:        agentID,
		ConversationID: conversationID,
		MessageID:      messageID,
		ParentSpanID:   rootSpanID,
		RootSpanID:     rootSpanID,
		SpanName:       spanName,
		SpanKind:       "internal",
		StartTime:      time.Now(),
		Status:         "running",
	}

	// 将traceID存入context
	ctx = context.WithValue(ctx, "trace_id", traceID)
	ctx = context.WithValue(ctx, "span_id", rootSpanID)

	return trace, ctx
}

// CreateSpan 创建子跨度
func (s *TraceService) CreateSpan(
	ctx context.Context,
	spanName, spanKind string,
	parentSpanID entity.SpanID,
) *entity.TraceSpan {
	_, _ = ctx.Value("trace_id").(entity.TraceID)
	spanID := entity.SpanID(uuid.New().String())

	if parentSpanID == "" {
		if currentSpanID, ok := ctx.Value("span_id").(entity.SpanID); ok {
			parentSpanID = currentSpanID
		}
	}

	span := &entity.TraceSpan{
		SpanID:       spanID,
		ParentSpanID: parentSpanID,
		SpanName:     spanName,
		StartTime:    time.Now(),
		Attributes:   make(map[string]string),
		Events:       make([]entity.TraceEvent, 0),
		Status:       "running",
	}

	// 更新context中的当前spanID
	ctx = context.WithValue(ctx, "span_id", spanID)

	return span
}

// EndSpan 结束跨度
func (s *TraceService) EndSpan(
	ctx context.Context,
	span *entity.TraceSpan,
	status string,
) {
	if span == nil {
		return
	}

	span.EndTime = time.Now()
	span.Duration = span.EndTime.Sub(span.StartTime)
	span.Status = status
}

// AddSpanEvent 添加跨度事件
func (s *TraceService) AddSpanEvent(
	span *entity.TraceSpan,
	eventName string,
	attributes map[string]string,
) {
	if span == nil {
		return
	}

	event := entity.TraceEvent{
		Timestamp: time.Now(),
		Name:      eventName,
		Attributes: attributes,
	}

	span.Events = append(span.Events, event)
}

// FinishTrace 完成链路追踪并保存
func (s *TraceService) FinishTrace(
	ctx context.Context,
	trace *entity.Trace,
	status, statusCode, statusMessage string,
) error {
	if trace == nil {
		return errorx.New(errno.ErrInvalidParamCode, errorx.KV("reason", "nil trace"))
	}

	trace.EndTime = time.Now()
	trace.Duration = int64(trace.EndTime.Sub(trace.StartTime).Milliseconds())
	trace.Status = status
	trace.StatusCode = statusCode
	trace.StatusMessage = statusMessage

	// 保存到数据库
	if err := s.traceRepo.Create(ctx, trace); err != nil {
		s.logger.Error("failed to save trace",
			zap.String("trace_id", string(trace.TraceID)),
			zap.Error(err))
		return errorx.Wrapf(err, errno.ErrTraceCreateFailedCode, errorx.KV("error", err.Error()))
	}

	s.logger.Debug("trace saved",
		zap.String("trace_id", string(trace.TraceID)),
		zap.String("status", status),
		zap.Int64("duration_ms", trace.Duration))

	return nil
}

// GetTrace 获取链路详情
func (s *TraceService) GetTrace(
	ctx context.Context,
	traceID entity.TraceID,
) (*entity.Trace, error) {
	if traceID == "" {
		return nil, errorx.New(errno.ErrInvalidParamCode, errorx.KV("reason", "empty trace_id"))
	}

	trace, err := s.traceRepo.GetByID(ctx, string(traceID))
	if err != nil {
		return nil, errorx.Wrapf(err, errno.ErrTraceNotFoundCode, errorx.KV("trace_id", string(traceID)))
	}

	return trace, nil
}

// QueryTraces 查询链路列表
func (s *TraceService) QueryTraces(
	ctx context.Context,
	filter *entity.TraceFilter,
) ([]*entity.Trace, int64, error) {
	if filter.TenantID == "" {
		return nil, 0, errorx.New(errno.ErrInvalidParamCode, errorx.KV("reason", "empty tenant_id"))
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	traces, total, err := s.traceRepo.Query(ctx, filter)
	if err != nil {
		return nil, 0, errorx.Wrapf(err, errno.ErrTraceQueryFailedCode, errorx.KV("error", err.Error()))
	}

	return traces, total, nil
}

// GetTraceTree 获取链路树(包含所有子跨度)
func (s *TraceService) GetTraceTree(
	ctx context.Context,
	traceID entity.TraceID,
) (*TraceTree, error) {
	trace, err := s.GetTrace(ctx, traceID)
	if err != nil {
		return nil, err
	}

	// 获取所有相关的spans
	spans, err := s.traceRepo.GetSpansByTraceID(ctx, string(traceID))
	if err != nil {
		return nil, errorx.Wrapf(err, errno.ErrTraceQueryFailedCode, errorx.KV("error", err.Error()))
	}

	// 构建树结构
	tree := &TraceTree{
		RootTrace: trace,
		Spans:     spans,
	}

	return tree, nil
}

// AnalyzePerformance 分析性能瓶颈
func (s *TraceService) AnalyalyzePerformance(
	ctx context.Context,
	traceID entity.TraceID,
) (*PerformanceAnalysis, error) {
	trace, err := s.GetTrace(ctx, traceID)
	if err != nil {
		return nil, err
	}

	analysis := &PerformanceAnalysis{
		TraceID:        traceID,
		TotalDuration:  trace.GetDurationMs(),
		Bottlenecks:    make([]string, 0),
		SlowOperations: make([]SlowOperation, 0),
	}

	// 获取所有spans
	spans, err := s.traceRepo.GetSpansByTraceID(ctx, string(traceID))
	if err != nil {
		return nil, err
	}

	// 分析每个span的性能
	totalDuration := trace.GetDurationMs()
	for _, span := range spans {
		spanDuration := int64(span.Duration.Milliseconds())
		percentage := float64(spanDuration) / float64(totalDuration) * 100

		// 识别慢操作(超过总时间的10%)
		if percentage > 10 {
			slowOp := SlowOperation{
				SpanName:     span.SpanName,
				Duration:     spanDuration,
				Percentage:   percentage,
				Attributes:   span.Attributes,
			}
			analysis.SlowOperations = append(analysis.SlowOperations, slowOp)

			bottleneck := fmt.Sprintf("%s 占用 %.1f%% 时间", span.SpanName, percentage)
			analysis.Bottlenecks = append(analysis.Bottlenecks, bottleneck)
		}
	}

	// 按持续时间排序
	for i := 0; i < len(analysis.SlowOperations); i++ {
		for j := i + 1; j < len(analysis.SlowOperations); j++ {
			if analysis.SlowOperations[j].Duration > analysis.SlowOperations[i].Duration {
				analysis.SlowOperations[i], analysis.SlowOperations[j] = analysis.SlowOperations[j], analysis.SlowOperations[i]
			}
		}
	}

	return analysis, nil
}

// BatchFinishTraces 批量完成链路追踪
func (s *TraceService) BatchFinishTraces(
	ctx context.Context,
	traces []*entity.Trace,
) error {
	if len(traces) == 0 {
		return nil
	}

	// 批量保存
	for _, trace := range traces {
		if err := s.traceRepo.Create(ctx, trace); err != nil {
			s.logger.Error("failed to save trace in batch",
				zap.String("trace_id", string(trace.TraceID)),
				zap.Error(err))
			continue
		}
	}

	s.logger.Debug("batch traces saved",
		zap.Int("count", len(traces)))

	return nil
}

// SetTraceAttributes 设置追踪属性
func (s *TraceService) SetTraceAttributes(
	trace *entity.Trace,
	attributes map[string]string,
) error {
	if trace == nil {
		return errorx.New(errno.ErrInvalidParamCode, errorx.KV("reason", "nil trace"))
	}

	// 序列化为JSON
	data, err := json.Marshal(attributes)
	if err != nil {
		return errorx.New(errno.InternalErrorCode, errorx.KV("reason", "failed to marshal attributes"))
	}

	trace.Attributes = string(data)
	return nil
}

// GetTraceStatistics 获取追踪统计信息
func (s *TraceService) GetTraceStatistics(
	ctx context.Context,
	tenantID string,
	startTime, endTime time.Time,
) (*TraceStatistics, error) {
	filter := &entity.TraceFilter{
		TenantID:  tenantID,
		StartTime: startTime,
		EndTime:   endTime,
	}

	traces, _, err := s.traceRepo.Query(ctx, filter)
	if err != nil {
		return nil, err
	}

	stats := &TraceStatistics{
		TotalTraces:    len(traces),
		SuccessTraces:  0,
		ErrorTraces:    0,
		TimeoutTraces:  0,
		AvgDuration:    0,
		MaxDuration:    int64(0),
		MinDuration:    int64(0),
		StatusDistribution: make(map[string]int),
	}

	if len(traces) == 0 {
		return stats, nil
	}

	totalDuration := int64(0)
	minDuration := int64(-1)
	maxDuration := int64(0)

	for _, trace := range traces {
		// 统计状态
		stats.StatusDistribution[trace.Status]++
		switch trace.Status {
		case "success":
			stats.SuccessTraces++
		case "error":
			stats.ErrorTraces++
		case "timeout":
			stats.TimeoutTraces++
		}

		// 统计时长
		duration := trace.Duration
		totalDuration += duration

		if minDuration == -1 || duration < minDuration {
			stats.MinDuration = duration
			minDuration = duration
		}
		if duration > maxDuration {
			stats.MaxDuration = duration
			maxDuration = duration
		}
	}

	stats.AvgDuration = totalDuration / int64(len(traces))

	return stats, nil
}

// TraceTree 链路树
type TraceTree struct {
	RootTrace *entity.Trace      `json:"root_trace"`
	Spans     []*entity.TraceSpan `json:"spans"`
}

// PerformanceAnalysis 性能分析结果
type PerformanceAnalysis struct {
	TraceID         entity.TraceID    `json:"trace_id"`
	TotalDuration   int64            `json:"total_duration"`   // 毫秒
	Bottlenecks     []string         `json:"bottlenecks"`      // 性能瓶颈列表
	SlowOperations  []SlowOperation  `json:"slow_operations"`  // 慢操作列表
}

// SlowOperation 慢操作
type SlowOperation struct {
	SpanName   string            `json:"span_name"`
	Duration   int64             `json:"duration"`    // 毫秒
	Percentage float64           `json:"percentage"`  // 占总时间的百分比
	Attributes map[string]string `json:"attributes"`
}

// TraceStatistics 追踪统计信息
type TraceStatistics struct {
	TotalTraces       int            `json:"total_traces"`
	SuccessTraces     int            `json:"success_traces"`
	ErrorTraces       int            `json:"error_traces"`
	TimeoutTraces     int            `json:"timeout_traces"`
	AvgDuration       int64          `json:"avg_duration"`       // 平均时长(毫秒)
	MaxDuration       int64          `json:"max_duration"`       // 最大时长(毫秒)
	MinDuration       int64          `json:"min_duration"`       // 最小时长(毫秒)
	StatusDistribution map[string]int `json:"status_distribution"` // 状态分布
}
