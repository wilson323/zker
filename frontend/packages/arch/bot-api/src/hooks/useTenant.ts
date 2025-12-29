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

import { useCallback, useEffect, useState } from 'react';

import { tenantApi } from '@coze-arch/bot-api';
import type {
  Tenant,
  Subscription,
  TenantStatus,
  SubscriptionTier,
  CreateTenantRequest,
  UpdateTenantRequest,
  ListTenantsRequest,
} from '@coze-arch/bot-api';

/**
 * 租户管理 Hook
 * 提供租户列表、详情、创建、更新等功能
 */
export const useTenant = (tenantId?: string) => {
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  // 获取租户详情
  const fetchTenant = useCallback(async () => {
    if (!tenantId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await tenantApi.getTenant({ tenant_id: tenantId });
      setTenant(response.tenant);
      setSubscription(response.subscription);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId]);

  useEffect(() => {
    fetchTenant();
  }, [fetchTenant]);

  return {
    tenant,
    subscription,
    loading,
    error,
    refetch: fetchTenant,
  };
};

/**
 * 租户列表 Hook
 * 支持分页、筛选、搜索
 */
export const useTenants = (params?: ListTenantsRequest) => {
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchTenants = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await tenantApi.listTenants(params || {});
      setTenants(response.tenants);
      setTotal(response.total);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [params]);

  useEffect(() => {
    fetchTenants();
  }, [fetchTenants]);

  return {
    tenants,
    total,
    loading,
    error,
    refetch: fetchTenants,
  };
};

/**
 * 创建租户 Hook
 */
export const useCreateTenant = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const createTenant = useCallback(async (data: CreateTenantRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await tenantApi.createTenant(data);
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    createTenant,
    loading,
    error,
  };
};

/**
 * 更新租户 Hook
 */
export const useUpdateTenant = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const updateTenant = useCallback(async (tenantId: string, data: UpdateTenantRequest) => {
    setLoading(true);
    setError(null);

    try {
      const response = await tenantApi.updateTenant({
        tenant_id: tenantId,
        ...data,
      });
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    updateTenant,
    loading,
    error,
  };
};

/**
 * 删除租户 Hook
 */
export const useDeleteTenant = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const deleteTenant = useCallback(async (tenantId: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await tenantApi.deleteTenant({ tenant_id: tenantId });
      return response;
    } catch (err) {
      setError(err as Error);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    deleteTenant,
    loading,
    error,
  };
};

/**
 * 订阅管理 Hook
 */
export const useSubscription = (tenantId?: string) => {
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchSubscription = useCallback(async () => {
    if (!tenantId) return;

    setLoading(true);
    setError(null);

    try {
      const response = await tenantApi.getSubscription({ tenant_id: tenantId });
      setSubscription(response.subscription);
    } catch (err) {
      setError(err as Error);
    } finally {
      setLoading(false);
    }
  }, [tenantId]);

  useEffect(() => {
    fetchSubscription();
  }, [fetchSubscription]);

  return {
    subscription,
    loading,
    error,
    refetch: fetchSubscription,
  };
};

/**
 * 升级订阅 Hook
 */
export const useUpgradeSubscription = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const upgradeSubscription = useCallback(
    async (
      tenantId: string,
      targetTier: SubscriptionTier,
      discountCode?: string,
    ) => {
      setLoading(true);
      setError(null);

      try {
        const response = await tenantApi.upgradeSubscription({
          tenant_id: tenantId,
          target_plan_tier: targetTier,
          discount_code: discountCode,
        });
        return response;
      } catch (err) {
        setError(err as Error);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  return {
    upgradeSubscription,
    loading,
    error,
  };
};

/**
 * 租户统计 Hook
 * 计算租户状态分布
 */
export const useTenantStats = (tenants: Tenant[]) => {
  const stats = {
    total: tenants.length,
    active: tenants.filter(t => t.status === TenantStatus.ACTIVE).length,
    suspended: tenants.filter(t => t.status === TenantStatus.SUSPENDED).length,
    draft: tenants.filter(t => t.status === TenantStatus.DRAFT).length,
    deleted: tenants.filter(t => t.status === TenantStatus.DELETED).length,
  };

  return stats;
};
