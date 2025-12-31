# ZKER 错误处理修复清单

**版本**: v1.0
**日期**: 2025-01-01
**作者**: ZKER企业级功能完善团队

---

## 📊 总体进度

- **总违规数**: 189 处
- **已修复**: 0 处 (0%)
- **待修复**: 189 处 (100%)

### 优先级分布

| 优先级 | 数量 | 已修复 | 待修复 | 完成率 |
|--------|------|--------|--------|--------|
| P0 - Domain Service层 | 80 | 0 | 80 | 0% |
| P1 - Application层 | 60 | 0 | 60 | 0% |
| P2 - Infra层和其他 | 49 | 0 | 49 | 0% |

### 模块分布

| 模块 | 数量 | 已修复 | 待修复 | 完成率 |
|------|------|--------|--------|--------|
| Knowledge | 35 | 0 | 35 | 0% |
| Workflow | 28 | 0 | 28 | 0% |
| Conversation | 18 | 0 | 18 | 0% |
| Permission | 12 | 0 | 12 | 0% |
| Bot | 15 | 0 | 15 | 0% |
| Tenant | 10 | 0 | 10 | 0% |
| Plugin | 12 | 0 | 12 | 0% |
| 其他 | 59 | 0 | 59 | 0% |

---

## 🎯 P0优先级 - Domain Service层

### Knowledge模块 (35处)

#### 文件: `backend/domain/knowledge/service/knowledge.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 322 | errors.New | `errors.New("document is empty")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |
| 458 | errors.New | `errors.New("invalid document type")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |
| 538 | errors.New | `errors.New("knowledge not found")` | `errorx.New(errno.ErrKnowledgeNotExistCode, ...)` | ⏳ | - |
| 612 | fmt.Errorf | `fmt.Errorf("parse failed: %w", err)` | `errorx.WrapByCode(err, errno.ErrKnowledgeParserParseFailCode, ...)` | ⏳ | - |
| 734 | errors.New | `errors.New("slice not found")` | `errorx.New(errno.ErrKnowledgeSliceNotExistCode, ...)` | ⏳ | - |

**小计**: 5处 | 已完成: 0 | 待修复: 5

#### 文件: `backend/domain/knowledge/service/datacopy.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 245 | errors.New | `errors.New("copy failed")` | `errorx.New(errno.ErrKnowledgeCopyFailCode, ...)` | ⏳ | - |
| 312 | fmt.Errorf | `fmt.Errorf("unsupported scene: %s", scene)` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |
| 456 | errors.New | `errors.New("source knowledge not found")` | `errorx.New(errno.ErrKnowledgeNotExistCode, ...)` | ⏳ | - |

**小计**: 3处 | 已完成: 0 | 待修复: 3

### Workflow模块 (28处)

#### 文件: `backend/domain/workflow/service/service_impl.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 156 | fmt.Errorf | `fmt.Errorf("workflow not found")` | `errorx.New(errno.ErrWorkflowNotFoundCode, ...)` | ⏳ | - |
| 234 | fmt.Errorf | `fmt.Errorf("node execution failed: %w", err)` | `errorx.WrapByCode(err, errno.ErrWorkflowExecutionFailedCode, ...)` | ⏳ | - |
| 458 | errors.New | `errors.New("invalid workflow")` | `errorx.New(errno.ErrWorkflowInvalidParamCode, ...)` | ⏳ | - |
| 512 | fmt.Errorf | `fmt.Errorf("workflow %s already exists", id)` | `errorx.New(errno.ErrWorkflowAlreadyExistsCode, ...)` | ⏳ | - |

**小计**: 4处 | 已完成: 0 | 待修复: 4

#### 文件: `backend/domain/workflow/service/executable_impl.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 89 | errors.New | `errors.New("workflow not executable")` | `errorx.New(errno.ErrWorkflowInvalidParamCode, ...)` | ⏳ | - |
| 123 | fmt.Errorf | `fmt.Errorf("execution failed: %w", err)` | `errorx.WrapByCode(err, errno.ErrWorkflowExecutionFailedCode, ...)` | ⏳ | - |

**小计**: 2处 | 已完成: 0 | 待修复: 2

### Conversation模块 (18处)

#### 文件: `backend/domain/conversation/service/conversation_impl.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 178 | errors.New | `errors.New("conversation not found")` | `errorx.New(errno.ErrConversationNotFoundCode, ...)` | ⏳ | - |
| 234 | fmt.Errorf | `fmt.Errorf("create message failed: %w", err)` | `errorx.WrapByCode(err, errno.ErrConversationMessageFailedCode, ...)` | ⏳ | - |
| 289 | errors.New | `errors.New("invalid conversation status")` | `errorx.New(errno.ErrConversationInvalidParamCode, ...)` | ⏳ | - |

**小计**: 3处 | 已完成: 0 | 待修复: 3

### Permission模块 (12处)

#### 文件: `backend/domain/permission/service/permission_checker.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 89 | errors.New | `errors.New("permission denied")` | `errorx.New(errno.ErrPermissionDeniedCode, ...)` | ⏳ | - |
| 145 | fmt.Errorf | `fmt.Errorf("user %d has no access", uid)` | `errorx.New(errno.ErrPermissionDeniedCode, ...)` | ⏳ | - |
| 201 | errors.New | `errors.New("role not found")` | `errorx.New(errno.ErrPermissionRoleNotFoundCode, ...)` | ⏳ | - |

**小计**: 3处 | 已完成: 0 | 待修复: 3

### Bot模块 (15处)

#### 文件: `backend/domain/bot/service/bot_service.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 123 | errors.New | `errors.New("bot not found")` | `errorx.New(errno.ErrBotNotFoundCode, ...)` | ⏳ | - |
| 189 | fmt.Errorf | `fmt.Errorf("invalid bot config: %w", err)` | `errorx.WrapByCode(err, errno.ErrBotInvalidConfigCode, ...)` | ⏳ | - |

**小计**: 2处 | 已完成: 0 | 待修复: 2

---

## 🔥 P1优先级 - Application层

### Workflow模块 (12处)

#### 文件: `backend/application/workflow/chatflow.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 530 | errors.New | `errors.New("the role of the last day message must be user")` | `errorx.New(errno.ErrConversationInvalidParamCode, ...)` | ⏳ | - |
| 581 | errors.New | `errors.New("project_id and bot_id cannot be set at the same time")` | `errorx.New(errno.ErrConversationInvalidParamCode, ...)` | ⏳ | - |
| 645 | fmt.Errorf | `fmt.Errorf("workflow execution failed: %w", err)` | `errorx.WrapByCode(err, errno.ErrWorkflowExecutionFailedCode, ...)` | ⏳ | - |

**小计**: 3处 | 已完成: 0 | 待修复: 3

### Knowledge模块 (8处)

#### 文件: `backend/application/knowledge/knowledge.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 64 | errors.New | `errors.New("unknown document type")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |
| 340 | errors.New | `errors.New("knowledge not found")` | `errorx.New(errno.ErrKnowledgeNotExistCode, ...)` | ⏳ | - |
| 350 | errors.New | `errors.New("document base is empty")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |
| 450 | errors.New | `errors.New("document ids is empty")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |

**小计**: 4处 | 已完成: 0 | 待修复: 4

### App模块 (6处)

#### 文件: `backend/application/app/app.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 123 | errors.New | `errors.New("app not found")` | `errorx.New(errno.ErrAppNotFoundCode, ...)` | ⏳ | - |
| 234 | fmt.Errorf | `fmt.Errorf("invalid app config: %w", err)` | `errorx.WrapByCode(err, errno.ErrAppInvalidConfigCode, ...)` | ⏳ | - |

**小计**: 2处 | 已完成: 0 | 待修复: 2

---

## ⚙️ P2优先级 - Infra层和其他

### Infra层 (20处)

#### 文件: `backend/infra/cache/redis_cluster.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 89 | errors.New | `errors.New("redis connection failed")` | `errorx.New(errno.ErrCacheErrorCode, ...)` | ⏳ | - |
| 123 | fmt.Errorf | `fmt.Errorf("cache set failed: %w", err)` | `errorx.WrapByCode(err, errno.ErrCacheErrorCode, ...)` | ⏳ | - |

**小计**: 2处 | 已完成: 0 | 待修复: 2

### Crossdomain层 (15处)

#### 文件: `backend/crossdomain/knowledge/impl/knowledge.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 76 | errors.New | `errors.New("document parsing strategy is required")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ⏳ | - |

**小计**: 1处 | 已完成: 0 | 待修复: 1

### API层 (14处)

#### 文件: `backend/api/handler/coze/workflow_service.go`

| 行号 | 类型 | Before | After | 状态 | 负责人 |
|------|------|--------|-------|------|--------|
| 234 | errors.New | `errors.New("invalid request")` | `errorx.New(errno.InvalidParamsCode, ...)` | ⏳ | - |
| 345 | fmt.Errorf | `fmt.Errorf("handler failed: %w", err)` | `errorx.WrapByCode(err, errno.InternalErrorCode, ...)` | ⏳ | - |

**小计**: 2处 | 已完成: 0 | 待修复: 2

---

## 📈 修复进度统计

### 按优先级

```
P0: ████████████████████░░░░░░░░ 0/80   (0%)
P1: ████████████████████░░░░░░░░ 0/60   (0%)
P2: ████████████████████░░░░░░░░ 0/49   (0%)
```

### 按模块

```
Knowledge:    ████████████████████░░░░░ 0/35   (0%)
Workflow:     ████████████████████░░░░░ 0/28   (0%)
Conversation: ████████████████████░░░░░ 0/18   (0%)
Permission:   ████████████████████░░░░░ 0/12   (0%)
Bot:          ████████████████████░░░░░ 0/15   (0%)
Tenant:       ████████████████████░░░░░ 0/10   (0%)
Plugin:       ████████████████████░░░░░ 0/12   (0%)
Other:        ████████████████████░░░░░ 0/59   (0%)
```

### 按错误类型

```
errors.New():  ████████████████████░░░░░ 0/124  (0%)
fmt.Errorf():  ████████████████████░░░░░ 0/58   (0%)
```

---

## 🔧 修复说明

### 状态标记

- ⏳ **待修复**: 已识别，等待修复
- 🔄 **修复中**: 正在修复
- ✅ **已完成**: 已修复并验证
- ❌ **有问题**: 修复遇到问题，需要review
- ⏭️ **跳过**: 已决定跳过（需说明原因）

### 修复流程

1. **识别**: 使用扫描工具识别违规
2. **分析**: 确定合适的错误码和KV对
3. **修复**: 修改代码
4. **测试**: 运行单元测试和集成测试
5. **验证**: 确保所有检查通过
6. **提交**: 更新清单并提交PR

### 更新清单

修复完成后，请更新此清单：

```markdown
| 322 | errors.New | `errors.New("document is empty")` | `errorx.New(errno.ErrKnowledgeInvalidParamCode, ...)` | ✅ | 张三 |
```

并更新总体进度统计。

---

## 📝 备注

### 高频错误模式

1. **资源不存在** (45处): 使用 `*NotFoundCode`
2. **参数无效** (38处): 使用 `*InvalidParamCode`
3. **执行失败** (28处): 使用 `*FailedCode` 或 `*DBCode`
4. **权限不足** (12处): 使用 `ErrPermissionDeniedCode`

### 修复建议

1. **批量修复**: 同一类型的错误可以批量修复
2. **优先处理**: 先修复P0优先级，影响最大
3. **模块并行**: 不同模块可以并行修复
4. **持续验证**: 每修复一批就验证一次

### 注意事项

1. **不要修改测试文件**: 测试代码中的 errors.New 可以保留
2. **保留上下文**: 确保KV对包含关键信息
3. **选择正确错误码**: 参考错误码规范文档
4. **运行测试**: 修复后必须运行测试

---

## 📚 相关文档

- [错误处理修复指南](./ZKER-错误处理修复指南_v1.0.md)
- [统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [企业级开发规范手册](./ZKER-企业级开发规范手册_v1.0.md)

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
**维护团队**: ZKER企业级功能完善团队
