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

package entity

import (
	"errors"
	"time"
)

// BotStoreReview Bot商店评论实体
type BotStoreReview struct {
	ReviewID  string     `json:"review_id" gorm:"column:review_id;primaryKey"`
	ItemID    string     `json:"item_id" gorm:"column:item_id;not null"`
	TenantID  string     `json:"tenant_id" gorm:"column:tenant_id;not null"`
	UserID    string     `json:"user_id" gorm:"column:user_id;not null"`
	Rating    int        `json:"rating" gorm:"column:rating;not null"` // 1-5
	Comment   string     `json:"comment" gorm:"column:comment;type:text"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

// TableName 指定表名
func (BotStoreReview) TableName() string {
	return "bot_store_reviews"
}

// ValidateRating 验证评分范围
func (r *BotStoreReview) ValidateRating() error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	return nil
}

// CanBeEdited 检查是否可以编辑（只能编辑自己的评论）
func (r *BotStoreReview) CanBeEdited(userID string) error {
	if r.UserID != userID {
		return errors.New("no permission to edit this review")
	}
	return nil
}

// CreateReviewRequest 创建评论请求
type CreateReviewRequest struct {
	ItemID  string `json:"item_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=1000"`
}

// Validate 验证请求
func (r *CreateReviewRequest) Validate() error {
	if r.ItemID == "" {
		return errors.New("item ID is required")
	}
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	if len(r.Comment) > 1000 {
		return errors.New("comment too long, maximum 1000 characters")
	}
	return nil
}

// UpdateReviewRequest 更新评论请求
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=1000"`
}

// Validate 验证请求
func (r *UpdateReviewRequest) Validate() error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	if len(r.Comment) > 1000 {
		return errors.New("comment too long, maximum 1000 characters")
	}
	return nil
}

// ListReviewsRequest 列表请求
type ListReviewsRequest struct {
	ItemID   string `json:"item_id" binding:"required"`
	Page     int    `json:"page" binding:"min=1"`
	PageSize int    `json:"page_size" binding:"min=1,max=100"`
	SortBy   string `json:"sort_by"` // latest, highest, lowest
}

// Validate 验证请求
func (r *ListReviewsRequest) Validate() error {
	if r.ItemID == "" {
		return errors.New("item ID is required")
	}
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 20
	}
	if r.PageSize > 100 {
		r.PageSize = 100
	}
	// 默认排序
	if r.SortBy == "" {
		r.SortBy = "latest"
	}
	return nil
}

// ListReviewsResponse 列表响应
type ListReviewsResponse struct {
	Reviews       []*BotStoreReview `json:"reviews"`
	Total         int               `json:"total"`
	Page          int               `json:"page"`
	PageSize      int               `json:"page_size"`
	TotalPages    int               `json:"total_pages"`
	AverageRating float64           `json:"average_rating"`
	RatingCount   int               `json:"rating_count"`
}

// ReviewStatistics 评论统计
type ReviewStatistics struct {
	AverageRating float64 `json:"average_rating"`
	RatingCount   int     `json:"rating_count"`
	Rating1Count  int     `json:"rating_1_count"`
	Rating2Count  int     `json:"rating_2_count"`
	Rating3Count  int     `json:"rating_3_count"`
	Rating4Count  int     `json:"rating_4_count"`
	Rating5Count  int     `json:"rating_5_count"`
}
