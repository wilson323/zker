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

// PublishBotRequest 发布Bot到商店请求
type PublishBotRequest struct {
	BotID        string   `json:"bot_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	Price        float64  `json:"price"`
	Screenshots  []string `json:"screenshots"`
	Version      string   `json:"version"`
	PublisherID  string   `json:"publisher_id"`
	TenantID     string   `json:"tenant_id"`
}

// UpdateBotStoreItemRequest 更新商店项目请求
type UpdateBotStoreItemRequest struct {
	ItemID       string   `json:"item_id"`
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Category     string   `json:"category,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Price        float64  `json:"price,omitempty"`
	Screenshots  []string `json:"screenshots,omitempty"`
	Version      string   `json:"version,omitempty"`
}

// BotStorePublisher Bot商店发布器接口
type BotStorePublisher interface {
	// PublishBot 发布Bot到商店
	PublishBot(ctx context.Context, req *PublishBotRequest) (*entity.BotStoreItem, error)

	// UnpublishBot 下架Bot
	UnpublishBot(ctx context.Context, itemID, userID string) error

	// UpdateBotStoreItem 更新商店中的Bot信息
	UpdateBotStoreItem(ctx context.Context, req *UpdateBotStoreItemRequest, userID string) error

	// GetPublishedBots 获取用户已发布的Bot列表
	GetPublishedBots(ctx context.Context, userID string, page, pageSize int) ([]*entity.BotStoreItem, int, error)
}
