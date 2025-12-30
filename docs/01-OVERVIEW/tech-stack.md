# ZKER 技术栈详解

**版本**: v3.0.0 | **更新**: 2025-01-03

---

## 目录

- [后端技术栈](#后端技术栈)
- [前端技术栈](#前端技术栈)
- [数据存储](#数据存储)
- [基础设施](#基础设施)
- [开发工具](#开发工具)
- [技术选型理由](#技术选型理由)

---

## 后端技术栈

### 编程语言

#### Go 1.24

**选择理由**:
- ✅ 高性能，适合高并发场景
- ✅ 强类型，编译时检查
- ✅ 简洁的并发模型（goroutine + channel）
- ✅ 丰富的标准库
- ✅ 快速编译，适合快速迭代

**核心特性应用**:

```go
// 1. 并发处理
func (s *Service) ProcessConversations(ctx context.Context, ids []string) error {
    errChan := make(chan error, len(ids))
    var wg sync.WaitGroup

    for _, id := range ids {
        wg.Add(1)
        go func(conversationID string) {
            defer wg.Done()
            if err := s.processOne(ctx, conversationID); err != nil {
                errChan <- err
            }
        }(id)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        if err != nil {
            return err
        }
    }
    return nil
}

// 2. 接口抽象
type TenantRepository interface {
    Create(ctx context.Context, tenant *Tenant) error
    FindByID(ctx context.Context, id string) (*Tenant, error)
    Update(ctx context.Context, tenant *Tenant) error
    Delete(ctx context.Context, id string) error
}

// 3. 错误处理
func (s *Service) CreateTenant(ctx context.Context, req *CreateRequest) (*Tenant, error) {
    if err := s.validate(req); err != nil {
        return nil, errorx.Wrapf(err, errno.InvalidArgument)
    }

    tenant, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, errorx.Wrapf(err, errno.TenantCreateFailed)
    }

    return tenant, nil
}
```

### Web 框架

#### Hertz (字节跳动开源)

**选择理由**:
- ✅ 高性能，基于 CloudWeGo 生态
- ✅ 企业级稳定性（字节内部大规模使用）
- ✅ 完善的中间件生态
- ✅ 优秀的性能表现

**核心特性**:

```go
// 路由定义
func main() {
    h := server.Default(server.WithHostPorts(":8888"))

    // 中间件
    h.Use(middleware.TenantIsolation())
    h.Use(middleware.PermissionCheck())
    h.Use(middleware.Session())

    // 路由组
    api := h.Group("/api/v1")
    {
        // 租户管理
        api.POST("/tenants", tenantHandler.Create)
        api.GET("/tenants", tenantHandler.List)
        api.GET("/tenants/:id", tenantHandler.Get)
        api.PUT("/tenants/:id", tenantHandler.Update)
        api.DELETE("/tenants/:id", tenantHandler.Delete)

        // 权限管理
        api.POST("/roles", roleHandler.Create)
        api.GET("/roles", roleHandler.List)
        api.PUT("/roles/:id", roleHandler.Update)
    }

    h.Spin()
}
```

**中间件开发**:

```go
// 中间件：租户隔离
func TenantIsolation() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 从会话中获取租户ID
        session := middleware.GetSession(c)
        if session == nil {
            c.JSON(401, ErrorResponse(errno.Unauthorized))
            c.Abort()
            return
        }

        // 2. 注入到上下文
        context.SetTenantID(ctx, session.TenantID)

        // 3. 继续处理
        c.Next(ctx)
    }
}
```

### ORM 框架

#### GORM

**选择理由**:
- ✅ 功能完善，API 简洁
- ✅ 支持多数据库
- ✅ Hook 机制完善
- ✅ 自动迁移

**核心用法**:

```go
// 模型定义
type Tenant struct {
    TenantID   string    `gorm:"primaryKey;size:36"`
    TenantName string    `gorm:"uniqueIndex;size:100;not null"`
    PlanType   string    `gorm:"size:20;not null;default:'free'"`
    Status     string    `gorm:"size:20;not null;default:'active'"`
    Quota      Quota     `gorm:"embedded;embeddedPrefix:quota_"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
    DeletedAt  *time.Time `gorm:"index"`
}

type Quota struct {
    MaxTokens     int64 `gorm:"default:1000000"`
    MaxAPICalls   int64 `gorm:"default:10000"`
    MaxStorageGB  int64 `gorm:"default:10"`
    UsedTokens    int64 `gorm:"default:0"`
    UsedAPICalls  int64 `gorm:"default:0"`
    UsedStorageGB int64 `gorm:"default:0"`
}

// 查询
func (r *TenantRepository) FindByID(ctx context.Context, id string) (*Tenant, error) {
    var tenant Tenant
    err := r.db.WithContext(ctx).
        Where("tenant_id = ? AND deleted_at IS NULL", id).
        First(&tenant).Error
    if err != nil {
        return nil, errorx.Wrapf(err, errno.TenantNotFound)
    }
    return &tenant, nil
}

// 创建
func (r *TenantRepository) Create(ctx context.Context, tenant *Tenant) error {
    return r.db.WithContext(ctx).Create(tenant).Error
}

// 更新
func (r *TenantRepository) Update(ctx context.Context, tenant *Tenant) error {
    return r.db.WithContext(ctx).Save(tenant).Error
}

// 软删除
func (r *TenantRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).
        Where("tenant_id = ?", id).
        Delete(&Tenant{}).Error
}
```

### 其他后端库

| 库 | 用途 | 版本 |
|------|------|------|
| **gin** | 路由（部分模块使用） | Latest |
| **wire** | 依赖注入 | Latest |
| **zap** | 日志 | Latest |
| **viper** | 配置管理 | Latest |
| **jwt-go** | JWT 认证 | Latest |
| **redis-go** | Redis 客户端 | Latest |
| **olivere/elastic** | Elasticsearch 客户端 | Latest |

---

## 前端技术栈

### 框架和语言

#### React 18 + TypeScript 5.x

**选择理由**:
- ✅ 生态丰富，组件库成熟
- ✅ Hooks 简化状态管理
- ✅ TypeScript 类型安全
- ✅ 社区活跃

**核心特性**:

```typescript
// 组件定义
import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Button, Table } from '@douyinfe/semi-ui';

interface Tenant {
  tenant_id: string;
  tenant_name: string;
  plan_type: string;
  status: string;
  quota: Quota;
}

interface Quota {
  max_tokens: number;
  used_tokens: number;
  max_api_calls: number;
  used_api_calls: number;
}

export const TenantList: React.FC = () => {
  // 状态管理
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    page: 1,
    pageSize: 10,
    total: 0,
  });

  // 数据获取
  const fetchTenants = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.getTenants({
        page: pagination.page,
        pageSize: pagination.pageSize,
      });
      setTenants(res.data);
      setPagination(prev => ({ ...prev, total: res.total }));
    } catch (error) {
      console.error('Failed to fetch tenants:', error);
    } finally {
      setLoading(false);
    }
  }, [pagination.page, pagination.pageSize]);

  useEffect(() => {
    fetchTenants();
  }, [fetchTenants]);

  // 计算属性
  const columns = useMemo(() => [
    {
      title: '租户名称',
      dataIndex: 'tenant_name',
      key: 'tenant_name',
    },
    {
      title: '套餐类型',
      dataIndex: 'plan_type',
      key: 'plan_type',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Badge type={status === 'active' ? 'success' : 'warning'}>
          {status === 'active' ? '激活' : '停用'}
        </Badge>
      ),
    },
    {
      title: '配额使用',
      key: 'quota',
      render: (_: any, record: Tenant) => (
        <Progress
          percent={(record.quota.used_tokens / record.quota.max_tokens) * 100}
          showInfo={false}
        />
      ),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Tenant) => (
        <Button onClick={() => handleEdit(record)}>编辑</Button>
      ),
    },
  ], []);

  return (
    <div className="tenant-list">
      <Table
        columns={columns}
        dataSource={tenants}
        loading={loading}
        pagination={pagination}
        onChange={setPagination}
      />
    </div>
  );
};
```

### UI 组件库

#### Semi Design (字节跳动)

**选择理由**:
- ✅ 设计规范统一
- ✅ 组件丰富（60+）
- ✅ TypeScript 支持
- ✅ 企业级稳定性

**核心组件**:

| 组件 | 用途 |
|------|------|
| **Layout** | 页面布局 |
| **Navigation** | 导航菜单 |
| **Table** | 数据表格 |
| **Form** | 表单 |
| **Modal** | 弹窗 |
| **Notification** | 通知 |
| **Button** | 按钮 |
| **Input** | 输入框 |

### 状态管理

#### Context API + Hooks

**选择理由**:
- ✅ 无需引入额外库
- ✅ React 内置，稳定
- ✅ 适合中小型应用

**示例**:

```typescript
// contexts/TenantContext.tsx
import React, { createContext, useContext, useState, ReactNode } from 'react';

interface Tenant {
  tenant_id: string;
  tenant_name: string;
}

interface TenantContextValue {
  currentTenant: Tenant | null;
  setCurrentTenant: (tenant: Tenant | null) => void;
}

const TenantContext = createContext<TenantContextValue | undefined>(undefined);

export const TenantProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [currentTenant, setCurrentTenant] = useState<Tenant | null>(null);

  return (
    <TenantContext.Provider value={{ currentTenant, setCurrentTenant }}>
      {children}
    </TenantContext.Provider>
  );
};

export const useTenant = () => {
  const context = useContext(TenantContext);
  if (!context) {
    throw new Error('useTenant must be used within TenantProvider');
  }
  return context;
};

// 使用
import { useTenant } from '@/contexts/TenantContext';

export const TenantHeader: React.FC = () => {
  const { currentTenant } = useTenant();

  return (
    <div>
      <h1>{currentTenant?.tenant_name}</h1>
    </div>
  );
};
```

### Monorepo 管理

#### Rush

**选择理由**:
- ✅ 微软开源，企业级
- ✅ 支持多种语言
- ✅ 增量构建
- ✅ 统一版本管理

**项目结构**:

```
frontend/
├── package.json          # Rush 配置
├── rush.json             # Rush 配置
├── common/
│   ├── package.json
│   └── ...
├── agent-ide/
│   ├── package.json
│   └── ...
└── apps/coze-studio/
    ├── package.json
    └── ...
```

---

## 数据存储

### MySQL 8.4.5

**选择理由**:
- ✅ 成熟稳定
- ✅ 支持事务
- ✅ 支持全文检索
- ✅ 丰富的索引类型

**配置优化**:

```sql
-- my.cnf
[mysqld]
# 连接数
max_connections = 1000

# 缓冲池大小（建议：物理内存的 70-80%）
innodb_buffer_pool_size = 4G

# 日志文件大小
innodb_log_file_size = 512M

# 刷盘策略（1：每秒刷盘，性能最好）
innodb_flush_log_at_trx_commit = 1

# 并发事务数
innodb_read_io_threads = 8
innodb_write_io_threads = 8

# 查询缓存
query_cache_size = 256M
query_cache_type = 1

# 字符集
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
```

### Redis 8.0

**选择理由**:
- ✅ 高性能
- ✅ 丰富的数据结构
- ✅ 持久化支持
- ✅ 集群支持

**使用场景**:

| 场景 | 数据结构 | 示例 |
|------|---------|------|
| **会话存储** | Hash | `session:{session_id}` |
| **缓存** | String | `tenant:{tenant_id}` |
| **分布式锁** | String | `lock:tenant:{tenant_id}` |
| **排行榜** | Sorted Set | `ranking:bot:{bot_id}` |
| **消息队列** | List | `queue:task` |

**配置优化**:

```conf
# redis.conf
# 内存
maxmemory 2gb
maxmemory-policy allkeys-lru

# 持久化
save 900 1
save 300 10
save 60 10000

# AOF
appendonly yes
appendfsync everysec

# 集群
cluster-enabled yes
cluster-config-file nodes.conf
```

### Elasticsearch 8.18

**选择理由**:
- ✅ 强大的全文检索
- ✅ 丰富的查询语法
- ✅ 聚合分析
- ✅ 分布式架构

**使用场景**:

| 场景 | 索引 | 字段 |
|------|------|------|
| **Bot 搜索** | bots | bot_name, description, tags |
| **对话搜索** | conversations | title, messages |
| **日志搜索** | logs | message, level, timestamp |

---

## 基础设施

### 容器化

#### Docker

**镜像构建**:

```dockerfile
# Dockerfile (后端)
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8888
CMD ["./main"]
```

```dockerfile
# Dockerfile (前端)
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### 编排

#### Docker Compose

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.4.5
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: zker
    volumes:
      - mysql-data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:8.0-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes

  elasticsearch:
    image: elasticsearch:8.18.0
    ports:
      - "9200:9200"
    environment:
      discovery.type: single-node
      ES_JAVA_OPTS: -Xms1g -Xmx1g
    volumes:
      - es-data:/usr/share/elasticsearch/data

  backend:
    build: ./backend
    ports:
      - "8888:8888"
    environment:
      MYSQL_HOST: mysql
      REDIS_HOST: redis
      ES_HOST: elasticsearch
    depends_on:
      - mysql
      - redis
      - elasticsearch

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  mysql-data:
  redis-data:
  es-data:
```

---

## 开发工具

### 代码质量

| 工具 | 用途 | 配置文件 |
|------|------|---------|
| **golangci-lint** | Go Lint | `.golangci.yml` |
| **ESLint** | JavaScript Lint | `.eslintrc.js` |
| **Prettier** | 代码格式化 | `.prettierrc` |
| **Husky** | Git Hooks | `.husky/` |

### 版本控制

#### Git

**分支策略**:

```
main          # 主分支，生产环境
├── develop   # 开发分支
│   ├── feature/xxx    # 功能分支
│   ├── bugfix/xxx     # 修复分支
│   └── hotfix/xxx     # 紧急修复分支
```

**提交规范**:

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型：
- feat: 新功能
- fix: 修复
- docs: 文档
- style: 格式
- refactor: 重构
- test: 测试
- chore: 构建/工具

### CI/CD

#### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.24
      - name: Lint
        run: golangci-lint run
      - name: Test
        run: go test ./... -cover

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Node
        uses: actions/setup-node@v3
        with:
          node-version: 18
      - name: Install
        run: npm ci
      - name: Lint
        run: npm run lint
      - name: Test
        run: npm test
```

---

## 技术选型理由

### 为什么选择 Go 而不是 Java/Node.js？

| 维度 | Go | Java | Node.js |
|------|-----|------|---------|
| **性能** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **并发** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **学习曲线** | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **生态** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **部署** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |

**结论**: Go 在性能和并发方面最优，适合高并发场景。

### 为什么选择 React 而不是 Vue/Angular？

| 维度 | React | Vue | Angular |
|------|-------|-----|---------|
| **生态** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **灵活性** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| **学习曲线** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ |
| **企业级** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

**结论**: React 生态最丰富，灵活性最高。

---

## 相关文档

| 文档 | 说明 |
|------|------|
| [architecture.md](architecture.md) | 系统架构详解 |
| [quick-start.md](quick-start.md) | 快速开始指南 |
| [02-SPECS/backend-dev-guide.md](../02-SPECS/backend-dev-guide.md) | 后端开发指南 |
| [02-SPECS/frontend-dev-guide.md](../02-SPECS/frontend-dev-guide.md) | 前端开发指南 |

---

**🎯 目标**: 选择最适合的技术栈，打造高性能、可维护的企业级应用！
