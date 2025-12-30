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

// ReportingType 汇报关系类型
type ReportingType string

const (
	ReportingTypeFunctional ReportingType = "FUNCTIONAL" // 职能汇报（直线）
	ReportingTypeProject    ReportingType = "PROJECT"    // 项目汇报（横向）
	ReportingTypeDotted     ReportingType = "DOTTED"     // 虚线汇报（非正式）
)

// MatrixReporting 矩阵汇报关系实体
type MatrixReporting struct {
	ID             int64         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         string        `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	ReportingType  ReportingType `json:"reporting_type" gorm:"type:enum('FUNCTIONAL','PROJECT','DOTTED');not null;index:idx_reporting_type"`
	SupervisorID   string        `json:"supervisor_id" gorm:"type:varchar(36);not null;index:idx_supervisor_id"`
	OrganizationID string        `json:"organization_id" gorm:"type:varchar(36);not null;index:idx_org_id"`
	EffectiveDate  int64         `json:"effective_date" gorm:"type:date;not null;index:idx_effective_date"` // 生效日期（毫秒时间戳）
	ExpiryDate     *int64        `json:"expiry_date,omitempty" gorm:"type:date;index:idx_expiry_date"`        // 失效日期（毫秒时间戳）
	IsPrimary      bool          `json:"is_primary" gorm:"type:bool;default:false;index:idx_is_primary"`     // 是否主要汇报关系
	CreatedAt      int64         `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt      int64         `json:"updated_at" gorm:"not null;default:0"`

	// 关联
	User       *Employee      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Supervisor *Employee      `json:"supervisor,omitempty" gorm:"foreignKey:SupervisorID"`
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
}

// TableName 指定表名
func (MatrixReporting) TableName() string {
	return "matrix_reportings"
}

// IsActive 是否有效
func (mr *MatrixReporting) IsActive() bool {
	now := time.Now()
	effectiveDate := time.Unix(mr.EffectiveDate/1000, 0)

	// 未到生效日期
	if now.Before(effectiveDate) {
		return false
	}

	// 已过期
	if mr.ExpiryDate != nil {
		expiryDate := time.Unix(*mr.ExpiryDate/1000, 0)
		if now.After(expiryDate) {
			return false
		}
	}

	return true
}

// GetEffectiveDateAsTime 获取生效时间
func (mr *MatrixReporting) GetEffectiveDateAsTime() time.Time {
	return time.Unix(mr.EffectiveDate/1000, 0)
}

// GetExpiryDateAsTime 获取失效时间
func (mr *MatrixReporting) GetExpiryDateAsTime() *time.Time {
	if mr.ExpiryDate == nil {
		return nil
	}
	t := time.Unix(*mr.ExpiryDate/1000, 0)
	return &t
}

// SupervisorInfo 上级信息
type SupervisorInfo struct {
	SupervisorID    string        `json:"supervisor_id"`
	SupervisorName  string        `json:"supervisor_name"`
	ReportingType   ReportingType `json:"reporting_type"`
	OrganizationID  string        `json:"organization_id"`
	OrgName         string        `json:"org_name"`
	IsPrimary       bool          `json:"is_primary"`
	EffectiveDate   int64         `json:"effective_date"`
	ExpiryDate      *int64        `json:"expiry_date,omitempty"`
}

// OrgChart 组织架构图节点
type OrgChart struct {
	EmployeeID    string       `json:"employee_id"`
	EmpName       string       `json:"emp_name"`
	Position      string       `json:"position"`
	Level         int          `json:"level"`
	Subordinates  []*OrgChart  `json:"subordinates,omitempty"`
	Reportings    []ReportingInfo `json:"reportings,omitempty"`
}

// ReportingInfo 汇报关系信息
type ReportingInfo struct {
	ReportingType ReportingType `json:"reporting_type"`
	OrgName       string        `json:"org_name"`
	IsPrimary     bool          `json:"is_primary"`
}

// PermissionPath 权限路径
type PermissionPath struct {
	Path         []string `json:"path"`         // 路径：user -> supervisor -> org
	OrgID        string   `json:"org_id"`       // 组织ID
	ReportingType string   `json:"reporting_type"` // 汇报类型
	Depth        int      `json:"depth"`        // 深度
}
