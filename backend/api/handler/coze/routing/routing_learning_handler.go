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

package routing

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	routingapp "github.com/coze-dev/coze-studio/backend/application/routing"
)

// RoutingLearningHandler 路由学习Handler
// 职责：处理路由学习相关的HTTP请求
type RoutingLearningHandler struct {
	learningService *routingapp.RoutingLearningService
	logger          *zap.Logger
}

// NewRoutingLearningHandler 创建路由学习Handler实例
func NewRoutingLearningHandler(
	learningService *routingapp.RoutingLearningService,
	logger *zap.Logger,
) *RoutingLearningHandler {
	return &RoutingLearningHandler{
		learningService: learningService,
		logger:          logger,
	}
}

// LearnFromLogsRequest 从日志学习请求
type LearnFromLogsRequest struct {
	TenantID    string  `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	StartTime   int64   `json:"start_time" binding:"required"` // Unix时间戳
	EndTime     int64   `json:"end_time" binding:"required"`   // Unix时间戳
	AnalysisDepth string `json:"analysis_depth" binding:"omitempty,oneof=quick standard deep"` // quick, standard, deep
}

// LearnFromLogsResponse 从日志学习响应
type LearnFromLogsResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    *LearningResultData   `json:"data,omitempty"`
}

// LearningResultData 学习结果数据
type LearningResultData struct {
	Insight             *RoutingInsightInfo      `json:"insight"`
	AgentPerformance    []AgentPerformanceStatsInfo `json:"agent_performance"`
	IntentDistribution  map[string]int            `json:"intent_distribution"`
	Recommendations     []string                 `json:"recommendations"`
	ExecutionTime       int64                    `json:"execution_time"` // 毫秒
}

// RoutingInsightInfo 路由洞察信息
type RoutingInsightInfo struct {
	TenantID           string                    `json:"tenant_id"`
	TimeRange          TimeRangeInfo             `json:"time_range"`
	TotalRequests      int                       `json:"total_requests"`
	AvgResponseTime    float64                   `json:"avg_response_time"`
	SuccessRate        float64                   `json:"success_rate"`
	IntentDistribution map[string]int            `json:"intent_distribution"`
	AgentPerformance   map[string]AgentPerformanceStatsInfo `json:"agent_performance"`
	Recommendations    []string                  `json:"recommendations"`
}

// TimeRangeInfo 时间范围信息
type TimeRangeInfo struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

// AgentPerformanceStatsInfo Agent性能统计信息
type AgentPerformanceStatsInfo struct {
	AgentID         string  `json:"agent_id"`
	RequestCount    int     `json:"request_count"`
	AvgResponseTime float64 `json:"avg_response_time"`
	SuccessRate     float64 `json:"success_rate"`
	AvgUserRating   float64 `json:"avg_user_rating"`
	ErrorRate       float64 `json:"error_rate"`
}

// LearnFromLogs 从日志学习
// @router /api/routing/learning/analyze [POST]
func (h *RoutingLearningHandler) LearnFromLogs(ctx context.Context, c *app.RequestContext) {
	var req LearnFromLogsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分析深度
	if req.AnalysisDepth == "" {
		req.AnalysisDepth = "standard"
	}

	// 调用Service层分析日志
	result, err := h.learningService.LearnFromLogs(ctx, req.TenantID, req.StartTime, req.EndTime, req.AnalysisDepth)
	if err != nil {
		h.logger.Error("failed to learn from logs",
			zap.String("tenant_id", req.TenantID),
			zap.Int64("start_time", req.StartTime),
			zap.Int64("end_time", req.EndTime),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	agentPerf := make([]AgentPerformanceStatsInfo, 0)
	if result.Insight != nil {
		for _, perf := range result.Insight.AgentPerformance {
			agentPerf = append(agentPerf, AgentPerformanceStatsInfo{
				AgentID:         perf.AgentID,
				RequestCount:    perf.RequestCount,
				AvgResponseTime: perf.AvgResponseTime,
				SuccessRate:     perf.SuccessRate,
				AvgUserRating:   perf.AvgUserRating,
				ErrorRate:       perf.ErrorRate,
			})
		}
	}

	insight := &RoutingInsightInfo{}
	if result.Insight != nil {
		agentPerfMap := make(map[string]AgentPerformanceStatsInfo)
		for _, perf := range agentPerf {
			agentPerfMap[perf.AgentID] = perf
		}

		insight = &RoutingInsightInfo{
			TenantID:           result.Insight.TenantID,
			TimeRange:          TimeRangeInfo{StartTime: result.Insight.TimeRange.StartTime, EndTime: result.Insight.TimeRange.EndTime},
			TotalRequests:      result.Insight.TotalRequests,
			AvgResponseTime:    result.Insight.AvgResponseTime,
			SuccessRate:        result.Insight.SuccessRate,
			IntentDistribution: result.Insight.IntentDistribution,
			AgentPerformance:   agentPerfMap,
			Recommendations:    result.Insight.Recommendations,
		}
	}

	c.JSON(http.StatusOK, LearnFromLogsResponse{
		Code:    0,
		Message: "success",
		Data: &LearningResultData{
			Insight:            insight,
			AgentPerformance:   agentPerf,
			IntentDistribution: result.IntentDistribution,
			Recommendations:    result.Recommendations,
			ExecutionTime:      result.ExecutionTime,
		},
	})
}

// SuggestRulesRequest 建议规则请求
type SuggestRulesRequest struct {
	TenantID      string  `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	StartTime     int64   `json:"start_time" binding:"required"`
	EndTime       int64   `json:"end_time" binding:"required"`
	MinConfidence float64 `json:"min_confidence" binding:"omitempty,min=0,max=1"` // 最小置信度
	MaxSuggestions int    `json:"max_suggestions" binding:"omitempty,min=1,max=50"` // 最大建议数
}

// SuggestRulesResponse 建议规则响应
type SuggestRulesResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    *RuleSuggestionData      `json:"data,omitempty"`
}

// RuleSuggestionData 规则建议数据
type RuleSuggestionData struct {
	TenantID       string                    `json:"tenant_id"`
	Suggestions    []RoutingRuleSuggestionInfo `json:"suggestions"`
	TotalAnalyzed  int                       `json:"total_analyzed"`
	ExecutionTime  int64                     `json:"execution_time"`
}

// RoutingRuleSuggestionInfo 路由规则建议信息
type RoutingRuleSuggestionInfo struct {
	Intent                string  `json:"intent"`
	SuggestedAgent        string  `json:"suggested_agent"`
	Confidence            float64 `json:"confidence"`
	Reason                string  `json:"reason"`
	ExpectedImprovement   float64 `json:"expected_improvement"`
	SampleSize            int     `json:"sample_size"`
	CurrentPerformance    float64 `json:"current_performance"`
	PredictedPerformance  float64 `json:"predicted_performance"`
}

// SuggestRules 建议规则
// @router /api/routing/learning/suggest-rules [POST]
func (h *RoutingLearningHandler) SuggestRules(ctx context.Context, c *app.RequestContext) {
	var req SuggestRulesRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.MinConfidence == 0 {
		req.MinConfidence = 0.7
	}
	if req.MaxSuggestions == 0 {
		req.MaxSuggestions = 10
	}

	// 调用Service层生成建议
	suggestions, err := h.learningService.SuggestRules(ctx, req.TenantID, req.StartTime, req.EndTime, req.MinConfidence, req.MaxSuggestions)
	if err != nil {
		h.logger.Error("failed to suggest rules",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	suggestionList := make([]RoutingRuleSuggestionInfo, len(suggestions.Suggestions))
	for i, sug := range suggestions.Suggestions {
		suggestionList[i] = RoutingRuleSuggestionInfo{
			Intent:               sug.Intent,
			SuggestedAgent:       sug.SuggestedAgent,
			Confidence:           sug.Confidence,
			Reason:               sug.Reason,
			ExpectedImprovement:  sug.ExpectedImprovement,
			SampleSize:           sug.SampleSize,
			CurrentPerformance:   sug.CurrentPerformance,
			PredictedPerformance: sug.PredictedPerformance,
		}
	}

	c.JSON(http.StatusOK, SuggestRulesResponse{
		Code:    0,
		Message: "success",
		Data: &RuleSuggestionData{
			TenantID:      suggestions.TenantID,
			Suggestions:   suggestionList,
			TotalAnalyzed: suggestions.TotalAnalyzed,
			ExecutionTime: suggestions.ExecutionTime,
		},
	})
}

// OptimizeStrategyRequest 优化策略请求
type OptimizeStrategyRequest struct {
	TenantID      string                 `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	OptimizationGoals []string           `json:"optimization_goals" binding:"required,min=1"` // 优化目标列表
	Constraints   map[string]interface{} `json:"constraints" binding:"omitempty"` // 约束条件
	ApplyChanges  bool                   `json:"apply_changes" binding:"omitempty"` // 是否应用更改
}

// OptimizeStrategyResponse 优化策略响应
type OptimizeStrategyResponse struct {
	Code    int                   `json:"code"`
	Message string                `json:"message"`
	Data    *OptimizationResultData `json:"data,omitempty"`
}

// OptimizationResultData 优化结果数据
type OptimizationResultData struct {
	StrategyID       string                 `json:"strategy_id"`
	TenantID         string                 `json:"strategy_id"`
	Applied          bool                   `json:"applied"`
	Changes          []StrategyChangeInfo   `json:"changes"`
	ExpectedMetrics  map[string]float64     `json:"expected_metrics"`
	ExecutionTime    int64                  `json:"execution_time"`
}

// StrategyChangeInfo 策略变更信息
type StrategyChangeInfo struct {
	ChangeType  string `json:"change_type"` // weight, threshold, rule, etc.
	Description string `json:"description"`
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	Reason      string `json:"reason"`
}

// OptimizeStrategy 优化策略
// @router /api/routing/learning/optimize-strategy [POST]
func (h *RoutingLearningHandler) OptimizeStrategy(ctx context.Context, c *app.RequestContext) {
	var req OptimizeStrategyRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层优化策略
	result, err := h.learningService.OptimizeStrategy(ctx, req.TenantID, req.OptimizationGoals, req.Constraints, req.ApplyChanges)
	if err != nil {
		h.logger.Error("failed to optimize strategy",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	changes := make([]StrategyChangeInfo, len(result.Changes))
	for i, change := range result.Changes {
		changes[i] = StrategyChangeInfo{
			ChangeType:  change.ChangeType,
			Description: change.Description,
			OldValue:    change.OldValue,
			NewValue:    change.NewValue,
			Reason:      change.Reason,
		}
	}

	c.JSON(http.StatusOK, OptimizeStrategyResponse{
		Code:    0,
		Message: "success",
		Data: &OptimizationResultData{
			StrategyID:      result.StrategyID,
			TenantID:        result.TenantID,
			Applied:         result.Applied,
			Changes:         changes,
			ExpectedMetrics: result.ExpectedMetrics,
			ExecutionTime:   result.ExecutionTime,
		},
	})
}

// GetLearningProgressRequest 获取学习进度请求
type GetLearningProgressRequest struct {
	TenantID string `form:"tenant_id" binding:"required" vd:"len($) > 0"`
	TaskID   string `form:"task_id" binding:"required"`
}

// GetLearningProgressResponse 获取学习进度响应
type GetLearningProgressResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    *LearningProgressInfo    `json:"data,omitempty"`
}

// LearningProgressInfo 学习进度信息
type LearningProgressInfo struct {
	TaskID        string  `json:"task_id"`
	TaskType      string  `json:"task_type"` // analyze, suggest, optimize
	Status        string  `json:"status"`    // pending, running, completed, failed
	Progress      float64 `json:"progress"`  // 0-100
	ProcessedCount int    `json:"processed_count"`
	TotalCount    int     `json:"total_count"`
	StartedAt     *int64  `json:"started_at,omitempty"`
	EstimatedEnd  *int64  `json:"estimated_end,omitempty"`
	Result        interface{} `json:"result,omitempty"`
	Error         string  `json:"error,omitempty"`
}

// GetLearningProgress 获取学习进度
// @router /api/routing/learning/progress [GET]
func (h *RoutingLearningHandler) GetLearningProgress(ctx context.Context, c *app.RequestContext) {
	var req GetLearningProgressRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层获取进度
	progress, err := h.learningService.GetLearningProgress(ctx, req.TenantID, req.TaskID)
	if err != nil {
		h.logger.Error("failed to get learning progress",
			zap.String("tenant_id", req.TenantID),
			zap.String("task_id", req.TaskID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, GetLearningProgressResponse{
		Code:    0,
		Message: "success",
		Data: &LearningProgressInfo{
			TaskID:         progress.TaskID,
			TaskType:       progress.TaskType,
			Status:         progress.Status,
			Progress:       progress.Progress,
			ProcessedCount: progress.ProcessedCount,
			TotalCount:     progress.TotalCount,
			StartedAt:      progress.StartedAt,
			EstimatedEnd:   progress.EstimatedEnd,
			Result:         progress.Result,
			Error:          progress.Error,
		},
	})
}

// ExportLearningReportRequest 导出学习报告请求
type ExportLearningReportRequest struct {
	TenantID  string  `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	StartTime int64   `json:"start_time" binding:"required"`
	EndTime   int64   `json:"end_time" binding:"required"`
	ReportFormat string `json:"report_format" binding:"omitempty,oneof=json pdf html"` // json, pdf, html
	IncludeCharts bool  `json:"include_charts" binding:"omitempty"`
}

// ExportLearningReportResponse 导出学习报告响应
type ExportLearningReportResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    *LearningReportData    `json:"data,omitempty"`
}

// LearningReportData 学习报告数据
type LearningReportData struct {
	ReportID    string `json:"report_id"`
	TenantID    string `json:"tenant_id"`
	TimeRange   TimeRangeInfo `json:"time_range"`
	Format      string `json:"format"`
	DownloadURL string `json:"download_url"`
	ExpiresAt   int64  `json:"expires_at"`
	Size        int64  `json:"size"` // 字节
}

// ExportLearningReport 导出学习报告
// @router /api/routing/learning/export-report [POST]
func (h *RoutingLearningHandler) ExportLearningReport(ctx context.Context, c *app.RequestContext) {
	var req ExportLearningReportRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.ReportFormat == "" {
		req.ReportFormat = "json"
	}

	// 调用Service层导出报告
	report, err := h.learningService.ExportLearningReport(ctx, req.TenantID, req.StartTime, req.EndTime, req.ReportFormat, req.IncludeCharts)
	if err != nil {
		h.logger.Error("failed to export learning report",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, ExportLearningReportResponse{
		Code:    0,
		Message: "success",
		Data: &LearningReportData{
			ReportID:    report.ReportID,
			TenantID:    report.TenantID,
			TimeRange:   TimeRangeInfo{StartTime: report.TimeRange.StartTime, EndTime: report.TimeRange.EndTime},
			Format:      report.Format,
			DownloadURL: report.DownloadURL,
			ExpiresAt:   report.ExpiresAt,
			Size:        report.Size,
		},
	})
}
