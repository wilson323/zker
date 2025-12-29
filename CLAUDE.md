# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

**📅 最后更新**: 2025-01-01
**🎯 企业级功能完善阶段**: 激活中

---

## 📖 快速导航

### 🚀 核心设计文档（企业级功能完善）
- **[实现差距分析与研发计划](./docs/企业级功能完善与统一性设计方案/ZKER-实现差距分析与研发计划_v1.0.md)** - 必读！识别12大差距，4人×8周详细计划
- **[企业级开发规范手册](./docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)** - 必读！200+页完整开发规范
- **[全局一致性检查清单](./docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md)** - 必读！4级检查体系，8大维度

### 👥 个人开发计划（4人并行）
- **[研发A - 后端架构师](./docs/企业级功能完善与统一性设计方案/研发A-后端架构师开发计划_v1.0.md)** - 多租户、RBAC、智能路由
- **[研发B - 后端工程师](./docs/企业级功能完善与统一性设计方案/研发B-后端工程师开发计划_v1.0.md)** - 错误码、性能测试、监控日志
- **[研发C - 前端工程师](./docs/企业级功能完善与统一性设计方案/研发C-前端工程师开发计划_v1.0.md)** - 组件库、业务页面、UX优化
- **[研发D - DevOps工程师](./docs/企业级功能完善与统一性设计方案/研发D-DevOps工程师开发计划_v1.0.md)** - Docker、K8s、CI/CD、灰度发布

### 📋 关键规范文档
- **[统一错误码定义规范](./docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)** - 300+错误码
- **[数据迁移方案](./docs/企业级功能完善与统一性设计方案/ZKER-数据迁移方案_v1.0.md)** - 零停机迁移
- **[灰度发布策略](./docs/企业级功能完善与统一性设计方案/ZKER-灰度发布策略_v1.0.md)** - 4阶段发布
- **[一键回滚方案](./docs/企业级功能完善与统一性设计方案/ZKER-一键回滚方案_v1.0.md)** - 5分钟快速回滚
- **[故障排查手册](./docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)** - 5大故障场景

### 🎓 技术参考文档
- **[技术组件清单与使用指南](./docs/企业级功能完善与统一性设计方案/ZKER-技术组件清单与使用指南(完整版).md)** - 100+页组件手册
- **[核心算法实现指南](./docs/企业级功能完善与统一性设计方案/ZKER-核心算法实现指南.md)** - 8大核心算法
- **[开发快速入门指南](./docs/企业级功能完善与统一性设计方案/ZKER-开发快速入门指南.md)** - 50+页快速入门

### 🏗️ 架构设计文档
- **[多租户SaaS架构](./docs/企业级功能完善与统一性设计方案/zker_MultiTenant_SaaS_完整架构设计文档.md)** - 完整架构设计
- **[数据库设计完整交付清单](./docs/企业级功能完善与统一性设计方案/数据库设计完整交付清单.md)** - 所有表结构
- **[API设计规范文档](./docs/企业级功能完善与统一性设计方案/API设计规范文档.md)** - RESTful规范

### 📚 完整文档索引
- **[文档索引](./docs/企业级功能完善与统一性设计方案/00-文档索引.md)** - 80+份文档完整索引

---

## 项目概述

Coze Studio 是一个一站式 AI Agent 开发平台，包含前端（React + TypeScript）和后端（Go）组件。项目采用 Rush.js 管理的复杂 monorepo 架构，拥有 135+ 个前端包，按层次依赖关系组织。

### 🎯 当前阶段：企业级功能完善

**核心目标**：将 ZKER 从单租户应用升级为企业级多租户 SaaS 平台

**关键任务**：
1. ✅ 多租户架构增强（tenant_id 迁移、订阅管理、配额强制）
2. ✅ RBAC 权限系统（5级数据权限 + 3级字段权限）
3. ✅ 智能路由引擎（混合意图匹配 + 评分路由）
4. ✅ 统一错误码系统（300+ 错误码，中英文双语）
5. ✅ 性能测试框架（JMeter/K6，完整性能基线）
6. ✅ 监控和日志（Prometheus + Grafana + ELK）
7. ✅ 数据迁移方案（双写 + 零停机迁移）
8. ✅ 灰度发布策略（4阶段发布，自动回滚）
9. ✅ 一键回滚系统（5分钟内快速回滚）
10. ✅ 故障排查手册（5大故障场景 SOP）

**研发团队**：4人并行开发，互不干扰
- 研发A（后端架构师）：多租户 + RBAC + 智能路由
- 研发B（后端工程师）：错误码 + 性能测试 + 监控
- 研发C（前端工程师）：组件库 + 业务页面 + UX优化
- 研发D（DevOps）：Docker + K8s + CI/CD + 灰度发布

**项目周期**：8周

---

## 🏗️ 架构概览

### 整体架构原则

#### 1. 核心设计原则（SOLID + KISS + DRY + YAGNI）

**SOLID 原则**:
- **S**ingle Responsibility：单一职责，每个组件只做一件事
- **O**pen/Closed：对扩展开放，对修改关闭
- **L**iskov Substitution：子类型可替换父类型
- **I**nterface Segregation：接口专一，避免"胖接口"
- **D**ependency Inversion：依赖抽象而非具体实现

**KISS (Keep It Simple, Stupid)**:
- 追求代码和设计的极致简洁
- 拒绝不必要的复杂性
- 优先选择最直观的解决方案

**DRY (Don't Repeat Yourself)**:
- 自动识别重复代码模式
- 主动建议抽象和复用
- 统一相似功能的实现方式

**YAGNI (You Aren't Gonna Need It)**:
- 仅实现当前明确所需的功能
- 抵制过度设计和未来特性预留
- 删除未使用的代码和依赖

#### 2. 后端架构（Go + DDD）

**分层架构**:
```
backend/
├── api/              # API层：HTTP处理器和路由
│   └── v1/          # API版本
├── application/      # 应用层：用例编排、事务边界
│   └── dto/         # 数据传输对象
├── domain/          # 领域层：核心业务逻辑（核心）
│   ├── entity/      # 实体
│   ├── repository/  # 仓储接口
│   └── service/     # 领域服务
└── infra/           # 基础设施层：技术实现
    ├── repository/  # 仓储实现
    └── client/      # 外部客户端
```

**关键架构模式**:
- **适配器模式**：层间解耦
- **接口隔离**：领域间清晰的契约定义
- **依赖倒置**：domain/ 不依赖任何外层
- **事件驱动**：使用消息队列异步通信

**模块所有权**（4人并行，互不干扰）:
```
backend/
├── domain/tenant/          # 研发A 独占
├── domain/permission/      # 研发A 独占
├── domain/routing/         # 研发A 独占
├── types/errno/           # 研发B 独占
├── tests/performance/     # 研发B 独占
├── infra/monitoring/      # 研发B 独占
├── infra/logging/         # 研发B 独占
└── infra/tracing/         # 研发B 独占
```

#### 3. 前端架构（React + TypeScript + Monorepo）

**Rush.js Monorepo 分层**:
```
frontend/packages/
├── arch/              # Level-1：核心基础设施
│   ├── ui-components/        # 研发C 独占 - 基础UI组件
│   └── business-components/  # 研发C 独占 - 业务组件
├── common/            # Level-2：共享组件和工具
│   ├── i18n/        # 研发C 独占 - 国际化
│   └── themes/      # 研发C 独占 - 主题系统
├── agent-ide/        # Level-3：Agent IDE功能域
├── workflow/         # Level-3：工作流功能域
├── studio/           # Level-3：Studio功能域
│   └── pages/
│       ├── tenant/      # 研发C 独占 - 租户管理
│       └── permission/  # 研发C 独占 - 权限管理
└── apps/             # Level-4：应用层
    └── coze-studio/  # 主应用
```

**关键架构模式**:
- **适配器模式**：广泛使用 `-adapter` 后缀包进行解耦
- **Base/Core 模式**：共享功能使用 `-base` 后缀包
- **工作区引用**：内部依赖使用 `workspace:*`

#### 4. 数据库架构（MySQL 8.4.5）

**核心设计原则**:
- 所有表都有 `tenant_id`（除系统配置表）
- 使用 `deleted_at` 软删除，避免硬删除
- 外键约束完整
- 索引设计遵循最左前缀原则
- 命名规范：`{table}_id`、`is_{property}`、`{action}_at`

**关键表分类**:
- **租户系统**：tenants、subscriptions、quotas、quota_usage
- **权限系统**：roles、data_permissions、field_permissions、user_roles
- **路由系统**：routing_rules、intent_matchers
- **业务表**：bots、conversations、knowledge、workflows（所有表都需要迁移）

---

## 🔐 核心开发规范

### 🎯 必读规范文档

**所有开发活动必须严格遵循以下规范**：

1. **[企业级开发规范手册](./docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)** ⭐⭐⭐
   - 后端开发规范（Go）：命名、函数设计、DDD分层、错误处理、并发安全
   - 前端开发规范（React + TypeScript）：组件设计、类型定义、性能优化
   - 数据库开发规范：表设计、索引优化、查询优化
   - API开发规范：RESTful设计、响应格式、错误码
   - Git工作流规范：分支策略、Conventional Commits
   - 代码审查清单：后端和前端

2. **[全局一致性检查清单](./docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md)** ⭐⭐⭐
   - 4级检查体系：L1个人 → L2模块 → L3集成 → L4发布
   - 8大维度：代码、架构、数据、API、前端、运维、文档、安全
   - 可执行的检查项，适用于每日工作

3. **[统一错误码定义规范](./docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)** ⭐⭐
   - 300+ 错误码完整定义
   - 错误响应格式规范
   - 中英文双语支持
   - 错误码使用示例

### 🔑 关键规范要点

#### 后端开发规范（Go）

**1. 命名规范**:
```go
// ✅ Good
package tenant       // 小写单数，描述性强
type TenantService struct {}
func (s *TenantService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error)
const MAX_BOT_COUNT = 100

// ❌ Bad
package tenants      // 复数
package util          // 过于通用
func create_bot() {}  // 下划线命名
```

**2. 函数设计规范**:
```go
// ✅ Good：函数 < 50行，参数 < 5个，错误包装
func (s *BotService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error) {
    if err := s.validateRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    // 业务逻辑...
    return bot, nil
}

// ❌ Bad：函数过长，参数过多
func (s *BotService) CreateBadExample(arg1, arg2, arg3, arg4, arg5, arg6 string) error {
    // 100+ 行代码...
}
```

**3. 并发安全规范**:
```go
// ✅ Good：限制并发数
func (s *BotService) ProcessBatch(ctx context.Context, botIDs []string) error {
    sem := make(chan struct{}, 10) // 最多10个并发
    var wg sync.WaitGroup
    for _, id := range botIDs {
        wg.Add(1)
        sem <- struct{}{}
        go func(botID string) {
            defer wg.Done()
            defer func() { <-sem }()
            s.processBot(ctx, botID)
        }(id)
    }
    wg.Wait()
    return nil
}
```

**4. 数据库操作规范**:
```go
// ✅ Good：使用事务，避免N+1查询
func (s *BotService) CreateBotWithConfig(ctx context.Context, bot *Bot, config *BotConfig) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(bot).Error; err != nil {
            return err
        }
        config.BotID = bot.BotID
        return tx.Create(config).Error
    })
}

// ✅ Good：使用Preload避免N+1
func (s *BotService) GetBotWithConfig(ctx context.Context, botID string) (*Bot, error) {
    var bot Bot
    err := s.db.Preload("Config").Preload("Creator").Where("bot_id = ?", botID).First(&bot).Error
    return &bot, err
}
```

#### 前端开发规范（React + TypeScript）

**1. 组件设计规范**:
```typescript
// ✅ Good：函数组件 + TypeScript + 清晰的Props
interface MessageBoxProps {
  message: {
    id: string;
    role: 'user' | 'assistant';
    content: string;
  };
}

export const MessageBox: React.FC<MessageBoxProps> = ({ message }) => {
  return <div className="message-box">{message.content}</div>;
};

// ❌ Bad：缺少类型，类组件
export class MessageBox extends React.Component {
  render() { return <div />; }
}
```

**2. 性能优化规范**:
```typescript
// ✅ Good：使用React.memo + useMemo + useCallback
export const BotCard: React.FC<BotCardProps> = React.memo(({ bot }) => {
  const sortedData = useMemo(() => {
    return data.sort((a, b) => a.createdAt - b.createdAt);
  }, [data]);

  return <div>{sortedData.map(item => <Item key={item.id} />)}</div>;
}, (prevProps, nextProps) => {
  return prevProps.bot.id === nextProps.bot.id;
});

// ❌ Bad：每次渲染都重新计算
export const BotCard: React.FC<BotCardProps> = ({ bot }) => {
  const sortedData = data.sort((a, b) => a.createdAt - b.createdAt);
  return <div>{sortedData.map(item => <Item />)}</div>;
};
```

**3. 状态管理规范**:
```typescript
// ✅ Good：Zustand store + 异步action
interface BotStore {
  bots: Bot[];
  loading: boolean;
  fetchBots: () => Promise<void>;
}

export const useBotStore = create<BotStore>((set, get) => ({
  bots: [],
  loading: false,
  fetchBots: async () => {
    set({ loading: true });
    try {
      const response = await botApi.list();
      set({ bots: response.data.bots, loading: false });
    } catch (error) {
      set({ loading: false });
      throw error;
    }
  },
}));
```

#### 数据库开发规范

**1. 表设计规范**:
```sql
-- ✅ Good：规范的命名和结构
CREATE TABLE bots (
    -- 主键：{entity}_id
    bot_id VARCHAR(36) PRIMARY KEY,

    -- 外键：{referenced_entity}_id
    tenant_id VARCHAR(36) NOT NULL,
    creator_id VARCHAR(36) NOT NULL,

    -- 布尔字段：is_{property}
    is_public BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    -- 时间戳：{action}_at
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    -- 索引
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status_created (status, created_at),
    UNIQUE KEY uk_tenant_name (tenant_id, name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**2. 查询优化规范**:
```sql
-- ✅ Good：避免SELECT *，使用游标分页
SELECT bot_id, name, status
FROM bots
WHERE tenant_id = 'xxx' AND created_at < :last_created_at
ORDER BY created_at DESC
LIMIT 21;

-- ✅ Good：使用JOIN避免N+1
SELECT b.bot_id, b.name, u.username
FROM bots b
LEFT JOIN users u ON b.creator_id = u.user_id
WHERE b.tenant_id = 'xxx';
```

---

## 🛠️ 开发工作流

### 环境配置

#### 方式一：Docker 快速启动（推荐）
```bash
# 1. 克隆项目
git clone https://github.com/coze-dev/coze-studio.git
cd coze-studio

# 2. 启动所有服务（MySQL、Redis、ES等）
cd docker
cp .env.example .env
docker compose up -d

# 3. 访问前端
open http://localhost:8888
```

#### 方式二：本地开发
```bash
# 1. 安装前端依赖
rush update

# 2. 启动中间件服务
make middleware

# 3. 启动Go后端
make server

# 4. 启动前端开发服务器
cd frontend/apps/coze-studio
npm run dev
```

### 常用命令

#### 构建命令
```bash
# 前端
rush build                    # 构建所有包
rush rebuild -o @coze-studio/app  # 构建特定包

# 后端
make build_server             # 构建Go服务器

# 完整构建
make web                      # Docker构建所有内容
```

#### 测试命令
```bash
# 前端测试（Vitest）
rush test                     # 运行所有测试
npm run test:cov             # 生成覆盖率报告

# 后端测试
cd backend && go test ./... -cover
```

#### 数据库命令
```bash
# 同步数据库架构
make sync_db

# 导出数据库架构
make dump_db

# 初始化SQL数据
make sql_init

# Atlas迁移管理
make atlas-hash
```

### Git 工作流

**分支策略**:
```
main         ← 生产环境，受保护
  ↑
develop      ← 开发环境，日常开发分支
  ↑
feature/*    ← 功能分支，从develop分出
bugfix/*     ← Bug修复分支，从develop分出
hotfix/*     ← 紧急修复分支，从main分出
```

**提交规范（Conventional Commits）**:
```bash
feat(scope): subject      # 新功能
fix(scope): subject       # Bug修复
docs(scope): subject      # 文档更新
test(scope): subject      # 测试相关
refactor(scope): subject  # 重构
perf(scope): subject      # 性能优化
style(scope): subject     # 代码格式（不影响功能）
chore(scope): subject     # 杂项（构建、依赖等）
```

**示例**:
```bash
git commit -m "feat(tenant): add tenant_id field to all tables"
git commit -m "fix(permission): resolve role check bug in data permission"
git commit -m "docs(api): update tenant management API documentation"
```

---

## 📦 关键技术栈

### 后端技术栈

| 类别 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **语言** | Go | 1.24.0 | 后端开发语言 |
| **HTTP框架** | Hertz | v0.10.2 | HTTP服务 |
| **ORM** | GORM | v1.25.11 | 数据库操作 |
| **数据库** | MySQL | 8.4.5 | 主数据库 |
| **缓存** | Redis | 8.0 | 缓存和会话 |
| **搜索** | Elasticsearch | 8.18.0 | 全文搜索 |
| **向量数据库** | Milvus | v2.5.10 | 向量存储 |
| **对象存储** | MinIO | latest | 文件存储 |
| **消息队列** | NSQ | v1.2.1 | 异步消息 |
| **配置中心** | etcd | 3.5 | 配置管理 |

### 前端技术栈

| 类别 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **语言** | TypeScript | 5.8.2 | 开发语言 |
| **框架** | React | 18.3.1 | UI框架 |
| **状态管理** | Zustand | latest | 全局状态 |
| **UI库** | Semi Design | latest | 组件库 |
| **样式** | Tailwind CSS | 3.3.3 | CSS框架 |
| **Monorepo** | Rush.js | latest | 包管理 |
| **构建工具** | Rsbuild | latest | 打包构建 |

### DevOps技术栈

| 类别 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **容器** | Docker | latest | 容器化 |
| **编排** | Kubernetes | latest | 容器编排 |
| **CI/CD** | GitHub Actions | latest | 持续集成 |
| **监控** | Prometheus | v2.45.0 | 指标采集 |
| **可视化** | Grafana | v10.0.0 | 监控大盘 |
| **日志** | ELK Stack | 8.x | 日志收集 |
| **追踪** | Jaeger | v1.50 | 分布式追踪 |
| **服务网格** | Istio | latest | 灰度发布 |

---

## 🎯 开发检查清单

### 代码提交前检查

**后端（Go）**:
- [ ] 代码符合 `ZKER-企业级开发规范手册_v1.0.md`
- [ ] `go test ./...` 通过，覆盖率 ≥ 80%
- [ ] `golangci-lint run` 通过，无警告
- [ ] 所有错误使用统一错误码
- [ ] 敏感信息不暴露到日志
- [ ] 并发安全（使用 go test -race）

**前端（React）**:
- [ ] 代码符合开发规范手册
- [ ] `rush test` 通过，覆盖率 ≥ 70%
- [ ] `rush lint` 通过，无警告
- [ ] TypeScript 类型完整，无 any
- [ ] 组件有 PropTypes/Interface
- [ ] 性能优化（React.memo、useMemo、useCallback）

**数据库**:
- [ ] 表命名符合规范（小写复数）
- [ ] 字段命名符合规范（`{table}_id`、`is_{property}`、`{action}_at`）
- [ ] 所有必要索引已创建
- [ ] 外键约束完整
- [ ] 使用 `deleted_at` 软删除

### Pull Request 前检查

- [ ] PR 描述清晰，包含变更说明和测试结果
- [ ] 所有 CI 检查通过
- [ ] 代码审查已通过
- [ ] 相关文档已更新
- [ ] 没有新的安全漏洞
- [ ] 没有明显的性能退化

---

## ⚠️ 常见问题与解决方案

### 前端开发

**问题**: 使用 `npm install` 失败
**解决**: 在根目录使用 `rush update`，Rush会处理包依赖

**问题**: 热重载不生效
**解决**: 检查 Rsbuild 配置，确保 `dev.hmr` 已启用

**问题**: 包之间循环依赖
**解决**: 使用 `rush check` 检查循环依赖，重构包结构

### 后端开发

**问题**: 数据库连接失败
**解决**:
```bash
# 1. 检查中间件服务是否运行
docker ps | grep mysql

# 2. 重启MySQL
docker compose restart mysql

# 3. 检查数据库连接配置
cat backend/conf/.env
```

**问题**: 模型调用失败
**解决**: 检查 `backend/conf/model/` 中是否配置了API密钥

**问题**: 测试失败
**解决**: 确保中间件服务运行，数据库已初始化

### Docker 开发

**问题**: 容器启动失败
**解决**:
```bash
# 检查日志
docker compose logs mysql

# 检查资源
docker stats

# 重置卷
make clean
docker compose up -d
```

**问题**: 端口冲突
**解决**: 修改 `docker-compose.yml` 中的端口映射

---

## 📖 详细文档链接

### 🚀 必读核心文档（企业级功能完善）
- [实现差距分析与研发计划](./docs/企业级功能完善与统一性设计方案/ZKER-实现差距分析与研发计划_v1.0.md) - 识别12大差距，4人×8周详细计划
- [企业级开发规范手册](./docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md) - 200+页完整开发规范
- [全局一致性检查清单](./docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md) - 4级检查体系

### 👥 个人开发计划
- [研发A - 后端架构师](./docs/企业级功能完善与统一性设计方案/研发A-后端架构师开发计划_v1.0.md)
- [研发B - 后端工程师](./docs/企业级功能完善与统一性设计方案/研发B-后端工程师开发计划_v1.0.md)
- [研发C - 前端工程师](./docs/企业级功能完善与统一性设计方案/研发C-前端工程师开发计划_v1.0.md)
- [研发D - DevOps工程师](./docs/企业级功能完善与统一性设计方案/研发D-DevOps工程师开发计划_v1.0.md)

### 📋 关键规范文档
- [统一错误码定义规范](./docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [数据迁移方案](./docs/企业级功能完善与统一性设计方案/ZKER-数据迁移方案_v1.0.md)
- [灰度发布策略](./docs/企业级功能完善与统一性设计方案/ZKER-灰度发布策略_v1.0.md)
- [一键回滚方案](./docs/企业级功能完善与统一性设计方案/ZKER-一键回滚方案_v1.0.md)
- [故障排查手册](./docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)

### 🎓 技术参考
- [技术组件清单与使用指南](./docs/企业级功能完善与统一性设计方案/ZKER-技术组件清单与使用指南(完整版).md)
- [核心算法实现指南](./docs/企业级功能完善与统一性设计方案/ZKER-核心算法实现指南.md)
- [开发快速入门指南](./docs/企业级功能完善与统一性设计方案/ZKER-开发快速入门指南.md)

### 🏗️ 架构设计
- [多租户SaaS架构](./docs/企业级功能完善与统一性设计方案/zker_MultiTenant_SaaS_完整架构设计文档.md)
- [数据库设计完整交付清单](./docs/企业级功能完善与统一性设计方案/数据库设计完整交付清单.md)
- [API设计规范文档](./docs/企业级功能完善与统一性设计方案/API设计规范文档.md)

### 📚 完整文档索引
- [文档索引](./docs/企业级功能完善与统一性设计方案/00-文档索引.md) - 80+份文档完整索引

### 🔍 其他开发规范
- [Coze Studio 开发规范](./docs/开发规范/开发规范总览.md) - 基础开发规范
- [开发规范/目录](./docs/开发规范/) - 详细开发规范目录

---

## 💡 获取帮助

### 使用 Skills

Claude Code 提供了多个技能（Skills）来帮助你快速完成常见任务：

```bash
# 查看可用技能
/skills

# 使用技能
/zcf:workflow      # 结构化六阶段开发工作流
/zcf:init-project  # 初始化项目AI上下文
/zcf:git-commit    # Git提交
/zcf:bmad-init     # 初始化项目
```

### 常见命令

| 命令 | 说明 |
|------|------|
| `/help` | 获取帮助信息 |
| `/commit` | 智能生成 Git 提交信息 |
| `/review` | 代码审查 |

---

**🎉 记住**: 所有开发活动必须严格遵循企业级开发规范，确保全局一致性！

**需要帮助？** 请查阅相关文档或使用 `/help` 命令。
