-- ================================================================================
-- ZKER 数据库DDL脚本 - 第四部分: 商业化与监控模块 (35张表)
-- ================================================================================
-- 版本: v1.0
-- 生成日期: 2025-01-01
-- 包含模块: 商业化、监控、多模态等
-- ================================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 商业化模块 (11张表)
-- ================================================================================

-- 4.1 订阅方案表 (subscription_plans)
CREATE TABLE IF NOT EXISTS subscription_plans (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '方案ID',
    plan_id VARCHAR(64) NOT NULL COMMENT '方案唯一标识',
    plan_name VARCHAR(100) NOT NULL COMMENT '方案名称',
    plan_type VARCHAR(50) NOT NULL COMMENT '方案类型',
    description TEXT NULL COMMENT '方案描述',

    -- 价格
    price DECIMAL(10,2) NOT NULL COMMENT '价格',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',
    billing_cycle VARCHAR(20) NOT NULL COMMENT '计费周期',

    -- 配额
    quotas JSON NOT NULL COMMENT '配额配置',

    -- 功能特性
    features JSON NOT NULL COMMENT '功能特性',

    -- 限制
    limits JSON NULL COMMENT '限制条件',

    -- 状态
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    is_public BOOLEAN DEFAULT TRUE COMMENT '是否公开',

    -- 显示顺序
    display_order INT DEFAULT 0 COMMENT '显示顺序',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_plan_id (plan_id),
    INDEX idx_plan_type (plan_type),
    INDEX idx_is_active (is_active),
    INDEX idx_is_public (is_public)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅方案表';

-- 4.2 订阅记录表 (subscriptions)
CREATE TABLE IF NOT EXISTS subscriptions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订阅ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    subscription_id VARCHAR(64) NOT NULL COMMENT '订阅唯一标识',
    plan_id VARCHAR(64) NOT NULL COMMENT '方案ID',

    -- 订阅周期
    start_date TIMESTAMP NOT NULL COMMENT '开始日期',
    end_date TIMESTAMP NULL COMMENT '结束日期',

    -- 自动续费
    auto_renew BOOLEAN DEFAULT FALSE COMMENT '是否自动续费',

    -- 状态
    status VARCHAR(20) DEFAULT 'active' COMMENT '订阅状态',

    -- 价格
    price DECIMAL(10,2) NOT NULL COMMENT '订阅价格',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_subscription_id (subscription_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_plan_id (plan_id),
    INDEX idx_status (status),
    INDEX idx_end_date (end_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订阅记录表';

-- 4.3 订单表 (orders)
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订单ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    order_id VARCHAR(64) NOT NULL COMMENT '订单唯一标识',
    order_no VARCHAR(64) NOT NULL COMMENT '订单号',

    -- 订单信息
    order_type VARCHAR(50) NOT NULL COMMENT '订单类型',
    plan_id VARCHAR(64) NULL COMMENT '方案ID',

    -- 金额
    subtotal DECIMAL(12,2) NOT NULL COMMENT '小计',
    tax DECIMAL(12,2) DEFAULT 0.00 COMMENT '税额',
    discount DECIMAL(12,2) DEFAULT 0.00 COMMENT '折扣',
    total DECIMAL(12,2) NOT NULL COMMENT '总额',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    -- 订单状态
    status VARCHAR(20) DEFAULT 'pending' COMMENT '订单状态',

    -- 用户信息
    customer_id BIGINT NOT NULL COMMENT '客户ID',

    -- 时间
    ordered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '下单时间',
    paid_at TIMESTAMP NULL COMMENT '支付时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_order_id (order_id),
    UNIQUE KEY uk_order_no (order_no),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status),
    INDEX idx_ordered_at (ordered_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- 4.4 订单项表 (order_items)
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '订单项ID',
    order_id VARCHAR(64) NOT NULL COMMENT '订单ID',
    item_id VARCHAR(64) NOT NULL COMMENT '订单项唯一标识',

    -- 商品信息
    product_type VARCHAR(50) NOT NULL COMMENT '商品类型',
    product_id VARCHAR(64) NOT NULL COMMENT '商品ID',
    product_name VARCHAR(255) NOT NULL COMMENT '商品名称',

    -- 数量和价格
    quantity INT DEFAULT 1 COMMENT '数量',
    unit_price DECIMAL(10,2) NOT NULL COMMENT '单价',
    total_price DECIMAL(12,2) NOT NULL COMMENT '总价',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_item_id (item_id),
    INDEX idx_order_id (order_id),
    INDEX idx_product (product_type, product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单项表';

-- 4.5 发票表 (invoices)
CREATE TABLE IF NOT EXISTS invoices (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '发票ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    invoice_id VARCHAR(64) NOT NULL COMMENT '发票唯一标识',
    invoice_no VARCHAR(64) NOT NULL COMMENT '发票号码',

    -- 账单周期
    billing_period_start DATE NOT NULL COMMENT '账期开始',
    billing_period_end DATE NOT NULL COMMENT '账期结束',

    -- 金额
    subtotal DECIMAL(12,2) NOT NULL COMMENT '小计',
    tax DECIMAL(12,2) DEFAULT 0.00 COMMENT '税额',
    total DECIMAL(12,2) NOT NULL COMMENT '总额',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    -- 状态
    status VARCHAR(20) DEFAULT 'unpaid' COMMENT '发票状态',

    -- 到期和支付
    due_date TIMESTAMP NOT NULL COMMENT '到期日期',
    paid_at TIMESTAMP NULL COMMENT '支付时间',
    payment_method VARCHAR(50) NULL COMMENT '支付方式',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_invoice_id (invoice_id),
    UNIQUE KEY uk_invoice_no (invoice_no),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发票表';

-- 4.6 发票项表 (invoice_items)
CREATE TABLE IF NOT EXISTS invoice_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '发票项ID',
    invoice_id VARCHAR(64) NOT NULL COMMENT '发票ID',

    -- 项目信息
    item_type VARCHAR(50) NOT NULL COMMENT '项目类型',
    item_name VARCHAR(255) NOT NULL COMMENT '项目名称',
    description TEXT NULL COMMENT '项目描述',

    -- 数量和价格
    quantity INT DEFAULT 1 COMMENT '数量',
    unit_price DECIMAL(10,2) NOT NULL COMMENT '单价',
    total_price DECIMAL(12,2) NOT NULL COMMENT '总价',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_invoice_id (invoice_id),
    INDEX idx_item_type (item_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发票项表';

-- 4.7 支付记录表 (payments)
CREATE TABLE IF NOT EXISTS payments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '支付ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    payment_id VARCHAR(64) NOT NULL COMMENT '支付唯一标识',
    payment_no VARCHAR(64) NOT NULL COMMENT '支付编号',

    -- 关联订单
    order_id VARCHAR(64) NOT NULL COMMENT '订单ID',
    invoice_id VARCHAR(64) NULL COMMENT '发票ID',

    -- 支付信息
    payment_method VARCHAR(50) NOT NULL COMMENT '支付方式',
    payment_channel VARCHAR(50) NULL COMMENT '支付渠道',

    -- 金额
    amount DECIMAL(12,2) NOT NULL COMMENT '支付金额',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    -- 状态
    status VARCHAR(20) DEFAULT 'pending' COMMENT '支付状态',

    -- 第三方信息
    transaction_id VARCHAR(100) NULL COMMENT '第三方交易ID',
    third_party_response JSON NULL COMMENT '第三方响应',

    -- 时间
    initiated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发起时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_payment_id (payment_id),
    UNIQUE KEY uk_payment_no (payment_no),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_order_id (order_id),
    INDEX idx_invoice_id (invoice_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付记录表';

-- 4.8 退款记录表 (refunds)
CREATE TABLE IF NOT EXISTS refunds (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '退款ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    refund_id VARCHAR(64) NOT NULL COMMENT '退款唯一标识',
    payment_id VARCHAR(64) NOT NULL COMMENT '支付ID',

    -- 退款信息
    refund_amount DECIMAL(12,2) NOT NULL COMMENT '退款金额',
    refund_reason TEXT NULL COMMENT '退款原因',

    -- 状态
    status VARCHAR(20) DEFAULT 'pending' COMMENT '退款状态',

    -- 第三方信息
    refund_transaction_id VARCHAR(100) NULL COMMENT '第三方退款交易ID',
    third_party_response JSON NULL COMMENT '第三方响应',

    -- 时间
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
    processed_at TIMESTAMP NULL COMMENT '处理时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_refund_id (refund_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_payment_id (payment_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款记录表';

-- 4.9 账单历史表 (billing_histories)
CREATE TABLE IF NOT EXISTS billing_histories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '历史ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 时间窗口
    billing_period_start DATE NOT NULL COMMENT '账期开始',
    billing_period_end DATE NOT NULL COMMENT '账期结束',

    -- 使用量
    usage_metrics JSON NOT NULL COMMENT '使用量指标',

    -- 费用
    subtotal DECIMAL(12,2) NOT NULL COMMENT '小计',
    tax DECIMAL(12,2) DEFAULT 0.00 COMMENT '税额',
    total DECIMAL(12,2) NOT NULL COMMENT '总额',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_billing_period (billing_period_start, billing_period_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账单历史表';

-- 4.10 预算设置表 (budget_settings)
CREATE TABLE IF NOT EXISTS budget_settings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '预算ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 预算配置
    budget_type VARCHAR(50) NOT NULL COMMENT '预算类型',
    budget_limit DECIMAL(12,2) NOT NULL COMMENT '预算限额',
    billing_cycle VARCHAR(20) NOT NULL COMMENT '计费周期',

    -- 告警阈值
    warning_threshold INT DEFAULT 80 COMMENT '警告阈值',
    critical_threshold INT DEFAULT 90 COMMENT '严重阈值',

    -- 通知配置
    notification_channels JSON NOT NULL COMMENT '通知渠道',

    -- 状态
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_budget_type (budget_type),
    INDEX idx_is_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预算设置表';

-- 4.11 预算告警表 (budget_alerts)
CREATE TABLE IF NOT EXISTS budget_alerts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '告警ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    budget_id BIGINT NOT NULL COMMENT '预算ID',

    -- 告警信息
    alert_type VARCHAR(20) NOT NULL COMMENT '告警类型',
    alert_level VARCHAR(20) NOT NULL COMMENT '告警级别',
    current_usage DECIMAL(12,2) NOT NULL COMMENT '当前使用量',
    budget_limit DECIMAL(12,2) NOT NULL COMMENT '预算限额',
    usage_percentage DECIMAL(5,2) NOT NULL COMMENT '使用百分比',

    -- 通知状态
    notification_sent BOOLEAN DEFAULT FALSE COMMENT '是否已通知',
    notification_sent_at TIMESTAMP NULL COMMENT '通知时间',

    -- 状态
    is_resolved BOOLEAN DEFAULT FALSE COMMENT '是否已解决',
    resolved_at TIMESTAMP NULL COMMENT '解决时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_budget_id (budget_id),
    INDEX idx_alert_type (alert_type),
    INDEX idx_alert_level (alert_level),
    INDEX idx_is_resolved (is_resolved)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预算告警表';

-- ================================================================================
-- 监控与评估模块 (12张表)
-- ================================================================================

-- 4.12 Agent执行记录表 (agent_executions)
CREATE TABLE IF NOT EXISTS agent_executions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '执行ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    agent_id VARCHAR(64) NOT NULL COMMENT 'Agent ID',
    execution_id VARCHAR(64) NOT NULL COMMENT '执行唯一标识',

    -- 执行输入
    input JSON NOT NULL COMMENT '执行输入',

    -- 执行状态
    status VARCHAR(20) DEFAULT 'running' COMMENT '执行状态',
    error_message TEXT NULL COMMENT '错误信息',

    -- 性能指标
    execution_time_ms INT NULL COMMENT '执行耗时',
    token_count BIGINT DEFAULT 0 COMMENT 'Token数量',

    -- 触发信息
    triggered_by BIGINT NOT NULL COMMENT '触发人',
    triggered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_execution_id (execution_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_status (status),
    INDEX idx_triggered_by (triggered_by),
    INDEX idx_status_created (status, triggered_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent执行记录表';

-- 4.13 Agent执行步骤表 (agent_execution_steps)
CREATE TABLE IF NOT EXISTS agent_execution_steps (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '步骤ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    execution_id VARCHAR(64) NOT NULL COMMENT '执行ID',

    -- 步骤信息
    step_number INT NOT NULL COMMENT '步骤序号',
    step_type VARCHAR(50) NOT NULL COMMENT '步骤类型',
    step_name VARCHAR(100) NOT NULL COMMENT '步骤名称',

    -- 步骤输入输出
    step_input JSON NULL COMMENT '步骤输入',
    step_output JSON NULL COMMENT '步骤输出',

    -- 状态
    status VARCHAR(20) DEFAULT 'pending' COMMENT '步骤状态',
    error_message TEXT NULL COMMENT '错误信息',

    -- 性能
    execution_time_ms INT NULL COMMENT '执行耗时',
    token_count BIGINT DEFAULT 0 COMMENT 'Token数量',

    -- 时间
    started_at TIMESTAMP NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_execution_id (execution_id),
    INDEX idx_step_number (step_number),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent执行步骤表';

-- 4.14 Agent性能指标表 (agent_performance_metrics)
CREATE TABLE IF NOT EXISTS agent_performance_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '指标ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    agent_id VARCHAR(64) NOT NULL COMMENT 'Agent ID',

    -- 时间窗口
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口',
    window_start TIMESTAMP NOT NULL COMMENT '窗口开始',
    window_end TIMESTAMP NOT NULL COMMENT '窗口结束',

    -- 执行统计
    total_executions BIGINT DEFAULT 0 COMMENT '总执行次数',
    successful_executions BIGINT DEFAULT 0 COMMENT '成功执行次数',
    failed_executions BIGINT DEFAULT 0 COMMENT '失败执行次数',

    -- 性能指标
    avg_execution_time_ms INT DEFAULT 0 COMMENT '平均执行时间',
    p95_execution_time_ms INT DEFAULT 0 COMMENT 'P95执行时间',
    p99_execution_time_ms INT DEFAULT 0 COMMENT 'P99执行时间',

    -- Token统计
    total_token_count BIGINT DEFAULT 0 COMMENT '总Token数',
    avg_token_count INT DEFAULT 0 COMMENT '平均Token数',

    -- 成功率
    success_rate DECIMAL(5,4) DEFAULT 1.0000 COMMENT '成功率',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_agent_time_window (agent_id, time_window, window_start),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_time_window (time_window, window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent性能指标表';

-- 4.15 Agent质量指标表 (agent_quality_metrics)
CREATE TABLE IF NOT EXISTS agent_quality_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '质量ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    agent_id VARCHAR(64) NOT NULL COMMENT 'Agent ID',

    -- 时间窗口
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口',
    window_start TIMESTAMP NOT NULL COMMENT '窗口开始',
    window_end TIMESTAMP NOT NULL COMMENT '窗口结束',

    -- 用户反馈
    total_feedbacks INT DEFAULT 0 COMMENT '总反馈数',
    positive_feedbacks INT DEFAULT 0 COMMENT '正面反馈数',
    negative_feedbacks INT DEFAULT 0 COMMENT '负面反馈数',

    -- 满意度
    avg_satisfaction_score DECIMAL(3,2) NULL COMMENT '平均满意度',

    -- 准确性
    accuracy_score DECIMAL(5,4) NULL COMMENT '准确性评分',

    -- 响应相关性
    relevance_score DECIMAL(5,4) NULL COMMENT '相关性评分',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_agent_time_window (agent_id, time_window, window_start),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent质量指标表';

-- 4.16 Agent告警规则表 (agent_alert_rules)
CREATE TABLE IF NOT EXISTS agent_alert_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '告警规则ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    agent_id VARCHAR(64) NOT NULL COMMENT 'Agent ID',
    rule_id VARCHAR(64) NOT NULL COMMENT '规则唯一标识',

    -- 监控指标
    metric_name VARCHAR(50) NOT NULL COMMENT '指标名称',
    condition_operator VARCHAR(10) NOT NULL COMMENT '条件操作符',
    threshold_value DECIMAL(15,4) NOT NULL COMMENT '阈值',

    -- 告警配置
    alert_level VARCHAR(20) DEFAULT 'warning' COMMENT '告警级别',
    notification_channels JSON NOT NULL COMMENT '通知渠道',

    -- 冷却配置
    cooldown_period INT DEFAULT 300 COMMENT '冷却期(秒)',

    -- 状态
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_rule_id (rule_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_agent_id (agent_id),
    INDEX idx_is_enabled (is_enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent告警规则表';

-- 4.17 Agent告警历史表 (agent_alert_history)
CREATE TABLE IF NOT EXISTS agent_alert_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '告警历史ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    rule_id VARCHAR(64) NOT NULL COMMENT '规则ID',
    alert_id VARCHAR(64) NOT NULL COMMENT '告警ID',

    -- 告警内容
    alert_title VARCHAR(255) NOT NULL COMMENT '告警标题',
    alert_message TEXT NOT NULL COMMENT '告警消息',
    alert_level VARCHAR(20) NOT NULL COMMENT '告警级别',

    -- 触发信息
    metric_name VARCHAR(50) NOT NULL COMMENT '指标名称',
    actual_value DECIMAL(15,4) NOT NULL COMMENT '实际值',
    threshold_value DECIMAL(15,4) NOT NULL COMMENT '阈值',

    -- 状态
    status VARCHAR(20) DEFAULT 'firing' COMMENT '告警状态',
    acknowledged_at TIMESTAMP NULL COMMENT '确认时间',
    resolved_at TIMESTAMP NULL COMMENT '解决时间',

    fired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_rule_id (rule_id),
    INDEX idx_status (status),
    INDEX idx_fired_at (fired_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Agent告警历史表';

-- 4.18 使用指标表 (usage_metrics)
CREATE TABLE IF NOT EXISTS usage_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '指标ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    metric_type VARCHAR(50) NOT NULL COMMENT '指标类型',
    metric_name VARCHAR(100) NOT NULL COMMENT '指标名称',

    -- 时间窗口
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口',
    window_size INT NOT NULL COMMENT '窗口大小',
    window_start TIMESTAMP NOT NULL COMMENT '窗口开始',
    window_end TIMESTAMP NOT NULL COMMENT '窗口结束',

    -- 指标值
    metric_value BIGINT NOT NULL COMMENT '指标值',

    -- 维度
    dimensions JSON NULL COMMENT '维度标签',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_metric_time_window (tenant_id, metric_type, time_window, window_start),
    INDEX idx_metric_type (metric_type),
    INDEX idx_time_window (time_window, window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='使用指标表';

-- 4.19 使用记录表 (usage_records)
CREATE TABLE IF NOT EXISTS usage_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NULL COMMENT '用户ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    resource_id VARCHAR(64) NOT NULL COMMENT '资源ID',

    -- 使用信息
    action VARCHAR(50) NOT NULL COMMENT '操作',
    usage_amount BIGINT DEFAULT 1 COMMENT '使用量',
    usage_unit VARCHAR(20) DEFAULT 'count' COMMENT '使用单位',

    -- 时间
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发生时间',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_action (action),
    INDEX idx_occurred_at (occurred_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='使用记录表';

-- 4.20 用户行为表 (user_behaviors)
CREATE TABLE IF NOT EXISTS user_behaviors (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '行为ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',

    -- 行为信息
    behavior_type VARCHAR(50) NOT NULL COMMENT '行为类型',
    resource_type VARCHAR(50) NULL COMMENT '资源类型',
    resource_id VARCHAR(64) NULL COMMENT '资源ID',

    -- 行为详情
    behavior_data JSON NULL COMMENT '行为数据',

    -- 上下文
    context JSON NULL COMMENT '上下文信息',

    -- 时间
    occurred_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发生时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_behavior_type (behavior_type),
    INDEX idx_occurred_at (occurred_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户行为表';

-- 4.21 成本分析表 (cost_analytics)
CREATE TABLE IF NOT EXISTS cost_analytics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分析ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 时间窗口
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口',
    window_start TIMESTAMP NOT NULL COMMENT '窗口开始',
    window_end TIMESTAMP NOT NULL COMMENT '窗口结束',

    -- 成本明细
    total_cost DECIMAL(12,2) NOT NULL COMMENT '总成本',
    compute_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT '计算成本',
    storage_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT '存储成本',
    network_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT '网络成本',
    token_cost DECIMAL(12,2) DEFAULT 0.00 COMMENT 'Token成本',

    -- 成本分析
    cost_by_service JSON NULL COMMENT '按服务分类成本',
    cost_by_user JSON NULL COMMENT '按用户分类成本',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_time_window (time_window, window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成本分析表';

-- 4.22 成本优化建议表 (cost_optimization_suggestions)
CREATE TABLE IF NOT EXISTS cost_optimization_suggestions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '建议ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',

    -- 建议信息
    suggestion_type VARCHAR(50) NOT NULL COMMENT '建议类型',
    suggestion_title VARCHAR(255) NOT NULL COMMENT '建议标题',
    suggestion_description TEXT NOT NULL COMMENT '建议描述',

    -- 预估节省
    estimated_savings DECIMAL(12,2) NOT NULL COMMENT '预估节省',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    -- 优先级
    priority VARCHAR(20) DEFAULT 'medium' COMMENT '优先级',

    -- 状态
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态',

    -- 实施信息
    implemented_at TIMESTAMP NULL COMMENT '实施时间',
    actual_savings DECIMAL(12,2) NULL COMMENT '实际节省',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_suggestion_type (suggestion_type),
    INDEX idx_priority (priority),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='成本优化建议表';

-- 4.23 收入记录表 (revenue_records)
CREATE TABLE IF NOT EXISTS revenue_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '收入ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    revenue_id VARCHAR(64) NOT NULL COMMENT '收入唯一标识',

    -- 收入信息
    revenue_type VARCHAR(50) NOT NULL COMMENT '收入类型',
    revenue_source VARCHAR(50) NOT NULL COMMENT '收入来源',
    amount DECIMAL(12,2) NOT NULL COMMENT '收入金额',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '货币',

    -- 时间
    period_start DATE NOT NULL COMMENT '周期开始',
    period_end DATE NOT NULL COMMENT '周期结束',
    recognized_at TIMESTAMP NOT NULL COMMENT '确认时间',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_revenue_id (revenue_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_revenue_type (revenue_type),
    INDEX idx_period (period_start, period_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收入记录表';

-- ================================================================================
-- 多模态模块 (7张表)
-- ================================================================================

-- 4.24 音频文件表 (audio_files)
CREATE TABLE IF NOT EXISTS audio_files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '音频ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    file_id VARCHAR(64) NOT NULL COMMENT '文件唯一标识',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型',
    file_size BIGINT NOT NULL COMMENT '文件大小',
    file_url VARCHAR(500) NOT NULL COMMENT '文件URL',
    duration INT NOT NULL COMMENT '时长(秒)',

    -- 音频属性
    sample_rate INT NULL COMMENT '采样率',
    channels INT NULL COMMENT '声道数',
    bit_rate INT NULL COMMENT '比特率',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_file_id (file_id),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='音频文件表';

-- 4.25 视频文件表 (video_files)
CREATE TABLE IF NOT EXISTS video_files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '视频ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    file_id VARCHAR(64) NOT NULL COMMENT '文件唯一标识',
    file_name VARCHAR(255) NOT NULL COMMENT '文件名',
    file_type VARCHAR(50) NOT NULL COMMENT '文件类型',
    file_size BIGINT NOT NULL COMMENT '文件大小',
    file_url VARCHAR(500) NOT NULL COMMENT '文件URL',
    duration INT NOT NULL COMMENT '时长(秒)',

    -- 视频属性
    width INT NULL COMMENT '宽度',
    height INT NULL COMMENT '高度',
    frame_rate DECIMAL(5,2) NULL COMMENT '帧率',
    bit_rate INT NULL COMMENT '比特率',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_file_id (file_id),
    INDEX idx_tenant_id (tenant_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频文件表';

-- 4.26 语音记录表 (speech_records)
CREATE TABLE IF NOT EXISTS speech_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '记录ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    audio_file_id VARCHAR(64) NOT NULL COMMENT '音频文件ID',

    -- 转录信息
    transcription_status VARCHAR(20) DEFAULT 'pending' COMMENT '转录状态',
    transcription_text TEXT NULL COMMENT '转录文本',
    language VARCHAR(10) DEFAULT 'zh' COMMENT '语言',

    -- 处理信息
    processed_at TIMESTAMP NULL COMMENT '处理时间',
    processing_time_ms INT NULL COMMENT '处理耗时',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_audio_file_id (audio_file_id),
    INDEX idx_transcription_status (transcription_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='语音记录表';

-- 4.27 多模态会话表 (multimodal_sessions)
CREATE TABLE IF NOT EXISTS multimodal_sessions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会话ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    session_id VARCHAR(64) NOT NULL COMMENT '会话唯一标识',
    user_id BIGINT NOT NULL COMMENT '用户ID',

    -- 模态类型
    modalities JSON NOT NULL COMMENT '模态类型列表',

    -- 状态
    status VARCHAR(20) DEFAULT 'active' COMMENT '会话状态',

    -- 时间
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
    ended_at TIMESTAMP NULL COMMENT '结束时间',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_session_id (session_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='多模态会话表';

-- 4.28 模态上下文表 (modality_contexts)
CREATE TABLE IF NOT EXISTS modality_contexts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '上下文ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    context_id VARCHAR(64) NOT NULL COMMENT '上下文唯一标识',
    modality_type VARCHAR(50) NOT NULL COMMENT '模态类型',

    -- 关联
    conversation_id VARCHAR(64) NULL COMMENT '会话ID',
    multimodal_session_id VARCHAR(64) NULL COMMENT '多模态会话ID',

    -- 上下文数据
    context_data JSON NOT NULL COMMENT '上下文数据',

    -- 状态
    status VARCHAR(20) DEFAULT 'active' COMMENT '状态',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_context_id (context_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_modality_type (modality_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模态上下文表';

-- 4.29 视频分析表 (video_analysis)
CREATE TABLE IF NOT EXISTS video_analysis (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '分析ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    video_file_id VARCHAR(64) NOT NULL COMMENT '视频文件ID',

    -- 分析类型
    analysis_type VARCHAR(50) NOT NULL COMMENT '分析类型',

    -- 分析结果
    analysis_result JSON NOT NULL COMMENT '分析结果',
    analysis_confidence DECIMAL(5,4) NULL COMMENT '分析置信度',

    -- 状态
    analysis_status VARCHAR(20) DEFAULT 'pending' COMMENT '分析状态',

    -- 时间
    analyzed_at TIMESTAMP NULL COMMENT '分析时间',
    processing_time_ms INT NULL COMMENT '处理耗时',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_video_file_id (video_file_id),
    INDEX idx_analysis_type (analysis_type),
    INDEX idx_analysis_status (analysis_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频分析表';

-- 4.30 视频帧表 (video_frames)
CREATE TABLE IF NOT EXISTS video_frames (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '帧ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    video_file_id VARCHAR(64) NOT NULL COMMENT '视频文件ID',
    frame_id VARCHAR(64) NOT NULL COMMENT '帧唯一标识',

    -- 帧信息
    frame_number INT NOT NULL COMMENT '帧序号',
    frame_time DECIMAL(10,3) NOT NULL COMMENT '帧时间(秒)',
    frame_url VARCHAR(500) NOT NULL COMMENT '帧URL',

    -- 帧属性
    width INT NOT NULL COMMENT '宽度',
    height INT NOT NULL COMMENT '高度',

    -- 分析结果
    analysis_data JSON NULL COMMENT '分析数据',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_frame_id (frame_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_video_file_id (video_file_id),
    INDEX idx_frame_number (frame_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频帧表';

SET FOREIGN_KEY_CHECKS = 1;
