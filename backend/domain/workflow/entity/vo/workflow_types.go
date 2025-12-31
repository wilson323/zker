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

package vo

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	workflowapi "github.com/coze-dev/coze-studio/backend/api/model/workflow"
)

// WorkFlowType 工作流类型（领域值对象）
// 从 api/model/workflow.WorkFlowType 迁移过来
type WorkFlowType int64

const (
	WorkFlowType_User     WorkFlowType = 0 // 用户定义
	WorkFlowType_GuanFang WorkFlowType = 1 // 官方模板
)

// String 返回类型字符串
func (p WorkFlowType) String() string {
	switch p {
	case WorkFlowType_User:
		return "User"
	case WorkFlowType_GuanFang:
		return "GuanFang"
	}
	return "<UNSET>"
}

// WorkFlowTypeFromString 从字符串解析
func WorkFlowTypeFromString(s string) (WorkFlowType, error) {
	switch s {
	case "User":
		return WorkFlowType_User, nil
	case "GuanFang":
		return WorkFlowType_GuanFang, nil
	}
	return WorkFlowType(0), fmt.Errorf("not a valid WorkFlowType string")
}

// Scan implements sql.Scanner
func (p *WorkFlowType) Scan(value interface{}) (err error) {
	var result sql.NullInt64
	err = result.Scan(value)
	*p = WorkFlowType(result.Int64)
	return
}

// Value implements driver.Valuer
func (p *WorkFlowType) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Validate 验证值对象
func (t WorkFlowType) Validate() error {
	if t != WorkFlowType_User && t != WorkFlowType_GuanFang {
		return ErrInvalidWorkflowType
	}
	return nil
}

// Tag 工作流标签（领域值对象）
// 从 api/model/workflow.Tag 迁移过来
type Tag int64

const (
	Tag_All           Tag = 1
	Tag_Hot           Tag = 2
	Tag_Information   Tag = 3
	Tag_Music         Tag = 4
	Tag_Picture       Tag = 5
	Tag_UtilityTool   Tag = 6
	Tag_Life          Tag = 7
	Tag_Traval        Tag = 8
	Tag_Network       Tag = 9
	Tag_System        Tag = 10
	Tag_Movie         Tag = 11
	Tag_Office        Tag = 12
	Tag_Shopping      Tag = 13
	Tag_Education     Tag = 14
	Tag_Health        Tag = 15
	Tag_Social        Tag = 16
	Tag_Entertainment Tag = 17
	Tag_Finance       Tag = 18
	Tag_Hidden        Tag = 100
)

// String 返回标签字符串
func (p Tag) String() string {
	switch p {
	case Tag_All:
		return "All"
	case Tag_Hot:
		return "Hot"
	case Tag_Information:
		return "Information"
	case Tag_Music:
		return "Music"
	case Tag_Picture:
		return "Picture"
	case Tag_UtilityTool:
		return "UtilityTool"
	case Tag_Life:
		return "Life"
	case Tag_Traval:
		return "Traval"
	case Tag_Network:
		return "Network"
	case Tag_System:
		return "System"
	case Tag_Movie:
		return "Movie"
	case Tag_Office:
		return "Office"
	case Tag_Shopping:
		return "Shopping"
	case Tag_Education:
		return "Education"
	case Tag_Health:
		return "Health"
	case Tag_Social:
		return "Social"
	case Tag_Entertainment:
		return "Entertainment"
	case Tag_Finance:
		return "Finance"
	case Tag_Hidden:
		return "Hidden"
	}
	return "<UNSET>"
}

// Scan implements sql.Scanner
func (p *Tag) Scan(value interface{}) (err error) {
	var result sql.NullInt64
	err = result.Scan(value)
	*p = Tag(result.Int64)
	return
}

// Value implements driver.Valuer
func (p *Tag) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Validate 验证值对象
func (t Tag) Validate() error {
	if t < Tag_All || t > Tag_Hidden {
		// 跳过中间的值，只检查边界
		if t != Tag_Hidden {
			return ErrInvalidTag
		}
	}
	return nil
}

// WorkflowMode 工作流模式（领域值对象）
// 从 api/model/workflow.WorkflowMode 迁移过来
type WorkflowMode int64

const (
	WorkflowMode_Workflow  WorkflowMode = 0
	WorkflowMode_Imageflow WorkflowMode = 1
	WorkflowMode_SceneFlow WorkflowMode = 2
	WorkflowMode_ChatFlow  WorkflowMode = 3
	WorkflowMode_All       WorkflowMode = 100
)

// String 返回模式字符串
func (p WorkflowMode) String() string {
	switch p {
	case WorkflowMode_Workflow:
		return "Workflow"
	case WorkflowMode_Imageflow:
		return "Imageflow"
	case WorkflowMode_SceneFlow:
		return "SceneFlow"
	case WorkflowMode_ChatFlow:
		return "ChatFlow"
	case WorkflowMode_All:
		return "All"
	}
	return "<UNSET>"
}

// WorkflowModeFromString 从字符串解析
func WorkflowModeFromString(s string) (WorkflowMode, error) {
	switch s {
	case "Workflow":
		return WorkflowMode_Workflow, nil
	case "Imageflow":
		return WorkflowMode_Imageflow, nil
	case "SceneFlow":
		return WorkflowMode_SceneFlow, nil
	case "ChatFlow":
		return WorkflowMode_ChatFlow, nil
	case "All":
		return WorkflowMode_All, nil
	}
	return WorkflowMode(0), fmt.Errorf("not a valid WorkflowMode string")
}

// Scan implements sql.Scanner
func (p *WorkflowMode) Scan(value interface{}) (err error) {
	var result sql.NullInt64
	err = result.Scan(value)
	*p = WorkflowMode(result.Int64)
	return
}

// Value implements driver.Valuer
func (p *WorkflowMode) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return int64(*p), nil
}

// Validate 验证值对象
func (m WorkflowMode) Validate() error {
	validModes := []WorkflowMode{
		WorkflowMode_Workflow,
		WorkflowMode_Imageflow,
		WorkflowMode_SceneFlow,
		WorkflowMode_ChatFlow,
		WorkflowMode_All,
	}
	for _, mode := range validModes {
		if m == mode {
			return nil
		}
	}
	return ErrInvalidWorkflowMode
}

// APIParameter API参数（领域值对象）
// 使用类型别名直接引用 api/model/workflow.APIParameter
// 这样可以避免重复定义，保持类型兼容性
// 在 DDD 架构中，对于简单的数据传输对象，可以直接使用 API 层的类型
type APIParameter = workflowapi.APIParameter
