# ZKER 监控系统部署使用指南

**版本**: v1.0.0
**最后更新**: 2025-01-01
**维护者**: ZKER DevOps Team

---

## 📋 目录

- [快速开始](#快速开始)
- [前置条件](#前置条件)
- [一键部署](#一键部署)
- [单独部署](#单独部署)
- [验证部署](#验证部署)
- [配置管理](#配置管理)
- [常见问题](#常见问题)
- [故障排查](#故障排查)

---

## 🚀 快速开始

### 最快5分钟部署

```bash
# 1. 进入项目目录
cd /path/to/coze-studio

# 2. 一键部署所有组件
sudo ./deploy/monitoring/scripts/deploy-all.sh

# 3. 等待部署完成，访问服务
# Prometheus: http://localhost:9090
# Grafana: http://localhost:3000 (admin/admin)
# AlertManager: http://localhost:9093
```

---

## 📦 前置条件

### 必需组件

| 组件 | 版本要求 | 检查命令 |
|------|---------|---------|
| **Docker** | ≥ 20.10 | `docker --version` |
| **Docker Compose** | ≥ 2.0 | `docker compose version` |
| **Bash** | ≥ 4.0 | `bash --version` |

### 安装Docker（如未安装）

**Ubuntu/Debian**:
```bash
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
```

**CentOS/RHEL**:
```bash
sudo yum install -y docker
sudo systemctl start docker
sudo systemctl enable docker
```

**macOS**:
```bash
brew install docker docker-compose
```

---

## 🎯 一键部署

### 完整部署（所有组件）

```bash
# 方式一：使用默认配置
sudo ./deploy/monitoring/scripts/deploy-all.sh

# 方式二：自定义端口
export PROMETHEUS_PORT=9090
export GRAFANA_PORT=3000
export ALERTMANAGER_PORT=9093
sudo ./deploy/monitoring/scripts/deploy-all.sh

# 方式三：自定义Grafana凭据
export GRAFANA_ADMIN_USER=admin
export GRAFANA_ADMIN_PASSWORD=your_secure_password
sudo ./deploy/monitoring/scripts/deploy-all.sh
```

### 部署流程

一键部署脚本将自动执行以下步骤：

1. ✅ 检查前置条件（Docker、网络）
2. ✅ 安装Prometheus（配置 + 启动）
3. ✅ 安装Grafana（Provisioning + 仪表盘）
4. ✅ 安装AlertManager（告警路由 + 配置）
5. ✅ 安装Exporters（Node、MySQL、Redis）
6. ✅ 验证所有服务启动
7. ✅ 显示访问信息

### 预期输出

```
==========================================
  ZKER 企业级监控系统 - 一键部署
==========================================
版本: v1.0.0
日期: 2025-01-01 12:00:00
==========================================

[STEP] 检查前置条件...
[INFO] ✅ Docker已安装
[INFO] ✅ Docker Compose已安装
[INFO] ✅ Docker网络已就绪

[STEP] 安装Prometheus...
[INFO] Prometheus容器启动成功
[INFO] ✅ Prometheus容器运行正常

[STEP] 安装Grafana...
[INFO] Grafana容器启动成功
[INFO] ✅ Grafana容器运行正常

[STEP] 安装AlertManager...
[INFO] AlertManager容器启动成功
[INFO] ✅ AlertManager容器运行正常

[STEP] 安装Exporters...
[INFO] Node Exporter安装完成
[INFO] MySQL Exporter安装完成
[INFO] Redis Exporter安装完成

[STEP] 验证部署...
[INFO] ✅ 所有服务已就绪

==========================================
  部署完成！
==========================================
```

---

## 🔧 单独部署

### 仅安装Prometheus

```bash
sudo ./deploy/monitoring/scripts/install-prometheus.sh
```

### 仅安装Grafana

```bash
sudo ./deploy/monitoring/scripts/install-grafana.sh
```

### 仅安装AlertManager

```bash
sudo ./deploy/monitoring/scripts/install-alertmanager.sh
```

### 仅安装Exporters

```bash
# 安装所有Exporters
sudo ./deploy/monitoring/scripts/install-exporters.sh all

# 仅安装Node Exporter
sudo ./deploy/monitoring/scripts/install-exporters.sh node

# 仅安装MySQL Exporter
sudo ./deploy/monitoring/scripts/install-exporters.sh mysql

# 仅安装Redis Exporter
sudo ./deploy/monitoring/scripts/install-exporters.sh redis
```

---

## ✅ 验证部署

### 1. 检查容器状态

```bash
docker ps -a | grep zker-
```

**预期输出**:
```
zker-prometheus       Up    0.0.0.0:9090->9090/tcp
zker-grafana          Up    0.0.0.0:3000->3000/tcp
zker-alertmanager     Up    0.0.0.0:9093->9093/tcp
zker-node-exporter    Up    0.0.0.0:9100->9100/tcp
zker-mysqld-exporter  Up    0.0.0.0:9104->9104/tcp
zker-redis-exporter   Up    0.0.0.0:9121->9121/tcp
```

### 2. 检查服务端点

```bash
# Prometheus健康检查
curl http://localhost:9090/-/healthy

# Grafana健康检查
curl http://localhost:3000/api/health

# AlertManager健康检查
curl http://localhost:9093/-/healthy

# Node Exporter指标
curl http://localhost:9100/metrics | head

# MySQL Exporter指标
curl http://localhost:9104/metrics | head

# Redis Exporter指标
curl http://localhost:9121/metrics | head
```

### 3. 访问Web UI

| 服务 | URL | 凭据 |
|------|-----|------|
| **Prometheus** | http://localhost:9090 | 无需登录 |
| **Grafana** | http://localhost:3000 | admin/admin |
| **AlertManager** | http://localhost:9093 | 无需登录 |

---

## ⚙️ 配置管理

### Prometheus配置

**配置文件位置**: `/etc/prometheus/prometheus.yml`

**修改配置**:
```bash
# 1. 编辑配置文件
sudo vim /etc/prometheus/prometheus.yml

# 2. 验证配置语法
docker exec zker-prometheus promtool check config /etc/prometheus/prometheus.yml

# 3. 热重载配置（无需重启）
docker exec zker-prometheus kill -HUP 1

# 4. 或者重启容器
docker restart zker-prometheus
```

**添加新的抓取目标**:
```yaml
scrape_configs:
  - job_name: 'my-service'
    static_configs:
      - targets: ['my-service:8080']
```

### Grafana配置

**配置文件位置**:
- 数据源: `/etc/grafana/provisioning/datasources/`
- 仪表盘: `/var/lib/grafana/dashboards/`

**修改管理员密码**:
```bash
# 方式一：通过环境变量（需要重建容器）
export GRAFANA_ADMIN_PASSWORD=new_password
docker rm -f zker-grafana
sudo ./deploy/monitoring/scripts/install-grafana.sh

# 方式二：通过Web UI（推荐）
# 1. 访问 http://localhost:3000
# 2. 登录后点击左下角头像 → Configuration → Change Password
```

**导入新的仪表盘**:
```bash
# 1. 准备JSON文件
cp my-dashboard.json /var/lib/grafana/dashboards/

# 2. 重启Grafana
docker restart zker-grafana

# 3. 或通过Web UI导入
# Dashboards → Import → Upload JSON file
```

### AlertManager配置

**配置文件位置**: `/etc/alertmanager/alertmanager.yml`

**修改配置**:
```bash
# 1. 编辑配置文件
sudo vim /etc/alertmanager/alertmanager.yml

# 2. 验证配置语法
docker exec zker-alertmanager amtool check-config /etc/alertmanager/alertmanager.yml

# 3. 重启容器
docker restart zker-alertmanager
```

**配置Slack通知**:
```yaml
global:
  slack_api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'

receivers:
  - name: 'slack-alerts'
    slack_configs:
      - channel: '#alerts'
        send_resolved: true
```

---

## 🐛 常见问题

### Q1: 端口被占用

**问题**: `Error: port 9090 already in use`

**解决方案**:
```bash
# 方式一：停止占用端口的进程
sudo lsof -ti:9090 | xargs kill -9

# 方式二：修改部署端口
export PROMETHEUS_PORT=9091
sudo ./deploy/monitoring/scripts/deploy-all.sh
```

### Q2: 容器无法启动

**问题**: 容器创建后立即退出

**排查步骤**:
```bash
# 1. 查看容器状态
docker ps -a | grep zker-

# 2. 查看容器日志
docker logs zker-prometheus
docker logs zker-grafana
docker logs zker-alertmanager

# 3. 检查配置文件语法
docker exec zker-prometheus promtool check config /etc/prometheus/prometheus.yml
docker exec zker-alertmanager amtool check-config /etc/alertmanager/alertmanager.yml
```

### Q3: Grafana无法连接Prometheus

**问题**: Grafana数据源测试失败

**解决方案**:
```bash
# 1. 确认Prometheus运行正常
curl http://localhost:9090/-/healthy

# 2. 确认数据源配置正确
cat /etc/grafana/provisioning/datasources/prometheus.yml

# 3. 检查网络连接
docker exec zker-grafana ping -c 3 zker-prometheus

# 4. 重启Grafana
docker restart zker-grafana
```

### Q4: Exporter无法连接数据库

**问题**: MySQL/Redis Exporter无数据

**解决方案**:
```bash
# 1. 检查环境变量
docker inspect zker-mysqld-exporter | grep -A 5 "Env"

# 2. 确认数据库可访问
docker exec zker-mysqld-exporter mysql -h zker-mysql -u root -p

# 3. 重新创建Exporter（使用正确的DSN）
export MYSQL_EXPORTER_DSN="root:password@(zker-mysql:3306)/"
docker rm -f zker-mysqld-exporter
sudo ./deploy/monitoring/scripts/install-exporters.sh mysql
```

---

## 🔍 故障排查

### 服务无法访问

**症状**: 浏览器无法打开服务页面

**排查流程**:
```bash
# 1. 检查容器状态
docker ps | grep zker-

# 2. 检查端口监听
netstat -tuln | grep -E "9090|3000|9093"

# 3. 检查防火墙
sudo firewall-cmd --list-ports
# 或
sudo ufw status

# 4. 测试本地连接
curl -v http://localhost:9090
```

### 告警未触发

**症状**: 触发条件满足但无告警

**排查流程**:
```bash
# 1. 检查Prometheus告警规则
curl http://localhost:9090/api/v1/rules | python -m json.tool

# 2. 检查当前活跃告警
curl http://localhost:9090/api/v1/alerts | python -m json.tool

# 3. 检查AlertManager接收到的告警
curl http://localhost:9093/api/v1/alerts | python -m json.tool

# 4. 检查AlertManager配置
docker exec zker-alertmanager amtool config show
```

### 数据缺失

**症状**: Grafana仪表盘无数据

**排查流程**:
```bash
# 1. 检查Prometheus targets
curl http://localhost:9090/api/v1/targets | python -m json.tool

# 2. 检查指标是否被抓取
curl 'http://localhost:9090/api/v1/query?query=up' | python -m json.tool

# 3. 检查Exporter输出
curl http://localhost:9100/metrics | grep node_cpu

# 4. 检查Grafana查询
# 在Grafana中打开Explore，手动执行PromQL查询
```

---

## 📊 性能优化

### Prometheus存储优化

```bash
# 1. 修改数据保留时间（默认30天）
docker run -d \
    --name zker-prometheus \
    --restart unless-stopped \
    -p 9090:9090 \
    -v /var/lib/prometheus:/prometheus \
    prom/prometheus:v2.45.0 \
    --storage.tsdb.retention.time=15d

# 2. 压缩历史数据
docker exec zker-prometheus promtool tsdb compact /prometheus
```

### Grafana性能优化

```bash
# 1. 增加内存限制
docker update --memory="2g" --memory-swap="2g" zker-grafana

# 2. 调整查询超时
# 在Grafana UI中：Configuration → Data Sources → Prometheus → Query Timeout: 120s
```

---

## 🧹 清理和卸载

### 停止所有服务

```bash
docker stop zker-prometheus zker-grafana zker-alertmanager zker-node-exporter zker-mysqld-exporter zker-redis-exporter
```

### 删除所有容器

```bash
docker rm -f zker-prometheus zker-grafana zker-alertmanager zker-node-exporter zker-mysqld-exporter zker-redis-exporter
```

### 删除数据目录

```bash
# ⚠️ 警告：此操作将删除所有监控数据！
sudo rm -rf /var/lib/prometheus
sudo rm -rf /var/lib/grafana
sudo rm -rf /var/lib/alertmanager
sudo rm -rf /etc/prometheus
sudo rm -rf /etc/grafana
sudo rm -rf /etc/alertmanager
```

### 完全卸载

```bash
# 1. 停止并删除容器
docker stop $(docker ps -q --filter "name=zker-")
docker rm $(docker ps -aq --filter "name=zker-")

# 2. 删除数据目录
sudo rm -rf /var/lib/prometheus /var/lib/grafana /var/lib/alertmanager
sudo rm -rf /etc/prometheus /etc/grafana /etc/alertmanager

# 3. 删除镜像（可选）
docker rmi prom/prometheus:v2.45.0
docker rmi grafana/grafana:10.0.0
docker rmi prom/alertmanager:v0.26.0
docker rmi prom/node-exporter:v1.7.0
docker rmi prom/mysqld-exporter:v0.15.1
docker rmi oliver006/redis_exporter:v1.55.0
```

---

## 📞 技术支持

如有问题，请联系：
- **DevOps团队**: devops@zker.example.com
- **紧急联系**: +86-xxx-xxxx-xxxx
- **GitHub Issues**: https://github.com/zker/coze-studio/issues

---

## 📚 相关文档

- [Prometheus告警规则说明](./README.md)
- [Grafana使用指南](./GRAFANA_GUIDE.md)
- [性能测试指南](../../tests/performance/README.md)
- [企业级开发规范手册](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)

---

**更新日志**:

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的部署流程和故障排查指南 | DevOps Team |
