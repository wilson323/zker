-- ============================================================
-- 租户监控运维系统 - 数据库表创建脚本
-- ============================================================
-- 版本: v1.0.0
-- 创建日期: 2025-01-01
-- 数据库: MySQL 8.4.5
-- 字符集: utf8mb4
-- ============================================================

-- ============================================================
-- 1. 租户指标表 (tenant_metrics)
-- ============================================================
CREATE TABLE IF NOT EXISTS tenant_metrics (
    -- 主键
    id VARCHAR(36) PRIMARY KEY COMMENT '指标ID (UUID)',

    -- 租户关联
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 指标类型和值
    metric_type VARCHAR(50) NOT NULL COMMENT '指标类型: qps, response_time, error_rate, concurrency, cpu_usage, memory_usage, disk_io, network_io',
    metric_value DECIMAL(15,4) NOT NULL COMMENT '指标值',
    metric_timestamp TIMESTAMP NOT NULL COMMENT '指标时间戳',

    -- 标签 (JSON格式)
    tags JSON COMMENT '指标标签，如 {"endpoint": "/api/bots", "method": "GET"}',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 索引
    INDEX idx_tenant_id (tenant_id) COMMENT '租户ID索引',
    INDEX idx_metric_type (metric_type) COMMENT '指标类型索引',
    INDEX idx_timestamp (metric_timestamp) COMMENT '时间戳索引',
    INDEX idx_tenant_type_time (tenant_id, metric_type, metric_timestamp) COMMENT '联合索引用于查询',

    -- 约束
    CONSTRAINT chk_metric_type CHECK (metric_type IN ('qps', 'response_time', 'error_rate', 'concurrency', 'cpu_usage', 'memory_usage', 'disk_io', 'network_io'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户监控指标表';

-- ============================================================
-- 2. Agent指标表 (agent_metrics)
-- ============================================================
CREATE TABLE IF NOT EXISTS agent_metrics (
    -- 主键
    id VARCHAR(36) PRIMARY KEY COMMENT '指标ID (UUID)',

    -- 租户和Agent关联
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    agent_id VARCHAR(36) NOT NULL COMMENT 'Agent ID',
    agent_name VARCHAR(200) NOT NULL COMMENT 'Agent名称',

    -- 指标类型和值
    metric_type VARCHAR(50) NOT NULL COMMENT 'Agent指标类型: qps, response_time, error_rate, satisfaction, token_usage, cost, conversation, active_users',
    metric_value DECIMAL(15,4) NOT NULL COMMENT '指标值',
    metric_timestamp TIMESTAMP NOT NULL COMMENT '指标时间戳',

    -- 元数据 (JSON格式)
    metadata JSON COMMENT '额外元数据',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 索引
    INDEX idx_tenant_id (tenant_id) COMMENT '租户ID索引',
    INDEX idx_agent_id (agent_id) COMMENT 'Agent ID索引',
    INDEX idx_tenant_agent (tenant_id, agent_id) COMMENT '租户-Agent联合索引',
    INDEX idx_metric_type (metric_type) COMMENT '指标类型索引',
    INDEX idx_timestamp (metric_timestamp) COMMENT '时间戳索引',

    -- 约束
    CONSTRAINT chk_agent_metric_type CHECK (metric_type IN ('qps', 'response_time', 'error_rate', 'satisfaction', 'token_usage', 'cost', 'conversation', 'active_users'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent监控指标表';

-- ============================================================
-- 3. 告警规则表 (alert_rules)
-- ============================================================
CREATE TABLE IF NOT EXISTS alert_rules (
    -- 主键
    id VARCHAR(36) PRIMARY KEY COMMENT '告警规则ID (UUID)',

    -- 租户关联
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 规则基本信息
    rule_name VARCHAR(100) NOT NULL COMMENT '规则名称',
    description TEXT COMMENT '规则描述',

    -- 告警条件
    metric_type VARCHAR(50) NOT NULL COMMENT '监控指标类型',
    threshold DECIMAL(15,4) NOT NULL COMMENT '告警阈值',
    comparison VARCHAR(10) NOT NULL COMMENT '比较操作符: gt(>), lt(<), eq(=), gte(>=), lte(<=)',

    -- 告警级别
    severity VARCHAR(20) NOT NULL COMMENT '告警严重级别: warning, critical, emergency',

    -- 通知配置
    notification_channels JSON COMMENT '通知渠道列表: ["email", "sms", "webhook", "slack", "dingtalk", "wechat"]',
    notification_config JSON COMMENT '通知配置 (JSON格式)',

    -- 评估配置
    evaluation_interval INT NOT NULL DEFAULT 60 COMMENT '评估间隔（秒）',
    for_duration INT NOT NULL DEFAULT 300 COMMENT '持续时间（秒），达到阈值持续多久才触发告警',

    -- 状态
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间（软删除）',

    -- 索引
    INDEX idx_tenant_id (tenant_id) COMMENT '租户ID索引',
    INDEX idx_deleted_at (deleted_at) COMMENT '软删除索引',

    -- 约束
    CONSTRAINT chk_alert_comparison CHECK (comparison IN ('gt', 'lt', 'eq', 'gte', 'lte')),
    CONSTRAINT chk_alert_severity CHECK (severity IN ('warning', 'critical', 'emergency')),
    CONSTRAINT chk_evaluation_interval CHECK (evaluation_interval > 0),
    CONSTRAINT chk_for_duration CHECK (for_duration >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';

-- ============================================================
-- 4. 告警历史表 (alert_history)
-- ============================================================
CREATE TABLE IF NOT EXISTS alert_history (
    -- 主键
    id VARCHAR(36) PRIMARY KEY COMMENT '告警历史ID (UUID)',

    -- 租户和规则关联
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    alert_rule_id VARCHAR(36) NOT NULL COMMENT '告警规则ID',
    alert_rule_name VARCHAR(100) NOT NULL COMMENT '告警规则名称',

    -- 告警信息
    severity VARCHAR(20) NOT NULL COMMENT '告警严重级别',
    alert_message TEXT NOT NULL COMMENT '告警消息',
    alert_data JSON COMMENT '告警数据 (JSON格式)，如 {"metric_value": 95.5, "threshold": 80}',

    -- Agent关联（可选）
    agent_id VARCHAR(36) COMMENT 'Agent ID',

    -- 告警状态
    status VARCHAR(20) NOT NULL COMMENT '告警状态: pending, acknowledged, resolved, silenced',

    -- 确认信息
    acknowledged_by VARCHAR(36) COMMENT '确认人ID',
    acknowledged_at TIMESTAMP NULL COMMENT '确认时间',

    -- 解决信息
    resolved_at TIMESTAMP NULL COMMENT '解决时间',
    resolved_by VARCHAR(36) COMMENT '解决人ID',
    resolution_note TEXT COMMENT '解决备注',

    -- 静默信息
    silenced_until TIMESTAMP NULL COMMENT '静默截止时间',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 索引
    INDEX idx_tenant_id (tenant_id) COMMENT '租户ID索引',
    INDEX idx_alert_rule (alert_rule_id) COMMENT '告警规则ID索引',
    INDEX idx_severity (severity) COMMENT '严重级别索引',
    INDEX idx_status (status) COMMENT '状态索引',
    INDEX idx_agent_id (agent_id) COMMENT 'Agent ID索引',
    INDEX idx_created_at (created_at) COMMENT '创建时间索引',

    -- 约束
    CONSTRAINT chk_alert_history_severity CHECK (severity IN ('warning', 'critical', 'emergency')),
    CONSTRAINT chk_alert_history_status CHECK (status IN ('pending', 'acknowledged', 'resolved', 'silenced'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警历史表';

-- ============================================================
-- 5. 性能报告表 (performance_reports)
-- ============================================================
CREATE TABLE IF NOT EXISTS performance_reports (
    -- 主键
    id VARCHAR(36) PRIMARY KEY COMMENT '报告ID (UUID)',

    -- 租户和Agent关联
    tenant_id VARCHAR(36) NOT NULL COMMENT '租户ID',
    agent_id VARCHAR(36) COMMENT 'Agent ID (可选，为空表示租户级报告)',

    -- 报告基本信息
    report_type VARCHAR(50) NOT NULL COMMENT '报告类型: daily, weekly, monthly',
    report_date VARCHAR(20) NOT NULL COMMENT '报告日期 (YYYY-MM-DD)',

    -- QPS指标 (JSON格式)
    qps_metrics JSON COMMENT 'QPS指标，如 {"avg": 1000, "p50": 950, "p95": 1100, "p99": 1500}',

    -- 响应时间指标 (JSON格式)
    response_time_metrics JSON COMMENT '响应时间指标，如 {"avg": 50, "p50": 45, "p95": 80, "p99": 120}',

    -- 错误率
    error_rate DECIMAL(5,4) NOT NULL COMMENT '错误率 (0-1)',

    -- 满意度
    satisfaction DECIMAL(3,2) COMMENT '用户满意度 (0.00-1.00)',

    -- Token使用和成本
    token_usage BIGINT COMMENT 'Token使用量',
    estimated_cost DECIMAL(15,4) COMMENT '估算成本',

    -- 报告内容
    recommendations TEXT COMMENT '优化建议 (JSON数组格式)',
    summary TEXT COMMENT '报告摘要',
    top_issues JSON COMMENT '顶级问题列表 (JSON格式)',
    improvements JSON COMMENT '改进点 (JSON格式)',

    -- 生成信息
    generated_by VARCHAR(36) COMMENT '生成者（用户ID或system）',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    -- 索引
    INDEX idx_tenant_id (tenant_id) COMMENT '租户ID索引',
    INDEX idx_agent_id (agent_id) COMMENT 'Agent ID索引',
    INDEX idx_report_type (report_type) COMMENT '报告类型索引',
    INDEX idx_report_date (report_date) COMMENT '报告日期索引',

    -- 约束
    CONSTRAINT chk_report_type CHECK (report_type IN ('daily', 'weekly', 'monthly')),
    CONSTRAINT chk_error_rate CHECK (error_rate >= 0 AND error_rate <= 1),
    CONSTRAINT chk_satisfaction CHECK (satisfaction IS NULL OR (satisfaction >= 0 AND satisfaction <= 1))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='性能报告表';

-- ============================================================
-- 初始化数据
-- ============================================================

-- 创建默认告警规则模板（可选）
-- INSERT INTO alert_rules (id, tenant_id, rule_name, metric_type, threshold, comparison, severity, notification_channels)
-- VALUES
-- ('alert-rule-1', 'system', 'High Error Rate', 'error_rate', 0.05, 'gt', 'critical', '["email"]'),
-- ('alert-rule-2', 'system', 'Slow Response Time', 'response_time', 1000, 'gt', 'warning', '["email"]'),
-- ('alert-rule-3', 'system', 'Low QPS', 'qps', 10, 'lt', 'warning', '["email"]');

-- ============================================================
-- 数据保留策略（建议通过定时任务执行）
-- ============================================================
-- 删除30天前的指标数据
-- DELETE FROM tenant_metrics WHERE metric_timestamp < DATE_SUB(NOW(), INTERVAL 30 DAY);
-- DELETE FROM agent_metrics WHERE metric_timestamp < DATE_SUB(NOW(), INTERVAL 30 DAY);

-- 删除90天前的告警历史
-- DELETE FROM alert_history WHERE created_at < DATE_SUB(NOW(), INTERVAL 90 DAY);

-- 删除1年前的性能报告
-- DELETE FROM performance_reports WHERE created_at < DATE_SUB(NOW(), INTERVAL 1 YEAR);

-- ============================================================
-- 权限设置（根据实际需求调整）
-- ============================================================
-- GRANT SELECT, INSERT, UPDATE, DELETE ON coze_studio.* TO 'coze_app'@'%';
-- FLUSH PRIVILEGES;

-- ============================================================
-- 验证脚本
-- ============================================================
-- 验证表是否创建成功
SELECT
    TABLE_NAME,
    TABLE_ROWS,
    CREATE_TIME,
    UPDATE_TIME
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
AND TABLE_NAME IN ('tenant_metrics', 'agent_metrics', 'alert_rules', 'alert_history', 'performance_reports');

-- 验证索引是否创建成功
SELECT
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    INDEX_TYPE
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = DATABASE()
AND TABLE_NAME IN ('tenant_metrics', 'agent_metrics', 'alert_rules', 'alert_history', 'performance_reports')
ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX;
