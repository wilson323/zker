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

// Package adapter 提供类型转换适配器
//
// 该包用于在 api/model（Thrift DTO）和 crossdomain/model（领域模型）之间转换，
// 使 bizpkg 层可以通过 crossdomain/model 接口访问数据，而底层存储仍使用 api/model 类型
// 以保持数据库兼容性。
package adapter

import (
	apiconfig "github.com/coze-dev/coze-studio/backend/api/model/admin/config"
	"github.com/coze-dev/coze-studio/backend/crossdomain/model"
)

// BasicConfigToDomain 将 api/model 的 BasicConfiguration 转换为 crossdomain/model 的 BasicConfiguration
func BasicConfigToDomain(apiConfig *apiconfig.BasicConfiguration) *model.BasicConfiguration {
	if apiConfig == nil {
		return nil
	}

	domainConfig := &model.BasicConfiguration{
		AdminEmails:             apiConfig.AdminEmails,
		DisableUserRegistration: apiConfig.DisableUserRegistration,
		AllowRegistrationEmail:  apiConfig.AllowRegistrationEmail,
		ServerHost:              apiConfig.ServerHost,
		CodeRunnerType:          model.CodeRunnerType(apiConfig.CodeRunnerType),
		PluginConfiguration:     PluginConfigToDomain(apiConfig.PluginConfiguration),
	}

	if apiConfig.SandboxConfig != nil {
		domainConfig.SandboxConfig = SandboxConfigToDomain(apiConfig.SandboxConfig)
	}

	return domainConfig
}

// BasicConfigFromDomain 将 crossdomain/model 的 BasicConfiguration 转换为 api/model 的 BasicConfiguration
func BasicConfigFromDomain(domainConfig *model.BasicConfiguration) *apiconfig.BasicConfiguration {
	if domainConfig == nil {
		return nil
	}

	apiCfg := &apiconfig.BasicConfiguration{
		AdminEmails:             domainConfig.AdminEmails,
		DisableUserRegistration: domainConfig.DisableUserRegistration,
		AllowRegistrationEmail:  domainConfig.AllowRegistrationEmail,
		ServerHost:              domainConfig.ServerHost,
		CodeRunnerType:          apiconfig.CodeRunnerType(domainConfig.CodeRunnerType),
		PluginConfiguration:     PluginConfigFromDomain(domainConfig.PluginConfiguration),
	}

	if domainConfig.SandboxConfig != nil {
		apiCfg.SandboxConfig = SandboxConfigFromDomain(domainConfig.SandboxConfig)
	}

	return apiCfg
}

// PluginConfigToDomain 将 api/model 的 PluginConfiguration 转换为 crossdomain/model 的 PluginConfiguration
func PluginConfigToDomain(apiConfig *apiconfig.PluginConfiguration) *model.PluginConfiguration {
	if apiConfig == nil {
		return nil
	}
	return &model.PluginConfiguration{
		CozeSaasPluginEnabled: apiConfig.CozeSaasPluginEnabled,
		CozeAPIToken:          apiConfig.CozeAPIToken,
		CozeSaasAPIBaseURL:    apiConfig.CozeSaasAPIBaseURL,
	}
}

// PluginConfigFromDomain 将 crossdomain/model 的 PluginConfiguration 转换为 api/model 的 PluginConfiguration
func PluginConfigFromDomain(domainConfig *model.PluginConfiguration) *apiconfig.PluginConfiguration {
	if domainConfig == nil {
		return nil
	}
	return &apiconfig.PluginConfiguration{
		CozeSaasPluginEnabled: domainConfig.CozeSaasPluginEnabled,
		CozeAPIToken:          domainConfig.CozeAPIToken,
		CozeSaasAPIBaseURL:    domainConfig.CozeSaasAPIBaseURL,
	}
}

// SandboxConfigToDomain 将 api/model 的 SandboxConfig 转换为 crossdomain/model 的 SandboxConfig
func SandboxConfigToDomain(apiConfig *apiconfig.SandboxConfig) *model.SandboxConfig {
	if apiConfig == nil {
		return nil
	}
	return &model.SandboxConfig{
		AllowEnv:       apiConfig.AllowEnv,
		AllowRead:      apiConfig.AllowRead,
		AllowWrite:     apiConfig.AllowWrite,
		AllowNet:       apiConfig.AllowNet,
		AllowRun:       apiConfig.AllowRun,
		AllowFfi:       apiConfig.AllowFfi,
		NodeModulesDir: apiConfig.NodeModulesDir,
		TimeoutSeconds: apiConfig.TimeoutSeconds,
		MemoryLimitMb:  apiConfig.MemoryLimitMb,
	}
}

// SandboxConfigFromDomain 将 crossdomain/model 的 SandboxConfig 转换为 api/model 的 SandboxConfig
func SandboxConfigFromDomain(domainConfig *model.SandboxConfig) *apiconfig.SandboxConfig {
	if domainConfig == nil {
		return nil
	}
	return &apiconfig.SandboxConfig{
		AllowEnv:       domainConfig.AllowEnv,
		AllowRead:      domainConfig.AllowRead,
		AllowWrite:     domainConfig.AllowWrite,
		AllowNet:       domainConfig.AllowNet,
		AllowRun:       domainConfig.AllowRun,
		AllowFfi:       domainConfig.AllowFfi,
		NodeModulesDir: domainConfig.NodeModulesDir,
		TimeoutSeconds: domainConfig.TimeoutSeconds,
		MemoryLimitMb:  domainConfig.MemoryLimitMb,
	}
}
