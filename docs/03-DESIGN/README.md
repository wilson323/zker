# ZKER 设计文档索引

**版本**: v3.0.0 | **更新**: 2025-01-03

---

## 📚 设计文档导航

本目录包含所有企业级功能的设计文档，涵盖多租户、RBAC、智能路由等核心模块。

### 🎯 核心设计文档

| 模块 | 文档 | 优先级 | 完成度 |
|------|------|--------|--------|
| **多租户架构** | [multi-tenant/](multi-tenant/) | P0 | ✅ 90% |
| **RBAC 权限** | [rbac/](rbac/) | P0 | ✅ 90% |
| **智能路由** | [routing/](routing/) | P0 | ✅ 90% |
| **Agent 监控** | [agent-monitoring/](agent-monitoring/) | P1 | 🚧 60% |
| **开发者平台** | [developer-platform/](developer-platform/) | P1 | 🚧 50% |

---

## 🏗️ 多租户架构设计

**目录**: [multi-tenant/](multi-tenant/)

### 核心文档

| 文档 | 说明 | 状态 |
|------|------|------|
| [architecture.md](multi-tenant/architecture.md) | 多租户架构设计 | ✅ |
| [data-model.md](multi-tenant/data-model.md) | 数据模型设计 | ✅ |
| [isolation.md](multi-tenant/isolation.md) | 租户隔离机制 | ✅ |
| [migration.md](multi-tenant/migration.md) | 数据迁移方案 | ✅ |
| [billing.md](multi-tenant/billing.md) | 订阅计费系统 | ✅ |
| [quota.md](multi-tenant/quota.md) | 配额管理系统 | ✅ |

### 核心特性

- ✅ **数据隔离**: 完整的租户数据隔离
- ✅ **配额管理**: 弹性配额，实时监控
- ✅ **订阅管理**: 多层级订阅计划
- ✅ **计费系统**: Token 计量，预算告警
- ✅ **安全合规**: 审计日志，合规要求

### 技术实现

```go
// 租户实体
type Tenant struct {
    TenantID     string    `gorm:"primaryKey"`
    TenantName   string    `gorm:"not null"`
    PlanType     string    // FREE, BASIC, PRO, ENTERPRISE
    Status       string    // ACTIVE, SUSPENDED, DELETED
    // ...
}

// 配额检查
func (s *QuotaService) CheckCreateBot(ctx context.Context, tenantID string) error {
    quota, err := s.GetByTenantID(ctx, tenantID)
    if quota.BotCount >= quota.MaxBots {
        return errorx.New(errno.QuotaExceeded)
    }
    return nil
}
```

---

## 🔐 RBAC 权限系统设计

**目录**: [rbac/](rbac/)

### 核心文档

| 文档 | 说明 | 状态 |
|------|------|------|
| [architecture.md](rbac/architecture.md) | RBAC 架构设计 | ✅ |
| [data-model.md](rbac/data-model.md) | 数据模型设计 | ✅ |
| [data-permission.md](rbac/data-permission.md) | 数据权限设计 | ✅ |
| [field-permission.md](rbac/field-permission.md) | 字段权限设计 | ✅ |
| [implementation.md](rbac/implementation.md) | 实施指南 | ✅ |

### 权限级别

#### 数据权限（5 级）

| 级别 | 说明 | 权限范围 |
|------|------|---------|
| **ALL** | 全部数据 | 可访问所有租户数据（仅管理员） |
| **DEPARTMENT** | 部门数据 | 可访问本部门及下级部门数据 |
| **OWN** | 个人数据 | 仅可访问自己创建的数据 |
| **CUSTOM** | 自定义 | 根据自定义规则过滤 |
| **NONE** | 无权限 | 不可访问任何数据 |

#### 字段权限（3 级）

| 级别 | 说明 |
|------|------|
| **hidden** | 隐藏，不可见 |
| **readonly** | 只读，可见但不可修改 |
| **editable** | 可编辑，完全可访问 |

### 技术实现

```go
// 数据权限检查
func (s *PermissionService) CheckDataPermission(
    ctx context.Context,
    userID string,
    resourceType string,
    resourceID string,
) error {
    permission, err := s.GetDataPermission(ctx, userID, resourceType)
    if err != nil {
        return err
    }

    switch permission.Level {
    case permission.DataPermissionAll:
        return nil // 管理员，允许访问
    case permission.DataPermissionOwn:
        return s.checkOwnership(ctx, userID, resourceID)
    case permission.DataPermissionNone:
        return errorx.New(errno.PermissionDenied)
    }
    return nil
}
```

---

## 🧠 智能路由引擎设计

**目录**: [routing/](routing/)

### 核心文档

| 文档 | 说明 | 状态 |
|------|------|------|
| [architecture.md](routing/architecture.md) | 路由引擎架构 | ✅ |
| [matchers.md](routing/matchers.md) | 匹配器设计 | ✅ |
| [rules.md](routing/rules.md) | 路由规则配置 | ✅ |
| [implementation.md](routing/implementation.md) | 实施指南 | ✅ |

### 匹配器类型

| 匹配器 | 说明 | 优先级 |
|--------|------|--------|
| **规则匹配器** | 关键词、正则、意图、分类 | 1 |
| **相似度匹配器** | 向量相似度计算 | 2 |
| **混合匹配器** | 多匹配器聚合 | 3 |
| **路由决策器** | 基于评分的智能路由 | 4 |

### 技术实现

```go
// 路由引擎
type RoutingEngine struct {
    ruleMatcher      *RuleMatcher
    similarityMatcher *SimilarityMatcher
    hybridMatcher    *HybridMatcher
    decisionMaker    *DecisionMaker
}

func (e *RoutingEngine) Route(
    ctx context.Context,
    query string,
    candidates []Agent,
) (*Agent, error) {
    // 1. 规则匹配
    ruleScores := e.ruleMatcher.Match(ctx, query, candidates)

    // 2. 相似度匹配
    simScores := e.similarityMatcher.Match(ctx, query, candidates)

    // 3. 混合评分
    finalScores := e.hybridMatcher.Combine(ruleScores, simScores)

    // 4. 决策路由
    return e.decisionMaker.Decide(ctx, candidates, finalScores)
}
```

---

## 📊 Agent 监控平台设计

**目录**: [agent-monitoring/](agent-monitoring/)

### 核心文档

| 文档 | 说明 | 状态 |
|------|------|------|
| [architecture.md](agent-monitoring/architecture.md) | 监控平台架构 | 🚧 |
| [metrics.md](agent-monitoring/metrics.md) | 监控指标设计 | 🚧 |
| [tracing.md](agent-monitoring/tracing.md) | 链路追踪设计 | 🚧 |
| [evaluation.md](agent-monitoring/evaluation.md) | 质量评估设计 | 🚧 |

### 监控指标

| 类别 | 指标 |
|------|------|
| **性能** | QPS、响应时间、错误率 |
| **质量** | 用户满意度、回答相关性、准确性 |
| **成本** | Token 使用量、API 调用次数 |
| **业务** | 活跃用户数、对话数、转化率 |

### 对标产品

- **LangSmith**: 全链路追踪、质量评估
- **Arize**: ML 模型监控
- **Weights & Biases**: 实验追踪

---

## 👨‍💻 开发者平台设计

**目录**: [developer-platform/](developer-platform/)

### 核心文档

| 文档 | 说明 | 状态 |
|------|------|------|
| [architecture.md](developer-platform/architecture.md) | 开发者平台架构 | 🚧 |
| [cli.md](developer-platform/cli.md) | CLI 工具设计 | 🚧 |
| [sdk.md](developer-platform/sdk.md) | SDK 设计 | 🚧 |
| [api-docs.md](developer-platform/api-docs.md) | API 文档设计 | 🚧 |

### 平台功能

- **CLI 工具**: 命令行工具，快速开发
- **多语言 SDK**: Python、JavaScript、Go
- **API 文档**: OpenAPI 规范，交互式文档
- **开发者社区**: 论坛、教程、示例

### 对标产品

- **Dify**: 开发者平台、API 文档
- **Langflow**: 可视化开发工具
- **FastGPT**: 快速部署能力

---

## 🔄 设计文档使用指南

### 阅读顺序建议

**架构师/技术负责人**:
1. 多租户架构设计
2. RBAC 架构设计
3. 智能路由架构设计
4. Agent 监控平台设计

**后端开发人员**:
1. 对应模块的数据模型
2. 对应模块的实施指南
3. 相关 API 文档

**前端开发人员**:
1. 对应模块的架构概览
2. 对应模块的 API 接口
3. 相关组件设计

### 文档规范

- ✅ **架构图**: 使用 Mermaid 图表
- ✅ **数据模型**: 使用 ER 图
- ✅ **API 接口**: 使用 OpenAPI 规范
- ✅ **代码示例**: 使用真实可运行代码

---

## 📊 设计文档统计

| 模块 | 文档数量 | 完成度 | 优先级 |
|------|---------|--------|--------|
| **多租户** | 6 | ✅ 90% | P0 |
| **RBAC** | 5 | ✅ 90% | P0 |
| **智能路由** | 4 | ✅ 90% | P0 |
| **Agent 监控** | 4 | 🚧 60% | P1 |
| **开发者平台** | 4 | 🚧 50% | P1 |
| **总计** | **23** | **82%** | - |

---

## 🎯 下一步计划

### P0 任务（本周）

- [ ] 完善多租户架构文档（10%）
- [ ] 完善 RBAC 架构文档（10%）
- [ ] 完善智能路由文档（10%）

### P1 任务（下周）

- [ ] 完成 Agent 监控平台设计（40%）
- [ ] 完成开发者平台设计（50%）

---

## 📞 反馈渠道

- **设计问题**: [GitHub Issues](https://github.com/coze-dev/coze-studio/issues)
- **设计建议**: [GitHub Discussions](https://github.com/coze-dev/coze-studio/discussions)
- **设计咨询**: zker-design@example.com

---

**🎯 目标**: 打造行业领先的企业级 AI 智能体工作台平台！
