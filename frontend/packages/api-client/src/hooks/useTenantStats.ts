// frontend/packages/api-client/src/hooks/useTenantStats.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { TenantStatsDTO } from '../types';

/**
 * 获取租户统计
 */
export function useTenantStats(
  tenantId: string,
  options?: Omit<UseQueryOptions<TenantStatsDTO>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['tenants', tenantId, 'stats'],
    queryFn: async () => {
      const response = await httpClient.get<TenantStatsDTO>(
        API_ENDPOINTS.TENANT.STATS(tenantId)
      );
      return response.data;
    },
    enabled: !!tenantId,
    ...options,
  });
}
