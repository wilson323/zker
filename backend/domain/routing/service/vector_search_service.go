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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// VectorSearchService 向量搜索服务
// 职责：提供基于Milvus的向量相似度搜索能力
type VectorSearchService struct {
	vectorStore   VectorStore
	intentRepo    repository.IntentRepository
	sampleRepo    repository.IntentSampleRepository
	logger        *zap.Logger

	// 搜索配置
	defaultTopK   int           // 默认返回Top-K结果
	scoreThreshold float64      // 相似度阈值
	cacheEnabled  bool          // 是否启用缓存
	cacheTTL      time.Duration // 缓存过期时间
}

// NewVectorSearchService 创建向量搜索服务实例
func NewVectorSearchService(
	vectorStore VectorStore,
	intentRepo repository.IntentRepository,
	sampleRepo repository.IntentSampleRepository,
	logger *zap.Logger,
) *VectorSearchService {
	return &VectorSearchService{
		vectorStore:    vectorStore,
		intentRepo:     intentRepo,
		sampleRepo:     sampleRepo,
		logger:         logger,
		defaultTopK:    10,
		scoreThreshold: 0.6,
		cacheEnabled:   true,
		cacheTTL:       5 * time.Minute,
	}
}

// SearchIntentByVector 根据向量搜索意图
func (s *VectorSearchService) SearchIntentByVector(
	ctx context.Context,
	tenantID string,
	embedding []float64,
	topK int,
) ([]*entity.IntentRecognitionResult, error) {
	if len(embedding) == 0 {
		return nil, errorx.New(errno.InvalidParamsCode, errorx.KV("reason", "empty embedding"))
	}

	if topK <= 0 {
		topK = s.defaultTopK
	}

	// 1. 向量相似度搜索
	results, err := s.vectorStore.Search(ctx, embedding, topK*3) // 搜索更多候选
	if err != nil {
		s.logger.Error("vector search failed",
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		return nil, errorx.WrapByCode(err, errno.ErrRouteSearchFailedCode, errorx.KV("error", err.Error()))
	}

	if len(results) == 0 {
		return nil, errorx.New(errno.ErrRouteNotFoundCode, errorx.KV("reason", "no similar vectors found"))
	}

	// 2. 按意图分组并计算得分
	intentScores := s.groupByIntent(results)

	// 3. 过滤低分结果
	filteredScores := make([]*entity.IntentRecognitionResult, 0)
	for _, score := range intentScores {
		if score.Confidence >= s.scoreThreshold {
			filteredScores = append(filteredScores, score)
		}
	}

	if len(filteredScores) == 0 {
		return nil, errorx.New(errno.ErrRouteNotFoundCode, errorx.KV("reason", "no intents above threshold"))
	}

	// 4. 限制返回数量
	if len(filteredScores) > topK {
		filteredScores = filteredScores[:topK]
	}

	return filteredScores, nil
}

// BatchSearchIntents 批量搜索意图
func (s *VectorSearchService) BatchSearchIntents(
	ctx context.Context,
	tenantID string,
	embeddings [][]float64,
	topK int,
) ([][]*entity.IntentRecognitionResult, error) {
	results := make([][]*entity.IntentRecognitionResult, len(embeddings))

	for i, embedding := range embeddings {
		result, err := s.SearchIntentByVector(ctx, tenantID, embedding, topK)
		if err != nil {
			s.logger.Warn("failed to search intent",
				zap.Int("index", i),
				zap.Error(err))
			results[i] = []*entity.IntentRecognitionResult{}
			continue
		}
		results[i] = result
	}

	return results, nil
}

// IndexIntentSamples 索引意图样本到向量数据库
func (s *VectorSearchService) IndexIntentSamples(
	ctx context.Context,
	intentID string,
) error {
	// 1. 获取意图信息
	intent, err := s.intentRepo.GetByID(ctx, intentID)
	if err != nil {
		return errorx.New(errno.ErrRouteNotFoundCode, errorx.KV("reason", "intent not found"), errorx.KV("intent_id", intentID))
	}

	// 2. 获取所有样本
	samples, err := s.sampleRepo.GetByIntentID(ctx, intentID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrIntentIndexingFailedCode, errorx.KV("error", err.Error()))
	}

	// 3. 批量索引
	for _, sample := range samples {
		embedding, err := sample.GetEmbedding()
		if err != nil {
			s.logger.Error("failed to get embedding",
				zap.String("sample_id", sample.SampleID),
				zap.Error(err))
			continue
		}

		metadata := map[string]string{
			"intent_id":   intentID,
			"tenant_id":   intent.TenantID,
			"sample_text": sample.Text,
		}

		if err := s.vectorStore.Insert(ctx, sample.SampleID, embedding, metadata); err != nil {
			s.logger.Error("failed to insert vector",
				zap.String("sample_id", sample.SampleID),
				zap.Error(err))
		}
	}

	s.logger.Info("indexed intent samples",
		zap.String("intent_id", intentID),
		zap.Int("count", len(samples)))

	return nil
}

// DeleteIntentVectors 删除意图向量
func (s *VectorSearchService) DeleteIntentVectors(
	ctx context.Context,
	intentID string,
) error {
	// 1. 获取所有样本ID
	samples, err := s.sampleRepo.GetByIntentID(ctx, intentID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrIntentIndexingFailedCode, errorx.KV("error", err.Error()))
	}

	// 2. 批量删除向量
	for _, sample := range samples {
		if err := s.vectorStore.Delete(ctx, sample.SampleID); err != nil {
			s.logger.Error("failed to delete vector",
				zap.String("sample_id", sample.SampleID),
				zap.Error(err))
		}
	}

	s.logger.Info("deleted intent vectors",
		zap.String("intent_id", intentID),
		zap.Int("count", len(samples)))

	return nil
}

// SearchSimilarQueries 搜索相似查询
func (s *VectorSearchService) SearchSimilarQueries(
	ctx context.Context,
	tenantID string,
	queryEmbedding []float64,
	limit int,
) ([]*SimilarQuery, error) {
	if limit <= 0 {
		limit = 5
	}

	// 向量搜索
	results, err := s.vectorStore.Search(ctx, queryEmbedding, limit)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrRouteSearchFailedCode, errorx.KV("error", err.Error()))
	}

	// 构建相似查询列表
	similarQueries := make([]*SimilarQuery, 0, len(results))
	for _, result := range results {
		query := &SimilarQuery{
			Text:       result.Metadata["sample_text"],
			Similarity: result.Score,
			IntentID:   result.Metadata["intent_id"],
		}
		similarQueries = append(similarQueries, query)
	}

	return similarQueries, nil
}

// groupByIntent 按意图分组并计算得分
func (s *VectorSearchService) groupByIntent(results []*VectorSearchResult) []*entity.IntentRecognitionResult {
	intentScores := make(map[string]*intentScoreInfo)

	// 分组并收集得分
	for _, result := range results {
		intentID := result.Metadata["intent_id"]
		if _, exists := intentScores[intentID]; !exists {
			intentScores[intentID] = &intentScoreInfo{
				IntentID: intentID,
				Scores:   make([]float64, 0),
			}
		}
		intentScores[intentID].Scores = append(intentScores[intentID].Scores, result.Score)
	}

	// 计算每个意图的最终得分
	scoredIntents := make([]*entity.IntentRecognitionResult, 0, len(intentScores))
	for intentID, info := range intentScores {
		// 使用加权平均：Top 3结果权重更高
		avgScore := s.calculateWeightedScore(info.Scores)

		result := &entity.IntentRecognitionResult{
			IntentID:    intentID,
			Confidence:  avgScore,
			MatchMethod: "vector",
		}

		scoredIntents = append(scoredIntents, result)
	}

	// 按得分降序排序
	for i := 0; i < len(scoredIntents); i++ {
		for j := i + 1; j < len(scoredIntents); j++ {
			if scoredIntents[j].Confidence > scoredIntents[i].Confidence {
				scoredIntents[i], scoredIntents[j] = scoredIntents[j], scoredIntents[i]
			}
		}
	}

	return scoredIntents
}

// calculateWeightedScore 计算加权得分
func (s *VectorSearchService) calculateWeightedScore(scores []float64) float64 {
	if len(scores) == 0 {
		return 0.0
	}

	// 对得分排序（降序）
	sortedScores := make([]float64, len(scores))
	copy(sortedScores, scores)
	for i := 0; i < len(sortedScores); i++ {
		for j := i + 1; j < len(sortedScores); j++ {
			if sortedScores[j] > sortedScores[i] {
				sortedScores[i], sortedScores[j] = sortedScores[j], sortedScores[i]
			}
		}
	}

	// Top 3加权平均
	weights := []float64{0.5, 0.3, 0.2} // 前3个权重
	sum := 0.0
	weightSum := 0.0

	for i, score := range sortedScores {
		if i >= len(weights) {
			break
		}
		sum += score * weights[i]
		weightSum += weights[i]
	}

	if weightSum == 0 {
		return 0.0
	}

	return sum / weightSum
}

// intentScoreInfo 意图得分信息
type intentScoreInfo struct {
	IntentID string
	Scores   []float64
}

// SimilarQuery 相似查询
type SimilarQuery struct {
	Text       string  `json:"text"`
	Similarity float64 `json:"similarity"`
	IntentID   string  `json:"intent_id"`
}

// GetIndexingStatus 获取索引状态
func (s *VectorSearchService) GetIndexingStatus(
	ctx context.Context,
	intentID string,
) (*IndexingStatus, error) {
	// 获取样本数量
	samples, err := s.sampleRepo.GetByIntentID(ctx, intentID)
	if err != nil {
		return nil, err
	}

	status := &IndexingStatus{
		IntentID:     intentID,
		TotalSamples: len(samples),
		IsIndexed:    len(samples) > 0,
		LastIndexed:  time.Now(),
	}

	return status, nil
}

// IndexingStatus 索引状态
type IndexingStatus struct {
	IntentID     string    `json:"intent_id"`
	TotalSamples int       `json:"total_samples"`
	IsIndexed    bool      `json:"is_indexed"`
	LastIndexed  time.Time `json:"last_indexed"`
}
