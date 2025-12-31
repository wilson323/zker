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

// VisualizationType 可视化类型
type VisualizationType string

const (
	VisualizationTypeHierarchy   VisualizationType = "hierarchy"   // 层级结构图
	VisualizationTypeMatrix      VisualizationType = "matrix"      // 矩阵组织图
	VisualizationTypeDepartment  VisualizationType = "department"  // 部门分布图
	VisualizationTypeNetwork     VisualizationType = "network"     // 关系网络图
)

// OrganizationNode 组织节点（用于可视化）
type OrganizationNode struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Code       string               `json:"code"`
	Type       OrganizationType     `json:"type"`
	Level      int                  `json:"level"`
	ParentID   *string              `json:"parent_id,omitempty"`
	Leader     *EmployeeSummary     `json:"leader,omitempty"`
	EmployeeCount int               `json:"employee_count"`
	Children   []*OrganizationNode  `json:"children,omitempty"`
	Metadata   NodeMetadata         `json:"metadata"`
}

// NodeMetadata 节点元数据
type NodeMetadata struct {
	Color      string                 `json:"color,omitempty"`      // 节点颜色
	Icon       string                 `json:"icon,omitempty"`       // 节点图标
	Size       int                    `json:"size,omitempty"`       // 节点大小
	Shape      string                 `json:"shape,omitempty"`      // 节点形状
	Properties map[string]interface{} `json:"properties,omitempty"` // 自定义属性
}

// VisualizationConfig 可视化配置
type VisualizationConfig struct {
	Type         VisualizationType       `json:"type"`
	Layout       string                  `json:"layout"`       // 布局方式：tree, force, circular
	ShowLeader   bool                    `json:"show_leader"`  // 是否显示负责人
	ShowCount    bool                    `json:"show_count"`   // 是否显示人数
	MaxLevel     int                     `json:"max_level"`    // 最大层级
	Colors       map[string]string       `json:"colors"`       // 颜色配置
	Animations   bool                    `json:"animations"`   // 是否启用动画
	Interactions bool                    `json:"interactions"` // 是否启用交互
}

// VisualizationResult 可视化结果
type VisualizationResult struct {
	Type        VisualizationType        `json:"type"`
	Config      VisualizationConfig      `json:"config"`
	Nodes       []*OrganizationNode      `json:"nodes"`
	Links       []*RelationshipLink      `json:"links,omitempty"`
	Statistics  VisualizationStatistics  `json:"statistics"`
	GeneratedAt int64                    `json:"generated_at"`
}

// RelationshipLink 关系链接（用于关系网络图）
type RelationshipLink struct {
	Source     string                 `json:"source"`     // 源节点ID
	Target     string                 `json:"target"`     // 目标节点ID
	Type       string                 `json:"type"`       // 关系类型
	Weight     float64                `json:"weight"`     // 关系权重
	Label      string                 `json:"label"`      // 关系标签
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// VisualizationStatistics 可视化统计信息
type VisualizationStatistics struct {
	TotalNodes       int              `json:"total_nodes"`
	TotalLinks       int              `json:"total_links"`
	MaxDepth         int              `json:"max_depth"`
	NodeTypeStats    map[string]int   `json:"node_type_stats"`
	LevelStats       map[int]int      `json:"level_stats"`
}

// GetGeneratedAtAsTime 获取生成时间
func (v *VisualizationResult) GetGeneratedAtAsTime() time.Time {
	return time.Unix(v.GeneratedAt/1000, 0)
}

// OrganizationComparison 组织对比结果
type OrganizationComparison struct {
	OrgIDs      []string              `json:"org_ids"`
	Dimensions  []ComparisonDimension  `json:"dimensions"`
	Results     map[string]DimensionResult `json:"results"`
	GeneratedAt int64                 `json:"generated_at"`
}

// ComparisonDimension 对比维度
type ComparisonDimension string

const (
	DimensionSize       ComparisonDimension = "size"        // 规模（人数）
	DimensionCost       ComparisonDimension = "cost"        // 成本
	DimensionUsage      ComparisonDimension = "usage"       // 使用率
	DimensionActivity   ComparisonDimension = "activity"    // 活跃度
	DimensionStructure  ComparisonDimension = "structure"   // 结构
)

// DimensionResult 维度对比结果
type DimensionResult struct {
	Dimension ComparisonDimension      `json:"dimension"`
	Values    map[string]interface{}   `json:"values"`
	Leader    string                   `json:"leader,omitempty"`
	Gap       float64                  `json:"gap,omitempty"`
	Unit      string                   `json:"unit,omitempty"`
}
