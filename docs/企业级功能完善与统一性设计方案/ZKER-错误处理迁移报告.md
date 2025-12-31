# ZKER 错误处理迁移报告

**生成时间**: 2025-12-31 08:56:28
**扫描目录**: D:\code\coze-studio\backend

---

## 📊 总体统计

- **总计违规**: 875 处
- **errors.New()**: 33 处
- **fmt.Errorf()**: 842 处

## 🎯 优先级分布

| 优先级 | 数量 | 说明 |
|--------|------|------|
| P0 | 0 | Domain Service层（核心业务逻辑）|
| P1 | 0 | Application层 + Domain其他层|
| P2 | 875 | Infra层和其他|

## 📦 模块分布

- **common**: 875 处

## 🔝 高频违规文件（Top 20）

| 文件 | 违规次数 |
|------|---------|
| api\handler\coze\validator.go | 27 |
| domain\tenant\service\tenant_registration_service.go | 17 |
| domain\workflow\service\service_impl.go | 16 |
| domain\billing\service\budget_management_service.go | 15 |
| domain\workflow\internal\nodes\selector\operator.go | 15 |
| domain\workflow\service\executable_impl.go | 14 |
| domain\monitoring\service\alert_service.go | 14 |
| application\queue\tasks.go | 13 |
| domain\humaninloop\service\collaboration_orchestrator.go | 13 |
| infra\queue\dead_letter_queue.go | 13 |
| tests\e2e\cmd\verify\main.go | 13 |
| domain\tenant\migration\cutover.go | 12 |
| infra\document\searchstore\impl\oceanbase\consts.go | 12 |
| infra\saga\orchestrator.go | 12 |
| domain\billing\service\payment_service.go | 11 |
| domain\tenant\migration\tenant_id_migration.go | 11 |
| infra\cache\permission_cache_unified.go | 11 |
| infra\document\searchstore\impl\milvus\milvus_manager.go | 11 |
| infra\document\searchstore\impl\oceanbase\convert.go | 11 |
| infra\eventbus\impl\nats\consumer.go | 11 |

---

## 📋 详细违规清单（按模块分类）

### COMMON 模块 (875 处)

| 文件 | 行号 | 类型 | 错误消息 | 建议错误码 | 优先级 |
|------|------|------|---------|-----------|--------|
| api\handler\coze\agent_run_service.go | 195 | errors.New | parameters field should be an object, no... | InternalErrorCode | P2 |
| api\handler\coze\message_service.go | 124 | errors.New | answer message id is required | InternalErrorCode | P2 |
| api\middleware\row_level_security.go | 422 | errors.New | resources must be a slice | InternalErrorCode | P2 |
| application\upload\icon.go | 192 | errors.New | tos key not exist | InternalErrorCode | P2 |
| application\upload\icon.go | 201 | errors.New | tos part is null | InternalErrorCode | P2 |
| application\upload\icon.go | 205 | errors.New | check parts fail | InternalErrorCode | P2 |
| application\upload\icon.go | 212 | errors.New | crc32 check fail | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_item.go | 68 | errors.New | bot already published | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_item.go | 71 | errors.New | bot is pending review | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 45 | errors.New | rating must be between 1 and 5 | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 53 | errors.New | no permission to edit this review | ForbiddenCode | P2 |
| domain\botstore\entity\bot_store_review.go | 68 | errors.New | item ID is required | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 71 | errors.New | rating must be between 1 and 5 | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 74 | errors.New | comment too long, maximum 1000 character... | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 88 | errors.New | rating must be between 1 and 5 | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 91 | errors.New | comment too long, maximum 1000 character... | InternalErrorCode | P2 |
| domain\botstore\entity\bot_store_review.go | 107 | errors.New | item ID is required | InternalErrorCode | P2 |
| domain\developer\entity\api_key.go | 69 | errors.New | cannot scan into APIKeyScopes | InternalErrorCode | P2 |
| domain\developer\entity\webhook.go | 74 | errors.New | cannot scan into WebhookEvents | InternalErrorCode | P2 |
| domain\developer\service\webhook_service.go | 175 | errors.New | webhook not found | NotFoundCode | P2 |
| domain\knowledge\service\datacopy.go | 538 | errors.New | invalid request | InvalidParamsCode | P2 |
| domain\knowledge\service\datacopy.go | 545 | errors.New | knowledge not found | NotFoundCode | P2 |
| domain\knowledge\service\knowledge.go | 324 | errors.New | document is empty | InternalErrorCode | P2 |
| domain\knowledge\service\knowledge.go | 329 | errors.New | auto caption type is not supported | InternalErrorCode | P2 |
| domain\workflow\internal\execute\callback.go | 754 | errors.New | state is nil | InternalErrorCode | P2 |
| domain\workflow\internal\execute\context.go | 96 | errors.New | state is nil | InternalErrorCode | P2 |
| domain\workflow\internal\execute\context.go | 137 | errors.New | state is nil | InternalErrorCode | P2 |
| domain\workflow\internal\execute\context.go | 187 | errors.New | state is nil | InternalErrorCode | P2 |
| domain\workflow\internal\execute\context.go | 246 | errors.New | state is nil | InternalErrorCode | P2 |
| domain\workflow\internal\execute\context.go | 300 | errors.New | state is nil | InternalErrorCode | P2 |
| domain\workflow\service\executable_impl.go | 323 | errors.New | CONVERSATION_NAME must be string | InternalErrorCode | P2 |
| domain\workflow\service\service_impl.go | 1658 | errors.New | conversation-related nodes are not suppo... | InternalErrorCode | P2 |
| pkg\security\xss.go | 311 | errors.New | invalid URL | InvalidParamsCode | P2 |
| main.go | 190 | fmt.Errorf | load env file(%s) failed, err=%w | InternalErrorCode | P2 |
| api\docs\generate_swagger.go | 615 | fmt.Errorf | error reading file %s: %w | InternalErrorCode | P2 |
| api\docs\generate_swagger.go | 722 | fmt.Errorf | error marshaling to YAML: %w | InternalErrorCode | P2 |
| api\docs\generate_swagger.go | 727 | fmt.Errorf | error marshaling to JSON: %w | InternalErrorCode | P2 |
| api\docs\generate_swagger.go | 732 | fmt.Errorf | error writing file: %w | InternalErrorCode | P2 |
| api\docs\generate_swagger.go | 742 | fmt.Errorf | missing openapi version | MissingParamCode | P2 |
| api\docs\generate_swagger.go | 746 | fmt.Errorf | missing API title | MissingParamCode | P2 |
| api\docs\generate_swagger.go | 750 | fmt.Errorf | missing API version | MissingParamCode | P2 |
| api\docs\generate_swagger.go | 755 | fmt.Errorf | no API paths defined | InternalErrorCode | P2 |
| api\docs\generate_swagger.go | 762 | fmt.Errorf | path %s has no operations | InternalErrorCode | P2 |
| api\handler\coze\validator.go | 29 | fmt.Errorf | tenant_id is required | InternalErrorCode | P2 |
| api\handler\coze\validator.go | 32 | fmt.Errorf | tenant_id length must be between 1 and 3... | InternalErrorCode | P2 |
| api\handler\coze\validator.go | 37 | fmt.Errorf | tenant_id contains invalid characters | InvalidParamsCode | P2 |
| api\handler\coze\validator.go | 45 | fmt.Errorf | uuid is required | InternalErrorCode | P2 |
| api\handler\coze\validator.go | 49 | fmt.Errorf | invalid UUID format | InvalidParamsCode | P2 |
| api\handler\coze\validator.go | 126 | fmt.Errorf | email is required | InternalErrorCode | P2 |
| api\handler\coze\validator.go | 129 | fmt.Errorf | email length exceeds maximum 254 charact... | InternalErrorCode | P2 |

---

## 🛠️ 修复指南

### 1. 快速开始

```bash
# 1. 查看完整报告
cat tools/error_handler_migration_report.md

# 2. 按优先级修复
# P0: Domain Service层（最优先）
# P1: Application层 + Domain其他层
# P2: Infra层和其他

# 3. 修复后验证
go test ./... -cover
golangci-lint run
```

### 2. 修复模板

#### errors.New() 修复模板

**Before:**
```go
return errors.New("document not found")
```

**After:**
```go
return errorx.New(errno.ErrKnowledgeDocumentNotExistCode,
    errorx.KV("document_id", documentID),
    errorx.KV("reason", "document not found"),
)
```

#### fmt.Errorf() 修复模板

**Before:**
```go
return fmt.Errorf("workflow %s not found", workflowID)
```

**After:**
```go
return errorx.New(errno.ErrWorkflowNotFoundCode,
    errorx.KV("workflow_id", workflowID),
    errorx.KV("reason", "workflow not found"),
)
```

### 3. 错误码映射表

详细的错误码映射请参考:
- `backend/types/errno/*.go` - 各模块错误码定义
- `docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md`

### 4. 注意事项

1. **不要修改测试文件**: 测试代码中的 errors.New 可以保留
2. **保留上下文信息**: 使用 errorx.KV() 添加关键上下文（ID、操作类型等）
3. **选择正确的错误码**: 参考错误码规范文档
4. **错误消息不要硬编码**: 将错误消息放在 errno 包的 i18n 中
5. **提交前验证**: 运行测试和lint确保修复正确

---

**报告生成完成**: 2025-12-31 08:56:28
