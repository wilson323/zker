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

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
)

// botStoreBrowser Bot商店浏览器实现
type botStoreBrowser struct {
	storeRepo    repository.BotStoreRepository
	categoryRepo repository.BotStoreCategoryRepository
}

// NewBotStoreBrowser 创建Bot商店浏览器
func NewBotStoreBrowser(
	storeRepo repository.BotStoreRepository,
	categoryRepo repository.BotStoreCategoryRepository,
) BotStoreBrowser {
	return &botStoreBrowser{
		storeRepo:    storeRepo,
		categoryRepo: categoryRepo,
	}
}

// ListBots 列出商店中的Bot（分页）
func (s *botStoreBrowser) ListBots(ctx context.Context, req *ListBotsRequest) (*ListBotsResponse, error) {
	// 1. 验证请求
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 2. 查询列表
	items, total, err := s.storeRepo.List(ctx, &repository.ListRequest{
		Category: req.Category,
		SortBy:   req.SortBy,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list bots: %w", err)
	}

	// 3. 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	// 4. 返回响应
	return &ListBotsResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// SearchBots 搜索Bot
func (s *botStoreBrowser) SearchBots(ctx context.Context, req *SearchBotsRequest) (*ListBotsResponse, error) {
	// 1. 验证请求
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	// 2. 执行搜索
	items, total, err := s.storeRepo.Search(ctx, &repository.SearchRequest{
		Query:    req.Query,
		Category: req.Category,
		Tags:     req.Tags,
		PriceMin: req.PriceMin,
		PriceMax: req.PriceMax,
		Status:   entity.BotStoreItemStatusPublished,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search bots: %w", err)
	}

	// 3. 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	// 4. 返回响应
	return &ListBotsResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetBotStoreItem 获取Bot详情
func (s *botStoreBrowser) GetBotStoreItem(ctx context.Context, itemID string) (*entity.BotStoreItem, error) {
	item, err := s.storeRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// 只返回已发布的项目
	if item.Status != entity.BotStoreItemStatusPublished {
		return nil, fmt.Errorf("bot store item not found")
	}

	return item, nil
}

// GetCategories 获取分类列表
func (s *botStoreBrowser) GetCategories(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	categories, err := s.categoryRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

// RecordView 记录浏览
func (s *botStoreBrowser) RecordView(ctx context.Context, itemID string) error {
	return s.storeRepo.IncrementViewCount(ctx, itemID)
}

// RecordDownload 记录下载
func (s *botStoreBrowser) RecordDownload(ctx context.Context, itemID string) error {
	return s.storeRepo.IncrementDownloadCount(ctx, itemID)
}
