# 知识图谱引擎 API 文档

## 📅 版本信息
- **版本**: v1.0
- **日期**: 2025-01-01
- **Base URL**: `/api/knowledge_graph`

---

## 📋 目录

1. [图谱构建API](#1-图谱构建api)
2. [图谱查询API](#2-图谱查询api)
3. [图谱推理API](#3-图谱推理api)
4. [图谱可视化API](#4-图谱可视化api)
5. [数据模型](#5-数据模型)
6. [错误码](#6-错误码)

---

## 1. 图谱构建API

### 1.1 从文本构建图谱

从文本中自动抽取实体和关系，构建完整的知识图谱。

**请求**
```http
POST /api/knowledge_graph/build_from_text
Content-Type: application/json
```

**请求体**
```json
{
  "text": "张三是阿里巴巴的高级工程师，负责人工智能算法研究。他在清华大学获得了计算机科学博士学位。"
}
```

**响应**
```json
{
  "code": 0,
  "message": "知识图谱构建成功"
}
```

**字段说明**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| text | string | 是 | 待处理的文本内容 |

**错误码**
- `400`: 文本内容为空
- `500`: LLM调用失败

---

### 1.2 抽取实体

从文本中抽取所有实体（不保存到数据库）。

**请求**
```http
POST /api/knowledge_graph/extract_entities
Content-Type: application/json
```

**请求体**
```json
{
  "text": "张三是阿里巴巴的高级工程师"
}
```

**响应**
```json
{
  "code": 0,
  "entities": [
    {
      "id": "uuid-1",
      "entity_type": "person",
      "entity_name": "张三",
      "properties": {}
    },
    {
      "id": "uuid-2",
      "entity_type": "organization",
      "entity_name": "阿里巴巴",
      "properties": {}
    }
  ],
  "count": 2
}
```

**实体类型（EntityType）**

| 类型 | 说明 | 示例 |
|------|------|------|
| `person` | 人物 | 张三、李四 |
| `organization` | 组织/机构 | 阿里巴巴、清华大学 |
| `location` | 地点 | 北京、上海 |
| `concept` | 概念 | 人工智能、机器学习 |
| `event` | 事件 | 2024年奥运会 |
| `product` | 产品 | iPhone、Tesla |
| `document` | 文档 | 论文、报告 |

---

### 1.3 抽取关系

从文本中抽取实体间的关系（不保存到数据库）。

**请求**
```http
POST /api/knowledge_graph/extract_relationships
Content-Type: application/json
```

**请求体**
```json
{
  "text": "张三是阿里巴巴的高级工程师",
  "entities": [
    {
      "id": "uuid-1",
      "entity_type": "person",
      "entity_name": "张三",
      "properties": {}
    },
    {
      "id": "uuid-2",
      "entity_type": "organization",
      "entity_name": "阿里巴巴",
      "properties": {}
    }
  ]
}
```

**响应**
```json
{
  "code": 0,
  "relationships": [
    {
      "id": "rel-uuid-1",
      "source_entity_id": "uuid-1",
      "target_entity_id": "uuid-2",
      "relation_type": "works_for",
      "properties": {
        "since": "2020"
      },
      "weight": 0.8
    }
  ],
  "count": 1
}
```

**关系类型（RelationType）**

| 类型 | 说明 | 适用场景 |
|------|------|----------|
| `knows` | 认识 | 人物关系 |
| `works_for` | 任职于 | 人物→组织 |
| `colleague_of` | 同事 | 人物关系 |
| `friend_of` | 朋友 | 人物关系 |
| `family_of` | 家人 | 人物关系 |
| `spouse_of` | 配偶 | 人物关系 |
| `child_of` | 子女 | 人物关系 |
| `parent_of` | 父母 | 人物关系 |
| `part_of` | 属于 | 组织、地点 |
| `subsidiary_of` | 子公司 | 组织关系 |
| `partner_of` | 合作伙伴 | 组织关系 |
| `competitor_of` | 竞争对手 | 组织关系 |
| `investor_in` | 投资者 | 组织关系 |
| `located_in` | 位于 | 地理关系 |
| `contains` | 包含 | 地理关系 |
| `borders` | 毗邻 | 地理关系 |
| `related_to` | 相关 | 通用关系 |
| `similar_to` | 相似 | 概念关系 |
| `instance_of` | 实例 | 概念关系 |
| `subclass_of` | 子类 | 概念关系 |
| `participated_in` | 参与 | 事件关系 |
| `caused` | 导致 | 事件关系 |
| `preceded` | 先于 | 事件关系 |
| `cites` | 引用 | 文档关系 |
| `references` | 参考 | 文档关系 |
| `about` | 关于 | 文档关系 |

---

### 1.4 添加实体

手动添加单个实体到知识图谱。

**请求**
```http
POST /api/knowledge_graph/add_entity
Content-Type: application/json
```

**请求体**
```json
{
  "entity_name": "张三",
  "entity_type": "person",
  "properties": {
    "age": 30,
    "gender": "男",
    "description": "高级工程师"
  }
}
```

**响应**
```json
{
  "code": 0,
  "entity_id": "uuid-xxx",
  "message": "实体添加成功"
}
```

---

### 1.5 添加关系

手动添加单个关系到知识图谱。

**请求**
```http
POST /api/knowledge_graph/add_relationship
Content-Type: application/json
```

**请求体**
```json
{
  "source_entity_id": "uuid-1",
  "target_entity_id": "uuid-2",
  "relation_type": "works_for",
  "properties": {
    "since": "2020",
    "position": "高级工程师"
  },
  "weight": 0.8
}
```

**响应**
```json
{
  "code": 0,
  "relationship_id": "rel-uuid-xxx",
  "message": "关系添加成功"
}
```

---

### 1.6 删除实体

删除指定实体及其所有相关关系。

**请求**
```http
POST /api/knowledge_graph/delete_entity
Content-Type: application/json
```

**请求体**
```json
{
  "entity_id": "uuid-xxx"
}
```

**响应**
```json
{
  "code": 0,
  "message": "实体删除成功"
}
```

---

### 1.7 获取实体统计

获取当前租户的实体统计信息。

**请求**
```http
POST /api/knowledge_graph/entity_stats
Content-Type: application/json
```

**请求体**
```json
{}
```

**响应**
```json
{
  "code": 0,
  "stats": {
    "person": 150,
    "organization": 50,
    "location": 30,
    "concept": 200,
    "event": 20,
    "product": 40,
    "document": 100
  }
}
```

---

## 2. 图谱查询API

### 2.1 查询实体

根据实体ID查询实体的详细信息。

**请求**
```http
POST /api/knowledge_graph/query_entity
Content-Type: application/json
```

**请求体**
```json
{
  "entity_id": "uuid-xxx"
}
```

**响应**
```json
{
  "code": 0,
  "entity": {
    "id": "uuid-xxx",
    "entity_type": "person",
    "entity_name": "张三",
    "properties": {
      "age": 30,
      "gender": "男",
      "description": "高级工程师"
    }
  }
}
```

---

### 2.2 查询邻居

查询实体的N度邻居（BFS遍历）。

**请求**
```http
POST /api/knowledge_graph/query_neighbors
Content-Type: application/json
```

**请求体**
```json
{
  "entity_id": "uuid-xxx",
  "depth": 2
}
```

**响应**
```json
{
  "code": 0,
  "center_node": {
    "id": "uuid-xxx",
    "entity_name": "张三",
    "entity_type": "person",
    "properties": {}
  },
  "neighbors": [
    {
      "id": "uuid-1",
      "entity_type": "organization",
      "entity_name": "阿里巴巴",
      "properties": {}
    },
    {
      "id": "uuid-2",
      "entity_type": "person",
      "entity_name": "李四",
      "properties": {}
    }
  ],
  "relationships": [
    {
      "id": "rel-uuid-1",
      "source_id": "uuid-xxx",
      "target_id": "uuid-1",
      "relation_type": "works_for",
      "properties": {}
    },
    {
      "id": "rel-uuid-2",
      "source_id": "uuid-xxx",
      "target_id": "uuid-2",
      "relation_type": "knows",
      "properties": {}
    }
  ],
  "depth": 2
}
```

**参数说明**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| entity_id | string | 是 | - | 中心实体ID |
| depth | int | 是 | - | 遍历深度（1-5） |

---

### 2.3 查询路径

查询两个实体之间的所有路径（BFS算法）。

**请求**
```http
POST /api/knowledge_graph/query_path
Content-Type: application/json
```

**请求体**
```json
{
  "source_entity_id": "uuid-1",
  "target_entity_id": "uuid-2",
  "max_depth": 5
}
```

**响应**
```json
{
  "code": 0,
  "paths": [
    {
      "nodes": ["uuid-1", "uuid-3", "uuid-4", "uuid-2"],
      "relationships": ["rel-1", "rel-2", "rel-3"],
      "length": 3,
      "weight": 0.75
    },
    {
      "nodes": ["uuid-1", "uuid-5", "uuid-2"],
      "relationships": ["rel-4", "rel-5"],
      "length": 2,
      "weight": 0.85
    }
  ],
  "count": 2
}
```

---

### 2.4 搜索实体

根据关键词或实体类型搜索实体。

**请求**
```http
POST /api/knowledge_graph/search_entities
Content-Type: application/json
```

**请求体**
```json
{
  "keyword": "张三",
  "entity_type": "person",
  "limit": 20
}
```

**响应**
```json
{
  "code": 0,
  "entities": [
    {
      "id": "uuid-xxx",
      "entity_type": "person",
      "entity_name": "张三",
      "properties": {}
    }
  ],
  "count": 1
}
```

**参数说明**

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| keyword | string | 否 | - | 搜索关键词（模糊匹配） |
| entity_type | string | 否 | - | 实体类型过滤 |
| limit | int | 是 | 20 | 返回结果数量限制（1-100） |

---

### 2.5 获取图谱统计

获取当前租户的图谱统计信息。

**请求**
```http
POST /api/knowledge_graph/statistics
Content-Type: application/json
```

**请求体**
```json
{}
```

**响应**
```json
{
  "code": 0,
  "statistics": {
    "node_count": 1000,
    "relationship_count": 5000,
    "avg_degree": 10.0,
    "density": 0.01,
    "connected_components": 5
  }
}
```

**字段说明**

| 字段 | 类型 | 说明 |
|------|------|------|
| node_count | int | 节点总数 |
| relationship_count | int | 关系总数 |
| avg_degree | float32 | 平均度数（每个节点的平均连接数） |
| density | float32 | 图密度（0-1之间） |
| connected_components | int | 连通分量数 |

---

### 2.6 查找最短路径

查找两个实体之间的最短路径。

**请求**
```http
POST /api/knowledge_graph/shortest_path
Content-Type: application/json
```

**请求体**
```json
{
  "source_entity_id": "uuid-1",
  "target_entity_id": "uuid-2",
  "max_depth": 10
}
```

**响应**
```json
{
  "code": 0,
  "path": {
    "nodes": ["uuid-1", "uuid-3", "uuid-2"],
    "relationships": ["rel-1", "rel-2"],
    "length": 2,
    "weight": 0.8
  }
}
```

---

## 3. 图谱推理API

### 3.1 推断关系

推断两个实体之间的潜在关系。

**请求**
```http
POST /api/knowledge_graph/infer_relationship
Content-Type: application/json
```

**请求体**
```json
{
  "entity1_id": "uuid-1",
  "entity2_id": "uuid-2"
}
```

**响应**
```json
{
  "code": 0,
  "relationship": {
    "id": "rel-uuid-xxx",
    "source_entity_id": "uuid-1",
    "target_entity_id": "uuid-2",
    "relation_type": "works_for",
    "weight": 0.75,
    "properties": {
      "inferred": true,
      "confidence": 0.75
    }
  }
}
```

---

### 3.2 推断实体类型

根据文本推断实体类型。

**请求**
```http
POST /api/knowledge_graph/infer_entity_type
Content-Type: application/json
```

**请求体**
```json
{
  "text": "张三是阿里巴巴的高级工程师"
}
```

**响应**
```json
{
  "code": 0,
  "entity_type": "person",
  "confidence": 0.9
}
```

---

### 3.3 推荐相关实体

基于图谱结构推荐相关实体（2跳邻居）。

**请求**
```http
POST /api/knowledge_graph/recommend_entities
Content-Type: application/json
```

**请求体**
```json
{
  "entity_id": "uuid-xxx",
  "top_k": 10
}
```

**响应**
```json
{
  "code": 0,
  "entities": [
    {
      "id": "uuid-1",
      "entity_type": "person",
      "entity_name": "李四",
      "properties": {}
    }
  ],
  "count": 10
}
```

---

### 3.4 查找相似实体

基于实体属性查找相似实体（Jaccard相似度）。

**请求**
```http
POST /api/knowledge_graph/find_similar
Content-Type: application/json
```

**请求体**
```json
{
  "entity_id": "uuid-xxx",
  "top_k": 10
}
```

**响应**
```json
{
  "code": 0,
  "entities": [
    {
      "id": "uuid-1",
      "entity_type": "person",
      "entity_name": "王五",
      "properties": {},
      "similarity": 0.85
    }
  ],
  "count": 5
}
```

---

## 4. 图谱可视化API

### 4.1 获取可视化数据

获取前端图谱可视化所需的JSON格式数据。

**请求**
```http
POST /api/knowledge_graph/visualization_data
Content-Type: application/json
```

**请求体**
```json
{
  "center_entity_id": "uuid-xxx",
  "depth": 2
}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "nodes": [
      {
        "id": "uuid-xxx",
        "label": "张三",
        "type": "person",
        "size": 25,
        "color": "#4ECDC4",
        "properties": {},
        "x": 100,
        "y": 200
      }
    ],
    "links": [
      {
        "source": "uuid-xxx",
        "target": "uuid-1",
        "label": "works_for",
        "type": "works_for",
        "weight": 0.8,
        "color": "#E0E0E0"
      }
    ]
  }
}
```

---

### 4.2 路径可视化

生成路径的可视化数据。

**请求**
```http
POST /api/knowledge_graph/path_visualization
Content-Type: application/json
```

**请求体**
```json
{
  "source_entity_id": "uuid-1",
  "target_entity_id": "uuid-2"
}
```

**响应**
```json
{
  "code": 0,
  "data": {
    "nodes": [...],
    "links": [...]
  }
}
```

---

## 5. 数据模型

### 5.1 GraphEntity（图实体）

```typescript
interface GraphEntity {
  id: string;                    // 实体ID（UUID）
  tenant_id: string;             // 租户ID
  entity_type: EntityType;       // 实体类型
  entity_name: string;           // 实体名称
  properties: EntityProperties;  // 实体属性（JSON对象）
  embedding?: number[];          // 向量嵌入（可选）
  created_at: number;            // 创建时间戳
  updated_at: number;            // 更新时间戳
  deleted_at?: number;           // 删除时间戳（软删除）
}
```

### 5.2 GraphRelationship（图关系）

```typescript
interface GraphRelationship {
  id: string;                       // 关系ID（UUID）
  tenant_id: string;                // 租户ID
  source_entity_id: string;         // 源实体ID
  target_entity_id: string;         // 目标实体ID
  relation_type: RelationType;      // 关系类型
  properties: RelationshipProperties; // 关系属性（JSON对象）
  weight: number;                   // 关系权重（0-1）
  created_at: number;               // 创建时间戳
  deleted_at?: number;              // 删除时间戳（软删除）
}
```

### 5.3 GraphQuery（图查询）

```typescript
interface GraphQuery {
  id: string;              // 查询ID（UUID）
  tenant_id: string;       // 租户ID
  query_name: string;      // 查询名称
  query_text: string;      // Cypher查询语句
  description?: string;    // 查询描述
  variables: QueryVariables; // 查询变量（JSON对象）
  is_public: boolean;      // 是否公开
  created_by: string;      // 创建者ID
  created_at: number;      // 创建时间戳
  updated_at: number;      // 更新时间戳
}
```

### 5.4 VisualizationData（可视化数据）

```typescript
interface VisualizationData {
  nodes: VisNode[];  // 节点数组
  links: VisLink[];  // 链接数组
}

interface VisNode {
  id: string;                    // 节点ID
  label: string;                 // 节点标签（显示文本）
  type: EntityType;              // 节点类型
  size?: number;                 // 节点大小
  color?: string;                // 节点颜色
  properties?: EntityProperties; // 节点属性
  x?: number;                    // X坐标
  y?: number;                    // Y坐标
}

interface VisLink {
  source: string;                // 源节点ID
  target: string;                // 目标节点ID
  label?: string;                // 链接标签
  type: RelationType;            // 关系类型
  weight?: number;               // 权重
  color?: string;                // 链接颜色
  properties?: RelationshipProperties; // 关系属性
}
```

---

## 6. 错误码

### 6.1 通用错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| `0` | 200 | 成功 |
| `400001` | 400 | 请求参数无效 |
| `500001` | 500 | 内部服务器错误 |

### 6.2 图谱构建错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| `205001` | 400 | 文本内容为空 |
| `205002` | 500 | LLM调用失败 |
| `205003` | 400 | 实体类型无效 |
| `205004` | 400 | 关系类型无效 |
| `205005` | 400 | 实体名称重复 |
| `205006` | 404 | 实体不存在 |

### 6.3 图谱查询错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| `206001` | 404 | 实体不存在 |
| `206002` | 400 | 查询深度超出限制 |
| `206003` | 404 | 未找到路径 |
| `206004` | 400 | 关键词过长 |

### 6.4 错误响应格式

```json
{
  "code": 205001,
  "message": "Text content is empty",
  "message_zh": "文本内容为空",
  "request_id": "uuid-xxx",
  "timestamp": "2025-01-01T00:00:00Z",
  "details": {
    "field": "text",
    "error": "required field is empty"
  }
}
```

---

## 📚 使用示例

### 示例1：构建人物关系图谱

```bash
curl -X POST http://localhost:8080/api/knowledge_graph/build_from_text \
  -H "Content-Type: application/json" \
  -d '{
    "text": "张三是阿里巴巴的高级工程师，负责人工智能算法研究。他的同事李四是清华大学毕业的博士，两人一起参与了多个AI项目。"
  }'
```

### 示例2：查询实体的2度邻居

```bash
curl -X POST http://localhost:8080/api/knowledge_graph/query_neighbors \
  -H "Content-Type: application/json" \
  -d '{
    "entity_id": "uuid-xxx",
    "depth": 2
  }'
```

### 示例3：查找两个实体之间的最短路径

```bash
curl -X POST http://localhost:8080/api/knowledge_graph/shortest_path \
  -H "Content-Type: application/json" \
  -d '{
    "source_entity_id": "uuid-1",
    "target_entity_id": "uuid-2"
  }'
```

---

## 🔐 认证说明

所有API都需要在请求头中包含认证信息：

```http
Authorization: Bearer {access_token}
X-Tenant-ID: {tenant_id}
```

---

## 📞 技术支持

如有问题，请参考：
- [实现总结文档](./知识图谱引擎实现总结_v1.0.md)
- [企业级开发规范手册](./企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](./企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
