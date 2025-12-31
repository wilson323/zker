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

package entity

import (
	"time"
)

// ExecutionStatus 执行状态
type ExecutionStatus string

const (
	ExecutionStatusRunning   ExecutionStatus = "running"   // 运行中
	ExecutionStatusCompleted ExecutionStatus = "completed" // 完成
	ExecutionStatusFailed    ExecutionStatus = "failed"    // 失败
	ExecutionStatusTimeout   ExecutionStatus = "timeout"   // 超时
	ExecutionStatusCancelled ExecutionStatus = "cancelled" // 已取消
)

// PluginExecution 插件执行记录
type PluginExecution struct {
	ID          int64             `json:"id"`
	ExecutionID string            `json:"execution_id"`
	PluginID    int64             `json:"plugin_id"`
	PluginName  string            `json:"plugin_name"`
	PluginVersion string          `json:"plugin_version"`
	Parameters  map[string]interface{} `json:"parameters"`
	Status      ExecutionStatus   `json:"status"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string            `json:"error,omitempty"`
	RetryCount  int               `json:"retry_count"`
	MaxRetries  int               `json:"max_retries"`
	Timeout     int               `json:"timeout"`       // 超时时间(秒)
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	DurationMs  int64             `json:"duration_ms"`   // 执行时长(毫秒)
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	TenantID    string            `json:"tenant_id"`
	UserID      string            `json:"user_id"`
}

// IsCompleted 判断是否完成
func (e *PluginExecution) IsCompleted() bool {
	return e.Status == ExecutionStatusCompleted ||
		e.Status == ExecutionStatusFailed ||
		e.Status == ExecutionStatusTimeout ||
		e.Status == ExecutionStatusCancelled
}

// CanRetry 判断是否可以重试
func (e *PluginExecution) CanRetry() bool {
	return !e.IsCompleted() && e.RetryCount < e.MaxRetries
}

// ExecutePluginRequest 执行插件请求
type ExecutePluginRequest struct {
	PluginName string                 `json:"plugin_name"`
	PluginVersion string              `json:"plugin_version,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
	Timeout     int                    `json:"timeout,omitempty"`      // 超时时间(秒),0表示使用默认值
	MaxRetries  int                    `json:"max_retries,omitempty"`  // 最大重试次数,默认3
	TenantID    string                 `json:"tenant_id"`
	UserID      string                 `json:"user_id"`
}

// ExecutePluginResponse 执行插件响应
type ExecutePluginResponse struct {
	ExecutionID string            `json:"execution_id"`
	Status      ExecutionStatus   `json:"status"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string            `json:"error,omitempty"`
}

// GetExecutionRequest 获取执行记录请求
type GetExecutionRequest struct {
	ExecutionID string `json:"execution_id"`
	TenantID    string `json:"tenant_id"`
}

// GetExecutionResponse 获取执行记录响应
type GetExecutionResponse struct {
	Execution   *PluginExecution `json:"execution"`
	Result      map[string]interface{} `json:"result,omitempty"`
	Error       string            `json:"error,omitempty"`
}

// PluginMetrics 插件指标
type PluginMetrics struct {
	PluginName         string  `json:"plugin_name"`
	TotalCalls         int     `json:"total_calls"`
	SuccessCount       int     `json:"success_count"`
	ErrorCount         int     `json:"error_count"`
	SuccessRate        float64 `json:"success_rate"`
	ErrorRate          float64 `json:"error_rate"`
	AvgResponseTimeMs  int64   `json:"avg_response_time_ms"`
	P95ResponseTimeMs  int64   `json:"p95_response_time_ms"`
	P99ResponseTimeMs  int64   `json:"p99_response_time_ms"`
	MinResponseTimeMs  int64   `json:"min_response_time_ms"`
	MaxResponseTimeMs  int64   `json:"max_response_time_ms"`
	ResponseTimes      []int64 `json:"-"` // 原始响应时间,用于计算分位数
}

// CalculateDuration 计算执行时长
func (e *PluginExecution) CalculateDuration() {
	if e.CompletedAt != nil && !e.CompletedAt.IsZero() {
		e.DurationMs = e.CompletedAt.Sub(e.StartedAt).Milliseconds()
	}
}

// SetResult 设置执行结果
func (e *PluginExecution) SetResult(result map[string]interface{}) {
	e.Result = result
	e.Status = ExecutionStatusCompleted
	now := time.Now()
	e.CompletedAt = &now
	e.CalculateDuration()
}

// SetError 设置错误
func (e *PluginExecution) SetError(err error) {
	e.Error = err.Error()
	e.Status = ExecutionStatusFailed
	now := time.Now()
	e.CompletedAt = &now
	e.CalculateDuration()
}

// IncrementRetry 增加重试次数
func (e *PluginExecution) IncrementRetry() {
	e.RetryCount++
}
