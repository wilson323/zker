# API路由注册指南

## 租户注册相关API路由

### 路由定义

以下路由需要在 `backend/api/router/register.go` 中注册：

#### 1. 发送验证码
```
POST /api/v1/tenants/send-verification-code
Handler: SendVerificationCode
```

#### 2. 检查企业名称可用性
```
GET /api/v1/tenants/check-company-name
Handler: CheckCompanyNameAvailability
Query: company_name (string)
```

#### 3. 检查子域名可用性
```
GET /api/v1/tenants/check-subdomain
Handler: CheckSubdomainAvailability
Query: subdomain (string)
```

### 注册示例代码

```go
import (
    handler "github.com/coze-dev/coze-studio/backend/api/handler/coze"
)

// 租户注册路由组
tenantRegistrationGroup := engine.Group("/api/v1/tenants")
{
    tenantRegistrationGroup.POST("/send-verification-code", handler.SendVerificationCode)
    tenantRegistrationGroup.GET("/check-company-name", handler.CheckCompanyNameAvailability)
    tenantRegistrationGroup.GET("/check-subdomain", handler.CheckSubdomainAvailability)
}
```

### 依赖注入

需要在容器中注册以下服务：

1. **TenantValidationService** - 租户验证服务
   - 依赖: TenantRepository

2. **VerificationCodeService** - 验证码服务
   - 依赖: Redis Client, EmailService

3. **EmailService** - 邮件服务
   - 依赖: 配置

### 测试API

#### 1. 检查企业名称可用性
```bash
curl -X GET "http://localhost:8888/api/v1/tenants/check-company-name?company_name=测试公司"
```

响应示例:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "available": true,
    "suggestions": []
  },
  "timestamp": 1704067200000
}
```

#### 2. 检查子域名可用性
```bash
curl -X GET "http://localhost:8888/api/v1/tenants/check-subdomain?subdomain=testcompany"
```

响应示例:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "available": true,
    "suggestions": [],
    "full_domain": "testcompany.saas.coze.com"
  },
  "timestamp": 1704067200000
}
```

#### 3. 发送验证码
```bash
curl -X POST "http://localhost:8888/api/v1/tenants/send-verification-code" \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

响应示例:
```json
{
  "code": 0,
  "message": "验证码已发送",
  "data": {
    "expires_in": 300,
    "message": "验证码有效期为5分钟"
  },
  "timestamp": 1704067200000
}
```

### 错误响应

#### 400 Bad Request - 参数错误
```json
{
  "code": 400,
  "message": "企业名称不能为空",
  "error": {
    "code": "INVALID_PARAMETER",
    "description": "..."
  }
}
```

#### 400 Bad Request - 子域名格式错误
```json
{
  "code": 400,
  "message": "子域名格式不正确，仅支持小写字母、数字、连字符",
  "error": {
    "code": "INVALID_SUBDOMAIN",
    "description": "子域名长度为3-63个字符，仅支持小写字母、数字和连字符"
  }
}
```

#### 500 Internal Server Error - 服务错误
```json
{
  "code": 500,
  "message": "检查企业名称失败",
  "error": {
    "code": "CHECK_FAILED",
    "description": "..."
  }
}
```

## 下一步

完成以下任务以使API完全可用：

1. ✅ Service层实现
2. ✅ Handler层实现
3. ✅ 单元测试
4. ⏳ 路由注册（需要在router/register.go中添加）
5. ⏳ 依赖注入配置（需要在容器中注册服务）
6. ⏳ 集成测试（需要真实数据库环境）
