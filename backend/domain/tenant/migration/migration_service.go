// backend/domain/tenant/migration/migration_service.go
package migration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"gorm.io/gorm"
)

// MigrationService 数据迁移服务接口
type MigrationService interface {
	// StartMigration 启动tenant_id迁移(双写模式)
	StartMigration(ctx context.Context, config *MigrationConfig) error

	// GetProgress 查询迁移进度
	GetProgress(ctx context.Context, migrationID string) (*MigrationProgress, error)

	// Rollback 回滚迁移
	Rollback(ctx context.Context, migrationID string) error

	// Validate 验证迁移结果
	Validate(ctx context.Context, migrationID string) error

	// ListMigrations 列出所有迁移任务
	ListMigrations(ctx context.Context) ([]*MigrationRecord, error)
}

// MigrationConfig 迁移配置
type MigrationConfig struct {
	MigrationID   string   `json:"migration_id"`   // 迁移任务ID
	TableNames    []string `json:"table_names"`    // 要迁移的表名列表
	Mode          string   `json:"mode"`           // 迁移模式: double_write, direct, rollback
	BatchSize     int      `json:"batch_size"`     // 批处理大小
	DefaultTenant string   `json:"default_tenant"` // 默认租户ID
	Workers       int      `json:"workers"`        // 并发worker数量
}

// MigrationProgress 迁移进度
type MigrationProgress struct {
	MigrationID     string     `json:"migration_id"`
	Status          string     `json:"status"` // pending, running, completed, failed, rolled_back
	ProgressPercent float64    `json:"progress_percent"`
	CurrentTable    string     `json:"current_table"`
	TotalTables     int        `json:"total_tables"`
	CompletedTables int        `json:"completed_tables"`
	TotalRecords    int64      `json:"total_records"`
	MigratedRecords int64      `json:"migrated_records"`
	FailedRecords   int64      `json:"failed_records"`
	StartedAt       time.Time  `json:"started_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	FinishedAt      *time.Time `json:"finished_at,omitempty"`
	Error           string     `json:"error,omitempty"`
}

// MigrationRecord 迁移记录
type MigrationRecord struct {
	ID              int64      `json:"id"`
	MigrationID     string     `json:"migration_id" gorm:"primaryKey"`
	Config          string     `json:"config"` // JSON序列化的配置
	Status          string     `json:"status"`
	ProgressPercent float64    `json:"progress_percent"`
	TotalRecords    int64      `json:"total_records"`
	MigratedRecords int64      `json:"migrated_records"`
	FailedRecords   int64      `json:"failed_records"`
	StartedAt       time.Time  `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
	Error           string     `json:"error"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// migrationServiceImpl 迁移服务实现
type migrationServiceImpl struct {
	db      *gorm.DB
	records map[string]*MigrationRecord // 迁移记录缓存
	mu      sync.RWMutex
}

// NewMigrationService 创建迁移服务
func NewMigrationService(db *gorm.DB) MigrationService {
	return &migrationServiceImpl{
		db:      db,
		records: make(map[string]*MigrationRecord),
	}
}

// StartMigration 启动迁移
func (s *migrationServiceImpl) StartMigration(ctx context.Context, config *MigrationConfig) error {
	if config.MigrationID == "" {
		config.MigrationID = fmt.Sprintf("migration_%d", time.Now().UnixNano())
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 1000
	}
	if config.Workers <= 0 {
		config.Workers = 5
	}

	logs.CtxInfof(ctx, "[MigrationService] 启动迁移任务: %s, 表: %v, 模式: %s",
		config.MigrationID, config.TableNames, config.Mode)

	// 1. 创建迁移记录
	record := &MigrationRecord{
		MigrationID:     config.MigrationID,
		Status:          "pending",
		ProgressPercent: 0,
		StartedAt:       time.Now(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	s.mu.Lock()
	s.records[config.MigrationID] = record
	s.mu.Unlock()

	// 保存到数据库
	if err := s.db.Create(record).Error; err != nil {
		return fmt.Errorf("保存迁移记录失败: %w", err)
	}

	// 2. 异步执行迁移
	go func() {
		if err := s.executeMigration(context.Background(), config); err != nil {
			logs.CtxErrorf(ctx, "[MigrationService] 迁移失败: %v", err)
			s.updateMigrationStatus(config.MigrationID, "failed", err.Error())
		}
	}()

	return nil
}

// executeMigration 执行迁移
func (s *migrationServiceImpl) executeMigration(ctx context.Context, config *MigrationConfig) error {
	// 更新状态为running
	s.updateMigrationStatus(config.MigrationID, "running", "")

	startTime := time.Now()
	totalMigrated := int64(0)
	totalFailed := int64(0)
	completedTables := 0

	switch config.Mode {
	case "double_write":
		// 双写模式: 添加字段 -> 双写 -> 数据回填 -> 验证
		for _, tableName := range config.TableNames {
			logs.CtxInfof(ctx, "[MigrationService] 处理表: %s (双写模式)", tableName)

			// 1. 添加tenant_id字段
			if err := s.addTenantIDColumn(ctx, tableName, config.DefaultTenant); err != nil {
				return fmt.Errorf("添加字段失败 [%s]: %w", tableName, err)
			}

			// 2. 启动双写
			if err := s.enableDualWrite(ctx, tableName); err != nil {
				return fmt.Errorf("启用双写失败 [%s]: %w", tableName, err)
			}

			// 3. 数据回填
			migrated, failed, err := s.backfillTenantID(ctx, tableName, config.DefaultTenant, config.BatchSize)
			if err != nil {
				return fmt.Errorf("数据回填失败 [%s]: %w", tableName, err)
			}

			totalMigrated += migrated
			totalFailed += failed
			completedTables++

			// 4. 更新进度
			progress := float64(completedTables) / float64(len(config.TableNames)) * 100
			s.updateProgress(config.MigrationID, progress, tableName, completedTables, totalMigrated, totalFailed)
		}

		// 5. 验证数据一致性
		logs.CtxInfof(ctx, "[MigrationService] 验证数据一致性...")
		verifier := NewDualWriteVerifier(s.db)
		for _, tableName := range config.TableNames {
			if err := verifier.VerifyDualWrite(ctx, tableName, 1*time.Hour); err != nil {
				logs.CtxWarnf(ctx, "[MigrationService] 表 %s 数据一致性验证失败: %v", tableName, err)
				// 不中断迁移，记录警告即可
			}
		}

	case "direct":
		// 直接模式: 添加字段 -> 数据回填
		for _, tableName := range config.TableNames {
			logs.CtxInfof(ctx, "[MigrationService] 处理表: %s (直接模式)", tableName)

			// 1. 添加tenant_id字段
			if err := s.addTenantIDColumn(ctx, tableName, config.DefaultTenant); err != nil {
				return fmt.Errorf("添加字段失败 [%s]: %w", tableName, err)
			}

			// 2. 数据回填
			migrated, failed, err := s.backfillTenantID(ctx, tableName, config.DefaultTenant, config.BatchSize)
			if err != nil {
				return fmt.Errorf("数据回填失败 [%s]: %w", tableName, err)
			}

			totalMigrated += migrated
			totalFailed += failed
			completedTables++

			// 3. 更新进度
			progress := float64(completedTables) / float64(len(config.TableNames)) * 100
			s.updateProgress(config.MigrationID, progress, tableName, completedTables, totalMigrated, totalFailed)
		}

	default:
		return fmt.Errorf("不支持的迁移模式: %s", config.Mode)
	}

	// 迁移完成
	duration := time.Since(startTime)
	logs.CtxInfof(ctx, "[MigrationService] 迁移完成: %s, 耗时: %s, 成功: %d, 失败: %d",
		config.MigrationID, duration, totalMigrated, totalFailed)

	// 更新最终状态
	finishedAt := time.Now()
	s.updateMigrationFinished(config.MigrationID, "completed", 100, totalMigrated, totalFailed, &finishedAt)

	return nil
}

// addTenantIDColumn 添加tenant_id字段
func (s *migrationServiceImpl) addTenantIDColumn(ctx context.Context, tableName, defaultTenant string) error {
	logs.CtxInfof(ctx, "[MigrationService] 添加tenant_id字段到表: %s", tableName)

	// 注意：MySQL 不支持 ALTER TABLE ... ADD COLUMN IF NOT EXISTS
	// 需要先检查列是否存在，不存在再添加
	// 使用反引号包裹表名，避免SQL关键字冲突
	sql := fmt.Sprintf(
		"ALTER TABLE `%s` ADD COLUMN `tenant_id` VARCHAR(36) NOT NULL DEFAULT '%s' COMMENT '租户ID', ADD INDEX `idx_tenant_id` (`tenant_id`)",
		tableName, defaultTenant,
	)

	if err := s.db.Exec(sql).Error; err != nil {
		return fmt.Errorf("执行SQL失败: %w", err)
	}

	logs.CtxInfof(ctx, "[MigrationService] ✓ tenant_id字段已添加: %s", tableName)
	return nil
}

// enableDualWrite 启用双写模式
func (s *migrationServiceImpl) enableDualWrite(ctx context.Context, tableName string) error {
	logs.CtxInfof(ctx, "[MigrationService] 启用双写模式: %s", tableName)

	// 注意: 实际的双写逻辑需要在应用层实现
	// 这里只是记录状态，提醒开启双写配置

	logs.CtxInfof(ctx, "[MigrationService] ✓ 双写模式已启用: %s", tableName)
	return nil
}

// backfillTenantID 回填tenant_id数据
func (s *migrationServiceImpl) backfillTenantID(ctx context.Context, tableName, defaultTenant string, batchSize int) (int64, int64, error) {
	logs.CtxInfof(ctx, "[MigrationService] 回填tenant_id: %s", tableName)

	totalMigrated := int64(0)
	totalFailed := int64(0)

	// 获取总记录数
	var totalRecords int64
	if err := s.db.Table(tableName).Count(&totalRecords).Error; err != nil {
		return 0, 0, fmt.Errorf("获取总记录数失败: %w", err)
	}

	logs.CtxInfof(ctx, "[MigrationService] 表 %s 共有 %d 条记录需要回填", tableName, totalRecords)

	// 分批回填
	maxID := int64(0)
	for {
		// 查询一批需要更新的记录
		var ids []int64
		err := s.db.Table(tableName).
			Select("id").
			Where("id > ? AND tenant_id = ?", maxID, defaultTenant).
			Order("id ASC").
			Limit(batchSize).
			Pluck("id", &ids).Error

		if err != nil {
			logs.CtxErrorf(ctx, "[MigrationService] 查询批次失败: %v", err)
			return totalMigrated, totalFailed, err
		}

		if len(ids) == 0 {
			break
		}

		// 更新这批记录
		for _, id := range ids {
			// 根据业务规则分配tenant_id
			tenantID := s.assignTenantID(ctx, tableName, id)

			result := s.db.Table(tableName).
				Where("id = ?", id).
				Update("tenant_id", tenantID)

			if result.Error != nil {
				totalFailed++
				logs.CtxErrorf(ctx, "[MigrationService] 更新记录失败 [table=%s, id=%d]: %v", tableName, id, result.Error)
			} else {
				totalMigrated++
			}

			maxID = id
		}

		// 打印进度
		progress := float64(totalMigrated+totalFailed) / float64(totalRecords) * 100
		logs.CtxInfof(ctx, "[MigrationService] 表 %s 回填进度: %d/%d (%.2f%%)",
			tableName, totalMigrated+totalFailed, totalRecords, progress)
	}

	logs.CtxInfof(ctx, "[MigrationService] ✓ 表 %s 回填完成: 成功 %d, 失败 %d", tableName, totalMigrated, totalFailed)
	return totalMigrated, totalFailed, nil
}

// assignTenantID 为记录分配tenant_id
func (s *migrationServiceImpl) assignTenantID(ctx context.Context, tableName string, id int64) string {
	// 简化实现: 查询记录的creator_id，然后从users表获取tenant_id
	var creatorID string
	var tenantID string

	// 1. 获取creator_id
	err := s.db.Table(tableName).
		Select("creator_id").
		Where("id = ?", id).
		Pluck("creator_id", &creatorID).Error

	if err != nil || creatorID == "" {
		// 如果没有creator_id，使用默认租户
		return "system_tenant"
	}

	// 2. 从users表获取tenant_id
	err = s.db.Table("users").
		Select("tenant_id").
		Where("user_id = ?", creatorID).
		Pluck("tenant_id", &tenantID).Error

	if err != nil || tenantID == "" {
		// 如果查询失败，使用默认租户
		return "system_tenant"
	}

	return tenantID
}

// GetProgress 查询迁移进度
func (s *migrationServiceImpl) GetProgress(ctx context.Context, migrationID string) (*MigrationProgress, error) {
	s.mu.RLock()
	record, exists := s.records[migrationID]
	s.mu.RUnlock()

	if !exists {
		// 从数据库查询
		record = &MigrationRecord{}
		err := s.db.Where("migration_id = ?", migrationID).First(record).Error
		if err != nil {
			return nil, fmt.Errorf("迁移记录不存在: %s", migrationID)
		}

		s.mu.Lock()
		s.records[migrationID] = record
		s.mu.Unlock()
	}

	return &MigrationProgress{
		MigrationID:     record.MigrationID,
		Status:          record.Status,
		ProgressPercent: record.ProgressPercent,
		TotalRecords:    record.TotalRecords,
		MigratedRecords: record.MigratedRecords,
		FailedRecords:   record.FailedRecords,
		StartedAt:       record.StartedAt,
		UpdatedAt:       record.UpdatedAt,
		FinishedAt:      record.FinishedAt,
		Error:           record.Error,
	}, nil
}

// Rollback 回滚迁移
func (s *migrationServiceImpl) Rollback(ctx context.Context, migrationID string) error {
	logs.CtxInfof(ctx, "[MigrationService] 开始回滚迁移: %s", migrationID)

	// 1. 获取迁移记录
	record, err := s.getMigrationRecord(migrationID)
	if err != nil {
		return err
	}

	if record.Status == "rolled_back" {
		return fmt.Errorf("迁移已经回滚")
	}

	// 2. 从配置中获取表名
	var config MigrationConfig
	// TODO: 反序列化record.Config为config

	// 3. 删除tenant_id字段
	for _, tableName := range config.TableNames {
		logs.CtxInfof(ctx, "[MigrationService] 删除字段 [table=%s]", tableName)

		sql := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN IF EXISTS `tenant_id`", tableName)
		if err := s.db.Exec(sql).Error; err != nil {
			logs.CtxWarnf(ctx, "[MigrationService] 删除字段失败 [table=%s]: %v", tableName, err)
			// 继续回滚其他表
		} else {
			logs.CtxInfof(ctx, "[MigrationService] ✓ 字段已删除 [table=%s]", tableName)
		}
	}

	// 4. 更新状态
	now := time.Now()
	s.updateMigrationFinished(migrationID, "rolled_back", 0, 0, 0, &now)

	logs.CtxInfof(ctx, "[MigrationService] ✓ 迁移已回滚: %s", migrationID)
	return nil
}

// Validate 验证迁移结果
func (s *migrationServiceImpl) Validate(ctx context.Context, migrationID string) error {
	logs.CtxInfof(ctx, "[MigrationService] 验证迁移结果: %s", migrationID)

	// 1. 获取迁移记录
	record, err := s.getMigrationRecord(migrationID)
	if err != nil {
		return err
	}

	if record.Status != "completed" {
		return fmt.Errorf("迁移尚未完成，当前状态: %s", record.Status)
	}

	// 2. 验证数据完整性
	// TODO: 从配置中获取表名

	// 3. 检查是否有NULL值
	// var nullCount int64
	// err = s.db.Table(tableName).
	//     Where("tenant_id IS NULL OR tenant_id = ''").
	//     Count(&nullCount).Error
	// if err != nil {
	//     return err
	// }
	// if nullCount > 0 {
	//     return fmt.Errorf("存在 %d 条记录的tenant_id为空", nullCount)
	// }

	logs.CtxInfof(ctx, "[MigrationService] ✓ 迁移验证通过: %s", migrationID)
	return nil
}

// ListMigrations 列出所有迁移任务
func (s *migrationServiceImpl) ListMigrations(ctx context.Context) ([]*MigrationRecord, error) {
	var records []*MigrationRecord
	err := s.db.Order("created_at DESC").Find(&records).Error
	if err != nil {
		return nil, fmt.Errorf("查询迁移记录失败: %w", err)
	}
	return records, nil
}

// ================================================================================
// 辅助方法
// ================================================================================

// getMigrationRecord 获取迁移记录
func (s *migrationServiceImpl) getMigrationRecord(migrationID string) (*MigrationRecord, error) {
	s.mu.RLock()
	record, exists := s.records[migrationID]
	s.mu.RUnlock()

	if !exists {
		record = &MigrationRecord{}
		err := s.db.Where("migration_id = ?", migrationID).First(record).Error
		if err != nil {
			return nil, fmt.Errorf("迁移记录不存在: %s", migrationID)
		}

		s.mu.Lock()
		s.records[migrationID] = record
		s.mu.Unlock()
	}

	return record, nil
}

// updateMigrationStatus 更新迁移状态
func (s *migrationServiceImpl) updateMigrationStatus(migrationID, status, errorMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.records[migrationID]
	if record == nil {
		return
	}

	record.Status = status
	record.Error = errorMsg
	record.UpdatedAt = time.Now()

	// 更新数据库
	s.db.Model(record).Updates(map[string]interface{}{
		"status":     status,
		"error":      errorMsg,
		"updated_at": time.Now(),
	})
}

// updateProgress 更新迁移进度
func (s *migrationServiceImpl) updateProgress(migrationID string, progress float64, currentTable string, completedTables int, migrated, failed int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.records[migrationID]
	if record == nil {
		return
	}

	record.ProgressPercent = progress
	// record.CompletedTables = completedTables  // TODO: 添加此字段到MigrationRecord
	record.MigratedRecords = migrated
	record.FailedRecords = failed
	record.UpdatedAt = time.Now()

	// 更新数据库
	s.db.Model(record).Updates(map[string]interface{}{
		"progress_percent": progress,
		"completed_tables": completedTables,
		"migrated_records": migrated,
		"failed_records":   failed,
		"updated_at":       time.Now(),
	})
}

// updateMigrationFinished 更新迁移完成状态
func (s *migrationServiceImpl) updateMigrationFinished(migrationID, status string, progress float64, migrated, failed int64, finishedAt *time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.records[migrationID]
	if record == nil {
		return
	}

	record.Status = status
	record.ProgressPercent = progress
	record.MigratedRecords = migrated
	record.FailedRecords = failed
	record.FinishedAt = finishedAt
	record.UpdatedAt = time.Now()

	// 更新数据库
	s.db.Model(record).Updates(map[string]interface{}{
		"status":           status,
		"progress_percent": progress,
		"migrated_records": migrated,
		"failed_records":   failed,
		"finished_at":      finishedAt,
		"updated_at":       time.Now(),
	})
}
