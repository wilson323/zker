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

package dal

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/internal/dal/model"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
)

// botStoreRepository Bot商店仓储实现
type botStoreRepository struct {
	db *gorm.DB
}

// NewBotStoreRepository 创建Bot商店仓储
func NewBotStoreRepository(db *gorm.DB) repository.BotStoreRepository {
	return &botStoreRepository{db: db}
}

// Create 创建商店项目
func (r *botStoreRepository) Create(ctx context.Context, item *entity.BotStoreItem) error {
	dalItem := r.entityToModel(item)
	if err := r.db.WithContext(ctx).Create(dalItem).Error; err != nil {
		return fmt.Errorf("failed to create bot store item: %w", err)
	}
	return nil
}

// Update 更新商店项目
func (r *botStoreRepository) Update(ctx context.Context, item *entity.BotStoreItem) error {
	dalItem := r.entityToModel(item)
	if err := r.db.WithContext(ctx).Save(dalItem).Error; err != nil {
		return fmt.Errorf("failed to update bot store item: %w", err)
	}
	return nil
}

// GetByID 根据ID获取商店项目
func (r *botStoreRepository) GetByID(ctx context.Context, itemID string) (*entity.BotStoreItem, error) {
	var dalItem model.BotStoreItem
	if err := r.db.WithContext(ctx).Where("item_id = ? AND deleted_at IS NULL", itemID).First(&dalItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bot store item not found")
		}
		return nil, fmt.Errorf("failed to get bot store item: %w", err)
	}
	return r.modelToEntity(&dalItem)
}

// GetByBotID 根据BotID获取商店项目
func (r *botStoreRepository) GetByBotID(ctx context.Context, botID string) (*entity.BotStoreItem, error) {
	var dalItem model.BotStoreItem
	if err := r.db.WithContext(ctx).Where("bot_id = ? AND deleted_at IS NULL", botID).First(&dalItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bot store item not found")
		}
		return nil, fmt.Errorf("failed to get bot store item by bot_id: %w", err)
	}
	return r.modelToEntity(&dalItem)
}

// Delete 删除商店项目（软删除）
func (r *botStoreRepository) Delete(ctx context.Context, itemID string) error {
	if err := r.db.WithContext(ctx).Where("item_id = ?", itemID).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to delete bot store item: %w", err)
	}
	return nil
}

// List 列出商店项目（分页）
func (r *botStoreRepository) List(ctx context.Context, req *repository.ListRequest) ([]*entity.BotStoreItem, int, error) {
	var dalItems []*model.BotStoreItem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BotStoreItem{}).Where("deleted_at IS NULL")

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		// 默认只显示已发布的
		query = query.Where("status = ?", entity.BotStoreItemStatusPublished)
	}

	// 分类过滤
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 排序
	switch req.SortBy {
	case "popular":
		query = query.Order("download_count DESC, view_count DESC")
	case "latest":
		query = query.Order("created_at DESC")
	case "rating":
		query = query.Order("rating DESC, rating_count DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count bot store items: %w", err)
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&dalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list bot store items: %w", err)
	}

	// 转换为实体
	items := make([]*entity.BotStoreItem, 0, len(dalItems))
	for _, dalItem := range dalItems {
		item, err := r.modelToEntity(dalItem)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, int(total), nil
}

// Search 搜索商店项目
func (r *botStoreRepository) Search(ctx context.Context, req *repository.SearchRequest) ([]*entity.BotStoreItem, int, error) {
	var dalItems []*model.BotStoreItem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BotStoreItem{}).Where("deleted_at IS NULL")

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		query = query.Where("status = ?", entity.BotStoreItemStatusPublished)
	}

	// 关键词搜索
	if req.Query != "" {
		searchPattern := "%" + req.Query + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// 分类过滤
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	// 价格范围
	if req.PriceMin > 0 {
		query = query.Where("price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		query = query.Where("price <= ?", req.PriceMax)
	}

	// 排序
	query = query.Order("rating DESC, download_count DESC")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// 分页
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&dalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search bot store items: %w", err)
	}

	// 转换为实体
	items := make([]*entity.BotStoreItem, 0, len(dalItems))
	for _, dalItem := range dalItems {
		item, err := r.modelToEntity(dalItem)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, int(total), nil
}

// GetPendingReviews 获取待审核列表
func (r *botStoreRepository) GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error) {
	var dalItems []*model.BotStoreItem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.BotStoreItem{}).Where(
		"status = ? AND deleted_at IS NULL",
		entity.BotStoreItemStatusPending,
	)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count pending reviews: %w", err)
	}

	// 分页
	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&dalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get pending reviews: %w", err)
	}

	// 转换为实体
	items := make([]*entity.BotStoreItem, 0, len(dalItems))
	for _, dalItem := range dalItems {
		item, err := r.modelToEntity(dalItem)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, int(total), nil
}

// IncrementViewCount 增加浏览次数
func (r *botStoreRepository) IncrementViewCount(ctx context.Context, itemID string) error {
	if err := r.db.WithContext(ctx).Model(&model.BotStoreItem{}).
		Where("item_id = ?", itemID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment view count: %w", err)
	}
	return nil
}

// IncrementDownloadCount 增加下载次数
func (r *botStoreRepository) IncrementDownloadCount(ctx context.Context, itemID string) error {
	if err := r.db.WithContext(ctx).Model(&model.BotStoreItem{}).
		Where("item_id = ?", itemID).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment download count: %w", err)
	}
	return nil
}

// UpdateRating 更新评分
func (r *botStoreRepository) UpdateRating(ctx context.Context, itemID string, rating float64) error {
	// 获取当前评分信息
	var item model.BotStoreItem
	if err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&item).Error; err != nil {
		return fmt.Errorf("failed to get item for rating update: %w", err)
	}

	// 计算新的平均评分
	newRatingCount := item.RatingCount + 1
	newRating := (item.Rating*float64(item.RatingCount) + rating) / float64(newRatingCount)

	// 更新评分
	if err := r.db.WithContext(ctx).Model(&model.BotStoreItem{}).
		Where("item_id = ?", itemID).
		Updates(map[string]interface{}{
			"rating":       newRating,
			"rating_count": newRatingCount,
		}).Error; err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}

	return nil
}

// entityToModel 实体转模型
func (r *botStoreRepository) entityToModel(item *entity.BotStoreItem) *model.BotStoreItem {
	tagsJSON, _ := json.Marshal(item.Tags)
	screenshotsJSON, _ := json.Marshal(item.Screenshots)

	return &model.BotStoreItem{
		ItemID:        item.ItemID,
		BotID:         item.BotID,
		TenantID:      item.TenantID,
		Name:          item.Name,
		Description:   item.Description,
		Category:      item.Category,
		Tags:          string(tagsJSON),
		Price:         item.Price,
		PublisherID:   item.PublisherID,
		Status:        string(item.Status),
		ViewCount:     item.ViewCount,
		DownloadCount: item.DownloadCount,
		Rating:        item.Rating,
		RatingCount:   item.RatingCount,
		Screenshots:   string(screenshotsJSON),
		Version:       item.Version,
		RejectReason:  item.RejectReason,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		DeletedAt:     item.DeletedAt,
	}
}

// modelToEntity 模型转实体
func (r *botStoreRepository) modelToEntity(dalItem *model.BotStoreItem) (*entity.BotStoreItem, error) {
	var tags []string
	var screenshots []string

	if dalItem.Tags != "" {
		if err := json.Unmarshal([]byte(dalItem.Tags), &tags); err != nil {
			tags = []string{}
		}
	}

	if dalItem.Screenshots != "" {
		if err := json.Unmarshal([]byte(dalItem.Screenshots), &screenshots); err != nil {
			screenshots = []string{}
		}
	}

	return &entity.BotStoreItem{
		ItemID:        dalItem.ItemID,
		BotID:         dalItem.BotID,
		TenantID:      dalItem.TenantID,
		Name:          dalItem.Name,
		Description:   dalItem.Description,
		Category:      dalItem.Category,
		Tags:          tags,
		Price:         dalItem.Price,
		PublisherID:   dalItem.PublisherID,
		Status:        entity.BotStoreItemStatus(dalItem.Status),
		ViewCount:     dalItem.ViewCount,
		DownloadCount: dalItem.DownloadCount,
		Rating:        dalItem.Rating,
		RatingCount:   dalItem.RatingCount,
		Screenshots:   screenshots,
		Version:       dalItem.Version,
		RejectReason:  dalItem.RejectReason,
		CreatedAt:     dalItem.CreatedAt,
		UpdatedAt:     dalItem.UpdatedAt,
		DeletedAt:     dalItem.DeletedAt,
	}, nil
}
