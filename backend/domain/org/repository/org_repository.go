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

// OrganizationRepository 组织仓储接口
type OrganizationRepository interface {
	// Create 创建组织
	Create(ctx context.Context, org *entity.Organization) error

	// GetByID 根据ID获取组织
	GetByID(ctx context.Context, orgID string) (*entity.Organization, error)

	// GetByCode 根据组织编码获取组织
	GetByCode(ctx context.Context, tenantID, code string) (*entity.Organization, error)

	// GetByTenantID 获取租户下的所有组织
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Organization, error)

	// GetTree 获取组织树
	GetTree(ctx context.Context, tenantID string) ([]*entity.Organization, error)

	// GetChildren 获取子组织
	GetChildren(ctx context.Context, parentID string) ([]*entity.Organization, error)

	// Update 更新组织
	Update(ctx context.Context, org *entity.Organization) error

	// Delete 软删除组织
	Delete(ctx context.Context, orgID string) error

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error)

	// List 分页查询组织列表
	List(ctx context.Context, filter *OrganizationFilter) ([]*entity.Organization, int64, error)
}

// OrganizationFilter 组织查询过滤器
type OrganizationFilter struct {
	TenantID   string
	OrgType    entity.OrganizationType
	Status     entity.OrganizationStatus
	ParentID   *string
	Level      *int
	Keyword    string // 搜索关键词（名称或编码）
	PageToken  string
	PageSize   int
}

// DepartmentRepository 部门仓储接口
type DepartmentRepository interface {
	// Create 创建部门
	Create(ctx context.Context, dept *entity.Department) error

	// GetByID 根据ID获取部门
	GetByID(ctx context.Context, deptID string) (*entity.Department, error)

	// GetByCode 根据部门编码获取部门
	GetByCode(ctx context.Context, tenantID, code string) (*entity.Department, error)

	// GetByTenantID 获取租户下的所有部门
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Department, error)

	// GetTree 获取部门树
	GetTree(ctx context.Context, tenantID string) ([]*entity.Department, error)

	// GetChildren 获取子部门
	GetChildren(ctx context.Context, parentID string) ([]*entity.Department, error)

	// Update 更新部门
	Update(ctx context.Context, dept *entity.Department) error

	// Delete 软删除部门
	Delete(ctx context.Context, deptID string) error

	// Move 移动部门到新的父部门下
	Move(ctx context.Context, deptID, newParentID string) error

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error)

	// List 分页查询部门列表
	List(ctx context.Context, filter *DepartmentFilter) ([]*entity.Department, int64, error)
}

// DepartmentFilter 部门查询过滤器
type DepartmentFilter struct {
	TenantID   string
	OrgID      string
	Status     entity.OrganizationStatus
	ParentID   *string
	Level      *int
	Keyword    string
	PageToken  string
	PageSize   int
}

// PositionRepository 岗位仓储接口
type PositionRepository interface {
	// Create 创建岗位
	Create(ctx context.Context, pos *entity.Position) error

	// GetByID 根据ID获取岗位
	GetByID(ctx context.Context, positionID string) (*entity.Position, error)

	// GetByCode 根据岗位编码获取岗位
	GetByCode(ctx context.Context, tenantID, code string) (*entity.Position, error)

	// GetByTenantID 获取租户下的所有岗位
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Position, error)

	// GetByDepartmentID 获取部门下的岗位
	GetByDepartmentID(ctx context.Context, deptID string) ([]*entity.Position, error)

	// Update 更新岗位
	Update(ctx context.Context, pos *entity.Position) error

	// Delete 软删除岗位
	Delete(ctx context.Context, positionID string) error

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error)

	// List 分页查询岗位列表
	List(ctx context.Context, filter *PositionFilter) ([]*entity.Position, int64, error)
}

// PositionFilter 岗位查询过滤器
type PositionFilter struct {
	TenantID    string
	DeptID      *string
	Category    string
	Status      entity.OrganizationStatus
	Level       *int
	Keyword     string
	PageToken   string
	PageSize    int
}

// EmployeeRepository 员工仓储接口
type EmployeeRepository interface {
	// Create 创建员工
	Create(ctx context.Context, emp *entity.Employee) error

	// GetByID 根据ID获取员工
	GetByID(ctx context.Context, empID string) (*entity.Employee, error)

	// GetByCode 根据工号获取员工
	GetByCode(ctx context.Context, tenantID, code string) (*entity.Employee, error)

	// GetByUserID 根据用户ID获取员工
	GetByUserID(ctx context.Context, userID string) (*entity.Employee, error)

	// GetByTenantID 获取租户下的所有员工
	GetByTenantID(ctx context.Context, tenantID string) ([]*entity.Employee, error)

	// GetByDepartmentID 获取部门下的员工
	GetByDepartmentID(ctx context.Context, deptID string) ([]*entity.Employee, error)

	// GetByOrgID 获取组织下的员工
	GetByOrgID(ctx context.Context, orgID string) ([]*entity.Employee, error)

	// Update 更新员工
	Update(ctx context.Context, emp *entity.Employee) error

	// UpdateStatus 更新员工状态
	UpdateStatus(ctx context.Context, empID string, status entity.EmployeeStatus) error

	// Delete 软删除员工
	Delete(ctx context.Context, empID string) error

	// ExistsByCode 检查工号是否存在
	ExistsByCode(ctx context.Context, tenantID, code string, excludeID string) (bool, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, tenantID, email string, excludeID string) (bool, error)

	// List 分页查询员工列表
	List(ctx context.Context, filter *EmployeeFilter) ([]*entity.Employee, int64, error)

	// Search 搜索员工（按姓名、工号、手机号、邮箱）
	Search(ctx context.Context, tenantID, keyword string, limit int) ([]*entity.Employee, error)

	// GetByPinyin 按拼音首字母查询员工
	GetByPinyin(ctx context.Context, tenantID, pinyin string) ([]*entity.Employee, error)
}

// EmployeeFilter 员工查询过滤器
type EmployeeFilter struct {
	TenantID      string
	OrgID         string
	DeptID        *string
	PositionID    *string
	EmployeeType  entity.EmployeeType
	EmployeeStatus entity.EmployeeStatus
	JobLevel      *int
	Status        entity.EmployeeStatus
	Keyword       string // 搜索关键词（姓名、工号、手机号、邮箱）
	PageToken     string
	PageSize      int
}

// EmployeeContractRepository 员工合同仓储接口
type EmployeeContractRepository interface {
	// Create 创建合同
	Create(ctx context.Context, contract *entity.EmployeeContract) error

	// GetByID 根据ID获取合同
	GetByID(ctx context.Context, contractID string) (*entity.EmployeeContract, error)

	// GetByEmployeeID 获取员工的所有合同
	GetByEmployeeID(ctx context.Context, empID string) ([]*entity.EmployeeContract, error)

	// GetActiveByEmployeeID 获取员工的生效合同
	GetActiveByEmployeeID(ctx context.Context, empID string) (*entity.EmployeeContract, error)

	// Update 更新合同
	Update(ctx context.Context, contract *entity.EmployeeContract) error

	// Delete 软删除合同
	Delete(ctx context.Context, contractID string) error

	// List 分页查询合同列表
	List(ctx context.Context, filter *ContractFilter) ([]*entity.EmployeeContract, int64, error)
}

// ContractFilter 合同查询过滤器
type ContractFilter struct {
	TenantID      string
	EmpID         string
	ContractType  string
	Status        string
	PageToken     string
	PageSize      int
}

// EmployeeTransferRepository 员工调岗仓储接口
type EmployeeTransferRepository interface {
	// Create 创建调岗记录
	Create(ctx context.Context, transfer *entity.EmployeeTransfer) error

	// GetByID 根据ID获取调岗记录
	GetByID(ctx context.Context, transferID string) (*entity.EmployeeTransfer, error)

	// GetByEmployeeID 获取员工的调岗记录
	GetByEmployeeID(ctx context.Context, empID string) ([]*entity.EmployeeTransfer, error)

	// List 分页查询调岗记录列表
	List(ctx context.Context, filter *TransferFilter) ([]*entity.EmployeeTransfer, int64, error)
}

// TransferFilter 调岗记录查询过滤器
type TransferFilter struct {
	TenantID     string
	EmpID        string
	TransferType string
	PageToken    string
	PageSize     int
}

// EmployeeResignationRepository 员工离职仓储接口
type EmployeeResignationRepository interface {
	// Create 创建离职记录
	Create(ctx context.Context, resignation *entity.EmployeeResignation) error

	// GetByID 根据ID获取离职记录
	GetByID(ctx context.Context, resignationID string) (*entity.EmployeeResignation, error)

	// GetByEmployeeID 获取员工的离职记录
	GetByEmployeeID(ctx context.Context, empID string) (*entity.EmployeeResignation, error)

	// Update 更新离职记录
	Update(ctx context.Context, resignation *entity.EmployeeResignation) error

	// List 分页查询离职记录列表
	List(ctx context.Context, filter *ResignationFilter) ([]*entity.EmployeeResignation, int64, error)
}

// ResignationFilter 离职记录查询过滤器
type ResignationFilter struct {
	TenantID        string
	EmpID           string
	ApprovalStatus  string
	ResignationType string
	PageToken       string
	PageSize        int
}

// OrganizationTreeRepository 组织树仓储接口（闭包表）
type OrganizationTreeRepository interface {
	// CreatePath 创建路径
	CreatePath(ctx context.Context, path *entity.OrganizationTree) error

	// CreatePaths 批量创建路径
	CreatePaths(ctx context.Context, paths []*entity.OrganizationTree) error

	// GetAncestors 获取祖先节点
	GetAncestors(ctx context.Context, descendantID string) ([]*entity.OrganizationTree, error)

	// GetDescendants 获取后代节点
	GetDescendants(ctx context.Context, ancestorID string) ([]*entity.OrganizationTree, error)

	// GetPath 获取路径（从根到指定节点）
	GetPath(ctx context.Context, nodeID string) ([]*entity.OrganizationTree, error)

	// DeletePaths 删除节点的所有路径
	DeletePaths(ctx context.Context, nodeID string) error

	// MoveSubtree 移动子树（更新路径）
	MoveSubtree(ctx context.Context, nodeID, newParentID string) error
}

// DepartmentTreeRepository 部门树仓储接口（闭包表）
type DepartmentTreeRepository interface {
	// CreatePath 创建路径
	CreatePath(ctx context.Context, path *entity.DepartmentTree) error

	// CreatePaths 批量创建路径
	CreatePaths(ctx context.Context, paths []*entity.DepartmentTree) error

	// GetAncestors 获取祖先节点
	GetAncestors(ctx context.Context, descendantID string) ([]*entity.DepartmentTree, error)

	// GetDescendants 获取后代节点
	GetDescendants(ctx context.Context, ancestorID string) ([]*entity.DepartmentTree, error)

	// GetPath 获取路径（从根到指定节点）
	GetPath(ctx context.Context, nodeID string) ([]*entity.DepartmentTree, error)

	// DeletePaths 删除节点的所有路径
	DeletePaths(ctx context.Context, nodeID string) error

	// MoveSubtree 移动子树（更新路径）
	MoveSubtree(ctx context.Context, nodeID, newParentID string) error
}
