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
package impl

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/api/model/admin/config"
	"github.com/coze-dev/coze-studio/backend/infra/document/rerank"
	"github.com/coze-dev/coze-studio/backend/infra/document/rerank/impl/crossencoder"
	"github.com/coze-dev/coze-studio/backend/infra/document/rerank/impl/rrf"
	"github.com/coze-dev/coze-studio/backend/infra/document/rerank/impl/vikingdb"
)

type Reranker = rerank.Reranker

func New(conf *config.KnowledgeConfig) Reranker {
	switch conf.RerankConfig.Type {
	case config.RerankType_VikingDB:
		return vikingdb.NewReranker(conf.RerankConfig.VikingdbConfig)
	case config.RerankType_RRF:
		return rrf.NewRRFReranker(0)
	case config.RerankType_CrossEncoder:
		// CrossEncoder Reranker (BGE-Reranker / Jina / Cohere)
		return crossencoder.NewSimpleCrossEncoderReranker()
	default:
		return rrf.NewRRFReranker(0)
	}
}

// NewCrossEncoderReranker 创建CrossEncoder Reranker
func NewCrossEncoderReranker(endpoint, model, apiKey string) (Reranker, error) {
	cfg := &crossencoder.CrossEncoderRerankerConfig{
		Endpoint: endpoint,
		Model:    model,
		APIKey:   apiKey,
	}
	reranker, err := crossencoder.NewCrossEncoderReranker(cfg)
	if err != nil {
		return nil, fmt.Errorf("create cross encoder reranker failed: %w", err)
	}
	return reranker, nil
}
