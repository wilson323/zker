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
	"github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// =====================================================
// 统一错误码规范 - Conversation模块
// 错误码段: 202 000 000 ~ 202 999 999
// =====================================================

const (
	// Conversation基础错误 (202 000 000 ~ 202 009 999)
	ErrConversationNotFoundCode        = 202000001 // 对话不存在
	ErrConversationAlreadyExistsCode   = 202000002 // 对话已存在
	ErrConversationInvalidParamCode    = 202000003 // 对话参数无效
	ErrConversationTitleTooLongCode    = 202000004 // 对话标题过长
	ErrConversationClosedCode          = 202000005 // 对话已关闭
	ErrConversationDeleteFailedCode    = 202000006 // 删除对话失败
	ErrConversationUpdateFailedCode    = 202000007 // 更新对话失败
	ErrConversationPermissionDeniedCode = 202000008 // 对话权限不足

	// Message相关 (202 010 000 ~ 202 019 999)
	ErrMessageNotFoundCode          = 202010001 // 消息不存在
	ErrMessageInvalidParamCode       = 202010002 // 消息参数无效
	ErrMessageContentTooLongCode     = 202010003 // 消息内容过长
	ErrMessageContentEmptyCode       = 202010004 // 消息内容为空
	ErrMessageSendFailedCode         = 202010005 // 发送消息失败
	ErrMessageDeleteFailedCode       = 202010006 // 删除消息失败
	ErrMessageUpdateFailedCode       = 202010007 // 更新消息失败
	ErrMessageBatchSendFailedCode    = 202010008 // 批量发送消息失败
	ErrMessageTypeNotSupportedCode   = 202010009 // 消息类型不支持

	// Agent运行相关 (202 020 000 ~ 202 029 999)
	ErrAgentRunNotFoundCode          = 202020001 // Agent运行不存在
	ErrAgentRunAlreadyCompletedCode  = 202020002 // Agent运行已完成
	ErrAgentRunFailedCode            = 202020003 // Agent运行失败
	ErrAgentRunTimeoutCode           = 202020004 // Agent运行超时
	ErrAgentRunCancelledCode         = 202020005 // Agent运行已取消
	ErrAgentRunStatusInvalidCode     = 202020006 // Agent运行状态无效
	ErrAgentRunCreateFailedCode      = 202020007 // 创建Agent运行失败
	ErrAgentRunRetryExceededCode     = 202020008 // Agent运行重试次数超限
	ErrAgentRunQuotaExceededCode     = 202020009 // Agent运行配额已超限

	// 流式输出相关 (202 030 000 ~ 202 039 999)
	ErrStreamConnectionLostCode      = 202030001 // 流式连接丢失
	ErrStreamTimeoutCode             = 202030002 // 流式输出超时
	ErrStreamInvalidFormatCode       = 202030003 // 流式输出格式无效
	ErrStreamCreateFailedCode        = 202030004 // 创建流式输出失败
	ErrStreamWriteFailedCode         = 202030005 // 写入流式输出失败

	// 上下文相关 (202 040 000 ~ 202 049 999)
	ErrContextNotFoundCode           = 202040001 // 上下文不存在
	ErrContextTooLargeCode           = 202040002 // 上下文过大
	ErrContextInvalidFormatCode      = 202040003 // 上下文格式无效
	ErrContextUpdateFailedCode       = 202040004 // 更新上下文失败
	ErrContextClearFailedCode        = 202040005 // 清除上下文失败

	// 附件相关 (202 050 000 ~ 202 059 999)
	ErrAttachmentNotFoundCode        = 202050001 // 附件不存在
	ErrAttachmentUploadFailedCode    = 202050002 // 上传附件失败
	ErrAttachmentInvalidTypeCode     = 202050003 // 附件类型无效
	ErrAttachmentSizeExceededCode    = 202050004 // 附件大小超限
	ErrAttachmentDeleteFailedCode    = 202050005 // 删除附件失败
	ErrAttachmentDownloadFailedCode  = 202050006 // 下载附件失败

	// 反馈相关 (202 060 000 ~ 202 069 999)
	ErrFeedbackNotFoundCode          = 202060001 // 反馈不存在
	ErrFeedbackInvalidTypeCode       = 202060002 // 反馈类型无效
	ErrFeedbackAlreadyExistsCode     = 202060003 // 反馈已存在
	ErrFeedbackCreateFailedCode      = 202060004 // 创建反馈失败
	ErrFeedbackUpdateFailedCode      = 202060005 // 更新反馈失败

	// 重试相关 (202 070 000 ~ 202 079 999)
	ErrRetryMaxAttemptsExceededCode  = 202070001 // 重试次数超限
	ErrRetryBackoffFailedCode        = 202070002 // 重试退避失败
	ErrRetryConditionNotMetCode      = 202070003 // 重试条件不满足
	ErrRetryConfigInvalidCode        = 202070004 // 重试配置无效

	// =====================================================
	// 向后兼容：保留旧的103段错误码（已废弃）
	// TODO: 逐步迁移到202段错误码，移除103段定义
	// =====================================================
	DeprecatedErrConversationInvalidParamCode = 103000000 // @deprecated 使用 ErrConversationInvalidParamCode (202000003)
	DeprecatedErrConversationPermissionCode   = 103000001 // @deprecated 使用 ErrConversationPermissionDeniedCode (202000008)
	DeprecatedErrConversationNotFound         = 103000002 // @deprecated 使用 ErrConversationNotFoundCode (202000001)
	DeprecatedErrConversationJsonMarshal      = 103000003 // @deprecated
	DeprecatedErrConversationAgentRunError    = 103100001 // @deprecated 使用 ErrAgentRunFailedCode (202020003)
	DeprecatedErrAgentNotExists               = 103100002 // @deprecated
	DeprecatedErrReplyUnknowEventType         = 103100003 // @deprecated
	DeprecatedErrUnknowInterruptType          = 103100004 // @deprecated
	DeprecatedErrInterruptDataEmpty           = 103100005 // @deprecated
	DeprecatedErrConversationMessageNotFound  = 103200001 // @deprecated 使用 ErrMessageNotFoundCode (202010001)
	DeprecatedErrAgentRun                     = 103200002 // @deprecated 使用 ErrAgentRunFailedCode (202020003)
	DeprecatedErrRecordNotFound               = 103200003 // @deprecated
	DeprecatedErrAgentRunWorkflowNotFound     = 103200004 // @deprecated
	DeprecatedErrInProgressCanNotCancel       = 103200005 // @deprecated
)

// 便捷错误变量
var (
	ErrConversationNotFound       = errorx.New(ErrConversationNotFoundCode)
	ErrConversationAlreadyExists  = errorx.New(ErrConversationAlreadyExistsCode)
	ErrConversationInvalidParam   = errorx.New(ErrConversationInvalidParamCode)
	ErrConversationClosed         = errorx.New(ErrConversationClosedCode)
	ErrMessageNotFound            = errorx.New(ErrMessageNotFoundCode)
	ErrMessageSendFailed          = errorx.New(ErrMessageSendFailedCode)
	ErrAgentRunFailed             = errorx.New(ErrAgentRunFailedCode)
	ErrAgentRunTimeout            = errorx.New(ErrAgentRunTimeoutCode)
	ErrStreamConnectionLost       = errorx.New(ErrStreamConnectionLostCode)
)

func init() {
	// =====================================================
	// 202段：标准错误码注册（Conversation模块）
	// =====================================================

	// Conversation基础错误注册
	code.Register(
		ErrConversationNotFoundCode,
		"Conversation not found: {conversation_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrConversationAlreadyExistsCode,
		"Conversation already exists: {conversation_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrConversationInvalidParamCode,
		"Invalid conversation parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrConversationTitleTooLongCode,
		"Conversation title too long: max {max_length} characters, got {actual_length}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrConversationClosedCode,
		"Conversation is closed: {conversation_id}, cannot send new messages",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrConversationDeleteFailedCode,
		"Failed to delete conversation: {conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrConversationUpdateFailedCode,
		"Failed to update conversation: {conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrConversationPermissionDeniedCode,
		"Permission denied for conversation: {conversation_id}, user_id: {user_id}",
		code.WithAffectStability(false),
	)

	// Message错误注册
	code.Register(
		ErrMessageNotFoundCode,
		"Message not found: {message_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrMessageInvalidParamCode,
		"Invalid message parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrMessageContentTooLongCode,
		"Message content too long: max {max_length} characters, got {actual_length}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrMessageContentEmptyCode,
		"Message content cannot be empty",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrMessageSendFailedCode,
		"Failed to send message: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrMessageDeleteFailedCode,
		"Failed to delete message: {message_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrMessageUpdateFailedCode,
		"Failed to update message: {message_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrMessageBatchSendFailedCode,
		"Failed to send batch messages: conversation_id={conversation_id}, count={count}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrMessageTypeNotSupportedCode,
		"Message type not supported: {message_type}",
		code.WithAffectStability(false),
	)

	// Agent运行错误注册
	code.Register(
		ErrAgentRunNotFoundCode,
		"Agent run not found: {run_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAgentRunAlreadyCompletedCode,
		"Agent run already completed: {run_id}, status: {status}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAgentRunFailedCode,
		"Agent run failed: run_id={run_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAgentRunTimeoutCode,
		"Agent run timeout: run_id={run_id}, timeout: {timeout_ms}ms",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAgentRunCancelledCode,
		"Agent run cancelled: run_id={run_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAgentRunStatusInvalidCode,
		"Invalid agent run status: {status}, must be one of: pending, running, completed, failed, cancelled",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAgentRunCreateFailedCode,
		"Failed to create agent run: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAgentRunRetryExceededCode,
		"Agent run retry exceeded: run_id={run_id}, max_retries={max_retries}, actual_retries={actual_retries}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAgentRunQuotaExceededCode,
		"Agent run quota exceeded: tenant_id={tenant_id}, resource_type={resource_type}",
		code.WithAffectStability(true),
	)

	// 流式输出错误注册
	code.Register(
		ErrStreamConnectionLostCode,
		"Stream connection lost: conversation_id={conversation_id}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrStreamTimeoutCode,
		"Stream timeout: conversation_id={conversation_id}, timeout: {timeout_ms}ms",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrStreamInvalidFormatCode,
		"Invalid stream format: {format}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrStreamCreateFailedCode,
		"Failed to create stream: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrStreamWriteFailedCode,
		"Failed to write to stream: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 上下文错误注册
	code.Register(
		ErrContextNotFoundCode,
		"Context not found: conversation_id={conversation_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrContextTooLargeCode,
		"Context too large: max {max_size} bytes, got {actual_size} bytes",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrContextInvalidFormatCode,
		"Invalid context format: {format}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrContextUpdateFailedCode,
		"Failed to update context: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrContextClearFailedCode,
		"Failed to clear context: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 附件错误注册
	code.Register(
		ErrAttachmentNotFoundCode,
		"Attachment not found: {attachment_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAttachmentUploadFailedCode,
		"Failed to upload attachment: conversation_id={conversation_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAttachmentInvalidTypeCode,
		"Invalid attachment type: {file_type}, allowed types: {allowed_types}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAttachmentSizeExceededCode,
		"Attachment size exceeded: max {max_size} bytes, got {actual_size} bytes",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAttachmentDeleteFailedCode,
		"Failed to delete attachment: {attachment_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAttachmentDownloadFailedCode,
		"Failed to download attachment: {attachment_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 反馈错误注册
	code.Register(
		ErrFeedbackNotFoundCode,
		"Feedback not found: {feedback_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrFeedbackInvalidTypeCode,
		"Invalid feedback type: {feedback_type}, must be one of: like, dislike, report",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrFeedbackAlreadyExistsCode,
		"Feedback already exists: message_id={message_id}, user_id={user_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrFeedbackCreateFailedCode,
		"Failed to create feedback: message_id={message_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrFeedbackUpdateFailedCode,
		"Failed to update feedback: {feedback_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// 重试错误注册
	code.Register(
		ErrRetryMaxAttemptsExceededCode,
		"Retry max attempts exceeded: max_attempts={max_attempts}, actual_attempts={actual_attempts}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrRetryBackoffFailedCode,
		"Retry backoff failed: attempt={attempt}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrRetryConditionNotMetCode,
		"Retry condition not met: condition={condition}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrRetryConfigInvalidCode,
		"Invalid retry config: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	// =====================================================
	// 103段：旧错误码注册（向后兼容，已废弃）
	// =====================================================
	code.Register(
		DeprecatedErrInProgressCanNotCancel,
		"in progress can not be cancelled",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrAgentRunWorkflowNotFound,
		"The chatflow is not configured. Please configure it and try again.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrRecordNotFound,
		"record not found or nothing to update",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrAgentRun,
		"Interal Server Error",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrConversationJsonMarshal,
		"json marshal failed",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrConversationNotFound,
		"conversation not found",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationPermissionCode,
		"unauthorized access : {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationInvalidParamCode,
		"invalid parameter : {msg}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationAgentRunError,
		"agent run error : {msg}",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrAgentNotExists,
		"agent not exists",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationMessageNotFound,
		"message not found",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrReplyUnknowEventType,
		"unknow event type",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrUnknowInterruptType,
		"unknow interrupt type",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrInterruptDataEmpty,
		"interrupt data is empty",
		code.WithAffectStability(true),
	)
}
