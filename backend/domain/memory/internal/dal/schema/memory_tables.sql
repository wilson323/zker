-- ============================================
-- 记忆引擎 - 数据库表结构
-- 版本: v1.0.0
-- 创建日期: 2025-01-03
-- ============================================

-- -------------------------------------------
-- 表1: conversation_memories (对话记忆表)
-- -------------------------------------------
CREATE TABLE IF NOT EXISTS conversation_memories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记忆ID',
    memory_id VARCHAR(36) NOT NULL COMMENT '记忆唯一标识(UUID)',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id VARCHAR(64) NOT NULL COMMENT '用户ID',
    conversation_id VARCHAR(64) NOT NULL COMMENT '会话ID',

    memory_type ENUM('SUMMARY', 'ENTITY', 'PREFERENCE', 'EVENT') NOT NULL COMMENT '记忆类型',
    content TEXT NOT NULL COMMENT '记忆内容',
    embedding VECTOR(1536) COMMENT '向量嵌入(1536维)',

    importance_score DECIMAL(3,2) DEFAULT 0.50 COMMENT '重要性评分(0-1)',
    access_count INT DEFAULT 0 COMMENT '访问次数',
    last_accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后访问时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',

    metadata JSON COMMENT '元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_memory_id (memory_id),
    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_conversation (conversation_id),
    INDEX idx_type (memory_type),
    INDEX idx_expires (expires_at),
    INDEX idx_importance (importance_score),
    INDEX idx_tenant_user_conv (tenant_id, user_id, conversation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='对话记忆表';

-- -------------------------------------------
-- 表2: knowledge_memories (知识记忆表)
-- -------------------------------------------
CREATE TABLE IF NOT EXISTS knowledge_memories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '知识ID',
    memory_id VARCHAR(36) NOT NULL COMMENT '知识唯一标识(UUID)',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',

    knowledge_type ENUM('DOCUMENT', 'FAQ', 'PROCEDURE', 'CONCEPT') NOT NULL COMMENT '知识类型',
    title VARCHAR(255) NOT NULL COMMENT '知识标题',
    content LONGTEXT NOT NULL COMMENT '知识内容',
    source_uri VARCHAR(512) COMMENT '来源URI',
    embedding VECTOR(1536) COMMENT '向量嵌入(1536维)',

    metadata JSON COMMENT '元数据',
    quality_score DECIMAL(3,2) DEFAULT 0.50 COMMENT '质量评分(0-1)',
    version INT DEFAULT 1 COMMENT '版本号',
    access_count INT DEFAULT 0 COMMENT '访问次数',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_memory_id (memory_id),
    INDEX idx_tenant_type (tenant_id, knowledge_type),
    INDEX idx_tenant_quality (tenant_id, quality_score),
    FULLTEXT INDEX idx_content (title, content) WITH PARSER ngram
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='知识记忆表';

-- -------------------------------------------
-- 表3: memory_associations (知识关联表)
-- -------------------------------------------
CREATE TABLE IF NOT EXISTS memory_associations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',
    source_memory_id VARCHAR(36) NOT NULL COMMENT '源记忆ID',
    target_memory_id VARCHAR(36) NOT NULL COMMENT '目标记忆ID',
    association_type VARCHAR(64) NOT NULL COMMENT '关联类型',
    strength DECIMAL(3,2) DEFAULT 0.50 COMMENT '关联强度(0-1)',
    metadata JSON COMMENT '元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_source_target (source_memory_id, target_memory_id),
    INDEX idx_source (source_memory_id),
    INDEX idx_target (target_memory_id),
    INDEX idx_type (association_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='知识关联表';

-- -------------------------------------------
-- 索引优化说明
-- -------------------------------------------
-- 1. conversation_memories:
--    - uk_memory_id: 唯一约束，通过memory_id快速查询
--    - idx_tenant_user: 租户+用户查询(多租户隔离)
--    - idx_conversation: 会话级别查询
--    - idx_type: 按类型过滤
--    - idx_expires: 过期清理查询
--    - idx_importance: 重要性排序
--
-- 2. knowledge_memories:
--    - uk_memory_id: 唯一约束
--    - idx_tenant_type: 租户+类型查询
--    - idx_tenant_quality: 租户+质量排序
--    - idx_content: 全文检索
--
-- 3. memory_associations:
--    - uk_source_target: 防止重复关联
--    - idx_source/target: 双向查询
--    - idx_type: 按关联类型过滤
