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
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// HealthStatus 健康状态
type HealthStatus string

const (
	HealthStatusUp      HealthStatus = "UP"       // 健康
	HealthStatusDown    HealthStatus = "DOWN"     // 不健康
	HealthStatusDegraded HealthStatus = "DEGRADED" // 降级（部分功能不可用）
)

// CheckResult 健康检查结果
type CheckResult struct {
	Name      string                 `json:"name"`      // 检查项名称
	Status    HealthStatus           `json:"status"`    // 健康状态
	Message   string                 `json:"message"`   // 详细消息
	Timestamp int64                  `json:"timestamp"` // 检查时间戳
	Details   map[string]interface{} `json:"details"`  // 额外详情
	Duration  int64                  `json:"duration"`  // 检查耗时（毫秒）
}

// HealthReport 健康报告
type HealthReport struct {
	Status    HealthStatus            `json:"status"`    // 整体健康状态
	Timestamp int64                   `json:"timestamp"` // 报告时间戳
	Checks    map[string]*CheckResult `json:"checks"`    // 所有检查项
	Version   string                  `json:"version"`   // 服务版本
	Uptime    int64                   `json:"uptime"`    // 运行时长（秒）
}

// HealthChecker 健康检查器接口
type HealthChecker interface {
	// Name 返回检查器名称
	Name() string

	// Check 执行健康检查
	Check(ctx context.Context) *CheckResult
}

// HealthService 健康检查服务
type HealthService struct {
	checkers map[string]HealthChecker
	mu       sync.RWMutex
	startTime time.Time
	version   string
}

// NewHealthService 创建健康检查服务
func NewHealthService(version string) *HealthService {
	return &HealthService{
		checkers:  make(map[string]HealthChecker),
		startTime: time.Now(),
		version:   version,
	}
}

// RegisterChecker 注册健康检查器
func (s *HealthService) RegisterChecker(checker HealthChecker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkers[checker.Name()] = checker
	logs.Infof("[HealthService] registered health checker: %s", checker.Name())
}

// UnregisterChecker 注销健康检查器
func (s *HealthService) UnregisterChecker(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.checkers, name)
	logs.Infof("[HealthService] unregistered health checker: %s", name)
}

// Check 执行健康检查
func (s *HealthService) Check(ctx context.Context) *HealthReport {
	s.mu.RLock()
	defer s.mu.RUnlock()

	report := &HealthReport{
		Status:    HealthStatusUp,
		Timestamp: time.Now().UnixMilli(),
		Checks:    make(map[string]*CheckResult),
		Version:   s.version,
		Uptime:    int64(time.Since(s.startTime).Seconds()),
	}

	// 并发执行所有健康检查
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, checker := range s.checkers {
		wg.Add(1)
		go func(name string, checker HealthChecker) {
			defer wg.Done()

			start := time.Now()
			result := checker.Check(ctx)
			result.Duration = time.Since(start).Milliseconds()

			mu.Lock()
			report.Checks[name] = result
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()

	// 计算整体健康状态
	s.calculateOverallStatus(report)

	return report
}

// CheckSingle 执行单个健康检查
func (s *HealthService) CheckSingle(ctx context.Context, name string) (*CheckResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	checker, ok := s.checkers[name]
	if !ok {
		return nil, fmt.Errorf("health checker not found: %s", name)
	}

	start := time.Now()
	result := checker.Check(ctx)
	result.Duration = time.Since(start).Milliseconds()

	return result, nil
}

// calculateOverallStatus 计算整体健康状态
func (s *HealthService) calculateOverallStatus(report *HealthReport) {
	downCount := 0
	degradedCount := 0

	for _, result := range report.Checks {
		switch result.Status {
		case HealthStatusDown:
			downCount++
		case HealthStatusDegraded:
			degradedCount++
		}
	}

	// 如果有任何关键检查项DOWN，整体状态为DOWN
	if s.hasCriticalDown(report) {
		report.Status = HealthStatusDown
		return
	}

	// 如果有任何检查项DEGRADED，整体状态为DEGRADED
	if degradedCount > 0 {
		report.Status = HealthStatusDegraded
		return
	}

	// 如果有非关键检查项DOWN，整体状态为DEGRADED
	if downCount > 0 {
		report.Status = HealthStatusDegraded
		return
	}

	// 所有检查项UP
	report.Status = HealthStatusUp
}

// hasCriticalDown 判断是否有关键检查项DOWN
func (s *HealthService) hasCriticalDown(report *HealthReport) bool {
	criticalCheckers := map[string]bool{
		"mysql":       true,
		"redis":       true,
		"database":    true,
		"main_db":     true,
	}

	for name, result := range report.Checks {
		if criticalCheckers[name] && result.Status == HealthStatusDown {
			return true
		}
	}

	return false
}

// ================================================================================
// MySQL健康检查器
// ================================================================================

// MySQLChecker MySQL健康检查器
type MySQLChecker struct {
	db *gorm.DB
}

// NewMySQLChecker 创建MySQL健康检查器
func NewMySQLChecker(db *gorm.DB) *MySQLChecker {
	return &MySQLChecker{db: db}
}

// Name 返回检查器名称
func (c *MySQLChecker) Name() string {
	return "mysql"
}

// Check 执行健康检查
func (c *MySQLChecker) Check(ctx context.Context) *CheckResult {
	result := &CheckResult{
		Name:      c.Name(),
		Timestamp: time.Now().UnixMilli(),
		Details:   make(map[string]interface{}),
	}

	// 获取SQL DB对象
	sqlDB, err := c.db.DB()
	if err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("Failed to get SQL DB: %v", err)
		return result
	}

	// 1. 检查数据库连接（Ping）
	start := time.Now()
	if err := sqlDB.PingContext(ctx); err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("Database ping failed: %v", err)
		result.Details["error"] = err.Error()
		return result
	}
	pingDuration := time.Since(start).Milliseconds()

	// 2. 检查连接池状态
	stats := sqlDB.Stats()
	result.Details["open_connections"] = stats.OpenConnections
	result.Details["in_use"] = stats.InUse
	result.Details["idle"] = stats.Idle
	result.Details["wait_count"] = stats.WaitCount
	result.Details["wait_duration_ms"] = stats.WaitDuration.Milliseconds()
	result.Details["max_idle_closed"] = stats.MaxIdleClosed
	result.Details["max_lifetime_closed"] = stats.MaxLifetimeClosed
	result.Details["ping_duration_ms"] = pingDuration

	// 3. 执行简单查询验证数据库可读性
	var version string
	if err := c.db.WithContext(ctx).Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("Database query failed: %v", err)
		result.Details["query_error"] = err.Error()
		return result
	}
	result.Details["version"] = version

	// 4. 检查关键表是否存在
	tables := []string{"tenants", "users", "roles"}
	missingTables := []string{}
	for _, table := range tables {
		if !c.tableExists(ctx, table) {
			missingTables = append(missingTables, table)
		}
	}

	if len(missingTables) > 0 {
		result.Status = HealthStatusDegraded
		result.Message = fmt.Sprintf("Missing tables: %v", missingTables)
		result.Details["missing_tables"] = missingTables
		return result
	}

	result.Status = HealthStatusUp
	result.Message = "Database is healthy"
	return result
}

// tableExists 检查表是否存在
func (c *MySQLChecker) tableExists(ctx context.Context, tableName string) bool {
	var count int64
	err := c.db.WithContext(ctx).
		Table("information_schema.tables").
		Where("table_schema = DATABASE() AND table_name = ?", tableName).
		Count(&count).Error
	return err == nil && count > 0
}

// ================================================================================
// Redis健康检查器
// ================================================================================

// RedisChecker Redis健康检查器
type RedisChecker struct {
	client RedisClient
}

// RedisClient Redis客户端接口
type RedisClient interface {
	Ping(ctx context.Context) error
}

// NewRedisChecker 创建Redis健康检查器
func NewRedisChecker(client RedisClient) *RedisChecker {
	return &RedisChecker{client: client}
}

// Name 返回检查器名称
func (c *RedisChecker) Name() string {
	return "redis"
}

// Check 执行健康检查
func (c *RedisChecker) Check(ctx context.Context) *CheckResult {
	result := &CheckResult{
		Name:      c.Name(),
		Timestamp: time.Now().UnixMilli(),
		Details:   make(map[string]interface{}),
	}

	if c.client == nil {
		result.Status = HealthStatusDown
		result.Message = "Redis client is not initialized"
		return result
	}

	// 1. Ping Redis
	start := time.Now()
	if err := c.client.Ping(ctx); err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("Redis ping failed: %v", err)
		result.Details["error"] = err.Error()
		return result
	}
	pingDuration := time.Since(start).Milliseconds()

	result.Status = HealthStatusUp
	result.Message = "Redis is healthy"
	result.Details["ping_duration_ms"] = pingDuration

	return result
}

// ================================================================================
// 数据库通用健康检查器（使用标准库database/sql）
// ================================================================================

// DatabaseChecker 数据库通用健康检查器
type DatabaseChecker struct {
	db    *sql.DB
	name  string
	query string // 用于验证数据库的查询
}

// NewDatabaseChecker 创建数据库健康检查器
func NewDatabaseChecker(db *sql.DB, name, query string) *DatabaseChecker {
	if query == "" {
		query = "SELECT 1"
	}
	return &DatabaseChecker{
		db:    db,
		name:  name,
		query: query,
	}
}

// Name 返回检查器名称
func (c *DatabaseChecker) Name() string {
	return c.name
}

// Check 执行健康检查
func (c *DatabaseChecker) Check(ctx context.Context) *CheckResult {
	result := &CheckResult{
		Name:      c.name,
		Timestamp: time.Now().UnixMilli(),
		Details:   make(map[string]interface{}),
	}

	if c.db == nil {
		result.Status = HealthStatusDown
		result.Message = "Database connection is not initialized"
		return result
	}

	// 1. Ping数据库
	start := time.Now()
	if err := c.db.PingContext(ctx); err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("Database ping failed: %v", err)
		result.Details["error"] = err.Error()
		return result
	}
	pingDuration := time.Since(start).Milliseconds()

	// 2. 执行查询验证数据库可读性
	if err := c.db.QueryRowContext(ctx, c.query).Scan(new(interface{})); err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("Database query failed: %v", err)
		result.Details["error"] = err.Error()
		result.Details["query"] = c.query
		return result
	}
	queryDuration := time.Since(start).Milliseconds()

	// 3. 获取连接池状态
	stats := c.db.Stats()
	result.Details["open_connections"] = stats.OpenConnections
	result.Details["in_use"] = stats.InUse
	result.Details["idle"] = stats.Idle
	result.Details["wait_count"] = stats.WaitCount
	result.Details["wait_duration_ms"] = stats.WaitDuration.Milliseconds()
	result.Details["ping_duration_ms"] = pingDuration
	result.Details["query_duration_ms"] = queryDuration - pingDuration

	result.Status = HealthStatusUp
	result.Message = "Database is healthy"
	return result
}

// ================================================================================
// HTTP健康检查器（检查外部服务）
// ================================================================================

// HTTPChecker HTTP健康检查器
type HTTPChecker struct {
	name    string
	url     string
	timeout time.Duration
	client  HTTPClient
}

// HTTPClient HTTP客户端接口
type HTTPClient interface {
	Get(ctx context.Context, url string) (int, string, error)
}

// NewHTTPChecker 创建HTTP健康检查器
func NewHTTPChecker(name, url string, timeout time.Duration, client HTTPClient) *HTTPChecker {
	return &HTTPChecker{
		name:    name,
		url:     url,
		timeout: timeout,
		client:  client,
	}
}

// Name 返回检查器名称
func (c *HTTPChecker) Name() string {
	return c.name
}

// Check 执行健康检查
func (c *HTTPChecker) Check(ctx context.Context) *CheckResult {
	result := &CheckResult{
		Name:      c.name,
		Timestamp: time.Now().UnixMilli(),
		Details:   make(map[string]interface{}),
	}

	if c.client == nil {
		result.Status = HealthStatusDown
		result.Message = "HTTP client is not initialized"
		return result
	}

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// 发送GET请求
	start := time.Now()
	statusCode, body, err := c.client.Get(ctx, c.url)
	duration := time.Since(start).Milliseconds()

	result.Details["url"] = c.url
	result.Details["duration_ms"] = duration
	result.Details["status_code"] = statusCode

	if err != nil {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("HTTP request failed: %v", err)
		result.Details["error"] = err.Error()
		return result
	}

	// 判断状态码
	if statusCode >= 200 && statusCode < 300 {
		result.Status = HealthStatusUp
		result.Message = "HTTP service is healthy"
	} else if statusCode >= 400 && statusCode < 500 {
		result.Status = HealthStatusDegraded
		result.Message = fmt.Sprintf("HTTP service returned client error: %d", statusCode)
		result.Details["body"] = body
	} else {
		result.Status = HealthStatusDown
		result.Message = fmt.Sprintf("HTTP service returned server error: %d", statusCode)
		result.Details["body"] = body
	}

	return result
}
