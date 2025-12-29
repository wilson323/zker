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

package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormresolver "gorm.io/plugin/dbresolver"
)

// ReadWriteSplitConfig 读写分离配置
type ReadWriteSplitConfig struct {
	// 写库配置
	WriteHost string `toml:"write_host"`
	WritePort int    `toml:"write_port"`

	// 读库配置
	ReadHosts []string `toml:"read_hosts"`
	ReadPort  int      `toml:"read_port"`

	// 通用配置
	User     string        `toml:"user"`
	Password string        `toml:"password"`
	Database string        `toml:"database"`
	Timeout  time.Duration `toml:"timeout"`

	// 连接池配置
	MaxOpenConns    int           `toml:"max_open_conns"`
	MaxIdleConns    int           `toml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `toml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `toml:"conn_max_idle_time"`

	// 负载均衡策略
	Policy string `toml:"policy"` // random, round_robin
}

// NewDBWithReadWriteSplit 创建带读写分离的数据库连接
// 使用GORM DBResolver插件实现自动读写路由
func NewDBWithReadWriteSplit(cfg *ReadWriteSplitConfig) (*gorm.DB, error) {
	// 构建写库DSN
	writeDSN := buildDSN(cfg.User, cfg.Password, cfg.WriteHost, cfg.WritePort, cfg.Database)

	// 打开写库连接
	db, err := gorm.Open(mysql.Open(writeDSN), &gorm.Config{
		SkipDefaultTransaction: true, // 禁用默认事务，提升读性能
		PrepareStmt:            true, // 启用预编译语句
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to write database: %w", err)
	}

	// 获取通用DB连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	// 配置读库连接
	var replicas []gorm.Dialector
	for _, host := range cfg.ReadHosts {
		readDSN := buildDSN(cfg.User, cfg.Password, host, cfg.ReadPort, cfg.Database)
		replicas = append(replicas, mysql.Open(readDSN))
	}

	// 选择负载均衡策略
	var policy gormresolver.Policy
	switch cfg.Policy {
	case "round_robin":
		policy = gormresolver.RoundRobinPolicy{}
	default:
		policy = gormresolver.RandomPolicy{}
	}

	// 使用DBResolver插件注册读写分离
	db.Use(gormresolver.Register(gormresolver.Config{
		Replicas: replicas,
		Policy:   policy,
		// 负载均衡策略
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping write database: %w", err)
	}

	return db, nil
}

// buildDSN 构建MySQL DSN
func buildDSN(user, password, host string, port int, dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s",
		user,
		password,
		host,
		port,
		dbname,
		"10s",
	)
}

// ========== 读写分离辅助函数 ==========

// Write 强制使用写库
// 适用场景: 读取刚写入的数据，需要强一致性
func Write(db *gorm.DB) *gorm.DB {
	return db.Clauses(gormresolver.Write)
}

// Read 强制使用读库
// 适用场景: 允许读取延迟数据，提升读性能
func Read(db *gorm.DB) *gorm.DB {
	return db.Clauses(gormresolver.Read)
}

// ========== 使用示例 ==========

// ExampleWriteThenRead 写后读示例
func ExampleWriteThenRead(db *gorm.DB) error {
	// 1. 写入数据 (自动路由到写库)
	tenant := &Tenant{
		TenantID: "tenant123",
		Name:     "Test Tenant",
	}
	if err := db.Create(tenant).Error; err != nil {
		return err
	}

	// 2. 立即读取 (强制使用写库，确保读到最新数据)
	var result Tenant
	if err := Write(db).Where("tenant_id = ?", "tenant123").First(&result).Error; err != nil {
		return err
	}

	// 3. 后续读取 (自动路由到读库，允许复制延迟)
	var list []Tenant
	if err := Read(db).Where("status = ?", "active").Find(&list).Error; err != nil {
		return err
	}

	return nil
}

// ExampleTransactionWithReadWrite 事务中读写示例
func ExampleTransactionWithReadWrite(db *gorm.DB) error {
	// 事务中所有操作都在写库执行
	return db.Transaction(func(tx *gorm.DB) error {
		// 创建租户
		tenant := &Tenant{
			TenantID: "tenant456",
			Name:     "Test Tenant 2",
		}
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}

		// 事务内读取，自动使用写库
		var result Tenant
		return tx.Where("tenant_id = ?", "tenant456").First(&result).Error
	})
}
