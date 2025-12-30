-- ============================================
-- Token Metering 数据表创建脚本
-- 执行时机: Week 5
-- 最后更新: 2025-12-30
-- ============================================

-- ============================================
-- 表1: token_usage - Token使用记录表
-- ============================================
CREATE TABLE IF NOT EXISTS token_usage (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    bot_id VARCHAR(36) NOT NULL COMMENT 'Bot ID',
    model_id VARCHAR(100) NOT NULL COMMENT '模型ID',
    prompt_tokens INT DEFAULT 0 COMMENT '提示词Token数',
    completion_tokens INT DEFAULT 0 COMMENT '完成Token数',
    total_tokens INT DEFAULT 0 COMMENT '总Token数',
    cost_usd DECIMAL(10, 6) DEFAULT 0.000000 COMMENT '成本（美元）',
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '请求时间',

    INDEX idx_tenant_user (tenant_id, user_id),
    INDEX idx_model_date (model_id, requested_at),
    INDEX idx_tenant_date (tenant_id, requested_at),
    INDEX idx_bot_date (bot_id, requested_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token使用记录表';


-- ============================================
-- 表2: token_budgets - Token预算表
-- ============================================
CREATE TABLE IF NOT EXISTS token_budgets (
    budget_id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '预算ID',
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    model_id VARCHAR(100) NOT NULL COMMENT '模型ID',
    budget_type ENUM('daily', 'weekly', 'monthly') NOT NULL COMMENT '预算类型',
    max_tokens INT DEFAULT 1000000 COMMENT '最大Token数',
    alert_threshold_percent INT DEFAULT 80 COMMENT '告警阈值（百分比）',
    auto_degradation_model VARCHAR(100) COMMENT '自动降级模型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_tenant_model_type (tenant_id, model_id, budget_type),
    INDEX idx_tenant_type (tenant_id, budget_type),
    INDEX idx_model (model_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Token预算表';


-- ============================================
-- 验证脚本
-- ============================================

-- 1. 检查表是否创建成功
SHOW TABLES LIKE 'token%';

-- 2. 检查token_usage表结构
DESC token_usage;

-- 3. 检查token_budgets表结构
DESC token_budgets;

-- 4. 检查索引
SHOW INDEX FROM token_usage;
SHOW INDEX FROM token_budgets;


-- ============================================
-- 初始化测试数据（可选）
-- ============================================

-- 插入测试预算配置
-- INSERT INTO token_budgets (tenant_id, model_id, budget_type, max_tokens, alert_threshold_percent)
-- VALUES
--     ('tenant_1', 'gpt-4', 'monthly', 1000000, 80),
--     ('tenant_1', 'gpt-3.5-turbo', 'monthly', 5000000, 80);
