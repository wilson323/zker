# 租户监控运维系统 (Tenant Monitoring & Operations)

## 📖 简介

完整的**企业级租户监控运维系统**，提供实时监控、智能告警、性能分析等能力。

### 核心功能

- ✅ **租户级监控**：QPS、响应时间、错误率、资源使用率
- ✅ **智能告警**：可配置规则、多通知渠道、告警确认/解决
- ✅ **健康评分**：0-100分健康度评估
- ✅ **Dashboard**：Grafana可视化大盘
- ✅ **Prometheus集成**：标准指标采集

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────────┐
│                       Frontend (React)                        │
│  ┌──────────────────┐  ┌──────────────────┐                 │
│  │ Monitoring Dashboard│  │ Alert Management │                 │
│  └──────────────────┘  └──────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      API Layer (Go)                          │
│  ┌───────────────────────────────────────────────────┐       │
│  │ MonitoringHandler                                 │       │
│  │  - GetTenantOverview                              │       │
│  │  - GetQPSMetrics                                  │       │
│  │  - CreateAlertRule                                │       │
│  └───────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     Service Layer (Go)                       │
│  ┌──────────────────────┐  ┌──────────────────────┐         │
│  │ TenantMonitoringSvc   │  │ AlertService          │         │
│  │  - CollectQPS         │  │  - CheckAlerts        │         │
│  │  - CollectResponseTime│  │  - SendAlert          │         │
│  │  - CalculateHealth    │  │  - AcknowledgeAlert   │         │
│  └──────────────────────┘  └──────────────────────┘         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer (Go)                      │
│  ┌──────────────────┐  ┌──────────────────┐                 │
│  │ TenantMetricsRepo│  │ AlertHistoryRepo │                 │
│  └──────────────────┘  └──────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Database (MySQL 8.4.5)                    │
│  tenant_metrics | agent_metrics | alert_rules               │
│  alert_history | performance_reports                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Monitoring Stack                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Prometheus   │  │ Grafana      │  │ AlertManager │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

## 📦 目录结构

```
backend/
├── domain/monitoring/
│   ├── entity/              # 实体定义
│   │   ├── tenant_metrics.go
│   │   ├── agent_metrics.go
│   │   └── alert.go
│   ├── repository/          # 仓储层
│   │   ├── tenant_metrics_repository.go
│   │   └── tenant_metrics_repository_impl.go
│   ├── service/             # 服务层
│   │   ├── tenant_monitoring_service.go
│   │   ├── alert_service.go
│   │   └── notification_service.go
│   └── migration/sql/       # 数据库迁移
│       └── create_monitoring_tables.sql
├── infra/monitoring/
│   ├── collector/           # 指标采集
│   │   └── metrics_collector.go
│   └── middleware/          # 中间件
│       └── prometheus_middleware.go
└── api/handler/coze/monitoring/
    └── monitoring_handler.go

frontend/apps/coze-studio/src/
├── pages/monitoring/
│   ├── TenantMonitoringDashboard.tsx
│   └── AlertManagementPage.tsx
├── types/monitoring.ts      # TypeScript类型
└── api/monitoring/index.ts  # API客户端

deploy/monitoring/
└── grafana/dashboards/
    └── tenant_monitoring_dashboard.json
```

## 🚀 快速开始

### 1. 数据库初始化

```bash
# 执行SQL脚本
mysql -u root -p coze_studio < backend/domain/monitoring/migration/sql/create_monitoring_tables.sql
```

### 2. 启动后端服务

```bash
# 设置环境变量
export MONITORING_ENABLED=true
export PROMETHEUS_ENABLED=true

# 启动服务
cd backend
go run cmd/server/main.go
```

### 3. 启动前端服务

```bash
cd frontend
npm run dev
```

### 4. 配置Prometheus

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'coze-studio'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s
```

### 5. 导入Grafana Dashboard

1. 登录Grafana (http://localhost:3000)
2. Dashboards -> Import
3. 上传 `deploy/monitoring/grafana/dashboards/tenant_monitoring_dashboard.json`

## 📊 API文档

### 租户监控

```
GET  /api/v1/monitoring/tenants/:tenant_id/overview          # 租户概览
GET  /api/v1/monitoring/tenants/:tenant_id/metrics/qps      # QPS指标
GET  /api/v1/monitoring/tenants/:tenant_id/metrics/response_time  # 响应时间
GET  /api/v1/monitoring/tenants/:tenant_id/metrics/error_rate      # 错误率
GET  /api/v1/monitoring/tenants/:tenant_id/health_score    # 健康评分
```

### 告警管理

```
POST /api/v1/monitoring/alert_rules                         # 创建告警规则
GET  /api/v1/monitoring/alerts/history                       # 告警历史
POST /api/v1/monitoring/alerts/:id/acknowledge               # 确认告警
POST /api/v1/monitoring/alerts/:id/resolve                   # 解决告警
POST /api/v1/monitoring/alerts/:id/silence                   # 静默告警
GET  /api/v1/monitoring/alerts/statistics                    # 告警统计
```

### 指标上报

```
POST /api/v1/monitoring/metrics/record                       # 记录单个指标
POST /api/v1/monitoring/metrics/batch                        # 批量记录指标
```

## 🔧 配置说明

### 告警规则配置

```json
{
  "tenant_id": "tenant-123",
  "rule_name": "高错误率告警",
  "metric_type": "error_rate",
  "threshold": 0.05,
  "comparison": "gt",
  "severity": "critical",
  "notification_channels": ["email", "webhook"],
  "evaluation_interval": 60,
  "for_duration": 300
}
```

### 通知渠道配置

```json
{
  "email_recipients": ["admin@example.com"],
  "sms_recipients": ["+8613800138000"],
  "webhook_url": "https://hooks.example.com/alerts",
  "headers": {
    "Authorization": "Bearer your-token"
  }
}
```

## 📈 监控指标

### Prometheus指标

```
coze_http_requests_total          # HTTP请求总数
coze_http_request_duration_seconds # 请求持续时间
coze_response_time_milliseconds   # 响应时间
coze_qps                          # QPS
coze_error_rate                   # 错误率
coze_health_score                 # 健康评分
```

### 数据库指标

- `tenant_metrics`: 租户级指标
- `agent_metrics`: Agent级指标
- `alert_history`: 告警历史

## 🧪 测试

### 单元测试

```bash
# Service层测试
cd backend/domain/monitoring/service
go test -v -cover

# Repository层测试
cd backend/domain/monitoring/repository
go test -v -cover
```

### 集成测试

```bash
# 使用testcontainers运行集成测试
cd backend
go test ./domain/monitoring/... -tags=integration
```

## 📚 文档

- [实现总结](../../docs/企业级功能完善与统一性设计方案/ZKER-租户监控运维系统-实现总结.md)
- [数据库设计](../../docs/企业级功能完善与统一性设计方案/数据库设计完整交付清单.md)
- [API规范](../../docs/企业级功能完善与统一性设计方案/API设计规范文档.md)

## 🔒 安全

- ✅ 租户隔离（所有表带tenant_id）
- ✅ API认证（JWT）
- ✅ 敏感数据加密
- ✅ SQL注入防护（GORM参数化查询）
- ✅ 访问控制（RBAC）

## 🚧 待完成

- [ ] Agent监控完善
- [ ] 性能报告生成
- [ ] 单元测试（覆盖率≥80%）
- [ ] 通知服务增强（邮件/短信集成）
- [ ] 监控指标丰富（CPU/内存/磁盘）

## 📝 更新日志

### v1.0.0 (2025-01-01)

- ✅ 初始版本发布
- ✅ 租户监控核心功能
- ✅ 告警系统
- ✅ Prometheus集成
- ✅ Grafana Dashboard
- ✅ 前端React组件

## 👥 贡献

欢迎提交PR和Issue！

## 📄 许可证

Apache License 2.0

---

**维护者**: ZKER研发团队
**联系方式**: dev@zker.com
