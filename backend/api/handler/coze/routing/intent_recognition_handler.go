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

package routing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	routingapp "github.com/coze-dev/coze-studio/backend/application/routing"
)

// IntentRecognitionHandler 意图识别Handler
// 职责：处理意图识别相关的HTTP请求
type IntentRecognitionHandler struct {
	intentService *routingapp.IntentRecognitionService
	logger        *zap.Logger
}

// NewIntentRecognitionHandler 创建意图识别Handler实例
func NewIntentRecognitionHandler(
	intentService *routingapp.IntentRecognitionService,
	logger *zap.Logger,
) *IntentRecognitionHandler {
	return &IntentRecognitionHandler{
		intentService: intentService,
		logger:        logger,
	}
}

// RecognizeIntentRequest 识别意图请求
type RecognizeIntentRequest struct {
	TenantID string `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	Text     string `json:"text" binding:"required" vd:"len($) > 0"`
	Strategy string `json:"strategy" binding:"omitempty,oneof=llm vector hybrid"` // llm, vector, hybrid
}

// RecognizeIntentResponse 识别意图响应
type RecognizeIntentResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    *IntentRecognitionResult  `json:"data,omitempty"`
}

// IntentRecognitionResult 意图识别结果
type IntentRecognitionResult struct {
	IntentID     string  `json:"intent_id"`
	IntentName   string  `json:"intent_name"`
	Confidence   float64 `json:"confidence"`
	AgentID      string  `json:"agent_id"`
	WorkflowID   string  `json:"workflow_id,omitempty"`
	MatchMethod  string  `json:"match_method"`  // llm, vector, hybrid
	Entities     map[string][]string `json:"entities,omitempty"`
}

// RecognizeIntent 识别意图
// @router /api/routing/intents/recognize [POST]
func (h *IntentRecognitionHandler) RecognizeIntent(ctx context.Context, c *app.RequestContext) {
	var req RecognizeIntentRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层识别意图
	result, err := h.intentService.RecognizeIntent(ctx, req.TenantID, req.Text, req.Strategy)
	if err != nil {
		h.logger.Error("failed to recognize intent",
			zap.String("tenant_id", req.TenantID),
			zap.String("text", req.Text),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	response := &IntentRecognitionResult{
		IntentID:    result.IntentID,
		IntentName:  result.IntentName,
		Confidence:  result.Confidence,
		AgentID:     result.AgentID,
		WorkflowID:  result.WorkflowID,
		MatchMethod: result.MatchMethod,
		Entities:    result.Entities,
	}

	c.JSON(http.StatusOK, RecognizeIntentResponse{
		Code:    0,
		Message: "success",
		Data:    response,
	})
}

// BatchRecognizeRequest 批量识别请求
type BatchRecognizeRequest struct {
	TenantID string   `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	Texts    []string `json:"texts" binding:"required,min=1,max=50"`
	Strategy string   `json:"strategy" binding:"omitempty,oneof=llm vector hybrid"`
}

// BatchRecognizeResponse 批量识别响应
type BatchRecognizeResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    []IntentRecognitionResult `json:"data,omitempty"`
}

// BatchRecognize 批量识别意图
// @router /api/routing/intents/batch [POST]
func (h *IntentRecognitionHandler) BatchRecognize(ctx context.Context, c *app.RequestContext) {
	var req BatchRecognizeRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层批量识别
	results, err := h.intentService.BatchRecognize(ctx, req.TenantID, req.Texts, req.Strategy)
	if err != nil {
		h.logger.Error("failed to batch recognize intents",
			zap.String("tenant_id", req.TenantID),
			zap.Int("count", len(req.Texts)),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	response := make([]IntentRecognitionResult, 0, len(results))
	for _, result := range results {
		response = append(response, IntentRecognitionResult{
			IntentID:    result.IntentID,
			IntentName:  result.IntentName,
			Confidence:  result.Confidence,
			AgentID:     result.AgentID,
			WorkflowID:  result.WorkflowID,
			MatchMethod: result.MatchMethod,
			Entities:    result.Entities,
		})
	}

	c.JSON(http.StatusOK, BatchRecognizeResponse{
		Code:    0,
		Message: "success",
		Data:    response,
	})
}

// TrainIntentModelRequest 训练意图模型请求
type TrainIntentModelRequest struct {
	IntentID  string   `json:"intent_id" binding:"required" vd:"len($) > 0"`
	TenantID  string   `json:"tenant_id" binding:"required" vd:"len($) > 0"`
	Samples   []string `json:"samples" binding:"required,min=5,max=1000"` // 至少5个训练样本
	Strategy  string   `json:"strategy" binding:"omitempty,oneof=llm vector hybrid"`
}

// TrainIntentModelResponse 训练意图模型响应
type TrainIntentModelResponse struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Data       *TrainResult `json:"data,omitempty"`
}

// TrainResult 训练结果
type TrainResult struct {
	IntentID      string  `json:"intent_id"`
	SampleCount   int     `json:"sample_count"`
	Accuracy      float64 `json:"accuracy"`
	TrainingTime  int64   `json:"training_time"` // 毫秒
}

// TrainIntentModel 训练意图模型
// @router /api/routing/intents/:id/train [POST]
func (h *IntentRecognitionHandler) TrainIntentModel(ctx context.Context, c *app.RequestContext) {
	intentID := c.Param("id")
	if intentID == "" {
		httputil.BadRequest(c, "intent_id is required")
		return
	}

	var req TrainIntentModelRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	req.IntentID = intentID

	// 调用Service层训练模型
	result, err := h.intentService.TrainIntentModel(ctx, req.TenantID, req.IntentID, req.Samples, req.Strategy)
	if err != nil {
		h.logger.Error("failed to train intent model",
			zap.String("tenant_id", req.TenantID),
			zap.String("intent_id", req.IntentID),
			zap.Int("sample_count", len(req.Samples)),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, TrainIntentModelResponse{
		Code:    0,
		Message: "success",
		Data: &TrainResult{
			IntentID:     result.IntentID,
			SampleCount:  result.SampleCount,
			Accuracy:     result.Accuracy,
			TrainingTime: result.TrainingTime,
		},
	})
}

// GetIntentsRequest 获取意图列表请求
type GetIntentsRequest struct {
	TenantID  string `form:"tenant_id" binding:"required" vd:"len($) > 0"`
	IsActive  *bool  `form:"is_active"`
	AgentID   string `form:"agent_id"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// GetIntentsResponse 获取意图列表响应
type GetIntentsResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *IntentListData `json:"data,omitempty"`
}

// IntentListData 意图列表数据
type IntentListData struct {
	Intents []IntentInfo `json:"intents"`
	Total   int          `json:"total"`
	Page    int          `json:"page"`
	PageSize int         `json:"page_size"`
}

// IntentInfo 意图信息
type IntentInfo struct {
	IntentID    string   `json:"intent_id"`
	IntentName  string   `json:"intent_name"`
	Description string   `json:"description"`
	AgentID     string   `json:"agent_id"`
	WorkflowID  string   `json:"workflow_id,omitempty"`
	Confidence  float64  `json:"confidence"`
	IsActive    bool     `json:"is_active"`
	Examples    []string `json:"examples"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

// GetIntents 获取意图列表
// @router /api/routing/intents [GET]
func (h *IntentRecognitionHandler) GetIntents(ctx context.Context, c *app.RequestContext) {
	var req GetIntentsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分页
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 调用Service层获取意图列表
	intents, total, err := h.intentService.GetIntents(ctx, req.TenantID, req.IsActive, req.AgentID, req.Page, req.PageSize)
	if err != nil {
		h.logger.Error("failed to get intents",
			zap.String("tenant_id", req.TenantID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 转换响应
	intentList := make([]IntentInfo, 0, len(intents))
	for _, intent := range intents {
		var examples []string
		if intent.Examples != "" {
			json.Unmarshal([]byte(intent.Examples), &examples)
		}

		intentList = append(intentList, IntentInfo{
			IntentID:    intent.IntentID,
			IntentName:  intent.IntentName,
			Description: intent.Description,
			AgentID:     intent.AgentID,
			WorkflowID:  intent.WorkflowID,
			Confidence:  intent.Confidence,
			IsActive:    intent.IsActive,
			Examples:    examples,
			CreatedAt:   intent.CreatedAt,
			UpdatedAt:   intent.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, GetIntentsResponse{
		Code:    0,
		Message: "success",
		Data: &IntentListData{
			Intents:  intentList,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	})
}

// UpdateIntentRequest 更新意图请求
type UpdateIntentRequest struct {
	IntentName  string   `json:"intent_name" binding:"omitempty,min=1,max=100"`
	Description string   `json:"description" binding:"omitempty,max=500"`
	Examples    []string `json:"examples" binding:"omitempty,min=1,max=100"`
	AgentID     string   `json:"agent_id" binding:"omitempty"`
	WorkflowID  string   `json:"workflow_id" binding:"omitempty"`
	Confidence  float64  `json:"confidence" binding:"omitempty,min=0,max=1"`
	IsActive    *bool    `json:"is_active" binding:"omitempty"`
}

// UpdateIntentResponse 更新意图响应
type UpdateIntentResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *IntentInfo `json:"data,omitempty"`
}

// UpdateIntent 更新意图
// @router /api/routing/intents/:id [PUT]
func (h *IntentRecognitionHandler) UpdateIntent(ctx context.Context, c *app.RequestContext) {
	intentID := c.Param("id")
	if intentID == "" {
		httputil.BadRequest(c, "intent_id is required")
		return
	}

	var req UpdateIntentRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 调用Service层更新意图
	intent, err := h.intentService.UpdateIntent(ctx, intentID, &req)
	if err != nil {
		h.logger.Error("failed to update intent",
			zap.String("intent_id", intentID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	// 解析示例句子
	var examples []string
	if intent.Examples != "" {
		json.Unmarshal([]byte(intent.Examples), &examples)
	}

	c.JSON(http.StatusOK, UpdateIntentResponse{
		Code:    0,
		Message: "success",
		Data: &IntentInfo{
			IntentID:    intent.IntentID,
			IntentName:  intent.IntentName,
			Description: intent.Description,
			AgentID:     intent.AgentID,
			WorkflowID:  intent.WorkflowID,
			Confidence:  intent.Confidence,
			IsActive:    intent.IsActive,
			Examples:    examples,
			CreatedAt:   intent.CreatedAt,
			UpdatedAt:   intent.UpdatedAt,
		},
	})
}

// DeleteIntentResponse 删除意图响应
type DeleteIntentResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DeleteIntent 删除意图
// @router /api/routing/intents/:id [DELETE]
func (h *IntentRecognitionHandler) DeleteIntent(ctx context.Context, c *app.RequestContext) {
	intentID := c.Param("id")
	if intentID == "" {
		httputil.BadRequest(c, "intent_id is required")
		return
	}

	// 调用Service层删除意图
	if err := h.intentService.DeleteIntent(ctx, intentID); err != nil {
		h.logger.Error("failed to delete intent",
			zap.String("intent_id", intentID),
			zap.Error(err))
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, DeleteIntentResponse{
		Code:    0,
		Message: "success",
	})
}
