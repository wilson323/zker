# API接口文档:Bot商店与模板库模块

**文档编号**: DE-API-2025-DOC-029
**版本**: v1.0.0
**创建日期**: 2025-12-30
**参考设计**: 《29-Bot商店与模板库.md》

---

## 📋 目录

1. [Bot商店发布API](#1-bot商店发布api)
2. [Bot商店浏览API](#2-bot商店浏览api)
3. [Bot搜索API](#3-bot搜索api)
4. [Bot一键复克API](#4-bot一键复克api)
5. [Bot评价系统API](#5-bot评价系统api)
6. [Bot排行榜API](#6-bot排行榜api)
7. [Bot模板库API](#7-bot模板库api)
8. [审核管理API](#8-审核管理api)

---

## 1. Bot商店发布API

### 1.1 创建商店条目

**接口描述**: 开发者创建Bot商店条目(草稿状态)

**请求方式**: `POST /api/v1/bot-store/listings`

**权限要求**: `bot_store:write`,Bot所有者

**请求体**:

```json
{
  "data": {
    "bot_id": 789,
    "title": "智能客服助手",
    "description": "7x24小时智能客服,支持多轮对话、知识库问答",
    "long_description": "## 功能介绍\n\n这是一个基于AI的智能客服助手...",
    "category": "客服支持",
    "subcategory": "在线客服",
    "tags": ["客服", "问答", "AI", "多轮对话"],
    "pricing_model": "FREE",
    "screenshots": [
      "https://cdn.coze.com/screenshots/789/1.png",
      "https://cdn.coze.com/screenshots/789/2.png",
      "https://cdn.coze.com/screenshots/789/3.png"
    ],
    "video_url": "https://cdn.coze.com/videos/789/demo.mp4",
    "language": "zh-CN",
    "features": [
      "智能问答",
      "多轮对话",
      "知识库集成",
      "情感分析"
    ],
    "use_cases": [
      "电商客服",
      "产品咨询",
      "售后支持"
    ],
    "version": "1.0.0",
    "changelog": "初始版本发布"
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| bot_id | Long | ✅ | Bot ID |
| title | String | ✅ | 标题(2-100字符) |
| description | String | ✅ | 简短描述(10-500字符) |
| long_description | String | ✅ | 详细描述(支持Markdown) |
| category | String | ✅ | 主分类 |
| subcategory | String | 否 | 子分类 |
| tags | Array | ✅ | 标签列表(最多10个) |
| pricing_model | Enum | ✅ | 定价模式(FREE/PAID) |
| screenshots | Array | ✅ | 截图URL列表(3-5张) |
| video_url | String | 否 | 演示视频URL |
| language | String | ✅ | 主要语言(zh-CN/en-US) |
| features | Array | ✅ | 功能特性列表 |
| use_cases | Array | 否 | 使用场景列表 |
| version | String | ✅ | 版本号 |
| changelog | String | ✅ | 更新日志 |

**响应示例**:

```json
{
  "code": 0,
  "message": "Listing created successfully",
  "data": {
    "listing_id": 123,
    "bot_id": 789,
    "title": "智能客服助手",
    "status": "DRAFT",
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-30T10:00:00Z"
  }
}
```

### 1.2 更新商店条目

**接口描述**: 更新Bot商店条目信息

**请求方式**: `PATCH /api/v1/bot-store/listings/{listing_id}`

**权限要求**: `bot_store:write`,条目所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "title": "智能客服助手 Pro",
    "description": "升级版智能客服,支持语音交互",
    "tags": ["客服", "问答", "AI", "语音交互"],
    "version": "1.1.0",
    "changelog": "- 新增语音交互功能\n- 优化知识库检索\n- 修复已知问题"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Listing updated successfully",
  "data": {
    "listing_id": 123,
    "title": "智能客服助手 Pro",
    "version": "1.1.0",
    "updated_at": "2025-12-30T11:00:00Z"
  }
}
```

### 1.3 提交审核

**接口描述**: 将草稿状态的商店条目提交审核

**请求方式**: `POST /api/v1/bot-store/listings/{listing_id}/submit`

**权限要求**: `bot_store:write`,条目所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "release_notes": "首个稳定版本,已完整测试"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Listing submitted for review",
  "data": {
    "listing_id": 123,
    "status": "PENDING_REVIEW",
    "submitted_at": "2025-12-30T12:00:00Z",
    "estimated_review_time": "1-3个工作日"
  }
}
```

### 1.4 下架商店条目

**接口描述**: 下架已发布的商店条目

**请求方式**: `POST /api/v1/bot-store/listings/{listing_id}/unpublish`

**权限要求**: `bot_store:write`,条目所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "reason": "需要更新版本"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Listing unpublished",
  "data": {
    "listing_id": 123,
    "status": "DRAFT",
    "unpublished_at": "2025-12-30T13:00:00Z"
  }
}
```

---

## 2. Bot商店浏览API

### 2.1 获取首页数据

**接口描述**: 获取Bot商店首页数据(精选、热门、新星、高分Bot)

**请求方式**: `GET /api/v1/bot-store/home`

**权限要求**: 无需认证或需要认证(个性化推荐)

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "featured": [
      {
        "listing_id": 101,
        "bot_id": 801,
        "title": "智能写作助手",
        "description": "AI驱动的文章创作工具",
        "avatar": "https://cdn.coze.com/bot/801/avatar.png",
        "developer": {
          "id": 101,
          "name": "Coze官方",
          "avatar": "https://cdn.coze.com/user/101/avatar.png"
        },
        "category": "生产力工具",
        "pricing_model": "FREE",
        "rating_avg": 4.8,
        "rating_count": 1250,
        "clone_count": 15000,
        "is_featured": true
      }
    ],
    "popular": [
      {
        "listing_id": 102,
        "title": "英语学习Bot",
        "clone_count": 25000,
        "rating_avg": 4.7
      }
    ],
    "new": [
      {
        "listing_id": 103,
        "title": "代码审查助手",
        "published_at": "2025-12-28T10:00:00Z",
        "rating_avg": 4.9
      }
    ],
    "top_rated": [
      {
        "listing_id": 104,
        "title": "数据分析专家",
        "rating_avg": 5.0,
        "rating_count": 89
      }
    ],
    "categories": [
      {
        "id": 1,
        "name": "客服支持",
        "icon": "https://cdn.coze.com/categories/cs.png",
        "count": 1250
      },
      {
        "id": 2,
        "name": "生产力工具",
        "icon": "https://cdn.coze.com/categories/prod.png",
        "count": 890
      }
    ]
  }
}
```

### 2.2 获取商店条目列表

**接口描述**: 获取商店条目列表,支持分页和筛选

**请求方式**: `GET /api/v1/bot-store/listings`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20,最大50) |
| category | String | 否 | 分类筛选 |
| pricing_model | String | 否 | 定价模式筛选(FREE/PAID) |
| sort_by | String | 否 | 排序字段(popular/newest/top_rated) |
| language | String | 否 | 语言筛选 |
| featured_only | Boolean | 否 | 仅显示精选(默认false) |

**请求示例**:

```http
GET /api/v1/bot-store/listings?page=1&page_size=20&category=客服支持&pricing_model=FREE&sort_by=popular
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "listing_id": 123,
      "bot_id": 789,
      "title": "智能客服助手",
      "description": "7x24小时智能客服,支持多轮对话、知识库问答",
      "avatar": "https://cdn.coze.com/bot/789/avatar.png",
      "developer": {
        "id": 101,
        "name": "张三",
        "avatar": "https://cdn.coze.com/user/101/avatar.png"
      },
      "category": "客服支持",
      "subcategory": "在线客服",
      "tags": ["客服", "问答", "AI"],
      "pricing_model": "FREE",
      "rating_avg": 4.7,
      "rating_count": 523,
      "clone_count": 8520,
      "view_count": 25000,
      "last_updated": "2025-12-28T15:30:00Z",
      "published_at": "2025-12-01T10:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 125,
      "total_pages": 7
    }
  }
}
```

### 2.3 获取商店条目详情

**接口描述**: 获取商店条目的完整详细信息

**请求方式**: `GET /api/v1/bot-store/listings/{listing_id}`

**权限要求**: 无需认证(登录用户可看到额外信息)

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "listing_id": 123,
    "bot_id": 789,
    "title": "智能客服助手",
    "description": "7x24小时智能客服,支持多轮对话、知识库问答",
    "long_description": "## 功能介绍\n\n这是一个基于AI的智能客服助手...",
    "avatar": "https://cdn.coze.com/bot/789/avatar.png",
    "developer": {
      "id": 101,
      "name": "张三",
      "avatar": "https://cdn.coze.com/user/101/avatar.png",
      "verified": true,
      "followers_count": 1250
    },
    "category": "客服支持",
    "subcategory": "在线客服",
    "tags": ["客服", "问答", "AI", "多轮对话"],
    "pricing_model": "FREE",
    "pricing": {
      "amount": 0,
      "currency": "CNY"
    },
    "screenshots": [
      "https://cdn.coze.com/screenshots/789/1.png",
      "https://cdn.coze.com/screenshots/789/2.png",
      "https://cdn.coze.com/screenshots/789/3.png"
    ],
    "video_url": "https://cdn.coze.com/videos/789/demo.mp4",
    "language": "zh-CN",
    "features": [
      "智能问答",
      "多轮对话",
      "知识库集成",
      "情感分析"
    ],
    "use_cases": [
      "电商客服",
      "产品咨询",
      "售后支持"
    ],
    "version": "1.0.0",
    "changelog": "初始版本发布",
    "stats": {
      "rating_avg": 4.7,
      "rating_count": 523,
      "clone_count": 8520,
      "view_count": 25000,
      "favorite_count": 1520
    },
    "rating_distribution": {
      "5_star": 420,
      "4_star": 85,
      "3_star": 15,
      "2_star": 2,
      "1_star": 1
    },
    "user_interaction": {
      "is_cloned": false,
      "is_favorited": false,
      "user_rating": null
    },
    "status": "PUBLISHED",
    "published_at": "2025-12-01T10:00:00Z",
    "last_updated": "2025-12-28T15:30:00Z",
    "created_at": "2025-11-25T10:00:00Z"
  }
}
```

### 2.4 获取分类列表

**接口描述**: 获取所有Bot分类及每个分类下的Bot数量

**请求方式**: `GET /api/v1/bot-store/categories`

**权限要求**: 无需认证

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "客服支持",
      "icon": "https://cdn.coze.com/categories/cs.png",
      "description": "智能客服、在线客服、售后支持",
      "count": 1250,
      "subcategories": [
        {
          "id": 11,
          "name": "在线客服",
          "count": 850
        },
        {
          "id": 12,
          "name": "售后支持",
          "count": 400
        }
      ]
    },
    {
      "id": 2,
      "name": "生产力工具",
      "icon": "https://cdn.coze.com/categories/prod.png",
      "description": "写作助手、数据分析、项目管理",
      "count": 890,
      "subcategories": []
    }
  ]
}
```

---

## 3. Bot搜索API

### 3.1 搜索Bot

**接口描述**: 根据关键词搜索Bot

**请求方式**: `GET /api/v1/bot-store/search`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| q | String | ✅ | 搜索关键词 |
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| category | String | 否 | 分类筛选 |
| pricing_model | String | 否 | 定价模式筛选 |
| language | String | 否 | 语言筛选 |
| rating_min | Float | 否 | 最低评分筛选 |

**请求示例**:

```http
GET /api/v1/bot-store/search?q=客服机器人&category=客服支持&rating_min=4.0&page=1&page_size=20
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "listing_id": 123,
        "bot_id": 789,
        "title": "智能客服助手",
        "description": "7x24小时智能客服,支持多轮对话、知识库问答",
        "avatar": "https://cdn.coze.com/bot/789/avatar.png",
        "developer": {
          "id": 101,
          "name": "张三"
        },
        "category": "客服支持",
        "tags": ["客服", "问答", "AI"],
        "pricing_model": "FREE",
        "rating_avg": 4.7,
        "clone_count": 8520,
        "highlight": {
          "title": "智能<em>客服</em>助手",
          "description": "7x24小时智能<em>客服</em>,支持多轮对话、知识库问答"
        }
      }
    ],
    "facets": {
      "categories": [
        {
          "value": "客服支持",
          "count": 20
        },
        {
          "value": "生产力工具",
          "count": 5
        }
      ],
      "pricing_models": [
        {
          "value": "FREE",
          "count": 18
        },
        {
          "value": "PAID",
          "count": 7
        }
      ],
      "ratings": [
        {
          "range": "4.0+",
          "count": 22
        },
        {
          "range": "3.0-4.0",
          "count": 3
        }
      ]
    }
  }
}
```

### 3.2 获取搜索建议

**接口描述**: 获取搜索关键词的自动补全建议

**请求方式**: `GET /api/v1/bot-store/search/suggestions`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| q | String | ✅ | 输入的关键词 |
| limit | Integer | 否 | 返回数量(默认10) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "query": "客服",
    "suggestions": [
      {
        "text": "客服机器人",
        "type": "keyword"
      },
      {
        "text": "智能客服助手",
        "type": "bot",
        "bot_id": 789,
        "avatar": "https://cdn.coze.com/bot/789/avatar.png"
      },
      {
        "text": "在线客服系统",
        "type": "keyword"
      }
    ]
  }
}
```

---

## 4. Bot一键复克API

### 4.1 一键复克Bot

**接口描述**: 用户一键复克商店中的Bot

**请求方式**: `POST /api/v1/bot-store/listings/{listing_id}/clone`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "new_bot_name": "我的客服助手",
    "include_knowledge_base": true,
    "include_workflow": true
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| new_bot_name | String | 否 | 新Bot名称(默认"源Bot名 副本") |
| include_knowledge_base | Boolean | 否 | 是否包含知识库(默认true) |
| include_workflow | Boolean | 否 | 是否包含工作流(默认true) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Bot cloned successfully",
  "data": {
    "new_bot_id": 890,
    "new_bot_name": "我的客服助手",
    "source_listing": {
      "listing_id": 123,
      "title": "智能客服助手",
      "version": "1.0.0"
    },
    "clone_summary": {
      "workflow_cloned": 3,
      "knowledge_base_cloned": 2,
      "total_documents": 156,
      "clone_time_seconds": 3.5
    },
    "created_at": "2025-12-30T14:00:00Z"
  }
}
```

### 4.2 获取复克记录

**接口描述**: 获取用户的Bot复克记录列表

**请求方式**: `GET /api/v1/bot-store/clones`

**权限要求**: 需要认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| listing_id | Long | 否 | 筛选指定条目 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "clone_id": 456,
      "source_listing": {
        "listing_id": 123,
        "title": "智能客服助手",
        "avatar": "https://cdn.coze.com/bot/789/avatar.png"
      },
      "new_bot": {
        "bot_id": 890,
        "name": "我的客服助手"
      },
      "cloned_at": "2025-12-30T14:00:00Z"
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

## 5. Bot评价系统API

### 5.1 创建Bot评价

**接口描述**: 用户对Bot进行评分和评论

**请求方式**: `POST /api/v1/bot-store/listings/{listing_id}/reviews`

**权限要求**: 需要认证,必须已复克该Bot

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "rating": 5,
    "title": "非常好用的客服助手",
    "content": "功能强大,设置简单,帮我们节省了大量人力成本",
    "pros": [
      "配置简单",
      "响应快速",
      "知识库集成方便"
    ],
    "cons": [
      "语音功能稍弱"
    ],
    "tags": ["易用", "推荐"]
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| rating | Integer | ✅ | 评分(1-5) |
| title | String | ✅ | 评价标题 |
| content | String | ✅ | 评价内容 |
| pros | Array | 否 | 优点列表 |
| cons | Array | 否 | 缺点列表 |
| tags | Array | 否 | 标签列表 |

**响应示例**:

```json
{
  "code": 0,
  "message": "Review created successfully",
  "data": {
    "review_id": 789,
    "listing_id": 123,
    "rating": 5,
    "title": "非常好用的客服助手",
    "content": "功能强大,设置简单,帮我们节省了大量人力成本",
    "user": {
      "id": 201,
      "name": "李四",
      "avatar": "https://cdn.coze.com/user/201/avatar.png"
    },
    "pros": [
      "配置简单",
      "响应快速",
      "知识库集成方便"
    ],
    "cons": [
      "语音功能稍弱"
    ],
    "tags": ["易用", "推荐"],
    "helpful_count": 0,
    "created_at": "2025-12-30T15:00:00Z"
  }
}
```

### 5.2 获取Bot评价列表

**接口描述**: 获取Bot的评价列表

**请求方式**: `GET /api/v1/bot-store/listings/{listing_id}/reviews`

**权限要求**: 无需认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认10) |
| sort_by | String | 否 | 排序方式(recent/helpful/highest_rated) |
| rating | Integer | 否 | 筛选评分(1-5) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "summary": {
      "rating_avg": 4.7,
      "rating_count": 523,
      "rating_distribution": {
        "5_star": 420,
        "4_star": 85,
        "3_star": 15,
        "2_star": 2,
        "1_star": 1
      }
    },
    "reviews": [
      {
        "review_id": 789,
        "listing_id": 123,
        "rating": 5,
        "title": "非常好用的客服助手",
        "content": "功能强大,设置简单,帮我们节省了大量人力成本",
        "user": {
          "id": 201,
          "name": "李四",
          "avatar": "https://cdn.coze.com/user/201/avatar.png"
        },
        "pros": ["配置简单", "响应快速"],
        "cons": ["语音功能稍弱"],
        "tags": ["易用", "推荐"],
        "helpful_count": 25,
        "is_helpful": false,
        "created_at": "2025-12-30T15:00:00Z"
      }
    ],
    "meta": {
      "pagination": {
        "page": 1,
        "page_size": 10,
        "total_count": 523,
        "total_pages": 53
      }
    }
  }
}
```

### 5.3 标记评价为有用

**接口描述**: 用户标记某个评价为有用

**请求方式**: `POST /api/v1/bot-store/reviews/{review_id}/helpful`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| review_id | Long | 评价ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Review marked as helpful",
  "data": {
    "review_id": 789,
    "helpful_count": 26,
    "is_helpful": true
  }
}
```

### 5.4 回复评价(开发者)

**接口描述**: Bot开发者回复用户评价

**请求方式**: `POST /api/v1/bot-store/reviews/{review_id}/reply`

**权限要求**: Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| review_id | Long | 评价ID |

**请求体**:

```json
{
  "data": {
    "content": "感谢您的好评!我们正在优化语音功能,下个版本会改进。"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Reply added",
  "data": {
    "reply_id": 101,
    "review_id": 789,
    "content": "感谢您的好评!我们正在优化语音功能,下个版本会改进。",
    "replier": {
      "id": 101,
      "name": "张三",
      "role": "developer"
    },
    "created_at": "2025-12-30T16:00:00Z"
  }
}
```

---

## 6. Bot排行榜API

### 6.1 获取热门Bot排行榜

**接口描述**: 获取热门Bot排行榜(按复克次数排序)

**请求方式**: `GET /api/v1/bot-store/leaderboards/popular`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | String | 否 | 统计周期(week/month/all,默认all) |
| limit | Integer | 否 | 返回数量(默认50) |
| category | String | 否 | 分类筛选 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "week",
    "updated_at": "2025-12-30T17:00:00Z",
    "leaderboard": [
      {
        "rank": 1,
        "listing_id": 101,
        "title": "智能写作助手",
        "avatar": "https://cdn.coze.com/bot/801/avatar.png",
        "developer": "Coze官方",
        "category": "生产力工具",
        "clone_count": 5200,
        "rating_avg": 4.8,
        "trend": "up"
      },
      {
        "rank": 2,
        "listing_id": 123,
        "title": "智能客服助手",
        "avatar": "https://cdn.coze.com/bot/789/avatar.png",
        "developer": "张三",
        "category": "客服支持",
        "clone_count": 4850,
        "rating_avg": 4.7,
        "trend": "same"
      }
    ]
  }
}
```

### 6.2 获取新星Bot排行榜

**接口描述**: 获取新星Bot排行榜(最近发布但增长快)

**请求方式**: `GET /api/v1/bot-store/leaderboards/rising`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | Integer | 否 | 统计天数(默认7天) |
| limit | Integer | 否 | 返回数量(默认50) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period_days": 7,
    "updated_at": "2025-12-30T17:00:00Z",
    "leaderboard": [
      {
        "rank": 1,
        "listing_id": 201,
        "title": "代码审查助手",
        "avatar": "https://cdn.coze.com/bot/901/avatar.png",
        "developer": "李四",
        "published_at": "2025-12-20T10:00:00Z",
        "clone_count": 3200,
        "clone_growth_rate": 320.0,
        "rating_avg": 4.9
      }
    ]
  }
}
```

### 6.3 获取高分Bot排行榜

**接口描述**: 获取高分Bot排行榜(按评分排序)

**请求方式**: `GET /api/v1/bot-store/leaderboards/top-rated`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| min_rating_count | Integer | 否 | 最少评价数(默认50) |
| limit | Integer | 否 | 返回数量(默认50) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "updated_at": "2025-12-30T17:00:00Z",
    "leaderboard": [
      {
        "rank": 1,
        "listing_id": 301,
        "title": "数据分析专家",
        "avatar": "https://cdn.coze.com/bot/1001/avatar.png",
        "developer": "王五",
        "rating_avg": 5.0,
        "rating_count": 89,
        "clone_count": 2100
      }
    ]
  }
}
```

---

## 7. Bot模板库API

### 7.1 获取模板列表

**接口描述**: 获取Bot模板列表

**请求方式**: `GET /api/v1/bot-templates`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category | String | 否 | 分类筛选 |
| source | String | 否 | 来源筛选(official/community) |
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "template_id": 1,
      "name": "客服机器人模板",
      "description": "快速搭建智能客服机器人",
      "category": "客服支持",
      "source": "official",
      "thumbnail": "https://cdn.coze.com/templates/1/thumb.png",
      "preview_url": "https://cdn.coze.com/templates/1/preview.png",
      "features": [
        "多轮对话",
        "知识库集成",
        "人工转接"
      ],
      "use_count": 15000,
      "rating_avg": 4.8,
      "difficulty": "EASY",
      "estimated_time": "10分钟"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 45,
      "total_pages": 3
    }
  }
}
```

### 7.2 获取模板详情

**接口描述**: 获取模板的详细信息

**请求方式**: `GET /api/v1/bot-templates/{template_id}`

**权限要求**: 无需认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| template_id | Long | 模板ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "template_id": 1,
    "name": "客服机器人模板",
    "description": "快速搭建智能客服机器人",
    "long_description": "## 模板介绍\n\n这是一个专为客服场景设计的Bot模板...",
    "category": "客服支持",
    "source": "official",
    "thumbnail": "https://cdn.coze.com/templates/1/thumb.png",
    "screenshots": [
      "https://cdn.coze.com/templates/1/s1.png",
      "https://cdn.coze.com/templates/1/s2.png"
    ],
    "features": [
      "多轮对话",
      "知识库集成",
      "人工转接",
      "情感分析"
    ],
    "components": [
      {
        "type": "workflow",
        "name": "客服工作流",
        "description": "完整的客服对话流程"
      },
      {
        "type": "knowledge_base",
        "name": "产品知识库",
        "description": "预置产品FAQ知识库"
      }
    ],
    "configuration_guide": [
      {
        "step": 1,
        "title": "配置知识库",
        "description": "上传您的产品FAQ文档"
      },
      {
        "step": 2,
        "title": "调整对话流程",
        "description": "根据业务场景修改工作流"
      }
    ],
    "use_count": 15000,
    "rating_avg": 4.8,
    "rating_count": 320,
    "difficulty": "EASY",
    "estimated_time": "10分钟",
    "author": {
      "id": 1,
      "name": "Coze官方"
    },
    "version": "1.2.0",
    "updated_at": "2025-12-20T10:00:00Z"
  }
}
```

### 7.3 使用模板创建Bot

**接口描述**: 使用模板快速创建Bot

**请求方式**: `POST /api/v1/bot-templates/{template_id}/use`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| template_id | Long | 模板ID |

**请求体**:

```json
{
  "data": {
    "bot_name": "我的客服机器人"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Bot created from template",
  "data": {
    "bot_id": 901,
    "bot_name": "我的客服机器人",
    "template": {
      "template_id": 1,
      "name": "客服机器人模板",
      "version": "1.2.0"
    },
    "components_included": [
      {
        "type": "workflow",
        "name": "客服工作流",
        "id": 101
      },
      {
        "type": "knowledge_base",
        "name": "产品知识库",
        "id": 51
      }
    ],
    "next_steps": [
      {
        "title": "配置知识库",
        "description": "上传您的产品FAQ文档",
        "action_url": "/bots/901/knowledge"
      },
      {
        "title": "调整对话流程",
        "description": "根据业务场景修改工作流",
        "action_url": "/bots/901/workflow"
      }
    ],
    "created_at": "2025-12-30T18:00:00Z"
  }
}
```

---

## 8. 审核管理API

### 8.1 获取待审核列表(管理员)

**接口描述**: 管理员获取待审核的商店条目列表

**请求方式**: `GET /api/v1/admin/bot-store/pending`

**权限要求**: `bot_store:review`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| category | String | 否 | 分类筛选 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "listing_id": 123,
      "bot_id": 789,
      "title": "智能客服助手",
      "description": "7x24小时智能客服",
      "avatar": "https://cdn.coze.com/bot/789/avatar.png",
      "developer": {
        "id": 101,
        "name": "张三",
        "email": "zhangsan@example.com"
      },
      "category": "客服支持",
      "tags": ["客服", "问答"],
      "pricing_model": "FREE",
      "screenshots": [
        "https://cdn.coze.com/screenshots/789/1.png"
      ],
      "submitted_at": "2025-12-29T10:00:00Z",
      "waiting_time_hours": 24
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

### 8.2 审核通过(管理员)

**接口描述**: 管理员审核通过商店条目

**请求方式**: `POST /api/v1/admin/bot-store/listings/{listing_id}/approve`

**权限要求**: `bot_store:review`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "notes": "审核通过,内容符合规范",
    "featured": true,
    "featured_until": "2025-12-31T23:59:59Z"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Listing approved",
  "data": {
    "listing_id": 123,
    "status": "PUBLISHED",
    "approved_at": "2025-12-30T19:00:00Z",
    "reviewer": {
      "id": 1,
      "name": "Admin"
    },
    "notes": "审核通过,内容符合规范",
    "featured": true,
    "featured_until": "2025-12-31T23:59:59Z"
  }
}
```

### 8.3 审核拒绝(管理员)

**接口描述**: 管理员审核拒绝商店条目

**请求方式**: `POST /api/v1/admin/bot-store/listings/{listing_id}/reject`

**权限要求**: `bot_store:review`

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| listing_id | Long | 条目ID |

**请求体**:

```json
{
  "data": {
    "reason": "描述内容与实际功能不符",
    "details": "实际测试后发现Bot无法实现描述中的语音交互功能"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Listing rejected",
  "data": {
    "listing_id": 123,
    "status": "REJECTED",
    "rejected_at": "2025-12-30T19:00:00Z",
    "reviewer": {
      "id": 1,
      "name": "Admin"
    },
    "reason": "描述内容与实际功能不符",
    "details": "实际测试后发现Bot无法实现描述中的语音交互功能",
    "can_resubmit": true
  }
}
```

### 8.4 获取审核历史(管理员)

**接口描述**: 获取商店条目的审核历史记录

**请求方式**: `GET /api/v1/admin/bot-store/listings/{listing_id}/review-history`

**权限要求**: `bot_store:review`

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "review_id": 501,
      "action": "APPROVE",
      "reviewer": {
        "id": 1,
        "name": "Admin"
      },
      "notes": "审核通过",
      "created_at": "2025-12-30T19:00:00Z"
    },
    {
      "review_id": 502,
      "action": "REJECT",
      "reviewer": {
        "id": 2,
        "name": "Admin2"
      },
      "reason": "描述不符",
      "created_at": "2025-12-28T10:00:00Z"
    }
  ]
}
```

---

## 附录

### A. 枚举类型定义

```typescript
// Bot状态
enum ListingStatus {
  DRAFT = 'DRAFT',                    // 草稿
  PENDING_REVIEW = 'PENDING_REVIEW',  // 待审核
  PUBLISHED = 'PUBLISHED',            // 已发布
  REJECTED = 'REJECTED',              // 已拒绝
  UNPUBLISHED = 'UNPUBLISHED'         // 已下架
}

// 定价模式
enum PricingModel {
  FREE = 'FREE',           // 免费
  PAID = 'PAID'            // 付费
}

// 模板来源
enum TemplateSource {
  OFFICIAL = 'official',   // 官方
  COMMUNITY = 'community'  // 社区
}

// 难度等级
enum Difficulty {
  EASY = 'EASY',           // 简单
  MEDIUM = 'MEDIUM',       // 中等
  HARD = 'HARD'            // 复杂
}

// 排序方式
enum SortBy {
  POPULAR = 'popular',           // 热门
  NEWEST = 'newest',             // 最新
  TOP_RATED = 'top_rated',       // 高分
  CLONE_COUNT = 'clone_count'    // 复克次数
}
```

### B. Go后端实现示例

```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// BotStoreHandler Bot商店处理器
type BotStoreHandler struct {
    storeService     BotStoreService
    cloneService     CloneService
    reviewService    ReviewService
    templateService  TemplateService
}

// CreateListing 创建商店条目
func (h *BotStoreHandler) CreateListing(ctx context.Context, c *app.RequestContext) {
    var req CreateListingRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    userID := getUserIDFromContext(ctx)

    listing, err := h.storeService.CreateListing(ctx, userID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(201, BotResponse{
        Code:    0,
        Message: "Listing created successfully",
        Data:    listing,
    })
}

// CloneBot 一键复克Bot
func (h *BotStoreHandler) CloneBot(ctx context.Context, c *app.RequestContext) {
    listingID := c.Param("listing_id")
    var req CloneBotRequest
    c.BindAndValidate(&req)

    userID := getUserIDFromContext(ctx)

    result, err := h.cloneService.CloneBot(ctx, userID, listingID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "Bot cloned successfully",
        Data:    result,
    })
}

// ApproveListing 审核通过(管理员)
func (h *BotStoreHandler) ApproveListing(ctx context.Context, c *app.RequestContext) {
    listingID := c.Param("listing_id")
    var req ApproveRequest
    c.BindAndValidate(&req)

    result, err := h.reviewService.ApproveListing(ctx, listingID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "Listing approved",
        Data:    result,
    })
}
```

### C. 前端TypeScript类型定义

```typescript
// types/bot-store.ts

// 商店条目
interface BotStoreListing {
  listing_id: number;
  bot_id: number;
  title: string;
  description: string;
  avatar: string;
  developer: {
    id: number;
    name: string;
    avatar: string;
    verified?: boolean;
  };
  category: string;
  tags: string[];
  pricing_model: 'FREE' | 'PAID';
  rating_avg: number;
  rating_count: number;
  clone_count: number;
  view_count: number;
  status: ListingStatus;
}

// Bot评价
interface BotReview {
  review_id: number;
  listing_id: number;
  rating: number;
  title: string;
  content: string;
  user: {
    id: number;
    name: string;
    avatar: string;
  };
  pros: string[];
  cons: string[];
  tags: string[];
  helpful_count: number;
  created_at: string;
}

// Bot模板
interface BotTemplate {
  template_id: number;
  name: string;
  description: string;
  category: string;
  source: 'official' | 'community';
  thumbnail: string;
  features: string[];
  use_count: number;
  rating_avg: number;
  difficulty: 'EASY' | 'MEDIUM' | 'HARD';
  estimated_time: string;
}

// API客户端
class BotStoreApiClient {
  async createListing(data: CreateListingRequest): Promise<BotStoreListing> {
    const response = await apiClient.post('/bot-store/listings', { data });
    return response.data.data;
  }

  async getListings(params?: GetListingsParams): Promise<BotStoreListing[]> {
    const response = await apiClient.get('/bot-store/listings', { params });
    return response.data.data;
  }

  async cloneBot(listingId: number, data: CloneBotRequest): Promise<CloneResult> {
    const response = await apiClient.post(`/bot-store/listings/${listingId}/clone`, {
      data,
    });
    return response.data.data;
  }

  async createReview(
    listingId: number,
    data: CreateReviewRequest
  ): Promise<BotReview> {
    const response = await apiClient.post(`/bot-store/listings/${listingId}/reviews`, {
      data,
    });
    return response.data.data;
  }

  async useTemplate(
    templateId: number,
    botName: string
  ): Promise<CreateBotFromTemplateResult> {
    const response = await apiClient.post(`/bot-templates/${templateId}/use`, {
      data: { bot_name: botName },
    });
    return response.data.data;
  }
}

export const botStoreApi = new BotStoreApiClient();
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,包含Bot商店发布、浏览、搜索、复克、评价、排行榜、模板库、审核管理等8个模块的完整API接口定义
