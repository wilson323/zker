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

package impl

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/crossdomain/botstore"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	botstoreService "github.com/coze-dev/coze-studio/backend/domain/botstore/service"
)

// botStoreImpl Bot商店跨域服务实现
type botStoreImpl struct {
	publisher botstoreService.BotStorePublisher
	browser   botstoreService.BotStoreBrowser
	reviewer  botstoreService.BotStoreReviewer
}

// InitDomainService 初始化Bot商店跨域服务
func InitDomainService(
	publisher botstoreService.BotStorePublisher,
	browser botstoreService.BotStoreBrowser,
	reviewer botstoreService.BotStoreReviewer,
) botstore.BotStore {
	return &botStoreImpl{
		publisher: publisher,
		browser:   browser,
		reviewer:  reviewer,
	}
}

// PublishBot 发布Bot到商店
func (s *botStoreImpl) PublishBot(ctx context.Context, req *botstoreService.PublishBotRequest) (*entity.BotStoreItem, error) {
	return s.publisher.PublishBot(ctx, req)
}

// UnpublishBot 下架Bot
func (s *botStoreImpl) UnpublishBot(ctx context.Context, itemID, userID string) error {
	return s.publisher.UnpublishBot(ctx, itemID, userID)
}

// UpdateBotStoreItem 更新商店中的Bot信息
func (s *botStoreImpl) UpdateBotStoreItem(ctx context.Context, req *botstoreService.UpdateBotStoreItemRequest, userID string) error {
	return s.publisher.UpdateBotStoreItem(ctx, req, userID)
}

// ListBots 列出商店中的Bot（分页）
func (s *botStoreImpl) ListBots(ctx context.Context, req *botstoreService.ListBotsRequest) (*botstoreService.ListBotsResponse, error) {
	return s.browser.ListBots(ctx, req)
}

// SearchBots 搜索Bot
func (s *botStoreImpl) SearchBots(ctx context.Context, req *botstoreService.SearchBotsRequest) (*botstoreService.ListBotsResponse, error) {
	return s.browser.SearchBots(ctx, req)
}

// GetBotStoreItem 获取Bot详情
func (s *botStoreImpl) GetBotStoreItem(ctx context.Context, itemID string) (*entity.BotStoreItem, error) {
	return s.browser.GetBotStoreItem(ctx, itemID)
}

// GetCategories 获取分类列表
func (s *botStoreImpl) GetCategories(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	return s.browser.GetCategories(ctx)
}

// RecordView 记录浏览
func (s *botStoreImpl) RecordView(ctx context.Context, itemID string) error {
	return s.browser.RecordView(ctx, itemID)
}

// RecordDownload 记录下载
func (s *botStoreImpl) RecordDownload(ctx context.Context, itemID string) error {
	return s.browser.RecordDownload(ctx, itemID)
}

// ReviewBot 审核Bot
func (s *botStoreImpl) ReviewBot(ctx context.Context, req *botstoreService.ReviewBotRequest) error {
	return s.reviewer.ReviewBot(ctx, req)
}

// GetPendingReviews 获取待审核列表
func (s *botStoreImpl) GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error) {
	return s.reviewer.GetPendingReviews(ctx, page, pageSize)
}
