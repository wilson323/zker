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

// Package testsetup provides testing utilities and helpers
package testsetup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestConfig 测试配置
type TestConfig struct {
	DatabaseURL string
	RedisAddr   string
}

// LoadTestConfig 加载测试配置
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL: getEnv("TEST_DATABASE_URL", "root:123456@tcp(127.0.0.1:3306)/coze_test?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:   getEnv("TEST_REDIS_ADDR", "127.0.0.1:6379"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// SetupTestDB 创建测试数据库连接
func SetupTestDB(t *testing.T) *gorm.DB {
	config := LoadTestConfig()

	db, err := gorm.Open(mysql.Open(config.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test database")

	sqlDB, err := db.DB()
	require.NoError(t, err)

	// 设置连接池
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)

	// 清理函数
	t.Cleanup(func() {
		sqlDB.Close()
	})

	return db
}

// CleanupTestDB 清理测试数据库
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	// 清理所有测试数据表
	tables := []string{
		"tenants",
		"users",
		"organizations",
		"departments",
		"employees",
		"roles",
		"permissions",
	}

	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE 1=1", table)).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logs.Warnf("Failed to cleanup table %s: %v", table, err)
		}
	}
}

// AssertDBError 断言数据库错误
func AssertDBError(t *testing.T, err error, expectedErrCode string) {
	assert.Error(t, err)
	if expectedErrCode != "" {
		assert.Contains(t, err.Error(), expectedErrCode)
	}
}

// AssertDBNoError 断言无数据库错误
func AssertDBNoError(t *testing.T, err error) {
	assert.NoError(t, err)
}

// AssertDBCount 断言数据库记录数量
func AssertDBCount(t *testing.T, db *gorm.DB, model interface{}, expectedCount int64) {
	var count int64
	err := db.Model(model).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, expectedCount, count)
}

// TruncateTable 截断表
func TruncateTable(t *testing.T, db *gorm.DB, tableName string) {
	err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
	require.NoError(t, err)
}

// SetupTestContext 设置测试上下文
func SetupTestContext(t *testing.T) context.Context {
	return context.Background()
}

// SkipIfDBUnavailable 如果数据库不可用则跳过测试
func SkipIfDBUnavailable(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Skip("Database not available")
	}

	err = sqlDB.Ping()
	if err != nil {
		t.Skipf("Database not available: %v", err)
	}
}
