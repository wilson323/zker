# 04-ZKER前台_发现_Discovery 详细设计说明书

**文档编号**: DE-DD-2025-004
**模块名称**: 发现_Discovery (Discovery)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03
**最后更新**: 2025-01-03

---

## 📋 文档修订历史

| 版本 | 日期 | 修订人 | 修订说明 |
|------|------|--------|----------|
| v1.0.0 | 2025-01-03 | ZKER Team | 初始版本，基于数据驱动设计 |

---

## 1. 模块概述

### 1.1 模块定位

**发现_Discovery** 是 ZKER 企业级 SaaS 平台的前台展示模块，通过**数据库驱动 + 通用展示组件**的方式，为企业用户提供Bot发现、浏览和快速访问能力。

**核心设计理念**：
- ✅ **数据驱动**：95% 通过数据库配置实现展示内容和顺序
- ✅ **通用组件**：5% 通用展示组件，复用性强
- ✅ **灵活配置**：管理员可配置展示哪些Bot、分类、排序规则
- ✅ **快速访问**：一键启动Bot对话

### 1.2 业务价值

| 受益者 | 价值 |
|--------|------|
| **企业用户** | 快速发现企业常用Bot，一键启动对话，提升工作效率 |
| **企业管理员** | 灵活配置发现页展示内容，推广重点Bot，提升Bot使用率 |
| **平台运营** | 通过数据分析优化Bot推荐，提升用户活跃度 |

### 1.3 与鲸智百应功能对齐

| 鲸智百应功能 | ZKER 实现方式 | 实现策略 |
|-------------|--------------|----------|
| Bot展示（卡片式） | ✅ 数据库配置 | discovery_items 表 |
| 分类浏览 | ✅ 数据库配置 | discovery_categories 表 |
| 搜索Bot | ✅ 通用搜索 | 复用 bots 表 |
| 推荐Bot | ✅ 配置推荐位 | is_featured 字段 |
| Bot使用量展示 | ✅ 统计数据 | 复用 bots.usage_count |
| 一键启动对话 | ✅ 跳转到会话页 | 前端路由跳转 |

**实现策略**：✅ 5% 编码（通用展示组件） + 95% 配置（数据库驱动）

---

## 2. 功能需求

### 2.1 核心功能清单

**F1 - Bot展示（核心功能）**
- F1.1 Bot卡片展示（名称、描述、图标、使用量）
- F1.2 分类展示（办公、营销、客服等）
- F1.3 搜索Bot（按名称、描述搜索）
- F1.4 筛选Bot（按分类、使用量、发布时间）
- F1.5 排序（推荐、最新、最热）

**F2 - 快速操作**
- F2.1 一键启动对话（点击Bot卡片打开对话）
- F2.2 收藏Bot（添加到个人收藏）
- F2.3 分享Bot（生成分享链接）
- F2.4 预览Bot（查看Bot详情）

**F3 - 管理配置（后台）**
- F3.1 配置展示Bot（选择哪些Bot显示在发现页）
- F3.2 配置分类（自定义分类）
- F3.3 配置推荐Bot（设置首页推荐位）
- F3.4 配置排序规则（自定义展示顺序）

### 2.2 非功能需求

| 需求类型 | 指标 | 说明 |
|---------|------|------|
| **性能** | 页面加载时间 | < 1 秒 |
| **可用性** | 系统可用性 | 99.9% （月度） |
| **可扩展性** | 新增展示Bot | 通过数据库配置，无需编码 |

---

## 3. 架构设计

### 3.1 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      前端层 (React 18)                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │ DiscoveryPage   │  │ BotCard         │  │ CategoryFilter  │  │
│  │ (发现页主容器)   │  │ (Bot卡片)       │  │ (分类筛选)      │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│  ┌─────────────────┐  ┌─────────────────┐                      │
│  │ SearchBar       │  │ FeaturedSection │                      │
│  │ (搜索栏)         │  │ (推荐区域)      │                      │
│  └─────────────────┘  └─────────────────┘                      │
└─────────────────────────────────────────────────────────────────┘
                              ↓ HTTP API
┌─────────────────────────────────────────────────────────────────┐
│                    业务服务层 (Go 1.23)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         DiscoveryService (简单的查询服务)                 │   │
│  │  - ListFeaturedBots()    (列出推荐Bot)                  │   │
│  │  - ListBotsByCategory()  (按分类列出Bot)                │   │
│  │  - SearchBots()          (搜索Bot)                      │   │
│  │  - GetCategories()       (获取分类)                     │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    数据访问层 (GORM)                             │
├─────────────────────────────────────────────────────────────────┤
│  discovery_categories | discovery_items                         │
│  bots (复用) | bot_favorites (复用)                             │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    基础设施层                                    │
├─────────────────────────────────────────────────────────────────┤
│  MySQL 8.4 | Redis 8.0 (缓存热门Bot)                            │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 核心流程

```
用户访问发现页
        ↓
前端加载分类（discovery_categories）
        ↓
前端加载推荐Bot（is_featured = true）
        ↓
前端加载所有Bot（按配置排序）
        ↓
用户交互：
  ├─ 点击分类 → 筛选Bot
  ├─ 搜索 → 关键词搜索Bot
  ├─ 点击Bot卡片 → 跳转到会话页
  └─ 收藏Bot → 添加到个人收藏
```

---

## 4. 数据库设计

### 4.1 核心表结构

#### 4.1.1 发现页分类表 (discovery_categories)

```sql
CREATE TABLE discovery_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台预置)',
    parent_id BIGINT DEFAULT NULL COMMENT '父分类ID',
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    description VARCHAR(200) COMMENT '分类描述',
    icon VARCHAR(10) COMMENT '图标 emoji',
    color VARCHAR(20) COMMENT '主题颜色',
    sort_order INT DEFAULT 0 COMMENT '排序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='发现页分类表';
```

#### 4.1.2 发现页展示配置表 (discovery_items)

**核心设计**：通过此表配置哪些Bot显示在发现页，以及展示顺序

```sql
CREATE TABLE discovery_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID (关联bots.id)',
    category_id BIGINT COMMENT '分类ID (关联discovery_categories.id)',

    -- 展示配置
    is_featured BOOLEAN DEFAULT FALSE COMMENT '是否推荐到首页',
    featured_position INT COMMENT '推荐位置 (1-10)',
    featured_expires_at DATETIME COMMENT '推荐过期时间',

    sort_order INT DEFAULT 0 COMMENT '展示顺序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否显示',

    -- 统计
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    click_count INT DEFAULT 0 COMMENT '点击次数',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY uk_tenant_bot (tenant_id, bot_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_category_id (category_id),
    INDEX idx_is_featured (is_featured),
    INDEX idx_featured_position (featured_position),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='发现页展示配置表';
```

**注意**：
- `discovery_items` 表是核心配置表，决定哪些Bot显示在发现页
- `bots` 表已在"07-数字员工管理"模块中定义，此处直接复用
- 通过 `bot_id` 关联 `bots` 表获取Bot详细信息

### 4.2 ER 图

```
┌──────────────────┐         ┌──────────────────┐
│discovery_categories│         │   bots (复用)     │
├──────────────────┤         ├──────────────────┤
│ id (PK)          │         │ id (PK)          │
│ tenant_id        │         │ tenant_id        │
│ name             │         │ name             │
│ icon             │         │ description      │
│ sort_order       │         │ avatar           │
└─────────┬────────┘         │ usage_count      │
          │                  └─────────┬────────┘
          │                            │
          │         ┌──────────────────┴──────────┐
          │         │                              │
          ↓         │                              │
┌──────────────────┴─────────────┐               │
│     discovery_items             │               │
├─────────────────────────────────┤               │
│ id (PK)                         │               │
│ tenant_id                       │               │
│ bot_id (FK)─────────────────────┘               │
│ category_id (FK)                               │
│ is_featured                                    │
│ sort_order                                     │
└────────────────────────────────────────────────┘
```

---

## 5. API 设计

### 5.1 REST API 列表

| 方法 | 路径 | 功能 |
|------|------|------|
| **发现页数据** |
| GET | /api/v1/discovery/categories | 获取分类列表 |
| GET | /api/v1/discovery/featured | 获取推荐Bot |
| GET | /api/v1/discovery/bots | 获取所有展示Bot |
| GET | /api/v1/discovery/bots/:id | 获取Bot详情 |
| **搜索与筛选** |
| GET | /api/v1/discovery/search | 搜索Bot |
| GET | /api/v1/discovery/category/:id/bots | 按分类获取Bot |

### 5.2 核心 API 详细设计

#### 5.2.1 获取发现页数据

**请求**：
```http
GET /api/v1/discovery/bots?category=5&featured=true&sort=featured&limit=20
Authorization: Bearer <token>
X-Tenant-ID: <tenant_id>
```

**响应**：
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
      }
    ],
    "featured": [
      {
        "id": "bot-001",
        "name": "邮件助手",
        "description": "智能撰写和回复邮件",
        "avatar": "https://cdn.zker.com/avatars/bot-email.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 1250,
        "rating": 4.8,
        "isFeatured": true,
        "featuredPosition": 1
      }
    ],
    "bots": [
      {
        "id": "bot-002",
        "name": "周报生成器",
        "description": "快速生成工作周报",
        "avatar": "https://cdn.zker.com/avatars/bot-report.png",
        "category": {
          "id": 1,
          "name": "办公效率"
        },
        "usageCount": 890,
        "rating": 4.6
      }
    ],
    "total": 150
  }
}
```

**Go 代码示例**：

```go
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
        return nil, err
    }

    // 2. 获取推荐Bot（缓存）
    var featuredBots []*DiscoveryBotItem
    if req.Featured {
        featuredBots, err = s.getCachedFeaturedBots(ctx, tenantID)
        if err != nil {
            s.logger.Warn("failed to get cached featured bots", zap.Error(err))
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
        return nil, err
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
```

#### 5.2.2 搜索Bot

**请求**：
```http
GET /api/v1/discovery/search?q=邮件&category=1&limit=10
Authorization: Bearer <token>
X-Tenant-ID: <tenant_id>
```

**响应**：
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
        "avatar": "https://cdn.zker.com/avatars/bot-email.png",
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
      }
    ],
    "total": 5
  }
}
```

---

## 6. 前端设计

### 6.1 核心组件

#### 6.1.1 DiscoveryPage（发现页主容器）

```tsx
import React, { useState, useEffect } from 'react';
import { Input, Tabs, Empty, Spin } from '@douyinfe/semi-ui';
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
      message.error('加载数据失败');
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
      message.error('搜索失败');
    }
  };

  const handleBotClick = (botId: string) => {
    navigate(`/chat?bot=${botId}`);
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
          onEnterPress={(e: any) => handleSearch(e.target.value)}
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

#### 6.1.2 BotCard（Bot卡片组件）

```tsx
import React from 'react';
import { Card, Tag, Button } from '@douyinfe/semi-ui';
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
          <Tag color="blue">{bot.category?.name}</Tag>
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

#### 6.1.3 FeaturedSection（推荐区域）

```tsx
import React from 'react';
import { Carousel } from '@douyinfe/semi-ui';
import { BotCard } from './BotCard';
import styles from './FeaturedSection.module.scss';

interface FeaturedSectionProps {
  bots: Bot[];
  onBotClick: (botId: string) => void;
}

export const FeaturedSection: React.FC<FeaturedSectionProps> = ({ bots, onBotClick }) => {
  return (
    <div className={styles.featuredSection}>
      <h2 className={styles.title}>✨ 精选推荐</h2>
      <Carousel
        className={styles.carousel}
        showArrow={true}
        autoplay={false}
        dots={false}
      >
        {bots.map(bot => (
          <div key={bot.id} className={styles.slide}>
            <BotCard bot={bot} onClick={() => onBotClick(bot.id)} />
          </div>
        ))}
      </Carousel>
    </div>
  );
};
```

### 6.2 样式设计

**SCSS 示例**（DiscoveryPage.module.scss）：

```scss
.discoveryPage {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;

  .header {
    margin-bottom: 32px;

    h1 {
      font-size: 32px;
      font-weight: 600;
      margin-bottom: 16px;
    }

    .searchBar {
      max-width: 600px;
    }
  }

  .botGrid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
    margin-top: 24px;
  }
}

.botCard {
  position: relative;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  }

  .content {
    .name {
      font-size: 18px;
      font-weight: 600;
      margin-bottom: 8px;
    }

    .description {
      color: #666;
      font-size: 14px;
      margin-bottom: 12px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .meta {
      display: flex;
      gap: 12px;
      align-items: center;
      font-size: 12px;
      color: #999;
    }
  }

  .featuredBadge {
    position: absolute;
    top: 12px;
    right: 12px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
  }
}
```

---

## 7. 配置系统设计

### 7.1 预置分类数据

```sql
-- 平台预置分类
INSERT INTO discovery_categories (tenant_id, name, description, icon, color, sort_order) VALUES
(NULL, '办公效率', '提升办公效率的Bot', '💼', '#1890ff', 1),
(NULL, '营销推广', '营销与推广Bot', '📢', '#52c41a', 2),
(NULL, '客户服务', '客服与支持Bot', '🎧', '#faad14', 3),
(NULL, '数据分析', '数据分析Bot', '📊', '#13c2c2', 4),
(NULL, '创意写作', '写作与创作Bot', '✍️', '#eb2f96', 5);
```

### 7.2 配置展示Bot

```sql
-- 配置Bot显示在发现页
INSERT INTO discovery_items (tenant_id, bot_id, category_id, is_featured, featured_position, sort_order) VALUES
-- 推荐Bot
('tenant-abc-001', 'bot-email-assistant', 1, TRUE, 1, 1),
('tenant-abc-001', 'bot-report-generator', 1, TRUE, 2, 2),
('tenant-abc-001', 'bot-social-poster', 2, TRUE, 3, 3),

-- 普通展示Bot
('tenant-abc-001', 'bot-calendar-helper', 1, FALSE, NULL, 10),
('tenant-abc-001', 'bot-seo-analyzer', 2, FALSE, NULL, 11),
('tenant-abc-001', 'bot-customer-service', 3, FALSE, NULL, 12);
```

### 7.3 企业自定义分类

```sql
-- 企业级自定义分类
INSERT INTO discovery_categories (tenant_id, name, icon, sort_order) VALUES
('tenant-abc-001', '销售支持', '🎯', 10),
('tenant-abc-001', '人力资源', '👥', 11);
```

---

## 8. 总结

### 8.1 实施策略总结

| 实施项 | 实施方式 | 工作量 | 说明 |
|-------|---------|--------|------|
| **前端展示组件** | 💻 通用组件 | 5% | BotCard、CategoryFilter等 |
| **后端查询API** | 💻 简单CRUD | 5% | 基础查询服务 |
| **展示配置** | 📊 数据库驱动 | 90% | discovery_items 表配置 |

**总计**：5% 编码（通用组件） + 95% 配置（数据库驱动）

### 8.2 与鲸智百应对齐情况

| 功能模块 | 对齐情况 | 实现方式 |
|---------|---------|----------|
| Bot展示 | ✅ 100% | 数据库配置 |
| 分类浏览 | ✅ 100% | 数据库配置 |
| 搜索Bot | ✅ 100% | 通用搜索 |
| 推荐Bot | ✅ 100% | is_featured 配置 |
| 一键启动 | ✅ 100% | 前端路由跳转 |

### 8.3 核心优势

- ✅ **配置驱动**：95% 通过数据库配置，灵活调整展示内容
- ✅ **通用组件**：可复用的Bot卡片、分类筛选组件
- ✅ **性能优化**：Redis 缓存热门Bot，加载速度快
- ✅ **扩展性强**：新增展示Bot仅需数据库配置，无需编码

---

**文档结束**
