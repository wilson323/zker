# 创建 API Handler Skill

## 技能描述

创建符合 Coze Studio 规范的 HTTP API Handler，包括请求处理、参数验证、响应封装和错误处理。

## 适用场景

- 需要暴露 RESTful API
- 需要处理 HTTP 请求
- 需要参数验证和绑定
- 需要统一响应格式

## 工作流程

### 1. 分析 API 需求

明确以下信息：
- **HTTP 方法**: GET/POST/PUT/PATCH/DELETE
- **路由路径**: API 路径设计
- **请求参数**: Query/Path/Body 参数
- **响应格式**: 统一响应结构
- **错误码**: 业务错误码定义

### 2. 定义 Handler 结构

**文件位置**: `api/handler/coze/{service}_handler.go` 或 `api/handler/{service}_handler.go`

```go
package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/errors"
)

// {Service}Handler {service} HTTP处理器
type {Service}Handler struct {
	// 依赖注入
	service service.{Service}Service
	// 其他依赖
	validator *validator.Validate
}

// New{Service}Handler 创建{service}处理器
func New{Service}Handler(
	service service.{Service}Service,
	validator *validator.Validate,
) *{Service}Handler {
	return &{Service}Handler{
		service:   service,
		validator: validator,
	}
}
```

### 3. 实现 Handler 方法

#### 3.1 创建资源

```go
// Create{Resource} 创建{resource}
//
// @router POST /api/v1/{resources} [Create{Resource}]
func (h *{Service}Handler) Create{Resource}(ctx context.Context, c *app.RequestContext) {
	// 1. 解析请求参数
	var req service.Create{Resource}Request
	if err := c.BindAndValidate(&req); err != nil {
		// 参数验证失败
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "请求参数不合法",
			Details: err.Error(),
		})
	}

	// 2. 调用服务层
	{resource}, err := h.service.Create{Resource}(ctx, &req)
	if err != nil {
		// 业务错误处理
		if errorx.Is(err, errno.Err{Resource}AlreadyExists) {
			return c.JSON(409, &ErrorResponse{
				Code:    ErrCode{Resource}AlreadyExists,
				Message: "{resource}已存在",
			})
		}

		// 其他错误
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	// 3. 返回成功响应
	return c.JSON(200, &SuccessResponse{
		Code:    0,
		Message: "success",
		Data:    {resource},
	})
}
```

#### 3.2 获取单个资源

```go
// Get{Resource}ByID 获取{resource}详情
//
// @router GET /api/v1/{resources}/:id [Get{Resource}ByID]
func (h *{Service}Handler) Get{Resource}ByID(ctx context.Context, c *app.RequestContext) {
	// 1. 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "无效的 ID",
		})
	}

	// 2. 调用服务层
	{resource}, err := h.service.Get{Resource}ByID(ctx, id)
	if err != nil {
		if errorx.Is(err, errno.Err{Resource}NotFound) {
			return c.JSON(404, &ErrorResponse{
				Code:    ErrCode{Resource}NotFound,
				Message: "{resource}不存在",
			})
		}
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	// 3. 返回成功响应
	return c.JSON(200, &SuccessResponse{
		Code:    0,
		Message: "success",
		Data:    {resource},
	})
}
```

#### 3.3 更新资源

```go
// Update{Resource} 更新{resource}
//
// @router PUT /api/v1/{resources}/:id [Update{Resource}]
func (h *{Service}Handler) Update{Resource}(ctx context.Context, c *app.RequestContext) {
	// 1. 解析参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "无效的 ID",
		})
	}

	var req service.Update{Resource}Request
	if err := c.BindAndValidate(&req); err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "请求参数不合法",
			Details: err.Error(),
		})
	}

	req.{Resource}ID = id

	// 2. 调用服务层
	if err := h.service.Update{Resource}(ctx, &req); err != nil {
		if errorx.Is(err, errno.Err{Resource}NotFound) {
			return c.JSON(404, &ErrorResponse{
				Code:    ErrCode{Resource}NotFound,
				Message: "{resource}不存在",
			})
		}
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	// 3. 返回成功响应
	return c.JSON(200, &SuccessResponse{
		Code:    0,
		Message: "success",
	})
}
```

#### 3.4 删除资源

```go
// Delete{Resource} 删除{resource}
//
// @router DELETE /api/v1/{resources}/:id [Delete{Resource}]
func (h *{Service}Handler) Delete{Resource}(ctx context.Context, c *app.RequestContext) {
	// 1. 解析路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "无效的 ID",
		})
	}

	// 2. 调用服务层
	if err := h.service.Delete{Resource}(ctx, id); err != nil {
		if errorx.Is(err, errno.Err{Resource}NotFound) {
			return c.JSON(404, &ErrorResponse{
				Code:    ErrCode{Resource}NotFound,
				Message: "{resource}不存在",
			})
		}
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	// 3. 返回成功响应
	return c.JSON(200, &SuccessResponse{
		Code:    0,
		Message: "success",
	})
}
```

#### 3.5 列表查询

```go
// List{Resources} 获取{resource}列表
//
// @router GET /api/v1/{resources} [List{Resources}]
func (h *{Service}Handler) List{Resources}(ctx context.Context, c *app.RequestContext) {
	// 1. 解析查询参数
	var req List{Resource}Query
	if err := c.BindAndValidate(&req); err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "请求参数不合法",
		})
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 2. 调用服务层
	{resources}, total, err := h.service.List{Resources}(ctx, &req)
	if err != nil {
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	// 3. 返回成功响应
	return c.JSON(200, &ListResponse{
		Code:    0,
		Message: "success",
		Data: ListData{
			Items:      {resources},
			TotalCount: total,
			Page:       req.Page,
			PageSize:   req.PageSize,
		},
	})
}
```

### 4. 请求参数定义

```go
// List{Resource}Query 列表查询参数
type List{Resource}Query struct {
	Page     int    `query:"page" vd:"$>=1"`
	PageSize int    `query:"page_size" vd:"$>=1&&<=100"`
	Keyword  string `query:"keyword"`
	Status   string `query:"status"`
	SpaceID  int64  `query:"space_id"`
}

// Create{Resource}Request 创建{resource}请求
type Create{Resource}Request struct {
	Name    string `json:"name" binding:"required,max=100"`
	Email   string `json:"email" binding:"required,email"`
	SpaceID int64  `json:"space_id" binding:"required"`
}

// Update{Resource}Request 更新{resource}请求
type Update{Resource}Request struct {
	Name    *string `json:"name" binding:"omitempty,max=100"`
	Email   *string `json:"email" binding:"omitempty,email"`
	Avatar  *string `json:"avatar" binding:"omitempty,url"`
}
```

### 5. 响应结构定义

```go
// SuccessResponse 成功响应
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ListData 列表数据
type ListData struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}
```

### 6. 错误码定义

```go
const (
	// 通用错误 1xxxx
	ErrCodeSuccess         = 0
	ErrCodeInvalidParams   = 10001
	ErrCodeUnauthorized    = 10002
	ErrCodeForbidden       = 10003
	ErrCodeNotFound        = 10004
	ErrCodeInternalError   = 10005

	// {Resource}错误 2xxxx
	ErrCode{Resource}NotFound      = 20001
	ErrCode{Resource}AlreadyExists  = 20002
	ErrCodeInvalid{Resource}Email  = 20003
)
```

### 7. 路由注册

```go
// router/{service}_router.go
package router

import (
	"{project}/api/handler/coze"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Register{Service}Routes(s *server.Hertz, handler *coze.{Service}Handler) {
	// {Resource}路由组
	{resource}Group := s.Group("/api/v1/{resources}")

	// 注册路由
	{resource}Group.POST("", handler.Create{Resource})
	{resource}Group.GET("", handler.List{Resources})
	{resource}Group.GET("/:id", handler.Get{Resource}ByID)
	{resource}Group.PUT("/:id", handler.Update{Resource})
	{resource}Group.DELETE("/:id", handler.Delete{Resource})
}
```

### 8. Handler 测试

```go
package coze_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"{project}/api/handler/coze"
	"{project}/domain/{domain}/service"
)

// Mock Service
type Mock{Service}Service struct {
	mock.Mock
}

func (m *Mock{Service}Service) Create{Resource}(ctx context.Context, req *service.Create{Resource}Request) (*service.{Resource}, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.{Resource}), args.Error(1)
}

func TestCreate{Resource}(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setup          func(*Mock{Service}Service)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			requestBody: map[string]interface{}{
				"name":    "Test",
				"email":   "test@example.com",
				"spaceID": float64(1),
			},
			setup: func(m *Mock{Service}Service) {
				m.On("Create{Resource}", mock.Anything, mock.Anything).
					Return(&service.{Resource}{ID: 123, Name: "Test"}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "success",
			},
		},
		{
			name: "invalid params",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			setup:          func(m *Mock{Service}Service) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(Mock{Service}Service)
			tt.setup(mockService)

			handler := coze.New{Service}Handler(mockService, nil)

			h := server.Default()
			h.POST("/api/v1/{resources}", handler.Create{Resource})

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/{resources}", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			h.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}
```

## 必须遵循的规范

### ✅ 命名规范

```go
// ✅ 正确的命名
type {Service}Handler struct {}
func New{Service}Handler() *{Service}Handler {}
func (h *{Service}Handler) Create{Resource}() {}
func (h *{Service}Handler) Get{Resource}ByID() {}

// ❌ 错误的命名
type {Service} struct {}                    ← 应使用 Handler 后缀
func New{Resource}Handler() {}               ← 应使用服务名
func (h *{Service}Handler) create() {}        ← 导出方法应大写
func (h *{Service}Handler) Get() {}           ← 应包含资源名
```

### ✅ HTTP 方法与路由规范

```go
// ✅ 正确的路由设计
GET    /api/v1/{resources}           // 列表
GET    /api/v1/{resources}/:id       // 详情
POST   /api/v1/{resources}           // 创建
PUT    /api/v1/{resources}/:id       // 更新（全部）
PATCH  /api/v1/{resources}/:id       // 更新（部分）
DELETE /api/v1/{resources}/:id       // 删除

// ✅ 正确的特殊操作
POST   /api/v1/{resources}/:id/publish   // 发布
POST   /api/v1/{resources}/:id/clone     // 克隆
POST   /api/v1/{resources}/:id/export    // 导出

// ❌ 错误的路由设计
GET    /api/v1/get{resources}            ← 不应包含动词
POST   /api/v1/create{resource}          ← 不应包含动词
GET    /api/v1/{resource}/:id            ← 应使用复数
GET    /api/v1/{resources}/get/:id      ← 路径不应包含动作
```

### ✅ 参数验证规范

```go
// ✅ 正确的参数验证
type Request struct {
	Name   string `json:"name" binding:"required,max=100"`
	Email  string `json:"email" binding:"required,email"`
	Age    int    `json:"age" binding:"required,gte=0,lte=150"`
}

// ❌ 错误的参数验证
type Request struct {
	Name   string `json:"name"`              ← 缺少验证规则
	Email  string `json:"email" binding:"email"` ← 缺少 required
	Age    int    `json:"age" binding:"min=0"`  ← 应使用 gte
}
```

### ✅ 错误处理规范

```go
// ✅ 正确的错误处理
if err := c.BindAndValidate(&req); err != nil {
	return c.JSON(400, &ErrorResponse{
		Code: ErrCodeInvalidParams,
		Message: "请求参数不合法",
	})
}

if err := h.service.Create{Resource}(ctx, &req); err != nil {
	if errorx.Is(err, errno.Err{Resource}AlreadyExists) {
		return c.JSON(409, &ErrorResponse{
			Code: ErrCode{Resource}AlreadyExists,
			Message: "{resource}已存在",
		})
	}
	return c.JSON(500, &ErrorResponse{
		Code: ErrCodeInternalError,
		Message: "内部服务错误",
	})
}

// ❌ 错误的错误处理
if err != nil {
	return c.JSON(500, err)              ← 未区分错误类型
	return c.String(500, err.Error())     ← 应使用 JSON 响应
}
```

### ✅ 响应格式规范

```go
// ✅ 正确的成功响应
{
	"code": 0,
	"message": "success",
	"data": {
		"id": 123,
		"name": "Test"
	}
}

// ✅ 正确的列表响应
{
	"code": 0,
	"message": "success",
	"data": {
		"items": [...],
		"total_count": 100,
		"page": 1,
		"page_size": 20
	}
}

// ✅ 正确的错误响应
{
	"code": 20001,
	"message": "resource not found",
	"details": "详细错误信息"
}

// ❌ 错误的响应格式
{ "data": {...} }                  ← 缺少 code 和 message
{ "error": "not found" }            ← 应使用 code 和 message
{ "success": false, "error": "..." }  ← 应使用统一格式
```

## 输出检查清单

创建完成后，必须确保：

- [ ] Handler 命名规范（{Service}Handler）
- [ ] 方法命名规范（动词-资源名）
- [ ] 路由符合 RESTful 规范
- [ ] 请求参数验证完整
- [ ] 响应格式统一
- [ ] 错误码定义完整
- [ ] 使用 context.Context
- [ ] 错误处理完善
- [ ] 包含单元测试
- [ ] 测试覆盖正常和异常场景
- [ ] 通过 gofmt 检查
- [ ] 通过 golangci-lint 检查

## 示例输出

**用户输入**：
```
创建用户管理 API Handler，包括 CRUD 操作
```

**AI 输出**：

**文件**: `api/handler/coze/user_service.go`

```go
package coze

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"

	"{project}/domain/user/service"
	"{project}/pkg/errorx"
	"{project}/types/errno"
)

type UserServiceHandler struct {
	service   service.UserService
	validator *validator.Validate
}

func NewUserServiceHandler(
	svc service.UserService,
	validator *validator.Validate,
) *UserServiceHandler {
	return &UserServiceHandler{
		service:   svc,
		validator: validator,
	}
}

// CreateUser 创建用户
// @router POST /api/v1/users [CreateUser]
func (h *UserServiceHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
	var req service.CreateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "请求参数不合法",
			Details: err.Error(),
		})
	}

	user, err := h.service.CreateUser(ctx, &req)
	if err != nil {
		if errorx.Is(err, errno.ErrUserAlreadyExists) {
			return c.JSON(409, &ErrorResponse{
				Code:    ErrCodeUserAlreadyExists,
				Message: "用户已存在",
			})
		}
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	return c.JSON(200, &SuccessResponse{
		Code:    0,
		Message: "success",
		Data:    user,
	})
}

// GetUserByID 获取用户详情
// @router GET /api/v1/users/:id [GetUserByID]
func (h *UserServiceHandler) GetUserByID(ctx context.Context, c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(400, &ErrorResponse{
			Code:    ErrCodeInvalidParams,
			Message: "无效的 ID",
		})
	}

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		if errorx.Is(err, errno.ErrUserNotFound) {
			return c.JSON(404, &ErrorResponse{
				Code:    ErrCodeUserNotFound,
				Message: "用户不存在",
			})
		}
		return c.JSON(500, &ErrorResponse{
			Code:    ErrCodeInternalError,
			Message: "内部服务错误",
		})
	}

	return c.JSON(200, &SuccessResponse{
		Code:    0,
		Message: "success",
		Data:    user,
	})
}
```

## 注意事项

1. **参数验证**: 使用 binding tag 进行参数验证
2. **错误处理**: 统一错误码和错误响应格式
3. **日志记录**: 记录关键操作和错误日志
4. **性能监控**: 添加性能监控指标
5. **安全性**: 验证权限、防止注入攻击
6. **API 文档**: 添加完整的 API 注释
7. **版本控制**: 路径包含版本号
8. **限流**: 对敏感操作添加限流
