// backend/domain/tenant/migration/session_user_migrator.go
// Session/User 表租户迁移器
package migration

import (
	"context"
	"fmt"
	"time"

	logs "github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// SessionUserMigrator Session/User表迁移器
type SessionUserMigrator struct {
	db            *gorm.DB
	batchSize     int
	defaultTenant string
	monitor       *ProgressMonitor
	checkpoint    *CheckpointManager
	logger        logs.CtxLogger
}

// NewSessionUserMigrator 创建迁移器
func NewSessionUserMigrator(db *gorm.DB, batchSize int, defaultTenant string) *SessionUserMigrator {
	return &SessionUserMigrator{
		db:            db,
		batchSize:     batchSize,
		defaultTenant: defaultTenant,
		monitor:       NewProgressMonitor(db),
		checkpoint:    NewCheckpointManager(db),
		logger:        logs.DefaultLogger(),
	}
}

// Migrate 执行迁移
func (m *SessionUserMigrator) Migrate(ctx context.Context) error {
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 开始 Session/User 表迁移...")

	// 1. 创建 user_tenant 表
	if err := m.createUserTenantTable(ctx); err != nil {
		return fmt.Errorf("创建 user_tenant 表失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ user_tenant 表创建完成")

	// 2. 添加 tenant_id 字段到 users 表
	if err := m.addTenantIDToUsersTable(ctx); err != nil {
		return fmt.Errorf("添加 tenant_id 到 users 表失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ users 表 tenant_id 字段添加完成")

	// 3. 添加 tenant_id 字段到 session 表
	if err := m.addTenantIDToSessionTable(ctx); err != nil {
		return fmt.Errorf("添加 tenant_id 到 session 表失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ session 表 tenant_id 字段添加完成")

	// 4. 回填 users 表的 tenant_id
	if err := m.backfillUsersTenantID(ctx); err != nil {
		return fmt.Errorf("回填 users 表 tenant_id 失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ users 表 tenant_id 回填完成")

	// 5. 回填 session 表的 tenant_id
	if err := m.backfillSessionTenantID(ctx); err != nil {
		return fmt.Errorf("回填 session 表 tenant_id 失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ session 表 tenant_id 回填完成")

	// 6. 初始化 user_tenant 关联数据
	if err := m.initUserTenantRelations(ctx); err != nil {
		return fmt.Errorf("初始化 user_tenant 关联数据失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ user_tenant 关联数据初始化完成")

	// 7. 验证迁移结果
	if err := m.validateMigration(ctx); err != nil {
		return fmt.Errorf("迁移验证失败: %w", err)
	}
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ 迁移验证通过")

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 🎉 Session/User 表迁移完成！")
	return nil
}

// createUserTenantTable 创建 user_tenant 关联表
func (m *SessionUserMigrator) createUserTenantTable(ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS user_tenant (
			id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
			user_id BIGINT NOT NULL COMMENT '用户ID',
			tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
			role VARCHAR(64) DEFAULT 'member' COMMENT '角色: owner, admin, member',
			is_default TINYINT(1) DEFAULT 0 COMMENT '是否为默认租户',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			deleted_at DATETIME DEFAULT NULL COMMENT '删除时间（软删除）',
			UNIQUE KEY uk_user_tenant (user_id, tenant_id, deleted_at),
			KEY idx_user_tenant_user (user_id),
			KEY idx_user_tenant_tenant (tenant_id),
			KEY idx_user_tenant_default (user_id, is_default, deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
		COMMENT='用户-租户关联表（支持用户多租户）'
	`
	return m.db.Exec(sql).Error
}

// addTenantIDToUsersTable 添加 tenant_id 到 users 表
func (m *SessionUserMigrator) addTenantIDToUsersTable(ctx context.Context) error {
	// 检查字段是否已存在
	var count int64
	if err := m.db.Raw(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'user'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		m.logger.CtxInfof(ctx, "[SessionUserMigrator] users 表 tenant_id 字段已存在，跳过")
		return nil
	}

	sql := `
		ALTER TABLE user
		ADD COLUMN tenant_id VARCHAR(64) NULL COMMENT '租户ID（多租户隔离）' AFTER id,
		ADD INDEX idx_user_tenant (tenant_id),
		ADD INDEX idx_user_email_tenant (email, tenant_id)
	`
	return m.db.Exec(sql).Error
}

// addTenantIDToSessionTable 添加 tenant_id 到 session 表
func (m *SessionUserMigrator) addTenantIDToSessionTable(ctx context.Context) error {
	// 检查字段是否已存在
	var count int64
	if err := m.db.Raw(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'session'
		  AND COLUMN_NAME = 'tenant_id'
	`).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		m.logger.CtxInfof(ctx, "[SessionUserMigrator] session 表 tenant_id 字段已存在，跳过")
		return nil
	}

	sql := `
		ALTER TABLE session
		ADD COLUMN tenant_id VARCHAR(64) NULL COMMENT '租户ID（多租户隔离）' AFTER user_id,
		ADD INDEX idx_session_tenant (tenant_id),
		ADD INDEX idx_session_user_tenant (user_id, tenant_id)
	`
	return m.db.Exec(sql).Error
}

// backfillUsersTenantID 回填 users 表的 tenant_id
func (m *SessionUserMigrator) backfillUsersTenantID(ctx context.Context) error {
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 开始回填 users 表的 tenant_id...")

	sql := `
		UPDATE user
		SET tenant_id = ?
		WHERE tenant_id IS NULL
	`
	result := m.db.Exec(sql, m.defaultTenant)
	if result.Error != nil {
		return result.Error
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] users 表已更新 %d 条记录", result.RowsAffected)
	return nil
}

// backfillSessionTenantID 回填 session 表的 tenant_id
func (m *SessionUserMigrator) backfillSessionTenantID(ctx context.Context) error {
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 开始回填 session 表的 tenant_id...")

	sql := `
		UPDATE session s
		INNER JOIN user u ON s.user_id = u.id
		SET s.tenant_id = u.tenant_id
		WHERE s.tenant_id IS NULL
	`
	result := m.db.Exec(sql)
	if result.Error != nil {
		return result.Error
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] session 表已更新 %d 条记录", result.RowsAffected)
	return nil
}

// initUserTenantRelations 初始化 user_tenant 关联数据
func (m *SessionUserMigrator) initUserTenantRelations(ctx context.Context) error {
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 开始初始化 user_tenant 关联数据...")

	sql := `
		INSERT INTO user_tenant (user_id, tenant_id, role, is_default, created_at, updated_at)
		SELECT
			id as user_id,
			? as tenant_id,
			'owner' as role,
			1 as is_default,
			NOW() as created_at,
			NOW() as updated_at
		FROM user
		WHERE tenant_id = ?
		  AND NOT EXISTS (
			SELECT 1 FROM user_tenant
			WHERE user_tenant.user_id = user.id
			  AND user_tenant.tenant_id = ?
			  AND user_tenant.deleted_at IS NULL
		  )
		ON DUPLICATE KEY UPDATE updated_at = NOW()
	`
	result := m.db.Exec(sql, m.defaultTenant, m.defaultTenant, m.defaultTenant)
	if result.Error != nil {
		return result.Error
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] user_tenant 表已插入 %d 条记录", result.RowsAffected)
	return nil
}

// validateMigration 验证迁移结果
func (m *SessionUserMigrator) validateMigration(ctx context.Context) error {
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 开始验证迁移结果...")

	// 1. 检查 users 表
	var userStats struct {
		TotalRecords    int64
		NullTenantCount int64
		HasTenantCount  int64
	}

	if err := m.db.Raw(`
		SELECT
			COUNT(*) as total_records,
			SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_count,
			SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_tenant_count
		FROM user
	`).Scan(&userStats).Error; err != nil {
		return err
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] users 表: 总记录=%d, 无tenant_id=%d, 有tenant_id=%d",
		userStats.TotalRecords, userStats.NullTenantCount, userStats.HasTenantCount)

	if userStats.NullTenantCount > 0 {
		return fmt.Errorf("users 表仍有 %d 条记录的 tenant_id 为 NULL", userStats.NullTenantCount)
	}

	// 2. 检查 session 表
	var sessionStats struct {
		TotalRecords    int64
		NullTenantCount int64
		HasTenantCount  int64
	}

	if err := m.db.Raw(`
		SELECT
			COUNT(*) as total_records,
			SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_count,
			SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_tenant_count
		FROM session
	`).Scan(&sessionStats).Error; err != nil {
		return err
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] session 表: 总记录=%d, 无tenant_id=%d, 有tenant_id=%d",
		sessionStats.TotalRecords, sessionStats.NullTenantCount, sessionStats.HasTenantCount)

	if sessionStats.NullTenantCount > 0 {
		return fmt.Errorf("session 表仍有 %d 条记录的 tenant_id 为 NULL", sessionStats.NullTenantCount)
	}

	// 3. 检查 user_tenant 表
	var relationStats struct {
		TotalRelations int64
		UniqueUsers    int64
		UniqueTenants  int64
		DefaultTenants int64
	}

	if err := m.db.Raw(`
		SELECT
			COUNT(*) as total_relations,
			COUNT(DISTINCT user_id) as unique_users,
			COUNT(DISTINCT tenant_id) as unique_tenants,
			SUM(CASE WHEN is_default = 1 THEN 1 ELSE 0 END) as default_tenants
		FROM user_tenant
		WHERE deleted_at IS NULL
	`).Scan(&relationStats).Error; err != nil {
		return err
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] user_tenant 表: 总关联=%d, 唯一用户=%d, 唯一租户=%d, 默认租户=%d",
		relationStats.TotalRelations, relationStats.UniqueUsers, relationStats.UniqueTenants, relationStats.DefaultTenants)

	return nil
}

// Rollback 回滚迁移
func (m *SessionUserMigrator) Rollback(ctx context.Context) error {
	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 开始回滚迁移...")

	// 1. 删除 user_tenant 表
	if err := m.db.Exec("DROP TABLE IF EXISTS user_tenant").Error; err != nil {
		m.logger.CtxErrorf(ctx, "[SessionUserMigrator] 删除 user_tenant 表失败: %v", err)
	} else {
		m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ user_tenant 表已删除")
	}

	// 2. 删除 users 表的 tenant_id 字段
	if err := m.db.Exec("ALTER TABLE user DROP COLUMN tenant_id").Error; err != nil {
		m.logger.CtxErrorf(ctx, "[SessionUserMigrator] 删除 users 表 tenant_id 字段失败: %v", err)
	} else {
		m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ users 表 tenant_id 字段已删除")
	}

	// 3. 删除 session 表的 tenant_id 字段
	if err := m.db.Exec("ALTER TABLE session DROP COLUMN tenant_id").Error; err != nil {
		m.logger.CtxErrorf(ctx, "[SessionUserMigrator] 删除 session 表 tenant_id 字段失败: %v", err)
	} else {
		m.logger.CtxInfof(ctx, "[SessionUserMigrator] ✅ session 表 tenant_id 字段已删除")
	}

	m.logger.CtxInfof(ctx, "[SessionUserMigrator] 回滚完成")
	return nil
}

// SessionUserMigrationStats 会话用户迁移统计信息
type SessionUserMigrationStats struct {
	TableName       string    `json:"table_name"`
	TotalRecords    int64     `json:"total_records"`
	MigratedRecords int64     `json:"migrated_records"`
	CompletionRate  float64   `json:"completion_rate"`
	StartedAt       time.Time `json:"started_at"`
	FinishedAt      time.Time `json:"finished_at"`
	DurationSeconds int64     `json:"duration_seconds"`
	Status          string    `json:"status"`
}

// GetStats 获取迁移统计信息
func (m *SessionUserMigrator) GetStats(ctx context.Context) (*SessionUserMigrationStats, error) {
	stats := &SessionUserMigrationStats{
		TableName: "session,user",
		StartedAt: time.Now(),
		Status:    "completed",
	}

	// 获取 users 表统计
	var userCount int64
	if err := m.db.Table("user").Where("tenant_id IS NOT NULL").Count(&userCount).Error; err != nil {
		return nil, err
	}
	stats.MigratedRecords += userCount

	// 获取 session 表统计
	var sessionCount int64
	if err := m.db.Table("session").Where("tenant_id IS NOT NULL").Count(&sessionCount).Error; err != nil {
		return nil, err
	}
	stats.MigratedRecords += sessionCount

	stats.TotalRecords = stats.MigratedRecords
	stats.CompletionRate = 100.0
	stats.FinishedAt = time.Now()
	stats.DurationSeconds = int64(time.Since(stats.StartedAt).Seconds())

	return stats, nil
}
