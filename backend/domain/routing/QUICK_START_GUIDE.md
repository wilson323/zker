# 智能路由引擎增强 - 快速入门指南

## 概述

本指南帮助开发者快速理解和使用智能路由引擎增强功能。

## 核心概念

### 1. 意图识别 (Intent Recognition)

意图识别是从用户输入中理解用户想要做什么的过程。

**示例**:
```
输入: "查询北京天气"
识别: Intent=query_weather, Confidence=0.95
路由: Agent=weather_agent
```

**三种识别策略**:
- **LLM策略**: 使用大语言模型理解语义
- **Vector策略**: 使用向量相似度匹配示例
- **Hybrid策略**: 混合两种策略，提高准确率

### 2. 实体提取 (Entity Extraction)

实体提取是从用户输入中提取关键信息。

**示例**:
```
输入: "明天从北京到上海的机票"
提取: [
  {time: 明天},
  {location: 北京},
  {location: 上海},
  {action: 机票}
]
```

### 3. 路由优化 (Routing Optimization)

路由优化根据多维度评分选择最优的Agent。

**评分维度**:
- 负载 (25%): Agent当前负载越低，分数越高
- 健康度 (20%): Agent成功率越高，分数越高
- 响应时间 (20%): 响应时间越短，分数越高
- 错误率 (15%): 错误率越低，分数越高
- 成本 (10%): 成本越低，分数越高
- 用户偏好 (10%): 用户历史使用偏好

### 4. A/B测试 (A/B Testing)

A/B测试用于对比不同路由策略的效果。

**流程**:
1. 创建测试（定义策略A和策略B）
2. 执行测试（按流量分配路由）
3. 收集数据（响应时间、用户评分）
4. 分析结果（统计显著性判断）
5. 选择获胜者

### 5. 路由学习 (Routing Learning)

路由学习从历史数据中学习，自动优化路由策略。

**学习能力**:
- 分析意图分布
- 评估Agent性能
- 生成路由建议
- 优化评分权重

## 快速开始

### 步骤1: 创建意图

```go
package main

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/domain/routing/entity"
)

// 创建意图
intent := &entity.Intent{
    TenantID:    "tenant-123",
    IntentName:  "query_weather",
    Description: "查询天气信息",
    AgentID:     "weather-agent-123",
    Confidence:  0.8,
}

// 设置训练样本
examples := []string{
    "北京天气怎么样",
    "上海今天下雨吗",
    "明天深圳气温多少",
}
intent.SetExamples(examples)

// 保存到数据库
err := intentRepo.Create(context.Background(), intent)
```

### 步骤2: 训练意图模型

```go
// 训练意图（生成向量嵌入）
err := intentRecognitionService.TrainIntentModel(
    context.Background(),
    intent.IntentID,
    examples,
)

// 训练完成后，意图可以用于识别
```

### 步骤3: 识别用户意图

```go
// 用户输入
text := "查询北京天气"

// 识别意图
result, err := intentRecognitionService.RecognizeIntent(
    context.Background(),
    "tenant-123",
    text,
)

if err != nil {
    // 处理错误
}

fmt.Printf("意图: %s, 置信度: %.2f\n",
    result.IntentName, result.Confidence)
fmt.Printf("路由到: %s\n", result.AgentID)
```

### 步骤4: 提取实体

```go
// 提取实体
entities, err := entityExtractor.ExtractEntities(
    context.Background(),
    "明天从北京到上海的机票",
)

for _, entity := range entities {
    fmt.Printf("%s: %s (置信度: %.2f)\n",
        entity.EntityName, entity.EntityValue, entity.Confidence)
}

// 输出:
// time: 明天 (置信度: 0.95)
// location: 北京 (置信度: 0.98)
// location: 上海 (置信度: 0.98)
```

### 步骤5: 优化路由决策

```go
// 构建路由请求
req := &RoutingRequest{
    TenantID: "tenant-123",
    UserID:   "user-123",
    Query:    "查询北京天气",
    Context: map[string]string{
        "device": "mobile",
        "region": "beijing",
    },
    Preferences: map[string]string{
        "preferred_agent": "fast-agent", // 用户偏好
    },
}

// 优化路由
decision, err := routingOptimizer.OptimizeRouting(
    context.Background(),
    req,
)

if err != nil {
    // 处理错误
}

fmt.Printf("最优Agent: %s\n", decision.AgentID)
fmt.Printf("得分: %.2f\n", decision.Score)
fmt.Printf("原因: %v\n", decision.Reasons)
```

### 步骤6: 创建A/B测试

```go
// 创建A/B测试
test := &entity.ABTest{
    TenantID:    "tenant-123",
    TestName:    "负载均衡 vs 评分路由",
    Description: "对比两种路由策略的效果",
    StrategyA:   "load_balance",
    StrategyB:   "score_based",
    TrafficSplit: 50, // 50:50分配
    SampleSize:  1000, // 需要收集1000个样本
}

err := abTestService.CreateABTest(context.Background(), test)

// 测试自动开始运行
```

### 步骤7: 执行A/B测试路由

```go
// 用户请求到来时，执行A/B测试
decision, err := abTestService.ExecuteABTest(
    context.Background(),
    test.TestID,
    req,
)

// decision.Strategy 会是 "load_balance" 或 "score_based"
// 根据策略进行路由
fmt.Printf("使用策略: %s\n", decision.Strategy)
```

### 步骤8: 分析A/B测试结果

```go
// 获取测试结果
result, err := abTestService.GetABTestResults(
    context.Background(),
    test.TestID,
)

fmt.Printf("策略A平均响应时间: %.2fms\n",
    result.StrategyAStats.AvgResponseTime)
fmt.Printf("策略B平均响应时间: %.2fms\n",
    result.StrategyBStats.AvgResponseTime)

fmt.Printf("统计显著性: %v\n",
    result.IsStatisticallySignificant)
fmt.Printf("推荐: %s\n", result.Recommendation)

// 如果测试完成，选择获胜者
if result.IsStatisticallySignificant {
    err = abTestService.ConcludeABTest(
        context.Background(),
        test.TestID,
        result.Winner,
    )
}
```

### 步骤9: 学习和优化

```go
// 分析最近30天的路由日志
insight, err := routingLearningService.LearnFromRoutingLogs(
    context.Background(),
    "tenant-123",
    TimeRange{
        StartTime: time.Now().Add(-30 * 24 * time.Hour).Unix(),
        EndTime:   time.Now().Unix(),
    },
)

fmt.Printf("总请求数: %d\n", insight.TotalRequests)
fmt.Printf("成功率: %.2f%%\n", insight.SuccessRate*100)
fmt.Printf("平均响应时间: %.2fms\n", insight.AvgResponseTime)

// 查看推荐
for _, rec := range insight.Recommendations {
    fmt.Printf("推荐: %s\n", rec)
}

// 获取路由规则建议
suggestions, err := routingLearningService.SuggestRoutingRules(
    context.Background(),
    "tenant-123",
)

for _, sug := range suggestions {
    fmt.Printf("意图 '%s' 建议路由到 %s (置信度: %.2f)\n",
        sug.Intent, sug.SuggestedAgent, sug.Confidence)
}
```

## 配置示例

### 意图识别服务配置

```go
service := NewIntentRecognitionService(
    intentRepo,      // 意图仓储
    sampleRepo,      // 样本仓储
    llmClient,       // LLM客户端（OpenAI/Claude）
    vectorStore,     // 向量存储（Milvus）
    logger,          // 日志记录器
)

// 设置识别策略
service.strategy = "hybrid"    // llm, vector, hybrid
service.hybridWeight = 0.6      // LLM权重60%
```

### 评分引擎配置

```go
engine := NewScoringEngine(
    loadMonitor,      // 负载监控
    serviceRegistry,  // 服务注册中心
    optLogRepo,       // 优化日志仓储
    logger,           // 日志记录器
)

// 设置评分权重
engine.SetWeights(
    0.25,  // load_weight
    0.20,  // health_weight
    0.20,  // response_weight
    0.15,  // error_weight
    0.10,  // cost_weight
    0.10,  // preference_weight
)
```

### 路由优化器配置

```go
optimizer := NewRoutingOptimizer(
    routingEngine,    // 路由引擎
    scoringEngine,    // 评分引擎
    ruleRepo,         // 规则仓储
    optLogRepo,       // 优化日志仓储
    logger,           // 日志记录器
)

// 启用特性
optimizer.enableLoadPrediction = true
optimizer.enableCostOptimization = true
optimizer.enableUserPreference = true
```

## 最佳实践

### 1. 意图设计

- ✅ 意图名称应简洁明了（如 `query_weather`, `book_hotel`）
- ✅ 每个意图至少提供10-20个训练样本
- ✅ 样本应覆盖各种表达方式
- ✅ 定期更新训练样本以提高准确率

### 2. 路由策略

- ✅ 使用混合策略（LLM + Vector）提高准确率
- ✅ 根据业务场景调整评分权重
- ✅ 启用负载预测避免过载
- ✅ 考虑用户偏好提升体验

### 3. A/B测试

- ✅ 样本量至少1000（统计显著性）
- ✅ 流量分配建议50:50
- ✅ 测试期间不要频繁更改配置
- ✅ 等待统计显著后再选择获胜者

### 4. 路由学习

- ✅ 定期分析路由日志（每周/每月）
- ✅ 关注低性能Agent并优化
- ✅ 根据建议自动生成路由规则
- ✅ 逐步调整评分权重

### 5. 性能优化

- ✅ 使用Redis缓存热点意图
- ✅ 批量处理路由日志
- ✅ 异步记录日志不阻塞请求
- ✅ 向量搜索使用索引加速

## 常见问题

### Q1: 意图识别准确率低怎么办？

**A**:
1. 增加训练样本数量和多样性
2. 使用混合策略（Hybrid）提高准确率
3. 调整置信度阈值
4. 定期重新训练模型

### Q2: 路由决策太慢怎么办？

**A**:
1. 启用Redis缓存热点意图
2. 使用向量搜索代替LLM
3. 批量预处理Agent健康状态
4. 优化数据库查询

### Q3: A/B测试多久能出结果？

**A**:
- 最小样本量：1000（500+500）
- 如果日请求量1000，需要约1天
- 建议收集至少1周数据以获得更稳定的结果

### Q4: 如何处理识别失败的情况？

**A**:
1. 返回置信度最高的结果（即使低于阈值）
2. 使用兜底路由规则
3. 记录失败日志用于后续优化
4. 提示用户重新描述需求

### Q5: 路由学习建议多久执行一次？

**A**:
- 轻负载：每周执行一次
- 重负载：每天执行一次
- 大规模调整后：立即执行
- 定期回顾：每月一次深度分析

## 下一步

1. **实现Repository层**: 使用GORM实现数据库操作
2. **实现API Handler**: 创建HTTP接口
3. **前端集成**: 开发管理界面
4. **监控告警**: 集成Prometheus + Grafana
5. **性能测试**: 使用JMeter/K6进行压力测试

## 参考文档

- [实现总结](./IMPLEMENTATION_SUMMARY.md)
- [API文档](./API_DOCUMENTATION.md)
- [企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
