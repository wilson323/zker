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

package internal

const (
	DefaultErrorMsg          = "Service Internal Error"
	DefaultIsAffectStability = true
)

var (
	ServiceInternalErrorCode int32 = 1
	CodeDefinitions                = make(map[int32]*CodeDefinition)
)

type CodeDefinition struct {
	Code              int32
	Message           string // 英文错误消息
	MessageZH         string // 中文错误消息（可选）
	IsAffectStability bool
}

type RegisterOption func(definition *CodeDefinition)

func WithAffectStability(affectStability bool) RegisterOption {
	return func(definition *CodeDefinition) {
		definition.IsAffectStability = affectStability
	}
}

// WithZhMessage 设置中文错误消息（可选）
func WithZhMessage(messageZH string) RegisterOption {
	return func(definition *CodeDefinition) {
		definition.MessageZH = messageZH
	}
}

func Register(code int32, msg string, opts ...RegisterOption) {
	definition := &CodeDefinition{
		Code:              code,
		Message:           msg,
		IsAffectStability: DefaultIsAffectStability,
	}

	for _, opt := range opts {
		opt(definition)
	}

	CodeDefinitions[code] = definition
}

func SetDefaultErrorCode(code int32) {
	ServiceInternalErrorCode = code
}
