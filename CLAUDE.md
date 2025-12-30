# ZKER - 企业级 AI 智能体工作台平台

> 基于鲸智百应平台的企业级增强版 | 目标：超越竞品，打造 SaaS 行业标杆

**版本**: v3.0.0 | **更新**: 2025-01-03 | **阶段**: 企业级功能完善

---

## 🎯 平台定位

ZKER 是一个**企业级多租户 SaaS AI Agent 开发平台**，对标 Dify、Langflow、FastGPT 等行业领先产品，具备以下核心能力：

| 能力 | 说明 | 对标产品 |
|------|------|---------|
| **多租户隔离** | 完整的数据隔离、配额管理、订阅计费 | AWS SaaS |
| **RBAC 权限** | 5级数据权限 + 3级字段权限 | 企业级 ERP |
| **智能路由** | 混合意图匹配 + 评分路由 | LangSmith |
| **Agent 监控** | 全链路追踪、质量评估、成本分析 | LangSmith |
| **开发者平台** | CLI、SDK、API 文档 | Dify |

**当前任务**: 从单租户应用升级为企业级多租户 SaaS 平台

---

## 📚 核心文档导航

### 🚀 必读文档（按角色）

| 角色 | 必读文档 | 用途 |
|------|---------|------|
| **架构师** | [03-DESIGN/README.md](docs/03-DESIGN/README.md) | 系统架构设计 |
| **后端开发** | [02-SPECS/backend-dev-guide.md](docs/02-SPECS/backend-dev-guide.md) | 后端开发规范 |
| **前端开发** | [02-SPECS/frontend-dev-guide.md](docs/02-SPECS/frontend-dev-guide.md) | 前端开发规范 |
| **DevOps** | [04-IMPLEMENTATION/deployment.md](docs/04-IMPLEMENTATION/deployment.md) | 部署运维指南 |
| **全员必读** | [02-SPECS/README.md](docs/02-SPECS/README.md) | 开发规范索引 |

### 📋 企业级功能文档

| 功能模块 | 文档 | 优先级 |
|---------|------|--------|
| **多租户架构** | [03-DESIGN/multi-tenant/](docs/03-DESIGN/multi-tenant/) | P0 |
| **RBAC 权限** | [03-DESIGN/rbac/](docs/03-DESIGN/rbac/) | P0 |
| **智能路由** | [03-DESIGN/routing/](docs/03-DESIGN/routing/) | P0 |
| **错误码系统** | [02-SPECS/error-handling.md](docs/02-SPECS/error-handling.md) | P0 |
| **监控告警** | [05-OPERATIONS/monitoring.md](docs/05-OPERATIONS/monitoring.md) | P1 |

### 📖 完整文档索引

详见: [docs/00-META/README.md](docs/00-META/README.md)

---

## 🛠️ 快速开始

### 环境启动（Docker）

```bash
# 1. 启动中间件服务
cd docker && docker compose up -d

# 2. 启动后端
cd .. && make server

# 3. 启动前端
rush build && cd frontend/apps/coze-studio && npm run dev

# 4. 访问
open http://localhost:8888
```

### 常用命令

```bash
# 构建
rush build              # 构建前端
make build_server       # 构建后端

# 测试
rush test               # 前端测试
go test ./... -cover    # 后端测试

# 检查
rush lint               # 前端检查
golangci-lint run       # 后端检查
```

---

## 🎨 可用 Skills

Claude Code 已集成企业级开发助手：

```bash
# 企业级功能
/enterprise-check      # 企业级规范检查 ⭐
/multitenant-dev       # 多租户开发助手 ⭐
/rbac-dev             # RBAC 权限系统开发 ⭐
/routing-dev          # 智能路由引擎开发 ⭐
/doc-query            # 文档查询 ⭐

# 通用开发
/create-component      # 创建 React 组件
/create-entity         # 创建领域实体
/check-standards       # 代码规范检查
```

详见: [.claude/skills/README.md](.claude/skills/README.md)

---

## 🏗️ 技术架构

### 后端架构（Go + DDD）

```
backend/
├── api/              # API层：HTTP 处理器
├── application/      # 应用层：用例编排
├── domain/          # 领域层：核心业务逻辑 ⭐
│   ├── tenant/      # 多租户（研发A 独占）
│   ├── permission/  # RBAC 权限（研发A 独占）
│   └── routing/     # 智能路由（研发A 独占）
└── infra/           # 基础设施层
```

**设计原则**: SOLID + KISS + DRY + YAGNI

### 前端架构（React + Monorepo）

```
frontend/packages/
├── arch/              # Level-1：基础设施
├── common/            # Level-2：共享组件
├── agent-ide/         # Level-3：Agent IDE
├── workflow/          # Level-3：工作流
├── studio/            # Level-3：Studio
│   └── pages/
│       ├── tenant/    # 租户管理（研发C 独占）
│       └── permission/# 权限管理（研发C 独占）
└── apps/coze-studio/  # Level-4：主应用
```

**设计模式**: 适配器模式 + Base/Core 模式

### 技术栈

| 类别 | 技术 |
|------|------|
| **后端** | Go 1.24 + Hertz + GORM |
| **前端** | React 18 + TypeScript + Semi Design |
| **数据库** | MySQL 8.4.5 + Redis 8.0 + Elasticsearch 8.18 |
| **运维** | Docker + Kubernetes + Prometheus + Grafana |

---

## 📋 提交前检查

### 代码质量

- [ ] 符合开发规范（后端/前端）
- [ ] 测试通过（覆盖率 ≥ 80%）
- [ ] Lint 检查通过
- [ ] 使用统一错误码

### 架构一致性

- [ ] DDD 分层清晰
- [ ] 多租户隔离完整
- [ ] RBAC 权限检查
- [ ] API 符合 RESTful 规范

### 文档完整性

- [ ] API 文档已更新
- [ ] 变更说明已添加
- [ ] 相关设计文档已同步

---

## 🔑 关键规范速查

### 命名规范

```go
// ✅ Good
package tenant
type TenantService struct {}
func (s *TenantService) CreateBot(ctx context.Context, req *CreateBotRequest) (*Bot, error)

// ❌ Bad
package tenants
func create_bot() {}
```

### 错误处理

```go
// ✅ Good：使用统一错误码
return errorx.Wrapf(err, errno.TenantNotFound)

// ❌ Bad：硬编码错误
return errors.New("tenant not found")
```

### 多租户隔离

```go
// ✅ Good：所有查询包含租户过滤
db.Where("tenant_id = ? AND deleted_at IS NULL", tenantID).Find(&bots)

// ❌ Bad：缺少租户隔离
db.Find(&bots)
```

---

## ⚡ 开发团队分工

| 角色 | 负责模块 | 独占目录 |
|------|---------|---------|
| **研发A（后端架构师）** | 多租户 + RBAC + 智能路由 | backend/domain/{tenant,permission,routing} |
| **研发B（后端工程师）** | 错误码 + 性能测试 + 监控 | backend/{types/errno,tests,infra/monitoring} |
| **研发C（前端工程师）** | 组件库 + 业务页面 | frontend/packages/{arch,business,studio/pages} |
| **研发D（DevOps）** | Docker + K8s + CI/CD | .github/workflows/, docker/, k8s/ |

---

## 💡 获取帮助

```bash
# 查看可用命令
/help

# 代码审查
/review

# Git 提交
/commit
```

**文档门户**: [docs/00-META/README.md](docs/00-META/README.md)
**问题反馈**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)

---

**🎯 目标**: 超越鲸智百应，打造企业级 SaaS AI 智能体工作台行业标杆！
