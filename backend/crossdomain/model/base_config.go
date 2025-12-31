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

package model

import "github.com/coze-dev/coze-studio/backend/types/errno"

// CodeRunnerType 代码运行器类型
// 与 api/model/admin/config.CodeRunnerType 保持兼容（int32）
type CodeRunnerType int32

const (
	CodeRunnerType_Local  CodeRunnerType = 0
	CodeRunnerType_Sandbox CodeRunnerType = 1
)

// String 返回代码运行器类型字符串
func (t CodeRunnerType) String() string {
	switch t {
	case CodeRunnerType_Sandbox:
		return "sandbox"
	case CodeRunnerType_Local:
		return "local"
	}
	return "unknown"
}

// BasicConfiguration 基础系统配置
//
// 该模型封装了系统的全局配置信息，与 api/model/admin/config.BasicConfiguration
// 保持字段兼容性。
type BasicConfiguration struct {
	// 用户管理
	AdminEmails             string
	DisableUserRegistration bool
	AllowRegistrationEmail  string

	// 服务器配置
	ServerHost string

	// 代码运行器配置
	CodeRunnerType CodeRunnerType
	SandboxConfig  *SandboxConfig

	// 插件配置
	PluginConfiguration *PluginConfiguration
}

// SandboxConfig 沙箱配置
// 与 api/model/admin/config.SandboxConfig 保持字段兼容
type SandboxConfig struct {
	AllowEnv       string
	AllowRead      string
	AllowWrite     string
	AllowNet       string
	AllowRun       string
	AllowFfi       string
	NodeModulesDir string
	TimeoutSeconds float64
	MemoryLimitMb  int64
}

// PluginConfiguration 插件配置
// 与 api/model/admin/config.PluginConfiguration 保持字段兼容
type PluginConfiguration struct {
	CozeSaasPluginEnabled bool
	CozeAPIToken          string
	CozeSaasAPIBaseURL    string
}

// Validate 验证基础配置
func (c *BasicConfiguration) Validate() error {
	if c.CodeRunnerType != CodeRunnerType_Local && c.CodeRunnerType != CodeRunnerType_Sandbox {
		return errno.ErrInvalidCodeRunnerType
	}
	if c.SandboxConfig != nil && c.SandboxConfig.TimeoutSeconds < 0 {
		return errno.ErrInvalidTimeoutSeconds
	}
	return nil
}

// GetServerHost 获取服务器主机地址，如果为空返回默认值
func (c *BasicConfiguration) GetServerHost(defaultHost string) string {
	if c == nil || c.ServerHost == "" {
		return defaultHost
	}
	return c.ServerHost
}
