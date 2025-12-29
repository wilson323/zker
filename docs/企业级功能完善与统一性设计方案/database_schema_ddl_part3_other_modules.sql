-- ================================================================================
-- ZKER 数据库DDL脚本 - 第三部分: 其他核心模块 (121张表)
-- ================================================================================
-- 版本: v1.0
-- 生成日期: 2025-01-01
-- 包含模块: 会话管理、Bot管理、知识管理、工作流、商业化等
-- ================================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 会话管理模块 (7张表)
-- ================================================================================

-- 3.1 会话表 (conversations)
CREATE TABLE IF NOT EXISTS conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会话ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话唯一标识',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    bot_id VARCHAR(64) NULL COMMENT 'Bot ID',
    title VARCHAR(500) NULL COMMENT '会话标题',
    summary TEXT NULL COMMENT '会话摘要',
    keywords JSON NULL COMMENT '会话关键词',
    sentiment_score DECIMAL(5,4) NULL COMMENT '情感评分',
    topic_category VARCHAR(50) NULL COMMENT '主题分类',

    -- 会话状态
    type VARCHAR(50) DEFAULT 'chat' COMMENT '会话类型',
    status VARCHAR(20) DEFAULT 'active' COMMENT '会话状态',
    modality_context_id VARCHAR(64) NULL COMMENT '多模态上下文ID',

    -- 性能指标
    message_count INT DEFAULT 0 COMMENT '消息数量',
    token_count BIGINT DEFAULT 0 COMMENT 'Token数量',
    first_response_time_ms INT NULL COMMENT '首次响应时间',
    resolution_status VARCHAR(20) DEFAULT 'unresolved' COMMENT '解决状态',
    satisfaction_score INT NULL COMMENT '满意度评分',

    -- 元数据
    conversation_tags JSON NULL COMMENT '会话标签',
    metadata JSON NULL COMMENT '元数据',

    -- 时间戳
    last_message_at TIMESTAMP NULL COMMENT '最后消息时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_conversation_id (conversation_id),
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_status (status),
    INDEX idx_topic_category (topic_category),
    INDEX idx_status_created (status, created_at),
    INDEX idx_tenant_status_created (tenant_id, status, created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

-- 3.2 消息表 (messages)
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '消息ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    message_id VARCHAR(64) NOT NULL COMMENT '消息唯一标识',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    user_id BIGINT NULL COMMENT '用户ID',

    -- 消息内容
    role VARCHAR(20) NOT NULL COMMENT '角色',
    content_type VARCHAR(50) DEFAULT 'text' COMMENT '内容类型',
    content TEXT NOT NULL COMMENT '消息内容',
    original_content TEXT NULL COMMENT '原始内容（编辑前）',

    -- 编辑历史
    edit_count INT DEFAULT 0 COMMENT '编辑次数',
    last_edited_at TIMESTAMP NULL COMMENT '最后编辑时间',

    -- 父子关系
    parent_message_id BIGINT NULL COMMENT '父消息ID',
    message_depth INT DEFAULT 0 COMMENT '消息深度',

    -- 状态和元数据
    status VARCHAR(20) DEFAULT 'sent' COMMENT '消息状态',
    metadata JSON NULL COMMENT '元数据',

    -- 交互数据
    read_by_users JSON NULL COMMENT '已读用户',
    reactions JSON NULL COMMENT '消息reactions',
    translated_content JSON NULL COMMENT '翻译内容',
    is_sensitive BOOLEAN DEFAULT FALSE COMMENT '是否敏感',
    processing_priority INT DEFAULT 5 COMMENT '处理优先级',

    -- Token统计
    token_count INT DEFAULT 0 COMMENT 'Token数量',

    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_message_id (message_id),
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_conversation_created (conversation_id, created_at),
    INDEX idx_user_id (user_id),
    INDEX idx_role (role),
    INDEX idx_content_type (content_type),
    INDEX idx_parent_message (parent_message_id),
    INDEX idx_message_depth (message_depth),
    INDEX idx_is_sensitive (is_sensitive),
    INDEX idx_processing_priority (processing_priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

-- 3.3 消息附件表 (message_attachments)
CREATE TABLE IF NOT EXISTS message_attachments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '附件ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    message_id VARCHAR(64) NOT NULL COMMENT '消息ID',
    attachment_id VARCHAR(64) NOT NULL COMMENT '附件唯一标识',

    -- 附件信息
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型',
    file_size BIGINT NOT NULL COMMENT '文件大小',
    file_url VARCHAR(500) NOT NULL COMMENT '文件URL',
    mime_type VARCHAR(100) NULL COMMENT 'MIME类型',

    -- 处理状态
    processing_status VARCHAR(20) DEFAULT 'pending' COMMENT '处理状态',
    processed_at TIMESTAMP NULL COMMENT '处理时间',

    -- 元数据
    metadata JSON NULL COMMENT '附件元数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_attachment_id (attachment_id),
    INDEX idx_message_id (message_id),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息附件表';

-- 3.4 会话反馈表 (conversation_feedback)
CREATE TABLE IF NOT EXISTS conversation_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '反馈ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',

    -- 反馈内容
    feedback_type VARCHAR(20) NOT NULL COMMENT '反馈类型',
    rating INT NULL COMMENT '评分(1-5)',
    comment TEXT NULL COMMENT '评论',
    suggestions TEXT NULL COMMENT '改进建议',

    -- 分类
    feedback_category VARCHAR(50) NULL COMMENT '反馈分类',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_user_id (user_id),
    INDEX idx_feedback_type (feedback_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话反馈表';

-- 3.5 会话记忆表 (conversation_memories)
CREATE TABLE IF NOT EXISTS conversation_memories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记忆ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    memory_type VARCHAR(50) NOT NULL COMMENT '记忆类型',

    -- 记忆内容
    memory_key VARCHAR(100) NOT NULL COMMENT '记忆键',
    memory_value TEXT NOT NULL COMMENT '记忆值',
    importance DECIMAL(5,4) DEFAULT 0.5000 COMMENT '重要性',

    -- 过期时间
    expires_at TIMESTAMP NULL COMMENT '过期时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_conversation_key (conversation_id, memory_key),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_memory_type (memory_type),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话记忆表';

-- 3.6 共享会话表 (shared_conversations)
CREATE TABLE IF NOT EXISTS shared_conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '共享ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    shared_by BIGINT NOT NULL COMMENT '分享人',

    -- 分享配置
    share_id VARCHAR(64) NOT NULL COMMENT '分享唯一标识',
    share_title VARCHAR(255) NULL COMMENT '分享标题',
    share_description TEXT NULL COMMENT '分享描述',

    -- 权限控制
    access_type VARCHAR(20) DEFAULT 'public' COMMENT '访问类型',
    password_hash VARCHAR(255) NULL COMMENT '访问密码',
    allowed_users JSON NULL COMMENT '允许访问的用户列表',

    -- 有效期
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    access_count INT DEFAULT 0 COMMENT '访问次数',

    -- 状态
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否有效',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_share_id (share_id),
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_expires_at (expires_at),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='共享会话表';

-- 3.7 会话日志归档表 (conversation_logs_archive)
CREATE TABLE IF NOT EXISTS conversation_logs_archive (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '日志ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    log_type VARCHAR(50) NOT NULL COMMENT '日志类型',
    log_content TEXT NOT NULL COMMENT '日志内容',
    metadata JSON NULL COMMENT '元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_log_type (log_type),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话日志归档表';

-- ================================================================================
-- Bot管理模块 (7张表)
-- ================================================================================

-- 3.8 Bot表 (bots)
CREATE TABLE IF NOT EXISTS bots (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT 'Bot ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot唯一标识',
    bot_alias VARCHAR(100) NULL COMMENT 'Bot别名',
    name VARCHAR(100) NOT NULL COMMENT 'Bot名称',
    description TEXT NULL COMMENT 'Bot描述',

    -- Bot类型
    type VARCHAR(50) NOT NULL COMMENT 'Bot类型',
    bot_category VARCHAR(50) DEFAULT 'general' COMMENT 'Bot分类',
    visibility VARCHAR(20) DEFAULT 'private' COMMENT '可见性',

    -- 配置
    config JSON NOT NULL COMMENT 'Bot配置',
    capabilities JSON NULL COMMENT '能力描述',
    icon_url VARCHAR(500) NULL COMMENT 'Bot图标',

    -- 版本管理
    version_id VARCHAR(64) NULL COMMENT '当前版本ID',
    version_number INT DEFAULT 1 COMMENT '版本号',

    -- 统计
    total_conversations BIGINT DEFAULT 0 COMMENT '总会话数',
    total_messages BIGINT DEFAULT 0 COMMENT '总消息数',
    avg_rating DECIMAL(3,2) NULL COMMENT '平均评分',

    -- 状态
    status VARCHAR(20) DEFAULT 'draft' COMMENT 'Bot状态',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否活跃',
    last_active_at TIMESTAMP NULL COMMENT '最后活跃时间',
    published_at TIMESTAMP NULL COMMENT '发布时间',

    -- 权限
    owner_id BIGINT NOT NULL COMMENT '所有者ID',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_bot_id (bot_id),
    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_bot_category (bot_category),
    INDEX idx_visibility (visibility),
    INDEX idx_status (status),
    INDEX idx_active (is_active, last_active_at),
    INDEX idx_owner_id (owner_id),
    FULLTEXT INDEX ft_name_description (name, description),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot表';

-- 3.9 Bot版本表 (bot_versions)
CREATE TABLE IF NOT EXISTS bot_versions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '版本ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    version_id VARCHAR(64) NOT NULL COMMENT '版本唯一标识',
    version_number INT NOT NULL COMMENT '版本号',
    version_name VARCHAR(100) NULL COMMENT '版本名称',

    -- 版本内容
    config JSON NOT NULL COMMENT 'Bot配置',
    diff_json JSON NULL COMMENT '版本差异',
    rollback_from_version_id BIGINT NULL COMMENT '回滚源版本',

    -- 标签和状态
    version_tags JSON NULL COMMENT '版本标签',
    deployment_env VARCHAR(20) DEFAULT 'development' COMMENT '部署环境',
    test_status VARCHAR(20) DEFAULT 'untested' COMMENT '测试状态',
    performance_baseline JSON NULL COMMENT '性能基准',

    -- 变更信息
    change_description TEXT NULL COMMENT '变更说明',
    change_log JSON NULL COMMENT '变更日志',

    -- 发布信息
    is_published BOOLEAN DEFAULT FALSE COMMENT '是否已发布',
    published_at TIMESTAMP NULL COMMENT '发布时间',
    published_by BIGINT NULL COMMENT '发布人',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_version_id (version_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_version_number (bot_id, version_number),
    INDEX idx_deployment_env (deployment_env),
    INDEX idx_test_status (test_status),
    INDEX idx_is_published (is_published)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot版本表';

-- 3.10 Bot技能表 (bot_skills)
CREATE TABLE IF NOT EXISTS bot_skills (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '技能ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    skill_id VARCHAR(64) NOT NULL COMMENT '技能唯一标识',
    skill_name VARCHAR(100) NOT NULL COMMENT '技能名称',
    skill_type VARCHAR(50) NOT NULL COMMENT '技能类型',

    -- 技能配置
    skill_config JSON NOT NULL COMMENT '技能配置',
    skill_parameters JSON NULL COMMENT '技能参数',

    -- 状态
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    skill_order INT DEFAULT 0 COMMENT '执行顺序',

    -- 描述
    description TEXT NULL COMMENT '技能描述',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_skill_id (skill_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_skill_type (skill_type),
    INDEX idx_is_enabled (is_enabled),
    INDEX idx_skill_order (skill_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot技能表';

-- 3.11 Bot访问控制表 (bot_access_controls)
CREATE TABLE IF NOT EXISTS bot_access_controls (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '访问控制ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    resource_id VARCHAR(64) NOT NULL COMMENT '资源ID',

    -- 访问权限
    access_level VARCHAR(20) NOT NULL COMMENT '访问级别',
    permissions JSON NOT NULL COMMENT '权限列表',

    -- 限制条件
    access_conditions JSON NULL COMMENT '访问条件',
    max_concurrent INT NULL COMMENT '最大并发数',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_access_level (access_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot访问控制表';

-- 3.12 Bot协作者表 (bot_collaborators)
CREATE TABLE IF NOT EXISTS bot_collaborators (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '协作者ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',

    -- 协作角色
    role VARCHAR(50) NOT NULL COMMENT '协作角色',
    permissions JSON NOT NULL COMMENT '权限列表',

    -- 状态
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',

    -- 时间
    invited_by BIGINT NULL COMMENT '邀请人',
    invited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '邀请时间',
    joined_at TIMESTAMP NULL COMMENT '加入时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_bot_user (bot_id, user_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role (role),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot协作者表';

-- 3.13 Bot知识库关联表 (bot_knowledge_bases)
CREATE TABLE IF NOT EXISTS bot_knowledge_bases (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    bot_id VARCHAR(64) NOT NULL COMMENT 'Bot ID',
    knowledge_base_id VARCHAR(64) NOT NULL COMMENT '知识库ID',

    -- 关联配置
    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否主知识库',
    priority INT DEFAULT 0 COMMENT '优先级',
    retrieval_config JSON NULL COMMENT '检索配置',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_bot_kb (bot_id, knowledge_base_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_knowledge_base_id (knowledge_base_id),
    INDEX idx_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot知识库关联表';

-- 3.14 Bot克隆记录表 (bot_clone_records)
CREATE TABLE IF NOT EXISTS bot_clone_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '克隆记录ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    source_bot_id VARCHAR(64) NOT NULL COMMENT '源Bot ID',
    target_bot_id VARCHAR(64) NOT NULL COMMENT '目标Bot ID',

    -- 克隆配置
    clone_config JSON NOT NULL COMMENT '克隆配置',
    clone_type VARCHAR(50) NOT NULL COMMENT '克隆类型',

    -- 状态
    clone_status VARCHAR(20) DEFAULT 'in_progress' COMMENT '克隆状态',
    progress DECIMAL(5,2) DEFAULT 0.00 COMMENT '克隆进度',

    -- 元数据
    error_message TEXT NULL COMMENT '错误信息',

    cloned_by BIGINT NOT NULL COMMENT '克隆人',
    cloned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '克隆时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_source_bot (source_bot_id),
    INDEX idx_target_bot (target_bot_id),
    INDEX idx_clone_status (clone_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot克隆记录表';

-- ================================================================================
-- 知识管理模块 (6张表)
-- ================================================================================

-- 3.15 知识库表 (knowledge_bases)
CREATE TABLE IF NOT EXISTS knowledge_bases (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '知识库ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    knowledge_base_id VARCHAR(64) NOT NULL COMMENT '知识库唯一标识',
    name VARCHAR(255) NOT NULL COMMENT '知识库名称',
    description TEXT NULL COMMENT '知识库描述',

    -- 知识库类型
    type VARCHAR(50) NOT NULL COMMENT '知识库类型',
    category VARCHAR(50) NULL COMMENT '知识库分类',

    -- RAG配置
    rag_config JSON NULL COMMENT 'RAG配置',

    -- 统计
    document_count INT DEFAULT 0 COMMENT '文档数量',
    total_chunks INT DEFAULT 0 COMMENT '总块数',
    total_tokens BIGINT DEFAULT 0 COMMENT '总Token数',
    avg_chunk_size INT DEFAULT 0 COMMENT '平均块大小',

    -- 同步状态
    sync_status VARCHAR(20) DEFAULT 'synced' COMMENT '同步状态',
    last_synced_at TIMESTAMP NULL COMMENT '最后同步时间',

    -- 权限
    owner_id BIGINT NOT NULL COMMENT '所有者ID',
    access_level VARCHAR(20) DEFAULT 'private' COMMENT '访问级别',

    -- 状态
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_knowledge_base_id (knowledge_base_id),
    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_type (type),
    INDEX idx_category (category),
    INDEX idx_sync_status (sync_status),
    INDEX idx_owner_id (owner_id),
    INDEX idx_status (status),
    FULLTEXT INDEX ft_name_description (name, description),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识库表';

-- 3.16 文档表 (documents)
CREATE TABLE IF NOT EXISTS documents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '文档ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    knowledge_base_id VARCHAR(64) NOT NULL COMMENT '知识库ID',
    document_id VARCHAR(64) NOT NULL COMMENT '文档唯一标识',
    title VARCHAR(500) NOT NULL COMMENT '文档标题',
    content LONGTEXT NULL COMMENT '文档内容',

    -- 文档信息
    file_name VARCHAR(255) NULL COMMENT '文件名',
    file_type VARCHAR(50) NULL COMMENT '文件类型',
    file_size BIGINT NULL COMMENT '文件大小',
    file_url VARCHAR(500) NULL COMMENT '文件URL',

    -- 处理状态
    processing_status VARCHAR(20) DEFAULT 'pending' COMMENT '处理状态',
    chunk_count INT DEFAULT 0 COMMENT '块数量',
    processed_at TIMESTAMP NULL COMMENT '处理时间',

    -- 元数据
    metadata JSON NULL COMMENT '文档元数据',
    tags JSON NULL COMMENT '标签',

    -- 权限
    owner_id BIGINT NULL COMMENT '所有者',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_document_id (document_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_knowledge_base_id (knowledge_base_id),
    INDEX idx_processing_status (processing_status),
    INDEX idx_owner_id (owner_id),
    FULLTEXT INDEX ft_title_content (title, content),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档表';

-- 3.17 文档块表 (document_chunks)
CREATE TABLE IF NOT EXISTS document_chunks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '块ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    knowledge_base_id VARCHAR(64) NOT NULL COMMENT '知识库ID',
    document_id VARCHAR(64) NOT NULL COMMENT '文档ID',
    chunk_id VARCHAR(64) NOT NULL COMMENT '块唯一标识',

    -- 块内容
    chunk_text TEXT NOT NULL COMMENT '块文本',
    chunk_type VARCHAR(50) DEFAULT 'text' COMMENT '块类型',
    chunk_size INT DEFAULT 0 COMMENT '块大小',
    chunk_language VARCHAR(10) DEFAULT 'zh' COMMENT '块语言',
    is_compressed BOOLEAN DEFAULT FALSE COMMENT '是否压缩',

    -- 嵌入向量
    embedding_vector BLOB NULL COMMENT '嵌入向量',
    embedding_model VARCHAR(100) NULL COMMENT '嵌入模型',
    quality_score DECIMAL(5,4) DEFAULT 0.5000 COMMENT '质量评分',

    -- 块哈希
    chunk_hash VARCHAR(64) NULL COMMENT '块哈希值',

    -- 位置信息
    chunk_index INT NOT NULL COMMENT '块索引',
    start_position INT NULL COMMENT '起始位置',
    end_position INT NULL COMMENT '结束位置',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_chunk_id (chunk_id),
    UNIQUE KEY uk_tenant_hash (tenant_id, chunk_hash),
    INDEX idx_document_id (document_id),
    INDEX idx_quality_score (quality_score)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档块表';

-- 3.18 文档解析器表 (document_parsers)
CREATE TABLE IF NOT EXISTS document_parsers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '解析器ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    parser_name VARCHAR(100) NOT NULL COMMENT '解析器名称',
    parser_type VARCHAR(50) NOT NULL COMMENT '解析器类型',

    -- 支持的文件类型
    supported_file_types JSON NOT NULL COMMENT '支持的文件类型列表',

    -- 解析配置
    parser_config JSON NOT NULL COMMENT '解析器配置',

    -- 状态
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_parser_name (parser_name),
    INDEX idx_parser_type (parser_type),
    INDEX idx_is_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文档解析器表';

-- 3.19 块策略表 (chunk_strategies)
CREATE TABLE IF NOT EXISTS chunk_strategies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '策略ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    strategy_name VARCHAR(100) NOT NULL COMMENT '策略名称',
    strategy_type VARCHAR(50) NOT NULL COMMENT '策略类型',

    -- 策略配置
    chunk_size INT DEFAULT 500 COMMENT '块大小',
    chunk_overlap INT DEFAULT 50 COMMENT '块重叠',
    strategy_config JSON NOT NULL COMMENT '策略配置',

    -- 适用范围
    applicable_file_types JSON NULL COMMENT '适用文件类型',

    -- 状态
    is_default BOOLEAN DEFAULT FALSE COMMENT '是否默认',
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_strategy_name (strategy_name),
    INDEX idx_strategy_type (strategy_type),
    INDEX idx_is_default (is_default)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='块策略表';

-- 3.20 知识块表 (knowledge_chunks)
CREATE TABLE IF NOT EXISTS knowledge_chunks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '知识块ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    knowledge_base_id VARCHAR(64) NOT NULL COMMENT '知识库ID',
    chunk_id VARCHAR(64) NOT NULL COMMENT '块唯一标识',

    -- 块内容
    content TEXT NOT NULL COMMENT '块内容',
    content_type VARCHAR(50) DEFAULT 'text' COMMENT '内容类型',

    -- 嵌入
    embedding_vector BLOB NULL COMMENT '嵌入向量',
    embedding_model VARCHAR(100) NULL COMMENT '嵌入模型',
    quality_score DECIMAL(5,4) DEFAULT 0.5000 COMMENT '质量评分',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',
    source_document_id VARCHAR(64) NULL COMMENT '来源文档',
    source_type VARCHAR(50) NULL COMMENT '来源类型',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_chunk_id (chunk_id),
    INDEX idx_knowledge_base_id (knowledge_base_id),
    INDEX idx_quality_score (quality_score)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='知识块表';

-- ================================================================================
-- 工作流模块 (5张表)
-- ================================================================================

-- 3.21 工作流元数据表 (workflow_meta)
CREATE TABLE IF NOT EXISTS workflow_meta (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '工作流ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    workflow_id VARCHAR(64) NOT NULL COMMENT '工作流唯一标识',
    name VARCHAR(255) NOT NULL COMMENT '工作流名称',
    description TEXT NULL COMMENT '工作流描述',

    -- 工作流类型
    mode VARCHAR(50) NOT NULL COMMENT '工作流模式',
    content_type VARCHAR(50) NOT NULL COMMENT '内容类型',
    category VARCHAR(50) NULL COMMENT '工作流分类',

    -- 标签
    tag VARCHAR(255) NULL COMMENT '标签',
    tags JSON NULL COMMENT '标签数组',

    -- 版本管理
    latest_version BIGINT DEFAULT 1 COMMENT '最新版本号',

    -- 统计
    complexity_score INT DEFAULT 0 COMMENT '复杂度评分',
    total_executions BIGINT DEFAULT 0 COMMENT '总执行次数',
    success_rate DECIMAL(5,4) DEFAULT 1.0000 COMMENT '成功率',
    avg_execution_time_ms INT DEFAULT 0 COMMENT '平均执行时间',

    -- 依赖
    dependencies JSON NULL COMMENT '依赖关系',

    -- 模板
    is_template BOOLEAN DEFAULT FALSE COMMENT '是否为模板',
    template_workflow_id BIGINT NULL COMMENT '模板来源',

    -- 状态
    status VARCHAR(20) DEFAULT 'draft' COMMENT '工作流状态',

    -- 权限
    owner_id BIGINT NOT NULL COMMENT '所有者ID',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_workflow_id (workflow_id),
    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_mode (mode),
    INDEX idx_category (category),
    INDEX idx_is_template (is_template),
    INDEX idx_status (status),
    INDEX idx_owner_id (owner_id),
    FULLTEXT INDEX ft_name_description (name, description),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流元数据表';

-- 3.22 工作流草稿表 (workflow_draft)
CREATE TABLE IF NOT EXISTS workflow_draft (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '草稿ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    workflow_id VARCHAR(64) NOT NULL COMMENT '工作流ID',
    version_number BIGINT NOT NULL COMMENT '版本号',

    -- 草稿内容
    content JSON NOT NULL COMMENT '工作流内容',
    draft_data JSON NOT NULL COMMENT '草稿数据',

    -- 自动保存
    is_autosave BOOLEAN DEFAULT FALSE COMMENT '是否自动保存',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    saved_by BIGINT NOT NULL COMMENT '保存人',
    saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '保存时间',

    UNIQUE KEY uk_workflow_version (workflow_id, version_number),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_saved_by (saved_by)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流草稿表';

-- 3.23 工作流版本表 (workflow_version)
CREATE TABLE IF NOT EXISTS workflow_version (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '版本ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    workflow_id VARCHAR(64) NOT NULL COMMENT '工作流ID',
    version_id VARCHAR(64) NOT NULL COMMENT '版本唯一标识',
    version_number BIGINT NOT NULL COMMENT '版本号',

    -- 版本内容
    content JSON NOT NULL COMMENT '工作流内容',
    version_hash VARCHAR(64) NOT NULL COMMENT '版本哈希',

    -- 变更说明
    change_description TEXT NULL COMMENT '变更说明',
    change_type VARCHAR(20) DEFAULT 'update' COMMENT '变更类型',

    -- 发布信息
    is_published BOOLEAN DEFAULT FALSE COMMENT '是否已发布',
    published_by BIGINT NULL COMMENT '发布人',
    published_at TIMESTAMP NULL COMMENT '发布时间',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_version_id (version_id),
    UNIQUE KEY uk_workflow_version (workflow_id, version_number),
    INDEX idx_is_published (is_published)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流版本表';

-- 3.24 工作流执行表 (workflow_execution)
CREATE TABLE IF NOT EXISTS workflow_execution (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '执行ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    workflow_id VARCHAR(64) NOT NULL COMMENT '工作流ID',
    execution_id VARCHAR(64) NOT NULL COMMENT '执行唯一标识',
    version_number BIGINT NOT NULL COMMENT '执行版本',

    -- 执行输入
    inputs JSON NOT NULL COMMENT '执行输入',

    -- 执行状态
    status VARCHAR(20) DEFAULT 'running' COMMENT '执行状态',
    error_message TEXT NULL COMMENT '错误信息',

    -- 性能指标
    execution_time_ms INT NULL COMMENT '执行耗时',
    token_count BIGINT DEFAULT 0 COMMENT 'Token数量',

    -- 执行人
    triggered_by BIGINT NOT NULL COMMENT '触发人',
    triggered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',

    UNIQUE KEY uk_execution_id (execution_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_workflow_id (workflow_id),
    INDEX idx_status (status),
    INDEX idx_triggered_by (triggered_by),
    INDEX idx_status_created (status, triggered_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流执行表';

-- 3.25 工作流快照表 (workflow_snapshot)
CREATE TABLE IF NOT EXISTS workflow_snapshot (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '快照ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    workflow_id VARCHAR(64) NOT NULL COMMENT '工作流ID',
    snapshot_id VARCHAR(64) NOT NULL COMMENT '快照唯一标识',

    -- 快照内容
    snapshot_data JSON NOT NULL COMMENT '快照数据',
    snapshot_type VARCHAR(50) NOT NULL COMMENT '快照类型',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_snapshot_id (snapshot_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_workflow_id (workflow_id),
    INDEX idx_snapshot_type (snapshot_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作流快照表';

-- 继续其他模块...

SET FOREIGN_KEY_CHECKS = 1;
