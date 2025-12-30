# 生产环境部署指南

**版本**: v1.0.0
**最后更新**: 2025-01-02
**维护者**: DevOps团队

---

## 📋 前置条件

- Kubernetes集群 v1.28+
- kubectl已配置并连接到集群
- Docker镜像仓库访问权限（GitHub Container Registry）
- 监控系统（Prometheus + Grafana）已部署
- 域名和SSL证书已配置

---

## 🚀 部署步骤

### 1. 准备配置文件

```bash
# 克隆配置仓库
git clone https://github.com/coze-studio/deploy-configs.git
cd deploy-configs

# 创建命名空间
kubectl create namespace coze-studio-prod
```

### 2. 创建Secrets

```bash
# 数据库Secret
kubectl create secret generic mysql-secret \
  --from-literal=password='YOUR_MYSQL_PASSWORD' \
  --from-literal=host='mysql.coze-studio.internal' \
  -n coze-studio-prod

# Redis Secret
kubectl create secret generic redis-secret \
  --from-literal=password='YOUR_REDIS_PASSWORD' \
  --from-literal=host='redis.coze-studio.internal' \
  -n coze-studio-prod

# API密钥Secret
kubectl create secret generic api-keys \
  --from-literal=openai-key='YOUR_OPENAI_KEY' \
  --from-literal=anthropic-key='YOUR_ANTHROPIC_KEY' \
  -n coze-studio-prod
```

### 3. 部署后端服务

```bash
# 应用K8s配置
kubectl apply -f k8s/backend/configmap.yaml -n coze-studio-prod
kubectl apply -f k8s/backend/deployment.yaml -n coze-studio-prod
kubectl apply -f k8s/backend/service.yaml -n coze-studio-prod
kubectl apply -f k8s/backend/hpa.yaml -n coze-studio-prod

# 等待部署完成
kubectl rollout status deployment/backend -n coze-studio-prod
```

### 4. 部署前端服务

```bash
kubectl apply -f k8s/frontend/deployment.yaml -n coze-studio-prod
kubectl apply -f k8s/frontend/service.yaml -n coze-studio-prod
kubectl apply -f k8s/frontend/ingress.yaml -n coze-studio-prod

kubectl rollout status deployment/frontend -n coze-studio-prod
```

### 5. 配置Ingress

```bash
kubectl apply -f k8s/ingress.yaml -n coze-studio-prod

# 验证Ingress
kubectl get ingress -n coze-studio-prod
```

### 6. 验证部署

```bash
# 检查Pod状态
kubectl get pods -n coze-studio-prod

# 检查Service状态
kubectl get svc -n coze-studio-prod

# 检查HPA状态
kubectl get hpa -n coze-studio-prod

# 健康检查
kubectl exec -it -n coze-studio-prod deployment/backend -- \
  curl http://localhost:8080/health

# 端口转发测试（可选）
kubectl port-forward -n coze-studio-prod svc/backend 8080:8080
```

---

## 🎯 灰度发布流程

### 1. 准备新版本

```bash
# 标记新版本
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0

# 触发CD流水线（自动）
```

### 2. 阶段1: 金丝雀发布（5%流量）

```bash
# 应用灰度配置
kubectl apply -f deploy/gray-release/canary.yaml -n coze-studio-prod

# 监控金丝雀指标
kubectl get canary backend -n coze-studio-prod -o yaml
```

### 3. 阶段2: 小范围灰度（20%流量）

```bash
# 更新流量权重（Flagger自动处理）
# 观察指标5分钟
```

### 4. 阶段3: 大范围灰度（50%流量）

```bash
# 继续观察指标10分钟
# 检查Grafana Dashboard
```

### 5. 阶段4: 全量发布（100%流量）

```bash
# Flagger自动完成全量发布
# 验证所有Pod健康
kubectl get pods -n coze-studio-prod
```

---

## 🔄 回滚流程

### 自动回滚（灰度失败）

灰度发布期间，如果以下条件触发，Flagger会自动回滚：
- 成功率 < 99%
- P99延迟 > 500ms
- 错误率 > 1%

### 手动回滚

```bash
# 方式1: 使用回滚脚本（推荐）
./scripts/rollback.sh coze-studio-prod backend

# 方式2: 使用kubectl直接回滚
kubectl rollout undo deployment/backend -n coze-studio-prod

# 等待回滚完成
kubectl rollout status deployment/backend -n coze-studio-prod

# 验证回滚结果
kubectl get pods -n coze-studio-prod -l app=backend
```

### 从备份恢复

```bash
# 找到备份文件
ls -lh deploy-backups/backend_*.yaml

# 应用备份
kubectl apply -f deploy-backups/backend_20250102_150000.yaml -n coze-studio-prod

# 验证
kubectl rollout status deployment/backend -n coze-studio-prod
```

---

## 🔧 故障排查

### 问题1: Pod启动失败

```bash
# 查看Pod事件
kubectl describe pod -n coze-studio-prod <pod-name>

# 查看Pod日志
kubectl logs -n coze-studio-prod <pod-name>
kubectl logs -n coze-studio-prod <pod-name> --previous

# 查看所有Pod状态
kubectl get pods -n coze-studio-prod --show-labels
```

### 问题2: 服务无法访问

```bash
# 检查Service
kubectl get svc -n coze-studio-prod -o wide

# 检查Endpoints
kubectl get endpoints -n coze-studio-prod

# 检查Ingress
kubectl get ingress -n coze-studio-prod -o yaml

# 测试Service连通性
kubectl run -it --rm debug --image=busybox --restart=Never -- \
  sh -c "wget -O- http://backend:8080/health"
```

### 问题3: 性能问题

```bash
# 查看资源使用
kubectl top pods -n coze-studio-prod
kubectl top nodes

# 查看HPA状态
kubectl describe hpa backend -n coze-studio-prod

# 查看Prometheus指标
open http://prometheus.coze-studio.com
```

### 问题4: 数据库连接失败

```bash
# 检查Secret
kubectl get secret mysql-secret -n coze-studio-prod -o yaml

# 测试数据库连接
kubectl run -it --rm mysql-client --image=mysql:8.4.5 --restart=Never -- \
  mysql -h mysql.coze-studio.internal -u root -p
```

---

## 📊 监控检查

### Grafana Dashboards

访问以下Dashboard监控系统状态：

- API Performance: http://grafana.coze-studio.com/d/api-performance
- Tenant Business: http://grafana.coze-studio.com/d/tenant-business
- Routing Performance: http://grafana.coze-studio.com/d/routing-performance
- System Resources: http://grafana.coze-studio.com/d/system-resources

### 关键指标

**API性能**:
- QPS > 100 req/s
- P95延迟 < 2s
- 错误率 < 5%

**系统资源**:
- CPU使用率 < 80%
- 内存使用率 < 85%
- 磁盘空间 > 20%

**数据库**:
- 连接数 < 80%
- 慢查询 < 10/s

---

## 🔔 告警通知

### 告警渠道

- Critical告警: #critical-alerts (Slack + PagerDuty)
- Backend告警: #backend-alerts (Slack)
- DevOps告警: #devops-alerts (Slack)

### 告警处理流程

1. 收到告警后，立即查看Grafana Dashboard
2. 确认问题严重程度和影响范围
3. 根据严重程度决定处理方式：
   - P4-P5: 记录问题，定期处理
   - P3: 工作时间内处理
   - P2: 立即处理
   - P1: 紧急处理，必要时回滚

---

## 🕐 维护窗口

**建议维护窗口**: 每周三凌晨2:00-4:00（UTC）

**维护公告**:
- 提前24小时发送维护通知
- 在#operations频道发布公告
- 更新状态页面

---

## 📞 联系方式

- **DevOps团队**: devops@coze-studio.com
- **紧急联系**: +86-xxx-xxxx-xxxx
- **文档**: https://docs.coze-studio.com
- **Runbooks**: https://docs.coze-studio.com/runbooks

---

## 📚 相关文档

- [灰度发布策略](../企业级功能完善与统一性设计方案/ZKER-灰度发布策略_v1.0.md)
- [一键回滚方案](../企业级功能完善与统一性设计方案/ZKER-一键回滚方案_v1.0.md)
- [故障排查手册](../企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)
- [监控告警阈值调优](../企业级功能完善与统一性设计方案/ZKER-监控告警阈值调优指南.md)

---

**祝部署顺利！** 🚀
