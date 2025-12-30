-- ============================================
-- Session/User 表 tenant_id 迁移 DDL
-- 执行时机: 维护窗口（建议凌晨2:00-4:00）
-- 执行时长: 预计5-10分钟
-- 最后更新: 2025-12-30
-- ============================================
--
-- ⚠️ 重要提示:
-- 1. Session表和User表需要添加tenant_id字段
-- 2. 创建user_tenant关联表支持用户多租户
-- 3. 先添加字段为NULL，迁移数据后改为NOT NULL
-- 4. 使用ALGORITHM=INPLACE降低锁影响
--
-- ============================================

-- ============================================
-- 步骤1: 为 session 表添加 tenant_id
-- ============================================
ALTER TABLE session
ADD COLUMN tenant_id VARCHAR(64) NULL
    COMMENT '租户ID（多租户隔离）'
    AFTER user_id,
ADD INDEX idx_session_tenant (tenant_id),
ADD INDEX idx_session_user_tenant (user_id, tenant_id);

-- 验证
SHOW COLUMNS FROM session LIKE 'tenant_id';
SHOW INDEX FROM session WHERE Key_name = 'idx_session_tenant';


-- ============================================
-- 步骤2: 为 users 表添加 tenant_id
-- ============================================
ALTER TABLE `user`
ADD COLUMN tenant_id VARCHAR(64) NULL
    COMMENT '租户ID（多租户隔离）'
    AFTER id,
ADD INDEX idx_user_tenant (tenant_id),
ADD INDEX idx_user_email_tenant (email, tenant_id);

-- 验证
SHOW COLUMNS FROM `user` LIKE 'tenant_id';
SHOW INDEX FROM `user` WHERE Key_name = 'idx_user_tenant';


-- ============================================
-- 步骤3: 创建 user_tenant 关联表
-- ============================================
CREATE TABLE IF NOT EXISTS user_tenant (
    id BIGINT PRIMARY KEY AUTO_INCREMENT
        COMMENT '主键ID',
    user_id BIGINT NOT NULL
        COMMENT '用户ID',
    tenant_id VARCHAR(64) NOT NULL
        COMMENT '租户ID',
    role VARCHAR(64) DEFAULT 'member'
        COMMENT '角色: owner, admin, member',
    is_default TINYINT(1) DEFAULT 0
        COMMENT '是否为默认租户',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        COMMENT '更新时间',
    deleted_at DATETIME DEFAULT NULL
        COMMENT '删除时间（软删除）',

    UNIQUE KEY uk_user_tenant (user_id, tenant_id, deleted_at),
    KEY idx_user_tenant_user (user_id),
    KEY idx_user_tenant_tenant (tenant_id),
    KEY idx_user_tenant_default (user_id, is_default, deleted_at)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci
COMMENT='用户-租户关联表（支持用户多租户）';

-- 验证
SHOW CREATE TABLE user_tenant;


-- ============================================
-- 步骤4: 回填 session 表的 tenant_id
-- ============================================
-- 策略: 通过user_id关联users表获取tenant_id
UPDATE session s
INNER JOIN `user` u ON s.user_id = u.id
SET s.tenant_id = u.tenant_id
WHERE s.tenant_id IS NULL;

-- 验证回填结果
SELECT
    COUNT(*) as total_sessions,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_count,
    SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_tenant_count
FROM session;


-- ============================================
-- 步骤5: 回填 users 表的 tenant_id
-- ============================================
-- 策略: 使用默认租户ID 'default'
-- TODO: 后续需要根据实际业务规则分配租户ID
UPDATE `user`
SET tenant_id = 'default'
WHERE tenant_id IS NULL;

-- 验证回填结果
SELECT
    COUNT(*) as total_users,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_count,
    SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_tenant_count
FROM `user`;


-- ============================================
-- 步骤6: 初始化 user_tenant 关联数据
-- ============================================
-- 为所有用户创建默认租户关联
INSERT INTO user_tenant (user_id, tenant_id, role, is_default, created_at, updated_at)
SELECT
    id as user_id,
    'default' as tenant_id,
    'owner' as role,
    1 as is_default,
    NOW() as created_at,
    NOW() as updated_at
FROM `user`
WHERE tenant_id = 'default'
  AND NOT EXISTS (
    SELECT 1 FROM user_tenant
    WHERE user_tenant.user_id = `user`.id
      AND user_tenant.tenant_id = 'default'
      AND user_tenant.deleted_at IS NULL
  );

-- 验证关联数据
SELECT
    COUNT(*) as total_user_tenants,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT tenant_id) as unique_tenants
FROM user_tenant
WHERE deleted_at IS NULL;


-- ============================================
-- 步骤7: 修改字段为 NOT NULL（可选）
-- ============================================
-- ⚠️ 仅在确认所有数据都已正确回填后执行
-- ALTER TABLE session MODIFY COLUMN tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID（多租户隔离）';
-- ALTER TABLE `user` MODIFY COLUMN tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID（多租户隔离）';


-- ============================================
-- 执行后验证脚本
-- ============================================

-- 1. 检查所有表是否都有tenant_id字段
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    COLUMN_TYPE,
    IS_NULLABLE,
    COLUMN_DEFAULT,
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND COLUMN_NAME = 'tenant_id'
  AND TABLE_NAME IN ('session', 'user', 'user_tenant')
ORDER BY TABLE_NAME;

-- 2. 检查所有索引是否创建
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    SEQ_IN_INDEX,
    INDEX_TYPE
FROM INFORMATION_SCHEMA.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN ('session', 'user', 'user_tenant')
  AND INDEX_NAME LIKE '%tenant%'
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;

-- 3. 检查数据完整性
SELECT
    'session' as table_name,
    COUNT(*) as total_records,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_count,
    SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_tenant_count,
    ROUND(SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) as completion_rate
FROM session
UNION ALL
SELECT
    'user' as table_name,
    COUNT(*) as total_records,
    SUM(CASE WHEN tenant_id IS NULL THEN 1 ELSE 0 END) as null_tenant_count,
    SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) as has_tenant_count,
    ROUND(SUM(CASE WHEN tenant_id IS NOT NULL THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2) as completion_rate
FROM `user`;

-- 4. 检查user_tenant关联完整性
SELECT
    COUNT(*) as total_relations,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT tenant_id) as unique_tenants,
    SUM(CASE WHEN is_default = 1 THEN 1 ELSE 0 END) as default_tenant_count
FROM user_tenant
WHERE deleted_at IS NULL;

-- 5. 检查孤儿数据（user在user表中不存在）
SELECT
    ut.user_id,
    ut.tenant_id,
    ut.role
FROM user_tenant ut
LEFT JOIN `user` u ON ut.user_id = u.id
WHERE u.id IS NULL
  AND ut.deleted_at IS NULL;


-- ============================================
-- 回滚脚本（⚠️ 仅在迁移失败时使用）
-- ============================================

/*
-- 回滚步骤1: 删除user_tenant表
DROP TABLE IF EXISTS user_tenant;

-- 回滚步骤2: 删除users表的tenant_id字段
ALTER TABLE `user` DROP COLUMN tenant_id;

-- 回滚步骤3: 删除session表的tenant_id字段
ALTER TABLE session DROP COLUMN tenant_id;
*/


-- ============================================
-- 执行完成标记
-- ============================================
-- 执行完成后，请在此处记录执行时间和结果
-- 执行时间: _______________
-- 执行结果: □ 成功  □ 失败
-- 备注说明: ___________________________________
-- 验证通过: □ 是  □ 否
-- ============================================
