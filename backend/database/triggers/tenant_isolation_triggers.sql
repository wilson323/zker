-- ====================================================================
-- 租户隔离触发器
-- ====================================================================
--
-- **功能**：防止跨租户数据更新
--
-- **触发器列表**：
-- 1. prevent_tenant_id_update - 防止修改tenant_id字段
-- 2. tenant_isolation_check - 插入时验证tenant_id
-- 3. prevent_cross_tenant_update - 更新时验证租户一致性
--
-- **使用方法**：
-- 为每个业务表执行以下SQL：
-- ```sql
-- CALL add_tenant_isolation_triggers('table_name');
-- ```
--
-- ====================================================================

-- --------------------------------------------------------------------
-- 1. 创建存储过程：为表添加租户隔离触发器
-- --------------------------------------------------------------------

DELIMITER $$

DROP PROCEDURE IF EXISTS add_tenant_isolation_triggers$$

CREATE PROCEDURE add_tenant_isolation_triggers(
    IN table_name VARCHAR(64)
)
BEGIN
    DECLARE trigger_name VARCHAR(64);

    -- 1. 防止修改tenant_id字段（UPDATE触发器）
    SET trigger_name = CONCAT('trg_prevent_tenant_id_update_', table_name);

    SET @sql = CONCAT('
        DROP TRIGGER IF EXISTS ', trigger_name, '$$

        CREATE TRIGGER ', trigger_name, '
        BEFORE UPDATE ON ', table_name, '
        FOR EACH ROW
        BEGIN
            -- 检查是否正在修改tenant_id
            IF OLD.tenant_id IS NOT NULL AND NEW.tenant_id IS NOT NULL AND OLD.tenant_id != NEW.tenant_id THEN
                SIGNAL SQLSTATE "45000"
                SET MESSAGE_TEXT = "Cannot modify tenant_id field for security reasons";
            END IF;

            -- 如果尝试将tenant_id设置为NULL
            IF NEW.tenant_id IS NULL AND OLD.tenant_id IS NOT NULL THEN
                SIGNAL SQLSTATE "45000"
                SET MESSAGE_TEXT = "Cannot set tenant_id to NULL";
            END IF;
        END$$
    ');

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 2. 插入时验证tenant_id（INSERT触发器）
    SET trigger_name = CONCAT('trg_tenant_isolation_check_', table_name);

    SET @sql = CONCAT('
        DROP TRIGGER IF EXISTS ', trigger_name, '$$

        CREATE TRIGGER ', trigger_name, '
        BEFORE INSERT ON ', table_name, '
        FOR EACH ROW
        BEGIN
            DECLARE tenant_exists INT;

            -- 验证tenant_id是否存在
            IF NEW.tenant_id IS NOT NULL THEN
                SELECT COUNT(*) INTO tenant_exists
                FROM tenants
                WHERE tenant_id = NEW.tenant_id
                  AND status = "active"
                  AND deleted_at IS NULL
                LIMIT 1;

                IF tenant_exists = 0 THEN
                    SIGNAL SQLSTATE "45000"
                    SET MESSAGE_TEXT = "Invalid tenant_id: tenant does not exist or is not active";
                END IF;
            END IF;
        END$$
    ');

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 3. 更新时验证租户一致性（UPDATE触发器）
    SET trigger_name = CONCAT('trg_prevent_cross_tenant_update_', table_name);

    SET @sql = CONCAT('
        DROP TRIGGER IF EXISTS ', trigger_name, '$$

        CREATE TRIGGER ', trigger_name, '
        BEFORE UPDATE ON ', table_name, '
        FOR EACH ROW
        BEGIN
            -- 确保tenant_id字段不被意外修改（双重保护）
            IF OLD.tenant_id IS NOT NULL AND NEW.tenant_id IS NOT NULL AND OLD.tenant_id != NEW.tenant_id THEN
                SIGNAL SQLSTATE "45000"
                SET MESSAGE_TEXT = "Cross-tenant update is not allowed";
            END IF;
        END$$
    ');

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

END$$

DELIMITER ;

-- --------------------------------------------------------------------
-- 2. 为核心业务表添加触发器
-- --------------------------------------------------------------------

-- 用户表
CALL add_tenant_isolation_triggers('users');

-- Bot表
CALL add_tenant_isolation_triggers('bots');

-- 对话表
CALL add_tenant_isolation_triggers('conversations');

-- 消息表
CALL add_tenant_isolation_triggers('messages');

-- 知识库表
CALL add_tenant_isolation_triggers('knowledge');

-- 工作流表
CALL add_tenant_isolation_triggers('workflows');

-- 数据库表
CALL add_tenant_isolation_triggers('databases');

-- 变量表
CALL add_tenant_isolation_triggers('variables');

-- 插件表
CALL add_tenant_isolation_triggers('plugins');

-- --------------------------------------------------------------------
-- 3. 创建存储过程：删除租户隔离触发器
-- --------------------------------------------------------------------

DELIMITER $$

DROP PROCEDURE IF EXISTS remove_tenant_isolation_triggers$$

CREATE PROCEDURE remove_tenant_isolation_triggers(
    IN table_name VARCHAR(64)
)
BEGIN
    DECLARE trigger_name VARCHAR(64);

    -- 删除触发器1
    SET trigger_name = CONCAT('trg_prevent_tenant_id_update_', table_name);
    SET @sql = CONCAT('DROP TRIGGER IF EXISTS ', trigger_name);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 删除触发器2
    SET trigger_name = CONCAT('trg_tenant_isolation_check_', table_name);
    SET @sql = CONCAT('DROP TRIGGER IF EXISTS ', trigger_name);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 删除触发器3
    SET trigger_name = CONCAT('trg_prevent_cross_tenant_update_', table_name);
    SET @sql = CONCAT('DROP TRIGGER IF EXISTS ', trigger_name);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

END$$

DELIMITER ;

-- --------------------------------------------------------------------
-- 4. 验证触发器是否正确创建
-- --------------------------------------------------------------------

-- 查看所有租户隔离触发器
SELECT
    TRIGGER_NAME,
    EVENT_MANIPULATION,
    EVENT_OBJECT_TABLE,
    ACTION_STATEMENT
FROM INFORMATION_SCHEMA.TRIGGERS
WHERE TRIGGER_NAME LIKE 'trg_%_tenant%'
ORDER BY EVENT_OBJECT_TABLE, TRIGGER_NAME;

-- --------------------------------------------------------------------
-- 5. 测试用例
-- --------------------------------------------------------------------

-- 测试1：尝试修改tenant_id（应该失败）
-- UPDATE users SET tenant_id = 'other-tenant-id' WHERE user_id = 123;
-- 预期：ERROR 1644 (45000): Cannot modify tenant_id field for security reasons

-- 测试2：尝试插入无效的tenant_id（应该失败）
-- INSERT INTO users (user_id, tenant_id, username) VALUES (999, 'invalid-tenant-id', 'test');
-- 预期：ERROR 1644 (45000): Invalid tenant_id: tenant does not exist or is not active

-- 测试3：正常插入（应该成功）
-- INSERT INTO users (user_id, tenant_id, username) VALUES (999, 'valid-tenant-id', 'test');
-- 预期：成功

-- 测试4：正常更新其他字段（应该成功）
-- UPDATE users SET username = 'newname' WHERE user_id = 123;
-- 预期：成功

-- --------------------------------------------------------------------
-- 6. 性能优化建议
-- --------------------------------------------------------------------

-- 1. 为所有表的tenant_id字段添加索引
-- ALTER TABLE users ADD INDEX idx_tenant_id (tenant_id);
-- ALTER TABLE bots ADD INDEX idx_tenant_id (tenant_id);
-- ALTER TABLE conversations ADD INDEX idx_tenant_id (tenant_id);

-- 2. 使用EXPLAIN检查查询执行计划
-- EXPLAIN SELECT * FROM users WHERE tenant_id = 'xxx';

-- 3. 监控触发器性能
-- SHOW PROFILE FOR QUERY 1;

-- --------------------------------------------------------------------
-- 7. 安全审计日志
-- --------------------------------------------------------------------

-- 创建审计日志表
CREATE TABLE IF NOT EXISTS tenant_isolation_audit_log (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    table_name VARCHAR(64) NOT NULL,
    operation VARCHAR(10) NOT NULL, -- INSERT/UPDATE/DELETE
    tenant_id VARCHAR(36) NOT NULL,
    record_id VARCHAR(36) NOT NULL,
    old_tenant_id VARCHAR(36),
    new_tenant_id VARCHAR(36),
    violation_type VARCHAR(64),
    error_message TEXT,
    attempted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_table_name (table_name),
    INDEX idx_attempted_at (attempted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='租户隔离违规审计日志';

-- --------------------------------------------------------------------
-- 8. 维护命令
-- --------------------------------------------------------------------

-- 为所有业务表批量添加触发器
-- （需要根据实际表列表调整）
/*
CALL add_tenant_isolation_triggers('users');
CALL add_tenant_isolation_triggers('bots');
CALL add_tenant_isolation_triggers('conversations');
CALL add_tenant_isolation_triggers('messages');
CALL add_tenant_isolation_triggers('knowledge');
CALL add_tenant_isolation_triggers('workflows');
CALL add_tenant_isolation_triggers('databases');
CALL add_tenant_isolation_triggers('variables');
CALL add_tenant_isolation_triggers('plugins');
CALL add_tenant_isolation_triggers('api_auth');
*/

-- 查看触发器执行统计
/*
SELECT
    EVENT_OBJECT_TABLE as table_name,
    EVENT_MANIPULATION as operation,
    COUNT(*) as trigger_count
FROM INFORMATION_SCHEMA.TRIGGERS
WHERE TRIGGER_NAME LIKE 'trg_%_tenant%'
GROUP BY EVENT_OBJECT_TABLE, EVENT_MANIPULATION
ORDER BY table_name, operation;
*/
