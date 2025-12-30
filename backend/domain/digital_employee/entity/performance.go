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
	"time"
)

// PerformancePeriod 统计周期
type PerformancePeriod string

const (
	PerformancePeriodDaily   PerformancePeriod = "daily"   // 日
	PerformancePeriodWeekly  PerformancePeriod = "weekly"  // 周
	PerformancePeriodMonthly PerformancePeriod = "monthly" // 月
)

// EmployeePerformance 员工绩效实体
type EmployeePerformance struct {
	ID              int64               `json:"id" gorm:"primaryKey;autoIncrement"`
	EmployeeID      string              `json:"employee_id" gorm:"type:varchar(36);not null;index:idx_employee_id"`
	TenantID        string              `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id"`
	TotalTasks      int                 `json:"total_tasks" gorm:"default:0"`
	CompletedTasks  int                 `json:"completed_tasks" gorm:"default:0"`
	FailedTasks     int                 `json:"failed_tasks" gorm:"default:0"`
	CompletionRate  float64             `json:"completion_rate" gorm:"type:decimal(5,2);default:0.00"`
	AvgResponseTime int                 `json:"avg_response_time" gorm:"default:0"` // 秒
	CustomerRating  float64             `json:"customer_rating" gorm:"type:decimal(3,2);default:0.00"` // 1-5分
	Period          PerformancePeriod   `json:"period" gorm:"type:enum('daily','weekly','monthly');not null"`
	Date            time.Time           `json:"date" gorm:"type:date;not null;index:idx_date"`
	UpdatedAt       int64               `json:"updated_at" gorm:"not null;default:0"`
}

// TableName 指定表名
func (EmployeePerformance) TableName() string {
	return "digital_employee_performance"
}

// GetUpdatedAtAsTime 获取更新时间
func (e *EmployeePerformance) GetUpdatedAtAsTime() time.Time {
	return time.Unix(e.UpdatedAt/1000, 0)
}

// CalculateCompletionRate 计算完成率
func (e *EmployeePerformance) CalculateCompletionRate() float64 {
	if e.TotalTasks == 0 {
		return 0.0
	}
	return (float64(e.CompletedTasks) / float64(e.TotalTasks)) * 100
}

// GetPerformanceRequest 获取员工绩效请求
type GetPerformanceRequest struct {
	EmployeeID string            `json:"employee_id" binding:"required,uuid"`
	Period     PerformancePeriod `json:"period" binding:"required,oneof=daily weekly monthly"`
	Date       *string           `json:"date" binding:"omitempty,datetime=2006-01-02"` // 格式: YYYY-MM-DD
}

// GetTeamPerformanceRequest 获取团队绩效请求
type GetTeamPerformanceRequest struct {
	TenantID string            `json:"tenant_id" binding:"required"`
	Period   PerformancePeriod `json:"period" binding:"required,oneof=daily weekly monthly"`
	Date     *string           `json:"date" binding:"omitempty,datetime=2006-01-02"`
	Page     int               `json:"page" binding:"required,min=1"`
	PageSize int               `json:"page_size" binding:"required,min=1,max=100"`
}

// GetTeamPerformanceResponse 获取团队绩效响应
type GetTeamPerformanceResponse struct {
	Performances []*EmployeePerformance `json:"performances"`
	Total        int64                  `json:"total"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"page_size"`
}
