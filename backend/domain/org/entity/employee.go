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

// EmployeeStatus 员工状态
type EmployeeStatus string

const (
	EmpStatusActive      EmployeeStatus = "active"       // 在职
	EmpStatusTrial       EmployeeStatus = "trial"        // 试用期
	EmpStatusProbation   EmployeeStatus = "probation"    // 试岗期
	EmpStatusResigned    EmployeeStatus = "resigned"     // 已离职
	EmpStatusRetired     EmployeeStatus = "retired"      // 退休
	EmpStatusSuspended   EmployeeStatus = "suspended"    // 停职
)

// EmployeeType 员工类型
type EmployeeType string

const (
	EmpTypeFullTime    EmployeeType = "full_time"    // 全职
	EmpTypePartTime    EmployeeType = "part_time"    // 兼职
	EmpTypeIntern      EmployeeType = "intern"        // 实习
	EmpTypeOutsourcing EmployeeType = "outsourcing"  // 外包
	EmpTypeContractor  EmployeeType = "contractor"   // 顾问
)

// Gender 性别
type Gender string

const (
	GenderMale   Gender = "male"   // 男
	GenderFemale Gender = "female" // 女
	GenderOther  Gender = "other"  // 其他
)

// Employee 员工实体
type Employee struct {
	EmpID         string         `json:"emp_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string         `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	UserID        *string        `json:"user_id,omitempty" gorm:"type:varchar(36);index:idx_user_id"` // 关联系统用户ID
	OrgID         string         `json:"org_id" gorm:"type:varchar(36);not null;index:idx_org_id"`
	DeptID        *string        `json:"dept_id,omitempty" gorm:"type:varchar(36);index:idx_dept_id"`
	PositionID    *string        `json:"position_id,omitempty" gorm:"type:varchar(36);index:idx_position_id"`

	// 基本信息
	EmpName      string         `json:"emp_name" gorm:"type:varchar(100);not null;index:idx_emp_name"`
	EmpCode      string         `json:"emp_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_tenant_code"` // 工号
	Gender       *Gender        `json:"gender,omitempty" gorm:"type:enum('male','female','other')"`
	Phone        *string        `json:"phone,omitempty" gorm:"type:varchar(20)"`
	Email        *string        `json:"email,omitempty" gorm:"type:varchar(100);index:idx_email"`

	// 职位信息
	EmployeeType EmployeeType   `json:"employee_type" gorm:"type:enum('full_time','part_time','intern','outsourcing','contractor');default:'full_time'"`
	EmployeeStatus EmployeeStatus `json:"employee_status" gorm:"type:enum('active','trial','probation','resigned','retired','suspended');default:'trial';index:idx_status"`
	JobLevel     int            `json:"job_level" gorm:"type:int;default:1;index:idx_job_level"` // 职级：P1-P10
	JobTitle     string         `json:"job_title" gorm:"type:varchar(100)"`                      // 职衔

	// 入职信息
	HireDate      int64          `json:"hire_date" gorm:"type:bigint;not null;index:idx_hire_date"`       // 入职日期（毫秒时间戳）
	RegularDate   *int64         `json:"regular_date,omitempty" gorm:"type:bigint;index:idx_regular_date"` // 转正日期
	ProbationDays int            `json:"probation_days" gorm:"type:int;default:90"`                       // 试用天数

	// 工作信息
	WorkLocation *string        `json:"work_location,omitempty" gorm:"type:varchar(200)"` // 工作地点
	DirectLeaderID *string      `json:"direct_leader_id,omitempty" gorm:"type:varchar(36)"` // 直接上级ID

	// 个人信息
	IDCard      *string        `json:"id_card,omitempty" gorm:"type:varchar(18)"`
	Birthday    *int64         `json:"birthday,omitempty" gorm:"type:bigint"` // 生日（毫秒时间戳）
	Address     *string        `json:"address,omitempty" gorm:"type:varchar(500)"`
	Education   *string        `json:"education,omitempty" gorm:"type:varchar(50)"` // 学历
	GraduateSchool *string     `json:"graduate_school,omitempty" gorm:"type:varchar(200)"` // 毕业院校
	Major       *string        `json:"major,omitempty" gorm:"type:varchar(100)"` // 专业

	// 紧急联系人
	EmergencyContact   *string `json:"emergency_contact,omitempty" gorm:"type:varchar(100)"`
	EmergencyPhone     *string `json:"emergency_phone,omitempty" gorm:"type:varchar(20)"`

	// 其他
	Status      EmployeeStatus `json:"status" gorm:"type:enum('active','trial','probation','resigned','retired','suspended');default:'trial';index:idx_status"`
	AvatarURL   *string        `json:"avatar_url,omitempty" gorm:"type:varchar(500)"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`

	// 时间戳
	CreatedAt   int64          `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt   int64          `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt   *int64         `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Organization *Organization `json:"organization,omitempty" gorm:"foreignKey:OrgID"`
	Department   *Department   `json:"department,omitempty" gorm:"foreignKey:DeptID"`
	Position     *Position     `json:"position,omitempty" gorm:"foreignKey:PositionID"`
	DirectLeader *Employee     `json:"direct_leader,omitempty" gorm:"foreignKey:DirectLeaderID"`
}

// TableName 指定表名
func (Employee) TableName() string {
	return "employees"
}

// IsActive 是否在职
func (e *Employee) IsActive() bool {
	return e.EmployeeStatus == EmpStatusActive || e.EmployeeStatus == EmpStatusTrial || e.EmployeeStatus == EmpStatusProbation
}

// IsDeleted 是否已删除
func (e *Employee) IsDeleted() bool {
	return e.DeletedAt != nil
}

// IsTrial 是否在试用期
func (e *Employee) IsTrial() bool {
	return e.EmployeeStatus == EmpStatusTrial
}

// GetHireDateAsTime 获取入职时间
func (e *Employee) GetHireDateAsTime() time.Time {
	return time.Unix(e.HireDate/1000, 0)
}

// GetRegularDateAsTime 获取转正时间
func (e *Employee) GetRegularDateAsTime() *time.Time {
	if e.RegularDate == nil {
		return nil
	}
	t := time.Unix(*e.RegularDate/1000, 0)
	return &t
}

// GetBirthdayAsTime 获取生日
func (e *Employee) GetBirthdayAsTime() *time.Time {
	if e.Birthday == nil {
		return nil
	}
	t := time.Unix(*e.Birthday/1000, 0)
	return &t
}

// GetCreatedAtAsTime 获取创建时间
func (e *Employee) GetCreatedAtAsTime() time.Time {
	return time.Unix(e.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (e *Employee) GetUpdatedAtAsTime() time.Time {
	return time.Unix(e.UpdatedAt/1000, 0)
}

// EmployeeContract 员工合同实体
type EmployeeContract struct {
	ContractID    string    `json:"contract_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	EmpID         string    `json:"emp_id" gorm:"type:varchar(36);not null;index:idx_emp_id"`

	ContractType  string    `json:"contract_type" gorm:"type:varchar(50);not null"` // 合同类型：劳动合同/实习协议/外包协议
	ContractNo    string    `json:"contract_no" gorm:"type:varchar(100);not null;uniqueIndex:uk_tenant_no"`

	StartDate     int64     `json:"start_date" gorm:"type:bigint;not null;index:idx_start_date"` // 合同开始日期
	EndDate       *int64    `json:"end_date,omitempty" gorm:"type:bigint;index:idx_end_date"`     // 合同结束日期（无固定期限则为NULL）

	Salary        *string   `json:"salary,omitempty" gorm:"type:varchar(50)"` // 薪资（密文字段）
	SalaryType    string    `json:"salary_type" gorm:"type:varchar(20);default:'monthly'"` // 薪资类型：monthly/yearly

	ProbationDays int       `json:"probation_days" gorm:"type:int;default:0"` // 试用天数
	ProbationSalary *string `json:"probation_salary,omitempty" gorm:"type:varchar(50)"` // 试用期薪资

	WorkHours     string    `json:"work_hours" gorm:"type:varchar(50);default:'8:00-17:00'"` // 工作时间
	WorkPlace     *string   `json:"work_place,omitempty" gorm:"type:varchar(200)"`            // 工作地点

	ContractFileURL *string `json:"contract_file_url,omitempty" gorm:"type:varchar(500)"` // 合同文件URL

	Status       string    `json:"status" gorm:"type:enum('draft','active','expired','terminated');default:'draft';index:idx_status"` // 草稿/生效/过期/终止

	SignedAt     *int64    `json:"signed_at,omitempty" gorm:"type:bigint"` // 签署时间

	// 时间戳
	CreatedAt    int64     `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64     `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt    *int64    `json:"deleted_at,omitempty" gorm:"index"`

	// 关联
	Employee     *Employee `json:"employee,omitempty" gorm:"foreignKey:EmpID"`
}

// TableName 指定表名
func (EmployeeContract) TableName() string {
	return "employee_contracts"
}

// IsDeleted 是否已删除
func (ec *EmployeeContract) IsDeleted() bool {
	return ec.DeletedAt != nil
}

// IsActive 是否生效
func (ec *EmployeeContract) IsActive() bool {
	return ec.Status == "active"
}

// GetStartDateAsTime 获取开始时间
func (ec *EmployeeContract) GetStartDateAsTime() time.Time {
	return time.Unix(ec.StartDate/1000, 0)
}

// GetEndDateAsTime 获取结束时间
func (ec *EmployeeContract) GetEndDateAsTime() *time.Time {
	if ec.EndDate == nil {
		return nil
	}
	t := time.Unix(*ec.EndDate/1000, 0)
	return &t
}

// EmployeeTransfer 员工调岗记录
type EmployeeTransfer struct {
	TransferID   string    `json:"transfer_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID     string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	EmpID        string    `json:"emp_id" gorm:"type:varchar(36);not null;index:idx_emp_id"`

	// 变更前
	OldDeptID    *string   `json:"old_dept_id,omitempty" gorm:"type:varchar(36)"`
	OldPositionID *string  `json:"old_position_id,omitempty" gorm:"type:varchar(36)"`
	OldJobLevel   int      `json:"old_job_level" gorm:"type:int;default:1"`
	OldJobTitle   string   `json:"old_job_title" gorm:"type:varchar(100)"`

	// 变更后
	NewDeptID    *string   `json:"new_dept_id,omitempty" gorm:"type:varchar(36)"`
	NewPositionID *string  `json:"new_position_id,omitempty" gorm:"type:varchar(36)"`
	NewJobLevel   int      `json:"new_job_level" gorm:"type:int;default:1"`
	NewJobTitle   string   `json:"new_job_title" gorm:"type:varchar(100)"`

	// 调岗信息
	TransferType string    `json:"transfer_type" gorm:"type:varchar(50);not null"` // 调岗/晋升/降职/平调
	TransferDate int64     `json:"transfer_date" gorm:"type:bigint;not null;index:idx_transfer_date"`
	Reason       string    `json:"reason" gorm:"type:text"` // 调岗原因

	ApproverID   *string   `json:"approver_id,omitempty" gorm:"type:varchar(36)"` // 审批人ID
	ApprovedAt   *int64    `json:"approved_at,omitempty" gorm:"type:bigint"`

	// 时间戳
	CreatedAt    int64     `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt    int64     `json:"updated_at" gorm:"not null;default:0"`

	// 关联
	Employee     *Employee `json:"employee,omitempty" gorm:"foreignKey:EmpID"`
}

// TableName 指定表名
func (EmployeeTransfer) TableName() string {
	return "employee_transfers"
}

// GetTransferDateAsTime 获取调岗时间
func (et *EmployeeTransfer) GetTransferDateAsTime() time.Time {
	return time.Unix(et.TransferDate/1000, 0)
}

// EmployeeResignation 员工离职记录
type EmployeeResignation struct {
	ResignationID string    `json:"resignation_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string    `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	EmpID         string    `json:"emp_id" gorm:"type:varchar(36);not null;index:idx_emp_id"`

	ResignationType string  `json:"resignation_type" gorm:"type:varchar(50);not null"` // 离职类型：主动离职/被动离职/合同到期/退休
	ResignationReason string `json:"resignation_reason" gorm:"type:text"` // 离职原因

	// 时间信息
	ApplyDate        int64   `json:"apply_date" gorm:"type:bigint;not null;index:idx_apply_date"` // 申请日期
	LastWorkDate     int64   `json:"last_work_date" gorm:"type:bigint;not null;index:idx_last_work_date"` // 最后工作日
	ResignationDate  int64   `json:"resignation_date" gorm:"type:bigint;not null;index:idx_resignation_date"` // 离职生效日期

	// 流程信息
	HandoverToID    *string `json:"handover_to_id,omitempty" gorm:"type:varchar(36)"` // 工作交接人
	HandoverStatus  string  `json:"handover_status" gorm:"type:varchar(20);default:'pending'"` // 交接状态：pending/completed

	ApprovalStatus  string  `json:"approval_status" gorm:"type:varchar(20);default:'pending';index:idx_approval"` // 审批状态：pending/approved/rejected
	ApproverID      *string `json:"approver_id,omitempty" gorm:"type:varchar(36)"` // 审批人ID
	ApprovedAt      *int64  `json:"approved_at,omitempty" gorm:"type:bigint"`
	ApprovalComment *string `json:"approval_comment,omitempty" gorm:"type:text"`

	// 离职后信息
	RehireEligible  bool    `json:"rehire_eligible" gorm:"type:bool;default:true"` // 是否可再录用
	Notes           *string `json:"notes,omitempty" gorm:"type:text"` // 备注

	// 时间戳
	CreatedAt       int64   `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt       int64   `json:"updated_at" gorm:"not null;default:0"`

	// 关联
	Employee        *Employee `json:"employee,omitempty" gorm:"foreignKey:EmpID"`
}

// TableName 指定表名
func (EmployeeResignation) TableName() string {
	return "employee_resignations"
}

// GetApplyDateAsTime 获取申请时间
func (er *EmployeeResignation) GetApplyDateAsTime() time.Time {
	return time.Unix(er.ApplyDate/1000, 0)
}

// GetLastWorkDateAsTime 获取最后工作日
func (er *EmployeeResignation) GetLastWorkDateAsTime() time.Time {
	return time.Unix(er.LastWorkDate/1000, 0)
}

// GetResignationDateAsTime 获取离职生效时间
func (er *EmployeeResignation) GetResignationDateAsTime() time.Time {
	return time.Unix(er.ResignationDate/1000, 0)
}
