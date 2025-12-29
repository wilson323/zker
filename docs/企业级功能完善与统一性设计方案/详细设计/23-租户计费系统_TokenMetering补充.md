# 23-租户计费系统_Token Metering补充文档

**模块**: 23-租户计费系统扩展
**扩展内容**: Token Metering、成本防护
**版本**: v2.0
**日期**: 2025-01-03

---

## 新增功能

### Token Usage Metrics 表

```sql
CREATE TABLE token_usage_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    user_id BIGINT,
    bot_id VARCHAR(64),
    conversation_id VARCHAR(64),
    message_id VARCHAR(64),

    input_tokens INT NOT NULL,
    output_tokens INT NOT NULL,
    total_tokens INT NOT NULL,

    model_provider VARCHAR(50) NOT NULL,
    model_name VARCHAR(50) NOT NULL,

    unit_price DECIMAL(10,6) NOT NULL COMMENT '每1K Token价格',
    cost DECIMAL(10,6) NOT NULL COMMENT '本次调用成本',

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 预算告警表

```sql
CREATE TABLE budget_alerts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    alert_type ENUM('threshold', 'hard_cap') NOT NULL,
    threshold_percent INT COMMENT '告警阈值百分比',
    is_sent BOOLEAN DEFAULT FALSE
);
```

### 自动降级策略

```go
type DowngradePolicy struct {
    Enabled          bool     `json:"enabled"`
    BudgetThreshold  int      `json:"budget_threshold"`  // 预算使用率阈值
    OriginalModel    string   `json:"original_model"`     // gpt-4
    FallbackModel    string   `json:"fallback_model"`     // gpt-3.5-turbo
}
```

---

## 实施工作量

**工期**: 2周
**成本**: ¥6万
**优先级**: P0
