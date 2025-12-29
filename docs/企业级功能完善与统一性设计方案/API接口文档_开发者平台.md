# API接口文档：开发者平台模块

**模块名称**: Developer Platform (开发者平台)
**设计文档**: 27-开发者平台.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P3

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 开发者管理API](#2-开发者管理api)
- [3. API密钥管理API](#3-api密钥管理api)
- [4. Bot项目管理API](#4-bot项目管理api)
- [5. 插件市场API](#5-插件市场api)
- [6. 统计分析API](#6-统计分析api)
- [7. 贡献管理API](#7-贡献管理api)
- [8. 数据模型](#8-数据模型)
- [9. 错误码定义](#9-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

开发者平台是对标 **Dify**、**Langflow** 的开发者生态建设平台，通过 **CLI 工具 + SDK + API 文档 + 开发者社区**，降低开发门槛，吸引第三方开发者：

- ✅ **开发者管理** - 开发者注册、认证、资料管理
- ✅ **API密钥管理** - 密钥生成、权限控制、使用限额
- ✅ **Bot项目管理** - Bot CRUD、版本管理、部署
- ✅ **插件市场** - 插件上传、审核、发布、搜索
- ✅ **统计分析** - 使用量统计、性能分析、成本分析
- ✅ **贡献激励** - 贡献者等级、徽章系统、积分奖励

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 认证: JWT + API Key
- CLI: Commander.js (Node.js)
- SDK: TypeScript/Python/Go/Java

**前端技术栈**:
- React 18 + TypeScript 5.6
- Semi Design UI组件库
- Swagger UI (API文档)

**CLI工具**:
- `@zeker/cli` - 命令行工具
- `@zeker/sdk` - 多语言SDK

**实现策略**: ✅ 30% 编码（工具框架） + 70% 内容建设（文档、教程、社区）

### 1.3 数据库表

| 表名 | 说明 |
|------|------|
| `developers` | 开发者表 |
| `api_keys` | API密钥表 |
| `bots` | Bot项目表 |
| `bot_versions` | Bot版本表 |
| `plugins_market` | 插件市场表 |
| `plugin_reviews` | 插件评价表 |
| `contributions` | 贡献记录表 |
| `badges` | 徽章表 |

---

## 2. 开发者管理API

### 2.1 开发者注册

**接口地址**: `POST /api/v1/developers/register`

**请求参数**:
```json
{
  "developer_name": "张三",
  "email": "zhangsan@example.com",
  "github_id": "zhangsan-github",
  "github_url": "https://github.com/zhangsan-github",
  "company": "示例公司",
  "position": "软件工程师"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "developer_id": 1001,
    "developer_name": "张三",
    "email": "zhangsan@example.com",
    "github_id": "zhangsan-github",
    "github_url": "https://github.com/zhangsan-github",
    "company": "示例公司",
    "position": "软件工程师",
    "level": "bronze",
    "reputation_score": 0,
    "created_at": "2025-01-15T10:00:00Z"
  }
}
```

### 2.2 获取开发者资料

**接口地址**: `GET /api/v1/developers/{developer_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "developer_id": 1001,
    "developer_name": "张三",
    "email": "zhangsan@example.com",
    "github_id": "zhangsan-github",
    "github_url": "https://github.com/zhangsan-github",
    "company": "示例公司",
    "position": "软件工程师",
    "level": "gold",
    "reputation_score": 1250,
    "total_commits": 85,
    "total_prs": 32,
    "total_issues_resolved": 28,
    "badges": [
      {
        "badge_id": "bug-hunter",
        "name": "Bug Hunter",
        "description": "发现并修复 5 个 Bug",
        "earned_at": "2025-01-10T15:30:00Z"
      },
      {
        "badge_id": "code-contributor",
        "name": "Code Contributor",
        "description": "提交 10 个 PR",
        "earned_at": "2025-01-12T10:20:00Z"
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

### 2.3 更新开发者资料

**接口地址**: `PUT /api/v1/developers/profile`

**权限**: 需要认证

**请求参数**:
```json
{
  "developer_name": "张三(更新)",
  "company": "新公司",
  "position": "高级软件工程师",
  "github_url": "https://github.com/zhangsan-github-updated"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "developer_id": 1001,
    "updated_at": "2025-01-15T11:00:00Z"
  }
}
```

### 2.4 获取开发者排行

**接口地址**: `GET /api/v1/developers/ranking`

**查询参数**:
- period: 时间周期 (week/month/year/all，默认all)
- limit: 返回数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "month",
    "items": [
      {
        "rank": 1,
        "developer_id": 1001,
        "developer_name": "张三",
        "level": "platinum",
        "reputation_score": 2850,
        "total_commits": 120,
        "total_prs": 45
      },
      {
        "rank": 2,
        "developer_id": 1002,
        "developer_name": "李四",
        "level": "gold",
        "reputation_score": 2100,
        "total_commits": 95,
        "total_prs": 38
      }
    ]
  }
}
```

### 2.5 升级开发者等级

**接口地址**: `POST /api/v1/developers/{developer_id}/level-up`

**权限**: 管理员

**请求参数**:
```json
{
  "new_level": "platinum",
  "reason": "连续3个月排名第一"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "升级成功",
  "data": {
    "developer_id": 1001,
    "old_level": "gold",
    "new_level": "platinum",
    "upgraded_at": "2025-01-15T12:00:00Z"
  }
}
```

---

## 3. API密钥管理API

### 3.1 创建API密钥

**接口地址**: `POST /api/v1/developers/api-keys`

**权限**: 需要认证

**请求参数**:
```json
{
  "name": "生产环境密钥",
  "description": "用于生产环境API调用",
  "scopes": ["bots:read", "bots:write", "conversations:read"],
  "rate_limit": 1000,
  "expires_at": "2026-01-15T00:00:00Z"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "key_id": "key-123456",
    "api_key": "sk_live_abcdef1234567890abcdef1234567890",
    "name": "生产环境密钥",
    "description": "用于生产环境API调用",
    "scopes": ["bots:read", "bots:write", "conversations:read"],
    "rate_limit": 1000,
    "expires_at": "2026-01-15T00:00:00Z",
    "created_at": "2025-01-15T14:00:00Z",
    "warning": "请妥善保管您的API密钥，此信息仅显示一次"
  }
}
```

### 3.2 获取API密钥列表

**接口地址**: `GET /api/v1/developers/api-keys`

**权限**: 需要认证

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 3,
    "items": [
      {
        "key_id": "key-123456",
        "name": "生产环境密钥",
        "api_key_preview": "sk_live_****************7890",
        "scopes": ["bots:read", "bots:write", "conversations:read"],
        "rate_limit": 1000,
        "last_used_at": "2025-01-15T10:30:00Z",
        "is_active": true,
        "expires_at": "2026-01-15T00:00:00Z",
        "created_at": "2025-01-01T00:00:00Z"
      },
      {
        "key_id": "key-789012",
        "name": "测试环境密钥",
        "api_key_preview": "sk_test_****************1234",
        "scopes": ["bots:read"],
        "rate_limit": 100,
        "last_used_at": "2025-01-14T16:20:00Z",
        "is_active": true,
        "expires_at": "2025-07-15T00:00:00Z",
        "created_at": "2025-01-10T00:00:00Z"
      }
    ]
  }
}
```

### 3.3 删除API密钥

**接口地址**: `DELETE /api/v1/developers/api-keys/{key_id}`

**权限**: 需要认证

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "key_id": "key-123456",
    "deleted_at": "2025-01-15T15:00:00Z"
  }
}
```

### 3.4 禁用/启用API密钥

**接口地址**: `PUT /api/v1/developers/api-keys/{key_id}/status`

**权限**: 需要认证

**请求参数**:
```json
{
  "is_active": false,
  "reason": "密钥泄露，临时禁用"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "状态更新成功",
  "data": {
    "key_id": "key-123456",
    "is_active": false,
    "updated_at": "2025-01-15T15:30:00Z"
  }
}
```

### 3.5 获取API密钥使用统计

**接口地址**: `GET /api/v1/developers/api-keys/{key_id}/stats`

**权限**: 需要认证

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "key_id": "key-123456",
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "summary": {
      "total_requests": 15230,
      "successful_requests": 14850,
      "failed_requests": 380,
      "success_rate": 0.975,
      "avg_response_time_ms": 120,
      "total_tokens_used": 1250000
    },
    "daily_usage": [
      {
        "date": "2025-01-01",
        "requests": 1050,
        "success_rate": 0.98,
        "avg_response_time_ms": 115
      },
      {
        "date": "2025-01-02",
        "requests": 1120,
        "success_rate": 0.97,
        "avg_response_time_ms": 125
      }
    ]
  }
}
```

---

## 4. Bot项目管理API

### 4.1 创建Bot

**接口地址**: `POST /api/v1/bots`

**权限**: 需要认证

**请求参数**:
```json
{
  "name": "客服助手",
  "description": "智能客服助手，支持常见问题解答",
  "avatar": "https://example.com/avatars/customer-service.png",
  "prompt_template": "你是一个专业的客服助手...",
  "model": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2000,
  "category": "customer_service"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "bot_id": "bot-123456",
    "name": "客服助手",
    "description": "智能客服助手，支持常见问题解答",
    "avatar": "https://example.com/avatars/customer-service.png",
    "version": "1.0.0",
    "status": "draft",
    "api_endpoint": "https://api.zeker.com/v1/bots/bot-123456",
    "created_at": "2025-01-15T16:00:00Z"
  }
}
```

### 4.2 获取Bot列表

**接口地址**: `GET /api/v1/bots`

**权限**: 需要认证

**查询参数**:
- keyword: 搜索关键词 (可选)
- category: 分类筛选 (可选)
- status: 状态筛选 (draft/published/archived，可选)
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
        "bot_id": "bot-123456",
        "name": "客服助手",
        "description": "智能客服助手",
        "avatar": "https://example.com/avatars/customer-service.png",
        "version": "1.0.0",
        "status": "published",
        "category": "customer_service",
        "total_conversations": 5250,
        "active_users": 890,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      },
      {
        "bot_id": "bot-789012",
        "name": "销售助手",
        "description": "智能销售助手",
        "avatar": "https://example.com/avatars/sales-assistant.png",
        "version": "2.1.0",
        "status": "published",
        "category": "sales",
        "total_conversations": 3200,
        "active_users": 560,
        "created_at": "2025-01-05T00:00:00Z",
        "updated_at": "2025-01-14T15:30:00Z"
      }
    ]
  }
}
```

### 4.3 获取Bot详情

**接口地址**: `GET /api/v1/bots/{bot_id}`

**权限**: 需要认证

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot-123456",
    "developer_id": 1001,
    "name": "客服助手",
    "description": "智能客服助手，支持常见问题解答",
    "avatar": "https://example.com/avatars/customer-service.png",
    "version": "1.0.0",
    "status": "published",
    "category": "customer_service",
    "prompt_template": "你是一个专业的客服助手...",
    "model": "gpt-4",
    "temperature": 0.7,
    "max_tokens": 2000,
    "api_endpoint": "https://api.zeker.com/v1/bots/bot-123456",
    "webhook_url": null,
    "total_conversations": 5250,
    "active_users": 890,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

### 4.4 更新Bot

**接口地址**: `PUT /api/v1/bots/{bot_id}`

**权限**: 需要认证，仅Bot所有者

**请求参数**:
```json
{
  "name": "客服助手(增强版)",
  "description": "智能客服助手，支持常见问题解答和多轮对话",
  "prompt_template": "你是一个专业的客服助手，请友好地回答用户问题...",
  "temperature": 0.8
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "bot_id": "bot-123456",
    "updated_at": "2025-01-15T17:00:00Z"
  }
}
```

### 4.5 发布Bot

**接口地址**: `POST /api/v1/bots/{bot_id}/publish`

**权限**: 需要认证，仅Bot所有者

**请求参数**:
```json
{
  "version": "1.0.0",
  "changelog": "首次发布"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "发布成功",
  "data": {
    "bot_id": "bot-123456",
    "version": "1.0.0",
    "status": "published",
    "published_at": "2025-01-15T18:00:00Z",
    "public_url": "https://app.zeker.com/bots/bot-123456"
  }
}
```

### 4.6 删除Bot

**接口地址**: `DELETE /api/v1/bots/{bot_id}`

**权限**: 需要认证，仅Bot所有者

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "bot_id": "bot-123456",
    "deleted_at": "2025-01-15T19:00:00Z"
  }
}
```

### 4.7 获取Bot版本历史

**接口地址**: `GET /api/v1/bots/{bot_id}/versions`

**权限**: 需要认证

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot-123456",
    "versions": [
      {
        "version_id": "ver-001",
        "version": "1.0.0",
        "changelog": "首次发布",
        "status": "published",
        "created_at": "2025-01-01T00:00:00Z",
        "created_by": {
          "developer_id": 1001,
          "developer_name": "张三"
        }
      },
      {
        "version_id": "ver-002",
        "version": "1.1.0",
        "changelog": "新增多轮对话支持",
        "status": "published",
        "created_at": "2025-01-10T00:00:00Z",
        "created_by": {
          "developer_id": 1001,
          "developer_name": "张三"
        }
      }
    ]
  }
}
```

---

## 5. 插件市场API

### 5.1 上传插件

**接口地址**: `POST /api/v1/plugins/submit`

**权限**: 需要认证

**请求参数**:
```json
{
  "name": "Web Search Plugin",
  "description": "提供网络搜索能力",
  "category": "search",
  "version": "1.0.0",
  "repository_url": "https://github.com/zhangsan/web-search-plugin",
  "homepage_url": "https://github.com/zhangsan/web-search-plugin#readme",
  "readme": "# Web Search Plugin\n\n这是一个强大的网络搜索插件...",
  "icon": "https://example.com/icons/web-search.png",
  "tags": ["search", "web", "google"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "提交成功，等待审核",
  "data": {
    "plugin_id": "plugin-123456",
    "name": "Web Search Plugin",
    "version": "1.0.0",
    "status": "pending",
    "submitted_at": "2025-01-15T20:00:00Z",
    "message": "您的插件已提交审核，预计1-3个工作日内完成"
  }
}
```

### 5.2 获取插件列表

**接口地址**: `GET /api/v1/plugins`

**查询参数**:
- keyword: 搜索关键词 (可选)
- category: 分类筛选 (可选)
- status: 状态筛选 (approved/pending/rejected，可选)
- sort: 排序方式 (downloads/stars/rating/updated，默认downloads)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 85,
    "items": [
      {
        "plugin_id": "plugin-web-search-001",
        "name": "Web Search Plugin",
        "description": "提供网络搜索能力，支持Google、Bing等搜索引擎",
        "category": "search",
        "version": "1.2.0",
        "icon": "https://example.com/icons/web-search.png",
        "developer": {
          "developer_id": 1001,
          "developer_name": "张三"
        },
        "repository_url": "https://github.com/zhangsan/web-search-plugin",
        "homepage_url": "https://github.com/zhangsan/web-search-plugin#readme",
        "downloads": 12500,
        "stars": 890,
        "rating": 4.8,
        "review_count": 125,
        "status": "approved",
        "tags": ["search", "web", "google"],
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      },
      {
        "plugin_id": "plugin-knowledge-base-001",
        "name": "Knowledge Base Plugin",
        "description": "知识库检索插件，支持向量搜索",
        "category": "knowledge",
        "version": "2.0.0",
        "icon": "https://example.com/icons/knowledge-base.png",
        "developer": {
          "developer_id": 1002,
          "developer_name": "李四"
        },
        "repository_url": "https://github.com/lisi/knowledge-base-plugin",
        "homepage_url": "https://github.com/lisi/knowledge-base-plugin#readme",
        "downloads": 8900,
        "stars": 650,
        "rating": 4.6,
        "review_count": 89,
        "status": "approved",
        "tags": ["knowledge", "vector", "search"],
        "created_at": "2025-01-05T00:00:00Z",
        "updated_at": "2025-01-14T15:30:00Z"
      }
    ]
  }
}
```

### 5.3 获取插件详情

**接口地址**: `GET /api/v1/plugins/{plugin_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "plugin_id": "plugin-web-search-001",
    "name": "Web Search Plugin",
    "description": "提供网络搜索能力，支持Google、Bing等搜索引擎，支持自定义搜索结果数量和语言",
    "category": "search",
    "version": "1.2.0",
    "icon": "https://example.com/icons/web-search.png",
    "readme": "# Web Search Plugin\n\n这是一个强大的网络搜索插件...",
    "developer": {
      "developer_id": 1001,
      "developer_name": "张三",
      "level": "gold",
      "avatar": "https://example.com/avatars/zhangsan.png"
    },
    "repository_url": "https://github.com/zhangsan/web-search-plugin",
    "homepage_url": "https://github.com/zhangsan/web-search-plugin#readme",
    "downloads": 12500,
    "stars": 890,
    "rating": 4.8,
    "review_count": 125,
    "status": "approved",
    "tags": ["search", "web", "google"],
    "installation_command": "zeker plugin install @zeker/plugin-web-search",
    "usage_example": {
      "title": "使用示例",
      "code": "import { WebSearchPlugin } from '@zeker/plugin-web-search';\n\nconst plugin = new WebSearchPlugin({\n  engine: 'google',\n  maxResults: 10\n});\n\nconst results = await plugin.search('AI技术趋势');"
    },
    "changelog": [
      {
        "version": "1.2.0",
        "date": "2025-01-15",
        "changes": [
          "新增Bing搜索引擎支持",
          "优化搜索结果排序算法"
        ]
      },
      {
        "version": "1.1.0",
        "date": "2025-01-10",
        "changes": [
          "新增自定义搜索结果数量",
          "修复多语言搜索问题"
        ]
      }
    ],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

### 5.4 审核插件

**接口地址**: `POST /api/v1/plugins/{plugin_id}/review`

**权限**: 管理员

**请求参数**:
```json
{
  "action": "approve",
  "comment": "插件功能完善，代码质量优秀，同意发布"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "审核完成",
  "data": {
    "plugin_id": "plugin-123456",
    "status": "approved",
    "reviewed_at": "2025-01-15T21:00:00Z",
    "reviewed_by": {
      "admin_id": 1,
      "admin_name": "管理员"
    },
    "comment": "插件功能完善，代码质量优秀，同意发布"
  }
}
```

### 5.5 更新插件版本

**接口地址**: `POST /api/v1/plugins/{plugin_id}/version`

**权限**: 需要认证，仅插件开发者

**请求参数**:
```json
{
  "version": "1.3.0",
  "changelog": "新增高级搜索过滤器",
  "repository_url": "https://github.com/zhangsan/web-search-plugin",
  "readme": "# Web Search Plugin v1.3.0\n\n更新内容..."
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "版本更新成功，等待审核",
  "data": {
    "plugin_id": "plugin-web-search-001",
    "new_version": "1.3.0",
    "status": "pending",
    "submitted_at": "2025-01-15T22:00:00Z"
  }
}
```

### 5.6 提交插件评价

**接口地址**: `POST /api/v1/plugins/{plugin_id}/reviews`

**权限**: 需要认证

**请求参数**:
```json
{
  "rating": 5,
  "title": "非常好用的搜索插件",
  "content": "插件功能强大，API设计简洁，文档完善，强烈推荐！",
  "is_recommended": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "评价提交成功",
  "data": {
    "review_id": "review-789012",
    "plugin_id": "plugin-web-search-001",
    "rating": 5,
    "title": "非常好用的搜索插件",
    "content": "插件功能强大，API设计简洁，文档完善，强烈推荐！",
    "is_recommended": true,
    "created_at": "2025-01-15T22:30:00Z"
  }
}
```

### 5.7 获取插件评价列表

**接口地址**: `GET /api/v1/plugins/{plugin_id}/reviews`

**查询参数**:
- rating: 评分筛选 (1-5，可选)
- sort: 排序方式 (recent/helpful，默认recent)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 125,
    "items": [
      {
        "review_id": "review-789012",
        "user": {
          "user_id": 2001,
          "username": "user001",
          "avatar": "https://example.com/avatars/user001.png"
        },
        "rating": 5,
        "title": "非常好用的搜索插件",
        "content": "插件功能强大，API设计简洁，文档完善，强烈推荐！",
        "is_recommended": true,
        "helpful_count": 25,
        "created_at": "2025-01-15T22:30:00Z",
        "developer_reply": {
          "content": "感谢您的支持！",
          "replied_at": "2025-01-16T08:00:00Z"
        }
      }
    ]
  }
}
```

---

## 6. 统计分析API

### 6.1 获取开发者统计数据

**接口地址**: `GET /api/v1/developers/stats`

**权限**: 需要认证

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)

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
      "total_bots": 12,
      "published_bots": 8,
      "total_plugins": 5,
      "published_plugins": 3,
      "total_api_calls": 152000,
      "total_tokens_used": 12500000,
      "total_conversations": 25800,
      "active_users": 3250
    },
    "bots_stats": [
      {
        "bot_id": "bot-123456",
        "bot_name": "客服助手",
        "conversations": 5250,
        "api_calls": 45000,
        "tokens_used": 3800000,
        "active_users": 890
      }
    ],
    "plugins_stats": [
      {
        "plugin_id": "plugin-web-search-001",
        "plugin_name": "Web Search Plugin",
        "downloads": 2500,
        "installations": 1250,
        "stars": 85
      }
    ]
  }
}
```

### 6.2 获取Bot使用统计

**接口地址**: `GET /api/v1/bots/{bot_id}/stats`

**权限**: 需要认证

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)
- granularity: 时间粒度 (day/week/month，默认day)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot-123456",
    "bot_name": "客服助手",
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "summary": {
      "total_conversations": 5250,
      "total_messages": 15800,
      "total_tokens_used": 3800000,
      "total_cost": 152.0,
      "avg_response_time_ms": 320,
      "satisfaction_rate": 0.92
    },
    "daily_stats": [
      {
        "date": "2025-01-01",
        "conversations": 350,
        "messages": 1050,
        "tokens_used": 250000,
        "cost": 10.0,
        "avg_response_time_ms": 315,
        "satisfaction_rate": 0.93
      },
      {
        "date": "2025-01-02",
        "conversations": 380,
        "messages": 1150,
        "tokens_used": 275000,
        "cost": 11.0,
        "avg_response_time_ms": 325,
        "satisfaction_rate": 0.91
      }
    ]
  }
}
```

### 6.3 获取插件下载统计

**接口地址**: `GET /api/v1/plugins/{plugin_id}/stats`

**权限**: 需要认证

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "plugin_id": "plugin-web-search-001",
    "plugin_name": "Web Search Plugin",
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "summary": {
      "total_downloads": 2500,
      "total_installations": 1250,
      "total_uninstalls": 25,
      "active_installations": 1225,
      "new_installs": 850,
      "stars": 85,
      "reviews": 25
    },
    "daily_stats": [
      {
        "date": "2025-01-01",
        "downloads": 150,
        "installations": 85,
        "uninstalls": 2,
        "stars": 5
      },
      {
        "date": "2025-01-02",
        "downloads": 180,
        "installations": 95,
        "uninstalls": 1,
        "stars": 8
      }
    ]
  }
}
```

### 6.4 获取平台总体统计

**接口地址**: `GET /api/v1/platform/stats`

**权限**: 管理员

**查询参数**:
- period: 统计周期 (week/month/year，默认month)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "month",
    "date_range": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "developers": {
      "total": 5250,
      "new_this_period": 350,
      "active": 1200
    },
    "bots": {
      "total": 12500,
      "new_this_period": 850,
      "published": 8900
    },
    "plugins": {
      "total": 850,
      "new_this_period": 45,
      "published": 520
    },
    "api_usage": {
      "total_requests": 15200000,
      "total_tokens": 1250000000,
      "success_rate": 0.985
    },
    "conversations": {
      "total": 2580000,
      "active_users": 325000
    }
  }
}
```

---

## 7. 贡献管理API

### 7.1 记录贡献

**接口地址**: `POST /api/v1/contributions`

**权限**: 系统内部或管理员

**请求参数**:
```json
{
  "developer_id": 1001,
  "contribution_type": "pull_request",
  "description": "修复Bot创建时的参数验证问题",
  "url": "https://github.com/zeker/platform/pull/123",
  "points": 10
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "贡献记录成功",
  "data": {
    "contribution_id": "contrib-001",
    "developer_id": 1001,
    "contribution_type": "pull_request",
    "description": "修复Bot创建时的参数验证问题",
    "url": "https://github.com/zeker/platform/pull/123",
    "points": 10,
    "recorded_at": "2025-01-15T23:00:00Z"
  }
}
```

### 7.2 获取贡献记录

**接口地址**: `GET /api/v1/developers/{developer_id}/contributions`

**查询参数**:
- type: 贡献类型筛选 (commit/pr/issue/tutorial，可选)
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
    "developer_id": 1001,
    "developer_name": "张三",
    "total_points": 1250,
    "total_contributions": 85,
    "items": [
      {
        "contribution_id": "contrib-001",
        "contribution_type": "pull_request",
        "description": "修复Bot创建时的参数验证问题",
        "url": "https://github.com/zeker/platform/pull/123",
        "points": 10,
        "created_at": "2025-01-15T10:00:00Z"
      },
      {
        "contribution_id": "contrib-002",
        "contribution_type": "tutorial",
        "description": "编写Bot开发入门教程",
        "url": "https://docs.zeker.com/tutorials/bot-101",
        "points": 25,
        "created_at": "2025-01-14T15:30:00Z"
      }
    ]
  }
}
```

### 7.3 授予徽章

**接口地址**: `POST /api/v1/developers/{developer_id}/badges`

**权限**: 管理员

**请求参数**:
```json
{
  "badge_type": "bug_hunter",
  "reason": "发现并修复5个Bug"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "徽章授予成功",
  "data": {
    "badge_id": "badge-001",
    "developer_id": 1001,
    "badge_type": "bug_hunter",
    "name": "Bug Hunter",
    "description": "发现并修复 5 个 Bug",
    "icon": "https://example.com/badges/bug-hunter.png",
    "earned_at": "2025-01-15T23:30:00Z"
  }
}
```

### 7.4 获取开发者徽章

**接口地址**: `GET /api/v1/developers/{developer_id}/badges`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "developer_id": 1001,
    "badges": [
      {
        "badge_id": "badge-001",
        "badge_type": "bug_hunter",
        "name": "Bug Hunter",
        "description": "发现并修复 5 个 Bug",
        "icon": "https://example.com/badges/bug-hunter.png",
        "earned_at": "2025-01-10T15:30:00Z"
      },
      {
        "badge_id": "badge-002",
        "badge_type": "code_contributor",
        "name": "Code Contributor",
        "description": "提交 10 个 PR",
        "icon": "https://example.com/badges/code-contributor.png",
        "earned_at": "2025-01-12T10:20:00Z"
      },
      {
        "badge_id": "badge-003",
        "badge_type": "tutorial_writer",
        "name": "Tutorial Writer",
        "description": "编写 5 篇教程",
        "icon": "https://example.com/badges/tutorial-writer.png",
        "earned_at": "2025-01-14T08:00:00Z"
      }
    ]
  }
}
```

---

## 8. 数据模型

### 8.1 Developer

```typescript
interface Developer {
  developer_id: number;
  user_id: number;
  developer_name: string;
  email: string;
  github_id?: string;
  github_url?: string;
  company?: string;
  position?: string;
  level: 'bronze' | 'silver' | 'gold' | 'platinum';
  reputation_score: number;
  total_commits: number;
  total_prs: number;
  total_issues_resolved: number;
  badges?: Badge[];
  created_at: Date;
  updated_at: Date;
}
```

### 8.2 APIKey

```typescript
interface APIKey {
  key_id: string;
  developer_id: number;
  api_key: string; // 仅创建时返回完整密钥
  api_key_preview?: string; // 列表时返回脱敏密钥
  name: string;
  description?: string;
  scopes: string[];
  rate_limit: number;
  is_active: boolean;
  expires_at?: Date;
  last_used_at?: Date;
  created_at: Date;
}
```

### 8.3 Bot

```typescript
interface Bot {
  bot_id: string;
  developer_id: number;
  name: string;
  description?: string;
  avatar?: string;
  version: string;
  status: 'draft' | 'published' | 'archived';
  category?: string;
  prompt_template?: string;
  model?: string;
  temperature?: number;
  max_tokens?: number;
  api_endpoint: string;
  webhook_url?: string;
  total_conversations?: number;
  active_users?: number;
  created_at: Date;
  updated_at: Date;
}

interface BotVersion {
  version_id: string;
  bot_id: string;
  version: string;
  changelog?: string;
  status: 'draft' | 'published';
  created_at: Date;
  created_by: {
    developer_id: number;
    developer_name: string;
  };
}
```

### 8.4 Plugin

```typescript
interface Plugin {
  plugin_id: string;
  developer_id: number;
  name: string;
  description: string;
  category: string;
  version: string;
  icon?: string;
  readme?: string;
  repository_url: string;
  homepage_url?: string;
  downloads: number;
  stars: number;
  rating: number;
  review_count: number;
  status: 'pending' | 'approved' | 'rejected';
  tags: string[];
  installation_command?: string;
  usage_example?: {
    title: string;
    code: string;
  };
  developer: {
    developer_id: number;
    developer_name: string;
    level?: string;
    avatar?: string;
  };
  created_at: Date;
  updated_at: Date;
}

interface PluginReview {
  review_id: string;
  plugin_id: string;
  user: {
    user_id: number;
    username: string;
    avatar?: string;
  };
  rating: number;
  title?: string;
  content?: string;
  is_recommended: boolean;
  helpful_count: number;
  created_at: Date;
  developer_reply?: {
    content: string;
    replied_at: Date;
  };
}
```

### 8.5 Contribution

```typescript
interface Contribution {
  contribution_id: string;
  developer_id: number;
  contribution_type: 'commit' | 'pull_request' | 'issue' | 'tutorial' | 'blog';
  description: string;
  url?: string;
  points: number;
  created_at: Date;
}

interface Badge {
  badge_id: string;
  badge_type: string;
  name: string;
  description: string;
  icon?: string;
  earned_at: Date;
}
```

### 8.6 Statistics

```typescript
interface DeveloperStats {
  period: {
    start: string;
    end: string;
  };
  summary: {
    total_bots: number;
    published_bots: number;
    total_plugins: number;
    published_plugins: number;
    total_api_calls: number;
    total_tokens_used: number;
    total_conversations: number;
    active_users: number;
  };
  bots_stats?: BotStats[];
  plugins_stats?: PluginStats[];
}

interface BotStats {
  bot_id: string;
  bot_name: string;
  conversations: number;
  api_calls: number;
  tokens_used: number;
  active_users: number;
}
```

---

## 9. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 40001 | 400 | 开发者名称已存在 |
| 40002 | 400 | 邮箱已注册 |
| 40003 | 400 | GitHub ID 已绑定 |
| 40004 | 403 | 无权限访问此资源 |
| 40005 | 404 | 开发者不存在 |
| 40006 | 400 | API密钥数量已达上限 |
| 40007 | 401 | API密钥无效或已过期 |
| 40008 | 403 | API密钥权限不足 |
| 40009 | 429 | API调用频率超限 |
| 40010 | 400 | Bot名称已存在 |
| 40011 | 404 | Bot不存在 |
| 40012 | 403 | 无权限操作此Bot |
| 40013 | 400 | Bot状态不允许此操作 |
| 40014 | 400 | 插件名称已存在 |
| 40015 | 404 | 插件不存在 |
| 40016 | 400 | 插件已提交审核，请勿重复提交 |
| 40017 | 400 | 插件仓库地址无效 |
| 40018 | 403 | 无权限操作此插件 |
| 40019 | 400 | 评价已存在 |
| 40020 | 400 | 徽章已授予 |
| 40021 | 400 | 贡献类型无效 |
| 40022 | 400 | 开发者等级无效 |

---

**文档结束**
