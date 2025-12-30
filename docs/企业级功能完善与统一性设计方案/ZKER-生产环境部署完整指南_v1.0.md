# ZKER企业级SaaS AI智能体平台 - 生产环境部署完整指南

**版本**: v1.0
**日期**: 2025-01-01
**部署目标**: 企业级多租户SaaS平台生产环境
**适用场景**: 公有云部署 / 私有云部署 / 混合云部署

---

## 📋 文档概述

本文档提供 ZKER 企业级 SaaS AI 智能体平台的**完整生产环境部署指南**，涵盖从基础设施准备到系统上线的全流程。

### 部署架构总览

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户接入层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   CDN/WAF    │  │   负载均衡    │  │   API网关     │          │
│  │  (Cloudflare)│  │  (ALB/SLB)   │  │  (Istio)     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                         应用服务层                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  前端静态资源  │  │   后端服务    │  │  Worker异步  │          │
│  │  (Nginx/OSS) │  │  (Go/K8s)    │  │  (NSQ消费)   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         3个实例            6个实例            3个实例             │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                         数据存储层                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │  MySQL   │  │  Redis   │  │    ES    │  │  Milvus  │        │
│  │ (MGR集群) │  │ (哨兵模式)│  │ (3节点)  │  │ (3节点)  │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                        │
│  │  MinIO   │  │   NSQ    │  │  etcd    │                        │
│  │ (分布式)  │  │ (集群)   │  │ (3节点)  │                        │
│  └──────────┘  └──────────┘  └──────────┘                        │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                         监控运维层                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │Prometheus│  │ Grafana  │  │   ELK    │  │  Jaeger  │        │
│  │ (监控采集) │  │ (可视化)  │  │ (日志)   │  │ (追踪)   │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────┐  ┌──────────┐                                        │
│  │ AlertMgr │  │ PagerDuty│                                        │
│  │ (告警)   │  │ (值班)   │                                        │
│  └──────────┘  └──────────┘                                        │
└─────────────────────────────────────────────────────────────────┘
```

### 部署规模建议

| 规模 | 并发用户 | 后端实例 | 数据库配置 | 适用场景 |
|------|---------|---------|-----------|---------|
| **小型** | < 500 | 3个 | 2核8G × 1 | 初创团队测试 |
| **中型** | 500-2000 | 6个 | 4核16G × 3 | 成长企业生产 |
| **大型** | 2000-10000 | 12个 | 8核32G × 5 | 大规模商用 |
| **超大型** | > 10000 | 24+个 | 16核64G × 9 | 企业级SaaS |

---

## 🔧 第一步：基础设施准备

### 1.1 服务器资源配置

#### 最小配置（测试环境）
```yaml
服务器数量: 3台
CPU: 4核
内存: 16GB
磁盘: 200GB SSD
网络: 10Mbps
```

#### 推荐配置（生产环境）
```yaml
服务器数量: 6台（应用服务器3台 + 数据库服务器3台）
CPU: 8核
内存: 32GB
磁盘: 500GB SSD（系统） + 1TB SSD（数据）
网络: 100Mbps
```

#### 企业级配置（高可用）
```yaml
服务器数量: 12台（应用6台 + 数据库3台 + 监控3台）
CPU: 16核
内存: 64GB
磁盘: 500GB SSD（系统） + 2TB NVMe（数据）
网络: 1Gbps
```

### 1.2 操作系统配置

#### 推荐系统
- **Ubuntu 22.04 LTS** / **CentOS 8 Stream** / **AlmaLinux 9**
- 内核版本 ≥ 5.15

#### 系统初始化脚本
```bash
#!/bin/bash
# system-init.sh - 系统初始化

# 1. 更新系统
apt update && apt upgrade -y

# 2. 安装基础工具
apt install -y curl wget vim git net-tools telnet htop iotop

# 3. 配置时区
timedatectl set-timezone Asia/Shanghai

# 4. 配置文件描述符限制
cat >> /etc/security/limits.conf << EOF
* soft nofile 65536
* hard nofile 65536
* soft nproc 65536
* hard nproc 65536
EOF

# 5. 配置内核参数
cat >> /etc/sysctl.conf << EOF
# 网络优化
net.core.somaxconn = 32768
net.ipv4.tcp_max_syn_backlog = 8192
net.core.netdev_max_backlog = 16384

# 连接跟踪
net.netfilter.nf_conntrack_max = 1000000
net.nf_conntrack_max = 1000000

# TCP优化
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_max_tw_buckets = 6000

# 虚拟内存
vm.swappiness = 10
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5
EOF

sysctl -p

# 6. 配置SSH（禁用密码登录，仅密钥）
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
systemctl restart sshd

# 7. 配置防火墙
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw allow 8080/tcp
ufw --force enable

echo "✅ 系统初始化完成！"
```

### 1.3 Docker & Docker Compose 安装

```bash
#!/bin/bash
# install-docker.sh

# 1. 安装Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# 2. 配置Docker守护进程
cat > /etc/docker/daemon.json << EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  },
  "storage-driver": "overlay2",
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com"
  ]
}
EOF

systemctl daemon-reload
systemctl enable docker
systemctl restart docker

# 3. 安装Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# 4. 验证安装
docker --version
docker-compose --version

echo "✅ Docker安装完成！"
```

### 1.4 Kubernetes 集群部署（可选）

如果使用 K8s 进行容器编排：

```bash
#!/bin/bash
# install-k8s.sh

# 1. 安装kubelet、kubeadm、kubectl
apt update
apt install -y apt-transport-https ca-certificates curl
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /' > /etc/apt/sources.list.d/kubernetes.list
apt update
apt install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl

# 2. 初始化Master节点
kubeadm init --pod-network-cidr=10.244.0.0/16

# 3. 配置kubectl
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config

# 4. 安装CNI插件（Calico）
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml

# 5. Worker节点加入（在worker节点执行）
# kubeadm join <master-ip>:6443 --token <token> --discovery-token-ca-cert-hash <hash>

echo "✅ K8s集群部署完成！"
```

---

## 🗄️ 第二步：数据库部署

### 2.1 MySQL 8.4.5 集群部署（MGR模式）

#### Docker Compose 配置
```yaml
# docker-compose-mysql.yml
version: '3.8'

services:
  mysql-primary:
    image: mysql:8.4.5
    container_name: zker-mysql-primary
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_REPLICATION_USER: repl
      MYSQL_REPLICATION_PASSWORD: ${MYSQL_REPL_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql-primary-data:/var/lib/mysql
      - ./mysql/conf/primary.cnf:/etc/mysql/conf.d/custom.cnf
    command:
      - "--server-id=1"
      - "--log-bin=mysql-bin"
      - "--binlog-format=ROW"
      - "--gtid-mode=ON"
      - "--enforce-gtid-consistency=ON"
    networks:
      - zker-network

  mysql-secondary-1:
    image: mysql:8.4.5
    container_name: zker-mysql-secondary-1
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    ports:
      - "3307:3306"
    volumes:
      - mysql-secondary-1-data:/var/lib/mysql
      - ./mysql/conf/secondary.cnf:/etc/mysql/conf.d/custom.cnf
    command:
      - "--server-id=2"
      - "--relay-log=relay-bin"
      - "--read-only=1"
    networks:
      - zker-network

  mysql-secondary-2:
    image: mysql:8.4.5
    container_name: zker-mysql-secondary-2
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    ports:
      - "3308:3306"
    volumes:
      - mysql-secondary-2-data:/var/lib/mysql
      - ./mysql/conf/secondary.cnf:/etc/mysql/conf.d/custom.cnf
    command:
      - "--server-id=3"
      - "--relay-log=relay-bin"
      - "--read-only=1"
    networks:
      - zker-network

  mysql-exporter:
    image: prom/mysqld-exporter:latest
    container_name: zker-mysql-exporter
    restart: always
    environment:
      DATA_SOURCE_NAME: "exporter:${MYSQL_EXPORTER_PASSWORD}@(zker-mysql-primary:3306)/"
    ports:
      - "9104:9104"
    networks:
      - zker-network

volumes:
  mysql-primary-data:
  mysql-secondary-1-data:
  mysql-secondary-2-data:

networks:
  zker-network:
    driver: bridge
```

#### MySQL 配置文件
```ini
# mysql/conf/primary.cnf
[mysqld]
# 连接配置
max_connections = 5000
max_connect_errors = 100000

# InnoDB配置
innodb_buffer_pool_size = 16G
innodb_log_file_size = 2G
innodb_flush_log_at_trx_commit = 1
innodb_flush_method = O_DIRECT

# 查询缓存（MySQL 8.0已移除，使用Redis代替）

# 慢查询日志
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

# 二进制日志
expire_logs_days = 7
max_binlog_size = 1G

# 字符集
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci

# 时区
default-time-zone = '+08:00'

# SQL模式
sql_mode = STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO
```

#### 启动MySQL集群
```bash
# 1. 创建环境变量文件
cat > .env.mysql << EOF
MYSQL_ROOT_PASSWORD=$(openssl rand -base64 32)
MYSQL_REPL_PASSWORD=$(openssl rand -base64 32)
MYSQL_EXPORTER_PASSWORD=$(openssl rand -base64 32)
EOF

# 2. 启动集群
docker-compose -f docker-compose-mysql.yml up -d

# 3. 等待MySQL就绪
sleep 30

# 4. 配置MGR（在MySQL主节点执行）
docker exec -it zker-mysql-primary mysql -uroot -p${MYSQL_ROOT_PASSWORD} << EOF
-- 配置MGR插件
INSTALL PLUGIN group_replication SONAME 'group_replication.so';

-- 创建复制用户
CREATE USER IF NOT EXISTS 'repl'@'%' IDENTIFIED BY '${MYSQL_REPL_PASSWORD}';
GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%';

-- 配置MGR通道
CHANGE MASTER TO MASTER_USER='repl', MASTER_PASSWORD='${MYSQL_REPL_PASSWORD}' FOR CHANNEL 'group_replication_recovery';

-- 启动MGR
START GROUP_REPLICATION;
EOF

# 5. 验证集群状态
docker exec zker-mysql-primary mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "SELECT * FROM performance_schema.replication_group_members;"

echo "✅ MySQL集群部署完成！"
```

### 2.2 数据库初始化

```bash
#!/bin/bash
# init-database.sh

# 1. 创建数据库
docker exec zker-mysql-primary mysql -uroot -p${MYSQL_ROOT_PASSWORD} << EOF
CREATE DATABASE IF NOT EXISTS zker CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS zker_logs CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

# 2. 执行初始化SQL脚本
cd backend
make sql_init

# 3. 验证表结构
docker exec zker-mysql-primary mysql -uroot -p${MYSQL_ROOT_PASSWORD} zker -e "SHOW TABLES;"

echo "✅ 数据库初始化完成！"
```

---

## 🔌 第三步：中间件部署

### 3.1 Redis 集群部署（哨兵模式）

```yaml
# docker-compose-redis.yml
version: '3.8'

services:
  redis-master:
    image: redis:7.0-alpine
    container_name: zker-redis-master
    restart: always
    command: redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}
    ports:
      - "6379:6379"
    volumes:
      - redis-master-data:/data
    networks:
      - zker-network

  redis-slave-1:
    image: redis:7.0-alpine
    container_name: zker-redis-slave-1
    restart: always
    command: redis-server --slaveof zker-redis-master 6379 --masterauth ${REDIS_PASSWORD} --requirepass ${REDIS_PASSWORD}
    ports:
      - "6380:6379"
    volumes:
      - redis-slave-1-data:/data
    networks:
      - zker-network

  redis-slave-2:
    image: redis:7.0-alpine
    container_name: zker-redis-slave-2
    restart: always
    command: redis-server --slaveof zker-redis-master 6379 --masterauth ${REDIS_PASSWORD} --requirepass ${REDIS_PASSWORD}
    ports:
      - "6381:6379"
    volumes:
      - redis-slave-2-data:/data
    networks:
      - zker-network

  redis-sentinel-1:
    image: redis:7.0-alpine
    container_name: zker-redis-sentinel-1
    restart: always
    command: redis-sentinel /etc/redis/sentinel.conf
    ports:
      - "26379:26379"
    volumes:
      - ./redis/sentinel.conf:/etc/redis/sentinel.conf
    networks:
      - zker-network

  redis-exporter:
    image: oliver006/redis_exporter:latest
    container_name: zker-redis-exporter
    restart: always
    environment:
      REDIS_ADDR: redis://zker-redis-master:6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
    ports:
      - "9121:9121"
    networks:
      - zker-network

volumes:
  redis-master-data:
  redis-slave-1-data:
  redis-slave-2-data:

networks:
  zker-network:
    driver: bridge
```

#### Redis Sentinel 配置
```ini
# redis/sentinel.conf
port 26379
sentinel monitor mymaster zker-redis-master 6379 2
sentinel auth-pass mymaster ${REDIS_PASSWORD}
sentinel down-after-milliseconds mymaster 5000
sentinel parallel-syncs mymaster 1
sentinel failover-timeout mymaster 10000
```

### 3.2 Elasticsearch 8.18.0 集群部署

```yaml
# docker-compose-es.yml
version: '3.8'

services:
  es-node-1:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.18.0
    container_name: zker-es-node-1
    restart: always
    environment:
      - node.name=es-node-1
      - cluster.name=zker-es-cluster
      - discovery.seed_hosts=es-node-2,es-node-3
      - cluster.initial_master_nodes=es-node-1,es-node-2,es-node-3
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms4g -Xmx4g"
      - xpack.security.enabled=true
      - ELASTIC_PASSWORD=${ES_PASSWORD}
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - es-node-1-data:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
    networks:
      - zker-network

  es-node-2:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.18.0
    container_name: zker-es-node-2
    restart: always
    environment:
      - node.name=es-node-2
      - cluster.name=zker-es-cluster
      - discovery.seed_hosts=es-node-1,es-node-3
      - cluster.initial_master_nodes=es-node-1,es-node-2,es-node-3
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms4g -Xmx4g"
      - xpack.security.enabled=true
      - ELASTIC_PASSWORD=${ES_PASSWORD}
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - es-node-2-data:/usr/share/elasticsearch/data
    ports:
      - "9201:9200"
    networks:
      - zker-network

  es-node-3:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.18.0
    container_name: zker-es-node-3
    restart: always
    environment:
      - node.name=es-node-3
      - cluster.name=zker-es-cluster
      - discovery.seed_hosts=es-node-1,es-node-2
      - cluster.initial_master_nodes=es-node-1,es-node-2,es-node-3
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms4g -Xmx4g"
      - xpack.security.enabled=true
      - ELASTIC_PASSWORD=${ES_PASSWORD}
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - es-node-3-data:/usr/share/elasticsearch/data
    ports:
      - "9202:9200"
    networks:
      - zker-network

  kibana:
    image: docker.elastic.co/kibana/kibana:8.18.0
    container_name: zker-kibana
    restart: always
    environment:
      - ELASTICSEARCH_HOSTS=http://es-node-1:9200
      - ELASTICSEARCH_USERNAME=kibana_system
      - ELASTICSEARCH_PASSWORD=${ES_PASSWORD}
    ports:
      - "5601:5601"
    networks:
      - zker-network
    depends_on:
      - es-node-1

volumes:
  es-node-1-data:
  es-node-2-data:
  es-node-3-data:

networks:
  zker-network:
    driver: bridge
```

### 3.3 Milvus 2.5.10 向量数据库部署

```yaml
# docker-compose-milvus.yml
version: '3.8'

services:
  etcd:
    image: quay.io/coreos/etcd:v3.5.0
    container_name: zker-etcd
    restart: always
    environment:
      - ETCD_AUTO_COMPACTION_MODE=revision
      - ETCD_AUTO_COMPACTION_RETENTION=1000
      - ETCD_QUOTA_BACKEND_BYTES=4294967296
    volumes:
      - etcd-data:/etcd
    command: etcd -advertise-client-urls=http://127.0.0.1:2379 -listen-client-urls http://0.0.0.0:2379
    networks:
      - zker-network

  minio:
    image: minio/minio:latest
    container_name: zker-minio
    restart: always
    environment:
      MINIO_ACCESS_KEY: ${MINIO_ACCESS_KEY}
      MINIO_SECRET_KEY: ${MINIO_SECRET_KEY}
    volumes:
      - minio-data:/minio_data
    command: minio server /minio_data
    ports:
      - "9000:9000"
      - "9001:9001"
    networks:
      - zker-network

  pulsar:
    image: apachepulsar/pulsar:2.10.0
    container_name: zker-pulsar
    restart: always
    command: bin/pulsar standalone
    networks:
      - zker-network

  milvus-standalone:
    image: milvusdb/milvus:v2.5.10-gpu
    container_name: zker-milvus
    restart: always
    environment:
      ETCD_ENDPOINTS: etcd:2379
      MINIO_ADDRESS: minio:9000
    volumes:
      - milvus-data:/var/lib/milvus
    ports:
      - "19530:19530"
    networks:
      - zker-network
    depends_on:
      - etcd
      - minio
      - pulsar

volumes:
  etcd-data:
  minio-data:
  milvus-data:

networks:
  zker-network:
    driver: bridge
```

### 3.4 启动所有中间件

```bash
#!/bin/bash
# start-middleware.sh

# 1. Redis
docker-compose -f docker-compose-redis.yml up -d

# 2. Elasticsearch
docker-compose -f docker-compose-es.yml up -d

# 3. Milvus
docker-compose -f docker-compose-milvus.yml up -d

# 4. 等待所有服务就绪
sleep 60

# 5. 验证服务状态
docker ps | grep -E "redis|es|milvus"

echo "✅ 所有中间件部署完成！"
```

---

## 🚀 第四步：后端服务部署

### 4.1 构建后端服务

```bash
#!/bin/bash
# build-backend.sh

cd backend

# 1. 编译Go二进制文件
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zker-api ./api

# 2. 构建Docker镜像
docker build -t zker-backend:v1.0 -f docker/Dockerfile .

# 3. 推送到镜像仓库（可选）
# docker tag zker-backend:v1.0 registry.example.com/zker-backend:v1.0
# docker push registry.example.com/zker-backend:v1.0

echo "✅ 后端服务构建完成！"
```

### 4.2 后端服务 Docker Compose 配置

```yaml
# docker-compose-backend.yml
version: '3.8'

services:
  zker-api-1:
    image: zker-backend:v1.0
    container_name: zker-api-1
    restart: always
    environment:
      # 数据库配置
      DB_HOST: zker-mysql-primary
      DB_PORT: 3306
      DB_USER: root
      DB_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      DB_NAME: zker

      # Redis配置
      REDIS_HOST: zker-redis-master
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}

      # Elasticsearch配置
      ES_HOSTS: http://zker-es-node-1:9200,zker-es-node-2:9200,zker-es-node-3:9200
      ES_USERNAME: elastic
      ES_PASSWORD: ${ES_PASSWORD}

      # Milvus配置
      MILVUS_HOST: zker-milvus
      MILVUS_PORT: 19530

      # 业务配置
      ENVIRONMENT: production
      LOG_LEVEL: info
      TZ: Asia/Shanghai
    ports:
      - "8081:8080"
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '1'
          memory: 2G

  zker-api-2:
    image: zker-backend:v1.0
    container_name: zker-api-2
    restart: always
    environment:
      # 同上...
    ports:
      - "8082:8080"
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G

  zker-api-3:
    image: zker-backend:v1.0
    container_name: zker-api-3
    restart: always
    environment:
      # 同上...
    ports:
      - "8083:8080"
    volumes:
      - ./logs:/app/logs
      - ./uploads:/app/uploads
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G

  zker-worker:
    image: zker-backend:v1.0
    container_name: zker-worker
    restart: always
    command: ./zker-worker
    environment:
      # 同API配置...
    volumes:
      - ./logs:/app/logs
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '4'
          memory: 8G

networks:
  zker-network:
    external: true
```

### 4.3 启动后端服务

```bash
#!/bin/bash
# start-backend.sh

# 1. 启动后端服务
docker-compose -f docker-compose-backend.yml up -d

# 2. 等待服务就绪
sleep 30

# 3. 健康检查
for i in {1..2}; do
  curl -f http://localhost:808$i/health || echo "API-$i 健康检查失败"
done

# 4. 验证数据库连接
docker exec zker-api-1 ./zker-api db:ping

echo "✅ 后端服务部署完成！"
```

---

## 🎨 第五步：前端部署

### 5.1 构建前端静态资源

```bash
#!/bin/bash
# build-frontend.sh

cd frontend

# 1. 安装依赖
rush update

# 2. 构建生产版本
rush build

# 3. 收集构建产物
mkdir -p dist
cp -r apps/coze-studio/dist/* dist/

# 4. 构建Docker镜像（Nginx）
docker build -t zker-frontend:v1.0 -f docker/Dockerfile .

echo "✅ 前端构建完成！"
```

### 5.2 Nginx 配置

```nginx
# nginx.conf
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 10240;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for" '
                    'rt=$request_time uct="$upstream_connect_time" '
                    'uht="$upstream_header_time" urt="$upstream_response_time"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 100M;

    # Gzip压缩
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript
               application/json application/javascript application/xml+rss
               application/rss+xml font/truetype font/opentype
               application/vnd.ms-fontobject image/svg+xml;

    # 前端服务器
    upstream frontend_servers {
        least_conn;
        server zker-frontend-1:80;
        server zker-frontend-2:80;
        server zker-frontend-3:80;
    }

    # 后端API服务器
    upstream api_servers {
        least_conn;
        server zker-api-1:8080;
        server zker-api-2:8080;
        server zker-api-3:8080;
    }

    server {
        listen 80;
        server_name www.zker.com;

        # 强制HTTPS（生产环境建议）
        # return 301 https://$server_name$request_uri;

        # 前端静态资源
        location / {
            root /usr/share/nginx/html;
            index index.html;
            try_files $uri $uri/ /index.html;

            # 缓存配置
            location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
                expires 30d;
                add_header Cache-Control "public, immutable";
            }
        }

        # API代理
        location /api/ {
            proxy_pass http://api_servers;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # 超时配置
            proxy_connect_timeout 60s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;

            # 缓冲配置
            proxy_buffering on;
            proxy_buffer_size 4k;
            proxy_buffers 8 4k;
        }

        # WebSocket支持
        location /ws/ {
            proxy_pass http://api_servers;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_read_timeout 86400;
        }

        # 健康检查
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }

    # HTTPS配置（生产环境启用）
    # server {
    #     listen 443 ssl http2;
    #     server_name www.zker.com;
    #
    #     ssl_certificate /etc/nginx/ssl/zker.com.crt;
    #     ssl_certificate_key /etc/nginx/ssl/zker.com.key;
    #     ssl_protocols TLSv1.2 TLSv1.3;
    #     ssl_ciphers HIGH:!aNULL:!MD5;
    #     ssl_prefer_server_ciphers on;
    #
    #     # 其他配置同上...
    # }
}
```

### 5.3 前端 Docker Compose 配置

```yaml
# docker-compose-frontend.yml
version: '3.8'

services:
  zker-frontend-1:
    image: zker-frontend:v1.0
    container_name: zker-frontend-1
    restart: always
    ports:
      - "3001:80"
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G

  zker-frontend-2:
    image: zker-frontend:v1.0
    container_name: zker-frontend-2
    restart: always
    ports:
      - "3002:80"
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G

  zker-frontend-3:
    image: zker-frontend:v1.0
    container_name: zker-frontend-3
    restart: always
    ports:
      - "3003:80"
    networks:
      - zker-network
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G

  nginx-lb:
    image: nginx:alpine
    container_name: zker-nginx-lb
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    networks:
      - zker-network
    depends_on:
      - zker-frontend-1
      - zker-frontend-2
      - zker-frontend-3

networks:
  zker-network:
    external: true
```

### 5.4 启动前端服务

```bash
#!/bin/bash
# start-frontend.sh

# 1. 启动前端服务
docker-compose -f docker-compose-frontend.yml up -d

# 2. 验证服务
curl -I http://localhost/

echo "✅ 前端服务部署完成！"
```

---

## 📊 第六步：监控系统部署

### 6.1 Prometheus + Grafana 部署

```yaml
# docker-compose-monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:v2.45.0
    container_name: zker-prometheus
    restart: always
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./monitoring/prometheus/rules:/etc/prometheus/rules
      - prometheus-data:/prometheus
    networks:
      - zker-network

  grafana:
    image: grafana/grafana:10.0.0
    container_name: zker-grafana
    restart: always
    environment:
      GF_SECURITY_ADMIN_USER: admin
      GF_SECURITY_ADMIN_PASSWORD: ${GRAFANA_ADMIN_PASSWORD}
      GF_INSTALL_PLUGINS: grafana-piechart-panel
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
      - ./monitoring/grafana/provisioning:/etc/grafana/provisioning
      - ./monitoring/grafana/dashboards:/var/lib/grafana/dashboards
    networks:
      - zker-network
    depends_on:
      - prometheus

  alertmanager:
    image: prom/alertmanager:v0.26.0
    container_name: zker-alertmanager
    restart: always
    command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/alertmanager'
    ports:
      - "9093:9093"
    volumes:
      - ./monitoring/alertmanager/alertmanager.yml:/etc/alertmanager/alertmanager.yml
      - alertmanager-data:/alertmanager
    networks:
      - zker-network

  node-exporter:
    image: prom/node-exporter:latest
    container_name: zker-node-exporter
    restart: always
    command:
      - '--path.procfs=/host/proc'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    ports:
      - "9100:9100"
    networks:
      - zker-network

  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    container_name: zker-cadvisor
    restart: always
    command:
      - '--housekeeping_interval=30s'
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:ro
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
    ports:
      - "8080:8080"
    networks:
      - zker-network

volumes:
  prometheus-data:
  grafana-data:
  alertmanager-data:

networks:
  zker-network:
    external: true
```

### 6.2 Prometheus 配置

```yaml
# monitoring/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: 'zker-prod'
    environment: 'production'

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

rule_files:
  - "/etc/prometheus/rules/*.yml"

scrape_configs:
  # API服务监控
  - job_name: 'zker-api'
    static_configs:
      - targets:
          - zker-api-1:8080
          - zker-api-2:8080
          - zker-api-3:8080
    metrics_path: '/metrics'

  # MySQL监控
  - job_name: 'mysql'
    static_configs:
      - targets:
          - mysql-exporter:9104

  # Redis监控
  - job_name: 'redis'
    static_configs:
      - targets:
          - redis-exporter:9121

  # Node Exporter
  - job_name: 'node'
    static_configs:
      - targets:
          - node-exporter:9100

  # cAdvisor
  - job_name: 'cadvisor'
    static_configs:
      - targets:
          - cadvisor:8080
```

### 6.3 告警规则配置

```yaml
# monitoring/prometheus/rules/alerts.yml
groups:
  - name: api_alerts
    interval: 30s
    rules:
      # API可用性告警
      - alert: APIServiceDown
        expr: up{job="zker-api"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "API服务不可用"
          description: "{{ $labels.instance }} API服务已宕机超过1分钟"

      # API高延迟告警
      - alert: APIHighLatency
        expr: histogram_quantile(0.95, http_request_duration_seconds_bucket) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API响应延迟过高"
          description: "P95延迟 {{ $value }}s 超过1s阈值"

      # API错误率告警
      - alert: APIHighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API错误率过高"
          description: "5xx错误率 {{ $value | humanizePercentage }} 超过5%"

  - name: database_alerts
    interval: 30s
    rules:
      # MySQL主从延迟告警
      - alert: MySQLReplicationLag
        expr: mysql_slave_lag_seconds > 60
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MySQL主从延迟过高"
          description: "主从延迟 {{ $value }}s 超过60s"

      # MySQL连接数告警
      - alert: MySQLHighConnections
        expr: mysql_global_status_threads_connected / mysql_global_variables_max_connections > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MySQL连接数过高"
          description: "连接数使用率 {{ $value | humanizePercentage }} 超过80%"

  - name: system_alerts
    interval: 30s
    rules:
      # CPU使用率告警
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "CPU使用率过高"
          description: "节点 {{ $labels.instance }} CPU使用率 {{ $value }}% 超过80%"

      # 内存使用率告警
      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "内存使用率过高"
          description: "节点 {{ $labels.instance }} 内存使用率 {{ $value }}% 超过85%"

      # 磁盘使用率告警
      - alert: HighDiskUsage
        expr: (1 - (node_filesystem_avail_bytes{fstype!~"tmpfs|fuse.*"} / node_filesystem_size_bytes)) * 100 > 85
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "磁盘使用率过高"
          description: "挂载点 {{ $labels.mountpoint }} 使用率 {{ $value }}% 超过85%"
```

### 6.4 启动监控系统

```bash
#!/bin/bash
# start-monitoring.sh

# 1. 启动监控系统
docker-compose -f docker-compose-monitoring.yml up -d

# 2. 等待服务就绪
sleep 30

# 3. 验证Prometheus
curl -f http://localhost:9090/-/healthy

# 4. 验证Grafana
curl -f http://localhost:3000/api/health

echo "✅ 监控系统部署完成！"
echo "📊 Prometheus: http://localhost:9090"
echo "📊 Grafana: http://localhost:3000 (admin/${GRAFANA_ADMIN_PASSWORD})"
```

---

## 🧪 第七步：部署验证

### 7.1 健康检查脚本

```bash
#!/bin/bash
# health-check.sh

echo "========================================="
echo "     ZKER生产环境健康检查"
echo "========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 检查函数
check_service() {
    local name=$1
    local url=$2
    local expected_code=${3:-200}

    echo -n "检查 $name ... "
    status_code=$(curl -s -o /dev/null -w "%{http_code}" $url)

    if [ "$status_code" == "$expected_code" ]; then
        echo -e "${GREEN}✓ 正常${NC} ($status_code)"
        return 0
    else
        echo -e "${RED}✗ 异常${NC} ($status_code)"
        return 1
    fi
}

# 1. 前端服务
check_service "前端服务" "http://localhost/health"

# 2. API服务
for i in {1..3}; do
    check_service "API服务-$i" "http://localhost:808$i/health"
done

# 3. 数据库
echo -n "检查 MySQL ... "
docker exec zker-mysql-primary mysqladmin ping -h localhost -uroot -p${MYSQL_ROOT_PASSWORD} > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 正常${NC}"
else
    echo -e "${RED}✗ 异常${NC}"
fi

# 4. Redis
echo -n "检查 Redis ... "
docker exec zker-redis-master redis-cli -a ${REDIS_PASSWORD} ping > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 正常${NC}"
else
    echo -e "${RED}✗ 异常${NC}"
fi

# 5. Elasticsearch
echo -n "检查 Elasticsearch ... "
curl -s -u elastic:${ES_PASSWORD} http://localhost:9200/_cluster/health | grep -q '"status":"green"\|"status":"yellow"'
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 正常${NC}"
else
    echo -e "${RED}✗ 异常${NC}"
fi

# 6. Prometheus
check_service "Prometheus" "http://localhost:9090/-/healthy" 200

# 7. Grafana
check_service "Grafana" "http://localhost:3000/api/health" 200

echo ""
echo "========================================="
echo "     健康检查完成"
echo "========================================="
```

### 7.2 功能测试

```bash
#!/bin/bash
# functional-test.sh

echo "开始功能测试..."

# 1. 用户注册
echo "1. 测试用户注册..."
curl -X POST http://localhost:8081/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@zker.com",
    "password": "Test123456!",
    "username": "testuser"
  }' | jq .

# 2. 用户登录
echo "2. 测试用户登录..."
TOKEN=$(curl -X POST http://localhost:8081/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@zker.com",
    "password": "Test123456!"
  }' | jq -r '.data.token')

echo "获取Token: $TOKEN"

# 3. 创建Bot
echo "3. 测试创建Bot..."
curl -X POST http://localhost:8081/api/v1/bots \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试Bot",
    "description": "这是一个测试Bot",
    "avatar_url": "https://example.com/avatar.png"
  }' | jq .

# 4. 获取Bot列表
echo "4. 测试获取Bot列表..."
curl -X GET http://localhost:8081/api/v1/bots \
  -H "Authorization: Bearer $TOKEN" | jq .

echo "✅ 功能测试完成！"
```

---

## 🔐 第八步：安全加固

### 8.1 SSL/TLS证书配置

```bash
#!/bin/bash
# setup-ssl.sh

DOMAIN="www.zker.com"

# 1. 生成自签名证书（测试环境）
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/${DOMAIN}.key \
  -out /etc/nginx/ssl/${DOMAIN}.crt \
  -subj "/C=CN/ST=Beijing/L=Beijing/O=ZKER/CN=${DOMAIN}"

# 2. 生产环境使用Let's Encrypt
# certbot certonly --webroot -w /var/www/html -d ${DOMAIN} -d www.${DOMAIN}

echo "✅ SSL证书配置完成！"
```

### 8.2 防火墙配置

```bash
#!/bin/bash
# setup-firewall.sh

# 清空规则
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X

# 默认策略
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# 允许本地回环
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# 允许已建立的连接
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 允许SSH
iptables -A INPUT -p tcp --dport 22 -j ACCEPT

# 允许HTTP/HTTPS
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# 保存规则
iptables-save > /etc/iptables/rules.v4

echo "✅ 防火墙配置完成！"
```

### 8.3 敏感信息管理

```bash
# 使用环境变量管理敏感信息
cat > .env.production << EOF
# 数据库
MYSQL_ROOT_PASSWORD=$(openssl rand -base64 32)
MYSQL_REPL_PASSWORD=$(openssl rand -base64 32)

# Redis
REDIS_PASSWORD=$(openssl rand -base64 32)

# Elasticsearch
ES_PASSWORD=$(openssl rand -base64 32)

# MinIO
MINIO_ACCESS_KEY=admin
MINIO_SECRET_KEY=$(openssl rand -base64 32)

# JWT
JWT_SECRET=$(openssl rand -base64 64)

# 加密密钥
ENCRYPTION_KEY=$(openssl rand -base64 32)

# Grafana
GRAFANA_ADMIN_PASSWORD=$(openssl rand -base64 16)

# 第三方服务
OPENAI_API_KEY=sk-xxx
ANTHROPIC_API_KEY=sk-ant-xxx
EOF

# 设置文件权限
chmod 600 .env.production
```

---

## 📦 第九步：备份恢复策略

### 9.1 数据库备份脚本

```bash
#!/bin/bash
# backup-database.sh

BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/zker_${DATE}.sql.gz"

# 创建备份目录
mkdir -p ${BACKUP_DIR}

# 备份数据库
docker exec zker-mysql-primary mysqldump \
  -uroot \
  -p${MYSQL_ROOT_PASSWORD} \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --all-databases | gzip > ${BACKUP_FILE}

# 保留最近7天的备份
find ${BACKUP_DIR} -name "zker_*.sql.gz" -mtime +7 -delete

# 上传到对象存储（可选）
# aws s3 cp ${BACKUP_FILE} s3://zker-backups/mysql/

echo "✅ 数据库备份完成: ${BACKUP_FILE}"
```

### 9.2 定时备份任务

```bash
# 添加到crontab
crontab -e

# 每天凌晨2点执行数据库备份
0 2 * * * /opt/scripts/backup-database.sh >> /var/log/backup.log 2>&1

# 每小时执行Redis备份
0 * * * * docker exec zker-redis-master redis-cli --rdb /data/dump_$(date +\%Y\%m\%d_\%H\%M).rdb

# 每天凌晨3点执行文件备份
0 3 * * * tar -czf /backup/uploads_$(date +\%Y\%m\%d).tar.gz /app/uploads
```

### 9.3 恢复脚本

```bash
#!/bin/bash
# restore-database.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

echo "开始恢复数据库: ${BACKUP_FILE}"

# 停止应用服务
docker-compose -f docker-compose-backend.yml stop

# 恢复数据库
gunzip < ${BACKUP_FILE} | docker exec -i zker-mysql-primary mysql -uroot -p${MYSQL_ROOT_PASSWORD}

# 启动应用服务
docker-compose -f docker-compose-backend.yml start

echo "✅ 数据库恢复完成！"
```

---

## 🚨 第十步：灰度发布配置

### 10.1 Istio 灰度发布配置

```yaml
# istio/vs-zker-api.yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: zker-api
spec:
  hosts:
    - "api.zker.com"
  gateways:
    - zker-gateway
  http:
    - match:
        - headers:
            x-gray-release:
              exact: "v2"
      route:
        - destination:
            host: zker-api
            subset: v2
          weight: 100
    - route:
        - destination:
            host: zker-api
            subset: v1
          weight: 90
        - destination:
            host: zker-api
            subset: v2
          weight: 10
```

### 10.2 灰度发布流程

```bash
#!/bin/bash
# gray-release.sh

VERSION=$1  # v2

echo "开始灰度发布 ${VERSION}..."

# 1. 部署新版本
docker-compose -f docker-compose-backend-${VERSION}.yml up -d

# 2. 配置灰度规则（10%流量）
kubectl apply -f istio/vs-zker-api-10.yaml

# 3. 观察指标
echo "观察新版本指标..."
# 检查错误率、延迟等指标

# 4. 逐步增加流量（根据观察结果）
# kubectl apply -f istio/vs-zker-api-50.yaml
# kubectl apply -f istio/vs-zker-api-100.yaml

# 5. 全量发布后删除旧版本
# docker-compose -f docker-compose-backend-v1.yml down

echo "✅ 灰度发布完成！"
```

---

## 📋 第十一步：运维手册

### 11.1 日常运维命令

```bash
# 查看服务状态
docker-compose -f docker-compose-backend.yml ps
docker-compose -f docker-compose-frontend.yml ps

# 查看日志
docker logs zker-api-1 --tail 100 -f
docker logs zker-frontend-1 --tail 100 -f

# 重启服务
docker-compose -f docker-compose-backend.yml restart zker-api-1

# 扩容服务
docker-compose -f docker-compose-backend.yml up -d --scale zker-api=5

# 查看资源使用
docker stats

# 进入容器调试
docker exec -it zker-api-1 sh
```

### 11.2 故障处理流程

#### 场景1：API服务无响应
```bash
# 1. 检查服务状态
docker ps | grep zker-api

# 2. 查看日志
docker logs zker-api-1 --tail 500

# 3. 检查数据库连接
docker exec zker-api-1 ./zker-api db:ping

# 4. 重启服务
docker-compose -f docker-compose-backend.yml restart

# 5. 如果仍无法解决，执行回滚
# ./rollback.sh v1
```

#### 场景2：数据库主从延迟
```bash
# 1. 检查主从状态
docker exec zker-mysql-primary mysql -e "SHOW SLAVE STATUS\G"

# 2. 检查慢查询
docker exec zker-mysql-primary mysql -e "SELECT * FROM mysql.slow_log ORDER BY start_time DESC LIMIT 10;"

# 3. 优化慢查询
# ...

# 4. 如果延迟过大，考虑提升从库为主库
# ...
```

#### 场景3：Redis连接数耗尽
```bash
# 1. 检查连接数
docker exec zker-redis-master redis-cli -a ${REDIS_PASSWORD} INFO clients

# 2. 查看慢查询
docker exec zker-redis-master redis-cli -a ${REDIS_PASSWORD} SLOWLOG GET 10

# 3. 清理过期键
docker exec zker-redis-master redis-cli -a ${REDIS_PASSWORD} --scan --pattern "session:*" | xargs docker exec zker-redis-master redis-cli -a ${REDIS_PASSWORD} DEL

# 4. 如果仍然无法解决，扩容Redis集群
# ...
```

### 11.3 监控告警处理

| 告警级别 | 响应时间 | 处理方式 |
|---------|---------|---------|
| **P0 - 严重** | 5分钟内 | 电话 + IM + 立即处理 |
| **P1 - 高** | 15分钟内 | IM + 1小时内处理 |
| **P2 - 中** | 1小时内 | 工单 + 4小时内处理 |
| **P3 - 低** | 1天内 | 工单 + 3天内处理 |

---

## 📊 附录

### A. 端口清单

| 服务 | 端口 | 说明 |
|------|------|------|
| Nginx | 80, 443 | 前端/网关 |
| API服务 | 8081-8083 | 后端API |
| MySQL | 3306-3308 | 数据库 |
| Redis | 6379-6381 | 缓存 |
| Elasticsearch | 9200-9202 | 搜索 |
| Milvus | 19530 | 向量数据库 |
| MinIO | 9000, 9001 | 对象存储 |
| Prometheus | 9090 | 监控采集 |
| Grafana | 3000 | 可视化 |
| AlertManager | 9093 | 告警 |

### B. 目录结构

```
/opt/zker/
├── backend/              # 后端代码
├── frontend/            # 前端代码
├── docker/              # Docker配置
│   ├── mysql/
│   ├── redis/
│   ├── es/
│   └── nginx/
├── monitoring/          # 监控配置
│   ├── prometheus/
│   ├── grafana/
│   └── alertmanager/
├── scripts/             # 运维脚本
│   ├── backup/
│   ├── deploy/
│   └── monitoring/
├── logs/                # 日志目录
├── uploads/             # 上传文件
└── backup/              # 备份目录
    ├── mysql/
    ├── redis/
    └── uploads/
```

### C. 环境变量清单

```bash
# .env.production 示例
ENVIRONMENT=production
TZ=Asia/Shanghai

# 数据库
DB_HOST=zker-mysql-primary
DB_PORT=3306
DB_USER=root
DB_PASSWORD=***
DB_NAME=zker

# Redis
REDIS_HOST=zker-redis-master
REDIS_PORT=6379
REDIS_PASSWORD=***

# Elasticsearch
ES_HOSTS=http://zker-es-node-1:9200
ES_USERNAME=elastic
ES_PASSWORD=***

# Milvus
MILVUS_HOST=zker-milvus
MILVUS_PORT=19530

# JWT
JWT_SECRET=***
JWT_EXPIRATION=24h

# 加密
ENCRYPTION_KEY=***

# 第三方服务
OPENAI_API_KEY=sk-xxx
ANTHROPIC_API_KEY=sk-ant-xxx

# 监控
GRAFANA_ADMIN_PASSWORD=***
```

---

## ✅ 部署检查清单

### 部署前检查

- [ ] 服务器资源配置满足要求
- [ ] 操作系统版本符合要求
- [ ] Docker & Docker Compose 已安装
- [ ] 网络配置正确（端口开放、DNS解析）
- [ ] 防火墙规则已配置
- [ ] SSL证书已准备（生产环境）
- [ ] 所有密码/密钥已生成并妥善保管

### 部署中检查

- [ ] MySQL集群主从复制正常
- [ ] Redis哨兵模式正常
- [ ] Elasticsearch集群健康状态为green/yellow
- [ ] Milvus可连接
- [ ] 后端服务健康检查通过
- [ ] 前端页面可访问
- [ ] 监控系统正常运行

### 部署后检查

- [ ] 所有服务健康检查通过
- [ ] 功能测试全部通过
- [ ] 监控大盘数据正常
- [ ] 告警规则已配置并测试
- [ ] 备份任务已配置
- [ ] 日志正常收集
- [ ] 性能基线已建立

---

## 🎯 总结

本文档提供了 ZKER 企业级 SaaS AI 智能体平台的**完整生产环境部署指南**，涵盖了：

✅ **基础设施准备** - 系统配置、Docker/K8s安装
✅ **数据库部署** - MySQL MGR集群、初始化
✅ **中间件部署** - Redis、ES、Milvus、MinIO、NSQ
✅ **后端部署** - Go服务、Worker异步任务
✅ **前端部署** - Nginx负载均衡、静态资源
✅ **监控部署** - Prometheus+Grafana+AlertManager
✅ **安全加固** - SSL/TLS、防火墙、敏感信息管理
✅ **备份恢复** - 数据库备份、文件备份、恢复流程
✅ **灰度发布** - Istio流量管理、灰度策略
✅ **运维手册** - 日常运维、故障处理、告警响应

按照本指南执行后，您将获得一个**企业级、高可用、可监控**的 ZKER SaaS 平台生产环境。

---

**文档版本**: v1.0
**最后更新**: 2025-01-01
**维护人**: DevOps团队
**联系方式**: ops@zker.com

**相关文档**:
- [ZKER-企业级SaaS平台最终验收报告_v1.0.md](./ZKER-企业级SaaS平台最终验收报告_v1.0.md)
- [ZKER-持续改进路线图_v1.0.md](./ZKER-持续改进路线图_v1.0.md)
- [ZKER-故障排查手册_v1.0.md](./ZKER-故障排查手册_v1.0.md)
