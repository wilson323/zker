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
	"time"

	"github.com/coze-dev/coze-studio/backend/infra/database"
	rediscli "github.com/coze-dev/coze-studio/backend/infra/storage/redis"
)

// ReadinessChecks 就绪检查结果
type ReadinessChecks struct {
	Database      bool `json:"database"`
	Redis         bool `json:"redis"`
	Elasticsearch bool `json:"elasticsearch"`
	MinIO         bool `json:"minio"`
}

// DetailedHealth 详细健康检查结果
type DetailedHealth struct {
	Components map[string]ComponentInfo `json:"components"`
}

// ComponentInfo 组件信息
type ComponentInfo struct {
	Status    string `json:"status"`    // "healthy", "degraded", "unhealthy"
	LatencyMs int64  `json:"latency_ms"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

// CheckReadiness 检查服务就绪状态
func CheckReadiness(ctx context.Context) ReadinessChecks {
	return ReadinessChecks{
		Database:      checkDatabase(ctx),
		Redis:         checkRedis(ctx),
		Elasticsearch: checkElasticsearch(ctx),
		MinIO:         checkMinIO(ctx),
	}
}

// CheckDetailed 详细健康检查
func CheckDetailed(ctx context.Context) DetailedHealth {
	components := make(map[string]ComponentInfo)

	// 检查数据库
	dbStatus := checkDatabaseDetailed(ctx)
	components["database"] = dbStatus

	// 检查 Redis
	redisStatus := checkRedisDetailed(ctx)
	components["redis"] = redisStatus

	// 检查 Elasticsearch
	esStatus := checkElasticsearchDetailed(ctx)
	components["elasticsearch"] = esStatus

	// 检查 MinIO
	minioStatus := checkMinIODetailed(ctx)
	components["minio"] = minioStatus

	return DetailedHealth{
		Components: components,
	}
}

// checkDatabase 检查数据库连接（简单检查）
func checkDatabase(ctx context.Context) bool {
	db := database.GetDB()
	if db == nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	// 设置1秒超时
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	return err == nil
}

// checkDatabaseDetailed 检查数据库连接（详细检查）
func checkDatabaseDetailed(ctx context.Context) ComponentInfo {
	start := time.Now()

	db := database.GetDB()
	if db == nil {
		return ComponentInfo{
			Status:    "unhealthy",
			LatencyMs: 0,
			Message:   "Database client not initialized",
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return ComponentInfo{
			Status:    "unhealthy",
			LatencyMs: 0,
			Error:     err.Error(),
		}
	}

	// 设置2秒超时
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return ComponentInfo{
			Status:    "unhealthy",
			LatencyMs: latency,
			Error:     err.Error(),
		}
	}

	// 执行简单查询验证数据库功能
	var result int
	err = db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		return ComponentInfo{
			Status:    "degraded",
			LatencyMs: latency,
			Error:     fmt.Sprintf("Query failed: %v", err),
		}
	}

	status := "healthy"
	if latency > 100 {
		status = "degraded"
	}

	return ComponentInfo{
		Status:    status,
		LatencyMs: latency,
		Message:   "Database is healthy",
	}
}

// checkRedis 检查 Redis 连接（简单检查）
func checkRedis(ctx context.Context) bool {
	client := rediscli.GetClient()
	if client == nil {
		return false
	}

	// 设置1秒超时
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	return err == nil
}

// checkRedisDetailed 检查 Redis 连接（详细检查）
func checkRedisDetailed(ctx context.Context) ComponentInfo {
	start := time.Now()

	client := rediscli.GetClient()
	if client == nil {
		return ComponentInfo{
			Status:    "unhealthy",
			LatencyMs: 0,
			Message:   "Redis client not initialized",
		}
	}

	// 设置2秒超时
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return ComponentInfo{
			Status:    "unhealthy",
			LatencyMs: latency,
			Error:     err.Error(),
		}
	}

	status := "healthy"
	if latency > 50 {
		status = "degraded"
	}

	return ComponentInfo{
		Status:    status,
		LatencyMs: latency,
		Message:   "Redis is healthy",
	}
}

// checkElasticsearch 检查 Elasticsearch 连接（简单检查）
func checkElasticsearch(ctx context.Context) bool {
	// TODO: 实现 Elasticsearch 健康检查
	// 目前返回 true，待集成 ES 后实现
	return true
}

// checkElasticsearchDetailed 检查 Elasticsearch 连接（详细检查）
func checkElasticsearchDetailed(ctx context.Context) ComponentInfo {
	// TODO: 实现 Elasticsearch 详细健康检查
	// 目前返回 healthy，待集成 ES 后实现
	return ComponentInfo{
		Status:    "healthy",
		LatencyMs: 0,
		Message:   "Elasticsearch check not implemented yet",
	}
}

// checkMinIO 检查 MinIO 连接（简单检查）
func checkMinIO(ctx context.Context) bool {
	// TODO: 实现 MinIO 健康检查
	// 目前返回 true，待集成 MinIO 后实现
	return true
}

// checkMinIODetailed 检查 MinIO 连接（详细检查）
func checkMinIODetailed(ctx context.Context) ComponentInfo {
	// TODO: 实现 MinIO 详细健康检查
	// 目前返回 healthy，待集成 MinIO 后实现
	return ComponentInfo{
		Status:    "healthy",
		LatencyMs: 0,
		Message:   "MinIO check not implemented yet",
	}
}
