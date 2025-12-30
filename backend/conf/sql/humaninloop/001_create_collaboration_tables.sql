-- =====================================================
-- 人机协同引擎 - 数据库表创建脚本
-- 版本: v1.0.0
-- 创建日期: 2025-01-03
-- =====================================================

-- =====================================================
-- 1. 协同任务表 (collaboration_tasks)
-- =====================================================
CREATE TABLE IF NOT EXISTS `collaboration_tasks` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID（UUID）',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 任务类型和来源
    `task_type` ENUM('REVIEW', 'CORRECTION', 'VALIDATION', 'ESCALATION') NOT NULL COMMENT '任务类型',
    `source` ENUM('AI', 'HUMAN', 'SYSTEM') NOT NULL COMMENT '任务来源',
    `priority` ENUM('LOW', 'MEDIUM', 'HIGH', 'URGENT') NOT NULL DEFAULT 'MEDIUM' COMMENT '任务优先级',
    `status` ENUM('PENDING', 'ASSIGNED', 'IN_PROGRESS', 'COMPLETED', 'ESCALATED', 'CANCELLED')
        NOT NULL DEFAULT 'PENDING' COMMENT '任务状态',

    -- 分配信息
    `assigned_to` VARCHAR(36) NULL COMMENT '分配给的审核人ID',

    -- SLA信息
    `sla_deadline` TIMESTAMP NOT NULL COMMENT 'SLA截止时间',

    -- 上下文和结果（JSON格式）
    `context` JSON NULL COMMENT '任务上下文',
    `result` JSON NULL COMMENT '审核结果',

    -- 关联信息
    `conversation_id` VARCHAR(36) NULL COMMENT '关联的会话ID',
    `message_id` VARCHAR(36) NULL COMMENT '关联的消息ID',
    `bot_id` VARCHAR(36) NULL COMMENT '关联的Bot ID',
    `workflow_id` VARCHAR(36) NULL COMMENT '关联的工作流ID',

    -- 时间戳
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',
    `completed_at` BIGINT NULL COMMENT '完成时间（毫秒时间戳）',

    -- 索引
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_task_id` (`task_id`),
    INDEX `idx_tenant_status` (`tenant_id`, `status`),
    INDEX `idx_assigned_to` (`assigned_to`),
    INDEX `idx_sla` (`sla_deadline`, `status`),
    INDEX `idx_priority` (`priority`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='协同任务表';

-- =====================================================
-- 2. 协同历史表 (collaboration_history)
-- =====================================================
CREATE TABLE IF NOT EXISTS `collaboration_history` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `task_id` VARCHAR(36) NOT NULL COMMENT '任务ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 操作者信息
    `actor_type` ENUM('AI', 'HUMAN', 'SYSTEM') NOT NULL COMMENT '操作者类型',
    `actor_id` VARCHAR(36) NOT NULL COMMENT '操作者ID',

    -- 操作信息
    `action` VARCHAR(64) NOT NULL COMMENT '操作类型',
    `details` JSON NULL COMMENT '操作详情',

    -- 时间戳
    `timestamp` BIGINT NOT NULL DEFAULT 0 COMMENT '操作时间（毫秒时间戳）',

    -- 索引
    PRIMARY KEY (`id`),
    INDEX `idx_task_id` (`task_id`),
    INDEX `idx_timestamp` (`timestamp`),
    INDEX `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='协同历史表';

-- =====================================================
-- 3. 协同配置表 (collaboration_configs)
-- =====================================================
CREATE TABLE IF NOT EXISTS `collaboration_configs` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `tenant_id` VARCHAR(36) NOT NULL COMMENT '租户ID',

    -- 自动审核阈值
    `auto_review_threshold` DECIMAL(3,2) NOT NULL DEFAULT 0.70 COMMENT '自动审核阈值',
    `auto_escalate_threshold` DECIMAL(3,2) NOT NULL DEFAULT 0.30 COMMENT '自动升级阈值',

    -- SLA配置
    `default_sla_minutes` INT NOT NULL DEFAULT 60 COMMENT '默认SLA（分钟）',
    `urgent_sla_minutes` INT NOT NULL DEFAULT 15 COMMENT '紧急任务SLA（分钟）',

    -- 审核人员池配置（JSON格式）
    `review_pool_filter` JSON NULL COMMENT '审核人员池过滤规则',
    `escalation_rules` JSON NULL COMMENT '升级规则',
    `priority_rules` JSON NULL COMMENT '优先级规则',

    -- 通知配置
    `notification_enabled` BOOLEAN NOT NULL DEFAULT TRUE COMMENT '是否启用通知',
    `notification_channels` JSON NULL COMMENT '通知渠道配置',

    -- 时间戳
    `created_at` BIGINT NOT NULL DEFAULT 0 COMMENT '创建时间（毫秒时间戳）',
    `updated_at` BIGINT NOT NULL DEFAULT 0 COMMENT '更新时间（毫秒时间戳）',

    -- 索引
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='协同配置表';

-- =====================================================
-- 4. 插入默认配置（平台级）
-- =====================================================

-- 注意：这里不插入平台级配置，因为表设计中没有 tenant_id = NULL 的记录
-- 租户配置会在租户创建时自动生成（使用默认值）

-- =====================================================
-- 5. 创建视图（可选）
-- =====================================================

-- 待审核任务视图
CREATE OR REPLACE VIEW `v_pending_tasks` AS
SELECT
    ct.task_id,
    ct.tenant_id,
    ct.task_type,
    ct.priority,
    ct.created_at,
    ct.sla_deadline,
    TIMESTAMPDIFF(MINUTE, NOW(), ct.sla_deadline) as remaining_minutes,
    CASE
        WHEN ct.sla_deadline < NOW() THEN 'overdue'
        WHEN TIMESTAMPDIFF(MINUTE, NOW(), ct.sla_deadline) < 15 THEN 'urgent'
        WHEN TIMESTAMPDIFF(MINUTE, NOW(), ct.sla_deadline) < 60 THEN 'warning'
        ELSE 'normal'
    END as sla_status
FROM collaboration_tasks ct
WHERE ct.status = 'PENDING';

-- 审核队列视图
CREATE OR REPLACE VIEW `v_review_queue` AS
SELECT
    ct.task_id,
    ct.tenant_id,
    ct.task_type,
    ct.priority,
    ct.status,
    ct.assigned_to,
    ct.created_at,
    ct.sla_deadline,
    (SELECT COUNT(*) FROM collaboration_history WHERE task_id = ct.task_id) as history_count
FROM collaboration_tasks ct
WHERE ct.status IN ('PENDING', 'ASSIGNED', 'IN_PROGRESS')
ORDER BY
    FIELD(ct.priority, 'URGENT', 'HIGH', 'MEDIUM', 'LOW'),
    ct.sla_deadline ASC;

-- =====================================================
-- 6. 创建存储过程（可选）
-- =====================================================

DELIMITER $$

-- 自动调整SLA优先级的存储过程
CREATE PROCEDURE IF NOT EXISTS `sp_adjust_priority_by_sla`()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE v_task_id VARCHAR(36);
    DECLARE v_priority VARCHAR(10);
    DECLARE v_remaining_minutes INT;

    DECLARE task_cursor CURSOR FOR
        SELECT task_id, priority, TIMESTAMPDIFF(MINUTE, NOW(), sla_deadline)
        FROM collaboration_tasks
        WHERE status IN ('PENDING', 'ASSIGNED', 'IN_PROGRESS')
          AND sla_deadline > NOW();

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    OPEN task_cursor;

    read_loop: LOOP
        FETCH task_cursor INTO v_task_id, v_priority, v_remaining_minutes;
        IF done THEN
            LEAVE read_loop;
        END IF;

        -- SLA剩余时间少于15分钟，提升到紧急
        IF v_remaining_minutes < 15 AND v_priority != 'URGENT' THEN
            UPDATE collaboration_tasks
            SET priority = 'URGENT', updated_at = UNIX_TIMESTAMP(NOW()) * 1000
            WHERE task_id = v_task_id;
        -- SLA剩余时间少于1小时，提升到高
        ELSEIF v_remaining_minutes < 60 AND v_priority = 'LOW' THEN
            UPDATE collaboration_tasks
            SET priority = 'HIGH', updated_at = UNIX_TIMESTAMP(NOW()) * 1000
            WHERE task_id = v_task_id;
        END IF;
    END LOOP;

    CLOSE task_cursor;
END$$

DELIMITER ;

-- =====================================================
-- 7. 创建触发器（可选）
-- =====================================================

DELIMITER $$

-- 更新时间戳触发器
CREATE TRIGGER IF NOT EXISTS `tr_collaboration_tasks_update`
BEFORE UPDATE ON `collaboration_tasks`
FOR EACH ROW
BEGIN
    SET NEW.updated_at = UNIX_TIMESTAMP(NOW()) * 1000;

    -- 当状态变更为完成时，设置完成时间
    IF NEW.status = 'COMPLETED' AND OLD.status != 'COMPLETED' THEN
        SET NEW.completed_at = UNIX_TIMESTAMP(NOW()) * 1000;
    END IF;
END$$

DELIMITER ;

-- =====================================================
-- 8. 性能优化建议
-- =====================================================

-- 对于大量查询，可以考虑添加以下复合索引：
-- ALTER TABLE collaboration_tasks ADD INDEX idx_tenant_priority_status (tenant_id, priority, status);
-- ALTER TABLE collaboration_tasks ADD INDEX idx_assigned_status (assigned_to, status);
-- ALTER TABLE collaboration_history ADD INDEX idx_task_action (task_id, action);

-- =====================================================
-- 结束
-- =====================================================
