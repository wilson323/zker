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

package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/migration/entity"
)

// MigrationAuditRepository 迁移审计仓储接口
type MigrationAuditRepository interface {
	// Create 创建审计记录
	Create(ctx context.Context, audit *entity.MigrationAudit) error

	// GetByID 根据ID获取审计记录
	GetByID(ctx context.Context, id int64) (*entity.MigrationAudit, error)

	// GetByMigrationID 根据迁移ID获取审计记录
	GetByMigrationID(ctx context.Context, migrationID string) (*entity.MigrationAudit, error)

	// Update 更新审计记录
	Update(ctx context.Context, audit *entity.MigrationAudit) error

	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, id int64, status entity.MigrationStatus) error

	// UpdateProgress 更新进度
	UpdateProgress(ctx context.Context, id int64, processedRecords, failedRecords int64) error

	// Complete 标记为完成
	Complete(ctx context.Context, id int64, processedRecords, failedRecords int64) error

	// Fail 标记为失败
	Fail(ctx context.Context, id int64, errorMessage string) error

	// List 查询审计记录列表
	List(ctx context.Context, limit, offset int) ([]*entity.MigrationAudit, int64, error)

	// GetLatestRunning 获取当前正在运行的迁移
	GetLatestRunning(ctx context.Context) (*entity.MigrationAudit, error)
}

// migrationAuditRepositoryImpl 迁移审计仓储实现
type migrationAuditRepositoryImpl struct {
	db *gorm.DB
}

// NewMigrationAuditRepository 创建迁移审计仓储实例
func NewMigrationAuditRepository(db *gorm.DB) MigrationAuditRepository {
	return &migrationAuditRepositoryImpl{db: db}
}

// Create 创建审计记录
func (r *migrationAuditRepositoryImpl) Create(ctx context.Context, audit *entity.MigrationAudit) error {
	if err := r.db.WithContext(ctx).Create(audit).Error; err != nil {
		return fmt.Errorf("failed to create migration audit: %w", err)
	}
	return nil
}

// GetByID 根据ID获取审计记录
func (r *migrationAuditRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.MigrationAudit, error) {
	var audit entity.MigrationAudit
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&audit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get migration audit by ID: %w", err)
	}
	return &audit, nil
}

// GetByMigrationID 根据迁移ID获取审计记录
func (r *migrationAuditRepositoryImpl) GetByMigrationID(ctx context.Context, migrationID string) (*entity.MigrationAudit, error) {
	var audit entity.MigrationAudit
	if err := r.db.WithContext(ctx).Where("migration_id = ?", migrationID).First(&audit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get migration audit by migration ID: %w", err)
	}
	return &audit, nil
}

// Update 更新审计记录
func (r *migrationAuditRepositoryImpl) Update(ctx context.Context, audit *entity.MigrationAudit) error {
	audit.UpdatedAt = time.Now().UnixMilli()
	if err := r.db.WithContext(ctx).Save(audit).Error; err != nil {
		return fmt.Errorf("failed to update migration audit: %w", err)
	}
	return nil
}

// UpdateStatus 更新状态
func (r *migrationAuditRepositoryImpl) UpdateStatus(ctx context.Context, id int64, status entity.MigrationStatus) error {
	now := time.Now().UnixMilli()
	updates := map[string]interface{}{
		"status":    status,
		"updated_at": now,
	}

	// 如果是开始执行，更新开始时间
	if status == entity.MigrationStatusRunning {
		updates["start_time"] = now
	}

	result := r.db.WithContext(ctx).
		Model(&entity.MigrationAudit{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update migration status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("migration audit not found: %d", id)
	}

	return nil
}

// UpdateProgress 更新进度
func (r *migrationAuditRepositoryImpl) UpdateProgress(ctx context.Context, id int64, processedRecords, failedRecords int64) error {
	now := time.Now().UnixMilli()
	result := r.db.WithContext(ctx).
		Model(&entity.MigrationAudit{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed_records": processedRecords,
			"failed_records":    failedRecords,
			"updated_at":        now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update migration progress: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("migration audit not found: %d", id)
	}

	return nil
}

// Complete 标记为完成
func (r *migrationAuditRepositoryImpl) Complete(ctx context.Context, id int64, processedRecords, failedRecords int64) error {
	now := time.Now().UnixMilli()
	var audit entity.MigrationAudit
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&audit).Error; err != nil {
		return fmt.Errorf("failed to get migration audit: %w", err)
	}

	endTime := now
	durationSeconds := (endTime - audit.StartTime) / 1000

	result := r.db.WithContext(ctx).
		Model(&entity.MigrationAudit{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":            entity.MigrationStatusCompleted,
			"end_time":          endTime,
			"duration_seconds":  durationSeconds,
			"processed_records": processedRecords,
			"failed_records":    failedRecords,
			"updated_at":        now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to complete migration: %w", result.Error)
	}

	return nil
}

// Fail 标记为失败
func (r *migrationAuditRepositoryImpl) Fail(ctx context.Context, id int64, errorMessage string) error {
	now := time.Now().UnixMilli()
	var audit entity.MigrationAudit
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&audit).Error; err != nil {
		return fmt.Errorf("failed to get migration audit: %w", err)
	}

	endTime := now
	durationSeconds := (endTime - audit.StartTime) / 1000

	result := r.db.WithContext(ctx).
		Model(&entity.MigrationAudit{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":           entity.MigrationStatusFailed,
			"end_time":         endTime,
			"duration_seconds": durationSeconds,
			"error_message":    errorMessage,
			"updated_at":       now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to mark migration as failed: %w", result.Error)
	}

	return nil
}

// List 查询审计记录列表
func (r *migrationAuditRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.MigrationAudit, int64, error) {
	var audits []*entity.MigrationAudit
	var total int64

	// 查询总数
	if err := r.db.WithContext(ctx).Model(&entity.MigrationAudit{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count migration audits: %w", err)
	}

	// 分页查询
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&audits).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list migration audits: %w", err)
	}

	return audits, total, nil
}

// GetLatestRunning 获取当前正在运行的迁移
func (r *migrationAuditRepositoryImpl) GetLatestRunning(ctx context.Context) (*entity.MigrationAudit, error) {
	var audit entity.MigrationAudit
	if err := r.db.WithContext(ctx).
		Where("status = ?", entity.MigrationStatusRunning).
		Order("start_time DESC").
		First(&audit).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get latest running migration: %w", err)
	}
	return &audit, nil
}
