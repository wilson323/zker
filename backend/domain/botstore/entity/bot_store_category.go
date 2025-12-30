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

import "time"

// BotStoreCategory Bot商店分类
type BotStoreCategory struct {
	CategoryID   string    `json:"category_id" gorm:"column:category_id;primaryKey"`
	Name         string    `json:"name" gorm:"column:name;not null;size:100"`
	Icon         string    `json:"icon" gorm:"column:icon;size:255"`
	Description  string    `json:"description" gorm:"column:description;type:text"`
	ParentID     *string   `json:"parent_id,omitempty" gorm:"column:parent_id"`
	SortOrder    int       `json:"sort_order" gorm:"column:sort_order;default:0"`
	IsActive     bool      `json:"is_active" gorm:"column:is_active;default:true"`
	BotCount     int       `json:"bot_count" gorm:"column:bot_count;default:0"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

// TableName 指定表名
func (BotStoreCategory) TableName() string {
	return "bot_store_categories"
}
