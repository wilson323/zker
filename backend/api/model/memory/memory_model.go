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

package memory

// StoreConversationMemoryRequest 存储对话记忆请求
type StoreConversationMemoryRequest struct {
	ConversationID string                 `json:"conversation_id" binding:"required"`
	UserID         string                 `json:"user_id" binding:"required"`
	TenantID       string                 `json:"tenant_id" binding:"required"`
	MemoryType     string                 `json:"memory_type" binding:"required,oneof=SUMMARY ENTITY PREFERENCE EVENT"`
	Content        string                 `json:"content" binding:"required"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// StoreConversationMemoryResponse 存储对话记忆响应
type StoreConversationMemoryResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      *MemoryData `json:"data,omitempty"`
	Timestamp int64  `json:"timestamp"`
	TraceID   string `json:"trace_id,omitempty"`
}

// MemoryData 记忆数据
type MemoryData struct {
	MemoryID    string `json:"memory_id"`
	VectorID    string `json:"vector_id"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
}

// RetrieveMemoriesRequest 检索记忆请求
type RetrieveMemoriesRequest struct {
	UserID    string   `json:"user_id" binding:"required"`
	Query     string   `json:"query" binding:"required"`
	TopK      int      `json:"top_k" binding:"min=1,max=20"`
	ScoreThreshold float64 `json:"score_threshold" binding:"min=0,max=1"`
}

// RetrieveMemoriesResponse 检索记忆响应
type RetrieveMemoriesResponse struct {
	Code      int                `json:"code"`
	Message   string             `json:"message"`
	Data      *MemoriesResult    `json:"data,omitempty"`
	Timestamp int64              `json:"timestamp"`
}

// MemoriesResult 记忆结果
type MemoriesResult struct {
	Memories []*MemoryItem `json:"memories"`
	Total    int           `json:"total"`
}

// MemoryItem 记忆项
type MemoryItem struct {
	MemoryID string                 `json:"memory_id"`
	Type     string                 `json:"type"`
	Content  string                 `json:"content"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// StoreKnowledgeRequest 存储知识请求
type StoreKnowledgeRequest struct {
	TenantID      string                 `json:"tenant_id" binding:"required"`
	KnowledgeType string                 `json:"knowledge_type" binding:"required,oneof=DOCUMENT FAQ PROCEDURE CONCEPT"`
	Title         string                 `json:"title" binding:"required"`
	Content       string                 `json:"content" binding:"required"`
	SourceURI     string                 `json:"source_uri"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// StoreKnowledgeResponse 存储知识响应
type StoreKnowledgeResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      *MemoryData `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// RetrieveKnowledgeRequest 检索知识请求
type RetrieveKnowledgeRequest struct {
	TenantID string  `json:"tenant_id" binding:"required"`
	Query    string  `json:"query" binding:"required"`
	TopK     int     `json:"top_k" binding:"min=1,max=20"`
}

// RetrieveKnowledgeResponse 检索知识响应
type RetrieveKnowledgeResponse struct {
	Code      int              `json:"code"`
	Message   string           `json:"message"`
	Data      *KnowledgeResult `json:"data,omitempty"`
	Timestamp int64            `json:"timestamp"`
}

// KnowledgeResult 知识结果
type KnowledgeResult struct {
	Knowledges []*KnowledgeItem `json:"knowledges"`
	Total      int              `json:"total"`
}

// KnowledgeItem 知识项
type KnowledgeItem struct {
	MemoryID    string                 `json:"memory_id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Score       float64                `json:"score"`
	QualityScore float64               `json:"quality_score"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
