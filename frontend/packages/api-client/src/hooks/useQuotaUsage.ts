// frontend/packages/api-client/src/hooks/useQuotaUsage.ts

import { useQuery, UseQueryOptions } from '@tanstack/react-query';
import { httpClient } from '../http/client';
import { API_ENDPOINTS } from '../utils/apiEndpoints';
import { QuotaUsageDTO, TenantID } from '../types';

/**
 * 获取租户配额使用量
 */
export function useQuotaUsage(
  tenantId: TenantID,
  options?: Omit<UseQueryOptions<QuotaUsageDTO[]>, 'queryKey' | 'queryFn'>
) {
  return useQuery({
    queryKey: ['quotas', tenantId, 'usage'],
    queryFn: async () => {
      const response = await httpClient.get<QuotaUsageDTO[]>(
        API_ENDPOINTS.QUOTA.USAGE(tenantId)
      );
      return response.data;
    },
    enabled: !!tenantId,
    ...options,
  });
}
