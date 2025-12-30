# ZKER 快速开始指南

**版本**: v3.0.0 | **更新**: 2025-01-03

---

## 目录

- [环境准备](#环境准备)
- [快速启动](#快速启动)
- [开发指南](#开发指南)
- [常见问题](#常见问题)
- [下一步](#下一步)

---

## 环境准备

### 系统要求

| 组件 | 最低要求 | 推荐配置 |
|------|---------|---------|
| **操作系统** | Windows 10 / macOS 11 / Ubuntu 20.04 | Windows 11 / macOS 13 / Ubuntu 22.04 |
| **CPU** | 4 核 | 8 核+ |
| **内存** | 8 GB | 16 GB+ |
| **磁盘** | 20 GB | 50 GB+ SSD |

### 软件依赖

#### 后端开发

```bash
# 检查 Go 版本（需要 1.24+）
go version

# 检查 Docker 版本（需要 20.10+）
docker --version
docker compose version

# 检查 Make 版本
make --version
```

#### 前端开发

```bash
# 检查 Node.js 版本（需要 18+）
node --version

# 检查 npm 版本（需要 9+）
npm --version

# 检查 Rush 版本（需要 5.80+）
rush --version
```

### 安装依赖

#### 安装 Go 1.24

**macOS**:
```bash
brew install go
```

**Linux**:
```bash
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**Windows**:
下载安装包: https://go.dev/dl/

#### 安装 Node.js 18

**macOS/Linux**:
```bash
# 使用 nvm 安装
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 18
nvm use 18
```

**Windows**:
下载安装包: https://nodejs.org/

#### 安装 Docker

**macOS/Windows**:
下载 Docker Desktop: https://www.docker.com/products/docker-desktop/

**Linux**:
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

#### 安装 Rush

```bash
npm install -g @microsoft/rush
```

---

## 快速启动

### 1. 克隆项目

```bash
git clone https://github.com/coze-dev/coze-studio.git
cd coze-studio
```

### 2. 启动中间件

```bash
cd docker
docker compose up -d

# 检查服务状态
docker compose ps
```

等待以下服务启动：
- MySQL (端口 3306)
- Redis (端口 6379)
- Elasticsearch (端口 9200)
- MinIO (端口 9000)

### 3. 配置后端

```bash
cd ../backend

# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件
vim configs/config.yaml
```

**配置文件示例**:

```yaml
# 服务器配置
server:
  port: 8888
  mode: debug  # debug/release

# 数据库配置
mysql:
  host: localhost
  port: 3306
  user: root
  password: password
  database: zker
  charset: utf8mb4

# Redis 配置
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

# Elasticsearch 配置
elasticsearch:
  hosts: ["http://localhost:9200"]
  username: ""
  password: ""

# 日志配置
log:
  level: info
  format: json
  output: stdout
```

### 4. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE zker CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行迁移脚本
go run cmd/migrate/main.go

# 导入初始数据
mysql -u root -p zker < docs/企业级功能完善与统一性设计方案/database_init_data.sql
```

### 5. 启动后端

```bash
# 安装依赖
go mod download

# 运行
make server

# 或直接运行
go run cmd/main.go
```

后端服务启动在: http://localhost:8888

### 6. 配置前端

```bash
cd ../frontend

# 安装 Rush
npm install -g @microsoft/rush

# 安装依赖
rush update

# 构建
rush build
```

### 7. 启动前端

```bash
cd apps/coze-studio

# 复制环境变量
cp .env.example .env.local

# 编辑环境变量
vim .env.local
```

**环境变量示例**:

```env
# API 地址
VITE_API_BASE_URL=http://localhost:8888

# WebSocket 地址
VITE_WS_BASE_URL=ws://localhost:8888

# 应用配置
VITE_APP_NAME=ZKER
VITE_APP_PORT=8888
```

```bash
# 启动开发服务器
npm run dev
```

前端服务启动在: http://localhost:8888

### 8. 访问应用

打开浏览器访问: http://localhost:8888

默认账号：
- 用户名: `admin`
- 密码: `admin123`

---

## 开发指南

### 后端开发

#### 创建新的 API

```bash
# 1. 在 domain/ 创建领域实体
# 例如: domain/tenant/entity/tenant.go

# 2. 在 domain/ 创建仓储接口
# 例如: domain/tenant/repository/tenant_repository.go

# 3. 在 application/ 创建应用服务
# 例如: application/tenant/tenant_application.go

# 4. 在 api/handler/ 创建处理器
# 例如: api/handler/coze/tenant_service.go

# 5. 在 api/router/ 注册路由
# 例如: api/router/coze/api.go
```

**示例**:

```go
// domain/tenant/entity/tenant.go
package entity

import "time"

type Tenant struct {
    TenantID   string    `gorm:"primaryKey"`
    TenantName string    `gorm:"uniqueIndex;size:100;not null"`
    Status     string    `gorm:"size:20;not null;default:'active'"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// api/handler/coze/tenant_service.go
package coze

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

type TenantHandler struct {
    tenantApp *application.TenantApplication
}

func (h *TenantHandler) CreateTenant(ctx context.Context, c *app.RequestContext) {
    var req CreateTenantRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, ErrorResponse(err))
        return
    }

    tenant, err := h.tenantApp.CreateTenant(ctx, &req)
    if err != nil {
        c.JSON(500, ErrorResponse(err))
        return
    }

    c.JSON(201, SuccessResponse(tenant))
}
```

#### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./domain/tenant/...

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### Lint 检查

```bash
# 运行 golangci-lint
golangci-lint run

# 自动修复
golangci-lint run --fix
```

### 前端开发

#### 创建新页面

```bash
# 1. 在对应的业务包中创建页面组件
# 例如: frontend/packages/studio/pages/tenant/TenantList.tsx

# 2. 创建路由配置
# 例如: frontend/apps/coze-studio/src/routes/tenant.ts

# 3. 在主路由中注册
# 例如: frontend/apps/coze-studio/src/routes/index.tsx
```

**示例**:

```typescript
// frontend/packages/studio/pages/tenant/TenantList.tsx
import React, { useState, useEffect } from 'react';
import { Table } from '@douyinfe/semi-ui';
import { useTenants } from '@/hooks/useTenants';

export const TenantList: React.FC = () => {
  const { data, loading, fetchTenants } = useTenants();

  useEffect(() => {
    fetchTenants();
  }, [fetchTenants]);

  return (
    <Table
      dataSource={data}
      loading={loading}
      columns={[
        { title: '租户名称', dataIndex: 'tenant_name' },
        { title: '状态', dataIndex: 'status' },
      ]}
    />
  );
};

// frontend/apps/coze-studio/src/routes/tenant.ts
export const tenantRoutes = [
  {
    path: '/tenant',
    component: () => import('@/pages/tenant/TenantList'),
  },
];
```

#### 运行测试

```bash
# 运行所有测试
rush test

# 运行特定包的测试
cd frontend/packages/common
npm test
```

#### Lint 检查

```bash
# 运行 ESLint
rush lint

# 自动修复
rush lint --fix
```

---

## 常见问题

### 后端问题

#### Q: MySQL 连接失败？

**A**: 检查 MySQL 是否启动，配置是否正确。

```bash
# 检查 MySQL 状态
docker compose ps mysql

# 查看 MySQL 日志
docker compose logs mysql

# 测试连接
mysql -u root -p -h localhost -P 3306
```

#### Q: Redis 连接失败？

**A**: 检查 Redis 是否启动。

```bash
# 检查 Redis 状态
docker compose ps redis

# 查看 Redis 日志
docker compose logs redis

# 测试连接
redis-cli ping
```

#### Q: Go 依赖下载失败？

**A**: 使用 Go 代理。

```bash
# 设置 Go 代理
go env -w GOPROXY=https://goproxy.cn,direct

# 清理缓存
go clean -modcache

# 重新下载
go mod download
```

### 前端问题

#### Q: npm install 失败？

**A**: 使用 npm 镜像。

```bash
# 设置 npm 镜像
npm config set registry https://registry.npmmirror.com

# 清理缓存
npm cache clean --force

# 重新安装
rush update
```

#### Q: 构建失败？

**A**: 检查 Node.js 版本。

```bash
# 检查版本
node --version

# 切换到正确版本
nvm use 18

# 清理并重新构建
rush clean
rush rebuild
```

#### Q: 前端页面空白？

**A**: 检查控制台错误，确认 API 地址正确。

```bash
# 检查环境变量
cat .env.local

# 检查 API 是否可访问
curl http://localhost:8888/health
```

### Docker 问题

#### Q: Docker 容器无法启动？

**A**: 检查端口是否被占用。

```bash
# 检查端口占用
lsof -i :3306  # MySQL
lsof -i :6379  # Redis
lsof -i :9200  # Elasticsearch

# 停止占用端口的服务
sudo systemctl stop mysql  # 停止系统 MySQL

# 重新启动 Docker
docker compose up -d
```

#### Q: Docker 内存不足？

**A**: 增加 Docker 内存限制。

1. 打开 Docker Desktop
2. 进入 Settings > Resources > Advanced
3. 增加 Memory 到至少 8 GB
4. 重启 Docker

---

## 下一步

### 学习资源

| 文档 | 说明 |
|------|------|
| [architecture.md](architecture.md) | 系统架构详解 |
| [tech-stack.md](tech-stack.md) | 技术栈详解 |
| [02-SPECS/](../02-SPECS/) | 开发规范 |
| [03-DESIGN/](../03-DESIGN/) | 设计文档 |

### 开发任务

根据角色选择任务：

#### 后端开发（研发A/B）

- [ ] 实现租户管理 API
- [ ] 实现 RBAC 权限 API
- [ ] 实现智能路由引擎
- [ ] 编写单元测试

#### 前端开发（研发C）

- [ ] 实现租户管理页面
- [ ] 实现权限管理页面
- [ ] 实现路由配置页面
- [ ] 编写组件测试

#### DevOps（研发D）

- [ ] 配置 CI/CD 流程
- [ ] 配置监控告警
- [ ] 编写部署文档
- [ ] 配置备份策略

### 参与贡献

欢迎贡献代码、文档、Bug 报告！

- **Bug 报告**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)
- **功能建议**: [GitHub Discussions](https://github.com/coze-dev/coze-studio/discussions)
- **代码贡献**: [CONTRIBUTING.md](../../CONTRIBUTING.md)

---

**🎯 目标**: 快速上手，高效开发！
