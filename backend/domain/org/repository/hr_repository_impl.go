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
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
)

// employeeContractRepository 员工合同仓储实现
type employeeContractRepository struct {
	db *gorm.DB
}

// NewEmployeeContractRepository 创建员工合同仓储实例
func NewEmployeeContractRepository(db *gorm.DB) EmployeeContractRepository {
	return &employeeContractRepository{db: db}
}

// Create 创建合同
func (r *employeeContractRepository) Create(ctx context.Context, contract *entity.EmployeeContract) error {
	contract.CreatedAt = time.Now().UnixMilli()
	contract.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(contract).Error
}

// GetByID 根据ID获取合同
func (r *employeeContractRepository) GetByID(ctx context.Context, contractID string) (*entity.EmployeeContract, error) {
	var contract entity.EmployeeContract
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Where("contract_id = ?", contractID).
		Where("deleted_at IS NULL").
		First(&contract).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// GetByEmployeeID 获取员工的所有合同
func (r *employeeContractRepository) GetByEmployeeID(ctx context.Context, empID string) ([]*entity.EmployeeContract, error) {
	var contracts []*entity.EmployeeContract
	err := r.db.WithContext(ctx).
		Where("emp_id = ?", empID).
		Where("deleted_at IS NULL").
		Order("start_date DESC").
		Find(&contracts).Error
	return contracts, err
}

// GetActiveByEmployeeID 获取员工的生效合同
func (r *employeeContractRepository) GetActiveByEmployeeID(ctx context.Context, empID string) (*entity.EmployeeContract, error) {
	var contract entity.EmployeeContract
	now := time.Now().UnixMilli()
	err := r.db.WithContext(ctx).
		Where("emp_id = ?", empID).
		Where("status = ?", "active").
		Where("start_date <= ?", now).
		Where("(end_date IS NULL OR end_date > ?)", now).
		Where("deleted_at IS NULL").
		Order("start_date DESC").
		First(&contract).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &contract, nil
}

// Update 更新合同
func (r *employeeContractRepository) Update(ctx context.Context, contract *entity.EmployeeContract) error {
	contract.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.EmployeeContract{}).
		Where("contract_id = ?", contract.ContractID).
		Updates(contract).Error
}

// Delete 软删除合同
func (r *employeeContractRepository) Delete(ctx context.Context, contractID string) error {
	now := time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.EmployeeContract{}).
		Where("contract_id = ?", contractID).
		Updates(map[string]interface{}{
			"deleted_at": &now,
			"updated_at": now,
		}).Error
}

// List 分页查询合同列表
func (r *employeeContractRepository) List(ctx context.Context, filter *ContractFilter) ([]*entity.EmployeeContract, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.EmployeeContract{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Joins("JOIN employees ON employees.emp_id = employee_contracts.emp_id").
			Where("employees.tenant_id = ?", filter.TenantID)
	}
	if filter.EmpID != "" {
		query = query.Where("emp_id = ?", filter.EmpID)
	}
	if filter.ContractType != "" {
		query = query.Where("contract_type = ?", filter.ContractType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 软删除过滤
	query = query.Where("deleted_at IS NULL")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parseInt(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	query = query.Order("start_date DESC")

	// 查询
	var contracts []*entity.EmployeeContract
	err := query.Preload("Employee").Find(&contracts).Error
	return contracts, total, err
}

// employeeTransferRepository 员工调岗仓储实现
type employeeTransferRepository struct {
	db *gorm.DB
}

// NewEmployeeTransferRepository 创建员工调岗仓储实例
func NewEmployeeTransferRepository(db *gorm.DB) EmployeeTransferRepository {
	return &employeeTransferRepository{db: db}
}

// Create 创建调岗记录
func (r *employeeTransferRepository) Create(ctx context.Context, transfer *entity.EmployeeTransfer) error {
	transfer.CreatedAt = time.Now().UnixMilli()
	transfer.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(transfer).Error
}

// GetByID 根据ID获取调岗记录
func (r *employeeTransferRepository) GetByID(ctx context.Context, transferID string) (*entity.EmployeeTransfer, error) {
	var transfer entity.EmployeeTransfer
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Where("transfer_id = ?", transferID).
		First(&transfer).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetByEmployeeID 获取员工的调岗记录
func (r *employeeTransferRepository) GetByEmployeeID(ctx context.Context, empID string) ([]*entity.EmployeeTransfer, error) {
	var transfers []*entity.EmployeeTransfer
	err := r.db.WithContext(ctx).
		Where("emp_id = ?", empID).
		Order("transfer_date DESC").
		Find(&transfers).Error
	return transfers, err
}

// List 分页查询调岗记录列表
func (r *employeeTransferRepository) List(ctx context.Context, filter *TransferFilter) ([]*entity.EmployeeTransfer, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.EmployeeTransfer{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Joins("JOIN employees ON employees.emp_id = employee_transfers.emp_id").
			Where("employees.tenant_id = ?", filter.TenantID)
	}
	if filter.EmpID != "" {
		query = query.Where("emp_id = ?", filter.EmpID)
	}
	if filter.TransferType != "" {
		query = query.Where("transfer_type = ?", filter.TransferType)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parseInt(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	query = query.Order("transfer_date DESC")

	// 查询
	var transfers []*entity.EmployeeTransfer
	err := query.Preload("Employee").Find(&transfers).Error
	return transfers, total, err
}

// employeeResignationRepository 员工离职仓储实现
type employeeResignationRepository struct {
	db *gorm.DB
}

// NewEmployeeResignationRepository 创建员工离职仓储实例
func NewEmployeeResignationRepository(db *gorm.DB) EmployeeResignationRepository {
	return &employeeResignationRepository{db: db}
}

// Create 创建离职记录
func (r *employeeResignationRepository) Create(ctx context.Context, resignation *entity.EmployeeResignation) error {
	resignation.CreatedAt = time.Now().UnixMilli()
	resignation.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).Create(resignation).Error
}

// GetByID 根据ID获取离职记录
func (r *employeeResignationRepository) GetByID(ctx context.Context, resignationID string) (*entity.EmployeeResignation, error) {
	var resignation entity.EmployeeResignation
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Where("resignation_id = ?", resignationID).
		First(&resignation).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &resignation, nil
}

// GetByEmployeeID 获取员工的离职记录
func (r *employeeResignationRepository) GetByEmployeeID(ctx context.Context, empID string) (*entity.EmployeeResignation, error) {
	var resignation entity.EmployeeResignation
	err := r.db.WithContext(ctx).
		Where("emp_id = ?", empID).
		Order("apply_date DESC").
		First(&resignation).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &resignation, nil
}

// Update 更新离职记录
func (r *employeeResignationRepository) Update(ctx context.Context, resignation *entity.EmployeeResignation) error {
	resignation.UpdatedAt = time.Now().UnixMilli()
	return r.db.WithContext(ctx).
		Model(&entity.EmployeeResignation{}).
		Where("resignation_id = ?", resignation.ResignationID).
		Updates(resignation).Error
}

// List 分页查询离职记录列表
func (r *employeeResignationRepository) List(ctx context.Context, filter *ResignationFilter) ([]*entity.EmployeeResignation, int64, error) {
	query := r.db.WithContext(ctx).Model(&entity.EmployeeResignation{})

	// 应用过滤条件
	if filter.TenantID != "" {
		query = query.Joins("JOIN employees ON employees.emp_id = employee_resignations.emp_id").
			Where("employees.tenant_id = ?", filter.TenantID)
	}
	if filter.EmpID != "" {
		query = query.Where("emp_id = ?", filter.EmpID)
	}
	if filter.ApprovalStatus != "" {
		query = query.Where("approval_status = ?", filter.ApprovalStatus)
	}
	if filter.ResignationType != "" {
		query = query.Where("resignation_type = ?", filter.ResignationType)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter.PageSize > 0 {
		offset := 0
		if filter.PageToken != "" {
			offset = parseInt(filter.PageToken)
		}
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	// 排序
	query = query.Order("apply_date DESC")

	// 查询
	var resignations []*entity.EmployeeResignation
	err := query.Preload("Employee").Find(&resignations).Error
	return resignations, total, err
}

// organizationTreeRepository 组织树仓储实现（闭包表）
type organizationTreeRepository struct {
	db *gorm.DB
}

// NewOrganizationTreeRepository 创建组织树仓储实例
func NewOrganizationTreeRepository(db *gorm.DB) OrganizationTreeRepository {
	return &organizationTreeRepository{db: db}
}

// CreatePath 创建路径
func (r *organizationTreeRepository) CreatePath(ctx context.Context, path *entity.OrganizationTree) error {
	return r.db.WithContext(ctx).Create(path).Error
}

// CreatePaths 批量创建路径
func (r *organizationTreeRepository) CreatePaths(ctx context.Context, paths []*entity.OrganizationTree) error {
	if len(paths) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(paths, 100).Error
}

// GetAncestors 获取祖先节点
func (r *organizationTreeRepository) GetAncestors(ctx context.Context, descendantID string) ([]*entity.OrganizationTree, error) {
	var paths []*entity.OrganizationTree
	err := r.db.WithContext(ctx).
		Where("descendant_id = ?", descendantID).
		Where("depth > 0").
		Order("depth ASC").
		Find(&paths).Error
	return paths, err
}

// GetDescendants 获取后代节点
func (r *organizationTreeRepository) GetDescendants(ctx context.Context, ancestorID string) ([]*entity.OrganizationTree, error) {
	var paths []*entity.OrganizationTree
	err := r.db.WithContext(ctx).
		Where("ancestor_id = ?", ancestorID).
		Where("depth > 0").
		Order("depth ASC").
		Find(&paths).Error
	return paths, err
}

// GetPath 获取路径（从根到指定节点）
func (r *organizationTreeRepository) GetPath(ctx context.Context, nodeID string) ([]*entity.OrganizationTree, error) {
	var paths []*entity.OrganizationTree
	err := r.db.WithContext(ctx).
		Where("descendant_id = ?", nodeID).
		Order("depth ASC").
		Find(&paths).Error
	return paths, err
}

// DeletePaths 删除节点的所有路径
func (r *organizationTreeRepository) DeletePaths(ctx context.Context, nodeID string) error {
	return r.db.WithContext(ctx).
		Where("ancestor_id = ? OR descendant_id = ?", nodeID, nodeID).
		Delete(&entity.OrganizationTree{}).Error
}

// MoveSubtree 移动子树（更新路径）
func (r *organizationTreeRepository) MoveSubtree(ctx context.Context, nodeID, newParentID string) error {
	// 1. 删除原节点的所有路径
	if err := r.DeletePaths(ctx, nodeID); err != nil {
		return err
	}

	// 2. 获取新父节点的所有路径
	parentPaths, err := r.GetPath(ctx, newParentID)
	if err != nil {
		return err
	}

	// 3. 创建新路径
	var newPaths []*entity.OrganizationTree
	depth := len(parentPaths)

	// 添加自己到自己的路径（深度为0）
	newPaths = append(newPaths, &entity.OrganizationTree{
		TenantID:     parentPaths[0].TenantID,
		AncestorID:   nodeID,
		DescendantID: nodeID,
		Depth:        0,
	})

	// 添加到所有祖先的路径
	for i, parentPath := range parentPaths {
		newPaths = append(newPaths, &entity.OrganizationTree{
			TenantID:     parentPath.TenantID,
			AncestorID:   parentPath.AncestorID,
			DescendantID: nodeID,
			Depth:        depth - i,
		})
	}

	// 4. 获取所有后代节点
	descendants, err := r.GetDescendants(ctx, nodeID)
	if err != nil {
		return err
	}

	// 5. 为每个后代节点创建新路径
	for _, descendant := range descendants {
		// 添加自己到自己的路径
		newPaths = append(newPaths, &entity.OrganizationTree{
			TenantID:     parentPaths[0].TenantID,
			AncestorID:   descendant.DescendantID,
			DescendantID: descendant.DescendantID,
			Depth:        0,
		})

		// 添加到所有祖先的路径
		for i, parentPath := range parentPaths {
			newPaths = append(newPaths, &entity.OrganizationTree{
				TenantID:     parentPath.TenantID,
				AncestorID:   parentPath.AncestorID,
				DescendantID: descendant.DescendantID,
				Depth:        depth - i + descendant.Depth,
			})
		}
	}

	return r.CreatePaths(ctx, newPaths)
}

// departmentTreeRepository 部门树仓储实现（闭包表）
type departmentTreeRepository struct {
	db *gorm.DB
}

// NewDepartmentTreeRepository 创建部门树仓储实例
func NewDepartmentTreeRepository(db *gorm.DB) DepartmentTreeRepository {
	return &departmentTreeRepository{db: db}
}

// CreatePath 创建路径
func (r *departmentTreeRepository) CreatePath(ctx context.Context, path *entity.DepartmentTree) error {
	return r.db.WithContext(ctx).Create(path).Error
}

// CreatePaths 批量创建路径
func (r *departmentTreeRepository) CreatePaths(ctx context.Context, paths []*entity.DepartmentTree) error {
	if len(paths) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(paths, 100).Error
}

// GetAncestors 获取祖先节点
func (r *departmentTreeRepository) GetAncestors(ctx context.Context, descendantID string) ([]*entity.DepartmentTree, error) {
	var paths []*entity.DepartmentTree
	err := r.db.WithContext(ctx).
		Where("descendant_id = ?", descendantID).
		Where("depth > 0").
		Order("depth ASC").
		Find(&paths).Error
	return paths, err
}

// GetDescendants 获取后代节点
func (r *departmentTreeRepository) GetDescendants(ctx context.Context, ancestorID string) ([]*entity.DepartmentTree, error) {
	var paths []*entity.DepartmentTree
	err := r.db.WithContext(ctx).
		Where("ancestor_id = ?", ancestorID).
		Where("depth > 0").
		Order("depth ASC").
		Find(&paths).Error
	return paths, err
}

// GetPath 获取路径（从根到指定节点）
func (r *departmentTreeRepository) GetPath(ctx context.Context, nodeID string) ([]*entity.DepartmentTree, error) {
	var paths []*entity.DepartmentTree
	err := r.db.WithContext(ctx).
		Where("descendant_id = ?", nodeID).
		Order("depth ASC").
		Find(&paths).Error
	return paths, err
}

// DeletePaths 删除节点的所有路径
func (r *departmentTreeRepository) DeletePaths(ctx context.Context, nodeID string) error {
	return r.db.WithContext(ctx).
		Where("ancestor_id = ? OR descendant_id = ?", nodeID, nodeID).
		Delete(&entity.DepartmentTree{}).Error
}

// MoveSubtree 移动子树（更新路径）
func (r *departmentTreeRepository) MoveSubtree(ctx context.Context, nodeID, newParentID string) error {
	// 实现逻辑与 organizationTreeRepository.MoveSubtree 相同
	// 1. 删除原节点的所有路径
	if err := r.DeletePaths(ctx, nodeID); err != nil {
		return err
	}

	// 2. 获取新父节点的所有路径
	parentPaths, err := r.GetPath(ctx, newParentID)
	if err != nil {
		return err
	}

	// 3. 创建新路径
	var newPaths []*entity.DepartmentTree
	depth := len(parentPaths)

	// 添加自己到自己的路径（深度为0）
	newPaths = append(newPaths, &entity.DepartmentTree{
		TenantID:     parentPaths[0].TenantID,
		AncestorID:   nodeID,
		DescendantID: nodeID,
		Depth:        0,
	})

	// 添加到所有祖先的路径
	for i, parentPath := range parentPaths {
		newPaths = append(newPaths, &entity.DepartmentTree{
			TenantID:     parentPath.TenantID,
			AncestorID:   parentPath.AncestorID,
			DescendantID: nodeID,
			Depth:        depth - i,
		})
	}

	// 4. 获取所有后代节点
	descendants, err := r.GetDescendants(ctx, nodeID)
	if err != nil {
		return err
	}

	// 5. 为每个后代节点创建新路径
	for _, descendant := range descendants {
		// 添加自己到自己的路径
		newPaths = append(newPaths, &entity.DepartmentTree{
			TenantID:     parentPaths[0].TenantID,
			AncestorID:   descendant.DescendantID,
			DescendantID: descendant.DescendantID,
			Depth:        0,
		})

		// 添加到所有祖先的路径
		for i, parentPath := range parentPaths {
			newPaths = append(newPaths, &entity.DepartmentTree{
				TenantID:     parentPath.TenantID,
				AncestorID:   parentPath.AncestorID,
				DescendantID: descendant.DescendantID,
				Depth:        depth - i + descendant.Depth,
			})
		}
	}

	return r.CreatePaths(ctx, newPaths)
}
