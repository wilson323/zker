-- ================================================================================
-- 003_create_tenant_isolation_tables.sql
-- 创建租户隔离策略表
-- ================================================================================
-- 用途: 支持Schema级和Database级隔离
-- ================================================================================

-- 租户隔离策略表
CREATE TABLE IF NOT EXISTS `tenant_isolation_policies` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `isolation_strategy` ENUM('row_level', 'schema_level', 'database_level') NOT NULL DEFAULT 'row_level' COMMENT '隔离策略',
    `schema_name` VARCHAR(64) NULL COMMENT 'Schema名称（schema_level时使用）',
    `db_instance_id` INT NULL COMMENT '数据库实例ID（database_level时使用）',
    `db_connection_config` JSON NULL COMMENT '数据库连接配置（database_level时使用）',
    `status` ENUM('pending', 'active', 'migrating', 'failed') NOT NULL DEFAULT 'pending' COMMENT '状态',
    `migration_progress` INT DEFAULT 0 COMMENT '迁移进度（0-100）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_tenant_id` (`tenant_id`),
    INDEX `idx_isolation_strategy` (`isolation_strategy`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户隔离策略表';

-- Schema隔离映射表
CREATE TABLE IF NOT EXISTS `tenant_schema_mappings` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `schema_name` VARCHAR(64) NOT NULL COMMENT 'Schema名称',
    `db_instance_id` INT NOT NULL COMMENT '数据库实例ID',
    `status` ENUM('creating', 'active', 'migrating', 'archived', 'failed') NOT NULL DEFAULT 'creating' COMMENT '状态',
    `table_count` INT DEFAULT 0 COMMENT '表数量',
    `data_size_mb` BIGINT DEFAULT 0 COMMENT '数据大小（MB）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_tenant_id` (`tenant_id`),
    UNIQUE KEY `uk_schema_name` (`schema_name`),
    INDEX `idx_db_instance` (`db_instance_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Schema隔离映射表';

-- 数据库实例表（用于database_level隔离）
CREATE TABLE IF NOT EXISTS `database_instances` (
    `id` INT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `instance_name` VARCHAR(100) NOT NULL COMMENT '实例名称',
    `host` VARCHAR(255) NOT NULL COMMENT '主机地址',
    `port` INT NOT NULL DEFAULT 3306 COMMENT '端口',
    `database_name` VARCHAR(100) NOT NULL COMMENT '数据库名',
    `username` VARCHAR(100) NOT NULL COMMENT '用户名',
    `password_encrypted` TEXT NOT NULL COMMENT '加密密码',
    `max_connections` INT DEFAULT 100 COMMENT '最大连接数',
    `status` ENUM('provisioning', 'active', 'maintenance', 'decommissioned') NOT NULL DEFAULT 'provisioning' COMMENT '状态',
    `capacity_gb` BIGINT DEFAULT 100 COMMENT '容量（GB）',
    `used_gb` BIGINT DEFAULT 0 COMMENT '已使用（GB）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_instance_name` (`instance_name`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库实例表';

-- 租户数据库绑定表
CREATE TABLE IF NOT EXISTS `tenant_database_bindings` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `db_instance_id` INT NOT NULL COMMENT '数据库实例ID',
    `binding_type` ENUM('exclusive', 'shared') NOT NULL DEFAULT 'exclusive' COMMENT '绑定类型',
    `priority` INT NOT NULL DEFAULT 0 COMMENT '优先级',
    `status` ENUM('pending', 'active', 'suspended', 'terminated') NOT NULL DEFAULT 'pending' COMMENT '状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_tenant_db` (`tenant_id`, `db_instance_id`),
    INDEX `idx_db_instance` (`db_instance_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户数据库绑定表';

-- 迁移任务表
CREATE TABLE IF NOT EXISTS `tenant_isolation_migrations` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',
    `migration_type` ENUM('to_schema', 'to_database', 'schema_to_db') NOT NULL COMMENT '迁移类型',
    `source_strategy` ENUM('row_level', 'schema_level', 'database_level') NOT NULL COMMENT '源策略',
    `target_strategy` ENUM('row_level', 'schema_level', 'database_level') NOT NULL COMMENT '目标策略',
    `status` ENUM('pending', 'running', 'completed', 'failed', 'rolled_back') NOT NULL DEFAULT 'pending' COMMENT '状态',
    `progress` INT DEFAULT 0 COMMENT '进度（0-100）',
    `total_tables` INT DEFAULT 0 COMMENT '总表数',
    `migrated_tables` INT DEFAULT 0 COMMENT '已迁移表数',
    `total_rows` BIGINT DEFAULT 0 COMMENT '总行数',
    `migrated_rows` BIGINT DEFAULT 0 COMMENT '已迁移行数',
    `error_message` TEXT NULL COMMENT '错误信息',
    `started_at` TIMESTAMP NULL COMMENT '开始时间',
    `finished_at` TIMESTAMP NULL COMMENT '完成时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_migration_type` (`migration_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='隔离策略迁移任务表';

-- ================================================================================
-- 初始化数据
-- ================================================================================

-- 为现有租户初始化row_level隔离策略
INSERT INTO tenant_isolation_policies (tenant_id, isolation_strategy, status, migration_progress)
SELECT
    tenant_id,
    'row_level' as isolation_strategy,
    'active' as status,
    100 as migration_progress
FROM tenants
WHERE NOT EXISTS (
    SELECT 1 FROM tenant_isolation_policies tip WHERE tip.tenant_id = tenants.tenant_id
);

-- ================================================================================
-- 验证脚本
-- ================================================================================

-- 检查表是否创建成功
SELECT
    TABLE_NAME,
    TABLE_COMMENT,
    CREATE_TIME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN (
    'tenant_isolation_policies',
    'tenant_schema_mappings',
    'database_instances',
    'tenant_database_bindings',
    'tenant_isolation_migrations'
  )
ORDER BY TABLE_NAME;

-- 检查索引是否创建成功
SELECT
    TABLE_NAME,
    INDEX_NAME,
    GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) as columns
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME IN (
    'tenant_isolation_policies',
    'tenant_schema_mappings',
    'database_instances',
    'tenant_database_bindings',
    'tenant_isolation_migrations'
  )
GROUP BY TABLE_NAME, INDEX_NAME
ORDER BY TABLE_NAME, INDEX_NAME;
