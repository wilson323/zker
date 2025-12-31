# ZKER NSQ任务队列使用指南

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**目标读者**: 开发者

---

## 📋 目录

- [1. 快速开始](#1-快速开始)
- [2. 基础使用](#2-基础使用)
- [3. 高级功能](#3-高级功能)
- [4. 业务任务封装](#4-业务任务封装)
- [5. 最佳实践](#5-最佳实践)
- [6. 故障排查](#6-故障排查)
- [7. 常见问题](#7-常见问题)

---

## 1. 快速开始

### 1.1 环境准备

**安装NSQ**:
```bash
# macOS
brew install nsq

# Linux
wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-1.2.1.linux-amd64.go1.16.6.tar.gz
tar -xzf nsq-1.2.1.linux-amd64.go1.16.6.tar.gz
```

**启动NSQ**:
```bash
# 启动NSQLookupd
nsqlookupd

# 启动NSQD
nsqd --lookupd-tcp-address=localhost:4160

# 启动NSQAdmin（可选）
nsqadmin --lookupd-http-address=localhost:4161
```

### 1.2 最小化示例

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/coze-dev/coze-studio/backend/infra/queue"
    "go.uber.org/zap"
)

func main() {
    // 1. 创建配置
    config := queue.DefaultNSQConfig()

    // 2. 创建日志
    logger, _ := zap.NewProduction()

    // 3. 创建队列管理器
    manager, err := queue.NewQueueManager(config, logger)
    if err != nil {
        log.Fatal(err)
    }

    // 4. 注册任务处理器
    handler := queue.TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
        logger.Info("processing task", zap.String("type", taskType))
        // 处理任务逻辑
        return nil
    })

    manager.RegisterTaskWithConsumer("test_task", handler)

    // 5. 启动队列管理器
    if err := manager.Start(); err != nil {
        log.Fatal(err)
    }

    // 6. 发布任务
    ctx := context.Background()
    task := map[string]interface{}{
        "id":   "123",
        "name": "test",
    }

    if err := manager.EnqueueTask(ctx, "test_task", task); err != nil {
        log.Fatal(err)
    }

    // 7. 保持运行
    time.Sleep(10 * time.Second)

    // 8. 优雅关闭
    manager.Stop()
}
```

---

## 2. 基础使用

### 2.1 配置队列管理器

**从配置文件加载**:
```go
import "gopkg.in/yaml.v2"

// 加载配置
data, err := os.ReadFile("conf/queue/nsq.yaml")
if err != nil {
    log.Fatal(err)
}

var config queue.NSQConfig
if err := yaml.Unmarshal(data, &config); err != nil {
    log.Fatal(err)
}

// 创建队列管理器
manager, err := queue.NewQueueManager(&config, logger)
```

**代码配置**:
```go
config := &queue.NSQConfig{
    NSQDAddresses:       []string{"localhost:4150"},
    NSQLookupdAddresses: []string{"localhost:4161"},
    ProducerTimeout:     5 * time.Second,
    ProducerRetry:       3,
    ConsumerMaxRetries:  3,
    ConsumerMaxInFlight: 10,
    DLQEnabled:          true,
}

manager, err := queue.NewQueueManager(config, logger)
```

### 2.2 发布任务

**简单任务**:
```go
task := map[string]interface{}{
    "document_id": "doc123",
    "action":      "parse",
}

err := manager.EnqueueTask(ctx, "knowledge_document", task)
```

**延迟任务**:
```go
// 延迟5分钟执行
err := manager.EnqueueTaskDelayed(ctx, "knowledge_document", 5*time.Minute, task)
```

**批量任务**:
```go
tasks := []interface{}{
    map[string]interface{}{"document_id": "doc1", "action": "parse"},
    map[string]interface{}{"document_id": "doc2", "action": "parse"},
    map[string]interface{}{"document_id": "doc3", "action": "parse"},
}

err := manager.EnqueueTaskBatch(ctx, "knowledge_document", tasks)
```

### 2.3 处理任务

**函数式处理器**:
```go
handler := queue.TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
    var task struct {
        DocumentID string `json:"document_id"`
        Action     string `json:"action"`
    }

    if err := json.Unmarshal(taskData, &task); err != nil {
        return err
    }

    // 处理任务
    switch task.Action {
    case "parse":
        return parseDocument(ctx, task.DocumentID)
    case "index":
        return indexDocument(ctx, task.DocumentID)
    }

    return nil
})

manager.RegisterTaskWithConsumer("knowledge_document", handler)
```

**结构化处理器**:
```go
type KnowledgeDocumentHandler struct {
    knowledgeService KnowledgeService
    logger          *zap.Logger
}

func (h *KnowledgeDocumentHandler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
    var task KnowledgeDocumentProcessTask
    if err := json.Unmarshal(taskData, &task); err != nil {
        return err
    }

    return h.knowledgeService.ProcessDocument(ctx, &task)
}

handler := &KnowledgeDocumentHandler{
    knowledgeService: svc,
    logger:          logger,
}

manager.RegisterTaskWithConsumer("knowledge_document", handler)
```

### 2.4 启动和停止

**启动**:
```go
// 注册所有任务处理器后启动
if err := manager.Start(); err != nil {
    log.Fatal(err)
}

// 监听关闭信号
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
<-sigChan

// 优雅关闭
manager.Stop()
```

**优雅关闭**:
```go
func (m *QueueManager) Stop() error {
    m.logger.Info("stopping queue manager")

    // 1. 停止接收新消息
    // 2. 等待正在处理的消息完成（最多30秒）
    // 3. 停止所有消费者
    // 4. 停止生产者

    m.logger.Info("queue manager stopped")
    return nil
}
```

---

## 3. 高级功能

### 3.1 自定义重试策略

**指数退避**:
```go
strategy := &queue.ExponentialBackoffStrategy{
    InitialDelay: 1 * time.Second,
    MaxDelay:     60 * time.Second,
    MaxAttempts:  5,
}

err := producer.PublishWithRetry(ctx, "topic", message, strategy)
```

**固定延迟**:
```go
strategy := &queue.FixedDelayStrategy{
    Delay:       5 * time.Second,
    MaxAttempts: 3,
}

err := producer.PublishWithRetry(ctx, "topic", message, strategy)
```

**自定义策略**:
```go
type CustomRetryStrategy struct{}

func (s *CustomRetryStrategy) ShouldRetry(ctx context.Context, attempt int, err error) (bool, time.Duration) {
    // 只重试网络错误
    if isNetworkError(err) && attempt < 3 {
        return true, 2 * time.Second
    }
    return false, 0
}
```

### 3.2 死信队列操作

**重放死信队列**:
```go
// 重放特定Topic的所有死信消息
err := manager.ReplayDLQTasks(ctx, "knowledge_document")
```

**重放单条消息**:
```go
// 获取DLQ消息ID
dlq := manager.GetDLQ()
messages, _ := dlq.storage.List(ctx, "knowledge_document-dlq", 100)

// 重放指定消息
for _, msg := range messages {
    if shouldReplay(msg) {
        err := manager.ReplayDLQTask(ctx, msg.MessageID)
        if err != nil {
            logger.Error("failed to replay message",
                zap.String("message_id", msg.MessageID),
                zap.Error(err),
            )
        }
    }
}
```

**删除死信消息**:
```go
dlq := manager.GetDLQ()
err := dlq.DeleteMessage(ctx, "message_id")
```

### 3.3 监控和指标

**获取统计信息**:
```go
stats := manager.GetStats()

// 生产者指标
fmt.Printf("Published: %d\n", stats.ProducerMetrics.PublishCount)
fmt.Printf("Success: %d\n", stats.ProducerMetrics.PublishSuccess)
fmt.Printf("Failed: %d\n", stats.ProducerMetrics.PublishFailed)

// 消费者指标
for taskType, consumerMetrics := range stats.ConsumerMetrics {
    fmt.Printf("Task: %s\n", taskType)
    fmt.Printf("  Received: %d\n", consumerMetrics.MessagesReceived)
    fmt.Printf("  Success: %d\n", consumerMetrics.MessagesSuccess)
    fmt.Printf("  Failed: %d\n", consumerMetrics.MessagesFailed)
    fmt.Printf("  To DLQ: %d\n", consumerMetrics.MessagesToDLQ)
}

// 死信队列指标
fmt.Printf("DLQ Stored: %d\n", stats.DLQMetrics.MessagesStored)
fmt.Printf("DLQ Replayed: %d\n", stats.DLQMetrics.MessagesReplayed)
```

**Prometheus集成**:
```go
import "github.com/prometheus/client_golang/prometheus"

// 暴露指标端点
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(":2112", nil)
```

### 3.4 动态消费者管理

**运行时添加消费者**:
```go
handler := queue.TaskHandlerFunc(func(ctx context.Context, taskType string, taskData []byte) error {
    // 处理逻辑
    return nil
})

err := manager.AddDynamicConsumer("new_task_type", handler)
```

**移除消费者**:
```go
err := manager.RemoveConsumer("old_task_type")
```

**获取所有消费者**:
```go
consumers := manager.GetAllConsumers()
for taskType, consumer := range consumers {
    if consumer.IsRunning() {
        fmt.Printf("Consumer %s is running\n", taskType)
    }
}
```

---

## 4. 业务任务封装

### 4.1 知识库文档处理

**任务定义**:
```go
type KnowledgeDocumentProcessTask struct {
    DocumentID string `json:"document_id"`
    KnowledgeID string `json:"knowledge_id"`
    Action     string `json:"action"` // "parse", "index", "delete"
    Priority   int    `json:"priority"`
    TenantID   string `json:"tenant_id"`
    UserID     string `json:"user_id"`
}
```

**使用示例**:
```go
// 发布任务
task := &KnowledgeDocumentProcessTask{
    DocumentID: "doc123",
    KnowledgeID: "kb456",
    Action:     "parse",
    Priority:   1,
    TenantID:   "tenant_001",
    UserID:     "user_001",
}

err := manager.EnqueueTask(ctx, TaskTypeKnowledgeDocument, task)
```

### 4.2 工作流执行

**任务定义**:
```go
type WorkflowExecutionTask struct {
    WorkflowID  string                 `json:"workflow_id"`
    VersionID   string                 `json:"version_id,omitempty"`
    Input       map[string]interface{} `json:"input"`
    TriggeredBy  string                 `json:"triggered_by"`
    TenantID    string                 `json:"tenant_id"`
    UserID      string                 `json:"user_id"`
    ScheduledAt string                 `json:"scheduled_at,omitempty"`
}
```

**使用示例**:
```go
// 立即执行
task := &WorkflowExecutionTask{
    WorkflowID: "wf123",
    Input: map[string]interface{}{
        "query": "分析销售数据",
        "date_range": map[string]string{
            "start": "2025-01-01",
            "end":   "2025-01-31",
        },
    },
    TriggeredBy: "user",
    TenantID:    "tenant_001",
    UserID:      "user_001",
}

err := manager.EnqueueTask(ctx, TaskTypeWorkflowExecution, task)

// 延迟执行（定时任务）
scheduledTime := time.Now().Add(24 * time.Hour)
task.ScheduledAt = scheduledTime.Format(time.RFC3339)

err := manager.EnqueueTaskDelayed(ctx, TaskTypeWorkflowExecution, 24*time.Hour, task)
```

### 4.3 Bot发布

**任务定义**:
```go
type BotPublishTask struct {
    BotID      string   `json:"bot_id"`
    VersionID  string   `json:"version_id"`
    PublishTo  []string `json:"publish_to"` // ["weixin", "feishu", "api"]
    Priority   int      `json:"priority"`
    TenantID   string   `json:"tenant_id"`
    UserID     string   `json:"user_id"`
    AutoEnable bool     `json:"auto_enable"`
}
```

**使用示例**:
```go
task := &BotPublishTask{
    BotID:     "bot123",
    VersionID: "v456",
    PublishTo: []string{"weixin", "feishu"},
    Priority:  1,
    TenantID:  "tenant_001",
    UserID:    "user_001",
    AutoEnable: true,
}

err := manager.EnqueueTask(ctx, TaskTypeBotPublish, task)
```

### 4.4 报告生成

**任务定义**:
```go
type ReportGenerationTask struct {
    ReportID   string `json:"report_id"`
    ReportType string `json:"report_type"` // "usage", "performance", "error"
    StartDate  string `json:"start_date"`
    EndDate    string `json:"end_date"`
    TenantID   string `json:"tenant_id"`
    UserID     string `json:"user_id"`
    Format     string `json:"format"` // "pdf", "xlsx", "csv"
    EmailTo    string `json:"email_to,omitempty"`
}
```

**使用示例**:
```go
task := &ReportGenerationTask{
    ReportID:   "report123",
    ReportType: "usage",
    StartDate:  "2025-01-01",
    EndDate:    "2025-01-31",
    TenantID:   "tenant_001",
    UserID:     "user_001",
    Format:     "pdf",
    EmailTo:    "admin@example.com",
}

err := manager.EnqueueTask(ctx, TaskTypeReportGeneration, task)
```

---

## 5. 最佳实践

### 5.1 任务设计

**幂等性**:
```go
func (h *Handler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
    var task Task
    json.Unmarshal(taskData, &task)

    // 检查是否已处理
    if h.isProcessed(task.ID) {
        return nil // 跳过重复处理
    }

    // 处理任务
    err := h.processTask(ctx, &task)
    if err != nil {
        return err
    }

    // 标记已处理
    return h.markProcessed(task.ID)
}
```

**任务ID生成**:
```go
import "github.com/google/uuid"

taskID := uuid.New().String()

// 或使用业务相关的ID
taskID := fmt.Sprintf("doc_%s_%d", documentID, time.Now().Unix())
```

**任务优先级**:
```go
// 高优先级任务（紧急）
task.Priority = 1

// 普通任务（默认）
task.Priority = 5

// 低优先级任务（批量）
task.Priority = 10
```

### 5.2 错误处理

**区分临时错误和永久错误**:
```go
func (h *Handler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
    err := h.processTask(ctx, task)

    // 临时错误（可重试）
    if isTemporaryError(err) {
        return err // NSQ会重试
    }

    // 永久错误（不重试，直接记录）
    if isPermanentError(err) {
        h.logPermanentError(ctx, task, err)
        return nil // 返回nil避免重试
    }

    return err
}
```

**超时处理**:
```go
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

err := h.processTaskWithContext(ctx, task)
if errors.Is(err, context.DeadlineExceeded) {
    // 超时错误，可重试
    return err
}
```

### 5.3 性能优化

**批量处理**:
```go
func (h *Handler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
    var tasks []Task
    json.Unmarshal(taskData, &tasks)

    // 批量处理
    for i := 0; i < len(tasks); i += batchSize {
        end := min(i+batchSize, len(tasks))
        batch := tasks[i:end]

        if err := h.processBatch(ctx, batch); err != nil {
            return err
        }
    }

    return nil
}
```

**并行处理**:
```go
func (h *Handler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
    var tasks []Task
    json.Unmarshal(taskData, &tasks)

    var wg sync.WaitGroup
    sem := make(chan struct{}, 10) // 限制并发数

    for _, task := range tasks {
        wg.Add(1)
        sem <- struct{}{}

        go func(t Task) {
            defer wg.Done()
            defer func() { <-sem }()

            h.processTask(ctx, &t)
        }(task)
    }

    wg.Wait()
    return nil
}
```

### 5.4 监控和日志

**结构化日志**:
```go
h.logger.Info("processing task",
    zap.String("task_type", taskType),
    zap.String("task_id", task.ID),
    zap.Duration("processing_time", time.Since(start)),
)
```

**指标上报**:
```go
// 计数器
taskCounter.Inc()

// 延迟直方图
taskDuration.Observe(time.Since(start).Seconds())

// 成功/失败
if err != nil {
    taskFailedCounter.Inc()
} else {
    taskSuccessCounter.Inc()
}
```

---

## 6. 故障排查

### 6.1 消息堆积

**症状**:
- Queue深度持续增长
- 消费延迟增加

**排查步骤**:
```bash
# 查看队列深度
curl http://localhost:4161/stats | jq '.data.topics[].depth'

# 查看消费者状态
curl http://localhost:4161/stats | jq '.data.topics[].clients'
```

**解决方案**:
1. 增加消费者数量（水平扩展）
2. 优化处理逻辑（提高吞吐）
3. 调整MaxInFlight参数

### 6.2 死信队列增长

**症状**:
- DLQ消息数快速增加
- 大量任务失败

**排查步骤**:
```go
// 查看DLQ消息
dlq := manager.GetDLQ()
messages, _ := dlq.List(ctx, "topic-dlq", 100)

// 分析失败原因
for _, msg := range messages {
    fmt.Printf("Message: %s, Error: %s\n", msg.MessageID, msg.LastError)
}
```

**解决方案**:
1. 分析失败原因，修复代码
2. 调整重试策略
3. 优化超时时间

### 6.3 消费缓慢

**症状**:
- 处理时间长
- CPU/内存占用高

**排查步骤**:
```go
// 查看处理时间分布
stats := manager.GetStats()
metrics := stats.ConsumerMetrics["task_type"]

for _, latency := range metrics.ProcessingTime {
    fmt.Printf("Latency: %v\n", latency)
}
```

**解决方案**:
1. 使用pprof分析性能瓶颈
2. 优化数据库查询
3. 添加缓存
4. 批量处理

---

## 7. 常见问题

### Q1: 如何确保消息不丢失？

**A**: NSQ默认持久化消息到磁盘。额外措施：
- 设置`--mem-queue-size=0`禁用内存队列
- 使用多个NSQD节点（集群模式）
- 启用死信队列

### Q2: 如何处理重复消息？

**A**: 实现幂等性：
```go
func (h *Handler) HandleTask(ctx context.Context, taskType string, taskData []byte) error {
    var task Task
    json.Unmarshal(taskData, &task)

    // 使用Redis去重
    key := fmt.Sprintf("task:%s", task.ID)
    if h.redis.Exists(ctx, key) {
        return nil
    }

    // 处理任务
    err := h.process(ctx, &task)
    if err != nil {
        return err
    }

    // 标记已处理（24小时过期）
    return h.redis.Set(ctx, key, "1", 24*time.Hour)
}
```

### Q3: 如何实现任务调度？

**A**: 使用延迟消息：
```go
// 每天凌晨2点执行
now := time.Now()
next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
if next.Before(now) {
    next = next.Add(24 * time.Hour)
}

delay := next.Sub(now)
manager.EnqueueTaskDelayed(ctx, "daily_report", delay, reportTask)
```

### Q4: 如何监控队列健康？

**A**: 设置告警规则：
```yaml
- alert: NSQQueueDepth
  expr: nsq_queue_depth > 10000
  for: 10m
  annotations:
    summary: "Queue depth too high"
```

### Q5: 如何优雅关闭消费者？

**A**:
```go
// 监听系统信号
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

<-sigChan

// 优雅关闭（等待最多30秒）
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

manager.Stop()
```

---

## 附录

### A. 配置示例

完整配置参考: [backend/conf/queue/nsq.yaml](../../backend/conf/queue/nsq.yaml)

### B. API参考

详细API文档: [ZKER-NSQ任务队列架构设计文档.md](./ZKER-NSQ任务队列架构设计文档_v1.0.md)

### C. 运维手册

运维操作手册: [ZKER-NSQ任务队列运维手册_v1.0.md](./ZKER-NSQ任务队列运维手册_v1.0.md)

---

**© 2025 ZKER Project. All rights reserved.**
