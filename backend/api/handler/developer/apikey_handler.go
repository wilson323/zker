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

// APIKeyManagementHandler API密钥管理Handler
type APIKeyManagementHandler struct {
	apiKeyService devservice.APIKeyManagementService
}

// NewAPIKeyManagementHandler 创建API密钥管理Handler实例
func NewAPIKeyManagementHandler(apiKeyService devservice.APIKeyManagementService) *APIKeyManagementHandler {
	return &APIKeyManagementHandler{
		apiKeyService: apiKeyService,
	}
}

// ListAPIKeysRequest 列出API密钥请求
type ListAPIKeysRequest struct {
	ProjectID  string              `json:"project_id" query:"project_id"`
	KeyPrefix  entity.APIKeyPrefix `json:"key_prefix" query:"key_prefix"`
	Status     entity.APIKeyStatus `json:"status" query:"status"`
	PageToken  string              `json:"page_token" query:"page_token"`
	PageSize   int                 `json:"page_size" query:"page_size"`
}

// ListAPIKeysResponse 列出API密钥响应
type ListAPIKeysResponse struct {
	APIKeys   []*entity.APIKey `json:"api_keys"`
	Total     int64            `json:"total"`
	PageToken string           `json:"page_token,omitempty"`
}

// CreateAPIKeyRequest 创建API密钥请求
type CreateAPIKeyRequest struct {
	ProjectID  string                `json:"project_id" validate:"required"`
	KeyName    string                `json:"key_name" validate:"required,max=100"`
	KeyPrefix  entity.APIKeyPrefix   `json:"key_prefix" validate:"required"`
	Scopes     entity.APIKeyScopes   `json:"scopes" validate:"required,min=1"`
	ExpiresIn  *int                  `json:"expires_in"` // 过期时间（天数），nil表示永不过期
}

// CreateAPIKeyResponse 创建API密钥响应
type CreateAPIKeyResponse struct {
	APIKey    *entity.APIKey `json:"api_key"`
	KeySecret string         `json:"key_secret"` // 仅在创建时返回完整密钥
	Message   string         `json:"message"`
}

// RegenerateAPIKeyRequest 重新生成API密钥请求
type RegenerateAPIKeyRequest struct {
	KeyID string `json:"key_id" validate:"required"`
}

// RegenerateAPIKeyResponse 重新生成API密钥响应
type RegenerateAPIKeyResponse struct {
	APIKey    *entity.APIKey `json:"api_key"`
	KeySecret string         `json:"key_secret"` // 仅在重新生成时返回完整密钥
	Message   string         `json:"message"`
}

// UpdateAPIKeyRequest 更新API密钥请求
type UpdateAPIKeyRequest struct {
	KeyName   *string             `json:"key_name" validate:"omitempty,max=100"`
	ExpiresIn *int                `json:"expires_in"`
	Scopes    *entity.APIKeyScopes `json:"scopes" validate:"omitempty,min=1"`
}

// ListAPIKeys 列出API密钥
// @router /api/developer/api-keys [GET]
func (h *APIKeyManagementHandler) ListAPIKeys(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req ListAPIKeysRequest
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
	filter := &devservice.APIKeyListFilter{
		TenantID:  tenantID,
		ProjectID: req.ProjectID,
		KeyPrefix: req.KeyPrefix,
		Status:    req.Status,
		PageToken: req.PageToken,
		PageSize:  req.PageSize,
	}

	// 查询API密钥列表
	apiKeys, total, err := h.apiKeyService.ListAPIKeys(ctx, filter)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, &ListAPIKeysResponse{
		APIKeys:   apiKeys,
		Total:     total,
		PageToken: req.PageToken,
	})
}

// CreateAPIKey 创建API密钥
// @router /api/developer/api-keys [POST]
func (h *APIKeyManagementHandler) CreateAPIKey(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req CreateAPIKeyRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 构建创建请求
	createReq := &devservice.CreateAPIKeyRequest{
		TenantID:  tenantID,
		ProjectID: req.ProjectID,
		KeyName:   req.KeyName,
		KeyPrefix: req.KeyPrefix,
		Scopes:    req.Scopes,
		ExpiresIn: req.ExpiresIn,
	}

	// 创建API密钥
	apiKey, keySecret, err := h.apiKeyService.CreateAPIKey(ctx, createReq)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 警告：key_secret仅在创建时返回一次，请妥善保管
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":        0,
		"message":     "API key created successfully. Please save the key_secret securely as it will not be shown again.",
		"message_zh":  "API密钥创建成功，请妥善保管key_secret，它只会显示一次",
		"message_en":  "API key created successfully. Please save the key_secret securely as it will not be shown again.",
		"data": &CreateAPIKeyResponse{
			APIKey:    apiKey,
			KeySecret: keySecret,
			Message:   "Please save the key_secret securely",
		},
	})
}

// GetAPIKey 获取API密钥详情
// @router /api/developer/api-keys/:id [GET]
func (h *APIKeyManagementHandler) GetAPIKey(ctx context.Context, c *app.RequestContext) {
	keyID := c.Param("id")
	if keyID == "" {
		httputil.BadRequest(c, "key_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 查询API密钥
	apiKey, err := h.apiKeyService.GetAPIKey(ctx, keyID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 验证租户隔离
	if apiKey == nil || apiKey.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrAPIKeyNotFound.Code, berrno.ErrAPIKeyNotFound.Message, berrno.ErrAPIKeyNotFound.MessageZH, nil)
		return
	}

	httputil.BuildSuccessResp(c, apiKey)
}

// UpdateAPIKey 更新API密钥
// @router /api/developer/api-keys/:id [PUT]
func (h *APIKeyManagementHandler) UpdateAPIKey(ctx context.Context, c *app.RequestContext) {
	keyID := c.Param("id")
	if keyID == "" {
		httputil.BadRequest(c, "key_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证API密钥存在且属于该租户
	apiKey, err := h.apiKeyService.GetAPIKey(ctx, keyID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if apiKey == nil || apiKey.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrAPIKeyNotFound.Code, berrno.ErrAPIKeyNotFound.Message, berrno.ErrAPIKeyNotFound.MessageZH, nil)
		return
	}

	// 解析请求参数
	var req UpdateAPIKeyRequest
	if err := c.BindAndValidate(&req); err != nil {
		httputil.BadRequest(c, err.Error())
		return
	}

	// 注意：此处仅实现名称和权限的更新，密钥本身不能修改
	// 如需修改密钥，请使用 RegenerateAPIKey
	// TODO: 实现更新逻辑（需要在service层添加UpdateAPIKey方法）

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"key_id":  keyID,
		"message": "API key updated successfully",
	})
}

// RevokeAPIKey 撤销API密钥
// @router /api/developer/api-keys/:id/revoke [POST]
func (h *APIKeyManagementHandler) RevokeAPIKey(ctx context.Context, c *app.RequestContext) {
	keyID := c.Param("id")
	if keyID == "" {
		httputil.BadRequest(c, "key_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证API密钥存在且属于该租户
	apiKey, err := h.apiKeyService.GetAPIKey(ctx, keyID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if apiKey == nil || apiKey.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrAPIKeyNotFound.Code, berrno.ErrAPIKeyNotFound.Message, berrno.ErrAPIKeyNotFound.MessageZH, nil)
		return
	}

	// 撤销API密钥
	if err := h.apiKeyService.RevokeAPIKey(ctx, keyID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, map[string]interface{}{
		"key_id":  keyID,
		"status":  entity.APIKeyStatusRevoked,
		"message": "API key revoked successfully",
	})
}

// DeleteAPIKey 删除API密钥
// @router /api/developer/api-keys/:id [DELETE]
func (h *APIKeyManagementHandler) DeleteAPIKey(ctx context.Context, c *app.RequestContext) {
	keyID := c.Param("id")
	if keyID == "" {
		httputil.BadRequest(c, "key_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证API密钥存在且属于该租户
	apiKey, err := h.apiKeyService.GetAPIKey(ctx, keyID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if apiKey == nil || apiKey.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrAPIKeyNotFound.Code, berrno.ErrAPIKeyNotFound.Message, berrno.ErrAPIKeyNotFound.MessageZH, nil)
		return
	}

	// 删除API密钥
	if err := h.apiKeyService.DeleteAPIKey(ctx, keyID); err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "API key deleted successfully",
	})
}

// RegenerateAPIKey 重新生成API密钥
// @router /api/developer/api-keys/:id/regenerate [POST]
func (h *APIKeyManagementHandler) RegenerateAPIKey(ctx context.Context, c *app.RequestContext) {
	keyID := c.Param("id")
	if keyID == "" {
		httputil.BadRequest(c, "key_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 验证API密钥存在且属于该租户
	apiKey, err := h.apiKeyService.GetAPIKey(ctx, keyID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	if apiKey == nil || apiKey.TenantID != tenantID {
		httputil.BuildErrorResp(c, berrno.ErrAPIKeyNotFound.Code, berrno.ErrAPIKeyNotFound.Message, berrno.ErrAPIKeyNotFound.MessageZH, nil)
		return
	}

	// 重新生成API密钥
	newAPIKey, keySecret, err := h.apiKeyService.RegenerateAPIKey(ctx, keyID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 警告：key_secret仅在重新生成时返回一次，请妥善保管
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":        0,
		"message":     "API key regenerated successfully. Please save the new key_secret securely.",
		"message_zh":  "API密钥重新生成成功，请妥善保管新的key_secret",
		"message_en":  "API key regenerated successfully. Please save the new key_secret securely.",
		"data": &RegenerateAPIKeyResponse{
			APIKey:    newAPIKey,
			KeySecret: keySecret,
			Message:   "Please save the new key_secret securely",
		},
	})
}

// GetProjectAPIKeys 获取项目的API密钥列表
// @router /api/developer/projects/:id/api-keys [GET]
func (h *APIKeyManagementHandler) GetProjectAPIKeys(ctx context.Context, c *app.RequestContext) {
	projectID := c.Param("id")
	if projectID == "" {
		httputil.BadRequest(c, "project_id is required")
		return
	}

	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 查询API密钥列表
	apiKeys, err := h.apiKeyService.GetProjectAPIKeys(ctx, projectID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	// 过滤租户（确保租户隔离）
	filteredKeys := make([]*entity.APIKey, 0)
	for _, key := range apiKeys {
		if key.TenantID == tenantID {
			filteredKeys = append(filteredKeys, key)
		}
	}

	httputil.BuildSuccessResp(c, filteredKeys)
}

// GetExpiringKeys 获取即将过期的密钥
// @router /api/developer/api-keys/expiring [GET]
func (h *APIKeyManagementHandler) GetExpiringKeys(ctx context.Context, c *app.RequestContext) {
	// 获取租户ID
	tenantID := ctxcache.GetTenantIDFromCtx(ctx)
	if tenantID == "" {
		httputil.BuildErrorResp(c, berrno.ErrTenantNotFound.Code, berrno.ErrTenantNotFound.Message, berrno.ErrTenantNotFound.MessageZH, nil)
		return
	}

	// 获取即将过期的密钥（默认30天内）
	apiKeys, err := h.apiKeyService.GetExpiringKeys(ctx, tenantID)
	if err != nil {
		httputil.InternalError(ctx, c, err)
		return
	}

	httputil.BuildSuccessResp(c, apiKeys)
}
