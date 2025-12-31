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
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	routingservice "github.com/coze-dev/coze-studio/backend/domain/routing/service"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/util"
)

// RoutingMonitoringHandler 路由监控Handler
type RoutingMonitoringHandler struct {
	h              *coze.Handler
	abTestService  *routingservice.ABTestService
	learningService *routingservice.RoutingLearningService
	loadMonitor    *routingservice.LoadMonitor
	logger         *routingservice.AdvancedLoadBalancer
}

// NewRoutingMonitoringHandler 创建路由监控Handler
func NewRoutingMonitoringHandler(
	h *coze.Handler,
	abTestService *routingservice.ABTestService,
	learningService *routingservice.RoutingLearningService,
	loadMonitor *routingservice.LoadMonitor,
	advancedBalancer *routingservice.AdvancedLoadBalancer,
) *RoutingMonitoringHandler {
	return &RoutingMonitoringHandler{
		h:               h,
		abTestService:   abTestService,
		learningService: learningService,
		loadMonitor:     loadMonitor,
		logger:          advancedBalancer,
	}
}

// RoutingMetricsRequest 路由指标请求
type RoutingMetricsRequest struct {
	TenantID  string `json:"tenant_id" query:"tenant_id"`
	StartTime int64  `json:"start_time" query:"start_time"`
	EndTime   int64  `json:"end_time" query:"end_time"`
	Granularity string `json:"granularity" query:"granularity"` // hour/day/week
}

// RoutingMetricsResponse 路由指标响应
type RoutingMetricsResponse struct {
	TenantID       string                    `json:"tenant_id"`
	TimeRange      routingservice.TimeRange  `json:"time_range"`
	TotalRequests  int                       `json:"total_requests"`
	AvgResponseTime float64                  `json:"avg_response_time"`
	SuccessRate    float64                   `json:"success_rate"`
	ErrorRate      float64                   `json:"error_rate"`
	IntentDistribution map[string]int         `json:"intent_distribution"`
	AgentPerformance map[string]*routingservice.AgentPerformanceStats `json:"agent_performance"`
	Timestamp      int64                     `json:"timestamp"`
}

// RoutingMetricsHandler 获取路由指标
func (h *RoutingMonitoringHandler) RoutingMetricsHandler(ctx context.Context, c *app.RequestContext) {
	var req RoutingMetricsRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// 设置默认时间范围（最近1小时）
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}
	if req.StartTime == 0 {
		req.StartTime = req.EndTime - 3600
	}

	// 构建时间范围
	timeRange := routingservice.TimeRange{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// 从路由学习服务获取洞察
	insight, err := h.learningService.LearnFromRoutingLogs(ctx, req.TenantID, timeRange)
	if err != nil {
		util.BuildErrorResp(c, err)
		return
	}

	// 构建响应
	resp := &RoutingMetricsResponse{
		TenantID:          req.TenantID,
		TimeRange:         timeRange,
		TotalRequests:     insight.TotalRequests,
		AvgResponseTime:   insight.AvgResponseTime,
		SuccessRate:       insight.SuccessRate,
		IntentDistribution: insight.IntentDistribution,
		AgentPerformance:  insight.AgentPerformance,
		Timestamp:         time.Now().Unix(),
	}

	// 计算错误率
	if len(insight.AgentPerformance) > 0 {
		totalErrorRate := 0.0
		for _, perf := range insight.AgentPerformance {
			totalErrorRate += perf.ErrorRate
		}
		resp.ErrorRate = totalErrorRate / float64(len(insight.AgentPerformance))
	}

	util.BuildSuccessResp(c, resp)
}

// RoutingPerformanceRequest 路由性能请求
type RoutingPerformanceRequest struct {
	TenantID string `json:"tenant_id" query:"tenant_id"`
	AgentID  string `json:"agent_id" query:"agent_id"`
	StartTime int64 `json:"start_time" query:"start_time"`
	EndTime  int64 `json:"end_time" query:"end_time"`
}

// RoutingPerformanceResponse 路由性能响应
type RoutingPerformanceResponse struct {
	TenantID         string                 `json:"tenant_id"`
	AgentID          string                 `json:"agent_id"`
	TotalRequests    int                    `json:"total_requests"`
	AvgResponseTime  float64                `json:"avg_response_time"`
	P95ResponseTime  float64                `json:"p95_response_time"`
	P99ResponseTime  float64                `json:"p99_response_time"`
	SuccessRate      float64                `json:"success_rate"`
	ErrorRate        float64                `json:"error_rate"`
	Throughput       float64                `json:"throughput"`       // requests per second
	ResourceUsage    ResourceUsageMetrics   `json:"resource_usage"`
	Timestamp        int64                  `json:"timestamp"`
}

// ResourceUsageMetrics 资源使用指标
type ResourceUsageMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`     // 0-1
	MemoryUsage float64 `json:"memory_usage"`  // 0-1
	DiskUsage   float64 `json:"disk_usage"`    // 0-1
	NetworkIn   float64 `json:"network_in"`    // bytes/s
	NetworkOut  float64 `json:"network_out"`   // bytes/s
}

// RoutingPerformanceHandler 获取路由性能
func (h *RoutingMonitoringHandler) RoutingPerformanceHandler(ctx context.Context, c *app.RequestContext) {
	var req RoutingPerformanceRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// 设置默认时间范围
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}
	if req.StartTime == 0 {
		req.StartTime = req.EndTime - 3600
	}

	// 从负载监控器获取实时指标
	metrics := h.loadMonitor.GetRealTimeMetrics(req.AgentID)

	// 构建响应
	resp := &RoutingPerformanceResponse{
		TenantID: req.TenantID,
		AgentID:  req.AgentID,
		Timestamp: time.Now().Unix(),
	}

	if metrics != nil {
		resp.TotalRequests = metrics.TotalRequests
		resp.AvgResponseTime = metrics.AvgResponseTime
		resp.SuccessRate = metrics.SuccessRate
		resp.ErrorRate = metrics.ErrorRate
		resp.ResourceUsage = ResourceUsageMetrics{
			CPUUsage:    metrics.CPUUsage,
			MemoryUsage: metrics.MemoryUsage,
		}
	}

	// 计算吞吐量
	timeRange := req.EndTime - req.StartTime
	if timeRange > 0 && resp.TotalRequests > 0 {
		resp.Throughput = float64(resp.TotalRequests) / float64(timeRange)
	}

	util.BuildSuccessResp(c, resp)
}

// ABTestListRequest A/B测试列表请求
type ABTestListRequest struct {
	TenantID string `json:"tenant_id" query:"tenant_id"`
	Status   string `json:"status" query:"status"`
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
}

// ABTestListResponse A/B测试列表响应
type ABTestListResponse struct {
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	ABTests    []*routingservice.ABTest  `json:"abtests"`
}

// ABTestListHandler 获取A/B测试列表
func (h *RoutingMonitoringHandler) ABTestListHandler(ctx context.Context, c *app.RequestContext) {
	var req ABTestListRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// 设置默认分页
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// TODO: 从仓储获取A/B测试列表
	// 这里需要实现仓储接口的List方法

	// 临时返回空列表
	resp := &ABTestListResponse{
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
		ABTests:  []*routingservice.ABTest{},
	}

	util.BuildSuccessResp(c, resp)
}

// ABTestResultsResponse A/B测试结果响应
type ABTestResultsResponse struct {
	TestID              string                           `json:"test_id"`
	StrategyAStats      *routingservice.TestStrategyStats `json:"strategy_a_stats"`
	StrategyBStats      *routingservice.TestStrategyStats `json:"strategy_b_stats"`
	Winner              string                           `json:"winner"`
	Confidence          float64                          `json:"confidence"`
	IsStatisticallySignificant bool                      `json:"is_statistically_significant"`
	Recommendation      string                           `json:"recommendation"`
	Timestamp           int64                            `json:"timestamp"`
}

// ABTestResultsHandler 获取A/B测试结果
func (h *RoutingMonitoringHandler) ABTestResultsHandler(ctx context.Context, c *app.RequestContext) {
	testID := c.Param("experiment_id")
	if testID == "" {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("reason", "experiment_id is required")))
		return
	}

	// 获取测试结果
	result, err := h.abTestService.GetABTestResults(ctx, testID)
	if err != nil {
		util.BuildErrorResp(c, err)
		return
	}

	// 构建响应
	resp := &ABTestResultsResponse{
		TestID:              result.TestID,
		StrategyAStats:      result.StrategyAStats,
		StrategyBStats:      result.StrategyBStats,
		Winner:              result.Winner,
		Confidence:          result.Confidence,
		IsStatisticallySignificant: result.IsStatisticallySignificant,
		Recommendation:      result.Recommendation,
		Timestamp:           time.Now().Unix(),
	}

	util.BuildSuccessResp(c, resp)
}

// RoutingHealthRequest 路由健康检查请求
type RoutingHealthRequest struct {
	TenantID string `json:"tenant_id" query:"tenant_id"`
}

// RoutingHealthResponse 路由健康检查响应
type RoutingHealthResponse struct {
	OverallHealth    string                 `json:"overall_health"`    // healthy/degraded/unhealthy
	AgentHealth      map[string]string      `json:"agent_health"`      // agent_id -> health_status
	CircuitBreakers  map[string]interface{} `json:"circuit_breakers"`  // 熔断器状态
	LoadBalancers    map[string]interface{} `json:"load_balancers"`    // 负载均衡器状态
	Timestamp        int64                  `json:"timestamp"`
}

// RoutingHealthHandler 路由健康检查
func (h *RoutingMonitoringHandler) RoutingHealthHandler(ctx context.Context, c *app.RequestContext) {
	var req RoutingHealthRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// TODO: 从健康检查服务获取实际健康状态
	// 这里需要集成健康检查服务

	resp := &RoutingHealthResponse{
		OverallHealth:   "healthy",
		AgentHealth:     make(map[string]string),
		CircuitBreakers: make(map[string]interface{}),
		LoadBalancers:   make(map[string]interface{}),
		Timestamp:       time.Now().Unix(),
	}

	util.BuildSuccessResp(c, resp)
}

// RoutingAlertsRequest 路由告警请求
type RoutingAlertsRequest struct {
	TenantID    string `json:"tenant_id" query:"tenant_id"`
	AlertType   string `json:"alert_type" query:"alert_type"`   // error/performance/capacity
	StartTime   int64  `json:"start_time" query:"start_time"`
	EndTime     int64  `json:"end_time" query:"end_time"`
	Severity    string `json:"severity" query:"severity"`       // critical/warning/info
	IsResolved  *bool  `json:"is_resolved" query:"is_resolved"`
	Page        int    `json:"page" query:"page"`
	PageSize    int    `json:"page_size" query:"page_size"`
}

// RoutingAlertsResponse 路由告警响应
type RoutingAlertsResponse struct {
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	Alerts     []*RoutingAlert        `json:"alerts"`
}

// RoutingAlert 路由告警
type RoutingAlert struct {
	AlertID     string    `json:"alert_id"`
	TenantID    string    `json:"tenant_id"`
	AlertType   string    `json:"alert_type"`
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AffectedAgents []string `json:"affected_agents"`
	Metrics     map[string]interface{} `json:"metrics"`
	IsResolved  bool      `json:"is_resolved"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// RoutingAlertsHandler 获取路由告警
func (h *RoutingMonitoringHandler) RoutingAlertsHandler(ctx context.Context, c *app.RequestContext) {
	var req RoutingAlertsRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// 设置默认分页
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// TODO: 从告警服务获取告警列表
	resp := &RoutingAlertsResponse{
		Total:    0,
		Page:     req.Page,
		PageSize: req.PageSize,
		Alerts:   []*RoutingAlert{},
	}

	util.BuildSuccessResp(c, resp)
}

// RoutingStatisticsRequest 路由统计请求
type RoutingStatisticsRequest struct {
	TenantID  string `json:"tenant_id" query:"tenant_id"`
	StartTime int64  `json:"start_time" query:"start_time"`
	EndTime   int64  `json:"end_time" query:"end_time"`
	GroupBy   string `json:"group_by" query:"group_by"` // agent/intent/hour
}

// RoutingStatisticsResponse 路由统计响应
type RoutingStatisticsResponse struct {
	TenantID    string                        `json:"tenant_id"`
	TimeRange   routingservice.TimeRange      `json:"time_range"`
	GroupBy     string                        `json:"group_by"`
	Statistics  map[string]interface{}        `json:"statistics"`
	Timestamp   int64                         `json:"timestamp"`
}

// RoutingStatisticsHandler 获取路由统计
func (h *RoutingMonitoringHandler) RoutingStatisticsHandler(ctx context.Context, c *app.RequestContext) {
	var req RoutingStatisticsRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// 设置默认时间范围
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}
	if req.StartTime == 0 {
		req.StartTime = req.EndTime - 86400 // 默认最近24小时
	}

	timeRange := routingservice.TimeRange{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// 获取统计信息
	insight, err := h.learningService.LearnFromRoutingLogs(ctx, req.TenantID, timeRange)
	if err != nil {
		util.BuildErrorResp(c, err)
		return
	}

	// 构建统计响应
	statistics := make(map[string]interface{})
	statistics["total_requests"] = insight.TotalRequests
	statistics["avg_response_time"] = insight.AvgResponseTime
	statistics["success_rate"] = insight.SuccessRate
	statistics["intent_distribution"] = insight.IntentDistribution

	if req.GroupBy == "agent" {
		statistics["agent_performance"] = insight.AgentPerformance
	}

	resp := &RoutingStatisticsResponse{
		TenantID:  req.TenantID,
		TimeRange: timeRange,
		GroupBy:   req.GroupBy,
		Statistics: statistics,
		Timestamp: time.Now().Unix(),
	}

	util.BuildSuccessResp(c, resp)
}

// RoutingComparisonRequest 路由对比请求
type RoutingComparisonRequest struct {
	TenantID      string   `json:"tenant_id"`
	Agents        []string `json:"agents"`         // 要对比的Agent ID列表
	StartTime     int64    `json:"start_time"`
	EndTime       int64    `json:"end_time"`
	Metrics       []string `json:"metrics"`        // 要对比的指标：response_time/success_rate/error_rate
}

// RoutingComparisonResponse 路由对比响应
type RoutingComparisonResponse struct {
	TenantID    string                        `json:"tenant_id"`
	Agents      []string                      `json:"agents"`
	TimeRange   routingservice.TimeRange      `json:"time_range"`
	Comparison  map[string]map[string]float64 `json:"comparison"` // agent_id -> metric_name -> value
	Winner      map[string]string             `json:"winner"`       // metric_name -> agent_id
	Timestamp   int64                         `json:"timestamp"`
}

// RoutingComparisonHandler 路由对比
func (h *RoutingMonitoringHandler) RoutingComparisonHandler(ctx context.Context, c *app.RequestContext) {
	var req RoutingComparisonRequest
	if err := c.BindAndValidate(&req); err != nil {
		util.BuildErrorResp(c, errorx.New(errno.ErrRouteFormatInvalidCode,
			errorx.KV("error", err.Error())))
		return
	}

	// 设置默认时间范围
	if req.EndTime == 0 {
		req.EndTime = time.Now().Unix()
	}
	if req.StartTime == 0 {
		req.StartTime = req.EndTime - 3600
	}

	timeRange := routingservice.TimeRange{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// 获取所有Agent的性能数据
	insight, err := h.learningService.LearnFromRoutingLogs(ctx, req.TenantID, timeRange)
	if err != nil {
		util.BuildErrorResp(c, err)
		return
	}

	// 构建对比数据
	comparison := make(map[string]map[string]float64)
	winner := make(map[string]string)

	for _, agentID := range req.Agents {
		if perf, exists := insight.AgentPerformance[agentID]; exists {
			comparison[agentID] = map[string]float64{
				"response_time": perf.AvgResponseTime,
				"success_rate":  perf.SuccessRate,
				"error_rate":    perf.ErrorRate,
			}
		}
	}

	// 确定每个指标的获胜者
	for _, metric := range req.Metrics {
		var bestAgent string
		var bestValue float64

		for agentID, metrics := range comparison {
			value, exists := metrics[metric]
			if !exists {
				continue
			}

			// 对于响应时间和错误率，越小越好；对于成功率，越大越好
			if bestAgent == "" {
				bestAgent = agentID
				bestValue = value
			} else if metric == "success_rate" {
				if value > bestValue {
					bestAgent = agentID
					bestValue = value
				}
			} else {
				if value < bestValue {
					bestAgent = agentID
					bestValue = value
				}
			}
		}

		if bestAgent != "" {
			winner[metric] = bestAgent
		}
	}

	resp := &RoutingComparisonResponse{
		TenantID:   req.TenantID,
		Agents:     req.Agents,
		TimeRange:  timeRange,
		Comparison: comparison,
		Winner:     winner,
		Timestamp:  time.Now().Unix(),
	}

	util.BuildSuccessResp(c, resp)
}
