# 研发B - 后端工程师开发计划

**负责人**: 研发B（后端工程师）
**开发周期**: 5周（4人并行）
**总工作量**: 29人天
**核心职责**: 数据迁移、配额检查中间件、错误码系统、性能测试
**最后更新**: 2025-12-30

---

## 📋 工作量总览

| 周次 | 主要任务 | 工作量 | 优先级 |
|------|---------|--------|--------|
| **Week 1-2** | tenant_id数据迁移 | 8人天 | **P0** ⚠️ |
| **Week 3** | 配额检查中间件 + 错误码系统 | 6人天 | **P0/P1** |
| **Week 4** | 端到端性能测试 | 6人天 | **P1** |
| **Week 5** | 性能基线文档 + Token Metering | 9人天 | **P2** |

**📌 注意**: 你的工作量最大，**tenant_id迁移是整个项目的最大风险项**，必须优先启动！

---

## ⚠️ 核心注意事项

### 1. 代码隔离原则

**你的职责范围**:
- ✅ 数据库迁移脚本和工具
- ✅ 配额检查中间件
- ✅ 统一错误码定义
- ✅ 性能测试脚本
- ✅ 性能基线文档
- ✅ Token Metering功能

**不要触碰**:
- ❌ 权限检查中间件（研发A负责）
- ❌ 前端代码（研发C负责）
- ❌ CI/CD和监控配置（研发D负责）
- ❌ 智能路由引擎（研发A负责）

### 2. 数据迁移红线

**必须遵守**:
1. 迁移脚本必须可回滚
2. 迁移必须支持幂等性
3. 每批处理不超过1000条
4. 迁移期间不锁表（使用低优先级锁）
5. 必须有数据一致性验证
6. 必须在测试环境验证通过后才能上生产

**禁止行为**:
1. 直接在生产环境执行未经测试的迁移
2. 使用不幂等的迁移脚本
3. 一次性迁移所有数据
4. 迁移期间没有监控
5. 迁移期间没有回滚方案

---

## 🎯 Week 1-2: tenant_id数据迁移 (P0 - 最大风险项)

### 任务1.1: 迁移准备和DDL变更 (Day 1)

**文件**: `backend/domain/tenant/migration/001_add_tenant_id_tables.sql`

**⚠️ 警告**: 这是整个项目最关键的任务，必须与研发A充分讨论方案后再实施！

**步骤1: 创建迁移计划文档**

创建文件：`backend/domain/tenant/migration/MIGRATION_PLAN.md`

```markdown
# tenant_id 迁移计划

## 迁移范围
- 表数量: 10+张业务表
- 数据量预估: 需要评估
- 预计停机时间: 0（零停机迁移）

## 迁移策略
1. DDL变更（添加字段）
2. 数据迁移（填充tenant_id）
3. 双写验证（同时读写新旧字段）
4. 切换（使用tenant_id）
5. 清理（移除旧逻辑）

## 回滚方案
- 每个阶段都有回滚脚本
- 验证失败立即回滚
```

**步骤2: DDL变更SQL脚本**

```sql
-- ============================================
-- tenant_id 迁移 DDL 变更
-- 执行时机: 维护窗口
-- ============================================

-- 1. bots表添加tenant_id
ALTER TABLE bots
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER bot_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 2. bot_configs表添加tenant_id
ALTER TABLE bot_configs
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER config_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 3. conversations表添加tenant_id
ALTER TABLE conversations
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER conversation_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_created (tenant_id, created_at),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 4. messages表添加tenant_id
ALTER TABLE messages
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER message_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_conversation (tenant_id, conversation_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 5. knowledge_bases表添加tenant_id
ALTER TABLE knowledge_bases
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER kb_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 6. knowledge_chunks表添加tenant_id
ALTER TABLE knowledge_chunks
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER chunk_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 7. workflows表添加tenant_id
ALTER TABLE workflows
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER workflow_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_status (tenant_id, status),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 8. workflow_executions表添加tenant_id
ALTER TABLE workflow_executions
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER execution_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD INDEX idx_tenant_workflow (tenant_id, workflow_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 9. single_agent_draft表添加tenant_id
ALTER TABLE single_agent_draft
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER draft_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;

-- 10. published_bots表添加tenant_id
ALTER TABLE published_bots
ADD COLUMN tenant_id VARCHAR(36) NULL AFTER published_id,
ADD INDEX idx_tenant_id (tenant_id),
ADD FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id) ON DELETE CASCADE;
```

**⚠️ 注意事项**:
1. **所有字段先设为NULL**，迁移完成后再改为NOT NULL
2. **外键约束使用ON DELETE CASCADE**，级联删除
3. **索引命名规范**: idx_tenant_id, idx_tenant_created等
4. **每张表的ALTER语句独立执行**，不要在一个事务中
5. **使用LOW_PRIORITY**降低锁影响：`ALTER TABLE ... LOCK=NONE`

**📖 开发规范**:
- SQL注释清晰
- 索引命名规范
- 外键约束完整
- 分批执行

**🔗 设计文档链接**:
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md) - **必读！**
- [Session和User表tenant_id字段迁移方案_v1.0.md](../Session和User表tenant_id字段迁移方案_v1.0.md)
- [数据库设计完整交付清单.md](../数据库设计完整交付清单.md)

---

### 任务1.2: 数据迁移脚本实现 (Day 2-4)

**文件**: `backend/domain/tenant/migration/migrate_tenant_id.go`

**实施步骤**:

```go
// backend/domain/tenant/migration/migrate_tenant_id.go
package migration

import (
    "context"
    "fmt"
    "time"

    "github.com/coze-studio/backend/pkg/logger"
    "gorm.io/gorm"
)

// TenantIDMigrator tenant_id迁移器
type TenantIDMigrator struct {
    db *gorm.DB
}

// NewTenantIDMigrator 创建迁移器
func NewTenantIDMigrator(db *gorm.DB) *TenantIDMigrator {
    return &TenantIDMigrator{db: db}
}

// MigrateTable 迁移单张表的tenant_id
func (m *TenantIDMigrator) MigrateTable(ctx context.Context, tableName, userIDField string) error {
    logger.CtxInfof(ctx, "[Migration] starting migration for table %s", tableName)

    batchSize := 1000
    offset := 0
    totalMigrated := 0

    for {
        // 1. 查询一批未迁移的数据
        var results []map[string]interface{}
        err := m.db.Table(tableName).
            Select(fmt.Sprintf("%s as user_id, id", userIDField)).
            Where("tenant_id IS NULL").
            Where("%s IS NOT NULL", userIDField).
            Limit(batchSize).
            Offset(offset).
            Find(&results).Error

        if err != nil {
            return fmt.Errorf("查询数据失败: %w", err)
        }

        if len(results) == 0 {
            break // 没有更多数据需要迁移
        }

        // 2. 批量更新tenant_id
        for _, row := range results {
            userID := row["user_id"].(string)
            id := row["id"].(string)

            // 从users表查找tenant_id
            var tenantID string
            err := m.db.Table("users").
                Select("tenant_id").
                Where("user_id = ?", userID).
                Pluck("tenant_id", &tenantID).Error

            if err != nil {
                logger.CtxErrorf(ctx, "[Migration] failed to get tenant_id for user %s: %v", userID, err)
                continue // 跳过这条数据，继续下一条
            }

            if tenantID == "" {
                // 用户没有租户，分配默认租户
                tenantID = "tenant_default"
                logger.CtxWarnf(ctx, "[Migration] user %s has no tenant, using default", userID)
            }

            // 更新tenant_id
            err = m.db.Table(tableName).
                Where("id = ?", id).
                Update("tenant_id", tenantID).Error

            if err != nil {
                logger.CtxErrorf(ctx, "[Migration] failed to update %s id %s: %v", tableName, id, err)
                continue
            }

            totalMigrated++
        }

        logger.CtxInfof(ctx, "[Migration] table %s: migrated %d records (batch %d)",
            tableName, totalMigrated, offset/batchSize+1)

        offset += batchSize

        // 避免锁表，每批之间休息100ms
        time.Sleep(100 * time.Millisecond)
    }

    logger.CtxInfof(ctx, "[Migration] table %s migration completed: total %d records",
        tableName, totalMigrated)

    return nil
}

// ValidateMigration 验证迁移结果
func (m *TenantIDMigrator) ValidateMigration(ctx context.Context, tableName string) error {
    // 1. 检查是否还有NULL的tenant_id
    var nullCount int64
    err := m.db.Table(tableName).
        Where("tenant_id IS NULL").
        Count(&nullCount).Error

    if err != nil {
        return fmt.Errorf("验证失败: %w", err)
    }

    if nullCount > 0 {
        return fmt.Errorf("表 %s 还有 %d 条记录的tenant_id为NULL", tableName, nullCount)
    }

    // 2. 检查外键完整性
    var invalidCount int64
    err = m.db.Table(tableName).
        Joins("LEFT JOIN tenants ON tenants.tenant_id = " + tableName + ".tenant_id").
        Where("tenants.tenant_id IS NULL").
        Count(&invalidCount).Error

    if err != nil {
        return fmt.Errorf("外键检查失败: %w", err)
    }

    if invalidCount > 0 {
        return fmt.Errorf("表 %s 有 %d 条记录的tenant_id不存在于tenants表", tableName, invalidCount)
    }

    logger.CtxInfof(ctx, "[Migration] table %s validation passed", tableName)
    return nil
}

// RollbackMigration 回滚迁移
func (m *TenantIDMigrator) RollbackMigration(ctx context.Context, tableName string) error {
    // 删除tenant_id字段的索引
    // 删除tenant_id字段
    // 详细代码省略...
    return nil
}
```

**⚠️ 注意事项**:
1. **分批处理**：每批1000条，避免锁表
2. **批次之间休息100ms**，降低数据库负载
3. **错误不中断**：单条失败记录日志，继续下一条
4. **详细日志**：每批记录进度
5. **默认租户**：用户没有租户时分配 `tenant_default`
6. **验证完整性**：检查外键完整性

**📖 开发规范**:
- 分批处理
- 错误处理优雅
- 日志详细
- 可回滚

**🔗 设计文档链接**:
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md) - 迁移策略章节
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)

---

### 任务1.3: 双写验证实现 (Day 5-6)

**文件**: `backend/domain/tenant/migration/dual_write.go`

**实施步骤**:

```go
// backend/domain/tenant/migration/dual_write.go
package migration

import (
    "context"
    "sync"

    "github.com/coze-studio/backend/pkg/logger"
)

// DualWriteVerifier 双写验证器
type DualWriteVerifier struct {
    db *gorm.DB
}

// VerifyDualWrite 验证双写数据一致性
func (v *DualWriteVerifier) VerifyDualWrite(ctx context.Context, tableName string, duration time.Duration) error {
    logger.CtxInfof(ctx, "[DualWrite] starting dual write verification for table %s", tableName)

    // 1. 查询指定时间段内的数据
    startTime := time.Now().Add(-duration)
    var records []map[string]interface{}

    err := v.db.Table(tableName).
        Where("updated_at >= ?", startTime).
        Find(&records).Error

    if err != nil {
        return fmt.Errorf("查询数据失败: %w", err)
    }

    // 2. 对比新旧字段数据
    inconsistencyCount := 0
    for _, record := range records {
        newField := record["tenant_id"]
        oldField := record["space_id"] // 假设旧字段是space_id

        if newField != oldField {
            inconsistencyCount++
            logger.CtxErrorf(ctx, "[DualWrite] data inconsistency found in record %v: new=%v, old=%v",
                record["id"], newField, oldField)
        }
    }

    // 3. 统计一致性
    consistencyRate := float64(len(records)-inconsistencyCount) / float64(len(records)) * 100

    logger.CtxInfof(ctx, "[DualWrite] verification completed: total=%d, inconsistent=%d, consistency=%.2f%%",
        len(records), inconsistencyCount, consistencyRate)

    if consistencyRate < 99.9 {
        return fmt.Errorf("数据一致性低于99.9%%: %.2f%%", consistencyRate)
    }

    return nil
}

// StartDualWriteMode 启动双写模式
func (v *DualWriteVerifier) StartDualWriteMode(ctx context.Context, tableName string) {
    // 修改CRUD操作，同时读写新旧字段
    // 详细实现省略...
}
```

**⚠️ 注意事项**:
1. **双写持续至少1周**，观察数据一致性
2. **一致性要求 ≥ 99.9%**
3. **不一致数据要有告警**
4. **双写期间要有监控**

**📖 开发规范**:
- 监控完整
- 告警及时
- 数据可追溯

**🔗 设计文档链接**:
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md) - 双写验证章节

---

### 任务1.4: 切换和清理 (Day 7-8)

**文件**: `backend/domain/tenant/migration/cutover.go`

**实施步骤**:

```go
// backend/domain/tenant/migration/cutover.go
package migration

// CutoverToTenantID 切换到使用tenant_id
func (m *TenantIDMigrator) CutoverToTenantID(ctx context.Context, tableName string) error {
    // 1. 将tenant_id字段改为NOT NULL
    err := m.db.Exec(fmt.Sprintf(`
        ALTER TABLE %s
        MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL
    `, tableName)).Error

    if err != nil {
        return fmt.Errorf("修改tenant_id为NOT NULL失败: %w", err)
    }

    logger.CtxInfof(ctx, "[Cutover] table %s: tenant_id is now NOT NULL", tableName)

    // 2. 移除旧的隔离逻辑代码
    // 详细实现省略...

    // 3. 更新查询，使用tenant_id而非旧字段
    // 详细实现省略...

    return nil
}

// CleanupLegacyData 清理旧数据
func (m *TenantIDMigrator) CleanupLegacyData(ctx context.Context, tableName string) error {
    // 1. 删除旧字段（谨慎！）
    // ALTER TABLE bots DROP COLUMN space_id;

    // 2. 删除旧索引

    logger.CtxInfof(ctx, "[Cleanup] table %s: legacy data cleaned up", tableName)
    return nil
}
```

**⚠️ 注意事项**:
1. **切换前必须验证完成**，一致性≥99.9%
2. **切换要选择低峰期**
3. **清理旧字段要谨慎**，建议保留1周后再删除
4. **切换后监控1周**，确认无问题

**📖 开发规范**:
- 谨慎操作
- 充分监控
- 可快速回滚

**🔗 设计文档链接**:
- [ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md) - 切换章节
- [ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)

---

## 🎯 Week 3: 配额检查中间件 + 错误码系统 (P0/P1)

### 任务2.1: 配额检查中间件 (Day 1-3, 3人天)

**文件**: `backend/api/middleware/quota_check.go`

**实施步骤**:

#### Day 1: 中间件核心逻辑

```go
// backend/api/middleware/quota_check.go
package middleware

import (
    "context"
    "strings"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
    "github.com/coze-studio/backend/application/tenant"
)

// QuotaCheckMiddleware 配额检查中间件
// 根据请求路径识别资源类型并检查配额
func QuotaCheckMiddleware(quotaService *tenant.TenantAppSVC) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 1. 提取tenant_id
        tenantID := ctx.GetString("tenant_id")
        if tenantID == "" {
            c.JSON(consts.StatusUnauthorized, map[string]interface{}{
                "code": "TENANT401",
                "message": "缺少租户信息",
                "message_zh": "缺少租户信息",
                "message_en": "Missing tenant information",
                "request_id": getRequestID(ctx),
            })
            c.Abort()
            return
        }

        // 2. 根据请求路径识别资源类型
        resourceType := getResourceTypeFromPath(c.Request.Path())
        if resourceType == "" {
            // 没有配额限制的路径，直接通过
            c.Next(ctx)
            return
        }

        // 3. 检查配额
        allowed, err := quotaService.CheckQuota(ctx, &tenant.CheckQuotaRequest{
            TenantID:     tenantID,
            ResourceType: resourceType,
        })

        if err != nil {
            c.JSON(consts.StatusInternalServerError, map[string]interface{}{
                "code": "QUOTA500",
                "message": "配额检查失败",
                "message_zh": "配额检查失败",
                "message_en": "Failed to check quota",
                "request_id": getRequestID(ctx),
            })
            c.Abort()
            return
        }

        if !allowed.Data.Allowed {
            c.JSON(consts.StatusPaymentRequired, map[string]interface{}{
                "code": "QUOTA402",
                "message": "配额已用完，请升级订阅",
                "message_zh": "配额已用完，请升级订阅",
                "message_en": "Quota exceeded, please upgrade your subscription",
                "resource_type": resourceType,
                "current_usage": allowed.Data.CurrentUsage,
                "max_limit": allowed.Data.MaxLimit,
                "request_id": getRequestID(ctx),
            })
            c.Abort()
            return
        }

        // 4. 配额充足，继续请求
        c.Next(ctx)
    }
}

// 资源类型映射
func getResourceTypeFromPath(path string) string {
    switch {
    case strings.HasPrefix(path, "/api/v1/bots"):
        return "bots"
    case strings.HasPrefix(path, "/api/v1/conversations"):
        return "messages"
    case strings.HasPrefix(path, "/api/v1/knowledge"):
        return "storage"
    case strings.HasPrefix(path, "/api/v1/workflows"):
        return "workflows"
    default:
        return "" // 没有配额限制
    }
}

func getRequestID(ctx context.Context) string {
    // 从context或header获取request_id
    return "req_" + time.Now().Format("20060102150405")
}
```

**⚠️ 注意事项**:
1. **不要阻塞所有请求**，只有资源创建才检查配额
2. **配额不足返回402状态码**（Payment Required）
3. **错误消息友好**，包含current_usage和max_limit
4. **request_id用于日志追踪**

**📖 开发规范**:
- 中间件简洁
- 错误消息清晰
- 包含上下文信息

**🔗 设计文档链接**:
- [API接口文档_租户计费系统.md](../API接口文档_租户计费系统.md)
- [23-租户计费系统_TokenMetering补充.md](../23-租户计费系统_TokenMetering补充.md)

#### Day 2: 单元测试

```go
// backend/api/middleware/quota_check_test.go
package middleware_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestQuotaCheckMiddleware_Allowed(t *testing.T) {
    // 测试配额充足场景
}

func TestQuotaCheckMiddleware_Exceeded(t *testing.T) {
    // 测试配额超限场景
}

func TestQuotaCheckMiddleware_NoQuotaPath(t *testing.T) {
    // 测试不需要配额检查的路径
}

func TestGetResourceTypeFromPath(t *testing.T) {
    // 测试资源类型识别
}
```

#### Day 3: 集成到路由

**修改文件**: `backend/api/router/register.go`

```go
// 在需要配额检查的路由上添加中间件
v1.POST("/bots",
    middleware.QuotaCheckMiddleware(tenantSvc),
    handler.CreateBot,
)

v1.POST("/conversations",
    middleware.QuotaCheckMiddleware(tenantSvc),
    handler.CreateConversation,
)
```

---

### 任务2.2: 统一错误码系统 (Day 4-6, 3人天)

**文件**: `backend/types/errno/*.go` (10个文件)

**实施步骤**:

#### Day 1: 错误码框架定义

**文件**: `backend/types/errno/common.go`

```go
// backend/types/errno/common.go
package errno

import (
    "net/http"
)

// ErrorCode 错误码接口
type ErrorCode interface {
    Code() string
    Message() string
    MessageZH() string
    MessageEN() string
    HTTPStatus() int
}

// BaseErrorCode 基础错误码实现
type BaseErrorCode struct {
    code       string
    message    string
    messageZH  string
    messageEN  string
    httpStatus int
}

func (e *BaseErrorCode) Code() string { return e.code }
func (e *BaseErrorCode) Message() string { return e.message }
func (e *BaseErrorCode) MessageZH() string { return e.messageZH }
func (e *BaseErrorCode) MessageEN() string { return e.messageEN }
func (e *BaseErrorCode) HTTPStatus() int { return e.httpStatus }

// 通用错误码（20个）
var (
    Success = &BaseErrorCode{"SUCCESS", "操作成功", "操作成功", "Operation successful", http.StatusOK}

    // 客户端错误 (2xx)
    InvalidParams = &BaseErrorCode{"COMMON201001", "参数错误", "参数错误", "Invalid parameters", http.StatusBadRequest}
    MissingParam  = &BaseErrorCode{"COMMON201002", "缺少必要参数", "缺少必要参数", "Missing required parameter", http.StatusBadRequest}
    InvalidFormat = &BaseErrorCode{"COMMON201003", "格式错误", "格式错误", "Invalid format", http.StatusBadRequest}

    // 认证错误 (4xx)
    Unauthorized = &BaseErrorCode{"COMMON401001", "未认证", "未认证", "Unauthorized", http.StatusUnauthorized}
    TokenExpired = &BaseErrorCode{"COMMON401002", "Token已过期", "Token已过期", "Token expired", http.StatusUnauthorized}
    InvalidToken = &BaseErrorCode{"COMMON401003", "无效的Token", "无效的Token", "Invalid token", http.StatusUnauthorized}

    // 权限错误 (4xx)
    Forbidden = &BaseErrorCode{"COMMON403001", "权限不足", "权限不足", "Permission denied", http.StatusForbidden}

    // 资源错误 (4xx)
    NotFound = &BaseErrorCode{"COMMON404001", "资源不存在", "资源不存在", "Resource not found", http.StatusNotFound}

    // 服务器错误 (5xx)
    InternalError = &BaseErrorCode{"COMMON500001", "服务器内部错误", "服务器内部错误", "Internal server error", http.StatusInternalServerError}
    DatabaseError = &BaseErrorCode{"COMMON500002", "数据库错误", "数据库错误", "Database error", http.StatusInternalServerError}
    ServiceUnavailable = &BaseErrorCode{"COMMON503001", "服务暂时不可用", "服务暂时不可用", "Service temporarily unavailable", http.StatusServiceUnavailable}
)
```

**⚠️ 注意事项**:
1. **错误码格式**: `{模块}{类型}{编号}` (如COMMON201001)
2. **必须支持中英文双语**
3. **HTTP状态码正确对应**
4. **错误码全局唯一**

**📖 开发规范**:
- 错误码命名规范
- 多语言支持
- 统一接口

**🔗 设计文档链接**:
- [ZKER-统一错误码定义规范.md](../ZKER-统一错误码定义规范.md) - **必读！**

#### Day 2: 各模块错误码定义

创建9个文件，每个模块定义25-40个错误码：

```go
// backend/types/errno/auth.go - 认证模块（40个错误码）
package errno

var (
    AUTH201001 = &BaseErrorCode{"AUTH201001", "用户名或密码错误", "用户名或密码错误", "Invalid username or password", http.StatusBadRequest}
    AUTH201002 = &BaseErrorCode{"AUTH201002", "用户名已存在", "用户名已存在", "Username already exists", http.StatusBadRequest}
    AUTH401001 = &BaseErrorCode{"AUTH401001", "未登录", "未登录", "Not logged in", http.StatusUnauthorized}
    // ... 共40个
)

// backend/types/errno/tenant.go - 租户模块（30个错误码）
package errno

var (
    TENANT401001 = &BaseErrorCode{"TENANT401001", "租户不存在", "租户不存在", "Tenant not found", http.StatusNotFound}
    TENANT403001 = &BaseErrorCode{"TENANT403001", "租户已暂停", "租户已暂停", "Tenant is suspended", http.StatusForbidden}
    // ... 共30个
)

// backend/types/errno/quota.go - 配额模块（30个错误码）
package errno

var (
    QUOTA402001 = &BaseErrorCode{"QUOTA402001", "Bot数量配额已用完", "Bot数量配额已用完", "Bot quota exceeded", http.StatusPaymentRequired}
    QUOTA402002 = &BaseErrorCode{"QUOTA402002", "消息配额已用完", "消息配额已用完", "Message quota exceeded", http.StatusPaymentRequired}
    // ... 共30个
)

// backend/types/errno/permission.go - 权限模块（40个错误码）
// backend/types/errno/routing.go - 路由模块（25个错误码）
// backend/types/errno/bot.go - Bot模块（35个错误码）
// backend/types/errno/chat.go - 对话模块（25个错误码）
// backend/types/errno/subscription.go - 订阅模块（25个错误码）
// backend/types/errno/user.go - 用户模块（30个错误码）
```

**⚠️ 注意事项**:
1. **每个文件25-40个错误码**
2. **错误码编号连续**
3. **所有错误码支持中英文**

#### Day 3: 错误码文档生成

**文件**: `backend/types/errno/ERROR_CODES.md`

```markdown
# ZKER 错误码完整列表

## 通用错误码 (COMMON)

| 错误码 | HTTP状态 | 中文消息 | 英文消息 |
|--------|---------|---------|---------|
| SUCCESS | 200 | 操作成功 | Operation successful |
| COMMON201001 | 400 | 参数错误 | Invalid parameters |
| ...

## 认证错误码 (AUTH)

| 错误码 | HTTP状态 | 中文消息 | 英文消息 |
|--------|---------|---------|---------|
| AUTH201001 | 400 | 用户名或密码错误 | Invalid username or password |
| ...

## 租户错误码 (TENANT)

...
```

**📖 开发规范**:
- 文档自动生成
- 错误码分类
- 索引完整

---

## 🎯 Week 4: 端到端性能测试 (P1)

### 任务3.1: K6性能测试脚本 (Day 1-4, 6人天)

**文件**: `tests/performance/k6/*.js`

**实施步骤**:

#### Day 1: 负载测试脚本

**文件**: `tests/performance/k6/load_test.js`

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
        { duration: '2m', target: 1000 },  // 爬坡到 1000 用户
        { duration: '5m', target: 1000 },  // 维持 1000 用户
        { duration: '2m', target: 0 },     // 爬坡到 0
    ],
    thresholds: {
        http_req_duration: ['p(95)<2000'],  // P95 < 2s
        http_req_duration: ['p(99)<5000'],  // P99 < 5s
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

    // 2. 获取Bot列表
    let listBots = http.get(`${BASE_URL}/api/v1/bots?page=1&page_size=20`, {
        headers: { 'Authorization': `Bearer ${token}` },
    });

    check(listBots, {
        'bot list status 200': (r) => r.status === 200,
        'has bots': (r) => r.json('data.bots.length') >= 0,
    });

    // 3. 创建对话
    let createConv = http.post(`${BASE_URL}/api/v1/conversations`, JSON.stringify({
        bot_id: 'test_bot_id',
        message: 'Hello',
    }), {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    });

    check(createConv, {
        'conversation created': (r) => r.status === 200,
    });

    sleep(1);
}
```

#### Day 2: 压力测试脚本

**文件**: `tests/performance/k6/stress_test.js`

```javascript
// tests/performance/k6/stress_test.js
import http from 'k6/http';

export const options = {
    stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 200 },
        { duration: '5m', target: 200 },
        { duration: '2m', target: 500 },
        { duration: '5m', target: 500 },
        { duration: '2m', target: 1000 },  // 压力峰值
        { duration: '3m', target: 1000 },
        { duration: '2m', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<3000'],  // P95 < 3s (压力下放宽要求)
        http_req_failed: ['rate<0.05'],     // 错误率 < 5% (压力下放宽要求)
    },
};
```

#### Day 3: 峰值测试脚本

**文件**: `tests/performance/k6/spike_test.js`

```javascript
// tests/performance/k6/spike_test.js
export const options = {
    stages: [
        { duration: '1m', target: 50 },
        { duration: '2m', target: 50 },
        { duration: '10s', target: 5000 },  // 突然飙升到5000
        { duration: '1m', target: 5000 },
        { duration: '10s', target: 50 },    // 快速回落
        { duration: '2m', target: 50 },
        { duration: '10s', target: 0 },
    ],
};
```

#### Day 4: 配额测试脚本

**文件**: `tests/performance/k6/quota_test.js`

```javascript
// tests/performance/k6/quota_test.js
// 测试配额检查中间件的性能影响
import http from 'k6/http';

export const options = {
    stages: [
        { duration: '1m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '1m', target: 0 },
    ],
};

export default function () {
    // 测试配额检查对性能的影响
    let res = http.post('http://localhost:8001/api/v1/bots', JSON.stringify({
        bot_name: 'test_bot',
    }), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer test_token',
        },
    });

    // 检查配额中间件是否影响响应时间
    check(res, {
        'status is 200 or 402': (r) => r.status === 200 || r.status === 402,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });
}
```

**⚠️ 注意事项**:
1. **不要在生产环境运行**压力测试
2. **测试环境数据充足**，避免影响测试结果
3. **测试前备份数据**
4. **监控数据库负载**

**📖 开发规范**:
- 脚本可配置
- 支持命令行参数
- 结果可重复

**🔗 设计文档链接**:
- [ZKER-性能测试计划_v1.0.md](../ZKER-性能测试计划_v1.0.md)
- [ZKER-性能基线文档.md](../ZKER-性能基线文档.md)

### 任务3.2: 测试执行和结果分析 (Day 5-6)

**文件**: `tests/performance/scripts/run_performance_test.sh`

```bash
#!/bin/bash
# tests/performance/scripts/run_performance_test.sh

echo "=== 运行性能测试 ==="

# 1. K6负载测试
echo "1. 运行K6负载测试..."
k6 run tests/performance/k6/load_test.js --out json=results/k6_load.json

# 2. K6压力测试
echo "2. 运行K6压力测试..."
k6 run tests/performance/k6/stress_test.js --out json=results/k6_stress.json

# 3. K6峰值测试
echo "3. 运行K6峰值测试..."
k6 run tests/performance/k6/spike_test.js --out json=results/k6_spike.json

# 4. 生成报告
echo "4. 生成测试报告..."
python3 tests/performance/scripts/generate_report.py

echo "=== 性能测试完成 ==="
```

---

## 🎯 Week 5: 性能基线文档 + Token Metering (P2)

### 任务4.1: 性能基线文档 (Day 1-3, 3人天)

**文件**: `docs/performance-baseline.md`

```markdown
# ZKER 性能基线文档

## API性能基线

| API端点 | P50 | P95 | P99 | 目标QPS |
|---------|-----|-----|-----|---------|
| POST /api/v1/auth/login | 50ms | 100ms | 200ms | 1000 |
| GET /api/v1/bots | 30ms | 80ms | 150ms | 2000 |
| POST /api/v1/bots | 100ms | 200ms | 500ms | 500 |
| POST /api/v1/conversations | 150ms | 300ms | 600ms | 800 |

## 数据库性能基线

| 操作类型 | 目标QPS | P95延迟 | 并发连接数 |
|---------|---------|--------|-----------|
| 读操作（SELECT） | 5000 | 50ms | 100 |
| 写操作（INSERT/UPDATE） | 1000 | 100ms | 50 |
| 复杂查询（JOIN） | 500 | 200ms | 20 |

## 系统资源基线

| 资源类型 | 正常范围 | 告警阈值 | 说明 |
|---------|---------|---------|------|
| CPU使用率 | < 50% | > 70% | 8核服务器 |
| 内存使用率 | < 60% | > 80% | 16GB内存 |
| 磁盘I/O | < 60% | > 80% | SSD磁盘 |
| 网络带宽 | < 50Mbps | > 80Mbps | 100Mbps带宽 |

## 并发用户基线

| 并发用户数 | API QPS | 平均响应时间 | 错误率 |
|-----------|--------|-------------|--------|
| 100 | 200 | 150ms | < 0.1% |
| 500 | 800 | 200ms | < 0.5% |
| 1000 | 1500 | 300ms | < 1% |

## 业务性能基线

| 业务指标 | 目标值 | 测量方法 |
|---------|--------|---------|
| 路由决策延迟 | < 100ms | P95延迟 |
| 意图识别准确率 | > 90% | 测试集评估 |
| 权限检查延迟 | < 10ms | P95延迟 |
| 配额检查延迟 | < 20ms | P95延迟 |
```

**⚠️ 注意事项**:
1. **基线数据基于实际测试结果**
2. **定期更新**（每月）
3. **与监控大盘联动**

**🔗 设计文档链接**:
- [ZKER-性能基线文档.md](../ZKER-性能基线文档.md)

---

### 任务4.2: Token Metering实现 (Day 4-9, 6人天)

**文件**:
- `backend/domain/billing/service/token_metering.go`
- `backend/domain/billing/repository/token_usage_repository.go`

**实施步骤**:

#### Day 1-2: 数据表设计

**文件**: `backend/domain/billing/migration/002_token_metering.sql`

```sql
-- Token使用记录表
CREATE TABLE token_usage (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    bot_id VARCHAR(36) NOT NULL,
    model_id VARCHAR(100) NOT NULL,
    prompt_tokens INT DEFAULT 0,
    completion_tokens INT DEFAULT 0,
    total_tokens INT DEFAULT 0,
    cost_usd DECIMAL(10, 6) DEFAULT 0.000000,
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_model_date (model_id, requested_at),
    INDEX idx_tenant_date (tenant_id, requested_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Token预算表
CREATE TABLE token_budgets (
    budget_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(36) NOT NULL,
    model_id VARCHAR(100) NOT NULL,
    budget_type ENUM('daily', 'weekly', 'monthly') NOT NULL,
    max_tokens INT DEFAULT 1000000,
    alert_threshold_percent INT DEFAULT 80,
    auto_degradation_model VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_tenant_model_type (tenant_id, model_id, budget_type),
    INDEX idx_tenant_type (tenant_id, budget_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### Day 3-5: Token计量服务

**文件**: `backend/domain/billing/service/token_metering.go`

```go
// backend/domain/billing/service/token_metering.go
package service

import (
    "context"
    "fmt"
)

type TokenMeteringService struct {
    tokenUsageRepo TokenUsageRepository
    budgetRepo     TokenBudgetRepository
}

// RecordTokenUsage 记录Token使用
func (s *TokenMeteringService) RecordTokenUsage(ctx context.Context, req *RecordTokenRequest) error {
    // 1. 计算成本
    cost := s.calculateCost(req.ModelID, req.TotalTokens)

    // 2. 记录使用
    usage := &TokenUsage{
        TenantID:        req.TenantID,
        UserID:          req.UserID,
        BotID:           req.BotID,
        ModelID:         req.ModelID,
        PromptTokens:    req.PromptTokens,
        CompletionTokens: req.CompletionTokens,
        TotalTokens:     req.TotalTokens,
        CostUSD:         cost,
    }

    if err := s.tokenUsageRepo.Create(ctx, usage); err != nil {
        return fmt.Errorf("记录Token使用失败: %w", err)
    }

    // 3. 检查预算
    return s.checkBudgetAlert(ctx, req.TenantID, req.ModelID)
}

// checkBudgetAlert 检查预算告警
func (s *TokenMeteringService) checkBudgetAlert(ctx context.Context, tenantID, modelID string) error {
    // 1. 获取预算配置
    budget, err := s.budgetRepo.GetByTenantAndModel(ctx, tenantID, modelID, "monthly")
    if err != nil {
        return err
    }

    // 2. 获取本月使用量
    usage, err := s.tokenUsageRepo.GetMonthlyUsage(ctx, tenantID, modelID)
    if err != nil {
        return err
    }

    // 3. 计算使用率
    usagePercent := float64(usage) / float64(budget.MaxTokens) * 100

    // 4. 触发告警
    if usagePercent >= float64(budget.AlertThresholdPercent) {
        s.sendBudgetAlert(ctx, tenantID, modelID, usagePercent, budget.MaxTokens, usage)
    }

    return nil
}

// sendBudgetAlert 发送预算告警
func (s *TokenMeteringService) sendBudgetAlert(ctx context.Context, tenantID, modelID string, usagePercent float64, maxTokens, usedTokens int) {
    // 发送告警通知
    logs.CtxWarnf(ctx, "[TokenMetering] budget alert: tenant=%s model=%s usage=%.2f%% (%d/%d)",
        tenantID, modelID, usagePercent, usedTokens, maxTokens)

    // TODO: 发送邮件/钉钉/企业微信通知
}
```

#### Day 6-7: 自动降级策略

```go
// GetRecommendedModel 获取推荐模型（预算超限时自动降级）
func (s *TokenMeteringService) GetRecommendedModel(ctx context.Context, tenantID, currentModelID string) (string, error) {
    // 1. 检查预算是否超限
    budget, err := s.budgetRepo.GetByTenantAndModel(ctx, tenantID, currentModelID, "monthly")
    if err != nil {
        return "", err
    }

    usage, err := s.tokenUsageRepo.GetMonthlyUsage(ctx, tenantID, currentModelID)
    if err != nil {
        return "", err
    }

    if usage < budget.MaxTokens {
        // 未超限，使用当前模型
        return currentModelID, nil
    }

    // 2. 预算已超限，返回降级模型
    if budget.AutoDegradationModel != "" {
        logs.CtxInfof(ctx, "[TokenMetering] auto degradation: %s -> %s",
            currentModelID, budget.AutoDegradationModel)
        return budget.AutoDegradationModel, nil
    }

    // 3. 没有配置降级模型，返回错误
    return "", fmt.Errorf("预算已超限且未配置降级模型")
}
```

#### Day 8-9: 单元测试和集成

**⚠️ 注意事项**:
1. **Token计量要准确**，每次API调用都要记录
2. **预算告警要及时**，超过阈值立即通知
3. **自动降级要可配置**，租户可选择是否启用
4. **成本计算要准确**，使用模型官方定价

**📖 开发规范**:
- 计量准确
- 告警及时
- 配置灵活

**🔗 设计文档链接**:
- [23-租户计费系统_TokenMetering补充.md](../23-租户计费系统_TokenMetering补充.md)

---

## 📊 关键文档链接汇总

### 必读文档（开发前必读）

1. **[ZKER-数据迁移方案_v1.0.md](../ZKER-数据迁移方案_v1.0.md)** ⭐⭐⭐
   - 零停机迁移方案
   - 双写验证策略
   - 回滚方案

2. **[ZKER-统一错误码定义规范.md](../ZKER-统一错误码定义规范.md)** ⭐⭐⭐
   - 300+错误码定义
   - 错误码格式规范
   - 多语言支持

3. **[ZKER-性能测试计划_v1.0.md](../ZKER-性能测试计划_v1.0.md)** ⭐⭐
   - 性能测试策略
   - 测试脚本设计

4. **[23-租户计费系统_TokenMetering补充.md](../23-租户计费系统_TokenMetering补充.md)** ⭐⭐
   - Token计量实现
   - 预算告警
   - 自动降级策略

### 参考文档

5. **[Session和User表tenant_id字段迁移方案_v1.0.md](../Session和User表tenant_id字段迁移方案_v1.0.md)** - 迁移示例
6. **[ZKER-一键回滚方案_v1.0.md](../ZKER-一键回滚方案_v1.0.md)** - 回滚策略
7. **[数据库设计完整交付清单.md](../数据库设计完整交付清单.md)** - 表结构设计

---

## ⚠️ 核心注意事项

### 1. 代码隔离原则（重申）

**你的职责范围**:
- ✅ 数据库迁移脚本和工具
- ✅ 配额检查中间件
- ✅ 统一错误码定义
- ✅ 性能测试脚本
- ✅ 性能基线文档
- ✅ Token Metering功能

**不要触碰**:
- ❌ 权限检查中间件（研发A负责）
- ❌ 前端代码（研发C负责）
- ❌ CI/CD和监控配置（研发D负责）

### 2. 开发规范红线

**必须遵守**:
1. 所有函数不超过50行
2. 所有错误使用统一错误码
3. 所有日志包含request_id、tenant_id
4. 单元测试覆盖率 ≥ 80%
5. 迁移脚本必须可回滚
6. 性能测试不在生产环境运行

**禁止行为**:
1. 在生产环境执行未经测试的迁移
2. 使用不幂等的迁移脚本
3. 一次性迁移所有数据（必须分批）
4. 迁移期间没有监控
5. 迁移期间没有回滚方案

### 3. Git提交规范

**Commit Message格式**:
```
feat(migration): implement tenant_id migration for bots table

- Add tenant_id column to bots table
- Implement batch migration (1000 records per batch)
- Add data validation and rollback scripts
- Migration completed for 10,000+ records

Refs: #123
```

---

## 📅 每日工作检查清单

### 开发前
- [ ] 阅读相关设计文档
- [ ] 理解迁移策略和风险
- [ ] 准备回滚方案
- [ ] 在测试环境验证

### 开发中
- [ ] 分批处理数据
- [ ] 详细的日志记录
- [ ] 错误不中断，记录并继续
- [ ] 定期备份

### 提交前
- [ ] 迁移脚本可回滚
- [ ] 数据一致性验证通过
- [ ] 测试环境验证通过
- [ ] 性能测试通过

---

## 🎯 成功标准

### Week 1-2结束时（tenant_id迁移）
- [ ] 10+张表添加tenant_id字段
- [ ] 所有数据迁移完成
- [ ] 数据一致性验证通过（100%）
- [ ] 回滚测试成功

### Week 3结束时（中间件和错误码）
- [ ] 配额检查中间件实现并集成
- [ ] 300+错误码全部定义
- [ ] 所有错误码支持中英文
- [ ] 单元测试覆盖率 ≥ 80%

### Week 4结束时（性能测试）
- [ ] 4种K6测试脚本全部实现
- [ ] 性能测试报告生成
- [ ] 性能基线数据明确
- [ ] 无性能瓶颈

### Week 5结束时（Token Metering）
- [ ] Token计量功能实现
- [ ] 预算告警功能实现
- [ ] 自动降级策略实现
- [ ] 文档完整交付

---

**🎉 你的工作量最大，但也是整个项目的基础。数据迁移是最大风险，务必谨慎！**
**遇到问题随时查阅设计文档或与研发A讨论架构方案。**
