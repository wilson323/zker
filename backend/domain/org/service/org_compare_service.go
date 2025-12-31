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

// OrgCompareService 组织对比服务
// 职责：多组织对比、对比维度分析、生成对比报告
type OrgCompareService struct {
	orgRepo     repository.OrganizationRepository
	employeeRepo repository.EmployeeRepository
	departmentRepo repository.DepartmentRepository
}

// NewOrgCompareService 创建组织对比服务实例
func NewOrgCompareService(
	orgRepo repository.OrganizationRepository,
	employeeRepo repository.EmployeeRepository,
	departmentRepo repository.DepartmentRepository,
) *OrgCompareService {
	return &OrgCompareService{
		orgRepo:     orgRepo,
		employeeRepo: employeeRepo,
		departmentRepo: departmentRepo,
	}
}

// CompareOrgsRequest 对比组织请求
type CompareOrgsRequest struct {
	OrgIDs     []string `json:"org_ids" binding:"required,min=2,max=10"`
	Dimensions []string `json:"dimensions" binding:"required"`
}

// OrgWithStats 带统计的组织
type OrgWithStats struct {
	Org   *entity.Organization `json:"org"`
	Stats *OrgStats            `json:"stats"`
}

// OrgStats 组织统计数据
type OrgStats struct {
	EmployeeCount    int64   `json:"employee_count"`
	DepartmentCount  int64   `json:"department_count"`
	ActiveCount      int64   `json:"active_count"`
	TrialCount       int64   `json:"trial_count"`
	AverageJobLevel  int64   `json:"average_job_level"`
	Cost             float64 `json:"cost,omitempty"`
	UsageRate        float64 `json:"usage_rate,omitempty"`
}

// CompareOrgsResponse 对比组织响应
type CompareOrgsResponse struct {
	Orgs       []*OrgWithStats   `json:"orgs"`
	Comparison *ComparisonReport `json:"comparison"`
}

// ComparisonReport 对比报告
type ComparisonReport struct {
	Dimensions    []string                      `json:"dimensions"`
	ComparisonData map[string]*DimensionCompare `json:"comparison_data"`
	Summary       string                        `json:"summary"`
}

// DimensionCompare 维度对比
type DimensionCompare struct {
	Dimension string       `json:"dimension"`
	Unit      string       `json:"unit"`
	Items     []*CompareItem `json:"items"`
	MaxValue  float64      `json:"max_value"`
	MinValue  float64      `json:"min_value"`
	AvgValue  float64      `json:"avg_value"`
}

// CompareItem 对比项
type CompareItem struct {
	OrgID    string  `json:"org_id"`
	OrgName  string  `json:"org_name"`
	Value    float64 `json:"value"`
	Rank     int     `json:"rank"`
	IsMax    bool    `json:"is_max"`
	IsMin    bool    `json:"is_min"`
}

// CompareOrgs 对比组织
func (s *OrgCompareService) CompareOrgs(
	ctx context.Context,
	req *CompareOrgsRequest,
) (*CompareOrgsResponse, error) {
	// 1. 参数验证
	if len(req.OrgIDs) < 2 || len(req.OrgIDs) > 10 {
		return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
			errorx.KV("reason", "InvalidOrgCount"),
		)
	}

	validDimensions := map[string]bool{
		"employee_count": true,
		"department_count": true,
		"active_count": true,
		"trial_count": true,
		"average_job_level": true,
		"cost": true,
		"usage_rate": true,
	}

	for _, dim := range req.Dimensions {
		if !validDimensions[dim] {
			return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
				errorx.KV("reason", "InvalidDimension"),
			)
		}
	}

	// 2. 查询各组织数据
	orgs := make([]*OrgWithStats, len(req.OrgIDs))
	for i, orgID := range req.OrgIDs {
		org, err := s.orgRepo.GetByID(ctx, orgID)
		if err != nil {
			return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
				errorx.KV("operation", "get organization"),
			)
		}
		if org == nil {
			return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
				errorx.KV("reason", "OrgNotFound"),
				errorx.KV("org_id", orgID),
			)
		}

		stats, err := s.getOrgStats(ctx, org, req.Dimensions)
		if err != nil {
			return nil, err
		}

		orgs[i] = &OrgWithStats{
			Org:   org,
			Stats: stats,
		}
	}

	// 3. 生成对比报告
	report := s.generateComparisonReport(orgs, req.Dimensions)

	return &CompareOrgsResponse{
		Orgs:       orgs,
		Comparison: report,
	}, nil
}

// getOrgStats 获取组织统计
func (s *OrgCompareService) getOrgStats(
	ctx context.Context,
	org *entity.Organization,
	dimensions []string,
) (*OrgStats, error) {
	stats := &OrgStats{}

	// 获取员工列表
	employees, err := s.employeeRepo.GetByOrgID(ctx, org.OrgID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get org employees"),
		)
	}

	// 获取部门列表
	departments, err := s.departmentRepo.GetByTenantID(ctx, org.TenantID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
			errorx.KV("operation", "get departments"),
		)
	}

	// 基础统计
	stats.EmployeeCount = int64(len(employees))
	stats.DepartmentCount = int64(len(departments))

	var totalJobLevel int64
	for _, emp := range employees {
		totalJobLevel += int64(emp.JobLevel)

		switch emp.EmployeeStatus {
		case entity.EmpStatusActive:
			stats.ActiveCount++
		case entity.EmpStatusTrial, entity.EmpStatusProbation:
			stats.TrialCount++
		}
	}

	if stats.EmployeeCount > 0 {
		stats.AverageJobLevel = totalJobLevel / stats.EmployeeCount
	}

	// 按维度计算其他统计
	for _, dim := range dimensions {
		switch dim {
		case "cost":
			// TODO: 实现成本统计
			stats.Cost = 0
		case "usage_rate":
			// TODO: 实现使用率统计
			stats.UsageRate = 0
		}
	}

	return stats, nil
}

// generateComparisonReport 生成对比报告
func (s *OrgCompareService) generateComparisonReport(
	orgs []*OrgWithStats,
	dimensions []string,
) *ComparisonReport {
	report := &ComparisonReport{
		Dimensions:     dimensions,
		ComparisonData: make(map[string]*DimensionCompare),
	}

	// 对每个维度进行对比
	for _, dim := range dimensions {
		compare := &DimensionCompare{
			Dimension: dim,
			Items:     make([]*CompareItem, len(orgs)),
		}

		// 设置单位
		switch dim {
		case "employee_count", "department_count", "active_count", "trial_count":
			compare.Unit = "人"
		case "average_job_level":
			compare.Unit = "级"
		case "cost":
			compare.Unit = "元"
		case "usage_rate":
			compare.Unit = "%"
		}

		// 收集数据
		values := make([]float64, len(orgs))
		for i, org := range orgs {
			value := s.getStatValue(org.Stats, dim)
			values[i] = value

			compare.Items[i] = &CompareItem{
				OrgID:   org.Org.OrgID,
				OrgName: org.Org.OrgName,
				Value:   value,
			}
		}

		// 计算最大值、最小值、平均值
		maxVal, minVal, avgVal := s.calculateValueStats(values)
		compare.MaxValue = maxVal
		compare.MinValue = minVal
		compare.AvgValue = avgVal

		// 排序并标记
		s.rankAndMarkItems(compare.Items)

		report.ComparisonData[dim] = compare
	}

	// 生成摘要
	report.Summary = s.generateSummary(report, orgs)

	return report
}

// getStatValue 获取统计值
func (s *OrgCompareService) getStatValue(stats *OrgStats, dimension string) float64 {
	switch dimension {
	case "employee_count":
		return float64(stats.EmployeeCount)
	case "department_count":
		return float64(stats.DepartmentCount)
	case "active_count":
		return float64(stats.ActiveCount)
	case "trial_count":
		return float64(stats.TrialCount)
	case "average_job_level":
		return float64(stats.AverageJobLevel)
	case "cost":
		return stats.Cost
	case "usage_rate":
		return stats.UsageRate
	default:
		return 0
	}
}

// calculateValueStats 计算值统计
func (s *OrgCompareService) calculateValueStats(values []float64) (max, min, avg float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}

	max = values[0]
	min = values[0]
	sum := 0.0

	for _, v := range values {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
		sum += v
	}

	avg = sum / float64(len(values))
	return
}

// rankAndMarkItems 排序并标记项
func (s *OrgCompareService) rankAndMarkItems(items []*CompareItem) {
	// 冒泡排序（按值降序）
	n := len(items)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if items[j].Value < items[j+1].Value {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}

	// 设置排名和标记
	maxVal := items[0].Value
	minVal := items[n-1].Value

	for i, item := range items {
		item.Rank = i + 1
		item.IsMax = item.Value == maxVal
		item.IsMin = item.Value == minVal
	}
}

// generateSummary 生成摘要
func (s *OrgCompareService) generateSummary(
	report *ComparisonReport,
	orgs []*OrgWithStats,
) string {
	if len(orgs) == 0 {
		return "无数据"
	}

	summary := fmt.Sprintf("共对比 %d 个组织，涵盖 %d 个维度。",
		len(orgs), len(report.Dimensions))

	// 找出各维度最优
	for dim, compare := range report.ComparisonData {
		if len(compare.Items) > 0 {
			topItem := compare.Items[0]
			dimName := s.getDimensionName(dim)
			summary += fmt.Sprintf("\n%s：%s 排名第一（%.2f %s）",
				dimName, topItem.OrgName, topItem.Value, compare.Unit)
		}
	}

	return summary
}

// getDimensionName 获取维度名称
func (s *OrgCompareService) getDimensionName(dimension string) string {
	names := map[string]string{
		"employee_count":   "员工总数",
		"department_count": "部门总数",
		"active_count":     "在职人数",
		"trial_count":      "试用期人数",
		"average_job_level": "平均职级",
		"cost":             "成本",
		"usage_rate":       "使用率",
	}

	if name, exists := names[dimension]; exists {
		return name
	}
	return dimension
}

// GetComparisonTrend 获取对比趋势（未来扩展）
func (s *OrgCompareService) GetComparisonTrend(
	ctx context.Context,
	orgIDs []string,
	dimension string,
	days int,
) (*ComparisonReport, error) {
	// TODO: 实现趋势对比功能
	return nil, errorx.New(errno.ErrPermissionInvalidParamCode,
		errorx.KV("reason", "NotImplemented"),
	)
}
