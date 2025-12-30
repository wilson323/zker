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

import "time"

// OrganizationType 组织类型
type OrganizationType string

const (
	OrgTypeCompany   OrganizationType = "company"   // 公司
	OrgTypeDivision  OrganizationType = "division"  // 分公司
	OrgTypeDepartment OrganizationType = "department" // 部门（虚拟组织）
	OrgTypeProject   OrganizationType = "project"   // 项目组
)

// OrganizationStatus 组织状态
type OrganizationStatus string

const (
	OrgStatusActive   OrganizationStatus = "active"   // 激活
	OrgStatusInactive OrganizationStatus = "inactive" // 停用
	OrgStatusFrozen   OrganizationStatus = "frozen"   // 冻结
)

// Organization 组织实体
type Organization struct {
	OrgID        string              `json:"org_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID     string              `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	OrgName      string              `json:"org_name" gorm:"type:varchar(200);not null"`
	OrgType      OrganizationType    `json:"org_type" gorm:"type:enum('company','division','department','project');not null"`
	ParentID     *string             `json:"parent_id,omitempty" gorm:"type:varchar(36)"`
	OrgCode      string              `json:"org_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_tenant_code"`
	Level        int                 `json:"level" gorm:"type:int;not null;default:1;index:idx_level"` // 层级：1-顶级，2-二级...
	Path         string              `json:"path" gorm:"type:varchar(500);not null"`                   // 组织路径：/1/2/3
	SortOrder    int                 `json:"sort_order" gorm:"type:int;default:0;index:idx_sort"`    // 排序
	Status       OrganizationStatus  `json:"status" gorm:"type:enum('active','inactive','frozen');default:'active';index:idx_status"`
	Description string              `json:"description" gorm:"type:text"`
	LeaderID    *string             `json:"leader_id,omitempty" gorm:"type:varchar(36)"` // 负责人ID
	CreatedAt   int64               `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt   int64               `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt   *int64              `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Parent   *Organization  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Organization `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Leader   *Employee      `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
}

// TableName 指定表名
func (Organization) TableName() string {
	return "organizations"
}

// IsActive 是否激活
func (o *Organization) IsActive() bool {
	return o.Status == OrgStatusActive
}

// IsDeleted 是否已删除
func (o *Organization) IsDeleted() bool {
	return o.DeletedAt != nil
}

// HasParent 是否有父组织
func (o *Organization) HasParent() bool {
	return o.ParentID != nil && *o.ParentID != ""
}

// GetLevel 获取层级
func (o *Organization) GetLevel() int {
	return o.Level
}

// GetCreatedAtAsTime 获取创建时间
func (o *Organization) GetCreatedAtAsTime() time.Time {
	return time.Unix(o.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (o *Organization) GetUpdatedAtAsTime() time.Time {
	return time.Unix(o.UpdatedAt/1000, 0)
}

// GetAncestorCount 获取祖先数量
func (o *Organization) GetAncestorCount() int {
	if o.Path == "/" {
		return 0
	}
	count := 0
	for _, c := range o.Path {
		if c == '/' {
			count++
		}
	}
	return count - 1
}

// Department 部门实体
type Department struct {
	DeptID       string             `json:"dept_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID     string             `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	OrgID        string             `json:"org_id" gorm:"type:varchar(36);not null;index:idx_org_id"` // 所属组织ID
	ParentID     *string            `json:"parent_id,omitempty" gorm:"type:varchar(36)"`              // 父部门ID
	DeptName     string             `json:"dept_name" gorm:"type:varchar(200);not null"`
	DeptCode     string             `json:"dept_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_tenant_code"`
	Level        int                `json:"level" gorm:"type:int;not null;default:1;index:idx_level"`
	Path         string             `json:"path" gorm:"type:varchar(500);not null"`
	SortOrder    int                `json:"sort_order" gorm:"type:int;default:0;index:idx_sort"`
	Status       OrganizationStatus `json:"status" gorm:"type:enum('active','inactive','frozen');default:'active';index:idx_status"`
	LeaderID     *string            `json:"leader_id,omitempty" gorm:"type:varchar(36)"`
	ParentLeader *string            `json:"parent_leader,omitempty" gorm:"type:varchar(36)"` // 上级部门负责人
	Description  string             `json:"description" gorm:"type:text"`
	CreatedAt    int64              `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64              `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt    *int64             `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrgID"`
	Parent       *Department   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children     []Department  `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Leader       *Employee     `json:"leader,omitempty" gorm:"foreignKey:LeaderID"`
}

// TableName 指定表名
func (Department) TableName() string {
	return "departments"
}

// IsActive 是否激活
func (d *Department) IsActive() bool {
	return d.Status == OrgStatusActive
}

// IsDeleted 是否已删除
func (d *Department) IsDeleted() bool {
	return d.DeletedAt != nil
}

// HasParent 是否有父部门
func (d *Department) HasParent() bool {
	return d.ParentID != nil && *d.ParentID != ""
}

// GetCreatedAtAsTime 获取创建时间
func (d *Department) GetCreatedAtAsTime() time.Time {
	return time.Unix(d.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (d *Department) GetUpdatedAtAsTime() time.Time {
	return time.Unix(d.UpdatedAt/1000, 0)
}

// Position 岗位实体
type Position struct {
	PositionID  string             `json:"position_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID    string             `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	DeptID      *string            `json:"dept_id,omitempty" gorm:"type:varchar(36);index:idx_dept_id"` // 所属部门
	PositionName string            `json:"position_name" gorm:"type:varchar(100);not null"`
	PositionCode string            `json:"position_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_tenant_code"`
	Level        int               `json:"level" gorm:"type:int;not null;default:1;index:idx_level"` // 职级：1-初级，2-中级，3-高级
	Category    string             `json:"category" gorm:"type:varchar(50);not null"` // 技术岗/管理岗/职能岗
	Responsibilities string        `json:"responsibilities" gorm:"type:text"` // 岗位职责
	Requirements   string         `json:"requirements" gorm:"type:text"`     // 任职要求
	Status       OrganizationStatus `json:"status" gorm:"type:enum('active','inactive','frozen');default:'active';index:idx_status"`
	SortOrder    int               `json:"sort_order" gorm:"type:int;default:0;index:idx_sort"`
	CreatedAt    int64             `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64             `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt    *int64            `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Department *Department `json:"department,omitempty" gorm:"foreignKey:DeptID"`
}

// TableName 指定表名
func (Position) TableName() string {
	return "positions"
}

// IsActive 是否激活
func (p *Position) IsActive() bool {
	return p.Status == OrgStatusActive
}

// IsDeleted 是否已删除
func (p *Position) IsDeleted() bool {
	return p.DeletedAt != nil
}

// GetCreatedAtAsTime 获取创建时间
func (p *Position) GetCreatedAtAsTime() time.Time {
	return time.Unix(p.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (p *Position) GetUpdatedAtAsTime() time.Time {
	return time.Unix(p.UpdatedAt/1000, 0)
}
