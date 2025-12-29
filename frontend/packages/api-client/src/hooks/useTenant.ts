// frontend/packages/api-client/src/hooks/useTenant.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { TenantDTO } from '../types';

/**
 * 获取单个租户
 */
export function useTenant(
  tenantId: string,
  options?: Omit<UseQueryOptions<TenantDTO>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['tenants', tenantId],
    queryFn: async () => {
      const response = await httpClient.get<TenantDTO>(
        API_ENDPOINTS.TENANT.GET(tenantId)
      );
      return response.data;
    },
    enabled: !!tenantId,
    ...options,
  });
}
