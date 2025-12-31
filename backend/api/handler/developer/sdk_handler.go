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

package developer

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	devservice "github.com/coze-dev/coze-studio/backend/domain/developer/service"
	"github.com/coze-dev/coze-studio/backend/domain/developer/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
)

// SDKGeneratorHandler SDK生成Handler
type SDKGeneratorHandler struct {
	sdkService devservice.SDKGeneratorService
}

// NewSDKGeneratorHandler 创建SDK生成Handler实例
func NewSDKGeneratorHandler(sdkService devservice.SDKGeneratorService) *SDKGeneratorHandler {
	return &SDKGeneratorHandler{
		sdkService: sdkService,
	}
}

// ListSDKsRequest 列出SDK请求
type ListSDKsRequest struct {
	ProjectID  string            `json:"project_id" query:"project_id"`
	Language   entity.SDKLanguage `json:"language" query:"language"`
	Status     entity.SDKStatus  `json:"status" query:"status"`
	PageToken  string            `json:"page_token" query:"page_token"`
	PageSize   int               `json:"page_size" query:"page_size"`
}

// ListSDKsResponse 列出SDK响应
type ListSDKsResponse struct {
	SDKs      []*entity.SDK `json:"sdks"`
	Total     int64         `json:"total"`
	PageToken string        `json:"page_token,omitempty"`
}

// GenerateSDKRequest 生成SDK请求
type GenerateSDKRequest struct {
	ProjectID   string             `json:"project_id" validate:"required"`
	SDKName     string             `json:"sdk_name" validate:"required,max=100"`
	Language    entity.SDKLanguage `json:"language" validate:"required"`
	Version     string             `json:"version" validate:"required"`
	Description string             `json:"description"`
}

// PublishSDKRequest 发布SDK请求
type PublishSDKRequest struct {
	SDKID string `json:"sdk_id" validate:"required"`
}

// ListSDKs 列出可用SDK
// @router /api/developer/sdks [GET]
func (h *SDKGeneratorHandler) ListSDKs(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req ListSDKsRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 设置默认分页参数
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 构建过滤器
	filter := &devservice.SDKListFilter{
		TenantID:  tenantID,
		ProjectID: req.ProjectID,
		Language:  req.Language,
		Status:    req.Status,
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	// 查询SDK列表
	sdks, total, err := h.sdkService.ListSDKs(ctx, filter)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, &ListSDKsResponse{
		SDKs:      sdks,
		Total:     total,
		PageToken: req.PageToken,
	})
}

// GenerateSDK 生成SDK
// @router /api/developer/sdks/generate [POST]
func (h *SDKGeneratorHandler) GenerateSDK(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 解析请求参数
	var req GenerateSDKRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 验证语言是否支持
	if !h.isValidLanguage(req.Language) {
		httputil.BuildErrorResp(c, berrno.ErrSDKLanguageInvalid.Int32Code(), berrno.ErrSDKLanguageInvalid.Message(), berrno.ErrSDKLanguageInvalid.MessageZH(), nil)
		return
	}

	// 构建生成请求
	generateReq := &devservice.GenerateSDKRequest{
		TenantID:    tenantID,
		ProjectID:   req.ProjectID,
		SDKName:     req.SDKName,
		Language:    req.Language,
		Version:     req.Version,
		Description: req.Description,
	}

	// 生成SDK
	sdk, err := h.sdkService.GenerateSDK(ctx, generateReq)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, sdk)
}

// GetSDK 获取SDK详情
// @router /api/developer/sdks/:id [GET]
func (h *SDKGeneratorHandler) GetSDK(ctx context.Context, c *app.RequestContext) {
	sdkID := c.Param("id")
	if sdkID == "" {
		httputil.BadRequest(c, "sdk_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询SDK
	sdk, err := h.sdkService.GetSDK(ctx, sdkID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 验证租户隔离
	if sdk == nil || sdk.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrSDKNotFound.Int32Code(), berrno.ErrSDKNotFound.Message(), berrno.ErrSDKNotFound.MessageZH(), nil)
		return
	}

	httputil.BuildSuccessResp(c, sdk)
}

// PublishSDK 发布SDK
// @router /api/developer/sdks/:id/publish [POST]
func (h *SDKGeneratorHandler) PublishSDK(ctx context.Context, c *app.RequestContext) {
	sdkID := c.Param("id")
	if sdkID == "" {
		httputil.BadRequest(c, "sdk_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证SDK存在且属于该租户
	sdk, err := h.sdkService.GetSDK(ctx, sdkID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if sdk == nil || sdk.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrSDKNotFound.Int32Code(), berrno.ErrSDKNotFound.Message(), berrno.ErrSDKNotFound.MessageZH(), nil)
		return
	}

	// 发布SDK
	if err := h.sdkService.PublishSDK(ctx, sdkID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"sdk_id":  sdkID,
		"status":  entity.SDKStatusPublished,
		"message": "SDK published successfully",
	})
}

// DeleteSDK 删除SDK
// @router /api/developer/sdks/:id [DELETE]
func (h *SDKGeneratorHandler) DeleteSDK(ctx context.Context, c *app.RequestContext) {
	sdkID := c.Param("id")
	if sdkID == "" {
		httputil.BadRequest(c, "sdk_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证SDK存在且属于该租户
	sdk, err := h.sdkService.GetSDK(ctx, sdkID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if sdk == nil || sdk.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrSDKNotFound.Int32Code(), berrno.ErrSDKNotFound.Message(), berrno.ErrSDKNotFound.MessageZH(), nil)
		return
	}

	// 检查SDK是否已发布（已发布的SDK不能删除，只能归档）
	if sdk.IsPublished() {
		httputil.BuildErrorResp(c, berrno.ErrSDKNotFound.Int32Code(), "Published SDK cannot be deleted", "已发布的SDK不能删除", map[string]interface{}{
			"sdk_id":    sdkID,
			"status":    sdk.Status,
			"suggestion": "Use archive operation instead",
		})
		return
	}

	// 删除SDK
	if err := h.sdkService.DeleteSDK(ctx, sdkID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "SDK deleted successfully",
	})
}

// DownloadSDK 下载SDK
// @router /api/developer/sdks/:id/download [GET]
func (h *SDKGeneratorHandler) DownloadSDK(ctx context.Context, c *app.RequestContext) {
	sdkID := c.Param("id")
	if sdkID == "" {
		httputil.BadRequest(c, "sdk_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 验证SDK存在且属于该租户
	sdk, err := h.sdkService.GetSDK(ctx, sdkID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if sdk == nil || sdk.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrSDKNotFound.Int32Code(), berrno.ErrSDKNotFound.Message(), berrno.ErrSDKNotFound.MessageZH(), nil)
		return
	}

	// 下载SDK包
	pkg, err := h.sdkService.DownloadSDK(ctx, sdkID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, pkg)
}

// GetSDKCode 获取SDK代码（预览）
// @router /api/developer/sdks/:id/code [GET]
func (h *SDKGeneratorHandler) GetSDKCode(ctx context.Context, c *app.RequestContext) {
	sdkID := c.Param("id")
	if sdkID == "" {
		httputil.BadRequest(c, "sdk_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询SDK
	sdk, err := h.sdkService.GetSDK(ctx, sdkID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 验证租户隔离
	if sdk == nil || sdk.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrSDKNotFound.Int32Code(), berrno.ErrSDKNotFound.Message(), berrno.ErrSDKNotFound.MessageZH(), nil)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"sdk_id":     sdkID,
		"language":   sdk.Language,
		"code":       sdk.Code,
		"readme":     sdk.Readme,
		"line_count": len([]rune(sdk.Code)),
	})
}

// GetProjectSDKs 获取项目的SDK列表
// @router /api/developer/projects/:id/sdks [GET]
func (h *SDKGeneratorHandler) GetProjectSDKs(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Int32Code(), berrno.ErrTenantNotFound.Message(), berrno.ErrTenantNotFound.MessageZH(), nil)
		return
	}

	// 查询SDK列表
	sdks, err := h.sdkService.GetProjectSDKs(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 过滤租户（确保租户隔离）
	filteredSDKs := make([]*entity.SDK, 0)
	for _, sdk := range sdks {
		if sdk.TenantID == tenantID {
			filteredSDKs = append(filteredSDKs, sdk)
		}
	}

	httputil.BuildSuccessResp(c, filteredSDKs)
}

// GetSupportedLanguages 获取支持的SDK语言列表
// @router /api/developer/sdks/languages [GET]
func (h *SDKGeneratorHandler) GetSupportedLanguages(ctx context.Context, c *app.RequestContext) {
	languages := []map[string]interface{}{
		{
			"language":     entity.SDKLanguagePython,
			"display_name": "Python",
			"display_name_zh": "Python",
			"description":  "Python 3.7+ SDK with async support",
			"icon":         "python",
		},
		{
			"language":     entity.SDKLanguageJavaScript,
			"display_name": "TypeScript",
			"display_name_zh": "TypeScript",
			"description":  "TypeScript/JavaScript SDK with full type definitions",
			"icon":         "javascript",
		},
		{
			"language":     entity.SDKLanguageGo,
			"display_name": "Go",
			"display_name_zh": "Go",
			"description":  "Go SDK with high performance",
			"icon":         "go",
		},
		{
			"language":     entity.SDKLanguageJava,
			"display_name": "Java",
			"display_name_zh": "Java",
			"description":  "Java SDK for enterprise applications",
			"icon":         "java",
		},
	}

	httputil.BuildSuccessResp(c, languages)
}

// isValidLanguage 验证语言是否支持
func (h *SDKGeneratorHandler) isValidLanguage(language entity.SDKLanguage) bool {
	switch language {
	case entity.SDKLanguagePython, entity.SDKLanguageJavaScript, entity.SDKLanguageGo, entity.SDKLanguageJava:
		return true
	default:
		return false
	}
}
