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

// ConversationMemoryDAO 对话记忆DAO
type ConversationMemoryDAO struct {
	ID              int64     `gorm:"primaryKey;autoIncrement;comment:记忆ID"`
	MemoryID        string    `gorm:"column:memory_id;type:varchar(36);not null;uniqueIndex:uk_memory_id;comment:记忆唯一标识"`
	TenantID        string    `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_tenant_user,idx_tenant_user_conv;comment:租户ID"`
	UserID          string    `gorm:"column:user_id;type:varchar(64);not null;index:idx_tenant_user,idx_tenant_user_conv;comment:用户ID"`
	ConversationID  string    `gorm:"column:conversation_id;type:varchar(64);not null;index:idx_conversation,idx_tenant_user_conv;comment:会话ID"`
	MemoryType      string    `gorm:"column:memory_type;type:enum('SUMMARY','ENTITY','PREFERENCE','EVENT');not null;index:idx_type;comment:记忆类型"`
	Content         string    `gorm:"column:content;type:text;not null;comment:记忆内容"`
	Embedding       []byte    `gorm:"column:embedding;type:vector(1536);comment:向量嵌入"`
	ImportanceScore float64   `gorm:"column:importance_score;type:decimal(3,2);default:0.50;index:idx_importance;comment:重要性评分"`
	AccessCount     int       `gorm:"column:access_count;type:int;default:0;comment:访问次数"`
	LastAccessedAt  time.Time `gorm:"column:last_accessed_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:最后访问时间"`
	ExpiresAt       *time.Time `gorm:"column:expires_at;type:timestamp;index:idx_expires;comment:过期时间"`
	Metadata        string    `gorm:"column:metadata;type:json;comment:元数据"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 指定表名
func (ConversationMemoryDAO) TableName() string {
	return "conversation_memories"
}
