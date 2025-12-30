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
