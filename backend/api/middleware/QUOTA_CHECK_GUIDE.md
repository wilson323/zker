# 配额检查中间件使用指南

> **版本**: v1.0  
> **最后更新**: 2025-01-01  
> **维护者**: 研发B - 后端工程师

---

## 📋 目录

- [功能概述](#功能概述)
- [资源类型映射](#资源类型映射)
- [集成步骤](#集成步骤)
- [错误响应格式](#错误响应格式)
- [性能优化](#性能优化)
- [测试验证](#测试验证)

---

## 功能概述

配额检查中间件（`QuotaCheckMiddleware`）用于在创建资源时自动检查租户配额，确保不会超出订阅限制。

### 核心特性

- ✅ **自动路径识别**: 根据请求路径自动识别资源类型
- ✅ **统一错误码**: 使用 errno 包的错误码（QUOTA402001-402004）
- ✅ **双语支持**: 中英文错误消息
- ✅ **性能监控**: 自动记录检查耗时
- ✅ **灵活跳过**: 无配额限制的路径自动跳过

### 性能要求

- **P95延迟**: < 20ms
- **P99延迟**: < 50ms
- **吞吐量**: > 1000 QPS

---

## 资源类型映射

| 资源类型 | 路径前缀 | 错误码 | 说明 |
|---------|---------|--------|------|
| `bots` | `/api/v1/bots` | QUOTA402001 | Bot数量配额 |
| `messages` | `/api/v1/conversations` | QUOTA402002 | 消息数量配额 |
| `storage` | `/api/v1/knowledge` | QUOTA402003 | 存储空间配额 |
| `workflows` | `/api/v1/workflows` | QUOTA402004 | 工作流数量配额 |

### 无需配额检查的路径

- `/api/health/*` - 健康检查
- `/api/v1/auth/*` - 认证相关
- `/api/v1/bots/:id` - Bot查询（GET/PUT/DELETE）
- `/api/passport/*` - 用户注册/登录

---

## 集成步骤

### 步骤 1: 创建 QuotaAppService 实例

```go
// backend/api/server.go 或 main.go
import (
    "github.com/coze-studio/backend/application/tenant"
    "github.com/coze-studio/backend/domain/tenant/repository"
)

// 创建配额应用服务
quotaRepo := repository.NewQuotaRepository(db)
quotaAppSvc := tenant.NewQuotaAppService(quotaRepo)
```

### 步骤 2: 在路由中注册中间件

由于项目使用 Hertz 自动生成路由（基于 IDL），有两种集成方式：

#### 方式 A: 修改 middleware.go（推荐）

在 `backend/api/router/coze/middleware.go` 中为需要配额检查的路由添加中间件：

```go
// backend/api/router/coze/middleware.go

// 为创建 Bot 添加配额检查
func _draftbotcreateMw() []app.HandlerFunc {
    return []app.HandlerFunc{
        middleware.QuotaCheckMiddleware(quotaAppSvc),
    }
}

// 为创建对话添加配额检查
func _createconversationMw() []app.HandlerFunc {
    return []app.HandlerFunc{
        middleware.QuotaCheckMiddleware(quotaAppSvc),
    }
}

// 为创建数据集添加配额检查
func _createdatasetMw() []app.HandlerFunc {
    return []app.HandlerFunc{
        middleware.QuotaCheckMiddleware(quotaAppSvc),
    }
}

// 为创建工作流添加配额检查
func _createworkflowMw() []app.HandlerFunc {
    return []app.HandlerFunc{
        middleware.QuotaCheckMiddleware(quotaAppSvc),
    }
}
```

#### 方式 B: 手动路由注册（适用于自定义路由）

如果需要手动注册路由，在 `backend/api/router/register.go` 中添加：

```go
// backend/api/router/register.go

import (
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/coze-studio/backend/api/handler"
    "github.com/coze-studio/backend/api/middleware"
)

func Register(r *server.Hertz) {
    // ... 其他路由
    
    // Bot 路由（带配额检查）
    botGroup := r.Group("/api/v1/bots")
    botGroup.POST("",
        middleware.QuotaCheckMiddleware(quotaAppSvc),
        handler.CreateBot,
    )
    
    // 对话路由（带配额检查）
    convGroup := r.Group("/api/v1/conversations")
    convGroup.POST("",
        middleware.QuotaCheckMiddleware(quotaAppSvc),
        handler.CreateConversation,
    )
    
    // 工作流路由（带配额检查）
    workflowGroup := r.Group("/api/v1/workflows")
    workflowGroup.POST("",
        middleware.QuotaCheckMiddleware(quotaAppSvc),
        handler.CreateWorkflow,
    )
}
```

### 步骤 3: 确保 tenant_id 在 Context 中

配额检查中间件依赖 `tenant_id` 从 Context 中提取。确保在认证中间件中设置了 `tenant_id`：

```go
// backend/api/middleware/auth.go

func AuthMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // ... 验证 Token
        
        // 从 Token 中提取 tenant_id
        tenantID := extractTenantIDFromToken(token)
        
        // 设置到 Context 中（重要！）
        ctx = context.WithValue(ctx, "tenant_id", tenantID)
        
        c.Next(ctx)
    }
}
```

---

## 错误响应格式

### 配额超限响应（402）

```json
{
  "code": "QUOTA402001",
  "message": "Bot数量配额已用完，请升级订阅",
  "message_zh": "Bot数量配额已用完，请升级订阅",
  "message_en": "Bot quota exceeded, please upgrade your subscription",
  "resource_type": "bots",
  "current_usage": 10,
  "max_limit": 10,
  "usage_percent": "100.00%",
  "request_id": "req_1735689600000000000"
}
```

### 配额服务错误响应（500）

```json
{
  "code": "QUOTA500001",
  "message": "配额检查失败",
  "message_zh": "配额检查失败",
  "message_en": "Failed to check quota"
}
```

---

## 性能优化

### 1. 路径缓存优化

资源类型识别使用 `strings.HasPrefix`，性能已经很好。但如果需要进一步优化：

```go
// backend/api/middleware/quota_check.go

// 使用 sync.Map 缓存路径-资源类型映射
var pathCache sync.Map

func getResourceTypeFromPathCached(path string) string {
    // 尝试从缓存获取
    if val, ok := pathCache.Load(path); ok {
        return val.(string)
    }
    
    // 缓存未命中，计算并缓存
    resourceType := getResourceTypeFromPath(path)
    pathCache.Store(path, resourceType)
    return resourceType
}
```

### 2. 配额缓存优化

在 `QuotaAppService` 中添加 Redis 缓存：

```go
// backend/application/tenant/quota_app.go

func (s *QuotaAppService) CheckQuota(ctx context.Context, req *CheckQuotaRequest) (*CheckQuotaResponse, error) {
    // 1. 尝试从 Redis 缓存获取
    cacheKey := fmt.Sprintf("quota:%s:%s", req.TenantID, req.ResourceType)
    cached, err := s.redis.Get(ctx, cacheKey)
    if err == nil && cached != nil {
        return cached, nil
    }
    
    // 2. 缓存未命中，从数据库查询
    quota, err := s.quotaRepo.GetByTenantAndResource(ctx, req.TenantID, entity.ResourceType(req.ResourceType))
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存（TTL: 60秒）
    s.redis.Set(ctx, cacheKey, quota, 60*time.Second)
    
    return s.buildResponse(quota), nil
}
```

### 3. 数据库索引优化

确保 `quotas` 表有以下索引：

```sql
-- 租户+资源类型联合索引（最重要的索引）
CREATE INDEX idx_tenant_resource ON quotas(tenant_id, resource_type);

-- 使用量索引（用于配额监控查询）
CREATE INDEX idx_used_count ON quotas(used_count);

-- 复合索引（用于复杂查询）
CREATE INDEX idx_tenant_resource_usage ON quotas(tenant_id, resource_type, used_count);
```

---

## 测试验证

### 单元测试

运行单元测试：

```bash
cd backend
go test ./api/middleware -v -run TestQuotaCheck
```

预期结果：
- ✅ 所有测试用例通过
- ✅ 测试覆盖率 ≥ 80%

### 集成测试

测试配额检查中间件是否正常工作：

```bash
# 1. 启动服务
make server

# 2. 发送创建 Bot 请求（配额充足）
curl -X POST http://localhost:8001/api/v1/bots \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"bot_name": "test_bot"}'

# 预期: 200 OK

# 3. 发送创建 Bot 请求（配额超限）
# 先将配额用满，然后再次请求
curl -X POST http://localhost:8001/api/v1/bots \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"bot_name": "test_bot"}'

# 预期: 402 Payment Required
# 响应体包含: {"code": "QUOTA402001", ...}
```

### 性能测试

使用 K6 进行性能测试：

```bash
cd tests/performance/k6
k6 run quota_test.js
```

预期结果：
- ✅ P95延迟 < 20ms
- ✅ 错误率 < 1%

---

## 常见问题

### Q1: 中间件不生效？

**原因**: `tenant_id` 未在 Context 中设置

**解决**: 确保认证中间件在配额检查中间件之前执行，并设置了 `tenant_id`

```go
// 错误的顺序
r.Use(QuotaCheckMiddleware)  // tenant_id 还未设置
r.Use(AuthMiddleware)        // 在这里才设置 tenant_id

// 正确的顺序
r.Use(AuthMiddleware)        // 先设置 tenant_id
r.Use(QuotaCheckMiddleware)  // 再检查配额
```

### Q2: 配额检查返回 500 错误？

**原因**: 配额服务初始化失败或数据库连接失败

**解决**: 检查日志，确认 `QuotaAppService` 正确初始化

```bash
# 查看日志
tail -f logs/coze-studio.log | grep QUOTA
```

### Q3: 性能未达到预期？

**原因**: 数据库查询慢或缺少索引

**解决**: 
1. 检查数据库慢查询日志
2. 确保索引已创建（见"数据库索引优化"章节）
3. 启用 Redis 缓存（见"配额缓存优化"章节）

---

## 维护清单

### 每日检查

- [ ] 监控配额检查中间件性能指标（P95延迟、错误率）
- [ ] 检查日志中的配额超限告警

### 每周检查

- [ ] 审查配额超限次数 TOP10 的租户
- [ ] 优化慢查询

### 每月检查

- [ ] 审查配额使用趋势
- [ ] 评估是否需要调整配额限制

---

## 附录

### 相关文档

- [ZKER-统一错误码定义规范.md](../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [ZKER-企业级开发规范手册_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [API接口文档_租户计费系统.md](../../docs/企业级功能完善与统一性设计方案/API接口文档_租户计费系统.md)

### 代码文件

- `backend/api/middleware/quota_check.go` - 中间件实现
- `backend/api/middleware/quota_check_test.go` - 单元测试
- `backend/application/tenant/quota_app.go` - 应用层服务
- `backend/domain/tenant/service/quota_service.go` - 领域层服务

---

**版权所有 © 2025 Coze Studio Enterprise Team**
