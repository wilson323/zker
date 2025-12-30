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

// MigrationStatus 迁移状态
type MigrationStatus string

const (
	MigrationStatusPending   MigrationStatus = "pending"    // 待执行
	MigrationStatusRunning   MigrationStatus = "running"    // 执行中
	MigrationStatusCompleted MigrationStatus = "completed"  // 已完成
	MigrationStatusFailed    MigrationStatus = "failed"     // 失败
	MigrationStatusRolledBack MigrationStatus = "rolled_back" // 已回滚
)

// MigrationType 迁移类型
type MigrationType string

const (
	MigrationTypeSchema      MigrationType = "schema"       // 结构迁移
	MigrationTypeData        MigrationType = "data"         // 数据迁移
	MigrationTypeBackfill    MigrationType = "backfill"     // 数据回填
	MigrationTypeCutover     MigrationType = "cutover"      // 割接
	MigrationTypeVerification MigrationType = "verification" // 验证
)

// MigrationAudit 迁移审计实体
type MigrationAudit struct {
	ID               int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	MigrationID      string          `json:"migration_id" gorm:"type:varchar(64);uniqueIndex;not null;comment:迁移批次ID"`
	MigrationName    string          `json:"migration_name" gorm:"type:varchar(200);not null;comment:迁移名称"`
	MigrationType    MigrationType   `json:"migration_type" gorm:"type:enum('schema','data','backfill','cutover','verification');not null;comment:迁移类型"`
	Status           MigrationStatus `json:"status" gorm:"type:enum('pending','running','completed','failed','rolled_back');not null;default:'pending';comment:迁移状态"`
	StartTime        int64           `json:"start_time" gorm:"not null;comment:开始时间"`
	EndTime          *int64          `json:"end_time,omitempty" gorm:"comment:结束时间"`
	DurationSeconds  *int64          `json:"duration_seconds,omitempty" gorm:"comment:耗时（秒）"`
	TotalRecords     *int64          `json:"total_records,omitempty" gorm:"comment:总记录数"`
	ProcessedRecords *int64          `json:"processed_records,omitempty" gorm:"comment:已处理记录数"`
	FailedRecords    *int64          `json:"failed_records,omitempty" gorm:"comment:失败记录数"`
	ErrorMessage     *string         `json:"error_message,omitempty" gorm:"type:text;comment:错误信息"`
	ExecutedBy       string          `json:"executed_by" gorm:"type:varchar(100);not null;comment:执行人"`
	Metadata         string          `json:"metadata" gorm:"type:json;comment:元数据（JSON格式）"`
	CreatedAt        int64           `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt        int64           `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (MigrationAudit) TableName() string {
	return "migration_audit"
}

// IsCompleted 是否已完成
func (m *MigrationAudit) IsCompleted() bool {
	return m.Status == MigrationStatusCompleted
}

// IsFailed 是否失败
func (m *MigrationAudit) IsFailed() bool {
	return m.Status == MigrationStatusFailed
}

// IsRunning 是否执行中
func (m *MigrationAudit) IsRunning() bool {
	return m.Status == MigrationStatusRunning
}

// GetDuration 获取耗时（秒）
func (m *MigrationAudit) GetDuration() int64 {
	if m.DurationSeconds != nil {
		return *m.DurationSeconds
	}
	if m.EndTime != nil {
		return *m.EndTime - m.StartTime
	}
	return 0
}

// GetProgress 获取进度百分比
func (m *MigrationAudit) GetProgress() float64 {
	if m.TotalRecords == nil || *m.TotalRecords == 0 {
		return 0
	}
	if m.ProcessedRecords == nil {
		return 0
	}
	return float64(*m.ProcessedRecords) / float64(*m.TotalRecords) * 100
}
