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

// ReleasedWorkflow 已发布工作流（领域值对象）
// 从 api/model/workflow.ReleasedWorkflow 迁移过来
// TODO: 完整迁移所有字段
type ReleasedWorkflow struct {
	PluginID    string
	WorkflowID  string
	SpaceID     string
	Name        string
	Desc        string
	Icon        string
	Inputs      string
	Outputs     string
	EndType     int32
	Type        int32
	Version     string
	CreateTime  int64
	UpdateTime  int64
	FlowMode    WorkflowMode
	FlowVersion string
}

// WorkflowDetailData 工作流详情数据（领域值对象）
// 从 api/model/workflow.WorkflowDetailData 迁移过来
// TODO: 完整迁移所有字段
type WorkflowDetailData struct {
	WorkflowID string
	SpaceID    string
	Name       string
	Desc       string
	Icon       string
	Inputs     string
	Outputs    string
	Version    string
	CreateTime int64
	UpdateTime int64
	ProjectID  string
	EndType    int32
	IconURI    string
	FlowMode   WorkflowMode
}

// WorkflowDetailInfoData 工作流详情信息数据（领域值对象）
// 从 api/model/workflow.WorkflowDetailInfoData 迁移过来
// TODO: 完整迁移所有字段
type WorkflowDetailInfoData struct {
	WorkflowID            string
	SpaceID               string
	Name                  string
	Desc                  string
	Icon                  string
	Inputs                string
	Outputs               string
	Version               string
	CreateTime            int64
	UpdateTime            int64
	ProjectID             string
	EndType               int32
	IconURI               string
	FlowMode              WorkflowMode
	PluginID              string
	Creator               *Creator
	FlowVersion           string
	FlowVersionDesc       string
	LatestFlowVersion     string
	LatestFlowVersionDesc string
	CommitID              string
	IsProject             bool
}

// SubWorkflow 已在 canvas.go 中定义，此处不再重复
// TODO: 检查是否有字段差异需要合并

// Creator 创建者信息（领域值对象）
// 从 api/model/workflow.Creator 迁移过来
type Creator struct {
	ID       string
	Username string
	Nickname string
	Avatar   string
	Self     bool
}

// NodeInfo 节点信息（领域值对象）
// 从 api/model/workflow.NodeInfo 迁移过来
type NodeInfo struct {
	NodeID   string
	NodeName string
	NodeType string
}

// Validate 验证值对象
func (r *ReleasedWorkflow) Validate() error {
	if r.WorkflowID == "" {
		return ErrInvalidWorkflowID
	}
	return nil
}

// Validate 验证值对象
func (w *WorkflowDetailData) Validate() error {
	if w.WorkflowID == "" {
		return ErrInvalidWorkflowID
	}
	return nil
}
