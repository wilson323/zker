# API接口文档：企业应用中心模块

**模块名称**: 企业应用中心 (EnterpriseAppCenter)
**设计文档**: 20-企业管理中台_企业应用中心.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 应用管理API](#2-应用管理api)
- [3. 应用分发API](#3-应用分发api)
- [4. 应用监控API](#4-应用监控api)
- [5. 单点登录API](#5-单点登录api)
- [6. 数据模型](#6-数据模型)
- [7. 错误码定义](#7-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

企业应用中心通过**应用上架 + 应用分发 + 应用管理**，为企业提供统一的应用入口和管理能力：

- ✅ **应用商店** - 企业内部应用商店
- ✅ **统一入口** - 单点登录(SSO)
- ✅ **应用管理** - 应用生命周期管理
- ✅ **应用分发** - 用户安装/卸载应用
- ✅ **使用监控** - 应用使用统计

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5
- 实现: 30% 编码（应用框架） + 70% 配置（应用数据）

**前端技术栈**:
- React 18 + TypeScript 5.6
- Semi Design UI组件库

**实现策略**: ✅ 30% 编码 + 70% 配置

### 1.3 数据库表

| 表名 | 说明 |
|------|------|
| `enterprise_apps` | 企业应用表 |
| `user_app_installations` | 用户应用安装表 |
| `app_usage_logs` | 应用使用日志表 |

---

## 2. 应用管理API

### 2.1 获取应用列表

**接口地址**: `GET /api/v1/apps`

**查询参数**:
- keyword: 搜索关键词 (可选)
- category: 分类筛选 (可选)
- is_active: 是否启用 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "id": "app-hr-system-001",
        "tenant_id": "tenant-001",
        "name": "HR人力资源系统",
        "icon": "https://example.com/icons/hr-system.png",
        "description": "企业人力资源管理平台",
        "app_url": "https://hr.company.com",
        "sso_enabled": true,
        "category": "企业管理",
        "is_active": true,
        "install_count": 1250,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      },
      {
        "id": "app-oa-system-001",
        "tenant_id": "tenant-001",
        "name": "OA办公系统",
        "icon": "https://example.com/icons/oa-system.png",
        "description": "企业办公自动化平台",
        "app_url": "https://oa.company.com",
        "sso_enabled": true,
        "category": "办公协作",
        "is_active": true,
        "install_count": 2100,
        "created_at": "2025-01-01T00:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ]
  }
}
```

### 2.2 创建应用

**接口地址**: `POST /api/v1/apps`

**权限**: 仅管理员

**请求参数**:
```json
{
  "name": "CRM客户管理系统",
  "icon": "https://example.com/icons/crm.png",
  "description": "企业客户关系管理平台",
  "app_url": "https://crm.company.com",
  "sso_enabled": true,
  "category": "企业管理"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": "app-crm-001",
    "name": "CRM客户管理系统",
    "created_at": "2025-01-15T14:00:00Z"
  }
}
```

### 2.3 更新应用

**接口地址**: `PUT /api/v1/apps/{app_id}`

**权限**: 仅管理员

**请求参数**:
```json
{
  "name": "CRM客户管理系统(更新版)",
  "description": "企业客户关系管理平台",
  "app_url": "https://crm2.company.com",
  "sso_enabled": true,
  "is_active": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": "app-crm-001",
    "updated_at": "2025-01-15T15:00:00Z"
  }
}
```

### 2.4 删除应用

**接口地址**: `DELETE /api/v1/apps/{app_id}`

**权限**: 仅管理员

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_app_id": "app-crm-001"
  }
}
```

### 2.5 获取应用详情

**接口地址**: `GET /api/v1/apps/{app_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "app-hr-system-001",
    "tenant_id": "tenant-001",
    "name": "HR人力资源系统",
    "icon": "https://example.com/icons/hr-system.png",
    "description": "企业人力资源管理平台，包括员工管理、薪资管理、考勤管理等功能",
    "app_url": "https://hr.company.com",
    "sso_enabled": true,
    "category": "企业管理",
    "is_active": true,
    "install_count": 1250,
    "active_users": 890,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

---

## 3. 应用分发API

### 3.1 安装应用

**接口地址**: `POST /api/v1/apps/{app_id}/install`

**请求参数**:
```json
{}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "安装成功",
  "data": {
    "installation_id": "install-001",
    "app_id": "app-hr-system-001",
    "app_name": "HR人力资源系统",
    "installed_at": "2025-01-15T16:00:00Z"
  }
}
```

### 3.2 卸载应用

**接口地址**: `DELETE /api/v1/apps/{app_id}/uninstall`

**响应示例**:
```json
{
  "code": 0,
  "message": "卸载成功",
  "data": {
    "installation_id": "install-001",
    "uninstalled_at": "2025-01-15T17:00:00Z"
  }
}
```

### 3.3 获取我的应用

**接口地址**: `GET /api/v1/my-apps`

**查询参数**:
- keyword: 搜索关键词 (可选)
- category: 分类筛选 (可选)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 8,
    "items": [
      {
        "installation_id": "install-001",
        "app_id": "app-hr-system-001",
        "app_name": "HR人力资源系统",
        "app_icon": "https://example.com/icons/hr-system.png",
        "app_url": "https://hr.company.com",
        "sso_enabled": true,
        "category": "企业管理",
        "installed_at": "2025-01-01T00:00:00Z",
        "last_accessed_at": "2025-01-15T10:30:00Z"
      },
      {
        "installation_id": "install-002",
        "app_id": "app-oa-system-001",
        "app_name": "OA办公系统",
        "app_icon": "https://example.com/icons/oa-system.png",
        "app_url": "https://oa.company.com",
        "sso_enabled": true,
        "category": "办公协作",
        "installed_at": "2025-01-05T00:00:00Z",
        "last_accessed_at": "2025-01-15T09:15:00Z"
      }
    ]
  }
}
```

### 3.4 批量安装应用

**接口地址**: `POST /api/v1/apps/batch-install`

**请求参数**:
```json
{
  "app_ids": [
    "app-hr-system-001",
    "app-oa-system-001",
    "app-crm-001"
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "批量安装成功",
  "data": {
    "successful": [
      {
        "app_id": "app-hr-system-001",
        "installation_id": "install-001"
      },
      {
        "app_id": "app-oa-system-001",
        "installation_id": "install-002"
      }
    ],
    "failed": [
      {
        "app_id": "app-crm-001",
        "error": "应用不存在或已停用"
      }
    ]
  }
}
```

---

## 4. 应用监控API

### 4.1 获取应用使用统计

**接口地址**: `GET /api/v1/apps/{app_id}/stats`

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)
- granularity: 时间粒度 (可选，day/week/month，默认day)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "app_id": "app-hr-system-001",
    "app_name": "HR人力资源系统",
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "summary": {
      "total_installs": 1250,
      "total_uninstalls": 15,
      "active_users": 890,
      "total_sessions": 15200,
      "avg_session_duration": 320
    },
    "daily_stats": [
      {
        "date": "2025-01-01",
        "new_installs": 85,
        "uninstalls": 1,
        "active_users": 820,
        "sessions": 1050,
        "avg_session_duration": 315
      },
      {
        "date": "2025-01-02",
        "new_installs": 92,
        "uninstalls": 0,
        "active_users": 850,
        "sessions": 1120,
        "avg_session_duration": 325
      }
    ]
  }
}
```

### 4.2 获取应用访问日志

**接口地址**: `GET /api/v1/apps/{app_id}/access-logs`

**权限**: 仅管理员

**查询参数**:
- user_id: 用户ID筛选 (可选)
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
        "log_id": "log-001",
        "user_id": 1001,
        "username": "张三",
        "app_id": "app-hr-system-001",
        "app_name": "HR人力资源系统",
        "action": "launch",
        "accessed_at": "2025-01-15T10:30:00Z",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0..."
      }
    ]
  }
}
```

### 4.3 获取用户应用使用排行

**接口地址**: `GET /api/v1/apps/user-ranking`

**权限**: 仅管理员

**查询参数**:
- start_date: 开始日期 (必填)
- end_date: 结束日期 (必填)
- limit: 返回数量 (默认10)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-01",
      "end": "2025-01-15"
    },
    "items": [
      {
        "user_id": 1001,
        "username": "张三",
        "department": "技术部",
        "total_sessions": 156,
        "total_duration": 12500,
        "apps_used": 8
      },
      {
        "user_id": 1002,
        "username": "李四",
        "department": "市场部",
        "total_sessions": 142,
        "total_duration": 11800,
        "apps_used": 7
      }
    ]
  }
}
```

---

## 5. 单点登录API

### 5.1 生成SSO登录令牌

**接口地址**: `POST /api/v1/apps/{app_id}/sso/token`

**请求参数**:
```json
{
  "redirect_url": "https://hr.company.com/sso/callback"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "sso_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "login_url": "https://hr.company.com/sso/login?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-01-15T11:00:00Z",
    "expires_in": 3600
  }
}
```

### 5.2 验证SSO令牌

**接口地址**: `POST /api/v1/sso/verify`

**说明**: 此接口供第三方应用调用，验证SSO令牌并获取用户信息

**请求参数**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "验证成功",
  "data": {
    "user_id": 1001,
    "username": "zhangsan",
    "email": "zhangsan@company.com",
    "name": "张三",
    "department": "技术部",
    "position": "软件工程师",
    "tenant_id": "tenant-001",
    "roles": ["user", "developer"],
    "expires_at": "2025-01-15T11:00:00Z"
  }
}
```

### 5.3 SSO登出

**接口地址**: `POST /api/v1/sso/logout`

**请求参数**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "logout_url": "https://hr.company.com/sso/logout"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登出成功",
  "data": {
    "logged_out_at": "2025-01-15T10:45:00Z"
  }
}
```

---

## 6. 数据模型

### 6.1 EnterpriseApp

```typescript
interface EnterpriseApp {
  id: string;
  tenant_id: string;
  name: string;
  icon: string;
  description: string;
  app_url: string;
  sso_enabled: boolean;
  category: string;
  is_active: boolean;
  install_count: number;
  active_users?: number;
  created_at: Date;
  updated_at: Date;
}
```

### 6.2 UserAppInstallation

```typescript
interface UserAppInstallation {
  id: number;
  tenant_id: string;
  user_id: number;
  app_id: string;
  installed_at: Date;
  last_accessed_at?: Date;
}

interface MyApp {
  installation_id: number;
  app_id: string;
  app_name: string;
  app_icon: string;
  app_url: string;
  sso_enabled: boolean;
  category: string;
  installed_at: Date;
  last_accessed_at?: Date;
}
```

### 6.3 AppStats

```typescript
interface AppStats {
  app_id: string;
  app_name: string;
  period: {
    start: string;
    end: string;
  };
  summary: {
    total_installs: number;
    total_uninstalls: number;
    active_users: number;
    total_sessions: number;
    avg_session_duration: number;
  };
  daily_stats: DailyStats[];
}

interface DailyStats {
  date: string;
  new_installs: number;
  uninstalls: number;
  active_users: number;
  sessions: number;
  avg_session_duration: number;
}
```

### 6.4 SSOToken

```typescript
interface SSOToken {
  sso_token: string;
  login_url: string;
  expires_at: Date;
  expires_in: number;
}

interface SSOUserInfo {
  user_id: number;
  username: string;
  email: string;
  name: string;
  department: string;
  position: string;
  tenant_id: string;
  roles: string[];
  expires_at: Date;
}
```

---

## 7. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 35001 | 404 | 应用不存在 |
| 35002 | 400 | 应用名称已存在 |
| 35003 | 400 | 应用URL无效 |
| 35004 | 400 | SSO配置无效 |
| 35005 | 403 | 无权限访问应用 |
| 35006 | 400 | 应用已安装 |
| 35007 | 404 | 应用未安装 |
| 35008 | 400 | 应用已停用，无法安装 |
| 35009 | 401 | SSO令牌无效或已过期 |
| 35010 | 400 | SSO验证失败 |

---

**文档结束**
