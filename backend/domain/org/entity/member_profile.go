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

// SkillLevel 技能熟练度等级
type SkillLevel string

const (
	SkillLevelBeginner     SkillLevel = "BEGINNER"     // 初级
	SkillLevelIntermediate SkillLevel = "INTERMEDIATE" // 中级
	SkillLevelAdvanced     SkillLevel = "ADVANCED"     // 高级
	SkillLevelExpert       SkillLevel = "EXPERT"       // 专家
)

// DegreeDegree 学位等级
type DegreeDegree string

const (
	DegreeHighSchool DegreeDegree = "HIGH_SCHOOL" // 高中
	DegreeBachelor   DegreeDegree = "BACHELOR"    // 本科
	DegreeMaster     DegreeDegree = "MASTER"      // 硕士
	DegreePhD        DegreeDegree = "PHD"         // 博士
)

// WorkExperience 工作经历实体
type WorkExperience struct {
	ID          int64   `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      string  `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID    string  `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	CompanyName string  `json:"company_name" gorm:"type:varchar(255);not null"`
	Position    string  `json:"position" gorm:"type:varchar(128);not null"`
	StartDate   int64   `json:"start_date" gorm:"type:date;not null;index:idx_start_date"` // 毫秒时间戳
	EndDate     *int64  `json:"end_date,omitempty" gorm:"type:date;index:idx_end_date"`    // 毫秒时间戳
	Description *string `json:"description,omitempty" gorm:"type:text"`
	CreatedAt   int64   `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt   int64   `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt   *int64  `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (WorkExperience) TableName() string {
	return "work_experiences"
}

// IsDeleted 是否已删除
func (we *WorkExperience) IsDeleted() bool {
	return we.DeletedAt != nil
}

// IsCurrent 是否当前工作
func (we *WorkExperience) IsCurrent() bool {
	return we.EndDate == nil
}

// GetStartDateAsTime 获取开始时间
func (we *WorkExperience) GetStartDateAsTime() time.Time {
	return time.Unix(we.StartDate/1000, 0)
}

// GetEndDateAsTime 获取结束时间
func (we *WorkExperience) GetEndDateAsTime() *time.Time {
	if we.EndDate == nil {
		return nil
	}
	t := time.Unix(*we.EndDate/1000, 0)
	return &t
}

// Education 教育经历实体
type Education struct {
	ID         int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     string       `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID   string       `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	SchoolName string       `json:"school_name" gorm:"type:varchar(255);not null"`
	Major      string       `json:"major" gorm:"type:varchar(128);not null"`
	Degree     DegreeDegree `json:"degree" gorm:"type:enum('HIGH_SCHOOL','BACHELOR','MASTER','PHD');not null"`
	StartDate  int64        `json:"start_date" gorm:"type:date;not null;index:idx_start_date"` // 毫秒时间戳
	EndDate    *int64       `json:"end_date,omitempty" gorm:"type:date;index:idx_end_date"`    // 毫秒时间戳
	CreatedAt  int64        `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt  int64        `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt  *int64       `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Education) TableName() string {
	return "educations"
}

// IsDeleted 是否已删除
func (edu *Education) IsDeleted() bool {
	return edu.DeletedAt != nil
}

// IsCurrent 是否当前在读
func (edu *Education) IsCurrent() bool {
	return edu.EndDate == nil
}

// GetStartDateAsTime 获取开始时间
func (edu *Education) GetStartDateAsTime() time.Time {
	return time.Unix(edu.StartDate/1000, 0)
}

// GetEndDateAsTime 获取结束时间
func (edu *Education) GetEndDateAsTime() *time.Time {
	if edu.EndDate == nil {
		return nil
	}
	t := time.Unix(*edu.EndDate/1000, 0)
	return &t
}

// MemberSkill 成员技能实体
type MemberSkill struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      string     `json:"user_id" gorm:"type:varchar(36);not null;index:idx_user_id"`
	TenantID    string     `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	SkillName   string     `json:"skill_name" gorm:"type:varchar(128);not null;index:idx_skill_name"`
	Proficiency SkillLevel `json:"proficiency" gorm:"type:enum('BEGINNER','INTERMEDIATE','ADVANCED','EXPERT');not null"`
	Certified   bool       `json:"certified" gorm:"type:bool;default:false"` // 是否认证
	CreatedAt   int64      `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt   int64      `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt   *int64     `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (MemberSkill) TableName() string {
	return "member_skills"
}

// IsDeleted 是否已删除
func (ms *MemberSkill) IsDeleted() bool {
	return ms.DeletedAt != nil
}

// Resume 简历实体（聚合员工档案所有信息）
type Resume struct {
	UserID          string            `json:"user_id"`
	TenantID        string            `json:"tenant_id"`
	Employee        *Employee         `json:"employee,omitempty"`
	WorkExperiences []*WorkExperience `json:"work_experiences,omitempty"`
	Educations      []*Education      `json:"educations,omitempty"`
	Skills          []*MemberSkill    `json:"skills,omitempty"`
	GeneratedAt     int64             `json:"generated_at"`
}

// GetTotalWorkYears 获取总工作年限
func (r *Resume) GetTotalWorkYears() int {
	if len(r.WorkExperiences) == 0 {
		return 0
	}

	totalDays := 0
	for _, exp := range r.WorkExperiences {
		start := exp.GetStartDateAsTime()
		end := time.Now()
		if exp.EndDate != nil {
			endTime := exp.GetEndDateAsTime()
			if endTime != nil {
				end = *endTime
			}
		}
		totalDays += int(end.Sub(start).Hours() / 24)
	}

	return totalDays / 365
}

// GetHighestDegree 获取最高学历
func (r *Resume) GetHighestDegree() DegreeDegree {
	if len(r.Educations) == 0 {
		return ""
	}

	degreeOrder := map[DegreeDegree]int{
		DegreeHighSchool: 1,
		DegreeBachelor:   2,
		DegreeMaster:     3,
		DegreePhD:        4,
	}

	highestDegree := r.Educations[0].Degree
	for _, edu := range r.Educations {
		if degreeOrder[edu.Degree] > degreeOrder[highestDegree] {
			highestDegree = edu.Degree
		}
	}

	return highestDegree
}

// GetExpertSkills 获取专家级技能
func (r *Resume) GetExpertSkills() []string {
	var expertSkills []string
	for _, skill := range r.Skills {
		if skill.Proficiency == SkillLevelExpert {
			expertSkills = append(expertSkills, skill.SkillName)
		}
	}
	return expertSkills
}
