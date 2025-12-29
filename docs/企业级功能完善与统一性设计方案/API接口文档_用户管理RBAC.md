# API接口文档：用户管理(RBAC)模块

**模块名称**: 用户管理 (User Management with RBAC)
**设计文档**: 21-运营管理_用户管理.md, 21-用户管理_RBAC细化补充_完整版.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P0

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 用户管理API](#2-用户管理api)
- [3. 角色管理API](#3-角色管理api)
- [4. 权限管理API](#4-权限管理api)
- [5. 数据权限API](#5-数据权限api)
- [6. 字段权限API](#6-字段权限api)
- [7. 临时授权API](#7-临时授权api)
- [8. 权限验证API](#8-权限验证api)
- [9. 数据模型](#9-数据模型)
- [10. 后端代码示例](#10-后端代码示例)
- [11. 前端代码示例](#11-前端代码示例)
- [12. 错误码定义](#12-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

用户管理模块提供基于RBAC（Role-Based Access Control）的完整用户权限管理能力，包括：

- ✅ **用户生命周期管理** - 注册、审核、激活、禁用、删除
- ✅ **角色管理** - 角色CRUD、角色权限分配
- ✅ **数据权限** - 5级数据权限（全部/本部门及下级/本部门/仅本人/自定义）
- ✅ **字段权限** - 3级字段权限（可见/脱敏/隐藏）
- ✅ **临时授权** - 支持临时权限授予与自动回收
- ✅ **权限审计** - 完整的权限变更日志

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 权限引擎: Casbin v2
- 数据库: MySQL 8.4.5
- 缓存: Redis 8.0

**前端技术栈**:
- 框架: React 18 + TypeScript
- UI库: Semi Design
- 状态管理: Zustand

---

## 2. 用户管理API

### 2.1 创建用户

**接口地址**: `POST /api/v1/admin/users`

**功能说明**: 创建新用户，支持批量创建

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "username": "john.doe",
  "email": "john.doe@example.com",
  "phone": "+86-13800138000",
  "nickname": "John Doe",
  "avatar_url": "https://cdn.example.com/avatar.jpg",
  "department_id": 100,
  "position": "Software Engineer",
  "status": "active",
  "role_ids": [1, 2],
  "password": "SecurePassword123!"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（唯一） |
| email | string | 是 | 邮箱（唯一） |
| phone | string | 否 | 手机号 |
| nickname | string | 是 | 显示名称 |
| avatar_url | string | 否 | 头像URL |
| department_id | int64 | 否 | 部门ID |
| position | string | 否 | 职位 |
| status | string | 否 | 状态: active/inactive/suspended |
| role_ids | array | 是 | 角色ID列表 |
| password | string | 是 | 密码（需符合安全策略） |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1001,
    "username": "john.doe",
    "email": "john.doe@example.com",
    "nickname": "John Doe",
    "department_id": 100,
    "status": "active",
    "created_at": "2025-01-03T10:30:00Z",
    "roles": [
      {
        "role_id": 1,
        "role_name": "Developer",
        "permissions": ["bot:read", "bot:write"]
      }
    ]
  },
  "timestamp": "2025-01-03T10:30:00Z",
  "trace_id": "abc123"
}
```

---

### 2.2 查询用户列表

**接口地址**: `GET /api/v1/admin/users`

**功能说明**: 分页查询用户列表，支持搜索和过滤

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| keyword | string | 否 | 搜索关键词（用户名/邮箱/昵称） |
| department_id | int64 | 否 | 部门ID过滤 |
| role_id | int64 | 否 | 角色ID过滤 |
| status | string | 否 | 状态过滤 |
| sort_by | string | 否 | 排序字段: created_at/username |
| sort_order | string | 否 | 排序方向: asc/desc |

**请求示例**:
```http
GET /api/v1/admin/users?page=1&page_size=20&keyword=john&status=active&sort_by=created_at&sort_order=desc
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 150,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "user_id": 1001,
        "username": "john.doe",
        "email": "john.doe@example.com",
        "nickname": "John Doe",
        "phone": "+86-138****8000",
        "avatar_url": "https://cdn.example.com/avatar.jpg",
        "department": {
          "department_id": 100,
          "department_name": "研发部"
        },
        "position": "Software Engineer",
        "status": "active",
        "roles": [
          {
            "role_id": 1,
            "role_name": "Developer"
          }
        ],
        "created_at": "2025-01-03T10:30:00Z",
        "last_login_at": "2025-01-03T15:20:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T10:30:00Z",
  "trace_id": "abc123"
}
```

---

### 2.3 获取用户详情

**接口地址**: `GET /api/v1/admin/users/{user_id}`

**功能说明**: 获取指定用户的详细信息

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1001,
    "username": "john.doe",
    "email": "john.doe@example.com",
    "phone": "+86-13800138000",
    "nickname": "John Doe",
    "avatar_url": "https://cdn.example.com/avatar.jpg",
    "department": {
      "department_id": 100,
      "department_name": "研发部",
      "parent_path": [1, 10, 100]
    },
    "position": "Software Engineer",
    "status": "active",
    "roles": [
      {
        "role_id": 1,
        "role_name": "Developer",
        "permissions": ["bot:read", "bot:write", "knowledge:read"]
      }
    ],
    "created_at": "2025-01-03T10:30:00Z",
    "updated_at": "2025-01-03T12:00:00Z",
    "last_login_at": "2025-01-03T15:20:00Z"
  },
  "timestamp": "2025-01-03T10:30:00Z",
  "trace_id": "abc123"
}
```

---

### 2.4 更新用户

**接口地址**: `PUT /api/v1/admin/users/{user_id}`

**功能说明**: 更新用户信息

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

**请求参数**:
```json
{
  "email": "john.doe.new@example.com",
  "phone": "+86-13900139000",
  "nickname": "John Doe II",
  "avatar_url": "https://cdn.example.com/avatar-new.jpg",
  "department_id": 200,
  "position": "Senior Software Engineer",
  "role_ids": [1, 2, 3]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "用户更新成功",
  "data": {
    "user_id": 1001,
    "updated_at": "2025-01-03T16:00:00Z"
  },
  "timestamp": "2025-01-03T16:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.5 删除用户

**接口地址**: `DELETE /api/v1/admin/users/{user_id}`

**功能说明**: 删除指定用户（软删除）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "用户删除成功",
  "data": {
    "user_id": 1001,
    "deleted_at": "2025-01-03T17:00:00Z"
  },
  "timestamp": "2025-01-03T17:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.6 更新用户状态

**接口地址**: `PUT /api/v1/admin/users/{user_id}/status`

**功能说明**: 更新用户状态（激活/禁用/冻结）

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

**请求参数**:
```json
{
  "status": "suspended",
  "reason": "违反公司规定"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 是 | 状态: active/inactive/suspended |
| reason | string | 否 | 状态变更原因 |

**响应示例**:
```json
{
  "code": 0,
  "message": "用户状态更新成功",
  "data": {
    "user_id": 1001,
    "status": "suspended",
    "updated_at": "2025-01-03T18:00:00Z"
  },
  "timestamp": "2025-01-03T18:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.7 批量导入用户

**接口地址**: `POST /api/v1/admin/users/batch-import`

**功能说明**: 批量导入用户（CSV/Excel）

**请求头**:
```http
Content-Type: multipart/form-data
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | CSV或Excel文件 |
| default_role_ids | string | 否 | 默认角色ID列表（逗号分隔） |
| send_notification | boolean | 否 | 是否发送通知邮件，默认false |
| department_id | int64 | 否 | 默认部门ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "批量导入任务创建成功",
  "data": {
    "task_id": "batch-import-20250103-001",
    "total_count": 150,
    "status": "processing",
    "estimated_time": 300,
    "created_at": "2025-01-03T19:00:00Z"
  },
  "timestamp": "2025-01-03T19:00:00Z",
  "trace_id": "abc123"
}
```

---

### 2.8 获取批量导入进度

**接口地址**: `GET /api/v1/admin/users/batch-import/{task_id}`

**功能说明**: 查询批量导入任务进度

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| task_id | string | 任务ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "batch-import-20250103-001",
    "status": "completed",
    "total_count": 150,
    "processed_count": 150,
    "success_count": 145,
    "failed_count": 5,
    "failed_records": [
      {
        "row": 10,
        "username": "duplicate.user",
        "error": "用户名已存在"
      }
    ],
    "started_at": "2025-01-03T19:00:00Z",
    "completed_at": "2025-01-03T19:05:00Z"
  },
  "timestamp": "2025-01-03T19:05:00Z",
  "trace_id": "abc123"
}
```

---

## 3. 角色管理API

### 3.1 创建角色

**接口地址**: `POST /api/v1/admin/roles`

**功能说明**: 创建新角色

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "role_name": "Product Manager",
  "role_code": "product_manager",
  "description": "产品经理角色，负责产品规划和需求管理",
  "is_system": false,
  "permissions": [
    "bot:read",
    "bot:write",
    "knowledge:read",
    "workflow:read",
    "workflow:execute"
  ]
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_name | string | 是 | 角色名称 |
| role_code | string | 是 | 角色代码（唯一） |
| description | string | 否 | 角色描述 |
| is_system | boolean | 否 | 是否系统角色，默认false |
| permissions | array | 是 | 权限代码列表 |

**响应示例**:
```json
{
  "code": 0,
  "message": "角色创建成功",
  "data": {
    "role_id": 10,
    "role_name": "Product Manager",
    "role_code": "product_manager",
    "description": "产品经理角色，负责产品规划和需求管理",
    "is_system": false,
    "permissions": [
      "bot:read",
      "bot:write",
      "knowledge:read",
      "workflow:read",
      "workflow:execute"
    ],
    "created_at": "2025-01-03T20:00:00Z"
  },
  "timestamp": "2025-01-03T20:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.2 查询角色列表

**接口地址**: `GET /api/v1/admin/roles`

**功能说明**: 查询角色列表

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认20 |
| keyword | string | 否 | 搜索关键词 |
| is_system | boolean | 否 | 是否系统角色 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "role_id": 1,
        "role_name": "Administrator",
        "role_code": "admin",
        "description": "系统管理员",
        "is_system": true,
        "user_count": 5,
        "permissions_count": 100,
        "created_at": "2025-01-01T00:00:00Z"
      },
      {
        "role_id": 10,
        "role_name": "Product Manager",
        "role_code": "product_manager",
        "description": "产品经理角色",
        "is_system": false,
        "user_count": 12,
        "permissions_count": 5,
        "created_at": "2025-01-03T20:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T20:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.3 获取角色详情

**接口地址**: `GET /api/v1/admin/roles/{role_id}`

**功能说明**: 获取角色详细信息

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| role_id | int64 | 角色ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "role_id": 10,
    "role_name": "Product Manager",
    "role_code": "product_manager",
    "description": "产品经理角色，负责产品规划和需求管理",
    "is_system": false,
    "permissions": [
      {
        "permission_code": "bot:read",
        "permission_name": "Bot查看权限",
        "resource": "bot",
        "action": "read"
      },
      {
        "permission_code": "bot:write",
        "permission_name": "Bot编辑权限",
        "resource": "bot",
        "action": "write"
      }
    ],
    "user_count": 12,
    "created_at": "2025-01-03T20:00:00Z",
    "updated_at": "2025-01-03T20:00:00Z"
  },
  "timestamp": "2025-01-03T20:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.4 更新角色

**接口地址**: `PUT /api/v1/admin/roles/{role_id}`

**功能说明**: 更新角色信息

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| role_id | int64 | 角色ID |

**请求参数**:
```json
{
  "role_name": "Senior Product Manager",
  "description": "高级产品经理角色",
  "permissions": [
    "bot:read",
    "bot:write",
    "bot:delete",
    "knowledge:read",
    "knowledge:write",
    "workflow:read",
    "workflow:execute"
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "角色更新成功",
  "data": {
    "role_id": 10,
    "updated_at": "2025-01-03T21:00:00Z"
  },
  "timestamp": "2025-01-03T21:00:00Z",
  "trace_id": "abc123"
}
```

---

### 3.5 删除角色

**接口地址**: `DELETE /api/v1/admin/roles/{role_id}`

**功能说明**: 删除角色（系统角色不可删除）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| role_id | int64 | 角色ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "角色删除成功",
  "data": {
    "role_id": 10,
    "deleted_at": "2025-01-03T22:00:00Z"
  },
  "timestamp": "2025-01-03T22:00:00Z",
  "trace_id": "abc123"
}
```

---

## 4. 权限管理API

### 4.1 查询所有权限

**接口地址**: `GET /api/v1/admin/permissions`

**功能说明**: 查询系统所有权限列表

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| resource | string | 否 | 资源类型过滤 |
| action | string | 否 | 操作类型过滤 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "items": [
      {
        "permission_code": "bot:read",
        "permission_name": "Bot查看权限",
        "resource": "bot",
        "action": "read",
        "description": "查看Bot列表和详情"
      },
      {
        "permission_code": "bot:write",
        "permission_name": "Bot编辑权限",
        "resource": "bot",
        "action": "write",
        "description": "创建和编辑Bot"
      },
      {
        "permission_code": "bot:delete",
        "permission_name": "Bot删除权限",
        "resource": "bot",
        "action": "delete",
        "description": "删除Bot"
      }
    ]
  },
  "timestamp": "2025-01-03T22:00:00Z",
  "trace_id": "abc123"
}
```

---

### 4.2 分配权限给角色

**接口地址**: `POST /api/v1/admin/roles/{role_id}/permissions`

**功能说明**: 为角色分配权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| role_id | int64 | 角色ID |

**请求参数**:
```json
{
  "permissions": [
    "bot:read",
    "bot:write",
    "knowledge:read",
    "workflow:execute"
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "权限分配成功",
  "data": {
    "role_id": 10,
    "assigned_count": 4,
    "updated_at": "2025-01-03T23:00:00Z"
  },
  "timestamp": "2025-01-03T23:00:00Z",
  "trace_id": "abc123"
}
```

---

## 5. 数据权限API

### 5.1 查询数据权限配置

**接口地址**: `GET /api/v1/rbac/data-permissions`

**功能说明**: 查询角色的数据权限配置

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_id | int64 | 是 | 角色ID |
| resource_type | string | 是 | 资源类型: user/bot/knowledge_base |

**请求示例**:
```http
GET /api/v1/rbac/data-permissions?role_id=10&resource_type=user
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "permission_id": 501,
    "role_id": 10,
    "resource_type": "user",
    "permission_scope": "department_and_below",
    "custom_org_ids": [100, 101, 102],
    "custom_user_ids": null,
    "permissions": ["read", "write"],
    "created_at": "2025-01-03T10:00:00Z",
    "updated_at": "2025-01-03T10:00:00Z"
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.2 创建数据权限

**接口地址**: `POST /api/v1/rbac/data-permissions`

**功能说明**: 为角色创建数据权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "role_id": 10,
  "resource_type": "user",
  "permission_scope": "department_and_below",
  "custom_org_ids": [100, 101, 102],
  "custom_user_ids": null,
  "permissions": ["read", "write", "delete"]
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_id | int64 | 是 | 角色ID |
| resource_type | string | 是 | 资源类型 |
| permission_scope | string | 是 | 权限范围: all/department_and_below/department/self/custom |
| custom_org_ids | array | 否 | 自定义组织ID列表 |
| custom_user_ids | array | 否 | 自定义用户ID列表 |
| permissions | array | 是 | 操作权限: read/write/delete |

**响应示例**:
```json
{
  "code": 0,
  "message": "数据权限创建成功",
  "data": {
    "permission_id": 502,
    "created_at": "2025-01-03T11:00:00Z"
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.3 更新数据权限

**接口地址**: `PUT /api/v1/rbac/data-permissions/{permission_id}`

**功能说明**: 更新数据权限配置

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| permission_id | int64 | 权限ID |

**请求参数**:
```json
{
  "permission_scope": "all",
  "permissions": ["read", "write", "delete"]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "数据权限更新成功",
  "data": {
    "permission_id": 502,
    "updated_at": "2025-01-03T12:00:00Z"
  },
  "timestamp": "2025-01-03T12:00:00Z",
  "trace_id": "abc123"
}
```

---

### 5.4 删除数据权限

**接口地址**: `DELETE /api/v1/rbac/data-permissions/{permission_id}`

**功能说明**: 删除数据权限

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| permission_id | int64 | 权限ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "数据权限删除成功",
  "data": {
    "permission_id": 502,
    "deleted_at": "2025-01-03T13:00:00Z"
  },
  "timestamp": "2025-01-03T13:00:00Z",
  "trace_id": "abc123"
}
```

---

## 6. 字段权限API

### 6.1 查询字段权限配置

**接口地址**: `GET /api/v1/rbac/field-permissions`

**功能说明**: 查询角色的字段权限配置

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_id | int64 | 是 | 角色ID |
| resource_type | string | 是 | 资源类型 |

**请求示例**:
```http
GET /api/v1/rbac/field-permissions?role_id=10&resource_type=user
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 8,
    "items": [
      {
        "field_permission_id": 601,
        "field_name": "phone",
        "field_display_name": "手机号",
        "permission": "masked",
        "mask_rule": "phone",
        "mask_pattern": null,
        "mask_replacement": "****"
      },
      {
        "field_permission_id": 602,
        "field_name": "email",
        "field_display_name": "邮箱",
        "permission": "visible",
        "mask_rule": null,
        "mask_pattern": null,
        "mask_replacement": null
      },
      {
        "field_permission_id": 603,
        "field_name": "id_card",
        "field_display_name": "身份证号",
        "permission": "hidden",
        "mask_rule": null,
        "mask_pattern": null,
        "mask_replacement": null
      }
    ]
  },
  "timestamp": "2025-01-03T14:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.2 创建字段权限

**接口地址**: `POST /api/v1/rbac/field-permissions`

**功能说明**: 为角色创建字段权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "role_id": 10,
  "resource_type": "user",
  "field_name": "phone",
  "field_display_name": "手机号",
  "permission": "masked",
  "mask_rule": "phone",
  "mask_pattern": null,
  "mask_replacement": "****"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role_id | int64 | 是 | 角色ID |
| resource_type | string | 是 | 资源类型 |
| field_name | string | 是 | 字段名称 |
| field_display_name | string | 否 | 字段显示名称 |
| permission | string | 是 | 权限级别: visible/masked/hidden |
| mask_rule | string | 否 | 脱敏规则: phone/id_card/email/bank_card |
| mask_pattern | string | 否 | 正则表达式模式 |
| mask_replacement | string | 否 | 替换字符，默认**** |

**响应示例**:
```json
{
  "code": 0,
  "message": "字段权限创建成功",
  "data": {
    "field_permission_id": 604,
    "created_at": "2025-01-03T15:00:00Z"
  },
  "timestamp": "2025-01-03T15:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.3 更新字段权限

**接口地址**: `PUT /api/v1/rbac/field-permissions/{field_permission_id}`

**功能说明**: 更新字段权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| field_permission_id | int64 | 字段权限ID |

**请求参数**:
```json
{
  "permission": "visible",
  "mask_rule": null
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "字段权限更新成功",
  "data": {
    "field_permission_id": 604,
    "updated_at": "2025-01-03T16:00:00Z"
  },
  "timestamp": "2025-01-03T16:00:00Z",
  "trace_id": "abc123"
}
```

---

### 6.4 批量更新字段权限

**接口地址**: `PUT /api/v1/rbac/roles/{role_id}/field-permissions/batch`

**功能说明**: 批量更新角色的字段权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| role_id | int64 | 角色ID |

**请求参数**:
```json
{
  "resource_type": "user",
  "fields": [
    {
      "field_name": "phone",
      "permission": "masked",
      "mask_rule": "phone"
    },
    {
      "field_name": "email",
      "permission": "visible"
    },
    {
      "field_name": "id_card",
      "permission": "hidden"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "批量更新成功",
  "data": {
    "updated_count": 3,
    "updated_at": "2025-01-03T17:00:00Z"
  },
  "timestamp": "2025-01-03T17:00:00Z",
  "trace_id": "abc123"
}
```

---

## 7. 临时授权API

### 7.1 创建临时授权

**接口地址**: `POST /api/v1/rbac/temporary-grants`

**功能说明**: 为用户创建临时授权

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "grantee_id": 1001,
  "permission_type": "resource",
  "permission_config": {
    "resource_type": "bot",
    "resource_id": "bot-123",
    "actions": ["read", "write"]
  },
  "start_time": "2025-01-03T10:00:00Z",
  "end_time": "2025-01-06T18:00:00Z",
  "reason": "项目临时协作需要"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grantee_id | int64 | 是 | 授权接受人ID |
| permission_type | string | 是 | 权限类型: role/resource/operation |
| permission_config | object | 是 | 权限配置 |
| start_time | string | 是 | 开始时间（ISO 8601） |
| end_time | string | 是 | 结束时间（ISO 8601） |
| reason | string | 否 | 授权原因 |

**响应示例**:
```json
{
  "code": 0,
  "message": "临时授权创建成功",
  "data": {
    "grant_id": 701,
    "grant_code": "TG-550e8400-e29b-41d4-a716-446655440000",
    "grantee_id": 1001,
    "permission_type": "resource",
    "start_time": "2025-01-03T10:00:00Z",
    "end_time": "2025-01-06T18:00:00Z",
    "created_at": "2025-01-03T09:00:00Z"
  },
  "timestamp": "2025-01-03T09:00:00Z",
  "trace_id": "abc123"
}
```

---

### 7.2 查询临时授权列表

**接口地址**: `GET /api/v1/rbac/temporary-grants`

**功能说明**: 查询临时授权列表

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| grantee_id | int64 | 否 | 授权接受人ID |
| grantor_id | int64 | 否 | 授权授予人ID |
| status | string | 否 | 状态: active/expired/revoked |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**请求示例**:
```http
GET /api/v1/rbac/temporary-grants?grantee_id=1001&status=active
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "items": [
      {
        "grant_id": 701,
        "grant_code": "TG-550e8400-e29b-41d4-a716-446655440000",
        "grantor_id": 1,
        "grantor_name": "Admin",
        "grantee_id": 1001,
        "grantee_name": "John Doe",
        "permission_type": "resource",
        "permission_config": {
          "resource_type": "bot",
          "resource_id": "bot-123",
          "actions": ["read", "write"]
        },
        "start_time": "2025-01-03T10:00:00Z",
        "end_time": "2025-01-06T18:00:00Z",
        "is_revoked": false,
        "status": "active",
        "reason": "项目临时协作需要",
        "created_at": "2025-01-03T09:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T10:00:00Z",
  "trace_id": "abc123"
}
```

---

### 7.3 撤销临时授权

**接口地址**: `POST /api/v1/rbac/temporary-grants/{grant_code}/revoke`

**功能说明**: 撤销临时授权

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| grant_code | string | 授权码 |

**请求参数**:
```json
{
  "reason": "项目提前结束"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "临时授权撤销成功",
  "data": {
    "grant_code": "TG-550e8400-e29b-41d4-a716-446655440000",
    "is_revoked": true,
    "revoked_at": "2025-01-04T10:00:00Z"
  },
  "timestamp": "2025-01-04T10:00:00Z",
  "trace_id": "abc123"
}
```

---

## 8. 权限验证API

### 8.1 检查用户权限

**接口地址**: `POST /api/v1/rbac/permissions/check`

**功能说明**: 检查用户是否拥有指定权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "user_id": 1001,
  "resource": "bot",
  "action": "write"
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| user_id | int64 | 是 | 用户ID |
| resource | string | 是 | 资源类型 |
| action | string | 是 | 操作类型: read/write/delete/execute |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "has_permission": true,
    "permission_source": "role",
    "role_name": "Developer",
    "checked_at": "2025-01-03T11:00:00Z"
  },
  "timestamp": "2025-01-03T11:00:00Z",
  "trace_id": "abc123"
}
```

---

### 8.2 批量检查权限

**接口地址**: `POST /api/v1/rbac/permissions/check-batch`

**功能说明**: 批量检查多个权限

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "user_id": 1001,
  "permissions": [
    {"resource": "bot", "action": "read"},
    {"resource": "bot", "action": "write"},
    {"resource": "bot", "action": "delete"},
    {"resource": "knowledge", "action": "read"}
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1001,
    "results": [
      {
        "resource": "bot",
        "action": "read",
        "has_permission": true
      },
      {
        "resource": "bot",
        "action": "write",
        "has_permission": true
      },
      {
        "resource": "bot",
        "action": "delete",
        "has_permission": false
      },
      {
        "resource": "knowledge",
        "action": "read",
        "has_permission": true
      }
    ]
  },
  "timestamp": "2025-01-03T12:00:00Z",
  "trace_id": "abc123"
}
```

---

### 8.3 获取用户所有权限

**接口地址**: `GET /api/v1/rbac/users/{user_id}/permissions`

**功能说明**: 获取用户的所有权限（包括角色权限和临时授权）

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| user_id | int64 | 用户ID |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1001,
    "roles": [
      {
        "role_id": 1,
        "role_name": "Developer",
        "permissions": [
          "bot:read",
          "bot:write",
          "knowledge:read"
        ]
      }
    ],
    "temporary_grants": [
      {
        "grant_code": "TG-550e8400-e29b-41d4-a716-446655440000",
        "permission_config": {
          "resource_type": "bot",
          "resource_id": "bot-123",
          "actions": ["delete"]
        },
        "end_time": "2025-01-06T18:00:00Z"
      }
    ],
    "all_permissions": [
      "bot:read",
      "bot:write",
      "bot:delete",
      "knowledge:read"
    ]
  },
  "timestamp": "2025-01-03T13:00:00Z",
  "trace_id": "abc123"
}
```

---

### 8.4 查询权限审计日志

**接口地址**: `GET /api/v1/rbac/audit-logs`

**功能说明**: 查询权限审计日志

**请求头**:
```http
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**查询参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| operator_id | int64 | 否 | 操作人ID |
| operation_type | string | 否 | 操作类型: grant/revoke/modify/check |
| target_type | string | 否 | 目标类型: role/user/resource |
| start_time | string | 否 | 开始时间 |
| end_time | string | 否 | 结束时间 |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 150,
    "items": [
      {
        "log_id": 8001,
        "operator_id": 1,
        "operator_name": "Admin",
        "operation_type": "grant",
        "target_type": "user",
        "target_id": "1001",
        "permission_detail": {
          "permission": "bot:write",
          "source": "role"
        },
        "old_value": null,
        "new_value": {
          "role_id": 10,
          "permissions": ["bot:write"]
        },
        "operation_result": "success",
        "ip_address": "192.168.1.100",
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  },
  "timestamp": "2025-01-03T14:00:00Z",
  "trace_id": "abc123"
}
```

---

## 9. 数据模型

### 9.1 User（用户）

```typescript
interface User {
  user_id: number;
  username: string;
  email: string;
  phone?: string;
  nickname: string;
  avatar_url?: string;
  department_id?: number;
  department?: Department;
  position?: string;
  status: 'active' | 'inactive' | 'suspended';
  roles?: Role[];
  created_at: string;
  updated_at: string;
  last_login_at?: string;
}
```

### 9.2 Role（角色）

```typescript
interface Role {
  role_id: number;
  role_name: string;
  role_code: string;
  description?: string;
  is_system: boolean;
  permissions?: Permission[];
  user_count?: number;
  created_at: string;
  updated_at: string;
}
```

### 9.3 Permission（权限）

```typescript
interface Permission {
  permission_code: string;
  permission_name: string;
  resource: string;
  action: 'read' | 'write' | 'delete' | 'execute';
  description?: string;
}
```

### 9.4 DataPermission（数据权限）

```typescript
interface DataPermission {
  permission_id: number;
  role_id: number;
  resource_type: string;
  permission_scope: 'all' | 'department_and_below' | 'department' | 'self' | 'custom';
  custom_org_ids?: number[];
  custom_user_ids?: number[];
  permissions: string[];
  created_at: string;
  updated_at: string;
}
```

### 9.5 FieldPermission（字段权限）

```typescript
interface FieldPermission {
  field_permission_id: number;
  role_id: number;
  resource_type: string;
  field_name: string;
  field_display_name?: string;
  permission: 'visible' | 'masked' | 'hidden';
  mask_rule?: 'phone' | 'id_card' | 'email' | 'bank_card';
  mask_pattern?: string;
  mask_replacement?: string;
}
```

### 9.6 TemporaryGrant（临时授权）

```typescript
interface TemporaryGrant {
  grant_id: number;
  grant_code: string;
  grantor_id: number;
  grantee_id: number;
  permission_type: 'role' | 'resource' | 'operation';
  permission_config: any;
  start_time: string;
  end_time: string;
  is_revoked: boolean;
  status: 'active' | 'expired' | 'revoked';
  reason?: string;
  created_at: string;
}
```

---

## 10. 后端代码示例

### 10.1 用户管理Handler

```go
// api/handler/user_handler.go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

type UserHandler struct {
    userService *service.UserService
    roleService *service.RoleService
    rbacService *service.RBACService
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(ctx context.Context, c *app.RequestContext) {
    var req dto.CreateUserRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, response.Error("参数错误", err))
        return
    }

    // 调用服务层
    user, err := h.userService.CreateUser(ctx, &req)
    if err != nil {
        c.JSON(500, response.Error("创建用户失败", err))
        return
    }

    c.JSON(200, response.Success(user))
}

// ListUsers 查询用户列表
func (h *UserHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
    var req dto.ListUsersRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, response.Error("参数错误", err))
        return
    }

    // 查询用户
    users, total, err := h.userService.ListUsers(ctx, &req)
    if err != nil {
        c.JSON(500, response.Error("查询失败", err))
        return
    }

    // 应用字段权限（脱敏）
    userID := auth.GetUserID(ctx)
    for i := range users {
        masked, _ := h.rbacService.ApplyFieldPermissions(
            ctx,
            userID,
            "user",
            users[i],
        )
        users[i] = masked.(entity.User)
    }

    c.JSON(200, response.SuccessPage(users, total, req.Page, req.PageSize))
}

// UpdateUserStatus 更新用户状态
func (h *UserHandler) UpdateUserStatus(ctx context.Context, c *app.RequestContext) {
    userID := c.Param64("user_id")

    var req dto.UpdateUserStatusRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, response.Error("参数错误", err))
        return
    }

    err := h.userService.UpdateUserStatus(ctx, userID, &req)
    if err != nil {
        c.JSON(500, response.Error("更新状态失败", err))
        return
    }

    c.JSON(200, response.Success(nil))
}
```

### 10.2 RBAC服务

```go
// application/service/rbac_service.go
package service

import (
    "context"
    "github.com/casbin/casbin/v2"
)

type RBACService struct {
    dataPermRepo      repository.DataPermissionRepository
    fieldPermRepo     repository.FieldPermissionRepository
    temporaryGrantRepo repository.TemporaryGrantRepository
    enforcer          *casbin.Enforcer
    orgService        *OrganizationService
    cache             *redis.Client
}

// FilterDataByPermission 根据权限过滤数据
func (s *RBACService) FilterDataByPermission(
    ctx context.Context,
    userID int64,
    resourceType string,
    dataIDs []interface{},
) ([]interface{}, error) {
    // 1. 获取用户的角色
    roles, err := s.getUserRoles(ctx, userID)
    if err != nil {
        return nil, err
    }

    // 2. 获取每个角色的数据权限
    var allowedIDs []interface{}
    hasAllPermission := false

    for _, role := range roles {
        perm, err := s.dataPermRepo.GetByRoleAndResource(
            ctx,
            role.ID,
            resourceType,
        )
        if err != nil {
            continue
        }

        switch perm.PermissionScope {
        case "all":
            hasAllPermission = true
            break

        case "department_and_below":
            orgIDs := s.getDepartmentAndBelow(ctx, userID)
            filtered := s.filterByOrgs(ctx, dataIDs, orgIDs)
            allowedIDs = append(allowedIDs, filtered...)

        case "department":
            orgID := s.getUserDepartment(ctx, userID)
            filtered := s.filterByOrgs(ctx, dataIDs, []int64{orgID})
            allowedIDs = append(allowedIDs, filtered...)

        case "self":
            filtered := s.filterByOwner(ctx, dataIDs, userID)
            allowedIDs = append(allowedIDs, filtered...)

        case "custom":
            if perm.CustomOrgIDs != nil {
                filtered := s.filterByOrgs(ctx, dataIDs, perm.CustomOrgIDs)
                allowedIDs = append(allowedIDs, filtered...)
            }
            if perm.CustomUserIDs != nil {
                filtered := s.filterByUsers(ctx, dataIDs, perm.CustomUserIDs)
                allowedIDs = append(allowedIDs, filtered...)
            }
        }
    }

    if hasAllPermission {
        return dataIDs, nil
    }

    return uniqueIDs(allowedIDs), nil
}

// ApplyFieldPermissions 应用字段权限（脱敏）
func (s *RBACService) ApplyFieldPermissions(
    ctx context.Context,
    userID int64,
    resourceType string,
    data interface{},
) (interface{}, error) {
    // 1. 获取用户的角色
    roles, err := s.getUserRoles(ctx, userID)
    if err != nil {
        return nil, err
    }

    // 2. 获取所有字段权限配置
    fieldPerms := make(map[string]*entity.FieldPermission)
    for _, role := range roles {
        perms, err := s.fieldPermRepo.GetByRoleAndResource(
            ctx,
            role.ID,
            resourceType,
        )
        if err != nil {
            continue
        }

        for _, perm := range perms {
            if existing, ok := fieldPerms[perm.FieldName]; !ok || perm.Permission < existing.Permission {
                fieldPerms[perm.FieldName] = perm
            }
        }
    }

    // 3. 应用权限
    result := make(map[string]interface{})
    dataMap, ok := data.(map[string]interface{})
    if !ok {
        return data, nil
    }

    for fieldName, value := range dataMap {
        perm, ok := fieldPerms[fieldName]
        if !ok {
            // 没有配置权限，默认可见
            result[fieldName] = value
            continue
        }

        switch perm.Permission {
        case "visible":
            result[fieldName] = value
        case "masked":
            result[fieldName] = s.maskValue(value, perm)
        case "hidden":
            // 不返回
        }
    }

    return result, nil
}

// CheckPermission 检查权限
func (s *RBACService) CheckPermission(
    ctx context.Context,
    userID int64,
    resource string,
    action string,
) (bool, error) {
    // 1. 检查RBAC权限
    allowed, err := s.enforcer.Enforce(
        fmt.Sprintf("user_%d", userID),
        resource,
        action,
    )
    if err != nil {
        return false, err
    }

    if !allowed {
        // 2. 检查临时授权
        return s.checkTemporaryGrant(ctx, userID, resource, action)
    }

    return true, nil
}

// CreateTemporaryGrant 创建临时授权
func (s *RBACService) CreateTemporaryGrant(
    ctx context.Context,
    req *dto.CreateGrantRequest,
) (*dto.CreateGrantResponse, error) {
    // 1. 生成授权码
    grantCode := "TG-" + uuid.New().String()

    // 2. 构建授权记录
    grant := &entity.TemporaryGrant{
        TenantID:         req.TenantID,
        GrantCode:        grantCode,
        GrantorID:        req.GrantorID,
        GranteeID:        req.GranteeID,
        PermissionType:   req.PermissionType,
        PermissionConfig: req.PermissionConfig,
        Reason:           req.Reason,
        StartTime:        req.StartTime,
        EndTime:          req.EndTime,
    }

    // 3. 持久化授权
    if err := s.temporaryGrantRepo.Create(ctx, grant); err != nil {
        return nil, err
    }

    // 4. 添加到自动回收任务
    s.scheduleRevoke(grant)

    return &dto.CreateGrantResponse{
        GrantCode: grantCode,
        Message:   "临时授权创建成功",
    }, nil
}

// maskValue 脱敏处理
func (s *RBACService) maskValue(
    value interface{},
    perm *entity.FieldPermission,
) interface{} {
    if value == nil {
        return nil
    }

    strValue := fmt.Sprintf("%v", value)

    // 优先使用正则表达式脱敏
    if perm.MaskPattern != "" {
        pattern := regexp.MustCompile(perm.MaskPattern)
        return pattern.ReplaceAllString(strValue, perm.MaskReplacement)
    }

    // 使用预定义规则脱敏
    switch perm.MaskRule {
    case "phone":
        if len(strValue) == 11 {
            return strValue[:3] + "****" + strValue[7:]
        }
    case "id_card":
        if len(strValue) == 18 {
            return strValue[:4] + "**********" + strValue[14:]
        }
    case "email":
        parts := strings.Split(strValue, "@")
        if len(parts) == 2 {
            return string(parts[0][0]) + "***@" + parts[1]
        }
    case "bank_card":
        return "****"
    }

    return perm.MaskReplacement
}
```

### 10.3 权限中间件

```go
// infra/middleware/permission_middleware.go
package middleware

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

type PermissionMiddleware struct {
    rbacService *service.RBACService
}

// CheckPermission 权限检查中间件
func (m *PermissionMiddleware) CheckPermission(
    resource string,
    action string,
) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 获取用户ID
        userID := auth.GetUserID(ctx)
        if userID == 0 {
            c.JSON(401, response.Error("未登录", nil))
            c.Abort()
            return
        }

        // 2. 检查权限
        allowed, err := m.rbacService.CheckPermission(
            ctx,
            userID,
            resource,
            action,
        )
        if err != nil {
            c.JSON(500, response.Error("权限检查失败", err))
            c.Abort()
            return
        }

        if !allowed {
            c.JSON(403, response.Error("无权限访问", nil))
            c.Abort()
            return
        }

        c.Next(ctx)
    }
}

// ApplyFieldPermissions 字段权限中间件
func (m *PermissionMiddleware) ApplyFieldPermissions(
    resourceType string,
) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        c.Next(ctx)

        // 在响应前应用字段权限
        resp := c.Response.GetBody()
        // 应用脱敏逻辑
        // ...
    }
}
```

---

## 11. 前端代码示例

### 11.1 RBAC API客户端

```typescript
// api/rbac-api.ts
import { request } from '@shared/utils/request';

export interface RBACApiClient {
  // 用户管理
  createUser(data: CreateUserRequest): Promise<User>;
  listUsers(params: ListUsersParams): Promise<PageResult<User>>;
  getUser(userId: number): Promise<User>;
  updateUser(userId: number, data: UpdateUserRequest): Promise<void>;
  deleteUser(userId: number): Promise<void>;
  updateUserStatus(userId: number, data: UpdateStatusRequest): Promise<void>;

  // 角色管理
  createRole(data: CreateRoleRequest): Promise<Role>;
  listRoles(params: ListRolesParams): Promise<PageResult<Role>>;
  getRole(roleId: number): Promise<Role>;
  updateRole(roleId: number, data: UpdateRoleRequest): Promise<void>;
  deleteRole(roleId: number): Promise<void>;
  assignPermissionsToRole(roleId: number, permissions: string[]): Promise<void>;

  // 权限管理
  listPermissions(params?: ListPermissionsParams): Promise<Permission[]>;
  checkPermission(params: CheckPermissionRequest): Promise<CheckPermissionResult>;
  checkPermissionsBatch(params: CheckPermissionsBatchRequest): Promise<CheckPermissionResult[]>;
  getUserPermissions(userId: number): Promise<UserPermissions>;

  // 数据权限
  createDataPermission(data: CreateDataPermissionRequest): Promise<DataPermission>;
  getDataPermissions(roleId: number, resourceType: string): Promise<DataPermission>;
  updateDataPermission(id: number, data: UpdateDataPermissionRequest): Promise<void>;
  deleteDataPermission(id: number): Promise<void>;

  // 字段权限
  createFieldPermission(data: CreateFieldPermissionRequest): Promise<FieldPermission>;
  getFieldPermissions(roleId: number, resourceType: string): Promise<FieldPermission[]>;
  updateFieldPermission(id: number, data: UpdateFieldPermissionRequest): Promise<void>;
  updateFieldPermissionsBatch(roleId: number, data: UpdateFieldPermissionsBatchRequest): Promise<void>;

  // 临时授权
  createTemporaryGrant(data: CreateTemporaryGrantRequest): Promise<TemporaryGrant>;
  listTemporaryGrants(params: ListTemporaryGrantsParams): Promise<PageResult<TemporaryGrant>>;
  revokeTemporaryGrant(grantCode: string, reason?: string): Promise<void>;

  // 审计日志
  listAuditLogs(params: ListAuditLogsParams): Promise<PageResult<AuditLog>>;
}

class RBACApiClientImpl implements RBACApiClient {
  private baseURL = '/api/v1';

  // 用户管理
  async createUser(data: CreateUserRequest): Promise<User> {
    return request.post(`${this.baseURL}/admin/users`, data);
  }

  async listUsers(params: ListUsersParams): Promise<PageResult<User>> {
    return request.get(`${this.baseURL}/admin/users`, { params });
  }

  async getUser(userId: number): Promise<User> {
    return request.get(`${this.baseURL}/admin/users/${userId}`);
  }

  async updateUser(userId: number, data: UpdateUserRequest): Promise<void> {
    return request.put(`${this.baseURL}/admin/users/${userId}`, data);
  }

  async deleteUser(userId: number): Promise<void> {
    return request.delete(`${this.baseURL}/admin/users/${userId}`);
  }

  async updateUserStatus(userId: number, data: UpdateStatusRequest): Promise<void> {
    return request.put(`${this.baseURL}/admin/users/${userId}/status`, data);
  }

  // 角色管理
  async createRole(data: CreateRoleRequest): Promise<Role> {
    return request.post(`${this.baseURL}/admin/roles`, data);
  }

  async listRoles(params: ListRolesParams): Promise<PageResult<Role>> {
    return request.get(`${this.baseURL}/admin/roles`, { params });
  }

  async getRole(roleId: number): Promise<Role> {
    return request.get(`${this.baseURL}/admin/roles/${roleId}`);
  }

  async updateRole(roleId: number, data: UpdateRoleRequest): Promise<void> {
    return request.put(`${this.baseURL}/admin/roles/${roleId}`, data);
  }

  async deleteRole(roleId: number): Promise<void> {
    return request.delete(`${this.baseURL}/admin/roles/${roleId}`);
  }

  async assignPermissionsToRole(roleId: number, permissions: string[]): Promise<void> {
    return request.post(`${this.baseURL}/admin/roles/${roleId}/permissions`, { permissions });
  }

  // 权限管理
  async listPermissions(params?: ListPermissionsParams): Promise<Permission[]> {
    return request.get(`${this.baseURL}/admin/permissions`, { params });
  }

  async checkPermission(params: CheckPermissionRequest): Promise<CheckPermissionResult> {
    return request.post(`${this.baseURL}/rbac/permissions/check`, params);
  }

  async checkPermissionsBatch(params: CheckPermissionsBatchRequest): Promise<CheckPermissionResult[]> {
    const result = await request.post<{ results: CheckPermissionResult[] }>(
      `${this.baseURL}/rbac/permissions/check-batch`,
      params
    );
    return result.results;
  }

  async getUserPermissions(userId: number): Promise<UserPermissions> {
    return request.get(`${this.baseURL}/rbac/users/${userId}/permissions`);
  }

  // 数据权限
  async createDataPermission(data: CreateDataPermissionRequest): Promise<DataPermission> {
    return request.post(`${this.baseURL}/rbac/data-permissions`, data);
  }

  async getDataPermissions(roleId: number, resourceType: string): Promise<DataPermission> {
    return request.get(`${this.baseURL}/rbac/data-permissions`, {
      params: { role_id: roleId, resource_type: resourceType }
    });
  }

  async updateDataPermission(id: number, data: UpdateDataPermissionRequest): Promise<void> {
    return request.put(`${this.baseURL}/rbac/data-permissions/${id}`, data);
  }

  async deleteDataPermission(id: number): Promise<void> {
    return request.delete(`${this.baseURL}/rbac/data-permissions/${id}`);
  }

  // 字段权限
  async createFieldPermission(data: CreateFieldPermissionRequest): Promise<FieldPermission> {
    return request.post(`${this.baseURL}/rbac/field-permissions`, data);
  }

  async getFieldPermissions(roleId: number, resourceType: string): Promise<FieldPermission[]> {
    const result = await request.get<{ items: FieldPermission[] }>(
      `${this.baseURL}/rbac/field-permissions`,
      { params: { role_id: roleId, resource_type: resourceType } }
    );
    return result.items;
  }

  async updateFieldPermission(id: number, data: UpdateFieldPermissionRequest): Promise<void> {
    return request.put(`${this.baseURL}/rbac/field-permissions/${id}`, data);
  }

  async updateFieldPermissionsBatch(
    roleId: number,
    data: UpdateFieldPermissionsBatchRequest
  ): Promise<void> {
    return request.put(`${this.baseURL}/rbac/roles/${roleId}/field-permissions/batch`, data);
  }

  // 临时授权
  async createTemporaryGrant(data: CreateTemporaryGrantRequest): Promise<TemporaryGrant> {
    return request.post(`${this.baseURL}/rbac/temporary-grants`, data);
  }

  async listTemporaryGrants(params: ListTemporaryGrantsParams): Promise<PageResult<TemporaryGrant>> {
    return request.get(`${this.baseURL}/rbac/temporary-grants`, { params });
  }

  async revokeTemporaryGrant(grantCode: string, reason?: string): Promise<void> {
    return request.post(`${this.baseURL}/rbac/temporary-grants/${grantCode}/revoke`, { reason });
  }

  // 审计日志
  async listAuditLogs(params: ListAuditLogsParams): Promise<PageResult<AuditLog>> {
    return request.get(`${this.baseURL}/rbac/audit-logs`, { params });
  }
}

export const rbacApi = new RBACApiClientImpl();
```

### 11.2 数据权限配置组件

```typescript
// components/rbac/DataPermissionConfig.tsx
import React, { useState, useEffect } from 'react';
import { Form, Select, TreeSelect, Button, Card, message } from '@douyinfe/semi-ui';
import { rbacApi } from '@api/rbac-api';

interface DataPermissionConfigProps {
  roleId: number;
  resourceType: string;
  onSuccess?: () => void;
}

export const DataPermissionConfig: React.FC<DataPermissionConfigProps> = ({
  roleId,
  resourceType,
  onSuccess
}) => {
  const [form] = Form.useForm();
  const [orgTree, setOrgTree] = useState([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchOrgTree();
    loadExistingPermission();
  }, [roleId, resourceType]);

  const fetchOrgTree = async () => {
    try {
      const response = await fetch('/api/v1/organizations/tree');
      const data = await response.json();
      setOrgTree(data.data);
    } catch (error) {
      message.error('加载组织树失败');
    }
  };

  const loadExistingPermission = async () => {
    try {
      const permission = await rbacApi.getDataPermissions(roleId, resourceType);
      if (permission) {
        form.setValues({
          permission_scope: permission.permission_scope,
          custom_org_ids: permission.custom_org_ids,
          custom_user_ids: permission.custom_user_ids,
          permissions: permission.permissions,
        });
      }
    } catch (error) {
      console.error('加载现有权限失败', error);
    }
  };

  const handleSubmit = async (values: any) => {
    setLoading(true);
    try {
      await rbacApi.createDataPermission({
        role_id: roleId,
        resource_type: resourceType,
        ...values,
      });
      message.success('数据权限配置保存成功');
      onSuccess?.();
    } catch (error) {
      message.error('保存失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card title="数据权限配置">
      <Form
        form={form}
        onSubmit={handleSubmit}
        initValues={{
          permission_scope: 'department',
          permissions: ['read'],
        }}
      >
        <Form.Select
          field="permission_scope"
          label="权限范围"
          style={{ width: 300 }}
          rules={[{ required: true }]}
        >
          <Select.Option value="all">全部数据</Select.Option>
          <Select.Option value="department_and_below">本部门及下级</Select.Option>
          <Select.Option value="department">本部门</Select.Option>
          <Select.Option value="self">仅本人</Select.Option>
          <Select.Option value="custom">自定义</Select.Option>
        </Form.Select>

        <Form.Dependency>
          {({ permission_scope }) =>
            permission_scope === 'custom' && (
              <Form.TreeSelect
                field="custom_org_ids"
                label="选择组织"
                treeData={orgTree}
                multiple
                placeholder="请选择允许访问的组织"
                style={{ width: 400 }}
              />
            )
          }
        </Form.Dependency>

        <Form.CheckboxGroup
          field="permissions"
          label="操作权限"
          options={[
            { label: '查看', value: 'read' },
            { label: '编辑', value: 'write' },
            { label: '删除', value: 'delete' },
          ]}
        />

        <Button
          type="primary"
          htmlType="submit"
          loading={loading}
        >
          保存配置
        </Button>
      </Form>
    </Card>
  );
};
```

### 11.3 字段权限配置组件

```typescript
// components/rbac/FieldPermissionConfig.tsx
import React, { useState, useEffect } from 'react';
import { Table, Form, Select, Input, Button, Card, message } from '@douyinfe/semi-ui';
import { rbacApi } from '@api/rbac-api';

interface FieldPermissionConfigProps {
  roleId: number;
  resourceType: string;
  fieldList: string[];
  onSuccess?: () => void;
}

export const FieldPermissionConfig: React.FC<FieldPermissionConfigProps> = ({
  roleId,
  resourceType,
  fieldList,
  onSuccess
}) => {
  const [loading, setLoading] = useState(false);
  const [fieldPermissions, setFieldPermissions] = useState<any[]>([]);

  useEffect(() => {
    loadFieldPermissions();
  }, [roleId, resourceType]);

  const loadFieldPermissions = async () => {
    try {
      const permissions = await rbacApi.getFieldPermissions(roleId, resourceType);
      setFieldPermissions(permissions);
    } catch (error) {
      console.error('加载字段权限失败', error);
    }
  };

  const handleBatchUpdate = async () => {
    setLoading(true);
    try {
      await rbacApi.updateFieldPermissionsBatch(roleId, {
        resource_type: resourceType,
        fields: fieldPermissions,
      });
      message.success('批量更新成功');
      onSuccess?.();
    } catch (error) {
      message.error('更新失败');
    } finally {
      setLoading(false);
    }
  };

  const columns = [
    {
      title: '字段名称',
      dataIndex: 'field_name',
      render: (text: string) => <span>{text}</span>,
    },
    {
      title: '权限级别',
      dataIndex: 'permission',
      render: (value: string, record: any, index: number) => (
        <Select
          value={value}
          onChange={(v) => {
            const newData = [...fieldPermissions];
            newData[index].permission = v;
            setFieldPermissions(newData);
          }}
          style={{ width: 150 }}
        >
          <Select.Option value="visible">可见</Select.Option>
          <Select.Option value="masked">脱敏</Select.Option>
          <Select.Option value="hidden">隐藏</Select.Option>
        </Select>
      ),
    },
    {
      title: '脱敏规则',
      dataIndex: 'mask_rule',
      render: (value: string, record: any, index: number) => (
        record.permission === 'masked' && (
          <Select
            value={value}
            onChange={(v) => {
              const newData = [...fieldPermissions];
              newData[index].mask_rule = v;
              setFieldPermissions(newData);
            }}
            style={{ width: 150 }}
          >
            <Select.Option value="phone">手机号</Select.Option>
            <Select.Option value="id_card">身份证</Select.Option>
            <Select.Option value="email">邮箱</Select.Option>
            <Select.Option value="bank_card">银行卡</Select.Option>
          </Select>
        )
      ),
    },
  ];

  return (
    <Card title="字段权限配置">
      <Table
        columns={columns}
        dataSource={fieldPermissions}
        rowKey="field_name"
        pagination={false}
      />

      <Button
        type="primary"
        onClick={handleBatchUpdate}
        loading={loading}
        style={{ marginTop: 16 }}
      >
        批量保存
      </Button>
    </Card>
  );
};
```

---

## 12. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 10001 | 400 | 用户名已存在 |
| 10002 | 400 | 邮箱已存在 |
| 10003 | 400 | 手机号已存在 |
| 10004 | 404 | 用户不存在 |
| 10005 | 403 | 无权限操作该用户 |
| 10006 | 400 | 密码不符合安全策略 |
| 10007 | 400 | 用户状态异常，无法操作 |
| 10101 | 400 | 角色名已存在 |
| 10102 | 404 | 角色不存在 |
| 10103 | 403 | 系统角色不可修改/删除 |
| 10104 | 400 | 角色仍有用户，无法删除 |
| 10201 | 400 | 权限代码不存在 |
| 10202 | 400 | 权限已分配给角色 |
| 10301 | 400 | 数据权限配置已存在 |
| 10302 | 404 | 数据权限配置不存在 |
| 10303 | 400 | 自定义组织ID列表为空 |
| 10401 | 400 | 字段权限配置已存在 |
| 10402 | 404 | 字段权限配置不存在 |
| 10501 | 400 | 临时授权已存在 |
| 10502 | 400 | 临时授权已过期 |
| 10503 | 403 | 无权限撤销该授权 |
| 10601 | 403 | 权限不足 |
| 10602 | 403 | 数据权限不足 |
| 10603 | 403 | 字段权限不足 |

---

**文档结束**
