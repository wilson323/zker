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
	"encoding/json"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// VisualizationService 可视化服务
// 职责：组织架构可视化、组织关系网络、组织对比
type VisualizationService struct {
	orgRepo    repository.OrganizationRepository
	visRepo    repository.VisualizationRepository
	empRepo    repository.EmployeeRepository
	deptRepo   repository.DepartmentRepository
	statsRepo  repository.OrganizationStatsRepository
}

// NewVisualizationService 创建可视化服务实例
func NewVisualizationService(
	orgRepo repository.OrganizationRepository,
	visRepo repository.VisualizationRepository,
	empRepo repository.EmployeeRepository,
	deptRepo repository.DepartmentRepository,
	statsRepo repository.OrganizationStatsRepository,
) *VisualizationService {
	return &VisualizationService{
		orgRepo:   orgRepo,
		visRepo:   visRepo,
		empRepo:   empRepo,
		deptRepo:  deptRepo,
		statsRepo: statsRepo,
	}
}

// GetVisualizationRequest 获取可视化请求
type GetVisualizationRequest struct {
	TenantID  string                      `json:"tenant_id" binding:"required"`
	OrgID     string                      `json:"org_id,omitempty"`
	Type      entity.VisualizationType    `json:"type" binding:"required,oneof=hierarchy matrix department network"`
	Config    entity.VisualizationConfig  `json:"config"`
}

// GetOrganizationTree 获取组织树可视化
func (s *VisualizationService) GetOrganizationTree(ctx context.Context, req *GetVisualizationRequest) (*entity.VisualizationResult, error) {
	// 1. 获取组织树
	orgs, err := s.visRepo.GetOrganizationTreeForVisualization(ctx, req.TenantID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get organization tree failed: tenant_id=%s", req.TenantID)
	}

	if len(orgs) == 0 {
		return nil, errorx.NewByErrorCode(errno.ErrOrgNotFound)
	}

	// 2. 批量获取员工数
	orgIDs := make([]string, 0, len(orgs))
	for _, org := range orgs {
		orgIDs = append(orgIDs, org.OrgID)
	}
	empCounts, err := s.visRepo.BatchGetEmployeeCount(ctx, orgIDs)
	if err != nil {
		return nil, errorx.Wrapf(err, "batch get employee count failed")
	}

	// 3. 构建可视化节点
	nodes := s.buildOrganizationNodes(orgs, empCounts, &req.Config)

	// 4. 构建统计信息
	stats := s.buildVisualizationStatistics(nodes)

	// 5. 返回可视化结果
	result := &entity.VisualizationResult{
		Type:        req.Type,
		Config:      req.Config,
		Nodes:       nodes,
		Statistics:  stats,
		GeneratedAt: time.Now().UnixMilli(),
	}

	return result, nil
}

// GetDepartmentTree 获取部门树可视化
func (s *VisualizationService) GetDepartmentTree(ctx context.Context, orgID string, config *entity.VisualizationConfig) (*entity.VisualizationResult, error) {
	// 1. 获取部门树
	departments, err := s.visRepo.GetDepartmentTreeForVisualization(ctx, orgID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get department tree failed: org_id=%s", orgID)
	}

	// 2. 批量获取员工数
	deptIDs := make([]string, 0, len(departments))
	for _, dept := range departments {
		deptIDs = append(deptIDs, dept.DeptID)
	}

	empCounts := make(map[string]int64)
	for _, deptID := range deptIDs {
		count, err := s.visRepo.GetEmployeeCountByDept(ctx, deptID)
		if err != nil {
			continue
		}
		empCounts[deptID] = count
	}

	// 3. 构建可视化节点
	nodes := s.buildDepartmentNodes(departments, empCounts, config)

	// 4. 构建统计信息
	stats := s.buildVisualizationStatistics(nodes)

	// 5. 返回可视化结果
	result := &entity.VisualizationResult{
		Type:        entity.VisualizationTypeDepartment,
		Config:      *config,
		Nodes:       nodes,
		Statistics:  stats,
		GeneratedAt: time.Now().UnixMilli(),
	}

	return result, nil
}

// GetRelationshipNetwork 获取关系网络图
func (s *VisualizationService) GetRelationshipNetwork(ctx context.Context, req *GetVisualizationRequest) (*entity.VisualizationResult, error) {
	// 1. 获取组织树
	orgs, err := s.visRepo.GetOrganizationTreeForVisualization(ctx, req.TenantID)
	if err != nil {
		return nil, errorx.Wrapf(err, "get organization tree for network failed: tenant_id=%s", req.TenantID)
	}

	// 2. 构建节点和链接
	nodes := make([]*entity.OrganizationNode, 0)
	links := make([]*entity.RelationshipLink, 0)

	for _, org := range orgs {
		node := &entity.OrganizationNode{
			ID:   org.OrgID,
			Name: org.OrgName,
			Code: org.OrgCode,
			Type: org.OrgType,
			Metadata: entity.NodeMetadata{
				Color: s.getColorByType(org.OrgType),
				Icon:  s.getIconByType(org.OrgType),
			},
		}
		nodes = append(nodes, node)

		// 添加父子关系链接
		if org.HasParent() && org.ParentID != nil {
			link := &entity.RelationshipLink{
				Source: *org.ParentID,
				Target: org.OrgID,
				Type:   "parent-child",
				Weight: 1.0,
				Label:  "上级-下级",
			}
			links = append(links, link)
		}
	}

	// 3. 构建统计信息
	stats := entity.VisualizationStatistics{
		TotalNodes:    len(nodes),
		TotalLinks:    len(links),
		NodeTypeStats: s.countNodesByType(nodes),
	}

	result := &entity.VisualizationResult{
		Type:        entity.VisualizationTypeNetwork,
		Config:      req.Config,
		Nodes:       nodes,
		Links:       links,
		Statistics:  stats,
		GeneratedAt: time.Now().UnixMilli(),
	}

	return result, nil
}

// buildOrganizationNodes 构建组织节点
func (s *VisualizationService) buildOrganizationNodes(
	orgs []*entity.Organization,
	empCounts map[string]int64,
	config *entity.VisualizationConfig,
) []*entity.OrganizationNode {
	nodes := make([]*entity.OrganizationNode, 0, len(orgs))

	for _, org := range orgs {
		node := &entity.OrganizationNode{
			ID:       org.OrgID,
			Name:     org.OrgName,
			Code:     org.OrgCode,
			Type:     org.OrgType,
			Level:    org.Level,
			ParentID: org.ParentID,
			Metadata: entity.NodeMetadata{
				Color: s.getColorByType(org.OrgType),
				Icon:  s.getIconByType(org.OrgType),
				Size:  s.getSizeByLevel(org.Level),
			},
		}

		// 设置员工数
		if count, ok := empCounts[org.OrgID]; ok {
			node.EmployeeCount = int(count)
		}

		// 设置负责人
		if config.ShowLeader && org.LeaderID != nil {
			if leader, err := s.empRepo.GetByID(context.Background(), *org.LeaderID); err == nil {
				leaderEmail := ""
				if leader.Email != nil {
					leaderEmail = *leader.Email
				}
				leaderPhone := ""
				if leader.Phone != nil {
					leaderPhone = *leader.Phone
				}
				node.Leader = &entity.EmployeeSummary{
					EmpID:  leader.EmpID,
					Name:   leader.EmpName,
					Email:  leaderEmail,
					Phone:  leaderPhone,
					Avatar: "", // Avatar field not in Employee entity
				}
			}
		}

		// 递归构建子节点
		if org.Children != nil && len(org.Children) > 0 {
			node.Children = s.buildOrganizationNodes(org.Children, empCounts, config)
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// buildDepartmentNodes 构建部门节点
func (s *VisualizationService) buildDepartmentNodes(
	departments []*entity.Department,
	empCounts map[string]int64,
	config *entity.VisualizationConfig,
) []*entity.OrganizationNode {
	nodes := make([]*entity.OrganizationNode, 0, len(departments))

	for _, dept := range departments {
		node := &entity.OrganizationNode{
			ID:      dept.DeptID,
			Name:    dept.DeptName,
			Code:    dept.DeptCode,
			Type:    entity.OrgTypeDepartment,
			Level:   dept.Level,
			Metadata: entity.NodeMetadata{
				Color: "#3498db",
				Icon:  "department",
				Size:  20 + dept.Level*5,
			},
		}

		// 设置员工数
		if count, ok := empCounts[dept.DeptID]; ok {
			node.EmployeeCount = int(count)
		}

		// 递归构建子节点
		if dept.Children != nil && len(dept.Children) > 0 {
			node.Children = s.buildDepartmentNodes(dept.Children, empCounts, config)
		}

		nodes = append(nodes, node)
	}

	return nodes
}

// buildVisualizationStatistics 构建可视化统计信息
func (s *VisualizationService) buildVisualizationStatistics(nodes []*entity.OrganizationNode) entity.VisualizationStatistics {
	stats := entity.VisualizationStatistics{
		TotalNodes:    len(nodes),
		NodeTypeStats: make(map[string]int),
		LevelStats:    make(map[int]int),
	}

	maxDepth := 0
	for _, node := range nodes {
		// 统计类型
		stats.NodeTypeStats[string(node.Type)]++

		// 统计层级
		stats.LevelStats[node.Level]++

		// 计算最大深度
		if depth := s.calculateNodeDepth(node); depth > maxDepth {
			maxDepth = depth
		}
	}

	stats.MaxDepth = maxDepth

	return stats
}

// calculateNodeDepth 计算节点深度
func (s *VisualizationService) calculateNodeDepth(node *entity.OrganizationNode) int {
	if node.Children == nil || len(node.Children) == 0 {
		return 1
	}

	maxChildDepth := 0
	for _, child := range node.Children {
		if depth := s.calculateNodeDepth(child); depth > maxChildDepth {
			maxChildDepth = depth
		}
	}

	return maxChildDepth + 1
}

// countNodesByType 按类型统计节点
func (s *VisualizationService) countNodesByType(nodes []*entity.OrganizationNode) map[string]int {
	stats := make(map[string]int)
	for _, node := range nodes {
		stats[string(node.Type)]++
	}
	return stats
}

// getColorByType 根据类型获取颜色
func (s *VisualizationService) getColorByType(orgType entity.OrganizationType) string {
	colors := map[entity.OrganizationType]string{
		entity.OrgTypeCompany:   "#e74c3c",
		entity.OrgTypeDivision:  "#f39c12",
		entity.OrgTypeDepartment: "#3498db",
		entity.OrgTypeProject:    "#2ecc71",
	}

	if color, ok := colors[orgType]; ok {
		return color
	}
	return "#95a5a6"
}

// getIconByType 根据类型获取图标
func (s *VisualizationService) getIconByType(orgType entity.OrganizationType) string {
	icons := map[entity.OrganizationType]string{
		entity.OrgTypeCompany:   "company",
		entity.OrgTypeDivision:  "division",
		entity.OrgTypeDepartment: "department",
		entity.OrgTypeProject:    "project",
	}

	if icon, ok := icons[orgType]; ok {
		return icon
	}
	return "organization"
}

// getSizeByLevel 根据层级获取大小
func (s *VisualizationService) getSizeByLevel(level int) int {
	// 基础大小30，每增加一级增加10
	return 30 + level*10
}

// CompareOrganizations 对比组织
func (s *VisualizationService) CompareOrganizations(ctx context.Context, orgIDs []string, dimensions []entity.ComparisonDimension) (*entity.OrganizationComparison, error) {
	if len(orgIDs) < 2 {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取组织信息
	orgs, err := s.visRepo.GetOrganizationForComparison(ctx, orgIDs)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.ErrOrgNotFound.Message()+", %v",
			errorx.KV("operation", "get organizations for comparison"),
		)
	}

	// 2. 获取统计数据
	stats, err := s.visRepo.GetOrganizationStatsForComparison(ctx, orgIDs)
	if err != nil {
		return nil, errorx.Wrapf(err, errno.ErrOrgUpdateFailed.Message()+", %v",
			errorx.KV("operation", "get organization stats"),
		)
	}

	// 3. 构建对比结果
	results := make(map[string]entity.DimensionResult)

	for _, dimension := range dimensions {
		result := entity.DimensionResult{
			Dimension: dimension,
			Values:    make(map[string]interface{}),
		}

		switch dimension {
		case entity.DimensionSize:
			for orgID := range stats {
				if totalEmps, ok := stat["total_employees"].(int64); ok {
					result.Values[orgID] = totalEmps
				}
			}
			result.Unit = "人"

		case entity.DimensionCost:
			for orgID := range stats {
				// 成本计算逻辑（此处简化）
				result.Values[orgID] = 0
			}
			result.Unit = "元"

		case entity.DimensionUsage:
			for orgID := range stats {
				// 使用率计算逻辑（此处简化）
				result.Values[orgID] = 0.0
			}
			result.Unit = "%"

		case entity.DimensionActivity:
			for orgID := range orgs {
				activityStats, err := s.statsRepo.GetActivityStats(ctx, orgID, 30)
				if err == nil {
					result.Values[orgID] = activityStats.AvgActivity
				}
			}
			result.Unit = "分"

		case entity.DimensionStructure:
			for orgID, org := range orgs {
				result.Values[orgID] = map[string]interface{}{
					"level": org.Level,
					"type":  org.OrgType,
				}
			}
		}

		results[string(dimension)] = result
	}

	comparison := &entity.OrganizationComparison{
		OrgIDs:      orgIDs,
		Dimensions:  dimensions,
		Results:     results,
		GeneratedAt: time.Now().UnixMilli(),
	}

	return comparison, nil
}

// ExportVisualization 导出可视化数据
func (s *VisualizationService) ExportVisualization(ctx context.Context, result *entity.VisualizationResult, format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(result)
	case "svg":
		return s.exportAsSVG(result)
	default:
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}
}

// exportAsSVG 导出为SVG格式（简化实现）
func (s *VisualizationService) exportAsSVG(result *entity.VisualizationResult) ([]byte, error) {
	svg := fmt.Sprintf(`<svg width="800" height="600" xmlns="http://www.w3.org/2000/svg">
		<!-- 组织可视化图 -->
		<!-- 节点数: %d -->
		<!-- 类型: %s -->
	</svg>`, result.Statistics.TotalNodes, result.Type)

	return []byte(svg), nil
}
