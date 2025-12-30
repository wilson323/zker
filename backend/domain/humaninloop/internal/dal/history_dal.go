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

package dal

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/humaninloop/entity"
)

// HistoryDAL 历史记录数据访问层
type HistoryDAL struct {
	db *gorm.DB
}

// NewHistoryDAL 创建历史记录DAL
func NewHistoryDAL(db *gorm.DB) *HistoryDAL {
	return &HistoryDAL{db: db}
}

// Create 创建历史记录
func (d *HistoryDAL) Create(ctx context.Context, history *entity.CollaborationHistory) error {
	history.Timestamp = time.Now().Unix() * 1000

	if err := d.db.WithContext(ctx).Create(history).Error; err != nil {
		return fmt.Errorf("failed to create collaboration history: %w", err)
	}
	return nil
}

// GetByTaskID 获取任务的历史记录
func (d *HistoryDAL) GetByTaskID(ctx context.Context, taskID string) ([]*entity.CollaborationHistory, error) {
	var histories []*entity.CollaborationHistory
	err := d.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Order("timestamp ASC").
		Find(&histories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get task history: %w", err)
	}
	return histories, nil
}

// GetByTaskIDWithPagination 分页获取任务历史
func (d *HistoryDAL) GetByTaskIDWithPagination(
	ctx context.Context,
	taskID string,
	page, pageSize int,
) ([]*entity.CollaborationHistory, int64, error) {
	var histories []*entity.CollaborationHistory
	var total int64

	query := d.db.WithContext(ctx).
		Model(&entity.CollaborationHistory{}).
		Where("task_id = ?", taskID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count history: %w", err)
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	// 执行查询
	if err := query.Order("timestamp DESC").Find(&histories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get task history: %w", err)
	}

	return histories, total, nil
}

// GetByAction 获取指定操作类型的历史记录
func (d *HistoryDAL) GetByAction(
	ctx context.Context,
	taskID string,
	action entity.HistoryAction,
) ([]*entity.CollaborationHistory, error) {
	var histories []*entity.CollaborationHistory
	err := d.db.WithContext(ctx).
		Where("task_id = ? AND action = ?", taskID, action).
		Order("timestamp ASC").
		Find(&histories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get history by action: %w", err)
	}
	return histories, nil
}

// DeleteOldHistory 删除旧的历史记录
func (d *HistoryDAL) DeleteOldHistory(ctx context.Context, beforeDate time.Time) (int64, error) {
	result := d.db.WithContext(ctx).
		Where("timestamp < ?", beforeDate.Unix()*1000).
		Delete(&entity.CollaborationHistory{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old history: %w", result.Error)
	}

	return result.RowsAffected, nil
}
