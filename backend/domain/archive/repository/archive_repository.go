package repository

import (
	"context"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/archive/entity"
	"gorm.io/gorm"
)

// ArchiveRepository 归档仓储接口
type ArchiveRepository interface {
	// Create 创建归档记录
	Create(ctx context.Context, archive *entity.Archive) error

	// FindByID 根据ID查找归档
	FindByID(ctx context.Context, archiveID string) (*entity.Archive, error)

	// FindByTenant 查找租户的归档记录
	FindByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.Archive, error)

	// FindByTimeRange 根据时间范围查找归档
	FindByTimeRange(ctx context.Context, tenantID, tableName string, startDate, endDate int64) ([]*entity.Archive, error)

	// UpdateStatus 更新归档状态
	UpdateStatus(ctx context.Context, archiveID, status string) error

	// ListPending 列出待处理的归档
	ListPending(ctx context.Context, limit int) ([]*entity.Archive, error)

	// GetStats 获取归档统计信息
	GetStats(ctx context.Context, tenantID string) (*entity.ArchiveStats, error)

	// Delete 删除归档记录
	Delete(ctx context.Context, archiveID string) error
}

// ArchiveTaskRepository 归档任务仓储接口
type ArchiveTaskRepository interface {
	// Create 创建归档任务
	Create(ctx context.Context, task *entity.ArchiveTask) error

	// FindByID 根据ID查找任务
	FindByID(ctx context.Context, taskID string) (*entity.ArchiveTask, error)

	// FindByStatus 根据状态查找任务
	FindByStatus(ctx context.Context, status string, limit int) ([]*entity.ArchiveTask, error)

	// UpdateProgress 更新任务进度
	UpdateProgress(ctx context.Context, taskID string, processed int64) error

	// UpdateStatus 更新任务状态
	UpdateStatus(ctx context.Context, taskID, status string, errorMsg string) error

	// ListRecent 列出最近的任务
	ListRecent(ctx context.Context, limit int) ([]*entity.ArchiveTask, error)
}

// ArchiveConfigRepository 归档配置仓储接口
type ArchiveConfigRepository interface {
	// Create 创建归档配置
	Create(ctx context.Context, config *entity.ArchiveConfig) error

	// FindByID 根据ID查找配置
	FindByID(ctx context.Context, configID string) (*entity.ArchiveConfig, error)

	// FindByTableName 根据表名查找配置
	FindByTableName(ctx context.Context, tableName string) (*entity.ArchiveConfig, error)

	// ListEnabled 列出启用的配置
	ListEnabled(ctx context.Context) ([]*entity.ArchiveConfig, error)

	// Update 更新配置
	Update(ctx context.Context, config *entity.ArchiveConfig) error

	// Delete 删除配置
	Delete(ctx context.Context, configID string) error
}

// archiveRepository 归档仓储实现
type archiveRepository struct {
	db *gorm.DB
}

// NewArchiveRepository 创建归档仓储
func NewArchiveRepository(db *gorm.DB) ArchiveRepository {
	return &archiveRepository{db: db}
}

func (r *archiveRepository) Create(ctx context.Context, archive *entity.Archive) error {
	return r.db.WithContext(ctx).Create(archive).Error
}

func (r *archiveRepository) FindByID(ctx context.Context, archiveID string) (*entity.Archive, error) {
	var archive entity.Archive
	err := r.db.WithContext(ctx).Where("archive_id = ?", archiveID).First(&archive).Error
	if err != nil {
		return nil, err
	}
	return &archive, nil
}

func (r *archiveRepository) FindByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*entity.Archive, error) {
	var archives []*entity.Archive
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&archives).Error
	return archives, err
}

func (r *archiveRepository) FindByTimeRange(ctx context.Context, tenantID, tableName string, startDate, endDate int64) ([]*entity.Archive, error) {
	var archives []*entity.Archive
	query := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)

	if tableName != "" {
		query = query.Where("table_name = ?", tableName)
	}

	err := query.Where(
		"(date_start >= ? AND date_end <= ?) OR "+
			"(date_start <= ? AND date_end >= ?) OR "+
			"(date_start >= ? AND date_start <= ?)",
		startDate, endDate,
		startDate, startDate,
		startDate, endDate,
	).Order("date_start ASC").Find(&archives).Error

	return archives, err
}

func (r *archiveRepository) UpdateStatus(ctx context.Context, archiveID, status string) error {
	return r.db.WithContext(ctx).
		Model(&entity.Archive{}).
		Where("archive_id = ?", archiveID).
		Update("status", status).Error
}

func (r *archiveRepository) ListPending(ctx context.Context, limit int) ([]*entity.Archive, error) {
	var archives []*entity.Archive
	err := r.db.WithContext(ctx).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&archives).Error
	return archives, err
}

func (r *archiveRepository) GetStats(ctx context.Context, tenantID string) (*entity.ArchiveStats, error) {
	var stats entity.ArchiveStats

	// 总归档数
	r.db.WithContext(ctx).Model(&entity.Archive{}).
		Where("tenant_id = ? AND status = 'completed'", tenantID).
		Count(&stats.TotalArchives)

	// 总记录数
	r.db.WithContext(ctx).Model(&entity.Archive{}).
		Where("tenant_id = ? AND status = 'completed'", tenantID).
		Select("COALESCE(SUM(record_count), 0)").
		Scan(&stats.TotalRecords)

	// 总大小（GB）
	r.db.WithContext(ctx).Model(&entity.Archive{}).
		Where("tenant_id = ? AND status = 'completed'", tenantID).
		Select("COALESCE(SUM(file_size) / 1024.0 / 1024.0 / 1024.0, 0)").
		Scan(&stats.TotalSize)

	// 压缩比（假设平均为4:1）
	stats.CompressionRatio = 4.0

	// 存储成本（S3 Glacier: ¥0.03/GB/月）
	stats.StorageCost = stats.TotalSize * 0.03

	return &stats, nil
}

func (r *archiveRepository) Delete(ctx context.Context, archiveID string) error {
	return r.db.WithContext(ctx).
		Where("archive_id = ?", archiveID).
		Delete(&entity.Archive{}).Error
}

// archiveTaskRepository 归档任务仓储实现
type archiveTaskRepository struct {
	db *gorm.DB
}

// NewArchiveTaskRepository 创建归档任务仓储
func NewArchiveTaskRepository(db *gorm.DB) ArchiveTaskRepository {
	return &archiveTaskRepository{db: db}
}

func (r *archiveTaskRepository) Create(ctx context.Context, task *entity.ArchiveTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *archiveTaskRepository) FindByID(ctx context.Context, taskID string) (*entity.ArchiveTask, error) {
	var task entity.ArchiveTask
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *archiveTaskRepository) FindByStatus(ctx context.Context, status string, limit int) ([]*entity.ArchiveTask, error) {
	var tasks []*entity.ArchiveTask
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at ASC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (r *archiveTaskRepository) UpdateProgress(ctx context.Context, taskID string, processed int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.ArchiveTask{}).
		Where("task_id = ?", taskID).
		Updates(map[string]interface{}{
			"processed_records": processed,
			"updated_at":        now.Unix() * 1000,
		}).Error
}

func (r *archiveTaskRepository) UpdateStatus(ctx context.Context, taskID, status string, errorMsg string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": now.Unix() * 1000,
	}

	if status == "completed" {
		updates["completed_at"] = now
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	return r.db.WithContext(ctx).
		Model(&entity.ArchiveTask{}).
		Where("task_id = ?", taskID).
		Updates(updates).Error
}

func (r *archiveTaskRepository) ListRecent(ctx context.Context, limit int) ([]*entity.ArchiveTask, error) {
	var tasks []*entity.ArchiveTask
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

// archiveConfigRepository 归档配置仓储实现
type archiveConfigRepository struct {
	db *gorm.DB
}

// NewArchiveConfigRepository 创建归档配置仓储
func NewArchiveConfigRepository(db *gorm.DB) ArchiveConfigRepository {
	return &archiveConfigRepository{db: db}
}

func (r *archiveConfigRepository) Create(ctx context.Context, config *entity.ArchiveConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *archiveConfigRepository) FindByID(ctx context.Context, configID string) (*entity.ArchiveConfig, error) {
	var config entity.ArchiveConfig
	err := r.db.WithContext(ctx).Where("config_id = ?", configID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *archiveConfigRepository) FindByTableName(ctx context.Context, tableName string) (*entity.ArchiveConfig, error) {
	var config entity.ArchiveConfig
	err := r.db.WithContext(ctx).Where("table_name = ?", tableName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *archiveConfigRepository) ListEnabled(ctx context.Context) ([]*entity.ArchiveConfig, error) {
	var configs []*entity.ArchiveConfig
	err := r.db.WithContext(ctx).
		Where("enabled = ?", true).
		Order("priority DESC, created_at ASC").
		Find(&configs).Error
	return configs, err
}

func (r *archiveConfigRepository) Update(ctx context.Context, config *entity.ArchiveConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *archiveConfigRepository) Delete(ctx context.Context, configID string) error {
	return r.db.WithContext(ctx).
		Where("config_id = ?", configID).
		Delete(&entity.ArchiveConfig{}).Error
}
