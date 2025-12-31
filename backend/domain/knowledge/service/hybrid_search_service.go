/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/infra/document/rerank"
	"github.com/coze-dev/coze-studio/backend/infra/document/searchstore"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// HybridSearchService 混合搜索引擎
// 结合向量搜索和关键词搜索,通过RRF融合和Reranker二次排序提升检索准确率
// 目标: 从65%提升到85%
type HybridSearchService struct {
	vectorStoreManager  searchstore.Manager
	keywordStoreManager searchstore.Manager
	advancedReranker     rerank.Reranker
	queryOptimizer       *QueryOptimizer
	cache                *SearchCache
}

// HybridSearchConfig 混合搜索配置
type HybridSearchConfig struct {
	VectorStoreManager   searchstore.Manager
	KeywordStoreManager  searchstore.Manager
	AdvancedReranker     rerank.Reranker
	QueryOptimizer       *QueryOptimizer
	Cache                *SearchCache
}

// NewHybridSearchService 创建混合搜索引擎
func NewHybridSearchService(config *HybridSearchConfig) (*HybridSearchService, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if config.VectorStoreManager == nil {
		return nil, fmt.Errorf("vector store manager is required")
	}
	if config.KeywordStoreManager == nil {
		return nil, fmt.Errorf("keyword store manager is required")
	}

	return &HybridSearchService{
		vectorStoreManager:  config.VectorStoreManager,
		keywordStoreManager: config.KeywordStoreManager,
		advancedReranker:     config.AdvancedReranker,
		queryOptimizer:       config.QueryOptimizer,
		cache:                config.Cache,
	}, nil
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query           string                 // 查询文本
	TopK            int                    // 返回结果数,默认10
	KnowledgeID     int64                  // 知识库ID
	DocumentIDs     []int64                // 文档ID列表
	MinScore        float64                // 最小相关度分数
	EnableCache     bool                   // 是否启用缓存
	SearchType      SearchType             // 搜索类型
	FilterMetadata  map[string]interface{} // 元数据过滤
}

// SearchType 搜索类型
type SearchType int

const (
	SearchTypeHybrid  SearchType = iota // 混合搜索(向量+关键词)
	SearchTypeVector                    // 纯向量搜索
	SearchTypeKeyword                   // 纯关键词搜索
)

// SearchResult 搜索结果
type SearchResult struct {
	Documents   []*schema.Document // 检索到的文档
	Query       string             // 优化后的查询
	SearchMeta  *SearchMetadata    // 搜索元数据
	TotalTimeMs int64              // 总耗时(毫秒)
}

// SearchMetadata 搜索元数据
type SearchMetadata struct {
	VectorCount      int                    // 向量搜索结果数
	KeywordCount     int                    // 关键词搜索结果数
	RerankedCount    int                    // 重排序后结果数
	CacheHit         bool                   // 是否命中缓存
	OriginalQuery    string                 // 原始查询
	ExpandedQueries  []string               // 扩展的查询列表
	VectorTimeMs     int64                  // 向量搜索耗时
	KeywordTimeMs    int64                  // 关键词搜索耗时
	RerankTimeMs     int64                  // 重排序耗时
	FusionStrategy   string                 // 融合策略
	ScoreDistribution map[string]float64    // 分数分布统计
}

// Search 执行混合搜索
func (s *HybridSearchService) Search(ctx context.Context, req *SearchRequest) (*SearchResult, error) {
	startTime := time.Now()

	if req == nil {
		return nil, fmt.Errorf("search request is required")
	}
	if req.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}
	if req.TopK <= 0 {
		req.TopK = 10 // 默认返回10条
	}

	// 1. 检查缓存
	if req.EnableCache && s.cache != nil {
		if cached := s.cache.Get(ctx, req.Query, req.KnowledgeID, req.TopK); cached != nil {
			logs.CtxInfof(ctx, "[HybridSearch] cache hit for query: %s", req.Query)
			cached.SearchMeta.CacheHit = true
			cached.TotalTimeMs = time.Since(startTime).Milliseconds()
			return cached, nil
		}
	}

	// 2. 查询优化(可选)
	optimizedQuery := req.Query
	expandedQueries := []string{}
	if s.queryOptimizer != nil {
		var err error
		optimizedQuery, expandedQueries, err = s.queryOptimizer.ExpandQuery(ctx, req.Query)
		if err != nil {
			logs.CtxWarnf(ctx, "[HybridSearch] query optimization failed: %v", err)
			// 降级使用原始查询
			optimizedQuery = req.Query
		}
	}

	// 3. 并行执行向量搜索和关键词搜索
	var wg sync.WaitGroup
	var vectorDocs []*schema.Document
	var keywordDocs []*schema.Document
	var vectorErr, keywordErr error
	var vectorTime, keywordTime int64

	wg.Add(2)

	// 向量搜索
	go func() {
		defer wg.Done()
		vStart := time.Now()
		vectorDocs, vectorErr = s.vectorSearch(ctx, optimizedQuery, req)
		vectorTime = time.Since(vStart).Milliseconds()
	}()

	// 关键词搜索
	go func() {
		defer wg.Done()
		kStart := time.Now()
		keywordDocs, keywordErr = s.keywordSearch(ctx, optimizedQuery, req)
		keywordTime = time.Since(kStart).Milliseconds()
	}()

	wg.Wait()

	// 处理搜索错误
	if vectorErr != nil {
		logs.CtxWarnf(ctx, "[HybridSearch] vector search failed: %v", vectorErr)
	}
	if keywordErr != nil {
		logs.CtxWarnf(ctx, "[HybridSearch] keyword search failed: %v", keywordErr)
	}

	// 如果两种搜索都失败,返回错误
	if vectorErr != nil && keywordErr != nil {
		return nil, fmt.Errorf("both vector and keyword search failed: vector=%v, keyword=%v", vectorErr, keywordErr)
	}

	// 4. 结果融合(Reciprocal Rank Fusion)
	fusedDocs := s.reciprocalRankFusion(vectorDocs, keywordDocs)

	// 5. Reranker二次排序(如果配置了高级Reranker)
	var rerankedDocs []*schema.Document
	var rerankTime int64

	if s.advancedReranker != nil && len(fusedDocs) > 0 {
		rStart := time.Now()
		var err error
		rerankedDocs, err = s.advancedReranking(ctx, optimizedQuery, fusedDocs, req.TopK*2)
		rerankTime = time.Since(rStart).Milliseconds()
		if err != nil {
			logs.CtxWarnf(ctx, "[HybridSearch] advanced reranking failed: %v", err)
			// 降级使用融合结果
			rerankedDocs = fusedDocs
		}
	} else {
		rerankedDocs = fusedDocs
	}

	// 6. 应用TopK和MinScore过滤
	finalDocs := s.applyFilters(rerankedDocs, req.TopK, req.MinScore)

	// 7. 构建搜索元数据
	searchMeta := &SearchMetadata{
		VectorCount:      len(vectorDocs),
		KeywordCount:     len(keywordDocs),
		RerankedCount:    len(rerankedDocs),
		CacheHit:         false,
		OriginalQuery:    req.Query,
		ExpandedQueries:  expandedQueries,
		VectorTimeMs:     vectorTime,
		KeywordTimeMs:    keywordTime,
		RerankTimeMs:     rerankTime,
		FusionStrategy:   "rrf",
		ScoreDistribution: s.calculateScoreDistribution(finalDocs),
	}

	result := &SearchResult{
		Documents:   finalDocs,
		Query:       optimizedQuery,
		SearchMeta:  searchMeta,
		TotalTimeMs: time.Since(startTime).Milliseconds(),
	}

	// 8. 缓存结果
	if req.EnableCache && s.cache != nil {
		s.cache.Set(ctx, req.Query, req.KnowledgeID, req.TopK, result)
	}

	// 9. 性能日志
	if result.TotalTimeMs > 100 {
		logs.CtxWarnf(ctx, "[HybridSearch] slow query detected: %s, took %dms",
			req.Query, result.TotalTimeMs)
	}

	return result, nil
}

// vectorSearch 向量搜索
func (s *HybridSearchService) vectorSearch(ctx context.Context, query string, req *SearchRequest) ([]*schema.Document, error) {
	if req.SearchType == SearchTypeKeyword {
		return nil, nil // 关键词搜索模式跳过向量搜索
	}

	collectionName := getCollectionName(req.KnowledgeID)
	searchStore, err := s.vectorStoreManager.GetSearchStore(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("get vector store failed: %w", err)
	}

	// 构建检索选项
	opts := s.buildSearchOptions(ctx, req)

	// 执行向量检索
	docs, err := searchStore.Retrieve(ctx, query, opts...)
	if err != nil {
		return nil, fmt.Errorf("vector retrieve failed: %w", err)
	}

	logs.CtxInfof(ctx, "[HybridSearch] vector search returned %d documents", len(docs))
	return docs, nil
}

// keywordSearch 关键词搜索
func (s *HybridSearchService) keywordSearch(ctx context.Context, query string, req *SearchRequest) ([]*schema.Document, error) {
	if req.SearchType == SearchTypeVector {
		return nil, nil // 向量搜索模式跳过关键词搜索
	}

	collectionName := getCollectionName(req.KnowledgeID)
	searchStore, err := s.keywordStoreManager.GetSearchStore(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("get keyword store failed: %w", err)
	}

	// 构建检索选项
	opts := s.buildSearchOptions(ctx, req)

	// 执行关键词检索
	docs, err := searchStore.Retrieve(ctx, query, opts...)
	if err != nil {
		return nil, fmt.Errorf("keyword retrieve failed: %w", err)
	}

	logs.CtxInfof(ctx, "[HybridSearch] keyword search returned %d documents", len(docs))
	return docs, nil
}

// reciprocalRankFusion Reciprocal Rank Fusion算法
// 公式: score(d) = Σ 1/(k + rank(d))
// 其中k是常数(通常为60),rank(d)是文档d在某个列表中的排名
func (s *HybridSearchService) reciprocalRankFusion(vectorDocs, keywordDocs []*schema.Document) []*schema.Document {
	const k = 60.0 // RRF常数

	// 文档ID -> 分数映射
	scores := make(map[string]float64)
	docMap := make(map[string]*schema.Document)

	// 处理向量搜索结果
	for rank, doc := range vectorDocs {
		if doc == nil {
			continue
		}
		score := 1.0 / (k + float64(rank+1))
		scores[doc.ID] += score
		docMap[doc.ID] = doc
	}

	// 处理关键词搜索结果
	for rank, doc := range keywordDocs {
		if doc == nil {
			continue
		}
		score := 1.0 / (k + float64(rank+1))
		scores[doc.ID] += score
		if docMap[doc.ID] == nil {
			docMap[doc.ID] = doc
		}
	}

	// 按分数排序
	type docScore struct {
		doc   *schema.Document
		score float64
	}

	var sortedDocs []docScore
	for id, score := range scores {
		sortedDocs = append(sortedDocs, docScore{
			doc:   docMap[id],
			score: score,
		})
	}

	sort.Slice(sortedDocs, func(i, j int) bool {
		return sortedDocs[i].score > sortedDocs[j].score
	})

	// 返回排序后的文档列表
	result := make([]*schema.Document, 0, len(sortedDocs))
	for _, item := range sortedDocs {
		// 更新文档的分数为RRF分数
		item.doc.WithScore(item.score)
		result = append(result, item.doc)
	}

	logs.CtxInfof(context.Background(), "[HybridSearch] RRF fused %d vector docs and %d keyword docs into %d unique docs",
		len(vectorDocs), len(keywordDocs), len(result))

	return result
}

// advancedReranking 高级重排序
func (s *HybridSearchService) advancedReranking(ctx context.Context, query string, docs []*schema.Document, topN int) ([]*schema.Document, error) {
	if len(docs) == 0 {
		return docs, nil
	}

	// 准备Reranker输入
	data := make([]*rerank.Data, 0, len(docs))
	for _, doc := range docs {
		data = append(data, &rerank.Data{
			Document: doc,
			Score:    doc.Score(),
		})
	}

	// 调用Reranker
	resp, err := s.advancedReranker.Rerank(ctx, &rerank.Request{
		Query: query,
		Data:  [][]*rerank.Data{data},
		TopN:  int64(topN),
	})
	if err != nil {
		return nil, fmt.Errorf("rerank failed: %w", err)
	}

	// 更新文档分数
	for _, item := range resp.SortedData {
		item.Document.WithScore(item.Score)
	}

	logs.CtxInfof(ctx, "[HybridSearch] reranked %d documents to %d", len(docs), len(resp.SortedData))
	return resp.SortedData, nil
}

// applyFilters 应用TopK和MinScore过滤
func (s *HybridSearchService) applyFilters(docs []*schema.Document, topK int, minScore float64) []*schema.Document {
	if len(docs) == 0 {
		return docs
	}

	// 过滤低分文档
	filtered := make([]*schema.Document, 0, len(docs))
	for _, doc := range docs {
		if doc.Score() >= minScore {
			filtered = append(filtered, doc)
		}
	}

	// 应用TopK
	if len(filtered) > topK {
		filtered = filtered[:topK]
	}

	return filtered
}

// buildSearchOptions 构建检索选项
func (s *HybridSearchService) buildSearchOptions(ctx context.Context, req *SearchRequest) []interface{} {
	// 这里根据实际需求构建检索选项
	// 例如: 分区键、过滤条件等
	opts := []interface{}{}

	// 可以添加文档过滤
	if len(req.DocumentIDs) > 0 {
		// TODO: 实现文档ID过滤
	}

	// 可以添加元数据过滤
	if len(req.FilterMetadata) > 0 {
		// TODO: 实现元数据过滤
	}

	return opts
}

// calculateScoreDistribution 计算分数分布统计
func (s *HybridSearchService) calculateScoreDistribution(docs []*schema.Document) map[string]float64 {
	if len(docs) == 0 {
		return map[string]float64{}
	}

	stats := map[string]float64{
		"min":   math.MaxFloat64,
		"max":   -math.MaxFloat64,
		"mean":  0,
		"count": float64(len(docs)),
	}

	sum := 0.0
	for _, doc := range docs {
		score := doc.Score()
		if score < stats["min"] {
			stats["min"] = score
		}
		if score > stats["max"] {
			stats["max"] = score
		}
		sum += score
	}

	stats["mean"] = sum / float64(len(docs))

	return stats
}

// ParallelSearch 并行搜索多个查询
func (s *HybridSearchService) ParallelSearch(ctx context.Context, queries []string, topK int) ([][]*schema.Document, error) {
	if len(queries) == 0 {
		return nil, fmt.Errorf("queries cannot be empty")
	}

	results := make([][]*schema.Document, len(queries))
	var wg sync.WaitGroup
	var mu sync.Mutex
	errGroup := make(chan error, len(queries))

	for i, query := range queries {
		wg.Add(1)
		go func(idx int, q string) {
			defer wg.Done()

			req := &SearchRequest{
				Query:       q,
				TopK:        topK,
				SearchType:  SearchTypeHybrid,
				EnableCache: true,
			}

			result, err := s.Search(ctx, req)
			if err != nil {
				logs.CtxErrorf(ctx, "[HybridSearch] parallel search failed for query '%s': %v", q, err)
				errGroup <- err
				return
			}

			mu.Lock()
			results[idx] = result.Documents
			mu.Unlock()
		}(i, query)
	}

	wg.Wait()
	close(errGroup)

	// 检查是否有错误
	var errs []error
	for err := range errGroup {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("%d out of %d searches failed: %v", len(errs), len(queries), errs[0])
	}

	return results, nil
}

// BatchSearch 批量搜索(优化版,复用连接)
func (s *HybridSearchService) BatchSearch(ctx context.Context, requests []*SearchRequest) ([]*SearchResult, error) {
	if len(requests) == 0 {
		return nil, fmt.Errorf("requests cannot be empty")
	}

	results := make([]*SearchResult, len(requests))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, req := range requests {
		wg.Add(1)
		go func(idx int, r *SearchRequest) {
			defer wg.Done()

			result, err := s.Search(ctx, r)
			if err != nil {
				logs.CtxErrorf(ctx, "[HybridSearch] batch search failed: %v", err)
				// 返回空结果而不是错误
				result = &SearchResult{Documents: []*schema.Document{}}
			}

			mu.Lock()
			results[idx] = result
			mu.Unlock()
		}(i, req)
	}

	wg.Wait()

	return results, nil
}
