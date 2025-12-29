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
  GetQuotasRequest,
  GetQuotasResponse,
  CheckQuotaRequest,
  CheckQuotaResponse,
  UpdateQuotaLimitRequest,
  UpdateQuotaLimitResponse,
  GetQuotaUsageRequest,
  GetQuotaUsageResponse,
} from './types/quota.types';

/**
 * 配额管理 API
 * 提供配额查询、检查、限额设置等功能
 */
class QuotaApiService {
  private readonly baseURL = '/api/v1/tenants';

  /**
   * 获取租户配额
   */
  async getQuotas(
    request: GetQuotasRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetQuotasResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/${request.tenant_id}/quotas`,
        params: {
          resource_types: request.resource_types?.join(','),
        },
      },
      config,
    );
  }

  /**
   * 检查配额
   * 在执行操作前检查是否有足够配额
   */
  async checkQuota(
    request: CheckQuotaRequest,
    config?: BotAPIRequestConfig,
  ): Promise<CheckQuotaResponse> {
    return axiosInstance.request(
      {
        method: 'POST',
        url: `${this.baseURL}/${request.tenant_id}/quotas/check`,
        data: request,
      },
      config,
    );
  }

  /**
   * 更新配额限制
   * 仅超级管理员可调用
   */
  async updateQuotaLimit(
    request: UpdateQuotaLimitRequest,
    config?: BotAPIRequestConfig,
  ): Promise<UpdateQuotaLimitResponse> {
    const { tenant_id, resource_type, ...updateData } = request;
    return axiosInstance.request(
      {
        method: 'PUT',
        url: `${this.baseURL}/${tenant_id}/quotas/${resource_type}`,
        data: updateData,
      },
      config,
    );
  }

  /**
   * 获取配额使用情况
   */
  async getQuotaUsage(
    request: GetQuotaUsageRequest,
    config?: BotAPIRequestConfig,
  ): Promise<GetQuotaUsageResponse> {
    return axiosInstance.request(
      {
        method: 'GET',
        url: `${this.baseURL}/${request.tenant_id}/quotas/${request.resource_type}/usage`,
        params: {
          from: request.from,
          to: request.to,
        },
      },
      config,
    );
  }

  /**
   * 批量检查配额
   * 一次性检查多个资源的配额
   */
  async batchCheckQuota(
    tenantId: string,
    checks: Array<{
      resource_type: string;
      required_count: number;
    }>,
    config?: BotAPIRequestConfig,
  ): Promise<
    Array<{
      resource_type: string;
      allowed: boolean;
      current_usage: number;
      limit: number;
      reason?: string;
    }>
  > {
    return Promise.all(
      checks.map(check =>
        this.checkQuota(
          {
            tenant_id: tenantId,
            ...check,
          },
          config,
        ),
      ),
    );
  }
}

/**
 * 导出单例实例
 */
export const quotaApi = new QuotaApiService();

/**
 * 导出类型
 */
export type {
  GetQuotasRequest,
  GetQuotasResponse,
  CheckQuotaRequest,
  CheckQuotaResponse,
  UpdateQuotaLimitRequest,
  UpdateQuotaLimitResponse,
  GetQuotaUsageRequest,
  GetQuotaUsageResponse,
};
