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

// FieldPermission 字段权限
type FieldPermission struct {
	PermissionID string       `json:"permission_id" gorm:"primaryKey;type:varchar(36)"`
	RoleID       string       `json:"role_id" gorm:"type:varchar(36);not null;index:idx_role_id"`
	ResourceType string       `json:"resource_type" gorm:"type:varchar(50);not null;index:idx_resource_type"`
	FieldName    string       `json:"field_name" gorm:"type:varchar(100);not null"`
	Permission   FieldPerm    `json:"permission" gorm:"type:enum('visible','editable','hidden');not null"`
	CreatedAt    int64        `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64        `json:"updated_at" gorm:"not null;default:0"`
}

// FieldPerm 字段权限级别
type FieldPerm string

const (
	FieldPermVisible  FieldPerm = "visible"  // 可见
	FieldPermEditable FieldPerm = "editable" // 可编辑
	FieldPermHidden   FieldPerm = "hidden"   // 隐藏
)

// TableName 指定表名
func (FieldPermission) TableName() string {
	return "field_permissions"
}

// DataPermissionLevel 数据权限级别
type DataPermissionLevel int

const (
	DataPermLevelAll          DataPermissionLevel = 5 // 全部数据
	DataPermLevelDeptAndBelow DataPermissionLevel = 4 // 本部门及下级部门
	DataPermLevelDept         DataPermissionLevel = 3 // 本部门
	DataPermLevelSelf         DataPermissionLevel = 2 // 本人数据
	DataPermLevelNone         DataPermissionLevel = 1 // 无数据权限
)

// OrganizationTemplate 组织模板
type OrganizationTemplate struct {
	TemplateID      string                `json:"template_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string                `json:"tenant_id" gorm:"type:varchar(36);index:idx_tenant_id"`
	TemplateName    string                `json:"template_name" gorm:"type:varchar(200);not null"`
	TemplateCode    string                `json:"template_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_code"`
	Category        string                `json:"category" gorm:"type:varchar(50);not null"` // 行业分类：tech, finance, retail等
	Description     string                `json:"description" gorm:"type:text"`
	Thumbnail       string                `json:"thumbnail" gorm:"type:varchar(500)"`
	IsPredefined    bool                  `json:"is_predefined" gorm:"default:false;index:idx_predefined"`
	Structure       string                `json:"structure" gorm:"type:longtext"` // JSON格式的组织结构
	Positions       string                `json:"positions" gorm:"type:longtext"`  // JSON格式的岗位配置
	Permissions     string                `json:"permissions" gorm:"type:longtext"` // JSON格式的权限配置
	Tags            string                `json:"tags" gorm:"type:varchar(500)"`
	UsageCount      int                   `json:"usage_count" gorm:"default:0;index:idx_usage"`
	Rating          float64               `json:"rating" gorm:"default:0"` // 模板评分
	Status          OrganizationStatus    `json:"status" gorm:"type:enum('active','inactive','frozen');default:'active'"`
	CreatedBy       string                `json:"created_by" gorm:"type:varchar(36)"`
	CreatedAt       int64                 `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt       int64                 `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt       *int64                `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (OrganizationTemplate) TableName() string {
	return "organization_templates"
}

// MigrationStatus 迁移状态
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"   // 待处理
	MigrationStatusRunning   MigrationStatus = "running"   // 进行中
	MigrationStatusCompleted MigrationStatus = "completed" // 已完成
	MigrationStatusFailed    MigrationStatus = "failed"    // 失败
	MigrationStatusRollback  MigrationStatus = "rollback"  // 回滚中
)

// OrganizationMigration 组织迁移
type OrganizationMigration struct {
	MigrationID     string             `json:"migration_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string             `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	SourceTenantID  string             `json:"source_tenant_id" gorm:"type:varchar(36);not null;index:idx_source"`
	TargetTenantID  string             `json:"target_tenant_id" gorm:"type:varchar(36);not null;index:idx_target"`
	MigrationType   string             `json:"migration_type" gorm:"type:varchar(50);not null"` // full, incremental
	SourceOrgID     string             `json:"source_org_id" gorm:"type:varchar(36)"`
	TargetOrgID     *string            `json:"target_org_id,omitempty" gorm:"type:varchar(36)"`
	TotalSteps      int                `json:"total_steps" gorm:"default:0"`
	CompletedSteps  int                `json:"completed_steps" gorm:"default:0"`
	Progress        int                `json:"progress" gorm:"default:0"` // 0-100
	Status          MigrationStatus    `json:"status" gorm:"type:enum('pending','running','completed','failed','rollback');default:'pending';index:idx_status"`
	ErrorMsg        string             `json:"error_msg,omitempty" gorm:"type:text"`
	StartedAt       *int64             `json:"started_at,omitempty"`
	CompletedAt     *int64             `json:"completed_at,omitempty"`
	CreatedBy       string             `json:"created_by" gorm:"type:varchar(36)"`
	CreatedAt       int64              `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt       int64              `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (OrganizationMigration) TableName() string {
	return "organization_migrations"
}

// OrganizationStats 组织统计
type OrganizationStats struct {
	OrgID              string  `json:"org_id"`
	TotalDepartments   int     `json:"total_departments"`
	TotalEmployees     int     `json:"total_employees"`
	ActiveEmployees    int     `json:"active_employees"`
	TotalPositions     int     `json:"total_positions"`
	ManagerCount       int     `json:"manager_count"`
	AvgJobLevel        float64 `json:"avg_job_level"`
	RecruitmentCount   int     `json:"recruitment_count"`
	ResignationCount   int     `json:"resignation_count"`
	TransferCount      int     `json:"transfer_count"`
	UpdatedAt          int64   `json:"updated_at"`
}

// DepartmentStats 部门统计
type DepartmentStats struct {
	DeptID          string  `json:"dept_id"`
	TotalEmployees  int     `json:"total_employees"`
	ActiveEmployees int     `json:"active_employees"`
	ManagerCount    int     `json:"manager_count"`
	AvgJobLevel     float64 `json:"avg_job_level"`
	UpdatedAt       int64   `json:"updated_at"`
}

// TenantStats 租户统计
type TenantStats struct {
	TenantID         string  `json:"tenant_id"`
	TotalOrgs        int     `json:"total_orgs"`
	TotalDepartments int     `json:"total_departments"`
	TotalEmployees   int     `json:"total_employees"`
	ActiveEmployees  int     `json:"active_employees"`
	TotalPositions   int     `json:"total_positions"`
	UpdatedAt        int64   `json:"updated_at"`
}

// ActivityStats 活跃度统计
type ActivityStats struct {
	OrgID          string                 `json:"org_id"`
	Days           int                    `json:"days"`
	DailyStats     map[string]int         `json:"daily_stats"`     // 每日活跃人数
	TopActiveUsers []string               `json:"top_active_users"` // 最活跃用户
	AvgActivity    float64                `json:"avg_activity"`
	PeakActivity   int                    `json:"peak_activity"`
	UpdatedAt      int64                  `json:"updated_at"`
}

// EmployeeSummary 员工摘要（用于通讯录、列表展示）
type EmployeeSummary struct {
	EmpID        string    `json:"emp_id"`
	Name         string    `json:"name"`
	EmpCode      string    `json:"emp_code"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	DeptName     string    `json:"dept_name"`
	OrgName      string    `json:"org_name"`
	PositionName string    `json:"position_name"`
	Status       EmployeeStatus `json:"status"`
}

// GetCreatedAtAsTime 获取创建时间
func (e *EmployeeSummary) GetCreatedAtAsTime() time.Time {
	return time.Unix(0, 0)
}
