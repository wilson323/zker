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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/billing/entity"
)

// ============================================================
// 计费系统健康检查
// ============================================================

// HealthChecker 健康检查器
type HealthChecker struct {
	db           *gorm.DB
	redisClient  interface{} // 可选的Redis客户端
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(db *gorm.DB, redisClient interface{}) *HealthChecker {
	return &HealthChecker{
		db:          db,
		redisClient: redisClient,
	}
}

// CheckResult 健康检查结果
type CheckResult struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"` // healthy, degraded, unhealthy
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// HealthStatus 系统健康状态
type HealthStatus struct {
	Status      string                 `json:"status"`
	Timestamp   time.Time              `json:"timestamp"`
	Version     string                 `json:"version"`
	Checks      map[string]CheckResult `json:"checks"`
	Uptime      time.Duration          `json:"uptime"`
}

// Check 执行完整健康检查
func (c *HealthChecker) Check(ctx context.Context) *HealthStatus {
	start := time.Now()
	checks := make(map[string]CheckResult)

	// 数据库健康检查
	checks["database"] = c.checkDatabase(ctx)

	// Redis健康检查 (如果配置)
	if c.redisClient != nil {
		checks["redis"] = c.checkRedis(ctx)
	}

	// Token日志表健康检查
	checks["token_log_table"] = c.checkTokenLogTable(ctx)

	// 预算配置表健康检查
	checks["budget_config_table"] = c.checkBudgetConfigTable(ctx)

	// 汇总状态
	status := c.calculateOverallStatus(checks)

	return &HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Checks:    checks,
		Uptime:    time.Since(start),
	}
}

// checkDatabase 检查数据库连接
func (c *HealthChecker) checkDatabase(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Name: "database",
	}

	// 执行简单的查询测试连接
	var count int64
	err := c.db.WithContext(ctx).Table("token_usage_logs").
		Where("created_at > ?", time.Now().AddDate(0, 0, -1)).
		Count(&count).Error

	result.Duration = time.Since(start)

	if err != nil {
		result.Status = "unhealthy"
		result.Error = fmt.Sprintf("database connection failed: %v", err)
		return result
	}

	result.Status = "healthy"
	result.Details = map[string]interface{}{
		"recent_records": count,
		"latency_ms":     result.Duration.Milliseconds(),
	}

	return result
}

// checkRedis 检查Redis连接
func (c *HealthChecker) checkRedis(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Name: "redis",
	}

	// TODO: 实现Redis PING 命令
	// 示例实现:
	// if err := c.redisClient.Ping(ctx).Err(); err != nil {
	//     result.Status = "unhealthy"
	//     result.Error = err.Error()
	//     return result
	// }

	result.Duration = time.Since(start)
	result.Status = "healthy"
	result.Details = map[string]interface{}{
		"latency_ms": result.Duration.Milliseconds(),
	}

	return result
}

// checkTokenLogTable 检查Token日志表健康度
func (c *HealthChecker) checkTokenLogTable(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Name: "token_log_table",
	}

	// 检查表是否存在
	if !c.db.Migrator().HasTable(&entity.TokenLog{}) {
		result.Status = "unhealthy"
		result.Error = "token_log table does not exist"
		return result
	}

	// 检查最近的记录
	var latestLog entity.TokenLog
	err := c.db.WithContext(ctx).
		Order("created_at DESC").
		First(&latestLog).Error

	result.Duration = time.Since(start)

	if err != nil {
		// 如果没有记录,不算错误
		if err == gorm.ErrRecordNotFound {
			result.Status = "healthy"
			result.Details = map[string]interface{}{
				"latest_record_age_ms": 0,
				"message":             "no records yet",
			}
			return result
		}

		result.Status = "unhealthy"
		result.Error = fmt.Sprintf("failed to query token_log: %v", err)
		return result
	}

	// 检查最新记录的年龄
	age := time.Since(latestLog.CreatedAt)
	result.Status = "healthy"
	result.Details = map[string]interface{}{
		"latest_record_age_ms": age.Milliseconds(),
		"latest_record_id":     latestLog.ID,
		"latest_record_date":   latestLog.CreatedAt,
	}

	// 如果最新记录超过24小时,标记为degraded
	if age > 24*time.Hour {
		result.Status = "degraded"
		result.Details["warning"] = "latest record is older than 24 hours"
	}

	return result
}

// checkBudgetConfigTable 检查预算配置表健康度
func (c *HealthChecker) checkBudgetConfigTable(ctx context.Context) CheckResult {
	start := time.Now()
	result := CheckResult{
		Name: "budget_config_table",
	}

	// 检查表是否存在
	if !c.db.Migrator().HasTable(&entity.BudgetConfig{}) {
		result.Status = "unhealthy"
		result.Error = "budget_config table does not exist"
		return result
	}

	// 统计活跃的预算配置
	var count int64
	err := c.db.WithContext(ctx).
		Model(&entity.BudgetConfig{}).
		Where("is_active = ?", true).
		Count(&count).Error

	result.Duration = time.Since(start)

	if err != nil {
		result.Status = "unhealthy"
		result.Error = fmt.Sprintf("failed to query budget_config: %v", err)
		return result
	}

	result.Status = "healthy"
	result.Details = map[string]interface{}{
		"active_budgets": count,
	}

	return result
}

// calculateOverallStatus 计算总体健康状态
func (c *HealthChecker) calculateOverallStatus(checks map[string]CheckResult) string {
	hasUnhealthy := false
	hasDegraded := false

	for _, check := range checks {
		switch check.Status {
		case "unhealthy":
			hasUnhealthy = true
		case "degraded":
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return "unhealthy"
	} else if hasDegraded {
		return "degraded"
	}
	return "healthy"
}

// CheckLiveness 活性检查 (快速检查)
func (c *HealthChecker) CheckLiveness(ctx context.Context) error {
	// 简单的数据库连接检查
	return c.db.WithContext(ctx).Exec("SELECT 1").Error
}

// CheckReadiness 就绪检查 (完整检查)
func (c *HealthChecker) CheckReadiness(ctx context.Context) error {
	status := c.Check(ctx)

	if status.Status == "unhealthy" {
		return fmt.Errorf("system is unhealthy: %d checks failed", c.countFailedChecks(status.Checks))
	}

	return nil
}

// countFailedCounts 统计失败的检查数
func (c *HealthChecker) countFailedChecks(checks map[string]CheckResult) int {
	count := 0
	for _, check := range checks {
		if check.Status == "unhealthy" {
			count++
		}
	}
	return count
}

// ============================================================
// 指标收集器
// ============================================================

// MetricsCollector 指标收集器
type MetricsCollector struct {
	healthChecker *HealthChecker
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector(healthChecker *HealthChecker) *MetricsCollector {
	return &MetricsCollector{
		healthChecker: healthChecker,
	}
}

// CollectMetrics 收集健康检查指标
func (c *MetricsCollector) CollectMetrics(ctx context.Context) map[string]interface{} {
	status := c.healthChecker.Check(ctx)

	metrics := make(map[string]interface{})

	// 整体状态
	metrics["health_status"] = status.Status
	metrics["uptime_seconds"] = status.Uptime.Seconds()

	// 各项检查结果
	for name, check := range status.Checks {
		metrics[fmt.Sprintf("check.%s.status", name)] = check.Status
		metrics[fmt.Sprintf("check.%s.duration_ms", name)] = check.Duration.Milliseconds()

		// 添加详细信息
		for k, v := range check.Details {
			metrics[fmt.Sprintf("check.%s.%s", name, k)] = v
		}
	}

	return metrics
}
