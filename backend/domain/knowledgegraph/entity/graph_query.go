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

// QueryVariables 查询变量（JSON类型）
type QueryVariables map[string]interface{}

// Scan 实现sql.Scanner接口
func (qv *QueryVariables) Scan(value interface{}) error {
	if value == nil {
		*qv = make(QueryVariables)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal QueryVariables value: %v", value)
	}

	return json.Unmarshal(bytes, qv)
}

// Value 实现driver.Valuer接口
func (qv QueryVariables) Value() (driver.Value, error) {
	if len(qv) == 0 {
		return "{}", nil
	}
	return json.Marshal(qv)
}

// GraphQuery 图查询（存储预定义的Cypher查询）
type GraphQuery struct {
	ID          string         `json:"id" gorm:"primaryKey;type:varchar(36)"`
	TenantID    string         `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	QueryName   string         `json:"query_name" gorm:"type:varchar(100);not null"`
	QueryText   string         `json:"query_text" gorm:"type:text;not null"` // Cypher查询语句
	Description string         `json:"description" gorm:"type:text"`          // 查询描述
	Variables   QueryVariables `json:"variables" gorm:"type:json"`           // 查询变量
	IsPublic    bool           `json:"is_public" gorm:"type:tinyint(1);default:false"` // 是否公开
	CreatedBy   string         `json:"created_by" gorm:"type:varchar(36);not null"` // 创建者ID
	CreatedAt   int64          `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt   int64          `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt   *int64         `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (GraphQuery) TableName() string {
	return "graph_queries"
}

// IsDeleted 是否已删除
func (gq *GraphQuery) IsDeleted() bool {
	return gq.DeletedAt != nil
}

// SetVariable 设置查询变量
func (gq *GraphQuery) SetVariable(key string, value interface{}) {
	if gq.Variables == nil {
		gq.Variables = make(QueryVariables)
	}
	gq.Variables[key] = value
}

// GetVariable 获取查询变量
func (gq *GraphQuery) GetVariable(key string) (interface{}, bool) {
	if gq.Variables == nil {
		return nil, false
	}
	val, ok := gq.Variables[key]
	return val, ok
}

// BeforeCreate GORM钩子：创建前
func (gq *GraphQuery) BeforeCreate() error {
	if gq.CreatedAt == 0 {
		gq.CreatedAt = time.Now().Unix()
	}
	if gq.UpdatedAt == 0 {
		gq.UpdatedAt = time.Now().Unix()
	}
	return nil
}

// BeforeUpdate GORM钩子：更新前
func (gq *GraphQuery) BeforeUpdate() error {
	gq.UpdatedAt = time.Now().Unix()
	return nil
}

// ========================================
// 图谱查询结果相关实体
// ========================================

// GraphResult 图谱查询结果
type GraphResult struct {
	Nodes         []*GraphNode        `json:"nodes"`
	Relationships []*GraphRelation    `json:"relationships"`
	Paths         []*GraphPath        `json:"paths,omitempty"`
	Statistics    *GraphStatistics    `json:"statistics,omitempty"`
}

// GraphNode 图节点
type GraphNode struct {
	ID          string            `json:"id"`
	EntityType  EntityType        `json:"entity_type"`
	EntityName  string            `json:"entity_name"`
	Properties  EntityProperties  `json:"properties"`
	Labels      []string          `json:"labels,omitempty"` // 节点标签
}

// GraphRelation 图关系（用于查询结果）
type GraphRelation struct {
	ID             string                 `json:"id"`
	SourceID       string                 `json:"source_id"`
	TargetID       string                 `json:"target_id"`
	RelationType   RelationType           `json:"relation_type"`
	Properties     RelationshipProperties `json:"properties"`
}

// GraphPath 图路径（用于最短路径查询等）
type GraphPath struct {
	Nodes         []string `json:"nodes"`          // 节点ID列表
	Relationships []string `json:"relationships"`  // 关系ID列表
	Length        int      `json:"length"`         // 路径长度
	Weight        float32  `json:"weight"`         // 路径权重
}

// GraphSubgraph 图子图（用于邻域查询）
type GraphSubgraph struct {
	CenterNode    *GraphNode          `json:"center_node"`
	Neighbors     []*GraphNode        `json:"neighbors"`
	Relationships []*GraphRelation    `json:"relationships"`
	Depth         int                 `json:"depth"`
}

// GraphStatistics 图谱统计信息
type GraphStatistics struct {
	NodeCount         int     `json:"node_count"`
	RelationshipCount int     `json:"relationship_count"`
	AvgDegree         float32 `json:"avg_degree"`          // 平均度数
	Density           float32 `json:"density"`             // 图密度
	ConnectedComponents int   `json:"connected_components"` // 连通分量数
}

// VisualizationData 可视化数据（用于前端渲染）
type VisualizationData struct {
	Nodes []*VisNode `json:"nodes"`
	Links []*VisLink `json:"links"`
}

// VisNode 可视化节点
type VisNode struct {
	ID          string                 `json:"id"`
	Label       string                 `json:"label"`
	Type        EntityType             `json:"type"`
	Size        int                    `json:"size,omitempty"`         // 节点大小
	Color       string                 `json:"color,omitempty"`        // 节点颜色
	Properties  EntityProperties       `json:"properties,omitempty"`
	X           float32                `json:"x,omitempty"`           // X坐标（用于布局）
	Y           float32                `json:"y,omitempty"`           // Y坐标
}

// VisLink 可视化链接
type VisLink struct {
	Source      string                 `json:"source"`               // 源节点ID
	Target      string                 `json:"target"`               // 目标节点ID
	Label       string                 `json:"label,omitempty"`      // 关系标签
	Type        RelationType           `json:"type"`                 // 关系类型
	Weight      float32                `json:"weight,omitempty"`     // 权重
	Color       string                 `json:"color,omitempty"`      // 链接颜色
	Properties  RelationshipProperties `json:"properties,omitempty"`
}
