# BudgetAlertService 单元测试设计文档

## 📋 测试概述

**测试文件**: `budget_alert_service_test.go`
**测试服务**: `BudgetAlertService`
**测试日期**: 2025-01-01
**状态**: ✅ 代码已完成，⏳ 等待repository接口完善后执行

## 🎯 测试目标

验证BudgetAlertService的以下核心功能：
1. 多级预算告警触发（80%、95%、100%）
2. 告警去重机制（防止同一天重复发送）
3. 硬性上限处理（暂停服务 vs 自动降级）
4. 多预算周期支持（月度、季度、年度）
5. 通知渠道解析和发送
6. 批量租户预算检查
7. 错误处理和边界情况

## 📊 测试用例清单

### 1. 基础场景测试 (5个)

| 测试用例 | 测试场景 | 预期结果 |
|---------|---------|----------|
| `TestBudgetAlertService_CheckBudget_NoBudget` | 未配置预算的租户 | 不发送告警，返回nil |
| `TestBudgetAlertService_CheckBudget_HardCapDisabled` | 未启用硬性上限的预算 | 不发送告警，返回nil |
| `TestBudgetAlertService_CheckBudget_FirstThreshold` | 使用率80%（一级告警阈值） | 发送Warning级别告警 |
| `TestBudgetAlertService_CheckBudget_SecondThreshold` | 使用率95%（二级告警阈值） | 发送Critical级别告警 |
| `TestBudgetAlertService_CheckBudget_HardCapExceeded` | 超过硬性上限100% | 发送Emergency级别告警 |

### 2. 自动降级测试 (2个)

| 测试用例 | 测试场景 | 预期结果 |
|---------|---------|----------|
| `TestBudgetAlertService_CheckBudget_AutoDowngrade` | 超限且启用自动降级 | 触发模型降级，发送降级告警 |
| `TestBudgetAlertService_AutoDowngradeWithoutModels` | 启用降级但未配置模型 | 不发送降级告警，记录警告日志 |

### 3. 告警去重测试 (1个)

| 测试用例 | 测试场景 | 预期结果 |
|---------|---------|----------|
| `TestBudgetAlertService_CheckBudget_AlertDeduplication` | 同一天内第二次触发相同阈值 | 不重复发送告警 |

### 4. 多预算周期测试 (2个)

| 测试用例 | 测试场景 | 预期结果 |
|---------|---------|----------|
| `TestBudgetAlertService_CheckBudget_QuarterlyBudget` | 季度预算，使用率80% | 从季度初计算使用量并发送告警 |
| `TestBudgetAlertService_CheckBudget_YearlyBudget` | 年度预算，使用率80% | 从年初计算使用量并发送告警 |

### 5. 错误处理测试 (1组)

| 测试用例 | 测试场景 | 预期结果 |
|---------|---------|----------|
| `TestBudgetAlertService_CheckBudget_RepositoryErrors` | Repository调用失败 | 返回错误，不发送告警 |
| ├── 获取预算配置失败 | 返回error | 返回error |
| └── 获取使用量汇总失败 | 返回error | 返回error |

### 6. 批量处理测试 (1个)

| 测试用例 | 测试场景 | 预期结果 |
|---------|---------|----------|
| `TestBudgetAlertService_CheckBudgetsForAllTenants` | 多个租户，部分启用硬性上限 | 批量检查所有租户，仅对启用的发送告警 |

### 7. 告警消息构建测试 (3个)

| 测试用例 | 告警级别 | Emoji | 提示文本 |
|---------|---------|-------|---------|
| `buildAlertMessage` | Warning | ⚠️ | "请注意控制使用量" |
| `buildAlertMessage` | Critical | 🚨 | "即将达到预算上限" |
| `buildAlertMessage` | Emergency | 🛑 | "服务已暂停或降级" |

### 8. 通知渠道解析测试 (4个)

| 测试用例 | 通知渠道配置 | 预期解析结果 |
|---------|-------------|------------|
| 单渠道 - email | `["email"]` | 1个渠道 |
| 多渠道 | `["email", "sms", "webhook"]` | 3个渠道 |
| 空字符串 | `""` | 默认email |
| 无效JSON | `"invalid json"` | 默认email |

## 🏗️ 测试架构

### Mock实现

```go
// Mock Repository实现
type mockBudgetSettingsRepository struct {
    budgets map[string]*entity.BudgetSettings
    getErr  error
    listErr error
}

type mockTokenUsageSummaryRepository struct {
    totalCost float64
    getErr    error
}

type mockBudgetAlertRepository struct {
    alerts          map[uint64]*entity.BudgetAlert
    alertsSentToday map[string]bool
    createErr       error
    checkErr        error
}

// Mock NotificationService
type mockNotificationService struct {
    sendAlertCalled bool
    sendAlertErr    error
    sentAlerts      []*entity.BudgetAlert
}
```

### 测试数据示例

```go
// 一级告警测试数据
budget := &entity.BudgetSettings{
    TenantID:              "tenant_test_001",
    BudgetType:            entity.BudgetTypeMonthly,
    BudgetAmount:          1000.0,
    HardCapEnabled:        true,
    AlertThreshold1:       80,
    AlertThreshold2:       95,
    NotificationChannels:  `["email"]`,
}
summaryRepo := &mockTokenUsageSummaryRepository{
    totalCost: 800.0, // 80% 使用率
}
```

## ✅ 预期测试覆盖率

| 功能模块 | 预期覆盖率 | 测试用例数 |
|---------|-----------|----------|
| CheckBudget | 100% | 10 |
| sendAlertIfNeeded | 100% | 5 |
| handleHardCap | 100% | 3 |
| executeDowngrade | 100% | 2 |
| buildAlertMessage | 100% | 3 |
| CheckBudgetsForAllTenants | 100% | 1 |
| **总计** | **~100%** | **24+** |

## 🔍 测试执行指南

### 前置条件

1. 确保repository接口完整（包含所有必需方法）
2. 修复token_metering.go中的编译错误
3. 安装测试依赖：
   ```bash
   go get github.com/stretchr/testify
   go get github.com/shopspring/decimal
   ```

### 执行测试

```bash
# 运行所有BudgetAlertService测试
cd backend
go test -v ./domain/billing/service -run TestBudgetAlertService

# 运行特定测试
go test -v ./domain/billing/service -run TestBudgetAlertService_CheckBudget_FirstThreshold

# 生成覆盖率报告
go test -coverprofile=coverage.out ./domain/billing/service -run TestBudgetAlertService
go tool cover -html=coverage.out -o coverage.html
```

## 🚨 已知问题和解决方案

### 问题1: Repository接口不完整

**问题描述**: `TokenBudgetRepository` 和 `TokenUsageRepository` 缺少以下方法：
- `GetByTenantAndModel(ctx, tenantID, modelID, period) (*entity.TokenBudget, error)`
- `GetMonthlyUsage(ctx, tenantID, modelID) (int, error)`
- `GetUsageByModel(ctx, tenantID) ([]*entity.ModelUsage, error)`

**影响**: token_metering.go无法编译，阻塞整个service包的测试

**解决方案**:
1. 补充repository接口定义
2. 或添加exclude build tag临时排除旧代码

### 问题2: Logger未定义

**问题描述**: token_metering_async.go中使用`logger`但未导入

**解决方案**: 已修复token_metering.go，待修复token_metering_async.go

### 问题3: RecordTokenRequest重复定义

**问题描述**: types.go和token_metering.go中都定义了RecordTokenRequest

**解决方案**: 已清理types.go，保留token_metering.go中的定义

## 📈 测试质量指标

| 指标 | 目标值 | 当前值 |
|-----|--------|--------|
| 代码覆盖率 | ≥80% | 待测试 |
| 测试通过率 | 100% | 待测试 |
| 测试执行时间 | <5秒 | 待测试 |
| Mock覆盖度 | 100% | ✅ 100% |

## 🎓 测试最佳实践

本测试严格遵循以下原则：

1. **SOLID原则**：
   - S: 单一职责 - 每个测试只验证一个功能点
   - O: 开放封闭 - 易于添加新测试，无需修改现有代码
   - D: 依赖倒置 - 使用Mock接口隔离依赖

2. **测试命名规范**：
   ```go
   func Test<Struct>_<Method>_<Scenario>(t *testing.T)
   ```

3. **AAA模式**：
   - Arrange（准备）：创建Mock对象和测试数据
   - Act（执行）：调用被测试方法
   - Assert（断言）：验证预期结果

4. **表驱动测试**：
   ```go
   tests := []struct {
       name     string
       input    InputType
       expected ExpectedType
   }{
       { /* 测试用例1 */ },
       { /* 测试用例2 */ },
   }
   ```

## 📝 测试检查清单

- [x] 创建完整的Mock实现
- [x] 编写所有测试用例
- [x] 添加测试文档注释
- [x] 遵循企业级开发规范
- [ ] 修复repository接口问题
- [ ] 执行测试并验证通过率
- [ ] 生成覆盖率报告
- [ ] 性能基准测试

## 🔄 持续改进

### 未来优化方向

1. **添加性能基准测试**：
   ```go
   func BenchmarkBudgetAlertService_CheckBudget(b *testing.B) {
       // 测试1000次预算检查的性能
   }
   ```

2. **添加并发测试**：
   ```go
   func TestBudgetAlertService_ConcurrentChecks(t *testing.T) {
       // 测试并发租户预算检查的线程安全性
   }
   ```

3. **添加模糊测试**：
   ```go
   func FuzzBudgetAlertService_CheckBudget(f *testing.F) {
       // 使用随机输入测试边界情况
   }
   ```

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
**维护者**: Coze Studio Backend Team
