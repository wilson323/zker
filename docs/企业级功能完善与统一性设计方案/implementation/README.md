# ZKER Enterprise - 开发者实施指南

**文档版本**: v1.0.0
**创建日期**: 2025-12-29
**目标读者**: 开发工程师、DevOps工程师
**预计阅读时间**: 30分钟

---

## 📋 文档说明

本文档面向参与ZKER Enterprise平台开发的所有工程师，提供从环境搭建到代码部署的完整实施指南。

---

## 目录

1. [快速开始](#1-快速开始)
2. [开发环境搭建](#2-开发环境搭建)
3. [项目结构详解](#3-项目结构详解)
4. [开发工作流程](#4-开发工作流程)
5. [代码规范](#5-代码规范)
6. [测试指南](#6-测试指南)
7. [部署流程](#7-部署流程)
8. [常见问题FAQ](#8-常见问题faq)
9. [最佳实践](#9-最佳实践)

---

## 1. 快速开始

### 1.1 前置要求

**硬件要求**:
- CPU: 4核以上
- 内存: 16GB以上
- 硬盘: 100GB以上SSD

**软件要求**:
- 操作系统: macOS 12+, Ubuntu 20.04+, Windows 10/11 (WSL2)
- Docker: 20.10+
- Docker Compose: 2.0+
- Go: 1.23+
- Node.js: 18.0+
- Rush: 全局安装

### 1.2 克隆项目

```bash
# 克隆代码仓库
git clone https://github.com/zker/zker-enterprise.git
cd zker-enterprise

# 安装Rush (如果还没安装)
npm install -g @microsoft/rush

# 安装依赖
rush update

# 构建项目
rush build
```

### 1.3 启动开发环境

```bash
# 启动中间件服务 (MySQL, Redis, ES等)
cd docker
cp .env.example .env
docker compose up -d

# 等待服务就绪...
# 检查服务状态
docker compose ps

# 启动Go后端 (新终端)
make server

# 启动前端开发服务器 (新终端)
cd frontend/apps/coze-studio
npm run dev
```

### 1.4 访问应用

- **前端应用**: http://localhost:8888
- **API文档**: http://localhost:8080/swagger
- **Grafana监控**: http://localhost:3000 (admin/admin)
- **Kibana日志**: http://localhost:5601

---

## 2. 开发环境搭建

### 2.1 Docker服务栈

项目使用Docker Compose管理所有中间件服务：

```yaml
# docker-compose.yaml 主要服务
services:
  mysql:
    image: mysql:8.4.5
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: zker

  redis:
    image: redis:8.0
    ports:
      - "6379:6379"

  elasticsearch:
    image: elasticsearch:8.18.0
    ports:
      - "9200:9200"
    environment:
      - "discovery.type=single-node"
      - "xpack.security.enabled=false"

  milvus:
    image: milvusdb/milvus:v2.5.10
    ports:
      - "19530:19530"

  nsqlookupd:
    image: nsqio/nsq:v1.2.1
    command: /nsqlookupd

  nsqd:
    image: nsqio/nsq:v1.2.1
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160

  etcd:
    image: quay.io/coreos/etcd:v3.5.0
```

**启动所有服务**:

```bash
cd docker
docker compose up -d

# 查看日志
docker compose logs -f

# 停止所有服务
docker compose down

# 重启单个服务
docker compose restart mysql
```

### 2.2 Go开发环境

**安装Go工具**:

```bash
# 安装Go 1.23
# macOS
brew install go

# Ubuntu
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Windows
# 下载安装包: https://go.dev/dl/

# 验证安装
go version
```

**配置Go环境变量**:

```bash
# ~/.bashrc 或 ~/.zshrc
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
```

**安装Go开发工具**:

```bash
# 安装golangci-lint (代码检查)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装air (热重载)
go install github.com/cosmtrek/air@latest

# 安装mockgen (Mock生成)
go install github.com/golang/mock/mockgen@latest
```

### 2.3 前端开发环境

**安装Node.js和Rush**:

```bash
# 安装Node.js 18 LTS
# macOS
brew install node@18

# Ubuntu
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 安装Rush
npm install -g @microsoft/rush

# 验证安装
node --version
rush --version
```

**配置前端开发工具**:

```bash
# 在项目根目录
rush update

# 安装前端依赖 (自动完成)
```

### 2.4 IDE配置

**推荐IDE**: VSCode, GoLand, WebStorm

**VSCode插件推荐**:

```json
// .vscode/extensions.json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-typescript-next",
    "dbaeumer.vscode-eslint",
    "esbenp.prettier-vscode",
    "eamodio.gitlens",
    "ms-azuretools.vscode-docker",
    "christian-kohler.path-intellisense"
  ]
}
```

**VSCode配置**:

```json
// .vscode/settings.json
{
  "go.useLanguageServer": true,
  "go.toolsManagement.autoUpdate": true,
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": "explicit"
  },
  "typescript.tsdk": "node_modules/typescript/lib"
}
```

---

## 3. 项目结构详解

### 3.1 整体目录结构

```
zker-enterprise/
├── backend/                    # Go后端
│   ├── api/                    # HTTP接口层
│   │   ├── router/             # 路由定义
│   │   └── handler/            # 处理器
│   ├── application/            # 应用层
│   │   ├── dto/                # 数据传输对象
│   │   └── service/            # 应用服务
│   ├── domain/                 # 领域层
│   │   ├── entity/             # 实体
│   │   ├── repository/         # 仓储接口
│   │   └── service/            # 领域服务
│   ├── infra/                  # 基础设施层
│   │   ├── db/                 # 数据库实现
│   │   ├── cache/              # 缓存实现
│   │   ├── mq/                 # 消息队列实现
│   │   └── es/                 # Elasticsearch实现
│   ├── config/                 # 配置
│   ├── pkg/                    # 工具包
│   └── main.go                 # 入口文件
│
├── frontend/                   # React前端 (Rush Monorepo)
│   ├── arch/                   # 架构层 (Level 1)
│   ├── common/                 # 公共层 (Level 2)
│   ├── agent-ide/              # Agent IDE (Level 3)
│   ├── workflow/               # 工作流 (Level 3)
│   └── apps/                   # 应用层 (Level 4)
│       └── coze-studio/        # 主应用
│           ├── src/
│           │   ├── pages/       # 页面组件
│           │   ├── components/  # 业务组件
│           │   ├── hooks/       # 自定义Hooks
│           │   ├── store/       # 状态管理
│           │   └── utils/       # 工具函数
│           └── package.json
│
├── docker/                     # Docker配置
│   ├── docker-compose.yaml
│   └── .env.example
│
├── docs/                       # 文档
│   └── 企业级功能完善与统一性设计方案/
│       ├── openapi/
│       │   └── openapi.yaml
│       ├── database/
│       │   └── schema.sql
│       └── uml/
│           └── sequences.md
│
├── scripts/                    # 脚本工具
│   ├── migration.sh            # 数据库迁移
│   ├── deploy.sh               # 部署脚本
│   └── test.sh                 # 测试脚本
│
├── Makefile                    # Make命令
└── README.md
```

### 3.2 后端分层架构

```
backend/
├── api/                        # 接口层
│   └── router/
│       └── v1/
│           ├── bot.go          # Bot路由
│           ├── conversation.go # 对话路由
│           └── auth.go         # 认证路由
│
├── application/                # 应用层
│   ├── service/
│   │   ├── bot_service.go     # Bot应用服务
│   │   └── auth_service.go    # 认证应用服务
│   └── dto/
│       ├── bot_request.go      # Bot请求DTO
│       └── bot_response.go     # Bot响应DTO
│
├── domain/                     # 领域层
│   ├── entity/
│   │   ├── bot.go             # Bot实体
│   │   ├── conversation.go    # 对话实体
│   │   └── user.go            # 用户实体
│   ├── repository/
│   │   ├── bot_repository.go  # Bot仓储接口
│   │   └── user_repository.go # 用户仓储接口
│   └── service/
│       ├── bot_domain_service.go    # Bot领域服务
│       └── auth_domain_service.go   # 认证领域服务
│
├── infra/                      # 基础设施层
│   ├── db/
│   │   ├── mysql.go           # MySQL实现
│   │   └── repository/
│   │       ├── bot_repo_impl.go
│   │       └── user_repo_impl.go
│   ├── cache/
│   │   └── redis.go           # Redis实现
│   └── mq/
│       └── nsq.go             # NSQ实现
│
└── main.go
```

### 3.3 前端模块结构

```
frontend/apps/coze-studio/src/
├── pages/                      # 页面
│   ├── login/
│   ├── bots/
│   │   ├── BotList.tsx
│   │   ├── BotDetail.tsx
│   │   └── BotBuilder.tsx
│   └── conversations/
│
├── components/                 # 组件
│   ├── common/                # 通用组件
│   │   ├── Button/
│   │   ├── Modal/
│   │   └── Table/
│   ├── bot/                   # Bot组件
│   │   ├── BotCard.tsx
│   │   └── BotConfig.tsx
│   └── conversation/          # 对话组件
│       ├── ChatBox.tsx
│       └── MessageList.tsx
│
├── hooks/                      # 自定义Hooks
│   ├── useAuth.ts
│   ├── useBot.ts
│   └── useConversation.ts
│
├── store/                      # 状态管理
│   ├── authStore.ts
│   ├── botStore.ts
│   └── conversationStore.ts
│
├── services/                   # API服务
│   ├── api.ts                 # API客户端
│   ├── botService.ts
│   └── authService.ts
│
└── utils/                      # 工具函数
    ├── request.ts
    ├── format.ts
    └── validate.ts
```

---

## 4. 开发工作流程

### 4.1 分支管理

```bash
# 主分支
main (生产环境)

# 开发分支
develop (开发环境)

# 功能分支
feature/bot-builder-enhancement
feature/rag-optimization

# 修复分支
hotfix/login-bug
hotfix/payment-issue

# 发布分支
release/v1.0.0
```

**分支工作流**:

```bash
# 1. 从develop创建功能分支
git checkout develop
git checkout -b feature/new-bot-feature

# 2. 开发并提交
git add .
git commit -m "feat: add new bot feature"

# 3. 推送到远程
git push origin feature/new-bot-feature

# 4. 创建Pull Request到develop
# 5. 代码审查通过后合并
# 6. 删除功能分支
git branch -d feature/new-bot-feature
```

### 4.2 提交规范

使用Conventional Commits规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型(type)**:
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具链相关

**示例**:

```bash
# 新功能
git commit -m "feat(bot): add bot cloning feature"

# Bug修复
git commit -m "fix(auth): resolve token refresh issue"

# 文档更新
git commit -m "docs(api): update OpenAPI specification"

# 重构
git commit -m "refactor(conversation): extract message service"
```

### 4.3 代码审查流程

```bash
# 1. 创建Pull Request
gh pr create --title "feat: add bot cloning feature" --body "PR description"

# 2. 请求审查
gh pr edit --add-reviewer @reviewer1 --add-reviewer @reviewer2

# 3. 修改后更新PR
git add .
git commit -m "fix: address review comments"
git push

# 4. 审查通过后合并
gh pr merge --squash
```

### 4.4 版本发布流程

```bash
# 1. 创建发布分支
git checkout develop
git checkout -b release/v1.0.0

# 2. 更新版本号
# backend/config/version.go
# frontend/package.json

# 3. 创建Tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 4. 合并到main并推送
git checkout main
git merge release/v1.0.0
git push origin main --tags

# 5. 部署到生产环境
make deploy
```

---

## 5. 代码规范

### 5.1 Go代码规范

**文件命名**:
- 文件名: `snake_case`
- 测试文件: `{filename}_test.go`
- 包名: `lowercase`

**代码格式化**:

```bash
# 格式化所有Go代码
go fmt ./...

# 或使用gofmt
gofmt -w -s .

# 检查代码
golangci-lint run
```

**命名规范**:

```go
// 包名: 小写单词
package user

// 接口名: 大写驼峰
type UserService interface {}

// 结构体: 大写驼峰
type User struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Email       string `json:"email"`
    CreatedAt   time.Time `json:"created_at"`
}

// 常量: 大写驼峰或全大写
const MaxRetryCount = 3
const API_BASE_URL = "https://api.zker.com"

// 私有变量: 小写驼峰
var privateVar string

// 公开变量: 大写驼峰
var PublicVar string

// 函数: 大写驼峰 (公开) / 小写驼峰 (私有)
func GetUser(id int64) (*User, error) {}
func createUser(user *User) error {}
```

**错误处理**:

```go
// 总是检查错误
user, err := userService.GetUser(ctx, userID)
if err != nil {
    return nil, fmt.Errorf("failed to get user: %w", err)
}

// 使用自定义错误类型
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

// 包装错误上下文
if err != nil {
    return fmt.Errorf("create user failed: %w", err)
}
```

### 5.2 TypeScript代码规范

**文件命名**:
- 组件文件: `PascalCase.tsx`
- 工具文件: `camelCase.ts`
- 类型文件: `camelCase.types.ts`
- Hook文件: `use{Purpose}.ts`

**代码格式化**:

```bash
# 格式化所有TypeScript代码
npm run format
# 或
prettier --write "src/**/*.ts"
```

**命名规范**:

```typescript
// 组件: 大写驼峰
function BotCard() {}
const BotConfig = () => {};

// 接口/类型: 大写驼峰
interface Bot {
  id: string;
  name: string;
}
type BotStatus = 'active' | 'inactive';

// 变量: 小写驼峰
const botName = 'My Bot';
let isActive = true;

// 常量: 全大写蛇形
const API_BASE_URL = 'https://api.zker.com';

// 私有成员: 下划线前缀
class BotService {
  private _httpClient: HttpClient;
}
```

**组件规范**:

```tsx
// 函数组件示例
import { useState, useEffect } from 'react';

interface BotCardProps {
  bot: Bot;
  onEdit: (bot: Bot) => void;
}

export function BotCard({ bot, onEdit }: BotCardProps) {
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // 副作用逻辑
  }, []);

  const handleEdit = () => {
    onEdit(bot);
  };

  return (
    <div className="bot-card">
      <h3>{bot.name}</h3>
      <button onClick={handleEdit}>编辑</button>
    </div>
  );
}

export default BotCard;
```

### 5.3 API规范

**RESTful API设计**:

```
GET    /api/v1/bots           # 获取Bot列表
POST   /api/v1/bots           # 创建Bot
GET    /api/v1/bots/{id}      # 获取Bot详情
PUT    /api/v1/bots/{id}      # 更新Bot
DELETE /api/v1/bots/{id}      # 删除Bot
```

**请求示例**:

```bash
# GET请求
curl -X GET "http://localhost:8080/api/v1/bots?page=1&page_size=20" \
  -H "Authorization: Bearer {token}" \
  -H "X-Tenant-ID: tenant_001"

# POST请求
curl -X POST "http://localhost:8080/api/v1/bots" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "客服Bot",
    "description": "智能客服助手",
    "type": "chat"
  }'
```

**响应格式**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "bot_123",
    "name": "客服Bot",
    "status": "active"
  }
}
```

---

## 6. 测试指南

### 6.1 后端测试

**单元测试**:

```go
// service_test.go示例
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestBotService_CreateBot(t *testing.T) {
    // Arrange
    mockRepo := new(MockBotRepository)
    service := NewBotService(mockRepo)
    bot := &entity.Bot{
        Name: "Test Bot",
        Type: "chat",
    }

    // Act
    err := service.CreateBot(context.Background(), bot)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, bot.ID)
}
```

**运行测试**:

```bash
# 运行所有测试
go test ./...

# 运行指定包的测试
go test ./application/service/

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool html -coverage.out -o coverage.html

# 查看覆盖率
go test -cover ./...
```

### 6.2 前端测试

**组件测试**:

```tsx
// BotCard.test.tsx示例
import { render, screen } from '@testing-library/react';
import BotCard from './BotCard';

describe('BotCard', () => {
  const mockBot = {
    id: 'bot_123',
    name: 'Test Bot',
    status: 'active',
  };

  test('renders bot name', () => {
    render(<BotCard bot={mockBot} onEdit={() => {}} />);
    expect(screen.getByText('Test Bot')).toBeInTheDocument();
  });

  test('calls onEdit when edit button is clicked', () => {
    const mockOnEdit = jest.fn();
    render(<BotCard bot={mockBot} onEdit={mockOnEdit} />);

    screen.getByText('编辑').click();
    expect(mockOnEdit).toHaveBeenCalledWith(mockBot);
  });
});
```

**运行测试**:

```bash
# 运行所有测试
npm test

# 运行指定文件的测试
npm test BotCard.test.tsx

# 生成覆盖率报告
npm run test:cov

# 交互式测试
npm test -- --coverage --watchAll=false
```

### 6.3 集成测试

**E2E测试** (Playwright):

```typescript
// e2e/bot.spec.ts
import { test, expect } from '@playwright/test';

test.describe('Bot Management', () => {
  test.beforeEach(async ({ page }) => {
    // 登录
    await page.goto('http://localhost:8888/login');
    await page.fill('input[name="email"]', 'admin@zker.com');
    await page.fill('input[name="password"]', 'password123');
    await page.click('button[type="submit"]');
    await expect(page).toHaveURL('http://localhost:8888/');
  });

  test('should create a new bot', async ({ page }) => {
    await page.goto('http://localhost:8888/bots');
    await page.click('button:has-text("新建Bot")');

    await page.fill('input[name="name"]', 'Test Bot');
    await page.fill('textarea[name="description"]', 'Test Description');
    await page.click('button:has-text("创建")');

    await expect(page.locator('.notification-success')).toContainText('Bot创建成功');
  });
});
```

---

## 7. 部署流程

### 7.1 本地部署

```bash
# 1. 构建前端
cd frontend/apps/coze-studio
npm run build

# 2. 构建后端
cd backend
go build -o bin/server ./main.go

# 3. 启动服务
./bin/server

# 4. 或使用Make命令
make build
make server
```

### 7.2 Docker部署

```bash
# 构建镜像
docker build -t zker/backend:latest ./backend
docker build -t zker/frontend:latest ./frontend/apps/coze-studio

# 推送到镜像仓库
docker tag zker/backend:latest registry.zker.com/zker/backend:latest
docker push registry.zker.com/zker/backend:latest

# 使用Docker Compose部署
cd docker
docker compose -f docker-compose.prod.yaml up -d
```

### 7.3 Kubernetes部署

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zker-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: zker-backend
  template:
    metadata:
      labels:
        app: zker-backend
    spec:
      containers:
      - name: backend
        image: registry.zker.com/zker/backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: mysql-service
        - name: REDIS_HOST
          value: redis-service
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: zker-backend
spec:
  selector:
    app: zker-backend
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

**部署到K8s**:

```bash
# 部署到K8s
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# 查看部署状态
kubectl get pods
kubectl logs -f deployment/zker-backend

# 扩容
kubectl scale deployment zker-backend --replicas=5
```

---

## 8. 常见问题FAQ

### 8.1 环境问题

**Q1: Docker服务启动失败?**

```bash
# 检查Docker服务状态
docker ps -a

# 查看服务日志
docker compose logs mysql
docker compose logs elasticsearch

# 重启服务
docker compose restart mysql
```

**Q2: 前端依赖安装失败?**

```bash
# 清理缓存重新安装
rush update --purge
rush update

# 或手动清理
rm -rf node_modules
rm -rf common/node_modules
rush update
```

**Q3: Go模块下载慢?**

```bash
# 使用国内代理
go env -w GOPROXY=https://goproxy.cn,direct

# 或使用Go Proxy
export GOPROXY=https://goproxy.io,direct
```

### 8.2 开发问题

**Q4: 前端热重载不生效?**

```bash
# 检查Rspack配置
cat frontend/apps/coze-studio/rsbuild.config.ts

# 重启开发服务器
cd frontend/apps/coze-studio
npm run dev
```

**Q5: 数据库连接失败?**

```bash
# 检查MySQL服务
docker ps | grep mysql

# 检查连接配置
cat docker/.env | grep MYSQL

# 测试连接
docker exec -it mysql zker -uroot -proot123
```

**Q6: API返回401 Unauthorized?**

```bash
# 检查Token是否有效
curl -H "Authorization: Bearer {token}" http://localhost:8080/api/v1/auth/validate

# 检查租户ID
curl -H "X-Tenant-ID: tenant_001" http://localhost:8080/api/v1/bots
```

---

## 9. 最佳实践

### 9.1 性能优化

**数据库优化**:

```sql
-- 添加索引
CREATE INDEX idx_tenant_bot_date ON bot_metrics(tenant_id, bot_id, metric_date);

-- 分区表(按月分区)
ALTER TABLE token_usage_metrics PARTITION BY RANGE (YEAR(created_at) * 100 + MONTH(created_at));

-- 慢查询优化
EXPLAIN SELECT * FROM bots WHERE tenant_id = 'tenant_001';
```

**缓存策略**:

```go
// Redis缓存示例
func (s *BotService) GetBot(ctx context.Context, botID string) (*entity.Bot, error) {
    // 1. 尝试从缓存获取
    cacheKey := fmt.Sprintf("bot:%s", botID)
    var bot entity.Bot
    err := s.cache.Get(ctx, cacheKey, &bot)
    if err == nil {
        return &bot, nil
    }

    // 2. 缓存未命中，从数据库查询
    bot, err = s.repository.GetByID(ctx, botID)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存 (TTL: 1小时)
    s.cache.Set(ctx, cacheKey, bot, time.Hour)

    return bot, nil
}
```

### 9.2 安全最佳实践

**SQL注入防护**:

```go
// 使用参数化查询
query := "SELECT * FROM users WHERE id = ?"
err := db.QueryRow(query, userID).Scan(&user)
```

**XSS防护**:

```tsx
// React默认转义，避免使用dangerouslySetInnerHTML
function MessageContent({ content }: { content: string }) {
  return <p>{content}</p>; // 自动转义
}
```

**CSRF防护**:

```go
// 中间件检查CSRF Token
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("X-CSRF-Token")
        if !validateCSRFToken(token) {
            c.JSON(403, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 9.3 可观测性

**日志规范**:

```go
// 结构化日志
logger.Info("Bot created",
    zap.String("bot_id", bot.ID),
    zap.String("tenant_id", bot.TenantID),
    zap.String("creator_id", bot.CreatorID),
)

logger.Error("Failed to create bot",
    zap.String("error", err.Error()),
    zap.String("tenant_id", tenantID),
)
```

**指标采集**:

```go
// Prometheus指标
var botCreationTotal = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "bot_creation_total",
        Help: "Total number of bots created",
    },
    []string{"tenant_id"},
)

botCreationTotal.WithLabelValues(tenantID).Inc()
```

---

## 📚 扩展阅读

- [Go最佳实践](https://go.dev/doc/effective_go)
- [React最佳实践](https://react.dev/learn)
- [Docker最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [Kubernetes最佳实践](https://kubernetes.io/docs/concepts/configuration/overview/)

---

## 🆘 获取帮助

遇到问题？联系我们：

- **技术支持**: support@zker.com
- **开发社区**: https://community.zker.com
- **问题追踪**: https://github.com/zker/zker-enterprise/issues

---

**文档维护**: ZKER Enterprise Team
**最后更新**: 2025-12-29
**版本**: v1.0.0
