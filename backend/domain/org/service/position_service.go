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

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// PositionService 岗位管理服务
// 职责：岗位的CRUD操作、岗位查询、业务规则验证
type PositionService struct {
	positionRepo repository.PositionRepository
	deptRepo     repository.DepartmentRepository
}

// NewPositionService 创建岗位服务实例
func NewPositionService(
	positionRepo repository.PositionRepository,
	deptRepo repository.DepartmentRepository,
) *PositionService {
	return &PositionService{
		positionRepo: positionRepo,
		deptRepo:     deptRepo,
	}
}

// CreatePositionRequest 创建岗位请求
type CreatePositionRequest struct {
	TenantID         string  `json:"tenant_id" binding:"required"`
	DeptID           *string `json:"dept_id,omitempty"`
	PositionName     string  `json:"position_name" binding:"required,min=1,max=100"`
	PositionCode     string  `json:"position_code" binding:"required,min=1,max=50"`
	Level            int     `json:"level" binding:"required,min=1,max=10"`
	Category         string  `json:"category" binding:"required,min=1,max=50"` // 技术岗/管理岗/职能岗
	Responsibilities string  `json:"responsibilities,omitempty"`
	Requirements     string  `json:"requirements,omitempty"`
	SortOrder        int     `json:"sort_order,omitempty"`
}

// UpdatePositionRequest 更新岗位请求
type UpdatePositionRequest struct {
	PositionID       string  `json:"position_id" binding:"required"`
	PositionName     string  `json:"position_name" binding:"omitempty,min=1,max=100"`
	Level            int     `json:"level" binding:"omitempty,min=1,max=10"`
	Category         string  `json:"category" binding:"omitempty,min=1,max=50"`
	Responsibilities string  `json:"responsibilities,omitempty"`
	Requirements     string  `json:"requirements,omitempty"`
	SortOrder        int     `json:"sort_order,omitempty"`
	Status           string  `json:"status" binding:"omitempty,oneof=active inactive frozen"`
}

// CreatePosition 创建岗位
func (s *PositionService) CreatePosition(ctx context.Context, req *CreatePositionRequest) (*entity.Position, error) {
	// 1. 如果有部门，验证部门存在
	if req.DeptID != nil && *req.DeptID != "" {
		dept, err := s.deptRepo.GetByID(ctx, *req.DeptID)
		if err != nil {
			return nil, fmt.Errorf("failed to get department: %w", err)
		}
		if dept == nil {
			return nil, errno.ErrDeptNotFound
		}
		if dept.TenantID != req.TenantID {
			return nil, errno.ErrTenantMismatch
		}
	}

	// 2. 验证编码唯一性
	exists, err := s.positionRepo.ExistsByCode(ctx, req.TenantID, req.PositionCode, "")
	if err != nil {
		return nil, fmt.Errorf("failed to check position code: %w", err)
	}
	if exists {
		return nil, errno.ErrPositionCodeAlreadyExists
	}

	// 3. 生成岗位ID
	positionID := generateUUID()

	// 4. 创建岗位实体
	position := &entity.Position{
		PositionID:      positionID,
		TenantID:        req.TenantID,
		DeptID:          req.DeptID,
		PositionName:    req.PositionName,
		PositionCode:    req.PositionCode,
		Level:           req.Level,
		Category:        req.Category,
		Responsibilities: req.Responsibilities,
		Requirements:    req.Requirements,
		Status:          entity.OrgStatusActive,
		SortOrder:       req.SortOrder,
		CreatedAt:       time.Now().UnixMilli(),
		UpdatedAt:       time.Now().UnixMilli(),
	}

	// 5. 创建岗位
	if err := s.positionRepo.Create(ctx, position); err != nil {
		return nil, fmt.Errorf("failed to create position: %w", err)
	}

	return position, nil
}

// GetPosition 获取岗位详情
func (s *PositionService) GetPosition(ctx context.Context, positionID string) (*entity.Position, error) {
	if positionID == "" {
		return nil, errno.ErrInvalidParam
	}

	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return nil, errno.ErrPositionNotFound
	}

	return position, nil
}

// UpdatePosition 更新岗位
func (s *PositionService) UpdatePosition(ctx context.Context, req *UpdatePositionRequest) (*entity.Position, error) {
	// 1. 获取岗位
	position, err := s.positionRepo.GetByID(ctx, req.PositionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return nil, errno.ErrPositionNotFound
	}

	// 2. 更新字段
	if req.PositionName != "" {
		position.PositionName = req.PositionName
	}
	if req.Level != 0 {
		position.Level = req.Level
	}
	if req.Category != "" {
		position.Category = req.Category
	}
	if req.Responsibilities != "" {
		position.Responsibilities = req.Responsibilities
	}
	if req.Requirements != "" {
		position.Requirements = req.Requirements
	}
	if req.SortOrder != 0 {
		position.SortOrder = req.SortOrder
	}
	if req.Status != "" {
		position.Status = entity.OrganizationStatus(req.Status)
	}

	// 3. 保存更新
	if err := s.positionRepo.Update(ctx, position); err != nil {
		return nil, fmt.Errorf("failed to update position: %w", err)
	}

	return position, nil
}

// DeletePosition 删除岗位
func (s *PositionService) DeletePosition(ctx context.Context, positionID string) error {
	if positionID == "" {
		return errno.ErrInvalidParam
	}

	// 1. 获取岗位
	position, err := s.positionRepo.GetByID(ctx, positionID)
	if err != nil {
		return fmt.Errorf("failed to get position: %w", err)
	}
	if position == nil {
		return errno.ErrPositionNotFound
	}

	// 2. TODO: 检查是否有员工使用此岗位
	// if hasEmployees {
	//     return errno.ErrPositionHasEmployees
	// }

	// 3. 软删除岗位
	if err := s.positionRepo.Delete(ctx, positionID); err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	return nil
}

// GetPositionByCode 根据编码获取岗位
func (s *PositionService) GetPositionByCode(ctx context.Context, tenantID, code string) (*entity.Position, error) {
	if tenantID == "" || code == "" {
		return nil, errno.ErrInvalidParam
	}

	position, err := s.positionRepo.GetByCode(ctx, tenantID, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get position by code: %w", err)
	}
	if position == nil {
		return nil, errno.ErrPositionNotFound
	}

	return position, nil
}

// ListPositions 分页查询岗位列表
func (s *PositionService) ListPositions(ctx context.Context, filter *repository.PositionFilter) ([]*entity.Position, int64, error) {
	if filter.TenantID == "" {
		return nil, 0, errno.ErrInvalidParam
	}

	positions, total, err := s.positionRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list positions: %w", err)
	}

	return positions, total, nil
}

// GetPositionsByDepartment 获取部门的所有岗位
func (s *PositionService) GetPositionsByDepartment(ctx context.Context, deptID string) ([]*entity.Position, error) {
	if deptID == "" {
		return nil, errno.ErrInvalidParam
	}

	positions, err := s.positionRepo.GetByDepartmentID(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions by department: %w", err)
	}

	return positions, nil
}

// GetPositionsByTenant 获取租户的所有岗位
func (s *PositionService) GetPositionsByTenant(ctx context.Context, tenantID string) ([]*entity.Position, error) {
	if tenantID == "" {
		return nil, errno.ErrInvalidParam
	}

	positions, err := s.positionRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions by tenant: %w", err)
	}

	return positions, nil
}

// GetPositionsByLevel 按职级获取岗位
func (s *PositionService) GetPositionsByLevel(ctx context.Context, tenantID string, level int) ([]*entity.Position, error) {
	if tenantID == "" || level <= 0 {
		return nil, errno.ErrInvalidParam
	}

	filter := &repository.PositionFilter{
		TenantID: tenantID,
		Level:    &level,
	}

	positions, _, err := s.positionRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions by level: %w", err)
	}

	return positions, nil
}

// GetPositionsByCategory 按类别获取岗位
func (s *PositionService) GetPositionsByCategory(ctx context.Context, tenantID, category string) ([]*entity.Position, error) {
	if tenantID == "" || category == "" {
		return nil, errno.ErrInvalidParam
	}

	filter := &repository.PositionFilter{
		TenantID: tenantID,
		Category: category,
	}

	positions, _, err := s.positionRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions by category: %w", err)
	}

	return positions, nil
}
