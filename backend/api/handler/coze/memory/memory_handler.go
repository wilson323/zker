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

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	convEntity "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/entity"
	convService "github.com/coze-dev/coze-studio/backend/domain/memory/conversation/service"
	knowEntity "github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
	knowService "github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/service"
	"github.com/coze-dev/coze-studio/backend/api/model/memory"
	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// MemoryHandler 记忆处理器
type MemoryHandler struct {
	convService convService.ConversationMemoryService
	knowService knowService.KnowledgeMemoryService
}

// NewMemoryHandler 创建记忆处理器
func NewMemoryHandler(
	convSvc convService.ConversationMemoryService,
	knowSvc knowService.KnowledgeMemoryService,
) *MemoryHandler {
	return &MemoryHandler{
		convService: convSvc,
		knowService: knowSvc,
	}
}

// StoreConversationMemory 存储对话记忆
// @router /api/v1/memory/conversation [POST]
func (h *MemoryHandler) StoreConversationMemory(ctx context.Context, c *app.RequestContext) {
	var req memory.StoreConversationMemoryRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryInvalidParamCode), err.Error(), "参数验证失败", nil)
		return
	}

	// 构建记忆实体
	mem := &convEntity.ConversationMemory{
		TenantID:       req.TenantID,
		UserID:         req.UserID,
		ConversationID: req.ConversationID,
		MemoryType:     convEntity.MemoryType(req.MemoryType),
		Content:        req.Content,
		Metadata:       req.Metadata,
		ImportanceScore: 0.5,
	}

	// 存储记忆
	if err := h.convService.StoreMemory(ctx, mem); err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryIDGenFailCode), err.Error(), "存储记忆失败", nil)
		return
	}

	// 返回成功响应
	httputil.BuildSuccessResp(c, &memory.MemoryData{
		MemoryID:  mem.MemoryID,
		VectorID:  mem.MemoryID,
		Type:      req.MemoryType,
		CreatedAt: mem.CreatedAt.Format(time.RFC3339),
	})
}

// RetrieveConversationMemories 检索对话记忆
// @router /api/v1/memory/conversation/search [POST]
func (h *MemoryHandler) RetrieveConversationMemories(ctx context.Context, c *app.RequestContext) {
	var req memory.RetrieveMemoriesRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryInvalidParamCode), err.Error(), "参数验证失败", nil)
		return
	}

	// 检索记忆
	memories, err := h.convService.RetrieveMemories(ctx, req.UserID, req.Query, req.TopK)
	if err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryIDGenFailCode), err.Error(), "检索记忆失败", nil)
		return
	}

	// 转换为响应格式
	items := make([]*memory.MemoryItem, len(memories))
	for i, mem := range memories {
		items[i] = &memory.MemoryItem{
			MemoryID: mem.Memory.MemoryID,
			Type:     string(mem.Memory.MemoryType),
			Content:  mem.Memory.Content,
			Score:    mem.Score,
			Metadata: mem.Memory.Metadata,
		}
	}

	httputil.BuildSuccessResp(c, &memory.MemoriesResult{
		Memories: items,
		Total:    len(items),
	})
}

// StoreKnowledge 存储知识
// @router /api/v1/memory/knowledge [POST]
func (h *MemoryHandler) StoreKnowledge(ctx context.Context, c *app.RequestContext) {
	var req memory.StoreKnowledgeRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryInvalidParamCode), err.Error(), "参数验证失败", nil)
		return
	}

	// 构建知识实体
	know := &knowEntity.KnowledgeMemory{
		TenantID:      req.TenantID,
		KnowledgeType: knowEntity.KnowledgeType(req.KnowledgeType),
		Title:         req.Title,
		Content:       req.Content,
		SourceURI:     req.SourceURI,
		Metadata:      req.Metadata,
		QualityScore:  0.5,
		Version:       1,
	}

	// 存储知识
	if err := h.knowService.StoreKnowledge(ctx, know); err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryIDGenFailCode), err.Error(), "存储知识失败", nil)
		return
	}

	// 返回成功响应
	httputil.BuildSuccessResp(c, &memory.MemoryData{
		MemoryID:  know.MemoryID,
		VectorID:  know.MemoryID,
		Type:      req.KnowledgeType,
		CreatedAt: know.CreatedAt.Format(time.RFC3339),
	})
}

// RetrieveKnowledge 检索知识
// @router /api/v1/memory/knowledge/search [POST]
func (h *MemoryHandler) RetrieveKnowledge(ctx context.Context, c *app.RequestContext) {
	var req memory.RetrieveKnowledgeRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryInvalidParamCode), err.Error(), "参数验证失败", nil)
		return
	}

	// 检索知识
	knowledges, err := h.knowService.RetrieveKnowledge(ctx, req.TenantID, req.Query, req.TopK)
	if err != nil {
		httputil.BuildErrorResp(c, int32(errno.ErrMemoryIDGenFailCode), err.Error(), "检索知识失败", nil)
		return
	}

	// 转换为响应格式
	items := make([]*memory.KnowledgeItem, len(knowledges))
	for i, know := range knowledges {
		items[i] = &memory.KnowledgeItem{
			MemoryID:     know.Knowledge.MemoryID,
			Type:         string(know.Knowledge.KnowledgeType),
			Title:        know.Knowledge.Title,
			Content:      know.Knowledge.Content,
			Score:        know.Score,
			QualityScore: know.Knowledge.QualityScore,
			Metadata:     know.Knowledge.Metadata,
		}
	}

	httputil.BuildSuccessResp(c, &memory.KnowledgeResult{
		Knowledges: items,
		Total:      len(items),
	})
}