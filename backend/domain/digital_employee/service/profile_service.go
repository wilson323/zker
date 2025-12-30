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

	"github.com/google/uuid"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/repository"
)

//go:generate mockgen -destination ../../../internal/mock/domain/digital_employee/service/profile_service_mock.go --package mockProfileService -source profile_service.go
type ProfileService interface {
	// CreateProfile 创建员工画像
	CreateProfile(ctx context.Context, req *entity.CreateProfileRequest) (*entity.EmployeeProfile, error)

	// UpdateProfile 更新员工画像
	UpdateProfile(ctx context.Context, req *entity.UpdateProfileRequest) error

	// GetProfile 获取员工画像
	GetProfile(ctx context.Context, employeeID string) (*entity.EmployeeProfile, error)

	// ListProfiles 列出员工画像
	ListProfiles(ctx context.Context, req *entity.ListProfilesRequest) (*entity.ListProfilesResponse, error)

	// DeleteProfile 删除员工画像
	DeleteProfile(ctx context.Context, employeeID string) error
}

type profileServiceImpl struct {
	profileRepo repository.ProfileRepository
}

// NewProfileService 创建员工画像服务
func NewProfileService(profileRepo repository.ProfileRepository) ProfileService {
	return &profileServiceImpl{
		profileRepo: profileRepo,
	}
}

// CreateProfile 创建员工画像
func (s *profileServiceImpl) CreateProfile(ctx context.Context, req *entity.CreateProfileRequest) (*entity.EmployeeProfile, error) {
	// 验证：检查名称是否已存在
	exists, err := s.profileRepo.ExistsByName(ctx, req.TenantID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check name existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("employee profile with name '%s' already exists", req.Name)
	}

	// 创建员工画像
	profile := &entity.EmployeeProfile{
		EmployeeID:     uuid.New().String(),
		TenantID:       req.TenantID,
		Name:           req.Name,
		Avatar:         req.Avatar,
		Role:           req.Role,
		Skills:         req.Skills,
		Personality:    req.Personality,
		KnowledgeBaseID: req.KnowledgeBaseID,
		BotID:          req.BotID,
		Status:         entity.EmployeeStatusActive,
	}

	if err := s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	return profile, nil
}

// UpdateProfile 更新员工画像
func (s *profileServiceImpl) UpdateProfile(ctx context.Context, req *entity.UpdateProfileRequest) error {
	// 获取现有画像
	profile, err := s.profileRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}

	// 更新字段
	if req.Name != "" {
		// 检查名称是否重复
		exists, err := s.profileRepo.ExistsByName(ctx, profile.TenantID, req.Name)
		if err != nil {
			return fmt.Errorf("failed to check name existence: %w", err)
		}
		if exists {
			return fmt.Errorf("employee profile with name '%s' already exists", req.Name)
		}
		profile.Name = req.Name
	}
	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	if req.Role != "" {
		profile.Role = req.Role
	}
	if req.Skills != nil {
		profile.Skills = req.Skills
	}
	if req.Personality != "" {
		profile.Personality = req.Personality
	}
	if req.KnowledgeBaseID != nil {
		profile.KnowledgeBaseID = req.KnowledgeBaseID
	}

	// 保存更新
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}

	return nil
}

// GetProfile 获取员工画像
func (s *profileServiceImpl) GetProfile(ctx context.Context, employeeID string) (*entity.EmployeeProfile, error) {
	profile, err := s.profileRepo.GetByID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile, nil
}

// ListProfiles 列出员工画像
func (s *profileServiceImpl) ListProfiles(ctx context.Context, req *entity.ListProfilesRequest) (*entity.ListProfilesResponse, error) {
	profiles, total, err := s.profileRepo.GetByTenantID(ctx, req.TenantID, req.Role, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}

	return &entity.ListProfilesResponse{
		Profiles: profiles,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// DeleteProfile 删除员工画像
func (s *profileServiceImpl) DeleteProfile(ctx context.Context, employeeID string) error {
	if err := s.profileRepo.Delete(ctx, employeeID); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}
	return nil
}
