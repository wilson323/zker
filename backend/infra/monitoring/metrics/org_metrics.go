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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ========== 组织管理指标 ==========

var (
	// OrganizationTotal 组织总数(Gauge)
	OrganizationTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "organization_total",
			Help: "Total number of organizations",
		},
		[]string{"tenant_id", "org_type", "status"}, // org_type: company, division, department, project; status: active, inactive, frozen
	)

	// OrganizationCreationTotal 组织创建总数(Counter)
	OrganizationCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "organization_creation_total",
			Help: "Total number of organization creations",
		},
		[]string{"tenant_id", "org_type", "result"}, // result: success, failure
	)

	// OrganizationDeletionTotal 组织删除总数(Counter)
	OrganizationDeletionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "organization_deletion_total",
			Help: "Total number of organization deletions",
		},
		[]string{"tenant_id", "org_type", "result"},
	)

	// OrganizationUpdateTotal 组织更新总数(Counter)
	OrganizationUpdateTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "organization_update_total",
			Help: "Total number of organization updates",
		},
		[]string{"tenant_id", "org_type", "update_type", "result"}, // update_type: info, leader, status, move
	)

	// OrganizationMoveTotal 组织移动总数(Counter)
	OrganizationMoveTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "organization_move_total",
			Help: "Total number of organization moves",
		},
		[]string{"tenant_id", "org_type", "result"},
	)

	// OrganizationDepthMax 组织最大层级深度(Gauge)
	OrganizationDepthMax = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "organization_depth_max",
			Help: "Maximum organization hierarchy depth",
		},
		[]string{"tenant_id"},
	)

	// OrganizationChildrenTotal 组织子节点数分布(Histogram)
	OrganizationChildrenTotal = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "organization_children_count",
			Help:    "Number of children per organization",
			Buckets: []float64{0, 1, 2, 5, 10, 20, 50, 100},
		},
		[]string{"tenant_id", "org_type"},
	)

	// OrganizationOperationDuration 组织操作延迟(Histogram)
	OrganizationOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "organization_operation_duration_seconds",
			Help:    "Organization operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "operation_type"}, // operation_type: create, update, delete, move, query
	)
)

// ========== 部门管理指标 ==========

var (
	// DepartmentTotal 部门总数(Gauge)
	DepartmentTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "department_total",
			Help: "Total number of departments",
		},
		[]string{"tenant_id", "org_id", "status"}, // status: active, inactive, frozen
	)

	// DepartmentCreationTotal 部门创建总数(Counter)
	DepartmentCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "department_creation_total",
			Help: "Total number of department creations",
		},
		[]string{"tenant_id", "org_id", "result"},
	)

	// DepartmentDeletionTotal 部门删除总数(Counter)
	DepartmentDeletionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "department_deletion_total",
			Help: "Total number of department deletions",
		},
		[]string{"tenant_id", "org_id", "result"},
	)

	// DepartmentMemberTotal 部门成员总数(Gauge)
	DepartmentMemberTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "department_member_total",
			Help: "Total number of department members",
		},
		[]string{"tenant_id", "dept_id"},
	)

	// DepartmentMemberCount 部门成员数分布(Histogram)
	DepartmentMemberCount = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "department_member_count",
			Help:    "Number of members per department",
			Buckets: []float64{0, 1, 5, 10, 20, 50, 100, 200, 500},
		},
		[]string{"tenant_id", "org_id"},
	)

	// DepartmentDepthMax 部门最大层级深度(Gauge)
	DepartmentDepthMax = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "department_depth_max",
			Help: "Maximum department hierarchy depth",
		},
		[]string{"tenant_id", "org_id"},
	)

	// DepartmentOperationDuration 部门操作延迟(Histogram)
	DepartmentOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "department_operation_duration_seconds",
			Help:    "Department operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "operation_type"},
	)
)

// ========== 员工管理指标 ==========

var (
	// EmployeeTotal 员工总数(Gauge)
	EmployeeTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "employee_total",
			Help: "Total number of employees",
		},
		[]string{"tenant_id", "org_id", "emp_status"}, // emp_status: active, inactive, resigned, locked
	)

	// EmployeeCreationTotal 员工创建总数(Counter)
	EmployeeCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_creation_total",
			Help: "Total number of employee creations",
		},
		[]string{"tenant_id", "org_id", "result"},
	)

	// EmployeeDeletionTotal 员工删除总数(Counter)
	EmployeeDeletionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_deletion_total",
			Help: "Total number of employee deletions",
		},
		[]string{"tenant_id", "org_id", "result"},
	)

	// EmployeeUpdateTotal 员工更新总数(Counter)
	EmployeeUpdateTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_update_total",
			Help: "Total number of employee updates",
		},
		[]string{"tenant_id", "update_type", "result"}, // update_type: info, status, transfer, promotion
	)

	// EmployeeTransferTotal 员工调转总数(Counter)
	EmployeeTransferTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_transfer_total",
			Help: "Total number of employee transfers",
		},
		[]string{"tenant_id", "transfer_type", "result"}, // transfer_type: dept, position
	)

	// EmployeePromotionTotal 员工晋升总数(Counter)
	EmployeePromotionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_promotion_total",
			Help: "Total number of employee promotions",
		},
		[]string{"tenant_id", "result"},
	)

	// EmployeeResignationTotal 员工离职总数(Counter)
	EmployeeResignationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_resignation_total",
			Help: "Total number of employee resignations",
		},
		[]string{"tenant_id", "org_id", "reason_type"}, // reason_type: voluntary, involuntary, retirement
	)

	// EmployeeOnboardingTotal 员工入职总数(Counter)
	EmployeeOnboardingTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_onboarding_total",
			Help: "Total number of employee onboardings",
		},
		[]string{"tenant_id", "org_id", "result"},
	)

	// EmployeeActiveDays 员工在职天数(Histogram)
	EmployeeActiveDays = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "employee_active_days",
			Help:    "Employee active days before resignation",
			Buckets: []float64{1, 30, 90, 180, 365, 730, 1095, 1825}, // 1天, 1月, 3月, 6月, 1年, 2年, 3年, 5年
		},
		[]string{"tenant_id", "org_id"},
	)

	// EmployeeOperationDuration 员工操作延迟(Histogram)
	EmployeeOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "employee_operation_duration_seconds",
			Help:    "Employee operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "operation_type"},
	)

	// EmployeeLoginTotal 员工登录总数(Counter)
	EmployeeLoginTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "employee_login_total",
			Help: "Total number of employee logins",
		},
		[]string{"tenant_id", "status"}, // status: success, failed
	)

	// EmployeeActiveDaily 每日活跃员工数(Gauge)
	EmployeeActiveDaily = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "employee_active_daily",
			Help: "Number of daily active employees",
		},
		[]string{"tenant_id", "org_id"},
	)

	// EmployeeActiveWeekly 每周活跃员工数(Gauge)
	EmployeeActiveWeekly = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "employee_active_weekly",
			Help: "Number of weekly active employees",
		},
		[]string{"tenant_id", "org_id"},
	)

	// EmployeeActiveMonthly 每月活跃员工数(Gauge)
	EmployeeActiveMonthly = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "employee_active_monthly",
			Help: "Number of monthly active employees",
		},
		[]string{"tenant_id", "org_id"},
	)
)

// ========== 岗位管理指标 ==========

var (
	// PositionTotal 岗位总数(Gauge)
	PositionTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "position_total",
			Help: "Total number of positions",
		},
		[]string{"tenant_id", "dept_id", "position_level", "status"}, // position_level: 1,2,3; status: active, inactive, frozen
	)

	// PositionCreationTotal 岗位创建总数(Counter)
	PositionCreationTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "position_creation_total",
			Help: "Total number of position creations",
		},
		[]string{"tenant_id", "position_level", "result"},
	)

	// PositionDeletionTotal 岗位删除总数(Counter)
	PositionDeletionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "position_deletion_total",
			Help: "Total number of position deletions",
		},
		[]string{"tenant_id", "position_level", "result"},
	)

	// PositionOccupiedTotal 岗位占用总数(Gauge)
	PositionOccupiedTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "position_occupied_total",
			Help: "Total number of occupied positions",
		},
		[]string{"tenant_id", "position_level"},
	)

	// PositionVacantTotal 岗位空缺总数(Gauge)
	PositionVacantTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "position_vacant_total",
			Help: "Total number of vacant positions",
		},
		[]string{"tenant_id", "position_level"},
	)

	// PositionFillRate 岗位填充率(Gauge)
	PositionFillRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "position_fill_rate",
			Help: "Position fill rate (percentage)",
		},
		[]string{"tenant_id", "position_level"},
	)

	// PositionOperationDuration 岗位操作延迟(Histogram)
	PositionOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "position_operation_duration_seconds",
			Help:    "Position operation latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "operation_type"},
	)
)

// ========== 组织架构指标 ==========

var (
	// OrgHierarchyDepth 组织架构层级深度(Gauge)
	OrgHierarchyDepth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "org_hierarchy_depth",
			Help: "Organization hierarchy depth",
		},
		[]string{"tenant_id", "hierarchy_type"}, // hierarchy_type: org, dept
	)

	// OrgTreeHealth 组织树健康度(Gauge)
	OrgTreeHealth = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "org_tree_health",
			Help: "Organization tree health score (0-100)",
		},
		[]string{"tenant_id", "health_dimension"}, // health_dimension: balance, depth, orphan
	)

	// OrgOrphanTotal 孤立组织节点数(Gauge)
	OrgOrphanTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "org_orphan_total",
			Help: "Total number of orphan organization nodes",
		},
		[]string{"tenant_id", "node_type"}, // node_type: org, dept
	)

	// OrgQueryDuration 组织查询延迟(Histogram)
	OrgQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "org_query_duration_seconds",
			Help:    "Organization query latency in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1},
		},
		[]string{"tenant_id", "query_type"}, // query_type: tree, path, ancestor, descendant
	)

	// OrgTraversalPerformance 组织遍历性能(Histogram)
	OrgTraversalPerformance = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "org_traversal_duration_seconds",
			Help:    "Organization tree traversal latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "traversal_type"}, // traversal_type: dfs, bfs, path
	)
)

// ========== HR生命周期指标 ==========

var (
	// HROnboardingDuration 入职办理延迟(Histogram)
	HROnboardingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hr_onboarding_duration_seconds",
			Help:    "HR onboarding process duration in seconds",
			Buckets: []float64{60, 300, 600, 1800, 3600}, // 1分钟, 5分钟, 10分钟, 30分钟, 1小时
		},
		[]string{"tenant_id", "step"}, // step: create, approve, notify
	)

	// HRResignationDuration 离职办理延迟(Histogram)
	HRResignationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hr_resignation_duration_seconds",
			Help:    "HR resignation process duration in seconds",
			Buckets: []float64{60, 300, 600, 1800, 3600},
		},
		[]string{"tenant_id", "step"},
	)

	// HRTransferDuration 调转办理延迟(Histogram)
	HRTransferDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hr_transfer_duration_seconds",
			Help:    "HR transfer process duration in seconds",
			Buckets: []float64{60, 300, 600, 1800, 3600},
		},
		[]string{"tenant_id", "transfer_type"},
	)

	// HRProcessTotal HR流程总数(Counter)
	HRProcessTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hr_process_total",
			Help: "Total number of HR processes",
		},
		[]string{"tenant_id", "process_type", "result"}, // process_type: onboarding, resignation, transfer, promotion
	)

	// HRApprovalDuration HR审批延迟(Histogram)
	HRApprovalDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hr_approval_duration_seconds",
			Help:    "HR approval latency in seconds",
			Buckets: []float64{60, 300, 600, 1800, 3600, 86400}, // 最长24小时
		},
		[]string{"tenant_id", "approval_type"},
	)
)

// ========== 组织目录指标 ==========

var (
	// DirectoryQueryDuration 目录查询延迟(Histogram)
	DirectoryQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "directory_query_duration_seconds",
			Help:    "Directory query latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "query_type"}, // query_type: org, dept, employee, search
	)

	// DirectoryCacheHitRate 目录缓存命中率(Gauge)
	DirectoryCacheHitRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "directory_cache_hit_rate",
			Help: "Directory cache hit rate (percentage)",
		},
		[]string{"tenant_id", "cache_type"}, // cache_type: org_tree, employee_list
	)

	// DirectorySearchDuration 目录搜索延迟(Histogram)
	DirectorySearchDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "directory_search_duration_seconds",
			Help:    "Directory search latency in seconds",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"tenant_id", "search_type"}, // search_type: name, code, phone, email
	)

	// DirectoryIndexSize 目录索引大小(Gauge)
	DirectoryIndexSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "directory_index_size_bytes",
			Help: "Directory index size in bytes",
		},
		[]string{"tenant_id", "index_type"},
	)
)

// ========== 辅助函数 ==========

// RecordOrganizationCreation 记录组织创建
func RecordOrganizationCreation(tenantID, orgType string, result string, duration float64) {
	OrganizationCreationTotal.WithLabelValues(tenantID, orgType, result).Inc()
	OrganizationOperationDuration.WithLabelValues(tenantID, "create").Observe(duration)
}

// RecordOrganizationDeletion 记录组织删除
func RecordOrganizationDeletion(tenantID, orgType string, result string, duration float64) {
	OrganizationDeletionTotal.WithLabelValues(tenantID, orgType, result).Inc()
	OrganizationOperationDuration.WithLabelValues(tenantID, "delete").Observe(duration)
}

// RecordOrganizationUpdate 记录组织更新
func RecordOrganizationUpdate(tenantID, orgType, updateType string, result string, duration float64) {
	OrganizationUpdateTotal.WithLabelValues(tenantID, orgType, updateType, result).Inc()
	OrganizationOperationDuration.WithLabelValues(tenantID, "update").Observe(duration)
}

// RecordOrganizationMove 记录组织移动
func RecordOrganizationMove(tenantID, orgType string, result string, duration float64) {
	OrganizationMoveTotal.WithLabelValues(tenantID, orgType, result).Inc()
	OrganizationOperationDuration.WithLabelValues(tenantID, "move").Observe(duration)
}

// UpdateOrganizationMetrics 更新组织指标
func UpdateOrganizationMetrics(tenantID string, orgType, status string, count int, depth int, childrenCount int) {
	OrganizationTotal.WithLabelValues(tenantID, orgType, status).Set(float64(count))
	OrganizationDepthMax.WithLabelValues(tenantID).Set(float64(depth))
	OrganizationChildrenTotal.WithLabelValues(tenantID, orgType).Observe(float64(childrenCount))
}

// RecordDepartmentCreation 记录部门创建
func RecordDepartmentCreation(tenantID, orgID string, result string, duration float64) {
	DepartmentCreationTotal.WithLabelValues(tenantID, orgID, result).Inc()
	DepartmentOperationDuration.WithLabelValues(tenantID, "create").Observe(duration)
}

// RecordDepartmentDeletion 记录部门删除
func RecordDepartmentDeletion(tenantID, orgID string, result string, duration float64) {
	DepartmentDeletionTotal.WithLabelValues(tenantID, orgID, result).Inc()
	DepartmentOperationDuration.WithLabelValues(tenantID, "delete").Observe(duration)
}

// UpdateDepartmentMetrics 更新部门指标
func UpdateDepartmentMetrics(tenantID, orgID string, status string, count int, memberCount int, depth int) {
	DepartmentTotal.WithLabelValues(tenantID, orgID, status).Set(float64(count))
	DepartmentMemberTotal.WithLabelValues(tenantID, "").Set(float64(memberCount))
	DepartmentMemberCount.WithLabelValues(tenantID, orgID).Observe(float64(memberCount))
	DepartmentDepthMax.WithLabelValues(tenantID, orgID).Set(float64(depth))
}

// RecordEmployeeCreation 记录员工创建
func RecordEmployeeCreation(tenantID, orgID string, result string, duration float64) {
	EmployeeCreationTotal.WithLabelValues(tenantID, orgID, result).Inc()
	EmployeeOperationDuration.WithLabelValues(tenantID, "create").Observe(duration)
}

// RecordEmployeeDeletion 记录员工删除
func RecordEmployeeDeletion(tenantID, orgID string, result string, duration float64) {
	EmployeeDeletionTotal.WithLabelValues(tenantID, orgID, result).Inc()
	EmployeeOperationDuration.WithLabelValues(tenantID, "delete").Observe(duration)
}

// RecordEmployeeTransfer 记录员工调转
func RecordEmployeeTransfer(tenantID, transferType string, result string, duration float64) {
	EmployeeTransferTotal.WithLabelValues(tenantID, transferType, result).Inc()
	EmployeeOperationDuration.WithLabelValues(tenantID, "transfer").Observe(duration)
}

// RecordEmployeePromotion 记录员工晋升
func RecordEmployeePromotion(tenantID string, result string, duration float64) {
	EmployeePromotionTotal.WithLabelValues(tenantID, result).Inc()
	EmployeeOperationDuration.WithLabelValues(tenantID, "promotion").Observe(duration)
}

// RecordEmployeeResignation 记录员工离职
func RecordEmployeeResignation(tenantID, orgID, reasonType string, activeDays int) {
	EmployeeResignationTotal.WithLabelValues(tenantID, orgID, reasonType).Inc()
	EmployeeActiveDays.WithLabelValues(tenantID, orgID).Observe(float64(activeDays))
}

// RecordEmployeeOnboarding 记录员工入职
func RecordEmployeeOnboarding(tenantID, orgID string, result string, duration float64) {
	EmployeeOnboardingTotal.WithLabelValues(tenantID, orgID, result).Inc()
	EmployeeOperationDuration.WithLabelValues(tenantID, "onboarding").Observe(duration)
}

// UpdateEmployeeMetrics 更新员工指标
func UpdateEmployeeMetrics(tenantID, orgID, empStatus string, count int) {
	EmployeeTotal.WithLabelValues(tenantID, orgID, empStatus).Set(float64(count))
}

// RecordEmployeeLogin 记录员工登录
func RecordEmployeeLogin(tenantID string, status string) {
	EmployeeLoginTotal.WithLabelValues(tenantID, status).Inc()
}

// UpdateEmployeeActiveCount 更新活跃员工数
func UpdateEmployeeActiveCount(tenantID, orgID string, period string, count int) {
	switch period {
	case "daily":
		EmployeeActiveDaily.WithLabelValues(tenantID, orgID).Set(float64(count))
	case "weekly":
		EmployeeActiveWeekly.WithLabelValues(tenantID, orgID).Set(float64(count))
	case "monthly":
		EmployeeActiveMonthly.WithLabelValues(tenantID, orgID).Set(float64(count))
	}
}

// RecordPositionCreation 记录岗位创建
func RecordPositionCreation(tenantID string, positionLevel string, result string, duration float64) {
	PositionCreationTotal.WithLabelValues(tenantID, positionLevel, result).Inc()
	PositionOperationDuration.WithLabelValues(tenantID, "create").Observe(duration)
}

// RecordPositionDeletion 记录岗位删除
func RecordPositionDeletion(tenantID string, positionLevel string, result string, duration float64) {
	PositionDeletionTotal.WithLabelValues(tenantID, positionLevel, result).Inc()
	PositionOperationDuration.WithLabelValues(tenantID, "delete").Observe(duration)
}

// UpdatePositionMetrics 更新岗位指标
func UpdatePositionMetrics(tenantID, deptID string, positionLevel, status string, total int, occupied int, vacant int, fillRate float64) {
	PositionTotal.WithLabelValues(tenantID, deptID, positionLevel, status).Set(float64(total))
	PositionOccupiedTotal.WithLabelValues(tenantID, positionLevel).Set(float64(occupied))
	PositionVacantTotal.WithLabelValues(tenantID, positionLevel).Set(float64(vacant))
	PositionFillRate.WithLabelValues(tenantID, positionLevel).Set(fillRate)
}

// UpdateOrgTreeHealth 更新组织树健康度
func UpdateOrgTreeHealth(tenantID, healthDimension string, score float64) {
	OrgTreeHealth.WithLabelValues(tenantID, healthDimension).Set(score)
}

// UpdateOrgOrphanCount 更新孤立节点数
func UpdateOrgOrphanCount(tenantID, nodeType string, count int) {
	OrgOrphanTotal.WithLabelValues(tenantID, nodeType).Set(float64(count))
}

// RecordOrgQuery 记录组织查询
func RecordOrgQuery(tenantID, queryType string, duration float64) {
	OrgQueryDuration.WithLabelValues(tenantID, queryType).Observe(duration)
}

// RecordOrgTraversal 记录组织遍历
func RecordOrgTraversal(tenantID, traversalType string, duration float64) {
	OrgTraversalPerformance.WithLabelValues(tenantID, traversalType).Observe(duration)
}

// RecordHROnboarding 记录入职流程
func RecordHROnboarding(tenantID, step string, result string, duration float64) {
	HRProcessTotal.WithLabelValues(tenantID, "onboarding", result).Inc()
	HROnboardingDuration.WithLabelValues(tenantID, step).Observe(duration)
}

// RecordHRResignation 记录离职流程
func RecordHRResignation(tenantID, step string, result string, duration float64) {
	HRProcessTotal.WithLabelValues(tenantID, "resignation", result).Inc()
	HRResignationDuration.WithLabelValues(tenantID, step).Observe(duration)
}

// RecordHRTransfer 记录调转流程
func RecordHRTransfer(tenantID, transferType string, result string, duration float64) {
	HRProcessTotal.WithLabelValues(tenantID, "transfer", result).Inc()
	HRTransferDuration.WithLabelValues(tenantID, transferType).Observe(duration)
}

// RecordHRApproval 记录HR审批
func RecordHRApproval(tenantID, approvalType string, duration float64) {
	HRApprovalDuration.WithLabelValues(tenantID, approvalType).Observe(duration)
}

// RecordDirectoryQuery 记录目录查询
func RecordDirectoryQuery(tenantID, queryType string, duration float64) {
	DirectoryQueryDuration.WithLabelValues(tenantID, queryType).Observe(duration)
}

// RecordDirectorySearch 记录目录搜索
func RecordDirectorySearch(tenantID, searchType string, duration float64) {
	DirectorySearchDuration.WithLabelValues(tenantID, searchType).Observe(duration)
}

// UpdateDirectoryCacheHitRate 更新目录缓存命中率
func UpdateDirectoryCacheHitRate(tenantID, cacheType string, hitRate float64) {
	DirectoryCacheHitRate.WithLabelValues(tenantID, cacheType).Set(hitRate)
}

// UpdateDirectoryIndexSize 更新目录索引大小
func UpdateDirectoryIndexSize(tenantID, indexType string, sizeBytes int) {
	DirectoryIndexSize.WithLabelValues(tenantID, indexType).Set(float64(sizeBytes))
}
