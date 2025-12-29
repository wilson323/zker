# 01-ZKER前台_会话_SuperChatbox 详细设计说明书

**文档编号**: DE-DD-2025-001
**模块名称**: 会话_SuperChatbox (SuperChatbox)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03
**最后更新**: 2025-01-03

---

## 📋 文档修订历史

| 版本 | 日期 | 修订人 | 修订说明 |
|------|------|--------|----------|
| v1.0.0 | 2025-01-03 | ZKER Team | 初始版本，基于独立编码设计 |

---

## 1. 模块概述

### 1.1 模块定位

**会话_SuperChatbox** 是 ZKER 企业级 SaaS 平台的核心前台模块，提供**企业级智能对话UI和交互能力**。作为用户与AI Bot交互的主要入口，支持多轮对话、流式输出、多模态输入、快捷指令等企业级特性。

**核心设计理念**：
- ✅ **独立编码**：60% 前端UI编码 + 40% 后端逻辑编码
- ✅ **配置化主题**：支持企业自定义对话主题和样式
- ✅ **实时交互**：WebSocket 双向通信，流式输出
- ✅ **多模态支持**：文本、图片、文件等多种输入方式
- ✅ **企业级**：会话隔离、权限控制、审计日志

### 1.2 业务价值

| 受益者 | 价值 |
|--------|------|
| **企业用户** | 自然流畅的对话体验，支持复杂任务，多轮对话无障碍 |
| **企业管理员** | 自定义企业对话主题，配置快捷指令，监控对话质量 |
| **平台运营** | 提升用户活跃度和留存率，收集对话数据优化Bot |

### 1.3 与鲸智百应功能对齐

| 鲸智百应功能 | ZKER 实现方式 | 实现策略 |
|-------------|--------------|----------|
| 多轮对话 | ✅ WebSocket + 会话管理 | 独立编码 |
| 流式输出 | ✅ Server-Sent Events | 独立编码 |
| 快捷指令 | ✅ quick_commands 表配置 | 数据库配置 |
| 对话历史 | ✅ messages 表存储 | 数据库存储 |
| 会话分享 | ✅ shared_conversations 表 | 数据库配置 |
| 多模态输入 | ✅ 文本/图片/文件上传 | 独立编码 |
| 提示词建议 | ✅ prompt_suggestions 表配置 | 数据库配置 |
| 引用来源 | ✅ RAG 引用展示 | 独立编码 |
| 评价反馈 | ✅ message_ratings 表 | 数据库存储 |

**实现策略**：✅ 60% 编码（UI组件 + WebSocket） + 40% 配置（主题、快捷指令）

---

## 2. 功能需求

### 2.1 核心功能清单

#### 2.1.1 对话交互（核心功能）

**F1 - 消息发送与接收**
- F1.1 文本消息输入（支持多行输入）
- F1.2 图片消息上传（拖拽/选择文件）
- F1.3 文件消息上传（支持PDF、Word、Excel等）
- F1.4 语音消息输入（可选，录音转文字）
- F1.5 消息编辑（编辑已发送消息）
- F1.6 消息删除（删除已发送消息）
- F1.7 消息重新生成（让AI重新回答）
- F1.8 消息复制（一键复制消息内容）

**F2 - 流式输出**
- F2.1 逐字输出（打字机效果）
- F2.2 停止生成（中断AI输出）
- F2.3 继续生成（从停止处继续）
- F2.4 Markdown 渲染（支持代码高亮、表格、公式）
- F2.5 引用来源展示（显示知识库引用）

**F3 - 多轮对话**
- F3.1 会话上下文管理（保留对话历史）
- F3.2 会话切换（在不同会话间切换）
- F3.3 会话新建/删除/重命名
- F3.4 会话归档（归档旧会话）
- F3.5 会话导出（导出为Markdown/PDF）

#### 2.1.2 快捷操作

**F4 - 快捷指令**
- F4.1 快捷指令面板（展示预置指令）
- F4.2 指令点击填充（一键填充到输入框）
- F4.3 自定义指令（用户/企业自定义）

**F5 - 提示词建议**
- F5.1 智能推荐（根据上下文推荐提示词）
- F5.2 提示词模板（常用提示词模板）

**F6 - 消息操作**
- F6.1 消息评价（点赞/点踩）
- F6.2 消息评论（对消息添加评论）
- F6.3 消息收藏（收藏重要消息）
- F6.4 消息引用（引用某条消息进行回复）

#### 2.1.3 高级功能

**F7 - 多模态支持**
- F7.1 图片识别（OCR提取文字）
- F7.2 文件解析（PDF/Word等文件内容提取）
- F7.3 语音输入（语音转文字）

**F8 - 协作功能**
- F8.1 会话分享（生成分享链接）
- F8.2 会话协作（多人共同对话）
- F8.3 消息转发（转发消息到其他会话）

### 2.2 非功能需求

| 需求类型 | 指标 | 说明 |
|---------|------|------|
| **性能** | 消息响应时间 | 首字输出 < 500ms |
| **性能** | 流式输出延迟 | 字间延迟 < 50ms |
| **可用性** | 系统可用性 | 99.9% （月度） |
| **并发** | 同时在线用户 | 单租户 100，系统 10000 |
| **兼容性** | 浏览器支持 | Chrome/Firefox/Safari/Edge 最新版 |
| **移动端** | 响应式设计 | 支持手机/平板访问 |

---

## 3. 架构设计

### 3.1 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      前端层 (React 18)                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │ Chatbox         │  │ MessageList     │  │ InputArea       │  │
│  │ (主容器)         │  │ (消息列表)      │  │ (输入区域)      │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  │
│  │ Sidebar         │  │ QuickCommands   │  │ AttachmentPanel │  │
│  │ (会话列表)       │  │ (快捷指令)      │  │ (附件面板)      │  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                    ↓ WebSocket + HTTP API
┌─────────────────────────────────────────────────────────────────┐
│                    API 网关层 (Hertz)                           │
├─────────────────────────────────────────────────────────────────┤
│  租户识别中间件 → 权限验证 → 会话权限检查 → 路由分发             │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    业务服务层 (Go 1.23)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         ConversationService (核心服务)                   │   │
│  │  - CreateConversation()    (创建会话)                   │   │
│  │  - SendMessage()           (发送消息)                   │   │
│  │  - StreamResponse()        (流式响应)                   │   │
│  │  - EditMessage()           (编辑消息)                   │   │
│  │  - DeleteMessage()         (删除消息)                   │   │
│  │  - RegenerateResponse()    (重新生成)                   │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         MessageHandler (消息处理器)                      │   │
│  │  - HandleTextMessage()     (处理文本消息)                │   │
│  │  - HandleImageMessage()    (处理图片消息)                │   │
│  │  - HandleFileMessage()     (处理文件消息)                │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│              Bot执行层 (复用 zker Bot 系统)                      │
├─────────────────────────────────────────────────────────────────┤
│  BotExecutor → PluginExecutor → LLMCaller → ResponseGenerator   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    数据访问层 (GORM)                             │
├─────────────────────────────────────────────────────────────────┤
│  conversations | messages | message_attachments                 │
│  quick_commands | prompt_suggestions | shared_conversations     │
│  message_ratings | message_comments                            │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    基础设施层                                    │
├─────────────────────────────────────────────────────────────────┤
│  MySQL 8.4 | Redis 8.0 | MinIO | WebSocket Server              │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 核心流程

#### 3.2.1 发送消息与流式输出流程

```
用户输入消息
        ↓
前端预验证（长度、格式）
        ↓
通过 WebSocket 发送消息
        ↓
后端接收消息
        ↓
创建用户消息记录（messages 表，role = 'user'）
        ↓
加载会话上下文（历史消息 N 条）
        ↓
调用 BotExecutor 执行 Bot
        ↓
Bot 处理流程：
  ├─ 检测意图
  ├─ 加载知识库（如果配置）
  ├─ 调用 Plugin（如果需要）
  ├─ 构建 LLM 提示词
  └─ 调用 LLM 生成响应
        ↓
流式返回 LLM 输出（WebSocket）
        ↓
前端逐字渲染（打字机效果）
        ↓
创建助手消息记录（messages 表，role = 'assistant'）
        ↓
更新会话最后活跃时间
```

#### 3.2.2 消息重新生成流程

```
用户点击"重新生成"
        ↓
前端通过 WebSocket 发送 regenerate 请求
        ↓
后端删除原助手消息（软删除）
        ↓
重新调用 BotExecutor（使用相同上下文）
        ↓
流式返回新响应
        ↓
前端渲染新响应
```

---

## 4. 数据库设计

### 4.1 核心表结构

#### 4.1.1 会话表 (conversations)

```sql
CREATE TABLE conversations (
    id VARCHAR(64) PRIMARY KEY COMMENT '会话ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    bot_id VARCHAR(64) COMMENT '关联的Bot ID (NULL表示自由对话)',
    title VARCHAR(200) COMMENT '会话标题',
    description VARCHAR(500) COMMENT '会话描述',

    -- 模型配置（可选，覆盖Bot默认配置）
    model_config JSON COMMENT '模型配置 {"temperature": 0.7, "max_tokens": 2000}',

    -- 统计
    message_count INT DEFAULT 0 COMMENT '消息数量',
    total_tokens INT DEFAULT 0 COMMENT '总Token数',
    last_message_at DATETIME COMMENT '最后消息时间',

    -- 状态
    status ENUM('active', 'archived', 'deleted') DEFAULT 'active' COMMENT '状态',
    is_pinned BOOLEAN DEFAULT FALSE COMMENT '是否置顶',

    -- 分享
    is_shared BOOLEAN DEFAULT FALSE COMMENT '是否分享',
    share_token VARCHAR(64) COMMENT '分享Token',
    share_expires_at DATETIME COMMENT '分享过期时间',

    -- 审计
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_status (status),
    INDEX idx_last_message_at (last_message_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='会话表';
```

#### 4.1.2 消息表 (messages)

```sql
CREATE TABLE messages (
    id VARCHAR(64) PRIMARY KEY COMMENT '消息ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    bot_id VARCHAR(64) COMMENT 'Bot ID',

    -- 消息内容
    role ENUM('user', 'assistant', 'system') NOT NULL COMMENT '角色',
    content TEXT NOT NULL COMMENT '消息内容 (Markdown格式)',
    content_type ENUM('text', 'image', 'file', 'mixed') DEFAULT 'text' COMMENT '内容类型',

    -- 元数据
    metadata JSON COMMENT '元数据 (token数、模型、引用等)',
    parent_id VARCHAR(64) COMMENT '父消息ID (用于重新生成)',

    -- 状态
    status ENUM('sending', 'sent', 'failed', 'deleted') DEFAULT 'sent' COMMENT '状态',
    error_message TEXT COMMENT '错误信息',

    -- 时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,

    INDEX idx_tenant_conv (tenant_id, conversation_id),
    INDEX idx_user_id (user_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='消息表';
```

**元数据 JSON 示例**：

```json
{
  "model": "gpt-4",
  "tokens": {
    "prompt": 150,
    "completion": 300,
    "total": 450
  },
  "latency_ms": 2340,
  "citations": [
    {
      "document_id": "doc-001",
      "chunk_id": 123,
      "text": "引用的文本片段...",
      "score": 0.92
    }
  ],
  "plugins": ["plugin-email-sender-001"]
}
```

#### 4.1.3 消息附件表 (message_attachments)

```sql
CREATE TABLE message_attachments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '附件ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    message_id VARCHAR(64) NOT NULL COMMENT '消息ID',
    file_type ENUM('image', 'file') NOT NULL COMMENT '文件类型',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    file_path VARCHAR(512) NOT NULL COMMENT '文件存储路径 (MinIO)',
    file_hash VARCHAR(64) COMMENT '文件哈希',
    mime_type VARCHAR(100) COMMENT 'MIME类型',
    metadata JSON COMMENT '附件元数据 (图片尺寸、页数等)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_msg (tenant_id, message_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='消息附件表';
```

#### 4.1.4 快捷指令表 (quick_commands)

**核心设计**：通过数据库配置实现快捷指令，支持用户/企业自定义

```sql
CREATE TABLE quick_commands (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '指令ID',
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台预置)',
    user_id BIGINT COMMENT '用户ID (NULL表示企业级，非NULL表示用户自定义)',
    bot_id VARCHAR(64) COMMENT '关联的Bot ID (NULL表示通用)',
    category VARCHAR(50) COMMENT '分类',

    -- 指令内容
    title VARCHAR(100) NOT NULL COMMENT '指令标题',
    description VARCHAR(200) COMMENT '指令描述',
    prompt TEXT NOT NULL COMMENT '提示词内容',
    variables JSON COMMENT '变量定义',

    -- 展示配置
    icon VARCHAR(10) COMMENT '图标 emoji',
    color VARCHAR(20) COMMENT '颜色',
    sort_order INT DEFAULT 0 COMMENT '排序',

    -- 统计
    usage_count INT DEFAULT 0 COMMENT '使用次数',

    -- 状态
    is_active BOOLEAN DEFAULT TRUE,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_user_id (user_id),
    INDEX idx_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='快捷指令表';
```

**快捷指令示例数据**：

```sql
-- 平台预置指令
INSERT INTO quick_commands (tenant_id, bot_id, category, title, description, prompt, icon, sort_order) VALUES
(NULL, NULL, 'writing', '写邮件', '快速撰写商务邮件', '请帮我写一封商务邮件，收件人：{{收件人}}，主题：{{主题}}，内容要点：\n{{内容}}', '📧', 1),
(NULL, NULL, 'writing', '总结文档', '总结文档要点', '请帮我总结以下文档的要点：\n{{文档内容}}', '📝', 2),
(NULL, NULL, 'coding', '写代码', '生成代码片段', '请帮我写一段代码，语言：{{语言}}，功能：{{功能描述}}', '💻', 3);

-- 企业自定义指令
INSERT INTO quick_commands (tenant_id, category, title, description, prompt, icon) VALUES
('tenant-abc-001', NULL, '写周报', '快速撰写工作周报', '请帮我写一份周报，本周完成的工作：\n{{工作内容}}', '📊');
```

#### 4.1.5 提示词建议表 (prompt_suggestions)

```sql
CREATE TABLE prompt_suggestions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '建议ID',
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台预置)',
    bot_id VARCHAR(64) COMMENT '关联的Bot ID',
    suggestion_type ENUM('contextual', 'template', 'history') NOT NULL COMMENT '建议类型',

    title VARCHAR(100) NOT NULL COMMENT '建议标题',
    prompt TEXT NOT NULL COMMENT '提示词内容',
    trigger_keywords JSON COMMENT '触发关键词列表',

    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_bot (tenant_id, bot_id),
    INDEX idx_type (suggestion_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='提示词建议表';
```

#### 4.1.6 消息评价表 (message_ratings)

```sql
CREATE TABLE message_ratings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '评价ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    message_id VARCHAR(64) NOT NULL COMMENT '消息ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    rating ENUM('thumbs_up', 'thumbs_down') NOT NULL COMMENT '评价',
    feedback TEXT COMMENT '反馈文本',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_message_user (message_id, user_id),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='消息评价表';
```

#### 4.1.7 会话分享表 (shared_conversations)

```sql
CREATE TABLE shared_conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分享ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    share_token VARCHAR(64) NOT NULL UNIQUE COMMENT '分享Token',
    shared_by BIGINT NOT NULL COMMENT '分享者ID',
    expires_at DATETIME COMMENT '过期时间 (NULL表示永久)',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_conv (tenant_id, conversation_id),
    INDEX idx_share_token (share_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='会话分享表';
```

---

## 5. API 设计

### 5.1 WebSocket API

#### 5.1.1 连接与消息

**连接 URL**：
```
wss://api.zker.com/v1/chat/ws?conversation_id={conv_id}&token={jwt_token}
```

**客户端消息格式**：

```json
// 发送文本消息
{
  "type": "message.send",
  "data": {
    "content": "你好，请帮我写一封邮件",
    "messageId": "msg-client-123"
  }
}

// 停止生成
{
  "type": "generation.stop",
  "data": {
    "messageId": "msg-assistant-456"
  }
}

// 重新生成
{
  "type": "message.regenerate",
  "data": {
    "parentMessageId": "msg-assistant-456"
  }
}

// 编辑消息
{
  "type": "message.edit",
  "data": {
    "messageId": "msg-user-123",
    "newContent": "修改后的内容"
  }
}
```

**服务端消息格式**：

```json
// 消息已接收
{
  "type": "message.received",
  "data": {
    "messageId": "msg-client-123",
    "tempId": "msg-server-789"
  }
}

// 开始流式输出
{
  "type": "generation.start",
  "data": {
    "messageId": "msg-assistant-456",
    "conversationId": "conv-001"
  }
}

// 流式内容片段
{
  "type": "generation.content",
  "data": {
    "messageId": "msg-assistant-456",
    "delta": "您好，",
    "isComplete": false
  }
}

// 流式输出结束
{
  "type": "generation.end",
  "data": {
    "messageId": "msg-assistant-456",
    "tokens": {
      "prompt": 150,
      "completion": 300,
      "total": 450
    },
    "latencyMs": 2340
  }
}

// 错误
{
  "type": "error",
  "data": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "请求过于频繁，请稍后再试"
  }
}
```

#### 5.1.2 Go WebSocket 服务实现

```go
package chat

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/gorilla/websocket"
)

// WebSocketMessage WebSocket消息
type WebSocketMessage struct {
    Type string                 `json:"type"`
    Data map[string]interface{} `json:"data"`
}

// ChatWebSocketHandler WebSocket处理器
type ChatWebSocketHandler struct {
    upgrader        websocket.Upgrader
    conversationSvc *ConversationService
    messageSvc      *MessageService
    botExecutor     *BotExecutor
    logger          *zap.Logger
}

// HandleWebSocket 处理WebSocket连接
func (h *ChatWebSocketHandler) HandleWebSocket(c *app.RequestContext) {
    // 1. 验证Token
    token := c.Query("token")
    claims, err := jwt.ParseToken(token)
    if err != nil {
        c.JSON(401, map[string]string{"error": "invalid token"})
        return
    }

    conversationID := c.Query("conversation_id")

    // 2. 升级到WebSocket
    conn, err := h.upgrader.Upgrade(c.Response.Writer(), c.Request(), nil)
    if err != nil {
        h.logger.Error("failed to upgrade websocket", zap.Error(err))
        return
    }
    defer conn.Close()

    // 3. 创建会话
    session := &ChatSession{
        Conn:           conn,
        TenantID:       claims.TenantID,
        UserID:         claims.UserID,
        ConversationID: conversationID,
        SendChan:       make(chan *WebSocketMessage, 100),
    }

    // 4. 启动发送协程
    go session.sendLoop()

    // 5. 启动接收协程
    session.receiveLoop()
}

// receiveLoop 接收消息循环
func (s *ChatSession) receiveLoop() {
    for {
        _, data, err := s.Conn.ReadMessage()
        if err != nil {
            s.logger.Error("failed to read message", zap.Error(err))
            break
        }

        var msg WebSocketMessage
        if err := json.Unmarshal(data, &msg); err != nil {
            s.logger.Error("failed to unmarshal message", zap.Error(err))
            continue
        }

        // 处理消息
        go s.handleMessage(&msg)
    }
}

// handleMessage 处理消息
func (s *ChatSession) handleMessage(msg *WebSocketMessage) {
    switch msg.Type {
    case "message.send":
        s.handleSendMessage(msg)
    case "generation.stop":
        s.handleStopGeneration(msg)
    case "message.regenerate":
        s.handleRegenerate(msg)
    case "message.edit":
        s.handleEditMessage(msg)
    default:
        s.logger.Warn("unknown message type", zap.String("type", msg.Type))
    }
}

// handleSendMessage 处理发送消息
func (s *ChatSession) handleSendMessage(msg *WebSocketMessage) {
    content := msg.Data["content"].(string)
    clientMsgID := msg.Data["messageId"].(string)

    // 1. 创建用户消息
    userMsg := &model.Message{
        ID:             generateMessageID(),
        TenantID:       s.TenantID,
        ConversationID: s.ConversationID,
        UserID:         s.UserID,
        Role:           "user",
        Content:        content,
        Status:         "sent",
    }

    if err := s.messageSvc.Create(context.Background(), userMsg); err != nil {
        s.SendError("failed to create message", err)
        return
    }

    // 2. 确认消息已接收
    s.Send(&WebSocketMessage{
        Type: "message.received",
        Data: map[string]interface{}{
            "messageId": clientMsgID,
            "tempId":    userMsg.ID,
        },
    })

    // 3. 加载会话上下文
    messages, err := s.messageSvc.ListByConversationID(context.Background(), s.ConversationID, 10)
    if err != nil {
        s.SendError("failed to load context", err)
        return
    }

    // 4. 调用Bot执行
    assistantMsg := &model.Message{
        ID:             generateMessageID(),
        TenantID:       s.TenantID,
        ConversationID: s.ConversationID,
        UserID:         s.UserID,
        Role:           "assistant",
        Status:         "sending",
    }

    // 5. 开始流式生成
    s.Send(&WebSocketMessage{
        Type: "generation.start",
        Data: map[string]interface{}{
            "messageId":      assistantMsg.ID,
            "conversationId": s.ConversationID,
        },
    })

    // 6. 执行Bot并流式输出
    stream := s.botExecutor.ExecuteStream(context.Background(), &BotExecuteRequest{
        ConversationID: s.ConversationID,
        Messages:       messages,
        BotID:          s.BotID,
    })

    var contentBuilder strings.Builder
    startTime := time.Now()

    for chunk := range stream {
        if chunk.Error != nil {
            s.SendError("generation failed", chunk.Error)
            assistantMsg.Status = "failed"
            assistantMsg.ErrorMessage = chunk.Error.Error()
            s.messageSvc.Update(context.Background(), assistantMsg)
            return
        }

        contentBuilder.WriteString(chunk.Delta)

        // 发送内容片段
        s.Send(&WebSocketMessage{
            Type: "generation.content",
            Data: map[string]interface{}{
                "messageId":  assistantMsg.ID,
                "delta":      chunk.Delta,
                "isComplete": false,
            },
        })
    }

    // 7. 保存助手消息
    assistantMsg.Content = contentBuilder.String()
    assistantMsg.Status = "sent"
    assistantMsg.Metadata = map[string]interface{}{
        "latency_ms": time.Since(startTime).Milliseconds(),
        "tokens":     chunk.Tokens,
    }

    if err := s.messageSvc.Create(context.Background(), assistantMsg); err != nil {
        s.SendError("failed to save message", err)
        return
    }

    // 8. 发送结束事件
    s.Send(&WebSocketMessage{
        Type: "generation.end",
        Data: map[string]interface{}{
            "messageId": assistantMsg.ID,
            "tokens":    chunk.Tokens,
            "latencyMs": time.Since(startTime).Milliseconds(),
        },
    })
}

// Send 发送消息
func (s *ChatSession) Send(msg *WebSocketMessage) error {
    data, err := json.Marshal(msg)
    if err != nil {
        return err
    }

    s.SendChan <- msg
    return nil
}

// SendError 发送错误
func (s *ChatSession) SendError(message string, err error) {
    s.Send(&WebSocketMessage{
        Type: "error",
        Data: map[string]interface{}{
            "code":    "INTERNAL_ERROR",
            "message": message,
            "detail":  err.Error(),
        },
    })
}

// sendLoop 发送循环
func (s *ChatSession) sendLoop() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case msg := <-s.SendChan:
            data, err := json.Marshal(msg)
            if err != nil {
                s.logger.Error("failed to marshal message", zap.Error(err))
                continue
            }

            if err := s.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
                s.logger.Error("failed to send message", zap.Error(err))
                return
            }

        case <-ticker.C:
            // 心跳
            if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                s.logger.Error("failed to send ping", zap.Error(err))
                return
            }
        }
    }
}
```

### 5.2 REST API

| 方法 | 路径 | 功能 |
|------|------|------|
| **会话管理** |
| POST | /api/v1/chat/conversations | 创建会话 |
| GET | /api/v1/chat/conversations | 列出会话 |
| GET | /api/v1/chat/conversations/:id | 查看会话详情 |
| PUT | /api/v1/chat/conversations/:id | 更新会话 |
| DELETE | /api/v1/chat/conversations/:id | 删除会话 |
| **消息管理** |
| GET | /api/v1/chat/conversations/:id/messages | 列出消息 |
| GET | /api/v1/chat/messages/:id | 查看消息详情 |
| PUT | /api/v1/chat/messages/:id | 编辑消息 |
| DELETE | /api/v1/chat/messages/:id | 删除消息 |
| POST | /api/v1/chat/messages/:id/regenerate | 重新生成 |
| **快捷指令** |
| GET | /api/v1/chat/quick-commands | 列出快捷指令 |
| POST | /api/v1/chat/quick-commands | 创建自定义指令 |
| **分享** |
| POST | /api/v1/chat/conversations/:id/share | 创建分享链接 |
| GET | /api/v1/chat/shared/:token | 查看分享会话 |

---

## 6. 前端设计

### 6.1 核心组件

#### 6.1.1 Chatbox（主容器）

```tsx
import React, { useState, useEffect, useRef } from 'react';
import { ChatboxSidebar } from './ChatboxSidebar';
import { MessageList } from './MessageList';
import { InputArea } from './InputArea';
import { QuickCommands } from './QuickCommands';
import { chatAPI } from '@/services/api';
import { useWebSocket } from '@/hooks/useWebSocket';
import styles from './Chatbox.module.scss';

interface ChatboxProps {
  botId?: string;
  conversationId?: string;
}

export const Chatbox: React.FC<ChatboxProps> = ({ botId, conversationId }) => {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [currentConversation, setCurrentConversation] = useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [isGenerating, setIsGenerating] = useState(false);

  // WebSocket连接
  const { sendMessage, lastMessage, readyState } = useWebSocket(
    chatAPI.getWebSocketURL(currentConversation?.id),
    {
      onMessage: (msg) => handleWebSocketMessage(msg),
    }
  );

  useEffect(() => {
    loadConversations();
    if (conversationId) {
      loadConversation(conversationId);
    }
  }, [conversationId]);

  const handleWebSocketMessage = (msg: WebSocketMessage) => {
    switch (msg.type) {
      case 'message.received':
        // 消息已接收确认
        break;
      case 'generation.start':
        setIsGenerating(true);
        // 创建空的助手消息
        setMessages(prev => [...prev, {
          id: msg.data.messageId,
          role: 'assistant',
          content: '',
          status: 'generating',
        }]);
        break;
      case 'generation.content':
        // 追加内容
        setMessages(prev => prev.map(msg =>
          msg.id === msg.data.messageId
            ? { ...msg, content: msg.content + msg.data.delta }
            : msg
        ));
        break;
      case 'generation.end':
        setIsGenerating(false);
        // 更新消息元数据
        setMessages(prev => prev.map(msg =>
          msg.id === msg.data.messageId
            ? { ...msg, status: 'sent', metadata: msg.data }
            : msg
        ));
        break;
      case 'error':
        message.error(msg.data.message);
        setIsGenerating(false);
        break;
    }
  };

  const handleSendMessage = (content: string, attachments?: File[]) => {
    const tempId = `temp-${Date.now()}`;

    // 立即显示用户消息
    setMessages(prev => [...prev, {
      id: tempId,
      role: 'user',
      content,
      status: 'sending',
      attachments,
    }]);

    // 通过WebSocket发送
    sendMessage({
      type: 'message.send',
      data: {
        content,
        messageId: tempId,
      },
    } as any);
  };

  const handleStopGeneration = () => {
    if (isGenerating && currentConversation) {
      sendMessage({
        type: 'generation.stop',
        data: { messageId: messages[messages.length - 1].id },
      } as any);
    }
  };

  const handleRegenerate = (messageId: string) => {
    sendMessage({
      type: 'message.regenerate',
      data: { parentMessageId: messageId },
    } as any);
  };

  return (
    <div className={styles.chatbox}>
      <ChatboxSidebar
        conversations={conversations}
        currentConversation={currentConversation}
        onSelectConversation={loadConversation}
        onCreateConversation={createConversation}
      />

      <div className={styles.main}>
        <MessageList
          messages={messages}
          isGenerating={isGenerating}
          onRegenerate={handleRegenerate}
          onEditMessage={handleEditMessage}
          onDeleteMessage={handleDeleteMessage}
        />

        <InputArea
          onSendMessage={handleSendMessage}
          onStopGeneration={handleStopGeneration}
          isGenerating={isGenerating}
        />

        <QuickCommands
          botId={botId}
          onSelectCommand={handleSelectCommand}
        />
      </div>
    </div>
  );
};
```

#### 6.1.2 MessageList（消息列表）

```tsx
import React, { useRef, useEffect } from 'react';
import { Avatar, Button, Dropdown } from '@douyinfe/semi-ui';
import { IconThumbsUp, IconThumbsDown, IconRefresh, IconEdit3, IconTrash } from '@douyinfe/semi-icons';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighting';
import styles from './MessageList.module.scss';

interface MessageListProps {
  messages: Message[];
  isGenerating: boolean;
  onRegenerate: (messageId: string) => void;
  onEditMessage: (messageId: string, newContent: string) => void;
  onDeleteMessage: (messageId: string) => void;
}

export const MessageList: React.FC<MessageListProps> = ({
  messages,
  isGenerating,
  onRegenerate,
  onEditMessage,
  onDeleteMessage,
}) => {
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // 自动滚动到底部
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages, isGenerating]);

  return (
    <div className={styles.messageList}>
      {messages.map((message, index) => (
        <div
          key={message.id}
          className={`${styles.message} ${styles[message.role]}`}
        >
          <Avatar size="small">
            {message.role === 'user' ? 'U' : 'AI'}
          </Avatar>

          <div className={styles.content}>
            <div className={styles.header}>
              <span className={styles.role}>
                {message.role === 'user' ? '用户' : 'AI助手'}
              </span>
              <span className={styles.time}>
                {new Date(message.createdAt).toLocaleTimeString()}
              </span>
            </div>

            <ReactMarkdown
              className={styles.markdown}
              components={{
                code({ node, inline, className, children, ...props }) {
                  const match = /language-(\w+)/.exec(className || '');
                  return !inline && match ? (
                    <SyntaxHighlighter
                      language={match[1]}
                      PreTag="div"
                      {...props}
                    >
                      {String(children).replace(/\n$/, '')}
                    </SyntaxHighlighter>
                  ) : (
                    <code className={className} {...props}>
                      {children}
                    </code>
                  );
                },
              }}
            >
              {message.content}
            </ReactMarkdown>

            {/* 引用来源 */}
            {message.metadata?.citations && (
              <div className={styles.citations}>
                <span>引用来源：</span>
                {message.metadata.citations.map((citation, i) => (
                  <a key={i} href="#" onClick={() => showCitation(citation)}>
                    [{i + 1}]
                  </a>
                ))}
              </div>
            )}

            {/* 操作按钮 */}
            {message.role === 'assistant' && message.status === 'sent' && (
              <div className={styles.actions}>
                <Button
                  icon={<IconThumbsUp />}
                  size="small"
                  theme="borderless"
                  onClick={() => handleRateMessage(message.id, 'thumbs_up')}
                />
                <Button
                  icon={<IconThumbsDown />}
                  size="small"
                  theme="borderless"
                  onClick={() => handleRateMessage(message.id, 'thumbs_down')}
                />
                <Button
                  icon={<IconRefresh />}
                  size="small"
                  theme="borderless"
                  onClick={() => onRegenerate(message.id)}
                />
                <Dropdown
                  trigger="click"
                  render={
                    <Dropdown.Menu>
                      <Dropdown.Item onClick={() => onEditMessage(message.id, message.content)}>
                        编辑
                      </Dropdown.Item>
                      <Dropdown.Item onClick={() => copyMessage(message.content)}>
                        复制
                      </Dropdown.Item>
                      <Dropdown.Item onClick={() => onDeleteMessage(message.id)}>
                        删除
                      </Dropdown.Item>
                    </Dropdown.Menu>
                  }
                >
                  <Button icon={<IconMore />} size="small" theme="borderless" />
                </Dropdown>
              </div>
            )}

            {/* 正在生成中... */}
            {message.status === 'generating' && (
              <div className={styles.generating}>
                <span>AI正在思考</span>
                <span className={styles.dots}>...</span>
              </div>
            )}
          </div>
        </div>
      ))}

      <div ref={messagesEndRef} />
    </div>
  );
};
```

#### 6.1.3 InputArea（输入区域）

```tsx
import React, { useState, useRef, useEffect } from 'react';
import { Button, Input, Upload } from '@douyinfe/semi-ui';
import { IconSend, IconStop, IconPaperclip, IconMicrophone } from '@douyinfe/semi-icons';
import { TextArea } from '@/components/TextArea';
import styles from './InputArea.module.scss';

interface InputAreaProps {
  onSendMessage: (content: string, attachments?: File[]) => void;
  onStopGeneration: () => void;
  isGenerating: boolean;
}

export const InputArea: React.FC<InputAreaProps> = ({
  onSendMessage,
  onStopGeneration,
  isGenerating,
}) => {
  const [content, setContent] = useState('');
  const [attachments, setAttachments] = useState<File[]>([]);
  const inputRef = useRef<any>();

  const handleSend = () => {
    if (content.trim() || attachments.length > 0) {
      onSendMessage(content, attachments);
      setContent('');
      setAttachments([]);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    // Enter发送，Shift+Enter换行
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <div className={styles.inputArea}>
      {/* 附件预览 */}
      {attachments.length > 0 && (
        <div className={styles.attachments}>
          {attachments.map((file, index) => (
            <div key={index} className={styles.attachment}>
              <span>{file.name}</span>
              <button
                onClick={() => setAttachments(attachments.filter((_, i) => i !== index))}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}

      {/* 输入框 */}
      <div className={styles.inputWrapper}>
        <Upload
          customRequest={({ file }) => setAttachments([...attachments, file])}
          showUploadList={false}
        >
          <Button icon={<IconPaperclip />} size="small" theme="borderless" />
        </Upload>

        <TextArea
          ref={inputRef}
          value={content}
          onChange={setContent}
          onKeyDown={handleKeyDown}
          placeholder="输入消息... (Enter发送，Shift+Enter换行)"
          autoSize={{ minRows: 1, maxRows: 10 }}
          className={styles.textarea}
        />

        {isGenerating ? (
          <Button
            icon={<IconStop />}
            size="small"
            type="danger"
            onClick={onStopGeneration}
          >
            停止
          </Button>
        ) : (
          <Button
            icon={<IconSend />}
            size="small"
            type="primary"
            onClick={handleSend}
            disabled={!content.trim() && attachments.length === 0}
          >
            发送
          </Button>
        )}
      </div>

      {/* 语音输入（可选） */}
      <div className={styles.footer}>
        <Button
          icon={<IconMicrophone />}
          size="small"
          theme="borderless"
          onClick={handleVoiceInput}
        >
          语音输入
        </Button>
      </div>
    </div>
  );
};
```

#### 6.1.4 QuickCommands（快捷指令）

```tsx
import React, { useState, useEffect } from 'react';
import { Card, Grid } from '@douyinfe/semi-ui';
import { chatAPI } from '@/services/api';
import styles from './QuickCommands.module.scss';

export const QuickCommands: React.FC<{ botId?: string; onSelectCommand: (prompt: string) => void }> = ({
  botId,
  onSelectCommand,
}) => {
  const [commands, setCommands] = useState<QuickCommand[]>([]);

  useEffect(() => {
    loadQuickCommands();
  }, [botId]);

  const loadQuickCommands = async () => {
    try {
      const resp = await chatAPI.listQuickCommands({ botId });
      setCommands(resp.data.items);
    } catch (error) {
      console.error('Failed to load quick commands:', error);
    }
  };

  return (
    <div className={styles.quickCommands}>
      <h4>快捷指令</h4>
      <Grid gutter={8} row={12} className={styles.grid}>
        {commands.map(cmd => (
          <Grid.Col span={6} key={cmd.id}>
            <Card
              className={styles.card}
              hoverable
              onClick={() => onSelectCommand(cmd.prompt)}
            >
              <div className={styles.icon}>{cmd.icon}</div>
              <div className={styles.title}>{cmd.title}</div>
              <div className={styles.description}>{cmd.description}</div>
            </Card>
          </Grid.Col>
        ))}
      </Grid>
    </div>
  );
};
```

### 6.2 样式设计

**SCSS 示例**（Chatbox.module.scss）：

```scss
.chatbox {
  display: flex;
  height: 100vh;
  background: #f5f5f5;

  .main {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
}

.messageList {
  flex: 1;
  overflow-y: auto;
  padding: 20px;

  .message {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;

    &.user {
      flex-direction: row-reverse;

      .content {
        background: #1890ff;
        color: white;
        border-radius: 12px 12px 0 12px;
      }
    }

    &.assistant {
      .content {
        background: white;
        border: 1px solid #e5e7eb;
        border-radius: 12px 12px 12px 0;
      }
    }
  }

  .content {
    max-width: 70%;
    padding: 12px 16px;

    .header {
      margin-bottom: 8px;
      font-size: 12px;
      color: #999;
    }

    .markdown {
      line-height: 1.6;

      pre {
        background: #f5f5f5;
        padding: 12px;
        border-radius: 6px;
        overflow-x: auto;
      }

      code {
        background: #f5f5f5;
        padding: 2px 6px;
        border-radius: 3px;
        font-family: 'Monaco', 'Consolas', monospace;
      }
    }

    .actions {
      margin-top: 8px;
      display: flex;
      gap: 4px;
      opacity: 0;
      transition: opacity 0.2s;

      &:hover {
        opacity: 1;
      }
    }

    &:hover .actions {
      opacity: 1;
    }
  }
}

.inputArea {
  padding: 16px 20px;
  background: white;
  border-top: 1px solid #e5e7eb;

  .inputWrapper {
    display: flex;
    gap: 12px;
    align-items: flex-end;
  }

  .textarea {
    flex: 1;
  }
}
```

---

## 7. 配置系统设计

### 7.1 快捷指令配置

**预置快捷指令**（`tenant_id IS NULL`）：

```sql
INSERT INTO quick_commands (tenant_id, bot_id, category, title, description, prompt, icon, sort_order) VALUES
-- 写作类
(NULL, NULL, 'writing', '写邮件', '快速撰写商务邮件', '请帮我写一封商务邮件，收件人：{{收件人}}，主题：{{主题}}，内容要点：\n{{内容}}', '📧', 1),
(NULL, NULL, 'writing', '写周报', '快速撰写工作周报', '请帮我写一份周报，本周完成的工作：\n{{工作内容}}', '📝', 2),
(NULL, NULL, 'writing', '总结文档', '总结文档要点', '请帮我总结以下文档的要点：\n{{文档内容}}', '📄', 3),

-- 编程类
(NULL, NULL, 'coding', '写代码', '生成代码片段', '请帮我写一段代码，语言：{{语言}}，功能：{{功能描述}}', '💻', 10),
(NULL, NULL, 'coding', '调试代码', '分析代码错误', '请帮我分析以下代码中的错误：\n{{代码}}', '🐛', 11),
(NULL, NULL, 'coding', '优化代码', '优化代码性能', '请帮我优化以下代码的性能：\n{{代码}}', '⚡', 12),

-- 分析类
(NULL, NULL, 'analysis', '数据分析', '分析数据趋势', '请帮我分析以下数据的趋势：\n{{数据}}', '📊', 20),
(NULL, NULL, 'analysis', '竞品分析', '分析竞品优劣势', '请帮我分析{{产品A}}相对于{{产品B}}的优劣势', '🔍', 21);
```

### 7.2 提示词建议配置

```sql
INSERT INTO prompt_suggestions (tenant_id, bot_id, suggestion_type, title, prompt, trigger_keywords) VALUES
(NULL, NULL, 'contextual', '继续生成', '请继续上述内容的生成...', '["继续", "还有呢", "go on"]'),
(NULL, NULL, 'contextual', '换个说法', '请用不同的方式重新表达上述内容', '["换个说法", "重新表达", "another way"]'),
(NULL, NULL, 'template', '翻译成英文', '请将以下内容翻译成英文：', '["翻译", "translate", "英文"]'),
(NULL, NULL, 'template', '简化解释', '请用更简单易懂的方式解释：', '["简化", "通俗", "小白能懂"]');
```

---

## 8. 核心代码实现

### 8.1 会话服务

```go
package chat

// ConversationService 会话服务
type ConversationService struct {
    conversationRepo repository.ConversationRepository
    messageRepo      repository.MessageRepository
    botExecutor      *BotExecutor
    logger           *zap.Logger
}

// CreateConversation 创建会话
func (s *ConversationService) CreateConversation(
    ctx context.Context,
    req *CreateConversationRequest,
) (*Conversation, error) {
    conversation := &model.Conversation{
        ID:       generateConversationID(),
        TenantID: getTenantID(ctx),
        UserID:   getUserID(ctx),
        BotID:    req.BotID,
        Title:    req.Title,
        Status:   "active",
    }

    if err := s.conversationRepo.Create(ctx, conversation); err != nil {
        return nil, err
    }

    return conversation, nil
}

// SendMessage 发送消息（HTTP接口，用于非WebSocket场景）
func (s *ConversationService) SendMessage(
    ctx context.Context,
    conversationID string,
    content string,
) (*Message, error) {
    // 1. 创建用户消息
    userMsg := &model.Message{
        ID:             generateMessageID(),
        TenantID:       getTenantID(ctx),
        ConversationID: conversationID,
        UserID:         getUserID(ctx),
        Role:           "user",
        Content:        content,
        Status:         "sent",
    }

    if err := s.messageRepo.Create(ctx, userMsg); err != nil {
        return nil, err
    }

    // 2. 加载上下文
    messages, _ := s.messageRepo.ListByConversationID(ctx, conversationID, 10)

    // 3. 执行Bot
    stream := s.botExecutor.ExecuteStream(ctx, &BotExecuteRequest{
        ConversationID: conversationID,
        Messages:       messages,
    })

    // 4. 收集完整响应
    var contentBuilder strings.Builder
    for chunk := range stream {
        if chunk.Error != nil {
            return nil, chunk.Error
        }
        contentBuilder.WriteString(chunk.Delta)
    }

    // 5. 保存助手消息
    assistantMsg := &model.Message{
        ID:             generateMessageID(),
        TenantID:       getTenantID(ctx),
        ConversationID: conversationID,
        UserID:         getUserID(ctx),
        Role:           "assistant",
        Content:        contentBuilder.String(),
        Status:         "sent",
    }

    if err := s.messageRepo.Create(ctx, assistantMsg); err != nil {
        return nil, err
    }

    return assistantMsg, nil
}
```

---

## 9. 测试用例

### 9.1 WebSocket连接测试

```javascript
// 使用 jest-websocket-mock 测试
import WS from 'jest-websocket-mock';

test('WebSocket连接和消息发送', async () => {
  const ws = new WS('ws://localhost:8888/v1/chat/ws?token=xxx');

  // 等待连接建立
  await ws.connected;

  // 发送消息
  ws.send(JSON.stringify({
    type: 'message.send',
    data: { content: '你好', messageId: 'msg-123' },
  }));

  // 等待响应
  await expect(ws).toReceiveMessage(JSON.stringify({
    type: 'message.received',
    data: { messageId: 'msg-123', tempId: expect.any(String) },
  }));

  await expect(ws).toReceiveMessage(JSON.stringify({
    type: 'generation.start',
    data: { messageId: expect.any(String) },
  }));
});
```

---

## 10. 部署方案

### 10.1 环境变量配置

```bash
# WebSocket
WS_ENABLED=true
WS_READ_BUFFER_SIZE=4096
WS_WRITE_BUFFER_SIZE=4096
WS_MAX_MESSAGE_SIZE=1048576

# 消息限制
MAX_MESSAGE_LENGTH=10000
MAX_CONTEXT_MESSAGES=20
MAX_CONCURRENT_SESSIONS=100

# 文件上传
MAX_UPLOAD_SIZE=52428800
ALLOWED_FILE_TYPES=.pdf,.doc,.docx,.txt,.md,.png,.jpg,.jpeg
```

---

## 11. 总结

### 11.1 实施策略总结

| 实施项 | 实施方式 | 工作量 | 说明 |
|-------|---------|--------|------|
| **前端UI组件** | 💻 独立编码 | 60% | React组件 + WebSocket客户端 |
| **WebSocket服务** | 💻 独立编码 | 20% | 实时双向通信 |
| **会话/消息存储** | 📊 数据库设计 | 10% | conversations + messages 表 |
| **快捷指令** | 📊 数据库配置 | 5% | quick_commands 表 |
| **提示词建议** | 📊 数据库配置 | 5% | prompt_suggestions 表 |

**总计**：60% 编码（UI + WebSocket） + 40% 配置（数据库）

### 11.2 与鲸智百应对齐情况

| 功能模块 | 对齐情况 | 实现方式 |
|---------|---------|----------|
| 多轮对话 | ✅ 100% | WebSocket + 会话管理 |
| 流式输出 | ✅ 100% | SSE + 打字机效果 |
| 快捷指令 | ✅ 100% | 数据库配置 |
| 对话历史 | ✅ 100% | messages 表存储 |
| 会话分享 | ✅ 100% | shared_conversations 表 |
| 多模态输入 | ✅ 100% | 附件上传 + 解析 |
| 引用来源 | ✅ 100% | RAG 引用展示 |

### 11.3 核心优势

- ✅ **实时交互**：WebSocket 双向通信，延迟低
- ✅ **流式输出**：打字机效果，用户体验好
- ✅ **多模态**：支持文本、图片、文件等多种输入
- ✅ **可配置**：快捷指令、提示词建议数据库配置
- ✅ **企业级**：会话隔离、权限控制、审计日志

---

**文档结束**
