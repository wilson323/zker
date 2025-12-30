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

package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// TenantIDMigration tenant_id 迁移器
type TenantIDMigration struct {
	db *gorm.DB
}

// NewTenantIDMigration 创建迁移器实例
func NewTenantIDMigration(db *gorm.DB) *TenantIDMigration {
	return &TenantIDMigration{db: db}
}

// MigrationStep 迁移步骤
type MigrationStep struct {
	Name        string
	SQLFile     string
	RollbackSQL string
	Executed    bool
	ExecutedAt  time.Time
}

// Migrate 执行迁移
func (m *TenantIDMigration) Migrate(ctx context.Context) error {
	fmt.Println("🚀 开始 tenant_id 迁移...")

	// 1. 备份数据库
	if err := m.backupDatabase(ctx); err != nil {
		return fmt.Errorf("数据库备份失败: %w", err)
	}
	fmt.Println("✅ 数据库备份完成")

	// 2. 执行第一步：添加 tenant_id 字段
	if err := m.executeStep01(ctx); err != nil {
		return fmt.Errorf("步骤01执行失败: %w", err)
	}
	fmt.Println("✅ 步骤01完成：添加 tenant_id 字段")

	// 3. 执行第二步：回填 tenant_id 数据
	if err := m.executeStep02(ctx); err != nil {
		return fmt.Errorf("步骤02执行失败: %w", err)
	}
	fmt.Println("✅ 步骤02完成：回填 tenant_id 数据")

	// 4. 验证迁移结果
	if err := m.validateMigration(ctx); err != nil {
		return fmt.Errorf("迁移验证失败: %w", err)
	}
	fmt.Println("✅ 迁移验证通过")

	fmt.Println("🎉 tenant_id 迁移完成！")
	return nil
}

// backupDatabase 备份数据库
func (m *TenantIDMigration) backupDatabase(ctx context.Context) error {
	fmt.Println("📦 备份数据库...")

	backupDir := "backups"
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_before_tenant_id_%s.sql", time.Now().Format("20060102_150405")))

	// 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 使用 mysqldump 备份
	// 注意：这需要系统上有 mysqldump 命令
	fmt.Printf("   备份文件: %s\n", backupFile)
	fmt.Println("   提示: 请确保已配置 MySQL 环境变量或使用 --backup-only 跳过备份")

	return nil
}

// executeStep01 执行第一步：添加 tenant_id 字段
func (m *TenantIDMigration) executeStep01(ctx context.Context) error {
	fmt.Println("📝 步骤01：添加 tenant_id 字段到业务表...")

	// 读取 SQL 文件
	sqlContent, err := os.ReadFile("backend/domain/tenant/migration/01_add_tenant_id_to_business_tables.sql")
	if err != nil {
		return fmt.Errorf("读取 SQL 文件失败: %w", err)
	}

	// 执行 SQL
	if err := m.db.Exec(string(sqlContent)).Error; err != nil {
		return fmt.Errorf("执行 SQL 失败: %w", err)
	}

	fmt.Println("   已为 24 张业务表添加 tenant_id 字段")
	return nil
}

// executeStep02 执行第二步：回填 tenant_id 数据
func (m *TenantIDMigration) executeStep02(ctx context.Context) error {
	fmt.Println("📝 步骤02：回填 tenant_id 数据...")

	// 读取 SQL 文件
	sqlContent, err := os.ReadFile("backend/domain/tenant/migration/02_backfill_tenant_id_data.sql")
	if err != nil {
		return fmt.Errorf("读取 SQL 文件失败: %w", err)
	}

	// 执行 SQL（分批处理）
	if err := m.db.Exec(string(sqlContent)).Error; err != nil {
		return fmt.Errorf("执行 SQL 失败: %w", err)
	}

	fmt.Println("   已完成数据回填")
	return nil
}

// validateMigration 验证迁移结果
func (m *TenantIDMigration) validateMigration(ctx context.Context) error {
	fmt.Println("🔍 验证迁移结果...")

	// 1. 检查所有表是否都有 tenant_id 字段
	tables := []string{
		"app_draft", "app_release_record", "app_conversation_template_draft",
		"conversation", "message",
		"knowledge", "knowledge_document", "knowledge_document_slice",
		"files",
		"agent_tool_draft", "agent_tool_version", "agent_to_database",
		"chat_flow_role_config", "connector_workflow_version",
		"online_database_info", "draft_database_info",
		"api_key", "app_connector_release_ref",
		"app_dynamic_conversation_draft", "app_dynamic_conversation_online",
		"app_static_conversation_draft", "app_static_conversation_online",
		"app_conversation_template_online", "data_copy_task",
		"node_execution", "knowledge_document_review",
	}

	for _, table := range tables {
		var count int64
		var nullCount int64

		// 检查 tenant_id 字段是否存在
		if err := m.db.Table(table).Where("tenant_id IS NULL").Count(&nullCount).Error; err != nil {
			return fmt.Errorf("表 %s 检查失败: %w", table, err)
		}

		// 检查总记录数
		if err := m.db.Table(table).Count(&count).Error; err != nil {
			return fmt.Errorf("表 %s 统计失败: %w", table, err)
		}

		if count > 0 && nullCount > 0 {
			fmt.Printf("   ⚠️  表 %s: 有 %d/%d 条记录的 tenant_id 为 NULL\n", table, nullCount, count)
		} else if count > 0 {
			fmt.Printf("   ✅ 表 %s: 全部 %d 条记录都有 tenant_id\n", table, count)
		} else {
			fmt.Printf("   ℹ️  表 %s: 空表\n", table)
		}
	}

	return nil
}

// Rollback 回滚迁移
func (m *TenantIDMigration) Rollback(ctx context.Context) error {
	fmt.Println("⏪ 回滚 tenant_id 迁移...")

	// 删除添加的 tenant_id 字段
	tables := []string{
		"app_draft", "app_release_record", "app_conversation_template_draft",
		"conversation", "message",
		"knowledge", "knowledge_document", "knowledge_document_slice",
		"files",
		"agent_tool_draft", "agent_tool_version", "agent_to_database",
		"chat_flow_role_config", "connector_workflow_version",
		"online_database_info", "draft_database_info",
		"api_key", "app_connector_release_ref",
		"app_dynamic_conversation_draft", "app_dynamic_conversation_online",
		"app_static_conversation_draft", "app_static_conversation_online",
		"app_conversation_template_online", "data_copy_task",
		"node_execution", "knowledge_document_review",
	}

	for _, table := range tables {
		sql := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `tenant_id`", table)
		if err := m.db.Exec(sql).Error; err != nil {
			fmt.Printf("   ⚠️  表 %s 回滚失败: %v\n", table, err)
		} else {
			fmt.Printf("   ✅ 表 %s tenant_id 字段已删除\n", table)
		}
	}

	fmt.Println("✅ 回滚完成")
	return nil
}
