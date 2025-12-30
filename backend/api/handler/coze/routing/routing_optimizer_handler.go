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

// RoutingOptimizerHandler 路由优化Handler
// 职责：处理路由优化相关的HTTP请求
type RoutingOptimizerHandler struct {
	optimizer *routingapp.RoutingOptimizerService
	logger    *zap.Logger
}

// NewRoutingOptimizerHandler 创建路由优化Handler实例
func NewRoutingOptimizerHandler(
	optimizer *routingapp.RoutingOptimizerService,
	logger *zap.Logger,
) *RoutingOptimizerHandler {
	return &RoutingOptimizerHandler{
		optimizer: optimizer,
		logger:    logger,
	}
}

// OptimizeRoutingRequest 优化路由请求
type OptimizeRoutingRequest struct {
	TenantID      string                  `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	UserInput     string                  `json:"user_input" binding:"required" vd:"len($) > 0"`
	Candidates    []AgentCandidate        `json:"candidates" binding:"required,min=1,max=10"`
	OptimizeLevel string                  `json:"optimize_level" binding:"omitempty,oneof=fast balanced accurate"` // fast, balanced, accurate
	Context       map[string]interface{}  `json:"context" binding:"omitempty"`
}

// AgentCandidate Agent候选者
type AgentCandidate struct {
	AgentID          string                 `json:"agent_id" binding:"required"`
	AgentName        string                 `json:"agent_name" binding:"required"`
	AgentType        string                 `json:"agent_type" binding:"required"` // single_agent, workflow
	Description      string                 `json:"description"`
	Capabilities     []string               `json:"capabilities"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// OptimizeRoutingResponse 优化路由响应
type OptimizeRoutingResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *OptimizeResult   `json:"data,omitempty"`
}

// OptimizeResult 优化结果
type OptimizeResult struct {
	SelectedAgent   *AgentMatchInfo   `json:"selected_agent"`
	AllScores       []AgentScoreInfo  `json:"all_scores"`
	OptimizationLog []OptimizationLog `json:"optimization_log"`
	ExecutionTime   int64             `json:"execution_time"` // 毫秒
}

// AgentMatchInfo Agent匹配信息
type AgentMatchInfo struct {
	AgentID        string                 `json:"agent_id"`
	AgentName      string                 `json:"agent_name"`
	AgentType      string                 `json:"agent_type"`
	MatchScore     float64               `json:"match_score"`
	MatchReason    string                `json:"match_reason"`
	Confidence     float64               `json:"confidence"`
	Recommendation bool                  `json:"recommendation"`
}

// AgentScoreInfo Agent评分信息
type AgentScoreInfo struct {
	AgentID        string                 `json:"agent_id"`
	AgentName      string                 `json:"agent_name"`
	TotalScore     float64               `json:"total_score"`
	ScoreBreakdown ScoreBreakdown        `json:"score_breakdown"`
}

// ScoreBreakdown 评分明细
type ScoreBreakdown struct {
	LoadScore      float64 `json:"load_score"`
	HealthScore    float64 `json:"health_score"`
	CapabilityScore float64 `json:"capability_score"`
	PerformanceScore float64 `json:"performance_score"`
	ContextScore   float64 `json:"context_score"`
}

// OptimizationLog 优化日志
type OptimizationLog struct {
	Step    string `json:"step"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// OptimizeRouting 优化路由
// @router /api/routing/optimize [POST]
func (h *RoutingOptimizerHandler) OptimizeRouting(ctx context.Context, c *app.RequestContext) {
	var req OptimizeRoutingRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认优化级别
	if req.OptimizeLevel == "" {
		req.OptimizeLevel = "balanced"
	}

	// 转换为Service层请求格式
	candidates := make([]routingapp.AgentCandidate, len(req.Candidates))
	for i, cand := range req.Candidates {
		candidates[i] = routingapp.AgentCandidate{
			AgentID:      cand.AgentID,
			AgentName:    cand.AgentName,
			AgentType:    cand.AgentType,
			Description:  cand.Description,
			Capabilities: cand.Capabilities,
			Metadata:     cand.Metadata,
		}
	}

	// 调用Service层优化路由
	result, err := h.optimizer.OptimizeRouting(ctx, req.TenantID, req.UserInput, candidates, req.OptimizeLevel, req.Context)
	if err != nil {
		h.logger.Error("failed to optimize routing",
			zap.String("tenant_id", req.TenantID),
			zap.String("user_input", req.UserInput),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	response := &OptimizeResult{
		SelectedAgent: &AgentMatchInfo{
			AgentID:        result.SelectedAgent.AgentID,
			AgentName:      result.SelectedAgent.AgentName,
			AgentType:      result.SelectedAgent.AgentType,
			MatchScore:     result.SelectedAgent.MatchScore,
			MatchReason:    result.SelectedAgent.MatchReason,
			Confidence:     result.SelectedAgent.Confidence,
			Recommendation: result.SelectedAgent.Recommendation,
		},
		ExecutionTime: result.ExecutionTime,
	}

	// 转换所有评分
	response.AllScores = make([]AgentScoreInfo, len(result.AllScores))
	for i, score := range result.AllScores {
		response.AllScores[i] = AgentScoreInfo{
			AgentID:    score.AgentID,
			AgentName:  score.AgentName,
			TotalScore: score.TotalScore,
			ScoreBreakdown: ScoreBreakdown{
				LoadScore:        score.ScoreBreakdown.LoadScore,
				HealthScore:      score.ScoreBreakdown.HealthScore,
				CapabilityScore:  score.ScoreBreakdown.CapabilityScore,
				PerformanceScore: score.ScoreBreakdown.PerformanceScore,
				ContextScore:     score.ScoreBreakdown.ContextScore,
			},
		}
	}

	// 转换优化日志
	response.OptimizationLog = make([]OptimizationLog, len(result.OptimizationLog))
	for i, log := range result.OptimizationLog {
		response.OptimizationLog[i] = OptimizationLog{
			Step:    log.Step,
			Message: log.Message,
			Details: log.Details,
		}
	}

	c.JSON(http.StatusOK, OptimizeRoutingResponse{
		Code:    0,
		Message: "success",
		Data:    response,
	})
}

// GetOptimalAgentRequest 获取最优Agent请求
type GetOptimalAgentRequest struct {
	TenantID     string                  `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	Intent       string                  `json:"intent" binding:"required" vd:"len($) > 0"`
	EntityType   string                  `json:"entity_type" binding:"required"` // agent, workflow
	Preferences  map[string]interface{}  `json:"preferences" binding:"omitempty"`
	ExcludeIDs   []string                `json:"exclude_ids" binding:"omitempty"`
}

// GetOptimalAgentResponse 获取最优Agent响应
type GetOptimalAgentResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    *AgentMatchInfo   `json:"data,omitempty"`
}

// GetOptimalAgent 获取最优Agent
// @router /api/routing/optimal-agent [POST]
func (h *RoutingOptimizerHandler) GetOptimalAgent(ctx context.Context, c *app.RequestContext) {
	var req GetOptimalAgentRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层获取最优Agent
	result, err := h.optimizer.GetOptimalAgent(ctx, req.TenantID, req.Intent, req.EntityType, req.Preferences, req.ExcludeIDs)
	if err != nil {
		h.logger.Error("failed to get optimal agent",
			zap.String("tenant_id", req.TenantID),
			zap.String("intent", req.Intent),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, GetOptimalAgentResponse{
		Code:    0,
		Message: "success",
		Data: &AgentMatchInfo{
			AgentID:        result.AgentID,
			AgentName:      result.AgentName,
			AgentType:      result.AgentType,
			MatchScore:     result.MatchScore,
			MatchReason:    result.MatchReason,
			Confidence:     result.Confidence,
			Recommendation: result.Recommendation,
		},
	})
}

// PredictAgentLoadRequest 预测Agent负载请求
type PredictAgentLoadRequest struct {
	TenantID    string   `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	AgentIDs    []string `json:"agent_ids" binding:"required,min=1,max=50"`
	TimeHorizon int      `json:"time_horizon" binding:"omitempty,min=1,max=60"` // 预测时间范围（分钟）
}

// PredictAgentLoadResponse 预测Agent负载响应
type PredictAgentLoadResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    *LoadPredictionData  `json:"data,omitempty"`
}

// LoadPredictionData 负载预测数据
type LoadPredictionData struct {
	Predictions []AgentLoadPrediction `json:"predictions"`
	Timestamp   int64                 `json:"timestamp"`
	TimeHorizon int                   `json:"time_horizon"`
}

// AgentLoadPrediction Agent负载预测
type AgentLoadPrediction struct {
	AgentID          string           `json:"agent_id"`
	AgentName        string           `json:"agent_name"`
	CurrentLoad      float64          `json:"current_load"`
	PredictedLoad    float64          `json:"predicted_load"`
	LoadTrend        string           `json:"load_trend"` // increasing, stable, decreasing
	Confidence       float64          `json:"confidence"`
	Recommendations  []string         `json:"recommendations"`
	TimeSeriesData   []TimeSeriesDataPoint `json:"time_series_data"`
}

// TimeSeriesDataPoint 时间序列数据点
type TimeSeriesDataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// PredictAgentLoad 预测Agent负载
// @router /api/routing/predict-load [POST]
func (h *RoutingOptimizerHandler) PredictAgentLoad(ctx context.Context, c *app.RequestContext) {
	var req PredictAgentLoadRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认预测时间范围
	if req.TimeHorizon == 0 {
		req.TimeHorizon = 30 // 默认30分钟
	}

	// 调用Service层预测负载
	predictions, err := h.optimizer.PredictAgentLoad(ctx, req.TenantID, req.AgentIDs, req.TimeHorizon)
	if err != nil {
		h.logger.Error("failed to predict agent load",
			zap.String("tenant_id", req.TenantID),
			zap.Strings("agent_ids", req.AgentIDs),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	predList := make([]AgentLoadPrediction, len(predictions))
	for i, pred := range predictions {
		timeSeries := make([]TimeSeriesDataPoint, len(pred.TimeSeriesData))
		for j, ts := range pred.TimeSeriesData {
			timeSeries[j] = TimeSeriesDataPoint{
				Timestamp: ts.Timestamp,
				Value:     ts.Value,
			}
		}

		predList[i] = AgentLoadPrediction{
			AgentID:         pred.AgentID,
			AgentName:       pred.AgentName,
			CurrentLoad:     pred.CurrentLoad,
			PredictedLoad:   pred.PredictedLoad,
			LoadTrend:       pred.LoadTrend,
			Confidence:      pred.Confidence,
			Recommendations: pred.Recommendations,
			TimeSeriesData:  timeSeries,
		}
	}

	c.JSON(http.StatusOK, PredictAgentLoadResponse{
		Code:    0,
		Message: "success",
		Data: &LoadPredictionData{
			Predictions: predList,
			Timestamp:   predictions[0].Timestamp,
			TimeHorizon: req.TimeHorizon,
		},
	})
}

// SetOptimizerConfigRequest 设置优化器配置请求
type SetOptimizerConfigRequest struct {
	TenantID     string                 `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	ScoringWeights ScoringWeightsConfig  `json:"scoring_weights" binding:"required"`
	OptimizationLevel string            `json:"optimization_level" binding:"omitempty,oneof=fast balanced accurate"`
	EnableCache  bool                   `json:"enable_cache" binding:"omitempty"`
	CacheTTL     int                    `json:"cache_ttl" binding:"omitempty,min=1,max=3600"` // 秒
}

// ScoringWeightsConfig 评分权重配置
type ScoringWeightsConfig struct {
	LoadWeight        float64 `json:"load_weight" binding:"required,min=0,max=1"`
	HealthWeight      float64 `json:"health_weight" binding:"required,min=0,max=1"`
	CapabilityWeight  float64 `json:"capability_weight" binding:"required,min=0,max=1"`
	PerformanceWeight float64 `json:"performance_weight" binding:"required,min=0,max=1"`
	ContextWeight     float64 `json:"context_weight" binding:"required,min=0,max=1"`
}

// SetOptimizerConfigResponse 设置优化器配置响应
type SetOptimizerConfigResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SetOptimizerConfig 设置优化器配置
// @router /api/routing/optimizer/config [POST]
func (h *RoutingOptimizerHandler) SetOptimizerConfig(ctx context.Context, c *app.RequestContext) {
	var req SetOptimizerConfigRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 转换权重配置
	weights := routingapp.ScoringWeightsConfig{
		LoadWeight:        req.ScoringWeights.LoadWeight,
		HealthWeight:      req.ScoringWeights.HealthWeight,
		CapabilityWeight:  req.ScoringWeights.CapabilityWeight,
		PerformanceWeight: req.ScoringWeights.PerformanceWeight,
		ContextWeight:     req.ScoringWeights.ContextWeight,
	}

	// 调用Service层设置配置
	if err := h.optimizer.SetConfig(ctx, req.TenantID, weights, req.OptimizationLevel, req.EnableCache, req.CacheTTL); err != nil {
		h.logger.Error("failed to set optimizer config",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, SetOptimizerConfigResponse{
		Code:    0,
		Message: "success",
	})
}

// GetOptimizerConfigRequest 获取优化器配置请求
type GetOptimizerConfigRequest struct {
	TenantID string `form:"tenant_id" binding:"required" vd:"len($) > 0"`
}

// GetOptimizerConfigResponse 获取优化器配置响应
type GetOptimizerConfigResponse struct {
	Code    int                      `json:"code"`
	Message string                   `json:"message"`
	Data    *OptimizerConfigInfo     `json:"data,omitempty"`
}

// OptimizerConfigInfo 优化器配置信息
type OptimizerConfigInfo struct {
	TenantID          string                 `json:"tenant_id"`
	ScoringWeights    ScoringWeightsConfig   `json:"scoring_weights"`
	OptimizationLevel string                 `json:"optimization_level"`
	EnableCache       bool                   `json:"enable_cache"`
	CacheTTL          int                    `json:"cache_ttl"`
}

// GetOptimizerConfig 获取优化器配置
// @router /api/routing/optimizer/config [GET]
func (h *RoutingOptimizerHandler) GetOptimizerConfig(ctx context.Context, c *app.RequestContext) {
	var req GetOptimizerConfigRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层获取配置
	config, err := h.optimizer.GetConfig(ctx, req.TenantID)
	if err != nil {
		h.logger.Error("failed to get optimizer config",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, GetOptimizerConfigResponse{
		Code:    0,
		Message: "success",
		Data: &OptimizerConfigInfo{
			TenantID:          config.TenantID,
			ScoringWeights: ScoringWeightsConfig{
				LoadWeight:        config.ScoringWeights.LoadWeight,
				HealthWeight:      config.ScoringWeights.HealthWeight,
				CapabilityWeight:  config.ScoringWeights.CapabilityWeight,
				PerformanceWeight: config.ScoringWeights.PerformanceWeight,
				ContextWeight:     config.ScoringWeights.ContextWeight,
			},
			OptimizationLevel: config.OptimizationLevel,
			EnableCache:       config.EnableCache,
			CacheTTL:          config.CacheTTL,
		},
	})
}
