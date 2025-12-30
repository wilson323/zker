# 计费系统集成测试指南

## 📋 目录

1. [测试概述](#测试概述)
2. [测试环境准备](#测试环境准备)
3. [测试执行](#测试执行)
4. [测试文件说明](#测试文件说明)
5. [测试覆盖率](#测试覆盖率)
6. [性能基准](#性能基准)
7. [常见问题](#常见问题)

---

## 📖 测试概述

本测试套件为ZKER计费系统提供完整的集成测试，包括：

- ✅ **Token计量API集成测试** - 7个核心场景
- ✅ **预算管理API集成测试** - 5个核心场景
- ✅ **端到端流程测试** - 完整的10步验证
- ✅ **并发测试** - 4个并发场景
- ✅ **错误场景测试** - 8个错误处理场景

### 测试覆盖的关键功能

#### Token计量API
- ✅ RecordTokenUsage - 单次记录
- ✅ BatchRecordTokenUsage - 批量记录（最多1000条）
- ✅ GetUsageStats - 使用统计查询
- ✅ GetModelUsageStats - 模型使用统计
- ✅ GetDailyUsageStats - 每日趋势统计

#### 预算管理API
- ✅ CreateBudget - 创建预算配置
- ✅ GetBudget - 获取预算配置
- ✅ UpdateBudget - 更新预算配置
- ✅ DeleteBudget - 删除预算配置
- ✅ GetBudgetUsage - 获取预算使用情况
- ✅ CheckBudget - 手动预算检查
- ✅ GetAlerts - 查询告警历史

#### 核心组件
- ✅ PricingEngine - 成本计算引擎
- ✅ TokenMeteringService - Token计量服务
- ✅ BudgetManagementService - 预算管理服务
- ✅ BudgetAlertService - 预算告警服务
- ✅ NotificationService - 通知服务（模拟）

---

## 🛠️ 测试环境准备

### 必需软件

```bash
# 1. Go 1.24.0+
go version

# 2. Docker 20.10+
docker --version

# 3. Docker Compose（可选）
docker-compose --version
```

### 环境检查

```bash
# 检查Go环境
cd backend
go mod tidy
go build ./...

# 检查Docker环境
docker ps
```

### 测试基础设施

本测试套件使用以下技术：

- **testcontainers-go**: MySQL 8.4.5测试容器
- **testify**: 断言和测试套件
- **GORM**: ORM框架
- **logrus**: 日志记录

---

## 🚀 测试执行

### 方式1: 使用测试脚本（推荐）

```bash
# 进入测试目录
cd backend/tests/integration/billing

# 给脚本添加执行权限
chmod +x run_tests.sh

# 运行所有测试
./run_tests.sh

# 运行测试并生成覆盖率报告
./run_tests.sh --coverage

# 运行测试并生成性能基准报告
./run_tests.sh --performance

# 详细输出模式
./run_tests.sh --verbose

# 查看帮助
./run_tests.sh --help
```

### 方式2: 使用Go命令

```bash
# 进入测试目录
cd backend/tests/integration/billing

# 运行所有集成测试
go test -v -timeout 30m ./...

# 运行特定测试套件
go test -v -timeout 30m ./... -run TestTokenMetering
go test -v -timeout 30m ./... -run TestBudgetManagement
go test -v -timeout 30m ./... -run TestBilling_E2E
go test -v -timeout 30m ./... -run TestConcurrent
go test -v -timeout 30m ./... -run TestErrorScenarios

# 运行特定测试用例
go test -v -timeout 30m ./... -run TestTokenMetering_CompleteFlow

# 生成覆盖率报告
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html

# 运行性能基准测试
go test -bench=. -benchmem ./...
```

### 方式3: 在IDE中运行

**VSCode**:
1. 打开测试文件
2. 点击测试函数上方的"Run Test"按钮
3. 或者使用命令面板: `Go: Test Function`

**GoLand**:
1. 打开测试文件
2. 右键点击测试函数
3. 选择 "Run 'TestXXX'"

---

## 📁 测试文件说明

### 测试文件清单

```
backend/tests/integration/billing/
├── README.md                              # 本文档
├── run_tests.sh                           # 测试执行脚本
├── token_metering_integration_test.go     # Token计量集成测试
├── budget_management_integration_test.go  # 预算管理集成测试
├── billing_e2e_test.go                    # 端到端集成测试
├── billing_concurrent_test.go             # 并发测试
├── billing_error_scenarios_test.go        # 错误场景测试
├── billing_repository_test.go             # Repository测试
└── results/                               # 测试结果目录
    └── {timestamp}/
        ├── summary.txt                    # 测试总结报告
        ├── token_metering.log             # Token计量测试日志
        ├── budget_management.log          # 预算管理测试日志
        ├── e2e.log                        # E2E测试日志
        ├── concurrent.log                 # 并发测试日志
        ├── error_scenarios.log            # 错误场景测试日志
        ├── coverage.out                   # 覆盖率原始数据
        ├── coverage.html                  # 覆盖率HTML报告
        ├── coverage_func.txt              # 覆盖率函数报告
        └── benchmark.log                  # 性能基准报告
```

### 测试场景详情

#### 1. Token计量集成测试 (token_metering_integration_test.go)

| 测试用例 | 测试场景 | 验证点 |
|---------|---------|--------|
| `TestTokenMetering_CompleteFlow` | 完整流程测试 | 日志记录、汇总更新、成本计算 |
| `TestTokenMetering_BatchRecord` | 批量记录测试 | 1000条记录、性能验证 |
| `TestTokenMetering_MultiModelStats` | 多模型统计测试 | gpt-4、claude-3、qwen-max |
| `TestTokenMetering_DailyTrend` | 每日趋势测试 | 7天数据、趋势计算 |
| `TestTokenMetering_BudgetAlertTrigger` | 预算告警触发测试 | 80%阈值告警 |

#### 2. 预算管理集成测试 (budget_management_integration_test.go)

| 测试用例 | 测试场景 | 验证点 |
|---------|---------|--------|
| `TestBudgetManagement_CRUD` | 预算CRUD测试 | 创建、查询、更新、删除 |
| `TestBudgetManagement_UsageCalculation` | 使用率计算测试 | 80%使用率、剩余金额 |
| `TestBudgetManagement_MultiPeriod` | 多周期预算测试 | 月度、季度、年度 |
| `TestBudgetManagement_AlertHistory` | 告警历史查询测试 | 分页、过滤功能 |
| `TestBudgetManagement_ManualCheck` | 手动预算检查测试 | 告警触发、返回结果 |

#### 3. 端到端集成测试 (billing_e2e_test.go)

| 步骤 | 测试场景 | 验证点 |
|-----|---------|--------|
| 步骤1 | 创建租户和Bot | 租户和Bot创建成功 |
| 步骤2 | 配置月度预算 | 预算配置正确 |
| 步骤3 | AI模型调用记录 | Token使用记录成功 |
| 步骤4 | 成本计算验证 | 成本计算准确 |
| 步骤5 | 汇总更新验证 | 汇总数据正确 |
| 步骤6 | 预算告警验证 | 告警触发正确 |
| 步骤7 | 通知发送验证 | 通知状态更新 |
| 步骤8 | 使用统计查询 | 统计数据准确 |
| 步骤9 | 预算使用查询 | 使用率计算正确 |
| 步骤10 | 告警历史查询 | 告警记录完整 |

#### 4. 并发测试 (billing_concurrent_test.go)

| 测试用例 | 并发度 | 验证点 |
|---------|--------|--------|
| `TestConcurrent_RecordTokenUsage` | 100 goroutines × 10 records | 无数据丢失、无竞态条件 |
| `TestConcurrent_BudgetCheck` | 50 goroutines | 告警去重机制 |
| `TestConcurrent_BudgetUpdate` | 20 goroutines | 数据一致性 |
| `TestConcurrent_StatsQuery` | 30 goroutines | 无死锁、结果一致 |

#### 5. 错误场景测试 (billing_error_scenarios_test.go)

| 测试用例 | 测试场景 | 预期行为 |
|---------|---------|---------|
| `TestErrorScenarios_InvalidParameters` | 7种无效参数 | 正确捕获错误 |
| `TestErrorScenarios_BudgetNotFound` | 预算不存在 | 返回明确错误 |
| `TestErrorScenarios_InvalidDateFormat` | 无效日期格式 | 拒绝无效日期 |
| `TestErrorScenarios_BatchPartialFailure` | 批量部分失败 | 有效记录保存 |
| `TestErrorScenarios_TenantIsolation` | 租户隔离验证 | 数据正确隔离 |
| `TestErrorScenarios_TransactionRollback` | 事务回滚 | 部分失败不影响其他 |
| `TestErrorScenarios_DatabaseConnectionFailure` | 数据库连接失败 | 返回连接错误 |
| `TestErrorScenarios_LargeBatchData` | 大批量数据 | 拒绝超限批量 |

---

## 📊 测试覆盖率

### 覆盖率目标

- **集成测试覆盖率**: ≥ 70%
- **关键路径覆盖率**: 100%

### 生成覆盖率报告

```bash
# 使用测试脚本
./run_tests.sh --coverage

# 或使用Go命令
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html
```

### 查看覆盖率报告

1. **HTML报告**: 在浏览器中打开 `coverage.html`
2. **函数报告**: 查看 `coverage_func.txt`

```
github.com/coze-dev/coze-studio/backend/domain/billing/service/token_metering_service.go:75:    NewTokenMeteringService          100.0%
github.com/coze-dev/coze-studio/backend/domain/billing/service/token_metering_service.go:160:   RecordTokenUsage                100.0%
github.com/coze-dev/coze-studio/backend/domain/billing/service/token_metering_service.go:258:   BatchRecordTokenUsage           100.0%
...
total:                                                                                           87.5%
```

---

## ⚡ 性能基准

### 性能基准目标

| 操作 | 性能要求 |
|-----|---------|
| 单次Token记录 | < 100ms |
| 批量记录（1000条） | < 1秒 |
| 预算检查 | < 50ms |
| 统计查询 | < 200ms |

### 运行性能基准测试

```bash
# 使用测试脚本
./run_tests.sh --performance

# 或使用Go命令
go test -bench=. -benchmem ./...
```

### 性能基准示例

```
BenchmarkRecordTokenUsage-8     1000    1.2 ms/op    512 B/op    10 allocs/op
BenchmarkBatchRecordTokenUsage-8   10   98.5 ms/op  10240 B/op   100 allocs/op
BenchmarkCheckBudget-8           5000    0.5 ms/op    256 B/op     5 allocs/op
```

---

## 🔧 常见问题

### Q1: 测试失败：Docker未运行

**错误信息**:
```
Failed to start MySQL container: context deadline exceeded
```

**解决方案**:
```bash
# 启动Docker
sudo systemctl start docker  # Linux
open -a Docker  # macOS

# 检查Docker状态
docker ps
```

### Q2: 测试超时

**错误信息**:
```
test timed out after 30m
```

**解决方案**:
```bash
# 增加超时时间
go test -timeout 60m ./...

# 或者只运行特定测试
go test -run TestTokenMetering_CompleteFlow ./...
```

### Q3: MySQL容器启动失败

**错误信息**:
```
container exited with code 1
```

**解决方案**:
```bash
# 清理旧的容器
docker system prune -f

# 检查端口占用
netstat -tuln | grep 3306

# 重新运行测试
go test ./...
```

### Q4: 测试数据未清理

**解决方案**:
```bash
# 测试容器会自动清理，如果未清理：
docker ps -a | grep mysql | awk '{print $1}' | xargs docker rm -f
```

### Q5: 并发测试偶尔失败

**解决方案**:
```bash
# 增加等待时间
# 或重新运行特定测试
go test -v -run TestConcurrent_RecordTokenUsage ./...
```

---

## 📞 获取帮助

如有问题，请查阅：

1. **项目文档**: `docs/企业级功能完善与统一性设计方案/`
2. **开发规范**: `docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md`
3. **GitHub Issues**: [提交问题](https://github.com/coze-dev/coze-studio/issues)

---

**最后更新**: 2025-01-01
**维护者**: 研发B（后端工程师）
