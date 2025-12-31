# Coze Studio - Kubernetes 部署完整指南

## 📋 目录

1. [环境准备](#环境准备)
2. [集群配置](#集群配置)
3. [部署步骤](#部署步骤)
4. [验证部署](#验证部署)
5. [扩展与高可用](#扩展与高可用)
6. [监控与日志](#监控与日志)
7. [故障排查](#故障排查)
8. [维护操作](#维护操作)
9. [安全最佳实践](#安全最佳实践)

---

## 🚀 环境准备

### 1.1 Kubernetes 集群要求

**最低配置（开发/测试环境）：**
- Kubernetes 版本：≥ 1.25
- 节点数：1 个 Master + 2 个 Worker
- 每个 Worker 节点资源：4 CPU, 16GB RAM, 100GB SSD

**推荐配置（生产环境）：**
- Kubernetes 版本：≥ 1.27
- 节点数：3 个 Master（HA）+ 5 个 Worker
- 每个 Worker 节点资源：8 CPU, 32GB RAM, 500GB SSD

**必需的 Kubernetes 组件：**
```bash
# 检查集群状态
kubectl version --short
kubectl cluster-info
kubectl get nodes

# 检查可用资源
kubectl top nodes
```

### 1.2 安装必需的组件

#### 1.2.1 安装 Ingress Controller（Nginx）

```bash
# 安装 NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml

# 验证安装
kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
```

#### 1.2.2 安装 cert-manager（自动 TLS 证书）

```bash
# 安装 cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.2/cert-manager.yaml

# 验证安装
kubectl get pods -n cert-manager
```

#### 1.2.3 安装 Prometheus Operator（可选，用于监控）

```bash
# 使用 Helm 安装
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

kubectl create namespace monitoring
helm install kube-prometheus-stack prometheus-community/kube-prometheus-stack -n monitoring
```

#### 1.2.4 安装 StorageClass（动态存储）

```bash
# 查看可用的 StorageClass
kubectl get storageclass

# 如果没有，创建一个（示例：使用 local-path）
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.26/deploy/local-path-storage.yaml
```

### 1.3 准备容器镜像

**构建并推送镜像到镜像仓库：**

```bash
# 后端 API 镜像
cd backend
docker build -t your-registry.com/coze-studio/backend-api:latest .
docker push your-registry.com/coze-studio/backend-api:latest

# 前端镜像
cd frontend/apps/coze-studio
docker build -t your-registry.com/coze-studio/frontend:latest .
docker push your-registry.com/coze-studio/frontend:latest

# Worker 镜像（如果需要单独构建）
docker build -t your-registry.com/coze-studio/backend-worker:latest -f Dockerfile.worker .
docker push your-registry.com/coze-studio/backend-worker:latest
```

**如果使用私有镜像仓库，创建 ImagePullSecret：**

```bash
kubectl create secret docker-registry docker-registry-secret \
  --docker-server=your-registry.com \
  --docker-username=your-username \
  --docker-password=your-password \
  --docker-email=your-email \
  -n coze-studio
```

---

## ⚙️ 集群配置

### 2.1 创建命名空间

```bash
kubectl apply -f deploy/k8s/base/namespace.yaml

# 验证
kubectl get namespace coze-studio
```

### 2.2 创建 ConfigMap

```bash
kubectl apply -f deploy/k8s/base/configmap.yaml

# 验证
kubectl get configmap -n coze-studio
kubectl describe configmap coze-studio-config -n coze-studio
```

### 2.3 生成 Secret

**方式一：使用交互式脚本生成（推荐）**

```bash
cd deploy/k8s
chmod +x generate-secret.sh
./generate-secret.sh

# 应用生成的 Secret
kubectl apply -f secret.yaml -n coze-studio

# 删除本地 Secret 文件（安全起见）
rm -f secret.yaml
```

**方式二：手动创建**

```bash
# 编辑 Secret 模板
cp deploy/k8s/base/secret-template.yaml secret.yaml

# 使用 base64 编码替换所有占位符
echo -n "your_password" | base64

# 应用 Secret
kubectl apply -f secret.yaml -n coze-studio

# 删除临时文件
rm -f secret.yaml
```

**验证 Secret：**

```bash
kubectl get secret coze-studio-secret -n coze-studio
kubectl describe secret coze-studio-secret -n coze-studio
```

---

## 🚢 部署步骤

### 3.1 部署顺序（按依赖关系）

**重要提示：** 必须按照以下顺序部署，否则会导致依赖服务无法启动。

```bash
# 1. 基础设施（已完成）
# Namespace、ConfigMap、Secret

# 2. 中间件（数据库、缓存、消息队列）
kubectl apply -f deploy/k8s/middleware/mysql.yaml
kubectl apply -f deploy/k8s/middleware/redis.yaml
kubectl apply -f deploy/k8s/middleware/elasticsearch.yaml
kubectl apply -f deploy/k8s/middleware/minio.yaml
kubectl apply -f deploy/k8s/middleware/nsq-etcd.yaml

# 3. 后端服务
kubectl apply -f deploy/k8s/backend/api-server.yaml
kubectl apply -f deploy/k8s/backend/worker.yaml

# 4. 前端服务
kubectl apply -f deploy/k8s/frontend/nginx.yaml

# 5. 网络策略（最后应用，避免部署期间阻塞流量）
kubectl apply -f deploy/k8s/base/network-policy.yaml
```

### 3.2 分步骤部署详解

#### 步骤 1：部署 MySQL

```bash
kubectl apply -f deploy/k8s/middleware/mysql.yaml

# 等待 MySQL 就绪
kubectl wait --for=condition=ready pod -l component=mysql -n coze-studio --timeout=300s

# 验证
kubectl get pods -l component=mysql -n coze-studio
kubectl get pvc -n coze-studio | grep mysql
```

#### 步骤 2：部署 Redis

```bash
kubectl apply -f deploy/k8s/middleware/redis.yaml

# 等待 Redis Master 就绪
kubectl wait --for=condition=ready pod -l component=redis,role=master -n coze-studio --timeout=300s

# 验证
kubectl get pods -l component=redis -n coze-studio
kubectl get pvc -n coze-studio | grep redis
```

#### 步骤 3：部署 Elasticsearch

```bash
kubectl apply -f deploy/k8s/middleware/elasticsearch.yaml

# 等待 ES 集群就绪（可能需要几分钟）
kubectl wait --for=condition=ready pod -l component=elasticsearch -n coze-studio --timeout=600s

# 验证集群健康状态
kubectl exec -it elasticsearch-0 -n coze-studio -- curl -s http://localhost:9200/_cluster/health
```

#### 步骤 4：部署其他中间件（MinIO、NSQ、etcd）

```bash
kubectl apply -f deploy/k8s/middleware/minio.yaml
kubectl apply -f deploy/k8s/middleware/nsq-etcd.yaml

# 等待所有中间件就绪
kubectl wait --for=condition=ready pod -l 'component in (minio,nsqd,nsqlookupd,nsqadmin,etcd)' -n coze-studio --timeout=300s

# 验证
kubectl get pods -n coze-studio
```

#### 步骤 5：部署后端 API Server

```bash
kubectl apply -f deploy/k8s/backend/api-server.yaml

# 等待 API Server 就绪
kubectl wait --for=condition=ready pod -l component=api-server -n coze-studio --timeout=300s

# 验证
kubectl get pods -l component=api-server -n coze-studio
kubectl get svc -l component=api-server -n coze-studio

# 测试 API 端点
kubectl port-forward svc/coze-studio-api 8888:8888 -n coze-studio
curl http://localhost:8888/api/health
```

#### 步骤 6：部署 Worker

```bash
kubectl apply -f deploy/k8s/backend/worker.yaml

# 验证
kubectl get pods -l component=worker -n coze-studio
kubectl logs -l component=worker -n coze-studio --tail=50
```

#### 步骤 7：部署前端

```bash
kubectl apply -f deploy/k8s/frontend/nginx.yaml

# 等待前端就绪
kubectl wait --for=condition=ready pod -l component=frontend -n coze-studio --timeout=300s

# 验证
kubectl get pods -l component=frontend -n coze-studio
kubectl get svc -l component=frontend -n coze-studio
```

#### 步骤 8：应用网络策略

```bash
kubectl apply -f deploy/k8s/base/network-policy.yaml

# 验证
kubectl get networkpolicy -n coze-studio
kubectl describe networkpolicy api-server-policy -n coze-studio
```

---

## ✅ 验证部署

### 4.1 检查所有 Pod 状态

```bash
# 查看所有 Pod
kubectl get pods -n coze-studio -o wide

# 预期输出（示例）：
# NAME                                READY   STATUS    RESTARTS   AGE
# coze-studio-api-xxx                 2/2     Running   0          5m
# coze-studio-worker-xxx              1/1     Running   0          4m
# coze-studio-frontend-xxx            1/1     Running   0          3m
# mysql-0                             1/1     Running   0          10m
# redis-master-0                      1/1     Running   0          9m
# redis-slave-0                       1/1     Running   0          9m
# redis-slave-1                       1/1     Running   0          9m
# redis-sentinel-xxx                  1/1     Running   0          9m
# elasticsearch-0                     1/1     Running   0          8m
# elasticsearch-1                     1/1     Running   0          8m
# elasticsearch-2                     1/1     Running   0          8m
# minio-xxx                           1/1     Running   0          7m
# nsqd-xxx                            1/1     Running   0          7m
# nsqlookupd-xxx                      1/1     Running   0          7m
# nsqadmin-xxx                        1/1     Running   0          7m
# etcd-0                              1/1     Running   0          7m
# etcd-1                              1/1     Running   0          7m
# etcd-2                              1/1     Running   0          7m

# 检查是否有 Pod 处于非 Running 状态
kubectl get pods -n coze-studio | grep -v Running | grep -v Completed
```

### 4.2 检查服务端点

```bash
# 查看所有 Service
kubectl get svc -n coze-studio

# 测试 API 健康检查
kubectl port-forward svc/coze-studio-api 8888:8888 -n coze-studio &
curl http://localhost:8888/api/health
curl http://localhost:8888/api/health/live
curl http://localhost:8888/api/health/ready
curl http://localhost:8888/api/health/detailed

# 测试前端
kubectl port-forward svc/coze-studio-frontend 8080:80 -n coze-studio &
curl http://localhost:8080
```

### 4.3 检查持久化存储

```bash
# 查看 PVC 绑定状态
kubectl get pvc -n coze-studio

# 预期输出：
# NAME              STATUS   VOLUME                                     CAPACITY   ACCESS MODES
# data-elasticsearch-0   Bound    pvc-xxx                                     200Gi      RWO
# data-mysql-0          Bound    pvc-xxx                                     500Gi      RWO
# data-redis-master-0   Bound    pvc-xxx                                     50Gi       RWO
# data-redis-slave-0    Bound    pvc-xxx                                     50Gi       RWO
# data-redis-slave-1    Bound    pvc-xxx                                     50Gi       RWO
# minio-pvc             Bound    pvc-xxx                                     1Ti        RWO
# nsqd-pvc              Bound    pvc-xxx                                     100Gi      RWO
```

### 4.4 检查 HPA 和 PDB

```bash
# 查看 HPA 状态
kubectl get hpa -n coze-studio

# 查看 PDB 状态
kubectl get pdb -n coze-studio
kubectl describe pdb coze-studio-api-pdb -n coze-studio
```

### 4.5 检查日志

```bash
# API Server 日志
kubectl logs -l component=api-server -n coze-studio --tail=100 -f

# Worker 日志
kubectl logs -l component=worker -n coze-studio --tail=100 -f

# MySQL 日志
kubectl logs mysql-0 -n coze-studio --tail=50

# Redis 日志
kubectl logs redis-master-0 -n coze-studio --tail=50

# 事件日志
kubectl get events -n coze-studio --sort-by='.lastTimestamp'
```

---

## 📊 扩展与高可用

### 5.1 手动扩展

```bash
# 扩展 API Server
kubectl scale deployment coze-studio-api -n coze-studio --replicas=5

# 扩展 Worker
kubectl scale deployment coze-studio-worker -n coze-studio --replicas=5

# 扩展 Frontend
kubectl scale deployment coze-studio-frontend -n coze-studio --replicas=3

# 验证
kubectl get pods -n coze-studio -l component=api-server
```

### 5.2 自动扩展（HPA）

HPA 已在 manifest 中配置，无需手动操作。

```bash
# 查看 HPA 状态
kubectl get hpa -n coze-studio -w

# 查看 HPA 详细信息
kubectl describe hpa coze-studio-api-hpa -n coze-studio

# 手动触发扩展（压力测试）
kubectl run -it --rm stress-test --image=busybox --restart=Never -- sh -c "while true; do wget -q -O- http://coze-studio-api/api/health; done"
```

### 5.3 节点维护（PDB）

**安全地维护节点（确保最小可用副本）：**

```bash
# 查看 PDB 状态
kubectl get pdb -n coze-studio

# 安全驱逐 Pod（遵循 PDB 约束）
kubectl cordon <node-name>  # 标记节点为不可调度
kubectl drain <node-name> --ignore-daemonsets --delete-emptydir-data

# 维护完成后恢复节点
kubectl uncordon <node-name>
```

---

## 📈 监控与日志

### 6.1 Prometheus 监控

**确保 Prometheus 采集指标：**

```bash
# 检查 ServiceMonitor（如果使用 Prometheus Operator）
kubectl get servicemonitor -n coze-studio

# 手动验证指标端点
kubectl port-forward svc/coze-studio-api 8888:8888 -n coze-studio &
curl http://localhost:8888/metrics
```

**访问 Prometheus UI：**

```bash
# 端口转发
kubectl port-forward svc/kube-prometheus-stack-prometheus 9090:9090 -n monitoring

# 浏览器访问：http://localhost:9090
```

### 6.2 Grafana 大盘

**导入预配置的 Dashboard：**

```bash
# 端口转发
kubectl port-forward svc/kube-prometheus-stack-grafana 3000:80 -n monitoring

# 浏览器访问：http://localhost:3000
# 默认用户名/密码：admin / prom-operator

# 导入 Dashboard（JSON 文件在 deploy/monitoring/dashboards/）
```

### 6.3 日志收集（ELK Stack）

**如果部署了 ELK：**

```bash
# 查看 Filebeat Pod
kubectl get pods -n kube-system | grep filebeat

# 查看 Kibana
kubectl port-forward svc/kibana 5601:5601 -n monitoring
```

---

## 🔧 故障排查

### 7.1 Pod 启动失败

**问题：** Pod 处于 `CrashLoopBackOff` 或 `Error` 状态

```bash
# 查看 Pod 状态
kubectl describe pod <pod-name> -n coze-studio

# 查看日志
kubectl logs <pod-name> -n coze-studio --previous

# 常见原因：
# 1. 镜像拉取失败 → 检查 ImagePullSecret
# 2. ConfigMap/Secret 缺失 → 检查配置是否挂载
# 3. 资源不足 → 检查 Node 资源
# 4. 健康检查失败 → 检查应用日志
```

### 7.2 服务无法访问

**问题：** 无法访问 Service 或 Pod

```bash
# 检查 Service 端点
kubectl get endpoints <service-name> -n coze-studio

# 测试 Pod 到 Pod 连通性
kubectl run -it --rm debug --image=busybox --restart=Never -- sh
# 在 Pod 内执行：
# wget -O- http://coze-studio-api:8888/api/health
# wget -O- http://10.x.x.x:8888/api/health  # 使用 Pod IP

# 检查 NetworkPolicy
kubectl get networkpolicy -n coze-studio
kubectl describe networkpolicy api-server-policy -n coze-studio
```

### 7.3 数据库连接失败

**问题：** 应用无法连接到 MySQL/Redis

```bash
# 检查数据库 Pod 状态
kubectl get pods -l component=mysql -n coze-studio
kubectl exec -it mysql-0 -n coze-studio -- mysql -uroot -p$(kubectl get secret coze-studio-secret -n coze-studio -o jsonpath='{.data.db-password}' | base64 -d)

# 检查密码
kubectl get secret coze-studio-secret -n coze-studio -o jsonpath='{.data.db-password}' | base64 -d

# 检查网络策略
kubectl describe networkpolicy mysql-policy -n coze-studio
```

### 7.4 PVC 无法绑定

**问题：** PVC 处于 `Pending` 状态

```bash
# 查看 PVC 状态
kubectl describe pvc <pvc-name> -n coze-studio

# 常见原因：
# 1. StorageClass 不存在 → kubectl get storageclass
# 2. 没有可用的 PV → 检查存储供应
# 3. 访问模式不匹配 → 检查 PVC accessModes
```

---

## 🛠️ 维护操作

### 8.1 更新镜像（滚动更新）

```bash
# 更新镜像标签
kubectl set image deployment/coze-studio-api api-server=coze-studio/backend-api:v1.2.3 -n coze-studio

# 查看滚动更新状态
kubectl rollout status deployment/coze-studio-api -n coze-studio

# 查看更新历史
kubectl rollout history deployment/coze-studio-api -n coze-studio

# 回滚到上一个版本
kubectl rollout undo deployment/coze-studio-api -n coze-studio

# 回滚到指定版本
kubectl rollout undo deployment/coze-studio-api --to-revision=2 -n coze-studio
```

### 8.2 备份与恢复

**备份：**

```bash
# 方法 1：使用 Velero（推荐）
velero backup create coze-studio-backup-$(date +%Y%m%d) --include-namespaces coze-studio

# 方法 2：手动备份 MySQL
kubectl exec mysql-0 -n coze-studio -- mysqldump -uroot -p${MYSQL_ROOT_PASSWORD} --all-databases > backup.sql

# 方法 3：备份 PVC 快照（云厂商支持）
kubectl get pvc -n coze-studio
# 使用云厂商 CLI 创建快照
```

**恢复：**

```bash
# Velero 恢复
velero restore create --from-backup coze-studio-backup-20250101

# MySQL 恢复
cat backup.sql | kubectl exec -i mysql-0 -n coze-studio -- mysql -uroot -p${MYSQL_ROOT_PASSWORD}
```

### 8.3 配置更新

```bash
# 更新 ConfigMap
kubectl edit configmap coze-studio-config -n coze-studio

# 重启相关 Pod 以使配置生效
kubectl rollout restart deployment/coze-studio-api -n coze-studio
kubectl rollout restart deployment/coze-studio-worker -n coze-studio

# 更新 Secret
kubectl create secret generic coze-studio-secret --dry-run=client --from-env-file=new.env -o yaml | kubectl apply -f -
kubectl rollout restart deployment/coze-studio-api -n coze-studio
```

### 8.4 清理资源

```bash
# 删除所有资源（谨慎操作！）
kubectl delete namespace coze-studio

# 删除特定资源
kubectl delete deployment -l app=coze-studio -n coze-studio
kubectl delete statefulset -l app=coze-studio -n coze-studio
kubectl delete pvc -n coze-studio --all
```

---

## 🔒 安全最佳实践

### 9.1 Secret 管理

**使用 External Secrets Operator：**

```bash
# 安装 External Secrets Operator
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace

# 配置 Vault/AWS Secrets Manager 集成
# 参考 deploy/k8s/base/secret-template.yaml 中的示例
```

### 9.2 RBAC 配置

**创建服务账号和角色：**

```bash
# 创建 ServiceAccount
kubectl create serviceaccount coze-studio-sa -n coze-studio

# 创建 Role（读权限）
kubectl create role coze-studio-reader --verb=get,list,watch --resource=pods,services,configmaps -n coze-studio

# 绑定 Role
kubectl create rolebinding coze-studio-reader-binding --role=coze-studio-reader --serviceaccount=coze-studio:coze-studio-sa -n coze-studio
```

### 9.3 Pod Security Standards

**使用 Pod Security Admission：**

```bash
# 标记 Namespace 为 privileged（仅用于需要特权的服务）
kubectl label --overwrite ns coze-studio pod-security.kubernetes.io/enforce=privileged

# 或使用 baseline（推荐）
kubectl label --overwrite ns coze-studio pod-security.kubernetes.io/enforce=baseline
```

### 9.4 网络隔离

**确保 NetworkPolicy 已启用：**

```bash
# 验证
kubectl get networkpolicy -n coze-studio

# 测试网络隔离
kubectl run -it --rm test --image=busybox --restart=Never -- sh
# 在 Pod 内尝试访问不允许的服务
```

---

## 📚 附录

### A. 快速部署脚本

```bash
#!/bin/bash
set -e

echo "🚀 开始部署 Coze Studio 到 Kubernetes..."

# 1. 创建命名空间
kubectl apply -f deploy/k8s/base/namespace.yaml

# 2. 生成 Secret
cd deploy/k8s
./generate-secret.sh
kubectl apply -f secret.yaml -n coze-studio
rm -f secret.yaml
cd ../..

# 3. 应用 ConfigMap
kubectl apply -f deploy/k8s/base/configmap.yaml

# 4. 部署中间件
kubectl apply -f deploy/k8s/middleware/mysql.yaml
kubectl apply -f deploy/k8s/middleware/redis.yaml
kubectl apply -f deploy/k8s/middleware/elasticsearch.yaml
kubectl apply -f deploy/k8s/middleware/minio.yaml
kubectl apply -f deploy/k8s/middleware/nsq-etcd.yaml

# 等待中间件就绪
echo "⏳ 等待中间件就绪..."
kubectl wait --for=condition=ready pod -l 'component in (mysql,redis,elasticsearch,minio,nsqd,etcd)' -n coze-studio --timeout=600s

# 5. 部署后端
kubectl apply -f deploy/k8s/backend/api-server.yaml
kubectl apply -f deploy/k8s/backend/worker.yaml

# 6. 部署前端
kubectl apply -f deploy/k8s/frontend/nginx.yaml

# 7. 应用网络策略
kubectl apply -f deploy/k8s/base/network-policy.yaml

echo "✅ 部署完成！"
echo ""
echo "验证部署："
kubectl get pods -n coze-studio
```

### B. 常用命令速查

```bash
# 查看 Pod 日志
kubectl logs -f <pod-name> -n coze-studio

# 进入 Pod
kubectl exec -it <pod-name> -n coze-studio -- sh

# 端口转发
kubectl port-forward svc/<service-name> 8080:80 -n coze-studio

# 查看资源使用
kubectl top pods -n coze-studio
kubectl top nodes

# 查看事件
kubectl get events -n coze-studio --sort-by='.lastTimestamp'

# 导出配置
kubectl get configmap coze-studio-config -n coze-studio -o yaml > backup-config.yaml
```

---

## 🆘 支持

如遇问题，请：

1. 查看日志：`kubectl logs -f <pod-name> -n coze-studio`
2. 查看事件：`kubectl get events -n coze-studio`
3. 参考文档：`docs/企业级功能完善与统一性设计方案/`
4. 联系 DevOps 团队：devops@coze-studio.com

---

**文档版本：** v1.0
**最后更新：** 2025-01-01
**维护者：** ZKER DevOps Team
