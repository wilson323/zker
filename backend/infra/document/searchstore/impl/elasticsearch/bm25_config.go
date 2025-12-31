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

package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

// BM25IndexConfig BM25索引配置
// 用于优化全文检索效果
type BM25IndexConfig struct {
	// 索引名称
	IndexName string
	// 分片数
	NumberOfShards int
	// 副本数
	NumberOfReplicas int
	// 分析器配置
	Analyzer *AnalyzerConfig
}

// AnalyzerConfig 分析器配置
type AnalyzerConfig struct {
	// 分词器类型: ik_smart (粗粒度), ik_max_word (细粒度)
	Tokenizer string
	// 是否使用IK分词器(中文优化)
	UseIK bool
}

// CreateBM25Index 创建BM25优化的索引
func (m *ElasticsearchSearchStore) CreateBM25Index(ctx context.Context, config *BM25IndexConfig) error {
	if config == nil {
		return fmt.Errorf("config is required")
	}
	if config.IndexName == "" {
		return fmt.Errorf("index name is required")
	}

	// 设置默认值
	if config.NumberOfShards <= 0 {
		config.NumberOfShards = 3
	}
	if config.NumberOfReplicas <= 0 {
		config.NumberOfReplicas = 1
	}
	if config.Analyzer == nil {
		config.Analyzer = &AnalyzerConfig{
			Tokenizer: "ik_max_word",
			UseIK:     true,
		}
	}

	// 构建索引映射
	mapping := m.buildBM25Mapping(config.Analyzer)

	// 构建索引设置
	settings := map[string]interface{}{
		"number_of_shards":   config.NumberOfShards,
		"number_of_replicas": config.NumberOfReplicas,
	}

	if config.Analyzer.UseIK {
		// 添加IK分词器配置
		settings["analysis"] = map[string]interface{}{
			"analyzer": map[string]interface{}{
				"ik_max_word": map[string]interface{}{
					"type":      "custom",
					"tokenizer": "ik_max_word",
				},
				"ik_smart": map[string]interface{}{
					"type":      "custom",
					"tokenizer": "ik_smart",
				},
			},
		}
	}

	// 构建完整索引定义
	indexDefinition := map[string]interface{}{
		"settings": settings,
		"mappings": mapping,
	}

	// 创建索引
	err := m.createIndex(ctx, config.IndexName, indexDefinition)
	if err != nil {
		return fmt.Errorf("create BM25 index failed: %w", err)
	}

	logs.CtxInfof(ctx, "[Elasticsearch] created BM25 index: %s", config.IndexName)

	return nil
}

// buildBM25Mapping 构建BM25映射
func (m *ElasticsearchSearchStore) buildBM25Mapping(config *AnalyzerConfig) map[string]interface{} {
	analyzer := "standard"
	tokenizer := "standard"

	if config.UseIK {
		analyzer = config.Tokenizer
		if analyzer == "" {
			analyzer = "ik_max_word"
		}
	}

	// 构建映射
	mapping := map[string]interface{}{
		"properties": map[string]interface{}{
			// 文档ID (必需)
			"document_id": map[string]interface{}{
				"type": "keyword",
			},
			// 知识库ID (必需)
			"knowledge_id": map[string]interface{}{
				"type": "keyword",
			},
			// 标题 (需要高权重)
			"title": map[string]interface{}{
				"type":     "text",
				"analyzer": analyzer,
				"search_analyzer": analyzer,
				"boost":    2.0, // 标题匹配权重翻倍
				"fields": map[string]interface{}{
					// 精确匹配字段
					"keyword": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
			// 内容 (主要检索字段)
			"content": map[string]interface{}{
				"type":            "text",
				"analyzer":        analyzer,
				"search_analyzer": analyzer,
				"similarity":      "BM25", // 显式使用BM25相似度算法
			},
			// 元数据 (用于过滤)
			"metadata": map[string]interface{}{
				"type": "object",
			},
			// 创建时间 (用于排序)
			"created_at": map[string]interface{}{
				"type": "date",
			},
			// 更新时间 (用于排序)
			"updated_at": map[string]interface{}{
				"type": "date",
			},
			// 切片序列号 (用于排序)
			"sequence": map[string]interface{}{
				"type": "long",
			},
			// 命中次数 (用于排序和热度)
			"hit_count": map[string]interface{}{
				"type": "long",
			},
			// 字符数 (用于过滤)
			"char_count": map[string]interface{}{
				"type": "long",
			},
		},
	}

	return mapping
}

// createIndex 创建索引
func (m *ElasticsearchSearchStore) createIndex(ctx context.Context, indexName string, definition map[string]interface{}) error {
	client, err := m.getClient()
	if err != nil {
		return fmt.Errorf("get client failed: %w", err)
	}

	// 检查索引是否存在
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("check index exists failed: %w", err)
	}

	if exists {
		logs.CtxInfof(ctx, "[Elasticsearch] index %s already exists, skipping creation", indexName)
		return nil
	}

	// 创建索引
	createResp, err := client.CreateIndex(indexName).BodyJson(definition).Do(ctx)
	if err != nil {
		return fmt.Errorf("create index request failed: %w", err)
	}

	if !createResp.Acknowledged {
		return fmt.Errorf("create index not acknowledged")
	}

	return nil
}

// UpdateIndexSettings 更新索引设置
func (m *ElasticsearchSearchStore) UpdateIndexSettings(ctx context.Context, indexName string, settings map[string]interface{}) error {
	client, err := m.getClient()
	if err != nil {
		return fmt.Errorf("get client failed: %w", err)
	}

	_, err = client.IndexPutSettings(indexName).BodyJson(settings).Do(ctx)
	if err != nil {
		return fmt.Errorf("update index settings failed: %w", err)
	}

	logs.CtxInfof(ctx, "[Elasticsearch] updated index settings: %s", indexName)

	return nil
}

// GetBM25IndexSettings 获取BM25索引的推荐设置
func GetBM25IndexSettings() map[string]interface{} {
	return map[string]interface{}{
		"index": map[string]interface{}{
			// BM25参数调优
			"similarity": map[string]interface{}{
				"custom_bm25": map[string]interface{}{
					"type":      "BM25",
					"k1":        1.2, // 词频饱和度参数 (1.2-2.0)
					"b":         0.75, // 长度归一化参数 (0.75)
				},
			},
		},
	}
}

// OptimizeIndexForBM25 优化现有索引为BM25
func (m *ElasticsearchSearchStore) OptimizeIndexForBM25(ctx context.Context, indexName string) error {
	client, err := m.getClient()
	if err != nil {
		return fmt.Errorf("get client failed: %w", err)
	}

	// 1. 获取当前索引设置
	settingsService := client.IndexGetSettings(indexName)
	settingsResp, err := settingsService.Do(ctx)
	if err != nil {
		return fmt.Errorf("get index settings failed: %w", err)
	}

	// 2. 检查是否需要更新BM25参数
	needsUpdate := false
	currentSettings := settingsResp[indexName].Settings

	// 检查similarity配置
	if _, ok := currentSettings["index.similarity.custom_bm25"]; !ok {
		needsUpdate = true
	}

	if needsUpdate {
		// 3. 更新索引设置
		recommendedSettings := GetBM25IndexSettings()
		err = m.UpdateIndexSettings(ctx, indexName, recommendedSettings)
		if err != nil {
			return fmt.Errorf("update index settings failed: %w", err)
		}

		logs.CtxInfof(ctx, "[Elasticsearch] optimized index %s for BM25", indexName)
	} else {
		logs.CtxInfof(ctx, "[Elasticsearch] index %s already optimized for BM25", indexName)
	}

	return nil
}

// BuildBM25SearchQuery 构建BM25搜索查询
func (m *ElasticsearchSearchStore) BuildBM25SearchQuery(query string, topK int, filter map[string]interface{}) map[string]interface{} {
	// 构建Bool查询
	boolQuery := map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []map[string]interface{}{
				// 标题匹配 (boost=2.0)
				{
					"multi_match": map[string]interface{}{
						"query":  query,
						"fields": []string{"title^2", "title.keyword^3"},
						"type":   "best_fields",
						"boost":  2.0,
					},
				},
				// 内容匹配
				{
					"match": map[string]interface{}{
						"content": map[string]interface{}{
							"query": query,
							"boost": 1.0,
						},
					},
				},
			},
			"minimum_should_match": 1,
		},
	}

	// 添加过滤条件
	if len(filter) > 0 {
		boolQuery["bool"].(map[string]interface{})["filter"] = []map[string]interface{}{
			{
				"term": filter,
			},
		}
	}

	// 构建完整查询
	searchQuery := map[string]interface{}{
		"query": boolQuery,
		"size":  topK,
		"sort": []map[string]interface{}{
			// 按分数排序
			{"_score": map[string]interface{}{"order": "desc"}},
			// 分数相同时按命中次数排序
			{"hit_count": map[string]interface{}{"order": "desc"}},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{
					"fragment_size":       100,
					"number_of_fragments": 3,
				},
				"content": map[string]interface{}{
					"fragment_size":       200,
					"number_of_fragments": 3,
				},
			},
		},
	}

	return searchQuery
}

// AnalyzeQuery 分析查询(用于调试)
func (m *ElasticsearchSearchStore) AnalyzeQuery(ctx context.Context, indexName, query, analyzer string) (map[string]interface{}, error) {
	client, err := m.getClient()
	if err != nil {
		return nil, fmt.Errorf("get client failed: %w", err)
	}

	if analyzer == "" {
		analyzer = "ik_max_word"
	}

	// 执行分析
	analyzeResp, err := client.IndexAnalyze().
		Index(indexName).
		Analyzer(analyzer).
		Text(query).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("analyze query failed: %w", err)
	}

	// 解析结果
	result := map[string]interface{}{
		"original": query,
		"analyzer": analyzer,
		"tokens":   analyzeResp.Tokens,
	}

	return result, nil
}

// GetIndexStats 获取索引统计信息
func (m *ElasticsearchSearchStore) GetIndexStats(ctx context.Context, indexName string) (map[string]interface{}, error) {
	client, err := m.getClient()
	if err != nil {
		return nil, fmt.Errorf("get client failed: %w", err)
	}

	statsResp, err := client.Index(indexName).Stats().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("get index stats failed: %w", err)
	}

	// 提取关键统计信息
	stats := map[string]interface{}{
		"document_count": statsResp.Indexes[indexName].Primaries.Docs.Count,
		"store_size":     statsResp.Indexes[indexName].Primaries.Store.SizeInBytes,
		"indexing_time":  statsResp.Indexes[indexName].Primaries.Indexing.IndexTimeInMillis,
		"search_time":    statsResp.Indexes[indexName].Primaries.Search.QueryTimeInMillis,
	}

	return stats, nil
}
