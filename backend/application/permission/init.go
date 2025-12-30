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

	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/domain/permission/repository"
	permissionservice "github.com/coze-dev/coze-studio/backend/domain/permission/service"
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

	// 2. 初始化领域服务
	roleSVC := permissionservice.NewRoleService(roleRepo, dataPermRepo, fieldPermRepo, userRoleRepo)
	departmentSVC := permissionservice.NewDepartmentService(departmentRepo, userDepartmentRepo)
	permissionChecker := permissionservice.NewPermissionChecker(
		userRoleRepo,
		dataPermRepo,
		fieldPermRepo,
		departmentRepo,
		userDepartmentRepo,
	)

	// 3. 初始化应用服务
	permissionAppSVC := NewPermissionApplicationService(
		permissionChecker,
		roleSVC,
		departmentSVC,
	)

	// 4. 设置全局变量
	PermissionAppSVC = permissionAppSVC

	// 🔧 P0修复：初始化权限中间件
	if permissionChecker != nil {
		middleware.InitPermissionMiddleware(permissionChecker)
	}

	return permissionAppSVC, nil
}
