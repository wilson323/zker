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

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
)

// ListBotsRequest 列表请求
type ListBotsRequest struct {
	Category string `json:"category"`
	SortBy   string `json:"sort_by"` // popular, latest, rating
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

// SearchBotsRequest 搜索请求
type SearchBotsRequest struct {
	Query    string   `json:"query"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	PriceMin float64  `json:"price_min"`
	PriceMax float64  `json:"price_max"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

// ListBotsResponse 列表响应
type ListBotsResponse struct {
	Items       []*entity.BotStoreItem `json:"items"`
	Total       int                    `json:"total"`
	Page        int                    `json:"page"`
	PageSize    int                    `json:"page_size"`
	TotalPages  int                    `json:"total_pages"`
}

// BotStoreBrowser Bot商店浏览器接口
type BotStoreBrowser interface {
	// ListBots 列出商店中的Bot（分页）
	ListBots(ctx context.Context, req *ListBotsRequest) (*ListBotsResponse, error)

	// SearchBots 搜索Bot
	SearchBots(ctx context.Context, req *SearchBotsRequest) (*ListBotsResponse, error)

	// GetBotStoreItem 获取Bot详情
	GetBotStoreItem(ctx context.Context, itemID string) (*entity.BotStoreItem, error)

	// GetCategories 获取分类列表
	GetCategories(ctx context.Context) ([]*entity.BotStoreCategory, error)

	// RecordView 记录浏览
	RecordView(ctx context.Context, itemID string) error

	// RecordDownload 记录下载
	RecordDownload(ctx context.Context, itemID string) error
}
