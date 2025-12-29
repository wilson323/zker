# API 接口文档 - 待办TaskCenter

**模块编号**: 05
**模块名称**: 待办TaskCenter (TaskCenter)
**API 版本**: v1.0.0
**基础路径**: `/api/v1`
**协议**: HTTPS + WebSocket
**数据格式**: JSON

---

## 📋 目录

- [1. 模块概述](#1-模块概述)
- [2. API 接口列表](#2-api-接口列表)
  - [2.1 待办管理](#21-待办管理)
  - [2.2 任务类型管理](#22-任务类型管理)
  - [2.3 动态字段管理](#23-动态字段管理)
  - [2.4 分类管理](#24-分类管理)
- [3. 数据模型](#3-数据模型)
- [4. 错误码定义](#4-错误码定义)
- [5. 调用示例](#5-调用示例)

---

## 1. 模块概述

### 1.1 功能说明

**待办TaskCenter** 是企业级 SaaS 平台的待办事项管理模块，通过**通用任务引擎 + 任务类型配置**的方式，提供灵活的待办管理能力。

**核心特性**：
- ✅ **通用任务引擎**：20% 编码实现通用任务处理框架
- ✅ **配置驱动**：80% 任务类型通过数据库配置实现
- ✅ **动态字段**：支持根据任务类型动态渲染表单
- ✅ **灵活扩展**：企业可自定义任务类型和字段

### 1.2 技术架构

**后端技术栈**：
- **框架**：CloudWeGo Hertz（Go）
- **数据库**：MySQL 8.4.5（主数据）、Redis 8.0（缓存）
- **验证**：JSON Schema 验证动态字段

**前端技术栈**：
- **框架**：React 18 + TypeScript 5.6
- **UI 库**：Semi Design
- **表单**：动态表单组件（基于字段配置渲染）

**核心设计理念**：
```
通用任务引擎（20%编码）
├── 任务 CRUD 操作
├── 动态字段验证
├── 状态流转管理
└── 权限控制

任务类型配置（80%配置）
├── 预置任务类型（审批、报告、会议）
├── 企业自定义类型
├── 字段定义（JSON配置）
└── 验证规则（JSON配置）
```

### 1.3 性能指标

| 指标 | 目标值 |
|------|--------|
| API 响应时间 | P95 < 100ms |
| 并发处理能力 | 10000 QPS |
| 任务创建延迟 | < 50ms |
| 数据库查询优化 | 单次 < 10ms |
| 缓存命中率 | > 90% |

---

## 2. API 接口列表

### 2.1 待办管理

#### 2.1.1 创建待办

**接口地址**：`POST /api/v1/tasks`

**功能描述**：创建新的待办事项，支持根据任务类型动态设置自定义字段。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| type_id | integer | 否 | 任务类型ID（不传则创建普通待办） |
| title | string | 是 | 任务标题（1-200字符） |
| description | string | 否 | 任务描述 |
| priority | string | 否 | 优先级：low/medium/high/urgent（默认：medium） |
| due_date | string | 否 | 截止日期（ISO 8601格式） |
| category_id | string | 否 | 分类ID |
| custom_fields | object | 否 | 自定义字段（根据任务类型定义） |

**请求示例**：
```json
{
  "type_id": 1,
  "title": "提交本周工作周报",
  "description": "总结本周完成的工作内容和下周计划",
  "priority": "high",
  "due_date": "2025-01-10T18:00:00Z",
  "category_id": "cat_work",
  "custom_fields": {
    "week_start": "2025-01-06",
    "achievements": "完成了用户管理模块的开发",
    "next_plan": "开始知识库模块的设计"
  }
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_001",
    "tenant_id": "tenant_001",
    "user_id": 10001,
    "type_id": 1,
    "title": "提交本周工作周报",
    "description": "总结本周完成的工作内容和下周计划",
    "status": "pending",
    "priority": "high",
    "due_date": "2025-01-10T18:00:00Z",
    "category_id": "cat_work",
    "custom_fields": {
      "week_start": "2025-01-06",
      "achievements": "完成了用户管理模块的开发",
      "next_plan": "开始知识库模块的设计"
    },
    "created_at": "2025-01-08T08:00:00Z",
    "updated_at": "2025-01-08T08:00:00Z"
  }
}
```

---

#### 2.1.2 列出待办

**接口地址**：`GET /api/v1/tasks`

**功能描述**：获取当前用户的待办列表，支持多维度筛选和排序。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**查询参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| status | string | 否 | 状态筛选：pending/in_progress/completed/cancelled |
| priority | string | 否 | 优先级筛选：low/medium/high/urgent |
| type_id | integer | 否 | 任务类型ID |
| category_id | string | 否 | 分类ID |
| keyword | string | 否 | 关键词搜索（标题或描述） |
| due_before | string | 否 | 截止日期上限（ISO 8601） |
| due_after | string | 否 | 截止日期下限（ISO 8601） |
| page | integer | 否 | 页码（默认：1） |
| page_size | integer | 否 | 每页数量（默认：20，最大：100） |
| sort_by | string | 否 | 排序字段：due_date/priority/created_at（默认：due_date） |
| sort_order | string | 否 | 排序方向：asc/desc（默认：asc） |

**请求示例**：
```http
GET /api/v1/tasks?status=pending&priority=high&due_before=2025-01-10T23:59:59Z&page=1&page_size=20&sort_by=priority&sort_order=desc
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 45,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": "task_001",
        "type_id": 1,
        "type_name": "周报",
        "type_icon": "📝",
        "type_color": "#52c41a",
        "title": "提交本周工作周报",
        "description": "总结本周完成的工作内容和下周计划",
        "status": "pending",
        "priority": "high",
        "due_date": "2025-01-10T18:00:00Z",
        "category_id": "cat_work",
        "category_name": "工作",
        "custom_fields": {
          "week_start": "2025-01-06",
          "achievements": "完成了用户管理模块的开发",
          "next_plan": "开始知识库模块的设计"
        },
        "created_at": "2025-01-08T08:00:00Z",
        "updated_at": "2025-01-08T08:00:00Z"
      }
    ]
  }
}
```

---

#### 2.1.3 获取待办详情

**接口地址**：`GET /api/v1/tasks/:id`

**功能描述**：获取指定待办的详细信息。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 待办ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_001",
    "tenant_id": "tenant_001",
    "user_id": 10001,
    "type_id": 1,
    "title": "提交本周工作周报",
    "description": "总结本周完成的工作内容和下周计划",
    "status": "pending",
    "priority": "high",
    "due_date": "2025-01-10T18:00:00Z",
    "completed_at": null,
    "category_id": "cat_work",
    "custom_fields": {
      "week_start": "2025-01-06",
      "achievements": "完成了用户管理模块的开发",
      "next_plan": "开始知识库模块的设计"
    },
    "created_at": "2025-01-08T08:00:00Z",
    "updated_at": "2025-01-08T08:00:00Z"
  }
}
```

---

#### 2.1.4 更新待办

**接口地址**：`PUT /api/v1/tasks/:id`

**功能描述**：更新待办信息，支持部分更新。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 待办ID |

**请求参数**（所有字段可选）：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| title | string | 否 | 任务标题 |
| description | string | 否 | 任务描述 |
| status | string | 否 | 状态：pending/in_progress/completed/cancelled |
| priority | string | 否 | 优先级 |
| due_date | string | 否 | 截止日期 |
| category_id | string | 否 | 分类ID |
| custom_fields | object | 否 | 自定义字段 |

**请求示例**：
```json
{
  "status": "in_progress",
  "priority": "urgent",
  "description": "总结本周完成的工作内容和下周计划（更新）"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_001",
    "title": "提交本周工作周报",
    "description": "总结本周完成的工作内容和下周计划（更新）",
    "status": "in_progress",
    "priority": "urgent",
    "updated_at": "2025-01-08T10:30:00Z"
  }
}
```

---

#### 2.1.5 完成待办

**接口地址**：`POST /api/v1/tasks/:id/complete`

**功能描述**：标记待办为已完成。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 待办ID |

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| completion_note | string | 否 | 完成备注 |

**请求示例**：
```json
{
  "completion_note": "已完成周报提交"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_001",
    "status": "completed",
    "completed_at": "2025-01-09T14:30:00Z",
    "completion_note": "已完成周报提交"
  }
}
```

---

#### 2.1.6 删除待办

**接口地址**：`DELETE /api/v1/tasks/:id`

**功能描述**：删除待办（软删除）。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 待办ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "Task deleted successfully"
}
```

---

#### 2.1.7 批量操作待办

**接口地址**：`POST /api/v1/tasks/batch`

**功能描述**：批量更新待办状态或删除待办。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| action | string | 是 | 操作类型：complete/delete/cancel |
| task_ids | array | 是 | 待办ID列表 |
| completion_note | string | 否 | 完成备注（action=complete时可用） |

**请求示例**：
```json
{
  "action": "complete",
  "task_ids": ["task_001", "task_002", "task_003"],
  "completion_note": "批量完成周任务"
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success_count": 3,
    "failed_count": 0,
    "failed_ids": []
  }
}
```

---

### 2.2 任务类型管理

#### 2.2.1 获取任务类型列表

**接口地址**：`GET /api/v1/task-types`

**功能描述**：获取可用的任务类型列表（包括平台预置和企业自定义）。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**查询参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| is_active | boolean | 否 | 是否启用（默认：true） |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "周报",
      "description": "工作周报提交",
      "icon": "📝",
      "color": "#52c41a",
      "is_system": true,
      "fields": [
        {
          "name": "week_start",
          "label": "周开始日期",
          "type": "date",
          "required": true
        },
        {
          "name": "achievements",
          "label": "本周成果",
          "type": "textarea",
          "required": true
        },
        {
          "name": "next_plan",
          "label": "下周计划",
          "type": "textarea",
          "required": false
        }
      ],
      "sort_order": 1,
      "is_active": true
    },
    {
      "id": 2,
      "name": "请假审批",
      "description": "员工请假申请审批",
      "icon": "📅",
      "color": "#1890ff",
      "is_system": true,
      "fields": [
        {
          "name": "applicant",
          "label": "申请人",
          "type": "text",
          "required": true
        },
        {
          "name": "days",
          "label": "请假天数",
          "type": "number",
          "required": true
        },
        {
          "name": "reason",
          "label": "请假原因",
          "type": "textarea",
          "required": true
        }
      ],
      "sort_order": 2,
      "is_active": true
    },
    {
      "id": 3,
      "name": "会议准备",
      "description": "会议相关准备事项",
      "icon": "📞",
      "color": "#faad14",
      "is_system": true,
      "fields": [
        {
          "name": "meeting_title",
          "label": "会议主题",
          "type": "text",
          "required": true
        },
        {
          "name": "attendees",
          "label": "参会人员",
          "type": "textarea",
          "required": false
        },
        {
          "name": "meeting_date",
          "label": "会议时间",
          "type": "datetime",
          "required": true
        }
      ],
      "sort_order": 3,
      "is_active": true
    }
  ]
}
```

---

#### 2.2.2 获取任务类型详情

**接口地址**：`GET /api/v1/task-types/:id`

**功能描述**：获取指定任务类型的详细配置信息。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | integer | 是 | 任务类型ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "tenant_id": null,
    "name": "周报",
    "description": "工作周报提交",
    "icon": "📝",
    "color": "#52c41a",
    "is_system": true,
    "fields": [
      {
        "name": "week_start",
        "label": "周开始日期",
        "type": "date",
        "required": true,
        "placeholder": "请选择日期"
      },
      {
        "name": "achievements",
        "label": "本周成果",
        "type": "textarea",
        "required": true,
        "placeholder": "请输入本周主要工作成果",
        "rows": 4
      },
      {
        "name": "next_plan",
        "label": "下周计划",
        "type": "textarea",
        "required": false,
        "placeholder": "请输入下周工作计划",
        "rows": 4
      }
    ],
    "validation_rules": {
      "achievements": {
        "min_length": 10,
        "max_length": 2000
      }
    },
    "sort_order": 1,
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### 2.2.3 创建自定义任务类型（企业）

**接口地址**：`POST /api/v1/task-types`

**功能描述**：创建企业自定义任务类型（需要企业管理员权限）。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 类型名称（1-50字符） |
| description | string | 否 | 类型描述（1-200字符） |
| icon | string | 否 | 图标 emoji |
| color | string | 否 | 主题颜色（十六进制） |
| fields | array | 是 | 字段定义 |
| validation_rules | object | 否 | 验证规则 |
| sort_order | integer | 否 | 排序顺序 |

**fields 数组项结构**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 字段名称（英文） |
| label | string | 是 | 字段标签 |
| type | string | 是 | 字段类型：text/textarea/number/date/datetime/select |
| required | boolean | 否 | 是否必填 |
| placeholder | string | 否 | 占位符 |
| options | array | 否 | 选项（type=select时必填） |
| rows | integer | 否 | 行数（type=textarea时） |
| default | any | 否 | 默认值 |

**请求示例**：
```json
{
  "name": "客户跟进",
  "description": "客户跟进记录",
  "icon": "👤",
  "color": "#722ed1",
  "fields": [
    {
      "name": "customer_name",
      "label": "客户名称",
      "type": "text",
      "required": true,
      "placeholder": "请输入客户名称"
    },
    {
      "name": "contact_date",
      "label": "联系日期",
      "type": "date",
      "required": true
    },
    {
      "name": "contact_method",
      "label": "联系方式",
      "type": "select",
      "required": true,
      "options": [
        {"value": "phone", "label": "电话"},
        {"value": "email", "label": "邮件"},
        {"value": "wechat", "label": "微信"},
        {"value": "meeting", "label": "面谈"}
      ]
    },
    {
      "name": "follow_content",
      "label": "跟进内容",
      "type": "textarea",
      "required": true,
      "rows": 5
    },
    {
      "name": "next_action",
      "label": "下一步计划",
      "type": "textarea",
      "required": false,
      "rows": 3
    }
  ],
  "sort_order": 10
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 100,
    "tenant_id": "tenant_001",
    "name": "客户跟进",
    "description": "客户跟进记录",
    "icon": "👤",
    "color": "#722ed1",
    "is_system": false,
    "fields": [
      {
        "name": "customer_name",
        "label": "客户名称",
        "type": "text",
        "required": true,
        "placeholder": "请输入客户名称"
      }
    ],
    "sort_order": 10,
    "is_active": true,
    "created_at": "2025-01-08T10:00:00Z"
  }
}
```

---

#### 2.2.4 更新任务类型

**接口地址**：`PUT /api/v1/task-types/:id`

**功能描述**：更新企业自定义任务类型（仅限企业创建的类型）。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | integer | 是 | 任务类型ID |

**请求参数**（所有字段可选）：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 类型名称 |
| description | string | 否 | 类型描述 |
| icon | string | 否 | 图标 |
| color | string | 否 | 主题颜色 |
| fields | array | 否 | 字段定义 |
| validation_rules | object | 否 | 验证规则 |
| is_active | boolean | 否 | 是否启用 |
| sort_order | integer | 否 | 排序顺序 |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 100,
    "name": "客户跟进（已更新）",
    "updated_at": "2025-01-08T11:00:00Z"
  }
}
```

---

#### 2.2.5 删除任务类型

**接口地址**：`DELETE /api/v1/task-types/:id`

**功能描述**：删除企业自定义任务类型（仅限企业创建的类型，且无关联待办）。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | integer | 是 | 任务类型ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "Task type deleted successfully"
}
```

---

### 2.3 动态字段管理

#### 2.3.1 验证动态字段

**接口地址**：`POST /api/v1/task-types/:id/validate`

**功能描述**：验证自定义字段数据是否符合任务类型的验证规则。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | integer | 是 | 任务类型ID |

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| custom_fields | object | 是 | 自定义字段数据 |

**请求示例**：
```json
{
  "custom_fields": {
    "week_start": "2025-01-06",
    "achievements": "完成了用户管理模块的开发",
    "next_plan": "开始知识库模块的设计"
  }
}
```

**响应示例（成功）**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "valid": true,
    "errors": []
  }
}
```

**响应示例（失败）**：
```json
{
  "code": 400,
  "message": "Validation failed",
  "data": {
    "valid": false,
    "errors": [
      {
        "field": "achievements",
        "message": "本周成果长度不能少于10个字符"
      },
      {
        "field": "week_start",
        "message": "周开始日期为必填项"
      }
    ]
  }
}
```

---

### 2.4 分类管理

#### 2.4.1 获取分类列表

**接口地址**：`GET /api/v1/task-categories`

**功能描述**：获取当前用户的任务分类列表。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "cat_work",
      "name": "工作",
      "icon": "💼",
      "color": "#1890ff",
      "sort_order": 1,
      "task_count": 25
    },
    {
      "id": "cat_personal",
      "name": "个人",
      "icon": "👤",
      "color": "#52c41a",
      "sort_order": 2,
      "task_count": 12
    },
    {
      "id": "cat_study",
      "name": "学习",
      "icon": "📚",
      "color": "#722ed1",
      "sort_order": 3,
      "task_count": 8
    }
  ]
}
```

---

#### 2.4.2 创建分类

**接口地址**：`POST /api/v1/task-categories`

**功能描述**：创建自定义任务分类。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**请求参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 是 | 分类名称（1-20字符） |
| icon | string | 否 | 图标 emoji |
| color | string | 否 | 主题颜色（十六进制） |
| sort_order | integer | 否 | 排序顺序 |

**请求示例**：
```json
{
  "name": "健康",
  "icon": "🏃",
  "color": "#fa541c",
  "sort_order": 4
}
```

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "cat_health",
    "name": "健康",
    "icon": "🏃",
    "color": "#fa541c",
    "sort_order": 4,
    "task_count": 0,
    "created_at": "2025-01-08T12:00:00Z"
  }
}
```

---

#### 2.4.3 更新分类

**接口地址**：`PUT /api/v1/task-categories/:id`

**功能描述**：更新任务分类。

**请求头**：
```http
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 分类ID |

**请求参数**（所有字段可选）：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| name | string | 否 | 分类名称 |
| icon | string | 否 | 图标 |
| color | string | 否 | 主题颜色 |
| sort_order | integer | 否 | 排序顺序 |

**响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "cat_health",
    "name": "健康（已更新）",
    "updated_at": "2025-01-08T12:30:00Z"
  }
}
```

---

#### 2.4.4 删除分类

**接口地址**：`DELETE /api/v1/task-categories/:id`

**功能描述**：删除任务分类（关联的待办将移至"未分类"）。

**请求头**：
```http
X-Tenant-ID: tenant_001
Authorization: Bearer {access_token}
```

**路径参数**：

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| id | string | 是 | 分类ID |

**响应示例**：
```json
{
  "code": 0,
  "message": "Category deleted successfully"
}
```

---

## 3. 数据模型

### 3.1 待办（Task）

```typescript
interface Task {
  // 基础信息
  id: string;                    // 待办ID（UUID）
  tenant_id: string;             // 租户ID
  user_id: number;               // 用户ID

  // 任务类型
  type_id?: number;              // 任务类型ID（可选）
  type_name?: string;            // 任务类型名称（查询时返回）
  type_icon?: string;            // 任务类型图标
  type_color?: string;           // 任务类型颜色

  // 基础字段
  title: string;                 // 任务标题
  description?: string;          // 任务描述
  status: TaskStatus;            // 状态
  priority: TaskPriority;        // 优先级

  // 时间字段
  due_date?: string;             // 截止日期（ISO 8601）
  completed_at?: string;         // 完成时间

  // 分类
  category_id?: string;          // 分类ID
  category_name?: string;        // 分类名称（查询时返回）

  // 自定义字段（JSON）
  custom_fields?: Record<string, any>;  // 自定义字段数据

  // 元数据
  created_at: string;            // 创建时间
  updated_at: string;            // 更新时间
}

type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled';
type TaskPriority = 'low' | 'medium' | 'high' | 'urgent';
```

---

### 3.2 任务类型（TaskType）

```typescript
interface TaskType {
  id: number;                    // 类型ID
  tenant_id?: string;            // 租户ID（NULL表示平台预置）
  name: string;                  // 类型名称
  description?: string;          // 类型描述
  icon?: string;                 // 图标 emoji
  color?: string;                // 主题颜色

  // 字段配置
  fields: FieldConfig[];         // 字段定义
  validation_rules?: ValidationRules;  // 验证规则

  // 元数据
  is_system: boolean;            // 是否系统预置
  sort_order: number;            // 排序顺序
  is_active: boolean;            // 是否启用
  created_at: string;            // 创建时间
}

// 字段配置
interface FieldConfig {
  name: string;                  // 字段名称
  label: string;                 // 字段标签
  type: FieldType;               // 字段类型
  required: boolean;             // 是否必填
  placeholder?: string;          // 占位符
  options?: SelectOption[];      // 选项（select类型）
  rows?: number;                 // 行数（textarea类型）
  default?: any;                 // 默认值
  min_length?: number;           // 最小长度
  max_length?: number;           // 最大长度
  min?: number;                  // 最小值（number类型）
  max?: number;                  // 最大值（number类型）
}

type FieldType = 'text' | 'textarea' | 'number' | 'date' | 'datetime' | 'select';

interface SelectOption {
  value: string;                 // 选项值
  label: string;                 // 选项标签
}

// 验证规则
interface ValidationRules {
  [fieldName: string]: {
    min_length?: number;
    max_length?: number;
    min?: number;
    max?: number;
    pattern?: string;            // 正则表达式
    custom_validator?: string;   // 自定义验证器名称
  };
}
```

---

### 3.3 任务分类（TaskCategory）

```typescript
interface TaskCategory {
  id: string;                    // 分类ID
  tenant_id: string;             // 租户ID
  user_id: number;               // 用户ID
  name: string;                  // 分类名称
  icon?: string;                 // 图标 emoji
  color?: string;                // 主题颜色
  sort_order: number;            // 排序顺序
  task_count: number;            // 待办数量
  created_at: string;            // 创建时间
}
```

---

### 3.4 创建待办请求（CreateTaskRequest）

```typescript
interface CreateTaskRequest {
  type_id?: number;              // 任务类型ID（可选）
  title: string;                 // 任务标题
  description?: string;          // 任务描述
  priority?: TaskPriority;       // 优先级（默认：medium）
  due_date?: string;             // 截止日期
  category_id?: string;          // 分类ID
  custom_fields?: Record<string, any>;  // 自定义字段
}
```

---

### 3.5 列出待办请求（ListTasksRequest）

```typescript
interface ListTasksRequest {
  // 筛选条件
  status?: TaskStatus;           // 状态筛选
  priority?: TaskPriority;       // 优先级筛选
  type_id?: number;              // 任务类型ID
  category_id?: string;          // 分类ID
  keyword?: string;              // 关键词搜索
  due_before?: string;           // 截止日期上限
  due_after?: string;            // 截止日期下限

  // 分页
  page?: number;                 // 页码（默认：1）
  page_size?: number;            // 每页数量（默认：20，最大：100）

  // 排序
  sort_by?: 'due_date' | 'priority' | 'created_at';  // 排序字段
  sort_order?: 'asc' | 'desc';   // 排序方向
}
```

---

### 3.6 列出待办响应（ListTasksResponse）

```typescript
interface ListTasksResponse {
  total: number;                 // 总数
  page: number;                  // 当前页
  page_size: number;             // 每页数量
  items: Task[];                 // 待办列表
}
```

---

## 4. 错误码定义

### 4.1 通用错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 0 | 成功 | 200 |
| 40001 | 参数错误 | 400 |
| 40002 | 参数缺失 | 400 |
| 40003 | 参数格式错误 | 400 |
| 40101 | 未授权 | 401 |
| 40301 | 无权限 | 403 |
| 40401 | 资源不存在 | 404 |
| 40901 | 资源冲突 | 409 |
| 50001 | 服务器内部错误 | 500 |
| 50002 | 数据库错误 | 500 |

---

### 4.2 待办相关错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 41001 | 待办不存在 | 404 |
| 41002 | 待办已完成，无法再次完成 | 400 |
| 41003 | 待办已取消，无法操作 | 400 |
| 41004 | 批量操作部分失败 | 207 |
| 41005 | 待办标题重复 | 409 |

---

### 4.3 任务类型相关错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 42001 | 任务类型不存在 | 404 |
| 42002 | 任务类型名称重复 | 409 |
| 42003 | 任务类型被使用，无法删除 | 409 |
| 42004 | 无法修改系统预置类型 | 403 |
| 42005 | 无法删除系统预置类型 | 403 |
| 42006 | 字段配置格式错误 | 400 |
| 42007 | 字段名称重复 | 400 |

---

### 4.4 分类相关错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 43001 | 分类不存在 | 404 |
| 43002 | 分类名称重复 | 409 |
| 43003 | 分类数量超限（最多20个） | 400 |

---

### 4.5 字段验证相关错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| 44001 | 必填字段缺失 | 400 |
| 44002 | 字段长度不符合要求 | 400 |
| 44003 | 字段值超出范围 | 400 |
| 44004 | 字段格式错误 | 400 |
| 44005 | 选项值无效 | 400 |
| 44006 | 日期格式错误 | 400 |

---

### 4.6 错误响应示例

**参数错误示例**：
```json
{
  "code": 40001,
  "message": "Invalid parameter",
  "details": {
    "field": "due_date",
    "reason": "日期格式错误，应为 ISO 8601 格式"
  }
}
```

**字段验证失败示例**：
```json
{
  "code": 44001,
  "message": "Required field missing",
  "details": {
    "field": "week_start",
    "message": "周开始日期为必填项"
  }
}
```

**批量操作部分失败示例**：
```json
{
  "code": 41004,
  "message": "Batch operation partially failed",
  "details": {
    "success_count": 8,
    "failed_count": 2,
    "failed_ids": ["task_003", "task_007"],
    "errors": [
      {
        "id": "task_003",
        "reason": "待办已完成，无法再次完成"
      },
      {
        "id": "task_007",
        "reason": "待办不存在"
      }
    ]
  }
}
```

---

## 5. 调用示例

### 5.1 创建普通待办

**场景**：用户创建一个简单的待办事项。

```http
POST /api/v1/tasks HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "title": "购买办公室用品",
  "description": "购买打印纸、笔、文件夹等办公用品",
  "priority": "medium",
  "due_date": "2025-01-12T18:00:00Z",
  "category_id": "cat_work"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_001",
    "title": "购买办公室用品",
    "status": "pending",
    "priority": "medium",
    "created_at": "2025-01-08T08:00:00Z"
  }
}
```

---

### 5.2 创建周报任务（带自定义字段）

**场景**：用户创建一个周报任务，填写自定义字段。

```http
POST /api/v1/tasks HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type_id": 1,
  "title": "提交本周工作周报",
  "description": "总结本周完成的工作内容和下周计划",
  "priority": "high",
  "due_date": "2025-01-10T18:00:00Z",
  "custom_fields": {
    "week_start": "2025-01-06",
    "achievements": "1. 完成了用户管理模块的开发\n2. 修复了3个线上bug\n3. 参与了代码评审",
    "next_plan": "1. 开始知识库模块的设计\n2. 完善单元测试\n3. 编写技术文档"
  }
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_002",
    "type_id": 1,
    "type_name": "周报",
    "type_icon": "📝",
    "type_color": "#52c41a",
    "title": "提交本周工作周报",
    "status": "pending",
    "priority": "high",
    "custom_fields": {
      "week_start": "2025-01-06",
      "achievements": "1. 完成了用户管理模块的开发\n2. 修复了3个线上bug\n3. 参与了代码评审",
      "next_plan": "1. 开始知识库模块的设计\n2. 完善单元测试\n3. 编写技术文档"
    },
    "created_at": "2025-01-08T08:00:00Z"
  }
}
```

---

### 5.3 查询待办列表（筛选高优先级待办）

**场景**：查询所有高优先级且未完成的待办，按截止日期升序排列。

```http
GET /api/v1/tasks?priority=high&status=pending&page=1&page_size=20&sort_by=due_date&sort_order=asc HTTP/1.1
Host: api.zker.com
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": "task_002",
        "type_name": "周报",
        "title": "提交本周工作周报",
        "status": "pending",
        "priority": "high",
        "due_date": "2025-01-10T18:00:00Z"
      },
      {
        "id": "task_003",
        "type_name": "会议准备",
        "title": "准备产品评审会议材料",
        "status": "pending",
        "priority": "high",
        "due_date": "2025-01-11T10:00:00Z"
      }
    ]
  }
}
```

---

### 5.4 完成待办

**场景**：用户完成待办并添加完成备注。

```http
POST /api/v1/tasks/task_002/complete HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "completion_note": "已完成周报提交并通过审核"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_002",
    "status": "completed",
    "completed_at": "2025-01-09T14:30:00Z",
    "completion_note": "已完成周报提交并通过审核"
  }
}
```

---

### 5.5 批量完成待办

**场景**：用户批量完成多个待办。

```http
POST /api/v1/tasks/batch HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "action": "complete",
  "task_ids": ["task_003", "task_004", "task_005"],
  "completion_note": "周度批量完成"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "success_count": 3,
    "failed_count": 0,
    "failed_ids": []
  }
}
```

---

### 5.6 获取任务类型列表

**场景**：用户选择创建任务时，获取可用的任务类型。

```http
GET /api/v1/task-types HTTP/1.1
Host: api.zker.com
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "周报",
      "description": "工作周报提交",
      "icon": "📝",
      "color": "#52c41a",
      "is_system": true,
      "fields": [
        {
          "name": "week_start",
          "label": "周开始日期",
          "type": "date",
          "required": true
        },
        {
          "name": "achievements",
          "label": "本周成果",
          "type": "textarea",
          "required": true
        },
        {
          "name": "next_plan",
          "label": "下周计划",
          "type": "textarea",
          "required": false
        }
      ]
    },
    {
      "id": 2,
      "name": "请假审批",
      "description": "员工请假申请审批",
      "icon": "📅",
      "color": "#1890ff",
      "is_system": true,
      "fields": [
        {
          "name": "applicant",
          "label": "申请人",
          "type": "text",
          "required": true
        },
        {
          "name": "days",
          "label": "请假天数",
          "type": "number",
          "required": true
        },
        {
          "name": "reason",
          "label": "请假原因",
          "type": "textarea",
          "required": true
        }
      ]
    },
    {
      "id": 3,
      "name": "会议准备",
      "description": "会议相关准备事项",
      "icon": "📞",
      "color": "#faad14",
      "is_system": true,
      "fields": [
        {
          "name": "meeting_title",
          "label": "会议主题",
          "type": "text",
          "required": true
        },
        {
          "name": "attendees",
          "label": "参会人员",
          "type": "textarea",
          "required": false
        },
        {
          "name": "meeting_date",
          "label": "会议时间",
          "type": "datetime",
          "required": true
        }
      ]
    }
  ]
}
```

---

### 5.7 创建企业自定义任务类型

**场景**：企业管理员创建"客户跟进"任务类型。

```http
POST /api/v1/task-types HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "name": "客户跟进",
  "description": "客户跟进记录",
  "icon": "👤",
  "color": "#722ed1",
  "fields": [
    {
      "name": "customer_name",
      "label": "客户名称",
      "type": "text",
      "required": true,
      "placeholder": "请输入客户名称"
    },
    {
      "name": "contact_date",
      "label": "联系日期",
      "type": "date",
      "required": true
    },
    {
      "name": "contact_method",
      "label": "联系方式",
      "type": "select",
      "required": true,
      "options": [
        {"value": "phone", "label": "电话"},
        {"value": "email", "label": "邮件"},
        {"value": "wechat", "label": "微信"},
        {"value": "meeting", "label": "面谈"}
      ]
    },
    {
      "name": "follow_content",
      "label": "跟进内容",
      "type": "textarea",
      "required": true,
      "rows": 5
    },
    {
      "name": "next_action",
      "label": "下一步计划",
      "type": "textarea",
      "required": false,
      "rows": 3
    }
  ],
  "sort_order": 10
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 100,
    "tenant_id": "tenant_001",
    "name": "客户跟进",
    "description": "客户跟进记录",
    "icon": "👤",
    "color": "#722ed1",
    "is_system": false,
    "fields": [
      {
        "name": "customer_name",
        "label": "客户名称",
        "type": "text",
        "required": true,
        "placeholder": "请输入客户名称"
      }
    ],
    "sort_order": 10,
    "is_active": true,
    "created_at": "2025-01-08T10:00:00Z"
  }
}
```

---

### 5.8 使用自定义任务类型创建待办

**场景**：用户使用刚创建的"客户跟进"类型创建待办。

```http
POST /api/v1/tasks HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type_id": 100,
  "title": "跟进ABC公司合作意向",
  "priority": "high",
  "due_date": "2025-01-15T18:00:00Z",
  "custom_fields": {
    "customer_name": "ABC科技有限公司",
    "contact_date": "2025-01-08",
    "contact_method": "meeting",
    "follow_content": "与ABC公司张总进行了面谈，对方对我们的企业级SaaS方案很感兴趣。讨论了实施周期和价格范围，对方表示预算在50-80万之间。需要尽快提供详细方案和报价。",
    "next_action": "1. 准备详细技术方案和实施计划\n2. 提供分阶段报价方案\n3. 安排产品演示"
  }
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "task_004",
    "type_id": 100,
    "type_name": "客户跟进",
    "type_icon": "👤",
    "type_color": "#722ed1",
    "title": "跟进ABC公司合作意向",
    "status": "pending",
    "priority": "high",
    "due_date": "2025-01-15T18:00:00Z",
    "custom_fields": {
      "customer_name": "ABC科技有限公司",
      "contact_date": "2025-01-08",
      "contact_method": "meeting",
      "follow_content": "与ABC公司张总进行了面谈，对方对我们的企业级SaaS方案很感兴趣。讨论了实施周期和价格范围，对方表示预算在50-80万之间。需要尽快提供详细方案和报价。",
      "next_action": "1. 准备详细技术方案和实施计划\n2. 提供分阶段报价方案\n3. 安排产品演示"
    },
    "created_at": "2025-01-08T14:00:00Z"
  }
}
```

---

### 5.9 创建任务分类

**场景**：用户创建"健康"分类。

```http
POST /api/v1/task-categories HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "name": "健康",
  "icon": "🏃",
  "color": "#fa541c",
  "sort_order": 4
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "cat_health",
    "name": "健康",
    "icon": "🏃",
    "color": "#fa541c",
    "sort_order": 4,
    "task_count": 0,
    "created_at": "2025-01-08T12:00:00Z"
  }
}
```

---

### 5.10 字段验证错误示例

**场景**：用户提交的任务数据不符合字段验证规则。

```http
POST /api/v1/task-types/1/validate HTTP/1.1
Host: api.zker.com
Content-Type: application/json
X-Tenant-ID: tenant_001
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "custom_fields": {
    "week_start": "2025-01-06",
    "achievements": "太短"
  }
}
```

**响应（验证失败）**：
```json
{
  "code": 400,
  "message": "Validation failed",
  "data": {
    "valid": false,
    "errors": [
      {
        "field": "achievements",
        "message": "本周成果长度不能少于10个字符"
      },
      {
        "field": "next_plan",
        "message": "下周计划为必填项"
      }
    ]
  }
}
```

---

## 附录

### A. 预置任务类型

系统预置以下任务类型：

1. **周报**（📝 #52c41a）
   - week_start: 周开始日期（date）
   - achievements: 本周成果（textarea）
   - next_plan: 下周计划（textarea）

2. **请假审批**（📅 #1890ff）
   - applicant: 申请人（text）
   - days: 请假天数（number）
   - reason: 请假原因（textarea）

3. **会议准备**（📞 #faad14）
   - meeting_title: 会议主题（text）
   - attendees: 参会人员（textarea）
   - meeting_date: 会议时间（datetime）

### B. 支持的字段类型

| 类型 | 说明 | 适用场景 |
|------|------|---------|
| text | 单行文本 | 短文本输入（名称、标题等） |
| textarea | 多行文本 | 长文本输入（描述、内容等） |
| number | 数字 | 数值输入（天数、金额等） |
| date | 日期 | 日期选择（日期） |
| datetime | 日期时间 | 日期时间选择（精确到分钟） |
| select | 下拉选择 | 固定选项选择（类型、方式等） |

### C. 前端动态表单实现

**动态表单渲染示例**（React + TypeScript）：

```tsx
import React, { useEffect, useState } from 'react';
import { Form, Input, Select, DatePicker } from '@douyinfe/semi-ui';

interface DynamicTaskFormProps {
  typeId: number;
  onSubmit: (data: any) => void;
}

export const DynamicTaskForm: React.FC<DynamicTaskFormProps> = ({
  typeId,
  onSubmit
}) => {
  const [taskType, setTaskType] = useState<TaskType | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadTaskType(typeId);
  }, [typeId]);

  const loadTaskType = async (typeId: number) => {
    const resp = await taskAPI.getTaskType(typeId);
    setTaskType(resp.data);
  };

  const renderField = (field: FieldConfig) => {
    switch (field.type) {
      case 'text':
        return (
          <Input
            placeholder={field.placeholder || `请输入${field.label}`}
          />
        );
      case 'textarea':
        return (
          <Input.TextArea
            rows={field.rows || 4}
            placeholder={field.placeholder || `请输入${field.label}`}
          />
        );
      case 'number':
        return (
          <Input
            type="number"
            placeholder={field.placeholder || `请输入${field.label}`}
          />
        );
      case 'date':
        return <DatePicker type="date" />;
      case 'datetime':
        return <DatePicker type="dateTime" />;
      case 'select':
        return (
          <Select
            placeholder={field.placeholder || `请选择${field.label}`}
            optionList={field.options || []}
          />
        );
      default:
        return null;
    }
  };

  const handleSubmit = (values: any) => {
    onSubmit(values);
  };

  return (
    <Form form={form} onSubmit={handleSubmit}>
      <Form.Input
        field="title"
        label="任务标题"
        required
        style={{ width: '100%' }}
      />
      <Form.TextArea
        field="description"
        label="描述"
        rows={3}
        style={{ width: '100%' }}
      />

      {/* 动态字段 */}
      {taskType?.fields.map(field => (
        <Form.Input
          key={field.name}
          field={`customFields.${field.name}`}
          label={field.label}
          required={field.required}
          style={{ width: '100%' }}
        >
          {renderField(field)}
        </Form.Input>
      ))}

      <Button type="primary" htmlType="submit">
        创建任务
      </Button>
    </Form>
  );
};
```

---

**文档版本**: v1.0.0
**最后更新**: 2025-01-08
**维护团队**: ZKER Enterprise Team

---

**文档结束**
