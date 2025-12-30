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
	"time"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// botStorePublisher Bot商店发布器实现
type botStorePublisher struct {
	storeRepo     repository.BotStoreRepository
	categoryRepo  repository.BotStoreCategoryRepository
}

// NewBotStorePublisher 创建Bot商店发布器
func NewBotStorePublisher(
	storeRepo repository.BotStoreRepository,
	categoryRepo repository.BotStoreCategoryRepository,
) BotStorePublisher {
	return &botStorePublisher{
		storeRepo:    storeRepo,
		categoryRepo: categoryRepo,
	}
}

// PublishBot 发布Bot到商店
func (s *botStorePublisher) PublishBot(ctx context.Context, req *PublishBotRequest) (*entity.BotStoreItem, error) {
	// 1. 验证请求
	if err := s.validatePublishRequest(req); err != nil {
		return nil, err
	}

	// 2. 检查Bot是否已发布
	existingItem, err := s.storeRepo.GetByBotID(ctx, req.BotID)
	if err == nil && existingItem != nil {
		// Bot已存在，检查状态
		if existingItem.Status == entity.BotStoreItemStatusPublished {
			return nil, errno.ErrAlreadyPublished
		}
		if existingItem.Status == entity.BotStoreItemStatusPending {
			return nil, errno.ErrPendingReview
		}
	}

	// 3. 验证分类是否存在
	if req.Category != "" {
		category, err := s.categoryRepo.GetByID(ctx, req.Category)
		if err != nil || category == nil {
			return nil, errno.ErrInvalidCategory
		}
	}

	// 4. 创建商店项目
	now := time.Now()
	item := &entity.BotStoreItem{
		ItemID:       uuid.New().String(),
		BotID:        req.BotID,
		TenantID:     req.TenantID,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Tags:         req.Tags,
		Price:        req.Price,
		PublisherID:  req.PublisherID,
		Status:       entity.BotStoreItemStatusPending, // 默认待审核
		Screenshots:  req.Screenshots,
		Version:      req.Version,
		ViewCount:    0,
		DownloadCount: 0,
		Rating:       0.0,
		RatingCount:  0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. 保存到数据库
	if err := s.storeRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create bot store item: %w", err)
	}

	// 6. 增加分类的Bot数量
	if req.Category != "" {
		_ = s.categoryRepo.IncrementBotCount(ctx, req.Category)
	}

	return item, nil
}

// UnpublishBot 下架Bot
func (s *botStorePublisher) UnpublishBot(ctx context.Context, itemID, userID string) error {
	// 1. 获取商店项目
	item, err := s.storeRepo.GetByID(ctx, itemID)
	if err != nil {
		return errno.ErrBotStoreItemNotFound
	}

	// 2. 验证权限
	if item.PublisherID != userID {
		return errno.ErrPermissionDenied
	}

	// 3. 检查状态
	if item.Status != entity.BotStoreItemStatusPublished {
		return errno.ErrInvalidStatus
	}

	// 4. 更新状态为下架
	item.Status = entity.BotStoreItemStatusOffline
	item.UpdatedAt = time.Now()

	if err := s.storeRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to unpublish bot: %w", err)
	}

	// 5. 减少分类的Bot数量
	if item.Category != "" {
		_ = s.categoryRepo.DecrementBotCount(ctx, item.Category)
	}

	return nil
}

// UpdateBotStoreItem 更新商店中的Bot信息
func (s *botStorePublisher) UpdateBotStoreItem(ctx context.Context, req *UpdateBotStoreItemRequest, userID string) error {
	// 1. 获取商店项目
	item, err := s.storeRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return errno.ErrBotStoreItemNotFound
	}

	// 2. 验证权限
	if item.PublisherID != userID {
		return errno.ErrPermissionDenied
	}

	// 3. 检查状态（只有草稿和已拒绝状态可以修改）
	if item.Status != entity.BotStoreItemStatusDraft &&
	   item.Status != entity.BotStoreItemStatusRejected {
		return errno.ErrCannotModifyPublishedItem
	}

	// 4. 更新字段
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Category != "" {
		// 验证分类是否存在
		category, err := s.categoryRepo.GetByID(ctx, req.Category)
		if err != nil || category == nil {
			return errno.ErrInvalidCategory
		}

		// 更新分类计数
		if item.Category != "" && item.Category != req.Category {
			_ = s.categoryRepo.DecrementBotCount(ctx, item.Category)
			_ = s.categoryRepo.IncrementBotCount(ctx, req.Category)
		}
		item.Category = req.Category
	}
	if req.Tags != nil {
		item.Tags = req.Tags
	}
	if req.Price >= 0 {
		item.Price = req.Price
	}
	if req.Screenshots != nil {
		item.Screenshots = req.Screenshots
	}
	if req.Version != "" {
		item.Version = req.Version
	}
	item.UpdatedAt = time.Now()

	// 5. 保存更新
	if err := s.storeRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to update bot store item: %w", err)
	}

	return nil
}

// GetPublishedBots 获取用户已发布的Bot列表
func (s *botStorePublisher) GetPublishedBots(ctx context.Context, userID string, page, pageSize int) ([]*entity.BotStoreItem, int, error) {
	items, total, err := s.storeRepo.List(ctx, &repository.ListRequest{
		Page:     page,
		PageSize: pageSize,
	})

	// 过滤出该用户的Bot
	filteredItems := make([]*entity.BotStoreItem, 0)
	for _, item := range items {
		if item.PublisherID == userID {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems, total, err
}

// validatePublishRequest 验证发布请求
func (s *botStorePublisher) validatePublishRequest(req *PublishBotRequest) error {
	if req.BotID == "" {
		return errno.ErrInvalidBotID
	}
	if req.Name == "" {
		return errno.ErrInvalidBotName
	}
	if req.TenantID == "" {
		return errno.ErrInvalidTenantID
	}
	if req.PublisherID == "" {
		return errno.ErrInvalidPublisherID
	}
	if req.Price < 0 {
		return errno.ErrInvalidPrice
	}
	if req.Category == "" {
		return errno.ErrInvalidCategory
	}
	return nil
}
