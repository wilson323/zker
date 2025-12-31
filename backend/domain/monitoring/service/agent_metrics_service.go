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
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
	"github.com/coze-dev/coze-studio/backend/domain/monitoring/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// AgentMetricsService Agent指标监控服务
type AgentMetricsService struct {
	metricsRepo   repository.AgentMetricsRepository
	prometheusSvc PrometheusClient
	logger        *zap.Logger
	cache         *agentMetricsCache
	mu            sync.RWMutex
}

// agentMetricsCache Agent指标缓存
type agentMetricsCache struct {
	data     map[string]*cachedAgentMetric
	expireAt time.Time
	mu       sync.RWMutex
}

type cachedAgentMetric struct {
	metric    interface{}
	timestamp time.Time
}

// PrometheusClient Prometheus客户端接口
type PrometheusClient interface {
	RecordMetric(metricName string, labels map[string]string, value float64) error
	RecordHistogram(metricName string, labels map[string]string, value float64) error
	RecordCounter(metricName string, labels map[string]string, value float64) error
}

// NewAgentMetricsService 创建Agent指标监控服务实例
func NewAgentMetricsService(
	metricsRepo repository.AgentMetricsRepository,
	prometheusSvc PrometheusClient,
	logger *zap.Logger,
) *AgentMetricsService {
	cache := &agentMetricsCache{
		data:     make(map[string]*cachedAgentMetric),
		expireAt: time.Now().Add(5 * time.Minute),
	}

	return &AgentMetricsService{
		metricsRepo:   metricsRepo,
		prometheusSvc: prometheusSvc,
		logger:        logger,
		cache:         cache,
	}
}

// CollectMetrics 收集Agent指标
func (s *AgentMetricsService) CollectMetrics(
	ctx context.Context,
	agentID string,
	timeRange *entity.TimeRange,
) (*AgentMetrics, error) {
	s.logger.Info("Collecting agent metrics",
		zap.String("agent_id", agentID),
		zap.Time("start", timeRange.StartTime),
		zap.Time("end", timeRange.EndTime),
	)

	// 并行收集各类指标
	var wg sync.WaitGroup
	var qpsMetrics *entity.AgentQPSMetrics
	var responseMetrics *entity.AgentResponseTimeMetrics
	var errorMetrics *entity.AgentErrorRateMetrics
	var satisfactionMetrics *entity.AgentSatisfactionMetrics
	var tokenMetrics *entity.AgentTokenUsageMetrics
	var qpsErr, responseErr, errorErr, satisfactionErr, tokenErr error

	wg.Add(5)

	go func() {
		defer wg.Done()
		qpsMetrics, qpsErr = s.metricsRepo.CalculateAgentQPSMetrics(ctx, agentID, timeRange.StartTime, timeRange.EndTime)
	}()

	go func() {
		defer wg.Done()
		responseMetrics, responseErr = s.metricsRepo.CalculateAgentResponseTimeMetrics(ctx, agentID, timeRange.StartTime, timeRange.EndTime)
	}()

	go func() {
		defer wg.Done()
		errorMetrics, errorErr = s.metricsRepo.CalculateAgentErrorRateMetrics(ctx, agentID, timeRange.StartTime, timeRange.EndTime)
	}()

	go func() {
		defer wg.Done()
		satisfactionMetrics, satisfactionErr = s.metricsRepo.CalculateAgentSatisfactionMetrics(ctx, agentID, timeRange.StartTime, timeRange.EndTime)
	}()

	go func() {
		defer wg.Done()
		tokenMetrics, tokenErr = s.metricsRepo.CalculateAgentTokenUsageMetrics(ctx, agentID, timeRange.StartTime, timeRange.EndTime)
	}()

	wg.Wait()

	if qpsErr != nil || responseErr != nil || errorErr != nil || satisfactionErr != nil || tokenErr != nil {
		return nil, errorx.New(errno.MetricsCollectionFailedCode,
			errorx.KV("qps_error", fmt.Sprintf("%v", qpsErr)),
			errorx.KV("response_error", fmt.Sprintf("%v", responseErr)),
			errorx.KV("error_error", fmt.Sprintf("%v", errorErr)),
			errorx.KV("satisfaction_error", fmt.Sprintf("%v", satisfactionErr)),
			errorx.KV("token_error", fmt.Sprintf("%v", tokenErr)))
	}

	metrics := &AgentMetrics{
		AgentID:            agentID,
		TimeRange:          *timeRange,
		QPSMetrics:         qpsMetrics,
		ResponseTimeMetrics: responseMetrics,
		ErrorRateMetrics:   errorMetrics,
		SatisfactionMetrics: satisfactionMetrics,
		TokenUsageMetrics:  tokenMetrics,
		CollectedAt:        time.Now(),
	}

	return metrics, nil
}

// GetMetrics 获取指标数据
func (s *AgentMetricsService) GetMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*MetricsResponse, error) {
	s.logger.Info("Getting agent metrics",
		zap.String("agent_id", req.AgentID),
		zap.String("metric_type", string(req.MetricType)),
	)

	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, errorx.WrapByCode(err, errno.InvalidParamsCode,
			errorx.KV("operation", "validate get metrics request"))
	}

	// 检查缓存
	cacheKey := fmt.Sprintf("metrics:%s:%s:%d:%d",
		req.AgentID, req.MetricType, req.TimeRange.StartTime.Unix(), req.TimeRange.EndTime.Unix())
	if cached := s.getFromCache(cacheKey); cached != nil {
		s.logger.Debug("Retrieved metrics from cache", zap.String("cache_key", cacheKey))
		return cached.(*MetricsResponse), nil
	}

	var response *MetricsResponse
	var err error

	// 根据指标类型获取数据
	switch req.MetricType {
	case entity.AgentMetricTypeQPS:
		response, err = s.getQPSMetrics(ctx, req)
	case entity.AgentMetricTypeResponseTime:
		response, err = s.getResponseTimeMetrics(ctx, req)
	case entity.AgentMetricTypeErrorRate:
		response, err = s.getErrorRateMetrics(ctx, req)
	case entity.AgentMetricTypeSatisfaction:
		response, err = s.getSatisfactionMetrics(ctx, req)
	case entity.AgentMetricTypeTokenUsage:
		response, err = s.getTokenUsageMetrics(ctx, req)
	default:
		return nil, errorx.New(errno.InvalidParamsCode,
			errorx.KV("metric_type", string(req.MetricType)),
			errorx.KV("reason", "unsupported metric type"))
	}

	if err != nil {
		return nil, err
	}

	// 缓存结果
	s.putToCache(cacheKey, response, 1*time.Minute)

	return response, nil
}

// CalculateMetrics 计算派生指标
func (s *AgentMetricsService) CalculateMetrics(
	ctx context.Context,
	metrics *RawMetrics,
) (*CalculatedMetrics, error) {
	s.logger.Info("Calculating derived metrics",
		zap.String("agent_id", metrics.AgentID),
	)

	calculated := &CalculatedMetrics{
		AgentID:       metrics.AgentID,
		CalculatedAt:  time.Now(),
		OverallScore:  0,
		PerformanceMetrics: &AgentPerformanceMetrics{
			Throughput:      metrics.Throughput,
			AvgResponseTime: metrics.ResponseTime,
			ErrorRate:       metrics.ErrorRate,
			Availability:    100 - metrics.ErrorRate*100,
		},
		QualityMetrics: &QualityMetrics{
			UserSatisfaction: metrics.UserSatisfaction,
			RelevanceScore:   metrics.RelevanceScore,
			AccuracyScore:    metrics.AccuracyScore,
		},
		CostMetrics: &CostMetrics{
			TokenUsage:      metrics.TokenUsage,
			APICallCount:    metrics.APICallCount,
			Cost:            metrics.Cost,
			CostPerRequest:  calculateCostPerRequest(metrics.Cost, metrics.APICallCount),
			CostPerToken:    calculateCostPerToken(metrics.Cost, metrics.TokenUsage),
		},
		BusinessMetrics: &BusinessMetrics{
			ActiveUsers:       metrics.ActiveUsers,
			ConversationCount: metrics.ConversationCount,
			ConversionRate:    metrics.ConversionRate,
		},
	}

	// 计算综合评分
	calculated.OverallScore = s.calculateOverallScore(calculated)

	// 生成优化建议
	calculated.Recommendations = s.generateRecommendations(calculated)

	return calculated, nil
}

// ExportMetrics 导出指标到Prometheus
func (s *AgentMetricsService) ExportMetrics(
	ctx context.Context,
	metrics *AgentMetrics,
) error {
	s.logger.Info("Exporting metrics to Prometheus",
		zap.String("agent_id", metrics.AgentID),
	)

	if s.prometheusSvc == nil {
		return errorx.New(errno.PrometheusNotAvailableCode,
			errorx.KV("reason", "Prometheus client not configured"))
	}

	labels := map[string]string{
		"agent_id": metrics.AgentID,
	}

	// 导出QPS指标
	if metrics.QPSMetrics != nil {
		if err := s.prometheusSvc.RecordMetric("zker_agent_qps", labels, metrics.QPSMetrics.Average); err != nil {
			s.logger.Error("Failed to export QPS metric", zap.Error(err))
			return errorx.WrapByCode(err, errno.MetricsExportFailedCode,
			errorx.KV("operation", "export metrics to prometheus"))
		}
		if err := s.prometheusSvc.RecordCounter("zker_agent_requests_total", labels, float64(metrics.QPSMetrics.Count)); err != nil {
			s.logger.Error("Failed to export requests total metric", zap.Error(err))
		}
	}

	// 导出响应时间指标
	if metrics.ResponseTimeMetrics != nil {
		if err := s.prometheusSvc.RecordHistogram("zker_agent_response_time_ms", labels, metrics.ResponseTimeMetrics.Average); err != nil {
			s.logger.Error("Failed to export response time metric", zap.Error(err))
			return errorx.WrapByCode(err, errno.MetricsExportFailedCode,
			errorx.KV("operation", "export metrics to prometheus"))
		}
	}

	// 导出错误率指标
	if metrics.ErrorRateMetrics != nil {
		if err := s.prometheusSvc.RecordMetric("zker_agent_error_rate", labels, metrics.ErrorRateMetrics.ErrorRate*100); err != nil {
			s.logger.Error("Failed to export error rate metric", zap.Error(err))
			return errorx.WrapByCode(err, errno.MetricsExportFailedCode,
			errorx.KV("operation", "export metrics to prometheus"))
		}
		if err := s.prometheusSvc.RecordCounter("zker_agent_errors_total", labels, float64(metrics.ErrorRateMetrics.ErrorRequests)); err != nil {
			s.logger.Error("Failed to export errors total metric", zap.Error(err))
		}
	}

	// 导出满意度指标
	if metrics.SatisfactionMetrics != nil {
		if err := s.prometheusSvc.RecordMetric("zker_agent_satisfaction_score", labels, metrics.SatisfactionMetrics.AverageScore); err != nil {
			s.logger.Error("Failed to export satisfaction metric", zap.Error(err))
			return errorx.WrapByCode(err, errno.MetricsExportFailedCode,
			errorx.KV("operation", "export metrics to prometheus"))
		}
	}

	// 导出Token使用指标
	if metrics.TokenUsageMetrics != nil {
		if err := s.prometheusSvc.RecordCounter("zker_agent_tokens_total", labels, float64(metrics.TokenUsageMetrics.TotalTokens)); err != nil {
			s.logger.Error("Failed to export tokens total metric", zap.Error(err))
		}
		if err := s.prometheusSvc.RecordMetric("zker_agent_cost", labels, metrics.TokenUsageMetrics.EstimatedCost); err != nil {
			s.logger.Error("Failed to export cost metric", zap.Error(err))
		}
	}

	s.logger.Info("Successfully exported metrics to Prometheus",
		zap.String("agent_id", metrics.AgentID),
	)

	return nil
}

// GetAgentHealthStatus 获取Agent健康状态
func (s *AgentMetricsService) GetAgentHealthStatus(
	ctx context.Context,
	tenantID, agentID string,
) (*entity.AgentHealthStatus, error) {
	s.logger.Info("Getting agent health status",
		zap.String("tenant_id", tenantID),
		zap.String("agent_id", agentID),
	)

	healthStatus, err := s.metricsRepo.GetAgentHealthStatus(ctx, tenantID, agentID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.AgentHealthCheckFailedCode,
			errorx.KV("operation", "get agent health status"))
	}

	return healthStatus, nil
}

// GetAgentPerformanceReport 获取Agent性能报告
func (s *AgentMetricsService) GetAgentPerformanceReport(
	ctx context.Context,
	agentID string,
	timeRange *entity.TimeRange,
) (*entity.AgentPerformanceReport, error) {
	s.logger.Info("Getting agent performance report",
		zap.String("agent_id", agentID),
	)

	report, err := s.metricsRepo.GetAgentPerformanceReport(ctx, agentID, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.PerformanceReportGenerationFailedCode,
			errorx.KV("operation", "get agent performance report"))
	}

	return report, nil
}

// ============================================================
// 私有方法
// ============================================================

// getQPSMetrics 获取QPS指标
func (s *AgentMetricsService) getQPSMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*MetricsResponse, error) {
	metrics, err := s.metricsRepo.CalculateAgentQPSMetrics(ctx, req.AgentID, req.TimeRange.StartTime, req.TimeRange.EndTime)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.MetricsQueryFailedCode,
			errorx.KV("operation", "calculate agent metrics"))
	}

	return &MetricsResponse{
		AgentID:     req.AgentID,
		MetricType:  req.MetricType,
		TimeRange:   *req.TimeRange,
		Data:        metrics,
		GeneratedAt: time.Now(),
	}, nil
}

// getResponseTimeMetrics 获取响应时间指标
func (s *AgentMetricsService) getResponseTimeMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*MetricsResponse, error) {
	metrics, err := s.metricsRepo.CalculateAgentResponseTimeMetrics(ctx, req.AgentID, req.TimeRange.StartTime, req.TimeRange.EndTime)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.MetricsQueryFailedCode,
			errorx.KV("operation", "calculate agent metrics"))
	}

	return &MetricsResponse{
		AgentID:     req.AgentID,
		MetricType:  req.MetricType,
		TimeRange:   *req.TimeRange,
		Data:        metrics,
		GeneratedAt: time.Now(),
	}, nil
}

// getErrorRateMetrics 获取错误率指标
func (s *AgentMetricsService) getErrorRateMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*MetricsResponse, error) {
	metrics, err := s.metricsRepo.CalculateAgentErrorRateMetrics(ctx, req.AgentID, req.TimeRange.StartTime, req.TimeRange.EndTime)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.MetricsQueryFailedCode,
			errorx.KV("operation", "calculate agent metrics"))
	}

	return &MetricsResponse{
		AgentID:     req.AgentID,
		MetricType:  req.MetricType,
		TimeRange:   *req.TimeRange,
		Data:        metrics,
		GeneratedAt: time.Now(),
	}, nil
}

// getSatisfactionMetrics 获取满意度指标
func (s *AgentMetricsService) getSatisfactionMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*MetricsResponse, error) {
	metrics, err := s.metricsRepo.CalculateAgentSatisfactionMetrics(ctx, req.AgentID, req.TimeRange.StartTime, req.TimeRange.EndTime)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.MetricsQueryFailedCode,
			errorx.KV("operation", "calculate agent metrics"))
	}

	return &MetricsResponse{
		AgentID:     req.AgentID,
		MetricType:  req.MetricType,
		TimeRange:   *req.TimeRange,
		Data:        metrics,
		GeneratedAt: time.Now(),
	}, nil
}

// getTokenUsageMetrics 获取Token使用指标
func (s *AgentMetricsService) getTokenUsageMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*MetricsResponse, error) {
	metrics, err := s.metricsRepo.CalculateAgentTokenUsageMetrics(ctx, req.AgentID, req.TimeRange.StartTime, req.TimeRange.EndTime)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.MetricsQueryFailedCode,
			errorx.KV("operation", "calculate agent metrics"))
	}

	return &MetricsResponse{
		AgentID:     req.AgentID,
		MetricType:  req.MetricType,
		TimeRange:   *req.TimeRange,
		Data:        metrics,
		GeneratedAt: time.Now(),
	}, nil
}

// calculateOverallScore 计算综合评分
func (s *AgentMetricsService) calculateOverallScore(metrics *CalculatedMetrics) float64 {
	performanceScore := 0.0
	qualityScore := 0.0
	costScore := 0.0
	businessScore := 0.0

	// 性能评分 (40%)
	p := metrics.PerformanceMetrics
	if p != nil {
		// 响应时间评分 (越低越好)
		rtScore := 100.0
		if p.AvgResponseTime < 100 {
			rtScore = 100
		} else if p.AvgResponseTime < 500 {
			rtScore = 80
		} else if p.AvgResponseTime < 1000 {
			rtScore = 60
		} else {
			rtScore = 30
		}

		// 错误率评分 (越低越好)
		erScore := 100.0 - p.ErrorRate*100

		// 可用性评分 (越高越好)
		avScore := p.Availability

		performanceScore = (rtScore + erScore + avScore) / 3.0
	}

	// 质量评分 (30%)
	q := metrics.QualityMetrics
	if q != nil {
		// 满意度评分 (1-5分，转换为0-100)
		satScore := q.UserSatisfaction * 20

		// 相关性和准确性评分
		relScore := q.RelevanceScore * 100
		accScore := q.AccuracyScore * 100

		qualityScore = (satScore + relScore + accScore) / 3.0
	}

	// 成本评分 (15%)
	c := metrics.CostMetrics
	if c != nil {
		// 成本评分 (越低越好)
		costScore = 100.0
		if c.CostPerRequest > 0.1 {
			costScore = 60
		} else if c.CostPerRequest > 0.05 {
			costScore = 80
		}
	}

	// 业务评分 (15%)
	b := metrics.BusinessMetrics
	if b != nil {
		// 转化率评分 (越高越好)
		businessScore = b.ConversionRate * 100
	}

	// 加权计算总分
	overallScore := performanceScore*0.4 + qualityScore*0.3 + costScore*0.15 + businessScore*0.15

	return overallScore
}

// generateRecommendations 生成优化建议
func (s *AgentMetricsService) generateRecommendations(metrics *CalculatedMetrics) []string {
	recommendations := []string{}

	p := metrics.PerformanceMetrics
	if p != nil {
		if p.AvgResponseTime > 1000 {
			recommendations = append(recommendations, "响应时间过长，建议优化提示词或使用更快的模型")
		}
		if p.ErrorRate > 0.05 {
			recommendations = append(recommendations, "错误率过高，建议检查错误日志并优化错误处理逻辑")
		}
		if p.Availability < 95 {
			recommendations = append(recommendations, "可用性较低，建议检查服务稳定性")
		}
	}

	q := metrics.QualityMetrics
	if q != nil {
		if q.UserSatisfaction < 3.0 {
			recommendations = append(recommendations, "用户满意度较低，建议优化回答质量")
		}
		if q.AccuracyScore < 0.7 {
			recommendations = append(recommendations, "准确性较低，建议增加知识库或优化提示词")
		}
	}

	c := metrics.CostMetrics
	if c != nil {
		if c.CostPerRequest > 0.1 {
			recommendations = append(recommendations, "单次请求成本较高，建议优化Token使用或切换更经济的模型")
		}
	}

	b := metrics.BusinessMetrics
	if b != nil {
		if b.ConversionRate < 0.3 {
			recommendations = append(recommendations, "转化率较低，建议优化对话流程和引导策略")
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Agent运行良好，继续保持")
	}

	return recommendations
}

// ============================================================
// 缓存方法
// ============================================================

func (s *AgentMetricsService) getFromCache(key string) interface{} {
	s.cache.mu.RLock()
	defer s.cache.mu.RUnlock()

	if time.Now().After(s.cache.expireAt) {
		return nil
	}

	cached, exists := s.cache.data[key]
	if !exists {
		return nil
	}

	if time.Since(cached.timestamp) > 5*time.Minute {
		return nil
	}

	return cached.metric
}

func (s *AgentMetricsService) putToCache(key string, metric interface{}, ttl time.Duration) {
	s.cache.mu.Lock()
	defer s.cache.mu.Unlock()

	s.cache.data[key] = &cachedAgentMetric{
		metric:    metric,
		timestamp: time.Now(),
	}
}

// ============================================================
// 辅助函数
// ============================================================

func calculateCostPerRequest(totalCost float64, apiCallCount int64) float64 {
	if apiCallCount == 0 {
		return 0
	}
	return totalCost / float64(apiCallCount)
}

func calculateCostPerToken(totalCost float64, tokenUsage int64) float64 {
	if tokenUsage == 0 {
		return 0
	}
	return totalCost / float64(tokenUsage)
}

// ============================================================
// 请求和响应类型
// ============================================================

// GetMetricsRequest 获取指标请求
type GetMetricsRequest struct {
	AgentID    string                      `json:"agent_id"`
	MetricType entity.AgentMetricType      `json:"metric_type"`
	TimeRange  *entity.TimeRange           `json:"time_range"`
	Aggregation string                     `json:"aggregation"` // avg, sum, min, max
}

// Validate 验证请求
func (r *GetMetricsRequest) Validate() error {
	if r.AgentID == "" {
		return fmt.Errorf("agent_id is required")
	}
	if err := entity.ValidateMetricType(r.MetricType); err != nil {
		return err
	}
	if r.TimeRange == nil {
		return fmt.Errorf("time_range is required")
	}
	return nil
}

// MetricsResponse 指标响应
type MetricsResponse struct {
	AgentID     string                 `json:"agent_id"`
	MetricType  entity.AgentMetricType `json:"metric_type"`
	TimeRange   entity.TimeRange       `json:"time_range"`
	Data        interface{}            `json:"data"`
	GeneratedAt time.Time              `json:"generated_at"`
}

// AgentMetrics Agent指标集合
type AgentMetrics struct {
	AgentID            string                         `json:"agent_id"`
	TimeRange          entity.TimeRange               `json:"time_range"`
	QPSMetrics         *entity.AgentQPSMetrics        `json:"qps_metrics,omitempty"`
	ResponseTimeMetrics *entity.AgentResponseTimeMetrics `json:"response_time_metrics,omitempty"`
	ErrorRateMetrics   *entity.AgentErrorRateMetrics  `json:"error_rate_metrics,omitempty"`
	SatisfactionMetrics *entity.AgentSatisfactionMetrics `json:"satisfaction_metrics,omitempty"`
	TokenUsageMetrics  *entity.AgentTokenUsageMetrics `json:"token_usage_metrics,omitempty"`
	CollectedAt        time.Time                      `json:"collected_at"`
}

// RawMetrics 原始指标
type RawMetrics struct {
	AgentID          string    `json:"agent_id"`
	ResponseTime     float64   `json:"response_time"`     // ms
	Throughput       float64   `json:"throughput"`        // QPS
	ErrorRate        float64   `json:"error_rate"`        // 0-1
	UserSatisfaction float64   `json:"user_satisfaction"` // 1-5
	RelevanceScore   float64   `json:"relevance_score"`   // 0-1
	AccuracyScore    float64   `json:"accuracy_score"`    // 0-1
	TokenUsage       int64     `json:"token_usage"`
	APICallCount     int64     `json:"api_call_count"`
	Cost             float64   `json:"cost"`
	ActiveUsers      int64     `json:"active_users"`
	ConversationCount int64    `json:"conversation_count"`
	ConversionRate   float64   `json:"conversion_rate"`
}

// CalculatedMetrics 计算后的指标
type CalculatedMetrics struct {
	AgentID             string                  `json:"agent_id"`
	OverallScore        float64                 `json:"overall_score"`        // 综合评分 0-100
	PerformanceMetrics  *AgentPerformanceMetrics `json:"performance_metrics"`
	QualityMetrics      *QualityMetrics         `json:"quality_metrics"`
	CostMetrics         *CostMetrics            `json:"cost_metrics"`
	BusinessMetrics     *BusinessMetrics        `json:"business_metrics"`
	Recommendations     []string                `json:"recommendations"`
	CalculatedAt        time.Time               `json:"calculated_at"`
}

// AgentPerformanceMetrics Agent性能指标（本地版本，与 performance_service 中的区分）
type AgentPerformanceMetrics struct {
	Throughput      float64 `json:"throughput"`       // QPS
	AvgResponseTime float64 `json:"avg_response_time"` // ms
	ErrorRate       float64 `json:"error_rate"`        // 0-1
	Availability    float64 `json:"availability"`      // 0-100
}

// QualityMetrics 质量指标
type QualityMetrics struct {
	UserSatisfaction float64 `json:"user_satisfaction"` // 1-5
	RelevanceScore   float64 `json:"relevance_score"`   // 0-1
	AccuracyScore    float64 `json:"accuracy_score"`    // 0-1
}

// CostMetrics 成本指标
type CostMetrics struct {
	TokenUsage     int64   `json:"token_usage"`
	APICallCount   int64   `json:"api_call_count"`
	Cost           float64 `json:"cost"`
	CostPerRequest float64 `json:"cost_per_request"`
	CostPerToken   float64 `json:"cost_per_token"`
}

// BusinessMetrics 业务指标
type BusinessMetrics struct {
	ActiveUsers       int64   `json:"active_users"`
	ConversationCount int64   `json:"conversation_count"`
	ConversionRate    float64 `json:"conversion_rate"` // 0-1
}
