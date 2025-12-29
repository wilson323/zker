# API接口文档:租户隔离机制模块

> **版本**: v1.0.0
> **创建日期**: 2025-12-30
> **对应设计文档**: [22-MultiTenant_SaaS核心_租户隔离机制.md](./详细设计/22-MultiTenant_SaaS核心_租户隔离机制.md)

---

## 📋 目录

1. [概述](#1-概述)
2. [隔离策略管理API](#2-隔离策略管理api)
3. [数据隔离验证API](#3-数据隔离验证api)
4. [租户资源监控API](#4-租户资源监控api)
5. [租户迁移API](#5-租户迁移api)
6. [错误代码](#6-错误代码)

---

## 1. 概述

### 1.1 模块说明

租户隔离机制是Multi-Tenant SaaS系统的核心技术基础,确保多个租户在同一物理基础设施上安全、独立地运行,实现数据、计算、缓存、存储的完全隔离。

### 1.2 隔离策略

| 策略 | 适用场景 | 说明 |
|------|---------|------|
| Row-Level隔离 | 小型租户(< 100用户) | 所有表添加tenant_id字段,成本最低 |
| Schema隔离 | 中型租户(100-1000用户) | 每个租户独立Schema,性能较好 |
| Database隔离 | 大型租户(> 1000用户) | 每个租户独立Database,完全物理隔离 |

### 1.3 基础信息

**Base URL**: `/api/v1`

**认证方式**: JWT Bearer Token (需要系统管理员权限)

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {token}
```

---

## 2. 隔离策略管理API

### 2.1 获取租户隔离策略

查询当前租户的隔离策略配置。

**接口描述**: 获取租户的隔离策略类型和配置

**请求方式**: `GET /tenants/current/isolation-strategy`

**权限要求**: 需要认证,租户管理员

**请求示例**:
```http
GET /api/v1/tenants/current/isolation-strategy
Authorization: Bearer {token}
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "isolation_strategy": {
      "type": "row_level",
      "name": "行级隔离",
      "description": "所有租户共享数据库,通过tenant_id字段隔离数据"
    },
    "tenant_size": {
      "size": "small",
      "user_count": 25,
      "bot_count": 3,
      "storage_gb": 5.5
    },
    "resource_allocation": {
      "database": "coze_enterprise (共享)",
      "schema": null,
      "cache_namespace": "tenant:tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "storage_bucket": "tenant-tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890"
    },
    "upgrade_threshold": {
      "users": 100,
      "storage_gb": 50
    },
    "can_upgrade": true,
    "upgrade_options": [
      {
        "target_strategy": "schema_level",
        "name": "Schema级隔离",
        "recommended_when": "用户数 > 100",
        "estimated_migration_time": "2-4小时"
      }
    ]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_001"
}
```

---

### 2.2 升级隔离策略(管理员)

升级租户的隔离策略。

**接口描述**: 将租户从低级隔离策略升级到高级策略

**请求方式**: `POST /admin/tenants/{tenant_id}/isolation-strategy/upgrade`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求体**:
```json
{
  "target_strategy": "schema_level",
  "scheduled_at": "2025-12-31T02:00:00Z",
  "notify": true
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_strategy | string | ✅ | 目标策略(schema_level/database_level) |
| scheduled_at | string | ❌ | 计划执行时间(ISO8601,默认立即执行) |
| notify | boolean | ❌ | 是否通知租户(默认true) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "隔离策略升级任务已创建",
  "data": {
    "migration_task_id": "migration-task-123456",
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "current_strategy": "row_level",
    "target_strategy": "schema_level",
    "status": "pending",
    "scheduled_at": "2025-12-31T02:00:00Z",
    "estimated_duration": "2-4小时",
    "estimated_downtime": "5-10分钟",
    "created_at": "2025-12-30T15:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_002"
}
```

**错误响应** (400 Bad Request):
```json
{
  "code": 40001,
  "message": "租户规模不满足升级条件",
  "error": {
    "code": "UPGRADE_CONDITION_NOT_MET",
    "current_size": "small",
    "required_size": "medium",
    "current_users": 25,
    "required_users": 100
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_003"
}
```

---

### 2.3 查询迁移任务状态

查询隔离策略迁移任务的执行状态。

**接口描述**: 查询迁移任务的进度和状态

**请求方式**: `GET /admin/tenants/{tenant_id}/isolation-strategy/migration/{task_id}`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |
| task_id | string | 迁移任务ID |

**请求示例**:
```http
GET /api/v1/admin/tenants/tenant-xxx/isolation-strategy/migration/migration-task-123456
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "migration-task-123456",
    "status": "in_progress",
    "progress": {
      "percentage": 65,
      "current_step": "正在迁移数据...",
      "completed_steps": [
        "创建目标Schema",
        "迁移表结构"
      ],
      "total_steps": [
        "创建目标Schema",
        "迁移表结构",
        "迁移数据",
        "验证数据完整性",
        "更新配置",
        "清理旧数据"
      ]
    },
    "started_at": "2025-12-31T02:00:00Z",
    "estimated_completion_at": "2025-12-31T04:30:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_004"
}
```

---

## 3. 数据隔离验证API

### 3.1 验证数据隔离

验证租户间的数据隔离是否正确。

**接口描述**: 检查是否存在跨租户数据泄露风险

**请求方式**: `POST /admin/tenants/{tenant_id}/isolation/verify`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| tenant_id | string | 租户ID |

**请求体**:
```json
{
  "verification_type": "data_leak",
  "tables": ["users", "bots", "conversations"]
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| verification_type | string | ✅ | 验证类型(data_leak/cross_tenant_access) |
| tables | string[] | ❌ | 要验证的表(默认全部表) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "验证完成",
  "data": {
    "tenant_id": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "verification_result": "passed",
    "verified_at": "2025-12-30T16:00:00Z",
    "details": {
      "users_table": {
        "status": "passed",
        "total_records": 25,
        "tenant_records": 25,
        "other_tenant_records": 0,
        "isolation_percentage": 100
      },
      "bots_table": {
        "status": "passed",
        "total_records": 3,
        "tenant_records": 3,
        "other_tenant_records": 0,
        "isolation_percentage": 100
      },
      "conversations_table": {
        "status": "passed",
        "total_records": 1250,
        "tenant_records": 1250,
        "other_tenant_records": 0,
        "isolation_percentage": 100
      }
    },
    "recommendations": [
      "所有表均通过数据隔离验证"
    ]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_005"
}
```

**发现问题时** (200 OK):
```json
{
  "code": 0,
  "message": "发现问题",
  "data": {
    "tenant_id": "tenant-xxx",
    "verification_result": "failed",
    "issues": [
      {
        "table": "conversations",
        "issue_type": "cross_tenant_leak",
        "severity": "high",
        "description": "发现3条不属于当前租户的记录",
        "leaked_records": [
          {
            "id": 101,
            "tenant_id": "tenant-other",
            "user_id": 5
          }
        ],
        "recommendation": "立即修复数据泄露问题"
      }
    ]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_006"
}
```

---

### 3.2 批量验证所有租户隔离

批量验证所有租户的数据隔离。

**接口描述**: 系统管理员验证所有租户的隔离状态

**请求方式**: `POST /admin/tenants/isolation/verify-all`

**权限要求**: 系统管理员

**请求体**:
```json
{
  "tenant_filter": {
    "status": "active",
    "isolation_strategy": "row_level"
  },
  "verification_type": "data_leak",
  "sample_size": 1000
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tenant_filter | object | ❌ | 租户筛选条件 |
| verification_type | string | ✅ | 验证类型 |
| sample_size | integer | ❌ | 每个租户抽样的记录数(默认1000) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "批量验证任务已创建",
  "data": {
    "verification_id": "verify-20251230-001",
    "status": "processing",
    "total_tenants": 150,
    "estimated_duration": "15-30分钟",
    "created_at": "2025-12-30T17:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_007"
}
```

---

### 3.3 查询批量验证结果

获取批量验证任务的结果。

**接口描述**: 查询批量验证任务的执行结果

**请求方式**: `GET /admin/tenants/isolation/verify-all/{verification_id}`

**权限要求**: 系统管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| verification_id | string | 验证任务ID |

**请求示例**:
```http
GET /api/v1/admin/tenants/isolation/verify-all/verify-20251230-001
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "验证完成",
  "data": {
    "verification_id": "verify-20251230-001",
    "status": "completed",
    "summary": {
      "total_tenants": 150,
      "passed": 148,
      "failed": 2,
      "skipped": 0,
      "pass_rate": 98.67
    },
    "failed_tenants": [
      {
        "tenant_id": "tenant-xxx",
        "company_name": "问题企业A",
        "issues": [
          {
            "table": "conversations",
            "issue_type": "cross_tenant_leak",
            "severity": "high",
            "affected_records": 3
          }
        ]
      }
    ],
    "completed_at": "2025-12-30T17:25:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_008"
}
```

---

## 4. 租户资源监控API

### 4.1 获取租户资源使用情况

获取租户的资源使用统计。

**接口描述**: 查询租户的数据库、缓存、存储资源使用情况

**请求方式**: `GET /tenants/current/resource-usage`

**权限要求**: 需要认证,租户管理员

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| resource_type | string | 否 | 资源类型(database/cache/storage) |

**请求示例**:
```http
GET /api/v1/tenants/current/resource-usage
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "database": {
      "isolation_strategy": "row_level",
      "shared_with": 149,
      "table_sizes": {
        "users": "12.5 KB",
        "bots": "45.2 KB",
        "conversations": "1.2 GB",
        "knowledge_bases": "856.3 MB"
      },
      "total_size": "2.01 GB",
      "connection_pool": {
        "active_connections": 5,
        "idle_connections": 15,
        "max_connections": 100
      },
      "query_performance": {
        "avg_query_time_ms": 15.5,
        "slow_queries_today": 3
      }
    },
    "cache": {
      "namespace": "tenant:tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "memory_usage": "256 MB",
      "key_count": 1250,
      "hit_rate": 0.92,
      "eviction_count": 5
    },
    "storage": {
      "bucket": "tenant-tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "total_size_gb": 5.5,
      "object_count": 350,
      "largest_object": "document_123.pdf (2.1 MB)",
      "upload_count_today": 15
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_009"
}
```

---

### 4.2 获取租户性能指标

获取租户的性能指标数据。

**接口描述**: 查询租户的QPS、响应时间等性能指标

**请求方式**: `GET /tenants/current/metrics`

**权限要求**: 需要认证,租户管理员

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | string | 否 | 统计周期(1h/6h/24h/7d/30d,默认24h) |
| metric_type | string | 否 | 指标类型(qps/latency/error_rate) |

**请求示例**:
```http
GET /api/v1/tenants/current/metrics?period=24h&metric_type=qps
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "24h",
    "metrics": {
      "qps": {
        "current": 45.2,
        "avg": 38.5,
        "peak": 125.3,
        "peak_time": "2025-12-30T14:30:00Z",
        "unit": "requests/second"
      },
      "latency": {
        "avg_ms": 45.5,
        "p50_ms": 42.3,
        "p95_ms": 68.2,
        "p99_ms": 95.8,
        "max_ms": 150
      },
      "error_rate": {
        "percentage": 0.12,
        "total_requests": 100000,
        "error_count": 120
      },
      "concurrent_users": {
        "current": 12,
        "max": 20
      }
    },
    "timeline": [
      {
        "timestamp": "2025-12-30T10:00:00Z",
        "qps": 25.3,
        "avg_latency_ms": 38.2
      },
      {
        "timestamp": "2025-12-30T11:00:00Z",
        "qps": 52.1,
        "avg_latency_ms": 48.5
      }
    ]
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_010"
}
```

---

### 4.3 获取资源使用趋势

获取资源使用的趋势数据。

**接口描述**: 查询资源使用的历史趋势

**请求方式**: `GET /tenants/current/resource-trends`

**权限要求**: 需要认证,租户管理员

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| resource_type | string | ✅ | 资源类型(database/cache/storage) |
| metric | string | ✅ | 指标(size/count/hit_rate) |
| period | string | ❌ | 统计周期(7d/30d/90d,默认30d) |
| granularity | string | ❌ | 时间粒度(day/week,默认day) |

**请求示例**:
```http
GET /api/v1/tenants/current/resource-trends?resource_type=database&metric=size&period=30d&granularity=day
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "resource_type": "database",
    "metric": "size",
    "period": "30d",
    "granularity": "day",
    "trends": [
      {
        "date": "2025-12-01",
        "value": 1.8,
        "unit": "GB"
      },
      {
        "date": "2025-12-02",
        "value": 1.85,
        "unit": "GB"
      },
      {
        "date": "2025-12-30",
        "value": 2.01,
        "unit": "GB"
      }
    ],
    "statistics": {
      "start_value": 1.5,
      "end_value": 2.01,
      "growth_rate": 34,
      "predicted_next_month": 2.7
    }
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_011"
}
```

---

## 5. 租户迁移API

### 5.1 创建数据导出任务

导出租户数据。

**接口描述**: 创建数据导出任务,支持全量导出和增量导出

**请求方式**: `POST /tenants/current/data-export`

**权限要求**: 需要认证,租户管理员

**请求体**:
```json
{
  "export_type": "full",
  "tables": ["users", "bots", "conversations"],
  "format": "sql",
  "compression": true,
  "include_attachments": true
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| export_type | string | ✅ | 导出类型(full增量) |
| tables | string[] | ✅ | 要导出的表 |
| format | string | ❌ | 导出格式(sql/csv/json,默认sql) |
| compression | boolean | ❌ | 是否压缩(默认true) |
| include_attachments | boolean | ❌ | 是否包含附件(默认false) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "导出任务已创建",
  "data": {
    "export_id": "export-20251230-001",
    "status": "processing",
    "estimated_size_gb": 2.1,
    "estimated_duration": "15-30分钟",
    "download_url_expires_at": "2025-12-31T17:00:00Z",
    "created_at": "2025-12-30T18:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_012"
}
```

---

### 5.2 查询导出任务状态

查询数据导出任务的执行状态。

**接口描述**: 查询导出任务的进度和下载链接

**请求方式**: `GET /tenants/current/data-export/{export_id}`

**权限要求**: 需要认证,租户管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| export_id | string | 导出任务ID |

**请求示例**:
```http
GET /api/v1/tenants/current/data-export/export-20251230-001
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "export_id": "export-20251230-001",
    "status": "completed",
    "progress": {
      "percentage": 100,
      "completed_tables": ["users", "bots", "conversations"],
      "total_tables": 3
    },
    "result": {
      "file_name": "tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890-export-20251230.sql.gz",
      "file_size_mb": 525.3,
      "download_url": "https://export.saas.coze.com/exports/export-20251230-001.sql.gz",
      "expires_at": "2025-12-31T17:00:00Z"
    },
    "completed_at": "2025-12-30T18:25:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_013"
}
```

---

### 5.3 创建数据导入任务

导入租户数据。

**接口描述** 创建数据导入任务,用于数据恢复或迁移

**请求方式**: `POST /tenants/current/data-import`

**权限要求**: 需要认证,租户管理员

**请求体**:
```json
{
  "source_url": "https://export.saas.coze.com/exports/export-20251230-001.sql.gz",
  "import_mode": "skip_existing",
  "validation_mode": "strict",
  "notify_on_completion": true
}
```

**字段说明**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| source_url | string | ✅ | 导入数据源URL |
| import_mode | string | ✅ | 导入模式(skip_existing/overwrite/merge) |
| validation_mode | string | ✅ | 验证模式(strict/lenient) |
| notify_on_completion | boolean | ❌ | 完成后通知(默认true) |

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "导入任务已创建",
  "data": {
    "import_id": "import-20251230-001",
    "status": "pending",
    "source_url": "https://export.saas.coze.com/exports/export-20251230-001.sql.gz",
    "estimated_duration": "20-40分钟",
    "created_at": "2025-12-30T18:30:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_014"
}
```

---

### 5.4 查询导入任务状态

查询数据导入任务的执行状态。

**接口描述**: 查询导入任务的进度和结果

**请求方式**: `GET /tenants/current/data-import/{import_id}`

**权限要求**: 需要认证,租户管理员

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| import_id | string | 导入任务ID |

**请求示例**:
```http
GET /api/v1/tenants/current/data-import/import-20251230-001
```

**成功响应** (200 OK):
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "import_id": "import-20251230-001",
    "status": "in_progress",
    "progress": {
      "percentage": 45,
      "current_step": "正在导入数据...",
      "completed_tables": ["users", "bots"],
      "total_tables": 5,
      "imported_records": 580,
      "total_records": 1280
    },
    "validation": {
      "validated_records": 580,
      "skipped_records": 0,
      "error_records": 0
    },
    "started_at": "2025-12-30T18:35:00Z",
    "estimated_completion_at": "2025-12-30T19:00:00Z"
  },
  "timestamp": 1735574400000,
  "trace_id": "trace_015"
}
```

---

## 6. 错误代码

| 错误代码 | HTTP状态码 | 说明 |
|---------|-----------|------|
| 40001 | 400 | 租户规模不满足升级条件 |
| 40002 | 400 | 导入数据验证失败 |
| 40003 | 400 | 导出数据格式不支持 |
| 40101 | 401 | 未认证 |
| 40301 | 403 | 无权限(需要系统管理员) |
| 40401 | 404 | 租户不存在 |
| 40901 | 409 | 隔离策略迁移任务已存在 |
| 42901 | 429 | 导出/导入任务并发数超限 |
| 50001 | 500 | 服务器内部错误 |
| 50002 | 500 | 数据库连接失败 |
| 50003 | 500 | 迁移任务失败 |

---

## 7. 代码示例

### 7.1 Go 后端实现示例

```go
// service/isolation_service.go
package service

import (
    "context"
    "gorm.io/gorm"
)

type IsolationService struct {
    db          *gorm.DB
    tenantService *TenantService
    migrationSvc  *MigrationService
}

// GetIsolationStrategy 获取租户隔离策略
func (s *IsolationService) GetIsolationStrategy(ctx context.Context, tenantID string) (*IsolationStrategyInfo, error) {
    tenant, err := s.tenantService.GetByID(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 判断租户规模
    size := s.determineTenantSize(ctx, tenantID)

    // 根据规模选择隔离策略
    strategy := s.selectIsolationStrategy(size)

    return &IsolationStrategyInfo{
        TenantID:          tenantID,
        IsolationStrategy: strategy,
        TenantSize:        size,
        ResourceAllocation: s.getResourceAllocation(tenantID, strategy),
        UpgradeThreshold:  s.getUpgradeThreshold(strategy),
        CanUpgrade:        s.canUpgrade(size),
    }, nil
}

// UpgradeStrategy 升级隔离策略
func (s *IsolationService) UpgradeStrategy(ctx context.Context, tenantID string, targetStrategy string) (*MigrationTask, error) {
    tenant, err := s.tenantService.GetByID(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 检查是否可以升级
    if !s.canUpgradeTo(tenant, targetStrategy) {
        return nil, errors.New("不满足升级条件")
    }

    // 创建迁移任务
    task := &MigrationTask{
        ID:             fmt.Sprintf("migration-%s", uuid.New().String()),
        TenantID:       tenantID,
        CurrentStrategy: tenant.IsolationStrategy,
        TargetStrategy: targetStrategy,
        Status:         "pending",
        CreatedAt:      time.Now(),
    }

    if err := s.db.Create(task).Error; err != nil {
        return nil, err
    }

    // 异步执行迁移
    go s.migrationSvc.ExecuteMigration(ctx, task)

    return task, nil
}

// VerifyDataIsolation 验证数据隔离
func (s *IsolationService) VerifyDataIsolation(ctx context.Context, tenantID string, tables []string) (*VerificationResult, error) {
    result := &VerificationResult{
        TenantID:     tenantID,
        VerifiedAt:   time.Now(),
    }

    passed := true
    details := make(map[string]*TableVerificationResult)

    for _, table := range tables {
        // 检查是否有跨租户数据泄露
        leakCount, err := s.checkCrossTenantLeak(ctx, tenantID, table)
        if err != nil {
            return nil, err
        }

        tableResult := &TableVerificationResult{
            TableName:             table,
            Status:                "passed",
            TotalRecords:          s.getTotalRecords(ctx, table),
            TenantRecords:        s.getTenantRecords(ctx, tenantID, table),
            OtherTenantRecords:    leakCount,
            IsolationPercentage: 100,
        }

        if leakCount > 0 {
            passed = false
            tableResult.Status = "failed"
            tableResult.IsolationPercentage = calculateIsolationPercentage(tableResult)
        }

        details[table] = tableResult
    }

    result.VerificationResult = map[bool]string{true: "passed", false: "failed"}[passed]
    result.Details = details

    return result, nil
}
```

### 7.2 前端 TypeScript 类型定义

```typescript
// types/isolation.ts

interface IsolationStrategy {
  type: 'row_level' | 'schema_level' | 'database_level';
  name: string;
  description: string;
}

interface TenantSize {
  size: 'small' | 'medium' | 'large' | 'enterprise';
  user_count: number;
  bot_count: number;
  storage_gb: number;
}

interface ResourceUsage {
  database: {
    isolation_strategy: IsolationStrategy;
    shared_with: number;
    table_sizes: Record<string, string>;
    total_size: string;
    connection_pool: {
      active_connections: number;
      idle_connections: number;
      max_connections: number;
    };
    query_performance: {
      avg_query_time_ms: number;
      slow_queries_today: number;
    };
  };
  cache: {
    namespace: string;
    memory_usage: string;
    key_count: number;
    hit_rate: number;
    eviction_count: number;
  };
  storage: {
    bucket: string;
    total_size_gb: number;
    object_count: number;
    largest_object: string;
    upload_count_today: number;
  };
}

interface MigrationTask {
  task_id: string;
  tenant_id: string;
  current_strategy: IsolationStrategy;
  target_strategy: IsolationStrategy;
  status: 'pending' | 'in_progress' | 'completed' | 'failed';
  progress?: {
    percentage: number;
    current_step: string;
    completed_steps: string[];
    total_steps: string[];
  };
  created_at: string;
  estimated_duration: string;
  estimated_completion_at?: string;
}

// API客户端
class IsolationApiClient {
  private baseURL = '/api/v1';

  async getIsolationStrategy(): Promise<IsolationStrategyInfo> {
    const response = await apiClient.get(`${this.baseURL}/tenants/current/isolation-strategy`);
    return response.data.data;
  }

  async upgradeStrategy(tenantId: string, targetStrategy: string, scheduledAt?: string): Promise<MigrationTask> {
    const response = await apiClient.post(
      `${this.baseURL}/admin/tenants/${tenantId}/isolation-strategy/upgrade`,
      {
        target_strategy: targetStrategy,
        scheduled_at: scheduledAt,
      }
    );
    return response.data.data;
  }

  async verifyDataIsolation(tenantId: string, tables?: string[]): Promise<VerificationResult> {
    const response = await apiClient.post(
      `${this.baseURL}/admin/tenants/${tenantId}/isolation/verify`,
      {
        tables: tables || ['users', 'bots', 'conversations'],
      }
    );
    return response.data.data;
  }

  async getResourceUsage(resourceType?: string): Promise<ResourceUsage> {
    const response = await apiClient.get(
      `${this.baseURL}/tenants/current/resource-usage`,
      { params: resourceType ? { resource_type: resourceType } : undefined }
    );
    return response.data.data;
  }

  async getMetrics(period = '24h', metricType?: string): Promise<MetricsData> {
    const response = await apiClient.get(
      `${this.baseURL}/tenants/current/metrics`,
      {
        params: { period, metric_type: metricType }
      }
    );
    return response.data.data;
  }
}

export const isolationApi = new IsolationApiClient();
```

---

**文档结束**

---

## 📊 文档统计

- **接口总数**: 10+
- **核心功能模块**: 5个
- **代码示例**: Go + TypeScript
- **错误代码**: 8个

**主要功能模块**:
1. 隔离策略管理 (3个接口)
2. 数据隔离验证 (3个接口)
3. 租户资源监控 (3个接口)
4. 租户迁移 (4个接口)
