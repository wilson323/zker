# API接口文档：会话SuperChatbox模块

**模块名称**: 会话SuperChatbox (SuperChatbox)
**设计文档**: 01-ZKER前台_会话_SuperChatbox.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P1

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. WebSocket通信协议](#2-websocket通信协议)
- [3. 会话管理API](#3-会话管理api)
- [4. 消息管理API](#4-消息管理api)
- [5. 快捷指令API](#5-快捷指令api)
- [6. 数据模型](#6-数据模型)
- [7. 错误码定义](#7-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

会话SuperChatbox是ZKER企业级SaaS平台的核心前台模块，提供企业级智能对话UI和交互能力：

- ✅ **实时对话** - WebSocket双向通信，流式输出
- ✅ **多轮对话** - 会话上下文管理，历史记录保留
- ✅ **多模态输入** - 文本、图片、文件等多种输入方式
- ✅ **快捷指令** - 预置和自定义快捷指令
- ✅ **会话管理** - 创建、删除、归档、分享会话
- ✅ **消息操作** - 编辑、删除、重新生成、评价消息

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 实时通信: WebSocket
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0
- 文件存储: MinIO

**前端技术栈**:
- 框架: React 18 + TypeScript
- WebSocket库: WebSocket API
- UI库: Semi Design
- Markdown渲染: react-markdown

---

## 2. WebSocket通信协议

### 2.1 连接建立

**WebSocket URL**:
```
wss://api.example.com/v1/chat/ws?conversation_id={conv_id}&token={jwt_token}
```

### 2.2 客户端消息格式

#### 发送文本消息
```json
{
  "type": "message.send",
  "data": {
    "content": "你好，请帮我写一封商务邮件",
    "messageId": "msg-client-123"
  }
}
```

#### 停止生成
```json
{
  "type": "generation.stop",
  "data": {
    "messageId": "msg-assistant-456"
  }
}
```

#### 重新生成
```json
{
  "type": "message.regenerate",
  "data": {
    "parentMessageId": "msg-assistant-456"
  }
}
```

### 2.3 服务端消息格式

#### 消息已接收
```json
{
  "type": "message.received",
  "data": {
    "messageId": "msg-client-123",
    "tempId": "msg-server-789"
  }
}
```

#### 开始流式输出
```json
{
  "type": "generation.start",
  "data": {
    "messageId": "msg-assistant-456",
    "conversationId": "conv-001"
  }
}
```

#### 流式内容片段
```json
{
  "type": "generation.content",
  "data": {
    "messageId": "msg-assistant-456",
    "delta": "您好，",
    "isComplete": false
  }
}
```

#### 流式输出结束
```json
{
  "type": "generation.end",
  "data": {
    "messageId": "msg-assistant-456",
    "tokens": {
      "prompt": 150,
      "completion": 300,
      "total": 450
    },
    "latencyMs": 2340
  }
}
```

---

## 3. 会话管理API

### 3.1 创建会话

**接口地址**: `POST /api/v1/chat/conversations`

**功能说明**: 创建新会话

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "bot_id": "bot-123",
  "title": "邮件写作助手",
  "description": "帮助我撰写商务邮件",
  "model_config": {
    "temperature": 0.7,
    "max_tokens": 2000
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "会话创建成功",
  "data": {
    "id": "conv-001",
    "tenant_id": "tenant-123",
    "user_id": 1001,
    "bot_id": "bot-123",
    "title": "邮件写作助手",
    "status": "active",
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

### 3.2 查询会话列表

**接口地址**: `GET /api/v1/chat/conversations`

**功能说明**: 获取当前用户的会话列表

**查询参数**:
- page: 页码，默认1
- page_size: 每页数量，默认20
- status: 状态过滤
- bot_id: Bot ID过滤

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "items": [
      {
        "id": "conv-001",
        "title": "邮件写作助手",
        "message_count": 15,
        "last_message_at": "2025-01-03T09:30:00Z",
        "status": "active"
      }
    ]
  }
}
```

### 3.3 删除会话

**接口地址**: `DELETE /api/v1/chat/conversations/{conversation_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "会话删除成功"
}
```

---

## 4. 消息管理API

### 4.1 获取会话消息列表

**接口地址**: `GET /api/v1/chat/conversations/{conversation_id}/messages`

**查询参数**:
- page: 页码
- page_size: 每页数量
- limit: 限制返回数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": "msg-001",
        "role": "user",
        "content": "你好，请帮我写一封商务邮件",
        "created_at": "2025-01-03T10:00:00Z"
      },
      {
        "id": "msg-002",
        "role": "assistant",
        "content": "好的，请问这封邮件的收件人是谁？",
        "created_at": "2025-01-03T10:00:02Z"
      }
    ]
  }
}
```

### 4.2 删除消息

**接口地址**: `DELETE /api/v1/chat/messages/{message_id}`

### 4.3 重新生成

**接口地址**: `POST /api/v1/chat/messages/{message_id}/regenerate`

---

## 5. 快捷指令API

### 5.1 获取快捷指令列表

**接口地址**: `GET /api/v1/chat/quick-commands`

**查询参数**:
- bot_id: Bot ID过滤
- category: 分类过滤

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "写邮件",
        "description": "快速撰写商务邮件",
        "prompt": "请帮我写一封商务邮件",
        "icon": "📧"
      }
    ]
  }
}
```

### 5.2 创建自定义快捷指令

**接口地址**: `POST /api/v1/chat/quick-commands`

**请求参数**:
```json
{
  "title": "写代码",
  "prompt": "请帮我写一段代码",
  "icon": "💻"
}
```

---

## 6. 数据模型

### 6.1 Conversation（会话）
```typescript
interface Conversation {
  id: string;
  tenant_id: string;
  user_id: number;
  bot_id?: string;
  title?: string;
  message_count: number;
  status: 'active' | 'archived' | 'deleted';
  created_at: string;
}
```

### 6.2 Message（消息）
```typescript
interface Message {
  id: string;
  conversation_id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  metadata?: {
    tokens?: { prompt: number; completion: number; total: number };
    latency_ms?: number;
  };
  created_at: string;
}
```

---

## 7. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 30001 | 400 | 会话不存在 |
| 30101 | 400 | 消息不存在 |
| 30107 | 429 | 发送消息过于频繁 |
| 30201 | 400 | 快捷指令不存在 |

---

**文档结束**
