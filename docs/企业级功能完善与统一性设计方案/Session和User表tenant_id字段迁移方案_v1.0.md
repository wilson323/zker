# Session和User表tenant_id字段迁移方案

## 📖 文档信息

- **版本**: v1.0
- **创建日期**: 2025-01-01
- **作者**: 研发B（后端工程师）
- **状态**: 📋 计划阶段
- **优先级**: P0（阻止生产部署）

---

## 🎯 迁移目标

### 核心问题

当前系统存在**架构缺陷**：
1. ❌ `Session`实体没有`TenantID`字段
2. ❌ `users`表没有`tenant_id`列
3. ❌ 用户与租户没有关联关系
4. ⚠️  无法从Session中提取tenant_id

**影响**：
- 租户隔离中间件只能从HTTP Header提取tenant_id
- 无法实现真正的多租户数据隔离
- **P0安全风险**：存在数据泄露可能性

### 迁移目标

✅ **Session添加TenantID字段**
- 修改`Session`实体结构
- 更新Session生成逻辑（登录时关联租户）
- 兼容现有Session（添加迁移逻辑）

✅ **users表添加tenant_id列**
- ALTER TABLE添加列
- 创建索引优化查询
- 数据迁移：为现有用户分配默认租户

✅ **创建user_tenant关联表**
- 支持用户多租户（未来扩展）
- 记录用户在租户中的角色
- 实现用户租户切换功能

---

## 📋 迁移计划

### 阶段划分

```
阶段1: 数据库变更（1-2天）
  ├─ 1.1 Session表添加tenant_id列
  ├─ 1.2 users表添加tenant_id列
  ├─ 1.3 创建user_tenants关联表
  └─ 1.4 创建索引和约束

阶段2: 代码变更（2-3天）
  ├─ 2.1 修改Session实体结构
  ├─ 2.2 更新Session生成逻辑
  ├─ 2.3 更新租户隔离中间件
  └─ 2.4 添加兼容性处理

阶段3: 数据迁移（1天）
  ├─ 3.1 为现有users分配默认租户
  ├─ 3.2 迁移现有Session数据
  └─ 3.3 验证数据完整性

阶段4: 测试验证（2-3天）
  ├─ 4.1 单元测试
  ├─ 4.2 集成测试
  ├─ 4.3 性能测试
  └─ 4.4 灰度发布

总计: 6-9天
```

---

## 🗄️ 阶段1: 数据库变更

### 1.1 Session表添加tenant_id列

**注意**: 当前系统Session存储在Redis中，没有独立的数据库表。
这个迁移实际上是**修改Session结构体**和**更新Redis数据**。

### 1.2 users表添加tenant_id列

#### 迁移SQL

```sql
-- =====================================================
-- 迁移脚本: users表添加tenant_id列
-- 版本: v1.0
-- 日期: 2025-01-01
-- =====================================================

-- Step 1: 添加tenant_id列（允许NULL，兼容现有数据）
ALTER TABLE `opencoze`.`users`
ADD COLUMN `tenant_id` VARCHAR(36) NULL COMMENT '租户ID（UUID）'
AFTER `session_key`;

-- Step 2: 创建索引（优化查询性能）
CREATE INDEX `idx_tenant_id` ON `opencoze`.`users` (`tenant_id`);

-- Step 3: 创建复合索引（优化租户+用户查询）
CREATE INDEX `idx_tenant_id_email` ON `opencoze`.`users` (`tenant_id`, `email`);

-- Step 4: 验证列是否添加成功
SELECT
    COLUMN_NAME,
    DATA_TYPE,
    IS_NULLABLE,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'users'
  AND COLUMN_NAME = 'tenant_id';
```

#### 回滚SQL

```sql
-- =====================================================
-- 回滚脚本: 删除users表的tenant_id列
-- 版本: v1.0
-- =====================================================

-- Step 1: 删除索引
DROP INDEX `idx_tenant_id` ON `opencoze`.`users`;
DROP INDEX `idx_tenant_id_email` ON `opencoze`.`users`;

-- Step 2: 删除列
ALTER TABLE `opencoze`.`users`
DROP COLUMN `tenant_id`;

-- Step 3: 验证列是否删除成功
SELECT COUNT(*) AS column_exists
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = 'opencoze'
  AND TABLE_NAME = 'users'
  AND COLUMN_NAME = 'tenant_id';
-- 期望结果: 0
```

---

### 1.3 创建user_tenants关联表

#### 建表SQL

```sql
-- =====================================================
-- 建表脚本: user_tenants关联表
-- 版本: v1.0
-- 说明: 支持用户多租户，记录用户在租户中的角色
-- =====================================================

CREATE TABLE `opencoze`.`user_tenants` (
    -- 主键
    `user_tenant_id` VARCHAR(36) PRIMARY KEY COMMENT '用户租户关联ID（UUID）',

    -- 外键
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 用户在租户中的角色
    `role` ENUM('owner', 'admin', 'member') DEFAULT 'member' COMMENT '角色',
    `is_default` BOOLEAN DEFAULT FALSE COMMENT '是否为默认租户',

    -- 时间戳
    `created_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒）',
    `updated_at` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒）',
    `deleted_at` BIGINT UNSIGNED NULL COMMENT '删除时间（毫秒）',

    -- 外键约束
    FOREIGN KEY (`user_id`) REFERENCES `opencoze`.`users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`tenant_id`) REFERENCES `opencoze`.`tenants`(`tenant_id`) ON DELETE CASCADE,

    -- 唯一约束：一个用户在一个租户中只能有一条记录
    UNIQUE INDEX `uk_user_tenant` (`user_id`, `tenant_id`),

    -- 索引
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_is_default` (`is_default`, `deleted_at`)
) CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户租户关联表';
```

#### 回滚SQL

```sql
-- =====================================================
-- 回滚脚本: 删除user_tenants表
-- 版本: v1.0
-- =====================================================

DROP TABLE IF EXISTS `opencoze`.`user_tenants`;
```

---

## 💻 阶段2: 代码变更

### 2.1 修改Session实体结构

#### 文件: `backend/domain/user/entity/session.go`

```go
package entity

import (
	"time"
)

const SessionKey = "session_key"

// Session 用户会话
type Session struct {
	UserID    int64
	TenantID  string    // 新增：租户ID
	Locale    string
	UserEmail string

	CreatedAt time.Time
	ExpiresAt time.Time
}

// GetTenantID 获取租户ID（兼容处理）
func (s *Session) GetTenantID() string {
	if s.TenantID != "" {
		return s.TenantID
	}
	// 兼容旧Session：返回空字符串，由调用方决定如何处理
	return ""
}

// HasTenantID 检查是否有租户ID
func (s *Session) HasTenantID() bool {
	return s.TenantID != ""
}
```

### 2.2 更新Session生成逻辑

#### 文件: `backend/application/user/...`（登录逻辑）

**修改点**: 用户登录时，将`tenant_id`写入Session

```go
// Login 用户登录
func (s *UserApplication) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 验证用户凭据
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// 2. 验证密码
	if !s.passwordService.Verify(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	// 3. 获取用户的默认租户（新增逻辑）
	tenantID, err := s.getTenantIDForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// 4. 创建Session（包含tenant_id）
	session := &entity.Session{
		UserID:    user.ID,
		TenantID:  tenantID, // 新增字段
		Locale:    user.Locale,
		UserEmail: user.Email,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30天
	}

	// 5. 存储Session到Redis
	sessionKey := generateSessionKey()
	if err := s.redis.Set(ctx, sessionKey, session, 30*24*time.Hour); err != nil {
		return nil, err
	}

	return &LoginResponse{
		SessionKey: sessionKey,
		User:       user,
		TenantID:   tenantID,
	}, nil
}

// getTenantIDForUser 获取用户的默认租户ID
func (s *UserApplication) getTenantIDForUser(ctx context.Context, userID int64) (string, error) {
	// 1. 尝试从user_tenants表获取默认租户
	var userTenant struct {
		TenantID string
	}
	err := s.db.WithContext(ctx).
		Table("user_tenants").
		Select("tenant_id").
		Where("user_id = ? AND is_default = TRUE AND deleted_at IS NULL", userID).
		First(&userTenant).Error

	if err == nil {
		return userTenant.TenantID, nil
	}

	// 2. 如果没有默认租户，获取用户的第一个租户
	err = s.db.WithContext(ctx).
		Table("user_tenants").
		Select("tenant_id").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at ASC").
		First(&userTenant).Error

	if err == nil {
		return userTenant.TenantID, nil
	}

	// 3. 如果用户没有任何租户，返回默认租户（兼容模式）
	// TODO: 数据迁移后移除，改为返回错误
	return "default", nil
}
```

### 2.3 更新租户隔离中间件

#### 文件: `backend/api/middleware/tenant_isolation.go`

**修改点**: 启用从Session提取tenant_id的逻辑

```go
// extractTenantID 从请求中提取tenant_id
//
// **优先级**：
// 1. HTTP Header: X-Tenant-ID（用于API调用）
// 2. Session: session.TenantID（✅ 数据迁移后启用）
// 3. JWT Token: claims.tenant_id（TODO：JWT改造后启用）
// 4. 默认值: default（兼容模式，数据迁移后移除）
func extractTenantID(c context.Context, ctx *app.RequestContext) (string, error) {
	// 1. 从HTTP Header获取（优先级最高）
	if tenantID := string(ctx.GetHeader(TenantIDHeader)); tenantID != "" {
		return tenantID, nil
	}

	// 2. 从Session获取（✅ 数据迁移后启用）
	if session, ok := ctxcache.Get[*entity.Session](c, consts.SessionDataKeyInCtx); ok {
		if session.HasTenantID() {
			return session.TenantID, nil
		}
	}

	// 3. 从JWT Token获取（TODO：JWT改造后启用）
	// if claims := getJWTClaims(ctx); claims != nil && claims.TenantID != "" {
	// 	return claims.TenantID, nil
	// }

	// 4. 兼容模式：使用默认租户ID（数据迁移后移除）
	logs.CtxWarnf(c, "[TenantIsolation] no tenant_id found, using default: %s", DefaultTenantID)
	return DefaultTenantID, nil
}
```

---

## 🔄 阶段3: 数据迁移

### 3.1 为现有users分配默认租户

#### 迁移策略

**方案A：为每个现有用户创建个人租户（推荐）**

```sql
-- =====================================================
-- 数据迁移: 为现有users创建个人租户并关联
-- 版本: v1.0
-- 策略: 为每个用户创建一个individual类型的租户
-- =====================================================

-- Step 1: 为每个用户创建个人租户
INSERT INTO `opencoze`.`tenants` (
    `tenant_id`,
    `tenant_name`,
    `tenant_type`,
    `status`,
    `subscription_tier`,
    `created_at`,
    `updated_at`
)
SELECT
    CONCAT('individual-user-', `id`) AS `tenant_id`,
    CONCAT(`name`, '\'s Personal Space') AS `tenant_name`,
    'individual' AS `tenant_type`,
    'active' AS `status`,
    'free' AS `subscription_tier`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `created_at`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `updated_at`
FROM `opencoze`.`users`
WHERE `deleted_at` IS NULL;

-- Step 2: 更新users表的tenant_id列
UPDATE `opencoze`.`users` u
SET u.`tenant_id` = CONCAT('individual-user-', u.`id`)
WHERE u.`deleted_at` IS NULL
  AND u.`tenant_id` IS NULL;

-- Step 3: 在user_tenants表中创建关联记录
INSERT INTO `opencoze`.`user_tenants` (
    `user_tenant_id`,
    `user_id`,
    `tenant_id`,
    `role`,
    `is_default`,
    `created_at`,
    `updated_at`
)
SELECT
    UUID() AS `user_tenant_id`,
    u.`id` AS `user_id`,
    u.`tenant_id` AS `tenant_id`,
    'owner' AS `role`,
    TRUE AS `is_default`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `created_at`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `updated_at`
FROM `opencoze`.`users` u
WHERE u.`deleted_at` IS NULL
  AND u.`tenant_id` IS NOT NULL;

-- Step 4: 验证数据迁移结果
SELECT
    (SELECT COUNT(*) FROM `opencoze`.`users` WHERE `deleted_at` IS NULL) AS total_users,
    (SELECT COUNT(*) FROM `opencoze`.`users` WHERE `tenant_id` IS NOT NULL AND `deleted_at` IS NULL) AS users_with_tenant,
    (SELECT COUNT(*) FROM `opencoze`.`tenants` WHERE `tenant_type` = 'individual' AND `status` = 'active') AS individual_tenants,
    (SELECT COUNT(*) FROM `opencoze`.`user_tenants` WHERE `deleted_at` IS NULL) AS user_tenant_relations;
```

**方案B：所有现有用户共享一个默认租户（简单但不推荐）**

```sql
-- =====================================================
-- 数据迁移: 所有现有用户共享默认租户
-- 版本: v1.0
-- 策略: 创建一个名为"default-org"的租户，所有用户关联到此租户
-- =====================================================

-- Step 1: 创建默认组织租户
INSERT INTO `opencoze`.`tenants` (
    `tenant_id`,
    `tenant_name`,
    `tenant_type`,
    `status`,
    `subscription_tier`,
    `created_at`,
    `updated_at`
) VALUES (
    'default-org',
    'Default Organization',
    'team',
    'active',
    'free',
    UNIX_TIMESTAMP(NOW()) * 1000,
    UNIX_TIMESTAMP(NOW()) * 1000
);

-- Step 2: 更新所有users的tenant_id
UPDATE `opencoze`.`users`
SET `tenant_id` = 'default-org'
WHERE `tenant_id` IS NULL
  AND `deleted_at` IS NULL;

-- Step 3: 在user_tenants表中创建关联记录
INSERT INTO `opencoze`.`user_tenants` (
    `user_tenant_id`,
    `user_id`,
    `tenant_id`,
    `role`,
    `is_default`,
    `created_at`,
    `updated_at`
)
SELECT
    UUID() AS `user_tenant_id`,
    `id` AS `user_id`,
    'default-org' AS `tenant_id`,
    'member' AS `role`,
    TRUE AS `is_default`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `created_at`,
    UNIX_TIMESTAMP(NOW()) * 1000 AS `updated_at`
FROM `opencoze`.`users`
WHERE `deleted_at` IS NULL;
```

### 3.2 迁移现有Session数据

#### Redis迁移脚本

**问题**: 现有Session存储在Redis中，没有`tenant_id`字段

**解决方案**: 两种策略

**策略A：惰性迁移（推荐）**

- 不主动更新Redis中的现有Session
- 等待用户重新登录时，自然生成包含`tenant_id`的新Session
- 兼容旧Session：如果没有`tenant_id`，使用默认值

**优点**:
- 无风险，不影响在线用户
- 逐步迁移，平滑过渡

**缺点**:
- 迁移周期长（需要等到Session过期，最多30天）

**策略B：立即迁移**

- 读取Redis中所有Session
- 为每个Session添加`tenant_id`
- 写回Redis

**优点**:
- 立即生效

**缺点**:
- 高风险（可能破坏现有Session）
- 性能影响（需要遍历所有Session）

**推荐**: 使用**策略A（惰性迁移）**

### 3.3 验证数据完整性

#### 验证SQL

```sql
-- =====================================================
-- 数据验证: 检查数据完整性
-- 版本: v1.0
-- =====================================================

-- 验证1: 所有users都应该有tenant_id
SELECT
    COUNT(*) AS users_without_tenant,
    'ERROR: Users without tenant_id' AS message
FROM `opencoze`.`users`
WHERE `tenant_id` IS NULL
  AND `deleted_at` IS NULL;

-- 期望结果: 0

-- 验证2: 所有users在user_tenants表中都应该有对应记录
SELECT
    COUNT(*) AS orphan_users,
    'ERROR: Users without user_tenant relation' AS message
FROM `opencoze`.`users` u
LEFT JOIN `opencoze`.`user_tenants` ut ON u.`id` = ut.`user_id` AND ut.`deleted_at` IS NULL
WHERE u.`deleted_at` IS NULL
  AND ut.`user_tenant_id` IS NULL;

-- 期望结果: 0

-- 验证3: user_tenants表中的所有tenant_id都应该在tenants表中存在
SELECT
    COUNT(*) AS invalid_tenants,
    'ERROR: User_tenants with non-existent tenants' AS message
FROM `opencoze`.`user_tenants` ut
LEFT JOIN `opencoze`.`tenants` t ON ut.`tenant_id` = t.`tenant_id`
WHERE ut.`deleted_at` IS NULL
  AND t.`tenant_id` IS NULL;

-- 期望结果: 0

-- 验证4: 每个用户至少有一个默认租户
SELECT
    u.`id` AS user_id,
    u.`email`,
    'WARNING: User without default tenant' AS message
FROM `opencoze`.`users` u
LEFT JOIN `opencoze`.`user_tenants` ut ON u.`id` = ut.`user_id` AND ut.`is_default` = TRUE AND ut.`deleted_at` IS NULL
WHERE u.`deleted_at` IS NULL
  AND ut.`user_tenant_id` IS NULL
LIMIT 10;

-- 期望结果: 0 rows

-- 验证5: 统计数据
SELECT
    'Total Users' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`users`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Users with Tenant ID' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`users`
WHERE `tenant_id` IS NOT NULL
  AND `deleted_at` IS NULL

UNION ALL

SELECT
    'Total Tenants' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`tenants`
WHERE `deleted_at` IS NULL

UNION ALL

SELECT
    'Individual Tenants' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`tenants`
WHERE `tenant_type` = 'individual'
  AND `deleted_at` IS NULL

UNION ALL

SELECT
    'User-Tenant Relations' AS metric,
    COUNT(*) AS value
FROM `opencoze`.`user_tenants`
WHERE `deleted_at` IS NULL;
```

---

## 🧪 阶段4: 测试验证

### 4.1 单元测试

#### 测试文件: `backend/domain/user/entity/session_test.go`

```go
package entity

import (
	"testing"
	"time"
)

func TestSession_GetTenantID(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected string
	}{
		{
			name: "Session with tenant_id",
			session: &Session{
				UserID:    123,
				TenantID:  "tenant-abc",
				UserEmail: "user@example.com",
			},
			expected: "tenant-abc",
		},
		{
			name: "Session without tenant_id (legacy)",
			session: &Session{
				UserID:    123,
				TenantID:  "",
				UserEmail: "user@example.com",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.GetTenantID()
			if result != tt.expected {
				t.Errorf("GetTenantID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSession_HasTenantID(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected bool
	}{
		{
			name: "Session with tenant_id",
			session: &Session{
				TenantID: "tenant-abc",
			},
			expected: true,
		},
		{
			name:     "Session without tenant_id",
			session:  &Session{TenantID: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.HasTenantID()
			if result != tt.expected {
				t.Errorf("HasTenantID() = %v, want %v", result, tt.expected)
			}
		})
	}
}
```

### 4.2 集成测试

#### 测试文件: `backend/api/middleware/tenant_isolation_integration_test.go`

```go
package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-studio/coze-studio/backend/domain/user/entity"
	"github.com/coze-studio/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-studio/coze-studio/backend/types/consts"
)

func TestTenantIsolation_WithSession(t *testing.T) {
	// 创建包含tenant_id的Session
	session := &entity.Session{
		UserID:    123,
		TenantID:  "test-tenant-456",
		UserEmail: "user@example.com",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	// 创建context并存储session
	ctx := context.Background()
	ctxcache.Store(ctx, consts.SessionDataKeyInCtx, session)

	// 测试extractTenantID函数
	// TODO: 添加完整的集成测试
}
```

### 4.3 性能测试

#### 测试SQL查询性能

```sql
-- 性能测试: 查询用户的租户
EXPLAIN SELECT
    ut.`tenant_id`,
    t.`tenant_name`,
    ut.`role`
FROM `opencoze`.`user_tenants` ut
JOIN `opencoze`.`tenants` t ON ut.`tenant_id` = t.`tenant_id`
WHERE ut.`user_id` = 123
  AND ut.`deleted_at` IS NULL;

-- 期望: 使用索引 idx_user_id
```

#### 性能基准

- 查询用户租户: < 10ms
- 验证租户状态: < 5ms
- Session提取tenant_id: < 1ms

---

## 🚀 部署计划

### 灰度发布策略

**阶段1: 内部测试（1天）**
- [ ] 开发环境验证
- [ ] 测试环境验证
- [ ] 性能测试

**阶段2: 灰度发布（2-3天）**
- [ ] 5% 流量
- [ ] 20% 流量
- [ ] 50% 流量
- [ ] 100% 流量

**阶段3: 监控观察（3天）**
- [ ] 错误率监控
- [ ] 性能指标监控
- [ ] 用户反馈收集

### 回滚方案

**触发条件**:
- 错误率 > 1%
- API响应时间 > 500ms
- 数据完整性检查失败

**回滚步骤**:
1. 执行回滚SQL（删除tenant_id列）
2. 代码回滚到迁移前版本
3. 重启服务
4. 验证系统恢复正常

---

## 📊 风险评估

### 高风险项

1. **数据迁移风险**
   - 风险: 迁移过程中可能破坏现有数据
   - 缓解: 充分测试，使用事务，准备回滚方案

2. **性能风险**
   - 风险: 添加tenant_id查询可能影响性能
   - 缓解: 创建索引，性能测试，优化查询

3. **兼容性风险**
   - 风险: 旧Session没有tenant_id，可能导致功能异常
   - 缓解: 兼容处理，惰性迁移，逐步过渡

### 低风险项

1. **代码变更风险**
   - 风险: 代码逻辑错误
   - 缓解: 单元测试，代码审查

---

## ✅ 验收标准

### 功能验收

- [ ] Session实体包含TenantID字段
- [ ] users表包含tenant_id列
- [ ] user_tenants表创建成功
- [ ] 租户隔离中间件可以从Session提取tenant_id
- [ ] 所有现有用户都有对应的租户
- [ ] 数据完整性验证通过

### 性能验收

- [ ] 登录接口响应时间 < 200ms
- [ ] 租户隔离中间件开销 < 1ms
- [ ] 数据库查询使用索引

### 安全验收

- [ ] 所有数据库查询包含tenant_id过滤
- [ ] 无法跨租户访问数据
- [ ] 审计日志包含tenant_id

---

## 📞 联系方式

- **负责人**: 研发B（后端工程师）
- **审核人**: 研发A（后端架构师）
- **最后更新**: 2025-01-01

---

**📅 预计完成时间**: 6-9天
