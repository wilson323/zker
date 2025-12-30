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
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/coze-dev/coze-studio/backend/domain/routing/entity"
	"github.com/coze-dev/coze-studio/backend/domain/routing/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"go.uber.org/zap"
)

// LLMClient LLM客户端接口
type LLMClient interface {
	// GenerateText 生成文本
	GenerateText(ctx context.Context, prompt string) (string, error)

	// GenerateEmbedding 生成向量嵌入
	GenerateEmbedding(ctx context.Context, text string) ([]float64, error)
}

// VectorStore 向量存储接口
type VectorStore interface {
	// Search 向量相似度搜索
	Search(ctx context.Context, embedding []float64, limit int) ([]*VectorSearchResult, error)

	// Insert 插入向量
	Insert(ctx context.Context, id string, embedding []float64, metadata map[string]string) error

	// Delete 删除向量
	Delete(ctx context.Context, id string) error
}

// VectorSearchResult 向量搜索结果
type VectorSearchResult struct {
	ID       string
	Score    float64
	Metadata map[string]string
}

// IntentRecognitionService 意图识别服务
// 职责：识别用户输入的意图，支持LLM和向量相似度两种方式
type IntentRecognitionService struct {
	intentRepo    repository.IntentRepository
	sampleRepo    repository.IntentSampleRepository
	llmClient     LLMClient
	vectorStore   VectorStore
	logger        *zap.Logger

	// 识别策略
	strategy      string // llm, vector, hybrid
	hybridWeight  float64 // 混合策略中LLM的权重
}

// NewIntentRecognitionService 创建意图识别服务实例
func NewIntentRecognitionService(
	intentRepo repository.IntentRepository,
	sampleRepo repository.IntentSampleRepository,
	llmClient LLMClient,
	vectorStore VectorStore,
	logger *zap.Logger,
) *IntentRecognitionService {
	return &IntentRecognitionService{
		intentRepo:   intentRepo,
		sampleRepo:   sampleRepo,
		llmClient:    llmClient,
		vectorStore:  vectorStore,
		logger:       logger,
		strategy:     "hybrid", // 默认使用混合策略
		hybridWeight: 0.6,      // LLM权重60%
	}
}

// RecognizeIntent 识别意图
func (s *IntentRecognitionService) RecognizeIntent(ctx context.Context, tenantID, text string) (*entity.IntentRecognitionResult, error) {
	if text == "" {
		return nil, errno.ROUTING201001.WithDetail("reason", "empty text")
	}

	// 获取租户的所有激活意图
	intents, err := s.intentRepo.GetActiveIntentsByTenant(ctx, tenantID)
	if err != nil {
		s.logger.Error("failed to get active intents",
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		return nil, errno.ROUTING500002.WithDetail("error", err.Error())
	}

	if len(intents) == 0 {
		return nil, errno.ROUTING404001.WithDetail("reason", "no active intents")
	}

	// 根据策略识别意图
	var result *entity.IntentRecognitionResult
	switch s.strategy {
	case "llm":
		result, err = s.recognizeByLLM(ctx, text, intents)
	case "vector":
		result, err = s.recognizeByVector(ctx, tenantID, text, intents)
	case "hybrid":
		result, err = s.recognizeByHybrid(ctx, tenantID, text, intents)
	default:
		return nil, errno.ROUTING201001.WithDetail("reason", "unknown recognition strategy")
	}

	if err != nil {
		return nil, err
	}

	// 检查置信度是否达标
	if result.Confidence < 0.5 {
		return nil, errno.ROUTING500002.WithDetail("reason", "low confidence", "confidence", result.Confidence)
	}

	return result, nil
}

// BatchRecognize 批量识别意图
func (s *IntentRecognitionService) BatchRecognize(ctx context.Context, tenantID string, texts []string) ([]*entity.IntentRecognitionResult, error) {
	results := make([]*entity.IntentRecognitionResult, 0, len(texts))

	for _, text := range texts {
		result, err := s.RecognizeIntent(ctx, tenantID, text)
		if err != nil {
			s.logger.Warn("failed to recognize intent",
				zap.String("text", text),
				zap.Error(err))
			// 继续处理其他文本
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// TrainIntentModel 训练意图模型
func (s *IntentRecognitionService) TrainIntentModel(ctx context.Context, intentID string, samples []string) error {
	// 1. 验证意图是否存在
	intent, err := s.intentRepo.GetByID(ctx, intentID)
	if err != nil {
		return errno.ROUTING404001.WithDetail("reason", "intent not found", "intent_id", intentID)
	}

	// 2. 删除旧样本
	if err := s.sampleRepo.DeleteByIntentID(ctx, intentID); err != nil {
		s.logger.Error("failed to delete old samples",
			zap.String("intent_id", intentID),
			zap.Error(err))
		return err
	}

	// 3. 批量创建新样本
	intentSamples := make([]*entity.IntentSample, 0, len(samples))
	for i, sampleText := range samples {
		// 生成向量嵌入
		embedding, err := s.llmClient.GenerateEmbedding(ctx, sampleText)
		if err != nil {
			s.logger.Error("failed to generate embedding",
				zap.String("sample", sampleText),
				zap.Error(err))
			continue
		}

		// 创建样本记录
		sample := &entity.IntentSample{
			SampleID: fmt.Sprintf("%s_%d", intentID, i),
			IntentID: intentID,
			Text:     sampleText,
		}
		if err := sample.SetEmbedding(embedding); err != nil {
			s.logger.Error("failed to set embedding",
				zap.String("sample", sampleText),
				zap.Error(err))
			continue
		}

		intentSamples = append(intentSamples, sample)

		// 插入向量存储
		metadata := map[string]string{
			"intent_id":   intentID,
			"tenant_id":   intent.TenantID,
			"sample_text": sampleText,
		}
		if err := s.vectorStore.Insert(ctx, sample.SampleID, embedding, metadata); err != nil {
			s.logger.Error("failed to insert vector",
				zap.String("sample_id", sample.SampleID),
				zap.Error(err))
		}
	}

	// 批量保存样本到数据库
	if err := s.sampleRepo.CreateBatch(ctx, intentSamples); err != nil {
		s.logger.Error("failed to create samples",
			zap.String("intent_id", intentID),
			zap.Error(err))
		return err
	}

	return nil
}

// UpdateIntentSamples 更新意图样本
func (s *IntentRecognitionService) UpdateIntentSamples(ctx context.Context, intentID string, samples []string) error {
	return s.TrainIntentModel(ctx, intentID, samples)
}

// recognizeByLLM 使用LLM识别意图
func (s *IntentRecognitionService) recognizeByLLM(ctx context.Context, text string, intents []*entity.Intent) (*entity.IntentRecognitionResult, error) {
	// 构建提示词
	prompt := s.buildLLMPrompt(text, intents)

	// 调用LLM
	response, err := s.llmClient.GenerateText(ctx, prompt)
	if err != nil {
		return nil, errno.ROUTING500002.WithDetail("error", err.Error())
	}

	// 解析LLM响应
	var llmResult struct {
		IntentName string  `json:"intent_name"`
		Confidence float64 `json:"confidence"`
		Reasoning  string  `json:"reasoning"`
	}
	if err := json.Unmarshal([]byte(response), &llmResult); err != nil {
		return nil, errno.ROUTING500002.WithDetail("error", "failed to parse llm response")
	}

	// 查找匹配的意图
	var matchedIntent *entity.Intent
	for _, intent := range intents {
		if intent.IntentName == llmResult.IntentName {
			matchedIntent = intent
			break
		}
	}

	if matchedIntent == nil {
		return nil, errno.ROUTING404001.WithDetail("reason", "intent not found", "intent_name", llmResult.IntentName)
	}

	return &entity.IntentRecognitionResult{
		IntentID:    matchedIntent.IntentID,
		IntentName:  matchedIntent.IntentName,
		Confidence:  llmResult.Confidence,
		AgentID:     matchedIntent.AgentID,
		WorkflowID:  matchedIntent.WorkflowID,
		MatchMethod: "llm",
	}, nil
}

// recognizeByVector 使用向量相似度识别意图
func (s *IntentRecognitionService) recognizeByVector(ctx context.Context, tenantID, text string, intents []*entity.Intent) (*entity.IntentRecognitionResult, error) {
	// 1. 生成文本的向量嵌入
	embedding, err := s.llmClient.GenerateEmbedding(ctx, text)
	if err != nil {
		return nil, errno.ROUTING500002.WithDetail("error", err.Error())
	}

	// 2. 向量相似度搜索
	results, err := s.vectorStore.Search(ctx, embedding, 10)
	if err != nil {
		return nil, errno.ROUTING500002.WithDetail("error", err.Error())
	}

	if len(results) == 0 {
		return nil, errno.ROUTING404001.WithDetail("reason", "no similar vectors found")
	}

	// 3. 按意图分组统计得分
	intentScores := make(map[string]*intentMatchScore)
	for _, result := range results {
		intentID := result.Metadata["intent_id"]
		if _, exists := intentScores[intentID]; !exists {
			intentScores[intentID] = &intentMatchScore{
				IntentID: intentID,
				Scores:   make([]float64, 0),
			}
		}
		intentScores[intentID].Scores = append(intentScores[intentID].Scores, result.Score)
	}

	// 4. 计算每个意图的平均得分
	for _, im := range intentScores {
		im.AvgScore = average(im.Scores)
	}

	// 5. 选择得分最高的意图
	var bestMatch *intentMatchScore
	for _, im := range intentScores {
		if bestMatch == nil || im.AvgScore > bestMatch.AvgScore {
			bestMatch = im
		}
	}

	// 6. 查找意图详情
	var matchedIntent *entity.Intent
	for _, intent := range intents {
		if intent.IntentID == bestMatch.IntentID {
			matchedIntent = intent
			break
		}
	}

	if matchedIntent == nil {
		return nil, errno.ROUTING404001.WithDetail("reason", "intent not found")
	}

	return &entity.IntentRecognitionResult{
		IntentID:    matchedIntent.IntentID,
		IntentName:  matchedIntent.IntentName,
		Confidence:  bestMatch.AvgScore,
		AgentID:     matchedIntent.AgentID,
		WorkflowID:  matchedIntent.WorkflowID,
		MatchMethod: "vector",
	}, nil
}

// recognizeByHybrid 使用混合策略识别意图
func (s *IntentRecognitionService) recognizeByHybrid(ctx context.Context, tenantID, text string, intents []*entity.Intent) (*entity.IntentRecognitionResult, error) {
	// 1. 使用LLM识别
	llmResult, err := s.recognizeByLLM(ctx, text, intents)
	if err != nil {
		s.logger.Warn("LLM recognition failed, fallback to vector", zap.Error(err))
		return s.recognizeByVector(ctx, tenantID, text, intents)
	}

	// 2. 使用向量识别
	vectorResult, err := s.recognizeByVector(ctx, tenantID, text, intents)
	if err != nil {
		s.logger.Warn("vector recognition failed, use LLM result", zap.Error(err))
		return llmResult, nil
	}

	// 3. 混合两种结果
	// 如果两者识别到同一个意图，取加权置信度
	if llmResult.IntentID == vectorResult.IntentID {
		hybridConfidence := llmResult.Confidence*s.hybridWeight + vectorResult.Confidence*(1-s.hybridWeight)
		return &entity.IntentRecognitionResult{
			IntentID:    llmResult.IntentID,
			IntentName:  llmResult.IntentName,
			Confidence:  hybridConfidence,
			AgentID:     llmResult.AgentID,
			WorkflowID:  llmResult.WorkflowID,
			MatchMethod: "hybrid",
		}, nil
	}

	// 如果识别到不同意图，选择置信度更高的
	if llmResult.Confidence >= vectorResult.Confidence {
		llmResult.MatchMethod = "hybrid_llm"
		return llmResult, nil
	}

	vectorResult.MatchMethod = "hybrid_vector"
	return vectorResult, nil
}

// buildLLMPrompt 构建LLM提示词
func (s *IntentRecognitionService) buildLLMPrompt(text string, intents []*entity.Intent) string {
	intentDescriptions := ""
	for _, intent := range intents {
		examples, _ := intent.GetExamples()
		intentDescriptions += fmt.Sprintf("- %s: %s\n", intent.IntentName, intent.Description)
		if len(examples) > 0 {
			intentDescriptions += fmt.Sprintf("  示例: %v\n", examples)
		}
	}

	prompt := fmt.Sprintf(`你是一个意图识别助手。请分析用户输入，识别最匹配的意图。

用户输入: %s

可用意图列表:
%s

请以JSON格式返回识别结果:
{
  "intent_name": "意图名称",
  "confidence": 0.95,
  "reasoning": "识别理由"
}
`, text, intentDescriptions)

	return prompt
}

// intentMatchScore 意图匹配得分
type intentMatchScore struct {
	IntentID string
	Scores   []float64
	AvgScore float64
}

// average 计算平均值
func average(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// sortAndLimit 排序并限制数量
func sortAndLimit(scores []float64, limit int) []float64 {
	sort.Float64s(scores)
	if len(scores) > limit {
		return scores[:limit]
	}
	return scores
}
