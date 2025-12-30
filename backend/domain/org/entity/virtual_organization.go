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

// VirtualOrgType 虚拟组织类型
type VirtualOrgType string

const (
	VirtualOrgTypeProject    VirtualOrgType = "PROJECT"     // 项目组
	VirtualOrgTypeTaskForce  VirtualOrgType = "TASK_FORCE"  // 任务组
	VirtualOrgTypeCommunity  VirtualOrgType = "COMMUNITY"   // 社区
	VirtualOrgTypeOther      VirtualOrgType = "OTHER"       // 其他
)

// VirtualOrgStatus 虚拟组织状态
type VirtualOrgStatus string

const (
	VirtualOrgStatusActive   VirtualOrgStatus = "ACTIVE"    // 激活
	VirtualOrgStatusArchived VirtualOrgStatus = "ARCHIVED"  // 归档
)

// VirtualOrgMemberRole 虚拟组织成员角色
type VirtualOrgMemberRole string

const (
	VirtualOrgRoleOwner   VirtualOrgMemberRole = "OWNER"   // 所有者
	VirtualOrgRoleAdmin   VirtualOrgMemberRole = "ADMIN"   // 管理员
	VirtualOrgRoleMember  VirtualOrgMemberRole = "MEMBER"  // 成员
)

// VirtualOrganization 虚拟组织实体
type VirtualOrganization struct {
	ID          int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	VirtualOrgID string           `json:"virtual_org_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_virtual_org_id"`
	TenantID     string           `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	Name         string           `json:"name" gorm:"type:varchar(128);not null"`
	Description  *string          `json:"description,omitempty" gorm:"type:text"`
	Type         VirtualOrgType   `json:"type" gorm:"type:enum('PROJECT','TASK_FORCE','COMMUNITY','OTHER');not null;index:idx_type"`
	OwnerID      string           `json:"owner_id" gorm:"type:varchar(36);not null;index:idx_owner_id"`
	Status       VirtualOrgStatus `json:"status" gorm:"type:enum('ACTIVE','ARCHIVED');default:'ACTIVE';index:idx_status"`
	CreatedAt    int64            `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64            `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt    *int64           `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Members []VirtualOrgMember `json:"members,omitempty" gorm:"foreignKey:VirtualOrgID;references:VirtualOrgID"`
	Tags    []VirtualOrgTag    `json:"tags,omitempty" gorm:"foreignKey:VirtualOrgID;references:VirtualOrgID"`
}

// TableName 指定表名
func (VirtualOrganization) TableName() string {
	return "virtual_organizations"
}

// IsDeleted 是否已删除
func (vo *VirtualOrganization) IsDeleted() bool {
	return vo.DeletedAt != nil
}

// IsActive 是否激活
func (vo *VirtualOrganization) IsActive() bool {
	return vo.Status == VirtualOrgStatusActive
}

// GetCreatedAtAsTime 获取创建时间
func (vo *VirtualOrganization) GetCreatedAtAsTime() time.Time {
	return time.Unix(vo.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (vo *VirtualOrganization) GetUpdatedAtAsTime() time.Time {
	return time.Unix(vo.UpdatedAt/1000, 0)
}

// VirtualOrgMember 虚拟组织成员实体
type VirtualOrgMember struct {
	ID            int64                `json:"id" gorm:"primaryKey;autoIncrement"`
	VirtualOrgID  string              `json:"virtual_org_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_org_user"`
	UserID        string              `json:"user_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_org_user"`
	Role          VirtualOrgMemberRole `json:"role" gorm:"type:enum('OWNER','ADMIN','MEMBER');not null"`
	JoinedAt      int64               `json:"joined_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (VirtualOrgMember) TableName() string {
	return "virtual_org_members"
}

// GetJoinedAtAsTime 获取加入时间
func (vom *VirtualOrgMember) GetJoinedAtAsTime() time.Time {
	return time.Unix(vom.JoinedAt/1000, 0)
}

// IsOwner 是否是所有者
func (vom *VirtualOrgMember) IsOwner() bool {
	return vom.Role == VirtualOrgRoleOwner
}

// IsAdmin 是否是管理员
func (vom *VirtualOrgMember) IsAdmin() bool {
	return vom.Role == VirtualOrgRoleAdmin
}

// VirtualOrgTag 虚拟组织标签实体
type VirtualOrgTag struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	VirtualOrgID string `json:"virtual_org_id" gorm:"type:varchar(36);not null;uniqueIndex:uk_org_tag"`
	TagName      string `json:"tag_name" gorm:"type:varchar(64);not null;uniqueIndex:uk_org_tag;index:idx_tag_name"`
}

// TableName 指定表名
func (VirtualOrgTag) TableName() string {
	return "virtual_org_tags"
}

// Permission 权限实体（用于权限继承）
type Permission struct {
	ID          int64  `json:"id"`
	PermissionID string `json:"permission_id"`
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	Action       string `json:"action"`
	Effect       string `json:"effect"` // allow/deny
}
