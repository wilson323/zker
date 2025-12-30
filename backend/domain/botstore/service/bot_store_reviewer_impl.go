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

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// botStoreReviewer Bot商店审核器实现
type botStoreReviewer struct {
	storeRepo repository.BotStoreRepository
}

// NewBotStoreReviewer 创建Bot商店审核器
func NewBotStoreReviewer(storeRepo repository.BotStoreRepository) BotStoreReviewer {
	return &botStoreReviewer{
		storeRepo: storeRepo,
	}
}

// ReviewBot 审核Bot
func (s *botStoreReviewer) ReviewBot(ctx context.Context, req *ReviewBotRequest) error {
	// 1. 获取待审核项目
	item, err := s.storeRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return errno.ErrBotStoreItemNotFound
	}

	// 2. 检查状态
	if err := item.CanBeReviewed(); err != nil {
		return err
	}

	// 3. 更新审核结果
	if req.Approved {
		item.Status = entity.BotStoreItemStatusPublished
		item.RejectReason = ""
	} else {
		item.Status = entity.BotStoreItemStatusRejected
		if req.Reason == "" {
			return errno.ErrRejectReasonRequired
		}
		item.RejectReason = req.Reason
	}
	item.UpdatedAt = time.Now()

	// 4. 保存更新
	if err := s.storeRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("failed to update review result: %w", err)
	}

	return nil
}

// GetPendingReviews 获取待审核列表
func (s *botStoreReviewer) GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error) {
	items, total, err := s.storeRepo.GetPendingReviews(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pending reviews: %w", err)
	}
	return items, total, nil
}

// GetReviewStatistics 获取审核统计
func (s *botStoreReviewer) GetReviewStatistics(ctx context.Context) (*ReviewStatistics, error) {
	// 获取所有状态的数量
	draftItems, _, _ := s.storeRepo.List(ctx, &repository.ListRequest{
		Status:   entity.BotStoreItemStatusDraft,
		Page:     1,
		PageSize: 1,
	})

	pendingItems, pendingTotal, _ := s.storeRepo.GetPendingReviews(ctx, 1, 1)

	publishedItems, publishedTotal, _ := s.storeRepo.List(ctx, &repository.ListRequest{
		Status:   entity.BotStoreItemStatusPublished,
		Page:     1,
		PageSize: 1,
	})

	rejectedItems, rejectedTotal, _ := s.storeRepo.List(ctx, &repository.ListRequest{
		Status:   entity.BotStoreItemStatusRejected,
		Page:     1,
		PageSize: 1,
	})

	return &ReviewStatistics{
		TotalPending:   pendingTotal,
		TotalPublished: publishedTotal,
		TotalRejected:  rejectedTotal,
		TotalDraft:     len(draftItems),
	}, nil
}
