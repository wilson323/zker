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

	"github.com/coze-dev/coze-studio/backend/domain/plugin/entity"
)

// PluginMonitorService 插件监控服务接口
type PluginMonitorService interface {
	// GetMetrics 获取插件指标
	GetMetrics(
		ctx context.Context,
		req *GetMetricsRequest,
	) (*GetMetricsResponse, error)

	// GetErrors 获取插件错误统计
	GetErrors(
		ctx context.Context,
		req *GetErrorsRequest,
	) (*GetErrorsResponse, error)

	// GetTrends 获取插件调用趋势
	GetTrends(
		ctx context.Context,
		req *GetTrendsRequest,
	) (*GetTrendsResponse, error)

	// RecordExecution 记录执行指标
	RecordExecution(
		ctx context.Context,
		execution *entity.PluginExecution,
	) error

	// GetRealTimeMetrics 获取实时指标
	GetRealTimeMetrics(
		ctx context.Context,
		pluginName string,
	) (*RealTimeMetrics, error)
}

// GetMetricsRequest 获取指标请求
type GetMetricsRequest struct {
	PluginName string `json:"plugin_name,omitempty"` // 可选,为空则返回所有插件
	StartTime  int64  `json:"start_time"`            // Unix时间戳
	EndTime    int64  `json:"end_time"`              // Unix时间戳
	TenantID   string `json:"tenant_id"`
}

// GetMetricsResponse 获取指标响应
type GetMetricsResponse struct {
	Metrics []*entity.PluginMetrics `json:"metrics"`
}

// GetErrorsRequest 获取错误请求
type GetErrorsRequest struct {
	PluginName string `json:"plugin_name"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	TenantID   string `json:"tenant_id"`
}

// GetErrorsResponse 获取错误响应
type GetErrorsResponse struct {
	Errors []*ErrorStat `json:"errors"`
}

// ErrorStat 错误统计
type ErrorStat struct {
	ErrorType   string  `json:"error_type"`   // 错误类型: timeout, connection_failed, error_response等
	Count       int     `json:"count"`         // 错误次数
	Percentage  float64 `json:"percentage"`    // 错误占比
	Examples    []string `json:"examples"`     // 错误示例
}

// GetTrendsRequest 获取趋势请求
type GetTrendsRequest struct {
	PluginName string `json:"plugin_name"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
	Interval   string `json:"interval"` // hour, day, week
	TenantID   string `json:"tenant_id"`
}

// GetTrendsResponse 获取趋势响应
type GetTrendsResponse struct {
	Trends *TrendData `json:"trends"`
}

// TrendData 趋势数据
type TrendData struct {
	TimeLabels     []string  `json:"time_labels"`      // 时间标签
	DailyCalls     []int     `json:"daily_calls"`      // 每日调用次数
	DailySuccess   []int     `json:"daily_success"`    // 每日成功次数
	DailyErrors    []int     `json:"daily_errors"`     // 每日错误次数
	SuccessRate    []float64 `json:"success_rate"`     // 成功率
	AvgResponseTime []int64  `json:"avg_response_time"` // 平均响应时间
}

// RealTimeMetrics 实时指标
type RealTimeMetrics struct {
	PluginName         string   `json:"plugin_name"`
	CurrentQPS         float64  `json:"current_qps"`
	ActiveExecutions   int      `json:"active_executions"`
	AvgResponseTimeMs  int64    `json:"avg_response_time_ms"`
	SuccessRate        float64  `json:"success_rate"`
	ErrorRate          float64  `json:"error_rate"`
	P95ResponseTimeMs  int64    `json:"p95_response_time_ms"`
	P99ResponseTimeMs  int64    `json:"p99_response_time_ms"`
}
