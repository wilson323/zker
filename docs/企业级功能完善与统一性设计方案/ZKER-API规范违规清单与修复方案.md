# ZKER API规范违规清单与修复方案

**📅 生成日期**: 2025-01-01
**🎯 用途**: 详细的API违规清单和可执行的修复方案
**📋 关联文档**: [ZKER-API规范深度审查报告_v1.0.md](./ZKER-API规范深度审查报告_v1.0.md)

---

## 📊 违规统计总览

| 违规类型 | 数量 | 严重程度 | 优先级 |
|---------|------|---------|--------|
| 响应格式不统一 | 601个 | 🔴 严重 | P0 |
| POST用于查询 | 30+个 | 🔴 严重 | P0 |
| POST用于删除 | 5个 | 🔴 严重 | P0 |
| POST用于更新 | 10个 | 🟡 中等 | P0 |
| 单数命名 | 23处 | 🟡 中等 | P1 |
| 驼峰命名 | 15+处 | 🟡 中等 | P1 |
| 缺少Swagger文档 | 422个 | 🔴 严重 | P1 |
| 缺少版本号 | 1个 | 🟢 轻微 | P2 |

---

## 1️⃣ 响应格式不统一 (601个) - P0

### 问题描述

704个JSON响应中,601个(85.4%)未使用企业级统一响应格式。

### 违规案例

#### 案例1: 直接返回实体 (最常见)

**文件**: `backend/api/handler/coze/agent_run_service.go`

```go
// ❌ 当前实现
func AgentRun(ctx context.Context, c *app.RequestContext) {
    resp := &AgentRunResponse{...}
    c.JSON(consts.StatusOK, resp)  // 直接返回,缺少code/message字段
}

// ✅ 修复方案
import "github.com/coze-dev/coze-studio/backend/api/internal/httputil"

func AgentRun(ctx context.Context, c *app.RequestContext) {
    resp := &AgentRunResponse{...}
    httputil.BuildSuccessResp(c, resp)  // 使用统一响应格式
}
```

#### 案例2: 错误响应格式不一致

**文件**: `backend/api/handler/coze/billing_handler.go`

```go
// ❌ 当前实现 - 自定义APIResponse,与标准不一致
type APIResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

c.JSON(http.StatusBadRequest, APIResponse{
    Code:    400,
    Message: "Invalid request",
})

// ✅ 修复方案 - 使用httputil.ErrorResponse
import "github.com/coze-dev/coze-studio/backend/api/internal/httputil"
import "github.com/coze-dev/coze-studio/backend/types/errno"

httputil.BuildErrorResp(c, errno.ErrInvalidParam, "Invalid request", "参数验证失败", nil)
```

### 批量修复脚本

```bash
#!/bin/bash
# scripts/fix_response_format.sh

echo "🔧 开始修复响应格式..."

# 步骤1: 查找所有违规文件
echo "📋 步骤1: 查找违规文件"
grep -rn "c\.JSON(consts\.StatusOK, resp)" backend/api/handler/ --include="*.go" -l > /tmp/violation_files.txt

# 步骤2: 添加httputil导入
echo "📦 步骤2: 添加httputil导入"
while read -r file; do
    if ! grep -q "backend/api/internal/httputil" "$file"; then
        # 在import区域添加导入
        sed -i '/^import (/a\\t"github.com/coze-dev/coze-studio/backend/api/internal/httputil"' "$file"
    fi
done < /tmp/violation_files.txt

# 步骤3: 替换响应调用
echo "🔄 步骤3: 替换响应调用"
find backend/api/handler -name "*.go" -type f -exec sed -i '
    s/c\.JSON(consts\.StatusOK, resp)/httputil.BuildSuccessResp(c, resp)/g
    s/c\.JSON(consts\.StatusOK, \&resp)/httputil.BuildSuccessResp(c, \&resp)/g
' {} \;

echo "✅ 修复完成! 请检查并提交更改。"
```

### 验证脚本

```bash
#!/bin/bash
# scripts/verify_response_format.sh

echo "🔍 验证响应格式..."

# 检查是否还有直接返回实体的情况
violations=$(grep -rn "c\.JSON(consts\.StatusOK, resp)" backend/api/handler/ --include="*.go" | wc -l)

if [ "$violations" -eq 0 ]; then
    echo "✅ 响应格式检查通过!"
else
    echo "❌ 发现 $violations 个违规响应,请继续修复"
    exit 1
fi
```

### 完整违规文件清单

```
backend/api/handler/coze/agent_run_service.go
backend/api/handler/coze/bot_open_api_service.go
backend/api/handler/coze/conversation_service.go
backend/api/handler/coze/intelligence_service.go
backend/api/handler/coze/knowledge_service.go
backend/api/handler/coze/memory_service.go
backend/api/handler/coze/message_service.go
backend/api/handler/coze/playground_service.go
backend/api/handler/coze/workflow_service.go
backend/api/handler/coze/config_service.go
backend/api/handler/coze/database_service.go
backend/api/handler/coze/developer_api_service.go
backend/api/handler/coze/knowledgegraph/graph_query_handler.go
backend/api/handler/coze/plugin_develop_service.go
backend/api/handler/coze/publish_service.go
backend/api/handler/coze/resource_service.go
backend/api/handler/coze/routing_service.go
backend/api/handler/coze/upload_service.go
... (共63个文件,601处违规)
```

---

## 2️⃣ HTTP方法误用 (45+个) - P0

### 问题描述

30+个查询接口使用POST,违反RESTful幂等性原则。

### 违规清单

#### 2.1 POST用于查询 (30+个)

| 序号 | 当前路径 | 当前方法 | 正确路径 | 正确方法 | 文件位置 |
|------|---------|---------|---------|---------|---------|
| 1 | `/bot/get_type_list` | POST | `/api/v1/bot-types` | GET | api.go:46 |
| 2 | `/conversation/get_message_list` | POST | `/api/v1/conversations/:id/messages` | GET | api.go:65 |
| 3 | `/developer/get_icon` | POST | `/api/v1/developers/:id/icon` | GET | api.go:69 |
| 4 | `/draftbot/get_display_info` | POST | `/api/v1/draft-bots/:id` | GET | api.go:77 |
| 5 | `/knowledge/detail` | POST | `/api/v1/knowledge/:id` | GET | api.go:119 |
| 6 | `/knowledge/document/list` | POST | `/api/v1/knowledge/:id/documents` | GET | api.go:126 |
| 7 | `/knowledge/document/progress/get` | POST | `/api/v1/documents/:id/progress` | GET | api.go:131 |
| 8 | `/knowledge/icon/get` | POST | `/api/v1/knowledge/:id/icon` | GET | api.go:136 |
| 9 | `/knowledge/slice/list` | POST | `/api/v1/knowledge/:id/slices` | GET | api.go:155 |
| 10 | `/knowledge/table_schema/get` | POST | `/api/v1/table-schemas/:id` | GET | api.go:160 |
| 11 | `/memory/database/get_by_id` | POST | `/api/v1/databases/:id` | GET | api.go:195 |
| 12 | `/memory/database/get_connector_name` | POST | `/api/v1/databases/:id/connector` | GET | api.go:196 |
| 13 | `/memory/database/get_online_database_id` | POST | `/api/v1/databases/:id/online-id` | GET | api.go:197 |
| 14 | `/memory/database/get_template` | POST | `/api/v1/database-templates/:id` | GET | api.go:198 |
| 15 | `/memory/database/list` | POST | `/api/v1/databases` | GET | api.go:199 |
| 16 | `/publish/publish_record_detail` | POST | `/api/v1/publish-records/:id` | GET | api.go:105 |
| 17 | `/publish/publish_record_list` | POST | `/api/v1/publish-records` | GET | api.go:106 |
| 18 | `/search/get_draft_intelligence_info` | POST | `/api/v1/draft-intelligences/:id` | GET | api.go:110 |
| 19 | `/search/get_draft_intelligence_list` | POST | `/api/v1/draft-intelligences` | GET | api.go:111 |
| 20 | `/playground/get_onboarding` | POST | `/api/v1/playground/onboarding` | GET | api.go:311 |

#### 2.2 POST用于删除 (5个)

| 序号 | 当前路径 | 当前方法 | 正确路径 | 正确方法 | 文件位置 |
|------|---------|---------|---------|---------|---------|
| 1 | `/knowledge/delete` | POST | `/api/v1/knowledge/:id` | DELETE | api.go:118 |
| 2 | `/knowledge/document/delete` | POST | `/api/v1/documents/:id` | DELETE | api.go:125 |
| 3 | `/draftbot/delete` | POST | `/api/v1/draft-bots/:id` | DELETE | api.go:75 |
| 4 | `/draft_project/delete` | POST | `/api/v1/draft-projects/:id` | DELETE | api.go:95 |
| 5 | `/knowledge/slice/delete` | POST | `/api/v1/slices/:id` | DELETE | api.go:154 |

#### 2.3 POST用于更新 (10+个)

| 序号 | 当前路径 | 当前方法 | 正确路径 | 正确方法 | 文件位置 |
|------|---------|---------|---------|---------|---------|
| 1 | `/knowledge/update` | POST | `/api/v1/knowledge/:id` | PUT | api.go:121 |
| 2 | `/knowledge/document/update` | POST | `/api/v1/documents/:id` | PUT | api.go:128 |
| 3 | `/draftbot/update_display_info` | POST | `/api/v1/draft-bots/:id` | PATCH | api.go:80 |
| 4 | `/draft_project/update` | POST | `/api/v1/draft-projects/:id` | PUT | api.go:97 |
| 5 | `/knowledge/slice/update` | POST | `/api/v1/slices/:id` | PUT | api.go:156 |

### 修复示例

#### 示例1: 查询接口

**当前实现**:
```go
// backend/api/router/coze/api.go
_knowledge.POST("/list", append(_listdatasetMw(), coze.ListDataset)...)

// backend/api/handler/coze/knowledge_service.go
func ListDataset(ctx context.Context, c *app.RequestContext) {
    var req ListDatasetRequest
    c.BindAndValidate(&req)  // 从body读取参数
}
```

**修复方案**:
```go
// backend/api/router/coze/api.go
_knowledge.GET("", append(_listdatasetMw(), coze.ListDataset)...)  // 改为GET

// backend/api/handler/coze/knowledge_service.go
func ListDataset(ctx context.Context, c *app.RequestContext) {
    // 从query参数读取
    page := c.Query("page", "1")
    pageSize := c.Query("page_size", "20")
    search := c.Query("search", "")
}
```

#### 示例2: 详情查询

**当前实现**:
```go
_knowledge.POST("/detail", append(_datasetdetailMw(), coze.DatasetDetail)...)

func DatasetDetail(ctx context.Context, c *app.RequestContext) {
    var req DatasetDetailRequest
    c.BindAndValidate(&req)
    knowledgeID := req.KnowledgeID
}
```

**修复方案**:
```go
_knowledge.GET("/:knowledge_id", append(_datasetdetailMw(), coze.DatasetDetail)...)

func DatasetDetail(ctx context.Context, c *app.RequestContext) {
    knowledgeID := c.Param("knowledge_id")  // 从路径参数获取
}
```

#### 示例3: 删除接口

**当前实现**:
```go
_knowledge.POST("/delete", append(_deletedatasetMw(), coze.DeleteDataset)...)

func DeleteDataset(ctx context.Context, c *app.RequestContext) {
    var req DeleteDatasetRequest
    c.BindAndValidate(&req)
    knowledgeID := req.KnowledgeID
}
```

**修复方案**:
```go
_knowledge.DELETE("/:knowledge_id", append(_deletedatasetMw(), coze.DeleteDataset)...)

func DeleteDataset(ctx context.Context, c *app.RequestContext) {
    knowledgeID := c.Param("knowledge_id")
}
```

### 批量修复工具

```bash
#!/bin/bash
# scripts/fix_http_methods.sh

echo "🔧 开始修复HTTP方法..."

# 步骤1: 查找所有POST查询
echo "📋 查找POST查询接口..."
grep -rn "\.POST.*\/list\|\.POST.*\/get\|\.POST.*\/detail" backend/api/router/ --include="*.go" > /tmp/post_get_violations.txt

# 步骤2: 查找所有POST删除
echo "📋 查找POST删除接口..."
grep -rn "\.POST.*\/delete" backend/api/router/ --include="*.go" > /tmp/post_delete_violations.txt

# 步骤3: 查找所有POST更新
echo "📋 查找POST更新接口..."
grep -rn "\.POST.*\/update" backend/api/router/ --include="*.go" > /tmp/post_update_violations.txt

echo "⚠️  请手动修改上述文件,建议修改顺序:"
echo "1. /tmp/post_get_violations.txt (POST查询改为GET)"
echo "2. /tmp/post_delete_violations.txt (POST删除改为DELETE)"
echo "3. /tmp/post_update_violations.txt (POST更新改为PUT/PATCH)"

echo "✅ 违规文件已导出,请逐个修改"
```

---

## 3️⃣ 路径命名不规范 (38+处) - P1

### 问题描述

使用单数命名、驼峰命名,不符合RESTful规范。

### 违规清单

#### 3.1 单数命名 (23处)

| 当前路径 | 正确路径 | 类型 |
|---------|---------|------|
| `/draftbot` | `/draft-bots` | 路由组 |
| `/intelligence_api` | `/intelligence-apis` | 路由组 |
| `/bot` | `/bots` | 路由组 |
| `/conversation` | `/conversations` | 路由组 |
| `/knowledge` | `/knowledge` | ✅ 正确 |
| `/memory` | `/memories` | 建议修正 |
| `/marketplace` | `/marketplaces` | 建议修正 |

#### 3.2 驼峰命名 (15+处)

| 当前路径 | 正确路径 | 类型 |
|---------|---------|------|
| `/draft_project` | `/draft-projects` | 路由组 |
| `/inner_task_list` | `/inner-tasks` | 路由 |
| `/get_type_list` | `/types` | 路由 |
| `/get_display_info` | (作为GET方法) | 方法名 |
| `/get_by_id` | (使用路径参数) | 方法名 |
| `/get_connector_name` | `/connector` | 路由 |
| `/get_template` | `/template` | 路由 |
| `/list_new` | `/list` | 路由 |

### 修复方案

#### 方案1: 路由组重命名

**当前实现**:
```go
// backend/api/router/coze/api.go
_draftbot := _api.Group("/draftbot")
_draft_project := _intelligence_api.Group("/draft_project")
```

**修复方案**:
```go
// backend/api/router/coze/api.go
_draft_bots := _api.Group("/draft-bots")
_draft_projects := _intelligence_api.Group("/draft-projects")
```

#### 方案2: 路由路径优化

**当前实现**:
```go
_bot.POST("/get_type_list", GetTypeList)
_draftbot.POST("/get_display_info", GetDisplayInfo)
_database.POST("/get_by_id", GetDatabaseByID)
_database.POST("/get_connector_name", GetConnectorName)
```

**修复方案**:
```go
_bot_types := _api.Group("/bot-types")
_bot_types.GET("", GetTypeList)

_draft_bots.GET("/:bot_id", GetDraftBotInfo)

_databases := _api.Group("/databases")
_databases.GET("/:database_id", GetDatabaseByID)
_databases.GET("/:database_id/connector", GetConnectorName)
```

### 路由设计最佳实践

#### 资源嵌套规则

```go
// ✅ 正确: 资源嵌套清晰
GET  /api/v1/bots/:bot_id/config
POST /api/v1/bots/:bot_id/publish
GET  /api/v1/bots/:bot_id/conversations
GET  /api/v1/conversations/:conversation_id/messages

// ❌ 错误: 扁平化设计
GET  /api/v1/bot_config?bot_id=xxx
POST /api/v1/bot_publish
GET  /api/v1/bot_conversation
GET  /api/v1/conversation_message
```

#### 路径参数vs查询参数

```go
// ✅ 正确: 路径参数用于必需的资源ID
GET /api/v1/bots/:bot_id

// ✅ 正确: 查询参数用于可选的过滤、排序、分页
GET /api/v1/bots?status=active&sort=created_at&page=1

// ❌ 错误: 所有参数都用POST body
POST /api/v1/bots/get_by_id
Body: {"bot_id": "xxx"}
```

---

## 4️⃣ Swagger文档缺失 (422个) - P1

### 问题描述

439个API中,仅17个有Swagger注释,覆盖率3.9%。

### 优先添加文档的API (Top 50)

#### 4.1 Bot管理 (10个)

| 序号 | API路径 | 方法 | 优先级 |
|------|---------|------|--------|
| 1 | `/api/v1/bots` | POST | P0 |
| 2 | `/api/v1/bots/:bot_id` | GET | P0 |
| 3 | `/api/v1/bots` | GET | P0 |
| 4 | `/api/v1/bots/:bot_id` | PUT | P0 |
| 5 | `/api/v1/bots/:bot_id` | DELETE | P0 |
| 6 | `/api/v1/bots/:bot_id/config` | GET | P1 |
| 7 | `/api/v1/bot-types` | GET | P1 |
| 8 | `/api/v1/draft-bots` | POST | P1 |
| 9 | `/api/v1/draft-bots/:bot_id` | GET | P1 |
| 10 | `/api/v1/draft-bots/:bot_id/publish` | POST | P1 |

#### 4.2 Conversation管理 (10个)

| 序号 | API路径 | 方法 | 优先级 |
|------|---------|------|--------|
| 11 | `/api/v1/conversations` | POST | P0 |
| 12 | `/api/v1/conversations/:conversation_id` | GET | P0 |
| 13 | `/api/v1/conversations` | GET | P0 |
| 14 | `/api/v1/conversations/:conversation_id` | DELETE | P0 |
| 15 | `/api/v1/conversations/:conversation_id/messages` | GET | P0 |
| 16 | `/api/v1/conversations/:conversation_id/messages` | POST | P0 |
| 17 | `/api/v1/conversations/:conversation_id/messages/:message_id` | DELETE | P1 |
| 18 | `/api/v1/conversations/:conversation_id/clear` | POST | P1 |
| 19 | `/api/v1/conversations/:conversation_id/break` | POST | P1 |

#### 4.3 Knowledge管理 (10个)

| 序号 | API路径 | 方法 | 优先级 |
|------|---------|------|--------|
| 20 | `/api/v1/knowledge` | POST | P0 |
| 21 | `/api/v1/knowledge/:knowledge_id` | GET | P0 |
| 22 | `/api/v1/knowledge` | GET | P0 |
| 23 | `/api/v1/knowledge/:knowledge_id` | PUT | P0 |
| 24 | `/api/v1/knowledge/:knowledge_id` | DELETE | P0 |
| 25 | `/api/v1/knowledge/:knowledge_id/documents` | GET | P1 |
| 26 | `/api/v1/knowledge/:knowledge_id/documents` | POST | P1 |
| 27 | `/api/v1/documents/:document_id` | GET | P1 |
| 28 | `/api/v1/documents/:document_id` | PUT | P1 |
| 29 | `/api/v1/documents/:document_id` | DELETE | P1 |

#### 4.4 Workflow管理 (10个)

| 序号 | API路径 | 方法 | 优先级 |
|------|---------|------|--------|
| 30 | `/api/v1/workflows` | POST | P0 |
| 31 | `/api/v1/workflows/:workflow_id` | GET | P0 |
| 32 | `/api/v1/workflows` | GET | P0 |
| 33 | `/api/v1/workflows/:workflow_id` | PUT | P0 |
| 34 | `/api/v1/workflows/:workflow_id` | DELETE | P0 |
| 35 | `/api/v1/workflows/:workflow_id/execute` | POST | P1 |
| 36 | `/api/v1/workflows/:workflow_id/nodes` | GET | P1 |

#### 4.5 其他核心API (10个)

| 序号 | API路径 | 方法 | 优先级 |
|------|---------|------|--------|
| 37 | `/api/v1/databases` | POST | P1 |
| 38 | `/api/v1/databases/:database_id` | GET | P1 |
| 39 | `/api/v1/databases` | GET | P1 |
| 40 | `/api/v1/files/upload` | POST | P1 |
| 41 | `/api/v1/plugins` | GET | P1 |
| 42 | `/api/v1/plugins/:plugin_id` | GET | P1 |
| 43 | `/api/v1/routing-rules` | POST | P1 |
| 44 | `/api/v1/routing-rules/:rule_id` | GET | P1 |
| 45 | `/api/v1/search` | GET | P2 |

### Swagger注释模板

```go
// CreateBot 创建Bot
// @Summary 创建Bot
// @Description 创建一个新的AI智能体,支持配置角色、提示词、工具等
// @Tags Bot管理
// @Accept json
// @Produce json
// @Param request body bot.CreateBotRequest true "创建Bot请求"
// @Success 200 {object} bot.CreateBotResponse "创建成功"
// @Failure 400 {object} errno.ErrorResponse "参数错误"
// @Failure 401 {object} errno.ErrorResponse "未授权"
// @Failure 500 {object} errno.ErrorResponse "服务器错误"
// @Router /api/v1/bots [post]
func (h *Handler) CreateBot(ctx context.Context, c *app.RequestContext) {
    // ...
}
```

### 批量生成工具

```bash
#!/bin/bash
# scripts/generate_swagger.sh

echo "🔧 开始生成Swagger文档..."

# 步骤1: 安装swag
echo "📦 安装swag工具..."
go install github.com/swaggo/swag/cmd/swag@latest

# 步骤2: 生成文档
echo "📝 生成Swagger文档..."
cd backend
swag init -g api/main.go -o docs --parseInternal --parseDepth 1

# 步骤3: 检查生成的文档
echo "🔍 检查生成的文档..."
if [ -f "docs/docs.go" ]; then
    echo "✅ Swagger文档生成成功!"
    echo "📖 访问文档: http://localhost:8080/swagger/index.html"
else
    echo "❌ Swagger文档生成失败"
    exit 1
fi
```

### Swagger注释检查清单

每个API必须包含的注释:

- [ ] @Summary - 简短描述
- [ ] @Description - 详细描述
- [ ] @Tags - API分组
- [ ] @Accept - 请求内容类型
- [ ] @Produce - 响应内容类型
- [ ] @Param - 请求参数
- [ ] @Success - 成功响应
- [ ] @Failure - 错误响应
- [ ] @Router - API路径和方法

---

## 5️⃣ 缺少版本号 (1个) - P2

### 违规案例

**文件**: `backend/api/router/coze/api.go:403`

```go
// ❌ 当前实现 - 缺少版本号
_workflow_api.GET("/apiDetail", append(_getapidetailMw(), coze.GetApiDetail)...)

// ✅ 修复方案 - 添加v1版本号
_workflow_api.GET("/v1/api-detail", append(_getapidetailMw(), coze.GetApiDetail)...)
```

---

## 🛠️ 修复优先级和时间表

### Phase 1: 紧急修复 (Week 1-2)

**目标**: 解决P0问题

| 任务 | 数量 | 工作量 | 负责人 |
|------|------|--------|--------|
| 统一响应格式 | 601个 | 3人日 | 后端工程师 |
| 修复HTTP方法误用 | 45个 | 2人日 | 后端工程师 |

**验收标准**:
- [ ] 100% API使用httputil.BuildSuccessResp/ErrorResp
- [ ] 0个POST用于查询
- [ ] 0个POST用于删除

### Phase 2: 规范优化 (Week 3-4)

**目标**: 解决P1问题

| 任务 | 数量 | 工作量 | 负责人 |
|------|------|--------|--------|
| 重命名路由 | 38处 | 2人日 | 后端工程师 |
| 补充Swagger文档 | 50个 | 5人日 | 后端工程师 |

**验收标准**:
- [ ] 100%资源使用复数
- [ ] 100%路径使用小写+连字符
- [ ] 核心API文档覆盖率≥80%

### Phase 3: 持续改进 (Week 5-8)

**目标**: 建立长期机制

| 任务 | 工作量 | 负责人 |
|------|--------|--------|
| 制定API设计规范 | 1人日 | 架构师 |
| 添加CI检查 | 2人日 | DevOps |
| 团队培训 | 1人日 | 架构师 |

**验收标准**:
- [ ] CI自动检查通过
- [ ] 新API必须通过审查
- [ ] 团队培训完成率100%

---

## 📊 进度跟踪

### 修复进度表

| 违规类型 | 总数 | 已修复 | 进行中 | 待修复 | 完成率 |
|---------|------|--------|--------|--------|--------|
| 响应格式不统一 | 601 | 0 | 0 | 601 | 0% |
| POST查询 | 30 | 0 | 0 | 30 | 0% |
| POST删除 | 5 | 0 | 0 | 5 | 0% |
| 单数命名 | 23 | 0 | 0 | 23 | 0% |
| 驼峰命名 | 15 | 0 | 0 | 15 | 0% |
| Swagger文档 | 422 | 0 | 0 | 422 | 0% |

### 每周更新

**Week 1 (2025-01-06)**:
- [ ] 统一响应格式: 0/601 (0%)
- [ ] 修复HTTP方法: 0/45 (0%)

**Week 2 (2025-01-13)**:
- [ ] 统一响应格式: ___/601 (___%)
- [ ] 修复HTTP方法: ___/45 (___%)

**Week 3-4 (2025-01-20)**:
- [ ] 重命名路由: ___/38 (___%)
- [ ] Swagger文档: ___/50 (___%)

**Week 5-8 (2025-02-03)**:
- [ ] CI检查: 待实现
- [ ] 团队培训: 待完成

---

## 🔗 相关资源

### 工具

1. **Swagger Editor**: https://editor.swagger.io/
2. **Swagger UI**: https://swagger.io/tools/swagger-ui/
3. **Postman**: https://www.postman.com/

### 文档

1. **[OpenAPI Specification](https://swagger.io/specification/)**
2. **[RESTful API设计指南](https://restfulapi.net/)**
3. **[ZKER-API设计规范文档](./API设计规范文档.md)**

### 示例

1. **GitHub REST API**: https://docs.github.com/en/rest
2. **Stripe API**: https://stripe.com/docs/api
3. **Shopify API**: https://shopify.dev/docs/api/admin-rest

---

**文档结束** 🎉

**⚠️ 重要提醒**: 所有P0问题必须在2周内全部解决!
