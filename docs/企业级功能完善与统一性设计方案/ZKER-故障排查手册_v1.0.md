# ZKER 智能路由与权限管理系统 - 故障排查手册

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**负责人**: 运维与开发团队
**审批人**: 技术架构委员会

---

## 📋 文档概述

### 手册目标

本故障排查手册提供系统常见故障的诊断流程和解决方案，帮助开发和运维人员：

**核心目标**:
- ✅ **快速定位**: 5 分钟内确定故障根因
- ✅ **标准化流程**: 统一的排查步骤和工具
- ✅ **预防为主**: 识别潜在问题，防患于未然
- ✅ **知识积累**: 建立故障知识库，持续改进

### 适用人群

**主要用户**:
- 运维工程师（SRE）
- 后端开发工程师
- 数据库管理员（DBA）
- 测试工程师

**次要用户**:
- 前端开发工程师
- 产品经理（了解故障影响）

---

## 🚨 故障分类与等级

### 故障分类

**按系统层级分类**:
```
L1 - 基础设施层（Infrastructure）
  ├─ 服务器硬件故障
  ├─ 网络故障
  ├─ 操作系统故障
  └─ 虚拟化平台故障

L2 - 平台层（Platform）
  ├─ Kubernetes 故障
  ├─ Docker 容器故障
  ├─ 服务网格故障
  └─ 存储系统故障

L3 - 应用层（Application）
  ├─ API 服务故障
  ├─ 业务逻辑错误
  ├─ 依赖服务故障
  └─ 第三方集成故障

L4 - 数据层（Data）
  ├─ 数据库故障
  ├─ 缓存故障
  ├─ 消息队列故障
  └─ 对象存储故障
```

**按故障类型分类**:
```
1. 服务不可用（Service Down）
   - Pod 崩溃
   - 服务端点无响应
   - 健康检查失败

2. 性能下降（Performance Degradation）
   - 响应时间增加
   - QPS 下降
   - 资源使用率过高

3. 功能异常（Functional Error）
   - API 返回错误
   - 业务逻辑错误
   - 数据不一致

4. 数据问题（Data Issue）
   - 数据丢失
   - 数据损坏
   - 数据不一致

5. 第三方服务故障（Third-party Failure）
   - AI 服务异常
   - 支付服务异常
   - 通知服务异常
```

### 故障等级

**P0 - 严重故障（Critical）** ⚠️⚠️⚠️
```
定义: 系统完全不可用，影响所有用户

响应时间: 立即（< 5 分钟）
解决时间: < 1 小时

示例:
- 所有 API 服务不可用
- 数据库完全宕机
- 核心业务流程中断

影响范围: 100% 用户
处理优先级: 最高
通知对象: 全员 + 管理层
```

**P1 - 高危故障（High）** ⚠️⚠️
```
定义: 核心功能不可用，影响大部分用户

响应时间: < 10 分钟
解决时间: < 4 小时

示例:
- Bot 创建失败
- 对话功能异常
- 知识库检索失败

影响范围: > 50% 用户
处理优先级: 高
通知对象: 技术团队 + 产品负责人
```

**P2 - 中等故障（Medium）** ⚠️
```
定义: 部分功能异常，影响少数用户

响应时间: < 30 分钟
解决时间: < 1 天

示例:
- 特定功能异常
- 部分用户无法登录
- 性能轻微下降

影响范围: < 50% 用户
处理优先级: 中
通知对象: 相关团队
```

**P3 - 低级故障（Low）
```
定义: 非核心功能问题，不影响主要业务

响应时间: < 2 小时
解决时间: < 1 周

示例:
- UI 显示问题
- 非关键功能 Bug
- 性能轻微优化

影响范围: < 10% 用户
处理优先级: 低
通知对象: 直接负责人
```

---

## 🔍 通用故障排查流程

### 故障排查五步法

```
┌─────────────────────────────────────────────────────────┐
│  步骤 1: 确认故障（Validate）                            │
│  ├─ 是否真的故障？                                       │
│  ├─ 故障现象是什么？                                     │
│  ├─ 影响范围多大？                                       │
│  └─ 严重程度如何？                                       │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  步骤 2: 定位根因（Diagnose）                             │
│  ├─ 查看监控和日志                                       │
│  ├─ 分析错误模式和趋势                                   │
│  ├─ 缩小问题范围                                         │
│  └─ 确定根本原因                                         │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  步骤 3: 临时解决（Mitigate）                             │
│  ├─ 快速恢复服务（如重启、回滚）                          │
│  ├─ 降低故障影响（如限流、降级）                          │
│  └─ 争取时间修复根因                                     │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  步骤 4: 根本修复（Resolve）                              │
│  ├─ 修复代码或配置                                       │
│  ├─ 部署修复版本                                         │
│  └─ 验证修复效果                                         │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│  步骤 5: 复盘改进（Improve）                              │
│  ├─ 记录故障详情                                         │
│  ├─ 分析根本原因                                         │
│  ├─ 制定改进措施                                         │
│  └─ 更新监控和文档                                       │
└─────────────────────────────────────────────────────────┘
```

### 快速诊断检查清单

**第一步: 确认故障现象（2 分钟）**
```bash
# 1.1 检查服务状态
kubectl get pods -A | grep zker
kubectl get svc | grep zker

# 1.2 检查错误率
curl -s 'http://prometheus:9090/api/v1/query?query=rate(http_requests_total{status=~"5.."}[5m])' | jq '.data.result[0].value[1]'

# 1.3 检查响应时间
curl -s 'http://prometheus:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket[5m]))' | jq '.data.result[0].value[1]'

# 1.4 检查用户反馈
curl -s http://feedback-api/api/v1/recent?severity=high | jq '.'
```

**第二步: 查看关键指标（3 分钟）**
```bash
# 2.1 CPU 使用率
kubectl top pods -l app=zker
kubectl top nodes

# 2.2 内存使用率
kubectl exec -it $(kubectl get pod -l app=zker -o jsonpath='{.items[0].metadata.name}') -- free -h

# 2.3 磁盘 I/O
kubectl exec -it $(kubectl get pod -l app=zker -o jsonpath='{.items[0].metadata.name}') -- iostat -x 1 5

# 2.4 网络连接
netstat -an | grep :8001 | wc -l
```

**第三步: 查看日志（5 分钟）**
```bash
# 3.1 应用日志
kubectl logs -f -l app=zker --tail=1000

# 3.2 错误日志
kubectl logs -l app=zker --tail=1000 | grep ERROR

# 3.3 慢查询日志
mysql -e "SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 10"

# 3.4 系统日志
journalctl -u kubelet -f
```

---

## 🛠️ 常见故障场景与解决方案

### 场景 1: API 服务不可用

**症状**:
- 所有 API 返回 502/503 错误
- 健康检查失败
- 服务端点无响应

**诊断步骤**:
```bash
# 1. 检查 Pod 状态
kubectl get pods -l app=zker-api

# 预期输出:
# NAME                          READY   STATUS    RESTARTS   AGE
# zker-api-7d9f8b5c6-k2m4n     0/1     Running   0          5m
# zker-api-7d9f8b5c6-x8p9q     0/1     Running   0          5m

# 如果 STATUS 不是 Running，继续诊断

# 2. 查看 Pod 详情
kubectl describe pod zker-api-7d9f8b5c6-k2m4n

# 关注:
# - Events 部分是否有错误信息
# - Last State 是否有退出原因
# - Restart Count 是否过多

# 3. 查看 Pod 日志
kubectl logs zker-api-7d9f8b5c6-k2m4n --previous  # 查看上一次启动的日志
kubectl logs zker-api-7d9f8b5c6-k2m4n -f           # 实时查看当前日志

# 4. 检查容器镜像
kubectl get pod zker-api-7d9f8b5c6-k2m4n -o jsonpath='{.spec.containers[*].image}'

# 5. 检查资源限制
kubectl get pod zker-api-7d9f8b5c6-k2m4n -o jsonpath='{.spec.containers[*].resources}'
```

**常见原因与解决方案**:

**原因 1: 容器启动失败**
```bash
# 诊断
kubectl logs ziker-api-xxx

# 常见错误:
# - "Error: ImagePullBackOff" → 镜像拉取失败
# - "Error: CrashLoopBackOff" → 容器启动后立即崩溃

# 解决方案:
# 1. 镜像拉取失败 → 检查镜像名称和仓库访问权限
kubectl edit deployment zker-api
# 修改 image 字段为正确的镜像地址

# 2. 应用启动失败 → 检查应用配置
kubectl get configmap zker-config -o yaml
kubectl edit configmap zker-config

# 3. 资源不足 → 增加资源限制
kubectl set resources deployment zker-api \
  --limits=cpu=2000m,memory=4Gi \
  --requests=cpu=1000m,memory=2Gi
```

**原因 2: OOMKilled（内存溢出）**
```bash
# 诊断
kubectl describe pod zker-api-xxx | grep OOMKilled

# 查看内存使用历史
kubectl top pod zker-api-xxx --containers

# 解决方案:
# 1. 增加内存限制
kubectl set resources deployment zker-api \
  --limits=memory=8Gi \
  --requests=memory=4Gi

# 2. 优化应用内存使用（代码层面）
# - 检查内存泄漏
# - 优化数据结构
# - 增加分页查询

# 3. 水平扩展 Pod
kubectl scale deployment zker-api --replicas=10
```

**原因 3: 健康检查失败**
```bash
# 诊断
kubectl describe pod zker-api-xxx | grep -A 10 Liveness

# 查看健康检查配置
kubectl get deployment zker-api -o yaml | grep -A 10 livenessProbe

# 手动测试健康检查
curl http://pod-ip:8001/health

# 常见错误:
# - "Get http://pod-ip:8001/health: dial tcp: connection refused"
#   → 应用未启动或端口配置错误

# 解决方案:
# 1. 调整健康检查配置
kubectl patch deployment zker-api -p '
{
  "spec": {
    "template": {
      "spec": {
        "containers": [{
          "name": "zker-api",
          "livenessProbe": {
            "httpGet": {
              "path": "/health",
              "port": 8001
            },
            "initialDelaySeconds": 30,
            "periodSeconds": 10,
            "timeoutSeconds": 5,
            "failureThreshold": 3
          }
        }]
      }
    }
  }
}'

# 2. 修复健康检查端点
# 确保应用监听正确的端口，/health 端点正常响应
```

---

### 场景 2: 数据库连接失败

**症状**:
- API 返回 "database connection error"
- 日志中出现 "too many connections"
- 数据库连接池耗尽

**诊断步骤**:
```bash
# 1. 检查数据库服务状态
kubectl get pods -l app=mysql

# 2. 检查数据库连接数
kubectl exec -it mysql-0 -- mysql -e "SHOW PROCESSLIST"

# 3. 检查数据库连接池配置
kubectl get configmap zker-config -o yaml | grep -A 10 DATABASE

# 4. 查看应用日志
kubectl logs -l app=zker-api | grep -i "database\|mysql\|connection"
```

**常见原因与解决方案**:

**原因 1: 连接数达到上限**
```bash
# 诊断
kubectl exec -it mysql-0 -- mysql -e "SHOW VARIABLES LIKE 'max_connections';"
kubectl exec -it mysql-0 -- mysql -e "SHOW STATUS LIKE 'Threads_connected';"

# 解决方案:
# 1. 临时增加连接数（立即生效）
kubectl exec -it mysql-0 -- mysql -e "SET GLOBAL max_connections = 500;"

# 2. 永久增加连接数（修改配置）
kubectl edit configmap mysql-config
# 修改 max_connections = 500

kubectl rollout restart deployment mysql

# 3. 优化应用连接池
# 修改应用配置，减少空闲连接
kubectl edit configmap zker-config
# 设置:
# max_open_conns: 100 (从 200 减少)
# max_idle_conns: 10  (从 50 减少)
# conn_max_lifetime: 300s (从 1h 减少)
```

**原因 2: 慢查询导致连接堆积**
```bash
# 诊断
# 1. 查看当前正在执行的查询
kubectl exec -it mysql-0 -- mysql -e "SHOW FULL PROCESSLIST;"

# 2. 查看慢查询日志
kubectl exec -it mysql-0 -- mysql -e "
  SELECT * FROM mysql.slow_log
  ORDER BY query_time DESC
  LIMIT 10;
"

# 3. 分析慢查询
kubectl exec -it mysql-0 -- mysql -e "
  EXPLAIN SELECT * FROM conversations
  WHERE tenant_id = 'xxx'
  ORDER BY created_at DESC
  LIMIT 20;
"

# 解决方案:
# 1. 终止长时间运行的查询
kubectl exec -it mysql-0 -- mysql -e "KILL <process_id>;"

# 2. 优化慢查询
# - 添加索引
# - 重写查询语句
# - 分页查询

# 3. 启用查询缓存
kubectl exec -it mysql-0 -- mysql -e "SET GLOBAL query_cache_type = ON;"
kubectl exec -it mysql-0 -- mysql -e "SET GLOBAL query_cache_size = 268435456;"  # 256MB
```

**原因 3: 网络问题**
```bash
# 诊断
# 1. 测试 Pod 到数据库的网络连通性
kubectl exec -it $(kubectl get pod -l app=zker-api -o jsonpath='{.items[0].metadata.name}') -- nc -zv mysql-service 3306

# 2. 检查 DNS 解析
kubectl exec -it $(kubectl get pod -l app=zker-api -o jsonpath='{.items[0].metadata.name}') -- nslookup mysql-service

# 3. 检查 Service 端点
kubectl get endpoints mysql-service

# 解决方案:
# 1. 修复 DNS 问题
kubectl edit configmap coredns
# 添加合适的上游 DNS

# 2. 修复网络策略
kubectl get networkpolicy -A
# 确保允许 API Pod 访问数据库

# 3. 重启网络组件
kubectl rollout restart deployment coredns -n kube-system
```

---

### 场景 3: API 响应缓慢

**症状**:
- P95 响应时间 > 5s
- 用户反馈系统卡顿
- API 超时

**诊断步骤**:
```bash
# 1. 查看响应时间监控
curl -s 'http://prometheus:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket{endpoint="/api/v1/bots"}[5m]))' | jq .

# 2. 查看慢请求分布
kubectl logs -l app=zker-api --tail=1000 | grep "duration" | awk '{print $NF}' | sort -n | tail -20

# 3. 分析应用性能
kubectl exec -it $(kubectl get pod -l app=zker-api -o jsonpath='{.items[0].metadata.name}') -- pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

**常见原因与解决方案**:

**原因 1: 数据库查询慢**
```bash
# 诊断（参考场景 2）

# 解决方案:
# 1. 添加索引
kubectl exec -it mysql-0 -- mysql -e "
  CREATE INDEX idx_tenant_created ON conversations(tenant_id, created_at);
"

# 2. 优化查询
# 避免使用 SELECT *
# 使用覆盖索引
# 分页查询

# 3. 使用缓存
# 在应用层添加 Redis 缓存
```

**原因 2: 外部 API 调用慢**
```bash
# 诊断
# 查看应用日志中的外部 API 调用记录
kubectl logs -l app=zker-api | grep "external_api"

# 解决方案:
# 1. 添加超时控制
# 修改代码，为外部 API 调用设置超时
httpClient.setTimeout(5000)  # 5 秒超时

# 2. 使用熔断器
# 当外部服务故障时，快速失败
circuitBreaker.open()

# 3. 使用缓存
# 缓存外部 API 的响应
```

**原因 3: CPU 密集型计算**
```bash
# 诊断
kubectl top pod -l app=zker-api

# 查看 CPU 使用率高的线程
kubectl exec -it $(kubectl get pod -l app=zker-api -o jsonpath='{.items[0].metadata.name}') -- top -H

# 解决方案:
# 1. 水平扩展
kubectl scale deployment zker-api --replicas=10

# 2. 优化算法
# 使用更高效的算法
# 减少不必要的计算

# 3. 异步处理
# 将耗时操作放入消息队列异步处理
```

---

### 场景 4: Redis 缓存故障

**症状**:
- 缓存命中率下降
- 数据库压力增加
- API 响应变慢

**诊断步骤**:
```bash
# 1. 检查 Redis 服务状态
kubectl get pods -l app=redis

# 2. 检查 Redis 连接
kubectl exec -it redis-0 -- redis-cli ping

# 3. 查看缓存命中率
kubectl exec -it redis-0 -- redis-cli info stats | grep keyspace_hits

# 4. 查看慢查询日志
kubectl exec -it redis-0 -- redis-cli slowlog get 10
```

**常见原因与解决方案**:

**原因 1: Redis 内存不足**
```bash
# 诊断
kubectl exec -it redis-0 -- redis-cli info memory | grep used_memory_human

# 解决方案:
# 1. 清理过期 key
kubectl exec -it redis-0 -- redis-cli --scan --pattern "session:*" | xargs redis-cli del

# 2. 设置 maxmemory-policy
kubectl exec -it redis-0 -- redis-cli CONFIG SET maxmemory-policy allkeys-lru

# 3. 增加 Redis 内存
kubectl edit statefulset redis
# 增加 resources.requests.memory
```

**原因 2: Redis 连接数过多**
```bash
# 诊断
kubectl exec -it redis-0 -- redis-cli info clients | grep connected_clients

# 解决方案:
# 1. 调整应用连接池配置
kubectl edit configmap zker-config
# 减少 Redis 连接池大小

# 2. 增加 Redis 最大连接数
kubectl exec -it redis-0 -- redis-cli CONFIG SET maxclients 10000
```

---

### 场景 5: 消息队列积压

**症状**:
- NSQ 消息积压
- 消费者处理速度慢
- 异步任务延迟

**诊断步骤**:
```bash
# 1. 查看队列深度
curl -s http://nsqadmin:4171/api/topics | jq '.topics[] | {name: .name, depth: .depth}'

# 2. 查看消费者状态
curl -s http://nsqadmin:4171/api/consumers | jq .

# 3. 查看消费者日志
kubectl logs -l app=worker --tail=1000
```

**常见原因与解决方案**:

**原因 1: 消费者处理慢**
```bash
# 解决方案:
# 1. 增加消费者数量
kubectl scale deployment worker --replicas=10

# 2. 优化消费逻辑
# 减少单个消息的处理时间
# 批量处理消息

# 3. 使用优先级队列
# 重要消息优先处理
```

**原因 2: 消费者崩溃**
```bash
# 诊断
kubectl logs -l app=worker | grep ERROR

# 解决方案:
# 1. 查看崩溃原因并修复 Bug

# 2. 增加重启策略
kubectl patch deployment worker -p '
{
  "spec": {
    "template": {
      "spec": {
        "containers": [{
          "name": "worker",
          "imagePullPolicy": "Always"
        }]
      }
    }
  }
}'

# 3. 增加健康检查
kubectl patch deployment worker -p '
{
  "spec": {
    "template": {
      "spec": {
        "containers": [{
          "name": "worker",
          "livenessProbe": {
            "exec": {
              "command": ["pgrep", "-f", "worker"]
            },
            "initialDelaySeconds": 30,
            "periodSeconds": 10
          }
        }]
      }
    }
  }
}'
```

---

## 🧰 故障排查工具集

### 系统工具

**1. kubectl** (Kubernetes 命令行工具)
```bash
# 查看资源状态
kubectl get all
kubectl get pods,svc,deploy

# 查看资源详情
kubectl describe pod <pod-name>
kubectl describe node <node-name>

# 查看日志
kubectl logs <pod-name> -f
kubectl logs <pod-name> --previous  # 查看上一次的日志

# 执行命令
kubectl exec -it <pod-name> -- /bin/bash

# 端口转发
kubectl port-forward <pod-name> 8080:80

# 应用配置
kubectl apply -f config.yaml
kubectl rollout restart deployment <deployment-name>

# 调试
kubectl debug -it <pod-name> --image=nicolaka/netshoot
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- bash
```

**2. curl / wget** (HTTP 请求工具)
```bash
# 测试 API 端点
curl http://api.zker.com/health
curl -X POST http://api.zker.com/api/v1/bots -H "Content-Type: application/json" -d '{"name": "test"}'

# 查看响应头
curl -I http://api.zker.com/health

# 测试响应时间
curl -w "@curl-format.txt" -o /dev/null -s http://api.zker.com/api/v1/bots

# curl-format.txt 内容:
# time_namelookup: %{time_namelookup}\n
# time_connect: %{time_connect}\n
# time_starttransfer: %{time_starttransfer}\n
# time_total: %{time_total}\n
```

**3. tcpdump / wireshark** (网络抓包)
```bash
# 抓包分析
tcpdump -i any host api.zker.com -w capture.pcap

# 实时查看 HTTP 请求
tcpdump -i any -A -s 0 'tcp port 8001 and (((ip[2:2] - ((ip[0]&0xf)<<2)) - ((tcp[12]&0xf0)>>2)) != 0)'

# 查看 TCP 连接
netstat -antp | grep :8001
ss -antp | grep :8001
```

**4. strace / ltrace** (系统调用追踪)
```bash
# 追踪系统调用
kubectl exec -it <pod-name> -- strace -p 1

# 追踪库函数调用
kubectl exec -it <pod-name> -- ltrace -p 1

# 查看文件访问
kubectl exec -it <pod-name> -- strace -e trace=open,openat,read,write -p 1
```

### 监控工具

**1. Prometheus + Grafana**
```bash
# Prometheus 查询
# 错误率
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# QPS
rate(http_requests_total[5m])

# P95 响应时间
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# CPU 使用率
rate(process_cpu_seconds_total[5m]) * 100

# 内存使用率
process_resident_memory_bytes / node_memory_MemTotal_bytes * 100

# Grafana Dashboard
# 访问: http://grafana.zker.com
# 关键 Dashboard:
# - ZKER API 监控
# - ZKER 数据库监控
# - ZKER 系统资源监控
```

**2. Jaeger** (分布式追踪)
```bash
# 访问 Jaeger UI
# http://jaeger.zker.com

# 查询 Trace
# - 选择服务: zker-api
# - 选择时间范围
# - 查看慢请求的 Trace

# 分析 Trace
# - 查找耗时最长的 Span
# - 分析调用链路
# - 定位瓶颈
```

### 日志分析工具

**1. Elasticsearch + Kibana**
```bash
# Kibana 查询语法
# 查看错误日志
level:ERROR

# 查看特定服务的日志
kubernetes.service.name:"zker-api"

# 查看慢请求
request_duration:>5000

# 组合查询
level:ERROR AND kubernetes.service.name:"zker-api"

# 时间范围查询
@timestamp:[now-1h TO now]
```

**2. jq** (JSON 处理)
```bash
# 格式化 JSON 输出
curl http://api.zker.com/health | jq .

# 提取特定字段
curl http://api.zker.com/api/v1/bots | jq '.bots[].name'

# 过滤数据
curl http://api.zker.com/api/v1/bots | jq '.bots[] | select(.status == "published")'

# 统计数据
kubectl get pods -o json | jq '.items[] | .spec.nodeName' | sort | uniq -c
```

---

## 📋 故障复盘模板

### 故障报告

**故障基本信息**:
```yaml
故障名称: [简短描述]
故障等级: [P0/P1/P2/P3]
故障时间: [开始时间 - 结束时间]
影响范围: [影响用户数、影响功能]
负责人: [主要处理人]
参与人员: [协作团队]
```

**故障时间线**:
```yaml
时间轴:
  - 时间: "2025-01-01 10:00"
    事件: 监控告警触发，错误率 > 10%
    操作: 自动告警通知到 on-call 工程师

  - 时间: "2025-01-01 10:02"
    事件: 工程师确认故障，开始排查
    操作: 查看 Grafana Dashboard，发现 API 服务响应超时

  - 时间: "2025-01-01 10:05"
    事件: 定位到数据库慢查询
    操作: 查看 MySQL slow log，发现 conversations 表查询无索引

  - 时间: "2025-01-01 10:10"
    事件: 执行临时修复，添加索引
    操作: 执行 CREATE INDEX，错误率下降到正常水平

  - 时间: "2025-01-01 10:15"
    事件: 验证修复效果，服务恢复正常
    操作: 查看监控指标，确认所有指标正常

  - 时间: "2025-01-01 11:00"
    事件: 故障复盘会议
    操作: 团队讨论根因和改进措施
```

**根本原因分析** (5 Whys):
```
为什么发生？
  1. 为什么 API 响应超时？
     → 数据库查询慢

  2. 为什么数据库查询慢？
     → 缺少索引，全表扫描

  3. 为什么缺少索引？
     → 上线时没有审查 SQL 性能

  4. 为什么没有审查 SQL 性能？
     → 缺少代码审查流程和性能测试

  5. 为什么缺少流程和测试？
     → 团队对性能重视不足，没有建立规范

根本原因: 缺少上线前性能审查流程
```

**改进措施**:
```yaml
短期措施 (1 周内):
  - task: 为所有慢查询添加索引
    owner: DBA 团队
    deadline: 2025-01-05

  - task: 增加数据库慢查询监控告警
    owner: 运维团队
    deadline: 2025-01-03

中期措施 (1 月内):
  - task: 建立上线前性能审查流程
    owner: 技术架构委员会
    deadline: 2025-01-31

  - task: 编写 SQL 性能优化规范
    owner: DBA 团队
    deadline: 2025-01-31

长期措施 (1 季度内):
  - task: 引入自动化性能测试工具
    owner: 测试团队
    deadline: 2025-03-31

  - task: 开展团队 SQL 性能培训
    owner: 培训团队
    deadline: 2025-03-31
```

**行动计划跟踪**:
```yaml
行动项清单:
  - id: ACTION-001
    task: 为 conversations 表添加索引
    owner: 张三 (DBA)
    status: ✅ 已完成
    due_date: 2025-01-05
    completed_date: 2025-01-03

  - id: ACTION-002
    task: 配置慢查询告警
    owner: 李四 (SRE)
    status: ⏳ 进行中
    due_date: 2025-01-03
    progress: 80%

  - id: ACTION-003
    task: 建立性能审查流程
    owner: 王五 (架构师)
    status: ⏳ 待开始
    due_date: 2025-01-31
```

---

## 📚 附录

### A. 常用命令速查

**Kubernetes**:
```bash
# 查看 Pod
kubectl get pods -A
kubectl get pods -l app=zker
kubectl get pods -o wide

# 查看日志
kubectl logs <pod-name>
kubectl logs <pod-name> -f
kubectl logs <pod-name> --tail=100
kubectl logs -l app=zker --tail=1000

# 执行命令
kubectl exec -it <pod-name> -- bash
kubectl exec -it <pod-name> -- sh

# 端口转发
kubectl port-forward <pod-name> 8080:80

# 资源管理
kubectl top pods
kubectl top nodes
kubectl describe pod <pod-name>
kubectl describe node <node-name>

# 扩缩容
kubectl scale deployment <name> --replicas=5
kubectl autoscale deployment <name> --min=3 --max=10 --cpu-percent=80
```

**数据库**:
```bash
# 连接数据库
mysql -h <host> -u <user> -p

# 查看连接数
SHOW PROCESSLIST;

# 查看慢查询
SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 10;

# 查看表大小
SELECT
  table_name,
  ROUND(data_length / 1024 / 1024, 2) as data_mb,
  ROUND(index_length / 1024 / 1024, 2) as index_mb
FROM information_schema.TABLES
WHERE table_schema = 'zker_production'
ORDER BY data_mb DESC;

# 分析查询
EXPLAIN SELECT * FROM conversations WHERE tenant_id = 'xxx';

# 查看索引
SHOW INDEX FROM conversations;

# 查看表结构
DESCRIBE conversations;
SHOW CREATE TABLE conversations;
```

**Redis**:
```bash
# 连接 Redis
redis-cli -h <host> -p 6379

# 查看 info
info
info stats
info memory

# 查看 key
keys *

# 查看 key 的值
get <key>
hgetall <hash-key>

# 查看 key 的 TTL
ttl <key>

# 慢查询
slowlog get 10

# 清理过期 key
scan 0 match session:* count 1000
```

### B. 监控指标参考

**应用指标**:
```yaml
# 请求指标
http_requests_total: 总请求数
http_requests_duration_seconds: 请求耗时
http_requests_by_path: 按 path 分组的请求数
http_requests_by_status: 按 status 分组的请求数

# 错误指标
http_requests_total{status=~"5.."}: 5xx 错误数
http_requests_total{status=~"4.."}: 4xx 错误数
exception_total: 异常总数

# 业务指标
bots_created_total: 创建 Bot 总数
conversations_created_total: 创建对话总数
messages_sent_total: 发送消息总数
knowledge_searches_total: 知识库检索总数
```

**系统指标**:
```yaml
# CPU
process_cpu_seconds_total: 进程 CPU 使用时间
node_cpu_seconds_total: 节点 CPU 使用时间

# 内存
process_resident_memory_bytes: 进程常驻内存
node_memory_MemTotal_bytes: 节点总内存
node_memory_MemAvailable_bytes: 节点可用内存

# 磁盘
node_filesystem_size_bytes: 文件系统总大小
node_filesystem_avail_bytes: 文件系统可用大小
node_disk_io_time_seconds_total: 磁盘 I/O 时间

# 网络
node_network_receive_bytes_total: 网络接收字节数
node_network_transmit_bytes_total: 网络发送字节数
```

**数据库指标**:
```yaml
# 连接
mysql_global_status_threads_connected: 当前连接数
mysql_global_status_max_connections: 最大连接数
mysql_global_status_threads_running: 正在运行的线程数

# 查询
mysql_global_status_questions: 总查询数
mysql_global_status_slow_queries: 慢查询数

# 复制
mysql_slave_status_seconds_behind_master: 从库延迟秒数
```

### C. 联系方式

| 角色 | 姓名 | 联系方式 | 职责 | 可联系时间 |
|-----|------|---------|------|-----------|
| On-call 工程师（周1） | [待填写] | [待填写] | 一线响应 | 24×7 |
| On-call 工程师（周2） | [待填写] | [待填写] | 一线响应 | 24×7 |
| On-call 工程师（周3） | [待填写] | [待填写] | 一线响应 | 24×7 |
| On-call 工程师（周4） | [待填写] | [待填写] | 一线响应 | 24×7 |
| On-call 工程师（周5） | [待填写] | [待填写] | 一线响应 | 24×7 |
| On-call 工程师（周6） | [待填写] | [待填写] | 一线响应 | 24×7 |
| On-call 工程师（周日） | [待填写] | [待填写] | 一线响应 | 24×7 |
| 技术负责人 | [待填写] | [待填写] | 二线支持 | 工作时间 |
| DBA 负责人 | [待填写] | [待填写] | 数据库专家 | 工作时间 |
| 运维负责人 | [待填写] | [待填写] | 基础设施专家 | 工作时间 |

**紧急联系流程**:
```
1. 一线响应: On-call 工程师（电话/Slack）
   ├─ 10 分钟内响应
   ├─ 30 分钟内给出初步诊断
   └─ 1 小时内解决或升级

2. 二线支持: 技术负责人（如果 P0 故障）
   ├─ 立即介入
   └─ 协调资源解决

3. 管理层通报: CTO/VP（如果影响 > 50% 用户）
   ├─ 每 30 分钟更新进度
   └─ 故障解除后发送事故报告
```

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 运维与开发团队 |

**审批记录**:

| 角色 | 姓名 | 审批意见 | 日期 |
|-----|------|---------|------|
| 技术架构委员会 | [待填写] | [待审批] | [待审批] |
| 运维负责人 | [待填写] | [待审批] | [待审批] |
| 开发负责人 | [待填写] | [待审批] | [待审批] |

---

**© 2025 ZKER Project. All rights reserved.**
