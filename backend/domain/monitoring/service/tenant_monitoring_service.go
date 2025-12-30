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
)

// TenantMonitoringService 租户监控服务
type TenantMonitoringService struct {
	metricsRepo   repository.TenantMetricsRepository
	logger        *zap.Logger
	metricsCache  *metricsCache
	mu            sync.RWMutex
}

// metricsCache 指标缓存
type metricsCache struct {
	data      map[string]*cachedMetric
	expireAt  time.Time
	mu        sync.RWMutex
}

type cachedMetric struct {
	metric    interface{}
	timestamp time.Time
}

// NewTenantMonitoringService 创建租户监控服务实例
func NewTenantMonitoringService(
	metricsRepo repository.TenantMetricsRepository,
	logger *zap.Logger,
) *TenantMonitoringService {
	cache := &metricsCache{
		data:     make(map[string]*cachedMetric),
		expireAt: time.Now().Add(5 * time.Minute),
	}

	return &TenantMonitoringService{
		metricsRepo:  metricsRepo,
		logger:       logger,
		metricsCache: cache,
	}
}

// CollectQPS 收集QPS指标
func (s *TenantMonitoringService) CollectQPS(
	ctx context.Context,
	tenantID string,
	timeRange *entity.TimeRange,
) (*entity.QPSMetrics, error) {
	s.logger.Info("Collecting QPS metrics",
		zap.String("tenant_id", tenantID),
		zap.Time("start", timeRange.StartTime),
		zap.Time("end", timeRange.EndTime),
	)

	// 检查缓存
	cacheKey := fmt.Sprintf("qps:%s:%d:%d", tenantID, timeRange.StartTime.Unix(), timeRange.EndTime.Unix())
	if cached := s.getFromCache(cacheKey); cached != nil {
		if metrics, ok := cached.(*entity.QPSMetrics); ok {
			return metrics, nil
		}
	}

	// 从数据库查询
	metrics, err := s.metricsRepo.CalculateQPSMetrics(ctx, tenantID, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		s.logger.Error("Failed to calculate QPS metrics",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to calculate QPS metrics: %w", err)
	}

	// 计算趋势
	trend, change := s.calculateTrend(ctx, tenantID, entity.MetricTypeQPS, metrics.Average, timeRange)
	metrics.Trend = trend
	metrics.TrendChange = change

	// 缓存结果
	s.putToCache(cacheKey, metrics, 1*time.Minute)

	return metrics, nil
}

// CollectResponseTime 收集响应时间指标
func (s *TenantMonitoringService) CollectResponseTime(
	ctx context.Context,
	tenantID string,
	timeRange *entity.TimeRange,
) (*entity.ResponseTimeMetrics, error) {
	s.logger.Info("Collecting response time metrics",
		zap.String("tenant_id", tenantID),
	)

	cacheKey := fmt.Sprintf("rt:%s:%d:%d", tenantID, timeRange.StartTime.Unix(), timeRange.EndTime.Unix())
	if cached := s.getFromCache(cacheKey); cached != nil {
		if metrics, ok := cached.(*entity.ResponseTimeMetrics); ok {
			return metrics, nil
		}
	}

	metrics, err := s.metricsRepo.CalculateResponseTimeMetrics(ctx, tenantID, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		s.logger.Error("Failed to calculate response time metrics",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to calculate response time metrics: %w", err)
	}

	// 计算趋势
	trend, change := s.calculateTrend(ctx, tenantID, entity.MetricTypeResponseTime, metrics.Average, timeRange)
	metrics.Trend = trend
	metrics.TrendChange = change

	s.putToCache(cacheKey, metrics, 1*time.Minute)

	return metrics, nil
}

// CollectErrorRate 收集错误率指标
func (s *TenantMonitoringService) CollectErrorRate(
	ctx context.Context,
	tenantID string,
	timeRange *entity.TimeRange,
) (*entity.ErrorRateMetrics, error) {
	s.logger.Info("Collecting error rate metrics",
		zap.String("tenant_id", tenantID),
	)

	cacheKey := fmt.Sprintf("er:%s:%d:%d", tenantID, timeRange.StartTime.Unix(), timeRange.EndTime.Unix())
	if cached := s.getFromCache(cacheKey); cached != nil {
		if metrics, ok := cached.(*entity.ErrorRateMetrics); ok {
			return metrics, nil
		}
	}

	metrics, err := s.metricsRepo.CalculateErrorRateMetrics(ctx, tenantID, timeRange.StartTime, timeRange.EndTime)
	if err != nil {
		s.logger.Error("Failed to calculate error rate metrics",
			zap.String("tenant_id", tenantID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to calculate error rate metrics: %w", err)
	}

	trend, change := s.calculateTrend(ctx, tenantID, entity.MetricTypeErrorRate, metrics.ErrorRate, timeRange)
	metrics.Trend = trend
	metrics.TrendChange = change

	s.putToCache(cacheKey, metrics, 1*time.Minute)

	return metrics, nil
}

// CalculateTenantHealthScore 计算租户健康评分
func (s *TenantMonitoringService) CalculateTenantHealthScore(
	ctx context.Context,
	tenantID string,
) (*entity.HealthScore, error) {
	s.logger.Info("Calculating tenant health score",
		zap.String("tenant_id", tenantID),
	)

	timeRange := entity.NewTimeRangeFromDuration(24 * time.Hour)

	// 并行收集各类指标
	var wg sync.WaitGroup
	var qpsMetrics *entity.QPSMetrics
	var responseMetrics *entity.ResponseTimeMetrics
	var errorMetrics *entity.ErrorRateMetrics
	var qpsErr, responseErr, errorErr error

	wg.Add(3)

	go func() {
		defer wg.Done()
		qpsMetrics, qpsErr = s.CollectQPS(ctx, tenantID, timeRange)
	}()

	go func() {
		defer wg.Done()
		responseMetrics, responseErr = s.CollectResponseTime(ctx, tenantID, timeRange)
	}()

	go func() {
		defer wg.Done()
		errorMetrics, errorErr = s.CollectErrorRate(ctx, tenantID, timeRange)
	}()

	wg.Wait()

	if qpsErr != nil || responseErr != nil || errorErr != nil {
		return nil, fmt.Errorf("failed to collect metrics: qps=%v, response=%v, error=%v",
			qpsErr, responseErr, errorErr)
	}

	// 计算各项评分
	qpsScore := s.calculateQPSScore(qpsMetrics)
	responseScore := s.calculateResponseScore(responseMetrics)
	errorScore := s.calculateErrorScore(errorMetrics)
	resourceScore := 80.0 // 简化处理，实际应查询资源使用率

	// 计算总分
	totalScore := (qpsScore + responseScore + errorScore + resourceScore) / 4.0

	healthScore := &entity.HealthScore{
		Score:         totalScore,
		Level:         entity.GetHealthLevel(totalScore),
		QPSScore:      qpsScore,
		ResponseScore: responseScore,
		ErrorScore:    errorScore,
		ResourceScore: resourceScore,
		Trend:         "stable",
		AssessedAt:    time.Now(),
	}

	// 评估趋势
	if qpsMetrics.Trend == "up" && responseMetrics.Trend == "down" && errorMetrics.Trend == "down" {
		healthScore.Trend = "improving"
	} else if qpsMetrics.Trend == "down" || responseMetrics.Trend == "up" || errorMetrics.Trend == "up" {
		healthScore.Trend = "degrading"
	}

	return healthScore, nil
}

// GetTenantOverview 获取租户概览
func (s *TenantMonitoringService) GetTenantOverview(
	ctx context.Context,
	tenantID, tenantName string,
) (*entity.TenantOverview, error) {
	s.logger.Info("Getting tenant overview",
		zap.String("tenant_id", tenantID),
	)

	timeRange := entity.NewTimeRangeFromDuration(1 * time.Hour)

	// 并行收集数据
	var wg sync.WaitGroup
	var healthScore *entity.HealthScore
	var qpsMetrics *entity.QPSMetrics
	var responseMetrics *entity.ResponseTimeMetrics
	var errorMetrics *entity.ErrorRateMetrics
	var healthErr, qpsErr, responseErr, errorErr error

	wg.Add(4)

	go func() {
		defer wg.Done()
		healthScore, healthErr = s.CalculateTenantHealthScore(ctx, tenantID)
	}()

	go func() {
		defer wg.Done()
		qpsMetrics, qpsErr = s.CollectQPS(ctx, tenantID, timeRange)
	}()

	go func() {
		defer wg.Done()
		responseMetrics, responseErr = s.CollectResponseTime(ctx, tenantID, timeRange)
	}()

	go func() {
		defer wg.Done()
		errorMetrics, errorErr = s.CollectErrorRate(ctx, tenantID, timeRange)
	}()

	wg.Wait()

	// 记录非关键错误（仅用于日志，不中断流程）
	if qpsErr != nil {
		s.logger.Warn("收集QPS指标失败", zap.Error(qpsErr))
	}
	if responseErr != nil {
		s.logger.Warn("收集响应时间指标失败", zap.Error(responseErr))
	}
	if errorErr != nil {
		s.logger.Warn("收集错误率指标失败", zap.Error(errorErr))
	}

	if healthErr != nil {
		return nil, healthErr
	}

	overview := &entity.TenantOverview{
		TenantID:        tenantID,
		TenantName:      tenantName,
		HealthScore:     *healthScore,
		CurrentQPS:      0,
		AvgResponseTime: 0,
		ErrorRate:       0,
		TotalRequests:   0,
		CPUUsage:        0,
		MemoryUsage:     0,
		DiskUsage:       0,
		AlertCount:      0,
		LastUpdated:     time.Now(),
	}

	if qpsMetrics != nil {
		overview.CurrentQPS = qpsMetrics.Average
		overview.TotalRequests = qpsMetrics.Count
	}

	if responseMetrics != nil {
		overview.AvgResponseTime = responseMetrics.Average
	}

	if errorMetrics != nil {
		overview.ErrorRate = errorMetrics.ErrorRate
	}

	return overview, nil
}

// RecordMetric 记录指标
func (s *TenantMonitoringService) RecordMetric(
	ctx context.Context,
	tenantID string,
	metricType entity.MetricType,
	value float64,
	tags entity.MetricTags,
) error {
	metric := entity.NewTenantMetrics(tenantID, metricType, value, tags)
	err := s.metricsRepo.Create(ctx, metric)
	if err != nil {
		s.logger.Error("Failed to record metric",
			zap.String("tenant_id", tenantID),
			zap.String("metric_type", string(metricType)),
			zap.Float64("value", value),
			zap.Error(err),
		)
		return fmt.Errorf("failed to record metric: %w", err)
	}

	// 清除相关缓存
	s.clearCacheForTenant(tenantID, metricType)

	return nil
}

// RecordMetricBatch 批量记录指标
func (s *TenantMonitoringService) RecordMetricBatch(
	ctx context.Context,
	metricsList []*entity.TenantMetrics,
) error {
	if len(metricsList) == 0 {
		return nil
	}

	err := s.metricsRepo.CreateBatch(ctx, metricsList)
	if err != nil {
		s.logger.Error("Failed to record metrics batch",
			zap.Int("count", len(metricsList)),
			zap.Error(err),
		)
		return fmt.Errorf("failed to record metrics batch: %w", err)
	}

	// 清除缓存
	tenantIDs := make(map[string]struct{})
	metricTypes := make(map[entity.MetricType]struct{})
	for _, m := range metricsList {
		tenantIDs[m.TenantID] = struct{}{}
		metricTypes[m.MetricType] = struct{}{}
	}

	for tenantID := range tenantIDs {
		for metricType := range metricTypes {
			s.clearCacheForTenant(tenantID, metricType)
		}
	}

	return nil
}

// ============================================================
// 私有方法
// ============================================================

// calculateTrend 计算趋势
func (s *TenantMonitoringService) calculateTrend(
	ctx context.Context,
	tenantID string,
	metricType entity.MetricType,
	currentValue float64,
	timeRange *entity.TimeRange,
) (string, float64) {
	// 获取上一个时间周期的数据
	previousDuration := timeRange.Duration()
	previousStartTime := timeRange.StartTime.Add(-previousDuration)
	previousEndTime := timeRange.StartTime

	var previousValue float64
	var err error

	switch metricType {
	case entity.MetricTypeQPS:
		var metrics *entity.QPSMetrics
		metrics, err = s.metricsRepo.CalculateQPSMetrics(ctx, tenantID, previousStartTime, previousEndTime)
		if metrics != nil {
			previousValue = metrics.Average
		}
	case entity.MetricTypeResponseTime:
		var metrics *entity.ResponseTimeMetrics
		metrics, err = s.metricsRepo.CalculateResponseTimeMetrics(ctx, tenantID, previousStartTime, previousEndTime)
		if metrics != nil {
			previousValue = metrics.Average
		}
	case entity.MetricTypeErrorRate:
		var metrics *entity.ErrorRateMetrics
		metrics, err = s.metricsRepo.CalculateErrorRateMetrics(ctx, tenantID, previousStartTime, previousEndTime)
		if metrics != nil {
			previousValue = metrics.ErrorRate
		}
	}

	if err != nil || previousValue == 0 {
		return "stable", 0
	}

	change := ((currentValue - previousValue) / previousValue) * 100

	switch {
	case change > 10:
		return "up", change
	case change < -10:
		return "down", change
	default:
		return "stable", change
	}
}

// calculateQPSScore 计算QPS评分
func (s *TenantMonitoringService) calculateQPSScore(metrics *entity.QPSMetrics) float64 {
	// 简化评分逻辑
	if metrics == nil {
		return 50.0
	}

	// 基于QPS值和趋势评分
	score := 60.0
	if metrics.Average > 100 {
		score = 90.0
	} else if metrics.Average > 50 {
		score = 75.0
	} else if metrics.Average > 10 {
		score = 60.0
	} else {
		score = 40.0
	}

	// 趋势调整
	if metrics.Trend == "up" {
		score = min(100, score+10)
	} else if metrics.Trend == "down" {
		score = max(0, score-10)
	}

	return score
}

// calculateResponseScore 计算响应时间评分
func (s *TenantMonitoringService) calculateResponseScore(metrics *entity.ResponseTimeMetrics) float64 {
	if metrics == nil {
		return 50.0
	}

	score := 60.0
	if metrics.Average < 100 {
		score = 90.0
	} else if metrics.Average < 500 {
		score = 75.0
	} else if metrics.Average < 1000 {
		score = 60.0
	} else {
		score = 30.0
	}

	if metrics.Trend == "down" {
		score = min(100, score+10)
	} else if metrics.Trend == "up" {
		score = max(0, score-10)
	}

	return score
}

// calculateErrorScore 计算错误率评分
func (s *TenantMonitoringService) calculateErrorScore(metrics *entity.ErrorRateMetrics) float64 {
	if metrics == nil {
		return 50.0
	}

	score := 100.0 - (metrics.ErrorRate * 100)
	if metrics.Trend == "down" {
		score = min(100, score+10)
	} else if metrics.Trend == "up" {
		score = max(0, score-10)
	}

	return score
}

// ============================================================
// 缓存方法
// ============================================================

func (s *TenantMonitoringService) getFromCache(key string) interface{} {
	s.metricsCache.mu.RLock()
	defer s.metricsCache.mu.RUnlock()

	if time.Now().After(s.metricsCache.expireAt) {
		return nil
	}

	cached, exists := s.metricsCache.data[key]
	if !exists {
		return nil
	}

	if time.Since(cached.timestamp) > 5*time.Minute {
		return nil
	}

	return cached.metric
}

func (s *TenantMonitoringService) putToCache(key string, metric interface{}, ttl time.Duration) {
	s.metricsCache.mu.Lock()
	defer s.metricsCache.mu.Unlock()

	s.metricsCache.data[key] = &cachedMetric{
		metric:    metric,
		timestamp: time.Now(),
	}
}

func (s *TenantMonitoringService) clearCacheForTenant(tenantID string, metricType entity.MetricType) {
	s.metricsCache.mu.Lock()
	defer s.metricsCache.mu.Unlock()

	// 简化处理：清除所有缓存
	s.metricsCache.data = make(map[string]*cachedMetric)
	s.metricsCache.expireAt = time.Now().Add(5 * time.Minute)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
