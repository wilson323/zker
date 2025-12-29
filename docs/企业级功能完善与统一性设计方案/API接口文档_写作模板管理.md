# API接口文档：写作模板管理模块

**模块名称**: 写作模板管理 (TemplateManagement)
**设计文档**: 10-企业管理中台_写作模板管理.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 模板管理API](#2-模板管理api)
- [3. 模板分类API](#3-模板分类api)
- [4. 数据模型](#4-数据模型)
- [5. 错误码定义](#5-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

写作模板管理是ZKER企业级SaaS平台的企业管理中台模块，通过**100%复用通用CRUD框架**的方式，提供写作模板管理能力：

- ✅ **模板管理** - 创建、编辑、删除、发布/下架模板
- ✅ **模板分类** - 自定义企业写作模板分类
- ✅ **零编码** - 完全复用CRUD框架，仅需数据库配置
- ✅ **企业定制** - 企业可自定义写作模板

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5 (复用AIWriting模块的writing_templates表)
- 实现: 通用CRUD框架 (0%编码)

**实现策略**: ✅ 0% 编码 + 100% 配置（复用CRUD框架）

### 1.3 表复用关系

| 表名 | 来源 | 用途 |
|------|------|------|
| `writing_templates` | 03-慧笔_AIWriting | 写作模板定义 |
| `template_categories` | 新增 | 模板分类 |

---

## 2. 模板管理API

### 2.1 获取模板列表

**接口地址**: `GET /api/v1/admin/templates`

**查询参数**:
- category: 分类过滤 (可选)
- keyword: 搜索关键词 (可选)
- is_public: 是否公开 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "id": "tpl-001",
        "tenant_id": "tenant-001",
        "category": "公文写作",
        "name": "公司内部通知",
        "description": "公司内部通知模板",
        "system_prompt": "请根据以下信息撰写一份公司内部通知：\\n\\n**标题**: {{title}}\\n**内容**: {{content}}",
        "variables": [
          {
            "name": "title",
            "type": "text",
            "label": "通知标题",
            "required": true
          },
          {
            "name": "content",
            "type": "textarea",
            "label": "通知内容",
            "required": true
          }
        ],
        "is_public": false,
        "usage_count": 89,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 2.2 创建模板

**接口地址**: `POST /api/v1/admin/templates`

**请求参数**:
```json
{
  "category": "公文写作",
  "name": "公司内部通知",
  "description": "公司内部通知模板",
  "system_prompt": "请根据以下信息撰写一份公司内部通知：\n\n**标题**: {{title}}\n**内容**: {{content}}",
  "variables": [
    {
      "name": "title",
      "type": "text",
      "label": "通知标题",
      "required": true
    },
    {
      "name": "content",
      "type": "textarea",
      "label": "通知内容",
      "required": true
    }
  ],
  "is_public": false
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "tpl-002",
    "name": "公司内部通知",
    "created_at": "2025-01-15T14:00:00Z"
  }
}
```

### 2.3 更新模板

**接口地址**: `PUT /api/v1/admin/templates/{template_id}`

**请求参数**:
```json
{
  "name": "公司内部通知(更新版)",
  "description": "公司内部通知模板",
  "system_prompt": "更新后的提示词",
  "variables": []
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": "tpl-001",
    "updated_at": "2025-01-15T15:00:00Z"
  }
}
```

### 2.4 删除模板

**接口地址**: `DELETE /api/v1/admin/templates/{template_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_template_id": "tpl-001"
  }
}
```

### 2.5 获取模板详情

**接口地址**: `GET /api/v1/admin/templates/{template_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "tpl-001",
    "tenant_id": "tenant-001",
    "category": "公文写作",
    "name": "公司内部通知",
    "description": "公司内部通知模板",
    "system_prompt": "请根据以下信息撰写一份公司内部通知...",
    "variables": [
      {
        "name": "title",
        "type": "text",
        "label": "通知标题",
        "required": true,
        "default": ""
      },
      {
        "name": "content",
        "type": "textarea",
        "label": "通知内容",
        "required": true,
        "default": ""
      }
    ],
    "is_public": false,
    "usage_count": 89,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

---

## 3. 模板分类API

### 3.1 获取分类列表

**接口地址**: `GET /api/v1/admin/template-categories`

**查询参数**:
- is_active: 是否启用 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 4,
    "items": [
      {
        "id": 1,
        "name": "公文写作",
        "description": "通知、报告、总结等公文模板",
        "sort_order": 1,
        "is_active": true,
        "template_count": 10
      },
      {
        "id": 2,
        "name": "营销文案",
        "description": "朋友圈、公众号等营销模板",
        "sort_order": 2,
        "is_active": true,
        "template_count": 8
      },
      {
        "id": 3,
        "name": "商务邮件",
        "description": "各类商务邮件模板",
        "sort_order": 3,
        "is_active": true,
        "template_count": 5
      },
      {
        "id": 4,
        "name": "报告总结",
        "description": "工作汇报、项目总结模板",
        "sort_order": 4,
        "is_active": true,
        "template_count": 2
      }
    ]
  }
}
```

### 3.2 创建分类

**接口地址**: `POST /api/v1/admin/template-categories`

**请求参数**:
```json
{
  "name": "技术文档",
  "description": "技术文档写作模板",
  "sort_order": 5
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 5,
    "name": "技术文档",
    "created_at": "2025-01-15T16:00:00Z"
  }
}
```

### 3.3 更新分类

**接口地址**: `PUT /api/v1/admin/template-categories/{category_id}`

**请求参数**:
```json
{
  "name": "技术文档(更新)",
  "description": "技术文档写作模板",
  "sort_order": 6,
  "is_active": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 5,
    "updated_at": "2025-01-15T17:00:00Z"
  }
}
```

### 3.4 删除分类

**接口地址**: `DELETE /api/v1/admin/template-categories/{category_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_category_id": 5
  }
}
```

---

## 4. 数据模型

### 4.1 WritingTemplate

```typescript
interface WritingTemplate {
  id: string;
  tenant_id: string;
  category: string;
  name: string;
  description?: string;
  system_prompt: string;
  variables: TemplateVariable[];
  is_public: boolean;
  usage_count: number;
  created_at: Date;
  updated_at: Date;
}

interface TemplateVariable {
  name: string;
  type: 'text' | 'textarea' | 'number' | 'date' | 'select';
  label: string;
  required: boolean;
  default?: string;
  options?: string[]; // for type='select'
}
```

### 4.2 TemplateCategory

```typescript
interface TemplateCategory {
  id: number;
  tenant_id: string; // NULL表示官方分类
  name: string;
  description?: string;
  sort_order: number;
  is_active: boolean;
  template_count?: number;
  created_at: Date;
  updated_at: Date;
}
```

---

## 5. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 32001 | 404 | 模板不存在 |
| 32002 | 400 | 模板名称已存在 |
| 32003 | 400 | 模板分类无效 |
| 32004 | 400 | 变量定义无效 |
| 32101 | 404 | 模板分类不存在 |
| 32102 | 400 | 分类名称已存在 |
| 32103 | 400 | 分类下存在模板,无法删除 |

---

**文档结束**
