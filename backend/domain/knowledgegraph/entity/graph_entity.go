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
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// EntityType 实体类型
type EntityType string

const (
	EntityTypePerson    EntityType = "person"     // 人物
	EntityTypeOrganization EntityType = "organization" // 组织
	EntityTypeLocation  EntityType = "location"   // 地点
	EntityTypeConcept   EntityType = "concept"    // 概念
	EntityTypeEvent     EntityType = "event"      // 事件
	EntityTypeProduct   EntityType = "product"    // 产品
	EntityTypeDocument  EntityType = "document"   // 文档
)

// IsValid 验证实体类型是否有效
func (e EntityType) IsValid() bool {
	switch e {
	case EntityTypePerson, EntityTypeOrganization, EntityTypeLocation,
		EntityTypeConcept, EntityTypeEvent, EntityTypeProduct, EntityTypeDocument:
		return true
	default:
		return false
	}
}

// EntityProperties 实体属性（JSON类型）
type EntityProperties map[string]interface{}

// Scan 实现sql.Scanner接口，用于从数据库读取JSON
func (ep *EntityProperties) Scan(value interface{}) error {
	if value == nil {
		*ep = make(EntityProperties)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal EntityProperties value: %v", value)
	}

	return json.Unmarshal(bytes, ep)
}

// Value 实现driver.Valuer接口，用于向数据库写入JSON
func (ep EntityProperties) Value() (driver.Value, error) {
	if len(ep) == 0 {
		return "{}", nil
	}
	return json.Marshal(ep)
}

// VectorEmbedding 向量嵌入（JSON类型存储）
type VectorEmbedding []float32

// Scan 实现sql.Scanner接口
func (ve *VectorEmbedding) Scan(value interface{}) error {
	if value == nil {
		*ve = make(VectorEmbedding, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal VectorEmbedding value: %v", value)
	}

	return json.Unmarshal(bytes, ve)
}

// Value 实现driver.Valuer接口
func (ve VectorEmbedding) Value() (driver.Value, error) {
	if len(ve) == 0 {
		return "[]", nil
	}
	return json.Marshal(ve)
}

// GraphEntity 图实体
type GraphEntity struct {
	ID         string            `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID   string            `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1;index:idx_tenant_type,priority:1"`
	EntityType EntityType        `json:"entity_type" gorm:"type:varchar(50);not null;index:idx_entity_type;index:idx_tenant_type,priority:2"`
	EntityName string            `json:"entity_name" gorm:"type:varchar(200);not null;index:idx_entity_name"`
	Properties EntityProperties  `json:"properties" gorm:"type:json"`
	Embedding  VectorEmbedding   `json:"embedding" gorm:"type:json"`
	CreatedAt  int64             `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt  int64             `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt  *int64            `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (GraphEntity) TableName() string {
	return "graph_entities"
}

// IsDeleted 是否已删除
func (ge *GraphEntity) IsDeleted() bool {
	return ge.DeletedAt != nil
}

// HasEmbedding 是否有向量嵌入
func (ge *GraphEntity) HasEmbedding() bool {
	return len(ge.Embedding) > 0
}

// GetProperty 获取属性值
func (ge *GraphEntity) GetProperty(key string) (interface{}, bool) {
	if ge.Properties == nil {
		return nil, false
	}
	val, ok := ge.Properties[key]
	return val, ok
}

// SetProperty 设置属性值
func (ge *GraphEntity) SetProperty(key string, value interface{}) {
	if ge.Properties == nil {
		ge.Properties = make(EntityProperties)
	}
	ge.Properties[key] = value
}

// BeforeCreate GORM钩子：创建前
func (ge *GraphEntity) BeforeCreate() error {
	if ge.CreatedAt == 0 {
		ge.CreatedAt = time.Now().Unix()
	}
	if ge.UpdatedAt == 0 {
		ge.UpdatedAt = time.Now().Unix()
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (ge *GraphEntity) BeforeUpdate() error {
	ge.UpdatedAt = time.Now().Unix()
	return nil
}
