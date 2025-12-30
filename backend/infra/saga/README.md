# Saga分布式事务框架使用指南

> **版本**: v1.0
> **创建日期**: 2025-12-30
> **作者**: ZKER架构团队

## 📖 目录

- [1. 概述](#1-概述)
- [2. 快速开始](#2-快速开始)
- [3. 核心概念](#3-核心概念)
- [4. 使用示例](#4-使用示例)
- [5. 最佳实践](#5-最佳实践)
- [6. 测试指南](#6-测试指南)

---

## 1. 概述

### 1.1 什么是Saga模式?

Saga模式是一种长活事务(Long Lived Transaction)模式,用于管理分布式事务。它将一个分布式事务拆分为多个本地事务,每个本地事务都有对应的补偿事务。

**核心思想**:
1. 将分布式事务拆分为多个本地事务
2. 按顺序执行每个本地事务
3. 如果某个本地事务失败,执行前面已完成的本地事务的补偿操作
4. 最终达到一致状态

### 1.2 适用场景

✅ **适合使用Saga的场景**:
- 跨服务业务流程(租户注册、Bot创建)
- 长事务(秒级到分钟级)
- 可以接受最终一致性
- 需要补偿回滚

❌ **不适合使用Saga的场景**:
- 强一致性要求(如支付)
- 短事务(毫秒级)
- 无法补偿的操作(如发送短信后无法回收)

---

## 2. 快速开始

### 2.1 初始化数据库

```sql
-- 执行迁移脚本
mysql -u root -p zker_db < backend/migrations/saga/001_create_saga_tables.sql
```

### 2.2 创建Saga协调器

```go
package main

import (
    "context"
    "go.uber.org/zap"
    "gorm.io/gorm"

    "your-project/backend/infra/saga"
)

func main() {
    // 1. 初始化依赖
    var db *gorm.DB
    logger := zap.NewExample()

    // 2. 创建Saga仓储
    repo := saga.NewMySQLSagaRepository(db)

    // 3. 自动迁移表结构
    if err := repo.AutoMigrate(); err != nil {
        logger.Fatal("迁移表结构失败", zap.Error(err))
    }

    // 4. 创建Saga协调器
    orchestrator := saga.NewSagaOrchestrator(repo, logger)

    // 5. 定义和执行Saga...
}
```

---

## 3. 核心概念

### 3.1 Saga定义

Saga是完整的分布式事务定义,包含:
- **ID**: 唯一标识
- **Name**: Saga名称
- **Steps**: 执行步骤列表
- **Compensations**: 补偿步骤列表
- **Timeout**: 超时时间
- **RetryPolicy**: 重试策略

### 3.2 Saga步骤

SagaStep是单个执行步骤,必须实现以下方法:
- `Execute(ctx, data)`: 执行步骤逻辑
- `Name()`: 返回步骤名称
- `Timeout()`: 返回步骤超时时间

### 3.3 补偿步骤

CompensationStep是补偿操作,必须实现以下方法:
- `Compensate(ctx, data)`: 执行补偿逻辑
- `Name()`: 返回补偿步骤名称
- `Timeout()`: 返回补偿超时时间

### 3.4 Saga执行记录

SagaExecution记录了Saga的执行过程:
- 执行状态
- 当前步骤
- 每个步骤的执行结果
- 输入输出数据
- 错误信息

---

## 4. 使用示例

### 4.1 定义简单Saga

```go
package main

import (
    "context"
    "fmt"
    "time"
    "your-project/backend/infra/saga"
)

// 步骤1: 创建订单
type CreateOrderStep struct{}

func (s *CreateOrderStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    req := data.(*CreateOrderRequest)
    order := &Order{
        ID:     generateID(),
        UserID: req.UserID,
        Amount: req.Amount,
    }
    // 保存订单到数据库
    if err := orderRepo.Save(ctx, order); err != nil {
        return nil, fmt.Errorf("创建订单失败: %w", err)
    }
    return order, nil
}

func (s *CreateOrderStep) Name() string {
    return "创建订单"
}

func (s *CreateOrderStep) Timeout() time.Duration {
    return 5 * time.Second
}

// 补偿1: 取消订单
type CancelOrderStep struct{}

func (s *CancelOrderStep) Compensate(ctx context.Context, data interface{}) error {
    order := data.(*Order)
    // 取消订单
    if err := orderRepo.Cancel(ctx, order.ID); err != nil {
        return fmt.Errorf("取消订单失败: %w", err)
    }
    return nil
}

func (s *CancelOrderStep) Name() string {
    return "取消订单"
}

func (s *CancelOrderStep) Timeout() time.Duration {
    return 5 * time.Second
}

// 步骤2: 扣减库存
type DeductInventoryStep struct{}

func (s *DeductInventoryStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    order := data.(*Order)
    // 扣减库存
    if err := inventoryRepo.Deduct(ctx, order.ProductID, order.Quantity); err != nil {
        return nil, fmt.Errorf("扣减库存失败: %w", err)
    }
    return order, nil
}

func (s *DeductInventoryStep) Name() string {
    return "扣减库存"
}

func (s *DeductInventoryStep) Timeout() time.Duration {
    return 5 * time.Second
}

// 补偿2: 恢复库存
type RestoreInventoryStep struct{}

func (s *RestoreInventoryStep) Compensate(ctx context.Context, data interface{}) error {
    order := data.(*Order)
    // 恢复库存
    if err := inventoryRepo.Restore(ctx, order.ProductID, order.Quantity); err != nil {
        return fmt.Errorf("恢复库存失败: %w", err)
    }
    return nil
}

func (s *RestoreInventoryStep) Name() string {
    return "恢复库存"
}

func (s *RestoreInventoryStep) Timeout() time.Duration {
    return 5 * time.Second
}

// 定义订单创建Saga
func NewOrderCreationSaga() *saga.Saga {
    return &saga.Saga{
        ID:          "saga-order-creation",
        Name:        "订单创建Saga",
        Description: "处理订单创建流程,包括创建订单、扣减库存",
        Steps: []saga.SagaStep{
            &CreateOrderStep{},
            &DeductInventoryStep{},
        },
        Compensations: []saga.CompensationStep{
            &CancelOrderStep{},
            &RestoreInventoryStep{},
        },
        Timeout:     30 * time.Second,
        RetryPolicy: saga.DefaultRetryPolicy(),
    }
}
```

### 4.2 执行Saga

```go
package main

import (
    "context"
    "log"
)

func main() {
    // 初始化协调器(见快速开始)
    orchestrator := setupOrchestrator()

    // 注册Saga
    orderSaga := NewOrderCreationSaga()
    if err := orchestrator.DefineSaga(orderSaga); err != nil {
        log.Fatalf("定义Saga失败: %v", err)
    }

    // 执行Saga
    ctx := context.Background()
    req := &CreateOrderRequest{
        UserID:    "user123",
        ProductID: "product456",
        Amount:    100,
        Quantity:  1,
    }

    execution, err := orchestrator.ExecuteSaga(ctx, "订单创建Saga", req)
    if err != nil {
        log.Printf("Saga执行失败: %v", err)
        // 检查是否已补偿
        if execution.Status == saga.SagaStatusCompensated {
            log.Printf("Saga已自动补偿,所有步骤已回滚")
        }
        return
    }

    log.Printf("Saga执行成功: %s", execution.ID)
    log.Printf("执行状态: %s", execution.Status)
    log.Printf("步骤数量: %d", len(execution.StepExecutions))
}
```

### 4.3 查询Saga状态

```go
// 查询执行状态
execution, err := orchestrator.GetStatus(ctx, executionID)
if err != nil {
    log.Printf("查询Saga状态失败: %v", err)
    return
}

log.Printf("Saga状态: %s", execution.Status)
log.Printf("当前步骤: %d", execution.CurrentStep)
log.Printf("开始时间: %s", execution.StartedAt)

// 打印每个步骤的执行情况
for i, stepExec := range execution.StepExecutions {
    log.Printf("步骤%d: %s, 状态: %s", i+1, stepExec.StepName, stepExec.Status)
    if stepExec.Error != nil {
        log.Printf("  错误: %v", stepExec.Error)
    }
}
```

---

## 5. 最佳实践

### 5.1 幂等性

**所有步骤和补偿操作都必须是幂等的**:

```go
// ✅ Good: 幂等操作
func (s *CreateOrderStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    order := &Order{
        ID: generateID(), // 使用唯一ID
        // ...
    }

    // 使用INSERT IGNORE或ON DUPLICATE KEY UPDATE
    err := orderRepo.SaveWithID(ctx, order.ID, order)
    if err != nil {
        return nil, err
    }
    return order, nil
}

// ❌ Bad: 非幂等操作
func (s *CreateOrderStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    // 每次执行都会创建新订单
    order := orderRepo.CreateNew()
    return order, nil
}
```

### 5.2 超时控制

**为每个步骤设置合理的超时时间**:

```go
func (s *CreateOrderStep) Timeout() time.Duration {
    // 根据实际业务需求设置
    // - 快速操作(数据库写): 1-5秒
    // - 中等操作(调用外部API): 10-30秒
    // - 慢速操作(批量处理): 1-5分钟
    return 5 * time.Second
}
```

### 5.3 补偿策略

**补偿操作要尽可能简单可靠**:

```go
// ✅ Good: 简单直接的补偿
func (s *CancelOrderStep) Compensate(ctx context.Context, data interface{}) error {
    order := data.(*Order)
    // 直接更新状态
    return orderRepo.UpdateStatus(ctx, order.ID, "cancelled")
}

// ❌ Bad: 复杂的补偿逻辑
func (s *CancelOrderStep) Compensate(ctx context.Context, data interface{}) error {
    order := data.(*Order)
    // 发送通知、更新日志、清理缓存等复杂操作
    // 补偿失败率高
}
```

### 5.4 错误处理

**明确区分可恢复和不可恢复的错误**:

```go
func (s *DeductInventoryStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    order := data.(*Order)

    // 检查库存
    available, err := inventoryRepo.GetAvailable(ctx, order.ProductID)
    if err != nil {
        // 数据库错误,可以重试
        return nil, &saga.SagaError{
            Code:        "DATABASE_ERROR",
            Message:     "查询库存失败",
            Cause:       err,
            Recoverable: true,
        }
    }

    if available < order.Quantity {
        // 库存不足,不可恢复
        return nil, &saga.SagaError{
            Code:        "INSUFFICIENT_INVENTORY",
            Message:     "库存不足",
            Recoverable: false,
        }
    }

    // 扣减库存
    return order, nil
}
```

### 5.5 日志记录

**记录详细的执行日志**:

```go
func (s *CreateOrderStep) Execute(ctx context.Context, data interface{}) (interface{}, error) {
    req := data.(*CreateOrderRequest)

    logger := zap.L().With(
        zap.String("step", "CreateOrder"),
        zap.String("user_id", req.UserID),
        zap.Float64("amount", req.Amount),
    )

    logger.Info("开始创建订单")

    order := &Order{
        ID:     generateID(),
        UserID: req.UserID,
        Amount: req.Amount,
    }

    if err := orderRepo.Save(ctx, order); err != nil {
        logger.Error("创建订单失败", zap.Error(err))
        return nil, err
    }

    logger.Info("创建订单成功", zap.String("order_id", order.ID))
    return order, nil
}
```

---

## 6. 测试指南

### 6.1 单元测试

```go
package saga_test

import (
    "context"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
    "your-project/backend/infra/saga"
)

func TestOrderCreationSaga_Success(t *testing.T) {
    // 准备测试环境
    repo := NewMockRepository()
    logger := zap.NewNop()
    orchestrator := saga.NewSagaOrchestrator(repo, logger)

    // 定义并注册Saga
    orderSaga := NewOrderCreationSaga()
    err := orchestrator.DefineSaga(orderSaga)
    assert.NoError(t, err)

    // 执行Saga
    ctx := context.Background()
    req := &CreateOrderRequest{
        UserID:    "user123",
        ProductID: "product456",
        Amount:    100,
        Quantity:  1,
    }

    execution, err := orchestrator.ExecuteSaga(ctx, "订单创建Saga", req)

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, saga.SagaStatusCompleted, execution.Status)
    assert.Equal(t, 2, len(execution.StepExecutions))

    // 验证每个步骤都成功
    for _, stepExec := range execution.StepExecutions {
        assert.Equal(t, saga.StepCompleted, stepExec.Status)
        assert.NoError(t, stepExec.Error)
    }
}

func TestOrderCreationSaga_Compensation(t *testing.T) {
    // 准备测试环境
    repo := NewMockRepository()
    logger := zap.NewNop()
    orchestrator := saga.NewSagaOrchestrator(repo, logger)

    // 定义会失败的Saga
    // ...

    // 执行并验证补偿
    // ...
}
```

### 6.2 集成测试

```bash
# 运行集成测试(需要MySQL)
go test -tags=integration -v ./backend/infra/saga/
```

### 6.3 性能测试

```go
func BenchmarkSagaExecution(b *testing.B) {
    orchestrator := setupOrchestrator()
    saga := NewOrderCreationSaga()
    orchestrator.DefineSaga(saga)

    ctx := context.Background()
    req := &CreateOrderRequest{...}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = orchestrator.ExecuteSaga(ctx, "订单创建Saga", req)
    }
}
```

---

## 7. 监控和运维

### 7.1 监控指标

建议监控以下指标:
- Saga执行总数
- 成功/失败/补偿的执行数量
- 平均执行时长
- 步骤失败率
- 补偿成功率

### 7.2 告警规则

建议配置以下告警:
- Saga失败率超过5%
- 补偿失败
- Saga执行超过预期时间
- 步骤超时次数过多

### 7.3 日志查询

查询Saga执行日志:
```sql
-- 查询失败的Saga执行记录
SELECT id, saga_id, status, error, started_at
FROM saga_executions
WHERE status = 'failed'
ORDER BY started_at DESC
LIMIT 100;

-- 查询特定Saga的所有执行记录
SELECT id, status, current_step, started_at, completed_at
FROM saga_executions
WHERE saga_id = 'saga-order-creation'
ORDER BY started_at DESC;
```

---

## 8. 故障排查

### 8.1 Saga卡在Running状态

**原因**: 可能是某个步骤执行超时或程序崩溃

**解决方案**:
```sql
-- 查询Running状态的执行记录
SELECT id, saga_id, current_step, started_at
FROM saga_executions
WHERE status = 'running'
AND started_at < NOW() - INTERVAL 30 MINUTE;

-- 手动标记为失败,触发人工介入
UPDATE saga_executions
SET status = 'failed',
    error = '手动标记: 执行超时',
    updated_at = NOW()
WHERE id = 'execution-id';
```

### 8.2 补偿执行失败

**原因**: 补偿步骤执行失败或超时

**解决方案**:
1. 检查补偿步骤的日志
2. 手动执行补偿操作
3. 修复后重新标记为Compensated

---

## 附录

### A. 常见问题

**Q1: Saga和2PC有什么区别?**
A: Saga是最终一致性,性能更好;2PC是强一致性,性能较差。

**Q2: 如何保证补偿一定成功?**
A: 通过重试机制+人工介入保证最终成功。

**Q3: 补偿步骤执行失败怎么办?**
A: 记录失败状态,发送告警,等待人工介入处理。

### B. 相关文档

- [Saga模式设计文档](../../docs/开发规范/分布式事务方案-Saga模式.md)
- [企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

**文档版本**: v1.0
**最后更新**: 2025-12-30
**作者**: ZKER架构团队
