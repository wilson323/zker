# 错误码映射表

**版本**: v1.0
**日期**: 2025-01-01
**维护团队**: ZKER企业级功能完善团队

---

## 📖 目录

1. [Knowledge模块](#knowledge模块)
2. [Workflow模块](#workflow模块)
3. [Conversation模块](#conversation模块)
4. [Bot模块](#bot模块)
5. [Permission模块](#permission模块)
6. [App模块](#app模块)
7. [Tenant模块](#tenant模块)
8. [通用错误码](#通用错误码)

---

## Knowledge模块

**错误码范围**: 105000000 ~ 105999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 105000000 | ErrKnowledgeInvalidParamCode | 参数无效 | 400 | 知识库参数无效 |
| 105000001 | ErrKnowledgePermissionCode | 权限不足 | 403 | 知识库权限不足 |
| 105000002 | ErrKnowledgeNonRetryableCode | 不可重试 | 400 | 不可重试的错误 |
| 105000003 | ErrKnowledgeDBCode | 数据库错误 | 500 | MySQL操作失败 |
| 105000004 | ErrKnowledgeSearchStoreCode | 搜索存储错误 | 500 | SearchStore操作失败 |
| 105000005 | ErrKnowledgeSystemCode | 系统错误 | 500 | 系统内部错误 |
| 105000006 | ErrKnowledgeCrossDomainCode | 跨域错误 | 500 | 跨域调用错误 |
| 105000007 | ErrKnowledgeEmbeddingCode | 向量化错误 | 500 | 向量化操作失败 |
| 105000008 | ErrKnowledgeIDGenCode | ID生成失败 | 500 | ID生成失败 |
| 105000009 | ErrKnowledgeMQSendFailCode | MQ发送失败 | 500 | MQ发送消息失败 |
| 105000010 | ErrKnowledgeDuplicateCode | 重复错误 | 409 | 知识库名称重复 |
| 105000011 | ErrKnowledgeNotExistCode | 资源不存在 | 404 | 知识库不存在 |
| 105000012 | ErrKnowledgeDocumentNotExistCode | 资源不存在 | 404 | 文档不存在 |
| 105000013 | ErrKnowledgeSemanticColumnValueEmptyCode | 参数无效 | 400 | 语义列值为空 |
| 105000014 | ErrKnowledgeParseJSONCode | 解析失败 | 500 | JSON解析失败 |
| 105000015 | ErrKnowledgeResegmentNotSupportedCode | 不支持的操作 | 400 | 处理中，不支持重新分段 |
| 105000016 | ErrKnowledgeFieldNameDuplicatedCode | 重复错误 | 409 | 字段名重复 |
| 105000017 | ErrKnowledgeDocOversizeCode | 参数无效 | 400 | 文档过大 |
| 105000018 | ErrKnowledgeDocNotReadyCode | 状态错误 | 400 | 文档未就绪 |
| 105000019 | ErrKnowledgeDownloadFailedCode | 下载失败 | 500 | 下载失败 |
| 105000020 | ErrKnowledgeTableInfoNotExistCode | 资源不存在 | 404 | 表信息不存在 |
| 105000021 | ErrKnowledgePutObjectFailCode | 存储失败 | 500 | PutObject失败 |
| 105000022 | ErrKnowledgeGetObjectURLFailCode | 获取URL失败 | 500 | 获取对象URL失败 |
| 105000023 | ErrKnowledgeGetDocProgressFailCode | 查询失败 | 500 | 获取文档进度失败 |
| 105000024 | ErrKnowledgeSliceInsertPositionIllegalCode | 参数无效 | 400 | 分片插入位置非法 |
| 105000025 | ErrKnowledgeSliceNotExistCode | 资源不存在 | 404 | 分片不存在 |
| 105000026 | ErrKnowledgeColumnParseFailCode | 解析失败 | 500 | 列解析失败 |
| 105000027 | ErrKnowledgeAutoAnnotationNotSupportedCode | 不支持的操作 | 400 | 不支持自动标注 |
| 105000028 | ErrKnowledgeGetParserFailCode | 获取失败 | 500 | 获取解析器失败 |
| 105000029 | ErrKnowledgeGetObjectFailCode | 获取失败 | 500 | 获取对象失败 |
| 105000030 | ErrKnowledgeParserParseFailCode | 解析失败 | 500 | 解析器解析失败 |
| 105000031 | ErrKnowledgeBuildRetrieveChainFailCode | 构建失败 | 500 | 构建检索链失败 |
| 105000032 | ErrKnowledgeRetrieveExecFailCode | 执行失败 | 500 | 检索执行失败 |
| 105000033 | ErrKnowledgeNL2SqlExecFailCode | 执行失败 | 500 | NL2SQL执行失败 |
| 105000034 | ErrKnowledgeCopyFailCode | 复制失败 | 500 | 复制失败 |
| 105000035 | ErrKnowledgeParseResultEmptyCode | 结果为空 | 400 | 解析结果为空 |
| 105000036 | ErrKnowledgeCacheClientSetFailCode | 缓存失败 | 500 | 缓存客户端设置失败 |
| 105000037 | ErrKnowledgeCheckTableSliceValidCode | 验证失败 | 400 | 分片内容验证失败 |

### 错误消息映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| document is empty | ErrKnowledgeInvalidParamCode | document_id, operation |
| document not found | ErrKnowledgeDocumentNotExistCode | document_id, operation |
| document not exist | ErrKnowledgeDocumentNotExistCode | document_id, operation |
| slice not found | ErrKnowledgeSliceNotExistCode | slice_id, operation |
| slice not exist | ErrKnowledgeSliceNotExistCode | slice_id, operation |
| knowledge not found | ErrKnowledgeNotExistCode | knowledge_id, operation |
| knowledge not exist | ErrKnowledgeNotExistCode | knowledge_id, operation |
| invalid document type | ErrKnowledgeInvalidParamCode | document_id, type, expected |
| document oversize | ErrKnowledgeDocOversizeCode | document_id, size, max_size |
| parse failed | ErrKnowledgeParserParseFailCode | document_id, parser_type, error |
| retrieve failed | ErrKnowledgeRetrieveExecFailCode | knowledge_id, query, error |
| copy failed | ErrKnowledgeCopyFailCode | source_id, target_id, error |

### 使用示例

```go
// 文档不存在
return errorx.New(errno.ErrKnowledgeDocumentNotExistCode,
    errorx.KV("document_id", docID),
    errorx.KV("operation", "get_document"),
)

// 文档为空
return errorx.New(errno.ErrKnowledgeInvalidParamCode,
    errorx.KV("reason", "document is empty"),
    errorx.KV("document_id", docID),
    errorx.KV("operation", "save_document"),
)

// 解析失败（包装错误）
return errorx.WrapByCode(err, errno.ErrKnowledgeParserParseFailCode,
    errorx.KV("document_id", docID),
    errorx.KV("parser_type", "pdf"),
    errorx.KV("error", err.Error()),
)
```

---

## Workflow模块

**错误码范围**: 203000000 ~ 203999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 203000000 | ErrWorkflowInvalidParamCode | 参数无效 | 400 | 工作流参数无效 |
| 203000001 | ErrWorkflowNotFoundCode | 资源不存在 | 404 | 工作流不存在 |
| 203000002 | ErrWorkflowAlreadyExistsCode | 重复错误 | 409 | 工作流已存在 |
| 203000050 | ErrWorkflowExecutionFailedCode | 执行失败 | 500 | 工作流执行失败 |
| 203000051 | ErrWorkflowNodeExecutionFailedCode | 执行失败 | 500 | 节点执行失败 |
| 203000052 | ErrWorkflowInvalidStateCode | 状态错误 | 400 | 工作流状态无效 |

### 错误消息映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| workflow not found | ErrWorkflowNotFoundCode | workflow_id, operation |
| workflow not exist | ErrWorkflowNotFoundCode | workflow_id, operation |
| invalid workflow | ErrWorkflowInvalidParamCode | workflow_id, reason |
| workflow execution failed | ErrWorkflowExecutionFailedCode | workflow_id, error |
| node execution failed | ErrWorkflowNodeExecutionFailedCode | workflow_id, node_id, node_type, error |
| workflow already exists | ErrWorkflowAlreadyExistsCode | workflow_id, name |

### 使用示例

```go
// 工作流不存在
return errorx.New(errno.ErrWorkflowNotFoundCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("operation", "get_workflow"),
)

// 节点执行失败
return errorx.WrapByCode(err, errno.ErrWorkflowNodeExecutionFailedCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("node_id", nodeID),
    errorx.KV("node_type", "llm"),
    errorx.KV("error", err.Error()),
)
```

---

## Conversation模块

**错误码范围**: 202000000 ~ 202999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 202000000 | ErrConversationInvalidParamCode | 参数无效 | 400 | 对话参数无效 |
| 202000001 | ErrConversationNotFoundCode | 资源不存在 | 404 | 对话不存在 |
| 202000002 | ErrConversationMessageFailedCode | 消息失败 | 500 | 消息操作失败 |

### 错误消息映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| conversation not found | ErrConversationNotFoundCode | conversation_id, operation |
| last message role must be user | ErrConversationInvalidParamCode | conversation_id, actual_role, expected_role |
| create message failed | ErrConversationMessageFailedCode | conversation_id, error |

### 使用示例

```go
// 对话不存在
return errorx.New(errno.ErrConversationNotFoundCode,
    errorx.KV("conversation_id", conversationID),
    errorx.KV("operation", "get_conversation"),
)

// 最后消息角色错误
return errorx.New(errno.ErrConversationInvalidParamCode,
    errorx.KV("reason", "last message role must be user"),
    errorx.KV("conversation_id", conversationID),
    errorx.KV("actual_role", lastMessage.Role),
    errorx.KV("expected_role", "user"),
)
```

---

## Bot模块

**错误码范围**: 201000000 ~ 201999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 201000000 | ErrBotInvalidParamCode | 参数无效 | 400 | Bot参数无效 |
| 201000001 | ErrBotNotFoundCode | 资源不存在 | 404 | Bot不存在 |
| 201000002 | ErrBotInvalidConfigCode | 配置无效 | 400 | Bot配置无效 |

### 错误消息映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| bot not found | ErrBotNotFoundCode | bot_id, operation |
| invalid bot config | ErrBotInvalidConfigCode | bot_id, config_key, reason |

### 使用示例

```go
// Bot不存在
return errorx.New(errno.ErrBotNotFoundCode,
    errorx.KV("bot_id", botID),
    errorx.KV("operation", "get_bot"),
)
```

---

## Permission模块

**错误码范围**: 108000000 ~ 108999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 108000040 | ErrPermissionDeniedCode | 权限不足 | 403 | 权限不足 |
| 108000041 | ErrPermissionRoleNotFoundCode | 资源不存在 | 404 | 角色不存在 |
| 108000042 | ErrPermissionInvalidParamCode | 参数无效 | 400 | 权限参数无效 |

### 错误消息映射

| 错误消息 | 错误码 | 必需KV对 |
|---------|--------|----------|
| permission denied | ErrPermissionDeniedCode | user_id, resource_type, resource_id, required_permission, reason |
| unauthorized access | ErrPermissionDeniedCode | user_id, resource, reason |
| role not found | ErrPermissionRoleNotFoundCode | role_id, operation |

### 使用示例

```go
// 权限不足
return errorx.New(errno.ErrPermissionDeniedCode,
    errorx.KV("user_id", userID),
    errorx.KV("resource_type", "knowledge"),
    errorx.KV("resource_id", knowledgeID),
    errorx.KV("required_permission", "write"),
    errorx.KV("reason", "insufficient privileges"),
)
```

---

## App模块

**错误码范围**: 204000000 ~ 204999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 204000000 | ErrAppInvalidParamCode | 参数无效 | 400 | 应用参数无效 |
| 204000001 | ErrAppNotFoundCode | 资源不存在 | 404 | 应用不存在 |
| 204000002 | ErrAppInvalidConfigCode | 配置无效 | 400 | 应用配置无效 |

### 使用示例

```go
// 应用不存在
return errorx.New(errno.ErrAppNotFoundCode,
    errorx.KV("app_id", appID),
    errorx.KV("operation", "get_app"),
)
```

---

## Tenant模块

**错误码范围**: 200000000 ~ 200999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 200000001 | ErrTenantNotFoundCode | 资源不存在 | 404 | 租户不存在 |
| 200000010 | ErrTenantAlreadyExistsCode | 重复错误 | 409 | 租户已存在 |
| 200000020 | ErrTenantInvalidParamCode | 参数无效 | 400 | 租户参数无效 |
| 200000040 | ErrTenantSuspendedCode | 状态错误 | 403 | 租户已暂停 |
| 200000041 | ErrTenantDeletedCode | 状态错误 | 410 | 租户已删除 |

### 使用示例

```go
// 租户不存在
return errorx.New(errno.ErrTenantNotFoundCode,
    errorx.KV("tenant_id", tenantID),
    errorx.KV("operation", "get_tenant"),
)
```

---

## 通用错误码

**错误码范围**: 100000000 ~ 199999999

### 错误码列表

| 错误码 | 常量名 | 错误类型 | HTTP状态码 | 说明 |
|--------|--------|---------|-----------|------|
| 201000001 | InvalidParamsCode | 参数无效 | 400 | 参数错误 |
| 201000002 | MissingParamCode | 缺少参数 | 400 | 缺少必要参数 |
| 201000003 | InvalidFormatCode | 格式错误 | 400 | 格式错误 |
| 401000001 | UnauthorizedCode | 未认证 | 401 | 未认证 |
| 401000002 | TokenExpiredCode | Token过期 | 401 | Token已过期 |
| 401000003 | InvalidTokenCode | Token无效 | 401 | 无效的Token |
| 403000001 | ForbiddenCode | 权限不足 | 403 | 权限不足 |
| 404000001 | NotFoundCode | 资源不存在 | 404 | 资源不存在 |
| 500000001 | InternalErrorCode | 内部错误 | 500 | 服务器内部错误 |
| 500000002 | DatabaseErrorCode | 数据库错误 | 500 | 数据库错误 |
| 500000003 | IDGenErrorCode | ID生成失败 | 500 | ID生成失败 |
| 500000004 | CacheErrorCode | 缓存错误 | 500 | 缓存操作失败 |

### 使用示例

```go
// 参数无效
return errorx.New(errno.InvalidParamsCode,
    errorx.KV("param", "limit"),
    errorx.KV("value", limit),
    errorx.KV("expected", "1-100"),
)

// 资源不存在
return errorx.New(errno.NotFoundCode,
    errorx.KV("resource_type", "user"),
    errorx.KV("resource_id", userID),
)

// 内部错误
return errorx.WrapByCode(err, errno.InternalErrorCode,
    errorx.KV("operation", "save_data"),
    errorx.KV("error", err.Error()),
)
```

---

## 📚 附录

### A. 错误码命名规范

```
Err{模块名}{错误类型}Code

示例:
- ErrKnowledgeInvalidParamCode  (Knowledge模块 + 参数无效)
- ErrWorkflowNotFoundCode       (Workflow模块 + 资源不存在)
- ErrPermissionDeniedCode       (Permission模块 + 权限不足)
```

### B. 错误码分配规则

```
{模块段}{类型段}{序号}

示例: 105000011
- 105: Knowledge模块
- 000: 固定
- 011: 序号（文档不存在）
```

### C. HTTP状态码映射

| 错误类型 | HTTP状态码 |
|---------|-----------|
| 资源不存在 (001-009) | 404 |
| 重复错误 (010-019) | 409 |
| 参数无效 (020-039) | 400 |
| 权限不足 (040-049) | 403 |
| 执行失败 (050-089) | 500 |
| 其他 (090-099) | 400 |

### D. 相关文档

- [统一错误码定义规范](../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [错误处理修复指南](../docs/企业级功能完善与统一性设计方案/ZKER-错误处理修复指南_v1.0.md)
- [企业级开发规范手册](../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
**维护团队**: ZKER企业级功能完善团队
