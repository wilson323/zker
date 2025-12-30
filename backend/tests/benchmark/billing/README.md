# Token计量和预算管理系统 - 性能基准测试文档

**版本**: v1.0.0
**最后更新**: 2025-12-30
**负责人**: 研发B

---

## 📋 目录

- [测试概述](#测试概述)
- [性能目标](#性能目标)
- [测试覆盖范围](#测试覆盖范围)
- [快速开始](#快速开始)
- [测试执行](#测试执行)
- [结果分析](#结果分析)
- [性能优化建议](#性能优化建议)
- [常见问题](#常见问题)

---

## 🎯 测试概述

### 测试目的

本性能基准测试套件旨在全面评估Token计量和预算管理系统的性能表现，确保：

1. ✅ **高性能**: 满足企业级应用的性能要求
2. ✅ **可扩展**: 支持租户数量和业务量增长
3. ✅ **稳定性**: 长时间运行无性能退化
4. ✅ **资源高效**: CPU和内存使用率在合理范围

### 测试范围

**测试模块**:
- Token计量服务（单次记录、批量记录）
- 定价引擎（成本计算、节省计算）
- 使用统计查询（不同数据量、时间范围）
- 预算管理服务（预算检查、使用情况）
- 数据库操作（批量插入、汇总查询）
- 并发性能（多goroutine并发）
- 内存泄漏检测

---

## 📊 性能目标

### 基线指标（Baseline）

| 操作 | 目标性能 | 说明 |
|------|---------|------|
| **RecordTokenUsage** | < 100ms/次 | 单次Token记录 |
| **BatchRecordTokenUsage (1000条)** | < 1秒 | 批量记录 |
| **GetUsageStats** | < 50ms | 使用统计查询 |
| **CheckBudget** | < 50ms | 预算检查 |
| **PricingEngine (10000次)** | < 100ms | 定价计算 |
| **并发吞吐量** | > 1000 req/s | 并发性能 |

### 资源使用限制

| 资源 | 限制 | 说明 |
|------|------|------|
| **CPU** | < 80% (单核) | 正常负载下 |
| **内存** | < 500MB | 测试过程中的内存增长 |
| **数据库连接** | < 100 | 并发连接数 |
| **goroutine数量** | < 1000 | 避免goroutine泄漏 |

---

## 🧪 测试覆盖范围

### Benchmark 1: Token计量单次记录性能

**测试函数**: `BenchmarkRecordTokenUsage`

**性能目标**:
- ops/s: > 10 次/秒
- 每次操作: < 100ms
- 内存分配: 最小化

**测试内容**:
- Token使用记录创建
- 成本计算
- 汇总数据更新
- 预算告警检查

### Benchmark 2: Token计量批量记录性能

**测试函数**: `BenchmarkBatchRecordTokenUsage`

**性能目标**:
- 批量大小: 10, 50, 100, 500, 1000
- 1000条记录: < 1秒
- 平均每条: < 1ms

**测试场景**:
- 不同批次大小的性能对比
- 批量操作的线性扩展性

### Benchmark 3: Token使用统计查询性能

**测试函数**: `BenchmarkGetUsageStats`

**性能目标**:
- ops/s: > 20 次/秒
- 每次查询: < 50ms

**测试场景**:
- 不同数据量: 1K, 10K, 100K条记录
- 聚合查询性能
- 时间范围查询性能

### Benchmark 4: 定价引擎性能测试

**测试函数**: `BenchmarkPricingEngine_CalculateCost`

**性能目标**:
- 10000次计算: < 100ms
- 平均每次: < 0.01ms
- 零内存分配（栈上计算）

**测试模型**: 13个AI模型
- OpenAI: GPT-4, GPT-3.5-Turbo
- Anthropic: Claude-3-Opus, Claude-3-Haiku
- Qwen: Qwen-Max, Qwen-Plus, Qwen-Turbo
- Baidu: ERNIE-Bot-4
- Zhipu: ChatGLM-Pro, ChatGLM-Turbo
- DeepSeek: DeepSeek-Chat
- Baichuan: Baichuan2-Turbo
- MiniMax: ABAB5.5-Chat

### Benchmark 5-7: 预算管理性能测试

**测试函数**:
- `BenchmarkCheckBudget`: 单租户预算检查
- `BenchmarkGetBudgetUsage`: 预算使用情况查询
- `BenchmarkCheckBudgetsForAllTenants`: 批量预算检查

**性能目标**:
- 单次检查: < 50ms
- 使用情况查询: < 100ms
- 100个租户批量检查: < 5秒

### Benchmark 8-9: 数据库操作性能测试

**测试函数**:
- `BenchmarkTokenUsageLogBatchCreate`: 批量插入性能
- `BenchmarkTokenUsageSummaryGetTotalCost`: 汇总查询性能

**性能目标**:
- 测试不同批次大小: 10, 50, 100, 500
- 不同时间范围: 日、周、月
- 验证索引效果

### Benchmark 10-11: 并发性能测试

**测试函数**:
- `BenchmarkConcurrentRecordTokenUsage`: 并发Token记录
- `BenchmarkConcurrentCheckBudget`: 并发预算检查

**性能目标**:
- 测试并发级别: 10, 50, 100, 500
- 无竞态条件
- 高吞吐量

### Profile 1-3: 性能分析

**测试函数**:
- `BenchmarkCPUProfile`: CPU性能分析
- `BenchmarkMemoryProfile`: 内存性能分析
- `TestMemoryLeak`: 内存泄漏检测

**分析工具**: go tool pprof

**用法**:
```bash
go test -cpuprofile=cpu.prof -bench=BenchmarkCPUProfile ./
go tool pprof cpu.prof
```

### Stress Test 1-2: 压力测试

**测试函数**:
- `TestStressTokenRecording`: 持续高负载测试
- `TestPeakLoad`: 峰值负载测试

**测试场景**:
- 持续1小时，100 QPS
- 短时间峰值，1000 QPS
- 验证系统稳定性和错误率

---

## 🚀 快速开始

### 前置要求

1. Go 1.24.0+
2. GNU Bash 4.0+
3. pprof (可选，用于profile分析)

### 安装pprof

```bash
go install github.com/google/pprof@latest
```

### 克隆项目

```bash
cd /path/to/coze-studio/backend/tests/benchmark/billing
```

---

## 📝 测试执行

### 方式1: 使用脚本（推荐）

**快速测试** (不包含压力测试):
```bash
./run_benchmark.sh quick
```

**完整测试** (包含压力测试，耗时较长):
```bash
./run_benchmark.sh full
```

**只生成Profile**:
```bash
./run_benchmark.sh profile
```

**只运行压力测试**:
```bash
./run_benchmark.sh stress
```

### 方式2: 手动执行

**运行所有基准测试**:
```bash
cd backend/tests/benchmark/billing

# 运行所有基准测试
go test -bench=. -benchmem -run=^$ ./

# 运行特定测试
go test -bench=BenchmarkRecordTokenUsage -benchmem ./
```

**生成CPU Profile**:
```bash
go test -cpuprofile=cpu.prof -bench=BenchmarkCPUProfile ./
go tool pprof cpu.prof
```

**生成Memory Profile**:
```bash
go test -memprofile=mem.prof -bench=BenchmarkMemoryProfile ./
go tool pprof mem.prof
```

**运行内存泄漏检测**:
```bash
go test -v -run=TestMemoryLeak ./
```

**运行压力测试**:
```bash
# 完整压力测试（包含持续高负载和峰值负载）
go test -v -run="TestStressTokenRecording|TestPeakLoad" ./

# 只运行持续高负载测试
go test -v -run=TestStressTokenRecording ./

# 只运行峰值负载测试
go test -v -run=TestPeakLoad ./

# 跳过压力测试（使用-short标志）
go test -short -v -run="TestStressTokenRecording|TestPeakLoad" ./
```

---

## 📈 结果分析

### 方式1: 使用分析脚本

```bash
./analyze_results.sh benchmark_20251230_120000.txt
```

### 方式2: 手动分析

**查看基准测试结果**:
```bash
# 运行测试并保存结果
go test -bench=. -benchmem . > results.txt 2>&1

# 查看结果
cat results.txt
```

**Profile分析**:
```bash
# CPU Profile分析
go tool pprof cpu.prof
# 在pprof交互式命令中:
(pprof) top10          # 查看Top 10热点
(pprof) list funcName  # 查看特定函数
(pprof) web            # 生成可视化图表

# Memory Profile分析
go tool pprof mem.prof
# 在pprof交互式命令中:
(pprof) top            # 查看Top 10内存分配
(pprof) list funcName  # 查看特定函数的内存分配
(pprof) web            # 生成可视化图表
```

### 性能报告解读

**基准测试输出格式**:
```
BenchmarkRecordTokenUsage-8      1000    1234567 ns/op    1234 B/op    12 allocs/op
```

**字段说明**:
- `-8`: GOMAXPROCS数量
- `1000`: 迭代次数
- `1234567 ns/op`: 每次操作耗时（纳秒） = 1.23ms/op
- `1234 B/op`: 每次操作内存分配（字节）
- `12 allocs/op`: 每次操作内存分配次数

**性能判断标准**:
- ✅ **优秀**: ns/op < 基线值
- ⚠️ **良好**: ns/op 在基线值的1-2倍之间
- ❌ **需要优化**: ns/op > 基线值的2倍

---

## 🔧 性能优化建议

### 1. 数据库优化

**问题**: 慢查询
```
症状: 某些API响应时间长，数据库日志出现慢查询

解决方案:
1. 添加合适的索引
   CREATE INDEX idx_tenant_created ON token_usage_logs(tenant_id, created_at);

2. 使用覆盖索引
   CREATE INDEX idx_cover ON token_usage_logs(tenant_id, total_cost, created_at);

3. 避免SELECT *
   SELECT id, tenant_id, total_cost FROM token_usage_logs WHERE ...;
```

**问题**: 连接池耗尽
```
症状: API返回 "Connection pool exhausted" 错误

解决方案:
1. 增加连接池大小
   db.SetMaxOpenConns(200)

2. 优化连接生命周期
   db.SetConnMaxLifetime(1 * time.Hour)

3. 检查连接泄漏
   defer db.Close()
```

### 2. 缓存优化

**实现多级缓存**:
```go
// L1: 本地缓存（内存）
// L2: Redis缓存
// L3: 数据库

func GetBot(ctx context.Context, botID string) (*Bot, error) {
    // L1: 本地缓存
    if bot, ok := localCache.Get(botID); ok {
        return bot, nil
    }

    // L2: Redis缓存
    val, err := redisClient.Get(ctx, "bot:"+botID).Result()
    if err == nil {
        var bot Bot
        json.Unmarshal([]byte(val), &bot)
        localCache.Set(botID, &bot, 5*time.Minute)
        return &bot, nil
    }

    // L3: 数据库
    bot, err := botRepo.FindByID(ctx, botID)
    if err != nil {
        return nil, err
    }

    // 回写缓存
    data, _ := json.Marshal(bot)
    redisClient.Set(ctx, "bot:"+botID, data, 30*time.Minute)
    localCache.Set(botID, bot, 5*time.Minute)

    return bot, nil
}
```

### 3. 并发优化

**限制并发数**:
```go
// 使用信号量限制并发
sem := make(chan struct{}, 100)
var wg sync.WaitGroup

for _, id := range botIDs {
    wg.Add(1)
    sem <- struct{}{}
    go func(botID string) {
        defer wg.Done()
        defer func() { <-sem }()
        processBot(botID)
    }(id)
}
wg.Wait()
```

### 4. 内存优化

**使用对象池**:
```go
var botPool = sync.Pool{
    New: func() interface{} {
        return &Bot{}
    },
}

func getBotFromPool() *Bot {
    return botPool.Get().(*Bot)
}

func putBotToPool(bot *Bot) {
    bot.Reset()
    botPool.Put(bot)
}
```

---

## ❓ 常见问题

### Q1: 测试失败怎么办？

**A**: 检查以下几点：
1. 确保所有依赖已安装（Go 1.24.0+）
2. 确保在正确的目录下执行测试
3. 查看错误日志，确定失败原因
4. 如果是数据库连接问题，检查数据库服务是否运行

### Q2: 压力测试耗时太长？

**A**: 压力测试设计为长时间运行，可以：
1. 使用`-short`标志跳过压力测试：`go test -short ./`
2. 修改测试参数，减少测试时间
3. 只运行特定的测试：`go test -run=BenchmarkXXX ./`

### Q3: 如何生成可视化图表？

**A**:
1. 使用pprof生成火焰图：`pprof -web cpu.prof`
2. 使用gnuplot绘制性能趋势图
3. 使用Grafana创建实时监控大盘

### Q4: 测试结果如何对比基线？

**A**:
1. 使用分析脚本：`./analyze_results.sh benchmark_result.txt`
2. 查看生成的基线对比报告
3. 重点关注未达标的测试项

### Q5: 如何添加新的基准测试？

**A**:
1. 在`billing_benchmark_test.go`中添加新的测试函数
2. 遵循命名规范：`BenchmarkXXX`
3. 使用`b.ResetTimer()`和`b.ReportAllocs()`
4. 添加性能目标注释

示例：
```go
// BenchmarkNewFeature 新功能性能测试
//
// 性能目标：
// - ops/s: > 100 次/秒
// - 每次操作: < 10ms
func BenchmarkNewFeature(b *testing.B) {
    // 初始化
    svc := setupService()
    ctx := context.Background()

    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        // 测试逻辑
        _, err := svc.DoSomething(ctx, i)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## 📚 相关文档

- [ZKER性能测试计划](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-性能测试计划_v1.0.md)
- [ZKER性能基线文档](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-性能基线文档.md)
- [企业级开发规范手册](../../../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

## 📝 更新日志

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|---------|------|
| v1.0.0 | 2025-12-30 | 初始版本，完整的性能基准测试框架 | 研发B |

---

**© 2025 Coze Studio. All rights reserved.**
