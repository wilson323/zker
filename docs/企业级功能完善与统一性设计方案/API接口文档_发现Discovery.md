# API接口文档：发现Discovery模块

**模块名称**: 发现Discovery (Discovery)
**设计文档**: 04-ZKER前台_发现_Discovery.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 发现页数据API](#2-发现页数据api)
- [3. Bot搜索API](#3-bot搜索api)
- [4. 分类管理API](#4-分类管理api)
- [5. Bot操作API](#5-bot操作api)
- [6. 配置管理API](#6-配置管理api)
- [7. 数据模型](#7-数据模型)
- [8. 后端代码示例](#8-后端代码示例)
- [9. 前端代码示例](#9-前端代码示例)
- [10. 错误码定义](#10-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

发现Discovery是ZKER企业级SaaS平台的前台Bot展示模块，通过数据库驱动和通用展示组件提供Bot发现、浏览和快速访问能力：

- ✅ **Bot展示** - 卡片式展示Bot名称、描述、图标、使用量
- ✅ **分类浏览** - 按办公、营销、客服等分类筛选Bot
- ✅ **搜索Bot** - 按名称、描述关键词搜索
- ✅ **推荐Bot** - 首页推荐位展示精选Bot
- ✅ **快速访问** - 一键启动Bot对话
- ✅ **灵活配置** - 管理员可配置展示Bot、分类、排序规则

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0 (缓存热门Bot)

**前端技术栈**:
- 框架: React 18 + TypeScript
- UI库: Semi Design
- 路由: React Router v6

### 1.3 核心设计理念

| 设计原则 | 实现方式 | 占比 |
|---------|---------|------|
| **配置驱动** | 数据库配置展示内容和顺序 | 95% |
| **通用组件** | 可复用的Bot卡片、分类筛选组件 | 5% |

---

## 2. 发现页数据API

### 2.1 获取发现页数据

**接口地址**: `GET /api/v1/discovery/bots`

**功能说明**: 获取发现页的完整数据（分类、推荐Bot、所有Bot）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:
- category: 分类ID（可选）
- featured: 是否只返回推荐Bot (true/false，默认false)
- sort: 排序方式 (featured/recent/popular，默认featured)
- limit: 每页数量，默认20
- offset: 偏移量，默认0

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "categories": [
      {
        "id": 1,
        "name": "办公效率",
        "icon": "💼",
        "color": "#1890ff",
        "botCount": 15
      },
      {
        "id": 2,
        "name": "营销推广",
        "icon": "📢",
        "color": "#52c41a",
        "botCount": 8
      },
      {
        "id": 3,
        "name": "客户服务",
        "icon": "🎧",
        "color": "#faad14",
        "botCount": 12
      }
    ],
    "featured": [
      {
        "id": "bot-001",
        "name": "邮件助手",
        "description": "智能撰写和回复邮件",
        "avatar": "https://cdn.example.com/avatars/bot-email.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 1250,
        "rating": 4.8,
        "isFeatured": true,
        "featuredPosition": 1
      },
      {
        "id": "bot-002",
        "name": "周报生成器",
        "description": "快速生成工作周报",
        "avatar": "https://cdn.example.com/avatars/bot-report.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 890,
        "rating": 4.6,
        "isFeatured": true,
        "featuredPosition": 2
      }
    ],
    "bots": [
      {
        "id": "bot-003",
        "name": "日程管理",
        "description": "智能安排日程和提醒",
        "avatar": "https://cdn.example.com/avatars/bot-calendar.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 567,
        "rating": 4.5,
        "isFeatured": false
      }
    ],
    "total": 150
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "trace-abc-123"
}
```

### 2.2 获取推荐Bot

**接口地址**: `GET /api/v1/discovery/featured`

**功能说明**: 仅获取首页推荐Bot列表

**查询参数**:
- limit: 返回数量，默认10，最大20

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "bot-001",
        "name": "邮件助手",
        "description": "智能撰写和回复邮件",
        "avatar": "https://cdn.example.com/avatars/bot-email.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 1250,
        "rating": 4.8,
        "featuredPosition": 1,
        "featuredExpiresAt": "2025-01-31T23:59:59Z"
      }
    ],
    "total": 10
  }
}
```

### 2.3 获取Bot详情

**接口地址**: `GET /api/v1/discovery/bots/{bot_id}`

**功能说明**: 获取Bot详细信息

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "bot-001",
    "name": "邮件助手",
    "description": "智能撰写和回复邮件，提升工作效率",
    "longDescription": "邮件助手是一款基于AI的智能邮件处理工具...",
    "avatar": "https://cdn.example.com/avatars/bot-email.png",
    "category": {
      "id": 1,
      "name": "办公效率",
      "icon": "💼",
      "color": "#1890ff"
    },
    "capabilities": [
      "撰写邮件",
      "回复邮件",
      "邮件润色",
      "邮件翻译"
    ],
    "usageCount": 1250,
    "rating": 4.8,
    "ratingCount": 85,
    "isFeatured": true,
    "viewCount": 3500,
    "clickCount": 1200,
    "createdAt": "2025-01-01T10:00:00Z",
    "updatedAt": "2025-01-03T10:00:00Z"
  }
}
```

---

## 3. Bot搜索API

### 3.1 搜索Bot

**接口地址**: `GET /api/v1/discovery/search`

**功能说明**: 按关键词搜索Bot（支持名称、描述搜索）

**查询参数**:
- q: 搜索关键词（必填）
- category: 分类ID过滤（可选）
- sort: 排序方式 (relevance/popular/recent，默认relevance)
- limit: 返回数量，默认20
- offset: 偏移量，默认0

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "results": [
      {
        "id": "bot-001",
        "name": "邮件助手",
        "description": "智能撰写和回复邮件",
        "avatar": "https://cdn.example.com/avatars/bot-email.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 1250,
        "rating": 4.8,
        "highlight": {
          "name": "<em>邮件</em>助手",
          "description": "智能撰写和回复<em>邮件</em>"
        }
      },
      {
        "id": "bot-004",
        "name": "邮件模板库",
        "description": "常用邮件模板集合",
        "avatar": "https://cdn.example.com/avatars/bot-email-templates.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 623,
        "rating": 4.5,
        "highlight": {
          "name": "<em>邮件</em>模板库",
          "description": "常用<em>邮件</em>模板集合"
        }
      }
    ],
    "total": 5,
    "query": "邮件"
  }
}
```

### 3.2 搜索建议

**接口地址**: `GET /api/v1/discovery/search/suggestions`

**功能说明**: 获取搜索关键词建议（自动补全）

**查询参数**:
- q: 输入的前缀（必填）
- limit: 返回数量，默认10

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "suggestions": [
      "邮件助手",
      "邮件模板",
      "邮件写作",
      "邮件翻译",
      "营销邮件"
    ]
  }
}
```

---

## 4. 分类管理API

### 4.1 获取分类列表

**接口地址**: `GET /api/v1/discovery/categories`

**功能说明**: 获取所有分类（包含Bot数量统计）

**查询参数**:
- include_count: 是否包含Bot数量统计 (默认true)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "办公效率",
        "description": "提升办公效率的Bot",
        "icon": "💼",
        "color": "#1890ff",
        "sortOrder": 1,
        "isActive": true,
        "botCount": 15
      },
      {
        "id": 2,
        "name": "营销推广",
        "description": "营销与推广Bot",
        "icon": "📢",
        "color": "#52c41a",
        "sortOrder": 2,
        "isActive": true,
        "botCount": 8
      }
    ],
    "total": 5
  }
}
```

### 4.2 按分类获取Bot

**接口地址**: `GET /api/v1/discovery/categories/{category_id}/bots`

**功能说明**: 获取指定分类下的Bot列表

**查询参数**:
- sort: 排序方式 (featured/recent/popular，默认featured)
- limit: 每页数量，默认20
- offset: 偏移量，默认0

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "category": {
      "id": 1,
      "name": "办公效率",
      "icon": "💼",
      "color": "#1890ff"
    },
    "items": [
      {
        "id": "bot-001",
        "name": "邮件助手",
        "description": "智能撰写和回复邮件",
        "avatar": "https://cdn.example.com/avatars/bot-email.png",
        "usageCount": 1250,
        "rating": 4.8
      }
    ],
    "total": 15
  }
}
```

### 4.3 创建分类（管理员）

**接口地址**: `POST /api/v1/discovery/categories`

**功能说明**: 创建新分类（仅管理员）

**请求参数**:
```json
{
  "name": "销售支持",
  "description": "销售相关Bot",
  "icon": "🎯",
  "color": "#f5222d",
  "sort_order": 10
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "分类创建成功",
  "data": {
    "id": 6,
    "name": "销售支持",
    "description": "销售相关Bot",
    "icon": "🎯",
    "color": "#f5222d",
    "sortOrder": 10,
    "isActive": true,
    "createdAt": "2025-01-03T10:00:00Z"
  }
}
```

### 4.4 更新分类（管理员）

**接口地址**: `PUT /api/v1/discovery/categories/{category_id}`

### 4.5 删除分类（管理员）

**接口地址**: `DELETE /api/v1/discovery/categories/{category_id}`

---

## 5. Bot操作API

### 5.1 收藏Bot

**接口地址**: `POST /api/v1/discovery/bots/{bot_id}/favorite`

**功能说明**: 收藏Bot到个人收藏列表

**响应示例**:
```json
{
  "code": 0,
  "message": "收藏成功"
}
```

### 5.2 取消收藏

**接口地址**: `DELETE /api/v1/discovery/bots/{bot_id}/favorite`

### 5.3 获取收藏列表

**接口地址**: `GET /api/v1/discovery/favorites`

**查询参数**:
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": "bot-001",
        "name": "邮件助手",
        "description": "智能撰写和回复邮件",
        "avatar": "https://cdn.example.com/avatars/bot-email.png",
        "favoritedAt": "2025-01-03T10:00:00Z"
      }
    ],
    "total": 5
  }
}
```

### 5.4 记录Bot浏览

**接口地址**: `POST /api/v1/discovery/bots/{bot_id}/view`

**功能说明**: 记录Bot浏览次数（统计用）

**响应示例**:
```json
{
  "code": 0,
  "message": "success"
}
```

### 5.5 记录Bot点击

**接口地址**: `POST /api/v1/discovery/bots/{bot_id}/click`

**功能说明**: 记录Bot点击次数（统计用）

**响应示例**:
```json
{
  "code": 0,
  "message": "success"
}
```

---

## 6. 配置管理API

### 6.1 配置展示Bot（管理员）

**接口地址**: `POST /api/v1/discovery/items`

**功能说明**: 配置Bot显示在发现页（管理员）

**请求参数**:
```json
{
  "bot_id": "bot-001",
  "category_id": 1,
  "is_featured": true,
  "featured_position": 1,
  "sort_order": 1
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "配置成功",
  "data": {
    "id": 1,
    "botId": "bot-001",
    "categoryId": 1,
    "isFeatured": true,
    "featuredPosition": 1,
    "sortOrder": 1,
    "createdAt": "2025-01-03T10:00:00Z"
  }
}
```

### 6.2 更新展示配置（管理员）

**接口地址**: `PUT /api/v1/discovery/items/{item_id}`

### 6.3 移除展示Bot（管理员）

**接口地址**: `DELETE /api/v1/discovery/items/{item_id}`

### 6.4 批量配置推荐Bot（管理员）

**接口地址**: `POST /api/v1/discovery/featured/batch`

**功能说明**: 批量设置推荐Bot及其顺序

**请求参数**:
```json
{
  "items": [
    {
      "bot_id": "bot-001",
      "featured_position": 1
    },
    {
      "bot_id": "bot-002",
      "featured_position": 2
    },
    {
      "bot_id": "bot-003",
      "featured_position": 3
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "批量配置成功",
  "data": {
    "successCount": 3,
    "failedCount": 0
  }
}
```

---

## 7. 数据模型

### 7.1 DiscoveryCategory（分类）
```typescript
interface DiscoveryCategory {
  id: number;
  tenant_id?: string;
  parent_id?: number;
  name: string;
  description?: string;
  icon?: string;
  color?: string;
  sort_order: number;
  is_active: boolean;
  bot_count?: number;
  created_at: string;
  updated_at: string;
}
```

### 7.2 DiscoveryItem（展示配置）
```typescript
interface DiscoveryItem {
  id: number;
  tenant_id: string;
  bot_id: string;
  category_id?: number;

  // 展示配置
  is_featured: boolean;
  featured_position?: number;
  featured_expires_at?: string;

  sort_order: number;
  is_active: boolean;

  // 统计
  view_count: number;
  click_count: number;

  created_at: string;
  updated_at: string;
}
```

### 7.3 DiscoveryBot（Bot展示项）
```typescript
interface DiscoveryBot {
  id: string;
  tenant_id: string;
  name: string;
  description: string;
  avatar: string;
  category?: {
    id: number;
    name: string;
    icon?: string;
    color?: string;
  };
  usage_count: number;
  rating: number;
  is_featured?: boolean;
  featured_position?: number;
  highlight?: {
    name?: string;
    description?: string;
  };
}
```

---

## 8. 后端代码示例

### 8.1 DiscoveryService

```go
package service

import (
	"context"
	"fmt"
	"time"
)

// DiscoveryService 发现服务
type DiscoveryService struct {
	discoveryRepo repository.DiscoveryRepository
	botRepo       repository.BotRepository
	cache         *redis.Client
	logger        *zap.Logger
}

// GetDiscoveryPageData 获取发现页数据
func (s *DiscoveryService) GetDiscoveryPageData(
	ctx context.Context,
	req *GetDiscoveryPageRequest,
) (*GetDiscoveryPageResponse, error) {
	tenantID := getTenantID(ctx)

	// 1. 获取分类
	categories, err := s.discoveryRepo.ListCategories(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("获取分类失败: %w", err)
	}

	// 2. 获取推荐Bot（缓存）
	var featuredBots []*DiscoveryBotItem
	if req.Featured {
		featuredBots, err = s.getCachedFeaturedBots(ctx, tenantID)
		if err != nil {
			s.logger.Warn("获取缓存推荐Bot失败", zap.Error(err))
			// 降级查询数据库
			featuredBots, err = s.discoveryRepo.ListFeaturedBots(ctx, tenantID, 10)
			if err != nil {
				return nil, err
			}
		}
	}

	// 3. 获取所有展示Bot
	bots, err := s.discoveryRepo.ListBots(ctx, &ListDiscoveryBotsRequest{
		TenantID:   tenantID,
		CategoryID: req.CategoryID,
		Sort:       req.Sort,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("获取Bot列表失败: %w", err)
	}

	// 4. 补充Bot详细信息
	for _, bot := range bots {
		botDetail, err := s.botRepo.GetByID(ctx, bot.BotID)
		if err != nil {
			continue
		}
		bot.Name = botDetail.Name
		bot.Description = botDetail.Description
		bot.Avatar = botDetail.Avatar
		bot.UsageCount = botDetail.UsageCount
		bot.Rating = botDetail.AvgRating
	}

	return &GetDiscoveryPageResponse{
		Categories: categories,
		Featured:   featuredBots,
		Bots:       bots,
		Total:      len(bots),
	}, nil
}

// getCachedFeaturedBots 获取缓存的推荐Bot
func (s *DiscoveryService) getCachedFeaturedBots(
	ctx context.Context,
	tenantID string,
) ([]*DiscoveryBotItem, error) {
	cacheKey := fmt.Sprintf("discovery:featured:%s", tenantID)

	// 尝试从Redis获取
	data, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var bots []*DiscoveryBotItem
		if err := json.Unmarshal([]byte(data), &bots); err == nil {
			return bots, nil
		}
	}

	// 缓存未命中，查询数据库
	bots, err := s.discoveryRepo.ListFeaturedBots(ctx, tenantID, 10)
	if err != nil {
		return nil, err
	}

	// 写入缓存（5分钟）
	jsonData, _ := json.Marshal(bots)
	s.cache.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	return bots, nil
}

// SearchBots 搜索Bot
func (s *DiscoveryService) SearchBots(
	ctx context.Context,
	req *SearchBotsRequest,
) (*SearchBotsResponse, error) {
	tenantID := getTenantID(ctx)

	// 使用LIKE模糊搜索
	bots, total, err := s.discoveryRepo.SearchBots(ctx, &SearchBotsRepoRequest{
		TenantID:   tenantID,
		Query:      req.Query,
		CategoryID: req.CategoryID,
		Sort:       req.Sort,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("搜索Bot失败: %w", err)
	}

	// 高亮关键词
	for _, bot := range bots {
		bot.Highlight = &Highlight{
			Name:        highlightKeyword(bot.Name, req.Query),
			Description: highlightKeyword(bot.Description, req.Query),
		}
	}

	return &SearchBotsResponse{
		Results: bots,
		Total:   total,
		Query:   req.Query,
	}, nil
}

// highlightKeyword 高亮关键词
func highlightKeyword(text, keyword string) string {
	// 简单实现：用<em>标签包裹
	return strings.ReplaceAll(
		text,
		keyword,
		fmt.Sprintf("<em>%s</em>", keyword),
	)
}
```

### 8.2 Repository层

```go
package repository

import (
	"context"
	"database/sql"
)

type DiscoveryRepository interface {
	// ListCategories 列出分类
	ListCategories(ctx context.Context, tenantID string) ([]*entity.DiscoveryCategory, error)

	// ListBots 列出展示Bot
	ListBots(ctx context.Context, req *ListDiscoveryBotsRequest) ([]*DiscoveryBotItem, error)

	// ListFeaturedBots 列出推荐Bot
	ListFeaturedBots(ctx context.Context, tenantID string, limit int) ([]*DiscoveryBotItem, error)

	// SearchBots 搜索Bot
	SearchBots(ctx context.Context, req *SearchBotsRepoRequest) ([]*DiscoveryBotItem, int, error)
}

type discoveryRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewDiscoveryRepository(db *gorm.DB, rdb *redis.Client) DiscoveryRepository {
	return &discoveryRepository{db: db, rdb: rdb}
}

// ListCategories 实现示例
func (r *discoveryRepository) ListCategories(
	ctx context.Context,
	tenantID string,
) ([]*entity.DiscoveryCategory, error) {
	var categories []*entity.DiscoveryCategory

	err := r.db.WithContext(ctx).
		Table("discovery_categories").
		Where("tenant_id IS NULL OR tenant_id = ?", tenantID).
		Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	// 统计每个分类的Bot数量
	for _, cat := range categories {
		var count int64
		r.db.Table("discovery_items").
			Where("category_id = ?", cat.ID).
			Where("is_active = ?", true).
			Count(&count)
		cat.BotCount = int(count)
	}

	return categories, nil
}

// ListBots 实现示例
func (r *discoveryRepository) ListBots(
	ctx context.Context,
	req *ListDiscoveryBotsRequest,
) ([]*DiscoveryBotItem, error) {
	var items []*entity.DiscoveryItem

	query := r.db.WithContext(ctx).
		Table("discovery_items").
		Where("tenant_id = ?", req.TenantID).
		Where("is_active = ?", true)

	if req.CategoryID != nil {
		query = query.Where("category_id = ?", *req.CategoryID)
	}

	// 排序
	switch req.Sort {
	case "featured":
		query = query.Order("is_featured DESC, featured_position ASC, sort_order ASC")
	case "recent":
		query = query.Order("created_at DESC")
	case "popular":
		query = query.Order("usage_count DESC")
	default:
		query = query.Order("sort_order ASC")
	}

	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}
	if req.Offset > 0 {
		query = query.Offset(req.Offset)
	}

	err := query.Find(&items).Error
	return items, err
}
```

---

## 9. 前端代码示例

### 9.1 DiscoveryPage组件

```typescript
// src/modules/discovery/components/DiscoveryPage.tsx
import React, { useState, useEffect } from 'react';
import { Input, Empty, Spin, Toast } from '@douyinfe/semi-ui';
import { IconSearch } from '@douyinfe/semi-icons';
import { useNavigate } from 'react-router-dom';
import { discoveryAPI } from '@/services/api';
import { BotCard } from './BotCard';
import { CategoryFilter } from './CategoryFilter';
import { FeaturedSection } from './FeaturedSection';
import styles from './DiscoveryPage.module.scss';

export const DiscoveryPage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [categories, setCategories] = useState<Category[]>([]);
  const [featured, setFeatured] = useState<Bot[]>([]);
  const [bots, setBots] = useState<Bot[]>([]);
  const [activeCategory, setActiveCategory] = useState<number | null>(null);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadDiscoveryData();
  }, [activeCategory]);

  const loadDiscoveryData = async () => {
    setLoading(true);
    try {
      const resp = await discoveryAPI.getDiscoveryPageData({
        category: activeCategory,
        sort: 'featured',
        limit: 50,
      });
      setCategories(resp.data.categories);
      setFeatured(resp.data.featured);
      setBots(resp.data.bots);
    } catch (error) {
      Toast.error('加载数据失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    if (!query.trim()) {
      loadDiscoveryData();
      return;
    }

    try {
      const resp = await discoveryAPI.searchBots({
        q: query,
        category: activeCategory,
        limit: 20,
      });
      setBots(resp.data.results);
    } catch (error) {
      Toast.error('搜索失败');
    }
  };

  const handleBotClick = (botId: string) => {
    // 跳转到会话页
    navigate(`/chat?bot=${botId}`);

    // 记录点击统计
    discoveryAPI.recordBotClick(botId).catch(console.error);
  };

  if (loading) {
    return (
      <div className={styles.loading}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div className={styles.discoveryPage}>
      {/* 头部搜索 */}
      <div className={styles.header}>
        <h1>发现</h1>
        <Input
          prefix={<IconSearch />}
          placeholder="搜索Bot..."
          size="large"
          value={searchQuery}
          onChange={setSearchQuery}
          onEnterPress={(e) => handleSearch(e.target.value)}
          className={styles.searchBar}
        />
      </div>

      {/* 推荐区域 */}
      {featured.length > 0 && !searchQuery && (
        <FeaturedSection
          bots={featured}
          onBotClick={handleBotClick}
        />
      )}

      {/* 分类筛选 */}
      <CategoryFilter
        categories={categories}
        activeCategory={activeCategory}
        onSelectCategory={setActiveCategory}
      />

      {/* Bot列表 */}
      <div className={styles.botGrid}>
        {bots.map(bot => (
          <BotCard
            key={bot.id}
            bot={bot}
            onClick={() => handleBotClick(bot.id)}
          />
        ))}
      </div>

      {bots.length === 0 && (
        <Empty
          title="暂无Bot"
          description="换个关键词试试吧"
        />
      )}
    </div>
  );
};
```

### 9.2 BotCard组件

```typescript
// src/modules/discovery/components/BotCard.tsx
import React from 'react';
import { Card, Tag } from '@douyinfe/semi-ui';
import { IconStar, IconUsers } from '@douyinfe/semi-icons';
import styles from './BotCard.module.scss';

interface BotCardProps {
  bot: Bot;
  onClick: () => void;
}

export const BotCard: React.FC<BotCardProps> = ({ bot, onClick }) => {
  return (
    <Card
      className={styles.botCard}
      hoverable
      onClick={onClick}
      cover={<img src={bot.avatar} alt={bot.name} />}
    >
      <div className={styles.content}>
        <h3 className={styles.name}>{bot.name}</h3>
        <p className={styles.description}>{bot.description}</p>

        <div className={styles.meta}>
          {bot.category && (
            <Tag color="blue">{bot.category.name}</Tag>
          )}
          <span className={styles.stat}>
            <IconUsers /> {bot.usageCount}
          </span>
          <span className={styles.rating}>
            <IconStar /> {bot.rating}
          </span>
        </div>
      </div>

      {bot.isFeatured && (
        <div className={styles.featuredBadge}>推荐</div>
      )}
    </Card>
  );
};
```

### 9.3 API服务

```typescript
// src/services/api/discovery.ts
import { request } from '@/utils/request';

export const discoveryAPI = {
  // 获取发现页数据
  getDiscoveryPageData: (params?: {
    category?: number;
    featured?: boolean;
    sort?: string;
    limit?: number;
    offset?: number;
  }) => {
    return request.get('/api/v1/discovery/bots', { params });
  },

  // 获取推荐Bot
  getFeaturedBots: (params?: { limit?: number }) => {
    return request.get('/api/v1/discovery/featured', { params });
  },

  // 获取Bot详情
  getBotDetail: (botId: string) => {
    return request.get(`/api/v1/discovery/bots/${botId}`);
  },

  // 搜索Bot
  searchBots: (params: {
    q: string;
    category?: number;
    sort?: string;
    limit?: number;
    offset?: number;
  }) => {
    return request.get('/api/v1/discovery/search', { params });
  },

  // 获取分类列表
  getCategories: (params?: { include_count?: boolean }) => {
    return request.get('/api/v1/discovery/categories', { params });
  },

  // 按分类获取Bot
  getBotsByCategory: (
    categoryId: number,
    params?: { sort?: string; limit?: number; offset?: number }
  ) => {
    return request.get(`/api/v1/discovery/categories/${categoryId}/bots`, {
      params,
    });
  },

  // 收藏Bot
  favoriteBot: (botId: string) => {
    return request.post(`/api/v1/discovery/bots/${botId}/favorite`);
  },

  // 取消收藏
  unfavoriteBot: (botId: string) => {
    return request.delete(`/api/v1/discovery/bots/${botId}/favorite`);
  },

  // 获取收藏列表
  getFavorites: (params?: { page?: number; page_size?: number }) => {
    return request.get('/api/v1/discovery/favorites', { params });
  },

  // 记录浏览
  recordBotView: (botId: string) => {
    return request.post(`/api/v1/discovery/bots/${botId}/view`);
  },

  // 记录点击
  recordBotClick: (botId: string) => {
    return request.post(`/api/v1/discovery/bots/${botId}/click`);
  },
};
```

---

## 10. 错误码定义

| 错误码 | HTTP状态码 | 说明 | 处理建议 |
|--------|-----------|------|----------|
| 40001 | 404 | Bot不存在 | 检查Bot ID是否正确 |
| 40002 | 404 | 分类不存在 | 检查分类ID是否正确 |
| 40003 | 400 | 搜索关键词为空 | 请输入搜索关键词 |
| 40004 | 400 | 分类名称重复 | 同租户下分类名称不能重复 |
| 40005 | 400 | 配置项已存在 | 该Bot已配置在发现页 |
| 40101 | 403 | 无权限操作 | 需要管理员权限 |
| 40102 | 403 | Bot未上架 | 该Bot未在发现页展示 |

---

**文档结束**
