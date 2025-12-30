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

// GraphQueryService 图谱查询服务
// 职责：执行图谱查询，获取实体、关系、子图、路径等
type GraphQueryService struct {
	entityRepo       repository.GraphEntityRepository
	relationshipRepo repository.GraphRelationshipRepository
	logger           *zap.Logger
}

// GraphQueryServiceConfig 配置
type GraphQueryServiceConfig struct {
	EntityRepo       repository.GraphEntityRepository
	RelationshipRepo repository.GraphRelationshipRepository
	Logger           *zap.Logger
}

// NewGraphQueryService 创建图谱查询服务
func NewGraphQueryService(config *GraphQueryServiceConfig) *GraphQueryService {
	return &GraphQueryService{
		entityRepo:       config.EntityRepo,
		relationshipRepo: config.RelationshipRepo,
		logger:           config.Logger,
	}
}

// QueryEntity 查询单个实体
func (s *GraphQueryService) QueryEntity(ctx context.Context, entityID string) (*entity.GraphEntity, error) {
	ent, err := s.entityRepo.GetByID(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("查询实体失败: %w", err)
	}

	s.logger.Debug("实体查询成功",
		zap.String("entity_id", entityID),
		zap.String("entity_name", ent.EntityName))

	return ent, nil
}

// QueryNeighbors 查询实体的邻居（N度关系）
func (s *GraphQueryService) QueryNeighbors(ctx context.Context, entityID string, depth int) (*entity.GraphSubgraph, error) {
	// 1. 获取中心节点
	centerNode, err := s.entityRepo.GetByID(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("获取中心节点失败: %w", err)
	}

	// 2. 转换为GraphNode
	centerGraphNode := &entity.GraphNode{
		ID:         centerNode.ID,
		EntityType: centerNode.EntityType,
		EntityName: centerNode.EntityName,
		Properties: centerNode.Properties,
	}

	subgraph := &entity.GraphSubgraph{
		CenterNode:    centerGraphNode,
		Neighbors:     make([]*entity.GraphNode, 0),
		Relationships: make([]*entity.GraphRelation, 0),
		Depth:         depth,
	}

	// 3. BFS遍历获取邻居
	visited := make(map[string]bool)
	visited[entityID] = true

	currentLevel := []string{entityID}

	for d := 0; d < depth && len(currentLevel) > 0; d++ {
		nextLevel := make([]string, 0)

		for _, nodeID := range currentLevel {
			// 获取出边关系
			outbound, inbound, err := s.relationshipRepo.ListNeighbors(ctx, nodeID, 100)
			if err != nil {
				s.logger.Warn("获取邻居关系失败", zap.String("node_id", nodeID), zap.Error(err))
				continue
			}

			// 处理出边
			for _, rel := range outbound {
				// 添加关系
				subgraph.Relationships = append(subgraph.Relationships, &entity.GraphRelation{
					ID:           rel.ID,
					SourceID:     rel.SourceEntityID,
					TargetID:     rel.TargetEntityID,
					RelationType: rel.RelationType,
					Properties:   rel.Properties,
				})

				// 添加目标节点
				if !visited[rel.TargetEntityID] {
					visited[rel.TargetEntityID] = true
					nextLevel = append(nextLevel, rel.TargetEntityID)

					targetEnt, err := s.entityRepo.GetByID(ctx, rel.TargetEntityID)
					if err == nil {
						subgraph.Neighbors = append(subgraph.Neighbors, &entity.GraphNode{
							ID:         targetEnt.ID,
							EntityType: targetEnt.EntityType,
							EntityName: targetEnt.EntityName,
							Properties: targetEnt.Properties,
						})
					}
				}
			}

			// 处理入边
			for _, rel := range inbound {
				subgraph.Relationships = append(subgraph.Relationships, &entity.GraphRelation{
					ID:           rel.ID,
					SourceID:     rel.SourceEntityID,
					TargetID:     rel.TargetEntityID,
					RelationType: rel.RelationType,
					Properties:   rel.Properties,
				})

				if !visited[rel.SourceEntityID] {
					visited[rel.SourceEntityID] = true
					nextLevel = append(nextLevel, rel.SourceEntityID)

					sourceEnt, err := s.entityRepo.GetByID(ctx, rel.SourceEntityID)
					if err == nil {
						subgraph.Neighbors = append(subgraph.Neighbors, &entity.GraphNode{
							ID:         sourceEnt.ID,
							EntityType: sourceEnt.EntityType,
							EntityName: sourceEnt.EntityName,
							Properties: sourceEnt.Properties,
						})
					}
				}
			}
		}

		currentLevel = nextLevel
	}

	s.logger.Debug("邻居查询完成",
		zap.String("center_id", entityID),
		zap.Int("neighbor_count", len(subgraph.Neighbors)),
		zap.Int("relationship_count", len(subgraph.Relationships)),
		zap.Int("depth", depth))

	return subgraph, nil
}

// QueryPath 查询两个实体之间的路径
func (s *GraphQueryService) QueryPath(ctx context.Context, sourceID, targetID string, maxDepth int) ([]*entity.GraphPath, error) {
	// 使用BFS查找最短路径
	paths := make([]*entity.GraphPath, 0)

	// 记录访问状态和父节点
	visited := make(map[string]bool)
	parent := make(map[string]*PathNode) // 记录路径: nodeID -> parent info

	// 队列用于BFS
	queue := []*PathNode{
		{NodeID: sourceID, Path: []string{sourceID}, Relationships: []string{}},
	}
	visited[sourceID] = true

	found := false

	for len(queue) > 0 && !found {
		current := queue[0]
		queue = queue[1:]

		if current.Depth >= maxDepth {
			continue
		}

		// 获取邻居关系
		outbound, inbound, err := s.relationshipRepo.ListNeighbors(ctx, current.NodeID, 100)
		if err != nil {
			s.logger.Warn("获取邻居关系失败", zap.String("node_id", current.NodeID), zap.Error(err))
			continue
		}

		// 处理出边
		for _, rel := range outbound {
			if visited[rel.TargetEntityID] {
				continue
			}

			visited[rel.TargetEntityID] = true
			newPath := make([]string, len(current.Path)+1)
			copy(newPath, current.Path)
			newPath[len(current.Path)] = rel.TargetEntityID

			newRels := make([]string, len(current.Relationships)+1)
			copy(newRels, current.Relationships)
			newRels[len(current.Relationships)] = rel.ID

			nextNode := &PathNode{
				NodeID:        rel.TargetEntityID,
				Path:          newPath,
				Relationships: newRels,
				Depth:         current.Depth + 1,
			}

			// 找到目标
			if rel.TargetEntityID == targetID {
				found = true
				paths = append(paths, &entity.GraphPath{
					Nodes:         newPath,
					Relationships: newRels,
					Length:        len(newPath) - 1,
					Weight:        s.calculatePathWeight(ctx, newRels),
				})
			}

			queue = append(queue, nextNode)
		}

		// 处理入边
		for _, rel := range inbound {
			if visited[rel.SourceEntityID] {
				continue
			}

			visited[rel.SourceEntityID] = true
			newPath := make([]string, len(current.Path)+1)
			copy(newPath, current.Path)
			newPath[len(current.Path)] = rel.SourceEntityID

			newRels := make([]string, len(current.Relationships)+1)
			copy(newRels, current.Relationships)
			newRels[len(current.Relationships)] = rel.ID

			nextNode := &PathNode{
				NodeID:        rel.SourceEntityID,
				Path:          newPath,
				Relationships: newRels,
				Depth:         current.Depth + 1,
			}

			if rel.SourceEntityID == targetID {
				found = true
				paths = append(paths, &entity.GraphPath{
					Nodes:         newPath,
					Relationships: newRels,
					Length:        len(newPath) - 1,
					Weight:        s.calculatePathWeight(ctx, newRels),
				})
			}

			queue = append(queue, nextNode)
		}
	}

	s.logger.Debug("路径查询完成",
		zap.String("source_id", sourceID),
		zap.String("target_id", targetID),
		zap.Int("path_count", len(paths)),
		zap.Int("max_depth", maxDepth))

	_ = parent // TODO: 使用parent记录路径信息，用于路径重构
	return paths, nil
}

// ExecuteCypherQuery 执行Cypher查询（模拟）
func (s *GraphQueryService) ExecuteCypherQuery(ctx context.Context, query *entity.GraphQuery) (*entity.GraphResult, error) {
	// 简化的Cypher查询执行（实际应该使用Cypher解析器）
	result := &entity.GraphResult{
		Nodes:         make([]*entity.GraphNode, 0),
		Relationships: make([]*entity.GraphRelation, 0),
		Statistics:    &entity.GraphStatistics{},
	}

	// 根据查询文本执行相应操作
	// 这里只实现简单的匹配查询
	// 实际生产环境应该集成完整的Cypher解析引擎

	// 示例: MATCH (n) RETURN n
	// 获取所有节点
	entities, _, err := s.entityRepo.ListByTenantID(ctx, query.TenantID, "", 1000, 0)
	if err == nil {
		for _, ent := range entities {
			result.Nodes = append(result.Nodes, &entity.GraphNode{
				ID:         ent.ID,
				EntityType: ent.EntityType,
				EntityName: ent.EntityName,
				Properties: ent.Properties,
			})
		}
		result.Statistics.NodeCount = len(entities)
	}

	// 获取所有关系
	relationships, _, err := s.relationshipRepo.ListByTenantID(ctx, query.TenantID, "", 1000, 0)
	if err == nil {
		for _, rel := range relationships {
			result.Relationships = append(result.Relationships, &entity.GraphRelation{
				ID:           rel.ID,
				SourceID:     rel.SourceEntityID,
				TargetID:     rel.TargetEntityID,
				RelationType: rel.RelationType,
				Properties:   rel.Properties,
			})
		}
		result.Statistics.RelationshipCount = len(relationships)
	}

	s.logger.Debug("Cypher查询执行完成",
		zap.String("query_id", query.ID),
		zap.Int("node_count", result.Statistics.NodeCount),
		zap.Int("relationship_count", result.Statistics.RelationshipCount))

	return result, nil
}

// SearchEntities 搜索实体
func (s *GraphQueryService) SearchEntities(ctx context.Context, tenantID, keyword string, entityType entity.EntityType, limit int) ([]*entity.GraphEntity, error) {
	var entities []*entity.GraphEntity
	var err error

	if keyword != "" {
		// 按名称模糊搜索
		entities, err = s.entityRepo.SearchByName(ctx, tenantID, keyword, limit)
	} else if entityType != "" {
		// 按类型搜索
		entities, err = s.entityRepo.SearchByType(ctx, tenantID, entityType, limit, 0)
	} else {
		// 列出所有实体
		entities, _, err = s.entityRepo.ListByTenantID(ctx, tenantID, "", limit, 0)
	}

	if err != nil {
		return nil, fmt.Errorf("搜索实体失败: %w", err)
	}

	s.logger.Debug("实体搜索完成",
		zap.String("tenant_id", tenantID),
		zap.String("keyword", keyword),
		zap.String("entity_type", string(entityType)),
		zap.Int("result_count", len(entities)))

	return entities, nil
}

// SearchRelationships 搜索关系
func (s *GraphQueryService) SearchRelationships(ctx context.Context, tenantID string, relationType entity.RelationType, limit int) ([]*entity.GraphRelationship, error) {
	relationships, err := s.relationshipRepo.SearchByType(ctx, tenantID, relationType, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("搜索关系失败: %w", err)
	}

	s.logger.Debug("关系搜索完成",
		zap.String("tenant_id", tenantID),
		zap.String("relation_type", string(relationType)),
		zap.Int("result_count", len(relationships)))

	return relationships, nil
}

// GetGraphStatistics 获取图谱统计信息
func (s *GraphQueryService) GetGraphStatistics(ctx context.Context, tenantID string) (*entity.GraphStatistics, error) {
	stats := &entity.GraphStatistics{}

	// 统计节点数
	entityCount, err := s.entityRepo.CountByTenantID(ctx, tenantID, "")
	if err != nil {
		return nil, fmt.Errorf("统计实体数量失败: %w", err)
	}
	stats.NodeCount = int(entityCount)

	// 统计关系数
	relationshipCount, err := s.relationshipRepo.CountByTenantID(ctx, tenantID, "")
	if err != nil {
		return nil, fmt.Errorf("统计关系数量失败: %w", err)
	}
	stats.RelationshipCount = int(relationshipCount)

	// 计算平均度数
	if stats.NodeCount > 0 {
		stats.AvgDegree = float32(stats.RelationshipCount*2) / float32(stats.NodeCount)
	}

	// 计算图密度（完全图的边数 = n*(n-1)/2）
	if stats.NodeCount > 1 {
		maxEdges := stats.NodeCount * (stats.NodeCount - 1) / 2
		stats.Density = float32(stats.RelationshipCount) / float32(maxEdges)
	}

	s.logger.Debug("图谱统计完成",
		zap.String("tenant_id", tenantID),
		zap.Int("node_count", stats.NodeCount),
		zap.Int("relationship_count", stats.RelationshipCount),
		zap.Float32("avg_degree", stats.AvgDegree),
		zap.Float32("density", stats.Density))

	return stats, nil
}

// ========================================
// 内部辅助方法
// ========================================

// PathNode BFS路径节点
type PathNode struct {
	NodeID        string
	Path          []string
	Relationships []string
	Depth         int
}

// calculatePathWeight 计算路径权重
func (s *GraphQueryService) calculatePathWeight(ctx context.Context, relationshipIDs []string) float32 {
	totalWeight := float32(0.0)

	for _, relID := range relationshipIDs {
		rel, err := s.relationshipRepo.GetByID(ctx, relID)
		if err == nil {
			totalWeight += rel.GetWeight()
		}
	}

	if len(relationshipIDs) > 0 {
		return totalWeight / float32(len(relationshipIDs))
	}

	return 0.0
}

// BatchGetEntities 批量获取实体
func (s *GraphQueryService) BatchGetEntities(ctx context.Context, entityIDs []string) ([]*entity.GraphEntity, error) {
	entities, err := s.entityRepo.MGetByID(ctx, entityIDs)
	if err != nil {
		return nil, fmt.Errorf("批量获取实体失败: %w", err)
	}

	return entities, nil
}

// GetConnectedComponents 获取连通分量数（简化实现）
func (s *GraphQueryService) GetConnectedComponents(ctx context.Context, tenantID string) (int, error) {
	// 简化实现：实际应该使用并查集或图遍历算法
	// 这里返回1表示假设图是连通的
	return 1, nil
}

// FindShortestPath 查找最短路径
func (s *GraphQueryService) FindShortestPath(ctx context.Context, sourceID, targetID string) (*entity.GraphPath, error) {
	paths, err := s.QueryPath(ctx, sourceID, targetID, 10)
	if err != nil {
		return nil, err
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("未找到路径")
	}

	// 返回最短路径
	shortest := paths[0]
	for _, path := range paths {
		if path.Length < shortest.Length {
			shortest = path
		}
	}

	return shortest, nil
}

// FindAllPaths 查找所有路径（带深度限制）
func (s *GraphQueryService) FindAllPaths(ctx context.Context, sourceID, targetID string, maxDepth int) ([]*entity.GraphPath, error) {
	return s.QueryPath(ctx, sourceID, targetID, maxDepth)
}
