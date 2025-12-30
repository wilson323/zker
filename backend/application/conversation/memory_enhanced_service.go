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

package conversation

import (
	"context"
	"fmt"
	"strings"

	convEntity "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	convService "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/service"
	knowService "github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/service"
)

// MemoryEnhancedChatRequest 增强对话请求
type MemoryEnhancedChatRequest struct {
	TenantID       string `json:"tenant_id"`
	UserID         string `json:"user_id"`
	ConversationID string `json:"conversation_id"`
	Query          string `json:"query"`
}

// MemoryEnhancedChatResponse 增强对话响应
type MemoryEnhancedChatResponse struct {
	Response       string `json:"response"`
	UsedMemories   int    `json:"used_memories"`
	UsedKnowledge  int    `json:"used_knowledge"`
}

// MemoryEnhancedConversationService 记忆增强的对话服务
type MemoryEnhancedConversationService struct {
	convService convService.ConversationMemoryService
	knowService knowService.KnowledgeMemoryService
	llmService  LLMService
}

// LLMService LLM服务接口
type LLMService interface {
	Chat(ctx context.Context, prompt string) (string, error)
}

// NewMemoryEnhancedConversationService 创建记忆增强对话服务
func NewMemoryEnhancedConversationService(
	convSvc convService.ConversationMemoryService,
	knowSvc knowService.KnowledgeMemoryService,
	llmSvc LLMService,
) *MemoryEnhancedConversationService {
	return &MemoryEnhancedConversationService{
		convService: convSvc,
		knowService: knowSvc,
		llmService:  llmSvc,
	}
}

// ChatWithMemory 带记忆的对话
func (s *MemoryEnhancedConversationService) ChatWithMemory(
	ctx context.Context,
	req *MemoryEnhancedChatRequest,
) (*MemoryEnhancedChatResponse, error) {
	// 1. 检索相关对话记忆
	convMemories, err := s.convService.RetrieveMemories(ctx, req.UserID, req.Query, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve conversation memories: %w", err)
	}

	// 2. 检索相关知识
	knowledges, err := s.knowService.RetrieveKnowledge(ctx, req.TenantID, req.Query, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve knowledge: %w", err)
	}

	// 3. 构建增强的Prompt
	enhancedPrompt := s.buildEnhancedPrompt(req.Query, convMemories, knowledges)

	// 4. 调用LLM
	response, err := s.llmService.Chat(ctx, enhancedPrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call LLM: %w", err)
	}

	// 5. 存储新的对话记忆
	err = s.convService.StoreMemory(ctx, &convEntity.ConversationMemory{
		TenantID:       req.TenantID,
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		MemoryType:     convEntity.MemoryTypeEntity,
		Content:        req.Query,
		ImportanceScore: 0.6,
	})
	if err != nil {
		// 记忆存储失败不影响对话
		fmt.Printf("Warning: failed to store memory: %v\n", err)
	}

	return &MemoryEnhancedChatResponse{
		Response:      response,
		UsedMemories:  len(convMemories),
		UsedKnowledge: len(knowledges),
	}, nil
}

// buildEnhancedPrompt 构建增强的Prompt
func (s *MemoryEnhancedConversationService) buildEnhancedPrompt(
	query string,
	convMemories []*convEntity.MemoryWithScore,
	knowledges []*knowEntity.KnowledgeWithScore,
) string {
	var builder strings.Builder

	// 1. 系统提示
	builder.WriteString("你是一个智能AI助手，请根据以下信息回答用户问题。\n\n")

	// 2. 相关对话记忆
	if len(convMemories) > 0 {
		builder.WriteString("【相关历史对话】\n")
		for i, mem := range convMemories {
			if mem.Score > 0.7 { // 只使用高相关性的记忆
				builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, mem.Memory.Content))
			}
		}
		builder.WriteString("\n")
	}

	// 3. 相关知识
	if len(knowledges) > 0 {
		builder.WriteString("【相关知识】\n")
		for i, know := range knowledges {
			if know.Score > 0.7 {
				builder.WriteString(fmt.Sprintf("%d. %s: %s\n", i+1, know.Knowledge.Title, know.Knowledge.Content))
			}
		}
		builder.WriteString("\n")
	}

	// 4. 用户问题
	builder.WriteString(fmt.Sprintf("【用户问题】\n%s\n", query))

	return builder.String()
}

// StoreUserPreference 存储用户偏好
func (s *MemoryEnhancedConversationService) StoreUserPreference(
	ctx context.Context,
	tenantID, userID, prefType, prefValue string,
) error {
	return s.convService.StoreMemory(ctx, &convEntity.ConversationMemory{
		TenantID:       tenantID,
		UserID:         userID,
		ConversationID: "preference",
		MemoryType:     convEntity.MemoryTypePreference,
		Content:        fmt.Sprintf("%s: %s", prefType, prefValue),
		ImportanceScore: 0.8,
	})
}

// GetUserPreferences 获取用户偏好
func (s *MemoryEnhancedConversationService) GetUserPreferences(
	ctx context.Context,
	userID string,
) (*convEntity.ExtractPreferenceResult, error) {
	return s.convService.ExtractPreferences(ctx, userID)
}

// StoreDocumentKnowledge 存储文档知识
func (s *MemoryEnhancedConversationService) StoreDocumentKnowledge(
	ctx context.Context,
	tenantID, docID, title, content string,
) error {
	_, err := s.knowService.StoreDocument(ctx, docID, title, content, 1000)
	return err
}

// GetConversationSummary 获取对话摘要
func (s *MemoryEnhancedConversationService) GetConversationSummary(
	ctx context.Context,
	conversationID string,
) (*convEntity.SummaryResult, error) {
	return s.convService.GenerateSummary(ctx, conversationID, 200)
}
