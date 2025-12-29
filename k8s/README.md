# ZKER Kubernetes 部署指南

> **最后更新**: 2025-01-01
> **维护人**: 研发 D - DevOps 工程师

---

## 📖 概述

本目录包含 ZKER 项目的 Kubernetes 部署配置,使用 Kustomize 管理多环境配置。

**目录结构**:
```
k8s/
├── base/              # 基础配置 (所有环境共享)
│   ├── namespace.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── hpa.yaml
│   └── kustomization.yaml
├── overlays/
│   ├── dev/           # 开发环境
│   ├── staging/       # 预发布环境
│   └── prod/          # 生产环境
└── README.md          # 本文档
```

---

## 🚀 快速开始

### 前置要求

- Kubernetes 集群 (v1.24+)
- kubectl 配置完成
- kustomize 安装 (`kubectl apply -k` 内置)

### 1. 创建命名空间

```bash
# 创建所有命名空间
kubectl apply -k k8s/base/
```

### 2. 配置密钥

**重要**: 生产环境必须修改 `secret.yaml` 中的默认值!

```bash
# 方式一: 直接编辑 Secret
kubectl edit secret db-secret -n zker
kubectl edit secret api-keys -n zker

# 方式二: 从文件创建
kubectl create secret generic db-secret \
  --from-literal=password=your_secure_password \
  --from-literal=user=zker \
  -n zker

# 方式三: 从环境变量创建
kubectl create secret generic api-keys \
  --from-literal=openai-api-key=$OPENAI_API_KEY \
  --from-literal=anthropic-api-key=$ANTHROPIC_API_KEY \
  -n zker
```

### 3. 部署到指定环境

```bash
# 开发环境
kubectl apply -k k8s/overlays/dev/

# 预发布环境
kubectl apply -k k8s/overlays/staging/

# 生产环境
REGISTRY=your-registry.com kubectl apply -k k8s/overlays/prod/
```

### 4. 验证部署

```bash
# 查看 Pod 状态
kubectl get pods -n zker-prod

# 查看 Service
kubectl get svc -n zker-prod

# 查看 HPA
kubectl get hpa -n zker-prod

# 查看日志
kubectl logs -f deployment/zker-api -n zker-prod

# 端口转发 (本地测试)
kubectl port-forward svc/zker-api 8080:80 -n zker-prod
```

---

## 🔧 环境配置差异

| 配置项 | 开发环境 | 预发布环境 | 生产环境 |
|--------|---------|-----------|---------|
| 副本数 | 1 | 2 | 3 (HPA: 3-20) |
| 日志级别 | debug | info | info |
| 日志格式 | text | json | json |
| 资源限制 | 128Mi/256Mi | 256Mi/512Mi | 256Mi/512Mi |
| HPA 最大副本 | 5 | 10 | 20 |
| Jaeger 采样率 | 100% | 50% | 10% |
| 速率限制 | 1000 req/s | 200 req/s | 100 req/s |

---

## 📊 监控和日志

### Prometheus 指标

API 服务暴露 `/metrics` 端点 (端口 8889):

```bash
# 访问指标
kubectl port-forward svc/zker-api 8889:8889 -n zker-prod
curl http://localhost:8889/metrics
```

### Grafana Dashboard

1. 访问 Grafana: `http://grafana.example.com`
2. 导入 Dashboard ID: `15757` (Kubernetes Cluster Monitoring)
3. 配置 Prometheus 数据源

### 日志查询

```bash
# 查看所有 Pod 日志
kubectl logs -l app=zker-api -n zker-prod --all-containers=true

# 查看特定 Pod 日志
kubectl logs -f pod/zker-api-xxx -n zker-prod

# 查看最近的日志
kubectl logs --tail=100 deployment/zker-api -n zker-prod

# 查看错误日志
kubectl logs deployment/zker-api -n zker-prod | grep -i error
```

---

## 🔄 更新和回滚

### 更新部署

```bash
# 更新镜像
kubectl set image deployment/zker-api \
  zker-api=zker/api:v1.0.1 \
  -n zker-prod

# 或者使用 Kustomize
kubectl apply -k k8s/overlays/prod/
```

### 查看更新状态

```bash
# 查看 Rollout 状态
kubectl rollout status deployment/zker-api -n zker-prod

# 查看更新历史
kubectl rollout history deployment/zker-api -n zker-prod
```

### 回滚部署

```bash
# 回滚到上一个版本
kubectl rollout undo deployment/zker-api -n zker-prod

# 回滚到指定版本
kubectl rollout undo deployment/zker-api --to-revision=3 -n zker-prod

# 使用回滚脚本 (见 scripts/rollback.sh)
./scripts/rollback.sh zker-prod v1.0.0
```

---

## 🔒 安全最佳实践

### 1. 最小权限原则

- 使用 `ServiceAccount` 限制 Pod 权限
- 启用 `PodSecurityPolicy` 或 `Pod Security Standards`
- 避免使用 `privileged` 模式

### 2. 网络策略

```yaml
# 示例: 只允许特定 Service 访问
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: zker-api-policy
spec:
  podSelector:
    matchLabels:
      app: zker-api
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: zker-frontend
    ports:
    - protocol: TCP
      port: 8888
```

### 3. 密钥管理

- **生产环境必须** 使用外部密钥管理系统:
  - AWS Secrets Manager
  - Azure Key Vault
  - HashiCorp Vault
- 使用 `External Secrets Operator` 同步密钥到 K8s

### 4. 镜像安全

```bash
# 扫描镜像漏洞
trivy image zker/api:v1.0.0

# 签名镜像
cosign sign zker/api:v1.0.0

# 验证签名
cosign verify zker/api:v1.0.0
```

---

## 📈 性能优化

### HPA 调优

```bash
# 查看 HPA 状态
kubectl get hpa -n zker-prod -w

# 查看 HPA 详细指标
kubectl describe hpa zker-api-hpa -n zker-prod

# 调整 HPA 参数
kubectl edit hpa zker-api-hpa -n zker-prod
```

### 资源配额

```yaml
# 示例: 限制命名空间资源
apiVersion: v1
kind: ResourceQuota
metadata:
  name: zker-quota
  namespace: zker-prod
spec:
  hard:
    requests.cpu: "10"
    requests.memory: "20Gi"
    limits.cpu: "20"
    limits.memory: "40Gi"
    persistentvolumeclaims: "5"
```

### Pod 反亲和性

```yaml
# 确保 Pod 分散在不同节点
spec:
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchExpressions:
            - key: app
              operator: In
              values:
              - zker-api
          topologyKey: kubernetes.io/hostname
```

---

## 🐛 故障排查

### Pod 无法启动

```bash
# 查看 Pod 事件
kubectl describe pod zker-api-xxx -n zker-prod

# 查看 Pod 日志
kubectl logs zker-api-xxx -n zker-prod

# 查看所有 Pod 状态
kubectl get pods -n zker-prod --show-all
```

### 镜像拉取失败

```bash
# 检查镜像名称和标签
kubectl get deployment zker-api -n zker-prod -o jsonpath='{.spec.template.spec.containers[0].image}'

# 创建 imagePullSecret
kubectl create secret docker-registry regcred \
  --docker-server=your-registry.com \
  --docker-username=your-username \
  --docker-password=your-password \
  -n zker-prod

# 在 Deployment 中引用
# imagePullSecrets:
# - name: regcred
```

### 服务无法访问

```bash
# 查看 Service Endpoints
kubectl get endpoints zker-api -n zker-prod

# 测试内部访问
kubectl run -it --rm debug --image=busybox --restart=Never -n zker-prod -- \
  wget -O- http://zker-api
```

---

## 📚 相关文档

- [企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [灰度发布策略](../../docs/企业级功能完善与统一性设计方案/ZKER-灰度发布策略_v1.0.md)
- [一键回滚方案](../../docs/企业级功能完善与统一性设计方案/ZKER-一键回滚方案_v1.0.md)
- [故障排查手册](../../docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)

---

**文档版本**: v1.0.0
**最后更新**: 2025-01-01
