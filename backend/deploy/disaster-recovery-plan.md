# 灾难恢复计划 (Disaster Recovery Plan)

## 📋 目录

1. [恢复目标](#恢复目标)
2. [备份策略](#备份策略)
3. [灾难场景](#灾难场景)
4. [恢复流程](#恢复流程)
5. [演练计划](#演练计划)

---

## 恢复目标

### RPO & RTO

| 服务级别 | RPO (数据丢失容忍度) | RTO (恢复时间目标) |
|---------|-------------------|------------------|
| P0 (核心数据) | 15分钟 | 1小时 |
| P1 (业务数据) | 1小时 | 4小时 |
| P2 (日志数据) | 24小时 | 24小时 |

### 服务优先级

**P0 (关键服务)**:
- 租户服务（tenant-service）
- 认证服务（auth-service）
- 数据库（MySQL集群）

**P1 (重要服务)**:
- 配额服务（quota-service）
- 订阅服务（subscription-service）
- Redis集群

**P2 (辅助服务)**:
- Bot服务（bot-service）
- 工作流服务（workflow-service）
- 日志和监控

---

## 备份策略

### 1. 数据库备份

#### 自动备份（每日）

```bash
#!/bin/bash
# 每日全量备份脚本
# 位置: backend/scripts/mysql-daily-backup.sh

BACKUP_DIR="/backup/mysql/daily"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# 创建备份目录
mkdir -p $BACKUP_DIR

# 全量备份
docker exec mysql-master mysqldump \
  -u root -p$MYSQL_PASSWORD \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --all-databases \
  --flush-logs \
  > $BACKUP_DIR/zker_full_$DATE.sql

# 压缩备份
gzip $BACKUP_DIR/zker_full_$DATE.sql

# 上传到S3
aws s3 cp $BACKUP_DIR/zker_full_$DATE.sql.gz \
  s3://coze-studio-backups/mysql/daily/

# 清理旧备份
find $BACKUP_DIR -name "zker_full_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: zker_full_$DATE.sql.gz"
```

#### Binlog备份（实时）

```bash
#!/bin/bash
# Binlog增量备份
# 每5分钟执行一次

BINLOG_DIR="/backup/mysql/binlog"
DATE=$(date +%Y%m%d)

# 刷新日志
docker exec mysql-master mysql -u root -p$MYSQL_PASSWORD -e "FLUSH LOGS;"

# 复制binlog到备份目录
docker cp mysql-master:/var/lib/mysql/mysql-bin.* $BINLOG_DIR/$DATE/

# 上传到S3
aws s3 sync $BINLOG_DIR/$DATE/ s3://coze-studio-backups/mysql/binlog/$DATE/
```

#### 备份验证

```bash
# 每周验证备份完整性
#!/bin/bash
# backend/scripts/verify-backup.sh

TEST_DB="zker_test"
BACKUP_FILE=$(ls -t /backup/mysql/daily/zker_full_*.sql.gz | head -n 1)

# 解压最新备份
gunzip -c $BACKUP_FILE > /tmp/verify_backup.sql

# 导入到测试数据库
mysql -u root -p$MYSQL_PASSWORD $TEST_DB < /tmp/verify_backup.sql

# 验证表结构
mysql -u root -p$MYSQL_PASSWORD $TEST_DB -e "
  SELECT COUNT(*) AS table_count
  FROM information_schema.tables
  WHERE table_schema = '$TEST_DB';
"

# 验证数据行数
mysql -u root -p$MYSQL_PASSWORD $TEST_DB -e "
  SELECT table_name, table_rows
  FROM information_schema.tables
  WHERE table_schema = '$TEST_DB'
  ORDER BY table_rows DESC;
"

echo "Backup verification completed"
```

### 2. Redis备份

```bash
#!/bin/bash
# Redis集群备份
# backend/scripts/redis-backup.sh

BACKUP_DIR="/backup/redis"
DATE=$(date +%Y%m%d_%H%M%S)

# 为每个节点创建RDB快照
for port in 7001 7002 7003 7004 7005 7006; do
  redis-cli -p $port BGSAVE

  # 等待BGSAVE完成
  while [ $(redis-cli -p $port LASTSAVE) -eq $(redis-cli -p $port LASTSAVE) ]; do
    sleep 1
  done

  # 复制RDB文件
  docker cp redis-cluster-$port:/data/dump.rdb \
    $BACKUP_DIR/dump_$port_$DATE.rdb
done

# 上传到S3
aws s3 sync $BACKUP_DIR/ s3://coze-studio-backups/redis/$DATE/
```

### 3. 配置备份

```bash
#!/bin/bash
# Kubernetes配置备份
# backend/scripts/k8s-config-backup.sh

BACKUP_DIR="/backup/k8s"
DATE=$(date +%Y%m%d_%H%M%S)

# 备份所有命名空间的资源
kubectl get all --all-namespaces -o yaml > $BACKUP_DIR/all_resources_$DATE.yaml

# 备份ConfigMaps和Secrets
kubectl get configmaps --all-namespaces -o yaml > $BACKUP_DIR/configmaps_$DATE.yaml
kubectl get secrets --all-namespaces -o yaml > $BACKUP_DIR/secrets_$DATE.yaml

# 备份Consul KV
curl -s http://consul-ui:8500/v1/kv/?recurse > $BACKUP_DIR/consul_kv_$DATE.json

# 上传到S3
aws s3 sync $BACKUP_DIR/ s3://coze-studio-backups/k8s/$DATE/
```

---

## 灾难场景

### 场景1: 数据中心完全中断

**影响**: 所有服务不可用，数据中心完全丢失

**恢复步骤**:

1. **评估影响**（0-15分钟）
   ```bash
   # 检查服务状态
   kubectl get nodes
   kubectl get pods --all-namespaces

   # 确认数据中心不可达
   ping datacenter-primary.example.com
   ```

2. **启用备用数据中心**（15-30分钟）
   ```bash
   # 切换DNS到备用数据中心
   # 使用GeoDNS或手动更新DNS记录

   # 更新Consul配置指向备用数据库
   kubectl set env deployment/tenant-service \
     DB_HOST=mysql-backup.example.com \
     -n coze-studio

   # 重启服务
   kubectl rollout restart deployment/tenant-service -n coze-studio
   ```

3. **恢复数据**（30分钟-2小时）
   ```bash
   # 从S3下载最新备份
   aws s3 cp s3://coze-studio-backups/mysql/daily/zker_full_20250101.sql.gz /tmp/

   # 恢复到备用数据库
   gunzip < /tmp/zger_full_20250101.sql.gz | \
     mysql -h mysql-backup -u root -p zker

   # 应用binlog增量恢复
   # 从备份时间点到现在
   ```

4. **验证服务**（2-3小时）
   ```bash
   # 执行数据一致性检查
   bash backend/scripts/data-consistency-check.sh

   # 执行烟雾测试
   bash backend/tests/smoke/run-smoke-tests.sh
   ```

### 场景2: 数据库主节点故障

**影响**: 数据写入中断，读取仍可用

**恢复步骤**:

1. **自动故障切换**（0-5分钟）
   ```bash
   # ProxySQL自动检测到Master故障
   # 提升Slave1为新Master

   # 或手动执行故障切换
   bash backend/infra/database/failover.sh
   ```

2. **验证新Master**（5-10分钟）
   ```bash
   # 检查新Master状态
   docker exec mysql-slave1 mysql -e "SHOW MASTER STATUS;"

   # 更新ProxySQL配置
   docker exec proxysql mysql -h 127.0.0.1 -P 6032 -u admin -padmin_password_2025 \
     -e "UPDATE mysql_servers SET hostgroup_id = 10 WHERE hostname = 'mysql-slave1';"
   ```

3. **重新配置旧Master为Slave**（10-30分钟）
   ```bash
   # 在旧Master上执行
   docker exec mysql-master mysql -e "
     STOP SLAVE;
     CHANGE MASTER TO
       MASTER_HOST='mysql-slave1',
       MASTER_USER='repl',
       MASTER_PASSWORD='repl_password',
       MASTER_LOG_FILE='mysql-bin.000001',
       MASTER_LOG_POS=154;
     START SLAVE;
   "
   ```

### 场景3: Redis集群故障

**影响**: 缓存丢失，会话中断，但数据不会丢失

**恢复步骤**:

1. **检查集群状态**（0-5分钟）
   ```bash
   redis-cli -c -p 7001 cluster info
   redis-cli -c -p 7001 cluster nodes
   ```

2. **重启失败节点**（5-15分钟）
   ```bash
   # 重启Redis Pod
   kubectl delete pod redis-cluster-0 -n database

   # 等待Pod重启
   kubectl wait --for=condition=ready pod -l app=redis-cluster -n database
   ```

3. **恢复数据（如需要）**（15-30分钟）
   ```bash
   # 从备份恢复RDB文件
   aws s3 cp s3://coze-studio-backups/redis/latest/dump_7001.rdb /tmp/

   # 复制到Redis容器
   docker cp /tmp/dump_7001.rdb redis-cluster-0:/data/dump.rdb

   # 重启Redis
   kubectl delete pod redis-cluster-0 -n database
   ```

### 场景4: Kubernetes集群故障

**影响**: 所有Pod不可用

**恢复步骤**:

1. **检查etcd健康**（0-10分钟）
   ```bash
   # 检查etcd成员
   kubectl get etcdmembers

   # 如果etcd损坏，从备份恢复
   ETCDCTL_API=3 etcdctl snapshot restore /backup/etcd/snapshot.db
   ```

2. **重建控制平面**（10-30分钟）
   ```bash
   # 如果控制平面损坏
   kubeadm init --config=kubeadm-config.yaml
   ```

3. **恢复工作负载**（30分钟-2小时）
   ```bash
   # 从备份恢复配置
   kubectl apply -f /backup/k8s/all_resources_latest.yaml
   kubectl apply -f /backup/k8s/configmaps_latest.yaml
   kubectl apply -f /backup/k8s/secrets_latest.yaml
   ```

### 场景5: 数据损坏或误删

**影响**: 数据不一致或丢失

**恢复步骤**:

1. **评估损坏范围**（0-15分钟）
   ```bash
   # 执行数据一致性检查
   bash backend/scripts/data-consistency-check.sh

   # 确定损坏的时间点
   # 查看binlog或审计日志
   ```

2. **时间点恢复（PITR）**（15分钟-2小时）
   ```bash
   # 1. 恢复全量备份
   gunzip < /backup/mysql/daily/zker_full_20250101.sql.gz | \
     mysql -u root -p zker_restored

   # 2. 应用binlog到损坏前的时间点
   mysqlbinlog \
     --start-datetime="2025-01-01 00:00:00" \
     --stop-datetime="2025-01-01 10:30:00" \
     /backup/mysql/binlog/mysql-bin.000001 | \
     mysql -u root -p zker_restored

   # 3. 验证恢复的数据
   mysql -u root -p zker_restored -e "SELECT COUNT(*) FROM tenants;"
   ```

3. **切换到恢复的数据库**（2-3小时）
   ```bash
   # 重命名数据库
   mysql -u root -p -e "
     RENAME DATABASE zker TO zker_corrupted;
     RENAME DATABASE zker_restored TO zker;
   "

   # 验证应用
   kubectl rollout restart deployment/tenant-service -n coze-studio
   ```

---

## 恢复流程

### 通用恢复流程图

```
灾难发生
    ↓
检测和评估（15分钟）
    ↓
选择恢复策略
    ↓
执行恢复
    ├→ 备份恢复
    ├→ 故障切换
    └→ 降级服务
    ↓
验证恢复（1小时）
    ↓
监控稳定性（24小时）
    ↓
事后分析
```

### 紧急联系清单

| 角色 | 姓名 | 电话 | 邮箱 |
|-----|------|------|------|
| On-Call工程师 | 工程师A | 138-xxxx-xxxx | a@coze-studio.com |
| 后端负责人 | 架构师B | 139-xxxx-xxxx | b@coze-studio.com |
| DevOps工程师 | 运维C | 136-xxxx-xxxx | c@coze-studio.com |
| CTO | 技术总监D | 135-xxxx-xxxx | d@coze-studio.com |

---

## 演练计划

### 季度演练（每季度一次）

#### Q1: 数据中心故障演练（3月）

```bash
# 模拟主数据中心故障
kubectl cordon node-primary-1
kubectl cordon node-primary-2
kubectl cordon node-primary-3

# 执行故障切换
bash backend/deploy/dr/failover-to-backup.sh

# 验证服务恢复
bash backend/tests/smoke/run-smoke-tests.sh

# 恢复原状
bash backend/deploy/dr/failback-to-primary.sh
```

#### Q2: 数据库故障演练（6月）

```bash
# 模拟Master故障
kubectl delete pod mysql-cluster-0 -n database

# 观察自动故障切换
# ProxySQL应该在30秒内切换到Slave

# 手动干预（如果自动切换失败）
bash backend/infra/database/failover.sh

# 验证数据一致性
bash backend/scripts/data-consistency-check.sh
```

#### Q3: 数据损坏恢复演练（9月）

```bash
# 模拟误删数据
kubectl exec -it mysql-cluster-0 -n database -- mysql -e "
  DELETE FROM zker.tenants WHERE tenant_id = 'test_tenant';
"

# 执行时间点恢复
bash backend/deploy/dr/pitr-recovery.sh \
  --database=zker \
  --until="2025-09-15 14:30:00"

# 验证恢复
kubectl exec -it mysql-cluster-0 -n database -- mysql -e "
  SELECT * FROM zker.tenants WHERE tenant_id = 'test_tenant';
"
```

#### Q4: 全系统灾难演练（12月）

```bash
# 关闭所有集群
kubectl scale deployment --all --replicas=0 --all-namespaces

# 执行完整恢复流程
bash backend/deploy/dr/full-recovery.sh

# 验证所有服务
bash backend/tests/integration/run-all-tests.sh

# 性能测试
k6 run backend/tests/performance/full-load-test.js
```

### 年度演练（每年一次）

**演练目标**: 验证异地多活切换

**演练内容**:
1. 完全切换到异地数据中心
2. 在异地数据中心运行24小时
3. 验证数据同步延迟 < 1秒
4. 验证所有核心功能正常
5. 执行回切流程

---

## 监控和告警

### DR相关指标

```yaml
# Prometheus告警规则
groups:
- name: disaster_recovery
  rules:
  # 备份失败告警
  - alert: BackupFailed
    expr: backup_success == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Backup failed for {{ $labels.service }}"
      description: "Backup has been failing for more than 5 minutes"

  # 数据复制延迟告警
  - alert: ReplicationLag
    expr: mysql_slave_lag_seconds > 60
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "MySQL replication lag: {{ $value }}s"

  # etcd集群不健康
  - alert: EtcdClusterUnhealthy
    expr: etcd_server_health_status < 1
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "etcd cluster member {{ $labels.instance }} is unhealthy"
```

---

## 文档维护

**更新频率**: 每季度或架构变更后

**维护责任人**: DevOps团队负责人

**最后更新**: 2025-01-01

**下次审核**: 2025-04-01
