# ZKER 监控系统 - 快速开始

**版本**: v1.0.0
**最后更新**: 2025-01-01

---

## 🚀 5分钟快速部署

### 前置条件

```bash
# 检查Docker
docker --version    # 要求 ≥ 20.10

# 检查Docker Compose
docker compose version  # 要求 ≥ 2.0
```

### 一键部署

```bash
# 进入项目目录
cd /path/to/coze-studio

# 一键部署所有监控组件
sudo ./deploy/monitoring/scripts/deploy-all.sh
```

### 访问服务

部署完成后，访问以下服务：

| 服务 | 地址 | 凭据 |
|------|------|------|
| **Prometheus** | http://localhost:9090 | 无需登录 |
| **Grafana** | http://localhost:3000 | admin/admin |
| **AlertManager** | http://localhost:9093 | 无需登录 |

---

## 📊 监控能力概览

### 监控维度（8大维度）

| 维度 | 覆盖内容 | 告警规则 |
|------|---------|---------|
| **系统资源** | CPU、内存、磁盘、网络 | ✅ 4条规则 |
| **API性能** | QPS、延迟、错误率 | ✅ 3条规则 |
| **数据库** | MySQL、Redis性能 | ✅ 5条规则 |
| **缓存** | 命中率、内存使用 | ✅ 2条规则 |
| **配额** | 租户配额使用率 | ✅ 3条规则 |
| **Bot服务** | 调用次数、Token使用 | ✅ 3条规则 |
| **权限系统** | 检查性能、拒绝率 | ✅ 2条规则 |
| **存储** | 使用率、成本 | ✅ 2条规则 |

### Grafana仪表盘（5个）

1. **API性能监控** - QPS、P95延迟、错误率
2. **系统资源监控** - CPU、内存、磁盘
3. **租户业务监控** - 配额、Bot调用、Token
4. **路由性能监控** - 决策延迟、意图识别
5. **计费监控** - 模型成本、API成本

### 告警规则（40+条）

- **Critical级别** - 15条（磁盘严重、API严重错误、配额超限）
- **Warning级别** - 25条（性能降级、配额接近上限）

---

## 📁 目录结构

```
deploy/monitoring/
├── scripts/                      # 部署脚本（SRP原则）
│   ├── deploy-all.sh            # 一键部署脚本
│   ├── install-prometheus.sh    # Prometheus安装
│   ├── install-grafana.sh       # Grafana安装
│   ├── install-alertmanager.sh  # AlertManager安装
│   └── install-exporters.sh     # Exporters安装
│
├── prometheus.yml                # Prometheus主配置
├── alerts.yml                    # 告警规则配置（40+条）
├── alertmanager/
│   └── config.yml               # AlertManager配置
│
├── grafana/
│   └── dashboards/              # Grafana仪表盘JSON
│       ├── api-performance.json
│       ├── system-resources.json
│       ├── tenant-business.json
│       ├── routing-performance.json
│       └── billing-dashboard.json
│
├── README.md                     # 告警规则说明
├── GRAFANA_GUIDE.md             # Grafana使用指南
├── DEPLOYMENT_GUIDE.md          # 部署使用指南
└── DEPLOYMENT_REPORT.md         # 部署报告
```

---

## 🔧 常用命令

### 查看服务状态

```bash
# 查看所有容器
docker ps | grep zker-

# 查看服务日志
docker logs -f zker-prometheus
docker logs -f zker-grafana
docker logs -f zker-alertmanager
```

### 管理服务

```bash
# 停止所有服务
docker stop zker-prometheus zker-grafana zker-alertmanager zker-node-exporter zker-mysqld-exporter zker-redis-exporter

# 启动所有服务
docker start zker-prometheus zker-grafana zker-alertmanager zker-node-exporter zker-mysqld-exporter zker-redis-exporter

# 重启服务
docker restart zker-prometheus
```

### 配置热重载

```bash
# Prometheus配置热重载（无需重启）
docker exec zker-prometheus kill -HUP 1

# AlertManager配置需要重启
docker restart zker-alertmanager
```

---

## 📖 详细文档

| 文档 | 说明 |
|------|------|
| [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) | 完整的部署使用指南 |
| [DEPLOYMENT_REPORT.md](./DEPLOYMENT_REPORT.md) | 详细的部署报告 |
| [README.md](./README.md) | 告警规则说明 |
| [GRAFANA_GUIDE.md](./GRAFANA_GUIDE.md) | Grafana使用指南 |

---

## 🎯 核心特性

### 单一职责原则（SRP）

✅ **每个脚本只负责一个组件**：
- `install-prometheus.sh` - 只负责Prometheus安装
- `install-grafana.sh` - 只负责Grafana安装
- `install-alertmanager.sh` - 只负责AlertManager安装
- `install-exporters.sh` - 只负责Exporters安装

✅ **每个配置文件只负责一个组件**：
- `prometheus.yml` - 只负责Prometheus配置
- `alerts.yml` - 只负责告警规则
- `alertmanager/config.yml` - 只负责告警路由

### 企业级特性

✅ **容器化部署** - 所有组件使用Docker容器
✅ **配置持久化** - 使用Volume持久化配置和数据
✅ **自动重启** - 使用`restart: unless-stopped`策略
✅ **健康检查** - 所有组件配置healthcheck
✅ **模块化设计** - 每个组件可独立安装和升级
✅ **完整文档** - 部署、使用、故障排查文档齐全

---

## 🐛 故障排查

### 服务无法访问

```bash
# 1. 检查容器状态
docker ps -a | grep zker-

# 2. 查看容器日志
docker logs zker-prometheus

# 3. 检查端口占用
netstat -tuln | grep -E "9090|3000|9093"
```

### 告警未触发

```bash
# 查看Prometheus告警规则
curl http://localhost:9090/api/v1/rules | python -m json.tool

# 查看当前活跃告警
curl http://localhost:9090/api/v1/alerts | python -m json.tool

# 查看AlertManager接收到的告警
curl http://localhost:9093/api/v1/alerts | python -m json.tool
```

### Grafana无数据

```bash
# 检查Prometheus targets
curl http://localhost:9090/api/v1/targets | python -m json.tool

# 检查Exporter输出
curl http://localhost:9100/metrics | head
```

---

## 📞 技术支持

如有问题，请联系：
- **DevOps团队**: devops@zker.example.com
- **紧急联系**: +86-xxx-xxxx-xxxx
- **GitHub Issues**: https://github.com/zker/coze-studio/issues

---

**更新日志**:

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2025-01-01 | v1.0.0 | 初始版本，完整的监控系统部署 | DevOps Team |
