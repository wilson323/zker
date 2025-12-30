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

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
)

// workExperienceRepository 工作经历仓储实现
type workExperienceRepository struct {
	db *gorm.DB
}

// NewWorkExperienceRepository 创建工作经历仓储
func NewWorkExperienceRepository(db *gorm.DB) WorkExperienceRepository {
	return &workExperienceRepository{db: db}
}

func (r *workExperienceRepository) Create(ctx context.Context, exp *entity.WorkExperience) error {
	return r.db.WithContext(ctx).Create(exp).Error
}

func (r *workExperienceRepository) Update(ctx context.Context, exp *entity.WorkExperience) error {
	return r.db.WithContext(ctx).Model(&entity.WorkExperience{}).
		Where("id = ?", exp.ID).
		Updates(exp).Error
}

func (r *workExperienceRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *workExperienceRepository) GetByID(ctx context.Context, id int64) (*entity.WorkExperience, error) {
	var exp entity.WorkExperience
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&exp).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *workExperienceRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.WorkExperience, error) {
	var exps []*entity.WorkExperience
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("start_date DESC").
		Find(&exps).Error
	return exps, err
}

// educationRepository 教育经历仓储实现
type educationRepository struct {
	db *gorm.DB
}

// NewEducationRepository 创建教育经历仓储
func NewEducationRepository(db *gorm.DB) EducationRepository {
	return &educationRepository{db: db}
}

func (r *educationRepository) Create(ctx context.Context, edu *entity.Education) error {
	return r.db.WithContext(ctx).Create(edu).Error
}

func (r *educationRepository) Update(ctx context.Context, edu *entity.Education) error {
	return r.db.WithContext(ctx).Model(&entity.Education{}).
		Where("id = ?", edu.ID).
		Updates(edu).Error
}

func (r *educationRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *educationRepository) GetByID(ctx context.Context, id int64) (*entity.Education, error) {
	var edu entity.Education
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&edu).Error
	if err != nil {
		return nil, err
	}
	return &edu, nil
}

func (r *educationRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Education, error) {
	var edus []*entity.Education
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("start_date DESC").
		Find(&edus).Error
	return edus, err
}

// memberSkillRepository 技能仓储实现
type memberSkillRepository struct {
	db *gorm.DB
}

// NewMemberSkillRepository 创建技能仓储
func NewMemberSkillRepository(db *gorm.DB) MemberSkillRepository {
	return &memberSkillRepository{db: db}
}

func (r *memberSkillRepository) Create(ctx context.Context, skill *entity.MemberSkill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *memberSkillRepository) Update(ctx context.Context, skill *entity.MemberSkill) error {
	return r.db.WithContext(ctx).Model(&entity.MemberSkill{}).
		Where("id = ?", skill.ID).
		Updates(skill).Error
}

func (r *memberSkillRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *memberSkillRepository) GetByID(ctx context.Context, id int64) (*entity.MemberSkill, error) {
	var skill entity.MemberSkill
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&skill).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *memberSkillRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.MemberSkill, error) {
	var skills []*entity.MemberSkill
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("proficiency DESC, skill_name ASC").
		Find(&skills).Error
	return skills, err
}

func (r *memberSkillRepository) ListByTenantID(ctx context.Context, tenantID string) ([]*entity.MemberSkill, error) {
	var skills []*entity.MemberSkill
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("proficiency DESC, skill_name ASC").
		Find(&skills).Error
	return skills, err
}

func (r *memberSkillRepository) SearchByName(ctx context.Context, tenantID, skillName string, limit int) ([]*entity.MemberSkill, error) {
	var skills []*entity.MemberSkill
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Where("skill_name LIKE ?", "%"+skillName+"%")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("proficiency DESC").Find(&skills).Error
	return skills, err
}

// virtualOrganizationRepository 虚拟组织仓储实现
type virtualOrganizationRepository struct {
	db *gorm.DB
}

// NewVirtualOrganizationRepository 创建虚拟组织仓储
func NewVirtualOrganizationRepository(db *gorm.DB) VirtualOrganizationRepository {
	return &virtualOrganizationRepository{db: db}
}

func (r *virtualOrganizationRepository) Create(ctx context.Context, org *entity.VirtualOrganization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *virtualOrganizationRepository) Update(ctx context.Context, org *entity.VirtualOrganization) error {
	return r.db.WithContext(ctx).Model(&entity.VirtualOrganization{}).
		Where("virtual_org_id = ?", org.VirtualOrgID).
		Updates(org).Error
}

func (r *virtualOrganizationRepository) Delete(ctx context.Context, virtualOrgID string) error {
	return r.db.WithContext(ctx).
		Where("virtual_org_id = ?", virtualOrgID).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *virtualOrganizationRepository) GetByID(ctx context.Context, virtualOrgID string) (*entity.VirtualOrganization, error) {
	var org entity.VirtualOrganization
	err := r.db.WithContext(ctx).
		Preload("Members").
		Preload("Tags").
		Where("virtual_org_id = ? AND deleted_at IS NULL", virtualOrgID).
		First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *virtualOrganizationRepository) ListByTenantID(ctx context.Context, tenantID string) ([]*entity.VirtualOrganization, error) {
	var orgs []*entity.VirtualOrganization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&orgs).Error
	return orgs, err
}

func (r *virtualOrganizationRepository) ListByOwner(ctx context.Context, ownerID string) ([]*entity.VirtualOrganization, error) {
	var orgs []*entity.VirtualOrganization
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND deleted_at IS NULL", ownerID).
		Order("created_at DESC").
		Find(&orgs).Error
	return orgs, err
}

func (r *virtualOrganizationRepository) SearchByType(ctx context.Context, tenantID string, orgType entity.VirtualOrgType) ([]*entity.VirtualOrganization, error) {
	var orgs []*entity.VirtualOrganization
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND type = ? AND deleted_at IS NULL", tenantID, orgType).
		Order("created_at DESC").
		Find(&orgs).Error
	return orgs, err
}

// virtualOrgMemberRepository 虚拟组织成员仓储实现
type virtualOrgMemberRepository struct {
	db *gorm.DB
}

// NewVirtualOrgMemberRepository 创建虚拟组织成员仓储
func NewVirtualOrgMemberRepository(db *gorm.DB) VirtualOrgMemberRepository {
	return &virtualOrgMemberRepository{db: db}
}

func (r *virtualOrgMemberRepository) AddMember(ctx context.Context, member *entity.VirtualOrgMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *virtualOrgMemberRepository) UpdateMemberRole(ctx context.Context, virtualOrgID, userID string, role entity.VirtualOrgMemberRole) error {
	return r.db.WithContext(ctx).Model(&entity.VirtualOrgMember{}).
		Where("virtual_org_id = ? AND user_id = ?", virtualOrgID, userID).
		Update("role", role).Error
}

func (r *virtualOrgMemberRepository) RemoveMember(ctx context.Context, virtualOrgID, userID string) error {
	return r.db.WithContext(ctx).
		Where("virtual_org_id = ? AND user_id = ?", virtualOrgID, userID).
		Delete(&entity.VirtualOrgMember{}).Error
}

func (r *virtualOrgMemberRepository) GetMember(ctx context.Context, virtualOrgID, userID string) (*entity.VirtualOrgMember, error) {
	var member entity.VirtualOrgMember
	err := r.db.WithContext(ctx).
		Where("virtual_org_id = ? AND user_id = ?", virtualOrgID, userID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *virtualOrgMemberRepository) ListMembers(ctx context.Context, virtualOrgID string) ([]*entity.VirtualOrgMember, error) {
	var members []*entity.VirtualOrgMember
	err := r.db.WithContext(ctx).
		Where("virtual_org_id = ?", virtualOrgID).
		Order("role ASC, joined_at ASC").
		Find(&members).Error
	return members, err
}

func (r *virtualOrgMemberRepository) ListUserOrgs(ctx context.Context, userID string) ([]*entity.VirtualOrgMember, error) {
	var memberships []*entity.VirtualOrgMember
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("joined_at DESC").
		Find(&memberships).Error
	return memberships, err
}

// virtualOrgTagRepository 虚拟组织标签仓储实现
type virtualOrgTagRepository struct {
	db *gorm.DB
}

// NewVirtualOrgTagRepository 创建虚拟组织标签仓储
func NewVirtualOrgTagRepository(db *gorm.DB) VirtualOrgTagRepository {
	return &virtualOrgTagRepository{db: db}
}

func (r *virtualOrgTagRepository) AddTag(ctx context.Context, tag *entity.VirtualOrgTag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *virtualOrgTagRepository) RemoveTag(ctx context.Context, virtualOrgID, tagName string) error {
	return r.db.WithContext(ctx).
		Where("virtual_org_id = ? AND tag_name = ?", virtualOrgID, tagName).
		Delete(&entity.VirtualOrgTag{}).Error
}

func (r *virtualOrgTagRepository) ListTags(ctx context.Context, virtualOrgID string) ([]*entity.VirtualOrgTag, error) {
	var tags []*entity.VirtualOrgTag
	err := r.db.WithContext(ctx).
		Where("virtual_org_id = ?", virtualOrgID).
		Order("tag_name ASC").
		Find(&tags).Error
	return tags, err
}

func (r *virtualOrgTagRepository) ListByTag(ctx context.Context, tenantID, tagName string) ([]*entity.VirtualOrgTag, error) {
	var tags []*entity.VirtualOrgTag
	err := r.db.WithContext(ctx).
		Joins("JOIN virtual_organizations ON virtual_organizations.virtual_org_id = virtual_org_tags.virtual_org_id").
		Where("virtual_organizations.tenant_id = ? AND virtual_org_tags.tag_name = ?", tenantID, tagName).
		Find(&tags).Error
	return tags, err
}

// matrixReportingRepository 矩阵汇报关系仓储实现
type matrixReportingRepository struct {
	db *gorm.DB
}

// NewMatrixReportingRepository 创建矩阵汇报关系仓储
func NewMatrixReportingRepository(db *gorm.DB) MatrixReportingRepository {
	return &matrixReportingRepository{db: db}
}

func (r *matrixReportingRepository) Create(ctx context.Context, reporting *entity.MatrixReporting) error {
	return r.db.WithContext(ctx).Create(reporting).Error
}

func (r *matrixReportingRepository) Update(ctx context.Context, reporting *entity.MatrixReporting) error {
	return r.db.WithContext(ctx).Model(&entity.MatrixReporting{}).
		Where("id = ?", reporting.ID).
		Updates(reporting).Error
}

func (r *matrixReportingRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.MatrixReporting{}).Error
}

func (r *matrixReportingRepository) GetByID(ctx context.Context, id int64) (*entity.MatrixReporting, error) {
	var reporting entity.MatrixReporting
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Supervisor").
		Preload("Organization").
		Where("id = ?", id).
		First(&reporting).Error
	if err != nil {
		return nil, err
	}
	return &reporting, nil
}

func (r *matrixReportingRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.MatrixReporting, error) {
	var reportings []*entity.MatrixReporting
	err := r.db.WithContext(ctx).
		Preload("Supervisor").
		Preload("Organization").
		Where("user_id = ?", userID).
		Order("is_primary DESC, effective_date DESC").
		Find(&reportings).Error
	return reportings, err
}

func (r *matrixReportingRepository) ListActiveByUserID(ctx context.Context, userID string) ([]*entity.MatrixReporting, error) {
	var reportings []*entity.MatrixReporting
	err := r.db.WithContext(ctx).
		Preload("Supervisor").
		Preload("Organization").
		Where("user_id = ?", userID).
		Where("(expiry_date IS NULL OR expiry_date > NOW())").
		Where("effective_date <= NOW()").
		Order("is_primary DESC, effective_date DESC").
		Find(&reportings).Error
	return reportings, err
}

func (r *matrixReportingRepository) ListBySupervisorID(ctx context.Context, supervisorID string) ([]*entity.MatrixReporting, error) {
	var reportings []*entity.MatrixReporting
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Organization").
		Where("supervisor_id = ?", supervisorID).
		Where("(expiry_date IS NULL OR expiry_date > NOW())").
		Where("effective_date <= NOW()").
		Order("is_primary DESC, effective_date DESC").
		Find(&reportings).Error
	return reportings, err
}

func (r *matrixReportingRepository) ListByOrgID(ctx context.Context, orgID string) ([]*entity.MatrixReporting, error) {
	var reportings []*entity.MatrixReporting
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Supervisor").
		Where("organization_id = ?", orgID).
		Where("(expiry_date IS NULL OR expiry_date > NOW())").
		Where("effective_date <= NOW()").
		Order("is_primary DESC, effective_date DESC").
		Find(&reportings).Error
	return reportings, err
}

func (r *matrixReportingRepository) GetPrimaryReporting(ctx context.Context, userID string) (*entity.MatrixReporting, error) {
	var reporting entity.MatrixReporting
	err := r.db.WithContext(ctx).
		Preload("Supervisor").
		Preload("Organization").
		Where("user_id = ? AND is_primary = ?", userID, true).
		Where("(expiry_date IS NULL OR expiry_date > NOW())").
		Where("effective_date <= NOW()").
		First(&reporting).Error
	if err != nil {
		return nil, err
	}
	return &reporting, nil
}
