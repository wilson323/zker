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

package repository

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
)

// VisualizationRepository 可视化仓储接口
type VisualizationRepository interface {
	// GetOrganizationTreeForVisualization 获取组织树用于可视化
	GetOrganizationTreeForVisualization(ctx context.Context, tenantID string) ([]*entity.Organization, error)

	// GetDepartmentTreeForVisualization 获取部门树用于可视化
	GetDepartmentTreeForVisualization(ctx context.Context, orgID string) ([]*entity.Department, error)

	// GetEmployeeCountByOrg 获取组织的员工数
	GetEmployeeCountByOrg(ctx context.Context, orgID string) (int64, error)

	// GetEmployeeCountByDept 获取部门的员工数
	GetEmployeeCountByDept(ctx context.Context, deptID string) (int64, error)

	// BatchGetEmployeeCount 批量获取员工数
	BatchGetEmployeeCount(ctx context.Context, orgIDs []string) (map[string]int64, error)

	// GetOrganizationForComparison 获取组织用于对比
	GetOrganizationForComparison(ctx context.Context, orgIDs []string) (map[string]*entity.Organization, error)

	// GetOrganizationStatsForComparison 获取组织统计用于对比
	GetOrganizationStatsForComparison(ctx context.Context, orgIDs []string) (map[string]map[string]interface{}, error)
}

// DataPermissionRepository 数据权限仓储接口
type DataPermissionRepository interface {
	// GetUserPermissionLevel 获取用户的数据权限级别
	GetUserPermissionLevel(ctx context.Context, userID, resourceType string) (int, error)

	// GetAccessibleDepartments 获取用户可访问的部门列表
	GetAccessibleDepartments(ctx context.Context, userID string) ([]string, error)

	// GetAccessibleOrganizations 获取用户可访问的组织列表
	GetAccessibleOrganizations(ctx context.Context, userID string) ([]string, error)

	// CheckDataPermission 检查数据权限
	CheckDataPermission(ctx context.Context, userID, resourceType, resourceID string) (bool, error)
}

// FieldPermissionRepository 字段权限仓储接口
type FieldPermissionRepository interface {
	// GetFieldPermissions 获取字段权限配置
	GetFieldPermissions(ctx context.Context, roleID, resourceType string) ([]*entity.FieldPermission, error)

	// GetUserFieldPermissions 获取用户的字段权限
	GetUserFieldPermissions(ctx context.Context, userID, resourceType string) ([]*entity.FieldPermission, error)

	// FilterFieldsByPermission 根据权限过滤字段
	FilterFieldsByPermission(ctx context.Context, userID, resourceType string, data map[string]interface{}) (map[string]interface{}, error)

	// BatchGetFieldPermissions 批量获取字段权限
	BatchGetFieldPermissions(ctx context.Context, userIDs []string, resourceType string) (map[string][]*entity.FieldPermission, error)
}

// OrganizationTemplateRepository 组织模板仓储接口
type OrganizationTemplateRepository interface {
	// Create 创建组织模板
	Create(ctx context.Context, template *entity.OrganizationTemplate) error

	// GetByID 根据ID获取模板
	GetByID(ctx context.Context, templateID string) (*entity.OrganizationTemplate, error)

	// GetByCategory 获取分类下的模板
	GetByCategory(ctx context.Context, category string) ([]*entity.OrganizationTemplate, error)

	// GetPredefinedTemplates 获取预置模板
	GetPredefinedTemplates(ctx context.Context) ([]*entity.OrganizationTemplate, error)

	// List 列出模板
	List(ctx context.Context, filter *TemplateFilter) ([]*entity.OrganizationTemplate, int64, error)

	// Update 更新模板
	Update(ctx context.Context, template *entity.OrganizationTemplate) error

	// Delete 软删除模板
	Delete(ctx context.Context, templateID string) error

	// IncrementUsage 增加使用次数
	IncrementUsage(ctx context.Context, templateID string) error
}

// TemplateFilter 模板查询过滤器
type TemplateFilter struct {
	TenantID  string
	Category  string
	IsPredefined bool
	Keyword   string
	PageToken string
	PageSize  int
}

// OrganizationMigrationRepository 组织迁移仓储接口
type OrganizationMigrationRepository interface {
	// Create 创建迁移任务
	Create(ctx context.Context, migration *entity.OrganizationMigration) error

	// GetByID 根据ID获取迁移任务
	GetByID(ctx context.Context, migrationID string) (*entity.OrganizationMigration, error)

	// GetByTenantID 获取租户的迁移任务列表
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.OrganizationMigration, error)

	// GetActiveMigration 获取进行中的迁移任务
	GetActiveMigration(ctx context.Context, sourceTenantID, targetTenantID string) (*entity.OrganizationMigration, error)

	// Update 更新迁移任务
	Update(ctx context.Context, migration *entity.OrganizationMigration) error

	// UpdateProgress 更新迁移进度
	UpdateProgress(ctx context.Context, migrationID string, progress int, status entity.MigrationStatus) error

	// List 列出迁移任务
	List(ctx context.Context, filter *MigrationFilter) ([]*entity.OrganizationMigration, int64, error)
}

// MigrationFilter 迁移任务查询过滤器
type MigrationFilter struct {
	TenantID      string
	SourceTenantID string
	TargetTenantID string
	Status        entity.MigrationStatus
	StartTime     int64
	EndTime       int64
	PageToken     string
	PageSize      int
}

// OrganizationStatsRepository 组织统计仓储接口
type OrganizationStatsRepository interface {
	// GetOrganizationStats 获取组织统计
	GetOrganizationStats(ctx context.Context, orgID string) (*entity.OrganizationStats, error)

	// GetDepartmentStats 获取部门统计
	GetDepartmentStats(ctx context.Context, deptID string) (*entity.DepartmentStats, error)

	// BatchGetOrganizationStats 批量获取组织统计
	BatchGetOrganizationStats(ctx context.Context, orgIDs []string) (map[string]*entity.OrganizationStats, error)

	// GetTenantStats 获取租户统计
	GetTenantStats(ctx context.Context, tenantID string) (*entity.TenantStats, error)

	// GetActivityStats 获取活跃度统计
	GetActivityStats(ctx context.Context, orgID string, days int) (*entity.ActivityStats, error)
}

// AddressBookRepository 通讯录仓储接口
type AddressBookRepository interface {
	// SearchEmployees 搜索员工（用于通讯录）
	SearchEmployees(ctx context.Context, req *AddressBookSearchRequest) ([]*entity.EmployeeSummary, int64, error)

	// GetDepartmentEmployees 获取部门员工（用于通讯录）
	GetDepartmentEmployees(ctx context.Context, deptID string) ([]*entity.EmployeeSummary, error)

	// GetOrganizationEmployees 获取组织员工（用于通讯录）
	GetOrganizationEmployees(ctx context.Context, orgID string) ([]*entity.EmployeeSummary, error)

	// GetAllEmployees 获取所有员工（用于通讯录）
	GetAllEmployees(ctx context.Context, tenantID string) ([]*entity.EmployeeSummary, error)
}

// AddressBookSearchRequest 通讯录搜索请求
type AddressBookSearchRequest struct {
	TenantID   string
	OrgID      string
	DeptID     string
	Keyword    string // 搜索关键词（姓名、手机号、邮箱）
	PositionID string
	PageToken  string
	PageSize   int
}
