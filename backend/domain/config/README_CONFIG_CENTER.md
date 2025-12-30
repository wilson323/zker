# 配置中心集成指南

## 概述

本配置中心基于 etcd 实现分布式配置管理，支持多租户配置隔离、动态配置更新、版本控制和一键回滚。

## 架构设计

### 核心组件

```
┌─────────────────────────────────────────────────────┐
│                   API Layer                         │
│  /api/config/* (创建、查询、更新、删除、回滚)        │
└──────────────────┬──────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────┐
│              Service Layer                          │
│  ConfigService                                      │
│  - CreateConfig()                                   │
│  - GetConfig()                                      │
│  - UpdateConfig()                                   │
│  - DeleteConfig()                                   │
│  - RollbackConfig()                                 │
└──────────────────┬──────────────────────────────────┘
                   │
       ┌───────────┴───────────┐
       │                       │
┌──────▼────────┐    ┌────────▼────────┐
│  Repository   │    │  etcd Provider  │
│  (Database)   │    │  (Dynamic)      │
│               │    │                 │
│ - MySQL       │    │ - Watch         │
│ - History     │    │ - Sync          │
└───────────────┘    └─────────────────┘
```

### 数据流程

1. **配置读取**：优先从 etcd 获取（实时），失败回退到数据库
2. **配置更新**：同时更新数据库和 etcd
3. **变更通知**：etcd watcher 自动推送配置变更到所有服务实例
4. **版本管理**：每次更新自动记录历史，支持一键回滚

## 安装和配置

### 1. 安装 etcd

```bash
# Docker 方式
docker run -d \
  --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  -e ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 \
  bitnami/etcd:latest

# 或使用 docker-compose
cat > docker-compose.yml <<EOF
version: '3'
services:
  etcd:
    image: bitnami/etcd:latest
    container_name: etcd
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379
    ports:
      - "2379:2379"
      - "2380:2380"
EOF

docker-compose up -d
```

### 2. 初始化数据库

```bash
# 执行迁移脚本
mysql -u root -p database_name < backend/domain/config/migration/01_init_config_tables.sql
```

### 3. 配置应用

在 `backend/conf/.env` 中添加 etcd 配置：

```env
# etcd 配置
ETCD_ENDPOINTS=http://localhost:2379
ETCD_NAMESPACE=coze-studio
ETCD_GROUP=production

# 可选：TLS 配置
ETCD_TLS_ENABLED=false
ETCD_CERT_FILE=
ETCD_KEY_FILE=
ETCD_CA_FILE=

# 可选：认证配置
ETCD_USERNAME=
ETCD_PASSWORD=
```

### 4. 初始化配置服务

```go
package main

import (
    "context"
    "github.com/coze-dev/coze-studio/backend/domain/config/service"
    "github.com/coze-dev/coze-studio/backend/domain/config/repository"
    configrepository "github.com/coze-dev/coze-studio/backend/infra/config/repository"
    etcdclient "github.com/coze-dev/coze-studio/backend/infra/dynconf/impl/etcd"
)

func initConfigCenter() {
    // 1. 创建 etcd 客户端
    provider := etcdclient.NewProvider([]string{"http://localhost:2379"})
    etcdCli, err := provider.Initialize(
        context.Background(),
        "coze-studio",  // namespace
        "production",   // group
    )
    if err != nil {
        log.Fatalf("init etcd failed: %v", err)
    }

    // 2. 创建仓储
    configRepo := configrepository.NewConfigRepository(db)
    historyRepo := configrepository.NewConfigHistoryRepository(db)

    // 3. 创建配置服务
    configService := service.NewConfigService(configRepo, historyRepo, etcdCli)

    // 4. 初始化 API handler
    coze.InitConfigCenter(configService)
}
```

## API 使用指南

### 1. 创建配置

```bash
POST /api/config/create
Content-Type: application/json

{
  "tenant_id": "tenant_123",
  "config_key": "feature.flag.smart_routing",
  "config_value": "true",
  "config_type": "bool",
  "description": "智能路由功能开关"
}
```

**响应**：
```json
{
  "config_id": "config_abc123"
}
```

### 2. 查询配置

```bash
GET /api/config/get?tenant_id=tenant_123&config_key=feature.flag.smart_routing
```

**响应**：
```json
{
  "config_id": "config_abc123",
  "config_key": "feature.flag.smart_routing",
  "config_value": "true",
  "config_type": "bool",
  "description": "智能路由功能开关",
  "updated_at": "2025-01-01 12:00:00"
}
```

### 3. 更新配置

```bash
POST /api/config/update
Content-Type: application/json

{
  "tenant_id": "tenant_123",
  "config_key": "feature.flag.smart_routing",
  "new_value": "false",
  "change_reason": "临时关闭智能路由进行测试"
}
```

**响应**：
```json
{
  "success": true
}
```

### 4. 删除配置

```bash
POST /api/config/delete
Content-Type: application/json

{
  "tenant_id": "tenant_123",
  "config_key": "feature.flag.smart_routing"
}
```

### 5. 列出租户所有配置

```bash
GET /api/config/list?tenant_id=tenant_123
```

**响应**：
```json
{
  "configs": [
    {
      "config_id": "config_001",
      "config_key": "feature.flag.smart_routing",
      "config_value": "true",
      "config_type": "bool",
      "description": "智能路由功能开关",
      "updated_at": "2025-01-01 12:00:00"
    }
  ],
  "total": 1
}
```

### 6. 查询配置变更历史

```bash
GET /api/config/history?config_id=config_abc123&limit=10&offset=0
```

**响应**：
```json
{
  "history": [
    {
      "history_id": "hist_001",
      "config_key": "feature.flag.smart_routing",
      "old_value": "true",
      "new_value": "false",
      "change_reason": "临时关闭智能路由进行测试",
      "changed_by": "user_123",
      "changed_at": "2025-01-01 12:30:00",
      "version_number": 2
    }
  ],
  "total": 1
}
```

### 7. 回滚配置

```bash
POST /api/config/rollback
Content-Type: application/json

{
  "tenant_id": "tenant_123",
  "config_key": "feature.flag.smart_routing",
  "version": 1
}
```

**响应**：
```json
{
  "success": true
}
```

## 配置类型支持

支持的配置类型：

- `string`: 字符串类型
- `int`: 整数类型
- `float`: 浮点数类型
- `bool`: 布尔类型（true/false）
- `json`: JSON 对象（需要符合 JSON 格式）

## 使用场景

### 1. 功能开关（Feature Flags）

```json
{
  "tenant_id": "tenant_123",
  "config_key": "feature.flag.new_ui",
  "config_value": "true",
  "config_type": "bool"
}
```

### 2. 性能调优

```json
{
  "tenant_id": "tenant_123",
  "config_key": "performance.max_concurrent_requests",
  "config_value": "100",
  "config_type": "int"
}
```

### 3. 第三方服务配置

```json
{
  "tenant_id": "tenant_123",
  "config_key": "integration.openai.api_key",
  "config_value": "sk-xxxxx",
  "config_type": "string",
  "is_encrypted": true
}
```

### 4. 路由规则配置

```json
{
  "tenant_id": "tenant_123",
  "config_key": "routing.rules.intent_match",
  "config_value": "{\"enabled\":true,\"threshold\":0.8}",
  "config_type": "json"
}
```

## 最佳实践

### 1. 配置命名规范

使用分层命名，使用 `.` 分隔：

```
feature.flag.{feature_name}
performance.{metric_name}
integration.{service}.{property}
routing.{component}.{setting}
monitoring.{component}.{property}
```

### 2. 配置分类

- **系统配置**（tenant_id = 'system'）：全局默认配置
- **租户配置**（tenant_id = 具体租户ID）：租户特定配置
- **业务配置**：与业务逻辑相关的配置
- **基础设施配置**：数据库、缓存、消息队列等

### 3. 变更管理

1. **变更前备份**：所有变更自动记录历史
2. **变更原因**：提供清晰的变更原因
3. **灰度发布**：先在小范围测试配置变更
4. **监控验证**：变更后监控系统指标

### 4. 安全性

- **敏感配置加密**：API Key、密钥等设置 `is_encrypted=true`
- **权限控制**：通过 RBAC 控制配置读写权限
- **审计日志**：所有配置变更记录操作人和时间

## 监控和运维

### 关键指标

1. **配置命中率**：etcd 命中率应 > 99%
2. **配置更新频率**：异常高频更新可能预示问题
3. **配置同步延迟**：etcd 到数据库的同步延迟应 < 1s
4. **Watcher 数量**：监控活跃 watcher 连接数

### 告警规则

```yaml
# 配置中心异常告警
- alert: ConfigEtcdDown
  expr: up{job="etcd"} == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "etcd 服务不可用"

- alert: ConfigUpdateTooFrequent
  expr: rate(config_updates_total[5m]) > 10
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "配置更新过于频繁"
```

## 故障排查

### 1. etcd 连接失败

```bash
# 检查 etcd 服务状态
docker ps | grep etcd

# 查看 etcd 日志
docker logs etcd

# 测试连接
etcdctl endpoint health
```

### 2. 配置未生效

```bash
# 检查配置是否存在于 etcd
etcdctl get /coze-studio/production/tenant_123/config/feature.flag.smart_routing

# 检查应用日志
grep "config center" /var/log/coze-studio/app.log

# 强制刷新配置（重启服务）
systemctl restart coze-studio
```

### 3. 回滚失败

```bash
# 查看历史版本
GET /api/config/history?config_id=config_abc123

# 手动回滚（更新为旧值）
POST /api/config/update
{
  "tenant_id": "tenant_123",
  "config_key": "feature.flag.smart_routing",
  "new_value": "true",
  "change_reason": "手动回滚"
}
```

## 性能优化

### 1. etcd 性能调优

```yaml
# etcd 配置优化
--quota-backend-bytes: 8GB        # 增加配额
--snapshot-count: 100000          # 快照计数
--auto-compaction-mode: periodic  # 自动压缩
--auto-compaction-retention: 5m   # 保留时间
```

### 2. 客户端缓存

```go
// 使用本地缓存减少 etcd 访问
type CachedConfigService struct {
    cache    *lru.Cache
    etcdCli  *EtcdClient
}

func (s *CachedConfigService) Get(key string) (string, error) {
    // 先查缓存
    if val, ok := s.cache.Get(key); ok {
        return val.(string), nil
    }

    // 缓存未命中，查询 etcd
    val, err := s.etcdCli.Get(ctx, key)
    if err != nil {
        return "", err
    }

    // 更新缓存
    s.cache.Add(key, val)
    return val, nil
}
```

## 参考资料

- [etcd 官方文档](https://etcd.io/docs/)
- [配置管理最佳实践](https://martinfowler.com/articles/externalized-configuration.html)
- [Feature Flags 实践](https://www.launchdarkly.com/blog/what-are-feature-flags/)
