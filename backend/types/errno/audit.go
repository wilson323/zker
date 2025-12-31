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
// 统一错误码规范 - Audit模块
// 错误码段: 208 000 000 ~ 208 999 999
// =====================================================

const (
	// Audit基础错误 (208 000 000 ~ 208 009 999)
	ErrAuditLogNotFoundCode         = 208000001 // 审计日志不存在
	ErrAuditLogSaveFailedCode       = 208000002 // 保存审计日志失败
	ErrAuditLogQueryFailedCode      = 208000003 // 查询审计日志失败
	ErrAuditLogDeleteFailedCode     = 208000004 // 删除审计日志失败
	ErrAuditLogExportFailedCode     = 208000005 // 导出审计日志失败
	ErrDataMaskingFailedCode        = 208000006 // 数据脱敏失败
	ErrSignatureGenerationFailedCode = 208000007 // 生成签名失败
	ErrSignatureVerificationFailedCode = 208000008 // 验证签名失败

	// Audit相关 (208 010 000 ~ 208 019 999)
	ErrAuditInvalidRequestCode      = 208010001 // 审计请求无效
	ErrAuditPermissionDeniedCode    = 208010002 // 审计权限不足
	ErrAuditArchiveFailedCode       = 208010003 // 审计日志归档失败
	ErrAuditRestoreFailedCode       = 208010004 // 审计日志恢复失败

	// 导出相关 (208 020 000 ~ 208 029 999)
	ErrExportLimitExceededCode      = 208020001 // 导出数量超限
	ErrInvalidExportFormatCode      = 208020002 // 无效的导出格式

	// 统计相关 (208 030 000 ~ 208 039 999)
	ErrAuditStatisticsQueryFailedCode = 208030001 // 审计统计查询失败

	// 基础设施相关 (208 040 000 ~ 208 049 999)
	ErrElasticsearchNotAvailableCode = 208040001 // Elasticsearch不可用
)

// 便捷错误码别名（不带Err前缀，用于直接引用int32错误码）
const (
	DataMaskingFailed              = ErrDataMaskingFailedCode
	SignatureGenerationFailed      = ErrSignatureGenerationFailedCode
	AuditLogSaveFailed             = ErrAuditLogSaveFailedCode
	InvalidRequest                 = ErrAuditInvalidRequestCode
	AuditLogQueryFailed            = ErrAuditLogQueryFailedCode
	ExportLimitExceeded            = ErrExportLimitExceededCode
	InvalidExportFormat            = ErrInvalidExportFormatCode
	AuditStatisticsQueryFailed     = ErrAuditStatisticsQueryFailedCode
	ElasticsearchNotAvailable      = ErrElasticsearchNotAvailableCode
)

// 便捷错误变量
var (
	ErrAuditLogNotFound             = errorx.New(ErrAuditLogNotFoundCode)
	ErrAuditLogSaveFailed           = errorx.New(ErrAuditLogSaveFailedCode)
	ErrAuditLogQueryFailed          = errorx.New(ErrAuditLogQueryFailedCode)
	ErrDataMaskingFailed            = errorx.New(ErrDataMaskingFailedCode)
	ErrSignatureGenerationFailed    = errorx.New(ErrSignatureGenerationFailedCode)
	ErrSignatureVerificationFailed  = errorx.New(ErrSignatureVerificationFailedCode)
	ErrAuditInvalidRequest          = errorx.New(ErrAuditInvalidRequestCode)
	ErrAuditPermissionDenied        = errorx.New(ErrAuditPermissionDeniedCode)
)

// 通用错误码（跨模块使用）
const (
	ErrInvalidRequestCode = 400000001 // 无效请求
	ErrRecordNotFoundCode = 404000001 // 记录不存在
	ErrIDGenErrorCode     = 500000001 // ID生成失败
)

// 通用便捷错误变量
var (
	ErrInvalidRequest = errorx.New(ErrInvalidRequestCode)
	ErrRecordNotFound = errorx.New(ErrRecordNotFoundCode)
	ErrIDGenError     = errorx.New(ErrIDGenErrorCode)
)

func init() {
	// =====================================================
	// 208段：标准错误码注册（Audit模块）
	// =====================================================

	// Audit基础错误注册
	code.Register(
		ErrAuditLogNotFoundCode,
		"Audit log not found: {log_id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAuditLogSaveFailedCode,
		"Failed to save audit log: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAuditLogQueryFailedCode,
		"Failed to query audit logs: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAuditLogDeleteFailedCode,
		"Failed to delete audit log: {log_id}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAuditLogExportFailedCode,
		"Failed to export audit logs: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrDataMaskingFailedCode,
		"Failed to mask sensitive data: {field}, error: {error}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrSignatureGenerationFailedCode,
		"Failed to generate signature: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrSignatureVerificationFailedCode,
		"Failed to verify signature: {log_id}, error: {error}",
		code.WithAffectStability(false),
	)

	// Audit相关错误注册
	code.Register(
		ErrAuditInvalidRequestCode,
		"Invalid audit request: {param}, value: {value}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAuditPermissionDeniedCode,
		"Permission denied for audit operation: user_id={user_id}, action={action}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrAuditArchiveFailedCode,
		"Failed to archive audit logs: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrAuditRestoreFailedCode,
		"Failed to restore audit logs: {error}",
		code.WithAffectStability(true),
	)

	// 导出相关错误注册
	code.Register(
		ErrExportLimitExceededCode,
		"Export limit exceeded: requested {requested}, maximum {maximum}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrInvalidExportFormatCode,
		"Invalid export format: {format}, supported formats: {supported_formats}",
		code.WithAffectStability(false),
	)

	// 统计相关错误注册
	code.Register(
		ErrAuditStatisticsQueryFailedCode,
		"Failed to query audit statistics: {error}",
		code.WithAffectStability(true),
	)

	// 基础设施相关错误注册
	code.Register(
		ErrElasticsearchNotAvailableCode,
		"Elasticsearch is not available or not enabled",
		code.WithAffectStability(false),
	)

	// =====================================================
	// 通用错误码注册（跨模块）
	// =====================================================

	code.Register(
		ErrInvalidRequestCode,
		"Invalid request: {message}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrRecordNotFoundCode,
		"Record not found: {resource}, id: {id}",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrIDGenErrorCode,
		"Failed to generate ID: {error}",
		code.WithAffectStability(true),
	)
}
