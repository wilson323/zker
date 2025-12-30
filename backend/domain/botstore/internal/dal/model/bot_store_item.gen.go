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

package model

import (
	"time"
)

// BotStoreItem bot_store_items表模型
type BotStoreItem struct {
	ItemID        string    `gorm:"column:item_id;primaryKey" db:"item_id"`
	BotID         string    `gorm:"column:bot_id" db:"bot_id"`
	TenantID      string    `gorm:"column:tenant_id" db:"tenant_id"`
	Name          string    `gorm:"column:name" db:"name"`
	Description   string    `gorm:"column:description" db:"description"`
	Category      string    `gorm:"column:category" db:"category"`
	Tags          string    `gorm:"column:tags" db:"tags"` // JSON字符串
	Price         float64   `gorm:"column:price" db:"price"`
	PublisherID   string    `gorm:"column:publisher_id" db:"publisher_id"`
	Status        string    `gorm:"column:status" db:"status"`
	ViewCount     int       `gorm:"column:view_count" db:"view_count"`
	DownloadCount int       `gorm:"column:download_count" db:"download_count"`
	Rating        float64   `gorm:"column:rating" db:"rating"`
	RatingCount   int       `gorm:"column:rating_count" db:"rating_count"`
	Screenshots   string    `gorm:"column:screenshots" db:"screenshots"` // JSON字符串
	Version       string    `gorm:"column:version" db:"version"`
	RejectReason  string    `gorm:"column:reject_reason" db:"reject_reason"`
	CreatedAt     time.Time `gorm:"column:created_at" db:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" db:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" db:"deleted_at"`
}

// TableName 表名
func (BotStoreItem) TableName() string {
	return "bot_store_items"
}

// BotStoreCategory bot_store_categories表模型
type BotStoreCategory struct {
	CategoryID  string     `gorm:"column:category_id;primaryKey" db:"category_id"`
	Name        string     `gorm:"column:name" db:"name"`
	Icon        string     `gorm:"column:icon" db:"icon"`
	Description string     `gorm:"column:description" db:"description"`
	ParentID    *string    `gorm:"column:parent_id" db:"parent_id"`
	SortOrder   int        `gorm:"column:sort_order" db:"sort_order"`
	IsActive    bool       `gorm:"column:is_active" db:"is_active"`
	BotCount    int        `gorm:"column:bot_count" db:"bot_count"`
	CreatedAt   time.Time  `gorm:"column:created_at" db:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" db:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" db:"deleted_at"`
}

// TableName 表名
func (BotStoreCategory) TableName() string {
	return "bot_store_categories"
}
