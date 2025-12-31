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

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	botstoreService "github.com/coze-dev/coze-studio/backend/domain/botstore/service"
)

// BotStore 跨域Bot商店服务接口
type BotStore interface {
	// PublishBot 发布Bot到商店
	PublishBot(ctx context.Context, req *botstoreService.PublishBotRequest) (*entity.BotStoreItem, error)

	// UnpublishBot 下架Bot
	UnpublishBot(ctx context.Context, itemID, userID string) error

	// UpdateBotStoreItem 更新商店中的Bot信息
	UpdateBotStoreItem(ctx context.Context, req *botstoreService.UpdateBotStoreItemRequest, userID string) error

	// ListBots 列出商店中的Bot（分页）
	ListBots(ctx context.Context, req *botstoreService.ListBotsRequest) (*botstoreService.ListBotsResponse, error)

	// SearchBots 搜索Bot
	SearchBots(ctx context.Context, req *botstoreService.SearchBotsRequest) (*botstoreService.ListBotsResponse, error)

	// GetBotStoreItem 获取Bot详情
	GetBotStoreItem(ctx context.Context, itemID string) (*entity.BotStoreItem, error)

	// GetCategories 获取分类列表
	GetCategories(ctx context.Context) ([]*entity.BotStoreCategory, error)

	// RecordView 记录浏览
	RecordView(ctx context.Context, itemID string) error

	// RecordDownload 记录下载
	RecordDownload(ctx context.Context, itemID string) error

	// ReviewBot 审核Bot
	ReviewBot(ctx context.Context, req *botstoreService.ReviewBotRequest) error

	// GetPendingReviews 获取待审核列表
	GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error)
}

var defaultBotStore BotStore

// SetDefaultSVC 设置默认服务
func SetDefaultSVC(svc BotStore) {
	defaultBotStore = svc
}

// DefaultSVC 获取默认服务
func DefaultSVC() BotStore {
	return defaultBotStore
}
