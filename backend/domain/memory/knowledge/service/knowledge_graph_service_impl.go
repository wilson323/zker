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
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
	"github.com/coze-dev/coze-studio/backend/pkg/redis"
)

const (
	// GraphEntityPrefix Redis key前缀
	GraphEntityPrefix = "knowledge_graph:entity:"
	// GraphRelationPrefix Redis key前缀
	GraphRelationPrefix = "knowledge_graph:relation:"
	// DefaultMaxHops 默认最大跳数
	DefaultMaxHops = 2
)

// NewKnowledgeGraphService 创建知识图谱服务
func NewKnowledgeGraphService(
	redisClient *redis.Client,
	llmClient llm.MemoryLLMClient,
) KnowledgeGraphService {
	return &knowledgeGraphServiceImpl{
		redisClient: redisClient,
		llmClient:   llmClient,
	}
}

type knowledgeGraphServiceImpl struct {
	redisClient *redis.Client
	llmClient   llm.MemoryLLMClient
}

// ExtractEntitiesAndRelations 从文本中抽取实体和关系
func (s *knowledgeGraphServiceImpl) ExtractEntitiesAndRelations(
	ctx context.Context,
	req *entity.ExtractEntitiesRequest,
) (*entity.ExtractEntitiesResponse, error) {
	// 1. 构建LLM提示
	prompt := s.buildExtractionPrompt(req.Text)

	// 2. 调用LLM抽取
	response, err := s.llmClient.Chat(ctx, &llm.ChatRequest{
		Messages: []*llm.Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 2000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to extract with LLM: %w", err)
	}

	// 3. 解析LLM响应
	result, err := s.parseExtractionResult(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse extraction result: %w", err)
	}

	// 4. 存储实体和关系到Redis
	for i := range result.Entities {
		if result.Entities[i].ID == "" {
			result.Entities[i].ID = uuid.New().String()
		}
		if result.Entities[i].CreatedAt.IsZero() {
			result.Entities[i].CreatedAt = time.Now()
		}
		err = s.AddEntity(ctx, &result.Entities[i])
		if err != nil {
			return nil, fmt.Errorf("failed to store entity: %w", err)
		}
	}

	for i := range result.Relations {
		if result.Relations[i].ID == "" {
			result.Relations[i].ID = uuid.New().String()
		}
		if result.Relations[i].CreatedAt.IsZero() {
			result.Relations[i].CreatedAt = time.Now()
		}
		err = s.AddRelation(ctx, &result.Relations[i])
		if err != nil {
			return nil, fmt.Errorf("failed to store relation: %w", err)
		}
	}

	return &entity.ExtractEntitiesResponse{
		Entities:   result.Entities,
		Relations:  result.Relations,
		Confidence: result.Confidence,
		ExtractedAt: time.Now(),
	}, nil
}

// Query 查询知识图谱
func (s *knowledgeGraphServiceImpl) Query(
	ctx context.Context,
	req *entity.QueryKnowledgeGraphRequest,
) (*entity.QueryKnowledgeGraphResponse, error) {
	startTime := time.Now()

	// 简化实现:从Redis查询相关实体和关系
	entities := make([]entity.GraphEntity, 0)
	relations := make([]entity.GraphRelation, 0)

	// TODO: 实际应该使用Neo4j Cypher查询
	// 这里简化为从Redis查询
	entityKeys, err := s.redisClient.Keys(ctx, GraphEntityPrefix+"*")
	if err == nil && len(entityKeys) > 0 {
		for _, key := range entityKeys {
			var entity entity.GraphEntity
			err := s.redisClient.Get(ctx, key, &entity)
			if err == nil {
				entities = append(entities, entity)
			}
		}
	}

	relationKeys, err := s.redisClient.Keys(ctx, GraphRelationPrefix+"*")
	if err == nil && len(relationKeys) > 0 {
		for _, key := range relationKeys {
			var relation entity.GraphRelation
			err := s.redisClient.Get(ctx, key, &relation)
			if err == nil {
				relations = append(relations, relation)
			}
		}
	}

	// 构建路径(简化)
	paths := make([]entity.GraphPath, 0)

	return &entity.QueryKnowledgeGraphResponse{
		Entities:  entities,
		Relations: relations,
		Paths:     paths,
		QueryTime: time.Since(startTime).Milliseconds(),
	}, nil
}

// Visualize 生成可视化数据
func (s *knowledgeGraphServiceImpl) Visualize(
	ctx context.Context,
	req *entity.VisualizeKnowledgeGraphRequest,
) (*entity.VisualizeKnowledgeGraphResponse, error) {
	// 1. 查询图谱
	queryReq := &entity.QueryKnowledgeGraphRequest{
		ConversationID: req.ConversationID,
		MaxHops:       DefaultMaxHops,
		Limit:         req.MaxNodes,
	}

	graph, err := s.Query(ctx, queryReq)
	if err != nil {
		return nil, fmt.Errorf("failed to query graph: %w", err)
	}

	// 2. 构建可视化节点
	nodes := make([]entity.GraphNode, len(graph.Entities))
	for i, entity := range graph.Entities {
		nodes[i] = entity.GraphNode{
			ID:         entity.ID,
			Label:      entity.Name,
			Type:       entity.Type,
			Properties: entity.Properties,
			Size:       s.calculateNodeSize(entity.Type),
			Color:      s.getNodeColor(entity.Type),
		}
	}

	// 3. 构建可视化边
	edges := make([]entity.GraphEdge, len(graph.Relations))
	for i, relation := range graph.Relations {
		edges[i] = entity.GraphEdge{
			ID:         relation.ID,
			From:       relation.FromEntity,
			To:         relation.ToEntity,
			Label:      relation.RelationType,
			Weight:     s.calculateEdgeWeight(relation.Confidence),
			Confidence: relation.Confidence,
		}
	}

	// 4. 计算统计信息
	stats := s.calculateStats(nodes, edges)

	return &entity.VisualizeKnowledgeGraphResponse{
		Nodes: nodes,
		Edges: edges,
		Stats: stats,
	}, nil
}

// AddEntity 添加实体
func (s *knowledgeGraphServiceImpl) AddEntity(
	ctx context.Context,
	entity *entity.GraphEntity,
) error {
	key := GraphEntityPrefix + entity.ID
	return s.redisClient.Set(ctx, key, entity, 0)
}

// AddRelation 添加关系
func (s *knowledgeGraphServiceImpl) AddRelation(
	ctx context.Context,
	relation *entity.GraphRelation,
) error {
	key := GraphRelationPrefix + relation.ID
	return s.redisClient.Set(ctx, key, relation, 0)
}

// DeleteEntity 删除实体
func (s *knowledgeGraphServiceImpl) DeleteEntity(
	ctx context.Context,
	entityID string,
) error {
	key := GraphEntityPrefix + entityID
	return s.redisClient.Del(ctx, key)
}

// DeleteRelation 删除关系
func (s *knowledgeGraphServiceImpl) DeleteRelation(
	ctx context.Context,
	relationID string,
) error {
	key := GraphRelationPrefix + relationID
	return s.redisClient.Del(ctx, key)
}

// GetEntity 获取实体
func (s *knowledgeGraphServiceImpl) GetEntity(
	ctx context.Context,
	entityID string,
) (*entity.GraphEntity, error) {
	key := GraphEntityPrefix + entityID
	var entity entity.GraphEntity
	err := s.redisClient.Get(ctx, key, &entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetEntityRelations 获取实体的所有关系
func (s *knowledgeGraphServiceImpl) GetEntityRelations(
	ctx context.Context,
	entityID string,
) ([]entity.GraphRelation, error) {
	// 查询所有关系
	relationKeys, err := s.redisClient.Keys(ctx, GraphRelationPrefix+"*")
	if err != nil {
		return nil, err
	}

	relations := make([]entity.GraphRelation, 0)
	for _, key := range relationKeys {
		var relation entity.GraphRelation
		err := s.redisClient.Get(ctx, key, &relation)
		if err != nil {
			continue
		}
		// 筛选与该实体相关的关系
		if relation.FromEntity == entityID || relation.ToEntity == entityID {
			relations = append(relations, relation)
		}
	}

	return relations, nil
}

// GetGraphStats 获取图谱统计
func (s *knowledgeGraphServiceImpl) GetGraphStats(
	ctx context.Context,
	conversationID string,
) (*entity.GraphStats, error) {
	// 查询所有实体和关系
	entityKeys, err := s.redisClient.Keys(ctx, GraphEntityPrefix+"*")
	if err != nil {
		return nil, err
	}

	relationKeys, err := s.redisClient.Keys(ctx, GraphRelationPrefix+"*")
	if err != nil {
		return nil, err
	}

	stats := &entity.GraphStats{
		TotalNodes: len(entityKeys),
		TotalEdges: len(relationKeys),
	}

	// 计算平均度数
	if stats.TotalNodes > 0 {
		stats.AvgDegree = float64(stats.TotalEdges*2) / float64(stats.TotalNodes)
	}

	// 计算密度
	if stats.TotalNodes > 1 {
		maxEdges := float64(stats.TotalNodes * (stats.TotalNodes - 1) / 2)
		stats.Density = float64(stats.TotalEdges) / maxEdges
	}

	return stats, nil
}

// buildExtractionPrompt 构建实体抽取提示
func (s *knowledgeGraphServiceImpl) buildExtractionPrompt(text string) string {
	return fmt.Sprintf(`请从以下文本中抽取实体和关系:

文本:%s

请以JSON格式返回,格式如下:
{
  "entities": [
    {"name": "实体名称", "type": "PERSON|ORG|LOCATION|CONCEPT", "properties": {}}
  ],
  "relations": [
    {"from_entity": "实体1", "to_entity": "实体2", "relation_type": "关系类型", "confidence": 0.9}
  ],
  "confidence": 0.85
}

注意:
1. 实体类型包括: PERSON(人物), ORG(组织), LOCATION(地点), CONCEPT(概念)
2. 关系类型包括: located_in(位于), founded_by(创立者), works_for(就职于), part_of(属于)等
3. confidence表示抽取置信度(0-1)
4. 只返回JSON,不要有任何额外说明
`, text)
}

// parseExtractionResult 解析抽取结果
func (s *knowledgeGraphServiceImpl) parseExtractionResult(response string) (*entity.ExtractEntitiesResponse, error) {
	var result struct {
		Entities   []entity.GraphEntity   `json:"entities"`
		Relations  []entity.GraphRelation `json:"relations"`
		Confidence float64                `json:"confidence"`
	}

	err := json.Unmarshal([]byte(response), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &entity.ExtractEntitiesResponse{
		Entities:   result.Entities,
		Relations:  result.Relations,
		Confidence: result.Confidence,
	}, nil
}

// calculateNodeSize 计算节点大小
func (s *knowledgeGraphServiceImpl) calculateNodeSize(entityType string) int {
	sizes := map[string]int{
		"PERSON":    30,
		"ORG":       40,
		"LOCATION":  35,
		"CONCEPT":   25,
	}
	if size, exists := sizes[entityType]; exists {
		return size
	}
	return 20
}

// getNodeColor 获取节点颜色
func (s *knowledgeGraphServiceImpl) getNodeColor(entityType string) string {
	colors := map[string]string{
		"PERSON":    "#FF6B6B",
		"ORG":       "#4ECDC4",
		"LOCATION":  "#45B7D1",
		"CONCEPT":   "#FFA07A",
	}
	if color, exists := colors[entityType]; exists {
		return color
	}
	return "#95E1D3"
}

// calculateEdgeWeight 计算边权重
func (s *knowledgeGraphServiceImpl) calculateEdgeWeight(confidence float64) float64 {
	return confidence * 5.0 // 权重范围0-5
}

// calculateStats 计算统计信息
func (s *knowledgeGraphServiceImpl) calculateStats(
	nodes []entity.GraphNode,
	edges []entity.GraphEdge,
) entity.GraphStats {
	stats := entity.GraphStats{
		TotalNodes: len(nodes),
		TotalEdges: len(edges),
	}

	// 计算平均度数
	if stats.TotalNodes > 0 {
		degrees := make(map[string]int)
		for _, edge := range edges {
			degrees[edge.From]++
			degrees[edge.To]++
		}

		totalDegree := 0
		maxDegree := 0
		for _, degree := range degrees {
			totalDegree += degree
			if degree > maxDegree {
				maxDegree = degree
			}
		}

		stats.AvgDegree = float64(totalDegree) / float64(stats.TotalNodes)
		stats.MaxDegree = maxDegree
	}

	// 计算密度
	if stats.TotalNodes > 1 {
		maxEdges := float64(stats.TotalNodes * (stats.TotalNodes - 1) / 2)
		stats.Density = float64(stats.TotalEdges) / maxEdges
	}

	return stats
}
