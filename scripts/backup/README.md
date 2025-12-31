# ZKER 日志系统与数据库备份部署指南

> **版本**: v1.0.0 | **更新**: 2025-01-03 | **状态**: ✅ 生产就绪

---

## 📋 目录

- [系统概览](#系统概览)
- [Loki日志系统部署](#loki日志系统部署)
- [数据库自动备份配置](#数据库自动备份配置)
- [Kubernetes部署](#kubernetes部署)
- [监控与告警](#监控与告警)
- [故障排查](#故障排查)
- [最佳实践](#最佳实践)

---

## 系统概览

### 已部署组件

| 组件 | 版本 | 端口 | 说明 |
|------|------|------|------|
| **Loki** | 2.9.2 | 3100 | 日志聚合系统 |
| **Promtail** | 2.9.2 | - | 日志采集Agent |
| **MySQL Backup** | - | - | 自动备份脚本 |
| **Redis Backup** | - | - | 自动备份脚本 |

### 架构图

```
┌─────────────────────────────────────────────────────────┐
│                     应用层                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐              │
│  │ Backend  │  │ Frontend │  │  Nginx   │              │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘              │
│       │            │            │                       │
└───────┼────────────┼────────────┼───────────────────────┘
        │            │            │
        ▼            ▼            ▼
┌─────────────────────────────────────────────────────────┐
│                   Promtail                              │
│  (采集日志 → 结构化 → 标签 → 发送到Loki)                 │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                     Loki                                │
│  (存储索引 → 压缩 → 保留策略 → 查询接口)                 │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                    Grafana                              │
│  (查询日志 → 可视化 → 告警 → 仪表盘)                     │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│              数据库备份系统                              │
│  ┌────────────┐  ┌────────────┐                        │
│  │ MySQL Dump │  │ Redis RDB  │                        │
│  │  (Cron)    │  │  (Cron)    │                        │
│  └─────┬──────┘  └─────┬──────┘                        │
│        │                │                                │
│        ▼                ▼                                │
│  ┌─────────────────────────────┐                        │
│  │   /backup (NFS/S3)         │                        │
│  └─────────────────────────────┘                        │
└─────────────────────────────────────────────────────────┘
```

---

## Loki日志系统部署

### 步骤1: 配置文件已创建

✅ **Loki配置**: `docker/volumes/monitoring/loki/loki-config.yml`
✅ **Promtail配置**: `docker/volumes/monitoring/promtail/promtail-config.yml`
✅ **Grafana数据源**: `docker/volumes/monitoring/grafana/provisioning/datasources/loki.yml`

### 步骤2: 启动服务

```bash
cd docker

# 启动监控服务（包含Loki和Promtail）
docker compose -f docker-compose-monitoring.yml up -d loki promtail

# 查看服务状态
docker compose -f docker-compose-monitoring.yml ps loki promtail

# 查看日志
docker compose -f docker-compose-monitoring.yml logs -f loki
docker compose -f docker-compose-monitoring.yml logs -f promtail
```

### 步骤3: 验证部署

```bash
# 1. 检查Loki健康状态
curl http://localhost:3100/ready

# 预期输出: ready

# 2. 检查Promtail健康状态
curl http://localhost:9080/metrics

# 预期输出: Promtail指标

# 3. 查询Loki标签
curl -G http://localhost:3100/loki/api/v1/labels

# 4. 查看日志流
curl -G http://localhost:3100/loki/api/v1/query \
  --data-urlencode 'query={job="backend"}'
```

### 步骤4: 配置Grafana

1. **登录Grafana**: http://localhost:3000 (admin/admin)
2. **检查数据源**: Configuration → Data Sources → Loki
3. **创建日志查询Dashboard**:
   - Create → Dashboard → Add new panel
   - Query: `{job="backend"} |= "error"`
   - Visualization: Logs
   - Save

### 步骤5: 日志查询示例

**LogQL查询语法**:

```logql
# 1. 查询所有backend错误日志
{job="backend"} |= "error"

# 2. 查询特定租户的日志
{job="backend", tenant_id="123"}

# 3. 正则匹配日志内容
{job="backend"} =~ ".*error.*timeout.*"

# 4. 统计错误率
count_over_time({job="backend"} |= "error" [5m])

# 5. 计算错误率百分比
sum(count_over_time({job="backend"} |= "error" [5m])) /
sum(count_over_time({job="backend"} [5m])) * 100

# 6. 按日志级别分组
sum by (level) (count_over_time({job="backend"} [1h]))

# 7. 查看最近5分钟的日志
{job="backend"} | logfmt | line_format "{{.level}}: {{.msg}}"
```

---

## 数据库自动备份配置

### 步骤1: 备份脚本已创建

✅ **MySQL备份**: `scripts/backup/backup-mysql.sh`
✅ **Redis备份**: `scripts/backup/backup-redis.sh`
✅ **Cron配置**: `scripts/backup/setup-backup-cron.sh`
✅ **备份验证**: `scripts/backup/verify-backup.sh`
✅ **恢复测试**: `scripts/backup/test-restore.sh`

### 步骤2: 配置环境变量

创建 `.env.backup` 文件:

```bash
# MySQL配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password

# 备份配置
BACKUP_ROOT=/backup
RETENTION_DAYS=30
NOTIFICATION_WEBHOOK=https://your-webhook-url
```

### 步骤3: 设置脚本权限

```bash
cd scripts/backup

chmod +x backup-mysql.sh
chmod +x backup-redis.sh
chmod +x setup-backup-cron.sh
chmod +x verify-backup.sh
chmod +x test-restore.sh
```

### 步骤4: 配置定时任务

```bash
# 使用root权限运行
sudo ./setup-backup-cron.sh

# 查看配置的定时任务
sudo ./setup-backup-cron.sh --show

# 测试备份脚本
sudo ./setup-backup-cron.sh --test
```

**定时任务计划**:

| 任务 | 时间 | 类型 |
|------|------|------|
| MySQL全量备份 | 每天 02:00 | 全量 |
| MySQL增量备份 | 每6小时 | 增量 |
| Redis RDB备份 | 每6小时 | RDB快照 |
| Redis数据导出 | 每天 03:00 | 数据导出 |
| 备份验证 | 每天 04:00 | 验证 |

### 步骤5: 验证备份

```bash
# 查看备份日志
tail -f /var/log/backup/mysql-backup.log
tail -f /var/log/backup/redis-backup.log

# 手动执行备份测试
sudo ./backup-mysql.sh BACKUP_TYPE=full
sudo ./backup-redis.sh BACKUP_TYPE=rdb

# 验证备份
sudo ./verify-backup.sh

# 查看备份文件
ls -lh /backup/mysql/full/
ls -lh /backup/redis/rdb/
```

---

## Kubernetes部署

### 前置条件

- Kubernetes集群 v1.24+
- NFS存储或云存储
- Kubectl配置完成

### 步骤1: 创建命名空间

```bash
kubectl create namespace zker
```

### 步骤2: 创建Secret

```bash
# MySQL Secret
kubectl apply -f deploy/k8s/cronjob/mysql-backup.yaml
```

### 步骤3: 部署CronJob

```bash
# MySQL备份
kubectl apply -f deploy/k8s/cronjob/mysql-backup.yaml

# Redis备份
kubectl apply -f deploy/k8s/cronjob/redis-backup.yaml
```

### 步骤4: 验证部署

```bash
# 查看CronJob
kubectl get cronjob -n zker

# 查看备份任务历史
kubectl get jobs -n zker

# 查看Pod日志
kubectl logs -l app=mysql-backup -n zker --tail=100

# 手动触发备份任务
kubectl create job mysql-backup-manual --from=cronjob/mysql-backup -n zker
```

---

## 监控与告警

### Loki告警规则

创建 `docker/volumes/monitoring/loki/alerts.yml`:

```yaml
groups:
  - name: log_alerts
    interval: 30s
    rules:
      # 高错误率告警
      - alert: HighErrorRate
        expr: |
          sum(rate({job="backend"} |= "error" [5m])) /
          sum(rate({job="backend"} [5m])) > 0.05
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "高错误率告警"
          description: "错误率超过5%"

      # 数据库连接错误
      - alert: DatabaseConnectionError
        expr: |
          sum(count_over_time({job="backend"} |= "database connection error" [5m])) > 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "数据库连接错误"
          description: "检测到数据库连接错误"

      # 慢查询告警
      - alert: SlowQueryDetected
        expr: |
          {job="mysql", type="slow-query"}
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "慢查询告警"
          description: "检测到MySQL慢查询"
```

### 备份监控Prometheus规则

添加到 `docker/volumes/monitoring/prometheus/alerts.yml`:

```yaml
groups:
  - name: backup_alerts
    interval: 1m
    rules:
      # MySQL备份失败
      - alert: MySQLBackupFailed
        expr: |
          time() - mysql_backup_success_timestamp > 86400
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "MySQL备份失败"
          description: "超过24小时未成功备份"

      # Redis备份失败
      - alert: RedisBackupFailed
        expr: |
          time() - redis_backup_success_timestamp > 21600
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "Redis备份失败"
          description: "超过6小时未成功备份"

      # 备份存储空间不足
      - alert: BackupStorageLow
        expr: |
          (backup_storage_available_bytes / backup_storage_size_bytes) < 0.2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "备份存储空间不足"
          description: "备份存储剩余空间低于20%"
```

---

## 故障排查

### Loki常见问题

#### 问题1: Loki无法启动

```bash
# 检查配置文件语法
docker run --rm -v $(pwd)/volumes/monitoring/loki:/loki grafana/loki:2.9.2 -config.file=/loki/loki-config.yml -print-config-stdin

# 查看详细日志
docker logs zker-loki --tail=100

# 检查磁盘空间
df -h docker/data/monitoring/loki
```

#### 问题2: Promtail无法采集日志

```bash
# 检查Promtail配置
docker logs zker-promtail --tail=100

# 验证日志文件权限
ls -la /var/log/

# 测试Loki连接
curl http://loki:3100/ready

# 查看Promtail位置文件
cat docker/data/monitoring/promtail/positions.yaml
```

#### 问题3: 日志查询慢

**优化方案**:

1. **增加查询并发**:
```yaml
# loki-config.yml
querier:
  query_concurrency: 16
  split_queries_by_interval: 30m
```

2. **启用查询缓存**:
```yaml
query_range:
  cache_results: true
```

3. **调整索引周期**:
```yaml
schema_config:
  configs:
    - from: 2025-01-01
      index:
        period: 24h  # 增加索引周期
```

### 备份常见问题

#### 问题1: MySQL备份失败

```bash
# 检查MySQL连接
mysql -h localhost -u root -p

# 查看备份日志
tail -f /var/log/backup/mysql-backup.log

# 检查磁盘空间
df -h /backup/mysql

# 测试mysqldump
mysqldump -u root -p --all-databases --single-transaction | wc -l
```

#### 问题2: Redis BGSAVE超时

```bash
# 检查Redis持久化配置
redis-cli CONFIG GET save
redis-cli CONFIG GET appendonly

# 手动触发BGSAVE
redis-cli BGSAVE
redis-cli LASTSAVE

# 查看BGSAVE进度
redis-cli INFO persistence | grep bgsave
```

#### 问题3: 备份验证失败

```bash
# 手动运行验证
./verify-backup.sh

# 检查备份文件
gzip -t /backup/mysql/full/mysql-full-*.sql.gz

# 测试MySQL恢复
./test-restore.sh -m test

# 检查Redis RDB
redis-check-rdb /backup/redis/rdb/redis-dump-*.rdb
```

---

## 最佳实践

### Loki日志管理

#### 1. 日志标签设计

**推荐标签**:
```yaml
labels:
  - job          # 任务名称
  - app          # 应用名称
  - env          # 环境 (production/staging/dev)
  - component    # 组件 (api/web/worker)
  - tenant_id    # 租户ID (可选)
  - level        # 日志级别
```

**避免高基数标签**:
```yaml
# ❌ Bad: 高基数
labels:
  - request_id   # 每个请求唯一
  - user_id      # 用户数量多
  - trace_id     # 追踪ID唯一

# ✅ Good: 低基数
labels:
  - service      # 服务名称
  - region       # 区域
  - team         # 团队
```

#### 2. 日志格式标准化

**结构化日志** (推荐):
```json
{
  "level": "info",
  "time": "2025-01-03T10:30:45Z",
  "msg": "Request completed",
  "tenant_id": "123",
  "user_id": "456",
  "request_id": "abc-123",
  "duration_ms": 150
}
```

**Promtail解析配置**:
```yaml
pipeline_stages:
  - json:
      expressions:
        - level: level
        - time: time
        - msg: msg
        - tenant_id: tenant_id

  - labels:
      level:
      tenant_id:

  - timestamp:
      source: time
      format: RFC3339
```

#### 3. 保留策略优化

```yaml
# loki-config.yml
limits_config:
  retention_period: 744h  # 31天

table_manager:
  retention_deletes_enabled: true
  retention_period: 744h

compactor:
  retention_enabled: true
  retention_delete_delay: 2h
  working_directory: /loki/boltdb-compact
```

### 数据库备份策略

#### 1. 3-2-1备份原则

- **3**份数据副本 (生产 + 备份 + 异地)
- **2**种存储介质 (本地 + 云存储)
- **1**份异地备份 (灾难恢复)

#### 2. 备份分级策略

| 级别 | 频率 | 类型 | 保留期 |
|------|------|------|--------|
| **L1** | 实时 | Binlog | 7天 |
| **L2** | 每小时 | 增量 | 30天 |
| **L3** | 每天 | 全量 | 90天 |
| **L4** | 每周 | 归档 | 永久 |

#### 3. 恢复时间目标 (RTO/RPO)

| 系统 | RTO | RPO | 策略 |
|------|-----|-----|------|
| MySQL | 1小时 | 15分钟 | 每小时增量 + 实时Binlog |
| Redis | 30分钟 | 6小时 | 每6小时RDB + AOF |
| Elasticsearch | 2小时 | 24小时 | 每天快照 |

#### 4. 异地备份

**使用Rsync同步到异地**:

```bash
# rsync异地备份脚本
#!/bin/bash
SOURCE="/backup"
DESTINATION="backup-server:/remote/backup"
LOG="/var/log/rsync-backup.log"

rsync -avz --delete \
  --progress \
  --log-file="$LOG" \
  "$SOURCE/" "$DESTINATION/"

# 验证同步
rsync -avz --dry-run "$SOURCE/" "$DESTINATION/"
```

**云存储备份** (AWS S3):

```bash
# 使用rclone备份到S3
rclone sync /backup s3://zker-backups/mysql \
  --progress \
  --log-file=/var/log/rclone-backup.log

# 设置生命周期策略 (自动归档)
aws s3api put-bucket-lifecycle-configuration \
  --bucket zker-backups \
  --lifecycle-configuration file:///s3-lifecycle.json
```

---

## 附录

### A. 完整文件清单

#### 日志系统文件

```
docker/volumes/monitoring/
├── loki/
│   └── loki-config.yml              # Loki配置
├── promtail/
│   └── promtail-config.yml          # Promtail配置
└── grafana/provisioning/
    └── datasources/
        └── loki.yml                 # Grafana数据源
```

#### 备份系统文件

```
scripts/backup/
├── backup-mysql.sh                  # MySQL备份脚本
├── backup-redis.sh                  # Redis备份脚本
├── setup-backup-cron.sh             # Cron配置脚本
├── verify-backup.sh                 # 备份验证脚本
├── test-restore.sh                  # 恢复测试脚本
└── README.md                        # 本文档

deploy/k8s/cronjob/
├── mysql-backup.yaml                # MySQL备份CronJob
└── redis-backup.yaml                # Redis备份CronJob
```

### B. 端口清单

| 服务 | 端口 | 用途 |
|------|------|------|
| Loki | 3100 | HTTP API |
| Promtail | 9080 | HTTP Metrics |
| Grafana | 3000 | Web UI |

### C. 参考文档

- [Loki官方文档](https://grafana.com/docs/loki/latest/)
- [Promtail配置](https://grafana.com/docs/loki/latest/clients/promtail/)
- [LogQL查询语法](https://grafana.com/docs/loki/latest/logql/)
- [MySQL备份最佳实践](https://dev.mysql.com/doc/refman/8.0/en/backup-strategy.html)
- [Redis持久化](https://redis.io/docs/manual/persistence/)

---

**创建日期**: 2025-01-03
**维护者**: ZKER DevOps Team
**文档版本**: v1.0.0
