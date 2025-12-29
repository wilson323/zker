# API接口文档：数据看板模块

**模块名称**: 数据看板 (Dashboard)
**设计文档**: 11-企业管理中台_数据看板.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 看板管理API](#2-看板管理api)
- [3. 看板数据API](#3-看板数据api)
- [4. 组件类型API](#4-组件类型api)
- [5. 数据分析API](#5-数据分析api)
- [6. 数据模型](#6-数据模型)
- [7. 错误码定义](#7-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

数据看板是ZKER企业级SaaS平台的企业管理中台模块，通过**拖拽式框架 + JSON配置**的方式，提供可视化数据看板能力：

- ✅ **看板管理** - 创建、编辑、删除、复制看板
- ✅ **组件配置** - 添加图表、指标卡片、表格等组件
- ✅ **数据源配置** - 支持SQL、API、WebSocket三种数据源
- ✅ **拖拽布局** - 基于React-Grid-Layout的拖拽式布局
- ✅ **数据分析** - Bot指标监控、用户行为追踪、性能监控、成本分析

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 实现: 10% 编码（拖拽框架） + 90% 配置（JSON配置）

**前端技术栈**:
- React 18 + TypeScript 5.6
- React-Grid-Layout (拖拽布局)
- ECharts / G2Plot (图表组件)
- Semi Design (UI组件库)

**实现策略**: ✅ 10% 编码（拖拽框架） + 90% 配置（JSON配置）

### 1.3 数据库表

| 表名 | 说明 |
|------|------|
| `dashboards` | 看板表 |
| `dashboard_component_types` | 组件类型表 |
| `bot_metrics` | Bot指标数据表 |
| `user_behaviors` | 用户行为追踪表 |
| `conversation_logs` | 对话记录表 |

---

## 2. 看板管理API

### 2.1 获取看板列表

**接口地址**: `GET /api/v1/dashboards`

**查询参数**:
- keyword: 搜索关键词 (可选)
- is_public: 是否公开 (可选)
- creator_id: 创建者ID (可选)
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
        "id": "dash-001",
        "tenant_id": "tenant-001",
        "creator_id": 1001,
        "name": "Bot使用统计",
        "description": "Bot使用情况统计看板",
        "layout": {
          "cols": 12,
          "rows": 10
        },
        "components": [
          {
            "id": "comp-1",
            "type": "metric-card",
            "position": { "x": 0, "y": 0, "w": 3, "h": 2 },
            "title": "总Bot数",
            "dataSource": {
              "type": "sql",
              "query": "SELECT COUNT(*) as value FROM bots WHERE tenant_id = ?"
            },
            "style": {
              "color": "#1890ff"
            }
          }
        ],
        "is_public": false,
        "allowed_roles": ["admin", "analyst"],
        "view_count": 1520,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 2.2 创建看板

**接口地址**: `POST /api/v1/dashboards`

**请求参数**:
```json
{
  "name": "Bot使用统计",
  "description": "Bot使用情况统计看板",
  "layout": {
    "cols": 12,
    "rows": 10
  },
  "components": [
    {
      "id": "comp-1",
      "type": "metric-card",
      "position": { "x": 0, "y": 0, "w": 3, "h": 2 },
      "title": "总Bot数",
      "dataSource": {
        "type": "sql",
        "query": "SELECT COUNT(*) as value FROM bots WHERE tenant_id = ?"
      },
      "style": {
        "color": "#1890ff"
      }
    },
    {
      "id": "comp-2",
      "type": "line-chart",
      "position": { "x": 3, "y": 0, "w": 6, "h": 4 },
      "title": "Bot使用趋势",
      "dataSource": {
        "type": "api",
        "url": "/api/v1/analytics/bot-usage-trend",
        "refreshInterval": 60
      }
    },
    {
      "id": "comp-3",
      "type": "table",
      "position": { "x": 0, "y": 4, "w": 12, "h": 6 },
      "title": "Top 10 Bot",
      "dataSource": {
        "type": "sql",
        "query": "SELECT name, usage_count FROM bots ORDER BY usage_count DESC LIMIT 10"
      }
    }
  ],
  "is_public": false,
  "allowed_roles": ["admin", "analyst"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "dash-002",
    "name": "Bot使用统计",
    "created_at": "2025-01-15T14:00:00Z"
  }
}
```

### 2.3 更新看板

**接口地址**: `PUT /api/v1/dashboards/{dashboard_id}`

**请求参数**:
```json
{
  "name": "Bot使用统计(更新版)",
  "description": "Bot使用情况统计看板",
  "layout": {
    "cols": 12,
    "rows": 10
  },
  "components": []
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": "dash-001",
    "updated_at": "2025-01-15T15:00:00Z"
  }
}
```

### 2.4 删除看板

**接口地址**: `DELETE /api/v1/dashboards/{dashboard_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_dashboard_id": "dash-001"
  }
}
```

### 2.5 复制看板

**接口地址**: `POST /api/v1/dashboards/{dashboard_id}/copy`

**请求参数**:
```json
{
  "name": "Bot使用统计(副本)"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "复制成功",
  "data": {
    "id": "dash-003",
    "name": "Bot使用统计(副本)",
    "created_at": "2025-01-15T16:00:00Z"
  }
}
```

### 2.6 获取看板详情

**接口地址**: `GET /api/v1/dashboards/{dashboard_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "dash-001",
    "tenant_id": "tenant-001",
    "creator_id": 1001,
    "creator_name": "张三",
    "name": "Bot使用统计",
    "description": "Bot使用情况统计看板",
    "layout": {
      "cols": 12,
      "rows": 10
    },
    "components": [
      {
        "id": "comp-1",
        "type": "metric-card",
        "position": { "x": 0, "y": 0, "w": 3, "h": 2 },
        "title": "总Bot数",
        "dataSource": {
          "type": "sql",
          "query": "SELECT COUNT(*) as value FROM bots WHERE tenant_id = ?"
        },
        "style": {
          "color": "#1890ff"
        }
      }
    ],
    "is_public": false,
    "allowed_roles": ["admin", "analyst"],
    "view_count": 1520,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

---

## 3. 看板数据API

### 3.1 获取看板数据

**接口地址**: `GET /api/v1/dashboards/{dashboard_id}/data`

**查询参数**:
- refresh: 是否强制刷新 (可选，默认false)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "dashboard_id": "dash-001",
    "timestamp": "2025-01-15T10:30:00Z",
    "components": {
      "comp-1": {
        "value": 156,
        "trend": "+12.5%",
        "updated_at": "2025-01-15T10:30:00Z"
      },
      "comp-2": {
        "data": [
          { "date": "2025-01-09", "value": 1200 },
          { "date": "2025-01-10", "value": 1350 },
          { "date": "2025-01-11", "value": 1280 },
          { "date": "2025-01-12", "value": 1420 },
          { "date": "2025-01-13", "value": 1560 },
          { "date": "2025-01-14", "value": 1490 },
          { "date": "2025-01-15", "value": 1680 }
        ],
        "updated_at": "2025-01-15T10:30:00Z"
      },
      "comp-3": {
        "data": [
          { "name": "客服助手", "usage_count": 3520 },
          { "name": "销售顾问", "usage_count": 2890 },
          { "name": "技术支持", "usage_count": 2150 }
        ],
        "updated_at": "2025-01-15T10:30:00Z"
      }
    }
  }
}
```

### 3.2 获取组件数据

**接口地址**: `GET /api/v1/dashboards/{dashboard_id}/components/{component_id}/data`

**查询参数**:
- refresh: 是否强制刷新 (可选，默认false)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "component_id": "comp-1",
    "value": 156,
    "trend": "+12.5%",
    "updated_at": "2025-01-15T10:30:00Z"
  }
}
```

### 3.3 WebSocket实时数据推送

**WebSocket地址**: `ws://api.example.com/api/v1/dashboards/{dashboard_id}/realtime`

**连接参数**:
- token: JWT认证令牌

**消息格式**:
```json
{
  "type": "data_update",
  "component_id": "comp-1",
  "data": {
    "value": 158,
    "trend": "+13.2%",
    "updated_at": "2025-01-15T10:31:00Z"
  }
}
```

---

## 4. 组件类型API

### 4.1 获取组件类型列表

**接口地址**: `GET /api/v1/dashboards/components/types`

**查询参数**:
- category: 组件分类 (可选，metric/chart/table/text/other)
- is_active: 是否启用 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 6,
    "items": [
      {
        "id": 1,
        "name": "metric-card",
        "category": "metric",
        "description": "指标卡片",
        "icon": "📊",
        "config_schema": {
          "title": { "type": "string", "required": true },
          "dataSource": { "type": "object", "required": true },
          "style": { "type": "object", "required": false }
        },
        "is_active": true
      },
      {
        "id": 2,
        "name": "line-chart",
        "category": "chart",
        "description": "折线图",
        "icon": "📈",
        "config_schema": {
          "title": { "type": "string", "required": true },
          "dataSource": { "type": "object", "required": true },
          "xField": { "type": "string", "required": true },
          "yField": { "type": "string", "required": true }
        },
        "is_active": true
      },
      {
        "id": 3,
        "name": "bar-chart",
        "category": "chart",
        "description": "柱状图",
        "icon": "📊",
        "is_active": true
      },
      {
        "id": 4,
        "name": "pie-chart",
        "category": "chart",
        "description": "饼图",
        "icon": "🍰",
        "is_active": true
      },
      {
        "id": 5,
        "name": "table",
        "category": "table",
        "description": "数据表格",
        "icon": "📋",
        "is_active": true
      },
      {
        "id": 6,
        "name": "text",
        "category": "text",
        "description": "文本组件",
        "icon": "📝",
        "is_active": true
      }
    ]
  }
}
```

### 4.2 获取预设模板

**接口地址**: `GET /api/v1/dashboards/templates`

**查询参数**:
- category: 模板分类 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 3,
    "items": [
      {
        "id": "tpl-bot-stats",
        "name": "Bot使用统计模板",
        "description": "Bot使用情况统计看板模板",
        "category": "bot_analytics",
        "thumbnail": "/templates/bot-stats.png",
        "layout": {
          "cols": 12,
          "rows": 10
        },
        "components": [
          {
            "id": "comp-total-bots",
            "type": "metric-card",
            "position": { "x": 0, "y": 0, "w": 3, "h": 2 },
            "title": "总Bot数",
            "dataSource": {
              "type": "sql",
              "query": "SELECT COUNT(*) as value FROM bots WHERE deleted_at IS NULL"
            }
          },
          {
            "id": "comp-usage-trend",
            "type": "line-chart",
            "position": { "x": 3, "y": 0, "w": 6, "h": 4 },
            "title": "近7日使用趋势",
            "dataSource": {
              "type": "api",
              "url": "/api/v1/analytics/bot-usage-trend?days=7"
            }
          }
        ],
        "is_official": true
      },
      {
        "id": "tpl-user-activity",
        "name": "用户活跃度模板",
        "description": "用户活跃度分析看板模板",
        "category": "user_analytics",
        "is_official": true
      },
      {
        "id": "tpl-system-perf",
        "name": "系统性能监控模板",
        "description": "系统性能监控看板模板",
        "category": "system_monitoring",
        "is_official": true
      }
    ]
  }
}
```

---

## 5. 数据分析API

### 5.1 获取Bot指标

**接口地址**: `GET /api/v1/bots/{bot_id}/metrics`

**查询参数**:
- start_date: 开始日期 (必填，如"2025-01-01")
- end_date: 结束日期 (必填，如"2025-01-31")
- granularity: 时间粒度 (可选，hour/day，默认day)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": 123,
    "bot_name": "客服助手",
    "metrics": [
      {
        "date": "2025-01-01",
        "conversation_count": 1520,
        "message_count": 4580,
        "active_users": 320,
        "new_users": 45,
        "avg_response_time": 1250.5,
        "p95_response_time": 2300.0,
        "p99_response_time": 3500.0,
        "error_rate": 0.012,
        "token_consumption": 1250000,
        "api_call_count": 1520,
        "resource_points_consumed": 15200,
        "satisfaction_score": 4.3,
        "feedback_count": 230
      },
      {
        "date": "2025-01-02",
        "conversation_count": 1680,
        "message_count": 5040,
        "active_users": 350,
        "new_users": 52,
        "avg_response_time": 1280.0,
        "p95_response_time": 2400.0,
        "p99_response_time": 3600.0,
        "error_rate": 0.010,
        "token_consumption": 1380000,
        "api_call_count": 1680,
        "resource_points_consumed": 16800,
        "satisfaction_score": 4.4,
        "feedback_count": 268
      }
    ],
    "summary": {
      "total_conversations": 45600,
      "total_messages": 136800,
      "avg_active_users": 350,
      "total_new_users": 1200,
      "avg_response_time": 1280.0,
      "avg_error_rate": 0.015,
      "total_token_consumption": 37500000,
      "total_resource_points": 456000,
      "avg_satisfaction_score": 4.2
    }
  }
}
```

### 5.2 获取分析报告

**接口地址**: `GET /api/v1/bots/{bot_id}/analytics`

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)
- metrics: 指标类型 (可选，逗号分隔，如"usage,performance,cost")

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "usage": {
      "conversation_trend": [
        { "date": "2025-01-01", "conversations": 1520, "messages": 4580 },
        { "date": "2025-01-02", "conversations": 1680, "messages": 5040 }
      ],
      "user_growth": [
        { "date": "2025-01-01", "new_users": 45, "active_users": 320 },
        { "date": "2025-01-02", "new_users": 52, "active_users": 350 }
      ],
      "peak_hours": [
        { "hour": 9, "conversations": 520 },
        { "hour": 14, "conversations": 680 },
        { "hour": 20, "conversations": 480 }
      ]
    },
    "performance": {
      "response_time_distribution": [
        { "range": "0-1s", "count": 32000 },
        { "range": "1-2s", "count": 10000 },
        { "range": "2-5s", "count": 3000 },
        { "range": ">5s", "count": 600 }
      ],
      "error_rate_trend": [
        { "date": "2025-01-01", "error_rate": 0.012 },
        { "date": "2025-01-02", "error_rate": 0.010 }
      ],
      "top_errors": [
        { "error": "Timeout", "count": 350, "percentage": 45.0 },
        { "error": "API Error", "count": 280, "percentage": 36.0 },
        { "error": "Rate Limit", "count": 145, "percentage": 19.0 }
      ]
    },
    "cost": {
      "token_cost_breakdown": [
        { "model": "GPT-4", "cost": 850.50, "percentage": 68.0 },
        { "model": "GPT-3.5", "cost": 320.20, "percentage": 26.0 },
        { "model": "Embedding", "cost": 78.30, "percentage": 6.0 }
      ],
      "daily_cost_trend": [
        { "date": "2025-01-01", "cost": 125.50 },
        { "date": "2025-01-02", "cost": 138.20 }
      ],
      "total_cost": 1250.50
    }
  }
}
```

### 5.3 导出报表

**接口地址**: `POST /api/v1/bots/{bot_id}/metrics/export`

**请求参数**:
```json
{
  "start_date": "2025-01-01",
  "end_date": "2025-01-31",
  "format": "excel",
  "include_charts": true,
  "metrics": ["usage", "performance", "cost", "satisfaction"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "导出任务已创建",
  "data": {
    "task_id": "export-task-123",
    "download_url": "/api/v1/downloads/export-task-123",
    "expires_at": "2025-01-02T00:00:00Z"
  }
}
```

### 5.4 获取对话记录

**接口地址**: `GET /api/v1/bots/{bot_id}/conversations`

**查询参数**:
- conversation_id: 会话ID (可选)
- user_id: 用户ID (可选)
- start_time: 开始时间 (可选)
- end_time: 结束时间 (可选)
- page: 页码 (可选，默认1)
- page_size: 每页数量 (可选，默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1000,
    "items": [
      {
        "conversation_id": "conv-abc-123",
        "user_id": 456,
        "user_nickname": "用户456",
        "message_count": 15,
        "duration": 320,
        "started_at": "2025-01-15T10:00:00Z",
        "ended_at": "2025-01-15T10:05:20Z",
        "feedback_score": 5,
        "feedback_tag": "helpful",
        "messages": [
          {
            "role": "user",
            "content_type": "text",
            "content": "你好",
            "created_at": "2025-01-15T10:00:00Z"
          },
          {
            "role": "assistant",
            "content_type": "text",
            "content": "您好,有什么我可以帮助您的吗?",
            "model_id": "gpt-4",
            "model_name": "GPT-4",
            "provider": "openai",
            "token_count": 25,
            "response_time_ms": 1200,
            "created_at": "2025-01-15T10:00:01.200Z"
          },
          {
            "role": "user",
            "content_type": "text",
            "content": "我想查询我的订单状态",
            "created_at": "2025-01-15T10:00:05.000Z"
          },
          {
            "role": "tool",
            "content_type": "text",
            "content": "{\"order_id\":\"ORD-12345\",\"status\":\"已发货\"}",
            "token_count": 15,
            "created_at": "2025-01-15T10:00:06.500Z"
          }
        ]
      }
    ]
  }
}
```

### 5.5 获取用户行为分析

**接口地址**: `GET /api/v1/bots/{bot_id}/user-behaviors`

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)
- action: 行为类型筛选 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_behaviors": 50000,
    "unique_users": 2500,
    "top_actions": [
      { "action": "bot_chat", "count": 35000, "percentage": 70.0 },
      { "action": "bot_view", "count": 10000, "percentage": 20.0 },
      { "action": "bot_share", "count": 5000, "percentage": 10.0 }
    ],
    "user_journey": [
      {
        "step": 1,
        "action": "bot_view",
        "count": 10000,
        "drop_off_rate": 0.2
      },
      {
        "step": 2,
        "action": "bot_chat",
        "count": 8000,
        "drop_off_rate": 0.1
      },
      {
        "step": 3,
        "action": "bot_share",
        "count": 7200,
        "drop_off_rate": 0.0
      }
    ],
    "retention": {
      "day_1": 0.65,
      "day_7": 0.35,
      "day_30": 0.15
    }
  }
}
```

---

## 6. 数据模型

### 6.1 Dashboard

```typescript
interface Dashboard {
  id: string;
  tenant_id: string;
  creator_id: number;
  name: string;
  description?: string;
  layout: DashboardLayout;
  components: DashboardComponent[];
  is_public: boolean;
  allowed_roles: string[];
  view_count: number;
  created_at: Date;
  updated_at: Date;
}

interface DashboardLayout {
  cols: number;
  rows: number;
}

interface DashboardComponent {
  id: string;
  type: string;
  position: ComponentPosition;
  title: string;
  dataSource: DataSource;
  style?: Record<string, any>;
}

interface ComponentPosition {
  x: number;
  y: number;
  w: number;
  h: number;
}

type DataSource = SQLDataSource | APIDataSource | WebSocketDataSource;

interface SQLDataSource {
  type: 'sql';
  query: string;
  params?: any[];
}

interface APIDataSource {
  type: 'api';
  url: string;
  method?: 'GET' | 'POST';
  headers?: Record<string, string>;
  refreshInterval?: number;
}

interface WebSocketDataSource {
  type: 'websocket';
  channel: string;
}
```

### 6.2 BotMetrics

```typescript
interface BotMetrics {
  bot_id: number;
  bot_name: string;
  metrics: MetricDataPoint[];
  summary: MetricSummary;
}

interface MetricDataPoint {
  date: string;
  hour?: number;
  conversation_count: number;
  message_count: number;
  active_users: number;
  new_users: number;
  avg_response_time: number;
  p95_response_time: number;
  p99_response_time: number;
  error_rate: number;
  token_consumption: number;
  api_call_count: number;
  resource_points_consumed: number;
  satisfaction_score: number;
  feedback_count: number;
}

interface MetricSummary {
  total_conversations: number;
  total_messages: number;
  avg_active_users: number;
  total_new_users: number;
  avg_response_time: number;
  avg_error_rate: number;
  total_token_consumption: number;
  total_resource_points: number;
  avg_satisfaction_score: number;
}
```

### 6.3 ConversationLog

```typescript
interface ConversationLog {
  conversation_id: string;
  user_id: number;
  user_nickname?: string;
  message_count: number;
  duration: number;
  started_at: Date;
  ended_at: Date;
  feedback_score?: number;
  feedback_tag?: string;
  messages: Message[];
}

interface Message {
  role: 'user' | 'assistant' | 'system' | 'tool';
  content_type: 'text' | 'image' | 'audio' | 'video' | 'file';
  content: string;
  model_id?: string;
  model_name?: string;
  provider?: string;
  token_count?: number;
  response_time_ms?: number;
  error_message?: string;
  created_at: Date;
}
```

### 6.4 UserBehavior

```typescript
interface UserBehavior {
  id: number;
  tenant_id: string;
  user_id: number;
  session_id?: string;
  action: string;
  entity_type?: string;
  entity_id?: number;
  content?: Record<string, any>;
  ip_address?: string;
  user_agent?: string;
  referrer?: string;
  device_type: 'desktop' | 'mobile' | 'tablet' | 'unknown';
  duration_ms?: number;
  created_at: Date;
}

interface UserBehaviorAnalysis {
  total_behaviors: number;
  unique_users: number;
  top_actions: ActionStat[];
  user_journey: JourneyStep[];
  retention: RetentionData;
}

interface ActionStat {
  action: string;
  count: number;
  percentage: number;
}

interface JourneyStep {
  step: number;
  action: string;
  count: number;
  drop_off_rate: number;
}

interface RetentionData {
  day_1: number;
  day_7: number;
  day_30: number;
}
```

---

## 7. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 33001 | 404 | 看板不存在 |
| 33002 | 400 | 看板名称已存在 |
| 33003 | 400 | 布局配置无效 |
| 33004 | 400 | 组件配置无效 |
| 33005 | 400 | 数据源配置无效 |
| 33006 | 400 | SQL查询错误 |
| 33007 | 403 | 无权限访问看板 |
| 33101 | 404 | 组件类型不存在 |
| 33102 | 400 | 组件类型已停用 |
| 33201 | 400 | 日期范围无效 |
| 33202 | 400 | 时间粒度无效 |
| 33203 | 404 | Bot不存在 |
| 33204 | 500 | 数据导出失败 |
| 33205 | 404 | 对话记录不存在 |

---

**文档结束**
