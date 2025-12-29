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

// Bot: 201 000 000 ~ 201 999 999
const (
	// Bot基础错误 (201 000 000 ~ 201 009 999)
	ErrBotNotFoundCode          = 201000001 // Bot不存在
	ErrBotAlreadyExistsCode     = 201000002 // Bot已存在
	ErrBotInvalidParamCode      = 201000003 // Bot参数无效
	ErrBotNameTooLongCode       = 201000004 // Bot名称过长
	ErrBotDescriptionTooLongCode = 201000005 // Bot描述过长
	ErrBotStatusInvalidCode     = 201000006 // Bot状态无效
	ErrBotOperationNotAllowedCode = 201000007 // 操作不允许
	ErrBotDeleteFailedCode      = 201000008 // 删除Bot失败
	ErrBotUpdateFailedCode      = 201000009 // 更新Bot失败

	// Bot发布相关 (201 010 000 ~ 201 019 999)
	ErrBotNotPublishedCode       = 201010001 // Bot未发布
	ErrBotAlreadyPublishedCode   = 201010002 // Bot已发布
	ErrBotPublishFailedCode      = 201010003 // 发布Bot失败
	ErrBotUnpublishFailedCode    = 201010004 // 取消发布失败
	ErrBotVersionNotFoundCode    = 201010005 // Bot版本不存在

	// Bot配置相关 (201 020 000 ~ 201 029 999)
	ErrBotConfigNotFoundCode     = 201020001 // Bot配置不存在
	ErrBotConfigInvalidCode      = 201020002 // Bot配置无效
	ErrBotConfigUpdateFailedCode = 201020003 // 更新Bot配置失败
	ErrBotVariableNotFoundCode   = 201020004 // Bot变量不存在
	ErrBotVariableInvalidCode    = 201020005 // Bot变量无效

	// Bot组件相关 (201 030 000 ~ 201 039 999)
	ErrBotComponentNotFoundCode   = 201030001 // Bot组件不存在
	ErrBotComponentInvalidCode    = 201030002 // Bot组件无效
	ErrBotComponentAddFailedCode   = 201030003 // 添加组件失败
	ErrBotComponentRemoveFailedCode = 201030004 // 删除组件失败
	ErrBotComponentConnectFailedCode = 201030005 // 连接组件失败

	// Bot插件相关 (201 040 000 ~ 201 049 999)
	ErrBotPluginNotFoundCode      = 201040001 // Bot插件不存在
	ErrBotPluginAlreadyAddedCode  = 201040002 // Bot插件已添加
	ErrBotPluginAddFailedCode     = 201040003 // 添加插件失败
	ErrBotPluginRemoveFailedCode  = 201040004 // 移除插件失败
	ErrBotPluginConfigFailedCode  = 201040005 // 配置插件失败

	// Bot知识库相关 (201 050 000 ~ 201 059 999)
	ErrBotKnowledgeNotFoundCode     = 201050001 // Bot知识库不存在
	ErrBotKnowledgeAlreadyAddedCode  = 201050002 // Bot知识库已添加
	ErrBotKnowledgeAddFailedCode     = 201050003 // 添加知识库失败
	ErrBotKnowledgeRemoveFailedCode  = 201050004 // 移除知识库失败
	ErrBotKnowledgeUpdateFailedCode  = 201050005 // 更新知识库失败

	// Bot工作流相关 (201 060 000 ~ 201 069 999)
	ErrBotWorkflowNotFoundCode    = 201060001 // Bot工作流不存在
	ErrBotWorkflowInvalidCode     = 201060002 // Bot工作流无效
	ErrBotWorkflowAddFailedCode    = 201060003 // 添加工作流失败
	ErrBotWorkflowRemoveFailedCode = 201060004 // 移除工作流失败
	ErrBotWorkflowUpdateFailedCode = 201060005 // 更新工作流失败

	// Bot执行相关 (201 070 000 ~ 201 079 999)
	ErrBotExecuteFailedCode       = 201070001 // 执行Bot失败
	ErrBotTimeoutCode              = 201070002 // Bot执行超时
	ErrBotQuotaExceededCode        = 201070003 // Bot配额已超限
	ErrBotPermissionDeniedCode     = 201070004 // Bot权限不足
	ErrBotDependencyNotFoundCode   = 201070005 // Bot依赖不存在

	// Bot测试相关 (201 080 000 ~ 201 089 999)
	ErrBotTestNotFoundCode         = 201080001 // Bot测试不存在
	ErrBotTestCreateFailedCode     = 201080002 // 创建Bot测试失败
	ErrBotTestExecuteFailedCode    = 201080003 // 执行Bot测试失败
	ErrBotTestTimeoutCode          = 201080004 // Bot测试超时
	ErrBotTestResultInvalidCode     = 201080005 // Bot测试结果无效

	// Bot版本管理 (201 090 000 ~ 201 099 999)
	ErrBotVersionConflictCode      = 201090001 // Bot版本冲突
	ErrBotRollbackFailedCode       = 201090002 // Bot回滚失败
	ErrBotRestoreFailedCode        = 201090003 // Bot恢复失败
	ErrBotVersionDeleteFailedCode   = 201090004 // 删除Bot版本失败
)

// 便捷错误变量
var (
	ErrBotNotFound          = errorx.New(ErrBotNotFoundCode)
	ErrBotAlreadyExists     = errorx.New(ErrBotAlreadyExistsCode)
	ErrBotInvalidParam      = errorx.New(ErrBotInvalidParamCode)
	ErrBotNotPublished       = errorx.New(ErrBotNotPublishedCode)
	ErrBotAlreadyPublished   = errorx.New(ErrBotAlreadyPublishedCode)
	ErrBotConfigInvalid      = errorx.New(ErrBotConfigInvalidCode)
	ErrBotExecuteFailed      = errorx.New(ErrBotExecuteFailedCode)
	ErrBotTimeout            = errorx.New(ErrBotTimeoutCode)
	ErrBotQuotaExceeded      = errorx.New(ErrBotQuotaExceededCode)
	ErrBotPermissionDenied    = errorx.New(ErrBotPermissionDeniedCode)
)

func init() {
	// Bot基础错误注册
	code.Register(
		ErrBotNotFoundCode,
		"Bot not found: {bot_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotAlreadyExistsCode,
		"Bot already exists: {bot_name}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotInvalidParamCode,
		"Invalid bot parameter: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotNameTooLongCode,
		"Bot name too long: max {max_length} characters, got {actual_length}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotStatusInvalidCode,
		"Invalid bot status: {status}, must be one of: draft, published, archived",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotOperationNotAllowedCode,
		"Operation not allowed for bot: {bot_id}, operation: {operation}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotDeleteFailedCode,
		"Failed to delete bot: {bot_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBotUpdateFailedCode,
		"Failed to update bot: {bot_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// Bot发布错误注册
	code.Register(
		ErrBotNotPublishedCode,
		"Bot is not published: {bot_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotAlreadyPublishedCode,
		"Bot already published: {bot_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrBotPublishFailedCode,
		"Failed to publish bot: {bot_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBotUnpublishFailedCode,
		"Failed to unpublish bot: {bot_id}, error: {error}",
		code.WithAffectStability(true),
	)

	// Bot执行错误注册
	code.Register(
		ErrBotExecuteFailedCode,
		"Failed to execute bot: {bot_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBotTimeoutCode,
		"Bot execution timeout: {bot_id}, timeout: {timeout_ms}ms",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBotQuotaExceededCode,
		"Bot quota exceeded: tenant_id={tenant_id}, bot_id={bot_id}, resource_type={resource_type}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrBotPermissionDeniedCode,
		"Permission denied for bot: {bot_id}, operation: {operation}",
		code.WithAffectStability(true),
	)
}
