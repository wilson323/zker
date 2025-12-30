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

// WorkExperienceRepository 工作经历仓储接口
type WorkExperienceRepository interface {
	Create(ctx context.Context, exp *entity.WorkExperience) error
	Update(ctx context.Context, exp *entity.WorkExperience) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.WorkExperience, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.WorkExperience, error)
}

// EducationRepository 教育经历仓储接口
type EducationRepository interface {
	Create(ctx context.Context, edu *entity.Education) error
	Update(ctx context.Context, edu *entity.Education) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.Education, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.Education, error)
}

// MemberSkillRepository 技能仓储接口
type MemberSkillRepository interface {
	Create(ctx context.Context, skill *entity.MemberSkill) error
	Update(ctx context.Context, skill *entity.MemberSkill) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.MemberSkill, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.MemberSkill, error)
	ListByTenantID(ctx context.Context, tenantID string) ([]*entity.MemberSkill, error)
	SearchByName(ctx context.Context, tenantID, skillName string, limit int) ([]*entity.MemberSkill, error)
}

// VirtualOrganizationRepository 虚拟组织仓储接口
type VirtualOrganizationRepository interface {
	Create(ctx context.Context, org *entity.VirtualOrganization) error
	Update(ctx context.Context, org *entity.VirtualOrganization) error
	Delete(ctx context.Context, virtualOrgID string) error
	GetByID(ctx context.Context, virtualOrgID string) (*entity.VirtualOrganization, error)
	ListByTenantID(ctx context.Context, tenantID string) ([]*entity.VirtualOrganization, error)
	ListByOwner(ctx context.Context, ownerID string) ([]*entity.VirtualOrganization, error)
	SearchByType(ctx context.Context, tenantID string, orgType entity.VirtualOrgType) ([]*entity.VirtualOrganization, error)
}

// VirtualOrgMemberRepository 虚拟组织成员仓储接口
type VirtualOrgMemberRepository interface {
	AddMember(ctx context.Context, member *entity.VirtualOrgMember) error
	UpdateMemberRole(ctx context.Context, virtualOrgID, userID string, role entity.VirtualOrgMemberRole) error
	RemoveMember(ctx context.Context, virtualOrgID, userID string) error
	GetMember(ctx context.Context, virtualOrgID, userID string) (*entity.VirtualOrgMember, error)
	ListMembers(ctx context.Context, virtualOrgID string) ([]*entity.VirtualOrgMember, error)
	ListUserOrgs(ctx context.Context, userID string) ([]*entity.VirtualOrgMember, error)
}

// VirtualOrgTagRepository 虚拟组织标签仓储接口
type VirtualOrgTagRepository interface {
	AddTag(ctx context.Context, tag *entity.VirtualOrgTag) error
	RemoveTag(ctx context.Context, virtualOrgID, tagName string) error
	ListTags(ctx context.Context, virtualOrgID string) ([]*entity.VirtualOrgTag, error)
	ListByTag(ctx context.Context, tenantID, tagName string) ([]*entity.VirtualOrgTag, error)
}

// MatrixReportingRepository 矩阵汇报关系仓储接口
type MatrixReportingRepository interface {
	Create(ctx context.Context, reporting *entity.MatrixReporting) error
	Update(ctx context.Context, reporting *entity.MatrixReporting) error
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*entity.MatrixReporting, error)
	ListByUserID(ctx context.Context, userID string) ([]*entity.MatrixReporting, error)
	ListActiveByUserID(ctx context.Context, userID string) ([]*entity.MatrixReporting, error)
	ListBySupervisorID(ctx context.Context, supervisorID string) ([]*entity.MatrixReporting, error)
	ListByOrgID(ctx context.Context, orgID string) ([]*entity.MatrixReporting, error)
	GetPrimaryReporting(ctx context.Context, userID string) (*entity.MatrixReporting, error)
}
