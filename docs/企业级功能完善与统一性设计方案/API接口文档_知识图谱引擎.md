# API接口文档：知识图谱引擎模块

**模块名称**: 知识图谱引擎 (KnowledgeGraphEngine)
**设计文档**: 18-五大AI引擎核心_知识图谱引擎.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 实体管理API](#2-实体管理api)
- [3. 关系管理API](#3-关系管理api)
- [4. 图谱构建API](#4-图谱构建api)
- [5. 图谱查询API](#5-图谱查询api)
- [6. 图谱推理API](#6-图谱推理api)
- [7. 数据模型](#7-数据模型)
- [8. 后端实现示例](#8-后端实现示例)
- [9. 前端调用示例](#9-前端调用示例)
- [10. 错误码定义](#10-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

知识图谱引擎负责构建和管理企业知识图谱，通过实体识别、关系抽取和图数据库，实现知识的结构化存储和推理查询：

- ✅ **实体识别** - 自动从文本中识别人名、地名、组织等实体
- ✅ **关系抽取** - 识别实体间的语义关系
- ✅ **知识融合** - 合并重复实体，维护一致性
- ✅ **图谱遍历** - 查找关联实体和关系
- ✅ **路径查询** - 查找实体间的关系路径
- ✅ **图谱推理** - 基于规则的推理查询

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 图数据库: Neo4j 5.x / NebulaGraph 3.x
- 辅助数据库: MySQL 8.4.5
- LLM: GPT-4 / 通义千问 (实体识别 + 关系抽取)

### 1.3 实体类型

| 实体类型 | 说明 | 示例属性 |
|---------|------|---------|
| **Person** | 人物 | name, title, department |
| **Organization** | 组织 | name, industry, location |
| **Location** | 地点 | name, country, region |
| **Product** | 产品 | name, category, version |
| **Event** | 事件 | name, date, participants |

---

## 2. 实体管理API

### 2.1 创建实体

**接口地址**: `POST /api/v1/knowledge-graph/entities`

**请求参数**:
```json
{
  "entity_type": "Person",
  "name": "张三",
  "properties": {
    "title": "产品经理",
    "department": "技术部",
    "email": "zhangsan@example.com"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "实体创建成功",
  "data": {
    "id": "entity-001",
    "type": "Person",
    "name": "张三",
    "properties": {
      "title": "产品经理",
      "department": "技术部",
      "email": "zhangsan@example.com"
    },
    "createdAt": "2025-01-03T11:00:00Z"
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "trace-001"
}
```

### 2.2 批量创建实体

**接口地址**: `POST /api/v1/knowledge-graph/entities/batch`

**请求参数**:
```json
{
  "entities": [
    {
      "entity_type": "Person",
      "name": "张三",
      "properties": {"title": "产品经理"}
    },
    {
      "entity_type": "Organization",
      "name": "某某科技公司",
      "properties": {"industry": "软件"}
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "批量创建成功",
  "data": {
    "created_count": 2,
    "entity_ids": ["entity-001", "entity-002"]
  }
}
```

### 2.3 获取实体详情

**接口地址**: `GET /api/v1/knowledge-graph/entities/{entity_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": "entity-001",
    "type": "Person",
    "name": "张三",
    "properties": {
      "title": "产品经理",
      "department": "技术部"
    },
    "relationships": {
      "incoming": 5,
      "outgoing": 3
    },
    "createdAt": "2025-01-03T11:00:00Z",
    "updatedAt": "2025-01-03T11:00:00Z"
  }
}
```

### 2.4 更新实体

**接口地址**: `PUT /api/v1/knowledge-graph/entities/{entity_id}`

**请求参数**:
```json
{
  "properties": {
    "title": "高级产品经理",
    "department": "产品部"
  }
}
```

### 2.5 删除实体

**接口地址**: `DELETE /api/v1/knowledge-graph/entities/{entity_id}`

**查询参数**:
- cascade: 是否级联删除关系 (默认false)

**响应示例**:
```json
{
  "code": 0,
  "message": "实体删除成功",
  "data": {
    "deleted_entity_id": "entity-001",
    "deleted_relationships_count": 8
  }
}
```

### 2.6 搜索实体

**接口地址**: `GET /api/v1/knowledge-graph/entities/search`

**查询参数**:
- query: 搜索关键词
- entity_type: 实体类型过滤
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "搜索成功",
  "data": {
    "entities": [
      {
        "id": "entity-001",
        "type": "Person",
        "name": "张三",
        "highlight": "张三"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

---

## 3. 关系管理API

### 3.1 创建关系

**接口地址**: `POST /api/v1/knowledge-graph/relationships`

**请求参数**:
```json
{
  "from_entity": "entity-001",
  "to_entity": "entity-002",
  "relation_type": "works_for",
  "properties": {
    "since": "2020-01-01",
    "role": "产品经理"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "关系创建成功",
  "data": {
    "id": "rel-001",
    "from_entity": "entity-001",
    "to_entity": "entity-002",
    "type": "works_for",
    "properties": {
      "since": "2020-01-01",
      "role": "产品经理"
    },
    "createdAt": "2025-01-03T11:00:00Z"
  }
}
```

### 3.2 批量创建关系

**接口地址**: `POST /api/v1/knowledge-graph/relationships/batch`

**请求参数**:
```json
{
  "relationships": [
    {
      "from_entity": "entity-001",
      "to_entity": "entity-002",
      "relation_type": "works_for",
      "properties": {}
    },
    {
      "from_entity": "entity-002",
      "to_entity": "entity-003",
      "relation_type": "located_in",
      "properties": {}
    }
  ]
}
```

### 3.3 获取关系详情

**接口地址**: `GET /api/v1/knowledge-graph/relationships/{relationship_id}`

### 3.4 更新关系

**接口地址**: `PUT /api/v1/knowledge-graph/relationships/{relationship_id}`

### 3.5 删除关系

**接口地址**: `DELETE /api/v1/knowledge-graph/relationships/{relationship_id}`

---

## 4. 图谱构建API

### 4.1 从文本抽取实体

**接口地址**: `POST /api/v1/knowledge-graph/extract/entities`

**请求参数**:
```json
{
  "text": "张三是某某科技公司的产品经理，负责技术部的产品规划工作。",
  "entity_types": ["Person", "Organization"],
  "confidence_threshold": 0.7
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "实体抽取成功",
  "data": {
    "entities": [
      {
        "text": "张三",
        "type": "Person",
        "confidence": 0.95,
        "start_pos": 0,
        "end_pos": 2
      },
      {
        "text": "某某科技公司",
        "type": "Organization",
        "confidence": 0.92,
        "start_pos": 4,
        "end_pos": 9
      }
    ],
    "total": 2
  }
}
```

### 4.2 从文本抽取关系

**接口地址**: `POST /api/v1/knowledge-graph/extract/relationships`

**请求参数**:
```json
{
  "text": "张三在某某科技公司担任产品经理。",
  "entities": [
    {"id": "entity-001", "text": "张三", "type": "Person"},
    {"id": "entity-002", "text": "某某科技公司", "type": "Organization"}
  ],
  "confidence_threshold": 0.7
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "关系抽取成功",
  "data": {
    "relationships": [
      {
        "from_entity": "entity-001",
        "to_entity": "entity-002",
        "type": "works_for",
        "confidence": 0.89,
        "evidence": "在...担任..."
      }
    ],
    "total": 1
  }
}
```

### 4.3 从文本构建图谱

**接口地址**: `POST /api/v1/knowledge-graph/build/from-text`

**请求参数**:
```json
{
  "text": "张三是某某科技公司的产品经理，负责技术部产品规划。",
  "auto_merge": true,
  "merge_threshold": 0.85
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "图谱构建成功",
  "data": {
    "created_entities": 3,
    "created_relationships": 2,
    "merged_entities": 0,
    "graph_id": "graph-001",
    "entities": [
      {
        "id": "entity-001",
        "type": "Person",
        "name": "张三"
      }
    ]
  }
}
```

### 4.4 知识融合

**接口地址**: `POST /api/v1/knowledge-graph/fusion`

**请求参数**:
```json
{
  "entity_ids": ["entity-001", "entity-002"],
  "strategy": "highest_confidence",
  "merge_properties": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "知识融合成功",
  "data": {
    "merged_entity_id": "entity-001",
    "merged_entity_ids": ["entity-002"],
    "merged_relationships": 5
  }
}
```

---

## 5. 图谱查询API

### 5.1 获取实体邻居

**接口地址**: `GET /api/v1/knowledge-graph/entities/{entity_id}/neighbors`

**查询参数**:
- direction: 方向 (all/incoming/outgoing)
- relation_types: 关系类型过滤
- depth: 深度 (默认1)
- limit: 数量限制

**响应示例**:
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "entity_id": "entity-001",
    "neighbors": [
      {
        "entity": {
          "id": "entity-002",
          "type": "Organization",
          "name": "某某科技公司"
        },
        "relationship": {
          "type": "works_for",
          "direction": "outgoing"
        }
      }
    ],
    "total": 1
  }
}
```

### 5.2 查找实体间路径

**接口地址**: `POST /api/v1/knowledge-graph/path/find`

**请求参数**:
```json
{
  "from_entity": "entity-001",
  "to_entity": "entity-003",
  "max_depth": 3,
  "path_type": "shortest"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "路径查找成功",
  "data": {
    "paths": [
      {
        "length": 2,
        "nodes": ["entity-001", "entity-002", "entity-003"],
        "edges": [
          {
            "from": "entity-001",
            "to": "entity-002",
            "type": "works_for"
          },
          {
            "from": "entity-002",
            "to": "entity-003",
            "type": "located_in"
          }
        ]
      }
    ],
    "total": 1
  }
}
```

### 5.3 执行Cypher查询

**接口地址**: `POST /api/v1/knowledge-graph/query/cypher`

**请求参数**:
```json
{
  "query": "MATCH (p:Person {name: '张三'})-[:works_for]->(o:Organization) RETURN p, o",
  "parameters": {}
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "查询成功",
  "data": {
    "results": [
      {
        "p": {
          "id": "entity-001",
          "type": "Person",
          "name": "张三"
        },
        "o": {
          "id": "entity-002",
          "type": "Organization",
          "name": "某某科技公司"
        }
      }
    ],
    "total": 1
  }
}
```

### 5.4 图谱统计

**接口地址**: `GET /api/v1/knowledge-graph/stats`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_entities": 1500,
    "total_relationships": 3500,
    "by_entity_type": {
      "Person": 500,
      "Organization": 300,
      "Location": 400,
      "Product": 300
    },
    "by_relation_type": {
      "works_for": 800,
      "located_in": 600,
      "related_to": 1200
    },
    "avg_degree": 4.67
  }
}
```

---

## 6. 图谱推理API

### 6.1 规则推理

**接口地址**: `POST /api/v1/knowledge-graph/inference/rule`

**请求参数**:
```json
{
  "rule": "IF Person A works_for Organization B AND Organization B located_in Location C THEN Person A located_in Location C",
  "entity_id": "entity-001"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "推理完成",
  "data": {
    "inferred_relationships": [
      {
        "from_entity": "entity-001",
        "to_entity": "entity-005",
        "type": "located_in",
        "confidence": 0.95,
        "reasoning": "通过works_for和located_in关系推理得出"
      }
    ],
    "total": 1
  }
}
```

### 6.2 推荐潜在关系

**接口地址**: `POST /api/v1/knowledge-graph/inference/suggest`

**请求参数**:
```json
{
  "entity_id": "entity-001",
  "max_suggestions": 10,
  "confidence_threshold": 0.6
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "推荐完成",
  "data": {
    "suggestions": [
      {
        "to_entity": "entity-010",
        "relation_type": "knows",
        "confidence": 0.82,
        "reasoning": "共同同事关系"
      }
    ],
    "total": 1
  }
}
```

---

## 7. 数据模型

### 7.1 Entity

```typescript
interface Entity {
  id: string;
  type: EntityType;
  name: string;
  properties?: Record<string, any>;
  created_at: Date;
  updated_at?: Date;
}

type EntityType =
  | 'Person'
  | 'Organization'
  | 'Location'
  | 'Product'
  | 'Event'
  | 'Concept';
```

### 7.2 Relationship

```typescript
interface Relationship {
  id: string;
  from_entity: string;
  to_entity: string;
  type: RelationType;
  properties?: Record<string, any>;
  confidence?: number;
  created_at: Date;
  updated_at?: Date;
}

type RelationType =
  | 'works_for'
  | 'located_in'
  | 'knows'
  | 'related_to'
  | 'part_of'
  | 'custom';
```

### 7.3 ExtractedEntity

```typescript
interface ExtractedEntity {
  text: string;
  type: EntityType;
  confidence: number;
  start_pos: number;
  end_pos: number;
  properties?: Record<string, any>;
}
```

### 7.4 ExtractedRelationship

```typescript
interface ExtractedRelationship {
  from_entity: string;
  to_entity: string;
  type: RelationType;
  confidence: number;
  evidence: string;
}
```

### 7.5 GraphPath

```typescript
interface GraphPath {
  length: number;
  nodes: string[];
  edges: Array<{
    from: string;
    to: string;
    type: RelationType;
  }>;
}
```

---

## 8. 后端实现示例

### 8.1 KnowledgeGraphEngine 核心实现

```go
package knowledgegraph

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type KnowledgeGraphEngine struct {
	neo4jDriver neo4j.DriverWithContext
	entityRepo  EntityRepository
	llmClient   LLMClient
}

// ExtractEntities 从文本抽取实体
func (e *KnowledgeGraphEngine) ExtractEntities(
	ctx context.Context,
	text string,
	entityTypes []string,
	confidenceThreshold float64,
) ([]*ExtractedEntity, error) {
	// 1. 构建LLM提示
	prompt := fmt.Sprintf(
		"从以下文本中识别实体（%v），返回JSON格式：\n%s",
		entityTypes, text,
	)

	// 2. 调用LLM
	resp, err := e.llmClient.Chat(ctx, &LLMChatRequest{
		Model: "gpt-4",
		Messages: []LLMMessage{
			{Role: "system", Content: "你是实体识别专家"},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	// 3. 解析LLM输出
	entities, err := e.parseEntities(resp.Content)
	if err != nil {
		return nil, err
	}

	// 4. 过滤低置信度实体
	result := make([]*ExtractedEntity, 0)
	for _, entity := range entities {
		if entity.Confidence >= confidenceThreshold {
			result = append(result, entity)
		}
	}

	return result, nil
}

// ExtractRelationships 从文本抽取关系
func (e *KnowledgeGraphEngine) ExtractRelationships(
	ctx context.Context,
	text string,
	entities []*Entity,
	confidenceThreshold float64,
) ([]*ExtractedRelationship, error) {
	// 1. 构建实体列表
	entityList := make([]string, len(entities))
	for i, e := range entities {
		entityList[i] = fmt.Sprintf("%s(%s)", e.Name, e.Type)
	}

	// 2. 构建LLM提示
	prompt := fmt.Sprintf(
		"识别实体间的关系。实体：%v\n文本：%s\n返回JSON格式",
		entityList, text,
	)

	// 3. 调用LLM
	resp, err := e.llmClient.Chat(ctx, &LLMChatRequest{
		Model: "gpt-4",
		Messages: []LLMMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return nil, err
	}

	// 4. 解析关系
	relationships, err := e.parseRelationships(resp.Content)
	if err != nil {
		return nil, err
	}

	// 5. 过滤低置信度关系
	result := make([]*ExtractedRelationship, 0)
	for _, rel := range relationships {
		if rel.Confidence >= confidenceThreshold {
			result = append(result, rel)
		}
	}

	return result, nil
}

// BuildGraph 构建图谱
func (e *KnowledgeGraphEngine) BuildGraph(
	ctx context.Context,
	entities []*Entity,
	relationships []*Relationship,
) error {
	session := e.neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	// 1. 创建实体节点
	for _, entity := range entities {
		_, err := session.Run(ctx,
			`MERGE (e:Entity {id: $id})
			 SET e.type = $type, e.name = $name, e.properties = $props`,
			map[string]interface{}{
				"id":    entity.ID,
				"type":  entity.Type,
				"name":  entity.Name,
				"props": entity.Properties,
			},
		)
		if err != nil {
			return fmt.Errorf("创建实体失败: %w", err)
		}
	}

	// 2. 创建关系边
	for _, rel := range relationships {
		_, err := session.Run(ctx,
			`MATCH (from:Entity {id: $fromId}), (to:Entity {id: $toId})
			 MERGE (from)-[r:RELATIONSHIP {type: $type}]->(to)
			 SET r.properties = $props`,
			map[string]interface{}{
				"fromId": rel.FromEntity,
				"toId":   rel.ToEntity,
				"type":   rel.RelationType,
				"props":  rel.Properties,
			},
		)
		if err != nil {
			return fmt.Errorf("创建关系失败: %w", err)
		}
	}

	return nil
}

// FindPath 查找实体间路径
func (e *KnowledgeGraphEngine) FindPath(
	ctx context.Context,
	fromEntity string,
	toEntity string,
	maxDepth int,
) ([]*GraphPath, error) {
	session := e.neo4jDriver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	// 执行最短路径查询
	result, err := session.Run(ctx,
		`MATCH path = shortestPath(
			 (from:Entity {id: $fromId})-[*1..$maxDepth]-(to:Entity {id: $toId})
		   )
		 RETURN path`,
		map[string]interface{}{
			"fromId":   fromEntity,
			"toId":     toEntity,
			"maxDepth": maxDepth,
		},
	)
	if err != nil {
		return nil, err
	}

	// 解析路径结果
	paths := make([]*GraphPath, 0)
	for result.Next(ctx) {
		record := result.Record()
		path := record.Values[0].(neo4j.Path)

		// 提取节点
		nodes := make([]string, len(path.Nodes))
		for i, node := range path.Nodes {
			nodes[i] = node.ElementId
		}

		// 提取边
		edges := make([]Edge, len(path.Relationships))
		for i, rel := range path.Relationships {
			edges[i] = Edge{
				From: rel.StartNodeElementId,
				To:   rel.EndNodeElementId,
				Type: rel.Type,
			}
		}

		paths = append(paths, &GraphPath{
			Length: len(path.Relationships),
			Nodes:  nodes,
			Edges:  edges,
		})
	}

	return paths, nil
}
```

### 8.2 HTTP Handler

```go
// ExtractEntitiesHandler 实体抽取
func ExtractEntitiesHandler(ctx context.Context, c *app.RequestContext) {
	var req ExtractEntitiesRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	entities, err := knowledgeGraphEngine.ExtractEntities(
		ctx,
		req.Text,
		req.EntityTypes,
		req.ConfidenceThreshold,
	)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "实体抽取失败",
		})
		return
	}

	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    0,
		Message: "实体抽取成功",
		Data: ExtractEntitiesResponse{
			Entities: entities,
			Total:    len(entities),
		},
	})
}

// FindPathHandler 路径查找
func FindPathHandler(ctx context.Context, c *app.RequestContext) {
	var req FindPathRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	paths, err := knowledgeGraphEngine.FindPath(
		ctx,
		req.FromEntity,
		req.ToEntity,
		req.MaxDepth,
	)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "路径查找失败",
		})
		return
	}

	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    0,
		Message: "路径查找成功",
		Data: FindPathResponse{
			Paths: paths,
			Total: len(paths),
		},
	})
}
```

---

## 9. 前端调用示例

### 9.1 创建实体

```typescript
import { request } from '@/utils/request';

export async function createEntity(params: {
  entityType: string;
  name: string;
  properties?: Record<string, any>;
}) {
  return request.post<Entity>('/api/v1/knowledge-graph/entities', {
    entity_type: params.entityType,
    name: params.name,
    properties: params.properties,
  });
}

// 使用示例
const entity = await createEntity({
  entityType: 'Person',
  name: '张三',
  properties: {
    title: '产品经理',
    department: '技术部',
  },
});
```

### 9.2 实体抽取

```typescript
export async function extractEntities(params: {
  text: string;
  entityTypes?: string[];
  confidenceThreshold?: number;
}) {
  return request.post<{
    entities: ExtractedEntity[];
    total: number;
  }>('/api/v1/knowledge-graph/extract/entities', {
    text: params.text,
    entity_types: params.entityTypes,
    confidence_threshold: params.confidenceThreshold || 0.7,
  });
}

// 使用示例
const result = await extractEntities({
  text: '张三是某某科技公司的产品经理',
  entityTypes: ['Person', 'Organization'],
});

console.log('抽取的实体:', result.data.entities);
```

### 9.3 查找路径

```typescript
export async function findPath(params: {
  fromEntity: string;
  toEntity: string;
  maxDepth?: number;
}) {
  return request.post<{
    paths: GraphPath[];
    total: number;
  }>('/api/v1/knowledge-graph/path/find', {
    from_entity: params.fromEntity,
    to_entity: params.toEntity,
    max_depth: params.maxDepth || 3,
  });
}

// 使用示例
const paths = await findPath({
  fromEntity: 'entity-001',
  toEntity: 'entity-003',
  maxDepth: 3,
});

console.log('找到的路径:', paths.data.paths);
```

### 9.4 React Hook 示例

```typescript
import { useState, useEffect } from 'react';
import { extractEntities, findPath } from '@/api/knowledge-graph';

export function useEntityExtraction() {
  const [entities, setEntities] = useState<ExtractedEntity[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const extract = async (text: string) => {
    setLoading(true);
    setError(null);

    try {
      const result = await extractEntities({ text });
      setEntities(result.data.entities);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  };

  return { entities, loading, error, extract };
}

// 使用示例
function EntityExtractionComponent() {
  const { entities, loading, extract } = useEntityExtraction();

  return (
    <div>
      <textarea
        placeholder="输入文本..."
        onChange={(e) => extract(e.target.value)}
      />
      {loading && <div>识别中...</div>}
      <ul>
        {entities.map((entity, index) => (
          <li key={index}>
            [{entity.type}] {entity.text} (置信度: {entity.confidence.toFixed(2)})
          </li>
        ))}
      </ul>
    </div>
  );
}
```

---

## 10. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 90001 | 404 | 实体不存在 |
| 90002 | 400 | 实体类型无效 |
| 90003 | 400 | 实体名称重复 |
| 90004 | 404 | 关系不存在 |
| 90005 | 400 | 关系类型无效 |
| 90006 | 400 | 实体不存在，无法创建关系 |
| 90007 | 500 | LLM实体抽取失败 |
| 90008 | 500 | LLM关系抽取失败 |
| 90009 | 400 | Cypher查询语法错误 |
| 90010 | 400 | 路径深度超限 |
| 90101 | 500 | Neo4j连接失败 |
| 90102 | 500 | 图数据库写入失败 |

---

**文档结束**
