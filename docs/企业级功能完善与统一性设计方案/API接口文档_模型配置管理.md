# API接口文档:模型配置管理系统

**文档编号**: DE-API-2025-DOC-032
**版本**: v1.0.0
**创建日期**: 2025-12-30
**参考设计**: 《32-模型配置管理系统.md》

---

## 📋 目录

1. [提供商管理API](#1-提供商管理api)
2. [模型配置管理API](#2-模型配置管理api)
3. [API密钥管理API](#3-api密钥管理api)
4. [路由规则管理API](#4-路由规则管理api)
5. [模型测试API](#5-模型测试api)
6. [模型统计监控API](#6-模型统计监控api)
7. [审计日志API](#7-审计日志api)

---

## 1. 提供商管理API

### 1.1 获取提供商列表

**接口描述**: 获取所有AI模型提供商列表

**请求方式**: `GET /api/v1/model-config/providers`

**权限要求**: `model_config:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | String | 否 | 状态筛选(active/deprecated/beta) |
| type | String | 否 | 类型筛选(openai/anthropic/azure/google) |

**请求示例**:

```http
GET /api/v1/model-config/providers?status=active
Authorization: Bearer {access_token}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "openai",
      "name": "openai",
      "display_name": "OpenAI",
      "type": "openai",
      "base_url": "https://api.openai.com/v1",
      "documentation_url": "https://platform.openai.com/docs",
      "status": "active",
      "supported_features": [
        "chat",
        "completion",
        "embedding",
        "image",
        "speech",
        "function_calling",
        "streaming",
        "vision"
      ],
      "config_schema": {
        "api_key_field": "Authorization",
        "auth_type": "bearer",
        "headers": {
          "Content-Type": "application/json"
        }
      },
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "anthropic",
      "name": "anthropic",
      "display_name": "Anthropic",
      "type": "anthropic",
      "base_url": "https://api.anthropic.com/v1",
      "documentation_url": "https://docs.anthropic.com",
      "status": "active",
      "supported_features": [
        "chat",
        "completion",
        "streaming",
        "vision"
      ],
      "config_schema": {
        "api_key_field": "x-api-key",
        "auth_type": "custom-header",
        "version": "2023-06-01"
      },
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### 1.2 获取提供商详情

**接口描述**: 获取指定提供商的详细信息

**请求方式**: `GET /api/v1/model-config/providers/{provider_id}`

**权限要求**: `model_config:read`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| provider_id | String | 提供商ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "openai",
    "name": "openai",
    "display_name": "OpenAI",
    "type": "openai",
    "base_url": "https://api.openai.com/v1",
    "documentation_url": "https://platform.openai.com/docs",
    "status": "active",
    "supported_features": ["chat", "completion", "embedding", "image", "speech"],
    "available_models": [
      {
        "model_id": "gpt-4-turbo",
        "model_name": "gpt-4-turbo-preview",
        "display_name": "GPT-4 Turbo",
        "status": "active"
      },
      {
        "model_id": "gpt-3.5-turbo",
        "model_name": "gpt-3.5-turbo",
        "display_name": "GPT-3.5 Turbo",
        "status": "active"
      }
    ],
    "rate_limits": {
      "requests_per_minute": 3500,
      "tokens_per_minute": 90000,
      "requests_per_day": 10000
    },
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

## 2. 模型配置管理API

### 2.1 获取模型配置列表

**接口描述**: 获取所有模型配置

**请求方式**: `GET /api/v1/model-config/models`

**权限要求**: `model_config:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| provider_id | String | 否 | 提供商筛选 |
| model_type | String | 否 | 类型筛选(chat/completion/embedding/image/speech) |
| status | String | 否 | 状态筛选(active/inactive/deprecated) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "gpt-4-turbo",
      "provider_id": "openai",
      "model_name": "gpt-4-turbo-preview",
      "display_name": "GPT-4 Turbo",
      "model_type": "chat",
      "capabilities": {
        "max_tokens": 4096,
        "context_length": 128000,
        "streaming": true,
        "function_calling": true,
        "vision": true
      },
      "pricing": {
        "input_price_per_1k_tokens": 0.01,
        "output_price_per_1k_tokens": 0.03,
        "currency": "USD"
      },
      "parameters": {
        "temperature": {
          "min": 0,
          "max": 2,
          "default": 0.7
        },
        "top_p": {
          "min": 0,
          "max": 1,
          "default": 1
        },
        "max_tokens": {
          "max": 4096,
          "default": 2048
        }
      },
      "status": "active",
      "is_default": false,
      "version": "2024-04-09",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### 2.2 创建模型配置

**接口描述**: 创建新的模型配置

**请求方式**: `POST /api/v1/model-config/models`

**权限要求**: `model_config:write`

**请求体**:

```json
{
  "data": {
    "provider_id": "openai",
    "model_name": "gpt-4-turbo-preview",
    "display_name": "GPT-4 Turbo",
    "model_type": "chat",
    "capabilities": {
      "max_tokens": 4096,
      "context_length": 128000,
      "streaming": true,
      "function_calling": true,
      "vision": true
    },
    "pricing": {
      "input_price_per_1k_tokens": 0.01,
      "output_price_per_1k_tokens": 0.03,
      "currency": "USD"
    },
    "parameters": {
      "temperature": {
        "min": 0,
        "max": 2,
        "default": 0.7
      },
      "max_tokens": {
        "max": 4096,
        "default": 2048
      }
    },
    "status": "active",
    "is_default": false
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Model config created",
  "data": {
    "id": "gpt-4-turbo",
    "provider_id": "openai",
    "model_name": "gpt-4-turbo-preview",
    "display_name": "GPT-4 Turbo",
    "model_type": "chat",
    "status": "active",
    "created_at": "2025-12-30T10:00:00Z"
  }
}
```

### 2.3 更新模型配置

**接口描述**: 更新模型配置

**请求方式**: `PATCH /api/v1/model-config/models/{model_id}`

**权限要求**: `model_config:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| model_id | String | 模型配置ID |

**请求体**:

```json
{
  "data": {
    "display_name": "GPT-4 Turbo (Updated)",
    "pricing": {
      "input_price_per_1k_tokens": 0.008,
      "output_price_per_1k_tokens": 0.025,
      "currency": "USD"
    },
    "status": "active"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Model config updated",
  "data": {
    "id": "gpt-4-turbo",
    "display_name": "GPT-4 Turbo (Updated)",
    "updated_at": "2025-12-30T11:00:00Z"
  }
}
```

### 2.4 删除模型配置

**接口描述**: 删除模型配置

**请求方式**: `DELETE /api/v1/model-config/models/{model_id}`

**权限要求**: `model_config:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| model_id | String | 模型配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Model config deleted",
  "data": {
    "id": "gpt-4-turbo",
    "deleted_at": "2025-12-30T12:00:00Z"
  }
}
```

### 2.5 设置默认模型

**接口描述**: 设置指定类型的默认模型

**请求方式**: `POST /api/v1/model-config/models/{model_id}/set-default`

**权限要求**: `model_config:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| model_id | String | 模型配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Default model updated",
  "data": {
    "model_type": "chat",
    "default_model_id": "gpt-4-turbo",
    "updated_at": "2025-12-30T12:00:00Z"
  }
}
```

---

## 3. API密钥管理API

### 3.1 获取API密钥列表

**接口描述**: 获取API密钥列表(脱敏显示)

**请求方式**: `GET /api/v1/model-config/api-keys`

**权限要求**: `api_key:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| provider_id | String | 否 | 提供商筛选 |
| is_active | Boolean | 否 | 是否激活 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "key_123456",
      "provider_id": "openai",
      "name": "生产环境密钥",
      "key_preview": "sk-xxxx...xxxx",
      "is_active": true,
      "quota_limit": 1000000,
      "quota_used": 15230,
      "quota_reset_at": "2025-12-31T00:00:00Z",
      "priority": 10,
      "last_used_at": "2025-12-30T09:30:00Z",
      "expires_at": "2026-12-30T00:00:00Z",
      "created_by": "admin@example.com",
      "created_at": "2025-01-01T00:00:00Z",
      "revoked_at": null
    }
  ]
}
```

### 3.2 创建API密钥

**接口描述**: 创建新的API密钥

**请求方式**: `POST /api/v1/model-config/api-keys`

**权限要求**: `api_key:write`

**请求体**:

```json
{
  "data": {
    "provider_id": "openai",
    "name": "生产环境密钥",
    "api_key": "sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "quota_limit": 1000000,
    "priority": 10
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| provider_id | String | ✅ | 提供商ID |
| name | String | ✅ | 密钥名称 |
| api_key | String | ✅ | API密钥(完整密钥,仅创建时可见) |
| quota_limit | Long | 否 | 配额限制(请求数/月) |
| priority | Integer | 否 | 优先级(数字越大优先级越高) |

**响应示例**:

```json
{
  "code": 0,
  "message": "API key created",
  "data": {
    "id": "key_123456",
    "provider_id": "openai",
    "name": "生产环境密钥",
    "key_preview": "sk-proj-xxxx...xxxx",
    "is_active": true,
    "quota_limit": 1000000,
    "quota_used": 0,
    "priority": 10,
    "created_at": "2025-12-30T13:00:00Z"
  }
}
```

### 3.3 撤销API密钥

**接口描述**: 撤销指定的API密钥

**请求方式**: `POST /api/v1/model-config/api-keys/{key_id}/revoke`

**权限要求**: `api_key:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| key_id | String | 密钥ID |

**请求体**:

```json
{
  "data": {
    "reason": "密钥已过期"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "API key revoked",
  "data": {
    "id": "key_123456",
    "is_active": false,
    "revoked_at": "2025-12-30T14:00:00Z"
  }
}
```

### 3.4 轮换API密钥

**接口描述**: 轮换API密钥(无缝切换)

**请求方式**: `POST /api/v1/model-config/api-keys/{key_id}/rotate`

**权限要求**: `api_key:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| key_id | String | 密钥ID |

**请求体**:

```json
{
  "data": {
    "new_api_key": "sk-proj-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
    "grace_period_seconds": 300
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "API key rotated",
  "data": {
    "old_key_id": "key_123456",
    "new_key_id": "key_789012",
    "rotated_at": "2025-12-30T15:00:00Z",
    "grace_period_ends_at": "2025-12-30T15:05:00Z"
  }
}
```

---

## 4. 路由规则管理API

### 4.1 获取路由规则列表

**接口描述**: 获取所有路由规则

**请求方式**: `GET /api/v1/model-config/routing-rules`

**权限要求**: `routing_rule:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| is_active | Boolean | 否 | 是否激活 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "cost-optimization",
      "name": "成本优化路由",
      "description": "简单查询使用低成本模型",
      "priority": 100,
      "conditions": {
        "complexity": "simple",
        "estimated_tokens": {
          "min": 0,
          "max": 1000
        }
      },
      "model_selection": {
        "primary_model": "gpt-3.5-turbo",
        "fallback_models": [
          "claude-3-haiku",
          "gpt-4-turbo-mini"
        ],
        "selection_strategy": "cost"
      },
      "rate_limits": {
        "max_requests_per_minute": 100,
        "max_tokens_per_minute": 50000
      },
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "quality-first",
      "name": "质量优先路由",
      "description": "复杂推理使用高质量模型",
      "priority": 50,
      "conditions": {
        "complexity": "complex",
        "user_tier": ["pro", "enterprise"]
      },
      "model_selection": {
        "primary_model": "gpt-4-turbo",
        "fallback_models": [
          "claude-3-opus",
          "gpt-4"
        ],
        "selection_strategy": "quality"
      },
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### 4.2 创建路由规则

**接口描述**: 创建新的路由规则

**请求方式**: `POST /api/v1/model-config/routing-rules`

**权限要求**: `routing_rule:write`

**请求体**:

```json
{
  "data": {
    "name": "VIP用户路由",
    "description": "VIP用户使用GPT-4",
    "priority": 75,
    "conditions": {
      "user_tier": ["enterprise"],
      "conversation_type": ["chat"]
    },
    "model_selection": {
      "primary_model": "gpt-4-turbo",
      "fallback_models": ["claude-3-opus"],
      "selection_strategy": "quality"
    },
    "rate_limits": {
      "max_requests_per_minute": 200,
      "max_tokens_per_minute": 100000
    },
    "is_active": true
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Routing rule created",
  "data": {
    "id": "vip-user-routing",
    "name": "VIP用户路由",
    "priority": 75,
    "is_active": true,
    "created_at": "2025-12-30T16:00:00Z"
  }
}
```

### 4.3 更新路由规则

**接口描述**: 更新路由规则

**请求方式**: `PATCH /api/v1/model-config/routing-rules/{rule_id}`

**权限要求**: `routing_rule:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| rule_id | String | 规则ID |

**请求体**:

```json
{
  "data": {
    "name": "VIP用户路由(已更新)",
    "priority": 70,
    "is_active": true
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Routing rule updated",
  "data": {
    "id": "vip-user-routing",
    "name": "VIP用户路由(已更新)",
    "updated_at": "2025-12-30T17:00:00Z"
  }
}
```

### 4.4 删除路由规则

**接口描述**: 删除路由规则

**请求方式**: `DELETE /api/v1/model-config/routing-rules/{rule_id}`

**权限要求**: `routing_rule:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| rule_id | String | 规则ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Routing rule deleted",
  "data": {
    "id": "vip-user-routing",
    "deleted_at": "2025-12-30T18:00:00Z"
  }
}
```

---

## 5. 模型测试API

### 5.1 测试模型

**接口描述**: 向指定模型发送测试请求

**请求方式**: `POST /api/v1/model-config/models/test`

**权限要求**: `model_config:test`

**请求体**:

```json
{
  "data": {
    "provider_id": "openai",
    "model_name": "gpt-4-turbo-preview",
    "test_prompt": "你好,请介绍一下你自己",
    "parameters": {
      "temperature": 0.7,
      "max_tokens": 500
    }
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| provider_id | String | ✅ | 提供商ID |
| model_name | String | ✅ | 模型名称 |
| test_prompt | String | ✅ | 测试提示词 |
| parameters | Object | 否 | 模型参数(可选) |

**响应示例** (成功):

```json
{
  "code": 0,
  "message": "Test successful",
  "data": {
    "test_id": "test_20251230123456",
    "provider_id": "openai",
    "model_name": "gpt-4-turbo-preview",
    "request": {
      "prompt": "你好,请介绍一下你自己",
      "parameters": {
        "temperature": 0.7,
        "max_tokens": 500
      }
    },
    "response": "你好!我是由OpenAI开发的语言模型GPT-4 Turbo...",
    "usage": {
      "input_tokens": 15,
      "output_tokens": 120,
      "total_tokens": 135
    },
    "latency_ms": 1250,
    "cost": 0.003825,
    "status": "success",
    "tested_at": "2025-12-30T19:00:00Z"
  }
}
```

**响应示例** (失败):

```json
{
  "code": 50003,
  "message": "Model API call failed",
  "data": {
    "test_id": "test_20251230123457",
    "error": {
      "code": "INVALID_API_KEY",
      "message": "提供的API密钥无效",
      "type": "authentication_error",
      "details": {
        "provider": "openai",
        "model": "gpt-4-turbo-preview",
        "suggestion": "请检查API密钥配置"
      }
    },
    "status": "error",
    "tested_at": "2025-12-30T19:00:05Z"
  }
}
```

### 5.2 批量测试模型

**接口描述**: 批量测试多个模型

**请求方式**: `POST /api/v1/model-config/models/batch-test`

**权限要求**: `model_config:test`

**请求体**:

```json
{
  "data": {
    "models": [
      {
        "provider_id": "openai",
        "model_name": "gpt-4-turbo"
      },
      {
        "provider_id": "anthropic",
        "model_name": "claude-3-opus"
      }
    ],
    "test_prompt": "解释什么是人工智能",
    "parallel": true
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Batch test completed",
  "data": {
    "test_batch_id": "batch_test_20251230123456",
    "total_models": 2,
    "results": [
      {
        "provider_id": "openai",
        "model_name": "gpt-4-turbo",
        "status": "success",
        "latency_ms": 1250,
        "cost": 0.0038
      },
      {
        "provider_id": "anthropic",
        "model_name": "claude-3-opus",
        "status": "success",
        "latency_ms": 1580,
        "cost": 0.0062
      }
    ],
    "completed_at": "2025-12-30T19:05:00Z"
  }
}
```

---

## 6. 模型统计监控API

### 6.1 获取模型统计概览

**接口描述**: 获取模型使用统计概览

**请求方式**: `GET /api/v1/model-config/stats/overview`

**权限要求**: `model_stats:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | String | ✅ | 统计周期(today/week/month/custom) |
| start_date | String | 否 | 自定义开始日期(ISO8601) |
| end_date | String | 否 | 自定义结束日期(ISO8601) |
| model_id | String | 否 | 模型筛选 |
| bot_id | String | 否 | Bot筛选 |
| group_by | String | 否 | 分组维度(model/bot/user) |

**请求示例**:

```http
GET /api/v1/model-config/stats/overview?period=week&group_by=model
Authorization: Bearer {access_token}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "week",
    "start_date": "2025-12-23T00:00:00Z",
    "end_date": "2025-12-30T23:59:59Z",
    "summary": {
      "total_calls": 152340,
      "successful_calls": 151850,
      "failed_calls": 490,
      "success_rate": 99.68,
      "total_tokens": 45678900,
      "total_cost": 1234.56,
      "currency": "USD"
    },
    "by_model": [
      {
        "model_id": "gpt-3.5-turbo",
        "model_name": "GPT-3.5 Turbo",
        "calls": 125800,
        "tokens": 35240000,
        "cost": 450.25,
        "avg_latency_ms": 350,
        "success_rate": 99.85
      },
      {
        "model_id": "gpt-4-turbo",
        "model_name": "GPT-4 Turbo",
        "calls": 26540,
        "tokens": 10438900,
        "cost": 784.31,
        "avg_latency_ms": 1200,
        "success_rate": 99.15
      }
    ],
    "daily_breakdown": [
      {
        "date": "2025-12-24",
        "calls": 21500,
        "tokens": 6750000,
        "cost": 175.50
      },
      {
        "date": "2025-12-25",
        "calls": 18200,
        "tokens": 5460000,
        "cost": 150.25
      }
    ]
  }
}
```

### 6.2 获取模型详细统计

**接口描述**: 获取指定模型的详细统计

**请求方式**: `GET /api/v1/model-config/stats/models/{model_id}`

**权限要求**: `model_stats:read`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| model_id | String | 模型配置ID |

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | String | ✅ | 开始日期 |
| end_date | String | ✅ | 结束日期 |
| granularity | String | 否 | 时间粒度(hour/day/week) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model_id": "gpt-4-turbo",
    "model_name": "GPT-4 Turbo",
    "period": {
      "start_date": "2025-12-23T00:00:00Z",
      "end_date": "2025-12-30T23:59:59Z"
    },
    "usage_stats": {
      "total_calls": 26540,
      "total_tokens": 10438900,
      "input_tokens": 7820000,
      "output_tokens": 2618900,
      "avg_tokens_per_call": 393,
      "total_cost": 784.31
    },
    "performance_stats": {
      "avg_latency_ms": 1200,
      "p50_latency_ms": 1100,
      "p95_latency_ms": 1800,
      "p99_latency_ms": 2500,
      "success_rate": 99.15,
      "error_rate": 0.85
    },
    "error_breakdown": [
      {
        "error_code": "rate_limit_exceeded",
        "count": 156,
        "percentage": 0.59
      },
      {
        "error_code": "timeout",
        "count": 68,
        "percentage": 0.26
      }
    ],
    "hourly_breakdown": [
      {
        "hour": "2025-12-30T10:00:00Z",
        "calls": 1250,
        "tokens": 490000,
        "avg_latency_ms": 1150,
        "cost": 36.75
      }
    ]
  }
}
```

### 6.3 获取模型健康状态

**接口描述**: 获取模型的实时健康状态

**请求方式**: `GET /api/v1/model-config/health`

**权限要求**: `model_stats:read`

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "model_id": "gpt-4-turbo",
      "model_name": "GPT-4 Turbo",
      "provider_id": "openai",
      "status": "healthy",
      "last_check_at": "2025-12-30T19:10:00Z",
      "metrics": {
        "last_5min": {
          "total_calls": 150,
          "success_rate": 99.33,
          "avg_latency_ms": 1200,
          "error_rate": 0.67
        },
        "last_1hour": {
          "total_calls": 1800,
          "success_rate": 99.12,
          "avg_latency_ms": 1250,
          "error_rate": 0.88
        }
      },
      "active_api_keys": 2,
      "available_quota": 85
    },
    {
      "model_id": "claude-3-opus",
      "model_name": "Claude 3 Opus",
      "provider_id": "anthropic",
      "status": "healthy",
      "last_check_at": "2025-12-30T19:10:00Z",
      "metrics": {
        "last_5min": {
          "total_calls": 85,
          "success_rate": 98.82,
          "avg_latency_ms": 1580,
          "error_rate": 1.18
        },
        "last_1hour": {
          "total_calls": 980,
          "success_rate": 99.08,
          "avg_latency_ms": 1620,
          "error_rate": 0.92
        }
      },
      "active_api_keys": 1,
      "available_quota": 100
    }
  ]
}
```

### 6.4 获取成本分析报告

**接口描述**: 获取成本分析报告

**请求方式**: `GET /api/v1/model-config/stats/cost-report`

**权限要求**: `model_stats:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | String | ✅ | 开始日期 |
| end_date | String | ✅ | 结束日期 |
| group_by | String | 否 | 分组维度(model/bot/user/team) |
| format | String | 否 | 报告格式(json/csv/excel) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start_date": "2025-12-01T00:00:00Z",
      "end_date": "2025-12-30T23:59:59Z"
    },
    "summary": {
      "total_cost": 5678.90,
      "currency": "USD",
      "total_calls": 452340,
      "total_tokens": 156789000
    },
    "by_model": [
      {
        "model_id": "gpt-4-turbo",
        "model_name": "GPT-4 Turbo",
        "calls": 85230,
        "tokens": 35240000,
        "cost": 2500.50,
        "percentage": 44.04
      },
      {
        "model_id": "gpt-3.5-turbo",
        "model_name": "GPT-3.5 Turbo",
        "calls": 345600,
        "tokens": 115200000,
        "cost": 2880.30,
        "percentage": 50.72
      }
    ],
    "by_bot": [
      {
        "bot_id": "bot_123",
        "bot_name": "客服助手",
        "calls": 125600,
        "cost": 1234.56,
        "percentage": 21.74
      }
    ],
    "by_user_tier": [
      {
        "user_tier": "enterprise",
        "calls": 85230,
        "cost": 3456.78,
        "percentage": 60.88
      },
      {
        "user_tier": "pro",
        "calls": 152340,
        "cost": 1890.12,
        "percentage": 33.29
      }
    ],
    "daily_trend": [
      {
        "date": "2025-12-01",
        "cost": 123.45
      },
      {
        "date": "2025-12-02",
        "cost": 145.67
      }
    ]
  }
}
```

---

## 7. 审计日志API

### 7.1 查询审计日志

**接口描述**: 查询模型配置相关的审计日志

**请求方式**: `GET /api/v1/model-config/audit-logs`

**权限要求**: `audit_log:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认50) |
| user_id | String | 否 | 用户筛选 |
| action | String | 否 | 操作类型筛选 |
| resource_type | String | 否 | 资源类型筛选 |
| resource_id | String | 否 | 资源ID筛选 |
| start_date | String | 否 | 开始日期 |
| end_date | String | 否 | 结束日期 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 123456,
      "user_id": "admin@example.com",
      "username": "Admin User",
      "action": "model_config.created",
      "resource_type": "model_config",
      "resource_id": "gpt-4-turbo",
      "resource_name": "GPT-4 Turbo",
      "old_value": null,
      "new_value": {
        "id": "gpt-4-turbo",
        "provider_id": "openai",
        "model_name": "gpt-4-turbo-preview",
        "display_name": "GPT-4 Turbo"
      },
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2025-12-30T10:00:00Z"
    },
    {
      "id": 123457,
      "user_id": "admin@example.com",
      "username": "Admin User",
      "action": "api_key.created",
      "resource_type": "api_key",
      "resource_id": "key_123456",
      "resource_name": "生产环境密钥",
      "old_value": null,
      "new_value": {
        "id": "key_123456",
        "provider_id": "openai",
        "name": "生产环境密钥",
        "key_preview": "sk-xxxx...xxxx"
      },
      "ip_address": "192.168.1.100",
      "created_at": "2025-12-30T11:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 50,
      "total_count": 1523,
      "total_pages": 31
    }
  }
}
```

### 7.2 导出审计日志

**接口描述**: 导出审计日志为文件

**请求方式**: `GET /api/v1/model-config/audit-logs/export`

**权限要求**: `audit_log:export`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | String | ✅ | 开始日期 |
| end_date | String | ✅ | 结束日期 |
| format | String | 否 | 文件格式(csv/excel/json) |
| fields | String | 否 | 包含字段(逗号分隔) |

**响应示例**:

```http
HTTP/1.1 200 OK
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="audit_logs_20251201_20251230.xlsx"

[二进制文件内容]
```

---

## 附录

### A. 错误码参考

| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 40001 | 请求参数错误 | 400 |
| 40101 | 未认证 | 401 |
| 40301 | 无权限访问 | 403 |
| 40302 | 无权限修改默认模型 | 403 |
| 40401 | 提供商不存在 | 404 |
| 40402 | 模型配置不存在 | 404 |
| 40403 | API密钥不存在 | 404 |
| 40404 | 路由规则不存在 | 404 |
| 40901 | 模型配置已存在 | 409 |
| 42201 | 模型配置验证失败 | 422 |
| 42202 | API密钥无效 | 422 |
| 42203 | 路由规则冲突 | 422 |
| 42204 | 不能删除默认模型 | 422 |
| 50001 | 模型API调用失败 | 500 |
| 50002 | Vault存储失败 | 500 |
| 50003 | 无可用模型 | 500 |
| 50004 | 密钥解密失败 | 500 |
| 50005 | 配额已用尽 | 429 |

### B. Webhook事件类型

**模型配置事件**:

```json
{
  "event": "model_config.created",
  "timestamp": "2025-12-30T10:00:00Z",
  "data": {
    "model_id": "gpt-4-turbo",
    "model_name": "GPT-4 Turbo",
    "provider_id": "openai",
    "created_by": "admin@example.com"
  }
}
```

```json
{
  "event": "model_config.updated",
  "timestamp": "2025-12-30T11:00:00Z",
  "data": {
    "model_id": "gpt-4-turbo",
    "changes": {
      "pricing": {
        "old": { "input_price_per_1k_tokens": 0.01 },
        "new": { "input_price_per_1k_tokens": 0.008 }
      }
    },
    "updated_by": "admin@example.com"
  }
}
```

```json
{
  "event": "model.test.completed",
  "timestamp": "2025-12-30T12:00:00Z",
  "data": {
    "test_id": "test_20251230123456",
    "model_id": "gpt-4-turbo",
    "status": "success",
    "latency_ms": 1250,
    "cost": 0.0038,
    "tested_by": "admin@example.com"
  }
}
```

**告警事件**:

```json
{
  "event": "alert.model.error_rate_high",
  "timestamp": "2025-12-30T15:00:00Z",
  "data": {
    "model_id": "gpt-4-turbo",
    "error_rate": 15.5,
    "threshold": 10,
    "last_5min_calls": 150,
    "recommendation": "建议检查API密钥或切换到备用模型"
  }
}
```

```json
{
  "event": "alert.api_key.expiring_soon",
  "timestamp": "2025-12-30T16:00:00Z",
  "data": {
    "key_id": "key_123456",
    "provider_id": "openai",
    "expires_at": "2025-12-31T00:00:00Z",
    "days_until_expiry": 1,
    "recommendation": "请尽快轮换密钥"
  }
}
```

```json
{
  "event": "alert.cost.budget_exceeded",
  "timestamp": "2025-12-30T17:00:00Z",
  "data": {
    "period": "month",
    "budget": 1000,
    "actual_cost": 1234.56,
    "exceeded_by": 234.56,
    "recommendation": "考虑启用成本优化路由规则"
  }
}
```

### C. Go后端实现示例

```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// ModelConfigHandler 模型配置处理器
type ModelConfigHandler struct {
    configService  *ModelConfigService
    keyManager     *APIKeyManager
    routingService *ModelRoutingService
    monitorService *ModelMonitorService
}

// ListProviders 获取提供商列表
func (h *ModelConfigHandler) ListProviders(ctx context.Context, c *app.RequestContext) {
    status := c.Query("status")
    if status == "" {
        status = "active"
    }

    providers, err := h.configService.ListProviders(ctx, &ListProvidersRequest{
        Status: status,
    })
    if err != nil {
        c.JSON(500, gin.H{"code": 50001, "message": "Internal server error"})
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
        "data":    providers,
    })
}

// CreateModelConfig 创建模型配置
func (h *ModelConfigHandler) CreateModelConfig(ctx context.Context, c *app.RequestContext) {
    var req CreateModelConfigRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    config, err := h.configService.CreateModelConfig(ctx, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(201, gin.H{
        "code":    0,
        "message": "Model config created",
        "data":    config,
    })
}

// TestModel 测试模型
func (h *ModelConfigHandler) TestModel(ctx context.Context, c *app.RequestContext) {
    var req TestModelRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    // 获取API密钥
    apiKey, err := h.keyManager.GetAPIKey(ctx, req.ProviderID)
    if err != nil {
        c.JSON(500, gin.H{"code": 50002, "message": "Failed to get API key"})
        return
    }

    // 调用模型API
    startTime := time.Now()
    response, err := h.callModelAPI(ctx, req.ProviderID, req.ModelName, apiKey, req.TestPrompt)
    latency := time.Since(startTime).Milliseconds()

    if err != nil {
        // 记录测试失败
        h.recordTestResult(ctx, &req, nil, err.Error(), latency)
        c.JSON(500, gin.H{
            "code":    50003,
            "message": "Model API call failed",
            "data": gin.H{
                "error": err.Error(),
            },
        })
        return
    }

    // 记录测试成功
    h.recordTestResult(ctx, &req, response, "", latency)

    c.JSON(200, gin.H{
        "code":    0,
        "message": "Test successful",
        "data": gin.H{
            "test_id":    generateTestID(),
            "response":   response,
            "latency_ms": latency,
            "cost":       calculateCost(req.ProviderID, req.ModelName,
                                    estimateTokens(req.TestPrompt),
                                    estimateTokens(response)),
        },
    })
}

// GetModelStats 获取模型统计
func (h *ModelConfigHandler) GetModelStats(ctx context.Context, c *app.RequestContext) {
    modelID := c.Query("model_id")
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    groupBy := c.Query("group_by")

    stats, err := h.monitorService.GetModelStats(ctx, &GetModelStatsRequest{
        ModelID:   modelID,
        StartDate: startDate,
        EndDate:   endDate,
        GroupBy:   groupBy,
    })
    if err != nil {
        c.JSON(500, gin.H{"code": 50001, "message": "Failed to get stats"})
        return
    }

    c.JSON(200, gin.H{
        "code":    0,
        "message": "success",
        "data":    stats,
    })
}
```

### D. 前端TypeScript类型定义

```typescript
// types/model-config.ts

// 提供商
interface Provider {
  id: string;
  name: string;
  display_name: string;
  type: 'openai' | 'anthropic' | 'azure' | 'google' | 'coze' | 'custom';
  base_url: string;
  documentation_url: string;
  status: 'active' | 'deprecated' | 'beta';
  supported_features: string[];
  config_schema: Record<string, any>;
  created_at: string;
  updated_at: string;
}

// 模型配置
interface ModelConfig {
  id: string;
  provider_id: string;
  model_name: string;
  display_name: string;
  model_type: 'chat' | 'completion' | 'embedding' | 'image' | 'speech';
  capabilities: {
    max_tokens: number;
    context_length: number;
    streaming: boolean;
    function_calling: boolean;
    vision: boolean;
  };
  pricing: {
    input_price_per_1k_tokens: number;
    output_price_per_1k_tokens: number;
    currency: string;
  };
  parameters: {
    temperature?: { min: number; max: number; default: number };
    top_p?: { min: number; max: number; default: number };
    max_tokens?: { max: number; default: number };
  };
  status: 'active' | 'inactive' | 'deprecated';
  is_default: boolean;
  version?: string;
  created_at: string;
  updated_at: string;
}

// API密钥
interface APIKey {
  id: string;
  provider_id: string;
  name: string;
  key_preview: string;
  is_active: boolean;
  quota_limit?: number;
  quota_used: number;
  quota_reset_at: string;
  priority: number;
  last_used_at?: string;
  expires_at?: string;
  created_by: string;
  created_at: string;
  revoked_at?: string;
}

// 路由规则
interface RoutingRule {
  id: string;
  name: string;
  description: string;
  priority: number;
  conditions: {
    bot_id?: string[];
    user_tier?: string[];
    conversation_type?: string[];
    complexity?: 'simple' | 'medium' | 'complex';
    estimated_tokens?: { min: number; max: number };
  };
  model_selection: {
    primary_model: string;
    fallback_models: string[];
    selection_strategy: 'cost' | 'quality' | 'speed' | 'custom';
  };
  rate_limits?: {
    max_requests_per_minute?: number;
    max_tokens_per_minute?: number;
  };
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// API客户端
class ModelConfigApiClient {
  private baseURL = '/api/v1/model-config';

  async listProviders(params?: { status?: string }): Promise<Provider[]> {
    const response = await apiClient.get(`${this.baseURL}/providers`, { params });
    return response.data.data;
  }

  async listModels(params?: {
    provider_id?: string;
    model_type?: string;
    status?: string;
  }): Promise<ModelConfig[]> {
    const response = await apiClient.get(`${this.baseURL}/models`, { params });
    return response.data.data;
  }

  async createModelConfig(data: CreateModelConfigRequest): Promise<ModelConfig> {
    const response = await apiClient.post(`${this.baseURL}/models`, { data });
    return response.data.data;
  }

  async updateModelConfig(
    modelId: string,
    data: Partial<ModelConfig>
  ): Promise<ModelConfig> {
    const response = await apiClient.patch(`${this.baseURL}/models/${modelId}`, { data });
    return response.data.data;
  }

  async deleteModelConfig(modelId: string): Promise<void> {
    await apiClient.delete(`${this.baseURL}/models/${modelId}`);
  }

  async testModel(data: TestModelRequest): Promise<TestResult> {
    const response = await apiClient.post(`${this.baseURL}/models/test`, { data });
    return response.data.data;
  }

  async createAPIKey(data: CreateAPIKeyRequest): Promise<APIKey> {
    const response = await apiClient.post(`${this.baseURL}/api-keys`, { data });
    return response.data.data;
  }

  async listAPIKeys(params?: {
    provider_id?: string;
    is_active?: boolean;
  }): Promise<APIKey[]> {
    const response = await apiClient.get(`${this.baseURL}/api-keys`, { params });
    return response.data.data;
  }

  async revokeAPIKey(keyId: string, reason?: string): Promise<void> {
    await apiClient.post(`${this.baseURL}/api-keys/${keyId}/revoke`, {
      data: { reason },
    });
  }

  async rotateAPIKey(
    keyId: string,
    newApiKey: string,
    gracePeriodSeconds?: number
  ): Promise<RotateResult> {
    const response = await apiClient.post(`${this.baseURL}/api-keys/${keyId}/rotate`, {
      data: {
        new_api_key: newApiKey,
        grace_period_seconds: gracePeriodSeconds,
      },
    });
    return response.data.data;
  }

  async listRoutingRules(params?: {
    is_active?: boolean;
  }): Promise<RoutingRule[]> {
    const response = await apiClient.get(`${this.baseURL}/routing-rules`, { params });
    return response.data.data;
  }

  async createRoutingRule(data: CreateRoutingRuleRequest): Promise<RoutingRule> {
    const response = await apiClient.post(`${this.baseURL}/routing-rules`, { data });
    return response.data.data;
  }

  async getModelStats(params: {
    model_id: string;
    start_date: string;
    end_date: string;
    group_by?: string;
  }): Promise<ModelStats> {
    const response = await apiClient.get(`${this.baseURL}/stats/models`, { params });
    return response.data.data;
  }

  async getCostReport(params: {
    start_date: string;
    end_date: string;
    group_by?: string;
  }): Promise<CostReport> {
    const response = await apiClient.get(`${this.baseURL}/stats/cost-report`, { params });
    return response.data.data;
  }

  async getModelHealth(): Promise<ModelHealth[]> {
    const response = await apiClient.get(`${this.baseURL}/health`);
    return response.data.data;
  }

  async getAuditLogs(params: {
    page?: number;
    page_size?: number;
    user_id?: string;
    action?: string;
    resource_type?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<PaginatedResponse<AuditLog>> {
    const response = await apiClient.get(`${this.baseURL}/audit-logs`, { params });
    return response.data.data;
  }

  async exportAuditLogs(params: {
    start_date: string;
    end_date: string;
    format?: 'csv' | 'excel' | 'json';
  }): Promise<Blob> {
    const response = await apiClient.get(`${this.baseURL}/audit-logs/export`, {
      params,
      responseType: 'blob',
    });
    return response.data;
  }
}

export const modelConfigApi = new ModelConfigApiClient();
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,包含提供商管理、模型配置、API密钥、路由规则、模型测试、统计监控、审计日志等7个模块的完整API接口定义
