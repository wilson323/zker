# ZKER 运维脚本使用手册

> **最后更新**: 2025-01-01
> **维护人**: 研发 D - DevOps 工程师

---

## 📖 概述

本目录包含 ZKER 项目的自动化运维脚本,用于简化日常开发和部署工作。

**脚本列表**:
1. `migrate.sh` - 数据库迁移自动化
2. `canary-deploy.sh` - 灰度发布 (基于 Istio)
3. `rollback.sh` - 一键回滚

---

## 🚀 快速开始

### 1. 数据库迁移 (`migrate.sh`)

**用途**: 自动化数据库 schema 迁移和版本管理

#### 执行迁移

```bash
# 执行所有迁移
./scripts/migrate.sh up

# 执行特定 schema 文件
./scripts/migrate.sh up ./docker/atlas/zker_schema.hcl

# 查看迁移状态
./scripts/migrate.sh status

# 显示数据库信息
./scripts/migrate.sh info
```

#### 创建新迁移

```bash
# 创建迁移模板
./scripts/migrate.sh create add_tenant_id_to_bots

# 编辑迁移文件
vi ./docker/atlas/migrations/20250130_120000_add_tenant_id_to_bots.sql
```

#### 清理旧版本

```bash
# 清理 30 天前的 Atlas 版本记录
./scripts/migrate.sh cleanup
```

---

### 2. 灰度发布 (`canary-deploy.sh`)

**用途**: 基于 Istio 的金丝雀部署,逐步切换流量

#### 手动模式

```bash
# 部署金丝雀版本 (10% 流量)
./scripts/canary-deploy.sh v1.0.1

# 部署金丝雀版本 (20% 流量)
./scripts/canary-deploy.sh v1.0.1 20

# 测试金丝雀版本
curl -H 'x-canary: true' http://zker-api/api/health

# 查看金丝雀日志
kubectl logs -f -n zker-prod -l app=zker-api,version=v1.0.1
```

#### 自动模式

```bash
# 自动逐步增加流量 (10% -> 20% -> ... -> 100%)
./scripts/canary-deploy.sh v1.0.1 auto

# 脚本会自动:
# 1. 部署金丝雀版本
# 2. 每 5 分钟增加 10% 流量
# 3. 监控错误率和 P95 延迟
# 4. 如果指标异常,自动回滚
```

#### 流量调整

手动调整流量比例:

```bash
# 更新到 50% 流量
./scripts/canary-deploy.sh v1.0.1 50

# 更新到 100% 流量 (完全切换)
./scripts/canary-deploy.sh v1.0.1 100
```

---

### 3. 一键回滚 (`rollback.sh`)

**用途**: 快速回滚到上一个稳定版本

#### 回滚到上一个版本

```bash
# 查看当前版本
kubectl get deployment zker-api -n zker-prod \
  -o jsonpath='{.spec.template.metadata.labels.version}'

# 回滚
./scripts/rollback.sh

# 或者指定命名空间和 Deployment
./scripts/rollback.sh -n zker-prod -d zker-api
```

#### 回滚到指定版本

```bash
# 列出可用版本
./scripts/rollback.sh --list

# 回滚到特定版本
./scripts/rollback.sh v1.0.0

# 强制回滚 (不询问确认)
./scripts/rollback.sh --force v1.0.0
```

#### 从备份恢复

```bash
# 列出备份文件
ls -lt ./rollback-backups/

# 从备份恢复
kubectl apply -f ./rollback-backups/zker-api_20250130_143022.yaml
```

---

## 🔧 高级用法

### 环境变量配置

所有脚本都支持环境变量配置:

```bash
# 数据库迁移
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=zker
./scripts/migrate.sh up

# 灰度发布
export NAMESPACE=zker-prod
export SERVICE=zker-api
export PROMETHEUS_URL=http://prometheus:9090
./scripts/canary-deploy.sh v1.0.1 auto

# 回滚
export NAMESPACE=zker-prod
export DEPLOYMENT=zker-api
export FORCE=false
./scripts/rollback.sh
```

### CI/CD 集成

#### GitHub Actions 示例

```yaml
- name: 数据库迁移
  run: ./scripts/migrate.sh up

- name: 灰度发布
  run: ./scripts/canary-deploy.sh ${{ github.ref_name }} auto

- name: 回滚 (失败时)
  if: failure()
  run: ./scripts/rollback.sh --force
```

---

## 📊 监控和日志

### 查看部署状态

```bash
# 查看 Deployment 状态
kubectl get deployments -n zker-prod

# 查看 Pod 状态
kubectl get pods -n zker-prod -l app=zker-api

# 查看 Rollout 历史
kubectl rollout history deployment/zker-api -n zker-prod

# 查看 HPA 状态
kubectl get hpa -n zker-prod
```

### 查看日志

```bash
# 查看所有 Pod 日志
kubectl logs -f -n zker-prod -l app=zker-api --all-containers=true

# 查看特定版本日志
kubectl logs -f -n zker-prod -l app=zker-api,version=v1.0.1

# 查看最近 100 行日志
kubectl logs --tail=100 deployment/zker-api -n zker-prod

# 查看错误日志
kubectl logs deployment/zker-api -n zker-prod | grep -i error
```

### Prometheus 查询

```bash
# 查询错误率
curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~\"5..\"}[5m])" | jq .

# 查询 P95 延迟
curl -s "http://prometheus:9090/api/v1/query?query=histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))" | jq .

# 查询 QPS
curl -s "http://prometheus:9090/api/v1/query?query=rate(http_requests_total[1m])" | jq .
```

---

## 🚨 故障排查

### 问题 1: 数据库迁移失败

**症状**: `migrate.sh up` 报错

**排查步骤**:

1. 检查数据库连接
   ```bash
   mysql -h localhost -u root -p
   ```

2. 检查迁移文件权限
   ```bash
   ls -la ./docker/atlas/migrations/
   ```

3. 手动执行 SQL
   ```bash
   mysql -u root -p zker < ./docker/atlas/migrations/xxx.sql
   ```

### 问题 2: 灰度发布卡住

**症状**: `canary-deploy.sh` 无法继续

**排查步骤**:

1. 检查 Istio 组件
   ```bash
   kubectl get pods -n istio-system
   ```

2. 检查 VirtualService
   ```bash
   kubectl get virtualservice zker-api -n zker-prod -o yaml
   ```

3. 检查金丝雀 Pod
   ```bash
   kubectl get pods -n zker-prod -l version=v1.0.1
   kubectl describe pod <pod-name> -n zker-prod
   ```

4. 查看金丝雀日志
   ```bash
   kubectl logs -f <pod-name> -n zker-prod
   ```

### 问题 3: 回滚失败

**症状**: `rollback.sh` 无法回滚

**排查步骤**:

1. 检查 ReplicaSets
   ```bash
   kubectl get replicasets -n zker-prod -l app=zker-api
   ```

2. 手动回滚
   ```bash
   kubectl rollout undo deployment/zker-api -n zker-prod
   ```

3. 从备份恢复
   ```bash
   ls -lt ./rollback-backups/
   kubectl apply -f ./rollback-backups/<latest-backup>.yaml
   ```

---

## 📚 相关文档

- [企业级开发规范手册](../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [灰度发布策略](../docs/企业级功能完善与统一性设计方案/ZKER-灰度发布策略_v1.0.md)
- [一键回滚方案](../docs/企业级功能完善与统一性设计方案/ZKER-一键回滚方案_v1.0.md)
- [故障排查手册](../docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)
- [Kubernetes 部署指南](../k8s/README.md)

---

## ⚠️ 注意事项

1. **生产环境操作前务必备份**
   - 数据库备份
   - K8s 资源备份
   - 配置文件备份

2. **灰度发布建议在低峰期进行**
   - 监控错误率和延迟
   - 准备快速回滚方案

3. **回滚是最后的手段**
   - 优先修复 Bug 重新发布
   - 回滚可能导致数据不一致

4. **定期测试回滚流程**
   - 确保回滚脚本可用
   - 验证回滚时间 < 5 分钟

---

**文档版本**: v1.0.0
**最后更新**: 2025-01-01
