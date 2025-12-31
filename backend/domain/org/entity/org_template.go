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

package entity

// OrgTemplate 组织模板
type OrgTemplate struct {
	TemplateID   string              `json:"template_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID     string              `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	Name         string              `json:"name" gorm:"type:varchar(200);not null"`
	Description  string              `json:"description" gorm:"type:text"`
	Config       *TemplateConfig     `json:"config" gorm:"type:json;not null"`
	IsBuiltin    bool                `json:"is_builtin" gorm:"type:bool;default:false"`
	Version      string              `json:"version" gorm:"type:varchar(20);not null"`
	IsEnabled    bool                `json:"is_enabled" gorm:"type:bool;default:true"`
	CreatedBy    string              `json:"created_by" gorm:"type:varchar(36);not null"`
	UpdatedBy    string              `json:"updated_by" gorm:"type:varchar(36);not null"`
	CreatedAt    int64               `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64               `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt    *int64              `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (OrgTemplate) TableName() string {
	return "org_templates"
}

// IsDeleted 是否已删除
func (t *OrgTemplate) IsDeleted() bool {
	return t.DeletedAt != nil
}

// IsEnabled 是否启用
func (t *OrgTemplate) IsEnabledTemplate() bool {
	return t.IsEnabled
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	Departments  []*DepartmentConfig  `json:"departments,omitempty"`
	Roles        []*RoleConfig        `json:"roles,omitempty"`
	Permissions  []*PermissionConfig  `json:"permissions,omitempty"`
	Workflows    []*WorkflowConfig    `json:"workflows,omitempty"`
}

// DepartmentConfig 部门配置
type DepartmentConfig struct {
	DeptName    string  `json:"dept_name" binding:"required"`
	DeptCode    string  `json:"dept_code" binding:"required"`
	ParentID    *string `json:"parent_id,omitempty"`
	Description string  `json:"description,omitempty"`
	LeaderID    *string `json:"leader_id,omitempty"`
	SortOrder   int     `json:"sort_order,omitempty"`
}

// RoleConfig 角色配置
type RoleConfig struct {
	RoleName        string            `json:"role_name" binding:"required"`
	RoleCode        string            `json:"role_code" binding:"required"`
	Description     string            `json:"description,omitempty"`
	IsSystemRole    bool              `json:"is_system_role,omitempty"`
	Permissions     []*string         `json:"permissions,omitempty"`
	DataPermissions []*DataPermConfig `json:"data_permissions,omitempty"`
}

// DataPermConfig 数据权限配置
type DataPermConfig struct {
	ResourceType string `json:"resource_type" binding:"required"`
	Scope        string `json:"scope" binding:"required,oneof=all self department org custom"`
	CustomScope  *string `json:"custom_scope,omitempty"`
}

// PermissionConfig 权限配置
type PermissionConfig struct {
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required,oneof=create read update delete manage"`
	Description string `json:"description,omitempty"`
}

// WorkflowConfig 工作流配置
type WorkflowConfig struct {
	WorkflowName string           `json:"workflow_name" binding:"required"`
	WorkflowCode string           `json:"workflow_code" binding:"required"`
	Description  string           `json:"description,omitempty"`
	Steps        []*WorkflowStep  `json:"steps,omitempty"`
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
	StepName    string `json:"step_name" binding:"required"`
	StepType    string `json:"step_type" binding:"required"`
	Assignee    string `json:"assignee,omitempty"`
	Description string `json:"description,omitempty"`
}
