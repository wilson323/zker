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

package service

import (
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// CreateReviewRequest 创建评论请求
type CreateReviewRequest struct {
	ItemID  string  `json:"item_id"`
	Rating  int     `json:"rating"`  // 1-5
	Comment string  `json:"comment"`
}

// Validate 验证请求
func (r *CreateReviewRequest) Validate() error {
	if r.ItemID == "" {
		return errno.ErrBotStoreItemNotFound
	}
	if r.Rating < 1 || r.Rating > 5 {
		return errno.ErrInvalidRating
	}
	return nil
}

// UpdateReviewRequest 更新评论请求
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating"`  // 1-5
	Comment *string `json:"comment"`
}

// Validate 验证请求
func (r *UpdateReviewRequest) Validate() error {
	if r.Rating != nil && (*r.Rating < 1 || *r.Rating > 5) {
		return errno.ErrInvalidRating
	}
	return nil
}

// ListReviewsRequest 获取评论列表请求
type ListReviewsRequest struct {
	ItemID   string `json:"item_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"` // "latest", "oldest", "highest", "lowest"
}

// Validate 验证请求
func (r *ListReviewsRequest) Validate() error {
	if r.ItemID == "" {
		return errno.ErrBotStoreItemNotFound
	}
	if r.Page < 1 {
		r.Page = 1
	}
	if r.PageSize < 1 || r.PageSize > 100 {
		r.PageSize = 20
	}
	return nil
}

// ListReviewsResponse 获取评论列表响应
type ListReviewsResponse struct {
	Reviews       []*entity.BotStoreReview `json:"reviews"`
	Total         int64                    `json:"total"`
	Page          int                      `json:"page"`
	PageSize      int                      `json:"page_size"`
	TotalPages    int                      `json:"total_pages"`
	AverageRating float64                  `json:"average_rating"`
	RatingCount   int64                    `json:"rating_count"`
}

// botStoreReviewService Bot商店评论服务实现
type botStoreReviewService struct {
	storeRepo    repository.BotStoreRepository
	reviewRepo   repository.BotStoreReviewRepository
}

// NewBotStoreReviewService 创建Bot商店评论服务
func NewBotStoreReviewService(
	storeRepo repository.BotStoreRepository,
	reviewRepo repository.BotStoreReviewRepository,
) *botStoreReviewService {
	return &botStoreReviewService{
		storeRepo:  storeRepo,
		reviewRepo: reviewRepo,
	}
}

// CreateReview 创建评论
func (s *botStoreReviewService) CreateReview(ctx context.Context, req *CreateReviewRequest, userID, tenantID string) (*entity.BotStoreReview, error) {
	// 1. 验证请求
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. 检查商品是否存在
	item, err := s.storeRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return nil, errno.ErrBotStoreItemNotFound
	}

	// 3. 检查商品是否已发布
	if !item.IsPublished() {
		return nil, errno.ErrInvalidStatus
	}

	// 4. 检查用户是否已评论
	hasReviewed, err := s.reviewRepo.HasUserReviewed(ctx, userID, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user reviewed: %w", err)
	}
	if hasReviewed {
		return nil, errno.ErrAlreadyReviewed
	}

	// 5. 创建评论
	review := &entity.BotStoreReview{
		ReviewID: uuid.New().String(),
		ItemID:   req.ItemID,
		TenantID: tenantID,
		UserID:   userID,
		Rating:   req.Rating,
		Comment:  req.Comment,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// 6. 更新商品评分统计
	if err := s.reviewRepo.UpdateItemRating(ctx, req.ItemID); err != nil {
		// 记录错误但不影响评论创建
		// 可以考虑异步更新或重试机制
		return review, nil
	}

	return review, nil
}

// UpdateReview 更新评论
func (s *botStoreReviewService) UpdateReview(ctx context.Context, reviewID string, req *UpdateReviewRequest, userID string) error {
	// 1. 验证请求
	if err := req.Validate(); err != nil {
		return err
	}

	// 2. 获取评论
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return errno.ErrReviewNotFound
	}

	// 3. 检查权限
	if err := review.CanBeEdited(userID); err != nil {
		return err
	}

	// 4. 更新评论
	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Comment != nil {
		review.Comment = *req.Comment
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	// 5. 更新商品评分统计
	_ = s.reviewRepo.UpdateItemRating(ctx, review.ItemID)

	return nil
}

// DeleteReview 删除评论
func (s *botStoreReviewService) DeleteReview(ctx context.Context, reviewID string, userID string) error {
	// 1. 获取评论
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return errno.ErrReviewNotFound
	}

	// 2. 检查权限
	if err := review.CanBeEdited(userID); err != nil {
		return err
	}

	// 3. 删除评论
	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	// 4. 更新商品评分统计
	_ = s.reviewRepo.UpdateItemRating(ctx, review.ItemID)

	return nil
}

// ListReviews 获取评论列表
func (s *botStoreReviewService) ListReviews(ctx context.Context, req *ListReviewsRequest) (*ListReviewsResponse, error) {
	// 1. 验证请求
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. 检查商品是否存在
	_, err := s.storeRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return nil, errno.ErrBotStoreItemNotFound
	}

	// 3. 获取评论列表
	reviews, total, err := s.reviewRepo.ListByItemID(ctx, req.ItemID, req.Page, req.PageSize, req.SortBy)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}

	// 4. 获取评论统计
	stats, err := s.reviewRepo.GetStatistics(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get review statistics: %w", err)
	}

	// 5. 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	// 6. 返回响应
	return &ListReviewsResponse{
		Reviews:       reviews,
		Total:         int64(total),
		Page:          req.Page,
		PageSize:      req.PageSize,
		TotalPages:    totalPages,
		AverageRating: stats.AverageRating,
		RatingCount:   int64(stats.RatingCount),
	}, nil
}

// GetReview 获取单个评论
func (s *botStoreReviewService) GetReview(ctx context.Context, reviewID string) (*entity.BotStoreReview, error) {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, errno.ErrReviewNotFound
	}
	return review, nil
}

// GetReviewStatistics 获取评论统计
func (s *botStoreReviewService) GetReviewStatistics(ctx context.Context, itemID string) (*entity.ReviewStatistics, error) {
	// 1. 检查商品是否存在
	_, err := s.storeRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, errno.ErrBotStoreItemNotFound
	}

	// 2. 获取统计信息
	stats, err := s.reviewRepo.GetStatistics(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get review statistics: %w", err)
	}

	return stats, nil
}
