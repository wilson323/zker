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

package org

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/domain/org/service"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// EnhancedOrgHandler 增强组织Handler
type EnhancedOrgHandler struct {
	profileService    *service.MemberProfileService
	virtualOrgService *service.VirtualOrganizationService
	matrixOrgService  *service.MatrixOrganizationService
}

// NewEnhancedOrgHandler 创建增强组织Handler实例
func NewEnhancedOrgHandler(
	profileService *service.MemberProfileService,
	virtualOrgService *service.VirtualOrganizationService,
	matrixOrgService *service.MatrixOrganizationService,
) *EnhancedOrgHandler {
	return &EnhancedOrgHandler{
		profileService:    profileService,
		virtualOrgService: virtualOrgService,
		matrixOrgService:  matrixOrgService,
	}
}

// ==================== 员工档案API ====================

// AddWorkExperience 添加工作经历
// @Router /api/v1/org/employees/:user_id/work_experiences [POST]
func (h *EnhancedOrgHandler) AddWorkExperience(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	var req entity.WorkExperience
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	req.UserID = userID
	if err := h.profileService.AddWorkExperience(ctx, &req); err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"message": "success"})
}

// ListWorkExperiences 列出工作经历
// @Router /api/v1/org/employees/:user_id/work_experiences [GET]
func (h *EnhancedOrgHandler) ListWorkExperiences(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	exps, err := h.profileService.ListWorkExperiences(ctx, userID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"work_experiences": exps})
}

// AddEducation 添加教育经历
// @Router /api/v1/org/employees/:user_id/educations [POST]
func (h *EnhancedOrgHandler) AddEducation(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	var req entity.Education
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	req.UserID = userID
	if err := h.profileService.AddEducation(ctx, &req); err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"message": "success"})
}

// ListEducations 列出教育经历
// @Router /api/v1/org/employees/:user_id/educations [GET]
func (h *EnhancedOrgHandler) ListEducations(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	edus, err := h.profileService.ListEducations(ctx, userID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"educations": edus})
}

// AddSkill 添加技能
// @Router /api/v1/org/employees/:user_id/skills [POST]
func (h *EnhancedOrgHandler) AddSkill(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	var req entity.MemberSkill
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	req.UserID = userID
	if err := h.profileService.AddSkill(ctx, &req); err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"message": "success"})
}

// ListSkills 列出技能
// @Router /api/v1/org/employees/:user_id/skills [GET]
func (h *EnhancedOrgHandler) ListSkills(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	skills, err := h.profileService.ListSkills(ctx, userID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"skills": skills})
}

// GenerateResume 生成简历
// @Router /api/v1/org/employees/:user_id/resume [GET]
func (h *EnhancedOrgHandler) GenerateResume(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	resume, err := h.profileService.GenerateResume(ctx, userID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, resume)
}

// ==================== 虚拟组织API ====================

// CreateVirtualOrg 创建虚拟组织
// @Router /api/v1/org/virtual-organizations [POST]
func (h *EnhancedOrgHandler) CreateVirtualOrg(ctx context.Context, c *app.RequestContext) {
	var req entity.VirtualOrganization
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	if err := h.virtualOrgService.CreateVirtualOrg(ctx, &req); err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantCreated, M{
		"virtual_org_id": req.VirtualOrgID,
		"message":        "success",
	})
}

// ListVirtualOrgs 列出虚拟组织
// @Router /api/v1/org/virtual-organizations [GET]
func (h *EnhancedOrgHandler) ListVirtualOrgs(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	orgType := c.Query("type")

	var orgTypePtr *entity.VirtualOrgType
	if orgType != "" {
		t := entity.VirtualOrgType(orgType)
		orgTypePtr = &t
	}

	orgs, err := h.virtualOrgService.ListVirtualOrgs(ctx, tenantID, orgTypePtr)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"virtual_organizations": orgs})
}

// AddVirtualMember 添加虚拟组织成员
// @Router /api/v1/org/virtual-organizations/:org_id/members [POST]
func (h *EnhancedOrgHandler) AddVirtualMember(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	role := entity.VirtualOrgMemberRole(req.Role)
	if err := h.virtualOrgService.AddMember(ctx, orgID, req.UserID, role); err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"message": "success"})
}

// ListVirtualMembers 列出虚拟组织成员
// @Router /api/v1/org/virtual-organizations/:org_id/members [GET]
func (h *EnhancedOrgHandler) ListVirtualMembers(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	members, err := h.virtualOrgService.ListMembers(ctx, orgID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"members": members})
}

// ==================== 矩阵组织API ====================

// CreateReporting 创建矩阵汇报关系
// @Router /api/v1/org/matrix-reportings [POST]
func (h *EnhancedOrgHandler) CreateReporting(ctx context.Context, c *app.RequestContext) {
	var req entity.MatrixReporting
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	if err := h.matrixOrgService.CreateReporting(ctx, &req); err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantCreated, M{
		"id":      req.ID,
		"message": "success",
	})
}

// GetSupervisors 获取所有上级
// @Router /api/v1/org/employees/:user_id/supervisors [GET]
func (h *EnhancedOrgHandler) GetSupervisors(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	supervisors, err := h.matrixOrgService.GetSupervisors(ctx, userID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"supervisors": supervisors})
}

// GetSubordinates 获取下属
// @Router /api/v1/org/employees/:user_id/subordinates [GET]
func (h *EnhancedOrgHandler) GetSubordinates(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	includeMatrixStr := c.Query("include_matrix")
	includeMatrix := includeMatrixStr == "true"

	subordinates, err := h.matrixOrgService.GetSubordinates(ctx, userID, includeMatrix)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, M{"subordinates": subordinates})
}

// GenerateOrgChart 生成组织架构图
// @Router /api/v1/org/organizations/:org_id/chart [GET]
func (h *EnhancedOrgHandler) GenerateOrgChart(ctx context.Context, c *app.RequestContext) {
	orgID := c.Param("org_id")
	if orgID == "" {
		c.JSON(constantBadRequest, errno.NewError(errno.ErrInvalidParam))
		return
	}

	orgChart, err := h.matrixOrgService.GenerateOrgChart(ctx, orgID)
	if err != nil {
		c.JSON(constantInternalServerError, errno.NewError(err))
		return
	}

	c.JSON(constantOK, orgChart)
}

// 辅助常量
const (
	constantOK                  = http.StatusOK
	constantBadRequest          = http.StatusBadRequest
	constantCreated             = http.StatusCreated
	constantInternalServerError = http.StatusInternalServerError
)

// M 简化的map类型
type M map[string]interface{}

// MarshalJSON 自定义JSON序列化
func (m M) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(m))
}
