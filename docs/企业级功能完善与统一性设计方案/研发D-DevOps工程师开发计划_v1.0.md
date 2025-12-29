# 研发 D - DevOps 工程师 8 周开发计划 v1.0

> **角色定位**: DevOps 工程师 - 负责 Docker 环境、K8s 部署、CI/CD 流水线、数据迁移、灰度发布、一键回滚
>
> **工作目标**: 构建自动化、可靠的交付和运维体系，确保零停机发布和快速回滚能力
>
> **核心原则**: 自动化优先、基础设施即代码、可回滚性、监控可观测

---

## 📋 个人职责概述

### 核心负责模块

```
docker/
├── docker-compose.yml        ✅ 专属负责 - 本地开发环境
├── docker-compose.prod.yml   ✅ 专属负责 - 生产环境
└── Dockerfile.*              ✅ 专属负责 - 各服务镜像

k8s/
├── base/                     ✅ 专属负责 - K8s 基础配置
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
├── overlays/
│   ├── dev/                  ✅ 专属负责 - 开发环境
│   ├── staging/              ✅ 专属负责 - 预发布环境
│   └── prod/                 ✅ 专属负责 - 生产环境

scripts/
├── deploy.sh                 ✅ 专属负责 - 部署脚本
├── rollback.sh               ✅ 专属负责 - 回滚脚本
├── migrate.sh                ✅ 专属负责 - 迁移脚本
└── health-check.sh            ✅ 专属负责 - 健康检查

.github/
└── workflows/                ✅ 专属负责 - CI/CD 流水线
    ├── ci.yml
    ├── cd.yml
    └── release.yml
```

### 协作接口

| 协作对象 | 协作内容 | 接口定义位置 | 依赖关系 |
|---------|---------|------------|---------|
| **研发 A** | 数据库迁移脚本、环境变量清单 | `migrations/`、`docker/.env.example` | 研发 A 提供 Schema → 研发 D 编写迁移 |
| **研发 B** | 监控服务部署、告警配置 | `docker/prometheus/`、`docker/grafana/` | 研发 B 提供配置 → 研发 D 部署 |
| **研发 C** | 前端构建配置、静态资源部署 | `frontend/rsbuild.config.ts` | 研发 C 提供构建产物 → 研发 D 部署 |

---

## 🎯 8 周详细开发计划

### Week 1-2: Docker 环境配置

#### Week 1: 本地开发环境 Docker化

**目标**: 完善 Docker Compose 配置，一键启动所有服务

##### Day 1-3: 完善服务配置

**docker-compose.yml**:
```yaml
version: '3.8'

services:
  # MySQL
  mysql:
    image: mysql:8.4.5
    container_name: zker-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: ${MYSQL_DATABASE}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./docker/mysql/conf.d:/etc/mysql/conf.d
      - ./docker/mysql/init:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - zker-network

  # Redis
  redis:
    image: redis:8.0-alpine
    container_name: zker-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./docker/redis/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - zker-network

  # Elasticsearch
  elasticsearch:
    image: elasticsearch:8.18.0
    container_name: zker-elasticsearch
    restart: unless-stopped
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost:9200/_cluster/health || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 5
    networks:
      - zker-network

  # Milvus (向量数据库)
  milvus:
    image: milvusdb/milvus:v2.5.10-gpu
    container_name: zker-milvus
    restart: unless-stopped
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio:9000
    ports:
      - "19530:19530"
    depends_on:
      - etcd
      - minio
    networks:
      - zker-network

  # MinIO (对象存储)
  minio:
    image: minio/minio:latest
    container_name: zker-minio
    restart: unless-stopped
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    command: server /data --console-address ":9001"
    networks:
      - zker-network

  # etcd (配置中心)
  etcd:
    image: quay.io/coreos/etcd:v3.5.9
    container_name: zker-etcd
    restart: unless-stopped
    ports:
      - "2379:2379"
    volumes:
      - etcd_data:/etcd-data
    command:
      - /usr/local/bin/etcd
      - --name=etcd0
      - --data-dir=/etcd-data
      - --listen-client-urls=http://0.0.0.0:2379
      - --advertise-client-urls=http://localhost:2379
    networks:
      - zker-network

  # NSQ (消息队列)
  nsqlookupd:
    image: nsqio/nsq:v1.2.1
    container_name: zker-nsqlookupd
    restart: unless-stopped
    ports:
      - "4160:4160"
      - "4161:4161"
    command: /nsqlookupd
    networks:
      - zker-network

  nsqd:
    image: nsqio/nsq:v1.2.1
    container_name: zker-nsqd
    restart: unless-stopped
    ports:
      - "4150:4150"
      - "4151:4151"
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    depends_on:
      - nsqlookupd
    networks:
      - zker-network

  # Prometheus (监控)
  prometheus:
    image: prom/prometheus:v2.45.0
    container_name: zker-prometheus
    restart: unless-stopped
    ports:
      - "9090:9090"
    volumes:
      - ./docker/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./docker/prometheus/alerts.yml:/etc/prometheus/alerts.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
    networks:
      - zker-network

  # Grafana (可视化)
  grafana:
    image: grafana/grafana:10.0.0
    container_name: zker-grafana
    restart: unless-stopped
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}
    volumes:
      - grafana_data:/var/lib/grafana
      - ./docker/grafana/provisioning:/etc/grafana/provisioning
    depends_on:
      - prometheus
    networks:
      - zker-network

  # Jaeger (分布式追踪)
  jaeger:
    image: jaegertracing/all-in-one:1.50
    container_name: zker-jaeger
    restart: unless-stopped
    ports:
      - "5775:5775/udp"
      - "6831:6831/udp"
      - "6832:6832/udp"
      - "5778:5778"
      - "16686:16686"
      - "14268:14268"
      - "14250:14250"
      - "9411:9411"
    environment:
      - COLLECTOR_ZIPKIN_HOST_PORT=:9411
    networks:
      - zker-network

volumes:
  mysql_data:
  redis_data:
  es_data:
  milvus_data:
  minio_data:
  etcd_data:
  prometheus_data:
  grafana_data:

networks:
  zker-network:
    driver: bridge
```

##### Day 4-5: 环境变量配置

**.env.example**:
```bash
# MySQL
MYSQL_ROOT_PASSWORD=root_password
MYSQL_DATABASE=zker

# MinIO
MINIO_ROOT_USER=admin
MINIO_ROOT_PASSWORD=minio_password

# Grafana
GRAFANA_ADMIN_PASSWORD=admin

# 应用配置
API_PORT=8080
ENVIRONMENT=development
LOG_LEVEL=debug
```

#### Week 2: 数据迁移自动化

**目标**: 实现数据库迁移的自动化和版本管理

##### Day 1-3: 迁移脚本框架

**migrate.sh**:
```bash
#!/bin/bash
# scripts/migrate.sh

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
MIGRATIONS_DIR="./migrations"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-root_password}"
DB_NAME="${DB_NAME:-zker}"

# 函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查迁移文件
check_migrations() {
    if [ ! -d "$MIGRATIONS_DIR" ]; then
        log_error "Migrations directory not found: $MIGRATIONS_DIR"
        exit 1
    fi

    log_info "Found migrations directory: $MIGRATIONS_DIR"
}

# 执行迁移
run_migration() {
    local migration_file=$1
    local migration_name=$(basename "$migration_file" .sql)

    log_info "Running migration: $migration_name"

    # 执行 SQL
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$migration_file"

    if [ $? -eq 0 ]; then
        log_info "Migration $migration_name completed successfully"
    else
        log_error "Migration $migration_name failed"
        exit 1
    fi
}

# 主函数
main() {
    local action=${1:-up}
    local version=${2:-latest}

    log_info "Starting database migration: $action"

    case $action in
        up)
            check_migrations
            for migration in $(ls -v "$MIGRATIONS_DIR"/*.sql 2>/dev/null || true); do
                run_migration "$migration"
            done
            ;;
        down)
            log_warn "Rollback to version: $version"
            # 实现回滚逻辑
            ;;
        status)
            log_info "Migration status:"
            # 查询迁移状态表
            mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" \
              -e "SELECT * FROM schema_migrations ORDER BY version;"
            ;;
        *)
            log_error "Unknown action: $action"
            echo "Usage: $0 {up|down|status} [version]"
            exit 1
            ;;
    esac

    log_info "Migration completed"
}

main "$@"
```

##### Day 4-5: 迁移文件示例

**migrations/001_create_tenants_table.sql**:
```sql
-- 001_create_tenants_table.sql
-- 租户表

CREATE TABLE IF NOT EXISTS tenants (
    tenant_id VARCHAR(36) PRIMARY KEY COMMENT '租户ID',
    tenant_name VARCHAR(200) NOT NULL COMMENT '租户名称',
    tenant_type ENUM('individual', 'team', 'enterprise') NOT NULL COMMENT '租户类型',
    status ENUM('active', 'suspended', 'deleted') DEFAULT 'active' COMMENT '状态',
    subscription_tier ENUM('free', 'pro', 'enterprise') DEFAULT 'free' COMMENT '订阅等级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_tenant_name (tenant_name, deleted_at),
    INDEX idx_status_type (status, tenant_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户表';

-- 记录迁移版本
INSERT INTO schema_migrations (version, name, applied_at)
VALUES (1, '001_create_tenants_table', NOW())
ON DUPLICATE KEY UPDATE applied_at = NOW();
```

---

### Week 3-4: K8s 集群部署

#### Week 3: K8s 基础配置

**目标**: 编写 K8s YAML 配置，实现声明式部署

##### Day 1-3: Deployment 配置

**k8s/base/deployment.yaml**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zker-api
  labels:
    app: zker-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: zker-api
  template:
    metadata:
      labels:
        app: zker-api
    spec:
      containers:
      - name: zker-api
        image: zker/api:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

**k8s/base/service.yaml**:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: zker-api
spec:
  selector:
    app: zker-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: zker-api-nodeport
spec:
  selector:
    app: zker-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: NodePort
```

##### Day 4-5: ConfigMap 和 Secret

**k8s/base/configmap.yaml**:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: zker-config
data:
  ENVIRONMENT: "production"
  LOG_LEVEL: "info"
  REDIS_HOST: "redis:6379"
  ES_HOST: "http://elasticsearch:9200"
```

**k8s/base/secret.yaml**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
stringData:
  host: "mysql.zker.svc.cluster.local"
  user: "root"
  password: "your_password_here"
  database: "zker"
---
apiVersion: v1
kind: Secret
metadata:
  name: api-keys
type: Opaque
stringData:
  openai-api-key: "your_openai_key"
  anthropic-api-key: "your_anthropic_key"
```

#### Week 4: Kustomize 配置

**目标**: 使用 Kustomize 管理多环境配置

##### Day 1-5: 环境配置

**k8s/overlays/dev/kustomization.yaml**:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: zker-dev

resources:
- ../../base

patchesStrategicMerge:
- deployment-patch.yaml

configMapGenerator:
- name: zker-config
  literals:
  - ENVIRONMENT=development
  - LOG_LEVEL=debug
```

**k8s/overlays/prod/kustomization.yaml**:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: zker-prod

resources:
- ../../base

patchesStrategicMerge:
- deployment-patch.yaml
- hpa-patch.yaml

configMapGenerator:
- name: zker-config
  literals:
  - ENVIRONMENT=production
  - LOG_LEVEL=info

images:
- name: zker/api
  newTag: v1.0.0
```

**k8s/overlays/prod/hpa-patch.yaml**:
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: zker-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: zker-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

---

### Week 5-6: CI/CD 流水线

#### Week 5: CI 流水线

**目标**: 实现代码提交后的自动构建和测试

##### Day 1-5: GitHub Actions CI

**.github/workflows/ci.yml**:
```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  # 后端测试
  backend-test:
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.4.5
        env:
          MYSQL_ROOT_PASSWORD: test_password
          MYSQL_DATABASE: test_zker
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Download dependencies
      working-directory: ./backend
      run: go mod download

    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        working-directory: ./backend

    - name: Run tests
      working-directory: ./backend
      run: |
        go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
        go tool cover -html=coverage.out -o coverage.html

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./backend/coverage.out

  # 前端测试
  frontend-test:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'

    - name: Install Rush
      run: npm install -g @microsoft/rush

    - name: Install dependencies
      run: rush update

    - name: Run linter
      run: rush lint

    - name: Run tests
      run: rush test

    - name: Build
      run: rush build

  # Docker 构建测试
  docker-build:
    runs-on: ubuntu-latest
    needs: [backend-test, frontend-test]

    steps:
    - uses: actions/checkout@v3

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Build API image
      uses: docker/build-push-action@v4
      with:
        context: ./backend
        file: ./backend/Dockerfile
        push: false
        tags: zker/api:test
```

#### Week 6: CD 流水线

**目标**: 实现自动部署到 K8s 集群

##### Day 1-5: GitHub Actions CD

**.github/workflows/cd.yml**:
```yaml
name: CD

on:
  push:
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Configure AWS credentials
      uses: aws-actions/configure-aws-credentials@v2
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-west-2

    - name: Login to Amazon ECR
      id: login-ecr
      uses: aws-actions/amazon-ecr-login@v1

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2

    - name: Build and push API image
      uses: docker/build-push-action@v4
      with:
        context: ./backend
        file: ./backend/Dockerfile
        push: true
        tags: |
          ${{ steps.login-ecr.outputs.registry }}/zker-api:${{ github.ref_name }}
          ${{ steps.login-ecr.outputs.registry }}/zker-api:latest
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Build and push Frontend image
      uses: docker/build-push-action@v4
      with:
        context: ./frontend
        push: true
        tags: |
          ${{ steps.login-ecr.outputs.registry }}/zker-frontend:${{ github.ref_name }}
          ${{ steps.login-ecr.outputs.registry }}/zker-frontend:latest

    - name: Setup kubectl
      uses: azure/setup-kubectl@v3

    - name: Update kubeconfig
      run: aws eks update-kubeconfig --name zker-prod --region us-west-2

    - name: Deploy to K8s
      run: |
        cd k8s/overlays/prod
        kustomize edit set image zker/api=${{ steps.login-ecr.outputs.registry }}/zker-api:${{ github.ref_name }}
        kustomize edit set image zker/frontend=${{ steps.login-ecr.outputs.registry }}/zker-frontend:${{ github.ref_name }}
        kubectl apply -k .

    - name: Verify deployment
      run: |
        kubectl rollout status deployment/zker-api -n zker-prod
        kubectl get pods -n zker-prod
```

---

### Week 7-8: 灰度发布 + 一键回滚

#### Week 7: 灰度发布策略

**目标**: 实现基于 Istio 的灰度发布

##### Day 1-3: Istio 配置

**VirtualService**:
```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  hosts:
  - zker-api
  http:
  - match:
    - headers:
        x-canary:
          exact: "true"
    route:
    - destination:
        host: zker-api
        subset: v2  # 新版本
      weight: 100
  - route:
    - destination:
        host: zker-api
        subset: v1  # 旧版本
      weight: 90
    - destination:
        host: zker-api
        subset: v2  # 新版本
      weight: 10   # 10% 流量到新版本
```

**DestinationRule**:
```yaml
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: zker-api
spec:
  host: zker-api
  subsets:
  - name: v1
    labels:
      version: v1.0.0
  - name: v2
    labels:
      version: v2.0.0
```

##### Day 4-5: 灰度发布脚本

**canary-deploy.sh**:
```bash
#!/bin/bash
# scripts/canary-deploy.sh

set -e

VERSION=$1
CANARY_PERCENT=${2:-10}

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version> [canary_percent]"
    exit 1
fi

log_info() {
    echo -e "\033[0;32m[INFO]\033[0m $1"
}

log_error() {
    echo -e "\033[0;31m[ERROR]\033[0m $1"
}

# 1. 部署新版本（金丝雀）
log_info "Deploying new version: $VERSION (canary)"
kubectl apply -f k8s/overlays/canary/

# 2. 等待新版本就绪
log_info "Waiting for new version to be ready..."
kubectl wait --for=condition=available --timeout=5m \
  deployment/zker-api-canary -n zker-prod

# 3. 更新 VirtualService，设置流量比例
log_info "Setting canary traffic to $CANARY_PERCENT%"
cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
  namespace: zker-prod
spec:
  http:
  - route:
    - destination:
        host: zker-api
        subset: v1
      weight: $((100 - CANARY_PERCENT))
    - destination:
        host: zker-api
        subset: v2
      weight: $CANARY_PERCENT
EOF

log_info "Canary deployment started. Monitor with:"
echo "  kubectl get pods -n zker-prod"
echo "  kubectl logs -f -n zker-prod -l app=zker-api,version=$VERSION"
```

#### Week 8: 一键回滚方案

**目标**: 实现 5 分钟内快速回滚

##### Day 1-3: 回滚脚本

**rollback.sh**:
```bash
#!/bin/bash
# scripts/rollback.sh

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo ""
    echo "Available versions:"
    kubectl get deployments -n zker-prod -o custom-columns="VERSION:.metadata.labels.version,.NAME:.metadata.name" | grep -v VERSION
    exit 1
fi

log_warn() {
    echo -e "\033[1;33m[WARN]\033[0m $1"
}

log_info() {
    echo -e "\033[0;32m[INFO]\033[0m $1"
}

# 1. 确认回滚
log_warn "This will rollback to version: $VERSION"
read -p "Are you sure? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Rollback cancelled"
    exit 0
fi

# 2. 记录当前版本
CURRENT_VERSION=$(kubectl get deployment zker-api -n zker-prod -o jsonpath='{.spec.template.metadata.labels.version}')
log_info "Current version: $CURRENT_VERSION"
log_info "Rolling back to: $VERSION"

# 3. 执行回滚
log_info "Rolling back deployment..."
kubectl rollout undo deployment/zker-api -n zker-prod

# 或者使用特定版本：
# kubectl set image deployment/zker-api zker-api=zker/api:$VERSION -n zker-prod

# 4. 等待回滚完成
log_info "Waiting for rollback to complete..."
kubectl rollout status deployment/zker-api -n zker-prod --timeout=5m

# 5. 验证
log_info "Verifying rollback..."
kubectl get pods -n zker-prod -l app=zker-api

log_info "Rollback completed successfully!"
```

##### Day 4-5: 自动化回滚

**自动回滚 CronJob**:
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: auto-rollback
spec:
  schedule: "*/5 * * * *"  # 每 5 分钟检查一次
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: auto-rollback
            image: zker/auto-rollback:latest
            env:
            - name: ERROR_THRESHOLD
              value: "5"  # 错误率阈值 5%
            - name: LATENCY_THRESHOLD
              value: "1000"  # 延迟阈值 1000ms
          restartPolicy: OnFailure
```

---

## 📦 专属开发规范

### 作为 DevOps 的特殊要求

#### 1. 基础设施即代码

✅ **必须**：
- 所有配置都通过 YAML 管理
- 版本控制所有配置文件
- 不可变基础设施

❌ **禁止**：
- 手动修改服务器配置
- 在生产环境直接操作
- 配置文件不进版本库

#### 2. 安全原则

✅ **必须**：
- 敏感信息使用 Secret/环境变量
- 最小权限原则
- 定期轮换密钥

```yaml
# ✅ Good
env:
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: db-secret
      key: password

# ❌ Bad
env:
- name: DB_PASSWORD
  value: "password123"  # 硬编码密码
```

---

## 🤝 协作接口定义

### 与研发 A 协作

**依赖研发 A**:
- 数据库 Schema
- 环境变量清单
- 健康检查端点

**向研发 A 提供**:
- 数据库迁移脚本模板
- 部署文档

### 与研发 B 协作

**依赖研发 B**:
- Prometheus 配置
- Grafana Dashboard
- 告警规则

**向研发 B 提供**:
- 监控服务部署
- 日志收集配置

### 与研发 C 协作

**依赖研发 C**:
- Dockerfile
- 构建产物

**向研发 C 提供**:
- 静态资源 CDN
- 部署状态 API

---

## 📅 里程碑和交付物

### Week 2 交付物

- [ ] Docker Compose 完整配置
- [ ] 数据库迁移脚本框架
- [ ] 环境变量配置模板
- [ ] 快速启动文档

### Week 4 交付物

- [ ] K8s 基础配置
- [ ] 多环境 Kustomize 配置
- [ ] 部署脚本
- [ ] K8s 部署文档

### Week 6 交付物

- [ ] CI/CD 流水线
- [ ] 自动化测试集成
- [ ] 自动部署到 K8s
- [ ] CI/CD 文档

### Week 8 交付物

- [ ] 灰度发布配置
- [ ] 一键回滚脚本
- [ ] 自动化回滚
- [ ] 运维手册

---

## ✅ 质量检查清单

### 部署前

- [ ] Docker 镜像构建成功
- [ ] 所有测试通过
- [ ] 环境变量配置完整
- [ ] 健康检查正常

### 发布前

- [ ] 灰度发布配置完成
- [ ] 监控告警配置完成
- [ ] 回滚方案准备就绪
- [ ] 发布文档完整

---

## 📊 关键指标

### 部署效率

- CI 流水线执行时间：< 10 分钟
- CD 部署时间：< 5 分钟
- 回滚时间：< 5 分钟

### 可靠性

- 部署成功率：> 99%
- 服务可用性：> 99.9%
- 零停机发布：100%

---

## 📚 参考资料

- [ZKER-企业级开发规范手册_v1.0.md](../ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md)
- [ZKER-灰度发布策略_v1.0.md](../ZKER-灰度发布策略_v1.0.md)
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)
- [ZKER-故障排查手册_v1.0.md](../ZKER-故障排查手册_v1.0.md)

---

**文档状态**: ✅ 已完成
**最后更新**: 2025-01-01
**责任人**: 研发 D - DevOps 工程师
**评审人**: 技术负责人
