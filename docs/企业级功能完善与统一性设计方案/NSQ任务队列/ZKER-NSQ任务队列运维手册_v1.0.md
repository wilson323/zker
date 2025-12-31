# ZKER NSQ任务队列运维手册

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**目标读者**: 运维工程师、SRE

---

## 📋 目录

- [1. 部署指南](#1-部署指南)
- [2. 配置管理](#2-配置管理)
- [3. 监控告警](#3-监控告警)
- [4. 日常运维](#4-日常运维)
- [5. 故障处理](#5-故障处理)
- [6. 性能调优](#6-性能调优)
- [7. 备份恢复](#7-备份恢复)
- [8. 安全加固](#8-安全加固)

---

## 1. 部署指南

### 1.1 Docker部署

**Docker Compose**:
```yaml
version: '3.8'

services:
  nsqlookupd:
    image: nsqio/nsq:v1.2.1
    command: /nsqlookupd
    ports:
      - "4160:4160"  # TCP
      - "4161:4161"  # HTTP
    volumes:
      - nsqlookupd-data:/data
    networks:
      - nsq

  nsqd:
    image: nsqio/nsq:v1.2.1
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    ports:
      - "4150:4150"  # TCP
      - "4151:4151"  # HTTP
    volumes:
      - nsqd-data:/data
    depends_on:
      - nsqlookupd
    networks:
      - nsq

  nsqadmin:
    image: nsqio/nsq:v1.2.1
    command: /nsqadmin --lookupd-http-address=nsqlookupd:4161
    ports:
      - "4171:4171"
    depends_on:
      - nsqlookupd
    networks:
      - nsq

volumes:
  nsqlookupd-data:
  nsqd-data:

networks:
  nsq:
    driver: bridge
```

**启动命令**:
```bash
docker-compose up -d
```

### 1.2 Kubernetes部署

**Namespace**:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: nsq
```

**NSQLookupD Deployment**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nsqlookupd
  namespace: nsq
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nsqlookupd
  template:
    metadata:
      labels:
        app: nsqlookupd
    spec:
      containers:
      - name: nsqlookupd
        image: nsqio/nsq:v1.2.1
        command:
          - /nsqlookupd
        ports:
        - containerPort: 4160  # TCP
          name: tcp
        - containerPort: 4161  # HTTP
          name: http
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /ping
            port: 4161
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ping
            port: 4161
          initialDelaySeconds: 5
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: nsqlookupd
  namespace: nsq
spec:
  selector:
    app: nsqlookupd
  ports:
  - port: 4160
    targetPort: 4160
    name: tcp
  - port: 4161
    targetPort: 4161
    name: http
  type: ClusterIP
```

**NSQD StatefulSet**:
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nsqd
  namespace: nsq
spec:
  serviceName: nsqd
  replicas: 3
  selector:
    matchLabels:
      app: nsqd
  template:
    metadata:
      labels:
        app: nsqd
    spec:
      containers:
      - name: nsqd
        image: nsqio/nsq:v1.2.1
        command:
          - /nsqd
          - --lookupd-tcp-address=nsqlookupd:4160
          - --broadcast-address=$(POD_IP)
        env:
          - name: POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
        ports:
        - containerPort: 4150
          name: tcp
        - containerPort: 4151
          name: http
        volumeMounts:
        - name: data
          mountPath: /data
        resources:
          requests:
            memory: "256Mi"
            cpu: "200m"
          limits:
            memory: "512Mi"
            cpu: "500m"
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: fast-ssd
      resources:
        requests:
          storage: 100Gi
```

### 1.3 集群部署

**多可用区部署**:
```
AZ-A: nsqlookupd-1, nsqd-1, nsqd-2
AZ-B: nsqlookupd-2, nsqd-3, nsqd-4
AZ-C: nsqlookupd-3, nsqd-5, nsqd-6
```

**负载均衡**:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: nsqlookupd-lb
  namespace: nsq
spec:
  selector:
    app: nsqlookupd
  ports:
  - port: 4160
    targetPort: 4160
  type: LoadBalancer
```

---

## 2. 配置管理

### 2.1 生产环境配置

**backend/conf/queue/nsq.yaml**:
```yaml
nsqd_addresses:
  - "nsqd-1.prod.example.com:4150"
  - "nsqd-2.prod.example.com:4150"
  - "nsqd-3.prod.example.com:4150"

nsqlookupd_addresses:
  - "nsqlookupd-1.prod.example.com:4161"
  - "nsqlookupd-2.prod.example.com:4161"
  - "nsqlookupd-3.prod.example.com:4161"

producer:
  timeout: 10s
  retry: 5
  retry_interval: 2s
  max_idle_conns: 50

consumer:
  max_retries: 5
  retry_delay: 2s
  max_in_flight: 20
  message_timeout: 300s
  heartbeat_interval: 30s

dlq:
  enabled: true
  topic_suffix: "-dlq"
  retention_duration: 168h  # 7天
  max_replays: 5

metrics:
  enabled: true
  address: ":2112"

logging:
  level: "warn"
  verbose: false
```

### 2.2 环境变量配置

**ConfigMap**:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nsq-config
  namespace: nsq
data:
  NSQD_ADDRESSES: "nsqd-1:4150,nsqd-2:4150,nsqd-3:4150"
  NSQLOOKUPD_ADDRESSES: "nsqlookupd-1:4161,nsqlookupd-2:4161,nsqlookupd-3:4161"
  PRODUCER_TIMEOUT: "10s"
  CONSUMER_MAX_RETRIES: "5"
  DLQ_ENABLED: "true"
```

### 2.3 配置热更新

**使用ConfigMap + Volume**:
```yaml
spec:
  containers:
  - name: app
    volumeMounts:
    - name: config
      mountPath: /etc/nsq
  volumes:
  - name: config
    configMap:
      name: nsq-config
```

**重载配置**:
```bash
# 更新ConfigMap
kubectl apply -f nsq-configmap.yaml

# 触发滚动更新
kubectl rollout restart deployment/queue-manager
```

---

## 3. 监控告警

### 3.1 Prometheus监控

**Prometheus配置**:
```yaml
scrape_configs:
  - job_name: 'nsq'
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - nsq
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: nsqd
      - source_labels: [__meta_kubernetes_pod_ip]
        target_label: __address__
        replacement: $1:4151
```

**关键指标**:
```promql
# 队列深度
nsq_queue_depth{topic=~".*"}

# 消息速率
rate(nsq_message_count[5m])

# 消费者健康度
nsq_clients_count{topic=~".*"}

# 死信队列大小
nsq_dlq_messages_stored_total
```

### 3.2 Grafana仪表板

**主要面板**:
1. **Overview**: 总体健康度
2. **Topics**: Topic详情和深度
3. **Producers**: 生产者指标
4. **Consumers**: 消费者指标
5. **DLQ**: 死信队列统计

**导入仪表板**:
```json
{
  "dashboard": {
    "title": "NSQ Dashboard",
    "panels": [
      {
        "title": "Queue Depth",
        "targets": [
          {
            "expr": "nsq_queue_depth"
          }
        ]
      },
      {
        "title": "Message Rate",
        "targets": [
          {
            "expr": "rate(nsq_message_count[5m])"
          }
        ]
      }
    ]
  }
}
```

### 3.3 告警规则

**alerting rules**:
```yaml
groups:
  - name: nsq_alerts
    interval: 30s
    rules:
      # 队列深度告警
      - alert: NSQQueueDepthHigh
        expr: nsq_queue_depth > 10000
        for: 10m
        labels:
          severity: warning
          team: platform
        annotations:
          summary: "NSQ queue depth too high"
          description: "Queue {{ $labels.topic }} depth is {{ $value }}"

      # 消费者离线
      - alert: NSQConsumerOffline
        expr: nsq_clients_count == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "No consumers for topic"
          description: "Topic {{ $labels.topic }} has no active consumers"

      # DLQ增长过快
      - alert: NSQDLQGrowthRate
        expr: rate(nsq_dlq_messages_stored_total[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "DLQ growing too fast"
          description: "DLQ for {{ $labels.topic }} growing at {{ $value }}/s"

      # 处理延迟过高
      - alert: NSQProcessingLatency
        expr: histogram_quantile(0.99, rate(nsq_consumer_processing_time_seconds_bucket[5m])) > 60
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Processing latency too high"
          description: "P99 latency is {{ $value }}s"
```

---

## 4. 日常运维

### 4.1 健康检查

**NSQD健康检查**:
```bash
curl http://localhost:4151/ping
# 输出: OK

curl http://localhost:4151/stats
# 查看详细统计信息
```

**应用健康检查**:
```bash
curl http://localhost:8080/health
# 返回: {"status":"healthy","queue_manager":{"started":true}}
```

### 4.2 日志管理

**查看NSQ日志**:
```bash
# Docker
docker logs nsqd-1 -f

# Kubernetes
kubectl logs -f deployment/nsqd -n nsq

# 查看特定时间范围
kubectl logs --since-time=2025-01-01T00:00:00Z -n nsq
```

**日志聚合**:
```yaml
# Fluentd配置
<source>
  @type tail
  path /var/log/nsq/*.log
  pos_file /var/log/fluentd-nsq.pos
  tag nsq.*
  <parse>
    @type json
  </parse>
</source>

<match nsq.**>
  @type elasticsearch
  host elasticsearch
  port 9200
  logstash_format true
  logstash_prefix nsq
</match>
```

### 4.3 数据清理

**清理过期数据**:
```bash
# NSQD数据目录
cd /var/lib/nsq

# 查看磁盘使用
du -sh *

# 删除旧Topic数据（谨慎操作）
rm -rf topics/old_topic
```

**死信队列清理**:
```go
// 定期清理过期DLQ消息
func CleanupExpiredDLQ(ctx context.Context) {
    dlq := manager.GetDLQ()
    err := dlq.CleanupExpiredMessages(ctx)
    if err != nil {
        logger.Error("failed to cleanup DLQ", zap.Error(err))
    }
}
```

### 4.4 性能检查

**查看队列状态**:
```bash
# 使用nsq_stat
go install github.com/nsqio/nsq/stat/nsq_stat@latest
nsq_stat --lookupd-http-address=localhost:4161

# 输出示例:
# +----------+-------------+----------+----------+
# | Topic    | Depth       | Messages | Clients  |
# +----------+-------------+----------+----------+
# | workflow | 1234        | 45678    | 3        |
# | bot      | 567         | 23456    | 2        |
# +----------+-------------+----------+----------+
```

**压测工具**:
```bash
# 使用nsq_tail
go install github.com/nsqio/nsq/apps/nsq_tail@latest
nsq_tail --topic=test_topic --channel=test --nsqd-tcp-address=localhost:4150

# 使用nsq_to_file
nsq_to_file --topic=test_topic --output=/tmp/test.log
```

---

## 5. 故障处理

### 5.1 NSQD节点故障

**症状**:
- 消息发布失败
- 部分消费者离线

**处理步骤**:
```bash
# 1. 确认故障节点
kubectl get pods -n nsq -o wide

# 2. 查看日志
kubectl logs nsqd-1 -n nsq --tail=100

# 3. 重启节点
kubectl delete pod nsqd-1 -n nsq

# 4. 验证恢复
kubectl get pods -n nsq
```

### 5.2 消息堆积

**症状**:
- Queue深度持续增长
- 消费延迟增加

**处理步骤**:
```bash
# 1. 确认堆积情况
curl http://localhost:4161/stats | jq '.data.topics[] | select(.depth > 1000)'

# 2. 检查消费者状态
curl http://localhost:4161/stats | jq '.data.topics[].clients'

# 3. 临时扩容消费者
kubectl scale deployment queue-consumer --replicas=10 -n nsq

# 4. 监控恢复
watch -n 5 'curl -s http://localhost:4161/stats | jq .depth'
```

### 5.3 死信队列爆满

**症状**:
- DLQ消息数快速增长
- 大量任务失败

**处理步骤**:
```bash
# 1. 分析失败原因
curl http://localhost:4161/stats | jq '.data.topics[] | select(.name | contains("-dlq"))'

# 2. 查看DLQ消息
# (需要实现DLQ查询接口)

# 3. 修复代码后重放
curl -X POST http://localhost:8080/api/v1/queue/replay-dlq \
  -H "Content-Type: application/json" \
  -d '{"topic": "knowledge_document"}'

# 4. 监控重放进度
curl http://localhost:8080/api/v1/queue/stats
```

### 5.4 磁盘空间不足

**症状**:
- NSQD写入失败
- 日志显示 "no space left on device"

**处理步骤**:
```bash
# 1. 检查磁盘使用
df -h /var/lib/nsq

# 2. 清理旧数据
# 方法1: 删除旧Topic
nsq_admin --topic=old_topic --empty

# 方法2: 调整数据保留
# 修改nsqd启动参数
--data-path=/mnt/large-disk/nsq

# 3. 扩容磁盘
kubectl patch pvc nsqd-data-nsqd-0 \
  -p '{"spec":{"resources":{"requests":{"storage":"200Gi"}}}}'
```

---

## 6. 性能调优

### 6.1 NSQD调优

**启动参数**:
```bash
nsqd \
  --lookupd-tcp-address=nsqlookupd:4160 \
  --mem-queue-size=10000 \
  --max-body-size=10485760 \
  --max-msg-size=1048576 \
  --max-msg-timeout=15m \
  --msg-timeout=60s \
  --e2e-processing-latency-percentile=99.9
```

**参数说明**:
- `--mem-queue-size`: 内存队列大小（0=全持久化）
- `--max-body-size`: 最大消息体大小（10MB）
- `--max-msg-size`: 最大消息大小（1MB）
- `--max-msg-timeout`: 最大消息超时（15分钟）
- `--msg-timeout`: 默认消息超时（60秒）

### 6.2 消费者调优

**配置优化**:
```yaml
consumer:
  max_in_flight: 20      # 增加并发度
  message_timeout: 300s  # 延长超时
  heartbeat_interval: 30s  # 心跳间隔
```

**水平扩展**:
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: queue-consumer
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: queue-consumer
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Pods
    pods:
      metric:
        name: queue_depth
      target:
        type: AverageValue
        averageValue: "100"
```

### 6.3 网络优化

**TCP调优**:
```bash
# /etc/sysctl.conf
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 8192
net.core.netdev_max_backlog = 16384

# 应用配置
sysctl -p
```

**KeepAlive**:
```yaml
nsqd:
  --tcp-max-addresses=3
  --tls-min-version=1.2
  --tls-required=true
```

---

## 7. 备份恢复

### 7.1 数据备份

**NSQD数据备份**:
```bash
#!/bin/bash
# backup_nsq.sh

BACKUP_DIR="/backup/nsq/$(date +%Y%m%d)"
NSQ_DATA_DIR="/var/lib/nsq"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 停止NSQD（可选，保证一致性）
# systemctl stop nsqd

# 备份数据
rsync -av $NSQ_DATA_DIR/ $BACKUP_DIR/

# 压缩
tar -czf $BACKUP_DIR.tar.gz $BACKUP_DIR/

# 上传到S3
aws s3 cp $BACKUP_DIR.tar.gz s3://backups/nsq/

# 清理本地
rm -rf $BACKUP_DIR

echo "Backup completed: $BACKUP_DIR.tar.gz"
```

**定时备份**:
```cron
# /etc/cron.d/nsq-backup
0 2 * * * root /usr/local/bin/backup_nsq.sh
```

### 7.2 数据恢复

**恢复步骤**:
```bash
#!/bin/bash
# restore_nsq.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: restore_nsq.sh <backup_file>"
    exit 1
fi

# 停止NSQD
systemctl stop nsqd

# 解压备份
tar -xzf $BACKUP_FILE -C /tmp/

# 恢复数据
rsync -av /tmp/$(basename $BACKUP_FILE .tar.gz)/ /var/lib/nsq/

# 启动NSQD
systemctl start nsqd

echo "Restore completed"
```

### 7.3 灾难恢复

**跨区域复制**:
```bash
# 方案1: 使用Rsync
rsync -avz -e ssh \
  /var/lib/nsq/ \
  user@dr-site:/var/lib/nsq/

# 方案2: 使用对象存储
aws s3 sync s3://nsq-prod/ s3://nsq-dr/
```

---

## 8. 安全加固

### 8.1 认证授权

**启用认证**:
```bash
nsqd \
  --auth-http-address=http://auth-service:8080 \
  --auth-http-request-header=X-Auth-Token
```

**认证服务**:
```go
func authHandler(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("X-Auth-Token")

    // 验证Token
    if !validateToken(token) {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    // 返回授权信息
    auth := map[string]interface{}{
        "ttl": 3600,
        "topics": []string{"tenant_123_.*"},
        "channels": []string{".*"},
    }

    json.NewEncoder(w).Encode(auth)
}
```

### 8.2 TLS加密

**生成证书**:
```bash
# CA证书
openssl genrsa -out ca-key.pem 2048
openssl req -x509 -new -nodes -key ca-key.pem \
  -days 1000 -out ca.pem -subj "/CN=nsq_ca"

# 服务器证书
openssl genrsa -out server-key.pem 2048
openssl req -new -key server-key.pem -out server.csr \
  -subj "/CN=nsq.example.com"
openssl x509 -req -in server.csr -CA ca.pem -CAkey ca-key.pem \
  -CAcreateserial -out server-cert.pem -days 365
```

**启用TLS**:
```bash
nsqd \
  --tls-cert=/path/to/server-cert.pem \
  --tls-key=/path/to/server-key.pem \
  --tls-client-auth-policy=require \
  --tls-root-ca-file=/path/to/ca.pem
```

### 8.3 网络隔离

**NetworkPolicy**:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: nsq-network-policy
  namespace: nsq
spec:
  podSelector:
    matchLabels:
      app: nsqd
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: backend
    ports:
    - protocol: TCP
      port: 4150
  egress:
  - to:
    - namespaceSelector:
        matchLabels:
          name: nsq
    ports:
    - protocol: TCP
      port: 4160
```

---

## 附录

### A. 常用命令

```bash
# 查看Topic列表
curl http://localhost:4161/stats | jq '.data.topics[].name'

# 查看特定Topic
curl http://localhost:4161/stats | jq '.data.topics[] | select(.name=="workflow")'

# 清空Topic
curl -X POST http://localhost:4151/topic_empty?topic=workflow

# 删除Topic
curl -X POST http://localhost:4151/topic_delete?topic=workflow

# 查看Channel
curl http://localhost:4161/stats | jq '.data.topics[].channels'

# 清空Channel
curl -X POST http://localhost:4151/channel_empty?topic=workflow&channel=worker1
```

### B. 监控脚本

```bash
#!/bin/bash
# monitor_nsq.sh

NSQLOOKUPD="http://localhost:4161"

while true; do
  clear
  echo "NSQ Queue Monitor - $(date)"
  echo "================================"

  curl -s $NSQLOOKUPD/stats | \
    jq -r '.data.topics[] | "\(.name | lpad(30)) Depth: \(.depth | lpad(8)) Messages: \(.message_count | lpad(10)) Clients: \(.clients | length | lpad(3))"'

  sleep 5
done
```

### C. 相关文档

- 架构设计: [ZKER-NSQ任务队列架构设计文档_v1.0.md](./ZKER-NSQ任务队列架构设计文档_v1.0.md)
- 使用指南: [ZKER-NSQ任务队列使用指南_v1.0.md](./ZKER-NSQ任务队列使用指南_v1.0.md)

---

**© 2025 ZKER Project. All rights reserved.**
