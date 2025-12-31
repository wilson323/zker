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

package testsetup

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MySQLContainer MySQL测试容器
type MySQLContainer struct {
	testcontainers.Container
	host     string
	port     string
	database string
	username string
	password string
}

// SetupMySQLContainer 启动MySQL测试容器
func SetupMySQLContainer(t *testing.T) *MySQLContainer {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.4",
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "test123",
			"MYSQL_DATABASE":      "coze_test",
		},
		WaitingFor: wait.ForLog("ready for connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate MySQL container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get MySQL container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("Failed to get MySQL container port: %v", err)
	}

	return &MySQLContainer{
		Container: container,
		host:      host,
		port:      port.Port(),
		database:  "coze_test",
		username:  "root",
		password:  "test123",
	}
}

// GetDSN 获取MySQL数据源名称
func (mc *MySQLContainer) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mc.username,
		mc.password,
		mc.host,
		mc.port,
		mc.database,
	)
}

// GetHost 获取MySQL主机地址
func (mc *MySQLContainer) GetHost() string {
	return mc.host
}

// GetPort 获取MySQL端口
func (mc *MySQLContainer) GetPort() string {
	return mc.port
}

// RedisContainer Redis测试容器
type RedisContainer struct {
	testcontainers.Container
	host string
	port string
}

// SetupRedisContainer 启动Redis测试容器
func SetupRedisContainer(t *testing.T) *RedisContainer {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "redis:8.0",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Redis container: %v", err)
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate Redis container: %v", err)
		}
	})

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get Redis container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get Redis container port: %v", err)
	}

	return &RedisContainer{
		Container: container,
		host:      host,
		port:      port.Port(),
	}
}

// GetAddr 获取Redis地址
func (rc *RedisContainer) GetAddr() string {
	return fmt.Sprintf("%s:%s", rc.host, rc.port)
}
