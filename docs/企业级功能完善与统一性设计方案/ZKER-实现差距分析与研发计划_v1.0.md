# ZKER 系统实现差距分析与研发计划

**文档版本**: v1.0
**创建日期**: 2025-01-01
**最后更新**: 2025-01-01
**负责人**: 技术架构委员会

---

## 📊 执行摘要

本报告基于对 Coze Studio 代码库的深度分析和所有 ZKER 设计文档的梳理，识别出现有实现与企业级需求之间的差距，并制定详细的 4 人并行研发计划。

**核心发现**:
- 🎯 **12 大功能模块**需要实现或增强
- 📈 **预计工作量**: 4 人 × 8 周 = 32 人周
- 🚀 **关键里程碑**: 8 周后完成所有企业级功能
- ✅ **优先级排序**: P0 (核心架构) → P1 (生产保障) → P2 (体验优化)

---

## 🔍 一、实现差距详细分析

### 1.1 多租户架构增强

#### 现状分析

**现有实现**:
```sql
-- 发现的租户隔离字段
space_id          -- 工作空间 ID（租户标识）
space_user        -- 空间成员表
  └── role_type   -- 角色: 1.owner 2.admin 3.member
```

**架构评估**:
- ✅ **数据隔离**: 已通过 `space_id` 实现基本租户隔离
- ⚠️ **租户元数据**: 缺少租户配置、订阅信息、资源配额
- ❌ **租户管理**: 缺少租户生命周期管理

#### 差距清单

| 功能 | 现状 | 设计要求 | 差距 | 优先级 |
|-----|------|---------|------|--------|
| **租户标识** | `space_id` | `tenant_id` UUID | ⚠️ 需统一 | P0 |
| **租户元数据** | 简单 `space` 表 | 完整 `tenants` 表 | ❌ 缺失 | P0 |
| **订阅管理** | 无 | `subscriptions` 表 | ❌ 缺失 | P0 |
| **资源配额** | 无 | `quotas` 表 + 检查逻辑 | ❌ 缺失 | P0 |
| **配额强制执行** | 无 | 应用层 + 数据库层 | ❌ 缺失 | P1 |
| **租户状态管理** | 无 | `active/suspended/deleted` | ❌ 缺失 | P1 |

#### 实现计划

**P0: 核心租户系统 (2 周)**

```sql
-- 1. 创建租户表
CREATE TABLE tenants (
  tenant_id VARCHAR(36) PRIMARY KEY,
  tenant_name VARCHAR(200) NOT NULL,
  tenant_type ENUM('individual', 'team', 'enterprise') NOT NULL,
  status ENUM('active', 'suspended', 'deleted') DEFAULT 'active',
  subscription_tier ENUM('free', 'pro', 'enterprise') DEFAULT 'free',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  INDEX idx_tenant_type (tenant_type),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 创建订阅表
CREATE TABLE subscriptions (
  subscription_id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  plan_tier ENUM('free', 'pro', 'enterprise') NOT NULL,
  status ENUM('active', 'past_due', 'canceled', 'expired') DEFAULT 'active',
  start_date DATE NOT NULL,
  end_date DATE,
  quota_bots INT DEFAULT 10,
  quota_messages_per_month INT DEFAULT 1000,
  quota_storage_gb INT DEFAULT 10,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
  INDEX idx_tenant_id (tenant_id),
  INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 创建配额使用表
CREATE TABLE quota_usage (
  tenant_id VARCHAR(36) NOT NULL,
  resource_type ENUM('bots', 'messages', 'storage') NOT NULL,
  used_count INT DEFAULT 0,
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  PRIMARY KEY (tenant_id, resource_type, period_start),
  FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
  INDEX idx_tenant_resource (tenant_id, resource_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. 数据迁移：从 space_id 迁移到 tenant_id
-- 执行数据迁移脚本（参考数据迁移方案文档）
```

**配额检查中间件 (Go)**:

```go
// backend/middleware/quota_check.go
package middleware

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

type QuotaCheckMiddleware struct {
    quotaService QuotaService
}

func (m *QuotaCheckMiddleware) QuotaCheck() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        tenantID := ctx.GetString("tenant_id")
        resourceType := getResourceTypeFromPath(ctx.Request.Path())

        // 检查配额
        allowed, err := m.quotaService.CheckQuota(c, tenantID, resourceType)
        if err != nil {
            ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
                "code": "QUOTA500",
                "message": "配额检查失败",
            })
            ctx.Abort()
            return
        }

        if !allowed {
            ctx.JSON(consts.StatusPaymentRequired, map[string]interface{}{
                "code": "QUOTA402",
                "message": "配额已用完，请升级订阅",
                "request_id": getRequestID(c),
            })
            ctx.Abort()
            return
        }

        ctx.Next(c)
    }
}
```

---

### 1.2 RBAC 权限系统增强

#### 现状分析

**现有实现**:
```go
// backend/crossdomain/permission/model/permission_check.go
type CheckAuthzData struct {
    UserID      string
    SpaceID     string
    ResourceType string
    ResourceID  string
    Action      string
}

// backend/domain/permission/service/permission.go
type Permission interface {
    CheckAuthz(ctx context.Context, req *CheckAuthzData) (*CheckAuthzResult, error)
}
```

**架构评估**:
- ✅ **基础框架**: 已有权限检查接口
- ⚠️ **权限粒度**: 仅 Space 级别，缺少资源级权限
- ❌ **动态权限**: 无动态策略引擎

#### 差距清单

| 功能 | 现状 | 设计要求 | 差距 | 优先级 |
|-----|------|---------|------|--------|
| **角色定义** | 3 种固定角色 | 可扩展角色系统 | ⚠️ 需增强 | P0 |
| **数据权限** | Space 级别 | 5 级数据权限 | ❌ 缺失 | P0 |
| **字段权限** | 无 | 3 级字段权限 | ❌ 缺失 | P1 |
| **动态策略** | 无 | 策略引擎 | ❌ 缺失 | P1 |
| **权限继承** | 无 | 角色继承 | ❌ 缺失 | P2 |

#### 实现计划

**P0: 5 级数据权限 + 可扩展角色 (3 周)**

```sql
-- 1. 角色表（增强版）
CREATE TABLE roles (
  role_id VARCHAR(36) PRIMARY KEY,
  tenant_id VARCHAR(36) NOT NULL,
  role_name VARCHAR(100) NOT NULL,
  role_type ENUM('system', 'custom') NOT NULL,
  parent_role_id VARCHAR(36),  -- 角色继承
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
  FOREIGN KEY (parent_role_id) REFERENCES roles(role_id),
  UNIQUE KEY uk_tenant_role (tenant_id, role_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 数据权限规则表（5级权限）
CREATE TABLE data_permissions (
  permission_id VARCHAR(36) PRIMARY KEY,
  role_id VARCHAR(36) NOT NULL,
  resource_type ENUM('bots', 'conversations', 'knowledge', 'workflows') NOT NULL,
  scope ENUM('ALL', 'DEPARTMENT', 'OWN', 'CUSTOM', 'NONE') NOT NULL,
  custom_filter JSON,  -- 自定义过滤条件
  FOREIGN KEY (role_id) REFERENCES roles(role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 字段权限表（3级权限）
CREATE TABLE field_permissions (
  permission_id VARCHAR(36) PRIMARY KEY,
  role_id VARCHAR(36) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  field_name VARCHAR(100) NOT NULL,
  permission_level ENUM('hidden', 'readonly', 'editable') NOT NULL,
  FOREIGN KEY (role_id) REFERENCES roles(role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4. 用户角色关联表
CREATE TABLE user_roles (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  tenant_id VARCHAR(36) NOT NULL,
  user_id VARCHAR(36) NOT NULL,
  role_id VARCHAR(36) NOT NULL,
  granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  granted_by VARCHAR(36),
  expires_at TIMESTAMP NULL,
  FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
  FOREIGN KEY (role_id) REFERENCES roles(role_id),
  UNIQUE KEY uk_tenant_user_role (tenant_id, user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5. 初始化系统角色
INSERT INTO roles (role_id, tenant_id, role_name, role_type) VALUES
('role_owner_all', 'system', 'Owner', 'system'),
('role_admin_all', 'system', 'Admin', 'system'),
('role_member_all', 'system', 'Member', 'system');
```

**权限检查服务实现**:

```go
// backend/domain/permission/service/rbac_permission.go
package service

import (
    "context"
    "fmt"
)

type RBACPermissionService struct {
    roleRepo       RoleRepository
    dataPermRepo   DataPermissionRepository
    fieldPermRepo  FieldPermissionRepository
}

// CheckDataPermission 检查数据权限（5级）
func (s *RBACPermissionService) CheckDataPermission(
    ctx context.Context,
    userID, tenantID, resourceType, resourceID, action string,
) (bool, error) {
    // 1. 获取用户所有角色
    roles, err := s.roleRepo.GetByUserAndTenant(ctx, userID, tenantID)
    if err != nil {
        return false, err
    }

    // 2. 检查每个角色的数据权限
    for _, role := range roles {
        dataPerm, err := s.dataPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType)
        if err != nil {
            continue
        }

        // 3. 根据权限范围判断
        allowed := s.evaluateDataScope(ctx, dataPerm.Scope, userID, resourceID, action)
        if allowed {
            return true, nil
        }
    }

    return false, nil
}

func (s *RBACPermissionService) evaluateDataScope(
    ctx context.Context,
    scope string,
    userID, resourceID, action string,
) bool {
    switch scope {
    case "ALL":
        return true  // 所有数据
    case "DEPARTMENT":
        // 检查资源是否属于用户部门
        return s.checkDepartmentAccess(ctx, userID, resourceID)
    case "OWN":
        // 检查资源是否由用户创建
        return s.checkOwnership(ctx, userID, resourceID)
    case "CUSTOM":
        // 执行自定义过滤逻辑
        return s.evaluateCustomFilter(ctx, userID, resourceID, action)
    case "NONE":
        return false
    default:
        return false
    }
}

// CheckFieldPermission 检查字段权限（3级）
func (s *RBACPermissionService) CheckFieldPermission(
    ctx context.Context,
    userID, tenantID, resourceType, fieldName string,
) (string, error) {
    roles, err := s.roleRepo.GetByUserAndTenant(ctx, userID, tenantID)
    if err != nil {
        return "hidden", err
    }

    maxLevel := "hidden"  // 默认隐藏
    for _, role := range roles {
        fieldPerm, err := s.fieldPermRepo.GetByRoleAndResource(ctx, role.RoleID, resourceType, fieldName)
        if err != nil {
            continue
        }

        // 权限级别: hidden < readonly < editable
        if s.compareFieldLevel(fieldPerm.PermissionLevel, maxLevel) > 0 {
            maxLevel = fieldPerm.PermissionLevel
        }
    }

    return maxLevel, nil
}

func (s *RBACPermissionService) compareFieldLevel(level1, level2 string) int {
    levels := map[string]int{
        "hidden":    0,
        "readonly":  1,
        "editable":  2,
    }
    return levels[level1] - levels[level2]
}
```

---

### 1.3 智能路由引擎

#### 现状分析

**现有实现**: ❌ **未发现智能路由引擎相关代码**

**设计要求**:
- 混合意图识别（规则 + 相似度 + 模型）
- 评分路由决策（多因素服务选择）
- 负载均衡与健康检查
- A/B 测试路由支持

#### 差距清单

| 模块 | 现状 | 设计要求 | 差距 | 优先级 |
|-----|------|---------|------|--------|
| **意图识别** | 无 | 混合意图匹配器 | ❌ 完全缺失 | P0 |
| **路由决策** | 无 | 评分路由器 | ❌ 完全缺失 | P0 |
| **服务发现** | 无 | 健康检查 + 注册中心 | ❌ 完全缺失 | P0 |
| **A/B 测试** | 无 | 实验路由 | ⚠️ 部分缺失 | P1 |
| **监控** | 无 | 路由决策日志 | ❌ 缺失 | P1 |

#### 实现计划

**P0: 智能路由引擎 (3 周)**

```go
// backend/domain/routing/service/intent_matcher.go
package service

import (
    "context"
    "math"
)

type IntentMatchInput struct {
    UserInput  string
    Context    map[string]interface{}
    Candidates []string  // 候选意图列表
}

type IntentMatchOutput struct {
    TopIntents    []IntentMatch
    MaxConfidence float64
    Entities      map[string]string
}

type IntentMatch struct {
    IntentName string
    Confidence  float64
    Entities    map[string]string
}

// HybridIntentMatcher 混合意图匹配器
type HybridIntentMatcher struct {
    ruleMatcher      RuleBasedMatcher
    similarityMatcher SimilarityMatcher
    modelMatcher     ModelBasedMatcher  // 可选
}

func (m *HybridIntentMatcher) Match(ctx context.Context, input *IntentMatchInput) (*IntentMatchOutput, error) {
    candidates := make(map[string]float64)

    // 1. 基于规则的匹配
    ruleResults := m.ruleMatcher.Match(ctx, input.UserInput)
    for _, result := range ruleResults {
        candidates[result.Intent] = result.Confidence
    }

    // 2. 基于相似度的匹配
    simResults, _ := m.similarityMatcher.Match(ctx, input.UserInput, input.Candidates)
    for _, result := range simResults {
        currentConf := candidates[result.Intent]
        if result.Confidence > currentConf {
            candidates[result.Intent] = result.Confidence
        }
    }

    // 3. 基于模型的匹配（可选）
    if m.modelMatcher != nil {
        modelResults, _ := m.modelMatcher.Match(ctx, input)
        for _, result := range modelResults {
            currentConf := candidates[result.Intent]
            if result.Confidence > currentConf {
                candidates[result.Intent] = result.Confidence
            }
        }
    }

    // 4. 聚合结果并排序
    topIntents := make([]IntentMatch, 0, len(candidates))
    for intent, conf := range candidates {
        topIntents = append(topIntents, IntentMatch{
            IntentName: intent,
            Confidence:  conf,
        })
    }

    sort.Slice(topIntents, func(i, j int) bool {
        return topIntents[i].Confidence > topIntents[j].Confidence
    })

    if len(topIntents) > 3 {
        topIntents = topIntents[:3]
    }

    return &IntentMatchOutput{
        TopIntents:    topIntents,
        MaxConfidence: topIntents[0].Confidence,
    }, nil
}

// RuleBasedMatcher 基于规则的匹配器
type RuleBasedMatcher struct {
    rules []IntentRule
}

type IntentRule struct {
    Intent     string
    Pattern    string  // 正则表达式
    Confidence float64
}

func (m *RuleBasedMatcher) Match(ctx context.Context, userInput string) []IntentMatch {
    results := make([]IntentMatch, 0)
    for _, rule := range m.rules {
        matched, _ := regexp.MatchString(rule.Pattern, userInput)
        if matched {
            results = append(results, IntentMatch{
                IntentName: rule.Intent,
                Confidence:  rule.Confidence,
            })
        }
    }
    return results
}

// SimilarityMatcher 基于相似度的匹配器
type SimilarityMatcher struct {
    embeddings map[string][]float64  // 意图嵌入向量
}

func (m *SimilarityMatcher) Match(ctx context.Context, userInput string, candidates []string) ([]IntentMatch, error) {
    // 1. 生成用户输入嵌入
    userEmbedding, err := m.generateEmbedding(ctx, userInput)
    if err != nil {
        return nil, err
    }

    // 2. 计算余弦相似度
    results := make([]IntentMatch, 0)
    for _, intent := range candidates {
        intentEmbedding := m.embeddings[intent]
        similarity := cosineSimilarity(userEmbedding, intentEmbedding)

        if similarity > 0.7 {
            results = append(results, IntentMatch{
                IntentName: intent,
                Confidence:  similarity,
            })
        }
    }
    return results, nil
}

func cosineSimilarity(a, b []float64) float64 {
    var dotProduct, normA, normB float64
    for i := range a {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
```

```go
// backend/domain/routing/service/score_router.go
package service

type ScoreBasedRouter struct {
    serviceRegistry ServiceRegistry
    intentMatcher   IntentMatcher
}

type RoutingDecision struct {
    SelectedService string
    DecisionReason  string
    Score           float64
}

func (r *ScoreBasedRouter) Route(ctx context.Context, input *RoutingInput) (*RoutingDecision, error) {
    // 1. 意图识别
    intentResult, err := r.intentMatcher.Match(ctx, &IntentMatchInput{
        UserInput:  input.UserInput,
        Context:    input.Context,
        Candidates: []string{"chat", "search", "workflow"},
    })
    if err != nil {
        return nil, err
    }

    // 2. 获取候选服务
    primaryIntent := intentResult.TopIntents[0].IntentName
    candidates := r.serviceRegistry.GetServicesByIntent(primaryIntent)

    // 3. 服务评分
    scoredServices := make([]ScoredService, 0)
    for _, service := range candidates {
        if service.Status != "available" {
            continue
        }

        score := r.calculateServiceScore(ctx, service, intentResult, input)
        scoredServices = append(scoredServices, ScoredService{
            Service: service,
            Score:    score,
        })
    }

    // 4. 选择最优服务
    sort.Slice(scoredServices, func(i, j int) bool {
        return scoredServices[i].Score > scoredServices[j].Score
    })

    selected := scoredServices[0]

    return &RoutingDecision{
        SelectedService: selected.Service.ID,
        DecisionReason:  fmt.Sprintf("Intent=%s, Score=%.2f", primaryIntent, selected.Score),
        Score:           selected.Score,
    }, nil
}

func (r *ScoreBasedRouter) calculateServiceScore(
    ctx context.Context,
    service *Service,
    intentResult *IntentMatchOutput,
    input *RoutingInput,
) float64 {
    score := 0.0

    // 1. 意图匹配度（权重 0.3）
    score += 0.3 * intentResult.MaxConfidence

    // 2. 服务负载（权重 0.2）
    loadScore := 1.0 - (float64(service.CurrentLoad) / float64(service.MaxLoad))
    score += 0.2 * loadScore

    // 3. 历史成功率（权重 0.2）
    score += 0.2 * service.SuccessRate

    // 4. 地域亲和性（权重 0.1）
    if service.Region == input.UserRegion {
        score += 0.1
    }

    // 5. 成本考虑（权重 0.1）
    normalizedCost := (service.MaxCost - service.Cost) / service.MaxCost
    score += 0.1 * normalizedCost

    return score
}
```

---

### 1.4 统一错误码系统

#### 现状分析

**现有实现**:
```go
// backend/types/errno/common.go
const (
    // 通用错误码
    Success            = 0
    ServerError        = 500

    // 业务错误码
    InvalidParam        = 400001
    Unauthorized        = 401001
    Forbidden          = 403001
    NotFound            = 404001
)
```

**架构评估**:
- ⚠️ **格式不统一**: 部分错误码不规范
- ❌ **覆盖不全**: 仅定义了少量错误码
- ❌ **多语言**: 无多语言支持
- ❌ **模块化**: 缺少按模块分类的错误码

#### 差距清单

| 功能 | 现状 | 设计要求 | 差距 | 优先级 |
|-----|------|---------|------|--------|
| **错误码格式** | 部分不规范 | `{模块}{类型}{编号}` | ⚠️ 需重构 | P0 |
| **错误码数量** | ~20 个 | 300+ 个 | ❌ 缺失 | P0 |
| **多语言支持** | 无 | 中文/英文/日文 | ❌ 缺失 | P1 |
| **错误码文档** | 无 | 完整文档 | ❌ 缺失 | P1 |

#### 实现计划

**P0: 统一错误码重构 (1 周)**

```go
// backend/types/errno/errors.go
package errno

import (
    "fmt"
)

// ErrorCode 错误码接口
type ErrorCode interface {
    Code() string
    Message() string
    MessageZH() string
    MessageEN() string
    MessageJA() string
    HTTPStatus() int
}

// BaseError 基础错误
type BaseError struct {
    code        string
    message     string
    messageZH   string
    messageEN   string
    messageJA   string
    httpStatus  int
}

func (e *BaseError) Code() string { return e.code }
func (e *BaseError) Message() string { return e.message }
func (e *BaseError) MessageZH() string { return e.messageZH }
func (e *BaseError) MessageEN() string { return e.messageEN }
func (e *BaseError) MessageJA() string { return e.messageJA }
func (e *BaseError) HTTPStatus() int { return e.httpStatus }

// NewError 创建错误
func NewError(code, message, messageZH, messageEN, messageJA string, httpStatus int) ErrorCode {
    return &BaseError{
        code:       code,
        message:    message,
        messageZH:  messageZH,
        messageEN:  messageEN,
        messageJA:  messageJA,
        httpStatus: httpStatus,
    }
}

// ============ 认证模块 (AUTH) ============
var (
    // 用户名或密码错误
    ErrInvalidCredentials = NewError(
        "AUTH201",
        "Invalid username or password",
        "用户名或密码错误",
        "Invalid username or password",
        "ユーザー名またはパスワードが間違っています",
        401,
    )

    // Token 过期
    ErrTokenExpired = NewError(
        "AUTH202",
        "Token expired",
        "Token 已过期",
        "Token expired",
        "トークンの有効期限が切れています",
        401,
    )

    // 权限不足
    ErrAccessDenied = NewError(
        "AUTH403",
        "Permission denied",
        "权限不足",
        "Permission denied",
        "権限が不足しています",
        403,
    )
)

// ============ Bot 模块 (BOT) ============
var (
    // Bot 不存在
    ErrBotNotFound = NewError(
        "BOT404",
        "Bot not found",
        "Bot 不存在",
        "Bot not found",
        "Botが存在しません",
        404,
    )

    // 配额超限
    ErrBotQuotaExceeded = NewError(
        "BOT422",
        "Bot quota exceeded",
        "Bot 数量达到配额上限",
        "Bot quota exceeded",
        "Botの数量がクオータを超過しました",
        422,
    )
)

// ============ 对话模块 (CONV) ============
var (
    // 对话不存在
    ErrConversationNotFound = NewError(
        "CONV404",
        "Conversation not found",
        "对话不存在",
        "Conversation not found",
        "会話が存在しません",
        404,
    )
)

// 更多模块错误码定义...
// 参考 ZKER-统一错误码定义规范.md
```

**错误响应中间件**:

```go
// backend/middleware/error_handler.go
package middleware

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

type ErrorResponse struct {
    Code       string `json:"code"`
    Message    string `json:"message"`
    MessageZH  string `json:"message_zh,omitempty"`
    MessageEN  string `json:"message_en,omitempty"`
    MessageJA  string `json:"message_ja,omitempty"`
    RequestID  string `json:"request_id"`
    Timestamp  string `json:"timestamp"`
}

func ErrorHandler() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        c.Next(ctx)

        // 检查是否有错误
        if len(c.Errors) == 0 {
            return
        }

        // 获取第一个错误
        err := c.Errors.Last()

        var response ErrorResponse

        // 类型断言，判断是否是自定义错误
        if errorCode, ok := err.Err.(errno.ErrorCode); ok {
            response = ErrorResponse{
                Code:      errorCode.Code(),
                Message:   errorCode.Message(),
                MessageZH: errorCode.MessageZH(),
                MessageEN: errorCode.MessageEN(),
                MessageJA: errorCode.MessageJA(),
                RequestID: getRequestID(ctx),
                Timestamp: time.Now().Format(time.RFC3339),
            }
            c.JSON(errorCode.HTTPStatus(), response)
        } else {
            // 未知错误
            response = ErrorResponse{
                Code:      "INTERNAL500",
                Message:   "Internal server error",
                MessageZH: "服务器内部错误",
                MessageEN: "Internal server error",
                RequestID: getRequestID(ctx),
                Timestamp: time.Now().Format(time.RFC3339),
            }
            c.JSON(500, response)
        }
    }
}
```

---

### 1.5 性能测试体系

#### 现状分析

**现有实现**: ⚠️ **未发现系统的性能测试框架**

**设计要求**:
- JMeter/K6 性能测试脚本
- 自动化性能基准测试
- 性能监控和告警

#### 差距清单

| 功能 | 现状 | 设计要求 | 差距 | 优先级 |
|-----|------|---------|------|--------|
| **性能测试脚本** | 无 | JMeter/K6 脚本库 | ❌ 缺失 | P0 |
| **性能基准** | 无 | 性能基线数据 | ❌ 缺失 | P1 |
| **性能监控** | Prometheus 部分监控 | 完整性能 Dashboard | ⚠️ 不完整 | P1 |
| **压力测试** | 无 | 定期压力测试 | ❌ 缺失 | P2 |

#### 实现计划

**P0: 性能测试框架搭建 (1 周)**

```
tests/
├── performance/
│   ├── jmeter/
│   │   ├── bot_api_test.jmx        # Bot API 性能测试
│   │   ├── conversation_test.jmx   # 对话性能测试
│   │   └── workflow_test.jmx       # 工作流性能测试
│   ├── k6/
│   │   ├── load_test.js            # 负载测试
│   │   ├── stress_test.js          # 压力测试
│   │   └── spike_test.js           # 峰值测试
│   └── scripts/
│       ├── run_performance_test.sh # 执行脚本
│       └── analyze_results.py      # 结果分析
```

**K6 测试脚本示例**:

```javascript
// tests/performance/k6/load_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = __ENV.API_URL || 'http://localhost:8001';

export const options = {
  stages: [
    { duration: '2m', target: 100 },   // 爬坡到 100 用户
    { duration: '5m', target: 100 },   // 维持 100 用户
    { duration: '2m', target: 500 },   // 爬坡到 500 用户
    { duration: '5m', target: 500 },   // 维持 500 用户
    { duration: '2m', target: 0 },     // 爬坡到 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],  // P95 < 2s
    http_req_failed: ['rate<0.01'],     // 错误率 < 1%
  },
};

export default function () {
  // 1. 登录
  let loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, JSON.stringify({
    username: 'test_user',
    password: 'test_password',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginRes, {
    'login successful': (r) => r.status === 200,
  });

  const token = loginRes.json('data.token');

  // 2. 获取 Bot 列表
  let listBots = http.get(`${BASE_URL}/api/v1/bots?page=1&page_size=20`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });

  check(listBots, {
    'bot list status 200': (r) => r.status === 200,
    'has bots': (r) => r.json('data.bots.length') > 0,
  });

  sleep(1);
}
```

---

### 1.6 其他关键差距

#### 监控与可观测性

**现状**: ⚠️ 部分监控（Prometheus + Grafana）

**缺失**:
- ❌ 分布式追踪（Jaeger/Zipkin）
- ❌ 统一日志收集（ELK Stack）
- ❌ 业务指标大盘

**实现计划**: 参考已生成的《性能测试计划》和《故障排查手册》

---

## 📋 二、4 人并行研发计划

### 2.1 团队组建

| 角色 | 姓名 | 职责 | 技能要求 |
|-----|------|------|---------|
| **研发 A** | 后端架构师 | 多租户 + RBAC + 路由引擎 | Go、DDD、架构设计 |
| **研发 B** | 后端工程师 | 错误码 + 性能 + 监控 | Go、性能优化、Prometheus |
| **研发 C** | 前端工程师 | 前端组件 + UI/UX | React、TypeScript、Semi Design |
| **研发 D** | DevOps 工程师 | CI/CD + 测试 + 部署 | Docker、K8s、JMeter |

### 2.2 并发开发策略

**原则**:
1. **模块隔离**: 按功能模块划分，减少代码冲突
2. **接口先行**: 先定义接口，再并行实现
3. **增量交付**: 分阶段交付，快速验证
4. **持续集成**: 每日合并，及早发现问题

### 2.3 8 周开发计划

#### 第 1-2 周：基础架构（P0）

**研发 A: 多租户核心 + RBAC 基础**
- [ ] Week 1:
  - 设计租户表结构
  - 实现租户管理 API
  - 数据迁移脚本
- [ ] Week 2:
  - 实现角色系统
  - 实现数据权限检查
  - 单元测试

**研发 B: 统一错误码 + API 规范**
- [ ] Week 1:
  - 重构错误码系统
  - 实现 300+ 错误码
  - 错误处理中间件
- [ ] Week 2:
  - OpenAPI 规范文档
  - API 接口定义
  - Postman 集合

**研发 C: 前端组件库 + 开发规范**
- [ ] Week 1:
  - 组件目录结构
  - 基础组件封装（Button、Input、Select）
  - 组件文档生成
- [ ] Week 2:
  - 业务组件开发（ChatInput、MessageBox）
  - Storybook 搭建
  - 组件使用指南

**研发 D: CI/CD + 测试环境**
- [ ] Week 1:
  - Docker 环境搭建
  - K8s 集群配置
  - CI/CD 流水线
- [ ] Week 2:
  - 性能测试框架
  - 监控 Dashboard
  - 自动化测试集成

#### 第 3-4 周：核心功能（P0）

**研发 A: 智能路由引擎**
- [ ] Week 3:
  - 意图识别器
  - 路由决策器
  - 服务发现
- [ ] Week 4:
  - 路由引擎集成
  - A/B 测试路由
  - 性能优化

**研发 B: 性能测试 + 监控**
- [ ] Week 3:
  - JMeter/K6 脚本开发
  - 性能基准测试
  - 性能 Dashboard
- [ ] Week 4:
  - 压力测试
  - 性能调优
  - 告警规则配置

**研发 C: 前端业务页面**
- [ ] Week 3:
  - 权限管理页面
  - 租户管理页面
  - 路由配置页面
- [ ] Week 4:
  - 监控 Dashboard 页面
  - 性能测试报告页面
  - 组件集成测试

**研发 D: 数据迁移 + 灰度发布**
- [ ] Week 3:
  - 数据迁移脚本
  - 双写迁移实现
  - 数据校验工具
- [ ] Week 4:
  - 灰度发布配置
  - 自动回滚系统
  - 部署文档

#### 第 5-6 周：生产保障（P1）

**研发 A: 高级 RBAC + 配额管理**
- [ ] Week 5:
  - 字段级权限控制
  - 动态权限策略
  - 权限继承
- [ ] Week 6:
  - 配额检查中间件
  - 配额使用统计
  - 订阅管理

**研发 B: 分布式追踪 + 日志**
- [ ] Week 5:
  - Jaeger 集成
  - 分布式追踪实现
  - Trace ID 传递
- [ ] Week 6:
  - ELK Stack 集成
  - 日志收集规范
  - 日志分析 Dashboard

**研发 C: 用户体验优化**
- [ ] Week 5:
  - 页面性能优化
  - 加载速度优化
  - 错误提示优化
- [ ] Week 6:
  - 多语言支持
  - 主题定制
  - 无障碍优化

**研发 D: 自动化测试 + 安全**
- [ ] Week 5:
  - 单元测试框架
  - 集成测试框架
  - E2E 测试框架
- [ ] Week 6:
  - 安全扫描集成
  - 渗透测试
  - 安全加固

#### 第 7-8 周：完善与上线（P2）

**全员联调与测试**:
- [ ] Week 7:
  - 端到端测试
  - 性能压测
  - 安全审计
- [ ] Week 8:
  - 灰度发布
  - 监控告警
  - 文档完善

---

### 2.4 任务分配矩阵（避免冲突）

#### 按模块划分

| 模块 | 负责人 | 协作方 | 依赖关系 |
|-----|-------|-------|---------|
| **多租户系统** | 研发 A | 研发 B、D | 独立模块，无冲突 |
| **RBAC 系统** | 研发 A | 研发 C | 依赖租户系统 |
| **智能路由** | 研发 A | 研发 B、D | 独立服务 |
| **错误码系统** | 研发 B | 全员 | 被所有模块依赖 |
| **性能测试** | 研发 B | 研发 A、C | 独立模块 |
| **前端组件** | 研发 C | 研发 A | 依赖后端 API |
| **前端页面** | 研发 C | 研发 A | 依赖组件库 |
| **CI/CD** | 研发 D | 全员 | 基础设施 |
| **数据迁移** | 研发 D | 研发 A | 独立脚本 |
| **监控告警** | 研发 D | 研发 B | 基础设施 |

#### 按代码库划分

**后端代码隔离**:
```
backend/
├── domain/
│   ├── tenant/          # 研发 A 负责
│   ├── permission/      # 研发 A 负责
│   └── routing/        # 研发 A 负责
├── types/
│   └── errno/          # 研发 B 负责（被所有模块依赖，优先开发）
├── middleware/          # 研发 A、B 负责
│   ├── quota_check.go
│   └── error_handler.go
└── tests/
    └── performance/    # 研发 B、D 负责
```

**前端代码隔离**:
```
frontend/packages/
├── arch/
│   └── bot-components/  # 研发 C 负责
├── studio/
│   ├── pages/tenant/    # 研发 C 负责
│   ├── pages/permission/ # 研发 C 负责
│   └── pages/routing/   # 研发 C 负责
└── common/
    └── ui-components/   # 研发 C 负责
```

**基础设施隔离**:
```
docker/
├── docker-compose.yml   # 研发 D 负责
├── k8s/                 # 研发 D 负责
└── scripts/             # 研发 D 负责
```

---

### 2.5 协作机制

**每日站会** (15 分钟):
- 时间: 每天上午 10:00
- 内容: 昨天完成、今天计划、遇到的阻碍

**每周技术评审** (1 小时):
- 时间: 每周五下午 3:00
- 内容: 代码审查、架构评审、风险识别

**双周冲刺回顾** (30 分钟):
- 时间: 每两周周五下午 4:00
- 内容: 进度回顾、指标评估、计划调整

---

## 📝 三、开发规范与注意事项

### 3.1 全局一致性保证

#### 代码风格统一

**Go 代码规范**:
```bash
# 使用 golangci-lint 进行代码检查
cd backend
golangci-lint run --enable-all

# 格式化代码
go fmt ./...

# 导入排序
goimports -w .
```

**TypeScript 代码规范**:
```bash
# 使用 ESLint 检查
cd frontend
rush lint

# 格式化代码
rush prettier

# 类型检查
rush typecheck
```

#### 命名规范

**数据库命名**:
```sql
-- 表名：小写 + 下划线
single_agent_draft
user_roles

-- 字段名：小写 + 下划线
tenant_id
created_at

-- 索引名：前缀 + 字段名
idx_tenant_id
uk_user_email  -- 唯一索引以 uk 开头
```

**Go 变量命名**:
```go
// 常量：大驼峰
const MaxRetries = 3

// 变量：大驼峰（导出）或 小驼峰（私有）
var UserService userService

// 函数：大驼峰
func GetUserByID(ctx context.Context, userID string) (*User, error)

// 接口：大驼峰 + er 后缀
type UserRepository interface{}
```

**TypeScript 命名**:
```typescript
// 组件：大驼峰
export const Button: React.FC<ButtonProps> = () => {}

// 变量：小驼峰
const userName: string = 'test';

// 类型/接口：大驼峰
interface UserProfile {
  userId: string;
}

// 常量：全大写 + 下划线
const MAX_RETRY_COUNT = 3;
```

#### 接口设计规范

**RESTful API**:
```yaml
# 资源命名：名词复数
GET    /api/v1/bots              # 列表
GET    /api/v1/bots/:bot_id       # 详情
POST   /api/v1/bots               # 创建
PUT    /api/v1/bots/:bot_id       # 更新
DELETE /api/v1/bots/:bot_id       # 删除

# 查询参数：蛇形命名
GET /api/v1/bots?page=1&page_size=20&status=published

# 响应格式：统一结构
{
  "code": "SUCCESS",
  "message": "操作成功",
  "data": {...},
  "request_id": "req_1234567890",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

### 3.2 并行开发注意事项

#### Git 分支策略

```
main (生产)
  ↑
develop (开发主分支)
  ↑
feature/tenant-system    # 研发 A 的功能分支
feature/rbac-system      # 研发 A 的功能分支
feature/routing-engine   # 研发 A 的功能分支
feature/error-codes      # 研发 B 的功能分支
feature/performance      # 研发 B 的功能分支
feature/ui-components    # 研发 C 的功能分支
feature/cicd            # 研发 D 的功能分支
```

**合并策略**:
1. **功能分支 → develop**: 通过 Pull Request 合并
2. **develop → main**: 每两周合并一次，打 tag
3. **代码审查**: 所有合并必须经过 Code Review

#### 接口契约管理

**IDL 定义优先**:
```thrift
// backend/idl/permission.thrift
// 定义权限接口，前后端共同遵守

struct CheckPermissionRequest {
    1: string user_id
    2: string tenant_id
    3: string resource_type
    4: string resource_id
    5: string action
}

struct CheckPermissionResponse {
    1: bool allowed
    2: string reason
}

service PermissionService {
    CheckPermissionResponse checkPermission(1: CheckPermissionRequest req)
}
```

**接口版本化**:
```
/api/v1/tenants           # 租户管理 API v1
/api/v2/tenants           # 租户管理 API v2（未来扩展）
```

#### 数据库迁移管理

**版本化迁移脚本**:
```sql
-- migrations/v1.0.0_add_tenant_tables.sql
-- migrations/v1.0.1_add_rbac_tables.sql
-- migrations/v1.1.0_add_routing_tables.sql
```

**迁移执行顺序**:
1. 研发 D 在测试环境执行迁移
2. 通知团队成员迁移已完成
3. 其他研发人员更新本地代码
4. 研发 A、B 更新业务逻辑以使用新表

### 3.3 质量保证机制

#### 单元测试要求

**后端测试覆盖率**:
```
domain/       层: 80% 覆盖率（核心业务逻辑）
application/ 层: 60% 覆盖率（应用服务）
api/          层: 40% 覆盖率（API 端点）
```

**测试示例**:
```go
// backend/domain/permission/service/rbac_permission_test.go
package service

import "testing"

func TestRBACPermissionService_CheckDataPermission(t *testing.T) {
    // Arrange
    service := setupTestService()
    ctx := context.Background()

    // Act
    allowed, err := service.CheckDataPermission(ctx, "user1", "tenant1", "bots", "bot1", "read")

    // Assert
    assert.NoError(t, err)
    assert.True(t, allowed)
}
```

**前端测试覆盖率**:
```
arch/        层: 80% 覆盖率（基础设施）
common/      层: 60% 覆盖率（通用组件）
agent-ide/   层: 40% 覆盖率（业务组件）
```

#### Code Review 清单

**每次 PR 必须检查**:
- [ ] 代码风格符合规范
- [ ] 单元测试通过
- [ ] 添加了必要的注释
- [ ] 更新了相关文档
- [ ] 没有引入安全漏洞
- [ ] 性能无明显退化
- [ ] 错误处理完善

#### 持续集成检查

**CI Pipeline 必须通过**:
```yaml
# .github/workflows/ci.yml
steps:
  - name: Lint
    run: rush lint

  - name: Type Check
    run: rush typecheck

  - name: Unit Tests
    run: rush test --coverage

  - name: Build
    run: rush build

  - name: Security Scan
    run: trivy scan . --security-checks vuln,config
```

---

## 📊 四、全局一致性检查清单

### 4.1 架构一致性

- [ ] **DDD 分层**: 所有新代码遵循 Controller → Service → Repository → DAO
- [ ] **模块边界**: 跨模块通信通过定义好的接口（Adapter 模式）
- [ ] **依赖方向**: 上层依赖下层，禁止反向依赖
- [ ] **事件驱动**: 跨模块异步通信使用 EventBus

### 4.2 数据一致性

- [ ] **租户隔离**: 所有业务表包含 `tenant_id` 字段
- [ ] **软删除**: 使用 `deleted_at` 而非 `DELETE`
- [ ] **审计字段**: `created_at`, `updated_at`, `created_by`, `updated_by`
- [ ] **索引策略**: 所有外键建索引，查询条件建索引
- [ ] **事务管理**: 跨表操作使用事务

### 4.3 API 一致性

- [ ] **路径规范**: `/api/v{version}/{resource}`
- [ ] **响应格式**: 统一的 `{code, message, data, request_id, timestamp}`
- [ ] **错误码**: 使用统一的错误码定义
- [ ] **分页参数**: `page`, `page_size` (默认 20)
- [ ] **认证**: Bearer Token 方式

### 4.4 前端一致性

- [ ] **组件复用**: 使用 `arch/` 和 `common/` 组件
- [ ] **状态管理**: 统一使用 Zustand store
- [ ] **类型安全**: 所有 Props 有 TypeScript 类型定义
- [ ] **样式方案**: 优先使用 Tailwind CSS 类名
- [ ] **错误处理**: 统一的错误提示组件

### 4.5 运维一致性

- [ ] **容器化**: 所有服务使用 Docker 镜像
- [ ] **配置管理**: 敏感信息使用环境变量或 ConfigMap
- [ ] **日志格式**: JSON 结构化日志
- [ ] **监控指标**: Prometheus 暴露 `/metrics` 端点
- [ ] **健康检查**: `/health` 端点

### 4.6 文档一致性

- [ ] **API 文档**: 所有 API 有 OpenAPI 规范
- [ ] **组件文档**: 所有组件有使用示例
- [ ] **部署文档**: 所有环境有部署指南
- [ ] **变更日志**: 每次发布更新 CHANGELOG
- [ ] **架构文档**: 关键决策有 ADR (Architecture Decision Record)

---

## 🎯 五、关键里程碑与验收标准

### Milestone 1: 基础架构完成 (Week 2)

**验收标准**:
- [x] 租户表结构设计完成
- [x] 角色权限表结构设计完成
- [x] 错误码系统重构完成
- [x] CI/CD 流水线搭建完成
- [x] 开发规范文档发布

**交付物**:
1. 数据库 DDL 脚本
2. 错误码定义文件
3. CI/CD 配置文件
4. 开发规范文档

### Milestone 2: 核心功能完成 (Week 4)

**验收标准**:
- [x] 多租户系统上线（支持 3 种租户类型）
- [x] RBAC 系统上线（5 级数据权限 + 3 级字段权限）
- [x] 智能路由引擎上线（意图识别 + 评分路由）
- [x] 前端组件库上线（30+ 组件）
- [x] 性能测试框架完成

**交付物**:
1. 租户管理 API + 前端页面
2. 权限管理 API + 前端页面
3. 路由引擎服务
4. 组件库文档站点
5. 性能测试报告

### Milestone 3: 生产保障完成 (Week 6)

**验收标准**:
- [x] 配额管理系统上线
- [x] 分布式追踪集成（Jaeger）
- [x] 日志收集系统（ELK）
- [x] 自动化测试框架（覆盖率 > 60%）
- [x] 安全扫描集成

**交付物**:
1. 配额检查中间件
2. 监控 Dashboard (Grafana)
3. 测试报告
4. 安全审计报告
5. 部署手册

### Milestone 4: 系统上线 (Week 8)

**验收标准**:
- [x] 所有 P0 功能完成
- [x] 性能测试通过（QPS 10000, P95 < 2s）
- [x] 灰度发布完成（4 个阶段）
- [x] 自动回滚系统可用
- [x] 完整文档交付

**交付物**:
1. 生产环境部署
2. 运维手册
3. 故障排查手册
4. 用户培训材料
5. 项目总结报告

---

## 📚 六、参考文档索引

所有已生成的 ZKER 设计文档：

1. [ZKER-统一错误码定义规范.md](./ZKER-统一错误码定义规范.md) - 300+ 错误码完整定义
2. [ZKER-技术组件清单与使用指南(完整版).md](./ZKER-技术组件清单与使用指南(完整版).md) - 前后端组件详细说明
3. [openapi/zker-api-v1-core-modules.yaml](./openapi/zker-api-v1-core-modules.yaml) - OpenAPI 3.0 规范
4. [ZKER-核心算法实现指南.md](./ZKER-核心算法实现指南.md) - 8 大核心算法伪代码
5. [ZKER-开发快速入门指南.md](./ZKER-开发快速入门指南.md) - 新人上手指南
6. [ZKER-性能测试计划_v1.0.md](./ZKER-性能测试计划_v1.0.md) - 完整性能测试方案
7. [ZKER-数据迁移方案_v1.0.md](./ZKER-数据迁移方案_v1.0.md) - 双写迁移方案
8. [ZKER-灰度发布策略_v1.0.md](./ZKER-灰度发布策略_v1.0.md) - 4 阶段灰度流程
9. [ZKER-一键回滚方案_v1.0.md](./ZKER-一键回滚方案_v1.0.md) - 自动化回滚系统
10. [ZKER-故障排查手册_v1.0.md](./ZKER-故障排查手册_v1.0.md) - 故障诊断流程

---

## 📞 七、联系与支持

**技术支持**:
- 架构问题: [技术架构委员会]
- 开发规范: [开发团队]
- 运维问题: [DevOps 团队]

**文档维护**:
- 版本: v1.0
- 最后更新: 2025-01-01
- 下次审查: 每两周

---

**© 2025 ZKER Project. All rights reserved.**
