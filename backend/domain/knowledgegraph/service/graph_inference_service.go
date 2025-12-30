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

// GraphInferenceService 图谱推理服务
// 职责：基于图谱进行推理,推断潜在关系、实体类型等
type GraphInferenceService struct {
	entityRepo       repository.GraphEntityRepository
	relationshipRepo repository.GraphRelationshipRepository
	llmClient        LLMClient
	logger           *zap.Logger
}

// GraphInferenceServiceConfig 配置
type GraphInferenceServiceConfig struct {
	EntityRepo       repository.GraphEntityRepository
	RelationshipRepo repository.GraphRelationshipRepository
	LLMClient        LLMClient
	Logger           *zap.Logger
}

// NewGraphInferenceService 创建图谱推理服务
func NewGraphInferenceService(config *GraphInferenceServiceConfig) *GraphInferenceService {
	return &GraphInferenceService{
		entityRepo:       config.EntityRepo,
		relationshipRepo: config.RelationshipRepo,
		llmClient:        config.LLMClient,
		logger:           config.Logger,
	}
}

// InferRelationship 推断两个实体之间的关系
func (s *GraphInferenceService) InferRelationship(ctx context.Context, entity1ID, entity2ID string) (*entity.GraphRelationship, error) {
	// 1. 检查关系是否已存在
	existingRel, err := s.relationshipRepo.GetByEntities(ctx, entity1ID, entity2ID, "")
	if err == nil {
		return existingRel, nil
	}

	// 2. 获取实体信息
	entity1, err := s.entityRepo.GetByID(ctx, entity1ID)
	if err != nil {
		return nil, fmt.Errorf("获取实体1失败: %w", err)
	}

	entity2, err := s.entityRepo.GetByID(ctx, entity2ID)
	if err != nil {
		return nil, fmt.Errorf("获取实体2失败: %w", err)
	}

	// 3. 使用LLM推断关系
	relType, confidence, err := s.inferRelationTypeWithLLM(ctx, entity1, entity2)
	if err != nil {
		s.logger.Warn("LLM推断关系失败", zap.Error(err))
		// 降级：使用规则推断
		relType, confidence = s.inferRelationTypeByRule(entity1, entity2)
	}

	if confidence < 0.5 {
		return nil, fmt.Errorf("无法确定关系类型 (confidence=%.2f)", confidence)
	}

	// 4. 创建推断的关系
	rel := &entity.GraphRelationship{
		TenantID:       entity1.TenantID,
		SourceEntityID: entity1ID,
		TargetEntityID: entity2ID,
		RelationType:   relType,
		Weight:         confidence,
		Properties:     make(entity.RelationshipProperties),
	}

	rel.SetProperty("inferred", true)
	rel.SetProperty("confidence", confidence)

	s.logger.Info("关系推断完成",
		zap.String("entity1", entity1.EntityName),
		zap.String("entity2", entity2.EntityName),
		zap.String("relation_type", string(relType)),
		zap.Float32("confidence", confidence))

	return rel, nil
}

// InferEntityType 推断实体类型
func (s *GraphInferenceService) InferEntityType(ctx context.Context, tenantID, text string) (entity.EntityType, error) {
	// 使用LLM推断实体类型
	prompt := fmt.Sprintf(`请根据以下文本推断实体的类型。

文本：%s

可选类型：
- person: 人物
- organization: 组织/机构
- location: 地点
- concept: 概念
- event: 事件
- product: 产品
- document: 文档

请只返回类型名称，不要其他内容。`, text)

	// 这里简化实现，实际应该调用LLM
	// 使用简单的规则推断
	_ = prompt // TODO: 实现LLM调用时使用prompt
	return s.inferEntityTypeByRule(text), nil
}

// CompleteSubgraph 补全子图（推断缺失的关系）
func (s *GraphInferenceService) CompleteSubgraph(ctx context.Context, entityID string) (*entity.GraphSubgraph, error) {
	// 1. 获取实体
	centerEntity, err := s.entityRepo.GetByID(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("获取中心实体失败: %w", err)
	}

	// 2. 获取现有邻居
	outbound, inbound, err := s.relationshipRepo.ListNeighbors(ctx, entityID, 100)
	if err != nil {
		return nil, fmt.Errorf("获取邻居失败: %w", err)
	}

	// 3. 转换为GraphSubgraph
	centerNode := &entity.GraphNode{
		ID:         centerEntity.ID,
		EntityType: centerEntity.EntityType,
		EntityName: centerEntity.EntityName,
		Properties: centerEntity.Properties,
	}

	subgraph := &entity.GraphSubgraph{
		CenterNode:    centerNode,
		Neighbors:     make([]*entity.GraphNode, 0),
		Relationships: make([]*entity.GraphRelation, 0),
		Depth:         1,
	}

	// 添加现有关系
	for _, rel := range outbound {
		subgraph.Relationships = append(subgraph.Relationships, &entity.GraphRelation{
			ID:           rel.ID,
			SourceID:     rel.SourceEntityID,
			TargetID:     rel.TargetEntityID,
			RelationType: rel.RelationType,
			Properties:   rel.Properties,
		})

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

	for _, rel := range inbound {
		subgraph.Relationships = append(subgraph.Relationships, &entity.GraphRelation{
			ID:           rel.ID,
			SourceID:     rel.SourceEntityID,
			TargetID:     rel.TargetEntityID,
			RelationType: rel.RelationType,
			Properties:   rel.Properties,
		})

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

	// 4. 推断缺失的关系（简化实现）
	// 实际应该基于图嵌入或规则引擎

	s.logger.Info("子图补全完成",
		zap.String("entity_id", entityID),
		zap.Int("neighbor_count", len(subgraph.Neighbors)),
		zap.Int("relationship_count", len(subgraph.Relationships)))

	return subgraph, nil
}

// RecommendEntities 推荐相关实体
func (s *GraphInferenceService) RecommendEntities(ctx context.Context, entityID string, topK int) ([]*entity.GraphEntity, error) {
	// 1. 获取实体的所有邻居
	outbound, inbound, err := s.relationshipRepo.ListNeighbors(ctx, entityID, 100)
	if err != nil {
		return nil, fmt.Errorf("获取邻居失败: %w", err)
	}

	// 2. 统计邻居的邻居（2跳节点）
	neighborScore := make(map[string]float32)

	// 处理出边
	for _, rel := range outbound {
		neighborID := rel.TargetEntityID
		weight := rel.GetWeight()

		// 获取邻居的邻居
		nextOutbound, _, _ := s.relationshipRepo.ListNeighbors(ctx, neighborID, 50)

		for _, nextRel := range nextOutbound {
			candidateID := nextRel.TargetEntityID
			if candidateID == entityID {
				continue
			}

			// 累加权重
			neighborScore[candidateID] += weight * nextRel.GetWeight() * 0.5
		}
	}

	// 处理入边
	for _, rel := range inbound {
		neighborID := rel.SourceEntityID
		weight := rel.GetWeight()

		_, nextInbound, _ := s.relationshipRepo.ListNeighbors(ctx, neighborID, 50)

		for _, nextRel := range nextInbound {
			candidateID := nextRel.SourceEntityID
			if candidateID == entityID {
				continue
			}

			neighborScore[candidateID] += weight * nextRel.GetWeight() * 0.5
		}
	}

	// 3. 排序并返回TopK
	type EntityScore struct {
		EntityID string
		Score    float32
	}

	scores := make([]EntityScore, 0, len(neighborScore))
	for id, score := range neighborScore {
		scores = append(scores, EntityScore{EntityID: id, Score: score})
	}

	// 简单排序（冒泡）
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].Score > scores[i].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// 获取TopK实体
	result := make([]*entity.GraphEntity, 0, topK)
	for i := 0; i < topK && i < len(scores); i++ {
		ent, err := s.entityRepo.GetByID(ctx, scores[i].EntityID)
		if err == nil {
			result = append(result, ent)
		}
	}

	s.logger.Info("实体推荐完成",
		zap.String("entity_id", entityID),
		zap.Int("recommended_count", len(result)))

	return result, nil
}

// FindSimilarEntities 查找相似实体
func (s *GraphInferenceService) FindSimilarEntities(ctx context.Context, entityID string, topK int) ([]*entity.GraphEntity, error) {
	// 1. 获取目标实体
	targetEntity, err := s.entityRepo.GetByID(ctx, entityID)
	if err != nil {
		return nil, fmt.Errorf("获取目标实体失败: %w", err)
	}

	// 2. 获取同类型的所有实体
	entities, _, err := s.entityRepo.ListByTenantID(ctx, targetEntity.TenantID, targetEntity.EntityType, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("获取同类实体失败: %w", err)
	}

	// 3. 计算相似度（基于Jaccard相似度）
	type EntitySimilarity struct {
		Entity     *entity.GraphEntity
		Similarity float32
	}

	similarities := make([]EntitySimilarity, 0)

	for _, ent := range entities {
		if ent.ID == entityID {
			continue
		}

		similarity := s.calculateEntitySimilarity(targetEntity, ent)
		similarities = append(similarities, EntitySimilarity{
			Entity:     ent,
			Similarity: similarity,
		})
	}

	// 4. 排序并返回TopK
	for i := 0; i < len(similarities); i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[j].Similarity > similarities[i].Similarity {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}

	result := make([]*entity.GraphEntity, 0, topK)
	for i := 0; i < topK && i < len(similarities); i++ {
		if similarities[i].Similarity > 0.3 { // 相似度阈值
			result = append(result, similarities[i].Entity)
		}
	}

	s.logger.Info("相似实体查找完成",
		zap.String("entity_id", entityID),
		zap.Int("similar_count", len(result)))

	return result, nil
}

// ========================================
// 内部辅助方法
// ========================================

// inferRelationTypeWithLLM 使用LLM推断关系类型
func (s *GraphInferenceService) inferRelationTypeWithLLM(ctx context.Context, entity1, entity2 *entity.GraphEntity) (entity.RelationType, float32, error) {
	// 简化实现：实际应该调用LLM
	// 这里使用规则作为降级方案
	relType, confidence := s.inferRelationTypeByRule(entity1, entity2)
	return relType, confidence, nil
}

// inferRelationTypeByRule 使用规则推断关系类型
func (s *GraphInferenceService) inferRelationTypeByRule(entity1, entity2 *entity.GraphEntity) (entity.RelationType, float32) {
	// 简单的规则引擎
	type1 := entity1.EntityType
	type2 := entity2.EntityType

	// Person -> Organization: works_for
	if type1 == entity.EntityTypePerson && type2 == entity.EntityTypeOrganization {
		return entity.RelationTypeWorksFor, 0.7
	}

	// Organization -> Person: part_of
	if type1 == entity.EntityTypeOrganization && type2 == entity.EntityTypePerson {
		return entity.RelationTypePartOf, 0.7
	}

	// Location -> Location: located_in
	if type1 == entity.EntityTypeLocation && type2 == entity.EntityTypeLocation {
		return entity.RelationTypeLocatedIn, 0.6
	}

	// 默认：related_to
	return entity.RelationTypeRelatedTo, 0.5
}

// inferEntityTypeByRule 使用规则推断实体类型
func (s *GraphInferenceService) inferEntityTypeByRule(text string) entity.EntityType {
	// 简化的规则引擎
	// 实际应该使用NLP或LLM

	keywords := map[string]entity.EntityType{
		"公司":  entity.EntityTypeOrganization,
		"企业":  entity.EntityTypeOrganization,
		"大学":  entity.EntityTypeOrganization,
		"医院":  entity.EntityTypeOrganization,
		"先生":  entity.EntityTypePerson,
		"女士":  entity.EntityTypePerson,
		"教授":  entity.EntityTypePerson,
		"博士":  entity.EntityTypePerson,
		"北京":  entity.EntityTypeLocation,
		"上海":  entity.EntityTypeLocation,
		"城市":  entity.EntityTypeLocation,
		"国家":  entity.EntityTypeLocation,
	}

	// 简单匹配关键词
	for keyword, entityType := range keywords {
		if contains(text, keyword) {
			return entityType
		}
	}

	// 默认：概念
	return entity.EntityTypeConcept
}

// calculateEntitySimilarity 计算实体相似度（Jaccard相似度）
func (s *GraphInferenceService) calculateEntitySimilarity(ent1, ent2 *entity.GraphEntity) float32 {
	// 1. 类型相似度
	typeSimilarity := float32(0.0)
	if ent1.EntityType == ent2.EntityType {
		typeSimilarity = 0.3
	}

	// 2. 属性相似度
	props1 := ent1.Properties
	props2 := ent2.Properties

	if len(props1) == 0 && len(props2) == 0 {
		return typeSimilarity
	}

	intersection := 0
	union := len(props1)

	for key := range props1 {
		if _, exists := props2[key]; exists {
			intersection++
		}
	}

	for key := range props2 {
		if _, exists := props1[key]; !exists {
			union++
		}
	}

	if union == 0 {
		return typeSimilarity
	}

	jaccard := float32(intersection) / float32(union)

	return typeSimilarity + jaccard*0.7
}

// contains 检查字符串是否包含子串（简单实现）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
