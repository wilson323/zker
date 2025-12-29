-- ================================================================================
-- ZKER 数据库DDL脚本 - 第二部分: 智能路由引擎 (18张表)
-- ================================================================================
-- 版本: v1.0
-- 生成日期: 2025-01-01
-- 优先级: P0 (核心功能)
-- ================================================================================

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ================================================================================
-- 智能路由引擎核心表 (18张表)
-- ================================================================================

-- 1. 路由规则表 (routing_rules)
CREATE TABLE IF NOT EXISTS routing_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    rule_id VARCHAR(64) NOT NULL COMMENT '规则唯一标识',
    version_id VARCHAR(64) NULL COMMENT '规则版本ID',
    name VARCHAR(100) NOT NULL COMMENT '规则名称',
    description TEXT NULL COMMENT '规则描述',
    priority INT NOT NULL DEFAULT 0 COMMENT '优先级，数字越大优先级越高',
    condition JSON NOT NULL COMMENT '触发条件',
    action JSON NOT NULL COMMENT '路由动作',
    rule_hash VARCHAR(64) NULL COMMENT '规则内容哈希值',
    effective_start_time TIMESTAMP NULL COMMENT '生效开始时间',
    effective_end_time TIMESTAMP NULL COMMENT '生效结束时间',
    status VARCHAR(20) DEFAULT 'enabled' COMMENT '状态',
    test_coverage DECIMAL(5,2) DEFAULT 0.00 COMMENT '测试覆盖率',
    last_tested_at TIMESTAMP NULL COMMENT '最后测试时间',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_rule_id (rule_id),
    UNIQUE KEY uk_tenant_name (tenant_id, name),
    INDEX idx_tenant_priority (tenant_id, priority),
    INDEX idx_version_id (version_id),
    INDEX idx_effective_time (effective_start_time, effective_end_time),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由规则表';

-- 2. 路由决策表 (routing_decisions)
CREATE TABLE IF NOT EXISTS routing_decisions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    decision_id VARCHAR(64) NOT NULL COMMENT '决策唯一标识',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    session_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    batch_id VARCHAR(64) NULL COMMENT '批处理ID',
    rule_id BIGINT NULL COMMENT '使用的规则ID',
    decision_version VARCHAR(64) NULL COMMENT '决策版本',

    -- 意图识别结果
    intent_type VARCHAR(50) NOT NULL COMMENT '意图类型',
    intent_confidence DECIMAL(5,4) NOT NULL COMMENT '意图置信度',
    entities JSON NULL COMMENT '提取的实体',

    -- 路由选择
    selected_service_id VARCHAR(64) NOT NULL COMMENT '选中的服务ID',
    selected_service_type VARCHAR(50) NOT NULL COMMENT '服务类型',
    candidates JSON NULL COMMENT '候选服务列表',
    decision_reason TEXT NULL COMMENT '决策理由',
    decision_explanation JSON NULL COMMENT '决策解释',
    decision_cost DECIMAL(10,4) DEFAULT 0.0000 COMMENT '决策成本',

    -- 性能指标
    decision_time_ms INT NOT NULL COMMENT '决策耗时(毫秒)',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_decision_id (decision_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_batch_id (batch_id),
    INDEX idx_intent_type (intent_type),
    INDEX idx_service_created (selected_service_id, created_at),
    INDEX idx_tenant_created (tenant_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由决策表';

-- 3. 服务候选表 (service_candidates)
CREATE TABLE IF NOT EXISTS service_candidates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    service_id VARCHAR(64) NOT NULL COMMENT '服务唯一标识',
    service_name VARCHAR(100) NOT NULL COMMENT '服务名称',
    service_type VARCHAR(50) NOT NULL COMMENT '服务类型',
    service_config JSON NOT NULL COMMENT '服务配置',
    service_tags JSON NULL COMMENT '服务标签',

    -- 健康状态
    health_status VARCHAR(20) DEFAULT 'healthy' COMMENT '健康状态',
    warmup_status VARCHAR(20) DEFAULT 'cold' COMMENT '预热状态',

    -- 容量管理
    capacity_limit INT DEFAULT 100 COMMENT '容量限制',
    current_load INT DEFAULT 0 COMMENT '当前负载',
    max_concurrent_requests INT DEFAULT 50 COMMENT '最大并发请求数',

    -- 性能指标
    avg_response_time_ms INT DEFAULT 0 COMMENT '平均响应时间',
    p95_response_time_ms INT DEFAULT 0 COMMENT 'P95响应时间',
    p99_response_time_ms INT DEFAULT 0 COMMENT 'P99响应时间',
    cold_start_time_ms INT DEFAULT 0 COMMENT '冷启动时间',

    -- 可用性
    is_available BOOLEAN DEFAULT TRUE COMMENT '是否可用',
    last_health_check_at TIMESTAMP NULL COMMENT '最后健康检查时间',

    -- 元数据
    description TEXT NULL COMMENT '服务描述',
    version VARCHAR(50) NULL COMMENT '服务版本',
    owner_id BIGINT NULL COMMENT '负责人ID',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_service_id (service_id),
    UNIQUE KEY uk_tenant_name (tenant_id, service_name),
    INDEX idx_service_type (service_type),
    INDEX idx_health_status (health_status),
    INDEX idx_warmup_status (warmup_status),
    INDEX idx_availability (is_available, deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务候选表';

-- 4. 路由反馈表 (routing_feedback)
CREATE TABLE IF NOT EXISTS routing_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    decision_id VARCHAR(64) NOT NULL COMMENT '关联的决策ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    session_id VARCHAR(64) NOT NULL COMMENT '会话ID',

    -- 反馈内容
    feedback_type VARCHAR(20) NOT NULL COMMENT '反馈类型',
    user_rating INT NULL COMMENT '用户评分(1-5)',
    user_comment TEXT NULL COMMENT '用户评论',
    is_correct BOOLEAN NULL COMMENT '路由是否正确',

    -- 分析数据
    actual_intent VARCHAR(50) NULL COMMENT '实际意图',
    should_route_to VARCHAR(64) NULL COMMENT '应该路由到的服务',
    improvement_suggestion TEXT NULL COMMENT '改进建议',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_decision_id (decision_id),
    INDEX idx_user_id (user_id),
    INDEX idx_user_rating_time (user_id, user_rating, created_at),
    INDEX idx_feedback_type (feedback_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由反馈表';

-- 5. 意图表 (intents)
CREATE TABLE IF NOT EXISTS intents (
    id VARCHAR(64) PRIMARY KEY COMMENT '意图唯一标识',
    type VARCHAR(50) NOT NULL UNIQUE COMMENT '意图类型',
    name VARCHAR(100) NOT NULL COMMENT '意图名称',
    display_name VARCHAR(100) NOT NULL COMMENT '显示名称',
    description TEXT NULL COMMENT '意图描述',
    category VARCHAR(50) NULL COMMENT '意图分类',
    parent_intent VARCHAR(50) NULL COMMENT '父意图',

    -- 阈值配置
    confidence_threshold DECIMAL(5,4) DEFAULT 0.7000 COMMENT '置信度阈值',

    -- 识别方式
    enable_llm BOOLEAN DEFAULT TRUE COMMENT '是否启用LLM识别',
    enable_rule BOOLEAN DEFAULT TRUE COMMENT '是否启用规则识别',
    enable_vector BOOLEAN DEFAULT TRUE COMMENT '是否启用向量识别',

    -- 识别配置
    llm_config JSON NULL COMMENT 'LLM识别配置',
    rule_config JSON NULL COMMENT '规则识别配置',
    vector_config JSON NULL COMMENT '向量识别配置',

    -- 示例
    examples JSON NULL COMMENT '示例列表',
    synonyms JSON NULL COMMENT '同义词列表',

    -- 元数据
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    INDEX idx_category (category),
    INDEX idx_parent_intent (parent_intent),
    INDEX idx_is_active (is_active),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='意图表';

-- 6. 意图样本表 (intent_samples)
CREATE TABLE IF NOT EXISTS intent_samples (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    intent_type VARCHAR(50) NOT NULL COMMENT '意图类型',
    sample_text TEXT NOT NULL COMMENT '样本文本',
    entities JSON NULL COMMENT '标注的实体',
    language VARCHAR(10) DEFAULT 'zh' COMMENT '语言',
    source VARCHAR(50) DEFAULT 'manual' COMMENT '样本来源',
    quality_score DECIMAL(5,4) NULL COMMENT '质量评分',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    INDEX idx_intent_type (intent_type),
    INDEX idx_language (language),
    INDEX idx_source (source),
    INDEX idx_quality_score (quality_score),
    INDEX idx_is_active (is_active),
    FULLTEXT INDEX ft_sample_text (sample_text)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='意图样本表';

-- 7. 意图模型表 (intent_models)
CREATE TABLE IF NOT EXISTS intent_models (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    model_id VARCHAR(64) NOT NULL COMMENT '模型唯一标识',
    model_name VARCHAR(100) NOT NULL COMMENT '模型名称',
    model_type VARCHAR(50) NOT NULL COMMENT '模型类型',
    model_config JSON NOT NULL COMMENT '模型配置',

    -- 训练数据
    trained_samples INT DEFAULT 0 COMMENT '训练样本数',
    trained_at TIMESTAMP NULL COMMENT '训练时间',

    -- 性能指标
    accuracy DECIMAL(5,4) NULL COMMENT '准确率',
    precision DECIMAL(5,4) NULL COMMENT '精确率',
    recall DECIMAL(5,4) NULL COMMENT '召回率',
    f1_score DECIMAL(5,4) NULL COMMENT 'F1分数',

    -- 部署信息
    is_deployed BOOLEAN DEFAULT FALSE COMMENT '是否已部署',
    deployed_at TIMESTAMP NULL COMMENT '部署时间',
    deployment_env VARCHAR(20) NULL COMMENT '部署环境',

    -- 版本信息
    version INT DEFAULT 1 COMMENT '版本号',
    parent_model_id VARCHAR(64) NULL COMMENT '父模型ID',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_model_id (model_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_model_type (model_type),
    INDEX idx_is_deployed (is_deployed),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='意图模型表';

-- 8. 路由日志表 (routing_logs)
CREATE TABLE IF NOT EXISTS routing_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    request_id VARCHAR(64) NOT NULL COMMENT '请求ID',
    session_id VARCHAR(64) NOT NULL COMMENT '会话ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',

    -- 请求信息
    request_text TEXT NOT NULL COMMENT '请求文本',
    request_metadata JSON NULL COMMENT '请求元数据',

    -- 处理信息
    processing_stage VARCHAR(50) NOT NULL COMMENT '处理阶段',
    processing_status VARCHAR(20) NOT NULL COMMENT '处理状态',
    processing_time_ms INT NOT NULL COMMENT '处理耗时',

    -- 结果信息
    intent_type VARCHAR(50) NULL COMMENT '识别的意图',
    confidence DECIMAL(5,4) NULL COMMENT '置信度',
    selected_service VARCHAR(64) NULL COMMENT '选中的服务',

    -- 错误信息
    error_code VARCHAR(50) NULL COMMENT '错误码',
    error_message TEXT NULL COMMENT '错误消息',

    -- 调试信息
    debug_info JSON NULL COMMENT '调试信息',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_request_id (request_id),
    INDEX idx_session_id (session_id),
    INDEX idx_user_id (user_id),
    INDEX idx_processing_status (processing_status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由日志表';

-- 9. 路由指标表 (routing_metrics)
CREATE TABLE IF NOT EXISTS routing_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    service_id VARCHAR(64) NOT NULL COMMENT '服务ID',
    intent_type VARCHAR(50) NOT NULL COMMENT '意图类型',

    -- 时间窗口
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口类型',
    window_start TIMESTAMP NOT NULL COMMENT '窗口开始',
    window_end TIMESTAMP NOT NULL COMMENT '窗口结束',

    -- 请求统计
    total_requests BIGINT DEFAULT 0 COMMENT '总请求数',
    successful_requests BIGINT DEFAULT 0 COMMENT '成功请求数',
    failed_requests BIGINT DEFAULT 0 COMMENT '失败请求数',

    -- 性能指标
    avg_response_time_ms INT DEFAULT 0 COMMENT '平均响应时间',
    p95_response_time_ms INT DEFAULT 0 COMMENT 'P95响应时间',
    p99_response_time_ms INT DEFAULT 0 COMMENT 'P99响应时间',

    -- 准确性指标
    accuracy DECIMAL(5,4) DEFAULT 0.0000 COMMENT '准确率',
    user_satisfaction DECIMAL(5,4) DEFAULT 0.0000 COMMENT '用户满意度',

    -- 负载指标
    current_load INT DEFAULT 0 COMMENT '当前负载',
    max_load INT DEFAULT 0 COMMENT '峰值负载',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_service_intent_time (service_id, intent_type, time_window, window_start),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_time_window (time_window, window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由指标表';

-- 10. 告警规则表 (alert_rules)
CREATE TABLE IF NOT EXISTS alert_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    rule_id VARCHAR(64) NOT NULL COMMENT '规则唯一标识',
    rule_name VARCHAR(100) NOT NULL COMMENT '规则名称',
    description TEXT NULL COMMENT '规则描述',

    -- 监控目标
    target_type VARCHAR(50) NOT NULL COMMENT '目标类型',
    target_id VARCHAR(64) NOT NULL COMMENT '目标ID',

    -- 触发条件
    metric_name VARCHAR(50) NOT NULL COMMENT '指标名称',
    condition_operator VARCHAR(10) NOT NULL COMMENT '条件操作符',
    threshold_value DECIMAL(15,4) NOT NULL COMMENT '阈值',
    evaluation_window INT NOT NULL COMMENT '评估窗口(秒)',

    -- 告警配置
    severity VARCHAR(20) DEFAULT 'warning' COMMENT '严重程度',
    notification_channels JSON NOT NULL COMMENT '通知渠道',

    -- 冷却配置
    cooldown_period INT DEFAULT 300 COMMENT '冷却期(秒)',

    -- 状态
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    last_triggered_at TIMESTAMP NULL COMMENT '最后触发时间',
    trigger_count INT DEFAULT 0 COMMENT '触发次数',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_rule_id (rule_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_target (target_type, target_id),
    INDEX idx_is_enabled (is_enabled),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';

-- 11. 告警历史表 (alert_history)
CREATE TABLE IF NOT EXISTS alert_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    alert_id VARCHAR(64) NOT NULL COMMENT '告警唯一标识',
    rule_id VARCHAR(64) NOT NULL COMMENT '触发规则ID',

    -- 告警内容
    alert_title VARCHAR(255) NOT NULL COMMENT '告警标题',
    alert_message TEXT NOT NULL COMMENT '告警消息',
    severity VARCHAR(20) NOT NULL COMMENT '严重程度',

    -- 触发信息
    metric_name VARCHAR(50) NOT NULL COMMENT '指标名称',
    actual_value DECIMAL(15,4) NOT NULL COMMENT '实际值',
    threshold_value DECIMAL(15,4) NOT NULL COMMENT '阈值',

    -- 状态
    status VARCHAR(20) DEFAULT 'firing' COMMENT '告警状态',
    acknowledged_at TIMESTAMP NULL COMMENT '确认时间',
    acknowledged_by BIGINT NULL COMMENT '确认人',
    resolved_at TIMESTAMP NULL COMMENT '解决时间',

    -- 通知信息
    notifications_sent JSON NULL COMMENT '发送的通知',
    notification_status JSON NULL COMMENT '通知状态',

    fired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_alert_id (alert_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_rule_id (rule_id),
    INDEX idx_status (status),
    INDEX idx_fired_at (fired_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警历史表';

-- 12. A/B测试表 (ab_tests)
CREATE TABLE IF NOT EXISTS ab_tests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    test_id VARCHAR(64) NOT NULL COMMENT '测试唯一标识',
    test_name VARCHAR(100) NOT NULL COMMENT '测试名称',
    description TEXT NULL COMMENT '测试描述',

    -- 测试配置
    target_type VARCHAR(50) NOT NULL COMMENT '测试目标类型',
    target_id VARCHAR(64) NOT NULL COMMENT '测试目标ID',
    test_config JSON NOT NULL COMMENT '测试配置',
    variants JSON NOT NULL COMMENT '测试变体',

    -- 流量分配
    traffic_allocation JSON NOT NULL COMMENT '流量分配比例',
    total_traffic BIGINT DEFAULT 0 COMMENT '总流量',

    -- 时间配置
    start_time TIMESTAMP NOT NULL COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',

    -- 状态
    status VARCHAR(20) DEFAULT 'draft' COMMENT '测试状态',
    winning_variant VARCHAR(50) NULL COMMENT '获胜变体',

    -- 统计结果
    results JSON NULL COMMENT '测试结果',
    statistical_significance DECIMAL(5,4) NULL COMMENT '统计显著性',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_test_id (test_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_status (status),
    INDEX idx_start_end_time (start_time, end_time),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='A/B测试表';

-- 13. 模型优化记录表 (model_optimization_records)
CREATE TABLE IF NOT EXISTS model_optimization_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    model_id VARCHAR(64) NOT NULL COMMENT '模型ID',
    optimization_id VARCHAR(64) NOT NULL COMMENT '优化唯一标识',

    -- 优化配置
    optimization_type VARCHAR(50) NOT NULL COMMENT '优化类型',
    optimization_config JSON NOT NULL COMMENT '优化配置',
    base_model_id VARCHAR(64) NULL COMMENT '基础模型ID',

    -- 优化结果
    optimization_status VARCHAR(20) NOT NULL COMMENT '优化状态',
    optimized_model_id VARCHAR(64) NULL COMMENT '优化后的模型ID',
    performance_improvement DECIMAL(5,4) NULL COMMENT '性能提升',

    -- 训练信息
    training_samples INT DEFAULT 0 COMMENT '训练样本数',
    training_time_ms INT DEFAULT 0 COMMENT '训练耗时',
    training_cost DECIMAL(10,4) DEFAULT 0.0000 COMMENT '训练成本',

    -- 评估结果
    evaluation_metrics JSON NULL COMMENT '评估指标',
    comparison_with_baseline JSON NULL COMMENT '与基线对比',

    started_at TIMESTAMP NOT NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_optimization_id (optimization_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_model_id (model_id),
    INDEX idx_optimization_status (optimization_status),
    INDEX idx_started_at (started_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型优化记录表';

-- 14. 意图训练历史表 (intent_training_history)
CREATE TABLE IF NOT EXISTS intent_training_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    training_id VARCHAR(64) NOT NULL COMMENT '训练唯一标识',
    model_id VARCHAR(64) NOT NULL COMMENT '模型ID',
    intent_type VARCHAR(50) NOT NULL COMMENT '意图类型',

    -- 训练配置
    training_type VARCHAR(50) NOT NULL COMMENT '训练类型',
    training_config JSON NOT NULL COMMENT '训练配置',
    training_samples INT NOT NULL COMMENT '训练样本数',

    -- 训练结果
    training_status VARCHAR(20) NOT NULL COMMENT '训练状态',
    training_progress DECIMAL(5,2) DEFAULT 0.00 COMMENT '训练进度',
    training_loss DECIMAL(10,4) NULL COMMENT '训练损失',
    validation_accuracy DECIMAL(5,4) NULL COMMENT '验证准确率',

    -- 时间信息
    started_at TIMESTAMP NOT NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    training_duration_ms INT NULL COMMENT '训练耗时',

    -- 资源使用
    gpu_hours DECIMAL(10,4) DEFAULT 0.0000 COMMENT 'GPU小时数',
    cpu_hours DECIMAL(10,4) DEFAULT 0.0000 COMMENT 'CPU小时数',
    training_cost DECIMAL(10,4) DEFAULT 0.0000 COMMENT '训练成本',

    -- 错误信息
    error_message TEXT NULL COMMENT '错误消息',
    error_details JSON NULL COMMENT '错误详情',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_training_id (training_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_model_id (model_id),
    INDEX idx_intent_type (intent_type),
    INDEX idx_training_status (training_status),
    INDEX idx_started_at (started_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='意图训练历史表';

-- 15. 服务健康表 (service_health)
CREATE TABLE IF NOT EXISTS service_health (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    service_id VARCHAR(64) NOT NULL COMMENT '服务ID',
    check_time TIMESTAMP NOT NULL COMMENT '检查时间',

    -- 健康状态
    health_status VARCHAR(20) NOT NULL COMMENT '健康状态',
    status_reason VARCHAR(255) NULL COMMENT '状态原因',

    -- 负载指标
    current_load INT DEFAULT 0 COMMENT '当前负载',
    current_connections INT DEFAULT 0 COMMENT '当前连接数',
    max_capacity INT DEFAULT 0 COMMENT '最大容量',

    -- 性能指标
    avg_response_time_ms INT DEFAULT 0 COMMENT '平均响应时间',
    p95_response_time_ms INT DEFAULT 0 COMMENT 'P95响应时间',
    p99_response_time_ms INT DEFAULT 0 COMMENT 'P99响应时间',

    -- 错误率
    error_rate DECIMAL(5,4) DEFAULT 0.0000 COMMENT '错误率',
    timeout_rate DECIMAL(5,4) DEFAULT 0.0000 COMMENT '超时率',

    -- 资源使用
    cpu_usage DECIMAL(5,2) DEFAULT 0.00 COMMENT 'CPU使用率',
    memory_usage DECIMAL(5,2) DEFAULT 0.00 COMMENT '内存使用率',
    disk_usage DECIMAL(5,2) DEFAULT 0.00 COMMENT '磁盘使用率',

    -- 自定义指标
    custom_metrics JSON NULL COMMENT '自定义指标',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    INDEX idx_tenant_id (tenant_id),
    INDEX idx_service_id (service_id),
    INDEX idx_check_time (check_time),
    INDEX idx_health_status (health_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务健康表';

-- 16. 负载统计表 (load_statistics)
CREATE TABLE IF NOT EXISTS load_statistics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    service_id VARCHAR(64) NOT NULL COMMENT '服务ID',

    -- 时间窗口
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口类型',
    window_start TIMESTAMP NOT NULL COMMENT '窗口开始',
    window_end TIMESTAMP NOT NULL COMMENT '窗口结束',

    -- 请求统计
    total_requests BIGINT DEFAULT 0 COMMENT '总请求数',
    successful_requests BIGINT DEFAULT 0 COMMENT '成功请求数',
    failed_requests BIGINT DEFAULT 0 COMMENT '失败请求数',

    -- 负载统计
    avg_load DECIMAL(10,2) DEFAULT 0.00 COMMENT '平均负载',
    peak_load INT DEFAULT 0 COMMENT '峰值负载',
    min_load INT DEFAULT 0 COMMENT '最小负载',

    -- 响应时间
    avg_response_time_ms INT DEFAULT 0 COMMENT '平均响应时间',
    min_response_time_ms INT DEFAULT 0 COMMENT '最小响应时间',
    max_response_time_ms INT DEFAULT 0 COMMENT '最大响应时间',

    -- 队列统计
    avg_queue_length DECIMAL(10,2) DEFAULT 0.00 COMMENT '平均队列长度',
    peak_queue_length INT DEFAULT 0 COMMENT '峰值队列长度',

    -- 并发统计
    avg_concurrent INT DEFAULT 0 COMMENT '平均并发数',
    peak_concurrent INT DEFAULT 0 COMMENT '峰值并发数',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    UNIQUE KEY uk_service_time_window (service_id, time_window, window_start),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_time_window (time_window, window_start, window_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='负载统计表';

-- 17. 容量限制表 (capacity_limits)
CREATE TABLE IF NOT EXISTS capacity_limits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    service_id VARCHAR(64) NOT NULL COMMENT '服务ID',

    -- 容量配置
    limit_type VARCHAR(50) NOT NULL COMMENT '限制类型',
    limit_name VARCHAR(100) NOT NULL COMMENT '限制名称',
    limit_value BIGINT NOT NULL COMMENT '限制值',
    limit_unit VARCHAR(20) NOT NULL COMMENT '限制单位',

    -- 时间配置
    time_window VARCHAR(20) NOT NULL COMMENT '时间窗口',
    reset_strategy VARCHAR(20) NOT NULL COMMENT '重置策略',

    -- 状态
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    current_usage BIGINT DEFAULT 0 COMMENT '当前使用量',

    -- 告警配置
    warning_threshold INT DEFAULT 80 COMMENT '告警阈值',
    critical_threshold INT DEFAULT 90 COMMENT '严重阈值',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    UNIQUE KEY uk_service_limit_type (service_id, limit_type),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_is_enabled (is_enabled),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='容量限制表';

-- 18. 路由规则版本表 (routing_rule_versions)
CREATE TABLE IF NOT EXISTS routing_rule_versions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键ID',
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    rule_id BIGINT NOT NULL COMMENT '规则ID',
    version_id VARCHAR(64) NOT NULL COMMENT '版本唯一标识',
    version_number INT NOT NULL COMMENT '版本号',

    -- 版本内容
    rule_name VARCHAR(100) NOT NULL COMMENT '规则名称',
    rule_description TEXT NULL COMMENT '规则描述',
    condition JSON NOT NULL COMMENT '触发条件',
    action JSON NOT NULL COMMENT '路由动作',
    priority INT NOT NULL COMMENT '优先级',

    -- 变更信息
    change_type VARCHAR(20) NOT NULL COMMENT '变更类型',
    change_description TEXT NULL COMMENT '变更描述',
    change_summary JSON NULL COMMENT '变更摘要',

    -- 测试信息
    test_coverage DECIMAL(5,2) DEFAULT 0.00 COMMENT '测试覆盖率',
    test_results JSON NULL COMMENT '测试结果',

    -- 发布信息
    is_published BOOLEAN DEFAULT FALSE COMMENT '是否已发布',
    published_at TIMESTAMP NULL COMMENT '发布时间',
    published_by BIGINT NULL COMMENT '发布人',

    -- 回滚信息
    is_rollback BOOLEAN DEFAULT FALSE COMMENT '是否为回滚版本',
    rollback_from_version_id VARCHAR(64) NULL COMMENT '从哪个版本回滚',

    created_by BIGINT NULL COMMENT '创建人',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',

    UNIQUE KEY uk_version_id (version_id),
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_rule_id (rule_id),
    INDEX idx_version_number (rule_id, version_number),
    INDEX idx_is_published (is_published)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='路由规则版本表';

SET FOREIGN_KEY_CHECKS = 1;
