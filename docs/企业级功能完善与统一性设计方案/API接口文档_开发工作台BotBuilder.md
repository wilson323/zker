# API接口文档：开发工作台BotBuilder模块

**模块名称**: 开发工作台BotBuilder (Low-Code Bot Builder)
**设计文档**: 13-百应开发平台_开发工作台_BotBuilder.md + 13-BotBuilder_提示词版本管理补充_完整版.md
**版本**: v2.0.0 (含提示词版本管理)
**最后更新**: 2025-01-03
**优先级**: P1

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. Bot管理API](#2-bot管理api)
- [3. 提示词版本管理API](#3-提示词版本管理api)
- [4. 版本对比与回滚API](#4-版本对比与回滚api)
- [5. A/B测试API](#5-ab测试api)
- [6. Bot测试API](#6-bot测试api)
- [7. 协作管理API](#7-协作管理api)
- [8. 模板管理API](#8-模板管理api)
- [9. 数据模型](#9-数据模型)
- [10. 错误码定义](#10-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

BotBuilder是ZKER企业级SaaS平台的低代码Bot构建平台,通过**可视化编辑器 + 企业级管理 + 提示词版本管理**,实现企业快速构建、测试、部署和迭代Bot:

- ✅ **可视化构建** - 拖拽式Bot设计器,零代码创建Bot
- ✅ **提示词管理** - 提示词编辑、版本控制、A/B测试
- ✅ **技能集成** - 关联知识库、插件、工作流
- ✅ **在线测试** - 实时预览、调试模式、日志查看
- ✅ **版本管理** - Git-like提示词版本控制、Diff对比、智能回滚
- ✅ **A/B测试** - 多版本并行测试,数据驱动决策
- ✅ **团队协作** - 多人协作编辑、评论审核、操作日志
- ✅ **模板市场** - Bot模板库,快速创建

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 版本对比: go-diff/diffmatchpatch
- 统计分析: 自研统计引擎

**前端技术栈**:
- 框架: React 18 + TypeScript
- UI库: Semi Design
- 编辑器: Monaco Editor
- Diff视图: react-diff-viewer

### 1.3 核心特性

| 特性类别 | 功能 | 说明 |
|---------|------|------|
| **Bot构建** | 可视化编辑器 | 提示词、技能、知识库配置 |
| **版本管理** | 提示词版本控制 | Git-like版本管理、Diff对比 |
| **测试** | 在线测试 | 实时预览、调试模式 |
| **部署** | 发布管理 | 版本管理、灰度发布 |
| **企业级** | 团队协作 | 多人协作、权限控制 |
| **优化** | A/B测试 | 多版本并行测试 |

---

## 2. Bot管理API

### 2.1 创建Bot

**接口地址**: `POST /api/v1/bot-builder/bots`

**功能说明**: 创建新Bot

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "name": "客服助手",
  "description": "智能客服Bot,用于解答用户问题",
  "avatar": "https://example.com/avatar.png",
  "prompt_config": {
    "system_prompt": "你是一个专业的客服助手...",
    "user_prompt_template": "用户问题:{query}",
    "temperature": 0.7,
    "max_tokens": 2000,
    "top_p": 0.9
  },
  "skill_ids": ["skill-001", "skill-002"],
  "knowledge_base_ids": ["kb-001"],
  "plugin_ids": ["plugin-001"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "Bot创建成功",
  "data": {
    "id": "bot-20250103-001",
    "tenant_id": "tenant-123",
    "creator_id": 1001,
    "name": "客服助手",
    "description": "智能客服Bot,用于解答用户问题",
    "avatar": "https://example.com/avatar.png",
    "status": "draft",
    "version": "v1.0",
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

### 2.2 查询Bot列表

**接口地址**: `GET /api/v1/bot-builder/bots`

**查询参数**:
- page: 页码,默认1
- page_size: 每页数量,默认20
- status: 状态过滤 (draft/published/archived)
- keyword: 搜索关键词

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": "bot-001",
        "name": "客服助手",
        "description": "智能客服Bot",
        "avatar": "https://example.com/avatar.png",
        "status": "published",
        "version": "v2.0",
        "total_conversations": 1250,
        "avg_rating": 4.5,
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 2.3 查看Bot详情

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "bot-001",
    "name": "客服助手",
    "description": "智能客服Bot",
    "avatar": "https://example.com/avatar.png",
    "status": "published",
    "version": "v2.0",
    "prompt_config": {
      "system_prompt": "你是一个专业的客服助手...",
      "temperature": 0.7,
      "max_tokens": 2000
    },
    "skills": [
      {
        "id": "skill-001",
        "name": "知识问答"
      }
    ],
    "knowledge_bases": [
      {
        "id": "kb-001",
        "name": "产品知识库"
      }
    ],
    "statistics": {
      "total_conversations": 1250,
      "avg_rating": 4.5,
      "success_rate": 0.92
    }
  }
}
```

### 2.4 更新Bot

**接口地址**: `PUT /api/v1/bot-builder/bots/{bot_id}`

**请求参数**:
```json
{
  "name": "客服助手(升级版)",
  "description": "升级后的描述",
  "prompt_config": {
    "system_prompt": "更新后的系统提示词..."
  }
}
```

**功能说明**: 更新Bot配置,自动创建新版本

### 2.5 发布Bot

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/publish`

**功能说明**: 发布Bot到生产环境

**请求参数**:
```json
{
  "publish_type": "full",
  "description": "发布v2.0版本"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "Bot发布成功",
  "data": {
    "bot_id": "bot-001",
    "version": "v2.0",
    "published_at": "2025-01-03T11:00:00Z"
  }
}
```

### 2.6 下架Bot

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/unpublish`

### 2.7 删除Bot

**接口地址**: `DELETE /api/v1/bot-builder/bots/{bot_id}`

---

## 3. 提示词版本管理API

### 3.1 创建提示词版本

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/prompt-versions`

**功能说明**: 为Bot创建新的提示词版本

**请求参数**:
```json
{
  "version_name": "优化专业术语解释",
  "version_description": "优化了对技术术语的解释,使回答更易懂",
  "system_prompt": "你是一个专业的技术助手,擅长用通俗易懂的语言解释技术概念...",
  "user_prompt_template": "用户问题:{query}\n上下文:{context}",
  "temperature": 0.7,
  "max_tokens": 2000,
  "top_p": 0.9
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "版本创建成功",
  "data": {
    "version_id": "ver-20250103-001",
    "bot_id": "bot-001",
    "version_number": "v2.0",
    "version_name": "优化专业术语解释",
    "is_active": false,
    "created_at": "2025-01-03T10:30:00Z"
  }
}
```

**Go后端代码示例**:

```go
// CreateVersion 创建新版本
func (s *PromptVersionService) CreateVersion(
    ctx context.Context,
    req *CreateVersionRequest,
) (*CreateVersionResponse, error) {
    // 1. 获取当前激活版本(作为父版本)
    activeVersion, _ := s.versionRepo.GetActiveVersion(ctx, req.BotID)

    // 2. 生成新版本号
    nextVersionNumber := s.generateNextVersion(ctx, req.BotID)

    // 3. 构建版本对象
    version := &entity.PromptVersion{
        ID:                 uuid.New().String(),
        BotID:              req.BotID,
        VersionNumber:      nextVersionNumber,
        VersionName:        req.VersionName,
        VersionDescription: req.Description,
        SystemPrompt:       req.SystemPrompt,
        UserPromptTemplate: req.UserPromptTemplate,
        Temperature:        req.Temperature,
        MaxTokens:          req.MaxTokens,
        TopP:               req.TopP,
        ParentVersionID:    getPtrID(activeVersion),
        IsActive:           false, // 新版本默认不激活
        CreatedBy:          req.CreatedBy,
    }

    // 4. 保存版本
    if err := s.versionRepo.Create(ctx, version); err != nil {
        return nil, err
    }

    // 5. 记录变更历史
    if activeVersion != nil {
        s.recordChangeHistory(ctx, activeVersion, version, req.CreatedBy)
    }

    return &CreateVersionResponse{
        VersionID:     version.ID,
        VersionNumber: nextVersionNumber,
        Message:       "版本创建成功",
    }, nil
}
```

### 3.2 查询提示词版本列表

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/prompt-versions`

**查询参数**:
- page: 页码
- page_size: 每页数量
- is_active: 是否激活版本

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "items": [
      {
        "id": "ver-001",
        "version_number": "v2.0",
        "version_name": "优化专业术语解释",
        "version_description": "优化了对技术术语的解释",
        "is_active": true,
        "total_conversations": 156,
        "avg_rating": 4.5,
        "success_rate": 0.92,
        "created_at": "2025-01-03T10:30:00Z",
        "activated_at": "2025-01-03T11:00:00Z"
      }
    ]
  }
}
```

### 3.3 查看版本详情

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/prompt-versions/{version_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "ver-001",
    "version_number": "v2.0",
    "version_name": "优化专业术语解释",
    "system_prompt": "你是一个专业的技术助手...",
    "user_prompt_template": "用户问题:{query}",
    "temperature": 0.7,
    "max_tokens": 2000,
    "top_p": 0.9,
    "statistics": {
      "total_conversations": 156,
      "avg_rating": 4.5,
      "success_rate": 0.92,
      "avg_response_time_ms": 2350,
      "total_evaluations": 45
    }
  }
}
```

### 3.4 激活版本

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/prompt-versions/{version_id}/activate`

**功能说明**: 激活指定版本,停用当前版本

**请求参数**:
```json
{
  "reason": "v2.0效果测试良好,正式发布"
}
```

**Go后端代码示例**:

```go
// ActivateVersion 激活版本
func (s *PromptVersionService) ActivateVersion(
    ctx context.Context,
    versionID string,
    operatedBy int64,
) error {
    // 1. 获取目标版本
    version, err := s.versionRepo.GetByID(ctx, versionID)
    if err != nil {
        return err
    }

    // 2. 获取当前激活版本
    activeVersion, _ := s.versionRepo.GetActiveVersion(ctx, version.BotID)

    // 3. 停用当前版本
    if activeVersion != nil {
        activeVersion.IsActive = false
        s.versionRepo.Update(ctx, activeVersion)
    }

    // 4. 激活新版本
    version.IsActive = true
    version.ActivatedAt = time.Now()
    if err := s.versionRepo.Update(ctx, version); err != nil {
        return err
    }

    // 5. 记录变更历史
    if activeVersion != nil {
        s.historyRepo.Create(ctx, &entity.ChangeHistory{
            BotID:         version.BotID,
            VersionIDFrom: activeVersion.ID,
            VersionIDTo:   version.ID,
            ChangeType:    "update",
            ChangeSummary: fmt.Sprintf("从 %s 激活到 %s", activeVersion.VersionNumber, version.VersionNumber),
            OperatedBy:    operatedBy,
            OperationSource: "manual",
        })
    }

    return nil
}
```

---

## 4. 版本对比与回滚API

### 4.1 对比两个版本

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/prompt-versions/compare`

**功能说明**: 对比两个版本的差异,生成Diff视图

**请求参数**:
```json
{
  "version_id_from": "ver-001",
  "version_id_to": "ver-002",
  "view_mode": "side_by_side"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version_from": {
      "version_number": "v1.0",
      "system_prompt": "你是一个客服助手..."
    },
    "version_to": {
      "version_number": "v2.0",
      "system_prompt": "你是一个专业的技术助手..."
    },
    "field_changes": {
      "system_prompt": {
        "field_name": "系统提示词",
        "from": "你是一个客服助手...",
        "to": "你是一个专业的技术助手...",
        "diff_type": "modified"
      },
      "temperature": {
        "field_name": "温度",
        "from": 0.5,
        "to": 0.7,
        "diff_type": "modified"
      }
    },
    "summary": "本次更新包含2处变更",
    "html_diff": "<div class='diff'>...</div>",
    "side_by_side_diff": {
      "system_prompt_diff": {
        "from_lines": ["你是一个客服助手"],
        "to_lines": ["你是一个专业的技术助手"],
        "diff_type": "modified"
      }
    }
  }
}
```

### 4.2 查看变更历史

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/change-history`

**查询参数**:
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "items": [
      {
        "id": 1,
        "change_type": "update",
        "change_summary": "从 v1.0 激活到 v2.0",
        "version_from": "v1.0",
        "version_to": "v2.0",
        "operated_by": 1001,
        "operated_at": "2025-01-03T11:00:00Z",
        "operation_source": "manual"
      }
    ]
  }
}
```

### 4.3 回滚到指定版本

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/prompt-versions/rollback`

**功能说明**: 回滚到指定版本,自动创建回滚快照

**请求参数**:
```json
{
  "target_version_id": "ver-001",
  "reason": "v2.0出现严重问题,回滚到v1.0"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "回滚成功",
  "data": {
    "rollback_version_id": "ver-003",
    "rollback_version_number": "v2.1",
    "target_version": "v1.0",
    "rolled_back_at": "2025-01-03T12:00:00Z"
  }
}
```

**Go后端代码示例**:

```go
// RollbackVersion 回滚到指定版本
func (s *PromptVersionService) RollbackVersion(
    ctx context.Context,
    targetVersionID string,
    operatedBy int64,
    reason string,
) error {
    // 1. 获取目标版本
    targetVersion, err := s.versionRepo.GetByID(ctx, targetVersionID)
    if err != nil {
        return err
    }

    // 2. 获取当前版本
    currentVersion, _ := s.versionRepo.GetActiveVersion(ctx, targetVersion.BotID)

    // 3. 创建回滚快照(基于目标版本创建新版本)
    rollbackVersion := &entity.PromptVersion{
        ID:                 uuid.New().String(),
        BotID:              targetVersion.BotID,
        VersionNumber:      s.generateNextVersion(ctx, targetVersion.BotID),
        VersionName:        fmt.Sprintf("回滚到 %s", targetVersion.VersionNumber),
        VersionDescription: reason,
        SystemPrompt:       targetVersion.SystemPrompt,
        UserPromptTemplate: targetVersion.UserPromptTemplate,
        Temperature:        targetVersion.Temperature,
        MaxTokens:          targetVersion.MaxTokens,
        TopP:               targetVersion.TopP,
        ParentVersionID:    getPtrID(currentVersion),
        IsActive:           true, // 回滚版本立即激活
        CreatedBy:          operatedBy,
    }

    // 4. 保存回滚版本
    if err := s.versionRepo.Create(ctx, rollbackVersion); err != nil {
        return err
    }

    // 5. 停用当前版本
    if currentVersion != nil {
        currentVersion.IsActive = false
        s.versionRepo.Update(ctx, currentVersion)
    }

    // 6. 记录变更历史
    s.historyRepo.Create(ctx, &entity.ChangeHistory{
        BotID:         targetVersion.BotID,
        VersionIDFrom: getPtrID(currentVersion),
        VersionIDTo:   rollbackVersion.ID,
        ChangeType:    "rollback",
        ChangeSummary: fmt.Sprintf("回滚到 %s: %s", targetVersion.VersionNumber, reason),
        OperatedBy:    operatedBy,
        OperationSource: "manual",
    })

    return nil
}
```

---

## 5. A/B测试API

### 5.1 创建A/B测试实验

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/ab-tests`

**功能说明**: 创建A/B测试实验,对比多个版本效果

**请求参数**:
```json
{
  "experiment_name": "v2.0 vs v3.0 效果对比",
  "experiment_description": "对比优化后的提示词效果",
  "start_time": "2025-01-03T12:00:00Z",
  "end_time": "2025-01-10T12:00:00Z",
  "min_sample_size": 100,
  "traffic_allocation": {
    "v2.0": 50,
    "v3.0": 50
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "A/B测试实验创建成功",
  "data": {
    "experiment_id": "exp-20250103-001",
    "bot_id": "bot-001",
    "experiment_name": "v2.0 vs v3.0 效果对比",
    "status": "draft",
    "created_at": "2025-01-03T11:30:00Z"
  }
}
```

### 5.2 启动A/B测试

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/ab-tests/{experiment_id}/start`

**响应示例**:
```json
{
  "code": 0,
  "message": "A/B测试已启动",
  "data": {
    "experiment_id": "exp-001",
    "status": "running",
    "started_at": "2025-01-03T12:00:00Z"
  }
}
```

### 5.3 查询A/B测试列表

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/ab-tests`

**查询参数**:
- status: 状态过滤 (draft/running/completed)
- page: 页码
- page_size: 每页数量

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 3,
    "items": [
      {
        "id": "exp-001",
        "experiment_name": "v2.0 vs v3.0 效果对比",
        "status": "running",
        "start_time": "2025-01-03T12:00:00Z",
        "end_time": "2025-01-10T12:00:00Z",
        "traffic_allocation": {
          "v2.0": 50,
          "v3.0": 50
        },
        "created_at": "2025-01-03T11:30:00Z"
      }
    ]
  }
}
```

### 5.4 查看A/B测试详情

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/ab-tests/{experiment_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "exp-001",
    "experiment_name": "v2.0 vs v3.0 效果对比",
    "status": "running",
    "start_time": "2025-01-03T12:00:00Z",
    "traffic_allocation": {
      "v2.0": 50,
      "v3.0": 50
    },
    "results": {
      "v2.0": {
        "version_number": "v2.0",
        "total_conversations": 125,
        "avg_rating": 4.2,
        "success_rate": 0.85,
        "avg_response_time_ms": 2100
      },
      "v3.0": {
        "version_number": "v3.0",
        "total_conversations": 130,
        "avg_rating": 4.5,
        "success_rate": 0.92,
        "avg_response_time_ms": 2350
      }
    },
    "winner": {
      "version_id": "ver-003",
      "version_number": "v3.0",
      "confidence_level": 95.0,
      "statistical_significance": true
    }
  }
}
```

### 5.5 结束A/B测试并选择获胜版本

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/ab-tests/{experiment_id}/complete`

**请求参数**:
```json
{
  "winner_version_id": "ver-003",
  "auto_activate": true
}
```

**功能说明**: 结束A/B测试,自动激活获胜版本

---

## 6. Bot测试API

### 6.1 在线测试Bot

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/test`

**功能说明**: 在线测试Bot,支持调试模式

**请求参数**:
```json
{
  "query": "你好,请介绍一下你的功能",
  "conversation_id": "conv-test-001",
  "debug_mode": true,
  "test_version_id": "ver-002"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "bot_id": "bot-001",
    "query": "你好,请介绍一下你的功能",
    "answer": "你好!我是客服助手,主要功能包括:\n1. 智能问答\n2. 知识库检索\n3. 工单处理...",
    "conversation_id": "conv-test-001",
    "message_id": "msg-test-001",
    "debug_info": {
      "version_used": "v2.0",
      "prompt_tokens": 150,
      "completion_tokens": 300,
      "total_tokens": 450,
      "response_time_ms": 2350,
      "skills_used": ["skill-001", "skill-002"],
      "knowledge_bases_used": ["kb-001"],
      "retrieved_chunks": 5
    }
  }
}
```

### 6.2 查看测试日志

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/test-logs`

**查询参数**:
- page: 页码
- page_size: 每页数量
- start_date: 开始日期
- end_date: 结束日期

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 50,
    "items": [
      {
        "id": "log-001",
        "query": "你好",
        "answer": "你好!我是客服助手...",
        "conversation_id": "conv-test-001",
        "test_version": "v2.0",
        "response_time_ms": 2350,
        "created_at": "2025-01-03T12:00:00Z"
      }
    ]
  }
}
```

---

## 7. 协作管理API

### 7.1 添加协作者

**接口地址**: `POST /api/v1/bot-builder/bots/{bot_id}/collaborators`

**请求参数**:
```json
{
  "user_id": 1002,
  "role": "editor"
}
```

**角色类型**:
- `owner`: 所有者(完全控制)
- `editor`: 编辑者(可修改Bot)
- `viewer`: 查看者(仅查看)

**响应示例**:
```json
{
  "code": 0,
  "message": "协作者添加成功",
  "data": {
    "id": 1,
    "bot_id": "bot-001",
    "user_id": 1002,
    "role": "editor",
    "created_at": "2025-01-03T12:00:00Z"
  }
}
```

### 7.2 查询协作者列表

**接口地址**: `GET /api/v1/bot-builder/bots/{bot_id}/collaborators`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "user_id": 1001,
        "user_name": "张三",
        "role": "owner",
        "created_at": "2025-01-03T10:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1002,
        "user_name": "李四",
        "role": "editor",
        "created_at": "2025-01-03T11:00:00Z"
      }
    ]
  }
}
```

### 7.3 更新协作者角色

**接口地址**: `PUT /api/v1/bot-builder/bots/{bot_id}/collaborators/{user_id}`

**请求参数**:
```json
{
  "role": "viewer"
}
```

### 7.4 移除协作者

**接口地址**: `DELETE /api/v1/bot-builder/bots/{bot_id}/collaborators/{user_id}`

---

## 8. 模板管理API

### 8.1 创建Bot模板

**接口地址**: `POST /api/v1/bot-builder/templates`

**功能说明**: 将Bot保存为模板,供快速创建

**请求参数**:
```json
{
  "name": "客服助手模板",
  "description": "标准客服Bot模板",
  "bot_id": "bot-001",
  "category": "customer_service",
  "tags": ["客服", "问答"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "模板创建成功",
  "data": {
    "id": "tpl-001",
    "name": "客服助手模板",
    "description": "标准客服Bot模板",
    "category": "customer_service",
    "tags": ["客服", "问答"],
    "created_at": "2025-01-03T12:00:00Z"
  }
}
```

### 8.2 查询模板列表

**接口地址**: `GET /api/v1/bot-builder/templates`

**查询参数**:
- category: 分类过滤
- keyword: 搜索关键词

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "items": [
      {
        "id": "tpl-001",
        "name": "客服助手模板",
        "description": "标准客服Bot模板",
        "category": "customer_service",
        "tags": ["客服", "问答"],
        "usage_count": 25,
        "created_at": "2025-01-03T12:00:00Z"
      }
    ]
  }
}
```

### 8.3 从模板创建Bot

**接口地址**: `POST /api/v1/bot-builder/templates/{template_id}/create-bot`

**请求参数**:
```json
{
  "name": "我的客服Bot",
  "description": "基于模板创建的客服Bot"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "Bot创建成功",
  "data": {
    "bot_id": "bot-002",
    "template_id": "tpl-001",
    "name": "我的客服Bot",
    "created_at": "2025-01-03T12:30:00Z"
  }
}
```

---

## 9. 数据模型

### 9.1 Bot(Bot)
```typescript
interface Bot {
  id: string;
  tenant_id: string;
  creator_id: number;
  team_id?: number;

  // 基本信息
  name: string;
  description?: string;
  avatar?: string;

  // 配置
  prompt_config: PromptConfig;
  skill_ids: string[];
  knowledge_base_ids: string[];
  plugin_ids: string[];

  // 状态
  status: 'draft' | 'published' | 'archived';
  version: string;

  // 模板
  is_template: boolean;
  template_id?: string;

  // 审核
  review_status: 'draft' | 'pending_review' | 'approved' | 'rejected';

  // 统计
  total_conversations: number;
  avg_rating: number;
  success_rate: number;

  // 审计
  created_at: string;
  updated_at: string;
}
```

### 9.2 PromptConfig(提示词配置)
```typescript
interface PromptConfig {
  system_prompt: string;
  user_prompt_template?: string;
  temperature: number;
  max_tokens: number;
  top_p: number;
}
```

### 9.3 PromptVersion(提示词版本)
```typescript
interface PromptVersion {
  id: string;
  bot_id: string;
  version_number: string; // v1.0, v2.0...
  version_name?: string;
  version_description?: string;

  // 提示词内容
  system_prompt: string;
  user_prompt_template?: string;
  temperature?: number;
  max_tokens?: number;
  top_p?: number;

  // 版本关系
  parent_version_id?: string;
  is_active: boolean;
  is_ab_test: boolean;

  // 效果统计
  total_conversations: number;
  avg_rating: number;
  success_rate: number;
  avg_response_time_ms: number;
  total_evaluations: number;

  // A/B测试数据
  ab_test_traffic_percent?: number;
  ab_test_start_at?: string;
  ab_test_end_at?: string;

  // 元数据
  created_by?: number;
  created_at: string;
  activated_at?: string;
  archived_at?: string;
  archived_reason?: string;
}
```

### 9.4 ABTest(A/B测试实验)
```typescript
interface ABTest {
  id: string;
  bot_id: string;
  experiment_name: string;
  experiment_description?: string;

  // 实验配置
  status: 'draft' | 'running' | 'paused' | 'completed' | 'cancelled';
  start_time?: string;
  end_time?: string;
  min_sample_size: number;

  // 流量分配
  traffic_allocation: Record<string, number>; // {"v2.0": 50, "v3.0": 50}

  // 实验结果
  winner_version_id?: string;
  confidence_level?: number;
  statistical_significance?: boolean;

  // 创建信息
  created_by?: number;
  created_at: string;
}
```

### 9.5 BotCollaborator(Bot协作者)
```typescript
interface BotCollaborator {
  id: number;
  bot_id: string;
  user_id: number;
  role: 'owner' | 'editor' | 'viewer';
  created_at: string;
}
```

### 9.6 BotTemplate(Bot模板)
```typescript
interface BotTemplate {
  id: string;
  tenant_id: string;
  name: string;
  description?: string;
  bot_id: string; // 基于哪个Bot创建
  category?: string;
  tags?: string[];
  usage_count: number;
  created_at: string;
}
```

---

## 10. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 50001 | 400 | Bot不存在 |
| 50002 | 400 | Bot名称已存在 |
| 50003 | 403 | 无Bot操作权限 |
| 50101 | 400 | 提示词版本不存在 |
| 50102 | 400 | 版本号已存在 |
| 50103 | 400 | 无法激活已归档的版本 |
| 50104 | 400 | 回滚目标版本无效 |
| 50201 | 400 | A/B测试实验不存在 |
| 50202 | 400 | 实验状态不允许该操作 |
| 50203 | 400 | 流量分配无效 |
| 50301 | 400 | 协作者已存在 |
| 50302 | 403 | 无协作者管理权限 |
| 50401 | 400 | Bot模板不存在 |
| 50501 | 400 | Bot测试失败 |

---

**文档结束**
