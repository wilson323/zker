# API接口文档：技能资源管理模块

**模块名称**: 技能资源管理 (Skill Management)
**设计文档**: 08-企业管理中台_技能资源管理.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 技能商店API](#2-技能商店api)
- [3. 技能管理API](#3-技能管理api)
- [4. 技能授权API](#4-技能授权api)
- [5. 技能编排API](#5-技能编排api)
- [6. 技能分类API](#6-技能分类api)
- [7. 技能统计API](#7-技能统计api)
- [8. 数据模型](#8-数据模型)
- [9. 错误码定义](#9-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

技能资源管理是ZKER企业级SaaS平台的核心模块，通过**100%复用zker Plugin系统**的方式，为企业提供技能资源的管理、编排、授权能力：

- ✅ **技能商店** - 浏览、搜索、安装、预览技能
- ✅ **技能管理** - 查看、配置、启用/禁用技能
- ✅ **技能授权** - 基于RBAC的技能权限控制
- ✅ **技能编排** - 可视化拖拽式技能流程设计
- ✅ **技能分类** - 企业自定义技能分类和标签
- ✅ **技能统计** - 使用量、成功率、错误日志监控

### 1.2 概念映射

| 鲸智百应概念 | zker 概念 | 实现方式 |
|-------------|-----------|----------|
| 技能 (Skill) | Plugin (插件) | 直接映射 |
| 技能编排 | Workflow (工作流) | 直接复用Workflow系统 |
| 技能商店 | Plugin Store | 直接复用 |
| 技能授权 | RBAC权限系统 | 通过role_permissions表配置 |
| 技能分类 | Plugin Tags | 通过skill_categories表配置 |

### 1.3 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5 (100%复用zker plugins表)
- 缓存: Redis 8.0
- 搜索: Elasticsearch

**实现策略**: ✅ 0% 编码 + 100% 配置/复用

---

## 2. 技能商店API

### 2.1 获取技能商店列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/store`

**查询参数**:
- category_id: 分类ID过滤 (可选)
- tag: 标签过滤 (可选)
- keyword: 搜索关键词 (可选)
- sort: 排序方式 (popular/latest/rating)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "items": [
      {
        "id": "plugin-001",
        "name": "天气查询",
        "description": "查询各地天气情况",
        "icon": "https://cdn.example.com/skills/weather.png",
        "category": "工具",
        "tags": ["查询", "天气", "生活"],
        "rating": 4.5,
        "review_count": 234,
        "install_count": 1523,
        "is_installed": false,
        "is_official": true,
        "author": "zker官方",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 2.2 获取技能详情

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/store/{skill_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "plugin-001",
    "name": "天气查询",
    "description": "查询各地天气情况，支持全球城市",
    "icon": "https://cdn.example.com/skills/weather.png",
    "screenshots": [
      "https://cdn.example.com/skills/weather/s1.png",
      "https://cdn.example.com/skills/weather/s2.png"
    ],
    "category": "工具",
    "tags": ["查询", "天气", "生活"],
    "version": "v1.2.0",
    "author": "zker官方",
    "is_official": true,
    "rating": 4.5,
    "review_count": 234,
    "install_count": 1523,
    "last_updated": "2025-01-10T00:00:00Z",
    "config_schema": {
      "type": "object",
      "properties": {
        "api_key": {
          "type": "string",
          "title": "API密钥"
        },
        "default_city": {
          "type": "string",
          "title": "默认城市"
        }
      }
    },
    "permissions": ["user_info", "network"],
    "is_installed": true,
    "installed_at": "2025-01-15T10:00:00Z"
  }
}
```

### 2.3 安装技能

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/store/{skill_id}/install`

**请求参数**:
```json
{
  "config": {
    "api_key": "sk-xxx",
    "default_city": "北京"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "安装成功",
  "data": {
    "skill_id": "plugin-001",
    "installed_at": "2025-01-15T14:00:00Z"
  }
}
```

### 2.4 卸载技能

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/{skill_id}/uninstall`

**响应示例**:
```json
{
  "code": 0,
  "message": "卸载成功",
  "data": {
    "skill_id": "plugin-001",
    "uninstalled_at": "2025-01-15T15:00:00Z"
  }
}
```

### 2.5 预览技能

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/store/{skill_id}/preview`

**请求参数**:
```json
{
  "action": "query_weather",
  "params": {
    "city": "北京"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "result": {
      "city": "北京",
      "temperature": "15°C",
      "condition": "晴",
      "humidity": "45%"
    }
  }
}
```

### 2.6 获取技能评价

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/{skill_id}/reviews`

**查询参数**:
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "average_rating": 4.5,
    "items": [
      {
        "id": "review-001",
        "user_name": "张三",
        "rating": 5,
        "comment": "非常好用的技能",
        "created_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

---

## 3. 技能管理API

### 3.1 获取已安装技能列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/installed`

**查询参数**:
- category_id: 分类ID过滤 (可选)
- status: 状态过滤 (enabled/disabled)
- keyword: 搜索关键词 (可选)
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "id": "plugin-001",
        "name": "天气查询",
        "description": "查询各地天气情况",
        "icon": "https://cdn.example.com/skills/weather.png",
        "category": "工具",
        "status": "enabled",
        "usage_count": 234,
        "success_rate": 0.98,
        "last_used_at": "2025-01-15T14:30:00Z",
        "installed_at": "2025-01-01T00:00:00Z"
      }
    ]
  }
}
```

### 3.2 配置技能

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/skills/{skill_id}/config`

**请求参数**:
```json
{
  "config": {
    "api_key": "sk-new-key",
    "default_city": "上海"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "配置更新成功",
  "data": {
    "skill_id": "plugin-001",
    "updated_at": "2025-01-15T16:00:00Z"
  }
}
```

### 3.3 启用技能

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/{skill_id}/enable`

**响应示例**:
```json
{
  "code": 0,
  "message": "技能已启用",
  "data": {
    "skill_id": "plugin-001",
    "status": "enabled"
  }
}
```

### 3.4 禁用技能

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/{skill_id}/disable`

**响应示例**:
```json
{
  "code": 0,
  "message": "技能已禁用",
  "data": {
    "skill_id": "plugin-001",
    "status": "disabled"
  }
}
```

---

## 4. 技能授权API

### 4.1 获取技能授权列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/{skill_id}/permissions`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "skill_id": "plugin-001",
    "skill_name": "天气查询",
    "permissions": [
      {
        "role_id": "role-001",
        "role_name": "普通员工",
        "permissions": ["use"],
        "granted_at": "2025-01-01T00:00:00Z",
        "granted_by": "管理员"
      },
      {
        "role_id": "role-002",
        "role_name": "管理员",
        "permissions": ["use", "config"],
        "granted_at": "2025-01-01T00:00:00Z",
        "granted_by": "超级管理员"
      }
    ]
  }
}
```

### 4.2 授权技能给角色

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/{skill_id}/permissions/grant`

**请求参数**:
```json
{
  "role_id": "role-003",
  "permissions": ["use", "config"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "授权成功",
  "data": {
    "role_id": "role-003",
    "skill_id": "plugin-001",
    "permissions": ["use", "config"],
    "granted_at": "2025-01-15T17:00:00Z"
  }
}
```

### 4.3 撤销技能授权

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skills/{skill_id}/permissions/revoke`

**请求参数**:
```json
{
  "role_id": "role-003"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "撤销成功",
  "data": {
    "role_id": "role-003",
    "skill_id": "plugin-001",
    "revoked_at": "2025-01-15T18:00:00Z"
  }
}
```

---

## 5. 技能编排API

### 5.1 获取技能编排列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skill-workflows`

**查询参数**:
- status: 状态过滤 (draft/published/archived)
- creator_id: 创建者过滤
- keyword: 搜索关键词
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
        "id": "workflow-001",
        "name": "客户信息查询流程",
        "description": "查询CRM客户信息并发送邮件",
        "status": "published",
        "nodes_count": 5,
        "usage_count": 89,
        "creator_name": "张三",
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 5.2 创建技能编排

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skill-workflows`

**请求参数**:
```json
{
  "name": "客户信息查询流程",
  "description": "查询CRM客户信息并发送邮件",
  "nodes": [
    {
      "id": "node-001",
      "type": "skill",
      "skill_id": "plugin-crm-query",
      "position": { "x": 100, "y": 100 },
      "config": {
        "customer_id": "${input.customer_id}"
      }
    },
    {
      "id": "node-002",
      "type": "skill",
      "skill_id": "plugin-email-send",
      "position": { "x": 300, "y": 100 },
      "config": {
        "to": "${node-001.data.email}",
        "subject": "客户信息",
        "body": "${node-001.data.content}"
      }
    }
  ],
  "edges": [
    {
      "id": "edge-001",
      "from": "node-001",
      "to": "node-002"
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
    "id": "workflow-001",
    "name": "客户信息查询流程",
    "status": "draft",
    "created_at": "2025-01-15T14:00:00Z"
  }
}
```

### 5.3 执行技能编排

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skill-workflows/{workflow_id}/execute`

**请求参数**:
```json
{
  "input": {
    "customer_id": "CUST-001"
  },
  "async": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "执行开始",
  "data": {
    "execution_id": "exec-001",
    "workflow_id": "workflow-001",
    "status": "running",
    "started_at": "2025-01-15T15:00:00Z"
  }
}
```

### 5.4 获取执行结果

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skill-workflows/executions/{execution_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "execution_id": "exec-001",
    "workflow_id": "workflow-001",
    "status": "completed",
    "result": {
      "node-001": {
        "status": "success",
        "data": {
          "customer_name": "张三",
          "email": "zhangsan@example.com"
        }
      },
      "node-002": {
        "status": "success",
        "data": {
          "sent": true
        }
      }
    },
    "started_at": "2025-01-15T15:00:00Z",
    "completed_at": "2025-01-15T15:00:02Z",
    "duration_ms": 2000
  }
}
```

---

## 6. 技能分类API

### 6.1 获取技能分类列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skill-categories`

**查询参数**:
- parent_id: 父分类ID (可选,不传则返回顶级分类)
- include_inactive: 是否包含已禁用分类 (默认false)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "items": [
      {
        "id": 1,
        "parent_id": null,
        "name": "办公",
        "description": "办公效率工具",
        "icon": "https://cdn.example.com/categories/office.png",
        "sort_order": 1,
        "is_active": true,
        "skill_count": 15,
        "children": [
          {
            "id": 11,
            "parent_id": 1,
            "name": "文档处理",
            "skill_count": 8
          }
        ]
      },
      {
        "id": 2,
        "parent_id": null,
        "name": "营销",
        "description": "营销推广工具",
        "icon": "https://cdn.example.com/categories/marketing.png",
        "sort_order": 2,
        "is_active": true,
        "skill_count": 12,
        "children": []
      }
    ]
  }
}
```

### 6.2 创建技能分类

**接口地址**: `POST /api/v1/tenants/{tenant_id}/skill-categories`

**请求参数**:
```json
{
  "parent_id": null,
  "name": "客服",
  "description": "客服相关技能",
  "icon": "https://cdn.example.com/categories/customer-service.png",
  "sort_order": 3
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 3,
    "name": "客服",
    "created_at": "2025-01-15T16:00:00Z"
  }
}
```

### 6.3 更新技能分类

**接口地址**: `PUT /api/v1/tenants/{tenant_id}/skill-categories/{category_id}`

**请求参数**:
```json
{
  "name": "客服中心",
  "description": "客服中心相关技能",
  "sort_order": 4
}
```

### 6.4 删除技能分类

**接口地址**: `DELETE /api/v1/tenants/{tenant_id}/skill-categories/{category_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_category_id": 3
  }
}
```

---

## 7. 技能统计API

### 7.1 获取技能使用统计

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/{skill_id}/stats`

**查询参数**:
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "skill_id": "plugin-001",
    "skill_name": "天气查询",
    "total_calls": 1523,
    "successful_calls": 1489,
    "failed_calls": 34,
    "success_rate": 0.977,
    "avg_response_time_ms": 235,
    "p50_response_time_ms": 200,
    "p95_response_time_ms": 350,
    "p99_response_time_ms": 500,
    "daily_stats": [
      {
        "date": "2025-01-09",
        "calls": 120,
        "success_rate": 0.98
      },
      {
        "date": "2025-01-10",
        "calls": 145,
        "success_rate": 0.97
      }
    ]
  }
}
```

### 7.2 获取企业技能概览

**接口地址**: `GET /api/v1/tenants/{tenant_id}/skills/overview`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_skills": 25,
    "enabled_skills": 22,
    "disabled_skills": 3,
    "total_calls_this_month": 15234,
    "total_calls_last_month": 12345,
    "growth_rate": 0.23,
    "top_skills": [
      {
        "skill_id": "plugin-001",
        "skill_name": "天气查询",
        "calls": 2341
      },
      {
        "skill_id": "plugin-002",
        "skill_name": "邮件发送",
        "calls": 1876
      }
    ],
    "category_distribution": [
      {
        "category": "办公",
        "skill_count": 10,
        "calls": 5000
      },
      {
        "category": "营销",
        "skill_count": 8,
        "calls": 3500
      }
    ]
  }
}
```

---

## 8. 数据模型

### 8.1 Skill (技能)

```typescript
interface Skill {
  id: string; // 对应plugins表
  name: string;
  description?: string;
  icon?: string;
  category?: string;
  tags: string[];
  version: string;
  author: string;
  is_official: boolean;
  rating: number;
  review_count: number;
  install_count: number;
  config_schema?: JSONSchema;
  permissions?: string[];
  is_installed: boolean;
  status: 'enabled' | 'disabled';
  installed_at?: Date;
}
```

### 8.2 SkillCategory (技能分类)

```typescript
interface SkillCategory {
  id: number;
  tenant_id: string;
  parent_id?: number;
  name: string;
  description?: string;
  icon?: string;
  sort_order: number;
  is_active: boolean;
  skill_count?: number;
  children?: SkillCategory[];
  created_at: Date;
  updated_at: Date;
}
```

### 8.3 SkillWorkflow (技能编排)

```typescript
interface SkillWorkflow {
  id: string;
  tenant_id: string;
  creator_id: number;
  name: string;
  description?: string;
  nodes: WorkflowNode[];
  edges: WorkflowEdge[];
  status: 'draft' | 'published' | 'archived';
  usage_count: number;
  created_at: Date;
  updated_at: Date;
}

interface WorkflowNode {
  id: string;
  type: 'skill' | 'logic' | 'variable';
  skill_id?: string;
  position: { x: number; y: number };
  config?: Record<string, any>;
}

interface WorkflowEdge {
  id: string;
  from: string; // node_id
  to: string; // node_id
}
```

### 8.4 SkillPermission (技能授权)

```typescript
interface SkillPermission {
  role_id: string;
  role_name: string;
  permissions: ('use' | 'config' | 'grant')[];
  granted_at: Date;
  granted_by: string;
}
```

### 8.5 SkillStats (技能统计)

```typescript
interface SkillStats {
  total_calls: number;
  successful_calls: number;
  failed_calls: number;
  success_rate: number;
  avg_response_time_ms: number;
  p50_response_time_ms: number;
  p95_response_time_ms: number;
  p99_response_time_ms: number;
  daily_stats: Array<{
    date: string;
    calls: number;
    success_rate: number;
  }>;
}
```

---

## 9. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 31001 | 404 | 技能不存在 |
| 31002 | 400 | 技能已安装 |
| 31003 | 400 | 技能未安装 |
| 31004 | 400 | 技能配置无效 |
| 31005 | 403 | 无权限使用此技能 |
| 31006 | 400 | 技能状态不允许此操作 |
| 31007 | 500 | 技能执行失败 |
| 31008 | 400 | 技能编排配置无效 |
| 31009 | 400 | 技能编排执行失败 |
| 31101 | 404 | 技能编排不存在 |
| 31102 | 400 | 技能分类已存在 |
| 31103 | 404 | 技能分类不存在 |
| 31104 | 400 | 技能分类下存在技能,无法删除 |
| 31201 | 400 | 技能授权已存在 |
| 31202 | 404 | 技能授权不存在 |

---

**文档结束**
