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
	"fmt"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/internal/dal/model"
)

// botStoreCategoryRepositoryImpl Bot商店分类仓储实现
type botStoreCategoryRepositoryImpl struct {
	db *gorm.DB
}

// NewBotStoreCategoryRepository 创建Bot商店分类仓储
func NewBotStoreCategoryRepository(db *gorm.DB) BotStoreCategoryRepository {
	return &botStoreCategoryRepositoryImpl{db: db}
}

// Create 创建分类
func (r *botStoreCategoryRepositoryImpl) Create(ctx context.Context, category *entity.BotStoreCategory) error {
	dalCategory := r.categoryEntityToModel(category)
	if err := r.db.WithContext(ctx).Create(dalCategory).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// Update 更新分类
func (r *botStoreCategoryRepositoryImpl) Update(ctx context.Context, category *entity.BotStoreCategory) error {
	dalCategory := r.categoryEntityToModel(category)
	if err := r.db.WithContext(ctx).Save(dalCategory).Error; err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

// GetByID 根据ID获取分类
func (r *botStoreCategoryRepositoryImpl) GetByID(ctx context.Context, categoryID string) (*entity.BotStoreCategory, error) {
	var dalCategory model.BotStoreCategory
	if err := r.db.WithContext(ctx).Where("category_id = ? AND deleted_at IS NULL", categoryID).First(&dalCategory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return r.categoryModelToEntity(&dalCategory)
}

// GetByCode 根据代码获取分类
func (r *botStoreCategoryRepositoryImpl) GetByCode(ctx context.Context, code string) (*entity.BotStoreCategory, error) {
	// 注意：当前实体模型没有 Code 字段，使用 CategoryID 代替
	var dalCategory model.BotStoreCategory
	if err := r.db.WithContext(ctx).Where("category_id = ? AND deleted_at IS NULL", code).First(&dalCategory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return r.categoryModelToEntity(&dalCategory)
}

// List 列出所有分类
func (r *botStoreCategoryRepositoryImpl) List(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	var dalCategories []*model.BotStoreCategory

	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("sort_order ASC, created_at DESC").
		Find(&dalCategories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	// 转换为实体
	categories := make([]*entity.BotStoreCategory, 0, len(dalCategories))
	for _, dalCategory := range dalCategories {
		category, err := r.categoryModelToEntity(dalCategory)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// ListAll 列出所有分类（已废弃，使用 List 替代）
// Deprecated: 使用 List 替代
func (r *botStoreCategoryRepositoryImpl) ListAll(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	return r.List(ctx)
}

// ListActive 列出活跃分类
func (r *botStoreCategoryRepositoryImpl) ListActive(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	var dalCategories []*model.BotStoreCategory

	if err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL AND is_active = ?", true).
		Order("sort_order ASC, created_at DESC").
		Find(&dalCategories).Error; err != nil {
		return nil, fmt.Errorf("failed to list active categories: %w", err)
	}

	// 转换为实体
	categories := make([]*entity.BotStoreCategory, 0, len(dalCategories))
	for _, dalCategory := range dalCategories {
		category, err := r.categoryModelToEntity(dalCategory)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Delete 删除分类
func (r *botStoreCategoryRepositoryImpl) Delete(ctx context.Context, categoryID string) error {
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// IncrementBotCount 增加Bot数量
func (r *botStoreCategoryRepositoryImpl) IncrementBotCount(ctx context.Context, categoryID string) error {
	if err := r.db.WithContext(ctx).Model(&model.BotStoreCategory{}).
		Where("category_id = ?", categoryID).
		UpdateColumn("bot_count", gorm.Expr("bot_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment bot count: %w", err)
	}
	return nil
}

// DecrementBotCount 减少Bot数量
func (r *botStoreCategoryRepositoryImpl) DecrementBotCount(ctx context.Context, categoryID string) error {
	if err := r.db.WithContext(ctx).Model(&model.BotStoreCategory{}).
		Where("category_id = ? AND bot_count > 0", categoryID).
		UpdateColumn("bot_count", gorm.Expr("bot_count - 1")).Error; err != nil {
		return fmt.Errorf("failed to decrement bot count: %w", err)
	}
	return nil
}

// categoryEntityToModel 分类实体转模型
func (r *botStoreCategoryRepositoryImpl) categoryEntityToModel(category *entity.BotStoreCategory) *model.BotStoreCategory {
	return &model.BotStoreCategory{
		CategoryID:  category.CategoryID,
		Name:        category.Name,
		Description: category.Description,
		Icon:        category.Icon,
		SortOrder:   category.SortOrder,
		BotCount:    category.BotCount,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		DeletedAt:   category.DeletedAt,
	}
}

// categoryModelToEntity 分类模型转实体
func (r *botStoreCategoryRepositoryImpl) categoryModelToEntity(dalCategory *model.BotStoreCategory) (*entity.BotStoreCategory, error) {
	return &entity.BotStoreCategory{
		CategoryID:  dalCategory.CategoryID,
		Name:        dalCategory.Name,
		Description: dalCategory.Description,
		Icon:        dalCategory.Icon,
		SortOrder:   dalCategory.SortOrder,
		BotCount:    dalCategory.BotCount,
		CreatedAt:   dalCategory.CreatedAt,
		UpdatedAt:   dalCategory.UpdatedAt,
		DeletedAt:   dalCategory.DeletedAt,
	}, nil
}
