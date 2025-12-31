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
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx/code"
)

// =====================================================
// 统一错误码规范 - Workflow模块
// 错误码段: 203 000 000 ~ 203 999 999
// =====================================================

const (
	// Workflow基础错误 (203 000 000 ~ 203 009 999)
	ErrWorkflowNotFoundCode            = 203000001 // 工作流不存在
	ErrWorkflowAlreadyExistsCode       = 203000002 // 工作流已存在
	ErrWorkflowInvalidParamCode        = 203000003 // 工作流参数无效
	ErrWorkflowNameTooLongCode         = 203000004 // 工作流名称过长
	ErrWorkflowNotPublishedCode        = 203000005 // 工作流未发布
	ErrWorkflowAlreadyPublishedCode    = 203000006 // 工作流已发布
	ErrWorkflowDeleteFailedCode        = 203000007 // 删除工作流失败
	ErrWorkflowUpdateFailedCode        = 203000008 // 更新工作流失败
	ErrWorkflowPermissionDeniedCode    = 203000009 // 工作流权限不足
	ErrWorkflowConflictParamCode       = 203000010 // 参数冲突：project_id和bot_id不能同时设置
	ErrWorkflowSpaceIDRequiredCode     = 203000011 // space_id必填
	ErrWorkflowSchemaRequiredCode      = 203000012 // Schema必填

	// Node相关 (203 010 000 ~ 203 019 999)
	ErrNodeNotFoundCode               = 203010001 // 节点不存在
	ErrNodeInvalidParamCode           = 203010002 // 节点参数无效
	ErrNodeConfigInvalidCode          = 203010003 // 节点配置无效
	ErrNodeExecuteFailedCode          = 203010004 // 节点执行失败
	ErrNodeTimeoutCode                = 203010005 // 节点执行超时
	ErrNodeCreateFailedCode           = 203010006 // 创建节点失败
	ErrNodeDeleteFailedCode           = 203010007 // 删除节点失败
	ErrNodeUpdateFailedCode           = 203010008 // 更新节点失败
	ErrNodeOutputParseFailedCode      = 203010009 // 节点输出解析失败
	ErrNodeTypeNotSupportedCode       = 203010010 // 节点类型不支持

	// Edge相关 (203 020 000 ~ 203 029 999)
	ErrEdgeNotFoundCode               = 203020001 // 边不存在
	ErrEdgeInvalidParamCode           = 203020002 // 边参数无效
	ErrEdgeConfigInvalidCode          = 203020003 // 边配置无效
	ErrEdgeCreateFailedCode           = 203020004 // 创建边失败
	ErrEdgeDeleteFailedCode           = 203020005 // 删除边失败
	ErrEdgeCycleDetectedCode          = 203020006 // 检测到循环依赖
	ErrEdgeSourceNotFoundCode         = 203020007 // 边的源节点不存在
	ErrEdgeTargetNotFoundCode         = 203020008 // 边的目标节点不存在

	// 执行相关 (203 030 000 ~ 203 039 999)
	ErrExecutionNotFoundCode          = 203030001 // 执行不存在
	ErrExecutionAlreadyCompletedCode  = 203030002 // 执行已完成
	ErrExecutionFailedCode            = 203030003 // 执行失败
	ErrExecutionTimeoutCode           = 203030004 // 执行超时
	ErrExecutionCancelledCode         = 203030005 // 执行已取消
	ErrExecutionStatusInvalidCode     = 203030006 // 执行状态无效
	ErrExecutionCreateFailedCode      = 203030007 // 创建执行失败
	ErrExecutionRetryExceededCode     = 203030008 // 执行重试次数超限
	ErrExecutionQuotaExceededCode     = 203030009 // 执行配额已超限

	// 变量相关 (203 040 000 ~ 203 049 999)
	ErrVariableNotFoundCode           = 203040001 // 变量不存在
	ErrVariableInvalidParamCode       = 203040002 // 变量参数无效
	ErrVariableTypeMismatchCode       = 203040003 // 变量类型不匹配
	ErrVariableReadOnlyCode           = 203040004 // 变量只读
	ErrVariableCreateFailedCode       = 203040005 // 创建变量失败
	ErrVariableDeleteFailedCode       = 203040006 // 删除变量失败
	ErrVariableUpdateFailedCode       = 203040007 // 更新变量失败

	// 触发器相关 (203 050 000 ~ 203 059 999)
	ErrTriggerNotFoundCode            = 203050001 // 触发器不存在
	ErrTriggerInvalidParamCode        = 203050002 // 触发器参数无效
	ErrTriggerConfigInvalidCode       = 203050003 // 触发器配置无效
	ErrTriggerCreateFailedCode        = 203050004 // 创建触发器失败
	ErrTriggerDeleteFailedCode        = 203050005 // 删除触发器失败
	ErrTriggerExecuteFailedCode       = 203050006 // 触发器执行失败

	// 版本管理 (203 060 000 ~ 203 069 999)
	ErrVersionNotFoundCode            = 203060001 // 版本不存在
	ErrVersionAlreadyExistsCode       = 203060002 // 版本已存在
	ErrVersionInvalidParamCode        = 203060003 // 版本参数无效
	ErrVersionNameInvalidCode         = 203060004 // 版本名称无效
	ErrVersionCreateFailedCode        = 203060005 // 创建版本失败
	ErrVersionDeleteFailedCode        = 203060006 // 删除版本失败
	ErrVersionRollbackFailedCode      = 203060007 // 版本回滚失败
	ErrVersionConflictCode            = 203060008 // 版本冲突

	// 调试相关 (203 070 000 ~ 203 079 999)
	ErrDebugModeFailedCode            = 203070001 // 调试模式失败
	ErrBreakpointInvalidCode          = 203070002 // 断点无效
	ErrBreakpointSetFailedCode        = 203070003 // 设置断点失败
	ErrStepExecutionFailedCode        = 203070004 // 单步执行失败
	ErrVariableWatchFailedCode        = 203070005 // 变量监视失败

	// 导入导出 (203 080 000 ~ 203 089 999)
	ErrImportFailedCode               = 203080001 // 导入失败
	ErrExportFailedCode               = 203080002 // 导出失败
	ErrInvalidFormatCode              = 203080003 // 格式无效
	ErrSchemaValidationFailedCode     = 203080004 // Schema验证失败
	ErrSerializationFailedCode        = 203080005 // 序列化失败

	// =====================================================
	// 向后兼容：保留旧的多段错误码（已废弃）
	// TODO: 逐步迁移到203段错误码，移除旧段定义
	// =====================================================

	// 720xxx段（旧Workflow错误码）
	DeprecatedErrWorkflowNotPublished             = 720702011 // @deprecated 使用 ErrWorkflowNotPublishedCode (203000005)
	DeprecatedErrMissingRequiredParam             = 720702002 // @deprecated 使用 ErrWorkflowInvalidParamCode (203000003)
	DeprecatedErrInterruptNotSupported            = 720702078 // @deprecated
	DeprecatedErrInvalidParameter                 = 720702001 // @deprecated 使用 ErrWorkflowInvalidParamCode (203000003)
	DeprecatedErrArrIndexOutOfRange               = 720712014 // @deprecated
	DeprecatedErrWorkflowExecuteFail              = 720701013 // @deprecated 使用 ErrExecutionFailedCode (203030003)
	DeprecatedErrQuestionOptionsEmpty             = 720712049 // @deprecated
	DeprecatedErrNodeOutputParseFail              = 720712023 // @deprecated 使用 ErrNodeOutputParseFailedCode (203010009)
	DeprecatedErrWorkflowTimeout                  = 720702085 // @deprecated 使用 ErrExecutionTimeoutCode (203030004)
	DeprecatedErrWorkflowNotFound                 = 720702004 // @deprecated 使用 ErrWorkflowNotFoundCode (203000001)
	DeprecatedErrSerializationDeserializationFail = 720701011 // @deprecated 使用 ErrSerializationFailedCode (203080005)
	DeprecatedErrInternalBadRequest               = 720701007 // @deprecated
	DeprecatedErrSchemaConversionFail             = 720702089 // @deprecated 使用 ErrSchemaValidationFailedCode (203080004)
	DeprecatedErrWorkflowCompileFail              = 720701003 // @deprecated
	DeprecatedErrPluginAPIErr                     = 720701004 // @deprecated
	DeprecatedErrConversationNameIsDuplicated     = 720702200 // @deprecated
	DeprecatedErrConversationOfAppNotFound        = 720702201 // @deprecated
	DeprecatedErrConversationNodeInvalidOperation = 720702250 // @deprecated
	DeprecatedErrOnlyDefaultConversationAllowInAgentScenario = 720712033 // @deprecated
	DeprecatedErrConversationNodesNotAvailable    = 702093204 // @deprecated

	// 777xxx段（旧Workflow错误码）
	DeprecatedErrConversationNodeOperationFail    = 777777782 // @deprecated
	DeprecatedErrMessageNodeOperationFail         = 777777781 // @deprecated
	DeprecatedErrChatFlowRoleOperationFail        = 777777780 // @deprecated
	DeprecatedErrConversationOfAppOperationFail   = 777777779 // @deprecated
	DeprecatedErrWorkflowSpecifiedVersionNotFound = 777777778 // @deprecated 使用 ErrVersionNotFoundCode (203060001)
	DeprecatedErrWorkflowCanceledByUser           = 777777777 // @deprecated 使用 ErrExecutionCancelledCode (203030005)
	DeprecatedErrNodeTimeout                      = 777777776 // @deprecated 使用 ErrNodeTimeoutCode (203010005)
	DeprecatedErrWorkflowOperationFail            = 777777775 // @deprecated
	DeprecatedErrIndexingNilArray                 = 777777774 // @deprecated
	DeprecatedErrLLMStructuredOutputParseFail     = 777777773 // @deprecated
	DeprecatedErrCreateNodeFail                   = 777777772 // @deprecated 使用 ErrNodeCreateFailedCode (203010006)
	DeprecatedErrWorkflowSnapshotNotFound         = 777777771 // @deprecated
	DeprecatedErrNotifyWorkflowResourceChangeErr  = 777777770 // @deprecated
	DeprecatedErrInvalidVersionName               = 777777769 // @deprecated 使用 ErrVersionNameInvalidCode (203060004)
	DeprecatedErrPluginIDNotFound                 = 777777768 // @deprecated
	DeprecatedErrTOSError                         = 777777767 // @deprecated
	DeprecatedErrToolIDNotFound                   = 777777766 // @deprecated
	DeprecatedErrAuthorizationRequired            = 777777765 // @deprecated
	DeprecatedErrVariablesAPIFail                 = 777777764 // @deprecated
	DeprecatedErrInputFieldMissing                = 777777763 // @deprecated
	DeprecatedErrConversationNotFoundForOperation = 777777762 // @deprecated

	// 其他段（稳定性问题）
	DeprecatedErrDatabaseError = 720700801 // @deprecated
	DeprecatedErrRedisError    = 720700803 // @deprecated
	DeprecatedErrIDGenError    = 720700808 // @deprecated
	DeprecatedErrCodeExecuteFail = 305000002 // @deprecated

	// OpenAPI错误码（保留用于映射）
	ErrOpenAPIWorkflowNotPublished  = 6031
	ErrOpenAPIBadRequest            = 4000
	ErrOpenAPIInterruptNotSupported = 6039
	ErrOpenAPIWorkflowTimeout       = 6023
)

// 便捷错误变量（用于错误处理）
var (
	// Workflow基础错误 (203 000 000 ~ 203 009 999)
	ErrWorkflowNotFound            = errorx.New(ErrWorkflowNotFoundCode)
	ErrWorkflowAlreadyExists       = errorx.New(ErrWorkflowAlreadyExistsCode)
	ErrWorkflowInvalidParam        = errorx.New(ErrWorkflowInvalidParamCode)
	ErrWorkflowNameTooLong         = errorx.New(ErrWorkflowNameTooLongCode)
	ErrWorkflowNotPublished        = errorx.New(ErrWorkflowNotPublishedCode)
	ErrWorkflowAlreadyPublished    = errorx.New(ErrWorkflowAlreadyPublishedCode)
	ErrWorkflowDeleteFailed        = errorx.New(ErrWorkflowDeleteFailedCode)
	ErrWorkflowUpdateFailed        = errorx.New(ErrWorkflowUpdateFailedCode)
	ErrWorkflowPermissionDenied    = errorx.New(ErrWorkflowPermissionDeniedCode)
	ErrWorkflowConflictParam       = errorx.New(ErrWorkflowConflictParamCode)
	ErrWorkflowSpaceIDRequired     = errorx.New(ErrWorkflowSpaceIDRequiredCode)
	ErrWorkflowSchemaRequired      = errorx.New(ErrWorkflowSchemaRequiredCode)

	// Node相关 (203 010 000 ~ 203 019 999)
	ErrNodeNotFound          = errorx.New(ErrNodeNotFoundCode)
	ErrNodeInvalidParam      = errorx.New(ErrNodeInvalidParamCode)
	ErrNodeConfigInvalid     = errorx.New(ErrNodeConfigInvalidCode)
	ErrNodeExecuteFailed     = errorx.New(ErrNodeExecuteFailedCode)
	ErrNodeTimeout           = errorx.New(ErrNodeTimeoutCode)
	ErrNodeCreateFailed      = errorx.New(ErrNodeCreateFailedCode)
	ErrNodeDeleteFailed      = errorx.New(ErrNodeDeleteFailedCode)
	ErrNodeUpdateFailed      = errorx.New(ErrNodeUpdateFailedCode)
	ErrNodeOutputParseFailed = errorx.New(ErrNodeOutputParseFailedCode)
	ErrNodeTypeNotSupported  = errorx.New(ErrNodeTypeNotSupportedCode)

	// Edge相关 (203 020 000 ~ 203 029 999)
	ErrEdgeNotFound       = errorx.New(ErrEdgeNotFoundCode)
	ErrEdgeInvalidParam   = errorx.New(ErrEdgeInvalidParamCode)
	ErrEdgeConfigInvalid  = errorx.New(ErrEdgeConfigInvalidCode)
	ErrEdgeCreateFailed   = errorx.New(ErrEdgeCreateFailedCode)
	ErrEdgeDeleteFailed   = errorx.New(ErrEdgeDeleteFailedCode)
	ErrEdgeCycleDetected  = errorx.New(ErrEdgeCycleDetectedCode)
	ErrEdgeSourceNotFound = errorx.New(ErrEdgeSourceNotFoundCode)
	ErrEdgeTargetNotFound = errorx.New(ErrEdgeTargetNotFoundCode)

	// 执行相关 (203 030 000 ~ 203 039 999)
	ErrExecutionNotFound         = errorx.New(ErrExecutionNotFoundCode)
	ErrExecutionAlreadyCompleted = errorx.New(ErrExecutionAlreadyCompletedCode)
	ErrExecutionFailed           = errorx.New(ErrExecutionFailedCode)
	ErrExecutionTimeout          = errorx.New(ErrExecutionTimeoutCode)
	ErrExecutionCancelled        = errorx.New(ErrExecutionCancelledCode)
	ErrExecutionStatusInvalid    = errorx.New(ErrExecutionStatusInvalidCode)
	ErrExecutionCreateFailed     = errorx.New(ErrExecutionCreateFailedCode)
	ErrExecutionRetryExceeded    = errorx.New(ErrExecutionRetryExceededCode)
	ErrExecutionQuotaExceeded    = errorx.New(ErrExecutionQuotaExceededCode)

	// 变量相关 (203 040 000 ~ 203 049 999)
	ErrVariableNotFound     = errorx.New(ErrVariableNotFoundCode)
	ErrVariableInvalidParam = errorx.New(ErrVariableInvalidParamCode)
	ErrVariableTypeMismatch = errorx.New(ErrVariableTypeMismatchCode)
	ErrVariableReadOnly     = errorx.New(ErrVariableReadOnlyCode)
	ErrVariableCreateFailed = errorx.New(ErrVariableCreateFailedCode)
	ErrVariableDeleteFailed = errorx.New(ErrVariableDeleteFailedCode)
	ErrVariableUpdateFailed = errorx.New(ErrVariableUpdateFailedCode)

	// 触发器相关 (203 050 000 ~ 203 059 999)
	ErrTriggerNotFound     = errorx.New(ErrTriggerNotFoundCode)
	ErrTriggerInvalidParam = errorx.New(ErrTriggerInvalidParamCode)
	ErrTriggerConfigInvalid = errorx.New(ErrTriggerConfigInvalidCode)
	ErrTriggerCreateFailed  = errorx.New(ErrTriggerCreateFailedCode)
	ErrTriggerDeleteFailed  = errorx.New(ErrTriggerDeleteFailedCode)
	ErrTriggerExecuteFailed = errorx.New(ErrTriggerExecuteFailedCode)

	// 版本管理 (203 060 000 ~ 203 069 999)
	ErrVersionNotFound       = errorx.New(ErrVersionNotFoundCode)
	ErrVersionAlreadyExists  = errorx.New(ErrVersionAlreadyExistsCode)
	ErrVersionInvalidParam   = errorx.New(ErrVersionInvalidParamCode)
	ErrVersionNameInvalid    = errorx.New(ErrVersionNameInvalidCode)
	ErrVersionCreateFailed   = errorx.New(ErrVersionCreateFailedCode)
	ErrVersionDeleteFailed   = errorx.New(ErrVersionDeleteFailedCode)
	ErrVersionRollbackFailed = errorx.New(ErrVersionRollbackFailedCode)
	ErrVersionConflict       = errorx.New(ErrVersionConflictCode)

	// 调试相关 (203 070 000 ~ 203 079 999)
	ErrDebugModeFailed      = errorx.New(ErrDebugModeFailedCode)
	ErrBreakpointInvalid    = errorx.New(ErrBreakpointInvalidCode)
	ErrBreakpointSetFailed  = errorx.New(ErrBreakpointSetFailedCode)
	ErrStepExecutionFailed  = errorx.New(ErrStepExecutionFailedCode)
	ErrVariableWatchFailed  = errorx.New(ErrVariableWatchFailedCode)

	// 导入导出 (203 080 000 ~ 203 089 999)
	ErrImportFailed           = errorx.New(ErrImportFailedCode)
	ErrExportFailed           = errorx.New(ErrExportFailedCode)
	ErrInvalidFormat          = errorx.New(ErrInvalidFormatCode)
	ErrSchemaValidationFailed = errorx.New(ErrSchemaValidationFailedCode)
	ErrSerializationFailed    = errorx.New(ErrSerializationFailedCode)

	// 已弃用错误码的便捷别名（向后兼容）
	// TODO: 逐步迁移到新的203段错误码，移除这些别名
	ErrVariablesAPIFail                        = errorx.New(DeprecatedErrVariablesAPIFail)
	ErrPluginIDNotFound                        = errorx.New(DeprecatedErrPluginIDNotFound)
	ErrPluginAPIErr                            = errorx.New(DeprecatedErrPluginAPIErr)
	ErrTOSError                                = errorx.New(DeprecatedErrTOSError)
	ErrSchemaConversionFail                    = errorx.New(DeprecatedErrSchemaConversionFail)
	ErrMissingRequiredParam                    = errorx.New(DeprecatedErrMissingRequiredParam)
	ErrInvalidParameter                        = errorx.New(DeprecatedErrInvalidParameter)
	ErrSerializationDeserializationFail        = errorx.New(DeprecatedErrSerializationDeserializationFail)
	ErrNodeOutputParseFail                     = errorx.New(DeprecatedErrNodeOutputParseFail)
	ErrInputFieldMissing                       = errorx.New(DeprecatedErrInputFieldMissing)
	ErrLLMStructuredOutputParseFail            = errorx.New(DeprecatedErrLLMStructuredOutputParseFail)
	ErrConversationNodesNotAvailable           = errorx.New(DeprecatedErrConversationNodesNotAvailable)
	ErrConversationNodeOperationFail           = errorx.New(DeprecatedErrConversationNodeOperationFail)
	ErrAuthorizationRequired                   = errorx.New(DeprecatedErrAuthorizationRequired)
	ErrOnlyDefaultConversationAllowInAgentScenario = errorx.New(DeprecatedErrOnlyDefaultConversationAllowInAgentScenario)
	ErrInterruptNotSupported                   = errorx.New(DeprecatedErrInterruptNotSupported) // @deprecated 使用 DeprecatedErrInterruptNotSupported
	ErrCreateNodeFail                          = errorx.New(DeprecatedErrCreateNodeFail)     // @deprecated 使用 ErrNodeCreateFailedCode (203010006)
	ErrDatabaseError                           = errorx.New(DeprecatedErrDatabaseError)       // @deprecated
	ErrRedisError                              = errorx.New(DeprecatedErrRedisError)          // @deprecated
	ErrWorkflowCanceledByUser                  = errorx.New(DeprecatedErrWorkflowCanceledByUser) // @deprecated 使用 ErrExecutionCancelledCode (203030005)

	// 新增错误码变量（用于向后兼容）
	ErrWorkflowExecuteFail              = errorx.New(DeprecatedErrWorkflowExecuteFail)              // @deprecated 使用 ErrExecutionFailedCode (203030003)
	ErrChatFlowRoleOperationFail        = errorx.New(DeprecatedErrChatFlowRoleOperationFail)        // @deprecated
	ErrWorkflowOperationFail            = errorx.New(DeprecatedErrWorkflowOperationFail)            // @deprecated
	ErrConversationOfAppOperationFail   = errorx.New(DeprecatedErrConversationOfAppOperationFail)   // @deprecated
	ErrConversationNotFoundForOperation = errorx.New(DeprecatedErrConversationNotFoundForOperation) // @deprecated

	// 补充缺失的已废弃错误码便捷变量（向后兼容）
	// 720xxx段
	ErrArrIndexOutOfRange             = errorx.New(DeprecatedErrArrIndexOutOfRange)           // @deprecated
	ErrQuestionOptionsEmpty           = errorx.New(DeprecatedErrQuestionOptionsEmpty)         // @deprecated
	ErrWorkflowTimeoutDeprecated      = errorx.New(DeprecatedErrWorkflowTimeout)              // @deprecated 使用 ErrExecutionTimeoutCode (203030004) - 注：与新版 ErrWorkflowTimeout 避免命名冲突
	ErrInternalBadRequest             = errorx.New(DeprecatedErrInternalBadRequest)           // @deprecated
	ErrWorkflowCompileFail            = errorx.New(DeprecatedErrWorkflowCompileFail)          // @deprecated
	ErrConversationNameIsDuplicated   = errorx.New(DeprecatedErrConversationNameIsDuplicated) // @deprecated
	ErrConversationOfAppNotFound      = errorx.New(DeprecatedErrConversationOfAppNotFound)    // @deprecated
	ErrConversationNodeInvalidOperation = errorx.New(DeprecatedErrConversationNodeInvalidOperation) // @deprecated

	// 777xxx段
	ErrWorkflowSpecifiedVersionNotFound = errorx.New(DeprecatedErrWorkflowSpecifiedVersionNotFound) // @deprecated 使用 ErrVersionNotFoundCode (203060001)
	ErrNodeTimeoutDeprecated            = errorx.New(DeprecatedErrNodeTimeout)                       // @deprecated 使用 ErrNodeTimeoutCode (203010005) - 注：与新版 ErrNodeTimeout 避免命名冲突
	ErrIndexingNilArray                 = errorx.New(DeprecatedErrIndexingNilArray)                  // @deprecated
	ErrWorkflowSnapshotNotFound         = errorx.New(DeprecatedErrWorkflowSnapshotNotFound)          // @deprecated
	ErrNotifyWorkflowResourceChangeErr  = errorx.New(DeprecatedErrNotifyWorkflowResourceChangeErr)   // @deprecated
	ErrInvalidVersionNameDeprecated     = errorx.New(DeprecatedErrInvalidVersionName)                // @deprecated 使用 ErrVersionNameInvalidCode (203060004) - 注：与新版 ErrVersionNameInvalid 避免命名冲突
	ErrToolIDNotFound                   = errorx.New(DeprecatedErrToolIDNotFound)                    // @deprecated
	ErrMessageNodeOperationFail         = errorx.New(DeprecatedErrMessageNodeOperationFail)          // @deprecated
)

func init() {
	// =====================================================
	// 203段：标准错误码注册（Workflow模块）
	// =====================================================

	// Workflow基础错误注册
	code.Register(
		ErrWorkflowNotFoundCode,
		"Workflow not found: {workflow_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowAlreadyExistsCode,
		"Workflow already exists: {workflow_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowInvalidParamCode,
		"Invalid workflow parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowNameTooLongCode,
		"Workflow name too long: max {max_length} characters, got {actual_length}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowNotPublishedCode,
		"Workflow not published: {workflow_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowAlreadyPublishedCode,
		"Workflow already published: {workflow_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowDeleteFailedCode,
		"Failed to delete workflow: {workflow_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrWorkflowUpdateFailedCode,
		"Failed to update workflow: {workflow_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrWorkflowPermissionDeniedCode,
		"Permission denied for workflow: {workflow_id}, user_id: {user_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowConflictParamCode,
		"project_id and bot_id cannot be set at the same time",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowSpaceIDRequiredCode,
		"space id is required",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrWorkflowSchemaRequiredCode,
		"validate tree schema is required",
		code.WithAffectStability(false),
	)

	// Node错误注册
	code.Register(
		ErrNodeNotFoundCode,
		"Node not found: {node_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrNodeInvalidParamCode,
		"Invalid node parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrNodeConfigInvalidCode,
		"Invalid node config: node_id={node_id}, error: {error}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrNodeExecuteFailedCode,
		"Node execution failed: node_id={node_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrNodeTimeoutCode,
		"Node timeout: node_id={node_id}, timeout: {timeout_ms}ms",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrNodeCreateFailedCode,
		"Failed to create node: {node_name}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrNodeDeleteFailedCode,
		"Failed to delete node: {node_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrNodeUpdateFailedCode,
		"Failed to update node: {node_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrNodeOutputParseFailedCode,
		"Node output parse failed: node_id={node_id}, warnings: {warnings}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrNodeTypeNotSupportedCode,
		"Node type not supported: {node_type}",
		code.WithAffectStability(false),
	)

	// Edge错误注册
	code.Register(
		ErrEdgeNotFoundCode,
		"Edge not found: {edge_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrEdgeInvalidParamCode,
		"Invalid edge parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrEdgeConfigInvalidCode,
		"Invalid edge config: edge_id={edge_id}, error: {error}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrEdgeCreateFailedCode,
		"Failed to create edge: error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrEdgeDeleteFailedCode,
		"Failed to delete edge: {edge_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrEdgeCycleDetectedCode,
		"Cycle detected in workflow: path={path}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrEdgeSourceNotFoundCode,
		"Edge source node not found: edge_id={edge_id}, source_id={source_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrEdgeTargetNotFoundCode,
		"Edge target node not found: edge_id={edge_id}, target_id={target_id}",
		code.WithAffectStability(false),
	)

	// Execution错误注册
	code.Register(
		ErrExecutionNotFoundCode,
		"Execution not found: {execution_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrExecutionAlreadyCompletedCode,
		"Execution already completed: {execution_id}, status: {status}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrExecutionFailedCode,
		"Execution failed: execution_id={execution_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrExecutionTimeoutCode,
		"Execution timeout: execution_id={execution_id}, timeout: {timeout_ms}ms",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrExecutionCancelledCode,
		"Execution cancelled: execution_id={execution_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrExecutionStatusInvalidCode,
		"Invalid execution status: {status}, must be one of: pending, running, completed, failed, cancelled",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrExecutionCreateFailedCode,
		"Failed to create execution: workflow_id={workflow_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrExecutionRetryExceededCode,
		"Execution retry exceeded: execution_id={execution_id}, max_retries={max_retries}, actual_retries={actual_retries}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrExecutionQuotaExceededCode,
		"Execution quota exceeded: tenant_id={tenant_id}, resource_type={resource_type}",
		code.WithAffectStability(true),
	)

	// Variable错误注册
	code.Register(
		ErrVariableNotFoundCode,
		"Variable not found: {variable_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVariableInvalidParamCode,
		"Invalid variable parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVariableTypeMismatchCode,
		"Variable type mismatch: variable_name={variable_name}, expected_type={expected_type}, actual_type={actual_type}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVariableReadOnlyCode,
		"Variable is read-only: {variable_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVariableCreateFailedCode,
		"Failed to create variable: {variable_name}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrVariableDeleteFailedCode,
		"Failed to delete variable: {variable_name}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrVariableUpdateFailedCode,
		"Failed to update variable: {variable_name}, error: {error}",
		code.WithAffectStability(true),
	)

	// Trigger错误注册
	code.Register(
		ErrTriggerNotFoundCode,
		"Trigger not found: {trigger_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTriggerInvalidParamCode,
		"Invalid trigger parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTriggerConfigInvalidCode,
		"Invalid trigger config: trigger_id={trigger_id}, error: {error}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrTriggerCreateFailedCode,
		"Failed to create trigger: error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTriggerDeleteFailedCode,
		"Failed to delete trigger: {trigger_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrTriggerExecuteFailedCode,
		"Trigger execution failed: trigger_id={trigger_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// Version错误注册
	code.Register(
		ErrVersionNotFoundCode,
		"Version not found: {version_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVersionAlreadyExistsCode,
		"Version already exists: workflow_id={workflow_id}, version_name={version_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVersionInvalidParamCode,
		"Invalid version parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVersionNameInvalidCode,
		"Invalid version name: {version_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrVersionCreateFailedCode,
		"Failed to create version: workflow_id={workflow_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrVersionDeleteFailedCode,
		"Failed to delete version: {version_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrVersionRollbackFailedCode,
		"Failed to rollback version: workflow_id={workflow_id}, version_id={version_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrVersionConflictCode,
		"Version conflict: workflow_id={workflow_id}, please resolve conflicts before committing",
		code.WithAffectStability(false),
	)

	// Debug错误注册
	code.Register(
		ErrDebugModeFailedCode,
		"Debug mode failed: workflow_id={workflow_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBreakpointInvalidCode,
		"Invalid breakpoint: node_id={node_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBreakpointSetFailedCode,
		"Failed to set breakpoint: node_id={node_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrStepExecutionFailedCode,
		"Step execution failed: execution_id={execution_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrVariableWatchFailedCode,
		"Variable watch failed: variable_name={variable_name}, error: {error}",
		code.WithAffectStability(true),
	)

	// Import/Export错误注册
	code.Register(
		ErrImportFailedCode,
		"Import failed: format={format}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrExportFailedCode,
		"Export failed: workflow_id={workflow_id}, format={format}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrInvalidFormatCode,
		"Invalid format: {format}, supported formats: {supported_formats}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSchemaValidationFailedCode,
		"Schema validation failed: error: {error}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSerializationFailedCode,
		"Serialization failed: error: {error}",
		code.WithAffectStability(true),
	)

	// =====================================================
	// 旧段错误码注册（向后兼容，已废弃）
	// =====================================================

	// 720xxx段错误码注册
	code.Register(
		DeprecatedErrWorkflowNotPublished,
		"Workflow not published. The requested operation cannot be performed on an unpublished workflow. Please publish the workflow and try again.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrMissingRequiredParam,
		"Missing required parameters：'{param}'. Please review the API documentation and ensure all mandatory fields are included in your request.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrInterruptNotSupported,
		"Synchronous requests do not support interruption. Please switch to asynchronous requests for interruptible operations.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrInvalidParameter,
		"Invalid request parameters. Please check your input and ensure all required fields are correctly formatted and within allowed ranges.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrArrIndexOutOfRange,
		"Array {arr_name} index out of bounds: The requested index {req_index} exceeds the array's length {arr_len}. Please ensure the index is within the valid range of the array. You can refer to debug_url for more details.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrWorkflowExecuteFail,
		"Workflow execution failure: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrQuestionOptionsEmpty,
		"question option is empty",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrNodeOutputParseFail,
		"node output parse fail: {warnings}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrWorkflowTimeout,
		"Workflow execution timed out. Please check for long-running operations, optimize if possible, or retry later.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrWorkflowNotFound,
		"workflow {id} not found, please check if the workflow exists and not deleted",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrSerializationDeserializationFail,
		"data serialization/deserialization fail, please contact support team",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrInternalBadRequest,
		"one of the request parameters for {scene} is invalid, please contact support team",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrSchemaConversionFail,
		"schema conversion failed, please contact support team",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrWorkflowCompileFail,
		"workflow compile failed, please contact support team",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrPluginAPIErr,
		"plugin api error",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationNameIsDuplicated,
		"conversation name {name} is duplicated",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationOfAppNotFound,
		"conversation not found, please check if the application conversation exists",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationNodeInvalidOperation,
		"Only conversation created through nodes are allowed to be modified or deleted.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrOnlyDefaultConversationAllowInAgentScenario,
		"Only default conversation allow in agent scenario",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationNodesNotAvailable,
		"Conversation nodes are unavailable in agent scenarios and require an app binding.",
		code.WithAffectStability(false),
	)

	// 777xxx段错误码注册
	code.Register(
		DeprecatedErrIndexingNilArray,
		"Array {arr_name} is nil: The requested index {req_index} cannot be extracted.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrWorkflowOperationFail,
		"Workflow operation failure: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrChatFlowRoleOperationFail,
		"ChatFlowRole operation failure: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrMessageNodeOperationFail,
		"Message node operation failure: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationNodeOperationFail,
		"Conversation node operation failure: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrCodeExecuteFail,
		"Function execution failed, please check the code of the function. Detail: {detail}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrWorkflowCanceledByUser,
		"workflow cancel by user",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationOfAppOperationFail,
		"Conversation management operation failure: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrNodeTimeout,
		"node timeout",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrLLMStructuredOutputParseFail,
		"parse LLM structured output failed, please refer to LLM's raw output for detail.",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrCreateNodeFail,
		"create node {node_name} failed: {cause}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrDatabaseError,
		"database operation failed",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrRedisError,
		"redis operation failed",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrIDGenError,
		"id generator failed",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrWorkflowSnapshotNotFound,
		"workflow {id} snapshot {commit_id} not found, please contact support team",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrNotifyWorkflowResourceChangeErr,
		"notify workflow resource change failed, please try again later",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrInvalidVersionName,
		"workflow version name is invalid",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrPluginIDNotFound,
		"plugin {id} not found",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrTOSError,
		"tos operation failed",
		code.WithAffectStability(true),
	)

	code.Register(
		DeprecatedErrToolIDNotFound,
		"tool {id} not found",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrAuthorizationRequired,
		"authorization required: {extra}",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrVariablesAPIFail,
		"variables API failed",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrInputFieldMissing,
		"input field {name} not found",
		code.WithAffectStability(false),
	)

	code.Register(
		DeprecatedErrConversationNotFoundForOperation,
		"Conversation not found. Please create a conversation before attempting to perform any related operations.",
		code.WithAffectStability(false),
	)
}

// OpenAPI错误码映射表（保留兼容性）
var errnoMap = map[int]int{
	// 203段 → OpenAPI映射（新增）
	ErrWorkflowNotPublishedCode: ErrOpenAPIWorkflowNotPublished,
	ErrWorkflowInvalidParamCode:  ErrOpenAPIBadRequest,
	ErrExecutionTimeoutCode:      ErrOpenAPIWorkflowTimeout,

	// 旧段 → OpenAPI映射（保留）
	DeprecatedErrWorkflowNotPublished:  ErrOpenAPIWorkflowNotPublished,
	DeprecatedErrMissingRequiredParam:  ErrOpenAPIBadRequest,
	DeprecatedErrInterruptNotSupported: ErrOpenAPIInterruptNotSupported,
	DeprecatedErrInvalidParameter:      ErrOpenAPIBadRequest,
	DeprecatedErrArrIndexOutOfRange:    ErrOpenAPIBadRequest,
	DeprecatedErrWorkflowTimeout:       ErrOpenAPIWorkflowTimeout,
}

// CodeForOpenAPI 将内部错误码转换为OpenAPI错误码
// 保留此函数以维持向后兼容性
func CodeForOpenAPI(err errorx.StatusError) int {
	if err == nil {
		return 0
	}

	if c, ok := errnoMap[int(err.Code())]; ok {
		return c
	}

	return int(err.Code())
}
