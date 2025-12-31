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

package health

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	// globalHealthService 全局健康检查服务实例
	globalHealthService *HealthService
)

// InitGlobalHealthService 初始化全局健康检查服务
func InitGlobalHealthService(version string) {
	globalHealthService = NewHealthService(version)
}

// GetGlobalHealthService 获取全局健康检查服务
func GetGlobalHealthService() *HealthService {
	return globalHealthService
}

// HealthCheckHandler 健康检查HTTP处理器
// GET /health
// 返回简化的健康状态（用于负载均衡器健康检查）
func HealthCheckHandler(ctx context.Context, c *app.RequestContext) {
	if globalHealthService == nil {
		c.String(consts.StatusServiceUnavailable, "DOWN")
		return
	}

	report := globalHealthService.Check(ctx)

	// 如果整体状态健康，返回200，否则返回503
	if report.Status == HealthStatusUp {
		c.String(consts.StatusOK, "UP")
	} else {
		c.String(consts.StatusServiceUnavailable, "DOWN")
	}
}

// HealthCheckDetailedHandler 详细健康检查HTTP处理器
// GET /health/detailed
// 返回完整的健康报告（用于监控和调试）
func HealthCheckDetailedHandler(ctx context.Context, c *app.RequestContext) {
	if globalHealthService == nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"status":    "DOWN",
			"error":     "health service not initialized",
			"timestamp": time.Now().UnixMilli(),
		})
		return
	}

	report := globalHealthService.Check(ctx)

	// 根据整体健康状态设置HTTP状态码
	var httpStatus int
	switch report.Status {
	case HealthStatusUp:
		httpStatus = consts.StatusOK
	case HealthStatusDegraded:
		httpStatus = 200 // 200 OK，但状态为DEGRADED
	case HealthStatusDown:
		httpStatus = consts.StatusServiceUnavailable
	default:
		httpStatus = consts.StatusInternalServerError
	}

	c.JSON(httpStatus, report)
}

// HealthCheckSingleHandler 单项健康检查HTTP处理器
// GET /health/check/{name}
// 返回指定检查项的健康状态
func HealthCheckSingleHandler(ctx context.Context, c *app.RequestContext) {
	if globalHealthService == nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error":     "health service not initialized",
			"timestamp": time.Now().UnixMilli(),
		})
		return
	}

	name := c.Param("name")
	if name == "" {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"error":     "checker name is required",
			"timestamp": time.Now().UnixMilli(),
		})
		return
	}

	result, err := globalHealthService.CheckSingle(ctx, name)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"error":     err.Error(),
			"timestamp": time.Now().UnixMilli(),
		})
		return
	}

	// 根据健康状态设置HTTP状态码
	var httpStatus int
	switch result.Status {
	case HealthStatusUp:
		httpStatus = consts.StatusOK
	case HealthStatusDegraded:
		httpStatus = 200 // 200 OK，但状态为DEGRADED
	case HealthStatusDown:
		httpStatus = consts.StatusServiceUnavailable
	default:
		httpStatus = consts.StatusInternalServerError
	}

	c.JSON(httpStatus, result)
}

// LivenessProbeHandler 存活探针处理器
// GET /health/live
// 用于Kubernetes liveness probe（检查服务是否存活）
func LivenessProbeHandler(ctx context.Context, c *app.RequestContext) {
	c.String(consts.StatusOK, "OK")
}

// ReadinessProbeHandler 就绪探针处理器
// GET /health/ready
// 用于Kubernetes readiness probe（检查服务是否就绪）
func ReadinessProbeHandler(ctx context.Context, c *app.RequestContext) {
	if globalHealthService == nil {
		c.String(consts.StatusServiceUnavailable, "NOT_READY")
		return
	}

	// 简单检查：只检查关键组件是否健康
	report := globalHealthService.Check(ctx)

	// 检查关键组件是否DOWN
	if globalHealthService.hasCriticalDown(report) {
		c.String(consts.StatusServiceUnavailable, "NOT_READY")
		return
	}

	c.String(consts.StatusOK, "READY")
}

// VersionHandler 版本信息处理器
// GET /version
// 返回服务版本信息
func VersionHandler(ctx context.Context, c *app.RequestContext) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"error": "failed to read build info",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"version":      globalHealthService.version,
		"go_version":   info.GoVersion,
		"build_time":   getBuildTime(*info),
		"git_commit":   getGitCommit(*info),
		"git_summary":  info.Main.Version,
		"uptime":       int64(time.Since(globalHealthService.startTime).Seconds()),
		"timestamp":    time.Now().UnixMilli(),
	})
}

// MetricsHandler 自定义指标处理器
// GET /health/metrics
// 返回自定义健康指标（用于Prometheus等监控系统）
func MetricsHandler(ctx context.Context, c *app.RequestContext) {
	if globalHealthService == nil {
		c.String(consts.StatusInternalServerError, "health_service_not_initialized")
		return
	}

	report := globalHealthService.Check(ctx)

	// 生成Prometheus格式的指标
	metrics := generatePrometheusMetrics(report)

	c.Header("Content-Type", "text/plain")
	c.String(consts.StatusOK, metrics)
}

// generatePrometheusMetrics 生成Prometheus格式的指标
func generatePrometheusMetrics(report *HealthReport) string {
	var output string

	// 整体健康状态
	status := 0
	if report.Status == HealthStatusUp {
		status = 1
	} else if report.Status == HealthStatusDegraded {
		status = 2
	}
	// DOWN = 3

	output += fmt.Sprintf("# HELP health_status Overall health status (0=unknown, 1=up, 2=degraded, 3=down)\n")
	output += fmt.Sprintf("# TYPE health_status gauge\n")
	output += fmt.Sprintf("health_status %d\n", status)

	// 各检查项状态
	output += fmt.Sprintf("\n# HELP health_check_status Health check status (0=unknown, 1=up, 2=degraded, 3=down)\n")
	output += fmt.Sprintf("# TYPE health_check_status gauge\n")

	for name, result := range report.Checks {
		checkStatus := 0
		if result.Status == HealthStatusUp {
			checkStatus = 1
		} else if result.Status == HealthStatusDegraded {
			checkStatus = 2
		} else if result.Status == HealthStatusDown {
			checkStatus = 3
		}

		output += fmt.Sprintf("health_check_status{name=\"%s\"} %d\n", name, checkStatus)

		// 检查耗时（毫秒）
		output += fmt.Sprintf("health_check_duration_ms{name=\"%s\"} %d\n", name, result.Duration)
	}

	// 运行时长
	output += fmt.Sprintf("\n# HELP app_uptime_seconds Application uptime in seconds\n")
	output += fmt.Sprintf("# TYPE app_uptime_seconds gauge\n")
	output += fmt.Sprintf("app_uptime_seconds %d\n", report.Uptime)

	return output
}

// getBuildTime 从build info获取构建时间
func getBuildTime(info debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "-buildmode" || setting.Key == "-ldflags" {
			// 尝试从ldflags中提取构建时间
			if setting.Value != "" {
				return setting.Value
			}
		}
	}
	return "unknown"
}

// getGitCommit 从build info获取Git提交
func getGitCommit(info debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}
	return "unknown"
}

// ================================================================================
// 健康检查配置
// ================================================================================

// HealthConfig 健康检查配置
type HealthConfig struct {
	EnableCache       bool          `json:"enable_cache"`        // 是否启用缓存
	CacheDuration     time.Duration `json:"cache_duration"`      // 缓存时长
	Timeout           time.Duration `json:"timeout"`             // 检查超时时间
	DisableCriticalCheck bool       `json:"disable_critical_check"` // 禁用关键检查
}

// DefaultHealthConfig 默认健康检查配置
var DefaultHealthConfig = HealthConfig{
	EnableCache:         false,
	CacheDuration:       10 * time.Second,
	Timeout:             30 * time.Second,
	DisableCriticalCheck: false,
}
