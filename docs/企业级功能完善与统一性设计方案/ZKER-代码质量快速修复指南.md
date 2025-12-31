# 代码质量快速修复指南

本指南提供快速修复ZKER项目代码质量问题的具体步骤和命令。

## 📋 目录

1. [Panic修复](#1-panic修复)
2. [输入验证添加](#2-输入验证添加)
3. [敏感信息脱敏](#3-敏感信息脱敏)
4. [验证和测试](#4-验证和测试)

---

## 1. Panic修复

### 1.1 快速扫描

```bash
# 运行panic扫描工具
cd /d/code/coze-studio
python tools/fix_panic.py --scan --report

# 查看报告
cat detailed_panic_report.txt
```

### 1.2 批量修复简单Panic

```bash
cd backend

# 1. 替换 "not implemented" panic
find . -name "*.go" -not -name "*_test.go" -not -name "main.go" \
    -exec sed -i 's/panic("not implemented")/return fmt.Errorf("not implemented")/g' {} \;

# 2. 替换简单的error wrap panic
find . -name "*.go" -not -name "*_test.go" \
    -exec sed -i 's/\tpanic(err)/\treturn err/g' {} \;

# 3. 替换 MustGetUIDFromCtx (需手动调整调用方)
# 先找到所有使用位置
grep -rn "MustGetUIDFromCtx" . --include="*.go" | grep -v "_test.go" > mustget_calls.txt
```

### 1.3 手动修复关键文件

#### 文件1: backend/application/base/ctxutil/session.go

**当前代码**:
```go
func MustGetUIDFromCtx(ctx context.Context) int64 {
    sessionData := GetUserSessionFromCtx(ctx)
    if sessionData == nil {
        panic("mustGetUIDFromCtx: sessionData is nil")
    }
    return sessionData.UserID
}
```

**修复后**:
```go
// GetUIDFromCtx 获取用户ID，如果未登录返回error
func GetUIDFromCtx(ctx context.Context) (int64, error) {
    sessionData := GetUserSessionFromCtx(ctx)
    if sessionData == nil {
        return 0, fmt.Errorf("session data is required")
    }
    return sessionData.UserID, nil
}

// MustGetUIDFromCtx 保留此函数用于向后兼容，但内部调用GetUIDFromCtx
// Deprecated: 使用 GetUIDFromCtx 替代
func MustGetUIDFromCtx(ctx context.Context) int64 {
    userID, err := GetUIDFromCtx(ctx)
    if err != nil {
        // 仅在极端情况下panic，记录详细日志
        log.Errorf("MustGetUIDFromCtx failed: %v", err)
        panic(err)
    }
    return userID
}
```

#### 文件2: backend/application/base/ctxutil/api_auth.go

**当前代码**:
```go
func MustGetUIDFromApiAuthCtx(ctx context.Context) int64 {
    apiKeyInfo := GetApiAuthFromCtx(ctx)
    if apiKeyInfo == nil {
        panic("mustGetUIDFromApiAuthCtx: apiKeyInfo is nil")
    }
    return apiKeyInfo.UserID
}
```

**修复后**:
```go
// GetUIDFromApiAuthCtx 从API Auth上下文获取用户ID
func GetUIDFromApiAuthCtx(ctx context.Context) (int64, error) {
    apiKeyInfo := GetApiAuthFromCtx(ctx)
    if apiKeyInfo == nil {
        return 0, fmt.Errorf("api auth info is required")
    }
    return apiKeyInfo.UserID, nil
}

// MustGetUIDFromApiAuthCtx 保留向后兼容
// Deprecated: 使用 GetUIDFromApiAuthCtx 替代
func MustGetUIDFromApiAuthCtx(ctx context.Context) int64 {
    userID, err := GetUIDFromApiAuthCtx(ctx)
    if err != nil {
        log.Errorf("MustGetUIDFromApiAuthCtx failed: %v", err)
        panic(err)
    }
    return userID
}
```

#### 文件3: backend/domain/workflow/internal/nodes/node.go

**当前代码**:
```go
func GetNodeAdaptor(et entity.NodeType) (NodeAdaptor, bool) {
    na, ok := nodeAdaptors[et]
    if !ok {
        panic(fmt.Sprintf("node type %s not registered", et))
    }
    return na(), ok
}
```

**修复后**:
```go
// GetNodeAdaptor 获取节点适配器，如果未注册返回error
func GetNodeAdaptor(et entity.NodeType) (NodeAdaptor, error) {
    na, ok := nodeAdaptors[et]
    if !ok {
        return nil, fmt.Errorf("node type %s not registered", et)
    }
    return na(), nil
}
```

然后更新所有调用方：
```bash
# 查找调用方
grep -rn "GetNodeAdaptor" backend/domain/workflow --include="*.go"

# 手动更新每个调用方处理error
# Before:
// adaptor, _ := nodes.GetNodeAdaptor(nodeType)

# After:
// adaptor, err := nodes.GetNodeAdaptor(nodeType)
// if err != nil {
//     return nil, fmt.Errorf("get node adaptor: %w", err)
// }
```

### 1.4 更新调用方（批量）

```bash
# 1. 找到所有 MustGet 调用
cd backend
grep -rn "MustGetUIDFromCtx\|MustGetUIDFromApiAuthCtx" . \
    --include="*.go" \
    --exclude-dir="vendor" \
    | grep -v "_test.go" \
    | grep -v "func MustGet" \
    > ../mustget_usage.txt

# 2. 逐个文件修复
# 示例：修复某个handler
# Before:
// userID := ctxutil.MustGetUIDFromCtx(ctx)

# After:
// userID, err := ctxutil.GetUIDFromCtx(ctx)
// if err != nil {
//     logs.Errorf("GetUIDFromCtx failed: %v", err)
//     internalServerErrorResponse(ctx, c, err)
//     return
// }
```

---

## 2. 输入验证添加

### 2.1 快速添加tenant_id验证

```bash
cd backend/api/handler/coze

# 1. 找到所有使用 tenant_id 但未验证的文件
grep -L "ValidateTenantID" *.go | xargs grep -l 'c\.Param("tenant_id")'

# 2. 在每个文件顶部添加导入（如果还没有）
# import "github.com/coze-dev/coze-studio/backend/api/handler/coze"

# 3. 在每个c.Param("tenant_id")后添加验证
```

**示例修复**:

```go
// ============================================
// 文件: backend/api/handler/coze/tenant_service.go
// ============================================

// ✅ 在文件顶部添加导入
import (
    "context"
    "fmt"
    "net/http"

    "github.com/cloudwego/hertz/pkg/app"

    "github.com/coze-dev/coze-studio/backend/api/model/tenant"
    tenantapp "github.com/coze-dev/coze-studio/backend/application/tenant"
)

// GetTenant 获取租户详情
// @router /api/tenants/:tenant_id [GET]
func GetTenant(ctx context.Context, c *app.RequestContext) {
    tenantID := c.Param("tenant_id")

    // ✅ 添加验证
    if err := ValidateTenantID(tenantID); err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    resp, err := tenantapp.TenantAppSVC.GetTenant(ctx, tenantID)
    if err != nil {
        internalServerErrorResponse(ctx, c, err)
        return
    }

    c.JSON(http.StatusOK, resp)
}

// UpdateTenant 更新租户信息
// @router /api/tenants/:tenant_id [PUT]
func UpdateTenant(ctx context.Context, c *app.RequestContext) {
    tenantID := c.Param("tenant_id")

    // ✅ 添加验证
    if err := ValidateTenantID(tenantID); err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    var err error
    var req tenant.UpdateTenantRequest
    err = c.BindAndValidate(&req)
    if err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    resp, err := tenantapp.TenantAppSVC.UpdateTenant(ctx, tenantID, &req)
    if err != nil {
        internalServerErrorResponse(ctx, c, err)
        return
    }

    c.JSON(http.StatusOK, resp)
}
```

### 2.2 批量添加limit/offset验证

```bash
cd backend/api/handler/coze

# 创建一个辅助脚本来批量添加验证
cat > add_validation.sh << 'EOF'
#!/bin/bash

# 在包含 limit 参数的函数中添加验证
FILES=$(grep -l 'c\.Query("limit")' *.go)

for file in $FILES; do
    echo "Processing $file..."

    # 检查是否已经有ValidateLimit
    if grep -q "ValidateLimit" "$file"; then
        echo "  ✓ Already has ValidateLimit"
        continue
    fi

    echo "  ⚠ Need to add ValidateLimit manually"
done
EOF

chmod +x add_validation.sh
./add_validation.sh
```

**示例修复**:

```go
// ============================================
// 文件: backend/api/handler/coze/monitoring/monitoring_handler.go
// ============================================

import (
    "github.com/coze-dev/coze-studio/backend/api/handler/coze"  // ✅ 添加导入
)

func GetMonitoringMetrics(ctx context.Context, c *app.RequestContext) {
    tenantID := c.Param("tenant_id")

    // ✅ 验证tenant_id
    if err := coze.ValidateTenantID(tenantID); err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    // ✅ 验证limit
    limit, err := coze.ValidateLimit(c.Query("limit"), 20)
    if err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    // ✅ 验证offset
    offset, err := coze.ValidateOffset(c.Query("offset"))
    if err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    // 使用验证后的参数...
}
```

### 2.3 添加日期验证

```go
// 示例：验证日期范围参数
func GetTokenUsage(ctx context.Context, c *app.RequestContext) {
    tenantID := c.Param("tenant_id")

    // 验证tenant_id
    if err := ValidateTenantID(tenantID); err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    // ✅ 验证日期格式
    startDate, err := ValidateDate(c.Query("start_date"))
    if err != nil {
        invalidParamRequestResponse(c,
            "invalid start_date format, expected YYYY-MM-DD")
        return
    }

    endDate, err := ValidateDate(c.Query("end_date"))
    if err != nil {
        invalidParamRequestResponse(c,
            "invalid end_date format, expected YYYY-MM-DD")
        return
    }

    // ✅ 验证日期范围
    if !endDate.IsZero() && !startDate.IsZero() && endDate.Before(startDate) {
        invalidParamRequestResponse(c,
            "end_date must be after start_date")
        return
    }

    // 使用验证后的日期...
}
```

---

## 3. 敏感信息脱敏

### 3.1 添加日志脱敏

```go
// ============================================
// 在需要记录日志的地方添加脱敏
// ============================================

import (
    "github.com/coze-dev/coze-studio/backend/api/middleware"
    logs "github.com/sirupsen/logrus"
)

// ❌ 修复前
func HandleLogin(ctx context.Context, c *app.RequestContext) {
    var req LoginRequest
    c.BindAndValidate(&req)

    // 记录敏感信息
    logs.Infof("User login: username=%s, password=%s",
        req.Username, req.Password)  // ❌
}

// ✅ 修复后
func HandleLogin(ctx context.Context, c *app.RequestContext) {
    var req LoginRequest
    c.BindAndValidate(&req)

    // 不记录密码
    logs.Infof("User login attempt: username=%s", req.Username)  // ✅

    // 或者使用脱敏工具
    auditLog := map[string]interface{}{
        "username": req.Username,
        "password": req.Password,
        "ip":       c.ClientIP(),
    }

    // 脱敏后记录
    sanitizedLog := middleware.SanitizeAuditLog(auditLog)
    logs.Infof("Login audit: %+v", sanitizedLog)
}
```

### 3.2 错误响应脱敏

```go
// ============================================
// 在错误处理中添加脱敏
// ============================================

import (
    "github.com/coze-dev/coze-studio/backend/api/middleware"
    logs "github.com/sirupsen/logrus"
)

func HandleDatabaseOperation(ctx context.Context, c *app.RequestContext) {
    result, err := database.Query(...)
    if err != nil {
        // ✅ 详细错误记录到日志
        logs.Errorf("Database query failed: %v", err)

        // ✅ 过滤敏感信息后返回给用户
        sanitizedErrMsg := middleware.SanitizeError(err)

        httputil.BuildErrorResp(c,
            errno.ErrInternalErrorCode,
            sanitizedErrMsg,  // 使用过滤后的消息
            "数据库操作失败",
            nil)
        return
    }

    c.JSON(http.StatusOK, result)
}
```

### 3.3 审计日志脱敏

```go
// ============================================
// 在审计日志中间件中添加脱敏
// ============================================

import (
    "github.com/coze-dev/coze-studio/backend/api/middleware"
)

func AuditLoggingMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 读取请求体
        body := c.Request.Body()
        requestURI := c.Request.URI().RequestURI()

        // ✅ 过滤敏感字段
        sanitizedBody := middleware.SanitizeRequest(body)

        // 记录审计日志
        auditLog := map[string]interface{}{
            "timestamp": time.Now().Format(time.RFC3339),
            "method":    c.Request.Method,
            "path":      requestURI,
            "body":      string(sanitizedBody),
            "ip":        c.ClientIP(),
        }

        // 使用脱敏后的日志
        auditLogger.WithFields(logrus.Fields{
            "audit": auditLog,
        }).Info("API request")

        c.Next(ctx)
    }
}
```

---

## 4. 验证和测试

### 4.1 代码检查

```bash
cd backend

# 1. Panic检查
echo "=== Panic检查 ==="
grep -rn "panic(" . --include="*.go" | \
    grep -v "_test.go" | \
    grep -v "main.go" | \
    wc -l

# 2. 输入验证检查
echo "=== tenant_id验证检查 ==="
grep -rn 'c\.Param("tenant_id")' api/handler --include="*.go" | \
    grep -v 'ValidateTenantID' | \
    wc -l

echo "=== limit验证检查 ==="
grep -rn 'c\.Query("limit")' api/handler --include="*.go" | \
    grep -v 'ValidateLimit' | \
    wc -l

# 3. 敏感信息检查
echo "=== 敏感信息日志检查 ==="
grep -rn 'logs\..*password\|logs\..*token' . --include="*.go" | \
    wc -l
```

### 4.2 运行测试

```bash
cd backend

# 1. 单元测试
echo "=== 运行单元测试 ==="
go test ./... -v -cover

# 2. 竞态检测
echo "=== 竞态检测 ==="
go test ./... -race

# 3. 覆盖率报告
echo "=== 生成覆盖率报告 ==="
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 4. Linter检查
echo "=== Linter检查 ==="
golangci-lint run --timeout=10m
```

### 4.3 构建和部署测试

```bash
# 1. 构建检查
cd backend
go build ./...

# 2. Docker构建测试
cd ..
docker build -f docker/Dockerfile.backend -t coze-studio:test .

# 3. 运行容器测试
docker-compose -f docker/docker-compose.test.yml up -d
```

---

## 5. 常见修复模式

### 5.1 MustGet → Get (带error处理)

```go
// ❌ Before
func SomeHandler(ctx context.Context, c *app.RequestContext) {
    userID := ctxutil.MustGetUIDFromCtx(ctx)
    // 使用userID...
}

// ✅ After
func SomeHandler(ctx context.Context, c *app.RequestContext) {
    userID, err := ctxutil.GetUIDFromCtx(ctx)
    if err != nil {
        logs.Errorf("GetUIDFromCtx failed: %v", err)
        internalServerErrorResponse(ctx, c, err)
        return
    }
    // 使用userID...
}
```

### 5.2 无验证 → 有验证

```go
// ❌ Before
func GetBot(ctx context.Context, c *app.RequestContext) {
    botID := c.Param("bot_id")
    // 直接使用botID...
}

// ✅ After
func GetBot(ctx context.Context, c *app.RequestContext) {
    botID := c.Param("bot_id")

    // 添加验证
    if err := ValidateUUID(botID); err != nil {
        invalidParamRequestResponse(c, err.Error())
        return
    }

    // 使用验证后的botID...
}
```

### 5.3 敏感日志 → 脱敏日志

```go
// ❌ Before
logs.Infof("Request: %+v", request)  // 可能包含密码

// ✅ After
sanitizedRequest := middleware.SanitizeMap(request)
logs.Infof("Request: %+v", sanitizedRequest)
```

---

## 6. IDE批量修复技巧

### VSCode

1. **正则查找和替换**:
   - 按 `Ctrl+H` 打开替换
   - 启用正则表达式模式
   - 查找: `panic\("not implemented"\)`
   - 替换: `return fmt.Errorf("not implemented")`

2. **多文件编辑**:
   - 按 `Ctrl+Shift+F` 全局搜索
   - 搜索: `c\.Param\("tenant_id"\)`
   - 在所有匹配的文件中添加验证

### GoLand

1. **结构化查找和替换**:
   - `Edit` → `Find` → `Replace Structurally`
   - 创建模板：`panic($msg$)` → `return fmt.Errorf($msg$)`

2. **意图检查**:
   - 右键 → `Run Inspection by Name`
   - 搜索: "Panic"
   - 批量修复所有问题

---

## 7. 进度跟踪

创建一个进度跟踪文件：

```bash
cat > backend/quality_fix_progress.md << 'EOF'
# 代码质量改进进度

## Panic修复进度

- [x] 创建验证函数库
- [ ] 修复 session_check panic (2处)
- [ ] 修复 registry_check panic (2处)
- [ ] 修复 not_implemented panic (1处)
- [ ] 修复 error_wrap panic (24处)
- [ ] 修复 data_validation panic (42处)

## 输入验证进度

- [x] 创建验证函数库
- [ ] tenant_id验证 (59处)
- [ ] bot_id验证 (45处)
- [ ] limit/offset验证 (63处)
- [ ] 日期验证 (15处)

## 敏感信息脱敏进度

- [x] 创建脱敏工具库
- [ ] 日志脱敏 (5处)
- [ ] 错误响应脱敏 (10处)
- [ ] 审计日志脱敏 (3处)
EOF
```

---

**下一步**: 开始修复第一组panic，预计完成时间：2小时
