# ZKER 分布式追踪配置

## OpenTelemetry + Jaeger

### 安装依赖

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/exporters/jaeger
go get go.opentelemetry.io/otel/semconv/v1.4.0
```

### 初始化追踪

```go
import (
    "context"
    "github.com/coze-studio/coze-studio/backend/infra/tracing"
)

func main() {
    // 初始化追踪
    err := tracing.Init(
        "zker-api",                    // 服务名称
        "http://localhost:14268/api/traces", // Jaeger endpoint
    )
    if err != nil {
        log.Fatal(err)
    }
    defer tracing.Shutdown(context.Background())

    // ... 启动服务
}
```

### 使用追踪

#### 1. HTTP请求追踪

```go
import (
    "github.com/cloudwego/hertz/pkg/app"
    tracing "github.com/coze-studio/coze-studio/backend/infra/tracing"
)

// 在路由中添加追踪中间件
r.Use(tracing.TracingMiddleware("zker-api"))
```

#### 2. 数据库查询追踪

```go
func (r *TenantRepository) GetByID(ctx context.Context, id string) (*Tenant, error) {
    var tenant Tenant

    err := tracing.TraceDBQuery(
        ctx,
        "mysql",        // dbSystem
        "zker",         // dbName
        "tenants",      // table
        "select",       // operation
        func() error {
            return r.db.WithContext(ctx).Where("tenant_id = ?", id).First(&tenant).Error
        },
    )

    return &tenant, err
}
```

#### 3. 缓存操作追踪

```go
func (s *CacheService) Get(ctx context.Context, key string) (interface{}, error) {
    return tracing.TraceCacheOperation(
        ctx,
        "redis",    // cacheType
        "get",      // operation
        key,        // key
        func() (interface{}, error) {
            return s.redis.Get(ctx, key)
        },
    )
}
```

#### 4. 业务操作追踪

```go
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*Tenant, error) {
    var createdTenant *Tenant

    err := tracing.TraceTenantOperation(
        ctx,
        "create_tenant", // operation
        req.TenantID,    // tenantID
        func() error {
            var err error
            createdTenant, err = s.repo.Create(ctx, req)
            return err
        },
    )

    return createdTenant, err
}
```

### 查看追踪数据

1. 启动Jaeger:
```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  jaegertracing/all-in-one:latest
```

2. 访问Jaeger UI:
```
http://localhost:16686
```

3. 查询追踪:
- 服务名: zker-api
- 操作名: tenant.create, db.query, cache.get
- Trace ID: 从日志中获取

### 采样率配置

开发环境: 100% 采样率
```go
tracing.InitWithSampler("zker-api", jaegerEndpoint, 1.0)
```

生产环境: 10% 采样率
```go
tracing.InitWithSampler("zker-api", jaegerEndpoint, 0.1)
```

### 环境变量

```bash
# OpenTelemetry
OTEL_SERVICE_NAME=zker-api
OTEL_EXPORTER_JAEGER_ENDPOINT=http://localhost:14268/api/traces
OTEL_TRACES_SAMPLER=0.1  # 10%采样率

# Jaeger
JAEGER_AGENT_HOST=localhost
JAEGER_AGENT_PORT=6831
JAEGER_SAMPLER_TYPE=probabilistic
JAEGER_SAMPLER_PARAM=0.1
```
