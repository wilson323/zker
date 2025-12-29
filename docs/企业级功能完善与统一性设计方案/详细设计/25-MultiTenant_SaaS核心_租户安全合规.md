# 25-MultiTenant_SaaS核心_租户安全合规 详细设计说明书

**文档编号**: DE-DD-2025-025
**模块名称**: 租户安全合规 (TenantSecurity)
**版本**: v1.0.0
**作者**: ZKER Enterprise Team
**创建日期**: 2025-01-03

---

## 1. 模块概述

**租户安全合规** 是 MultiTenant SaaS 的安全合规模块，通过**多层级安全防护 + 合规审计**，确保平台满足等保三级、GDPR、SOC2等合规要求。

**核心设计理念**：
- ✅ **多层防护**：身份认证、权限控制、数据加密、审计日志
- ✅ **合规认证**：等保三级、GDPR、SOC2
- ✅ **租户隔离**：严格的多租户数据隔离

**实现策略**：✅ 50% 编码（安全框架） + 50% 配置（合规规则）

---

## 2. 核心功能

### 2.1 身份认证

**F1 - 认证方式**
- F1.1 用户名密码登录
- F1.2 手机验证码登录
- F1.3 第三方登录（OAuth2：企业微信、钉钉、飞书）
- F1.4 SSO单点登录（SAML）
- F1.5 MFA多因素认证

**F2 - 会话管理**
- F2.1 JWT Token
- F2.2 Token刷新
- F2.3 强制登出
- F2.4 异地登录检测

### 2.2 权限控制

**F3 - RBAC权限模型**
- F3.1 角色管理
- F3.2 权限管理
- F3.3 用户授权
- F3.4 资源级权限

**F4 - 数据隔离**
- F4.1 行级隔离（Row-Level Security）
- F4.2 列级脱敏
- F4.3 审计日志

### 2.3 合规审计

**F5 - 审计日志**
- F5.1 操作日志（CRUD）
- F5.2 访问日志
- F5.3 敏感操作日志
- F5.4 日志归档

---

## 3. 数据库设计

### 3.1 核心表结构

#### 3.1.1 审计日志表 (audit_logs)

```sql
CREATE TABLE audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT COMMENT '用户ID',
    actor_type ENUM('user', 'system', 'api') NOT NULL COMMENT '操作者类型',
    actor_id VARCHAR(64) COMMENT '操作者ID',

    -- 操作信息
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型 (bot/message/user/...)',
    resource_id VARCHAR(64) COMMENT '资源ID',
    action VARCHAR(50) NOT NULL COMMENT '操作类型 (create/read/update/delete/login/...)',
    action_detail TEXT COMMENT '操作详情',

    -- 上下文
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT 'User-Agent',
    request_id VARCHAR(64) COMMENT '请求ID',

    -- 结果
    status ENUM('success', 'failure', 'partial') NOT NULL COMMENT '操作状态',
    error_code VARCHAR(50) COMMENT '错误码',
    error_message TEXT COMMENT '错误信息',

    -- 敏感级别
    sensitivity ENUM('low', 'medium', 'high', 'critical') DEFAULT 'low' COMMENT '敏感级别',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at_ms INT COMMENT '毫秒部分',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at),
    INDEX idx_sensitivity (sensitivity)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='审计日志表';
```

#### 3.1.2 登录日志表 (login_logs)

```sql
CREATE TABLE login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id VARCHAR(64) COMMENT '租户ID',
    user_id BIGINT COMMENT '用户ID',
    login_method ENUM('password', 'sms', 'oauth', 'saml') NOT NULL COMMENT '登录方式',

    -- 登录信息
    ip_address VARCHAR(45) NOT NULL COMMENT 'IP地址',
    location JSON COMMENT '地理位置 {"country":"CN","province":"Beijing"}',
    device_info JSON COMMENT '设备信息 {"type":"desktop","os":"Windows"}',
    user_agent VARCHAR(500) COMMENT 'User-Agent',

    -- 结果
    status ENUM('success', 'failure', 'blocked') NOT NULL COMMENT '登录状态',
    failure_reason VARCHAR(200) COMMENT '失败原因',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_user_id (user_id),
    INDEX idx_ip_address (ip_address),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='登录日志表';
```

#### 3.1.3 敏感操作审计表 (sensitive_operations)

```sql
CREATE TABLE sensitive_operations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '操作ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    operator_id BIGINT NOT NULL COMMENT '操作人ID',
    operator_type ENUM('user', 'admin', 'system') NOT NULL,

    operation_type ENUM('delete_bot', 'delete_user', 'export_data', 'change_permissions', 'access_sensitive_data') NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(64) NOT NULL,
    reason TEXT COMMENT '操作原因',
    approval_id VARCHAR(64) COMMENT '审批ID（如果需要审批）',

    ip_address VARCHAR(45),
    status ENUM('pending', 'approved', 'rejected', 'executed', 'failed') DEFAULT 'pending',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_operator_id (operator_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='敏感操作审计表';
```

---

## 4. 安全框架设计

### 4.1 JWT认证中间件

```go
package auth

// JWTAuthMiddleware JWT认证中间件
type JWTAuthMiddleware struct {
    jwtParser *jwt.Parser
    userRepo  repository.UserRepository
    logger    *zap.Logger
}

// Authenticate JWT认证
func (m *JWTAuthMiddleware) Authenticate(
    ctx context.Context,
    tokenString string,
) (*Claims, error) {
    // 1. 解析Token
    token, err := m.jwtParser.Parse(tokenString)
    if err != nil {
        return nil, err
    }

    // 2. 验证签名
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        // 3. 检查Token是否过期
        if exp, ok := claims["exp"].(float64); ok {
            if time.Now().Unix() > int64(exp) {
                return nil, errors.New("token expired")
            }
        }

        // 4. 检查用户是否存在且活跃
        userID := int64(claims["user_id"].(float64))
        user, err := m.userRepo.GetByID(ctx, userID)
        if err != nil || user.Status != "active" {
            return nil, errors.New("user not found or inactive")
        }

        // 5. 记录审计日志
        m.recordAuditLog(ctx, &AuditLog{
            TenantID:    user.TenantID,
            UserID:      userID,
            ActorType:   "user",
            ActionType:  "authenticate",
            Status:      "success",
            IPAddress:   getIPAddress(ctx),
        })

        return &Claims{
            TenantID: user.TenantID,
            UserID:   userID,
            Role:     claims["role"].(string),
        }, nil
    }

    return nil, errors.New("invalid token")
}
```

### 4.2 RBAC权限检查

```go
package auth

// RBACChecker RBAC权限检查器
type RBACChecker struct {
    roleRepo       repository.RoleRepository
    permissionRepo repository.PermissionRepository
}

// CheckPermission 检查权限
func (c *RBACChecker) CheckPermission(
    ctx context.Context,
    userID int64,
    resourceType, resourceID, permission string,
) (bool, error) {
    // 1. 获取用户角色
    roles, err := c.roleRepo.GetByUserID(ctx, userID)
    if err != nil {
        return false, err
    }

    // 2. 检查每个角色的权限
    for _, role := range roles {
        // 检查角色是否有该资源的权限
        perm, err := c.permissionRepo.GetByRoleAndResource(
            ctx, role.ID, resourceType, resourceID,
        )
        if err == nil && perm != nil {
            // 检查具体权限
            if hasPermission(perm.Permissions, permission) {
                return true, nil
            }
        }
    }

    return false, nil
}

// hasPermission 检查权限列表
func hasPermission(permissions []string, required string) bool {
    for _, perm := range permissions {
        if perm == "*" || perm == required {
            return true
        }
    }
    return false
}
```

### 4.3 数据脱敏

```go
package security

// DataMasking 数据脱敏
type DataMasking struct{}

// MaskEmail 邮箱脱敏
func (m *DataMasking) MaskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }

    username := parts[0]
    domain := parts[1]

    if len(username) <= 3 {
        return strings.Repeat("*", len(username)) + "@" + domain
    }

    return username[:2] + strings.Repeat("*", len(username)-2) + "@" + domain
}

// MaskPhone 手机号脱敏
func (m *DataMasking) MaskPhone(phone string) string {
    if len(phone) != 11 {
        return phone
    }

    return phone[:3] + "****" + phone[7:]
}

// MaskIDCard 身份证脱敏
func (m *DataMasking) MaskIDCard(idcard string) string {
    if len(idcard) != 18 {
        return idcard
    }

    return idcard[:6] + "********" + idcard[14:]
}
```

---

## 5. 合规设计

### 5.1 等保三级合规

**技术要求**：
- ✅ 身份鉴别：双因素认证、密码复杂度策略
- ✅ 访问控制：RBAC权限模型、最小权限原则
- ✅ 安全审计：完整的审计日志（保留6个月以上）
- ✅ 数据加密：传输加密（TLS）、存储加密（AES-256）
- ✅ 入侵防范：WAF、IDS、漏洞扫描

**实施清单**：

```sql
-- 1. 密码复杂度策略
ALTER TABLE users ADD COLUMN password_changed_at DATETIME;
CREATE TRIGGER check_password_complexity
  BEFORE INSERT ON users
  FOR EACH ROW
BEGIN
  -- 密码长度至少8位，包含大小写字母、数字、特殊字符
  -- 每90天必须修改密码
END;

-- 2. 审计日志保留策略
CREATE EVENT retention_audit_logs
ON SCHEDULE EVERY 1 DAY
DO
  DELETE FROM audit_logs WHERE created_at < DATE_SUB(NOW(), INTERVAL 2 YEAR);

-- 3. 敏感操作审批流程
INSERT INTO sensitive_operations_config (operation_type, requires_approval) VALUES
('delete_bot', true),
('export_all_data', true),
('change_permissions', true);
```

### 5.2 GDPR合规

**核心要求**：
- ✅ 数据最小化：仅收集必要数据
- ✅ 用户权利：数据访问、删除、可携带权
- ✅ 数据保护：加密、匿名化
- ✅ 违规通知：72小时内通知数据泄露

**实施清单**：

```sql
-- 1. 用户同意记录表
CREATE TABLE user_consents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    consent_type ENUM('privacy_policy', 'terms_of_service', 'data_processing') NOT NULL,
    granted BOOLEAN NOT NULL,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    withdrawn_at DATETIME,

    UNIQUE KEY uk_user_type (user_id, consent_type),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='用户同意记录表';

-- 2. 数据删除请求表
CREATE TABLE data_deletion_requests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    request_type ENUM('account_deletion', 'data_erasure') NOT NULL,
    status ENUM('pending', 'processing', 'completed', 'rejected') DEFAULT 'pending',
    requested_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,

    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='数据删除请求表';
```

---

## 6. API 设计

| 方法 | 路径 | 功能 |
|------|------|------|
| **认证** |
| POST | /api/v1/auth/login | 用户登录 |
| POST | /api/v1/auth/logout | 用户登出 |
| POST | /api/v1/auth/refresh | 刷新Token |
| **审计** |
| GET | /api/v1/audit/logs | 查询审计日志 |
| GET | /api/v1/audit/login-logs | 查询登录日志 |
| GET | /api/v1/audit/sensitive-operations | 查询敏感操作 |
| **合规** |
| POST | /api/v1/gdpr/data-export | 导出个人数据 |
| POST | /api/v1/gdpr/data-delete | 删除个人数据 |
| GET | /api/v1/gdpr/consents | 查看用户同意记录 |

---

## 7. 总结

### 7.1 实施策略总结

| 实施项 | 实施方式 | 工作量 |
|-------|---------|--------|
| **安全框架** | 💻 独立编码 | 50% |
| **合规规则** | 📊 配置驱动 | 50% |

**总计**：50% 编码 + 50% 配置

### 7.2 核心优势

- ✅ **多层防护**：认证、授权、加密、审计
- ✅ **合规认证**：等保三级、GDPR、SOC2
- ✅ **租户隔离**：严格的多租户数据隔离

---

**文档结束**
