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

// Package model 提供跨领域的通用数据模型
//
// 这些模型在 domain、application、bizpkg 层共享使用，
// 避免直接依赖 api/model 层，符合 DDD 分层架构原则。
package model

import "time"

// BotConfig 跨领域的Bot配置模型
//
// 该模型封装了Bot的核心配置信息，用于在domain层和application层之间传递数据。
// 不包含API层的特定字段，保持领域模型的纯净性。
type BotConfig struct {
	// 基础信息
	BotID       string
	BotName     string
	BotDesc     string
	BotIcon     string
	IsPublic    bool
	IsActive    bool
	PublishTime *time.Time

	// 模型配置
	ModelID   *int64
	ModelName string

	// LLM参数
	Temperature      float32
	MaxTokens        int32
	TopP             float32
	FrequencyPenalty float32
	PresencePenalty  float32

	// 高级配置
	ThinkingType     int32  // 0:Default, 1:Enable, 2:Disable, 3:Auto
	EnableSearch     bool
	EnableMemory     bool
	EnableKnowledge  bool
	EnableWorkflow   bool

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// SimpleModelInfo 简化的模型信息（用于 BotConfig）
type SimpleModelInfo struct {
	ModelID      *int64
	ModelName    string
	Temperature  float32
	MaxTokens    int32
	TopP         float32
	EnableSearch bool
}

// ToModelInfo 转换为 SimpleModelInfo
func (c *BotConfig) ToModelInfo() *SimpleModelInfo {
	if c == nil {
		return nil
	}
	return &SimpleModelInfo{
		ModelID:      c.ModelID,
		ModelName:    c.ModelName,
		Temperature:  c.Temperature,
		MaxTokens:    c.MaxTokens,
		TopP:         c.TopP,
		EnableSearch: c.EnableSearch,
	}
}

// Validate 验证配置有效性
func (c *BotConfig) Validate() error {
	if c.BotID == "" {
		return ErrInvalidBotID
	}
	if c.ModelID == nil && c.ModelName == "" {
		return ErrInvalidModelConfig
	}
	if c.Temperature < 0 || c.Temperature > 2 {
		return ErrInvalidTemperature
	}
	if c.MaxTokens < 0 {
		return ErrInvalidMaxTokens
	}
	if c.TopP < 0 || c.TopP > 1 {
		return ErrInvalidTopP
	}
	return nil
}

// ShortMemoryPolicy 短期记忆策略
type ShortMemoryPolicy struct {
	Limit        int32   // 记忆条数限制
	Threshold    float64 // 相似度阈值
	Strategy     string  // 策略类型
	MemoryWindow int32   // 记忆窗口大小
}

// BotModelInfo Bot LLM 模型配置信息（完整版本，兼容 bot_common.ModelInfo）
//
// 该结构用于描述 Bot 使用的 LLM 模型配置参数，
// 包含温度、最大 Token、Top P 等生成参数。
type BotModelInfo struct {
	// Model ID
	ModelId *int64

	// Temperature, model output randomness, the larger the value, the more random, the smaller the more conservative (0-1]
	Temperature *float64

	// Maximum Token Reply
	MaxTokens *int32

	// Another model's output randomness, the larger the value, the more random [0, 1]
	TopP *float64

	// Frequency penalty, adjust the frequency of words in the generated content, the fewer positive words are [-1.0, 1.0]
	FrequencyPenalty *float64

	// There is a penalty, adjust the frequency of new words in the generated content, avoid repeating words with positive values, and use new words [-1.0, 1.0]
	PresencePenalty *float64

	// contextual policy
	ShortMemoryPolicy *ShortMemoryPolicy

	// When generating, sample the size of the candidate set
	TopK *int32

	// Stop sequence
	Stop *string

	// Whether to return the prompt
	ReturnPrompt *bool

	// Return number of results
	N *int32

	// Whether to enable search tool
	EnableSearch *bool
}

// GetModelId 获取模型 ID，如果未设置返回 0
func (m *BotModelInfo) GetModelId() int64 {
	if m == nil || m.ModelId == nil {
		return 0
	}
	return *m.ModelId
}

// GetTemperature 获取温度参数，如果未设置返回默认值 0.7
func (m *BotModelInfo) GetTemperature() float64 {
	if m == nil || m.Temperature == nil {
		return 0.7
	}
	return *m.Temperature
}

// GetMaxTokens 获取最大 Token 数，如果未设置返回默认值 2048
func (m *BotModelInfo) GetMaxTokens() int32 {
	if m == nil || m.MaxTokens == nil {
		return 2048
	}
	return *m.MaxTokens
}

// GetTopP 获取 TopP 参数，如果未设置返回默认值 1.0
func (m *BotModelInfo) GetTopP() float64 {
	if m == nil || m.TopP == nil {
		return 1.0
	}
	return *m.TopP
}

// GetFrequencyPenalty 获取频率惩罚参数，如果未设置返回默认值 0.0
func (m *BotModelInfo) GetFrequencyPenalty() float64 {
	if m == nil || m.FrequencyPenalty == nil {
		return 0.0
	}
	return *m.FrequencyPenalty
}

// GetPresencePenalty 获取存在惩罚参数，如果未设置返回默认值 0.0
func (m *BotModelInfo) GetPresencePenalty() float64 {
	if m == nil || m.PresencePenalty == nil {
		return 0.0
	}
	return *m.PresencePenalty
}

// GetTopK 获取 TopK 参数，如果未设置返回默认值 0
func (m *BotModelInfo) GetTopK() int32 {
	if m == nil || m.TopK == nil {
		return 0
	}
	return *m.TopK
}

// IsSearchEnabled 检查是否启用搜索工具
func (m *BotModelInfo) IsSearchEnabled() bool {
	if m == nil || m.EnableSearch == nil {
		return false
	}
	return *m.EnableSearch
}

// Validate 验证模型配置有效性
func (m *BotModelInfo) Validate() error {
	if m.GetModelId() <= 0 {
		return ErrInvalidModelID
	}
	temp := m.GetTemperature()
	if temp < 0 || temp > 2 {
		return ErrInvalidTemperature
	}
	maxTokens := m.GetMaxTokens()
	if maxTokens < 0 {
		return ErrInvalidMaxTokens
	}
	topP := m.GetTopP()
	if topP < 0 || topP > 1 {
		return ErrInvalidTopP
	}
	return nil
}

// 错误定义
var (
	ErrInvalidBotID         = &ModelError{Message: "bot_id cannot be empty"}
	ErrInvalidModelConfig   = &ModelError{Message: "model config is invalid"}
	ErrInvalidModelID       = &ModelError{Message: "model_id must be positive"}
	ErrInvalidModelName     = &ModelError{Message: "model_name cannot be empty"}
	ErrInvalidTemperature   = &ModelError{Message: "temperature must be between 0 and 2"}
	ErrInvalidMaxTokens     = &ModelError{Message: "max_tokens must be non-negative"}
	ErrInvalidTopP          = &ModelError{Message: "top_p must be between 0 and 1"}
	ErrInvalidProviderID    = &ModelError{Message: "provider_id must be positive"}
	ErrInvalidProviderName  = &ModelError{Message: "provider_name cannot be empty"}
	ErrInvalidModelClass    = &ModelError{Message: "model_class cannot be empty"}
	ErrInvalidAPIEndpoint   = &ModelError{Message: "api_endpoint cannot be empty"}
)

// ModelError 模型相关错误
type ModelError struct {
	Message string
}

func (e *ModelError) Error() string {
	return e.Message
}
