# E2E测试文档

## 概述

本目录包含Coze Studio的完整端到端(E2E)测试套件，模拟真实用户场景，从API请求到数据库操作的完整流程。

## 设计原则

### 1. 单一职责原则(SRP)

**每个E2E测试文件只负责一个用户旅程**:
- ✅ `tenant_registration_journey_test.go` - 租户注册旅程
- ✅ `bot_creation_approval_journey_test.go` - Bot创建+审核旅程
- ✅ `organization_permission_journey_test.go` - 组织管理+权限控制旅程
- ✅ `conversation_memory_journey_test.go` - 对话+记忆旅程
- ❌ 禁止一个测试文件混合多个用户旅程

**每个E2E测试只负责一个完整的业务流程**:
- ✅ `TestE2E_TenantRegistrationJourney` - 完整的租户注册流程
- ✅ `TestE2E_BotCreationWithApprovalJourney` - 完整的Bot创建+审核流程
- ❌ 禁止一个E2E测试拆分成多个不相关的流程

### 2. 真实环境模拟

所有E2E测试必须:
- 使用testcontainers启动真实的MySQL 8.4.5容器
- 启动完整的HTTP服务器
- 执行真实的数据库操作
- 模拟完整的用户交互流程

### 3. 完整流程覆盖

每个E2E测试必须覆盖:
1. HTTP请求 → 2. 业务逻辑 → 3. 数据库操作 → 4. 数据验证

## 测试套件结构

```
backend/tests/e2e/
├── test_application.go              # 测试应用封装（完整的服务栈）
├── helpers/
│   └── test_helpers.go              # 测试辅助函数（单一职责）
├── tenant_registration_journey_test.go       # 租户注册旅程
├── bot_creation_approval_journey_test.go     # Bot创建+审核旅程
├── organization_permission_journey_test.go   # 组织+权限旅程
├── conversation_memory_journey_test.go       # 对话+记忆旅程
├── run_e2e_tests.sh                 # E2E测试执行脚本
└── README.md                        # 本文档
```

## 核心组件

### 1. TestApplication

封装完整的应用栈，提供统一的测试环境:

```go
app := e2e.SetupTestApplication(t)
defer app.Terminate(ctx)

// 发送HTTP请求
app.PostJSON("/api/v1/tenants", req, &resp)
app.GetJSON("/api/v1/tenants/" + id, &resp)

// 访问数据库
app.DB.Table("tenants").Where("tenant_id = ?", id).First(&tenant)
```

### 2. 测试辅助函数

每个函数只负责一个创建或验证任务:

```go
// 创建测试数据
tenant := e2eHelpers.CreateTestTenant(t, app, "测试公司")
user := e2eHelpers.CreateTestUser(t, app, tenantID, "admin")
bot := e2eHelpers.CreateTestBot(t, app, tenantID, "Bot名称", 0.95)

// 验证数据
e2eHelpers.AssertTenantExists(t, app.DB, tenantID)
e2eHelpers.AssertBotStatus(t, app.DB, botID, "published")
e2eHelpers.AssertReviewTaskExists(t, app.DB, tenantID, "pending")
```

## 测试场景

### 场景1: 租户注册旅程

**流程**: 提交注册信息 → 创建租户 → 初始化RBAC → 分配默认配额 → 验证欢迎邮件

**验证点**:
- ✅ 租户已创建
- ✅ 默认角色已创建
- ✅ 默认配额已分配
- ✅ 订阅信息正确
- ✅ 并发注册稳定性

### 场景2: Bot创建+审核旅程

**流程**: 创建Bot(AI置信度低) → 创建审核任务 → 审核人分配 → 人工审核 → 审核通过 → Bot发布

**验证点**:
- ✅ 低置信度Bot进入草稿状态
- ✅ 审核任务自动创建
- ✅ 审核任务可分配
- ✅ 审核通过后Bot发布
- ✅ 审核历史完整记录
- ✅ 高置信度Bot直接发布
- ✅ 审核拒绝流程

### 场景3: 组织管理+权限控制旅程

**流程**: 创建部门 → 添加成员 → 分配部门权限 → 验证权限缓存 → 验证数据权限过滤

**验证点**:
- ✅ 组织结构正确创建
- ✅ 用户可添加到组织
- ✅ 数据权限正确分配
- ✅ 权限缓存已更新
- ✅ 数据权限过滤生效
- ✅ 多级组织层级支持

### 场景4: 对话+记忆旅程

**流程**: 用户对话 → 提取实体记忆 → 后续对话检索记忆 → 验证记忆增强效果

**验证点**:
- ✅ 实体自动提取
- ✅ 记忆存储正确
- ✅ 记忆检索生效
- ✅ 回答风格符合偏好
- ✅ 记忆重要性评分
- ✅ 记忆访问计数
- ✅ 多实体提取
- ✅ 对话上下文连续性

## 执行测试

### 方式1: 使用脚本执行（推荐）

```bash
# 执行所有E2E测试
./backend/tests/e2e/run_e2e_tests.sh

# 执行特定测试
./backend/tests/e2e/run_e2e_tests.sh -test.v -test.run TestE2E_TenantRegistrationJourney

# 生成覆盖率报告
./backend/tests/e2e/run_e2e_tests.sh -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### 方式2: 使用Go测试命令

```bash
# 进入E2E测试目录
cd backend/tests/e2e

# 执行所有E2E测试
go test -v -timeout 30m

# 执行特定测试文件
go test -v -run TestE2E_TenantRegistrationJourney -timeout 10m

# 跳过短测试
go test -v -timeout 30m

# 生成覆盖率
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 方式3: 并发执行（加快测试速度）

```bash
# 并发执行E2E测试（注意：可能需要数据库隔离）
go test -v -parallel 4 -timeout 30m
```

## 测试覆盖率目标

| 类型 | 目标覆盖率 | 当前状态 |
|------|-----------|---------|
| E2E整体覆盖率 | ≥ 70% | 待验证 |
| 用户旅程覆盖 | 100% | 4/4 |
| 关键业务流程 | 100% | 待验证 |

## 测试数据管理

### 数据隔离

每个测试使用独立的数据:
- 每个测试自动创建独立的租户
- 每个测试使用独立的用户
- 测试结束后自动清理数据

### 数据清理

测试完成后自动清理:
```go
t.Cleanup(func() {
    e2eHelpers.CleanupTenantData(t, ctx, app.DB, tenant.TenantID)
})
```

## 关键约束

### 1. 单一职责强制

**每个测试文件/函数只负责一个旅程/流程**:
- ❌ 禁止一个测试文件混合多个用户旅程
- ❌ 禁止一个测试函数包含多个不相关的流程

### 2. 真实环境模拟

**必须启动完整的HTTP服务器和数据库**:
- ✅ 使用testcontainers + MySQL 8.4.5
- ✅ 使用httptest.Server
- ❌ 禁止使用Mock Handler（除非必要）

### 3. 完整流程覆盖

**从HTTP请求到数据库操作的所有步骤**:
- ✅ API请求验证
- ✅ 业务逻辑验证
- ✅ 数据库操作验证
- ✅ 数据一致性验证

### 4. 异步处理

**正确处理异步任务和轮询**:
```go
// 使用轮询等待异步操作完成
e2eHelpers.WaitForCondition(t, func() bool {
    var status string
    app.DB.Table("bots").Where("bot_id = ?", botID).Select("status").Scan(&status)
    return status == "published"
}, 30*time.Second, time.Second)
```

## 故障排查

### MySQL容器启动失败

**问题**: testcontainers无法启动MySQL容器
**解决**:
```bash
# 检查Docker是否运行
docker ps

# 检查端口占用
netstat -an | grep 3306

# 重启Docker
# macOS
restart-docker-desktop

# Linux
sudo systemctl restart docker
```

### 测试超时

**问题**: 测试执行超时
**解决**:
```bash
# 增加超时时间
go test -v -timeout 60m

# 或者并行执行
go test -v -parallel 4 -timeout 30m
```

### 数据库连接失败

**问题**: 无法连接到testcontainers MySQL
**解决**:
```go
// 增加等待时间
wait.ForLog("ready for connections").WithOccurrence(2).WithStartupTimeout(10*time.Minute)
```

## 最佳实践

### 1. 测试命名

遵循清晰的命名规范:
```go
// Good: 描述性命名
TestE2E_TenantRegistrationJourney
TestE2E_BotCreationWithApprovalJourney

// Bad: 不明确的命名
TestE2E_1
TestE2E_Case
```

### 2. 测试结构

使用子测试组织复杂流程:
```go
func TestE2E_TenantRegistrationJourney(t *testing.T) {
    app := e2e.SetupTestApplication(t)

    t.Run("Step1_提交注册信息", func(t *testing.T) { ... })
    t.Run("Step2_验证租户已创建", func(t *testing.T) { ... })
    t.Run("Step3_验证默认角色已创建", func(t *testing.T) { ... })
}
```

### 3. 断言使用

优先使用require:
```go
// require: 失败立即停止
require.NoError(t, err, "创建租户应该成功")
require.NotEmpty(t, tenantID, "租户ID不应为空")

// assert: 失败继续执行（用于多个验证）
assert.Equal(t, "active", status, "状态应为active")
assert.Greater(t, count, 0, "数量应大于0")
```

## 贡献指南

添加新的E2E测试时，请遵循:

1. **单一职责**: 一个测试文件只负责一个用户旅程
2. **完整流程**: 覆盖从HTTP请求到数据库操作的所有步骤
3. **真实环境**: 使用testcontainers和真实HTTP服务器
4. **清晰命名**: 使用描述性的测试和步骤名称
5. **数据隔离**: 确保测试数据独立，不互相影响

## 参考

- [Go测试最佳实践](https://go.dev/doc/tutorial/add-a-test)
- [testcontainers-go文档](https://golang.testcontainers.org/)
- [testify文档](https://github.com/stretchr/testify)
