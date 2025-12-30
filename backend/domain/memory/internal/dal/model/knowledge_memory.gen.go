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

// KnowledgeMemoryDAO 知识记忆DAO
type KnowledgeMemoryDAO struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;comment:知识ID"`
	MemoryID     string    `gorm:"column:memory_id;type:varchar(36);not null;uniqueIndex:uk_memory_id;comment:知识唯一标识"`
	TenantID     string    `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_tenant_type,idx_tenant_quality;comment:租户ID"`
	KnowledgeType string   `gorm:"column:knowledge_type;type:enum('DOCUMENT','FAQ','PROCEDURE','CONCEPT');not null;index:idx_tenant_type;comment:知识类型"`
	Title        string    `gorm:"column:title;type:varchar(255);not null;comment:知识标题"`
	Content      string    `gorm:"column:content;type:longtext;not null;comment:知识内容"`
	SourceURI    string    `gorm:"column:source_uri;type:varchar(512);comment:来源URI"`
	Embedding    []byte    `gorm:"column:embedding;type:vector(1536);comment:向量嵌入"`
	Metadata     string    `gorm:"column:metadata;type:json;comment:元数据"`
	QualityScore float64   `gorm:"column:quality_score;type:decimal(3,2);default:0.50;index:idx_tenant_quality;comment:质量评分"`
	Version      int       `gorm:"column:version;type:int;default:1;comment:版本号"`
	AccessCount  int       `gorm:"column:access_count;type:int;default:0;comment:访问次数"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间"`
}

// TableName 指定表名
func (KnowledgeMemoryDAO) TableName() string {
	return "knowledge_memories"
}

// MemoryAssociationDAO 知识关联DAO
type MemoryAssociationDAO struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;comment:关联ID"`
	SourceMemoryID   string    `gorm:"column:source_memory_id;type:varchar(36);not null;uniqueIndex:uk_source_target,index:idx_source;comment:源记忆ID"`
	TargetMemoryID   string    `gorm:"column:target_memory_id;type:varchar(36);not null;uniqueIndex:uk_source_target,index:idx_target;comment:目标记忆ID"`
	AssociationType  string    `gorm:"column:association_type;type:varchar(64);not null;index:idx_type;comment:关联类型"`
	Strength         float64   `gorm:"column:strength;type:decimal(3,2);default:0.50;comment:关联强度"`
	Metadata         string    `gorm:"column:metadata;type:json;comment:元数据"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

// TableName 指定表名
func (MemoryAssociationDAO) TableName() string {
	return "memory_associations"
}
