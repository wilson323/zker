# API接口文档：数字员工管理模块

**模块名称**: 数字员工管理 (Bot Management)
**设计文档**: 07-企业管理中台_数字员工管理.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. Bot管理API](#2-bot管理api)
- [3. Bot技能关联API](#3-bot技能关联api)
- [4. Bot知识库关联API](#4-bot知识库关联api)
- [5. Bot发布与版本API](#5-bot发布与版本api)
- [6. Bot分析API](#6-bot分析api)
- [7. Bot预览API](#7-bot预览api)
- [8. Bot模板API](#8-bot模板api)
- [9. 数据模型](#9-数据模型)
- [10. 错误码定义](#10-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

数字员工管理是ZKER企业AI智能体的统一管理平台，提供Bot的创建、配置、发布、监控和数据分析：

- ✅ **Bot管理** - 创建、编辑、删除、复制Bot
- ✅ **三种类型** - 问答型(Q&A)、操作型(Action)、综合型(Comprehensive)
- ✅ **技能配置** - 关联技能、排序、启用/禁用
- ✅ **知识库配置** - 关联知识库、设置优先级
- ✅ **发布管理** - 发布、下架、版本管理
- ✅ **数据分析** - 使用统计、效果分析、日志查询

### 1.2 Bot类型

| Bot类型 | 说明 | 必需配置 | 典型场景 |
|---------|------|---------|---------|
| **问答型** | 基于知识库回答问题 | 至少1个知识库 | FAQ、规章制度、产品知识 |
| **操作型** | 调用技能执行操作 | 至少1个技能 | 查询数据、发送通知、调用API |
| **综合型** | 同时具备问答和操作能力 | 至少1个知识库或技能 | 智能客服、助理Bot |

### 1.3 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0
- 分析: ClickHouse

**前端技术栈**:
- React 18 + TypeScript
- Semi Design UI
- React Flow (可视化配置)

---

## 2. Bot管理API

### 2.1 获取Bot列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots`

**查询参数**:
- type: Bot类型筛选 (qa/action/comprehensive)
- status: 状态筛选 (draft/published/archived)
- keyword: 搜索关键词(名称/描述)
- creator_id: 创建者ID筛选
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "items": [
      {
        "id": "bot-001",
        "tenant_id": "tenant-a",
        "creator_id": 1,
        "creator_name": "张三",
        "name": "企业FAQ助手",
        "description": "回答企业常见问题",
        "avatar": "https://cdn.example.com/avatar/bot-001.png",
        "type": "qa",
        "status": "published",
        "is_listed": true,
        "usage_count": 1523,
        "active_users": 89,
        "skill_count": 0,
        "knowledge_base_count": 3,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:30:00Z"
      }
    ]
  }
}
```

### 2.2 创建Bot

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots`

**请求参数**:
```json
{
  "name": "客服助手",
  "description": "智能客服助手,支持问答和操作",
  "avatar": "https://cdn.example.com/avatar/new-bot.png",
  "type": "comprehensive",
  "prompt": "你是一个专业的客服助手...",
  "welcome_message": "您好,我是您的专属客服助手",
  "conversation_style": {
    "tone": "friendly",
    "address_user": "您",
    "length": "moderate"
  },
  "temperature": 0.7,
  "max_tokens": 2000,
  "stream_output": true,
  "skills": [
    {
      "skill_id": "skill-001",
      "is_enabled": true,
      "sort_order": 1,
      "config": {}
    }
  ],
  "knowledge_bases": [
    {
      "knowledge_base_id": "kb-001",
      "priority": 8,
      "is_enabled": true
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "bot-002",
    "name": "客服助手",
    "type": "comprehensive",
    "status": "draft",
    "created_at": "2025-01-15T14:30:00Z"
  }
}
```

### 2.3 获取Bot详情

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "bot-001",
    "name": "企业FAQ助手",
    "type": "qa",
    "status": "published",
    "is_listed": true,
    "prompt": "你是一个专业的企业FAQ助手...",
    "welcome_message": "您好,我是企业FAQ助手",
    "conversation_style": {
      "tone": "formal",
      "address_user": "您",
      "length": "moderate"
    },
    "temperature": 0.7,
    "max_tokens": 2000,
    "usage_count": 1523,
    "active_users": 89,
    "skills": [],
    "knowledge_bases": [
      {
        "knowledge_base_id": "kb-001",
        "knowledge_base_name": "企业规章制度",
        "priority": 8,
        "is_enabled": true
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

### 2.4 更新Bot

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/bots/{bot_id}`

**请求参数**:
```json
{
  "name": "企业FAQ助手v2",
  "description": "回答企业常见问题(更新版)",
  "prompt": "你是一个专业的企业FAQ助手...",
  "conversation_style": {
    "tone": "friendly"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": "bot-001",
    "name": "企业FAQ助手v2",
    "updated_at": "2025-01-15T15:00:00Z"
  }
}
```

### 2.5 删除Bot

**接口地址**: `DELETE /api/v1/tenants/{tenant_id}/bots/{bot_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_id": "bot-001"
  }
}
```

### 2.6 复制Bot

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/copy`

**请求参数**:
```json
{
  "name": "企业FAQ助手副本",
  "copy_skills": true,
  "copy_knowledge_bases": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "复制成功",
  "data": {
    "id": "bot-003",
    "name": "企业FAQ助手副本",
    "created_at": "2025-01-15T16:00:00Z"
  }
}
```

### 2.7 发布Bot

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/publish`

**请求参数**:
```json
{
  "create_version": true,
  "version_description": "首次发布"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "发布成功",
  "data": {
    "id": "bot-001",
    "status": "published",
    "published_at": "2025-01-15T17:00:00Z",
    "version": {
      "id": "version-001",
      "version_number": "v1.0.0"
    }
  }
}
```

### 2.8 下架Bot

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/unpublish`

**响应示例**:
```json
{
  "code": 0,
  "message": "下架成功",
  "data": {
    "id": "bot-001",
    "status": "draft"
  }
}
```

### 2.9 上架到发现页

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/list`

**响应示例**:
```json
{
  "code": 0,
  "message": "上架成功",
  "data": {
    "id": "bot-001",
    "is_listed": true
  }
}
```

### 2.10 从发现页撤下

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/unlist`

**响应示例**:
```json
{
  "code": 0,
  "message": "撤下成功",
  "data": {
    "id": "bot-001",
    "is_listed": false
  }
}
```

---

## 3. Bot技能关联API

### 3.1 添加技能

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/skills`

**请求参数**:
```json
{
  "skill_id": "skill-002",
  "is_enabled": true,
  "sort_order": 2,
  "config": {
    "timeout": 30
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "添加成功",
  "data": {
    "id": "bot-skill-002",
    "bot_id": "bot-001",
    "skill_id": "skill-002"
  }
}
```

### 3.2 移除技能

**接口地址**: `DELETE /api/v1/tenants/{tenant_id}/bots/{bot_id}/skills/{skill_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "移除成功"
}
```

### 3.3 更新技能配置

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/bots/{bot_id}/skills/{skill_id}`

**请求参数**:
```json
{
  "is_enabled": false,
  "sort_order": 3,
  "config": {
    "timeout": 60
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功"
}
```

---

## 4. Bot知识库关联API

### 4.1 关联知识库

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/knowledge-bases`

**请求参数**:
```json
{
  "knowledge_base_id": "kb-002",
  "priority": 7,
  "is_enabled": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "关联成功",
  "data": {
    "id": "bot-kb-002",
    "bot_id": "bot-001",
    "knowledge_base_id": "kb-002"
  }
}
```

### 4.2 移除知识库

**接口地址**: `DELETE /api/v1/tenants/{tenant_id}/bots/{bot_id}/knowledge-bases/{kb_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "移除成功"
}
```

---

## 5. Bot发布与版本API

### 5.1 创建Bot版本

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/versions`

**请求参数**:
```json
{
  "description": "优化了回复速度"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "版本创建成功",
  "data": {
    "id": "version-002",
    "bot_id": "bot-001",
    "version_number": "v1.0.1",
    "description": "优化了回复速度",
    "created_at": "2025-01-15T18:00:00Z"
  }
}
```

### 5.2 获取Bot版本列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/versions`

**响应示例**:
```json
{
  "code": 0,
  "message": "获取成功",
  "data": {
    "total": 2,
    "items": [
      {
        "id": "version-001",
        "version_number": "v1.0.0",
        "description": "首次发布",
        "created_by": "张三",
        "created_at": "2025-01-10T10:00:00Z"
      },
      {
        "id": "version-002",
        "version_number": "v1.0.1",
        "description": "优化了回复速度",
        "created_by": "李四",
        "created_at": "2025-01-15T18:00:00Z"
      }
    ]
  }
}
```

### 5.3 回滚到指定版本

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/versions/{version_id}/rollback`

**响应示例**:
```json
{
  "code": 0,
  "message": "回滚成功",
  "data": {
    "bot_id": "bot-001",
    "rolled_back_to": "v1.0.0",
    "new_version": "v1.0.2",
    "rolled_back_at": "2025-01-15T19:00:00Z"
  }
}
```

---

## 6. Bot分析API

### 6.1 获取Bot统计概览

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/analytics/overview`

**查询参数**:
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_conversations": 1523,
    "total_messages": 8456,
    "unique_users": 89,
    "avg_response_time": 1250,
    "avg_tokens_per_message": 156,
    "avg_rating": 4.5,
    "positive_rate": 0.92,
    "skill_calls": {
      "skill-001": 234,
      "skill-002": 156
    },
    "trend": [
      {
        "date": "2025-01-09",
        "conversations": 120,
        "messages": 654
      },
      {
        "date": "2025-01-10",
        "conversations": 145,
        "messages": 823
      }
    ]
  }
}
```

### 6.2 获取Bot使用日志

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/analytics/logs`

**查询参数**:
- start_date: 开始日期
- end_date: 结束日期
- user_id: 用户ID筛选 (可选)
- action: 操作类型筛选 (可选)
- page: 页码 (可选)
- page_size: 每页数量 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1523,
    "items": [
      {
        "id": 1,
        "bot_id": "bot-001",
        "user_id": 10,
        "user_name": "李四",
        "conversation_id": "conv-001",
        "action": "send_message",
        "details": {
          "message": "今天天气怎么样?",
          "response_time": 1250
        },
        "created_at": "2025-01-15T10:30:00Z"
      }
    ]
  }
}
```

---

## 7. Bot预览API

### 7.1 预览对话

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/preview`

**请求参数**:
```json
{
  "message": "你好",
  "conversation_id": "preview-conv-001",
  "user_id": 1
}
```

**响应示例** (流式):
```
Content-Type: text/event-stream

data: {"text": "你"}
data: {"text": "好"}
data: {"text": "！"}
data: {"text": "有"}
data: {"text": "什"}
data: {"text": "么"}
data: {"text": "可"}
data: {"text": "以"}
data: {"text": "帮"}
data: {"text": "助"}
data: {"text": "你"}
data: {"text": "的"}
data: {"text": "？"}
data: [DONE]
```

---

## 8. Bot模板API

### 8.1 获取模板列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bot-templates`

**查询参数**:
- category: 分类筛选 (可选)
- type: Bot类型筛选 (可选)
- source: 来源筛选 (official/custom) (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 20,
    "items": [
      {
        "id": "template-001",
        "tenant_id": null,
        "name": "客服助手模板",
        "description": "专业的客服助手模板",
        "category": "客服",
        "type": "comprehensive",
        "thumbnail": "https://cdn.example.com/template/tpl-001.png",
        "usage_count": 152,
        "is_official": true
      }
    ]
  }
}
```

### 8.2 使用模板创建Bot

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bot-templates/{template_id}/use`

**请求参数**:
```json
{
  "name": "我的客服助手",
  "customize": {
    "prompt": "自定义提示词..."
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "bot_id": "bot-003",
    "name": "我的客服助手"
  }
}
```

---

## 9. 数据模型

### 9.1 Bot

```typescript
interface Bot {
  id: string;
  tenant_id: string;
  creator_id: number;
  name: string;
  description?: string;
  avatar?: string;
  type: 'qa' | 'action' | 'comprehensive';
  prompt?: string;
  welcome_message?: string;
  conversation_style?: ConversationStyle;
  temperature?: number;
  max_tokens?: number;
  stream_output?: boolean;
  status: 'draft' | 'published' | 'archived';
  is_listed: boolean;
  published_at?: Date;
  usage_count: number;
  active_users: number;
  created_at: Date;
  updated_at: Date;
}

interface ConversationStyle {
  tone: 'formal' | 'casual' | 'friendly' | 'humorous';
  address_user: '您' | '你';
  length: 'concise' | 'moderate' | 'detailed';
  personality?: string;
}
```

### 9.2 BotSkill

```typescript
interface BotSkill {
  id: string;
  bot_id: string;
  skill_id: string;
  sort_order: number;
  is_enabled: boolean;
  config?: Record<string, any>;
  created_at: Date;
}
```

### 9.3 BotKnowledgeBase

```typescript
interface BotKnowledgeBase {
  id: string;
  bot_id: string;
  knowledge_base_id: string;
  priority: number; // 1-10
  is_enabled: boolean;
  created_at: Date;
}
```

### 9.4 BotVersion

```typescript
interface BotVersion {
  id: string;
  bot_id: string;
  version_number: string; // v1.0.0
  description?: string;
  config: BotConfig;
  created_by: string;
  created_at: Date;
}

interface BotConfig {
  basic_info: {
    name: string;
    description?: string;
    avatar?: string;
  };
  skills: Array<{
    skill_id: string;
    is_enabled: boolean;
    config?: Record<string, any>;
  }>;
  knowledge_bases: Array<{
    knowledge_base_id: string;
    priority: number;
    is_enabled: boolean;
  }>;
  prompt?: string;
  conversation_style?: ConversationStyle;
  temperature?: number;
  max_tokens?: number;
}
```

### 9.5 BotStatistics

```typescript
interface BotStatistics {
  total_conversations: number;
  total_messages: number;
  unique_users: number;
  avg_response_time: number;
  avg_tokens_per_message: number;
  avg_rating: number;
  positive_rate: number;
  skill_calls: Record<string, number>;
}
```

---

## 10. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 30001 | 404 | Bot不存在 |
| 30002 | 400 | Bot名称已存在 |
| 30003 | 400 | Bot类型无效 |
| 30004 | 400 | 问答型Bot必须关联至少1个知识库 |
| 30005 | 400 | 操作型Bot必须关联至少1个技能 |
| 30006 | 400 | 综合型Bot至少需要关联1个技能或知识库 |
| 30007 | 400 | Bot状态不允许此操作 |
| 30008 | 404 | 技能不存在 |
| 30009 | 404 | 知识库不存在 |
| 30010 | 400 | Bot已发布,无法修改 |
| 30101 | 404 | Bot版本不存在 |
| 30102 | 400 | 无法回滚到当前版本 |
| 30201 | 404 | Bot模板不存在 |

---

**文档结束**
