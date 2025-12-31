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

import "time"

// WorkflowConfig 跨领域的工作流配置模型
type WorkflowConfig struct {
	// 基础信息
	WorkflowID   string
	WorkflowName string
	WorkflowDesc string
	IsPublic     bool
	IsActive     bool

	// 模型配置
	ModelID   *int64
	ModelName string

	// LLM参数
	Temperature float32
	MaxTokens   int32

	// 工作流特定配置
	TimeoutSeconds int32

	// 时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// WorkflowNode 工作流节点
type WorkflowNode struct {
	NodeID      string
	NodeType    string
	WorkflowID  string
	IsEnabled   bool
	ConfigData  []byte // JSON编码的节点配置
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Validate 验证工作流配置
func (c *WorkflowConfig) Validate() error {
	if c.WorkflowID == "" {
		return ErrInvalidWorkflowID
	}
	if c.ModelID == nil && c.ModelName == "" {
		return ErrInvalidModelConfig
	}
	if c.Temperature < 0 || c.Temperature > 2 {
		return ErrInvalidTemperature
	}
	return nil
}

// 错误定义
var (
	ErrInvalidWorkflowID = &ModelError{Message: "workflow_id cannot be empty"}
)
