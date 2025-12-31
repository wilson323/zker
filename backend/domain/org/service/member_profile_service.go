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
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}
	if exp.CompanyName == "" || exp.Position == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 2. 验证员工存在
	emp, err := s.empRepo.GetByUserID(ctx, exp.UserID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}
	if emp.TenantID != exp.TenantID {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "TenantMismatch"),
            )
	}

	// 3. 验证日期逻辑
	if exp.EndDate != nil && *exp.EndDate <= exp.StartDate {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidDateRange"),
                )
	}

	// 4. 创建工作经历
	exp.CreatedAt = time.Now().UnixMilli()
	exp.UpdatedAt = time.Now().UnixMilli()
	if err := s.workExpRepo.Create(ctx, exp); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create work experience"),
            )
	}

	return nil
}

// UpdateWorkExperience 更新工作经历
func (s *MemberProfileService) UpdateWorkExperience(ctx context.Context, id int64, exp *entity.WorkExperience) error {
	// 1. 获取现有记录
	existing, err := s.workExpRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get work experience"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "WorkExperienceNotFound"),
            )
	}

	// 2. 验证权限
	if existing.UserID != exp.UserID {
		return errorx.New(errno.ErrPermissionDeniedCode,
                errorx.KV("reason", "PermissionDenied"),
            )
	}

	// 3. 验证日期逻辑
	if exp.EndDate != nil && *exp.EndDate <= exp.StartDate {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidDateRange"),
                )
	}

	// 4. 更新
	exp.ID = id
	exp.UpdatedAt = time.Now().UnixMilli()
	if err := s.workExpRepo.Update(ctx, exp); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update work experience"),
            )
	}

	return nil
}

// DeleteWorkExperience 删除工作经历
func (s *MemberProfileService) DeleteWorkExperience(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.workExpRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get work experience"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "WorkExperienceNotFound"),
            )
	}

	// 2. 删除
	if err := s.workExpRepo.Delete(ctx, id); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "delete work experience"),
            )
	}

	return nil
}

// ListWorkExperiences 列出工作经历
func (s *MemberProfileService) ListWorkExperiences(ctx context.Context, userID string) ([]*entity.WorkExperience, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	exps, err := s.workExpRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list work experiences"),
            )
	}

	return exps, nil
}

// AddEducation 添加教育经历
func (s *MemberProfileService) AddEducation(ctx context.Context, edu *entity.Education) error {
	// 1. 验证参数
	if edu.UserID == "" || edu.TenantID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}
	if edu.SchoolName == "" || edu.Major == "" || edu.Degree == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 2. 验证员工存在
	emp, err := s.empRepo.GetByUserID(ctx, edu.UserID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}
	if emp.TenantID != edu.TenantID {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "TenantMismatch"),
            )
	}

	// 3. 验证日期逻辑
	if edu.EndDate != nil && *edu.EndDate <= edu.StartDate {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidDateRange"),
                )
	}

	// 4. 创建教育经历
	edu.CreatedAt = time.Now().UnixMilli()
	edu.UpdatedAt = time.Now().UnixMilli()
	if err := s.eduRepo.Create(ctx, edu); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create education"),
            )
	}

	return nil
}

// UpdateEducation 更新教育经历
func (s *MemberProfileService) UpdateEducation(ctx context.Context, id int64, edu *entity.Education) error {
	// 1. 获取现有记录
	existing, err := s.eduRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get education"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "EducationNotFound"),
            )
	}

	// 2. 验证权限
	if existing.UserID != edu.UserID {
		return errorx.New(errno.ErrPermissionDeniedCode,
                errorx.KV("reason", "PermissionDenied"),
            )
	}

	// 3. 验证日期逻辑
	if edu.EndDate != nil && *edu.EndDate <= edu.StartDate {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidDateRange"),
                )
	}

	// 4. 更新
	edu.ID = id
	edu.UpdatedAt = time.Now().UnixMilli()
	if err := s.eduRepo.Update(ctx, edu); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update education"),
            )
	}

	return nil
}

// DeleteEducation 删除教育经历
func (s *MemberProfileService) DeleteEducation(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.eduRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get education"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "EducationNotFound"),
            )
	}

	// 2. 删除
	if err := s.eduRepo.Delete(ctx, id); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "delete education"),
            )
	}

	return nil
}

// ListEducations 列出教育经历
func (s *MemberProfileService) ListEducations(ctx context.Context, userID string) ([]*entity.Education, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	edus, err := s.eduRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list educations"),
            )
	}

	return edus, nil
}

// AddSkill 添加技能
func (s *MemberProfileService) AddSkill(ctx context.Context, skill *entity.MemberSkill) error {
	// 1. 验证参数
	if skill.UserID == "" || skill.TenantID == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}
	if skill.SkillName == "" || skill.Proficiency == "" {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "InvalidParam"),
                )
	}

	// 2. 验证员工存在
	emp, err := s.empRepo.GetByUserID(ctx, skill.UserID)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                    errorx.KV("reason", "EmployeeNotFound"),
                )
	}
	if emp.TenantID != skill.TenantID {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "TenantMismatch"),
            )
	}

	// 3. 创建技能
	skill.CreatedAt = time.Now().UnixMilli()
	skill.UpdatedAt = time.Now().UnixMilli()
	if err := s.skillRepo.Create(ctx, skill); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "create skill"),
            )
	}

	return nil
}

// UpdateSkill 更新技能
func (s *MemberProfileService) UpdateSkill(ctx context.Context, id int64, skill *entity.MemberSkill) error {
	// 1. 获取现有记录
	existing, err := s.skillRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get skill"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "SkillNotFound"),
            )
	}

	// 2. 验证权限
	if existing.UserID != skill.UserID {
		return errorx.New(errno.ErrPermissionDeniedCode,
                errorx.KV("reason", "PermissionDenied"),
            )
	}

	// 3. 更新
	skill.ID = id
	skill.UpdatedAt = time.Now().UnixMilli()
	if err := s.skillRepo.Update(ctx, skill); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "update skill"),
            )
	}

	return nil
}

// DeleteSkill 删除技能
func (s *MemberProfileService) DeleteSkill(ctx context.Context, id int64) error {
	// 1. 获取现有记录
	existing, err := s.skillRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get skill"),
            )
	}
	if existing == nil {
		return errorx.New(errno.ErrPermissionInvalidParamCode,
                errorx.KV("reason", "SkillNotFound"),
            )
	}

	// 2. 删除
	if err := s.skillRepo.Delete(ctx, id); err != nil {
		return errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "delete skill"),
            )
	}

	return nil
}

// ListSkills 列出技能
func (s *MemberProfileService) ListSkills(ctx context.Context, userID string) ([]*entity.MemberSkill, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	skills, err := s.skillRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list skills"),
            )
	}

	return skills, nil
}

// GenerateResume 生成简历
func (s *MemberProfileService) GenerateResume(ctx context.Context, userID string) (*entity.Resume, error) {
	if userID == "" {
		return nil, errorx.NewByErrorCode(errno.ErrInvalidParam)
	}

	// 1. 获取员工信息
	emp, err := s.empRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "get employee"),
            )
	}
	if emp == nil {
		return nil, errorx.NewByErrorCode(errno.ErrEmployeeNotFound)
	}

	// 2. 获取工作经历
	workExps, err := s.workExpRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list work experiences"),
            )
	}

	// 3. 获取教育经历
	educations, err := s.eduRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list educations"),
            )
	}

	// 4. 获取技能
	skills, err := s.skillRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrPermissionCheckFailedCode,
                errorx.KV("operation", "list skills"),
            )
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
