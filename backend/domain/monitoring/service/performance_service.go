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
	"math"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/monitoring/entity"
	"github.com/coze-dev/coze-studio/backend/domain/monitoring/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// PerformanceService 性能分析服务
type PerformanceService struct {
	metricsRepo  repository.AgentMetricsRepository
	alertService AlertService
	logger       *zap.Logger
}

// AlertService 告警服务接口
type AlertService interface {
	CheckThreshold(ctx context.Context, tenantID, agentID string, metrics *PerformanceMetrics) error
	CreateAlert(ctx context.Context, alert *PerformanceAlert) error
}

// NewPerformanceService 创建性能分析服务实例
func NewPerformanceService(
	metricsRepo repository.AgentMetricsRepository,
	alertService AlertService,
	logger *zap.Logger,
) *PerformanceService {
	return &PerformanceService{
		metricsRepo:  metricsRepo,
		alertService: alertService,
		logger:       logger,
	}
}

// GetPerformanceReport 获取性能报告
func (s *PerformanceService) GetPerformanceReport(
	ctx context.Context,
	agentID string,
	timeRange *entity.TimeRange,
) (*PerformanceReport, error) {
	s.logger.Info("Getting performance report",
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
		return nil, errorx.Wrapf(fmt.Errorf("qps=%v, response=%v, error=%v, satisfaction=%v, token=%v",
			qpsErr, responseErr, errorErr, satisfactionErr, tokenErr), errno.MetricsQueryFailed)
	}

	// 计算综合评分
	overallScore := s.calculateOverallPerformanceScore(qpsMetrics, responseMetrics, errorMetrics, satisfactionMetrics)

	// 生成性能等级
	performanceLevel := entity.GetPerformanceLevel(overallScore)

	// 生成优化建议
	recommendations := s.generatePerformanceRecommendations(qpsMetrics, responseMetrics, errorMetrics, satisfactionMetrics, tokenMetrics)

	// 获取历史对比数据
	previousPeriodTimeRange := &entity.TimeRange{
		StartTime: timeRange.StartTime.Add(-timeRange.Duration()),
		EndTime:   timeRange.StartTime,
	}

	// 生成对比数据
	comparisonData := s.generateComparisonData(ctx, agentID, timeRange, previousPeriodTimeRange)

	report := &PerformanceReport{
		AgentID:            agentID,
		ReportDate:         time.Now().Format("2006-01-02"),
		TimeRange:          *timeRange,
		QPSMetrics:         qpsMetrics,
		ResponseTimeMetrics: responseMetrics,
		ErrorRateMetrics:   errorMetrics,
		SatisfactionMetrics: satisfactionMetrics,
		TokenUsageMetrics:  tokenMetrics,
		OverallScore:       overallScore,
		PerformanceLevel:   performanceLevel,
		Recommendations:    recommendations,
		ComparedWithPeriod: previousPeriodTimeRange.StartTime.Format("2006-01-02") + " ~ " + previousPeriodTimeRange.EndTime.Format("2006-01-02"),
		Improvements:       comparisonData.Improvements,
		Concerns:           comparisonData.Concerns,
		GeneratedAt:        time.Now(),
	}

	return report, nil
}

// AnalyzePerformance 分析性能瓶颈
func (s *PerformanceService) AnalyzePerformance(
	ctx context.Context,
	agentID string,
) (*PerformanceAnalysis, error) {
	s.logger.Info("Analyzing performance bottlenecks",
		zap.String("agent_id", agentID),
	)

	// 获取最近24小时的指标
	timeRange := entity.NewTimeRangeFromDuration(24 * time.Hour)

	report, err := s.GetPerformanceReport(ctx, agentID, timeRange)
	if err != nil {
		return nil, err
	}

	analysis := &PerformanceAnalysis{
		AgentID:        agentID,
		AnalyzedAt:     time.Now(),
		TimeRange:      *timeRange,
		OverallScore:   report.OverallScore,
		PerformanceLevel: report.PerformanceLevel,
		Bottlenecks:    s.identifyBottlenecks(report),
		StrongPoints:   s.identifyStrongPoints(report),
		RiskFactors:    s.identifyRiskFactors(report),
		OptimizationPlan: s.createOptimizationPlan(report),
	}

	return analysis, nil
}

// GetPerformanceTrends 获取性能趋势
func (s *PerformanceService) GetPerformanceTrends(
	ctx context.Context,
	req *TrendsRequest,
) (*TrendsResponse, error) {
	s.logger.Info("Getting performance trends",
		zap.String("agent_id", req.AgentID),
		zap.String("metric_type", string(req.MetricType)),
		zap.Int("days", req.Days),
	)

	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, errorx.Wrapf(err, errno.InvalidRequest)
	}

	// 生成时间序列数据点
	interval := 24 * time.Hour
	trendData := make([]*TrendDataPoint, 0, req.Days)

	for i := req.Days - 1; i >= 0; i-- {
		endTime := time.Now().AddDate(0, 0, -i)
		startTime := endTime.Add(-interval)

		var value float64
		var err error

		// 根据指标类型获取数据
		switch req.MetricType {
		case entity.AgentMetricTypeQPS:
			metrics, mErr := s.metricsRepo.CalculateAgentQPSMetrics(ctx, req.AgentID, startTime, endTime)
			if mErr == nil && metrics != nil {
				value = metrics.Average
			}
			err = mErr
		case entity.AgentMetricTypeResponseTime:
			metrics, mErr := s.metricsRepo.CalculateAgentResponseTimeMetrics(ctx, req.AgentID, startTime, endTime)
			if mErr == nil && metrics != nil {
				value = metrics.Average
			}
			err = mErr
		case entity.AgentMetricTypeErrorRate:
			metrics, mErr := s.metricsRepo.CalculateAgentErrorRateMetrics(ctx, req.AgentID, startTime, endTime)
			if mErr == nil && metrics != nil {
				value = metrics.ErrorRate * 100
			}
			err = mErr
		case entity.AgentMetricTypeSatisfaction:
			metrics, mErr := s.metricsRepo.CalculateAgentSatisfactionMetrics(ctx, req.AgentID, startTime, endTime)
			if mErr == nil && metrics != nil {
				value = metrics.AverageScore
			}
			err = mErr
		case entity.AgentMetricTypeTokenUsage:
			metrics, mErr := s.metricsRepo.CalculateAgentTokenUsageMetrics(ctx, req.AgentID, startTime, endTime)
			if mErr == nil && metrics != nil {
				value = float64(metrics.TotalTokens)
			}
			err = mErr
		default:
			return nil, errorx.Wrapf(fmt.Errorf("unsupported metric type: %s", req.MetricType), errno.InvalidRequest)
		}

		if err != nil {
			s.logger.Warn("Failed to get metrics for trend",
				zap.String("agent_id", req.AgentID),
				zap.Time("date", startTime),
				zap.Error(err),
			)
		}

		trendData = append(trendData, &TrendDataPoint{
			Timestamp: startTime,
			Date:      startTime.Format("2006-01-02"),
			Value:     value,
		})
	}

	// 计算趋势统计
	trendStats := s.calculateTrendStatistics(trendData)

	response := &TrendsResponse{
		AgentID:     req.AgentID,
		MetricType:  req.MetricType,
		TimeRange:   *entity.NewTimeRangeFromDuration(time.Duration(req.Days) * 24 * time.Hour),
		DataPoints:  trendData,
		Statistics:  trendStats,
		GeneratedAt: time.Now(),
	}

	return response, nil
}

// ComparePerformance 性能对比
func (s *PerformanceService) ComparePerformance(
	ctx context.Context,
	tenantID string,
	agentIDs []string,
	timeRange *entity.TimeRange,
) (*ComparisonReport, error) {
	s.logger.Info("Comparing agent performance",
		zap.String("tenant_id", tenantID),
		zap.Int("agent_count", len(agentIDs)),
	)

	if len(agentIDs) == 0 {
		return nil, errorx.Wrapf(fmt.Errorf("agent_ids is required"), errno.InvalidRequest)
	}

	// 并行收集所有Agent的指标
	var wg sync.WaitGroup
	metricsMap := make(map[string]*entity.AgentPerformanceReport)
	metricsMapMu := sync.Mutex{}
	errorsMap := make(map[string]error)
	errorsMapMu := sync.Mutex{}

	for _, agentID := range agentIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			report, err := s.metricsRepo.GetAgentPerformanceReport(ctx, id, timeRange.StartTime, timeRange.EndTime)
			if err != nil {
				errorsMapMu.Lock()
				errorsMap[id] = err
				errorsMapMu.Unlock()
				return
			}

			metricsMapMu.Lock()
			metricsMap[id] = report
			metricsMapMu.Unlock()
		}(agentID)
	}

	wg.Wait()

	if len(metricsMap) == 0 {
		return nil, errorx.Wrapf(fmt.Errorf("no valid metrics found"), errno.MetricsQueryFailed)
	}

	// 生成对比数据
	comparisons := s.generateComparisons(metricsMap)

	// 计算排名
	rankings := s.calculateRankings(comparisons)

	// 生成对比报告
	report := &ComparisonReport{
		TenantID:       tenantID,
		AgentIDs:       agentIDs,
		TimeRange:      *timeRange,
		Comparisons:    comparisons,
		Rankings:       rankings,
		GeneratedAt:    time.Now(),
	}

	return report, nil
}

// GetTopPerformingAgents 获取表现最佳的Agent
func (s *PerformanceService) GetTopPerformingAgents(
	ctx context.Context,
	tenantID string,
	limit int,
	timeRange *entity.TimeRange,
) ([]*TopAgent, error) {
	s.logger.Info("Getting top performing agents",
		zap.String("tenant_id", tenantID),
		zap.Int("limit", limit),
	)

	comparisons, err := s.metricsRepo.CompareAgents(ctx, tenantID, []string{}, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.MetricsQueryFailed)
	}

	// 按综合评分排序
	sort.Slice(comparisons, func(i, j int) bool {
		return s.calculateComparisonScore(comparisons[i]) > s.calculateComparisonScore(comparisons[j])
	})

	// 取前N个
	topAgents := make([]*TopAgent, 0, limit)
	for i := 0; i < min(limit, len(comparisons)); i++ {
		c := comparisons[i]
		topAgents = append(topAgents, &TopAgent{
			AgentID:       c.AgentID,
			AgentName:     c.AgentName,
			OverallRank:   c.OverallRanking,
			OverallScore:  s.calculateComparisonScore(c),
			Highlights:    s.getAgentHighlights(c),
		})
	}

	return topAgents, nil
}

// ============================================================
// 私有方法
// ============================================================

// calculateOverallPerformanceScore 计算综合性能评分
func (s *PerformanceService) calculateOverallPerformanceScore(
	qps *entity.AgentQPSMetrics,
	response *entity.AgentResponseTimeMetrics,
	error *entity.AgentErrorRateMetrics,
	satisfaction *entity.AgentSatisfactionMetrics,
) float64 {
	score := 0.0
	weightSum := 0.0

	// QPS评分 (权重: 25%)
	if qps != nil {
		qpsScore := 60.0
		if qps.Average > 100 {
			qpsScore = 90
		} else if qps.Average > 50 {
			qpsScore = 75
		} else if qps.Average > 10 {
			qpsScore = 60
		} else {
			qpsScore = 40
		}
		score += qpsScore * 0.25
		weightSum += 0.25
	}

	// 响应时间评分 (权重: 30%)
	if response != nil {
		rtScore := 60.0
		if response.Average < 100 {
			rtScore = 90
		} else if response.Average < 500 {
			rtScore = 75
		} else if response.Average < 1000 {
			rtScore = 60
		} else {
			rtScore = 30
		}
		score += rtScore * 0.30
		weightSum += 0.30
	}

	// 错误率评分 (权重: 25%)
	if error != nil {
		erScore := 100.0 - error.ErrorRate*100
		score += erScore * 0.25
		weightSum += 0.25
	}

	// 满意度评分 (权重: 20%)
	if satisfaction != nil {
		satScore := satisfaction.AverageScore * 20 // 1-5分转换为0-100
		score += satScore * 0.20
		weightSum += 0.20
	}

	if weightSum > 0 {
		return score / weightSum
	}
	return 0
}

// generatePerformanceRecommendations 生成性能建议
func (s *PerformanceService) generatePerformanceRecommendations(
	qps *entity.AgentQPSMetrics,
	response *entity.AgentResponseTimeMetrics,
	error *entity.AgentErrorRateMetrics,
	satisfaction *entity.AgentSatisfactionMetrics,
	token *entity.AgentTokenUsageMetrics,
) []string {
	recommendations := []string{}

	// QPS建议
	if qps != nil {
		if qps.Average < 10 {
			recommendations = append(recommendations, "QPS较低，可以考虑增加推广或优化Agent功能")
		}
		if qps.Trend == "down" {
			recommendations = append(recommendations, "QPS呈下降趋势，建议分析用户行为变化")
		}
	}

	// 响应时间建议
	if response != nil {
		if response.Average > 1000 {
			recommendations = append(recommendations, "响应时间过长(>1s)，建议优化提示词或使用更快的模型")
		} else if response.Average > 500 {
			recommendations = append(recommendations, "响应时间偏长，建议检查知识库检索或工具调用耗时")
		}
		if response.P95 > response.Average*2 {
			recommendations = append(recommendations, "P95响应时间远高于平均值，存在长尾延迟问题")
		}
	}

	// 错误率建议
	if error != nil {
		if error.ErrorRate > 0.05 {
			recommendations = append(recommendations, "错误率过高(>5%)，建议检查错误日志并优化错误处理")
		} else if error.ErrorRate > 0.01 {
			recommendations = append(recommendations, "错误率偏高(>1%)，建议分析错误类型并针对性优化")
		}
		if len(error.ErrorBreakdown) > 0 {
			maxErrorType := ""
			maxCount := int64(0)
			for errType, count := range error.ErrorBreakdown {
				if count > maxCount {
					maxCount = count
					maxErrorType = errType
				}
			}
			if maxErrorType != "" {
				recommendations = append(recommendations, fmt.Sprintf("最常见错误类型: %s，建议优先解决", maxErrorType))
			}
		}
	}

	// 满意度建议
	if satisfaction != nil {
		if satisfaction.AverageScore < 3.0 {
			recommendations = append(recommendations, "用户满意度较低(<3分)，建议优化回答质量和用户体验")
		} else if satisfaction.AverageScore < 4.0 {
			recommendations = append(recommendations, "用户满意度中等，建议收集反馈并持续改进")
		}
		oneStarRate := float64(satisfaction.OneStarCount) / float64(satisfaction.TotalRatings) * 100
		if oneStarRate > 20 {
			recommendations = append(recommendations, "1星评分占比过高，建议分析差评原因")
		}
	}

	// Token使用建议
	if token != nil {
		if token.EstimatedCost > 100 {
			recommendations = append(recommendations, "成本较高，建议优化提示词以减少Token使用")
		}
		if token.AveragePerReq > 2000 {
			recommendations = append(recommendations, "平均Token使用量较高，建议精简提示词或优化上下文")
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "性能表现良好，继续保持")
	}

	return recommendations
}

// identifyBottlenecks 识别性能瓶颈
func (s *PerformanceService) identifyBottlenecks(report *PerformanceReport) []string {
	bottlenecks := []string{}

	if report.ResponseTimeMetrics != nil {
		if report.ResponseTimeMetrics.Average > 1000 {
			bottlenecks = append(bottlenecks, "响应时间: 严重瓶颈(>1s)")
		} else if report.ResponseTimeMetrics.Average > 500 {
			bottlenecks = append(bottlenecks, "响应时间: 中度瓶颈(>500ms)")
		}
	}

	if report.ErrorRateMetrics != nil {
		if report.ErrorRateMetrics.ErrorRate > 0.05 {
			bottlenecks = append(bottlenecks, "错误率: 严重瓶颈(>5%)")
		} else if report.ErrorRateMetrics.ErrorRate > 0.01 {
			bottlenecks = append(bottlenecks, "错误率: 轻度瓶颈(>1%)")
		}
	}

	if report.SatisfactionMetrics != nil {
		if report.SatisfactionMetrics.AverageScore < 3.0 {
			bottlenecks = append(bottlenecks, "用户满意度: 严重问题(<3分)")
		}
	}

	if len(bottlenecks) == 0 {
		bottlenecks = append(bottlenecks, "无明显瓶颈")
	}

	return bottlenecks
}

// identifyStrongPoints 识别优势
func (s *PerformanceService) identifyStrongPoints(report *PerformanceReport) []string {
	strongPoints := []string{}

	if report.QPSMetrics != nil && report.QPSMetrics.Average > 100 {
		strongPoints = append(strongPoints, "高吞吐量")
	}

	if report.ResponseTimeMetrics != nil && report.ResponseTimeMetrics.Average < 200 {
		strongPoints = append(strongPoints, "快速响应")
	}

	if report.ErrorRateMetrics != nil && report.ErrorRateMetrics.ErrorRate < 0.005 {
		strongPoints = append(strongPoints, "低错误率")
	}

	if report.SatisfactionMetrics != nil && report.SatisfactionMetrics.AverageScore >= 4.5 {
		strongPoints = append(strongPoints, "高用户满意度")
	}

	if len(strongPoints) == 0 {
		strongPoints = append(strongPoints, "暂无明显优势")
	}

	return strongPoints
}

// identifyRiskFactors 识别风险因素
func (s *PerformanceService) identifyRiskFactors(report *PerformanceReport) []string {
	riskFactors := []string{}

	if report.ResponseTimeMetrics != nil {
		if report.ResponseTimeMetrics.P95 > report.ResponseTimeMetrics.Average*3 {
			riskFactors = append(riskFactors, "长尾延迟风险")
		}
		if report.ResponseTimeMetrics.Trend == "up" && report.ResponseTimeMetrics.TrendChange > 20 {
			riskFactors = append(riskFactors, "响应时间持续恶化")
		}
	}

	if report.ErrorRateMetrics != nil {
		if report.ErrorRateMetrics.Trend == "up" && report.ErrorRateMetrics.TrendChange > 20 {
			riskFactors = append(riskFactors, "错误率持续上升")
		}
	}

	if report.SatisfactionMetrics != nil {
		if report.SatisfactionMetrics.Trend == "down" && report.SatisfactionMetrics.TrendChange < -20 {
			riskFactors = append(riskFactors, "用户满意度持续下降")
		}
	}

	if len(riskFactors) == 0 {
		riskFactors = append(riskFactors, "无明显风险")
	}

	return riskFactors
}

// createOptimizationPlan 创建优化计划
func (s *PerformanceService) createOptimizationPlan(report *PerformanceReport) *OptimizationPlan {
	plan := &OptimizationPlan{
		PriorityActions:   []string{},
		ImprovementActions: []string{},
		LongTermActions:    []string{},
	}

	// 高优先级操作
	if report.ResponseTimeMetrics != nil && report.ResponseTimeMetrics.Average > 1000 {
		plan.PriorityActions = append(plan.PriorityActions, "立即优化响应时间：检查模型选择和提示词复杂度")
	}
	if report.ErrorRateMetrics != nil && report.ErrorRateMetrics.ErrorRate > 0.05 {
		plan.PriorityActions = append(plan.PriorityActions, "立即降低错误率：检查日志并修复高频错误")
	}

	// 改进操作
	if report.ResponseTimeMetrics != nil && report.ResponseTimeMetrics.Average > 500 {
		plan.ImprovementActions = append(plan.ImprovementActions, "优化知识库检索速度")
	}
	if report.SatisfactionMetrics != nil && report.SatisfactionMetrics.AverageScore < 4.0 {
		plan.ImprovementActions = append(plan.ImprovementActions, "收集用户反馈并优化回答质量")
	}

	// 长期操作
	plan.LongTermActions = append(plan.LongTermActions, "建立性能监控和告警机制")
	plan.LongTermActions = append(plan.LongTermActions, "定期进行性能评估和优化")

	return plan
}

// generateComparisonData 生成对比数据
func (s *PerformanceService) generateComparisonData(
	ctx context.Context,
	agentID string,
	currentRange *entity.TimeRange,
	previousRange *entity.TimeRange,
) *ComparisonData {
	// 获取当前周期数据
	currentQPS, _ := s.metricsRepo.CalculateAgentQPSMetrics(ctx, agentID, currentRange.StartTime, currentRange.EndTime)
	currentRT, _ := s.metricsRepo.CalculateAgentResponseTimeMetrics(ctx, agentID, currentRange.StartTime, currentRange.EndTime)
	currentError, _ := s.metricsRepo.CalculateAgentErrorRateMetrics(ctx, agentID, currentRange.StartTime, currentRange.EndTime)

	// 获取上一周期数据
	previousQPS, _ := s.metricsRepo.CalculateAgentQPSMetrics(ctx, agentID, previousRange.StartTime, previousRange.EndTime)
	previousRT, _ := s.metricsRepo.CalculateAgentResponseTimeMetrics(ctx, agentID, previousRange.StartTime, previousRange.EndTime)
	previousError, _ := s.metricsRepo.CalculateAgentErrorRateMetrics(ctx, agentID, previousRange.StartTime, previousRange.EndTime)

	improvements := []string{}
	concerns := []string{}

	// 对比QPS
	if currentQPS != nil && previousQPS != nil {
		if currentQPS.Average > previousQPS.Average*1.1 {
			improvements = append(improvements, fmt.Sprintf("QPS提升 %.1f%%", (currentQPS.Average/previousQPS.Average-1)*100))
		} else if currentQPS.Average < previousQPS.Average*0.9 {
			concerns = append(concerns, fmt.Sprintf("QPS下降 %.1f%%", (1-currentQPS.Average/previousQPS.Average)*100))
		}
	}

	// 对比响应时间
	if currentRT != nil && previousRT != nil {
		if currentRT.Average < previousRT.Average*0.9 {
			improvements = append(improvements, fmt.Sprintf("响应时间改善 %.1f%%", (1-currentRT.Average/previousRT.Average)*100))
		} else if currentRT.Average > previousRT.Average*1.1 {
			concerns = append(concerns, fmt.Sprintf("响应时间恶化 %.1f%%", (currentRT.Average/previousRT.Average-1)*100))
		}
	}

	// 对比错误率
	if currentError != nil && previousError != nil {
		if currentError.ErrorRate < previousError.ErrorRate*0.9 {
			improvements = append(improvements, fmt.Sprintf("错误率降低 %.1f%%", (1-currentError.ErrorRate/previousError.ErrorRate)*100))
		} else if currentError.ErrorRate > previousError.ErrorRate*1.1 {
			concerns = append(concerns, fmt.Sprintf("错误率上升 %.1f%%", (currentError.ErrorRate/previousError.ErrorRate-1)*100))
		}
	}

	return &ComparisonData{
		Improvements: improvements,
		Concerns:     concerns,
	}
}

// generateComparisons 生成对比数据
func (s *PerformanceService) generateComparisons(
	reports map[string]*entity.AgentPerformanceReport,
) []*entity.AgentComparison {
	comparisons := make([]*entity.AgentComparison, 0, len(reports))

	for agentID, report := range reports {
		c := &entity.AgentComparison{
			AgentID:     agentID,
			AgentName:   report.AgentName,
			QPSRanking:  0,
			ResponseTimeRanking: 0,
			ErrorRateRanking: 0,
			SatisfactionRanking: 0,
			OverallRanking: 0,
		}
		comparisons = append(comparisons, c)
	}

	return comparisons
}

// calculateRankings 计算排名
func (s *PerformanceService) calculateRankings(
	comparisons []*entity.AgentComparison,
) *RankingData {
	// 按各项指标排序并计算排名
	qpsRanking := make([]string, len(comparisons))
	rtRanking := make([]string, len(comparisons))
	errorRanking := make([]string, len(comparisons))
	satRanking := make([]string, len(comparisons))

	// QPS排名 (降序)
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].QPSPercentile > comparisons[j].QPSPercentile
	})
	for i, c := range comparisons {
		c.QPSRanking = i + 1
		qpsRanking[i] = c.AgentID
	}

	// 响应时间排名 (升序，percentile越小越好)
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].ResponseTimePercentile < comparisons[j].ResponseTimePercentile
	})
	for i, c := range comparisons {
		c.ResponseTimeRanking = i + 1
		rtRanking[i] = c.AgentID
	}

	// 错误率排名 (升序)
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].ErrorRatePercentile < comparisons[j].ErrorRatePercentile
	})
	for i, c := range comparisons {
		c.ErrorRateRanking = i + 1
		errorRanking[i] = c.AgentID
	}

	// 满意度排名 (降序)
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].SatisfactionPercentile > comparisons[j].SatisfactionPercentile
	})
	for i, c := range comparisons {
		c.SatisfactionRanking = i + 1
		satRanking[i] = c.AgentID
	}

	// 综合排名
	sort.Slice(comparisons, func(i, j int) bool {
		return s.calculateComparisonScore(comparisons[i]) > s.calculateComparisonScore(comparisons[j])
	})
	for i, c := range comparisons {
		c.OverallRanking = i + 1
	}

	return &RankingData{
		QPSRanking:         qpsRanking,
		ResponseTimeRanking: rtRanking,
		ErrorRateRanking:   errorRanking,
		SatisfactionRanking: satRanking,
	}
}

// calculateComparisonScore 计算对比评分
func (s *PerformanceService) calculateComparisonScore(c *entity.AgentComparison) float64 {
	// 简化的综合评分计算
	qpsScore := 100 - float64(c.QPSRanking)*5
	rtScore := 100 - float64(c.ResponseTimeRanking)*5
	errorScore := 100 - float64(c.ErrorRateRanking)*5
	satScore := 100 - float64(c.SatisfactionRanking)*5

	return (qpsScore + rtScore + errorScore + satScore) / 4.0
}

// calculateTrendStatistics 计算趋势统计
func (s *PerformanceService) calculateTrendStatistics(data []*TrendDataPoint) *TrendStatistics {
	if len(data) == 0 {
		return &TrendStatistics{}
	}

	values := make([]float64, len(data))
	for i, d := range data {
		values[i] = d.Value
	}

	// 计算基本统计
	sum := 0.0
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	avg := sum / float64(len(values))

	// 计算标准差
	variance := 0.0
	for _, v := range values {
		diff := v - avg
		variance += diff * diff
	}
	variance /= float64(len(values))
	stdDev := math.Sqrt(variance)

	// 计算趋势
	trend := "stable"
	trendChange := 0.0
	if len(data) >= 2 {
		firstHalf := values[:len(values)/2]
		secondHalf := values[len(values)/2:]

		firstAvg := average(firstHalf)
		secondAvg := average(secondHalf)

		if secondAvg > 0 {
			trendChange = ((secondAvg - firstAvg) / firstAvg) * 100
		}

		if trendChange > 10 {
			trend = "up"
		} else if trendChange < -10 {
			trend = "down"
		}
	}

	return &TrendStatistics{
		Average:       avg,
		Min:           min,
		Max:           max,
		StdDev:        stdDev,
		Trend:         trend,
		TrendChange:   trendChange,
		DataPointCount: len(data),
	}
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// getAgentHighlights 获取Agent亮点
func (s *PerformanceService) getAgentHighlights(c *entity.AgentComparison) []string {
	highlights := []string{}

	if c.QPSRanking <= 3 {
		highlights = append(highlights, "QPS排名前列")
	}
	if c.ResponseTimeRanking <= 3 {
		highlights = append(highlights, "响应时间优秀")
	}
	if c.ErrorRateRanking <= 3 {
		highlights = append(highlights, "错误率低")
	}
	if c.SatisfactionRanking <= 3 {
		highlights = append(highlights, "用户满意度高")
	}

	if len(highlights) == 0 {
		highlights = append(highlights, "表现稳定")
	}

	return highlights
}

// ============================================================
// 数据类型定义
// ============================================================

// PerformanceReport 性能报告
type PerformanceReport struct {
	AgentID             string                         `json:"agent_id"`
	ReportDate          string                         `json:"report_date"`
	TimeRange           entity.TimeRange               `json:"time_range"`
	QPSMetrics          *entity.AgentQPSMetrics        `json:"qps_metrics,omitempty"`
	ResponseTimeMetrics *entity.AgentResponseTimeMetrics `json:"response_time_metrics,omitempty"`
	ErrorRateMetrics    *entity.AgentErrorRateMetrics  `json:"error_rate_metrics,omitempty"`
	SatisfactionMetrics *entity.AgentSatisfactionMetrics `json:"satisfaction_metrics,omitempty"`
	TokenUsageMetrics   *entity.AgentTokenUsageMetrics `json:"token_usage_metrics,omitempty"`
	OverallScore        float64                        `json:"overall_score"`
	PerformanceLevel    string                         `json:"performance_level"`
	Recommendations     []string                       `json:"recommendations"`
	ComparedWithPeriod  string                         `json:"compared_with_period"`
	Improvements        []string                       `json:"improvements"`
	Concerns            []string                       `json:"concerns"`
	GeneratedAt         time.Time                      `json:"generated_at"`
}

// PerformanceAnalysis 性能分析
type PerformanceAnalysis struct {
	AgentID           string             `json:"agent_id"`
	AnalyzedAt        time.Time          `json:"analyzed_at"`
	TimeRange         entity.TimeRange   `json:"time_range"`
	OverallScore      float64            `json:"overall_score"`
	PerformanceLevel  string             `json:"performance_level"`
	Bottlenecks       []string           `json:"bottlenecks"`
	StrongPoints      []string           `json:"strong_points"`
	RiskFactors       []string           `json:"risk_factors"`
	OptimizationPlan  *OptimizationPlan  `json:"optimization_plan"`
}

// OptimizationPlan 优化计划
type OptimizationPlan struct {
	PriorityActions    []string `json:"priority_actions"`
	ImprovementActions []string `json:"improvement_actions"`
	LongTermActions    []string `json:"long_term_actions"`
}

// ComparisonData 对比数据
type ComparisonData struct {
	Improvements []string `json:"improvements"`
	Concerns     []string `json:"concerns"`
}

// TrendsRequest 趋势查询请求
type TrendsRequest struct {
	AgentID    string                      `json:"agent_id"`
	MetricType entity.AgentMetricType      `json:"metric_type"`
	Days       int                         `json:"days"`
}

// Validate 验证请求
func (r *TrendsRequest) Validate() error {
	if r.AgentID == "" {
		return fmt.Errorf("agent_id is required")
	}
	if err := entity.ValidateMetricType(r.MetricType); err != nil {
		return err
	}
	if r.Days < 1 || r.Days > 90 {
		return fmt.Errorf("days must be between 1 and 90")
	}
	return nil
}

// TrendsResponse 趋势查询响应
type TrendsResponse struct {
	AgentID     string                      `json:"agent_id"`
	MetricType  entity.AgentMetricType      `json:"metric_type"`
	TimeRange   entity.TimeRange            `json:"time_range"`
	DataPoints  []*TrendDataPoint           `json:"data_points"`
	Statistics  *TrendStatistics            `json:"statistics"`
	GeneratedAt time.Time                   `json:"generated_at"`
}

// TrendDataPoint 趋势数据点
type TrendDataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Date      string    `json:"date"`
	Value     float64   `json:"value"`
}

// TrendStatistics 趋势统计
type TrendStatistics struct {
	Average        float64 `json:"average"`
	Min            float64 `json:"min"`
	Max            float64 `json:"max"`
	StdDev         float64 `json:"std_dev"`
	Trend          string  `json:"trend"`          // up, down, stable
	TrendChange    float64 `json:"trend_change"`   // 百分比变化
	DataPointCount int     `json:"data_point_count"`
}

// ComparisonReport 对比报告
type ComparisonReport struct {
	TenantID    string                     `json:"tenant_id"`
	AgentIDs    []string                   `json:"agent_ids"`
	TimeRange   entity.TimeRange           `json:"time_range"`
	Comparisons []*entity.AgentComparison  `json:"comparisons"`
	Rankings    *RankingData               `json:"rankings"`
	GeneratedAt time.Time                  `json:"generated_at"`
}

// RankingData 排名数据
type RankingData struct {
	QPSRanking          []string `json:"qps_ranking"`
	ResponseTimeRanking []string `json:"response_time_ranking"`
	ErrorRateRanking    []string `json:"error_rate_ranking"`
	SatisfactionRanking []string `json:"satisfaction_ranking"`
}

// TopAgent Top Agent
type TopAgent struct {
	AgentID      string   `json:"agent_id"`
	AgentName    string   `json:"agent_name"`
	OverallRank  int      `json:"overall_rank"`
	OverallScore float64  `json:"overall_score"`
	Highlights   []string `json:"highlights"`
}

// PerformanceAlert 性能告警
type PerformanceAlert struct {
	AlertID      string                 `json:"alert_id"`
	AgentID      string                 `json:"agent_id"`
	AlertType    string                 `json:"alert_type"`
	Severity     string                 `json:"severity"`
	Metrics      *PerformanceMetrics    `json:"metrics"`
	Message      string                 `json:"message"`
	CreatedAt    time.Time              `json:"created_at"`
}

// PerformanceMetrics 性能指标
type PerformanceMetrics struct {
	QPS            float64 `json:"qps"`
	ResponseTime   float64 `json:"response_time"`
	ErrorRate      float64 `json:"error_rate"`
	Satisfaction   float64 `json:"satisfaction"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
