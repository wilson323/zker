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

package coze

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/model/health"
	healthcheck "github.com/coze-dev/coze-studio/backend/infra/health"
)

// Health 基本健康检查
// @router /api/health [GET]
func Health(ctx context.Context, c *app.RequestContext) {
	// 执行基本健康检查
	checks := healthcheck.CheckReadiness(ctx)

	// 判断服务是否健康
	isHealthy := checks.Database && checks.Redis

	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, &health.HealthResponse{
		Status:    checks,
		Timestamp: time.Now().UnixMilli(),
	})
}

// Live 存活检查
// @router /api/health/live [GET]
func Live(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, &health.LiveResponse{
		Status:  "ok",
		Message: "Service is alive",
	})
}

// Ready 就绪检查
// @router /api/health/ready [GET]
func Ready(ctx context.Context, c *app.RequestContext) {
	checks := healthcheck.CheckReadiness(ctx)

	allReady := checks.Database && checks.Redis && checks.Elasticsearch && checks.MinIO

	status := "ready"
	if !allReady {
		status = "not_ready"
	}

	c.JSON(http.StatusOK, &health.ReadyResponse{
		Status: status,
		Checks: checks,
	})
}

// Detailed 详细健康检查
// @router /api/health/detailed [GET]
func Detailed(ctx context.Context, c *app.RequestContext) {
	detailedHealth := healthcheck.CheckDetailed(ctx)

	// 确定整体状态
	overallStatus := "healthy"
	hasDegraded := false
	hasUnhealthy := false

	for _, component := range detailedHealth.Components {
		switch component.Status {
		case "degraded":
			hasDegraded = true
		case "unhealthy":
			hasUnhealthy = true
		}
	}

	if hasUnhealthy {
		overallStatus = "unhealthy"
	} else if hasDegraded {
		overallStatus = "degraded"
	}

	c.JSON(http.StatusOK, &health.DetailedHealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now().UnixMilli(),
		Components: detailedHealth.Components,
	})
}
