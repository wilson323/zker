# 跨模块集成测试文档

## 📋 概述

本节包含ZKER项目的跨模块集成测试，确保各核心模块间的接口正确性、数据流转正常和事务一致性。

## 🎯 测试覆盖范围

### 1. 租户+权限系统集成测试 (`tenant_permission_integration_test.go`)

**职责**: 验证租户管理与RBAC权限系统的集成

**测试场景**:
- ✅ 租户创建时自动初始化RBAC权限
- ✅ 租户隔离+数据权限过滤
- ✅ 角色分配与权限继承
- ✅ 字段级权限控制
- ✅ 租户删除时权限清理

**关键验证点**:
- 创建租户后自动创建TenantAdmin角色
- 不同租户的数据严格隔离
- 用户分配角色后权限正确继承
- 字段权限正确控制字段可见性

### 2. Saga+业务系统集成测试 (`saga_business_integration_test.go`)

**职责**: 验证Saga编排引擎与业务系统的集成

**测试场景**:
- ✅ Bot创建Saga完整流程
- ✅ Bot创建Saga失败时的补偿
- ✅ Saga重试机制
- ✅ Saga超时机制
- ✅ Saga并行执行
- ✅ Saga状态查询
- ✅ Saga定义验证

**关键验证点**:
- Saga执行成功后所有业务数据都已创建
- Saga步骤失败时自动触发补偿
- 步骤失败时正确重试
- 超时步骤正确处理

### 3. 人机协同+Bot系统集成测试 (`humaninloop_bot_integration_test.go`)

**职责**: 验证人机协同流程与Bot创建的集成

**测试场景**:
- ✅ Bot创建需要人工审核（低置信度）
- ✅ Bot创建无需审核（高置信度）
- ✅ 审核通过后Bot自动发布
- ✅ 审核拒绝后Bot保持草稿状态
- ✅ 批量审核Bot
- ✅ 被拒绝Bot重新提交
- ✅ 自动审核阈值配置
- ✅ 审核升级（复杂Bot）

**关键验证点**:
- AI置信度低于阈值时触发审核
- 高置信度Bot直接发布
- 审核通过后Bot自动发布
- 审核拒绝后Bot保持草稿状态

### 4. 记忆+对话系统集成测试 (`memory_conversation_integration_test.go`)

**职责**: 验证记忆系统与对话的集成

**测试场景**:
- ✅ 对话时注入记忆
- ✅ 对话后自动创建记忆
- ✅ 按相关性检索记忆
- ✅ 记忆过期和清理
- ✅ 记忆重要性衰减
- ✅ 跨会话记忆共享
- ✅ 记忆隐私和隔离
- ✅ 记忆更新和合并

**关键验证点**:
- 对话时自动注入相关记忆
- 对话后自动提取并存储记忆
- 记忆按相关性排序
- 过期记忆自动清理
- 不同用户记忆严格隔离

## 🛠️ 技术栈

### 测试框架
- **Go testing**: Go原生测试框架
- **testify/suite**: 测试套件支持
- **testcontainers-go**: 容器化测试环境

### 测试工具
- **testcontainers**: MySQL 8.4.5测试容器
- **Mock对象**: 自定义Mock服务
  - MockBotRepository
  - MockKnowledgeRepository
  - MockLLMClient
  - MockMilvusClient

### 覆盖率工具
- **go test -cover**: Go原生覆盖率
- **go tool cover**: 覆盖率报告生成

## 🚀 快速开始

### 运行所有跨模块集成测试

**Linux/macOS**:
```bash
cd backend/tests/integration
chmod +x run-integration-tests.sh
./run-integration-tests.sh
```

**Windows**:
```cmd
cd backend\tests\integration
run-integration-tests.bat
```

### 运行特定测试套件

```bash
# 租户+权限系统
./run-integration-tests.sh -s tenant-permission

# Saga+业务系统
./run-integration-tests.sh -s saga-business

# 人机协同+Bot系统
./run-integration-tests.sh -s humaninloop-bot

# 记忆+对话系统
./run-integration-tests.sh -s memory-conversation
```

### 生成覆盖率报告

```bash
# 运行测试并生成覆盖率
go test -coverprofile=coverage.out -covermode=atomic ./tests/integration/...

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 查看总体覆盖率
go tool cover -func=coverage.out | grep total
```

## 📁 目录结构

```
tests/integration/
├── fixtures/                              # 测试fixtures和helpers
│   ├── test_data.go                       # 测试数据生成函数
│   ├── mock_repositories.go               # Mock仓储实现
│   └── mock_services.go                   # Mock服务实现
│
├── tenant_permission_integration_test.go  # 租户+权限集成测试
├── saga_business_integration_test.go      # Saga+业务集成测试
├── humaninloop_bot_integration_test.go    # 人机协同+Bot集成测试
├── memory_conversation_integration_test.go # 记忆+对话集成测试
│
├── run-integration-tests.sh               # 测试执行脚本 (Linux/macOS)
└── run-integration-tests.bat              # 测试执行脚本 (Windows)
```

## 🔧 核心设计原则

### 1. 单一职责原则 (SRP)

**每个测试文件只负责一个集成场景**
- ✅ `tenant_permission_integration_test.go` - 仅测试租户+权限集成
- ❌ 禁止一个测试文件测试多个不相关的集成场景

**每个测试函数只负责一个验证点**
- ✅ `TestTenantCreationWithRBAC` - 仅测试租户创建+RBAC集成
- ❌ 禁止一个测试函数测试多个不相关的流程

### 2. 测试隔离原则

**数据隔离**:
- 每个测试使用独立的测试数据
- 测试前清理所有相关表数据
- 测试结束后自动清理容器

**时间隔离**:
- 测试可以按任意顺序执行
- 测试之间无依赖关系
- 支持并发执行

### 3. 真实环境原则

**必须使用testcontainers**:
- ✅ MySQL 8.4.5 testcontainers
- ❌ 禁止使用本地MySQL
- ❌ 禁止使用Mock数据库

**外部依赖使用Mock**:
- ✅ LLM调用使用MockLLMClient
- ✅ 向量数据库使用MockMilvusClient

## 📊 验证清单

### 测试编写检查清单

- [ ] 每个测试文件只负责一个集成场景
- [ ] 每个测试函数只负责一个验证点
- [ ] 所有测试使用testcontainers + MySQL 8.4.5
- [ ] 测试可以独立运行
- [ ] 测试数据准备完整
- [ ] 测试数据清理完整
- [ ] 断言清晰明确
- [ ] 错误处理完整

### 测试执行检查清单

- [ ] 集成测试覆盖率 ≥ 80%
- [ ] 所有测试可以独立运行
- [ ] 测试执行时间 < 5分钟
- [ ] 测试数据清理完整
- [ ] 无竞态条件
- [ ] 无内存泄漏

## 🐛 常见问题

### 1. testcontainers启动失败

**问题**: Docker容器启动超时

**解决方案**:
```bash
# 检查Docker是否运行
docker info

# 手动拉取镜像
docker pull mysql:8.4.5
```

### 2. 测试数据未清理

**问题**: 测试间数据污染

**解决方案**:
```go
// 确保每个测试前清理
func (s *TestSuite) SetupTest() {
    tables := []string{...}
    CleanupTestData(s.db, tables)
}
```

## 📈 性能指标

### 目标性能指标

| 指标 | 目标值 |
|------|--------|
| 测试覆盖率 | ≥ 80% |
| 测试执行时间 | < 5分钟 |
| 测试并发数 | ≥ 4 |

## 📚 参考文档

- [Go Testing指南](https://golang.org/pkg/testing/)
- [testcontainers-go文档](https://golang.testcontainers.org/)
- [企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

## 联系方式

如有问题，请联系：
- 集成测试负责人: 研发A（后端架构师）
