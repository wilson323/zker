// backend/domain/tenant/migration/session_user_validator.go
// Session/User 表迁移验证器
package migration

import (
	"context"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// SessionUserValidator Session/User表迁移验证器
type SessionUserValidator struct {
	db     *gorm.DB
	logger logs.CtxLogger
}

// NewSessionUserValidator 创建验证器
func NewSessionUserValidator(db *gorm.DB) *SessionUserValidator {
	return &SessionUserValidator{
		db:     db,
		logger: logs.DefaultLogger(),
	}
}

// ValidationResult 验证结果
type ValidationResult struct {
	Item      string `json:"item"`
	Passed    bool   `json:"passed"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
	Timestamp string `json:"timestamp"`
}

// ValidationReport 验证报告
type ValidationReport struct {
	ValidatedAt time.Time           `json:"validated_at"`
	TotalItems  int                 `json:"total_items"`
	PassedItems int                 `json:"passed_items"`
	FailedItems int                 `json:"failed_items"`
	PassRate    float64             `json:"pass_rate"`
	Results     []*ValidationResult `json:"results"`
	Status      string              `json:"status"` // passed, failed, warning
}

// ValidateMigration 执行完整的迁移验证
func (v *SessionUserValidator) ValidateMigration(ctx context.Context) (*ValidationReport, error) {
	v.logger.CtxInfof(ctx, "[SessionUserValidator] 开始迁移验证...")

	report := &ValidationReport{
		ValidatedAt: time.Now(),
		Results:     make([]*ValidationResult, 0),
		Status:      "passed",
	}

	// 1. 验证 users 表 tenant_id 字段
	result1 := v.validateUsersTableSchema(ctx)
	report.Results = append(report.Results, result1)

	// 2. 验证 session 表 tenant_id 字段
	result2 := v.validateSessionTableSchema(ctx)
	report.Results = append(report.Results, result2)

	// 3. 验证 user_tenant 表存在
	result3 := v.validateUserTenantTable(ctx)
	report.Results = append(report.Results, result3)

	// 4. 验证 users 表数据完整性
	result4 := v.validateUsersData(ctx)
	report.Results = append(report.Results, result4)

	// 5. 验证 session 表数据完整性
	result5 := v.validateSessionData(ctx)
	report.Results = append(report.Results, result5)

	// 6. 验证 user_tenant 关联完整性
	result6 := v.validateUserTenantRelations(ctx)
	report.Results = append(report.Results, result6)

	// 7. 验证索引是否创建
	result7 := v.validateIndexes(ctx)
	report.Results = append(report.Results, result7)

	// 8. 验证孤儿数据
	result8 := v.validateOrphanData(ctx)
	report.Results = append(report.Results, result8)

	// 统计结果
	for _, result := range report.Results {
		if result.Passed {
			report.PassedItems++
		} else {
			report.FailedItems++
			report.Status = "failed"
		}
	}

	report.TotalItems = len(report.Results)
	report.PassRate = float64(report.PassedItems) / float64(report.TotalItems) * 100

	if report.Status == "failed" {
		v.logger.CtxErrorf(ctx, "[SessionUserValidator] 验证失败: 通过率=%.2f%% (%d/%d)",
			report.PassRate, report.PassedItems, report.TotalItems)
		return report, fmt.Errorf("迁移验证失败")
	}

	v.logger.CtxInfof(ctx, "[SessionUserValidator] 验证通过: 通过率=%.2f%% (%d/%d)",
		report.PassRate, report.PassedItems, report.TotalItems)

	return report, nil
}

// validateUsersTableSchema 验证 users 表结构
func (v *SessionUserValidator) validateUsersTableSchema(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "users.tenant_id字段",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var count int64
	err := v.db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'user'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&count).Error

	if err != nil {
		result.Passed = false
		result.Message = "查询字段失败"
		result.Details = err.Error()
		return result
	}

	if count == 0 {
		result.Passed = false
		result.Message = "tenant_id字段不存在"
		return result
	}

	result.Passed = true
	result.Message = "tenant_id字段存在"
	return result
}

// validateSessionTableSchema 验证 session 表结构
func (v *SessionUserValidator) validateSessionTableSchema(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "session.tenant_id字段",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var count int64
	err := v.db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'session'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&count).Error

	if err != nil {
		result.Passed = false
		result.Message = "查询字段失败"
		result.Details = err.Error()
		return result
	}

	if count == 0 {
		result.Passed = false
		result.Message = "tenant_id字段不存在"
		return result
	}

	result.Passed = true
	result.Message = "tenant_id字段存在"
	return result
}

// validateUserTenantTable 验证 user_tenant 表
func (v *SessionUserValidator) validateUserTenantTable(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "user_tenant表",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var count int64
	err := v.db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'user_tenant'
	`).Scan(&count).Error

	if err != nil {
		result.Passed = false
		result.Message = "查询表失败"
		result.Details = err.Error()
		return result
	}

	if count == 0 {
		result.Passed = false
		result.Message = "user_tenant表不存在"
		return result
	}

	result.Passed = true
	result.Message = "user_tenant表存在"
	return result
}

// validateUsersData 验证 users 表数据
func (v *SessionUserValidator) validateUsersData(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "users表数据完整性",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var stats struct {
		Total    int64
		NullCnt  int64
		HasValue int64
	}

	err := v.db.Raw(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_cnt,
			SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_value
		FROM user
	`).Scan(&stats).Error

	if err != nil {
		result.Passed = false
		result.Message = "查询数据失败"
		result.Details = err.Error()
		return result
	}

	if stats.NullCnt > 0 {
		result.Passed = false
		result.Message = fmt.Sprintf("存在%d条记录的tenant_id为NULL", stats.NullCnt)
		result.Details = fmt.Sprintf("总记录=%d, NULL=%d, 有值=%d", stats.Total, stats.NullCnt, stats.HasValue)
		return result
	}

	result.Passed = true
	result.Message = "所有记录都有tenant_id"
	result.Details = fmt.Sprintf("总记录=%d", stats.Total)
	return result
}

// validateSessionData 验证 session 表数据
func (v *SessionUserValidator) validateSessionData(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "session表数据完整性",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var stats struct {
		Total    int64
		NullCnt  int64
		HasValue int64
	}

	err := v.db.Raw(`
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_cnt,
			SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_value
		FROM session
	`).Scan(&stats).Error

	if err != nil {
		result.Passed = false
		result.Message = "查询数据失败"
		result.Details = err.Error()
		return result
	}

	if stats.NullCnt > 0 {
		result.Passed = false
		result.Message = fmt.Sprintf("存在%d条记录的tenant_id为NULL", stats.NullCnt)
		result.Details = fmt.Sprintf("总记录=%d, NULL=%d, 有值=%d", stats.Total, stats.NullCnt, stats.HasValue)
		return result
	}

	result.Passed = true
	result.Message = "所有记录都有tenant_id"
	result.Details = fmt.Sprintf("总记录=%d", stats.Total)
	return result
}

// validateUserTenantRelations 验证 user_tenant 关联
func (v *SessionUserValidator) validateUserTenantRelations(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "user_tenant关联完整性",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var stats struct {
		UserCount     int64
		RelationCount int64
	}

	// 获取用户数
	v.db.Table("user").Count(&stats.UserCount)

	// 获取关联数
	v.db.Raw(`
		SELECT COUNT(DISTINCT user_id)
		FROM user_tenant
		WHERE deleted_at IS NULL
	`).Scan(&stats.RelationCount)

	if stats.RelationCount < stats.UserCount {
		result.Passed = false
		result.Message = "部分用户没有租户关联"
		result.Details = fmt.Sprintf("用户数=%d, 关联数=%d, 缺失=%d",
			stats.UserCount, stats.RelationCount, stats.UserCount-stats.RelationCount)
		return result
	}

	result.Passed = true
	result.Message = "所有用户都有租户关联"
	result.Details = fmt.Sprintf("用户数=%d, 关联数=%d", stats.UserCount, stats.RelationCount)
	return result
}

// validateIndexes 验证索引
func (v *SessionUserValidator) validateIndexes(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "tenant_id索引",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 检查 users 表索引
	var userIndexCount int64
	v.db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'user'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&userIndexCount)

	// 检查 session 表索引
	var sessionIndexCount int64
	v.db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'session'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&sessionIndexCount)

	if userIndexCount == 0 || sessionIndexCount == 0 {
		result.Passed = false
		result.Message = "tenant_id索引缺失"
		result.Details = fmt.Sprintf("users表索引=%d, session表索引=%d", userIndexCount, sessionIndexCount)
		return result
	}

	result.Passed = true
	result.Message = "tenant_id索引存在"
	result.Details = fmt.Sprintf("users表索引=%d, session表索引=%d", userIndexCount, sessionIndexCount)
	return result
}

// validateOrphanData 验证孤儿数据
func (v *SessionUserValidator) validateOrphanData(ctx context.Context) *ValidationResult {
	result := &ValidationResult{
		Item:      "孤儿数据检查",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 检查 user_tenant 中是否有无效的 user_id
	var invalidCount int64
	v.db.Raw(`
		SELECT COUNT(*)
		FROM user_tenant ut
		LEFT JOIN user u ON ut.user_id = u.id
		WHERE u.id IS NULL AND ut.deleted_at IS NULL
	`).Scan(&invalidCount)

	if invalidCount > 0 {
		result.Passed = false
		result.Message = "存在无效的user_tenant关联"
		result.Details = fmt.Sprintf("无效关联数=%d", invalidCount)
		return result
	}

	result.Passed = true
	result.Message = "无孤儿数据"
	return result
}

// PrintReport 打印验证报告
func (v *SessionUserValidator) PrintReport(report *ValidationReport) {
	fmt.Println("\n========================================")
	fmt.Println("   Session/User 表迁移验证报告")
	fmt.Println("========================================")
	fmt.Printf("验证时间: %s\n", report.ValidatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("验证项目: %d\n", report.TotalItems)
	fmt.Printf("通过项目: %d\n", report.PassedItems)
	fmt.Printf("失败项目: %d\n", report.FailedItems)
	fmt.Printf("通过率: %.2f%%\n", report.PassRate)
	fmt.Printf("状态: %s\n", report.Status)
	fmt.Println("----------------------------------------")

	for _, result := range report.Results {
		status := "✅ PASS"
		if !result.Passed {
			status = "❌ FAIL"
		}
		fmt.Printf("%s [%s] %s\n", status, result.Item, result.Message)
		if result.Details != "" {
			fmt.Printf("       详情: %s\n", result.Details)
		}
	}

	fmt.Println("========================================\n")
}
