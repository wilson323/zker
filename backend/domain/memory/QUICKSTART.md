# 记忆引擎 - 快速开始指南

## 🚀 快速开始

### 1. 数据库初始化

执行SQL脚本创建表结构：

```bash
mysql -u root -p coze_studio < backend/domain/memory/internal/dal/schema/memory_tables.sql
```

### 2. 配置Milvus向量数据库

确保Milvus服务已启动：

```bash
docker-compose up -d milvus
```

### 3. 启动应用

```bash
cd backend
make server
```

### 4. 测试API

#### 存储对话记忆
```bash
curl -X POST http://localhost:8888/api/v1/memory/conversation \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "tenant-001",
    "user_id": "user-001",
    "conversation_id": "conv-001",
    "memory_type": "ENTITY",
    "content": "用户是HR经理"
  }'
```

#### 检索对话记忆
```bash
curl -X POST http://localhost:8888/api/v1/memory/conversation/search \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-001",
    "query": "用户的工作",
    "top_k": 5
  }'
```

## 📚 代码示例

### 基础使用

```go
package main

import (
    "context"
    "fmt"

    convEntity "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
    convService "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/service"
    "github.com/coze-dev/coze-studio/backend/infra/llm"
    vectordb "github.com/coze-dev/coze-studio/backend/infra/vectordb"
    "gorm.io/gorm"
)

func main() {
    // 1. 初始化依赖
    db := initDB() // 初始化MySQL连接
    vectorClient := initVectorClient() // 初始化Milvus客户端
    llmClient := llm.NewMockMemoryLLMClient() // 使用Mock LLM

    // 2. 创建服务
    memoryService := convService.NewConversationMemoryService(db, vectorClient, llmClient)

    // 3. 存储记忆
    ctx := context.Background()
    memory := &convEntity.ConversationMemory{
        TenantID:       "tenant-001",
        UserID:         "user-001",
        ConversationID: "conv-001",
        MemoryType:     convEntity.MemoryTypeEntity,
        Content:        "用户是HR经理，负责技术岗位招聘",
        ImportanceScore: 0.8,
    }

    err := memoryService.StoreMemory(ctx, memory)
    if err != nil {
        panic(err)
    }

    // 4. 检索记忆
    memories, err := memoryService.RetrieveMemories(ctx, "user-001", "用户工作", 5)
    if err != nil {
        panic(err)
    }

    for _, mem := range memories {
        fmt.Printf("记忆: %s (相似度: %.2f)\n", mem.Memory.Content, mem.Score)
    }
}
```

### 记忆增强对话

```go
package main

import (
    "context"

    "github.com/coze-dev/coze-studio/backend/application/conversation"
)

func chatWithMemoryExample() {
    // 初始化服务
    enhancedService := conversation.NewMemoryEnhancedConversationService(
        convSvc,
        knowSvc,
        llmSvc,
    )

    // 记忆增强对话
    ctx := context.Background()
    response, err := enhancedService.ChatWithMemory(ctx, &conversation.MemoryEnhancedChatRequest{
        TenantID:       "tenant-001",
        UserID:         "user-001",
        ConversationID: "conv-001",
        Query:          "如何筛选候选人?",
    })

    if err != nil {
        panic(err)
    }

    fmt.Printf("AI回复: %s\n", response.Response)
    fmt.Printf("使用了 %d 条记忆, %d 条知识\n", response.UsedMemories, response.UsedKnowledge)
}
```

## 🔧 配置说明

### 环境变量

```bash
# MySQL配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_DATABASE=coze_studio

# Milvus配置
MILVUS_HOST=localhost
MILVUS_PORT=19530
MILVUS_USERNAME=
MILVUS_PASSWORD=

# LLM配置
LLM_PROVIDER=openai
LLM_API_KEY=sk-xxx
LLM_EMBEDDING_MODEL=text-embedding-ada-002
```

### 依赖服务

```yaml
services:
  mysql:
    image: mysql:8.4.5
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: coze_studio

  milvus:
    image: milvusdb/milvus:v2.5.10
    ports:
      - "19530:19530"
    depends_on:
      - etcd
      - minio

  redis:
    image: redis:8.0
    ports:
      - "6379:6379"
```

## 🧪 运行测试

```bash
# 单元测试
cd backend/domain/memory/conversation/service
go test -v

# 集成测试(需要testcontainers)
go test -v -tags=integration

# 性能测试
go test -bench=. -benchmem
```

## 📖 常见问题

### Q: 如何切换向量数据库？
A: 实现`VectorClient`接口，支持Milvus、Weaviate、Qdrant等。

### Q: 如何使用真实的LLM？
A: 替换`MockMemoryLLMClient`为真实的LLM客户端实现。

### Q: 如何实现记忆缓存？
A: 在DAL层添加Redis缓存，提高检索性能。

### Q: 如何扩展记忆类型？
A: 在`MemoryType`枚举中添加新类型，并更新相关逻辑。

## 🔗 相关文档

- [完整设计文档](../docs/企业级功能完善与统一性设计方案/详细设计/17-五大AI引擎核心_记忆引擎.md)
- [API接口文档](../docs/企业级功能完善与统一性设计方案/API接口文档_记忆引擎.md)
- [企业级开发规范](../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

## 💡 最佳实践

1. **合理设置重要性评分**: 根据记忆类型动态调整重要性
2. **定期清理过期记忆**: 使用定时任务清理过期数据
3. **监控向量检索性能**: 确保检索响应时间 < 200ms
4. **使用缓存优化**: 热点记忆使用Redis缓存
5. **多租户隔离**: 确保所有查询都包含tenant_id过滤

---

**最后更新**: 2025-01-03
