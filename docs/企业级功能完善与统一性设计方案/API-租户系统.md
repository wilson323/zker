# 租户系统 API 文档

## 概述

本文档描述了 ZKER 企业级多租户 SaaS 平台的租户管理系统 API，包括租户管理、订阅管理、配额管理和计费统计。

## 基础信息

- **Base URL**: `https://api.zker.com/v1`
- **Content-Type**: `application/json`
- **认证方式**: Bearer Token (JWT)

## 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

### 错误响应
```json
{
  "code": "TENANT_NOT_FOUND",
  "message": "租户不存在",
  "details": {
    "tenant_id": "tenant_123"
  },
  "request_id": "req_abc123",
  "timestamp": 1704067200000
}
```

---

## 1. 租户管理

### 1.1 创建租户

**接口**: `POST /tenants`

**请求参数**:
```json
{
  "tenant_name": "示例企业",
  "tenant_type": "enterprise",
  "contact_email": "admin@example.com",
  "contact_phone": "+86-13800138000"
}
```

**字段说明**:
- `tenant_name`: 租户名称（必填，1-200字符）
- `tenant_type`: 租户类型（必填）
  - `individual`: 个人
  - `team`: 团队
  - `enterprise`: 企业
- `contact_email`: 联系邮箱（必填，需符合邮箱格式）
- `contact_phone`: 联系电话（可选）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_abc123",
    "tenant_name": "示例企业",
    "tenant_type": "enterprise",
    "status": "active",
    "subscription_tier": "free",
    "created_at": 1704067200000,
    "updated_at": 1704067200000
  }
}
```

### 1.2 获取租户信息

**接口**: `GET /tenants/{tenant_id}`

**路径参数**:
- `tenant_id`: 租户ID

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_abc123",
    "tenant_name": "示例企业",
    "tenant_type": "enterprise",
    "status": "active",
    "subscription_tier": "pro",
    "quotas": [
      {
        "resource_type": "bots",
        "max_limit": 100,
        "used_count": 25,
        "reset_cycle": "monthly"
      }
    ],
    "created_at": 1704067200000,
    "updated_at": 1704067200000
  }
}
```

### 1.3 更新租户信息

**接口**: `PUT /tenants/{tenant_id}`

**请求参数**:
```json
{
  "tenant_name": "新企业名称",
  "contact_email": "newadmin@example.com",
  "contact_phone": "+86-13900139000"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_abc123",
    "tenant_name": "新企业名称",
    "updated_at": 1704153600000
  }
}
```

### 1.4 升级订阅等级

**接口**: `POST /tenants/{tenant_id}/upgrade`

**请求参数**:
```json
{
  "target_tier": "enterprise"
}
```

**字段说明**:
- `target_tier`: 目标订阅等级（必填）
  - `free`: 免费版
  - `pro`: 专业版
  - `enterprise`: 企业版

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_abc123",
    "subscription_tier": "enterprise",
    "quotas": [
      {
        "resource_type": "bots",
        "max_limit": 1000,
        "used_count": 25
      }
    ]
  }
}
```

### 1.5 列出租户

**接口**: `GET /tenants`

**查询参数**:
- `status`: 租户状态（可选）- `active`, `suspended`, `deleted`
- `tenant_type`: 租户类型（可选）- `individual`, `team`, `enterprise`
- `subscription_tier`: 订阅等级（可选）- `free`, `pro`, `enterprise`
- `page_token`: 分页token（可选）
- `page_size`: 每页数量（可选，默认20，最大100）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenants": [
      {
        "tenant_id": "tenant_abc123",
        "tenant_name": "示例企业",
        "status": "active",
        "subscription_tier": "pro"
      }
    ],
    "total_count": 100,
    "next_page_token": "page_token_xyz"
  }
}
```

---

## 2. 订阅管理

### 2.1 获取订阅信息

**接口**: `GET /tenants/{tenant_id}/subscription`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "subscription_id": "sub_abc123",
    "tenant_id": "tenant_abc123",
    "plan_tier": "pro",
    "billing_cycle": "monthly",
    "status": "active",
    "started_at": 1704067200000,
    "expires_at": 1706745600000,
    "auto_renew": true
  }
}
```

**字段说明**:
- `plan_tier`: 订阅等级
- `billing_cycle`: 计费周期
  - `monthly`: 月付
  - `yearly`: 年付
- `status`: 订阅状态
  - `active`: 活跃
  - `suspended`: 暂停
  - `cancelled`: 已取消

### 2.2 更新订阅

**接口**: `PUT /tenants/{tenant_id}/subscription`

**请求参数**:
```json
{
  "billing_cycle": "yearly",
  "auto_renew": true
}
```

---

## 3. 配额管理

### 3.1 获取配额状态

**接口**: `GET /tenants/{tenant_id}/quotas`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "quotas": [
      {
        "resource_type": "bots",
        "max_limit": 100,
        "used_count": 25,
        "remaining": 75,
        "usage_percent": 25.0,
        "reset_cycle": "monthly",
        "last_reset_at": 1704067200000,
        "alert_level": "normal"
      },
      {
        "resource_type": "messages",
        "max_limit": 10000,
        "used_count": 8500,
        "remaining": 1500,
        "usage_percent": 85.0,
        "reset_cycle": "daily",
        "last_reset_at": 1704067200000,
        "alert_level": "warning"
      }
    ]
  }
}
```

**资源类型**:
- `bots`: Bot数量
- `workflows`: 工作流数量
- `messages`: 消息数量
- `storage`: 存储空间（字节）
- `team_members`: 团队成员数

**告警级别**:
- `normal`: 正常（< 80%）
- `warning`: 警告（80% - 95%）
- `critical`: 严重（95% - 100%）
- `exceeded`: 超限（> 100%）

### 3.2 检查配额

**接口**: `POST /tenants/{tenant_id}/quotas/check`

**请求参数**:
```json
{
  "resource_type": "bots",
  "required_count": 5
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "allowed": true,
    "current_usage": 25,
    "max_limit": 100,
    "remaining": 75
  }
}
```

**错误示例**（配额不足）:
```json
{
  "code": "QUOTA_EXCEEDED",
  "message": "配额不足",
  "details": {
    "tenant_id": "tenant_abc123",
    "resource_type": "bots",
    "current_usage": 98,
    "max_limit": 100,
    "required_count": 5
  }
}
```

### 3.3 消费配额

**接口**: `POST /tenants/{tenant_id}/quotas/consume`

**请求参数**:
```json
{
  "resource_type": "messages",
  "count": 100
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "resource_type": "messages",
    "previous_count": 8500,
    "current_count": 8600,
    "remaining": 1400
  }
}
```

### 3.4 回滚配额

**接口**: `POST /tenants/{tenant_id}/quotas/rollback`

**请求参数**:
```json
{
  "resource_type": "messages",
  "count": 100
}
```

---

## 4. 计费管理

### 4.1 记录使用量

**接口**: `POST /tenants/{tenant_id}/usage`

**请求参数**:
```json
{
  "resource_type": "messages",
  "action": "create",
  "quantity": 100,
  "metadata": {
    "bot_id": "bot_123",
    "conversation_id": "conv_456"
  }
}
```

**操作类型**:
- `create`: 创建
- `update`: 更新
- `delete`: 删除
- `query`: 查询

### 4.2 获取使用量汇总

**接口**: `GET /tenants/{tenant_id}/usage/summary`

**查询参数**:
- `start_date`: 开始日期（必填，格式：YYYY-MM-DD）
- `end_date`: 结束日期（必填，格式：YYYY-MM-DD）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "summaries": [
      {
        "resource_type": "messages",
        "total_quantity": 50000,
        "action_counts": {
          "create": 30000,
          "update": 15000,
          "delete": 5000
        },
        "first_usage": "2024-01-01T00:00:00Z",
        "last_usage": "2024-01-31T23:59:59Z"
      }
    ]
  }
}
```

### 4.3 生成账单

**接口**: `POST /tenants/{tenant_id}/invoices/generate`

**请求参数**:
```json
{
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "invoice_id": "inv_abc123",
    "tenant_id": "tenant_abc123",
    "billing_cycle": "monthly",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "total_usage": 50000,
    "total_amount": 50.00,
    "currency": "CNY",
    "status": "pending",
    "due_date": "2024-02-15T00:00:00Z",
    "created_at": "2024-02-01T00:00:00Z"
  }
}
```

### 4.4 获取账单详情

**接口**: `GET /tenants/{tenant_id}/invoices/{invoice_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "invoice_id": "inv_abc123",
    "status": "paid",
    "total_amount": 50.00,
    "paid_at": "2024-02-05T10:30:00Z"
  }
}
```

### 4.5 列出账单

**接口**: `GET /tenants/{tenant_id}/invoices`

**查询参数**:
- `status`: 账单状态（可选）- `pending`, `paid`, `overdue`, `cancelled`
- `limit`: 返回数量（可选，默认20，最大100）
- `offset`: 偏移量（可选，默认0）

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "invoices": [
      {
        "invoice_id": "inv_abc123",
        "billing_cycle": "monthly",
        "total_amount": 50.00,
        "status": "paid",
        "created_at": "2024-02-01T00:00:00Z"
      }
    ],
    "total_count": 12
  }
}
```

### 4.6 支付账单

**接口**: `POST /tenants/{tenant_id}/invoices/{invoice_id}/pay`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "invoice_id": "inv_abc123",
    "status": "paid",
    "paid_at": 1704494400000
  }
}
```

---

## 5. 配额监控

### 5.1 获取配额监控状态

**接口**: `GET /tenants/{tenant_id}/monitoring`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant_abc123",
    "quotas": [
      {
        "resource_type": "bots",
        "max_limit": 100,
        "used_count": 95,
        "usage_percent": 95.0,
        "alert_level": "critical",
        "should_reset": false,
        "last_reset_at": 1704067200000
      }
    ],
    "alerts": [
      {
        "resource_type": "bots",
        "alert_type": "critical",
        "message": "配额使用率已达95%"
      }
    ]
  }
}
```

---

## 错误码列表

| 错误码 | 说明 | HTTP状态码 |
|--------|------|-----------|
| `TENANT_NOT_FOUND` | 租户不存在 | 404 |
| `TENANT_ALREADY_EXISTS` | 租户已存在 | 409 |
| `TENANT_NAME_INVALID` | 租户名称无效 | 400 |
| `QUOTA_EXCEEDED` | 配额超限 | 403 |
| `QUOTA_NOT_FOUND` | 配额不存在 | 404 |
| `SUBSCRIPTION_NOT_FOUND` | 订阅不存在 | 404 |
| `SUBSCRIPTION_EXPIRED` | 订阅已过期 | 403 |
| `INVOICE_NOT_FOUND` | 账单不存在 | 404 |
| `INVOICE_ALREADY_PAID` | 账单已支付 | 409 |
| `INVALID_REQUEST` | 请求参数无效 | 400 |
| `INTERNAL_ERROR` | 内部错误 | 500 |

---

## 使用示例

### cURL 示例

#### 创建租户
```bash
curl -X POST https://api.zker.com/v1/tenants \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_name": "示例企业",
    "tenant_type": "enterprise",
    "contact_email": "admin@example.com"
  }'
```

#### 获取配额状态
```bash
curl -X GET https://api.zker.com/v1/tenants/tenant_abc123/quotas \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 检查配额
```bash
curl -X POST https://api.zker.com/v1/tenants/tenant_abc123/quotas/check \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "resource_type": "bots",
    "required_count": 5
  }'
```

### JavaScript (Axios) 示例

```javascript
// 创建租户
const createTenant = async () => {
  try {
    const response = await axios.post(
      'https://api.zker.com/v1/tenants',
      {
        tenant_name: '示例企业',
        tenant_type: 'enterprise',
        contact_email: 'admin@example.com'
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      }
    );
    return response.data;
  } catch (error) {
    console.error('创建租户失败:', error.response.data);
  }
};

// 检查配额
const checkQuota = async (tenantId, resourceType, requiredCount) => {
  try {
    const response = await axios.post(
      `https://api.zker.com/v1/tenants/${tenantId}/quotas/check`,
      {
        resource_type: resourceType,
        required_count: requiredCount
      },
      {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      }
    );
    return response.data;
  } catch (error) {
    if (error.response?.data?.code === 'QUOTA_EXCEEDED') {
      console.error('配额不足:', error.response.data.details);
    }
  }
};
```

---

## 附录

### A. 订阅等级对比

| 特性 | Free | Pro | Enterprise |
|------|------|-----|------------|
| Bot数量 | 3 | 100 | 1000 |
| 工作流数量 | 5 | 100 | 1000 |
| 消息数/月 | 1,000 | 100,000 | 1,000,000 |
| 存储空间 | 1GB | 100GB | 1TB |
| 团队成员 | 1 | 10 | 100 |
| 价格 | 免费 | ¥99/月 | ¥999/月 |

### B. 配额重置策略

| 资源类型 | 重置周期 |
|---------|---------|
| bots | monthly（每月1日） |
| workflows | monthly（每月1日） |
| messages | daily（每日0点） |
| storage | never（永不重置） |
| team_members | never（永不重置） |

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
**维护者**: 研发A团队
