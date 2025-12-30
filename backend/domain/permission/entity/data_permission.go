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

// ResourceType 资源类型
type ResourceType string

const (
	ResourceTypeBots         ResourceType = "bots"         // Bot
	ResourceTypeConversations ResourceType = "conversations" // 对话
	ResourceTypeKnowledge    ResourceType = "knowledge"    // 知识库
	ResourceTypeWorkflows    ResourceType = "workflows"    // 工作流
	ResourceTypePlugins      ResourceType = "plugins"      // 插件
)

// DataPermissionScope 数据权限范围
type DataPermissionScope string

const (
	DataPermissionScopeAll                 DataPermissionScope = "ALL"                  // 全部数据
	DataPermissionScopeDepartment          DataPermissionScope = "DEPARTMENT"           // 本部门数据
	DataPermissionScopeDepartmentAndSub    DataPermissionScope = "DEPARTMENT_AND_SUB"   // 本部门及子部门数据
	DataPermissionScopeOwn                 DataPermissionScope = "OWN"                  // 仅自己的数据
	DataPermissionScopeCustom              DataPermissionScope = "CUSTOM"               // 自定义过滤
	DataPermissionScopeNone                DataPermissionScope = "NONE"                 // 无权限
)

// DataPermission 数据权限实体
type DataPermission struct {
	PermissionID  string              `json:"permission_id" gorm:"primaryKey;type:varchar(36)"`
	RoleID        string              `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
	ResourceType  ResourceType        `json:"resource_type" gorm:"type:enum('bots','conversations','knowledge','workflows','plugins');not null;index:idx_resource_type"`
	Scope         DataPermissionScope `json:"scope" gorm:"type:enum('ALL','DEPARTMENT','OWN','CUSTOM','NONE');not null"`
	CustomFilter  string              `json:"custom_filter,omitempty" gorm:"type:json"`
	CreatedAt     int64               `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt     int64               `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (DataPermission) TableName() string {
	return "data_permissions"
}

// HasPermission 是否有权限（不是NONE）
func (dp *DataPermission) HasPermission() bool {
	return dp.Scope != DataPermissionScopeNone
}

// IsAllPermission 是否是全部数据权限
func (dp *DataPermission) IsAllPermission() bool {
	return dp.Scope == DataPermissionScopeAll
}

// IsOwnPermission 是否是仅自己数据权限
func (dp *DataPermission) IsOwnPermission() bool {
	return dp.Scope == DataPermissionScopeOwn
}

// IsDepartmentPermission 是否是部门数据权限
func (dp *DataPermission) IsDepartmentPermission() bool {
	return dp.Scope == DataPermissionScopeDepartment
}

// IsCustomPermission 是否是自定义过滤权限
func (dp *DataPermission) IsCustomPermission() bool {
	return dp.Scope == DataPermissionScopeCustom
}

// Department 部门实体
type Department struct {
	DepartmentID       string  `json:"department_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID           string  `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	DepartmentName     string  `json:"department_name" gorm:"type:varchar(200);not null"`
	ParentDepartmentID *string `json:"parent_department_id,omitempty" gorm:"type:varchar(36)"`
	CreatedAt          int64   `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt          int64   `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt          *int64  `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	ParentDepartment *Department `json:"parent_department,omitempty" gorm:"foreignKey:ParentDepartmentID"`
	UserDepartments  []UserDepartment `json:"user_departments,omitempty" gorm:"foreignKey:DepartmentID"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}

// HasParent 是否有父部门
func (d *Department) HasParent() bool {
	return d.ParentDepartmentID != nil && *d.ParentDepartmentID != ""
}

// UserDepartment 用户部门关联实体
type UserDepartment struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID       string `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID     string `json:"tenant_id" gorm:"type:varchar(36);not null"`
	DepartmentID string `json:"department_id" gorm:"type:varchar(36);not null;index:idx_department_id"`
	IsLeader     bool   `json:"is_leader" gorm:"default:false"`
	CreatedAt    int64  `json:"created_at" gorm:"not null;default:0"`

	// 关联
	Department Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}

// TableName 指定表名
func (UserDepartment) TableName() string {
	return "user_departments"
}
