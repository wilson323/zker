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

// webhookDLQRepository Webhook死信队列仓储实现
type webhookDLQRepository struct {
	db *gorm.DB
}

// NewWebhookDLQRepository 创建Webhook死信队列仓储实例
func NewWebhookDLQRepository(db *gorm.DB) WebhookDLQRepository {
	return &webhookDLQRepository{db: db}
}

// Create 创建死信队列条目
func (r *webhookDLQRepository) Create(ctx context.Context, entry *entity.WebhookDLQEntry) error {
	entry.CreatedAt = time.Now().UnixMilli()

	// 设置过期时间（默认7天）
	if entry.ExpiresAt == nil {
		expiry := time.Now().Add(7 * 24 * time.Hour).UnixMilli()
		entry.ExpiresAt = &expiry
	}

	return r.db.WithContext(ctx).Create(entry).Error
}

// GetByID 根据ID获取死信队列条目
func (r *webhookDLQRepository) GetByID(ctx context.Context, dlqID string) (*entity.WebhookDLQEntry, error) {
	var entry entity.WebhookDLQEntry
	err := r.db.WithContext(ctx).
		Where("dlq_id = ?", dlqID).
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

// GetByWebhookID 根据Webhook ID获取死信队列条目列表
func (r *webhookDLQRepository) GetByWebhookID(ctx context.Context, webhookID string, limit, offset int) ([]*entity.WebhookDLQEntry, error) {
	var entries []*entity.WebhookDLQEntry

	query := r.db.WithContext(ctx).
		Where("webhook_id = ?", webhookID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&entries).Error
	return entries, err
}

// GetExpired 获取过期的死信队列条目（用于清理）
func (r *webhookDLQRepository) GetExpired(ctx context.Context) ([]*entity.WebhookDLQEntry, error) {
	var entries []*entity.WebhookDLQEntry

	// 查找7天前创建的条目
	expiry := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()

	err := r.db.WithContext(ctx).
		Where("created_at < ?", expiry).
		Find(&entries).Error

	return entries, err
}

// GetRetryable 获取可重试的条目
func (r *webhookDLQRepository) GetRetryable(ctx context.Context, maxRetries int) ([]*entity.WebhookDLQEntry, error) {
	var entries []*entity.WebhookDLQEntry

	// 查找重试次数未超过上限且未过期的条目
	expiry := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()

	err := r.db.WithContext(ctx).
		Where("retry_count < ?", maxRetries).
		Where("created_at > ?", expiry).
		Order("created_at ASC").
		Find(&entries).Error

	return entries, err
}

// UpdateRetryCount 更新重试次数
func (r *webhookDLQRepository) UpdateRetryCount(ctx context.Context, dlqID string, retryCount int) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.WebhookDLQEntry{}).
		Where("dlq_id = ?", dlqID).
		Updates(map[string]interface{}{
			"retry_count":   retryCount,
			"last_retry_at": now,
		}).Error
}

// Delete 根据ID删除死信队列条目
func (r *webhookDLQRepository) Delete(ctx context.Context, dlqID string) error {
	return r.db.WithContext(ctx).
		Where("dlq_id = ?", dlqID).
		Delete(&entity.WebhookDLQEntry{}).Error
}

// DeleteExpired 删除过期的死信队列条目
func (r *webhookDLQRepository) DeleteExpired(ctx context.Context) (int64, error) {
	var count int64

	// 删除7天前的条目
	expiry := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()

	result := r.db.WithContext(ctx).
		Where("created_at < ?", expiry).
		Delete(&entity.WebhookDLQEntry{})

	if err := result.Error; err != nil {
		return 0, err
	}

	count = result.RowsAffected
	return count, nil
}

// CountByWebhookID 统计Webhook的死信队列条目数量
func (r *webhookDLQRepository) CountByWebhookID(ctx context.Context, webhookID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.WebhookDLQEntry{}).
		Where("webhook_id = ?", webhookID).
		Count(&count).Error

	return count, err
}
