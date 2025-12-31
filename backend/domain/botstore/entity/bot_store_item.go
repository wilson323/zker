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

package entity

import (
	"errors"
	"fmt"
	"time"
)

// BotStoreItemStatus 商店项目状态
type BotStoreItemStatus string

const (
	BotStoreItemStatusDraft     BotStoreItemStatus = "draft"     // 草稿
	BotStoreItemStatusPending   BotStoreItemStatus = "pending"   // 待审核
	BotStoreItemStatusPublished BotStoreItemStatus = "published" // 已发布
	BotStoreItemStatusRejected  BotStoreItemStatus = "rejected"  // 已拒绝
	BotStoreItemStatusOffline   BotStoreItemStatus = "offline"   // 已下架
)

// BotStoreItem Bot商店项目实体
type BotStoreItem struct {
	ItemID         string             `json:"item_id" gorm:"column:item_id;primaryKey"`
	BotID          string             `json:"bot_id" gorm:"column:bot_id;not null"`
	TenantID       string             `json:"tenant_id" gorm:"column:tenant_id;not null"`
	Name           string             `json:"name" gorm:"column:name;not null;size:255"`
	Description    string             `json:"description" gorm:"column:description;type:text"`
	Category       string             `json:"category" gorm:"column:category;size:100"`
	Tags           []string           `json:"tags" gorm:"column:tags;type:json"`
	Price          float64            `json:"price" gorm:"column:price;type:decimal(10,2);default:0.00"`
	PublisherID    string             `json:"publisher_id" gorm:"column:publisher_id;not null"`
	Status         BotStoreItemStatus `json:"status" gorm:"column:status;size:20;default:pending"`
	ViewCount      int                `json:"view_count" gorm:"column:view_count;default:0"`
	DownloadCount  int                `json:"download_count" gorm:"column:download_count;default:0"`
	Rating         float64            `json:"rating" gorm:"column:rating;type:decimal(3,2);default:0.00"`
	RatingCount    int                `json:"rating_count" gorm:"column:rating_count;default:0"`
	Screenshots    []string           `json:"screenshots" gorm:"column:screenshots;type:json"`
	Version        string             `json:"version" gorm:"column:version;size:50"`
	RejectReason   string             `json:"reject_reason" gorm:"column:reject_reason;type:text"`
	CreatedAt      time.Time          `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time          `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

// TableName 指定表名
func (BotStoreItem) TableName() string {
	return "bot_store_items"
}

// CanBePublished 检查是否可以发布
func (b *BotStoreItem) CanBePublished() error {
	if b.Status == BotStoreItemStatusPublished {
		return errors.New("bot already published")
	}
	if b.Status == BotStoreItemStatusPending {
		return errors.New("bot is pending review")
	}
	return nil
}

// IsPublished 检查是否已发布
func (b *BotStoreItem) IsPublished() bool {
	return b.Status == BotStoreItemStatusPublished
}

// CanBeReviewed 检查是否可以审核
func (b *BotStoreItem) CanBeReviewed() error {
	if b.Status != BotStoreItemStatusPending {
		return fmt.Errorf("invalid status for review: %s", b.Status)
	}
	return nil
}
