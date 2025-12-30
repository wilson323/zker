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

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ==================== 测试上下文 ====================

/**
 * 测试上下文
 *
 * 提供测试所需的共享资源：数据库连接、配置等
 */
type TestContext struct {
	Container testcontainers.Container
	DB        *sql.DB
	Config    TestConfig
	Cleanup   func() error
}

/**
 * 测试配置
 */
type TestConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DSN        string
}

// ==================== 全局测试设置 ====================

/**
 * TestMain - 测试主入口
 *
 * 在所有测试运行前执行一次，用于初始化测试环境
 */
func TestMain(m *testing.M) {
	// TODO: 全局测试初始化
	// 例如：设置日志级别、初始化配置等

	log.Println("集成测试环境初始化...")

	// 运行测试
	code := m.Run()

	// 清理资源
	log.Println("集成测试环境清理...")

	// 退出
	os.Exit(code)
}

// ==================== 工具函数 ====================

/**
 * 创建MySQL测试容器
 *
 * 遵循单一职责原则：只负责创建MySQL容器
 *
 * @param t 测试对象
 * @return *TestContext 测试上下文
 * @return error 错误信息
 */
func SetupMySQLTestContainer(t *testing.T) (*TestContext, error) {
	ctx := context.Background()

	// MySQL 8.4.5 容器配置
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.4.5",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "test123",
			"MYSQL_DATABASE":      "coze_test",
			"MYSQL_USER":          "coze_test",
			"MYSQL_PASSWORD":      "coze_test123",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(5*time.Minute),
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(5*time.Minute),
			wait.ForLog("ready for connections").WithOccurrence(2).WithStartupTimeout(5*time.Minute),
		),
		Cmd: []string{
			"--character-set-server=utf8mb4",
			"--collation-server=utf8mb4_unicode_ci",
			"--default-authentication-plugin=mysql_native_password",
		},
	}

	// 启动容器
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start MySQL container: %w", err)
	}

	// 获取容器信息
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	// 构造DSN
	dsn := fmt.Sprintf("coze_test:coze_test123@tcp(%s:%s)/coze_test?parseTime=true&loc=Local", host, port.Port())

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 创建测试上下文
	config := TestConfig{
		DBHost:     host,
		DBPort:     port.Port(),
		DBUser:     "coze_test",
		DBPassword: "coze_test123",
		DBName:     "coze_test",
		DSN:        dsn,
	}

	testCtx := &TestContext{
		Container: container,
		DB:        db,
		Config:    config,
		Cleanup: func() error {
			// 清理函数
			if err := db.Close(); err != nil {
				return fmt.Errorf("failed to close database: %w", err)
			}
			if err := container.Terminate(ctx); err != nil {
				return fmt.Errorf("failed to terminate container: %w", err)
			}
			return nil
		},
	}

	// 注册清理函数
	t.Cleanup(func() {
		if err := testCtx.Cleanup(); err != nil {
			t.Logf("Cleanup error: %v", err)
		}
	})

	return testCtx, nil
}

/**
 * 等待数据库就绪
 *
 * @param db 数据库连接
 * @param maxAttempts 最大尝试次数
 * @return error 错误信息
 */
func WaitForDatabaseReady(db *sql.DB, maxAttempts int) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		log.Printf("数据库未就绪，重试 %d/%d: %v", i+1, maxAttempts, err)
		time.Sleep(time.Second)
	}
	return fmt.Errorf("数据库在 %d 次尝试后仍未就绪: %w", maxAttempts, err)
}

/**
 * 清理测试数据
 *
 * 遵循测试隔离原则：每个测试后清理数据
 *
 * @param db 数据库连接
 * @param tables 要清理的表名列表
 * @return error 错误信息
 */
func CleanupTestData(db *sql.DB, tables []string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // 安全回滚

	// 禁用外键检查
	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// 清理表数据
	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	// 启用外键检查
	if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS=1"); err != nil {
		return fmt.Errorf("failed to enable foreign key checks: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

/**
 * 执行SQL脚本文件
 *
 * @param db 数据库连接
 * @param scriptFile SQL脚本文件路径
 * @return error 错误信息
 */
func ExecuteSQLScript(db *sql.DB, scriptFile string) error {
	// TODO: 实现SQL脚本文件读取和执行
	return nil
}
