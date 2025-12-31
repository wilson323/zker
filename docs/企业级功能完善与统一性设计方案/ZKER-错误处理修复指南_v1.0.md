# ZKER 错误处理修复指南

**版本**: v1.0
**日期**: 2025-01-01
**作者**: ZKER企业级功能完善团队

---

## 📋 目录

1. [概述](#概述)
2. [问题分析](#问题分析)
3. [错误处理规范](#错误处理规范)
4. [修复策略](#修复策略)
5. [模块修复指南](#模块修复指南)
6. [测试验证](#测试验证)
7. [常见问题](#常见问题)

---

## 概述

### 背景

ZKER项目中发现**189处**未使用统一错误码系统的错误处理违规。这些违规主要分为3类：

1. **直接使用 `errors.New()`** (124处, 65.6%)
2. **使用 `fmt.Errorf()`** (58处, 30.7%)
3. **缺少上下文信息** (112处, 59.3%)

### 目标

统一所有错误处理格式，使用企业级错误码系统，提供：
- ✅ 统一的错误码（中英文双语）
- ✅ 丰富的上下文信息（KV对）
- ✅ 一致的错误响应格式
- ✅ 便于监控和排查

### 修复范围

- **优先级 P0**: Domain Service层（核心业务逻辑，80处）
- **优先级 P1**: Application层 + Domain其他层（60处）
- **优先级 P2**: Infra层和其他（49处）

---

## 问题分析

### 错误类型1: 直接使用 errors.New()

**❌ 错误代码**:
```go
// backend/domain/knowledge/service/knowledge.go:322
return errors.New("document is empty")

// backend/application/workflow/chatflow.go:530
return nil, errors.New("the role of the last day message must be user")

// backend/domain/workflow/service/service_impl.go:458
return errors.New("workflow not found")
```

**问题**:
- ❌ 没有使用统一的错误码
- ❌ 没有提供上下文信息（如ID、操作类型）
- ❌ 错误消息硬编码，不支持国际化
- ❌ 难以监控和统计

**✅ 正确代码**:
```go
return errorx.New(errno.ErrKnowledgeInvalidParamCode,
    errorx.KV("reason", "document is empty"),
    errorx.KV("document_id", docID),
    errorx.KV("operation", "save_document"),
)

return errorx.New(errno.ErrConversationInvalidParamCode,
    errorx.KV("reason", "last message role must be user"),
    errorx.KV("conversation_id", convID),
)

return errorx.New(errno.ErrWorkflowNotFoundCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("operation", "get_workflow"),
)
```

### 错误类型2: 使用 fmt.Errorf()

**❌ 错误代码**:
```go
// backend/domain/datacopy/service/datacopy.go:538
return fmt.Errorf("unsupported copy scene '%s'", metaInfo.scene)

// backend/domain/developer/service/api_key_management_service.go:245
return fmt.Errorf("user %d does not have access to space %d", uid, spaceID)

// backend/application/workflow/workflow.go:678
return fmt.Errorf("workflow execution failed: %w", err)
```

**问题**:
- ❌ 错误消息动态拼接，难以统一处理
- ❌ 没有结构化的错误码
- ❌ 上下文信息不完整

**✅ 正确代码**:
```go
return errorx.New(errno.ErrDataCopyInvalidSceneCode,
    errorx.KV("scene", metaInfo.scene),
    errorx.KV("supported_scenes", []string{"full", "incremental"}),
)

return errorx.New(errno.ErrPermissionDeniedCode,
    errorx.KV("user_id", uid),
    errorx.KV("space_id", spaceID),
    errorx.KV("reason", "user does not have access to space"),
)

return errorx.WrapByCode(err, errno.ErrWorkflowExecutionFailedCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("node_id", nodeID),
)
```

### 错误类型3: 缺少上下文信息

**❌ 错误代码**:
```go
// 直接返回，没有上下文
return err

// 只有简单包装
return fmt.Errorf("failed: %w", err)

// 错误码正确但缺少KV对
return errorx.New(errno.ErrKnowledgeNotFoundCode)
```

**问题**:
- ❌ 缺少关键信息（ID、操作类型、参数值）
- ❌ 难以定位问题根因
- ❌ 无法进行精准监控

**✅ 正确代码**:
```go
return errorx.WrapByCode(err, errno.ErrKnowledgeDBCode,
    errorx.KV("operation", "save_document"),
    errorx.KV("document_id", docID),
    errorx.KV("knowledge_id", knowledgeID),
)

return errorx.WrapByCode(err, errno.ErrWorkflowExecutionFailedCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("node_id", nodeID),
    errorx.KV("node_type", "llm"),
    errorx.KV("error", err.Error()),
)
```

---

## 错误处理规范

### 核心原则

1. **统一错误码**: 所有错误必须使用 `errno` 包中定义的错误码
2. **结构化上下文**: 使用 `errorx.KV()` 提供关键上下文信息
3. **错误包装**: 对底层错误使用 `errorx.WrapByCode()` 包装
4. **双语支持**: 错误消息支持中英文（在errno包的i18n中定义）
5. **禁止硬编码**: 错误消息不要硬编码在业务代码中

### 标准格式

#### 1. 新建错误（无底层错误）

```go
// 格式
return errorx.New(errno.ErrXxxCode,
    errorx.KV("key1", "value1"),
    errorx.KV("key2", "value2"),
)

// 示例
return errorx.New(errno.ErrKnowledgeNotFoundCode,
    errorx.KV("knowledge_id", knowledgeID),
    errorx.KV("operation", "get_knowledge"),
)
```

#### 2. 包装错误（有底层错误）

```go
// 格式
return errorx.WrapByCode(err, errno.ErrXxxCode,
    errorx.KV("key1", "value1"),
    errorx.KV("key2", "value2"),
)

// 示例
return errorx.WrapByCode(err, errno.ErrKnowledgeDBCode,
    errorx.KV("operation", "save_document"),
    errorx.KV("document_id", docID),
)
```

### 必需的上下文信息

#### 1. 资源ID相关错误

```go
// ✅ Good：包含资源ID
return errorx.New(errno.ErrKnowledgeNotFoundCode,
    errorx.KV("knowledge_id", knowledgeID),
)

// ❌ Bad：缺少ID
return errorx.New(errno.ErrKnowledgeNotFoundCode)
```

#### 2. 操作相关错误

```go
// ✅ Good：包含操作类型
return errorx.WrapByCode(err, errno.ErrKnowledgeDBCode,
    errorx.KV("operation", "save_document"),
    errorx.KV("document_id", docID),
)

// ❌ Bad：没有操作信息
return errorx.WrapByCode(err, errno.ErrKnowledgeDBCode)
```

#### 3. 权限相关错误

```go
// ✅ Good：包含用户ID、资源ID、权限类型
return errorx.New(errno.ErrPermissionDeniedCode,
    errorx.KV("user_id", userID),
    errorx.KV("resource_type", "knowledge"),
    errorx.KV("resource_id", knowledgeID),
    errorx.KV("required_permission", "write"),
)

// ❌ Bad：信息不足
return errorx.New(errno.ErrPermissionDeniedCode)
```

#### 4. 参数验证错误

```go
// ✅ Good：包含参数名、参数值、错误原因
return errorx.New(errno.ErrKnowledgeInvalidParamCode,
    errorx.KV("param", "document_type"),
    errorx.KV("value", docType),
    errorx.KV("expected", "pdf,docx,txt"),
)

// ❌ Bad：没有参数详情
return errorx.New(errno.ErrKnowledgeInvalidParamCode)
```

### 常用KV键名

| 键名 | 说明 | 示例 |
|------|------|------|
| `operation` | 操作类型 | `save_document`, `get_workflow` |
| `resource_id` | 资源ID | `knowledge_id`, `workflow_id` |
| `user_id` | 用户ID | `12345` |
| `tenant_id` | 租户ID | `67890` |
| `param` | 参数名 | `document_type`, `limit` |
| `value` | 参数值 | `pdf`, `100` |
| `reason` | 错误原因 | `document is empty` |
| `error` | 原始错误 | `connection timeout` |
| `node_id` | 工作流节点ID | `node_123` |
| `node_type` | 节点类型 | `llm`, `code` |

---

## 修复策略

### 修复优先级

#### P0 - 最高优先级（Domain Service层）

**范围**: `backend/domain/*/service/*.go`

**影响**: 核心业务逻辑，直接影响用户体验

**修复顺序**:
1. Knowledge模块 (`domain/knowledge/service/`)
2. Workflow模块 (`domain/workflow/service/`)
3. Conversation模块 (`domain/conversation/service/`)
4. Permission模块 (`domain/permission/service/`)
5. Bot模块 (`domain/bot/service/`, `domain/agent/service/`)

**预计修复**: 80处

#### P1 - 高优先级（Application层）

**范围**: `backend/application/**/*.go`

**影响**: 应用编排层，影响多个领域服务

**修复顺序**:
1. Workflow模块 (`application/workflow/`)
2. Knowledge模块 (`application/knowledge/`)
3. Conversation模块 (`application/conversation/`)
4. Tenant模块 (`application/tenant/`)

**预计修复**: 60处

#### P2 - 中优先级（其他层）

**范围**: `backend/infra/**/*.go`, `backend/crossdomain/**/*.go`

**影响**: 基础设施层，跨领域层

**修复顺序**:
1. Infra层
2. Crossdomain层
3. API层

**预计修复**: 49处

### 修复流程

#### 步骤1: 识别违规

```bash
# 使用自动化工具扫描
python tools/error_handler_migrator_enhanced.py

# 查看报告
cat docs/企业级功能完善与统一性设计方案/ZKER-错误处理迁移报告.md
```

#### 步骤2: 选择错误码

参考 `backend/types/errno/*.go` 中的错误码定义：

```go
// Knowledge模块错误码 (105xxxxx)
const (
    ErrKnowledgeInvalidParamCode     = 105000000
    ErrKnowledgePermissionCode       = 105000001
    ErrKnowledgeDBCode               = 105000003
    ErrKnowledgeNotExistCode         = 105000011
    ErrKnowledgeDocumentNotExistCode = 105000012
    // ...
)

// Workflow模块错误码 (203xxxxx)
const (
    ErrWorkflowInvalidParamCode    = 203000000
    ErrWorkflowNotFoundCode        = 203000001
    ErrWorkflowExecutionFailedCode = 203000050
    // ...
)
```

#### 步骤3: 修复代码

**示例1: Knowledge模块**

```go
// Before
return errors.New("document not found")

// After
return errorx.New(errno.ErrKnowledgeDocumentNotExistCode,
    errorx.KV("document_id", docID),
    errorx.KV("operation", "get_document"),
)
```

**示例2: Workflow模块**

```go
// Before
return fmt.Errorf("workflow %s execution failed: %w", workflowID, err)

// After
return errorx.WrapByCode(err, errno.ErrWorkflowExecutionFailedCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("node_id", nodeID),
    errorx.KV("node_type", nodeType),
)
```

#### 步骤4: 添加import

确保文件顶部有必要的import：

```go
import (
    "github.com/coze-dev/coze-studio/backend/pkg/errorx"
    "github.com/coze-dev/coze-studio/backend/types/errno"
)
```

#### 步骤5: 测试验证

```bash
# 单元测试
cd backend/domain/knowledge/service
go test -v -run TestKnowledge

# 集成测试
cd backend
go test ./... -cover

# Lint检查
golangci-lint run
```

---

## 模块修复指南

### Knowledge模块修复指南

**文件**: `backend/domain/knowledge/service/knowledge.go`

**错误数量**: 15处

#### 修复清单

| 行号 | Before | After | 状态 |
|------|--------|-------|------|
| 322 | `errors.New("document is empty")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ |
| 458 | `errors.New("invalid document type")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ |
| 538 | `errors.New("knowledge not found")` | `errorx.New(errno.ErrKnowledgeNotExistCode, ...)` | ⏳ |

#### 修复示例

```go
// 文件: backend/domain/knowledge/service/knowledge.go:322

// ❌ Before
if document == nil || document.Content == "" {
    return errors.New("document is empty")
}

// ✅ After
if document == nil || document.Content == "" {
    return errorx.New(errno.ErrKnowledgeInvalidParamCode,
        errorx.KV("reason", "document is empty"),
        errorx.KV("document_id", documentID),
        errorx.KV("operation", "save_document"),
    )
}
```

#### 错误码映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| document is empty | ErrKnowledgeInvalidParamCode | document_id, operation |
| document not found | ErrKnowledgeDocumentNotExistCode | document_id, operation |
| invalid document type | ErrKnowledgeInvalidParamCode | document_id, type, expected |
| knowledge not found | ErrKnowledgeNotExistCode | knowledge_id, operation |
| parse failed | ErrKnowledgeParserParseFailCode | document_id, parser_type |

### Workflow模块修复指南

**文件**: `backend/domain/workflow/service/service_impl.go`

**错误数量**: 18处

#### 修复清单

| 行号 | Before | After | 状态 |
|------|--------|-------|------|
| 156 | `fmt.Errorf("workflow not found")` | `errorx.New(errno.ErrWorkflowNotFoundCode, ...)` | ⏳ |
| 234 | `fmt.Errorf("node execution failed")` | `errorx.WrapByCode(err, ErrWorkflowExecutionFailedCode, ...)` | ⏳ |
| 458 | `errors.New("invalid workflow")` | `errorx.New(errno.ErrWorkflowInvalidParamCode, ...)` | ⏳ |

#### 修复示例

```go
// 文件: backend/domain/workflow/service/service_impl.go:156

// ❌ Before
workflow, err := s.repo.Get(ctx, workflowID)
if err != nil {
    return nil, fmt.Errorf("workflow %s not found", workflowID)
}

// ✅ After
workflow, err := s.repo.Get(ctx, workflowID)
if err != nil {
    return nil, errorx.New(errno.ErrWorkflowNotFoundCode,
        errorx.KV("workflow_id", workflowID),
        errorx.KV("operation", "get_workflow"),
    )
}
```

#### 错误码映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| workflow not found | ErrWorkflowNotFoundCode | workflow_id, operation |
| invalid workflow | ErrWorkflowInvalidParamCode | workflow_id, reason |
| node execution failed | ErrWorkflowExecutionFailedCode | workflow_id, node_id, node_type |
| workflow execution failed | ErrWorkflowExecutionFailedCode | workflow_id, error |

### Application模块修复指南

**文件**: `backend/application/workflow/chatflow.go`

**错误数量**: 12处

#### 修复示例

```go
// 文件: backend/application/workflow/chatflow.go:530

// ❌ Before
if lastMessage.Role != "user" {
    return nil, errors.New("the role of the last day message must be user")
}

// ✅ After
if lastMessage.Role != "user" {
    return nil, errorx.New(errno.ErrConversationInvalidParamCode,
        errorx.KV("reason", "last message role must be user"),
        errorx.KV("conversation_id", conversationID),
        errorx.KV("actual_role", lastMessage.Role),
        errorx.KV("expected_role", "user"),
    )
}
```

### Permission模块修复指南

**文件**: `backend/domain/permission/service/permission_checker.go`

**错误数量**: 9处

#### 修复示例

```go
// ❌ Before
if !hasPermission {
    return errors.New("permission denied")
}

// ✅ After
if !hasPermission {
    return errorx.New(errno.ErrPermissionDeniedCode,
        errorx.KV("user_id", userID),
        errorx.KV("resource_type", resourceType),
        errorx.KV("resource_id", resourceID),
        errorx.KV("required_permission", requiredPerm),
        errorx.KV("reason", "insufficient privileges"),
    )
}
```

---

## 测试验证

### 单元测试

#### 测试错误码正确性

```go
func TestKnowledgeService_CreateDocument_Error(t *testing.T) {
    // 测试空文档错误
    _, err := service.CreateDocument(ctx, &CreateDocumentRequest{
        DocumentID: "",
        Content:    "",
    })

    // 验证错误码
    require.Error(t, err)
    statusErr, ok := err.(errorx.StatusError)
    require.True(t, ok, "错误必须实现StatusError接口")
    assert.Equal(t, errno.ErrKnowledgeInvalidParamCode, statusErr.Code())

    // 验证KV对
    extra := statusErr.Extra()
    assert.Contains(t, extra, "reason")
    assert.Contains(t, extra, "operation")
}
```

#### 测试上下文信息

```go
func TestWorkflowService_Execute_Error(t *testing.T) {
    // 测试节点执行失败
    _, err := service.Execute(ctx, workflowID)

    require.Error(t, err)
    statusErr := err.(errorx.StatusError)

    // 验证必需的上下文信息
    extra := statusErr.Extra()
    assert.Contains(t, extra, "workflow_id", "缺少workflow_id")
    assert.Contains(t, extra, "node_id", "缺少node_id")
    assert.Contains(t, extra, "node_type", "缺少node_type")
}
```

### 集成测试

```bash
# 运行所有测试
cd backend
go test ./... -cover -v

# 按模块测试
go test ./domain/knowledge/... -v
go test ./domain/workflow/... -v
go test ./application/... -v
```

### Lint检查

```bash
# 运行golangci-lint
golangci-lint run --enable-all

# 检查特定文件
golangci-lint run backend/domain/knowledge/service/knowledge.go
```

### 验证清单

修复完成后，确保：

- [ ] 所有错误使用统一错误码
- [ ] 所有错误包含必需的KV对
- [ ] 错误消息没有硬编码（在errno包的i18n中定义）
- [ ] 所有测试通过
- [ ] golangci-lint无警告
- [ ] 错误响应格式符合规范

---

## 常见问题

### Q1: 如何选择合适的错误码？

**A**: 参考以下规则：

1. **按模块选择**: Knowledge模块使用 `ErrKnowledge*`，Workflow使用 `ErrWorkflow*`
2. **按错误类型**:
   - 资源不存在 → `*NotFoundCode`
   - 参数无效 → `*InvalidParamCode`
   - 权限不足 → `*PermissionCode` 或 `ErrPermissionDeniedCode`
   - 数据库错误 → `*DBCode`
   - 执行失败 → `*FailedCode`

### Q2: KV对应该包含哪些信息？

**A**: 至少包含：

1. **资源标识**: resource_id, document_id, workflow_id等
2. **操作类型**: operation (get, save, delete等)
3. **错误原因**: reason
4. **用户/租户ID** (如适用): user_id, tenant_id

根据具体场景添加其他关键信息。

### Q3: 是否需要修改测试文件？

**A**: 不需要。测试文件中的 `errors.New()` 可以保留，但如果测试涉及错误码验证，建议也使用统一格式。

### Q4: 如何处理第三方库返回的错误？

**A**: 使用 `errorx.WrapByCode()` 包装：

```go
// 第三方库错误
err := thirdPartyClient.DoSomething()
if err != nil {
    return errorx.WrapByCode(err, errno.ErrKnowledgeSystemCode,
        errorx.KV("operation", "third_party_call"),
        errorx.KV("library", "third-party"),
        errorx.KV("error", err.Error()),
    )
}
```

### Q5: 错误消息应该放在哪里？

**A**: 错误消息模板定义在 `backend/pkg/errorx/code/` 和 `backend/types/errno/` 的 `init()` 函数中：

```go
// backend/types/errno/knowledge.go
func init() {
    code.Register(
        ErrKnowledgeInvalidParamCode,
        "invalid parameter : {msg}",
        code.WithAffectStability(false),
    )
}
```

业务代码中只需要传递KV对，不要硬编码错误消息。

### Q6: 如何处理错误国际化？

**A**: 错误码系统已支持中英文双语。错误消息在 `code.Register()` 中注册：

```go
code.Register(
    ErrKnowledgeNotFoundCode,
    "knowledge not exist: {msg}",  // 英文
    code.WithMessageZH("知识库不存在: {msg}"),  // 中文
)
```

响应时会根据请求头自动选择语言。

### Q7: 修复后如何验证没有破坏现有功能？

**A**:

1. **运行完整测试套件**: `go test ./... -cover`
2. **检查CI/CD**: 确保所有CI检查通过
3. **手动测试**: 关键业务流程测试
4. **灰度发布**: 先在部分用户中验证
5. **监控错误指标**: 观察错误码统计是否正常

---

## 附录

### A. 完整错误码列表

详见: `backend/types/errno/`

### B. 错误处理示例代码

详见: `backend/pkg/errorx/`

### C. 相关文档

- [统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)
- [全局一致性检查清单](./ZKER-全局一致性检查清单_v1.0.md)

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
**维护团队**: ZKER企业级功能完善团队
