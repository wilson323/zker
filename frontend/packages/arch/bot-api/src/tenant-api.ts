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

import { axiosInstance, type BotAPIRequestConfig } from './axios';
import type {
  CreateTenantRequest,
  CreateTenantResponse,
  GetTenantRequest,
  GetTenantResponse,
  UpdateTenantRequest,
  UpdateTenantResponse,
  ListTenantsRequest,
  ListTenantsResponse,
  DeleteTenantRequest,
  DeleteTenantResponse,
  GetSubscriptionRequest,
  GetSubscriptionResponse,
  UpdateSubscriptionRequest,
  UpdateSubscriptionResponse,
  UpgradeSubscriptionRequest,
  UpgradeSubscriptionResponse,
} from './types/tenant.types';

/**
 * 租户管理 API
 * 提供租户、订阅的完整 CRUD 操作
 */
class TenantApiService {
  private readonly baseURL = '/api/v1/tenants';

  /**
   * 创建租户
   */
  async createTenant(
    request: CreateTenantRequest,
    config?: BotAPIRequestConfig,
  ): Promise<CreateTenantResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: this.baseURL,
        data: request,
      },
      config,
    );
  }

  /**
   * 获取租户详情
   */
  async getTenant(
    request: GetTenantRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetTenantResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/${request.tenant_id}`,
      },
      config,
    );
  }

  /**
   * 更新租户信息
   */
  async updateTenant(
    request: UpdateTenantRequest,
    config?: BotAPIRequestConfig,
  ): Promise<UpdateTenantResponse> {
    const { tenant_id, ...updateData } = request;
    return axiosInstance.request(
      {
        method: 'PUT',
        url: `${this.baseURL}/${tenant_id}`,
        data: updateData,
      },
      config,
    );
  }

  /**
   * 列出租户（分页）
   */
  async listTenants(
    request?: ListTenantsRequest,
    config?: BotAPIRequestConfig,
  ): Promise<ListTenantsResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: this.baseURL,
        params: request,
      },
      config,
    );
  }

  /**
   * 删除租户（软删除）
   */
  async deleteTenant(
    request: DeleteTenantRequest,
    config?: BotAPIRequestConfig,
  ): Promise<DeleteTenantResponse> {
    return axiosInstance.request(
      {
        method: 'DELETE',
        url: `${this.baseURL}/${request.tenant_id}`,
      },
      config,
    );
  }

  /**
   * 获取订阅信息
   */
  async getSubscription(
    request: GetSubscriptionRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetSubscriptionResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/${request.tenant_id}/subscription`,
      },
      config,
    );
  }

  /**
   * 更新订阅信息
   */
  async updateSubscription(
    request: UpdateSubscriptionRequest,
    config?: BotAPIRequestConfig,
  ): Promise<UpdateSubscriptionResponse> {
    const { tenant_id, ...updateData } = request;
    return axiosInstance.request(
      {
        method: 'PUT',
        url: `${this.baseURL}/${tenant_id}/subscription`,
        data: updateData,
      },
      config,
    );
  }

  /**
   * 升级订阅套餐
   */
  async upgradeSubscription(
    request: UpgradeSubscriptionRequest,
    config?: BotAPIRequestConfig,
  ): Promise<UpgradeSubscriptionResponse> {
    const { tenant_id, ...upgradeData } = request;
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/${tenant_id}/subscription/upgrade`,
        data: upgradeData,
      },
      config,
    );
  }
}

/**
 * 导出单例实例
 */
export const tenantApi = new TenantApiService();

/**
 * 导出类型
 */
export type {
  CreateTenantRequest,
  CreateTenantResponse,
  GetTenantRequest,
  GetTenantResponse,
  UpdateTenantRequest,
  UpdateTenantResponse,
  ListTenantsRequest,
  ListTenantsResponse,
  DeleteTenantRequest,
  DeleteTenantResponse,
  GetSubscriptionRequest,
  GetSubscriptionResponse,
  UpdateSubscriptionRequest,
  UpdateSubscriptionResponse,
  UpgradeSubscriptionRequest,
  UpgradeSubscriptionResponse,
};
