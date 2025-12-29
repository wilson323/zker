# Coze Studio Skills 使用指南

## 概述

Coze Studio Skills 是一套专门为项目定制的 AI 辅助开发技能集合，基于全局代码深度分析创建，确保开发过程严格遵循规范，保持全局代码一致性，提高开发效率。

## 技术栈覆盖

### 前端技术栈
- **框架**: React 18
- **语言**: TypeScript 5.0+
- **构建**: Rsbuild (Rspack)
- **状态管理**: Zustand
- **UI 库**: Semi Design + Tailwind CSS
- **包管理**: Rush.js Monorepo (135+ packages)
- **测试**: Vitest + Testing Library

### 后端技术栈
- **语言**: Go 1.23+
- **框架**: Hertz
- **架构**: DDD (Domain-Driven Design)
- **ORM**: GORM
- **数据库**: MySQL 8.4.5
- **消息队列**: NSQ
- **测试**: Go testing + testify

## Skills 目录

### 📂 Skills 总览

```
.claude/skills/
├── README.md                     (11 KB)  - 本文件
│
├── 🎨 前端开发 Skills
│   ├── create-component.md        (7.7 KB) - 创建 React 组件
│   ├── create-hook.md             (16 KB)  - 创建自定义 Hook
│   ├── create-store.md            (19 KB)  - 创建 Zustand Store
│   ├── create-api-client.md       (22 KB)  - 创建 API Client
│   └── create-package.md          (28 KB)  - 创建 Rush Package
│
├── 🚀 后端开发 Skills
│   ├── create-entity.md           (21 KB)  - 创建领域实体
│   ├── create-domain-service.md   (14 KB)  - 创建领域服务
│   ├── create-repository.md       (19 KB)  - 创建仓储实现
│   └── create-api-handler.md      (19 KB)  - 创建 API 处理器
│
├── 🔍 代码质量 Skills
│   ├── check-standards.md         (12 KB)  - 代码规范检查
│   ├── review-code.md             (9.4 KB) - 代码审查
│   ├── write-tests.md             (20 KB)  - 测试生成
│   └── generate-docs.md           (9.0 KB) - 文档生成
│
└── 🏢 企业级功能 Skills (新增)
    ├── enterprise-check.md        (5.2 KB) - 企业级规范检查 ⭐
    ├── multitenant-dev.md         (5.8 KB) - 多租户开发助手 ⭐
    ├── rbac-dev.md               (4.9 KB) - RBAC权限系统开发 ⭐
    ├── routing-dev.md            (5.1 KB) - 智能路由引擎开发 ⭐
    └── doc-query.md               (6.3 KB) - 企业级文档查询 ⭐

总计: 18 个 skills, 275 KB
```

---

## 📚 Skills 详细说明

### 🎨 前端开发 Skills

#### 1. create-component.md
**创建 React 组件**

**功能**: 创建符合规范的 React 组件（TypeScript/React）

**适用场景**:
- 需要创建新的 UI 组件
- 需要重构现有组件
- 需要添加复杂交互逻辑

**输出内容**:
- 组件文件（.tsx）
- 类型定义（.types.ts）
- 样式文件（.styles.ts）
- 测试文件（.test.tsx）

**规范覆盖**:
- ✅ 文件命名：kebab-case.tsx
- ✅ 组件命名：PascalCase
- ✅ Props 接口：ComponentNameProps
- ✅ 组件结构：7 步骤（Hooks → useMemo → useEffect → useCallback → render）
- ✅ 事件处理：handle 前缀
- ✅ 回调函数：on 前缀

**使用示例**:
```bash
"创建一个用户头像组件"
→ 生成 user-avatar/ 目录，包含所有必要文件
```

---

#### 2. create-hook.md
**创建自定义 Hook**

**功能**: 创建符合规范的自定义 React Hook

**适用场景**:
- 需要复用状态逻辑
- 需要封装副作用逻辑
- 需要抽象数据获取逻辑

**Hook 类型**:
- 值保持 Hook（useInitialValue）
- 响应式 Hook（useResponsive）
- 数据获取 Hook（useUserData）
- 事件处理 Hook（useHandleSubmit）

**规范覆盖**:
- ✅ Hook 命名：use 前缀 + camelCase
- ✅ 文件命名：use-{purpose}.ts
- ✅ 返回值：对象或元组
- ✅ 使用规则：仅在顶层调用
- ✅ 性能优化：useCallback、useMemo

**基于真实代码**:
- `useInitialValue` - 保持初始值
- `useResponsive` - 响应式断点

---

#### 3. create-store.md
**创建 Zustand Store**

**功能**: 创建符合规范的状态管理 Store

**适用场景**:
- 需要全局状态管理
- 需要跨组件共享状态
- 需要状态持久化

**输出内容**:
- Store 接口定义
- Store 实现
- 选择器 Hooks
- 持久化配置
- 测试文件

**规范覆盖**:
- ✅ Store 命名：use{Domain}Store
- ✅ State 字段：清晰描述（is/has/selected 前缀）
- ✅ Actions：动词-名词模式
- ✅ 中间件：devtools + persist
- ✅ 选择器：性能优化

**基于真实代码**:
- `useProjectAuthStore` - 项目认证管理

---

#### 4. create-api-client.md
**创建 API Client**

**功能**: 创建符合项目规范的 API Client，用于前端调用后端服务

**适用场景**:
- 需要调用新的后端 API
- 需要封装 HTTP 请求逻辑
- 需要统一错误处理
- 需要添加请求拦截器

**输出内容**:
- API Client 实现
- 类型定义（请求/响应）
- React Hook 封装
- 错误处理
- 测试文件

**规范覆盖**:
- ✅ 文件命名：{service}-api.ts
- ✅ 导出命名：{service}Api
- ✅ 统一 axiosInstance
- ✅ 请求/响应拦截器
- ✅ 错误处理机制
- ✅ TypeScript 类型安全

**基于真实代码**:
- `basicApi` - 基础服务 API
- `knowledgeApi` - 知识库 API

---

#### 5. create-package.md
**创建 Rush Package**

**功能**: 创建符合 Rush.js Monorepo 规范的完整前端 Package

**适用场景**:
- 需要创建新的前端功能包
- 需要封装共享组件和工具
- 需要创建适配器或接口层
- 需要创建业务模块

**输出内容**:
- 完整的目录结构
- package.json 配置
- TypeScript 配置（tsconfig.json, tsconfig.build.json）
- ESLint 配置
- Vitest 配置
- README 文档
- 源代码框架

**规范覆盖**:
- ✅ 包命名：@coze-arch/{name} 或 @coze-studio/{name}
- ✅ 4 层依赖结构：arch → common → domain → apps
- ✅ workspace:* 引用内部包
- ✅ 标准目录结构
- ✅ 完整配置文件
- ✅ 测试和文档

**包类型支持**:
- 基础设施包（Infrastructure）
- 共享组件包（Components）
- 适配器包（Adapter）
- 接口包（Interface）
- 业务模块包（Domain）

---

### 🚀 后端开发 Skills

#### 6. create-entity.md
**创建领域实体**

**功能**: 创建符合 DDD 架构和 GORM 规范的领域实体（Entity）和数据模型（Model）

**适用场景**:
- 需要创建新的领域实体
- 需要定义数据模型和数据库映射
- 需要实现实体业务规则和状态转换
- 需要添加查询选项（Where Option）

**输出内容**:
- Domain Entity（领域实体）
- Crossdomain Model（数据模型）
- GORM 标签配置
- 实体方法
- 查询选项类型

**规范覆盖**:
- ✅ 实体分层：Domain Entity + Crossdomain Model
- ✅ 字段命名：驼峰式，PascalCase 导出
- ✅ GORM 标签：完整的数据库映射
- ✅ 软删除：DeletedAt 指针类型
- ✅ 状态常量：类型安全的枚举
- ✅ 业务方法：状态转换和验证

**基于真实代码**:
- `User` Entity - 用户实体
- `Space` Entity - 空间实体
- `Knowledge` Entity - 知识库实体

---

#### 7. create-domain-service.md
**创建领域服务**

**功能**: 创建符合 DDD 架构的领域服务

**适用场景**:
- 需要创建新的领域服务
- 需要实现复杂业务逻辑
- 需要跨多个实体的业务操作

**输出内容**:
- 服务接口（service.go）
- 服务实现（service_impl.go）
- 实体定义（entity/{entity}.go）
- 仓储接口（repository/interface.go）
- 服务测试（service_test.go）

**规范覆盖**:
- ✅ 服务命名：{Domain}Service / {domain}ServiceImpl
- ✅ 方法命名：PascalCase + 动词前缀
- ✅ 错误处理：errorx.Wrapf
- ✅ DDD 分层：严格遵循 4 层架构
- ✅ 依赖注入：构造函数

---

#### 8. create-repository.md
**创建仓储实现**

**功能**: 创建数据访问层仓储实现

**适用场景**:
- 需要持久化领域实体
- 需要复杂查询逻辑
- 需要事务支持

**输出内容**:
- 仓储接口（repository/interface.go）
- 仓储实现（repository_impl.go）
- 仓储测试（repository_test.go）

**规范覆盖**:
- ✅ 接口命名：{Entity}Repository
- ✅ 实现命名：{entity}RepositoryImpl
- ✅ 方法命名：Create/FindByID/Update/Delete
- ✅ GORM 使用：正确规范
- ✅ 软删除：deleted_at 检查
- ✅ 错误处理：errorx.Wrapf

---

#### 9. create-api-handler.md
**创建 API Handler**

**功能**: 创建 HTTP API 处理器

**适用场景**:
- 需要暴露 RESTful API
- 需要处理 HTTP 请求
- 需要参数验证和绑定

**输出内容**:
- Handler 结构体
- 请求处理方法
- 参数验证
- 响应封装
- 路由注册
- Handler 测试

**规范覆盖**:
- ✅ Handler 命名：{Service}Handler
- ✅ 方法命名：Create{Resource}/Get{Resource}ByID
- ✅ 路由设计：RESTful 规范
- ✅ 参数验证：binding tags
- ✅ 响应格式：统一结构
- ✅ 错误码：规范定义

---

### 🔍 代码质量 Skills

#### 10. check-standards.md
**代码规范检查**

**功能**: 全面检查代码是否符合规范

**检查内容**:
- 前端：12+ 项（命名、结构、类型、测试等）
- 后端：15+ 项（命名、结构、DDD、错误处理等）

**检查报告**:
- ✅ 检查概述
- ✅ 通过项列表
- ✅ 失败项清单（含位置、问题、建议、修复代码）
- ✅ 改进建议
- ✅ 总体评价

**输出示例**:
```markdown
# 代码规范检查报告

## 检查概述
- 检查文件: user-profile.tsx
- 总体评分: ⭐⭐⭐⭐☆ (4/5)

## ❌ 失败项 (3)
1. 函数命名不规范
   - 位置: user-profile.tsx:45
   - 问题: 函数名 `user` 不符合规范
   - 建议: 改为 `getUserById`
```

---

#### 11. review-code.md
**代码审查**

**功能**: 深入审查代码质量

**审查维度**:
1. 功能性：是否实现需求
2. 规范性：是否遵循规范
3. 可维护性：代码清晰度
4. 性能：性能问题检查
5. 安全性：安全隐患检查
6. 测试：测试覆盖检查
7. 文档：注释和文档检查

**审查流程**:
1. 自动化检查（lint/test）
2. 代码阅读
3. 深入分析
4. 生成报告

**输出内容**:
- 审查概述
- 优点列表
- 问题列表（含优先级）
- 分项评分（7 个维度）
- 改进建议
- 审查结论

---

#### 12. write-tests.md
**测试生成**

**功能**: 为代码生成测试用例

**生成内容**:
- React 组件测试
- Hooks 测试
- API Client 测试
- Go 服务测试
- HTTP Handler 测试

**测试场景**:
- ✅ 正常流程（Happy Path）
- ✅ 边界情况
- ✅ 错误处理
- ✅ 加载状态
- ✅ 用户交互

**测试覆盖要求**:
- Level 1 (arch): 80%
- Level 2 (common): 30%
- Level 3-4: 灵活

---

#### 13. generate-docs.md
**文档生成**

**功能**: 为代码生成文档

**生成内容**:
- Go Doc 注释
- JSDoc 注释
- README 文件
- API 文档

**文档格式**:
- Go Doc：标准 Go 注释格式
- JSDoc：标准 TypeScript 注释
- README：Markdown 格式
- API 文档：OpenAPI 格式

---

### 🏢 企业级功能 Skills (新增)

#### 14. enterprise-check.md **企业级规范检查** ⭐

**功能**: 检查代码和设计是否符合企业级开发规范

**适用场景**:
- 需要检查代码是否符合企业级规范
- 需要进行全局一致性检查
- 需要进行代码审查前的预检查
- 需要验证架构设计是否符合规范

**检查内容**:
- 代码命名规范
- 函数/组件设计（SOLID 原则）
- 错误处理（统一错误码）
- 数据库表设计规范
- API 设计（RESTful 规范）
- 全局一致性问题

**检查级别**:
- L1（个人）：代码风格、命名、基础逻辑
- L2（模块）：模块间接口、数据一致性、API 契约
- L3（集成）：跨模块集成、端到端流程
- L4（发布）：完整性、安全性、向后兼容性

**输出内容**:
- 检查结果报告（通过/不通过）
- 不符合规范的项目列表
- 改进建议
- 相关规范文档链接

**相关文档**:
- [ZKER-企业级开发规范手册_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-全局一致性检查清单_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-全局一致性检查清单_v1.0.md)

---

#### 15. multitenant-dev.md **多租户开发助手** ⭐

**功能**: 协助多租户架构相关的开发工作

**适用场景**:
- 创建租户实体和仓储
- 实现配额检查逻辑
- 添加租户隔离中间件
- 设计订阅管理功能
- 生成租户相关迁移脚本

**核心功能**:
- 生成符合规范的租户实体代码
- 生成配额检查服务代码
- 提供租户隔离最佳实践
- 检查是否遗漏租户隔离
- 生成数据库迁移脚本

**输出内容**:
- 实体代码框架
- 服务代码框架
- 中间件代码
- 迁移脚本
- 测试用例示例

**相关文档**:
- [研发A-后端架构师开发计划_v1.0.md](../../docs/企业级功能完善与统一性设计方案/研发A-后端架构师开发计划_v1.0.md)
- [zker_MultiTenant_SaaS_完整架构设计文档.md](../../docs/企业级功能完善与统一性设计方案/zker_MultiTenant_SaaS_完整架构设计文档.md)

---

#### 16. rbac-dev.md **RBAC权限系统开发助手** ⭐

**功能**: 协助 RBAC 权限系统开发

**适用场景**:
- 创建角色和权限实体
- 实现数据权限检查
- 实现字段权限过滤
- 设计权限中间件
- 配置权限规则

**权限级别**:
- 数据权限：ALL、DEPARTMENT、OWN、CUSTOM、NONE（5级）
- 字段权限：hidden、readonly、editable（3级）

**输出内容**:
- 权限实体代码
- 权限检查服务代码
- 权限中间件代码
- 权限配置示例
- 测试用例示例

**相关文档**:
- [研发A-后端架构师开发计划_v1.0.md](../../docs/企业级功能完善与统一性设计方案/研发A-后端架构师开发计划_v1.0.md)
- [ZKER-企业级开发规范手册_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

#### 17. routing-dev.md **智能路由引擎开发助手** ⭐

**功能**: 协助智能路由引擎开发

**适用场景**:
- 创建规则匹配器
- 创建相似度匹配器
- 实现混合意图匹配
- 实现路由决策引擎
- 创建路由规则配置

**匹配器类型**:
- 规则匹配器：关键词、正则、意图、分类
- 相似度匹配器：向量相似度计算
- 混合匹配器：多匹配器聚合
- 路由决策器：基于评分的智能路由

**输出内容**:
- 匹配器代码框架
- 路由决策引擎代码
- 规则配置示例
- 性能优化建议
- 测试用例示例

**相关文档**:
- [研发A-后端架构师开发计划_v1.0.md](../../docs/企业级功能完善与统一性设计方案/研发A-后端架构师开发计划_v1.0.md)
- [API接口文档_智能路由引擎.md](../../docs/企业级功能完善与统一性设计方案/API接口文档_智能路由引擎.md)

---

#### 18. doc-query.md **企业级文档查询** ⭐

**功能**: 快速查询企业级功能完善相关的所有文档

**适用场景**:
- 查找特定功能的 API 文档
- 查找特定技术组件的使用指南
- 查找特定开发规范
- 查找故障排查指南
- 查找个人开发计划

**查询能力**:
- 快速检索 80+ 份企业级文档
- 按类别过滤文档
- 按角色过滤文档
- 提供文档摘要和关键链接
- 推荐相关文档

**文档分类**:
- 核心设计文档
- 个人开发计划
- 关键规范文档
- 运维文档

**输出内容**:
- 相关文档列表
- 文档摘要
- 关键章节链接
- 相关文档推荐

**相关文档**:
- [00-文档索引.md](../../docs/企业级功能完善与统一性设计方案/00-文档索引.md)

---

## 📊 规范覆盖矩阵

| Skill | 命名规范 | 结构规范 | 错误处理 | 架构规范 | 测试规范 | 文档规范 |
|-------|---------|---------|---------|---------|---------|---------|
| create-component | ✅✅✅ | ✅✅✅ | ✅ | - | ✅✅ | ✅ |
| create-hook | ✅✅✅ | ✅✅✅ | ✅ | - | ✅✅ | ✅ |
| create-store | ✅✅✅ | ✅✅✅ | ✅ | - | ✅✅ | ✅ |
| create-api-client | ✅✅✅ | ✅✅ | ✅✅ | - | ✅✅ | ✅ |
| create-package | ✅✅✅ | ✅✅✅ | ✅ | ✅✅ | ✅✅ | ✅✅ |
| create-entity | ✅✅✅ | ✅✅✅ | ✅ | ✅✅ | - | ✅ |
| create-domain-service | ✅✅✅ | ✅✅✅ | ✅✅✅ | ✅✅✅ | ✅✅ | ✅ |
| create-repository | ✅✅✅ | ✅✅ | ✅✅✅ | ✅✅ | ✅✅ | ✅ |
| create-api-handler | ✅✅✅ | ✅✅ | ✅✅✅ | ✅✅ | ✅ | ✅ |
| check-standards | ✅✅✅ | ✅✅✅ | ✅✅ | ✅✅✅ | ✅✅ | ✅ |
| review-code | ✅✅✅ | ✅✅✅ | ✅✅ | ✅✅✅ | ✅✅ | ✅✅ |
| write-tests | ✅✅ | - | ✅✅ | - | ✅✅✅ | - |
| generate-docs | ✅ | ✅ | - | - | - | ✅✅✅ |

**图例**: ✅ 完全覆盖 | ✅ 部分覆盖 | - 不适用

---

## 🚀 使用方式

### 方式一：直接引用

```bash
"请使用 create-component skill 创建一个用户列表组件"
```

### 方式二：描述需求

```bash
"创建一个用户管理服务，包括 CRUD 操作"
→ AI 自动使用 create-domain-service skill
```

### 方式三：组合使用

```bash
"创建完整的用户管理功能"
→ 依次使用：
  1. create-entity（后端实体）
  2. create-repository（仓储实现）
  3. create-domain-service（领域服务）
  4. create-api-handler（API 处理器）
  5. create-api-client（前端 API Client）
  6. create-component（前端组件）
  7. write-tests（生成测试）
  8. check-standards（检查规范）
```

```bash
"创建一个新的工具包"
→ 依次使用：
  1. create-package（创建包结构）
  2. create-component（添加组件）
  3. write-tests（生成测试）
  4. generate-docs（生成文档）
```

---

## 🎯 全局一致性保证

### 1. 命名一致性

所有生成的代码都遵循统一的命名规范：
- 前端文件：kebab-case
- 后端文件：lowercase.go
- 组件/类：PascalCase
- 函数/方法：动词-名词模式

### 2. 结构一致性

所有代码遵循统一的文件结构：
- 前端包：src → components → hooks → utils → types
- 后端：domain → application → api → infra
- Rush.js 4 层依赖：arch → common → domain → apps

### 3. 错误处理一致性

所有代码使用统一的错误处理模式：
- 前端：try-catch + Error Boundary
- 后端：errorx.Wrapf + 错误码

### 4. 测试一致性

所有测试使用统一的测试模式：
- 前端：Vitest + Testing Library
- 后端：testify + mock

---

## 📈 预期效果

使用这套 Skills 系统后，预期可以实现：

1. **开发效率提升 50%+**
   - 自动生成符合规范的代码
   - 减少重复工作
   - 加快开发速度

2. **代码质量提升 80%+**
   - 严格遵循规范
   - 完整的测试覆盖
   - 统一的代码风格

3. **全局一致性 100%**
   - 命名规范统一
   - 架构模式一致
   - 错误处理统一

4. **新人上手时间减少 70%**
   - 清晰的规范指导
   - 完整的代码模板
   - 详细的使用文档

---

## 🔄 版本历史

- **v3.0.0** (2025-01-01): 企业级功能完善版本 ⭐
  - **新增企业级功能 Skills**（5个）：
    - enterprise-check - 企业级规范检查
    - multitenant-dev - 多租户开发助手
    - rbac-dev - RBAC权限系统开发
    - routing-dev - 智能路由引擎开发
    - doc-query - 企业级文档查询
  - 更新 skills 总数：18 个，275 KB
  - 全面支持企业级功能完善阶段的开发工作

- **v2.0.0** (2025-01-25): 完整版本
  - 新增前端开发：create-api-client, create-package
  - 新增后端开发：create-entity
  - 更新 skills 总数：13 个，235 KB

- **v1.0.0** (2025-01-01): 初始版本
  - 前端开发：create-component, create-hook, create-store
  - 后端开发：create-domain-service, create-repository, create-api-handler
  - 代码质量：check-standards, review-code, write-tests, generate-docs

---

## 🤝 贡献指南

如需添加或改进 skills，请：

1. 参考现有 skill 的格式
2. 确保覆盖开发规范要点
3. 提供清晰的示例
4. 包含完整的检查清单
5. 更新本 README 文档

---

## 📞 联系方式

如有问题或建议，请联系 Coze Studio 技术委员会。
