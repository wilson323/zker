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

package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/entity"
	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/repository"
)

// GraphVisualizationService 图谱可视化服务
// 职责：生成前端可视化所需的数据
type GraphVisualizationService struct {
	entityRepo       repository.GraphEntityRepository
	relationshipRepo repository.GraphRelationshipRepository
	logger           *zap.Logger
}

// GraphVisualizationServiceConfig 配置
type GraphVisualizationServiceConfig struct {
	EntityRepo       repository.GraphEntityRepository
	RelationshipRepo repository.GraphRelationshipRepository
	Logger           *zap.Logger
}

// NewGraphVisualizationService 创建图谱可视化服务
func NewGraphVisualizationService(config *GraphVisualizationServiceConfig) *GraphVisualizationService {
	return &GraphVisualizationService{
		entityRepo:       config.EntityRepo,
		relationshipRepo: config.RelationshipRepo,
		logger:           config.Logger,
	}
}

// GenerateGraphData 生成可视化数据
func (s *GraphVisualizationService) GenerateGraphData(ctx context.Context, tenantID string, centerEntityID string, depth int) (*entity.VisualizationData, error) {
	if centerEntityID == "" {
		// 生成整个租户的图谱
		return s.generateFullGraph(ctx, tenantID, 100)
	}

	// 生成以中心节点为起点的子图
	return s.generateSubGraph(ctx, centerEntityID, depth)
}

// generateFullGraph 生成完整图谱
func (s *GraphVisualizationService) generateFullGraph(ctx context.Context, tenantID string, limit int) (*entity.VisualizationData, error) {
	visData := &entity.VisualizationData{
		Nodes: make([]*entity.VisNode, 0),
		Links: make([]*entity.VisLink, 0),
	}

	// 1. 获取所有实体
	entities, _, err := s.entityRepo.ListByTenantID(ctx, tenantID, "", limit, 0)
	if err != nil {
		return nil, fmt.Errorf("获取实体失败: %w", err)
	}

	// 转换为可视化节点
	for _, ent := range entities {
		visData.Nodes = append(visData.Nodes, s.convertToVisNode(ent))
	}

	// 2. 获取所有关系
	relationships, _, err := s.relationshipRepo.ListByTenantID(ctx, tenantID, "", limit*2, 0)
	if err != nil {
		return nil, fmt.Errorf("获取关系失败: %w", err)
	}

	// 转换为可视化链接
	for _, rel := range relationships {
		visData.Links = append(visData.Links, s.convertToVisLink(rel))
	}

	s.logger.Debug("完整图谱数据生成完成",
		zap.String("tenant_id", tenantID),
		zap.Int("node_count", len(visData.Nodes)),
		zap.Int("link_count", len(visData.Links)))

	return visData, nil
}

// generateSubGraph 生成子图
func (s *GraphVisualizationService) generateSubGraph(ctx context.Context, centerEntityID string, depth int) (*entity.VisualizationData, error) {
	visData := &entity.VisualizationData{
		Nodes: make([]*entity.VisNode, 0),
		Links: make([]*entity.VisLink, 0),
	}

	// 1. 获取中心节点
	centerNode, err := s.entityRepo.GetByID(ctx, centerEntityID)
	if err != nil {
		return nil, fmt.Errorf("获取中心节点失败: %w", err)
	}

	visData.Nodes = append(visData.Nodes, s.convertToVisNode(centerNode))

	// 2. BFS遍历获取邻居节点
	visited := make(map[string]bool)
	visited[centerEntityID] = true

	currentLevel := []string{centerEntityID}

	for d := 0; d < depth && len(currentLevel) > 0; d++ {
		nextLevel := make([]string, 0)

		for _, nodeID := range currentLevel {
			// 获取邻居关系
			outbound, inbound, err := s.relationshipRepo.ListNeighbors(ctx, nodeID, 50)
			if err != nil {
				s.logger.Warn("获取邻居关系失败", zap.String("node_id", nodeID), zap.Error(err))
				continue
			}

			// 处理出边
			for _, rel := range outbound {
				// 添加链接
				visData.Links = append(visData.Links, s.convertToVisLink(rel))

				// 添加目标节点
				if !visited[rel.TargetEntityID] {
					visited[rel.TargetEntityID] = true
					nextLevel = append(nextLevel, rel.TargetEntityID)

					targetEnt, err := s.entityRepo.GetByID(ctx, rel.TargetEntityID)
					if err == nil {
						visData.Nodes = append(visData.Nodes, s.convertToVisNode(targetEnt))
					}
				}
			}

			// 处理入边
			for _, rel := range inbound {
				visData.Links = append(visData.Links, s.convertToVisLink(rel))

				if !visited[rel.SourceEntityID] {
					visited[rel.SourceEntityID] = true
					nextLevel = append(nextLevel, rel.SourceEntityID)

					sourceEnt, err := s.entityRepo.GetByID(ctx, rel.SourceEntityID)
					if err == nil {
						visData.Nodes = append(visData.Nodes, s.convertToVisNode(sourceEnt))
					}
				}
			}
		}

		currentLevel = nextLevel
	}

	s.logger.Debug("子图数据生成完成",
		zap.String("center_id", centerEntityID),
		zap.Int("node_count", len(visData.Nodes)),
		zap.Int("link_count", len(visData.Links)),
		zap.Int("depth", depth))

	return visData, nil
}

// GeneratePathVisualization 生成路径可视化数据
func (s *GraphVisualizationService) GeneratePathVisualization(ctx context.Context, path *entity.GraphPath) (*entity.VisualizationData, error) {
	visData := &entity.VisualizationData{
		Nodes: make([]*entity.VisNode, 0),
		Links: make([]*entity.VisLink, 0),
	}

	// 1. 获取路径上的所有节点
	entities, err := s.entityRepo.MGetByID(ctx, path.Nodes)
	if err != nil {
		return nil, fmt.Errorf("获取路径节点失败: %w", err)
	}

	// 转换为可视化节点
	for _, ent := range entities {
		visNode := s.convertToVisNode(ent)
		// 路径节点颜色高亮
		visNode.Color = "#FF6B6B"
		visData.Nodes = append(visData.Nodes, visNode)
	}

	// 2. 获取路径上的所有关系
	for _, relID := range path.Relationships {
		rel, err := s.relationshipRepo.GetByID(ctx, relID)
		if err != nil {
			continue
		}

		visLink := s.convertToVisLink(rel)
		// 路径链接颜色高亮
		visLink.Color = "#FF6B6B"
		visData.Links = append(visData.Links, visLink)
	}

	s.logger.Debug("路径可视化数据生成完成",
		zap.Int("node_count", len(visData.Nodes)),
		zap.Int("link_count", len(visData.Links)),
		zap.Int("path_length", path.Length))

	return visData, nil
}

// GenerateForceLayoutData 生成力导向布局数据
func (s *GraphVisualizationService) GenerateForceLayoutData(ctx context.Context, tenantID string, limit int) (*entity.VisualizationData, error) {
	// 获取图谱数据
	visData, err := s.generateFullGraph(ctx, tenantID, limit)
	if err != nil {
		return nil, err
	}

	// 计算节点大小（基于度数）
	nodeDegree := make(map[string]int)
	for _, link := range visData.Links {
		nodeDegree[link.Source]++
		nodeDegree[link.Target]++
	}

	// 设置节点大小
	for _, node := range visData.Nodes {
		degree := nodeDegree[node.ID]
		if degree > 0 {
			// 节点大小与度数成正比
			node.Size = 10 + degree*2
			if node.Size > 50 {
				node.Size = 50
			}
		} else {
			node.Size = 10
		}
	}

	return visData, nil
}

// ========================================
// 内部辅助方法
// ========================================

// convertToVisNode 转换为可视化节点
func (s *GraphVisualizationService) convertToVisNode(ent *entity.GraphEntity) *entity.VisNode {
	node := &entity.VisNode{
		ID:         ent.ID,
		Label:      ent.EntityName,
		Type:       ent.EntityType,
		Size:       20, // 默认大小
		Color:      s.getNodeColor(ent.EntityType),
		Properties: ent.Properties,
	}

	return node
}

// convertToVisLink 转换为可视化链接
func (s *GraphVisualizationService) convertToVisLink(rel *entity.GraphRelationship) *entity.VisLink {
	link := &entity.VisLink{
		Source:     rel.SourceEntityID,
		Target:     rel.TargetEntityID,
		Label:      string(rel.RelationType),
		Type:       rel.RelationType,
		Weight:     rel.Weight,
		Color:      s.getLinkColor(rel.RelationType),
		Properties: rel.Properties,
	}

	return link
}

// getNodeColor 根据实体类型获取节点颜色
func (s *GraphVisualizationService) getNodeColor(entityType entity.EntityType) string {
	colors := map[entity.EntityType]string{
		entity.EntityTypePerson:       "#4ECDC4", // 青色
		entity.EntityTypeOrganization: "#FF6B6B", // 红色
		entity.EntityTypeLocation:     "#95E1D3", // 绿色
		entity.EntityTypeConcept:      "#F38181", // 粉色
		entity.EntityTypeEvent:        "#AA96DA", // 紫色
		entity.EntityTypeProduct:      "#FCBAD3", // 粉紫
		entity.EntityTypeDocument:     "#FFFFD2", // 黄色
	}

	if color, ok := colors[entityType]; ok {
		return color
	}

	return "#CCCCCC" // 默认灰色
}

// getLinkColor 根据关系类型获取链接颜色
func (s *GraphVisualizationService) getLinkColor(relationType entity.RelationType) string {
	// 根据关系类型返回不同的颜色
	// 默认使用浅灰色
	return "#E0E0E0"
}

// GetGraphOverview 获取图谱概览
func (s *GraphVisualizationService) GetGraphOverview(ctx context.Context, tenantID string) (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// 获取实体统计
	entityCount, err := s.entityRepo.CountByTenantID(ctx, tenantID, "")
	if err == nil {
		overview["entity_count"] = entityCount
	}

	// 按类型统计
	typeStats := make(map[string]int64)
	types := []entity.EntityType{
		entity.EntityTypePerson,
		entity.EntityTypeOrganization,
		entity.EntityTypeLocation,
		entity.EntityTypeConcept,
		entity.EntityTypeEvent,
		entity.EntityTypeProduct,
		entity.EntityTypeDocument,
	}

	for _, entityType := range types {
		count, err := s.entityRepo.CountByTenantID(ctx, tenantID, entityType)
		if err == nil {
			typeStats[string(entityType)] = count
		}
	}
	overview["entity_types"] = typeStats

	// 获取关系统计
	relationshipCount, err := s.relationshipRepo.CountByTenantID(ctx, tenantID, "")
	if err == nil {
		overview["relationship_count"] = relationshipCount
	}

	s.logger.Debug("图谱概览获取完成",
		zap.String("tenant_id", tenantID),
		zap.Any("overview", overview))

	return overview, nil
}

// GenerateClusterData 生成聚类数据（简化实现）
func (s *GraphVisualizationService) GenerateClusterData(ctx context.Context, tenantID string) (map[string][]string, error) {
	// 简化实现：按实体类型聚类
	clusters := make(map[string][]string)

	entities, _, err := s.entityRepo.ListByTenantID(ctx, tenantID, "", 10000, 0)
	if err != nil {
		return nil, fmt.Errorf("获取实体失败: %w", err)
	}

	for _, ent := range entities {
		clusterName := string(ent.EntityType)
		clusters[clusterName] = append(clusters[clusterName], ent.ID)
	}

	return clusters, nil
}
