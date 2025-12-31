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
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
)

// Prometheus指标
var (
	pluginExecutionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "plugin_execution_duration_seconds",
			Help:    "Plugin execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"plugin_name", "status"},
	)

	pluginExecutionTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plugin_executions_total",
			Help: "Total number of plugin executions",
		},
		[]string{"plugin_name", "status"},
	)

	pluginExecutionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "plugin_execution_errors_total",
			Help: "Total number of plugin execution errors",
		},
		[]string{"plugin_name", "error_type"},
	)
)

func init() {
	prometheus.MustRegister(pluginExecutionDuration)
	prometheus.MustRegister(pluginExecutionTotal)
	prometheus.MustRegister(pluginExecutionErrors)
}

// NewPluginMonitorService 创建插件监控服务
func NewPluginMonitorService(db *gorm.DB) PluginMonitorService {
	return &pluginMonitorServiceImpl{
		db:            db,
		metricsCache:  make(map[string]*entity.PluginMetrics),
		realTimeStats: make(map[string]*RealTimeStats),
		cacheMutex:    &sync.RWMutex{},
	}
}

type pluginMonitorServiceImpl struct {
	db            *gorm.DB
	metricsCache  map[string]*entity.PluginMetrics
	realTimeStats map[string]*RealTimeStats
	cacheMutex    *sync.RWMutex
}

// RealTimeStats 实时统计
type RealTimeStats struct {
	TotalCalls       int64
	SuccessCount     int64
	ErrorCount       int64
	TotalDurationMs  int64
	ResponseTimes    []int64
	WindowStartTime  time.Time
}

// GetMetrics 获取插件指标
func (s *pluginMonitorServiceImpl) GetMetrics(
	ctx context.Context,
	req *GetMetricsRequest,
) (*GetMetricsResponse, error) {
	// 1. 查询执行记录
	var executions []entity.PluginExecution
	query := s.db.Table("plugin_executions").
		Where("created_at >= ? AND created_at <= ?", req.StartTime, req.EndTime)

	if req.PluginName != "" {
		query = query.Where("plugin_name = ?", req.PluginName)
	}

	if err := query.Find(&executions).Error; err != nil {
		return nil, fmt.Errorf("failed to query executions: %w", err)
	}

	// 2. 计算指标
	metricsMap := s.calculateMetricsMap(executions)

	// 3. 转换为切片
	metrics := make([]*entity.PluginMetrics, 0, len(metricsMap))
	for _, metric := range metricsMap {
		metrics = append(metrics, metric)
	}

	return &GetMetricsResponse{
		Metrics: metrics,
	}, nil
}

// GetErrors 获取插件错误统计
func (s *pluginMonitorServiceImpl) GetErrors(
	ctx context.Context,
	req *GetErrorsRequest,
) (*GetErrorsResponse, error) {
	// 1. 查询失败的执行记录
	var executions []entity.PluginExecution
	err := s.db.Table("plugin_executions").
		Where("plugin_name = ? AND created_at >= ? AND created_at <= ? AND status = ?",
			req.PluginName, req.StartTime, req.EndTime, entity.ExecutionStatusFailed).
		Find(&executions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query failed executions: %w", err)
	}

	// 2. 统计错误类型
	errorMap := make(map[string]*ErrorStat)
	totalErrors := len(executions)

	for _, exec := range executions {
		errorType := s.classifyError(exec.Error)
		if _, exists := errorMap[errorType]; !exists {
			errorMap[errorType] = &ErrorStat{
				ErrorType:  errorType,
				Count:      0,
				Examples:   make([]string, 0, 3),
			}
		}
		stat := errorMap[errorType]
		stat.Count++
		if len(stat.Examples) < 3 {
			stat.Examples = append(stat.Examples, exec.Error)
		}
	}

	// 3. 计算百分比
	errors := make([]*ErrorStat, 0, len(errorMap))
	for _, stat := range errorMap {
		stat.Percentage = float64(stat.Count) / float64(totalErrors)
		errors = append(errors, stat)
	}

	return &GetErrorsResponse{
		Errors: errors,
	}, nil
}

// GetTrends 获取插件调用趋势
func (s *pluginMonitorServiceImpl) GetTrends(
	ctx context.Context,
	req *GetTrendsRequest,
) (*GetTrendsResponse, error) {
	// 1. 确定时间间隔
	var interval string
	switch req.Interval {
	case "hour", "day", "week":
		interval = req.Interval
	default:
		interval = "day"
	}

	// 2. 查询执行记录
	var executions []entity.PluginExecution
	err := s.db.Table("plugin_executions").
		Where("plugin_name = ? AND created_at >= ? AND created_at <= ?",
			req.PluginName, req.StartTime, req.EndTime).
		Order("created_at ASC").
		Find(&executions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query executions: %w", err)
	}

	// 3. 按时间分组统计
	trends := s.groupByTime(executions, interval)

	return &GetTrendsResponse{
		Trends: trends,
	}, nil
}

// RecordExecution 记录执行指标
func (s *pluginMonitorServiceImpl) RecordExecution(
	ctx context.Context,
	execution *entity.PluginExecution,
) error {
	// 1. 更新Prometheus指标
	duration := float64(execution.DurationMs) / 1000.0 // 转换为秒
	pluginExecutionDuration.WithLabelValues(
		execution.PluginName,
		string(execution.Status),
	).Observe(duration)

	pluginExecutionTotal.WithLabelValues(
		execution.PluginName,
		string(execution.Status),
	).Inc()

	if execution.Status == entity.ExecutionStatusFailed {
		errorType := s.classifyError(execution.Error)
		pluginExecutionErrors.WithLabelValues(
			execution.PluginName,
			errorType,
		).Inc()
	}

	// 2. 更新实时统计
	s.updateRealTimeStats(execution)

	return nil
}

// GetRealTimeMetrics 获取实时指标
func (s *pluginMonitorServiceImpl) GetRealTimeMetrics(
	ctx context.Context,
	pluginName string,
) (*RealTimeMetrics, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	stats, exists := s.realTimeStats[pluginName]
	if !exists {
		return &RealTimeMetrics{
			PluginName: pluginName,
		}, nil
	}

	// 计算指标
	totalCalls := stats.TotalCalls
	successRate := 0.0
	if totalCalls > 0 {
		successRate = float64(stats.SuccessCount) / float64(totalCalls)
	}

	errorRate := 0.0
	if totalCalls > 0 {
		errorRate = float64(stats.ErrorCount) / float64(totalCalls)
	}

	avgResponseTime := int64(0)
	if totalCalls > 0 {
		avgResponseTime = stats.TotalDurationMs / totalCalls
	}

	// 计算P95和P99
	p95 := int64(0)
	p99 := int64(0)
	if len(stats.ResponseTimes) > 0 {
		sort.Slice(stats.ResponseTimes, func(i, j int) bool {
			return stats.ResponseTimes[i] < stats.ResponseTimes[j]
		})
		p95Index := int(float64(len(stats.ResponseTimes)) * 0.95)
		p99Index := int(float64(len(stats.ResponseTimes)) * 0.99)
		if p95Index < len(stats.ResponseTimes) {
			p95 = stats.ResponseTimes[p95Index]
		}
		if p99Index < len(stats.ResponseTimes) {
			p99 = stats.ResponseTimes[p99Index]
		}
	}

	return &RealTimeMetrics{
		PluginName:        pluginName,
		CurrentQPS:        0, // TODO: 计算QPS
		ActiveExecutions:  0, // TODO: 从执行缓存获取
		AvgResponseTimeMs: avgResponseTime,
		SuccessRate:       successRate,
		ErrorRate:         errorRate,
		P95ResponseTimeMs: p95,
		P99ResponseTimeMs: p99,
	}, nil
}

// calculateMetricsMap 计算指标映射
func (s *pluginMonitorServiceImpl) calculateMetricsMap(
	executions []entity.PluginExecution,
) map[string]*entity.PluginMetrics {
	metricsMap := make(map[string]*entity.PluginMetrics)

	for _, exec := range executions {
		if _, exists := metricsMap[exec.PluginName]; !exists {
			metricsMap[exec.PluginName] = &entity.PluginMetrics{
				PluginName:    exec.PluginName,
				TotalCalls:    0,
				SuccessCount:  0,
				ErrorCount:    0,
				ResponseTimes: make([]int64, 0),
			}
		}

		metric := metricsMap[exec.PluginName]
		metric.TotalCalls++

		if exec.Status == entity.ExecutionStatusCompleted {
			metric.SuccessCount++
		} else {
			metric.ErrorCount++
		}

		if exec.DurationMs > 0 {
			metric.ResponseTimes = append(metric.ResponseTimes, exec.DurationMs)
		}
	}

	// 计算统计值
	for _, metric := range metricsMap {
		if metric.TotalCalls > 0 {
			metric.SuccessRate = float64(metric.SuccessCount) / float64(metric.TotalCalls)
			metric.ErrorRate = float64(metric.ErrorCount) / float64(metric.TotalCalls)
		}

		if len(metric.ResponseTimes) > 0 {
			sort.Slice(metric.ResponseTimes, func(i, j int) bool {
				return metric.ResponseTimes[i] < metric.ResponseTimes[j]
			})

			total := int64(0)
			for _, t := range metric.ResponseTimes {
				total += t
			}
			metric.AvgResponseTimeMs = total / int64(len(metric.ResponseTimes))

			minIndex := 0
			maxIndex := len(metric.ResponseTimes) - 1
			metric.MinResponseTimeMs = metric.ResponseTimes[minIndex]
			metric.MaxResponseTimeMs = metric.ResponseTimes[maxIndex]

			p95Index := int(float64(len(metric.ResponseTimes)) * 0.95)
			p99Index := int(float64(len(metric.ResponseTimes)) * 0.99)
			if p95Index < len(metric.ResponseTimes) {
				metric.P95ResponseTimeMs = metric.ResponseTimes[p95Index]
			}
			if p99Index < len(metric.ResponseTimes) {
				metric.P99ResponseTimeMs = metric.ResponseTimes[p99Index]
			}
		}
	}

	return metricsMap
}

// classifyError 分类错误
func (s *pluginMonitorServiceImpl) classifyError(errorMessage string) string {
	// 简化实现:根据错误消息分类
	if contains(errorMessage, "timeout") || contains(errorMessage, "deadline") {
		return "timeout"
	}
	if contains(errorMessage, "connection") || contains(errorMessage, "connect") {
		return "connection_failed"
	}
	if contains(errorMessage, "401") || contains(errorMessage, "403") {
		return "auth_failed"
	}
	if contains(errorMessage, "500") || contains(errorMessage, "502") || contains(errorMessage, "503") {
		return "server_error"
	}
	return "unknown"
}

// groupByTime 按时间分组
func (s *pluginMonitorServiceImpl) groupByTime(
	executions []entity.PluginExecution,
	interval string,
) *TrendData {
	// 简化实现:返回基本趋势
	trends := &TrendData{
		TimeLabels:      make([]string, 0),
		DailyCalls:      make([]int, 0),
		DailySuccess:    make([]int, 0),
		DailyErrors:     make([]int, 0),
		SuccessRate:     make([]float64, 0),
		AvgResponseTime: make([]int64, 0),
	}

	// TODO: 实现时间分组逻辑
	// 这里简化为返回空数据
	return trends
}

// updateRealTimeStats 更新实时统计
func (s *pluginMonitorServiceImpl) updateRealTimeStats(execution *entity.PluginExecution) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if _, exists := s.realTimeStats[execution.PluginName]; !exists {
		s.realTimeStats[execution.PluginName] = &RealTimeStats{
			ResponseTimes:   make([]int64, 0, 1000),
			WindowStartTime: time.Now(),
		}
	}

	stats := s.realTimeStats[execution.PluginName]
	stats.TotalCalls++

	if execution.Status == entity.ExecutionStatusCompleted {
		stats.SuccessCount++
	} else {
		stats.ErrorCount++
	}

	if execution.DurationMs > 0 {
		stats.TotalDurationMs += execution.DurationMs
		stats.ResponseTimes = append(stats.ResponseTimes, execution.DurationMs)
	}

	// 滑动窗口:保留最近1000条记录
	if len(stats.ResponseTimes) > 1000 {
		stats.ResponseTimes = stats.ResponseTimes[len(stats.ResponseTimes)-1000:]
	}

	// 重置窗口(每小时)
	if time.Since(stats.WindowStartTime) > time.Hour {
		stats.TotalCalls = 0
		stats.SuccessCount = 0
		stats.ErrorCount = 0
		stats.TotalDurationMs = 0
		stats.ResponseTimes = make([]int64, 0, 1000)
		stats.WindowStartTime = time.Now()
	}
}

// contains 检查字符串包含
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
