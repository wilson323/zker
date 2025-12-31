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

package botstore

// PublishBotToStoreRequest 发布Bot到商店请求
type PublishBotToStoreRequest struct {
	BotID       string   `json:"bot_id" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Category    string   `json:"category" binding:"required"`
	Tags        []string `json:"tags"`
	Price       float64  `json:"price" binding:"required,min=0"`
	Screenshots []string `json:"screenshots"`
	Version     string   `json:"version"`
}

// PublishBotToStoreResponse 发布Bot到商店响应
type PublishBotToStoreResponse struct {
	BaseResponse
	Data *BotStoreItemDTO `json:"data"`
}

// BotStoreItemDTO Bot商店项目DTO
type BotStoreItemDTO struct {
	ItemID        string    `json:"item_id"`
	BotID         string    `json:"bot_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	Tags          []string  `json:"tags"`
	Price         float64   `json:"price"`
	PublisherID   string    `json:"publisher_id"`
	Status        string    `json:"status"`
	ViewCount     int       `json:"view_count"`
	DownloadCount int       `json:"download_count"`
	Rating        float64   `json:"rating"`
	RatingCount   int       `json:"rating_count"`
	Screenshots   []string  `json:"screenshots"`
	Version       string    `json:"version"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

// UnpublishBotRequest 下架Bot请求
type UnpublishBotRequest struct {
	ItemID string `json:"item_id" binding:"required"`
}

// UnpublishBotResponse 下架Bot响应
type UnpublishBotResponse struct {
	BaseResponse
}

// UpdateBotStoreItemRequest 更新商店项目请求
type UpdateBotStoreItemRequest struct {
	ItemID      string   `json:"item_id" binding:"required"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Price       float64  `json:"price,min=0"`
	Screenshots []string `json:"screenshots"`
	Version     string   `json:"version"`
}

// UpdateBotStoreItemResponse 更新商店项目响应
type UpdateBotStoreItemResponse struct {
	BaseResponse
}

// ListBotStoreItemsRequest 列出商店项目请求
type ListBotStoreItemsRequest struct {
	Category string `query:"category"`
	SortBy   string `query:"sort_by" binding:"omitempty,oneof=popular latest rating"`
	Page     int    `query:"page" binding:"omitempty,min=1"`
	PageSize int    `query:"page_size" binding:"omitempty,min=1,max=100"`
}

// ListBotStoreItemsResponse 列出商店项目响应
type ListBotStoreItemsResponse struct {
	BaseResponse
	Data *BotStoreItemListData `json:"data"`
}

// BotStoreItemListData 商店项目列表数据
type BotStoreItemListData struct {
	Items      []*BotStoreItemDTO `json:"items"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

// SearchBotStoreItemsRequest 搜索商店项目请求
type SearchBotStoreItemsRequest struct {
	Query    string   `query:"query" binding:"required"`
	Category string   `query:"category"`
	Tags     []string `query:"tags"`
	PriceMin float64  `query:"price_min" binding:"omitempty,min=0"`
	PriceMax float64  `query:"price_max" binding:"omitempty,min=0"`
	Page     int      `query:"page" binding:"omitempty,min=1"`
	PageSize int      `query:"page_size" binding:"omitempty,min=1,max=100"`
}

// SearchBotStoreItemsResponse 搜索商店项目响应
type SearchBotStoreItemsResponse struct {
	BaseResponse
	Data *BotStoreItemListData `json:"data"`
}

// GetBotStoreItemRequest 获取商店项目请求
type GetBotStoreItemRequest struct {
	ItemID string `path:"item_id" binding:"required"`
}

// GetBotStoreItemResponse 获取商店项目响应
type GetBotStoreItemResponse struct {
	BaseResponse
	Data *BotStoreItemDTO `json:"data"`
}

// GetBotCategoriesResponse 获取分类列表响应
type GetBotCategoriesResponse struct {
	BaseResponse
	data []*BotStoreCategoryDTO `json:"data"`
}

// BotStoreCategoryDTO Bot商店分类DTO
type BotStoreCategoryDTO struct {
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	ParentID    string `json:"parent_id,omitempty"`
	SortOrder   int    `json:"sort_order"`
	BotCount    int    `json:"bot_count"`
}

// ReviewBotStoreItemRequest 审核Bot商店项目请求
type ReviewBotStoreItemRequest struct {
	Approved bool   `json:"approved" binding:"required"`
	Reason   string `json:"reason"`
}

// ReviewBotStoreItemResponse 审核Bot商店项目响应
type ReviewBotStoreItemResponse struct {
	BaseResponse
}

// GetPendingReviewsRequest 获取待审核列表请求
type GetPendingReviewsRequest struct {
	Page     int `query:"page" binding:"omitempty,min=1"`
	PageSize int `query:"page_size" binding:"omitempty,min=1,max=100"`
}

// GetPendingReviewsResponse 获取待审核列表响应
type GetPendingReviewsResponse struct {
	BaseResponse
	Data *BotStoreItemListData `json:"data"`
}

// CreateReviewRequest 创建评论请求
type CreateReviewRequest struct {
	ItemID  string `json:"item_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=1000"`
}

// CreateReviewResponse 创建评论响应
type CreateReviewResponse struct {
	BaseResponse
	Data *BotStoreReviewDTO `json:"data"`
}

// BotStoreReviewDTO Bot商店评论DTO
type BotStoreReviewDTO struct {
	ReviewID  string `json:"review_id"`
	ItemID    string `json:"item_id"`
	UserID    string `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UpdateReviewRequest 更新评论请求
type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"max=1000"`
}

// UpdateReviewResponse 更新评论响应
type UpdateReviewResponse struct {
	BaseResponse
}

// DeleteReviewResponse 删除评论响应
type DeleteReviewResponse struct {
	BaseResponse
}

// GetReviewsRequest 获取评论列表请求
type GetReviewsRequest struct {
	ItemID   string `path:"item_id" binding:"required"`
	Page     int    `query:"page" binding:"omitempty,min=1"`
	PageSize int    `query:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `query:"sort_by" binding:"omitempty,oneof=latest highest lowest"`
}

// GetReviewsResponse 获取评论列表响应
type GetReviewsResponse struct {
	BaseResponse
	Data *ReviewListData `json:"data"`
}

// ReviewListData 评论列表数据
type ReviewListData struct {
	Reviews       []*BotStoreReviewDTO `json:"reviews"`
	Total         int                  `json:"total"`
	Page          int                  `json:"page"`
	PageSize      int                  `json:"page_size"`
	TotalPages    int                  `json:"total_pages"`
	AverageRating float64              `json:"average_rating"`
	RatingCount   int                  `json:"rating_count"`
}

// GetReviewStatisticsResponse 获取评论统计响应
type GetReviewStatisticsResponse struct {
	BaseResponse
	Data *ReviewStatisticsDTO `json:"data"`
}

// ReviewStatisticsDTO 评论统计DTO
type ReviewStatisticsDTO struct {
	AverageRating float64 `json:"average_rating"`
	RatingCount   int     `json:"rating_count"`
	Rating1Count  int     `json:"rating_1_count"`
	Rating2Count  int     `json:"rating_2_count"`
	Rating3Count  int     `json:"rating_3_count"`
	Rating4Count  int     `json:"rating_4_count"`
	Rating5Count  int     `json:"rating_5_count"`
}

// BaseResponse 基础响应
type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
