-- Saga分布式事务框架表结构
-- 版本: v1.0
-- 创建日期: 2025-12-30
-- 说明: 支持Saga编排式分布式事务模式

-- ============================================
-- 表1: Saga定义表 (saga_definitions)
-- 说明: 存储Saga流程定义
-- ============================================
CREATE TABLE IF NOT EXISTS saga_definitions (
    -- 主键
    id VARCHAR(64) PRIMARY KEY COMMENT 'Saga唯一标识',

    -- 基本信息
    name VARCHAR(128) NOT NULL UNIQUE COMMENT 'Saga名称',
    description TEXT COMMENT 'Saga描述',

    -- 步骤定义(JSON格式,存储SagaStep接口的实现)
    steps TEXT NOT NULL COMMENT '执行步骤定义(JSON)',

    -- 补偿步骤定义(JSON格式,存储CompensationStep接口的实现)
    compensations TEXT NOT NULL COMMENT '补偿步骤定义(JSON)',

    -- 重试策略(JSON格式)
    retry_policy TEXT COMMENT '重试策略配置(JSON)',

    -- 超时配置
    timeout_seconds BIGINT NOT NULL DEFAULT 300 COMMENT '超时时间(秒)',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 索引
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Saga定义表';

-- ============================================
-- 表2: Saga执行记录表 (saga_executions)
-- 说明: 存储Saga执行记录
-- ============================================
CREATE TABLE IF NOT EXISTS saga_executions (
    -- 主键
    id VARCHAR(64) PRIMARY KEY COMMENT '执行记录ID',

    -- 关联信息
    saga_id VARCHAR(64) NOT NULL COMMENT 'Saga定义ID',

    -- 执行状态
    status VARCHAR(32) NOT NULL COMMENT '执行状态: pending/running/completed/failed/compensating/compensated',
    current_step INT NOT NULL DEFAULT 0 COMMENT '当前步骤索引',

    -- 数据(JSON格式)
    input_data TEXT COMMENT '输入数据(JSON)',
    output_data TEXT COMMENT '输出数据(JSON)',

    -- 错误信息
    error TEXT COMMENT '错误信息',

    -- 步骤执行记录(JSON格式)
    step_executions TEXT COMMENT '步骤执行记录列表(JSON)',

    -- 时间戳
    started_at TIMESTAMP NOT NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

    -- 外键约束
    FOREIGN KEY (saga_id) REFERENCES saga_definitions(id) ON DELETE CASCADE,

    -- 索引
    INDEX idx_saga_id (saga_id),
    INDEX idx_status (status),
    INDEX idx_started_at (started_at),
    INDEX idx_saga_status (saga_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
COMMENT='Saga执行记录表';

-- ============================================
-- 初始化数据
-- ============================================

-- 示例: 租户注册Saga定义(实际使用时通过代码注册)
-- INSERT INTO saga_definitions (id, name, description, steps, compensations, timeout_seconds)
-- VALUES (
--     'saga-tenant-registration',
--     '租户注册Saga',
--     '处理企业租户注册流程',
--     '[]',
--     '[]',
--     300
-- );
