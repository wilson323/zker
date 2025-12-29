# 代码示例库

## 概述

本代码示例库提供 ZKER Enterprise 平台的核心功能实现示例，遵循项目的开发规范和架构模式。

## 目录结构

```
examples/
├── backend/                    # Go 后端示例
│   ├── domain/                # 领域层示例
│   │   ├── entity.go         # 实体定义
│   │   ├── repository.go     # 仓储接口
│   │   └── service.go        # 领域服务
│   ├── application/           # 应用层示例
│   │   ├── dto.go            # 数据传输对象
│   │   └── service.go        # 应用服务
│   ├── api/                   # API层示例
│   │   └── handler.go        # HTTP处理器
│   ├── infra/                 # 基础设施示例
│   │   ├── repository.go     # 仓储实现
│   │   ├── middleware/       # 中间件
│   │   └── rag/              # RAG实现
│   └── crossdomain/           # 跨域关注点
│       ├── auth/             # 认证
│       ├── permission/       # 权限
│       └── metering/         # Token计量
└── frontend/                  # React 前端示例
    ├── hooks/                # 自定义Hooks
    ├── services/             # API服务
    ├── components/           # 组件示例
    └── stores/               # 状态管理
```

## 快速导航

### 后端示例（Go）

| 示例 | 说明 | 文件路径 |
|-----|------|---------|
| 实体定义 | 租户实体领域模型 | `backend/domain/entity/tenant.go` |
| 仓储模式 | 数据访问抽象层 | `backend/domain/repository/bot_repository.go` |
| 领域服务 | 复杂业务逻辑封装 | `backend/domain/service/conversation_service.go` |
| 应用服务 | 用例编排 | `backend/application/service/bot_service.go` |
| API处理器 | HTTP请求处理 | `backend/api/handler/bot_handler.go` |
| 认证中间件 | JWT认证 | `backend/infra/middleware/auth.go` |
| 权限中间件 | RBAC权限检查 | `backend/infra/middleware/permission.go` |
| Token计量 | 使用量统计和告警 | `backend/crossdomain/metering/token_service.go` |
| RAG实现 | 混合检索和重排序 | `backend/infra/rag/hybrid_search.go` |

### 前端示例（React + TypeScript）

| 示例 | 说明 | 文件路径 |
|-----|------|---------|
| 自定义Hook | 认证状态管理 | `frontend/hooks/useAuth.ts` |
| API服务 | HTTP客户端封装 | `frontend/services/api.ts` |
| 表单组件 | Bot创建表单 | `frontend/components/BotForm.tsx` |
| 状态管理 | Zustand Store | `frontend/stores/botStore.ts` |

## 使用方法

### 复制代码到项目

```bash
# 假设当前在 examples/ 目录
cd ../..

# 复制后端示例
cp -r docs/企业级功能完善与统一性设计方案/examples/backend/* backend/

# 复制前端示例
cp -r docs/企业级功能完善与统一性设计方案/examples/frontend/* frontend/
```

### 代码示例特点

- **完整性**: 每个示例都是可运行的完整代码
- **注释详尽**: 关键逻辑都有详细的中英文注释
- **遵循规范**: 严格遵守项目开发规范
- **最佳实践**: 展示SOLID、DRY、KISS原则的应用
- **生产就绪**: 包含错误处理、日志、性能优化

## 架构模式说明

### DDD四层架构

```
┌─────────────────────────────────────┐
│         API Layer (接口层)           │  HTTP/gRPC处理器
├─────────────────────────────────────┤
│    Application Layer (应用层)        │  用例编排、DTO
├─────────────────────────────────────┤
│      Domain Layer (领域层)           │  实体、仓储接口、领域服务
├─────────────────────────────────────┤
│   Infrastructure Layer (基础设施层)  │  仓储实现、外部服务
└─────────────────────────────────────┘
```

### 依赖方向

- API → Application → Domain
- Infrastructure → Domain（实现接口）
- 所有层都可以依赖 CrossDomain（跨域关注点）

## 代码规范要点

### Go 代码规范

1. **包命名**: 小写单词，不使用下划线或驼峰
2. **接口命名**: 动词+名词，以 `er` 或 `or` 结尾
3. **错误处理**: 使用自定义错误类型，始终处理错误
4. **日志**: 使用结构化日志，包含请求追踪ID

### TypeScript 代码规范

1. **组件命名**: PascalCase，描述性名称
2. **文件命名**: kebab-case 或 PascalCase（组件）
3. **类型定义**: 优先使用 interface，复杂类型使用 type
4. **Hooks**: 以 `use` 开头，遵循React Hooks规则

## 扩展和定制

所有示例代码都可以作为模板进行扩展：

1. **重命名实体**: 替换 `Bot` 为你的业务实体
2. **调整字段**: 根据业务需求增减字段
3. **添加验证**: 在DTO层增加自定义验证规则
4. **扩展查询**: 在Repository层添加复杂查询方法
5. **定制权限**: 根据业务场景调整权限粒度

## 反馈和改进

如果你发现代码示例有问题或有改进建议，请通过以下方式反馈：

1. 在项目仓库提 Issue
2. 提交 Pull Request
3. 联系架构团队

## 版本历史

- v1.0.0 (2025-01-15): 初始版本，包含核心功能示例
