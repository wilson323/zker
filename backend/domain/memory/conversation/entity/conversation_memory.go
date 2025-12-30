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

package entity

import (
	"encoding/json"
	"time"
)

// MemoryType 记忆类型
type MemoryType string

const (
	MemoryTypeSummary    MemoryType = "SUMMARY"    // 对话摘要
	MemoryTypeEntity     MemoryType = "ENTITY"     // 实体提取
	MemoryTypePreference MemoryType = "PREFERENCE" // 用户偏好
	MemoryTypeEvent      MemoryType = "EVENT"      // 事件记录
)

// ConversationMemory 对话记忆实体
type ConversationMemory struct {
	ID              int64              `json:"id"`
	MemoryID        string             `json:"memory_id"`
	TenantID        string             `json:"tenant_id"`
	UserID          string             `json:"user_id"`
	ConversationID  string             `json:"conversation_id"`
	MemoryType      MemoryType         `json:"memory_type"`
	Content         string             `json:"content"`
	Embedding       []float32          `json:"embedding"`        // 向量嵌入(1536维)
	ImportanceScore float64            `json:"importance_score"` // 重要性评分(0-1)
	AccessCount     int                `json:"access_count"`
	LastAccessedAt  time.Time          `json:"last_accessed_at"`
	ExpiresAt       *time.Time         `json:"expires_at"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Metadata        map[string]interface{} `json:"metadata"` // 元数据
}

// GetTenantID 获取租户ID(用于多租户隔离)
func (m *ConversationMemory) GetTenantID() string {
	return m.TenantID
}

// HasTenantID 检查是否有租户ID
func (m *ConversationMemory) HasTenantID() bool {
	return m.TenantID != ""
}

// IsExpired 检查记忆是否过期
func (m *ConversationMemory) IsExpired() bool {
	if m.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*m.ExpiresAt)
}

// ShouldCompress 判断是否需要压缩(重要性低且访问次数少)
func (m *ConversationMemory) ShouldCompress() bool {
	return m.ImportanceScore < 0.3 && m.AccessCount < 5
}

// MemoryWithScore 带相似度分数的记忆
type MemoryWithScore struct {
	Memory *ConversationMemory `json:"memory"`
	Score  float64              `json:"score"` // 相似度分数(0-1)
}

// SummaryResult 对话摘要结果
type SummaryResult struct {
	OriginalCount     int       `json:"original_count"`
	CompressedCount   int       `json:"compressed_count"`
	Summary           string    `json:"summary"`
	KeyFacts          []string  `json:"key_facts"`
	CompressionRatio  float64   `json:"compression_ratio"`
	NewMemoryID       string    `json:"new_memory_id"`
}

// ExtractEntityResult 实体提取结果
type ExtractEntityResult struct {
	Entities      []Entity `json:"entities"`
	Confidence    float64  `json:"confidence"`
	ExtractedAt   time.Time `json:"extracted_at"`
}

// Entity 提取的实体
type Entity struct {
	EntityType   string   `json:"entity_type"`   // 实体类型: PERSON, ORG, LOCATION等
	EntityValue  string   `json:"entity_value"`  // 实体值
	Attributes   map[string]string `json:"attributes"` // 实体属性
	Confidence   float64  `json:"confidence"`    // 置信度
}

// ExtractPreferenceResult 偏好提取结果
type ExtractPreferenceResult struct {
	Preferences []Preference `json:"preferences"`
	Confidence  float64      `json:"confidence"`
	ExtractedAt time.Time    `json:"extracted_at"`
}

// Preference 用户偏好
type Preference struct {
	PreferenceType string `json:"preference_type"` // 偏好类型: communication_style, response_length等
	PreferenceValue string `json:"preference_value"` // 偏好值
	Source       string `json:"source"`        // 来源: explicit, implicit, inferred
	Confidence   float64 `json:"confidence"`   // 置信度
}

// GetMetadataJSON 获取元数据JSON
func (m *ConversationMemory) GetMetadataJSON() (string, error) {
	if m.Metadata == nil {
		return "{}", nil
	}
	data, err := json.Marshal(m.Metadata)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetMetadataFromJSON 从JSON设置元数据
func (m *ConversationMemory) SetMetadataFromJSON(jsonStr string) error {
	if jsonStr == "" || jsonStr == "{}" {
		m.Metadata = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), &m.Metadata)
}
