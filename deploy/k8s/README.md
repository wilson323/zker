# Coze Studio - Kubernetes 部署清单

## 📁 目录结构

```
deploy/k8s/
├── base/                           # 基础设施配置
│   ├── namespace.yaml              # 命名空间定义
│   ├── configmap.yaml              # 通用配置 ConfigMap
│   ├── secret-template.yaml        # Secret 模板（敏感信息）
│   ├── network-policy.yaml         # 网络策略（安全隔离）
│   └── generate-secret.sh          # Secret 生成辅助脚本
├── backend/                        # 后端服务
│   ├── api-server.yaml             # API Server Deployment + Service + HPA + PDB
│   └── worker.yaml                 # Worker Deployment + Service + PDB
├── frontend/                       # 前端服务
│   └── nginx.yaml                  # Nginx Deployment + Service + Ingress + HPA
├── middleware/                     # 中间件
│   ├── mysql.yaml                  # MySQL StatefulSet + Service + ConfigMap
│   ├── redis.yaml                  # Redis Master/Slave + Sentinel (哨兵高可用)
│   ├── elasticsearch.yaml          # Elasticsearch StatefulSet (3 节点集群)
│   ├── minio.yaml                  # MinIO Deployment + Service + PVC
│   └── nsq-etcd.yaml               # NSQ (消息队列) + etcd (配置中心)
├── KUBERNETES_DEPLOYMENT_GUIDE.md  # 完整部署指南
└── README.md                       # 本文件
```

---

## 🚀 快速开始

### 前置要求

- Kubernetes 集群 ≥ 1.25
- kubectl 已配置
- 镜像已推送到镜像仓库
- StorageClass 已配置

### 一键部署（推荐）

```bash
# 1. 生成 Secret
cd deploy/k8s
chmod +x generate-secret.sh
./generate-secret.sh
kubectl apply -f secret.yaml -n coze-studio
rm -f secret.yaml

# 2. 应用所有配置
kubectl apply -f base/namespace.yaml
kubectl apply -f base/configmap.yaml
kubectl apply -f middleware/
kubectl apply -f backend/
kubectl apply -f frontend/
kubectl apply -f base/network-policy.yaml

# 3. 验证部署
kubectl get pods -n coze-studio
```

**详细步骤请参考：** [KUBERNETES_DEPLOYMENT_GUIDE.md](./KUBERNETES_DEPLOYMENT_GUIDE.md)

---

## 📦 已部署组件清单

### 应用层

| 组件 | 类型 | 副本数 | 端口 | 说明 |
|------|------|--------|------|------|
| **API Server** | Deployment | 3 (HPA) | 8888 (HTTP), 9000 (gRPC) | 后端 API 服务，支持自动扩展 |
| **Worker** | Deployment | 2 | 9091 (metrics) | 后台任务处理器 |
| **Frontend** | Deployment | 2 (HPA) | 80 (HTTP) | 前端 Nginx 服务 |

### 中间件层

| 组件 | 类型 | 副本数 | 端口 | 说明 |
|------|------|--------|------|------|
| **MySQL** | StatefulSet | 1 | 3306 | 主数据库 |
| **Redis Master** | StatefulSet | 1 | 6379 | Redis 主节点 |
| **Redis Slave** | StatefulSet | 2 | 6379 | Redis 从节点 |
| **Redis Sentinel** | Deployment | 3 | 26379 | Redis 哨兵（高可用） |
| **Elasticsearch** | StatefulSet | 3 | 9200 (HTTP), 9300 (Transport) | 搜索引擎集群 |
| **MinIO** | Deployment | 1 | 9000 (API), 9001 (Console) | 对象存储 |
| **NSQD** | Deployment | 1 | 4150 (TCP), 4151 (HTTP) | 消息队列守护进程 |
| **NSQ Lookupd** | Deployment | 1 | 4160 (TCP), 4161 (HTTP) | 服务发现 |
| **NSQ Admin** | Deployment | 1 | 4171 (HTTP) | Web 管理界面 |
| **etcd** | StatefulSet | 3 | 2379 (Client), 2380 (Peer) | 配置中心（3 节点集群） |

### 高可用组件

- **HPA (Horizontal Pod Autoscaler)**: API Server、Frontend
- **PDB (PodDisruptionBudget)**: 所有核心组件
- **NetworkPolicy**: 网络隔离（最小权限原则）

---

## 🔧 配置说明

### 环境变量

所有配置通过 `ConfigMap` 和 `Secret` 管理：

- **ConfigMap** (`coze-studio-config`): 非敏感配置
  - 数据库连接地址
  - Redis/ES/MinIO 端点
  - 性能参数（连接池大小、超时等）

- **Secret** (`coze-studio-secret`): 敏感配置
  - 数据库密码
  - API 密钥
  - JWT 签名密钥

### 存储需求

| 组件 | 存储大小 | StorageClass |
|------|----------|--------------|
| MySQL | 500Gi | fast-ssd |
| Redis Master | 50Gi | fast-ssd |
| Redis Slave | 50Gi × 2 | fast-ssd |
| Elasticsearch | 200Gi × 3 | fast-ssd |
| MinIO | 1Ti | fast-ssd |
| NSQD | 100Gi | fast-ssd |
| etcd | 50Gi × 3 | fast-ssd |
| **总计** | **~3.5Ti** | - |

### 资源配置（每个 Pod）

| 组件 | CPU (Request/Limit) | Memory (Request/Limit) |
|------|---------------------|------------------------|
| API Server | 500m / 2000m | 512Mi / 2Gi |
| Worker | 200m / 1000m | 256Mi / 1Gi |
| Frontend | 100m / 500m | 128Mi / 512Mi |
| MySQL | 1000m / 4000m | 2Gi / 8Gi |
| Redis | 500m / 2000m | 1Gi / 4Gi |
| Elasticsearch | 1000m / 4000m | 4Gi / 8Gi |
| MinIO | 500m / 2000m | 1Gi / 4Gi |

---

## 📊 监控与日志

### Prometheus 监控

所有 Pod 都配置了 Prometheus 注解：

```yaml
prometheus.io/scrape: "true"
prometheus.io/port: "8888"
prometheus.io/path: "/metrics"
```

### 健康检查

- **Liveness Probe**: `/api/health/live`
- **Readiness Probe**: `/api/health/ready`
- **Startup Probe**: `/api/health/live`

### 日志查看

```bash
# API Server
kubectl logs -l component=api-server -n coze-studio -f

# Worker
kubectl logs -l component=worker -n coze-studio -f

# MySQL
kubectl logs mysql-0 -n coze-studio -f
```

---

## 🔒 安全特性

1. **网络隔离**: NetworkPolicy 限制 Pod 间通信
2. **密钥管理**: Secret 存储敏感信息，支持 External Secrets Operator
3. **RBAC**: 基于角色的访问控制
4. **Pod Security**: Pod Security Admission 策略
5. **TLS 支持**: Ingress TLS 配置（使用 cert-manager）

---

## 🛠️ 常用操作

### 扩展服务

```bash
# 手动扩展
kubectl scale deployment coze-studio-api -n coze-studio --replicas=5

# HPA 自动扩展（已配置）
kubectl get hpa -n coze-studio
```

### 更新镜像

```bash
kubectl set image deployment/coze-studio-api api-server=coze-studio/backend-api:v1.2.3 -n coze-studio
kubectl rollout status deployment/coze-studio-api -n coze-studio
```

### 回滚

```bash
kubectl rollout undo deployment/coze-studio-api -n coze-studio
```

### 查看日志

```bash
kubectl logs -f <pod-name> -n coze-studio
kubectl logs -f <pod-name> -n coze-studio --previous  # 上一个版本的日志
```

### 进入 Pod

```bash
kubectl exec -it <pod-name> -n coze-studio -- sh
```

---

## 📖 相关文档

- **[KUBERNETES_DEPLOYMENT_GUIDE.md](./KUBERNETES_DEPLOYMENT_GUIDE.md)** - 完整部署指南
- **[生产环境部署指南](../docs/deployment/PRODUCTION_DEPLOYMENT_GUIDE.md)** - 生产环境最佳实践
- **[灰度发布策略](../../docs/企业级功能完善与统一性设计方案/ZKER-灰度发布策略_v1.0.md)** - 灰度发布方案
- **[监控运维指南](../docs/deployment/MONITORING_OPS_GUIDE.md)** - Prometheus + Grafana 使用指南

---

## 🆘 故障排查

### Pod 无法启动

```bash
kubectl describe pod <pod-name> -n coze-studio
kubectl logs <pod-name> -n coze-studio
```

### 服务无法访问

```bash
kubectl get endpoints <service-name> -n coze-studio
kubectl run -it --rm debug --image=busybox -- sh -c "wget -O- http://<service-name>:8888/api/health"
```

### PVC 无法绑定

```bash
kubectl describe pvc <pvc-name> -n coze-studio
kubectl get storageclass
```

**详细故障排查请参考：** [KUBERNETES_DEPLOYMENT_GUIDE.md](./KUBERNETES_DEPLOYMENT_GUIDE.md) 第 7 节

---

## 📝 更新日志

### v1.0 (2025-01-01)

- ✅ 创建完整的 Kubernetes 部署清单
- ✅ 支持所有中间件（MySQL、Redis、ES、MinIO、NSQ、etcd）
- ✅ 配置 HPA 自动扩展
- ✅ 配置 PDB 保障高可用
- ✅ 配置 NetworkPolicy 网络隔离
- ✅ 完整的部署文档和操作指南

---

## 👥 维护者

ZKER DevOps Team

---

**文档版本：** v1.0
**最后更新：** 2025-01-01
