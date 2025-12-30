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
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
)

// webhookRepository Webhook仓储实现
type webhookRepository struct {
	db *gorm.DB
}

// NewWebhookRepository 创建Webhook仓储实例
func NewWebhookRepository(db *gorm.DB) WebhookRepository {
	return &webhookRepository{db: db}
}

// Create 创建Webhook
func (r *webhookRepository) Create(ctx context.Context, webhook *entity.Webhook) error {
	webhook.CreatedAt = time.Now().UnixMilli()
	webhook.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(webhook).Error
}

// GetByID 根据ID获取Webhook
func (r *webhookRepository) GetByID(ctx context.Context, webhookID string) (*entity.Webhook, error) {
	var webhook entity.Webhook
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("webhook_id = ?", webhookID).
		Where("deleted_at IS NULL").
		First(&webhook).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// GetByProjectID 根据项目ID获取Webhook列表
func (r *webhookRepository) GetByProjectID(ctx context.Context, projectID string) ([]*entity.Webhook, error) {
	var webhooks []*entity.Webhook
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&webhooks).Error

	return webhooks, err
}

// GetByTenantID 根据租户ID获取Webhook列表
func (r *webhookRepository) GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Webhook, error) {
	var webhooks []*entity.Webhook
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Find(&webhooks).Error

	return webhooks, err
}

// Update 更新Webhook
func (r *webhookRepository) Update(ctx context.Context, webhook *entity.Webhook) error {
	webhook.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Webhook{}).
		Where("webhook_id = ?", webhook.WebhookID).
		Updates(webhook).Error
}

// Delete 软删除Webhook
func (r *webhookRepository) Delete(ctx context.Context, webhookID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Webhook{}).
		Where("webhook_id = ?", webhookID).
		Update("deleted_at", now).Error
}

// List 分页查询Webhook列表
func (r *webhookRepository) List(ctx context.Context, filter *WebhookFilter) ([]*entity.Webhook, int64, error) {
	var webhooks []*entity.Webhook
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Webhook{}).Where("deleted_at IS NULL")

	// 应用过滤器
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.ProjectID != "" {
		query = query.Where("project_id = ?", filter.ProjectID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := 0
	limit := filter.PageSize
	if limit <= 0 {
		limit = 20
	}

	err := query.
		Preload("Project").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&webhooks).Error

	return webhooks, total, err
}

// UpdateStatus 更新Webhook状态
func (r *webhookRepository) UpdateStatus(ctx context.Context, webhookID string, status entity.WebhookStatus) error {
	return r.db.WithContext(ctx).
		Model(&entity.Webhook{}).
		Where("webhook_id = ?", webhookID).
		Update("status", status).Error
}

// UpdateStats 更新统计信息
func (r *webhookRepository) UpdateStats(ctx context.Context, webhookID string, success bool) error {
	updates := map[string]interface{}{
		"updated_at": time.Now().UnixMilli(),
	}

	if success {
		updates["success_count"] = gorm.Expr("success_count + 1")
	} else {
		updates["failure_count"] = gorm.Expr("failure_count + 1")
	}

	return r.db.WithContext(ctx).
		Model(&entity.Webhook{}).
		Where("webhook_id = ?", webhookID).
		Updates(updates).Error
}

// UpdateLastTriggerAt 更新最后触发时间
func (r *webhookRepository) UpdateLastTriggerAt(ctx context.Context, webhookID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.Webhook{}).
		Where("webhook_id = ?", webhookID).
		Update("last_triggered_at", now).Error
}

// CountByProjectID 统计项目的Webhook数量
func (r *webhookRepository) CountByProjectID(ctx context.Context, projectID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Webhook{}).
		Where("project_id = ?", projectID).
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

// webhookLogRepository Webhook日志仓储实现
type webhookLogRepository struct {
	db *gorm.DB
}

// NewWebhookLogRepository 创建Webhook日志仓储实例
func NewWebhookLogRepository(db *gorm.DB) WebhookLogRepository {
	return &webhookLogRepository{db: db}
}

// Create 创建Webhook日志
func (r *webhookLogRepository) Create(ctx context.Context, log *entity.WebhookLog) error {
	log.CreatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(log).Error
}

// GetByWebhookID 根据Webhook ID获取日志列表
func (r *webhookLogRepository) GetByWebhookID(ctx context.Context, webhookID string, limit, offset int) ([]*entity.WebhookLog, error) {
	var logs []*entity.WebhookLog

	query := r.db.WithContext(ctx).
		Where("webhook_id = ?", webhookID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.
		Order("created_at DESC").
		Find(&logs).Error

	return logs, err
}

// DeleteOldLogs 删除旧日志（保留N天）
func (r *webhookLogRepository) DeleteOldLogs(ctx context.Context, days int) (int64, error) {
	cutoffTime := time.Now().AddDate(0, 0, -days).UnixMilli()

	result := r.db.WithContext(ctx).
		Where("created_at < ?", cutoffTime).
		Delete(&entity.WebhookLog{})

	return result.RowsAffected, result.Error
}

// GetStats 获取统计信息
func (r *webhookLogRepository) GetStats(ctx context.Context, webhookID string, startTime, endTime int64) (*WebhookStats, error) {
	var stats WebhookStats

	err := r.db.WithContext(ctx).
		Model(&entity.WebhookLog{}).
		Where("webhook_id = ?", webhookID).
		Where("created_at >= ?", startTime).
		Where("created_at <= ?", endTime).
		Select(
			"COUNT(*) as total_calls",
			"SUM(CASE WHEN success = true THEN 1 ELSE 0 END) as success_calls",
			"SUM(CASE WHEN success = false THEN 1 ELSE 0 END) as failure_calls",
			"AVG(duration_ms) as avg_duration",
		).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	// 计算成功率
	if stats.TotalCalls > 0 {
		stats.SuccessRate = float64(stats.SuccessCalls) / float64(stats.TotalCalls) * 100
	}

	return &stats, nil
}
