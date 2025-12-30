-- ============================================================
-- 计费系统数据库表结构 (MySQL 8.4.5)
-- 基于设计文档: 23-MultiTenant SaaS核心_租户计费系统.md
-- ============================================================

-- ============================================================
-- 1. 计费账户表 (billing_accounts)
-- ============================================================

CREATE TABLE IF NOT EXISTS billing_accounts (
    -- 主键
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '账户ID',

    -- 租户信息
    tenant_id VARCHAR(64) NOT NULL UNIQUE COMMENT '租户ID',

    -- 账户状态
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '账户状态: active, suspended, closed',

    -- 信用额度
    credit_limit DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '信用额度',
    available_credit DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '可用信用额度',
    currency VARCHAR(3) NOT NULL DEFAULT 'CNY' COMMENT '货币',
    auto_recharge BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否自动充值',
    recharge_threshold DECIMAL(12,2) NULL COMMENT '自动充值阈值',
    recharge_amount DECIMAL(12,2) NULL COMMENT '自动充值金额',

    -- 账单周期
    billing_cycle VARCHAR(20) NOT NULL DEFAULT 'monthly' COMMENT '计费周期: monthly, quarterly, yearly',
    billing_day_of_month TINYINT NULL COMMENT '每月账单日(1-28)',

    -- 支付方式
    payment_methods JSON NULL COMMENT '支付方式列表',
    default_payment_method VARCHAR(64) NULL COMMENT '默认支付方式ID',

    -- 备注
    notes TEXT NULL COMMENT '备注',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    -- 索引
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='计费账户表';


-- ============================================================
-- 2. 发票表 (invoices)
-- ============================================================

CREATE TABLE IF NOT EXISTS invoices (
    -- 主键
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '发票ID',

    -- 租户信息
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    billing_account_id BIGINT UNSIGNED NOT NULL COMMENT '计费账户ID',

    -- 发票编号
    invoice_number VARCHAR(64) NOT NULL UNIQUE COMMENT '发票编号',
    invoice_type VARCHAR(20) NOT NULL COMMENT '发票类型: subscription, overage, one_time, credit',

    -- 账单周期
    period_start TIMESTAMP NOT NULL COMMENT '账期开始',
    period_end TIMESTAMP NOT NULL COMMENT '账期结束',

    -- 金额明细
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '小计金额',
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '税额',
    discount_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '折扣金额',
    total_amount DECIMAL(12,2) NOT NULL COMMENT '总金额',
    currency VARCHAR(3) NOT NULL DEFAULT 'CNY' COMMENT '货币',

    -- 费用明细
    line_items JSON NULL COMMENT '费用明细',

    -- 发票状态
    status VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '发票状态: draft, sent, viewed, paid, partial_paid, overdue, cancelled, void',

    -- 支付信息
    due_date TIMESTAMP NULL COMMENT '到期日期',
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 COMMENT '已支付金额',
    paid_at TIMESTAMP NULL COMMENT '支付时间',

    -- 发票PDF
    pdf_url TEXT NULL COMMENT '发票PDF地址',
    pdf_generated_at TIMESTAMP NULL COMMENT 'PDF生成时间',

    -- 备注
    metadata JSON NULL COMMENT '元数据',
    notes TEXT NULL COMMENT '备注',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    -- 外键约束
    CONSTRAINT fk_invoices_billing_account FOREIGN KEY (billing_account_id)
        REFERENCES billing_accounts(id) ON DELETE RESTRICT ON UPDATE CASCADE,

    -- 索引
    INDEX idx_tenant_status (tenant_id, status),
    INDEX idx_tenant_period (tenant_id, period_start, period_end),
    INDEX idx_account (billing_account_id),
    INDEX idx_status_period (status, due_date),
    INDEX idx_amount (total_amount),
    INDEX idx_period_start (period_start),
    INDEX idx_period_end (period_end),
    INDEX idx_due_date (due_date),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发票表';


-- ============================================================
-- 3. 支付记录表 (payments)
-- ============================================================

CREATE TABLE IF NOT EXISTS payments (
    -- 主键
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '支付ID',

    -- 租户信息
    tenant_id VARCHAR(64) NOT NULL COMMENT '租户ID',
    billing_account_id BIGINT UNSIGNED NOT NULL COMMENT '计费账户ID',
    invoice_id BIGINT UNSIGNED NULL COMMENT '关联发票ID',

    -- 支付单号
    payment_number VARCHAR(64) NOT NULL UNIQUE COMMENT '支付单号',
    transaction_id VARCHAR(128) NULL COMMENT '第三方交易ID',

    -- 支付金额
    amount DECIMAL(12,2) NOT NULL COMMENT '支付金额',
    currency VARCHAR(3) NOT NULL DEFAULT 'CNY' COMMENT '货币',

    -- 支付方式
    payment_method VARCHAR(50) NOT NULL COMMENT '支付方式: alipay, wechat, credit_card, bank_transfer, paypal',
    payment_method_id VARCHAR(64) NULL COMMENT '支付方式ID',

    -- 支付状态
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '支付状态: pending, processing, success, failed, cancelled, refunded',

    -- 第三方信息
    provider VARCHAR(50) NOT NULL COMMENT '支付提供商: alipay, wechat, stripe, paypal',
    provider_response TEXT NULL COMMENT '第三方响应',

    -- 退款信息
    refund_amount DECIMAL(12,2) NULL COMMENT '退款金额',
    refunded_at TIMESTAMP NULL COMMENT '退款时间',
    refund_reason TEXT NULL COMMENT '退款原因',

    -- 失败原因
    failure_reason TEXT NULL COMMENT '失败原因',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',
    notes TEXT NULL COMMENT '备注',

    -- 时间戳
    processed_at TIMESTAMP NULL COMMENT '处理时间',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    -- 外键约束
    CONSTRAINT fk_payments_billing_account FOREIGN KEY (billing_account_id)
        REFERENCES billing_accounts(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_payments_invoice FOREIGN KEY (invoice_id)
        REFERENCES invoices(id) ON DELETE SET NULL ON UPDATE CASCADE,

    -- 索引
    INDEX idx_tenant_created (tenant_id, created_at),
    INDEX idx_account (billing_account_id),
    INDEX idx_invoice (invoice_id),
    INDEX idx_amount (amount),
    INDEX idx_status (status),
    INDEX idx_method (payment_method),
    INDEX idx_transaction (transaction_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付记录表';


-- ============================================================
-- 4. 发票明细表 (invoice_line_items)
-- ============================================================

CREATE TABLE IF NOT EXISTS invoice_line_items (
    -- 主键
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '明细ID',
    invoice_id BIGINT UNSIGNED NOT NULL COMMENT '发票ID',

    -- 明细信息
    description TEXT NOT NULL COMMENT '描述',
    quantity INT NOT NULL DEFAULT 1 COMMENT '数量',
    unit_price DECIMAL(12,2) NOT NULL COMMENT '单价',
    amount DECIMAL(12,2) NOT NULL COMMENT '金额',

    -- 分类
    item_type VARCHAR(50) NOT NULL COMMENT '项目类型: subscription, overage_tokens, overage_storage, setup_fee, discount',
    item_code VARCHAR(50) NOT NULL COMMENT '项目代码',

    -- 元数据
    metadata JSON NULL COMMENT '元数据',

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',

    -- 外键约束
    CONSTRAINT fk_line_items_invoice FOREIGN KEY (invoice_id)
        REFERENCES invoices(id) ON DELETE CASCADE ON UPDATE CASCADE,

    -- 索引
    INDEX idx_invoice (invoice_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='发票明细表';


-- ============================================================
-- 初始化数据
-- ============================================================

-- 插入示例计费账户（可选）
-- INSERT INTO billing_accounts (tenant_id, status, credit_limit, available_credit, currency, billing_cycle)
-- VALUES
--     ('tenant_001', 'active', 10000.00, 10000.00, 'CNY', 'monthly'),
--     ('tenant_002', 'active', 50000.00, 50000.00, 'CNY', 'monthly');


-- ============================================================
-- 性能优化：添加复合索引（根据查询模式）
-- ============================================================

-- 发票查询优化：租户+状态+创建时间
CREATE INDEX idx_invoices_tenant_status_created ON invoices(tenant_id, status, created_at DESC);

-- 发票查询优化：租户+账期
CREATE INDEX idx_invoices_tenant_period_range ON invoices(tenant_id, period_start, period_end);

-- 支付查询优化：租户+状态+创建时间
CREATE INDEX idx_payments_tenant_status_created ON payments(tenant_id, status, created_at DESC);


-- ============================================================
-- 数据清理策略（可选）
-- ============================================================

-- 创建软删除数据清理存储过程
DELIMITER $$

CREATE PROCEDURE CleanupOldBillingData()
BEGIN
    DECLARE deleted_count INT;

    -- 清理90天前的软删除支付记录
    DELETE FROM payments WHERE deleted_at < DATE_SUB(NOW(), INTERVAL 90 DAY);

    SET deleted_count = ROW_COUNT();
    IF deleted_count > 0 THEN
        SELECT CONCAT('Cleaned up ', deleted_count, ' old payment records') AS result;
    END IF;

    -- 清理180天前的软删除发票明细
    DELETE FROM invoice_line_items WHERE deleted_at < DATE_SUB(NOW(), INTERVAL 180 DAY);

    SET deleted_count = ROW_COUNT();
    IF deleted_count > 0 THEN
        SELECT CONCAT('Cleaned up ', deleted_count, ' old line item records') AS result;
    END IF;

    -- 注意：不自动清理invoices和billing_accounts，这些是核心业务数据
END$$

DELIMITER ;

-- 设置定时任务（需要event scheduler开启）
-- SET GLOBAL event_scheduler = ON;
-- CREATE EVENT IF NOT EXISTS cleanup_billing_data
-- ON SCHEDULE EVERY 1 DAY
-- STARTS CONCAT(CURDATE() + INTERVAL 1 DAY, ' 03:00:00')
-- DO CALL CleanupOldBillingData();


-- ============================================================
-- 视图：租户计费汇总（可选）
-- ============================================================

CREATE OR REPLACE VIEW v_tenant_billing_summary AS
SELECT
    ba.tenant_id,
    ba.status AS account_status,
    ba.credit_limit,
    ba.available_credit,
    ba.currency,
    ba.billing_cycle,
    COUNT(DISTINCT i.id) AS total_invoices,
    COALESCE(SUM(i.total_amount), 0) AS total_invoiced_amount,
    COALESCE(SUM(i.paid_amount), 0) AS total_paid_amount,
    COALESCE(SUM(p.amount), 0) AS total_payment_amount,
    COUNT(DISTINCT p.id) AS total_payments
FROM
    billing_accounts ba
LEFT JOIN
    invoices i ON ba.id = i.billing_account_id AND i.deleted_at IS NULL
LEFT JOIN
    payments p ON ba.id = p.billing_account_id AND p.deleted_at IS NULL
WHERE
    ba.deleted_at IS NULL
GROUP BY
    ba.tenant_id, ba.id;


-- ============================================================
-- 完成
-- ============================================================
