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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
	"github.com/coze-dev/coze-studio/backend/infra/llm"
)

// MockRedisClientForGraph Mock Redis客户端(用于知识图谱)
type MockRedisClientForGraph struct {
	mock.Mock
}

func (m *MockRedisClientForGraph) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisClientForGraph) Get(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockRedisClientForGraph) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockRedisClientForGraph) Keys(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), nil
}

// MockLLMClientForGraph Mock LLM客户端(用于知识图谱)
type MockLLMClientForGraph struct {
	mock.Mock
}

func (m *MockLLMClientForGraph) Chat(ctx context.Context, req *llm.ChatRequest) (string, error) {
	args := m.Called(ctx, req)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return args.String(0), nil
}

func (m *MockLLMClientForGraph) Embed(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), nil
}

// TestKnowledgeGraphService_ExtractEntitiesAndRelations 测试实体和关系抽取
func TestKnowledgeGraphService_ExtractEntitiesAndRelations(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	// 设置Mock期望
	mockLLM.On("Chat", mock.Anything, mock.Anything).Return(`{
		"entities": [
			{"name": "苹果公司", "type": "ORG", "properties": {"industry": "科技"}},
			{"name": "加利福尼亚州", "type": "LOCATION", "properties": {}}
		],
		"relations": [
			{"from_entity": "苹果公司", "to_entity": "加利福尼亚州", "relation_type": "located_in", "confidence": 0.95}
		],
		"confidence": 0.9
	}`, nil)

	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	req := &entity.ExtractEntitiesRequest{
		Text:           "苹果公司是一家位于加利福尼亚州的科技公司",
		ConversationID: "conv-001",
		TenantID:       "tenant-001",
		UserID:         "user-001",
	}

	resp, err := service.ExtractEntitiesAndRelations(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Entities))
	assert.Equal(t, 1, len(resp.Relations))
	assert.Greater(t, resp.Confidence, 0.8)
	mockLLM.AssertExpectations(t)
}

// TestKnowledgeGraphService_AddEntity 测试添加实体
func TestKnowledgeGraphService_AddEntity(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	entity := &entity.GraphEntity{
		ID:   "entity-001",
		Name: "测试实体",
		Type: "CONCEPT",
		Properties: map[string]interface{}{
			"key": "value",
		},
		CreatedAt: time.Now(),
	}

	err := service.AddEntity(context.Background(), entity)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

// TestKnowledgeGraphService_AddRelation 测试添加关系
func TestKnowledgeGraphService_AddRelation(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	relation := &entity.GraphRelation{
		ID:          "relation-001",
		FromEntity:  "entity-001",
		ToEntity:    "entity-002",
		RelationType: "related_to",
		Confidence:  0.9,
		CreatedAt:   time.Now(),
	}

	err := service.AddRelation(context.Background(), relation)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

// TestKnowledgeGraphService_GetEntity 测试获取实体
func TestKnowledgeGraphService_GetEntity(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	testEntity := &entity.GraphEntity{
		ID:        "entity-001",
		Name:      "测试实体",
		Type:      "CONCEPT",
		CreatedAt: time.Now(),
	}

	mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		value := args.Get(2).(*entity.GraphEntity)
		*value = *testEntity
	}).Return(nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	entity, err := service.GetEntity(context.Background(), "entity-001")

	assert.NoError(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, "entity-001", entity.ID)
	assert.Equal(t, "测试实体", entity.Name)
	mockRedis.AssertExpectations(t)
}

// TestKnowledgeGraphService_DeleteEntity 测试删除实体
func TestKnowledgeGraphService_DeleteEntity(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	mockRedis.On("Del", mock.Anything, mock.Anything).Return(nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	err := service.DeleteEntity(context.Background(), "entity-001")

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

// TestKnowledgeGraphService_Visualize 测试可视化
func TestKnowledgeGraphService_Visualize(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	// Mock Keys调用
	mockRedis.On("Keys", mock.Anything, mock.Anything).Return([]string{
		"knowledge_graph:entity:entity-001",
		"knowledge_graph:entity:entity-002",
		"knowledge_graph:relation:relation-001",
	}, nil)

	// Mock Get调用
	mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		key := args.String(1)
		value := args.Get(2)

		if key == "knowledge_graph:entity:entity-001" {
			entity := value.(*entity.GraphEntity)
			entity.ID = "entity-001"
			entity.Name = "实体1"
			entity.Type = "PERSON"
		} else if key == "knowledge_graph:entity:entity-002" {
			entity := value.(*entity.GraphEntity)
			entity.ID = "entity-002"
			entity.Name = "实体2"
			entity.Type = "ORG"
		} else if key == "knowledge_graph:relation:relation-001" {
			relation := value.(*entity.GraphRelation)
			relation.ID = "relation-001"
			relation.FromEntity = "entity-001"
			relation.ToEntity = "entity-002"
			relation.RelationType = "works_for"
			relation.Confidence = 0.9
		}
	}).Return(nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	req := &entity.VisualizeKnowledgeGraphRequest{
		ConversationID: "conv-001",
		MaxNodes:       50,
		MinConfidence:  0.5,
	}

	resp, err := service.Visualize(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Nodes))
	assert.Equal(t, 1, len(resp.Edges))
	assert.Equal(t, 2, resp.Stats.TotalNodes)
	assert.Equal(t, 1, resp.Stats.TotalEdges)
	mockRedis.AssertExpectations(t)
}

// TestKnowledgeGraphService_GetGraphStats 测试获取图谱统计
func TestKnowledgeGraphService_GetGraphStats(t *testing.T) {
	mockRedis := new(MockRedisClientForGraph)
	mockLLM := new(MockLLMClientForGraph)

	mockRedis.On("Keys", mock.Anything, "knowledge_graph:entity:*").Return([]string{
		"knowledge_graph:entity:entity-001",
		"knowledge_graph:entity:entity-002",
	}, nil)
	mockRedis.On("Keys", mock.Anything, "knowledge_graph:relation:*").Return([]string{
		"knowledge_graph:relation:relation-001",
	}, nil)

	service := NewKnowledgeGraphService(mockRedis, mockLLM)

	stats, err := service.GetGraphStats(context.Background(), "conv-001")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, 2, stats.TotalNodes)
	assert.Equal(t, 1, stats.TotalEdges)
	assert.Greater(t, stats.AvgDegree, 0.0)
	mockRedis.AssertExpectations(t)
}
