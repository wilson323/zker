// backend/domain/tenant/migration/scheduled_checker.go
// 定时校验任务 - 定期检查数据一致性
package migration

import (
	"context"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// ScheduledChecker 定时校验器
type ScheduledChecker struct {
	db               *gorm.DB
	logger           logs.CtxLogger
	checker          *DataConsistencyChecker
	compensationMgr  *CompensationManager
	checkInterval    time.Duration
	enabled          bool
	autoFix          bool
	lastCheckTime    time.Time
	lastCheckResults map[string]*ConsistencyCheckResult
}

// NewScheduledChecker 创建定时校验器
func NewScheduledChecker(db *gorm.DB, checkInterval time.Duration) *ScheduledChecker {
	return &ScheduledChecker{
		db:               db,
		logger:           logs.DefaultLogger(),
		checker:          NewDataConsistencyChecker(db),
		compensationMgr:  NewCompensationManager(db),
		checkInterval:    checkInterval,
		enabled:          true,
		autoFix:          false, // 默认不自动修复
		lastCheckResults: make(map[string]*ConsistencyCheckResult),
	}
}

// SetAutoFix 设置是否自动修复
func (c *ScheduledChecker) SetAutoFix(enabled bool) {
	c.autoFix = enabled
	if enabled {
		logs.Info("[Checker] 自动修复已启用")
	} else {
		logs.Info("[Checker] 自动修复已禁用")
	}
}

// Enable 启用校验
func (c *ScheduledChecker) Enable() {
	c.enabled = true
	logs.Info("[Checker] 定时校验已启用")
}

// Disable 禁用校验
func (c *ScheduledChecker) Disable() {
	c.enabled = false
	logs.Info("[Checker] 定时校验已禁用")
}

// StartPeriodicCheck 启动定时校验
func (c *ScheduledChecker) StartPeriodicCheck(ctx context.Context, tables []string) {
	if !c.enabled {
		logs.Warn("[Checker] 校验器未启用，跳过定时任务")
		return
	}

	c.logger.CtxInfof(ctx, "[Checker] 启动定时校验任务 (间隔: %v)", c.checkInterval)

	ticker := time.NewTicker(c.checkInterval)
	defer ticker.Stop()

	// 立即执行一次
	c.performCheck(ctx, tables)

	for {
		select {
		case <-ticker.C:
			c.performCheck(ctx, tables)

		case <-ctx.Done():
			logs.Info("[Checker] 定时校验任务已停止")
			return
		}
	}
}

// performCheck 执行校验
func (c *ScheduledChecker) performCheck(ctx context.Context, tables []string) {
	c.logger.CtxInfof(ctx, "[Checker] ========== 开始定时校验 (%s) ==========", time.Now().Format("2006-01-02 15:04:05"))

	startTime := time.Now()

	// 执行一致性检查
	results := c.checker.CheckAllTables(ctx, tables)

	// 保存结果
	c.lastCheckTime = time.Now()
	for tableName, result := range results {
		c.lastCheckResults[tableName] = result
	}

	// 统计
	totalTables := len(results)
	passedTables := 0
	failedTables := 0
	totalIssues := 0

	for _, result := range results {
		if result.AllChecksPassed {
			passedTables++
		} else {
			failedTables++
			totalIssues += len(result.Issues)

			// 自动修复
			if c.autoFix && len(result.Issues) > 0 {
				c.logger.CtxWarnf(ctx, "[Checker] 检测到问题，开始自动修复: %s", result.TableName)
				fixed, err := c.checker.FixInconsistentData(ctx, result.TableName, "default_tenant")
				if err != nil {
					c.logger.CtxErrorf(ctx, "[Checker] 自动修复失败: %s, 错误: %v", result.TableName, err)
				} else {
					c.logger.CtxInfof(ctx, "[Checker] 自动修复完成: %s, 修复 %d 条记录", result.TableName, fixed)
				}
			}
		}
	}

	duration := time.Since(startTime)

	c.logger.CtxInfof(ctx, "[Checker] 校验完成: %d/%d 表通过, 耗时 %v, 发现 %d 个问题",
		passedTables, totalTables, duration, totalIssues)

	// 告警
	if failedTables > 0 {
		c.logger.CtxWarnf(ctx, "[Checker] 警告: %d 张表未通过校验，可能需要人工介入", failedTables)
	}
}

// ManualCheck 手动触发校验
func (c *ScheduledChecker) ManualCheck(ctx context.Context, tables []string) map[string]*ConsistencyCheckResult {
	logs.Info("[Checker] 手动触发校验")
	c.performCheck(ctx, tables)
	return c.lastCheckResults
}

// GetLastCheckResults 获取最后一次校验结果
func (c *ScheduledChecker) GetLastCheckResults() map[string]*ConsistencyCheckResult {
	return c.lastCheckResults
}

// GetLastCheckTime 获取最后一次校验时间
func (c *ScheduledChecker) GetLastCheckTime() time.Time {
	return c.lastCheckTime
}

// GenerateReport 生成校验报告
func (c *ScheduledChecker) GenerateReport(ctx context.Context, tables []string) *ConsistencyReport {
	return c.checker.GenerateConsistencyReport(ctx, tables)
}

// CheckStatus 校验状态
type CheckStatus struct {
	Enabled       bool          `json:"enabled"`
	AutoFix       bool          `json:"auto_fix"`
	LastCheckTime time.Time     `json:"last_check_time"`
	CheckInterval time.Duration `json:"check_interval"`
	NextCheckTime time.Time     `json:"next_check_time"`
	PendingTables int           `json:"pending_tables"`
}

// GetStatus 获取校验器状态
func (c *ScheduledChecker) GetStatus() *CheckStatus {
	status := &CheckStatus{
		Enabled:       c.enabled,
		AutoFix:       c.autoFix,
		LastCheckTime: c.lastCheckTime,
		CheckInterval: c.checkInterval,
	}

	if !c.lastCheckTime.IsZero() {
		status.NextCheckTime = c.lastCheckTime.Add(c.checkInterval)
	}

	status.PendingTables = len(c.lastCheckResults)

	return status
}
