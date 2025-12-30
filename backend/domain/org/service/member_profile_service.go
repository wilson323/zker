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

// MemberProfileService 员工档案服务
type MemberProfileService struct {
	workExpRepo repository.WorkExperienceRepository
	eduRepo     repository.EducationRepository
	skillRepo   repository.MemberSkillRepository
	empRepo     repository.EmployeeRepository
}

// NewMemberProfileService 创建员工档案服务实例
func NewMemberProfileService(
	workExpRepo repository.WorkExperienceRepository,
	eduRepo repository.EducationRepository,
	skillRepo repository.MemberSkillRepository,
	empRepo repository.EmployeeRepository,
) *MemberProfileService {
	return &MemberProfileService{
		workExpRepo: workExpRepo,
		eduRepo:     eduRepo,
		skillRepo:   skillRepo,
		empRepo:     empRepo,
	}
}

// AddWorkExperience 添加工作经历
func (s *MemberProfileService) AddWorkExperience(ctx context.Context, exp *entity.WorkExperience) error {
	// 1. 验证参数
	if exp.UserID == "" || exp.TenantID == "" {
		return errno.ErrInvalidParam
	}
	if exp.CompanyName == "" || exp.Position == "" {
		return errno.ErrInvalidParam
	}

	// 2. 验证员工存在
	emp, err := s.empRepo.GetByUserID(ctx, exp.UserID)
	if err != nil {
		return fmt.Errorf("failed to get employee: %w", err)
	}
	if emp == nil {
		return errno.ErrEmployeeNotFound
	}
	if emp.TenantID != exp.TenantID {
		return errno.ErrTenantMismatch
	}

	// 3. 验证日期逻辑
	if exp.EndDate != nil && *exp.EndDate <= exp.StartDate {
		return errno.ErrInvalidDateRange
	}

	// 4. 创建工作经历
	exp.CreatedAt = time.Now().UnixMilli()
	exp.UpdatedAt = time.Now().UnixMilli()
	if err := s.workExpRepo.Create(ctx, exp); err != nil {
		return fmt.Errorf("failed to create work experience: %w", err)
	}

	return nil
}

// UpdateWorkExperience 更新工作经历
func (s *MemberProfileService) UpdateWorkExperience(ctx context.Context, id int64, exp *entity.WorkExperience) error {
	// 1. 获取现有记录
	existing, err := s.workExpRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get work experience: %w", err)
	}
	if existing == nil {
		return errno.ErrWorkExperienceNotFound
	}

	// 2. 验证权限
	if existing.UserID != exp.UserID {
		return errno.ErrPermissionDenied
	}

	// 3. 验证日期逻辑
	if exp.EndDate != nil && *exp.EndDate <= exp.StartDate {
		return errno.ErrInvalidDateRange
	}

	// 4. 更新
	exp.ID = id
	exp.UpdatedAt = time.Now().UnixMilli()
	if err := s.workExpRepo.Update(ctx, exp); err != nil {
		return fmt.Errorf("failed to update work experience: %w", err)
	}

	return nil
}

// DeleteWorkExperience 删除工作经历
func (s *MemberProfileService) DeleteWorkExperience(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.workExpRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get work experience: %w", err)
	}
	if existing == nil {
		return errno.ErrWorkExperienceNotFound
	}

	// 2. 删除
	if err := s.workExpRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete work experience: %w", err)
	}

	return nil
}

// ListWorkExperiences 列出工作经历
func (s *MemberProfileService) ListWorkExperiences(ctx context.Context, userID string) ([]*entity.WorkExperience, error) {
	if userID == "" {
		return nil, errno.ErrInvalidParam
	}

	exps, err := s.workExpRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list work experiences: %w", err)
	}

	return exps, nil
}

// AddEducation 添加教育经历
func (s *MemberProfileService) AddEducation(ctx context.Context, edu *entity.Education) error {
	// 1. 验证参数
	if edu.UserID == "" || edu.TenantID == "" {
		return errno.ErrInvalidParam
	}
	if edu.SchoolName == "" || edu.Major == "" || edu.Degree == "" {
		return errno.ErrInvalidParam
	}

	// 2. 验证员工存在
	emp, err := s.empRepo.GetByUserID(ctx, edu.UserID)
	if err != nil {
		return fmt.Errorf("failed to get employee: %w", err)
	}
	if emp == nil {
		return errno.ErrEmployeeNotFound
	}
	if emp.TenantID != edu.TenantID {
		return errno.ErrTenantMismatch
	}

	// 3. 验证日期逻辑
	if edu.EndDate != nil && *edu.EndDate <= edu.StartDate {
		return errno.ErrInvalidDateRange
	}

	// 4. 创建教育经历
	edu.CreatedAt = time.Now().UnixMilli()
	edu.UpdatedAt = time.Now().UnixMilli()
	if err := s.eduRepo.Create(ctx, edu); err != nil {
		return fmt.Errorf("failed to create education: %w", err)
	}

	return nil
}

// UpdateEducation 更新教育经历
func (s *MemberProfileService) UpdateEducation(ctx context.Context, id int64, edu *entity.Education) error {
	// 1. 获取现有记录
	existing, err := s.eduRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get education: %w", err)
	}
	if existing == nil {
		return errno.ErrEducationNotFound
	}

	// 2. 验证权限
	if existing.UserID != edu.UserID {
		return errno.ErrPermissionDenied
	}

	// 3. 验证日期逻辑
	if edu.EndDate != nil && *edu.EndDate <= edu.StartDate {
		return errno.ErrInvalidDateRange
	}

	// 4. 更新
	edu.ID = id
	edu.UpdatedAt = time.Now().UnixMilli()
	if err := s.eduRepo.Update(ctx, edu); err != nil {
		return fmt.Errorf("failed to update education: %w", err)
	}

	return nil
}

// DeleteEducation 删除教育经历
func (s *MemberProfileService) DeleteEducation(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.eduRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get education: %w", err)
	}
	if existing == nil {
		return errno.ErrEducationNotFound
	}

	// 2. 删除
	if err := s.eduRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	return nil
}

// ListEducations 列出教育经历
func (s *MemberProfileService) ListEducations(ctx context.Context, userID string) ([]*entity.Education, error) {
	if userID == "" {
		return nil, errno.ErrInvalidParam
	}

	edus, err := s.eduRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list educations: %w", err)
	}

	return edus, nil
}

// AddSkill 添加技能
func (s *MemberProfileService) AddSkill(ctx context.Context, skill *entity.MemberSkill) error {
	// 1. 验证参数
	if skill.UserID == "" || skill.TenantID == "" {
		return errno.ErrInvalidParam
	}
	if skill.SkillName == "" || skill.Proficiency == "" {
		return errno.ErrInvalidParam
	}

	// 2. 验证员工存在
	emp, err := s.empRepo.GetByUserID(ctx, skill.UserID)
	if err != nil {
		return fmt.Errorf("failed to get employee: %w", err)
	}
	if emp == nil {
		return errno.ErrEmployeeNotFound
	}
	if emp.TenantID != skill.TenantID {
		return errno.ErrTenantMismatch
	}

	// 3. 创建技能
	skill.CreatedAt = time.Now().UnixMilli()
	skill.UpdatedAt = time.Now().UnixMilli()
	if err := s.skillRepo.Create(ctx, skill); err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}

// UpdateSkill 更新技能
func (s *MemberProfileService) UpdateSkill(ctx context.Context, id int64, skill *entity.MemberSkill) error {
	// 1. 获取现有记录
	existing, err := s.skillRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get skill: %w", err)
	}
	if existing == nil {
		return errno.ErrSkillNotFound
	}

	// 2. 验证权限
	if existing.UserID != skill.UserID {
		return errno.ErrPermissionDenied
	}

	// 3. 更新
	skill.ID = id
	skill.UpdatedAt = time.Now().UnixMilli()
	if err := s.skillRepo.Update(ctx, skill); err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	return nil
}

// DeleteSkill 删除技能
func (s *MemberProfileService) DeleteSkill(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.skillRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get skill: %w", err)
	}
	if existing == nil {
		return errno.ErrSkillNotFound
	}

	// 2. 删除
	if err := s.skillRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}

// ListSkills 列出技能
func (s *MemberProfileService) ListSkills(ctx context.Context, userID string) ([]*entity.MemberSkill, error) {
	if userID == "" {
		return nil, errno.ErrInvalidParam
	}

	skills, err := s.skillRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}

	return skills, nil
}

// GenerateResume 生成简历
func (s *MemberProfileService) GenerateResume(ctx context.Context, userID string) (*entity.Resume, error) {
	if userID == "" {
		return nil, errno.ErrInvalidParam
	}

	// 1. 获取员工信息
	emp, err := s.empRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}
	if emp == nil {
		return nil, errno.ErrEmployeeNotFound
	}

	// 2. 获取工作经历
	workExps, err := s.workExpRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list work experiences: %w", err)
	}

	// 3. 获取教育经历
	educations, err := s.eduRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list educations: %w", err)
	}

	// 4. 获取技能
	skills, err := s.skillRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}

	// 5. 构建简历
	resume := &entity.Resume{
		UserID:          userID,
		TenantID:        emp.TenantID,
		Employee:        emp,
		WorkExperiences: workExps,
		Educations:      educations,
		Skills:          skills,
		GeneratedAt:     time.Now().UnixMilli(),
	}

	return resume, nil
}
