-- Bot商店表创建脚本
-- 创建时间: 2025-12-30
-- 描述: 创建Bot商店相关的表

-- 创建bot_store_items表
CREATE TABLE IF NOT EXISTS bot_store_items (
    item_id VARCHAR(36) PRIMARY KEY COMMENT '商店项目ID',
    bot_id VARCHAR(36) NOT NULL COMMENT 'Bot ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    name VARCHAR(255) NOT NULL COMMENT 'Bot名称',
    description TEXT COMMENT 'Bot描述',
    category VARCHAR(100) COMMENT '分类',
    tags JSON COMMENT '标签列表',
    price DECIMAL(10,2) DEFAULT 0.00 COMMENT '价格（0表示免费）',
    publisher_id VARCHAR(36) NOT NULL COMMENT '发布者ID',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态：draft, pending, published, rejected, offline',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    download_count INT DEFAULT 0 COMMENT '下载次数',
    rating DECIMAL(3,2) DEFAULT 0.00 COMMENT '评分（0.00-5.00）',
    rating_count INT DEFAULT 0 COMMENT '评分人数',
    screenshots JSON COMMENT '截图列表',
    version VARCHAR(50) COMMENT '版本号',
    reject_reason TEXT COMMENT '拒绝原因',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间（软删除）',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_bot_id (bot_id),
    INDEX idx_status (status),
    INDEX idx_category (category),
    INDEX idx_publisher_id (publisher_id),
    INDEX idx_status_created (status, created_at),
    INDEX idx_rating_download (rating DESC, download_count DESC),
    UNIQUE KEY uk_bot_id (bot_id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (publisher_id) REFERENCES opencoze.user(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot商店项目表';

-- 创建bot_store_categories表
CREATE TABLE IF NOT EXISTS bot_store_categories (
    category_id VARCHAR(36) PRIMARY KEY COMMENT '分类ID',
    name VARCHAR(100) NOT NULL COMMENT '分类名称',
    icon VARCHAR(255) COMMENT '分类图标URL',
    description TEXT COMMENT '分类描述',
    parent_id VARCHAR(36) COMMENT '父分类ID',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否活跃',
    bot_count INT DEFAULT 0 COMMENT '该分类下的Bot数量',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间（软删除）',

    INDEX idx_parent_id (parent_id),
    INDEX idx_is_active (is_active),
    INDEX idx_sort_order (sort_order),
    UNIQUE KEY uk_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Bot商店分类表';

-- 插入默认分类数据
INSERT INTO bot_store_categories (category_id, name, icon, description, sort_order, is_active, bot_count) VALUES
('cat_productivity', '效率工具', '/icons/productivity.png', '提升工作效率的Bot', 1, TRUE, 0),
('cat_entertainment', '娱乐休闲', '/icons/entertainment.png', '娱乐和休闲类Bot', 2, TRUE, 0),
('cat_education', '教育学习', '/icons/education.png', '教育和学习类Bot', 3, TRUE, 0),
('cat_business', '商业金融', '/icons/business.png', '商业和金融相关Bot', 4, TRUE, 0),
('cat_development', '开发工具', '/icons/development.png', '软件开发相关Bot', 5, TRUE, 0),
('cat_health', '健康生活', '/icons/health.png', '健康和生活方式Bot', 6, TRUE, 0),
('cat_creative', '创意设计', '/icons/creative.png', '创意和设计类Bot', 7, TRUE, 0),
('cat_other', '其他', '/icons/other.png', '其他类别Bot', 99, TRUE, 0);
