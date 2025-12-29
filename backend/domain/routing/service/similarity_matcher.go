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
	"math"
	"sort"

	"github.com/coze-studio/backend/domain/routing/entity"
)

// EmbeddingClient 向量嵌入客户端接口
type EmbeddingClient interface {
	GetEmbedding(ctx context.Context, text string) ([]float32, error)
}

// BotRepository Bot仓储接口（用于获取Bot列表）
type BotRepository interface {
	ListByTenant(ctx context.Context, tenantID string, filter interface{}) ([]*BotInfo, error)
}

// BotInfo Bot信息（简化版）
type BotInfo struct {
	BotID     string
	Name      string
	Desc      string
	Embedding  []float32 // 预计算的向量嵌入
	IsActive  bool
}

// SimilarityMatcher 相似度匹配器
type SimilarityMatcher struct {
	botRepo        BotRepository
	embeddingClient EmbeddingClient
	similarityThreshold float32 // 相似度阈值
}

// NewSimilarityMatcher 创建相似度匹配器实例
func NewSimilarityMatcher(
	botRepo BotRepository,
	embeddingClient EmbeddingClient,
	similarityThreshold float32,
) *SimilarityMatcher {
	if similarityThreshold <= 0 {
		similarityThreshold = 0.7 // 默认阈值
	}

	return &SimilarityMatcher{
		botRepo:             botRepo,
		embeddingClient:      embeddingClient,
		similarityThreshold:  similarityThreshold,
	}
}

// Match 相似度匹配
func (m *SimilarityMatcher) Match(ctx context.Context, input *MatchInput) ([]*MatchOutput, error) {
	// 1. 获取用户输入的嵌入向量
	userEmbedding, err := m.embeddingClient.GetEmbedding(ctx, input.UserInput)
	if err != nil {
		return nil, err
	}

	// 2. 获取租户下所有启用的Bot
	bots, err := m.botRepo.ListByTenant(ctx, input.TenantID, nil)
	if err != nil {
		return nil, err
	}

	results := make([]*MatchOutput, 0)

	// 3. 计算相似度
	for _, bot := range bots {
		if !bot.IsActive || bot.Embedding == nil {
			continue
		}

		similarity := cosineSimilarity(userEmbedding, bot.Embedding)

		if similarity >= m.similarityThreshold {
			results = append(results, &MatchOutput{
				BotID:      bot.BotID,
				Confidence: similarity,
				MatchType:  "similarity",
			})
		}
	}

	// 4. 按相似度排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})

	// 5. 返回 Top-3
	if len(results) > 3 {
		results = results[:3]
	}

	return results, nil
}

// cosineSimilarity 余弦相似度计算
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	denominator := float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB)))
	if denominator == 0 {
		return 0
	}

	return float64(dotProduct / denominator)
}
