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
	"sort"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// HybridRetrievalService 混合检索服务
// 职责：实现关键词检索和向量检索的融合(RRF算法)
type HybridRetrievalService struct {
	elasticClient ElasticSearchClient
	vectorClient  VectorClient
	logger        *zap.Logger

	// RRF配置
	rrfK          int     // RRF常数,默认60
	keywordWeight float64 // 关键词检索权重
	vectorWeight  float64 // 向量检索权重
	topK          int     // 最终返回数量
}

// NewHybridRetrievalService 创建混合检索服务实例
func NewHybridRetrievalService(
	elasticClient ElasticSearchClient,
	vectorClient VectorClient,
	logger *zap.Logger,
) *HybridRetrievalService {
	return &HybridRetrievalService{
		elasticClient: elasticClient,
		vectorClient:  vectorClient,
		logger:        logger,
		rrfK:          60,
		keywordWeight: 0.5,
		vectorWeight:  0.5,
		topK:          10,
	}
}

// Retrieve 混合检索
func (s *HybridRetrievalService) Retrieve(
	ctx context.Context,
	query string,
	embedding []float64,
	tenantID string,
	index string,
) ([]*RetrievalResult, error) {
	// 1. 并行执行关键词检索和向量检索
	type resultChan struct {
		results []*RetrievalResult
		err     error
		source  string
	}

	keywordChan := make(chan resultChan)
	vectorChan := make(chan resultChan)

	// 关键词检索
	go func() {
		results, err := s.elasticClient.Search(ctx, index, query, tenantID, s.topK*2)
		keywordChan <- resultChan{results: results, err: err, source: "keyword"}
	}()

	// 向量检索
	go func() {
		results, err := s.vectorClient.Search(ctx, index, embedding, s.topK*2)
		vectorChan <- resultChan{results: results, err: err, source: "vector"}
	}()

	// 等待两个检索完成
	keywordResult := <-keywordChan
	vectorResult := <-vectorChan

	if keywordResult.err != nil {
		s.logger.Warn("keyword search failed", zap.Error(keywordResult.err))
	}
	if vectorResult.err != nil {
		s.logger.Warn("vector search failed", zap.Error(vectorResult.err))
	}

	// 2. 使用RRF算法融合结果
	fusedResults := s.fuseByRRF(keywordResult.results, vectorResult.results)

	// 3. 限制返回数量
	if len(fusedResults) > s.topK {
		fusedResults = fusedResults[:s.topK]
	}

	return fusedResults, nil
}

// fuseByRRF 使用RRF(Reciprocal Rank Fusion)算法融合结果
func (s *HybridRetrievalService) fuseByRRF(
	keywordResults []*RetrievalResult,
	vectorResults []*RetrievalResult,
) []*RetrievalResult {
	// 构建文档ID到RRF得分的映射
	scoreMap := make(map[string]*RRFScore)

	// 处理关键词检索结果
	for i, result := range keywordResults {
		if _, exists := scoreMap[result.DocID]; !exists {
			scoreMap[result.DocID] = &RRFScore{
				DocID:     result.DocID,
				Content:   result.Content,
				Metadata:  result.Metadata,
				KeywordRank: i + 1,
				VectorRank:  -1,
			}
		} else {
			scoreMap[result.DocID].KeywordRank = i + 1
		}
	}

	// 处理向量检索结果
	for i, result := range vectorResults {
		if _, exists := scoreMap[result.DocID]; !exists {
			scoreMap[result.DocID] = &RRFScore{
				DocID:     result.DocID,
				Content:   result.Content,
				Metadata:  result.Metadata,
				KeywordRank: -1,
				VectorRank:  i + 1,
			}
		} else {
			scoreMap[result.DocID].VectorRank = i + 1
		}
	}

	// 计算RRF得分
	for _, rrfScore := range scoreMap {
		rrfScore.CalculateRRF(s.rrfK, s.keywordWeight, s.vectorWeight)
	}

	// 转换为列表并排序
	results := make([]*RetrievalResult, 0, len(scoreMap))
	for _, rrfScore := range scoreMap {
		result := &RetrievalResult{
			DocID:    rrfScore.DocID,
			Content:  rrfScore.Content,
			Score:    rrfScore.RRFScore,
			Metadata: rrfScore.Metadata,
			Source:   "hybrid",
		}
		results = append(results, result)
	}

	// 按RRF得分降序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// ElasticSearchClient Elasticsearch客户端接口
type ElasticSearchClient interface {
	Search(ctx context.Context, index, query, tenantID string, topK int) ([]*RetrievalResult, error)
}

// VectorClient 向量客户端接口
type VectorClient interface {
	Search(ctx context.Context, index string, embedding []float64, topK int) ([]*RetrievalResult, error)
}

// RetrievalResult 检索结果
type RetrievalResult struct {
	DocID    string                 `json:"doc_id"`
	Content  string                 `json:"content"`
	Score    float64                `json:"score"`
	Metadata map[string]string      `json:"metadata"`
	Source   string                 `json:"source"` // keyword, vector, hybrid
}

// RRFScore RRF得分
type RRFScore struct {
	DocID       string
	Content     string
	Metadata    map[string]string
	KeywordRank int // 关键词检索排名,-1表示未出现
	VectorRank  int // 向量检索排名,-1表示未出现
	RRFScore    float64
}

// CalculateRRF 计算RRF得分
func (r *RRFScore) CalculateRRF(k int, keywordWeight, vectorWeight float64) {
	score := 0.0

	// 关键词检索得分
	if r.KeywordRank > 0 {
		score += keywordWeight * (1.0 / float64(k+r.KeywordRank))
	}

	// 向量检索得分
	if r.VectorRank > 0 {
		score += vectorWeight * (1.0 / float64(k+r.VectorRank))
	}

	r.RRFScore = score
}
