# 服务健康检查集成指南

**文档版本**: v1.0
**创建日期**: 2025-01-01
**负责人**: 后端架构团队

---

## 📋 功能概述

健康检查服务用于监控应用及其依赖项（数据库、缓存、消息队列等）的健康状态，支持：

- ✅ **多种健康检查器**: MySQL、Redis、HTTP服务、通用数据库
- ✅ **三种健康状态**: UP（健康）、DOWN（不可用）、DEGRADED（降级）
- ✅ **并发检查**: 并发执行所有检查项，提高效率
- ✅ **Kubernetes兼容**: 支持liveness和readiness探针
- ✅ **Prometheus指标**: 输出Prometheus格式的监控指标
- ✅ **详细诊断信息**: 每个检查项都包含详细的诊断信息

---

## 🚀 快速开始

### 1. 初始化健康检查服务

在应用启动时初始化健康检查服务：

```go
// backend/main.go 或 backend/application/init.go

package main

import (
    "context"
    "github.com/coze-studio/backend/pkg/health"
    "gorm.io/gorm"
    "github.com/redis/go-redis/v9"
)

func initHealthCheck(db *gorm.DB, redisClient *redis.Client) {
    // 1. 获取服务版本
    version := getServiceVersion()

    // 2. 创建健康检查服务
    healthSvc := health.NewHealthService(version)

    // 3. 注册MySQL检查器
    mysqlChecker := health.NewMySQLChecker(db)
    healthSvc.RegisterChecker(mysqlChecker)

    // 4. 注册Redis检查器
    redisChecker := health.NewRedisChecker(redisClient)
    healthSvc.RegisterChecker(redisChecker)

    // 5. 注册HTTP服务检查器（可选）
    // httpClient := NewHTTPClient()
    // httpChecker := health.NewHTTPChecker("external_api", "https://api.example.com/health", 10*time.Second, httpClient)
    // healthSvc.RegisterChecker(httpChecker)

    // 6. 初始化全局健康检查服务
    health.InitGlobalHealthService(version)

    logs.Info("Health check service initialized")
}

func getServiceVersion() string {
    // 从环境变量、构建信息或配置文件读取版本号
    version := os.Getenv("APP_VERSION")
    if version == "" {
        version = "1.0.0" // 默认版本
    }
    return version
}
```

### 2. 注册健康检查路由

```go
// backend/router/router.go

package router

import (
    "github.com/cloudwego/hertz/pkg/server"
    "github.com/coze-studio/backend/pkg/health"
)

// RegisterHealthRoutes 注册健康检查路由
func RegisterHealthRoutes(r *server.Hertz) {
    // 简单健康检查（用于负载均衡器）
    r.GET("/health", health.HealthCheckHandler)

    // 详细健康报告（用于监控和调试）
    r.GET("/health/detailed", health.HealthCheckDetailedHandler)

    // 单项健康检查
    r.GET("/health/check/:name", health.HealthCheckSingleHandler)

    // Kubernetes探针
    r.GET("/health/live", health.LivenessProbeHandler)   // 存活探针
    r.GET("/health/ready", health.ReadinessProbeHandler) // 就绪探针

    // 版本信息
    r.GET("/version", health.VersionHandler)

    // Prometheus指标
    r.GET("/health/metrics", health.MetricsHandler)
}
```

### 3. 在main.go中集成

```go
// backend/main.go

func main() {
    ctx := context.Background()

    // ... 其他初始化代码 ...

    // 初始化健康检查
    initHealthCheck(application.GetDB(), application.GetRedis())

    // 启动HTTP服务器
    startHttpServer()
}

func startHttpServer() {
    // ... 服务器配置 ...

    s := server.Default(opts...)

    // 注册路由
    router.RegisterHealthRoutes(s)
    router.RegisterAPIRoutes(s)

    // ... 启动服务器 ...
}
```

---

## 🔧 使用示例

### 示例 1: MySQL健康检查

```go
import (
    "github.com/coze-studio/backend/pkg/health"
    "gorm.io/gorm"
)

func setupMySQLHealthCheck(db *gorm.DB) {
    mysqlChecker := health.NewMySQLChecker(db)

    healthSvc := health.GetGlobalHealthService()
    healthSvc.RegisterChecker(mysqlChecker)
}
```

**检查结果示例**：
```json
{
  "name": "mysql",
  "status": "UP",
  "message": "Database is healthy",
  "timestamp": 1735689600000,
  "duration": 15,
  "details": {
    "open_connections": 10,
    "in_use": 3,
    "idle": 7,
    "wait_count": 0,
    "wait_duration_ms": 0,
    "max_idle_closed": 1,
    "max_lifetime_closed": 0,
    "ping_duration_ms": 5,
    "version": "8.4.5"
  }
}
```

### 示例 2: Redis健康检查

```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/coze-studio/backend/pkg/health"
)

func setupRedisHealthCheck(client *redis.Client) {
    redisChecker := health.NewRedisChecker(client)

    healthSvc := health.GetGlobalHealthService()
    healthSvc.RegisterChecker(redisChecker)
}
```

**检查结果示例**：
```json
{
  "name": "redis",
  "status": "UP",
  "message": "Redis is healthy",
  "timestamp": 1735689600000,
  "duration": 3,
  "details": {
    "ping_duration_ms": 3
  }
}
```

### 示例 3: HTTP服务健康检查

```go
import (
    "context"
    "net/http"
    "time"
    "github.com/coze-studio/backend/pkg/health"
)

// HTTPClient 实现
type MyHTTPClient struct{}

func (c *MyHTTPClient) Get(ctx context.Context, url string) (int, string, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return 0, "", err
    }
    defer resp.Body.Close()

    body := new(strings.Builder)
    io.Copy(body, resp.Body)

    return resp.StatusCode, body.String(), nil
}

func setupHTTPHealthCheck() {
    client := &MyHTTPClient{}
    httpChecker := health.NewHTTPChecker(
        "external_api",
        "https://api.example.com/health",
        10*time.Second,
        client,
    )

    healthSvc := health.GetGlobalHealthService()
    healthSvc.RegisterChecker(httpChecker)
}
```

### 示例 4: 自定义健康检查器

```go
import (
    "context"
    "github.com/coze-studio/backend/pkg/health"
)

// MyCustomChecker 自定义健康检查器
type MyCustomChecker struct {
    name string
}

func NewMyCustomChecker(name string) *MyCustomChecker {
    return &MyCustomChecker{name: name}
}

func (c *MyCustomChecker) Name() string {
    return c.name
}

func (c *MyCustomChecker) Check(ctx context.Context) *health.CheckResult {
    result := &health.CheckResult{
        Name:      c.name,
        Timestamp: time.Now().UnixMilli(),
        Details:   make(map[string]interface{}),
    }

    // 执行自定义检查逻辑
    if err := checkMyService(ctx); err != nil {
        result.Status = health.HealthStatusDown
        result.Message = fmt.Sprintf("Service check failed: %v", err)
        result.Details["error"] = err.Error()
        return result
    }

    result.Status = health.HealthStatusUp
    result.Message = "Service is healthy"
    return result
}

func checkMyService(ctx context.Context) error {
    // 自定义检查逻辑
    return nil
}

// 注册自定义检查器
func setupCustomHealthCheck() {
    checker := NewMyCustomChecker("my_service")

    healthSvc := health.GetGlobalHealthService()
    healthSvc.RegisterChecker(checker)
}
```

---

## 📊 API端点说明

### GET /health
**说明**: 简单健康检查（用于负载均衡器）
**响应格式**: 纯文本
**示例**:
```bash
# 健康状态
curl http://localhost:8888/health
# 返回: UP

# 不健康状态
curl http://localhost:8888/health
# 返回: DOWN
```

### GET /health/detailed
**说明**: 详细健康报告
**响应格式**: JSON
**示例**:
```bash
curl http://localhost:8888/health/detailed
```

**响应示例**：
```json
{
  "status": "UP",
  "timestamp": 1735689600000,
  "checks": {
    "mysql": {
      "name": "mysql",
      "status": "UP",
      "message": "Database is healthy",
      "timestamp": 1735689600000,
      "duration": 15,
      "details": {
        "open_connections": 10,
        "in_use": 3,
        "idle": 7
      }
    },
    "redis": {
      "name": "redis",
      "status": "UP",
      "message": "Redis is healthy",
      "timestamp": 1735689600000,
      "duration": 3,
      "details": {
        "ping_duration_ms": 3
      }
    }
  },
  "version": "1.0.0",
  "uptime": 3600
}
```

**HTTP状态码映射**:
- `200 OK`: 整体健康（UP或DEGRADED）
- `503 Service Unavailable`: 整体不健康（DOWN）

### GET /health/check/:name
**说明**: 单项健康检查
**参数**: `name` - 检查器名称
**示例**:
```bash
# 检查MySQL
curl http://localhost:8888/health/check/mysql

# 检查Redis
curl http://localhost:8888/health/check/redis
```

### GET /health/live
**说明**: Kubernetes存活探针
**响应**: `OK` 或 `NOT_OK`
**用途**: 检查服务是否存活（是否响应）

### GET /health/ready
**说明**: Kubernetes就绪探针
**响应**: `READY` 或 `NOT_READY`
**用途**: 检查服务是否就绪（关键组件是否健康）

### GET /version
**说明**: 版本信息
**响应格式**: JSON
**示例**:
```bash
curl http://localhost:8888/version
```

**响应示例**：
```json
{
  "version": "1.0.0",
  "go_version": "go1.24.0",
  "build_time": "2025-01-01T00:00:00Z",
  "git_commit": "abc123",
  "git_summary": "v1.0.0",
  "uptime": 3600,
  "timestamp": 1735689600000
}
```

### GET /health/metrics
**说明**: Prometheus格式指标
**响应格式**: 纯文本（Prometheus格式）
**示例**:
```bash
curl http://localhost:8888/health/metrics
```

**响应示例**：
```
# HELP health_status Overall health status (0=unknown, 1=up, 2=degraded, 3=down)
# TYPE health_status gauge
health_status 1

# HELP health_check_status Health check status (0=unknown, 1=up, 2=degraded, 3=down)
# TYPE health_check_status gauge
health_check_status{name="mysql"} 1
health_check_status{name="redis"} 1
health_check_duration_ms{name="mysql"} 15
health_check_duration_ms{name="redis"} 3

# HELP app_uptime_seconds Application uptime in seconds
# TYPE app_uptime_seconds gauge
app_uptime_seconds 3600
```

---

## 🔌 Kubernetes集成

### Liveness Probe配置

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: coze-studio-backend
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: coze-studio-backend
    spec:
      containers:
      - name: backend
        image: coze-studio/backend:latest
        ports:
        - containerPort: 8888
        # 存活探针：检查服务是否响应
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8888
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        # 就绪探针：检查服务是否就绪
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8888
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
```

### HPA（水平Pod自动扩缩容）配置

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: coze-studio-backend-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: coze-studio-backend
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metric:
        name: health_status
        target:
          type: AverageValue
          averageValue: "1"  # 保持健康状态为1
```

---

## 📈 Prometheus监控集成

### Prometheus配置

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'coze-studio-backend'
    static_configs:
      - targets: ['backend:8888']
    metrics_path: '/health/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard查询

**整体健康状态**:
```promql
health_status
```

**各组件健康状态**:
```promql
health_check_status
```

**检查耗时**:
```promql
rate(health_check_duration_ms[5m])
```

**服务运行时长**:
```promql
app_uptime_seconds
```

---

## ⚠️ 注意事项

### 1. 性能考虑

- ✅ **并发检查**: 所有检查项并发执行，提高效率
- ✅ **超时控制**: 每个检查都有独立的超时控制
- ✅ **避免频繁检查**: 建议间隔10-30秒

### 2. 关键组件标记

以下组件标记为**关键组件**，其DOWN状态会导致整体状态DOWN：
- MySQL/数据库
- Redis/缓存

其他组件（如外部HTTP服务）DOWN不会导致整体DOWN，但会标记为DEGRADED。

### 3. 健康状态判断

**整体状态计算逻辑**：
1. 如果任何关键组件DOWN → 整体DOWN
2. 如果任何组件DEGRADED → 整体DEGRADED
3. 如果非关键组件DOWN → 整体DEGRADED
4. 所有组件UP → 整体UP

### 4. 错误处理

- ✅ 健康检查失败不会影响主服务
- ✅ 检查器执行错误会被捕获并记录
- ✅ 每个检查器独立运行，互不影响

---

## 🧪 测试

### 手动测试

```bash
# 1. 测试简单健康检查
curl http://localhost:8888/health

# 2. 测试详细健康报告
curl http://localhost:8888/health/detailed | jq

# 3. 测试单项检查
curl http://localhost:8888/health/check/mysql | jq

# 4. 测试版本信息
curl http://localhost:8888/version | jq

# 5. 测试Prometheus指标
curl http://localhost:8888/health/metrics
```

### 模拟故障测试

```bash
# 1. 停止MySQL，观察健康检查
docker stop mysql
curl http://localhost:8888/health/detailed | jq

# 2. 停止Redis，观察健康检查
docker stop redis
curl http://localhost:8888/health/detailed | jq

# 3. 重启服务，观察恢复
docker start mysql
curl http://localhost:8888/health/detailed | jq
```

---

## 📚 相关文档

- [ZKER-企业级开发规范手册_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-企业级开发规范手册_v1.0.md)
- [ZKER-故障排查手册_v1.0.md](../../docs/企业级功能完善与统一性设计方案/ZKER-故障排查手册_v1.0.md)
- [Kubernetes探针文档](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)

---

## 🐛 常见问题

### Q1: 健康检查超时怎么办？

**A**: 增加超时时间或优化检查逻辑：

```go
// 为单个检查器设置超时
checker := health.NewHTTPChecker(
    "api",
    "https://api.example.com/health",
    30*time.Second,  // 增加超时时间
    client,
)
```

### Q2: 如何禁用某个健康检查器？

**A**: 注销检查器：

```go
healthSvc := health.GetGlobalHealthService()
healthSvc.UnregisterChecker("mysql")
```

### Q3: 健康检查影响性能怎么办？

**A**:
1. 减少检查频率（负载均衡器）
2. 优化检查逻辑（缓存查询结果）
3. 使用异步检查（不阻塞主线程）

### Q4: 如何添加自定义指标到Prometheus输出？

**A**: 修改`generatePrometheusMetrics`函数，添加自定义指标：

```go
output += fmt.Sprintf("\n# TYPE my_custom_metric gauge\n")
output += fmt.Sprintf("my_custom_metric{name=\"%s\"} %d\n", name, value)
```

---

**文档变更历史**:

| 版本 | 日期 | 变更内容 | 作者 |
|-----|------|---------|------|
| v1.0 | 2025-01-01 | 初始版本 | 后端架构团队 |

---

**© 2025 ZKER Project. All rights reserved.**
