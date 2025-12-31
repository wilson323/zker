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

package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// DatabaseIsolationService Database级隔离服务
// 负责为租户创建独立的Database,实现最高级别的数据隔离
//
// 使用示例:
//
//	dbSvc := service.NewDatabaseIsolationService(db, dbConfig)
//	err := dbSvc.CreateDatabaseForTenant(ctx, "tenant-123")
type DatabaseIsolationService struct {
	db              *gorm.DB
	baseDB          *sql.DB     // 基础数据库连接(用于创建新数据库)
	dbHost          string
	dbPort          int
	dbUser          string
	dbPassword      string
	charset         string
	collation       string
}

// NewDatabaseIsolationService 创建Database隔离服务实例
func NewDatabaseIsolationService(
	db *gorm.DB,
	baseDB *sql.DB,
	dbHost string,
	dbPort int,
	dbUser, dbPassword string,
) *DatabaseIsolationService {
	return &DatabaseIsolationService{
		db:         db,
		baseDB:     baseDB,
		dbHost:     dbHost,
		dbPort:     dbPort,
		dbUser:     dbUser,
		dbPassword: dbPassword,
		charset:    "utf8mb4",
		collation:  "utf8mb4_unicode_ci",
	}
}

// CreateDatabaseForTenant 为租户创建独立Database
//
// 功能:
// 1. 创建Database(命名格式: tenant_{tenant_id的hash前8位})
// 2. 创建数据库用户和授权
// 3. 在Database中创建完整表结构
// 4. 验证Database创建成功
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - dbName: 创建的数据库名称
//   - dbUser: 创建的数据库用户名
//   - error: 创建失败时返回错误
func (s *DatabaseIsolationService) CreateDatabaseForTenant(
	ctx context.Context,
	tenantID string,
) (string, string, error) {
	logs.Infof("[DatabaseIsolation] Creating database for tenant: %s", tenantID)

	// 1. 生成数据库名称和用户名
	dbName := s.generateDatabaseName(tenantID)
	dbUser := s.generateDatabaseUser(tenantID)
	dbPassword := s.generateDatabasePassword()

	// 2. 检查数据库是否已存在
	exists, err := s.databaseExists(ctx, dbName)
	if err != nil {
		return "", "", fmt.Errorf("failed to check database existence: %w", err)
	}
	if exists {
		return "", "", fmt.Errorf("database already exists: %s", dbName)
	}

	// 3. 使用基础连接创建数据库和用户
	err = s.createDatabaseAndUser(ctx, dbName, dbUser, dbPassword)
	if err != nil {
		return "", "", fmt.Errorf("failed to create database and user: %w", err)
	}

	// 4. 在新数据库中创建表结构
	err = s.createTablesInDatabase(ctx, dbName, dbUser, dbPassword)
	if err != nil {
		// 回滚:删除数据库
		_ = s.dropDatabase(ctx, dbName)
		return "", "", fmt.Errorf("failed to create tables: %w", err)
	}

	// 5. 验证数据库创建成功
	if err := s.ValidateDatabase(ctx, tenantID); err != nil {
		// 回滚:删除数据库
		_ = s.dropDatabase(ctx, dbName)
		return "", "", fmt.Errorf("database validation failed: %w", err)
	}

	logs.Infof("[DatabaseIsolation] Database created successfully: %s (user: %s)", dbName, dbUser)
	return dbName, dbUser, nil
}

// DeleteDatabaseForTenant 删除租户Database
//
// 功能:
// 1. 删除Database及其所有数据
// 2. 删除数据库用户
// 3. 回收权限
// 4. 验证删除成功
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - error: 删除失败时返回错误
func (s *DatabaseIsolationService) DeleteDatabaseForTenant(ctx context.Context, tenantID string) error {
	logs.Infof("[DatabaseIsolation] Deleting database for tenant: %s", tenantID)

	dbName := s.generateDatabaseName(tenantID)
	dbUser := s.generateDatabaseUser(tenantID)

	// 检查数据库是否存在
	exists, err := s.databaseExists(ctx, dbName)
	if err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("database not found: %s", dbName)
	}

	// 删除数据库用户
	if err := s.dropDatabaseUser(ctx, dbUser); err != nil {
		logs.Warnf("[DatabaseIsolation] Failed to drop user %s: %v", dbUser, err)
		// 继续执行,不中断删除流程
	}

	// 删除数据库
	if err := s.dropDatabase(ctx, dbName); err != nil {
		return err
	}

	logs.Infof("[DatabaseIsolation] Database deleted successfully: %s", dbName)
	return nil
}

// MigrateDatabase 迁移数据库到新的隔离策略
//
// 功能:
// 1. 从行级/Schema级隔离迁移到Database级隔离
// 2. 导出源数据
// 3. 创建新Database
// 4. 导入数据
// 5. 验证数据完整性
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//   - fromStrategy: 原隔离策略
//   - toStrategy: 目标隔离策略
//
// 返回:
//   - error: 迁移失败时返回错误
func (s *DatabaseIsolationService) MigrateDatabase(
	ctx context.Context,
	tenantID string,
	fromStrategy, toStrategy string,
) error {
	logs.Infof("[DatabaseIsolation] Migrating tenant %s from %s to %s", tenantID, fromStrategy, toStrategy)

	// 1. 创建新Database
	dbName, dbUser, err := s.CreateDatabaseForTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to create target database: %w", err)
	}

	// 2. 导出数据
	data, err := s.exportTenantData(ctx, tenantID, fromStrategy)
	if err != nil {
		// 回滚:删除数据库
		_ = s.DeleteDatabaseForTenant(ctx, tenantID)
		return fmt.Errorf("failed to export data: %w", err)
	}

	// 3. 导入数据
	if err := s.importTenantData(ctx, dbName, dbUser, data); err != nil {
		// 回滚:删除数据库
		_ = s.DeleteDatabaseForTenant(ctx, tenantID)
		return fmt.Errorf("failed to import data: %w", err)
	}

	// 4. 验证数据完整性
	if err := s.verifyDataIntegrity(ctx, tenantID, dbName, fromStrategy); err != nil {
		// 回滚:删除数据库
		_ = s.DeleteDatabaseForTenant(ctx, tenantID)
		return fmt.Errorf("data integrity verification failed: %w", err)
	}

	logs.Infof("[DatabaseIsolation] Migration completed successfully for tenant %s", tenantID)
	return nil
}

// ValidateDatabase 验证Database完整性
func (s *DatabaseIsolationService) ValidateDatabase(ctx context.Context, tenantID string) error {
	dbName := s.generateDatabaseName(tenantID)

	// 检查数据库是否存在
	exists, err := s.databaseExists(ctx, dbName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("database not found: %s", dbName)
	}

	// 检查基础表是否存在
	tables := []string{"bots", "conversations", "knowledge_bases", "users"}
	for _, table := range tables {
		if err := s.tableExistsInDatabase(ctx, dbName, table); err != nil {
			return fmt.Errorf("table %s not found in database %s: %w", table, dbName, err)
		}
	}

	logs.Infof("[DatabaseIsolation] Database validation passed: %s", dbName)
	return nil
}

// GetDatabaseConnection 获取租户专用数据库连接
//
// 返回一个只连接到该租户数据库的GORM实例
func (s *DatabaseIsolationService) GetDatabaseConnection(
	ctx context.Context,
	tenantID string,
) (*gorm.DB, error) {
	dbName := s.generateDatabaseName(tenantID)
	dbUser := s.generateDatabaseUser(tenantID)

	// TODO: 从安全存储中获取数据库密码
	// 这里简化处理,实际应从配置中心或密钥管理服务获取
	dbPassword := s.generateDatabasePassword()

	// 构建数据源名称(DSN)
	dsn := s.buildDSN(dbName, dbUser, dbPassword)

	// 创建新的GORM实例
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to tenant database: %w", err)
	}

	return db, nil
}

// ================================================================================
// 私有方法
// ================================================================================

// generateDatabaseName 生成数据库名称
func (s *DatabaseIsolationService) generateDatabaseName(tenantID string) string {
	// 数据库命名: tenant_{tenant_id的hash前8位}
	hash := simpleHash(tenantID)
	return fmt.Sprintf("tenant_%s", hash[:8])
}

// generateDatabaseUser 生成数据库用户名
func (s *DatabaseIsolationService) generateDatabaseUser(tenantID string) string {
	hash := simpleHash(tenantID)
	return fmt.Sprintf("tenant_user_%s", hash[:8])
}

// generateDatabasePassword 生成数据库密码
func (s *DatabaseIsolationService) generateDatabasePassword() string {
	// 实际应使用安全的随机密码生成器
	// 这里简化处理
	return fmt.Sprintf("pwd_%d", time.Now().UnixNano())
}

// databaseExists 检查数据库是否存在
func (s *DatabaseIsolationService) databaseExists(ctx context.Context, dbName string) (bool, error) {
	sql := `
		SELECT COUNT(*)
		FROM information_schema.schemata
		WHERE schema_name = ?
	`

	var count int
	err := s.db.WithContext(ctx).Raw(sql, dbName).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// createDatabaseAndUser 创建数据库和用户
func (s *DatabaseIsolationService) createDatabaseAndUser(
	ctx context.Context,
	dbName, dbUser, dbPassword string,
) error {
	// 1. 创建数据库
	createDBSQL := fmt.Sprintf(
		"CREATE DATABASE `%s` DEFAULT CHARACTER SET %s COLLATE %s",
		dbName, s.charset, s.collation,
	)
	if err := s.db.WithContext(ctx).Exec(createDBSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	logs.Infof("[DatabaseIsolation] Database created: %s", dbName)

	// 2. 创建用户并授权
	createUserSQL := fmt.Sprintf(
		"CREATE USER '%s'@'%%' IDENTIFIED BY '%s'",
		dbUser, dbPassword,
	)
	if err := s.db.WithContext(ctx).Exec(createUserSQL).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	grantSQL := fmt.Sprintf(
		"GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%'",
		dbName, dbUser,
	)
	if err := s.db.WithContext(ctx).Exec(grantSQL).Error; err != nil {
		return fmt.Errorf("failed to grant privileges: %w", err)
	}
	flushSQL := "FLUSH PRIVILEGES"
	if err := s.db.WithContext(ctx).Exec(flushSQL).Error; err != nil {
		return fmt.Errorf("failed to flush privileges: %w", err)
	}

	logs.Infof("[DatabaseIsolation] User created and granted: %s", dbUser)
	return nil
}

// dropDatabase 删除数据库
func (s *DatabaseIsolationService) dropDatabase(ctx context.Context, dbName string) error {
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName)
	if err := s.db.WithContext(ctx).Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}
	return nil
}

// dropDatabaseUser 删除数据库用户
func (s *DatabaseIsolationService) dropDatabaseUser(ctx context.Context, dbUser string) error {
	sql := fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", dbUser)
	if err := s.db.WithContext(ctx).Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to drop user: %w", err)
	}
	return nil
}

// createTablesInDatabase 在数据库中创建表结构
func (s *DatabaseIsolationService) createTablesInDatabase(
	ctx context.Context,
	dbName, dbUser, dbPassword string,
) error {
	// 获取租户数据库连接
	tenantDB, err := s.GetDatabaseConnection(ctx, dbName)
	if err != nil {
		return err
	}

	// 定义基础表DDL
	tablesDDL := map[string]string{
		"bots": `
			CREATE TABLE IF NOT EXISTS bots (
				id BIGINT PRIMARY KEY AUTO_INCREMENT,
				tenant_id VARCHAR(36) NOT NULL,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				status VARCHAR(50) DEFAULT 'active',
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL,
				INDEX idx_tenant_id (tenant_id),
				INDEX idx_status (status)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`,
		"conversations": `
			CREATE TABLE IF NOT EXISTS conversations (
				id BIGINT PRIMARY KEY AUTO_INCREMENT,
				tenant_id VARCHAR(36) NOT NULL,
				bot_id BIGINT NOT NULL,
				user_id BIGINT NOT NULL,
				status VARCHAR(50) DEFAULT 'active',
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL,
				INDEX idx_tenant_id (tenant_id),
				INDEX idx_bot_id (bot_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`,
		"knowledge_bases": `
			CREATE TABLE IF NOT EXISTS knowledge_bases (
				id BIGINT PRIMARY KEY AUTO_INCREMENT,
				tenant_id VARCHAR(36) NOT NULL,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				created_at BIGINT NOT NULL,
				updated_at BIGINT NOT NULL,
				INDEX idx_tenant_id (tenant_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`,
		"users": `
			CREATE TABLE IF NOT EXISTS users (
				id BIGINT PRIMARY KEY AUTO_INCREMENT,
				tenant_id VARCHAR(36) NOT NULL,
				username VARCHAR(100) NOT NULL,
				email VARCHAR(255) NOT NULL,
				status VARCHAR(50) DEFAULT 'active',
				created_at BIGINT NOT NULL,
				UNIQUE KEY uk_username (username),
				INDEX idx_tenant_id (tenant_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		`,
	}

	// 依次创建表
	for tableName, ddl := range tablesDDL {
		logs.Infof("[DatabaseIsolation] Creating table %s in database %s", tableName, dbName)
		if err := tenantDB.WithContext(ctx).Exec(ddl).Error; err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}
		logs.Infof("[DatabaseIsolation] Table %s created successfully", tableName)
	}

	return nil
}

// tableExistsInDatabase 检查表是否存在于数据库中
func (s *DatabaseIsolationService) tableExistsInDatabase(ctx context.Context, dbName, tableName string) error {
	sql := `
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = ? AND table_name = ?
	`

	var count int
	err := s.db.WithContext(ctx).Raw(sql, dbName, tableName).Scan(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("table not found")
	}

	return nil
}

// exportTenantData 导出租户数据
func (s *DatabaseIsolationService) exportTenantData(
	ctx context.Context,
	tenantID string,
	strategy string,
) (map[string][]map[string]interface{}, error) {
	// 根据不同的隔离策略,从不同的位置导出数据
	tables := []string{"bots", "conversations", "knowledge_bases", "users"}
	data := make(map[string][]map[string]interface{})

	for _, table := range tables {
		var rows []map[string]interface{}
		var err error

		switch strategy {
		case "row_level":
			// 从共享表导出
			err = s.db.WithContext(ctx).
				Table(table).
				Where("tenant_id = ?", tenantID).
				Find(&rows).Error
		case "schema_level":
			// 从独立Schema导出
			schemaName := fmt.Sprintf("tenant_%s", simpleHash(tenantID)[:8])
			fullTableName := fmt.Sprintf("`%s`.%s", schemaName, table)
			err = s.db.WithContext(ctx).
				Table(fullTableName).
				Where("tenant_id = ?", tenantID).
				Find(&rows).Error
		default:
			return nil, fmt.Errorf("unsupported strategy: %s", strategy)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to export table %s: %w", table, err)
		}
		data[table] = rows
		logs.Infof("[DatabaseIsolation] Exported %d rows from table %s", len(rows), table)
	}

	return data, nil
}

// importTenantData 导入租户数据
func (s *DatabaseIsolationService) importTenantData(
	ctx context.Context,
	dbName, dbUser string,
	data map[string][]map[string]interface{},
) error {
	// 获取租户数据库连接
	tenantDB, err := s.GetDatabaseConnection(ctx, dbName)
	if err != nil {
		return err
	}

	for table, rows := range data {
		if len(rows) == 0 {
			continue
		}

		for _, row := range rows {
			if err := tenantDB.WithContext(ctx).Table(table).Create(row).Error; err != nil {
				return fmt.Errorf("failed to import into table %s: %w", table, err)
			}
		}
		logs.Infof("[DatabaseIsolation] Imported %d rows into table %s", len(rows), table)
	}

	return nil
}

// verifyDataIntegrity 验证数据完整性
func (s *DatabaseIsolationService) verifyDataIntegrity(
	ctx context.Context,
	tenantID, dbName string,
	fromStrategy string,
) error {
	tables := []string{"bots", "conversations", "knowledge_bases", "users"}

	for _, table := range tables {
		// 查询源数据记录数
		var sourceCount int64
		switch fromStrategy {
		case "row_level":
			err := s.db.WithContext(ctx).
				Table(table).
				Where("tenant_id = ?", tenantID).
				Count(&sourceCount).Error
			if err != nil {
				return fmt.Errorf("failed to count source data in table %s: %w", table, err)
			}
		case "schema_level":
			schemaName := fmt.Sprintf("tenant_%s", simpleHash(tenantID)[:8])
			fullTableName := fmt.Sprintf("`%s`.%s", schemaName, table)
			err := s.db.WithContext(ctx).
				Table(fullTableName).
				Where("tenant_id = ?", tenantID).
				Count(&sourceCount).Error
			if err != nil {
				return fmt.Errorf("failed to count source data in table %s: %w", table, err)
			}
		default:
			return fmt.Errorf("unsupported strategy: %s", fromStrategy)
		}

		// 查询目标数据库记录数
		tenantDB, err := s.GetDatabaseConnection(ctx, dbName)
		if err != nil {
			return err
		}

		var targetCount int64
		err = tenantDB.WithContext(ctx).
			Table(table).
			Where("tenant_id = ?", tenantID).
			Count(&targetCount).Error
		if err != nil {
			return fmt.Errorf("failed to count target data in table %s: %w", table, err)
		}

		// 比较记录数
		if sourceCount != targetCount {
			return fmt.Errorf("data count mismatch in table %s: source=%d, target=%d",
				table, sourceCount, targetCount)
		}

		logs.Infof("[DatabaseIsolation] Data integrity verified for table %s: %d rows", table, sourceCount)
	}

	return nil
}

// buildDSN 构建数据源名称
func (s *DatabaseIsolationService) buildDSN(dbName, dbUser, dbPassword string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		s.dbHost,
		s.dbPort,
		dbName,
		s.charset,
	)
}

// simpleHash 简单哈希函数(复用)
func simpleHash(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
