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

package model

import "time"

// ModelType 模型类型
type ModelType int32

const (
	ModelType_LLM           ModelType = 0
	ModelType_TextEmbedding ModelType = 1
	ModelType_Rerank        ModelType = 2
)

// String 返回模型类型字符串
func (t ModelType) String() string {
	switch t {
	case ModelType_LLM:
		return "LLM"
	case ModelType_TextEmbedding:
		return "TextEmbedding"
	case ModelType_Rerank:
		return "Rerank"
	}
	return "UNKNOWN"
}

// ModelClass 模型分类
type ModelClass string

const (
	ModelClass_SEED      ModelClass = "SEED"
	ModelClass_GPT       ModelClass = "GPT"
	ModelClass_Claude    ModelClass = "Claude"
	ModelClass_DeekSeek  ModelClass = "DeekSeek"
	ModelClass_Gemini    ModelClass = "Gemini"
	ModelClass_Llama     ModelClass = "Llama"
	ModelClass_QWen      ModelClass = "QWen"
)

// ModelConfig 跨领域的模型配置
//
// 该模型用于LLM配置管理，在domain、application、bizpkg层共享
type ModelConfig struct {
	// 基础信息
	ModelID   int64
	ModelName string
	ModelType ModelType
	IsActive  bool

	// 提供商信息
	ProviderID     int64
	ProviderName   string
	ProviderClass  ModelClass
	ProviderAPIKey string

	// 连接配置
	APIEndpoint string
	TimeoutSec  int32

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// ModelProvider 模型提供商信息
type ModelProvider struct {
	ProviderID    int64
	ProviderName  string
	ModelClass    ModelClass
	APIEndpoint   string
	APIKey        string
	IsEnabled     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Validate 验证模型配置有效性
func (c *ModelConfig) Validate() error {
	if c.ModelID <= 0 {
		return ErrInvalidModelID
	}
	if c.ModelName == "" {
		return ErrInvalidModelName
	}
	if c.ProviderID <= 0 {
		return ErrInvalidProviderID
	}
	if c.APIEndpoint == "" {
		return ErrInvalidAPIEndpoint
	}
	return nil
}

// Validate 验证提供商配置有效性
func (p *ModelProvider) Validate() error {
	if p.ProviderID <= 0 {
		return ErrInvalidProviderID
	}
	if p.ProviderName == "" {
		return ErrInvalidProviderName
	}
	if p.ModelClass == "" {
		return ErrInvalidModelClass
	}
	if p.APIEndpoint == "" {
		return ErrInvalidAPIEndpoint
	}
	return nil
}
