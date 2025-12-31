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

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
)

// botStoreReviewRepository Bot商店评论仓储实现
type botStoreReviewRepository struct {
	db *gorm.DB
}

// NewBotStoreReviewRepository 创建Bot商店评论仓储
func NewBotStoreReviewRepository(db *gorm.DB) repository.BotStoreReviewRepository {
	return &botStoreReviewRepository{db: db}
}

// Create 创建评论
func (r *botStoreReviewRepository) Create(ctx context.Context, review *entity.BotStoreReview) error {
	if err := r.db.WithContext(ctx).Create(review).Error; err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

// Update 更新评论
func (r *botStoreReviewRepository) Update(ctx context.Context, review *entity.BotStoreReview) error {
	if err := r.db.WithContext(ctx).Save(review).Error; err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}
	return nil
}

// Delete 删除评论（软删除）
func (r *botStoreReviewRepository) Delete(ctx context.Context, reviewID string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&entity.BotStoreReview{}).
		Where("review_id = ?", reviewID).
		Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}

// GetByID 根据ID获取评论
func (r *botStoreReviewRepository) GetByID(ctx context.Context, reviewID string) (*entity.BotStoreReview, error) {
	var review entity.BotStoreReview
	if err := r.db.WithContext(ctx).
		Where("review_id = ? AND deleted_at IS NULL", reviewID).
		First(&review).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}
	return &review, nil
}

// GetByUserAndItem 根据用户ID和商品ID获取评论
func (r *botStoreReviewRepository) GetByUserAndItem(ctx context.Context, userID, itemID string) (*entity.BotStoreReview, error) {
	var review entity.BotStoreReview
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND item_id = ? AND deleted_at IS NULL", userID, itemID).
		First(&review).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 未找到评论，返回nil而不是错误
		}
		return nil, fmt.Errorf("failed to get review by user and item: %w", err)
	}
	return &review, nil
}

// ListByItemID 根据商品ID列出评论（分页）
func (r *botStoreReviewRepository) ListByItemID(ctx context.Context, itemID string, page, pageSize int, sortBy string) ([]*entity.BotStoreReview, int, error) {
	var reviews []*entity.BotStoreReview
	var total int64

	// 构建基础查询
	query := r.db.WithContext(ctx).
		Model(&entity.BotStoreReview{}).
		Where("item_id = ? AND deleted_at IS NULL", itemID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	// 应用排序
	orderBy := "created_at DESC"
	switch sortBy {
	case "highest":
		orderBy = "rating DESC, created_at DESC"
	case "lowest":
		orderBy = "rating ASC, created_at DESC"
	case "latest":
		orderBy = "created_at DESC"
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order(orderBy).Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list reviews: %w", err)
	}

	return reviews, int(total), nil
}

// ListByUserID 根据用户ID列出评论（分页）
func (r *botStoreReviewRepository) ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*entity.BotStoreReview, int, error) {
	var reviews []*entity.BotStoreReview
	var total int64

	query := r.db.WithContext(ctx).
		Model(&entity.BotStoreReview{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews by user: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list reviews by user: %w", err)
	}

	return reviews, int(total), nil
}

// GetStatistics 获取商品评论统计
func (r *botStoreReviewRepository) GetStatistics(ctx context.Context, itemID string) (*entity.ReviewStatistics, error) {
	type RatingCount struct {
		Rating int
		Count  int
	}

	var ratingCounts []*RatingCount
	if err := r.db.WithContext(ctx).
		Model(&entity.BotStoreReview{}).
		Select("rating, count(*) as count").
		Where("item_id = ? AND deleted_at IS NULL", itemID).
		Group("rating").
		Find(&ratingCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to get rating statistics: %w", err)
	}

	stats := &entity.ReviewStatistics{
		Rating1Count: 0,
		Rating2Count: 0,
		Rating3Count: 0,
		Rating4Count: 0,
		Rating5Count: 0,
	}

	totalRating := 0
	totalCount := 0

	for _, rc := range ratingCounts {
		totalCount += rc.Count
		totalRating += rc.Rating * rc.Count

		switch rc.Rating {
		case 1:
			stats.Rating1Count = rc.Count
		case 2:
			stats.Rating2Count = rc.Count
		case 3:
			stats.Rating3Count = rc.Count
		case 4:
			stats.Rating4Count = rc.Count
		case 5:
			stats.Rating5Count = rc.Count
		}
	}

	stats.RatingCount = totalCount
	if totalCount > 0 {
		stats.AverageRating = float64(totalRating) / float64(totalCount)
	}

	return stats, nil
}

// UpdateItemRating 更新商品评分统计
func (r *botStoreReviewRepository) UpdateItemRating(ctx context.Context, itemID string) error {
	stats, err := r.GetStatistics(ctx, itemID)
	if err != nil {
		return err
	}

	// 更新bot_store_items表中的评分统计
	if err := r.db.WithContext(ctx).
		Model(&entity.BotStoreItem{}).
		Where("item_id = ?", itemID).
		Updates(map[string]interface{}{
			"rating":        stats.AverageRating,
			"rating_count":  stats.RatingCount,
		}).Error; err != nil {
		return fmt.Errorf("failed to update item rating: %w", err)
	}

	return nil
}

// CountByItemID 统计商品评论数
func (r *botStoreReviewRepository) CountByItemID(ctx context.Context, itemID string) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.BotStoreReview{}).
		Where("item_id = ? AND deleted_at IS NULL", itemID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count reviews: %w", err)
	}
	return int(count), nil
}

// HasUserReviewed 检查用户是否已评论
func (r *botStoreReviewRepository) HasUserReviewed(ctx context.Context, userID, itemID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.BotStoreReview{}).
		Where("user_id = ? AND item_id = ? AND deleted_at IS NULL", userID, itemID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user reviewed: %w", err)
	}
	return count > 0, nil
}
