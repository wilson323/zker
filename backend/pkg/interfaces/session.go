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

package interfaces

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
)

// SessionService Session服务接口（解决API层→应用层依赖）
type SessionService interface {
	// ValidateSession 验证会话并返回会话数据
	ValidateSession(ctx context.Context, sessionKey string) (*entity.Session, error)
}

// PermissionService 权限服务接口（解决API层→应用层依赖）
type PermissionService interface {
	// CheckPermission 检查权限
	CheckPermission(ctx context.Context, userID int64, resource string, action string) bool
}

// QuotaService 配额服务接口（解决API层→应用层依赖）
type QuotaService interface {
	// CheckQuota 检查配额
	CheckQuota(ctx context.Context, tenantID string, quotaType string) (bool, error)
}

// 全局服务实例（由application层在初始化时注入）
var (
	GlobalSessionService    SessionService
	GlobalPermissionService PermissionService
	GlobalQuotaService      QuotaService
)

// SetSessionService 设置Session服务实例
func SetSessionService(svc SessionService) {
	GlobalSessionService = svc
}

// SetPermissionService 设置权限服务实例
func SetPermissionService(svc PermissionService) {
	GlobalPermissionService = svc
}

// SetQuotaService 设置配额服务实例
func SetQuotaService(svc QuotaService) {
	GlobalQuotaService = svc
}
