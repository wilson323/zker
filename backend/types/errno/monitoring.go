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

package errno

import (
	"net/http"
)

// 监控模块错误码
var (
	// MetricsCollectionFailed 指标收集失败
	MetricsCollectionFailed = &BaseErrorCode{
		code:       "MONITOR500001",
		message:    "Failed to collect metrics",
		messageZH:  "指标收集失败",
		messageEN:  "Failed to collect metrics",
		httpStatus: http.StatusInternalServerError,
	}

	// PrometheusNotAvailable Prometheus 不可用
	PrometheusNotAvailable = &BaseErrorCode{
		code:       "MONITOR503001",
		message:    "Prometheus service is not available",
		messageZH:  "Prometheus 服务不可用",
		messageEN:  "Prometheus service is not available",
		httpStatus: http.StatusServiceUnavailable,
	}

	// MetricsExportFailed 指标导出失败
	MetricsExportFailed = &BaseErrorCode{
		code:       "MONITOR500002",
		message:    "Failed to export metrics",
		messageZH:  "指标导出失败",
		messageEN:  "Failed to export metrics",
		httpStatus: http.StatusInternalServerError,
	}

	// MetricsQueryFailed 指标查询失败
	MetricsQueryFailed = &BaseErrorCode{
		code:       "MONITOR500003",
		message:    "Failed to query metrics",
		messageZH:  "指标查询失败",
		messageEN:  "Failed to query metrics",
		httpStatus: http.StatusInternalServerError,
	}

	// AlertRuleNotFound 告警规则不存在
	AlertRuleNotFound = &BaseErrorCode{
		code:       "MONITOR404001",
		message:    "Alert rule not found",
		messageZH:  "告警规则不存在",
		messageEN:  "Alert rule not found",
		httpStatus: http.StatusNotFound,
	}

	// InvalidAlertConfiguration 无效的告警配置
	InvalidAlertConfiguration = &BaseErrorCode{
		code:       "MONITOR400001",
		message:    "Invalid alert configuration",
		messageZH:  "无效的告警配置",
		messageEN:  "Invalid alert configuration",
		httpStatus: http.StatusBadRequest,
	}
)

// 监控模块错误码 int32 常量
const (
	// 指标收集失败
	MetricsCollectionFailedCode int32 = 55000001
	// Prometheus 不可用
	PrometheusNotAvailableCode int32 = 55300001
	// 指标导出失败
	MetricsExportFailedCode int32 = 55000002
	// 指标查询失败
	MetricsQueryFailedCode int32 = 55000003
	// 告警规则不存在
	AlertRuleNotFoundCode int32 = 55400001
	// 无效的告警配置
	InvalidAlertConfigurationCode int32 = 55400001
	// Agent 健康检查失败
	AgentHealthCheckFailedCode int32 = 55000004
	// 性能报告生成失败
	PerformanceReportGenerationFailedCode int32 = 55000005
)
