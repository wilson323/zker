# API接口文档：租户安全合规模块

**模块名称**: 租户安全合规 (TenantSecurity)
**设计文档**: 25-MultiTenant_SaaS核心_租户安全合规.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 身份认证API](#2-身份认证api)
- [3. 审计日志API](#3-审计日志api)
- [4. 敏感操作API](#4-敏感操作api)
- [5. GDPR合规API](#5-gdpr合规api)
- [6. 数据模型](#6-数据模型)
- [7. 错误码定义](#7-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

租户安全合规是MultiTenant SaaS的安全合规模块，通过**多层级安全防护 + 合规审计**，确保平台满足等保三级、GDPR、SOC2等合规要求：

- ✅ **身份认证** - 用户名密码、手机验证码、OAuth2、SSO、MFA
- ✅ **会话管理** - JWT Token、Token刷新、强制登出、异地登录检测
- ✅ **权限控制** - RBAC权限模型、数据隔离、审计日志
- ✅ **安全审计** - 操作日志、访问日志、敏感操作日志、日志归档
- ✅ **合规认证** - 等保三级、GDPR、SOC2
- ✅ **数据脱敏** - 邮箱、手机号、身份证等敏感信息脱敏

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 加密: AES-256 (存储)、TLS 1.3 (传输)
- 认证: JWT、OAuth2、SAML
- 实现: 50% 编码（安全框架） + 50% 配置（合规规则）

**实现策略**: ✅ 50% 编码 + 50% 配置

### 1.3 数据库表

| 表名 | 说明 |
|------|------|
| `audit_logs` | 审计日志表 |
| `login_logs` | 登录日志表 |
| `sensitive_operations` | 敏感操作审计表 |
| `user_consents` | 用户同意记录表 |
| `data_deletion_requests` | 数据删除请求表 |

---

## 2. 身份认证API

### 2.1 用户登录

**接口地址**: `POST /api/v1/auth/login`

**请求参数**:
```json
{
  "login_method": "password",
  "username": "user@example.com",
  "password": "encrypted_password",
  "mfa_code": "123456",
  "remember_me": false
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1001,
      "username": "zhangsan",
      "email": "zhangsan@example.com",
      "name": "张三",
      "tenant_id": "tenant-001",
      "roles": ["user"],
      "mfa_enabled": true
    }
  }
}
```

### 2.2 手机验证码登录

**接口地址**: `POST /api/v1/auth/login/sms`

**请求参数**:
```json
{
  "phone": "+8613812345678",
  "sms_code": "123456",
  "remember_me": false
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1001,
      "username": "zhangsan",
      "phone": "+86138****5678",
      "name": "张三",
      "tenant_id": "tenant-001",
      "roles": ["user"]
    }
  }
}
```

### 2.3 OAuth2第三方登录

**接口地址**: `POST /api/v1/auth/login/oauth`

**请求参数**:
```json
{
  "provider": "wechat_work",
  "code": "authorization_code_from_provider",
  "redirect_uri": "https://example.com/oauth/callback"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 1001,
      "username": "zhangsan",
      "email": "zhangsan@example.com",
      "name": "张三",
      "tenant_id": "tenant-001",
      "roles": ["user"]
    }
  }
}
```

### 2.4 刷新Token

**接口地址**: `POST /api/v1/auth/refresh`

**请求参数**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "刷新成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

### 2.5 用户登出

**接口地址**: `POST /api/v1/auth/logout`

**请求头**: `Authorization: Bearer {access_token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "登出成功",
  "data": {
    "logged_out_at": "2025-01-15T10:30:00Z"
  }
}
```

### 2.6 强制登出所有设备

**接口地址**: `POST /api/v1/auth/logout/all`

**请求头**: `Authorization: Bearer {access_token}`

**响应示例**:
```json
{
  "code": 0,
  "message": "已强制登出所有设备",
  "data": {
    "logged_out_devices": 5
  }
}
```

### 2.7 发送验证码

**接口地址**: `POST /api/v1/auth/sms/send`

**请求参数**:
```json
{
  "phone": "+8613812345678",
  "purpose": "login"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "验证码已发送",
  "data": {
    "expires_in": 300,
    "resend_after": 60
  }
}
```

---

## 3. 审计日志API

### 3.1 查询审计日志

**接口地址**: `GET /api/v1/audit/logs`

**权限**: 管理员或审计员

**查询参数**:
- tenant_id: 租户ID (可选)
- user_id: 用户ID (可选)
- resource_type: 资源类型 (可选)
- resource_id: 资源ID (可选)
- action: 操作类型 (可选)
- sensitivity: 敏感级别 (可选，low/medium/high/critical)
- status: 操作状态 (可选，success/failure/partial)
- start_time: 开始时间 (可选)
- end_time: 结束时间 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 1520,
    "items": [
      {
        "id": 100001,
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "user_id": 1001,
        "username": "zhangsan",
        "actor_type": "user",
        "actor_id": "1001",
        "resource_type": "bot",
        "resource_id": "bot-123",
        "resource_name": "客服助手",
        "action": "update",
        "action_detail": "更新Bot配置：修改系统提示词",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "request_id": "req-abc-123",
        "status": "success",
        "error_code": null,
        "error_message": null,
        "sensitivity": "medium",
        "created_at": "2025-01-15T10:30:00.123Z"
      },
      {
        "id": 100002,
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "user_id": null,
        "username": null,
        "actor_type": "system",
        "actor_id": "system-001",
        "resource_type": "bot",
        "resource_id": "bot-123",
        "resource_name": "客服助手",
        "action": "delete",
        "action_detail": "系统自动清理过期对话记录",
        "ip_address": null,
        "user_agent": null,
        "request_id": "req-sys-456",
        "status": "success",
        "sensitivity": "low",
        "created_at": "2025-01-15T10:00:00.000Z"
      }
    ]
  }
}
```

### 3.2 查询登录日志

**接口地址**: `GET /api/v1/audit/login-logs`

**权限**: 管理员或审计员

**查询参数**:
- tenant_id: 租户ID (可选)
- user_id: 用户ID (可选)
- login_method: 登录方式 (可选，password/sms/oauth/saml)
- status: 登录状态 (可选，success/failure/blocked)
- start_time: 开始时间 (可选)
- end_time: 结束时间 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 320,
    "items": [
      {
        "id": 5001,
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "user_id": 1001,
        "username": "zhangsan",
        "login_method": "password",
        "ip_address": "192.168.1.100",
        "location": {
          "country": "CN",
          "province": "Beijing",
          "city": "Beijing"
        },
        "device_info": {
          "type": "desktop",
          "os": "Windows 10",
          "browser": "Chrome 120"
        },
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...",
        "status": "success",
        "failure_reason": null,
        "created_at": "2025-01-15T10:30:00Z"
      },
      {
        "id": 5002,
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "user_id": 1001,
        "username": "zhangsan",
        "login_method": "password",
        "ip_address": "192.168.1.101",
        "location": {
          "country": "CN",
          "province": "Shanghai",
          "city": "Shanghai"
        },
        "device_info": {
          "type": "mobile",
          "os": "iOS 17",
          "browser": "Safari 17"
        },
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X)...",
        "status": "failure",
        "failure_reason": "密码错误",
        "created_at": "2025-01-15T10:25:00Z"
      }
    ]
  }
}
```

### 3.3 获取安全事件统计

**接口地址**: `GET /api/v1/audit/security-events`

**权限**: 管理员或审计员

**查询参数**:
- tenant_id: 租户ID (可选)
- start_time: 开始时间 (必填)
- end_time: 结束时间 (必填)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-01T00:00:00Z",
      "end": "2025-01-15T23:59:59Z"
    },
    "summary": {
      "total_events": 15200,
      "successful_logins": 12500,
      "failed_logins": 2500,
      "suspicious_activities": 120,
      "blocked_attempts": 80
    },
    "by_type": {
      "login_success": 12500,
      "login_failure": 2500,
      "permission_denied": 150,
      "suspicious_location": 45,
      "multiple_failures": 30
    },
    "top_users": {
      "by_login_attempts": [
        {
          "user_id": 1001,
          "username": "zhangsan",
          "attempts": 320
        }
      ],
      "by_failed_attempts": [
        {
          "user_id": 1002,
          "username": "lisi",
          "failed_attempts": 25
        }
      ]
    },
    "suspicious_locations": [
      {
        "ip_address": "203.0.113.1",
        "location": "Unknown",
        "attempts": 15
      }
    ]
  }
}
```

---

## 4. 敏感操作API

### 4.1 查询敏感操作列表

**接口地址**: `GET /api/v1/audit/sensitive-operations`

**权限**: 管理员或审计员

**查询参数**:
- tenant_id: 租户ID (可选)
- operator_id: 操作人ID (可选)
- operation_type: 操作类型 (可选)
- status: 状态 (可选，pending/approved/rejected/executed/failed)
- start_time: 开始时间 (可选)
- end_time: 结束时间 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 85,
    "items": [
      {
        "id": 3001,
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "operator_id": 1001,
        "operator_name": "管理员",
        "operator_type": "admin",
        "operation_type": "delete_bot",
        "resource_type": "bot",
        "resource_id": "bot-123",
        "resource_name": "客服助手",
        "reason": "Bot已废弃，需要删除",
        "approval_id": "approval-001",
        "ip_address": "192.168.1.100",
        "status": "approved",
        "created_at": "2025-01-15T10:00:00Z",
        "approved_at": "2025-01-15T10:15:00Z",
        "approved_by": 1002,
        "approved_by_name": "超级管理员"
      },
      {
        "id": 3002,
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "operator_id": 1003,
        "operator_name": "数据分析师",
        "operator_type": "user",
        "operation_type": "export_data",
        "resource_type": "conversation",
        "resource_id": "all",
        "resource_name": "所有对话记录",
        "reason": "数据分析需要导出对话数据",
        "approval_id": null,
        "ip_address": "192.168.1.105",
        "status": "pending",
        "created_at": "2025-01-15T11:00:00Z"
      }
    ]
  }
}
```

### 4.2 提交敏感操作申请

**接口地址**: `POST /api/v1/audit/sensitive-operations`

**请求参数**:
```json
{
  "operation_type": "export_data",
  "resource_type": "conversation",
  "resource_id": "all",
  "reason": "数据分析需要导出对话数据"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "申请已提交，等待审批",
  "data": {
    "id": 3003,
    "operation_type": "export_data",
    "resource_type": "conversation",
    "resource_id": "all",
    "status": "pending",
    "created_at": "2025-01-15T12:00:00Z"
  }
}
```

### 4.3 审批敏感操作

**接口地址**: `POST /api/v1/audit/sensitive-operations/{operation_id}/approve`

**权限**: 仅管理员或审批人

**请求参数**:
```json
{
  "action": "approve",
  "comment": "同意导出，请确保数据安全"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "审批成功",
  "data": {
    "id": 3003,
    "status": "approved",
    "approved_by": 1001,
    "approved_by_name": "管理员",
    "approved_at": "2025-01-15T12:15:00Z",
    "comment": "同意导出，请确保数据安全"
  }
}
```

### 4.4 撤销敏感操作申请

**接口地址**: `POST /api/v1/audit/sensitive-operations/{operation_id}/cancel`

**权限**: 仅申请本人或管理员

**响应示例**:
```json
{
  "code": 0,
  "message": "申请已撤销",
  "data": {
    "id": 3003,
    "status": "cancelled",
    "cancelled_at": "2025-01-15T12:10:00Z"
  }
}
```

---

## 5. GDPR合规API

### 5.1 导出个人数据

**接口地址**: `POST /api/v1/gdpr/data-export`

**请求参数**:
```json
{
  "include_conversations": true,
  "include_bots": true,
  "include_analytics": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "数据导出任务已创建",
  "data": {
    "export_id": "export-20250115-001",
    "status": "processing",
    "estimated_completion_time": "2025-01-15T13:00:00Z",
    "download_url": null,
    "expires_at": "2025-01-22T00:00:00Z"
  }
}
```

### 5.2 查询导出任务状态

**接口地址**: `GET /api/v1/gdpr/data-export/{export_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "export_id": "export-20250115-001",
    "user_id": 1001,
    "status": "completed",
    "file_size": 15728640,
    "file_count": 5,
    "created_at": "2025-01-15T12:00:00Z",
    "completed_at": "2025-01-15T12:30:00Z",
    "download_url": "/api/v1/gdpr/data-export/export-20250115-001/download",
    "expires_at": "2025-01-22T00:00:00Z"
  }
}
```

### 5.3 下载导出数据

**接口地址**: `GET /api/v1/gdpr/data-export/{export_id}/download`

**响应**: 二进制文件流

### 5.4 删除个人数据

**接口地址**: `POST /api/v1/gdpr/data-delete`

**请求参数**:
```json
{
  "reason": "根据GDPR权利，要求删除个人数据",
  "delete_conversations": true,
  "delete_bots": true,
  "delete_analytics": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "删除请求已提交",
  "data": {
    "request_id": "delete-20250115-001",
    "status": "pending",
    "estimated_processing_time": "2025-01-15T14:00:00Z"
  }
}
```

### 5.5 查询删除请求状态

**接口地址**: `GET /api/v1/gdpr/data-delete/{request_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "request_id": "delete-20250115-001",
    "user_id": 1001,
    "request_type": "data_erasure",
    "status": "completed",
    "created_at": "2025-01-15T13:00:00Z",
    "processed_at": "2025-01-15T13:30:00Z",
    "deleted_records": {
      "conversations": 1520,
      "bots": 5,
      "analytics": 320
    }
  }
}
```

### 5.6 查看用户同意记录

**接口地址**: `GET /api/v1/gdpr/consents`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1001,
    "consents": [
      {
        "id": 1,
        "consent_type": "privacy_policy",
        "consent_name": "隐私政策",
        "granted": true,
        "granted_at": "2025-01-01T00:00:00Z",
        "withdrawn_at": null,
        "version": "v2.0"
      },
      {
        "id": 2,
        "consent_type": "terms_of_service",
        "consent_name": "服务条款",
        "granted": true,
        "granted_at": "2025-01-01T00:00:00Z",
        "withdrawn_at": null,
        "version": "v3.0"
      },
      {
        "id": 3,
        "consent_type": "data_processing",
        "consent_name": "数据处理同意",
        "granted": true,
        "granted_at": "2025-01-01T00:00:00Z",
        "withdrawn_at": null,
        "version": "v1.0"
      }
    ]
  }
}
```

### 5.7 撤回同意

**接口地址**: `POST /api/v1/gdpr/consents/{consent_id}/withdraw`

**请求参数**:
```json
{
  "reason": "不再同意数据处理"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "同意已撤回",
  "data": {
    "id": 3,
    "consent_type": "data_processing",
    "granted": false,
    "withdrawn_at": "2025-01-15T14:00:00Z"
  }
}
```

### 5.8 给予同意

**接口地址**: `POST /api/v1/gdpr/consents/grant`

**请求参数**:
```json
{
  "consent_type": "data_processing",
  "granted": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "同意已记录",
  "data": {
    "id": 4,
    "consent_type": "data_processing",
    "granted": true,
    "granted_at": "2025-01-15T14:05:00Z"
  }
}
```

---

## 6. 数据模型

### 6.1 AuditLog

```typescript
interface AuditLog {
  id: number;
  tenant_id: string;
  tenant_name?: string;
  user_id?: number;
  username?: string;
  actor_type: 'user' | 'system' | 'api';
  actor_id?: string;
  resource_type: string;
  resource_id?: string;
  resource_name?: string;
  action: string;
  action_detail?: string;
  ip_address?: string;
  user_agent?: string;
  request_id?: string;
  status: 'success' | 'failure' | 'partial';
  error_code?: string;
  error_message?: string;
  sensitivity: 'low' | 'medium' | 'high' | 'critical';
  created_at: Date;
  created_at_ms?: number;
}
```

### 6.2 LoginLog

```typescript
interface LoginLog {
  id: number;
  tenant_id?: string;
  tenant_name?: string;
  user_id?: number;
  username?: string;
  login_method: 'password' | 'sms' | 'oauth' | 'saml';
  ip_address: string;
  location?: {
    country: string;
    province: string;
    city: string;
  };
  device_info?: {
    type: 'desktop' | 'mobile' | 'tablet';
    os: string;
    browser: string;
  };
  user_agent?: string;
  status: 'success' | 'failure' | 'blocked';
  failure_reason?: string;
  created_at: Date;
}
```

### 6.3 SensitiveOperation

```typescript
interface SensitiveOperation {
  id: number;
  tenant_id: string;
  tenant_name?: string;
  operator_id: number;
  operator_name: string;
  operator_type: 'user' | 'admin' | 'system';
  operation_type: 'delete_bot' | 'delete_user' | 'export_data' | 'change_permissions' | 'access_sensitive_data';
  resource_type: string;
  resource_id: string;
  resource_name?: string;
  reason?: string;
  approval_id?: string;
  ip_address?: string;
  status: 'pending' | 'approved' | 'rejected' | 'executed' | 'failed';
  created_at: Date;
  approved_at?: Date;
  approved_by?: number;
  approved_by_name?: string;
}
```

### 6.4 UserConsent

```typescript
interface UserConsent {
  id: number;
  user_id: number;
  consent_type: 'privacy_policy' | 'terms_of_service' | 'data_processing';
  consent_name: string;
  granted: boolean;
  granted_at?: Date;
  withdrawn_at?: Date;
  version: string;
}
```

### 6.5 DataDeletionRequest

```typescript
interface DataDeletionRequest {
  request_id: string;
  user_id: number;
  request_type: 'account_deletion' | 'data_erasure';
  status: 'pending' | 'processing' | 'completed' | 'rejected';
  created_at: Date;
  processed_at?: Date;
  deleted_records?: {
    conversations?: number;
    bots?: number;
    analytics?: number;
  };
}
```

### 6.6 DataExportRequest

```typescript
interface DataExportRequest {
  export_id: string;
  user_id: number;
  status: 'processing' | 'completed' | 'failed';
  file_size?: number;
  file_count?: number;
  created_at: Date;
  completed_at?: Date;
  download_url?: string;
  expires_at: Date;
}
```

---

## 7. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 37001 | 401 | 用户名或密码错误 |
| 37002 | 401 | 验证码错误或已过期 |
| 37003 | 401 | Token无效或已过期 |
| 37004 | 403 | 账户已被禁用 |
| 37005 | 429 | 登录尝试过多，已被暂时锁定 |
| 37006 | 401 | MFA验证码错误 |
| 37007 | 401 | Refresh Token无效或已过期 |
| 37008 | 403 | 无权限访问该资源 |
| 37009 | 403 | 无权限执行该操作 |
| 37101 | 400 | 审计日志查询参数无效 |
| 37102 | 403 | 无权限查询审计日志 |
| 37103 | 400 | 时间范围超过限制 |
| 37201 | 403 | 无权限执行敏感操作 |
| 37202 | 400 | 敏感操作已被拒绝 |
| 37203 | 400 | 敏感操作已在审批中 |
| 37204 | 403 | 无权限审批敏感操作 |
| 37301 | 404 | 导出任务不存在 |
| 37302 | 400 | 导出任务处理中，请稍后再试 |
| 37303 | 404 | 删除请求不存在 |
| 37304 | 400 | 删除请求已在处理中 |
| 37305 | 400 | 导出文件已过期 |
| 37306 | 404 | 用户同意记录不存在 |
| 37307 | 400 | 该同意项不允许撤回 |

---

**文档结束**
