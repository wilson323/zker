# API接口文档：智能路由引擎模块

**模块名称**: 智能路由引擎 (IntelligentRoutingEngine)
**设计文档**: 16-五大AI引擎核心_智能路由引擎.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P1

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 路由决策API](#2-路由决策api)
- [3. 意图识别API](#3-意图识别api)
- [4. 路由规则管理API](#4-路由规则管理api)
- [5. 负载均衡配置API](#5-负载均衡配置api)
- [6. 学习优化API](#6-学习优化api)
- [7. 监控与分析API](#7-监控与分析api)
- [8. 反馈收集API](#8-反馈收集api)
- [9. 数据模型](#9-数据模型)
- [10. 错误码定义](#10-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

智能路由引擎是ZKER企业级SaaS平台的"大脑中枢"，负责AI助手的智能路由：

- ✅ **意图识别** - NLP意图分类、实体抽取、关键词提取
- ✅ **路由决策** - 多因素综合决策，选择最优服务
- ✅ **负载均衡** - 5种负载均衡策略
- ✅ **学习优化** - 基于反馈数据持续优化
- ✅ **监控分析** - 9大监控指标、实时告警

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 意图识别: 通义千问/GPT-4
- 实体抽取: spaCy/HanLP
- 规则引擎: Drools/Govaluate
- 向量检索: Milvus
- 缓存: Redis 8.0
- 消息队列: NSQ
- 分析数据库: ClickHouse

### 1.3 性能指标

| 指标 | 目标值 |
|------|--------|
| 路由准确率 | ≥ 95% |
| 路由决策耗时 | < 200ms |
| 系统可用性 | ≥ 99.9% |
| 并发处理能力 | 10000 QPS |

---

## 2. 路由决策API

### 2.1 执行路由

**接口地址**: `POST /api/v1/routing/route`

**功能说明**: 根据用户输入智能路由到最优服务

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "input": "XX产品的价格是多少？",
  "context": {
    "conversation_id": "conv-001",
    "history": [
      {
        "role": "user",
        "content": "你好"
      },
      {
        "role": "assistant",
        "content": "您好，有什么可以帮您？"
      }
    ]
  },
  "options": {
    "enable_cache": true,
    "timeout_ms": 5000
  }
}
```

**参数说明**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| input | string | 是 | 用户输入 |
| context | object | 否 | 上下文信息 |
| options | object | 否 | 路由选项 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "intent": {
      "type": "knowledge_query",
      "confidence": 0.95,
      "entities": {
        "product": "XX产品"
      },
      "metadata": {
        "recognized_by": "llm",
        "model": "qwen-turbo"
      }
    },
    "entities": {
      "product": "XX产品",
      "intent_keywords": ["价格", "多少"]
    },
    "candidates": [
      {
        "service_id": "bot-001",
        "service_name": "产品知识库Bot",
        "service_type": "bot",
        "score": 0.92,
        "reasons": [
          "用户匹配度: 0.85",
          "历史成功率: 0.90",
          "负载评分: 0.95"
        ]
      },
      {
        "service_id": "bot-002",
        "service_name": "通用客服Bot",
        "service_type": "bot",
        "score": 0.78,
        "reasons": []
      }
    ],
    "selected": {
      "service_id": "bot-001",
      "service_name": "产品知识库Bot",
      "service_type": "bot",
      "endpoint": "https://api.example.com/bots/bot-001/chat"
    },
    "confidence": 0.92,
    "reason": "根据意图识别（知识查询）和用户画像（销售岗位），选择产品知识库Bot，历史成功率90%",
    "routing_time_ms": 125
  }
}
```

### 2.2 批量路由

**接口地址**: `POST /api/v1/routing/route/batch`

**功能说明**: 批量路由多个请求

**请求参数**:
```json
{
  "requests": [
    {
      "input": "产品价格是多少？",
      "user_id": "user-001"
    },
    {
      "input": "创建销售订单",
      "user_id": "user-002"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "results": [
      {
        "request_id": 0,
        "selected_service_id": "bot-001",
        "intent": "knowledge_query"
      },
      {
        "request_id": 1,
        "selected_service_id": "workflow-001",
        "intent": "task_execution"
      }
    ]
  }
}
```

---

## 3. 意图识别API

### 3.1 识别意图

**接口地址**: `POST /api/v1/routing/intent/recognize`

**功能说明**: 单独识别用户输入的意图

**请求参数**:
```json
{
  "input": "帮我创建一个销售订单",
  "options": {
    "enable_llm": true,
    "enable_rule": true,
    "enable_vector": true,
    "confidence_threshold": 0.7
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "intent": {
      "type": "task_execution",
      "confidence": 0.92,
      "entities": {
        "action": "创建",
        "target": "销售订单"
      },
      "metadata": {
        "recognized_by": "llm",
        "model": "gpt-4",
        "processing_time_ms": 85
      }
    },
    "alternative_intents": [
      {
        "type": "knowledge_query",
        "confidence": 0.15
      }
    ]
  }
}
```

### 3.2 提取实体

**接口地址**: `POST /api/v1/routing/intent/extract-entities`

**请求参数**:
```json
{
  "input": "张三明天要去上海分公司开会",
  "entity_types": ["person", "location", "organization", "time"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "entities": {
      "person": "张三",
      "location": "上海",
      "organization": "上海分公司",
      "time": "明天"
    },
    "spans": [
      {
        "text": "张三",
        "label": "PERSON",
        "start": 0,
        "end": 2
      },
      {
        "text": "明天",
        "label": "TIME",
        "start": 3,
        "end": 5
      }
    ]
  }
}
```

### 3.3 获取意图列表

**接口地址**: `GET /api/v1/routing/intent/types`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "intents": [
      {
        "type": "knowledge_query",
        "name": "知识查询",
        "description": "用户询问知识库中的信息",
        "examples": [
          "XX产品的价格是多少？",
          "如何配置API密钥？"
        ]
      },
      {
        "type": "task_execution",
        "name": "任务执行",
        "description": "用户请求执行特定操作",
        "examples": [
          "帮我创建一个销售订单",
          "发送邮件给客户"
        ]
      },
      {
        "type": "data_analysis",
        "name": "数据分析",
        "description": "用户请求分析数据",
        "examples": [
          "上个月销售额",
          "统计用户增长趋势"
        ]
      },
      {
        "type": "content_generation",
        "name": "内容生成",
        "description": "用户请求生成内容",
        "examples": [
          "写一份产品周报",
          "生成营销文案"
        ]
      }
    ]
  }
}
```

---

## 4. 路由规则管理API

### 4.1 创建路由规则

**接口地址**: `POST /api/v1/routing/rules`

**请求参数**:
```json
{
  "name": "销售部门优先路由到CRM Bot",
  "description": "销售部门的用户优先路由到CRM Bot",
  "priority": 1,
  "enabled": true,
  "condition": {
    "intent": ["knowledge_query", "task_execution"],
    "keyword": ["客户", "订单", "销售"],
    "user_roles": ["sales", "sales_manager"],
    "departments": ["销售部"],
    "time_range": {
      "start": "09:00",
      "end": "18:00"
    },
    "custom_expr": "user.profile.department == '销售部' && input.contains('订单')"
  },
  "action": {
    "target_service_id": "bot-crm-001",
    "filter_type": "service_type",
    "filter_value": "bot",
    "priority_boost": 10
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "路由规则创建成功",
  "data": {
    "rule_id": "rule-001",
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

### 4.2 获取路由规则列表

**接口地址**: `GET /api/v1/routing/rules`

**查询参数**:
- page: 页码
- page_size: 每页数量
- enabled: 是否启用
- priority: 优先级过滤

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": "rule-001",
        "name": "销售部门优先路由到CRM Bot",
        "priority": 1,
        "enabled": true,
        "hit_count": 1500,
        "last_hit_at": "2025-01-03T10:00:00Z",
        "created_at": "2025-01-01T10:00:00Z"
      }
    ]
  }
}
```

### 4.3 更新路由规则

**接口地址**: `PUT /api/v1/routing/rules/{rule_id}`

### 4.4 删除路由规则

**接口地址**: `DELETE /api/v1/routing/rules/{rule_id}`

### 4.5 启用/禁用路由规则

**接口地址**: `PUT /api/v1/routing/rules/{rule_id}/status`

**请求参数**:
```json
{
  "enabled": false
}
```

### 4.6 测试路由规则

**接口地址**: `POST /api/v1/routing/rules/test`

**功能说明**: 测试路由规则是否匹配

**请求参数**:
```json
{
  "rule": {
    "condition": {
      "intent": ["knowledge_query"],
      "departments": ["销售部"]
    }
  },
  "test_case": {
    "input": "查询客户信息",
    "intent": "knowledge_query",
    "user_profile": {
      "department": "销售部"
    }
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "matched": true,
    "match_details": {
      "intent_matched": true,
      "department_matched": true
    }
  }
}
```

---

## 5. 负载均衡配置API

### 5.1 获取负载均衡策略

**接口地址**: `GET /api/v1/routing/load-balance/strategy`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "strategy": "weighted_round_robin",
    "config": {
      "weights": {
        "bot-001": 10,
        "bot-002": 5,
        "bot-003": 3
      }
    }
  }
}
```

### 5.2 更新负载均衡策略

**接口地址**: `PUT /api/v1/routing/load-balance/strategy`

**请求参数**:
```json
{
  "strategy": "response_time_priority",
  "config": {
    "enable_sticky_session": true,
    "sticky_session_ttl": 3600
  }
}
```

**strategy枚举**:
- `round_robin`: 轮询
- `weighted_round_robin`: 加权轮询
- `least_connections`: 最少连接
- `response_time_priority`: 响应时间优先
- `sticky_session`: 粘性会话

### 5.3 获取服务实时状态

**接口地址**: `GET /api/v1/routing/services/status`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "services": [
      {
        "service_id": "bot-001",
        "service_name": "产品知识库Bot",
        "status": "online",
        "current_connections": 25,
        "current_load": 0.35,
        "avg_response_time_ms": 150,
        "error_rate": 0.01,
        "availability": 0.99
      },
      {
        "service_id": "bot-002",
        "service_name": "通用客服Bot",
        "status": "online",
        "current_connections": 10,
        "current_load": 0.15,
        "avg_response_time_ms": 120,
        "error_rate": 0.005,
        "availability": 0.995
      }
    ]
  }
}
```

---

## 6. 学习优化API

### 6.1 触发模型优化

**接口地址**: `POST /api/v1/routing/learning/optimize-model`

**功能说明**: 触发意图识别模型的优化训练

**请求参数**:
```json
{
  "optimization_type": "prompt_optimization",
  "target_intents": ["knowledge_query"],
  "time_range_days": 7
}
```

**optimization_type枚举**:
- `prompt_optimization`: 提示词优化
- `add_training_samples`: 增加训练样本
- `fine_tune`: 模型微调

**响应示例**:
```json
{
  "code": 0,
  "message": "模型优化任务已提交",
  "data": {
    "task_id": "opt-task-001",
    "status": "pending",
    "estimated_time_minutes": 15
  }
}
```

### 6.2 获取优化任务状态

**接口地址**: `GET /api/v1/routing/learning/optimize-model/{task_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "opt-task-001",
    "status": "completed",
    "progress": 100,
    "result": {
      "before_accuracy": 0.88,
      "after_accuracy": 0.94,
      "improvement": 0.06
    },
    "started_at": "2025-01-03T10:00:00Z",
    "completed_at": "2025-01-03T10:15:00Z"
  }
}
```

### 6.3 创建A/B测试

**接口地址**: `POST /api/v1/routing/learning/ab-tests`

**请求参数**:
```json
{
  "name": "新提示词效果测试",
  "description": "测试优化后的提示词效果",
  "control_config": {
    "strategy": "current_prompt"
  },
  "treatment_config": {
    "strategy": "new_prompt",
    "prompt_version": "v2.0"
  },
  "traffic_split": 50,
  "min_sample_size": 500
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "A/B测试创建成功",
  "data": {
    "test_id": "ab-test-001",
    "status": "running",
    "started_at": "2025-01-03T10:00:00Z"
  }
}
```

---

## 7. 监控与分析API

### 7.1 获取路由指标

**接口地址**: `GET /api/v1/routing/metrics`

**查询参数**:
- time_range: 时间范围（1h/24h/7d/30d）
- tenant_id: 租户ID（可选）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "time_range": "24h",
    "metrics": {
      "intent_accuracy": 0.94,
      "intent_confidence_avg": 0.89,
      "routing_success_rate": 0.96,
      "avg_routing_time_ms": 125,
      "service_availability": 0.995,
      "avg_response_time_ms": 350,
      "error_rate": 0.008,
      "user_satisfaction_avg": 4.5,
      "positive_feedback_rate": 0.92,
      "load_balance_efficiency": 0.88
    }
  }
}
```

### 7.2 获取意图分布

**接口地址**: `GET /api/v1/routing/analytics/intent-distribution`

**查询参数**:
- time_range: 时间范围
- group_by: 分组方式（intent/hour/day）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "time_range": "24h",
    "distribution": [
      {
        "intent": "knowledge_query",
        "count": 3500,
        "percentage": 45.0,
        "avg_confidence": 0.92
      },
      {
        "intent": "task_execution",
        "count": 2800,
        "percentage": 36.0,
        "avg_confidence": 0.89
      },
      {
        "intent": "data_analysis",
        "count": 950,
        "percentage": 12.2,
        "avg_confidence": 0.94
      }
    ]
  }
}
```

### 7.3 获取服务路由统计

**接口地址**: `GET /api/v1/routing/analytics/service-routing`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "services": [
      {
        "service_id": "bot-001",
        "service_name": "产品知识库Bot",
        "routing_count": 5000,
        "success_count": 4800,
        "success_rate": 0.96,
        "avg_response_time_ms": 150,
        "user_satisfaction_avg": 4.6
      }
    ]
  }
}
```

### 7.4 获取告警列表

**接口地址**: `GET /api/v1/routing/alerts`

**查询参数**:
- status: 状态过滤（active/resolved）
- severity: 严重程度（critical/warning/info）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "alerts": [
      {
        "id": "alert-001",
        "type": "low_routing_success",
        "severity": "critical",
        "message": "路由成功率过低: 92.50% (阈值: 95%)",
        "tenant_id": "tenant-001",
        "status": "active",
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

---

## 8. 反馈收集API

### 8.1 提交路由反馈

**接口地址**: `POST /api/v1/routing/feedback`

**请求参数**:
```json
{
  "session_id": "session-001",
  "input": "产品价格是多少？",
  "intent": "knowledge_query",
  "selected_service_id": "bot-001",
  "success": true,
  "user_rating": 5,
  "response_time_ms": 150,
  "comment": "回答很准确"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "反馈提交成功"
}
```

### 8.2 批量提交反馈

**接口地址**: `POST /api/v1/routing/feedback/batch`

**请求参数**:
```json
{
  "feedbacks": [
    {
      "session_id": "session-001",
      "input": "产品价格",
      "intent": "knowledge_query",
      "selected_service_id": "bot-001",
      "success": true,
      "user_rating": 5
    },
    {
      "session_id": "session-002",
      "input": "创建订单",
      "intent": "task_execution",
      "selected_service_id": "workflow-001",
      "success": false,
      "user_rating": 2
    }
  ]
}
```

### 8.3 获取反馈统计

**接口地址**: `GET /api/v1/routing/feedback/statistics`

**查询参数**:
- time_range: 时间范围
- service_id: 服务过滤

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_feedbacks": 5000,
    "success_rate": 0.94,
    "avg_rating": 4.3,
    "rating_distribution": {
      "5": 3000,
      "4": 1200,
      "3": 500,
      "2": 200,
      "1": 100
    }
  }
}
```

---

## 9. 数据模型

### 9.1 RouteRequest（路由请求）
```typescript
interface RouteRequest {
  tenant_id: string;
  user_id: string;
  input: string;
  context?: {
    conversation_id?: string;
    history?: Array<{
      role: 'user' | 'assistant';
      content: string;
    }>;
  };
  options?: {
    enable_cache?: boolean;
    timeout_ms?: number;
  };
}
```

### 9.2 RouteResponse（路由响应）
```typescript
interface RouteResponse {
  intent: Intent;
  entities: Record<string, string>;
  candidates: Candidate[];
  selected: SelectedService;
  confidence: number;
  reason: string;
  routing_time_ms: number;
}
```

### 9.3 Intent（意图）
```typescript
interface Intent {
  type: string;
  confidence: number;
  entities: Record<string, any>;
  metadata?: {
    recognized_by?: string;
    model?: string;
    processing_time_ms?: number;
  };
}
```

### 9.4 Candidate（候选服务）
```typescript
interface Candidate {
  service_id: string;
  service_name: string;
  service_type: 'bot' | 'workflow' | 'plugin';
  score: number;
  reasons: string[];
}
```

### 9.5 RoutingRule（路由规则）
```typescript
interface RoutingRule {
  id: string;
  name: string;
  description?: string;
  priority: number;
  enabled: boolean;
  condition: {
    intent?: string[];
    keyword?: string[];
    user_roles?: string[];
    departments?: string[];
    time_range?: {
      start: string;
      end: string;
    };
    custom_expr?: string;
  };
  action: {
    target_service_id?: string;
    filter_type?: string;
    filter_value?: string;
    priority_boost?: number;
  };
  hit_count: number;
  last_hit_at?: string;
  created_at: string;
}
```

### 9.6 RoutingMetrics（路由指标）
```typescript
interface RoutingMetrics {
  intent_accuracy: number;
  intent_confidence_avg: number;
  routing_success_rate: number;
  avg_routing_time_ms: number;
  service_availability: number;
  avg_response_time_ms: number;
  error_rate: number;
  user_satisfaction_avg: number;
  positive_feedback_rate: number;
  load_balance_efficiency: number;
}
```

---

## 10. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 60001 | 400 | 用户输入为空 |
| 60002 | 400 | 输入超过最大长度限制 |
| 60003 | 500 | 意图识别失败 |
| 60004 | 500 | 实体抽取失败 |
| 60101 | 404 | 未找到可用的路由服务 |
| 60102 | 500 | 路由决策超时 |
| 60103 | 500 | 路由决策失败 |
| 60201 | 400 | 路由规则不存在 |
| 60202 | 400 | 路由规则配置无效 |
| 60203 | 400 | 路由规则优先级冲突 |
| 60301 | 400 | 负载均衡策略不支持 |
| 60302 | 500 | 服务状态检查失败 |
| 60401 | 400 | 优化任务不存在 |
| 60402 | 400 | 优化任务正在执行中 |
| 60403 | 500 | 优化任务执行失败 |

---

**文档结束**
