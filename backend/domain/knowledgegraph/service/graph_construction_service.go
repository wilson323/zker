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

	"github.com/bytedance/sonic"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/entity"
	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/repository"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/pkg/llmclient"
)

// LLMClient LLM客户端接口（用于实体和关系抽取）
type LLMClient interface {
	Chat(ctx context.Context, req *llmclient.ChatRequest) (*llmclient.ChatResponse, error)
}

// VectorStore 向量存储接口（用于存储实体嵌入）
type VectorStore interface {
	Insert(ctx context.Context, collection string, vectors []Vector) error
	Search(ctx context.Context, collection string, vector []float32, topK int) ([]SearchResult, error)
}

// Vector 向量数据
type Vector struct {
	ID      string
	Vector  []float32
	Metadata map[string]interface{}
}

// SearchResult 向量搜索结果
type SearchResult struct {
	ID       string
	Score    float32
	Metadata map[string]interface{}
}

// GraphConstructionService 图谱构建服务
// 职责：从文本中抽取实体和关系,构建知识图谱
type GraphConstructionService struct {
	entityRepo       repository.GraphEntityRepository
	relationshipRepo repository.GraphRelationshipRepository
	llmClient        LLMClient
	vectorStore      VectorStore
	idgen            idgen.IDGenerator
	logger           *zap.Logger
}

// GraphConstructionServiceConfig 配置
type GraphConstructionServiceConfig struct {
	EntityRepo       repository.GraphEntityRepository
	RelationshipRepo repository.GraphRelationshipRepository
	LLMClient        LLMClient
	VectorStore      VectorStore
	IDGen            idgen.IDGenerator
	Logger           *zap.Logger
}

// NewGraphConstructionService 创建图谱构建服务
func NewGraphConstructionService(config *GraphConstructionServiceConfig) *GraphConstructionService {
	return &GraphConstructionService{
		entityRepo:       config.EntityRepo,
		relationshipRepo: config.RelationshipRepo,
		llmClient:        config.LLMClient,
		vectorStore:      config.VectorStore,
		idgen:            config.IDGen,
		logger:           config.Logger,
	}
}

// ExtractEntities 从文本中抽取实体
func (s *GraphConstructionService) ExtractEntities(ctx context.Context, tenantID, text string) ([]*entity.GraphEntity, error) {
	s.logger.Info("开始抽取实体",
		zap.String("tenant_id", tenantID),
		zap.Int("text_length", len(text)))

	// 构建LLM提示词
	prompt := s.buildEntityExtractionPrompt(text)

	// 调用LLM
	req := &llmclient.ChatRequest{
		Messages: []*llmclient.Message{
			{Role: "system", Content: "你是一个专业的知识图谱构建助手,擅长从文本中识别和抽取实体。"},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
	}

	resp, err := s.llmClient.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM调用失败: %w", err)
	}

	// 解析LLM响应
	entities, err := s.parseEntityExtractionResponse(resp.Content, tenantID)
	if err != nil {
		return nil, fmt.Errorf("解析实体抽取响应失败: %w", err)
	}

	s.logger.Info("实体抽取完成",
		zap.String("tenant_id", tenantID),
		zap.Int("entity_count", len(entities)))

	return entities, nil
}

// ExtractRelationships 从文本中抽取关系
func (s *GraphConstructionService) ExtractRelationships(ctx context.Context, tenantID string, text string, entities []*entity.GraphEntity) ([]*entity.GraphRelationship, error) {
	s.logger.Info("开始抽取关系",
		zap.String("tenant_id", tenantID),
		zap.Int("entity_count", len(entities)))

	// 构建实体映射（用于LLM引用）
	entityMap := make(map[string]string)
	for _, ent := range entities {
		entityMap[ent.EntityName] = ent.ID
	}

	// 构建LLM提示词
	prompt := s.buildRelationshipExtractionPrompt(text, entities)

	// 调用LLM
	req := &llmclient.ChatRequest{
		Messages: []*llmclient.Message{
			{Role: "system", Content: "你是一个专业的知识图谱构建助手,擅长从文本中识别和抽取实体间的关系。"},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.3,
	}

	resp, err := s.llmClient.Chat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("LLM调用失败: %w", err)
	}

	// 解析LLM响应
	relationships, err := s.parseRelationshipExtractionResponse(resp.Content, tenantID, entityMap)
	if err != nil {
		return nil, fmt.Errorf("解析关系抽取响应失败: %w", err)
	}

	s.logger.Info("关系抽取完成",
		zap.String("tenant_id", tenantID),
		zap.Int("relationship_count", len(relationships)))

	return relationships, nil
}

// BuildGraphFromText 从文本构建完整图谱
func (s *GraphConstructionService) BuildGraphFromText(ctx context.Context, tenantID, text string) error {
	s.logger.Info("开始从文本构建图谱",
		zap.String("tenant_id", tenantID),
		zap.Int("text_length", len(text)))

	// 1. 抽取实体
	entities, err := s.ExtractEntities(ctx, tenantID, text)
	if err != nil {
		return fmt.Errorf("实体抽取失败: %w", err)
	}

	// 2. 批量保存实体
	if err := s.AddEntities(ctx, entities); err != nil {
		return fmt.Errorf("保存实体失败: %w", err)
	}

	// 3. 抽取关系
	relationships, err := s.ExtractRelationships(ctx, tenantID, text, entities)
	if err != nil {
		return fmt.Errorf("关系抽取失败: %w", err)
	}

	// 4. 批量保存关系
	if err := s.AddRelationships(ctx, relationships); err != nil {
		return fmt.Errorf("保存关系失败: %w", err)
	}

	s.logger.Info("图谱构建完成",
		zap.String("tenant_id", tenantID),
		zap.Int("entity_count", len(entities)),
		zap.Int("relationship_count", len(relationships)))

	return nil
}

// AddEntity 添加单个实体
func (s *GraphConstructionService) AddEntity(ctx context.Context, ent *entity.GraphEntity) error {
	if ent.ID == "" {
		ent.ID = s.idgen.GenerateID()
	}

	if err := s.entityRepo.Create(ctx, ent); err != nil {
		return fmt.Errorf("创建实体失败: %w", err)
	}

	s.logger.Debug("实体添加成功",
		zap.String("entity_id", ent.ID),
		zap.String("entity_name", ent.EntityName),
		zap.String("entity_type", string(ent.EntityType)))

	return nil
}

// AddEntities 批量添加实体
func (s *GraphConstructionService) AddEntities(ctx context.Context, entities []*entity.GraphEntity) error {
	if len(entities) == 0 {
		return nil
	}

	// 为没有ID的实体生成ID
	for _, ent := range entities {
		if ent.ID == "" {
			ent.ID = s.idgen.GenerateID()
		}
	}

	if err := s.entityRepo.CreateBatch(ctx, entities); err != nil {
		return fmt.Errorf("批量创建实体失败: %w", err)
	}

	s.logger.Debug("批量添加实体成功", zap.Int("count", len(entities)))
	return nil
}

// AddRelationship 添加单个关系
func (s *GraphConstructionService) AddRelationship(ctx context.Context, rel *entity.GraphRelationship) error {
	if rel.ID == "" {
		rel.ID = s.idgen.GenerateID()
	}

	// 验证关系有效性
	if !rel.IsValid() {
		return fmt.Errorf("无效的关系: source=%s, target=%s, type=%s",
			rel.SourceEntityID, rel.TargetEntityID, rel.RelationType)
	}

	if err := s.relationshipRepo.Create(ctx, rel); err != nil {
		return fmt.Errorf("创建关系失败: %w", err)
	}

	s.logger.Debug("关系添加成功",
		zap.String("relationship_id", rel.ID),
		zap.String("source_id", rel.SourceEntityID),
		zap.String("target_id", rel.TargetEntityID),
		zap.String("relation_type", string(rel.RelationType)))

	return nil
}

// AddRelationships 批量添加关系
func (s *GraphConstructionService) AddRelationships(ctx context.Context, relationships []*entity.GraphRelationship) error {
	if len(relationships) == 0 {
		return nil
	}

	// 为没有ID的关系生成ID
	for _, rel := range relationships {
		if rel.ID == "" {
			rel.ID = s.idgen.GenerateID()
		}
	}

	if err := s.relationshipRepo.CreateBatch(ctx, relationships); err != nil {
		return fmt.Errorf("批量创建关系失败: %w", err)
	}

	s.logger.Debug("批量添加关系成功", zap.Int("count", len(relationships)))
	return nil
}

// ========================================
// 内部辅助方法
// ========================================

// buildEntityExtractionPrompt 构建实体抽取提示词
func (s *GraphConstructionService) buildEntityExtractionPrompt(text string) string {
	return fmt.Sprintf(`请从以下文本中抽取所有重要的实体。实体类型包括：
- person: 人物
- organization: 组织/机构
- location: 地点
- concept: 概念
- event: 事件
- product: 产品
- document: 文档

文本内容：
%s

请以JSON格式返回结果，格式如下：
[
  {
    "entity_name": "实体名称",
    "entity_type": "实体类型",
    "properties": {
      "description": "简短描述",
      "其他属性": "属性值"
    }
  }
]

只返回JSON，不要其他内容。`, text)
}

// buildRelationshipExtractionPrompt 构建关系抽取提示词
func (s *GraphConstructionService) buildRelationshipExtractionPrompt(text string, entities []*entity.GraphEntity) string {
	// 构建实体列表
	entityList := ""
	for i, ent := range entities {
		entityList += fmt.Sprintf("%d. %s (ID: %s, 类型: %s)\n", i+1, ent.EntityName, ent.ID, ent.EntityType)
	}

	return fmt.Sprintf(`请从以下文本中抽取实体之间的关系。

已知实体列表：
%s

文本内容：
%s

关系类型包括：
- knows: 认识
- works_for: 任职于
- located_in: 位于
- part_of: 属于
- related_to: 相关
- cites: 引用
- about: 关于
等...

请以JSON格式返回结果，格式如下：
[
  {
    "source_entity": "源实体名称",
    "target_entity": "目标实体名称",
    "relation_type": "关系类型",
    "properties": {
      "weight": 0.8,
      "其他属性": "属性值"
    }
  }
]

只返回JSON，不要其他内容。`, entityList, text)
}

// parseEntityExtractionResponse 解析实体抽取响应
func (s *GraphConstructionService) parseEntityExtractionResponse(content string, tenantID string) ([]*entity.GraphEntity, error) {
	var extractedEntities []struct {
		EntityName  string                 `json:"entity_name"`
		EntityType  string                 `json:"entity_type"`
		Properties  map[string]interface{} `json:"properties"`
	}

	if err := sonic.Unmarshal([]byte(content), &extractedEntities); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	entities := make([]*entity.GraphEntity, 0, len(extractedEntities))
	for _, extracted := range extractedEntities {
		// 验证实体类型
		entType := entity.EntityType(extracted.EntityType)
		if !entType.IsValid() {
			s.logger.Warn("跳过无效实体类型",
				zap.String("entity_name", extracted.EntityName),
				zap.String("entity_type", extracted.EntityType))
			continue
		}

		ent := &entity.GraphEntity{
			ID:         s.idgen.GenerateID(),
			TenantID:   tenantID,
			EntityType: entType,
			EntityName: extracted.EntityName,
			Properties: extracted.Properties,
		}

		entities = append(entities, ent)
	}

	return entities, nil
}

// parseRelationshipExtractionResponse 解析关系抽取响应
func (s *GraphConstructionService) parseRelationshipExtractionResponse(content string, tenantID string, entityMap map[string]string) ([]*entity.GraphRelationship, error) {
	var extractedRelationships []struct {
		SourceEntity string                 `json:"source_entity"`
		TargetEntity string                 `json:"target_entity"`
		RelationType string                 `json:"relation_type"`
		Properties   map[string]interface{} `json:"properties"`
	}

	if err := sonic.Unmarshal([]byte(content), &extractedRelationships); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	relationships := make([]*entity.GraphRelationship, 0, len(extractedRelationships))
	for _, extracted := range extractedRelationships {
		// 查找实体ID
		sourceID, ok := entityMap[extracted.SourceEntity]
		if !ok {
			s.logger.Warn("源实体不存在", zap.String("source_entity", extracted.SourceEntity))
			continue
		}

		targetID, ok := entityMap[extracted.TargetEntity]
		if !ok {
			s.logger.Warn("目标实体不存在", zap.String("target_entity", extracted.TargetEntity))
			continue
		}

		relType := entity.RelationType(extracted.RelationType)

		// 从属性中获取权重
		weight := float32(1.0)
		if weightVal, ok := extracted.Properties["weight"]; ok {
			if w, ok := weightVal.(float64); ok {
				weight = float32(w)
			}
		}

		rel := &entity.GraphRelationship{
			ID:             s.idgen.GenerateID(),
			TenantID:       tenantID,
			SourceEntityID: sourceID,
			TargetEntityID: targetID,
			RelationType:   relType,
			Properties:     extracted.Properties,
			Weight:         weight,
		}

		// 验证关系有效性
		if !rel.IsValid() {
			s.logger.Warn("跳过无效关系",
				zap.String("source_id", sourceID),
				zap.String("target_id", targetID),
				zap.String("relation_type", string(relType)))
			continue
		}

		relationships = append(relationships, rel)
	}

	return relationships, nil
}

// GenerateEntityEmbedding 为实体生成向量嵌入
func (s *GraphConstructionService) GenerateEntityEmbedding(ctx context.Context, ent *entity.GraphEntity) error {
	// 调用嵌入模型生成向量
	// 这里需要实现具体的嵌入逻辑
	return nil
}

// ========================================
// 数据结构
// ========================================

// EntityExtractionResult 实体抽取结果
type EntityExtractionResult struct {
	Entities []*entity.GraphEntity `json:"entities"`
}

// RelationshipExtractionResult 关系抽取结果
type RelationshipExtractionResult struct {
	Relationships []*entity.GraphRelationship `json:"relationships"`
}

// GraphBuildResult 图谱构建结果
type GraphBuildResult struct {
	EntityCount       int `json:"entity_count"`
	RelationshipCount int `json:"relationship_count"`
}

// ExtractEntityAndRelationship 从文本抽取实体和关系（一步完成）
func (s *GraphConstructionService) ExtractEntityAndRelationship(ctx context.Context, tenantID, text string) ([]*entity.GraphEntity, []*entity.GraphRelationship, error) {
	s.logger.Info("开始抽取实体和关系",
		zap.String("tenant_id", tenantID),
		zap.Int("text_length", len(text)))

	// 1. 抽取实体
	entities, err := s.ExtractEntities(ctx, tenantID, text)
	if err != nil {
		return nil, nil, fmt.Errorf("实体抽取失败: %w", err)
	}

	// 2. 抽取关系
	relationships, err := s.ExtractRelationships(ctx, tenantID, text, entities)
	if err != nil {
		return nil, nil, fmt.Errorf("关系抽取失败: %w", err)
	}

	return entities, relationships, nil
}

// MergeEntity 合并实体（去重）
func (s *GraphConstructionService) MergeEntity(ctx context.Context, tenantID string, ent *entity.GraphEntity) (*entity.GraphEntity, error) {
	// 检查是否已存在同名实体
	existing, err := s.entityRepo.GetByTenantIDAndName(ctx, tenantID, ent.EntityName)
	if err == nil {
		// 实体已存在，合并属性
		for k, v := range ent.Properties {
			existing.SetProperty(k, v)
		}

		if err := s.entityRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("更新实体失败: %w", err)
		}

		return existing, nil
	}

	// 新实体，创建
	if err := s.AddEntity(ctx, ent); err != nil {
		return nil, fmt.Errorf("创建实体失败: %w", err)
	}

	return ent, nil
}

// UpsertEntity 创建或更新实体
func (s *GraphConstructionService) UpsertEntity(ctx context.Context, ent *entity.GraphEntity) error {
	return s.entityRepo.Upsert(ctx, ent)
}

// DeleteEntity 删除实体及其所有关系
func (s *GraphConstructionService) DeleteEntity(ctx context.Context, entityID string) error {
	// 1. 删除所有相关关系
	if err := s.relationshipRepo.DeleteByEntity(ctx, entityID); err != nil {
		return fmt.Errorf("删除关系失败: %w", err)
	}

	// 2. 删除实体
	if err := s.entityRepo.Delete(ctx, entityID); err != nil {
		return fmt.Errorf("删除实体失败: %w", err)
	}

	s.logger.Info("实体删除成功", zap.String("entity_id", entityID))
	return nil
}

// GetEntityStats 获取实体统计信息
func (s *GraphConstructionService) GetEntityStats(ctx context.Context, tenantID string) (map[entity.EntityType]int64, error) {
	stats := make(map[entity.EntityType]int64)

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
		if err != nil {
			return nil, fmt.Errorf("统计实体数量失败: %w", err)
		}
		stats[entityType] = count
	}

	return stats, nil
}
