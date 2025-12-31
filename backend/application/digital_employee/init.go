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

package digital_employee

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/repository"
	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/service"
	"gorm.io/gorm"
)

// ServiceComponents 数字员工服务组件
type ServiceComponents struct {
	DB *gorm.DB
}

// ApplicationService 数字员工应用服务
type ApplicationService struct {
	DomainSVC *service.Container
}

var (
	digitalEmployeeAppSvc *ApplicationService
)

// InitService 初始化数字员工服务
func InitService(ctx context.Context, components *ServiceComponents) (*ApplicationService, error) {
	// 初始化Repository（使用domain层提供的factory）
	repos := repository.NewRepositories(components.DB)

	// 初始化Domain Service Container
	domainContainer := service.NewContainer(
		repos.Profile,
		repos.Task,
		repos.Performance,
	)

	// 创建应用服务
	digitalEmployeeAppSvc = &ApplicationService{
		DomainSVC: domainContainer,
	}

	return digitalEmployeeAppSvc, nil
}

// GetService 获取数字员工应用服务
func GetService() *ApplicationService {
	return digitalEmployeeAppSvc
}

// GetDomainService 获取领域服务容器（供handler使用）
func GetDomainService() *service.Container {
	if digitalEmployeeAppSvc == nil {
		return nil
	}
	return digitalEmployeeAppSvc.DomainSVC
}
