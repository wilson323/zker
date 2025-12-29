# API接口文档：AI Agent监控与评估平台

**模块名称**: AI Agent Monitoring & Evaluation Platform (AgentMonitor)
**设计文档**: 26-五大AI引擎核心_Agent监控与评估平台.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P3

**对标产品**: [LangSmith](https://www.ibm.com/cn-zh/think/topics/langsmith)

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. Agent监控API](#2-agent监控api)
- [3. 评估API](#3-评估api)
- [4. 提示词版本管理API](#4-提示词版本管理api)
- [5. A/B测试API](#5-ab测试api)
- [6. 成本分析API](#6-成本分析api)
- [7. 数据模型](#7-数据模型)
- [8. 错误码定义](#8-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

AI Agent监控与评估平台是对标**LangSmith**的生产级AI Agent应用监控与评估系统，提供**全流程可观测性 + 质量评估 + 成本分析**能力：

- ✅ **全流程可观测性** - 可视化追踪树，清晰展现每一步执行逻辑
- ✅ **实时监控** - 性能追踪、错误追踪、用户行为分析
- ✅ **质量评估** - 12种评估方法（精确匹配、基于事实问答等）
- ✅ **版本管理** - 提示词版本管理、A/B测试、效果对比
- ✅ **成本分析** - 实时成本追踪、成本预测、优化建议

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 时序数据库: InfluxDB (监控指标)
- 实现: 40% 编码（监控框架） + 60% 配置（评估规则）

**前端技术栈**:
- React 18 + TypeScript 5.6
- 可视化: D3.js / ECharts (追踪树)
- Semi Design UI组件库

**实现策略**: ✅ 40% 编码 + 60% 配置

### 1.3 数据库表

| 表名 | 说明 |
|------|------|
| `agent_executions` | Agent执行追踪表 |
| `agent_execution_steps` | Agent执行步骤表 |
| `user_feedback` | 用户反馈表 |
| `prompt_versions` | 提示词版本表 |
| `ab_tests` | A/B测试表 |
| `cost_analytics` | 成本分析表 |

---

## 2. Agent监控API

### 2.1 获取执行列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/executions`

**查询参数**:
- bot_id: Bot ID (可选)
- user_id: 用户ID (可选)
- status: 状态筛选 (可选，pending/running/success/failed/timeout)
- start_date: 开始日期 (可选)
- end_date: 结束日期 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1000,
    "items": [
      {
        "id": "exec-001",
        "tenant_id": "tenant-001",
        "bot_id": "bot-123",
        "bot_name": "客服助手",
        "user_id": 101,
        "user_name": "张三",
        "conversation_id": "conv-456",
        "message_id": "msg-789",
        "input_text": "怎么重置密码？",
        "output_text": "您可以按以下步骤重置密码：1.登录账号...",
        "status": "success",
        "input_tokens": 50,
        "output_tokens": 200,
        "total_tokens": 250,
        "duration_ms": 1500,
        "model_provider": "openai",
        "model_name": "gpt-4",
        "cost": 0.005,
        "start_time": "2025-01-03T10:00:00Z",
        "end_time": "2025-01-03T10:00:01.5Z"
      }
    ]
  }
}
```

### 2.2 获取执行详情

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/executions/{execution_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "exec-001",
    "tenant_id": "tenant-001",
    "bot_id": "bot-123",
    "bot_name": "客服助手",
    "user_id": 101,
    "user_name": "张三",
    "conversation_id": "conv-456",
    "input_text": "怎么重置密码？",
    "output_text": "您可以按以下步骤重置密码：1.登录账号...",
    "status": "success",
    "total_tokens": 250,
    "duration_ms": 1500,
    "cost": 0.005,
    "model_provider": "openai",
    "model_name": "gpt-4",
    "start_time": "2025-01-03T10:00:00Z",
    "end_time": "2025-01-03T10:00:01.5Z",
    "trace_tree": {
      "id": "step-1",
      "name": "意图识别",
      "type": "llm_call",
      "input": "用户输入: 怎么重置密码？",
      "output": "意图: account_management",
      "duration_ms": 200,
      "tokens": 50,
      "model_name": "gpt-3.5-turbo",
      "status": "success",
      "children": [
        {
          "id": "step-2",
          "name": "知识库检索",
          "type": "tool_call",
          "input": "关键词: 重置 密码",
          "output": "检索到3篇文档",
          "duration_ms": 300,
          "status": "success",
          "children": []
        },
        {
          "id": "step-3",
          "name": "生成回复",
          "type": "llm_call",
          "input": "基于上下文生成回复",
          "output": "您可以按以下步骤...",
          "duration_ms": 1000,
          "tokens": 200,
          "model_name": "gpt-4",
          "status": "success",
          "children": []
        }
      ]
    }
  }
}
```

### 2.3 获取性能指标

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/metrics`

**查询参数**:
- bot_id: Bot ID (可选)
- start_date: 开始日期 (可选)
- end_date: 结束日期 (可选)
- granularity: 时间粒度 (可选，hour/day/week，默认day)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot-123",
    "bot_name": "客服助手",
    "period": {
      "start": "2025-01-01T00:00:00Z",
      "end": "2025-01-03T23:59:59Z",
      "granularity": "day"
    },
    "metrics": [
      {
        "timestamp": "2025-01-01T00:00:00Z",
        "total_requests": 500,
        "success_requests": 490,
        "failed_requests": 10,
        "success_rate": 0.98,
        "avg_duration_ms": 1200,
        "p50_duration_ms": 1100,
        "p95_duration_ms": 2000,
        "p99_duration_ms": 3500,
        "total_tokens": 125000,
        "avg_tokens_per_request": 250,
        "total_cost": 25.00
      },
      {
        "timestamp": "2025-01-02T00:00:00Z",
        "total_requests": 520,
        "success_requests": 515,
        "failed_requests": 5,
        "success_rate": 0.99,
        "avg_duration_ms": 1150,
        "p50_duration_ms": 1050,
        "p95_duration_ms": 1900,
        "p99_duration_ms": 3200,
        "total_tokens": 130000,
        "avg_tokens_per_request": 250,
        "total_cost": 26.00
      }
    ],
    "summary": {
      "total_requests": 1520,
      "total_tokens": 380000,
      "total_cost": 76.50,
      "avg_success_rate": 0.985,
      "avg_duration_ms": 1175
    }
  }
}
```

### 2.4 获取错误统计

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/errors`

**查询参数**:
- bot_id: Bot ID (可选)
- start_date: 开始日期 (可选)
- end_date: 结束日期 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-01T00:00:00Z",
      "end": "2025-01-03T23:59:59Z"
    },
    "total_errors": 120,
    "error_rate": 0.0789,
    "by_type": {
      "timeout": 45,
      "api_failure": 35,
      "llm_error": 25,
      "validation_error": 10,
      "other": 5
    },
    "top_errors": [
      {
        "error_type": "timeout",
        "error_message": "Request timeout after 30s",
        "count": 45,
        "percentage": 37.5
      },
      {
        "error_type": "api_failure",
        "error_message": "Knowledge base API error",
        "count": 35,
        "percentage": 29.2
      }
    ],
    "errors_by_bot": [
      {
        "bot_id": "bot-123",
        "bot_name": "客服助手",
        "total_errors": 50,
        "error_rate": 0.10
      }
    ]
  }
}
```

---

## 3. 评估API

### 3.1 创建评估任务

**接口地址**: `POST /api/v1/tenants/{tenant_id}/agent-monitor/evaluations`

**请求参数**:
```json
{
  "name": "客服助手准确性评估",
  "bot_id": "bot-123",
  "evaluation_type": "accuracy",
  "dataset_id": "ds-001",
  "evaluation_methods": ["exact_match", "similarity", "llm_as_judge"],
  "ground_truth": [
    {
      "input": "怎么重置密码？",
      "expected_output": "您可以按以下步骤重置密码：1.登录账号..."
    },
    {
      "input": "如何联系客服？",
      "expected_output": "您可以通过以下方式联系客服..."
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "评估任务创建成功",
  "data": {
    "evaluation_id": "eval-001",
    "status": "running",
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

### 3.2 获取评估结果

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/evaluations/{evaluation_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "evaluation_id": "eval-001",
    "name": "客服助手准确性评估",
    "bot_id": "bot-123",
    "bot_name": "客服助手",
    "evaluation_type": "accuracy",
    "status": "completed",
    "created_at": "2025-01-03T10:00:00Z",
    "completed_at": "2025-01-03T10:05:00Z",
    "overall_score": 85.5,
    "passed": true,
    "threshold": 80.0,
    "scores": {
      "exact_match": 70.0,
      "similarity": 88.5,
      "llm_as_judge": 88.0
    },
    "detailed_results": [
      {
        "case_id": 1,
        "input": "怎么重置密码？",
        "actual_output": "您可以按以下步骤重置密码...",
        "expected_output": "您可以按以下步骤重置密码...",
        "score": 90.0,
        "passed": true,
        "method_scores": {
          "exact_match": 85.0,
          "similarity": 92.0,
          "llm_as_judge": 93.0
        }
      },
      {
        "case_id": 2,
        "input": "如何联系客服？",
        "actual_output": "您可以通过以下方式联系客服...",
        "expected_output": "您可以通过以下方式联系客服...",
        "score": 81.0,
        "passed": true,
        "method_scores": {
          "exact_match": 55.0,
          "similarity": 85.0,
          "llm_as_judge": 83.0
        }
      }
    ],
    "issues": [
      {
        "type": "low_accuracy",
        "description": "在处理复杂问题时准确率偏低",
        "affected_cases": 15,
        "percentage": 15.0,
        "suggestion": "建议优化知识库检索策略，增加上下文窗口"
      }
    ]
  }
}
```

### 3.3 提交用户反馈

**接口地址**: `POST /api/v1/tenants/{tenant_id}/agent-monitor/feedback`

**请求参数**:
```json
{
  "execution_id": "exec-001",
  "rating": 5,
  "is_helpful": true,
  "feedback_reason": "回答准确",
  "feedback_text": "很好地解决了我的问题",
  "tags": ["准确", "快速"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "反馈提交成功",
  "data": {
    "feedback_id": 12345,
    "created_at": "2025-01-03T10:30:00Z"
  }
}
```

### 3.4 获取反馈统计

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/feedback/statistics`

**查询参数**:
- bot_id: Bot ID (可选)
- start_date: 开始日期 (可选)
- end_date: 结束日期 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-01T00:00:00Z",
      "end": "2025-01-03T23:59:59Z"
    },
    "total_feedbacks": 520,
    "avg_rating": 4.5,
    "helpful_count": 490,
    "helpful_rate": 0.942,
    "rating_distribution": {
      "5": 380,
      "4": 100,
      "3": 30,
      "2": 8,
      "1": 2
    },
    "top_feedback_reasons": [
      {
        "reason": "回答准确",
        "count": 250
      },
      {
        "reason": "响应快速",
        "count": 180
      }
    ]
  }
}
```

---

## 4. 提示词版本管理API

### 4.1 创建提示词版本

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions`

**请求参数**:
```json
{
  "version": "v1.1.0",
  "system_prompt": "你是一个专业的客服助手，请友好、准确地回答用户问题...",
  "user_prompt_template": "用户问题: {user_input}\n上下文: {context}",
  "parameters": {
    "type": "object",
    "properties": {
      "user_input": {
        "type": "string",
        "description": "用户输入"
      },
      "context": {
        "type": "string",
        "description": "上下文信息"
      }
    },
    "required": ["user_input"]
  },
  "parent_version_id": "v1.0.0",
  "change_description": "优化了回复的专业性和友好度"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "版本创建成功",
  "data": {
    "id": "v1.1.0",
    "tenant_id": "tenant-001",
    "bot_id": "bot-123",
    "version": "v1.1.0",
    "is_active": false,
    "created_at": "2025-01-03T11:00:00Z"
  }
}
```

### 4.2 获取提示词版本列表

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions`

**查询参数**:
- is_active: 是否仅显示激活版本 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "items": [
      {
        "id": "v1.0.0",
        "tenant_id": "tenant-001",
        "bot_id": "bot-123",
        "version": "v1.0.0",
        "is_active": true,
        "change_description": "初始版本",
        "avg_rating": 4.2,
        "total_evaluations": 120,
        "success_rate": 0.92,
        "created_at": "2025-01-01T00:00:00Z"
      },
      {
        "id": "v1.1.0",
        "tenant_id": "tenant-001",
        "bot_id": "bot-123",
        "version": "v1.1.0",
        "is_active": false,
        "parent_version_id": "v1.0.0",
        "change_description": "优化了回复的专业性和友好度",
        "avg_rating": 4.5,
        "total_evaluations": 50,
        "success_rate": 0.95,
        "created_at": "2025-01-03T11:00:00Z"
      }
    ]
  }
}
```

### 4.3 激活提示词版本

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions/{version_id}/activate`

**响应示例**:
```json
{
  "code": 0,
  "message": "版本已激活",
  "data": {
    "id": "v1.1.0",
    "version": "v1.1.0",
    "is_active": true,
    "activated_at": "2025-01-03T12:00:00Z"
  }
}
```

### 4.4 版本对比

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/prompt-versions/compare`

**查询参数**:
- version_id_1: 版本1 ID (必填)
- version_id_2: 版本2 ID (必填)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version_1": {
      "id": "v1.0.0",
      "version": "v1.0.0",
      "system_prompt": "你是一个客服助手..."
    },
    "version_2": {
      "id": "v1.1.0",
      "version": "v1.1.0",
      "system_prompt": "你是一个专业的客服助手，请友好、准确地..."
    },
    "diff": {
      "system_prompt": {
        "type": "modified",
        "old_value": "你是一个客服助手",
        "new_value": "你是一个专业的客服助手，请友好、准确地"
      }
    }
  }
}
```

---

## 5. A/B测试API

### 5.1 创建A/B测试

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/ab-tests`

**请求参数**:
```json
{
  "test_name": "提示词优化A/B测试",
  "control_version_id": "v1.0.0",
  "treatment_version_ids": ["v1.1.0", "v1.2.0"],
  "traffic_split": {
    "v1.0.0": 34,
    "v1.1.0": 33,
    "v1.2.0": 33
  },
  "start_time": "2025-01-03T00:00:00Z",
  "end_time": "2025-01-10T00:00:00Z"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "A/B测试创建成功",
  "data": {
    "id": "ab-test-001",
    "status": "draft",
    "created_at": "2025-01-03T13:00:00Z"
  }
}
```

### 5.2 启动A/B测试

**接口地址**: `POST /api/v1/tenants/{tenant_id}/bots/{bot_id}/ab-tests/{ab_test_id}/start`

**响应示例**:
```json
{
  "code": 0,
  "message": "A/B测试已启动",
  "data": {
    "id": "ab-test-001",
    "status": "running",
    "started_at": "2025-01-03T14:00:00Z"
  }
}
```

### 5.3 获取A/B测试结果

**接口地址**: `GET /api/v1/tenants/{tenant_id}/bots/{bot_id}/ab-tests/{ab_test_id}/results`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "ab-test-001",
    "test_name": "提示词优化A/B测试",
    "status": "completed",
    "start_time": "2025-01-03T00:00:00Z",
    "end_time": "2025-01-10T00:00:00Z",
    "total_participants": 15000,
    "results_by_version": [
      {
        "version_id": "v1.0.0",
        "traffic_percentage": 34,
        "participants": 5100,
        "avg_rating": 4.2,
        "success_rate": 0.92,
        "avg_duration_ms": 1200
      },
      {
        "version_id": "v1.1.0",
        "traffic_percentage": 33,
        "participants": 4950,
        "avg_rating": 4.5,
        "success_rate": 0.95,
        "avg_duration_ms": 1150
      },
      {
        "version_id": "v1.2.0",
        "traffic_percentage": 33,
        "participants": 4950,
        "avg_rating": 4.3,
        "success_rate": 0.93,
        "avg_duration_ms": 1180
      }
    ],
    "winner": {
      "version_id": "v1.1.0",
      "improvement_over_control": 0.3,
      "statistical_significance": 0.95
    }
  }
}
```

---

## 6. 成本分析API

### 6.1 获取成本统计

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/costs`

**查询参数**:
- dimension_type: 维度类型 (可选，total/by_bot/by_user/by_department)
- dimension_id: 维度ID (可选)
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-01T00:00:00Z",
      "end": "2025-01-03T23:59:59Z"
    },
    "total_cost": 2500.50,
    "cost_breakdown": {
      "by_bot": [
        {
          "bot_id": "bot-123",
          "bot_name": "客服助手",
          "total_requests": 15000,
          "total_tokens": 3750000,
          "total_cost": 750.00,
          "percentage": 30.0
        },
        {
          "bot_id": "bot-124",
          "bot_name": "销售助手",
          "total_requests": 25000,
          "total_tokens": 6250000,
          "total_cost": 1250.00,
          "percentage": 50.0
        }
      ],
      "by_model": [
        {
          "model_name": "gpt-4",
          "total_tokens": 5000000,
          "total_cost": 2000.00,
          "percentage": 80.0
        },
        {
          "model_name": "gpt-3.5-turbo",
          "total_tokens": 5000000,
          "total_cost": 500.50,
          "percentage": 20.0
        }
      ]
    },
    "trend": [
      {
        "date": "2025-01-01",
        "total_cost": 800.00
      },
      {
        "date": "2025-01-02",
        "total_cost": 850.00
      },
      {
        "date": "2025-01-03",
        "total_cost": 850.50
      }
    ]
  }
}
```

### 6.2 获取成本预测

**接口地址**: `GET /api/v1/tenants/{tenant_id}/agent-monitor/costs/forecast`

**查询参数**:
- bot_id: Bot ID (可选)
- forecast_days: 预测天数 (默认7天)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot-123",
    "forecast_days": 7,
    "current_daily_cost": 250.00,
    "forecast": [
      {
        "date": "2025-01-04",
        "predicted_cost": 255.00,
        "confidence": "high"
      },
      {
        "date": "2025-01-05",
        "predicted_cost": 260.00,
        "confidence": "high"
      }
    ],
    "total_forecasted_cost": 1785.00,
    "optimization_suggestions": [
      {
        "type": "cache_optimization",
        "description": "启用缓存可降低20%成本",
        "potential_savings": 50.00
      },
      {
        "type": "model_selection",
        "description": "对简单查询使用GPT-3.5可降低30%成本",
        "potential_savings": 75.00
      }
    ]
  }
}
```

---

## 7. 数据模型

### 7.1 AgentExecution

```typescript
interface AgentExecution {
  id: string;
  tenant_id: string;
  user_id: number;
  bot_id: string;
  bot_name?: string;
  conversation_id: string;
  message_id: string;
  input_text: string;
  output_text?: string;
  status: 'pending' | 'running' | 'success' | 'failed' | 'timeout';
  input_tokens: number;
  output_tokens: number;
  total_tokens: number;
  duration_ms?: number;
  model_provider?: string;
  model_name?: string;
  error_type?: string;
  error_message?: string;
  error_stack?: string;
  cost?: number;
  start_time: Date;
  end_time?: Date;
  metadata?: Record<string, any>;
}

interface TraceTreeNode {
  id: string;
  name: string;
  type: 'llm_call' | 'tool_call' | 'routing' | 'validation';
  input?: any;
  output?: any;
  duration_ms: number;
  tokens?: number;
  model_name?: string;
  status: 'success' | 'failed';
  error_message?: string;
  children: TraceTreeNode[];
}
```

### 7.2 Evaluation

```typescript
interface Evaluation {
  evaluation_id: string;
  name: string;
  bot_id: string;
  bot_name?: string;
  evaluation_type: 'accuracy' | 'relevance' | 'safety' | 'hallucination';
  status: 'running' | 'completed' | 'failed';
  overall_score?: number;
  passed?: boolean;
  threshold?: number;
  scores?: {
    [method: string]: number;
  };
  detailed_results?: EvaluationCaseResult[];
  issues?: EvaluationIssue[];
  created_at: Date;
  completed_at?: Date;
}

interface EvaluationCaseResult {
  case_id: number;
  input: string;
  actual_output: string;
  expected_output: string;
  score: number;
  passed: boolean;
  method_scores?: {
    [method: string]: number;
  };
}

interface EvaluationIssue {
  type: string;
  description: string;
  affected_cases: number;
  percentage?: number;
  suggestion: string;
}
```

### 7.3 PromptVersion

```typescript
interface PromptVersion {
  id: string;
  tenant_id: string;
  bot_id: string;
  version: string;
  system_prompt: string;
  user_prompt_template?: string;
  parameters?: JSONSchema;
  parent_version_id?: string;
  is_active: boolean;
  change_description?: string;
  avg_rating?: number;
  total_evaluations?: number;
  success_rate?: number;
  created_by: number;
  created_at: Date;
}
```

### 7.4 ABTest

```typescript
interface ABTest {
  id: string;
  tenant_id: string;
  bot_id: string;
  test_name: string;
  control_version_id: string;
  treatment_version_ids: string[];
  traffic_split: {
    [version_id: string]: number;
  };
  status: 'draft' | 'running' | 'paused' | 'completed';
  start_time?: Date;
  end_time?: Date;
  statistical_significance?: number;
  winner_version_id?: string;
  created_by: number;
  created_at: Date;
}
```

### 7.5 CostAnalytics

```typescript
interface CostAnalytics {
  total_cost: number;
  cost_breakdown: {
    by_bot: CostByDimension[];
    by_model: CostByDimension[];
    by_user: CostByDimension[];
  };
  trend: CostTrendPoint[];
}

interface CostByDimension {
  dimension_id: string;
  dimension_name?: string;
  total_requests: number;
  total_tokens: number;
  total_cost: number;
  percentage: number;
}

interface CostTrendPoint {
  date: string;
  total_cost: number;
}
```

---

## 8. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 38001 | 404 | 执行记录不存在 |
| 38002 | 400 | 无效的时间范围 |
| 38003 | 403 | 无权限访问该租户数据 |
| 38004 | 500 | 追踪树生成失败 |
| 38101 | 404 | 评估任务不存在 |
| 38102 | 400 | 无效的评估类型 |
| 38103 | 400 | 无效的评估方法 |
| 38104 | 400 | 测试数据集不存在 |
| 38105 | 400 | 基准答案格式错误 |
| 38201 | 400 | 提示词版本已存在 |
| 38202 | 404 | 提示词版本不存在 |
| 38203 | 400 | 无法激活当前激活的版本 |
| 38204 | 403 | 无权限修改提示词版本 |
| 38301 | 400 | A/B测试已存在 |
| 38302 | 404 | A/B测试不存在 |
| 38303 | 400 | 流量分配无效（总和必须为100） |
| 38304 | 400 | 测试状态不允许启动 |
| 38401 | 400 | 无效的维度类型 |
| 38402 | 400 | 预测天数无效（1-30天） |
| 38403 | 500 | 成本预测失败 |

---

**文档结束**
