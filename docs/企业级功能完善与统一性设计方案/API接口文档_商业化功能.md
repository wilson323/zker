# API接口文档：商业化功能模块

**文档编号**: DE-API-2025-DOC-030
**版本**: v1.0.0
**创建日期**: 2025-12-30
**参考设计**: 《30-商业化功能_订阅与付费.md》

---

## 📋 目录

1. [订阅管理API](#1-订阅管理api)
2. [Bot付费访问API](#2-bot付费访问api)
3. [支付订单API](#3-支付订单api)
4. [退款管理API](#4-退款管理api)
5. [收益统计API](#5-收益统计api)
6. [提现管理API](#6-提现管理api)
7. [优惠券API](#7-优惠券api)
8. [账单管理API](#8-账单管理api)

---

## 1. 订阅管理API

### 1.1 获取订阅套餐列表

**接口描述**: 获取所有可用的订阅套餐

**请求方式**: `GET /api/v1/subscription/plans`

**权限要求**: 无需认证或 `subscription:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| is_active | Boolean | 否 | 是否启用(默认true) |
| is_public | Boolean | 否 | 是否公开(默认true) |

**请求示例**:

```http
GET /api/v1/subscription/plans?is_active=true&is_public=true
Authorization: Bearer {access_token}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "免费版",
      "code": "FREE",
      "description": "体验基础功能",
      "pricing": {
        "monthly": 0,
        "quarterly": 0,
        "yearly": 0
      },
      "limits": {
        "max_bots": 3,
        "max_conversations_per_month": 1000,
        "max_knowledge_base_size_mb": 100,
        "max_api_calls_per_month": 10000
      },
      "features": ["basic_workflow", "community_support"],
      "support_level": "COMMUNITY",
      "is_active": true,
      "is_public": true,
      "sort_order": 1
    },
    {
      "id": 2,
      "name": "专业版",
      "code": "PRO",
      "description": "适合专业开发者",
      "pricing": {
        "monthly": 99.00,
        "quarterly": 280.00,
        "yearly": 990.00
      },
      "limits": {
        "max_bots": 50,
        "max_conversations_per_month": 50000,
        "max_knowledge_base_size_mb": 10240,
        "max_api_calls_per_month": 500000
      },
      "features": [
        "advanced_workflow",
        "webhook",
        "api_access",
        "email_support"
      ],
      "support_level": "EMAIL",
      "is_active": true,
      "is_public": true,
      "sort_order": 2
    }
  ]
}
```

### 1.2 获取当前订阅信息

**接口描述**: 获取当前用户的订阅状态和套餐信息

**请求方式**: `GET /api/v1/subscription/current`

**权限要求**: `subscription:read`

**请求示例**:

```http
GET /api/v1/subscription/current
Authorization: Bearer {access_token}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "user_id": 456,
    "plan": {
      "id": 2,
      "name": "专业版",
      "code": "PRO",
      "features": ["advanced_workflow", "webhook", "api_access"]
    },
    "status": "ACTIVE",
    "billing_cycle": "MONTHLY",
    "current_period_start": "2025-12-01T00:00:00Z",
    "current_period_end": "2026-01-01T00:00:00Z",
    "auto_renew": true,
    "usage": {
      "bots_count": 15,
      "conversations_this_month": 12500,
      "api_calls_this_month": 150000,
      "knowledge_base_size_mb": 512
    },
    "created_at": "2025-11-01T10:00:00Z",
    "updated_at": "2025-12-15T11:30:00Z"
  }
}
```

### 1.3 购买订阅套餐

**接口描述**: 购买指定订阅套餐,创建订单并返回支付参数

**请求方式**: `POST /api/v1/subscriptions/purchase`

**权限要求**: `subscription:write`

**请求体**:

```json
{
  "data": {
    "plan_id": 2,
    "billing_cycle": "YEARLY",
    "payment_method": "WECHAT_PAY",
    "coupon_code": "NEWYEAR2025"
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| plan_id | Long | ✅ | 套餐ID |
| billing_cycle | Enum | ✅ | 计费周期(MONTHLY/QUARTERLY/YEARLY) |
| payment_method | Enum | ✅ | 支付方式(WECHAT_PAY/ALIPAY/STRIPE) |
| coupon_code | String | 否 | 优惠券代码 |

**响应示例**:

```json
{
  "code": 0,
  "message": "Subscription purchase created",
  "data": {
    "subscription_id": 124,
    "order_id": "ORD20251230123456",
    "amount": 990.00,
    "discount": 99.00,
    "final_amount": 891.00,
    "currency": "CNY",
    "payment": {
      "method": "WECHAT_PAY",
      "provider": "WECHAT",
      "qr_code_url": "weixin://wxpay/bizpayurl?pr=xxxxx",
      "prepay_id": "wx26164353382146xxxxx",
      "code_url": "HTTPS://QR.ALIPAY.COM/xxxxx"
    },
    "expires_at": "2025-12-30T11:00:00Z"
  }
}
```

### 1.4 取消订阅

**接口描述**: 取消当前订阅,停止自动续费

**请求方式**: `POST /api/v1/subscriptions/{id}/cancel`

**权限要求**: `subscription:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | Long | 订阅ID |

**请求体**:

```json
{
  "data": {
    "reason": "不再需要",
    "cancel_immediately": false
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reason | String | 否 | 取消原因 |
| cancel_immediately | Boolean | 否 | 是否立即取消(默认false,周期结束后取消) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Subscription cancelled successfully",
  "data": {
    "id": 123,
    "status": "ACTIVE",
    "auto_renew": false,
    "expires_at": "2026-01-01T00:00:00Z",
    "will_expire_at_period_end": true
  }
}
```

### 1.5 订阅升降级

**接口描述**: 更改订阅套餐(升级或降级)

**请求方式**: `POST /api/v1/subscriptions/{id}/change-plan`

**权限要求**: `subscription:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | Long | 订阅ID |

**请求体**:

```json
{
  "data": {
    "target_plan_id": 3,
    "billing_cycle": "MONTHLY",
    "prorate": true
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_plan_id | Long | ✅ | 目标套餐ID |
| billing_cycle | Enum | ✅ | 新计费周期 |
| prorate | Boolean | 否 | 是否按比例计费(默认true) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Plan change initiated",
  "data": {
    "subscription_id": 123,
    "old_plan": {
      "id": 2,
      "name": "专业版"
    },
    "new_plan": {
      "id": 3,
      "name": "企业版"
    },
    "change_type": "UPGRADE",
    "prorated_amount": 900.00,
    "order_id": "ORD20251230234567",
    "effective_at": "2025-12-31T00:00:00Z"
  }
}
```

### 1.6 恢复订阅

**接口描述**: 恢复已取消的订阅

**请求方式**: `POST /api/v1/subscriptions/{id}/resume`

**权限要求**: `subscription:write`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| id | Long | 订阅ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Subscription resumed successfully",
  "data": {
    "id": 123,
    "status": "ACTIVE",
    "auto_renew": true,
    "resumed_at": "2025-12-30T10:30:00Z"
  }
}
```

---

## 2. Bot付费访问API

### 2.1 检查Bot访问权限

**接口描述**: 检查用户是否有权访问指定Bot

**请求方式**: `POST /api/v1/bots/{bot_id}/access/check`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |

**请求体**:

```json
{
  "data": {
    "check_type": "full"
  }
}
```

**响应示例** (需要付费):

```json
{
  "code": 0,
  "message": "Access check completed",
  "data": {
    "bot_id": 789,
    "has_access": false,
    "access_type": "PUBLIC_PAID",
    "pricing": {
      "model": "ONE_TIME",
      "amount": 9.90,
      "currency": "CNY",
      "original_price": 19.90,
      "discount": 50
    },
    "purchase_required": {
      "is_required": true,
      "purchase_url": "/api/v1/bots/789/purchase",
      "trial_available": false
    },
    "subscription_override": {
      "eligible_plans": ["PRO", "ENTERPRISE"],
      "has_subscription": false
    }
  }
}
```

**响应示例** (已有访问权限):

```json
{
  "code": 0,
  "message": "Access granted",
  "data": {
    "bot_id": 789,
    "has_access": true,
    "access_type": "PUBLIC_PAID",
    "access_source": "PURCHASE",
    "purchased_at": "2025-12-01T10:00:00Z",
    "expires_at": null
  }
}
```

### 2.2 购买Bot访问权限

**接口描述**: 购买Bot的一次性访问权限

**请求方式**: `POST /api/v1/bots/{bot_id}/purchase`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |

**请求体**:

```json
{
  "data": {
    "payment_method": "WECHAT_PAY",
    "coupon_code": "BOT50OFF"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Bot purchase created",
  "data": {
    "purchase_id": 456,
    "order_id": "ORD20251230345678",
    "bot": {
      "id": 789,
      "name": "高级客服助手",
      "avatar": "https://cdn.coze.com/bot/789/avatar.png"
    },
    "amount": 9.90,
    "discount": 10.00,
    "final_amount": 8.91,
    "currency": "CNY",
    "payment": {
      "method": "WECHAT_PAY",
      "qr_code_url": "weixin://wxpay/bizpayurl?pr=yyyyy",
      "expires_at": "2025-12-30T11:00:00Z"
    }
  }
}
```

### 2.3 设置Bot付费模式

**接口描述**: Bot开发者设置Bot的访问控制模式

**请求方式**: `PUT /api/v1/bots/{bot_id}/access-control`

**权限要求**: `bot:write`,必须是Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |

**请求体**:

```json
{
  "data": {
    "access_type": "PUBLIC_PAID",
    "pricing_model": "ONE_TIME",
    "pricing": {
      "amount": 9.90,
      "currency": "CNY",
      "trial_enabled": true,
      "trial_duration_days": 7,
      "trial_messages_limit": 10
    },
    "subscription_plans_allowed": ["PRO", "ENTERPRISE"],
    "revenue_share_opt_in": true
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| access_type | Enum | ✅ | 访问类型(PUBLIC_FREE/PUBLIC_PAID/SUBSCRIPTION_ONLY/PRIVATE) |
| pricing_model | Enum | ✅ | 定价模型(ONE_TIME/SUBSCRIPTION/USAGE_BASED/TIERED) |
| pricing | Object | ✅ | 定价配置 |
| subscription_plans_allowed | Array | 否 | 允许访问的订阅套餐列表 |
| revenue_share_opt_in | Boolean | 否 | 是否参与收益分成 |

**响应示例**:

```json
{
  "code": 0,
  "message": "Bot access control updated",
  "data": {
    "bot_id": 789,
    "access_type": "PUBLIC_PAID",
    "pricing_model": "ONE_TIME",
    "pricing": {
      "amount": 9.90,
      "currency": "CNY",
      "trial_enabled": true
    },
    "updated_at": "2025-12-30T10:30:00Z"
  }
}
```

### 2.4 获取Bot购买记录

**接口描述**: 获取用户购买的所有付费Bot

**请求方式**: `GET /api/v1/bot-purchases`

**权限要求**: 需要认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| bot_id | Long | 否 | 筛选指定Bot |
| status | String | 否 | 状态筛选(ACTIVE/EXPIRED) |

**请求示例**:

```http
GET /api/v1/bot-purchases?page=1&page_size=20&status=ACTIVE
Authorization: Bearer {access_token}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 456,
      "bot": {
        "id": 789,
        "name": "高级客服助手",
        "avatar": "https://cdn.coze.com/bot/789/avatar.png"
      },
      "pricing_model": "ONE_TIME",
      "amount_paid": 9.90,
      "purchased_at": "2025-12-01T10:00:00Z",
      "status": "ACTIVE",
      "expires_at": null,
      "usage_stats": {
        "conversations_count": 25,
        "messages_count": 342
      }
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 1,
      "total_pages": 1
    }
  }
}
```

---

## 3. 支付订单API

### 3.1 创建支付订单

**接口描述**: 创建支付订单并返回支付参数

**请求方式**: `POST /api/v1/payments/orders`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "order_type": "SUBSCRIPTION",
    "resource_type": "subscription",
    "resource_id": 124,
    "amount": 990.00,
    "currency": "CNY",
    "payment_method": "WECHAT_PAY",
    "client_ip": "123.45.67.89",
    "return_url": "https://coze.com/payment/return",
    "cancel_url": "https://coze.com/payment/cancel"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Order created",
  "data": {
    "order_id": "ORD20251230456789",
    "order_type": "SUBSCRIPTION",
    "amount": 990.00,
    "currency": "CNY",
    "status": "PENDING",
    "payment": {
      "method": "WECHAT_PAY",
      "qr_code_url": "weixin://wxpay/bizpayurl?pr=zzzzz",
      "prepay_id": "wx26164353382146zzzzz",
      "expires_at": "2025-12-30T11:00:00Z"
    },
    "created_at": "2025-12-30T10:00:00Z"
  }
}
```

### 3.2 查询订单状态

**接口描述**: 查询订单的支付状态

**请求方式**: `GET /api/v1/payments/orders/{order_id}`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| order_id | String | 订单ID |

**响应示例** (支付成功):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "order_id": "ORD20251230456789",
    "order_type": "SUBSCRIPTION",
    "amount": 990.00,
    "discount": 99.00,
    "final_amount": 891.00,
    "currency": "CNY",
    "status": "PAID",
    "payment_method": "WECHAT_PAY",
    "paid_at": "2025-12-30T10:15:00Z",
    "transaction_id": "4200001234567890123456789",
    "subscription": {
      "id": 124,
      "status": "ACTIVE",
      "current_period_end": "2026-01-30T10:15:00Z"
    },
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-30T10:15:00Z"
  }
}
```

### 3.3 取消订单

**接口描述**: 取消待支付订单

**请求方式**: `POST /api/v1/payments/orders/{order_id}/cancel`

**权限要求**: 需要认证,订单所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| order_id | String | 订单ID |

**请求体**:

```json
{
  "data": {
    "reason": "放弃购买"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Order cancelled",
  "data": {
    "order_id": "ORD20251230456789",
    "status": "CANCELLED",
    "cancelled_at": "2025-12-30T10:05:00Z",
    "reason": "放弃购买"
  }
}
```

### 3.4 获取订单列表

**接口描述**: 获取用户的订单列表

**请求方式**: `GET /api/v1/payments/orders`

**权限要求**: 需要认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| status | String | 否 | 状态筛选(PENDING/PAID/FAILED/CANCELLED) |
| order_type | String | 否 | 订单类型(SUBSCRIPTION/BOT_PURCHASE) |
| start_date | String | 否 | 开始日期(ISO8601) |
| end_date | String | 否 | 结束日期(ISO8601) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "order_id": "ORD20251230456789",
      "order_type": "SUBSCRIPTION",
      "amount": 990.00,
      "final_amount": 891.00,
      "currency": "CNY",
      "status": "PAID",
      "paid_at": "2025-12-30T10:15:00Z",
      "created_at": "2025-12-30T10:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 15,
      "total_pages": 1
    }
  }
}
```

---

## 4. 退款管理API

### 4.1 创建退款申请(管理员)

**接口描述**: 管理员为订单创建退款申请

**请求方式**: `POST /api/v1/admin/refunds`

**权限要求**: `refund:write`

**请求体**:

```json
{
  "data": {
    "order_id": "ORD20251230456789",
    "refund_amount": 891.00,
    "reason": "用户申请退款",
    "refund_type": "FULL"
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_id | String | ✅ | 订单ID |
| refund_amount | Decimal | ✅ | 退款金额 |
| reason | String | ✅ | 退款原因 |
| refund_type | Enum | 否 | 退款类型(FULL/PARTIAL) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Refund created",
  "data": {
    "refund_id": "REF20251230123456",
    "order_id": "ORD20251230456789",
    "refund_amount": 891.00,
    "currency": "CNY",
    "status": "PENDING",
    "reason": "用户申请退款",
    "created_at": "2025-12-30T11:00:00Z"
  }
}
```

### 4.2 查询退款状态

**接口描述**: 查询退款申请的状态

**请求方式**: `GET /api/v1/refunds/{refund_id}`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| refund_id | String | 退款ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "refund_id": "REF20251230123456",
    "order_id": "ORD20251230456789",
    "refund_amount": 891.00,
    "currency": "CNY",
    "status": "SUCCESS",
    "reason": "用户申请退款",
    "refunded_at": "2025-12-30T11:30:00Z",
    "transaction_id": "REF4200001234567890",
    "subscription_cancelled": true,
    "created_at": "2025-12-30T11:00:00Z",
    "updated_at": "2025-12-30T11:30:00Z"
  }
}
```

### 4.3 获取退款列表

**接口描述**: 获取退款申请列表

**请求方式**: `GET /api/v1/admin/refunds`

**权限要求**: `refund:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| status | String | 否 | 状态筛选(PENDING/SUCCESS/FAILED) |
| start_date | String | 否 | 开始日期 |
| end_date | String | 否 | 结束日期 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "refund_id": "REF20251230123456",
      "order_id": "ORD20251230456789",
      "user": {
        "id": 456,
        "name": "张三",
        "email": "zhangsan@example.com"
      },
      "refund_amount": 891.00,
      "currency": "CNY",
      "status": "PENDING",
      "reason": "用户申请退款",
      "created_at": "2025-12-30T11:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 5,
      "total_pages": 1
    }
  }
}
```

---

## 5. 收益统计API

### 5.1 获取收益概览

**接口描述**: Bot开发者查看收益概览

**请求方式**: `GET /api/v1/revenue/overview`

**权限要求**: `revenue:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | String | 否 | 统计周期(today/week/month/year/all) |

**请求示例**:

```http
GET /api/v1/revenue/overview?period=month
Authorization: Bearer {access_token}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "month",
    "total_revenue": 5000.00,
    "platform_fee": 500.00,
    "net_revenue": 4500.00,
    "available_balance": 3500.00,
    "pending_balance": 1000.00,
    "withdrawn_amount": 2000.00,
    "currency": "CNY",
    "stats": {
      "total_sales": 125,
      "total_bots_sold": 8,
      "average_order_value": 40.00,
      "revenue_change_percent": 25.5
    },
    "breakdown": {
      "by_bot": [
        {
          "bot_id": 789,
          "bot_name": "高级客服助手",
          "sales_count": 50,
          "revenue": 2000.00
        }
      ],
      "by_date": [
        {
          "date": "2025-12-29",
          "revenue": 500.00
        }
      ]
    }
  }
}
```

### 5.2 获取Bot收益详情

**接口描述**: 查看指定Bot的收益详情

**请求方式**: `GET /api/v1/bots/{bot_id}/revenue`

**权限要求**: `revenue:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | String | 否 | 开始日期 |
| end_date | String | 否 | 结束日期 |
| granularity | String | 否 | 时间粒度(day/week/month) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": 789,
    "bot_name": "高级客服助手",
    "total_revenue": 2000.00,
    "total_sales": 50,
    "net_revenue": 1800.00,
    "platform_fee": 200.00,
    "currency": "CNY",
    "daily_breakdown": [
      {
        "date": "2025-12-29",
        "sales": 10,
        "revenue": 200.00,
        "net_revenue": 180.00
      }
    ],
    "recent_transactions": [
      {
        "transaction_id": "TXN20251229001",
        "order_id": "ORD20251229001",
        "amount": 9.90,
        "net_amount": 8.91,
        "timestamp": "2025-12-29T15:30:00Z"
      }
    ]
  }
}
```

### 5.3 获取收益报表

**接口描述**: 导出收益报表

**请求方式**: `GET /api/v1/revenue/report`

**权限要求**: `revenue:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | String | ✅ | 开始日期 |
| end_date | String | ✅ | 结束日期 |
| format | String | 否 | 报表格式(csv/excel/pdf) |
| bot_id | Long | 否 | 筛选指定Bot |

**响应示例**:

```http
HTTP/1.1 200 OK
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="revenue_report_2025-12-01_to_2025-12-30.xlsx"

[二进制文件内容]
```

---

## 6. 提现管理API

### 6.1 创建提现申请

**接口描述**: Bot开发者创建提现申请

**请求方式**: `POST /api/v1/revenue/withdrawals`

**权限要求**: `revenue:write`

**请求体**:

```json
{
  "data": {
    "amount": 1000.00,
    "withdrawal_method": "BANK_TRANSFER",
    "bank_account": {
      "bank_name": "中国工商银行",
      "account_number": "6222021234567890123",
      "account_holder_name": "张三",
      "bank_branch": "北京分行朝阳支行"
    },
    "notes": "月度提现"
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| amount | Decimal | ✅ | 提现金额 |
| withdrawal_method | Enum | ✅ | 提现方式(BANK_TRANSFER/ALIPAY/WECHAT) |
| bank_account | Object | ✅ | 银行账户信息 |
| notes | String | 否 | 备注 |

**响应示例**:

```json
{
  "code": 0,
  "message": "Withdrawal request created",
  "data": {
    "withdrawal_id": "WTH20251230123456",
    "amount": 1000.00,
    "currency": "CNY",
    "status": "PENDING",
    "withdrawal_method": "BANK_TRANSFER",
    "bank_account": {
      "bank_name": "中国工商银行",
      "account_number": "6222021234****123",
      "account_holder_name": "张三"
    },
    "estimated_processing_days": 3,
    "created_at": "2025-12-30T12:00:00Z"
  }
}
```

### 6.2 获取提现记录

**接口描述**: 获取提现申请记录列表

**请求方式**: `GET /api/v1/revenue/withdrawals`

**权限要求**: `revenue:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| status | String | 否 | 状态筛选(PENDING/APPROVED/REJECTED/COMPLETED) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "withdrawal_id": "WTH20251230123456",
      "amount": 1000.00,
      "currency": "CNY",
      "status": "PENDING",
      "withdrawal_method": "BANK_TRANSFER",
      "bank_account": {
        "bank_name": "中国工商银行",
        "account_number": "6222021234****123"
      },
      "created_at": "2025-12-30T12:00:00Z"
    },
    {
      "withdrawal_id": "WTH20251220111222",
      "amount": 1500.00,
      "currency": "CNY",
      "status": "COMPLETED",
      "withdrawal_method": "BANK_TRANSFER",
      "approved_at": "2025-12-21T10:00:00Z",
      "completed_at": "2025-12-23T15:30:00Z",
      "created_at": "2025-12-20T14:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 2,
      "total_pages": 1
    }
  }
}
```

### 6.3 审批提现申请(管理员)

**接口描述**: 管理员审批提现申请

**请求方式**: `POST /api/v1/admin/withdrawals/{withdrawal_id}/approve`

**权限要求**: `withdrawal:approve`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| withdrawal_id | String | 提现ID |

**请求体**:

```json
{
  "data": {
    "action": "APPROVE",
    "notes": "审核通过",
    "receipt_reference": "TXN2025123099999"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Withdrawal approved",
  "data": {
    "withdrawal_id": "WTH20251230123456",
    "status": "APPROVED",
    "amount": 1000.00,
    "currency": "CNY",
    "approved_at": "2025-12-30T14:00:00Z",
    "approved_by": {
      "admin_id": 1,
      "name": "Admin User"
    },
    "notes": "审核通过",
    "estimated_completion_at": "2025-12-31T14:00:00Z"
  }
}
```

### 6.4 拒绝提现申请(管理员)

**接口描述**: 管理员拒绝提现申请

**请求方式**: `POST /api/v1/admin/withdrawals/{withdrawal_id}/reject`

**权限要求**: `withdrawal:approve`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| withdrawal_id | String | 提现ID |

**请求体**:

```json
{
  "data": {
    "reason": "银行账户信息不完整",
    "notes": "请补全开户行信息后重新申请"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Withdrawal rejected",
  "data": {
    "withdrawal_id": "WTH20251230123456",
    "status": "REJECTED",
    "amount": 1000.00,
    "rejected_at": "2025-12-30T14:00:00Z",
    "rejected_by": {
      "admin_id": 1,
      "name": "Admin User"
    },
    "reason": "银行账户信息不完整"
  }
}
```

---

## 7. 优惠券API

### 7.1 验证优惠券

**接口描述**: 验证优惠券代码是否有效

**请求方式**: `POST /api/v1/coupons/validate`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "coupon_code": "NEWYEAR2025",
    "order_type": "SUBSCRIPTION",
    "order_amount": 990.00
  }
}
```

**响应示例** (优惠券有效):

```json
{
  "code": 0,
  "message": "Coupon valid",
  "data": {
    "coupon_code": "NEWYEAR2025",
    "coupon_type": "PERCENTAGE",
    "discount_value": 10,
    "discount_amount": 99.00,
    "final_amount": 891.00,
    "currency": "CNY",
    "valid_from": "2025-01-01T00:00:00Z",
    "valid_until": "2025-12-31T23:59:59Z",
    "usage_limit": {
      "max_uses": 1000,
      "used_count": 234,
      "remaining_uses": 766,
      "user_max_uses": 1,
      "user_used_count": 0
    },
    "applicable_products": ["SUBSCRIPTION", "BOT_PURCHASE"],
    "min_order_amount": 50.00,
    "is_valid": true
  }
}
```

**响应示例** (优惠券无效):

```json
{
  "code": 42201,
  "message": "Coupon validation failed",
  "errors": [
    {
      "field": "coupon_code",
      "message": "Coupon code expired"
    }
  ]
}
```

### 7.2 获取可用优惠券列表

**接口描述**: 获取当前用户可用的优惠券

**请求方式**: `GET /api/v1/coupons/available`

**权限要求**: 需要认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| order_type | String | 否 | 订单类型筛选 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "coupon_code": "NEWYEAR2025",
      "coupon_name": "新年优惠",
      "coupon_type": "PERCENTAGE",
      "discount_value": 10,
      "description": "全场9折优惠",
      "valid_until": "2025-12-31T23:59:59Z",
      "applicable_products": ["SUBSCRIPTION", "BOT_PURCHASE"]
    },
    {
      "coupon_code": "WELCOME10",
      "coupon_name": "新人优惠",
      "coupon_type": "FIXED",
      "discount_value": 10.00,
      "description": "新用户立减10元",
      "valid_until": "2025-06-30T23:59:59Z",
      "applicable_products": ["SUBSCRIPTION"]
    }
  ]
}
```

### 7.3 创建优惠券(管理员)

**接口描述**: 管理员创建优惠券

**请求方式**: `POST /api/v1/admin/coupons`

**权限要求**: `coupon:write`

**请求体**:

```json
{
  "data": {
    "coupon_code": "NEWSALE2025",
    "coupon_name": "春季促销",
    "coupon_type": "PERCENTAGE",
    "discount_value": 15,
    "valid_from": "2025-03-01T00:00:00Z",
    "valid_until": "2025-03-31T23:59:59Z",
    "usage_limits": {
      "max_uses": 500,
      "user_max_uses": 1
    },
    "applicable_products": ["SUBSCRIPTION", "BOT_PURCHASE"],
    "min_order_amount": 100.00,
    "description": "春季85折优惠"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Coupon created",
  "data": {
    "id": 123,
    "coupon_code": "NEWSALE2025",
    "coupon_name": "春季促销",
    "coupon_type": "PERCENTAGE",
    "discount_value": 15,
    "status": "ACTIVE",
    "created_at": "2025-12-30T15:00:00Z"
  }
}
```

---

## 8. 账单管理API

### 8.1 获取账单列表

**接口描述**: 获取用户的账单列表

**请求方式**: `GET /api/v1/billing/invoices`

**权限要求**: `billing:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| status | String | 否 | 状态筛选(PAID/UNPAID/OVERDUE) |
| start_date | String | 否 | 开始日期 |
| end_date | String | 否 | 结束日期 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "invoice_id": "INV202512-001",
      "subscription_id": 123,
      "billing_period": {
        "start": "2025-12-01T00:00:00Z",
        "end": "2026-01-01T00:00:00Z"
      },
      "amount": 990.00,
      "discount": 99.00,
      "tax": 0.00,
      "total_amount": 891.00,
      "currency": "CNY",
      "status": "PAID",
      "paid_at": "2025-12-01T10:15:00Z",
      "due_date": "2025-12-15T23:59:59Z",
      "download_url": "https://cdn.coze.com/invoices/INV202512-001.pdf"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 12,
      "total_pages": 1
    }
  }
}
```

### 8.2 获取账单详情

**接口描述**: 获取账单详细信息

**请求方式**: `GET /api/v1/billing/invoices/{invoice_id}`

**权限要求**: `billing:read`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| invoice_id | String | 账单ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "invoice_id": "INV202512-001",
    "subscription": {
      "id": 123,
      "plan_name": "专业版",
      "billing_cycle": "YEARLY"
    },
    "billing_period": {
      "start": "2025-12-01T00:00:00Z",
      "end": "2026-01-01T00:00:00Z"
    },
    "line_items": [
      {
        "description": "专业版 - 年付",
        "quantity": 1,
        "unit_price": 990.00,
        "amount": 990.00
      }
    ],
    "discount": {
      "coupon_code": "NEWYEAR2025",
      "discount_type": "PERCENTAGE",
      "discount_value": 10,
      "amount": 99.00
    },
    "subtotal": 990.00,
    "tax": 0.00,
    "total_amount": 891.00,
    "currency": "CNY",
    "status": "PAID",
    "paid_at": "2025-12-01T10:15:00Z",
    "due_date": "2025-12-15T23:59:59Z",
    "payment_method": "WECHAT_PAY",
    "transaction_id": "4200001234567890123456789",
    "created_at": "2025-12-01T00:00:00Z"
  }
}
```

### 8.3 下载账单PDF

**接口描述**: 下载账单PDF文件

**请求方式**: `GET /api/v1/billing/invoices/{invoice_id}/download`

**权限要求**: `billing:read`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| invoice_id | String | 账单ID |

**响应示例**:

```http
HTTP/1.1 200 OK
Content-Type: application/pdf
Content-Disposition: attachment; filename="INV202512-001.pdf"

[PDF二进制内容]
```

---

## 附录

### A. 错误码参考

| 错误码 | 说明 |
|--------|------|
| 40001 | 请求参数错误 |
| 40101 | Token缺失 |
| 40103 | Token过期 |
| 40301 | 无权限访问 |
| 40401 | 资源不存在 |
| 40901 | 资源已存在(如重复购买) |
| 42201 | 业务规则验证失败 |
| 42202 | 订阅已过期 |
| 42203 | 配额已用尽 |
| 42901 | 超过API调用频率限制 |

### B. Webhook事件类型

支付Webhook事件通知:

```json
{
  "event": "payment.success",
  "timestamp": "2025-12-30T10:15:00Z",
  "data": {
    "order_id": "ORD20251230456789",
    "payment": {
      "method": "WECHAT_PAY",
      "amount": 891.00,
      "currency": "CNY",
      "transaction_id": "4200001234567890123456789"
    },
    "subscription": {
      "id": 124,
      "status": "ACTIVE"
    }
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

// SubscriptionHandler 订阅处理器
type SubscriptionHandler struct {
    subscriptionService SubscriptionService
    paymentService       PaymentService
}

// GetSubscriptionPlans 获取订阅套餐列表
func (h *SubscriptionHandler) GetSubscriptionPlans(ctx context.Context, c *app.RequestContext) {
    var req GetPlansRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    plans, err := h.subscriptionService.GetPlans(ctx, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "success",
        Data:    plans,
    })
}

// PurchaseSubscription 购买订阅
func (h *SubscriptionHandler) PurchaseSubscription(ctx context.Context, c *app.RequestContext) {
    var req PurchaseSubscriptionRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    userID := getUserIDFromContext(ctx)

    result, err := h.subscriptionService.Purchase(ctx, userID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "Subscription purchase created",
        Data:    result,
    })
}
```

### D. 前端TypeScript类型定义

```typescript
// types/commercial.ts

// 订阅相关类型
interface SubscriptionPlan {
  id: number;
  name: string;
  code: 'FREE' | 'PRO' | 'ENTERPRISE';
  pricing: {
    monthly: number;
    quarterly: number;
    yearly: number;
  };
  limits: {
    max_bots: number;
    max_conversations_per_month: number;
    max_knowledge_base_size_mb: number;
    max_api_calls_per_month: number;
  };
  features: string[];
}

interface Subscription {
  id: number;
  user_id: number;
  plan: SubscriptionPlan;
  status: 'ACTIVE' | 'CANCELLED' | 'EXPIRED';
  billing_cycle: 'MONTHLY' | 'QUARTERLY' | 'YEARLY';
  current_period_start: string;
  current_period_end: string;
  auto_renew: boolean;
}

// 订单相关类型
interface Order {
  order_id: string;
  order_type: 'SUBSCRIPTION' | 'BOT_PURCHASE';
  amount: number;
  final_amount: number;
  currency: string;
  status: 'PENDING' | 'PAID' | 'FAILED' | 'CANCELLED';
  payment_method: 'WECHAT_PAY' | 'ALIPAY' | 'STRIPE';
}

// 收益相关类型
interface RevenueOverview {
  period: string;
  total_revenue: number;
  platform_fee: number;
  net_revenue: number;
  available_balance: number;
  pending_balance: number;
  withdrawn_amount: number;
}

// API客户端方法
class CommercialApiClient {
  async purchaseSubscription(
    planId: number,
    billingCycle: string,
    paymentMethod: string,
    couponCode?: string
  ): Promise<Order> {
    const response = await apiClient.post('/subscriptions/purchase', {
      data: {
        plan_id: planId,
        billing_cycle: billingCycle,
        payment_method: paymentMethod,
        coupon_code: couponCode,
      },
    });
    return response.data.data;
  }

  async getRevenueOverview(period?: string): Promise<RevenueOverview> {
    const response = await apiClient.get('/revenue/overview', {
      params: { period },
    });
    return response.data.data;
  }

  async createWithdrawal(
    amount: number,
    bankAccount: BankAccount
  ): Promise<Withdrawal> {
    const response = await apiClient.post('/revenue/withdrawals', {
      data: {
        amount,
        withdrawal_method: 'BANK_TRANSFER',
        bank_account: bankAccount,
      },
    });
    return response.data.data;
  }
}
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,包含订阅管理、Bot付费访问、支付订单、退款管理、收益统计、提现管理、优惠券、账单管理等8个模块的完整API接口定义
