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

// KnowledgeType 知识类型
type KnowledgeType string

const (
	KnowledgeTypeDocument  KnowledgeType = "DOCUMENT"  // 文档
	KnowledgeTypeFAQ       KnowledgeType = "FAQ"       // 常见问题
	KnowledgeTypeProcedure KnowledgeType = "PROCEDURE" // 流程
	KnowledgeTypeConcept   KnowledgeType = "CONCEPT"   // 概念
)

// KnowledgeMemory 知识记忆实体
type KnowledgeMemory struct {
	ID            int64              `json:"id"`
	MemoryID      string             `json:"memory_id"`
	TenantID      string             `json:"tenant_id"`
	KnowledgeType KnowledgeType      `json:"knowledge_type"`
	Title         string             `json:"title"`
	Content       string             `json:"content"`
	SourceURI     string             `json:"source_uri"`      // 来源URI
	Embedding     []float32          `json:"embedding"`       // 向量嵌入
	Metadata      map[string]interface{} `json:"metadata"`    // 元数据
	QualityScore  float64            `json:"quality_score"`   // 质量评分(0-1)
	Version       int                `json:"version"`        // 版本号
	AccessCount   int                `json:"access_count"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

// GetTenantID 获取租户ID(用于多租户隔离)
func (k *KnowledgeMemory) GetTenantID() string {
	return k.TenantID
}

// HasTenantID 检查是否有租户ID
func (k *KnowledgeMemory) HasTenantID() bool {
	return k.TenantID != ""
}

// IsHighQuality 判断是否为高质量知识
func (k *KnowledgeMemory) IsHighQuality() bool {
	return k.QualityScore >= 0.7
}

// KnowledgeWithScore 带相似度分数的知识
type KnowledgeWithScore struct {
	Knowledge *KnowledgeMemory `json:"knowledge"`
	Score     float64          `json:"score"` // 相似度分数(0-1)
}

// KnowledgeGraph 知识图谱
type KnowledgeGraph struct {
	Nodes    []KnowledgeNode    `json:"nodes"`
	Edges    []KnowledgeEdge    `json:"edges"`
	Stats    GraphStats         `json:"stats"`
}

// KnowledgeNode 知识节点
type KnowledgeNode struct {
	NodeID    string                 `json:"node_id"`
	Type      string                 `json:"type"`      // ENTITY, CONCEPT等
	Label     string                 `json:"label"`     // 节点标签
	Properties map[string]interface{} `json:"properties"` // 节点属性
	MemoryID  string                 `json:"memory_id"` // 关联的记忆ID
}

// KnowledgeEdge 知识关系边
type KnowledgeEdge struct {
	EdgeID        string    `json:"edge_id"`
	SourceNodeID  string    `json:"source_node_id"`
	TargetNodeID  string    `json:"target_node_id"`
	RelationType  string    `json:"relation_type"` // 关系类型: RELATED_TO, PART_OF等
	Strength      float64   `json:"strength"`     // 关系强度(0-1)
	CreatedAt     time.Time `json:"created_at"`
}

// GraphStats 图谱统计
type GraphStats struct {
	NodeCount      int     `json:"node_count"`
	EdgeCount      int     `json:"edge_count"`
	Density        float64 `json:"density"`         // 图密度
	AvgDegree      float64 `json:"avg_degree"`      // 平均度数
	MaxDegree      int     `json:"max_degree"`      // 最大度数
	ConnectedComponents int  `json:"connected_components"` // 连通分量数
}

// KnowledgeAssociation 知识关联
type KnowledgeAssociation struct {
	ID               int64     `json:"id"`
	SourceMemoryID   string    `json:"source_memory_id"`
	TargetMemoryID   string    `json:"target_memory_id"`
	AssociationType  string    `json:"association_type"` // 关联类型
	Strength         float64   `json:"strength"`         // 关联强度(0-1)
	CreatedAt        time.Time `json:"created_at"`
}

// GetMetadataJSON 获取元数据JSON
func (k *KnowledgeMemory) GetMetadataJSON() (string, error) {
	if k.Metadata == nil {
		return "{}", nil
	}
	data, err := json.Marshal(k.Metadata)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetMetadataFromJSON 从JSON设置元数据
func (k *KnowledgeMemory) SetMetadataFromJSON(jsonStr string) error {
	if jsonStr == "" || jsonStr == "{}" {
		k.Metadata = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &k.Metadata)
}

// DocumentChunk 文档分块(用于文档存储)
type DocumentChunk struct {
	ChunkID    string    `json:"chunk_id"`
	DocumentID string    `json:"document_id"`
	Index      int       `json:"index"` // 分块序号
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding"`
	CreatedAt  time.Time `json:"created_at"`
}
