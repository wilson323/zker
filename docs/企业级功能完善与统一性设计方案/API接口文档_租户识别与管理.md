# API接口文档:租户识别与管理模块

> **版本**: v1.0.0
> **创建日期**: 2025-12-30
> **对应设计文档**: [21-MultiTenant_SaaS核心_租户识别与管理.md](./详细设计/21-MultiTenant_SaaS核心_租户识别与管理.md)

---

## 📋 目录

1. [概述](#1-概述)
2. [租户注册API](#2-租户注册api)
3. [租户信息管理API](#3-租户信息管理api)
4. [租户识别API](#4-租户识别api)
5. [租户配额管理API](#5-租户配额管理api)
6. [租户套餐管理API](#6-租户套餐管理api)
7. [租户生命周期管理API](#7-租户生命周期管理api)
8. [租户管理后台API](#8-租户管理后台api)
9. [错误代码](#9-错误代码)
10. [代码示例](#10-代码示例)

---

## 1. 概述

### 1.1 模块说明

租户识别与管理是Multi-Tenant SaaS系统的核心基础模块,提供租户注册、识别、信息管理、配额管理、套餐管理和生命周期管理功能。

### 1.2 基础信息

**Base URL**: `/api/v1`

**认证方式**:
- 租户注册相关接口: 无需认证
- 租户管理接口: JWT Bearer Token

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {token}  // 需要认证的接口
X-Tenant-ID: {tenant_id}      // API调用时传递租户ID
```

### 1.3 统一响应格式

参见 [API设计规范文档](./API设计规范文档.md#3-统一响应结构)

---

## 2. 租户注册API

### 2.1 检查企业名称是否可用

检查企业名称是否已被注册。

**接口描述**: 验证企业名称的唯一性

**请求方式**: `GET /tenants/check-company-name`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| company_name | string | ✅ | 企业名称 |

**请求示例**:
```http
GET /api/v1/tenants/check-company-name?company_name=示例企业
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "available": true,
    "suggestions": []
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_001"
}
```

**企业名称已存在时** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "available": false,
    "suggestions": [
      "示例企业(科技)",
      "示例企业-2025",
      "示例企业有限公司"
    ]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_002"
}
```

---

### 2.2 检查子域名是否可用

检查子域名是否已被占用。

**接口描述**: 验证子域名的唯一性

**请求方式**: `GET /tenants/check-subdomain`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| subdomain | string | ✅ | 子域名(3-63个字符,仅小写字母、数字、连字符) |

**请求示例**:
```http
GET /api/v1/tenants/check-subdomain?subdomain=example-company
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "available": true,
    "suggestions": [],
    "full_domain": "example-company.saas.coze.com"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_003"
}
```

**子域名已占用时** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "available": false,
    "suggestions": [
      "example-company-2025",
      "example-company-tech",
      "example-company-corp"
    ]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_004"
}
```

---

### 2.3 发送注册验证码

发送邮箱验证码,用于注册验证。

**接口描述**: 向指定邮箱发送6位数字验证码

**请求方式**: `POST /tenants/send-verification-code`

**权限要求**: 无需认证

**请求体**:
```json
{
  "email": "contact@example.com"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | ✅ | 邮箱地址 |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "验证码已发送",
  "data": {
    "expires_in": 300,
    "message": "验证码有效期为5分钟"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_005"
}
```

**错误响应** (429 Too Many Requests):
```json
{
  "code": 42901,
  "message": "验证码发送过于频繁,请1分钟后重试",
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "retry_after": 60
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_006"
}
```

---

### 2.4 注册租户

创建新的租户账户。

**接口描述**: 自助注册租户,5分钟完成入驻

**请求方式**: `POST /tenants/register`

**权限要求**: 无需认证

**请求体**:
```json
{
  "company_name": "示例企业",
  "contact_name": "张三",
  "contact_email": "zhangsan@example.com",
  "contact_phone": "13800138000",
  "verification_code": "123456",
  "industry": "互联网",
  "employee_count": 50,
  "deployment_type": "public",
  "subdomain": "example-company",
  "plan_id": "plan-free",
  "timezone": "Asia/Shanghai",
  "language": "zh-CN"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| company_name | string | ✅ | 企业名称 |
| contact_name | string | ✅ | 联系人姓名 |
| contact_email | string | ✅ | 联系人邮箱 |
| contact_phone | string | ✅ | 联系人电话 |
| verification_code | string | ✅ | 邮箱验证码 |
| industry | string | ✅ | 行业(互联网/金融/教育/医疗/制造业/零售/其他) |
| employee_count | integer | ✅ | 员工数量(1-100000) |
| deployment_type | string | ✅ | 部署方式(public/private) |
| subdomain | string | ✅ | 子域名(3-63个字符) |
| plan_id | string | ✅ | 套餐ID(plan-free/plan-professional/plan-enterprise) |
| timezone | string | ❌ | 时区(默认Asia/Shanghai) |
| language | string | ❌ | 语言(默认zh-CN) |

**成功响应** (201 Created):
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "company_name": "示例企业",
    "subdomain": "example-company",
    "status": "trial",
    "trial_end_date": "2025-02-12T23:59:59Z",
    "admin_user": {
      "user_id": "usr_1234567890",
      "username": "zhangsan@example.com",
      "initial_password": "Abc12345",
      "name": "张三"
    },
    "login_url": "https://example-company.saas.coze.com/login",
    "plan": {
      "plan_id": "plan-free",
      "plan_name": "免费版",
      "max_users": 5,
      "max_bots": 1,
      "max_storage_gb": 1,
      "max_api_calls_per_day": 100
    },
    "created_at": "2025-12-30T10:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_007"
}
```

**错误响应**:

400 Bad Request (企业名称已存在):
```json
{
  "code": 40001,
  "message": "企业名称已被注册",
  "error": {
    "code": "COMPANY_NAME_EXISTS",
    "field": "company_name",
    "suggestions": ["示例企业(科技)", "示例企业-2025"]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_008"
}
```

400 Bad Request (验证码错误):
```json
{
  "code": 40002,
  "message": "验证码错误或已过期",
  "error": {
    "code": "INVALID_VERIFICATION_CODE"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_009"
}
```

---

## 3. 租户信息管理API

### 3.1 获取当前租户信息

获取当前登录用户的租户详细信息。

**接口描述**: 获取租户的完整信息,包括套餐、配额、使用情况

**请求方式**: `GET /tenants/current`

**权限要求**: 需要认证

**请求示例**:
```http
GET /api/v1/tenants/current
Authorization: Bearer {token}
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "company_name": "示例企业",
    "subdomain": "example-company",
    "custom_domain": null,
    "logo": "https://cdn.saas.coze.com/logo/tenant-001.png",
    "contact": {
      "name": "张三",
      "email": "zhangsan@example.com",
      "phone": "13800138000"
    },
    "company_info": {
      "industry": "互联网",
      "employee_count": 50,
      "deployment_type": "public"
    },
    "status": "trial",
    "trial_period": {
      "start_date": "2025-12-30T00:00:00Z",
      "end_date": "2025-02-12T23:59:59Z",
      "days_remaining": 44
    },
    "plan": {
      "plan_id": "plan-free",
      "plan_name": "免费版",
      "price": 0,
      "duration": 12,
      "features": {
        "max_users": 5,
        "max_bots": 1,
        "max_knowledge_bases": 1,
        "max_documents": 10,
        "max_storage_gb": 1,
        "max_api_calls_per_day": 100,
        "custom_domain": false,
        "ssl_certificate": false,
        "api_access": false,
        "sso": false
      }
    },
    "usage": {
      "users": 3,
      "bots": 1,
      "knowledge_bases": 1,
      "documents": 5,
      "storage_gb": 0.5,
      "api_calls_today": 45
    },
    "settings": {
      "timezone": "Asia/Shanghai",
      "language": "zh-CN"
    },
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-30T10:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_010"
}
```

---

### 3.2 更新租户信息

更新租户的基本信息。

**接口描述**: 更新联系人信息、Logo、时区等

**请求方式**: `PUT /tenants/current`

**权限要求**: 需要认证,租户管理员

**请求体**:
```json
{
  "contact_name": "李四",
  "contact_phone": "13900139000",
  "logo": "https://cdn.saas.coze.com/logo/new-logo.png",
  "custom_domain": "www.example.com",
  "timezone": "Asia/Shanghai",
  "language": "zh-CN"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| contact_name | string | ❌ | 联系人姓名 |
| contact_phone | string | ❌ | 联系人电话 |
| logo | string | ❌ | 企业Logo URL |
| custom_domain | string | ❌ | 自定义域名(需套餐支持) |
| timezone | string | ❌ | 时区 |
| language | string | ❌ | 语言(zh-CN/en-US) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "updated_fields": [
      "contact_name",
      "contact_phone",
      "logo"
    ],
    "updated_at": "2025-12-30T14:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_011"
}
```

**错误响应** (403 Forbidden):
```json
{
  "code": 40301,
  "message": "当前套餐不支持自定义域名",
  "error": {
    "code": "FEATURE_NOT_SUPPORTED",
    "feature": "custom_domain",
    "current_plan": "plan-free",
    "required_plan": "plan-professional",
    "upgrade_url": "https://example-company.saas.coze.com/settings/plans"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_012"
}
```

---

## 4. 租户识别API

### 4.1 通过子域名识别租户(内部)

**接口描述**: 系统内部通过子域名识别租户,由中间件自动处理

**识别方式**:
1. 子域名: `{tenant_id}.saas.coze.com`
2. 路径: `saas.coze.com/{tenant_id}/`
3. Header: `X-Tenant-ID: {tenant_id}` (API调用)

**成功时**: 中间件将 `tenant_id` 和 `tenant` 信息注入到请求上下文中

**失败响应** (400 Bad Request):
```json
{
  "code": 40001,
  "message": "无法识别租户ID",
  "error": {
    "code": "TENANT_NOT_IDENTIFIED",
    "hint": "请确保URL包含租户信息: {tenant_id}.saas.coze.com 或 saas.coze.com/{tenant_id}/"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_013"
}
```

**租户不存在** (404 Not Found):
```json
{
  "code": 40401,
  "message": "租户不存在",
  "error": {
    "code": "TENANT_NOT_FOUND",
    "tenant_id": "example-company"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_014"
}
```

**租户已停用** (403 Forbidden):
```json
{
  "code": 40301,
  "message": "租户已停用",
  "error": {
    "code": "TENANT_SUSPENDED",
    "tenant_id": "example-company",
    "reason": "请联系客服",
    "support_email": "support@saas.coze.com"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_015"
}
```

---

## 5. 租户配额管理API

### 5.1 获取租户配额

获取当前租户的所有配额信息。

**接口描述**: 查询租户的各类配额及使用情况

**请求方式**: `GET /tenants/current/quotas`

**权限要求**: 需要认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| quota_type | string | 否 | 配额类型筛选(users/bots/storage/api_calls/concurrent) |

**请求示例**:
```http
GET /api/v1/tenants/current/quotas
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": {
      "quota_type": "users",
      "max_value": 5,
      "current_value": 3,
      "remaining": 2,
      "unit": "个",
      "usage_percentage": 60,
      "status": "normal"
    },
    "bots": {
      "quota_type": "bots",
      "max_value": 1,
      "current_value": 1,
      "remaining": 0,
      "unit": "个",
      "usage_percentage": 100,
      "status": "warning",
      "upgrade_url": "https://example-company.saas.coze.com/settings/plans"
    },
    "storage": {
      "quota_type": "storage",
      "max_value": 1,
      "current_value": 0.5,
      "remaining": 0.5,
      "unit": "GB",
      "usage_percentage": 50,
      "status": "normal"
    },
    "api_calls": {
      "quota_type": "api_calls",
      "max_value": 100,
      "current_value": 45,
      "remaining": 55,
      "unit": "次/天",
      "usage_percentage": 45,
      "status": "normal",
      "reset_at": "2025-12-31T00:00:00Z"
    },
    "concurrent": {
      "quota_type": "concurrent",
      "max_value": 2,
      "current_value": 1,
      "remaining": 1,
      "unit": "个",
      "usage_percentage": 50,
      "status": "normal"
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_016"
}
```

---

### 5.2 检查配额

检查指定配额是否有足够余额。

**接口描述**: 在执行操作前检查配额是否充足

**请求方式**: `POST /tenants/current/quotas/check`

**权限要求**: 需要认证

**请求体**:
```json
{
  "quota_type": "users",
  "amount": 1
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| quota_type | string | ✅ | 配额类型 |
| amount | integer | ✅ | 需要的数量 |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "配额充足",
  "data": {
    "available": true,
    "quota_type": "users",
    "current": 3,
    "max": 5,
    "remaining": 2,
    "requested": 1
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_017"
}
```

**配额不足** (400 Bad Request):
```json
{
  "code": 40003,
  "message": "配额不足: 用户数已达上限(5/5),请升级套餐",
  "error": {
    "code": "QUOTA_EXCEEDED",
    "quota_type": "users",
    "current": 5,
    "max": 5,
    "requested": 1,
    "upgrade_url": "https://example-company.saas.coze.com/settings/plans"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_018"
}
```

---

### 5.3 获取配额使用历史

获取配额使用历史记录。

**接口描述**: 查询配额历史使用情况

**请求方式**: `GET /tenants/current/quotas/history`

**权限要求**: 需要认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| quota_type | string | 否 | 配额类型筛选 |
| start_date | string | 否 | 开始日期(ISO8601) |
| end_date | string | 否 | 结束日期(ISO8601) |
| page | integer | 否 | 页码(默认1) |
| page_size | integer | 否 | 每页大小(默认20) |

**请求示例**:
```http
GET /api/v1/tenants/current/quotas/history?quota_type=api_calls&start_date=2025-12-01&end_date=2025-12-30
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "date": "2025-12-29",
        "quota_type": "api_calls",
        "usage_value": 85,
        "max_value": 100,
        "usage_percentage": 85
      },
      {
        "date": "2025-12-28",
        "quota_type": "api_calls",
        "usage_value": 92,
        "max_value": 100,
        "usage_percentage": 92
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 30,
      "total_pages": 2,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_019"
}
```

---

## 6. 租户套餐管理API

### 6.1 获取所有套餐列表

获取所有可用的套餐信息。

**接口描述**: 查询所有套餐及价格

**请求方式**: `GET /tenant-plans`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| include_inactive | boolean | 否 | 是否包含已下架套餐(默认false) |

**请求示例**:
```http
GET /api/v1/tenant-plans
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "plan-free",
      "name": "免费版",
      "description": "体验基础功能",
      "price": 0,
      "duration": 12,
      "is_popular": false,
      "features": {
        "max_users": 5,
        "max_bots": 1,
        "max_knowledge_bases": 1,
        "max_documents": 10,
        "max_storage_gb": 1,
        "max_api_calls_per_day": 100,
        "max_api_calls_per_month": 3000,
        "max_concurrent_users": 2,
        "max_concurrent_requests": 10,
        "custom_domain": false,
        "ssl_certificate": false,
        "api_access": false,
        "sso": false,
        "white_label": false,
        "priority_support": false,
        "sla": "无保障"
      }
    },
    {
      "id": "plan-professional",
      "name": "专业版",
      "description": "适合中小型团队",
      "price": 299,
      "duration": 1,
      "is_popular": true,
      "features": {
        "max_users": 50,
        "max_bots": 10,
        "max_knowledge_bases": 10,
        "max_documents": 500,
        "max_storage_gb": 50,
        "max_api_calls_per_day": 5000,
        "max_api_calls_per_month": 150000,
        "max_concurrent_users": 20,
        "max_concurrent_requests": 100,
        "custom_domain": true,
        "ssl_certificate": true,
        "api_access": true,
        "sso": false,
        "white_label": false,
        "priority_support": true,
        "sla": "99.5%可用性"
      }
    },
    {
      "id": "plan-enterprise",
      "name": "企业版",
      "description": "适合大型企业定制需求",
      "price": 999,
      "duration": 1,
      "is_popular": false,
      "features": {
        "max_users": 500,
        "max_bots": 100,
        "max_knowledge_bases": 100,
        "max_documents": 10000,
        "max_storage_gb": 500,
        "max_api_calls_per_day": 50000,
        "max_api_calls_per_month": 1500000,
        "max_concurrent_users": 200,
        "max_concurrent_requests": 1000,
        "custom_domain": true,
        "ssl_certificate": true,
        "api_access": true,
        "sso": true,
        "white_label": true,
        "priority_support": true,
        "sla": "99.9%可用性"
      }
    }
  ],
  "timestamp": 1735574400000,
  "trace_id": "trace_020"
}
```

---

### 6.2 获取套餐对比

对比不同套餐的功能差异。

**接口描述**: 直观展示套餐间的差异

**请求方式**: `GET /tenant-plans/compare`

**权限要求**: 无需认证

**请求示例**:
```http
GET /api/v1/tenant-plans/compare
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "comparison": [
      {
        "feature": "用户数",
        "plan-free": "5个",
        "plan-professional": "50个",
        "plan-enterprise": "500个"
      },
      {
        "feature": "Bot数量",
        "plan-free": "1个",
        "plan-professional": "10个",
        "plan-enterprise": "100个"
      },
      {
        "feature": "存储空间",
        "plan-free": "1GB",
        "plan-professional": "50GB",
        "plan-enterprise": "500GB"
      },
      {
        "feature": "API调用/天",
        "plan-free": "100次",
        "plan-professional": "5000次",
        "plan-enterprise": "50000次"
      },
      {
        "feature": "自定义域名",
        "plan-free": "❌",
        "plan-professional": "✓",
        "plan-enterprise": "✓"
      },
      {
        "feature": "SSO单点登录",
        "plan-free": "❌",
        "plan-professional": "❌",
        "plan-enterprise": "✓"
      },
      {
        "feature": "月费",
        "plan-free": "¥0",
        "plan-professional": "¥299",
        "plan-enterprise": "¥999"
      }
    ],
    "recommendation": {
      "plan_id": "plan-professional",
      "reason": "专业版提供了最佳性价比,适合大多数中小型团队"
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_021"
}
```

---

### 6.3 购买/订阅套餐

购买或升级套餐。

**接口描述**: 创建订阅订单并返回支付信息

**请求方式**: `POST /tenants/current/plans/subscribe`

**权限要求**: 需要认证,租户管理员

**请求体**:
```json
{
  "plan_id": "plan-professional",
  "duration": 12,
  "payment_method": "alipay",
  "coupon_code": "NEWYEAR2025"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| plan_id | string | ✅ | 目标套餐ID |
| duration | integer | ✅ | 订阅时长(月) |
| payment_method | string | ✅ | 支付方式(alipay/wechat_pay/stripe) |
| coupon_code | string | ❌ | 优惠券代码 |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "订单创建成功",
  "data": {
    "order_id": "order-1234567890",
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "old_plan": {
      "plan_id": "plan-free",
      "plan_name": "免费版"
    },
    "new_plan": {
      "plan_id": "plan-professional",
      "plan_name": "专业版"
    },
    "billing": {
      "duration": 12,
      "original_amount": 3588,
      "discount": 358.8,
      "final_amount": 3229.2,
      "currency": "CNY"
    },
    "payment": {
      "method": "alipay",
      "payment_url": "https://openapi.alipay.com/gateway.do?...",
      "qr_code_url": "HTTPS://QR.ALIPAY.COM/...",
      "expires_at": "2025-12-30T11:00:00Z"
    },
    "created_at": "2025-12-30T10:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_022"
}
```

**错误响应** (400 Bad Request):
```json
{
  "code": 40004,
  "message": "不允许降级套餐,请先联系客服",
  "error": {
    "code": "PLAN_DOWNGRADE_NOT_ALLOWED",
    "current_plan": "plan-professional",
    "target_plan": "plan-free",
    "support_email": "sales@saas.coze.com"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_023"
}
```

---

### 6.4 获取当前套餐详情

获取当前订阅的套餐详细信息。

**接口描述**: 查询当前套餐、到期时间、续费信息等

**请求方式**: `GET /tenants/current/subscription`

**权限要求**: 需要认证

**请求示例**:
```http
GET /api/v1/tenants/current/subscription
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "plan": {
      "plan_id": "plan-professional",
      "plan_name": "专业版",
      "price": 299,
      "duration": 1
    },
    "subscription": {
      "status": "active",
      "start_date": "2025-01-01T00:00:00Z",
      "end_date": "2026-01-01T00:00:00Z",
      "auto_renew": true,
      "days_remaining": 2
    },
    "billing": {
      "next_billing_date": "2026-01-01T00:00:00Z",
      "next_billing_amount": 299,
      "currency": "CNY",
      "payment_method": "alipay"
    },
    "upgrade_options": [
      {
        "plan_id": "plan-enterprise",
        "plan_name": "企业版",
        "price": 999,
        "prorated_amount": 700
      }
    ],
    "cancel_url": "https://example-company.saas.coze.com/settings/subscription/cancel"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_024"
}
```

---

### 6.5 取消订阅

取消当前订阅,停止自动续费。

**接口描述**: 取消订阅,当前周期结束后不再续费

**请求方式**: `POST /tenants/current/subscription/cancel`

**权限要求**: 需要认证,租户管理员

**请求体**:
```json
{
  "reason": "成本优化",
  "feedback": "功能超出需求",
  "cancel_immediately": false
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reason | string | ❌ | 取消原因 |
| feedback | string | ❌ | 详细反馈 |
| cancel_immediately | boolean | ❌ | 是否立即取消(默认false,周期结束时取消) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "订阅已取消",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "subscription": {
      "status": "active",
      "auto_renew": false,
      "end_date": "2026-01-01T00:00:00Z",
      "will_expire_at_period_end": true
    },
    "message": "订阅将在当前计费周期结束时取消,您仍可使用服务至 2026-01-01",
    "data_retention_days": 30
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_025"
}
```

---

## 7. 租户生命周期管理API

### 7.1 激活租户(管理员)

激活已停用的租户。

**接口描述**: 系统管理员激活停用的租户

**请求方式**: `POST /admin/tenants/{tenant_id}/activate`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求体**:
```json
{
  "reason": "已完成欠费补缴"
}
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "租户已激活",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "status": "active",
    "activated_at": "2025-12-30T14:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_026"
}
```

---

### 7.2 停用租户(管理员)

停用违规或欠费的租户。

**接口描述**: 系统管理员停用租户

**请求方式**: `POST /admin/tenants/{tenant_id}/suspend`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求体**:
```json
{
  "reason": "违反服务条款",
  "details": "发送垃圾邮件",
  "notify": true
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reason | string | ✅ | 停用原因 |
| details | string | ❌ | 详细说明 |
| notify | boolean | ❌ | 是否发送通知邮件(默认true) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "租户已停用",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "status": "suspended",
    "suspended_at": "2025-12-30T14:00:00Z",
    "reason": "违反服务条款"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_027"
}
```

---

### 7.3 删除租户(管理员)

删除租户,支持软删除和硬删除。

**接口描述**: 系统管理员删除租户

**请求方式**: `DELETE /admin/tenants/{tenant_id}`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| hard_delete | boolean | 否 | 是否硬删除(永久删除数据,默认false) |

**请求示例**:
```http
DELETE /api/v1/admin/tenants/tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890?hard_delete=false
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "租户已删除",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "deleted": true,
    "hard_delete": false,
    "data_retention_days": 30,
    "deleted_at": "2025-12-30T14:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_028"
}
```

---

## 8. 租户管理后台API

### 8.1 获取租户列表(管理员)

获取所有租户的列表。

**接口描述**: 系统管理员查看所有租户

**请求方式**: `GET /admin/tenants`

**权限要求**: 系统管理员

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 状态筛选(trial/active/suspended/deleted) |
| plan_id | string | 否 | 套餐筛选 |
| keyword | string | 否 | 搜索关键词(企业名称) |
| page | integer | 否 | 页码(默认1) |
| page_size | integer | 否 | 每页大小(默认20) |
| sort_by | string | 否 | 排序字段(created_at/updated_at/company_name) |
| sort_order | string | 否 | 排序方向(asc/desc,默认desc) |

**请求示例**:
```http
GET /api/v1/admin/tenants?status=trial&page=1&page_size=20
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "items": [
      {
        "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "company_name": "示例企业",
        "subdomain": "example-company",
        "contact_email": "zhangsan@example.com",
        "industry": "互联网",
        "status": "trial",
        "plan_name": "免费版",
        "trial_start_date": "2025-12-30T00:00:00Z",
        "trial_end_date": "2025-02-12T23:59:59Z",
        "trial_days_remaining": 44,
        "usage": {
          "users": 3,
          "bots": 1,
          "storage_gb": 0.5
        },
        "created_at": "2025-12-30T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_029"
}
```

---

### 8.2 获取租户统计(管理员)

获取租户的统计数据。

**接口描述**: 系统管理员查看租户统计

**请求方式**: `GET /admin/tenants/stats`

**权限要求**: 系统管理员

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | string | 否 | 统计周期(today/week/month/year/all,默认month) |

**请求示例**:
```http
GET /api/v1/admin/tenants/stats?period=month
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "month",
    "summary": {
      "total_tenants": 500,
      "new_tenants": 50,
      "trial_tenants": 120,
      "active_tenants": 350,
      "suspended_tenants": 30,
      "deleted_tenants": 0
    },
    "by_status": [
      {
        "status": "trial",
        "count": 120,
        "percentage": 24
      },
      {
        "status": "active",
        "count": 350,
        "percentage": 70
      }
    ],
    "by_plan": [
      {
        "plan_name": "免费版",
        "count": 150,
        "percentage": 30
      },
      {
        "plan_name": "专业版",
        "count": 280,
        "percentage": 56
      },
      {
        "plan_name": "企业版",
        "count": 70,
        "percentage": 14
      }
    ],
    "conversion": {
      "trial_to_paid_rate": 58.3,
      "avg_trial_days": 15
    },
    "revenue": {
      "mrr": 295000,
      "arr": 3540000,
      "currency": "CNY"
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_030"
}
```

---

### 8.3 获取试用即将到期的租户(管理员)

获取试用即将到期的租户列表。

**接口描述**: 查询7天内试用到期的租户,用于跟进转化

**请求方式**: `GET /admin/tenants/trial-expiring`

**权限要求**: 系统管理员

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | integer | 否 | 天数范围(默认7天) |
| page | integer | 否 | 页码(默认1) |
| page_size | integer | 否 | 每页大小(默认20) |

**请求示例**:
```http
GET /api/v1/admin/tenants/trial-expiring?days=7
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "tenant_id": "tenant-b2c3d4e5-f6g7-8901-bcde-f23456789012",
        "company_name": "测试企业A",
        "subdomain": "test-company-a",
        "contact_email": "contact@test-a.com",
        "contact_phone": "13900139000",
        "trial_end_date": "2025-01-05T23:59:59Z",
        "days_remaining": 6,
        "usage": {
          "users": 5,
          "bots": 1,
          "api_calls_this_week": 350
        },
        "engagement_score": 85,
        "conversion_probability": "high"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 25,
      "total_pages": 2,
      "has_next": true,
      "has_prev": false
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_031"
}
```

---

## 9. 错误代码

| 错误代码 | HTTP状态码 | 说明 |
|---------|-----------|------|
| 40001 | 400 | 企业名称已被注册 |
| 40002 | 400 | 验证码错误或已过期 |
| 40003 | 400 | 配额不足 |
| 40004 | 400 | 不允许降级套餐 |
| 40101 | 401 | 未认证 |
| 40301 | 403 | 租户已停用 |
| 40302 | 403 | 功能不支持当前套餐 |
| 40401 | 404 | 租户不存在 |
| 40901 | 409 | 子域名已被占用 |
| 41001 | 410 | 租户已注销 |
| 42901 | 429 | 验证码发送过于频繁 |
| 42902 | 429 | API调用次数已达上限 |
| 50001 | 500 | 服务器内部错误 |

---

## 10. 代码示例

### 10.1 Go 后端实现示例

```go
// api/handler/tenant_handler.go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

type TenantHandler struct {
    tenantService *TenantService
    quotaService  *TenantQuotaService
    planService   *TenantPlanService
}

// CheckCompanyName 检查企业名称是否可用
//
//	@Summary		检查企业名称是否可用
//	@Description	验证企业名称的唯一性
//	@Tags			tenants
//	@Accept			json
//	@Produce		json
//	@Param			company_name	query		string	true	"企业名称"
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Router			/tenants/check-company-name [get]
func (h *TenantHandler) CheckCompanyName(ctx context.Context, c *app.RequestContext) {
    companyName := c.Query("company_name")

    if companyName == "" {
        response.Error(c, 400, 40001, "INVALID_PARAMETER", "企业名称不能为空", nil)
        return
    }

    available, suggestions, err := h.tenantService.CheckCompanyNameAvailability(ctx, companyName)
    if err != nil {
        response.Error(c, 500, 50001, "INTERNAL_ERROR", "服务器内部错误", nil)
        return
    }

    response.Success(c, &struct {
        Available   bool     `json:"available"`
        Suggestions []string `json:"suggestions"`
    }{
        Available:   available,
        Suggestions: suggestions,
    })
}

// RegisterTenant 注册租户
//
//	@Summary		注册租户
//	@Description	创建新的租户账户
//	@Tags			tenants
//	@Accept			json
//	@Produce		json
//	@Param			request	body		request.RegisterTenantRequest	true	"注册请求"
//	@Success		201		{object}	response.SuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Router			/tenants/register [post]
func (h *TenantHandler) RegisterTenant(ctx context.Context, c *app.RequestContext) {
    var req RegisterTenantRequest
    if err := c.BindAndValidate(&req); err != nil {
        response.Error(c, 400, 40001, "INVALID_PARAMETER", "请求参数错误", err)
        return
    }

    tenant, err := h.tenantService.Register(ctx, &req)
    if err != nil {
        // 处理业务错误
        if err.Error() == "企业名称已被注册" {
            response.Error(c, 400, 40001, "COMPANY_NAME_EXISTS", "企业名称已被注册", nil)
        } else if err.Error() == "子域名已被占用" {
            response.Error(c, 400, 40901, "SUBDOMAIN_EXISTS", "子域名已被占用", nil)
        } else {
            response.Error(c, 500, 50001, "INTERNAL_ERROR", "服务器内部错误", nil)
        }
        return
    }

    response.Success(c, tenant)
}

// GetCurrentTenant 获取当前租户信息
//
//	@Summary		获取当前租户信息
//	@Description	获取当前登录用户的租户详细信息
//	@Tags			tenants
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		401		{object}	response.ErrorResponse
//	@Router			/tenants/current [get]
func (h *TenantHandler) GetCurrentTenant(ctx context.Context, c *app.RequestContext) {
    tenantID := ctx.Value("tenant_id").(string)

    tenant, err := h.tenantService.GetByID(ctx, tenantID)
    if err != nil {
        response.Error(c, 404, 40401, "TENANT_NOT_FOUND", "租户不存在", nil)
        return
    }

    // 获取配额信息
    quotas, err := h.quotaService.GetTenantQuotas(ctx, tenantID)
    if err != nil {
        response.Error(c, 500, 50001, "INTERNAL_ERROR", "服务器内部错误", nil)
        return
    }

    // 获取套餐信息
    plan, err := h.planService.GetByID(ctx, tenant.PlanID)
    if err != nil {
        response.Error(c, 500, 50001, "INTERNAL_ERROR", "服务器内部错误", nil)
        return
    }

    response.Success(c, &struct {
        Tenant *Tenant       `json:"tenant"`
        Quotas []TenantQuota  `json:"quotas"`
        Plan   *TenantPlan    `json:"plan"`
    }{
        Tenant: tenant,
        Quotas: quotas,
        Plan:   plan,
    })
}
```

### 10.2 前端 TypeScript 类型定义

```typescript
// types/tenant.ts

// 租户相关类型
interface Tenant {
  tenant_id: string;
  company_name: string;
  subdomain: string;
  custom_domain?: string;
  logo?: string;
  contact: {
    name: string;
    email: string;
    phone: string;
  };
  company_info: {
    industry: string;
    employee_count: number;
    deployment_type: 'public' | 'private';
  };
  status: 'trial' | 'active' | 'suspended' | 'deleted';
  trial_period?: {
    start_date: string;
    end_date: string;
    days_remaining: number;
  };
  plan: TenantPlan;
  usage: TenantUsage;
  settings: {
    timezone: string;
    language: string;
  };
  created_at: string;
  updated_at: string;
}

interface TenantPlan {
  plan_id: string;
  plan_name: string;
  price: number;
  duration: number;
  features: {
    max_users: number;
    max_bots: number;
    max_knowledge_bases: number;
    max_documents: number;
    max_storage_gb: number;
    max_api_calls_per_day: number;
    max_api_calls_per_month: number;
    max_concurrent_users: number;
    max_concurrent_requests: number;
    custom_domain: boolean;
    ssl_certificate: boolean;
    api_access: boolean;
    sso: boolean;
    white_label: boolean;
    priority_support: boolean;
    sla: string;
  };
}

interface TenantQuota {
  quota_type: string;
  max_value: number;
  current_value: number;
  remaining: number;
  unit: string;
  usage_percentage: number;
  status: 'normal' | 'warning' | 'exceeded';
}

interface TenantUsage {
  users: number;
  bots: number;
  knowledge_bases: number;
  documents: number;
  storage_gb: number;
  api_calls_today: number;
}

// API客户端方法
class TenantApiClient {
  private baseURL = '/api/v1';

  // 检查企业名称
  async checkCompanyName(companyName: string): Promise<{
    available: boolean;
    suggestions: string[];
  }> {
    const response = await apiClient.get(
      `${this.baseURL}/tenants/check-company-name`,
      { params: { company_name: companyName } }
    );
    return response.data.data;
  }

  // 检查子域名
  async checkSubdomain(subdomain: string): Promise<{
    available: boolean;
    suggestions: string[];
    full_domain?: string;
  }> {
    const response = await apiClient.get(
      `${this.baseURL}/tenants/check-subdomain`,
      { params: { subdomain } }
    );
    return response.data.data;
  }

  // 发送验证码
  async sendVerificationCode(email: string): Promise<void> {
    await apiClient.post(`${this.baseURL}/tenants/send-verification-code`, {
      email,
    });
  }

  // 注册租户
  async register(data: RegisterTenantRequest): Promise<Tenant> {
    const response = await apiClient.post(`${this.baseURL}/tenants/register`, data);
    return response.data.data;
  }

  // 获取当前租户信息
  async getCurrentTenant(): Promise<Tenant> {
    const response = await apiClient.get(`${this.baseURL}/tenants/current`);
    return response.data.data;
  }

  // 获取配额
  async getQuotas(quotaType?: string): Promise<Record<string, TenantQuota>> {
    const response = await apiClient.get(`${this.baseURL}/tenants/current/quotas`, {
      params: quotaType ? { quota_type: quotaType } : undefined,
    });
    return response.data.data;
  }

  // 检查配额
  async checkQuota(quotaType: string, amount: number): Promise<void> {
    await apiClient.post(`${this.baseURL}/tenants/current/quotas/check`, {
      quota_type: quotaType,
      amount,
    });
  }

  // 获取套餐列表
  async getPlans(): Promise<TenantPlan[]> {
    const response = await apiClient.get(`${this.baseURL}/tenant-plans`);
    return response.data.data;
  }

  // 订阅套餐
  async subscribePlan(data: SubscribePlanRequest): Promise<SubscribeResult> {
    const response = await apiClient.post(
      `${this.baseURL}/tenants/current/plans/subscribe`,
      data
    );
    return response.data.data;
  }

  // 取消订阅
  async cancelSubscription(reason?: string, cancelImmediately = false): Promise<void> {
    await apiClient.post(`${this.baseURL}/tenants/current/subscription/cancel`, {
      reason,
      cancel_immediately: cancelImmediately,
    });
  }
}

export const tenantApi = new TenantApiClient();
```

---

**文档结束**

---

## 📊 文档统计

- **接口总数**: 20+
- **核心功能模块**: 8个
- **代码示例**: Go + TypeScript
- **错误代码**: 13个

**主要功能模块**:
1. 租户注册 (4个接口)
2. 租户信息管理 (2个接口)
3. 租户识别 (中间件自动处理)
4. 租户配额管理 (3个接口)
5. 租户套餐管理 (5个接口)
6. 租户生命周期管理 (3个接口)
7. 租户管理后台 (3个接口)
