# 研发 B - 后端工程师 8 周开发计划 v1.0

> **角色定位**: 后端工程师 - 负责统一错误码系统、性能测试框架、监控和日志、分布式追踪等工程质量保障
>
> **工作目标**: 构建完善的工程基础设施，确保系统可观测性、高性能、高可靠性
>
> **核心原则**: 质量优先、可观测性、自动化、数据驱动

---

## 📋 个人职责概述

### 核心负责模块

```
backend/
├── types/errno/              ✅ 专属负责 - 错误码定义
│   ├── errno.go
│   ├── errors.go
│   └── error_codes/
├── tests/performance/        ✅ 专属负责 - 性能测试
│   ├── bot_test.go
│   ├── routing_test.go
│   └── quota_test.go
├── infra/monitoring/         ✅ 专属负责 - 监控
│   ├── metrics/
│   ├── prometheus/
│   └── grafana/
├── infra/logging/            ✅ 专属负责 - 日志
│   ├── logger.go
│   ├── formatter.go
│   └── middleware.go
├── infra/tracing/            ✅ 专属负责 - 分布式追踪
│   ├── tracer.go
│   └── middleware.go
└── infra/cache/              ✅ 专属负责 - 缓存
    ├── redis.go
    └── middleware.go
```

### 协作接口

| 协作对象 | 协作内容 | 接口定义位置 | 依赖关系 |
|---------|---------|------------|---------|
| **研发 A** | 错误码使用、性能测试对象、监控埋点 | `types/errno/`、`tests/performance/` | 研发 B 定义 → 研发 A 使用 |
| **研发 C** | 前端错误码映射、监控数据展示 | `frontend/packages/api-client/` | 研发 B 提供 API → 研发 C 集成 |
| **研发 D** | 监控服务部署、日志采集、性能测试环境 | `docker/`、`scripts/` | 研发 B 提供配置 → 研发 D 部署 |

---

## 🎯 8 周详细开发计划

### Week 1-2: 统一错误码系统 + API 规范

#### Week 1: 错误码定义系统

**目标**: 建立完整的错误码体系，覆盖所有业务场景

##### Day 1-2: 错误码基础架构

**错误码结构设计** (`types/errno/errors.go`):
```go
// errno/errors.go
package errno

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// Errno 统一错误码结构
type Errno struct {
    Code       string                 `json:"code"`                    // 错误码
    Message    string                 `json:"message"`                 // 错误消息（英文）
    MessageZH  string                 `json:"message_zh,omitempty"`   // 错误消息（中文）
    HTTPStatus int                    `json:"-"`                       // HTTP 状态码
    Details    map[string]interface{} `json:"details,omitempty"`       // 错误详情
    RequestID  string                 `json:"request_id,omitempty"`   // 请求ID
    Timestamp  string                 `json:"timestamp,omitempty"`     // 时间戳
}

// Error 实现 error 接口
func (e *Errno) Error() string {
    if e.MessageZH != "" {
        return fmt.Sprintf("[%s] %s ( %s )", e.Code, e.Message, e.MessageZH)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// WithDetails 添加错误详情
func (e *Errno) WithDetails(details map[string]interface{}) *Errno {
    e.Details = details
    return e
}

// WithRequestID 添加请求ID
func (e *Errno) WithRequestID(requestID string) *Errno {
    e.RequestID = requestID
    return e
}

// ToJSON 转换为 JSON
func (e *Errno) ToJSON() []byte {
    data, _ := json.Marshal(e)
    return data
}

// HTTPStatus 获取 HTTP 状态码
func (e *Errno) GetHTTPStatus() int {
    if e.HTTPStatus != 0 {
        return e.HTTPStatus
    }
    return http.StatusInternalServerError
}

// New 创建新的错误码
func New(code, message, messageZH string, httpStatus int) *Errno {
    return &Errno{
        Code:       code,
        Message:    message,
        MessageZH:  messageZH,
        HTTPStatus: httpStatus,
    }
}

// Wrap 包装错误
func Wrap(err error, code, message, messageZH string) *Errno {
    return &Errno{
        Code:       code,
        Message:    message,
        MessageZH:  messageZH,
        HTTPStatus: http.StatusInternalServerError,
        Details: map[string]interface{}{
            "underlying_error": err.Error(),
        },
    }
}
```

**错误码常量定义** (`types/errno/error_codes/tenant.go`):
```go
// error_codes/tenant.go
package error_codes

const (
    // 租户相关错误码 (TENANT_*)
    TENANT_NOT_FOUND = "TENANT_NOT_FOUND"
    TENANT_ALREADY_EXISTS = "TENANT_ALREADY_EXISTS"
    TENANT_SUSPENDED = "TENANT_SUSPENDED"
    TENANT_DELETED = "TENANT_DELETED"
)

// TenantErrors 租户错误码实例
var (
    ErrTenantNotFound = &errno.Errno{
        Code:       TENANT_NOT_FOUND,
        Message:    "Tenant not found",
        MessageZH:  "租户不存在",
        HTTPStatus: 404,
    }
    ErrTenantAlreadyExists = &errno.Errno{
        Code:       TENANT_ALREADY_EXISTS,
        Message:    "Tenant already exists",
        MessageZH:  "租户已存在",
        HTTPStatus: 409,
    }
    ErrTenantSuspended = &errno.Errno{
        Code:       TENANT_SUSPENDED,
        Message:    "Tenant is suspended",
        MessageZH:  "租户已暂停",
        HTTPStatus: 403,
    }
)
```

**错误码分类** (`types/errno/error_codes/categories.go`):
```go
// error_codes/categories.go
package error_codes

// 错误码分类
const (
    // 通用错误 (1xxx)
    SUCCESS = "SUCCESS"
    UNKNOWN_ERROR = "UNKNOWN_ERROR"
    INVALID_PARAMS = "INVALID_PARAMS"
    UNAUTHORIZED = "UNAUTHORIZED"
    FORBIDDEN = "FORBIDDEN"
    NOT_FOUND = "NOT_FOUND"
    CONFLICT = "CONFLICT"
    RATE_LIMIT_EXCEEDED = "RATE_LIMIT_EXCEEDED"
    INTERNAL_ERROR = "INTERNAL_ERROR"
    SERVICE_UNAVAILABLE = "SERVICE_UNAVAILABLE"

    // 租户错误 (2xxx)
    TENANT_NOT_FOUND = "TENANT_NOT_FOUND"
    TENANT_ALREADY_EXISTS = "TENANT_ALREADY_EXISTS"
    TENANT_SUSPENDED = "TENANT_SUSPENDED"

    // 配额错误 (3xxx)
    QUOTA_EXCEEDED = "QUOTA_EXCEEDED"
    QUOTA_INVALID = "QUOTA_INVALID"
    QUOTA_RESET_FAILED = "QUOTA_RESET_FAILED"

    // 权限错误 (4xxx)
    PERMISSION_DENIED = "PERMISSION_DENIED"
    ROLE_NOT_FOUND = "ROLE_NOT_FOUND"
    ROLE_ALREADY_EXISTS = "ROLE_ALREADY_EXISTS"
    DATA_PERMISSION_DENIED = "DATA_PERMISSION_DENIED"
    FIELD_PERMISSION_DENIED = "FIELD_PERMISSION_DENIED"

    // Bot 错误 (5xxx)
    BOT_NOT_FOUND = "BOT_NOT_FOUND"
    BOT_ALREADY_EXISTS = "BOT_ALREADY_EXISTS"
    BOT_INVALID_CONFIG = "BOT_INVALID_CONFIG"
    BOT_VERSION_CONFLICT = "BOT_VERSION_CONFLICT"

    // 工作流错误 (6xxx)
    WORKFLOW_NOT_FOUND = "WORKFLOW_NOT_FOUND"
    WORKFLOW_INVALID = "WORKFLOW_INVALID"
    WORKFLOW_EXECUTION_FAILED = "WORKFLOW_EXECUTION_FAILED"
    NODE_EXECUTION_FAILED = "NODE_EXECUTION_FAILED"

    // 知识库错误 (7xxx)
    KNOWLEDGE_NOT_FOUND = "KNOWLEDGE_NOT_FOUND"
    DOCUMENT_UPLOAD_FAILED = "DOCUMENT_UPLOAD_FAILED"
    DOCUMENT_PARSE_FAILED = "DOCUMENT_PARSE_FAILED"
    VECTOR_SEARCH_FAILED = "VECTOR_SEARCH_FAILED"

    // 路由错误 (8xxx)
    ROUTING_NO_MATCH = "ROUTING_NO_MATCH"
    ROUTING_INVALID_RULE = "ROUTING_INVALID_RULE"
    ROUTING_SERVICE_UNHEALTHY = "ROUTING_SERVICE_UNHEALTHY"

    // 订阅和计费错误 (9xxx)
    SUBSCRIPTION_EXPIRED = "SUBSCRIPTION_EXPIRED"
    SUBSCRIPTION_NOT_FOUND = "SUBSCRIPTION_NOT_FOUND"
    BILLING_FAILED = "BILLING_FAILED"
    PAYMENT_REQUIRED = "PAYMENT_REQUIRED"
)
```

##### Day 3-5: 错误码使用规范和工具

**错误码中间件** (`api/middleware/error_middleware.go`):
```go
// middleware/error_middleware.go
package middleware

import (
    "bytes"
    "encoding/json"
    "io"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-studio/types/errno"
    "github.com/google/uuid"
)

// ErrorHandler 错误处理中间件
func ErrorHandler() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 生成请求 ID
        requestID := uuid.New().String()
        c.Set("request_id", requestID)

        // 执行请求
        c.Next(ctx)

        // 处理错误
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            handleErr := err.Err

            // 转换为 Errno
            var e *errno.Errno
            switch t := handleErr.(type) {
            case *errno.Errno:
                e = t
            default:
                e = &errno.Errno{
                    Code:       errno.INTERNAL_ERROR,
                    Message:    "Internal server error",
                    MessageZH:  "服务器内部错误",
                    HTTPStatus: 500,
                    Details: map[string]interface{}{
                        "underlying_error": t.Error(),
                    },
                }
            }

            // 添加请求 ID 和时间戳
            e.WithRequestID(requestID)
            e.Timestamp = time.Now().Format(time.RFC3339)

            // 记录错误日志
            logError(ctx, c, e)

            // 返回错误响应
            c.JSON(e.GetHTTPStatus(), e)
        }
    }
}

// logError 记录错误日志
func logError(ctx context.Context, c *app.RequestContext, e *errno.Errno) {
    // 读取请求体（用于调试）
    var body string
    if c.Request.Body() != nil {
        bodyBytes, _ := io.ReadAll(c.Request.Body())
        c.Request.SetBody(bytes.NewReader(bodyBytes)) // 恢复 body
        body = string(bodyBytes)
    }

    // 结构化日志
    logEntry := map[string]interface{}{
        "level":       "error",
        "timestamp":   e.Timestamp,
        "request_id":  e.RequestID,
        "code":        e.Code,
        "message":     e.Message,
        "message_zh":  e.MessageZH,
        "http_method": c.Request.Method(),
        "http_path":   c.Request.URI().Path(),
        "http_status": e.HTTPStatus,
        "client_ip":   c.ClientIP(),
        "user_agent":  c.Request.Header.Get("User-Agent"),
        "request_body": body,
        "details":     e.Details,
    }

    // 输出到日志（使用 JSON 格式）
    logJSON, _ := json.Marshal(logEntry)
    println(string(logJSON))
}
```

**错误码测试工具** (`types/errno/error_codes/test_tool.go`):
```go
// test_tool/main.go
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/coze-studio/types/errno"
    "github.com/coze-studio/types/errno/error_codes"
)

func main() {
    // 生成错误码文档
    generateErrorDoc()
}

func generateErrorDoc() {
    errors := []*errno.Errno{
        error_codes.ErrTenantNotFound,
        error_codes.ErrTenantAlreadyExists,
        error_codes.ErrTenantSuspended,
        // ... 所有错误码
    }

    // 生成 Markdown 表格
    fmt.Println("# 错误码清单\n")
    fmt.Println("| 错误码 | 英文消息 | 中文消息 | HTTP 状态码 |")
    fmt.Println("|--------|---------|---------|------------|")

    for _, e := range errors {
        fmt.Printf("| %s | %s | %s | %d |\n",
            e.Code, e.Message, e.MessageZH, e.HTTPStatus)
    }

    // 生成 JSON 文件（供前端使用）
    errorJSON, _ := json.MarshalIndent(errors, "", "  ")
    os.WriteFile("error_codes.json", errorJSON, 0644)
}
```

#### Week 2: API 规范和文档

**目标**: 完善 OpenAPI 规范，生成 API 文档

##### Day 1-3: OpenAPI 规范定义

**OpenAPI 规范文件** (`openapi/zker-api-v1-tenant.yaml`):
```yaml
openapi: 3.0.0
info:
  title: ZKER Tenant API
  version: 1.0.0
  description: ZKER 租户管理 API
  contact:
    name: API Support
    email: api-support@zker.com

servers:
  - url: http://localhost:8080/api/v1
    description: 开发环境
  - url: https://api.zker.com/api/v1
    description: 生产环境

tags:
  - name: tenants
    description: 租户管理
  - name: subscriptions
    description: 订阅管理
  - name: quotas
    description: 配额管理

paths:
  /tenants:
    get:
      summary: 获取租户列表
      tags: [tenants]
      parameters:
        - name: page_token
          in: query
          schema:
            type: string
          description: 分页令牌
        - name: page_size
          in: query
          schema:
            type: integer
            default: 20
            minimum: 1
            maximum: 100
          description: 每页数量
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TenantListResponse'
        '400':
          description: 参数错误
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: 未认证
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
    post:
      summary: 创建租户
      tags: [tenants]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateTenantRequest'
      responses:
        '201':
          description: 创建成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TenantResponse'
        '400':
          description: 参数错误
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '409':
          description: 租户已存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /tenants/{tenant_id}:
    get:
      summary: 获取租户详情
      tags: [tenants]
      parameters:
        - name: tenant_id
          in: path
          required: true
          schema:
            type: string
          description: 租户 ID
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TenantResponse'
        '404':
          description: 租户不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
    put:
      summary: 更新租户
      tags: [tenants]
      parameters:
        - name: tenant_id
          in: path
          required: true
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateTenantRequest'
      responses:
        '200':
          description: 更新成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/TenantResponse'
        '404':
          description: 租户不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
    delete:
      summary: 删除租户
      tags: [tenants]
      parameters:
        - name: tenant_id
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: 删除成功
        '404':
          description: 租户不存在
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  schemas:
    Tenant:
      type: object
      properties:
        tenant_id:
          type: string
          description: 租户 ID
        tenant_name:
          type: string
          description: 租户名称
        tenant_type:
          type: string
          enum: [individual, team, enterprise]
          description: 租户类型
        status:
          type: string
          enum: [active, suspended, deleted]
          description: 状态
        subscription_tier:
          type: string
          enum: [free, pro, enterprise]
          description: 订阅等级
        created_at:
          type: string
          format: date-time
          description: 创建时间
        updated_at:
          type: string
          format: date-time
          description: 更新时间

    CreateTenantRequest:
      type: object
      required:
        - tenant_name
        - tenant_type
      properties:
        tenant_name:
          type: string
          minLength: 1
          maxLength: 200
          description: 租户名称
        tenant_type:
          type: string
          enum: [individual, team, enterprise]
          description: 租户类型

    UpdateTenantRequest:
      type: object
      properties:
        tenant_name:
          type: string
          minLength: 1
          maxLength: 200
        status:
          type: string
          enum: [active, suspended, deleted]

    TenantResponse:
      type: object
      properties:
        code:
          type: string
          example: SUCCESS
        message:
          type: string
          example: Operation successful
        message_zh:
          type: string
          example: 操作成功
        message_en:
          type: string
          example: Operation successful
        data:
          $ref: '#/components/schemas/Tenant'
        request_id:
          type: string
          description: 请求 ID
        timestamp:
          type: string
          format: date-time

    TenantListResponse:
      type: object
      properties:
        code:
          type: string
          example: SUCCESS
        message:
          type: string
        data:
          type: object
          properties:
            tenants:
              type: array
              items:
                $ref: '#/components/schemas/Tenant'
            next_page_token:
              type: string
            total_count:
              type: integer
        request_id:
          type: string
        timestamp:
          type: string
          format: date-time

    ErrorResponse:
      type: object
      properties:
        code:
          type: string
          example: TENANT_NOT_FOUND
        message:
          type: string
          example: Tenant not found
        message_zh:
          type: string
          example: 租户不存在
        message_en:
          type: string
          example: Tenant not found
        details:
          type: object
        request_id:
          type: string
        timestamp:
          type: string
          format: date-time
```

##### Day 4-5: API 文档生成和验证

**API 文档生成工具** (`scripts/generate_api_docs.sh`):
```bash
#!/bin/bash
# 生成 API 文档

# 1. 使用 swagger 生成文档
docker run --rm -v ${PWD}:/local \
  swaggerapi/swagger-codegen-cli-v3 generate \
  -i /local/openapi/zker-api-v1-tenant.yaml \
  -g markdown \
  -o /local/docs/api/tenant

# 2. 生成 TypeScript 类型定义
docker run --rm -v ${PWD}:/local \
  openapitools/openapi-generator-cli generate \
  -i /local/openapi/zker-api-v1-tenant.yaml \
  -g typescript-axios \
  -o /local/frontend/packages/api-client/src/gen/tenant

echo "API 文档生成完成"
```

**API 契约测试** (`tests/api/contract_test.go`):
```go
// api/contract_test.go
package api_test

import (
    "context"
    "encoding/json"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/coze-studio/api/v1/tenant"
)

func TestTenantAPI_Contract(t *testing.T) {
    // 1. 创建租户
    createReq := map[string]interface{}{
        "tenant_name": "Test Tenant",
        "tenant_type": "enterprise",
    }

    resp := callAPI("POST", "/api/v1/tenants", createReq)

    // 验证响应格式
    assert.Equal(t, 201, resp.StatusCode)

    var result map[string]interface{}
    json.Unmarshal(resp.Body, &result)

    // 验证必填字段
    assert.Contains(t, result, "code")
    assert.Contains(t, result, "message")
    assert.Contains(t, result, "data")
    assert.Contains(t, result, "request_id")
    assert.Contains(t, result, "timestamp")

    // 验证数据结构
    data := result["data"].(map[string]interface{})
    assert.Contains(t, data, "tenant_id")
    assert.Contains(t, data, "tenant_name")
    assert.Contains(t, data, "tenant_type")
    assert.Contains(t, data, "status")
    assert.Contains(t, data, "created_at")

    // 2. 获取租户列表
    listResp := callAPI("GET", "/api/v1/tenants", nil)
    assert.Equal(t, 200, listResp.StatusCode)
}

// 验证错误响应格式
func TestTenantAPI_ErrorResponse(t *testing.T) {
    resp := callAPI("GET", "/api/v1/tenants/nonexistent-id", nil)

    assert.Equal(t, 404, resp.StatusCode)

    var errResp map[string]interface{}
    json.Unmarshal(resp.Body, &errResp)

    // 验证错误响应必填字段
    assert.Contains(t, errResp, "code")
    assert.Contains(t, errResp, "message")
    assert.Contains(t, errResp, "request_id")
    assert.Contains(t, errResp, "timestamp")

    // 验证错误码
    assert.Equal(t, "TENANT_NOT_FOUND", errResp["code"])
}
```

---

### Week 3-4: 性能测试框架

#### Week 3: 性能测试工具集成

**目标**: 集成 JMeter/K6，建立性能测试基准

##### Day 1-3: K6 性能测试脚本

**K6 测试脚本** (`tests/performance/tenant_load_test.js`):
```javascript
// tests/performance/tenant_load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const latency = new Trend('latency');

// 测试配置
export const options = {
    stages: [
        { duration: '1m', target: 10 },   // 1分钟爬坡到10用户
        { duration: '3m', target: 10 },   // 维持10用户3分钟
        { duration: '1m', target: 50 },   // 1分钟爬坡到50用户
        { duration: '3m', target: 50 },   // 维持50用户3分钟
        { duration: '1m', target: 100 },  // 1分钟爬坡到100用户
        { duration: '5m', target: 100 },  // 维持100用户5分钟
        { duration: '2m', target: 0 },    // 2分钟降到0
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'],  // 95% 请求响应时间 < 500ms
        http_req_failed: ['rate<0.01'],    // 错误率 < 1%
        errors: ['rate<0.01'],
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export function setup() {
    // 创建测试租户
    const createResp = http.post(
        `${BASE_URL}/api/v1/tenants`,
        JSON.stringify({
            tenant_name: 'Performance Test Tenant',
            tenant_type: 'enterprise',
        }),
        {
            headers: { 'Content-Type': 'application/json' },
        }
    );

    const tenant = JSON.parse(createResp.body).data;
    return { tenant_id: tenant.tenant_id };
}

export default function (data) {
    // 1. 获取租户列表
    const listResp = http.get(`${BASE_URL}/api/v1/tenants`);

    check(listResp, {
        'list status is 200': (r) => r.status === 200,
        'list response time < 500ms': (r) => r.timings.duration < 500,
    }) || errorRate.add(1);

    latency.add(listResp.timings.duration);

    // 2. 获取租户详情
    const detailResp = http.get(
        `${BASE_URL}/api/v1/tenants/${data.tenant_id}`
    );

    check(detailResp, {
        'detail status is 200': (r) => r.status === 200,
        'detail has tenant_id': (r) => {
            const body = JSON.parse(r.body);
            return body.data.tenant_id === data.tenant_id;
        },
    }) || errorRate.add(1);

    sleep(1); // 每个VU每秒执行一次
}

export function teardown(data) {
    // 清理测试数据
    http.del(`${BASE_URL}/api/v1/tenants/${data.tenant_id}`);
}
```

**配额检查性能测试** (`tests/performance/quota_load_test.js`):
```javascript
// tests/performance/quota_load_test.js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
    stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<100'],  // 配额检查应该很快
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TENANT_ID = __ENV.TENANT_ID || 'test-tenant-id';

export default function () {
    // 测试配额检查性能
    const checkResp = http.get(
        `${BASE_URL}/api/v1/tenants/${TENANT_ID}/quotas/bots/check`
    );

    check(checkResp, {
        'quota check status is 200': (r) => r.status === 200,
        'quota check time < 50ms': (r) => r.timings.duration < 50,
    });
}
```

##### Day 4-5: 性能基准测试

**基准测试用例** (`tests/performance/benchmark_test.go`):
```go
// performance/benchmark_test.go
package performance

import (
    "context"
    "testing"
    "time"

    "github.com/coze-studio/domain/tenant/service"
    "github.com/coze-studio/infra/db"
)

func BenchmarkQuotaCheck(b *testing.B) {
    // Setup
    db := setupTestDB()
    quotaRepo :=.NewQuotaRepository(db)
    quotaService := service.NewQuotaService(quotaRepo)
    ctx := context.Background()

    // Reset timer
    b.ResetTimer()

    // Run benchmark
    for i := 0; i < b.N; i++ {
        quotaService.CheckQuota(ctx, "test-tenant-id", "bots", 1)
    }
}

func BenchmarkPermissionCheck(b *testing.B) {
    // Setup
    db := setupTestDB()
    permChecker := setupPermissionChecker(db)
    ctx := context.Background()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        permChecker.CheckDataPermission(ctx, "user-id", "tenant-id", "bots", "bot-id")
    }
}

func BenchmarkRoutingDecision(b *testing.B) {
    // Setup
    router := setupRoutingEngine()
    ctx := context.Background()
    input := &MatchInput{
        UserInput: "我要创建一个聊天机器人",
        TenantID:  "test-tenant-id",
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        router.Route(ctx, input)
    }
}

// 并发基准测试
func BenchmarkQuotaCheck_Parallel(b *testing.B) {
    db := setupTestDB()
    quotaRepo :=.NewQuotaRepository(db)
    quotaService := service.NewQuotaService(quotaRepo)
    ctx := context.Background()

    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            quotaService.CheckQuota(ctx, "test-tenant-id", "bots", 1)
        }
    })
}
```

#### Week 4: 性能优化和监控

**目标**: 优化性能瓶颈，建立监控大盘

##### Day 1-3: 性能分析工具

**pprof 集成** (`infra/monitoring/pprof.go`):
```go
// monitoring/pprof.go
package monitoring

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"

    "github.com/cloudwego/hertz/pkg/app"
)

// StartPprofServer 启动 pprof 服务器
func StartPprofServer(port int) {
    mux := http.NewServeMux()

    // pprof 端点
    mux.HandleFunc("/debug/pprof/", http.Index)
    mux.HandleFunc("/debug/pprof/cmdline", http.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", http.Profile)
    mux.HandleFunc("/debug/pprof/symbol", http.Symbol)
    mux.HandleFunc("/debug/pprof/trace", http.Trace)

    go func() {
        fmt.Printf("pprof server listening on :%d\n", port)
        http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
    }()
}

// Hertz 中间件
func PprofMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 在生产环境中，需要添加认证
        // 只允许管理员访问
        c.Next(ctx)
    }
}
```

**使用 pprof 分析性能**:
```bash
# 1. 启动服务时开启 pprof
go run main.go --enable-pprof

# 2. 采集 CPU profile
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof

# 3. 分析 CPU profile
go tool pprof cpu.prof

# 4. 采集内存 profile
curl http://localhost:6060/debug/pprof/heap > heap.prof

# 5. 分析内存 profile
go tool pprof heap.prof

# 6. 采集 goroutine
curl http://localhost:6060/debug/pprof/goroutine > goroutine.prof
```

##### Day 4-5: Prometheus 监控集成

**Prometheus 指标定义** (`infra/monitoring/metrics/metrics.go`):
```go
// metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP 请求总数
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    // HTTP 请求延迟
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request latency in seconds",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
        },
        []string{"method", "endpoint"},
    )

    // 配额检查延迟
    QuotaCheckDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "quota_check_duration_seconds",
            Help:    "Quota check latency in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
        },
    )

    // 配额使用量
    QuotaUsage = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "quota_usage",
            Help: "Current quota usage",
        },
        []string{"tenant_id", "resource_type"},
    )

    // 权限检查延迟
    PermissionCheckDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "permission_check_duration_seconds",
            Help:    "Permission check latency in seconds",
            Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
        },
    )

    // 路由决策延迟
    RoutingDecisionDuration = promauto.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "routing_decision_duration_seconds",
            Help:    "Routing decision latency in seconds",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1},
        },
    )

    // 数据库连接池使用率
    DBPoolUsage = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "db_pool_usage",
            Help: "Database connection pool usage",
        },
        []string{"database"},
    )

    // 缓存命中率
    CacheHitRate = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cache_operations_total",
            Help: "Total cache operations",
        },
        []string{"operation", "cache_type"},
    )
)

// RecordHTTPRequest 记录 HTTP 请求
func RecordHTTPRequest(method, endpoint string, status int, duration float64) {
    HTTPRequestsTotal.WithLabelValues(method, endpoint, fmt.Sprintf("%d", status)).Inc()
    HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

// RecordQuotaCheck 记录配额检查
func RecordQuotaCheck(duration float64) {
    QuotaCheckDuration.Observe(duration)
}

// RecordPermissionCheck 记录权限检查
func RecordPermissionCheck(duration float64) {
    PermissionCheckDuration.Observe(duration)
}
```

**Prometheus 中间件** (`api/middleware/prometheus_middleware.go`):
```go
// middleware/prometheus_middleware.go
package middleware

import (
    "context"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-studio/infra/monitoring/metrics"
)

func PrometheusMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        start := time.Now()

        // 执行请求
        c.Next(ctx)

        // 记录指标
        duration := time.Since(start).Seconds()
        method := string(c.Request.Method())
        endpoint := string(c.Request.URI().Path())
        status := c.Response.StatusCode()

        metrics.RecordHTTPRequest(method, endpoint, status, duration)
    }
}
```

**Grafana 监控大盘配置** (`infra/monitoring/grafana/dashboards/api-dashboard.json`):
```json
{
  "dashboard": {
    "title": "ZKER API Performance Dashboard",
    "panels": [
      {
        "title": "Request Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "P95 Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Error Rate",
        "targets": [
          {
            "expr": "rate(http_requests_total{status=~\"5..\"}[5m]) / rate(http_requests_total[5m])"
          }
        ]
      },
      {
        "title": "Quota Check Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(quota_check_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Permission Check Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(permission_check_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Routing Decision Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(routing_decision_duration_seconds_bucket[5m]))"
          }
        ]
      }
    ]
  }
}
```

---

### Week 5-6: 监控和日志系统

#### Week 5: 结构化日志

**目标**: 建立统一的日志规范和收集系统

##### Day 1-3: 结构化日志库

**日志库封装** (`infra/logging/logger.go`):
```go
// logging/logger.go
package logging

import (
    "context"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var (
    logger *zap.Logger
)

// Init 初始化日志
func Init(env string) error {
    var config zap.Config

    if env == "production" {
        // 生产环境：JSON 格式，输出到文件
        config = zap.Config{
            Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
            Development: false,
            Sampling: &zap.SamplingConfig{
                Initial:    100,
                Thereafter: 100,
            },
            Encoding: "json",
            EncoderConfig: zapcore.EncoderConfig{
                TimeKey:        "timestamp",
                LevelKey:       "level",
                NameKey:        "logger",
                CallerKey:      "caller",
                MessageKey:     "message",
                StacktraceKey:  "stacktrace",
                LineEnding:     zapcore.DefaultLineEnding,
                EncodeLevel:    zapcore.LowercaseLevelEncoder,
                EncodeTime:     zapcore.ISO8601TimeEncoder,
                EncodeDuration: zapcore.SecondsDurationEncoder,
                EncodeCaller:   zapcore.ShortCallerEncoder,
            },
            OutputPaths:      []string{"/var/log/zker/app.log"},
            ErrorOutputPaths: []string{"/var/log/zker/error.log"},
        }
    } else {
        // 开发环境：Console 格式，彩色输出
        config = zap.Config{
            Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
            Development: true,
            Encoding:    "console",
            EncoderConfig: zapcore.EncoderConfig{
                TimeKey:        "T",
                LevelKey:       "L",
                NameKey:        "N",
                CallerKey:      "C",
                FunctionKey:    zapcore.OmitKey,
                MessageKey:     "M",
                StacktraceKey:  "S",
                LineEnding:     zapcore.DefaultLineEnding,
                EncodeLevel:    zapcore.CapitalColorLevelEncoder,
                EncodeTime:     zapcore.ISO8601TimeEncoder,
                EncodeDuration: zapcore.StringDurationEncoder,
                EncodeCaller:   zapcore.ShortCallerEncoder,
            },
            OutputPaths:      []string{"stdout"},
            ErrorOutputPaths: []string{"stderr"},
        }
    }

    var err error
    logger, err = config.Build()
    return err
}

// Info 记录 Info 日志
func Info(msg string, fields ...zap.Field) {
    logger.Info(msg, fields...)
}

// Error 记录 Error 日志
func Error(msg string, fields ...zap.Field) {
    logger.Error(msg, fields...)
}

// Debug 记录 Debug 日志
func Debug(msg string, fields ...zap.Field) {
    logger.Debug(msg, fields...)
}

// Warn 记录 Warn 日志
func Warn(msg string, fields ...zap.Field) {
    logger.Warn(msg, fields...)
}

// Fatal 记录 Fatal 日志并退出
func Fatal(msg string, fields ...zap.Field) {
    logger.Fatal(msg, fields...)
}

// WithContext 从 context 中提取信息并创建 logger
func WithContext(ctx context.Context) *zap.Logger {
    fields := make([]zap.Field, 0)

    // 从 context 中提取 request_id
    if requestID := ctx.Value("request_id"); requestID != nil {
        fields = append(fields, zap.String("request_id", requestID.(string)))
    }

    // 从 context 中提取 user_id
    if userID := ctx.Value("user_id"); userID != nil {
        fields = append(fields, zap.String("user_id", userID.(string)))
    }

    // 从 context 中提取 tenant_id
    if tenantID := ctx.Value("tenant_id"); tenantID != nil {
        fields = append(fields, zap.String("tenant_id", tenantID.(string)))
    }

    return logger.With(fields...)
}
```

**日志中间件** (`api/middleware/logging_middleware.go`):
```go
// middleware/logging_middleware.go
package middleware

import (
    "context"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-studio/infra/logging"
    "go.uber.org/zap"
)

func LoggingMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        start := time.Now()

        // 记录请求开始
        logging.WithContext(ctx).Info("Request started",
            zap.String("method", string(c.Request.Method())),
            zap.String("path", string(c.Request.URI().Path())),
            zap.String("query", string(c.Request.URI().QueryString())),
            zap.String("client_ip", c.ClientIP()),
        )

        // 执行请求
        c.Next(ctx)

        // 记录请求完成
        duration := time.Since(start)
        logging.WithContext(ctx).Info("Request completed",
            zap.Int("status", c.Response.StatusCode()),
            zap.Duration("duration", duration),
            zap.Int("response_size", len(c.Response.Body())),
        )
    }
}
```

##### Day 4-5: 日志收集和查询

**Filebeat 配置** (`infra/logging/filebeat.yml`):
```yaml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/zker/*.log
    json.keys_under_root: true
    json.add_error_key: true

    fields:
      app: zker
      env: ${ENV:development}

    fields_under_root: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  indices:
    - index: "zker-logs-%{+yyyy.MM.dd}"
      when.equals:
        app: "zker"

setup.template.name: "zker-logs"
setup.template.pattern: "zeker-logs-*"

processors:
  - drop_event:
      when:
        equals:
          level: "DEBUG"
```

**日志查询工具** (`scripts/query_logs.sh`):
```bash
#!/bin/bash
# 查询日志

INDEX=$1
QUERY=$2

curl -X GET "localhost:9200/${INDEX}/_search?pretty" \
  -H 'Content-Type: application/json' \
  -d'
{
  "query": {
    "bool": {
      "must": [
        { "match": { "message": "'"$QUERY"'" }}
      ]
    }
  },
  "size": 100,
  "sort": [
    { "@timestamp": "desc" }
  ]
}
'
```

#### Week 6: 分布式追踪

**目标**: 集成 OpenTelemetry，实现端到端追踪

##### Day 1-3: OpenTelemetry 集成

**Tracer 初始化** (`infra/tracing/tracer.go`):
```go
// tracing/tracer.go
package tracing

import (
    "context"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Init 初始化 tracer
func Init(serviceName, jaegerEndpoint string) error {
    // 创建 Jaeger exporter
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
        jaeger.WithEndpoint(jaegerEndpoint),
    ))
    if err != nil {
        return err
    }

    // 创建 resource
    res, err := resource.Merge(
        resource.Default(),
        resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
        ),
    )
    if err != nil {
        return err
    }

    // 创建 tracer provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(res),
        trace.WithSampler(trace.TraceIDRatioBased(0.1)), // 采样率 10%
    )

    otel.SetTracerProvider(tp)
    return nil
}

// StartSpan 启动 span
func StartSpan(ctx context.Context, tracerName, spanName string) (context.Context, *trace.Span) {
    tracer := otel.Tracer(tracerName)
    return tracer.Start(ctx, spanName)
}
```

**追踪中间件** (`api/middleware/tracing_middleware.go`):
```go
// middleware/tracing_middleware.go
package middleware

import (
    "context"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-studio/infra/tracing"
    "go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(serviceName string) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 从请求头中提取 trace context
        carrier := make propagation.HeaderCarrier)
        carrier.Set("traceparent", c.Request.Header.Get("traceparent"))

        // 提取或创建 span context
        spanContext := trace.SpanContext{}
        if spanContext.IsValid() {
            ctx = trace.ContextWithSpanContext(ctx, spanContext)
        }

        // 创建 span
        ctx, span := tracing.StartSpan(
            ctx,
            serviceName,
            string(c.Request.Method())+" "+string(c.Request.URI().Path()),
        )

        // 添加 span 属性
        span.SetAttributes(
            attribute.String("http.method", string(c.Request.Method())),
            attribute.String("http.url", string(c.Request.URI().String())),
            attribute.String("http.host", c.Request.Host()),
            attribute.String("http.remote_addr", c.ClientIP()),
        )

        // 执行请求
        c.Next(ctx)

        // 结束 span
        span.SetAttributes(
            attribute.Int("http.status_code", c.Response.StatusCode()),
        )
        span.End()
    }
}
```

**数据库追踪** (`infra/db/tracing.go`):
```go
// db/tracing.go
package db

import (
    "context"

    "gorm.io/gorm"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
    "go.opentelemetry.io/otel/trace"
)

// TracedPlugin GORM 追踪插件
type TracedPlugin struct{}

func (p *TracedPlugin) Name() string {
    return "traced"
}

func (p *TracedPlugin) Initialize(db *gorm.DB) error {
    // 注册回调
    db.Callback().Query().Before("gorm:query").Register("traced:before_query", p.beforeQuery)
    db.Callback().Query().After("gorm:query").Register("traced:after_query", p.afterQuery)
    db.Callback().Create().Before("gorm:create").Register("traced:before_create", p.beforeCreate)
    db.Callback().Create().After("gorm:create").Register("traced:after_create", p.afterCreate)
    return nil
}

func (p *TracedPlugin) beforeQuery(db *gorm.DB) {
    tracer := otel.Tracer("gorm")
    ctx := db.Statement.Context
    _, span := tracer.Start(ctx, "gorm:query")

    span.SetAttributes(
        attribute.String("sql.table", db.Statement.Table),
        attribute.String("sql.query", db.Statement.SQL.String()),
    )
}

func (p *TracedPlugin) afterQuery(db *gorm.DB) {
    // 结束 span
}
```

##### Day 4-5: 追踪数据查询和可视化

**查询追踪数据** (`scripts/query_traces.sh`):
```bash
#!/bin/bash
# 查询 Jaeger 追踪数据

TRACE_ID=$1

curl "http://localhost:16686/api/traces/${TRACE_ID}"
```

**Grafana Trace Dashboard**:
```json
{
  "dashboard": {
    "title": "ZKER Distributed Tracing",
    "panels": [
      {
        "title": "Request Latency by Endpoint",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Trace Search",
        "type": "table",
        "targets": [
          {
            "query": "service.name=\"zker-api\""
          }
        ]
      }
    ]
  }
}
```

---

### Week 7-8: 高级监控 + 优化

#### Week 7: 告警系统和日志分析

**目标**: 建立告警规则，实现自动化日志分析

##### Day 1-3: Prometheus 告警规则

**告警规则配置** (`infra/monitoring/prometheus/alerts.yml`):
```yaml
groups:
  - name: api_alerts
    interval: 30s
    rules:
      # API 错误率告警
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value | humanizePercentage }} for {{ $labels.endpoint }}"

      # API 延迟告警
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High latency detected"
          description: "P95 latency is {{ $value }}s for {{ $labels.endpoint }}"

      # 配额即将耗尽告警
      - alert: QuotaNearLimit
        expr: quota_usage / quota_limit > 0.8
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Quota near limit"
          description: "Quota usage is {{ $value | humanizePercentage }} for {{ $labels.tenant_id }}"

      # 数据库连接池告警
      - alert: DBPoolHighUsage
        expr: db_pool_usage / db_pool_max > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Database pool high usage"
          description: "Pool usage is {{ $value | humanizePercentage }}"

      # 缓存命中率低告警
      - alert: LowCacheHitRate
        expr: rate(cache_operations_total{operation="hit"}[5m]) / rate(cache_operations_total[5m]) < 0.7
        for: 10m
        labels:
          severity: info
        annotations:
          summary: "Low cache hit rate"
          description: "Cache hit rate is {{ $value | humanizePercentage }}"
```

**告警通知** (`infra/monitoring/alertmanager/alertmanager.yml`):
```yaml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'

  routes:
    - match:
        severity: critical
      receiver: 'critical'

receivers:
  - name: 'default'
    webhook_configs:
      - url: 'http://localhost:5001/alerts'

  - name: 'critical'
    webhook_configs:
      - url: 'http://localhost:5001/alerts'
    email_configs:
      - to: 'oncall@zker.com'
        from: 'alertmanager@zker.com'
        smarthost: 'smtp.gmail.com:587'
```

##### Day 4-5: 日志分析和异常检测

**日志异常检测** (`scripts/detect_anomalies.sh`):
```bash
#!/bin/bash
# 检测日志中的异常模式

# 1. 检测错误激增
ERROR_COUNT=$(curl -s "localhost:9200/zker-logs-*/_search" \
  -H 'Content-Type: application/json' \
  -d'
  {
    "query": {
      "range": {
        "@timestamp": {
          "gte": "now-5m"
        }
      }
    },
    "size": 0,
    "aggs": {
      "error_count": {
        "filter": {
          "term": { "level": "error" }
        }
      }
    }
  }
' | jq '.aggregations.error_count.doc_count')

# 如果错误数超过阈值，发送告警
if [ $ERROR_COUNT -gt 100 ]; then
  echo "High error count detected: $ERROR_COUNT"
  # 发送告警...
fi

# 2. 检测异常关键词
ANOMALY_PATTERNS=("panic" "fatal" "OutOfMemory" "Deadlock")

for pattern in "${ANOMALY_PATTERNS[@]}"; do
  COUNT=$(curl -s "localhost:9200/zker-logs-*/_count" \
    -H 'Content-Type: application/json' \
    -d'
    {
      "query": {
        "bool": {
          "must": [
            { "range": { "@timestamp": { "gte": "now-1h" } } },
            { "match": { "message": "'"$pattern"'" } }
          ]
        }
      }
    }
  ' | jq '.count')

  if [ $COUNT -gt 0 ]; then
    echo "Anomaly detected: $pattern ($COUNT occurrences)"
  fi
done
```

#### Week 8: 性能优化总结 + 文档

**目标**: 优化总结，编写运维手册

##### Day 1-3: 性能优化总结

**性能优化报告** (`docs/performance/optimization_report.md`):
```markdown
# 性能优化报告

## 优化前

- API P95 延迟: 800ms
- 配额检查延迟: 150ms
- 权限检查延迟: 200ms
- 数据库连接池使用率: 90%

## 优化后

- API P95 延迟: 450ms (↓ 43%)
- 配额检查延迟: 40ms (↓ 73%)
- 权限检查延迟: 80ms (↓ 60%)
- 数据库连接池使用率: 65% (↓ 28%)

## 优化措施

### 1. 数据库优化
- 添加索引: `idx_tenant_id`, `idx_status_created`
- 使用 Preload 避免 N+1 查询
- 使用连接池（最大连接数 100）

### 2. 缓存优化
- Redis 缓存配额数据（TTL 5 分钟）
- Redis 缓存权限数据（TTL 10 分钟）
- 缓存命中率: 85%

### 3. 并发优化
- 使用 goroutine 池处理并发请求
- 信号量限制并发数（最大 100）
- 减少锁竞争

### 4. 代码优化
- 减少不必要的日志输出
- 优化 JSON 序列化
- 使用对象池复用对象

## 监控指标

- 错误率: < 1%
- 可用性: > 99.9%
- P95 延迟: < 500ms
- P99 延迟: < 1s
```

##### Day 4-5: 运维手册

**运维手册** (`docs/operations/monitoring_guide.md`):
```markdown
# 监控和告警运维手册

## 监控架构

- Prometheus: 指标采集和存储
- Grafana: 可视化大盘
- Alertmanager: 告警路由
- Jaeger: 分布式追踪
- Elasticsearch + Filebeat: 日志收集

## 关键指标

### API 指标
- QPS: `rate(http_requests_total[1m])`
- P95 延迟: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))`
- 错误率: `rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])`

### 业务指标
- 配额检查延迟: `histogram_quantile(0.95, rate(quota_check_duration_seconds_bucket[5m]))`
- 权限检查延迟: `histogram_quantile(0.95, rate(permission_check_duration_seconds_bucket[5m]))`
- 路由决策延迟: `histogram_quantile(0.95, rate(routing_decision_duration_seconds_bucket[5m]))`

### 基础设施指标
- CPU 使用率: `rate(process_cpu_seconds_total[5m])`
- 内存使用量: `process_resident_memory_bytes`
- 数据库连接池: `db_pool_usage / db_pool_max`

## 告警规则

### Critical 告警
- API 错误率 > 5%
- P95 延迟 > 1s
- 数据库连接池使用率 > 90%

### Warning 告警
- API 错误率 > 1%
- P95 延迟 > 500ms
- 配额使用率 > 80%

## 故障排查

### 1. API 延迟高

**排查步骤**:
1. 查看 Grafana 延迟图表，定位慢接口
2. 使用 Jaeger 查看调用链，找到慢操作
3. 检查数据库慢查询日志
4. 使用 pprof 分析 CPU profile

**常见原因**:
- 数据库慢查询
- 缓存未命中
- 外部服务调用慢
- 锁竞争

### 2. 错误率突增

**排查步骤**:
1. 查看错误日志，找到错误类型
2. 检查依赖服务状态
3. 检查数据库连接
4. 检查内存和磁盘

**常见原因**:
- 数据库连接耗尽
- 依赖服务不可用
- 配额/权限拒绝
- 代码 Bug

### 3. 内存泄漏

**排查步骤**:
1. 使用 pprof 采集 heap profile
2. 分析内存占用
3. 检查 goroutine 数量
4. 检查缓存大小

## 日常运维

### 每日
- 检查 Grafana 大盘
- 处理 Critical 告警
- 查看错误日志

### 每周
- 分析性能趋势
- 检查告警规则是否需要调整
- 审查日志存储空间

### 每月
- 性能测试
- 容量规划
- 优化建议
```

---

## 📦 专属开发规范

### 作为质量保障工程师的特殊要求

#### 1. 可观测性优先

✅ **必须**：
- 所有 API 都有 Prometheus 指标
- 所有操作都有日志记录
- 关键路径都有分布式追踪

❌ **禁止**：
- "黑盒"代码无法追踪
- 日志级别错误（如生产环境 DEBUG）
- 缺少关键指标

#### 2. 错误处理规范

✅ **必须**：
- 所有错误使用统一错误码
- 错误日志包含上下文
- 敏感信息不暴露到错误消息

```go
// ✅ Good
if err != nil {
    logging.WithContext(ctx).Error("Failed to create bot",
        zap.String("tenant_id", tenantID),
        zap.Error(err),
    )
    return error_codes.ErrBotCreationFailed.WithDetails(map[string]interface{}{
        "tenant_id": tenantID,
    })
}

// ❌ Bad
if err != nil {
    log.Println(err) // 缺少上下文
    return err // 未使用统一错误码
}
```

#### 3. 性能测试规范

✅ **必须**：
- 性能测试用例与单元测试分离
- 使用真实数据量进行测试
- 记录性能基线

❌ **禁止**：
- 在生产环境运行性能测试
- 性能测试数据污染生产数据库

---

## 🤝 协作接口定义

### 与研发 A 协作

**向研发 A 提供**:
- 统一错误码定义（`types/errno/`）
- 性能测试脚本（`tests/performance/`）
- 监控指标定义（`infra/monitoring/metrics/`）

**依赖研发 A**:
- 数据库 Schema（用于优化查询）
- API 接口定义（用于性能测试）

### 与研发 C 协作

**向研发 C 提供**:
- 错误码 JSON 文件（前端映射）
- API 文档（OpenAPI 规范）
- 监控数据展示接口

**依赖研发 C**:
- 前端性能指标收集
- 前端错误日志上报

### 与研发 D 协作

**向研发 D 提供**:
- Prometheus 配置
- Grafana Dashboard 配置
- 日志收集配置
- 告警规则配置

**依赖研发 D**:
- 监控服务部署
- 日志服务部署
- 告警通知渠道配置

---

## 📅 里程碑和交付物

### Week 2 交付物

- [ ] 统一错误码定义（300+ 错误码）
- [ ] 错误码中间件
- [ ] API 契约测试
- [ ] OpenAPI 规范文件

### Week 4 交付物

- [ ] K6 性能测试脚本（5+ 场景）
- [ ] 性能基准测试
- [ ] Prometheus 指标定义
- [ ] Grafana 监控大盘

### Week 6 交付物

- [ ] 结构化日志库
- [ ] 日志收集配置
- [ ] OpenTelemetry 追踪
- [ ] 追踪可视化

### Week 8 交付物

- [ ] 告警规则配置
- [ ] 性能优化报告
- [ ] 运维手册
- [ ] 故障排查指南

---

## ✅ 质量检查清单

### 代码提交前

- [ ] 错误码符合规范
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 性能测试通过
- [ ] 日志级别正确
- [ ] Prometheus 指标定义

### 集成前

- [ ] 与其他模块接口测试通过
- [ ] 监控大盘正常展示
- [ ] 告警规则生效
- [ ] 日志正常收集

### 发布前

- [ ] 性能基线达标
- [ ] 监控告警配置完成
- [ ] 日志采集正常
- [ ] 运维手册完整

---

## 📊 关键指标

### 代码质量

- 单元测试覆盖率：≥ 80%
- 性能测试覆盖率：100%（核心接口）
- 错误码覆盖率：100%（所有 API）

### 性能指标

- API P95 延迟：< 500ms
- 配额检查延迟：< 50ms
- 权限检查延迟：< 100ms
- 路由决策延迟：< 200ms

### 监控指标

- 监控覆盖率：100%（所有服务）
- 日志采集率：100%
- 告警准确率：> 95%

---

## 📚 参考资料

- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-全局一致性检查清单_v1.0.md](../ZKER-全局一致性检查清单_v1.0.md)
- [ZKER-实现差距分析与研发计划_v1.0.md](../ZKER-实现差距分析与研发计划_v1.0.md)
- [ZKER-统一错误码定义规范.md](../ZKER-统一错误码定义规范.md)
- [ZKER-故障排查手册_v1.0.md](../ZKER-故障排查手册_v1.0.md)

---

**文档状态**: ✅ 已完成
**最后更新**: 2025-01-01
**责任人**: 研发 B - 后端工程师
**评审人**: 技术负责人
