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

package permission

import (
	"context"

	permissionservice "github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/pkg/interfaces"
)

// PermissionCheckerAdapter 权限检查器适配器
// 将 PermissionChecker 适配为 interfaces.PermissionService 接口
type PermissionCheckerAdapter struct {
	checker *permissionservice.PermissionChecker
}

// NewPermissionCheckerAdapter 创建权限检查器适配器
func NewPermissionCheckerAdapter(checker *permissionservice.PermissionChecker) *PermissionCheckerAdapter {
	return &PermissionCheckerAdapter{
		checker: checker,
	}
}

// CheckPermission 实现 interfaces.PermissionService 接口
// 将 int64 userID 转换为 string，并简化参数
func (a *PermissionCheckerAdapter) CheckPermission(ctx context.Context, userID int64, resource string, action string) bool {
	userIDStr := string(rune(userID))

	// 简化实现：这里应该根据 resource 和 action 解析出 resourceType 和 resourceID
	// 当前为基本实现
	allowed, _ := a.checker.CheckDataPermission(
		ctx,
		"", // tenantID 需要从 context 中获取
		userIDStr,
		"", // resourceType 从 resource 解析
		action,
		"", // resourceID 从 resource 解析
	)

	return allowed
}

// Ensure PermissionCheckerAdapter implements interfaces.PermissionService
var _ interfaces.PermissionService = (*PermissionCheckerAdapter)(nil)
