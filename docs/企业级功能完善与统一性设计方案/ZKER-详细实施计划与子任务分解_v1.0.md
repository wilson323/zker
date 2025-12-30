# ZKER 企业级功能完善 - 详细实施计划与子任务分解

**版本**: v1.0
**创建日期**: 2025-01-01
**项目周期**: 8周（4人并行开发）
**文档状态**: 正式版

---

## 📊 **执行摘要**

### 总体进度

| 阶段 | 状态 | 工作量 | 完成时间 |
|------|------|--------|----------|
| **P0任务** | ✅ 已完成 | 15人天 | Day 1-3 |
| **P1任务** | 🔄 待执行 | 38人天 | Day 4-15 |
| **P2任务** | ⏳ 待执行 | 23人天 | Day 16-30 |
| **总计** | 进行中 | 76人天 | 30个工作日 |

### 人员分配

| 角色 | 负责领域 | 主要任务 |
|------|----------|----------|
| **研发A（后端架构师）** | 多租户、权限、路由 | P0-1✅、P1-2、P1-7、P2-1 |
| **研发B（后端工程师）** | 错误码、性能、监控 | P0-2✅、P1-1、P1-8、P1-9、P2-2、P2-3 |
| **研发C（前端工程师）** | 组件库、业务页面 | P1-3、P2-4 |
| **研发D（DevOps）** | Docker、K8s、CI/CD | P1-4、P1-10、P2-5、P2-6 |

---

## ✅ **P0任务：核心基础设施（已完成）**

### P0-1: tenant_id 迁移（✅ 已完成）

**工作量**: 8人天 | **完成时间**: Day 1-2 | **负责人**: 研发A

#### 完成情况
- ✅ 26张表成功添加 `tenant_id` 字段
- ✅ 创建Go迁移执行器（支持回滚）
- ✅ 创建Bash迁移包装脚本
- ✅ 数据库备份和验证机制
- ✅ 空数据库验证通过

#### 交付物
1. `backend/domain/tenant/migration/tenant_id_migration.go`
2. `scripts/migrate_tenant_id.sh`
3. `backend/domain/tenant/migration/02_backfill_tenant_id_data_simple.sql`
4. 数据库备份文件

---

### P0-2: 配额中间件实现（✅ 已完成）

**工作量**: 4人天 | **完成时间**: Day 2-3 | **负责人**: 研发B

#### 完成情况
- ✅ QuotaRepository 接口和实现
- ✅ QuotaService（检查、消费、回滚）
- ✅ QuotaCheck 中间件（3种模式）
- ✅ QuotaMonitoringService（监控告警）
- ✅ QuotaInitializer（初始化服务）

#### 交付物
1. `backend/domain/tenant/repository/quota_repository.go`
2. `backend/domain/tenant/repository/quota_repository_impl.go`
3. `backend/domain/tenant/service/quota_monitoring.go`
4. `backend/domain/tenant/service/quota_init.go`

---

### P0-3: 权限中间件实现（✅ 已完成）

**工作量**: 3人天 | **完成时间**: Day 3 | **负责人**: 研发A

#### 完成情况
- ✅ Permission 接口定义
- ✅ AuthzChecker 授权检查器
- ✅ PermissionCheck 中间件
- ✅ RBAC模型（角色、数据权限、字段权限）

#### 交付物
1. `backend/domain/permission/permission.go`
2. `backend/domain/permission/authz_checker.go`
3. `backend/api/middleware/permission_check.go`

---

## 🔄 **P1任务：核心功能完善（38人天）**

### P1-1: 统一错误码系统完善（5人天）

**工作量**: 5人天 | **计划时间**: Day 4-8 | **负责人**: 研发B

#### 子任务分解

##### 1.1 错误码体系完善（2人天）
- [ ] **1.1.1** 审查现有错误码定义
  - 检查 `backend/types/errno/` 下所有错误码文件
  - 验证错误码分类和编号规则
  - 输出：错误码审查报告

- [ ] **1.1.2** 补充缺失的错误码
  - 租户模块错误码（10个）
  - 权限模块错误码（15个）
  - 路由模块错误码（10个）
  - 配置模块错误码（5个）
  - 输出：完整的错误码定义文件

- [ ] **1.1.3** 错误码国际化支持
  - 创建 `backend/pkg/errorx/i18n/` 目录
  - 实现中英文错误消息映射
  - 输出：i18n 错误消息文件

##### 1.2 错误响应格式标准化（1人天）
- [ ] **1.2.1** 定义统一响应格式
  ```go
  type ErrorResponse struct {
      Code      int                    `json:"code"`
      Message   string                 `json:"message"`
      Details   map[string]interface{} `json:"details,omitempty"`
      RequestID string                 `json:"request_id"`
      Timestamp int64                  `json:"timestamp"`
  }
  ```

- [ ] **1.2.2** 创建响应构造工具
  - `pkg/errorx/response.go`
  - 支持多语言错误消息
  - 自动生成 request_id

- [ ] **1.2.3** 更新所有API错误处理
  - 审查所有 Handler 的错误返回
  - 统一使用新的响应格式

##### 1.3 错误码文档生成（1人天）
- [ ] **1.3.1** 自动生成错误码文档
  - 从代码注释提取错误码信息
  - 生成 Markdown 文档
  - 输出：`docs/错误码完整列表_v1.0.md`

- [ ] **1.3.2** 创建错误码使用指南
  - 如何添加新错误码
  - 如何正确处理错误
  - 最佳实践

##### 1.4 错误码测试（1人天）
- [ ] **1.4.1** 错误码单元测试
  - 测试所有错误码注册
  - 测试错误消息格式化
  - 测试i18n支持

- [ ] **1.4.2** 集成测试
  - 测试API错误响应格式
  - 测试错误码唯一性

#### 验收标准
- ✅ 所有模块都有完整的错误码定义（300+错误码）
- ✅ 错误响应格式统一，支持中英文
- ✅ 错误码文档完整且自动生成
- ✅ 单元测试覆盖率 ≥ 90%

#### 依赖关系
- 依赖于：P0-2（配额中间件已实现部分错误码）
- 被依赖：所有P1、P2任务

#### 风险评估
- 🟡 **中风险**: 错误码编号冲突
  - 缓解措施：建立错误码注册表，自动检测重复

---

### P1-2: 系统角色初始化（2人天）

**工作量**: 2人天 | **计划时间**: Day 4-5 | **负责人**: 研发A

#### 子任务分解

##### 2.1 系统角色定义（0.5人天）
- [ ] **2.1.1** 定义4个系统角色
  1. **TenantAdmin** (租户管理员)
     - 全部数据权限
     - 所有字段可编辑
     - 可管理用户和角色

  2. **Developer** (开发者)
     - 部门数据权限
     - 开发相关字段可编辑
     - 可创建和管理Bot

  3. **Operator** (运营人员)
     - 仅自己数据权限
     - 只读大部分字段
     - 可查看数据统计

  4. **Viewer** (查看者)
     - 仅自己数据权限
     - 只读所有字段
     - 无编辑权限

- [ ] **2.1.2** 创建角色初始化SQL
  - 文件：`backend/domain/permission/migration/01_init_system_roles.sql`
  - 包含4个系统角色的INSERT语句
  - 包含默认的数据权限和字段权限

##### 2.2 角色权限配置（1人天）
- [ ] **2.2.1** 为每个角色配置数据权限
  ```sql
  -- TenantAdmin: ALL for all resources
  -- Developer: DEPARTMENT for bots/workflows, OWN for others
  -- Operator: OWN for all resources
  -- Viewer: OWN for all resources (readonly)
  ```

- [ ] **2.2.2** 为每个角色配置字段权限
  ```sql
  -- 定义敏感字段：api_key, secret_key, webhook_url
  -- TenantAdmin: editable for all
  -- Developer: readonly for sensitive fields
  -- Operator/Viewer: hidden for sensitive fields
  ```

- [ ] **2.2.3** 创建Role初始化服务
  - 文件：`backend/domain/tenant/service/role_init.go`
  - 函数：`InitializeSystemRoles(tenantID string) error`
  - 确保每个租户创建时自动初始化4个角色

##### 2.3 默认用户分配（0.5人天）
- [ ] **2.3.1** 租户创建者自动分配TenantAdmin角色
- [ ] **2.3.2** 创建用户角色关联逻辑
- [ ] **2.3.3** 测试角色权限验证

#### 验收标准
- ✅ 4个系统角色SQL脚本可执行
- ✅ 每个角色都有明确的数据权限和字段权限
- ✅ 新租户自动初始化系统角色
- ✅ 租户创建者自动获得TenantAdmin角色

#### 依赖关系
- 依赖于：P0-3（权限中间件）
- 被依赖：P1-3（前端权限管理页面）

#### 风险评估
- 🟢 **低风险**: 角色定义清晰，逻辑简单

---

### P1-3: 前端TypeScript API客户端（6人天）

**工作量**: 6人天 | **计划时间**: Day 6-11 | **负责人**: 研发C

#### 子任务分解

##### 3.1 API客户端基础设施（2人天）
- [ ] **3.1.1** 创建API客户端基础架构
  - 目录：`frontend/packages/client/api/`
  - 文件结构：
    ```
    api/
    ├── base.ts            # 基础HTTP客户端
    ├── types/             # TypeScript类型定义
    ├── tenant/            # 租户API
    ├── permission/        # 权限API
    ├── quota/             # 配额API
    └── index.ts           # 统一导出
    ```

- [ ] **3.1.2** 实现HTTP请求拦截器
  - 自动添加 tenant_id header
  - 自动添加 authorization header
  - 统一错误处理
  - 请求/响应日志

- [ ] **3.1.3** 类型定义生成
  - 从后端API定义生成TypeScript类型
  - 使用 OpenAPI/Swagger 规范
  - 输出：`api/types/*.ts`

##### 3.2 租户API客户端（1人天）
- [ ] **3.2.1** 租户管理API
  ```typescript
  // tenant/api.ts
  export class TenantAPI {
      listTenants(params: ListTenantsRequest): Promise<ListTenantsResponse>
      getTenant(tenantId: string): Promise<Tenant>
      createTenant(data: CreateTenantRequest): Promise<Tenant>
      updateTenant(tenantId: string, data: UpdateTenantRequest): Promise<Tenant>
      deleteTenant(tenantId: string): Promise<void>
  }
  ```

- [ ] **3.2.2** 订阅管理API
  ```typescript
  export class SubscriptionAPI {
      getSubscription(tenantId: string): Promise<Subscription>
      upgradeSubscription(tenantId: string, tier: SubscriptionTier): Promise<Subscription>
      cancelSubscription(tenantId: string): Promise<void>
  }
  ```

- [ ] **3.2.3** 配额查询API
  ```typescript
  export class QuotaAPI {
      getQuotas(tenantId: string): Promise<Quota[]>
      getQuotaUsage(tenantId: string, resourceType: string): Promise<QuotaUsage>
  }
  ```

##### 3.3 权限API客户端（1人天）
- [ ] **3.3.1** 角色管理API
  ```typescript
  export class RoleAPI {
      listRoles(tenantId: string): Promise<Role[]>
      createRole(data: CreateRoleRequest): Promise<Role>
      updateRole(roleId: string, data: UpdateRoleRequest): Promise<Role>
      deleteRole(roleId: string): Promise<void>
      assignRole(userId: string, roleId: string): Promise<void>
      revokeRole(userId: string, roleId: string): Promise<void>
  }
  ```

- [ ] **3.3.2** 权限检查API
  ```typescript
  export class PermissionAPI {
      checkPermission(req: CheckPermissionRequest): Promise<CheckPermissionResponse>
      getUserPermissions(userId: string): Promise<Permission[]>
  }
  ```

##### 3.4 用户API客户端（1人天）
- [ ] **3.4.1** 用户管理API
- [ ] **3.4.2** 用户角色关联API
- [ ] **3.4.3** 用户部门关联API

##### 3.5 API客户端测试（1人天）
- [ ] **3.5.1** 单元测试
  - Mock HTTP响应
  - 测试所有API方法

- [ ] **3.5.2** 集成测试
  - 连接真实后端API
  - 测试错误处理

#### 验收标准
- ✅ TypeScript类型完整，无any类型
- ✅ 所有API都有对应的客户端方法
- ✅ 自动添加认证和租户信息
- ✅ 统一错误处理和类型定义
- ✅ 单元测试覆盖率 ≥ 80%

#### 依赖关系
- 依赖于：P1-1（错误码系统）
- 被依赖：P1-3（前端业务页面）

#### 风险评估
- 🟡 **中风险**: TypeScript类型定义可能不完整
  - 缓解措施：使用工具自动生成类型定义

---

### P1-4: 服务健康检查（3人天）

**工作量**: 3人天 | **计划时间**: Day 6-8 | **负责人**: 研发B

#### 子任务分解

##### 4.1 健康检查接口设计（0.5人天）
- [ ] **4.1.1** 定义健康检查标准
  ```go
  type HealthStatus string
  const (
      HealthStatusHealthy   HealthStatus = "healthy"
      HealthStatusDegraded  HealthStatus = "degraded"
      HealthStatusUnhealthy HealthStatus = "unhealthy"
  )

  type HealthCheckResult struct {
      Service    string                 `json:"service"`
      Status     HealthStatus           `json:"status"`
      Version    string                 `json:"version"`
      Timestamp  time.Time              `json:"timestamp"`
      Checks     map[string]CheckResult `json:"checks"`
  }
  ```

- [ ] **4.1.2** 设计健康检查端点
  - `GET /health` - 简要健康状态
  - `GET /health/detailed` - 详细健康信息
  - `GET /health/ready` - 就绪探针（K8s）
  - `GET /health/live` - 存活探针（K8s）

##### 4.2 核心服务健康检查（1人天）
- [ ] **4.2.1** MySQL健康检查
  ```go
  func CheckMySQL(ctx context.Context) CheckResult {
      // 检查连接
      // 检查响应时间
      // 检查连接池状态
  }
  ```

- [ ] **4.2.2** Redis健康检查
  ```go
  func CheckRedis(ctx context.Context) CheckResult {
      // 检查连接
      // 检查响应时间
      // 检查内存使用
  }
  ```

- [ ] **4.2.3** Elasticsearch健康检查
- [ ] **4.2.4** MinIO健康检查
- [ ] **4.2.5** Milvus健康检查

##### 4.3 应用级健康检查（1人天）
- [ ] **4.3.1** 检查租户服务状态
- [ ] **4.3.2** 检查权限服务状态
- [ ] **4.3.3** 检查配额服务状态
- [ ] **4.3.4** 检查路由服务状态
- [ ] **4.3.5** 检查队列积压情况

##### 4.4 健康检查中间件（0.5人天）
- [ ] **4.4.1** 创建健康检查Handler
  - 文件：`backend/api/handler/health.go`
  - 实现所有健康检查端点

- [ ] **4.4.2** 集成到主路由
  ```go
  r.GET("/health", handler.Health)
  r.GET("/health/detailed", handler.DetailedHealth)
  r.GET("/health/ready", handler.Ready)
  r.GET("/health/live", handler.Live)
  ```

#### 验收标准
- ✅ 所有外部依赖服务都有健康检查
- ✅ 响应时间 < 1秒
- ✅ 支持K8s探针
- ✅ 返回详细的健康信息

#### 依赖关系
- 依赖于：所有P0任务
- 被依赖：P1-10（监控告警）

#### 风险评估
- 🟢 **低风险**: 健康检查逻辑简单

---

### P1-5: 智能路由引擎增强（8人天）

**工作量**: 8人天 | **计划时间**: Day 9-16 | **负责人**: 研发A

#### 子任务分解

##### 5.1 路由规则引擎（3人天）
- [ ] **5.1.1** 意图识别增强
  ```go
  type IntentMatcher interface {
      Match(input string) (*Intent, error)
      Train(examples []TrainingExample) error
  }

  // 实现多种匹配策略
  - KeywordMatcher: 关键词匹配
  - RegexMatcher: 正则表达式匹配
  - SemanticMatcher: 语义相似度匹配（使用Embedding）
  ```

- [ ] **5.1.2** 路由规则优先级
  ```go
  type RoutingRule struct {
      RuleID      string
      TenantID    string
      Priority    int    // 优先级（数字越大越优先）
      Conditions  []Condition
      Actions     []Action
      Enabled     bool
  }
  ```

- [ ] **5.1.3** 规则匹配引擎
  - 多条件组合（AND/OR）
  - 通配符支持
  - 正则表达式支持

##### 5.2 评分路由机制（2人天）
- [ ] **5.2.1** 评分算法实现
  ```go
  type Scorer interface {
      Score(ctx context.Context, req *RoutingRequest) (float64, error)
  }

  // 多维度评分
  - CapabilityScorer: 能力匹配评分
  - PerformanceScorer: 性能评分
  - CostScorer: 成本评分
  - AvailabilityScorer: 可用性评分
  ```

- [ ] **5.2.2** 评分聚合策略
  - 加权平均
  - 最高分
  - 票选机制

##### 5.3 路由缓存优化（1人天）
- [ ] **5.3.1** 路由结果缓存
  - Redis缓存路由决策
  - TTL控制
  - 缓存失效策略

- [ ] **5.3.2** 路由规则缓存
  - 内存缓存规则
  - 热更新机制

##### 5.4 A/B测试路由（1人天）
- [ ] **5.4.1** A/B测试支持
  ```go
  type ABTestRule struct {
      TestID      string
      TrafficSplit map[string]int  // agent_id: percentage
      Metrics      []string        // 跟踪的指标
  }
  ```

- [ ] **5.4.2** 流量分配算法
- [ ] **5.4.3** 结果统计分析

##### 5.5 路由监控和分析（1人天）
- [ ] **5.5.1** 路由决策日志
  - 记录每次路由决策
  - 包含输入、规则、评分、结果

- [ ] **5.5.2** 路由效果分析
  - 成功率统计
  - 响应时间分析
  - 成本分析

#### 验收标准
- ✅ 支持3种意图匹配策略
- ✅ 支持4维评分路由
- ✅ 路由延迟 < 50ms
- ✅ 缓存命中率 > 80%
- ✅ 支持A/B测试

#### 依赖关系
- 依赖于：P0-1（tenant_id迁移）
- 被依赖：P2-1（前端路由配置页面）

#### 风险评估
- 🟡 **中风险**: 语义匹配准确度
  - 缓解措施：提供多种匹配策略，支持人工干预

---

### P1-6: 配置中心集成（4人天）

**工作量**: 4人天 | **计划时间**: Day 12-15 | **负责人**: 研发B

#### 子任务分解

##### 6.1 etcd集成（2人天）
- [ ] **6.1.1** etcd客户端封装
  ```go
  package config

  type EtcdClient struct {
      client *clientv3.Client
  }

  func NewEtcdClient(endpoints []string) (*EtcdClient, error)
  func (c *EtcdClient) Get(key string) (string, error)
  func (c *EtcdClient) Put(key, value string) error
  func (c *EtcdClient) Watch(key string) <-chan string
  ```

- [ ] **6.1.2** 配置热更新
  - Watch etcd key变化
  - 自动更新内存配置
  - 通知相关服务

##### 6.2 配置管理（1人天）
- [ ] **6.2.1** 配置结构定义
  ```go
  type AppConfig struct {
      Tenant    TenantConfig
      Permission PermissionConfig
      Quota     QuotaConfig
      Routing   RoutingConfig
  }
  ```

- [ ] **6.2.2** 配置验证
  - 启动时验证配置
  - 配置更新时验证
  - 提供配置修复建议

##### 6.3 配置API（1人天）
- [ ] **6.3.1** 配置查询API
- [ ] **6.3.2** 配置更新API
- [ ] **6.3.3** 配置历史版本

#### 验收标准
- ✅ 配置变更实时生效
- ✅ 配置变更无需重启服务
- ✅ 配置有完整的验证机制

#### 依赖关系
- 依赖于：P0任务
- 被依赖：所有需要动态配置的功能

#### 风险评估
- 🟢 **低风险**: etcd是成熟方案

---

### P1-7: 租户管理后台API（4人天）

**工作量**: 4人天 | **计划时间**: Day 12-15 | **负责人**: 研发A

#### 子任务分解

##### 7.1 租户CRUD API（1人天）
- [ ] **7.1.1** 创建租户
  ```go
  POST /api/tenants
  Request: {
      "tenant_name": "Acme Corp",
      "tenant_type": "enterprise",
      "billing_cycle": "monthly",
      ...
  }
  Response: Tenant + 初始配额 + 系统角色
  ```

- [ ] **7.1.2** 查询租户列表
- [ ] **7.1.3** 获取租户详情
- [ ] **7.1.4** 更新租户
- [ ] **7.1.5** 删除/停用租户

##### 7.2 订阅管理API（1人天）
- [ ] **7.2.1** 订阅升级/降级
- [ ] **7.2.2** 订阅续费
- [ ] **7.2.3** 订阅取消
- [ ] **7.2.4** 订阅历史查询

##### 7.3 配额管理API（1人天）
- [ ] **7.3.1** 查询配额使用情况
- [ ] **7.3.2** 调整配额限制
- [ ] **7.3.3** 配额使用统计
- [ ] **7.3.4** 配额告警配置

##### 7.4 租户数据隔离API（1人天）
- [ ] **7.4.1** 租户数据导出
- [ ] **7.4.2** 租户数据删除
- [ ] **7.4.3** 租户数据迁移

#### 验收标准
- ✅ API符合RESTful规范
- ✅ 有完整的权限检查
- ✅ 有完整的错误处理
- ✅ API文档完整

#### 依赖关系
- 依赖于：P0-1（tenant_id）、P0-2（配额）、P1-1（错误码）

#### 风险评估
- 🟢 **低风险**: 标准CRUD操作

---

### P1-8: 性能测试框架（6人天）

**工作量**: 6人天 | **计划时间**: Day 16-21 | **负责人**: 研发B

#### 子任务分解

##### 8.1 JMeter测试脚本（2人天）
- [ ] **8.1.1** 核心API测试脚本
  - 租户API（创建、查询、更新）
  - Bot API（创建、对话）
  - 知识库API（上传、查询）
  - 工作流API（创建、执行）

- [ ] **8.1.2** 场景测试脚本
  - 并发用户场景
  - 高负载场景
  - 长时间稳定性测试

- [ ] **8.1.3** 数据准备脚本
  - 生成测试数据
  - 清理测试数据

##### 8.2 K6测试脚本（2人天）
- [ ] **8.2.1** 轻量级性能测试
  ```javascript
  import http from 'k6/http';
  import { check, sleep } from 'k6';

  export let options = {
      vus: 100,
      duration: '30s',
  };

  export default function() {
      let res = http.get('https://api.example.com/bots');
      check(res, {
          'status is 200': (r) => r.status === 200,
          'response time < 500ms': (r) => r.timings.duration < 500,
      });
      sleep(1);
  }
  ```

- [ ] **8.2.2** 性能基准测试
- [ ] **8.2.3** 回归测试脚本

##### 8.3 性能基线建立（1人天）
- [ ] **8.3.1** 定义性能指标
  - 响应时间（P50, P95, P99）
  - 吞吐量（QPS）
  - 错误率
  - 资源使用率

- [ ] **8.3.2** 执行基线测试
- [ ] **8.3.3** 生成性能报告
- [ ] **8.3.4** 建立性能基线

##### 8.4 性能监控系统（1人天）
- [ ] **8.4.1** 测试结果收集
  - Prometheus指标导出
  - 测试报告生成

- [ ] **8.4.2** 性能趋势分析
  - 历史数据对比
  - 性能退化检测

#### 验收标准
- ✅ 覆盖所有核心API
- ✅ 建立性能基线
- ✅ 自动化测试脚本
- ✅ 完整的性能报告

#### 依赖关系
- 依赖于：所有P0、P1核心功能

#### 风险评估
- 🟡 **中风险**: 测试数据准备可能复杂
  - 缓解措施：编写数据生成脚本

---

### P1-9: 日志收集和分析（4人天）

**工作量**: 4人天 | **计划时间**: Day 18-21 | **负责人**: 研发B

#### 子任务分解

##### 9.1 日志规范化（1人天）
- [ ] **9.1.1** 定义日志格式
  ```go
  type LogEntry struct {
      Level       string                 `json:"level"`
      Timestamp   string                 `json:"timestamp"`
      Service     string                 `json:"service"`
      TenantID    string                 `json:"tenant_id,omitempty"`
      UserID      string                 `json:"user_id,omitempty"`
      RequestID   string                 `json:"request_id"`
      Message     string                 `json:"message"`
      Fields      map[string]interface{} `json:"fields"`
  }
  ```

- [ ] **9.1.2** 日志级别规范
  - ERROR: 错误（需要立即处理）
  - WARN: 警告（需要关注）
  - INFO: 重要信息
  - DEBUG: 调试信息

- [ ] **9.1.3** 审查现有日志
  - 统一日志格式
  - 添加必要的上下文

##### 9.2 ELK集成（2人天）
- [ ] **9.2.1** Filebeat配置
  ```yaml
  filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /var/log/zker/*.log
    json.keys_under_root: true
    json.add_error_key: true

  output.elasticsearch:
    hosts: ["elasticsearch:9200"]
    index: "zker-logs-%{+yyyy.MM.dd}"
  ```

- [ ] **9.2.2** Logstash配置
  - 日志解析
  - 日志过滤
  - 日志增强

- [ ] **9.2.3** Kibana仪表盘
  - 创建日志索引模式
  - 创建日志查询仪表盘
  - 创建错误监控仪表盘

##### 9.3 日志查询API（1人天）
- [ ] **9.3.1** 日志查询接口
  ```go
  GET /api/logs/query
  Query: {
      "start_time": "...",
      "end_time": "...",
      "level": "ERROR",
      "tenant_id": "...",
      "query": "..."
  }
  ```

- [ ] **9.3.2** 日志统计API
- [ ] **9.3.3** 错误日志聚合

#### 验收标准
- ✅ 所有日志统一格式
- ✅ 日志可实时查询
- ✅ Kibana仪表盘完整
- ✅ 日志查询API可用

#### 依赖关系
- 依赖于：P1-10（监控告警）

#### 风险评估
- 🟡 **中风险**: ELK资源消耗大
  - 缓解措施：合理配置日志保留期

---

### P1-10: Prometheus + Grafana监控（4人天）

**工作量**: 4人天 | **计划时间**: Day 18-21 | **负责人**: 研发B

#### 子任务分解

##### 10.1 Prometheus指标定义（1人天）
- [ ] **10.1.1** 业务指标
  ```go
  var (
      tenantCount = prometheus.NewGaugeVec(
          prometheus.GaugeOpts{
              Name: "zker_tenants_total",
              Help: "Total number of tenants",
          },
          []string{"status"},
      )

      quotaUsage = prometheus.NewGaugeVec(
          prometheus.GaugeOpts{
              Name: "zker_quota_usage_percent",
              Help: "Quota usage percentage",
          },
          []string{"tenant_id", "resource_type"},
      )
  )
  ```

- [ ] **10.1.2** HTTP指标
  - 请求总数
  - 响应时间（P50, P95, P99）
  - 错误率

- [ ] **10.1.3** 数据库指标
  - 连接池使用率
  - 查询响应时间
  - 慢查询数量

##### 10.2 Prometheus配置（1人天）
- [ ] **10.2.1** Prometheus配置文件
  ```yaml
  scrape_configs:
    - job_name: 'zker-backend'
      static_configs:
        - targets: ['backend:8080']
      metrics_path: '/metrics'
      scrape_interval: 15s
  ```

- [ ] **10.2.2** 告警规则
  ```yaml
  groups:
  - name: zker_alerts
    rules:
    - alert: HighErrorRate
      expr: rate(zerk_http_errors_total[5m]) > 0.05
      for: 5m
      annotations:
        summary: "High error rate detected"
    - alert: QuotaExceeded
      expr: zker_quota_usage_percent > 90
      annotations:
        summary: "Quota usage exceeded 90%"
  ```

##### 10.3 Grafana仪表盘（1.天）
- [ ] **10.3.1** 系统概览仪表盘
  - QPS
  - 响应时间
  - 错误率
  - 租户数量

- [ ] **10.3.2** 租户详情仪表盘
  - 租户配额使用
  - API调用统计
  - 资源使用情况

- [ ] **10.3.3** 数据库监控仪表盘
- [ ] **10.3.4** 告警管理仪表盘

##### 10.4 告警通知（1人天）
- [ ] **10.4.1** AlertManager配置
  - 邮件通知
  - 钉钉通知
  - Slack通知（可选）

- [ ] **10.4.2** 告警分级
  - P0: 立即处理（电话+短信）
  - P1: 尽快处理（邮件+钉钉）
  - P2: 关注即可（邮件）

#### 验收标准
- ✅ 所有关键指标都有Prometheus指标
- ✅ Grafana仪表盘完整且美观
- ✅ 告警规则完整
- ✅ 告警通知正常工作

#### 依赖关系
- 依赖于：P1-4（健康检查）

#### 风险评估
- 🟢 **低风险**: Prometheus是成熟方案

---

## ⏳ **P2任务：前端和运维（23人天）**

### P2-1: 租户管理前端页面（6人天）

**工作量**: 6人天 | **计划时间**: Day 22-27 | **负责人**: 研发C

#### 子任务分解

##### 1.1 租户列表页面（1人天）
- [ ] **1.1.1** 页面布局
  - 表格展示租户列表
  - 搜索和筛选
  - 分页

- [ ] **1.1.2** 列表数据
  - 租户名称、类型、状态
  - 订阅等级、到期时间
  - 配额使用情况

##### 1.2 租户详情页面（1人天）
- [ ] **1.2.1** 基本信息展示
- [ ] **1.2.2** 订阅信息展示
- [ ] **1.2.3** 配额使用展示
- [ ] **1.2.4** 操作日志

##### 1.3 租户创建/编辑页面（1人天）
- [ ] **1.3.1** 表单设计
  - 租户名称、类型
  - 订阅等级、计费周期
  - 初始配额设置

- [ ] **1.3.2** 表单验证
- [ ] **1.3.3** 提交和确认

##### 1.4 配额管理页面（2人天）
- [ ] **1.4.1** 配额概览
  - 所有租户的配额使用情况
  - 图表展示（趋势图、饼图）

- [ ] **1.4.2** 配额调整
  - 调整配额限制
  - 配额变更历史

- [ ] **1.4.3** 配额告警
  - 告警规则配置
  - 告警历史

##### 1.5 权限管理页面（1人天）
- [ ] **1.5.1** 角色列表
- [ ] **1.5.2** 角色创建/编辑
- [ ] **1.5.3** 用户角色分配

#### 验收标准
- ✅ UI符合设计规范
- ✅ 所有API集成完成
- ✅ 表单验证完整
- ✅ 响应式设计

#### 依赖关系
- 依赖于：P1-3（TypeScript API客户端）、P1-7（租户管理API）

#### 风险评估
- 🟢 **低风险**: 标准CRUD页面

---

### P2-2: 权限管理前端页面（4人天）

**工作量**: 4人天 | **计划时间**: Day 24-27 | **负责人**: 研发C

#### 子任务分解

##### 2.1 角色管理页面（1.5人天）
- [ ] **2.1.1** 角色列表
  - 系统角色、自定义角色
  - 角色权限概览

- [ ] **2.1.2** 角色创建/编辑
  - 角色基本信息
  - 数据权限配置
  - 字段权限配置

##### 2.2 用户管理页面（1.5人天）
- [ ] **2.2.1** 用户列表
- [ ] **2.2.2** 用户创建/编辑
- [ ] **2.2.3** 角色分配
- [ ] **2.2.4** 部门分配

##### 2.3 权限测试页面（1人天）
- [ ] **2.3.1** 权限检查工具
  - 输入用户、资源、操作
  - 显示权限检查结果

- [ ] **2.3.2** 权限可视化
  - 权限矩阵展示
  - 权限依赖关系

#### 验收标准
- ✅ UI直观易用
- ✅ 权限配置清晰
- ✅ 支持批量操作

#### 依赖关系
- 依赖于：P1-3（TypeScript API客户端）

#### 风险评估
- 🟡 **中风险**: 权限配置复杂
  - 缓解措施：提供权限模板和向导

---

### P2-3: 路由配置前端页面（4人天）

**工作量**: 4人天 | **计划时间**: Day 26-29 | **负责人**: 研发C

#### 子任务分解

##### 3.1 路由规则列表（1人天）
- [ ] **3.1.1** 规则列表展示
  - 规则名称、优先级
  - 匹配条件
  - 路由目标
  - 启用状态

- [ ] **3.1.2** 搜索和筛选
- [ ] **3.1.3** 批量操作

##### 3.2 路由规则编辑器（2人天）
- [ ] **3.2.1** 条件配置
  - 关键词匹配
  - 正则表达式匹配
  - 语义相似度匹配

- [ ] **3.2.2** 评分配置
  - 评分维度
  - 权重设置

- [ ] **3.2.3** 规则测试
  - 输入测试文本
  - 显示路由结果
  - A/B测试支持

##### 3.3 路由分析页面（1人天）
- [ ] **3.3.1** 路由统计
  - 路由次数
  - 路由成功率
  - 平均响应时间

- [ ] **3.3.2** 可视化图表
  - 路由趋势图
  - Agent使用分布

#### 验收标准
- ✅ 配置界面直观
- ✅ 实时预览路由结果
- ✅ 支持规则导入导出

#### 依赖关系
- 依赖于：P1-5（智能路由引擎）

#### 风险评估
- 🟡 **中风险**: 路由规则配置复杂
  - 缓解措施：提供规则模板和测试工具

---

### P2-4: Docker和K8s配置完善（3人天）

**工作量**: 3人天 | **计划时间**: Day 22-24 | **负责人**: 研发D

#### 子任务分解

##### 4.1 Docker优化（1人天）
- [ ] **4.1.1** 多阶段构建
  ```dockerfile
  # 构建阶段
  FROM golang:1.24 AS builder
  WORKDIR /app
  COPY . .
  RUN go build -o zker-backend ./backend

  # 运行阶段
  FROM alpine:latest
  COPY --from=builder /app/zker-backend /usr/local/bin/
  ```

- [ ] **4.1.2** 镜像优化
  - 减小镜像体积
  - 安全扫描

##### 4.2 Kubernetes配置（1人天）
- [ ] **4.2.1** Deployment配置
  ```yaml
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: zker-backend
  spec:
    replicas: 3
    selector:
      matchLabels:
        app: zker-backend
    template:
      metadata:
        labels:
          app: zker-backend
      spec:
        containers:
        - name: backend
          image: zker-backend:latest
          ports:
          - containerPort: 8080
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
  ```

- [ ] **4.2.2** Service配置
- [ ] **4.2.3** ConfigMap和Secret配置
- [ ] **4.2.4** HPA配置（自动扩缩容）

##### 4.3 Helm Chart（1人天）
- [ ] **4.3.1** 创建Helm Chart
  ```
  helm/zker/
  ├── Chart.yaml
  ├── values.yaml
  ├── templates/
  │   ├── deployment.yaml
  │   ├── service.yaml
  │   ├── ingress.yaml
  │   └── configmap.yaml
  └── README.md
  ```

- [ ] **4.3.2** 编写Helm Chart文档
- [ ] **4.3.3** 测试Helm安装和升级

#### 验收标准
- ✅ Docker镜像 < 200MB
- ✅ K8s配置完整且可部署
- ✅ Helm Chart可一键部署
- ✅ 健康检查正常工作

#### 依赖关系
- 依赖于：P1-4（健康检查）

#### 风险评估
- 🟢 **低风险**: Docker和K8s是成熟技术

---

### P2-5: CI/CD流水线（3人天）

**工作量**: 3人天 | **计划时间**: Day 25-27 | **负责人**: 研发D

#### 子任务分解

##### 5.1 GitHub Actions配置（1.5人天）
- [ ] **5.1.1** CI流水线
  ```yaml
  name: CI

  on:
    push:
      branches: [ main, develop ]
    pull_request:
      branches: [ main ]

  jobs:
    test:
      runs-on: ubuntu-latest
      steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.24'
      - name: Run tests
        run: |
          cd backend
          go test ./... -cover
      - name: Run linter
        run: |
          cd backend
          golangci-lint run
  ```

- [ ] **5.1.2** CD流水线
  ```yaml
  name: CD

  on:
    push:
      tags:
        - 'v*'

  jobs:
    build-and-deploy:
      runs-on: ubuntu-latest
      steps:
      - uses: actions/checkout@v3
      - name: Build Docker image
        run: docker build -t zker-backend:${{ github.ref_name }} .
      - name: Push to Registry
        run: docker push zker-backend:${{ github.ref_name }}
      - name: Deploy to K8s
        run: |
          kubectl set image deployment/zker-backend \
            backend=zker-backend:${{ github.ref_name }}
  ```

##### 5.2 自动化测试集成（1人天）
- [ ] **5.2.1** 单元测试自动化
- [ ] **5.2.2** 集成测试自动化
- [ ] **5.2.3** 性能测试自动化（触发式）

##### 5.3 代码质量检查（0.5人天）
- [ ] **5.3.1** SonarQube集成
- [ ] **5.3.2** 代码覆盖率检查
- [ ] **5.3.3** 安全扫描

#### 验收标准
- ✅ CI/CD完全自动化
- ✅ 代码合并自动触发CI
- ✅ 测试通过自动部署
- ✅ 代码质量门禁完整

#### 依赖关系
- 依赖于：P1-8（性能测试）、P2-4（Docker/K8s）

#### 风险评估
- 🟡 **中风险**: CI/CD配置复杂
  - 缓解措施：使用成熟模板，分阶段实施

---

### P2-6: 灰度发布（3人天）

**工作量**: 3人天 | **计划时间**: Day 28-30 | **负责人**: 研发D

#### 子任务分解

##### 6.1 Istio配置（1.5人天）
- [ ] **6.1.1** 安装Istio
  ```bash
  istioctl install --set profile=demo
  ```

- [ ] **6.1.2** VirtualService配置
  ```yaml
  apiVersion: networking.istio.io/v1beta1
  kind: VirtualService
  metadata:
    name: zker-backend
  spec:
    hosts:
    - zker.example.com
    http:
    - match:
      - headers:
          x-canary:
            exact: "true"
      route:
      - destination:
          host: zker-backend
          subset: v2  # 新版本
        weight: 100
    - route:
      - destination:
          host: zker-backend
          subset: v1  # 旧版本
        weight: 90
      - destination:
          host: zker-backend
          subset: v2
        weight: 10
  ```

- [ ] **6.1.3** DestinationRule配置
  ```yaml
  apiVersion: networking.istio.io/v1beta1
  kind: DestinationRule
  metadata:
    name: zker-backend
  spec:
    host: zker-backend
    subsets:
    - name: v1
      labels:
        version: v1.0.0
    - name: v2
      labels:
        version: v2.0.0
  ```

##### 6.2 灰度发布流程（1人天）
- [ ] **6.2.1** 定义灰度阶段
  - 阶段1: 内部测试（5%流量）
  - 阶段2: 部分用户（20%流量）
  - 阶段3: 全量发布（100%流量）

- [ ] **6.2.2** 监控指标
  - 错误率
  - 响应时间
  - 资源使用

- [ ] **6.2.3** 回滚机制
  - 自动回滚触发条件
  - 手动回滚流程

##### 6.3 灰度发布脚本（0.5人天）
- [ ] **6.3.1** 创建灰度发布脚本
  ```bash
  #!/bin/bash
  # canary-release.sh <version> <traffic_percentage>

  VERSION=$1
  TRAFFIC=$2

  # 更新VirtualService
  kubectl apply -f canary-vs-${VERSION}.yaml

  # 等待观察
  echo "Canary release started: ${TRAFFIC}% traffic to v${VERSION}"
  echo "Monitor metrics and adjust traffic as needed"
  ```

- [ ] **6.3.2** 创建回滚脚本
  ```bash
  #!/bin/bash
  # rollback.sh <version>

  VERSION=$1

  # 恢复到旧版本
  kubectl apply -f stable-vs.yaml

  echo "Rolled back to stable version"
  ```

#### 验收标准
- ✅ Istio正常工作
- ✅ 灰度发布流程完整
- ✅ 支持流量调整
- ✅ 支持一键回滚

#### 依赖关系
- 依赖于：P2-4（K8s配置）

#### 风险评估
- 🟡 **中风险**: Istio配置复杂
  - 缓解措施：先在测试环境验证

---

## 📅 **总体时间表**

```
Week 1 (Day 1-5): P0任务完成 + P1任务启动
├─ Day 1-2: P0-1 tenant_id迁移
├─ Day 2-3: P0-2 配额中间件
├─ Day 3:   P0-3 权限中间件
├─ Day 4-5: P1-1 错误码系统 + P1-2 系统角色

Week 2 (Day 6-10): P1任务继续
├─ Day 6-8: P1-3 API客户端 + P1-4 健康检查
├─ Day 9-10: P1-5 路由引擎（开始）

Week 3 (Day 11-15): P1任务完成
├─ Day 11-12: P1-5 路由引擎（继续）
├─ Day 12-15: P1-6 配置中心 + P1-7 租户管理API

Week 4 (Day 16-21): P1任务完成 + P2任务启动
├─ Day 16-21: P1-8 性能测试 + P1-9 日志 + P1-10 监控

Week 5 (Day 22-27): P2任务执行
├─ Day 22-24: P2-1 租户管理前端 + P2-4 Docker/K8s
├─ Day 24-27: P2-2 权限管理前端 + P2-3 路由配置前端

Week 6 (Day 28-30): P2任务完成
├─ Day 25-27: P2-5 CI/CD
├─ Day 28-30: P2-6 灰度发布
```

---

## 📊 **资源分配矩阵**

| 任务 | 研发A | 研发B | 研发C | 研发D | 总人天 |
|------|-------|-------|-------|-------|---------|
| P0任务 | 8d | 4d | 0d | 0d | 12d |
| P1-1 错误码 | 0d | 5d | 0d | 0d | 5d |
| P1-2 系统角色 | 2d | 0d | 0d | 0d | 2d |
| P1-3 API客户端 | 0d | 0d | 6d | 0d | 6d |
| P1-4 健康检查 | 0d | 3d | 0d | 0d | 3d |
| P1-5 路由引擎 | 8d | 0d | 0d | 0d | 8d |
| P1-6 配置中心 | 0d | 4d | 0d | 0d | 4d |
| P1-7 租户管理API | 4d | 0d | 0d | 0d | 4d |
| P1-8 性能测试 | 0d | 6d | 0d | 0d | 6d |
| P1-9 日志收集 | 0d | 4d | 0d | 0d | 4d |
| P1-10 监控告警 | 0d | 4d | 0d | 0d | 4d |
| P2-1 租户前端 | 0d | 0d | 6d | 0d | 6d |
| P2-2 权限前端 | 0d | 0d | 4d | 0d | 4d |
| P2-3 路由前端 | 0d | 0d | 4d | 0d | 4d |
| P2-4 Docker/K8s | 0d | 0d | 0d | 3d | 3d |
| P2-5 CI/CD | 0d | 0d | 0d | 3d | 3d |
| P2-6 灰度发布 | 0d | 0d | 0d | 3d | 3d |
| **总计** | **22d** | **30d** | **20d** | **9d** | **76d** |

---

## ✅ **验收标准总结**

### 功能完整性
- ✅ 所有P0、P1、P2功能按计划实现
- ✅ 所有功能都有完整的测试
- ✅ 所有功能都有完整的文档

### 代码质量
- ✅ 代码符合企业级开发规范
- ✅ 单元测试覆盖率 ≥ 80%
- ✅ 集成测试覆盖率 ≥ 70%
- ✅ 无代码重复（DRY原则）
- ✅ 接口设计简洁（KISS原则）

### 性能指标
- ✅ API响应时间 P95 < 500ms
- ✅ API响应时间 P99 < 1s
- ✅ 系统吞吐量 ≥ 1000 QPS
- ✅ 数据库查询优化（无N+1）

### 安全性
- ✅ 所有API都有权限检查
- ✅ 租户数据完全隔离
- ✅ 配额强制执行
- ✅ 敏感数据加密存储

### 可运维性
- ✅ 完善的日志记录
- ✅ 完善的监控告警
- ✅ 自动化部署
- ✅ 灰度发布能力

---

## 🎯 **下一步行动**

### 立即开始（Day 4）
1. **研发B**: 开始P1-1错误码系统完善
2. **研发A**: 开始P1-2系统角色初始化
3. **研发C**: 准备P1-3 API客户端开发

### 本周目标（Day 4-10）
1. 完成P1-1错误码系统（5人天）
2. 完成P1-2系统角色（2人天）
3. 完成P1-3 API客户端（6人天）
4. 完成P1-4健康检查（3人天）

### 关键里程碑
- **Week 2 End**: 所有P1核心功能完成
- **Week 4 End**: 所有P1和P2功能完成
- **Week 6 End**: 系统上线

---

**文档维护**: 本文档应在每个任务完成后更新进度
**审批流程**: 每个P1任务完成后需要团队Review
**风险上报**: 遇到阻塞问题立即上报项目经理
