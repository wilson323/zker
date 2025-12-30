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
	"time"
)

// ConfigItem 配置项
type ConfigItem struct {
	ConfigID    string     `gorm:"primaryKey;column:config_id" json:"config_id"`
	TenantID    string     `gorm:"column:tenant_id;index:idx_tenant_key" json:"tenant_id"`
	ConfigKey   string     `gorm:"column:config_key;index:idx_tenant_key;uniqueIndex:uk_tenant_key" json:"config_key"`
	ConfigValue string     `gorm:"type:text;column:config_value" json:"config_value"`
	ConfigType  string     `gorm:"column:config_type" json:"config_type"` // string, int, float, bool, json
	Description string     `gorm:"type:text;column:description" json:"description"`
	IsEncrypted bool       `gorm:"column:is_encrypted;default:false" json:"is_encrypted"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy   string     `gorm:"column:updated_by" json:"updated_by"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

// TableName 表名
func (ConfigItem) TableName() string {
	return "config_items"
}

// ConfigHistory 配置变更历史
type ConfigHistory struct {
	HistoryID     string     `gorm:"primaryKey;column:history_id" json:"history_id"`
	ConfigID      string     `gorm:"column:config_id;index:idx_config_id" json:"config_id"`
	TenantID      string     `gorm:"column:tenant_id;index:idx_tenant_id" json:"tenant_id"`
	ConfigKey     string     `gorm:"column:config_key" json:"config_key"`
	OldValue      string     `gorm:"type:text;column:old_value" json:"old_value"`
	NewValue      string     `gorm:"type:text;column:new_value" json:"new_value"`
	ChangeReason  string     `gorm:"type:text;column:change_reason" json:"change_reason"`
	ChangedBy     string     `gorm:"column:changed_by" json:"changed_by"`
	ChangedAt     time.Time  `gorm:"column:changed_at" json:"changed_at"`
	VersionNumber int        `gorm:"column:version_number" json:"version_number"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

// TableName 表名
func (ConfigHistory) TableName() string {
	return "config_history"
}
