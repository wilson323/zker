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

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
)

// BotStoreRepository Bot商店仓储接口
type BotStoreRepository interface {
	// Create 创建商店项目
	Create(ctx context.Context, item *entity.BotStoreItem) error

	// Update 更新商店项目
	Update(ctx context.Context, item *entity.BotStoreItem) error

	// GetByID 根据ID获取商店项目
	GetByID(ctx context.Context, itemID string) (*entity.BotStoreItem, error)

	// GetByBotID 根据BotID获取商店项目
	GetByBotID(ctx context.Context, botID string) (*entity.BotStoreItem, error)

	// Delete 删除商店项目（软删除）
	Delete(ctx context.Context, itemID string) error

	// List 列出商店项目（分页）
	List(ctx context.Context, req *ListRequest) ([]*entity.BotStoreItem, int, error)

	// Search 搜索商店项目
	Search(ctx context.Context, req *SearchRequest) ([]*entity.BotStoreItem, int, error)

	// GetPendingReviews 获取待审核列表
	GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error)

	// IncrementViewCount 增加浏览次数
	IncrementViewCount(ctx context.Context, itemID string) error

	// IncrementDownloadCount 增加下载次数
	IncrementDownloadCount(ctx context.Context, itemID string) error

	// UpdateRating 更新评分
	UpdateRating(ctx context.Context, itemID string, rating float64) error
}

// ListRequest 列表请求
type ListRequest struct {
	Category   string
	Status     entity.BotStoreItemStatus
	SortBy     string // popular, latest, rating
	Page       int
	PageSize   int
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query    string
	Category string
	Tags     []string
	PriceMin float64
	PriceMax float64
	Status   entity.BotStoreItemStatus
	Page     int
	PageSize int
}

// BotStoreCategoryRepository Bot商店分类仓储接口
type BotStoreCategoryRepository interface {
	// Create 创建分类
	Create(ctx context.Context, category *entity.BotStoreCategory) error

	// Update 更新分类
	Update(ctx context.Context, category *entity.BotStoreCategory) error

	// GetByID 根据ID获取分类
	GetByID(ctx context.Context, categoryID string) (*entity.BotStoreCategory, error)

	// List 列出所有分类
	List(ctx context.Context) ([]*entity.BotStoreCategory, error)

	// ListActive 列出活跃分类
	ListActive(ctx context.Context) ([]*entity.BotStoreCategory, error)

	// Delete 删除分类
	Delete(ctx context.Context, categoryID string) error

	// IncrementBotCount 增加Bot数量
	IncrementBotCount(ctx context.Context, categoryID string) error

	// DecrementBotCount 减少Bot数量
	DecrementBotCount(ctx context.Context, categoryID string) error
}
