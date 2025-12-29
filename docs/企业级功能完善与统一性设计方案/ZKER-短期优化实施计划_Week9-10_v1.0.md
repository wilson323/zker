# ZKER 短期优化实施计划 (Week 9-10)

> **文档版本**: v1.0
> **制定日期**: 2025-01-01
> **执行周期**: Week 9-10 (2周)
> **负责人**: 研发B - 后端工程师

---

## 📋 目录

1. [执行摘要](#执行摘要)
2. [优化目标](#优化目标)
3. [优化项目清单](#优化项目清单)
4. [实施计划](#实施计划)
5. [风险评估](#风险评估)
6. [验收标准](#验收标准)

---

## 执行摘要

本计划涵盖Week 9-10的短期优化措施，旨在进一步提升系统性能和可扩展性：

| 优化项目 | 预期收益 | 优先级 | 工作量 |
|---------|---------|--------|--------|
| 数据库读写分离 | 读QPS +200%, 查询延迟 -33% | P0 | 5天 |
| Redis集群化 | 缓存容量 +200%, 可用性 99.99% | P0 | 3天 |
| API网关优化 | 后端压力 -30%, 恶意请求拦截100% | P1 | 2天 |

**总工作量**: 10个工作日

---

## 优化目标

### 性能目标

| 指标 | 当前值 | 目标值 | 提升幅度 |
|------|--------|--------|---------|
| 数据库读QPS | 2,000 | 6,000 | +200% |
| 数据库查询P95 | 120ms | 80ms | -33% |
| 缓存容量 | 4GB | 12GB | +200% |
| 缓存可用性 | 99.9% | 99.99% | +0.09% |
| API响应时间 | 280ms | 200ms | -29% |
| 后端服务压力 | 100% | 70% | -30% |

### 可靠性目标

- ✅ 数据库主从切换时间 < 30秒
- ✅ Redis集群故障自动恢复 < 10秒
- ✅ API网关故障隔离 < 5秒
- ✅ 整体系统可用性 ≥ 99.95%

---

## 优化项目清单

### 1. 数据库读写分离

#### 1.1 架构设计

```
┌─────────────┐
│  API Gateway│
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  ProxySQL       │
│  (读写路由)      │
└──────┬──────────┘
       │
       ├─────────┬─────────┐
       ▼         ▼         ▼
   ┌───────┐ ┌───────┐ ┌───────┐
   │Master │ │Slave1 │ │Slave2 │
   │(Write)│ │(Read) │ │(Read) │
   └───────┘ └───────┘ └───────┘
```

**组件说明**:
- **ProxySQL**: 智能代理，自动路由读写请求
- **Master**: 主库，处理所有写操作
- **Slave1/Slave2**: 从库，处理所有读操作
- **复制**: MySQL Binlog异步复制

#### 1.2 技术方案

**ProxySQL配置**:
```bash
# 安装ProxySQL
docker run -d \
  --name proxysql \
  -p 6033:6033 \
  -p 6032:6032 \
  proxysql/proxysql:latest

# 配置后端服务器
INSERT INTO mysql_servers (
    hostgroup_id, hostname, port
) VALUES (
    10, 'mysql-master', 3306
), (
    20, 'mysql-slave1', 3306
), (
    20, 'mysql-slave2', 3306
);

# 配置读写路由规则
INSERT INTO mysql_query_rules (
    rule_id, active, match_pattern, destination_hostgroup, apply
) VALUES (
    1, 1, '^SELECT.*FOR UPDATE$', 10, 1
), (
    2, 1, '^SELECT', 20, 1
), (
    3, 1, '.*', 10, 1
);
```

**MySQL主从复制配置**:
```sql
-- Master配置 (my.cnf)
[mysqld]
server-id = 1
log-bin = mysql-bin
binlog-format = ROW
binlog-do-db = zker

-- Slave配置 (my.cnf)
[mysqld]
server-id = 2  # Slave1使用2, Slave2使用3
relay-log = mysql-relay-bin
read-only = 1
```

#### 1.3 代码修改

**GORM配置**:
```go
// config/database.go
type DatabaseConfig struct {
    WriteHost string `toml:"write_host"`
    ReadHost  string `toml:"read_host"`
    Port      int    `toml:"port"`
    User      string `toml:"user"`
    Password  string `toml:"password"`
    Database  string `toml:"database"`
}

// 初始化DB连接
func InitDB(cfg *DatabaseConfig) (*gorm.DB, error) {
    // 写库连接
    writeDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.WriteHost, cfg.Port, cfg.Database)

    // 读库连接
    readDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.ReadHost, cfg.Port, cfg.Database)

    // 使用DBResolver插件
    db, err := gorm.Open(mysql.Open(writeDSN), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // 配置读写分离
    db.Use(dbresolver.Register(dbresolver.Config{
        Replicas: []gorm.Dialector{
            mysql.Open(readDSN),
        },
        Policy: dbresolver.RandomPolicy{},
    }))

    return db, nil
}
```

**读写分离使用示例**:
```go
// 自动路由到写库
db.Create(&tenant)

// 自动路由到读库
db.Where("tenant_id = ?", id).First(&tenant)

// 强制使用写库
db.Clauses(dbresolver.Write).Where("tenant_id = ?", id).First(&tenant)
```

#### 1.4 实施步骤

**Day 1-2: 环境准备**
1. 部署MySQL从库 (2个)
2. 配置主从复制
3. 验证数据同步

**Day 3-4: ProxySQL部署**
1. 安装ProxySQL
2. 配置读写路由规则
3. 测试路由逻辑

**Day 5: 代码集成和测试**
1. 修改数据库配置
2. 集成GORM DBResolver
3. 压力测试验证

#### 1.5 监控指标

```go
// 新增Prometheus指标
var (
    DBWriteQPS = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_write_qps_total",
            Help: "Database write queries per second",
        },
        []string{"database"},
    )

    DBReadQPS = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_read_qps_total",
            Help: "Database read queries per second",
        },
        []string{"database"},
    )

    DBReplicationLag = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "db_replication_lag_seconds",
            Help: "Database replication lag in seconds",
        },
        []string{"master", "slave"},
    )
)
```

#### 1.6 告警规则

```yaml
# prometheus/alerts.yml
- alert: DBReplicationLagHigh
  expr: db_replication_lag_seconds > 10
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Database replication lag is too high"
    description: "Replication lag {{ $value }}s for {{ $labels.master }} -> {{ $labels.slave }}"

- alert: DBSlaveDown
  expr: up{job="mysql-slave"} == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "MySQL slave is down"
    description: "Slave {{ $labels.instance }} is not responding"
```

---

### 2. Redis集群化

#### 2.1 架构设计

```
                    ┌─────────────┐
                    │  Application│
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Redis Cluster│
                    │  (3主3从)    │
                    └──────┬──────┘
                           │
      ┌────────────────────┼────────────────────┐
      │                    │                    │
┌─────▼─────┐       ┌─────▼─────┐       ┌─────▼─────┐
│  Shard 1  │       │  Shard 2  │       │  Shard 3  │
│ Master A  │       │ Master B  │       │ Master C  │
│ Slave A'  │       │ Slave B'  │       │ Slave C'  │
└───────────┘       └───────────┘       └───────────┘
   Slot 0-5460       Slot 5461-10922     Slot 10923-16383
```

**组件说明**:
- **3个主节点**: 分片存储数据
- **3个从节点**: 主节点副本，自动故障转移
- **16384个槽位**: 数据分片
- **Gossip协议**: 节点间通信
- **Smart Client**: 客户端路由

#### 2.2 技术方案

**Redis Cluster配置**:
```bash
# 创建Redis配置文件
# redis-cluster-1.conf (Shard 1 Master)
port 7001
cluster-enabled yes
cluster-config-file nodes-7001.conf
cluster-node-timeout 5000
appendonly yes
appendfilename "appendonly-7001.aof"
dbfilename "dump-7001.rdb"

# redis-cluster-4.conf (Shard 1 Slave)
port 7004
cluster-enabled yes
cluster-config-file nodes-7004.conf
cluster-node-timeout 5000
slaveof redis-cluster-1 7001
```

**Docker部署**:
```yaml
# docker/redis-cluster.yml
version: '3.8'

services:
  redis-cluster-1:
    image: redis:7-alpine
    command: redis-server /usr/local/etc/redis/redis-cluster-1.conf
    ports:
      - "7001:7001"
    volumes:
      - ./redis/redis-cluster-1.conf:/usr/local/etc/redis/redis-cluster-1.conf
      - redis-cluster-1-data:/data
    networks:
      - redis-cluster

  redis-cluster-4:
    image: redis:7-alpine
    command: redis-server /usr/local/etc/redis/redis-cluster-4.conf
    ports:
      - "7004:7004"
    volumes:
      - ./redis/redis-cluster-4.conf:/usr/local/etc/redis/redis-cluster-4.conf
      - redis-cluster-4-data:/data
    networks:
      - redis-cluster

  # ... 其他节点配置

networks:
  redis-cluster:
    driver: bridge

volumes:
  redis-cluster-1-data:
  redis-cluster-4-data:
```

**集群初始化**:
```bash
# 创建集群 (3主3从)
redis-cli --cluster create \
  127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 \
  127.0.0.1:7004 127.0.0.1:7005 127.0.0.1:7006 \
  --cluster-replicas 1

# 验证集群状态
redis-cli -c -p 7001 cluster info
redis-cli -c -p 7001 cluster nodes
```

#### 2.3 代码修改

**Go Redis客户端配置**:
```go
// infra/cache/redis_cluster.go
package cache

import (
    "context"
    "github.com/redis/go-redis/v9"
)

type RedisClusterConfig struct {
    Addrs []string `toml:"addrs"` // ["localhost:7001", "localhost:7002", ...]
    MaxRetries int    `toml:"max_retries"`
    PoolSize   int    `toml:"pool_size"`
}

func NewRedisCluster(cfg *RedisClusterConfig) (*redis.ClusterClient, error) {
    client := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs:     cfg.Addrs,
        MaxRetries: cfg.MaxRetries,
        PoolSize:   cfg.PoolSize,
        ReadOnly:   true,  // 允许从节点读取
        RouteRandomly: true, // 随机路由到从节点
    })

    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := client.Ping(ctx).Err()
    if err != nil {
        return nil, err
    }

    return client, nil
}
```

**使用示例**:
```go
// 自动路由到正确的分片
err := redisCluster.Set(ctx, "tenant:123", tenantData, 1*time.Hour).Err()

// 自动路由到从节点读取
val, err := redisCluster.Get(ctx, "tenant:123").Result()
```

#### 2.4 实施步骤

**Day 1: Redis集群部署**
1. 准备6个Redis配置文件
2. 部署3主3从集群
3. 验证集群健康状态

**Day 2: 代码集成**
1. 更新Redis客户端配置
2. 集成go-redis集群客户端
3. 功能测试

**Day 3: 数据迁移和监控**
1. 数据从单机迁移到集群
2. 配置集群监控
3. 性能测试验证

#### 2.5 监控指标

```go
// 新增Redis集群监控指标
var (
    RedisClusterMemoryBytes = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "redis_cluster_memory_bytes",
            Help: "Redis cluster memory usage in bytes",
        },
        []string{"node", "shard"},
    )

    RedisClusterKeyCount = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "redis_cluster_keys_total",
            Help: "Total number of keys in Redis cluster",
        },
        []string{"node", "shard"},
    )

    RedisClusterOpsPerSecond = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "redis_cluster_ops_total",
            Help: "Redis cluster operations per second",
        },
        []string{"node", "shard", "operation"},
    )
)
```

#### 2.6 告警规则

```yaml
# prometheus/alerts.yml
- alert: RedisClusterNodeDown
  expr: redis_up{job="redis-cluster"} == 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Redis cluster node is down"
    description: "Node {{ $labels.instance }} in shard {{ $labels.shard }} is not responding"

- alert: RedisClusterMemoryHigh
  expr: redis_cluster_memory_bytes / redis_cluster_max_memory > 0.9
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Redis cluster memory usage is high"
    description: "Node {{ $labels.instance }} memory usage is {{ $value }}%"
```

---

### 3. API网关优化

#### 3.1 架构设计

```
┌──────────┐
│  Client  │
└────┬─────┘
     │
┌────▼───────────────────────┐
│      Kong API Gateway       │
│  ┌───────────────────────┐  │
│  │ Rate Limiting         │  │
│  │ Circuit Breaker       │  │
│  │ Request Transformation │  │
│  │ Authentication        │  │
│  └───────────────────────┘  │
└────┬───────────────────────┘
     │
     ├─────────┬─────────┬─────────┐
     ▼         ▼         ▼         ▼
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│Tenant  │ │Quota   │ │Subs    │ │Auth    │
│Service │ │Service │ │Service │ │Service │
└────────┘ └────────┘ └────────┘ └────────┘
```

**功能特性**:
- **限流**: 令牌桶算法，防止过载
- **熔断**: 自动隔离故障服务
- **降级**: 返回缓存数据或默认响应
- **认证**: 统一JWT验证
- **监控**: 自动收集请求指标

#### 3.2 技术方案

**Kong网关部署**:
```yaml
# docker/kong.yml
version: '3.8'

services:
  kong:
    image: kong:3.4-alpine
    environment:
      KONG_DATABASE: "off"
      KONG_PROXY_ACCESS_LOG: /dev/stdout
      KONG_ADMIN_ACCESS_LOG: /dev/stdout
      KONG_PROXY_ERROR_LOG: /dev/stderr
      KONG_ADMIN_ERROR_LOG: /dev/stderr
      KONG_ADMIN_LISTEN: "0.0.0.0:8001"
    ports:
      - "8000:8000"   # Proxy
      - "8443:8443"   # Proxy SSL
      - "8001:8001"   # Admin
    networks:
      - kong-net

  konga:
    image: pantsel/konga:latest
    environment:
      NODE_ENV: production
    ports:
      - "1337:1337"
    networks:
      - kong-net

networks:
  kong-net:
    driver: bridge
```

**Kong配置**:
```bash
# 1. 添加上游服务
curl -X POST http://localhost:8001/services \
  --data name=tenant-service \
  --data url=http://tenant-service:8080

# 2. 添加路由规则
curl -X POST http://localhost:8001/services/tenant-service/routes \
  --data paths[]=/api/v1/tenants \
  --data strip_path=false

# 3. 添加限流插件
curl -X POST http://localhost:8001/services/tenant-service/plugins \
  --data name=rate-limiting \
  --data config.minute=100 \
  --data config.hour=1000 \
  --data config.policy=local

# 4. 添加熔断插件
curl -X POST http://localhost:8001/services/tenant-service/plugins \
  --data name=circuit-breaker \
  --data config.error_threshold=50 \
  --data config.volume_threshold=10 \
  --data config.half_open_timeout=30000

# 5. 添加JWT认证插件
curl -X POST http://localhost:8001/services/tenant-service/plugins \
  --data name=jwt

# 6. 添加Prometheus监控插件
curl -X POST http://localhost:8001/plugins \
  --data name=prometheus \
  --data config.per_consumer=true
```

#### 3.3 实施步骤

**Day 1: Kong网关部署**
1. 部署Kong和Konga (管理UI)
2. 配置服务发现
3. 验证代理功能

**Day 2: 插件配置**
1. 配置限流插件
2. 配置熔断插件
3. 配置JWT认证
4. 配置监控插件

#### 3.4 限流策略

**按租户限流**:
```yaml
# 不同套餐的限流策略
free_tier:
  requests_per_minute: 60
  requests_per_hour: 1000

basic_tier:
  requests_per_minute: 100
  requests_per_hour: 2000

premium_tier:
  requests_per_minute: 200
  requests_per_hour: 5000

enterprise_tier:
  requests_per_minute: 500
  requests_per_hour: 10000
```

**按API端点限流**:
```yaml
# 配额检查API (高频调用)
/api/v1/quota/check:
  requests_per_minute: 200

# 租户管理API (低频调用)
/api/v1/tenants:
  requests_per_minute: 100
```

#### 3.5 熔断配置

**熔断触发条件**:
- 错误率 > 50% (10个请求内)
- 响应时间 > 5s
- 连续失败次数 > 5

**熔断状态**:
1. **Closed**: 正常状态
2. **Open**: 熔断开启，拒绝请求
3. **Half-Open**: 半开状态，尝试恢复

**降级策略**:
- 返回缓存数据
- 返回默认响应
- 返回友好错误信息

#### 3.6 监控指标

```yaml
# Kong Prometheus插件自动暴露指标
# kong_http_requests_total
# kong_latency_ms
# kong_upstream_latency_ms
```

---

## 实施计划

### Week 9 日程

| 日期 | 任务 | 负责人 | 工作量 |
|------|------|--------|--------|
| Day 1-2 | 数据库读写分离 - 环境准备 | 研发B + 研发D | 2天 |
| Day 3-4 | 数据库读写分离 - ProxySQL部署 | 研发B + 研发D | 2天 |
| Day 5 | 数据库读写分离 - 代码集成和测试 | 研发B | 1天 |

### Week 10 日程

| 日期 | 任务 | 负责人 | 工作量 |
|------|------|--------|--------|
| Day 1 | Redis集群化 - 部署和配置 | 研发B + 研发D | 1天 |
| Day 2 | Redis集群化 - 代码集成 | 研发B | 1天 |
| Day 3 | Redis集群化 - 数据迁移和测试 | 研发B | 1天 |
| Day 4 | API网关优化 - 部署和插件配置 | 研发B + 研发D | 1天 |
| Day 5 | 集成测试和性能验证 | 研发B | 1天 |

---

## 风险评估

### 高风险项

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|---------|
| 主从数据不一致 | 数据错误 | 中 | 实施半同步复制，监控复制延迟 |
| Redis集群分片迁移失败 | 服务中断 | 低 | 充分测试，准备回滚方案 |
| Kong网关单点故障 | 服务不可用 | 中 | 部署多实例 + 负载均衡 |

### 回滚方案

**数据库读写分离回滚**:
```bash
# 1. 停止ProxySQL
docker stop proxysql

# 2. 修改应用配置，直连Master
# config.toml
read_host = "mysql-master"  # 改为写库地址

# 3. 重启应用
make restart
```

**Redis集群回滚**:
```bash
# 1. 停止集群应用
# 2. 修改应用配置，连接单机Redis
# config.toml
redis_mode = "standalone"
redis_addr = "localhost:6379"

# 3. 重启应用
make restart
```

**Kong网关回滚**:
```bash
# 1. 停止Kong
docker stop kong

# 2. 修改Nginx配置，直连后端服务
# 3. 重启Nginx
nginx -s reload
```

---

## 验收标准

### 功能验收

- [ ] 数据库读写分离正常工作
- [ ] Redis集群自动分片和故障转移
- [ ] API网关限流、熔断、认证功能正常

### 性能验收

- [ ] 数据库读QPS ≥ 6,000 (当前2,000)
- [ ] 数据库查询P95 ≤ 80ms (当前120ms)
- [ ] Redis容量 ≥ 12GB (当前4GB)
- [ ] API响应时间 ≤ 200ms (当前280ms)

### 可靠性验收

- [ ] 主从切换时间 < 30秒
- [ ] Redis集群故障恢复 < 10秒
- [ ] 系统可用性 ≥ 99.95%

### 监控验收

- [ ] 所有新增指标正常采集
- [ ] 告警规则正常触发
- [ ] Grafana大盘正常显示

---

## 附录

### 参考文档

- [ProxySQL官方文档](https://proxysql.com/documentation/)
- [Redis Cluster官方文档](https://redis.io/docs/manual/scaling/)
- [Kong网关官方文档](https://docs.konghq.com/)
- [GORM读写分离文档](https://gorm.io/docs/dbresolver.html)
- [go-redis集群客户端](https://redis.uptrace.dev/cluster)

### 相关配置文件

- `docker/proxysql.yml` - ProxySQL Docker配置
- `docker/redis-cluster.yml` - Redis集群配置
- `docker/kong.yml` - Kong网关配置
- `backend/config/database.toml` - 数据库配置
- `backend/config/cache.toml` - 缓存配置

---

**制定人**: 研发B - 后端工程师
**审核人**: 技术负责人
**执行周期**: Week 9-10
**下次更新**: Week 10结束时
