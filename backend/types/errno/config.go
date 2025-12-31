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
// 统一错误码规范 - Config模块
// 错误码段: 207 000 000 ~ 207 999 999
// =====================================================

const (
	// 配置加载错误 (207 000 000 ~ 207 009 999)
	ErrConfigLoadFailedCode      = 207000001 // 配置加载失败
	ErrConfigParseFailedCode     = 207000002 // 配置解析失败
	ErrConfigNotFoundCode        = 207000003 // 配置不存在

	// 配置验证错误 (207 010 000 ~ 207 019 999)
	ErrConfigValidationFailedCode = 207010001 // 配置验证失败
	ErrConfigInvalidValueCode    = 207010002 // 配置值无效

	// 代码运行器配置错误 (207 020 000 ~ 207 029 999)
	ErrInvalidCodeRunnerTypeCode = 207020001 // 代码运行器类型无效
	ErrInvalidTimeoutSecondsCode = 207020002 // 超时秒数无效
)

// 便捷错误变量
var (
	ErrConfigLoadFailed       = errorx.New(ErrConfigLoadFailedCode)
	ErrConfigParseFailed      = errorx.New(ErrConfigParseFailedCode)
	ErrConfigNotFound         = errorx.New(ErrConfigNotFoundCode)
	ErrConfigValidationFailed = errorx.New(ErrConfigValidationFailedCode)
	ErrConfigInvalidValue     = errorx.New(ErrConfigInvalidValueCode)
	ErrInvalidCodeRunnerType  = errorx.New(ErrInvalidCodeRunnerTypeCode)
	ErrInvalidTimeoutSeconds  = errorx.New(ErrInvalidTimeoutSecondsCode)
)

func init() {
	// =====================================================
	// 207段：标准错误码注册（Config模块）
	// =====================================================

	// 配置加载错误注册
	code.Register(
		ErrConfigLoadFailedCode,
		"Failed to load configuration: config_key={config_key}, source={source}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrConfigParseFailedCode,
		"Failed to parse configuration: config_key={config_key}, format={format}, error: {error}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrConfigNotFoundCode,
		"Configuration not found: config_key={config_key}",
		code.WithAffectStability(false),
	)

	// 配置验证错误注册
	code.Register(
		ErrConfigValidationFailedCode,
		"Configuration validation failed: config_key={config_key}, validation_errors={validation_errors}",
		code.WithAffectStability(true),
	)

	code.Register(
		ErrConfigInvalidValueCode,
		"Invalid configuration value: config_key={config_key}, value={value}, reason: {reason}",
		code.WithAffectStability(false),
	)

	// 代码运行器配置错误注册
	code.Register(
		ErrInvalidCodeRunnerTypeCode,
		"Invalid code runner type: {type}, must be one of: local (0), sandbox (1)",
		code.WithAffectStability(false),
	)

	code.Register(
		ErrInvalidTimeoutSecondsCode,
		"Invalid timeout seconds: {seconds}, must be >= 0",
		code.WithAffectStability(false),
	)
}
