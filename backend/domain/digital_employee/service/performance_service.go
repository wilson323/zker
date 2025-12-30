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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/repository"
)

//go:generate mockgen -destination ../../../internal/mock/domain/digital_employee/service/performance_service_mock.go --package mockPerformanceService -source performance_service.go
type PerformanceService interface {
	// GetPerformance 获取员工绩效
	GetPerformance(ctx context.Context, req *entity.GetPerformanceRequest) (*entity.EmployeePerformance, error)

	// GetTeamPerformance 获取团队绩效
	GetTeamPerformance(ctx context.Context, req *entity.GetTeamPerformanceRequest) (*entity.GetTeamPerformanceResponse, error)

	// CalculateCompletionRate 计算完成率
	CalculateCompletionRate(ctx context.Context, employeeID string) (float64, error)

	// CalculateAvgResponseTime 计算平均响应时间
	CalculateAvgResponseTime(ctx context.Context, employeeID string) (int, error)

	// UpdatePerformance 更新绩效数据（任务完成后调用）
	UpdatePerformance(ctx context.Context, employeeID, tenantID string, taskCompleted bool, responseTime int) error
}

type performanceServiceImpl struct {
	performanceRepo repository.PerformanceRepository
	taskRepo        repository.TaskAssignmentRepository
	profileRepo     repository.ProfileRepository
}

// NewPerformanceService 创建绩效服务
func NewPerformanceService(
	performanceRepo repository.PerformanceRepository,
	taskRepo repository.TaskAssignmentRepository,
	profileRepo repository.ProfileRepository,
) PerformanceService {
	return &performanceServiceImpl{
		performanceRepo: performanceRepo,
		taskRepo:        taskRepo,
		profileRepo:     profileRepo,
	}
}

// GetPerformance 获取员工绩效
func (s *performanceServiceImpl) GetPerformance(ctx context.Context, req *entity.GetPerformanceRequest) (*entity.EmployeePerformance, error) {
	// 确定日期
	date := req.Date
	if date == nil {
		today := time.Now().Format("2006-01-02")
		date = &today
	}

	// 获取员工画像以获取租户ID
	profile, err := s.profileRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee profile: %w", err)
	}

	// 获取或创建绩效记录
	performance, err := s.performanceRepo.GetOrCreate(ctx, req.EmployeeID, profile.TenantID, req.Period, *date)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance: %w", err)
	}

	return performance, nil
}

// GetTeamPerformance 获取团队绩效
func (s *performanceServiceImpl) GetTeamPerformance(ctx context.Context, req *entity.GetTeamPerformanceRequest) (*entity.GetTeamPerformanceResponse, error) {
	performances, total, err := s.performanceRepo.GetByTenantID(ctx, req.TenantID, req.Period, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get team performance: %w", err)
	}

	return &entity.GetTeamPerformanceResponse{
		Performances: performances,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}, nil
}

// CalculateCompletionRate 计算完成率
func (s *performanceServiceImpl) CalculateCompletionRate(ctx context.Context, employeeID string) (float64, error) {
	total, completed, _, err := s.taskRepo.CountByEmployeeID(ctx, employeeID)
	if err != nil {
		return 0.0, fmt.Errorf("failed to count tasks: %w", err)
	}

	if total == 0 {
		return 0.0, nil
	}

	return (float64(completed) / float64(total)) * 100, nil
}

// CalculateAvgResponseTime 计算平均响应时间
func (s *performanceServiceImpl) CalculateAvgResponseTime(ctx context.Context, employeeID string) (int, error) {
	// TODO: 实现平均响应时间计算逻辑
	// 需要从任务分配表中计算平均响应时间
	return 0, nil
}

// UpdatePerformance 更新绩效数据
func (s *performanceServiceImpl) UpdatePerformance(ctx context.Context, employeeID, tenantID string, taskCompleted bool, responseTime int) error {
	// 确定周期和日期
	now := time.Now()
	period := entity.PerformancePeriodDaily
	date := now.Format("2006-01-02")

	// 获取或创建绩效记录
	performance, err := s.performanceRepo.GetOrCreate(ctx, employeeID, tenantID, period, date)
	if err != nil {
		return fmt.Errorf("failed to get or create performance: %w", err)
	}

	// 更新统计数据
	performance.TotalTasks++
	if taskCompleted {
		performance.CompletedTasks++
	} else {
		performance.FailedTasks++
	}

	// 计算完成率
	performance.CompletionRate = performance.CalculateCompletionRate()

	// 更新平均响应时间（简单平均）
	if responseTime > 0 {
		if performance.AvgResponseTime == 0 {
			performance.AvgResponseTime = responseTime
		} else {
			// 加权平均
			performance.AvgResponseTime = (performance.AvgResponseTime + responseTime) / 2
		}
	}

	// 保存更新
	if err := s.performanceRepo.Update(ctx, performance); err != nil {
		return fmt.Errorf("failed to update performance: %w", err)
	}

	return nil
}
