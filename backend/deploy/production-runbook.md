# 生产环境部署手册 (Production Deployment Runbook)

## 📋 目录

1. [部署前检查](#部署前检查)
2. [分步部署指南](#分步部署指南)
3. [回滚流程](#回滚流程)
4. [故障排查](#故障排查)

---

## 部署前检查

### 1. 环境检查清单

**基础设施**:
- [ ] Kubernetes集群已就绪（至少3个节点）
- [ ] Helm 3.x 已安装
- [ ] kubectl 已配置并可访问集群
- [ ] Docker镜像仓库可访问
- [ ] 监控系统（Prometheus + Grafana）已部署
- [ ] 日志系统（ELK）已部署
- [ ] 分布式追踪（Jaeger）已部署

**配置检查**:
- [ ] 所有环境变量已配置（.env文件）
- [ ] 数据库连接配置已验证
- [ ] Redis集群配置已验证
- [ ] Consul配置已验证
- [ ] SSL证书已准备

**数据检查**:
- [ ] 数据库备份已完成
- [ ] Redis数据已备份（如需要）
- [ ] 数据迁移脚本已准备

### 2. 健康检查端点

所有服务必须实现以下健康检查端点：

```
GET /health/live  - 存活探针（检查服务是否运行）
GET /health/ready - 就绪探针（检查服务是否可以接收流量）
GET /health/startup - 启动探针（检查服务是否启动成功）
```

响应示例：
```json
{
  "status": "healthy",
  "timestamp": "2025-01-01T10:00:00Z",
  "checks": {
    "database": "ok",
    "redis": "ok",
    "consul": "ok"
  }
}
```

---

## 分步部署指南

### 阶段1: 基础设施部署（第1-2天）

#### 步骤1.1: 部署Consul集群

```bash
# 添加Consul Helm仓库
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# 创建命名空间
kubectl create namespace consul

# 部署Consul集群（3个节点）
helm install consul hashicorp/consul \
  --namespace consul \
  --set global.name=consul \
  --set server.replicas=3 \
  --set ui.enabled=true \
  --set ui.service.type=LoadBalancer \
  --set connectInject.enabled=true \
  --set client.enabled=true \
  --set meshGateway.enabled=true

# 验证部署
kubectl get pods -n consul
kubectl port-forward svc/consul-ui 8500:80 -n consul
# 访问: http://localhost:8500
```

#### 步骤1.2: 部署MySQL主从集群

```bash
# 创建命名空间
kubectl create namespace database

# 部署MySQL Operator
helm repo presslabs https://presslabs.github.io/charts
helm install presslabs-mysql presslabs/mysql-operator \
  --namespace database

# 部署MySQL主从集群
cat <<EOF | kubectl apply -f -
apiVersion: mysql.presslabs.org/v1alpha1
kind: MysqlCluster
metadata:
  name: mysql-cluster
  namespace: database
spec:
  replicas: 3
  secretName: mysql-secret
  podSpec:
    image: mysql:8.4.5
    resources:
      requests:
        memory: 2Gi
        cpu: 1000m
  volumeSpec:
    persistentVolumeClaim:
      storageClassName: fast-ssd
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 100Gi
EOF

# 等待集群就绪
kubectl wait --for=condition=ready mysqlcluster/mysql-cluster -n database --timeout=600s

# 验证主从复制
kubectl exec -n database mysql-cluster-0 -- mysql -e "SHOW MASTER STATUS"
kubectl exec -n database mysql-cluster-1 -- mysql -e "SHOW SLAVE STATUS"
```

#### 步骤1.3: 部署Redis集群

```bash
# 部署Redis集群（3主3从）
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis-cluster bitnami/redis-cluster \
  --namespace database \
  --set auth.enabled=true \
  --set auth.password=your_redis_password \
  --set cluster.nodes=6 \
  --set cluster.replicas.replicationFactor=2 \
  --set persistence.enabled=true \
  --set persistence.size=50Gi \
  --set metrics.enabled=true

# 验证集群状态
kubectl exec -it -n database redis-cluster-0 -- redis-cli -c cluster info
```

#### 步骤1.4: 部署ProxySQL（读写分离）

```bash
# 部署ProxySQL
helm install proxysql ./.helm/proxysql \
  --namespace database \
  --set mysql.servers=mysql-cluster-0.mysql-cluster.database.svc.cluster.local,mysql-cluster-1.mysql-cluster.database.svc.cluster.local,mysql-cluster-2.mysql-cluster.database.svc.cluster.local \
  --set mysql.user=root \
  --set mysql.password=your_mysql_password

# 验证读写分离配置
kubectl exec -n database proxysql-0 -- mysql -h 127.0.0.1 -P 6032 -u admin -padmin_password_2025 -e "SELECT * FROM mysql_servers;"
```

### 阶段2: 监控和日志系统部署（第3天）

#### 步骤2.1: 部署Prometheus

```bash
# 添加Prometheus Operator
helm repo add prometheus-community https://prometheus-community.github.io/helm
helm repo update

# 创建命名空间
kubectl create namespace monitoring

# 部署Prometheus Operator
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring \
  --set grafana.adminPassword=admin \
  --set prometheus.prometheusSpec.retention=30d \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.storageClassName=fast-ssd \
  --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=200Gi

# 访问Grafana
kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring
# 访问: http://localhost:3000 (admin/admin)
```

#### 步骤2.2: 部署Jaeger（分布式追踪）

```bash
# 部署Jaeger Operator
kubectl create namespace observability
kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.50.0/jaeger-operator.yaml -n observability

# 部署Jaeger实例
cat <<EOF | kubectl apply -f -
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger-prod
  namespace: observability
spec:
  strategy: production
  storage:
    type: elasticsearch
    elasticsearch:
      nodeCount: 3
      storage: 100Gi
      resources:
        requests:
          cpu: 200m
          memory: 2Gi
  ingress:
    enabled: true
    security: none
EOF

# 访问Jaeger UI
kubectl port-forward svc/jaeger-prod-query 16686:16686 -n observability
# 访问: http://localhost:16686
```

#### 步骤2.3: 部署ELK Stack（日志）

```bash
# 部署Elasticsearch
helm install elasticsearch elastic/elasticsearch \
  --namespace logging \
  --set replicas=3 \
  --set volumeClaimTemplate.resources.requests.storage=200Gi \
  --set resources.requests.cpu=2 \
  --set resources.requests.memory=8Gi

# 部署Kibana
helm install kibana elastic/kibana \
  --namespace logging \
  --set service.type=LoadBalancer

# 部署Fluentd（日志收集）
helm install fluentd fluent/fluentd \
  --namespace logging \
  --set image.tag=v1.16-1 \
  --set flushInterval=10 \
  --set resources.requests.cpu=200m \
  --set resources.requests.memory=500Mi
```

### 阶段3: 微服务部署（第4-5天）

#### 步骤3.1: 构建并推送Docker镜像

```bash
# 构建镜像
docker build -t coze-studio/tenant-service:v1.0.0 -f docker/tenant-service/Dockerfile .
docker build -t coze-studio/quota-service:v1.0.0 -f docker/quota-service/Dockerfile .
docker build -t coze-studio/subscription-service:v1.0.0 -f docker/subscription-service/Dockerfile .
docker build -t coze-studio/auth-service:v1.0.0 -f docker/auth-service/Dockerfile .

# 推送到镜像仓库
docker push coze-studio/tenant-service:v1.0.0
docker push coze-studio/quota-service:v1.0.0
docker push coze-studio/subscription-service:v1.0.0
docker push coze-studio/auth-service:v1.0.0
```

#### 步骤3.2: 部署租户服务

```bash
# 部署租户服务
kubectl apply -f k8s/tenant-service-deployment.yaml

# 监控部署状态
kubectl rollout status deployment/tenant-service -n coze-studio

# 查看Pod状态
kubectl get pods -n coze-studio -l app=tenant-service

# 查看日志
kubectl logs -f -n coze-studio -l app=tenant-service --all-containers=true

# 验证服务
kubectl port-forward svc/tenant-service 8001:8001 -n coze-studio
curl http://localhost:8001/health/ready
```

#### 步骤3.3: 部署其他服务

```bash
# 依次部署所有服务
kubectl apply -f k8s/quota-service-deployment.yaml
kubectl apply -f k8s/subscription-service-deployment.yaml
kubectl apply -f k8s/auth-service-deployment.yaml
kubectl apply -f k8s/bot-service-deployment.yaml

# 监控所有部署
kubectl get deployments -n coze-studio
kubectl get pods -n coze-studio
```

#### 步骤3.4: 部署Kong API网关

```bash
# 部署Kong
helm repo add kong https://charts.konghq.com
helm repo update

helm install kong kong/kong \
  --namespace coze-studio \
  --set env.database=postgres \
  --set postgresql.enabled=true \
  --set postgresql.postgresqlPassword=kong_postgres_password \
  --set ingressController.enabled=true \
  --set ingressController.installCRDs=false

# 应用Kong配置
kubectl apply -f docker/kong/kong-services.yml

# 验证Kong
kubectl port-forward svc/kong-proxy 8080:80 -n coze-studio
curl http://localhost:8080/health
```

### 阶段4: 流量切换（第6天）

#### 步骤4.1: 灰度发布（Canary Deployment）

```bash
# 使用Istio进行灰度发布
kubectl apply -f k8s/istio-canary.yaml

# 配置灰度规则（10%流量到新版本）
cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: tenant-service-canary
  namespace: coze-studio
spec:
  hosts:
  - tenant-service
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: tenant-service
        subset: v2
      weight: 100
  - route:
    - destination:
        host: tenant-service
        subset: v1
      weight: 90
    - destination:
        host: tenant-service
        subset: v2
      weight: 10
EOF

# 监控新版本指标
kubectl port-forward svc/prometheus-grafana 3000:80 -n monitoring
# 访问Grafana查看服务指标
```

#### 步骤4.2: 逐步增加流量

```bash
# 50%流量
kubectl apply -f k8s/canary-50%.yaml

# 等待15分钟观察

# 100%流量
kubectl apply -f k8s/canary-100%.yaml

# 移除旧版本
kubectl delete deployment tenant-service-v1 -n coze-studio
```

### 阶段5: 部署后验证（第7天）

#### 步骤5.1: 数据一致性检查

```bash
# 执行数据一致性检查脚本
bash backend/scripts/data-consistency-check.sh
```

#### 步骤5.2: 性能测试

```bash
# 执行性能测试
k6 run backend/tests/performance/api-load-test.js
k6 run backend/tests/performance/db_readwrite_benchmark.sh
```

#### 步骤5.3: 监控验证

```bash
# 检查Prometheus指标
curl http://prometheus-server:9090/api/v1/query?query=up{job="tenant-service"}

# 检查Jaeger追踪
curl http://jaeger-query:16686/api/traces?service=tenant-service&limit=10
```

---

## 回滚流程

### 自动回滚（K8s原生）

```bash
# 回滚到上一个版本
kubectl rollout undo deployment/tenant-service -n coze-studio

# 回滚到特定版本
kubectl rollout undo deployment/tenant-service -n coze-studio --to-revision=3

# 查看回滚历史
kubectl rollout history deployment/tenant-service -n coze-studio
```

### 手动回滚（灰度失败时）

```bash
# 1. 立即将流量切换回旧版本
kubectl apply -f k8s/canary-0%.yaml

# 2. 验证旧版本服务正常
kubectl get pods -n coze-studio -l app=tenant-service,version=v1
kubectl logs -f -n coze-studio -l app=tenant-service,version=v1

# 3. 删除新版本
kubectl delete deployment tenant-service-v2 -n coze-studio

# 4. 分析失败原因
kubectl logs -n coze-studio -l app=tenant-service,version=v2 --previous
kubectl describe pod -n coze-studio -l app=tenant-service,version=v2
```

### 数据库回滚

```bash
# 使用数据库备份恢复
mysql -h mysql-master -u root -p zker < backup/zker_backup_20250101.sql

# 或使用Atlas迁移回滚
atlas migrate diff \
  --env local \
  --dir "file://backend/migrations" \
  --to "mysql://root:password@localhost:3306/zker" \
  --format '{{ sql . }}' | mysql -h mysql-master -u root -p zker
```

---

## 故障排查

### 常见问题及解决方案

#### 问题1: Pod启动失败

**症状**: `CrashLoopBackOff` 或 `ImagePullBackOff`

**排查**:
```bash
# 查看Pod详情
kubectl describe pod <pod-name> -n coze-studio

# 查看Pod日志
kubectl logs <pod-name> -n coze-studio

# 常见原因：
# - 镜像拉取失败 -> 检查镜像仓库访问和镜像名称
# - 资源不足 -> 检查节点资源（kubectl top nodes）
# - 配置错误 -> 检查ConfigMap和Secret
```

#### 问题2: 服务不可达

**症状**: `Connection refused` 或 `502 Bad Gateway`

**排查**:
```bash
# 检查Service
kubectl get svc -n coze-studio

# 检查Endpoint
kubectl get endpoints <service-name> -n coze-studio

# 检查Pod标签
kubectl get pods -n coze-studio --show-labels

# 常见原因：
# - Service选择器错误 -> 修复selector
# - Pod未就绪 -> 检查readinessProbe
# - 网络策略阻塞 -> 检查NetworkPolicy
```

#### 问题3: 数据库连接失败

**症状**: `Can't connect to MySQL server`

**排查**:
```bash
# 检查数据库Pod
kubectl get pods -n database

# 测试数据库连接
kubectl run -it --rm mysql-client --image=mysql:8.4.5 --restart=Never -- \
  mysql -h mysql-cluster-0.mysql-cluster.database.svc.cluster.local -u root -p

# 常见原因：
# - 数据库未就绪 -> 等待数据库启动
# - 密码错误 -> 检查Secret配置
# - 网络不通 -> 检查Service和NetworkPolicy
```

#### 问题4: 高延迟或超时

**症状**: API响应慢或超时

**排查**:
```bash
# 查看Prometheus指标
kubectl port-forward svc/prometheus-kube-prometheus-prometheus 9090:9090 -n monitoring
# 访问: http://localhost:9090
# 查询: http_request_duration_seconds_bucket{service="tenant-service"}

# 查看Jaeger追踪
kubectl port-forward svc/jaeger-prod-query 16686:16686 -n observability
# 访问: http://localhost:16686
# 搜索: tenant-service

# 常见原因：
# - 数据库慢查询 -> 查看慢查询日志，优化SQL或添加索引
# - Redis缓存未命中 -> 检查缓存策略
# - 内存不足 -> 检查Pod内存使用，增加资源限制
```

#### 问题5: 数据不一致

**症状**: 主从数据不同步

**排查**:
```bash
# 执行数据一致性检查
bash backend/scripts/data-consistency-check.sh

# 手动检查MySQL主从状态
kubectl exec -n database mysql-cluster-0 -- mysql -e "SHOW MASTER STATUS"
kubectl exec -n database mysql-cluster-1 -- mysql -e "SHOW SLAVE STATUS\G"

# 常见原因：
# - 复制延迟 -> 等待同步完成
# - 主从断连 -> 检查网络连接
# - 数据损坏 -> 使用备份恢复
```

### 应急响应流程

#### P0级别（完全中断）

1. **立即响应**（5分钟内）
   - 通知on-call工程师
   - 发送P0告警到所有渠道（Slack、邮件、短信）
   - 打开紧急战争室（Zoom会议）

2. **快速恢复**（15分钟内）
   - 执行自动回滚：`kubectl rollout undo`
   - 切换流量到备用数据中心（如果有）
   - 启用降级模式（只读模式）

3. **根本分析**（1小时内）
   - 收集日志和指标
   - 复现问题
   - 确定根本原因

4. **永久修复**（24小时内）
   - 实施修复
   - 更新测试用例
   - 进行事后分析（Postmortem）

#### P1级别（严重影响）

1. **响应**（15分钟内）
   - 通知团队负责人
   - 发送P1告警

2. **修复**（2小时内）
   - 评估影响范围
   - 实施热修复或配置更改

3. **验证**（4小时内）
   - 确认修复有效
   - 监控回归

---

## 联系方式

**On-Call工程师**: 138-xxxx-xxxx
**技术支持邮箱**: support@coze-studio.com
**Slack频道**: #coze-studio-oncall
**紧急战争室**: https://zoom.us/j/xxxxxxxxx
