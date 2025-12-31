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

package crossencoder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/infra/document/rerank"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// CrossEncoderRerankerConfig CrossEncoder Reranker配置
type CrossEncoderRerankerConfig struct {
	// Reranker服务端点
	Endpoint string
	// 模型名称(如bge-reranker-v2-m3)
	Model string
	// 批处理大小
	BatchSize int
	// 超时时间
	Timeout time.Duration
	// API密钥(如果需要)
	APIKey string
}

// CrossEncoderReranker 交叉编码器Reranker
// 支持BGE-Reranker、Cohere Rerank等模型
type CrossEncoderReranker struct {
	config *CrossEncoderRerankerConfig
	client *http.Client
}

// NewCrossEncoderReranker 创建CrossEncoderReranker
func NewCrossEncoderReranker(config *CrossEncoderRerankerConfig) (rerank.Reranker, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}
	if config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	// 设置默认值
	if config.Model == "" {
		config.Model = "bge-reranker-v2-m3"
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 32
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	return &CrossEncoderReranker{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}, nil
}

// Rerank 执行二次排序
func (r *CrossEncoderReranker) Rerank(ctx context.Context, req *rerank.Request) (*rerank.Response, error) {
	if req == nil || req.Data == nil || len(req.Data) == 0 {
		return nil, fmt.Errorf("invalid request: no data provided")
	}

	// 展平所有数据
	var allDocs []*rerank.Data
	for _, resultList := range req.Data {
		allDocs = append(allDocs, resultList...)
	}

	if len(allDocs) == 0 {
		return &rerank.Response{SortedData: []*rerank.Data{}}, nil
	}

	// 调用Reranker API
	scores, err := r.callRerankerAPI(ctx, req.Query, allDocs)
	if err != nil {
		logs.CtxWarnf(ctx, "[CrossEncoderReranker] API call failed: %v, falling back to original scores", err)
		// 降级:使用原始分数
		return r.fallbackRerank(allDocs, req.TopN), nil
	}

	// 更新文档分数
	for i, doc := range allDocs {
		if i < len(scores) {
			doc.Score = scores[i]
		}
	}

	// 按新分数排序
	sort.Slice(allDocs, func(i, j int) bool {
		return allDocs[i].Score > allDocs[j].Score
	})

	// 应用TopN
	topN := int64(len(allDocs))
	if req.TopN != nil && *req.TopN > 0 && *req.TopN < topN {
		topN = *req.TopN
	}

	logs.CtxInfof(ctx, "[CrossEncoderReranker] reranked %d documents to %d", len(allDocs), topN)

	return &rerank.Response{
		SortedData: allDocs[:topN],
	}, nil
}

// callRerankerAPI 调用Reranker API
func (r *CrossEncoderReranker) callRerankerAPI(ctx context.Context, query string, docs []*rerank.Data) ([]float64, error) {
	// 方案1: 调用本地BGE Reranker服务
	// 方案2: 调用Jina Reranker API
	// 方案3: 调用Cohere Rerank API

	// 这里实现一个通用的HTTP客户端,支持多种API
	switch r.config.Model {
	case "bge-reranker-v2-m3", "bge-reranker-large":
		return r.callBGEReranker(ctx, query, docs)
	case "jina-reranker-v1-base-en":
		return r.callJinaReranker(ctx, query, docs)
	case "cohere-rerank-v2", "cohere-rerank-english-v2.0":
		return r.callCohereReranker(ctx, query, docs)
	default:
		return r.callBGEReranker(ctx, query, docs)
	}
}

// callBGEReranker 调用BGE Reranker服务
// BGE Reranker通常部署为本地服务或通过API调用
func (r *CrossEncoderReranker) callBGEReranker(ctx context.Context, query string, docs []*rerank.Data) ([]float64, error) {
	// 准备请求体
	type Document struct {
		Text string `json:"text"`
	}

	requestBody := map[string]interface{}{
		"model":  r.config.Model,
		"query":  query,
		"documents": make([]string, len(docs)),
	}

	for i, doc := range docs {
		requestBody["documents"].([]string)[i] = doc.Document.Content
	}

	// 发送HTTP请求
	reqData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.config.Endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if r.config.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+r.config.APIKey)
	}

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response struct {
		Results []struct {
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relevance_score"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	// 提取分数
	scores := make([]float64, len(docs))
	for _, result := range response.Results {
		if result.Index >= 0 && result.Index < len(scores) {
			scores[result.Index] = result.RelevanceScore
		}
	}

	return scores, nil
}

// callJinaReranker 调用Jina Reranker API
func (r *CrossEncoderReranker) callJinaReranker(ctx context.Context, query string, docs []*rerank.Data) ([]float64, error) {
	// Jina Reranker API: https://jina.ai/reranker
	endpoint := "https://api.jina.ai/v1/rerank"

	type Document struct {
		Text string `json:"text"`
	}

	requestBody := map[string]interface{}{
		"model":      r.config.Model,
		"query":      query,
		"documents":  make([]Document, len(docs)),
		"top_n":      len(docs),
	}

	for i, doc := range docs {
		requestBody["documents"].([]Document)[i] = Document{Text: doc.Document.Content}
	}

	reqData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+r.config.APIKey)

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Results []struct {
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relevance_score"`
			Document       struct {
				Text string `json:"text"`
			} `json:"document"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	scores := make([]float64, len(docs))
	for _, result := range response.Results {
		if result.Index >= 0 && result.Index < len(scores) {
			scores[result.Index] = result.RelevanceScore
		}
	}

	return scores, nil
}

// callCohereReranker 调用Cohere Rerank API
func (r *CrossEncoderReranker) callCohereReranker(ctx context.Context, query string, docs []*rerank.Data) ([]float64, error) {
	// Cohere Rerank API: https://docs.cohere.com/reference/rerank
	endpoint := "https://api.cohere.ai/v1/rerank"

	type Document struct {
		Text string `json:"text"`
	}

	requestBody := map[string]interface{}{
		"model":      r.config.Model,
		"query":      query,
		"documents":  make([]Document, len(docs)),
		"top_n":      len(docs),
		"return_documents": false,
	}

	for i, doc := range docs {
		requestBody["documents"].([]Document)[i] = Document{Text: doc.Document.Content}
	}

	reqData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+r.config.APIKey)
	httpReq.Header.Set("X-Client-Name", "coze-studio")

	resp, err := r.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Results []struct {
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relevance_score"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	scores := make([]float64, len(docs))
	for _, result := range response.Results {
		if result.Index >= 0 && result.Index < len(scores) {
			scores[result.Index] = result.RelevanceScore
		}
	}

	return scores, nil
}

// fallbackRerank 降级排序(使用原始分数)
func (r *CrossEncoderReranker) fallbackRerank(docs []*rerank.Data, topN *int64) *rerank.Response {
	// 按原始分数排序
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Score > docs[j].Score
	})

	n := int64(len(docs))
	if topN != nil && *topN > 0 && *topN < n {
		n = *topN
	}

	return &rerank.Response{
		SortedData: docs[:n],
	}
}

// ValidateEndpoint 验证端点是否可用
func (r *CrossEncoderReranker) ValidateEndpoint(ctx context.Context) error {
	endpointURL, err := url.Parse(r.config.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}

	// 发送健康检查请求
	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL.String(), nil)
	if err != nil {
		return fmt.Errorf("create health check request failed: %w", err)
	}

	if r.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+r.config.APIKey)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}

// SimpleCrossEncoderReranker 简化版CrossEncoder Reranker
// 用于不需要外部API的场景,使用基于规则的相关性计算
type SimpleCrossEncoderReranker struct {
	// 配置参数
	titleWeight    float64 // 标题权重
	contentWeight  float64 // 内容权重
	originalWeight float64 // 原始分数权重
}

// NewSimpleCrossEncoderReranker 创建简化版CrossEncoder Reranker
func NewSimpleCrossEncoderReranker() rerank.Reranker {
	return &SimpleCrossEncoderReranker{
		titleWeight:    0.4,
		contentWeight:  0.3,
		originalWeight: 0.3,
	}
}

// Rerank 执行简化版二次排序
func (r *SimpleCrossEncoderReranker) Rerank(ctx context.Context, req *rerank.Request) (*rerank.Response, error) {
	if req == nil || req.Data == nil || len(req.Data) == 0 {
		return nil, fmt.Errorf("invalid request: no data provided")
	}

	// 展平所有数据
	var allDocs []*rerank.Data
	for _, resultList := range req.Data {
		allDocs = append(allDocs, resultList...)
	}

	if len(allDocs) == 0 {
		return &rerank.Response{SortedData: []*rerank.Data{}}, nil
	}

	// 计算新的相关性分数
	queryTokens := tokenize(req.Query)
	for _, doc := range allDocs {
		newScore := r.calculateRelevance(queryTokens, doc)
		doc.Score = newScore
	}

	// 按新分数排序
	sort.Slice(allDocs, func(i, j int) bool {
		return allDocs[i].Score > allDocs[j].Score
	})

	// 应用TopN
	topN := int64(len(allDocs))
	if req.TopN != nil && *req.TopN > 0 && *req.TopN < topN {
		topN = *req.TopN
	}

	logs.CtxInfof(ctx, "[SimpleCrossEncoderReranker] reranked %d documents to %d", len(allDocs), topN)

	return &rerank.Response{
		SortedData: allDocs[:topN],
	}, nil
}

// calculateRelevance 计算相关性分数
func (r *SimpleCrossEncoderReranker) calculateRelevance(queryTokens []string, doc *rerank.Data) float64 {
	// 提取文档的标题和内容
	var title, content string
	if doc.Document != nil {
		if titleVal, ok := doc.Document.MetaData["title"]; ok {
			title = toString(titleVal)
		}
		content = doc.Document.Content
	}

	// 计算标题匹配分数
	titleTokens := tokenize(title)
	titleScore := calculateOverlap(queryTokens, titleTokens)

	// 计算内容匹配分数
	contentTokens := tokenize(content)
	contentScore := calculateOverlap(queryTokens, contentTokens)

	// 结合原始分数
	originalScore := doc.Score

	// 加权组合
	finalScore := titleScore*r.titleWeight + contentScore*r.contentWeight + originalScore*r.originalWeight

	return finalScore
}

// tokenize 分词
func tokenize(text string) []string {
	// 简化分词:按空格和标点符号分割
	// 实际应该使用更复杂的分词器(如jieba)
	tokens := make([]string, 0)
	currentToken := ""

	for _, ch := range text {
		if isWordChar(ch) {
			currentToken += string(ch)
		} else {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
		}
	}

	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

// isWordChar 判断是否为单词字符
func isWordChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || (ch >= 0x4E00 && ch <= 0x9FFF)
}

// calculateOverlap 计算重叠度
func calculateOverlap(tokens1, tokens2 []string) float64 {
	if len(tokens1) == 0 || len(tokens2) == 0 {
		return 0.0
	}

	// 创建token频率映射
	freq := make(map[string]int)
	for _, token := range tokens1 {
		freq[token]++
	}

	// 计算匹配的token数
	matches := 0
	for _, token := range tokens2 {
		if freq[token] > 0 {
			matches++
			freq[token]--
		}
	}

	// 计算重叠度分数
	return float64(matches) / float64(len(tokens1))
}

// toString 安全转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
