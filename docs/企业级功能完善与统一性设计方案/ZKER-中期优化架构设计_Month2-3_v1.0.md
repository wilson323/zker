# ZKER 中期优化架构设计 (Month 2-3)

> **文档版本**: v1.0
> **制定日期**: 2025-01-01
> **执行周期**: Month 2-3 (8周)
> **负责人**: 研发B + 研发A + 研发D

---

## 📋 目录

1. [执行摘要](#执行摘要)
2. [微服务拆分](#微服务拆分)
3. [消息队列异步化](#消息队列异步化)
4. [CDN加速](#cdn加速)
5. [实施计划](#实施计划)

---

## 执行摘要

本方案涵盖Month 2-3的中期优化措施，进一步提升系统的可扩展性、性能和可靠性：

| 优化项目 | 预期收益 | 优先级 | 工作量 |
|---------|---------|--------|--------|
| 微服务拆分 | 独立部署/扩容，故障隔离 | P0 | 3周 |
| 消息队列异步化 | API响应-50%，吞吐量+100% | P0 | 2周 |
| CDN加速 | 资源加载<100ms，带宽成本-60% | P1 | 1周 |

**总工作量**: 6周 (Month 2-3)

---

## 微服务拆分

### 1. 当前架构

```
┌─────────────┐
│  Monolith   │
│  (单体应用)   │
│             │
│ - Tenant    │
│ - Quota     │
│ - Subs      │
│ - Auth      │
└─────────────┘
```

**问题**:
- 代码耦合严重
- 无法独立部署和扩容
- 故障影响范围大
- 技术栈无法独立演进

### 2. 目标架构

```
                    ┌──────────────┐
                    │ API Gateway  │
                    │  (Kong)      │
                    └──────┬───────┘
                           │
      ┌────────────────────┼────────────────────┐
      │                    │                    │
┌─────▼─────┐      ┌───────▼──────┐    ┌──────▼─────┐
│   Tenant   │      │    Quota     │    │   Auth     │
│  Service   │      │   Service    │    │  Service   │
│            │      │              │    │            │
│ - CRUD     │      │ - Check      │    │ - JWT      │
│ - List     │      │ - Update     │    │ - RBAC     │
│ - Mgmt     │      │ - Mgmt       │    │ - SSO      │
└─────┬──────┘      └──────┬───────┘    └──────┬─────┘
      │                    │                    │
      └────────────────────┼────────────────────┘
                           │
                    ┌──────▼──────┐
                    │  Shared DB  │
                    │  (MySQL)    │
                    └─────────────┘
```

### 3. 服务边界定义

#### 3.1 租户服务 (Tenant Service)

**职责**:
- 租户CRUD操作
- 租户列表查询
- 租户状态管理
- 租户配额关联

**API端点**:
```
POST   /api/v1/tenants              # 创建租户
GET    /api/v1/tenants              # 租户列表
GET    /api/v1/tenants/:id          # 租户详情
PUT    /api/v1/tenants/:id          # 更新租户
DELETE /api/v1/tenants/:id          # 删除租户
PATCH  /api/v1/tenants/:id/status   # 更新状态
```

**数据库**:
```sql
-- 独立Schema
tenant.tenants
tenant.subscriptions
```

**技术栈**:
- Go + Hertz
- GORM
- MySQL (独立数据库)

#### 3.2 配额服务 (Quota Service)

**职责**:
- 配额检查
- 配额使用更新
- 配额限制管理
- 配额统计报表

**API端点**:
```
POST   /api/v1/quota/check          # 检查配额
POST   /api/v1/quota/update         # 更新使用量
GET    /api/v1/quota/usage          # 查询使用量
GET    /api/v1/quota/limit          # 查询限制
```

**数据库**:
```sql
quota.quotas
quota.quota_usage
quota.quota_history
```

**技术栈**:
- Go + Hertz
- GORM
- Redis (配额缓存)

#### 3.3 认证服务 (Auth Service)

**职责**:
- 用户认证 (JWT)
- 权限验证 (RBAC)
- 会话管理
- 单点登录 (SSO)

**API端点**:
```
POST   /api/v1/auth/login           # 用户登录
POST   /api/v1/auth/logout          # 用户登出
POST   /api/v1/auth/refresh         # 刷新Token
GET    /api/v1/auth/permissions    # 查询权限
POST   /api/v1/auth/verify          # 验证Token
```

**数据库**:
```sql
auth.users
auth.roles
auth.permissions
auth.user_roles
auth.sessions
```

**技术栈**:
- Go + Hertz
- JWT (RSA签名)
- Redis (会话缓存)

### 4. 服务间通信

#### 4.1 同步通信 (HTTP/REST)

**适用场景**:
- 需要实时响应的查询
- 强一致性要求的操作

**示例**:
```go
// 租户服务调用配额服务检查配额
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) error {
    // 1. 检查租户配额
    ok, err := s.quotaClient.CheckQuota(ctx, "tenant", req.TenantID)
    if err != nil {
        return err
    }
    if !ok {
        return errors.New("quota exceeded")
    }

    // 2. 创建租户
    return s.repo.Create(ctx, req)
}
```

#### 4.2 异步通信 (消息队列)

**适用场景**:
- 非核心业务操作
- 允许最终一致性
- 高并发场景

**示例**:
```go
// 发布租户创建事件
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) error {
    // 1. 创建租户
    if err := s.repo.Create(ctx, req); err != nil {
        return err
    }

    // 2. 发布事件
    event := TenantCreatedEvent{
        TenantID: req.TenantID,
        Name:     req.Name,
        CreatedAt: time.Now(),
    }
    s.pub.Publish("tenant.created", event)

    return nil
}

// 配额服务订阅事件
func (s *QuotaService) OnTenantCreated(event TenantCreatedEvent) {
    // 为新租户初始化配额
    s.InitQuota(context.Background(), event.TenantID)
}
```

### 5. 数据一致性

#### 5.1 分布式事务

**Saga模式**:
```go
// 创建租户 + 初始化配额的分布式事务
func (s *TenantService) CreateTenantWithQuota(ctx context.Context, req *CreateTenantRequest) error {
    saga := saga.NewSaga("create-tenant-with-quota")

    // 步骤1: 创建租户
    saga.AddStep(
        func() error {
            return s.repo.Create(ctx, req)
        },
        func() error {
            return s.repo.Delete(ctx, req.TenantID)  // 补偿操作
        },
    )

    // 步骤2: 初始化配额
    saga.AddStep(
        func() error {
            return s.quotaClient.InitQuota(ctx, req.TenantID)
        },
        func() error {
            return s.quotaClient.RemoveQuota(ctx, req.TenantID)  // 补偿操作
        },
    )

    return saga.Execute(ctx)
}
```

### 6. 服务发现与注册

#### 6.1 Consul注册中心

```go
// 服务注册
func RegisterService(consulAddr, serviceName, serviceAddr string) error {
    config := api.DefaultConfig()
    config.Address = consulAddr

    client, err := api.NewClient(config)
    if err != nil {
        return err
    }

    registration := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%s", serviceName, serviceAddr),
        Name:    serviceName,
        Address: serviceAddr,
        Port:    8080,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s/health", serviceAddr),
            Interval:                       "10s",
            Timeout:                        "3s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }

    return client.Agent().ServiceRegister(registration)
}

// 服务发现
func DiscoverService(consulAddr, serviceName string) (string, error) {
    config := api.DefaultConfig()
    config.Address = consulAddr

    client, err := api.NewClient(config)
    if err != nil {
        return "", err
    }

    services, _, err := client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return "", err
    }

    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances found for %s", serviceName)
    }

    // 随机选择一个实例
    service := services[rand.Intn(len(services))]
    return fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port), nil
}
```

---

## 消息队列异步化

### 1. 架构设计

```
┌──────────┐         ┌──────────────┐         ┌──────────┐
│ Producer │────────▶│    NSQ       │────────▶│ Consumer │
│ (API)    │         │   (Topic)    │         │ (Worker) │
└──────────┘         └──────────────┘         └──────────┘
                            │
                            ├─ topic: tenant-events
                            ├─ topic: quota-events
                            └─ topic: subscription-events
```

### 2. NSQ集群配置

```yaml
# docker/nsq.yml
version: '3.8'

services:
  nsqlookupd:
    image: nsqio/nsq:v1.2.1
    command: /nsqlookupd
    ports:
      - "4160:4160"
      - "4161:4161"

  nsqd:
    image: nsqio/nsq:v1.2.1
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    depends_on:
      - nsqlookupd
    ports:
      - "4150:4150"
      - "4151:4151"

  nsqadmin:
    image: nsqio/nsq:v1.2.1
    command: /nsqadmin --lookupd-http-address=nsqlookupd:4161
    ports:
      - "4171:4171"
```

### 3. 消息定义

```go
// 租户事件
type TenantEvent struct {
    EventType string    `json:"event_type"` // created, updated, deleted
    TenantID  string    `json:"tenant_id"`
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
}

// 配额事件
type QuotaEvent struct {
    EventType    string  `json:"event_type"` // checked, updated, exceeded
    TenantID     string  `json:"tenant_id"`
    ResourceType string  `json:"resource_type"`
    Amount       int     `json:"amount"`
    CurrentUsage int     `json:"current_usage"`
    Limit        int     `json:"limit"`
    Timestamp    time.Time `json:"timestamp"`
}
```

### 4. Producer实现

```go
// 发布消息到NSQ
func (p *Producer) PublishTenantEvent(event TenantEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    err = p.producer.Publish("tenant-events", data)
    if err != nil {
        return fmt.Errorf("failed to publish tenant event: %w", err)
    }

    return nil
}
```

### 5. Consumer实现

```go
// 处理租户事件
func (c *Consumer) HandleTenant(message *nsq.Message) error {
    var event TenantEvent
    if err := json.Unmarshal(message.Body, &event); err != nil {
        return err
    }

    switch event.EventType {
    case "created":
        // 处理租户创建事件
        return c.onTenantCreated(event)
    case "updated":
        // 处理租户更新事件
        return c.onTenantUpdated(event)
    case "deleted":
        // 处理租户删除事件
        return c.onTenantDeleted(event)
    default:
        return fmt.Errorf("unknown event type: %s", event.EventType)
    }
}
```

---

## CDN加速

### 1. 静态资源CDN

**资源类型**:
- 前端JS/CSS文件
- 图片和视频
- 字体文件
- 下载文件

**CDN配置**:
```nginx
# Nginx配置
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
    add_header X-Content-Type-Options nosniff;

    # CDN回源
    proxy_pass https://cdn.zker.com;
    proxy_set_header Host cdn.zker.com;
}
```

### 2. API响应缓存

**Redis缓存策略**:
```go
// 缓存租户详情 (5分钟)
func (s *TenantService) GetTenant(ctx context.Context, id string) (*Tenant, error) {
    // 1. 查询缓存
    cacheKey := fmt.Sprintf("tenant:%s", id)
    cached, err := s.redis.Get(ctx, cacheKey)
    if err == nil {
        var tenant Tenant
        json.Unmarshal([]byte(cached), &tenant)
        return &tenant, nil
    }

    // 2. 查询数据库
    tenant, err := s.repo.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. 写入缓存
    data, _ := json.Marshal(tenant)
    s.redis.Set(ctx, cacheKey, data, 5*time.Minute)

    return tenant, nil
}
```

### 3. CDN预热

```bash
#!/bin/bash
# CDN预热脚本

# 预热静态资源
urls=(
    "https://cdn.zker.com/static/js/main.js"
    "https://cdn.zker.com/static/css/main.css"
    "https://cdn.zker.com/images/logo.png"
)

for url in "${urls[@]}"; do
    echo "Warming up: $url"
    curl -X PURGE "$url"
done
```

---

## 实施计划

### Month 2 (Week 1-4)

| 周次 | 任务 | 负责人 | 工作量 |
|------|------|--------|--------|
| Week 1 | 服务边界定义 + 接口设计 | 研发A + 研发B | 1周 |
| Week 2-3 | 租户服务拆分 + 独立部署 | 研发B + 研发D | 2周 |
| Week 4 | 配额服务拆分 + 独立部署 | 研发B + 研发D | 1周 |

### Month 3 (Week 5-8)

| 周次 | 任务 | 负责人 | 工作量 |
|------|------|--------|--------|
| Week 5 | 认证服务拆分 + JWT集成 | 研发A + 研发B | 1周 |
| Week 6-7 | NSQ消息队列集成 + 异步化 | 研发B + 研发D | 2周 |
| Week 8 | CDN配置 + 性能优化 | 研发B + 研发D | 1周 |

---

**制定人**: 研发B - 后端工程师
**审核人**: 技术负责人
**执行周期**: Month 2-3
