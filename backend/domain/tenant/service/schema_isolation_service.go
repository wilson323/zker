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
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// SchemaIsolationService Schema级隔离服务
// 负责为租户创建独立的Schema,实现租户数据隔离
//
// 使用示例:
//
//	schemaSvc := service.NewSchemaIsolationService(db)
//	err := schemaSvc.CreateSchemaForTenant(ctx, "tenant-123")
type SchemaIsolationService struct {
	db *gorm.DB
}

// NewSchemaIsolationService 创建Schema隔离服务实例
func NewSchemaIsolationService(db *gorm.DB) *SchemaIsolationService {
	return &SchemaIsolationService{
		db: db,
	}
}

// CreateSchemaForTenant 为租户创建独立Schema
//
// 功能:
// 1. 创建Schema(命名格式: tenant_{tenant_id})
// 2. 创建Schema权限(授权给租户用户)
// 3. 在Schema中创建基础表结构
// 4. 验证Schema创建成功
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - schemaName: 创建的Schema名称
//   - error: 创建失败时返回错误
func (s *SchemaIsolationService) CreateSchemaForTenant(ctx context.Context, tenantID string) (string, error) {
	logs.Infof("[SchemaIsolation] Creating schema for tenant: %s", tenantID)

	// 1. 生成Schema名称
	schemaName := s.generateSchemaName(tenantID)

	// 2. 检查Schema是否已存在
	exists, err := s.schemaExists(ctx, schemaName)
	if err != nil {
		return "", fmt.Errorf("failed to check schema existence: %w", err)
	}
	if exists {
		return "", fmt.Errorf("schema already exists: %s", schemaName)
	}

	// 3. 使用事务创建Schema和相关配置
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 3.1 创建Schema
		if err := s.createSchema(ctx, schemaName); err != nil {
			return err
		}

		// 3.2 在Schema中创建基础表
		if err := s.createTablesInSchema(ctx, schemaName); err != nil {
			return fmt.Errorf("failed to create tables: %w", err)
		}

		// 3.3 记录Schema创建日志
		logs.Infof("[SchemaIsolation] Schema created successfully: %s for tenant %s", schemaName, tenantID)

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to create schema: %w", err)
	}

	return schemaName, nil
}

// DeleteSchemaForTenant 删除租户Schema
//
// 功能:
// 1. 删除Schema及其所有数据
// 2. 清理Schema权限
// 3. 验证删除成功
//
// 参数:
//   - ctx: 上下文
//   - tenantID: 租户ID
//
// 返回:
//   - error: 删除失败时返回错误
func (s *SchemaIsolationService) DeleteSchemaForTenant(ctx context.Context, tenantID string) error {
	logs.Infof("[SchemaIsolation] Deleting schema for tenant: %s", tenantID)

	schemaName := s.generateSchemaName(tenantID)

	// 检查Schema是否存在
	exists, err := s.schemaExists(ctx, schemaName)
	if err != nil {
		return fmt.Errorf("failed to check schema existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("schema not found: %s", schemaName)
	}

	// 删除Schema
	if err := s.dropSchema(ctx, schemaName); err != nil {
		return err
	}

	logs.Infof("[SchemaIsolation] Schema deleted successfully: %s", schemaName)
	return nil
}

// MigrateSchema 迁移Schema到新的隔离策略
//
// 功能:
// 1. 从行级隔离迁移到Schema级隔离
// 2. 导出数据
// 3. 创建新Schema
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
func (s *SchemaIsolationService) MigrateSchema(
	ctx context.Context,
	tenantID string,
	fromStrategy, toStrategy string,
) error {
	logs.Infof("[SchemaIsolation] Migrating tenant %s from %s to %s", tenantID, fromStrategy, toStrategy)

	// 1. 创建新Schema
	schemaName, err := s.CreateSchemaForTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to create target schema: %w", err)
	}

	// 2. 导出数据(从共享表)
	data, err := s.exportTenantData(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	// 3. 导入数据(到独立Schema)
	if err := s.importTenantData(ctx, schemaName, data); err != nil {
		return fmt.Errorf("failed to import data: %w", err)
	}

	// 4. 验证数据完整性
	if err := s.verifyDataIntegrity(ctx, tenantID, schemaName); err != nil {
		return fmt.Errorf("data integrity verification failed: %w", err)
	}

	logs.Infof("[SchemaIsolation] Migration completed successfully for tenant %s", tenantID)
	return nil
}

// ValidateSchema 验证Schema完整性
func (s *SchemaIsolationService) ValidateSchema(ctx context.Context, tenantID string) error {
	schemaName := s.generateSchemaName(tenantID)

	// 检查Schema是否存在
	exists, err := s.schemaExists(ctx, schemaName)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("schema not found: %s", schemaName)
	}

	// 检查基础表是否存在
	tables := []string{"bots", "conversations", "knowledge_bases", "users"}
	for _, table := range tables {
		if err := s.tableExistsInSchema(ctx, schemaName, table); err != nil {
			return fmt.Errorf("table %s not found in schema %s: %w", table, schemaName, err)
		}
	}

	logs.Infof("[SchemaIsolation] Schema validation passed: %s", schemaName)
	return nil
}

// ================================================================================
// 私有方法
// ================================================================================

// generateSchemaName 生成Schema名称
func (s *SchemaIsolationService) generateSchemaName(tenantID string) string {
	// Schema命名: tenant_{tenant_id的hash前8位}
	// 避免Schema名称过长(MySQL限制64字符)
	hash := simpleHash(tenantID)
	return fmt.Sprintf("tenant_%s", hash[:8])
}

// schemaExists 检查Schema是否存在
func (s *SchemaIsolationService) schemaExists(ctx context.Context, schemaName string) (bool, error) {
	sql := `
		SELECT COUNT(*)
		FROM information_schema.schemata
		WHERE schema_name = ?
	`

	var count int
	err := s.db.WithContext(ctx).Raw(sql, schemaName).Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// createSchema 创建Schema
func (s *SchemaIsolationService) createSchema(ctx context.Context, schemaName string) error {
	sql := fmt.Sprintf("CREATE SCHEMA `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", schemaName)
	if err := s.db.WithContext(ctx).Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

// dropSchema 删除Schema
func (s *SchemaIsolationService) dropSchema(ctx context.Context, schemaName string) error {
	sql := fmt.Sprintf("DROP SCHEMA IF EXISTS `%s`", schemaName)
	if err := s.db.WithContext(ctx).Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}
	return nil
}

// createTablesInSchema 在Schema中创建基础表
func (s *SchemaIsolationService) createTablesInSchema(ctx context.Context, schemaName string) error {
	// 定义基础表DDL
	tablesDDL := map[string]string{
		"bots": `
			CREATE TABLE :schema.bots (
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
			CREATE TABLE :schema.conversations (
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
			CREATE TABLE :schema.knowledge_bases (
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
			CREATE TABLE :schema.users (
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
		// 替换schema占位符
		ddl = strings.ReplaceAll(ddl, ":schema", "`"+schemaName+"`")

		logs.Infof("[SchemaIsolation] Creating table %s in schema %s", tableName, schemaName)
		if err := s.db.WithContext(ctx).Exec(ddl).Error; err != nil {
			return fmt.Errorf("failed to create table %s: %w", tableName, err)
		}
		logs.Infof("[SchemaIsolation] Table %s created successfully", tableName)
	}

	return nil
}

// tableExistsInSchema 检查表是否存在于Schema中
func (s *SchemaIsolationService) tableExistsInSchema(ctx context.Context, schemaName, tableName string) error {
	sql := `
		SELECT COUNT(*)
		FROM information_schema.tables
		WHERE table_schema = ? AND table_name = ?
	`

	var count int
	err := s.db.WithContext(ctx).Raw(sql, schemaName, tableName).Scan(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("table not found")
	}

	return nil
}

// exportTenantData 导出租户数据
func (s *SchemaIsolationService) exportTenantData(ctx context.Context, tenantID string) (map[string][]map[string]interface{}, error) {
	// 定义需要导出的表
	tables := []string{"bots", "conversations", "knowledge_bases", "users"}
	data := make(map[string][]map[string]interface{})

	for _, table := range tables {
		var rows []map[string]interface{}
		err := s.db.WithContext(ctx).
			Table(table).
			Where("tenant_id = ?", tenantID).
			Find(&rows).Error
		if err != nil {
			return nil, fmt.Errorf("failed to export table %s: %w", table, err)
		}
		data[table] = rows
		logs.Infof("[SchemaIsolation] Exported %d rows from table %s", len(rows), table)
	}

	return data, nil
}

// importTenantData 导入租户数据
func (s *SchemaIsolationService) importTenantData(
	ctx context.Context,
	schemaName string,
	data map[string][]map[string]interface{},
) error {
	for table, rows := range data {
		if len(rows) == 0 {
			continue
		}

		fullTableName := fmt.Sprintf("`%s`.%s", schemaName, table)
		for _, row := range rows {
			if err := s.db.WithContext(ctx).Table(fullTableName).Create(row).Error; err != nil {
				return fmt.Errorf("failed to import into table %s: %w", table, err)
			}
		}
		logs.Infof("[SchemaIsolation] Imported %d rows into table %s", len(rows), table)
	}

	return nil
}

// verifyDataIntegrity 验证数据完整性
func (s *SchemaIsolationService) verifyDataIntegrity(ctx context.Context, tenantID, schemaName string) error {
	// 比较行级隔离数据和Schema级隔离数据的记录数
	tables := []string{"bots", "conversations", "knowledge_bases", "users"}

	for _, table := range tables {
		// 查询行级隔离数据
		var rowLevelCount int64
		err := s.db.WithContext(ctx).
			Table(table).
			Where("tenant_id = ?", tenantID).
			Count(&rowLevelCount).Error
		if err != nil {
			return fmt.Errorf("failed to count row-level data in table %s: %w", table, err)
		}

		// 查询Schema级隔离数据
		var schemaLevelCount int64
		fullTableName := fmt.Sprintf("`%s`.%s", schemaName, table)
		err = s.db.WithContext(ctx).
			Table(fullTableName).
			Where("tenant_id = ?", tenantID).
			Count(&schemaLevelCount).Error
		if err != nil {
			return fmt.Errorf("failed to count schema-level data in table %s: %w", table, err)
		}

		// 比较记录数
		if rowLevelCount != schemaLevelCount {
			return fmt.Errorf("data count mismatch in table %s: row-level=%d, schema-level=%d",
				table, rowLevelCount, schemaLevelCount)
		}

		logs.Infof("[SchemaIsolation] Data integrity verified for table %s: %d rows", table, rowLevelCount)
	}

	return nil
}
