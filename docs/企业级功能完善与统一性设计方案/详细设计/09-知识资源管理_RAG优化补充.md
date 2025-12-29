# 09-知识资源管理_RAG优化补充文档

**模块**: 09-知识资源管理扩展
**扩展内容**: RAG深度优化
**版本**: v2.0
**日期**: 2025-01-03

---

## 新增功能

### F5 - 智能分块引擎

**ChunkingStrategy**:
- `fixed`: 固定大小分块（如512 tokens）
- `semantic`: 基于语义边界的智能分块
- `recursive`: 递归分割（先大后小）
- `hybrid`: 混合策略

### F6 - 混合检索引擎

```go
type HybridSearchConfig struct {
    VectorWeight    float64  `json:"vector_weight"`    // 向量检索权重
    KeywordWeight   float64  `json:"keyword_weight"`   // 关键词检索权重
    RerankerEnabled bool     `json:"reranker_enabled"` // 是否启用重排序
    TopK            int      `json:"top_k"`            // 召回数量
}
```

### F7 - 查询优化引擎

- Query Rewriting（查询重写）
- Query Expansion（查询扩展）
- HyDE（假设性文档嵌入）

### F8 - 引用标注系统

```sql
CREATE TABLE knowledge_citations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    conversation_id VARCHAR(64) NOT NULL,
    message_id VARCHAR(64) NOT NULL,
    chunk_id BIGINT NOT NULL,
    confidence DECIMAL(3,2) COMMENT '置信度 (0-1)',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## 实施工作量

**工期**: 2周
**成本**: ¥6万
**优先级**: P0
