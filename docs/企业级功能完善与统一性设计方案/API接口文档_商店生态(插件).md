# API接口文档：商店生态(插件)模块

**模块名称**: 商店生态 (PluginStore)
**设计文档**: 12-百应开发平台_商店生态.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 插件商店API](#2-插件商店api)
- [3. 插件管理API](#3-插件管理api)
- [4. 插件审核API](#4-插件审核api)
- [5. 评价系统API](#5-评价系统api)
- [6. 收益管理API](#6-收益管理api)
- [7. 数据模型](#7-数据模型)
- [8. 错误码定义](#8-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

商店生态是ZKER的插件/技能市场模块，通过**复用zker Plugin Store** + 商业化功能扩展，实现插件上架、审核、计费、评价等完整生态：

- ✅ **插件浏览** - 浏览插件商店、搜索插件、分类筛选
- ✅ **插件安装** - 免费插件安装、付费插件购买
- ✅ **插件管理** - 上传插件、更新版本、查看数据
- ✅ **审核机制** - 插件上架审核流程
- ✅ **评价系统** - 插件评价、评分、推荐算法
- ✅ **收益管理** - 收入明细、提现申请、收益报表

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 实现: 100%复用zker Plugin Store + 配置扩展

**前端技术栈**:
- React 18 + TypeScript 5.6
- Semi Design UI组件库

**实现策略**: ✅ 100%复用 + 商业化扩展

### 1.3 复用关系

| 功能 | 实现方式 | 说明 |
|------|---------|------|
| 插件定义 | 🔁 复用plugins表 | 直接复用zker |
| 插件版本 | 🔁 复用plugin_versions表 | 直接复用zker |
| 插件评价 | 🔁 复用plugin_reviews表 | 直接复用zker |
| 下载记录 | 🔁 复用plugin_downloads表 | 直接复用zker |
| 商业化字段 | 📊 扩展plugins表 | pricing_type, price等 |
| 审核流程 | 📊 新增plugin_reviews_audit表 | 审核管理 |
| 使用统计 | 📊 新增plugin_usage_stats表 | 统计分析 |

---

## 2. 插件商店API

### 2.1 获取插件列表

**接口地址**: `GET /api/v1/store/plugins`

**查询参数**:
- category: 分类筛选 (可选)
- keyword: 搜索关键词 (可选)
- pricing_type: 定价类型 (可选，free/paid/freemium)
- sort: 排序方式 (可选，recommended/popular/newest/rating)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 150,
    "items": [
      {
        "id": "plugin-google-search-001",
        "name": "Google搜索",
        "description": "使用Google搜索网页内容",
        "icon": "https://example.com/icons/google-search.png",
        "category": "搜索工具",
        "pricing_type": "free",
        "price": null,
        "rating": 4.5,
        "review_count": 120,
        "install_count": 15200,
        "author": "ZKER Team",
        "version": "1.0.0",
        "status": "approved",
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      },
      {
        "id": "plugin-email-sender-001",
        "name": "邮件发送",
        "description": "发送邮件、邮件模板",
        "icon": "https://example.com/icons/email-sender.png",
        "category": "通信工具",
        "pricing_type": "paid",
        "price": 99.00,
        "rating": 4.8,
        "review_count": 85,
        "install_count": 8900,
        "author": "Developer Inc.",
        "version": "2.1.0",
        "status": "approved",
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 2.2 获取插件详情

**接口地址**: `GET /api/v1/store/plugins/{plugin_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "plugin-google-search-001",
    "name": "Google搜索",
    "description": "使用Google搜索网页内容，支持搜索结果摘要",
    "long_description": "<p>强大的Google搜索插件，支持：</p><ul><li>网页搜索</li><li>搜索结果摘要</li><li>自定义搜索参数</li></ul>",
    "icon": "https://example.com/icons/google-search.png",
    "screenshots": [
      "https://example.com/screenshots/s1.png",
      "https://example.com/screenshots/s2.png"
    ],
    "category": "搜索工具",
    "pricing_type": "free",
    "price": null,
    "rating": 4.5,
    "review_count": 120,
    "install_count": 15200,
    "active_users": 3200,
    "author": "ZKER Team",
    "author_id": 1001,
    "version": "1.0.0",
    "status": "approved",
    "config_schema": [
      {
        "name": "api_key",
        "type": "string",
        "required": true,
        "label": "API Key",
        "description": "Google Custom Search API Key"
      },
      {
        "name": "cx",
        "type": "string",
        "required": true,
        "label": "搜索引擎ID",
        "description": "Custom Search Engine ID"
      }
    ],
    "permissions": [
      "http:api.google.com",
      "storage:kv"
    ],
    "usage_example": "1. 安装插件\n2. 配置API Key和CX\n3. 在Bot中使用Google搜索工具",
    "changelog": [
      {
        "version": "1.0.0",
        "date": "2025-01-01",
        "changes": "初始版本发布"
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

### 2.3 获取推荐插件

**接口地址**: `GET /api/v1/store/plugins/recommended`

**查询参数**:
- limit: 返回数量 (默认10，最大50)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "plugin-weather-001",
        "name": "天气查询",
        "description": "实时天气、天气预报",
        "icon": "https://example.com/icons/weather.png",
        "category": "生活工具",
        "pricing_type": "free",
        "rating": 4.7,
        "review_count": 210,
        "install_count": 25600,
        "score": 92.5,
        "reason": "推荐理由：高评分、高安装量、活跃用户多"
      }
    ]
  }
}
```

### 2.4 获取热门插件

**接口地址**: `GET /api/v1/store/plugins/popular`

**查询参数**:
- limit: 返回数量 (默认10，最大50)
- period: 统计周期 (可选，week/month/all，默认week)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "week",
    "items": [
      {
        "id": "plugin-weather-001",
        "name": "天气查询",
        "install_count": 1250,
        "trend": "+15.2%"
      }
    ]
  }
}
```

### 2.5 获取新上架插件

**接口地址**: `GET /api/v1/store/plugins/new`

**查询参数**:
- limit: 返回数量 (默认10，最大50)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "plugin-stock-query-001",
        "name": "股票查询",
        "description": "股票行情、财务数据",
        "icon": "https://example.com/icons/stock.png",
        "category": "金融工具",
        "pricing_type": "free",
        "rating": 0,
        "review_count": 0,
        "install_count": 125,
        "published_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 2.6 获取插件分类

**接口地址**: `GET /api/v1/store/categories`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "cat-search",
        "name": "搜索工具",
        "icon": "🔍",
        "plugin_count": 15,
        "sort_order": 1
      },
      {
        "id": "cat-communication",
        "name": "通信工具",
        "icon": "📧",
        "plugin_count": 8,
        "sort_order": 2
      },
      {
        "id": "cat-productivity",
        "name": "效率工具",
        "icon": "⚡",
        "plugin_count": 22,
        "sort_order": 3
      }
    ]
  }
}
```

### 2.7 搜索插件

**接口地址**: `GET /api/v1/store/plugins/search`

**查询参数**:
- q: 搜索关键词 (必填)
- category: 分类筛选 (可选)
- pricing_type: 定价类型 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "query": "天气",
    "items": [
      {
        "id": "plugin-weather-001",
        "name": "天气查询",
        "description": "实时天气、天气预报",
        "icon": "https://example.com/icons/weather.png",
        "category": "生活工具",
        "pricing_type": "free",
        "rating": 4.7,
        "review_count": 210,
        "install_count": 25600
      }
    ]
  }
}
```

---

## 3. 插件管理API

### 3.1 安装插件

**接口地址**: `POST /api/v1/store/plugins/{plugin_id}/install`

**请求参数**:
```json
{
  "config": {
    "api_key": "AIzaSyXXXXXXXXXXXX",
    "cx": "012345678901234567890:abcde"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "安装成功",
  "data": {
    "installation_id": "install-001",
    "plugin_id": "plugin-google-search-001",
    "plugin_name": "Google搜索",
    "installed_at": "2025-01-15T14:00:00Z"
  }
}
```

### 3.2 购买付费插件

**接口地址**: `POST /api/v1/store/plugins/{plugin_id}/purchase`

**请求参数**:
```json
{
  "payment_method": "alipay",
  "quantity": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "订单创建成功",
  "data": {
    "order_id": "order-20250115-001",
    "plugin_id": "plugin-email-sender-001",
    "plugin_name": "邮件发送",
    "amount": 99.00,
    "currency": "CNY",
    "payment_url": "https://pay.example.com/checkout/order-20250115-001",
    "expires_at": "2025-01-15T15:00:00Z"
  }
}
```

### 3.3 获取我的插件

**接口地址**: `GET /api/v1/store/my-plugins`

**查询参数**:
- status: 状态筛选 (可选，installed/purchased/published)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "installation_id": "install-001",
        "plugin_id": "plugin-google-search-001",
        "plugin_name": "Google搜索",
        "plugin_icon": "https://example.com/icons/google-search.png",
        "status": "installed",
        "installed_at": "2025-01-01T00:00:00Z",
        "last_used_at": "2025-01-15T10:00:00Z"
      },
      {
        "installation_id": "install-002",
        "plugin_id": "plugin-email-sender-001",
        "plugin_name": "邮件发送",
        "plugin_icon": "https://example.com/icons/email-sender.png",
        "status": "purchased",
        "purchased_at": "2025-01-10T00:00:00Z",
        "amount": 99.00
      }
    ]
  }
}
```

### 3.4 卸载插件

**接口地址**: `DELETE /api/v1/store/my-plugins/{installation_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "卸载成功",
  "data": {
    "installation_id": "install-001",
    "uninstalled_at": "2025-01-15T15:00:00Z"
  }
}
```

### 3.5 上传插件

**接口地址**: `POST /api/v1/store/plugins/upload`

**请求参数** (multipart/form-data):
```json
{
  "plugin_file": "<binary>",
  "icon": "<binary>",
  "screenshots[]": ["<binary>", "<binary>"],
  "name": "My Plugin",
  "description": "Plugin description",
  "category": "搜索工具",
  "pricing_type": "free",
  "price": null,
  "long_description": "<p>Detailed description</p>",
  "config_schema": "[{\"name\":\"api_key\",\"type\":\"string\",\"required\":true}]",
  "permissions": ["http:api.example.com"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "plugin_id": "plugin-my-plugin-001",
    "name": "My Plugin",
    "version": "1.0.0",
    "status": "draft",
    "created_at": "2025-01-15T16:00:00Z"
  }
}
```

### 3.6 更新插件版本

**接口地址**: `POST /api/v1/store/plugins/{plugin_id}/versions`

**请求参数** (multipart/form-data):
```json
{
  "plugin_file": "<binary>",
  "version": "1.1.0",
  "changelog": "修复了XX问题，增加了XX功能"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "版本创建成功",
  "data": {
    "version_id": "ver-001",
    "plugin_id": "plugin-my-plugin-001",
    "version": "1.1.0",
    "changelog": "修复了XX问题，增加了XX功能",
    "status": "draft",
    "created_at": "2025-01-15T17:00:00Z"
  }
}
```

---

## 4. 插件审核API

### 4.1 提交审核

**接口地址**: `POST /api/v1/store/plugins/{plugin_id}/submit-for-review`

**请求参数**:
```json
{
  "version": "1.0.0",
  "changelog": "初始版本发布"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "提交成功，等待审核",
  "data": {
    "audit_id": 456,
    "plugin_id": "plugin-my-plugin-001",
    "version": "1.0.0",
    "status": "pending",
    "submitted_at": "2025-01-15T18:00:00Z"
  }
}
```

### 4.2 获取审核状态

**接口地址**: `GET /api/v1/store/plugins/{plugin_id}/audit-status`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "audit_id": 456,
    "plugin_id": "plugin-my-plugin-001",
    "version": "1.0.0",
    "status": "pending",
    "review_comments": null,
    "reviewed_at": null,
    "submitted_at": "2025-01-15T18:00:00Z"
  }
}
```

### 4.3 管理员审核插件

**接口地址**: `POST /api/v1/store/admin/plugins/{plugin_id}/audit`

**权限**: 仅管理员

**请求参数**:
```json
{
  "audit_id": 456,
  "status": "approved",
  "comments": "审核通过，功能完整，无违规内容"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "审核完成",
  "data": {
    "audit_id": 456,
    "plugin_id": "plugin-my-plugin-001",
    "status": "approved",
    "review_comments": "审核通过，功能完整，无违规内容",
    "reviewer_id": 1001,
    "reviewer_name": "Admin",
    "reviewed_at": "2025-01-15T19:00:00Z"
  }
}
```

### 4.4 获取待审核列表

**接口地址**: `GET /api/v1/store/admin/audits/pending`

**权限**: 仅管理员

**查询参数**:
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 12,
    "items": [
      {
        "audit_id": 456,
        "plugin_id": "plugin-my-plugin-001",
        "plugin_name": "My Plugin",
        "version": "1.0.0",
        "submitter_id": 2001,
        "submitter_name": "Developer",
        "status": "pending",
        "submitted_at": "2025-01-15T18:00:00Z"
      }
    ]
  }
}
```

---

## 5. 评价系统API

### 5.1 提交评价

**接口地址**: `POST /api/v1/store/plugins/{plugin_id}/reviews`

**请求参数**:
```json
{
  "rating": 5,
  "title": "非常好用的插件",
  "content": "功能强大，易于使用，强烈推荐！",
  "verified": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "评价成功",
  "data": {
    "id": 123,
    "plugin_id": "plugin-weather-001",
    "user_id": 3001,
    "username": "用户A",
    "avatar": "https://example.com/avatars/user-a.png",
    "rating": 5,
    "title": "非常好用的插件",
    "content": "功能强大，易于使用，强烈推荐！",
    "verified": true,
    "helpful_count": 0,
    "created_at": "2025-01-15T20:00:00Z"
  }
}
```

### 5.2 获取插件评价

**接口地址**: `GET /api/v1/store/plugins/{plugin_id}/reviews`

**查询参数**:
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)
- rating: 评分筛选 (可选，1-5)
- sort: 排序方式 (可选，newest/helpful/verified)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 150,
    "average_rating": 4.5,
    "rating_distribution": {
      "5": 100,
      "4": 30,
      "3": 15,
      "2": 3,
      "1": 2
    },
    "items": [
      {
        "id": 123,
        "user_id": 3001,
        "username": "用户A",
        "avatar": "https://example.com/avatars/user-a.png",
        "rating": 5,
        "title": "非常好用的插件",
        "content": "功能强大，易于使用",
        "verified": true,
        "helpful_count": 25,
        "reply": "感谢您的评价！",
        "replied_at": "2025-01-15T21:00:00Z",
        "created_at": "2025-01-03T00:00:00Z"
      }
    ]
  }
}
```

### 5.3 标记有用

**接口地址**: `POST /api/v1/store/plugins/reviews/{review_id}/helpful`

**响应示例**:
```json
{
  "code": 0,
  "message": "标记成功",
  "data": {
    "review_id": 123,
    "helpful_count": 26
  }
}
```

### 5.4 开发者回复

**接口地址**: `POST /api/v1/store/plugins/reviews/{review_id}/reply`

**权限**: 仅插件开发者

**请求参数**:
```json
{
  "reply": "感谢您的评价！如有问题请随时联系我们。"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "回复成功",
  "data": {
    "review_id": 123,
    "reply": "感谢您的评价！如有问题请随时联系我们。",
    "replied_at": "2025-01-15T22:00:00Z"
  }
}
```

---

## 6. 收益管理API

### 6.1 获取收益概览

**接口地址**: `GET /api/v1/store/earnings`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_revenue": 12500.50,
    "pending_revenue": 2500.00,
    "withdrawn_revenue": 10000.50,
    "current_month_revenue": 3200.20,
    "currency": "CNY",
    "plugins": [
      {
        "plugin_id": "plugin-email-sender-001",
        "plugin_name": "邮件发送",
        "total_revenue": 8500.00,
        "current_month_revenue": 2200.00,
        "sales_count": 86
      },
      {
        "plugin_id": "plugin-schedule-001",
        "plugin_name": "日程管理",
        "total_revenue": 4000.50,
        "current_month_revenue": 1000.20,
        "sales_count": 45
      }
    ]
  }
}
```

### 6.2 获取收益明细

**接口地址**: `GET /api/v1/store/earnings/transactions`

**查询参数**:
- plugin_id: 插件ID筛选 (可选)
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
    "total": 156,
    "items": [
      {
        "transaction_id": "txn-20250115-001",
        "plugin_id": "plugin-email-sender-001",
        "plugin_name": "邮件发送",
        "type": "sale",
        "amount": 99.00,
        "revenue_share": 0.70,
        "developer_revenue": 69.30,
        "platform_revenue": 29.70,
        "currency": "CNY",
        "status": "completed",
        "purchaser_id": 4001,
        "purchaser_name": "用户B",
        "created_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 6.3 申请提现

**接口地址**: `POST /api/v1/store/earnings/withdraw`

**请求参数**:
```json
{
  "amount": 5000.00,
  "method": "bank_transfer",
  "bank_account": {
    "bank_name": "中国工商银行",
    "account_number": "6222021234567890",
    "account_holder": "张三"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "提现申请已提交",
  "data": {
    "withdrawal_id": "withdraw-20250115-001",
    "amount": 5000.00,
    "fee": 5.00,
    "actual_amount": 4995.00,
    "currency": "CNY",
    "status": "pending",
    "estimated_arrival_date": "2025-01-18T00:00:00Z",
    "created_at": "2025-01-15T23:00:00Z"
  }
}
```

### 6.4 获取提现记录

**接口地址**: `GET /api/v1/store/earnings/withdrawals`

**查询参数**:
- status: 状态筛选 (可选，pending/completed/rejected)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 8,
    "items": [
      {
        "withdrawal_id": "withdraw-20250115-001",
        "amount": 5000.00,
        "fee": 5.00,
        "actual_amount": 4995.00,
        "currency": "CNY",
        "method": "bank_transfer",
        "status": "pending",
        "reject_reason": null,
        "estimated_arrival_date": "2025-01-18T00:00:00Z",
        "created_at": "2025-01-15T23:00:00Z"
      },
      {
        "withdrawal_id": "withdraw-20250110-001",
        "amount": 3000.00,
        "fee": 3.00,
        "actual_amount": 2997.00,
        "currency": "CNY",
        "method": "bank_transfer",
        "status": "completed",
        "reject_reason": null,
        "completed_at": "2025-01-13T00:00:00Z",
        "created_at": "2025-01-10T10:00:00Z"
      }
    ]
  }
}
```

### 6.5 获取收益报表

**接口地址**: `GET /api/v1/store/earnings/report`

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)
- granularity: 时间粒度 (可选，day/week/month，默认day)
- plugin_id: 插件ID筛选 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "summary": {
      "total_revenue": 5200.50,
      "total_sales": 52,
      "avg_revenue_per_sale": 100.01
    },
    "daily_data": [
      {
        "date": "2025-01-01",
        "revenue": 350.00,
        "sales": 3
      },
      {
        "date": "2025-01-02",
        "revenue": 495.00,
        "sales": 5
      }
    ],
    "plugin_breakdown": [
      {
        "plugin_id": "plugin-email-sender-001",
        "plugin_name": "邮件发送",
        "revenue": 3200.00,
        "sales": 32,
        "percentage": 61.5
      }
    ]
  }
}
```

---

## 7. 数据模型

### 7.1 Plugin

```typescript
interface Plugin {
  id: string;
  tenant_id: string;
  name: string;
  description: string;
  long_description?: string;
  icon: string;
  screenshots?: string[];
  category: string;
  pricing_type: 'free' | 'paid' | 'freemium';
  price?: number;
  revenue_share_rate: number;
  total_revenue: number;
  status: 'draft' | 'pending_review' | 'approved' | 'rejected' | 'suspended';
  rating: number;
  review_count: number;
  install_count: number;
  active_users: number;
  author: string;
  author_id: number;
  version: string;
  config_schema: ConfigSchema[];
  permissions: string[];
  usage_example?: string;
  score?: number;
  created_at: Date;
  updated_at: Date;
}

interface ConfigSchema {
  name: string;
  type: 'string' | 'number' | 'boolean' | 'select' | 'textarea';
  required: boolean;
  label: string;
  description?: string;
  default?: any;
  options?: string[]; // for type='select'
}
```

### 7.2 PluginReview

```typescript
interface PluginReview {
  id: number;
  plugin_id: string;
  user_id: number;
  username?: string;
  avatar?: string;
  rating: number; // 1-5
  title?: string;
  content: string;
  verified: boolean;
  helpful_count: number;
  reply?: string;
  replied_at?: Date;
  created_at: Date;
}

interface ReviewSummary {
  total: number;
  average_rating: number;
  rating_distribution: {
    '5': number;
    '4': number;
    '3': number;
    '2': number;
    '1': number;
  };
}
```

### 7.3 PluginAudit

```typescript
interface PluginAudit {
  id: number;
  plugin_id: string;
  version: string;
  submitter_id: number;
  reviewer_id?: number;
  status: 'pending' | 'approved' | 'rejected';
  review_comments?: string;
  reviewed_at?: Date;
  submitted_at: Date;
}
```

### 7.4 Earning

```typescript
interface Earning {
  transaction_id: string;
  plugin_id: string;
  plugin_name: string;
  type: 'sale' | 'refund';
  amount: number;
  revenue_share: number;
  developer_revenue: number;
  platform_revenue: number;
  currency: string;
  status: 'pending' | 'completed' | 'cancelled';
  purchaser_id?: number;
  purchaser_name?: string;
  created_at: Date;
}

interface EarningSummary {
  total_revenue: number;
  pending_revenue: number;
  withdrawn_revenue: number;
  current_month_revenue: number;
  currency: string;
  plugins: PluginEarning[];
}

interface PluginEarning {
  plugin_id: string;
  plugin_name: string;
  total_revenue: number;
  current_month_revenue: number;
  sales_count: number;
}
```

### 7.5 Withdrawal

```typescript
interface Withdrawal {
  withdrawal_id: string;
  amount: number;
  fee: number;
  actual_amount: number;
  currency: string;
  method: 'bank_transfer' | 'alipay' | 'wechat';
  status: 'pending' | 'completed' | 'rejected';
  reject_reason?: string;
  estimated_arrival_date?: Date;
  completed_at?: Date;
  created_at: Date;
}
```

---

## 8. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 34001 | 404 | 插件不存在 |
| 34002 | 400 | 插件名称已存在 |
| 34003 | 400 | 插件版本无效 |
| 34004 | 400 | 配置Schema无效 |
| 34005 | 400 | 权限声明无效 |
| 34006 | 403 | 无权限操作插件 |
| 34007 | 400 | 插件文件格式错误 |
| 34008 | 400 | 插件已安装 |
| 34009 | 404 | 插件未安装 |
| 34101 | 400 | 评价已存在 |
| 34102 | 404 | 评价不存在 |
| 34103 | 400 | 评分无效(1-5) |
| 34104 | 403 | 无权限回复评价 |
| 34201 | 400 | 余额不足 |
| 34202 | 400 | 提现金额低于最低限额 |
| 34203 | 400 | 提现申请过于频繁 |
| 34204 | 400 | 银行账户信息无效 |
| 34205 | 403 | 无权限提现 |
| 34301 | 400 | 插件已在审核中 |
| 34302 | 400 | 插件未提交审核 |
| 34303 | 403 | 无权限审核插件 |

---

**文档结束**
