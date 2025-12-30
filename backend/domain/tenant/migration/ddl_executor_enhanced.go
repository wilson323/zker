// backend/domain/tenant/migration/ddl_executor_enhanced.go
// DDL执行器增强版 - 优化错误处理和重试机制
package migration

import (
	"context"
	"fmt"
	"strings"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// DDLErrorType DDL错误类型
type DDLErrorType int

const (
	DDLErrorUnknown DDLErrorType = iota
	DDLErrorSyntax
	DDLErrorLockTimeout
	DDLErrorPermission
	DDLErrorDeadlock
	DDLErrorTableExists
)

// DDLExecutorEnhanced 增强版DDL执行器
type DDLExecutorEnhanced struct {
	db         *gorm.DB
	logger     logs.CtxLogger
	maxRetries int
	retryDelay time.Duration
	stats      *DDLStats
}

// DDLStats DDL执行统计
type DDLStats struct {
	TotalCommands int
	SuccessCount  int
	RetryCount    int
	FailureCount  int
	TotalDuration time.Duration
}

// NewDDLExecutorEnhanced 创建增强版DDL执行器
func NewDDLExecutorEnhanced(db *gorm.DB) *DDLExecutorEnhanced {
	return &DDLExecutorEnhanced{
		db:         db,
		logger:     logs.DefaultLogger(),
		maxRetries: 3,
		retryDelay: 1 * time.Second,
		stats:      &DDLStats{},
	}
}

// ExecuteWithRetry 执行DDL并自动重试
func (e *DDLExecutorEnhanced) ExecuteWithRetry(ctx context.Context, sql string) error {
	e.logger.CtxInfof(ctx, "[DDLExecutor] 执行DDL: %s", sql)

	var lastErr error
	for attempt := 0; attempt <= e.maxRetries; attempt++ {
		if attempt > 0 {
			e.logger.CtxWarnf(ctx, "[DDLExecutor] 重试第 %d 次", attempt)
			time.Sleep(e.retryDelay * time.Duration(attempt))
		}

		err := e.db.Exec(sql).Error
		if err == nil {
			e.stats.SuccessCount++
			if attempt > 0 {
				e.stats.RetryCount++
			}
			return nil
		}

		lastErr = err
		errorType := e.classifyError(err)

		// 不可重试的错误
		if errorType == DDLErrorSyntax || errorType == DDLErrorPermission || errorType == DDLErrorTableExists {
			return fmt.Errorf("不可重试错误 (类型: %d): %w", errorType, err)
		}

		// 可重试的错误
		if errorType == DDLErrorLockTimeout || errorType == DDLErrorDeadlock {
			e.logger.CtxWarnf(ctx, "[DDLExecutor] 检测到锁等待，重试中...")
			continue
		}
	}

	return lastErr
}

// classifyError 分类错误类型
func (e *DDLExecutorEnhanced) classifyError(err error) DDLErrorType {
	if err == nil {
		return DDLErrorUnknown
	}

	errMsg := strings.ToLower(err.Error())

	if strings.Contains(errMsg, "syntax error") {
		return DDLErrorSyntax
	}
	if strings.Contains(errMsg, "access denied") {
		return DDLErrorPermission
	}
	if strings.Contains(errMsg, "table already exists") {
		return DDLErrorTableExists
	}
	if strings.Contains(errMsg, "lock wait timeout") {
		return DDLErrorLockTimeout
	}
	if strings.Contains(errMsg, "deadlock") {
		return DDLErrorDeadlock
	}

	return DDLErrorUnknown
}

// GetStats 获取统计信息
func (e *DDLExecutorEnhanced) GetStats() *DDLStats {
	return e.stats
}
