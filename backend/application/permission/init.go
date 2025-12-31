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
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	permissionservice "github.com/coze-dev/coze-studio/backend/domain/permission/service"
	"github.com/coze-dev/coze-studio/backend/pkg/interfaces"
)

var PermissionAppSVC *PermissionApplicationService

// ServiceComponents 权限应用服务组件
type ServiceComponents struct {
	DB *gorm.DB
}

// InitService 初始化权限应用服务
func InitService(c *ServiceComponents) (*PermissionApplicationService, error) {
	// 1. 初始化领域仓储
	roleRepo := repository.NewRoleRepository(c.DB)
	dataPermRepo := repository.NewDataPermissionRepository(c.DB)
	fieldPermRepo := repository.NewFieldPermissionRepository(c.DB)
	userRoleRepo := repository.NewUserRoleRepository(c.DB)
	departmentRepo := repository.NewDepartmentRepository(c.DB)
	userDepartmentRepo := repository.NewUserDepartmentRepository(c.DB)

	// 2. 初始化权限检查器（首先创建，因为 RoleService 需要它）
	permissionChecker := permissionservice.NewPermissionChecker(
		c.DB,
		roleRepo,
		dataPermRepo,
		fieldPermRepo,
		userRoleRepo,
		userDepartmentRepo,
		departmentRepo,
	)

	// 3. 初始化角色服务（现在需要 permissionChecker 参数）
	roleSVC := permissionservice.NewRoleService(
		roleRepo,
		dataPermRepo,
		fieldPermRepo,
		userRoleRepo,
		permissionChecker,
	)

	// 4. 初始化部门服务适配器
	// 注意：这是一个临时的实现，实际应该从 org 模块注入真实的部门服务
	departmentSVC := NewDefaultDepartmentAdapter()

	// 5. 初始化应用服务
	permissionAppSVC := NewPermissionApplicationService(
		permissionChecker,
		roleSVC,
		departmentSVC,
	)

	// 6. 设置全局变量
	PermissionAppSVC = permissionAppSVC

	// 7. 将权限服务注册到接口层，供API层调用（避免应用层依赖API层）
	if permissionChecker != nil {
		// 使用适配器将 PermissionChecker 转换为 interfaces.PermissionService
		permissionAdapter := NewPermissionCheckerAdapter(permissionChecker)
		interfaces.SetPermissionService(permissionAdapter)
	}

	return permissionAppSVC, nil
}
