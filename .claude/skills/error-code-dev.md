# 统一错误码开发助手

**版本**: v3.0.0 | **更新**: 2025-01-03

协助开发者正确使用 ZKER 的 300+ 统一错误码系统。

---

## 🎯 使用场景

- 查找错误码
- 创建新错误码
- 错误处理
- 错误包装
- 国际化错误

---

## 📦 错误码体系

### 错误码格式

```
格式: {模块编号}{错误类型}{三位数字}
示例: 2001001 (租户不存在), 500101 (Bot配额超限)
```

### 常用错误码

| 错误码 | 常量名 | 说明 | HTTP状态码 |
|--------|--------|------|-----------|
| 2001001 | ErrTenantNotFoundCode | 租户不存在 | 404 |
| 2004001 | ErrTenantQuotaExceededCode | 租户配额已用尽 | 403 |
| 5000001 | ErrQuotaExceededCode | 通用配额超限 | 429 |
| 5001001 | ErrQuotaBotExceededCode | Bot配额超限 | 429 |

---

## 💻 代码模板

### 1. 返回错误响应

```go
// ✅ Good
func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*Tenant, error) {
    tenant, err := s.repo.GetByTenantID(ctx, tenantID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errorx.New(errno.ErrTenantNotFoundCode)
        }
        return nil, errorx.Wrapf(err, errno.ErrTenantCheckFailed)
    }
    return tenant, nil
}
```

### 2. 创建新错误码

```go
// 1. 定义错误码常量
const (
    ErrTenantInvalidStatusCode = 2002040
)

// 2. 定义错误变量
var (
    ErrTenantInvalidStatus = &BaseErrorCode{
        code:       "TENANT_INVALID_STATUS",
        message:    "Invalid tenant status",
        messageZH:  "租户状态无效",
        httpStatus: http.StatusBadRequest,
    }
)
```

---

## 📋 检查清单

### 错误码定义

- [ ] 错误码遵循命名规范
- [ ] 包含中英文双语消息
- [ ] 设置正确的 HTTP 状态码

### 错误处理

- [ ] 使用统一错误码常量
- [ ] 错误包装使用 errorx.Wrapf
- [ ] 添加上下文信息

---

## 📖 相关文档

- [ZKER-统一错误码定义规范](../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [02-SPECS/error-handling.md](../../docs/02-SPECS/error-handling.md)

---

**🎯 目标**: 确保 300+ 错误码正确使用，实现全局一致性！
