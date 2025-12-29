# ZKER核心业务流程UML时序图

> **文档编号**: ARCH-UML-2025-001
> **文档类型**: 架构设计文档
> **创建日期**: 2025-12-30
> **版本**: v1.0
> **密级**: 内部公开

---

## 📋 目录

- [1. 租户注册流程](#1-租户注册流程)
- [2. Bot创建流程](#2-bot创建流程)
- [3. 对话交互流程](#3-对话交互流程)
- [4. 知识库检索流程](#4-知识库检索流程)
- [5. 权限验证流程](#5-权限验证流程)
- [6. 补充指南](#6-补充指南)

---

## 1. 租户注册流程

### 1.1 流程概述

**业务场景**：企业用户自助注册ZKER平台

**涉及服务**：
- API网关
- 租户服务 (tenant-service)
- 组织服务 (organization-service)
- 用户服务 (user-service)
- 通知服务 (notification-service)
- MySQL数据库
- Redis缓存

### 1.2 时序图

```mermaid
sequenceDiagram
    autonumber
    participant User as 用户
    participant API as API网关
    participant Tenant as 租户服务
    participant Org as 组织服务
    participant UserSvc as 用户服务
    participant Notification as 通知服务
    participant MySQL as MySQL
    participant Redis as Redis

    User->>API: POST /api/v1/tenants/register
    Note over User,API: 提交注册信息<br/>（企业名称、联系人、子域名）

    API->>API: 租户识别中间件<br/>（跳过，注册时无租户）
    API->>Tenant: RegisterTenantCommand

    Tenant->>MySQL: INSERT INTO tenants
    Note over Tenant,MySQL: 创建租户记录<br/>状态：trial（试用）

    MySQL-->>Tenant: 返回租户ID
    Tenant->>Redis: 缓存租户信息
    Note over Tenant,Redis: key: tenant:{tenant_id}<br/>ttl: 1小时

    Tenant->>Org: CreateDefaultOrganizationCommand
    Note over Tenant,Org: 租户ID + 企业名称

    Org->>MySQL: INSERT INTO organizations
    Note over Org,MySQL: 创建默认组织<br/>类型：company

    MySQL-->>Org: 返回组织ID
    Org->>Redis: 缓存组织信息
    Org-->>Tenant: 返回组织

    Tenant->>UserSvc: CreateAdminUserCommand
    Note over Tenant,UserSvc: 租户ID + 组织ID<br/>+ 联系人信息

    UserSvc->>MySQL: INSERT INTO users
    Note over UserSvc,MySQL: 创建管理员账户<br/>角色：admin

    MySQL-->>UserSvc: 返回用户ID
    UserSvc->>Redis: 缓存用户信息
    UserSvc->>MySQL: INSERT INTO user_roles
    Note over UserSvc,MySQL: 分配管理员角色

    UserSvc-->>Tenant: 返回管理员

    Tenant->>Notification: SendWelcomeEmailCommand
    Note over Tenant,Notification: 用户邮箱 + 租户信息

    Notification->>Notification: 生成欢迎邮件<br/>+ 登录链接
    Notification->>Notification: 发送邮件（SMTP）
    Notification-->>Tenant: 邮件发送成功

    Tenant-->>API: TenantDTO
    API-->>User: 200 OK<br/>{tenant_id, subdomain, admin_user_id}

    Note over User: 注册成功！<br/>14天试用期
```

### 1.3 状态机图

```mermaid
stateDiagram-v2
    [*] --> 注册中: 提交注册信息
    注册中 --> 试用中: 创建成功
    注册中 --> 注册失败: 创建失败
    注册失败 --> [*]

    试用中 --> 付费中: 购买套餐
    试用中 --> 已停用: 违规/欠费
    试用中 --> 已删除: 主动删除

    付费中 --> 正常: 支付成功
    付费中 --> 试用中: 支付失败

    正常 --> 已停用: 欠费/违规
    正常 --> 已删除: 主动删除

    已停用 --> 正常: 补缴费用/解除违规
    已停用 --> 已删除: 超过保留期

    已删除 --> [*]
```

### 1.4 关键节点说明

| 节点 | 说明 | 超时时间 |
|------|------|----------|
| **创建租户** | 在MySQL中创建租户记录 | 5秒 |
| **缓存租户** | 将租户信息写入Redis | 1秒 |
| **创建组织** | 创建默认组织（公司级别） | 5秒 |
| **创建管理员** | 创建管理员账户并分配角色 | 5秒 |
| **发送邮件** | 异步发送欢迎邮件 | 30秒（异步） |

### 1.5 异常处理

| 异常场景 | 处理策略 |
|---------|---------|
| **子域名冲突** | 返回400错误，提示用户更换子域名 |
| **MySQL连接失败** | 重试3次，间隔1秒，仍失败则返回500错误 |
| **邮件发送失败** | 不阻塞注册流程，后台重试 |
| **Redis失败** | 降级，直接从MySQL查询 |

---

## 2. Bot创建流程

### 2.1 流程概述

**业务场景**：用户在ZKER平台创建新的AI Bot

**涉及服务**：
- API网关
- Bot服务 (bot-service)
- 知识库服务 (knowledge-service)
- 权限服务 (permission-service)
- 审计服务 (audit-service)

### 2.2 时序图

```mermaid
sequenceDiagram
    autonumber
    participant User as 用户
    participant API as API网关
    participant Tenant as 租户识别中间件
    participant Bot as Bot服务
    participant KB as 知识库服务
    participant Perm as 权限服务
    participant Audit as 审计服务
    participant MySQL as MySQL

    User->>API: POST /api/v1/bots
    Note over User,API: Authorization: Bearer {token}

    API->>Tenant: 提取租户上下文
    Note over API,Tenant: 从Token解析tenant_id

    Tenant->>MySQL: SELECT * FROM tenants<br/>WHERE id = ?
    MySQL-->>Tenant: 租户信息
    Tenant->>API: 注入租户上下文
    Note over Tenant,API: ctx.Set("tenant_id", xxx)

    API->>Bot: CreateBotCommand<br/>(tenant_id, user_id, bot_data)

    Bot->>Bot: 验证Bot名称唯一性
    Bot->>MySQL: SELECT COUNT(*) FROM bots<br/>WHERE tenant_id = ? AND name = ?

    Bot->>MySQL: INSERT INTO bots
    Note over Bot,MySQL: 创建Bot记录<br/>状态：draft

    MySQL-->>Bot: 返回bot_id
    Bot->>KB: CreateDefaultKnowledgeBaseCommand
    Note over Bot,KB: bot_id + tenant_id

    KB->>MySQL: INSERT INTO knowledge_bases
    Note over KB,MySQL: 创建默认知识库<br/>类型：private

    MySQL-->>KB: 返回kb_id
    KB->>MySQL: INSERT INTO kb_chunks
    Note over KB,MySQL: 初始化空知识库

    KB-->>Bot: 返回知识库
    Bot->>Perm: GrantBotPermissionCommand
    Note over Bot,Perm: user_id, bot_id, role="owner"

    Perm->>MySQL: INSERT INTO permissions
    Note over Perm,MySQL: 授予Bot所有者权限

    MySQL-->>Perm: 返回permission_id
    Perm-->>Bot: 权限授予成功

    Bot->>Audit: LogBotCreationEvent
    Note over Bot,Audit: 记录审计日志

    Audit->>MySQL: INSERT INTO audit_logs
    Note over Audit,MySQL: 操作类型：create<br/>资源类型：bot

    Audit-->>Bot: 审计记录成功

    Bot-->>API: BotDTO
    API-->>User: 200 OK<br/>{bot_id, name, kb_id, status}
```

### 2.3 Bot状态机图

```mermaid
stateDiagram-v2
    [*] --> 草稿: 创建Bot
    草稿 --> 配置中: 配置提示词
    草稿 --> 已删除: 删除Bot

    配置中 --> 测试中: 保存并测试
    配置中 --> 草稿: 重置

    测试中 --> 已发布: 发布Bot
    测试中 --> 配置中: 修改配置

    已发布 --> 已停用: 停用Bot
    已发布 --> 已发布: 更新配置

    已停用 --> 已发布: 重新启用
    已停用 --> 已删除: 删除Bot

    已删除 --> [*]
```

---

## 3. 对话交互流程

### 3.1 流程概述

**业务场景**：用户与ZKER AI Bot进行对话

**涉及服务**：
- API网关
- 对话服务 (conversation-service)
- 智能路由引擎 (routing-engine)
- LLM服务 (llm-service)
- 插件调度引擎 (plugin-engine)
- 记忆引擎 (memory-engine)
- 知识库服务 (knowledge-service)

### 3.2 时序图

```mermaid
sequenceDiagram
    autonumber
    participant User as 用户
    participant API as API网关
    participant Conv as 对话服务
    participant Routing as 智能路由引擎
    participant Memory as 记忆引擎
    participant KB as 知识库服务
    participant LLM as LLM服务
    participant Plugin as 插件引擎
    participant MySQL as MySQL
    participant Redis as Redis

    User->>API: POST /api/v1/conversations/{conv_id}/messages
    Note over User,API: {message: "今天天气怎么样？"}

    API->>Conv: SendMessageCommand
    Note over API,Conv: conversation_id + user_message

    Conv->>MySQL: SELECT * FROM conversations<br/>WHERE id = ?
    MySQL-->>Conv: 对话信息

    Conv->>Routing: AnalyzeIntent
    Note over Conv,Routing: 分析用户意图

    Routing->>Routing: 提取关键词<br/>（天气）
    Routing-->>Conv: Intent: SMALLTALK<br/>Confidence: 0.95

    Conv->>Memory: LoadConversationHistory
    Note over Conv,Memory: 加载对话历史

    Memory->>Redis: GET conversation:{conv_id}
    Redis-->>Memory: 对话历史（最近10轮）
    Memory-->>Conv: 历史消息

    Conv->>KB: SearchRelevantKnowledge
    Note over Conv,KB: 查询相关知识

    KB->>KB: 向量检索
    Note over KB,KB: query: "今天天气"<br/>top_k: 3
    KB-->>Conv: 相关知识片段

    Conv->>LLM: GenerateCompletion
    Note over Conv,LLM: prompt + 历史记录<br/>+ 知识片段

    LLM->>LLM: 调用LLM API<br/>（OpenAI/Claude/通义千问）

    LLM->>LLM: Token计数
    Note over LLM: 统计输入/输出Token数

    LLM-->>Conv: LLMResponse<br/>{reply, tokens_used, model}

    Conv->>MySQL: INSERT INTO messages
    Note over Conv,MySQL: 保存用户消息

    Conv->>MySQL: INSERT INTO messages
    Note over Conv,MySQL: 保存Bot回复

    Conv->>MySQL: UPDATE conversations<br/>SET updated_at = NOW()
    Note over Conv,MySQL: 更新对话时间

    Conv->>Redis: SET conversation:{conv_id}
    Note over Conv,Redis: 缓存对话历史<br/>ttl: 24小时

    Conv-->>API: MessageDTO
    API-->>User: 200 OK<br/>{reply, message_id, tokens_used}

    Note over User: Bot回复：<br/>"对不起，我无法获取实时天气信息。"
```

### 3.3 对话状态机图

```mermaid
stateDiagram-v2
    [*] --> 进行中: 创建对话
    进行中 --> 暂停: 用户长时间无响应
    暂停 --> 进行中: 用户发送消息
    进行中 --> 已结束: 用户主动结束

    已结束 --> 已归档: 超过30天
    已归档 --> [*]

    note right of 进行中
        可发送消息
        可调用插件
        可使用知识库
    end note
```

---

## 4. 知识库检索流程

### 4.1 流程概述

**业务场景**：RAG检索增强生成

**涉及服务**：
- 对话服务
- 知识库服务
- 向量数据库 (Milvus)
- Elasticsearch
- LLM服务

### 4.2 时序图

```mermaid
sequenceDiagram
    autonumber
    participant User as 用户消息
    participant KB as 知识库服务
    participant Milvus as Milvus向量库
    participant ES as Elasticsearch
    participant Reranker as 重排序引擎
    participant LLM as LLM服务

    User->>KB: 查询知识<br/>"如何配置租户权限？"

    KB->>KB: Query预处理
    Note over KB: 分词、去停用词<br/>提取关键词

    KB->>ES: 关键词检索
    Note over KB,ES: query: "租户 权限 配置"<br/>top_k: 50

    ES-->>KB: BM25文档列表<br/>{doc_id, score}

    KB->>KB: 生成查询向量
    Note over KB: 使用Embedding模型<br/>text-embedding-3-small

    KB->>Milvus: 向量检索
    Note over KB,Milvus: vector: [0.1, 0.2, ...]<br/>top_k: 50<br/>metric: COSINE

    Milvus-->>KB: 向量检索结果<br/>{doc_id, distance}

    KB->>KB: 混合检索合并
    Note over KB: BM25 + 向量检索<br/>加权融合

    KB->>Reranker: 重排序
    Note over KB,Reranker: top 100文档<br/>→ top 10

    Reranker->>Reranker: 计算相关性分数
    Note over Reranker: 使用Cross-Encoder<br/>模型：bge-reranker-v2

    Reranker-->>KB: 重排序结果<br/>top 10文档

    KB-->>User: 返回知识片段<br/>[{"content", "score"}]
```

---

## 5. 权限验证流程

### 5.1 流程概述

**业务场景**：用户访问资源时的权限验证

**涉及服务**：
- API网关
- 认证中间件
- 权限中间件
- 权限服务
- Redis缓存

### 5.2 时序图

```mermaid
sequenceDiagram
    autonumber
    participant User as 用户请求
    participant API as API网关
    participant Auth as 认证中间件
    participant Perm as 权限中间件
    participant PermSvc as 权限服务
    participant Redis as Redis
    participant MySQL as MySQL

    User->>API: GET /api/v1/bots/{bot_id}
    Note over User,API: Authorization: Bearer {token}

    API->>Auth: 验证Token
    Note over API,Auth: JWT签名验证

    Auth->>Redis: GET token:{token}
    Redis-->>Auth: Token有效或无效

    alt Token无效
        Auth-->>API: 401 Unauthorized
        API-->>User: 401 Unauthorized
    end

    Auth->>Auth: 解析Token获取user_id
    Auth-->>API: 注入用户上下文<br/>ctx.Set("user_id", xxx)

    API->>Perm: 检查权限
    Note over API,Perm: required_permission: "bot:read"

    Perm->>Redis: GET user_permissions:{user_id}
    Note over Perm,Redis: 尝试从缓存获取

    alt 缓存命中
        Redis-->>Perm: 用户权限列表
    else 缓存未命中
        Perm->>PermSvc: LoadUserPermissions
        Note over Perm,PermSvc: user_id + resource_type

        PermSvc->>MySQL: SELECT * FROM permissions<br/>WHERE user_id = ?
        MySQL-->>PermSvc: 权限列表

        PermSvc->>MySQL: SELECT * FROM roles<br/>WHERE id IN (role_ids)
        MySQL-->>PermSvc: 角色权限

        PermSvc-->>Perm: 合并后的权限列表

        Perm->>Redis: SET user_permissions:{user_id}
        Note over Perm,Redis: 缓存权限列表<br/>ttl: 1小时
    end

    Perm->>Perm: 检查是否有"bot:read"权限
    Note over Perm: 遍历权限列表<br/>查找匹配项

    alt 无权限
        Perm-->>API: 403 Forbidden
        API-->>User: 403 Forbidden<br/>(缺少权限)
    end

    Perm->>Perm: 检查数据权限
    Note over Perm: 是否有访问该bot_id的权限

    alt 无数据权限
        Perm-->>API: 403 Forbidden<br/>(无权访问该Bot)
        API-->>User: 403 Forbidden
    end

    Perm-->>API: 权限验证通过
    API->>API: 执行业务逻辑
    API-->>User: 200 OK<br/>(返回数据)
```

### 5.3 权限验证流程图

```mermaid
flowchart TD
    Start([用户请求]) --> Auth[认证中间件]
    Auth --> TokenValid{Token有效?}
    TokenValid -->|否| 401[返回401]
    TokenValid -->|是| Perm[权限中间件]

    Perm --> CacheHit{缓存命中?}
    CacheHit -->|是| CheckPerm[检查缓存权限]
    CacheHit -->|否| LoadPerm[加载权限]
    LoadPerm --> CachePerm[写入缓存]
    CachePerm --> CheckPerm

    CheckPerm --> HasPerm{有权限?}
    HasPerm -->|否| 403[返回403]
    HasPerm -->|是| DataPerm{数据权限?}

    DataPerm -->|否| Business[执行业务逻辑]
    DataPerm -->|是| CheckData{有数据权限?}
    CheckData -->|否| 403
    CheckData -->|是| Business

    Business --> Response([返回响应])
    401 --> End([结束])
    403 --> End
    Response --> End
```

---

## 6. 补充指南

### 6.1 时序图绘制规范

**工具推荐**：
- **Mermaid Live Editor**: https://mermaid.live/
- **PlantUML**: https://plantuml.com/
- **Draw.io**: https://app.diagrams.net/

**Mermaid语法示例**：

```mermaid
sequenceDiagram
    autonumber  # 自动编号
    participant A as 服务A
    participant B as 服务B

    A->>B: 同步调用
    B-->>A: 返回结果

    A->>B: 异步调用
    B--xA: 后台处理

    Note over A,B: 注释说明

    alt 成功分支
        B-->>A: 成功响应
    else 失败分支
        B-->>A: 失败响应
    end
```

**命名规范**：
- 参与者：使用中文名称 + 英文服务名
- 消息：使用动宾结构，如"创建Bot"、"查询用户"
- 注释：说明关键信息，如数据结构、超时时间

### 6.2 需要补充时序图的模块

#### P0优先级（核心流程）

| 序号 | 模块名称 | 流程名称 | 优先级 |
|------|---------|---------|--------|
| 1 | 租户管理 | ✅ 租户注册流程 | P0 |
| 2 | 租户管理 | 租户续费流程 | P0 |
| 3 | 租户管理 | 租户停用流程 | P0 |
| 4 | Bot管理 | ✅ Bot创建流程 | P0 |
| 5 | Bot管理 | Bot发布流程 | P0 |
| 6 | 对话服务 | ✅ 对话交互流程 | P0 |
| 7 | 知识库服务 | ✅ 知识库检索流程 | P0 |
| 8 | 权限服务 | ✅ 权限验证流程 | P0 |
| 9 | 用户管理 | 用户登录流程 | P0 |
| 10 | 组织管理 | 组织架构调整流程 | P0 |

#### P1优先级（重要流程）

| 序号 | 模块名称 | 流程名称 | 优先级 |
|------|---------|---------|--------|
| 11 | 工作流引擎 | 工作流执行流程 | P1 |
| 12 | 插件引擎 | 插件调度流程 | P1 |
| 13 | 记忆引擎 | 对话记忆存储流程 | P1 |
| 14 | 智能路由 | 意图识别流程 | P1 |
| 15 | 人机协同 | 人工审核流程 | P1 |
| 16 | 计费服务 | Token计量流程 | P1 |
| 17 | 监控服务 | Agent监控流程 | P1 |
| 18 | 审计服务 | 审计日志记录流程 | P1 |

#### P2优先级（辅助流程）

| 序号 | 模块名称 | 流程名称 | 优先级 |
|------|---------|---------|--------|
| 19 | 数据看板 | 报表生成流程 | P2 |
| 20 | 文件上传 | 文件上传流程 | P2 |
| 21 | 批量导入 | Excel导入流程 | P2 |
| 22 | 通知服务 | 消息推送流程 | P2 |
| 23 | 定时任务 | 定时任务执行流程 | P2 |

### 6.3 时序图模板

```mermaid
sequenceDiagram
    autonumber
    participant User as 用户
    participant API as API网关
    participant Service as 业务服务
    participant DB as MySQL/Redis
    participant Other as 其他服务

    Note over User,Other: 流程概述说明

    User->>API: 请求
    Note over User,API: 请求数据说明

    API->>Service: 业务命令

    Service->>DB: 查询/写入
    DB-->>Service: 返回结果

    Service->>Other: 调用其他服务
    Other-->>Service: 返回结果

    Service-->>API: 返回数据
    API-->>User: 响应

    Note over User: 流程结束说明
```

### 6.4 补充步骤

**步骤1**：选择模块和流程
- 确定要补充的模块
- 列出核心业务流程列表

**步骤2**：绘制时序图
- 使用Mermaid语法
- 遵循命名规范
- 包含关键节点和异常处理

**步骤3**：补充说明
- 流程概述
- 关键节点说明
- 异常处理策略

**步骤4**：审核检查
- 检查时序图完整性
- 验证逻辑正确性
- 确保符合规范

---

## 7. 附录

### 7.1 Mermaid速查表

| 语法 | 说明 | 示例 |
|------|------|------|
| `sequenceDiagram` | 创建时序图 | `sequenceDiagram` |
| `participant` | 定义参与者 | `participant A as 服务A` |
| `->` | 同步调用 | `A->>B: 请求` |
| `-->` | 异步调用 | `A->>B: 异步消息` |
| `-->>` | 返回响应 | `B-->>A: 响应` |
| `autonumber` | 自动编号 | `autonumber` |
| `Note over` | 注释 | `Note over A,B: 说明` |
| `alt/else` | 条件分支 | `alt 条件` |
| `loop` | 循环 | `loop 循环条件` |
| `opt` | 可选流程 | `opt 可选条件` |

### 7.2 PlantUML速查表

```plantuml
@startuml
actor 用户
participant "API网关" as API
database "MySQL" as DB

用户 -> API : 请求
API -> DB : 查询
DB --> API : 返回
API --> 用户 : 响应
@enduml
```

---

**文档版本**: v1.0
**最后更新**: 2025-12-30
**作者**: ZKER架构团队
