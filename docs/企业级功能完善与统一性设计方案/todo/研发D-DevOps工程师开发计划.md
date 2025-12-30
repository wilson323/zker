# 研发D - DevOps工程师开发计划

**负责人**: 研发D（DevOps工程师）
**开发周期**: 5周（4人并行）
**总工作量**: 10人天
**核心职责**: CI/CD、Grafana监控、告警规则、灰度发布、部署文档
**最后更新**: 2025-12-30

---

## 📋 工作量总览

| 周次 | 主要任务 | 工作量 | 优先级 |
|------|---------|--------|--------|
| **Week 1-2** | CI/CD流水线优化 + 测试环境 | 2人天 | **P1** |
| **Week 3-4** | Grafana监控大盘 + 告警规则 | 5人天 | **P1** |
| **Week 5** | 灰度发布配置 + 部署文档 | 3人天 | **P2** |

**📌 注意**: 你的工作量最小，但是**生产环境稳定性的保障者**！

---

## ⚠️ 核心注意事项

### 1. 代码隔离原则

**你的职责范围**:
- ✅ CI/CD流水线配置（GitHub Actions）
- ✅ Docker镜像构建优化
- ✅ Kubernetes部署配置
- ✅ Grafana监控大盘配置
- ✅ Prometheus告警规则
- ✅ 灰度发布策略配置
- ✅ 部署文档和运维手册
- ✅ 监控指标设计

**不要触碰**:
- ❌ 后端业务代码（研发A、研发B负责）
- ❌ 前端代码（研发C负责）
- ❌ 数据库迁移脚本（研发B负责）
- ❌ 功能实现代码

### 2. DevOps工作红线

**必须遵守**:
1. **所有配置文件都要版本化**
2. **配置变更要有Review**
3. **生产环境变更要有审批**
4. **变更前必须备份**
5. **变更后必须验证**
6. **所有变更可回滚**

**禁止行为**:
1. 直接在生产环境修改配置
2. 使用未经测试的配置
3. 配置变更没有文档记录
4. 监控告警缺失
5. 变更没有回滚方案

---

## 🎯 Week 1-2: CI/CD流水线优化 (P1)

### 任务1.1: GitHub Actions工作流配置 (Day 1, 1人天)

**文件**: `.github/workflows/ci.yml`

```yaml
# .github/workflows/ci.yml
name: CI - 持续集成

on:
  push:
    branches: [main, develop, 'feature/**']
  pull_request:
    branches: [main, develop]

jobs:
  # 后端测试
  backend-test:
    name: Backend Tests
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.4.5
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: test_db
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

      redis:
        image: redis:8.0
        ports:
          - 6379:6379
        options: >-
          --health-cmd="redis-cli ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
      - name: Checkout代码
        uses: actions/checkout@v3

      - name: 设置Go环境
        uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'
          cache: true

      - name: 下载依赖
        working-directory: ./backend
        run: go mod download

      - name: 运行单元测试
        working-directory: ./backend
        run: |
          go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
          go tool cover -func=coverage.out | grep total

      - name: 上传覆盖率
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out
          flags: backend

      - name: golangci-lint检查
        working-directory: ./backend
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
          golangci-lint run --timeout=5m

      - name: 安全扫描
        uses: securego/gosec@master
        with:
          args: ./backend/...

  # 前端测试
  frontend-test:
    name: Frontend Tests
    runs-on: ubuntu-latest

    steps:
      - name: Checkout代码
        uses: actions/checkout@v3

      - name: 设置Node.js环境
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json

      - name: 安装Rush
        working-directory: .
        run: npm install -g @microsoft/rush

      - name: 安装依赖
        working-directory: frontend
        run: rush update

      - name: TypeScript类型检查
        working-directory: frontend
        run: rushx check:types

      - name: ESLint检查
        working-directory: frontend
        run: rushx lint

      - name: 运行单元测试
        working-directory: frontend
        run: rushx test --coverage

      - name: 构建检查
        working-directory: frontend
        run: rush build

  # Docker镜像构建
  docker-build:
    name: Build Docker Images
    runs-on: ubuntu-latest
    needs: [backend-test, frontend-test]
    if: github.event_name == 'push' && github.ref == 'refs/heads/main'

    steps:
      - name: Checkout代码
        uses: actions/checkout@v3

      - name: 设置Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: 登录Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: 构建并推送后端镜像
        uses: docker/build-push-action@v4
        with:
          context: ./backend
          push: true
          tags: |
            coze-studio/backend:${{ github.sha }}
            coze-studio/backend:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: 构建并推送前端镜像
        uses: docker/build-push-action@v4
        with:
          context: ./frontend/apps/coze-studio
          push: true
          tags: |
            coze-studio/frontend:${{ github.sha }}
            coze-studio/frontend:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

**⚠️ 注意事项**:
1. **测试环境使用MySQL和Redis服务**
2. **缓存Go modules和Node modules**
3. **代码覆盖率上传到Codecov**
4. **golangci-lint检查必须通过**
5. **Docker镜像使用多阶段构建**
6. **使用BuildKit缓存加速构建**

**📖 开发规范**:
- CI流水线标准化
- 缓存策略优化
- 测试覆盖率要求

**🔗 设计文档链接**:
- [ZKER-灰度发布策略_v1.0.md](../ZKER-灰度发布策略_v1.0.md)
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)

---

### 任务1.2: CD部署流水线 (Day 2, 1人天)

**文件**: `.github/workflows/cd.yml`

```yaml
# .github/workflows/cd.yml
name: CD - 持续部署

on:
  push:
    tags:
      - 'v*'

jobs:
  # 部署到测试环境
  deploy-staging:
    name: Deploy to Staging
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.coze-studio.com

    steps:
      - name: Checkout代码
        uses: actions/checkout@v3

      - name: 设置kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'v1.28.0'

      - name: 配置Kubeconfig
        run: |
          mkdir -p ~/.kube
          echo "${{ secrets.KUBE_CONFIG_STAGING }}" | base64 -d > ~/.kube/config

      - name: 更新后端Deployment
        run: |
          kubectl set image deployment/backend \
            backend=coze-studio/backend:${{ github.sha }} \
            -n coze-studio-staging

      - name: 等待后端就绪
        run: |
          kubectl rollout status deployment/backend -n coze-studio-staging

      - name: 更新前端Deployment
        run: |
          kubectl set image deployment/frontend \
            frontend=coze-studio/frontend:${{ github.sha }} \
            -n coze-studio-staging

      - name: 等待前端就绪
        run: |
          kubectl rollout status deployment/frontend -n coze-studio-staging

      - name: 健康检查
        run: |
          kubectl wait --for=condition=available pod \
            -l app=backend -n coze-studio-staging --timeout=300s

      - name: 通知部署结果
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: 'Staging环境部署完成'
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}

  # 部署到生产环境（需手动批准）
  deploy-production:
    name: Deploy to Production
    runs-on: ubuntu-latest
    needs: deploy-staging
    environment:
      name: production
      url: https://coze-studio.com

    steps:
      - name: Checkout代码
        uses: actions/checkout@v3

      - name: 设置kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'v1.28.0'

      - name: 配置Kubeconfig
        run: |
          mkdir -p ~/.kube
          echo "${{ secrets.KUBE_CONFIG_PROD }}" | base64 -d > ~/.kube/config

      - name: 创建备份
        run: |
          kubectl get deployment backend -n coze-studio-prod -o yaml > backup-backend-${{ github.sha }}.yaml

      - name: 更新后端Deployment（金丝雀）
        run: |
          kubectl set image deployment/backend \
            backend=coze-studio/backend:${{ github.sha }} \
            -n coze-studio-prod

      - name: 等待金丝雀就绪
        run: |
          kubectl rollout status deployment/backend -n coze-studio-prod

      - name: 金丝雀健康检查
        run: |
          # 检查金丝雀版本的错误率
          sleep 60  # 等待1分钟收集指标
          ERROR_RATE=$(kubectl exec -n coze-studio-prof deployment/prometheus -- \
            promql --query='rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])')

          if (( $(echo "$ERROR_RATE > 0.05" | bc -l) )); then
            echo "金丝雀版本错误率过高，回滚"
            kubectl rollout undo deployment/backend -n coze-studio-prod
            exit 1
          fi

      - name: 全量发布
        run: |
          # 金丝雀版本健康，全量发布
          kubectl scale deployment backend --replicas=10 -n coze-studio-prod

      - name: 通知部署结果
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          text: '生产环境部署完成'
          webhook_url: ${{ secrets.SLACK_WEBHOOK }}

      - name: 部署失败回滚
        if: failure()
        run: |
          kubectl rollout undo deployment/backend -n coze-studio-prod
```

**⚠️ 注意事项**:
1. **生产环境部署需要手动批准**
2. **金丝雀发布：先发布5%流量**
3. **健康检查：错误率>5%自动回滚**
4. **部署前创建备份**
5. **部署结果通知到Slack**
6. **失败自动回滚**

**📖 开发规范**:
- 生产环境保护
- 金丝雀发布策略
- 自动回滚机制

**🔗 设计文档链接**:
- [ZKER-灰度发布策略_v1.0.md](../ZKER-灰度发布策略_v1.0.md) - 4阶段灰度发布
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md) - 快速回滚方案

---

### 任务1.3: 测试环境搭建 (Day 2, 完成)

**文件**: `docker/docker-compose.test.yml`

```yaml
# docker/docker-compose.test.yml
version: '3.8'

services:
  # MySQL数据库
  mysql:
    image: mysql:8.4.5
    container_name: coze-mysql-test
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: coze_studio_test
      TZ: Asia/Shanghai
    ports:
      - "3307:3306"
    volumes:
      - mysql_test_data:/var/lib/mysql
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
      - --default-authentication-plugin=mysql_native_password
    networks:
      - test-network

  # Redis缓存
  redis:
    image: redis:8.0
    container_name: coze-redis-test
    ports:
      - "6380:6379"
    volumes:
      - redis_test_data:/data
    command: redis-server --appendonly yes
    networks:
      - test-network

  # Elasticsearch
  elasticsearch:
    image: elasticsearch:8.18.0
    container_name: coze-es-test
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - xpack.security.enabled=false
    ports:
      - "9201:9200"
    volumes:
      - es_test_data:/usr/share/elasticsearch/data
    networks:
      - test-network

  # Prometheus监控
  prometheus:
    image: prom/prometheus:v2.45.0
    container_name: coze-prometheus-test
    ports:
      - "9091:9090"
    volumes:
      - ./monitoring/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_test_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
    networks:
      - test-network

  # Grafana可视化
  grafana:
    image: grafana/grafana:10.0.0
    container_name: coze-grafana-test
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_test_data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
    networks:
      - test-network

volumes:
  mysql_test_data:
  redis_test_data:
  es_test_data:
  prometheus_test_data:
  grafana_test_data:

networks:
  test-network:
    driver: bridge
```

**启动命令**:

```bash
# 启动测试环境
docker-compose -f docker/docker-compose.test.yml up -d

# 查看服务状态
docker-compose -f docker/docker-compose.test.yml ps

# 查看日志
docker-compose -f docker/docker-compose.test.yml logs -f

# 停止测试环境
docker-compose -f docker/docker-compose.test.yml down
```

**⚠️ 注意事项**:
1. **测试环境端口与生产环境不同**（避免冲突）
2. **使用独立的数据卷**
3. **Prometheus和Grafana用于测试监控**
4. **环境变量使用测试配置**

**📖 开发规范**:
- 环境隔离
- 配置清晰
- 易于启动

**🔗 设计文档链接**:
- [ZKER-开发快速入门指南.md](../ZKER-开发快速入门指南.md) - 环境配置章节

---

## 🎯 Week 3-4: Grafana监控大盘 + 告警规则 (P1)

### 任务2.1: Grafana Dashboard配置 (Day 1-4, 5人天)

#### Day 1: API性能Dashboard

**文件**: `deploy/monitoring/grafana/dashboards/api-performance.json`

```json
{
  "dashboard": {
    "title": "API Performance Dashboard",
    "tags": ["api", "performance"],
    "timezone": "browser",
    "panels": [
      {
        "title": "Request QPS",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{job=\"coze-studio\"}[1m]))",
            "legendFormat": "QPS"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 0, "y": 0}
      },
      {
        "title": "Response Time (P95)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "P95"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 12, "y": 0}
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))",
            "legendFormat": "Error Rate"
          }
        ],
        "gridPos": {"h": 8, "w": 12, "x": 0, "y": 8}
      },
      {
        "title": "API Response Time by Endpoint",
        "type": "table",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) by (endpoint)",
            "format": "table"
          }
        ],
        "gridPos": {"h": 8, "w": 24, "x": 0, "y": 16}
      }
    ]
  }
}
```

#### Day 2: 租户业务Dashboard

**文件**: `deploy/monitoring/grafana/dashboards/tenant-business.json`

```json
{
  "dashboard": {
    "title": "Tenant Business Dashboard",
    "tags": ["tenant", "business"],
    "panels": [
      {
        "title": "Total Tenants",
        "type": "stat",
        "targets": [
          {
            "expr": "count(tenant_info{status=\"active\"})"
          }
        ]
      },
      {
        "title": "Tenant Subscription Distribution",
        "type": "piechart",
        "targets": [
          {
            "expr": "count by (subscription_tier) (tenant_info)"
          }
        ]
      },
      {
        "title": "Quota Usage Distribution",
        "type": "heatmap",
        "targets": [
          {
            "expr": "quota_usage_percent by (tenant_id, resource_type)"
          }
        ]
      }
    ]
  }
}
```

#### Day 3: 路由性能Dashboard

**文件**: `deploy/monitoring/grafana/dashboards/routing-performance.json`

```json
{
  "dashboard": {
    "title": "Routing Performance Dashboard",
    "tags": ["routing", "performance"],
    "panels": [
      {
        "title": "Routing Decision Latency",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(routing_decision_duration_seconds[5m]))"
          }
        ]
      },
      {
        "title": "Intent Match Confidence",
        "type": "gauge",
        "targets": [
          {
            "expr": "avg(routing_intent_confidence)"
          }
        ]
      },
      {
        "title": "Bot Health Status",
        "type": "table",
        "targets": [
          {
            "expr": "bot_health_status by (bot_id)"
          }
        ]
      }
    ]
  }
}
```

#### Day 4: 系统资源Dashboard

**文件**: `deploy/monitoring/grafana/dashboards/system-resources.json`

```json
{
  "dashboard": {
    "title": "System Resources Dashboard",
    "tags": ["system", "resources"],
    "panels": [
      {
        "title": "CPU Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "100 - (avg by (instance) (irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) * 100)"
          }
        ]
      },
      {
        "title": "Memory Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100"
          }
        ]
      },
      {
        "title": "Disk I/O",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(node_disk_io_time_seconds_total[5m])"
          }
        ]
      },
      {
        "title": "Network Traffic",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(node_network_receive_bytes_total[5m])"
          }
        ]
      }
    ]
  }
}
```

**导入Dashboard**:

```bash
# 使用Grafana API导入Dashboard
curl -X POST http://localhost:3000/api/dashboards/import \
  -H "Content-Type: application/json" \
  -u admin:admin \
  -d @deploy/monitoring/grafana/dashboards/api-performance.json
```

**⚠️ 注意事项**:
1. **每个Dashboard都要有清晰的标题和描述**
2. **Panel布局合理**（使用gridPos）
3. **查询语句优化**（避免过多数据）
4. **使用变量模板**（如$tenant_id）
5. **Dashboard支持自动刷新**

**📖 开发规范**:
- Grafana Dashboard最佳实践
- PromQL查询优化
- 可视化设计

**🔗 设计文档链接**:
- [ZKER-监控告警阈值调优指南.md](../ZKER-监控告警阈值调优指南.md)
- [24-租户监控运维_Agent监控补充.md](../24-租户监控运维_Agent监控补充.md)

---

### 任务2.2: Prometheus告警规则 (Day 5, 1人天)

**文件**: `deploy/monitoring/prometheus/alerts.yml`

```yaml
# deploy/monitoring/prometheus/alerts.yml
groups:
  # API性能告警
  - name: api_performance_alerts
    interval: 30s
    rules:
      - alert: APIHighErrorRate
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m]))
          / sum(rate(http_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: critical
          team: backend
        annotations:
          summary: "API错误率过高"
          description: "API错误率为 {{ $value | humanizePercentage }}，超过5%阈值"

      - alert: APIHighLatency
        expr: |
          histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 3
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "API响应时间过长"
          description: "API P95响应时间为 {{ $value }}s，超过3s阈值"

      - alert: APILowQPS
        expr: sum(rate(http_requests_total[5m])) < 100
        for: 10m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "API QPS过低"
          description: "API QPS为 {{ $value }}/s，低于100/s阈值"

  # 租户业务告警
  - name: tenant_business_alerts
    interval: 1m
    rules:
      - alert: TenantQuotaNearLimit
        expr: quota_usage_percent > 80
        for: 1m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "租户配额即将用完"
          description: "租户 {{ $labels.tenant_id }} 的 {{ $labels.resource_type }} 配额使用率为 {{ $value }}%"

      - alert: TenantQuotaExceeded
        expr: quota_usage_percent >= 95
        for: 1m
        labels:
          severity: critical
          team: backend
        annotations:
          summary: "租户配额已超限"
          description: "租户 {{ $labels.tenant_id }} 的 {{ $labels.resource_type }} 配额使用率为 {{ $value }}%，建议升级订阅"

  # 系统资源告警
  - name: system_resource_alerts
    interval: 1m
    rules:
      - alert: HighCPUUsage
        expr: |
          100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
          team: devops
        annotations:
          summary: "CPU使用率过高"
          description: "实例 {{ $labels.instance }} CPU使用率为 {{ $value }}%"

      - alert: HighMemoryUsage
        expr: |
          (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 5m
        labels:
          severity: warning
          team: devops
        annotations:
          summary: "内存使用率过高"
          description: "实例 {{ $labels.instance }} 内存使用率为 {{ $value }}%"

      - alert: HighDiskIO
        expr: rate(node_disk_io_time_seconds_total[5m]) > 0.8
        for: 5m
        labels:
          severity: warning
          team: devops
        annotations:
          summary: "磁盘I/O过高"
          description: "实例 {{ $labels.instance }} 磁盘I/O使用率为 {{ $value }}%"

      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100 < 20
        for: 5m
        labels:
          severity: critical
          team: devops
        annotations:
          summary: "磁盘空间不足"
          description: "实例 {{ $labels.instance }} 磁盘剩余空间为 {{ $value }}%"

  # 路由性能告警
  - name: routing_performance_alerts
    interval: 1m
    rules:
      - alert: RoutingHighLatency
        expr: |
          histogram_quantile(0.95, rate(routing_decision_duration_seconds[5m])) > 0.2
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "路由决策延迟过高"
          description: "路由决策P95延迟为 {{ $value }}s"

      - alert: BotUnhealthy
        expr: bot_health_status == 0
        for: 2m
        labels:
          severity: critical
          team: backend
        annotations:
          summary: "Bot不健康"
          description: "Bot {{ $labels.bot_id }} 不健康，成功率过低或负载过高"

      - alert: LowIntentConfidence
        expr: avg_over_time(routing_intent_confidence[5m]) < 0.7
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "意图识别置信度过低"
          description: "平均意图置信度为 {{ $value }}，低于70%阈值"

  # 数据库告警
  - name: database_alerts
    interval: 30s
    rules:
      - alert: MySQLHighConnections
        expr: mysql_global_status_threads_connected / mysql_global_variables_max_connections * 100 > 80
        for: 5m
        labels:
          severity: warning
          team: devops
        annotations:
          summary: "MySQL连接数过高"
          description: "MySQL连接使用率为 {{ $value }}%"

      - alert: MySQLSlowQueries
        expr: rate(mysql_global_status_slow_queries[5m]) > 10
        for: 5m
        labels:
          severity: warning
          team: backend
        annotations:
          summary: "MySQL慢查询过多"
          description: "MySQL慢查询率为 {{ $value }}/s"

      - alert: RedisHighMemory
        expr: redis_memory_used_bytes / redis_memory_max_bytes * 100 > 80
        for: 5m
        labels:
          severity: warning
          team: devops
        annotations:
          summary: "Redis内存使用率过高"
          description: "Redis内存使用率为 {{ $value }}%"
```

**⚠️ 注意事项**:
1. **告警规则分级**（warning/critical）
2. **告警持续时间**（for: 5m，避免瞬时告警）
3. **告警消息清晰**，包含关键信息
4. **告警标签完整**（severity/team）
5. **告警阈值合理**（基于性能基线）

**📖 开发规范**:
- 告警规则命名规范
- PromQL查询优化
- 告警消息模板

**🔗 设计文档链接**:
- [ZKER-监控告警阈值调优指南.md](../ZKER-监控告警阈值调优指南.md)
- [ZKER-性能基线文档.md](../ZKER-性能基线文档.md)

---

### 任务2.3: 告警通知配置 (Day 5)

**文件**: `deploy/monitoring/alertmanager/config.yml`

```yaml
# deploy/monitoring/alertmanager/config.yml
global:
  resolve_timeout: 5m
  slack_api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'

route:
  group_by: ['alertname', 'cluster', 'service']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 12h
  receiver: 'default'
  routes:
    # Critical告警立即通知
    - match:
        severity: critical
      receiver: 'critical-alerts'
      continue: true

    # Backend告警发送到后端团队
    - match:
        team: backend
      receiver: 'backend-team'

    # DevOps告警发送到DevOps团队
    - match:
        team: devops
      receiver: 'devops-team'

receivers:
  - name: 'default'
    slack_configs:
      - channel: '#alerts'
        send_resolved: true

  - name: 'critical-alerts'
    slack_configs:
      - channel: '#critical-alerts'
        send_resolved: true
        title: '{{ .GroupLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'

  - name: 'backend-team'
    slack_configs:
      - channel: '#backend-alerts'
        send_resolved: true

  - name: 'devops-team'
    slack_configs:
      - channel: '#devops-alerts'
        send_resolved: true
        pagerduty_configs:
          - service: 'coze-studio-oncall'
```

**⚠️ 注意事项**:
1. **Critical告警立即通知**
2. **按团队分路由**（backend/devops）
3. **告警解决后也要通知**
4. **告警聚合**（相同告警合并）

**📖 开发规范**:
- 告警路由清晰
- 通知渠道合理
- 告警聚合优化

**🔗 设计文档链接**:
- [ZKER-监控告警阈值调优指南.md](../ZKER-监控告警阈值调优指南.md)

---

## 🎯 Week 5: 灰度发布配置 + 部署文档 (P2)

### 任务3.1: 灰度发布策略 (Day 1, 1人天)

**文件**: `deploy/gray-release/strategy.yaml`

```yaml
# deploy/gray-release/strategy.yaml
apiVersion: flagger.app/v1beta1
kind: CanaryStrategy
metadata:
  name: coze-studio-canary
  namespace: coze-studio-prod
spec:
  # 分析目标
  analysis:
    interval: 1m
    threshold: 5
    maxWeight: 50
    stepWeight: 10
    metrics:
      - name: request-success-rate
        thresholdRange:
          min: 99
          max: 100
        query: |
          sum(rate(http_request_duration_seconds_bucket{status!~"5.."}[5m]))
          /
          sum(rate(http_request_duration_seconds_bucket[5m]))

      - name: request-duration
        thresholdRange:
          min: 0
          max: 500
        query: |
          histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

      - name: request-error-rate
        thresholdRange:
          min: 0
          max: 1
        query: |
          sum(rate(http_requests_total{status=~"5.."}[5m]))
          /
          sum(rate(http_requests_total[5m]))

  # 灰度阶段
  steps:
    # 阶段1: 金丝雀（5%流量）
    - setWeight: 5
      pause: { duration: 2m }
      metrics:
        - name: request-success-rate
          thresholdRange: { min: 99, max: 100 }

    # 阶段2: 小范围灰度（20%流量）
    - setWeight: 20
      pause: { duration: 5m }
      metrics:
        - name: request-success-rate
          thresholdRange: { min: 99, max: 100 }

    # 阶段3: 大范围灰度（50%流量）
    - setWeight: 50
      pause: { duration: 10m }
      metrics:
        - name: request-success-rate
          thresholdRange: { min: 99, max: 100 }

    # 阶段4: 全量发布（100%流量）
    - setWeight: 100
      pause: { duration: 5m }

  # 自动回滚条件
  webhook:
    - url: http://flagger-prometheus/prometheus
      timeout: 5s
```

**⚠️ 注意事项**:
1. **4阶段灰度**：5% → 20% → 50% → 100%
2. **每个阶段有持续时间**（观察效果）
3. **自动回滚**：错误率>5%或延迟>500ms
4. **金丝雀发布**：先发布到小范围用户

**📖 开发规范**:
- 灰度发布策略
- 自动回滚机制
- 监控指标完整

**🔗 设计文档链接**:
- [ZKER-灰度发布策略_v1.0.md](../ZKER-灰度发布策略_v1.0.md) - **必读！**

---

### 任务3.2: 一键回滚脚本 (Day 2, 1人天)

**文件**: `scripts/rollback.sh`

```bash
#!/bin/bash
# scripts/rollback.sh - 一键回滚脚本

set -e

# 配置
NAMESPACE="${1:-coze-studio-prod}"
DEPLOYMENT="${2:-backend}"
BACKUP_FILE="${3:-backup-latest.yaml}"

echo "=========================================="
echo "  ZKER 一键回滚脚本"
echo "=========================================="
echo "命名空间: $NAMESPACE"
echo "部署: $DEPLOYMENT"
echo "备份文件: $BACKUP_FILE"
echo ""

# 1. 确认回滚
read -p "确定要回滚吗？(yes/no): " confirm
if [ "$confirm" != "yes" ]; then
  echo "取消回滚"
  exit 0
fi

# 2. 检查备份文件
if [ ! -f "$BACKUP_FILE" ]; then
  echo "错误: 备份文件不存在: $BACKUP_FILE"
  exit 1
fi

# 3. 回滚Deployment
echo "开始回滚..."
kubectl rollout undo deployment/$DEPLOYMENT -n $NAMESPACE

# 4. 等待回滚完成
echo "等待回滚完成..."
kubectl rollout status deployment/$DEPLOYMENT -n $NAMESPACE

# 5. 验证回滚
echo "验证回滚结果..."
kubectl wait --for=condition=available pod \
  -l app=$DEPLOYMENT -n $NAMESPACE --timeout=300s

# 6. 显示回滚版本
REVISION=$(kubectl rollout history deployment/$DEPLOYMENT -n $NAMESPACE | tail -1)
echo "当前版本:"
echo "$REVISION"

echo ""
echo "=========================================="
echo "  回滚完成！"
echo "=========================================="
```

**使用方法**:

```bash
# 回滚后端
./scripts/rollback.sh coze-studio-prod backend backup-backend-abc123.yaml

# 回滚前端
./scripts/rollback.sh coze-studio-prod frontend backup-frontend-abc123.yaml
```

**⚠️ 注意事项**:
1. **回滚前必须确认**
2. **回滚要有备份文件**
3. **回滚后验证Pod状态**
4. **显示当前版本信息**

**📖 开发规范**:
- 脚本要有确认步骤
- 错误处理完整
- 日志输出清晰

**🔗 设计文档链接**:
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)

---

### 任务3.3: 部署文档 (Day 3, 1人天)

**文件**: `docs/deployment/PRODUCTION_DEPLOYMENT_GUIDE.md`

```markdown
# 生产环境部署指南

## 前置条件

- Kubernetes集群 v1.28+
- Docker镜像仓库访问权限
- kubectl配置完成
- 监控系统（Prometheus + Grafana）已部署

## 部署步骤

### 1. 准备配置文件

```bash
# 克隆配置仓库
git clone https://github.com/coze-studio/deploy-configs.git
cd deploy-configs
```

### 2. 创建命名空间

```bash
kubectl create namespace coze-studio-prod
```

### 3. 创建Secrets

```bash
# 数据库Secret
kubectl create secret generic mysql-secret \
  --from-literal=password='YOUR_MYSQL_PASSWORD' \
  -n coze-studio-prod

# Redis Secret
kubectl create secret generic redis-secret \
  --from-literal=password='YOUR_REDIS_PASSWORD' \
  -n coze-studio-prod
```

### 4. 部署应用

```bash
# 部署后端
kubectl apply -f k8s/backend/deployment.yaml -n coze-studio-prod

# 部署前端
kubectl apply -f k8s/frontend/deployment.yaml -n coze-studio-prod

# 部署Service
kubectl apply -f k8s/backend/service.yaml -n coze-studio-prod
kubectl apply -f k8s/frontend/service.yaml -n coze-studio-prod

# 部署Ingress
kubectl apply -f k8s/ingress.yaml -n coze-studio-prod
```

### 5. 验证部署

```bash
# 检查Pod状态
kubectl get pods -n coze-studio-prod

# 检查服务状态
kubectl get svc -n coze-studio-prod

# 检查Ingress
kubectl get ingress -n coze-studio-prod

# 健康检查
kubectl exec -it -n coze-studio-prod deployment/backend -- curl http://localhost:8001/health
```

### 6. 监控检查

```bash
# 访问Grafana
open http://grafana.coze-studio.com

# 访问Prometheus
open http://prometheus.coze-studio.com

# 检查告警
kubectl get prometheusrules -n coze-studio-prod
```

## 灰度发布

### 1. 金丝雀发布（5%流量）

```bash
kubectl apply -f deploy/gray-release/canary-5pct.yaml
```

### 2. 小范围灰度（20%流量）

```bash
kubectl apply -f deploy/gray-release/canary-20pct.yaml
```

### 3. 大范围灰度（50%流量）

```bash
kubectl apply -f deploy/gray-release/canary-50pct.yaml
```

### 4. 全量发布（100%流量）

```bash
kubectl apply -f deploy/gray-release/canary-100pct.yaml
```

## 回滚

### 自动回滚（灰度失败）

灰度发布期间，如果错误率>5%，自动回滚到上一个版本。

### 手动回滚

```bash
# 使用回滚脚本
./scripts/rollback.sh coze-studio-prod backend backup-latest.yaml

# 或使用kubectl
kubectl rollout undo deployment/backend -n coze-studio-prod
```

## 故障排查

### 问题1: Pod启动失败

```bash
# 查看Pod事件
kubectl describe pod -n coze-studio-prod <pod-name>

# 查看Pod日志
kubectl logs -n coze-studio-prod <pod-name>
```

### 问题2: 服务无法访问

```bash
# 检查Service
kubectl get svc -n coze-studio-prod

# 检查Endpoints
kubectl get endpoints -n coze-studio-prod

# 检查Ingress
kubectl get ingress -n coze-studio-prod
```

### 问题3: 性能问题

```bash
# 查看资源使用
kubectl top pods -n coze-studio-prod
kubectl top nodes

# 查看Prometheus指标
open http://prometheus.coze-studio.com
```

## 监控告警

### Grafana Dashboards

- API Performance Dashboard: http://grafana.coze-studio.com/d/api-performance
- Tenant Business Dashboard: http://grafana.coze-studio.com/d/tenant-business
- System Resources Dashboard: http://grafana.coze-studio.com/d/system-resources

### 告警规则

- Critical告警: #critical-alerts
- Backend告警: #backend-alerts
- DevOps告警: #devops-alerts

## 维护窗口

建议维护窗口: 每周三凌晨2:00-4:00（UTC）

## 联系方式

- DevOps团队: devops@coze-studio.com
- 紧急联系: +86-xxx-xxxx-xxxx
```

**⚠️ 注意事项**:
1. **部署步骤清晰详细**
2. **包含故障排查章节**
3. **监控和告警链接**
4. **回滚步骤完整**
5. **联系信息完整**

**📖 开发规范**:
- 文档结构清晰
- 步骤可执行
- 命令示例完整
- 故障排查详细

**🔗 设计文档链接**:
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)
- [ZKER-灰度发布策略_v1.0.md](../ZKER-灰度发布策略_v1.0.md)
- [ZKER-故障排查手册_v1.0.md](../ZKER-故障排查手册_v1.0.md)

---

## 📊 关键文档链接汇总

### 必读文档（开发前必读）

1. **[ZKER-灰度发布策略_v1.0.md](../ZKER-灰度发布策略_v1.0.md)** ⭐⭐⭐
   - 4阶段灰度发布策略
   - 自动回滚机制

2. **[ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)** ⭐⭐⭐
   - 5分钟快速回滚
   - 回滚脚本完整

3. **[ZKER-监控告警阈值调优指南.md](../ZKER-监控告警阈值调优指南.md)** ⭐⭐⭐
   - Prometheus告警规则
   - Grafana Dashboard配置

4. **[ZKER-故障排查手册_v1.0.md](../ZKER-故障排查手册_v1.0.md)** ⭐⭐
   - 5大故障场景SOP
   - 故障排查流程

### 参考文档

5. **[24-租户监控运维_Agent监控补充.md](../24-租户监控运维_Agent监控补充.md)** - Agent监控指标
6. **[ZKER-性能基线文档.md](../ZKER-性能基线文档.md)** - 性能基线数据
7. **[ZKER-开发快速入门指南.md](../ZKER-开发快速入门指南.md)** - 环境配置

---

## ⚠️ 核心注意事项（重申）

### 1. 代码隔离原则（重申）

**你的职责范围**:
- ✅ CI/CD流水线配置
- ✅ Docker镜像构建优化
- ✅ Kubernetes部署配置
- ✅ Grafana监控大盘
- ✅ Prometheus告警规则
- ✅ 灰度发布配置
- ✅ 部署文档和运维手册
- ✅ 监控指标设计

**不要触碰**:
- ❌ 后端业务代码（研发A、研发B负责）
- ❌ 前端代码（研发C负责）
- ❌ 数据库迁移脚本（研发B负责）

### 2. DevOps工作红线（重申）

**必须遵守**:
1. 所有配置文件都要版本化
2. 配置变更要有Review
3. 生产环境变更要有审批
4. 变更前必须备份
5. 变更后必须验证
6. 所有变更可回滚

**禁止行为**:
1. 直接在生产环境修改配置
2. 使用未经测试的配置
3. 配置变更没有文档记录
4. 监控告警缺失
5. 变更没有回滚方案

### 3. Git提交规范

**Commit Message格式**:
```
feat(devops): configure Grafana dashboards for monitoring

- Add API Performance Dashboard
- Add Tenant Business Dashboard
- Add System Resources Dashboard
- Configure Prometheus alert rules

Refs: #123
```

---

## 📅 每日工作检查清单

### 开发前
- [ ] 阅读相关设计文档
- [ ] 理解架构和部署策略
- [ ] 准备测试环境
- [ ] 备份当前配置

### 开发中
- [ ] 配置文件版本化
- [ ] 监控指标完整
- [ ] 告警规则合理
- [ ] 文档同步更新

### 提交前
- [ ] 配置在测试环境验证
- [ ] 监控和告警正常
- [ ] 文档完整
- [ ] 回滚方案确认

---

## 🎯 成功标准

### Week 1-2结束时（CI/CD）
- [ ] CI流水线正常运行
- [ ] 单元测试自动运行
- [ ] Docker镜像自动构建
- [ ] CD流水线部署到测试环境

### Week 3-4结束时（监控告警）
- [ ] 4个Grafana Dashboard创建完成
- [ ] 50+条告警规则配置完成
- [ ] 告警通知到Slack正常
- [ ] 监控指标数据正常

### Week 5结束时（灰度+文档）
- [ ] 4阶段灰度发布配置完成
- [ ] 一键回滚脚本完成
- [ ] 部署文档完整
- [ ] 运维手册完整

---

## 🚀 快速启动指南

### 第一天任务

1. **阅读必读文档**（2小时）
   - ZKER-灰度发布策略_v1.0.md
   - ZKER-一键回滚方案_v1.0.md
   - ZKER-监控告警阈值调优指南.md

2. **设置CI/CD流水线**（4小时）
   - 创建GitHub Actions工作流
   - 配置Docker镜像构建
   - 测试流水线运行

3. **配置Grafana**（2小时）
   - 部署Grafana到测试环境
   - 导入API Performance Dashboard
   - 验证监控数据

### 第一周目标

- [ ] CI/CD流水线正常运行
- [ ] 测试环境搭建完成
- [ ] 至少1个Grafana Dashboard完成

---

**🎉 你的工作保障生产环境稳定性，责任重大！**
**所有配置变更都要谨慎，充分测试后再上生产！**
**遇到问题随时查阅设计文档或与团队讨论。**
