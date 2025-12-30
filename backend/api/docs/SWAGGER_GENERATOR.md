# API 文档生成器使用指南

## 概述

这是一个企业级 OpenAPI 3.0 规范文档生成器，用于自动生成 RESTful API 的 Swagger 文档。

## 功能特性

- ✅ **完整的 OpenAPI 3.0 规范支持**
- ✅ **自动扫描 Handler 文件**，从注释中提取 API 定义
- ✅ **同时生成 YAML 和 JSON 格式**
- ✅ **企业级错误码规范**
- ✅ **多环境配置支持**
- ✅ **内置 Schema 定义**
- ✅ **JWT 和 API Key 认证方案**
- ✅ **自动验证生成的规范**

## 快速开始

### 1. 基本使用

```bash
# 使用默认配置（从 backend/api/docs 目录运行）
cd backend/api/docs
go run generate_swagger.go

# 输出：
# 📂 Scanning handler files in: D:\code\coze-studio\backend\api\handler\coze\org
# ✅ Found 27 API paths
# ✅ OpenAPI spec generated successfully: swagger\swagger.yaml
# ✅ Also generated: swagger\swagger.json
```

### 2. 使用自定义配置

创建配置文件 `config.json`：

```json
{
  "input_dir": "D:\\code\\coze-studio\\backend\\api\\handler\\coze\\org",
  "output_dir": "./swagger",
  "format": "yaml",
  "version": "v1.0.0",
  "base_path": "/api"
}
```

然后运行：

```bash
go run generate_swagger.go config.json
```

### 3. 编译可执行文件

```bash
# 编译
go build -o generate_swagger.exe generate_swagger.go

# 运行
./generate_swagger.exe
```

## 配置说明

| 配置项 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| `input_dir` | string | Handler 文件目录 | `../handler/coze/org` |
| `output_dir` | string | 输出目录 | `./swagger` |
| `format` | string | 主输出格式（yaml/json） | `yaml` |
| `version` | string | API 版本 | `v1.0.0` |
| `base_path` | string | API 基础路径 | `/api` |

## Handler 注释规范

为了自动生成 API 文档，需要在 Handler 文件中添加特定的注释：

### 基本格式

```go
// @router GET /api/organizations
// @summary 获取组织列表
// @description 查询所有组织，支持分页和筛选
// @tags 组织管理
// @param query page int false "页码"
// @param query page_size int false "每页数量"
// @success 200 {object} OrganizationListResponse
// @failure 400 {object} ErrorResponse
// @failure 401 {object} ErrorResponse
func (h *OrganizationHandler) GetOrganizations(ctx context.Context, req *GetOrganizationsRequest) (*GetOrganizationsResponse, error) {
    // 实现...
}
```

### 支持的注释标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `@router` | 定义 HTTP 方法和路径 | `@router GET /api/organizations` |
| `@summary` | API 简短描述 | `@summary 获取组织列表` |
| `@description` | 详细说明 | `@description 查询所有组织...` |
| `@tags` | API 分组标签 | `@tags 组织管理` |
| `@param` | 参数定义 | `@param query page int false "页码"` |
| `@success` | 成功响应 | `@success 200 {object} Response` |
| `@failure` | 失败响应 | `@failure 400 {object} ErrorResponse` |

### 参数定义格式

```
@param {位置} {名称} {类型} {必需} {描述}
```

- **位置**: `query`、`path`、`header`、`body`
- **类型**: `string`、`int`、`bool`、`object` 等
- **必需**: `true`、`false`

示例：

```go
// @param query page int false "页码，默认1"
// @param query page_size int false "每页数量，默认20"
// @param path id string true "组织ID"
// @param body request body CreateOrganizationRequest true "创建请求"
```

## 输出文件

生成器会在输出目录中创建以下文件：

```
swagger/
├── swagger.yaml    # OpenAPI 3.0 YAML 格式
└── swagger.json    # OpenAPI 3.0 JSON 格式
```

## 查看文档

### 1. Swagger Editor（在线）

访问 [Swagger Editor](https://editor.swagger.io/)，导入生成的 `swagger.yaml` 或 `swagger.json` 文件。

### 2. Redoc（在线）

访问 [Redoc](https://redocly.github.io/redoc/)，导入生成的文件。

### 3. 本地 Swagger UI

使用 Docker 运行 Swagger UI：

```bash
docker run -p 8080:8080 \
  -e SWAGGER_JSON=/swagger/swagger.yaml \
  -v $(pwd)/swagger:/swagger \
  swaggerapi/swagger-ui
```

然后访问 `http://localhost:8080`。

### 4. VS Code 插件

安装 [OpenAPI (Swagger) Editor](https://marketplace.visualstudio.com/items?items=42Crash.vscode-openapi) 插件，直接在 VS Code 中预览。

## 生成的文档结构

### Info 元数据

```yaml
info:
  title: 组织中心管理API
  version: 1.0.0
  description: 企业级组织管理系统的RESTful API文档
  contact:
    name: API Support
    email: api-support@coze.com
  license:
    name: Apache 2.0
    url: https://www.apache.org/licenses/LICENSE-2.0.html
```

### 服务器配置

```yaml
servers:
  - url: http://localhost:8080
    description: 开发环境
  - url: https://api-test.coze.com
    description: 测试环境
  - url: https://api.coze.com
    description: 生产环境
```

### 认证方案

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT Token认证
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
      description: API Key认证
```

### Schema 定义

生成的文档包含以下预定义 Schema：

- **ErrorResponse**: 统一错误响应格式
  - `code`: 错误码
  - `message`: 错误信息（中文）
  - `message_en`: 错误信息（英文）
  - `data`: 额外错误详情

- **OrganizationData**: 组织数据结构
  - `org_id`: 组织ID
  - `tenant_id`: 租户ID
  - `org_name`: 组织名称
  - `org_type`: 组织类型
  - 等等...

## 错误码规范

所有 API 错误响应遵循统一格式：

```json
{
  "code": 100001,
  "message": "租户不存在",
  "message_en": "Tenant not found",
  "data": {
    "tenant_id": "123456"
  }
}
```

错误码规范详见：[统一错误码定义规范](../../企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)

## 最佳实践

### 1. API 设计原则

- ✅ 使用 RESTful 风格的 URL 设计
- ✅ 使用 HTTP 方法表达操作意图（GET/POST/PUT/DELETE）
- ✅ 资源命名使用名词复数形式
- ✅ 使用标准 HTTP 状态码
- ✅ 统一的错误响应格式

### 2. 注释规范

- ✅ 每个公开 API 都必须有完整注释
- ✅ `@summary` 应简洁明了，不超过 50 字
- ✅ `@description` 详细说明功能、参数、注意事项
- ✅ 使用 `@tags` 进行合理的 API 分组
- ✅ 所有参数必须标注类型和是否必需

### 3. 版本管理

- 在 API 路径中包含版本号：`/api/v1/organizations`
- 生成文档时更新 `version` 字段
- 重大变更时递增主版本号

### 4. 安全性

- 所有业务 API 都需要认证
- 敏感操作需要权限检查
- 使用 HTTPS 传输
- 敏感信息不暴露在 URL 中

## 故障排查

### 问题：编译错误 "yaml redeclared"

**原因**：同时导入了多个 yaml 包

**解决**：只保留 `gopkg.in/yaml.v3`

```go
import (
    "gopkg.in/yaml.v3"  // ✅ Good
)
```

### 问题：路径没有生成

**原因**：Handler 文件中缺少 `@router` 注释

**解决**：添加完整的注释

```go
// @router GET /api/organizations
```

### 问题：生成的路径为空 `{}`

**原因**：注释格式不正确或解析失败

**解决**：
1. 检查 `@router` 格式是否正确
2. 确保方法（GET/POST/PUT/DELETE）大写
3. 路径使用引号包裹

## 集成到 CI/CD

### GitHub Actions 示例

```yaml
name: Generate API Docs

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/api/handler/**'

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'

      - name: Generate Swagger
        run: |
          cd backend/api/docs
          go run generate_swagger.go

      - name: Commit docs
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add backend/api/docs/swagger/
          git commit -m "docs: auto-generate API documentation"
          git push
```

## 进阶功能

### 1. 自定义 Schema

在 `generateSchemas()` 函数中添加自定义 Schema：

```go
func generateSchemas() map[string]Schema {
    return map[string]Schema{
        "MyCustomType": {
            Type: "object",
            Required: []string{"id", "name"},
            Properties: map[string]Schema{
                "id": {
                    Type:        "string",
                    Description: "ID",
                },
                "name": {
                    Type:        "string",
                    Description: "名称",
                },
            },
        },
    }
}
```

### 2. 添加通用响应

在 `generateCommonResponses()` 中添加：

```go
"Created": {
    Description: "资源创建成功",
    Content: map[string]MediaType{
        "application/json": {
            Schema: &Schema{Ref: "#/components/schemas/Resource"},
        },
    },
},
```

### 3. 多服务器配置

修改 `generateOpenAPISpec()` 中的 Servers 配置：

```go
Servers: []Server{
    {
        URL:         "http://{env}.example.com",
        Description: "动态环境",
        Variables: map[string]Variable{
            "env": {
                Default:     "dev",
                Description: "环境",
                Enum:        []string{"dev", "test", "prod"},
            },
        },
    },
},
```

## 参考资源

- [OpenAPI 3.0 规范](https://swagger.io/specification/)
- [Swagger 官方文档](https://swagger.io/docs/)
- [Redoc 文档](https://github.com/Redocly/redoc)
- [企业级开发规范手册](../../企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](../../企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)

## 贡献指南

如需改进文档生成器，请：

1. Fork 本项目
2. 创建特性分支
3. 提交变更
4. 推送到分支
5. 创建 Pull Request

## 许可证

Apache License 2.0

---

**生成器版本**: v1.0.0
**最后更新**: 2025-01-01
**维护者**: API Team
