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

import (
	"encoding/json"
	"time"
)

// EmployeeRole 员工角色类型
type EmployeeRole string

const (
	EmployeeRoleCustomerService EmployeeRole = "customer_service" // 客服
	EmployeeRoleSales           EmployeeRole = "sales"            // 销售
	EmployeeRoleTechSupport     EmployeeRole = "tech_support"     // 技术支持
	EmployeeRoleConsultant      EmployeeRole = "consultant"       // 顾问
	EmployeeRoleTrainer         EmployeeRole = "trainer"          // 培训师
)

// EmployeeStatus 员工状态
type EmployeeStatus string

const (
	EmployeeStatusActive   EmployeeStatus = "active"   // 激活
	EmployeeStatusInactive EmployeeStatus = "inactive" // 停用
	EmployeeStatusDeleted  EmployeeStatus = "deleted"  // 已删除
)

// EmployeeProfile 员工画像实体
type EmployeeProfile struct {
	EmployeeID    string         `json:"employee_id" gorm:"primaryKey;type:varchar(36)"`
	TenantID      string         `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	Name          string         `json:"name" gorm:"type:varchar(100);not null"`
	Avatar        string         `json:"avatar" gorm:"type:varchar(255)"`
	Role          EmployeeRole   `json:"role" gorm:"type:enum('customer_service','sales','tech_support','consultant','trainer');not null;index:idx_role"`
	Skills        []string       `json:"skills" gorm:"type:json"` // JSON数组
	Personality   string         `json:"personality" gorm:"type:varchar(100)"`
	KnowledgeBaseID *string      `json:"knowledge_base_id,omitempty" gorm:"type:varchar(36)"`
	BotID         *string        `json:"bot_id,omitempty" gorm:"type:varchar(36)"`
	Status        EmployeeStatus `json:"status" gorm:"type:enum('active','inactive','deleted');default:'active';index:idx_status"`
	CreatedAt     int64          `json:"created_at" gorm:"not null;default:0"`
	UpdatedAt     int64          `json:"updated_at" gorm:"not null;default:0"`
	DeletedAt     *int64         `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (EmployeeProfile) TableName() string {
	return "digital_employee_profiles"
}

// IsActive 是否激活
func (e *EmployeeProfile) IsActive() bool {
	return e.Status == EmployeeStatusActive
}

// IsDeleted 是否已删除
func (e *EmployeeProfile) IsDeleted() bool {
	return e.Status == EmployeeStatusDeleted || e.DeletedAt != nil
}

// GetCreatedAtAsTime 获取创建时间
func (e *EmployeeProfile) GetCreatedAtAsTime() time.Time {
	return time.Unix(e.CreatedAt/1000, 0)
}

// GetUpdatedAtAsTime 获取更新时间
func (e *EmployeeProfile) GetUpdatedAtAsTime() time.Time {
	return time.Unix(e.UpdatedAt/1000, 0)
}

// SkillsToString 将技能数组转为JSON字符串
func (e *EmployeeProfile) SkillsToString() (string, error) {
	if e.Skills == nil {
		return "[]", nil
	}
	data, err := json.Marshal(e.Skills)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// StringToSkills 从JSON字符串解析技能数组
func (e *EmployeeProfile) StringToSkills(s string) error {
	if s == "" || s == "[]" {
		e.Skills = []string{}
		return nil
	}
	return json.Unmarshal([]byte(s), &e.Skills)
}

// CreateProfileRequest 创建员工画像请求
type CreateProfileRequest struct {
	TenantID      string       `json:"tenant_id" binding:"required"`
	Name          string       `json:"name" binding:"required,min=1,max=100"`
	Avatar        string       `json:"avatar" binding:"omitempty,max=255"`
	Role          EmployeeRole `json:"role" binding:"required,oneof=customer_service sales tech_support consultant trainer"`
	Skills        []string     `json:"skills" binding:"required"`
	Personality   string       `json:"personality" binding:"omitempty,max=100"`
	KnowledgeBaseID *string    `json:"knowledge_base_id,omitempty" binding:"omitempty,uuid"`
	BotID         *string      `json:"bot_id,omitempty" binding:"omitempty,uuid"`
}

// UpdateProfileRequest 更新员工画像请求
type UpdateProfileRequest struct {
	EmployeeID    string       `json:"employee_id" binding:"required,uuid"`
	Name          string       `json:"name" binding:"omitempty,min=1,max=100"`
	Avatar        string       `json:"avatar" binding:"omitempty,max=255"`
	Role          EmployeeRole `json:"role" binding:"omitempty,oneof=customer_service sales tech_support consultant trainer"`
	Skills        []string     `json:"skills" binding:"omitempty"`
	Personality   string       `json:"personality" binding:"omitempty,max=100"`
	KnowledgeBaseID *string    `json:"knowledge_base_id,omitempty" binding:"omitempty,uuid"`
}

// ListProfilesRequest 列出员工画像请求
type ListProfilesRequest struct {
	TenantID string         `json:"tenant_id" binding:"required"`
	Role     *EmployeeRole  `json:"role" binding:"omitempty,oneof=customer_service sales tech_support consultant trainer"`
	Status   *EmployeeStatus `json:"status" binding:"omitempty,oneof=active inactive deleted"`
	Page     int            `json:"page" binding:"required,min=1"`
	PageSize int            `json:"page_size" binding:"required,min=1,max=100"`
}

// ListProfilesResponse 列出员工画像响应
type ListProfilesResponse struct {
	Profiles []*EmployeeProfile `json:"profiles"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}
