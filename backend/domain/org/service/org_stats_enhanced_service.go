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

package service

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// OrgStatsEnhancedService 组织统计增强服务
// 职责：部门维度统计、角色维度统计、活跃度统计、成本统计
type OrgStatsEnhancedService struct {
	employeeRepo  repository.EmployeeRepository
	departmentRepo repository.DepartmentRepository
	positionRepo  repository.PositionRepository
}

// NewOrgStatsEnhancedService 创建组织统计增强服务实例
func NewOrgStatsEnhancedService(
	employeeRepo repository.EmployeeRepository,
	departmentRepo repository.DepartmentRepository,
	positionRepo repository.PositionRepository,
) *OrgStatsEnhancedService {
	return &OrgStatsEnhancedService{
		employeeRepo:  employeeRepo,
		departmentRepo: departmentRepo,
		positionRepo:  positionRepo,
	}
}

// TimeRange 时间范围
type TimeRange struct {
	StartDate int64 `json:"start_date"`
	EndDate   int64 `json:"end_date"`
}

// GetEnhancedStatsRequest 获取增强统计请求
type GetEnhancedStatsRequest struct {
	TenantID   string      `json:"tenant_id" binding:"required"`
	OrgID      *string     `json:"org_id,omitempty"`
	Dimension  string      `json:"dimension" binding:"required,oneof=department role activity cost"`
	TimeRange  *TimeRange  `json:"time_range,omitempty"`
}

// EnhancedStatsResponse 增强统计响应
type EnhancedStatsResponse struct {
	Dimension string                 `json:"dimension"`
	Stats     map[string]interface{} `json:"stats"`
	Metadata  *StatsMetadata         `json:"metadata,omitempty"`
}

// StatsMetadata 统计元数据
type StatsMetadata struct {
	TotalCount  int64  `json:"total_count"`
	GeneratedAt int64  `json:"generated_at"`
}

// DepartmentStats 部门统计
type DepartmentStats struct {
	DeptID          string `json:"dept_id"`
	DeptName        string `json:"dept_name"`
	EmployeeCount   int64  `json:"employee_count"`
	ActiveCount     int64  `json:"active_count"`
	TrialCount      int64  `json:"trial_count"`
	ResignedCount   int64  `json:"resigned_count"`
	AverageJobLevel int64  `json:"average_job_level"`
}

// RoleStats 角色统计
type RoleStats struct {
	PositionID     string `json:"position_id"`
	PositionName   string `json:"position_name"`
	EmployeeCount  int64  `json:"employee_count"`
	AverageJobLevel int64 `json:"average_job_level"`
	MinJobLevel    int    `json:"min_job_level"`
	MaxJobLevel    int    `json:"max_job_level"`
}

// ActivityStats 活跃度统计
type ActivityStats struct {
	Date           string `json:"date"`
	ActiveUsers    int64  `json:"active_users"`
	NewHires       int64  `json:"new_hires"`
	Resignations   int64  `json:"resignations"`
	Transfers      int64  `json:"transfers"`
}

// CostStats 成本统计
type CostStats struct {
	Category       string  `json:"category"`
	TotalCost      float64 `json:"total_cost"`
	AverageCost    float64 `json:"average_cost"`
	EmployeeCount  int64   `json:"employee_count"`
}

// GetEnhancedStats 获取增强统计
func (s *OrgStatsEnhancedService) GetEnhancedStats(
	ctx context.Context,
	req *GetEnhancedStatsRequest,
) (*EnhancedStatsResponse, error) {
	// 1. 参数验证
	if req.TenantID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 2. 根据维度获取统计
	switch req.Dimension {
	case "department":
		return s.getDepartmentStats(ctx, req)
	case "role":
		return s.getRoleStats(ctx, req)
	case "activity":
		return s.getActivityStats(ctx, req)
	case "cost":
		return s.getCostStats(ctx, req)
	default:
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "InvalidDimension"),
		)
	}
}

// getDepartmentStats 获取部门统计
func (s *OrgStatsEnhancedService) getDepartmentStats(
	ctx context.Context,
	req *GetEnhancedStatsRequest,
) (*EnhancedStatsResponse, error) {
	// 1. 获取所有部门
	var departments []*entity.Department
	var err error

	if req.OrgID != nil && *req.OrgID != "" {
		// TODO: 需要添加GetByOrgID方法
		departments, err = s.departmentRepo.GetByTenantID(ctx, req.TenantID)
	} else {
		departments, err = s.departmentRepo.GetByTenantID(ctx, req.TenantID)
	}

	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get departments"),
		)
	}

	// 2. 统计各部门数据
	deptStats := make([]*DepartmentStats, 0, len(departments))
	var totalCount int64

	for _, dept := range departments {
		// 获取部门员工
		employees, err := s.employeeRepo.GetByDepartmentID(ctx, dept.DeptID)
		if err != nil {
			continue
		}

		stats := &DepartmentStats{
			DeptID:   dept.DeptID,
			DeptName: dept.DeptName,
		}

		var totalJobLevel int64
		for _, emp := range employees {
			stats.EmployeeCount++
			totalJobLevel += int64(emp.JobLevel)

			switch emp.EmployeeStatus {
			case entity.EmpStatusActive:
				stats.ActiveCount++
			case entity.EmpStatusTrial, entity.EmpStatusProbation:
				stats.TrialCount++
			case entity.EmpStatusResigned:
				stats.ResignedCount++
			}
		}

		if stats.EmployeeCount > 0 {
			stats.AverageJobLevel = totalJobLevel / stats.EmployeeCount
		}

		totalCount += stats.EmployeeCount
		deptStats = append(deptStats, stats)
	}

	// 3. 构建响应
	return &EnhancedStatsResponse{
		Dimension: "department",
		Stats: map[string]interface{}{
			"departments": deptStats,
		},
		Metadata: &StatsMetadata{
			TotalCount:  totalCount,
			GeneratedAt: getOrgStatsCurrentTimestamp(),
		},
	}, nil
}

// getRoleStats 获取角色统计
func (s *OrgStatsEnhancedService) getRoleStats(
	ctx context.Context,
	req *GetEnhancedStatsRequest,
) (*EnhancedStatsResponse, error) {
	// 1. 获取所有岗位
	positions, err := s.positionRepo.GetByTenantID(ctx, req.TenantID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get positions"),
		)
	}

	// 2. 统计各岗位数据
	roleStats := make([]*RoleStats, 0, len(positions))
	var totalCount int64

	for _, position := range positions {
		// 查询该岗位的员工
		filter := &repository.EmployeeFilter{
			TenantID:   req.TenantID,
			PositionID: &position.PositionID,
			PageToken:  "1",
			PageSize:   1000,
		}

		employees, total, err := s.employeeRepo.List(ctx, filter)
		if err != nil {
			continue
		}

		if total == 0 {
			continue
		}

		stats := &RoleStats{
			PositionID:   position.PositionID,
			PositionName: position.PositionName,
			EmployeeCount: total,
		}

		var totalJobLevel int64
		minJobLevel := 100
		maxJobLevel := 0

		for _, emp := range employees {
			totalJobLevel += int64(emp.JobLevel)

			if emp.JobLevel < minJobLevel {
				minJobLevel = emp.JobLevel
			}
			if emp.JobLevel > maxJobLevel {
				maxJobLevel = emp.JobLevel
			}
		}

		stats.AverageJobLevel = totalJobLevel / total
		stats.MinJobLevel = minJobLevel
		stats.MaxJobLevel = maxJobLevel

		totalCount += total
		roleStats = append(roleStats, stats)
	}

	// 3. 构建响应
	return &EnhancedStatsResponse{
		Dimension: "role",
		Stats: map[string]interface{}{
			"roles": roleStats,
		},
		Metadata: &StatsMetadata{
			TotalCount:  totalCount,
			GeneratedAt: getOrgStatsCurrentTimestamp(),
		},
	}, nil
}

// getActivityStats 获取活跃度统计
func (s *OrgStatsEnhancedService) getActivityStats(
	ctx context.Context,
	req *GetEnhancedStatsRequest,
) (*EnhancedStatsResponse, error) {
	// 1. 获取所有员工（用于统计）
	filter := &repository.EmployeeFilter{
		TenantID: req.TenantID,
		PageToken: "1",
		PageSize:  10000,
	}

	employees, total, err := s.employeeRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "list employees"),
		)
	}

	// 2. 统计活跃度数据
	activeCount := int64(0)
	trialCount := int64(0)
	resignedCount := int64(0)

	for _, emp := range employees {
		switch emp.EmployeeStatus {
		case entity.EmpStatusActive:
			activeCount++
		case entity.EmpStatusTrial, entity.EmpStatusProbation:
			trialCount++
		case entity.EmpStatusResigned:
			resignedCount++
		}
	}

	// 3. 构建响应
	activityStats := &ActivityStats{
		Date:         getCurrentDate(),
		ActiveUsers:  activeCount,
		NewHires:     trialCount,
		Resignations: resignedCount,
		Transfers:    0, // TODO: 需要调岗记录支持
	}

	return &EnhancedStatsResponse{
		Dimension: "activity",
		Stats: map[string]interface{}{
			"activity": activityStats,
		},
		Metadata: &StatsMetadata{
			TotalCount:  total,
			GeneratedAt: getOrgStatsCurrentTimestamp(),
		},
	}, nil
}

// getCostStats 获取成本统计
func (s *OrgStatsEnhancedService) getCostStats(
	ctx context.Context,
	req *GetEnhancedStatsRequest,
) (*EnhancedStatsResponse, error) {
	// 1. 获取所有员工
	filter := &repository.EmployeeFilter{
		TenantID: req.TenantID,
		PageToken: "1",
		PageSize:  10000,
	}

	employees, total, err := s.employeeRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "list employees"),
		)
	}

	// 2. 按职级统计成本
	costByLevel := make(map[int]*CostStats)
	for _, emp := range employees {
		if _, exists := costByLevel[emp.JobLevel]; !exists {
			costByLevel[emp.JobLevel] = &CostStats{
				Category:      fmt.Sprintf("P%d", emp.JobLevel),
				EmployeeCount: 0,
			}
		}

		costByLevel[emp.JobLevel].EmployeeCount++
		// TODO: 需要实际薪资数据支持
	}

	// 3. 转换为切片
	costStats := make([]*CostStats, 0, len(costByLevel))
	for _, stats := range costByLevel {
		costStats = append(costStats, stats)
	}

	// 4. 构建响应
	return &EnhancedStatsResponse{
		Dimension: "cost",
		Stats: map[string]interface{}{
			"costs": costStats,
		},
		Metadata: &StatsMetadata{
			TotalCount:  total,
			GeneratedAt: getOrgStatsCurrentTimestamp(),
		},
	}, nil
}

// GetDepartmentStatsDetail 获取部门统计详情
func (s *OrgStatsEnhancedService) GetDepartmentStatsDetail(
	ctx context.Context,
	tenantID, deptID string,
) (*DepartmentStats, error) {
	// 1. 验证部门存在
	dept, err := s.departmentRepo.GetByID(ctx, deptID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get department"),
		)
	}
	if dept == nil {
		return nil, errorx.NewByErrorCode(errno.ErrDeptNotFound)
	}
	if dept.TenantID != tenantID {
		return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
	}

	// 2. 获取部门员工
	employees, err := s.employeeRepo.GetByDepartmentID(ctx, deptID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get department employees"),
		)
	}

	// 3. 统计数据
	stats := &DepartmentStats{
		DeptID:   dept.DeptID,
		DeptName: dept.DeptName,
	}

	var totalJobLevel int64
	for _, emp := range employees {
		stats.EmployeeCount++
		totalJobLevel += int64(emp.JobLevel)

		switch emp.EmployeeStatus {
		case entity.EmpStatusActive:
			stats.ActiveCount++
		case entity.EmpStatusTrial, entity.EmpStatusProbation:
			stats.TrialCount++
		case entity.EmpStatusResigned:
			stats.ResignedCount++
		}
	}

	if stats.EmployeeCount > 0 {
		stats.AverageJobLevel = totalJobLevel / stats.EmployeeCount
	}

	return stats, nil
}

// GetRoleStatsDetail 获取角色统计详情
func (s *OrgStatsEnhancedService) GetRoleStatsDetail(
	ctx context.Context,
	tenantID, positionID string,
) (*RoleStats, error) {
	// 1. 验证岗位存在
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get position"),
		)
	}
	if position == nil {
		return nil, errorx.NewByErrorCode(errno.ErrPositionNotFound)
	}
	if position.TenantID != tenantID {
		return nil, errorx.NewByErrorCode(errno.ErrTenantMismatch)
	}

	// 2. 查询该岗位的员工
	filter := &repository.EmployeeFilter{
		TenantID:   tenantID,
		PositionID: &positionID,
		PageToken:  "1",
		PageSize:   1000,
	}

	employees, total, err := s.employeeRepo.List(ctx, filter)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "list employees by position"),
		)
	}

	// 3. 统计数据
	stats := &RoleStats{
		PositionID:   position.PositionID,
		PositionName: position.PositionName,
		EmployeeCount: total,
	}

	if total == 0 {
		return stats, nil
	}

	var totalJobLevel int64
	minJobLevel := 100
	maxJobLevel := 0

	for _, emp := range employees {
		totalJobLevel += int64(emp.JobLevel)

		if emp.JobLevel < minJobLevel {
			minJobLevel = emp.JobLevel
		}
		if emp.JobLevel > maxJobLevel {
			maxJobLevel = emp.JobLevel
		}
	}

	stats.AverageJobLevel = totalJobLevel / total
	stats.MinJobLevel = minJobLevel
	stats.MaxJobLevel = maxJobLevel

	return stats, nil
}

// 辅助函数
func getOrgStatsCurrentTimestamp() int64 {
	return 0 // TODO: 实现获取当前时间戳
}

func getCurrentDate() string {
	return "" // TODO: 实现获取当前日期
}
