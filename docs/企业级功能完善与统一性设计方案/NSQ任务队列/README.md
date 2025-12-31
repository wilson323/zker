# ZKER NSQ任务队列系统

**版本**: v1.0
**创建日期**: 2025-01-01
**状态**: ✅ 已完成

---

## 📖 项目简介

ZKER NSQ任务队列系统是面向ZKER平台的企业级异步任务处理解决方案，提供高性能、高可靠性、可观测的异步任务处理能力。

### 核心特性

- ✅ **高可靠性**: 消息持久化、失败重试、死信队列
- ⚡ **高性能**: 批量处理、连接池、并发控制
- 📊 **可观测性**: Prometheus指标、结构化日志、Grafana大盘
- 🔐 **安全性**: TLS加密、认证授权、租户隔离
- 🔧 **易用性**: 统一接口、类型安全、易于集成

---

## 🚀 快速开始

### 安装NSQ

```bash
# macOS
brew install nsq

# 启动NSQ
nsqlookupd &
nsqd --lookupd-tcp-address=localhost:4160 &
```

### 基础使用

```go
import "github.com/coze-dev/coze-studio/backend/infra/queue"

// 1. 创建队列管理器
manager, _ := queue.NewQueueManager(config, logger)

// 2. 注册任务处理器
handler := queue.TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
    // 处理任务
    return nil
})
manager.RegisterTaskWithConsumer("my_task", handler)

// 3. 启动
manager.Start()

// 4. 发布任务
manager.EnqueueTask(ctx, "my_task", map[string]interface{}{"id": "123"})
```

---

## 📂 项目结构

```
backend/
├── infra/queue/           # 基础设施层
│   ├── nsq_config.go      # 配置模块
│   ├── nsq_producer.go    # 生产者
│   ├── nsq_consumer.go    # 消费者
│   ├── dead_letter_queue.go # 死信队列
│   ├── manager.go         # 队列管理器
│   └── nsq_test.go        # 测试文件
│
└── application/queue/     # 应用层
    └── tasks.go          # 业务任务封装

backend/conf/queue/
└── nsq.yaml             # 配置示例

docs/企业级功能完善与统一性设计方案/NSQ任务队列/
├── README.md            # 本文件
├── ZKER-NSQ任务队列架构设计文档_v1.0.md
├── ZKER-NSQ任务队列使用指南_v1.0.md
├── ZKER-NSQ任务队列运维手册_v1.0.md
└── ZKER-NSQ任务队列实现报告_v1.0.md
```

---

## 📚 文档索引

### 核心文档

1. **[架构设计文档](./ZKER-NSQ任务队列架构设计文档_v1.0.md)**
   - 系统架构设计
   - 核心组件说明
   - 数据流设计
   - 可靠性保障
   - 性能优化策略

2. **[使用指南](./ZKER-NSQ任务队列使用指南_v1.0.md)**
   - 快速开始教程
   - 基础使用方法
   - 高级功能说明
   - 业务任务封装
   - 最佳实践
   - 常见问题

3. **[运维手册](./ZKER-NSQ任务队列运维手册_v1.0.md)**
   - 部署指南（Docker/K8s）
   - 配置管理
   - 监控告警
   - 故障处理
   - 性能调优
   - 备份恢复

4. **[实现报告](./ZKER-NSQ任务队列实现报告_v1.0.md)**
   - 项目概述
   - 交付清单
   - 核心功能
   - 测试覆盖
   - 验收标准

---

## 🎯 核心功能

### 1. 任务发布

```go
// 简单任务
manager.EnqueueTask(ctx, "task_type", task)

// 延迟任务
manager.EnqueueTaskDelayed(ctx, "task_type", 5*time.Minute, task)

// 批量任务
manager.EnqueueTaskBatch(ctx, "task_type", tasks)
```

### 2. 任务处理

```go
// 函数式处理器
handler := queue.TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
    // 处理逻辑
    return nil
})

// 类型化处理器
handler := queue.NewTypedMessageHandler(func(ctx context.Context, msg *MyTask) error {
    // 类型安全的处理
    return nil
})
```

### 3. 死信队列

```go
// 自动转入DLQ
// 当重试次数超过限制时自动转入

// 手动重放
manager.ReplayDLQTasks(ctx, "task_type")

// 重放单条消息
manager.ReplayDLQTask(ctx, "message_id")
```

### 4. 监控指标

```go
stats := manager.GetStats()

// 生产者指标
fmt.Printf("Published: %d\n", stats.ProducerMetrics.PublishCount)

// 消费者指标
fmt.Printf("Success: %d\n", stats.ConsumerMetrics["task_type"].MessagesSuccess)

// 死信队列指标
fmt.Printf("DLQ: %d\n", stats.DLQMetrics.MessagesStored)
```

---

## 🔧 配置示例

```yaml
# backend/conf/queue/nsq.yaml

nsqd_addresses:
  - "localhost:4150"

nsqlookupd_addresses:
  - "localhost:4161"

producer:
  timeout: 5s
  retry: 3
  retry_interval: 1s

consumer:
  max_retries: 3
  retry_delay: 1s
  max_in_flight: 10
  message_timeout: 60s

dlq:
  enabled: true
  topic_suffix: "-dlq"
  retention_duration: 72h
  max_replays: 3

metrics:
  enabled: true
  address: ":2112"

logging:
  level: "info"
  verbose: false
```

---

## 📦 预定义任务类型

系统提供6种预定义业务任务：

1. **KnowledgeDocumentProcessTask** - 知识库文档处理
2. **WorkflowExecutionTask** - 工作流执行
3. **BotPublishTask** - Bot发布
4. **ReportGenerationTask** - 报告生成
5. **TokenUsageTask** - Token计量
6. **EmailNotificationTask** - 邮件通知

详细使用方法请参考[使用指南](./ZKER-NSQ任务队列使用指南_v1.0.md)。

---

## 🧪 测试

### 运行测试

```bash
# 单元测试
cd backend/infra/queue
go test -v -cover

# 并发测试
go test -v -race

# 性能测试
go test -bench=. -benchmem
```

### 测试覆盖

- ✅ 配置验证测试
- ✅ 重试策略测试
- ✅ 存储CRUD测试
- ✅ 并发安全测试
- ✅ 性能Benchmark测试

---

## 📊 性能指标

| 指标 | 目标值 | 备注 |
|-----|--------|------|
| 吞吐量 | >10000 msg/s | 单个NSQD节点 |
| 延迟 | P99 < 100ms | 端到端延迟 |
| 可用性 | 99.9% | 集群部署 |
| 并发度 | >1000 | 消费者并发 |

---

## 🔐 安全特性

- **TLS加密**: 传输层数据加密
- **认证授权**: 基于HTTP的认证
- **租户隔离**: Topic级别的租户隔离
- **审计日志**: 完整的操作审计

---

## 🎨 架构亮点

### SOLID设计
- **单一职责**: 每个组件职责明确
- **接口隔离**: 清晰的接口定义
- **依赖倒置**: 依赖抽象而非具体实现

### 并发安全
- **互斥锁**: 保护共享状态
- **原子操作**: 计数器更新
- **并发控制**: 限制goroutine数量

### 优雅关闭
- **等待处理完成**: 最多30秒
- **资源释放**: 连接、文件、内存
- **状态保持**: 保证数据一致性

---

## 📈 后续优化

- [ ] 支持优先级队列
- [ ] 支持任务调度（Cron）
- [ ] 增强消息追踪（Trace ID传递）
- [ ] 支持消息压缩
- [ ] 集成testcontainers进行集成测试

---

## 🤝 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交变更
4. 发起Pull Request

### 代码规范

- 遵循Go最佳实践
- 添加单元测试
- 更新相关文档
- 通过所有CI检查

---

## 📞 支持与联系

- **技术支持**: 基础设施团队
- **文档维护**: 技术架构委员会
- **Bug反馈**: GitHub Issues
- **功能建议**: 产品需求池

---

## 📄 许可证

Copyright © 2025 ZKER Project. All rights reserved.

---

**最后更新**: 2025-01-01
