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

package httputil

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/pkg/conv"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

type data struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

// ErrorResponse 企业级统一错误响应格式
// 遵循ZKER统一错误码定义规范
type ErrorResponse struct {
	Code       int32                  `json:"code"`                 // 错误码
	Message    string                 `json:"message"`              // 英文错误消息
	MessageZH  string                 `json:"message_zh,omitempty"` // 中文错误消息
	MessageEN  string                 `json:"message_en,omitempty"` // 英文错误消息(冗余字段,便于前端使用)
	Details    map[string]interface{} `json:"details,omitempty"`    // 错误详细信息
	RequestID  string                 `json:"request_id,omitempty"` // 请求追踪ID
	Timestamp  string                 `json:"timestamp,omitempty"`  // 错误发生时间
	TraceID    string                 `json:"trace_id,omitempty"`   // 分布式追踪ID
	TenantID   string                 `json:"tenant_id,omitempty"`  // 租户ID(企业级特性)
}

// BadRequest 400错误响应
func BadRequest(c *app.RequestContext, errMsg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, data{Code: http.StatusBadRequest, Msg: errMsg})
}

// Unauthorized 401错误响应
func Unauthorized(c *app.RequestContext, errMsg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, data{Code: http.StatusUnauthorized, Msg: errMsg})
}

// InternalError 500错误响应(兼容旧版本)
func InternalError(ctx context.Context, c *app.RequestContext, err error) {
	var customErr errorx.StatusError

	if errors.As(err, &customErr) && customErr.Code() != 0 {
		logs.CtxWarnf(ctx, "[ErrorX] error:  %v %v \n", customErr.Code(), err)
		c.AbortWithStatusJSON(http.StatusOK, data{Code: customErr.Code(), Msg: customErr.Msg()})
		return
	}

	logs.CtxErrorf(ctx, "[InternalError]  error: %v \n", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, data{Code: 500, Msg: "internal server error"})
}

// BuildErrorResp 构建企业级统一错误响应(新版本)
// 遵循KISS原则:函数<50行,职责单一
// 遵循DRY原则:复用EnhancedError逻辑
func BuildErrorResp(c *app.RequestContext, errCode int32, message, messageZH string, details map[string]interface{}) {
	response := &ErrorResponse{
		Code:      errCode,
		Message:   message,
		MessageZH: messageZH,
		MessageEN: message,
		Details:   details,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 从context获取请求ID和追踪ID
	if reqID := c.Query("request_id"); reqID != "" {
		response.RequestID = reqID
	}
	if traceID := conv.BytesToString(c.GetHeader("X-Trace-ID")); traceID != "" {
		response.TraceID = traceID
	}
	if tenantID := conv.BytesToString(c.GetHeader("X-Tenant-ID")); tenantID != "" {
		response.TenantID = tenantID
	}

	// 获取HTTP状态码
	httpStatus := errno.GetHTTPStatusForError(errCode)
	c.AbortWithStatusJSON(httpStatus, response)
}

// BuildErrorRespFromEnhanced 从EnhancedError构建错误响应
// 遵循KISS原则:简单转换,不添加额外逻辑
func BuildErrorRespFromEnhanced(c *app.RequestContext, enhancedErr *errno.EnhancedError) {
	response := &ErrorResponse{
		Code:       enhancedErr.Code,
		Message:    enhancedErr.Message,
		MessageZH:  enhancedErr.MessageZH,
		MessageEN:  enhancedErr.Message,
		Details:    enhancedErr.Details,
		RequestID:  enhancedErr.RequestID,
		Timestamp:  enhancedErr.Timestamp,
		TraceID:    enhancedErr.TraceID,
		TenantID:   enhancedErr.TenantID,
	}

	// 获取HTTP状态码
	httpStatus := enhancedErr.GetHTTPStatus()
	c.AbortWithStatusJSON(httpStatus, response)
}

// BuildSuccessResp 构建成功响应
// 遵循KISS原则:简单直接
func BuildSuccessResp(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":        0,
		"message":     "SUCCESS",
		"message_zh":  "操作成功",
		"message_en":  "Operation successful",
		"data":        data,
		"timestamp":   time.Now().Format(time.RFC3339),
	})
}
