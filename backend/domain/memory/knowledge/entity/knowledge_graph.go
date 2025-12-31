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
	"encoding/json"
	"time"
)

// GraphEntity 图谱实体
type GraphEntity struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`       // PERSON, ORG, LOCATION, CONCEPT等
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
}

// GraphRelation 图谱关系
type GraphRelation struct {
	ID         string                 `json:"id"`
	FromEntity string                 `json:"from_entity"`
	ToEntity   string                 `json:"to_entity"`
	RelationType string               `json:"relation_type"` // located_in, founded_by, works_for等
	Properties map[string]interface{} `json:"properties"`
	Confidence float64                `json:"confidence"`   // 关系置信度(0-1)
	CreatedAt  time.Time              `json:"created_at"`
}

// ExtractEntitiesRequest 实体抽取请求
type ExtractEntitiesRequest struct {
	Text           string `json:"text"`
	ConversationID string `json:"conversation_id"`
	TenantID       string `json:"tenant_id"`
	UserID         string `json:"user_id"`
}

// ExtractEntitiesResponse 实体抽取响应
type ExtractEntitiesResponse struct {
	Entities   []GraphEntity `json:"entities"`
	Relations  []GraphRelation `json:"relations"`
	Confidence float64        `json:"confidence"`
	ExtractedAt time.Time     `json:"extracted_at"`
}

// QueryKnowledgeGraphRequest 知识图谱查询请求
type QueryKnowledgeGraphRequest struct {
	Entity        string   `json:"entity"`          // 起始实体
	RelationTypes []string `json:"relation_types"`  // 关系类型过滤
	MaxHops       int      `json:"max_hops"`        // 最大跳数
	Limit         int      `json:"limit"`           // 结果数量限制
	ConversationID string  `json:"conversation_id"`
}

// QueryKnowledgeGraphResponse 知识图谱查询响应
type QueryKnowledgeGraphResponse struct {
	Entities  []GraphEntity     `json:"entities"`
	Relations []GraphRelation   `json:"relations"`
	Paths     []GraphPath       `json:"paths"`
	QueryTime int64             `json:"query_time_ms"`
}

// GraphPath 图谱路径
type GraphPath struct {
	Entities []GraphEntity   `json:"entities"`
	Relations []GraphRelation `json:"relations"`
	Length   int             `json:"length"`
	Score    float64         `json:"score"`
}

// VisualizeKnowledgeGraphRequest 可视化知识图谱请求
type VisualizeKnowledgeGraphRequest struct {
	ConversationID string `json:"conversation_id"`
	MaxNodes       int    `json:"max_nodes"`
	MinConfidence  float64 `json:"min_confidence"`
}

// VisualizeKnowledgeGraphResponse 可视化知识图谱响应
type VisualizeKnowledgeGraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
	Stats GraphStats  `json:"stats"`
}

// GraphNode 图谱节点(用于可视化)
type GraphNode struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Size       int                    `json:"size"`    // 节点大小
	Color      string                 `json:"color"`   // 节点颜色
}

// GraphEdge 图谱边(用于可视化)
type GraphEdge struct {
	ID          string  `json:"id"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Label       string  `json:"label"`
	Weight      float64 `json:"weight"`  // 边权重
	Confidence  float64 `json:"confidence"`
}

// GraphStats 图谱统计
type GraphStats struct {
	TotalNodes    int     `json:"total_nodes"`
	TotalEdges    int     `json:"total_edges"`
	AvgDegree     float64 `json:"avg_degree"`
	MaxDegree     int     `json:"max_degree"`
	ClusterCount  int     `json:"cluster_count"`
	Density       float64 `json:"density"`
}

// GetPropertiesJSON 获取属性JSON
func (e *GraphEntity) GetPropertiesJSON() (string, error) {
	if e.Properties == nil {
		return "{}", nil
	}
	data, err := json.Marshal(e.Properties)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetPropertiesFromJSON 从JSON设置属性
func (e *GraphEntity) SetPropertiesFromJSON(jsonStr string) error {
	if jsonStr == "" || jsonStr == "{}" {
		e.Properties = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &e.Properties)
}

// GetPropertiesJSON 获取关系属性JSON
func (r *GraphRelation) GetPropertiesJSON() (string, error) {
	if r.Properties == nil {
		return "{}", nil
	}
	data, err := json.Marshal(r.Properties)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetPropertiesFromJSON 从JSON设置关系属性
func (r *GraphRelation) SetPropertiesFromJSON(jsonStr string) error {
	if jsonStr == "" || jsonStr == "{}" {
		r.Properties = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &r.Properties)
}
