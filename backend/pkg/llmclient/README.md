# llmclient 包

## 概述

`llmclient` 包为知识图谱服务提供统一的 LLM 客户端接口，支持从文本中抽取实体和关系。

## 接口定义

### LLMClient

```go
type LLMClient interface {
    // Chat 发送聊天请求
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

    // ChatStream 发送流式聊天请求
    ChatStream(ctx context.Context, req *ChatRequest) (<-chan *ChatChunk, error)

    // Close 关闭客户端
    Close() error
}
```

### 数据结构

#### ChatRequest

```go
type ChatRequest struct {
    Model       string     // 模型名称
    Messages    []*Message // 消息列表
    Temperature float64    // 温度参数（可选）
    MaxTokens   int        // 最大token数（可选）
}
```

#### Message

```go
type Message struct {
    Role    string // 角色：system/user/assistant
    Content string // 内容
}
```

#### ChatResponse

```go
type ChatResponse struct {
    ID      string   // 响应ID
    Object  string   // 对象类型
    Created int64    // 创建时间戳
    Model   string   // 模型名称
    Content string   // 响应内容（便捷字段）
    Choices []Choice // 选择列表
    Usage   Usage    // 使用情况
}
```

#### Choice

```go
type Choice struct {
    Index        int     // 选择索引
    Message      Message // 消息内容
    FinishReason string  // 结束原因
}
```

#### Usage

```go
type Usage struct {
    PromptTokens     int // 提示词token数
    CompletionTokens int // 完成token数
    TotalTokens      int // 总token数
}
```

### 配置

```go
type Config struct {
    APIKey  string // API密钥
    BaseURL string // API基础URL
    Model   string // 模型名称
}
```

## 使用示例

### 创建客户端

```go
import "github.com/coze-dev/coze-studio/backend/pkg/llmclient"

config := &llmclient.Config{
    APIKey:  "your-api-key",
    BaseURL: "https://api.example.com",
    Model:   "gpt-3.5-turbo",
}

client, err := llmclient.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 发送聊天请求

```go
req := &llmclient.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []*llmclient.Message{
        {Role: "system", Content: "You are a helpful assistant"},
        {Role: "user", Content: "Hello, how are you?"},
    },
    Temperature: 0.7,
    MaxTokens:   100,
}

resp, err := client.Chat(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Response:", resp.Content)
```

### 流式聊天

```go
req := &llmclient.ChatRequest{
    Model: "gpt-3.5-turbo",
    Messages: []*llmclient.Message{
        {Role: "user", Content: "Tell me a story"},
    },
}

chunkChan, err := client.ChatStream(context.Background(), req)
if err != nil {
    log.Fatal(err)
}

for chunk := range chunkChan {
    if len(chunk.Choices) > 0 {
        fmt.Print(chunk.Choices[0].Delta.Content)
    }
}
```

## 在知识图谱服务中的使用

```go
// 从文本中抽取实体
func (s *GraphConstructionService) ExtractEntities(ctx context.Context, tenantID, text string) ([]*entity.GraphEntity, error) {
    req := &llmclient.ChatRequest{
        Messages: []*llmclient.Message{
            {Role: "system", Content: "你是一个专业的知识图谱构建助手"},
            {Role: "user", Content: text},
        },
        Temperature: 0.3,
    }

    resp, err := s.llmClient.Chat(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("LLM调用失败: %w", err)
    }

    // 解析响应...
    return entities, nil
}
```

## TODO

当前实现为 `mockClient`，用于编译通过。实际使用时需要实现具体的 LLM 客户端：

1. **OpenAI 客户端**：集成 OpenAI API
2. **Claude 客户端**：集成 Anthropic Claude API
3. **本地模型客户端**：集成本地部署的模型（如 Ollama）
4. **自定义客户端**：根据业务需求定制

## 测试

```bash
cd backend/pkg/llmclient
go test -v
```

## 注意事项

1. 客户端实现需要是并发安全的
2. 应该实现连接池和超时控制
3. 需要实现重试机制和错误处理
4. 建议添加 Prometheus 监控指标
5. 敏感信息（API Key）应该通过配置中心或环境变量传入
