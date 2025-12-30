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

package config

// CreateConfigReq 创建配置请求
type CreateConfigReq struct {
	TenantID    string `json:"tenant_id" binding:"required"`
	ConfigKey   string `json:"config_key" binding:"required"`
	ConfigValue string `json:"config_value" binding:"required"`
	ConfigType  string `json:"config_type" binding:"required,oneof=string int float bool json"`
	Description string `json:"description"`
}

// CreateConfigResp 创建配置响应
type CreateConfigResp struct {
	ConfigID string `json:"config_id"`
}

// GetConfigReq 获取配置请求
type GetConfigReq struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	ConfigKey string `json:"config_key" binding:"required"`
}

// GetConfigResp 获取配置响应
type GetConfigResp struct {
	ConfigID    string `json:"config_id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
}

// UpdateConfigReq 更新配置请求
type UpdateConfigReq struct {
	TenantID     string `json:"tenant_id" binding:"required"`
	ConfigKey    string `json:"config_key" binding:"required"`
	NewValue     string `json:"new_value" binding:"required"`
	ChangeReason string `json:"change_reason"`
}

// UpdateConfigResp 更新配置响应
type UpdateConfigResp struct {
	Success bool `json:"success"`
}

// DeleteConfigReq 删除配置请求
type DeleteConfigReq struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	ConfigKey string `json:"config_key" binding:"required"`
}

// DeleteConfigResp 删除配置响应
type DeleteConfigResp struct {
	Success bool `json:"success"`
}

// ListConfigsReq 列出配置请求
type ListConfigsReq struct {
	TenantID string `json:"tenant_id" binding:"required"`
}

// ListConfigsResp 列出配置响应
type ListConfigsResp struct {
	Configs []ConfigItem `json:"configs"`
	Total   int          `json:"total"`
}

// ConfigItem 配置项
type ConfigItem struct {
	ConfigID    string `json:"config_id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
}

// GetConfigHistoryReq 获取配置历史请求
type GetConfigHistoryReq struct {
	ConfigID string `json:"config_id" binding:"required"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
}

// GetConfigHistoryResp 获取配置历史响应
type GetConfigHistoryResp struct {
	History []ConfigHistoryItem `json:"history"`
	Total   int                 `json:"total"`
}

// ConfigHistoryItem 配置历史项
type ConfigHistoryItem struct {
	HistoryID     string `json:"history_id"`
	ConfigKey     string `json:"config_key"`
	OldValue      string `json:"old_value"`
	NewValue      string `json:"new_value"`
	ChangeReason  string `json:"change_reason"`
	ChangedBy     string `json:"changed_by"`
	ChangedAt     string `json:"changed_at"`
	VersionNumber int    `json:"version_number"`
}

// RollbackConfigReq 回滚配置请求
type RollbackConfigReq struct {
	TenantID string `json:"tenant_id" binding:"required"`
	ConfigKey string `json:"config_key" binding:"required"`
	Version  int    `json:"version" binding:"required"`
}

// RollbackConfigResp 回滚配置响应
type RollbackConfigResp struct {
	Success bool `json:"success"`
}
