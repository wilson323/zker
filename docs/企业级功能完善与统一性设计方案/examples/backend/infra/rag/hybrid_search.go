// Package rag 提供RAG（检索增强生成）实现
//
// 本文件展示混合检索的完整实现，包括：
// - 向量检索（Milvus）
// - 关键词检索（Elasticsearch）
// - 结果融合（RRF算法）
// - 重排序（Reranker）
// - 上下文构建
package rag

import (
	"context"
	"fmt"
	"sort"
)

// HybridSearchService 混合检索服务
type HybridSearchService struct {
	vectorStore    VectorStore
	keywordStore   KeywordSearchStore
	reranker       Reranker
	embeddingModel EmbeddingModel
}

// NewHybridSearchService 创建混合检索服务
func NewHybridSearchService(
	vectorStore VectorStore,
	keywordStore KeywordSearchStore,
	reranker Reranker,
	embeddingModel EmbeddingModel,
) *HybridSearchService {
	return &HybridSearchService{
		vectorStore:    vectorStore,
		keywordStore:   keywordStore,
		reranker:       reranker,
		embeddingModel: embeddingModel,
	}
}

// Search 执行混合检索
//
// 本方法展示完整的RAG检索流程：
// 1. 生成查询向量
// 2. 并行执行向量和关键词检索
// 3. 使用RRF算法融合结果
// 4. 使用Reranker重排序
// 5. 返回Top-K结果
func (s *HybridSearchService) Search(
	ctx context.Context,
	req *SearchRequest,
) (*SearchResult, error) {
	// 步骤1: 生成查询向量
	queryVector, err := s.embeddingModel.Embed(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// 步骤2: 并行执行向量和关键词检索
	type resultPair struct {
		results []DocumentChunk
		err     error
	}

	vectorChan := make(chan resultPair, 1)
	keywordChan := make(chan resultPair, 1)

	// 向量检索
	go func() {
		results, err := s.vectorStore.Search(ctx, &VectorSearchRequest{
			Vector:        queryVector,
			TopK:          req.VectorTopK,
			Filter:        req.Filter,
			KnowledgeBaseID: req.KnowledgeBaseID,
		})
		vectorChan <- resultPair{results, err}
	}()

	// 关键词检索
	go func() {
		results, err := s.keywordStore.Search(ctx, &KeywordSearchRequest{
			Query:          req.Query,
			TopK:           req.KeywordTopK,
			Filter:         req.Filter,
			KnowledgeBaseID: req.KnowledgeBaseID,
		})
		keywordChan <- resultPair{results, err}
	}()

	// 等待两个检索完成
	vectorRes := <-vectorChan
	keywordRes := <-keywordChan

	if vectorRes.err != nil {
		return nil, fmt.Errorf("vector search failed: %w", vectorRes.err)
	}
	if keywordRes.err != nil {
		return nil, fmt.Errorf("keyword search failed: %w", keywordRes.err)
	}

	// 步骤3: 融合结果（使用RRF算法）
	fusedResults := s.fuseByRRF(
		vectorRes.results,
		keywordRes.results,
		req.RRFK,
	)

	// 步骤4: 重排序
	var finalResults []DocumentChunk
	if req.EnableRerank && s.reranker != nil {
		rerankedResults, err := s.reranker.Rerank(ctx, req.Query, fusedResults, req.TopK)
		if err != nil {
			// Rerank失败，返回融合结果
			finalResults = s.topK(fusedResults, req.TopK)
		} else {
			finalResults = rerankedResults
		}
	} else {
		finalResults = s.topK(fusedResults, req.TopK)
	}

	// 步骤5: 构建返回结果
	result := &SearchResult{
		Chunks:       finalResults,
		TotalVector:  len(vectorRes.results),
		TotalKeyword: len(keywordRes.results),
		TotalFused:   len(fusedResults),
		Query:        req.Query,
	}

	return result, nil
}

// fuseByRRF 使用Reciprocal Rank Fusion算法融合结果
//
// RRF公式：score(d) = sum(1 / (k + rank_d))
//
// 参数：
//   - vectorResults: 向量检索结果
//   - keywordResults: 关键词检索结果
//   - k: 平滑参数（通常为60）
//
// 返回：融合后的结果列表
func (s *HybridSearchService) fuseByRRF(
	vectorResults []DocumentChunk,
	keywordResults []DocumentChunk,
	k int,
) []DocumentChunk {
	// 使用map存储每个文档的融合得分
	scores := make(map[string]*FusedChunk)
	chunkMap := make(map[string]DocumentChunk)

	// 处理向量检索结果
	for rank, chunk := range vectorResults {
		id := chunk.ID
		score := 1.0 / float64(k+rank+1)

		if _, exists := scores[id]; !exists {
			scores[id] = &FusedChunk{
				ID:       id,
				Score:    0,
				RefCount: 0,
			}
			chunkMap[id] = chunk
		}

		scores[id].Score += score
		scores[id].RefCount++
		scores[id].VectorRank = rank
	}

	// 处理关键词检索结果
	for rank, chunk := range keywordResults {
		id := chunk.ID
		score := 1.0 / float64(k+rank+1)

		if _, exists := scores[id]; !exists {
			scores[id] = &FusedChunk{
				ID:       id,
				Score:    0,
				RefCount: 0,
			}
			chunkMap[id] = chunk
		}

		scores[id].Score += score
		scores[id].RefCount++
		scores[id].KeywordRank = rank
	}

	// 转换为切片并排序
	fusedChunks := make([]*FusedChunk, 0, len(scores))
	for _, fc := range scores {
		fusedChunks = append(fusedChunks, fc)
	}

	// 按得分降序排序
	sort.Slice(fusedChunks, func(i, j int) bool {
		// 优先比较得分
		if fusedChunks[i].Score != fusedChunks[j].Score {
			return fusedChunks[i].Score > fusedChunks[j].Score
		}
		// 得相同时，优先选择在两个列表中都出现的文档
		return fusedChunks[i].RefCount > fusedChunks[j].RefCount
	})

	// 构建结果列表
	results := make([]DocumentChunk, 0, len(fusedChunks))
	for _, fc := range fusedChunks {
		chunk := chunkMap[fc.ID]
		chunk.FusedScore = fc.Score
		chunk.VectorRank = fc.VectorRank
		chunk.KeywordRank = fc.KeywordRank
		results = append(results, chunk)
	}

	return results
}

// topK 提取Top-K结果
func (s *HybridSearchService) topK(results []DocumentChunk, k int) []DocumentChunk {
	if k <= 0 || k >= len(results) {
		return results
	}
	return results[:k]
}

// ============ 数据结构定义 ============

// SearchRequest 检索请求
type SearchRequest struct {
	// Query 查询文本
	Query string `json:"query"`
	// KnowledgeBaseID 知识库ID
	KnowledgeBaseID string `json:"knowledge_base_id"`
	// TopK 最终返回的Top-K数量
	TopK int `json:"top_k"`
	// VectorTopK 向量检索的Top-K数量（默认50）
	VectorTopK int `json:"vector_top_k"`
	// KeywordTopK 关键词检索的Top-K数量（默认50）
	KeywordTopK int `json:"keyword_top_k"`
	// RRFK RRF算法的k参数（默认60）
	RRFK int `json:"rrf_k"`
	// EnableRerank 是否启用重排序（默认true）
	EnableRerank bool `json:"enable_rerank"`
	// Filter 过滤条件
	Filter *SearchFilter `json:"filter"`
}

// SearchFilter 检索过滤条件
type SearchFilter struct {
	// DocumentType 文档类型
	DocumentType *string `json:"document_type"`
	// Tags 标签过滤
	Tags []string `json:"tags"`
	// CreatedAfter 创建时间范围-起始
	CreatedAfter *int64 `json:"created_after"`
	// CreatedBefore 创建时间范围-结束
	CreatedBefore *int64 `json:"created_before"`
}

// SearchResult 检索结果
type SearchResult struct {
	// Chunks 检索到的文档块列表
	Chunks []DocumentChunk `json:"chunks"`
	// TotalVector 向量检索总数
	TotalVector int `json:"total_vector"`
	// TotalKeyword 关键词检索总数
	TotalKeyword int `json:"total_keyword"`
	// TotalFused 融合后总数
	TotalFused int `json:"total_fused"`
	// Query 查询文本
	Query string `json:"query"`
}

// DocumentChunk 文档块
type DocumentChunk struct {
	// ID 唯一标识
	ID string `json:"id"`
	// DocumentID 所属文档ID
	DocumentID string `json:"document_id"`
	// KnowledgeBaseID 所属知识库ID
	KnowledgeBaseID string `json:"knowledge_base_id"`
	// Content 块内容
	Content string `json:"content"`
	// Metadata 元数据
	Metadata map[string]interface{} `json:"metadata"`
	// Score 相似度得分（0-1）
	Score float64 `json:"score"`
	// FusedScore 融合后的得分
	FusedScore float64 `json:"fused_score,omitempty"`
	// VectorRank 在向量检索中的排名
	VectorRank int `json:"vector_rank,omitempty"`
	// KeywordRank 在关键词检索中的排名
	KeywordRank int `json:"keyword_rank,omitempty"`
}

// FusedChunk 融合块（内部使用）
type FusedChunk struct {
	ID          string  `json:"id"`
	Score       float64 `json:"score"`
	RefCount    int     `json:"ref_count"`
	VectorRank  int     `json:"vector_rank"`
	KeywordRank int     `json:"keyword_rank"`
}

// ============ 接口定义 ============

// VectorStore 向量存储接口
type VectorStore interface {
	// Search 执行向量检索
	Search(ctx context.Context, req *VectorSearchRequest) ([]DocumentChunk, error)
	// Insert 插入向量
	Insert(ctx context.Context, chunks []DocumentChunk) error
	// Delete 删除向量
	Delete(ctx context.Context, ids []string) error
}

// VectorSearchRequest 向量检索请求
type VectorSearchRequest struct {
	Vector          []float32 `json:"vector"`
	TopK            int       `json:"top_k"`
	Filter          *SearchFilter `json:"filter"`
	KnowledgeBaseID string    `json:"knowledge_base_id"`
}

// KeywordSearchStore 关键词存储接口
type KeywordSearchStore interface {
	// Search 执行关键词检索
	Search(ctx context.Context, req *KeywordSearchRequest) ([]DocumentChunk, error)
	// Insert 插入文档
	Insert(ctx context.Context, chunks []DocumentChunk) error
	// Delete 删除文档
	Delete(ctx context.Context, ids []string) error
}

// KeywordSearchRequest 关键词检索请求
type KeywordSearchRequest struct {
	Query          string `json:"query"`
	TopK           int    `json:"top_k"`
	Filter         *SearchFilter `json:"filter"`
	KnowledgeBaseID string `json:"knowledge_base_id"`
}

// Reranker 重排序器接口
type Reranker interface {
	// Rerank 重排序
	Rerank(ctx context.Context, query string, chunks []DocumentChunk, topK int) ([]DocumentChunk, error)
}

// EmbeddingModel 向量化模型接口
type EmbeddingModel interface {
	// Embed 将文本转换为向量
	Embed(ctx context.Context, text string) ([]float32, error)
	// BatchEmbed 批量向量化
	BatchEmbed(ctx context.Context, texts []string) ([][]float32, error)
}
