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
	"testing"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHybridSearchService_Search 测试混合搜索
func TestHybridSearchService_Search(t *testing.T) {
	// Skip in CI environment
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// 创建混合搜索引擎
	engine, err := setupTestHybridSearchEngine(ctx, t)
	require.NoError(t, err, "setup test engine failed")
	require.NotNil(t, engine, "engine should not be nil")

	// 测试用例
	tests := []struct {
		name    string
		req     *SearchRequest
		wantErr bool
		check   func(t *testing.T, result *SearchResult)
	}{
		{
			name: "正常搜索",
			req: &SearchRequest{
				Query:       "如何使用ZKER平台",
				TopK:        10,
				KnowledgeID: 1,
				SearchType:  SearchTypeHybrid,
				EnableCache: false,
			},
			wantErr: false,
			check: func(t *testing.T, result *SearchResult) {
				assert.NotNil(t, result, "result should not be nil")
				assert.LessOrEqual(t, len(result.Documents), 10, "should return at most 10 documents")
				assert.Greater(t, result.TotalTimeMs, int64(0), "should have taken some time")

				// 验证分数排序
				for i := 1; i < len(result.Documents); i++ {
					assert.GreaterOrEqual(t, result.Documents[i-1].Score(), result.Documents[i].Score(),
						"documents should be sorted by score descending")
				}
			},
		},
		{
			name: "空查询",
			req: &SearchRequest{
				Query:       "",
				TopK:        10,
				KnowledgeID: 1,
			},
			wantErr: true,
			check:   nil,
		},
		{
			name: "TopK过滤",
			req: &SearchRequest{
				Query:       "租户管理",
				TopK:        5,
				KnowledgeID: 1,
				SearchType:  SearchTypeHybrid,
			},
			wantErr: false,
			check: func(t *testing.T, result *SearchResult) {
				assert.LessOrEqual(t, len(result.Documents), 5, "should return at most 5 documents")
			},
		},
		{
			name: "MinScore过滤",
			req: &SearchRequest{
				Query:       "不存在的查询内容xyz123",
				TopK:        10,
				KnowledgeID: 1,
				MinScore:    0.5,
				SearchType:  SearchTypeHybrid,
			},
			wantErr: false,
			check: func(t *testing.T, result *SearchResult) {
				// 所有文档分数应该 >= MinScore
				for _, doc := range result.Documents {
					assert.GreaterOrEqual(t, doc.Score(), 0.5, "document score should be >= MinScore")
				}
			},
		},
		{
			name: "纯向量搜索",
			req: &SearchRequest{
				Query:       "语义搜索",
				TopK:        10,
				KnowledgeID: 1,
				SearchType:  SearchTypeVector,
			},
			wantErr: false,
			check: func(t *testing.T, result *SearchResult) {
				assert.NotNil(t, result, "result should not be nil")
				assert.NotNil(t, result.SearchMeta, "search meta should not be nil")
				assert.Equal(t, "rrf", result.SearchMeta.FusionStrategy, "fusion strategy should be rrf")
			},
		},
		{
			name: "纯关键词搜索",
			req: &SearchRequest{
				Query:       "关键词匹配",
				TopK:        10,
				KnowledgeID: 1,
				SearchType:  SearchTypeKeyword,
			},
			wantErr: false,
			check: func(t *testing.T, result *SearchResult) {
				assert.NotNil(t, result, "result should not be nil")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.Search(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err, "should return error")
				return
			}

			assert.NoError(t, err, "should not return error")
			if tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

// TestHybridSearchService_Cache 测试缓存功能
func TestHybridSearchService_Cache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// 创建带缓存的搜索引擎
	engine := setupTestHybridSearchServiceWithCache(t)
	require.NotNil(t, engine, "engine should not be nil")

	req := &SearchRequest{
		Query:       "缓存测试查询",
		TopK:        5,
		KnowledgeID: 1,
		SearchType:  SearchTypeHybrid,
		EnableCache: true,
	}

	// 第一次查询 - 应该缓存未命中
	start1 := time.Now()
	result1, err := engine.Search(ctx, req)
	require.NoError(t, err, "first search should not fail")
	duration1 := time.Since(start1)

	assert.NotNil(t, result1, "first result should not be nil")
	assert.False(t, result1.SearchMeta.CacheHit, "first search should not hit cache")

	// 第二次查询 - 应该缓存命中
	start2 := time.Now()
	result2, err := engine.Search(ctx, req)
	require.NoError(t, err, "second search should not fail")
	duration2 := time.Since(start2)

	assert.NotNil(t, result2, "second result should not be nil")
	assert.True(t, result2.SearchMeta.CacheHit, "second search should hit cache")
	assert.Less(t, duration2, duration1, "cached search should be faster")

	// 检查统计信息
	stats := engine.GetCacheStats()
	assert.NotNil(t, stats, "stats should not be nil")
	assert.Greater(t, stats.Hits, int64(0), "should have cache hits")
	assert.Greater(t, stats.HitRate, 0.0, "hit rate should be greater than 0")
}

// TestReciprocalRankFusion 测试RRF算法
func TestReciprocalRankFusion(t *testing.T) {
	// 创建测试数据
	vectorDocs := createMockDocuments([]string{"doc1", "doc2", "doc3", "doc4", "doc5"}, 0.9)
	keywordDocs := createMockDocuments([]string{"doc3", "doc1", "doc6", "doc2", "doc4"}, 0.8)

	engine := &HybridSearchService{}

	// 执行RRF融合
	fusedDocs := engine.reciprocalRankFusion(vectorDocs, keywordDocs)

	// 验证结果
	assert.NotNil(t, fusedDocs, "fused docs should not be nil")
	assert.Equal(t, 6, len(fusedDocs), "should have 6 unique documents")

	// 验证分数排序
	for i := 1; i < len(fusedDocs); i++ {
		assert.GreaterOrEqual(t, fusedDocs[i-1].Score(), fusedDocs[i].Score(),
			"documents should be sorted by score descending")
	}

	// doc3和doc1应该排名靠前(出现在两个列表中)
	top2IDs := []string{fusedDocs[0].ID, fusedDocs[1].ID}
	assert.Contains(t, top2IDs, "doc1", "doc1 should be in top 2")
	assert.Contains(t, top2IDs, "doc3", "doc3 should be in top 2")
}

// TestSearchCache 测试搜索缓存
func TestSearchCache(t *testing.T) {
	ctx := context.Background()
	cache := NewSearchCache(10, 5*time.Minute)

	// 测试Set和Get
	result := &SearchResult{
		Documents: createMockDocuments([]string{"doc1", "doc2"}, 0.9),
		SearchMeta: &SearchMetadata{
			VectorCount:  2,
			KeywordCount: 2,
		},
	}

	cache.Set(ctx, "test query", 1, 10, result)
	cached := cache.Get(ctx, "test query", 1, 10)

	assert.NotNil(t, cached, "cached result should not be nil")
	assert.Equal(t, 2, len(cached.Documents), "should have 2 documents")

	// 测试未命中
	miss := cache.Get(ctx, "non-existent query", 1, 10)
	assert.Nil(t, miss, "miss should return nil")

	// 测试统计信息
	stats := cache.GetStats()
	assert.Equal(t, 1, stats.Size, "cache size should be 1")
	assert.Equal(t, int64(1), stats.Hits, "should have 1 hit")
	assert.Equal(t, int64(1), stats.Misses, "should have 1 miss")
	assert.Greater(t, stats.HitRate, 0.0, "hit rate should be greater than 0")
}

// TestHybridSearchService_ParallelSearch 测试并行搜索
func TestHybridSearchService_ParallelSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	engine, err := setupTestHybridSearchEngine(ctx, t)
	require.NoError(t, err)

	queries := []string{
		"查询1",
		"查询2",
		"查询3",
	}

	results, err := engine.ParallelSearch(ctx, queries, 5)
	require.NoError(t, err, "parallel search should not fail")
	assert.Equal(t, len(queries), len(results), "should return results for all queries")

	for i, result := range results {
		assert.NotNil(t, result, "result %d should not be nil", i)
		assert.LessOrEqual(t, len(result), 5, "result %d should have at most 5 documents", i)
	}
}

// setupTestHybridSearchEngine 设置测试用的混合搜索引擎
func setupTestHybridSearchEngine(ctx context.Context, t *testing.T) (*HybridSearchService, error) {
	// TODO: 实现测试环境设置
	// 1. 创建测试用的向量存储
	// 2. 创建测试用的关键词存储
	// 3. 创建Reranker
	// 4. 创建QueryOptimizer
	// 5. 创建HybridSearchService

	return nil, nil
}

// setupTestHybridSearchServiceWithCache 设置带缓存的测试搜索引擎
func setupTestHybridSearchServiceWithCache(t *testing.T) *CachedSearchEngine {
	cache := NewSearchCache(100, 5*time.Minute)

	engine := &HybridSearchService{}

	return NewCachedSearchEngine(engine, cache)
}

// createMockDocuments 创建模拟文档
func createMockDocuments(ids []string, baseScore float64) []*schema.Document {
	docs := make([]*schema.Document, len(ids))
	for i, id := range ids {
		docs[i] = &schema.Document{
			ID:      id,
			Content: "This is the content of document " + id,
			MetaData: map[string]interface{}{
				"title": "Document " + id,
			},
		}
		docs[i].WithScore(baseScore - float64(i)*0.1)
	}
	return docs
}

// BenchmarkHybridSearch 混合搜索性能基准测试
func BenchmarkHybridSearch(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark test")
	}

	ctx := context.Background()
	engine, err := setupTestHybridSearchEngine(ctx, &testing.T{})
	if err != nil {
		b.Fatalf("setup test engine failed: %v", err)
	}

	req := &SearchRequest{
		Query:       "性能测试查询",
		TopK:        10,
		KnowledgeID: 1,
		SearchType:  SearchTypeHybrid,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Search(ctx, req)
		if err != nil {
			b.Fatalf("search failed: %v", err)
		}
	}
}

// BenchmarkSearchCache 缓存性能基准测试
func BenchmarkSearchCache(b *testing.B) {
	ctx := context.Background()
	cache := NewSearchCache(10000, 5*time.Minute)

	result := &SearchResult{
		Documents: createMockDocuments([]string{"doc1", "doc2", "doc3"}, 0.9),
	}

	// 预热缓存
	cache.Set(ctx, "test query", 1, 10, result)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(ctx, "test query", 1, 10)
	}
}
