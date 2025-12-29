# API接口文档：WorkflowEditor工作流编辑器模块

**模块名称**: WorkflowEditor (WorkflowEditor)
**设计文档**: 14-百应开发平台_开发工作台_WorkflowEditor.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 工作流管理API](#2-工作流管理api)
- [3. 节点管理API](#3-节点管理api)
- [4. 工作流执行API](#4-工作流执行api)
- [5. 调试与日志API](#5-调试与日志api)
- [6. 模板管理API](#6-模板管理api)
- [7. 数据模型](#7-数据模型)
- [8. 错误码定义](#8-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

WorkflowEditor是ZKER的可视化工作流编辑器，通过100%复用zker Workflow系统，实现复杂业务流程的可视化编排：

- ✅ **可视化编排** - 拖拽式节点编辑，可视化连接配置
- ✅ **多种节点类型** - 触发器、技能、逻辑、数据节点
- ✅ **调试执行** - 单步调试、断点调试、执行日志
- ✅ **版本管理** - 工作流版本控制、回滚
- ✅ **模板库** - 预置模板、自定义模板

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 工作流引擎: zker Workflow（100%复用）
- 数据库: MySQL 8.4.5

**前端技术栈**:
- 框架: React 18 + TypeScript
- 画布库: React Flow / X6
- UI库: Semi Design

### 1.3 实现策略

| 实施项 | 实现方式 | 占比 |
|-------|---------|------|
| **核心功能** | 🔁 100%复用zker Workflow | 0% |
| **企业扩展** | 📊 配置扩展 | 10% |

---

## 2. 工作流管理API

### 2.1 获取工作流列表

**接口地址**: `GET /api/v1/workflows`

**功能说明**: 获取当前租户的所有工作流

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:
- page: 页码，默认1
- page_size: 每页数量，默认20
- status: 状态过滤 (draft/active/archived)
- keyword: 搜索关键词

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "id": "wf-001",
        "name": "客户信息同步流程",
        "description": "从CRM同步客户信息到知识库",
        "status": "active",
        "version": "3",
        "nodesCount": 8,
        "executionsCount": 1250,
        "lastExecutedAt": "2025-01-03T09:30:00Z",
        "createdAt": "2025-01-01T10:00:00Z",
        "updatedAt": "2025-01-03T10:00:00Z"
      },
      {
        "id": "wf-002",
        "name": "智能客服路由",
        "description": "根据用户问题路由到不同技能组",
        "status": "active",
        "version": "1",
        "nodesCount": 5,
        "executionsCount": 3500,
        "lastExecutedAt": "2025-01-03T10:15:00Z",
        "createdAt": "2025-01-02T10:00:00Z",
        "updatedAt": "2025-01-02T10:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "trace-abc-123"
}
```

### 2.2 创建工作流

**接口地址**: `POST /api/v1/workflows`

**功能说明**: 创建新工作流

**请求参数**:
```json
{
  "name": "客户信息同步流程",
  "description": "从CRM同步客户信息到知识库",
  "category": "integration",
  "nodes": [
    {
      "id": "node-001",
      "type": "trigger",
      "subtype": "http",
      "name": "HTTP触发器",
      "config": {
        "path": "/webhook/crm-sync",
        "method": "POST"
      },
      "position": {
        "x": 100,
        "y": 100
      }
    },
    {
      "id": "node-002",
      "type": "action",
      "subtype": "bot",
      "name": "调用Bot",
      "config": {
        "bot_id": "bot-001",
        "input_mapping": {
          "question": "{{trigger.body.question}}"
        }
      },
      "position": {
        "x": 300,
        "y": 100
      }
    }
  ],
  "edges": [
    {
      "id": "edge-001",
      "source": "node-001",
      "target": "node-002",
      "sourceHandle": "output",
      "targetHandle": "input"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "工作流创建成功",
  "data": {
    "id": "wf-001",
    "name": "客户信息同步流程",
    "description": "从CRM同步客户信息到知识库",
    "category": "integration",
    "status": "draft",
    "version": 1,
    "nodes": [...],
    "edges": [...],
    "createdAt": "2025-01-03T10:00:00Z"
  }
}
```

### 2.3 获取工作流详情

**接口地址**: `GET /api/v1/workflows/{workflow_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "wf-001",
    "name": "客户信息同步流程",
    "description": "从CRM同步客户信息到知识库",
    "category": "integration",
    "status": "active",
    "version": 3,
    "nodes": [
      {
        "id": "node-001",
        "type": "trigger",
        "subtype": "http",
        "name": "HTTP触发器",
        "config": {
          "path": "/webhook/crm-sync",
          "method": "POST"
        },
        "position": {
          "x": 100,
          "y": 100
        }
      }
    ],
    "edges": [
      {
        "id": "edge-001",
        "source": "node-001",
        "target": "node-002"
      }
    ],
    "executionsCount": 1250,
    "lastExecutedAt": "2025-01-03T09:30:00Z",
    "createdAt": "2025-01-01T10:00:00Z",
    "updatedAt": "2025-01-03T10:00:00Z"
  }
}
```

### 2.4 更新工作流

**接口地址**: `PUT /api/v1/workflows/{workflow_id}`

**请求参数**: 同创建工作流

**响应示例**:
```json
{
  "code": 0,
  "message": "工作流更新成功",
  "data": {
    "id": "wf-001",
    "version": 4,
    "updatedAt": "2025-01-03T11:00:00Z"
  }
}
```

### 2.5 删除工作流

**接口地址**: `DELETE /api/v1/workflows/{workflow_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "工作流删除成功"
}
```

### 2.6 发布/取消发布工作流

**接口地址**: `POST /api/v1/workflows/{workflow_id}/publish`

**请求参数**:
```json
{
  "action": "publish"
}
```

**action枚举**: publish（发布）/ unpublish（取消发布）

**响应示例**:
```json
{
  "code": 0,
  "message": "工作流发布成功",
  "data": {
    "id": "wf-001",
    "status": "active",
    "publishedAt": "2025-01-03T11:00:00Z"
  }
}
```

---

## 3. 节点管理API

### 3.1 获取节点类型列表

**接口地址**: `GET /api/v1/workflows/node-types`

**功能说明**: 获取所有可用的节点类型

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "categories": [
      {
        "name": "触发器",
        "types": [
          {
            "type": "trigger",
            "subtype": "http",
            "name": "HTTP触发器",
            "description": "通过HTTP请求触发工作流",
            "icon": "webhook",
            "configSchema": {
              "type": "object",
              "properties": {
                "path": {
                  "type": "string",
                  "title": "请求路径"
                },
                "method": {
                  "type": "string",
                  "enum": ["GET", "POST", "PUT", "DELETE"],
                  "title": "HTTP方法"
                }
              }
            }
          },
          {
            "type": "trigger",
            "subtype": "timer",
            "name": "定时器",
            "description": "按计划定时触发工作流",
            "icon": "clock",
            "configSchema": {...}
          }
        ]
      },
      {
        "name": "技能",
        "types": [
          {
            "type": "action",
            "subtype": "bot",
            "name": "调用Bot",
            "description": "调用指定的Bot处理请求",
            "icon": "robot"
          },
          {
            "type": "action",
            "subtype": "plugin",
            "name": "调用插件",
            "description": "调用插件执行特定功能",
            "icon": "plugin"
          }
        ]
      },
      {
        "name": "逻辑",
        "types": [
          {
            "type": "logic",
            "subtype": "condition",
            "name": "条件判断",
            "description": "根据条件执行不同分支",
            "icon": "branch"
          },
          {
            "type": "logic",
            "subtype": "loop",
            "name": "循环",
            "description": "重复执行一组节点",
            "icon": "loop"
          }
        ]
      },
      {
        "name": "数据",
        "types": [
          {
            "type": "data",
            "subtype": "database",
            "name": "数据库操作",
            "description": "执行数据库查询或更新",
            "icon": "database"
          },
          {
            "type": "data",
            "subtype": "api",
            "name": "API调用",
            "description": "调用外部API接口",
            "icon": "api"
          }
        ]
      }
    ]
  }
}
```

### 3.2 添加节点

**接口地址**: `POST /api/v1/workflows/{workflow_id}/nodes`

**请求参数**:
```json
{
  "type": "action",
  "subtype": "bot",
  "name": "调用邮件Bot",
  "config": {
    "bot_id": "bot-email-assistant",
    "input_mapping": {
      "question": "{{trigger.body.question}}"
    }
  },
  "position": {
    "x": 300,
    "y": 100
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "节点添加成功",
  "data": {
    "id": "node-003",
    "type": "action",
    "subtype": "bot",
    "name": "调用邮件Bot",
    "position": {
      "x": 300,
      "y": 100
    },
    "createdAt": "2025-01-03T11:00:00Z"
  }
}
```

### 3.3 更新节点

**接口地址**: `PUT /api/v1/workflows/{workflow_id}/nodes/{node_id}`

### 3.4 删除节点

**接口地址**: `DELETE /api/v1/workflows/{workflow_id}/nodes/{node_id}`

---

## 4. 工作流执行API

### 4.1 执行工作流

**接口地址**: `POST /api/v1/workflows/{workflow_id}/execute`

**功能说明**: 手动触发执行工作流

**请求参数**:
```json
{
  "input": {
    "customer_id": "cust-001",
    "action": "sync"
  },
  "options": {
    "debug_mode": true,
    "async": false
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "工作流执行成功",
  "data": {
    "executionId": "exec-001",
    "workflowId": "wf-001",
    "status": "success",
    "input": {...},
    "output": {
      "result": "同步成功",
      "synced_records": 25
    },
    "nodesExecuted": 5,
    "duration": 1250,
    "startedAt": "2025-01-03T11:00:00Z",
    "completedAt": "2025-01-03T11:00:01.25Z"
  }
}
```

### 4.2 获取执行记录列表

**接口地址**: `GET /api/v1/workflows/{workflow_id}/executions`

**查询参数**:
- page: 页码
- page_size: 每页数量
- status: 状态过滤 (success/failed/running)
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1250,
    "items": [
      {
        "id": "exec-001",
        "workflowId": "wf-001",
        "status": "success",
        "input": {...},
        "output": {
          "result": "同步成功",
          "synced_records": 25
        },
        "nodesExecuted": 5,
        "duration": 1250,
        "startedAt": "2025-01-03T11:00:00Z",
        "completedAt": "2025-01-03T11:00:01.25Z"
      },
      {
        "id": "exec-002",
        "workflowId": "wf-001",
        "status": "failed",
        "error": "节点执行超时",
        "nodesExecuted": 3,
        "duration": 30000,
        "startedAt": "2025-01-03T10:55:00Z",
        "completedAt": "2025-01-03T10:55:30Z"
      }
    ]
  }
}
```

### 4.3 获取执行详情

**接口地址**: `GET /api/v1/workflows/executions/{execution_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "exec-001",
    "workflowId": "wf-001",
    "status": "success",
    "input": {
      "customer_id": "cust-001",
      "action": "sync"
    },
    "output": {
      "result": "同步成功",
      "synced_records": 25
    },
    "nodeExecutions": [
      {
        "nodeId": "node-001",
        "nodeName": "HTTP触发器",
        "status": "success",
        "input": {...},
        "output": {...},
        "startedAt": "2025-01-03T11:00:00Z",
        "completedAt": "2025-01-03T11:00:00.05Z"
      },
      {
        "nodeId": "node-002",
        "nodeName": "调用Bot",
        "status": "success",
        "input": {...},
        "output": {...},
        "startedAt": "2025-01-03T11:00:00.05Z",
        "completedAt": "2025-01-03T11:00:01Z"
      }
    ],
    "nodesExecuted": 5,
    "duration": 1250,
    "startedAt": "2025-01-03T11:00:00Z",
    "completedAt": "2025-01-03T11:00:01.25Z"
  }
}
```

---

## 5. 调试与日志API

### 5.1 单步调试

**接口地址**: `POST /api/v1/workflows/{workflow_id}/debug`

**功能说明**: 进入调试模式，单步执行工作流

**请求参数**:
```json
{
  "input": {
    "customer_id": "cust-001"
  },
  "breakpoints": ["node-002"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "进入调试模式",
  "data": {
    "debugSessionId": "debug-001",
    "workflowId": "wf-001",
    "currentNode": {
      "id": "node-001",
      "name": "HTTP触发器",
      "status": "completed",
      "output": {...}
    },
    "nextNode": {
      "id": "node-002",
      "name": "调用Bot",
      "status": "waiting",
      "input": {...}
    },
    "variables": {
      "customer_id": "cust-001",
      "action": "sync"
    }
  }
}
```

### 5.2 继续执行下一步

**接口地址**: `POST /api/v1/workflows/debug/{debug_session_id}/step`

**响应示例**:
```json
{
  "code": 0,
  "message": "执行下一步",
  "data": {
    "debugSessionId": "debug-001",
    "currentNode": {
      "id": "node-002",
      "name": "调用Bot",
      "status": "running"
    }
  }
}
```

### 5.3 停止调试

**接口地址**: `POST /api/v1/workflows/debug/{debug_session_id}/stop`

**响应示例**:
```json
{
  "code": 0,
  "message": "调试已停止"
}
```

---

## 6. 模板管理API

### 6.1 获取模板列表

**接口地址**: `GET /api/v1/workflows/templates`

**查询参数**:
- category: 分类过滤
- keyword: 搜索关键词

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "tpl-001",
        "name": "HTTP API模板",
        "description": "快速创建HTTP API工作流",
        "category": "integration",
        "thumbnail": "https://cdn.example.com/templates/http-api.png",
        "nodes": [
          {
            "type": "trigger",
            "subtype": "http",
            "name": "HTTP触发器"
          },
          {
            "type": "action",
            "subtype": "bot",
            "name": "调用Bot"
          }
        ],
        "usageCount": 350
      },
      {
        "id": "tpl-002",
        "name": "定时任务模板",
        "description": "创建定时执行的工作流",
        "category": "automation",
        "thumbnail": "https://cdn.example.com/templates/scheduler.png",
        "nodes": [
          {
            "type": "trigger",
            "subtype": "timer",
            "name": "定时器"
          }
        ],
        "usageCount": 220
      }
    ],
    "total": 15
  }
}
```

### 6.2 从模板创建工作流

**接口地址**: `POST /api/v1/workflows/templates/{template_id}/create`

**请求参数**:
```json
{
  "name": "我的定时任务",
  "description": "每天早上9点执行数据同步"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "工作流创建成功",
  "data": {
    "id": "wf-003",
    "name": "我的定时任务",
    "description": "每天早上9点执行数据同步",
    "status": "draft",
    "nodes": [...],
    "createdAt": "2025-01-03T11:00:00Z"
  }
}
```

---

## 7. 数据模型

### 7.1 Workflow（工作流）
```typescript
interface Workflow {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  category?: string;
  status: 'draft' | 'active' | 'archived';
  version: number;
  nodes: WorkflowNode[];
  edges: WorkflowEdge[];
  executions_count: number;
  last_executed_at?: string;
  created_at: string;
  updated_at: string;
}

interface WorkflowNode {
  id: string;
  type: 'trigger' | 'action' | 'logic' | 'data';
  subtype: string;
  name: string;
  config: Record<string, any>;
  position: {
    x: number;
    y: number;
  };
}

interface WorkflowEdge {
  id: string;
  source: string;
  target: string;
  sourceHandle?: string;
  targetHandle?: string;
  condition?: string;
}
```

### 7.2 WorkflowExecution（执行记录）
```typescript
interface WorkflowExecution {
  id: string;
  workflow_id: string;
  status: 'success' | 'failed' | 'running';
  input: Record<string, any>;
  output?: Record<string, any>;
  error?: string;
  nodes_executed: number;
  duration: number;
  started_at: string;
  completed_at?: string;
  node_executions: NodeExecution[];
}

interface NodeExecution {
  node_id: string;
  node_name: string;
  status: 'success' | 'failed' | 'running';
  input?: Record<string, any>;
  output?: Record<string, any>;
  error?: string;
  started_at: string;
  completed_at?: string;
}
```

---

## 8. 错误码定义

| 错误码 | HTTP状态码 | 说明 | 处理建议 |
|--------|-----------|------|----------|
| 60001 | 404 | 工作流不存在 | 检查工作流ID |
| 60002 | 400 | 工作流定义无效 | 检查节点配置是否正确 |
| 60003 | 400 | 工作流已发布 | 请先取消发布再修改 |
| 60004 | 400 | 节点配置错误 | 检查节点参数配置 |
| 60005 | 400 | 节点连接循环 | 工作流不能包含循环依赖 |
| 60006 | 400 | 工作流正在执行 | 请等待执行完成 |
| 60101 | 404 | 执行记录不存在 | 检查执行ID |
| 60102 | 400 | 执行超时 | 优化工作流或增加超时时间 |
| 60103 | 400 | 节点执行失败 | 检查节点配置和输入数据 |
| 60201 | 404 | 模板不存在 | 检查模板ID |
| 60202 | 400 | 模板格式错误 | 联系管理员 |

---

**文档结束**
