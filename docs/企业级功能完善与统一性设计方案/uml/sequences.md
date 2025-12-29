# ZKER Enterprise - UML序列图集合

**文档版本**: v1.0.0
**创建日期**: 2025-12-29
**作者**: ZKER Enterprise Team

---

## 📋 文档说明

本文档包含ZKER Enterprise平台关键业务流程的UML序列图，使用Mermaid格式绘制。

---

## 目录

1. [用户认证流程](#1-用户认证流程)
2. [AI对话流程](#2-ai对话流程)
3. [Bot创建与配置流程](#3-bot创建与配置流程)
4. [知识库文档上传与向量化流程](#4-知识库文档上传与向量化流程)
5. [RAG检索流程](#5-rag检索流程)
6. [租户注册流程](#6-租户注册流程)
7. [Token计量与预算告警流程](#7-token计量与预算告警流程)
8. [权限检查流程](#8-权限检查流程)

---

## 1. 用户认证流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端应用
    participant API as API网关
    participant AuthService as 认证服务
    participant UserService as 用户服务
    participant TenantService as 租户服务
    participant Redis as Redis缓存

    User->>Frontend: 输入用户名密码
    Frontend->>API: POST /auth/login
    Note over API: 提取租户ID (subdomain/header)

    API->>AuthService: 验证用户凭证
    AuthService->>UserService: 查询用户信息
    UserService->>TenantService: 验证租户状态

    alt 租户存在且状态正常
        TenantService-->>UserService: 租户信息
        UserService-->>AuthService: 用户数据

        alt 密码验证成功
            AuthService->>Redis: 生成JWT Token
            Note over Redis: access_token (2h)<br/>refresh_token (7d)
            Redis-->>AuthService: Token生成成功

            AuthService-->>API: 返回Token + 用户信息
            API-->>Frontend: 返回登录成功

            Frontend->>Frontend: 保存Token到localStorage
            Frontend->>Frontend: 跳转到首页
        else 密码错误
            AuthService-->>API: 返回认证失败
            API-->>Frontend: 401 Unauthorized
            Frontend-->>User: 显示密码错误提示
        end
    else 租户不存在或已禁用
        TenantService-->>AuthService: 租户无效
        AuthService-->>API: 返回租户错误
        API-->>Frontend: 403 Forbidden
        Frontend-->>User: 显示租户已禁用提示
    end
```

---

## 2. AI对话流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端应用
    participant API as API网关
    participant ConversationService as 对话服务
    participant BotService as Bot服务
    participant LLMService as LLM服务
    participant MemoryService as 记忆服务
    participant KnowledgeService as 知识服务
    participant PluginService as 插件服务

    User->>Frontend: 发送消息
    Frontend->>API: POST /conversations/{id}/messages
    Note over API: Authorization: Bearer {token}<br/>X-Tenant-ID: {tenant_id}

    API->>ConversationService: 创建消息记录
    ConversationService->>BotService: 获取Bot配置

    BotService->>MemoryService: 加载对话历史
    MemoryService-->>BotService: 返回历史消息

    BotService->>BotService: 构建Prompt

    alt 启用RAG
        BotService->>KnowledgeService: 检索相关知识
        KnowledgeService->>KnowledgeService: 混合检索 (向量+关键词)
        KnowledgeService-->>BotService: 返回相关知识块
        BotService->>BotService: 增强Prompt
    end

    alt 需要调用插件
        BotService->>PluginService: 调用插件
        PluginService-->>BotService: 返回插件结果
        BotService->>BotService: 更新Prompt
    end

    BotService->>LLMService: 发送请求到LLM
    Note over LLMService: model: gpt-4<br/>temperature: 0.7<br/>max_tokens: 2000

    LLMService-->>BotService: 返回AI响应
    Note over BotService: stream: true/false

    BotService->>ConversationService: 保存AI响应

    alt 流式响应
        API-->>Frontend: SSE流式返回
        Frontend->>User: 实时显示响应
    else 一次性响应
        API-->>Frontend: 返回完整响应
        Frontend->>User: 显示响应
    end

    ConversationService->>ConversationService: 记录Token使用量
    ConversationService->>ConversationService: 更新Bot指标
```

---

## 3. Bot创建与配置流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端应用
    participant API as API网关
    participant BotService as Bot服务
    participant KnowledgeService as 知识库服务
    participant PluginService as 插件服务
    participant DB as 数据库

    User->>Frontend: 创建新Bot
    Frontend->>API: POST /bots
    Note over API: Request Body:<br/>{name, description, type, ...}

    API->>BotService: 创建Bot
    BotService->>DB: INSERT INTO bots
    DB-->>BotService: 返回bot_id

    BotService-->>API: 返回Bot基本信息
    API-->>Frontend: Bot创建成功

    User->>Frontend: 配置Bot参数
    Frontend->>API: PUT /bots/{id}

    alt 配置提示词
        API->>BotService: 更新提示词模板
        BotService->>BotService: 验证提示词格式
        BotService->>DB: UPDATE bots SET prompt_template
    end

    alt 配置模型参数
        API->>BotService: 更新模型配置
        BotService->>DB: UPDATE bots SET model_config
        Note over DB: {model, temperature, max_tokens}
    end

    alt 关联知识库
        API->>KnowledgeService: 关联知识库
        KnowledgeService->>DB: INSERT INTO bot_knowledge_bases
        Note over DB: bot_id, knowledge_base_id
    end

    alt 安装插件
        API->>PluginService: 安装插件
        PluginService->>DB: INSERT INTO bot_plugins
        Note over DB: bot_id, plugin_id, config
    end

    BotService-->>API: 配置更新成功
    API-->>Frontend: 返回完整Bot配置
    Frontend->>User: 显示配置成功
```

---

## 4. 知识库文档上传与向量化流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端应用
    participant API as API网关
    participant KnowledgeService as 知识库服务
    participant Storage as 对象存储
    participant Worker as 向量化Worker
    participant EmbeddingService as 向量化服务
    participant VectorDB as 向量数据库(Milvus)
    participant DB as 数据库(MySQL)

    User->>Frontend: 上传文档
    Frontend->>API: POST /knowledge/{id}/documents
    Note over API: Content-Type: multipart/form-data

    API->>Storage: 上传文件到MinIO/S3
    Storage-->>API: 返回file_url

    API->>KnowledgeService: 创建文档记录
    KnowledgeService->>DB: INSERT INTO knowledge_documents
    Note over DB: status = 'pending'<br/>file_url, file_size, file_type

    KnowledgeService-->>API: 返回document_id
    API-->>Frontend: 文档上传成功，开始处理

    Note over Worker: 异步Worker轮询待处理文档
    Worker->>DB: SELECT * FROM knowledge_documents WHERE status = 'pending'

    Worker->>KnowledgeService: 开始处理文档
    KnowledgeService->>DB: UPDATE status = 'processing'

    KnowledgeService->>Storage: 下载文档内容
    Storage-->>KnowledgeService: 返回文档内容

    alt 分块策略 = hybrid
        KnowledgeService->>KnowledgeService: 智能分块
        Note over KnowledgeService: 1. 语义分块<br/>2. 递归分块<br/>3. 固定大小兜底
    end

    KnowledgeService->>DB: 批量插入分块记录
    Note over DB: INSERT INTO knowledge_chunks<br/>(chunk_text, chunk_type)

    loop 每个分块
        KnowledgeService->>EmbeddingService: 生成向量
        EmbeddingService->>EmbeddingService: 调用OpenAI API
        EmbeddingService-->>KnowledgeService: 返回1536维向量

        KnowledgeService->>VectorDB: 存储向量
        VectorDB->>VectorDB: 创建HNSW索引
        Note over VectorDB: Collection: knowledge_chunks
    end

    KnowledgeService->>DB: UPDATE document SET status = 'completed'
    KnowledgeService->>DB: UPDATE knowledge_base SET chunk_count = chunk_count + N

    API-->>Frontend: WebSocket通知处理完成
    Frontend->>User: 显示"文档已成功向量化"
```

---

## 5. RAG检索流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant BotService as Bot服务
    participant KnowledgeService as 知识库服务
    participant VectorDB as 向量数据库
    participant KeywordSearch as 关键词搜索引擎
    participant RerankerService as 重排序服务
    participant LLMService as LLM服务

    User->>BotService: 发送查询
    BotService->>BotService: 查询重写
    Note over BotService: 纠错、扩展、HyDE

    BotService->>KnowledgeService: 检索知识
    KnowledgeService->>KnowledgeService: 查询向量化
    KnowledgeService->>VectorDB: 向量检索
    Note over VectorDB: Top-K: 10<br/>Metric: Cosine Similarity
    VectorDB-->>KnowledgeService: 返回Top-K结果

    par 并行检索
        KnowledgeService->>KeywordSearch: 关键词检索
        Note over KeywordSearch: BM25算法<br/>Top-K: 10
    and
        KnowledgeService->>KnowledgeService: 精确匹配
    end

    KeywordSearch-->>KnowledgeService: 返回关键词结果

    KnowledgeService->>KnowledgeService: RRF融合
    Note over KnowledgeService: Reciprocal Rank Fusion<br/>score = 1/(k+rank_vector) + 1/(k+rank_keyword)<br/>k = 60

    alt 启用Reranker
        KnowledgeService->>RerankerService: 重排序
        Note over RerankerService: Cohere Rerank API<br/>或 Cross-Encoder
        RerankerService-->>KnowledgeService: 返回重排序结果
    end

    KnowledgeService-->>BotService: 返回Top-3知识块

    BotService->>BotService: 构建增强Prompt
    Note over BotService: System Prompt<br/>+ 知识块<br/>+ 历史对话<br/>+ 用户查询

    BotService->>LLMService: 发送到LLM
    LLMService-->>BotService: 返回响应

    BotService->>BotService: 引用标注
    Note over BotService: 高亮引用的知识来源<br/>添加置信度分数

    BotService-->>User: 返回带引用的响应
```

---

## 6. 租户注册流程

```mermaid
sequenceDiagram
    participant Admin as 超级管理员
    participant Frontend as 前端应用
    participant API as API网关
    participant TenantService as 租户服务
    participant UserService as 用户服务
    participant DB as 数据库
    participant EmailService as 邮件服务

    Admin->>Frontend: 创建新租户
    Frontend->>API: POST /tenants
    Note over API: {name, domain, admin_email, plan}

    API->>TenantService: 验证域名可用性
    TenantService->>DB: SELECT * FROM tenants WHERE domain = ?

    alt 域名已存在
        TenantService-->>API: 返回域名冲突
        API-->>Frontend: 409 Conflict
        Frontend-->>Admin: 显示"域名已被使用"
    else 域名可用
        API->>TenantService: 创建租户
        TenantService->>DB: INSERT INTO tenants
        Note over DB: id = tenant_{uuid}<br/>status = 'active'<br/>plan = 'basic'

        TenantService->>DB: 创建默认部门
        Note over DB: INSERT INTO departments (tenant_id, name, parent_id)

        TenantService->>UserService: 创建管理员账号
        UserService->>DB: INSERT INTO users
        Note over DB: username=admin, role=admin

        TenantService->>EmailService: 发送欢迎邮件
        Note over EmailService: 包含登录链接<br/>临时密码

        TenantService-->>API: 返回租户信息
        API-->>Frontend: 租户创建成功
        Frontend->>Admin: 显示租户详情和管理链接
    end
```

---

## 7. Token计量与预算告警流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant BotService as Bot服务
    participant LLMService as LLM服务
    participant TokenMetering as Token计量服务
    participant BudgetAlert as 预算告警服务
    participant DB as 数据库
    participant NotificationService as 通知服务

    User->>BotService: 发起对话

    BotService->>TokenMetering: 检查预算状态
    TokenMetering->>DB: SELECT * FROM budget_alerts WHERE bot_id = ?
    DB-->>TokenMetering: 返回预算配置

    alt 预算已耗尽
        TokenMetering-->>BotService: 预算不足
        BotService->>BotService: 启用降级模型
        Note over BotService: model: gpt-4 → gpt-3.5-turbo

        BotService->>NotificationService: 发送告警
        NotificationService-->>User: "预算已耗尽，已切换到低成本模型"
    end

    BotService->>LLMService: 调用LLM API
    LLMService-->>BotService: 返回响应 + Token使用量
    Note over BotService: input_tokens: 100<br/>output_tokens: 500

    BotService->>TokenMetering: 记录Token使用
    TokenMetering->>DB: INSERT INTO token_usage_metrics
    Note over DB: tenant_id, bot_id, input_tokens, output_tokens<br/>total_cost = (input + output) * unit_price

    TokenMetering->>BudgetAlert: 检查预算使用率
    BudgetAlert->>DB: SELECT SUM(total_tokens) FROM token_usage_metrics

    alt 使用率 >= 100%
        BudgetAlert->>NotificationService: 发送紧急告警
        NotificationService-->>User: "预算已耗尽！"
        NotificationService-->>Admin: "用户 {user} 预算已耗尽"

    else 使用率 >= 90%
        BudgetAlert->>NotificationService: 发送警告
        NotificationService-->>User: "预算使用率已达90%"

    else 使用率 >= 80%
        BudgetAlert->>NotificationService: 发送提醒
        NotificationService-->>User: "预算使用率已达80%"
    end

    BotService-->>User: 返回响应
```

---

## 8. 权限检查流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant Frontend as 前端应用
    participant API as API网关
    participant AuthMiddleware as 认证中间件
    participant PermissionService as 权限服务
    participant DB as 数据库
    participant Cache as Redis缓存

    User->>Frontend: 请求访问资源
    Frontend->>API: PUT /bots/{bot_id}
    Note over API: Authorization: Bearer {token}<br/>X-Tenant-ID: {tenant_id}

    API->>AuthMiddleware: 验证Token
    AuthMiddleware->>Cache: GET token:{token}

    alt Token有效且在缓存中
        Cache-->>AuthMiddleware: 返回user_info
    else Token不在缓存
        AuthMiddleware->>DB: SELECT * FROM users WHERE id = ?
        DB-->>AuthMiddleware: 返回用户信息
        AuthMiddleware->>Cache: SET token:{token} user_info (TTL: 2h)
    end

    alt Token无效或过期
        AuthMiddleware-->>API: 401 Unauthorized
        API-->>Frontend: 请重新登录
        Frontend->>User: 跳转到登录页
    end

    API->>PermissionService: 检查权限
    Note over PermissionService: resource: bot<br/>action: edit

    PermissionService->>Cache: GET permission:{user_id}:bot:edit

    alt 权限在缓存中
        Cache-->>PermissionService: 返回权限结果
    else 权限不在缓存
        PermissionService->>DB: 查询用户角色
        DB-->>PermissionService: 返回role_id

        PermissionService->>DB: 查询角色权限
        DB-->>PermissionService: 返回permissions列表

        PermissionService->>PermissionService: 检查是否有编辑权限

        alt 拥有权限
            PermissionService->>Cache: SET permission:{user_id}:bot:edit true (TTL: 1h)
            PermissionService-->>API: 权限检查通过
        else 无权限
            PermissionService-->>API: 权限不足
            API-->>Frontend: 403 Forbidden
            Frontend-->>User: 显示"权限不足"
        end
    end

    alt 检查数据权限
        PermissionService->>DB: SELECT * FROM data_permissions WHERE user_id = ?
        DB-->>PermissionService: 返回数据权限配置

        alt permission_level = 'department'
            PermissionService->>DB: 检查用户部门与Bot创建者部门
            alt 同部门
                PermissionService-->>API: 权限检查通过
            else 不同部门
                PermissionService-->>API: 403 Forbidden
            end
        end
    end

    API->>API: 执行业务逻辑
    API-->>Frontend: 返回操作结果
```

---

## 9. Bot发布到多渠道流程

```mermaid
sequenceDiagram
    participant User as 用户
    participant BotService as Bot服务
    participant ChannelService as 渠道服务
    participant WechatAdapter as 微信适配器
    participant FeishuAdapter as 飞书适配器
    participant DiscordAdapter as Discord适配器

    User->>BotService: 发布Bot到渠道
    Note over BotService: channels: [wechat, feishu, discord]

    BotService->>ChannelService: 创建发布任务

    par 发布到微信公众号
        ChannelService->>WechatAdapter: 配置公众号
        WechatAdapter->>WechatAdapter: 创建菜单
        WechatAdapter->>WechatAdapter: 配置自动回复
        WechatAdapter->>WechatAdapter: 设置Webhook
        WechatAdapter-->>ChannelService: 发布成功
    end

    par 发布到飞书
        ChannelService->>FeishuAdapter: 创建飞书应用
        FeishuAdapter->>FeishuAdapter: 配置事件订阅
        FeishuAdapter->>FeishuAdapter: 配置机器人卡片
        FeishuAdapter-->>ChannelService: 发布成功
    end

    par 发布到Discord
        ChannelService->>DiscordAdapter: 创建Discord Bot
        DiscordAdapter->>DiscordAdapter: 注册命令
        DiscordAdapter->>DiscordAdapter: 配置权限
        DiscordAdapter-->>ChannelService: 发布成功
    end

    ChannelService-->>BotService: 所有渠道发布完成
    BotService-->>User: 显示发布结果

    Note over User: 微信: ✅<br/>飞书: ✅<br/>Discord: ❌ (配置错误)
```

---

## 📊 序列图使用说明

### 如何在文档中使用

在Markdown文档中引用这些序列图：

```markdown
## 用户登录流程

详见 [UML序列图 - 用户认证流程](./uml-sequences.md#1-用户认证流程)
```

### 在Mermaid Live Editor中查看

1. 访问 https://mermaid.live/
2. 将序列图代码复制粘贴到编辑器
3. 实时查看渲染效果

### 在支持Mermaid的Markdown查看器中查看

- GitHub: 原生支持
- GitLab: 原生支持
- VS Code: 安装 Mermaid Preview 插件
- Obsidian: 原生支持

---

## 🎯 扩展建议

如需添加更多业务流程的序列图，建议按以下模板创建：

```mermaid
sequenceDiagram
    participant Actor1 as 角色名称1
    participant Actor2 as 角色名称2
    participant System as 系统名称

    Actor1->>Actor2: 操作描述
    Actor2->>System: 系统调用
    System-->>Actor2: 返回结果
    Actor2-->>Actor1: 响应结果
```

---

**文档维护**: ZKER Enterprise Team
**最后更新**: 2025-12-29
