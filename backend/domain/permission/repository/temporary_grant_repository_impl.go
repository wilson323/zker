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

	"github.com/coze-dev/coze-studio/backend/domain/permission/entity"
)

// temporaryGrantRepository 临时授权仓储实现
type temporaryGrantRepository struct {
	db *gorm.DB
}

// NewTemporaryGrantRepository 创建临时授权仓储实例
func NewTemporaryGrantRepository(db *gorm.DB) TemporaryGrantRepository {
	return &temporaryGrantRepository{db: db}
}

// Create 创建临时授权
func (r *temporaryGrantRepository) Create(ctx context.Context, grant *entity.TemporaryGrant) error {
	grant.CreatedAt = time.Now().UnixMilli()
	grant.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(grant).Error
}

// GetByID 根据ID获取临时授权
func (r *temporaryGrantRepository) GetByID(ctx context.Context, grantID int64) (*entity.TemporaryGrant, error) {
	var grant entity.TemporaryGrant
	err := r.db.WithContext(ctx).
		Where("grant_id = ?", grantID).
		First(&grant).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &grant, nil
}

// GetByGrantCode 根据授权码获取临时授权
func (r *temporaryGrantRepository) GetByGrantCode(ctx context.Context, grantCode string) (*entity.TemporaryGrant, error) {
	var grant entity.TemporaryGrant
	err := r.db.WithContext(ctx).
		Where("grant_code = ?", grantCode).
		First(&grant).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &grant, nil
}

// Update 更新临时授权
func (r *temporaryGrantRepository) Update(ctx context.Context, grant *entity.TemporaryGrant) error {
	grant.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.TemporaryGrant{}).
		Where("grant_id = ?", grant.GrantID).
		Updates(grant).Error
}

// Delete 删除临时授权
func (r *temporaryGrantRepository) Delete(ctx context.Context, grantID int64) error {
	return r.db.WithContext(ctx).
		Delete(&entity.TemporaryGrant{}, grantID).Error
}

// List 获取租户的所有临时授权
func (r *temporaryGrantRepository) List(ctx context.Context, filter *TemporaryGrantFilter) ([]*entity.TemporaryGrant, int64, error) {
	var grants []*entity.TemporaryGrant
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TemporaryGrant{})

	// 应用过滤器
	if filter.TenantID != "" {
		query = query.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.GranteeID != "" {
		query = query.Where("grantee_id = ?", filter.GranteeID)
	}
	if filter.GrantorID != "" {
		query = query.Where("grantor_id = ?", filter.GrantorID)
	}
	if filter.PermissionType != nil {
		query = query.Where("permission_type = ?", *filter.PermissionType)
	}
	if filter.IsUsed != nil {
		query = query.Where("is_used = ?", *filter.IsUsed)
	}
	if filter.IsRevoked != nil {
		query = query.Where("is_revoked = ?", *filter.IsRevoked)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	query = query.Order("created_at DESC")
	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize)
	}
	if filter.PageToken != "" {
		query = query.Where("grant_id > ?", filter.PageToken)
	}

	if err := query.Find(&grants).Error; err != nil {
		return nil, 0, err
	}

	return grants, total, nil
}

// ListExpired 获取所有过期的临时授权
func (r *temporaryGrantRepository) ListExpired(ctx context.Context, expiresBefore int64) ([]*entity.TemporaryGrant, error) {
	var grants []*entity.TemporaryGrant
	err := r.db.WithContext(ctx).
		Where("expires_at < ?", expiresBefore).
		Where("is_used = ?", false).
		Where("is_revoked = ?", false).
		Find(&grants).Error

	if err != nil {
		return nil, err
	}
	return grants, nil
}

// DeleteExpired 删除过期的临时授权
func (r *temporaryGrantRepository) DeleteExpired(ctx context.Context, expiresBefore int64) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", expiresBefore).
		Delete(&entity.TemporaryGrant{})

	return result.RowsAffected, result.Error
}

// GetByGrantee 获取被授权人的所有临时授权
func (r *temporaryGrantRepository) GetByGrantee(ctx context.Context, granteeID string) ([]*entity.TemporaryGrant, error) {
	var grants []*entity.TemporaryGrant
	err := r.db.WithContext(ctx).
		Where("grantee_id = ?", granteeID).
		Order("created_at DESC").
		Find(&grants).Error

	if err != nil {
		return nil, err
	}
	return grants, nil
}

// GetByGrantor 获取授权人的所有临时授权
func (r *temporaryGrantRepository) GetByGrantor(ctx context.Context, grantorID string) ([]*entity.TemporaryGrant, error) {
	var grants []*entity.TemporaryGrant
	err := r.db.WithContext(ctx).
		Where("grantor_id = ?", grantorID).
		Order("created_at DESC").
		Find(&grants).Error

	if err != nil {
		return nil, err
	}
	return grants, nil
}

// temporaryGrantHistoryRepository 临时授权历史仓储实现
type temporaryGrantHistoryRepository struct {
	db *gorm.DB
}

// NewTemporaryGrantHistoryRepository 创建临时授权历史仓储实例
func NewTemporaryGrantHistoryRepository(db *gorm.DB) TemporaryGrantHistoryRepository {
	return &temporaryGrantHistoryRepository{db: db}
}

// Create 创建历史记录
func (r *temporaryGrantHistoryRepository) Create(ctx context.Context, history *entity.TemporaryGrantHistory) error {
	history.CreatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(history).Error
}

// GetByGrantID 获取临时授权的所有历史记录
func (r *temporaryGrantHistoryRepository) GetByGrantID(ctx context.Context, grantID int64) ([]*entity.TemporaryGrantHistory, error) {
	var histories []*entity.TemporaryGrantHistory
	err := r.db.WithContext(ctx).
		Where("grant_id = ?", grantID).
		Order("created_at DESC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}
	return histories, nil
}

// GetByGrantCode 根据授权码获取历史记录
func (r *temporaryGrantHistoryRepository) GetByGrantCode(ctx context.Context, grantCode string) ([]*entity.TemporaryGrantHistory, error) {
	var histories []*entity.TemporaryGrantHistory
	err := r.db.WithContext(ctx).
		Where("grant_code = ?", grantCode).
		Order("created_at DESC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}
	return histories, nil
}

// List 获取租户的所有历史记录
func (r *temporaryGrantHistoryRepository) List(ctx context.Context, filter *HistoryFilter) ([]*entity.TemporaryGrantHistory, int64, error) {
	var histories []*entity.TemporaryGrantHistory
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TemporaryGrantHistory{})

	// 应用过滤器
	if filter.GrantID > 0 {
		query = query.Where("grant_id = ?", filter.GrantID)
	}
	if filter.GrantCode != "" {
		query = query.Where("grant_code = ?", filter.GrantCode)
	}
	if filter.ActionType != nil {
		query = query.Where("action_type = ?", *filter.ActionType)
	}
	if filter.OperatorID != "" {
		query = query.Where("operator_id = ?", filter.OperatorID)
	}
	if filter.Before > 0 {
		query = query.Where("created_at < ?", filter.Before)
	}
	if filter.After > 0 {
		query = query.Where("created_at > ?", filter.After)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	query = query.Order("created_at DESC")
	if filter.PageSize > 0 {
		query = query.Limit(filter.PageSize)
	}
	if filter.PageToken != "" {
		query = query.Where("history_id > ?", filter.PageToken)
	}

	if err := query.Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// DeleteByGrantID 删除临时授权的所有历史记录
func (r *temporaryGrantHistoryRepository) DeleteByGrantID(ctx context.Context, grantID int64) error {
	return r.db.WithContext(ctx).
		Where("grant_id = ?", grantID).
		Delete(&entity.TemporaryGrantHistory{}).Error
}

// DeleteExpired 删除过期的历史记录（例如90天前）
func (r *temporaryGrantHistoryRepository) DeleteExpired(ctx context.Context, before int64) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("created_at < ?", before).
		Delete(&entity.TemporaryGrantHistory{})

	return result.RowsAffected, result.Error
}
