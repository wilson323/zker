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
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
)

// MatrixOrganizationService 矩阵组织服务
type MatrixOrganizationService struct {
	reportingRepo repository.MatrixReportingRepository
	empRepo       repository.EmployeeRepository
	orgRepo       repository.OrganizationRepository
}

// NewMatrixOrganizationService 创建矩阵组织服务实例
func NewMatrixOrganizationService(
	reportingRepo repository.MatrixReportingRepository,
	empRepo repository.EmployeeRepository,
	orgRepo repository.OrganizationRepository,
) *MatrixOrganizationService {
	return &MatrixOrganizationService{
		reportingRepo: reportingRepo,
		empRepo:       empRepo,
		orgRepo:       orgRepo,
	}
}

// CreateReporting 创建矩阵汇报关系
func (s *MatrixOrganizationService) CreateReporting(ctx context.Context, reporting *entity.MatrixReporting) error {
	// 1. 验证参数
	if reporting.UserID == "" || reporting.SupervisorID == "" || reporting.OrganizationID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}
	if reporting.ReportingType == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 2. 验证用户存在
	user, err := s.empRepo.GetByID(ctx, reporting.UserID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get user"),
            )
	}
	if user == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}

	// 3. 验证上级存在
	supervisor, err := s.empRepo.GetByID(ctx, reporting.SupervisorID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get supervisor"),
            )
	}
	if supervisor == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}

	// 4. 验证组织存在
	org, err := s.orgRepo.GetByID(ctx, reporting.OrganizationID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get organization"),
            )
	}
	if org == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "organization not found"),
                errorx.KV("resource_type", "organization"),
            )
	}

	// 5. 验证日期逻辑
	if reporting.ExpiryDate != nil && *reporting.ExpiryDate <= reporting.EffectiveDate {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidDateRange"),
                )
	}

	// 6. 如果设置为主要汇报关系,先取消其他主要关系
	if reporting.IsPrimary {
		primary, err := s.reportingRepo.GetPrimaryReporting(ctx, reporting.UserID)
		if err == nil && primary != nil {
			primary.IsPrimary = false
			if err := s.reportingRepo.Update(ctx, primary); err != nil {
				return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update old primary"),
            )
			}
		}
	}

	// 7. 创建汇报关系
	reporting.CreatedAt = time.Now().UnixMilli()
	reporting.UpdatedAt = time.Now().UnixMilli()
	if err := s.reportingRepo.Create(ctx, reporting); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create reporting"),
            )
	}

	return nil
}

// UpdateReporting 更新矩阵汇报关系
func (s *MatrixOrganizationService) UpdateReporting(ctx context.Context, id int64, reporting *entity.MatrixReporting) error {
	// 1. 获取现有记录
	existing, err := s.reportingRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get reporting"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "ReportingNotFound"),
                )
	}

	// 2. 验证日期逻辑
	if reporting.ExpiryDate != nil && *reporting.ExpiryDate <= reporting.EffectiveDate {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidDateRange"),
                )
	}

	// 3. 更新
	reporting.ID = id
	reporting.UpdatedAt = time.Now().UnixMilli()
	if err := s.reportingRepo.Update(ctx, reporting); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update reporting"),
            )
	}

	return nil
}

// DeleteReporting 删除矩阵汇报关系
func (s *MatrixOrganizationService) DeleteReporting(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.reportingRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get reporting"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "ReportingNotFound"),
                )
	}

	// 2. 删除
	if err := s.reportingRepo.Delete(ctx, id); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "delete reporting"),
            )
	}

	return nil
}

// GetReportings 获取用户的所有汇报关系
func (s *MatrixOrganizationService) GetReportings(ctx context.Context, userID string) ([]*entity.MatrixReporting, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	reportings, err := s.reportingRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list reportings"),
            )
	}

	return reportings, nil
}

// GetSubordinates 获取下属(包括矩阵式)
func (s *MatrixOrganizationService) GetSubordinates(ctx context.Context, supervisorID string, includeMatrix bool) ([]*entity.Employee, error) {
	if supervisorID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取所有汇报关系
	reportings, err := s.reportingRepo.ListBySupervisorID(ctx, supervisorID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list reportings"),
            )
	}

	// 2. 提取下级用户ID
	userIDs := make(map[string]bool)
	for _, r := range reportings {
		if includeMatrix || r.ReportingType == entity.ReportingTypeFunctional {
			userIDs[r.UserID] = true
		}
	}

	// 3. 获取员工信息
	var employees []*entity.Employee
	for userID := range userIDs {
		emp, err := s.empRepo.GetByID(ctx, userID)
		if err != nil {
			continue
		}
		if emp != nil {
			employees = append(employees, emp)
		}
	}

	return employees, nil
}

// GetSupervisors 获取所有上级(包括矩阵式)
func (s *MatrixOrganizationService) GetSupervisors(ctx context.Context, userID string) ([]*entity.SupervisorInfo, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取所有汇报关系
	reportings, err := s.reportingRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list reportings"),
            )
	}

	// 2. 构建上级信息
	var supervisors []*entity.SupervisorInfo
	for _, r := range reportings {
		if r.Supervisor == nil {
			continue
		}

		orgName := ""
		if r.Organization != nil {
			orgName = r.Organization.OrgName
		}

		info := &entity.SupervisorInfo{
			SupervisorID:   r.SupervisorID,
			SupervisorName: r.Supervisor.EmpName,
			ReportingType:  r.ReportingType,
			OrganizationID: r.OrganizationID,
			OrgName:        orgName,
			IsPrimary:      r.IsPrimary,
			EffectiveDate:  r.EffectiveDate,
			ExpiryDate:     r.ExpiryDate,
		}
		supervisors = append(supervisors, info)
	}

	return supervisors, nil
}

// GenerateOrgChart 生成组织架构图
func (s *MatrixOrganizationService) GenerateOrgChart(ctx context.Context, orgID string) (*entity.OrgChart, error) {
	// 1. 获取组织
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get org"),
            )
	}
	if org == nil {
		return nil, errorx.NewByErrorCode(errno.ErrOrgNotFound)
	}

	// 2. 获取组织负责人作为根节点
	if org.LeaderID == nil {
		return nil, errorx.NewByErrorCode(errno.ErrOrgLeaderNotFound)
	}

	leader, err := s.empRepo.GetByID(ctx, *org.LeaderID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get leader"),
            )
	}
	if leader == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	// 3. 构建组织架构图
	orgChart := &entity.OrgChart{
		EmployeeID: leader.EmpID,
		EmpName:    leader.EmpName,
		Position:   leader.JobTitle,
		Level:      1,
	}

	// 4. 递归获取下属
	subordinates, err := s.buildOrgChartRecursive(ctx, *org.LeaderID, orgID, 2)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "build org chart"),
            )
	}
	orgChart.Subordinates = subordinates

	return orgChart, nil
}

// buildOrgChartRecursive 递归构建组织架构图
func (s *MatrixOrganizationService) buildOrgChartRecursive(ctx context.Context, supervisorID, orgID string, level int) ([]*entity.OrgChart, error) {
	// 获取该上级在指定组织下的所有直属下级
	reportings, err := s.reportingRepo.ListBySupervisorID(ctx, supervisorID)
	if err != nil {
		return nil, err
	}

	var children []*entity.OrgChart
	for _, r := range reportings {
		// 只包含指定组织的汇报关系
		if r.OrganizationID != orgID {
			continue
		}

		// 构建子节点
		child := &entity.OrgChart{
			EmployeeID: r.UserID,
			EmpName:    r.User.EmpName,
			Position:   r.User.JobTitle,
			Level:      level,
		}

		// 递归获取下级
		subChildren, err := s.buildOrgChartRecursive(ctx, r.UserID, orgID, level+1)
		if err != nil {
			continue
		}
		child.Subordinates = subChildren

		children = append(children, child)
	}

	return children, nil
}

// CalculatePermissionPaths 计算权限路径
func (s *MatrixOrganizationService) CalculatePermissionPaths(ctx context.Context, userID string) ([]*entity.PermissionPath, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取所有汇报关系
	reportings, err := s.reportingRepo.ListActiveByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list reportings"),
            )
	}

	// 2. 构建权限路径
	var paths []*entity.PermissionPath
	for _, r := range reportings {
		path := &entity.PermissionPath{
			Path:         []string{r.UserID, r.SupervisorID, r.OrganizationID},
			OrgID:        r.OrganizationID,
			ReportingType: string(r.ReportingType),
			Depth:        1,
		}
		paths = append(paths, path)
	}

	return paths, nil
}
