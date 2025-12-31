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

// Package integration 提供集成测试辅助函数和工具
package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// =====================================================
// Testcontainers辅助函数
// =====================================================

// SetupMySQLContainer 启动MySQL测试容器
func SetupMySQLContainer(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	container, err := mysql.RunContainer(ctx,
		testcontainers.WithImage("mysql:8.4.5"),
		mysql.WithDatabase("test_db"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("ready for connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	require.NoError(t, err, "启动MySQL容器失败")

	// 获取连接信息
	connectionString, err := container.ConnectionString(ctx)
	require.NoError(t, err, "获取连接字符串失败")

	// 创建GORM连接
	db, err := gorm.Open(gormMysql.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "连接数据库失败")

	// 清理函数
	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("终止MySQL容器失败: %v", err)
		}
	}

	return db, cleanup
}

// SetupRedisContainer 启动Redis测试容器
func SetupRedisContainer(t *testing.T) (string, func()) {
	ctx := context.Background()

	container, err := redis.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(3*time.Minute),
		),
	)
	require.NoError(t, err, "启动Redis容器失败")

	// 获取连接地址
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "6379")
	require.NoError(t, err)

	connectionString := fmt.Sprintf("%s:%s", host, port.Port())

	// 清理函数
	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("终止Redis容器失败: %v", err)
		}
	}

	return connectionString, cleanup
}

// =====================================================
// 数据库迁移辅助函数
// =====================================================

// RunMigrations 运行数据库迁移
func RunMigrations(t *testing.T, db *gorm.DB, models ...interface{}) {
	err := db.AutoMigrate(models...)
	require.NoError(t, err, "数据库迁移失败")
}

// TruncateTables 清空表数据
func TruncateTables(t *testing.T, db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error
		require.NoError(t, err)
	}
}

// =====================================================
// 测试数据辅助函数
// =====================================================

// SeedTenant 插入测试租户数据
func SeedTenant(t *testing.T, db *gorm.DB, tenantID, name string) {
	sql := `INSERT INTO tenants (tenant_id, tenant_name, tenant_type, status, subscription_tier, created_at, updated_at)
	        VALUES (?, ?, 'enterprise', 'active', 'pro', NOW(), NOW())`
	err := db.Exec(sql, tenantID, name).Error
	require.NoError(t, err)
}

// SeedUser 插入测试用户数据
func SeedUser(t *testing.T, db *gorm.DB, userID, tenantID, username, email string) {
	sql := `INSERT INTO users (user_id, tenant_id, username, email, created_at, updated_at)
	        VALUES (?, ?, ?, ?, NOW(), NOW())`
	err := db.Exec(sql, userID, tenantID, username, email).Error
	require.NoError(t, err)
}

// =====================================================
// HTTP测试辅助函数
// =====================================================

// APITestCase API测试用例
type APITestCase struct {
	Name           string
	Method         string
	URL            string
	Body           interface{}
	ExpectedStatus int
	ExpectedBody   interface{}
	Setup          func() error
	Teardown       func() error
}

// RunAPITests 运行API测试
func RunAPITests(t *testing.T, tests []APITestCase, testFunc func(t *testing.T, tc APITestCase)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Setup != nil {
				require.NoError(t, tt.Setup())
				defer func() {
					if tt.Teardown != nil {
						require.NoError(t, tt.Teardown())
					}
				}()
			}
			testFunc(t, tt)
		})
	}
}

// =====================================================
// 并发测试辅助函数
// =====================================================

// RunConcurrentTest 运行并发测试
func RunConcurrentTest(t *testing.T, concurrency int, testFunc func(t *testing.T, workerID int)) {
	t.Helper()

	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			defer func() {
				done <- true
			}()
			t.Run(fmt.Sprintf("Worker-%d", workerID), func(t *testing.T) {
				testFunc(t, workerID)
			})
		}(i)
	}

	for i := 0; i < concurrency; i++ {
		<-done
	}
}

// =====================================================
// 性能测试辅助函数
// =====================================================

// MeasureExecutionTime 测量执行时间
func MeasureExecutionTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// AssertDuration 断言执行时间在范围内
func AssertDuration(t *testing.T, duration time.Duration, min, max time.Duration) {
	t.Helper()
	if duration < min {
		t.Errorf("执行时间 %v 小于最小期望值 %v", duration, min)
	}
	if duration > max {
		t.Errorf("执行时间 %v 大于最大期望值 %v", duration, max)
	}
}

// =====================================================
// 事务辅助函数
// =====================================================

// WithTransaction 在事务中执行操作
func WithTransaction(t *testing.T, db *gorm.DB, fn func(tx *gorm.DB) error) {
	t.Helper()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		require.NoError(t, err)
	} else {
		require.NoError(t, tx.Commit().Error)
	}
}

// =====================================================
// 等待辅助函数
// =====================================================

// WaitForCondition 等待条件满足
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, interval time.Duration) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("等待条件超时")
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}

// WaitForDBUpdate 等待数据库更新
func WaitForDBUpdate(t *testing.T, db *gorm.DB, query string, args ...interface{}) {
	WaitForCondition(t, func() bool {
		var count int64
		err := db.Raw(query, args...).Count(&count).Error
		return err == nil && count > 0
	}, 10*time.Second, 100*time.Millisecond)
}

// =====================================================
// 断言辅助函数
// =====================================================

// AssertDBCount 断言数据库记录数
func AssertDBCount(t *testing.T, db *gorm.DB, table string, expectedCount int64) {
	t.Helper()

	var count int64
	err := db.Table(table).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

// AssertDBExists 断言记录存在
func AssertDBExists(t *testing.T, db *gorm.DB, query string, args ...interface{}) {
	t.Helper()

	var count int64
	err := db.Raw(query, args...).Count(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0))
}

// AssertDBNotExists 断言记录不存在
func AssertDBNotExists(t *testing.T, db *gorm.DB, query string, args ...interface{}) {
	t.Helper()

	var count int64
	err := db.Raw(query, args...).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
