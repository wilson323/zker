# 统一错误码系统文档

> **版本**: v1.0
> **最后更新**: 2025-01-01
> **维护者**: 研发B - 后端工程师

## 📋 目录

- [概述](#概述)
- [错误码设计规范](#错误码设计规范)
- [通用错误码](#通用错误码)
- [认证模块错误码 (AUTH)](#认证模块错误码-auth)
- [权限模块错误码 (PERMISSION)](#权限模块错误码-permission)
- [租户模块错误码 (TENANT)](#租户模块错误码-tenant)
- [配额模块错误码 (QUOTA)](#配额模块错误码-quota)
- [路由模块错误码 (ROUTING)](#路由模块错误码-routing)
- [Bot模块错误码 (BOT)](#bot模块错误码-bot)
- [对话模块错误码 (CHAT)](#对话模块错误码-chat)
- [订阅模块错误码 (SUBSCRIPTION)](#订阅模块错误码-subscription)
- [用户模块错误码 (USER)](#用户模块错误码-user)
- [错误响应格式](#错误响应格式)
- [使用示例](#使用示例)

---

## 概述

本系统采用统一的错误码设计，覆盖所有业务模块，支持中英文双语，提供清晰的错误信息和HTTP状态码映射。

### 错误码覆盖范围

| 模块 | 错误码数量 | 文件位置 |
|------|-----------|----------|
| 通用错误码 | 15 | common.go |
| 认证模块 | 40+ | auth.go |
| 权限模块 | 40+ | permission.go |
| 租户模块 | 20+ | tenant.go |
| 配额模块 | 15+ | quota.go |
| 路由模块 | 10+ | routing.go |
| Bot模块 | 15+ | bot.go |
| 对话模块 | 10+ | chat.go |
| 订阅模块 | 10+ | subscription.go |
| 用户模块 | 15+ | user.go |
| **总计** | **300+** | - |

---

## 错误码设计规范

### 错误码格式

```
{MODULE}{HTTP_STATUS}{SEQUENTIAL_NUMBER}
```

**示例**: `AUTH401001`
- `AUTH`: 模块标识（认证模块）
- `401`: HTTP状态码（未授权）
- `001`: 序号

### HTTP状态码映射

| HTTP状态 | 说明 | 使用场景 |
|---------|------|---------|
| 200 | 成功 | 操作成功 |
| 201 | 参数错误 | 客户端请求参数有误 |
| 400 | 错误请求 | 一般性客户端错误 |
| 401 | 未授权 | 未登录或Token无效 |
| 402 | 需要付费 | 配额用完，需升级订阅 |
| 403 | 禁止访问 | 权限不足 |
| 404 | 未找到 | 资源不存在 |
| 409 | 冲突 | 资源已存在或冲突 |
| 422 | 无法处理 | 业务逻辑验证失败 |
| 423 | 锁定 | 资源被锁定 |
| 424 | 依赖失败 | 依赖操作失败 |
| 426 | 需要升级 | 协议或版本需要升级 |
| 500 | 服务器错误 | 服务器内部错误 |

### 模块标识

| 标识 | 模块名称 | 说明 |
|------|---------|------|
| COMMON | 通用 | 通用错误码 |
| AUTH | 认证 | 登录、Token、会话管理 |
| PERMISSION | 权限 | RBAC、数据权限 |
| TENANT | 租户 | 多租户管理 |
| QUOTA | 配额 | 资源配额管理 |
| ROUTING | 路由 | 智能路由、意图识别 |
| BOT | Bot | Bot管理 |
| CHAT | 对话 | 对话、消息管理 |
| SUBSCRIPTION | 订阅 | 订阅、计费 |
| USER | 用户 | 用户管理 |

---

## 通用错误码

### 成功响应 (200)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| SUCCESS | 操作成功 | Operation successful |

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| COMMON201001 | 参数错误 | Invalid parameters |
| COMMON201002 | 参数缺失 | Missing parameters |
| COMMON201003 | 参数格式错误 | Invalid parameter format |
| COMMON201004 | 必填参数缺失 | Required parameter missing |
| COMMON201005 | 参数值超出范围 | Parameter value out of range |

### 请求错误 (400)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| COMMON400001 | 请求体格式错误 | Invalid request body format |
| COMMON400002 | Content-Type错误 | Invalid Content-Type |
| COMMON400003 | 请求方法不支持 | Method not allowed |

### 未授权错误 (401)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| COMMON401001 | 未登录 | Not logged in |
| COMMON401002 | Token已过期 | Token has expired |
| COMMON401003 | 无效的Token | Invalid token |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| COMMON403001 | 权限不足 | Permission denied |
| COMMON403002 | 访问被拒绝 | Access denied |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| COMMON404001 | 资源不存在 | Resource not found |
| COMMON404002 | 接口不存在 | API endpoint not found |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| COMMON500001 | 服务器内部错误 | Internal server error |
| COMMON500002 | 服务暂不可用 | Service temporarily unavailable |
| COMMON500003 | 数据库错误 | Database error |

---

## 认证模块错误码 (AUTH)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH201001 | 用户名或密码错误 | Invalid username or password |
| AUTH201002 | 用户名已存在 | Username already exists |
| AUTH201003 | 邮箱格式错误 | Invalid email format |
| AUTH201004 | 手机号格式错误 | Invalid phone number format |
| AUTH201005 | 密码格式错误 | Invalid password format |
| AUTH201006 | 验证码错误 | Invalid verification code |
| AUTH201007 | 验证码已过期 | Verification code has expired |
| AUTH201008 | 注册参数缺失 | Missing registration parameters |
| AUTH201009 | 用户名长度错误 | Username length must be between 4-20 characters |
| AUTH201010 | 密码过于简单 | Password is too simple |

### 认证错误 (401)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH401001 | 未登录 | Not logged in |
| AUTH401002 | Token已过期 | Token has expired |
| AUTH401003 | 无效的Token | Invalid token |
| AUTH401004 | Token缺失 | Missing token |
| AUTH401005 | Refresh Token无效 | Invalid refresh token |
| AUTH401006 | Refresh Token已过期 | Refresh token has expired |
| AUTH401007 | 登录已过期 | Login has expired |
| AUTH401008 | 会话已失效 | Session has been invalidated |
| AUTH401009 | 重复登录 | Duplicate login detected |
| AUTH401010 | 登录失败次数过多 | Too many failed login attempts |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH403001 | 账号已禁用 | Account has been disabled |
| AUTH403002 | 账号已锁定 | Account has been locked |
| AUTH403003 | 账号已封禁 | Account has been banned |
| AUTH403004 | 权限不足 | Permission denied |
| AUTH403005 | 需要管理员权限 | Admin permission required |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH404001 | 用户不存在 | User not found |
| AUTH404002 | 角色不存在 | Role not found |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH500001 | 认证服务异常 | Authentication service error |
| AUTH500002 | Token生成失败 | Failed to generate token |
| AUTH500003 | 密码加密失败 | Failed to encrypt password |
| AUTH500004 | 会话创建失败 | Failed to create session |
| AUTH500005 | 验证码发送失败 | Failed to send verification code |
| AUTH500006 | 用户创建失败 | Failed to create user |
| AUTH500007 | 用户更新失败 | Failed to update user |
| AUTH500008 | 用户删除失败 | Failed to delete user |

### 第三方登录错误 (422)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH422001 | 第三方登录失败 | Third-party login failed |
| AUTH422002 | 第三方授权失败 | Third-party authorization failed |
| AUTH422003 | 未绑定账号 | Account not linked |
| AUTH422004 | 已绑定其他账号 | Already linked to another account |

### OAuth相关错误 (426)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH426001 | 无效的OAuth客户端 | Invalid OAuth client |
| AUTH426002 | OAuth授权码无效 | Invalid OAuth authorization code |
| AUTH426003 | OAuth scope无效 | Invalid OAuth scope |

### 账号安全相关错误 (423)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH423001 | 密码过于简单 | Password is too simple |
| AUTH423002 | 密码与旧密码相同 | New password must be different from old password |
| AUTH423003 | 旧密码错误 | Old password is incorrect |
| AUTH423004 | 两次密码不一致 | Passwords do not match |
| AUTH423005 | 需要验证旧密码 | Old password verification required |
| AUTH423006 | 密码重置次数过多 | Too many password reset attempts |
| AUTH423007 | 验证码错误次数过多 | Too many incorrect verification code attempts |

### 邮箱验证相关错误 (424)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| AUTH424001 | 邮箱未验证 | Email has not been verified |
| AUTH424002 | 邮箱已验证 | Email has already been verified |
| AUTH424003 | 验证链接已失效 | Verification link has expired |
| AUTH424004 | 验证链接无效 | Invalid verification link |

---

## 权限模块错误码 (PERMISSION)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| PERMISSION201001 | 角色名称格式错误 | Invalid role name format |
| PERMISSION201002 | 权限ID格式错误 | Invalid permission ID format |
| PERMISSION201003 | 数据权限范围错误 | Invalid data permission scope |
| PERMISSION201004 | 角色已存在 | Role already exists |
| PERMISSION201005 | 权限已存在 | Permission already exists |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| PERMISSION403001 | 无权限访问 | No permission to access this resource |
| PERMISSION403002 | 无权限操作 | No permission to perform this operation |
| PERMISSION403003 | 无权限查看 | No permission to view this data |
| PERMISSION403004 | 无权限创建 | No permission to create resource |
| PERMISSION403005 | 无权限编辑 | No permission to edit resource |
| PERMISSION403006 | 无权限删除 | No permission to delete resource |
| PERMISSION403007 | 无权限导出 | No permission to export data |
| PERMISSION403008 | 无权限导入 | No permission to import data |
| PERMISSION403009 | 超出数据权限范围 | Data permission scope exceeded |
| PERMISSION403010 | 租户隔离限制 | Cross-tenant access denied |
| PERMISSION403011 | 需要管理员权限 | Admin permission required |
| PERMISSION403012 | 需要所有者权限 | Owner permission required |
| PERMISSION403020 | 只能访问自己的数据 | Can only access own data |
| PERMISSION403021 | 只能访问部门数据 | Can only access department data |
| PERMISSION403022 | 只能访问团队数据 | Can only access team data |
| PERMISSION403023 | 无权查看敏感字段 | No permission to view sensitive field |
| PERMISSION403024 | 无权编辑敏感字段 | No permission to edit sensitive field |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| PERMISSION404001 | 角色不存在 | Role not found |
| PERMISSION404002 | 权限不存在 | Permission not found |
| PERMISSION404003 | 用户角色不存在 | User role not found |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| PERMISSION500001 | 权限检查失败 | Failed to check permission |
| PERMISSION500002 | 角色创建失败 | Failed to create role |
| PERMISSION500003 | 角色更新失败 | Failed to update role |
| PERMISSION500004 | 角色删除失败 | Failed to delete role |
| PERMISSION500005 | 权限分配失败 | Failed to assign permission |
| PERMISSION500006 | 权限回收失败 | Failed to revoke permission |
| PERMISSION500007 | 用户角色创建失败 | Failed to create user role |
| PERMISSION500008 | 用户角色删除失败 | Failed to delete user role |

### 角色管理相关错误 (409)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| PERMISSION409001 | 系统角色不可删除 | System role cannot be deleted |
| PERMISSION409002 | 角色正在使用中 | Role is currently in use |
| PERMISSION409003 | 无法删除最后一个管理员 | Cannot delete the last admin |
| PERMISSION409004 | 无法移除所有权限 | Role must have at least one permission |
| PERMISSION409005 | 权限冲突 | Permission conflict |

---

## 租户模块错误码 (TENANT)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| TENANT201001 | 租户名称格式错误 | Invalid tenant name format |
| TENANT201002 | 租户描述过长 | Tenant description too long |
| TENANT201003 | 租户配置参数错误 | Invalid tenant configuration parameters |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| TENANT404001 | 租户不存在 | Tenant not found |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| TENANT403001 | 租户已暂停 | Tenant is suspended |
| TENANT403002 | 租户已过期 | Tenant is expired |
| TENANT403003 | 无权访问该租户 | No permission to access this tenant |

### 冲突错误 (409)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| TENANT409001 | 租户名称已存在 | Tenant name already exists |
| TENANT409002 | 租户域名已存在 | Tenant domain already exists |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| TENANT500001 | 租户创建失败 | Failed to create tenant |
| TENANT500002 | 租户更新失败 | Failed to update tenant |
| TENANT500003 | 租户删除失败 | Failed to delete tenant |
| TENANT500004 | 租户配置更新失败 | Failed to update tenant configuration |

---

## 配额模块错误码 (QUOTA)

### 配额不足错误 (402)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| QUOTA402001 | Bot数量配额已用完 | Bot quota exceeded, please upgrade your subscription |
| QUOTA402002 | 消息配额已用完 | Message quota exceeded, please upgrade your subscription |
| QUOTA402003 | 存储配额已用完 | Storage quota exceeded, please upgrade your subscription |
| QUOTA402004 | 工作流配额已用完 | Workflow quota exceeded, please upgrade your subscription |
| QUOTA402005 | Token配额已用完 | Token quota exceeded, please upgrade your subscription |
| QUOTA402006 | API调用配额已用完 | API call quota exceeded, please upgrade your subscription |
| QUOTA402007 | 并发用户配额已用完 | Concurrent user quota exceeded |
| QUOTA402008 | 知识库配额已用完 | Knowledge base quota exceeded |

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| QUOTA201001 | 配额类型错误 | Invalid quota type |
| QUOTA201002 | 配额数量格式错误 | Invalid quota amount format |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| QUOTA500001 | 配额检查失败 | Failed to check quota |
| QUOTA500002 | 配额更新失败 | Failed to update quota |
| QUOTA500003 | 配额使用记录失败 | Failed to record quota usage |

---

## 路由模块错误码 (ROUTING)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| ROUTING201001 | 路由规则格式错误 | Invalid routing rule format |
| ROUTING201002 | 路由配置参数错误 | Invalid routing configuration parameters |
| ROUTING201003 | 意图配置错误 | Invalid intent configuration |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| ROUTING404001 | 路由不存在 | Route not found |
| ROUTING404002 | 路由规则不存在 | Routing rule not found |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| ROUTING500001 | 路由决策失败 | Routing decision failed |
| ROUTING500002 | 意图识别失败 | Intent recognition failed |
| ROUTING500003 | 路由规则加载失败 | Failed to load routing rules |

---

## Bot模块错误码 (BOT)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| BOT201001 | Bot名称格式错误 | Invalid bot name format |
| BOT201002 | Bot描述过长 | Bot description too long |
| BOT201003 | Bot配置参数错误 | Invalid bot configuration parameters |
| BOT201004 | Bot头像格式错误 | Invalid bot avatar format |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| BOT404001 | Bot不存在 | Bot not found |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| BOT403001 | 无权访问Bot | No permission to access bot |
| BOT403002 | Bot已发布 | Bot has been published |
| BOT403003 | Bot已下线 | Bot has been offline |

### 冲突错误 (409)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| BOT409001 | Bot名称已存在 | Bot name already exists |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| BOT500001 | Bot创建失败 | Failed to create bot |
| BOT500002 | Bot更新失败 | Failed to update bot |
| BOT500003 | Bot删除失败 | Failed to delete bot |
| BOT500004 | Bot发布失败 | Failed to publish bot |

---

## 对话模块错误码 (CHAT)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| CHAT201001 | 消息内容为空 | Message content is empty |
| CHAT201002 | 消息过长 | Message too long |
| CHAT201003 | 消息格式错误 | Invalid message format |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| CHAT404001 | 对话不存在 | Conversation not found |
| CHAT404002 | 消息不存在 | Message not found |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| CHAT403001 | 无权访问对话 | No permission to access conversation |
| CHAT403002 | 无权发送消息 | No permission to send message |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| CHAT500001 | 消息发送失败 | Failed to send message |
| CHAT500002 | 对话创建失败 | Failed to create conversation |
| CHAT500003 | 消息流式传输失败 | Failed to stream message |

---

## 订阅模块错误码 (SUBSCRIPTION)

### 配额/付费错误 (402)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| SUBSCRIPTION402001 | 订阅已过期 | Subscription has expired |
| SUBSCRIPTION402002 | 订阅配额已用完 | Subscription quota exhausted |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| SUBSCRIPTION404001 | 订阅不存在 | Subscription not found |

### 冲突错误 (409)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| SUBSCRIPTION409001 | 订阅已存在 | Subscription already exists |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| SUBSCRIPTION500001 | 订阅创建失败 | Failed to create subscription |

---

## 用户模块错误码 (USER)

### 客户端错误 (2xx)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| USER201001 | 用户名格式错误 | Invalid username format |
| USER201002 | 邮箱已存在 | Email already exists |

### 资源错误 (404)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| USER404001 | 用户不存在 | User not found |

### 权限错误 (403)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| USER403001 | 账号已禁用 | Account has been disabled |

### 服务器错误 (500)

| 错误码 | 中文说明 | 英文说明 |
|--------|---------|---------|
| USER500001 | 用户创建失败 | Failed to create user |
| USER500002 | 用户更新失败 | Failed to update user |

---

## 错误响应格式

### 标准错误响应

```json
{
  "code": "AUTH401001",
  "message": "未登录",
  "message_zh": "未登录",
  "message_en": "Not logged in",
  "details": {
    "field": "token",
    "reason": "Token已过期"
  },
  "timestamp": "2025-01-01T12:00:00Z",
  "request_id": "req-1234567890"
}
```

### 成功响应

```json
{
  "code": "SUCCESS",
  "message": "操作成功",
  "data": {
    // 业务数据
  },
  "timestamp": "2025-01-01T12:00:00Z"
}
```

---

## 使用示例

### Go后端使用

```go
package handler

import (
    "your-project/types/errno"
    "github.com/cloudwego/hertz/pkg/app"
)

func Login(ctx context.Context, c *app.RequestContext) {
    // 参数验证
    if req.Username == "" {
        c.JSON(400, map[string]interface{}{
            "code": errno.AUTH201001.Code(),
            "message": errno.AUTH201001.Message(),
            "message_en": errno.AUTH201001.MessageEN(),
        })
        return
    }

    // Token验证
    if token == "" {
        c.JSON(401, map[string]interface{}{
            "code": errno.AUTH401004.Code(),
            "message": errno.AUTH401004.Message(),
        })
        return
    }
}
```

### 错误码助手函数

```go
// 使用error_helper.go中的辅助函数
import "your-project/types/errno"

// 返回错误响应
c.JSON(errno.AUTH401001.HTTPStatus(), errno.AUTH401001.ToResponse())

// 包装错误
err := fmt.Errorf("login failed: %w", errno.AUTH401001)
```

---

## 维护指南

### 新增错误码步骤

1. **确定模块和HTTP状态码**
   - 选择合适的模块标识（如 AUTH、PERMISSION）
   - 选择合适的HTTP状态码（如 401、403）

2. **分配序号**
   - 查看现有错误码，避免重复
   - 使用连续序号（如 001、002、003）

3. **添加错误码定义**
   - 在对应的模块文件中添加错误码
   - 提供中英文双语说明

4. **更新文档**
   - 在 ERROR_CODES.md 中添加错误码说明
   - 确保文档与代码同步

### 错误码命名规范

- **大写字母**：错误码变量名必须大写
- **模块前缀**：必须以模块标识开头
- **HTTP状态**：紧跟HTTP状态码
- **序号**：3位数字序号，不足补零

**示例**:
```go
AUTH401001 = &BaseErrorCode{"AUTH401001", "未登录", "未登录", "Not logged in", 401}
```

---

## 附录

### 相关文档

- [ZKER-统一错误码定义规范.md](../../docs/企业级功能完善与统一性设计方案/ZKER-统一错误码定义规范.md)
- [ZKER-企业级开发规范手册_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

### 更新日志

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.0 | 2025-01-01 | 初始版本，包含300+错误码 |

---

**版权所有 © 2025 Coze Studio Enterprise Team**
