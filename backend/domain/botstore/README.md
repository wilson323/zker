# Bot商店模块

## 概述

Bot商店模块实现了Bot的发布、浏览、搜索和审核功能，采用DDD（领域驱动设计）架构。

## 架构

### 分层架构

```
domain/botstore/
├── entity/                    # 实体层
│   ├── bot_store_item.go      # 商店项目实体
│   └── bot_store_category.go  # 分类实体
├── repository/                # 仓储接口
│   └── repository.go          # 仓储接口定义
├── service/                   # 服务层
│   ├── publisher.go           # 发布器接口
│   ├── browser.go             # 浏览器接口
│   ├── reviewer.go            # 审核器接口
│   ├── bot_store_publisher_impl.go
│   ├── bot_store_browser_impl.go
│   └── bot_store_reviewer_impl.go
└── internal/                  # 内部实现
    └── dal/                   # 数据访问层
        ├── model/             # 数据模型
        └── *.go               # 仓储实现
```

### API层

```
api/
├── model/botstore/            # API模型
│   └── botstore.go
├── handler/coze/              # API处理器
│   └── bot_store_service.go
└── router/coze/               # 路由注册
    └── bot_store.go
```

## 核心功能

### 1. Bot发布

**功能**：将Bot发布到商店

**API**：`POST /api/v1/bot-store/publish`

**请求参数**：
```json
{
  "bot_id": "bot123",
  "name": "智能助手",
  "description": "一个强大的AI助手",
  "category": "cat_productivity",
  "tags": ["效率", "AI"],
  "price": 0.0,
  "screenshots": ["screenshot1.jpg"],
  "version": "1.0.0"
}
```

**流程**：
1. 验证请求参数
2. 检查Bot是否已发布
3. 验证分类是否存在
4. 创建商店项目（状态：pending）
5. 增加分类的Bot数量

### 2. Bot浏览

**功能**：列出商店中的Bot（分页）

**API**：`GET /api/v1/bot-store/list`

**查询参数**：
- `category`: 分类（可选）
- `sort_by`: 排序方式（popular/latest/rating）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20，最大100）

**排序规则**：
- `popular`: 下载次数降序 → 浏览次数降序
- `latest`: 创建时间降序
- `rating`: 评分降序 → 评分人数降序

### 3. Bot搜索

**功能**：搜索Bot

**API**：`GET /api/v1/bot-store/search`

**查询参数**：
- `query`: 搜索关键词（必填）
- `category`: 分类（可选）
- `tags`: 标签列表（可选）
- `price_min`: 最低价格（可选）
- `price_max`: 最高价格（可选）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20，最大100）

**搜索范围**：Bot名称和描述

### 4. Bot审核

**功能**：管理员审核Bot

**API**：`POST /api/v1/bot-store/admin/:item_id/review`

**请求参数**：
```json
{
  "approved": true,
  "reason": ""
}
```

**审核规则**：
- 只有`pending`状态的Bot可以审核
- 通过：状态变更为`published`
- 拒绝：状态变更为`rejected`，必须提供拒绝原因

## 状态机

```
draft ──────> pending ──────> published
 │                │               │
 │                │               │
 └────────────────┴───────────────┘
                  │
                  v
             rejected ───> pending (重新提交)

published ─────────> offline (下架)
```

## 数据库表

### bot_store_items

| 字段 | 类型 | 说明 |
|------|------|------|
| item_id | VARCHAR(36) | 主键 |
| bot_id | VARCHAR(36) | Bot ID（唯一） |
| tenant_id | VARCHAR(36) | 租户ID |
| name | VARCHAR(255) | Bot名称 |
| description | TEXT | Bot描述 |
| category | VARCHAR(100) | 分类 |
| tags | JSON | 标签列表 |
| price | DECIMAL(10,2) | 价格（0表示免费） |
| publisher_id | VARCHAR(36) | 发布者ID |
| status | VARCHAR(20) | 状态 |
| view_count | INT | 浏览次数 |
| download_count | INT | 下载次数 |
| rating | DECIMAL(3,2) | 评分（0.00-5.00） |
| rating_count | INT | 评分人数 |
| screenshots | JSON | 截图列表 |
| version | VARCHAR(50) | 版本号 |
| reject_reason | TEXT | 拒绝原因 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 删除时间（软删除） |

**索引**：
- PRIMARY KEY (item_id)
- UNIQUE KEY uk_bot_id (bot_id)
- INDEX idx_tenant_id (tenant_id)
- INDEX idx_status (status)
- INDEX idx_category (category)
- INDEX idx_status_created (status, created_at)
- INDEX idx_rating_download (rating DESC, download_count DESC)

### bot_store_categories

| 字段 | 类型 | 说明 |
|------|------|------|
| category_id | VARCHAR(36) | 主键 |
| name | VARCHAR(100) | 分类名称（唯一） |
| icon | VARCHAR(255) | 分类图标URL |
| description | TEXT | 分类描述 |
| parent_id | VARCHAR(36) | 父分类ID |
| sort_order | INT | 排序顺序 |
| is_active | BOOLEAN | 是否活跃 |
| bot_count | INT | 该分类下的Bot数量 |
| created_at | TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | 更新时间 |
| deleted_at | TIMESTAMP | 删除时间（软删除） |

**默认分类**：
- cat_productivity: 效率工具
- cat_entertainment: 娱乐休闲
- cat_education: 教育学习
- cat_business: 商业金融
- cat_development: 开发工具
- cat_health: 健康生活
- cat_creative: 创意设计
- cat_other: 其他

## 错误码

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| BSTORE400001 | 400 | 请求参数无效 |
| BSTORE400002 | 400 | Bot ID无效 |
| BSTORE400003 | 400 | Bot名称无效 |
| BSTORE400004 | 400 | 租户ID无效 |
| BSTORE400005 | 400 | 发布者ID无效 |
| BSTORE400006 | 400 | 价格无效 |
| BSTORE400007 | 400 | 分类无效 |
| BSTORE403001 | 403 | 无权访问Bot商店项目 |
| BSTORE403005 | 403 | Bot已发布 |
| BSTORE403006 | 403 | Bot待审核中 |
| BSTORE403007 | 403 | 无法修改已发布的项目 |
| BSTORE404001 | 404 | Bot商店项目不存在 |
| BSTORE404002 | 404 | 分类不存在 |
| BSTORE409001 | 409 | Bot状态无效 |
| BSTORE409002 | 409 | 拒绝原因必填 |
| BSTORE500001 | 500 | 创建Bot商店项目失败 |
| BSTORE500002 | 500 | 更新Bot商店项目失败 |
| BSTORE500003 | 500 | 删除Bot商店项目失败 |
| BSTORE500005 | 500 | 审核Bot商店项目失败 |

## 使用示例

### 发布Bot

```go
req := &service.PublishBotRequest{
    BotID:       "bot123",
    Name:        "智能助手",
    Description: "一个强大的AI助手",
    Category:    "cat_productivity",
    Tags:        []string{"效率", "AI"},
    Price:       0.0,
    PublisherID: "user123",
    TenantID:    "tenant123",
}

item, err := publisher.PublishBot(ctx, req)
if err != nil {
    // 处理错误
}
```

### 浏览Bot

```go
req := &service.ListBotsRequest{
    Category: "cat_productivity",
    SortBy:   "popular",
    Page:     1,
    PageSize: 20,
}

resp, err := browser.ListBots(ctx, req)
if err != nil {
    // 处理错误
}

for _, item := range resp.Items {
    fmt.Printf("Bot: %s, 评分: %.1f\n", item.Name, item.Rating)
}
```

### 审核Bot

```go
req := &service.ReviewBotRequest{
    ItemID:     "item123",
    Approved:   true,
    Reason:     "",
    ReviewerID: "admin123",
}

err := reviewer.ReviewBot(ctx, req)
if err != nil {
    // 处理错误
}
```

## 测试

运行单元测试：

```bash
cd backend/domain/botstore
go test ./... -cover
```

运行性能测试：

```bash
go test ./... -bench=. -benchmem
```

## 扩展功能

未来可以添加的功能：

1. **Bot评分系统**：用户可以对Bot进行评分
2. **Bot收藏**：用户可以收藏喜欢的Bot
3. **Bot评论**：用户可以对Bot发表评论
4. **Bot推荐**：基于用户行为的Bot推荐
5. **Bot统计**：发布者可以查看Bot的详细统计数据
6. **Bot版本管理**：支持多个版本的Bot
7. **Bot更新通知**：Bot更新时通知用户
8. **Bot订阅**：用户可以订阅付费Bot
9. **Bot分成**：与发布者进行收益分成
10. **Bot排行榜**：各种维度的Bot排行榜

## 贡献指南

1. 遵循DDD架构
2. 编写单元测试（覆盖率≥80%）
3. 使用统一的错误码
4. 添加Swagger注释
5. 更新文档

## 作者

Bot商店实现团队

## 许可证

Apache License 2.0
