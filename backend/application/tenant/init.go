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

package tenant

import (
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	tenantservice "github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
)

var TenantAppSVC *TenantApplicationService

// ServiceComponents 租户应用服务组件
type ServiceComponents struct {
	DB    *gorm.DB
	Cache cache.Cmdable
}

// InitService 初始化租户应用服务
func InitService(c *ServiceComponents) (*TenantApplicationService, error) {
	// 1. 初始化领域仓储
	tenantRepo := repository.NewTenantRepository(c.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(c.DB)
	quotaRepo := repository.NewQuotaRepository(c.DB)

	// 2. 初始化领域服务
	tenantSvc := tenantservice.NewTenantManagementService(c.DB, tenantRepo, quotaRepo, subscriptionRepo)
	subscriptionSvc := tenantservice.NewSubscriptionService(subscriptionRepo, quotaRepo)
	quotaSvc := tenantservice.NewQuotaService(quotaRepo)

	// 3. 创建应用服务适配器
	tenantAdapter := NewTenantServiceAdapter(tenantSvc)
	subscriptionAdapter := NewSubscriptionServiceAdapter(subscriptionSvc)
	quotaAdapter := NewQuotaServiceAdapter(quotaSvc)

	// 4. 初始化配额监控器（无需计费服务）
	quotaMonitor := NewQuotaMonitorOptimized(quotaRepo, nil, 10)

	// 5. 初始化应用服务
	tenantAppSvc := NewTenantApplicationService(
		tenantAdapter,
		subscriptionAdapter,
		quotaAdapter,
		nil, // billingSvc - 暂时不使用
		quotaMonitor,
		nil, // roleSvc - 初始化为nil，将在外部注入
	)

	// 6. 设置全局变量
	TenantAppSVC = tenantAppSvc

	// 7. 将配额服务注册到接口层，供API层调用（避免应用层依赖API层）
	// 注意：暂时注释掉，因为接口签名不匹配，需要后续修复interfaces.QuotaService接口
	// if quotaSvc != nil {
	// 	interfaces.SetQuotaService(quotaSvc)
	// }

	return tenantAppSvc, nil
}
