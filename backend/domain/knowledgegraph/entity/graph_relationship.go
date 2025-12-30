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

// RelationType 关系类型
type RelationType string

const (
	// 人物关系
	RelationTypeKnows          RelationType = "knows"           // 认识
	RelationTypeWorksFor       RelationType = "works_for"       // 任职于
	RelationTypeColleagueOf    RelationType = "colleague_of"    // 同事
	RelationTypeFriendOf       RelationType = "friend_of"       // 朋友
	RelationTypeFamilyOf       RelationType = "family_of"       // 家人
	RelationTypeSpouseOf       RelationType = "spouse_of"       // 配偶
	RelationTypeChildOf        RelationType = "child_of"        // 子女
	RelationTypeParentOf       RelationType = "parent_of"       // 父母

	// 组织关系
	RelationTypePartOf         RelationType = "part_of"         // 属于
	RelationTypeSubsidiaryOf   RelationType = "subsidiary_of"   // 子公司
	RelationTypePartnerOf      RelationType = "partner_of"      // 合作伙伴
	RelationTypeCompetitorOf   RelationType = "competitor_of"   // 竞争对手
	RelationTypeInvestorIn     RelationType = "investor_in"     // 投资者

	// 地理关系
	RelationTypeLocatedIn      RelationType = "located_in"      // 位于
	RelationTypeContains       RelationType = "contains"        // 包含
	RelationTypeBorders        RelationType = "borders"         // 毗邻

	// 概念关系
	RelationTypeRelatedTo      RelationType = "related_to"      // 相关
	RelationTypeSimilarTo      RelationType = "similar_to"      // 相似
	RelationTypeInstanceOf     RelationType = "instance_of"     // 实例
	RelationTypeSubclassOf     RelationType = "subclass_of"     // 子类

	// 事件关系
	RelationTypeParticipatedIn RelationType = "participated_in" // 参与
	RelationTypeCaused         RelationType = "caused"          // 导致
	RelationTypePreceded       RelationType = "preceded"        // 先于

	// 文档关系
	RelationTypeCites          RelationType = "cites"           // 引用
	RelationTypeReferences     RelationType = "references"      // 参考
	RelationTypeAbout          RelationType = "about"           // 关于
)

// RelationshipProperties 关系属性（JSON类型）
type RelationshipProperties map[string]interface{}

// Scan 实现sql.Scanner接口
func (rp *RelationshipProperties) Scan(value interface{}) error {
	if value == nil {
		*rp = make(RelationshipProperties)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal RelationshipProperties value: %v", value)
	}

	return json.Unmarshal(bytes, rp)
}

// Value 实现driver.Valuer接口
func (rp RelationshipProperties) Value() (driver.Value, error) {
	if len(rp) == 0 {
		return "{}", nil
	}
	return json.Marshal(rp)
}

// GraphRelationship 图关系
type GraphRelationship struct {
	ID             string                   `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID       string                   `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id,priority:1;index:idx_tenant_source,priority:1;index:idx_tenant_target,priority:1;index:idx_tenant_relation,priority:1"`
	SourceEntityID string                   `json:"source_entity_id" gorm:"type:varchar(36);not null;index:idx_source;index:idx_tenant_source,priority:2"`
	TargetEntityID string                   `json:"target_entity_id" gorm:"type:varchar(36);not null;index:idx_target;index:idx_tenant_target,priority:2"`
	RelationType   RelationType             `json:"relation_type" gorm:"type:varchar(50);not null;index:idx_relation;index:idx_tenant_relation,priority:2"`
	Properties     RelationshipProperties   `json:"properties" gorm:"type:json"`
	Weight         float32                  `json:"weight" gorm:"type:float;default:1.0"` // 关系权重（用于图谱推理）
	CreatedAt      int64                    `json:"created_at" gorm:"not null;default:0"`
	DeletedAt      *int64                   `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (GraphRelationship) TableName() string {
	return "graph_relationships"
}

// IsDeleted 是否已删除
func (gr *GraphRelationship) IsDeleted() bool {
	return gr.DeletedAt != nil
}

// GetProperty 获取属性值
func (gr *GraphRelationship) GetProperty(key string) (interface{}, bool) {
	if gr.Properties == nil {
		return nil, false
	}
	val, ok := gr.Properties[key]
	return val, ok
}

// SetProperty 设置属性值
func (gr *GraphRelationship) SetProperty(key string, value interface{}) {
	if gr.Properties == nil {
		gr.Properties = make(RelationshipProperties)
	}
	gr.Properties[key] = value
}

// GetWeight 获取权重（如果属性中有weight字段，优先使用）
func (gr *GraphRelationship) GetWeight() float32 {
	if weightVal, ok := gr.GetProperty("weight"); ok {
		if weight, ok := weightVal.(float64); ok {
			return float32(weight)
		}
	}
	return gr.Weight
}

// BeforeCreate GORM钩子：创建前
func (gr *GraphRelationship) BeforeCreate() error {
	if gr.CreatedAt == 0 {
		gr.CreatedAt = time.Now().Unix()
	}
	return nil
}

// IsSymmetric 是否对称关系
func (gr *GraphRelationship) IsSymmetric() bool {
	switch gr.RelationType {
	case RelationTypeKnows, RelationTypeColleagueOf, RelationTypeFriendOf,
		RelationTypePartnerOf, RelationTypeCompetitorOf, RelationTypeRelatedTo,
		RelationTypeSimilarTo, RelationTypeBorders:
		return true
	default:
		return false
	}
}

// IsValid 验证关系是否有效
func (gr *GraphRelationship) IsValid() bool {
	// 检查源和目标实体ID不能相同
	if gr.SourceEntityID == gr.TargetEntityID {
		return false
	}

	// 检查关系类型是否有效
	switch gr.RelationType {
	case RelationTypeKnows, RelationTypeWorksFor, RelationTypeColleagueOf,
		RelationTypeFriendOf, RelationTypeFamilyOf, RelationTypeSpouseOf,
		RelationTypeChildOf, RelationTypeParentOf, RelationTypePartOf,
		RelationTypeSubsidiaryOf, RelationTypePartnerOf, RelationTypeCompetitorOf,
		RelationTypeInvestorIn, RelationTypeLocatedIn, RelationTypeContains,
		RelationTypeBorders, RelationTypeRelatedTo, RelationTypeSimilarTo,
		RelationTypeInstanceOf, RelationTypeSubclassOf, RelationTypeParticipatedIn,
		RelationTypeCaused, RelationTypePreceded, RelationTypeCites,
		RelationTypeReferences, RelationTypeAbout:
		return true
	default:
		return false
	}
}
