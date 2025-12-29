# 百应前台 - 会话 (Super Chatbox) 详细设计文档

> **模块编号**: 01
> **模块名称**: 会话 (Super Chatbox)
> **功能定位**: 员工使用前台的核心入口,多模态AI对话交互
> **对齐产品**: 鲸智百应 - 会话功能
> **版本**: v1.0

---

## 📋 模块概述

### 功能定义

**会话 (Super Chatbox)** 是企业员工与 AI 数字员工交互的核心入口,提供多模态输入、智能路由、多轮对话、上下文记忆等能力。

### 核心价值

- 🎯 **零学习成本** - 自然语言交互,无需学习
- 🎯 **多模态输入** - 支持文字/语音/图片/文件
- 🎯 **智能路由** - 自动识别意图并分配给合适的数字员工
- 🎯 **上下文记忆** - 四级记忆体系,越用越智能
- 🎯 **实时响应** - 流式输出,< 2s 响应延迟

---

## 一、功能架构设计

### 1.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Super Chatbox 架构                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  【前端交互层】                                               │
│  ├── 输入区域 (文字/语音/图片/文件)                          │
│  ├── 对话区域 (消息流/Markdown渲染/代码高亮)                │
│  └── 快捷操作 (引导问题/常用功能/历史记录)                  │
│                                                              │
│  【服务层】                                                   │
│  ├── 会话管理 (ChatService)                                 │
│  ├── 消息管理 (MessageService)                              │
│  ├── 上下文管理 (ContextService)                            │
│  └── 实时通信 (WebSocket/SSE)                               │
│                                                              │
│  【AI引擎层】                                                 │
│  ├── 智能路由引擎 (RoutingEngine)                           │
│  ├── 意图识别 (IntentRecognizer)                            │
│  ├── 数字员工调度 (BotDispatcher)                           │
│  └── 插件调度 (PluginScheduler)                             │
│                                                              │
│  【记忆引擎层】                                               │
│  ├── 会话记忆 (SessionMemory)                               │
│  ├── 个人记忆 (PersonalMemory)                             │
│  ├── 组织记忆 (OrgMemory)                                  │
│  └── 企业记忆 (EnterpriseMemory)                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 核心组件

```typescript
// 会话核心组件
interface SuperChatboxComponents {
  // 输入组件
  InputArea: {
    TextInput: TextInput          // 文字输入
    VoiceInput: VoiceInput        // 语音输入 (自动转文字)
    ImageInput: ImageInput        // 图片输入 (OCR识别)
    FileInput: FileInput          // 文件上传 (PDF/Word/Excel)
  }

  // 对话组件
  ChatArea: {
    MessageList: MessageList      // 消息列表
    MessageRenderer: MessageRenderer  // 消息渲染 (Markdown/代码)
    TypingIndicator: TypingIndicator  // 输入指示器
  }

  // 快捷操作
  QuickActions: {
    SuggestedQuestions: SuggestedQuestions  // 引导问题
    QuickCommands: QuickCommands    // 快捷命令
    History: History                // 历史记录
  }
}
```

---

## 二、前端设计 (React + TypeScript)

### 2.1 页面结构

```typescript
// src/pages/ChatboxPage.tsx
import React, { useState, useRef, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useChatbox } from '@/hooks/useChatbox';

interface ChatboxPageProps {
  tenantId: string;
  botId?: string;           // 可选,如果指定则直接进入该数字员工
}

const ChatboxPage: React.FC<ChatboxPageProps> = ({ tenantId, botId }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isRecording, setIsRecording] = useState<boolean>(false);

  // WebSocket 连接
  const { ws, send, disconnect } = useWebSocket(
    `/api/v1/tenants/${tenantId}/chat/stream`
  );

  // 消息处理
  const handleSend = async () => {
    if (!input.trim()) return;

    // 添加用户消息
    const userMessage: Message = {
      id: uuid(),
      role: 'user',
      content: input,
      type: 'text',
      timestamp: new Date(),
    };
    setMessages(prev => [...prev, userMessage]);

    // 发送到后端
    send({
      type: 'message',
      data: {
        content: input,
        bot_id: botId,
        conversation_id: conversationId,
      },
    });

    setInput('');
  };

  return (
    <div className="chatbox-page">
      {/* 头部 */}
      <ChatboxHeader tenantId={tenantId} botId={botId} />

      {/* 消息列表 */}
      <MessageList messages={messages} isLoading={isLoading} />

      {/* 引导问题 */}
      {messages.length === 0 && (
        <SuggestedQuestions tenantId={tenantId} onSelect={setInput} />
      )}

      {/* 输入区域 */}
      <InputArea
        value={input}
        onChange={setInput}
        onSend={handleSend}
        onVoiceStart={() => setIsRecording(true)}
        onVoiceStop={() => setIsRecording(false)}
        onFileUpload={handleFileUpload}
      />
    </div>
  );
};
```

### 2.2 输入区域组件

```typescript
// src/components/chatbox/InputArea.tsx
interface InputAreaProps {
  value: string;
  onChange: (value: string) => void;
  onSend: () => void;
  onVoiceStart: () => void;
  onVoiceStop: () => void;
  onFileUpload: (file: File) => void;
}

const InputArea: React.FC<InputAreaProps> = ({
  value,
  onChange,
  onSend,
  onVoiceStart,
  onVoiceStop,
  onFileUpload,
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null);

  return (
    <div className="input-area">
      {/* 工具栏 */}
      <div className="toolbar">
        {/* 文件上传 */}
        <Button
          icon={<IconUpload />}
          onClick={() => fileInputRef.current?.click()}
        >
          上传文件
        </Button>
        <input
          ref={fileInputRef}
          type="file"
          hidden
          accept=".pdf,.doc,.docx,.ppt,.pptx,.txt,.md,.jpg,.png"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) onFileUpload(file);
          }}
        />

        {/* 语音输入 */}
        <VoiceRecorder
          onStart={onVoiceStart}
          onStop={onVoiceStop}
        />

        {/* 快捷命令 */}
        <QuickCommandMenu />
      </div>

      {/* 文字输入框 */}
      <div className="textarea-container">
        <TextArea
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="输入消息,Shift+Enter 换行,Enter 发送..."
          autoSize={{ minRows: 1, maxRows: 10 }}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
              e.preventDefault();
              onSend();
            }
          }}
        />
      </div>

      {/* 发送按钮 */}
      <Button
        type="primary"
        icon={<IconSend />}
        onClick={onSend}
        disabled={!value.trim()}
      >
        发送
      </Button>
    </div>
  );
};
```

### 2.3 消息渲染组件

```typescript
// src/components/chatbox/MessageRenderer.tsx
interface MessageRendererProps {
  message: Message;
  tenantId: string;
}

const MessageRenderer: React.FC<MessageRendererProps> = ({ message, tenantId }) => {
  return (
    <div className={`message message-${message.role}`}>
      {/* 头像 */}
      <Avatar
        src={message.role === 'assistant' ? message.botAvatar : message.userAvatar}
        size="large"
      />

      {/* 消息内容 */}
      <div className="message-content">
        {/* 用户名/机器人名 */}
        <div className="message-sender">
          {message.role === 'assistant' ? message.botName : message.userName}
        </div>

        {/* 消息体 */}
        <div className="message-body">
          {/* 文字消息 (支持 Markdown) */}
          {message.type === 'text' && (
            <MarkdownRenderer content={message.content} />
          )}

          {/* 代码消息 */}
          {message.type === 'code' && (
            <CodeBlock
              language={message.language}
              code={message.content}
            />
          )}

          {/* 图片消息 */}
          {message.type === 'image' && (
            <ImageViewer src={message.content} />
          )}

          {/* 文件消息 */}
          {message.type === 'file' && (
            <FileAttachment
              fileName={message.fileName}
              fileUrl={message.fileUrl}
              fileSize={message.fileSize}
            />
          )}

          {/* 工具调用消息 */}
          {message.type === 'tool_call' && (
            <ToolInvocation
              toolName={message.toolName}
              toolInput={message.toolInput}
              toolOutput={message.toolOutput}
            />
          )}
        </div>

        {/* 时间戳 */}
        <div className="message-time">
          {formatTime(message.timestamp)}
        </div>

        {/* 操作按钮 */}
        {message.role === 'assistant' && (
          <div className="message-actions">
            <CopyButton content={message.content} />
            <RegenerateButton message={message} />
            <FeedbackButton message={message} />
          </div>
        )}
      </div>
    </div>
  );
};
```

### 2.4 Markdown 渲染器

```typescript
// src/components/markdown/MarkdownRenderer.tsx
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import remarkGfm from 'remark-gfm';

interface MarkdownRendererProps {
  content: string;
}

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ content }) => {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      components={{
        // 代码块渲染
        code({ node, inline, className, children, ...props }) {
          const match = /language-(\w+)/.exec(className || '');
          const language = match ? match[1] : '';

          return !inline && language ? (
            <SyntaxHighlighter
              language={language}
              PreTag="div"
              className="code-block"
            >
              {String(children).replace(/\n$/, '')}
            </SyntaxHighlighter>
          ) : (
            <code className={className} {...props}>
              {children}
            </code>
          );
        },

        // 表格渲染
        table({ children }) {
          return (
            <div className="table-container">
              <table>{children}</table>
            </div>
          );
        },

        // 链接渲染
        a({ children, href }) {
          return (
            <a
              href={href}
              target="_blank"
              rel="noopener noreferrer"
            >
              {children}
            </a>
          );
        },
      }}
    >
      {content}
    </ReactMarkdown>
  );
};
```

### 2.5 语音输入组件

```typescript
// src/components/chatbox/VoiceRecorder.tsx
import { useVoiceRecorder } from '@/hooks/useVoiceRecorder';

interface VoiceRecorderProps {
  onStart: () => void;
  onStop: (text: string) => void;
}

const VoiceRecorder: React.FC<VoiceRecorderProps> = ({ onStart, onStop }) => {
  const { isRecording, startRecording, stopRecording } = useVoiceRecorder({
    onTranscribe: (text) => {
      onStop(text);
    },
  });

  return (
    <Button
      type={isRecording ? 'danger' : 'default'}
      icon={<IconMicrophone />}
      onClick={() => {
        if (isRecording) {
          stopRecording();
        } else {
          startRecording();
          onStart();
        }
      }}
    >
      {isRecording ? '停止录音' : '语音输入'}
    </Button>
  );
};
```

---

## 三、后端设计 (Go)

### 3.1 API 接口设计

```go
// api/v1/chatbox.go
package v1

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// ChatboxRouter 聊天路由
type ChatboxRouter struct {
    chatService *service.ChatService
}

// NewChatboxRouter 创建聊天路由
func NewChatboxRouter(chatService *service.ChatService) *ChatboxRouter {
    return &ChatboxRouter{
        chatService: chatService,
    }
}

// Routes 注册路由
func (r *ChatboxRouter) Routes() []app.Route {
    return []app.Route{
        // 发送消息
        {
            Method: http.MethodPost,
            Path:   "/tenants/:tenant_id/chat/messages",
            Handler: r.SendMessage,
        },

        // 流式对话 (WebSocket/SSE)
        {
            Method: http.MethodGet,
            Path:   "/tenants/:tenant_id/chat/stream",
            Handler: r.StreamChat,
        },

        // 获取对话历史
        {
            Method: http.MethodGet,
            Path:   "/tenants/:tenant_id/chat/conversations/:conversation_id/messages",
            Handler: r.GetMessages,
        },

        // 上传文件
        {
            Method: http.MethodPost,
            Path:   "/tenants/:tenant_id/chat/files",
            Handler: r.UploadFile,
        },

        // 语音转文字
        {
            Method: http.MethodPost,
            Path:   "/tenants/:tenant_id/chat/audio/transcribe",
            Handler: r.TranscribeAudio,
        },

        // 停止生成
        {
            Method: http.MethodPost,
            Path:   "/tenants/:tenant_id/chat/generations/:generation_id/stop",
            Handler: r.StopGeneration,
        },
    }
}
```

### 3.2 发送消息接口

```go
// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
    ConversationID string `json:"conversation_id,omitempty"` // 对话ID,可选
    BotID          string `json:"bot_id,omitempty"`          // 数字员工ID,可选
    Content        string `json:"content" binding:"required"` // 消息内容
    Type           string `json:"type,omitempty"`           // 消息类型 text/image/file
    FileID         string `json:"file_id,omitempty"`         // 文件ID
    Stream         bool   `json:"stream,omitempty"`         // 是否流式输出
}

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
    MessageID      string `json:"message_id"`
    ConversationID string `json:"conversation_id"`
    GenerationID   string `json:"generation_id"`   // 生成ID,用于停止
    Answer         string `json:"answer,omitempty"`      // 如果非流式,直接返回答案
    Sources        []Source `json:"sources,omitempty"`    // 知识来源
    ToolCalls       []ToolCall `json:"tool_calls,omitempty"` // 工具调用
}

// Source 知识来源
type Source struct {
    Type     string `json:"type"`     // knowledge_base/database
    ID       string `json:"id"`       // 知识库ID
    Name     string `json:"name"`     // 知识库名
    Document string `json:"document"` // 文档名
    ChunkID  string `json:"chunk_id"` // 分段ID
    Score    float64 `json:"score"`   // 相似度分数
}

// ToolCall 工具调用
type ToolCall struct {
    ID       string                 `json:"id"`
    Name     string                 `json:"name"`
    Input    map[string]interface{} `json:"input"`
    Output   string                 `json:"output,omitempty"`
    Error    string                 `json:"error,omitempty"`
}

// SendMessage 发送消息
func (r *ChatboxRouter) SendMessage(ctx context.Context, c *app.RequestContext) {
    var req SendMessageRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, ErrorResponse(err))
        return
    }

    // 从 Context 中获取租户ID和用户ID
    tenantID := c.Param("tenant_id")
    userID := ctx.Value("user_id").(string)

    // 调用服务
    result, err := r.chatService.SendMessage(ctx, &service.SendMessageParams{
        TenantID:       tenantID,
        UserID:         userID,
        ConversationID: req.ConversationID,
        BotID:          req.BotID,
        Content:        req.Content,
        Type:           req.Type,
        FileID:         req.FileID,
        Stream:         req.Stream,
    })

    if err != nil {
        c.JSON(500, ErrorResponse(err))
        return
    }

    c.JSON(200, SuccessResponse(result))
}
```

### 3.3 流式对话接口 (SSE)

```go
// StreamChat 流式对话 (SSE)
func (r *ChatboxRouter) StreamChat(ctx context.Context, c *app.RequestContext) {
    // 设置 SSE 响应头
    c.Response.Header.Set("Content-Type", "text/event-stream")
    c.Response.Header.Set("Cache-Control", "no-cache")
    c.Response.Header.Set("Connection", "keep-alive")

    // 获取参数
    tenantID := c.Param("tenant_id")
    userID := ctx.Value("user_id").(string)
    botID := c.Query("bot_id")
    conversationID := c.Query("conversation_id")
    content := c.Query("content")

    // 创建流式响应
    streamChan := make(chan *service.StreamChunk, 10)
    errChan := make(chan error, 1)

    // 异步处理
    go func() {
        defer close(streamChan)
        defer close(errChan)

        err := r.chatService.StreamChat(ctx, &service.StreamChatParams{
            TenantID:       tenantID,
            UserID:         userID,
            BotID:          botID,
            ConversationID: conversationID,
            Content:        content,
            OutputChan:     streamChan,
        })

        if err != nil {
            errChan <- err
        }
    }()

    // 发送 SSE 事件
    c.SetStreamWriter(true)
    streamWriter := c.StreamWriter()

    for {
        select {
        case chunk, ok := <-streamChan:
            if !ok {
                // 发送结束事件
                streamWriter.Write([]byte("event: done\ndata: {}\n\n"))
                return
            }

            // 发送内容事件
            data, _ := json.Marshal(chunk)
            streamWriter.Write([]byte(fmt.Sprintf("event: message\ndata: %s\n\n", data)))

        case err := <-errChan:
            if err != nil {
                // 发送错误事件
                errorData, _ := json.Marshal(map[string]string{"error": err.Error()})
                streamWriter.Write([]byte(fmt.Sprintf("event: error\ndata: %s\n\n", errorData)))
            }
            return

        case <-ctx.Done():
            return
        }
    }
}
```

### 3.4 聊天服务实现

```go
// service/chat_service.go
package service

import (
    "context"
    "time"
)

type ChatService struct {
    db               *gorm.DB
    routingEngine    *engine.RoutingEngine
    memoryEngine     *engine.MemoryEngine
    botService       *BotService
    knowledgeService *KnowledgeService
    pluginService    *PluginService
}

// SendMessageParams 发送消息参数
type SendMessageParams struct {
    TenantID       string
    UserID         string
    ConversationID string
    BotID          string
    Content        string
    Type           string
    FileID         string
    Stream         bool
}

// SendMessage 发送消息
func (s *ChatService) SendMessage(ctx context.Context, params *SendMessageParams) (*SendMessageResponse, error) {
    // 1. 智能路由 (如果没有指定 BotID)
    botID := params.BotID
    if botID == "" {
        routing, err := s.routingEngine.Route(ctx, &engine.RouteRequest{
            TenantID: params.TenantID,
            UserID:   params.UserID,
            Input:    params.Content,
        })
        if err != nil {
            return nil, err
        }
        botID = routing.BotID
    }

    // 2. 获取对话ID
    conversationID := params.ConversationID
    if conversationID == "" {
        // 创建新对话
        conv, err := s.createConversation(ctx, params.TenantID, params.UserID, botID)
        if err != nil {
            return nil, err
        }
        conversationID = conv.ID
    }

    // 3. 保存用户消息
    userMsg := &Message{
        TenantID:       params.TenantID,
        ConversationID: conversationID,
        Role:           "user",
        Content:        params.Content,
        Type:           params.Type,
        CreatedAt:      time.Now(),
    }
    if err := s.db.Create(userMsg).Error; err != nil {
        return nil, err
    }

    // 4. 获取数字员工配置
    bot, err := s.botService.GetBot(ctx, params.TenantID, botID)
    if err != nil {
        return nil, err
    }

    // 5. 加载上下文 (会话记忆)
    context, err := s.memoryEngine.LoadSession(ctx, conversationID, 10)
    if err != nil {
        return nil, err
    }
    context = append(context, &engine.Memory{
        Role:    "user",
        Content: params.Content,
    })

    // 6. 调用数字员工
    response, err := s.invokeBot(ctx, bot, context)
    if err != nil {
        return nil, err
    }

    // 7. 保存助手消息
    assistantMsg := &Message{
        TenantID:       params.TenantID,
        ConversationID: conversationID,
        Role:           "assistant",
        Content:        response.Answer,
        Type:           "text",
        Sources:        response.Sources,
        ToolCalls:      response.ToolCalls,
        CreatedAt:      time.Now(),
    }
    if err := s.db.Create(assistantMsg).Error; err != nil {
        return nil, err
    }

    // 8. 更新会话记忆
    s.memoryEngine.SaveSession(ctx, conversationID, append(context, &engine.Memory{
        Role:    "assistant",
        Content: response.Answer,
    }))

    // 9. 返回响应
    return &SendMessageResponse{
        MessageID:      assistantMsg.ID,
        ConversationID: conversationID,
        Answer:         response.Answer,
        Sources:        response.Sources,
        ToolCalls:      response.ToolCalls,
    }, nil
}

// BotResponse 数字员工响应
type BotResponse struct {
    Answer    string     `json:"answer"`
    Sources   []Source   `json:"sources,omitempty"`
    ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// invokeBot 调用数字员工
func (s *ChatService) invokeBot(ctx context.Context, bot *Bot, context []*engine.Memory) (*BotResponse, error) {
    switch bot.Type {
    case "qa":
        // 问答型数字员工
        return s.invokeQABot(ctx, bot, context)

    case "operation":
        // 操作型数字员工
        return s.invokeOperationBot(ctx, bot, context)

    case "comprehensive":
        // 综合型数字员工
        return s.invokeComprehensiveBot(ctx, bot, context)

    default:
        return nil, errors.New("unsupported bot type")
    }
}

// invokeQABot 调用问答型数字员工
func (s *ChatService) invokeQABot(ctx context.Context, bot *Bot, history []*engine.Memory) (*BotResponse, error) {
    // 获取最后一条用户消息
    lastMessage := history[len(history)-1]
    query := lastMessage.Content

    // 从知识库检索
    config := bot.Config.(map[string]interface{})
    knowledgeBaseIDs := config["knowledge_bases"].([]string)

    searchResults, err := s.knowledgeService.Search(ctx, &KnowledgeSearchParams{
        TenantID:         bot.TenantID,
        KnowledgeBaseIDs: knowledgeBaseIDs,
        Query:            query,
        TopK:             5,
    })
    if err != nil {
        return nil, err
    }

    // 构建提示词
    prompt := s.buildQAPrompt(query, searchResults, history)

    // 调用 LLM
    answer, err := s.callLLM(ctx, bot, prompt)
    if err != nil {
        return nil, err
    }

    // 构建来源
    sources := make([]Source, 0, len(searchResults))
    for _, result := range searchResults {
        sources = append(sources, Source{
            Type:     "knowledge_base",
            ID:       result.KnowledgeBaseID,
            Name:     result.KnowledgeBaseName,
            Document: result.DocumentName,
            ChunkID:  result.ChunkID,
            Score:    result.Score,
        })
    }

    return &BotResponse{
        Answer:  answer,
        Sources: sources,
    }, nil
}

// invokeOperationBot 调用操作型数字员工
func (s *ChatService) invokeOperationBot(ctx context.Context, bot *Bot, history []*engine.Memory) (*BotResponse, error) {
    // 获取最后一条用户消息
    lastMessage := history[len(history)-1]
    input := lastMessage.Content

    // 意图识别 + 任务分解
    tasks, err := s.routingEngine.DecomposeTask(ctx, input)
    if err != nil {
        return nil, err
    }

    // 执行任务
    toolCalls := make([]ToolCall, 0)
    for _, task := range tasks {
        // 调用插件/工具
        result, err := s.pluginService.Invoke(ctx, bot.TenantID, task.PluginID, task.Input)
        if err != nil {
            return nil, err
        }

        toolCalls = append(toolCalls, ToolCall{
            ID:     uuid.New().String(),
            Name:   task.PluginID,
            Input:  task.Input,
            Output: result,
        })
    }

    // 生成总结
    prompt := s.buildOperationSummaryPrompt(input, toolCalls, history)
    answer, err := s.callLLM(ctx, bot, prompt)
    if err != nil {
        return nil, err
    }

    return &BotResponse{
        Answer:    answer,
        ToolCalls: toolCalls,
    }, nil
}

// invokeComprehensiveBot 调用综合型数字员工
func (s *ChatService) invokeComprehensiveBot(ctx context.Context, bot *Bot, history []*engine.Memory) (*BotResponse, error) {
    // 先尝试问答
    qaResponse, err := s.invokeQABot(ctx, bot, history)
    if err == nil && qaResponse.Answer != "" {
        return qaResponse, nil
    }

    // 如果问答无法解决,尝试操作
    return s.invokeOperationBot(ctx, bot, history)
}
```

---

## 四、智能路由引擎

### 4.1 意图识别

```go
// engine/intent_recognizer.go
package engine

type IntentRecognizer struct {
    llmClient *llm.Client
}

// Intent 意图
type Intent struct {
    Type       string  `json:"type"`       // knowledge_query / task_execution / data_analysis
    Confidence float64 `json:"confidence"` // 置信度
    Entities   map[string]interface{} `json:"entities"`   // 提取的实体
}

// Recognize 识别意图
func (r *IntentRecognizer) Recognize(ctx context.Context, input string) (*Intent, error) {
    prompt := fmt.Sprintf(`
你是一个意图识别专家。请分析用户输入,识别其意图。

用户输入: %s

请从以下意图中选择一个:
1. knowledge_query - 知识查询(如: "XX产品的价格是多少?")
2. task_execution - 任务执行(如: "帮我创建一个销售订单")
3. data_analysis - 数据分析(如: "上个月销售额是多少?")

请以JSON格式返回:
{
  "type": "意图类型",
  "confidence": 置信度(0-1),
  "entities": {"提取的实体": "值"}
}
`, input)

    response, err := r.llmClient.Call(ctx, &llm.CallRequest{
        Model: "gpt-3.5-turbo",
        Messages: []llm.Message{
            {Role: "system", Content: "你是一个意图识别专家"},
            {Role: "user", Content: prompt},
        },
        Temperature: 0.1,
    })
    if err != nil {
        return nil, err
    }

    var intent Intent
    if err := json.Unmarshal([]byte(response), &intent); err != nil {
        return nil, err
    }

    return &intent, nil
}
```

### 4.2 路由决策

```go
// engine/routing_engine.go
type RoutingEngine struct {
    intentRecognizer *IntentRecognizer
    botService       *BotService
    pluginService    *PluginService
}

// RouteRequest 路由请求
type RouteRequest struct {
    TenantID string
    UserID   string
    Input    string
}

// RoutingDecision 路由决策
type RoutingDecision struct {
    BotID    string   // 目标数字员工ID
    Reason   string   // 路由原因
    Confidence float64 // 置信度
}

// Route 路由
func (e *RoutingEngine) Route(ctx context.Context, req *RouteRequest) (*RoutingDecision, error) {
    // 1. 意图识别
    intent, err := e.intentRecognizer.Recognize(ctx, req.Input)
    if err != nil {
        return nil, err
    }

    // 2. 根据意图选择数字员工
    var botID string
    switch intent.Type {
    case "knowledge_query":
        // 查找问答型数字员工
        bot, err := e.botService.FindBestQABot(ctx, req.TenantID, intent.Entities)
        if err != nil {
            return nil, err
        }
        botID = bot.ID

    case "task_execution":
        // 查找操作型数字员工
        bot, err := e.botService.FindBestOperationBot(ctx, req.TenantID, intent.Entities)
        if err != nil {
            return nil, err
        }
        botID = bot.ID

    case "data_analysis":
        // 查找数据分析数字员工
        bot, err := e.botService.FindBestAnalyticsBot(ctx, req.TenantID)
        if err != nil {
            return nil, err
        }
        botID = bot.ID

    default:
        // 默认数字员工
        bot, err := e.botService.GetDefaultBot(ctx, req.TenantID)
        if err != nil {
            return nil, err
        }
        botID = bot.ID
    }

    return &RoutingDecision{
        BotID:      botID,
        Reason:     fmt.Sprintf("意图识别为: %s, 置信度: %.2f", intent.Type, intent.Confidence),
        Confidence: intent.Confidence,
    }, nil
}
```

---

## 五、记忆引擎设计

### 5.1 四级记忆模型

```go
// engine/memory_engine.go
type MemoryLevel int

const (
    SessionLevel MemoryLevel = iota    // 会话记忆
    PersonalLevel                       // 个人记忆
    OrgLevel                            // 组织记忆
    EnterpriseLevel                     // 企业记忆
)

// Memory 记忆
type Memory struct {
    Level    MemoryLevel `json:"level"`
    Role     string      `json:"role"`     // user / assistant / system
    Content  string      `json:"content"`
    Metadata map[string]interface{} `json:"metadata"`
    Time     time.Time   `json:"time"`
}

// MemoryEngine 记忆引擎
type MemoryEngine struct {
    sessionCache     *redis.Client
    personalDB       *gorm.DB
    orgDB            *gorm.DB
    enterpriseDB     *gorm.DB
}

// LoadSession 加载会话记忆
func (e *MemoryEngine) LoadSession(ctx context.Context, conversationID string, limit int) ([]*Memory, error) {
    key := fmt.Sprintf("session:%s:messages", conversationID)

    // 从 Redis 获取
    data, err := e.sessionCache.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return []*Memory{}, nil
    }
    if err != nil {
        return nil, err
    }

    var messages []*Memory
    if err := json.Unmarshal(data, &messages); err != nil {
        return nil, err
    }

    // 只返回最近 N 条
    if len(messages) > limit {
        messages = messages[len(messages)-limit:]
    }

    return messages, nil
}

// SaveSession 保存会话记忆
func (e *MemoryEngine) SaveSession(ctx context.Context, conversationID string, messages []*Memory) error {
    key := fmt.Sprintf("session:%s:messages", conversationID)

    data, err := json.Marshal(messages)
    if err != nil {
        return err
    }

    // 保存到 Redis (1小时过期)
    return e.sessionCache.Set(ctx, key, data, 1*time.Hour).Err()
}

// SavePersonal 保存个人记忆
func (e *MemoryEngine) SavePersonal(ctx context.Context, userID string, memory *PersonalMemory) error {
    memory.UserID = userID
    return e.personalDB.WithContext(ctx).Create(memory).Error
}

// SaveOrg 保存组织记忆
func (e *MemoryEngine) SaveOrg(ctx context.Context, orgID string, memory *OrgMemory) error {
    memory.OrganizationID = orgID
    return e.orgDB.WithContext(ctx).Create(memory).Error
}

// SaveEnterprise 保存企业记忆
func (e *MemoryEngine) SaveEnterprise(ctx context.Context, memory *EnterpriseMemory) error {
    return e.enterpriseDB.WithContext(ctx).Create(memory).Error
}

// PersonalMemory 个人记忆
type PersonalMemory struct {
    ID        uint      `gorm:"primarykey"`
    UserID    string    `gorm:"index"`
    Type      string    `gorm:"index"` // preference / history / habit
    Content   string    `gorm:"type:json"`
    CreatedAt time.Time
}

// OrgMemory 组织记忆
type OrgMemory struct {
    ID             uint   `gorm:"primarykey"`
    OrganizationID string `gorm:"index"`
    Type           string `gorm:"index"` // knowledge / experience / best_practice
    Content        string `gorm:"type:json"`
    CreatedAt      time.Time
}

// EnterpriseMemory 企业记忆
type EnterpriseMemory struct {
    ID        uint      `gorm:"primarykey"`
    Type      string    `gorm:"index"` // strategy / culture / wisdom
    Content   string    `gorm:"type:json"`
    Version   int
    CreatedAt time.Time
}
```

---

## 六、数据库设计

### 6.1 对话相关表

```sql
-- 对话表
CREATE TABLE conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL,
    bot_id BIGINT NOT NULL,
    type ENUM('chat', 'workflow', 'multi_agent') DEFAULT 'chat',
    status ENUM('active', 'archived') DEFAULT 'active',
    title VARCHAR(512),
    message_count INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 消息表
CREATE TABLE messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    conversation_id BIGINT NOT NULL,
    role ENUM('user', 'assistant', 'system') NOT NULL,
    content TEXT,
    type ENUM('text', 'image', 'file', 'code', 'tool_call') DEFAULT 'text',
    metadata JSON,

    -- 知识来源
    sources JSON,

    -- 工具调用
    tool_calls JSON,

    -- 反馈
    rating INT,
    feedback TEXT,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    INDEX idx_tenant_conv (tenant_id, conversation_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 文件表
CREATE TABLE chat_files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    conversation_id BIGINT,
    user_id BIGINT NOT NULL,
    file_name VARCHAR(512) NOT NULL,
    file_type VARCHAR(64),
    file_size BIGINT,
    file_url VARCHAR(1024),
    storage_key VARCHAR(1024),
    status ENUM('uploading', 'completed', 'failed') DEFAULT 'uploading',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE SET NULL,
    INDEX idx_tenant_conv (tenant_id, conversation_id),
    INDEX idx_tenant_user (tenant_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 七、前端Hook设计

### 7.1 useChatbox Hook

```typescript
// src/hooks/useChatbox.ts
import { useState, useCallback, useRef } from 'react';
import { useWebSocket } from './useWebSocket';

interface UseChatboxOptions {
  tenantId: string;
  botId?: string;
  onMessage?: (message: Message) => void;
  onError?: (error: Error) => void;
}

export function useChatbox(options: UseChatboxOptions) {
  const [messages, setMessages] = useState<Message[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isStreaming, setIsStreaming] = useState(false);
  const generationIdRef = useRef<string>();

  // WebSocket 连接
  const { ws, send, disconnect } = useWebSocket(
    `/api/v1/tenants/${options.tenantId}/chat/stream`
  );

  // 发送消息
  const sendMessage = useCallback(async (input: string) => {
    if (!input.trim()) return;

    setIsLoading(true);
    setIsStreaming(true);

    // 添加用户消息
    const userMessage: Message = {
      id: uuid(),
      role: 'user',
      content: input,
      timestamp: new Date(),
    };
    setMessages(prev => [...prev, userMessage]);

    try {
      // 发送到后端
      send({
        type: 'message',
        data: {
          content: input,
          bot_id: options.botId,
        },
      });
    } catch (error) {
      options.onError?.(error as Error);
      setIsLoading(false);
      setIsStreaming(false);
    }
  }, [options.tenantId, options.botId]);

  // 停止生成
  const stopGeneration = useCallback(async () => {
    if (!generationIdRef.current) return;

    await fetch(`/api/v1/tenants/${options.tenantId}/chat/generations/${generationIdRef.current}/stop`, {
      method: 'POST',
    });

    setIsStreaming(false);
    setIsLoading(false);
  }, [options.tenantId]);

  // 重新生成
  const regenerate = useCallback(async (messageId: string) => {
    // 实现重新生成逻辑
  }, [options.tenantId]);

  return {
    messages,
    isLoading,
    isStreaming,
    sendMessage,
    stopGeneration,
    regenerate,
  };
}
```

---

## 八、性能优化

### 8.1 消息列表虚拟滚动

```typescript
// src/components/chatbox/VirtualizedMessageList.tsx
import { useVirtualizer } from '@tanstack/react-virtual';

interface VirtualizedMessageListProps {
  messages: Message[];
}

const VirtualizedMessageList: React.FC<VirtualizedMessageListProps> = ({ messages }) => {
  const parentRef = useRef<HTMLDivElement>(null);

  const rowVirtualizer = useVirtualizer({
    count: messages.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 100, // 估算每条消息高度
    overscan: 5,
  });

  return (
    <div ref={parentRef} className="message-list">
      <div
        style={{
          height: `${rowVirtualizer.getTotalSize()}px`,
          width: '100%',
          position: 'relative',
        }}
      >
        {rowVirtualizer.getVirtualItems().map(virtualRow => (
          <div
            key={virtualRow.key}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: `${virtualRow.size}px`,
              transform: `translateY(${virtualRow.start}px)`,
            }}
          >
            <MessageRenderer message={messages[virtualRow.index]} />
          </div>
        ))}
      </div>
    </div>
  );
};
```

### 8.2 消息缓存

```typescript
// src/utils/messageCache.ts
class MessageCache {
  private cache = new Map<string, Message[]>();
  private maxCacheSize = 100;

  get(conversationId: string): Message[] | undefined {
    return this.cache.get(conversationId);
  }

  set(conversationId: string, messages: Message[]) {
    // LRU 缓存淘汰
    if (this.cache.size >= this.maxCacheSize) {
      const firstKey = this.cache.keys().next().value;
      this.cache.delete(firstKey);
    }

    this.cache.set(conversationId, messages);
  }

  clear() {
    this.cache.clear();
  }
}

export const messageCache = new MessageCache();
```

---

## 九、测试用例

### 9.1 单元测试

```typescript
// src/components/chatbox/__tests__/InputArea.test.tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { InputArea } from '../InputArea';

describe('InputArea', () => {
  it('should send message on Enter key press', () => {
    const mockOnSend = jest.fn();
    const { getByRole } = render(
      <InputArea
        value="Hello"
        onChange={() => {}}
        onSend={mockOnSend}
        onVoiceStart={() => {}}
        onVoiceStop={() => {}}
        onFileUpload={() => {}}
      />
    );

    const textarea = getByRole('textbox');
    fireEvent.keyDown(textarea, { key: 'Enter', code: 'Enter' });

    expect(mockOnSend).toHaveBeenCalledTimes(1);
  });

  it('should not send message on Shift+Enter', () => {
    const mockOnSend = jest.fn();
    const { getByRole } = render(
      <InputArea
        value="Hello"
        onChange={() => {}}
        onSend={mockOnSend}
        onVoiceStart={() => {}}
        onVoiceStop={() => {}}
        onFileUpload={() => {}}
      />
    );

    const textarea = getByRole('textbox');
    fireEvent.keyDown(textarea, { key: 'Enter', shiftKey: true });

    expect(mockOnSend).not.toHaveBeenCalled();
  });
});
```

---

## 十、总结

### 核心特性

- ✅ 多模态输入 (文字/语音/图片/文件)
- ✅ 智能路由 (自动识别意图并分配数字员工)
- ✅ 四级记忆 (会话/个人/组织/企业)
- ✅ 实时响应 (流式输出, < 2s)
- ✅ 上下文理解 (多轮对话)

### 技术栈

**前端**:
- React 18 + TypeScript
- WebSocket/SSE 实时通信
- Monaco Editor 代码编辑
- react-markdown Markdown 渲染
- TanStack Virtual 虚拟滚动

**后端**:
- Go 1.23 + Hertz 框架
- GORM 数据库 ORM
- Redis 会话缓存
- GPT-3.5/4 意图识别
- Milvus 向量检索

### 工作量评估

- 前端开发: 4 周
- 后端开发: 4 周
- 测试: 2 周
- **总计**: 10 周

---

**文档结束**
