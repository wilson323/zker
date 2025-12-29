# API接口文档：记忆引擎模块

**模块名称**: 记忆引擎 (MemoryEngine)
**设计文档**: 17-五大AI引擎核心_记忆引擎.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 记忆存储API](#2-记忆存储api)
- [3. 记忆检索API](#3-记忆检索api)
- [4. 记忆管理API](#4-记忆管理api)
- [5. 记忆分析API](#5-记忆分析api)
- [6. 数据模型](#6-数据模型)
- [7. 后端实现示例](#7-后端实现示例)
- [8. 前端调用示例](#8-前端调用示例)
- [9. 错误码定义](#9-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

记忆引擎负责存储和检索对话历史、用户偏好和知识信息，实现长期记忆和语义检索：

- ✅ **分层记忆** - 对话记忆、用户偏好记忆、知识记忆
- ✅ **语义检索** - 基于向量相似度的智能检索
- ✅ **记忆压缩** - 长对话自动摘要压缩
- ✅ **多维度搜索** - 语义、关键词、时间范围
- ✅ **记忆过期** - 支持TTL自动清理

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 向量数据库: Milvus v2.5.10
- LLM: 通义千问/GPT-4 (embedding + summarization)

### 1.3 记忆类型

| 记忆类型 | 说明 | 存储时长 | 示例 |
|---------|------|---------|------|
| **对话记忆** | 当前对话上下文 | 会话周期 | 最近10轮对话 |
| **用户偏好** | 长期用户偏好 | 永久 | 喜欢简洁的回答 |
| **知识记忆** | 提取的知识点 | 永久 | 用户是HR经理 |

---

## 2. 记忆存储API

### 2.1 存储对话记忆

**接口地址**: `POST /api/v1/memory/conversation`

**请求参数**:
```json
{
  "conversation_id": "conv-001",
  "role": "user",
  "content": "我是HR经理，负责招聘工作",
  "metadata": {
    "timestamp": "2025-01-03T11:00:00Z",
    "message_type": "user_input"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "对话记忆存储成功",
  "data": {
    "memory_id": "mem-001",
    "vector_id": "vec-001",
    "type": "conversation",
    "createdAt": "2025-01-03T11:00:00Z"
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "trace-001"
}
```

### 2.2 存储用户偏好

**接口地址**: `POST /api/v1/memory/preference`

**请求参数**:
```json
{
  "user_id": "user-001",
  "preference_type": "communication_style",
  "preference_value": "简洁专业",
  "confidence": 0.9,
  "source": "explicit_feedback",
  "metadata": {
    "conversation_count": 15,
    "last_updated": "2025-01-03T11:00:00Z"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "用户偏好存储成功",
  "data": {
    "memory_id": "pref-001",
    "vector_id": "vec-002",
    "type": "preference",
    "createdAt": "2025-01-03T11:00:00Z"
  }
}
```

### 2.3 存储知识记忆

**接口地址**: `POST /api/v1/memory/knowledge`

**请求参数**:
```json
{
  "user_id": "user-001",
  "knowledge_type": "user_fact",
  "knowledge_content": "用户是某科技公司HR经理，负责技术岗位招聘",
  "entities": {
    "role": "HR经理",
    "industry": "科技",
    "responsibility": "技术岗位招聘"
  },
  "confidence": 0.85,
  "source_conversation_id": "conv-001"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "知识记忆存储成功",
  "data": {
    "memory_id": "know-001",
    "vector_id": "vec-003",
    "type": "knowledge",
    "createdAt": "2025-01-03T11:00:00Z"
  }
}
```

### 2.4 批量存储记忆

**接口地址**: `POST /api/v1/memory/batch`

**请求参数**:
```json
{
  "memories": [
    {
      "type": "conversation",
      "conversation_id": "conv-001",
      "role": "user",
      "content": "你好"
    },
    {
      "type": "conversation",
      "conversation_id": "conv-001",
      "role": "assistant",
      "content": "您好，有什么可以帮助您？"
    }
  ]
}
```

---

## 3. 记忆检索API

### 3.1 语义检索记忆

**接口地址**: `POST /api/v1/memory/search/semantic`

**请求参数**:
```json
{
  "query": "用户的工作是什么",
  "user_id": "user-001",
  "memory_types": ["knowledge", "preference"],
  "top_k": 5,
  "score_threshold": 0.7
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "检索成功",
  "data": {
    "memories": [
      {
        "memory_id": "know-001",
        "type": "knowledge",
        "content": "用户是某科技公司HR经理，负责技术岗位招聘",
        "score": 0.92,
        "metadata": {
          "createdAt": "2025-01-03T11:00:00Z"
        }
      }
    ],
    "total": 1
  }
}
```

### 3.2 获取对话历史

**接口地址**: `GET /api/v1/memory/conversation/{conversation_id}`

**查询参数**:
- limit: 返回数量 (默认10)
- offset: 偏移量
- order: 排序方式 (asc/desc)

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "conversation_id": "conv-001",
    "messages": [
      {
        "memory_id": "mem-001",
        "role": "user",
        "content": "你好",
        "timestamp": "2025-01-03T11:00:00Z"
      },
      {
        "memory_id": "mem-002",
        "role": "assistant",
        "content": "您好，有什么可以帮助您？",
        "timestamp": "2025-01-03T11:00:01Z"
      }
    ],
    "total": 2
  }
}
```

### 3.3 获取用户偏好

**接口地址**: `GET /api/v1/memory/user/{user_id}/preferences`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "preferences": [
      {
        "memory_id": "pref-001",
        "type": "communication_style",
        "value": "简洁专业",
        "confidence": 0.9,
        "updatedAt": "2025-01-03T11:00:00Z"
      }
    ]
  }
}
```

### 3.4 获取用户知识图谱

**接口地址**: `GET /api/v1/memory/user/{user_id}/knowledge`

**查询参数**:
- entity_type: 实体类型过滤
- limit: 返回数量

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "knowledge_graph": {
      "entities": [
        {
          "entity": "HR经理",
          "type": "role",
          "relations": [
            {
              "relation": "负责",
              "target": "技术岗位招聘"
            }
          ]
        }
      ],
      "total": 1
    }
  }
}
```

### 3.5 时间范围检索

**接口地址**: `GET /api/v1/memory/search/time-range`

**查询参数**:
- user_id: 用户ID
- start_time: 开始时间
- end_time: 结束时间
- memory_types: 记忆类型列表

**响应示例**:
```json
{
  "code": 0,
  "message": "检索成功",
  "data": {
    "memories": [],
    "total": 0
  }
}
```

---

## 4. 记忆管理API

### 4.1 压缩对话记忆

**接口地址**: `POST /api/v1/memory/conversation/{conversation_id}/compress`

**请求参数**:
```json
{
  "strategy": "llm_summary",
  "target_length": 200,
  "keep_key_facts": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "记忆压缩成功",
  "data": {
    "original_count": 50,
    "compressed_count": 1,
    "summary": "用户是HR经理，讨论了招聘流程和候选人筛选标准...",
    "compression_ratio": 0.98,
    "key_facts": [
      "用户负责技术岗位招聘",
      "关注候选人技术能力"
    ]
  }
}
```

### 4.2 更新记忆

**接口地址**: `PUT /api/v1/memory/{memory_id}`

**请求参数**:
```json
{
  "content": "更新后的内容",
  "confidence": 0.95,
  "metadata": {}
}
```

### 4.3 删除记忆

**接口地址**: `DELETE /api/v1/memory/{memory_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "记忆删除成功",
  "data": {
    "deleted_memory_id": "mem-001"
  }
}
```

### 4.4 批量删除记忆

**接口地址**: `DELETE /api/v1/memory/batch`

**请求参数**:
```json
{
  "memory_ids": ["mem-001", "mem-002"],
  "reason": "用户请求删除"
}
```

### 4.5 设置记忆TTL

**接口地址**: `PUT /api/v1/memory/{memory_id}/ttl`

**请求参数**:
```json
{
  "ttl_seconds": 2592000,
  "reason": "30天后过期"
}
```

---

## 5. 记忆分析API

### 5.1 记忆统计

**接口地址**: `GET /api/v1/memory/stats`

**查询参数**:
- user_id: 用户ID
- memory_type: 记忆类型

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total_memories": 150,
    "by_type": {
      "conversation": 100,
      "preference": 30,
      "knowledge": 20
    },
    "storage_used_mb": 2.5,
    "oldest_memory": "2025-01-01T00:00:00Z",
    "newest_memory": "2025-01-03T11:00:00Z"
  }
}
```

### 5.2 记忆热度分析

**接口地址**: `GET /api/v1/memory/user/{user_id}/hot-memories`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "hot_memories": [
      {
        "memory_id": "know-001",
        "type": "knowledge",
        "content": "用户是HR经理",
        "access_count": 25,
        "last_accessed": "2025-01-03T10:55:00Z"
      }
    ]
  }
}
```

### 5.3 记忆质量评估

**接口地址**: `POST /api/v1/memory/quality/evaluate`

**请求参数**:
```json
{
  "memory_ids": ["know-001", "pref-001"],
  "evaluation_dimensions": [
    "relevance",
    "accuracy",
    "completeness"
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "评估完成",
  "data": {
    "results": [
      {
        "memory_id": "know-001",
        "scores": {
          "relevance": 0.92,
          "accuracy": 0.88,
          "completeness": 0.75
        },
        "overall_score": 0.85
      }
    ]
  }
}
```

---

## 6. 数据模型

### 6.1 Memory

```typescript
interface Memory {
  id: string;
  vector_id: string;
  user_id: string;
  conversation_id?: string;
  type: 'conversation' | 'preference' | 'knowledge';
  content: string;
  metadata?: Record<string, any>;
  confidence?: number;
  created_at: Date;
  updated_at?: Date;
  expires_at?: Date;
}
```

### 6.2 ConversationMemory

```typescript
interface ConversationMemory extends Memory {
  type: 'conversation';
  role: 'user' | 'assistant' | 'system';
  message_type?: 'text' | 'image' | 'file';
}
```

### 6.3 PreferenceMemory

```typescript
interface PreferenceMemory extends Memory {
  type: 'preference';
  preference_type: string;
  preference_value: any;
  source: 'explicit' | 'implicit' | 'inferred';
}
```

### 6.4 KnowledgeMemory

```typescript
interface KnowledgeMemory extends Memory {
  type: 'knowledge';
  knowledge_type: string;
  entities?: Record<string, any>;
  source_conversation_id?: string;
}
```

### 6.5 MemorySearchRequest

```typescript
interface MemorySearchRequest {
  query: string;
  user_id: string;
  memory_types?: MemoryType[];
  top_k?: number;
  score_threshold?: number;
  time_range?: {
    start: Date;
    end: Date;
  };
}
```

### 6.6 MemorySearchResult

```typescript
interface MemorySearchResult {
  memories: Array<{
    memory_id: string;
    type: MemoryType;
    content: string;
    score: number;
    metadata?: Record<string, any>;
  }>;
  total: number;
}
```

---

## 7. 后端实现示例

### 7.1 MemoryEngine 核心实现

```go
package memory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type MemoryEngine struct {
	memoryRepo   MemoryRepository
	vectorStore  VectorStore
	llmClient    LLMClient
}

// StoreMemory 存储记忆
func (e *MemoryEngine) StoreMemory(
	ctx context.Context,
	req *StoreMemoryRequest,
) (*Memory, error) {
	// 1. 生成向量
	vector, err := e.llmClient.Embed(ctx, req.Content)
	if err != nil {
		return nil, err
	}

	// 2. 存储到向量数据库
	vectorID, err := e.vectorStore.Insert(ctx, vector, req.Metadata)
	if err != nil {
		return nil, err
	}

	// 3. 存储记忆元数据
	memory := &Memory{
		ID:            generateMemoryID(),
		VectorID:      vectorID,
		UserID:        req.UserID,
		ConversationID: req.ConversationID,
		Type:          req.Type,
		Content:       req.Content,
		Metadata:      req.Metadata,
		Confidence:    req.Confidence,
		CreatedAt:     time.Now(),
	}

	err = e.memoryRepo.Create(ctx, memory)
	if err != nil {
		return nil, err
	}

	return memory, nil
}

// RetrieveMemories 语义检索记忆
func (e *MemoryEngine) RetrieveMemories(
	ctx context.Context,
	query string,
	userID string,
	memoryTypes []MemoryType,
	topK int,
) ([]*MemoryWithScore, error) {
	// 1. 向量化查询
	vector, err := e.llmClient.Embed(ctx, query)
	if err != nil {
		return nil, err
	}

	// 2. 向量检索
	vectorIDs, scores, err := e.vectorStore.Search(
		ctx,
		vector,
		topK,
		0.7, // score_threshold
	)
	if err != nil {
		return nil, err
	}

	// 3. 加载记忆内容
	memories, err := e.memoryRepo.GetByVectorIDs(ctx, vectorIDs)
	if err != nil {
		return nil, err
	}

	// 4. 组合结果
	result := make([]*MemoryWithScore, len(memories))
	for i, mem := range memories {
		result[i] = &MemoryWithScore{
			Memory: mem,
			Score:  scores[i],
		}
	}

	return result, nil
}

// SummarizeConversation 压缩对话记忆
func (e *MemoryEngine) SummarizeConversation(
	ctx context.Context,
	conversationID string,
	targetLength int,
) (*SummaryResult, error) {
	// 1. 获取对话历史
	memories, err := e.memoryRepo.GetByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// 2. 构建对话文本
	conversation := buildConversationText(memories)

	// 3. LLM摘要
	summary, err := e.llmClient.Summarize(ctx, conversation, targetLength)
	if err != nil {
		return nil, err
	}

	// 4. 提取关键事实
	keyFacts, err := e.llmClient.ExtractKeyFacts(ctx, conversation)
	if err != nil {
		keyFacts = []string{}
	}

	// 5. 存储摘要
	summaryMemory := &Memory{
		ID:            generateMemoryID(),
		UserID:        memories[0].UserID,
		ConversationID: conversationID,
		Type:          MemoryTypeConversation,
		Content:       summary.Text,
		Metadata: map[string]interface{}{
			"is_summary":    true,
			"key_facts":     keyFacts,
			"original_count": len(memories),
		},
		CreatedAt: time.Now(),
	}

	err = e.memoryRepo.Create(ctx, summaryMemory)
	if err != nil {
		return nil, err
	}

	// 6. 删除原始记忆
	err = e.memoryRepo.DeleteByConversationID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	return &SummaryResult{
		OriginalCount:    len(memories),
		CompressedCount:  1,
		Summary:         summary.Text,
		KeyFacts:        keyFacts,
		CompressionRatio: 1.0 / float64(len(memories)),
	}, nil
}
```

### 7.2 HTTP Handler

```go
// StoreMemoryHandler 存储记忆
func StoreMemoryHandler(ctx context.Context, c *app.RequestContext) {
	var req StoreMemoryRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	memory, err := memoryEngine.StoreMemory(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "存储记忆失败",
		})
		return
	}

	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    0,
		Message: "记忆存储成功",
		Data:    memory,
	})
}

// SearchMemoriesHandler 语义检索记忆
func SearchMemoriesHandler(ctx context.Context, c *app.RequestContext) {
	var req MemorySearchRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, ErrorResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	memories, err := memoryEngine.RetrieveMemories(
		ctx,
		req.Query,
		req.UserID,
		req.MemoryTypes,
		req.TopK,
	)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, ErrorResponse{
			Code:    500,
			Message: "检索记忆失败",
		})
		return
	}

	c.JSON(consts.StatusOK, SuccessResponse{
		Code:    0,
		Message: "检索成功",
		Data: MemorySearchResult{
			Memories: memories,
			Total:    len(memories),
		},
	})
}
```

### 7.3 数据库模型

```go
type Memory struct {
	ID             string    `gorm:"primaryKey;type:varchar(64)"`
	VectorID       string    `gorm:"index;type:varchar(64)"`
	UserID         string    `gorm:"index;type:varchar(64)"`
	ConversationID *string   `gorm:"index;type:varchar(64)"`
	Type           string    `gorm:"index;type:varchar(32)"`
	Content        string    `gorm:"type:text"`
	Metadata       *string   `gorm:"type:json"`
	Confidence     *float64  `gorm:"type:decimal(3,2)"`
	CreatedAt      time.Time `gorm:"index"`
	UpdatedAt      *time.Time
	ExpiresAt      *time.Time `gorm:"index"`
}

func (Memory) TableName() string {
	return "conversation_memories"
}
```

---

## 8. 前端调用示例

### 8.1 存储对话记忆

```typescript
import { request } from '@/utils/request';

export async function storeConversationMemory(params: {
  conversationId: string;
  role: 'user' | 'assistant';
  content: string;
  metadata?: Record<string, any>;
}) {
  return request.post('/api/v1/memory/conversation', {
    conversation_id: params.conversationId,
    role: params.role,
    content: params.content,
    metadata: {
      ...params.metadata,
      timestamp: new Date().toISOString(),
    },
  });
}

// 使用示例
await storeConversationMemory({
  conversationId: 'conv-001',
  role: 'user',
  content: '我是HR经理，负责招聘工作',
});
```

### 8.2 语义检索记忆

```typescript
export async function searchMemories(params: {
  query: string;
  userId: string;
  memoryTypes?: Array<'conversation' | 'preference' | 'knowledge'>;
  topK?: number;
  scoreThreshold?: number;
}) {
  return request.post<MemorySearchResult>('/api/v1/memory/search/semantic', {
    query: params.query,
    user_id: params.userId,
    memory_types: params.memoryTypes,
    top_k: params.topK || 5,
    score_threshold: params.scoreThreshold || 0.7,
  });
}

// 使用示例
const result = await searchMemories({
  query: '用户的工作是什么',
  userId: 'user-001',
  memoryTypes: ['knowledge', 'preference'],
  topK: 5,
});

console.log('检索到的记忆:', result.data.memories);
```

### 8.3 获取对话历史

```typescript
export async function getConversationHistory(params: {
  conversationId: string;
  limit?: number;
  offset?: number;
}) {
  return request.get<{
    conversation_id: string;
    messages: Array<{
      memory_id: string;
      role: string;
      content: string;
      timestamp: string;
    }>;
    total: number;
  }>(`/api/v1/memory/conversation/${params.conversationId}`, {
    params: {
      limit: params.limit || 10,
      offset: params.offset || 0,
    },
  });
}

// 使用示例
const history = await getConversationHistory({
  conversationId: 'conv-001',
  limit: 20,
});

console.log('对话历史:', history.data.messages);
```

### 8.4 React Hook 示例

```typescript
import { useState, useEffect } from 'react';
import { searchMemories, getConversationHistory } from '@/api/memory';

export function useMemorySearch(userId: string) {
  const [memories, setMemories] = useState<MemoryWithScore[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const search = async (query: string) => {
    setLoading(true);
    setError(null);

    try {
      const result = await searchMemories({
        query,
        userId,
        topK: 5,
      });
      setMemories(result.data.memories);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  };

  return { memories, loading, error, search };
}

// 使用示例
function MemorySearchComponent() {
  const { memories, loading, search } = useMemorySearch('user-001');

  return (
    <div>
      <input
        type="text"
        placeholder="搜索记忆..."
        onChange={(e) => search(e.target.value)}
      />
      {loading && <div>搜索中...</div>}
      <ul>
        {memories.map((mem) => (
          <li key={mem.memory_id}>
            [{mem.type}] {mem.content} (相似度: {mem.score.toFixed(2)})
          </li>
        ))}
      </ul>
    </div>
  );
}
```

---

## 9. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 80001 | 404 | 记忆不存在 |
| 80002 | 400 | 记忆类型无效 |
| 80003 | 400 | 向量化失败 |
| 80004 | 500 | 向量存储失败 |
| 80005 | 400 | 检索参数无效 |
| 80006 | 500 | LLM摘要失败 |
| 80007 | 400 | 记忆已过期 |
| 80008 | 403 | 无权访问此记忆 |
| 80101 | 400 | 压缩策略无效 |
| 80102 | 500 | 批量操作部分失败 |

---

**文档结束**
