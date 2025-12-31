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

package coze

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/model/config"
	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/domain/config/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var configService *service.ConfigService

// InitConfigCenter 初始化配置中心服务
func InitConfigCenter(cfgService *service.ConfigService) {
	configService = cfgService
}

// CreateConfig 创建配置
// @router /api/config/create [POST]
func CreateConfig(ctx context.Context, c *app.RequestContext) {
	var req config.CreateConfigReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	// 获取用户ID（从context或header）
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "system"
	}

	result, err := configService.CreateConfig(
		ctx,
		req.TenantID,
		req.ConfigKey,
		req.ConfigValue,
		req.ConfigType,
		req.Description,
		userID,
	)
	if err != nil {
		logs.CtxErrorf(ctx, "create config failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("create config failed: %v", err))
		return
	}

	resp := &config.CreateConfigResp{
		ConfigID: result.ConfigID,
	}
	httputil.BuildSuccessResp(c, resp)
}

// GetConfig 获取配置
// @router /api/config/get [GET]
func GetConfig(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	configKey := c.Query("config_key")

	if tenantID == "" || configKey == "" {
		c.String(consts.StatusBadRequest, "tenant_id and config_key are required")
		return
	}

	result, err := configService.GetConfig(ctx, tenantID, configKey)
	if err != nil {
		logs.CtxErrorf(ctx, "get config failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("get config failed: %v", err))
		return
	}

	resp := &config.GetConfigResp{
		ConfigID:    result.ConfigID,
		ConfigKey:   result.ConfigKey,
		ConfigValue: result.ConfigValue,
		ConfigType:  result.ConfigType,
		Description: result.Description,
		UpdatedAt:   result.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	httputil.BuildSuccessResp(c, resp)
}

// UpdateConfig 更新配置
// @router /api/config/update [POST]
func UpdateConfig(ctx context.Context, c *app.RequestContext) {
	var req config.UpdateConfigReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	// 获取用户ID
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "system"
	}

	err := configService.UpdateConfig(
		ctx,
		req.TenantID,
		req.ConfigKey,
		req.NewValue,
		req.ChangeReason,
		userID,
	)
	if err != nil {
		logs.CtxErrorf(ctx, "update config failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("update config failed: %v", err))
		return
	}

	resp := &config.UpdateConfigResp{
		Success: true,
	}
	httputil.BuildSuccessResp(c, resp)
}

// DeleteConfig 删除配置
// @router /api/config/delete [POST]
func DeleteConfig(ctx context.Context, c *app.RequestContext) {
	var req config.DeleteConfigReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	err := configService.DeleteConfig(ctx, req.TenantID, req.ConfigKey)
	if err != nil {
		logs.CtxErrorf(ctx, "delete config failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("delete config failed: %v", err))
		return
	}

	resp := &config.DeleteConfigResp{
		Success: true,
	}
	httputil.BuildSuccessResp(c, resp)
}

// ListConfigs 列出租户所有配置
// @router /api/config/list [GET]
func ListConfigs(ctx context.Context, c *app.RequestContext) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.String(consts.StatusBadRequest, "tenant_id is required")
		return
	}

	configs, err := configService.ListConfigs(ctx, tenantID)
	if err != nil {
		logs.CtxErrorf(ctx, "list configs failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("list configs failed: %v", err))
		return
	}

	// 转换为响应格式
	items := make([]config.ConfigItem, 0, len(configs))
	for _, cfg := range configs {
		items = append(items, config.ConfigItem{
			ConfigID:    cfg.ConfigID,
			ConfigKey:   cfg.ConfigKey,
			ConfigValue: cfg.ConfigValue,
			ConfigType:  cfg.ConfigType,
			Description: cfg.Description,
			UpdatedAt:   cfg.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	resp := &config.ListConfigsResp{
		Configs: items,
		Total:   len(items),
	}
	httputil.BuildSuccessResp(c, resp)
}

// GetConfigHistory 获取配置变更历史
// @router /api/config/history [GET]
func GetConfigHistory(ctx context.Context, c *app.RequestContext) {
	configID := c.Query("config_id")
	if configID == "" {
		c.String(consts.StatusBadRequest, "config_id is required")
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	histories, err := configService.GetConfigHistory(ctx, configID, limit, offset)
	if err != nil {
		logs.CtxErrorf(ctx, "get config history failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("get config history failed: %v", err))
		return
	}

	// 转换为响应格式
	items := make([]config.ConfigHistoryItem, 0, len(histories))
	for _, h := range histories {
		items = append(items, config.ConfigHistoryItem{
			HistoryID:     h.HistoryID,
			ConfigKey:     h.ConfigKey,
			OldValue:      h.OldValue,
			NewValue:      h.NewValue,
			ChangeReason:  h.ChangeReason,
			ChangedBy:     h.ChangedBy,
			ChangedAt:     h.ChangedAt.Format("2006-01-02 15:04:05"),
			VersionNumber: h.VersionNumber,
		})
	}

	resp := &config.GetConfigHistoryResp{
		History: items,
		Total:   len(items),
	}
	httputil.BuildSuccessResp(c, resp)
}

// RollbackConfig 回滚配置到指定版本
// @router /api/config/rollback [POST]
func RollbackConfig(ctx context.Context, c *app.RequestContext) {
	var req config.RollbackConfigReq
	if err := c.BindAndValidate(&req); err != nil {
		c.String(consts.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
		return
	}

	// 获取用户ID
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "system"
	}

	err := configService.RollbackConfig(
		ctx,
		req.TenantID,
		req.ConfigKey,
		req.Version,
		userID,
	)
	if err != nil {
		logs.CtxErrorf(ctx, "rollback config failed: %v", err)
		c.String(consts.StatusInternalServerError, fmt.Sprintf("rollback config failed: %v", err))
		return
	}

	resp := &config.RollbackConfigResp{
		Success: true,
	}
	httputil.BuildSuccessResp(c, resp)
}
