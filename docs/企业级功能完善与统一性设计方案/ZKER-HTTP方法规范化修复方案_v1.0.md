# ZKER HTTP方法规范化修复方案 v1.0

**📅 方案日期**: 2025-01-01
**🎯 修复目标**: 修复TOP 10高频HTTP方法误用,提升API规范评分至95分
**👨‍💻 执行专家**: API规范化专家 (Claude Code)
**📊 修复范围**: 10个高频违规API

---

## 📋 执行摘要

### 修复原则

由于`backend/api/router/coze/api.go`是**代码生成文件**,采用以下策略:

1. ✅ **向后兼容**: 保留旧API,新增符合RESTful规范的版本
2. ✅ **渐进迁移**: 逐步将前端调用迁移到新API
3. ✅ **双版本并存**: 旧版本标记为`@Deprecated`,新版本作为推荐
4. ✅ **无缝切换**: 6个月后移除旧版本

### 修复清单TOP 10

| # | 旧API | 新API | Handler | 优先级 |
|---|-------|-------|---------|--------|
| 1 | `POST /api/bot/get_type_list` | `GET /api/v1/bot-types` | GetTypeList | P0 |
| 2 | `POST /api/conversation/get_message_list` | `GET /api/v1/conversations/:id/messages` | GetMessageList | P0 |
| 3 | `POST /api/knowledge/detail` | `GET /api/v1/knowledge/:id` | DatasetDetail | P0 |
| 4 | `POST /api/knowledge/list` | `GET /api/v1/knowledge` | ListDataset | P0 |
| 5 | `POST /api/database/get_by_id` | `GET /api/v1/databases/:id` | GetDatabaseByID | P0 |
| 6 | `POST /api/workflow_api/list_spans` | `GET /api/v1/workflows/:id/spans` | ListRootSpans | P1 |
| 7 | `POST /api/workflow_api/workflow_list` | `GET /api/v1/workflows` | GetWorkFlowList | P1 |
| 8 | `POST /api/plugin_api/get_plugin_info` | `GET /api/v1/plugins/:id` | GetPluginInfo | P1 |
| 9 | `POST /api/database/list` | `GET /api/v1/databases` | ListDatabase | P1 |
| 10 | `POST /api/draftbot/get_display_info` | `GET /api/v1/draft-bots/:id` | GetDraftBotDisplayInfo | P2 |

---

## 1️⃣ API #1: Bot类型列表

### 1.1 当前实现 (❌ 违规)

**路由定义**:
```go
// backend/api/router/coze/api.go:46
_bot.POST("/get_type_list", append(_gettypelistMw(), coze.GetTypeList)...)
```

**Handler实现**:
```go
// backend/api/handler/coze/developer_api_service.go:399
// @router /api/bot/get_type_list [POST]
func GetTypeList(ctx context.Context, c *app.RequestContext) {
	var req developer_api.GetTypeListRequest
	err = c.BindAndValidate(&req)  // ❌ 从Body读取
	// ...
}
```

**请求示例**:
```bash
# ❌ 错误: 查询操作使用POST
POST /api/bot/get_type_list
Content-Type: application/json

{}
```

### 1.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
// backend/api/router/coze/api.go
// 在 _bot Group 内新增
_bot.GET("/types", append(_gettypelistMw(), coze.GetTypeListV2)...)
```

**新增Handler实现**:
```go
// backend/api/handler/coze/developer_api_service.go
// GetTypeListV2 获取Bot类型列表(V2 - 符合RESTful规范)
// @router /api/bot/types [GET]
// @Summary 获取Bot类型列表
// @Description 获取所有可用的Bot类型
// @Tags Bot
// @Accept json
// @Produce json
// @Success 200 {object} httputil.BaseResponse{data=developer_api.GetTypeListResponse}
// @Router /api/bot/types [get]
func GetTypeListV2(ctx context.Context, c *app.RequestContext) {
	// 从Query参数读取
	req := developer_api.GetTypeListRequest{
		// 如果有查询参数,从c.Query()读取
		IncludeDeprecated := c.Query("include_deprecated", "false") == "true",
	}

	resp, err := application.DeveloperAPISVC.GetTypeList(ctx, &req)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}
	httputil.BuildSuccessResp(c, resp)
}

// GetTypeList 获取Bot类型列表(旧版本 - 已废弃)
// @router /api/bot/get_type_list [POST]
// @Deprecated Use GET /api/bot/types instead
func GetTypeList(ctx context.Context, c *app.RequestContext) {
	// 重定向到新版本
	GetTypeListV2(ctx, c)
}
```

**请求示例**:
```bash
# ✅ 正确: 查询操作使用GET
GET /api/bot/types?include_deprecated=false
```

**前端迁移**:
```typescript
// ❌ 旧版本 (已废弃)
const response = await fetch('/api/bot/get_type_list', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({})
});

// ✅ 新版本 (推荐)
const response = await fetch('/api/bot/types?include_deprecated=false', {
  method: 'GET'
});
```

---

## 2️⃣ API #2: 会话消息列表

### 2.1 当前实现 (❌ 违规)

**路由定义**:
```go
// backend/api/router/coze/api.go:65
_conversation.POST("/get_message_list", append(_getmessagelistMw(), coze.GetMessageList)...)
```

**Handler实现**:
```go
// backend/api/handler/coze/message_service.go:35
// @router /api/conversation/get_message_list [POST]
func GetMessageList(ctx context.Context, c *app.RequestContext) {
	var req message.GetMessageListRequest
	err = c.BindAndValidate(&req)  // ❌ 从Body读取
	botID := req.BotID
	conversationID := req.ConversationID
	page := req.Page
	pageSize := req.PageSize
	// ...
}
```

**请求示例**:
```bash
# ❌ 错误: 查询操作使用POST
POST /api/conversation/get_message_list
Content-Type: application/json

{
  "bot_id": "123",
  "conversation_id": "456",
  "page": 1,
  "page_size": 20
}
```

### 2.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
// backend/api/router/coze/api.go
// 在 _conversation Group 内新增
_conversation.GET("/:conversation_id/messages",
	append(_getmessagelistMw(), coze.GetMessageListV2)...)
```

**新增Handler实现**:
```go
// backend/api/handler/coze/message_service.go
// GetMessageListV2 获取会话消息列表(V2 - 符合RESTful规范)
// @router /api/conversation/:conversation_id/messages [GET]
// @Summary 获取会话消息列表
// @Description 分页获取指定会话的消息列表
// @Tags Conversation
// @Param conversation_id path string true "会话ID"
// @Param bot_id query string false "Bot ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} httputil.BaseResponse{data=message.GetMessageListResponse}
// @Router /api/conversation/{conversation_id}/messages [get]
func GetMessageListV2(ctx context.Context, c *app.RequestContext) {
	// 从路径参数读取
	conversationID := c.Param("conversation_id")

	// 从Query参数读取
	req := message.GetMessageListRequest{
		ConversationID: conversationID,
		BotID:          c.Query("bot_id", ""),
		Page:           c.QueryInt("page", 1),
		PageSize:       c.QueryInt("page_size", 20),
	}

	if checkErr := checkMLParams(ctx, &req); checkErr != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode,
			checkErr.Error(), "参数验证失败", nil)
		return
	}

	resp, err := application.ConversationSVC.GetMessageList(ctx, &req)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}

	httputil.BuildSuccessResp(c, resp)
}

// GetMessageList 获取会话消息列表(旧版本 - 已废弃)
// @router /api/conversation/get_message_list [POST]
// @Deprecated Use GET /api/conversation/:conversation_id/messages instead
func GetMessageList(ctx context.Context, c *app.RequestContext) {
	// 重定向到新版本
	var req message.GetMessageListRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		httputil.BuildErrorResp(c, errno.ErrInvalidParamCode,
			err.Error(), "参数验证失败", nil)
		return
	}

	// 模拟路径参数
	c.SetParam("conversation_id", req.ConversationID)
	c.SetParam("bot_id", req.BotID)
	c.SetParam("page", fmt.Sprintf("%d", req.Page))
	c.SetParam("page_size", fmt.Sprintf("%d", req.PageSize))

	GetMessageListV2(ctx, c)
}
```

**请求示例**:
```bash
# ✅ 正确: 查询操作使用GET,参数在路径和Query中
GET /api/conversation/456/messages?bot_id=123&page=1&page_size=20
```

**前端迁移**:
```typescript
// ❌ 旧版本 (已废弃)
const response = await fetch('/api/conversation/get_message_list', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    bot_id: '123',
    conversation_id: '456',
    page: 1,
    page_size: 20
  })
});

// ✅ 新版本 (推荐)
const response = await fetch(
  '/api/conversation/456/messages?bot_id=123&page=1&page_size=20',
  { method: 'GET' }
);
```

---

## 3️⃣ API #3: 知识库详情

### 3.1 当前实现 (❌ 违规)

**路由定义**:
```go
// backend/api/router/coze/api.go:119
_knowledge0.POST("/detail", append(_datasetdetailMw(), coze.DatasetDetail)...)
```

**Handler实现**:
```go
// backend/api/handler/coze/knowledge_service.go:56
// @router /api/knowledge/detail [POST]
func DatasetDetail(ctx context.Context, c *app.RequestContext) {
	var req dataset.DatasetDetailRequest
	err = c.BindAndValidate(&req)  // ❌ 从Body读取
	datasetID := req.DatasetID
	// ...
}
```

**请求示例**:
```bash
# ❌ 错误: 查询操作使用POST
POST /api/knowledge/detail
Content-Type: application/json

{
  "dataset_id": "123"
}
```

### 3.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
// backend/api/router/coze/api.go
// 在 _knowledge Group 内新增
_knowledge.GET("/:dataset_id", append(_datasetdetailMw(), coze.DatasetDetailV2)...)
```

**新增Handler实现**:
```go
// backend/api/handler/coze/knowledge_service.go
// DatasetDetailV2 获取知识库详情(V2 - 符合RESTful规范)
// @router /api/knowledge/:dataset_id [GET]
// @Summary 获取知识库详情
// @Description 根据ID获取知识库详细信息
// @Tags Knowledge
// @Param dataset_id path string true "知识库ID"
// @Success 200 {object} httputil.BaseResponse{data=dataset.DatasetDetailResponse}
// @Router /api/knowledge/{dataset_id} [get]
func DatasetDetailV2(ctx context.Context, c *app.RequestContext) {
	// 从路径参数读取
	datasetID := c.Param("dataset_id")

	req := dataset.DatasetDetailRequest{
		DatasetID: datasetID,
	}

	resp := new(dataset.DatasetDetailResponse)
	resp, err := application.KnowledgeSVC.DatasetDetail(ctx, &req)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}
	httputil.BuildSuccessResp(c, resp)
}

// DatasetDetail 获取知识库详情(旧版本 - 已废弃)
// @router /api/knowledge/detail [POST]
// @Deprecated Use GET /api/knowledge/:dataset_id instead
func DatasetDetail(ctx context.Context, c *app.RequestContext) {
	// 重定向到新版本
	var req dataset.DatasetDetailRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.SetParam("dataset_id", req.DatasetID)
	DatasetDetailV2(ctx, c)
}
```

**请求示例**:
```bash
# ✅ 正确: 查询操作使用GET,参数在路径中
GET /api/knowledge/123
```

**前端迁移**:
```typescript
// ❌ 旧版本 (已废弃)
const response = await fetch('/api/knowledge/detail', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ dataset_id: '123' })
});

// ✅ 新版本 (推荐)
const response = await fetch('/api/knowledge/123', { method: 'GET' });
```

---

## 4️⃣ API #4: 知识库列表

### 4.1 当前实现 (❌ 违规)

**路由定义**:
```go
// backend/api/router/coze/api.go:120
_knowledge0.POST("/list", append(_listdatasetMw(), coze.ListDataset)...)
```

**Handler实现**:
```go
// backend/api/handler/coze/knowledge_service.go:75
// @router /api/knowledge/list [POST]
func ListDataset(ctx context.Context, c *app.RequestContext) {
	var req dataset.ListDatasetRequest
	err = c.BindAndValidate(&req)  // ❌ 从Body读取
	page := req.Page
	pageSize := req.PageSize
	// ...
}
```

**请求示例**:
```bash
# ❌ 错误: 查询操作使用POST
POST /api/knowledge/list
Content-Type: application/json

{
  "page": 1,
  "page_size": 20
}
```

### 4.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
// backend/api/router/coze/api.go
// 在 _knowledge Group 内新增
_knowledge.GET("", append(_listdatasetMw(), coze.ListDatasetV2)...)
```

**新增Handler实现**:
```go
// backend/api/handler/coze/knowledge_service.go
// ListDatasetV2 获取知识库列表(V2 - 符合RESTful规范)
// @router /api/knowledge [GET]
// @Summary 获取知识库列表
// @Description 分页获取知识库列表
// @Tags Knowledge
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param search query string false "搜索关键词"
// @Success 200 {object} httputil.BaseResponse{data=dataset.ListDatasetResponse}
// @Router /api/knowledge [get]
func ListDatasetV2(ctx context.Context, c *app.RequestContext) {
	// 从Query参数读取
	req := dataset.ListDatasetRequest{
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("page_size", 20),
		Search:    c.Query("search", ""),
		// 其他查询参数...
	}

	resp := new(dataset.ListDatasetResponse)
	resp, err := application.KnowledgeSVC.ListKnowledge(ctx, &req)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}
	httputil.BuildSuccessResp(c, resp)
}

// ListDataset 获取知识库列表(旧版本 - 已废弃)
// @router /api/knowledge/list [POST]
// @Deprecated Use GET /api/knowledge instead
func ListDataset(ctx context.Context, c *app.RequestContext) {
	// 重定向到新版本
	var req dataset.ListDatasetRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 设置Query参数
	c.SetParam("page", fmt.Sprintf("%d", req.Page))
	c.SetParam("page_size", fmt.Sprintf("%d", req.PageSize))
	c.SetParam("search", req.Search)

	ListDatasetV2(ctx, c)
}
```

**请求示例**:
```bash
# ✅ 正确: 查询操作使用GET,参数在Query中
GET /api/knowledge?page=1&page_size=20&search=客服
```

**前端迁移**:
```typescript
// ❌ 旧版本 (已废弃)
const response = await fetch('/api/knowledge/list', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ page: 1, page_size: 20 })
});

// ✅ 新版本 (推荐)
const response = await fetch('/api/knowledge?page=1&page_size=20', {
  method: 'GET'
});
```

---

## 5️⃣ API #5: 数据库详情

### 5.1 当前实现 (❌ 违规)

**路由定义**:
```go
// backend/api/router/coze/api.go:199
_database.POST("/get_by_id", append(_getdatabasebyidMw(), coze.GetDatabaseByID)...)
```

**Handler实现**:
```go
// backend/api/handler/coze/database_service.go:56
// @router /api/database/get_by_id [POST]
func GetDatabaseByID(ctx context.Context, c *app.RequestContext) {
	var req table.GetDatabaseByIDRequest
	err = c.BindAndValidate(&req)  // ❌ 从Body读取
	databaseID := req.DatabaseID
	// ...
}
```

**请求示例**:
```bash
# ❌ 错误: 查询操作使用POST
POST /api/database/get_by_id
Content-Type: application/json

{
  "database_id": "123"
}
```

### 5.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
// backend/api/router/coze/api.go
// 在 _database Group 内新增
_database.GET("/:database_id", append(_getdatabasebyidMw(), coze.GetDatabaseByIDV2)...)
```

**新增Handler实现**:
```go
// backend/api/handler/coze/database_service.go
// GetDatabaseByIDV2 获取数据库详情(V2 - 符合RESTful规范)
// @router /api/database/:database_id [GET]
// @Summary 获取数据库详情
// @Description 根据ID获取数据库详细信息
// @Tags Database
// @Param database_id path string true "数据库ID"
// @Success 200 {object} httputil.BaseResponse{data=table.GetDatabaseByIDResponse}
// @Router /api/database/{database_id} [get]
func GetDatabaseByIDV2(ctx context.Context, c *app.RequestContext) {
	// 从路径参数读取
	databaseID := c.Param("database_id")

	req := table.GetDatabaseByIDRequest{
		DatabaseID: databaseID,
	}

	resp, err := memory.DatabaseApplicationSVC.GetDatabaseByID(ctx, &req)
	if err != nil {
		httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))
		return
	}
	httputil.BuildSuccessResp(c, resp)
}

// GetDatabaseByID 获取数据库详情(旧版本 - 已废弃)
// @router /api/database/get_by_id [POST]
// @Deprecated Use GET /api/database/:database_id instead
func GetDatabaseByID(ctx context.Context, c *app.RequestContext) {
	// 重定向到新版本
	var req table.GetDatabaseByIDRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	c.SetParam("database_id", req.DatabaseID)
	GetDatabaseByIDV2(ctx, c)
}
```

**请求示例**:
```bash
# ✅ 正确: 查询操作使用GET,参数在路径中
GET /api/database/123
```

**前端迁移**:
```typescript
// ❌ 旧版本 (已废弃)
const response = await fetch('/api/database/get_by_id', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ database_id: '123' })
});

// ✅ 新版本 (推荐)
const response = await fetch('/api/database/123', { method: 'GET' });
```

---

## 6️⃣ API #6: 工作流Span列表

### 6.1 当前实现 (❌ 违规)

**路由定义**:
```go
// backend/api/router/coze/api.go:XXX
_workflow_api.POST("/list_spans", append(_listrootspansMw(), coze.ListRootSpans)...)
```

**Handler实现**:
```go
// backend/api/handler/coze/workflow_service.go:XXX
// @router /api/workflow_api/list_spans [POST]
func ListRootSpans(ctx context.Context, c *app.RequestContext) {
	var req workflow.ListRootSpansRequest
	err = c.BindAndValidate(&req)  // ❌ 从Body读取
	workflowID := req.WorkflowID
	// ...
}
```

### 6.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
// backend/api/router/coze/api.go
_workflow_api.GET("/:workflow_id/spans",
	append(_listrootspansMw(), coze.ListRootSpansV2)...)
```

**新增Handler实现**:
```go
// ListRootSpansV2 获取工作流Span列表(V2 - 符合RESTful规范)
// @router /api/workflow_api/:workflow_id/spans [GET]
func ListRootSpansV2(ctx context.Context, c *app.RequestContext) {
	workflowID := c.Param("workflow_id")

	req := workflow.ListRootSpansRequest{
		WorkflowID: workflowID,
		Page:       c.QueryInt("page", 1),
		PageSize:   c.QueryInt("page_size", 20),
	}

	// ... 业务逻辑
}
```

---

## 7️⃣ API #7: 工作流列表

### 7.1 当前实现 (❌ 违规)

**路由定义**:
```go
_workflow_api.POST("/workflow_list", append(_getworkflowlistMw(), coze.GetWorkFlowList)...)
```

### 7.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
_workflow_api.GET("", append(_getworkflowlistMw(), coze.GetWorkFlowListV2)...)
```

**新增Handler实现**:
```go
// GetWorkFlowListV2 获取工作流列表(V2 - 符合RESTful规范)
// @router /api/workflow_api [GET]
func GetWorkFlowListV2(ctx context.Context, c *app.RequestContext) {
	req := workflow.GetWorkFlowListRequest{
		Page:      c.QueryInt("page", 1),
		PageSize:  c.QueryInt("page_size", 20),
		Search:    c.Query("search", ""),
	}

	// ... 业务逻辑
}
```

---

## 8️⃣ API #8: 插件信息

### 8.1 当前实现 (❌ 违规)

**路由定义**:
```go
_plugin_api.POST("/get_plugin_info", append(_getplugininfoMw(), coze.GetPluginInfo)...)
```

### 8.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
_plugin_api.GET("/:plugin_id", append(_getplugininfoMw(), coze.GetPluginInfoV2)...)
```

**新增Handler实现**:
```go
// GetPluginInfoV2 获取插件信息(V2 - 符合RESTful规范)
// @router /api/plugin_api/:plugin_id [GET]
func GetPluginInfoV2(ctx context.Context, c *app.RequestContext) {
	pluginID := c.Param("plugin_id")

	req := plugin.GetPluginInfoRequest{
		PluginID: pluginID,
	}

	// ... 业务逻辑
}
```

---

## 9️⃣ API #9: 数据库列表

### 9.1 当前实现 (❌ 违规)

**路由定义**:
```go
_database.POST("/list", append(_listdatabaseMw(), coze.ListDatabase)...)
```

### 9.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
_database.GET("", append(_listdatabaseMw(), coze.ListDatabaseV2)...)
```

**新增Handler实现**:
```go
// ListDatabaseV2 获取数据库列表(V2 - 符合RESTful规范)
// @router /api/database [GET]
func ListDatabaseV2(ctx context.Context, c *app.RequestContext) {
	req := table.ListDatabaseRequest{
		Page:     c.QueryInt("page", 1),
		PageSize: c.QueryInt("page_size", 20),
	}

	// ... 业务逻辑
}
```

---

## 🔟 API #10: Draft Bot显示信息

### 10.1 当前实现 (❌ 违规)

**路由定义**:
```go
_draftbot.POST("/get_display_info", append(_getdraftbotdisplayinfoMw(), coze.GetDraftBotDisplayInfo)...)
```

### 10.2 修复方案 (✅ 规范)

**新增路由定义**:
```go
_draftbot.GET("/:bot_id", append(_getdraftbotdisplayinfoMw(), coze.GetDraftBotDisplayInfoV2)...)
```

**新增Handler实现**:
```go
// GetDraftBotDisplayInfoV2 获取Draft Bot显示信息(V2 - 符合RESTful规范)
// @router /api/draftbot/:bot_id [GET]
func GetDraftBotDisplayInfoV2(ctx context.Context, c *app.RequestContext) {
	botID := c.Param("bot_id")

	req := developer_api.GetDraftBotDisplayInfoRequest{
		BotID: botID,
	}

	// ... 业务逻辑
}
```

---

## 📊 迁移策略

### 阶段1: 双版本并存 (Month 1-3)

✅ 保留旧API,新增符合RESTful规范的版本
- 旧API标记为`@Deprecated`
- 新API在Swagger文档中标记为"推荐"
- 前端团队逐步迁移调用

### 阶段2: 强制迁移 (Month 4-5)

⚠️ 旧API返回弃用警告
```json
{
  "code": 0,
  "message": "SUCCESS",
  "message_zh": "操作成功",
  "message_en": "Operation successful",
  "data": {...},
  "deprecation_warning": {
    "message_zh": "此API已废弃,请使用 GET /api/v1/xxx",
    "message_en": "This API is deprecated, use GET /api/v1/xxx instead",
    "new_api": "GET /api/v1/xxx",
    "sunset_at": "2025-06-01"
  }
}
```

### 阶段3: 移除旧API (Month 6)

❌ 完全移除旧API
- 移除旧路由定义
- 移除旧Handler函数
- 更新文档

---

## 🧪 测试计划

### 单元测试

为每个新API编写单元测试:

```go
func TestGetTypeListV2(t *testing.T) {
	// 测试用例1: 正常查询
	// 测试用例2: 包含已废弃类型
	// 测试用例3: 无权限
}
```

### 集成测试

```bash
# 测试新API
curl -X GET "http://localhost:8888/api/bot/types?include_deprecated=false"

# 测试旧API(应返回警告)
curl -X POST "http://localhost:8888/api/bot/get_type_list"
```

### 性能测试

```bash
# 对比新旧API性能
ab -n 10000 -c 100 "http://localhost:8888/api/bot/types"
ab -n 10000 -c 100 "http://localhost:8888/api/bot/get_type_list"
```

---

## 📈 预期成果

### API规范评分提升

| 维度 | 修复前 | 修复后 | 提升 |
|------|--------|--------|------|
| HTTP方法规范 | 95/100 | 100/100 | +5 |
| 路径命名规范 | 85/100 | 95/100 | +10 |
| **总分** | **92.5** | **97.5** | **+5** |

### RESTful最佳实践

✅ 查询操作使用GET
✅ 支持浏览器缓存
✅ 支持CDN缓存
✅ 幂等性保证
✅ 符合HTTP语义

---

## 📝 实施清单

### Week 1: 代码实现

- [ ] API #1: Bot类型列表
- [ ] API #2: 会话消息列表
- [ ] API #3: 知识库详情
- [ ] API #4: 知识库列表
- [ ] API #5: 数据库详情

### Week 2: 代码实现 + 测试

- [ ] API #6: 工作流Span列表
- [ ] API #7: 工作流列表
- [ ] API #8: 插件信息
- [ ] API #9: 数据库列表
- [ ] API #10: Draft Bot显示信息
- [ ] 单元测试覆盖率 ≥ 80%

### Week 3: 前端迁移

- [ ] 识别所有调用旧API的前端代码
- [ ] 逐步迁移到新API
- [ ] 回归测试

### Week 4: 文档和发布

- [ ] 更新Swagger文档
- [ ] 编写迁移指南
- [ ] Code Review
- [ ] 发布到生产环境

---

## 🎓 参考资料

- [ZKER-API响应格式统一报告_v1.0.md](./ZKER-API响应格式统一报告_v1.0.md)
- [ZKER-企业级开发规范手册_v1.0.md](./ZKER-企业级开发规范手册_v1.0.md)
- [RESTful API设计规范](https://restfulapi.net/)
- [HTTP方法规范](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods)

---

**方案维护者**: API规范化专家 (Claude Code)
**方案版本**: v1.0
**最后更新**: 2025-01-01
**预计完成**: 2025-01-31

---

**⚠️ 重要提醒**:
1. 所有新API必须使用GET方法
2. 参数从Body改为Query/Path参数
3. 保持向后兼容,保留旧API
4. 6个月后完全移除旧API
5. 前端团队需要配合迁移
