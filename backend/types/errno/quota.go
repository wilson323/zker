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

package errno

import (
	"net/http"

	"github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// Quota: 300 000 000 ~ 300 999 999
const (
	// 配额基础错误 (300 000 000 ~ 300 009 999)
	ErrQuotaNotFoundCode         = 300000001 // 配额不存在
	ErrQuotaInvalidCode          = 300000002 // 配额无效
	ErrQuotaExceededCode         = 300000003 // 配额已超限
	ErrQuotaCheckFailedCode      = 300000004 // 配额检查失败
	ErrQuotaResetFailedCode      = 300000005 // 配额重置失败
	ErrQuotaUpdateFailedCode     = 300000006 // 配额更新失败
	ErrQuotaDeleteFailedCode     = 300000007 // 配额删除失败
	ErrQuotaCreateFailedCode     = 300000008 // 配额创建失败
	ErrQuotaInvalidResourceCode  = 300000009 // 资源类型无效
	ErrQuotaInvalidLimitCode     = 300000010 // 限制值无效
	ErrQuotaConsumeFailedCode    = 300000011 // 配额消费失败
	ErrQuotaRollbackFailedCode   = 300000012 // 配额回滚失败

	// Bot 配额 (300 010 000 ~ 300 019 999)
	ErrQuotaBotExceededCode           = 300010001 // Bot 数量超限
	ErrQuotaBotCheckFailedCode        = 300010002 // Bot 配额检查失败
	ErrQuotaBotLimitReachedCode       = 300010003 // Bot 数量已达上限
	ErrQuotaBotCreationDeniedCode     = 300010004 // Bot 创建被拒绝
	ErrQuotaBotInvalidQuotaTypeCode   = 300010005 // 无效的 Bot 配额类型

	// 知识库配额 (300 020 000 ~ 300 029 999)
	ErrQuotaKnowledgeExceededCode       = 300020001 // 知识库数量超限
	ErrQuotaKnowledgeSizeExceededCode   = 300020002 // 知识库存储超限
	ErrQuotaKnowledgeDocumentExceededCode = 300020003 // 文档数量超限
	ErrQuotaKnowledgeCheckFailedCode    = 300020004 // 知识库配额检查失败

	// 工作流配额 (300 030 000 ~ 300 039 999)
	ErrQuotaWorkflowExceededCode      = 300030001 // 工作流数量超限
	ErrQuotaWorkflowNodeExceededCode  = 300030002 // 工作流节点数超限
	ErrQuotaWorkflowExecutionExceededCode = 300030003 // 工作流执行次数超限
	ErrQuotaWorkflowCheckFailedCode   = 300030004 // 工作流配额检查失败

	// API 调用配额 (300 040 000 ~ 300 049 999)
	ErrQuotaAPICallExceededCode    = 300040001 // API 调用次数超限
	ErrQuotaAPICallRateExceededCode = 300040002 // API 调用速率超限
	ErrQuotaAPICallDailyExceededCode = 300040003 // API 日调用次数超限
	ErrQuotaAPICallMonthlyExceededCode = 300040004 // API 月调用次数超限

	// 存储配额 (300 050 000 ~ 300 059 999)
	ErrQuotaStorageExceededCode    = 300050001 // 存储空间超限
	ErrQuotaStorageCheckFailedCode = 300050002 // 存储配额检查失败
	ErrQuotaStorageInvalidPathCode = 300050003 // 无效的存储路径

	// 并发配额 (300 060 000 ~ 300 069 999)
	ErrQuotaConcurrentExceededCode    = 300060001 // 并发数超限
	ErrQuotaConcurrentCheckFailedCode = 300060002 // 并发配额检查失败

	// ===== 便捷错误码常量（向后兼容） =====
	ErrQuotaCalculateFailedCode        = 300050004 // 计算配额失败
	ErrQuotaBotInvalidCode             = 300010006 // Bot配额无效
	ErrQuotaKnowledgeDocExceededCode   = 300020003 // 文档数量超限（同 ErrQuotaKnowledgeDocumentExceededCode）
	ErrQuotaKnowledgeInvalidCode       = 300020005 // 知识库配额无效
	ErrQuotaWorkflowExecExceededCode   = 300030003 // 工作流执行次数超限（同 ErrQuotaWorkflowExecutionExceededCode）
	ErrQuotaWorkflowInvalidCode        = 300030005 // 工作流配额无效
	ErrQuotaAPICallInvalidCode         = 300040005 // API调用配额无效
	ErrQuotaAPICallCheckFailedCode     = 300040006 // API调用配额检查失败
	ErrQuotaStorageFileCountExceededCode = 300050004 // 存储文件数量超限
	ErrQuotaStorageInvalidCode         = 300050005 // 存储配额无效
)

func init() {
	// 配额基础错误注册
	code.Register(
		ErrQuotaNotFoundCode,
		"Quota not found: tenant_id={tenant_id}, resource_type={resource_type}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrQuotaInvalidCode,
		"Invalid quota: tenant_id={tenant_id}, resource_type={resource_type}, reason: {reason}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrQuotaExceededCode,
		"Quota exceeded: tenant_id={tenant_id}, resource_type={resource_type}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaCheckFailedCode,
		"Failed to check quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaResetFailedCode,
		"Failed to reset quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaUpdateFailedCode,
		"Failed to update quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaDeleteFailedCode,
		"Failed to delete quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaCreateFailedCode,
		"Failed to create quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaInvalidResourceCode,
		"Invalid resource type: {resource_type}, must be one of: bots, knowledge, workflows, api_calls, storage, concurrent",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrQuotaInvalidLimitCode,
		"Invalid quota limit: resource_type={resource_type}, limit={limit}, must be greater than 0",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrQuotaConsumeFailedCode,
		"Failed to consume quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaRollbackFailedCode,
		"Failed to rollback quota: tenant_id={tenant_id}, resource_type={resource_type}, error: {error}",
		code.WithAffectStability(true),
	)

	// Bot 配额错误注册
	code.Register(
		ErrQuotaBotExceededCode,
		"Bot quota exceeded: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaBotCheckFailedCode,
		"Failed to check bot quota: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaBotLimitReachedCode,
		"Bot limit reached: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaBotCreationDeniedCode,
		"Bot creation denied: tenant_id={tenant_id}, reason: {reason}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaBotInvalidQuotaTypeCode,
		"Invalid bot quota type: {quota_type}",
		code.WithAffectStability(false),
	)

	// 知识库配额错误注册
	code.Register(
		ErrQuotaKnowledgeExceededCode,
		"Knowledge base quota exceeded: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaKnowledgeSizeExceededCode,
		"Knowledge base storage exceeded: tenant_id={tenant_id}, current_size={current_size}MB, limit={limit}MB",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaKnowledgeDocumentExceededCode,
		"Document count exceeded: tenant_id={tenant_id}, knowledge_id={knowledge_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaKnowledgeCheckFailedCode,
		"Failed to check knowledge quota: tenant_id={tenant_id}, knowledge_id={knowledge_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 工作流配额错误注册
	code.Register(
		ErrQuotaWorkflowExceededCode,
		"Workflow quota exceeded: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaWorkflowNodeExceededCode,
		"Workflow node quota exceeded: workflow_id={workflow_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaWorkflowExecutionExceededCode,
		"Workflow execution quota exceeded: tenant_id={tenant_id}, workflow_id={workflow_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaWorkflowCheckFailedCode,
		"Failed to check workflow quota: tenant_id={tenant_id}, workflow_id={workflow_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// API 调用配额错误注册
	code.Register(
		ErrQuotaAPICallExceededCode,
		"API call quota exceeded: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaAPICallRateExceededCode,
		"API call rate exceeded: tenant_id={tenant_id}, current_rate={current_rate}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaAPICallDailyExceededCode,
		"Daily API call quota exceeded: tenant_id={tenant_id}, date={date}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaAPICallMonthlyExceededCode,
		"Monthly API call quota exceeded: tenant_id={tenant_id}, month={month}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	// 存储配额错误注册
	code.Register(
		ErrQuotaStorageExceededCode,
		"Storage quota exceeded: tenant_id={tenant_id}, current_size={current_size}MB, limit={limit}MB",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaStorageCheckFailedCode,
		"Failed to check storage quota: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaStorageInvalidPathCode,
		"Invalid storage path: tenant_id={tenant_id}, path={path}",
		code.WithAffectStability(false),
	)

	// 并发配额错误注册
	code.Register(
		ErrQuotaConcurrentExceededCode,
		"Concurrent request quota exceeded: tenant_id={tenant_id}, current={current}, limit={limit}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrQuotaConcurrentCheckFailedCode,
		"Failed to check concurrent quota: tenant_id={tenant_id}, error: {error}",
		code.WithAffectStability(true),
	)
}

// 配额限制错误 (402xxx)
const (
	ErrQuotaBotLimitExceededCode    = 300040201 // Bot 数量超限
	ErrQuotaMessageLimitExceededCode = 300040202 // 消息数量超限
	ErrQuotaKnowledgeLimitExceededCode = 300040203 // 知识库存储超限
	ErrQuotaWorkflowLimitExceededCode = 300040204 // 工作流数量超限
)

// Quota 便捷错误码变量（用于直接使用）
var (
	// 配额基础错误
	ErrQuotaNotFound = &BaseErrorCode{
		code:       "QUOTA_NOT_FOUND",
		message:    "Quota not found",
		messageZH:  "配额不存在",
		messageEN:  "Quota not found",
		httpStatus: http.StatusNotFound,
	}
	ErrQuotaInvalid = &BaseErrorCode{
		code:       "QUOTA_INVALID",
		message:    "Invalid quota",
		messageZH:  "配额无效",
		messageEN:  "Invalid quota",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaExceeded = &BaseErrorCode{
		code:       "QUOTA_EXCEEDED",
		message:    "Quota exceeded",
		messageZH:  "配额已超限",
		messageEN:  "Quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaCheckFailed = &BaseErrorCode{
		code:       "QUOTA_CHECK_FAILED",
		message:    "Quota check failed",
		messageZH:  "配额检查失败",
		messageEN:  "Quota check failed",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaResetFailed = &BaseErrorCode{
		code:       "QUOTA_RESET_FAILED",
		message:    "Failed to reset quota",
		messageZH:  "配额重置失败",
		messageEN:  "Failed to reset quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaUpdateFailed = &BaseErrorCode{
		code:       "QUOTA_UPDATE_FAILED",
		message:    "Failed to update quota",
		messageZH:  "配额更新失败",
		messageEN:  "Failed to update quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaDeleteFailed = &BaseErrorCode{
		code:       "QUOTA_DELETE_FAILED",
		message:    "Failed to delete quota",
		messageZH:  "配额删除失败",
		messageEN:  "Failed to delete quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaCreateFailed = &BaseErrorCode{
		code:       "QUOTA_CREATE_FAILED",
		message:    "Failed to create quota",
		messageZH:  "配额创建失败",
		messageEN:  "Failed to create quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaInvalidResource = &BaseErrorCode{
		code:       "QUOTA_INVALID_RESOURCE",
		message:    "Invalid resource type",
		messageZH:  "资源类型无效",
		messageEN:  "Invalid resource type",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaInvalidLimit = &BaseErrorCode{
		code:       "QUOTA_INVALID_LIMIT",
		message:    "Invalid quota limit",
		messageZH:  "限制值无效",
		messageEN:  "Invalid quota limit",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaConsumeFailed = &BaseErrorCode{
		code:       "QUOTA_CONSUME_FAILED",
		message:    "Failed to consume quota",
		messageZH:  "配额消费失败",
		messageEN:  "Failed to consume quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaRollbackFailed = &BaseErrorCode{
		code:       "QUOTA_ROLLBACK_FAILED",
		message:    "Failed to rollback quota",
		messageZH:  "配额回滚失败",
		messageEN:  "Failed to rollback quota",
		httpStatus: http.StatusInternalServerError,
	}

	// Bot 配额错误
	ErrQuotaBotExceeded = &BaseErrorCode{
		code:       "QUOTA_BOT_EXCEEDED",
		message:    "Bot quota exceeded",
		messageZH:  "Bot 数量超限",
		messageEN:  "Bot quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaBotCheckFailed = &BaseErrorCode{
		code:       "QUOTA_BOT_CHECK_FAILED",
		message:    "Failed to check bot quota",
		messageZH:  "Bot 配额检查失败",
		messageEN:  "Failed to check bot quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaBotLimitReached = &BaseErrorCode{
		code:       "QUOTA_BOT_LIMIT_REACHED",
		message:    "Bot limit reached",
		messageZH:  "Bot 数量已达上限",
		messageEN:  "Bot limit reached",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaBotCreationDenied = &BaseErrorCode{
		code:       "QUOTA_BOT_CREATION_DENIED",
		message:    "Bot creation denied",
		messageZH:  "Bot 创建被拒绝",
		messageEN:  "Bot creation denied",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaBotInvalidQuotaType = &BaseErrorCode{
		code:       "QUOTA_BOT_INVALID_QUOTA_TYPE",
		message:    "Invalid bot quota type",
		messageZH:  "无效的 Bot 配额类型",
		messageEN:  "Invalid bot quota type",
		httpStatus: http.StatusBadRequest,
	}

	// 知识库配额错误
	ErrQuotaKnowledgeExceeded = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_EXCEEDED",
		message:    "Knowledge base quota exceeded",
		messageZH:  "知识库数量超限",
		messageEN:  "Knowledge base quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaKnowledgeSizeExceeded = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_SIZE_EXCEEDED",
		message:    "Knowledge base storage exceeded",
		messageZH:  "知识库存储超限",
		messageEN:  "Knowledge base storage exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaKnowledgeDocumentExceeded = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_DOCUMENT_EXCEEDED",
		message:    "Document count exceeded",
		messageZH:  "文档数量超限",
		messageEN:  "Document count exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaKnowledgeCheckFailed = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_CHECK_FAILED",
		message:    "Failed to check knowledge quota",
		messageZH:  "知识库配额检查失败",
		messageEN:  "Failed to check knowledge quota",
		httpStatus: http.StatusInternalServerError,
	}

	// 工作流配额错误
	ErrQuotaWorkflowExceeded = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_EXCEEDED",
		message:    "Workflow quota exceeded",
		messageZH:  "工作流数量超限",
		messageEN:  "Workflow quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaWorkflowNodeExceeded = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_NODE_EXCEEDED",
		message:    "Workflow node quota exceeded",
		messageZH:  "工作流节点数超限",
		messageEN:  "Workflow node quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaWorkflowExecutionExceeded = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_EXECUTION_EXCEEDED",
		message:    "Workflow execution quota exceeded",
		messageZH:  "工作流执行次数超限",
		messageEN:  "Workflow execution quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaWorkflowCheckFailed = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_CHECK_FAILED",
		message:    "Failed to check workflow quota",
		messageZH:  "工作流配额检查失败",
		messageEN:  "Failed to check workflow quota",
		httpStatus: http.StatusInternalServerError,
	}

	// API 调用配额错误
	ErrQuotaAPICallExceeded = &BaseErrorCode{
		code:       "QUOTA_API_CALL_EXCEEDED",
		message:    "API call quota exceeded",
		messageZH:  "API 调用次数超限",
		messageEN:  "API call quota exceeded",
		httpStatus: http.StatusTooManyRequests,
	}
	ErrQuotaAPICallRateExceeded = &BaseErrorCode{
		code:       "QUOTA_API_CALL_RATE_EXCEEDED",
		message:    "API call rate exceeded",
		messageZH:  "API 调用速率超限",
		messageEN:  "API call rate exceeded",
		httpStatus: http.StatusTooManyRequests,
	}
	ErrQuotaAPICallDailyExceeded = &BaseErrorCode{
		code:       "QUOTA_API_CALL_DAILY_EXCEEDED",
		message:    "Daily API call quota exceeded",
		messageZH:  "API 日调用次数超限",
		messageEN:  "Daily API call quota exceeded",
		httpStatus: http.StatusTooManyRequests,
	}
	ErrQuotaAPICallMonthlyExceeded = &BaseErrorCode{
		code:       "QUOTA_API_CALL_MONTHLY_EXCEEDED",
		message:    "Monthly API call quota exceeded",
		messageZH:  "API 月调用次数超限",
		messageEN:  "Monthly API call quota exceeded",
		httpStatus: http.StatusTooManyRequests,
	}

	// 存储配额错误
	ErrQuotaStorageExceeded = &BaseErrorCode{
		code:       "QUOTA_STORAGE_EXCEEDED",
		message:    "Storage quota exceeded",
		messageZH:  "存储空间超限",
		messageEN:  "Storage quota exceeded",
		httpStatus: http.StatusInsufficientStorage,
	}
	ErrQuotaStorageCheckFailed = &BaseErrorCode{
		code:       "QUOTA_STORAGE_CHECK_FAILED",
		message:    "Failed to check storage quota",
		messageZH:  "存储配额检查失败",
		messageEN:  "Failed to check storage quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaStorageInvalidPath = &BaseErrorCode{
		code:       "QUOTA_STORAGE_INVALID_PATH",
		message:    "Invalid storage path",
		messageZH:  "无效的存储路径",
		messageEN:  "Invalid storage path",
		httpStatus: http.StatusBadRequest,
	}

	// 并发配额错误
	ErrQuotaConcurrentExceeded = &BaseErrorCode{
		code:       "QUOTA_CONCURRENT_EXCEEDED",
		message:    "Concurrent request quota exceeded",
		messageZH:  "并发数超限",
		messageEN:  "Concurrent request quota exceeded",
		httpStatus: http.StatusTooManyRequests,
	}
	ErrQuotaConcurrentCheckFailed = &BaseErrorCode{
		code:       "QUOTA_CONCURRENT_CHECK_FAILED",
		message:    "Failed to check concurrent quota",
		messageZH:  "并发配额检查失败",
		messageEN:  "Failed to check concurrent quota",
		httpStatus: http.StatusInternalServerError,
	}

	// 便捷错误码变量
	ErrQuotaCalculateFailed = &BaseErrorCode{
		code:       "QUOTA_CALCULATE_FAILED",
		message:    "Failed to calculate quota",
		messageZH:  "计算配额失败",
		messageEN:  "Failed to calculate quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaBotInvalid = &BaseErrorCode{
		code:       "QUOTA_BOT_INVALID",
		message:    "Invalid bot quota",
		messageZH:  "Bot配额无效",
		messageEN:  "Invalid bot quota",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaKnowledgeDocExceeded = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_DOC_EXCEEDED",
		message:    "Document count exceeded",
		messageZH:  "文档数量超限",
		messageEN:  "Document count exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaKnowledgeInvalid = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_INVALID",
		message:    "Invalid knowledge quota",
		messageZH:  "知识库配额无效",
		messageEN:  "Invalid knowledge quota",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaWorkflowExecExceeded = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_EXEC_EXCEEDED",
		message:    "Workflow execution quota exceeded",
		messageZH:  "工作流执行次数超限",
		messageEN:  "Workflow execution quota exceeded",
		httpStatus: http.StatusForbidden,
	}
	ErrQuotaWorkflowInvalid = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_INVALID",
		message:    "Invalid workflow quota",
		messageZH:  "工作流配额无效",
		messageEN:  "Invalid workflow quota",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaAPICallInvalid = &BaseErrorCode{
		code:       "QUOTA_API_CALL_INVALID",
		message:    "Invalid API call quota",
		messageZH:  "API调用配额无效",
		messageEN:  "Invalid API call quota",
		httpStatus: http.StatusBadRequest,
	}
	ErrQuotaAPICallCheckFailed = &BaseErrorCode{
		code:       "QUOTA_API_CALL_CHECK_FAILED",
		message:    "Failed to check API call quota",
		messageZH:  "API调用配额检查失败",
		messageEN:  "Failed to check API call quota",
		httpStatus: http.StatusInternalServerError,
	}
	ErrQuotaStorageFileCountExceeded = &BaseErrorCode{
		code:       "QUOTA_STORAGE_FILE_COUNT_EXCEEDED",
		message:    "Storage file count exceeded",
		messageZH:  "存储文件数量超限",
		messageEN:  "Storage file count exceeded",
		httpStatus: http.StatusInsufficientStorage,
	}
	ErrQuotaStorageInvalid = &BaseErrorCode{
		code:       "QUOTA_STORAGE_INVALID",
		message:    "Invalid storage quota",
		messageZH:  "存储配额无效",
		messageEN:  "Invalid storage quota",
		httpStatus: http.StatusBadRequest,
	}

	// 配额限制错误 (402xxx)
	ErrQuotaBotLimitExceeded = &BaseErrorCode{
		code:       "QUOTA_BOT_LIMIT_EXCEEDED",
		message:    "Bot quota exceeded",
		messageZH:  "Bot数量超限",
		messageEN:  "Bot quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
	ErrQuotaMessageLimitExceeded = &BaseErrorCode{
		code:       "QUOTA_MESSAGE_LIMIT_EXCEEDED",
		message:    "Message quota exceeded",
		messageZH:  "消息数量超限",
		messageEN:  "Message quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
	ErrQuotaKnowledgeLimitExceeded = &BaseErrorCode{
		code:       "QUOTA_KNOWLEDGE_LIMIT_EXCEEDED",
		message:    "Knowledge storage quota exceeded",
		messageZH:  "知识库存储超限",
		messageEN:  "Knowledge storage quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
	ErrQuotaWorkflowLimitExceeded = &BaseErrorCode{
		code:       "QUOTA_WORKFLOW_LIMIT_EXCEEDED",
		message:    "Workflow quota exceeded",
		messageZH:  "工作流数量超限",
		messageEN:  "Workflow quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
)

// QUOTA402001-QUOTA402004 错误码变量（向后兼容）
var (
	QUOTA402001 = &BaseErrorCode{
		code:       "QUOTA402001",
		message:    "Bot quota exceeded",
		messageZH:  "Bot数量超限",
		messageEN:  "Bot quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
	QUOTA402002 = &BaseErrorCode{
		code:       "QUOTA402002",
		message:    "Message quota exceeded",
		messageZH:  "消息数量超限",
		messageEN:  "Message quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
	QUOTA402003 = &BaseErrorCode{
		code:       "QUOTA402003",
		message:    "Knowledge storage quota exceeded",
		messageZH:  "知识库存储超限",
		messageEN:  "Knowledge storage quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
	QUOTA402004 = &BaseErrorCode{
		code:       "QUOTA402004",
		message:    "Workflow quota exceeded",
		messageZH:  "工作流数量超限",
		messageEN:  "Workflow quota exceeded",
		httpStatus: http.StatusPaymentRequired,
	}
)
