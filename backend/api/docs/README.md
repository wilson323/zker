# 组织中心API文档

欢迎使用组织中心管理API文档！本目录包含完整的API文档、示例和工具，帮助您快速集成和使用组织中心API。

## 📁 文档结构

```
backend/api/docs/
├── swagger/                    # Swagger/OpenAPI文档
│   ├── swagger.yaml            # OpenAPI 3.0规范文件
│   ├── index.html              # Swagger UI交互式文档
│   └── README.md               # Swagger文档说明
├── postman/                    # Postman Collection
│   └── OrgCenter_Collection.json  # 完整的Postman集合
├── API_GUIDE.md                # API使用指南（必读）
├── CLIENT_SDK.md               # 客户端SDK文档
├── generate_swagger.go         # Swagger文档生成工具
└── README.md                   # 本文件
```

---

## 🚀 快速开始

### 1. 在线查看API文档

打开Swagger UI文档：

```bash
# 方式一：直接在浏览器中打开（需要先启动HTTP服务器）
open backend/api/docs/swagger/index.html

# 方式二：使用Python启动本地服务器
cd backend/api/docs/swagger
python3 -m http.server 8080
# 然后访问 http://localhost:8080
```

Swagger UI提供：
- ✅ 完整的API端点列表
- ✅ 在线测试API（Try it out）
- ✅ 请求/响应示例
- ✅ Schema定义查看
- ✅ 环境切换功能

### 2. 导入Postman Collection

1. 打开Postman
2. 点击 `Import` 按钮
3. 选择 `backend/api/docs/postman/OrgCenter_Collection.json`
4. 开始测试API！

Postman Collection包含：
- ✅ 63个API请求配置
- ✅ 预配置的环境变量
- ✅ 自动化测试脚本
- ✅ 完整的请求示例

### 3. 阅读API使用指南

```bash
# 查看API使用指南
cat backend/api/docs/API_GUIDE.md
```

API使用指南包含：
- ✅ API概览和认证说明
- ✅ 通用请求/响应格式
- ✅ 完整的错误码列表
- ✅ 63个API端点详细说明
- ✅ cURL/JavaScript/Python请求示例
- ✅ 最佳实践和FAQ

### 4. 使用客户端SDK

查看SDK文档，选择适合您的语言：

```bash
# 查看SDK文档
cat backend/api/docs/CLIENT_SDK.md
```

支持的SDK：
- ✅ Go SDK
- ✅ JavaScript/TypeScript SDK
- ✅ Python SDK
- ✅ Java SDK

---

## 📚 核心文档

### API使用指南 (API_GUIDE.md)

**必读文档！** 包含完整的API使用说明：

| 章节 | 内容 |
|------|------|
| API概览 | 基础信息、服务端点、核心功能 |
| 认证授权 | Token获取、租户隔离、权限控制 |
| 通用格式 | 请求格式、响应格式、分页参数 |
| 错误码说明 | HTTP状态码、业务错误码 |
| API端点 | 63个API详细说明 |
| 请求示例 | cURL、JavaScript、Python示例 |
| 最佳实践 | 错误处理、Token管理、分页查询等 |
| FAQ | 常见问题解答 |

### 客户端SDK文档 (CLIENT_SDK.md)

详细的SDK使用文档：

| 语言 | 特性 |
|------|------|
| Go | • 完整的类型定义<br>• 自动重试机制<br>• 并发安全 |
| JavaScript/TypeScript | • React Hook支持<br>• TypeScript类型安全<br>• Axios/Fetch兼容 |
| Python | • 同步/异步API<br>• 上下文管理器<br>• 自动Token刷新 |
| Java | • Builder模式<br>• RxJava支持<br>• Spring Boot集成 |

---

## 🛠️ 开发工具

### Swagger文档生成器

企业级 OpenAPI 3.0 文档生成器，支持自动扫描 Handler 文件并生成标准 API 文档。

#### 快速开始

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

#### 使用 Makefile（推荐）

```bash
# 查看所有可用命令
make help

# 常用命令
make run          # 生成文档
make build        # 编译生成器
make serve        # 启动本地 Swagger UI
make clean        # 清理生成的文件
make validate     # 验证 OpenAPI 规范
```

#### 配置文件

创建 `config.json` 自定义配置：

```json
{
  "input_dir": "D:\\code\\coze-studio\\backend\\api\\handler\\coze\\org",
  "output_dir": "./swagger",
  "format": "yaml",
  "version": "v1.0.0",
  "base_path": "/api"
}
```

#### 输出文件

生成器会同时创建 YAML 和 JSON 两种格式：

```
swagger/
├── swagger.yaml    # OpenAPI 3.0 YAML 格式
└── swagger.json    # OpenAPI 3.0 JSON 格式
```

#### 在线查看文档

1. **Swagger Editor**: https://editor.swagger.io/
   - 导入生成的 `swagger.yaml` 或 `swagger.json`
   - 可视化编辑和测试

2. **Redoc**: https://redocly.github.io/redoc/
   - 更美观的文档展示
   - 支持三栏布局

3. **本地 Swagger UI**:
   ```bash
   make serve
   # 访问 http://localhost:8080
   ```

### Swagger注释规范

在Handler中添加Swagger注释以自动生成文档：

```go
// CreateOrganization 创建组织
// @router POST /api/organizations
// @summary 创建组织
// @description 创建一个新的组织，支持设置组织类型、父组织等
// @tags 组织管理
// @param body request body CreateOrganizationRequest true "创建组织请求"
// @success 200 {object} CreateOrganizationResponse
// @failure 400 {object} ErrorResponse "请求参数错误"
// @failure 401 {object} ErrorResponse "未授权"
// @failure 500 {object} ErrorResponse "服务器内部错误"
func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (*CreateOrganizationResponse, error) {
    // ...
}
```

#### 支持的注释标签

| 标签 | 说明 | 示例 |
|------|------|------|
| `@router` | 定义 HTTP 方法和路径 | `@router GET /api/organizations` |
| `@summary` | API 简短描述 | `@summary 获取组织列表` |
| `@description` | 详细说明 | `@description 查询所有组织...` |
| `@tags` | API 分组标签 | `@tags 组织管理` |
| `@param` | 参数定义 | `@param query page int false "页码"` |
| `@success` | 成功响应 | `@success 200 {object} Response` |
| `@failure` | 失败响应 | `@failure 400 {object} ErrorResponse` |

#### 参数定义格式

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

### 详细文档

查看完整的生成器文档：[SWAGGER_GENERATOR.md](./SWAGGER_GENERATOR.md)

包含内容：
- 完整的使用指南
- Handler 注释规范
- CI/CD 集成方案
- 故障排查
- 进阶功能（自定义 Schema、多服务器配置等）

---

## 🔐 认证测试

### 获取Token

```bash
# 使用curl登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "your_password"
  }'

# 响应示例
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200
  }
}
```

### 使用Token

```bash
# 在请求Header中携带Token
curl -X GET http://localhost:8080/api/organizations \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Tenant-ID: tenant_001"
```

---

## 📊 API统计

| 模块 | 接口数量 | 说明 |
|------|---------|------|
| 组织管理 | 10 | 创建、查询、更新、删除、树形结构等 |
| 部门管理 | 10 | 完整的部门层级管理 |
| 员工管理 | 15 | 员工全生命周期管理 |
| 岗位管理 | 10 | 岗位定义和查询 |
| 通讯录服务 | 8 | 组织目录和员工搜索 |
| 认证 | 2 | 登录和Token刷新 |
| **合计** | **63** | 完整的API覆盖 |

---

## 🔗 相关链接

### 项目文档

- [实现差距分析与研发计划](../../../docs/企业级功能完善与统一性设计方案/ZKER-实现差距分析与研发计划_v1.0.md)
- [企业级开发规范手册](../../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [统一错误码定义规范](../../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [API设计规范文档](../../../docs/企业级功能完善与统一性设计方案/API设计规范文档.md)
- [数据库设计完整交付清单](../../../docs/企业级功能完善与统一性设计方案/数据库设计完整交付清单.md)

### 外部资源

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [Postman Learning Center](https://learning.postman.com/)

---

## 💡 使用建议

### 开发环境

1. **本地开发**：使用Swagger UI在线测试
2. **接口调试**：使用Postman Collection
3. **自动化测试**：使用Postman的Collection Runner

### 生产环境

1. **阅读API_GUIDE.md**：了解完整的API规范
2. **使用SDK**：集成客户端SDK简化开发
3. **错误处理**：参考最佳实践章节

---

## 🆘 技术支持

如有问题，请联系：

- 📧 **邮箱**: api-support@coze.com
- 📚 **文档**: https://docs.coze.com
- 🐛 **Bug报告**: https://github.com/coze-dev/coze-studio/issues
- 💬 **讨论**: https://github.com/coze-dev/coze-studio/discussions

---

## 📝 更新日志

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0.0 | 2025-01-01 | 初始版本发布，包含完整的API文档 |

---

## ⭐ 贡献指南

欢迎贡献文档改进！

1. Fork本项目
2. 创建您的特性分支 (`git checkout -b feature/docs-improvement`)
3. 提交您的更改 (`git commit -m 'Add some documentation'`)
4. 推送到分支 (`git push origin feature/docs-improvement`)
5. 开启一个Pull Request

---

**🎉 祝您使用愉快！**
