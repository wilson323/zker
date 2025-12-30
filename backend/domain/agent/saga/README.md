# Bot创建Saga分布式事务系统

## 概述

本包实现了企业级分布式事务系统，使用**Saga模式**保证Bot创建过程中的**最终一致性**。

### 核心特性

✅ **完整的Saga模式实现** - 长事务编排
✅ **自动补偿机制** - 失败时自动回滚已执行的步骤
✅ **事务状态追踪** - 完整的执行状态监控
✅ **幂等性保证** - 支持重试和故障恢复
✅ **超时控制** - 步骤级和Saga级超时控制
✅ **指数退避重试** - 智能的重试策略

## 业务流程

Bot创建Saga包含以下3个步骤：

```
1. 创建Bot记录 (create_bot)
   ↓
2. 创建默认知识库 (create_knowledge_base)
   ↓
3. 分配Bot权限 (assign_permissions)
```

### 补偿流程（任一步骤失败时触发）

```
失败发生
   ↓
逆序执行补偿：
3. 撤销Bot权限 (assign_permissions_compensation)
   ↓
2. 删除知识库 (create_knowledge_base_compensation)
   ↓
1. 删除Bot记录 (create_bot_compensation)
```

## 架构设计

### 分层架构

```
domain/agent/saga/          # 领域层：业务Saga定义
├── bot_creation_saga.go    # Bot创建Saga实现
└── bot_creation_saga_test.go

infra/saga/                 # 基础设施层：Saga框架
├── definition.go           # Saga定义（状态、接口、错误）
├── orchestrator.go         # Saga协调器
├── coordinator.go          # Saga事件协调器
├── repository.go           # Saga持久化仓储
└── utils.go               # 工具函数（ID生成等）
```

### 核心概念

#### 1. Saga（事务定义）

```go
type Saga struct {
    ID            string            // Saga唯一标识
    Name          string            // Saga名称
    Description   string            // Saga描述
    Steps         []SagaStep        // 执行步骤
    Compensations []CompensationStep // 补偿步骤
    Timeout       time.Duration     // 超时时间
    RetryPolicy   *RetryPolicy      // 重试策略
}
```

#### 2. SagaStep（执行步骤）

```go
type SagaStep interface {
    Execute(ctx context.Context, data interface{}) (interface{}, error)
    Name() string
    Timeout() time.Duration
}
```

#### 3. CompensationStep（补偿步骤）

```go
type CompensationStep interface {
    Compensate(ctx context.Context, data interface{}) error
    Name() string
    Timeout() time.Duration
}
```

#### 4. SagaExecution（执行记录）

```go
type SagaExecution struct {
    ID             string           // 执行记录ID
    SagaID         string           // Saga定义ID
    Status         SagaStatus       // 执行状态
    CurrentStep    int              // 当前步骤索引
    InputData      interface{}      // 输入数据
    OutputData     interface{}      // 输出数据
    Error          error            // 错误信息
    StartedAt      time.Time        // 开始时间
    CompletedAt    *time.Time       // 完成时间
    StepExecutions []StepExecution  // 步骤执行记录
}
```

### 状态机

```
pending → running → completed
           ↓
        failed → compensating → compensated
```

## 使用示例

### 1. 定义Saga

```go
import (
    "github.com/coze-dev/coze-studio/backend/domain/agent/saga"
    "github.com/coze-dev/coze-studio/backend/infra/saga"
)

// 创建Saga定义
sagaDef := saga.NewBotCreationSaga()

// 输出：
// ID: saga-bot-creation
// Name: bot-creation-saga
// Steps: 3
// Compensations: 3
// Timeout: 5m0s
```

### 2. 注册Saga

```go
import (
    "go.uber.org/zap"
    "gorm.io/gorm"
)

// 创建仓储
repo := saga.NewMySQLSagaRepository(db)

// 创建协调器
logger := zap.NewExample()
orchestrator := saga.NewSagaOrchestrator(repo, logger)

// 注册Saga定义
if err := orchestrator.DefineSaga(sagaDef); err != nil {
    log.Fatal("定义Saga失败", err)
}
```

### 3. 执行Saga

```go
// 准备输入数据
cmd := &saga.BotCreationCommand{
    TenantID:    "tenant-123",
    Name:        "客服助手",
    Description: "智能客服Bot",
    Avatar:      "avatar.png",
    Type:        "chatbot",
    CreatorID:   "user-456",
}

// 执行Saga
ctx := context.Background()
execution, err := orchestrator.ExecuteSaga(ctx, "bot-creation-saga", cmd)
if err != nil {
    // Saga执行失败，已自动补偿
    log.Error("Saga执行失败", zap.Error(err))
    return
}

// Saga执行成功
log.Info("Saga执行成功",
    zap.String("execution_id", execution.ID),
    zap.String("status", string(execution.Status)))
```

### 4. 查询状态

```go
// 查询执行状态
execution, err := orchestrator.GetStatus(ctx, executionID)
if err != nil {
    log.Error("查询状态失败", zap.Error(err))
    return
}

log.Info("Saga状态",
    zap.String("execution_id", execution.ID),
    zap.String("status", string(execution.Status)),
    zap.Int("current_step", execution.CurrentStep))
```

## 测试

### 运行测试

```bash
# 测试Bot创建Saga
cd backend
go test -v ./domain/agent/saga/...

# 测试Saga框架
go test -v ./infra/saga/...

# 运行所有测试
go test -v ./domain/agent/saga/... ./infra/saga/...
```

### 测试覆盖

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./domain/agent/saga/...
go tool cover -html=coverage.out
```

### 测试输出示例

```
=== RUN   TestNewBotCreationSaga
    ✓ Saga定义验证通过
      - ID: saga-bot-creation
      - Name: bot-creation-saga
      - Description: Bot创建流程：创建Bot记录 → 创建默认知识库 → 分配权限
      - Steps: 3
      - Compensations: 3
      - Timeout: 5m0s
--- PASS: TestNewBotCreationSaga (0.00s)

=== RUN   TestCreateBotStep
[执行] 创建Bot: 测试Bot (ID: 4b854e63-9929-4fb7-b47d-57fa8a2dec64)
    ✓ 创建Bot步骤执行成功
      - Bot ID: 4b854e63-9929-4fb7-b47d-57fa8a2dec64
      - Bot Name: 测试Bot
      - Bot Status: draft
--- PASS: TestCreateBotStep (0.00s)

=== RUN   TestCreateBotCompensation
[补偿] 删除Bot: 测试Bot (ID: 66b93cc2-1e49-4b95-a813-d882c41929c3)
    ✓ 创建Bot补偿执行成功
      - 已删除 Bot ID: 66b93cc2-1e49-4b95-a813-d882c41929c3
--- PASS: TestCreateBotCompensation (0.00s)

PASS
ok      github.com/coze-dev/coze-studio/backend/domain/agent/saga    4.976s
```

## 企业级特性

### 1. 完整的补偿机制

- ✅ 自动检测步骤失败
- ✅ 逆序执行补偿事务
- ✅ 保证最终一致性

### 2. 可靠性保证

- ✅ 幂等性：支持安全的重试
- ✅ 超时控制：防止步骤无限等待
- ✅ 状态追踪：完整的执行历史记录

### 3. 性能优化

- ✅ 指数退避重试：智能的重试策略
- ✅ 并发安全：使用互斥锁保护共享状态
- ✅ 事件驱动：异步事件通知

### 4. 可观测性

- ✅ 结构化日志：使用zap记录关键操作
- ✅ 状态监控：实时追踪Saga执行状态
- ✅ 错误追踪：详细的错误信息和堆栈

## 错误处理

### 错误类型

```go
// Saga错误
var (
    ErrSagaNotFound         = &SagaError{Code: "SAGA_NOT_FOUND"}
    ErrSagaTimeout          = &SagaError{Code: "SAGA_TIMEOUT"}
    ErrStepFailed           = &SagaError{Code: "STEP_FAILED"}
    ErrCompensationFailed   = &SagaError{Code: "COMPENSATION_FAILED"}
    ErrInvalidSagaDefinition = &SagaError{Code: "INVALID_SAGA_DEFINITION"}
)
```

### 错误示例

```go
execution, err := orchestrator.ExecuteSaga(ctx, "bot-creation-saga", cmd)
if err != nil {
    var sagaErr *saga.SagaError
    if errors.As(err, &sagaErr) {
        log.Error("Saga执行失败",
            zap.String("code", sagaErr.Code),
            zap.String("message", sagaErr.Message),
            zap.String("saga_id", sagaErr.SagaID),
            zap.String("step_name", sagaErr.StepName),
            zap.Bool("recoverable", sagaErr.Recoverable))
    }
    return
}
```

## 性能指标

### 基准测试

```bash
# ID生成性能测试
go test -bench=. -benchmem ./domain/agent/saga/...
```

### 预期性能

- ID生成: ~1,000 ns/op
- 步骤执行: < 100ms（实际业务逻辑）
- 补偿执行: < 100ms（实际业务逻辑）
- 整体超时: 5分钟（可配置）

## 最佳实践

### 1. 步骤设计

- ✅ 每个步骤应该是幂等的
- ✅ 步骤之间通过返回值传递数据
- ✅ 避免步骤之间有强依赖

### 2. 补偿设计

- ✅ 补偿逻辑应该简单可靠
- ✅ 补偿失败不应阻止后续补偿
- ✅ 补偿应该是幂等的

### 3. 超时设置

```go
// 步骤超时：30秒
stepTimeout := 30 * time.Second

// Saga超时：5分钟
sagaTimeout := 5 * time.Minute
```

### 4. 重试策略

```go
// 默认重试策略
&RetryPolicy{
    MaxAttempts:     3,                // 最多重试3次
    InitialInterval: 1 * time.Second,  // 初始间隔1秒
    MaxInterval:     10 * time.Second, // 最大间隔10秒
    Multiplier:      2.0,              // 指数退避倍数2
}
```

## 扩展指南

### 添加新的Saga

1. 定义业务实体
2. 实现SagaStep接口
3. 实现CompensationStep接口
4. 创建Saga定义函数
5. 编写单元测试

示例：`bot_creation_saga.go`

### 自定义仓储

实现 `Repository` 接口以支持不同的存储后端：

```go
type Repository interface {
    SaveSaga(ctx context.Context, saga *Saga) error
    FindSagaByName(ctx context.Context, name string) (*Saga, error)
    SaveExecution(ctx context.Context, execution *SagaExecution) error
    // ...
}
```

## 故障排查

### 常见问题

#### 1. Saga执行超时

**症状**: Saga长时间处于`running`状态

**解决方案**:
- 检查步骤超时设置
- 检查数据库连接
- 检查外部服务可用性

#### 2. 补偿失败

**症状**: Saga状态为`failed`，部分数据未回滚

**解决方案**:
- 检查补偿逻辑是否正确
- 实现补偿操作的幂等性
- 添加手动补偿脚本

#### 3. 状态不一致

**症状**: 执行记录与实际状态不符

**解决方案**:
- 检查数据库事务
- 实现状态校验逻辑
- 添加定期清理任务

## 相关文档

- [企业级开发规范手册](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [分布式事务设计](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-核心算法实现指南.md)
- [API设计规范](../../../../../docs/企业级功能完善与统一性设计方案/API设计规范文档.md)

## 许可证

Copyright © 2024 Coze Studio. All rights reserved.
