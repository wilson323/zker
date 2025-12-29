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

	"github.com/coze-studio/backend/domain/tenant/repository"
	tenantservice "github.com/coze-studio/backend/domain/tenant/service"
	"github.com/coze-studio/backend/infra/cache"
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
	quotaUsageRepo := repository.NewQuotaUsageRepository(c.DB)
	invoiceRepo := repository.NewInvoiceRepository(c.DB)

	// 2. 初始化领域服务
	tenantSvc := tenantservice.NewTenantService(tenantRepo)
	subscriptionSvc := tenantservice.NewSubscriptionService(subscriptionRepo, tenantRepo)
	quotaSvc := tenantservice.NewQuotaService(quotaRepo)
	billingSvc := tenantservice.NewBillingService(quotaUsageRepo, invoiceRepo, quotaRepo)
	quotaMonitor := tenantservice.NewQuotaMonitorOptimized(quotaRepo, billingSvc, 10)

	// 3. 初始化应用服务
	tenantAppSvc := NewTenantApplicationService(
		tenantSvc,
		subscriptionSvc,
		quotaSvc,
		billingSvc,
		quotaMonitor,
	)

	// 4. 设置全局变量
	TenantAppSVC = tenantAppSvc

	return tenantAppSvc, nil
}
