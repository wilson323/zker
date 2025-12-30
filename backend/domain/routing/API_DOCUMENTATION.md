# 智能路由引擎增强 - API文档

## 概述

本文档描述了智能路由引擎增强功能的API接口设计。这些API尚未实现Handler层，这里提供的是API设计规范。

## 基础路径

```
/api/v1/routing
```

## 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "SUCCESS",
  "message_zh": "操作成功",
  "data": {},
  "request_id": "req-123",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

### 错误响应
```json
{
  "code": 209050,
  "message": "Routing optimization failed",
  "message_zh": "路由优化失败",
  "request_id": "req-123",
  "timestamp": "2025-01-01T00:00:00Z",
  "details": {}
}
```

## API接口

### 1. 意图管理

#### 1.1 创建意图
```
POST /api/v1/routing/intents
```

**请求体**:
```json
{
  "intent_name": "query_weather",
  "description": "查询天气",
  "examples": [
    "北京天气怎么样",
    "上海今天下雨吗",
    "明天深圳气温多少"
  ],
  "agent_id": "agent-123",
  "confidence": 0.8
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "intent_id": "intent-123",
    "intent_name": "query_weather",
    "status": "created"
  }
}
```

#### 1.2 训练意图模型
```
POST /api/v1/routing/intents/{intent_id}/train
```

**请求体**:
```json
{
  "samples": [
    "查询北京天气",
    "上海今天冷吗",
    "深圳明天会下雨吗"
  ]
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "intent_id": "intent-123",
    "sample_count": 3,
    "training_status": "completed"
  }
}
```

#### 1.3 识别意图
```
POST /api/v1/routing/intents/recognize
```

**请求体**:
```json
{
  "text": "查询北京天气",
  "tenant_id": "tenant-123"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "intent_id": "intent-123",
    "intent_name": "query_weather",
    "confidence": 0.95,
    "agent_id": "agent-123",
    "workflow_id": "",
    "entities": [
      {
        "entity_name": "location",
        "entity_value": "北京",
        "confidence": 0.98
      },
      {
        "entity_name": "action",
        "entity_value": "查询",
        "confidence": 0.95
      }
    ],
    "match_method": "hybrid"
  }
}
```

#### 1.4 批量识别意图
```
POST /api/v1/routing/intents/batch-recognize
```

**请求体**:
```json
{
  "texts": [
    "查询北京天气",
    "预订上海酒店",
    "明天深圳会下雨吗"
  ],
  "tenant_id": "tenant-123"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "results": [
      {
        "text": "查询北京天气",
        "intent_id": "intent-123",
        "confidence": 0.95
      },
      {
        "text": "预订上海酒店",
        "intent_id": "intent-456",
        "confidence": 0.92
      }
    ]
  }
}
```

### 2. 路由优化

#### 2.1 优化路由
```
POST /api/v1/routing/optimize
```

**请求体**:
```json
{
  "tenant_id": "tenant-123",
  "user_id": "user-123",
  "query": "查询北京天气",
  "context": {
    "device": "mobile",
    "region": "beijing"
  },
  "preferences": {
    "preferred_agent": "agent-123"
  }
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "agent_id": "agent-123",
    "confidence": 0.95,
    "score": 0.87,
    "reasons": [
      "Load: 0.70",
      "Health: 0.95",
      "Response: 0.75",
      "Preference: 0.90"
    ],
    "match_type": "optimized",
    "strategy": "multi_dimension_score",
    "predicted_load": 0.30
  }
}
```

#### 2.2 预测Agent负载
```
GET /api/v1/routing/agents/{agent_id}/load-prediction
```

**查询参数**:
- `start_time`: 开始时间（Unix时间戳）
- `end_time`: 结束时间（Unix时间戳）

**响应**:
```json
{
  "code": 0,
  "data": {
    "agent_id": "agent-123",
    "predicted_load": 0.65,
    "confidence": 0.85,
    "time_window": {
      "start_time": 1735660800,
      "end_time": 1735747200
    }
  }
}
```

### 3. A/B测试

#### 3.1 创建A/B测试
```
POST /api/v1/routing/ab-tests
```

**请求体**:
```json
{
  "tenant_id": "tenant-123",
  "test_name": "路由策略对比测试",
  "description": "对比负载均衡和评分策略",
  "strategy_a": "load_balance",
  "strategy_b": "score_based",
  "traffic_split": 50,
  "sample_size": 1000
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "test_id": "test-123",
    "status": "running",
    "start_time": "2025-01-01T00:00:00Z"
  }
}
```

#### 3.2 执行A/B测试路由
```
POST /api/v1/routing/ab-tests/{test_id}/execute
```

**请求体**:
```json
{
  "tenant_id": "tenant-123",
  "user_id": "user-123",
  "query": "查询天气"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "test_id": "test-123",
    "strategy": "load_balance",
    "agent_id": "agent-123"
  }
}
```

#### 3.3 获取A/B测试结果
```
GET /api/v1/routing/ab-tests/{test_id}/results
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "test_id": "test-123",
    "strategy_a_stats": {
      "strategy": "load_balance",
      "sample_count": 500,
      "avg_response_time": 1520.5,
      "avg_user_rating": 4.2,
      "success_rate": 0.92,
      "error_rate": 0.08
    },
    "strategy_b_stats": {
      "strategy": "score_based",
      "sample_count": 500,
      "avg_response_time": 1480.3,
      "avg_user_rating": 4.5,
      "success_rate": 0.95,
      "error_rate": 0.05
    },
    "winner": "score_based",
    "confidence": 0.98,
    "is_statistically_significant": true,
    "recommendation": "推荐策略B: 成功率较高 (95.00% vs 92.00%)"
  }
}
```

#### 3.4 结束A/B测试
```
POST /api/v1/routing/ab-tests/{test_id}/conclude
```

**请求体**:
```json
{
  "winner": "score_based"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "test_id": "test-123",
    "winner": "score_based",
    "confidence": 0.98,
    "status": "completed",
    "end_time": "2025-01-01T12:00:00Z"
  }
}
```

#### 3.5 暂停/恢复A/B测试
```
POST /api/v1/routing/ab-tests/{test_id}/pause
POST /api/v1/routing/ab-tests/{test_id}/resume
```

### 4. 路由学习

#### 4.1 学习路由日志
```
GET /api/v1/routing/learning/insights
```

**查询参数**:
- `tenant_id`: 租户ID
- `start_time`: 开始时间（Unix时间戳）
- `end_time`: 结束时间（Unix时间戳）

**响应**:
```json
{
  "code": 0,
  "data": {
    "tenant_id": "tenant-123",
    "time_range": {
      "start_time": 1735056000,
      "end_time": 1735660800
    },
    "total_requests": 10000,
    "avg_response_time": 1520.5,
    "success_rate": 0.92,
    "intent_distribution": {
      "query_weather": 3500,
      "book_hotel": 2500,
      "query_route": 2000,
      "other": 2000
    },
    "agent_performance": {
      "agent-123": {
        "agent_id": "agent-123",
        "request_count": 5000,
        "avg_response_time": 1450.0,
        "success_rate": 0.95,
        "avg_user_rating": 4.5,
        "error_rate": 0.05
      }
    },
    "recommendations": [
      "Agent agent-456 错误率较高(15.00%)，建议检查配置或减少负载",
      "Agent agent-789 响应时间较长(3500.00ms)，建议优化"
    ]
  }
}
```

#### 4.2 建议路由规则
```
GET /api/v1/routing/learning/suggestions
```

**查询参数**:
- `tenant_id`: 租户ID

**响应**:
```json
{
  "code": 0,
  "data": {
    "suggestions": [
      {
        "intent": "query_weather",
        "suggested_agent": "agent-123",
        "confidence": 0.92,
        "reason": "该Agent在意图'query_weather'上的表现最佳（150次调用）",
        "expected_improvement": 0.15
      },
      {
        "intent": "book_hotel",
        "suggested_agent": "agent-456",
        "confidence": 0.88,
        "reason": "该Agent在意图'book_hotel'上的表现最佳（120次调用）",
        "expected_improvement": 0.12
      }
    ]
  }
}
```

#### 4.3 优化路由策略
```
POST /api/v1/routing/learning/optimize-strategy
```

**请求体**:
```json
{
  "strategy_id": "strategy-123"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "strategy_id": "strategy-123",
    "strategy_name": "optimized_strategy",
    "config": {
      "load_weight": 0.25,
      "health_weight": 0.20,
      "response_weight": 0.20,
      "error_weight": 0.15,
      "cost_weight": 0.10,
      "preference_weight": 0.10
    },
    "expected_metrics": {
      "avg_response_time": 1500.0,
      "success_rate": 0.92,
      "user_satisfaction": 4.3
    }
  }
}
```

### 5. 实体提取

#### 5.1 提取实体
```
POST /api/v1/routing/entities/extract
```

**请求体**:
```json
{
  "text": "明天从北京到上海的机票"
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "entities": [
      {
        "entity_name": "time",
        "entity_value": "明天",
        "confidence": 0.95,
        "start_position": 0,
        "end_position": 2
      },
      {
        "entity_name": "location",
        "entity_value": "北京",
        "confidence": 0.98,
        "start_position": 4,
        "end_position": 6
      },
      {
        "entity_name": "location",
        "entity_value": "上海",
        "confidence": 0.98,
        "start_position": 8,
        "end_position": 10
      },
      {
        "entity_name": "action",
        "entity_value": "机票",
        "confidence": 0.90,
        "start_position": 11,
        "end_position": 13
      }
    ]
  }
}
```

#### 5.2 批量提取实体
```
POST /api/v1/routing/entities/batch-extract
```

**请求体**:
```json
{
  "texts": [
    "明天从北京到上海的机票",
    "预订深圳的酒店"
  ]
}
```

**响应**:
```json
{
  "code": 0,
  "data": {
    "results": [
      {
        "text": "明天从北京到上海的机票",
        "entities": [...]
      },
      {
        "text": "预订深圳的酒店",
        "entities": [...]
      }
    ]
  }
}
```

## 错误码说明

| 错误码 | HTTP状态 | 说明 |
|--------|----------|------|
| 209001 | 404 | 路由规则不存在 |
| 209002 | 404 | 意图不存在 |
| 209003 | 404 | A/B测试不存在 |
| 209010 | 409 | 路由规则已存在 |
| 209011 | 409 | 意图已存在 |
| 209020 | 400 | 路由参数无效 |
| 209021 | 400 | 意图参数无效 |
| 209022 | 400 | A/B测试参数无效 |
| 209023 | 400 | 路由策略无效 |
| 209050 | 500 | 路由优化失败 |
| 209051 | 500 | 实体提取失败 |
| 209052 | 500 | 向量搜索失败 |
| 209053 | 500 | A/B测试执行失败 |
| 209054 | 500 | 路由学习失败 |
| 209055 | 500 | 评分计算失败 |

## 权限要求

所有API都需要以下HTTP头：
```
Authorization: Bearer {access_token}
X-Tenant-ID: {tenant_id}
```

## 限流规则

- 意图识别: 100次/分钟
- 路由优化: 200次/分钟
- 实体提取: 150次/分钟
- A/B测试: 50次/分钟
- 路由学习: 20次/分钟

## 注意事项

1. **租户隔离**: 所有操作都需要提供 `tenant_id`，确保数据隔离
2. **异步操作**: 意图训练等耗时操作会异步执行，返回任务ID
3. **缓存**: 意图识别结果会缓存5分钟，提高响应速度
4. **监控**: 所有路由决策都会记录日志，用于后续学习优化
