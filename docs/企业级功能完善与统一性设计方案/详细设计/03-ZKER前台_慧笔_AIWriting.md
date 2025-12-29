# 03-慧笔_AIWriting 详细设计说明书

**文档编号**: DE-DD-2025-003
**模块名称**: 慧笔_AIWriting
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03
**最后更新**: 2025-01-03

---

## 📋 文档修订历史

| 版本 | 日期 | 修订人 | 修订说明 |
|------|------|--------|----------|
| v1.0.0 | 2025-01-03 | ZKER Team | 初始版本，基于配置驱动设计 |

---

## 1. 模块概述

### 1.1 模块定位

**慧笔_AIWriting** 是 ZKER 企业级 SaaS 平台的核心前台模块之一，通过**通用AI写作引擎 + 数据库模板配置**的方式，为企业用户提供智能内容创作能力。

**核心设计理念**：
- ✅ **配置优先**：80% 功能通过数据库配置实现
- ✅ **场景驱动**：支持 100+ 写作场景，通过模板配置扩展
- ✅ **企业级**：多租户隔离、权限控制、审计日志
- ✅ **低代码**：企业可自定义写作模板，无需编码

### 1.2 业务价值

| 受益者 | 价值 |
|--------|------|
| **企业用户** | 写作效率提升 10 倍，内容质量标准化，支持批量生成 |
| **企业管理员** | 通过模板库管控企业对外内容风格，支持品牌一致性 |
| **平台运营** | 预置 50+ 行业模板，企业可自建模板，形成模板生态 |

### 1.3 与鲸智百应功能对齐

| 鲸智百应功能 | ZKER 实现方式 | 实现策略 |
|-------------|--------------|----------|
| 公文写作（通知、报告、总结） | ✅ 数据库模板配置 | `writing_templates` 表 |
| 营销文案（朋友圈、小红书、公众号） | ✅ 模板 + 风格配置 | `writing_styles` 表 |
| 邮件写作（商务、求职、请假） | ✅ 模板 + 变量系统 | `template_variables` 表 |
| 多语言支持（中英日韩） | ✅ LLM 多语言能力 | 无需额外配置 |
| 智能润色（改写、扩写、缩写） | ✅ 通用写作引擎 | 后端服务实现 |
| 批量生成（一次生成多版本） | ✅ 引擎批量调用 | API 并发调用 |
| 历史记录与复用 | ✅ 数据库存储 | `writing_histories` 表 |

**实现策略**：✅ 20% 编码（通用引擎） + 80% 配置（数据库模板）

---

## 2. 功能需求

### 2.1 核心功能清单

#### 2.1.1 用户端功能（前台）

**F1 - 智能写作（核心功能）**
- F1.1 选择写作场景（公文/营销/邮件/报告/创意等）
- F1.2 输入写作需求（自然语言描述）
- F1.3 填充模板变量（如：收件人、主题、关键要点）
- F1.4 AI 生成初稿（支持流式输出）
- F1.5 交互式优化（继续扩写、缩写、改写、润色）
- F1.6 多版本生成（一次生成 3 个版本供选择）
- F1.7 导出文档（Markdown / Word / PDF）

**F2 - 模板库**
- F2.1 浏览模板（分类、搜索、收藏）
- F2.2 使用模板（一键加载模板变量）
- F2.3 评价模板（点赞、评分、评论）
- F2.4 自定义模板（企业管理员可创建企业私有模板）

**F3 - 历史记录**
- F3.1 查看历史（按时间、场景、模板筛选）
- F3.2 复用历史（基于历史内容二次创作）
- F3.3 删除历史（软删除，保留 30 天）

#### 2.1.2 管理端功能（后台）

**F4 - 模板管理（企业管理员）**
- F4.1 创建模板（定义变量、提示词、示例）
- F4.2 编辑模板（修改模板内容、变量）
- F4.3 发布/下架模板（控制模板可用性）
- F4.4 模板分类管理（公文/营销/邮件/报告等）
- F4.5 模板使用统计（查看使用次数、用户评价）

**F5 - 风格管理**
- F5.1 预置风格库（正式/活泼/专业/亲切等）
- F5.2 自定义风格（配置语气、用词偏好、句式风格）
- F5.3 风格与模板关联（模板默认风格）

**F6 - 内容审核**
- F6.1 敏感词过滤（预置 + 自定义敏感词库）
- F6.2 内容质量检测（重复率、连贯性）
- F6.3 审核日志（记录违规内容、处理结果）

### 2.2 非功能需求

| 需求类型 | 指标 | 说明 |
|---------|------|------|
| **性能** | 响应时间 | 首次生成 < 3 秒，流式输出延迟 < 500ms |
| **性能** | 并发能力 | 单租户 10 并发，系统 1000 并发 |
| **可用性** | 系统可用性 | 99.9% （月度） |
| **可扩展性** | 新增场景 | 通过数据库配置，无需编码 |
| **安全性** | 数据隔离 | 严格的多租户数据隔离 |
| **安全性** | 内容安全 | 敏感词过滤 + 人工审核 |

---

## 3. 架构设计

### 3.1 系统架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      前端层 (React 18)                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │ 写作编辑器   │  │ 模板库页面  │  │ 历史记录    │              │
│  │ (Monaco)    │  │ (展示页面)  │  │ (列表页)    │              │
│  └─────────────┘  └─────────────┘  └─────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                              ↓ HTTP/WebSocket
┌─────────────────────────────────────────────────────────────────┐
│                      API 网关层 (Hertz)                         │
├─────────────────────────────────────────────────────────────────┤
│  租户识别中间件 → 权限验证 → 限流控制 → 路由分发                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    业务服务层 (Go 1.23)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │            AIWritingService (核心服务)                    │   │
│  │  - ProcessWritingRequest()     (处理写作请求)             │   │
│  │  - OptimizeContent()          (内容优化)                  │   │
│  │  - BatchGenerate()            (批量生成)                  │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │            TemplateService (模板服务)                     │   │
│  │  - GetTemplate()              (获取模板)                  │   │
│  │  - ListTemplates()            (列出模板)                  │   │
│  │  - ValidateVariables()        (验证变量)                  │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │            StyleService (风格服务)                        │   │
│  │  - GetStyle()                 (获取风格)                  │   │
│  │  - ApplyStyle()               (应用风格)                  │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    核心引擎层 (配置驱动)                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         AIWritingEngine (通用AI写作引擎)                  │   │
│  │  - LoadTemplate()          (从数据库加载模板)             │   │
│  │  - FillVariables()         (填充模板变量)                 │   │
│  │  - BuildPrompt()           (构建LLM提示词)                │   │
│  │  - CallLLM()               (调用大模型)                   │   │
│  │  - StreamOutput()          (流式输出)                     │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         ContentOptimizer (内容优化器)                     │   │
│  │  - Rewrite()                (改写)                        │   │
│  │  - Expand()                 (扩写)                        │   │
│  │  - Summarize()              (缩写)                        │   │
│  │  - Polish()                 (润色)                        │   │
│  └──────────────────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │         ContentValidator (内容验证器)                     │   │
│  │  - CheckSensitiveWords()   (敏感词检测)                   │   │
│  │  - CheckQuality()          (质量检测)                     │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    数据访问层 (GORM)                             │
├─────────────────────────────────────────────────────────────────┤
│  writing_templates | writing_histories | writing_styles        │
│  template_variables | sensitive_words | audit_logs             │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                    基础设施层                                    │
├─────────────────────────────────────────────────────────────────┤
│  MySQL 8.4 | Redis 8.0 | LLM API | MinIO | Elasticsearch        │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 核心流程

#### 3.2.1 智能写作主流程

```
用户输入写作需求
        ↓
选择写作模板 (从 writing_templates 加载)
        ↓
填充模板变量 (template_variables)
        ↓
加载写作风格 (可选，从 writing_styles 加载)
        ↓
构建 LLM 提示词 (模板系统提示词 + 变量 + 风格)
        ↓
调用 LLM 生成 (stream 模式)
        ↓
内容验证 (敏感词检测 + 质量检测)
        ↓
保存历史记录 (writing_histories)
        ↓
返回生成结果 (流式输出给前端)
```

#### 3.2.2 内容优化流程

```
用户请求优化（扩写/缩写/改写/润色）
        ↓
加载优化提示词模板 (预定义在代码中)
        ↓
构建优化提示词 (原始内容 + 优化指令)
        ↓
调用 LLM 优化
        ↓
对比优化结果 (质量评分)
        ↓
返回优化后的内容
```

---

## 4. 数据库设计

### 4.1 ER 图

```
┌──────────────────┐         ┌──────────────────┐
│ writing_templates│         │ writing_styles   │
├──────────────────┤         ├──────────────────┤
│ id (PK)          │         │ id (PK)          │
│ tenant_id (FK)   │         │ tenant_id (FK)   │
│ category         │         │ name             │
│ name             │         │ config (JSON)    │
│ description      │         └──────────────────┘
│ system_prompt    │                 ↑
│ variables (JSON) │         ┌───────┴───────┐
│ example_output   │         │               │
│ style_id (FK)    │         │               │
│ is_public        │         │               │
│ usage_count      │         │               │
└─────────┬────────┘         │               │
          │                  │               │
          │                  │               │
          ↓          ┌───────┴───────┐       │
┌──────────────────┐│               │       │
│template_variables││               │       │
├──────────────────┤│               │       │
│ id (PK)          ││               │       │
│ template_id (FK) ││               │       │
│ name             ││               │       │
│ type             ││               │       │
│ required         ││               │       │
│ default_value    ││               │       │
│ options (JSON)   ││               │       │
└──────────────────┘│               │       │
                     │               │       │
                     │               │       │
┌────────────────────┴───────────────┴───────┴────────┐
│            writing_histories                        │
├─────────────────────────────────────────────────────┤
│ id (PK)                                            │
│ tenant_id (FK)                                     │
│ user_id (FK)                                       │
│ template_id (FK)                                   │
│ style_id (FK)                                      │
│ title                                             │
│ input_variables (JSON)                            │
│ generated_content                                  │
│ optimized_content                                  │
│ status (draft/published/archived)                  │
│ word_count                                        │
│ quality_score                                     │
│ created_at                                        │
└─────────────────────────────────────────────────────┘
```

### 4.2 核心表结构

#### 4.2.1 写作模板表 (writing_templates)

**核心设计**：通过数据库存储实现 80% 场景配置

```sql
CREATE TABLE writing_templates (
    id VARCHAR(64) PRIMARY KEY COMMENT '模板ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID (NULL表示平台预置模板)',
    category VARCHAR(50) NOT NULL COMMENT '分类: 公文/营销/邮件/报告/创意',
    name VARCHAR(100) NOT NULL COMMENT '模板名称',
    description VARCHAR(500) COMMENT '模板描述',
    icon VARCHAR(512) COMMENT '模板图标URL',

    -- 核心：模板配置
    system_prompt TEXT NOT NULL COMMENT '系统提示词模板 (支持{{变量}}占位符)',
    variables JSON NOT NULL COMMENT '模板变量定义',
    example_input JSON COMMENT '示例输入值',
    example_output TEXT COMMENT '示例输出',

    -- 风格关联
    style_id VARCHAR(64) COMMENT '默认风格ID',

    -- 控制字段
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否公开到模板市场',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    version INT DEFAULT 1 COMMENT '版本号',

    -- 统计字段
    usage_count INT DEFAULT 0 COMMENT '使用次数',
    favorite_count INT DEFAULT 0 COMMENT '收藏次数',
    avg_rating DECIMAL(3,2) DEFAULT 0.00 COMMENT '平均评分 (1-5)',

    -- 审计字段
    created_by BIGINT COMMENT '创建者ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL COMMENT '软删除时间',

    UNIQUE KEY uk_tenant_category_name (tenant_id, category, name),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_category (category),
    INDEX idx_is_public (is_public),
    INDEX idx_is_active (is_active),
    INDEX idx_usage_count (usage_count),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='写作模板表 - 核心配置表';
```

**变量定义 JSON 格式示例**：

```json
{
  "variables": [
    {
      "name": "recipient",
      "type": "text",
      "label": "收件人",
      "required": true,
      "placeholder": "请输入收件人姓名或职位",
      "default": "",
      "maxLength": 50
    },
    {
      "name": "purpose",
      "type": "select",
      "label": "邮件目的",
      "required": true,
      "options": [
        {"value": "inquiry", "label": "咨询"},
        {"value": "complaint", "label": "投诉"},
        {"value": "appreciation", "label": "感谢"},
        {"value": "invitation", "label": "邀请"}
      ],
      "default": "inquiry"
    },
    {
      "name": "keyPoints",
      "type": "textarea",
      "label": "关键要点",
      "required": true,
      "placeholder": "请输入邮件的关键要点，每行一个",
      "default": "",
      "maxLength": 500
    },
    {
      "name": "tone",
      "type": "radio",
      "label": "语气",
      "required": false,
      "options": [
        {"value": "formal", "label": "正式"},
        {"value": "friendly", "label": "友好"},
        {"value": "urgent", "label": "紧急"}
      ],
      "default": "formal"
    }
  ]
}
```

**系统提示词模板示例**（商务邮件）：

```
你是一位专业的商务邮件写作助手。请根据以下信息撰写一封商务邮件：

**收件人**: {{recipient}}
**邮件目的**: {{purpose}}
**关键要点**:
{{keyPoints}}

**语气要求**: {{tone}}

要求：
1. 邮件结构清晰，包含明确的主题、开头、正文、结尾
2. 语言简洁专业，避免冗余
3. 根据语气要求调整措辞
4. 字数控制在 {{tone === 'formal' ? '300-500' : '200-400'}} 字

请直接输出邮件内容，不要包含任何说明性文字。
```

#### 4.2.2 写作历史表 (writing_histories)

```sql
CREATE TABLE writing_histories (
    id VARCHAR(64) PRIMARY KEY COMMENT '历史记录ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    template_id VARCHAR(64) COMMENT '使用的模板ID',
    style_id VARCHAR(64) COMMENT '使用的风格ID',

    title VARCHAR(200) COMMENT '作品标题 (用户自定义)',
    input_variables JSON NOT NULL COMMENT '输入变量快照',
    generated_content LONGTEXT COMMENT 'AI生成的内容',
    optimized_content LONGTEXT COMMENT '优化后的内容',

    -- 元数据
    word_count INT COMMENT '字数',
    quality_score DECIMAL(3,2) COMMENT '质量评分 (0-1)',
    generation_time INT COMMENT '生成耗时 (毫秒)',
    llm_model VARCHAR(50) COMMENT '使用的LLM模型',
    llm_tokens INT COMMENT '消耗的Token数',

    -- 状态管理
    status ENUM('draft', 'published', 'archived') DEFAULT 'draft' COMMENT '状态',
    is_favorited BOOLEAN DEFAULT FALSE COMMENT '是否收藏',

    -- 审计
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_template_id (template_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='写作历史表';
```

#### 4.2.3 写作风格表 (writing_styles)

```sql
CREATE TABLE writing_styles (
    id VARCHAR(64) PRIMARY KEY COMMENT '风格ID',
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台预置风格)',
    name VARCHAR(50) NOT NULL COMMENT '风格名称',
    description VARCHAR(200) COMMENT '风格描述',

    -- 核心：风格配置
    config JSON NOT NULL COMMENT '风格配置',

    -- 统计
    usage_count INT DEFAULT 0 COMMENT '使用次数',

    -- 审计
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,

    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='写作风格表';
```

**风格配置 JSON 示例**：

```json
{
  "tone": "professional",
  "toneDescription": "专业、正式、严谨",
  "vocabulary": {
    "level": "advanced",
    "preferProfessional": true,
    "avoidSlang": true
  },
  "sentence": {
    "avgLength": "20-30",
    "structure": "varied",
    "avoidRepetition": true
  },
  "punctuation": {
    "formal": true
  },
  "emoji": {
    "allow": false
  },
  "systemPromptAddition": "\n\n**风格要求**: \n- 使用专业、正式的语言\n- 句式多样化，避免重复\n- 不使用表情符号和口语化表达\n- 保持客观、严谨的语气"
}
```

#### 4.2.4 模板变量表 (template_variables)

**可选设计**：用于高级场景，变量需要独立管理时使用

```sql
CREATE TABLE template_variables (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    template_id VARCHAR(64) NOT NULL COMMENT '模板ID',
    name VARCHAR(50) NOT NULL COMMENT '变量名称',
    type ENUM('text', 'textarea', 'select', 'radio', 'checkbox', 'date', 'number') NOT NULL COMMENT '变量类型',
    label VARCHAR(100) NOT NULL COMMENT '显示标签',
    required BOOLEAN DEFAULT TRUE COMMENT '是否必填',
    placeholder VARCHAR(200) COMMENT '占位符',
    default_value VARCHAR(500) COMMENT '默认值',
    options JSON COMMENT '选项配置 (select/radio/checkbox)',
    validation_rules JSON COMMENT '验证规则',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_template_name (template_id, name),
    INDEX idx_template_id (template_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='模板变量表 (可选)';
```

#### 4.2.5 敏感词表 (sensitive_words)

```sql
CREATE TABLE sensitive_words (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) COMMENT '租户ID (NULL表示平台级敏感词)',
    word VARCHAR(100) NOT NULL COMMENT '敏感词',
    category ENUM('politics', 'porn', 'violence', 'discrimination', 'custom') NOT NULL COMMENT '分类',
    severity ENUM('high', 'medium', 'low') DEFAULT 'medium' COMMENT '严重程度',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_word (word),
    INDEX idx_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='敏感词表';
```

### 4.3 索引设计

| 索引类型 | 表名 | 索引名 | 字段 | 用途 |
|---------|------|--------|------|------|
| **主键** | writing_templates | PRIMARY | id | 唯一标识 |
| **唯一** | writing_templates | uk_tenant_category_name | (tenant_id, category, name) | 防止重复 |
| **普通** | writing_templates | idx_tenant_id | tenant_id | 租户隔离查询 |
| **普通** | writing_templates | idx_category | category | 分类筛选 |
| **普通** | writing_templates | idx_is_public | is_public | 公开模板查询 |
| **普通** | writing_templates | idx_usage_count | usage_count | 热门模板排序 |
| **全文** | writing_templates | ft_name_desc | name, description | 模板搜索 |

---

## 5. API 设计

### 5.1 REST API 列表

#### 5.1.1 写作相关 API

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| POST | /api/v1/writing/generate | 智能写作（核心） | user |
| POST | /api/v1/writing/optimize | 内容优化 | user |
| POST | /api/v1/writing/batch-generate | 批量生成 | user |
| GET | /api/v1/writing/histories | 获取历史记录 | user |
| GET | /api/v1/writing/histories/:id | 获取历史详情 | user |
| DELETE | /api/v1/writing/histories/:id | 删除历史记录 | user |
| POST | /api/v1/writing/histories/:id/favorite | 收藏/取消收藏 | user |

#### 5.1.2 模板相关 API

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | /api/v1/writing/templates | 列出模板 | user |
| GET | /api/v1/writing/templates/:id | 获取模板详情 | user |
| GET | /api/v1/writing/templates/:id/variables | 获取模板变量 | user |
| POST | /api/v1/writing/templates | 创建模板（管理员） | admin |
| PUT | /api/v1/writing/templates/:id | 更新模板（管理员） | admin |
| DELETE | /api/v1/writing/templates/:id | 删除模板（管理员） | admin |
| POST | /api/v1/writing/templates/:id/publish | 发布/下架模板 | admin |
| GET | /api/v1/writing/templates/categories | 获取模板分类 | user |

#### 5.1.3 风格相关 API

| 方法 | 路径 | 功能 | 权限 |
|------|------|------|------|
| GET | /api/v1/writing/styles | 列出风格 | user |
| GET | /api/v1/writing/styles/:id | 获取风格详情 | user |
| POST | /api/v1/writing/styles | 创建风格（管理员） | admin |
| PUT | /api/v1/writing/styles/:id | 更新风格（管理员） | admin |
| DELETE | /api/v1/writing/styles/:id | 删除风格（管理员） | admin |

### 5.2 核心 API 详细设计

#### 5.2.1 智能写作 API

**请求**：
```http
POST /api/v1/writing/generate
Authorization: Bearer <token>
X-Tenant-ID: <tenant_id>
Content-Type: application/json

{
  "templateId": "tpl-20250103-001",
  "variables": {
    "recipient": "张总",
    "purpose": "inquiry",
    "keyPoints": "1. 询问产品价格\n2. 了解交付周期\n3. 索要产品目录",
    "tone": "formal"
  },
  "styleId": "style-professional",
  "options": {
    "stream": true,
    "maxTokens": 2000,
    "temperature": 0.7
  }
}
```

**响应（流式）**：
```
Content-Type: text/event-stream

data: {"type":"start","id":"gen-20250103-123456"}

data: {"type":"content","delta":"您好张总，"}

data: {"type":"content","delta":"我是XX公司的李明。"}

data: {"type":"content","delta":"特此致信，"}

...

data: {"type":"end","id":"gen-20250103-123456","wordCount":387,"tokens":512,"qualityScore":0.92}
```

**Go 代码示例**：

```go
type WritingService struct {
    templateRepo     *repository.TemplateRepository
    historyRepo      *repository.HistoryRepository
    styleRepo        *repository.StyleRepository
    llmClient        *llm.Client
    validator        *validator.ContentValidator
    logger           *zap.Logger
}

// GenerateRequest 生成请求
type GenerateRequest struct {
    TemplateID string                 `json:"templateId" binding:"required"`
    Variables  map[string]interface{} `json:"variables" binding:"required"`
    StyleID    string                 `json:"styleId"`
    Options    GenerateOptions        `json:"options"`
}

type GenerateOptions struct {
    Stream      bool    `json:"stream"`
    MaxTokens   int     `json:"maxTokens"`
    Temperature float64 `json:"temperature"`
}

// Generate 智能写作
func (s *WritingService) Generate(ctx context.Context, req *GenerateRequest) (<-chan StreamEvent, error) {
    // 1. 加载模板
    template, err := s.templateRepo.GetByID(ctx, req.TemplateID)
    if err != nil {
        return nil, err
    }

    // 2. 验证变量
    if err := s.validateVariables(template, req.Variables); err != nil {
        return nil, err
    }

    // 3. 加载风格（可选）
    var style *model.Style
    if req.StyleID != "" {
        style, _ = s.styleRepo.GetByID(ctx, req.StyleID)
    }

    // 4. 构建提示词
    systemPrompt := s.buildSystemPrompt(template, style)
    userPrompt := s.buildUserPrompt(template, req.Variables)

    // 5. 调用 LLM
    llmReq := &llm.ChatRequest{
        Model:       "gpt-4",
        Messages: []llm.Message{
            {Role: "system", Content: systemPrompt},
            {Role: "user", Content: userPrompt},
        },
        Temperature: req.Options.Temperature,
        MaxTokens:   req.Options.MaxTokens,
        Stream:      req.Options.Stream,
    }

    streamChan := make(chan StreamEvent, 100)

    go func() {
        defer close(streamChan)

        // 流式调用 LLM
        llmStream, err := s.llmClient.ChatStream(ctx, llmReq)
        if err != nil {
            streamChan <- StreamEvent{Type: "error", Data: err.Error()}
            return
        }

        var fullContent strings.Builder
        tokenCount := 0

        for chunk := range llmStream {
            if chunk.Error != nil {
                streamChan <- StreamEvent{Type: "error", Data: chunk.Error.Error()}
                return
            }

            fullContent.WriteString(chunk.Delta)
            tokenCount++

            streamChan <- StreamEvent{
                Type:  "content",
                Delta: chunk.Delta,
            }
        }

        // 6. 内容验证
        content := fullContent.String()
        if violations := s.validator.CheckSensitiveWords(content); len(violations) > 0 {
            streamChan <- StreamEvent{
                Type: "sensitive_warning",
                Data: violations,
            }
        }

        // 7. 保存历史
        history := &model.WritingHistory{
            ID:              generateID(),
            TenantID:        getTenantID(ctx),
            UserID:          getUserID(ctx),
            TemplateID:      req.TemplateID,
            StyleID:         req.StyleID,
            InputVariables:  req.Variables,
            GeneratedContent: content,
            WordCount:       countWords(content),
            LLMModel:        llmReq.Model,
            LLMTokens:       tokenCount,
            QualityScore:    s.validator.CalculateQuality(content),
        }
        s.historyRepo.Create(ctx, history)

        // 8. 发送结束事件
        streamChan <- StreamEvent{
            Type:       "end",
            WordCount:  history.WordCount,
            Tokens:     tokenCount,
            QualityScore: history.QualityScore,
        }
    }()

    return streamChan, nil
}

// buildSystemPrompt 构建系统提示词
func (s *WritingService) buildSystemPrompt(template *model.Template, style *model.Style) string {
    prompt := template.SystemPrompt

    // 应用风格
    if style != nil && style.Config.SystemPromptAddition != "" {
        prompt += style.Config.SystemPromptAddition
    }

    return prompt
}

// buildUserPrompt 构建用户提示词
func (s *WritingService) buildUserPrompt(template *model.Template, variables map[string]interface{}) string {
    // 从模板的 example_input 中提取提示词格式
    prompt := fmt.Sprintf("**模板**: %s\n\n", template.Name)

    // 填充变量
    for key, value := range variables {
        prompt += fmt.Sprintf("**%s**: %v\n", key, value)
    }

    return prompt
}

// validateVariables 验证变量
func (s *WritingService) validateVariables(template *model.Template, variables map[string]interface{}) error {
    var vars []TemplateVariable
    if err := json.Unmarshal(template.Variables, &vars); err != nil {
        return err
    }

    for _, v := range vars {
        value, exists := variables[v.Name]

        // 检查必填
        if v.Required && !exists {
            return fmt.Errorf("变量 %s 为必填项", v.Label)
        }

        // 类型验证
        if exists {
            if err := validateVariableType(v, value); err != nil {
                return err
            }
        }
    }

    return nil
}
```

#### 5.2.2 内容优化 API

**请求**：
```http
POST /api/v1/writing/optimize
Authorization: Bearer <token>
X-Tenant-ID: <tenant_id>
Content-Type: application/json

{
  "content": "原始内容...",
  "operation": "expand",
  "instructions": "请将这段话扩展到500字，重点补充数据支撑"
}
```

**响应**：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "optimizedContent": "优化后的内容...",
    "originalWordCount": 150,
    "optimizedWordCount": 487,
    "improvements": [
      "增加了数据支撑",
      "优化了语言表达",
      "补充了结构框架"
    ],
    "qualityScore": 0.91
  }
}
```

**Go 代码示例**：

```go
// OptimizeContent 内容优化
func (s *WritingService) OptimizeContent(ctx context.Context, req *OptimizeRequest) (*OptimizeResponse, error) {
    // 构建优化提示词
    prompt := s.buildOptimizePrompt(req.Operation, req.Content, req.Instructions)

    // 调用 LLM
    llmReq := &llm.ChatRequest{
        Model: "gpt-4",
        Messages: []llm.Message{
            {Role: "system", Content: "你是一位专业的内容编辑，擅长优化和改进文本质量。"},
            {Role: "user", Content: prompt},
        },
        Temperature: 0.7,
        MaxTokens:   3000,
    }

    resp, err := s.llmClient.Chat(ctx, llmReq)
    if err != nil {
        return nil, err
    }

    return &OptimizeResponse{
        OptimizedContent:     resp.Content,
        OriginalWordCount:    countWords(req.Content),
        OptimizedWordCount:   countWords(resp.Content),
        QualityScore:         s.validator.CalculateQuality(resp.Content),
    }, nil
}

// buildOptimizePrompt 构建优化提示词
func (s *WritingService) buildOptimizePrompt(operation, content, instructions string) string {
    operationMap := map[string]string{
        "expand":   "扩写",
        "summarize": "缩写",
        "rewrite":  "改写",
        "polish":   "润色",
    }

    operationName := operationMap[operation]

    prompt := fmt.Sprintf(`请对以下内容进行**%s**:

**原文**:
%s

`, operationName, content)

    if instructions != "" {
        prompt += fmt.Sprintf("**特殊要求**:\n%s\n\n", instructions)
    }

    prompt += `请直接输出优化后的内容，不要包含任何说明性文字。`

    return prompt
}
```

### 5.3 WebSocket API

用于实时协作编辑场景（高级功能，可选）

```javascript
// 客户端连接
const ws = new WebSocket('wss://api.zker.com/v1/writing/collaborate/session-id')

// 发送编辑操作
ws.send(JSON.stringify({
  type: 'edit',
  operation: 'insert',
  position: 120,
  content: '新增内容'
}))

// 接收协作更新
ws.onmessage = (event) => {
  const data = JSON.parse(event.data)
  if (data.type === 'remote_edit') {
    // 应用远程编辑
    applyRemoteEdit(data)
  }
}
```

---

## 6. 前端设计

### 6.1 页面结构

```
/writing
├── /editor              # 写作编辑器（主页）
│   ├── /templates-panel # 模板选择面板
│   ├── /variables-form  # 变量填写表单
│   └── /output-panel    # 输出结果面板
├── /templates           # 模板库
│   ├── /categories      # 分类浏览
│   └── /detail          # 模板详情
├── /histories           # 历史记录
└── /settings            # 设置（管理员）
    ├── /templates-manage
    └── /styles-manage
```

### 6.2 核心组件设计

#### 6.2.1 WritingEditor（写作编辑器）

**TypeScript 组件**：

```tsx
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Tabs, Button, Form, Input, Select, message } from '@douyinfe/semi-ui';
import { Editor } from '@monaco-editor/react';
import { writingAPI } from '@/services/api';
import styles from './WritingEditor.module.scss';

interface Template {
  id: string;
  name: string;
  description: string;
  category: string;
  systemPrompt: string;
  variables: TemplateVariable[];
}

interface TemplateVariable {
  name: string;
  type: 'text' | 'textarea' | 'select' | 'radio';
  label: string;
  required: boolean;
  placeholder?: string;
  default?: string;
  options?: Array<{ value: string; label: string }>;
}

export const WritingEditor: React.FC = () => {
  const [template, setTemplate] = useState<Template | null>(null);
  const [variables, setVariables] = useState<Record<string, any>>({});
  const [generatedContent, setGeneratedContent] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [activeTab, setActiveTab] = useState('input');

  // 加载模板
  useEffect(() => {
    const templateId = localStorage.getItem('selectedTemplateId');
    if (templateId) {
      loadTemplate(templateId);
    }
  }, []);

  const loadTemplate = async (templateId: string) => {
    try {
      const resp = await writingAPI.getTemplate(templateId);
      setTemplate(resp.data);

      // 初始化变量值
      const initialValues: Record<string, any> = {};
      resp.data.variables.forEach((v: TemplateVariable) => {
        if (v.default !== undefined) {
          initialValues[v.name] = v.default;
        }
      });
      setVariables(initialValues);
    } catch (error) {
      message.error('加载模板失败');
    }
  };

  const handleGenerate = async () => {
    if (!template) return;

    setIsGenerating(true);
    setGeneratedContent('');
    setActiveTab('output');

    try {
      const stream = await writingAPI.generate({
        templateId: template.id,
        variables: variables,
        styleId: localStorage.getItem('selectedStyleId') || undefined,
        options: {
          stream: true,
          maxTokens: 2000,
          temperature: 0.7,
        },
      });

      // 流式接收
      const reader = stream.getReader();
      const decoder = new TextDecoder();

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        const lines = chunk.split('\n');

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = JSON.parse(line.slice(6));
            if (data.type === 'content') {
              setGeneratedContent(prev => prev + data.delta);
            } else if (data.type === 'end') {
              setIsGenerating(false);
              message.success(`生成完成！字数: ${data.wordCount}`);
            } else if (data.type === 'sensitive_warning') {
              message.warning('检测到敏感词，请注意');
            }
          }
        }
      }
    } catch (error) {
      message.error('生成失败');
      setIsGenerating(false);
    }
  };

  const renderVariableInput = (variable: TemplateVariable) => {
    switch (variable.type) {
      case 'text':
        return (
          <Input
            placeholder={variable.placeholder}
            value={variables[variable.name] || ''}
            onChange={value => setVariables({ ...variables, [variable.name]: value })}
          />
        );
      case 'textarea':
        return (
          <Input.TextArea
            placeholder={variable.placeholder}
            rows={4}
            value={variables[variable.name] || ''}
            onChange={value => setVariables({ ...variables, [variable.name]: value })}
          />
        );
      case 'select':
        return (
          <Select
            placeholder={`请选择${variable.label}`}
            value={variables[variable.name]}
            onChange={value => setVariables({ ...variables, [variable.name]: value })}
            optionList={variable.options || []}
          />
        );
      default:
        return null;
    }
  };

  return (
    <div className={styles.editor}>
      <div className={styles.header}>
        <h2>{template?.name || 'AI 智能写作'}</h2>
        <Button
          theme="solid"
          type="primary"
          loading={isGenerating}
          onClick={handleGenerate}
          disabled={!template}
        >
          {isGenerating ? '生成中...' : '生成内容'}
        </Button>
      </div>

      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={[
          {
            itemKey: 'input',
            tab: '输入',
            content: (
              <div className={styles.inputPanel}>
                {template && (
                  <Form
                    labelPosition="left"
                    labelAlign="right"
                    labelWidth="120px"
                  >
                    {template.variables.map(variable => (
                      <Form.InputField
                        key={variable.name}
                        label={variable.label}
                        required={variable.required}
                      >
                        {renderVariableInput(variable)}
                      </Form.InputField>
                    ))}
                  </Form>
                )}
              </div>
            ),
          },
          {
            itemKey: 'output',
            tab: '生成结果',
            content: (
              <div className={styles.outputPanel}>
                <Editor
                  height="500px"
                  language="markdown"
                  value={generatedContent}
                  theme="vs-dark"
                  options={{
                    readOnly: true,
                    minimap: { enabled: false },
                    wordWrap: 'on',
                  }}
                />
                {generatedContent && (
                  <div className={styles.actions}>
                    <Button onClick={() => handleOptimize('expand')}>扩写</Button>
                    <Button onClick={() => handleOptimize('summarize')}>缩写</Button>
                    <Button onClick={() => handleOptimize('rewrite')}>改写</Button>
                    <Button onClick={() => handleOptimize('polish')}>润色</Button>
                    <Button onClick={() => handleExport('markdown')}>导出 Markdown</Button>
                    <Button onClick={() => handleExport('word')}>导出 Word</Button>
                  </div>
                )}
              </div>
            ),
          },
        ]}
      />
    </div>
  );
};
```

#### 6.2.2 TemplateLibrary（模板库）

```tsx
import React, { useState, useEffect } from 'react';
import { Card, Input, Select, Button, Empty, Tag } from '@douyinfe/semi-ui';
import { writingAPI } from '@/services/api';
import styles from './TemplateLibrary.module.scss';

export const TemplateLibrary: React.FC = () => {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [category, setCategory] = useState<string>('all');
  const [search, setSearch] = useState('');

  useEffect(() => {
    loadTemplates();
  }, [category, search]);

  const loadTemplates = async () => {
    try {
      const resp = await writingAPI.listTemplates({
        category: category === 'all' ? undefined : category,
        search,
      });
      setTemplates(resp.data.items);
    } catch (error) {
      message.error('加载模板失败');
    }
  };

  return (
    <div className={styles.library}>
      <div className={styles.filters}>
        <Input
          placeholder="搜索模板..."
          prefix={<IconSearch />}
          value={search}
          onChange={setSearch}
          style={{ width: 300 }}
        />
        <Select
          value={category}
          onChange={setCategory}
          optionList={[
            { value: 'all', label: '全部分类' },
            { value: '公文', label: '公文' },
            { value: '营销', label: '营销' },
            { value: '邮件', label: '邮件' },
            { value: '报告', label: '报告' },
            { value: '创意', label: '创意' },
          ]}
          style={{ width: 200 }}
        />
      </div>

      <div className={styles.grid}>
        {templates.map(template => (
          <Card
            key={template.id}
            className={styles.card}
            title={template.name}
            headerExtraContent={
              <Button
                size="small"
                onClick={() => handleUseTemplate(template.id)}
              >
                使用模板
              </Button>
            }
          >
            <p>{template.description}</p>
            <div className={styles.meta}>
              <Tag color="blue">{template.category}</Tag>
              <span>使用 {template.usageCount} 次</span>
              <span>评分 {template.avgRating} ⭐</span>
            </div>
          </Card>
        ))}
      </div>

      {templates.length === 0 && (
        <Empty
          title="暂无模板"
          description="请尝试切换分类或搜索关键词"
        />
      )}
    </div>
  );
};
```

### 6.3 样式设计

**SCSS 示例**（WritingEditor.module.scss）：

```scss
.editor {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f8f9fa;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    background: white;
    border-bottom: 1px solid #e5e7eb;

    h2 {
      margin: 0;
      font-size: 20px;
      font-weight: 600;
    }
  }

  .inputPanel {
    padding: 24px;
    background: white;
    min-height: 400px;
  }

  .outputPanel {
    position: relative;
    padding: 24px;
    background: #1e1e1e;

    .actions {
      display: flex;
      gap: 8px;
      padding: 12px 24px;
      background: white;
      border-top: 1px solid #e5e7eb;
    }
  }
}
```

---

## 7. 配置系统设计

### 7.1 预置模板数据

**平台级预置模板**（`tenant_id IS NULL`）：

```sql
-- 公文类模板
INSERT INTO writing_templates (id, tenant_id, category, name, description, system_prompt, variables, is_public) VALUES
('tpl-gongwen-notice', NULL, '公文', '通知公告', '用于发布公司内部通知、公告等正式文件',
'你是一位专业的公文写作助手。请根据以下信息撰写一份**通知**：

**标题**: {{title}}
**发文部门**: {{department}}
**主送**: {{recipients}}
**事由**: {{reason}}
**具体内容**:
{{content}}

**发文日期**: {{date}}

要求：
1. 严格按照公文格式：标题、主送机关、正文、落款、日期
2. 语言庄重、准确、简洁
3. 结构清晰，逻辑严谨
4. 字数控制在 {{length}} 字以内

请直接输出通知全文。',
'[
  {"name": "title", "type": "text", "label": "通知标题", "required": true, "placeholder": "例如：关于XX的通知"},
  {"name": "department", "type": "text", "label": "发文部门", "required": true, "default": "公司办公室"},
  {"name": "recipients", "type": "textarea", "label": "主送单位/人员", "required": true, "placeholder": "例如：公司各部门、全体员工"},
  {"name": "reason", "type": "textarea", "label": "发文事由", "required": true, "placeholder": "简要说明发文原因"},
  {"name": "content", "type": "textarea", "label": "具体内容", "required": true, "placeholder": "详细描述通知内容"},
  {"name": "date", "type": "date", "label": "发文日期", "required": true, "default": "today"},
  {"name": "length", "type": "select", "label": "字数要求", "required": false, "options": [{"value": "300", "label": "300字以内"}, {"value": "500", "label": "500字以内"}, {"value": "800", "label": "800字以内"}], "default": "500"}
]',
TRUE),

-- 营销类模板
('tpl-marketing-wechat', NULL, '营销', '微信公众号文章', '用于撰写微信公众号营销文案',
'你是一位资深的微信公众号运营专家。请根据以下信息撰写一篇公众号文章：

**文章主题**: {{topic}}
**目标读者**: {{audience}}
**核心卖点**: {{sellingPoints}}
**品牌/产品**: {{brand}}
**行动号召**: {{cta}}

要求：
1. 标题吸引眼球，使用数字、疑问、热点等方式
2. 开头3秒抓住读者注意力
3. 正文采用痛点-分析-解决方案结构
4. 语言生动有趣，贴近读者
5. 适当使用emoji增强可读性
6. 结尾明确行动号召
7. 字数 {{wordCount}} 字

请直接输出文章内容（包含标题）。',
'[
  {"name": "topic", "type": "text", "label": "文章主题", "required": true, "placeholder": "例如：如何通过AI提升工作效率"},
  {"name": "audience", "type": "select", "label": "目标读者", "required": true, "options": [{"value": "职场新人", "label": "职场新人"}, {"value": "企业管理者", "label": "企业管理者"}, {"value": "创业者", "label": "创业者"}, {"value": "技术人员", "label": "技术人员"}]},
  {"name": "sellingPoints", "type": "textarea", "label": "核心卖点", "required": true, "placeholder": "每行一个卖点"},
  {"name": "brand", "type": "text", "label": "品牌/产品名称", "required": true},
  {"name": "cta", "type": "text", "label": "行动号召", "required": true, "placeholder": "例如：立即点击下方链接免费试用"},
  {"name": "wordCount", "type": "select", "label": "字数", "required": false, "options": [{"value": "800", "label": "800字左右"}, {"value": "1500", "label": "1500字左右"}, {"value": "2000", "label": "2000字左右"}], "default": "1500"}
]',
TRUE),

-- 邮件类模板
('tpl-email-business', NULL, '邮件', '商务邮件', '用于撰写各类商务邮件',
'你是一位专业的商务邮件写作助手。请根据以下信息撰写一封商务邮件：

**收件人**: {{recipient}}
**发件人**: {{sender}}
**邮件目的**: {{purpose}}
**核心内容**:
{{keyPoints}}

**语气**: {{tone}}

要求：
1. 主题行简洁明确，包含核心信息
2. 开头得体，使用礼貌称呼
3. 正文结构清晰，分段论述
4. 语言简洁专业，避免冗余
5. 结尾明确下一步行动或期待
6. 严格按照所选语气调整措辞
7. 字数 200-400 字

请输出完整邮件（包含主题行）。',
'[
  {"name": "recipient", "type": "text", "label": "收件人", "required": true, "placeholder": "例如：张总 / 李经理"},
  {"name": "sender", "type": "text", "label": "发件人", "required": true, "placeholder": "你的姓名"},
  {"name": "purpose", "type": "select", "label": "邮件目的", "required": true, "options": [{"value": "inquiry", "label": "咨询/询价"}, {"value": "proposal", "label": "提案/合作"}, {"value": "followup", "label": "跟进"}, {"value": "appreciation", "label": "感谢"}, {"value": "apology", "label": "道歉"}]},
  {"name": "keyPoints", "type": "textarea", "label": "核心内容", "required": true, "placeholder": "每行一个要点"},
  {"name": "tone", "type": "radio", "label": "语气", "required": true, "options": [{"value": "formal", "label": "正式"}, {"value": "friendly", "label": "友好"}, {"value": "urgent", "label": "紧急"}], "default": "formal"}
]',
TRUE);
```

### 7.2 预置风格数据

```sql
-- 平台级预置风格
INSERT INTO writing_styles (id, tenant_id, name, description, config) VALUES
('style-professional', NULL, '专业正式', '适用于商务公文、报告等正式场合',
'{
  "tone": "professional",
  "toneDescription": "专业、正式、严谨",
  "vocabulary": {
    "level": "advanced",
    "preferProfessional": true,
    "avoidSlang": true
  },
  "sentence": {
    "avgLength": "20-30",
    "structure": "varied",
    "avoidRepetition": true
  },
  "punctuation": {
    "formal": true
  },
  "emoji": {
    "allow": false
  },
  "systemPromptAddition": "\n\n**风格要求**: \n- 使用专业、正式的语言\n- 句式多样化，避免重复\n- 不使用表情符号和口语化表达\n- 保持客观、严谨的语气"
}'),

('style-friendly', NULL, '友好亲切', '适用于内部沟通、客户服务等场景',
'{
  "tone": "friendly",
  "toneDescription": "友好、亲切、温暖",
  "vocabulary": {
    "level": "intermediate",
    "preferProfessional": false,
    "avoidSlang": false
  },
  "sentence": {
    "avgLength": "15-25",
    "structure": "simple",
    "avoidRepetition": true
  },
  "punctuation": {
    "formal": false
  },
  "emoji": {
    "allow": true,
    "maxCount": 3
  },
  "systemPromptAddition": "\n\n**风格要求**: \n- 使用友好、亲切的语言\n- 适当使用emoji增强亲和力\n- 句式简洁明了\n- 保持温暖、积极的语气"
}'),

('style-creative', NULL, '创意活泼', '适用于营销文案、社交媒体等创意场景',
'{
  "tone": "creative",
  "toneDescription": "创意、活泼、有趣",
  "vocabulary": {
    "level": "intermediate",
    "preferProfessional": false,
    "avoidSlang": false
  },
  "sentence": {
    "avgLength": "10-20",
    "structure": "varied",
    "avoidRepetition": false
  },
  "punctuation": {
    "formal": false
  },
  "emoji": {
    "allow": true,
    "maxCount": 10
  },
  "systemPromptAddition": "\n\n**风格要求**: \n- 使用创意、活泼的语言\n- 鼓励使用网络流行语\n- 多用emoji和特殊符号\n- 保持有趣、吸引人的语气\n- 可以使用反问、感叹等修辞手法"
}');
```

### 7.3 敏感词配置

```sql
-- 平台级敏感词（示例）
INSERT INTO sensitive_words (tenant_id, word, category, severity) VALUES
(NULL, '涉政敏感词1', 'politics', 'high'),
(NULL, '涉政敏感词2', 'politics', 'high'),
(NULL, '色情敏感词1', 'porn', 'high'),
(NULL, '色情敏感词2', 'porn', 'high'),
(NULL, '暴力敏感词1', 'violence', 'medium'),
(NULL, '歧视敏感词1', 'discrimination', 'medium');
```

---

## 8. 核心代码实现

### 8.1 通用 AI 写作引擎

**Go 后端核心引擎**：

```go
package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/zeromicro/go-zero/core/logx"
)

// AIWritingEngine 通用AI写作引擎
type AIWritingEngine struct {
	templateRepo TemplateRepository
	llmClient    LLMClient
	validator    ContentValidator
	logger       logx.Logger
}

// NewAIWritingEngine 创建引擎实例
func NewAIWritingEngine(
	templateRepo TemplateRepository,
	llmClient LLMClient,
	validator ContentValidator,
) *AIWritingEngine {
	return &AIWritingEngine{
		templateRepo: templateRepo,
		llmClient:    llmClient,
		validator:    validator,
		logger:       logx.WithContext(context.Background()),
	}
}

// ProcessRequest 处理写作请求（核心方法）
func (e *AIWritingEngine) ProcessRequest(
	ctx context.Context,
	req *WritingRequest,
) (*WritingResponse, error) {
	// 1. 加载模板（从数据库）
	tmpl, err := e.templateRepo.GetByID(ctx, req.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	// 2. 验证变量
	if err := e.validateVariables(tmpl, req.Variables); err != nil {
		return nil, fmt.Errorf("variable validation failed: %w", err)
	}

	// 3. 填充模板变量
	systemPrompt, err := e.fillTemplate(tmpl.SystemPrompt, req.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fill template: %w", err)
	}

	// 4. 应用风格（可选）
	if req.StyleID != "" {
		style, err := e.templateRepo.GetStyleByID(ctx, req.StyleID)
		if err == nil {
			systemPrompt += e.extractStylePrompt(style)
		}
	}

	// 5. 构建 LLM 请求
	userPrompt := e.buildUserPrompt(tmpl, req.Variables)
	llmReq := &LLMRequest{
		Model:       req.Model,
		Messages:    []Message{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userPrompt}},
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      req.Stream,
	}

	// 6. 调用 LLM
	if req.Stream {
		return e.streamGenerate(ctx, llmReq, tmpl)
	}

	return e.generate(ctx, llmReq, tmpl)
}

// generate 非流式生成
func (e *AIWritingEngine) generate(
	ctx context.Context,
	llmReq *LLMRequest,
	tmpl *Template,
) (*WritingResponse, error) {
	resp, err := e.llmClient.Chat(ctx, llmReq)
	if err != nil {
		return nil, err
	}

	// 内容验证
	if violations := e.validator.CheckSensitiveWords(resp.Content); len(violations) > 0 {
		e.logger.Infof("Detected sensitive words: %v", violations)
	}

	return &WritingResponse{
		Content:      resp.Content,
		WordCount:    countWords(resp.Content),
		TokensUsed:   resp.TotalTokens,
		QualityScore: e.validator.CalculateQuality(resp.Content),
	}, nil
}

// streamGenerate 流式生成
func (e *AIWritingEngine) streamGenerate(
	ctx context.Context,
	llmReq *LLMRequest,
	tmpl *Template,
) (*WritingResponse, error) {
	stream, err := e.llmClient.ChatStream(ctx, llmReq)
	if err != nil {
		return nil, err
	}

	// 流式处理逻辑
	// ...（略）

	return &WritingResponse{}, nil
}

// fillTemplate 填充模板变量
func (e *AIWritingEngine) fillTemplate(
	tmplText string,
	vars map[string]interface{},
) (string, error) {
	// 使用 Go template 语法填充
	t, err := template.New("writing").Parse(tmplText)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildUserPrompt 构建用户提示词
func (e *AIWritingEngine) buildUserPrompt(
	tmpl *Template,
	vars map[string]interface{},
) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("**模板**: %s\n\n", tmpl.Name))

	// 填充变量
	for key, value := range vars {
		prompt.WriteString(fmt.Sprintf("**%s**: %v\n", key, value))
	}

	return prompt.String()
}

// validateVariables 验证变量
func (e *AIWritingEngine) validateVariables(
	tmpl *Template,
	vars map[string]interface{},
) error {
	var variables []TemplateVariable
	if err := json.Unmarshal(tmpl.Variables, &variables); err != nil {
		return err
	}

	for _, v := range variables {
		value, exists := vars[v.Name]

		// 检查必填
		if v.Required && !exists {
			return fmt.Errorf("变量 %s 为必填项", v.Label)
		}

		// 类型验证
		if exists {
			if err := validateVariableType(v, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateVariableType 变量类型验证
func validateVariableType(v TemplateVariable, value interface{}) error {
	switch v.Type {
	case "text", "textarea":
		if str, ok := value.(string); ok {
			if v.MaxLength > 0 && len(str) > v.MaxLength {
				return fmt.Errorf("%s 长度超限", v.Label)
			}
		}
	case "select", "radio":
		// 验证是否在选项列表中
		// ...
	}

	return nil
}

// extractStylePrompt 提取风格提示词
func (e *AIWritingEngine) extractStylePrompt(style *Style) string {
	var config StyleConfig
	if err := json.Unmarshal(style.Config, &config); err != nil {
		return ""
	}

	return config.SystemPromptAddition
}
```

### 8.2 内容优化器

```go
package engine

// ContentOptimizer 内容优化器
type ContentOptimizer struct {
	llmClient LLMClient
}

// OptimizeOptions 优化选项
type OptimizeOptions struct {
	Operation    string  // expand, summarize, rewrite, polish
	Instructions string  // 自定义指令
	Temperature  float64 // 温度值
}

// Optimize 优化内容
func (o *ContentOptimizer) Optimize(
	ctx context.Context,
	content string,
	opts OptimizeOptions,
) (*OptimizeResult, error) {
	prompt := o.buildOptimizePrompt(opts.Operation, content, opts.Instructions)

	llmReq := &LLMRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "system", Content: "你是一位专业的内容编辑，擅长优化和改进文本质量。"},
			{Role: "user", Content: prompt},
		},
		Temperature: opts.Temperature,
		MaxTokens:   3000,
	}

	resp, err := o.llmClient.Chat(ctx, llmReq)
	if err != nil {
		return nil, err
	}

	return &OptimizeResult{
		OptimizedContent:   resp.Content,
		OriginalWordCount:  countWords(content),
		OptimizedWordCount: countWords(resp.Content),
		Improvements:       o.analyzeImprovements(content, resp.Content),
		QualityScore:       o.calculateQuality(resp.Content),
	}, nil
}

// buildOptimizePrompt 构建优化提示词
func (o *ContentOptimizer) buildOptimizePrompt(operation, content, instructions string) string {
	operationMap := map[string]string{
		"expand":   "扩写",
		"summarize": "缩写",
		"rewrite":  "改写",
		"polish":   "润色",
	}

	operationName := operationMap[operation]

	prompt := fmt.Sprintf(`请对以下内容进行**%s**:

**原文**:
%s

`, operationName, content)

	if instructions != "" {
		prompt += fmt.Sprintf("**特殊要求**:\n%s\n\n", instructions)
	}

	prompt += `请直接输出优化后的内容，不要包含任何说明性文字。`

	return prompt
}
```

---

## 9. 测试用例

### 9.1 单元测试

**Go 测试示例**：

```go
package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAIWritingEngine_ValidateVariables(t *testing.T) {
	engine := &AIWritingEngine{}

	tmpl := &Template{
		Variables: []byte(`[
			{"name": "title", "type": "text", "label": "标题", "required": true},
			{"name": "content", "type": "textarea", "label": "内容", "required": true}
		]`),
	}

	tests := []struct {
		name      string
		vars      map[string]interface{}
		wantError bool
	}{
		{
			name: "正常输入",
			vars: map[string]interface{}{
				"title":   "测试标题",
				"content": "测试内容",
			},
			wantError: false,
		},
		{
			name: "缺少必填字段",
			vars: map[string]interface{}{
				"title": "测试标题",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := engine.validateVariables(tmpl, tt.vars)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
```

### 9.2 集成测试

```go
func TestAIWritingEngine_ProcessRequest_Integration(t *testing.T) {
	// 跳过单元测试（需要真实 LLM）
	if testing.Short() {
		t.Skip()
	}

	// 使用 mock LLM client
	mockLLM := &MockLLMClient{}
	engine := NewAIWritingEngine(mockLLM, nil)

	req := &WritingRequest{
		TemplateID: "tpl-test-001",
		Variables: map[string]interface{}{
			"title":   "测试",
			"content": "测试内容",
		},
	}

	resp, err := engine.ProcessRequest(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Content)
	assert.Greater(t, resp.WordCount, 0)
}
```

---

## 10. 部署方案

### 10.1 环境变量配置

```bash
# .env.production

# 数据库
DB_HOST=mysql.internal
DB_PORT=3306
DB_NAME=zker_production
DB_USER=zker_user
DB_PASS=your_password

# Redis
REDIS_HOST=redis.internal
REDIS_PORT=6379
REDIS_PASS=your_password

# LLM API
LLM_PROVIDER=openai
LLM_API_KEY=sk-xxx
LLM_MODEL=gpt-4
LLM_BASE_URL=https://api.openai.com/v1

# 服务端口
HTTP_PORT=8888
GRPC_PORT=8889

# 日志级别
LOG_LEVEL=info
```

### 10.2 Kubernetes 部署

**Deployment YAML**：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zker-writing-service
  namespace: zker-production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: writing-service
  template:
    metadata:
      labels:
        app: writing-service
    spec:
      containers:
      - name: writing-service
        image: zker/writing-service:v1.0.0
        ports:
        - containerPort: 8888
          name: http
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: zker-config
              key: db.host
        - name: LLM_API_KEY
          valueFrom:
            secretKeyRef:
              name: zker-secrets
              key: llm.api.key
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8888
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8888
          initialDelaySeconds: 10
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: zker-writing-service
  namespace: zker-production
spec:
  selector:
    app: writing-service
  ports:
  - port: 80
    targetPort: 8888
  type: ClusterIP
```

### 10.3 数据库迁移

```bash
# 使用 Atlas 进行数据库迁移
atlas migrate diff writing_templates \
  --dir "file://migrations" \
  --to "file://schema/writing_templates.sql" \
  --dev-url "docker://mysql/8/dev"

# 执行迁移
atlas migrate apply \
  --dir "file://migrations" \
  --url "mysql://user:pass@localhost:3306/zker_production"
```

---

## 11. 监控与运维

### 11.1 关键指标监控

| 指标 | 说明 | 告警阈值 |
|------|------|----------|
| **API 响应时间** | P50, P95, P99 延迟 | P99 > 5s |
| **API 错误率** | 请求失败率 | > 1% |
| **LLM Token 消耗** | 每日 Token 用量 | 超过配额 |
| **并发请求数** | 当前并发数 | > 80% 容量 |
| **模板使用率** | Top 10 模板使用次数 | - |

### 11.2 日志规范

```json
{
  "level": "info",
  "timestamp": "2025-01-03T10:30:45Z",
  "tenant_id": "tenant-123",
  "user_id": 456,
  "template_id": "tpl-email-001",
  "action": "generate",
  "duration_ms": 2340,
  "tokens_used": 512,
  "word_count": 387,
  "status": "success"
}
```

---

## 12. 扩展性设计

### 12.1 如何新增写作场景

**方式1：通过数据库配置（推荐）**

```sql
INSERT INTO writing_templates (id, category, name, description, system_prompt, variables, is_public) VALUES
('tpl-custom-001', '营销', '小红书种草文案', '用于撰写小红书种草笔记',
'你是一位小红书种草达人。请根据以下信息撰写一篇种草笔记：

**产品**: {{product}}
**使用场景**: {{scenario}}
**核心卖点**: {{sellingPoints}}
**个人体验**: {{experience}}

要求：
1. 标题使用emoji + 数字 + 疑问句式
2. 内容真实感强，使用"我发现"、"亲测"等词汇
3. 多用emoji增加趣味性
4. 首图关键信息突出
5. 语气亲切，像和朋友分享
6. 字数 300-500 字

请直接输出笔记内容（包含标题）。',
'[...变量定义...]', TRUE);
```

**方式2：通过管理后台界面**

1. 登录企业管理后台
2. 进入"模板管理" → "创建模板"
3. 填写模板信息（名称、描述、分类）
4. 配置模板变量（变量名、类型、是否必填）
5. 编写系统提示词（支持 `{{变量名}}` 占位符）
6. 保存并发布

### 12.2 如何新增写作风格

```sql
INSERT INTO writing_styles (id, tenant_id, name, description, config) VALUES
('style-humor', NULL, '幽默风趣', '适用于轻松幽默的写作场景',
'{
  "tone": "humor",
  "toneDescription": "幽默、风趣、轻松",
  "vocabulary": {
    "level": "intermediate",
    "preferProfessional": false,
    "avoidSlang": false
  },
  "sentence": {
    "avgLength": "10-20",
    "structure": "varied",
    "avoidRepetition": false
  },
  "emoji": {
    "allow": true,
    "maxCount": 15
  },
  "systemPromptAddition": "\n\n**风格要求**: \n- 使用幽默、风趣的语言\n- 适当使用网络梗和流行语\n- 多用emoji和特殊符号\n- 保持轻松、搞笑的语气\n- 可以使用夸张、自嘲等修辞手法"
}');
```

---

## 13. 总结

### 13.1 实施策略总结

| 实施项 | 实施方式 | 工作量 | 说明 |
|-------|---------|--------|------|
| **通用AI写作引擎** | 💻 独立编码 | 20% | 核心引擎，一次开发长期复用 |
| **预置模板库** | 📊 数据库配置 | 50% | 50+ 预置模板数据初始化 |
| **企业自定义模板** | 📊 数据库配置 | 10% | 企业管理员通过后台配置 |
| **风格系统** | 📊 数据库配置 | 10% | 预置风格 + 企业自定义风格 |
| **前端编辑器** | 💻 独立编码 | 10% | Monaco Editor + 表单组件 |

**总计**：20% 编码 + 80% 配置

### 13.2 与鲸智百应对齐情况

| 功能模块 | 对齐情况 | 实现方式 |
|---------|---------|----------|
| 公文写作 | ✅ 100% | 预置模板（通知、报告、总结等） |
| 营销文案 | ✅ 100% | 预置模板（微信、小红书、朋友圈等） |
| 邮件写作 | ✅ 100% | 预置模板（商务、求职、请假等） |
| 多语言支持 | ✅ 100% | LLM 原生能力 |
| 智能润色 | ✅ 100% | 内容优化引擎 |
| 批量生成 | ✅ 100% | API 并发调用 |
| 历史记录 | ✅ 100% | 数据库存储 |

### 13.3 企业级特性

- ✅ **多租户隔离**：严格的租户级数据隔离
- ✅ **权限控制**：基于 RBAC 的细粒度权限
- ✅ **审计日志**：完整的操作记录和追溯
- ✅ **内容安全**：敏感词过滤 + 质量检测
- ✅ **可扩展性**：通过数据库配置新增场景，无需编码
- ✅ **高性能**：流式输出 + 并发处理 + Redis 缓存

---

**文档结束**
