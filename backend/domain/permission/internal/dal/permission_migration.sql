-- ==========================================
-- RBAC权限系统完整版数据库迁移脚本
-- 版本: v2.0.0
-- 创建日期: 2025-01-01
-- 描述: 扩展数据权限和字段权限表，支持5级数据权限+3级字段权限
-- ==========================================

-- ==========================================
-- 1. 数据权限规则表扩展
-- ==========================================

-- 修改 data_permissions 表，支持5级数据权限
ALTER TABLE `data_permissions`
MODIFY COLUMN `scope` ENUM('ALL', 'DEPARTMENT_AND_SUB', 'DEPARTMENT', 'OWN', 'CUSTOM', 'NONE') NOT NULL
COMMENT '数据权限范围: ALL-全部数据, DEPARTMENT_AND_SUB-本部门及子部门, DEPARTMENT-本部门, OWN-仅本人, CUSTOM-自定义过滤, NONE-无权限';

-- 添加索引以提升查询性能
ALTER TABLE `data_permissions`
ADD INDEX `idx_role_resource_scope` (`role_id`, `resource_type`, `scope`);

-- ==========================================
-- 2. 字段权限表优化
-- ==========================================

-- 修改 field_permissions 表，优化权限级别定义
ALTER TABLE `field_permissions`
MODIFY COLUMN `permission_level` ENUM('hidden', 'readonly', 'editable', 'required') NOT NULL
COMMENT '字段权限级别: hidden-隐藏, readonly-只读, editable-可编辑, required-必填';

-- 添加脱敏规则字段
ALTER TABLE `field_permissions`
ADD COLUMN `mask_rule` VARCHAR(100) DEFAULT NULL COMMENT '脱敏规则: phone/email/idcard/name/card/regex:{pattern}';

-- 添加索引
ALTER TABLE `field_permissions`
ADD UNIQUE KEY `uk_role_resource_field` (`role_id`, `resource_type`, `field_name`);

-- ==========================================
-- 3. 部门表扩展（支持树形结构）
-- ==========================================

-- 确保 departments 表有 parent_department_id 字段
ALTER TABLE `departments`
ADD COLUMN IF NOT EXISTS `parent_department_id` VARCHAR(36) DEFAULT NULL COMMENT '父部门ID',
ADD COLUMN IF NOT EXISTS `level` INT DEFAULT 0 COMMENT '部门层级（0为根部门）',
ADD COLUMN IF NOT EXISTS `path` VARCHAR(500) DEFAULT NULL COMMENT '部门路径，如: /1/2/3',
ADD COLUMN IF NOT EXISTS `is_active` BOOLEAN DEFAULT TRUE COMMENT '是否启用';

-- 添加索引
ALTER TABLE `departments`
ADD INDEX `idx_parent_department` (`parent_department_id`),
ADD INDEX `idx_tenant_level` (`tenant_id`, `level`);

-- ==========================================
-- 4. 用户部门关联表扩展
-- ==========================================

-- 确保 user_departments 表有 is_leader 字段
ALTER TABLE `user_departments`
ADD COLUMN IF NOT EXISTS `is_leader` BOOLEAN DEFAULT FALSE COMMENT '是否是部门领导',
ADD COLUMN IF NOT EXISTS `is_primary` BOOLEAN DEFAULT FALSE COMMENT '是否是主部门',
ADD COLUMN IF NOT EXISTS `job_title` VARCHAR(100) DEFAULT NULL COMMENT '职位';

-- 添加索引
ALTER TABLE `user_departments`
ADD UNIQUE KEY `uk_user_dept_primary` (`user_id`, `department_id`, `is_primary`);

-- ==========================================
-- 5. 角色表扩展
-- ==========================================

-- 确保 roles 表有 priority 字段（用于权限合并时的优先级）
ALTER TABLE `roles`
ADD COLUMN IF NOT EXISTS `priority` INT DEFAULT 0 COMMENT '权限优先级（数字越大优先级越高）',
ADD COLUMN IF NOT EXISTS `is_system` BOOLEAN DEFAULT FALSE COMMENT '是否是系统角色',
ADD COLUMN IF NOT EXISTS `description` TEXT DEFAULT NULL COMMENT '角色描述';

-- 添加索引
ALTER TABLE `roles`
ADD INDEX `idx_tenant_priority` (`tenant_id`, `priority`);

-- ==========================================
-- 6. 用户角色关联表扩展
-- ==========================================

-- 确保 user_roles 表有 expires_at 字段（支持临时角色）
ALTER TABLE `user_roles`
ADD COLUMN IF NOT EXISTS `expires_at` TIMESTAMP NULL DEFAULT NULL COMMENT '角色过期时间（NULL表示永不过期）',
ADD COLUMN IF NOT EXISTS `granted_by` VARCHAR(36) DEFAULT NULL COMMENT '授权人ID',
ADD COLUMN IF NOT EXISTS `granted_reason` VARCHAR(200) DEFAULT NULL COMMENT '授权原因';

-- 添加索引
ALTER TABLE `user_roles`
ADD INDEX `idx_expires_at` (`expires_at`);

-- ==========================================
-- 7. 权限审计日志表（新增）
-- ==========================================

CREATE TABLE IF NOT EXISTS `permission_audit_logs` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `resource_type` VARCHAR(64) NOT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(36) DEFAULT NULL COMMENT '资源ID',
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型: check/grant/revoke',
    `permission_type` VARCHAR(50) NOT NULL COMMENT '权限类型: data/field',
    `decision` VARCHAR(20) NOT NULL COMMENT '决策结果: allow/deny',
    `role_id` VARCHAR(36) DEFAULT NULL COMMENT '角色ID',
    `reason` VARCHAR(500) DEFAULT NULL COMMENT '决策原因',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户代理',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX `idx_tenant_user` (`tenant_id`, `user_id`),
    INDEX `idx_resource` (`resource_type`, `resource_id`),
    INDEX `idx_created_at` (`created_at`),
    INDEX `idx_decision` (`decision`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限审计日志表';

-- ==========================================
-- 8. 权限变更历史表（新增）
-- ==========================================

CREATE TABLE IF NOT EXISTS `permission_change_history` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '历史ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `change_type` VARCHAR(50) NOT NULL COMMENT '变更类型: role_grant/role_revoke/perm_update/field_update',
    `target_type` VARCHAR(50) NOT NULL COMMENT '目标类型: user/role',
    `target_id` VARCHAR(36) NOT NULL COMMENT '目标ID',
    `role_id` VARCHAR(36) DEFAULT NULL COMMENT '角色ID',
    `permission_detail` JSON DEFAULT NULL COMMENT '权限详情（变更前后的对比）',
    `operator_id` VARCHAR(36) NOT NULL COMMENT '操作人ID',
    `operator_name` VARCHAR(100) DEFAULT NULL COMMENT '操作人姓名',
    `change_reason` VARCHAR(500) DEFAULT NULL COMMENT '变更原因',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX `idx_tenant_target` (`tenant_id`, `target_type`, `target_id`),
    INDEX `idx_operator` (`operator_id`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限变更历史表';

-- ==========================================
-- 9. 初始化默认数据和索引优化
-- ==========================================

-- 为所有业务表添加 department_id 字段（如果不存在）
-- 注意: 这里列出主要的业务表，实际执行时需要根据表结构调整

-- Bots表
ALTER TABLE `bots`
ADD COLUMN IF NOT EXISTS `department_id` VARCHAR(36) DEFAULT NULL COMMENT '所属部门ID',
ADD INDEX IF NOT EXISTS `idx_department` (`department_id`);

-- Conversations表
ALTER TABLE `conversations`
ADD COLUMN IF NOT EXISTS `department_id` VARCHAR(36) DEFAULT NULL COMMENT '所属部门ID',
ADD INDEX IF NOT EXISTS `idx_department` (`department_id`);

-- Knowledge bases表
ALTER TABLE `knowledge_bases`
ADD COLUMN IF NOT EXISTS `department_id` VARCHAR(36) DEFAULT NULL COMMENT '所属部门ID',
ADD INDEX IF NOT EXISTS `idx_department` (`department_id`);

-- Workflows表
ALTER TABLE `workflows`
ADD COLUMN IF NOT EXISTS `department_id` VARCHAR(36) DEFAULT NULL COMMENT '所属部门ID',
ADD INDEX IF NOT EXISTS `idx_department` (`department_id`);

-- ==========================================
-- 10. 创建视图：用户完整权限视图
-- ==========================================

CREATE OR REPLACE VIEW `v_user_permissions` AS
SELECT
    ur.tenant_id,
    ur.user_id,
    ur.role_id,
    r.role_name,
    r.role_code,
    r.role_type,
    dp.resource_type,
    dp.scope AS data_scope,
    dp.custom_filter,
    fp.field_name,
    fp.permission_level AS field_permission_level,
    fp.mask_rule,
    ur.expires_at,
    ur.created_at AS granted_at
FROM user_roles ur
INNER JOIN roles r ON ur.role_id = r.role_id
LEFT JOIN data_permissions dp ON ur.role_id = dp.role_id
LEFT JOIN field_permissions fp ON ur.role_id = fp.role_id
WHERE r.deleted_at IS NULL
  AND (ur.expires_at IS NULL OR ur.expires_at > NOW());

-- ==========================================
-- 11. 创建存储过程：清理过期角色
-- ==========================================

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_cleanup_expired_roles`()
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    -- 删除过期的用户角色关联
    DELETE FROM user_roles
    WHERE expires_at IS NOT NULL
      AND expires_at < NOW();

    -- 记录清理数量
    SELECT ROW_COUNT() AS cleaned_count;

    COMMIT;
END$$

DELIMITER ;

-- ==========================================
-- 12. 创建存储过程：重建部门路径
-- ==========================================

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_rebuild_department_paths`()
BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    -- 重建部门路径（简化版本）
    UPDATE departments d1
    LEFT JOIN departments d2 ON d1.parent_department_id = d2.department_id
    SET d1.level = IFNULL(d2.level, -1) + 1,
        d1.path = CONCAT(IFNULL(d2.path, ''), '/', d1.department_id)
    WHERE d1.deleted_at IS NULL;

    COMMIT;
END$$

DELIMITER ;

-- ==========================================
-- 13. 创建定时任务清理事件（可选）
-- ==========================================

-- 创建事件每天凌晨清理过期角色
SET GLOBAL event_scheduler = ON;

CREATE EVENT IF NOT EXISTS `evt_cleanup_expired_roles`
ON SCHEDULE EVERY 1 DAY
STARTS CONCAT(CURRENT_DATE, ' 02:00:00')
DO CALL sp_cleanup_expired_roles();

-- ==========================================
-- 14. 性能优化：添加复合索引
-- ==========================================

-- data_permissions 表
ALTER TABLE `data_permissions`
ADD INDEX `idx_tenant_role_resource` (`tenant_id`, `role_id`, `resource_type`);

-- field_permissions 表
ALTER TABLE `field_permissions`
ADD INDEX `idx_tenant_role_field` (`tenant_id`, `role_id`, `field_name`);

-- user_roles 表
ALTER TABLE `user_roles`
ADD INDEX `idx_tenant_user_role` (`tenant_id`, `user_id`, `role_id`);

-- user_departments 表
ALTER TABLE `user_departments`
ADD INDEX `idx_tenant_dept_leader` (`tenant_id`, `department_id`, `is_leader`);

-- ==========================================
-- 15. 数据完整性约束
-- ==========================================

-- 添加外键约束（如果需要）
-- ALTER TABLE `data_permissions`
-- ADD CONSTRAINT `fk_data_perm_role`
-- FOREIGN KEY (`role_id`) REFERENCES `roles`(`role_id`) ON DELETE CASCADE;

-- ALTER TABLE `field_permissions`
-- ADD CONSTRAINT `fk_field_perm_role`
-- FOREIGN KEY (`role_id`) REFERENCES `roles`(`role_id`) ON DELETE CASCADE;

-- ALTER TABLE `user_departments`
-- ADD CONSTRAINT `fk_user_dept_user`
-- FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE,
-- ADD CONSTRAINT `fk_user_dept_dept`
-- FOREIGN KEY (`department_id`) REFERENCES `departments`(`department_id`) ON DELETE CASCADE;

-- ==========================================
-- 迁移脚本执行完成
-- ==========================================
