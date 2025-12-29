-- ============================================================================
-- ZKER Enterprise - 数据库完整SQL脚本
-- ============================================================================
-- 版本: v1.0.0
-- 创建日期: 2025-12-29
-- 数据库: MySQL 8.4.5+
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci
--
-- 说明:
-- 1. 本脚本包含所有核心表的完整DDL定义
-- 2. 所有表都包含 tenant_id 字段用于多租户隔离
-- 3. 所有表都包含标准的审计字段 (created_at, updated_at, deleted_at)
-- 4. 所有表都使用 InnoDB 引擎
-- 5. 主键使用 BIGINT 自增或 VARCHAR(64) UUID
-- ============================================================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================================
-- 1. 租户管理 (Tenants)
-- ============================================================================

DROP TABLE IF EXISTS `tenants`;
CREATE TABLE `tenants` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '租户ID',
    `name` VARCHAR(255) NOT NULL COMMENT '租户名称',
    `domain` VARCHAR(100) NOT NULL UNIQUE COMMENT '租户域名',
    `plan` ENUM('free', 'basic', 'professional', 'enterprise') NOT NULL DEFAULT 'free' COMMENT '订阅计划',
    `status` ENUM('active', 'suspended', 'expired') NOT NULL DEFAULT 'active' COMMENT '租户状态',

    -- 配额限制
    `max_users` INT NOT NULL DEFAULT 100 COMMENT '最大用户数',
    `max_bots` INT NOT NULL DEFAULT 10 COMMENT '最大Bot数',
    `max_storage_gb` INT NOT NULL DEFAULT 10 COMMENT '最大存储空间(GB)',

    -- 配置
    `settings` JSON COMMENT '租户配置',
    `logo_url` VARCHAR(512) COMMENT 'Logo URL',
    `custom_domain` VARCHAR(100) COMMENT '自定义域名',
    `sso_enabled` BOOLEAN DEFAULT FALSE COMMENT '是否启用SSO',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_domain` (`domain`),
    INDEX `idx_status` (`status`),
    INDEX `idx_plan` (`plan`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- ============================================================================
-- 2. 用户管理 (Users)
-- ============================================================================

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 基本信息
    `username` VARCHAR(100) NOT NULL COMMENT '用户名',
    `email` VARCHAR(255) NOT NULL COMMENT '邮箱',
    `phone` VARCHAR(20) COMMENT '手机号',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
    `avatar_url` VARCHAR(512) COMMENT '头像URL',

    -- 角色权限
    `role` ENUM('super_admin', 'admin', 'user', 'viewer') NOT NULL DEFAULT 'user' COMMENT '角色',
    `status` ENUM('active', 'inactive', 'suspended') NOT NULL DEFAULT 'active' COMMENT '用户状态',

    -- 个人资料
    `nickname` VARCHAR(100) COMMENT '昵称',
    `department_id` BIGINT COMMENT '部门ID',
    `job_title` VARCHAR(100) COMMENT '职位',

    -- 登录信息
    `last_login_at` DATETIME COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(64) COMMENT '最后登录IP',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY `uk_tenant_username` (`tenant_id`, `username`),
    UNIQUE KEY `uk_tenant_email` (`tenant_id`, `email`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_email` (`email`),
    INDEX `idx_status` (`status`),
    INDEX `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================================================
-- 3. 组织架构 (Departments)
-- ============================================================================

DROP TABLE IF EXISTS `departments`;
CREATE TABLE `departments` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '部门ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `parent_id` BIGINT COMMENT '父部门ID',

    -- 基本信息
    `name` VARCHAR(255) NOT NULL COMMENT '部门名称',
    `code` VARCHAR(100) COMMENT '部门编码',
    `description` TEXT COMMENT '部门描述',
    `level` INT NOT NULL DEFAULT 1 COMMENT '部门层级',
    `path` VARCHAR(500) COMMENT '部门路径',
    `sort_order` INT DEFAULT 0 COMMENT '排序',

    -- 统计信息
    `member_count` INT DEFAULT 0 COMMENT '成员数量',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_parent_id` (`parent_id`),
    INDEX `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';

-- ============================================================================
-- 4. Bot管理 (Bots)
-- ============================================================================

DROP TABLE IF EXISTS `bots`;
CREATE TABLE `bots` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT 'Bot ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `creator_id` BIGINT NOT NULL COMMENT '创建者ID',

    -- 基本信息
    `name` VARCHAR(255) NOT NULL COMMENT 'Bot名称',
    `description` TEXT COMMENT 'Bot描述',
    `avatar_url` VARCHAR(512) COMMENT 'Bot头像',
    `status` ENUM('active', 'inactive', 'archived') NOT NULL DEFAULT 'active' COMMENT 'Bot状态',
    `type` ENUM('chat', 'workflow', 'agent') NOT NULL DEFAULT 'chat' COMMENT 'Bot类型',

    -- 配置信息
    `prompt_template` TEXT COMMENT '提示词模板',
    `model_config` JSON COMMENT '模型配置',
    `welcome_message` VARCHAR(500) COMMENT '欢迎语',

    -- 统计信息
    `total_conversations` INT DEFAULT 0 COMMENT '总对话数',
    `total_messages` INT DEFAULT 0 COMMENT '总消息数',
    `active_users` INT DEFAULT 0 COMMENT '活跃用户数',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_creator_id` (`creator_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_type` (`type`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot表';

-- ============================================================================
-- 5. 对话管理 (Conversations)
-- ============================================================================

DROP TABLE IF EXISTS `conversations`;
CREATE TABLE `conversations` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '对话ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `bot_id` VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',

    -- 状态信息
    `status` ENUM('active', 'closed', 'archived') NOT NULL DEFAULT 'active' COMMENT '对话状态',
    `title` VARCHAR(255) COMMENT '对话标题',

    -- 统计信息
    `message_count` INT DEFAULT 0 COMMENT '消息数量',
    `token_usage` INT DEFAULT 0 COMMENT 'Token使用量',

    -- 元数据
    `metadata` JSON COMMENT '元数据',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_bot_id` (`bot_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话表';

-- ============================================================================
-- 6. 消息管理 (Messages)
-- ============================================================================

DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '消息ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `conversation_id` VARCHAR(64) NOT NULL COMMENT '对话ID',

    -- 消息内容
    `role` ENUM('user', 'assistant', 'system') NOT NULL COMMENT '角色',
    `content` TEXT NOT NULL COMMENT '消息内容',
    `content_type` ENUM('text', 'image', 'audio', 'video', 'file') NOT NULL DEFAULT 'text' COMMENT '内容类型',

    -- Token使用
    `input_tokens` INT DEFAULT 0 COMMENT '输入Token数',
    `output_tokens` INT DEFAULT 0 COMMENT '输出Token数',
    `total_tokens` INT DEFAULT 0 COMMENT '总Token数',

    -- 性能指标
    `response_time` INT COMMENT '响应时间(ms)',
    `model` VARCHAR(100) COMMENT '使用的模型',

    -- 反馈
    `feedback_score` INT COMMENT '反馈评分(1-5)',
    `feedback_tag` VARCHAR(50) COMMENT '反馈标签',

    -- 元数据
    `metadata` JSON COMMENT '元数据',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_conversation_id` (`conversation_id`),
    INDEX `idx_role` (`role`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

-- ============================================================================
-- 7. 知识库管理 (Knowledge Bases)
-- ============================================================================

DROP TABLE IF EXISTS `knowledge_bases`;
CREATE TABLE `knowledge_bases` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '知识库ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `creator_id` BIGINT NOT NULL COMMENT '创建者ID',

    -- 基本信息
    `name` VARCHAR(255) NOT NULL COMMENT '知识库名称',
    `description` TEXT COMMENT '知识库描述',
    `type` ENUM('document', 'qa', 'dataset') NOT NULL DEFAULT 'document' COMMENT '知识库类型',

    -- 配置信息
    `chunking_strategy` ENUM('fixed', 'semantic', 'recursive', 'hybrid') NOT NULL DEFAULT 'hybrid' COMMENT '分块策略',
    `chunk_size` INT DEFAULT 500 COMMENT '分块大小',
    `chunk_overlap` INT DEFAULT 50 COMMENT '分块重叠',
    `embedding_model` VARCHAR(100) NOT NULL DEFAULT 'text-embedding-ada-002' COMMENT '向量化模型',

    -- 统计信息
    `document_count` INT DEFAULT 0 COMMENT '文档数量',
    `chunk_count` INT DEFAULT 0 COMMENT '分块数量',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_creator_id` (`creator_id`),
    INDEX `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识库表';

-- ============================================================================
-- 8. 知识库文档 (Knowledge Documents)
-- ============================================================================

DROP TABLE IF EXISTS `knowledge_documents`;
CREATE TABLE `knowledge_documents` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '文档ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `knowledge_base_id` VARCHAR(64) NOT NULL COMMENT '知识库ID',

    -- 文档信息
    `filename` VARCHAR(255) NOT NULL COMMENT '文件名',
    `file_size` BIGINT COMMENT '文件大小(字节)',
    `file_type` VARCHAR(50) COMMENT '文件类型',
    `file_url` VARCHAR(512) COMMENT '文件URL',

    -- 处理状态
    `status` ENUM('pending', 'processing', 'completed', 'failed') NOT NULL DEFAULT 'pending' COMMENT '处理状态',
    `chunk_count` INT DEFAULT 0 COMMENT '分块数量',
    `error_message` TEXT COMMENT '错误信息',

    -- 元数据
    `metadata` JSON COMMENT '文档元数据',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_knowledge_base_id` (`knowledge_base_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识库文档表';

-- ============================================================================
-- 9. 知识库分块 (Knowledge Chunks) - RAG优化
-- ============================================================================

DROP TABLE IF EXISTS `knowledge_chunks`;
CREATE TABLE `knowledge_chunks` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分块ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `knowledge_base_id` VARCHAR(64) NOT NULL COMMENT '知识库ID',
    `document_id` VARCHAR(64) NOT NULL COMMENT '文档ID',

    -- 分块内容
    `chunk_text` TEXT NOT NULL COMMENT '分块文本',
    `chunk_type` ENUM('fixed', 'semantic', 'recursive', 'hybrid') NOT NULL COMMENT '分块类型',
    `chunk_index` INT NOT NULL COMMENT '分块序号',

    -- 向量数据
    `embedding_vector` VECTOR(1536) NOT NULL COMMENT '向量嵌入(1536维)',

    -- 元数据
    `metadata` JSON COMMENT '元数据',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_knowledge_base_id` (`knowledge_base_id`),
    INDEX `idx_document_id` (`document_id`),
    INDEX `idx_embedding_vector` (`embedding_vector`),
    FULLTEXT INDEX `ft_chunk_text` (`chunk_text`) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识库分块表';

-- ============================================================================
-- 10. 工作流管理 (Workflows)
-- ============================================================================

DROP TABLE IF EXISTS `workflows`;
CREATE TABLE `workflows` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '工作流ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `creator_id` BIGINT NOT NULL COMMENT '创建者ID',

    -- 基本信息
    `name` VARCHAR(255) NOT NULL COMMENT '工作流名称',
    `description` TEXT COMMENT '工作流描述',
    `status` ENUM('draft', 'active', 'archived') NOT NULL DEFAULT 'draft' COMMENT '工作流状态',

    -- 工作流定义
    `workflow_definition` JSON NOT NULL COMMENT '工作流定义(DAG)',
    `version` INT DEFAULT 1 COMMENT '版本号',

    -- 统计信息
    `total_executions` INT DEFAULT 0 COMMENT '总执行次数',
    `success_executions` INT DEFAULT 0 COMMENT '成功执行次数',
    `failed_executions` INT DEFAULT 0 COMMENT '失败执行次数',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_creator_id` (`creator_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流表';

-- ============================================================================
-- 11. 插件管理 (Plugins)
-- ============================================================================

DROP TABLE IF EXISTS `plugins`;
CREATE TABLE `plugins` (
    `id` VARCHAR(64) PRIMARY KEY COMMENT '插件ID',

    -- 基本信息
    `name` VARCHAR(255) NOT NULL COMMENT '插件名称',
    `description` TEXT COMMENT '插件描述',
    `category` VARCHAR(50) NOT NULL COMMENT '插件分类',
    `version` VARCHAR(20) NOT NULL COMMENT '插件版本',
    `author` VARCHAR(100) COMMENT '作者',

    -- 图标和媒体
    `icon_url` VARCHAR(512) COMMENT '插件图标',
    `screenshot_urls` JSON COMMENT '截图URL列表',

    -- 配置schema
    `config_schema` JSON COMMENT '配置Schema',
    `api_endpoint` VARCHAR(512) COMMENT 'API端点',

    -- 统计信息
    `install_count` INT DEFAULT 0 COMMENT '安装次数',
    `rating` DECIMAL(3,2) DEFAULT 0.00 COMMENT '评分(0-5)',
    `review_count` INT DEFAULT 0 COMMENT '评论数',

    -- 状态
    `status` ENUM('pending', 'active', 'rejected') NOT NULL DEFAULT 'pending' COMMENT '审核状态',
    `is_official` BOOLEAN DEFAULT FALSE COMMENT '是否官方插件',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    INDEX `idx_category` (`category`),
    INDEX `idx_status` (`status`),
    INDEX `idx_rating` (`rating`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='插件表';

-- ============================================================================
-- 12. Bot插件关联 (Bot Plugins)
-- ============================================================================

DROP TABLE IF EXISTS `bot_plugins`;
CREATE TABLE `bot_plugins` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `bot_id` VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    `plugin_id` VARCHAR(64) NOT NULL COMMENT '插件ID',

    -- 配置信息
    `config` JSON COMMENT '插件配置',
    `enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY `uk_bot_plugin` (`bot_id`, `plugin_id`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_bot_id` (`bot_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot插件关联表';

-- ============================================================================
-- 13. Bot指标数据 (Bot Metrics) - 数据分析
-- ============================================================================

DROP TABLE IF EXISTS `bot_metrics`;
CREATE TABLE `bot_metrics` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `bot_id` VARCHAR(64) NOT NULL COMMENT 'Bot ID',

    -- 时间维度
    `metric_date` DATE NOT NULL COMMENT '指标日期',
    `metric_hour` TINYINT COMMENT '指标小时(0-23)',

    -- 对话指标
    `conversation_count` INT DEFAULT 0 COMMENT '对话次数',
    `message_count` INT DEFAULT 0 COMMENT '消息数量',
    `active_users` INT DEFAULT 0 COMMENT '活跃用户数',
    `new_users` INT DEFAULT 0 COMMENT '新增用户数',

    -- 性能指标
    `avg_response_time` INT COMMENT '平均响应时间(ms)',
    `p95_response_time` INT COMMENT 'P95响应时间(ms)',
    `p99_response_time` INT COMMENT 'P99响应时间(ms)',
    `error_rate` DECIMAL(5,4) COMMENT '错误率',

    -- 消耗指标
    `token_consumption` BIGINT DEFAULT 0 COMMENT 'Token消耗量',
    `api_call_count` INT DEFAULT 0 COMMENT 'API调用次数',

    -- 满意度指标
    `satisfaction_score` DECIMAL(3,2) COMMENT '满意度评分(0-5)',

    INDEX `idx_tenant_bot_date` (`tenant_id`, `bot_id`, `metric_date`),
    INDEX `idx_metric_date` (`metric_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot指标数据表';

-- ============================================================================
-- 14. Token使用计量 (Token Usage Metrics) - Token计量补充
-- ============================================================================

DROP TABLE IF EXISTS `token_usage_metrics`;
CREATE TABLE `token_usage_metrics` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `bot_id` VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    `user_id` BIGINT COMMENT '用户ID',
    `department_id` BIGINT COMMENT '部门ID',

    -- 时间维度
    `usage_date` DATE NOT NULL COMMENT '使用日期',

    -- Token使用量
    `input_tokens` BIGINT DEFAULT 0 COMMENT '输入Token数',
    `output_tokens` BIGINT DEFAULT 0 COMMENT '输出Token数',
    `total_tokens` BIGINT DEFAULT 0 COMMENT '总Token数',

    -- 成本信息
    `model` VARCHAR(100) NOT NULL COMMENT '模型名称',
    `unit_cost` DECIMAL(10,6) COMMENT '单价(Token/美元)',
    `total_cost` DECIMAL(10,4) COMMENT '总成本(美元)',

    INDEX `idx_tenant_date` (`tenant_id`, `usage_date`),
    INDEX `idx_bot_date` (`bot_id`, `usage_date`),
    INDEX `idx_user_date` (`user_id`, `usage_date`),
    INDEX `idx_department_date` (`department_id`, `usage_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Token使用计量表';

-- ============================================================================
-- 15. 预算告警 (Budget Alerts) - Token计量补充
-- ============================================================================

DROP TABLE IF EXISTS `budget_alerts`;
CREATE TABLE `budget_alerts` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `bot_id` VARCHAR(64) COMMENT 'Bot ID(为空表示租户级)',
    `user_id` BIGINT COMMENT '用户ID(为空表示部门级)',
    `department_id` BIGINT COMMENT '部门ID',

    -- 预算配置
    `alert_type` ENUM('token', 'cost') NOT NULL COMMENT '告警类型',
    `budget_threshold` BIGINT NOT NULL COMMENT '预算阈值',
    `alert_threshold` INT NOT NULL COMMENT '告警阈值(百分比)',

    -- 告警状态
    `current_usage` BIGINT DEFAULT 0 COMMENT '当前使用量',
    `alert_triggered` BOOLEAN DEFAULT FALSE COMMENT '是否已触发告警',
    `last_alert_at` DATETIME COMMENT '最后告警时间',

    -- 自动降级配置
    `downgrade_enabled` BOOLEAN DEFAULT FALSE COMMENT '是否启用自动降级',
    `downgrade_model` VARCHAR(100) COMMENT '降级后的模型',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_bot_id` (`bot_id`),
    INDEX `idx_alert_triggered` (`alert_triggered`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预算告警表';

-- ============================================================================
-- 16. 权限管理 (Permissions) - RBAC细化补充
-- ============================================================================

DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 基本信息
    `name` VARCHAR(100) NOT NULL COMMENT '角色名称',
    `code` VARCHAR(100) NOT NULL COMMENT '角色编码',
    `description` TEXT COMMENT '角色描述',
    `level` INT NOT NULL DEFAULT 1 COMMENT '角色级别(1-5,数字越小权限越大)',

    -- 是否系统角色
    `is_system` BOOLEAN DEFAULT FALSE COMMENT '是否系统角色',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
    INDEX `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `role_id` BIGINT NOT NULL COMMENT '角色ID',

    -- 权限信息
    `resource` VARCHAR(100) NOT NULL COMMENT '资源(bot/workflow/knowledge等)',
    `action` VARCHAR(50) NOT NULL COMMENT '操作(view/edit/delete/publish等)',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY `uk_role_resource_action` (`role_id`, `resource`, `action`),
    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

DROP TABLE IF EXISTS `data_permissions`;
CREATE TABLE `data_permissions` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `role_id` BIGINT COMMENT '角色ID(对角色生效则user_id为空)',

    -- 数据权限范围
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `permission_level` ENUM('all', 'department', 'department_and_sub', 'self', 'custom') NOT NULL COMMENT '权限级别',
    `custom_filter` JSON COMMENT '自定义过滤条件',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据权限表';

DROP TABLE IF EXISTS `temporary_grants`;
CREATE TABLE `temporary_grants` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(64) NOT NULL COMMENT '资源ID',

    -- 授权信息
    `permission` VARCHAR(50) NOT NULL COMMENT '权限',
    `granted_by` BIGINT NOT NULL COMMENT '授权人ID',
    `reason` VARCHAR(500) COMMENT '授权原因',

    -- 时间限制
    `granted_at` DATETIME NOT NULL COMMENT '授权时间',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',

    -- 状态
    `status` ENUM('active', 'expired', 'revoked') NOT NULL DEFAULT 'active' COMMENT '状态',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_expires_at` (`expires_at`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='临时授权表';

-- ============================================================================
-- 17. 监控日志 (Audit Logs)
-- ============================================================================

DROP TABLE IF EXISTS `audit_logs`;
CREATE TABLE `audit_logs` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',

    -- 操作信息
    `resource_type` VARCHAR(50) NOT NULL COMMENT '资源类型',
    `resource_id` VARCHAR(64) NOT NULL COMMENT '资源ID',
    `action` VARCHAR(50) NOT NULL COMMENT '操作类型',
    `action_detail` VARCHAR(100) COMMENT '操作详情',

    -- 请求信息
    `request_method` VARCHAR(10) COMMENT '请求方法',
    `request_url` VARCHAR(512) COMMENT '请求URL',
    `request_ip` VARCHAR(64) COMMENT '请求IP',
    `user_agent` VARCHAR(512) COMMENT 'User-Agent',

    -- 响应信息
    `response_status` INT COMMENT '响应状态码',
    `response_time` INT COMMENT '响应时间(ms)',
    `error_message` TEXT COMMENT '错误信息',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_resource` (`resource_type`, `resource_id`),
    INDEX `idx_action` (`action`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志表';

-- ============================================================================
-- 18. 订阅管理 (Subscriptions)
-- ============================================================================

DROP TABLE IF EXISTS `subscriptions`;
CREATE TABLE `subscriptions` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'ID',
    `tenant_id` VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 订阅信息
    `plan` ENUM('free', 'basic', 'professional', 'enterprise') NOT NULL COMMENT '订阅计划',
    `status` ENUM('active', 'expired', 'cancelled', 'suspended') NOT NULL DEFAULT 'active' COMMENT '订阅状态',

    -- 计费周期
    `billing_cycle` ENUM('monthly', 'quarterly', 'yearly') NOT NULL COMMENT '计费周期',
    `price` DECIMAL(10,2) NOT NULL COMMENT '价格',

    -- 时间范围
    `start_date` DATE NOT NULL COMMENT '开始日期',
    `end_date` DATE NOT NULL COMMENT '结束日期',

    -- 自动续费
    `auto_renew` BOOLEAN DEFAULT TRUE COMMENT '是否自动续费',

    -- 审计字段
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX `idx_tenant_id` (`tenant_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_end_date` (`end_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅表';

-- ============================================================================
-- 初始化数据
-- ============================================================================

-- 插入默认租户
INSERT INTO `tenants` (`id`, `name`, `domain`, `plan`, `status`, `max_users`, `max_bots`) VALUES
('tenant_default', '默认租户', 'default', 'free', 'active', 100, 10);

-- 插入超级管理员用户
INSERT INTO `users` (`id`, `tenant_id`, `username`, `email`, `password_hash`, `role`, `status`) VALUES
(1, 'tenant_default', 'admin', 'admin@zker.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'super_admin', 'active');

-- ============================================================================
-- 恢复外键检查
-- ============================================================================
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 说明
-- ============================================================================
-- 1. 密码哈希示例: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy (对应密码: admin123)
-- 2. 所有表都包含tenant_id用于多租户隔离
-- 3. 所有表都包含标准审计字段
-- 4. 向量字段需要MySQL 8.0.17+或使用插件支持
-- 5. 建议在生产环境中启用binlog和慢查询日志
-- 6. 定期执行OPTIMIZE TABLE优化表性能
-- ============================================================================
