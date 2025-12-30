-- =====================================================
-- 组织中心增强功能 - 数据库表结构
-- 版本: v1.0
-- 创建日期: 2025-01-01
-- 说明: 包含员工档案扩展、虚拟组织、矩阵组织等增强功能
-- =====================================================

-- =====================================================
-- 1. 员工档案扩展表
-- =====================================================

-- 工作经历表
CREATE TABLE IF NOT EXISTS `work_experiences` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `company_name` VARCHAR(255) NOT NULL COMMENT '公司名称',
    `position` VARCHAR(128) NOT NULL COMMENT '职位',
    `start_date` DATE NOT NULL COMMENT '开始日期',
    `end_date` DATE DEFAULT NULL COMMENT '结束日期(空表示当前在职)',
    `description` TEXT COMMENT '工作描述',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间(软删除)',

    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_start_date` (`start_date`),
    INDEX `idx_end_date` (`end_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作经历表';

-- 教育经历表
CREATE TABLE IF NOT EXISTS `educations` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `school_name` VARCHAR(255) NOT NULL COMMENT '学校名称',
    `major` VARCHAR(128) NOT NULL COMMENT '专业',
    `degree` ENUM('HIGH_SCHOOL', 'BACHELOR', 'MASTER', 'PHD') NOT NULL COMMENT '学位',
    `start_date` DATE NOT NULL COMMENT '开始日期',
    `end_date` DATE DEFAULT NULL COMMENT '结束日期(空表示当前在读)',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间(软删除)',

    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_start_date` (`start_date`),
    INDEX `idx_end_date` (`end_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教育经历表';

-- 成员技能表
CREATE TABLE IF NOT EXISTS `member_skills` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `skill_name` VARCHAR(128) NOT NULL COMMENT '技能名称',
    `proficiency` ENUM('BEGINNER', 'INTERMEDIATE', 'ADVANCED', 'EXPERT') NOT NULL COMMENT '熟练度',
    `certified` BOOLEAN DEFAULT FALSE COMMENT '是否认证',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间(软删除)',

    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_skill_name` (`skill_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成员技能表';

-- =====================================================
-- 2. 虚�拟组织表
-- =====================================================

-- 虚拟组织表
CREATE TABLE IF NOT EXISTS `virtual_organizations` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `virtual_org_id` VARCHAR(36) NOT NULL UNIQUE COMMENT '虚拟组织ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `name` VARCHAR(128) NOT NULL COMMENT '组织名称',
    `description` TEXT DEFAULT NULL COMMENT '组织描述',
    `type` ENUM('PROJECT', 'TASK_FORCE', 'COMMUNITY', 'OTHER') NOT NULL COMMENT '组织类型',
    `owner_id` VARCHAR(36) NOT NULL COMMENT '所有者ID',
    `status` ENUM('ACTIVE', 'ARCHIVED') DEFAULT 'ACTIVE' COMMENT '状态',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间(软删除)',

    INDEX `idx_virtual_org_id` (`virtual_org_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_owner_id` (`owner_id`),
    INDEX `idx_type` (`type`),
    INDEX `idx_status` (`status`),
    UNIQUE KEY `uk_virtual_org_id` (`virtual_org_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='虚拟组织表';

-- 虚拟组织成员表
CREATE TABLE IF NOT EXISTS `virtual_org_members` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `virtual_org_id` VARCHAR(36) NOT NULL COMMENT '虚拟组织ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `role` ENUM('OWNER', 'ADMIN', 'MEMBER') NOT NULL COMMENT '角色',
    `joined_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',

    UNIQUE KEY `uk_org_user` (`virtual_org_id`, `user_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='虚拟组织成员表';

-- 虚拟组织标签表
CREATE TABLE IF NOT EXISTS `virtual_org_tags` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `virtual_org_id` VARCHAR(36) NOT NULL COMMENT '虚拟组织ID',
    `tag_name` VARCHAR(64) NOT NULL COMMENT '标签名称',

    UNIQUE KEY `uk_org_tag` (`virtual_org_id`, `tag_name`),
    INDEX `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='虚拟组织标签表';

-- =====================================================
-- 3. 矩阵组织表
-- =====================================================

-- 矩阵汇报关系表
CREATE TABLE IF NOT EXISTS `matrix_reportings` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(36) NOT NULL COMMENT '用户ID',
    `reporting_type` ENUM('FUNCTIONAL', 'PROJECT', 'DOTTED') NOT NULL COMMENT '汇报关系类型',
    `supervisor_id` VARCHAR(36) NOT NULL COMMENT '上级ID',
    `organization_id` VARCHAR(36) NOT NULL COMMENT '组织ID',
    `effective_date` DATE NOT NULL COMMENT '生效日期',
    `expiry_date` DATE DEFAULT NULL COMMENT '失效日期',
    `is_primary` BOOLEAN DEFAULT FALSE COMMENT '是否主要汇报关系',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_supervisor_id` (`supervisor_id`),
    INDEX `idx_org_id` (`organization_id`),
    INDEX `idx_reporting_type` (`reporting_type`),
    INDEX `idx_effective_date` (`effective_date`),
    INDEX `idx_expiry_date` (`expiry_date`),
    INDEX `idx_is_primary` (`is_primary`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='矩阵汇报关系表';

-- =====================================================
-- 4. 初始化数据
-- =====================================================

-- 创建示例虚拟组织类型枚举说明
-- 这仅用于文档说明，不作为实际SQL执行

-- PROJECT: 项目组 - 临时性项目团队
-- TASK_FORCE: 任务组 - 针对特定任务的突击队
-- COMMUNITY: 社区 - 兴趣社区或实践社区
-- OTHER: 其他 - 未分类的虚拟组织

-- 创建示例汇报关系类型枚举说明
-- 这仅用于文档说明，不作为实际SQL执行

-- FUNCTIONAL: 职能汇报 - 传统直线汇报关系
-- PROJECT: 项目汇报 - 项目期间的横向汇报
-- DOTTED: 虚线汇报 - 非正式的咨询关系

-- =====================================================
-- 5. 性能优化建议
-- =====================================================

-- 1. 根据实际查询模式，可能需要添加复合索引
-- 例如: ALTER TABLE work_experiences ADD INDEX idx_user_start (user_id, start_date);

-- 2. 对于大表，考虑分区策略
-- 例如: 按tenant_id或created_at分区

-- 3. 对于全文搜索需求，可考虑添加全文索引
-- 例如: ALTER TABLE member_skills ADD FULLTEXT INDEX ft_skill_name (skill_name);

-- =====================================================
-- 6. 数据迁移说明
-- =====================================================

-- 如果需要从旧系统迁移数据，请参考以下步骤:

-- 1. 备份现有数据
-- mysqldump -u root -p database_name > backup_$(date +%Y%m%d).sql

-- 2. 执行表结构创建
-- mysql -u root -p database_name < enhanced_org_schema.sql

-- 3. 运行数据迁移脚本(需单独编写)
-- mysql -u root -p database_name < migration_script.sql

-- 4. 验证数据完整性
-- SELECT COUNT(*) FROM work_experiences;
-- SELECT COUNT(*) FROM educations;
-- SELECT COUNT(*) FROM member_skills;
-- SELECT COUNT(*) FROM virtual_organizations;
-- SELECT COUNT(*) FROM virtual_org_members;
-- SELECT COUNT(*) FROM matrix_reportings;

-- =====================================================
-- 7. 权限设置
-- =====================================================

-- 创建只读用户(用于报表查询)
-- CREATE USER 'org_readonly'@'%' IDENTIFIED BY 'password';
-- GRANT SELECT ON database_name.work_experiences TO 'org_readonly'@'%';
-- GRANT SELECT ON database_name.educations TO 'org_readonly'@'%';
-- GRANT SELECT ON database_name.member_skills TO 'org_readonly'@'%';
-- GRANT SELECT ON database_name.virtual_organizations TO 'org_readonly'@'%';
-- GRANT SELECT ON database_name.virtual_org_members TO 'org_readonly'@'%';
-- GRANT SELECT ON database_name.matrix_reportings TO 'org_readonly'@'%';
-- FLUSH PRIVILEGES;

-- =====================================================
-- 8. 监控查询
-- =====================================================

-- 查看表大小
-- SELECT
--     table_name AS 'Table',
--     ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)'
-- FROM information_schema.TABLES
-- WHERE table_schema = 'database_name'
--     AND table_name IN ('work_experiences', 'educations', 'member_skills',
--                        'virtual_organizations', 'virtual_org_members',
--                        'matrix_reportings')
-- ORDER BY (data_length + index_length) DESC;

-- 查看索引使用情况
-- SHOW INDEX FROM work_experiences;
-- SHOW INDEX FROM educations;
-- SHOW INDEX FROM member_skills;
-- SHOW INDEX FROM virtual_organizations;
-- SHOW INDEX FROM virtual_org_members;
-- SHOW INDEX FROM matrix_reportings;

-- =====================================================
-- End of SQL Schema
-- =====================================================
