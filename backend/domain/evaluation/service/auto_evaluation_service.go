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

	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// AutoEvaluationService 自动化评估服务
// 职责：实现BLEU/ROUGE评分、语义相似度评分、准确性评分
type AutoEvaluationService struct {
	llmClient     LLMClientForEval
	vectorClient  VectorClientForEval
	logger        *zap.Logger
}

// NewAutoEvaluationService 创建自动化评估服务实例
func NewAutoEvaluationService(
	llmClient LLMClientForEval,
	vectorClient VectorClientForEval,
	logger *zap.Logger,
) *AutoEvaluationService {
	return &AutoEvaluationService{
		llmClient:    llmClient,
		vectorClient: vectorClient,
		logger:       logger,
	}
}

// EvaluateResponse 评估响应质量
func (s *AutoEvaluationService) EvaluateResponse(
	ctx context.Context,
	req *EvaluationRequest,
) (*EvaluationResult, error) {
	result := &EvaluationResult{
		Query:     req.Query,
		Response:  req.Response,
		Reference: req.Reference,
		Metrics:   make(map[string]float64),
	}

	// 1. BLEU评分
	bleuScore, err := s.calculateBLEU(req.Response, req.Reference)
	if err != nil {
		s.logger.Warn("failed to calculate BLEU", zap.Error(err))
	} else {
		result.Metrics["bleu"] = bleuScore
	}

	// 2. ROUGE评分
	rougeScore, err := s.calculateROUGE(req.Response, req.Reference)
	if err != nil {
		s.logger.Warn("failed to calculate ROUGE", zap.Error(err))
	} else {
		result.Metrics["rouge"] = rougeScore
	}

	// 3. 语义相似度评分
	similarityScore, err := s.calculateSemanticSimilarity(ctx, req.Response, req.Reference)
	if err != nil {
		s.logger.Warn("failed to calculate semantic similarity", zap.Error(err))
	} else {
		result.Metrics["semantic_similarity"] = similarityScore
	}

	// 4. 准确性评分
	accuracyScore, err := s.calculateAccuracy(ctx, req)
	if err != nil {
		s.logger.Warn("failed to calculate accuracy", zap.Error(err))
	} else {
		result.Metrics["accuracy"] = accuracyScore
	}

	// 5. 计算综合评分
	result.OverallScore = s.calculateOverallScore(result.Metrics)

	return result, nil
}

// calculateBLEU 计算BLEU评分
func (s *AutoEvaluationService) calculateBLEU(response, reference string) (float64, error) {
	// 简化实现：计算1-gram到4-gram的精度
	// 实际应用中应该使用完整的BLEU算法

	responseWords := tokenize(response)
	referenceWords := tokenize(reference)

	if len(responseWords) == 0 {
		return 0.0, nil
	}

	// 计算匹配数
	matches := 0
	for _, word := range responseWords {
		for _, refWord := range referenceWords {
			if word == refWord {
				matches++
				break
			}
		}
	}

	precision := float64(matches) / float64(len(responseWords))

	// 简化的BLEU分数(仅使用1-gram)
	return precision, nil
}

// calculateROUGE 计算ROUGE评分
func (s *AutoEvaluationService) calculateROUGE(response, reference string) (float64, error) {
	// 简化实现：计算ROUGE-L(最长公共子序列)
	lcs := longestCommonSubsequence(response, reference)
	referenceLen := len(reference)

	if referenceLen == 0 {
		return 0.0, nil
	}

	rougeL := float64(lcs) / float64(referenceLen)
	return rougeL, nil
}

// calculateSemanticSimilarity 计算语义相似度
func (s *AutoEvaluationService) calculateSemanticSimilarity(
	ctx context.Context,
	response, reference string,
) (float64, error) {
	// 生成向量嵌入
	responseEmbedding, err := s.vectorClient.GetEmbedding(ctx, response)
	if err != nil {
		return 0.0, err
	}

	referenceEmbedding, err := s.vectorClient.GetEmbedding(ctx, reference)
	if err != nil {
		return 0.0, err
	}

	// 计算余弦相似度
	similarity := cosineSimilarity(responseEmbedding, referenceEmbedding)
	return similarity, nil
}

// calculateAccuracy 计算准确性评分
func (s *AutoEvaluationService) calculateAccuracy(
	ctx context.Context,
	req *EvaluationRequest,
) (float64, error) {
	// 简化实现：使用LLM评估准确性
	prompt := buildAccuracyPrompt(req.Query, req.Response, req.Reference)

	accuracy, err := s.llmClient.EvaluateAccuracy(ctx, prompt)
	if err != nil {
		return 0.0, err
	}

	return accuracy, nil
}

// calculateOverallScore 计算综合评分
func (s *AutoEvaluationService) calculateOverallScore(metrics map[string]float64) float64 {
	weights := map[string]float64{
		"bleu":              0.25,
		"rouge":             0.25,
		"semantic_similarity": 0.30,
		"accuracy":          0.20,
	}

	totalScore := 0.0
	totalWeight := 0.0

	for metric, weight := range weights {
		if score, exists := metrics[metric]; exists {
			totalScore += score * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0.0
	}

	return totalScore / totalWeight
}

// BatchEvaluate 批量评估
func (s *AutoEvaluationService) BatchEvaluate(
	ctx context.Context,
	requests []*EvaluationRequest,
) ([]*EvaluationResult, error) {
	results := make([]*EvaluationResult, len(requests))

	for i, req := range requests {
		result, err := s.EvaluateResponse(ctx, req)
		if err != nil {
			s.logger.Warn("failed to evaluate response",
				zap.Int("index", i),
				zap.Error(err))
			results[i] = nil
			continue
		}
		results[i] = result
	}

	return results, nil
}

// LLMClientForEval LLM评估客户端接口
type LLMClientForEval interface {
	EvaluateAccuracy(ctx context.Context, prompt string) (float64, error)
}

// VectorClientForEval 向量评估客户端接口
type VectorClientForEval interface {
	GetEmbedding(ctx context.Context, text string) ([]float64, error)
}

// EvaluationRequest 评估请求
type EvaluationRequest struct {
	Query     string `json:"query"`
	Response  string `json:"response"`
	Reference string `json:"reference"` // 参考答案
}

// EvaluationResult 评估结果
type EvaluationResult struct {
	Query       string              `json:"query"`
	Response    string              `json:"response"`
	Reference   string              `json:"reference"`
	Metrics     map[string]float64  `json:"metrics"`
	OverallScore float64            `json:"overall_score"`
}

// tokenize 分词
func tokenize(text string) []string {
	// 简化实现：按空格分词
	words := make([]string, 0)
	currentWord := ""

	for _, char := range text {
		if char == ' ' || char == '\t' || char == '\n' {
			if currentWord != "" {
				words = append(words, currentWord)
				currentWord = ""
			}
		} else {
			currentWord += string(char)
		}
	}

	if currentWord != "" {
		words = append(words, currentWord)
	}

	return words
}

// longestCommonSubsequence 最长公共子序列长度
func longestCommonSubsequence(s1, s2 string) int {
	m := len(s1)
	n := len(s2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

// cosineSimilarity 余弦相似度
func cosineSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return 0.0
	}

	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		norm1 += vec1[i] * vec1[i]
		norm2 += vec2[i] * vec2[i]
	}

	if norm1 == 0 || norm2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// max 返回最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// buildAccuracyPrompt 构建准确性评估提示词
func buildAccuracyPrompt(query, response, reference string) string {
	return fmt.Sprintf(`请评估以下AI响应的准确性(0-1分):

用户问题: %s
AI响应: %s
参考答案: %s

请仅返回一个0到1之间的数字，表示响应的准确性。`)
}
