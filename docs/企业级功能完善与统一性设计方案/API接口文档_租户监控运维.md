# API接口文档：租户监控运维模块

**模块名称**: 租户监控运维 (TenantMonitoring)
**设计文档**: 24-MultiTenant_SaaS核心_租户监控运维.md
**版本**: v1.0.0
**最后更新**: 2025-01-03
**优先级**: P2

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 监控指标API](#2-监控指标api)
- [3. 告警规则API](#3-告警规则api)
- [4. 告警历史API](#4-告警历史api)
- [5. 监控大盘API](#5-监控大盘api)
- [6. 数据模型](#6-数据模型)
- [7. 错误码定义](#7-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

租户监控运维是MultiTenant SaaS的核心运维模块，通过**多维度监控 + 自动告警**，确保多租户系统的稳定性、性能和安全性：

- ✅ **性能监控** - API响应时间(P50/P95/P99)、错误率、并发数、队列长度
- ✅ **资源监控** - CPU使用率、内存使用率、磁盘IO、网络带宽
- ✅ **业务监控** - 活跃用户数、API调用量、Token消耗量、Bot使用量
- ✅ **告警管理** - 阈值告警、趋势告警、异常检测
- ✅ **通知渠道** - 邮件、短信、Webhook、企业微信/钉钉
- ✅ **可视化大盘** - Grafana集成，租户级监控大盘

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 数据库: MySQL 8.4.5 (告警规则、告警历史)
- 时序数据库: InfluxDB (监控指标存储)
- 可视化: Grafana
- 实现: 50% 编码（监控框架） + 50% 配置（告警规则）

**实现策略**: ✅ 50% 编码 + 50% 配置

### 1.3 数据库表

| 表名 | 说明 | 存储引擎 |
|------|------|---------|
| `monitoring_metrics` | 监控指标表 | InfluxDB (时序数据库) |
| `alert_rules` | 告警规则表 | MySQL |
| `alert_history` | 告警历史表 | MySQL |

---

## 2. 监控指标API

### 2.1 查询监控指标

**接口地址**: `GET /api/v1/monitoring/metrics`

**查询参数**:
- tenant_id: 租户ID (可选，不填则查询所有租户)
- metric_name: 指标名称 (可选)
- start_time: 开始时间戳 (必填)
- end_time: 结束时间戳 (必填)
- aggregation: 聚合方式 (可选，avg/sum/max/min/count，默认avg)
- interval: 时间间隔 (可选，如1m, 5m, 1h, 1d，默认5m)
- group_by: 分组字段 (可选，如tenant_id, metric_name)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "query": {
      "tenant_id": "tenant-001",
      "metric_name": "api_calls",
      "start_time": "2025-01-15T00:00:00Z",
      "end_time": "2025-01-15T23:59:59Z",
      "aggregation": "sum",
      "interval": "1h"
    },
    "series": [
      {
        "tenant_id": "tenant-001",
        "metric_name": "api_calls",
        "tags": {
          "endpoint": "/api/v1/bots",
          "method": "GET"
        },
        "datapoints": [
          {
            "timestamp": "2025-01-15T00:00:00Z",
            "value": 1250
          },
          {
            "timestamp": "2025-01-15T01:00:00Z",
            "value": 1180
          },
          {
            "timestamp": "2025-01-15T02:00:00Z",
            "value": 1320
          }
        ]
      }
    ],
    "summary": {
      "total": 15200,
      "avg": 633.33,
      "max": 1320,
      "min": 1180
    }
  }
}
```

### 2.2 获取实时指标

**接口地址**: `GET /api/v1/monitoring/metrics/realtime`

**查询参数**:
- tenant_id: 租户ID (可选)
- metric_names: 指标名称列表 (可选，逗号分隔)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "timestamp": "2025-01-15T10:30:00Z",
    "metrics": [
      {
        "tenant_id": "tenant-001",
        "metric_name": "api_response_time_p50",
        "value": 125.5,
        "unit": "ms"
      },
      {
        "tenant_id": "tenant-001",
        "metric_name": "api_response_time_p95",
        "value": 320.8,
        "unit": "ms"
      },
      {
        "tenant_id": "tenant-001",
        "metric_name": "api_response_time_p99",
        "value": 850.2,
        "unit": "ms"
      },
      {
        "tenant_id": "tenant-001",
        "metric_name": "api_error_rate",
        "value": 0.012,
        "unit": "%"
      },
      {
        "tenant_id": "tenant-001",
        "metric_name": "active_users",
        "value": 320,
        "unit": "count"
      }
    ]
  }
}
```

### 2.3 获取租户指标概览

**接口地址**: `GET /api/v1/monitoring/tenants/{tenant_id}/overview`

**查询参数**:
- start_time: 开始时间 (可选，默认24小时前)
- end_time: 结束时间 (可选，默认当前时间)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "tenant_id": "tenant-001",
    "tenant_name": "示例企业",
    "period": {
      "start": "2025-01-14T10:30:00Z",
      "end": "2025-01-15T10:30:00Z"
    },
    "performance": {
      "api_calls": 152000,
      "avg_response_time": 128.5,
      "p95_response_time": 320.8,
      "p99_response_time": 850.2,
      "error_rate": 0.015
    },
    "resources": {
      "cpu_usage_percent": 45.2,
      "memory_usage_percent": 62.8,
      "disk_usage_percent": 55.3,
      "network_in_mb": 1250.5,
      "network_out_mb": 2100.8
    },
    "business": {
      "active_users": 1250,
      "bot_calls": 85000,
      "token_usage": 12500000,
      "knowledge_searches": 3200
    }
  }
}
```

### 2.4 获取指标列表

**接口地址**: `GET /api/v1/monitoring/metrics/catalog`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "categories": [
      {
        "name": "性能指标",
        "metrics": [
          {
            "name": "api_response_time_p50",
            "display_name": "API响应时间P50",
            "unit": "ms",
            "description": "50%请求的响应时间"
          },
          {
            "name": "api_response_time_p95",
            "display_name": "API响应时间P95",
            "unit": "ms",
            "description": "95%请求的响应时间"
          },
          {
            "name": "api_response_time_p99",
            "display_name": "API响应时间P99",
            "unit": "ms",
            "description": "99%请求的响应时间"
          },
          {
            "name": "api_error_rate",
            "display_name": "API错误率",
            "unit": "%",
            "description": "API请求错误率"
          }
        ]
      },
      {
        "name": "资源指标",
        "metrics": [
          {
            "name": "cpu_usage_percent",
            "display_name": "CPU使用率",
            "unit": "%",
            "description": "租户CPU使用率"
          },
          {
            "name": "memory_usage_percent",
            "display_name": "内存使用率",
            "unit": "%",
            "description": "租户内存使用率"
          }
        ]
      },
      {
        "name": "业务指标",
        "metrics": [
          {
            "name": "active_users",
            "display_name": "活跃用户数",
            "unit": "count",
            "description": "活跃用户数量"
          },
          {
            "name": "bot_calls",
            "display_name": "Bot调用次数",
            "unit": "count",
            "description": "Bot API调用次数"
          }
        ]
      }
    ]
  }
}
```

---

## 3. 告警规则API

### 3.1 获取告警规则列表

**接口地址**: `GET /api/v1/monitoring/alert-rules`

**查询参数**:
- tenant_id: 租户ID (可选，NULL表示平台级规则)
- metric_name: 监控指标 (可选)
- severity: 严重程度 (可选，info/warning/critical)
- is_active: 是否启用 (可选)
- page: 页码 (默认1)
- page_size: 每页数量 (默认20)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": 1,
        "tenant_id": null,
        "name": "API响应时间过长",
        "description": "当P99响应时间超过5秒时触发告警",
        "metric_name": "api_response_time_p99",
        "condition": "gt",
        "threshold": 5000,
        "duration": 300,
        "severity": "critical",
        "notification_channels": {
          "email": ["ops@zker.com"],
          "webhook": ["https://hooks.example.com/critical"]
        },
        "is_active": true,
        "created_at": "2025-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "tenant_id": "tenant-001",
        "name": "Bot调用次数异常",
        "description": "每分钟Bot调用次数超过1000次时触发告警",
        "metric_name": "bot_calls_per_minute",
        "condition": "gt",
        "threshold": 1000,
        "duration": 60,
        "severity": "warning",
        "notification_channels": {
          "wechat": {
            "corpId": "ww12345",
            "agentId": 123456,
            "toUser": ["admin"]
          }
        },
        "is_active": true,
        "created_at": "2025-01-10T00:00:00Z"
      }
    ]
  }
}
```

### 3.2 创建告警规则

**接口地址**: `POST /api/v1/monitoring/alert-rules`

**权限**: 管理员或租户管理员

**请求参数**:
```json
{
  "name": "内存使用率过高",
  "description": "当内存使用率超过90%时触发告警",
  "metric_name": "memory_usage_percent",
  "condition": "gt",
  "threshold": 90,
  "duration": 300,
  "severity": "warning",
  "notification_channels": {
    "email": ["admin@example.com"],
    "wechat": {
      "corpId": "ww12345",
      "agentId": 123456,
      "toUser": ["admin"]
    }
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 16,
    "name": "内存使用率过高",
    "created_at": "2025-01-15T14:00:00Z"
  }
}
```

### 3.3 更新告警规则

**接口地址**: `PUT /api/v1/monitoring/alert-rules/{rule_id}`

**权限**: 仅规则创建者或管理员

**请求参数**:
```json
{
  "name": "内存使用率过高(更新)",
  "description": "当内存使用率超过85%时触发告警",
  "threshold": 85,
  "is_active": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 16,
    "updated_at": "2025-01-15T15:00:00Z"
  }
}
```

### 3.4 删除告警规则

**接口地址**: `DELETE /api/v1/monitoring/alert-rules/{rule_id}`

**权限**: 仅规则创建者或管理员

**响应示例**:
```json
{
  "code": 0,
  "message": "删除成功",
  "data": {
    "deleted_rule_id": 16
  }
}
```

### 3.5 获取告警规则详情

**接口地址**: `GET /api/v1/monitoring/alert-rules/{rule_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "tenant_id": null,
    "name": "API响应时间过长",
    "description": "当P99响应时间超过5秒时触发告警",
    "metric_name": "api_response_time_p99",
    "condition": "gt",
    "threshold": 5000,
    "duration": 300,
    "severity": "critical",
    "notification_channels": {
      "email": ["ops@zker.com"],
      "webhook": ["https://hooks.example.com/critical"],
      "sms": ["+86138xxxxxxxx"]
    },
    "is_active": true,
    "triggered_count": 15,
    "last_triggered_at": "2025-01-15T10:00:00Z",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-10T00:00:00Z"
  }
}
```

---

## 4. 告警历史API

### 4.1 获取告警历史列表

**接口地址**: `GET /api/v1/monitoring/alerts`

**查询参数**:
- tenant_id: 租户ID (可选)
- rule_id: 告警规则ID (可选)
- severity: 严重程度 (可选)
- status: 状态 (可选，fired/resolved/acknowledged)
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
    "total": 125,
    "items": [
      {
        "id": 1001,
        "rule_id": 1,
        "rule_name": "API响应时间过长",
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "metric_name": "api_response_time_p99",
        "actual_value": 5250.5,
        "threshold": 5000,
        "severity": "critical",
        "status": "fired",
        "fired_at": "2025-01-15T10:00:00Z",
        "resolved_at": null,
        "acknowledged_by": null,
        "duration_minutes": 30
      },
      {
        "id": 1002,
        "rule_id": 2,
        "rule_name": "Bot调用次数异常",
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "metric_name": "bot_calls_per_minute",
        "actual_value": 1250,
        "threshold": 1000,
        "severity": "warning",
        "status": "resolved",
        "fired_at": "2025-01-15T09:00:00Z",
        "resolved_at": "2025-01-15T09:30:00Z",
        "acknowledged_by": 1001,
        "acknowledged_by_name": "管理员",
        "duration_minutes": 30
      }
    ]
  }
}
```

### 4.2 确认告警

**接口地址**: `POST /api/v1/monitoring/alerts/{alert_id}/acknowledge`

**请求参数**:
```json
{
  "comment": "正在处理，已通知相关团队"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "确认成功",
  "data": {
    "id": 1001,
    "status": "acknowledged",
    "acknowledged_by": 1001,
    "acknowledged_by_name": "管理员",
    "acknowledged_at": "2025-01-15T10:30:00Z",
    "comment": "正在处理，已通知相关团队"
  }
}
```

### 4.3 解决告警

**接口地址**: `POST /api/v1/monitoring/alerts/{alert_id}/resolve`

**请求参数**:
```json
{
  "comment": "问题已修复，API响应时间恢复正常"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "告警已解决",
  "data": {
    "id": 1001,
    "status": "resolved",
    "resolved_at": "2025-01-15T11:00:00Z",
    "comment": "问题已修复，API响应时间恢复正常"
  }
}
```

### 4.4 获取告警统计

**接口地址**: `GET /api/v1/monitoring/alerts/statistics`

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
      "total": 1250,
      "by_severity": {
        "critical": 45,
        "warning": 520,
        "info": 685
      },
      "by_status": {
        "fired": 120,
        "acknowledged": 80,
        "resolved": 1050
      },
      "avg_resolution_time_minutes": 45.5
    },
    "top_rules": [
      {
        "rule_id": 1,
        "rule_name": "API响应时间过长",
        "triggered_count": 150
      },
      {
        "rule_id": 2,
        "rule_name": "Bot调用次数异常",
        "triggered_count": 125
      }
    ],
    "top_tenants": [
      {
        "tenant_id": "tenant-001",
        "tenant_name": "示例企业",
        "alert_count": 85
      }
    ]
  }
}
```

---

## 5. 监控大盘API

### 5.1 获取大盘数据

**接口地址**: `GET /api/v1/monitoring/dashboard`

**查询参数**:
- tenant_id: 租户ID (可选，不填则显示平台级大盘)
- start_time: 开始时间 (可选，默认24小时前)
- end_time: 结束时间 (可选，默认当前时间)

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": {
      "start": "2025-01-14T10:30:00Z",
      "end": "2025-01-15T10:30:00Z"
    },
    "overview": {
      "total_tenants": 150,
      "active_tenants": 125,
      "total_api_calls": 15200000,
      "avg_response_time": 128.5,
      "overall_error_rate": 0.012
    },
    "metrics": {
      "api_calls": {
        "current": 15200000,
        "previous": 14500000,
        "change": "+4.83%"
      },
      "active_users": {
        "current": 12500,
        "previous": 11800,
        "change": "+5.93%"
      },
      "avg_response_time": {
        "current": 128.5,
        "previous": 135.2,
        "change": "-4.95%"
      },
      "error_rate": {
        "current": 0.012,
        "previous": 0.015,
        "change": "-20.00%"
      }
    },
    "top_tenants": {
      "by_api_calls": [
        {
          "tenant_id": "tenant-001",
          "tenant_name": "示例企业",
          "api_calls": 1520000
        }
      ],
      "by_active_users": [
        {
          "tenant_id": "tenant-001",
          "tenant_name": "示例企业",
          "active_users": 1250
        }
      ]
    },
    "alerts": {
      "fired": 25,
      "acknowledged": 15,
      "resolved": 350,
      "critical": 5,
      "warning": 20
    }
  }
}
```

### 5.2 导出监控报告

**接口地址**: `POST /api/v1/monitoring/export`

**请求参数**:
```json
{
  "tenant_id": "tenant-001",
  "start_time": "2025-01-01T00:00:00Z",
  "end_time": "2025-01-15T23:59:59Z",
  "format": "pdf",
  "include_metrics": ["performance", "resources", "business"],
  "include_alerts": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "报告生成中",
  "data": {
    "report_id": "report-20250115-001",
    "download_url": "/api/v1/monitoring/reports/report-20250115-001/download",
    "expires_at": "2025-01-16T00:00:00Z"
  }
}
```

---

## 6. 数据模型

### 6.1 MonitoringMetric

```typescript
interface MonitoringMetric {
  id: number;
  tenant_id: string;
  metric_name: string;
  metric_value: number;
  metric_tags?: Record<string, string>;
  timestamp: number;
}

interface MetricSeries {
  tenant_id?: string;
  metric_name: string;
  tags?: Record<string, string>;
  datapoints: Datapoint[];
}

interface Datapoint {
  timestamp: string;
  value: number;
}
```

### 6.2 AlertRule

```typescript
interface AlertRule {
  id: number;
  tenant_id?: string;
  name: string;
  description?: string;
  metric_name: string;
  condition: 'gt' | 'lt' | 'eq' | 'gte' | 'lte';
  threshold: number;
  duration: number;
  severity: 'info' | 'warning' | 'critical';
  notification_channels: NotificationChannels;
  is_active: boolean;
  triggered_count?: number;
  last_triggered_at?: Date;
  created_at: Date;
  updated_at: Date;
}

interface NotificationChannels {
  email?: string[];
  sms?: string[];
  webhook?: string[];
  wechat?: {
    corpId: string;
    agentId: number;
    toUser: string[];
  };
  dingtalk?: {
    webhook: string;
    atMobiles?: string[];
  };
}
```

### 6.3 Alert

```typescript
interface Alert {
  id: number;
  rule_id: number;
  rule_name: string;
  tenant_id?: string;
  tenant_name?: string;
  metric_name: string;
  actual_value: number;
  threshold: number;
  severity: 'info' | 'warning' | 'critical';
  status: 'fired' | 'resolved' | 'acknowledged';
  fired_at: Date;
  resolved_at?: Date;
  acknowledged_by?: number;
  acknowledged_by_name?: string;
  comment?: string;
  duration_minutes?: number;
}
```

### 6.4 Dashboard

```typescript
interface Dashboard {
  period: {
    start: string;
    end: string;
  };
  overview: {
    total_tenants: number;
    active_tenants: number;
    total_api_calls: number;
    avg_response_time: number;
    overall_error_rate: number;
  };
  metrics: {
    [key: string]: {
      current: number;
      previous: number;
      change: string;
    };
  };
  top_tenants: {
    by_api_calls: TopTenant[];
    by_active_users: TopTenant[];
  };
  alerts: {
    fired: number;
    acknowledged: number;
    resolved: number;
    critical: number;
    warning: number;
  };
}

interface TopTenant {
  tenant_id: string;
  tenant_name: string;
  api_calls?: number;
  active_users?: number;
}
```

---

## 7. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 36001 | 400 | 监控指标名称无效 |
| 36002 | 400 | 时间范围无效 |
| 36003 | 400 | 聚合方式无效 |
| 36004 | 400 | 时间间隔无效 |
| 36101 | 404 | 告警规则不存在 |
| 36102 | 400 | 告警规则名称已存在 |
| 36103 | 400 | 阈值无效 |
| 36104 | 400 | 条件无效 |
| 36105 | 400 | 通知渠道配置无效 |
| 36106 | 403 | 无权限操作告警规则 |
| 36201 | 404 | 告警不存在 |
| 36202 | 400 | 告警状态不允许此操作 |
| 36203 | 403 | 无权限操作告警 |
| 36301 | 400 | 报告格式无效 |
| 36302 | 500 | 报告生成失败 |

---

**文档结束**
