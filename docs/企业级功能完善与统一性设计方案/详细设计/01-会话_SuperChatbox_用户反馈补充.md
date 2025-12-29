# 01-会话模块_用户反馈补充文档

**模块**: 01-会话_SuperChatbox扩展
**扩展内容**: 用户反馈系统
**版本**: v2.0
**日期**: 2025-01-03

---

## 新增功能

### 用户反馈收集

**前端组件**:
```typescript
<FeedbackPanel
  conversationId="conv-123"
  onFeedback={(rating, reason) => {
    // 上报反馈
  }}
/>
```

**数据库表**:
```sql
CREATE TABLE conversation_feedback (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id VARCHAR(64) NOT NULL,
    conversation_id VARCHAR(64) NOT NULL,
    message_id VARCHAR(64) NOT NULL,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    is_helpful BOOLEAN,
    feedback_reason VARCHAR(100),
    feedback_text TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**反馈分析**:
- 平均评分
- 有助率
- 常见问题Top 10
- 优化建议

---

## 实施工作量

**工期**: 1周
**成本**: ¥3万
**优先级**: P1
