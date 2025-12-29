# ZKER 开发快速入门指南

> **文档类型**: 开发指南 / 快速入门
> **版本**: v1.0
> **生成日期**: 2025-01-01
> **目标读者**: 新加入的开发人员
> **学习目标**: 3-5 天内能够开始独立开发

---

## 📋 目录

1. [欢迎加入 ZKER](#欢迎加入-zker)
2. [环境搭建](#环境搭建)
3. [核心概念](#核心概念)
4. [前端开发指南](#前端开发指南)
5. [后端开发指南](#后端开发指南)
6. [开发规范](#开发规范)
7. [常见问题](#常见问题)
8. [快速参考](#快速参考)

---

## 欢迎加入 ZKER

### 1.1 ZKER 是什么?

**ZKER** 是一个一站式 AI Agent 开发平台,采用微服务架构,支持:

✅ **智能 Bot 开发**: 可视化创建和配置 AI Bot
✅ **多模态交互**: 支持文本、图片、语音、视频
✅ **知识库管理**: 企业级知识库,支持向量和关键词检索
✅ **工作流编排**: 可视化设计复杂工作流
✅ **智能路由**: 意图识别 + 路由决策 + 负载均衡
✅ **多租户 SaaS**: 企业级多租户隔离

### 1.2 技术栈概览

**前端技术栈**:
- React 18 + TypeScript
- Semi Design UI 框架
- Tailwind CSS
- Zustand 状态管理
- React Query 数据请求
- Rsbuild 快速构建

**后端技术栈**:
- Go 1.21+
- Hertz HTTP 框架
- GORM ORM 框架
- DDD 领域驱动设计
- NSQ 消息队列
- MySQL 8.4.5 + Redis 8.0

**基础设施**:
- Docker + Docker Compose
- etcd 服务发现
- MinIO 对象存储
- Milvus 向量数据库
- Elasticsearch 搜索引擎

### 1.3 本文档学习路径

**Day 1**: 环境搭建 + 核心概念
**Day 2**: 前端开发入门
**Day 3**: 后端开发入门
**Day 4-5**: 实战项目

---

## 环境搭建

### 2.1 前置要求

**必需软件**:
- Git
- Docker Desktop (Windows/Mac) 或 Docker (Linux)
- Node.js 18+
- Go 1.21+
- Make 工具
- 代码编辑器: VS Code (推荐)

**推荐软件**:
- Postman: API 测试
- Redis Desktop Manager: Redis 可视化
- TablePlus: 数据库管理

### 2.2 克隆代码

```bash
# 1. 克隆仓库
git clone https://github.com/coze-studio/zker.git
cd zker

# 2. 查看分支
git branch -a

# 3. 切换到开发分支
git checkout develop

# 4. 查看最近提交
git log --oneline -10
```

### 2.3 启动基础设施 (Docker 方式)

**推荐**: 使用 Docker 快速启动所有中间件服务

```bash
# 1. 进入 docker 目录
cd docker

# 2. 复制环境配置文件
cp .env.example .env

# 3. 编辑 .env 文件,配置必要参数
# - MySQL root 密码
# - Redis 密码
# - API Keys (可选)

# 4. 启动所有服务
docker compose up -d

# 5. 查看服务状态
docker compose ps

# 预期看到以下服务都在运行:
# - mysql (端口 3306)
# - redis (端口 6379)
# - elasticsearch (端口 9200)
# - milvus (端口 19530)
# - minio (端口 9000)
# - etcd (端口 2379)
# - nsqlookupd (端口 4160)
# - nsqd (端口 4150)
```

**验证服务**:

```bash
# 检查 MySQL
docker exec -it zker-mysql mysql -uroot -p
# 输入密码 (默认: root123456)
SHOW DATABASES;

# 检查 Redis
docker exec -it zker-redis redis-cli
PING
# 应该返回 PONG

# 检查 Elasticsearch
curl http://localhost:9200/_cluster/health
```

### 2.4 初始化数据库

```bash
# 回到项目根目录
cd ..

# 方式1: 使用 Make 命令
make sync_db    # 同步数据库结构
make sql_init   # 初始化基础数据

# 方式2: 手动执行
# 连接数据库
mysql -h127.0.0.1 -P3306 -uroot -p

# 执行 DDL 脚本
source docs/企业级功能完善与统一性设计方案/database_schema_ddl.sql
source docs/企业级功能完善与统一性设计方案/database_schema_ddl_part2_intelligent_routing.sql
source docs/企业级功能完善与统一性设计方案/database_schema_ddl_part3_other_modules.sql
source docs/企业级功能完善与统一性设计方案/database_schema_ddl_part4_billing_monitoring.sql

# 执行初始化数据脚本
source docs/企业级功能完善与统一性设计方案/database_init_data.sql
```

**验证数据库**:

```sql
-- 检查表数量 (应该是 157 张表)
SELECT COUNT(*) FROM information_schema.tables
WHERE table_schema = 'zker';

-- 检查默认租户
SELECT * FROM tenants WHERE tenant_id = 'default';

-- 检查管理员用户
SELECT * FROM users WHERE user_id = 'admin';
# 默认密码: admin123
```

### 2.5 安装前端依赖

```bash
# 使用 Rush 安装所有前端包依赖
cd frontend
rush update

# 预计耗时: 5-10 分钟 (取决于网络速度)

# 验证安装
ls common/    # 应该有 node_modules
ls arch/     # 应该有 node_modules
```

### 2.6 配置环境变量

```bash
# 在项目根目录创建 .env.local 文件
cat > .env.local << 'EOF'
# API 地址
VITE_API_BASE_URL=http://localhost:8001

# WebSocket 地址
VITE_WS_BASE_URL=ws://localhost:8001

# 应用配置
VITE_APP_NAME=ZKER
VITE_APP_PORT=8888

# 第三方服务 (可选)
OPENAI_API_KEY=your_key_here
EOF

# 前端环境变量需要以 VITE_ 开头
```

### 2.7 启动开发服务

**启动后端服务**:

```bash
# 方式1: 使用 Make 命令 (推荐)
make server

# 方式2: 手动启动
cd backend
go run cmd/main.go

# 后端默认运行在 :8001
```

**启动前端服务**:

```bash
# 方式1: 使用 Make 命令
make fe

# 方式2: 手动启动
cd frontend/apps/coze-studio
npm run dev

# 前端默认运行在 :8888
```

**验证服务**:

```bash
# 检查后端健康状态
curl http://localhost:8001/health

# 访问前端
open http://localhost:8888

# 默认登录账号:
# 用户名: admin
# 密码: admin123
```

### 2.8 环境搭建检查清单

使用以下清单验证环境是否搭建成功:

- [ ] Docker 所有服务正常运行
- [ ] MySQL 数据库初始化完成 (157张表)
- [ ] Redis 连接正常
- [ ] 后端服务启动无错误
- [ ] 前端页面可以访问
- [ ] 可以正常登录系统
- [ ] API 请求正常响应

---

## 核心概念

### 3.1 系统架构

**整体架构**:

```
┌─────────────────────────────────────────────────┐
│                   前端层                       │
│  React 18 + TypeScript + Semi Design          │
│  (Rush.js Monorepo - 135+ 包)                │
└──────────────────┬──────────────────────────────┘
                   │ HTTP/WebSocket
┌──────────────────┴──────────────────────────────┐
│                 API 网关                        │
│              (Nginx / Kong)                     │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────┐
│                微服务层                         │
│  ┌──────────┬──────────┬──────────┬─────────┐ │
│  │ 用户服务  │ Bot服务  │ 会话服务 │  ...    │ │
│  │ :8001    │ :8003    │ :8004    │         │ │
│  └──────────┴──────────┴──────────┴─────────┘ │
└──────────────────┬──────────────────────────────┘
                   │
┌──────────────────┴──────────────────────────────┐
│              基础设施层                         │
│  MySQL │ Redis │ Milvus │ ES │ MinIO │ NSQ   │
└─────────────────────────────────────────────────┘
```

### 3.2 前端架构

**Monorepo 结构** (Rush.js 管理):

```
frontend/
├── arch/          # Level 1: 核心基础设施 (不依赖业务)
│   ├── bot-hooks/           # Bot 相关 Hooks
│   ├── bot-utils/           # Bot 工具函数
│   ├── hooks/               # 通用 Hooks
│   └── logger/              # 日志工具
│
├── common/        # Level 2: 共享组件和工具
│   ├── chat-area/           # 聊天区域组件
│   ├── biz-components/      # 业务组件
│   └── uploader-adapter/    # 上传器适配器
│
├── agent-ide/     # Level 3: Agent IDE 功能域
│   ├── bot-config-area/     # Bot 配置区域
│   ├── bot-plugin/          # Bot 插件系统
│   └── space-bot/           # 空间 Bot
│
├── workflow/      # Level 3: 工作流功能域
│   ├── playground/          # 工作流 Playground
│   ├── render/              # 工作流渲染
│   └── components/          # 工作流组件
│
└── apps/          # Level 4: 应用入口
    └── coze-studio/         # 主应用
```

**关键设计模式**:

1. **适配器模式**: 大量使用 `-adapter` 后缀包实现层间解耦
2. **Base/Core 模式**: 共享功能放在 `-base` 或 `-core` 包
3. **Context 模式**: 使用 React Context 管理组件树状态
4. **Store 模式**: 使用 Zustand 管理全局状态

### 3.3 后端架构

**DDD 分层架构**:

```
backend/
├── domain/        # 领域层 - 实体和业务规则
│   ├── user/               # 用户领域
│   │   ├── user.go        # 用户实体
│   │   ├── repository.go  # 仓储接口
│   │   └── service.go     # 领域服务
│   └── bot/                # Bot 领域
│
├── application/   # 应用层 - 用例编排
│   ├── user/
│   │   └── usecase.go     # 用户用例
│   └── bot/
│
├── api/           # API 层 - HTTP 处理器
│   ├── user/
│   │   └── handler.go     # 用户 HTTP 处理器
│   └── bot/
│
├── infra/         # 基础设施层 - 技术实现
│   ├── database/           # 数据库实现
│   │   ├── mysql.go       # MySQL 实现
│   │   └── repository.go  # 仓储实现
│   └── cache/              # 缓存实现
│
└── crossdomain/   # 跨领域 - 关注点横切
    ├── auth/               # 认证授权
    ├── logging/            # 日志记录
    └── metrics/            # 监控指标
```

**关键设计原则**:

1. **依赖倒置**: 领域层不依赖基础设施层
2. **接口隔离**: 定义清晰的 Repository 接口
3. **单一职责**: 每层只负责自己的职责
4. **CQRS**: 读写分离 (复杂场景)

### 3.4 多租户架构

**租户隔离策略**:

- **数据库**: 共享数据库,共享 Schema,通过 `tenant_id` 字段隔离
- **应用层**: 自动添加 `tenant_id` 过滤条件
- **缓存**: Key 中包含 `tenant:{tenant_id}` 前缀

**租户 ID 传递**:

```go
// HTTP Header
X-Tenant-ID: tenant_123

// JWT Token Claims
{
  "user_id": "user_456",
  "tenant_id": "tenant_123"
}

// 上下文
ctx = context.WithValue(ctx, "tenant_id", "tenant_123")
```

### 3.5 RBAC 权限模型

**五级数据权限**:

1. **全部** (ALL): 可访问所有数据
2. **组织** (ORGANIZATION): 只能访问本组织数据
3. **部门** (DEPARTMENT): 只能访问本部门数据
4. **自己** (OWN): 只能访问自己的数据
5. **自定义** (CUSTOM): 自定义规则

**三级字段权限**:

1. **可见** (VISIBLE): 可以查看字段
2. **可编辑** (EDITABLE): 可以修改字段
3. **隐藏** (HIDDEN): 字段不可见

---

## 前端开发指南

### 4.1 快速开始

**创建一个新页面**:

```bash
# 1. 在 apps/coze-studio/src/pages/ 下创建新页面
cd frontend/apps/coze-studio/src/pages
mkdir my-page

# 2. 创建页面组件
cd my-page
touch index.tsx
```

```typescript
// index.tsx
import React from 'react';
import { Button } from '@coze-studio/ui';

export const MyPage: React.FC = () => {
  return (
    <div className="p-4">
      <h1>我的页面</h1>
      <Button type="primary">点击我</Button>
    </div>
  );
};

export default MyPage;
```

**添加路由**:

```typescript
// 在 apps/coze-studio/src/routes/index.tsx 中添加路由

import { MyPage } from '@/pages/my-page';

// 添加路由
{
  path: 'my-page',
  element: <MyPage />,
}
```

### 4.2 使用组件

**使用 UI 组件**:

```typescript
import { Button, Input, Modal, Table } from '@coze-studio/ui';

function MyComponent() {
  const [visible, setVisible] = useState(false);

  return (
    <div>
      <Button type="primary" onClick={() => setVisible(true)}>
        打开对话框
      </Button>

      <Modal
        visible={visible}
        title="标题"
        onOk={() => setVisible(false)}
        onCancel={() => setVisible(false)}
      >
        <Input placeholder="请输入" />
      </Modal>
    </div>
  );
}
```

**使用业务组件**:

```typescript
import { BotSelector } from '@coze-studio/bot';
import { ChatInput } from '@coze-studio/chat';

function MyComponent() {
  const [botId, setBotId] = useState('');

  return (
    <div>
      <BotSelector
        value={botId}
        onChange={setBotId}
        placeholder="选择 Bot"
      />

      <ChatInput
        onSend={(content) => console.log(content)}
      />
    </div>
  );
}
```

### 4.3 状态管理

**使用 Zustand Store**:

```typescript
// 创建 Store (stores/my-store.ts)
import { create } from 'zustand';

interface MyStore {
  count: number;
  increment: () => void;
  decrement: () => void;
}

export const useMyStore = create<MyStore>((set) => ({
  count: 0,
  increment: () => set((state) => ({ count: state.count + 1 })),
  decrement: () => set((state) => ({ count: state.count - 1 })),
}));

// 在组件中使用
import { useMyStore } from '@/stores/my-store';

function MyComponent() {
  const { count, increment, decrement } = useMyStore();

  return (
    <div>
      <p>Count: {count}</p>
      <Button onClick={increment}>+</Button>
      <Button onClick={decrement}>-</Button>
    </div>
  );
}
```

**使用 React Query**:

```typescript
import { useQuery, useMutation } from '@tanstack/react-query';
import { botApi } from '@/api/bot';

function BotList() {
  // 查询数据
  const { data, isLoading, error } = useQuery({
    queryKey: ['bots'],
    queryFn: () => botApi.list({ page: 1, pageSize: 20 }),
  });

  // 变更数据
  const createMutation = useMutation({
    mutationFn: botApi.create,
    onSuccess: () => {
      // 刷新列表
      queryClient.invalidateQueries(['bots']);
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      {data.items.map((bot) => (
        <div key={bot.bot_id}>{bot.name}</div>
      ))}

      <Button onClick={() => createMutation.mutate({ name: 'New Bot' })}>
        创建 Bot
      </Button>
    </div>
  );
}
```

### 4.4 API 调用

**定义 API**:

```typescript
// api/bot.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
});

// 自动添加 Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const botApi = {
  // 获取 Bot 列表
  list: (params: { page: number; pageSize: number }) =>
    api.get('/api/v1/bots', { params }).then((res) => res.data),

  // 创建 Bot
  create: (data: { name: string; description: string }) =>
    api.post('/api/v1/bots', data).then((res) => res.data),

  // 获取 Bot 详情
  get: (botId: string) =>
    api.get(`/api/v1/bots/${botId}`).then((res) => res.data),

  // 更新 Bot
  update: (botId: string, data: any) =>
    api.put(`/api/v1/bots/${botId}`, data).then((res) => res.data),

  // 删除 Bot
  delete: (botId: string) =>
    api.delete(`/api/v1/bots/${botId}`).then((res) => res.data),
};
```

### 4.5 样式开发

**使用 Tailwind CSS**:

```tsx
function MyComponent() {
  return (
    <div className="p-4 bg-white rounded-lg shadow">
      <h1 className="text-2xl font-bold text-gray-900">
        标题
      </h1>
      <p className="mt-2 text-gray-600">
        内容
      </p>
      <Button className="mt-4" type="primary">
        按钮
      </Button>
    </div>
  );
}
```

**使用 CSS Modules**:

```less
// styles.module.less
.container {
  padding: 16px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  .title {
    font-size: 24px;
    font-weight: bold;
    color: #111827;
  }

  .content {
    margin-top: 8px;
    color: #6B7280;
  }
}
```

```tsx
import styles from './styles.module.less';

function MyComponent() {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>标题</h1>
      <p className={styles.content}>内容</p>
    </div>
  );
}
```

---

## 后端开发指南

### 5.1 项目结构

**创建新的 API 接口**:

```go
// 1. 定义领域实体 (domain/bot/bot.go)
package bot

type Bot struct {
    BotID      string    `json:"bot_id" gorm:"primaryKey"`
    TenantID   string    `json:"tenant_id" gorm:"index"`
    Name       string    `json:"name" gorm:"not null"`
    Status     string    `json:"status" gorm:"default:'draft'"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}

// 2. 定义 Repository 接口 (domain/bot/repository.go)
type Repository interface {
    Create(ctx context.Context, bot *Bot) error
    FindByID(ctx context.Context, botID string) (*Bot, error)
    Update(ctx context.Context, bot *Bot) error
    Delete(ctx context.Context, botID string) error
    List(ctx context.Context, tenantID string, page, pageSize int) ([]*Bot, int64, error)
}

// 3. 实现 Repository (infra/database/repository.go)
type repositoryImpl struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) bot.Repository {
    return &repositoryImpl{db: db}
}

func (r *repositoryImpl) Create(ctx context.Context, bot *Bot) error {
    return r.db.WithContext(ctx).Create(bot).Error
}

func (r *repositoryImpl) FindByID(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    err := r.db.WithContext(ctx).Where("bot_id = ?", botID).First(&bot).Error
    if err != nil {
        return nil, err
    }
    return &bot, nil
}

// 4. 定义 Service (domain/bot/service.go)
type Service struct {
    repo bot.Repository
}

func NewService(repo bot.Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    // 参数验证
    if req.Name == "" {
        return nil, errors.New("bot name cannot be empty")
    }

    // 创建实体
    bot := &Bot{
        BotID:    generateUUID(),
        TenantID: getTenantID(ctx),
        Name:     req.Name,
        Status:   "draft",
    }

    // 保存
    if err := s.repo.Create(ctx, bot); err != nil {
        return nil, err
    }

    return bot, nil
}

// 5. 定义 HTTP Handler (api/bot/handler.go)
type Handler struct {
    service *bot.Service
}

func NewHandler(service *bot.Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) CreateBot(c *app.RequestContext) {
    var req CreateBotRequest
    if err := c.Bind(&req); err != nil {
        c.JSON(400, gin.H{"code": "SYS100", "message": "请求参数错误"})
        return
    }

    bot, err := h.service.CreateBot(c.Context(), &req)
    if err != nil {
        c.JSON(500, gin.H{"code": "BOT500", "message": err.Error()})
        return
    }

    c.JSON(201, gin.H{
        "code":    "SUCCESS",
        "message": "Bot 创建成功",
        "data":    bot,
    })
}

// 6. 注册路由 (api/router.go)
func RegisterRoutes(r *hertz.Engine) {
    botService := bot.NewService(botRepo)
    botHandler := bot.NewHandler(botService)

    botGroup := r.Group("/api/v1/bots")
    {
        botGroup.POST("", botHandler.CreateBot)
        botGroup.GET("", botHandler.ListBots)
        botGroup.GET("/:bot_id", botHandler.GetBot)
        botGroup.PUT("/:bot_id", botHandler.UpdateBot)
        botGroup.DELETE("/:bot_id", botHandler.DeleteBot)
    }
}
```

### 5.2 数据库操作

**使用 GORM**:

```go
// 查询单条
var bot Bot
err := db.Where("bot_id = ?", botID).First(&bot).Error
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrBotNotFound
    }
    return nil, err
}

// 查询列表
var bots []Bot
err := db.Where("tenant_id = ? AND status = ?", tenantID, "published").
    Limit(pageSize).
    Offset((page - 1) * pageSize).
    Find(&bots).Error

// 统计总数
var total int64
db.Model(&Bot{}).Where("tenant_id = ? AND status = ?", tenantID, "published").Count(&total)

// 创建
if err := db.Create(&bot).Error; err != nil {
    return err
}

// 更新
if err := db.Model(&bot).Updates(map[string]interface{}{
    "name":   newName,
    "status": newStatus,
}).Error; err != nil {
    return err
}

// 软删除
if err := db.Delete(&bot).Error; err != nil {
    return err
}
```

### 5.3 错误处理

**定义错误类型**:

```go
// pkg/errors/errors.go
package errors

type AppError struct {
    Code       string
    Message    string
    StatusCode int
}

func (e *AppError) Error() string {
    return e.Message
}

// 预定义错误
var (
    ErrBotNotFound      = &AppError{"BOT300", "Bot 不存在", 404}
    ErrBotNameConflict  = &AppError{"BOT402", "Bot 名称已存在", 409}
    ErrBotQuotaExceeded = &AppError{"BOT401", "Bot 数量超限", 422}
)

// 使用错误
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrBotNotFound
    }
    return nil, err
}
```

**统一错误响应**:

```go
// api/middleware/error_handler.go
func ErrorHandler() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        c.Next(ctx)

        if len(c.Errors) > 0 {
            err := c.Errors.Last()

            if appErr, ok := err.Err.(*errors.AppError); ok {
                c.JSON(appErr.StatusCode, gin.H{
                    "code":    appErr.Code,
                    "message": appErr.Message,
                })
            } else {
                c.JSON(500, gin.H{
                    "code":    "SYS500",
                    "message": "系统内部错误",
                })
            }
        }
    }
}
```

### 5.4 中间件使用

**认证中间件**:

```go
// api/middleware/auth.go
func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 获取 Token
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"code": "AUTH200", "message": "未登录"})
            c.Abort()
            return
        }

        // 2. 验证 Token
        claims, err := ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"code": "AUTH202", "message": "Token 无效"})
            c.Abort()
            return
        }

        // 3. 存储到上下文
        c.Set("user_id", claims.UserID)
        c.Set("tenant_id", claims.TenantID)

        c.Next(ctx)
    }
}

// 使用
botGroup := r.Group("/api/v1/bots")
botGroup.Use(AuthMiddleware())
{
    botGroup.POST("", botHandler.CreateBot)
    botGroup.GET("", botHandler.ListBots)
}
```

**租户隔离中间件**:

```go
// api/middleware/tenant.go
func TenantIsolationMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            // 从 JWT Token 中获取
            tenantID = c.GetString("tenant_id")
        }

        if tenantID == "" {
            c.JSON(400, gin.H{"code": "SYS101", "message": "缺少租户 ID"})
            c.Abort()
            return
        }

        // 检查租户状态
        tenant, err := tenantService.GetByID(ctx, tenantID)
        if err != nil || tenant.Status != "active" {
            c.JSON(403, gin.H{"code": "TARG300", "message": "租户不存在或已停用"})
            c.Abort()
            return
        }

        c.Set("tenant_id", tenantID)
        c.Next(ctx)
    }
}
```

---

## 开发规范

### 6.1 Git 工作流

**分支策略**:

```
main (生产环境)
  ↑
  develop (开发环境)
    ↑
    feature/xxx (功能分支)
    bugfix/xxx (Bug 修复分支)
    hotfix/xxx (紧急修复分支)
```

**分支命名规范**:
- `feature/bot-management` - 新功能开发
- `bugfix/login-error` - Bug 修复
- `hotfix/security-patch` - 紧急修复

**提交流程**:

```bash
# 1. 创建功能分支
git checkout develop
git pull origin develop
git checkout -b feature/my-feature

# 2. 开发并提交
git add .
git commit -m "feat: 添加 Bot 管理功能"

# 3. 推送到远程
git push origin feature/my-feature

# 4. 创建 Pull Request
# 在 GitHub 上创建 PR,请求合并到 develop

# 5. 代码审查通过后合并
# 删除本地分支
git branch -d feature/my-feature
```

**提交信息规范** (Conventional Commits):

```
feat: 新功能
fix: Bug 修复
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
perf: 性能优化
test: 测试相关
chore: 构建/工具相关
```

**示例**:
```bash
git commit -m "feat: 添加 Bot 保存接口"
git commit -m "fix: 修复登录 Token 过期问题"
git commit -m "docs: 更新 API 文档"
```

### 6.2 代码审查清单

**前端代码审查**:

- [ ] 组件命名清晰
- [ ] Props 类型定义完整
- [ ] 状态管理逻辑合理
- [ ] 没有性能问题(不必要的渲染)
- [ ] 样式使用 Tailwind 或 CSS Modules
- [ ] 错误处理完善
- [ ] 代码格式符合 ESLint 规则

**后端代码审查**:

- [ ] 遵循 DDD 分层架构
- [ ] 接口定义清晰
- [ ] 错误处理完善
- [ ] 数据库操作有事务保护
- [ ] 没有 SQL 注入风险
- [ ] 代码格式符合 golangci-lint 规则
- [ ] 有必要的注释

### 6.3 测试规范

**前端测试**:

```typescript
// 组件测试
import { render, screen } from '@testing-library/react';
import { Button } from './Button';

test('renders button with text', () => {
  render(<Button>Click me</Button>);
  expect(screen.getByText('Click me')).toBeInTheDocument();
});

test('calls onClick when clicked', () => {
  const handleClick = jest.fn();
  render(<Button onClick={handleClick}>Click me</Button>);

  screen.getByText('Click me').click();

  expect(handleClick).toHaveBeenCalledTimes(1);
});
```

**后端测试**:

```go
func TestBotService_CreateBot(t *testing.T) {
    // Arrange
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)
    req := &CreateBotRequest{Name: "Test Bot"}

    // Act
    bot, err := service.CreateBot(context.Background(), req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, bot)
    assert.Equal(t, "Test Bot", bot.Name)
    assert.Equal(t, "draft", bot.Status)
}
```

---

## 常见问题

### 7.1 环境问题

**Q: Docker 服务启动失败?**

A: 检查端口是否被占用

```bash
# 检查端口占用
lsof -i :3306  # macOS
netstat -ano | findstr :3306  # Windows

# 停止占用端口的服务或修改 docker-compose.yml 中的端口映射
```

**Q: 前端依赖安装失败?**

A: 使用国内镜像

```bash
# 配置 npm 国内镜像
npm config set registry https://registry.npmmirror.com

# 清除缓存重新安装
rush clean
rush update
```

**Q: 后端启动报错 "database connection failed"?**

A: 检查 MySQL 是否启动

```bash
# 检查 Docker 服务
docker compose ps

# 重启 MySQL
docker compose restart mysql

# 检查数据库连接
mysql -h127.0.0.1 -P3306 -uroot -p
```

### 7.2 开发问题

**Q: 前端页面空白?**

A: 检查控制台错误

```bash
# 打开浏览器开发者工具 (F12)
# 查看 Console 标签页的错误信息
# 查看 Network 标签页的 API 请求是否正常
```

**Q: API 请求跨域?**

A: 配置代理

```typescript
// vite.config.ts
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8001',
        changeOrigin: true,
      },
    },
  },
});
```

**Q: 组件导入报错?**

A: 检查包路径

```typescript
// ❌ 错误
import { Button } from './button';

// ✅ 正确
import { Button } from '@coze-studio/ui';
```

### 7.3 调试技巧

**前端调试**:

```typescript
// 使用 console.log
console.log('Debug:', data);

// 使用 React DevTools
// 安装 React Developer Tools 浏览器扩展

// 使用 debugger
debugger;  // 代码会在这里暂停
```

**后端调试**:

```go
// 使用 fmt.Println
fmt.Printf("Debug: %+v\n", bot)

// 使用 Delve 调试器
dlv debug cmd/main.go

// 设置断点
// 在代码中添加 breakpoint 标记
runtime.Breakpoint()
```

---

## 快速参考

### 8.1 常用命令

```bash
# 基础设施
make middleware      # 启动所有中间件
make stop           # 停止所有中间件
make clean          # 清理中间件数据

# 后端
make server         # 启动后端服务
make build_server   # 构建后端
make test_backend   # 运行后端测试

# 前端
make fe             # 启动前端服务
make build_fe      # 构建前端
make test_fe        # 运行前端测试
make lint_fe        # 前端代码检查

# 数据库
make sync_db        # 同步数据库结构
make dump_db        # 导出数据库结构
make reset_db       # 重置数据库

# 全部
make debug          # 启动完整开发环境 (前端+后端+中间件)
make build          # 构建所有
make test           # 运行所有测试
make lint           # 运行所有检查
```

### 8.2 默认账号密码

**系统管理员**:
- 用户名: `admin`
- 密码: `admin123`
- 权限: 超级管理员

**数据库**:
- 用户名: `root`
- 密码: `root123456`
- 数据库: `zker`

**Redis**:
- 密码: (无密码)

### 8.3 端口清单

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 8888 | React 开发服务器 |
| 后端 API | 8001 | Go HTTP 服务 |
| MySQL | 3306 | 数据库 |
| Redis | 6379 | 缓存 |
| Elasticsearch | 9200 | 搜索引擎 |
| Milvus | 19530 | 向量数据库 |
| MinIO Console | 9001 | 对象存储管理 |
| etcd | 2379 | 服务发现 |
| NSQ Lookupd | 4160 | 消息队列查找 |
| NSQD | 4150 | 消息队列 |

### 8.4 关键文档链接

- [统一错误码定义](./ZKER-统一错误码定义规范.md)
- [技术组件清单与使用指南](./ZKER-技术组件清单与使用指南(完整版).md)
- [核心算法实现指南](./ZKER-核心算法实现指南.md)
- [OpenAPI 规范](./openapi/zker-api-v1-core-modules.yaml)
- [数据库设计清单](./数据库设计完整交付清单.md)

---

## 下一步

🎉 **恭喜你完成了快速入门!**

接下来建议:

1. **深入学习**: 阅读详细设计文档,理解系统架构
2. **实践项目**: 尝试修复一个简单的 Bug 或添加一个小功能
3. **代码审查**: 参与 Pull Request 代码审查,学习最佳实践
4. **团队协作**: 参加团队 Daily Standup,了解项目进展

**需要帮助?**

- 查看 [常见问题](#常见问题)
- 联系导师或团队负责人
- 在团队群里提问

---

**文档维护**: ZKER 架构团队
**更新频率**: 每季度更新
**反馈渠道**: 提交 Issue 或 PR
