-- ================================================================================
-- ZKER 企业级AI Agent平台 - 完整数据库DDL脚本
-- ================================================================================
-- 版本: v1.0
-- 生成日期: 2025-01-01
-- 表总数: 157张
-- 字符集: utf8mb4_unicode_ci
-- 引擎: InnoDB
--
-- 使用说明:
-- 1. 本脚本包含所有数据库表的完整定义
-- 2. 执行前请确保MySQL版本 >= 8.0
-- 3. 建议使用Atlas管理数据库版本
-- 4. 执行顺序: 按照模块顺序执行，确保依赖关系正确
-- ================================================================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 第一部分: 多租户核心表 (5张表) - P0优先级
-- ================================================================================

-- 1.1 租户表 (tenants)
CREATE TABLE IF NOT EXISTS tenants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '租户ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户唯一标识',
    name VARCHAR(255) NOT NULL COMMENT '租户名称',
    industry VARCHAR(50) NULL COMMENT '租户所属行业',
    company_size VARCHAR(20) NULL COMMENT '公司规模',
    region VARCHAR(50) NULL COMMENT '租户所在地区',
    timezone VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '租户时区',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '租户货币',
    logo_url VARCHAR(500) NULL COMMENT '租户Logo URL',
    custom_domain VARCHAR(255) NULL COMMENT '租户自定义域名',
    subscription_plan VARCHAR(50) DEFAULT 'free' COMMENT '订阅方案',
    status VARCHAR(20) DEFAULT 'active' COMMENT '租户状态',
    expires_at TIMESTAMP NULL COMMENT '租户过期时间',
    quota JSON COMMENT '配额配置',
    settings JSON COMMENT '租户设置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_tenant_id (tenant_id),
    UNIQUE KEY uk_custom_domain (custom_domain),
    INDEX idx_industry (industry),
    INDEX idx_region (region),
    INDEX idx_subscription_plan (subscription_plan),
    INDEX idx_expires_at (expires_at),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';

-- 1.2 租户订阅方案表 (tenant_plans)
CREATE TABLE IF NOT EXISTS tenant_plans (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    plan_id VARCHAR(64) NOT NULL COMMENT '方案ID',
    plan_name VARCHAR(100) NOT NULL COMMENT '方案名称',
    plan_type VARCHAR(50) NOT NULL COMMENT '方案类型',
    billing_cycle VARCHAR(20) NOT NULL COMMENT '计费周期',
    price DECIMAL(10,2) NOT NULL COMMENT '价格',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    features JSON COMMENT '功能特性',
    quotas JSON COMMENT '配额限制',
    start_date TIMESTAMP NOT NULL COMMENT '开始日期',
    end_date TIMESTAMP NULL COMMENT '结束日期',
    auto_renew BOOLEAN DEFAULT FALSE COMMENT '是否自动续费',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_plan_id (plan_id),
    INDEX idx_status (status),
    INDEX idx_end_date (end_date),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户订阅方案表';

-- 1.3 租户使用量表 (tenant_usage)
CREATE TABLE IF NOT EXISTS tenant_usage (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    metric_type VARCHAR(50) NOT NULL COMMENT '指标类型',
    metric_name VARCHAR(100) NOT NULL COMMENT '指标名称',
    usage_value BIGINT NOT NULL DEFAULT 0 COMMENT '使用量值',
    warning_threshold INT DEFAULT 80 COMMENT '预警阈值',
    limit_value BIGINT NULL COMMENT '使用量限制',
    usage_unit VARCHAR(20) DEFAULT 'count' COMMENT '使用量单位',
    reset_cycle VARCHAR(20) DEFAULT 'monthly' COMMENT '重置周期',
    next_reset_at TIMESTAMP NULL COMMENT '下次重置时间',
    usage_trend JSON COMMENT '使用量趋势',
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口',
    window_size INT NOT NULL COMMENT '窗口大小',
    dimensions JSON COMMENT '维度标签',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_tenant_metric_time (tenant_id, metric_type, time_window),
    INDEX idx_metric_type (metric_type),
    INDEX idx_next_reset_at (next_reset_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户使用量表';

-- 1.4 租户配额表 (tenant_quotas)
CREATE TABLE IF NOT EXISTS tenant_quotas (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    quota_type VARCHAR(50) NOT NULL COMMENT '配额类型',
    quota_name VARCHAR(100) NOT NULL COMMENT '配额名称',
    quota_limit BIGINT NOT NULL COMMENT '配额限制',
    quota_used BIGINT DEFAULT 0 COMMENT '已使用量',
    quota_remaining BIGINT GENERATED ALWAYS AS (quota_limit - quota_used) COMMENT '剩余量',
    reset_cycle VARCHAR(20) DEFAULT 'monthly' COMMENT '重置周期',
    last_reset_at TIMESTAMP NULL COMMENT '最后重置时间',
    next_reset_at TIMESTAMP NULL COMMENT '下次重置时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_tenant_quota_type (tenant_id, quota_type),
    INDEX idx_next_reset_at (next_reset_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户配额表';

-- 1.5 租户发票表 (tenant_invoices)
CREATE TABLE IF NOT EXISTS tenant_invoices (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '发票ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    invoice_no VARCHAR(64) NOT NULL COMMENT '发票号码',
    billing_period_start DATE NOT NULL COMMENT '账期开始',
    billing_period_end DATE NOT NULL COMMENT '账期结束',
    subtotal DECIMAL(12,2) NOT NULL COMMENT '小计',
    tax DECIMAL(12,2) DEFAULT 0.00 COMMENT '税额',
    total DECIMAL(12,2) NOT NULL COMMENT '总额',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    status VARCHAR(20) DEFAULT 'unpaid' COMMENT '发票状态',
    due_date TIMESTAMP NOT NULL COMMENT '到期日期',
    paid_at TIMESTAMP NULL COMMENT '支付时间',
    payment_method VARCHAR(50) NULL COMMENT '支付方式',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_invoice_no (invoice_no),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户发票表';

-- ================================================================================
-- 第二部分: 用户与权限管理表 (10张表)
-- ================================================================================

-- 2.1 用户表 (users)
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    user_id VARCHAR(64) NOT NULL COMMENT '用户唯一标识',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    username VARCHAR(100) NOT NULL COMMENT '用户名',
    email VARCHAR(255) NOT NULL COMMENT '邮箱',
    phone VARCHAR(20) NULL COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    avatar_url VARCHAR(500) NULL COMMENT '头像URL',
    status VARCHAR(20) DEFAULT 'active' COMMENT '用户状态',
    email_verified BOOLEAN DEFAULT FALSE COMMENT '邮箱是否验证',
    phone_verified BOOLEAN DEFAULT FALSE COMMENT '手机是否验证',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) NULL COMMENT '最后登录IP',
    settings JSON COMMENT '用户设置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_user_id (user_id),
    UNIQUE KEY uk_tenant_username (tenant_id, username),
    UNIQUE KEY uk_tenant_email (tenant_id, email),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2.2 用户资料表 (user_profiles)
CREATE TABLE IF NOT EXISTS user_profiles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    display_name VARCHAR(100) NULL COMMENT '显示名称',
    first_name VARCHAR(50) NULL COMMENT '名',
    last_name VARCHAR(50) NULL COMMENT '姓',
    gender VARCHAR(10) NULL COMMENT '性别',
    birthday DATE NULL COMMENT '生日',
    country VARCHAR(50) NULL COMMENT '国家',
    province VARCHAR(50) NULL COMMENT '省份',
    city VARCHAR(50) NULL COMMENT '城市',
    address TEXT NULL COMMENT '详细地址',
    bio TEXT NULL COMMENT '个人简介',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_user_id (user_id),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户资料表';

-- 2.3 角色表 (roles)
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '角色ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role_id VARCHAR(64) NOT NULL COMMENT '角色唯一标识',
    name VARCHAR(100) NOT NULL COMMENT '角色名称',
    code VARCHAR(100) NOT NULL COMMENT '角色编码',
    description TEXT NULL COMMENT '角色描述',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否系统角色',
    level INT DEFAULT 0 COMMENT '角色级别',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_role_id (role_id),
    UNIQUE KEY uk_tenant_code (tenant_id, code),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 2.4 权限表 (permissions)
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '权限ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    parent_id BIGINT NULL COMMENT '父权限ID',
    code VARCHAR(128) NOT NULL COMMENT '权限编码',
    name VARCHAR(128) NOT NULL COMMENT '权限名称',
    type VARCHAR(20) NOT NULL COMMENT '权限类型',
    resource_type VARCHAR(50) NULL COMMENT '资源类型',
    action VARCHAR(50) NULL COMMENT '操作',
    description TEXT NULL COMMENT '权限描述',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_tenant_code (tenant_id, code),
    INDEX idx_parent_id (parent_id),
    INDEX idx_resource_action (resource_type, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 2.5 角色权限关联表 (role_permissions)
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role_id VARCHAR(64) NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 2.6 用户角色关联表 (user_roles)
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    role_id VARCHAR(64) NOT NULL COMMENT '角色ID',
    granted_by VARCHAR(64) NULL COMMENT '授权人',
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_user_role (user_id, role_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_role_id (role_id),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 2.7 数据权限表 (data_permissions)
CREATE TABLE IF NOT EXISTS data_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role_id VARCHAR(64) NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    permission_scope VARCHAR(20) NOT NULL COMMENT '权限范围',
    custom_org_ids JSON NULL COMMENT '自定义组织ID列表',
    custom_user_ids JSON NULL COMMENT '自定义用户ID列表',
    permissions JSON NOT NULL COMMENT '权限操作',
    inherits_from_id BIGINT NULL COMMENT '继承自哪个权限',
    priority INT DEFAULT 0 COMMENT '权限优先级',
    conditions JSON NULL COMMENT '权限生效条件',
    cache_key VARCHAR(255) NULL COMMENT '缓存键',
    permission_status VARCHAR(20) DEFAULT 'active' COMMENT '权限状态',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_role_id (role_id),
    INDEX idx_resource_type (resource_type),
    INDEX idx_priority (priority),
    INDEX idx_cache_key (cache_key),
    INDEX idx_permission_status (permission_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据权限表';

-- 2.8 字段权限表 (field_permissions)
CREATE TABLE IF NOT EXISTS field_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    role_id VARCHAR(64) NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    field_name VARCHAR(100) NOT NULL COMMENT '字段名',
    field_display_name VARCHAR(100) NULL COMMENT '字段显示名称',
    permission ENUM('visible', 'masked', 'hidden') NOT NULL COMMENT '权限级别',
    mask_rule VARCHAR(50) NULL COMMENT '脱敏规则',
    mask_pattern VARCHAR(255) NULL COMMENT '脱敏正则',
    mask_replacement VARCHAR(10) DEFAULT '*' COMMENT '替换字符',
    permission_source VARCHAR(50) DEFAULT 'manual' COMMENT '权限来源',
    field_group_name VARCHAR(100) NULL COMMENT '字段组名称',
    dependent_fields JSON NULL COMMENT '依赖字段',
    validation_rules JSON NULL COMMENT '验证规则',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_role_id (role_id),
    INDEX idx_resource_field (resource_type, field_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字段权限表';

-- 2.9 临时授权表 (temporary_grants)
CREATE TABLE IF NOT EXISTS temporary_grants (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    grantor_id BIGINT NOT NULL COMMENT '授权人ID',
    grantee_id BIGINT NOT NULL COMMENT '被授权人ID',
    approver_id BIGINT NULL COMMENT '审批人ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    resource_id BIGINT NOT NULL COMMENT '资源ID',
    permission_config JSON NOT NULL COMMENT '权限配置',
    start_time TIMESTAMP NOT NULL COMMENT '开始时间',
    end_time TIMESTAMP NOT NULL COMMENT '结束时间',
    approval_status VARCHAR(20) DEFAULT 'pending' COMMENT '审批状态',
    request_reason TEXT NULL COMMENT '申请理由',
    approval_reason TEXT NULL COMMENT '审批理由',
    max_usage_count INT NULL COMMENT '最大使用次数',
    usage_count INT DEFAULT 0 COMMENT '已使用次数',
    auto_renew BOOLEAN DEFAULT FALSE COMMENT '是否自动续期',
    notify_grantee BOOLEAN DEFAULT TRUE COMMENT '是否通知被授权人',
    notify_grantor BOOLEAN DEFAULT TRUE COMMENT '是否通知授权人',
    is_revoked BOOLEAN DEFAULT FALSE COMMENT '是否已撤销',
    revoked_at TIMESTAMP NULL COMMENT '撤销时间',
    revoked_by BIGINT NULL COMMENT '撤销人ID',
    revoke_reason TEXT NULL COMMENT '撤销理由',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_grantor_id (grantor_id),
    INDEX idx_grantee_id (grantee_id),
    INDEX idx_approval_status (approval_status),
    INDEX idx_start_end_time (start_time, end_time),
    INDEX idx_is_revoked (is_revoked)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='临时授权表';

-- 2.10 登录日志表 (login_logs)
CREATE TABLE IF NOT EXISTS login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    username VARCHAR(100) NOT NULL COMMENT '用户名',
    login_type VARCHAR(20) NOT NULL COMMENT '登录类型',
    login_status VARCHAR(20) NOT NULL COMMENT '登录状态',
    failure_reason VARCHAR(255) NULL COMMENT '失败原因',
    ip_address VARCHAR(50) NOT NULL COMMENT 'IP地址',
    user_agent TEXT NULL COMMENT 'User-Agent',
    device_type VARCHAR(50) NULL COMMENT '设备类型',
    browser_type VARCHAR(50) NULL COMMENT '浏览器类型',
    location_country VARCHAR(50) NULL COMMENT '国家',
    location_province VARCHAR(50) NULL COMMENT '省份',
    location_city VARCHAR(50) NULL COMMENT '城市',
    login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_login_status (login_status),
    INDEX idx_login_at (login_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ================================================================================
-- 第三部分: 组织管理表 (12张表) - P1优先级
-- ================================================================================

-- 3.1 组织表 (organizations)
CREATE TABLE IF NOT EXISTS organizations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '组织ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    organization_id VARCHAR(64) NOT NULL COMMENT '组织唯一标识',
    parent_id BIGINT NULL COMMENT '父组织ID',
    name VARCHAR(255) NOT NULL COMMENT '组织名称',
    code VARCHAR(100) NOT NULL COMMENT '组织编码',
    type VARCHAR(50) DEFAULT 'department' COMMENT '组织类型',
    path VARCHAR(500) NULL COMMENT '组织路径',
    level INT DEFAULT 0 COMMENT '组织层级',
    sort_order INT DEFAULT 0 COMMENT '排序',
    leader_id BIGINT NULL COMMENT '负责人ID',
    description TEXT NULL COMMENT '组织描述',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_organization_id (organization_id),
    UNIQUE KEY uk_tenant_code (tenant_id, code),
    INDEX idx_parent_id (parent_id),
    INDEX idx_path (path(255)),
    INDEX idx_level (level),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';

-- 3.2 岗位表 (positions)
CREATE TABLE IF NOT EXISTS positions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '岗位ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    organization_id BIGINT NOT NULL COMMENT '所属组织ID',
    position_id VARCHAR(64) NOT NULL COMMENT '岗位唯一标识',
    name VARCHAR(100) NOT NULL COMMENT '岗位名称',
    code VARCHAR(100) NOT NULL COMMENT '岗位编码',
    level VARCHAR(50) NULL COMMENT '岗位级别',
    category VARCHAR(50) NULL COMMENT '岗位类别',
    description TEXT NULL COMMENT '岗位描述',
    requirements TEXT NULL COMMENT '岗位要求',
    responsibilities TEXT NULL COMMENT '岗位职责',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_position_id (position_id),
    UNIQUE KEY uk_org_code (organization_id, code),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='岗位表';

-- 3.3 成员表 (members)
CREATE TABLE IF NOT EXISTS members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '成员ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    organization_id BIGINT NOT NULL COMMENT '所属组织ID',
    member_id VARCHAR(64) NOT NULL COMMENT '成员唯一标识',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    position_id BIGINT NULL COMMENT '岗位ID',
    member_type VARCHAR(20) DEFAULT 'employee' COMMENT '成员类型',
    employee_no VARCHAR(50) NULL COMMENT '工号',
    email VARCHAR(255) NULL COMMENT '邮箱',
    phone VARCHAR(20) NULL COMMENT '手机号',
    work_experiences JSON NULL COMMENT '工作经历',
    educations JSON NULL COMMENT '教育背景',
    skills JSON NULL COMMENT '技能标签',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    joined_at TIMESTAMP NULL COMMENT '加入时间',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_member_id (member_id),
    UNIQUE KEY uk_user_org (user_id, organization_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_position_id (position_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成员表';

-- 3.4 成员角色关联表 (member_roles)
CREATE TABLE IF NOT EXISTS member_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    member_id VARCHAR(64) NOT NULL COMMENT '成员ID',
    role_id VARCHAR(64) NOT NULL COMMENT '角色ID',
    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否主角色',
    granted_by BIGINT NULL COMMENT '授权人',
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_member_role (member_id, role_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_role_id (role_id),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成员角色关联表';

-- 3.5 成员技能表 (member_skills)
CREATE TABLE IF NOT EXISTS member_skills (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    member_id VARCHAR(64) NOT NULL COMMENT '成员ID',
    skill_name VARCHAR(100) NOT NULL COMMENT '技能名称',
    skill_level VARCHAR(20) DEFAULT 'intermediate' COMMENT '技能水平',
    certified BOOLEAN DEFAULT FALSE COMMENT '是否认证',
    years_experience INT DEFAULT 0 COMMENT '年限',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_member_id (member_id),
    INDEX idx_skill_level (skill_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成员技能表';

-- 3.6 工作经历表 (work_experiences)
CREATE TABLE IF NOT EXISTS work_experiences (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    member_id VARCHAR(64) NOT NULL COMMENT '成员ID',
    company_name VARCHAR(255) NOT NULL COMMENT '公司名称',
    position VARCHAR(100) NOT NULL COMMENT '职位',
    industry VARCHAR(50) NULL COMMENT '行业',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NULL COMMENT '结束日期',
    is_current BOOLEAN DEFAULT FALSE COMMENT '是否当前',
    description TEXT NULL COMMENT '工作描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_member_id (member_id),
    INDEX idx_start_date (start_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作经历表';

-- 3.7 教育背景表 (educations)
CREATE TABLE IF NOT EXISTS educations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    member_id VARCHAR(64) NOT NULL COMMENT '成员ID',
    school_name VARCHAR(255) NOT NULL COMMENT '学校名称',
    major VARCHAR(100) NOT NULL COMMENT '专业',
    degree VARCHAR(50) NOT NULL COMMENT '学位',
    education_level VARCHAR(20) NOT NULL COMMENT '学历',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NULL COMMENT '结束日期',
    description TEXT NULL COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_member_id (member_id),
    INDEX idx_education_level (education_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教育背景表';

-- 3.8 岗位变更表 (position_changes)
CREATE TABLE IF NOT EXISTS position_changes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    member_id VARCHAR(64) NOT NULL COMMENT '成员ID',
    from_position_id BIGINT NULL COMMENT '原岗位',
    to_position_id BIGINT NOT NULL COMMENT '新岗位',
    from_organization_id BIGINT NULL COMMENT '原组织',
    to_organization_id BIGINT NOT NULL COMMENT '新组织',
    change_type VARCHAR(20) NOT NULL COMMENT '变更类型',
    change_reason TEXT NULL COMMENT '变更原因',
    effective_date DATE NOT NULL COMMENT '生效日期',
    approved_by BIGINT NULL COMMENT '审批人',
    approved_at TIMESTAMP NULL COMMENT '审批时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_member_id (member_id),
    INDEX idx_effective_date (effective_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='岗位变更表';

-- 3.9 虚拟组织表 (virtual_organizations)
CREATE TABLE IF NOT EXISTS virtual_organizations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    org_id VARCHAR(64) NOT NULL COMMENT '虚拟组织ID',
    name VARCHAR(255) NOT NULL COMMENT '组织名称',
    type VARCHAR(50) NOT NULL COMMENT '组织类型',
    description TEXT NULL COMMENT '组织描述',
    owner_id BIGINT NOT NULL COMMENT '负责人',
    members JSON NULL COMMENT '成员ID列表',
    permissions JSON NULL COMMENT '组织权限',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_org_id (org_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_owner_id (owner_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='虚拟组织表';

-- 3.10 虚拟组织成员表 (virtual_org_members)
CREATE TABLE IF NOT EXISTS virtual_org_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    virtual_org_id VARCHAR(64) NOT NULL COMMENT '虚拟组织ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role VARCHAR(50) DEFAULT 'member' COMMENT '角色',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    invited_by BIGINT NULL COMMENT '邀请人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_org_user (virtual_org_id, user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='虚拟组织成员表';

-- 3.11 组织委托表 (organization_delegates)
CREATE TABLE IF NOT EXISTS organization_delegates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    delegator_id BIGINT NOT NULL COMMENT '委托人ID',
    delegate_id BIGINT NOT NULL COMMENT '被委托人ID',
    organization_id BIGINT NOT NULL COMMENT '组织ID',
    delegate_permissions JSON NOT NULL COMMENT '委托权限',
    start_time TIMESTAMP NOT NULL COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否生效',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_delegator_id (delegator_id),
    INDEX idx_delegate_id (delegate_id),
    INDEX idx_org_id (organization_id),
    INDEX idx_start_end_time (start_time, end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织委托表';

-- 3.12 团队成员表 (team_members)
CREATE TABLE IF NOT EXISTS team_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '成员ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    team_id BIGINT NOT NULL COMMENT '团队ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role VARCHAR(50) DEFAULT 'member' COMMENT '角色',
    permissions JSON NULL COMMENT '权限',
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_team_user (team_id, user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='团队成员表';

-- 继续下一部分...
-- 由于脚本过长，将分为多个文件

SET FOREIGN_KEY_CHECKS = 1;
