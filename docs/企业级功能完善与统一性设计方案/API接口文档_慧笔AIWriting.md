# API接口文档：慧笔AIWriting模块

**模块名称**: 慧笔AIWriting (AIWriting)
**设计文档**: 03-ZKER前台_慧笔_AIWriting.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 智能写作API](#2-智能写作api)
- [3. 内容优化API](#3-内容优化api)
- [4. 模板管理API](#4-模板管理api)
- [5. 风格管理API](#5-风格管理api)
- [6. 历史记录API](#6-历史记录api)
- [7. 内容审核API](#7-内容审核api)
- [8. 数据模型](#8-数据模型)
- [9. 错误码定义](#9-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

慧笔AIWriting是ZKER企业级SaaS平台的智能内容创作模块：

- ✅ **智能写作** - 基于模板的AI内容生成，支持100+场景
- ✅ **内容优化** - 扩写、缩写、改写、润色
- ✅ **模板库** - 公文、营销、邮件、报告、创意等50+预置模板
- ✅ **风格管理** - 正式、活泼、专业、亲切等多种写作风格
- ✅ **批量生成** - 一次生成多版本供选择
- ✅ **历史记录** - 保存、复用、收藏历史作品

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- AI模型: 通义千问/GPT-4
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0

**前端技术栈**:
- 框架: React 18 + TypeScript
- 编辑器: Monaco Editor
- UI库: Semi Design

### 1.3 核心特性

| 特性 | 说明 |
|------|------|
| **配置优先** | 80%功能通过数据库配置实现 |
| **场景驱动** | 支持100+写作场景，可扩展 |
| **企业级** | 多租户隔离、权限控制、审计日志 |
| **低代码** | 企业可自定义模板，无需编码 |

---

## 2. 智能写作API

### 2.1 智能写作（流式）

**接口地址**: `POST /api/v1/writing/generate`

**功能说明**: 基于模板生成内容，支持流式输出

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "templateId": "tpl-business-email-001",
  "variables": {
    "recipient": "张总",
    "purpose": "inquiry",
    "keyPoints": "1. 询问产品价格\n2. 了解交付周期\n3. 索要产品目录",
    "tone": "formal"
  },
  "styleId": "style-professional",
  "options": {
    "stream": true,
    "maxTokens": 2000,
    "temperature": 0.7
  }
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| templateId | string | 是 | 模板ID |
| variables | object | 是 | 模板变量值 |
| styleId | string | 否 | 风格ID |
| options | object | 否 | 生成选项 |

**响应类型**: `text/event-stream` (Server-Sent Events)

**事件流格式**:
```
event: start
data: {"type":"start","id":"gen-20250103-123456","timestamp":"2025-01-03T10:00:00Z"}

event: content
data: {"type":"content","delta":"您好张总，"}

event: content
data: {"type":"content","delta":"我是XX公司的李明。"}

event: content
data: {"type":"content","delta":"特此致信，"}

event: end
data: {"type":"end","id":"gen-20250103-123456","wordCount":387,"tokens":512,"qualityScore":0.92}
```

### 2.2 批量生成

**接口地址**: `POST /api/v1/writing/batch-generate`

**功能说明**: 一次生成多个版本供选择

**请求参数**:
```json
{
  "templateId": "tpl-business-email-001",
  "variables": {
    "recipient": "张总",
    "purpose": "inquiry",
    "keyPoints": "询问产品价格",
    "tone": "formal"
  },
  "count": 3,
  "options": {
    "maxTokens": 2000,
    "temperature": 0.7
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "taskId": "batch-task-001",
    "status": "processing",
    "versions": []
  }
}
```

**获取批量生成结果**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "taskId": "batch-task-001",
    "status": "completed",
    "versions": [
      {
        "id": "ver-001",
        "content": "您好张总，我是XX公司的李明...",
        "wordCount": 387,
        "qualityScore": 0.92
      },
      {
        "id": "ver-002",
        "content": "尊敬的张总，您好！",
        "wordCount": 412,
        "qualityScore": 0.88
      },
      {
        "id": "ver-003",
        "content": "张总，您好！",
        "wordCount": 365,
        "qualityScore": 0.90
      }
    ]
  }
}
```

---

## 3. 内容优化API

### 3.1 内容优化

**接口地址**: `POST /api/v1/writing/optimize`

**请求参数**:
```json
{
  "content": "原始内容...",
  "operation": "expand",
  "instructions": "请将这段话扩展到500字，重点补充数据支撑"
}
```

**operation枚举**:
- `expand`: 扩写
- `summarize`: 缩写
- `rewrite`: 改写
- `polish`: 润色

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "optimizedContent": "优化后的内容...",
    "originalWordCount": 150,
    "optimizedWordCount": 487,
    "improvements": [
      "增加了数据支撑",
      "优化了语言表达",
      "补充了结构框架"
    ],
    "qualityScore": 0.91,
    "tokens": 650
  }
}
```

### 3.2 继续写作

**接口地址**: `POST /api/v1/writing/continue`

**功能说明**: 基于现有内容继续写作

**请求参数**:
```json
{
  "content": "现有内容...",
  "instructions": "请继续补充产品功能介绍部分"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "continuedContent": "补充的内容...",
    "originalWordCount": 300,
    "newWordCount": 450,
    "totalWordCount": 750
  }
}
```

---

## 4. 模板管理API

### 4.1 获取模板列表

**接口地址**: `GET /api/v1/writing/templates`

**查询参数**:
- page: 页码
- page_size: 每页数量
- category: 分类过滤（公文/营销/邮件/报告/创意）
- keyword: 搜索关键词
- is_public: 是否公开模板

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "items": [
      {
        "id": "tpl-business-email-001",
        "category": "邮件",
        "name": "商务咨询邮件",
        "description": "用于商务咨询场景的正式邮件模板",
        "icon": "https://cdn.example.com/icons/email.png",
        "variables": [
          {
            "name": "recipient",
            "type": "text",
            "label": "收件人",
            "required": true,
            "placeholder": "请输入收件人姓名"
          },
          {
            "name": "purpose",
            "type": "select",
            "label": "邮件目的",
            "required": true,
            "options": [
              {"value": "inquiry", "label": "咨询"},
              {"value": "complaint", "label": "投诉"}
            ]
          }
        ],
        "styleId": "style-professional",
        "usageCount": 1250,
        "favoriteCount": 85,
        "avgRating": 4.5,
        "isPublic": true,
        "createdBy": 1001,
        "createdAt": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 4.2 获取模板详情

**接口地址**: `GET /api/v1/writing/templates/{template_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "tpl-business-email-001",
    "category": "邮件",
    "name": "商务咨询邮件",
    "description": "用于商务咨询场景的正式邮件模板",
    "systemPrompt": "你是一位专业的商务邮件写作助手...",
    "variables": [
      {
        "name": "recipient",
        "type": "text",
        "label": "收件人",
        "required": true,
        "placeholder": "请输入收件人姓名",
        "maxLength": 50
      }
    ],
    "exampleInput": {
      "recipient": "张总",
      "purpose": "inquiry"
    },
    "exampleOutput": "您好张总，我是XX公司的李明...",
    "styleId": "style-professional",
    "usageCount": 1250,
    "avgRating": 4.5
  }
}
```

### 4.3 创建模板（管理员）

**接口地址**: `POST /api/v1/writing/templates`

**请求参数**:
```json
{
  "category": "营销",
  "name": "产品推广文案",
  "description": "用于产品推广的营销文案模板",
  "icon": "https://cdn.example.com/icons/marketing.png",
  "systemPrompt": "你是一位专业的营销文案写作助手...",
  "variables": [
    {
      "name": "productName",
      "type": "text",
      "label": "产品名称",
      "required": true
    }
  ],
  "exampleInput": {
    "productName": "XX智能助手"
  },
  "exampleOutput": "隆重介绍XX智能助手...",
  "styleId": "style-creative"
}
```

### 4.4 更新模板（管理员）

**接口地址**: `PUT /api/v1/writing/templates/{template_id}`

### 4.5 删除模板（管理员）

**接口地址**: `DELETE /api/v1/writing/templates/{template_id}`

### 4.6 发布/下架模板

**接口地址**: `POST /api/v1/writing/templates/{template_id}/publish`

**请求参数**:
```json
{
  "action": "publish"
}
```

**action枚举**: publish（发布）/ unpublish（下架）

### 4.7 获取模板分类

**接口地址**: `GET /api/v1/writing/templates/categories`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "categories": [
      {
        "id": "official",
        "name": "公文",
        "description": "通知、报告、总结等正式公文",
        "icon": "document",
        "templateCount": 25
      },
      {
        "id": "marketing",
        "name": "营销",
        "description": "朋友圈、小红书、公众号营销文案",
        "icon": "megaphone",
        "templateCount": 30
      }
    ]
  }
}
```

---

## 5. 风格管理API

### 5.1 获取风格列表

**接口地址**: `GET /api/v1/writing/styles`

**查询参数**:
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": "style-professional",
        "name": "专业正式",
        "description": "专业、正式、严谨",
        "config": {
          "tone": "professional",
          "toneDescription": "专业、正式、严谨",
          "vocabulary": {
            "level": "advanced",
            "preferProfessional": true,
            "avoidSlang": true
          },
          "sentence": {
            "avgLength": "20-30",
            "structure": "varied",
            "avoidRepetition": true
          },
          "punctuation": {
            "formal": true
          },
          "emoji": {
            "allow": false
          }
        },
        "usageCount": 3500
      }
    ]
  }
}
```

### 5.2 获取风格详情

**接口地址**: `GET /api/v1/writing/styles/{style_id}`

### 5.3 创建风格（管理员）

**接口地址**: `POST /api/v1/writing/styles`

**请求参数**:
```json
{
  "name": "亲切活泼",
  "description": "亲切、活泼、自然",
  "config": {
    "tone": "friendly",
    "toneDescription": "亲切、活泼、自然",
    "vocabulary": {
      "level": "intermediate",
      "preferProfessional": false,
      "avoidSlang": false
    },
    "sentence": {
      "avgLength": "15-25",
      "structure": "simple",
      "avoidRepetition": false
    },
    "punctuation": {
      "formal": false
    },
    "emoji": {
      "allow": true
    }
  }
}
```

### 5.4 更新风格（管理员）

**接口地址**: `PUT /api/v1/writing/styles/{style_id}`

### 5.5 删除风格（管理员）

**接口地址**: `DELETE /api/v1/writing/styles/{style_id}`

---

## 6. 历史记录API

### 6.1 获取历史记录列表

**接口地址**: `GET /api/v1/writing/histories`

**查询参数**:
- page: 页码
- page_size: 每页数量
- template_id: 模板过滤
- status: 状态过滤（draft/published/archived）
- keyword: 搜索关键词
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 85,
    "items": [
      {
        "id": "hist-001",
        "title": "商务咨询邮件",
        "templateId": "tpl-business-email-001",
        "templateName": "商务咨询邮件",
        "styleId": "style-professional",
        "generatedContent": "您好张总，我是XX公司的李明...",
        "wordCount": 387,
        "qualityScore": 0.92,
        "status": "published",
        "isFavorited": true,
        "createdAt": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 6.2 获取历史记录详情

**接口地址**: `GET /api/v1/writing/histories/{history_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "hist-001",
    "title": "商务咨询邮件",
    "templateId": "tpl-business-email-001",
    "styleId": "style-professional",
    "inputVariables": {
      "recipient": "张总",
      "purpose": "inquiry",
      "keyPoints": "询问产品价格"
    },
    "generatedContent": "完整内容...",
    "optimizedContent": "优化后的内容...",
    "wordCount": 387,
    "qualityScore": 0.92,
    "generationTime": 1250,
    "llmModel": "gpt-4",
    "llmTokens": 512,
    "status": "published",
    "isFavorited": true,
    "createdAt": "2025-01-03T10:00:00Z",
    "updatedAt": "2025-01-03T11:00:00Z"
  }
}
```

### 6.3 收藏/取消收藏

**接口地址**: `POST /api/v1/writing/histories/{history_id}/favorite`

**请求参数**:
```json
{
  "action": "favorite"
}
```

**action枚举**: favorite（收藏）/ unfavorite（取消收藏）

### 6.4 删除历史记录

**接口地址**: `DELETE /api/v1/writing/histories/{history_id}`

### 6.5 基于历史二次创作

**接口地址**: `POST /api/v1/writing/histories/{history_id}/recreate`

**请求参数**:
```json
{
  "templateId": "tpl-business-email-001",
  "variables": {
    "recipient": "李总",
    "purpose": "invitation"
  }
}
```

---

## 7. 内容审核API

### 7.1 敏感词检测

**接口地址**: `POST /api/v1/writing/audit/sensitive-words`

**请求参数**:
```json
{
  "content": "待检测的内容..."
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "hasViolation": true,
    "violations": [
      {
        "word": "敏感词",
        "category": "politics",
        "severity": "high",
        "position": {
          "start": 10,
          "end": 14
        }
      }
    ],
    "suggestion": "请修改内容中的敏感词汇"
  }
}
```

### 7.2 内容质量检测

**接口地址**: `POST /api/v1/writing/audit/quality`

**请求参数**:
```json
{
  "content": "待检测的内容..."
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "qualityScore": 0.85,
    "metrics": {
      "coherence": 0.88,
      "clarity": 0.82,
      "grammar": 0.90,
      "vocabulary": 0.80
    },
    "suggestions": [
      "建议补充数据支撑",
      "部分语句过长，建议分段"
    ]
  }
}
```

---

## 8. 数据模型

### 8.1 Template（模板）
```typescript
interface Template {
  id: string;
  tenant_id: string;
  category: 'official' | 'marketing' | 'email' | 'report' | 'creative';
  name: string;
  description?: string;
  icon?: string;
  system_prompt: string;
  variables: TemplateVariable[];
  example_input?: Record<string, any>;
  example_output?: string;
  style_id?: string;
  is_public: boolean;
  is_active: boolean;
  version: number;
  usage_count: number;
  favorite_count: number;
  avg_rating: number;
  created_by?: number;
  created_at: string;
  updated_at: string;
}
```

### 8.2 TemplateVariable（模板变量）
```typescript
interface TemplateVariable {
  name: string;
  type: 'text' | 'textarea' | 'select' | 'radio' | 'checkbox' | 'date' | 'number';
  label: string;
  required: boolean;
  placeholder?: string;
  default_value?: string;
  options?: Array<{
    value: string;
    label: string;
  }>;
  max_length?: number;
  validation_rules?: any;
}
```

### 8.3 Style（风格）
```typescript
interface Style {
  id: string;
  tenant_id?: string;
  name: string;
  description?: string;
  config: {
    tone: string;
    tone_description: string;
    vocabulary: {
      level: 'basic' | 'intermediate' | 'advanced';
      preferProfessional: boolean;
      avoidSlang: boolean;
    };
    sentence: {
      avgLength: string;
      structure: 'simple' | 'varied' | 'complex';
      avoidRepetition: boolean;
    };
    punctuation: {
      formal: boolean;
    };
    emoji: {
      allow: boolean;
    };
    system_prompt_addition?: string;
  };
  usage_count: number;
  created_at: string;
  updated_at: string;
}
```

### 8.4 WritingHistory（历史记录）
```typescript
interface WritingHistory {
  id: string;
  tenant_id: string;
  user_id: number;
  template_id?: string;
  style_id?: string;
  title?: string;
  input_variables: Record<string, any>;
  generated_content: string;
  optimized_content?: string;
  word_count?: number;
  quality_score?: number;
  generation_time?: number;
  llm_model?: string;
  llm_tokens?: number;
  status: 'draft' | 'published' | 'archived';
  is_favorited: boolean;
  created_at: string;
  updated_at: string;
}
```

---

## 9. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 80001 | 400 | 模板不存在 |
| 80002 | 400 | 模板变量验证失败 |
| 80003 | 400 | 模板已下架 |
| 80101 | 400 | 风格不存在 |
| 80102 | 400 | 风格配置无效 |
| 80201 | 400 | 写作生成失败 |
| 80202 | 400 | LLM调用失败 |
| 80203 | 400 | 内容审核未通过 |
| 80204 | 429 | 调用频率超限 |
| 80301 | 404 | 历史记录不存在 |
| 80302 | 403 | 无权限访问该历史记录 |
| 80401 | 400 | 内容包含敏感词 |
| 80402 | 400 | 内容质量不达标 |

---

**文档结束**
