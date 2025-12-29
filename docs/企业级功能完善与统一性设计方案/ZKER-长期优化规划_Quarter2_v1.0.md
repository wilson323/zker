# ZKER 长期优化规划 (Quarter 2)

> **文档版本**: v1.0
> **制定日期**: 2025-01-01
> **执行周期**: Quarter 2 (12周)
> **负责人**: 技术团队

---

## 📋 目录

1. [执行摘要](#执行摘要)
2. [多地域部署](#多地域部署)
3. [数据库分库分表](#数据库分库分表)
4. [服务网格 (Istio)](#服务网格-istio)
5. [实施路线图](#实施路线图)

---

## 执行摘要

本方案涵盖Quarter 2的长期优化措施，实现全球化部署和百万级租户支持：

| 优化项目 | 预期收益 | 优先级 | 工作量 |
|---------|---------|--------|--------|
| 多地域部署 | 全国延迟<50ms，RTO<5min | P0 | 4周 |
| 数据库分库分表 | 支持10万租户 | P0 | 4周 |
| 服务网格 | 安全性100%，运维效率+50% | P1 | 4周 |

**总工作量**: 12周 (Quarter 2)

---

## 多地域部署

### 1. 当前架构

```
     全国用户
         │
         ▼
  ┌────────────┐
  │  单机房    │
  │  (北京)    │
  └────────────┘
```

**问题**:
- 南方用户延迟高 (>100ms)
- 单点故障风险
- 无法应对地域级灾难

### 2. 目标架构

```
     全国用户
         │
         ▼
  ┌─────────────┐
  │  智能DNS    │
  │ (GeoDNS)    │
  └──────┬──────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼───┐ ┌─────┐
│ 华东  │ │ 华北 │ │华南 │
│ 上海  │ │ 北京 │ │广州 │
└───┬───┘ └──┬───┘ └──┬──┘
    │         │        │
    └────┬────┴────┬───┘
         │         │
    ┌────▼─────────▼────┐
    │  数据同步中心     │
    │  (CDC/Canal)     │
    └──────────────────┘
```

### 3. 地域选择

| 地域 | 城市 | 机房 | 覆盖范围 |
|------|------|------|---------|
| 华东 | 上海 | 阿里云 | 上海、江苏、浙江 |
| 华北 | 北京 | 腾讯云 | 北京、河北、天津 |
| 华南 | 广州 | 华为云 | 广东、福建、广西 |

### 4. 智能DNS路由

**GeoDNS配置**:
```bash
# BIND配置示例
$ORIGIN zker.com.

@       IN  SOA ns1.zker.com. admin.zker.com. (
        2025010101  ; serial
        3600        ; refresh
        1800        ; retry
        604800      ; expire
        86400 )     ; minimum

; 华东用户路由到上海机房
api     IN  A   192.168.1.1   ; 上海
        IN  A   192.168.1.2

; 华北用户路由到北京机房
api     IN  A   192.168.2.1   ; 北京
        IN  A   192.168.2.2

; 华南用户路由到广州机房
api     IN  A   192.168.3.1   ; 广州
        IN  A   192.168.3.2

; 使用GeoIP实现智能路由
```

### 5. 数据同步

**Canal CDC同步**:
```
MySQL Master (北京)
      │
      ├─ Binlog
      ▼
  Canal Server
      │
      ├─► MySQL Slave (上海)
      └─► MySQL Slave (广州)
```

**配置示例**:
```yaml
# Canal配置
canal.destinations=zker_east,zker_south
canal.instance.zker_east.filter=zker\\.tenants,zker\\.bots
canal.instance.zker_east.canalSlaveId=1234
canal.instance.zker_south.canalSlaveId=5678
```

---

## 数据库分库分表

### 1. 当前架构

```
┌─────────────┐
│  单库单表   │
│  (MySQL)    │
│             │
│ 1亿+ 记录   │
└─────────────┘
```

**问题**:
- 单表数据量过大
- 查询性能下降
- 备份和恢复困难
- 无法水平扩展

### 2. 目标架构

```
                  ┌──────────────┐
                  │  应用层      │
                  └──────┬───────┘
                         │
                ┌────────▼─────────┐
                │  分片路由层      │
                │  (Vitess Proxy)  │
                └──────┬───────────┘
                       │
      ┌────────────────┼────────────────┐
      │                │                │
┌─────▼─────┐   ┌─────▼─────┐   ┌─────▼─────┐
│  Shard 1  │   │  Shard 2  │   │  Shard 3  │
│  tenant0  │   │  tenant1  │   │  tenant2  │
│  (MySQL)  │   │  (MySQL)  │   │  (MySQL)  │
└───────────┘   └───────────┘   └───────────┘
```

### 3. 分片策略

#### 3.1 按租户ID分库

**分库规则**:
```go
// 计算分片ID (0-1023)
func GetShardID(tenantID string) int {
    // 使用CRC32哈希
    hash := crc32.ChecksumIEEE([]byte(tenantID))
    return int(hash % 1024)  // 1024个分片
}

// 计算数据库名称
func GetDatabaseName(shardID int) string {
    // 每个数据库包含4个分片
    dbID := shardID / 4
    return fmt.Sprintf("zker_%04d", dbID)  // zker_0000 - zker_0255
}

// 计算表名
func GetTableName(shardID int) string {
    // 每个分片对应一个表
    tableID := shardID % 4
    return fmt.Sprintf("tenants_%d", tableID)  // tenants_0 - tenants_3
}
```

**分片示例**:
```
租户ID: tenant_123456
├─ CRC32哈希: 0x8f9a2b3c
├─ 分片ID: 932
├─ 数据库: zker_0233  (932/4=233)
└─ 表名: tenants_0    (932%4=0)
```

#### 3.2 分库分表映射

```
分片ID范围    数据库        表
────────────────────────────────────
0-3          zker_0000     tenants_0 - tenants_3
4-7          zker_0001     tenants_0 - tenants_3
8-11         zker_0002     tenants_0 - tenants_3
...
1020-1023    zker_0255     tenants_0 - tenants_3
```

### 4. Vitess部署

```yaml
# docker/vitess.yml
version: '3.8'

services:
  # Vitess集群拓扑
  etcd:
    image: vitess/etcd:vital
    ports:
      - "2379:2379"
      - "2380:2380"

  vtctld:
    image: vitess/vtctld:vital
    depends_on:
      - etcd

  vtgate:
    image: vitess/vgate:vital
    ports:
      - "15306:15306"  # MySQL协议
      - "15999:15999"  # gRPC
    depends_on:
      - etcd

  vttablet:
    image: vitess/vtablet:vital
    ports:
      - "15100:15100"  # MySQL协议
      - "16100:16100"  # gRPC
```

### 5. 数据迁移

**双写迁移方案**:
```
阶段1: 双写 (新数据写入旧库和新库)
┌─────────┐     ┌─────────┐
│ 旧库    │◄───►│ 新库    │
└─────────┘     └─────────┘
       │
       └──► 读操作仍从旧库

阶段2: 切换读 (读操作切换到新库)
┌─────────┐     ┌─────────┐
│ 旧库    │     │ 新库    │◄───读
└─────────┘     └─────────┘
       │
       └──► 停止写入

阶段3: 下线旧库 (只保留新库)
┌─────────┐     ┌─────────┐
│ 旧库    │     │ 新库    │◄───读写
│(已下线) │     └─────────┘
└─────────┘
```

---

## 服务网格 (Istio)

### 1. 架构设计

```
┌─────────────────────────────────────┐
│         Istio Service Mesh         │
│  ┌───────────────────────────────┐  │
│  │  Control Plane                 │  │
│  │  - Istiod (Pilot)             │  │
│  │  - Citadel (CA)               │  │
│  │  - Galley (Config)            │  │
│  └───────────────────────────────┘  │
└───────────────┬─────────────────────┘
                │
                ▼
┌─────────────────────────────────────┐
│         Data Plane (Envoy)         │
│                                     │
│  ┌──────────┐  ┌──────────┐       │
│  │ Tenant   │  │ Quota    │       │
│  │ Service  │  │ Service  │       │
│  │ + Envoy  │  │ + Envoy  │       │
│  └──────────┘  └──────────┘       │
└─────────────────────────────────────┘
```

### 2. 核心功能

#### 2.1 流量管理

```yaml
# VirtualService配置
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: tenant-service
spec:
  hosts:
  - tenant-service
  http:
  - match:
    - headers:
        x-user-region:
          exact: east
    route:
    - destination:
        host: tenant-service
        subset: v2  # 华东版本
      weight: 100
  - route:
    - destination:
        host: tenant-service
        subset: v1  # 默认版本
      weight: 100
```

#### 2.2 安全通信

```yaml
# PeerAuthentication配置
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
spec:
  mtls:
    mode: STRICT  # 强制mTLS加密
```

#### 2.3 可观测性

```yaml
# ServiceEntry配置
apiVersion: networking.istio.io/v1beta1
kind: ServiceEntry
metadata:
  name: external-mysql
spec:
  hosts:
  - mysql.external.com
  ports:
  - number: 3306
    name: mysql
    protocol: TCP
  location: MESH_EXTERNAL
```

### 3. 部署架构

```bash
# 安装Istio
istioctl install --set profile=demo

# 注入Sidecar
kubectl label namespace default istio-injection=enabled

# 部署应用
kubectl apply -f tenant-service.yaml
```

### 4. 监控集成

**Prometheus + Grafana + Kiali**:
```yaml
# Kiali配置
apiVersion: kiali.io/v1alpha1
kind: Kiali
metadata:
  name: kiali
spec:
  auth:
    strategy: anonymous
  external_services:
    prometheus:
      url: "http://prometheus:9090"
    grafana:
      url: "http://grafana:3000"
```

---

## 实施路线图

### Month 4 (Week 1-4)

| 周次 | 任务 | 负责人 | 里程碑 |
|------|------|--------|--------|
| Week 1 | 地域机房选择 + 基础设施准备 | 研发D | 完成机房评估 |
| Week 2-3 | 多地域部署实施 | 研发D + 研发B | 完成3地域部署 |
| Week 4 | GeoDNS配置 + 故障切换测试 | 研发D | 实现智能路由 |

### Month 5 (Week 5-8)

| 周次 | 任务 | 负责人 | 里程碑 |
|------|------|--------|--------|
| Week 5 | Vitess部署 + 分片方案设计 | 研发B + DBA | 完成分片设计 |
| Week 6-7 | 数据迁移实施 | 研发B + DBA | 迁移10万租户 |
| Week 8 | 性能测试 + 验收 | 研发B | 达成性能目标 |

### Month 6 (Week 9-12)

| 周次 | 任务 | 负责人 | 里程碑 |
|------|------|--------|--------|
| Week 9 | Istio安装 + 基础配置 | 研发D | 服务网格就绪 |
| Week 10-11 | Sidecar注入 + 流量管理 | 研发B + 研发D | 所有服务接入网格 |
| Week 12 | mTLS启用 + 监控集成 | 研发B | 安全性100% |

---

## 验收标准

### 多地域部署

- [ ] 全国访问延迟 < 50ms
- [ ] 地域故障RTO < 5分钟
- [ ] 数据一致性 ≥ 99.99%
- [ ] 故障自动切换成功率 = 100%

### 分库分表

- [ ] 支持10万租户
- [ ] 查询性能保持稳定 (P95 < 100ms)
- [ ] 迁移过程零停机
- [ ] 分片路由正确率 = 100%

### 服务网格

- [ ] 所有服务间通信加密 (mTLS)
- [ ] 流量管理覆盖率 = 100%
- [ ] 服务拓扑可视化完整
- [ ] 故障定位时间 < 1分钟

---

**制定人**: 研发B - 后端工程师
**审核人**: 技术负责人
**执行周期**: Quarter 2
