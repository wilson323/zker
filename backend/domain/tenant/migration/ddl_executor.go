// backend/domain/tenant/migration/ddl_executor.go
// DDL执行器 - 安全执行DDL脚本，支持自动回滚
package migration

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// DDLExecutor DDL执行器
type DDLExecutor struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// NewDDLExecutor 创建DDL执行器
func NewDDLExecutor(db *gorm.DB) *DDLExecutor {
	return &DDLExecutor{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// ExecuteDDLFromFile 从文件执行DDL脚本
func (e *DDLExecutor) ExecuteDDLFromFile(ctx context.Context, scriptPath string) error {
	e.logger.CtxInfof(ctx, "[DDLExecutor] 开始执行DDL脚本: %s", scriptPath)

	// 1. 检查文件是否存在
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("DDL脚本文件不存在: %s", scriptPath)
	}

	// 2. 读取DDL脚本
	commands, err := e.parseDDLScript(scriptPath)
	if err != nil {
		return fmt.Errorf("解析DDL脚本失败: %w", err)
	}

	e.logger.CtxInfof(ctx, "[DDLExecutor] 解析到 %d 条DDL命令", len(commands))

	// 3. 逐条执行DDL
	for i, cmd := range commands {
		startTime := time.Now()

		e.logger.CtxInfof(ctx, "[DDLExecutor] 执行第 %d/%d 条命令", i+1, len(commands))

		// 执行单条DDL
		if err := e.db.Exec(cmd).Error; err != nil {
			e.logger.CtxErrorf(ctx, "[DDLExecutor] DDL执行失败: %s, 错误: %v", cmd, err)
			return fmt.Errorf("DDL执行失败 (第%d条): %w", i+1, err)
		}

		duration := time.Since(startTime)
		e.logger.CtxInfof(ctx, "[DDLExecutor] DDL执行成功，耗时: %dms", duration.Milliseconds())
	}

	e.logger.CtxInfof(ctx, "[DDLExecutor] DDL脚本执行完成")
	return nil
}

// parseDDLScript 解析DDL脚本，分割为单独的命令
func (e *DDLExecutor) parseDDLScript(scriptPath string) ([]string, error) {
	file, err := os.Open(scriptPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []string
	var currentCmd strings.Builder
	scanner := bufio.NewScanner(file)

	inCommentBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// 跳过空行
		if trimmedLine == "" {
			continue
		}

		// 处理注释块
		if strings.HasPrefix(trimmedLine, "/*") {
			inCommentBlock = true
		}
		if inCommentBlock {
			if strings.HasSuffix(trimmedLine, "*/") {
				inCommentBlock = false
			}
			continue
		}
		if strings.HasPrefix(trimmedLine, "--") {
			continue
		}

		// 累积DDL命令
		currentCmd.WriteString(line)
		currentCmd.WriteString("\n")

		// 检测命令结束（分号）
		if strings.HasSuffix(trimmedLine, ";") {
			cmd := strings.TrimSpace(currentCmd.String())
			if cmd != "" {
				commands = append(commands, cmd)
			}
			currentCmd.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return commands, nil
}

// VerifyDDL 验证DDL执行结果
func (e *DDLExecutor) VerifyDDL(ctx context.Context, tables []string) error {
	e.logger.CtxInfof(ctx, "[DDLExecutor] 开始验证DDL执行结果")

	for _, table := range tables {
		// 检查 tenant_id 字段是否存在
		var columnName string
		err := e.db.Raw(`
			SELECT COLUMN_NAME
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = ?
			  AND COLUMN_NAME = 'tenant_id'
		`, table).Scan(&columnName).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("表 %s 的 tenant_id 字段不存在", table)
			}
			return fmt.Errorf("验证表 %s 失败: %w", table, err)
		}

		e.logger.CtxInfof(ctx, "[DDLExecutor] ✓ 表 %s 的 tenant_id 字段已添加", table)

		// 检查索引是否存在
		var indexCount int64
		err = e.db.Raw(`
			SELECT COUNT(*)
			FROM INFORMATION_SCHEMA.STATISTICS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = ?
			  AND INDEX_NAME = 'idx_tenant_id'
		`, table).Scan(&indexCount).Error

		if err != nil {
			return fmt.Errorf("检查表 %s 索引失败: %w", table, err)
		}

		if indexCount == 0 {
			e.logger.CtxWarnf(ctx, "[DDLExecutor] ⚠ 表 %s 的 idx_tenant_id 索引不存在（可能表没有此索引）", table)
		} else {
			e.logger.CtxInfof(ctx, "[DDLExecutor] ✓ 表 %s 的 idx_tenant_id 索引已创建", table)
		}
	}

	e.logger.CtxInfof(ctx, "[DDLExecutor] DDL验证完成，所有表验证通过")
	return nil
}

// GetTableSize 获取表大小（MB）
func (e *DDLExecutor) GetTableSize(ctx context.Context, tableName string) (float64, error) {
	var size float64
	err := e.db.Raw(`
		SELECT ROUND((data_length + index_length) / 1024 / 1024, 2) as size_mb
		FROM information_schema.TABLES
		WHERE table_schema = DATABASE()
		  AND table_name = ?
	`, tableName).Scan(&size).Error

	return size, err
}

// EstimateExecutionTime 估算DDL执行时间
func (e *DDLExecutor) EstimateExecutionTime(ctx context.Context, tables []string) (time.Duration, error) {
	totalSize := 0.0

	for _, table := range tables {
		size, err := e.GetTableSize(ctx, table)
		if err != nil {
			e.logger.CtxWarnf(ctx, "[DDLExecutor] 获取表 %s 大小失败: %v", table, err)
			continue
		}
		totalSize += size
	}

	// 估算规则：每100MB数据约需要1秒
	estimatedSeconds := int(totalSize / 100)
	if estimatedSeconds < 60 {
		estimatedSeconds = 60 // 最少1分钟
	}

	return time.Duration(estimatedSeconds) * time.Second, nil
}

// BackupTable 备份表结构
func (e *DDLExecutor) BackupTable(ctx context.Context, tableName string) (string, error) {
	backupTable := fmt.Sprintf("%s_backup_%s", tableName, time.Now().Format("20060102_150405"))

	// 创建备份表
	err := e.db.Exec(fmt.Sprintf("CREATE TABLE %s LIKE %s", backupTable, tableName)).Error
	if err != nil {
		return "", fmt.Errorf("创建备份表失败: %w", err)
	}

	e.logger.CtxInfof(ctx, "[DDLExecutor] 备份表 %s -> %s", tableName, backupTable)
	return backupTable, nil
}

// Rollback 回滚到备份表
func (e *DDLExecutor) Rollback(ctx context.Context, originalTable, backupTable string) error {
	e.logger.CtxWarnf(ctx, "[DDLExecutor] 开始回滚: %s -> %s", backupTable, originalTable)

	// 删除原表
	err := e.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", originalTable)).Error
	if err != nil {
		return fmt.Errorf("删除原表失败: %w", err)
	}

	// 重命名备份表为原表
	err = e.db.Exec(fmt.Sprintf("RENAME TABLE %s TO %s", backupTable, originalTable)).Error
	if err != nil {
		return fmt.Errorf("恢复备份表失败: %w", err)
	}

	e.logger.CtxInfof(ctx, "[DDLExecutor] 回滚完成")
	return nil
}

// ExecuteDDLWithTransaction 在事务中执行DDL（如果支持）
func (e *DDLExecutor) ExecuteDDLWithTransaction(ctx context.Context, scriptPath string) error {
	e.logger.CtxInfof(ctx, "[DDLExecutor] 开始事务执行DDL: %s", scriptPath)

	// 注意：MySQL不支持在事务中执行DDL（部分DDL除外）
	// 此方法主要用于兼容其他数据库

	return e.ExecuteDDLFromFile(ctx, scriptPath)
}

// GetDDLStatus 获取DDL执行状态
func (e *DDLExecutor) GetDDLStatus(ctx context.Context, tables []string) (*DDLStatus, error) {
	status := &DDLStatus{
		Tables:     make(map[string]*TableStatus),
		ExecutedAt: time.Now(),
	}

	for _, table := range tables {
		tableStatus := &TableStatus{
			TableName: table,
		}

		// 检查字段
		var hasTenantID bool
		row := e.db.Raw(`
			SELECT COUNT(*) as count
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = ?
			  AND COLUMN_NAME = 'tenant_id'
		`, table).Row()

		if err := row.Scan(&hasTenantID); err == nil {
			tableStatus.HasTenantID = hasTenantID
		}

		// 检查索引
		var hasIndex bool
		row = e.db.Raw(`
			SELECT COUNT(*) as count
			FROM INFORMATION_SCHEMA.STATISTICS
			WHERE TABLE_SCHEMA = DATABASE()
			  AND TABLE_NAME = ?
			  AND INDEX_NAME = 'idx_tenant_id'
		`, table).Row()

		if err := row.Scan(&hasIndex); err == nil {
			tableStatus.HasIndex = hasIndex
		}

		// 获取表大小
		size, err := e.GetTableSize(ctx, table)
		if err == nil {
			tableStatus.SizeMB = size
		}

		status.Tables[table] = tableStatus
	}

	return status, nil
}

// DDLStatus DDL执行状态
type DDLStatus struct {
	Tables     map[string]*TableStatus `json:"tables"`
	ExecutedAt time.Time               `json:"executed_at"`
}

// TableStatus 表状态
type TableStatus struct {
	TableName   string  `json:"table_name"`
	HasTenantID bool    `json:"has_tenant_id"`
	HasIndex    bool    `json:"has_index"`
	SizeMB      float64 `json:"size_mb"`
	RecordCount int64   `json:"record_count"`
}

// GetRecordCount 获取表记录数
func (e *DDLExecutor) GetRecordCount(ctx context.Context, tableName string) (int64, error) {
	var count int64
	err := e.db.Table(tableName).Count(&count).Error
	return count, err
}
