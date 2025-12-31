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

// MigrationTaskStatus 迁移任务状态
type MigrationTaskStatus string

const (
	MigrationTaskStatusPending   MigrationTaskStatus = "pending"     // 待执行
	MigrationTaskStatusRunning   MigrationTaskStatus = "running"     // 执行中
	MigrationTaskStatusCompleted MigrationTaskStatus = "completed"  // 已完成
	MigrationTaskStatusFailed    MigrationTaskStatus = "failed"      // 失败
	MigrationTaskStatusRolledBack MigrationTaskStatus = "rolled_back" // 已回滚
)

// MigrationType 迁移类型
type MigrationType string

const (
	MigrationTypeFull         MigrationType = "full"         // 全量迁移
	MigrationTypeIncremental  MigrationType = "incremental"  // 增量迁移
)

// MigrationDataType 迁移数据类型
type MigrationDataType string

const (
	MigrationDataEmployees    MigrationDataType = "employees"    // 员工数据
	MigrationDataDepartments  MigrationDataType = "departments"  // 部门数据
	MigrationDataRoles        MigrationDataType = "roles"        // 角色数据
	MigrationDataPermissions  MigrationDataType = "permissions"  // 权限数据
)

// OrgMigrationTask 组织迁移任务
type OrgMigrationTask struct {
	TaskID          string              `json:"task_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID        string              `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	SourceOrgID     string              `json:"source_org_id" gorm:"type:varchar(36);not null;index:idx_source_org"`
	TargetTenantID  string              `json:"target_tenant_id" gorm:"type:varchar(36);not null;index:idx_target_tenant"`
	TargetOrgID     *string             `json:"target_org_id,omitempty" gorm:"type:varchar(36);index:idx_target_org"`

	Type            MigrationType       `json:"type" gorm:"type:enum('full','incremental');not null"`
	DataTypes       []MigrationDataType `json:"data_types" gorm:"type:json;not null"`
	Status          MigrationTaskStatus `json:"status" gorm:"type:enum('pending','running','completed','failed','rolled_back');default:'pending';index:idx_status"`

	Progress        int                 `json:"progress" gorm:"type:int;default:0"`        // 进度 0-100
	TotalItems      int                 `json:"total_items" gorm:"type:int;default:0"`
	MigratedItems   int                 `json:"migrated_items" gorm:"type:int;default:0"`
	FailedItems     int                 `json:"failed_items" gorm:"type:int;default:0"`

	ErrorMessages   []string            `json:"error_messages,omitempty" gorm:"type:json"`

	StartedAt       *int64              `json:"started_at,omitempty" gorm:"type:bigint"`
	CompletedAt     *int64              `json:"completed_at,omitempty" gorm:"type:bigint"`
	RolledBackAt    *int64              `json:"rolled_back_at,omitempty" gorm:"type:bigint"`

	CreatedBy       string              `json:"created_by" gorm:"type:varchar(36);not null"`
	CreatedAt       int64               `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt       int64               `json:"updated_at" gorm:"not null;default:0"`

	// 关联
	SourceOrg       *Organization       `json:"source_org,omitempty" gorm:"foreignKey:SourceOrgID"`
}

// TableName 指定表名
func (OrgMigrationTask) TableName() string {
	return "org_migration_tasks"
}

// IsPending 是否待执行
func (t *OrgMigrationTask) IsPending() bool {
	return t.Status == MigrationTaskStatusPending
}

// IsRunning 是否执行中
func (t *OrgMigrationTask) IsRunning() bool {
	return t.Status == MigrationTaskStatusRunning
}

// IsCompleted 是否已完成
func (t *OrgMigrationTask) IsCompleted() bool {
	return t.Status == MigrationTaskStatusCompleted
}

// IsFailed 是否失败
func (t *OrgMigrationTask) IsFailed() bool {
	return t.Status == MigrationTaskStatusFailed
}

// IsRolledBack 是否已回滚
func (t *OrgMigrationTask) IsRolledBack() bool {
	return t.Status == MigrationTaskStatusRolledBack
}

// GetProgressPercent 获取进度百分比
func (t *OrgMigrationTask) GetProgressPercent() int {
	if t.TotalItems == 0 {
		return 0
	}
	return int((float64(t.MigratedItems) / float64(t.TotalItems)) * 100)
}

// OrgMigrationLog 迁移日志
type OrgMigrationLog struct {
	LogID       string              `json:"log_id" gorm:"primaryKey;type:varchar(36)"`
	TaskID      string              `json:"task_id" gorm:"type:varchar(36);not null;index:idx_task_id"`
	DataType    MigrationDataType    `json:"data_type" gorm:"type:enum('employees','departments','roles','permissions');not null"`
	Action      string              `json:"action" gorm:"type:varchar(50);not null"` // migrate|rollback
	EntityType  string              `json:"entity_type" gorm:"type:varchar(50);not null"`
	EntityID    string              `json:"entity_id" gorm:"type:varchar(36);not null"`
	Status      string              `json:"status" gorm:"type:varchar(20);not null"` // success|failed
	Message     string              `json:"message" gorm:"type:text"`
	CreatedAt   int64               `json:"created_at" gorm:"not null;default:0"`

	// 关联
	Task        *OrgMigrationTask    `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}

// TableName 指定表名
func (OrgMigrationLog) TableName() string {
	return "org_migration_logs"
}
